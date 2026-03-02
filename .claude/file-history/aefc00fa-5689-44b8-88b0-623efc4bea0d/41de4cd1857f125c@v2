// GRIE Memory Leak Analyzer
// Detects memory leaks using Go's runtime.MemStats
// Monitors: heap size, allocations, GC frequency, goroutine count
package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// MemorySnapshot - Record memory state at a point in time
type MemorySnapshot struct {
	time           time.Time
	heapAlloc      uint64
	heapSys        uint64
	heapIdle       uint64
	heapInuse      uint64
	heapReleased   uint64
	mallocs        uint64
	frees          uint64
	liveObjects    uint64
	gcRuns         uint32
	goroutines     int
	pauseNs        uint64
}

// MemoryAnalyzer - Track memory over time
type MemoryAnalyzer struct {
	snapshots  []*MemorySnapshot
	startTime  time.Time
	mu         sync.RWMutex
	opsCount   int64
}

// Snapshot - Take a memory snapshot
func (ma *MemoryAnalyzer) snapshot() *MemorySnapshot {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return &MemorySnapshot{
		time:           time.Now(),
		heapAlloc:      stats.HeapAlloc,
		heapSys:        stats.HeapSys,
		heapIdle:       stats.HeapIdle,
		heapInuse:      stats.HeapInuse,
		heapReleased:   stats.HeapReleased,
		mallocs:        stats.Mallocs,
		frees:          stats.Frees,
		liveObjects:    stats.Mallocs - stats.Frees,
		gcRuns:         stats.NumGC,
		goroutines:     runtime.NumGoroutine(),
		pauseNs:        stats.PauseNs[(stats.NumGC+255)%256],
	}
}

// Record - Record current memory state
func (ma *MemoryAnalyzer) record() *MemorySnapshot {
	snap := ma.snapshot()

	ma.mu.Lock()
	ma.snapshots = append(ma.snapshots, snap)
	ma.mu.Unlock()

	return snap
}

// SimulateWorkload - Execute operations that might leak memory
func (ma *MemoryAnalyzer) simulateWorkload(durationSeconds int, opsPerSecond int64) {
	ma.startTime = time.Now()
	deadline := ma.startTime.Add(time.Duration(durationSeconds) * time.Second)

	operationCount := int64(0)
	tickerInterval := time.Duration(1000000000 / opsPerSecond) * time.Nanosecond

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		select {
		case <-ticker.C:
			// Simulate operation: create/destroy objects
			_ = make([]byte, 1024) // Allocate 1KB
			atomic.AddInt64(&ma.opsCount, 1)
			operationCount++

			// Periodic GC hint (not forcing, just hint)
			if operationCount%100000 == 0 {
				runtime.GC()
			}
		}
	}

	atomic.StoreInt64(&ma.opsCount, operationCount)
}

// AnalyzeLeaks - Detect potential memory leaks
func (ma *MemoryAnalyzer) analyzeLeaks() {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	if len(ma.snapshots) < 2 {
		fmt.Println("❌ Not enough snapshots for analysis")
		return
	}

	first := ma.snapshots[0]
	last := ma.snapshots[len(ma.snapshots)-1]

	heapGrowth := int64(last.heapAlloc) - int64(first.heapAlloc)
	heapGrowthRate := float64(heapGrowth) / float64(first.heapAlloc) * 100.0

	allocDiff := last.mallocs - first.mallocs
	freeDiff := last.frees - first.frees
	netAllocs := allocDiff - freeDiff

	gcGrowth := int32(last.gcRuns) - int32(first.gcRuns)
	goroutineGrowth := last.goroutines - first.goroutines

	elapsed := last.time.Sub(first.time)

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("MEMORY LEAK ANALYSIS")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("\n📊 HEAP ANALYSIS:\n")
	fmt.Printf("  Initial Heap:      %v\n", formatBytes(first.heapAlloc))
	fmt.Printf("  Final Heap:        %v\n", formatBytes(last.heapAlloc))
	fmt.Printf("  Growth:            %v (%+.2f%%)\n", formatBytesDelta(heapGrowth), heapGrowthRate)
	fmt.Printf("  Per Operation:     %.2f bytes/op\n", float64(heapGrowth)/float64(atomic.LoadInt64(&ma.opsCount)))

	fmt.Printf("\n💾 ALLOCATION ANALYSIS:\n")
	fmt.Printf("  Total Allocations: %d\n", allocDiff)
	fmt.Printf("  Total Frees:       %d\n", freeDiff)
	fmt.Printf("  Net Allocations:   %d (%+.2f%%)\n", netAllocs,
		float64(netAllocs)/float64(freeDiff)*100.0)
	fmt.Printf("  Alloc/Free Ratio:  %.4f (ideal: 1.0)\n", float64(allocDiff)/float64(freeDiff))

	fmt.Printf("\n🔄 GC & RUNTIME:\n")
	fmt.Printf("  GC Runs:           %d\n", gcGrowth)
	fmt.Printf("  GC Frequency:      %.2f runs/second\n", float64(gcGrowth)/elapsed.Seconds())
	fmt.Printf("  Goroutines:        %d → %d (%+d)\n", first.goroutines, last.goroutines, goroutineGrowth)

	fmt.Printf("\n🔍 LEAK DETECTION:\n")

	// Heuristic 1: Heap growth rate
	if heapGrowthRate > 50.0 {
		fmt.Printf("  ⚠️  HIGH heap growth: %.2f%% (threshold: 50%%)\n", heapGrowthRate)
	} else {
		fmt.Printf("  ✅ Heap growth normal: %.2f%%\n", heapGrowthRate)
	}

	// Heuristic 2: Allocation/Free balance
	freeRatio := float64(freeDiff) / float64(allocDiff)
	if freeRatio < 0.95 {
		fmt.Printf("  ⚠️  UNBALANCED alloc/free: %.2f%% freed (threshold: 95%%)\n", freeRatio*100.0)
	} else {
		fmt.Printf("  ✅ Balanced alloc/free: %.2f%% freed\n", freeRatio*100.0)
	}

	// Heuristic 3: Goroutine leak
	if goroutineGrowth > 10 {
		fmt.Printf("  ⚠️  GOROUTINE LEAK: +%d goroutines (threshold: 10)\n", goroutineGrowth)
	} else {
		fmt.Printf("  ✅ Goroutine count stable: %+d\n", goroutineGrowth)
	}

	// Heuristic 4: GC frequency
	gcPerSecond := float64(gcGrowth) / elapsed.Seconds()
	if gcPerSecond > 10.0 {
		fmt.Printf("  ⚠️  EXCESSIVE GC: %.1f runs/sec (threshold: 10)\n", gcPerSecond)
	} else {
		fmt.Printf("  ✅ GC frequency normal: %.2f runs/sec\n", gcPerSecond)
	}

	// Overall assessment
	fmt.Printf("\n" + strings.Repeat("-", 70))
	leakScore := calculateLeakScore(heapGrowthRate, freeRatio, goroutineGrowth, gcPerSecond)
	assessmentRisk := assessRisk(leakScore)

	fmt.Printf("LEAK RISK SCORE: %.2f/100 → %s\n", leakScore, assessmentRisk)
	fmt.Println(strings.Repeat("=", 70))
}

// calculateLeakScore - Calculate memory leak risk score (0-100)
func calculateLeakScore(heapGrowthRate float64, freeRatio float64, goroutineGrowth int, gcPerSecond float64) float64 {
	score := 0.0

	// Heap growth factor (0-40 points)
	if heapGrowthRate > 0 {
		heapScore := math.Min(40.0, (heapGrowthRate / 100.0) * 40.0)
		score += heapScore
	}

	// Alloc/Free balance (0-30 points)
	if freeRatio < 1.0 {
		freeScore := (1.0 - freeRatio) * 30.0 * 100.0
		score += math.Min(30.0, freeScore)
	}

	// Goroutine leak (0-20 points)
	if goroutineGrowth > 0 {
		goroutineScore := float64(goroutineGrowth) * 2.0
		score += math.Min(20.0, goroutineScore)
	}

	// GC frequency (0-10 points)
	if gcPerSecond > 10.0 {
		gcScore := (gcPerSecond - 10.0) / 10.0 * 10.0
		score += math.Min(10.0, gcScore)
	}

	return math.Min(100.0, score)
}

// assessRisk - Assess risk level based on score
func assessRisk(score float64) string {
	if score < 10.0 {
		return "✅ LOW RISK (No leak detected)"
	} else if score < 25.0 {
		return "⚠️  MEDIUM RISK (Potential issues)"
	} else if score < 50.0 {
		return "🔴 HIGH RISK (Likely leak)"
	} else {
		return "🔴🔴 CRITICAL (Definite leak)"
	}
}

// formatBytes - Format bytes as human-readable string
func formatBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024.0)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024.0/1024.0)
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/1024.0/1024.0/1024.0)
	}
}

// formatBytesDelta - Format byte delta (can be negative)
func formatBytesDelta(delta int64) string {
	if delta == 0 {
		return "0 B"
	}
	sign := "+"
	bytes := uint64(delta)
	if delta < 0 {
		sign = "-"
		bytes = uint64(-delta)
	}
	return fmt.Sprintf("%s%s", sign, formatBytes(bytes))
}

// PrintRealTimeMonitoring - Print real-time memory stats
func (ma *MemoryAnalyzer) printRealtimeStats() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("REAL-TIME MEMORY MONITORING")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-15s %-15s %-15s %-15s\n",
		"Time", "Heap Alloc", "Growth", "Ops Count")
	fmt.Println(strings.Repeat("-", 70))

	ma.mu.RLock()
	defer ma.mu.RUnlock()

	var prevHeap uint64 = 0
	for i, snap := range ma.snapshots {
		elapsed := snap.time.Sub(ma.startTime)
		growth := ""
		if i > 0 && prevHeap > 0 {
			delta := int64(snap.heapAlloc) - int64(prevHeap)
			growth = formatBytesDelta(delta)
		}

		fmt.Printf("%-15s %-15s %-15s %-15d\n",
			elapsed.String(),
			formatBytes(snap.heapAlloc),
			growth,
			snap.liveObjects,
		)
		prevHeap = snap.heapAlloc
	}
}

func main() {
	durationFlag := flag.Int("duration", 30, "Test duration in seconds")
	rateFlag := flag.Int64("rate", 10000, "Operations per second")
	intervalFlag := flag.Int("interval", 5, "Snapshot interval in seconds")
	flag.Parse()

	analyzer := &MemoryAnalyzer{
		snapshots: make([]*MemorySnapshot, 0),
	}

	fmt.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Println("GRIE MEMORY LEAK ANALYZER")
	fmt.Println("═══════════════════════════════════════════════════════════════════\n")

	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("  Duration:        %d seconds\n", *durationFlag)
	fmt.Printf("  Target Rate:     %d ops/sec\n", *rateFlag)
	fmt.Printf("  Sample Interval: %d seconds\n", *intervalFlag)
	fmt.Printf("  CPUs:            %d\n\n", runtime.NumCPU())

	// Sampling goroutine
	stopChan := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(*intervalFlag) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				snap := analyzer.record()
				fmt.Printf("📊 [%v] Heap: %v | Allocs: %d | GCs: %d | Goroutines: %d\n",
					snap.time.Sub(analyzer.startTime).Round(1*time.Second),
					formatBytes(snap.heapAlloc),
					snap.liveObjects,
					snap.gcRuns,
					snap.goroutines,
				)
			}
		}
	}()

	// Initial snapshot
	analyzer.record()

	fmt.Println("🚀 Starting memory leak detection...\n")

	// Run workload
	analyzer.simulateWorkload(*durationFlag, *rateFlag)

	// Final snapshot
	analyzer.record()

	// Stop sampling
	close(stopChan)
	time.Sleep(100 * time.Millisecond)

	// Print analysis
	analyzer.printRealtimeStats()
	analyzer.analyzeLeaks()

	opsCompleted := atomic.LoadInt64(&analyzer.opsCount)
	if opsCompleted > 0 {
		fmt.Printf("\n✅ Analysis complete: %d operations processed\n", opsCompleted)
	}
}
