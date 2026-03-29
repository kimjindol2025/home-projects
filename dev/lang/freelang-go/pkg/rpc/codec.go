package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Codec encodes/decodes Go values to/from []byte.
// The framework ships one implementation (BinaryCodec); callers can
// substitute their own (e.g., JSON for debugging).
type Codec interface {
	// Marshal serialises v into dst (grow as needed). Returns sub-slice used.
	Marshal(v interface{}, dst []byte) ([]byte, error)
	// Unmarshal deserialises src into v (v must be a pointer).
	Unmarshal(src []byte, v interface{}) error
}

// BinaryCodec is a hand-rolled binary codec supporting common types:
// string, []byte, bool, int64, nil, error.
//
// Type tag byte layout:
//   0x01 string   4-byte len (BE) + UTF-8 bytes
//   0x02 []byte   4-byte len (BE) + raw bytes
//   0x03 bool     1 byte (0x00=false, 0x01=true)
//   0x04 int64    8 bytes (BE)
//   0x05 nil      0 additional bytes
//   0x06 error    4-byte len (BE) + message bytes
type BinaryCodec struct{}

// NewBinaryCodec creates a new BinaryCodec.
func NewBinaryCodec() *BinaryCodec {
	return &BinaryCodec{}
}

// Marshal serialises v into dst using the binary codec.
func (c *BinaryCodec) Marshal(v interface{}, dst []byte) ([]byte, error) {
	if v == nil {
		return append(dst, 0x05), nil
	}

	switch val := v.(type) {
	case string:
		return c.marshalString(val, dst)
	case []byte:
		return c.marshalBytes(val, dst)
	case bool:
		return c.marshalBool(val, dst)
	case int64:
		return c.marshalInt64(val, dst)
	case error:
		return c.marshalError(val, dst)
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

// Unmarshal deserialises src into v using the binary codec.
func (c *BinaryCodec) Unmarshal(src []byte, v interface{}) error {
	if len(src) == 0 {
		return fmt.Errorf("empty buffer")
	}

	tag := src[0]
	payload := src[1:]

	switch tag {
	case 0x01:
		return c.unmarshalString(payload, v)
	case 0x02:
		return c.unmarshalBytes(payload, v)
	case 0x03:
		return c.unmarshalBool(payload, v)
	case 0x04:
		return c.unmarshalInt64(payload, v)
	case 0x05:
		return c.unmarshalNil(v)
	case 0x06:
		return c.unmarshalError(payload, v)
	default:
		return fmt.Errorf("unknown type tag: 0x%02x", tag)
	}
}

func (c *BinaryCodec) marshalString(s string, dst []byte) ([]byte, error) {
	dst = append(dst, 0x01)
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(s)))
	return append(dst, []byte(s)...), nil
}

func (c *BinaryCodec) unmarshalString(payload []byte, v interface{}) error {
	if len(payload) < 4 {
		return fmt.Errorf("string payload too short")
	}
	length := binary.BigEndian.Uint32(payload[0:4])
	if uint32(len(payload)-4) < length {
		return fmt.Errorf("string length mismatch")
	}

	ptr, ok := v.(*string)
	if !ok {
		return fmt.Errorf("unmarshal string: expected *string, got %T", v)
	}
	*ptr = string(payload[4 : 4+length])
	return nil
}

func (c *BinaryCodec) marshalBytes(b []byte, dst []byte) ([]byte, error) {
	dst = append(dst, 0x02)
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(b)))
	return append(dst, b...), nil
}

func (c *BinaryCodec) unmarshalBytes(payload []byte, v interface{}) error {
	if len(payload) < 4 {
		return fmt.Errorf("bytes payload too short")
	}
	length := binary.BigEndian.Uint32(payload[0:4])
	if uint32(len(payload)-4) < length {
		return fmt.Errorf("bytes length mismatch")
	}

	ptr, ok := v.(*[]byte)
	if !ok {
		return fmt.Errorf("unmarshal bytes: expected *[]byte, got %T", v)
	}
	*ptr = payload[4 : 4+length]
	return nil
}

func (c *BinaryCodec) marshalBool(b bool, dst []byte) ([]byte, error) {
	dst = append(dst, 0x03)
	if b {
		return append(dst, 0x01), nil
	}
	return append(dst, 0x00), nil
}

func (c *BinaryCodec) unmarshalBool(payload []byte, v interface{}) error {
	if len(payload) < 1 {
		return fmt.Errorf("bool payload too short")
	}

	ptr, ok := v.(*bool)
	if !ok {
		return fmt.Errorf("unmarshal bool: expected *bool, got %T", v)
	}
	*ptr = payload[0] != 0x00
	return nil
}

func (c *BinaryCodec) marshalInt64(i int64, dst []byte) ([]byte, error) {
	dst = append(dst, 0x04)
	return binary.BigEndian.AppendUint64(dst, uint64(i)), nil
}

func (c *BinaryCodec) unmarshalInt64(payload []byte, v interface{}) error {
	if len(payload) < 8 {
		return fmt.Errorf("int64 payload too short")
	}

	ptr, ok := v.(*int64)
	if !ok {
		return fmt.Errorf("unmarshal int64: expected *int64, got %T", v)
	}
	*ptr = int64(binary.BigEndian.Uint64(payload[0:8]))
	return nil
}

func (c *BinaryCodec) unmarshalNil(v interface{}) error {
	// For nil, we expect v to be nil or an interface{} that we can't really set.
	// This is mostly for symmetry. In practice, unmarshalNil is called only
	// when we know the value is nil from the tag.
	return nil
}

func (c *BinaryCodec) marshalError(err error, dst []byte) ([]byte, error) {
	msg := err.Error()
	dst = append(dst, 0x06)
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(msg)))
	return append(dst, []byte(msg)...), nil
}

func (c *BinaryCodec) unmarshalError(payload []byte, v interface{}) error {
	if len(payload) < 4 {
		return fmt.Errorf("error payload too short")
	}
	length := binary.BigEndian.Uint32(payload[0:4])
	if uint32(len(payload)-4) < length {
		return fmt.Errorf("error length mismatch")
	}

	ptr, ok := v.(*error)
	if !ok {
		return fmt.Errorf("unmarshal error: expected *error, got %T", v)
	}
	msg := string(payload[4 : 4+length])
	*ptr = errors.New(msg)
	return nil
}
