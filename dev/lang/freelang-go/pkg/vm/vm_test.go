package vm

import (
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/checker"
	"github.com/kimjindol2025/freelang-go/pkg/compiler"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
	"github.com/kimjindol2025/freelang-go/pkg/object"
	"github.com/kimjindol2025/freelang-go/pkg/parser"
)

func TestIntegerArithmetic(t *testing.T) {
	input := "5 + 10;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	if result == nil {
		t.Errorf("expected result, got nil")
		return
	}

	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Errorf("expected integer, got %T", result)
		return
	}

	if intResult.Value != 15 {
		t.Errorf("expected 15, got %d", intResult.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" + " " + "world";`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	strResult, ok := result.(*object.String)
	if !ok {
		t.Errorf("expected string, got %T", result)
		return
	}

	if strResult.Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", strResult.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "5 > 3;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	boolResult, ok := result.(*object.Boolean)
	if !ok {
		t.Errorf("expected boolean, got %T", result)
		return
	}

	if !boolResult.Value {
		t.Errorf("expected true, got false")
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2, 3];"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	arrResult, ok := result.(*object.Array)
	if !ok {
		t.Errorf("expected array, got %T", result)
		return
	}

	if len(arrResult.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arrResult.Elements))
		return
	}

	for i, expected := range []int64{1, 2, 3} {
		elem, ok := arrResult.Elements[i].(*object.Integer)
		if !ok {
			t.Errorf("expected integer at index %d, got %T", i, arrResult.Elements[i])
			continue
		}
		if elem.Value != expected {
			t.Errorf("expected %d at index %d, got %d", expected, i, elem.Value)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{a: 1, b: 2};`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	hashResult, ok := result.(*object.Hash)
	if !ok {
		t.Errorf("expected hash, got %T", result)
		return
	}

	if len(hashResult.Pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(hashResult.Pairs))
	}
}

func TestArrayIndex(t *testing.T) {
	input := "[1, 2, 3][1];"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Errorf("expected integer, got %T", result)
		return
	}

	if intResult.Value != 2 {
		t.Errorf("expected 2, got %d", intResult.Value)
	}
}

func TestLogicalOperators(t *testing.T) {
	input := "true && false;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	boolResult, ok := result.(*object.Boolean)
	if !ok {
		t.Errorf("expected boolean, got %T", result)
		return
	}

	if boolResult.Value {
		t.Errorf("expected false, got true")
	}
}

func TestNotOperator(t *testing.T) {
	input := "!true;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	boolResult, ok := result.(*object.Boolean)
	if !ok {
		t.Errorf("expected boolean, got %T", result)
		return
	}

	if boolResult.Value {
		t.Errorf("expected false, got true")
	}
}

func TestNegation(t *testing.T) {
	input := "-5;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Errorf("expected integer, got %T", result)
		return
	}

	if intResult.Value != -5 {
		t.Errorf("expected -5, got %d", intResult.Value)
	}
}

func TestMultiplication(t *testing.T) {
	input := "3 * 4;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Errorf("expected integer, got %T", result)
		return
	}

	if intResult.Value != 12 {
		t.Errorf("expected 12, got %d", intResult.Value)
	}
}

func TestDivision(t *testing.T) {
	input := "20 / 4;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	ch := checker.New()
	ch.Check(program)

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	vm := New(c.Bytecode(), c.Constants())
	if err := vm.Run(); err != nil {
		t.Fatalf("VM error: %v", err)
	}

	result := vm.LastPoppedStackElem()
	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Errorf("expected integer, got %T", result)
		return
	}

	if intResult.Value != 5 {
		t.Errorf("expected 5, got %d", intResult.Value)
	}
}
