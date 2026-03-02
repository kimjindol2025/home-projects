package dispatcher

import (
	"testing"
	"time"
)

func TestDispatcher_Create(t *testing.T) {
	d, err := NewDispatcher(4, 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(d.workers) != 4 {
		t.Fatalf("Worker count = %d, want 4", len(d.workers))
	}

	if d.ringBuffer.Capacity() != 1024 {
		t.Fatalf("Ring buffer capacity = %d, want 1024", d.ringBuffer.Capacity())
	}
}

func TestDispatcher_Submit(t *testing.T) {
	d, _ := NewDispatcher(2, 64)

	jobID, err := d.Submit([]byte("test payload"), 5)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if jobID != 1 {
		t.Fatalf("Job ID = %d, want 1", jobID)
	}

	stats := d.Stats()
	if stats["total_received"] != uint64(1) {
		t.Fatalf("Total received = %v, want 1", stats["total_received"])
	}
}

func TestDispatcher_Submit_Multiple(t *testing.T) {
	d, _ := NewDispatcher(4, 256)

	for i := 0; i < 10; i++ {
		_, err := d.Submit([]byte{byte(i)}, 5)
		if err != nil {
			t.Fatalf("Submit %d failed: %v", i, err)
		}
	}

	stats := d.Stats()
	if stats["total_received"] != uint64(10) {
		t.Fatalf("Total received = %v, want 10", stats["total_received"])
	}

	if stats["queue_size"] != uint64(10) {
		t.Fatalf("Queue size = %v, want 10", stats["queue_size"])
	}
}

func TestDispatcher_FindNextWorker(t *testing.T) {
	d, _ := NewDispatcher(4, 64)

	// All workers should be idle initially
	for i := 0; i < 4; i++ {
		idx := d.findNextWorker()
		if idx < 0 {
			t.Fatal("Should find idle worker")
		}
	}

	// Set one worker to busy
	d.workers[0].SetState("BUSY")

	// Should skip busy worker
	idx := d.findNextWorker()
	if idx == 0 {
		t.Fatal("Should skip busy worker")
	}
}

func TestDispatcher_Backpressure(t *testing.T) {
	d, err := NewDispatcher(2, 16) // Small buffer (power of 2) to trigger backpressure
	if err != nil {
		t.Fatalf("NewDispatcher failed: %v", err)
	}

	// maxQueueLen = 16 * 0.8 = 12.8 ≈ 12
	// Backpressure activates when queueSize > 12 (i.e., at 13+)
	// Submit first 13 successfully (queue size 13 triggers backpressure AFTER)
	for i := 0; i < 13; i++ {
		_, submitErr := d.Submit([]byte{byte(i)}, 5)
		if submitErr != nil {
			t.Fatalf("Submit %d should succeed, got: %v", i, submitErr)
		}
	}

	// Queue should now be backpressured
	if !d.IsBackpressured() {
		t.Fatal("Should be backpressured after 13 submissions (> 12.8 threshold)")
	}

	// 14th submission should fail due to backpressure check
	_, submitErr := d.Submit([]byte("extra"), 5)
	if submitErr == nil {
		t.Fatal("Submit 14 should fail due to backpressure")
	}

	stats := d.Stats()
	if stats["total_dropped"] != uint64(1) {
		t.Fatalf("Total dropped = %v, want 1", stats["total_dropped"])
	}
}

func TestDispatcher_Stats(t *testing.T) {
	d, _ := NewDispatcher(3, 128)

	d.Submit([]byte("job1"), 5)
	d.Submit([]byte("job2"), 3)

	stats := d.Stats()

	if stats["total_workers"] != 3 {
		t.Fatalf("Total workers = %v, want 3", stats["total_workers"])
	}

	if stats["idle_workers"] != 3 {
		t.Fatalf("Idle workers = %v, want 3", stats["idle_workers"])
	}

	if stats["queue_size"] != uint64(2) {
		t.Fatalf("Queue size = %v, want 2", stats["queue_size"])
	}

	if stats["total_received"] != uint64(2) {
		t.Fatalf("Total received = %v, want 2", stats["total_received"])
	}

	if stats["backpressured"] != false {
		t.Fatal("Should not be backpressured yet")
	}
}

func TestDispatcher_WorkerState(t *testing.T) {
	d, _ := NewDispatcher(1, 64)
	worker := d.workers[0]

	if !worker.IsIdle() {
		t.Fatal("Worker should be idle initially")
	}

	worker.SetState("BUSY")
	if worker.IsIdle() {
		t.Fatal("Worker should not be idle when set to BUSY")
	}

	worker.SetState("IDLE")
	if !worker.IsIdle() {
		t.Fatal("Worker should be idle again")
	}
}

func TestDispatcher_WorkerLastSeen(t *testing.T) {
	d, _ := NewDispatcher(1, 64)
	worker := d.workers[0]

	before := time.Now()
	worker.UpdateLastSeen()
	after := time.Now()

	lastSeen := worker.GetLastSeen()
	if lastSeen.Before(before) || lastSeen.After(after.Add(1*time.Second)) {
		t.Fatal("Last seen time is incorrect")
	}
}

func TestDispatcher_WorkerProcessed(t *testing.T) {
	d, _ := NewDispatcher(1, 64)
	worker := d.workers[0]

	count := worker.IncrementProcessed()
	if count != 1 {
		t.Fatalf("First increment = %d, want 1", count)
	}

	count = worker.IncrementProcessed()
	if count != 2 {
		t.Fatalf("Second increment = %d, want 2", count)
	}
}

func BenchmarkDispatcher_Submit(b *testing.B) {
	d, _ := NewDispatcher(8, 65536)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Submit([]byte("benchmark payload"), 5)
		if d.IsBackpressured() {
			// Reset if backpressured
			d.ringBuffer.Reset()
		}
	}
}

func BenchmarkDispatcher_FindWorker(b *testing.B) {
	d, _ := NewDispatcher(8, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.findNextWorker()
	}
}
