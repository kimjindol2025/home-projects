package kernel

import (
	"testing"

	"grie-engine/internal/shm"
)

func TestKernelBridge_Create(t *testing.T) {
	// Create shared memory first
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Failed to create SHM: %v", err)
	}
	defer mgr.Close()

	// Create bridge
	kb, err := NewKernelBridge(mgr.Path())
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}
	defer kb.Close()

	if kb.shmManager == nil {
		t.Fatal("Bridge SHM manager is nil")
	}

	if kb.header == nil {
		t.Fatal("Bridge header is nil")
	}
}

func TestKernelBridge_SubmitJob(t *testing.T) {
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

	payload := []byte("Test payload for Zig kernel")

	err = kb.SubmitJob(payload)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// Verify data was written
	dataSlice := mgr.DataSlice()
	if string(dataSlice[:len(payload)]) != string(payload) {
		t.Fatalf("Payload mismatch")
	}
}

func TestKernelBridge_ExecuteSync(t *testing.T) {
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

	// Manually set up the response (simulating Zig kernel)
	// In real scenario, Zig would do this
	payload := []byte("Input data")

	err = kb.SubmitJob(payload)
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}

	// Simulate Zig processing (manually mark as IDLE after reading)
	mgr.Header().StoreState(2) // StateReady

	// In a real test, we'd wait for Zig to process
	// For now, just verify the structure is correct

	stats := kb.Stats()
	if stats["shm_path"] != mgr.Path() {
		t.Fatalf("Stats path mismatch")
	}
}

func TestKernelBridge_Stats(t *testing.T) {
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

	stats := kb.Stats()

	if stats["total_signals"] != uint64(0) {
		t.Fatalf("Initial signals should be 0")
	}

	if stats["shm_path"] == "" {
		t.Fatalf("SHM path missing from stats")
	}

	if stats["avg_latency_ns"] != uint64(0) {
		t.Fatalf("Initial latency should be 0")
	}
}

func TestKernelBridge_Close(t *testing.T) {
	mgr, err := shm.Create(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("Create SHM: %v", err)
	}
	defer mgr.Close()

	kb, err := NewKernelBridge(mgr.Path())
	if err != nil {
		t.Fatalf("Create bridge: %v", err)
	}

	err = kb.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func BenchmarkKernelBridge_SubmitJob(b *testing.B) {
	mgr, _ := shm.Create(10 * 1024 * 1024)
	defer mgr.Close()

	kb, _ := NewKernelBridge(mgr.Path())
	defer kb.Close()

	payload := []byte("benchmark payload")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Header().StoreState(0) // Reset state
		_ = kb.SubmitJob(payload)
	}
}

func BenchmarkKernelBridge_WaitResult(b *testing.B) {
	mgr, _ := shm.Create(10 * 1024 * 1024)
	defer mgr.Close()

	kb, _ := NewKernelBridge(mgr.Path())
	defer kb.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Set state to READY (simulating job completion)
		mgr.Header().StoreState(2) // StateReady

		// Try to wait (will succeed immediately)
		_, _ = kb.WaitResult(100)
	}
}

func BenchmarkKernelBridge_ExecuteSync(b *testing.B) {
	mgr, _ := shm.Create(10 * 1024 * 1024)
	defer mgr.Close()

	kb, _ := NewKernelBridge(mgr.Path())
	defer kb.Close()

	payload := []byte("sync payload")

	// Simulate async Zig processing (but don't use goroutine in benchmark)
	// Instead, manually simulate the response
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.Header().StoreState(0) // Reset to IDLE

		// Directly simulate submit + immediate response
		_ = kb.SubmitJob(payload)

		// Simulate Zig completing instantly
		mgr.Header().StoreState(0) // StateIdle
	}
}
