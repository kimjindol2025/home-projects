package rpc

import (
	"fmt"
	"io"
	"sync"
)

// HandlerFunc is a registered RPC handler.
// args is the raw body bytes from the request (codec-encoded argument).
// It returns result bytes to write back, or an error.
type HandlerFunc func(args []byte) (result []byte, err error)

// frameBufPool reuses frame buffers across requests.
var frameBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 512)
	},
}

// Server dispatches incoming RPC frames to registered handlers.
type Server struct {
	codec    Codec
	handlers map[string]HandlerFunc
	mu       sync.RWMutex
	errors   []string
	closed   bool
}

// NewServer creates a Server using the given codec.
func NewServer(codec Codec) *Server {
	return &Server{
		codec:    codec,
		handlers: make(map[string]HandlerFunc),
		errors:   make([]string, 0),
	}
}

// Register associates a method name with a handler.
// Returns error if the name is already taken.
func (s *Server) Register(method string, fn HandlerFunc) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[method]; exists {
		return fmt.Errorf("method already registered: %s", method)
	}

	s.handlers[method] = fn
	return nil
}

// Serve reads frames from t until t is closed or an unrecoverable error occurs.
// It dispatches each request synchronously in the calling goroutine.
// For concurrent serving, the caller wraps each accepted connection in a goroutine.
func (s *Server) Serve(t Transport) error {
	headerBuf := make([]byte, HeaderSize)

	for {
		// Read 20-byte header.
		if _, err := io.ReadFull(t, headerBuf); err != nil {
			if err == io.EOF {
				return nil // normal close
			}
			return err
		}

		// Decode header and get payload sizes.
		var msg Message
		methodLen, errorLen, bodyLen, err := DecodeHeader(headerBuf, &msg)
		if err != nil {
			s.mu.Lock()
			s.errors = append(s.errors, fmt.Sprintf("DecodeHeader error: %v", err))
			s.mu.Unlock()
			continue // try next frame
		}

		// Read payload.
		payloadLen := methodLen + errorLen + bodyLen
		payloadBuf := make([]byte, payloadLen)
		if _, err := io.ReadFull(t, payloadBuf); err != nil {
			return err
		}

		// Decode body (zero-copy sub-slices).
		DecodeBody(payloadBuf, &msg, methodLen, errorLen, bodyLen)

		// Dispatch request.
		s.handleRequest(t, &msg)
	}
}

// handleRequest dispatches a single request and sends a response.
func (s *Server) handleRequest(t Transport, req *Message) {
	// Look up handler.
	s.mu.RLock()
	handler, ok := s.handlers[string(req.Method)]
	s.mu.RUnlock()

	var respMsg Message
	respMsg.ID = req.ID
	respMsg.Type = MsgResponse
	respMsg.Flags = 0

	if !ok {
		// Unknown method: send MsgError.
		respMsg.Type = MsgError
		errorMsg := "method not found"
		respMsg.Error = []byte(errorMsg)

		s.mu.Lock()
		s.errors = append(s.errors, fmt.Sprintf("unknown method: %s", string(req.Method)))
		s.mu.Unlock()
	} else {
		// Call handler.
		result, err := handler(req.Body)
		if err != nil {
			// Handler error: send MsgError.
			respMsg.Type = MsgError
			respMsg.Error = []byte(err.Error())
		} else {
			// Success: send MsgResponse.
			respMsg.Body = result
		}
	}

	// Encode and send response using buffer pool.
	frameBuf := frameBufPool.Get().([]byte)[:0]
	frameBuf = respMsg.Encode(frameBuf)
	_, _ = t.Write(frameBuf) // Ignore write errors; connection will fail on next read.
	frameBufPool.Put(frameBuf[:0])
}

// ServeAsync spawns a goroutine to serve t, returning immediately.
func (s *Server) ServeAsync(t Transport) {
	go func() {
		_ = s.Serve(t)
	}()
}

// Errors returns accumulated non-fatal errors.
func (s *Server) Errors() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.errors...)
}
