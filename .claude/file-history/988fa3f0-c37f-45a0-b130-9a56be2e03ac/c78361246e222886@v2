package compiler

import (
	"testing"

	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/parser"
)

// Helper to parse and compile
func parseAndCompile(t *testing.T, src string) *IRProgram {
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, err := p.Parse()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	compiler := NewCompiler()
	program, err := compiler.Compile(module)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	return program
}

// Test simple variable declaration
func TestCompileVariableDeclaration(t *testing.T) {
	src := `let x = 42`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check that instructions were generated
	if len(program.Instructions) == 0 {
		t.Errorf("expected instructions, got none")
	}

	// Check opcodes
	found := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpLoadConst || ins.Opcode == OpStoreVar {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected LOAD_CONST or STORE_VAR instruction")
	}
}

// Test function compilation
func TestCompileFunctionStatement(t *testing.T) {
	src := `fn add(x, y) { return x + y }`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check function was registered
	if _, ok := program.Functions["add"]; !ok {
		t.Errorf("expected function 'add' in program")
	}
}

// Test arithmetic operations
func TestCompileArithmetic(t *testing.T) {
	src := `let result = 10 + 5`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for arithmetic opcode
	found := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpAdd {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected ADD instruction")
	}
}

// Test if statement compilation
func TestCompileIfStatement(t *testing.T) {
	src := `
let x = 10
if (x > 5) {
    let y = 20
}`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for jump instructions
	jumpFound := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpJumpIfFalse || ins.Opcode == OpJump {
			jumpFound = true
			break
		}
	}

	if !jumpFound {
		t.Errorf("expected jump instruction")
	}
}

// Test while loop compilation
func TestCompileWhileLoop(t *testing.T) {
	src := `
let x = 0
while (x < 10) {
    let y = x + 1
}`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for jump instructions
	jumpFound := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpJump {
			jumpFound = true
			break
		}
	}

	if !jumpFound {
		t.Errorf("expected jump instruction in while loop")
	}
}

// Test function call compilation
func TestCompileCallExpression(t *testing.T) {
	src := `let result = add(1, 2)`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for call instruction
	callFound := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpCall {
			callFound = true
			break
		}
	}

	if !callFound {
		t.Errorf("expected CALL instruction")
	}
}

// Test array compilation
func TestCompileArrayExpression(t *testing.T) {
	src := `let arr = [1, 2, 3]`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for array create instruction
	arrayFound := false
	for _, ins := range program.Instructions {
		if ins.Opcode == OpArrayCreate {
			arrayFound = true
			break
		}
	}

	if !arrayFound {
		t.Errorf("expected ARRAY_CREATE instruction")
	}
}

// Test opcode string representation
func TestOpcodeString(t *testing.T) {
	tests := []struct {
		opcode   Opcode
		expected string
	}{
		{OpLoadConst, "LOAD_CONST"},
		{OpAdd, "ADD"},
		{OpSub, "SUB"},
		{OpReturn, "RETURN"},
		{OpCall, "CALL"},
	}

	for _, tt := range tests {
		if got := tt.opcode.String(); got != tt.expected {
			t.Errorf("Opcode.String() = %v, want %v", got, tt.expected)
		}
	}
}

// Test instruction string representation
func TestInstructionString(t *testing.T) {
	ins := NewInstruction(OpLoadConst, 42)
	if ins.String() == "" {
		t.Errorf("expected non-empty instruction string")
	}
}

// Test IR program creation
func TestIRProgramCreation(t *testing.T) {
	program := NewIRProgram()

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	if len(program.Instructions) != 0 {
		t.Errorf("expected empty instructions, got %d", len(program.Instructions))
	}

	if len(program.Functions) != 0 {
		t.Errorf("expected empty functions, got %d", len(program.Functions))
	}
}

// Test IR function creation
func TestIRFunctionCreation(t *testing.T) {
	fn := NewIRFunction("test")

	if fn.Name != "test" {
		t.Errorf("expected function name 'test', got %q", fn.Name)
	}

	if len(fn.Instructions) != 0 {
		t.Errorf("expected empty instructions, got %d", len(fn.Instructions))
	}
}

// Benchmark compiler performance
func BenchmarkCompiler(b *testing.B) {
	src := `
fn fibonacci(n) {
    if (n <= 1) {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

let result = fibonacci(10)
`

	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compiler := NewCompiler()
		_, _ = compiler.Compile(module)
	}
}

// Test optimizer
func TestOptimizer(t *testing.T) {
	program := NewIRProgram()
	program.Instructions = append(program.Instructions, NewInstruction(OpLoadConst, 1))
	program.Instructions = append(program.Instructions, NewInstruction(OpLoadConst, 2))
	program.Instructions = append(program.Instructions, NewInstruction(OpAdd))

	optimizer := NewOptimizer(program)
	err := optimizer.Optimize()

	if err != nil {
		t.Errorf("optimization failed: %v", err)
	}

	if len(program.Instructions) < 3 {
		t.Errorf("expected at least 3 instructions after optimization")
	}
}

// Test multiple statements
func TestCompileMultipleStatements(t *testing.T) {
	src := `
let x = 10
let y = 20
let z = x + y
`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	if len(program.Instructions) < 5 {
		t.Errorf("expected at least 5 instructions, got %d", len(program.Instructions))
	}
}

// Test nested blocks
func TestCompileNestedBlocks(t *testing.T) {
	src := `
let x = 10
if (x > 5) {
    let y = 20
    if (y > 15) {
        let z = 30
    }
}`
	program := parseAndCompile(t, src)

	if program == nil {
		t.Errorf("expected program, got nil")
	}

	// Check for multiple jump instructions (nested if)
	jumpCount := 0
	for _, ins := range program.Instructions {
		if ins.Opcode == OpJumpIfFalse || ins.Opcode == OpJump {
			jumpCount++
		}
	}

	if jumpCount < 2 {
		t.Errorf("expected at least 2 jump instructions, got %d", jumpCount)
	}
}
