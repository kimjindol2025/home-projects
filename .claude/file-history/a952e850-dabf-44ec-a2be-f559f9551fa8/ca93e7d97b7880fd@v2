package monitoring

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Counter is a monotonically increasing metric
type Counter struct {
	name  string
	help  string
	value uint64 // atomic
}

// NewCounter creates a new counter
func NewCounter(name, help string) *Counter {
	return &Counter{name: name, help: help}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

// Add increments the counter by n
func (c *Counter) Add(n uint64) {
	atomic.AddUint64(&c.value, n)
}

// Value returns the current value
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Gauge is a metric that can increase and decrease
type Gauge struct {
	name  string
	help  string
	value int64 // atomic
}

// NewGauge creates a new gauge
func NewGauge(name, help string) *Gauge {
	return &Gauge{name: name, help: help}
}

// Set sets the gauge to a specific value
func (g *Gauge) Set(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// Inc increments the gauge by 1
func (g *Gauge) Inc() {
	atomic.AddInt64(&g.value, 1)
}

// Dec decrements the gauge by 1
func (g *Gauge) Dec() {
	atomic.AddInt64(&g.value, -1)
}

// Add increments the gauge by n
func (g *Gauge) Add(n int64) {
	atomic.AddInt64(&g.value, n)
}

// Value returns the current value
func (g *Gauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
}

// Registry holds all metrics
type Registry struct {
	// Job metrics
	JobsReceived    *Counter
	JobsDropped     *Counter
	JobsProcessed   *Counter
	JobsFailed      *Counter
	JobsRetried     *Counter

	// Queue metrics
	QueueSize      *Gauge
	ActiveWorkers  *Gauge

	// Circuit breaker metrics
	CBAllowed       *Counter
	CBRejected      *Counter

	// Shared memory metrics
	SHMWrites       *Counter
	SHMReads        *Counter
	SHMErrors       *Counter

	// Retry metrics
	RetryAttempts   *Counter
	RetrySuccess    *Counter
	RetryExhausted  *Counter

	// Timing
	startTime time.Time
}

// NewRegistry creates a new metric registry
func NewRegistry() *Registry {
	return &Registry{
		JobsReceived:   NewCounter("grie_jobs_received_total", "Total jobs received"),
		JobsDropped:    NewCounter("grie_jobs_dropped_total", "Total jobs dropped due to backpressure"),
		JobsProcessed:  NewCounter("grie_jobs_processed_total", "Total jobs successfully processed"),
		JobsFailed:     NewCounter("grie_jobs_failed_total", "Total jobs that failed"),
		JobsRetried:    NewCounter("grie_jobs_retried_total", "Total retry attempts"),
		QueueSize:      NewGauge("grie_queue_size", "Current job queue size"),
		ActiveWorkers:  NewGauge("grie_active_workers", "Number of active workers"),
		CBAllowed:      NewCounter("grie_circuit_breaker_allowed_total", "Requests allowed by circuit breaker"),
		CBRejected:     NewCounter("grie_circuit_breaker_rejected_total", "Requests rejected by circuit breaker"),
		SHMWrites:      NewCounter("grie_shm_writes_total", "Total shared memory writes"),
		SHMReads:       NewCounter("grie_shm_reads_total", "Total shared memory reads"),
		SHMErrors:      NewCounter("grie_shm_errors_total", "Total shared memory errors"),
		RetryAttempts:  NewCounter("grie_retry_attempts_total", "Total retry attempts"),
		RetrySuccess:   NewCounter("grie_retry_success_total", "Successful retries"),
		RetryExhausted: NewCounter("grie_retry_exhausted_total", "Retries exhausted"),
		startTime:      time.Now(),
	}
}

// Uptime returns the registry uptime in seconds
func (r *Registry) Uptime() float64 {
	return time.Since(r.startTime).Seconds()
}

// WritePrometheusFormat writes metrics in Prometheus text format
func (r *Registry) WritePrometheusFormat(w io.Writer) error {
	metrics := []struct {
		name  string
		help  string
		typ   string
		value interface{}
	}{
		// Job metrics
		{"grie_jobs_received_total", "Total jobs received", "counter", r.JobsReceived.Value()},
		{"grie_jobs_dropped_total", "Total jobs dropped due to backpressure", "counter", r.JobsDropped.Value()},
		{"grie_jobs_processed_total", "Total jobs successfully processed", "counter", r.JobsProcessed.Value()},
		{"grie_jobs_failed_total", "Total jobs that failed", "counter", r.JobsFailed.Value()},
		{"grie_jobs_retried_total", "Total retry attempts", "counter", r.JobsRetried.Value()},
		// Queue metrics
		{"grie_queue_size", "Current job queue size", "gauge", r.QueueSize.Value()},
		{"grie_active_workers", "Number of active workers", "gauge", r.ActiveWorkers.Value()},
		// Circuit breaker metrics
		{"grie_circuit_breaker_allowed_total", "Requests allowed by circuit breaker", "counter", r.CBAllowed.Value()},
		{"grie_circuit_breaker_rejected_total", "Requests rejected by circuit breaker", "counter", r.CBRejected.Value()},
		// Shared memory metrics
		{"grie_shm_writes_total", "Total shared memory writes", "counter", r.SHMWrites.Value()},
		{"grie_shm_reads_total", "Total shared memory reads", "counter", r.SHMReads.Value()},
		{"grie_shm_errors_total", "Total shared memory errors", "counter", r.SHMErrors.Value()},
		// Retry metrics
		{"grie_retry_attempts_total", "Total retry attempts", "counter", r.RetryAttempts.Value()},
		{"grie_retry_success_total", "Successful retries", "counter", r.RetrySuccess.Value()},
		{"grie_retry_exhausted_total", "Retries exhausted", "counter", r.RetryExhausted.Value()},
		// System metrics
		{"grie_uptime_seconds", "Uptime in seconds", "gauge", r.Uptime()},
	}

	// Write HELP and TYPE lines
	for _, m := range metrics {
		fmt.Fprintf(w, "# HELP %s %s\n", m.name, m.help)
		fmt.Fprintf(w, "# TYPE %s %s\n", m.name, m.typ)
	}

	// Write values
	for _, m := range metrics {
		fmt.Fprintf(w, "%s %v\n", m.name, m.value)
	}

	return nil
}

// StatsJSON returns metrics in JSON-compatible format
func (r *Registry) StatsJSON() map[string]interface{} {
	return map[string]interface{}{
		"jobs": map[string]interface{}{
			"received":  r.JobsReceived.Value(),
			"dropped":   r.JobsDropped.Value(),
			"processed": r.JobsProcessed.Value(),
			"failed":    r.JobsFailed.Value(),
			"retried":   r.JobsRetried.Value(),
		},
		"queue": map[string]interface{}{
			"size":           r.QueueSize.Value(),
			"active_workers": r.ActiveWorkers.Value(),
		},
		"circuit_breaker": map[string]interface{}{
			"allowed":  r.CBAllowed.Value(),
			"rejected": r.CBRejected.Value(),
		},
		"shm": map[string]interface{}{
			"writes": r.SHMWrites.Value(),
			"reads":  r.SHMReads.Value(),
			"errors": r.SHMErrors.Value(),
		},
		"retry": map[string]interface{}{
			"attempts":  r.RetryAttempts.Value(),
			"success":   r.RetrySuccess.Value(),
			"exhausted": r.RetryExhausted.Value(),
		},
		"uptime_seconds": r.Uptime(),
	}
}

// DefaultRegistry is the global registry instance
var DefaultRegistry = NewRegistry()
