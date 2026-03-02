package circuitbreaker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// State represents the circuit breaker state
type State int32

const (
	StateClosed   State = 0 // Normal operation
	StateOpen     State = 1 // Failing, reject requests
	StateHalfOpen State = 2 // Testing if recovered
)

var stateNames = map[State]string{
	StateClosed:   "CLOSED",
	StateOpen:     "OPEN",
	StateHalfOpen: "HALF_OPEN",
}

// String implements fmt.Stringer
func (s State) String() string {
	if name, ok := stateNames[s]; ok {
		return name
	}
	return "UNKNOWN"
}

// ErrCircuitOpen is returned when circuit is open
var ErrCircuitOpen = fmt.Errorf("circuit breaker is open")

// BreakerConfig holds circuit breaker configuration
type BreakerConfig struct {
	FailureThreshold   int           // Consecutive failures to open (default: 5)
	RecoveryTimeout    time.Duration // Time to wait before half-open (default: 30s)
	SuccessThreshold   int           // Consecutive successes to close (default: 2)
	HalfOpenMaxCalls   int           // Concurrent requests allowed in half-open (default: 1)
}

// DefaultBreakerConfig returns standard configuration
func DefaultBreakerConfig() BreakerConfig {
	return BreakerConfig{
		FailureThreshold: 5,
		RecoveryTimeout:  30 * time.Second,
		SuccessThreshold: 2,
		HalfOpenMaxCalls: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu sync.RWMutex

	state State

	// Closed state
	consecutiveFailures int32
	lastFailureTime     int64 // Unix nano

	// Open state
	openedAt int64 // Unix nano when transitioned to open

	// Half-open state
	consecutiveSuccesses int32
	halfOpenAttempts     int32

	// Configuration
	cfg BreakerConfig

	// Metrics
	totalAllowed  uint64
	totalRejected uint64
}

// NewCircuitBreaker creates a new circuit breaker with default config
func NewCircuitBreaker() *CircuitBreaker {
	return NewCircuitBreakerWithConfig(DefaultBreakerConfig())
}

// NewCircuitBreakerWithConfig creates a circuit breaker with custom config
func NewCircuitBreakerWithConfig(cfg BreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.RecoveryTimeout <= 0 {
		cfg.RecoveryTimeout = 30 * time.Second
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.HalfOpenMaxCalls <= 0 {
		cfg.HalfOpenMaxCalls = 1
	}

	return &CircuitBreaker{
		state: StateClosed,
		cfg:   cfg,
	}
}

// Allow checks if a request is allowed
// Returns ErrCircuitOpen if circuit is open
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.allowLocked()
}

func (cb *CircuitBreaker) allowLocked() error {
	now := time.Now().UnixNano()

	switch cb.state {
	case StateClosed:
		atomic.AddUint64(&cb.totalAllowed, 1)
		return nil

	case StateOpen:
		// Check if recovery timeout has elapsed
		openedAt := atomic.LoadInt64(&cb.openedAt)
		if now-openedAt >= int64(cb.cfg.RecoveryTimeout) {
			// Transition to half-open
			cb.state = StateHalfOpen
			atomic.StoreInt32(&cb.consecutiveSuccesses, 0)
			atomic.StoreInt32(&cb.halfOpenAttempts, 0)
			atomic.AddUint64(&cb.totalAllowed, 1)
			return nil
		}
		atomic.AddUint64(&cb.totalRejected, 1)
		return ErrCircuitOpen

	case StateHalfOpen:
		// Allow limited requests in half-open
		attempts := atomic.AddInt32(&cb.halfOpenAttempts, 1)
		if attempts > int32(cb.cfg.HalfOpenMaxCalls) {
			atomic.AddInt32(&cb.halfOpenAttempts, -1) // Undo increment
			atomic.AddUint64(&cb.totalRejected, 1)
			return ErrCircuitOpen
		}
		atomic.AddUint64(&cb.totalAllowed, 1)
		return nil
	}

	return ErrCircuitOpen
}

// RecordSuccess marks a request as successful
// Can transition from half-open to closed
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// Reset failure counter
		atomic.StoreInt32(&cb.consecutiveFailures, 0)

	case StateHalfOpen:
		// Count towards closing
		successes := atomic.AddInt32(&cb.consecutiveSuccesses, 1)
		if successes >= int32(cb.cfg.SuccessThreshold) {
			cb.state = StateClosed
			atomic.StoreInt32(&cb.consecutiveFailures, 0)
			atomic.StoreInt32(&cb.consecutiveSuccesses, 0)
			atomic.StoreInt32(&cb.halfOpenAttempts, 0)
		}
	}
}

// RecordFailure marks a request as failed
// Can transition from closed to open or half-open back to open
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now().UnixNano()

	switch cb.state {
	case StateClosed:
		failures := atomic.AddInt32(&cb.consecutiveFailures, 1)
		atomic.StoreInt64(&cb.lastFailureTime, now)

		if failures >= int32(cb.cfg.FailureThreshold) {
			// Transition to open
			cb.state = StateOpen
			atomic.StoreInt64(&cb.openedAt, now)
		}

	case StateHalfOpen:
		// Any failure in half-open goes back to open
		cb.state = StateOpen
		atomic.StoreInt64(&cb.openedAt, now)
		atomic.StoreInt32(&cb.consecutiveFailures, 1)
		atomic.StoreInt32(&cb.consecutiveSuccesses, 0)
		atomic.StoreInt32(&cb.halfOpenAttempts, 0)
	}
}

// State returns the current circuit breaker state
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats returns current statistics
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	stats := map[string]interface{}{
		"state":            cb.state.String(),
		"total_allowed":    atomic.LoadUint64(&cb.totalAllowed),
		"total_rejected":   atomic.LoadUint64(&cb.totalRejected),
		"consecutive_fail": atomic.LoadInt32(&cb.consecutiveFailures),
		"consecutive_succ": atomic.LoadInt32(&cb.consecutiveSuccesses),
	}

	switch cb.state {
	case StateOpen:
		openedAt := time.Unix(0, atomic.LoadInt64(&cb.openedAt))
		stats["opened_at"] = openedAt.Format(time.RFC3339Nano)
		stats["time_until_recovery"] = (cb.cfg.RecoveryTimeout - time.Since(openedAt)).String()

	case StateHalfOpen:
		stats["half_open_attempts"] = atomic.LoadInt32(&cb.halfOpenAttempts)
		stats["half_open_max_calls"] = cb.cfg.HalfOpenMaxCalls
	}

	return stats
}

// Reset clears all state (useful for testing)
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	atomic.StoreInt32(&cb.consecutiveFailures, 0)
	atomic.StoreInt32(&cb.consecutiveSuccesses, 0)
	atomic.StoreInt32(&cb.halfOpenAttempts, 0)
	atomic.StoreInt64(&cb.lastFailureTime, 0)
	atomic.StoreInt64(&cb.openedAt, 0)
	atomic.StoreUint64(&cb.totalAllowed, 0)
	atomic.StoreUint64(&cb.totalRejected, 0)
}
