package compiler

import (
	"fmt"
)

// SIMDOptimizer performs automatic vectorization of numerical operations
type SIMDOptimizer struct {
	program *IRProgram
	stats   SIMDOptimizationStats
}

// SIMDOptimizationStats tracks SIMD optimization metrics
type SIMDOptimizationStats struct {
	VectorizableLoops   int64
	VectorizedLoops     int64
	InstructionsFused   int64
	VectorInstructions  int64
	MemoryBandwidth     float64 // bytes/cycle
	VectorizationRatio  float64 // 0-1 scale
}

// VectorType represents SIMD vector types
type VectorType int

const (
	Vec4F   VectorType = iota // 4 floats
	Vec2F64                    // 2 doubles
	Vec8I   // 8 ints
	Vec4I64 // 4 int64s
)

// VectorOpcode represents SIMD operations
type VectorOpcode int

const (
	VecAdd VectorOpcode = iota
	VecSub
	VecMul
	VecDiv
	VecLoad
	VecStore
	VecShuffle
)

// NewSIMDOptimizer creates a new SIMD optimizer
func NewSIMDOptimizer(program *IRProgram) *SIMDOptimizer {
	return &SIMDOptimizer{
		program: program,
	}
}

// OptimizeSIMD performs SIMD optimizations
func (so *SIMDOptimizer) OptimizeSIMD() error {
	if so.program == nil || len(so.program.Instructions) == 0 {
		return nil
	}

	so.detectVectorizablePatterns()
	so.emitVectorInstructions()
	so.analyzeMemoryBandwidth()

	return nil
}

// detectVectorizablePatterns identifies loops suitable for vectorization
func (so *SIMDOptimizer) detectVectorizablePatterns() {
	if len(so.program.Instructions) < 5 {
		return
	}

	i := 0
	for i < len(so.program.Instructions)-4 {
		// Pattern: Loop with uniform array operations
		// LOAD loop_counter, LOAD limit, CMP, JMP_IF_FALSE, ... LOAD array[i], op, STORE array[i], JUMP

		if so.isLoopStart(i) {
			loopEnd := so.findLoopEnd(i)
			if loopEnd > i && so.canVectorize(i, loopEnd) {
				so.stats.VectorizableLoops++
				i = loopEnd
				continue
			}
		}

		i++
	}
}

// isLoopStart checks if instruction at index i starts a loop
func (so *SIMDOptimizer) isLoopStart(idx int) bool {
	if idx+3 >= len(so.program.Instructions) {
		return false
	}

	// Pattern: LOAD (counter) → LOAD (limit) → CMP/LESS → JMP_IF_FALSE
	ins1 := so.program.Instructions[idx]
	ins2 := so.program.Instructions[idx+1]
	ins3 := so.program.Instructions[idx+2]
	ins4 := so.program.Instructions[idx+3]

	return ins1.Opcode == OpLoadVar &&
		ins2.Opcode == OpLoadConst &&
		(ins3.Opcode == OpLess || ins3.Opcode == OpEqual) &&
		ins4.Opcode == OpJumpIfFalse
}

// findLoopEnd locates the end of a loop
func (so *SIMDOptimizer) findLoopEnd(start int) int {
	// Look for the matching JUMP back to loop start
	for i := start + 1; i < len(so.program.Instructions); i++ {
		ins := so.program.Instructions[i]
		if ins.Opcode == OpJump {
			if len(ins.Args) > 0 {
				if target, ok := ins.Args[0].(int); ok {
					if target == start {
						return i
					}
				}
			}
		}
	}
	return -1
}

// canVectorize checks if a loop can be vectorized
func (so *SIMDOptimizer) canVectorize(start, end int) bool {
	// Check for data dependencies and uniform operations
	hasArrayOps := false
	hasComplexOps := false

	for i := start; i <= end; i++ {
		if i >= len(so.program.Instructions) {
			break
		}

		ins := so.program.Instructions[i]

		// Detect array operations
		if ins.Opcode == OpLoadVar || ins.Opcode == OpStoreVar {
			hasArrayOps = true
		}

		// Check for non-vectorizable operations
		switch ins.Opcode {
		case OpCall, OpJumpIfTrue, OpJumpIfFalse:
			if ins.Opcode == OpCall {
				hasComplexOps = true
			}
		}
	}

	// Vectorizable if has array ops and no complex ops
	return hasArrayOps && !hasComplexOps
}

// emitVectorInstructions converts scalar operations to vector operations
func (so *SIMDOptimizer) emitVectorInstructions() {
	newInstructions := []*Instruction{}
	i := 0

	for i < len(so.program.Instructions) {
		ins := so.program.Instructions[i]

		// Pattern: LOAD a[i] → LOAD b[i] → ADD → STORE c[i]
		// Convert to single VADD instruction
		if ins.Opcode == OpLoadVar && i+3 < len(so.program.Instructions) {
			next1 := so.program.Instructions[i+1]
			next2 := so.program.Instructions[i+2]
			next3 := so.program.Instructions[i+3]

			if next1.Opcode == OpLoadVar &&
				next2.Opcode == OpAdd &&
				next3.Opcode == OpStoreVar {

				// Emit vectorized instruction
				newInstructions = append(newInstructions,
					NewInstruction(OpLoadConst, "VADD_4F"))
				so.stats.VectorInstructions++
				so.stats.InstructionsFused += 4
				so.stats.VectorizedLoops++
				i += 4
				continue
			}

			// Pattern: LOAD a[i] → LOAD b[i] → MUL → STORE c[i]
			if next1.Opcode == OpLoadVar &&
				next2.Opcode == OpMul &&
				next3.Opcode == OpStoreVar {

				newInstructions = append(newInstructions,
					NewInstruction(OpLoadConst, "VMUL_4F"))
				so.stats.VectorInstructions++
				so.stats.InstructionsFused += 4
				i += 4
				continue
			}
		}

		newInstructions = append(newInstructions, ins)
		i++
	}

	so.program.Instructions = newInstructions
}

// analyzeMemoryBandwidth estimates memory bandwidth utilization
func (so *SIMDOptimizer) analyzeMemoryBandwidth() {
	if len(so.program.Instructions) == 0 {
		so.stats.MemoryBandwidth = 0
		return
	}

	// Count memory operations
	loadCount := 0
	storeCount := 0
	vectorOps := 0

	for _, ins := range so.program.Instructions {
		switch ins.Opcode {
		case OpLoadVar, OpLoadConst:
			loadCount++
		case OpStoreVar:
			storeCount++
		}

		// Count vector operations (marked with VADD, VMUL, etc.)
		if ins.Opcode == OpLoadConst {
			if s, ok := ins.Args[0].(string); ok {
				if len(s) > 1 && s[0] == 'V' {
					vectorOps++
				}
			}
		}
	}

	totalMemOps := loadCount + storeCount
	if totalMemOps > 0 {
		// Assume 4 elements per vector op = 4x bandwidth improvement
		bandwidth := float64(vectorOps*4) / float64(totalMemOps)
		if bandwidth > 1.0 {
			bandwidth = 1.0
		}
		so.stats.MemoryBandwidth = bandwidth
	}

	// Vectorization ratio
	if len(so.program.Instructions) > 0 {
		so.stats.VectorizationRatio = float64(so.stats.VectorInstructions) /
			float64(len(so.program.Instructions))
	}
}

// GetStatistics returns SIMD optimization statistics
func (so *SIMDOptimizer) GetStatistics() SIMDOptimizationStats {
	return so.stats
}

// Profile returns a profile of SIMD optimizations
func (so *SIMDOptimizer) Profile() string {
	return fmt.Sprintf(`SIMD Optimization Profile:
  Vectorizable Loops: %d
  Vectorized Loops: %d
  Vector Instructions: %d
  Instructions Fused: %d
  Memory Bandwidth: %.2f
  Vectorization Ratio: %.2f
`,
		so.stats.VectorizableLoops,
		so.stats.VectorizedLoops,
		so.stats.VectorInstructions,
		so.stats.InstructionsFused,
		so.stats.MemoryBandwidth,
		so.stats.VectorizationRatio)
}

// GetSupportedTypes returns supported vector types
func GetSupportedTypes() []VectorType {
	return []VectorType{Vec4F, Vec2F64, Vec8I, Vec4I64}
}

// TypeString returns string representation of vector type
func (vt VectorType) String() string {
	switch vt {
	case Vec4F:
		return "vec4f"
	case Vec2F64:
		return "vec2f64"
	case Vec8I:
		return "vec8i"
	case Vec4I64:
		return "vec4i64"
	default:
		return "unknown"
	}
}

// OpcodeString returns string representation of vector opcode
func (vo VectorOpcode) String() string {
	switch vo {
	case VecAdd:
		return "VADD"
	case VecSub:
		return "VSUB"
	case VecMul:
		return "VMUL"
	case VecDiv:
		return "VDIV"
	case VecLoad:
		return "VLOAD"
	case VecStore:
		return "VSTORE"
	case VecShuffle:
		return "VSHUFFLE"
	default:
		return "UNKNOWN"
	}
}

// CanVectorizeOperation checks if an opcode can be vectorized
func CanVectorizeOperation(op Opcode) bool {
	switch op {
	case OpAdd, OpSub, OpMul, OpDiv, OpMod:
		return true
	case OpLoadVar, OpStoreVar:
		return true
	default:
		return false
	}
}

// EstimateSpeedup calculates expected speedup from vectorization
func EstimateSpeedup(vectorLength int) float64 {
	// Assume overhead of 10% per vectorization
	overhead := 0.1
	return float64(vectorLength) / (1.0 + overhead)
}
