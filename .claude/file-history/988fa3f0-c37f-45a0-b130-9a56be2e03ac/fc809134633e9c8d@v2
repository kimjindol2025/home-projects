package runtime

import (
	"testing"

	"github.com/freelang-ai/gofree/internal/compiler"
)

// Test VM optimizer creation
func TestNewVMOptimizer(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{},
		Constants:    []interface{}{},
	}
	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)

	if optimizer == nil {
		t.Errorf("expected optimizer, got nil")
	}

	if optimizer.vm != vm {
		t.Errorf("optimizer should reference the VM")
	}
}

// Test instruction optimization
func TestOptimizeInstructions(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadConst, Args: []interface{}{5.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: compiler.OpAdd},
			{Opcode: compiler.OpReturn},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeInstructions()

	stats := optimizer.GetStatistics()

	if stats.TotalInstructions != 4 {
		t.Errorf("expected 4 total instructions, got %d", stats.TotalInstructions)
	}
}

// Test stack depth estimation
func TestOptimizeStackUsage(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadConst, Args: []interface{}{1.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{2.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: compiler.OpAdd},
			{Opcode: compiler.OpReturn},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeStackUsage()

	// Stack should be pre-allocated
	if len(vm.stack) == 0 {
		t.Errorf("expected stack to be pre-allocated")
	}
}

// Test cache frequent opcodes
func TestCacheFrequentOpcodes(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadConst, Args: []interface{}{1.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{2.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{3.0}},
			{Opcode: compiler.OpAdd},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeInstructions() // Need to populate stats first
	optimizer.CacheFrequentOpcodes()

	// Should complete without error
	stats := optimizer.GetStatistics()
	if stats.TotalInstructions != 4 {
		t.Errorf("expected 4 instructions, got %d", stats.TotalInstructions)
	}
}

// Test statistics update
func TestUpdateStatistics(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{},
		Constants:    []interface{}{},
	}
	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)

	stats := optimizer.GetStatistics()
	initialExecuted := stats.ExecutedInstructions

	optimizer.UpdateStatistics()

	stats = optimizer.GetStatistics()
	if stats.ExecutedInstructions <= initialExecuted {
		t.Errorf("expected executed count to increase")
	}
}

// Test memory access optimization
func TestOptimizeMemoryAccess(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadVar, Args: []interface{}{0}},
			{Opcode: compiler.OpStoreVar, Args: []interface{}{0}},
			{Opcode: compiler.OpReturn},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeMemoryAccess()

	// Should complete without error
	if len(program.Instructions) == 0 {
		t.Errorf("expected instructions to remain")
	}
}

// Test array operations optimization
func TestOptimizeArrayOperations(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadVar, Args: []interface{}{0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{0.0}},
			{Opcode: compiler.OpAdd, Args: []interface{}{}},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeArrayOperations()

	// Should complete without error
	if len(program.Instructions) != 3 {
		t.Errorf("expected 3 instructions")
	}
}

// Test profile output
func TestProfile(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadConst, Args: []interface{}{5.0}},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.OptimizeInstructions()
	optimizer.UpdateStatistics()

	profile := optimizer.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 20 {
		t.Errorf("expected detailed profile output")
	}
}

// Test VM stats structure
func TestVMStats(t *testing.T) {
	stats := VMStats{
		TotalInstructions:    100,
		ExecutedInstructions: 50,
		StackOperations:      200,
		FunctionCalls:        5,
		CacheHits:            80,
		CacheMisses:          20,
		AverageStackDepth:    3.5,
	}

	if stats.TotalInstructions != 100 {
		t.Errorf("expected 100 total instructions")
	}

	if stats.AverageStackDepth != 3.5 {
		t.Errorf("expected average depth 3.5")
	}
}

// Test unroll loops
func TestUnrollLoops(t *testing.T) {
	program := &compiler.IRProgram{
		Instructions: []*compiler.Instruction{
			{Opcode: compiler.OpLoadConst, Args: []interface{}{0.0}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{10.0}},
			{Opcode: compiler.OpLess},
			{Opcode: compiler.OpJumpIfFalse, Args: []interface{}{10}},
			{Opcode: compiler.OpLoadConst, Args: []interface{}{1.0}},
			{Opcode: compiler.OpAdd},
			{Opcode: compiler.OpJump, Args: []interface{}{1}},
		},
		Constants: []interface{}{},
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)
	optimizer.UnrollLoops()

	// Should complete without error
	if len(program.Instructions) == 0 {
		t.Errorf("expected instructions")
	}
}

// Benchmark VM optimization
func BenchmarkVMOptimization(b *testing.B) {
	program := &compiler.IRProgram{
		Instructions: make([]*compiler.Instruction, 100),
		Constants:    []interface{}{},
	}

	for i := 0; i < 100; i++ {
		program.Instructions[i] = &compiler.Instruction{
			Opcode: compiler.OpLoadConst,
			Args:   []interface{}{float64(i)},
		}
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.OptimizeStackUsage()
		optimizer.CacheFrequentOpcodes()
		optimizer.UnrollLoops()
	}
}

// Benchmark stack estimation
func BenchmarkEstimateStackDepth(b *testing.B) {
	program := &compiler.IRProgram{
		Instructions: make([]*compiler.Instruction, 1000),
		Constants:    []interface{}{},
	}

	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			program.Instructions[i] = &compiler.Instruction{
				Opcode: compiler.OpLoadConst,
				Args:   []interface{}{float64(i)},
			}
		} else {
			program.Instructions[i] = &compiler.Instruction{
				Opcode: compiler.OpAdd,
			}
		}
	}

	vm := NewVM(program)
	optimizer := NewVMOptimizer(vm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optimizer.estimateStackDepth()
	}
}
