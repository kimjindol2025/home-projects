package rpc

import (
	"encoding/binary"
	"fmt"
)

// Message is the in-memory representation of one RPC frame.
// All []byte fields are sub-slices of a single backing buffer when decoded
// (zero-copy reads). Callers must not retain these slices past the next Read.
type Message struct {
	Version byte
	Type    MessageType
	Flags   uint16
	ID      uint64
	Method  []byte // sub-slice or owned bytes (see Encode)
	Error   []byte // non-nil on MsgError frames
	Body    []byte // argument bytes (request) or result bytes (response)
}

// Encode serialises msg into dst, growing it as needed.
// Returns the sub-slice of dst containing the frame.
// Zero heap allocations if cap(dst) >= HeaderSize + len(Method) + len(Error) + len(Body).
func (m *Message) Encode(dst []byte) []byte {
	totalPayload := len(m.Method) + len(m.Error) + len(m.Body)
	needed := HeaderSize + totalPayload

	// Grow dst if necessary.
	if cap(dst) < needed {
		dst = make([]byte, needed)
	} else if len(dst) < needed {
		dst = dst[:needed]
	}

	// Encode 20-byte header.
	dst[0] = ProtoVersion
	dst[1] = byte(m.Type)
	binary.BigEndian.PutUint16(dst[2:4], m.Flags)
	binary.BigEndian.PutUint64(dst[4:12], m.ID)
	binary.BigEndian.PutUint16(dst[12:14], uint16(len(m.Method)))
	binary.BigEndian.PutUint16(dst[14:16], uint16(len(m.Error)))
	binary.BigEndian.PutUint32(dst[16:20], uint32(totalPayload))

	// Copy payload.
	offset := HeaderSize
	copy(dst[offset:], m.Method)
	offset += len(m.Method)
	copy(dst[offset:], m.Error)
	offset += len(m.Error)
	copy(dst[offset:], m.Body)

	return dst[:needed]
}

// DecodeHeader reads the 20-byte header from src into m.
// src must be exactly HeaderSize bytes.
// Returns (methodLen, errorLen, bodyLen) so the caller knows how many
// more bytes to read before calling DecodeBody.
// Returns an error if the frame is invalid.
func DecodeHeader(src []byte, m *Message) (methodLen, errorLen, bodyLen uint32, err error) {
	if len(src) < HeaderSize {
		return 0, 0, 0, ErrInvalidFrame
	}

	m.Version = src[0]
	if m.Version != ProtoVersion {
		return 0, 0, 0, ErrInvalidFrame
	}

	m.Type = MessageType(src[1])
	m.Flags = binary.BigEndian.Uint16(src[2:4])
	m.ID = binary.BigEndian.Uint64(src[4:12])

	mLen := uint32(binary.BigEndian.Uint16(src[12:14]))
	eLen := uint32(binary.BigEndian.Uint16(src[14:16]))
	payloadLen := binary.BigEndian.Uint32(src[16:20])

	if payloadLen > MaxPayloadSize {
		return 0, 0, 0, ErrPayloadTooLarge
	}

	// Validate consistency: payloadLen should equal mLen + eLen + bodyLen
	if payloadLen < mLen+eLen {
		return 0, 0, 0, ErrInvalidFrame
	}

	bodyLen = payloadLen - mLen - eLen
	return mLen, eLen, bodyLen, nil
}

// DecodeBody populates m.Method, m.Error, m.Body as zero-copy sub-slices of payload.
// payload must be exactly methodLen+errorLen+bodyLen bytes (from DecodeHeader).
func DecodeBody(payload []byte, m *Message, methodLen, errorLen, bodyLen uint32) {
	offset := 0
	m.Method = payload[offset : offset+int(methodLen)]
	offset += int(methodLen)

	if errorLen > 0 {
		m.Error = payload[offset : offset+int(errorLen)]
	} else {
		m.Error = nil
	}
	offset += int(errorLen)

	if bodyLen > 0 {
		m.Body = payload[offset : offset+int(bodyLen)]
	} else {
		m.Body = nil
	}
}

// String returns a human-readable representation of the message.
func (m *Message) String() string {
	return fmt.Sprintf("Message{ID=%d, Type=%d, Method=%q, ErrorLen=%d, BodyLen=%d}",
		m.ID, m.Type, string(m.Method), len(m.Error), len(m.Body))
}
