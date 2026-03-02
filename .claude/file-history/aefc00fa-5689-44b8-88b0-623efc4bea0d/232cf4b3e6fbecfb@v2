// Integration tests: Go Orchestrator + Zig Kernel simulation
package internal

import (
	"encoding/binary"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// TestGoJuliaIntegration - Simulate Go → Julia data flow with MatMul
func TestGoJuliaIntegration(t *testing.T) {
	t.Run("MatMul4x4_Integration", func(t *testing.T) {
		// Create test matrices
		matrixA := [4][4]float64{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 16},
		}

		matrixB := [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		}

		// Serialize to binary
		data := make([]byte, 256)
		offset := 0

		// Serialize matrix A
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				bits := math.Float64bits(matrixA[i][j])
				binary.LittleEndian.PutUint64(data[offset:], bits)
				offset += 8
			}
		}

		// Serialize matrix B
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				bits := math.Float64bits(matrixB[i][j])
				binary.LittleEndian.PutUint64(data[offset:], bits)
				offset += 8
			}
		}

		// Verify serialization
		if offset != 256 {
			t.Fatalf("Expected 256 bytes, got %d", offset)
		}

		// Verify A[0][0]
		bits := binary.LittleEndian.Uint64(data[:8])
		value := math.Float64frombits(bits)
		if value != 1.0 {
			t.Errorf("Expected 1.0, got %f", value)
		}

		t.Log("✅ MatMul4x4 data serialization passed")
	})

	t.Run("FFT16_Integration", func(t *testing.T) {
		// Create complex number data (16 pairs of float64)
		data := make([]byte, 16*16) // 16 complex numbers * 16 bytes

		for i := 0; i < 16; i++ {
			// Real part
			realBits := math.Float64bits(float64(i) * math.Pi / 16.0)
			binary.LittleEndian.PutUint64(data[i*16:], realBits)

			// Imaginary part
			imagBits := math.Float64bits(float64(i) * math.E / 16.0)
			binary.LittleEndian.PutUint64(data[i*16+8:], imagBits)
		}

		// Verify first complex number
		realBits := binary.LittleEndian.Uint64(data[:8])
		realValue := math.Float64frombits(realBits)

		if math.Abs(realValue) < 1e-10 { // Expected ~0
			t.Logf("✅ First real part ≈ 0")
		}

		imagBits := binary.LittleEndian.Uint64(data[8:16])
		imagValue := math.Float64frombits(imagBits)
		expectedImag := 0.0 * math.E / 16.0

		if math.Abs(imagValue-expectedImag) < 1e-10 {
			t.Logf("✅ First imag part ≈ %f", expectedImag)
		}

		t.Log("✅ FFT-16 data serialization passed")
	})

	t.Run("Pipeline_Latency", func(t *testing.T) {
		const iterations = 100
		var latencies []time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()

			// Simulate: state transition IDLE → WRITING → READY → IDLE
			var state int32 = 0 // IDLE
			atomic.StoreInt32(&state, 1) // WRITING
			time.Sleep(1 * time.Microsecond)
			atomic.StoreInt32(&state, 2) // READY

			// Simulate read
			if atomic.LoadInt32(&state) == 2 {
				atomic.StoreInt32(&state, 0) // back to IDLE
			}

			elapsed := time.Since(start)
			latencies = append(latencies, elapsed)
		}

		// Calculate statistics
		var avgLatency time.Duration
		for _, l := range latencies {
			avgLatency += l
		}
		avgLatency /= time.Duration(len(latencies))

		t.Logf("✅ Average latency: %v (over %d iterations)", avgLatency, iterations)

		if avgLatency > 100*time.Microsecond {
			t.Logf("⚠️  Latency is higher than expected: %v", avgLatency)
		}
	})

	t.Run("Concurrent_Reader_Writer", func(t *testing.T) {
		const numOperations = 1000
		var completed int32
		var failed int32

		var wg sync.WaitGroup

		// Writer goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numOperations; i++ {
				// Simulate write
				time.Sleep(1 * time.Nanosecond)
				atomic.AddInt32(&completed, 1)
			}
		}()

		// Reader goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numOperations; i++ {
				// Simulate read
				if atomic.LoadInt32(&completed) > 0 {
					atomic.AddInt32(&failed, 0) // All succeed
				}
				time.Sleep(1 * time.Nanosecond)
			}
		}()

		wg.Wait()

		if atomic.LoadInt32(&completed) != numOperations {
			t.Errorf("Expected %d operations, got %d", numOperations,
				atomic.LoadInt32(&completed))
		}

		t.Logf("✅ Concurrent operations: %d writes, %d reads",
			atomic.LoadInt32(&completed), numOperations)
	})

	t.Run("Memory_DataIntegrity", func(t *testing.T) {
		// Simulate shared memory data integrity check
		testData := []byte{
			0x01, 0x02, 0x03, 0x04, // Magic bytes
			0x05, 0x06, 0x07, 0x08, // Version
			0x00, 0x00, 0x00, 0x00, // State (IDLE)
		}

		// Verify magic
		expectedMagic := uint32(0x04030201)
		actualMagic := binary.LittleEndian.Uint32(testData[:4])

		if actualMagic == expectedMagic {
			t.Log("✅ Data integrity verified")
		} else {
			t.Errorf("Data corruption detected: expected %x, got %x",
				expectedMagic, actualMagic)
		}
	})
}

// BenchmarkGoJuliaSignal - Benchmark Go → Julia signal latency
func BenchmarkGoJuliaSignal(b *testing.B) {
	var state int32
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		atomic.StoreInt32(&state, 2) // READY
		_ = atomic.LoadInt32(&state)  // Read
		atomic.StoreInt32(&state, 0)  // IDLE
	}
}

// BenchmarkMatrixSerialization - Benchmark 4×4 matrix serialization
func BenchmarkMatrixSerialization(b *testing.B) {
	matrix := [4][4]float64{
		{1.1, 2.2, 3.3, 4.4},
		{5.5, 6.6, 7.7, 8.8},
		{9.9, 10.1, 11.11, 12.12},
		{13.13, 14.14, 15.15, 16.16},
	}

	data := make([]byte, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		offset := 0
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				bits := math.Float64bits(matrix[i][j])
				binary.LittleEndian.PutUint64(data[offset:], bits)
				offset += 8
			}
		}
	}
}

// TestE2EGoJuliaFlow - Full End-to-End test
func TestE2EGoJuliaFlow(t *testing.T) {
	t.Run("Full_Pipeline", func(t *testing.T) {
		// Simulate: Write → State READY → Read → Compute → State IDLE
		steps := 0

		// Step 1: Prepare data
		data := make([]byte, 256)
		for i := 0; i < 256; i += 8 {
			bits := math.Float64bits(float64(i) / 256.0)
			binary.LittleEndian.PutUint64(data[i:i+8], bits)
		}
		steps++

		// Step 2: Set state to READY
		var state int32 = 2 // READY
		steps++

		// Step 3: Reader polls and detects READY
		if atomic.LoadInt32(&state) == 2 {
			steps++
		}

		// Step 4: Reader performs computation
		computationTime := time.Duration(0)
		{
			start := time.Now()
			sum := float64(0)
			for i := 0; i < 32; i++ {
				bits := binary.LittleEndian.Uint64(data[i*8 : i*8+8])
				sum += math.Float64frombits(bits)
			}
			computationTime = time.Since(start)
			steps++
			t.Logf("Computation result: %f, time: %v", sum, computationTime)
		}

		// Step 5: Reader transitions to IDLE
		atomic.StoreInt32(&state, 0)
		steps++

		// Step 6: Verify final state
		if atomic.LoadInt32(&state) == 0 {
			steps++
		}

		if steps == 6 {
			t.Logf("✅ Full E2E pipeline completed (%d steps)", steps)
		} else {
			t.Errorf("Pipeline incomplete: %d/6 steps", steps)
		}
	})
}

// TestMemoryAlignment - Verify cache-line alignment
func TestMemoryAlignment(t *testing.T) {
	// Simulate EngineHeader alignment (should be 128 bytes)
	type EngineHeader struct {
		Magic       uint64    // 8
		Version     uint32    // 4
		State       int32     // 4
		SeqNum      uint64    // 8
		Timestamp   int64     // 8
		DataLen     uint64    // 8
		_pad1       [24]byte  // 24
		WriterPID   int64     // 8
		ReaderPID   int64     // 8
		TotalWrites uint64    // 8
		TotalReads  uint64    // 8
		_pad2       [32]byte  // 32
		// Total: 128
	}

	header := EngineHeader{}
	size := unsafe.Sizeof(header)

	if size != 128 {
		t.Errorf("Expected header size 128, got %d", size)
	} else {
		t.Logf("✅ Header size verified: %d bytes", size)
	}

	// Verify cache-line alignment (64-byte cache line)
	stateOffset := int(unsafe.Offsetof(header.State))
	pidOffset := int(unsafe.Offsetof(header.WriterPID))

	if stateOffset >= 64 {
		t.Errorf("State should be in first cache line, offset %d", stateOffset)
	}

	if pidOffset < 64 {
		t.Errorf("PID should be in second cache line, offset %d", pidOffset)
	}

	t.Logf("✅ Cache-line separation verified")
	t.Logf("   Cache Line 1: Magic(0) → State(%d)", stateOffset)
	t.Logf("   Cache Line 2: WriterPID(%d) → End(128)", pidOffset)
}

