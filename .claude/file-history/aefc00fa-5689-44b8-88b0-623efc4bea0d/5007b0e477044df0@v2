package backpressure

import (
	"testing"
	"time"
)

func TestBackpressureController_Create(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	if bc.GetQueueSize() != 0 {
		t.Fatalf("Initial queue size = %d, want 0", bc.GetQueueSize())
	}

	if bc.maxQueueSize != 1000 {
		t.Fatalf("Max queue size = %d, want 1000", bc.maxQueueSize)
	}
}

func TestBackpressureController_UpdateQueueSize(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(500)
	if bc.GetQueueSize() != 500 {
		t.Fatalf("Queue size = %d, want 500", bc.GetQueueSize())
	}

	bc.UpdateQueueSize(800)
	if bc.GetQueueSize() != 800 {
		t.Fatalf("Queue size = %d, want 800", bc.GetQueueSize())
	}
}

func TestBackpressureController_IsOverloaded(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(750) // 75% - not overloaded
	if bc.IsOverloaded() {
		t.Fatal("Should not be overloaded at 75%")
	}

	bc.UpdateQueueSize(850) // 85% - overloaded
	if !bc.IsOverloaded() {
		t.Fatal("Should be overloaded at 85%")
	}
}

func TestBackpressureController_ShouldDrop_LowPriority(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(850) // Overloaded

	// Low priority job should be dropped
	if !bc.ShouldDrop(1) {
		t.Fatal("Should drop low priority job when overloaded")
	}

	// High priority job should not be dropped
	if bc.ShouldDrop(8) {
		t.Fatal("Should not drop high priority job")
	}
}

func TestBackpressureController_ShouldDrop_NotOverloaded(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(500) // Not overloaded

	// Should not drop any job when not overloaded
	if bc.ShouldDrop(1) {
		t.Fatal("Should not drop when not overloaded")
	}
}

func TestBackpressureController_ShouldBuffer(t *testing.T) {
	bc := NewBackpressureController(StrategyBuffer, 1000)

	bc.UpdateQueueSize(500)
	if !bc.ShouldBuffer(5) {
		t.Fatal("Should buffer when space available")
	}

	bc.UpdateQueueSize(1000) // Full
	if bc.ShouldBuffer(5) {
		t.Fatal("Should not buffer when full")
	}
}

func TestBackpressureController_ShouldThrottle(t *testing.T) {
	bc := NewBackpressureController(StrategyThrottle, 1000)

	bc.UpdateQueueSize(500)
	if bc.ShouldThrottle() {
		t.Fatal("Should not throttle when not overloaded")
	}

	bc.UpdateQueueSize(850)
	if !bc.ShouldThrottle() {
		t.Fatal("Should throttle when overloaded")
	}
}

func TestBackpressureController_GetThrottleDuration(t *testing.T) {
	bc := NewBackpressureController(StrategyThrottle, 1000)

	bc.UpdateQueueSize(750) // 75% - not overloaded
	duration := bc.GetThrottleDuration()
	if duration != 0 {
		t.Fatalf("Duration at 75 percent got %v, want 0", duration)
	}

	bc.UpdateQueueSize(850) // 85% - overloaded
	duration = bc.GetThrottleDuration()
	if duration <= 0 {
		t.Fatalf("Duration at 85 percent should be > 0, got %v", duration)
	}
}

func TestBackpressureController_Stats(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(850) // 85% - overloaded
	bc.ShouldDrop(1)        // Low priority, should be dropped
	bc.ShouldDrop(1)        // Low priority, should be dropped

	stats := bc.Stats()

	if stats["queue_size"] != uint64(850) {
		t.Fatalf("Stats queue_size = %v, want 850", stats["queue_size"])
	}

	if stats["dropped_jobs"] != uint64(2) {
		t.Fatalf("Stats dropped_jobs = %v, want 2", stats["dropped_jobs"])
	}

	if utilization, ok := stats["utilization"].(float64); !ok || utilization != 0.85 {
		t.Fatalf("Stats utilization want 0.85, got %v", stats["utilization"])
	}
}

func TestBackpressureController_Reset(t *testing.T) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	bc.UpdateQueueSize(850) // 85% - overloaded
	bc.ShouldDrop(1)        // Should be dropped

	if dropped := bc.Stats()["dropped_jobs"]; dropped != uint64(1) {
		t.Fatalf("Before reset: dropped_jobs = %v, want 1", dropped)
	}

	bc.Reset()

	if dropped := bc.Stats()["dropped_jobs"]; dropped != uint64(0) {
		t.Fatalf("After reset: dropped_jobs = %v, want 0", dropped)
	}
}

func TestBackpressureStrategy_String(t *testing.T) {
	tests := []struct {
		strategy BackpressureStrategy
		want     string
	}{
		{StrategyDrop, "DROP"},
		{StrategyBuffer, "BUFFER"},
		{StrategyThrottle, "THROTTLE"},
		{StrategyHybrid, "HYBRID"},
	}

	for _, test := range tests {
		if got := test.strategy.String(); got != test.want {
			t.Fatalf("String() = %s, want %s", got, test.want)
		}
	}
}

func TestBackpressureController_HybridStrategy(t *testing.T) {
	bc := NewBackpressureController(StrategyHybrid, 1000)

	bc.UpdateQueueSize(850) // Overloaded

	// Hybrid should support both drop and throttle
	if !bc.ShouldDrop(1) {
		t.Fatal("Hybrid should support drop")
	}

	if !bc.ShouldThrottle() {
		t.Fatal("Hybrid should support throttle")
	}
}

func TestBackpressureController_ThrottleWindow(t *testing.T) {
	bc := NewBackpressureController(StrategyThrottle, 1000)
	bc.throttleWindow = 100 * time.Millisecond

	bc.UpdateQueueSize(850)

	// First throttle should succeed
	if !bc.ShouldThrottle() {
		t.Fatal("First throttle should succeed")
	}

	// Second throttle within window should fail
	if bc.ShouldThrottle() {
		t.Fatal("Second throttle within window should fail")
	}

	// Wait for window to pass
	time.Sleep(110 * time.Millisecond)

	// Third throttle should succeed
	if !bc.ShouldThrottle() {
		t.Fatal("Third throttle after window should succeed")
	}
}

func BenchmarkBackpressureController_ShouldDrop(b *testing.B) {
	bc := NewBackpressureController(StrategyDrop, 1000)
	bc.UpdateQueueSize(850)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.ShouldDrop(5)
	}
}

func BenchmarkBackpressureController_IsOverloaded(b *testing.B) {
	bc := NewBackpressureController(StrategyDrop, 1000)
	bc.UpdateQueueSize(850)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.IsOverloaded()
	}
}

func BenchmarkBackpressureController_UpdateQueueSize(b *testing.B) {
	bc := NewBackpressureController(StrategyDrop, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.UpdateQueueSize(uint64(i % 1000))
	}
}
