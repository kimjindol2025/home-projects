package protocol

import (
	"sync/atomic"
	"unsafe"
)

const (
	// MagicNumber: "GRIE" in hex + version 1
	MagicNumber uint64 = 0x47524945_00000001
	HeaderSize         = 128
	DataOffset         = 128
)

// StateFlag represents the state machine state
type StateFlag int32

const (
	StateIdle    StateFlag = 0
	StateWriting StateFlag = 1
	StateReady   StateFlag = 2
	StateReading StateFlag = 3
)

// String representation for debugging
func (s StateFlag) String() string {
	switch s {
	case StateIdle:
		return "IDLE"
	case StateWriting:
		return "WRITING"
	case StateReady:
		return "READY"
	case StateReading:
		return "READING"
	default:
		return "UNKNOWN"
	}
}

// EngineHeader: 공유 메모리 제어 영역 (128바이트 = 2 × 64B Cache Line)
// Layout (tight packing, 8-byte aligned):
// [0:8]    Magic        (uint64)
// [8:12]   Version      (uint32)
// [12:16]  State        (int32, atomic)
// [16:24]  SeqNum       (uint64, atomic)
// [24:32]  Timestamp    (int64)
// [32:40]  DataLen      (uint64)
// [40:48]  WriterPID    (int64)
// [48:56]  ReaderPID    (int64)
// [56:64]  TotalWrites  (uint64, atomic)
// [64:72]  TotalReads   (uint64, atomic)
// [72:128] _pad         (56 bytes)
type EngineHeader struct {
	Magic       uint64  // [0:8]
	Version     uint32  // [8:12]
	State       int32   // [12:16] atomic
	SeqNum      uint64  // [16:24] atomic
	Timestamp   int64   // [24:32]
	DataLen     uint64  // [32:40]
	WriterPID   int64   // [40:48]
	ReaderPID   int64   // [48:56]
	TotalWrites uint64  // [56:64] atomic
	TotalReads  uint64  // [64:72] atomic
	_pad        [56]byte // [72:128]
}

// Runtime check in tests - see header_test.go TestEngineHeaderSize_Is128Bytes

// CASState atomically compares and swaps the state
func (h *EngineHeader) CASState(old, new StateFlag) bool {
	statePtr := (*int32)(unsafe.Pointer(&h.State))
	return atomic.CompareAndSwapInt32(statePtr, int32(old), int32(new))
}

// LoadState atomically loads the state
func (h *EngineHeader) LoadState() StateFlag {
	statePtr := (*int32)(unsafe.Pointer(&h.State))
	return StateFlag(atomic.LoadInt32(statePtr))
}

// StoreState atomically stores the state
func (h *EngineHeader) StoreState(s StateFlag) {
	statePtr := (*int32)(unsafe.Pointer(&h.State))
	atomic.StoreInt32(statePtr, int32(s))
}

// LoadSeqNum atomically loads the sequence number
func (h *EngineHeader) LoadSeqNum() uint64 {
	seqPtr := (*uint64)(unsafe.Pointer(&h.SeqNum))
	return atomic.LoadUint64(seqPtr)
}

// AddSeqNum atomically increments the sequence number
func (h *EngineHeader) AddSeqNum() uint64 {
	seqPtr := (*uint64)(unsafe.Pointer(&h.SeqNum))
	return atomic.AddUint64(seqPtr, 1)
}

// LoadTotalWrites atomically loads total writes counter
func (h *EngineHeader) LoadTotalWrites() uint64 {
	ptr := (*uint64)(unsafe.Pointer(&h.TotalWrites))
	return atomic.LoadUint64(ptr)
}

// AddTotalWrites atomically increments total writes
func (h *EngineHeader) AddTotalWrites() uint64 {
	ptr := (*uint64)(unsafe.Pointer(&h.TotalWrites))
	return atomic.AddUint64(ptr, 1)
}

// LoadTotalReads atomically loads total reads counter
func (h *EngineHeader) LoadTotalReads() uint64 {
	ptr := (*uint64)(unsafe.Pointer(&h.TotalReads))
	return atomic.LoadUint64(ptr)
}

// AddTotalReads atomically increments total reads
func (h *EngineHeader) AddTotalReads() uint64 {
	ptr := (*uint64)(unsafe.Pointer(&h.TotalReads))
	return atomic.AddUint64(ptr, 1)
}
