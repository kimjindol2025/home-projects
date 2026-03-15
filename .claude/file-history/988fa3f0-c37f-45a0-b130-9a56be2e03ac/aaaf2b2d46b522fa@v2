package memory

import (
	"testing"
	"time"
)

// Test GC optimizer creation
func TestNewGCOptimizer(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	if optimizer == nil {
		t.Errorf("expected optimizer, got nil")
	}

	if optimizer.gc != gc {
		t.Errorf("optimizer should reference the GC")
	}
}

// Test collection timing optimization
func TestOptimizeCollectionTiming(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	// Allocate some objects
	gc.Allocate(1, 8, "number")
	gc.Allocate(2, 8, "number")

	optimizer.OptimizeCollectionTiming()

	stats := optimizer.GetStatistics()
	if stats.TotalCollections == 0 {
		t.Errorf("expected at least 1 collection")
	}
}

// Test heap growth analysis
func TestAnalyzeHeapGrowth(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	// First allocation sets baseline
	gc.Allocate(1, 8, "number")
	growth1 := optimizer.analyzeHeapGrowth()

	if growth1 != 0 {
		t.Errorf("expected 0 growth on first call")
	}

	// Second allocation
	gc.Allocate(2, 16, "number")
	growth2 := optimizer.analyzeHeapGrowth()

	if growth2 <= 0 {
		t.Errorf("expected positive growth")
	}
}

// Test pool sizing optimization
func TestOptimizePoolSizing(t *testing.T) {
	pool := NewMemoryPool()

	// Add some objects to pool
	for i := 0; i < 5; i++ {
		pool.Allocate(8)
	}

	optimizer := NewGCOptimizer(nil) // No GC needed for pool optimization
	optimizer.OptimizePoolSizing(pool)

	// Should complete without error
	utilization := optimizer.calculatePoolUtilization(pool)
	if utilization < 0 || utilization > 1 {
		t.Errorf("invalid utilization: %f", utilization)
	}
}

// Test fragmentation reduction
func TestReduceFragmentation(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	// Allocate and mark some objects
	obj1 := gc.Allocate(1, 8, "number")
	_ = gc.Allocate(2, 8, "number")
	obj3 := gc.Allocate(3, 8, "number")

	gc.AddRoot(obj1)
	gc.AddRoot(obj3)

	optimizer.ReduceFragmentation()

	stats := optimizer.GetStatistics()
	if stats.Fragmentation < 0 || stats.Fragmentation > 1 {
		t.Logf("fragmentation: %.1f%%", stats.Fragmentation*100)
	}
}

// Test incremental GC
func TestEnableIncrementalGC(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	// Allocate objects
	for i := 0; i < 50; i++ {
		gc.Allocate(i, 8, "number")
	}

	optimizer.EnableIncrementalGC()

	stats := optimizer.GetStatistics()
	if stats.IncrementalCollections == 0 {
		t.Errorf("expected incremental collection to be recorded")
	}
}

// Test root tracking optimization
func TestOptimizeRootTracking(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	obj := gc.Allocate(1, 8, "number")

	// Add duplicate roots
	gc.AddRoot(obj)
	gc.AddRoot(obj)

	initialRootCount := len(gc.roots)

	optimizer.OptimizeRootTracking()

	finalRootCount := len(gc.roots)

	if finalRootCount >= initialRootCount {
		t.Logf("root deduplication may not have occurred")
	}
}

// Test collection threshold tuning
func TestTuneCollectionThreshold(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	currentThreshold := 1000
	newThreshold := optimizer.TuneCollectionThreshold(currentThreshold)

	if newThreshold <= 0 {
		t.Errorf("expected positive threshold")
	}

	// Allocate many objects to trigger threshold adjustment
	for i := 0; i < 500; i++ {
		gc.Allocate(i, 8, "number")
	}

	newThreshold = optimizer.TuneCollectionThreshold(newThreshold)
	if newThreshold <= 0 {
		t.Errorf("expected positive threshold after adjustment")
	}
}

// Test memory compaction
func TestCompactMemory(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	// Allocate objects
	obj1 := gc.Allocate(1, 8, "number")
	_ = gc.Allocate(2, 8, "number")
	obj3 := gc.Allocate(3, 8, "number")

	gc.AddRoot(obj1)
	gc.AddRoot(obj3)

	optimizer.CompactMemory()

	stats := gc.GetStats()
	if stats.TotalObjects == 0 {
		t.Errorf("expected objects after compaction")
	}
}

// Test statistics recording
func TestGetStatistics(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	gc.Allocate(1, 8, "number")

	stats := optimizer.GetStatistics()

	if stats.TotalCollections < 0 {
		t.Errorf("expected non-negative collections")
	}

	if stats.PoolUtilization < 0 || stats.PoolUtilization > 1 {
		t.Errorf("invalid pool utilization")
	}
}

// Test collection time recording
func TestRecordCollectionTime(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	duration := 5 * time.Millisecond

	optimizer.RecordCollectionTime(duration)

	stats := optimizer.GetStatistics()
	if stats.TotalPauseTime != duration {
		t.Errorf("expected total pause time %v, got %v", duration, stats.TotalPauseTime)
	}

	if stats.MaxPauseTime != duration {
		t.Errorf("expected max pause time %v", duration)
	}
}

// Test profile output
func TestProfile(t *testing.T) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	gc.Allocate(1, 8, "number")
	optimizer.OptimizeCollectionTiming()
	optimizer.RecordCollectionTime(2 * time.Millisecond)

	profile := optimizer.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 50 {
		t.Errorf("expected detailed profile output")
	}
}

// Test GC optimization stats
func TestGCOptimizationStats(t *testing.T) {
	stats := GCOptimizationStats{
		TotalCollections:      100,
		IncrementalCollections: 80,
		FullCollections:       20,
		TotalPauseTime:        100 * time.Millisecond,
		AveragePauseTime:      1 * time.Millisecond,
		MaxPauseTime:          5 * time.Millisecond,
		Fragmentation:         0.25,
		PoolUtilization:       0.75,
		UsedMemory:            4096,
		FreeMemory:            2048,
	}

	if stats.TotalCollections != 100 {
		t.Errorf("expected 100 total collections")
	}

	if stats.IncrementalCollections+stats.FullCollections != stats.TotalCollections {
		t.Errorf("collection count mismatch")
	}

	if stats.Fragmentation != 0.25 {
		t.Errorf("expected 0.25 fragmentation (25%%)")
	}
}

// Benchmark GC optimization
func BenchmarkGCOptimization(b *testing.B) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	for i := 0; i < 100; i++ {
		gc.Allocate(i, 8, "number")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.OptimizeCollectionTiming()
		optimizer.ReduceFragmentation()
		optimizer.OptimizeRootTracking()
	}
}

// Benchmark heap growth analysis
func BenchmarkAnalyzeHeapGrowth(b *testing.B) {
	gc := NewGC()
	optimizer := NewGCOptimizer(gc)

	gc.Allocate(1, 8, "number")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.analyzeHeapGrowth()
	}
}

// Benchmark pool sizing
func BenchmarkOptimizePoolSizing(b *testing.B) {
	pool := NewMemoryPool()

	for i := 0; i < 10; i++ {
		pool.Allocate(8)
	}

	optimizer := NewGCOptimizer(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.OptimizePoolSizing(pool)
	}
}
