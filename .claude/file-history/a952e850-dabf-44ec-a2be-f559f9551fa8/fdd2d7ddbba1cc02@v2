package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// ErrorType classifies errors for retry decision
type ErrorType int

const (
	ErrorTransient ErrorType = iota // Transient error (retry eligible)
	ErrorPermanent                   // Permanent error (do not retry)
)

// RetryableError wraps an error with retry classification
type RetryableError struct {
	Err  error
	Type ErrorType
}

// Error implements the error interface
func (r *RetryableError) Error() string {
	return r.Err.Error()
}

// Unwrap returns the underlying error
func (r *RetryableError) Unwrap() error {
	return r.Err
}

// RetryConfig holds retry behavior configuration
type RetryConfig struct {
	MaxAttempts  int           // Maximum retry attempts (includes first attempt)
	BaseDelay    time.Duration // Initial retry delay
	MaxDelay     time.Duration // Maximum retry delay (cap)
	Multiplier   float64       // Exponential backoff multiplier
	JitterFactor float64       // Jitter factor (0.0-1.0)
}

// DefaultRetryConfig returns standard retry configuration
// 3 attempts with exponential backoff: 10ms -> 20ms -> 40ms
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		BaseDelay:    10 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		JitterFactor: 0.1,
	}
}

// NoRetryConfig returns configuration that never retries
func NoRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 1,
	}
}

// Do executes fn with exponential backoff retry logic
// Returns the first non-transient error or context cancellation
func Do(ctx context.Context, cfg RetryConfig, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		// Check context before attempting
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if retErr, ok := err.(*RetryableError); ok {
			if retErr.Type == ErrorPermanent {
				return err
			}
		}

		// Don't retry on last attempt
		if attempt == cfg.MaxAttempts-1 {
			return lastErr
		}

		// Calculate backoff delay
		delay := calculateBackoff(attempt, cfg)

		// Wait with context awareness
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}

// calculateBackoff computes exponential backoff with jitter
func calculateBackoff(attempt int, cfg RetryConfig) time.Duration {
	// Exponential: baseDelay * (multiplier ^ attempt)
	exponent := math.Pow(cfg.Multiplier, float64(attempt))
	delay := time.Duration(float64(cfg.BaseDelay) * exponent)

	// Cap at max delay
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}

	// Add jitter: +/- (jitterFactor * delay)
	if cfg.JitterFactor > 0 {
		jitter := time.Duration(rand.Float64() * 2 * float64(delay) * cfg.JitterFactor)
		jitter -= time.Duration(float64(delay) * cfg.JitterFactor)
		delay += jitter
		if delay < 0 {
			delay = 0
		}
	}

	return delay
}

// IsTransient checks if an error is transient
func IsTransient(err error) bool {
	if retErr, ok := err.(*RetryableError); ok {
		return retErr.Type == ErrorTransient
	}
	// Default: assume transient (safe assumption for most cases)
	return true
}

// IsTimeout checks if error is a timeout error
func IsTimeout(err error) bool {
	type timeoutError interface {
		Timeout() bool
	}
	te, ok := err.(timeoutError)
	return ok && te.Timeout()
}

// AsTransient wraps error as transient
func AsTransient(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*RetryableError); ok {
		return err
	}
	return &RetryableError{Err: err, Type: ErrorTransient}
}

// AsPermanent wraps error as permanent
func AsPermanent(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*RetryableError); ok {
		return err
	}
	return &RetryableError{Err: err, Type: ErrorPermanent}
}
