package compiler

// Opcode represents IR instruction types
type Opcode int

const (
	// Constants and Variables
	OpLoadConst Opcode = iota
	OpLoadVar
	OpStoreVar
	OpLoadGlobal
	OpStoreGlobal

	// Arithmetic Operations
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNeg

	// Comparison Operations
	OpEqual
	OpNotEqual
	OpLess
	OpLessOrEqual
	OpGreater
	OpGreaterOrEqual

	// Logical Operations
	OpAnd
	OpOr
	OpNot

	// Control Flow
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpReturn
	OpYield

	// Function Calls
	OpCall
	OpCallBuiltin

	// Array Operations
	OpArrayCreate
	OpArrayAccess
	OpArraySet
	OpArrayLen

	// Object Operations
	OpObjectCreate
	OpObjectAccess
	OpObjectSet

	// Type Operations
	OpTypeOf
	OpInstanceOf

	// Exception Handling
	OpThrow
	OpTryEnter
	OpTryExit

	// Other
	OpPop
	OpDuplicate
	OpNoOp

	// Async
	OpAwait
)

var opcodeNames = map[Opcode]string{
	OpLoadConst:      "LOAD_CONST",
	OpLoadVar:        "LOAD_VAR",
	OpStoreVar:       "STORE_VAR",
	OpLoadGlobal:     "LOAD_GLOBAL",
	OpStoreGlobal:    "STORE_GLOBAL",
	OpAdd:            "ADD",
	OpSub:            "SUB",
	OpMul:            "MUL",
	OpDiv:            "DIV",
	OpMod:            "MOD",
	OpNeg:            "NEG",
	OpEqual:          "EQUAL",
	OpNotEqual:       "NOT_EQUAL",
	OpLess:           "LESS",
	OpLessOrEqual:    "LESS_OR_EQUAL",
	OpGreater:        "GREATER",
	OpGreaterOrEqual: "GREATER_OR_EQUAL",
	OpAnd:            "AND",
	OpOr:             "OR",
	OpNot:            "NOT",
	OpJump:           "JUMP",
	OpJumpIfFalse:    "JUMP_IF_FALSE",
	OpJumpIfTrue:     "JUMP_IF_TRUE",
	OpReturn:         "RETURN",
	OpYield:          "YIELD",
	OpCall:           "CALL",
	OpCallBuiltin:    "CALL_BUILTIN",
	OpArrayCreate:    "ARRAY_CREATE",
	OpArrayAccess:    "ARRAY_ACCESS",
	OpArraySet:       "ARRAY_SET",
	OpArrayLen:       "ARRAY_LEN",
	OpObjectCreate:   "OBJECT_CREATE",
	OpObjectAccess:   "OBJECT_ACCESS",
	OpObjectSet:      "OBJECT_SET",
	OpTypeOf:         "TYPE_OF",
	OpInstanceOf:     "INSTANCE_OF",
	OpThrow:          "THROW",
	OpTryEnter:       "TRY_ENTER",
	OpTryExit:        "TRY_EXIT",
	OpPop:            "POP",
	OpDuplicate:      "DUPLICATE",
	OpNoOp:           "NOOP",
	OpAwait:          "AWAIT",
}

// String returns the string representation of an opcode
func (op Opcode) String() string {
	if name, ok := opcodeNames[op]; ok {
		return name
	}
	return "UNKNOWN"
}

// Instruction represents a single IR instruction
type Instruction struct {
	Opcode Opcode
	Args   []interface{} // Variable arguments (constants, variable indices, jump targets, etc.)
	Line   int           // Source line number for error reporting
	Column int           // Source column number
}

// NewInstruction creates a new instruction
func NewInstruction(opcode Opcode, args ...interface{}) *Instruction {
	return &Instruction{
		Opcode: opcode,
		Args:   args,
		Line:   0,
		Column: 0,
	}
}

// String returns the string representation of an instruction
func (ins *Instruction) String() string {
	result := ins.Opcode.String()
	if len(ins.Args) > 0 {
		result += " "
		for i, arg := range ins.Args {
			if i > 0 {
				result += ", "
			}
			result += formatArg(arg)
		}
	}
	return result
}

// formatArg formats an argument for display
func formatArg(arg interface{}) string {
	switch v := arg.(type) {
	case string:
		return "\"" + v + "\""
	case int:
		return string(rune(v + 48)) // Simple number formatting
	case float64:
		return "f" // Placeholder
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []interface{}:
		return "[...]"
	default:
		return "?"
	}
}

// IRProgram represents a compiled IR program
type IRProgram struct {
	Instructions []*Instruction
	Constants    []interface{}        // Constant values (strings, numbers)
	Functions    map[string]*IRFunction
	GlobalVars   map[string]int       // Global variable names to index
	SourceMap    map[int]SourceLocation // Instruction index to source location
}

// SourceLocation represents a location in source code
type SourceLocation struct {
	Line   int
	Column int
	File   string
}

// IRFunction represents a compiled function
type IRFunction struct {
	Name         string
	ParamCount   int
	LocalVarCount int
	Instructions []*Instruction
	SourceMap    map[int]SourceLocation
}

// NewIRProgram creates a new IR program
func NewIRProgram() *IRProgram {
	return &IRProgram{
		Instructions: []*Instruction{},
		Constants:    []interface{}{},
		Functions:    make(map[string]*IRFunction),
		GlobalVars:   make(map[string]int),
		SourceMap:    make(map[int]SourceLocation),
	}
}

// NewIRFunction creates a new IR function
func NewIRFunction(name string) *IRFunction {
	return &IRFunction{
		Name:         name,
		ParamCount:   0,
		LocalVarCount: 0,
		Instructions: []*Instruction{},
		SourceMap:    make(map[int]SourceLocation),
	}
}
