package benchmark

import (
	"fmt"
	stdruntime "runtime"
	"time"

	"github.com/freelang-ai/gofree/internal/compiler"
	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/memory"
	"github.com/freelang-ai/gofree/internal/parser"
	"github.com/freelang-ai/gofree/internal/runtime"
)

// BenchmarkResult holds a single benchmark measurement
type BenchmarkResult struct {
	Name      string
	Time      time.Duration
	MemBefore uint64
	MemAfter  uint64
	Iterations int
}

// BenchmarkComparison compares before/after optimization results
type BenchmarkComparison struct {
	Metric         string
	Before         BenchmarkResult
	After          BenchmarkResult
	ImprovementPct float64 // Positive = improvement
}

// BenchmarkSuite manages comprehensive performance testing
type BenchmarkSuite struct {
	results       []BenchmarkResult
	comparisons   []BenchmarkComparison
	targetImprovements map[string]float64 // Target improvements by metric
}

// NewBenchmarkSuite creates a new benchmark suite with targets
func NewBenchmarkSuite() *BenchmarkSuite {
	return &BenchmarkSuite{
		results: make([]BenchmarkResult, 0),
		comparisons: make([]BenchmarkComparison, 0),
		targetImprovements: map[string]float64{
			"compile_time": 0.20, // 20% improvement target
			"exec_time":    0.15, // 15% improvement target
			"memory":       0.25, // 25% improvement target
		},
	}
}

// MeasureCompileTime benchmarks the compilation phase
func (bs *BenchmarkSuite) MeasureCompileTime(src string, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "compile_time",
		Iterations: iterations,
	}

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		l := lexer.NewLexer(src)
		tokens := l.Tokenize()

		p := parser.NewParser(tokens)
		module, _ := p.Parse()

		c := compiler.NewCompiler()

		start := time.Now()
		_, _ = c.Compile(module)
		elapsed := time.Since(start)

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	return result
}

// MeasureExecutionTime benchmarks VM execution
func (bs *BenchmarkSuite) MeasureExecutionTime(src string, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "exec_time",
		Iterations: iterations,
	}

	// Compile once
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()
	c := compiler.NewCompiler()
	program, _ := c.Compile(module)

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		vm := runtime.NewVM(program)

		start := time.Now()
		_ = vm.Run()
		elapsed := time.Since(start)

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	return result
}

// MeasureMemory benchmarks memory allocation and GC
func (bs *BenchmarkSuite) MeasureMemory(src string, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "memory",
		Iterations: iterations,
	}

	// Compile once
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()
	c := compiler.NewCompiler()
	program, _ := c.Compile(module)

	var memAfter uint64
	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		var m1 stdruntime.MemStats
		stdruntime.ReadMemStats(&m1)
		result.MemBefore = m1.Alloc

		vm := runtime.NewVM(program)

		start := time.Now()
		_ = vm.Run()
		elapsed := time.Since(start)

		var m2 stdruntime.MemStats
		stdruntime.ReadMemStats(&m2)
		memAfter = m2.Alloc

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	result.MemAfter = memAfter
	return result
}

// MeasureCompilerOptimization benchmarks optimizer performance
func (bs *BenchmarkSuite) MeasureCompilerOptimization(src string, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "compiler_optimization",
		Iterations: iterations,
	}

	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()

	c := compiler.NewCompiler()
	program, _ := c.Compile(module)

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		opt := compiler.NewOptimizer(program)

		start := time.Now()
		_ = opt.Optimize()
		elapsed := time.Since(start)

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	return result
}

// MeasureMemoryGC benchmarks garbage collection optimization
func (bs *BenchmarkSuite) MeasureMemoryGC(objectCount int, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "memory_gc",
		Iterations: iterations,
	}

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		gc := memory.NewGC()
		optimizer := memory.NewGCOptimizer(gc)

		// Allocate objects
		for j := 0; j < objectCount; j++ {
			gc.Allocate(j, 8, "number")
		}

		start := time.Now()
		optimizer.OptimizeCollectionTiming()
		optimizer.ReduceFragmentation()
		optimizer.OptimizeRootTracking()
		elapsed := time.Since(start)

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	return result
}

// MeasureVMOptimization benchmarks VM optimization
func (bs *BenchmarkSuite) MeasureVMOptimization(src string, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Name:       "vm_optimization",
		Iterations: iterations,
	}

	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()
	c := compiler.NewCompiler()
	program, _ := c.Compile(module)

	vm := runtime.NewVM(program)

	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		opt := runtime.NewVMOptimizer(vm)

		start := time.Now()
		opt.OptimizeInstructions()
		opt.OptimizeStackUsage()
		opt.CacheFrequentOpcodes()
		elapsed := time.Since(start)

		totalTime += elapsed
	}

	result.Time = totalTime / time.Duration(iterations)
	return result
}

// Compare creates a before/after comparison
func (bs *BenchmarkSuite) Compare(metric string, before, after BenchmarkResult) BenchmarkComparison {
	improvement := (before.Time.Seconds() - after.Time.Seconds()) / before.Time.Seconds()
	comparison := BenchmarkComparison{
		Metric:         metric,
		Before:         before,
		After:          after,
		ImprovementPct: improvement,
	}

	bs.comparisons = append(bs.comparisons, comparison)
	return comparison
}

// RecordResult adds a benchmark result
func (bs *BenchmarkSuite) RecordResult(result BenchmarkResult) {
	bs.results = append(bs.results, result)
}

// GetComparisons returns all comparisons
func (bs *BenchmarkSuite) GetComparisons() []BenchmarkComparison {
	return bs.comparisons
}

// Report generates a comprehensive benchmark report
func (bs *BenchmarkSuite) Report() string {
	report := "========== FreeLang Benchmark Report ==========\n\n"

	// Results summary
	report += fmt.Sprintf("Individual Measurements (%d total):\n", len(bs.results))
	for _, r := range bs.results {
		report += fmt.Sprintf("  %s: %.3f ms (iterations: %d)\n",
			r.Name, r.Time.Seconds()*1000, r.Iterations)
		if r.MemAfter > 0 {
			memDiff := int64(r.MemAfter) - int64(r.MemBefore)
			report += fmt.Sprintf("    Memory change: %+d bytes\n", memDiff)
		}
	}

	report += "\nOptimization Improvements:\n"
	for _, comp := range bs.comparisons {
		target := bs.targetImprovements[comp.Metric]
		achieved := comp.ImprovementPct
		status := "✓"
		if achieved < target {
			status = "✗"
		}

		report += fmt.Sprintf("  %s: %s %.1f%% achieved (target: %.1f%%)\n",
			comp.Metric, status, achieved*100, target*100)
	}

	report += "\nTarget Improvements:\n"
	for metric, target := range bs.targetImprovements {
		report += fmt.Sprintf("  %s: %.1f%%\n", metric, target*100)
	}

	return report
}

// MemStats wraps Go's memory statistics for our use
type MemStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// ReadMemStats captures memory statistics
func ReadMemStats(stats *MemStats) {
	var m stdruntime.MemStats
	stdruntime.ReadMemStats(&m)
	stats.Alloc = m.Alloc
	stats.TotalAlloc = m.TotalAlloc
	stats.Sys = m.Sys
	stats.NumGC = m.NumGC
}

// BenchmarkProfile provides detailed profiling information
type BenchmarkProfile struct {
	Name           string
	Duration       time.Duration
	AllocatedBytes uint64
	GCRuns         uint32
	Instructions   int64
	StackDepth     int64
}

// ProfileResult captures profiling information for a benchmark run
func (bs *BenchmarkSuite) ProfileResult(name string, duration time.Duration,
	allocBytes uint64, gcRuns uint32) BenchmarkProfile {
	return BenchmarkProfile{
		Name:           name,
		Duration:       duration,
		AllocatedBytes: allocBytes,
		GCRuns:         gcRuns,
	}
}

// ProfileReport generates a detailed profile report
func (prof *BenchmarkProfile) Report() string {
	return fmt.Sprintf(`Profile: %s
  Duration: %.3f ms
  Allocated: %d bytes
  GC Runs: %d
  Stack Depth: %d
`,
		prof.Name,
		prof.Duration.Seconds()*1000,
		prof.AllocatedBytes,
		prof.GCRuns,
		prof.StackDepth)
}
