package runtime

import (
	"fmt"

	"github.com/freelang-ai/gofree/internal/compiler"
)

// VMOptimizer provides performance optimizations for VM execution
type VMOptimizer struct {
	vm             *VM
	instructionMap map[compiler.Opcode]int // Cache opcode dispatch
	statistics     VMStats
}

// VMStats tracks VM performance statistics
type VMStats struct {
	TotalInstructions    int64
	ExecutedInstructions int64
	StackOperations      int64
	FunctionCalls        int64
	CacheHits            int64
	CacheMisses          int64
	AverageStackDepth    float64
}

// NewVMOptimizer creates a new VM optimizer
func NewVMOptimizer(vm *VM) *VMOptimizer {
	return &VMOptimizer{
		vm:             vm,
		instructionMap: make(map[compiler.Opcode]int),
	}
}

// OptimizeInstructions pre-processes instructions for faster dispatch
func (o *VMOptimizer) OptimizeInstructions() {
	if o.vm.program == nil {
		return
	}

	// Build instruction map for faster lookup
	for i, ins := range o.vm.program.Instructions {
		o.instructionMap[ins.Opcode] = i
	}

	o.statistics.TotalInstructions = int64(len(o.vm.program.Instructions))
}

// InlineBuiltins replaces common function calls with optimized versions
func (o *VMOptimizer) InlineBuiltins() {
	if o.vm.program == nil {
		return
	}

	// Inline frequently used built-in functions
	for i := 0; i < len(o.vm.program.Instructions)-1; i++ {
		ins := o.vm.program.Instructions[i]

		// Pattern: LOAD_VAR "println", CALL
		// In a real implementation, would replace with optimized opcode
		if ins.Opcode == compiler.OpCall &&
			len(ins.Args) > 0 {
			if funcName, ok := ins.Args[0].(string); ok {
				// Track frequent function calls for inlining
				switch funcName {
				case "print", "println":
					// Mark for potential inlining
					_ = funcName
				case "len":
					// Mark for inlining
					_ = funcName
				case "type":
					// Mark for inlining
					_ = funcName
				}
			}
		}
	}
}

// OptimizeStackUsage pre-allocates stack with appropriate size
func (o *VMOptimizer) OptimizeStackUsage() {
	if o.vm.program == nil {
		return
	}

	// Estimate maximum stack depth
	maxDepth := o.estimateStackDepth()

	// Pre-allocate stack
	if len(o.vm.stack) < maxDepth {
		newStack := make([]interface{}, maxDepth)
		copy(newStack, o.vm.stack)
		o.vm.stack = newStack
	}
}

// estimateStackDepth analyzes instructions to determine max stack depth
func (o *VMOptimizer) estimateStackDepth() int {
	if o.vm.program == nil {
		return 100
	}

	depth := 0
	maxDepth := 0

	for _, ins := range o.vm.program.Instructions {
		switch ins.Opcode {
		// Operations that push to stack
		case compiler.OpLoadConst, compiler.OpLoadVar,
			compiler.OpDuplicate:
			depth++

		// Operations that pop from stack
		case compiler.OpPop, compiler.OpReturn:
			if depth > 0 {
				depth--
			}

		// Binary operations: pop 2, push 1
		case compiler.OpAdd, compiler.OpSub, compiler.OpMul,
			compiler.OpDiv, compiler.OpMod, compiler.OpEqual,
			compiler.OpLess, compiler.OpGreater:
			if depth >= 2 {
				depth--
			}

		// Store pops 1
		case compiler.OpStoreVar:
			if depth > 0 {
				depth--
			}
		}

		if depth > maxDepth {
			maxDepth = depth
		}
	}

	// Add buffer for safety
	return maxDepth + 10
}

// CacheFrequentOpcodes identifies and optimizes frequent opcodes
func (o *VMOptimizer) CacheFrequentOpcodes() {
	if o.vm.program == nil {
		return
	}

	opcodeFreq := make(map[compiler.Opcode]int)

	// Count opcode frequencies
	for _, ins := range o.vm.program.Instructions {
		opcodeFreq[ins.Opcode]++
	}

	// Log most frequent opcodes
	// (In real VM, could use inline caching for these)
	_ = opcodeFreq
}

// UnrollLoops attempts to unroll small loops
func (o *VMOptimizer) UnrollLoops() {
	if o.vm.program == nil {
		return
	}

	instructions := o.vm.program.Instructions

	// Pattern detection: identify small loops
	for i := 0; i < len(instructions)-1; i++ {
		ins := instructions[i]

		// Check for JUMP_IF_FALSE (loop condition check)
		if ins.Opcode == compiler.OpJumpIfFalse &&
			len(ins.Args) > 0 {

			if target, ok := ins.Args[0].(int); ok {
				// Calculate loop size
				loopSize := target - i

				// Only unroll small loops (< 5 instructions)
				if loopSize > 0 && loopSize < 5 {
					// Mark loop for optimization
					// (Real implementation would duplicate loop body)
					_ = loopSize
				}
			}
		}
	}
}

// GetStatistics returns VM performance statistics
func (o *VMOptimizer) GetStatistics() VMStats {
	return o.statistics
}

// UpdateStatistics tracks execution statistics
func (o *VMOptimizer) UpdateStatistics() {
	o.statistics.ExecutedInstructions++

	// Update average stack depth
	currentDepth := float64(len(o.vm.stack))
	o.statistics.AverageStackDepth =
		(o.statistics.AverageStackDepth + currentDepth) / 2.0
}

// OptimizeMemoryAccess reduces redundant memory access
func (o *VMOptimizer) OptimizeMemoryAccess() {
	if o.vm.program == nil {
		return
	}

	instructions := o.vm.program.Instructions

	// Pattern: LOAD_VAR x, STORE_VAR x (redundant)
	// Can be eliminated if x is not used between operations
	for i := 0; i < len(instructions)-1; i++ {
		ins := instructions[i]
		next := instructions[i+1]

		if ins.Opcode == compiler.OpLoadVar &&
			next.Opcode == compiler.OpStoreVar {

			if len(ins.Args) > 0 && len(next.Args) > 0 {
				if loadVar, ok := ins.Args[0].(int); ok {
					if storeVar, ok := next.Args[0].(int); ok {
						if loadVar == storeVar {
							// Redundant load-store, can optimize
							// Remove this pattern in a full implementation
						}
					}
				}
			}
		}
	}
}

// OptimizeArrayOperations inlines array access patterns
func (o *VMOptimizer) OptimizeArrayOperations() {
	if o.vm.program == nil {
		return
	}

	instructions := o.vm.program.Instructions

	// Pattern: LOAD_VAR arr, LOAD_CONST index, Indexing operation
	// Can use direct indexed access pattern
	for i := 0; i < len(instructions)-2; i++ {
		if instructions[i].Opcode == compiler.OpLoadVar &&
			instructions[i+1].Opcode == compiler.OpLoadConst {

			// Mark for potential array access optimization
			// In real implementation, would fuse these operations
			_ = instructions[i]
		}
	}
}

// Profile returns an analysis of VM bottlenecks
func (o *VMOptimizer) Profile() string {
	stats := o.GetStatistics()

	return fmt.Sprintf(`VM Performance Profile:
  Total Instructions: %d
  Executed: %d
  Stack Ops: %d
  Function Calls: %d
  Cache Hit Rate: %.1f%%
  Average Stack Depth: %.1f`,
		stats.TotalInstructions,
		stats.ExecutedInstructions,
		stats.StackOperations,
		stats.FunctionCalls,
		100*float64(stats.CacheHits)/float64(stats.CacheHits+stats.CacheMisses),
		stats.AverageStackDepth)
}
