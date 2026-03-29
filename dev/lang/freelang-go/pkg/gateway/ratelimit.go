package gateway

import (
	"sync"
	"time"
)

// TokenBucket implements rate limiting using a token bucket with lazy refill.
// No background goroutine is needed; tokens are refilled on each Allow() call.
type TokenBucket struct {
	mu           sync.Mutex
	capacity     float64   // Maximum tokens
	tokens       float64   // Current tokens
	refillRate   float64   // Tokens per second
	lastRefillAt time.Time // Last time tokens were refilled
}

// NewTokenBucket creates a new token bucket.
func NewTokenBucket(capacity float64, refillRatePerSecond float64) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity, // Start full
		refillRate:   refillRatePerSecond,
		lastRefillAt: time.Now(),
	}
}

// Allow checks if one token is available, and if so, consumes it and returns true.
// Otherwise returns false without consuming tokens.
// Tokens are lazily refilled based on elapsed time since last Allow() call.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillAt).Seconds()
	tb.lastRefillAt = now

	// Refill tokens based on elapsed time
	refilled := elapsed * tb.refillRate
	tb.tokens = min(tb.capacity, tb.tokens+refilled)

	// Try to consume one token
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// Remaining returns the current token count.
func (tb *TokenBucket) Remaining() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillAt).Seconds()
	refilled := elapsed * tb.refillRate
	return min(tb.capacity, tb.tokens+refilled)
}

// Reset resets the bucket to full capacity.
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = tb.capacity
	tb.lastRefillAt = time.Now()
}

// min returns the smaller of a and b.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
