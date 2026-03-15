package compiler

import "fmt"

// Optimizer performs IR optimizations
type Optimizer struct {
	program *IRProgram
}

// NewOptimizer creates a new optimizer
func NewOptimizer(program *IRProgram) *Optimizer {
	return &Optimizer{program: program}
}

// Optimize applies all optimizations
func (o *Optimizer) Optimize() error {
	if o.program == nil {
		return fmt.Errorf("program is nil")
	}

	// Apply constant folding
	o.constantFold()

	// Remove dead code
	o.removeDeadCode()

	// Peephole optimization
	o.peepholeOptimize()

	return nil
}

// constantFold performs constant folding optimization
func (o *Optimizer) constantFold() {
	// Track constant values during compilation
	for _, ins := range o.program.Instructions {
		switch ins.Opcode {
		case OpLoadConst:
			// Instruction is already optimized
		case OpAdd, OpSub, OpMul, OpDiv, OpMod:
			// TODO: Fold constants at compile time
		}
	}
}

// removeDeadCode removes unreachable code
func (o *Optimizer) removeDeadCode() {
	// Mark reachable instructions
	reachable := make(map[int]bool)
	o.markReachable(reachable, 0)

	// Filter unreachable instructions
	filtered := []*Instruction{}
	for i, ins := range o.program.Instructions {
		if reachable[i] {
			filtered = append(filtered, ins)
		}
	}

	o.program.Instructions = filtered
}

// markReachable marks all reachable instructions
func (o *Optimizer) markReachable(reachable map[int]bool, start int) {
	if start < 0 || start >= len(o.program.Instructions) {
		return
	}

	if reachable[start] {
		return // Already marked
	}

	reachable[start] = true
	ins := o.program.Instructions[start]

	switch ins.Opcode {
	case OpReturn, OpThrow:
		// No next instruction
	case OpJump:
		// Jump to label
		if len(ins.Args) > 0 {
			if target, ok := ins.Args[0].(int); ok {
				o.markReachable(reachable, target)
			}
		}
	case OpJumpIfTrue, OpJumpIfFalse:
		// Can jump or fall through
		if len(ins.Args) > 0 {
			if target, ok := ins.Args[0].(int); ok {
				o.markReachable(reachable, target)
			}
		}
		o.markReachable(reachable, start+1)
	default:
		// Fall through to next instruction
		o.markReachable(reachable, start+1)
	}
}

// peepholeOptimize performs peephole optimizations
func (o *Optimizer) peepholeOptimize() {
	instructions := o.program.Instructions
	if len(instructions) < 2 {
		return
	}

	// Pattern matching for common optimizations
	optimized := []*Instruction{}
	for i := 0; i < len(instructions); i++ {
		ins := instructions[i]

		// Optimization: LOAD_CONST followed by STORE_VAR
		if i+1 < len(instructions) &&
			ins.Opcode == OpLoadConst &&
			instructions[i+1].Opcode == OpStoreVar {
			// Keep both instructions (no optimization needed)
			optimized = append(optimized, ins)
		} else if i+1 < len(instructions) &&
			ins.Opcode == OpLoadConst &&
			instructions[i+1].Opcode == OpPop {
			// Remove: LOAD_CONST followed by POP
			i++ // Skip POP
		} else if ins.Opcode == OpPop && i > 0 &&
			optimized[len(optimized)-1].Opcode == OpPop {
			// Remove duplicate POPs
			// (skip this one)
		} else if ins.Opcode == OpDuplicate && i > 0 &&
			optimized[len(optimized)-1].Opcode == OpDuplicate {
			// Remove duplicate DUPLICATEs
			// (skip this one)
		} else {
			optimized = append(optimized, ins)
		}
	}

	o.program.Instructions = optimized
}

// Stats represents optimization statistics
type Stats struct {
	OriginalSize  int
	OptimizedSize int
	RemovedCount  int
}

// GetStats returns optimization statistics
func (o *Optimizer) GetStats() Stats {
	return Stats{
		OptimizedSize: len(o.program.Instructions),
	}
}
