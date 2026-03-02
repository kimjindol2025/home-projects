package shm

import (
	"testing"

	"grie-engine/internal/shm"
)

// Benchmark: Zero-Copy Read Latency
// Goal: < 100ns per operation
func BenchmarkZeroCopy_Read_Latency(b *testing.B) {
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	// Pre-populate with data
	testData := []byte("benchmark data for zero-copy latency test")
	mgr.WriteData(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset to READY for next iteration
		mgr.Header().StoreState(0) // IDLE
		mgr.Header().StoreState(2) // READY

		// Measure just the read operation
		_, _ = mgr.ReadData(10)
	}
}

// Benchmark: Atomic State Polling
// Goal: < 50ns per poll
func BenchmarkAtomicState_Polling(b *testing.B) {
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.Header().LoadState()
	}
}

// Benchmark: End-to-End Throughput
// Goal: > 100MB/s
func BenchmarkEnd_To_End_Throughput(b *testing.B) {
	mgr, err := shm.Create(100 * 1024 * 1024) // 100MB for larger operations
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	// Create a reasonably sized data chunk
	dataSize := 1024 // 1KB per message
	testData := make([]byte, dataSize)
	for i := 0; i < len(testData); i++ {
		testData[i] = byte(i % 256)
	}

	b.ResetTimer()
	b.SetBytes(int64(dataSize))

	for i := 0; i < b.N; i++ {
		// Reset state
		mgr.Header().StoreState(0) // IDLE

		// Write
		mgr.WriteData(testData)

		// Read
		mgr.Header().StoreState(2) // READY
		_, _ = mgr.ReadData(10)
	}
}

// Benchmark: Atomic CAS (Compare-And-Swap)
func BenchmarkAtomicCAS_Performance(b *testing.B) {
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Header().CASState(0, 1)
		mgr.Header().CASState(1, 0)
	}
}

// Benchmark: Data Copy Throughput
func BenchmarkData_Copy_Throughput(b *testing.B) {
	mgr, err := shm.Create(100 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Create failed: %v", err)
	}
	defer mgr.Close()

	dataSize := 4096 // 4KB per write
	testData := make([]byte, dataSize)

	b.ResetTimer()
	b.SetBytes(int64(dataSize))

	for i := 0; i < b.N; i++ {
		mgr.Header().StoreState(0)
		_ = mgr.WriteData(testData)
	}
}
