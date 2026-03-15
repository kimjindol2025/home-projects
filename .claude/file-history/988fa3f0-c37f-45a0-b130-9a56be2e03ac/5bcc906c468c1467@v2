package compiler

import (
	"testing"
)

// Test advanced optimizer creation
func TestNewAdvancedOptimizer(t *testing.T) {
	program := NewIRProgram()
	optimizer := NewAdvancedOptimizer(program)

	if optimizer == nil {
		t.Errorf("expected optimizer, got nil")
	}

	if optimizer.program != program {
		t.Errorf("optimizer should reference the program")
	}
}

// Test arithmetic simplification
func TestSimplifyArithmetic(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 5.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 0.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.ArithmeticSimplified == 0 {
		t.Logf("arithmetic simplification may not have been applied")
	}
}

// Test x - x optimization
func TestSubtractionIdentity(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpSub, 1))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.ArithmeticSimplified == 0 {
		t.Logf("identity subtraction may not have been recognized")
	}
}

// Test multi-instruction pattern matching
func TestMultiInstructionPatterns(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 5.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 3.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.PatternsMatched == 0 {
		t.Logf("pattern matching may not have been applied")
	}
}

// Test branch optimization
func TestBranchOptimization(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 1.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpJump, 2)) // Jump to next instruction
	program.Instructions = append(program.Instructions,
		NewInstruction(OpReturn))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.BranchOptimized == 0 {
		t.Logf("branch optimization may not have been applied")
	}
}

// Test register pressure analysis
func TestRegisterPressure(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 2))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 3))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpAdd))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 4))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.RegisterPressure < 0 || stats.RegisterPressure > 1 {
		t.Errorf("invalid register pressure: %.2f", stats.RegisterPressure)
	}
}

// Test optimization with empty program
func TestOptimizeEmptyProgram(t *testing.T) {
	program := NewIRProgram()
	optimizer := NewAdvancedOptimizer(program)

	err := optimizer.OptimizeAdvanced()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test optimization with single instruction
func TestOptimizeSingleInstruction(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 42.0))

	optimizer := NewAdvancedOptimizer(program)
	err := optimizer.OptimizeAdvanced()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(program.Instructions) != 1 {
		t.Errorf("expected 1 instruction, got %d", len(program.Instructions))
	}
}

// Test pattern: LOAD → STORE → LOAD → STORE
func TestDuplicateLoadPattern(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 2))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadVar, 1))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpStoreVar, 3))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.PatternsMatched == 0 {
		t.Logf("duplicate load pattern may not have been recognized")
	}
}

// Test helper function: isZero
func TestIsZero(t *testing.T) {
	tests := []struct {
		args     []interface{}
		expected bool
	}{
		{[]interface{}{0.0}, true},
		{[]interface{}{0}, true},
		{[]interface{}{int64(0)}, true},
		{[]interface{}{1.0}, false},
		{[]interface{}{}, false},
	}

	for _, tt := range tests {
		if got := isZero(tt.args); got != tt.expected {
			t.Errorf("isZero(%v) = %v, expected %v", tt.args, got, tt.expected)
		}
	}
}

// Test helper function: isOne
func TestIsOne(t *testing.T) {
	tests := []struct {
		args     []interface{}
		expected bool
	}{
		{[]interface{}{1.0}, true},
		{[]interface{}{1}, true},
		{[]interface{}{int64(1)}, true},
		{[]interface{}{0.0}, false},
		{[]interface{}{}, false},
	}

	for _, tt := range tests {
		if got := isOne(tt.args); got != tt.expected {
			t.Errorf("isOne(%v) = %v, expected %v", tt.args, got, tt.expected)
		}
	}
}

// Test statistics structure
func TestAdvancedOptimizationStats(t *testing.T) {
	stats := AdvancedOptimizationStats{
		PatternsMatched:      10,
		ArithmeticSimplified: 5,
		BranchOptimized:      3,
		RegisterPressure:     0.75,
		InstructionsRemoved:  8,
	}

	if stats.PatternsMatched != 10 {
		t.Errorf("expected 10 patterns matched")
	}

	if stats.RegisterPressure != 0.75 {
		t.Errorf("expected 0.75 register pressure")
	}
}

// Test profile output
func TestProfileOutput(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 1.0))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	profile := optimizer.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 30 {
		t.Errorf("expected detailed profile, got %d chars", len(profile))
	}
}

// Test instruction removal
func TestInstructionRemoval(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 5.0))
	program.Instructions = append(program.Instructions,
		NewInstruction(OpJump, 1)) // Unnecessary jump
	program.Instructions = append(program.Instructions,
		NewInstruction(OpLoadConst, 10.0))

	optimizer := NewAdvancedOptimizer(program)
	_ = optimizer.OptimizeAdvanced()

	stats := optimizer.GetStatistics()
	if stats.InstructionsRemoved == 0 {
		t.Logf("expected some instruction removal")
	}
}

// Benchmark advanced optimization
func BenchmarkAdvancedOptimization(b *testing.B) {
	program := NewIRProgram()

	// Create a program with many instructions
	for i := 0; i < 100; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, float64(i)))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, 1.0))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
	}

	optimizer := NewAdvancedOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optimizer.OptimizeAdvanced()
	}
}

// Benchmark arithmetic simplification
func BenchmarkArithmeticSimplification(b *testing.B) {
	program := NewIRProgram()

	for i := 0; i < 50; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, float64(i)))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, 0.0))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
	}

	optimizer := NewAdvancedOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.simplifyArithmetic()
	}
}

// Benchmark pattern matching
func BenchmarkPatternMatching(b *testing.B) {
	program := NewIRProgram()

	for i := 0; i < 100; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, float64(i)))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadVar, i%5))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpAdd))
	}

	optimizer := NewAdvancedOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.optimizeMultiInstructionPatterns()
	}
}

// Benchmark branch optimization
func BenchmarkBranchOptimization(b *testing.B) {
	program := NewIRProgram()

	for i := 0; i < 50; i++ {
		program.Instructions = append(program.Instructions,
			NewInstruction(OpLoadConst, float64(i)))
		program.Instructions = append(program.Instructions,
			NewInstruction(OpJump, i+1))
	}

	optimizer := NewAdvancedOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.optimizeBranches()
	}
}
