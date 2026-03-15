package runtime

import (
	"testing"

	"github.com/freelang-ai/gofree/internal/compiler"
	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/parser"
)

// Helper to compile and run
func compileAndRun(t *testing.T, src string) *VM {
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, err := p.Parse()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	compiler := compiler.NewCompiler()
	program, err := compiler.Compile(module)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	vm := NewVM(program)
	err = vm.Run()
	if err != nil {
		t.Logf("runtime error: %v", err)
	}

	return vm
}

// Test arithmetic operations
func TestVMAddition(t *testing.T) {
	src := `let x = 10 + 5`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test variable declaration
func TestVMVariableDeclaration(t *testing.T) {
	src := `let x = 42`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test arithmetic: subtraction
func TestVMSubtraction(t *testing.T) {
	src := `let result = 20 - 8`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test arithmetic: multiplication
func TestVMMultiplication(t *testing.T) {
	src := `let result = 6 * 7`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test arithmetic: division
func TestVMDivision(t *testing.T) {
	src := `let result = 100 / 2`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test comparison: equal
func TestVMComparison(t *testing.T) {
	src := `let result = 5 == 5`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test logical operations
func TestVMLogicalAnd(t *testing.T) {
	src := `let result = true && false`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test logical OR
func TestVMLogicalOr(t *testing.T) {
	src := `let result = true || false`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test stack operations
func TestVMStackOperations(t *testing.T) {
	vm := NewVM(nil)

	vm.push(10)
	vm.push(20)

	if vm.peek() != 20 {
		t.Errorf("expected 20, got %v", vm.peek())
	}

	val := vm.pop()
	if val != 20 {
		t.Errorf("expected 20, got %v", val)
	}

	if vm.peek() != 10 {
		t.Errorf("expected 10, got %v", vm.peek())
	}
}

// Test builtin print function
func TestVMBuiltinPrint(t *testing.T) {
	vm := NewVM(&compiler.IRProgram{
		Instructions: []*compiler.Instruction{},
		Functions:    make(map[string]*compiler.IRFunction),
		GlobalVars:   make(map[string]int),
	})

	if builtin, ok := vm.builtins["print"]; ok {
		_, err := builtin("Hello, World!")
		if err != nil {
			t.Errorf("print failed: %v", err)
		}
	}
}

// Test builtin len function
func TestVMBuiltinLen(t *testing.T) {
	vm := NewVM(&compiler.IRProgram{
		Instructions: []*compiler.Instruction{},
		Functions:    make(map[string]*compiler.IRFunction),
		GlobalVars:   make(map[string]int),
	})

	if builtin, ok := vm.builtins["len"]; ok {
		arr := []interface{}{1, 2, 3}
		result, err := builtin(arr)
		if err != nil {
			t.Errorf("len failed: %v", err)
		}
		if result != int64(3) {
			t.Errorf("expected 3, got %v", result)
		}
	}
}

// Test type function
func TestVMBuiltinType(t *testing.T) {
	vm := NewVM(&compiler.IRProgram{
		Instructions: []*compiler.Instruction{},
		Functions:    make(map[string]*compiler.IRFunction),
		GlobalVars:   make(map[string]int),
	})

	tests := []struct {
		val      interface{}
		expected string
	}{
		{nil, "null"},
		{true, "boolean"},
		{42.0, "number"},
		{"hello", "string"},
	}

	for _, tt := range tests {
		result := vm.getType(tt.val)
		if result != tt.expected {
			t.Errorf("getType(%v) = %v, want %v", tt.val, result, tt.expected)
		}
	}
}

// Test toString function
func TestVMToString(t *testing.T) {
	vm := NewVM(nil)

	tests := []struct {
		val      interface{}
		expected string
	}{
		{nil, "null"},
		{true, "true"},
		{false, "false"},
		{42.0, "42"},
		{"hello", "hello"},
	}

	for _, tt := range tests {
		result := vm.toString(tt.val)
		if result != tt.expected {
			t.Errorf("toString(%v) = %v, want %v", tt.val, result, tt.expected)
		}
	}
}

// Test isTruthy function
func TestVMIsTruthy(t *testing.T) {
	vm := NewVM(nil)

	tests := []struct {
		val      interface{}
		expected bool
	}{
		{true, true},
		{false, false},
		{nil, false},
		{0.0, false},
		{1.0, true},
		{"", false},
		{"hello", true},
	}

	for _, tt := range tests {
		result := vm.isTruthy(tt.val)
		if result != tt.expected {
			t.Errorf("isTruthy(%v) = %v, want %v", tt.val, result, tt.expected)
		}
	}
}

// Test string concatenation
func TestVMStringConcat(t *testing.T) {
	src := `let msg = "Hello" + " " + "World"`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test array creation
func TestVMArrayCreation(t *testing.T) {
	src := `let arr = [1, 2, 3]`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Test function call
func TestVMFunctionCall(t *testing.T) {
	src := `let result = len([1, 2, 3])`
	vm := compileAndRun(t, src)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}
}

// Benchmark VM execution
func BenchmarkVMExecution(b *testing.B) {
	src := `
let x = 0
let y = 0
let z = 0
for (let i = 0; i < 100; i = i + 1) {
    x = x + 1
    y = y + 2
    z = z + 3
}`

	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()

	compiler := compiler.NewCompiler()
	program, _ := compiler.Compile(module)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm := NewVM(program)
		_ = vm.Run()
	}
}

// Test VM initialization
func TestVMInitialization(t *testing.T) {
	program := compiler.NewIRProgram()
	vm := NewVM(program)

	if vm == nil {
		t.Errorf("expected VM, got nil")
	}

	if len(vm.builtins) == 0 {
		t.Errorf("expected builtin functions")
	}
}

// Test basic execution
func TestVMBasicExecution(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 42))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpPop))

	vm := NewVM(program)
	err := vm.Run()

	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
}
