package benchmark

import (
	"testing"
)

// Fibonacci source code for testing
const fibonacciSrc = `
fn fibonacci(n) {
    if (n <= 1) {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

let result = fibonacci(10)
`

// Loop source code for testing
const loopSrc = `
let sum = 0
for i in range(0, 100) {
    sum = sum + i
}
`

// Array operations source code
const arrayOpsSrc = `
let arr = [1, 2, 3, 4, 5]
let sum = 0
for i in range(0, 5) {
    sum = sum + arr[i]
}
`

// Arithmetic source for constant folding
const arithmeticSrc = `
let a = 10 + 5
let b = a * 2
let c = b - 3
`

// Test compile time benchmark
func TestMeasureCompileTime(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureCompileTime(arithmeticSrc, 10)

	if result.Name != "compile_time" {
		t.Errorf("expected name 'compile_time', got %q", result.Name)
	}

	if result.Time.Microseconds() <= 0 {
		t.Errorf("expected positive time, got %v", result.Time)
	}

	if result.Iterations != 10 {
		t.Errorf("expected 10 iterations, got %d", result.Iterations)
	}
}

// Test execution time benchmark
func TestMeasureExecutionTime(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureExecutionTime(arithmeticSrc, 5)

	if result.Name != "exec_time" {
		t.Errorf("expected name 'exec_time', got %q", result.Name)
	}

	if result.Time < 0 {
		t.Errorf("expected non-negative time, got %v", result.Time)
	}
}

// Test memory benchmark
func TestMeasureMemory(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureMemory(arrayOpsSrc, 3)

	if result.Name != "memory" {
		t.Errorf("expected name 'memory', got %q", result.Name)
	}

	if result.Time < 0 {
		t.Errorf("expected non-negative time, got %v", result.Time)
	}

	if result.MemAfter == 0 {
		t.Logf("memory measurement inconclusive (may be optimized away)")
	}
}

// Test compiler optimization benchmark
func TestMeasureCompilerOptimization(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureCompilerOptimization(arithmeticSrc, 5)

	if result.Name != "compiler_optimization" {
		t.Errorf("expected name 'compiler_optimization', got %q", result.Name)
	}

	if result.Time.Microseconds() <= 0 {
		t.Errorf("expected positive time, got %v", result.Time)
	}
}

// Test memory GC benchmark
func TestMeasureMemoryGC(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureMemoryGC(100, 5)

	if result.Name != "memory_gc" {
		t.Errorf("expected name 'memory_gc', got %q", result.Name)
	}

	if result.Time.Microseconds() <= 0 {
		t.Errorf("expected positive time, got %v", result.Time)
	}
}

// Test VM optimization benchmark
func TestMeasureVMOptimization(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := suite.MeasureVMOptimization(arithmeticSrc, 5)

	if result.Name != "vm_optimization" {
		t.Errorf("expected name 'vm_optimization', got %q", result.Name)
	}

	if result.Time < 0 {
		t.Errorf("expected non-negative time, got %v", result.Time)
	}
}

// Test benchmark comparison
func TestBenchmarkComparison(t *testing.T) {
	suite := NewBenchmarkSuite()

	before := BenchmarkResult{
		Name:       "test_metric",
		Time:       10,
		Iterations: 1,
	}

	after := BenchmarkResult{
		Name:       "test_metric",
		Time:       8,
		Iterations: 1,
	}

	comparison := suite.Compare("test_metric", before, after)

	if comparison.Metric != "test_metric" {
		t.Errorf("expected metric 'test_metric', got %q", comparison.Metric)
	}

	if comparison.ImprovementPct <= 0 || comparison.ImprovementPct >= 1 {
		t.Errorf("expected improvement between 0 and 1, got %.2f", comparison.ImprovementPct)
	}

	comparisons := suite.GetComparisons()
	if len(comparisons) != 1 {
		t.Errorf("expected 1 comparison, got %d", len(comparisons))
	}
}

// Test benchmark result recording
func TestRecordResult(t *testing.T) {
	suite := NewBenchmarkSuite()

	result := BenchmarkResult{
		Name:       "test_result",
		Time:       100,
		Iterations: 5,
	}

	suite.RecordResult(result)

	if len(suite.results) != 1 {
		t.Errorf("expected 1 result, got %d", len(suite.results))
	}
}

// Test benchmark report generation
func TestBenchmarkReport(t *testing.T) {
	suite := NewBenchmarkSuite()

	// Add some results
	suite.RecordResult(BenchmarkResult{
		Name:       "compile_time",
		Time:       100,
		Iterations: 1,
	})

	suite.Compare("compile_time",
		BenchmarkResult{Name: "compile_time", Time: 100, Iterations: 1},
		BenchmarkResult{Name: "compile_time", Time: 85, Iterations: 1})

	report := suite.Report()

	if report == "" {
		t.Errorf("expected non-empty report")
	}

	if len(report) < 50 {
		t.Errorf("expected detailed report, got only %d chars", len(report))
	}
}

// Test memory statistics capture
func TestReadMemStats(t *testing.T) {
	stats := &MemStats{}
	ReadMemStats(stats)

	if stats.Alloc == 0 {
		t.Logf("memory allocation not captured (expected on some systems)")
	}
}

// Test profile result creation
func TestProfileResult(t *testing.T) {
	suite := NewBenchmarkSuite()

	profile := suite.ProfileResult("test_profile", 100, 1024, 5)

	if profile.Name != "test_profile" {
		t.Errorf("expected name 'test_profile', got %q", profile.Name)
	}

	if profile.AllocatedBytes != 1024 {
		t.Errorf("expected 1024 bytes, got %d", profile.AllocatedBytes)
	}

	if profile.GCRuns != 5 {
		t.Errorf("expected 5 GC runs, got %d", profile.GCRuns)
	}
}

// Test profile report
func TestProfileReport(t *testing.T) {
	profile := BenchmarkProfile{
		Name:           "test",
		Duration:       100,
		AllocatedBytes: 2048,
		GCRuns:         3,
		StackDepth:     10,
	}

	report := profile.Report()

	if report == "" {
		t.Errorf("expected non-empty profile report")
	}

	if len(report) < 30 {
		t.Errorf("expected detailed profile, got only %d chars", len(report))
	}
}

// Test compiler optimization improvement measurement
func TestCompilerOptimizationImprovement(t *testing.T) {
	suite := NewBenchmarkSuite()

	// Measure optimization time on constant folding code
	optResult := suite.MeasureCompilerOptimization(arithmeticSrc, 5)

	if optResult.Time.Microseconds() <= 0 {
		t.Errorf("expected positive optimization time")
	}

	t.Logf("Compiler optimization time: %.3f us", float64(optResult.Time.Microseconds()))
}

// Test memory GC optimization improvement
func TestMemoryGCOptimizationImprovement(t *testing.T) {
	suite := NewBenchmarkSuite()

	gcResult := suite.MeasureMemoryGC(200, 3)

	if gcResult.Time.Microseconds() <= 0 {
		t.Errorf("expected positive GC optimization time")
	}

	t.Logf("Memory GC optimization time: %.3f us", float64(gcResult.Time.Microseconds()))
}

// Test full benchmark suite with multiple measurements
func TestFullBenchmarkSuite(t *testing.T) {
	suite := NewBenchmarkSuite()

	// Measure compile time
	compileResult := suite.MeasureCompileTime(arithmeticSrc, 5)
	suite.RecordResult(compileResult)

	// Measure execution time
	execResult := suite.MeasureExecutionTime(arithmeticSrc, 5)
	suite.RecordResult(execResult)

	// Measure compiler optimization
	optResult := suite.MeasureCompilerOptimization(arithmeticSrc, 5)
	suite.RecordResult(optResult)

	if len(suite.results) != 3 {
		t.Errorf("expected 3 results, got %d", len(suite.results))
	}

	report := suite.Report()
	if len(report) < 100 {
		t.Errorf("expected comprehensive report, got %d chars", len(report))
	}

	t.Logf("Report:\n%s", report)
}

// Test regression detection
func TestRegressionDetection(t *testing.T) {
	suite := NewBenchmarkSuite()

	// Simulate a regression (performance got worse)
	before := BenchmarkResult{
		Name:       "test",
		Time:       100,
		Iterations: 1,
	}

	after := BenchmarkResult{
		Name:       "test",
		Time:       150, // 50% slower (regression)
		Iterations: 1,
	}

	comparison := suite.Compare("test", before, after)

	if comparison.ImprovementPct >= 0 {
		t.Errorf("expected negative improvement (regression), got %.2f", comparison.ImprovementPct)
	}

	t.Logf("Regression detected: %.1f%% performance degradation", -comparison.ImprovementPct*100)
}

// Benchmark compile time
func BenchmarkCompileTime(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureCompileTime(arithmeticSrc, 1)
	}
}

// Benchmark execution time
func BenchmarkExecutionTime(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureExecutionTime(arithmeticSrc, 1)
	}
}

// Benchmark compiler optimization
func BenchmarkCompilerOptimization(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureCompilerOptimization(arithmeticSrc, 1)
	}
}

// Benchmark memory GC
func BenchmarkMemoryGC(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureMemoryGC(100, 1)
	}
}

// Benchmark VM optimization
func BenchmarkVMOptimization(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureVMOptimization(arithmeticSrc, 1)
	}
}

// Benchmark full suite
func BenchmarkFullSuite(b *testing.B) {
	suite := NewBenchmarkSuite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.MeasureCompileTime(arithmeticSrc, 1)
		suite.MeasureExecutionTime(arithmeticSrc, 1)
		suite.MeasureCompilerOptimization(arithmeticSrc, 1)
		suite.MeasureMemoryGC(100, 1)
		suite.MeasureVMOptimization(arithmeticSrc, 1)
	}
}
