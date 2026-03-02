package backpressure

import (
	"sync/atomic"
	"time"
)

// BackpressureStrategy defines how to handle overload
type BackpressureStrategy int

const (
	// StrategyDrop: Drop low-priority jobs when overloaded
	StrategyDrop BackpressureStrategy = iota
	// StrategyBuffer: Buffer jobs in memory until capacity available
	StrategyBuffer
	// StrategyThrottle: Slow down producer when overloaded
	StrategyThrottle
	// StrategyHybrid: Combination of strategies
	StrategyHybrid
)

// BackpressureController manages backpressure handling
type BackpressureController struct {
	strategy         BackpressureStrategy
	queueSize        uint64 // atomic
	maxQueueSize     uint64
	droppedJobs      uint64 // atomic
	bufferedJobs     uint64 // atomic
	throttledRequests uint64 // atomic
	lastThrottleTime  int64  // atomic: Unix nano
	throttleWindow    time.Duration
	recoveryThreshold float64 // Queue utilization to recover (0.5 = 50%)
}

// NewBackpressureController creates a new backpressure controller
func NewBackpressureController(strategy BackpressureStrategy, maxQueueSize uint64) *BackpressureController {
	return &BackpressureController{
		strategy:          strategy,
		queueSize:         0,
		maxQueueSize:      maxQueueSize,
		throttleWindow:    100 * time.Millisecond,
		recoveryThreshold: 0.5,
	}
}

// UpdateQueueSize updates the current queue size
func (bc *BackpressureController) UpdateQueueSize(size uint64) {
	atomic.StoreUint64(&bc.queueSize, size)
}

// GetQueueSize returns the current queue size
func (bc *BackpressureController) GetQueueSize() uint64 {
	return atomic.LoadUint64(&bc.queueSize)
}

// IsOverloaded checks if the queue is overloaded
func (bc *BackpressureController) IsOverloaded() bool {
	size := bc.GetQueueSize()
	utilization := float64(size) / float64(bc.maxQueueSize)
	return utilization > 0.8 // 80% threshold
}

// ShouldDrop determines if a job should be dropped (Drop Strategy)
func (bc *BackpressureController) ShouldDrop(priority int) bool {
	if bc.strategy != StrategyDrop && bc.strategy != StrategyHybrid {
		return false
	}

	if !bc.IsOverloaded() {
		return false
	}

	// Drop low-priority jobs (priority 0-2 out of 0-10)
	if priority < 3 {
		atomic.AddUint64(&bc.droppedJobs, 1)
		return true
	}

	return false
}

// ShouldBuffer determines if a job should be buffered (Buffer Strategy)
func (bc *BackpressureController) ShouldBuffer(priority int) bool {
	if bc.strategy != StrategyBuffer && bc.strategy != StrategyHybrid {
		return false
	}

	size := bc.GetQueueSize()
	if size < bc.maxQueueSize {
		atomic.AddUint64(&bc.bufferedJobs, 1)
		return true
	}

	return false
}

// ShouldThrottle determines if producer should be throttled (Throttle Strategy)
func (bc *BackpressureController) ShouldThrottle() bool {
	if bc.strategy != StrategyThrottle && bc.strategy != StrategyHybrid {
		return false
	}

	if !bc.IsOverloaded() {
		return false
	}

	// Prevent throttling too frequently
	now := time.Now().UnixNano()
	lastTime := atomic.LoadInt64(&bc.lastThrottleTime)

	if now-lastTime < bc.throttleWindow.Nanoseconds() {
		return false
	}

	atomic.StoreInt64(&bc.lastThrottleTime, now)
	atomic.AddUint64(&bc.throttledRequests, 1)

	return true
}

// GetThrottleDuration returns how long producer should wait
func (bc *BackpressureController) GetThrottleDuration() time.Duration {
	utilization := float64(bc.GetQueueSize()) / float64(bc.maxQueueSize)

	// Scale wait time based on utilization
	// At 80% utilization: 10ms
	// At 100% utilization: 50ms
	baseDuration := 10 * time.Millisecond
	if utilization > 0.8 {
		scaledDuration := baseDuration + time.Duration(float64(baseDuration)*(utilization-0.8)/0.2)
		return scaledDuration
	}

	return 0
}

// Stats returns backpressure statistics
func (bc *BackpressureController) Stats() map[string]interface{} {
	return map[string]interface{}{
		"strategy":             bc.strategy,
		"queue_size":           bc.GetQueueSize(),
		"max_queue_size":       bc.maxQueueSize,
		"utilization":          float64(bc.GetQueueSize()) / float64(bc.maxQueueSize),
		"is_overloaded":        bc.IsOverloaded(),
		"dropped_jobs":         atomic.LoadUint64(&bc.droppedJobs),
		"buffered_jobs":        atomic.LoadUint64(&bc.bufferedJobs),
		"throttled_requests":   atomic.LoadUint64(&bc.throttledRequests),
		"throttle_window_ms":   bc.throttleWindow.Milliseconds(),
		"recovery_threshold":   bc.recoveryThreshold,
	}
}

// Reset clears all counters
func (bc *BackpressureController) Reset() {
	atomic.StoreUint64(&bc.droppedJobs, 0)
	atomic.StoreUint64(&bc.bufferedJobs, 0)
	atomic.StoreUint64(&bc.throttledRequests, 0)
	atomic.StoreInt64(&bc.lastThrottleTime, 0)
}

// String returns strategy name
func (s BackpressureStrategy) String() string {
	switch s {
	case StrategyDrop:
		return "DROP"
	case StrategyBuffer:
		return "BUFFER"
	case StrategyThrottle:
		return "THROTTLE"
	case StrategyHybrid:
		return "HYBRID"
	default:
		return "UNKNOWN"
	}
}
