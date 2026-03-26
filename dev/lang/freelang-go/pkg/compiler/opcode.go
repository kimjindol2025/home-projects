package compiler

import (
	"fmt"
)

// OpCode represents a bytecode instruction
type OpCode byte

const (
	// Constants and Literals
	OpConstant OpCode = iota
	OpTrue
	OpFalse
	OpNull

	// Variables
	OpGetLocal
	OpSetLocal
	OpGetGlobal
	OpSetGlobal

	// Arrays and Hashes
	OpArray
	OpHash
	OpIndex
	OpIndexAssign

	// Arithmetic operators
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpPower

	// Comparison operators
	OpEqual
	OpNotEqual
	OpLessThan
	OpLessThanOrEqual
	OpGreaterThan
	OpGreaterThanOrEqual

	// Logical operators
	OpAnd
	OpOr
	OpNot
	OpNegate

	// Bitwise operators
	OpBitAnd
	OpBitOr
	OpBitXor
	OpBitNot
	OpLeftShift
	OpRightShift

	// Control flow
	OpJump
	OpJumpNotTruthy
	OpLoop
	OpBreak
	OpContinue

	// Functions
	OpFunction
	OpCall
	OpReturn
	OpReturnValue

	// Stack manipulation
	OpPop
	OpDup

	// Built-in functions
	OpLen
	OpPrint
	OpType

	// Scope
	OpEnterScope
	OpExitScope
)

var opCodeNames = map[OpCode]string{
	OpConstant:           "CONSTANT",
	OpTrue:               "TRUE",
	OpFalse:              "FALSE",
	OpNull:               "NULL",
	OpGetLocal:           "GET_LOCAL",
	OpSetLocal:           "SET_LOCAL",
	OpGetGlobal:          "GET_GLOBAL",
	OpSetGlobal:          "SET_GLOBAL",
	OpArray:              "ARRAY",
	OpHash:               "HASH",
	OpIndex:              "INDEX",
	OpIndexAssign:        "INDEX_ASSIGN",
	OpAdd:                "ADD",
	OpSubtract:           "SUBTRACT",
	OpMultiply:           "MULTIPLY",
	OpDivide:             "DIVIDE",
	OpModulo:             "MODULO",
	OpPower:              "POWER",
	OpEqual:              "EQUAL",
	OpNotEqual:           "NOT_EQUAL",
	OpLessThan:           "LESS_THAN",
	OpLessThanOrEqual:    "LESS_THAN_OR_EQUAL",
	OpGreaterThan:        "GREATER_THAN",
	OpGreaterThanOrEqual: "GREATER_THAN_OR_EQUAL",
	OpAnd:                "AND",
	OpOr:                 "OR",
	OpNot:                "NOT",
	OpNegate:             "NEGATE",
	OpBitAnd:             "BIT_AND",
	OpBitOr:              "BIT_OR",
	OpBitXor:             "BIT_XOR",
	OpBitNot:             "BIT_NOT",
	OpLeftShift:          "LEFT_SHIFT",
	OpRightShift:         "RIGHT_SHIFT",
	OpJump:               "JUMP",
	OpJumpNotTruthy:      "JUMP_NOT_TRUTHY",
	OpLoop:               "LOOP",
	OpBreak:              "BREAK",
	OpContinue:           "CONTINUE",
	OpFunction:           "FUNCTION",
	OpCall:               "CALL",
	OpReturn:             "RETURN",
	OpReturnValue:        "RETURN_VALUE",
	OpPop:                "POP",
	OpDup:                "DUP",
	OpLen:                "LEN",
	OpPrint:              "PRINT",
	OpType:               "TYPE",
	OpEnterScope:         "ENTER_SCOPE",
	OpExitScope:          "EXIT_SCOPE",
}

// String returns the string representation of the opcode
func (op OpCode) String() string {
	if name, ok := opCodeNames[op]; ok {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", op)
}

// Instruction represents a bytecode instruction
type Instruction struct {
	OpCode   OpCode
	Operands []int
}

// Instructions represents a sequence of instructions
type Instructions []byte

// String returns a string representation of instructions
func (insts Instructions) String() string {
	var result string
	i := 0
	for i < len(insts) {
		opcode := OpCode(insts[i])
		result += fmt.Sprintf("%d %s", i, opCodeNames[opcode])

		// Parse operands based on opcode
		operands, width := DecodeInstructions(insts[i:])
		for _, op := range operands {
			result += fmt.Sprintf(" %d", op)
		}
		result += "\n"

		i += width
	}
	return result
}

// EncodeInstructions encodes an instruction with operands
func EncodeInstructions(opcode OpCode, operands ...int) []byte {
	var result []byte
	result = append(result, byte(opcode))

	for _, operand := range operands {
		// Encode as 4-byte big-endian integer
		result = append(result, byte((operand>>24)&0xFF))
		result = append(result, byte((operand>>16)&0xFF))
		result = append(result, byte((operand>>8)&0xFF))
		result = append(result, byte(operand&0xFF))
	}

	return result
}

// DecodeInstructions decodes an instruction and returns operands and width
func DecodeInstructions(bytecode []byte) ([]int, int) {
	opcode := OpCode(bytecode[0])

	// Determine number of operands based on opcode
	operandCount := getOperandCount(opcode)

	var operands []int
	offset := 1

	for i := 0; i < operandCount; i++ {
		operand := (int(bytecode[offset]) << 24) |
			(int(bytecode[offset+1]) << 16) |
			(int(bytecode[offset+2]) << 8) |
			int(bytecode[offset+3])
		operands = append(operands, operand)
		offset += 4
	}

	return operands, offset
}

// getOperandCount returns the number of operands for an opcode
func getOperandCount(opcode OpCode) int {
	switch opcode {
	case OpConstant, OpGetLocal, OpSetLocal, OpGetGlobal, OpSetGlobal,
		OpArray, OpHash, OpIndex, OpIndexAssign, OpJump, OpJumpNotTruthy,
		OpLoop, OpFunction, OpCall:
		return 1

	case OpTrue, OpFalse, OpNull, OpAdd, OpSubtract, OpMultiply, OpDivide,
		OpModulo, OpPower, OpEqual, OpNotEqual, OpLessThan, OpLessThanOrEqual,
		OpGreaterThan, OpGreaterThanOrEqual, OpAnd, OpOr, OpNot, OpNegate, OpBitAnd,
		OpBitOr, OpBitXor, OpBitNot, OpLeftShift, OpRightShift, OpBreak,
		OpContinue, OpReturn, OpReturnValue, OpPop, OpDup, OpLen, OpPrint,
		OpType, OpEnterScope, OpExitScope:
		return 0

	default:
		return 0
	}
}
