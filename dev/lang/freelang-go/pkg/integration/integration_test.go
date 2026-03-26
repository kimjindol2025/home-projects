package integration

import (
	"fmt"
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/checker"
	"github.com/kimjindol2025/freelang-go/pkg/compiler"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
	"github.com/kimjindol2025/freelang-go/pkg/object"
	"github.com/kimjindol2025/freelang-go/pkg/parser"
	"github.com/kimjindol2025/freelang-go/pkg/vm"
)

// evalCode compiles and executes code, returning the result
func evalCode(source string) (object.Object, error) {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse error: %s", p.Errors()[0])
	}

	ch := checker.New()
	ch.Check(program)

	if len(ch.Errors()) > 0 {
		return nil, fmt.Errorf("check error: %s", ch.Errors()[0])
	}

	c := compiler.New(ch.GlobalScope())
	c.Compile(program)

	if len(c.Errors()) > 0 {
		return nil, fmt.Errorf("compile error: %s", c.Errors()[0])
	}

	v := vm.New(c.Bytecode(), c.Constants())
	if err := v.Run(); err != nil {
		return nil, err
	}

	return v.LastPoppedStackElem(), nil
}

// Test Case 1: Simple Expressions
func TestSimpleExpressions(t *testing.T) {
	tests := []struct {
		code   string
		expect interface{}
	}{
		{"5;", int64(5)},
		{"10;", int64(10)},
		{`"hello";`, "hello"},
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		switch expected := tt.expect.(type) {
		case int64:
			intResult, ok := result.(*object.Integer)
			if !ok || intResult.Value != expected {
				t.Errorf("code '%s': expected %d, got %v", tt.code, expected, result)
			}
		case string:
			strResult, ok := result.(*object.String)
			if !ok || strResult.Value != expected {
				t.Errorf("code '%s': expected %s, got %v", tt.code, expected, result)
			}
		case bool:
			boolResult, ok := result.(*object.Boolean)
			if !ok || boolResult.Value != expected {
				t.Errorf("code '%s': expected %v, got %v", tt.code, expected, result)
			}
		}
	}
}

// Test Case 2: Arithmetic Operations
func TestArithmetic(t *testing.T) {
	tests := []struct {
		code   string
		expect int64
	}{
		{"1 + 2;", 3},
		{"5 - 3;", 2},
		{"4 * 5;", 20},
		{"20 / 4;", 5},
		{"10 % 3;", 1},
		{"2 ** 3;", 8},
		{"-5;", -5},
		{"+10;", 10},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		intResult, ok := result.(*object.Integer)
		if !ok || intResult.Value != tt.expect {
			t.Errorf("code '%s': expected %d, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 3: Comparisons
func TestComparisons(t *testing.T) {
	tests := []struct {
		code   string
		expect bool
	}{
		{"5 == 5;", true},
		{"5 != 5;", false},
		{"5 < 10;", true},
		{"5 > 10;", false},
		{"5 <= 5;", true},
		{"5 >= 5;", true},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		boolResult, ok := result.(*object.Boolean)
		if !ok || boolResult.Value != tt.expect {
			t.Errorf("code '%s': expected %v, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 4: Logical Operations
func TestLogical(t *testing.T) {
	tests := []struct {
		code   string
		expect bool
	}{
		{"true && true;", true},
		{"true && false;", false},
		{"false && true;", false},
		{"true || false;", true},
		{"false || false;", false},
		{"!true;", false},
		{"!false;", true},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		boolResult, ok := result.(*object.Boolean)
		if !ok || boolResult.Value != tt.expect {
			t.Errorf("code '%s': expected %v, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 5: String Operations
func TestStringOps(t *testing.T) {
	tests := []struct {
		code   string
		expect string
	}{
		{`"hello" + " " + "world";`, "hello world"},
		{`"foo" + "bar";`, "foobar"},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		strResult, ok := result.(*object.String)
		if !ok || strResult.Value != tt.expect {
			t.Errorf("code '%s': expected %s, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 6: Arrays
func TestArrays(t *testing.T) {
	code := "[1, 2, 3];"
	result, err := evalCode(code)
	if err != nil {
		t.Fatalf("code '%s': %v", code, err)
	}

	arrResult, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("expected array, got %T", result)
	}

	if len(arrResult.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arrResult.Elements))
	}

	for i, expected := range []int64{1, 2, 3} {
		elem, ok := arrResult.Elements[i].(*object.Integer)
		if !ok || elem.Value != expected {
			t.Errorf("element %d: expected %d, got %v", i, expected, arrResult.Elements[i])
		}
	}
}

// Test Case 7: Array Indexing
func TestArrayIndexing(t *testing.T) {
	tests := []struct {
		code   string
		expect int64
	}{
		{"[1, 2, 3][0];", 1},
		{"[1, 2, 3][1];", 2},
		{"[1, 2, 3][2];", 3},
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		intResult, ok := result.(*object.Integer)
		if !ok || intResult.Value != tt.expect {
			t.Errorf("code '%s': expected %d, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 8: Hashes
func TestHashes(t *testing.T) {
	code := `{a: 1, b: 2, c: 3};`
	result, err := evalCode(code)
	if err != nil {
		t.Fatalf("code '%s': %v", code, err)
	}

	hashResult, ok := result.(*object.Hash)
	if !ok {
		t.Fatalf("expected hash, got %T", result)
	}

	if len(hashResult.Pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(hashResult.Pairs))
	}

	expected := map[string]int64{"a": 1, "b": 2, "c": 3}
	for key, expectedVal := range expected {
		val, ok := hashResult.Pairs[key]
		if !ok {
			t.Errorf("expected key %q in hash", key)
			continue
		}

		intVal, ok := val.(*object.Integer)
		if !ok || intVal.Value != expectedVal {
			t.Errorf("hash[%q]: expected %d, got %v", key, expectedVal, val)
		}
	}
}

// Test Case 9: Hash Indexing
func TestHashIndexing(t *testing.T) {
	code := `{a: 1, b: 2}["a"];`
	result, err := evalCode(code)
	if err != nil {
		t.Fatalf("code '%s': %v", code, err)
	}

	intResult, ok := result.(*object.Integer)
	if !ok || intResult.Value != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

// Test Case 10: Operator Precedence
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		code   string
		expect int64
	}{
		{"2 + 3 * 4;", 14},      // 3 * 4 first
		{"(2 + 3) * 4;", 20},    // 2 + 3 first
		{"10 - 5 - 2;", 3},      // left-to-right
		{"2 ** 3;", 8},          // power
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		intResult, ok := result.(*object.Integer)
		if !ok || intResult.Value != tt.expect {
			t.Errorf("code '%s': expected %d, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 11: Bitwise Operations
func TestBitwise(t *testing.T) {
	tests := []struct {
		code   string
		expect int64
	}{
		{"5 & 3;", 1},           // 101 & 011 = 001
		{"5 | 3;", 7},           // 101 | 011 = 111
		{"5 ^ 3;", 6},           // 101 ^ 011 = 110
		{"~5;", -6},             // bitwise not
		{"4 << 1;", 8},          // left shift
		{"8 >> 1;", 4},          // right shift
	}

	for _, tt := range tests {
		result, err := evalCode(tt.code)
		if err != nil {
			t.Errorf("code '%s': %v", tt.code, err)
			continue
		}

		intResult, ok := result.(*object.Integer)
		if !ok || intResult.Value != tt.expect {
			t.Errorf("code '%s': expected %d, got %v", tt.code, tt.expect, result)
		}
	}
}

// Test Case 12: Mixed Operations
func TestMixed(t *testing.T) {
	code := "5 + 3 * 2 - 1;"
	result, err := evalCode(code)
	if err != nil {
		t.Fatalf("code '%s': %v", code, err)
	}

	intResult, ok := result.(*object.Integer)
	if !ok || intResult.Value != 10 {
		t.Errorf("expected 10, got %v", result)
	}
}
