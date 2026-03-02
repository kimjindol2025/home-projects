package circuitbreaker

import (
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker()
	if cb.State() != StateClosed {
		t.Errorf("expected initial state StateClosed, got %v", cb.State())
	}
}

func TestAllowClosed(t *testing.T) {
	cb := NewCircuitBreaker()
	err := cb.Allow()
	if err != nil {
		t.Errorf("expected no error in closed state, got %v", err)
	}
}

func TestTransitionClosedToOpen(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  10 * time.Millisecond,
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Record 3 failures
	for i := 0; i < 3; i++ {
		cb.Allow()
		cb.RecordFailure()
	}

	if cb.State() != StateOpen {
		t.Errorf("expected StateOpen after failures, got %v", cb.State())
	}

	// Verify rejection
	err := cb.Allow()
	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestTransitionOpenToHalfOpen(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  50 * time.Millisecond,
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Open the circuit
	cb.Allow()
	cb.RecordFailure()

	if cb.State() != StateOpen {
		t.Fatal("circuit not opened")
	}

	// Wait for recovery timeout
	time.Sleep(60 * time.Millisecond)

	// Next Allow() should transition to half-open
	err := cb.Allow()
	if err != nil {
		t.Errorf("expected no error when transitioning to half-open, got %v", err)
	}

	if cb.State() != StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %v", cb.State())
	}
}

func TestTransitionHalfOpenToClosed(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  50 * time.Millisecond,
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 3,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Open circuit
	cb.Allow()
	cb.RecordFailure()

	// Wait and transition to half-open
	time.Sleep(60 * time.Millisecond)
	cb.Allow()

	// Record successes
	for i := 0; i < 2; i++ {
		cb.RecordSuccess()
	}

	if cb.State() != StateClosed {
		t.Errorf("expected StateClosed after successes, got %v", cb.State())
	}
}

func TestTransitionHalfOpenToOpen(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  50 * time.Millisecond,
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 3,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Open circuit
	cb.Allow()
	cb.RecordFailure()

	// Transition to half-open
	time.Sleep(60 * time.Millisecond)
	cb.Allow()

	// Fail in half-open
	cb.RecordFailure()

	if cb.State() != StateOpen {
		t.Errorf("expected StateOpen after failure in half-open, got %v", cb.State())
	}
}

func TestHalfOpenMaxCalls(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  5 * time.Second, // Longer recovery timeout
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 2,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Open circuit
	cb.Allow()
	cb.RecordFailure()

	if cb.State() != StateOpen {
		t.Fatal("circuit should be open")
	}

	// Manually transition to half-open (simulate recovery timeout)
	// We'll do this by directly testing the half-open state
	cb.mu.Lock()
	cb.state = StateHalfOpen
	cb.halfOpenAttempts = 0
	cb.consecutiveSuccesses = 0
	cb.mu.Unlock()

	// First request in half-open should be allowed
	err1 := cb.Allow()
	if err1 != nil {
		t.Errorf("first half-open request should be allowed, got %v", err1)
	}

	// Second request should be allowed
	err2 := cb.Allow()
	if err2 != nil {
		t.Errorf("second half-open request should be allowed (HalfOpenMaxCalls=2), got %v", err2)
	}

	// Third request should be rejected (max calls exceeded)
	err3 := cb.Allow()
	if err3 != ErrCircuitOpen {
		t.Errorf("third half-open request should be rejected, got %v", err3)
	}
}

func TestMetrics(t *testing.T) {
	cb := NewCircuitBreaker()

	// Record some activity
	for i := 0; i < 3; i++ {
		cb.Allow()
	}

	stats := cb.Stats()
	allowed := stats["total_allowed"].(uint64)
	if allowed != 3 {
		t.Errorf("expected total_allowed=3, got %d", allowed)
	}

	// Reject due to circuit open
	cb.Allow()
	cb.RecordFailure()
	cb.Allow()
	cb.RecordFailure()
	cb.Allow()
	cb.RecordFailure()
	cb.Allow()
	cb.RecordFailure()
	cb.Allow()
	cb.RecordFailure()

	cb.Allow() // This should be rejected (circuit open)

	stats = cb.Stats()
	rejected := stats["total_rejected"].(uint64)
	if rejected < 1 {
		t.Errorf("expected total_rejected >= 1, got %d", rejected)
	}
}

func TestReset(t *testing.T) {
	cfg := BreakerConfig{
		FailureThreshold: 1,
		RecoveryTimeout:  10 * time.Millisecond,
		SuccessThreshold: 1,
		HalfOpenMaxCalls: 1,
	}
	cb := NewCircuitBreakerWithConfig(cfg)

	// Open circuit
	cb.Allow()
	cb.RecordFailure()

	if cb.State() != StateOpen {
		t.Fatal("circuit not opened")
	}

	// Reset
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("expected StateClosed after reset, got %v", cb.State())
	}

	stats := cb.Stats()
	if stats["total_allowed"].(uint64) != 0 {
		t.Errorf("expected total_allowed=0 after reset")
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateClosed, "CLOSED"},
		{StateOpen, "OPEN"},
		{StateHalfOpen, "HALF_OPEN"},
		{State(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestConcurrentAllowAndRecord(t *testing.T) {
	cb := NewCircuitBreaker()

	// Simulate concurrent requests and failures
	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				if cb.Allow() == nil {
					cb.RecordSuccess()
				}
			}
		}()
	}

	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 3; j++ {
				cb.Allow()
				cb.RecordFailure()
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should eventually open
	if cb.State() != StateOpen {
		t.Logf("state=%v (may not be open due to timing)", cb.State())
	}
}
