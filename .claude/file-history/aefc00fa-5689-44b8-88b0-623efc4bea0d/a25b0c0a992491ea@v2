package shm

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"grie-engine/internal/protocol"
)

// SharedMemManager manages zero-copy shared memory via mmap
type SharedMemManager struct {
	path     string                   // /tmp/grie_shm_XXXXX
	file     *os.File                 // file handle
	data     []byte                   // mmap'ed data slice
	header   *protocol.EngineHeader   // pointer to header in mmap
}

// Create creates and initializes a new shared memory segment
// Size should be at least 10MB (10*1024*1024)
func Create(size int) (*SharedMemManager, error) {
	if size < protocol.DataOffset {
		return nil, fmt.Errorf("size must be at least %d bytes", protocol.DataOffset)
	}

	// Generate unique temp file path
	timestamp := time.Now().UnixNano()
	path := filepath.Join("/tmp", fmt.Sprintf("grie_shm_%d", timestamp))

	// Create file with write+read permissions
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	// Truncate to desired size (fill with zeros)
	if err := file.Truncate(int64(size)); err != nil {
		return nil, fmt.Errorf("failed to truncate file: %w", err)
	}

	// mmap with MAP_SHARED flag (visible to other processes)
	data, err := syscall.Mmap(int(file.Fd()), 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %w", err)
	}

	// Get header pointer (first 128 bytes)
	headerPtr := (*protocol.EngineHeader)(unsafe.Pointer(&data[0]))

	// Initialize header
	headerPtr.Magic = protocol.MagicNumber
	headerPtr.Version = 1
	headerPtr.State = int32(protocol.StateIdle)
	headerPtr.SeqNum = 0
	headerPtr.Timestamp = time.Now().UnixNano()
	headerPtr.DataLen = 0
	headerPtr.WriterPID = int64(os.Getpid())
	headerPtr.ReaderPID = 0
	headerPtr.TotalWrites = 0
	headerPtr.TotalReads = 0

	mgr := &SharedMemManager{
		path:   path,
		file:   file,
		data:   data,
		header: headerPtr,
	}

	return mgr, nil
}

// Open opens an existing shared memory file and maps it
func Open(path string) (*SharedMemManager, error) {
	// Open in read-write mode
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	size := int(stat.Size())

	// mmap the file (MAP_SHARED for inter-process visibility)
	data, err := syscall.Mmap(int(file.Fd()), 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %w", err)
	}

	// Get header pointer
	headerPtr := (*protocol.EngineHeader)(unsafe.Pointer(&data[0]))

	// Verify magic number
	if headerPtr.Magic != protocol.MagicNumber {
		syscall.Munmap(data)
		return nil, fmt.Errorf("invalid magic number: got %#x, want %#x",
			headerPtr.Magic, protocol.MagicNumber)
	}

	mgr := &SharedMemManager{
		path:   path,
		file:   file,
		data:   data,
		header: headerPtr,
	}

	return mgr, nil
}

// Close unmaps and closes the shared memory
func (m *SharedMemManager) Close() error {
	if m.data != nil {
		if err := syscall.Munmap(m.data); err != nil {
			return fmt.Errorf("munmap failed: %w", err)
		}
		m.data = nil
		m.header = nil
	}

	if m.file != nil {
		if err := m.file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
		m.file = nil
	}

	// Optionally remove the file (caller should decide)
	// os.Remove(m.path)

	return nil
}

// Path returns the file path of the shared memory
func (m *SharedMemManager) Path() string {
	return m.path
}

// Size returns the total size of the shared memory (including header)
func (m *SharedMemManager) Size() int {
	return len(m.data)
}

// Header returns a pointer to the EngineHeader
func (m *SharedMemManager) Header() *protocol.EngineHeader {
	return m.header
}

// DataSlice returns the data portion (after header, starting at byte 128)
func (m *SharedMemManager) DataSlice() []byte {
	if len(m.data) <= protocol.DataOffset {
		return nil
	}
	return m.data[protocol.DataOffset:]
}

// WriteData writes data to the shared memory and updates the header
// This follows the state machine: IDLE -> WRITING -> READY
func (m *SharedMemManager) WriteData(data []byte) error {
	if len(data) > len(m.DataSlice()) {
		return fmt.Errorf("data too large: %d > %d", len(data), len(m.DataSlice()))
	}

	// Transition: IDLE -> WRITING
	if !m.header.CASState(protocol.StateIdle, protocol.StateWriting) {
		return fmt.Errorf("cannot acquire write lock (current state: %v)", m.header.LoadState())
	}

	// Copy data to shared memory
	copy(m.DataSlice(), data)

	// Update metadata
	m.header.DataLen = uint64(len(data))
	m.header.Timestamp = time.Now().UnixNano()
	m.header.AddSeqNum()
	m.header.AddTotalWrites()

	// Transition: WRITING -> READY
	m.header.StoreState(protocol.StateReady)

	return nil
}

// ReadData reads data from shared memory (polls until READY, then reads)
// This follows the state machine: polling READY -> READING -> IDLE
// timeout: max time to wait for READY state (0 = no wait)
func (m *SharedMemManager) ReadData(timeoutMs int) ([]byte, error) {
	start := time.Now()
	timeout := time.Duration(timeoutMs) * time.Millisecond

	// Poll for READY state
	for {
		state := m.header.LoadState()
		if state == protocol.StateReady {
			break
		}

		if timeoutMs > 0 && time.Since(start) > timeout {
			return nil, fmt.Errorf("timeout waiting for READY state (current: %v)", state)
		}

		time.Sleep(1 * time.Millisecond)
	}

	// Transition: READY -> READING (atomic CAS)
	if !m.header.CASState(protocol.StateReady, protocol.StateReading) {
		return nil, fmt.Errorf("cannot acquire read lock")
	}

	// Read data
	dataLen := m.header.DataLen
	if dataLen > uint64(len(m.DataSlice())) {
		m.header.StoreState(protocol.StateIdle)
		return nil, fmt.Errorf("invalid data length: %d", dataLen)
	}

	result := make([]byte, dataLen)
	copy(result, m.DataSlice()[:dataLen])

	// Update metadata
	m.header.AddTotalReads()

	// Transition: READING -> IDLE
	m.header.StoreState(protocol.StateIdle)

	return result, nil
}

// GetStats returns current header statistics
func (m *SharedMemManager) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"state":        m.header.LoadState().String(),
		"seq_num":      m.header.LoadSeqNum(),
		"data_len":     m.header.DataLen,
		"writer_pid":   m.header.WriterPID,
		"reader_pid":   m.header.ReaderPID,
		"total_writes": m.header.LoadTotalWrites(),
		"total_reads":  m.header.LoadTotalReads(),
	}
}
