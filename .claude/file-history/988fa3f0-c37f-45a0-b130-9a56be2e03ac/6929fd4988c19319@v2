package compiler

import (
	"fmt"
	"sync"
	"time"
)

// JITCompiler manages just-in-time compilation
type JITCompiler struct {
	program          *IRProgram
	hotspotDetector  *HotspotDetector
	jitCache         *JITCache
	compilationStats JITCompilationStats
	mu               sync.RWMutex
}

// JITCompilationStats tracks JIT compilation metrics
type JITCompilationStats struct {
	HotspotsFunctions    int64
	CompiledFunctions    int64
	CacheHits            int64
	CacheMisses          int64
	CompilationTime      int64 // nanoseconds
	NativeCodeSize       int64 // bytes
	CacheSize            int64
	EstimatedSpeedup     float64 // theoretical speedup from JIT
	ExecutionCount       int64
	AverageExecutionTime float64 // nanoseconds
}

// HotspotDetector detects hot functions for JIT compilation
type HotspotDetector struct {
	executionCounts map[string]int64
	threshold       int64 // execution count threshold for JIT
	mu              sync.RWMutex
}

// ExecutionProfile tracks execution statistics
type ExecutionProfile struct {
	FunctionName    string
	ExecutionCount  int64
	TotalTime       int64 // nanoseconds
	AverageTime     float64
	LastExecuted    time.Time
	IsHotspot       bool
	CompileAttempts int
}

// JITCache stores compiled code
type JITCache struct {
	cache map[string]*JITCode
	mu    sync.RWMutex
}

// JITCode represents compiled native code
type JITCode struct {
	FunctionName   string
	Instructions   []*Instruction
	NativeCode     []byte
	CompileTime    int64 // nanoseconds
	ExecutionCount int64
	CacheTime      time.Time
}

// NewJITCompiler creates a new JIT compiler
func NewJITCompiler(program *IRProgram) *JITCompiler {
	return &JITCompiler{
		program:         program,
		hotspotDetector: NewHotspotDetector(100), // threshold: 100 executions
		jitCache:        NewJITCache(),
		compilationStats: JITCompilationStats{},
	}
}

// NewHotspotDetector creates a new hotspot detector
func NewHotspotDetector(threshold int64) *HotspotDetector {
	return &HotspotDetector{
		executionCounts: make(map[string]int64),
		threshold:       threshold,
	}
}

// NewJITCache creates a new JIT code cache
func NewJITCache() *JITCache {
	return &JITCache{
		cache: make(map[string]*JITCode),
	}
}

// RecordExecution records a function execution
func (hd *HotspotDetector) RecordExecution(funcName string) {
	hd.mu.Lock()
	defer hd.mu.Unlock()

	hd.executionCounts[funcName]++
}

// IsHotspot checks if a function is a hotspot
func (hd *HotspotDetector) IsHotspot(funcName string) bool {
	hd.mu.RLock()
	defer hd.mu.RUnlock()

	return hd.executionCounts[funcName] >= hd.threshold
}

// GetExecutionCount returns execution count for a function
func (hd *HotspotDetector) GetExecutionCount(funcName string) int64 {
	hd.mu.RLock()
	defer hd.mu.RUnlock()

	return hd.executionCounts[funcName]
}

// GetHotspots returns all functions exceeding threshold
func (hd *HotspotDetector) GetHotspots() map[string]int64 {
	hd.mu.RLock()
	defer hd.mu.RUnlock()

	hotspots := make(map[string]int64)
	for name, count := range hd.executionCounts {
		if count >= hd.threshold {
			hotspots[name] = count
		}
	}
	return hotspots
}

// SetThreshold updates the hotspot threshold
func (hd *HotspotDetector) SetThreshold(threshold int64) {
	hd.mu.Lock()
	defer hd.mu.Unlock()

	hd.threshold = threshold
}

// Store stores compiled code in cache
func (jc *JITCache) Store(funcName string, code *JITCode) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	code.CacheTime = time.Now()
	jc.cache[funcName] = code
}

// Get retrieves compiled code from cache
func (jc *JITCache) Get(funcName string) (*JITCode, bool) {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	code, exists := jc.cache[funcName]
	return code, exists
}

// Size returns cache size (number of entries)
func (jc *JITCache) Size() int {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	return len(jc.cache)
}

// Clear clears the entire cache
func (jc *JITCache) Clear() {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	jc.cache = make(map[string]*JITCode)
}

// GetMemorySize returns total memory used by cache
func (jc *JITCache) GetMemorySize() int64 {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	var totalSize int64
	for _, code := range jc.cache {
		totalSize += int64(len(code.NativeCode))
	}
	return totalSize
}

// CompileForJIT performs JIT compilation on a function
func (jc *JITCompiler) CompileForJIT(funcName string) error {
	if jc.program == nil {
		return fmt.Errorf("program is nil")
	}

	// Check cache first
	if cachedCode, exists := jc.jitCache.Get(funcName); exists {
		jc.mu.Lock()
		jc.compilationStats.CacheHits++
		jc.mu.Unlock()
		cachedCode.ExecutionCount++
		return nil
	}

	jc.mu.Lock()
	jc.compilationStats.CacheMisses++
	jc.mu.Unlock()

	// Get function from program
	fn, exists := jc.program.Functions[funcName]
	if !exists {
		return fmt.Errorf("function %s not found", funcName)
	}

	// Compile function to native code
	startTime := time.Now()
	nativeCode := jc.compileToNativeCode(fn)
	compilationTime := time.Since(startTime).Nanoseconds()

	// Create JIT code object
	jitCode := &JITCode{
		FunctionName:   funcName,
		Instructions:   fn.Instructions,
		NativeCode:     nativeCode,
		CompileTime:    compilationTime,
		ExecutionCount: 0,
	}

	// Store in cache
	jc.jitCache.Store(funcName, jitCode)

	jc.mu.Lock()
	jc.compilationStats.CompiledFunctions++
	jc.compilationStats.CompilationTime += compilationTime
	jc.compilationStats.NativeCodeSize += int64(len(nativeCode))
	jc.mu.Unlock()

	return nil
}

// compileToNativeCode converts IR instructions to native code (stub)
func (jc *JITCompiler) compileToNativeCode(fn *IRFunction) []byte {
	// Simplified: Generate x86-64 assembly-like code
	// In real implementation, would use LLVM or similar
	var code []byte

	for _, ins := range fn.Instructions {
		switch ins.Opcode {
		case OpLoadConst:
			// MOV instruction (simplified representation)
			code = append(code, 0x48, 0xB8) // MOV RAX, imm64
		case OpAdd:
			code = append(code, 0x48, 0x01, 0xC1) // ADD RAX, RCX
		case OpReturn:
			code = append(code, 0xC3) // RET
		}
	}

	return code
}

// ProfileFunction profiles a function execution
func (jc *JITCompiler) ProfileFunction(funcName string) ExecutionProfile {
	hd := jc.hotspotDetector
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	execCount := hd.GetExecutionCount(funcName)
	profile := ExecutionProfile{
		FunctionName:   funcName,
		ExecutionCount: execCount,
		IsHotspot:      hd.IsHotspot(funcName),
		LastExecuted:   time.Now(),
	}

	// Calculate average execution time if compiled
	if jitCode, exists := jc.jitCache.Get(funcName); exists {
		if jitCode.ExecutionCount > 0 {
			profile.AverageTime = float64(jitCode.CompileTime) / float64(jitCode.ExecutionCount)
		}
	}

	return profile
}

// OptimizeHotspots performs JIT compilation on all detected hotspots
func (jc *JITCompiler) OptimizeHotspots() error {
	if jc.program == nil {
		return fmt.Errorf("program is nil")
	}

	hotspots := jc.hotspotDetector.GetHotspots()

	jc.mu.Lock()
	jc.compilationStats.HotspotsFunctions = int64(len(hotspots))
	jc.mu.Unlock()

	for funcName := range hotspots {
		if err := jc.CompileForJIT(funcName); err != nil {
			return err
		}
	}

	return nil
}

// EstimateSpeedup calculates theoretical speedup from JIT compilation
func (jc *JITCompiler) EstimateSpeedup() float64 {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	if jc.compilationStats.CompiledFunctions == 0 {
		return 1.0
	}

	// Simplified: assume 2-3x speedup per compiled function
	// + cache hit benefit
	compiledRatio := float64(jc.compilationStats.CompiledFunctions) /
		float64(jc.compilationStats.HotspotsFunctions + 1)
	cacheHitRatio := float64(jc.compilationStats.CacheHits) /
		float64(jc.compilationStats.CacheHits + jc.compilationStats.CacheMisses + 1)

	speedup := 1.0 + (compiledRatio * 2.5) + (cacheHitRatio * 0.5)
	return speedup
}

// GetStatistics returns JIT compilation statistics
func (jc *JITCompiler) GetStatistics() JITCompilationStats {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	stats := jc.compilationStats
	stats.CacheSize = int64(jc.jitCache.Size())
	stats.EstimatedSpeedup = jc.EstimateSpeedup()

	return stats
}

// Profile returns a profile of JIT compilation
func (jc *JITCompiler) Profile() string {
	stats := jc.GetStatistics()

	return fmt.Sprintf(`JIT Compilation Profile:
  Hotspot Functions: %d
  Compiled Functions: %d
  Cache Size: %d
  Cache Hits: %d
  Cache Misses: %d
  Native Code Size: %d bytes
  Compilation Time: %d ns
  Estimated Speedup: %.2fx
`,
		stats.HotspotsFunctions,
		stats.CompiledFunctions,
		stats.CacheSize,
		stats.CacheHits,
		stats.CacheMisses,
		stats.NativeCodeSize,
		stats.CompilationTime,
		stats.EstimatedSpeedup)
}

// ResetStatistics clears all statistics
func (jc *JITCompiler) ResetStatistics() {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	jc.compilationStats = JITCompilationStats{}
}

// ClearCache clears the JIT cache
func (jc *JITCompiler) ClearCache() {
	jc.jitCache.Clear()
}
