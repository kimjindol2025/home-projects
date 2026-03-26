package compiler

import (
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/checker"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
	"github.com/kimjindol2025/freelang-go/pkg/parser"
)

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	if len(compiler.Constants()) != 1 {
		t.Errorf("expected 1 constant, got %d", len(compiler.Constants()))
	}

	if compiler.Constants()[0] != int64(5) {
		t.Errorf("expected constant 5, got %v", compiler.Constants()[0])
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello";`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	if len(compiler.Constants()) != 1 {
		t.Errorf("expected 1 constant, got %d", len(compiler.Constants()))
	}

	if compiler.Constants()[0] != "hello" {
		t.Errorf("expected constant 'hello', got %v", compiler.Constants()[0])
	}
}

func TestBooleanLiteral(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Boolean literals don't add to constants pool
	if len(compiler.Bytecode()) == 0 {
		t.Errorf("expected bytecode for boolean literal")
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2, 3];"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constants for 1, 2, 3
	if len(compiler.Constants()) != 3 {
		t.Errorf("expected 3 constants, got %d", len(compiler.Constants()))
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{a: 1, b: 2};`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constants for keys and values
	if len(compiler.Constants()) < 2 {
		t.Errorf("expected at least 2 constants, got %d", len(compiler.Constants()))
	}
}

func TestArithmeticExpression(t *testing.T) {
	input := "1 + 2;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constants for 1 and 2
	if len(compiler.Constants()) != 2 {
		t.Errorf("expected 2 constants, got %d", len(compiler.Constants()))
	}

	// Should have ADD instruction in bytecode
	if len(compiler.Bytecode()) == 0 {
		t.Errorf("expected bytecode")
	}
}

func TestLetStatement(t *testing.T) {
	input := "let x = 5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constant for 5
	if len(compiler.Constants()) != 1 {
		t.Errorf("expected 1 constant, got %d", len(compiler.Constants()))
	}
}

func TestOpcodeString(t *testing.T) {
	tests := []struct {
		opcode   OpCode
		expected string
	}{
		{OpConstant, "CONSTANT"},
		{OpTrue, "TRUE"},
		{OpFalse, "FALSE"},
		{OpAdd, "ADD"},
		{OpSubtract, "SUBTRACT"},
		{OpCall, "CALL"},
	}

	for _, tt := range tests {
		if tt.opcode.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.opcode.String())
		}
	}
}

func TestEncodeDecodeInstructions(t *testing.T) {
	tests := []struct {
		opcode   OpCode
		operands []int
	}{
		{OpConstant, []int{65534}},
		{OpAdd, []int{}},
		{OpCall, []int{2}},
	}

	for _, tt := range tests {
		encoded := EncodeInstructions(tt.opcode, tt.operands...)

		decoded, width := DecodeInstructions(encoded)

		if width != len(encoded) {
			t.Errorf("expected width %d, got %d", len(encoded), width)
		}

		if len(decoded) != len(tt.operands) {
			t.Errorf("expected %d operands, got %d", len(tt.operands), len(decoded))
		}

		for i, op := range decoded {
			if op != tt.operands[i] {
				t.Errorf("operand %d: expected %d, got %d", i, tt.operands[i], op)
			}
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	input := "-5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constant for 5
	if len(compiler.Constants()) != 1 {
		t.Errorf("expected 1 constant, got %d", len(compiler.Constants()))
	}
}

func TestComparisonExpression(t *testing.T) {
	input := "1 == 2;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Should have constants for 1 and 2
	if len(compiler.Constants()) != 2 {
		t.Errorf("expected 2 constants, got %d", len(compiler.Constants()))
	}
}

func TestLogicalExpression(t *testing.T) {
	input := "true && false;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	compiler := New(ch.GlobalScope())
	compiler.Compile(program)

	if len(compiler.Errors()) > 0 {
		t.Fatalf("compilation errors: %v", compiler.Errors())
	}

	// Boolean literals don't add to constants
	if len(compiler.Bytecode()) == 0 {
		t.Errorf("expected bytecode")
	}
}
