package runtime

import (
	"fmt"
	"strconv"

	"github.com/freelang-ai/gofree/internal/compiler"
)

// VM represents the FreeLang virtual machine
type VM struct {
	program      *compiler.IRProgram
	stack        []interface{}     // Value stack
	frames       []*Frame          // Call stack
	globals      map[string]interface{} // Global variables
	builtins     map[string]BuiltinFunc // Built-in functions
	memory       *Memory           // Memory manager
	retValue     interface{}       // Return value
	isRunning    bool
	errors       []RuntimeError
}

// Frame represents a function call frame
type Frame struct {
	fn          *compiler.IRFunction
	pc          int           // Program counter
	basePointer int           // Base pointer for local variables
	locals      map[string]interface{}
}

// BuiltinFunc represents a built-in function
type BuiltinFunc func(args ...interface{}) (interface{}, error)

// RuntimeError represents a runtime error
type RuntimeError struct {
	Line    int
	Column  int
	Message string
}

// Memory represents memory management
type Memory struct {
	objects map[int]interface{}
	nextID  int
}

// NewVM creates a new virtual machine
func NewVM(program *compiler.IRProgram) *VM {
	vm := &VM{
		program:  program,
		stack:    []interface{}{},
		frames:   []*Frame{},
		globals:  make(map[string]interface{}),
		builtins: make(map[string]BuiltinFunc),
		memory:   &Memory{objects: make(map[int]interface{}), nextID: 0},
		errors:   []RuntimeError{},
	}

	vm.registerBuiltins()
	return vm
}

// registerBuiltins registers built-in functions
func (vm *VM) registerBuiltins() {
	// I/O functions
	vm.builtins["print"] = func(args ...interface{}) (interface{}, error) {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(vm.toString(arg))
		}
		fmt.Println()
		return nil, nil
	}

	vm.builtins["println"] = func(args ...interface{}) (interface{}, error) {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(vm.toString(arg))
		}
		fmt.Println()
		return nil, nil
	}

	// Type functions
	vm.builtins["type"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("type() requires 1 argument")
		}
		return vm.getType(args[0]), nil
	}

	vm.builtins["len"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("len() requires 1 argument")
		}
		switch v := args[0].(type) {
		case string:
			return int64(len(v)), nil
		case []interface{}:
			return int64(len(v)), nil
		case map[string]interface{}:
			return int64(len(v)), nil
		default:
			return nil, fmt.Errorf("len() unsupported type")
		}
	}

	// Conversion functions
	vm.builtins["tostring"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return "", nil
		}
		return vm.toString(args[0]), nil
	}

	vm.builtins["tonumber"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return 0.0, nil
		}
		switch v := args[0].(type) {
		case float64:
			return v, nil
		case int64:
			return float64(v), nil
		case string:
			f, err := strconv.ParseFloat(v, 64)
			return f, err
		default:
			return 0.0, fmt.Errorf("cannot convert to number")
		}
	}

	// Array functions
	vm.builtins["push"] = func(args ...interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("push() requires 2 arguments")
		}
		if arr, ok := args[0].([]interface{}); ok {
			return append(arr, args[1]), nil
		}
		return nil, fmt.Errorf("push() first argument must be array")
	}

	vm.builtins["pop"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("pop() requires 1 argument")
		}
		if arr, ok := args[0].([]interface{}); ok && len(arr) > 0 {
			return arr[len(arr)-1], nil
		}
		return nil, fmt.Errorf("pop() cannot pop from empty array")
	}

	// Math functions
	vm.builtins["abs"] = func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return 0, nil
		}
		switch v := args[0].(type) {
		case float64:
			if v < 0 {
				return -v, nil
			}
			return v, nil
		case int64:
			if v < 0 {
				return -v, nil
			}
			return v, nil
		default:
			return 0, fmt.Errorf("abs() unsupported type")
		}
	}
}

// Run executes the IR program
func (vm *VM) Run() error {
	if vm.program == nil {
		return fmt.Errorf("program is nil")
	}

	vm.isRunning = true
	defer func() { vm.isRunning = false }()

	// Execute main instructions
	pc := 0
	for pc < len(vm.program.Instructions) && vm.isRunning {
		ins := vm.program.Instructions[pc]
		nextPC, err := vm.executeInstruction(ins, pc)
		if err != nil {
			return err
		}
		pc = nextPC
	}

	if len(vm.errors) > 0 {
		return fmt.Errorf("%d runtime errors", len(vm.errors))
	}

	return nil
}

// executeInstruction executes a single instruction
func (vm *VM) executeInstruction(ins *compiler.Instruction, pc int) (int, error) {
	switch ins.Opcode {
	case compiler.OpLoadConst:
		if len(ins.Args) > 0 {
			vm.push(ins.Args[0])
		}

	case compiler.OpLoadVar:
		if len(ins.Args) > 0 {
			if index, ok := ins.Args[0].(int); ok {
				if index >= 0 && index < len(vm.stack) {
					vm.push(vm.stack[index])
				}
			}
		}

	case compiler.OpStoreVar:
		if len(ins.Args) > 0 {
			if index, ok := ins.Args[0].(int); ok {
				val := vm.pop()
				if index < len(vm.stack) {
					vm.stack[index] = val
				}
			}
		}

	case compiler.OpAdd:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.add(left, right))

	case compiler.OpSub:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.subtract(left, right))

	case compiler.OpMul:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.multiply(left, right))

	case compiler.OpDiv:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.divide(left, right))

	case compiler.OpMod:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.modulo(left, right))

	case compiler.OpEqual:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.equal(left, right))

	case compiler.OpNotEqual:
		right := vm.pop()
		left := vm.pop()
		vm.push(!vm.equal(left, right).(bool))

	case compiler.OpLess:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.less(left, right))

	case compiler.OpGreater:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.greater(left, right))

	case compiler.OpAnd:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.isTruthy(left) && vm.isTruthy(right))

	case compiler.OpOr:
		right := vm.pop()
		left := vm.pop()
		vm.push(vm.isTruthy(left) || vm.isTruthy(right))

	case compiler.OpNot:
		val := vm.pop()
		vm.push(!vm.isTruthy(val))

	case compiler.OpJump:
		if len(ins.Args) > 0 {
			if target, ok := ins.Args[0].(int); ok {
				return target, nil
			}
		}

	case compiler.OpJumpIfFalse:
		val := vm.pop()
		if len(ins.Args) > 0 && !vm.isTruthy(val) {
			if target, ok := ins.Args[0].(int); ok {
				return target, nil
			}
		}

	case compiler.OpJumpIfTrue:
		val := vm.pop()
		if len(ins.Args) > 0 && vm.isTruthy(val) {
			if target, ok := ins.Args[0].(int); ok {
				return target, nil
			}
		}

	case compiler.OpCall:
		if len(ins.Args) >= 2 {
			fnName := ins.Args[0].(string)
			argCount := ins.Args[1].(int)

			// Get arguments from stack
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}

			// Call function
			if builtin, ok := vm.builtins[fnName]; ok {
				result, err := builtin(args...)
				if err != nil {
					return pc + 1, err
				}
				vm.push(result)
			}
		}

	case compiler.OpArrayCreate:
		if len(ins.Args) > 0 {
			count := ins.Args[0].(int)
			arr := make([]interface{}, count)
			for i := count - 1; i >= 0; i-- {
				arr[i] = vm.pop()
			}
			vm.push(arr)
		}

	case compiler.OpPop:
		vm.pop()

	case compiler.OpReturn:
		vm.retValue = vm.pop()
		vm.isRunning = false

	default:
		return pc + 1, fmt.Errorf("unknown opcode: %v", ins.Opcode)
	}

	return pc + 1, nil
}

// Helper methods

func (vm *VM) push(val interface{}) {
	vm.stack = append(vm.stack, val)
}

func (vm *VM) pop() interface{} {
	if len(vm.stack) == 0 {
		return nil
	}
	val := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return val
}

func (vm *VM) peek() interface{} {
	if len(vm.stack) == 0 {
		return nil
	}
	return vm.stack[len(vm.stack)-1]
}

func (vm *VM) add(left, right interface{}) interface{} {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l + r
		}
	case string:
		if r, ok := right.(string); ok {
			return l + r
		}
	}
	return 0
}

func (vm *VM) subtract(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok {
			return l - r
		}
	}
	return 0
}

func (vm *VM) multiply(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok {
			return l * r
		}
	}
	return 0
}

func (vm *VM) divide(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok && r != 0 {
			return l / r
		}
	}
	return 0
}

func (vm *VM) modulo(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok && r != 0 {
			return float64(int64(l) % int64(r))
		}
	}
	return 0
}

func (vm *VM) equal(left, right interface{}) interface{} {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l == r
		}
	case string:
		if r, ok := right.(string); ok {
			return l == r
		}
	case bool:
		if r, ok := right.(bool); ok {
			return l == r
		}
	case nil:
		return right == nil
	}
	return false
}

func (vm *VM) less(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok {
			return l < r
		}
	}
	return false
}

func (vm *VM) greater(left, right interface{}) interface{} {
	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok {
			return l > r
		}
	}
	return false
}

func (vm *VM) isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case nil:
		return false
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return true
	}
}

func (vm *VM) toString(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%f", v)
	case string:
		return v
	case []interface{}:
		return fmt.Sprintf("[array with %d elements]", len(v))
	default:
		return "[object]"
	}
}

func (vm *VM) getType(val interface{}) string {
	switch val.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

// ReturnValue returns the last return value
func (vm *VM) ReturnValue() interface{} {
	return vm.retValue
}

// Errors returns runtime errors
func (vm *VM) Errors() []RuntimeError {
	return vm.errors
}
