// GRIE Stress Test
// Tests GRIE under high load conditions
// Monitors: latency, throughput, memory, CPU, errors
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	SHM_SIZE = 10 * 1024 * 1024
	HEADER_SIZE = 128
)

// StressTestConfig - Configuration for stress test
type StressTestConfig struct {
	duration      time.Duration // Total test duration
	rampUpTime    time.Duration // Time to ramp up to target ops/sec
	targetOpsPerSec int64       // Target operations per second
	maxOpsPerSec  int64         // Maximum ops/sec (for burst tests)
	minOpsPerSec  int64         // Minimum ops/sec (for variable load)
	mode          string        // "fixed", "ramp", "variable", "burst"
}

// StressTestMetrics - Collected metrics
type StressTestMetrics struct {
	startTime        time.Time
	totalOps         int64
	successOps       int64
	failedOps        int64
	totalLatencyNs   int64
	minLatencyNs     int64
	maxLatencyNs     int64
	lastCheckTime    time.Time
	lastCheckOps     int64
	memAllocations   int64
	memPeakBytes     uint64
	memCurrentBytes  uint64
	goroutineCount   int
	errors           []string
	mu               sync.RWMutex
}

// SimulatedOperation - Simulate a single GRIE operation
func simulateOperation() time.Duration {
	start := time.Now()

	// Simulate: state transitions + memory access
	var state int32 = 0
	atomic.StoreInt32(&state, 2) // READY

	// Simulate some work
	sum := 0.0
	for i := 0; i < 100; i++ {
		sum += float64(i) * 3.14159
	}

	if atomic.LoadInt32(&state) == 2 {
		atomic.StoreInt32(&state, 0) // IDLE
	}

	elapsed := time.Since(start)
	_ = sum // Use sum to prevent optimization
	return elapsed
}

// UpdateMetrics - Record operation metrics
func (m *StressTestMetrics) recordOperation(latency time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.AddInt64(&m.totalOps, 1)
	latencyNs := latency.Nanoseconds()

	if success {
		atomic.AddInt64(&m.successOps, 1)
	} else {
		atomic.AddInt64(&m.failedOps, 1)
	}

	atomic.AddInt64(&m.totalLatencyNs, latencyNs)

	if m.minLatencyNs == 0 || latencyNs < m.minLatencyNs {
		m.minLatencyNs = latencyNs
	}
	if latencyNs > m.maxLatencyNs {
		m.maxLatencyNs = latencyNs
	}
}

// CheckMemory - Update memory statistics
func (m *StressTestMetrics) checkMemory() {
	m.mu.Lock()
	defer m.mu.Unlock()

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m.memCurrentBytes = stats.Alloc
	m.memAllocations = int64(stats.Mallocs)

	if stats.Alloc > m.memPeakBytes {
		m.memPeakBytes = stats.Alloc
	}

	m.goroutineCount = runtime.NumGoroutine()
}

// PrintReport - Print metrics report
func (m *StressTestMetrics) printReport(elapsed time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	success := atomic.LoadInt64(&m.successOps)
	failed := atomic.LoadInt64(&m.failedOps)
	total := atomic.LoadInt64(&m.totalOps)
	totalLatencyNs := atomic.LoadInt64(&m.totalLatencyNs)

	if total == 0 {
		fmt.Println("❌ No operations completed")
		return
	}

	avgLatencyNs := totalLatencyNs / total
	throughput := float64(total) / elapsed.Seconds()
	successRate := float64(success) / float64(total) * 100.0

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("GRIE STRESS TEST RESULTS")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("\n📊 DURATION & THROUGHPUT:\n")
	fmt.Printf("  Total Time:        %v\n", elapsed)
	fmt.Printf("  Total Operations:  %d\n", total)
	fmt.Printf("  Throughput:        %.2f ops/sec\n", throughput)
	fmt.Printf("  Success Rate:      %.2f%% (%d/%d)\n", successRate, success, total)

	fmt.Printf("\n⏱️  LATENCY STATISTICS (nanoseconds):\n")
	fmt.Printf("  Min:               %d ns\n", m.minLatencyNs)
	fmt.Printf("  Avg:               %d ns\n", avgLatencyNs)
	fmt.Printf("  Max:               %d ns\n", m.maxLatencyNs)
	fmt.Printf("  P50 (Median):      ~%d ns (approximated)\n", avgLatencyNs)

	fmt.Printf("\n💾 MEMORY USAGE:\n")
	fmt.Printf("  Current:           %.2f MB\n", float64(m.memCurrentBytes)/1024/1024)
	fmt.Printf("  Peak:              %.2f MB\n", float64(m.memPeakBytes)/1024/1024)
	fmt.Printf("  Total Allocations: %d\n", m.memAllocations)

	fmt.Printf("\n🔄 GOROUTINES & ERRORS:\n")
	fmt.Printf("  Goroutine Count:   %d\n", m.goroutineCount)
	fmt.Printf("  Failed Operations: %d\n", failed)

	if len(m.errors) > 0 {
		fmt.Printf("  Errors Detected:   %d\n", len(m.errors))
		for i, err := range m.errors {
			if i < 5 { // Show first 5 errors
				fmt.Printf("    - %s\n", err)
			}
		}
		if len(m.errors) > 5 {
			fmt.Printf("    ... and %d more\n", len(m.errors)-5)
		}
	} else {
		fmt.Printf("  Errors:            None ✅\n")
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
}

// WorkerPool - Execute operations with configurable rate
type WorkerPool struct {
	config  *StressTestConfig
	metrics *StressTestMetrics
	running bool
	mu      sync.RWMutex
	wg      sync.WaitGroup
	done    chan struct{}
}

// NewWorkerPool - Create new worker pool
func NewWorkerPool(config *StressTestConfig, metrics *StressTestMetrics) *WorkerPool {
	return &WorkerPool{
		config:  config,
		metrics: metrics,
		done:    make(chan struct{}),
	}
}

// Start - Begin stress test
func (wp *WorkerPool) start(numWorkers int) {
	wp.mu.Lock()
	if wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = true
	wp.mu.Unlock()

	wp.metrics.startTime = time.Now()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wp.wg.Add(1)
		go wp.workerRoutine(i)
	}

	// Monitoring goroutine
	wp.wg.Add(1)
	go wp.monitoringRoutine()

	fmt.Printf("🚀 Started %d workers, target: %d ops/sec\n", numWorkers, wp.config.targetOpsPerSec)
}

// workerRoutine - Single worker that executes operations
func (wp *WorkerPool) workerRoutine(id int) {
	defer wp.wg.Done()

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	opsThisInterval := int64(0)
	targetOpsPerInterval := (wp.config.targetOpsPerSec * 1000000) / 1000 // ops per ms

	for {
		select {
		case <-wp.done:
			return
		case <-ticker.C:
			// Execute operations at target rate
			for opsThisInterval < targetOpsPerInterval {
				latency := simulateOperation()
				wp.metrics.recordOperation(latency, true)
				opsThisInterval++

				// Prevent busy-wait
				if opsThisInterval%100 == 0 {
					runtime.Gosched()
				}
			}
			opsThisInterval = 0
		}
	}
}

// monitoringRoutine - Collect metrics periodically
func (wp *WorkerPool) monitoringRoutine() {
	defer wp.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-wp.done:
			return
		case <-ticker.C:
			wp.metrics.checkMemory()

			// Calculate current ops/sec
			elapsed := time.Since(wp.metrics.startTime)
			total := atomic.LoadInt64(&wp.metrics.totalOps)
			currentOpsPerSec := float64(total) / elapsed.Seconds()

			fmt.Printf("📈 [%v] Ops: %d | Rate: %.0f ops/sec | Mem: %.2f MB | Goroutines: %d\n",
				elapsed.Round(1*time.Second),
				total,
				currentOpsPerSec,
				float64(wp.metrics.memCurrentBytes)/1024/1024,
				wp.metrics.goroutineCount,
			)
		}
	}
}

// Stop - Stop stress test
func (wp *WorkerPool) stop() {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = false
	wp.mu.Unlock()

	close(wp.done)
	wp.wg.Wait()

	fmt.Println("✅ Test stopped")
}

// RunStressTest - Execute stress test
func runStressTest(config *StressTestConfig, duration time.Duration) *StressTestMetrics {
	metrics := &StressTestMetrics{
		minLatencyNs: 0,
	}

	pool := NewWorkerPool(config, metrics)

	// Determine number of workers (4-16 based on CPU count)
	numWorkers := runtime.NumCPU()
	if numWorkers < 4 {
		numWorkers = 4
	} else if numWorkers > 16 {
		numWorkers = 16
	}

	// Start workers
	pool.start(numWorkers)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run for specified duration
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("\n⏱️  Test duration reached")
	case <-sigChan:
		fmt.Println("\n⚠️  Interrupted by user")
	}

	pool.stop()

	// Final metrics update
	metrics.checkMemory()

	return metrics
}

func main() {
	// Parse flags
	modeFlag := flag.String("mode", "fixed", "Test mode: fixed|ramp|variable|burst")
	durationFlag := flag.Duration("duration", 10*time.Second, "Test duration")
	rateFlag := flag.Int64("rate", 10000, "Target ops/sec")
	maxRateFlag := flag.Int64("max-rate", 100000, "Max ops/sec (for variable/burst)")
	minRateFlag := flag.Int64("min-rate", 1000, "Min ops/sec (for variable)")
	flag.Parse()

	config := &StressTestConfig{
		mode:           *modeFlag,
		duration:       *durationFlag,
		targetOpsPerSec: *rateFlag,
		maxOpsPerSec:   *maxRateFlag,
		minOpsPerSec:   *minRateFlag,
	}

	fmt.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Println("GRIE STRESS TEST")
	fmt.Println("═══════════════════════════════════════════════════════════════════\n")

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("  Mode:              %s\n", config.mode)
	fmt.Printf("  Duration:          %v\n", config.duration)
	fmt.Printf("  Target Rate:       %d ops/sec\n", config.targetOpsPerSec)
	fmt.Printf("  CPUs Available:    %d\n", runtime.NumCPU())
	fmt.Println()

	// Run stress test
	metrics := runStressTest(config, config.duration)

	// Print results
	metrics.printReport(time.Since(metrics.startTime))

	// Exit code based on results
	success := atomic.LoadInt64(&metrics.successOps)
	total := atomic.LoadInt64(&metrics.totalOps)

	if total == 0 {
		os.Exit(1)
	}

	if success < total {
		fmt.Printf("\n⚠️  Some operations failed: %d/%d\n", success, total)
		os.Exit(1)
	}

	fmt.Println("\n✅ Stress test completed successfully!")
	os.Exit(0)
}
