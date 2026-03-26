package checker

import (
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/lexer"
	"github.com/kimjindol2025/freelang-go/pkg/parser"
)

func TestLetStatement(t *testing.T) {
	input := `let x = 5; let y = x;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}

	// Verify that x and y are defined
	symbol, ok := c.symtab.Resolve("x")
	if !ok {
		t.Fatalf("variable x not found")
	}

	if !TypesEqual(symbol.Type, TypeInt) {
		t.Errorf("variable x has type %s, expected int", symbol.Type.String())
	}
}

func TestDuplicateVariableError(t *testing.T) {
	input := `let x = 5; let x = 10;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for duplicate variable")
	}

	if len(c.Errors()) != 1 {
		t.Errorf("expected 1 error, got %d", len(c.Errors()))
	}
}

func TestUndefinedVariableError(t *testing.T) {
	input := `x;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for undefined variable")
	}
}

func TestArrayTypeInference(t *testing.T) {
	input := `let arr = [1, 2, 3];`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}

	symbol, ok := c.symtab.Resolve("arr")
	if !ok {
		t.Fatalf("variable arr not found")
	}

	arrayType, ok := symbol.Type.(*ArrayType)
	if !ok {
		t.Fatalf("variable arr should be array type, got %T", symbol.Type)
	}

	if !TypesEqual(arrayType.ElementType, TypeInt) {
		t.Errorf("array element type is %s, expected int", arrayType.ElementType.String())
	}
}

func TestArrayTypeIncompatibleError(t *testing.T) {
	input := `let arr = [1, "string", 3];`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for incompatible array element types")
	}
}

func TestHashTypeInference(t *testing.T) {
	input := `let hash = {a: 1, b: 2};`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}

	symbol, ok := c.symtab.Resolve("hash")
	if !ok {
		t.Fatalf("variable hash not found")
	}

	hashType, ok := symbol.Type.(*HashType)
	if !ok {
		t.Fatalf("variable hash should be hash type, got %T", symbol.Type)
	}

	if !TypesEqual(hashType.KeyType, TypeString) && !TypesEqual(hashType.KeyType, TypeAny) {
		t.Errorf("hash key type is %s", hashType.KeyType.String())
	}

	if !TypesEqual(hashType.ValueType, TypeInt) {
		t.Errorf("hash value type is %s, expected int", hashType.ValueType.String())
	}
}

func TestIfConditionTypeCheck(t *testing.T) {
	tests := []struct {
		input      string
		shouldFail bool
	}{
		{"if (true) { 1; }", false},
		{"if (1 < 2) { 1; }", false},
		{"if (5) { 1; }", true}, // Number as condition should fail
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		c := New()
		c.Check(program)

		hasError := len(c.Errors()) > 0
		if hasError != tt.shouldFail {
			if tt.shouldFail {
				t.Errorf("expected error for: %s", tt.input)
			} else {
				t.Errorf("unexpected error for: %s\nerrors: %v", tt.input, c.Errors())
			}
		}
	}
}

func TestForLoopTypeCheck(t *testing.T) {
	input := `let arr = [1, 2, 3]; for x in arr { x; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}

func TestFunctionTypeInference(t *testing.T) {
	input := `let add = func(a, b) { a + b; };`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}

	symbol, ok := c.symtab.Resolve("add")
	if !ok {
		t.Fatalf("variable add not found")
	}

	funcType, ok := symbol.Type.(*FunctionType)
	if !ok {
		t.Fatalf("variable add should be function type, got %T", symbol.Type)
	}

	if len(funcType.ParameterTypes) != 2 {
		t.Errorf("function has %d parameters, expected 2", len(funcType.ParameterTypes))
	}
}

func TestArithmeticOperatorTypeCheck(t *testing.T) {
	tests := []struct {
		input      string
		shouldFail bool
	}{
		{"5 + 3;", false},
		{"5 + 3.5;", false},
		{`"hello" + 5;`, true}, // String + number should fail
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		c := New()
		c.Check(program)

		hasError := len(c.Errors()) > 0
		if hasError != tt.shouldFail {
			if tt.shouldFail {
				t.Errorf("expected error for: %s", tt.input)
			} else {
				t.Errorf("unexpected error for: %s\nerrors: %v", tt.input, c.Errors())
			}
		}
	}
}

func TestCallExpressionTypeCheck(t *testing.T) {
	input := `let add = func(a, b) { a + b; }; add(1, 2);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}

func TestCallExpressionArgumentCountError(t *testing.T) {
	input := `let add = func(a, b) { a + b; }; add(1);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for wrong argument count")
	}
}

func TestIndexExpressionTypeCheck(t *testing.T) {
	input := `let arr = [1, 2, 3]; arr[0];`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}

func TestIndexExpressionInvalidIndexError(t *testing.T) {
	input := `let arr = [1, 2, 3]; arr["0"];`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for string index on array")
	}
}

func TestAssignmentTypeCheck(t *testing.T) {
	input := `let x = 5; x = 10;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}

func TestAssignmentTypeIncompatibleError(t *testing.T) {
	input := `let x = 5; x = "string";`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) == 0 {
		t.Fatalf("expected error for incompatible assignment type")
	}
}

func TestBuiltinFunctions(t *testing.T) {
	input := `len([1, 2, 3]); print(5);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}

func TestScopeChaining(t *testing.T) {
	input := `let x = 5; if (true) { let y = x; y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := New()
	c.Check(program)

	if len(c.Errors()) > 0 {
		t.Fatalf("type errors: %v", c.Errors())
	}
}
