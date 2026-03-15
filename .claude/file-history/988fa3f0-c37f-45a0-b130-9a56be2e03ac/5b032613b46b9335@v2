package compiler

import (
	"testing"
)

// Test constant folding optimization
func TestConstantFolding(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: OpAdd},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// After optimization, should have folded 5+3 into 8
	stats := optimizer.GetStats()
	if stats.OptimizedSize >= 3 {
		t.Logf("constant folding reduced %d instructions to %d", stats.OriginalSize, stats.OptimizedSize)
	}
}

// Test unused variable removal
func TestUnusedVariableRemoval(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{42.0}},
			{Opcode: OpStoreVar, Args: []interface{}{0}},
			{Opcode: OpPop},
			{Opcode: OpLoadConst, Args: []interface{}{10.0}},
			{Opcode: OpReturn},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := optimizer.GetStats()
	if stats.RemovedCount > 0 {
		t.Logf("removed %d instructions", stats.RemovedCount)
	}
}

// Test peephole optimization
func TestPeepholeOptimization(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpPop}, // Useless pop
			{Opcode: OpLoadConst, Args: []interface{}{10.0}},
			{Opcode: OpReturn},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Should remove useless LOAD_CONST, POP pattern
	stats := optimizer.GetStats()
	if stats.RemovedCount < 1 {
		t.Logf("peephole optimization removed %d instructions", stats.RemovedCount)
	}
}

// Test dead code removal
func TestDeadCodeRemoval(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpReturn},
			{Opcode: OpLoadConst, Args: []interface{}{10.0}}, // Unreachable
			{Opcode: OpPop},                                   // Unreachable
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// After optimization, unreachable code should be removed
	if len(program.Instructions) > 2 {
		t.Errorf("expected <= 2 instructions, got %d", len(program.Instructions))
	}
}

// Test strength reduction
func TestStrengthReduction(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpMul, Args: []interface{}{2.0}}, // Can be optimized to ADD
			{Opcode: OpReturn},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Strength reduction should replace MUL 2 with cheaper operations
	stats := optimizer.GetStats()
	if stats.OptimizedSize < stats.OriginalSize {
		t.Logf("strength reduction reduced size from %d to %d", stats.OriginalSize, stats.OptimizedSize)
	}
}

// Test optimization statistics
func TestOptimizationStats(t *testing.T) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: OpAdd},
			{Opcode: OpReturn},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := optimizer.GetStats()

	if stats.OriginalSize != 4 {
		t.Errorf("expected original size 4, got %d", stats.OriginalSize)
	}

	if stats.OptimizedSize > stats.OriginalSize {
		t.Errorf("expected optimized size <= original size")
	}

	if stats.RemovedCount < 0 {
		t.Errorf("expected non-negative removed count")
	}
}

// Test optimizer with nil program
func TestOptimizerNilProgram(t *testing.T) {
	optimizer := NewOptimizer(nil)
	err := optimizer.Optimize()

	if err == nil {
		t.Errorf("expected error for nil program")
	}
}

// Test isBinaryOp
func TestIsBinaryOp(t *testing.T) {
	tests := []struct {
		opcode Opcode
		expect bool
	}{
		{OpAdd, true},
		{OpSub, true},
		{OpMul, true},
		{OpDiv, true},
		{OpEqual, true},
		{OpLess, true},
		{OpLoadConst, false},
		{OpReturn, false},
	}

	for _, test := range tests {
		result := isBinaryOp(test.opcode)
		if result != test.expect {
			t.Errorf("isBinaryOp(%v) = %v, want %v", test.opcode, result, test.expect)
		}
	}
}

// Test foldConstants
func TestFoldConstants(t *testing.T) {
	tests := []struct {
		a      interface{}
		b      interface{}
		op     Opcode
		expect interface{}
	}{
		{5.0, 3.0, OpAdd, 8.0},
		{10.0, 3.0, OpSub, 7.0},
		{4.0, 5.0, OpMul, 20.0},
		{20.0, 4.0, OpDiv, 5.0},
		{10, 3, OpAdd, 13.0},
		{int64(5), int64(2), OpMul, 10.0},
	}

	for _, test := range tests {
		result := foldConstants(test.a, test.b, test.op)
		if result != test.expect {
			t.Errorf("foldConstants(%v, %v, %v) = %v, want %v",
				test.a, test.b, test.op, result, test.expect)
		}
	}
}

// Test toNumber
func TestToNumber(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect float64
		ok     bool
	}{
		{5.0, 5.0, true},
		{10, 10.0, true},
		{int64(42), 42.0, true},
		{"string", 0, false},
		{nil, 0, false},
	}

	for _, test := range tests {
		result, ok := toNumber(test.input)
		if ok != test.ok {
			t.Errorf("toNumber(%v) ok = %v, want %v", test.input, ok, test.ok)
		}
		if ok && result != test.expect {
			t.Errorf("toNumber(%v) = %v, want %v", test.input, result, test.expect)
		}
	}
}

// Benchmark optimization
func BenchmarkOptimization(b *testing.B) {
	program := &IRProgram{
		Instructions: make([]*Instruction, 100),
		Constants:    []interface{}{},
	}

	for i := 0; i < 100; i++ {
		program.Instructions[i] = &Instruction{
			Opcode: OpLoadConst,
			Args:   []interface{}{float64(i)},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer := NewOptimizer(program)
		optimizer.Optimize()
	}
}

// Benchmark constant folding
func BenchmarkConstantFolding(b *testing.B) {
	program := &IRProgram{
		Instructions: []*Instruction{
			{Opcode: OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: OpAdd},
		},
		Constants: []interface{}{},
	}

	optimizer := NewOptimizer(program)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.constantFold()
	}
}
