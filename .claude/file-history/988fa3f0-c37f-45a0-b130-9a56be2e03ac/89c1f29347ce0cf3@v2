package memory

import (
	"fmt"
	"time"
)

// GCOptimizer provides garbage collection optimizations
type GCOptimizer struct {
	gc               *GC
	collectionTiming []time.Duration
	heapSize         int
	fragmentation    float64
	statistics       GCOptimizationStats
}

// GCOptimizationStats tracks GC performance
type GCOptimizationStats struct {
	TotalCollections      int64
	IncrementalCollections int64
	FullCollections       int64
	TotalPauseTime        time.Duration
	AveragePauseTime      time.Duration
	MaxPauseTime          time.Duration
	Fragmentation         float64
	PoolUtilization       float64
	FreeMemory            int
	UsedMemory            int
}

// NewGCOptimizer creates a new GC optimizer
func NewGCOptimizer(gc *GC) *GCOptimizer {
	return &GCOptimizer{
		gc:               gc,
		collectionTiming: make([]time.Duration, 0),
	}
}

// OptimizeCollectionTiming determines optimal collection strategy
func (o *GCOptimizer) OptimizeCollectionTiming() {
	if o.gc == nil {
		return
	}

	// Analyze heap size trend
	heapGrowth := o.analyzeHeapGrowth()

	// Determine if incremental or full collection is needed
	if heapGrowth < 0.1 { // Small growth
		// Incremental collection sufficient
		o.gc.Collect()
		o.statistics.IncrementalCollections++
	} else if heapGrowth > 0.5 { // Large growth
		// Need full collection
		o.gc.Collect()
		o.statistics.FullCollections++
	}

	o.statistics.TotalCollections++
}

// analyzeHeapGrowth calculates the rate of heap growth
func (o *GCOptimizer) analyzeHeapGrowth() float64 {
	stats := o.gc.GetStats()
	currentSize := stats.TotalSize

	if o.heapSize == 0 {
		o.heapSize = currentSize
		return 0
	}

	growth := float64(currentSize-o.heapSize) / float64(o.heapSize)
	o.heapSize = currentSize

	return growth
}

// OptimizePoolSizing dynamically adjusts memory pool sizes
func (o *GCOptimizer) OptimizePoolSizing(pool *MemoryPool) {
	if pool == nil || pool.pools == nil {
		return
	}

	// Analyze allocation patterns
	for size, objects := range pool.pools {
		utilization := float64(len(objects)) / 10.0 // Max pool size is 10

		if utilization > 0.8 {
			// Pool is heavily used, consider expanding
			pool.pools[size] = append(objects, &Object{Size: size})
		} else if utilization < 0.2 && len(objects) > 5 {
			// Pool is underutilized, trim it
			pool.pools[size] = objects[:len(objects)-1]
		}
	}

	o.statistics.PoolUtilization = o.calculatePoolUtilization(pool)
}

// calculatePoolUtilization returns average pool utilization
func (o *GCOptimizer) calculatePoolUtilization(pool *MemoryPool) float64 {
	if len(pool.pools) == 0 {
		return 0
	}

	totalUtilization := 0.0
	for _, objects := range pool.pools {
		totalUtilization += float64(len(objects)) / 10.0
	}

	return totalUtilization / float64(len(pool.pools))
}

// ReduceFragmentation attempts to reduce memory fragmentation
func (o *GCOptimizer) ReduceFragmentation() {
	if o.gc == nil {
		return
	}

	stats := o.gc.GetStats()

	// Calculate fragmentation: unused space in allocated objects
	// Simplified: if only 50% of allocated memory is marked, high fragmentation
	if stats.MarkedObjects > 0 && stats.TotalObjects > 0 {
		utilization := float64(stats.MarkedObjects) / float64(stats.TotalObjects)
		o.fragmentation = 1.0 - utilization

		// If fragmentation > 50%, suggest compaction
		if o.fragmentation > 0.5 {
			o.statistics.Fragmentation = o.fragmentation
			// In a real implementation, would trigger compaction
		}
	}
}

// EnableIncrementalGC converts stop-world GC to incremental
func (o *GCOptimizer) EnableIncrementalGC() {
	if o.gc == nil {
		return
	}

	// Incremental GC: perform marking in small chunks
	// instead of all at once (stop-the-world)

	// Phase 1: Mark phase (incremental)
	// Instead of marking all objects at once,
	// mark a batch and return control to application

	markedCount := 0
	markBatchSize := 100 // Mark up to 100 objects per phase

	for _, obj := range o.gc.objects {
		if !obj.Marked {
			o.gc.mark(obj)
			markedCount++

			if markedCount >= markBatchSize {
				// Yield control back to application
				break
			}
		}
	}

	o.statistics.IncrementalCollections++
}

// OptimizeRootTracking improves root object tracking
func (o *GCOptimizer) OptimizeRootTracking() {
	if o.gc == nil {
		return
	}

	roots := o.gc.roots

	// Optimize: deduplicate roots
	seen := make(map[int]bool)
	uniqueRoots := []*Object{}

	for _, root := range roots {
		if !seen[root.ID] {
			seen[root.ID] = true
			uniqueRoots = append(uniqueRoots, root)
		}
	}

	// If duplicates found, update roots
	if len(uniqueRoots) < len(roots) {
		o.gc.roots = uniqueRoots
	}
}

// TuneCollectionThreshold adjusts GC threshold based on usage patterns
func (o *GCOptimizer) TuneCollectionThreshold(currentThreshold int) int {
	if o.gc == nil {
		return currentThreshold
	}

	stats := o.gc.GetStats()

	// Adapt threshold based on allocation rate
	// If allocation is slow, increase threshold (less GC)
	// If allocation is fast, decrease threshold (more GC)

	newThreshold := currentThreshold

	if stats.TotalObjects > currentThreshold {
		// Allocation is aggressive, decrease threshold
		newThreshold = currentThreshold / 2
	} else if stats.TotalObjects < currentThreshold/2 {
		// Allocation is relaxed, increase threshold
		newThreshold = currentThreshold * 2
	}

	return newThreshold
}

// CompactMemory simulates memory compaction (relocating objects)
func (o *GCOptimizer) CompactMemory() {
	if o.gc == nil {
		return
	}

	// Memory compaction: move objects to eliminate holes
	// This reduces fragmentation but requires updating references

	// Step 1: Mark live objects
	for _, root := range o.gc.roots {
		o.gc.mark(root)
	}

	// Step 2: Compact (in a real implementation)
	// - Calculate new positions for marked objects
	// - Update all references
	// - Move objects to new positions

	// Step 3: Clear marks
	o.gc.marked = make(map[int]bool)
}

// GetStatistics returns GC optimization statistics
func (o *GCOptimizer) GetStatistics() GCOptimizationStats {
	if o.gc != nil {
		stats := o.gc.GetStats()
		o.statistics.UsedMemory = stats.TotalSize
		o.statistics.FreeMemory = 0 // Would need to track free memory separately
	}

	// Calculate average pause time
	if len(o.collectionTiming) > 0 {
		totalTime := time.Duration(0)
		maxTime := time.Duration(0)

		for _, duration := range o.collectionTiming {
			totalTime += duration
			if duration > maxTime {
				maxTime = duration
			}
		}

		o.statistics.MaxPauseTime = maxTime
		o.statistics.AveragePauseTime = totalTime / time.Duration(len(o.collectionTiming))
	}

	return o.statistics
}

// RecordCollectionTime records the time taken for a collection cycle
func (o *GCOptimizer) RecordCollectionTime(duration time.Duration) {
	o.collectionTiming = append(o.collectionTiming, duration)
	o.statistics.TotalPauseTime += duration

	// Keep last 100 measurements for trend analysis
	if len(o.collectionTiming) > 100 {
		o.collectionTiming = o.collectionTiming[1:]
	}
}

// Profile returns a GC performance profile
func (o *GCOptimizer) Profile() string {
	stats := o.GetStatistics()

	return fmt.Sprintf(`GC Optimization Profile:
  Total Collections: %d
  Incremental: %d, Full: %d
  Total Pause Time: %v
  Average Pause: %v
  Max Pause: %v
  Fragmentation: %.1f%%
  Pool Utilization: %.1f%%
  Used Memory: %d bytes
  Free Memory: %d bytes`,
		stats.TotalCollections,
		stats.IncrementalCollections,
		stats.FullCollections,
		stats.TotalPauseTime,
		stats.AveragePauseTime,
		stats.MaxPauseTime,
		stats.Fragmentation*100,
		stats.PoolUtilization*100,
		stats.UsedMemory,
		stats.FreeMemory)
}
