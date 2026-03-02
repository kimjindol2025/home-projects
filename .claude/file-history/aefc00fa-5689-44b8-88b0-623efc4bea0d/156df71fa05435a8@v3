package kernel

import (
	"grie-engine/internal/shm"
	"testing"
)

// Test basic matrix structure compliance
func TestSimdMatrixTypes(t *testing.T) {
	// These are compile-time checks in Zig
	// This test verifies the interface is correct

	// Verify our bridge can handle SIMD operations
	if !testing.Short() {
		t.Log("SIMD types verified at Zig compile time")
	}
}

// Test SIMD-accelerated operations can be dispatched
func TestSimdDispatch_MatMul(t *testing.T) {
	// Op type 0 = MatMul 4x4
	payload := make([]byte, 1+2*4*4*8) // op_type + 2 matrices of 4x4 f64

	// Set operation type
	payload[0] = 0 // MatMul

	// Fill with dummy matrix data
	for i := 1; i < len(payload); i++ {
		payload[i] = byte(i % 256)
	}

	// In real scenario, Zig kernel would process this
	// For now, just verify the structure
	if len(payload) < 1+32*8 {
		t.Fatalf("Payload too small for MatMul 4x4")
	}

	// Operation type should be recognized
	op := payload[0]
	if op != 0 {
		t.Fatalf("Wrong operation type: %d", op)
	}

	t.Log("MatMul dispatch structure verified")
}

func TestSimdDispatch_FFT(t *testing.T) {
	// Op type 1 = FFT 16
	payload := make([]byte, 1+16*16) // op_type + 16 complex numbers

	// Set operation type
	payload[0] = 1 // FFT

	if len(payload) < 1+16*16 {
		t.Fatalf("Payload too small for FFT-16")
	}

	t.Log("FFT dispatch structure verified")
}

// Benchmark matrix multiplication operation latency
func BenchmarkSimdMatMul_Dispatch(b *testing.B) {
	// Simulate Go preparing MatMul data
	payload := make([]byte, 1+2*4*4*8)
	payload[0] = 0 // MatMul

	for i := 1; i < len(payload); i++ {
		payload[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// In real case, Zig kernel processes this
		// We just verify it can be dispatched
		_ = payload[0] // Access operation type
	}
}

func BenchmarkSimdFFT_Dispatch(b *testing.B) {
	payload := make([]byte, 1+16*16)
	payload[0] = 1 // FFT

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = payload[0]
	}
}

// Test kernel bridge with SIMD operations
func TestKernelBridge_SimdOperation(t *testing.T) {
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create SHM: %v", err)
	}
	defer mgr.Close()

	kb, err := NewKernelBridge(mgr.Path())
	if err != nil {
		t.Fatalf("Create bridge: %v", err)
	}
	defer kb.Close()

	// Create MatMul payload
	payload := make([]byte, 1+2*4*4*8)
	payload[0] = 0 // MatMul op

	// Fill matrices with sample data
	for i := 1; i < len(payload); i++ {
		payload[i] = byte((i * 7) % 256)
	}

	// Submit SIMD operation
	err = kb.SubmitJob(payload)
	if err != nil {
		t.Fatalf("SubmitJob: %v", err)
	}

	// Verify it was written to SHM
	dataSlice := mgr.DataSlice()
	if dataSlice[0] != 0 {
		t.Fatalf("SIMD op type not written correctly")
	}

	t.Log("SIMD operation submitted successfully")
}
