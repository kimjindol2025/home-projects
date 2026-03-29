package rpc

import "errors"

var (
	ErrTimeout        = errors.New("rpc: call timed out")
	ErrMethodNotFound = errors.New("rpc: method not found")
	ErrServerClosed   = errors.New("rpc: server is closed")
	ErrConnectionLost = errors.New("rpc: connection lost")
	ErrInvalidFrame   = errors.New("rpc: invalid frame")
	ErrPayloadTooLarge = errors.New("rpc: payload exceeds limit")
)
