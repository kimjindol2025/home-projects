package rpc

import (
	"io"
	"net"
)

// Transport abstracts over the byte stream.
// For tests: InProcessTransport (a pair of net.Pipe() connections).
// For real use: a net.Conn (TCP/Unix socket).
type Transport interface {
	// Read blocks until exactly len(p) bytes are available or error.
	Read(p []byte) (n int, err error)
	// Write sends all of p or returns an error.
	Write(p []byte) (n int, err error)
	// Close tears down the underlying connection.
	Close() error
}

// ConnTransport wraps a net.Conn to satisfy Transport.
type ConnTransport struct {
	conn net.Conn
}

// NewConnTransport creates a Transport from a net.Conn.
func NewConnTransport(c net.Conn) *ConnTransport {
	return &ConnTransport{conn: c}
}

// Read reads exactly len(p) bytes or returns an error.
func (t *ConnTransport) Read(p []byte) (n int, err error) {
	return io.ReadFull(t.conn, p)
}

// Write writes all of p or returns an error.
func (t *ConnTransport) Write(p []byte) (n int, err error) {
	return t.conn.Write(p)
}

// Close closes the underlying connection.
func (t *ConnTransport) Close() error {
	return t.conn.Close()
}

// NewInProcessPair returns two transports wired together via net.Pipe().
// The first is the client side, the second the server side.
func NewInProcessPair() (client Transport, server Transport, err error) {
	clientConn, serverConn := net.Pipe()
	return NewConnTransport(clientConn), NewConnTransport(serverConn), nil
}
