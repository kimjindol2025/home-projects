package compiler

import "fmt"

// Optimizer performs IR optimizations
type Optimizer struct {
	program   *IRProgram
	stats     Stats
	constants map[interface{}]int // Map constant values to indices
}

// NewOptimizer creates a new optimizer
func NewOptimizer(program *IRProgram) *Optimizer {
	return &Optimizer{
		program:   program,
		constants: make(map[interface{}]int),
	}
}

// Optimize applies all optimizations
func (o *Optimizer) Optimize() error {
	if o.program == nil {
		return fmt.Errorf("program is nil")
	}

	o.stats.OriginalSize = len(o.program.Instructions)

	// Apply optimizations in sequence
	o.constantFold()
	o.removeDeadCode()
	o.removeUnusedVariables()
	o.peepholeOptimize()
	o.strengthReduction()

	o.stats.OptimizedSize = len(o.program.Instructions)
	o.stats.RemovedCount = o.stats.OriginalSize - o.stats.OptimizedSize

	return nil
}

// constantFold performs constant folding optimization
func (o *Optimizer) constantFold() {
	instructions := o.program.Instructions
	optimized := []*Instruction{}

	for i := 0; i < len(instructions); i++ {
		ins := instructions[i]

		// Pattern: LOAD_CONST a, LOAD_CONST b, binary_op
		if i+2 < len(instructions) &&
			instructions[i].Opcode == OpLoadConst &&
			instructions[i+1].Opcode == OpLoadConst &&
			isBinaryOp(instructions[i+2].Opcode) {

			valA := instructions[i].Args[0]
			valB := instructions[i+1].Args[0]
			result := foldConstants(valA, valB, instructions[i+2].Opcode)

			if result != nil {
				// Replace with single constant load
				optimized = append(optimized, &Instruction{
					Opcode: OpLoadConst,
					Args:   []interface{}{result},
				})
				i += 2 // Skip next two instructions
				continue
			}
		}

		optimized = append(optimized, ins)
	}

	o.program.Instructions = optimized
}

// isBinaryOp checks if opcode is a binary operation
func isBinaryOp(opcode Opcode) bool {
	return opcode == OpAdd || opcode == OpSub ||
		opcode == OpMul || opcode == OpDiv ||
		opcode == OpMod || opcode == OpEqual ||
		opcode == OpLess || opcode == OpGreater
}

// foldConstants evaluates constant expressions
func foldConstants(a, b interface{}, op Opcode) interface{} {
	aNum, okA := toNumber(a)
	bNum, okB := toNumber(b)

	if !okA || !okB {
		return nil
	}

	switch op {
	case OpAdd:
		return aNum + bNum
	case OpSub:
		return aNum - bNum
	case OpMul:
		return aNum * bNum
	case OpDiv:
		if bNum == 0 {
			return nil
		}
		return aNum / bNum
	case OpMod:
		if bNum == 0 {
			return nil
		}
		return int64(aNum) % int64(bNum)
	}
	return nil
}

// toNumber converts interface to float64
func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

// removeDeadCode removes unreachable code
func (o *Optimizer) removeDeadCode() {
	reachable := make(map[int]bool)
	o.markReachable(reachable, 0)

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
		return
	}

	reachable[start] = true
	ins := o.program.Instructions[start]

	switch ins.Opcode {
	case OpReturn, OpThrow:
		// Terminal instruction
	case OpJump:
		if len(ins.Args) > 0 {
			if target, ok := ins.Args[0].(int); ok {
				o.markReachable(reachable, target)
			}
		}
	case OpJumpIfTrue, OpJumpIfFalse:
		if len(ins.Args) > 0 {
			if target, ok := ins.Args[0].(int); ok {
				o.markReachable(reachable, target)
			}
		}
		o.markReachable(reachable, start+1)
	default:
		o.markReachable(reachable, start+1)
	}
}

// removeUnusedVariables removes stores to variables that are never loaded
func (o *Optimizer) removeUnusedVariables() {
	instructions := o.program.Instructions

	// Track variable usage
	used := make(map[int]bool)
	stored := make(map[int]bool)

	// First pass: find used and stored variables
	for _, ins := range instructions {
		switch ins.Opcode {
		case OpLoadVar:
			if len(ins.Args) > 0 {
				if varIdx, ok := ins.Args[0].(int); ok {
					used[varIdx] = true
				}
			}
		case OpStoreVar:
			if len(ins.Args) > 0 {
				if varIdx, ok := ins.Args[0].(int); ok {
					stored[varIdx] = true
				}
			}
		}
	}

	// Second pass: remove stores to unused variables
	optimized := []*Instruction{}
	for i := 0; i < len(instructions); i++ {
		ins := instructions[i]

		// Remove STORE_VAR x, POP (where x is unused)
		if ins.Opcode == OpStoreVar && len(ins.Args) > 0 {
			if varIdx, ok := ins.Args[0].(int); ok {
				if !used[varIdx] && i+1 < len(instructions) &&
					instructions[i+1].Opcode == OpPop {
					// Skip both STORE_VAR and POP
					i++ // Skip next POP
					continue
				}
			}
		}

		optimized = append(optimized, ins)
	}

	o.program.Instructions = optimized
}

// peepholeOptimize performs peephole optimizations
func (o *Optimizer) peepholeOptimize() {
	instructions := o.program.Instructions
	optimized := []*Instruction{}

	for i := 0; i < len(instructions); i++ {
		ins := instructions[i]
		skip := false

		// Pattern 1: LOAD_CONST, POP (useless)
		if i+1 < len(instructions) &&
			ins.Opcode == OpLoadConst &&
			instructions[i+1].Opcode == OpPop {
			i++ // Skip POP
			skip = true
		}

		// Pattern 2: DUPLICATE, POP, POP (can simplify)
		if i+2 < len(instructions) &&
			ins.Opcode == OpDuplicate &&
			instructions[i+1].Opcode == OpPop &&
			instructions[i+2].Opcode == OpPop {
			i += 2 // Skip both POPs
			skip = true
		}

		// Pattern 3: Multiple POPs -> reduce
		popCount := 0
		j := i
		for j < len(instructions) && instructions[j].Opcode == OpPop {
			popCount++
			j++
		}
		if popCount > 1 && ins.Opcode == OpPop {
			// Replace multiple POPs with single combined instruction
			optimized = append(optimized, &Instruction{
				Opcode: OpPop,
				Args:   []interface{}{popCount},
			})
			i += popCount - 1
			skip = true
		}

		if !skip {
			optimized = append(optimized, ins)
		}
	}

	o.program.Instructions = optimized
}

// strengthReduction replaces expensive operations with cheaper ones
func (o *Optimizer) strengthReduction() {
	instructions := o.program.Instructions
	optimized := []*Instruction{}

	for i := 0; i < len(instructions); i++ {
		ins := instructions[i]

		// Pattern: MUL by 2 -> add itself (x * 2 = x + x)
		if ins.Opcode == OpMul && len(ins.Args) > 0 {
			if val, ok := ins.Args[0].(float64); ok && val == 2 {
				// Replace with cheaper instruction pattern
				optimized = append(optimized, &Instruction{
					Opcode: OpDuplicate,
					Args:   nil,
				})
				optimized = append(optimized, &Instruction{
					Opcode: OpAdd,
					Args:   nil,
				})
				continue
			}
		}

		// Pattern: DIV by 1 -> no-op (just remove)
		if ins.Opcode == OpDiv && len(ins.Args) > 0 {
			if val, ok := ins.Args[0].(float64); ok && val == 1 {
				// Remove division by 1
				continue
			}
		}

		optimized = append(optimized, ins)
	}

	o.program.Instructions = optimized
}

// Stats represents optimization statistics
type Stats struct {
	OriginalSize       int
	OptimizedSize      int
	RemovedCount       int
	ConstantsFolded    int
	DeadCodeRemoved    int
	UnusedVarsRemoved  int
	PeepholeOptimized  int
	StrengthReduced    int
}

// GetStats returns optimization statistics
func (o *Optimizer) GetStats() Stats {
	return o.stats
}
