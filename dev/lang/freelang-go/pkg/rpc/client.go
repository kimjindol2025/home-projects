package rpc

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// pending tracks an in-flight request waiting for its response.
type pending struct {
	ch chan *Message // buffered 1
}

// Buffer pools for performance optimization.
var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 512)
		},
	}
	pendingPool = sync.Pool{
		New: func() interface{} {
			return &pending{ch: make(chan *Message, 1)}
		},
	}
)

// Client sends RPC requests over a Transport and matches responses by ID.
type Client struct {
	codec      Codec
	transport  Transport
	nextID     uint64
	inflight   map[uint64]*pending
	mu         sync.Mutex
	writeMu    sync.Mutex // serializes writes to prevent frame interleaving
	errors     []string
	done       chan struct{}
	closeOnce  sync.Once
}

// NewClient creates a Client and starts the background read loop on t.
func NewClient(codec Codec, t Transport) *Client {
	c := &Client{
		codec:     codec,
		transport: t,
		nextID:    1,
		inflight:  make(map[uint64]*pending),
		errors:    make([]string, 0),
		done:      make(chan struct{}),
	}

	// Start background read loop.
	go c.readLoop()

	return c
}

// readLoop continuously reads responses from the transport and dispatches them
// to waiting callers.
func (c *Client) readLoop() {
	defer close(c.done)

	headerBuf := make([]byte, HeaderSize)

	for {
		// Read 20-byte header.
		if _, err := io.ReadFull(c.transport, headerBuf); err != nil {
			if err != io.EOF {
				c.mu.Lock()
				c.errors = append(c.errors, fmt.Sprintf("readLoop ReadFull error: %v", err))
				c.mu.Unlock()
			}

			// Connection closed or error; unblock all pending calls.
			c.mu.Lock()
			for _, pend := range c.inflight {
				select {
				case pend.ch <- &Message{Type: 0xFF}: // sentinel for connection lost
				default:
				}
			}
			c.inflight = make(map[uint64]*pending)
			c.mu.Unlock()
			return
		}

		// Decode header.
		var msg Message
		methodLen, errorLen, bodyLen, err := DecodeHeader(headerBuf, &msg)
		if err != nil {
			c.mu.Lock()
			c.errors = append(c.errors, fmt.Sprintf("DecodeHeader error: %v", err))
			c.mu.Unlock()
			continue
		}

		// Read payload.
		payloadLen := methodLen + errorLen + bodyLen
		payloadBuf := make([]byte, payloadLen)
		if _, err := io.ReadFull(c.transport, payloadBuf); err != nil {
			c.mu.Lock()
			c.errors = append(c.errors, fmt.Sprintf("readLoop payload read error: %v", err))
			c.mu.Unlock()
			return
		}

		// Decode body.
		DecodeBody(payloadBuf, &msg, methodLen, errorLen, bodyLen)

		// Dispatch response to waiting caller.
		c.mu.Lock()
		pend, ok := c.inflight[msg.ID]
		if ok {
			delete(c.inflight, msg.ID)
		}
		c.mu.Unlock()

		if ok {
			// Send response (buffered channel, so non-blocking).
			select {
			case pend.ch <- &msg:
			default:
				// Receiver already timed out; drop the response.
			}
		}

		// If message type is MsgPing, respond with MsgPong.
		if msg.Type == MsgPing {
			pongMsg := Message{Type: MsgPong, ID: msg.ID}
			frameBuf := make([]byte, 0, HeaderSize)
			frameBuf = pongMsg.Encode(frameBuf)

			c.writeMu.Lock()
			_, _ = c.transport.Write(frameBuf)
			c.writeMu.Unlock()
		}
	}
}

// Call invokes method with the given args value, waits for the response,
// and decodes it into result (must be a pointer).
// Blocks indefinitely — use CallWithTimeout for production use.
func (c *Client) Call(method string, args interface{}, result interface{}) error {
	return c.CallWithTimeout(method, args, result, 0)
}

// CallWithTimeout is like Call but returns ErrTimeout if no response
// arrives within d. If d is 0, waits indefinitely.
func (c *Client) CallWithTimeout(method string, args interface{}, result interface{}, d time.Duration) error {
	// Allocate and encode args using buffer pool.
	argsBuf := bufPool.Get().([]byte)[:0]
	defer bufPool.Put(argsBuf[:0])
	var err error
	argsBuf, err = c.codec.Marshal(args, argsBuf)
	if err != nil {
		return fmt.Errorf("marshal args: %w", err)
	}

	// Build request message.
	c.mu.Lock()
	msgID := c.nextID
	c.nextID++
	c.mu.Unlock()

	req := Message{
		Type:   MsgRequest,
		ID:     msgID,
		Method: []byte(method),
		Body:   argsBuf,
	}

	// Register pending response handler using pool.
	pend := pendingPool.Get().(*pending)
	c.mu.Lock()
	c.inflight[msgID] = pend
	c.mu.Unlock()

	// Send request using buffer pool.
	reqBuf := bufPool.Get().([]byte)[:0]
	defer bufPool.Put(reqBuf[:0])
	reqBuf = req.Encode(reqBuf)

	c.writeMu.Lock()
	_, err = c.transport.Write(reqBuf)
	c.writeMu.Unlock()

	if err != nil {
		c.mu.Lock()
		delete(c.inflight, msgID)
		c.mu.Unlock()
		pendingPool.Put(pend)
		return fmt.Errorf("write request: %w", err)
	}

	// Wait for response with optional timeout.
	var respMsg *Message
	if d > 0 {
		// Use time.NewTimer instead of time.After to avoid memory leaks.
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case respMsg = <-pend.ch:
		case <-timer.C:
			c.mu.Lock()
			delete(c.inflight, msgID)
			c.mu.Unlock()
			pendingPool.Put(pend)
			return ErrTimeout
		}
	} else {
		respMsg = <-pend.ch
	}

	// Check for connection lost sentinel.
	if respMsg.Type == 0xFF {
		pendingPool.Put(pend)
		return ErrConnectionLost
	}

	// Handle error response.
	if respMsg.Type == MsgError {
		pendingPool.Put(pend)
		return fmt.Errorf("rpc error: %s", string(respMsg.Error))
	}

	// Unmarshal result.
	if result != nil {
		if err := c.codec.Unmarshal(respMsg.Body, result); err != nil {
			pendingPool.Put(pend)
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	pendingPool.Put(pend)
	return nil
}

// Close shuts down the underlying transport and unblocks all pending calls
// with ErrConnectionLost.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.transport.Close()
		<-c.done // wait for readLoop to finish
	})
	return err
}

// Errors returns accumulated non-fatal errors.
func (c *Client) Errors() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.errors...)
}
