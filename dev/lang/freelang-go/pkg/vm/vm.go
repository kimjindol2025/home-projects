package vm

import (
	"fmt"

	"github.com/kimjindol2025/freelang-go/pkg/compiler"
	"github.com/kimjindol2025/freelang-go/pkg/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

// Frame represents a function call frame
type Frame struct {
	Bytecode []byte
	IP       int // Instruction pointer
}

// VM is a stack-based virtual machine for executing bytecode
type VM struct {
	constants []interface{}
	stack     []object.Object
	sp        int // Stack pointer (points to next free slot)
	globals   []object.Object
	frames    []Frame
	fp        int // Frame pointer
}

// New creates a new virtual machine
func New(bytecode []byte, constants []interface{}) *VM {
	vm := &VM{
		constants: constants,
		stack:     make([]object.Object, StackSize),
		sp:        0,
		globals:   make([]object.Object, GlobalsSize),
		frames:    make([]Frame, MaxFrames),
		fp:        0,
	}

	// Initialize first frame
	vm.frames[0] = Frame{Bytecode: bytecode, IP: 0}
	vm.fp = 1

	return vm
}

// Run executes the bytecode
func (vm *VM) Run() error {
	frame := &vm.frames[0]

	for frame.IP < len(frame.Bytecode) {
		opcode := compiler.OpCode(frame.Bytecode[frame.IP])
		frame.IP++

		switch opcode {
			case compiler.OpConstant:
				idx := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				vm.push(vm.objectFromConstant(vm.constants[idx]))

			case compiler.OpTrue:
				vm.push(&object.Boolean{Value: true})

			case compiler.OpFalse:
				vm.push(&object.Boolean{Value: false})

			case compiler.OpNull:
				vm.push(&object.Null{})

			case compiler.OpPop:
				vm.sp-- // Don't actually pop, just decrement pointer for next access

			case compiler.OpAdd:
				right := vm.pop()
				left := vm.pop()
				result := vm.add(left, right)
				vm.push(result)

			case compiler.OpSubtract:
				right := vm.pop()
				left := vm.pop()
				result := vm.subtract(left, right)
				vm.push(result)

			case compiler.OpMultiply:
				right := vm.pop()
				left := vm.pop()
				result := vm.multiply(left, right)
				vm.push(result)

			case compiler.OpDivide:
				right := vm.pop()
				left := vm.pop()
				result := vm.divide(left, right)
				vm.push(result)

			case compiler.OpModulo:
				right := vm.pop()
				left := vm.pop()
				result := vm.modulo(left, right)
				vm.push(result)

			case compiler.OpPower:
				right := vm.pop()
				left := vm.pop()
				result := vm.power(left, right)
				vm.push(result)

			case compiler.OpEqual:
				right := vm.pop()
				left := vm.pop()
				result := vm.equal(left, right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpNotEqual:
				right := vm.pop()
				left := vm.pop()
				result := !vm.equal(left, right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpLessThan:
				right := vm.pop()
				left := vm.pop()
				result := vm.lessThan(left, right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpLessThanOrEqual:
				right := vm.pop()
				left := vm.pop()
				result := vm.lessThan(left, right) || vm.equal(left, right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpGreaterThan:
				right := vm.pop()
				left := vm.pop()
				result := vm.lessThan(right, left)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpGreaterThanOrEqual:
				right := vm.pop()
				left := vm.pop()
				result := vm.lessThan(right, left) || vm.equal(left, right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpAnd:
				right := vm.pop()
				left := vm.pop()
				result := vm.isTruthy(left) && vm.isTruthy(right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpOr:
				right := vm.pop()
				left := vm.pop()
				result := vm.isTruthy(left) || vm.isTruthy(right)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpNot:
				operand := vm.pop()
				result := !vm.isTruthy(operand)
				vm.push(&object.Boolean{Value: result})

			case compiler.OpNegate:
				operand := vm.pop()
				if intObj, ok := operand.(*object.Integer); ok {
					vm.push(&object.Integer{Value: -intObj.Value})
				} else {
					vm.push(&object.Null{})
				}

			case compiler.OpBitAnd:
				right := vm.pop()
				left := vm.pop()
				result := vm.bitAnd(left, right)
				vm.push(result)

			case compiler.OpBitOr:
				right := vm.pop()
				left := vm.pop()
				result := vm.bitOr(left, right)
				vm.push(result)

			case compiler.OpBitXor:
				right := vm.pop()
				left := vm.pop()
				result := vm.bitXor(left, right)
				vm.push(result)

			case compiler.OpBitNot:
				operand := vm.pop()
				result := vm.bitNot(operand)
				vm.push(result)

			case compiler.OpLeftShift:
				right := vm.pop()
				left := vm.pop()
				result := vm.leftShift(left, right)
				vm.push(result)

			case compiler.OpRightShift:
				right := vm.pop()
				left := vm.pop()
				result := vm.rightShift(left, right)
				vm.push(result)

			case compiler.OpArray:
				numElements := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				elements := make([]object.Object, numElements)
				for i := int(numElements) - 1; i >= 0; i-- {
					elements[i] = vm.pop()
				}
				vm.push(&object.Array{Elements: elements})

			case compiler.OpHash:
				numPairs := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				pairs := make(map[string]object.Object)
				// Pop in reverse order (stack is LIFO)
				for i := 0; i < int(numPairs); i++ {
					value := vm.pop()
					key := vm.pop()
					keyStr := vm.objectToString(key)
					pairs[keyStr] = value
				}
				vm.push(&object.Hash{Pairs: pairs})

			case compiler.OpIndex:
				index := vm.pop()
				left := vm.pop()
				result := vm.index(left, index)
				vm.push(result)

			case compiler.OpJump:
				pos := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP = int(pos)

			case compiler.OpJumpNotTruthy:
				pos := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				condition := vm.pop()
				if !vm.isTruthy(condition) {
					frame.IP = int(pos)
				}

			case compiler.OpGetLocal:
				localIdx := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				// Get local variable from stack
				if vm.sp > int(localIdx) {
					vm.push(vm.stack[vm.sp-int(localIdx)-1])
				}

			case compiler.OpSetLocal:
				_ = vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				// Set local variable on stack
				value := vm.pop()
				if vm.sp >= 0 && vm.sp < len(vm.stack) {
					vm.stack[vm.sp] = value
				}

			case compiler.OpGetGlobal:
				idx := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				vm.push(vm.globals[idx])

			case compiler.OpSetGlobal:
				idx := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				vm.globals[idx] = vm.pop()

			case compiler.OpReturn:
				vm.sp--

			case compiler.OpReturnValue:
				vm.fp++
				vm.frames[vm.fp] = Frame{Bytecode: frame.Bytecode, IP: len(frame.Bytecode)}

			case compiler.OpCall:
				numArgs := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				_ = numArgs // TODO: function call implementation

			case compiler.OpFunction:
				idx := vm.readUint32(frame.Bytecode, frame.IP)
				frame.IP += 4
				vm.push(&object.Function{Name: fmt.Sprintf("fn_%d", idx)})

			case compiler.OpPrint:
				obj := vm.pop()
				fmt.Println(vm.objectToString(obj))

			case compiler.OpLen:
				obj := vm.pop()
				length := vm.len(obj)
				vm.push(&object.Integer{Value: int64(length)})

			case compiler.OpType:
				obj := vm.pop()
				typeStr := vm.getType(obj)
				vm.push(&object.String{Value: typeStr})

			default:
				return fmt.Errorf("unknown opcode: %v", opcode)
			}
		}

	return nil
}

// LastPoppedStackElem returns the last popped element
func (vm *VM) LastPoppedStackElem() object.Object {
	if vm.sp >= 0 && vm.sp < len(vm.stack) {
		return vm.stack[vm.sp]
	}
	if vm.sp-1 >= 0 && vm.sp-1 < len(vm.stack) {
		return vm.stack[vm.sp-1]
	}
	return nil
}

// Helper methods

func (vm *VM) push(obj object.Object) {
	vm.stack[vm.sp] = obj
	vm.sp++
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj
}

func (vm *VM) readUint32(bytecode []byte, offset int) uint32 {
	return uint32(bytecode[offset])<<24 | uint32(bytecode[offset+1])<<16 | uint32(bytecode[offset+2])<<8 | uint32(bytecode[offset+3])
}

func (vm *VM) objectFromConstant(constant interface{}) object.Object {
	switch v := constant.(type) {
	case int64:
		return &object.Integer{Value: v}
	case string:
		return &object.String{Value: v}
	default:
		return &object.Null{}
	}
}

func (vm *VM) objectToString(obj object.Object) string {
	switch o := obj.(type) {
	case *object.Integer:
		return fmt.Sprintf("%d", o.Value)
	case *object.String:
		return o.Value
	case *object.Boolean:
		return fmt.Sprintf("%v", o.Value)
	default:
		return ""
	}
}

func (vm *VM) isTruthy(obj object.Object) bool {
	switch o := obj.(type) {
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return false
	case *object.Integer:
		return o.Value != 0
	default:
		return true
	}
}

func (vm *VM) add(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value + r.Value}
	}

	ls, ok1 := left.(*object.String)
	rs, ok2 := right.(*object.String)
	if ok1 && ok2 {
		return &object.String{Value: ls.Value + rs.Value}
	}

	return &object.Null{}
}

func (vm *VM) subtract(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value - r.Value}
	}
	return &object.Null{}
}

func (vm *VM) multiply(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value * r.Value}
	}
	return &object.Null{}
}

func (vm *VM) divide(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 && r.Value != 0 {
		return &object.Integer{Value: l.Value / r.Value}
	}
	return &object.Null{}
}

func (vm *VM) modulo(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 && r.Value != 0 {
		return &object.Integer{Value: l.Value % r.Value}
	}
	return &object.Null{}
}

func (vm *VM) power(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		result := int64(1)
		for i := int64(0); i < r.Value; i++ {
			result *= l.Value
		}
		return &object.Integer{Value: result}
	}
	return &object.Null{}
}

func (vm *VM) equal(left, right object.Object) bool {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return l.Value == r.Value
	}

	lb, ok1 := left.(*object.Boolean)
	rb, ok2 := right.(*object.Boolean)
	if ok1 && ok2 {
		return lb.Value == rb.Value
	}

	ls, ok1 := left.(*object.String)
	rs, ok2 := right.(*object.String)
	if ok1 && ok2 {
		return ls.Value == rs.Value
	}

	_, ok1 = left.(*object.Null)
	_, ok2 = right.(*object.Null)
	if ok1 && ok2 {
		return true
	}

	return false
}

func (vm *VM) lessThan(left, right object.Object) bool {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return l.Value < r.Value
	}
	return false
}

func (vm *VM) bitAnd(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value & r.Value}
	}
	return &object.Null{}
}

func (vm *VM) bitOr(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value | r.Value}
	}
	return &object.Null{}
}

func (vm *VM) bitXor(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value ^ r.Value}
	}
	return &object.Null{}
}

func (vm *VM) bitNot(operand object.Object) object.Object {
	o, ok := operand.(*object.Integer)
	if ok {
		return &object.Integer{Value: ^o.Value}
	}
	return &object.Null{}
}

func (vm *VM) leftShift(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 && r.Value >= 0 {
		return &object.Integer{Value: l.Value << uint(r.Value)}
	}
	return &object.Null{}
}

func (vm *VM) rightShift(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 && r.Value >= 0 {
		return &object.Integer{Value: l.Value >> uint(r.Value)}
	}
	return &object.Null{}
}

func (vm *VM) index(left, index object.Object) object.Object {
	switch l := left.(type) {
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if ok && idx.Value >= 0 && int(idx.Value) < len(l.Elements) {
			return l.Elements[int(idx.Value)]
		}
		return &object.Null{}

	case *object.Hash:
		key := vm.objectToString(index)
		if val, ok := l.Pairs[key]; ok {
			return val
		}
		return &object.Null{}

	default:
		return &object.Null{}
	}
}

func (vm *VM) len(obj object.Object) int {
	switch o := obj.(type) {
	case *object.String:
		return len(o.Value)
	case *object.Array:
		return len(o.Elements)
	case *object.Hash:
		return len(o.Pairs)
	default:
		return 0
	}
}

func (vm *VM) getType(obj object.Object) string {
	switch obj.(type) {
	case *object.Integer:
		return "int"
	case *object.String:
		return "string"
	case *object.Boolean:
		return "bool"
	case *object.Array:
		return "array"
	case *object.Hash:
		return "hash"
	case *object.Null:
		return "null"
	default:
		return "unknown"
	}
}
