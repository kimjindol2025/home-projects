package rpc

// ProtoVersion is the protocol version tag written into every frame header.
const ProtoVersion byte = 0x01

// MaxPayloadSize guards against malformed/malicious frames (16 MiB).
const MaxPayloadSize = 16 * 1024 * 1024

// HeaderSize is the fixed byte count of the frame header.
const HeaderSize = 20

// MessageType discriminates requests from responses and control frames.
type MessageType uint8

const (
	MsgRequest  MessageType = 0x01
	MsgResponse MessageType = 0x02
	MsgError    MessageType = 0x03
	MsgPing     MessageType = 0x04
	MsgPong     MessageType = 0x05
)

// Wire Protocol Frame Layout (20-byte fixed header + variable payload):
//
// Offset  Size  Field        Encoding           Description
// ------  ----  -----------  ----------------   ----------------------------
// 0       1     Version      uint8              Protocol version = 0x01
// 1       1     MsgType      uint8              See MessageType constants
// 2       2     Flags        uint16 big-endian  Reserved, must be 0x0000
// 4       8     MsgID        uint64 big-endian  Monotonic request ID (client-issued)
// 12      2     MethodLen    uint16 big-endian  Length of method name in payload
// 14      2     ErrorLen     uint16 big-endian  Length of error string in payload (0 on request)
// 16      4     PayloadLen   uint32 big-endian  Total payload bytes = MethodLen + ErrorLen + BodyLen
// 20      ...   Method       []byte             UTF-8 method name (MethodLen bytes)
// 20+M    ...   ErrorString  []byte             Error text (ErrorLen bytes, 0 on request)
// 20+M+E  ...   Body         []byte             Raw argument/result bytes (rest of payload)
//
// Zero-copy principle: receiver allocates single payloadBuf of PayloadLen bytes,
// then takes sub-slices as Method, Error, and Body. The headerBuf is reused
// across loop iterations outside the message reading code.
