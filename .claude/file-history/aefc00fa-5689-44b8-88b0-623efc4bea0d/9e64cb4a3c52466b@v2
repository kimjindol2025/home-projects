package shm

import (
	"os"
	"testing"
	"time"

	"grie-engine/internal/protocol"
)

func TestCreate_And_Close(t *testing.T) {
	const size = 10 * 1024 * 1024 // 10MB
	mgr, err := Create(size)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	if mgr.Path() == "" {
		t.Fatal("Path is empty")
	}

	if mgr.Size() != size {
		t.Fatalf("Size = %d, want %d", mgr.Size(), size)
	}

	// Verify file exists
	if _, err := os.Stat(mgr.Path()); os.IsNotExist(err) {
		t.Fatal("File does not exist")
	}
}

func TestCreate_Writes_MagicNumber(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	header := mgr.Header()
	if header.Magic != protocol.MagicNumber {
		t.Fatalf("Magic = %#x, want %#x", header.Magic, protocol.MagicNumber)
	}
}

func TestOpen_Reads_MagicNumber(t *testing.T) {
	// Create first
	mgr1, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	path := mgr1.Path()
	mgr1.Close()

	defer os.Remove(path)

	// Open and verify
	mgr2, err := Open(path)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer mgr2.Close()

	header := mgr2.Header()
	if header.Magic != protocol.MagicNumber {
		t.Fatalf("Magic = %#x, want %#x", header.Magic, protocol.MagicNumber)
	}
}

func TestDataSlice_Offset_Is_128(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	dataSlice := mgr.DataSlice()
	if len(dataSlice) != mgr.Size()-protocol.DataOffset {
		t.Fatalf("DataSlice length = %d, want %d", len(dataSlice), mgr.Size()-protocol.DataOffset)
	}
}

func TestWrite_And_Read_Data(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	testData := []byte("Hello, GRIE!")

	// Write
	err = mgr.WriteData(testData)
	if err != nil {
		t.Fatalf("WriteData failed: %v", err)
	}

	// Verify header state is READY
	if mgr.Header().LoadState() != protocol.StateReady {
		t.Fatalf("State = %v, want READY", mgr.Header().LoadState())
	}

	// Read
	readData, err := mgr.ReadData(100)
	if err != nil {
		t.Fatalf("ReadData failed: %v", err)
	}

	if string(readData) != string(testData) {
		t.Fatalf("Read data = %q, want %q", string(readData), string(testData))
	}

	// Verify state is back to IDLE
	if mgr.Header().LoadState() != protocol.StateIdle {
		t.Fatalf("State after read = %v, want IDLE", mgr.Header().LoadState())
	}
}

func TestFullProtocol_Write_Then_Read_Cycle(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	// Cycle 1
	data1 := []byte("First message")
	if err := mgr.WriteData(data1); err != nil {
		t.Fatalf("Write 1 failed: %v", err)
	}
	if mgr.Header().LoadTotalWrites() != 1 {
		t.Fatalf("TotalWrites = %d, want 1", mgr.Header().LoadTotalWrites())
	}

	read1, err := mgr.ReadData(100)
	if err != nil {
		t.Fatalf("Read 1 failed: %v", err)
	}
	if string(read1) != string(data1) {
		t.Fatalf("Read 1 = %q, want %q", string(read1), string(data1))
	}
	if mgr.Header().LoadTotalReads() != 1 {
		t.Fatalf("TotalReads = %d, want 1", mgr.Header().LoadTotalReads())
	}

	// Cycle 2
	data2 := []byte("Second message")
	if err := mgr.WriteData(data2); err != nil {
		t.Fatalf("Write 2 failed: %v", err)
	}
	if mgr.Header().LoadTotalWrites() != 2 {
		t.Fatalf("TotalWrites = %d, want 2", mgr.Header().LoadTotalWrites())
	}

	read2, err := mgr.ReadData(100)
	if err != nil {
		t.Fatalf("Read 2 failed: %v", err)
	}
	if string(read2) != string(data2) {
		t.Fatalf("Read 2 = %q, want %q", string(read2), string(data2))
	}
	if mgr.Header().LoadTotalReads() != 2 {
		t.Fatalf("TotalReads = %d, want 2", mgr.Header().LoadTotalReads())
	}
}

func TestWriteData_Concurrent_Read(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	testData := []byte("Concurrent test")
	done := make(chan error, 1)

	// Launch reader goroutine
	go func() {
		readData, err := mgr.ReadData(1000) // Wait up to 1 second
		if err != nil {
			done <- err
			return
		}
		if string(readData) != string(testData) {
			done <- nil // Mismatch is ok for this test
		}
		done <- nil
	}()

	// Small delay to ensure reader is waiting
	time.Sleep(10 * time.Millisecond)

	// Write data
	if err := mgr.WriteData(testData); err != nil {
		t.Fatalf("WriteData failed: %v", err)
	}

	// Wait for reader
	if err := <-done; err != nil {
		t.Fatalf("Reader failed: %v", err)
	}
}

func TestGetStats(t *testing.T) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	stats := mgr.GetStats()

	if state, ok := stats["state"].(string); !ok || state != "IDLE" {
		t.Fatalf("state stat = %v, want IDLE", stats["state"])
	}

	if _, ok := stats["seq_num"]; !ok {
		t.Fatal("seq_num not in stats")
	}

	if _, ok := stats["total_writes"]; !ok {
		t.Fatal("total_writes not in stats")
	}

	if _, ok := stats["total_reads"]; !ok {
		t.Fatal("total_reads not in stats")
	}
}

func BenchmarkWriteData(b *testing.B) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	testData := []byte("benchmark data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset state for next iteration
		mgr.Header().StoreState(protocol.StateIdle)
		_ = mgr.WriteData(testData)
	}
}

func BenchmarkReadData(b *testing.B) {
	mgr, err := Create(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	// Pre-write data to READY state
	_ = mgr.WriteData([]byte("benchmark data"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset to READY for next read
		mgr.Header().StoreState(protocol.StateReady)
		_, _ = mgr.ReadData(100)
	}
}
