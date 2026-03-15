package compiler

import (
	"testing"
)

// Test JIT compiler creation
func TestNewJITCompiler(t *testing.T) {
	program := NewIRProgram()
	jit := NewJITCompiler(program)

	if jit == nil {
		t.Errorf("expected JIT compiler, got nil")
	}

	if jit.program != program {
		t.Errorf("JIT compiler should reference the program")
	}

	if jit.hotspotDetector == nil {
		t.Errorf("expected hotspot detector, got nil")
	}

	if jit.jitCache == nil {
		t.Errorf("expected JIT cache, got nil")
	}
}

// Test hotspot detector creation
func TestNewHotspotDetector(t *testing.T) {
	hd := NewHotspotDetector(100)

	if hd == nil {
		t.Errorf("expected hotspot detector, got nil")
	}

	if hd.threshold != 100 {
		t.Errorf("expected threshold 100, got %d", hd.threshold)
	}
}

// Test execution counting
func TestRecordExecution(t *testing.T) {
	hd := NewHotspotDetector(5)

	hd.RecordExecution("func1")
	hd.RecordExecution("func1")
	hd.RecordExecution("func1")

	count := hd.GetExecutionCount("func1")
	if count != 3 {
		t.Errorf("expected execution count 3, got %d", count)
	}
}

// Test hotspot detection
func TestIsHotspot(t *testing.T) {
	hd := NewHotspotDetector(3)

	// Record executions below threshold
	hd.RecordExecution("func1")
	hd.RecordExecution("func1")
	if hd.IsHotspot("func1") {
		t.Errorf("func1 should not be hotspot yet")
	}

	// Record execution reaching threshold
	hd.RecordExecution("func1")
	if !hd.IsHotspot("func1") {
		t.Errorf("func1 should be hotspot at threshold")
	}
}

// Test multiple hotspots
func TestGetHotspots(t *testing.T) {
	hd := NewHotspotDetector(2)

	hd.RecordExecution("func1")
	hd.RecordExecution("func1")
	hd.RecordExecution("func2")
	hd.RecordExecution("func2")
	hd.RecordExecution("func3") // Below threshold

	hotspots := hd.GetHotspots()
	if len(hotspots) != 2 {
		t.Errorf("expected 2 hotspots, got %d", len(hotspots))
	}

	if count, ok := hotspots["func1"]; !ok || count != 2 {
		t.Errorf("expected func1 in hotspots with count 2")
	}

	if count, ok := hotspots["func2"]; !ok || count != 2 {
		t.Errorf("expected func2 in hotspots with count 2")
	}

	if _, ok := hotspots["func3"]; ok {
		t.Errorf("func3 should not be in hotspots")
	}
}

// Test threshold adjustment
func TestSetThreshold(t *testing.T) {
	hd := NewHotspotDetector(10)

	hd.RecordExecution("func1")
	hd.RecordExecution("func1")
	hd.RecordExecution("func1")

	if hd.IsHotspot("func1") {
		t.Errorf("func1 should not be hotspot at threshold 10")
	}

	hd.SetThreshold(2)
	if !hd.IsHotspot("func1") {
		t.Errorf("func1 should be hotspot at threshold 2")
	}
}

// Test JIT cache creation
func TestNewJITCache(t *testing.T) {
	cache := NewJITCache()

	if cache == nil {
		t.Errorf("expected JIT cache, got nil")
	}

	if cache.Size() != 0 {
		t.Errorf("expected empty cache, got %d entries", cache.Size())
	}
}

// Test cache storage and retrieval
func TestCacheStoreAndGet(t *testing.T) {
	cache := NewJITCache()

	code := &JITCode{
		FunctionName: "func1",
		NativeCode:   []byte{0x48, 0xB8, 0xC3},
	}

	cache.Store("func1", code)

	retrieved, exists := cache.Get("func1")
	if !exists {
		t.Errorf("expected code in cache")
	}

	if retrieved.FunctionName != "func1" {
		t.Errorf("expected func1, got %s", retrieved.FunctionName)
	}
}

// Test cache size
func TestCacheSize(t *testing.T) {
	cache := NewJITCache()

	for i := 1; i <= 5; i++ {
		code := &JITCode{
			FunctionName: "func" + string(rune('0'+i)),
			NativeCode:   []byte{0x48, 0xB8},
		}
		cache.Store("func"+string(rune('0'+i)), code)
	}

	if cache.Size() != 5 {
		t.Errorf("expected cache size 5, got %d", cache.Size())
	}
}

// Test cache clear
func TestCacheClear(t *testing.T) {
	cache := NewJITCache()

	code := &JITCode{
		FunctionName: "func1",
		NativeCode:   []byte{0x48, 0xB8},
	}
	cache.Store("func1", code)

	if cache.Size() != 1 {
		t.Errorf("expected cache size 1, got %d", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("expected empty cache after clear")
	}
}

// Test cache memory size
func TestCacheMemorySize(t *testing.T) {
	cache := NewJITCache()

	code1 := &JITCode{
		FunctionName: "func1",
		NativeCode:   []byte{0x48, 0xB8, 0xC3},
	}

	code2 := &JITCode{
		FunctionName: "func2",
		NativeCode:   []byte{0x48, 0x01, 0xC1, 0xC3},
	}

	cache.Store("func1", code1)
	cache.Store("func2", code2)

	memSize := cache.GetMemorySize()
	expectedSize := int64(len(code1.NativeCode) + len(code2.NativeCode))
	if memSize != expectedSize {
		t.Errorf("expected memory size %d, got %d", expectedSize, memSize)
	}
}

// Test JIT compilation
func TestCompileForJIT(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
		NewInstruction(OpReturn),
	)

	jit := NewJITCompiler(program)
	err := jit.CompileForJIT("func1")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := jit.GetStatistics()
	if stats.CompiledFunctions != 1 {
		t.Errorf("expected 1 compiled function, got %d", stats.CompiledFunctions)
	}
}

// Test JIT compilation with invalid function
func TestCompileForJITInvalidFunction(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	jit := NewJITCompiler(program)
	err := jit.CompileForJIT("nonexistent")

	if err == nil {
		t.Errorf("expected error for nonexistent function")
	}
}

// Test cache hits and misses
func TestCacheHitsAndMisses(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)

	jit := NewJITCompiler(program)

	// First compilation should be a cache miss
	jit.CompileForJIT("func1")
	stats1 := jit.GetStatistics()
	if stats1.CacheMisses != 1 {
		t.Errorf("expected 1 cache miss, got %d", stats1.CacheMisses)
	}

	// Second compilation should be a cache hit
	jit.CompileForJIT("func1")
	stats2 := jit.GetStatistics()
	if stats2.CacheHits != 1 {
		t.Errorf("expected 1 cache hit, got %d", stats2.CacheHits)
	}
}

// Test optimize hotspots
func TestOptimizeHotspots(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)
	program.Functions["func2"] = NewIRFunction("func2")
	program.Functions["func2"].Instructions = append(
		program.Functions["func2"].Instructions,
		NewInstruction(OpLoadConst, 43.0),
	)

	jit := NewJITCompiler(program)

	// Record executions to create hotspots
	for i := 0; i < 100; i++ {
		jit.hotspotDetector.RecordExecution("func1")
	}
	for i := 0; i < 150; i++ {
		jit.hotspotDetector.RecordExecution("func2")
	}

	err := jit.OptimizeHotspots()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := jit.GetStatistics()
	if stats.HotspotsFunctions != 2 {
		t.Errorf("expected 2 hotspot functions, got %d", stats.HotspotsFunctions)
	}

	if stats.CompiledFunctions != 2 {
		t.Errorf("expected 2 compiled functions, got %d", stats.CompiledFunctions)
	}
}

// Test profile function
func TestProfileFunction(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)

	jit := NewJITCompiler(program)

	// Record some executions
	for i := 0; i < 150; i++ {
		jit.hotspotDetector.RecordExecution("func1")
	}

	profile := jit.ProfileFunction("func1")
	if profile.FunctionName != "func1" {
		t.Errorf("expected func1, got %s", profile.FunctionName)
	}

	if profile.ExecutionCount != 150 {
		t.Errorf("expected execution count 150, got %d", profile.ExecutionCount)
	}

	if !profile.IsHotspot {
		t.Errorf("expected func1 to be hotspot")
	}
}

// Test JIT estimate speedup
func TestJITEstimateSpeedup(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)

	jit := NewJITCompiler(program)

	// Record hotspot and compile
	for i := 0; i < 100; i++ {
		jit.hotspotDetector.RecordExecution("func1")
	}

	jit.OptimizeHotspots()

	speedup := jit.EstimateSpeedup()
	if speedup <= 1.0 {
		t.Errorf("expected speedup > 1.0, got %.2f", speedup)
	}

	if speedup > 4.0 {
		t.Errorf("expected reasonable speedup, got %.2f", speedup)
	}
}

// Test statistics structure
func TestJITCompilationStats(t *testing.T) {
	stats := JITCompilationStats{
		HotspotsFunctions:  5,
		CompiledFunctions:  4,
		CacheHits:          20,
		CacheMisses:        10,
		NativeCodeSize:     2048,
		EstimatedSpeedup:   2.5,
	}

	if stats.HotspotsFunctions != 5 {
		t.Errorf("expected 5 hotspots")
	}

	if stats.EstimatedSpeedup != 2.5 {
		t.Errorf("expected 2.5 speedup")
	}
}

// Test JIT profile output
func TestJITProfileOutput(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")

	jit := NewJITCompiler(program)
	for i := 0; i < 100; i++ {
		jit.hotspotDetector.RecordExecution("func1")
	}
	jit.OptimizeHotspots()

	profile := jit.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 50 {
		t.Errorf("expected detailed profile, got %d chars", len(profile))
	}
}

// Test reset statistics
func TestResetStatistics(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)

	jit := NewJITCompiler(program)

	for i := 0; i < 100; i++ {
		jit.hotspotDetector.RecordExecution("func1")
	}
	jit.OptimizeHotspots()

	stats := jit.GetStatistics()
	if stats.CompiledFunctions == 0 {
		t.Errorf("expected compiled functions before reset")
	}

	jit.ResetStatistics()
	stats = jit.GetStatistics()
	if stats.CompiledFunctions != 0 {
		t.Errorf("expected 0 compiled functions after reset")
	}
}

// Test clear cache
func TestJITClearCache(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(
		program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0),
	)

	jit := NewJITCompiler(program)
	jit.CompileForJIT("func1")

	stats := jit.GetStatistics()
	if stats.CacheSize != 1 {
		t.Errorf("expected cache size 1, got %d", stats.CacheSize)
	}

	jit.ClearCache()

	stats = jit.GetStatistics()
	if stats.CacheSize != 0 {
		t.Errorf("expected cache size 0 after clear")
	}
}

// Test JIT code object
func TestJITCode(t *testing.T) {
	code := &JITCode{
		FunctionName:   "func1",
		NativeCode:     []byte{0x48, 0xB8, 0xC3},
		CompileTime:    1000000,
		ExecutionCount: 0,
	}

	if code.FunctionName != "func1" {
		t.Errorf("expected func1, got %s", code.FunctionName)
	}

	if len(code.NativeCode) != 3 {
		t.Errorf("expected 3 bytes, got %d", len(code.NativeCode))
	}
}

// Test nil program handling
func TestJITNilProgram(t *testing.T) {
	jit := NewJITCompiler(nil)
	err := jit.OptimizeHotspots()

	if err == nil {
		t.Errorf("expected error for nil program")
	}
}

// Test execution profile
func TestExecutionProfile(t *testing.T) {
	profile := ExecutionProfile{
		FunctionName:   "func1",
		ExecutionCount: 100,
		IsHotspot:      true,
		AverageTime:    150.5,
	}

	if profile.FunctionName != "func1" {
		t.Errorf("expected func1")
	}

	if profile.ExecutionCount != 100 {
		t.Errorf("expected execution count 100")
	}

	if !profile.IsHotspot {
		t.Errorf("expected hotspot flag true")
	}
}

// Benchmark JIT compilation
func BenchmarkJITCompilation(b *testing.B) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	for i := 0; i < 20; i++ {
		name := "func" + string(rune('0'+i%10))
		if program.Functions[name] == nil {
			program.Functions[name] = NewIRFunction(name)
		}
		program.Functions[name].Instructions = append(
			program.Functions[name].Instructions,
			NewInstruction(OpLoadConst, float64(i)))
	}

	jit := NewJITCompiler(program)

	// Record hotspots
	for name := range program.Functions {
		for i := 0; i < 100; i++ {
			jit.hotspotDetector.RecordExecution(name)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jit.OptimizeHotspots()
	}
}

// Benchmark hotspot detection
func BenchmarkHotspotDetection(b *testing.B) {
	hd := NewHotspotDetector(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			hd.RecordExecution("func" + string(rune('0'+j%10)))
		}
	}
}

// Benchmark cache operations
func BenchmarkCacheOperations(b *testing.B) {
	cache := NewJITCache()

	codes := make([]*JITCode, 20)
	for i := 0; i < 20; i++ {
		codes[i] = &JITCode{
			FunctionName: "func" + string(rune('0'+i%10)),
			NativeCode:   []byte{0x48, 0xB8, 0xC3},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			cache.Store(codes[j].FunctionName, codes[j])
			_, _ = cache.Get(codes[j].FunctionName)
		}
	}
}

// Benchmark profiling
func BenchmarkProfiling(b *testing.B) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	for i := 0; i < 20; i++ {
		name := "func" + string(rune('0'+i%10))
		if program.Functions[name] == nil {
			program.Functions[name] = NewIRFunction(name)
		}
	}

	jit := NewJITCompiler(program)

	// Pre-record some executions
	for name := range program.Functions {
		for i := 0; i < 100; i++ {
			jit.hotspotDetector.RecordExecution(name)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for name := range program.Functions {
			_ = jit.ProfileFunction(name)
		}
	}
}
