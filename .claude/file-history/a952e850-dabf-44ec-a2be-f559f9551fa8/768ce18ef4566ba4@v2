package config

import "time"

// TimeoutConfig holds all timeout configurations for GRIE
type TimeoutConfig struct {
	JobProcessTimeout time.Duration // Job processing timeout
	SHMWriteTimeout   time.Duration // Shared memory write timeout
	SHMReadTimeout    time.Duration // Shared memory read timeout
	KernelWaitTimeout time.Duration // Kernel response wait timeout
	ShutdownTimeout   time.Duration // Graceful shutdown timeout
}

// DefaultTimeoutConfig returns production-grade timeout configuration
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		JobProcessTimeout: 5 * time.Second,
		SHMWriteTimeout:   1 * time.Second,
		SHMReadTimeout:    100 * time.Millisecond,
		KernelWaitTimeout: 5 * time.Second,
		ShutdownTimeout:   30 * time.Second,
	}
}

// FastTimeoutConfig returns testing-optimized timeout configuration
func FastTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		JobProcessTimeout: 500 * time.Millisecond,
		SHMWriteTimeout:   100 * time.Millisecond,
		SHMReadTimeout:    10 * time.Millisecond,
		KernelWaitTimeout: 500 * time.Millisecond,
		ShutdownTimeout:   5 * time.Second,
	}
}
