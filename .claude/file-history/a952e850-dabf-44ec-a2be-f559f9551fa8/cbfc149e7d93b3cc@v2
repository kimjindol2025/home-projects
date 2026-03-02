package retry

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 10*time.Millisecond {
		t.Errorf("expected BaseDelay=10ms, got %v", cfg.BaseDelay)
	}
}

func TestNoRetryConfig(t *testing.T) {
	cfg := NoRetryConfig()
	if cfg.MaxAttempts != 1 {
		t.Errorf("expected MaxAttempts=1, got %d", cfg.MaxAttempts)
	}
}

func TestDoSuccess(t *testing.T) {
	cfg := DefaultRetryConfig()
	attempts := 0

	err := Do(context.Background(), cfg, func() error {
		attempts++
		if attempts < 2 {
			return AsTransient(fmt.Errorf("transient error"))
		}
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestDoPermanentError(t *testing.T) {
	cfg := DefaultRetryConfig()
	attempts := 0

	err := Do(context.Background(), cfg, func() error {
		attempts++
		return AsPermanent(fmt.Errorf("permanent error"))
	})

	if err == nil {
		t.Error("expected error")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt on permanent error, got %d", attempts)
	}
}

func TestDoExhausted(t *testing.T) {
	cfg := DefaultRetryConfig()
	attempts := 0

	err := Do(context.Background(), cfg, func() error {
		attempts++
		return AsTransient(fmt.Errorf("always fails"))
	})

	if err == nil {
		t.Error("expected error after exhausting retries")
	}
	if attempts != cfg.MaxAttempts {
		t.Errorf("expected %d attempts, got %d", cfg.MaxAttempts, attempts)
	}
}

func TestDoContextCancelled(t *testing.T) {
	cfg := DefaultRetryConfig()
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, cfg, func() error {
		attempts++
		return AsTransient(fmt.Errorf("always fails"))
	})

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if attempts > cfg.MaxAttempts {
		t.Errorf("expected <= %d attempts with context cancellation, got %d", cfg.MaxAttempts, attempts)
	}
}

func TestCalculateBackoff(t *testing.T) {
	cfg := RetryConfig{
		BaseDelay:    10 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		JitterFactor: 0.0, // No jitter for deterministic test
	}

	delays := make([]time.Duration, 3)
	for i := 0; i < 3; i++ {
		delays[i] = calculateBackoff(i, cfg)
	}

	// Expected: 10ms, 20ms, 40ms
	expected := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		40 * time.Millisecond,
	}

	for i, exp := range expected {
		if delays[i] != exp {
			t.Errorf("attempt %d: expected %v, got %v", i, exp, delays[i])
		}
	}
}

func TestCalculateBackoffMaxCap(t *testing.T) {
	cfg := RetryConfig{
		BaseDelay:    10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   10.0,
		JitterFactor: 0.0,
	}

	delay := calculateBackoff(5, cfg)
	if delay > cfg.MaxDelay {
		t.Errorf("expected delay <= %v, got %v", cfg.MaxDelay, delay)
	}
}

func TestIsTransient(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "transient wrapped",
			err:      AsTransient(fmt.Errorf("test")),
			expected: true,
		},
		{
			name:     "permanent wrapped",
			err:      AsPermanent(fmt.Errorf("test")),
			expected: false,
		},
		{
			name:     "unwrapped error",
			err:      fmt.Errorf("test"),
			expected: true, // Default to transient
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsTransient(tt.err) != tt.expected {
				t.Errorf("IsTransient(%v) = %v, expected %v", tt.err, IsTransient(tt.err), tt.expected)
			}
		})
	}
}
