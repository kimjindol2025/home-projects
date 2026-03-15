package compiler

import (
	"testing"
)

// Test SIMD optimizer creation
func TestNewSIMDOptimizer(t *testing.T) {
	program := NewIRProgram()
	optimizer := NewSIMDOptimizer(program)

	if optimizer == nil {
		t.Errorf("expected optimizer, got nil")
	}

	if optimizer.program != program {
		t.Errorf("optimizer should reference the program")
	}
}

// Test SIMD optimization
func TestOptimizeSIMD(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 2))

	optimizer := NewSIMDOptimizer(program)
	err := optimizer.OptimizeSIMD()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test vectorizable pattern detection
func TestDetectVectorizablePatterns(t *testing.T) {
	program := NewIRProgram()

	// Simulate a loop pattern
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0)) // loop counter
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 100.0)) // limit
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLess))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpJumpIfFalse, 10)) // loop end

	optimizer := NewSIMDOptimizer(program)
	optimizer.detectVectorizablePatterns()

	stats := optimizer.GetStatistics()
	if stats.VectorizableLoops < 0 {
		t.Errorf("expected non-negative vectorizable loops")
	}
}

// Test vector instruction emission
func TestEmitVectorInstructions(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 2))

	optimizer := NewSIMDOptimizer(program)
	optimizer.emitVectorInstructions()

	stats := optimizer.GetStatistics()
	if stats.VectorInstructions == 0 {
		t.Logf("vector instructions may not have been emitted")
	}
}

// Test memory bandwidth analysis
func TestAnalyzeMemoryBandwidth(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 2))

	optimizer := NewSIMDOptimizer(program)
	optimizer.analyzeMemoryBandwidth()

	stats := optimizer.GetStatistics()
	if stats.MemoryBandwidth < 0 || stats.MemoryBandwidth > 1 {
		t.Errorf("invalid memory bandwidth: %.2f", stats.MemoryBandwidth)
	}
}

// Test vector type string representation
func TestVectorTypeString(t *testing.T) {
	tests := []struct {
		vt       VectorType
		expected string
	}{
		{Vec4F, "vec4f"},
		{Vec2F64, "vec2f64"},
		{Vec8I, "vec8i"},
		{Vec4I64, "vec4i64"},
	}

	for _, tt := range tests {
		if got := tt.vt.String(); got != tt.expected {
			t.Errorf("VectorType.String() = %v, expected %v", got, tt.expected)
		}
	}
}

// Test vector opcode string representation
func TestVectorOpcodeString(t *testing.T) {
	tests := []struct {
		vo       VectorOpcode
		expected string
	}{
		{VecAdd, "VADD"},
		{VecSub, "VSUB"},
		{VecMul, "VMUL"},
		{VecDiv, "VDIV"},
		{VecLoad, "VLOAD"},
		{VecStore, "VSTORE"},
		{VecShuffle, "VSHUFFLE"},
	}

	for _, tt := range tests {
		if got := tt.vo.String(); got != tt.expected {
			t.Errorf("VectorOpcode.String() = %v, expected %v", got, tt.expected)
		}
	}
}

// Test can vectorize operation
func TestCanVectorizeOperation(t *testing.T) {
	tests := []struct {
		op       Opcode
		expected bool
	}{
		{OpAdd, true},
		{OpSub, true},
		{OpMul, true},
		{OpDiv, true},
		{OpLoadVar, true},
		{OpStoreVar, true},
		{OpJump, false},
		{OpCall, false},
	}

	for _, tt := range tests {
		if got := CanVectorizeOperation(tt.op); got != tt.expected {
			t.Errorf("CanVectorizeOperation(%v) = %v, expected %v", tt.op, got, tt.expected)
		}
	}
}

// Test supported vector types
func TestGetSupportedTypes(t *testing.T) {
	types := GetSupportedTypes()

	if len(types) != 4 {
		t.Errorf("expected 4 supported types, got %d", len(types))
	}

	expected := []VectorType{Vec4F, Vec2F64, Vec8I, Vec4I64}
	for i, et := range expected {
		if types[i] != et {
			t.Errorf("type at %d = %v, expected %v", i, types[i], et)
		}
	}
}

// Test estimate speedup
func TestEstimateSpeedup(t *testing.T) {
	speedup := EstimateSpeedup(4)

	if speedup < 3 || speedup > 4 {
		t.Errorf("expected speedup around 3-4 for 4-element vectors, got %.2f", speedup)
	}

	if speedup <= 1 {
		t.Errorf("expected speedup > 1, got %.2f", speedup)
	}
}

// Test statistics structure
func TestSIMDOptimizationStats(t *testing.T) {
	stats := SIMDOptimizationStats{
		VectorizableLoops:  10,
		VectorizedLoops:    8,
		VectorInstructions: 20,
		InstructionsFused:  80,
		MemoryBandwidth:    0.85,
		VectorizationRatio: 0.50,
	}

	if stats.VectorizableLoops != 10 {
		t.Errorf("expected 10 vectorizable loops")
	}

	if stats.MemoryBandwidth != 0.85 {
		t.Errorf("expected 0.85 memory bandwidth")
	}
}

// Test SIMD profile output
func TestSIMDProfileOutput(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0))

	optimizer := NewSIMDOptimizer(program)
	_ = optimizer.OptimizeSIMD()

	profile := optimizer.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 40 {
		t.Errorf("expected detailed profile, got %d chars", len(profile))
	}
}

// Test loop start detection
func TestIsLoopStart(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 100.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLess))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpJumpIfFalse, 10))

	optimizer := NewSIMDOptimizer(program)

	if !optimizer.isLoopStart(0) {
		t.Errorf("expected loop start detection")
	}

	if optimizer.isLoopStart(1) {
		t.Errorf("should not detect loop start at position 1")
	}
}

// Test SIMD with empty program
func TestSIMDEmptyProgram(t *testing.T) {
	program := NewIRProgram()
	optimizer := NewSIMDOptimizer(program)

	err := optimizer.OptimizeSIMD()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := optimizer.GetStatistics()
	if stats.VectorInstructions != 0 {
		t.Errorf("expected 0 vector instructions")
	}
}

// Benchmark SIMD optimization
func BenchmarkSIMDOptimization(b *testing.B) {
	program := NewIRProgram()

	// Create a loop-like pattern
	for i := 0; i < 50; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, 0))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, 1))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpStoreVar, 2))
	}

	optimizer := NewSIMDOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optimizer.OptimizeSIMD()
	}
}

// Benchmark vectorization detection
func BenchmarkDetectVectorizable(b *testing.B) {
	program := NewIRProgram()

	for i := 0; i < 100; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, 0))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, float64(i)))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
	}

	optimizer := NewSIMDOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.detectVectorizablePatterns()
	}
}

// Benchmark memory bandwidth analysis
func BenchmarkMemoryBandwidth(b *testing.B) {
	program := NewIRProgram()

	for i := 0; i < 100; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, i%10))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, (i+1)%10))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpStoreVar, (i+2)%10))
	}

	optimizer := NewSIMDOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.analyzeMemoryBandwidth()
	}
}
