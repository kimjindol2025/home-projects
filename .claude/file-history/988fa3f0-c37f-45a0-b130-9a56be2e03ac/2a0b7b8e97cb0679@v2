package compiler

import (
	"fmt"
)

// AdvancedOptimizer performs sophisticated instruction pattern recognition
type AdvancedOptimizer struct {
	program *IRProgram
	stats   AdvancedOptimizationStats
}

// AdvancedOptimizationStats tracks optimization metrics
type AdvancedOptimizationStats struct {
	PatternsMatched      int64
	ArithmeticSimplified int64
	BranchOptimized      int64
	RegisterPressure     float64 // 0-1 scale
	InstructionsRemoved  int64
}

// NewAdvancedOptimizer creates an advanced optimizer
func NewAdvancedOptimizer(program *IRProgram) *AdvancedOptimizer {
	return &AdvancedOptimizer{
		program: program,
	}
}

// OptimizeAdvanced performs all advanced optimizations
func (ao *AdvancedOptimizer) OptimizeAdvanced() error {
	if ao.program == nil || len(ao.program.Instructions) == 0 {
		return nil
	}

	// Apply optimizations in sequence
	ao.simplifyArithmetic()
	ao.optimizeMultiInstructionPatterns()
	ao.optimizeBranches()
	ao.analyzeRegisterPressure()

	return nil
}

// simplifyArithmetic eliminates algebraic identity operations
func (ao *AdvancedOptimizer) simplifyArithmetic() {
	if len(ao.program.Instructions) < 2 {
		return
	}

	newInstructions := []*Instruction{}

	for i := 0; i < len(ao.program.Instructions); i++ {
		ins := ao.program.Instructions[i]

		// Pattern: LOAD_CONST 0 → ADD
		// Becomes: Skip ADD, use first operand
		if ins.Opcode == OpAdd && i >= 2 {
			prev := ao.program.Instructions[i-1]

			// Check if loading 0 and adding
			if prev.Opcode == OpLoadConst && isZero(prev.Args) {
				// Skip the zero add
				ao.stats.ArithmeticSimplified++
				continue
			}

			// Check if multiplying by 1
			if ins.Opcode == OpMul && isOne(prev.Args) {
				ao.stats.ArithmeticSimplified++
				continue
			}
		}

		// Pattern: x - x → LOAD_CONST 0
		if ins.Opcode == OpSub && i >= 1 {
			prev := ao.program.Instructions[i-1]
			if prev.Opcode == OpLoadVar && len(ins.Args) > 0 && len(prev.Args) > 0 {
				if prev.Args[0] == ins.Args[0] {
					// Same variable, result is 0
					newInstructions = append(newInstructions,
						NewInstruction(OpLoadConst, 0))
					ao.stats.ArithmeticSimplified++
					continue
				}
			}
		}

		// Pattern: x * 2 → x << 1
		if ins.Opcode == OpMul && len(ins.Args) > 0 {
			if multiplier, ok := ins.Args[0].(float64); ok && multiplier == 2.0 {
				// Replace with bit shift (simulated)
				newInstructions = append(newInstructions,
					NewInstruction(OpLoadConst, int64(1)))
				newInstructions = append(newInstructions,
					NewInstruction(OpShiftLeft)) // Pseudo-opcode
				ao.stats.ArithmeticSimplified++
				continue
			}
		}

		newInstructions = append(newInstructions, ins)
	}

	if len(newInstructions) < len(ao.program.Instructions) {
		ao.stats.InstructionsRemoved += int64(len(ao.program.Instructions) - len(newInstructions))
		ao.program.Instructions = newInstructions
	}
}

// optimizeMultiInstructionPatterns recognizes sequences of instructions
func (ao *AdvancedOptimizer) optimizeMultiInstructionPatterns() {
	if len(ao.program.Instructions) < 3 {
		return
	}

	newInstructions := []*Instruction{}
	i := 0

	for i < len(ao.program.Instructions) {
		ins := ao.program.Instructions[i]

		// Pattern: LOAD a → LOAD b → ADD
		// Check if we can fuse them
		if ins.Opcode == OpLoadConst && i+2 < len(ao.program.Instructions) {
			next1 := ao.program.Instructions[i+1]
			next2 := ao.program.Instructions[i+2]

			if next1.Opcode == OpLoadConst && next2.Opcode == OpAdd {
				// Can fuse: fold the addition at compile time
				if a, ok := ins.Args[0].(float64); ok {
					if b, ok := next1.Args[0].(float64); ok {
						result := a + b
						newInstructions = append(newInstructions,
							NewInstruction(OpLoadConst, result))
						ao.stats.PatternsMatched++
						i += 3
						continue
					}
				}
			}
		}

		// Pattern: LOAD → STORE → LOAD (same var) → STORE (different var)
		// Can be optimized to: LOAD → STORE → DUP (instead of LOAD again)
		if ins.Opcode == OpLoadVar && i+3 < len(ao.program.Instructions) {
			next1 := ao.program.Instructions[i+1]
			next2 := ao.program.Instructions[i+2]
			next3 := ao.program.Instructions[i+3]

			if next1.Opcode == OpStoreVar &&
				next2.Opcode == OpLoadVar &&
				next3.Opcode == OpStoreVar {

				if ins.Args[0] == next2.Args[0] { // Same variable loaded twice
					// Replace second LOAD with DUP
					newInstructions = append(newInstructions, ins)
					newInstructions = append(newInstructions, next1)
					newInstructions = append(newInstructions,
						NewInstruction(OpDuplicate))
					newInstructions = append(newInstructions, next3)
					ao.stats.PatternsMatched++
					i += 4
					continue
				}
			}
		}

		newInstructions = append(newInstructions, ins)
		i++
	}

	ao.program.Instructions = newInstructions
}

// optimizeBranches removes unnecessary jumps
func (ao *AdvancedOptimizer) optimizeBranches() {
	if len(ao.program.Instructions) < 2 {
		return
	}

	newInstructions := []*Instruction{}

	for i := 0; i < len(ao.program.Instructions); i++ {
		ins := ao.program.Instructions[i]

		// Pattern: JMP x → (next instruction is already at x)
		if (ins.Opcode == OpJump || ins.Opcode == OpJumpIfFalse) &&
			len(ins.Args) > 0 {
			target := -1
			if t, ok := ins.Args[0].(int); ok {
				target = t
			}

			// Check if jump target is immediately after
			if target == i+1 {
				ao.stats.BranchOptimized++
				continue // Skip unnecessary jump
			}
		}

		// Pattern: JMP x → JMP y → simplify to JMP y
		if ins.Opcode == OpJump && i+1 < len(ao.program.Instructions) {
			next := ao.program.Instructions[i+1]
			if next.Opcode == OpJump && len(next.Args) > 0 {
				// Jump to a jump, redirect to final target
				if target, ok := next.Args[0].(int); ok {
					newInstructions = append(newInstructions,
						NewInstruction(OpJump, target))
					ao.stats.BranchOptimized++
					i++ // Skip the intermediate jump
					continue
				}
			}
		}

		newInstructions = append(newInstructions, ins)
	}

	ao.program.Instructions = newInstructions
}

// analyzeRegisterPressure calculates how many variables are live simultaneously
func (ao *AdvancedOptimizer) analyzeRegisterPressure() {
	if len(ao.program.Instructions) == 0 {
		ao.stats.RegisterPressure = 0
		return
	}

	maxLive := 0
	currentLive := 0

	// Simple heuristic: count variable loads/stores
	for _, ins := range ao.program.Instructions {
		switch ins.Opcode {
		case OpLoadVar:
			currentLive++
			if currentLive > maxLive {
				maxLive = currentLive
			}

		case OpStoreVar:
			if currentLive > 0 {
				currentLive--
			}
		}
	}

	// Normalize to 0-1 scale (assuming 16 registers)
	ao.stats.RegisterPressure = float64(maxLive) / 16.0
	if ao.stats.RegisterPressure > 1.0 {
		ao.stats.RegisterPressure = 1.0
	}
}

// GetStatistics returns optimization statistics
func (ao *AdvancedOptimizer) GetStatistics() AdvancedOptimizationStats {
	return ao.stats
}

// Profile returns a profile of optimizations applied
func (ao *AdvancedOptimizer) Profile() string {
	return fmt.Sprintf(`Advanced Optimization Profile:
  Patterns Matched: %d
  Arithmetic Simplified: %d
  Branches Optimized: %d
  Instructions Removed: %d
  Register Pressure: %.2f
`,
		ao.stats.PatternsMatched,
		ao.stats.ArithmeticSimplified,
		ao.stats.BranchOptimized,
		ao.stats.InstructionsRemoved,
		ao.stats.RegisterPressure)
}

// Helper functions

// isZero checks if argument is zero
func isZero(args []interface{}) bool {
	if len(args) == 0 {
		return false
	}

	switch v := args[0].(type) {
	case float64:
		return v == 0
	case int:
		return v == 0
	case int64:
		return v == 0
	}

	return false
}

// isOne checks if argument is one
func isOne(args []interface{}) bool {
	if len(args) == 0 {
		return false
	}

	switch v := args[0].(type) {
	case float64:
		return v == 1
	case int:
		return v == 1
	case int64:
		return v == 1
	}

	return false
}

// OpShiftLeft is a pseudo-opcode for bit shift (not in standard set)
const OpShiftLeft = Opcode(999)
