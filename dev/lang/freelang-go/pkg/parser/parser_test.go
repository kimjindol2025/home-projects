package parser

import (
	"fmt"
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/ast"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
)

func TestLetStatement(t *testing.T) {
	input := "let x = 5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	if !testLetStatement(t, stmt, "x") {
		return
	}

	val := stmt.(*ast.LetStatement).Value
	if !testIntegerLiteral(t, val, 5) {
		return
	}
}

func TestReturnStatement(t *testing.T) {
	input := "return 5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	retStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ReturnStatement. got=%T", stmt)
	}

	if !testIntegerLiteral(t, retStmt.ReturnValue, 5) {
		return
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	if !testIdentifier(t, expStmt.Expression, "foobar") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	if !testIntegerLiteral(t, expStmt.Expression, 5) {
		return
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	str, ok := expStmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.StringLiteral. got=%T", expStmt.Expression)
	}

	if str.Value != "hello world" {
		t.Errorf("str.Value is not 'hello world'. got=%q", str.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	if !testBooleanLiteral(t, expStmt.Expression, true) {
		return
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!true", "!", 1}, // true → 1 for simplicity
		{"-5", "-", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		prefixExpr, ok := expStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expr is not *ast.PrefixExpression. got=%T", expStmt.Expression)
		}

		if prefixExpr.Operator != tt.operator {
			t.Errorf("prefixExpr.Operator is not %s. got=%s", tt.operator, prefixExpr.Operator)
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		infixExpr, ok := expStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expr is not *ast.InfixExpression. got=%T", expStmt.Expression)
		}

		if !testIntegerLiteral(t, infixExpr.Left, tt.leftValue) {
			return
		}

		if infixExpr.Operator != tt.operator {
			t.Errorf("infixExpr.Operator is not %s. got=%s", tt.operator, infixExpr.Operator)
		}

		if !testIntegerLiteral(t, infixExpr.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + 5 * 5", "(5 + (5 * 5))"},
		{"5 + 5 + 5", "((5 + 5) + 5)"},
		{"5 * 5 + 5", "((5 * 5) + 5)"},
		{"5 + 5 * 2 + 6 / 2 - 3", "(((5 + (5 * 2)) + (6 / 2)) - 3)"},
		{"-5 + 5", "((-5) + 5)"},
		{"!true == false", "((!true) == false)"},
		{"3 > 5 == false", "((3 > 5) == false)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.Statements[0].(*ast.ExpressionStatement).Expression.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x > 5) { return x; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt := program.Statements[0]
	ifStmt, ok := stmt.(*ast.IfStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.IfStatement. got=%T", stmt)
	}

	if ifStmt.Condition == nil {
		t.Fatalf("ifStmt.Condition is nil")
	}

	if ifStmt.Consequence == nil {
		t.Fatalf("ifStmt.Consequence is nil")
	}

	if len(ifStmt.Consequence.Statements) != 1 {
		t.Errorf("consequence has wrong number of statements. got=%d", len(ifStmt.Consequence.Statements))
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x > 5) { return x; } else { return 0; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt := program.Statements[0]
	ifStmt, ok := stmt.(*ast.IfStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.IfStatement. got=%T", stmt)
	}

	if ifStmt.Alternative == nil {
		t.Fatalf("ifStmt.Alternative is nil")
	}

	if len(ifStmt.Alternative.Statements) != 1 {
		t.Errorf("alternative has wrong number of statements. got=%d", len(ifStmt.Alternative.Statements))
	}
}

func TestForStatement(t *testing.T) {
	input := `for x in arr { x + 1; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	forStmt, ok := stmt.(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ForStatement. got=%T", stmt)
	}

	if forStmt.LoopVar.Value != "x" {
		t.Errorf("forStmt.LoopVar.Value is not 'x'. got=%q", forStmt.LoopVar.Value)
	}

	if forStmt.Body == nil {
		t.Fatalf("forStmt.Body is nil")
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `func(x, y) { return x + y; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	funcLit, ok := expStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionLiteral. got=%T", expStmt.Expression)
	}

	if len(funcLit.Parameters) != 2 {
		t.Errorf("funcLit has wrong number of parameters. got=%d", len(funcLit.Parameters))
	}

	if funcLit.Parameters[0].Value != "x" {
		t.Errorf("first parameter is not 'x'. got=%q", funcLit.Parameters[0].Value)
	}

	if funcLit.Parameters[1].Value != "y" {
		t.Errorf("second parameter is not 'y'. got=%q", funcLit.Parameters[1].Value)
	}
}

func TestFunctionCall(t *testing.T) {
	input := `add(2, 3);`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	callExpr, ok := expStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpression. got=%T", expStmt.Expression)
	}

	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}

	if len(callExpr.Arguments) != 2 {
		t.Errorf("callExpr has wrong number of arguments. got=%d", len(callExpr.Arguments))
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2, 3]`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	arrLit, ok := expStmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.ArrayLiteral. got=%T", expStmt.Expression)
	}

	if len(arrLit.Elements) != 3 {
		t.Errorf("arrLit has wrong number of elements. got=%d", len(arrLit.Elements))
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{a: 1, b: 2}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	hashLit, ok := expStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.HashLiteral. got=%T", expStmt.Expression)
	}

	if len(hashLit.Pairs) != 2 {
		t.Errorf("hashLit has wrong number of pairs. got=%d", len(hashLit.Pairs))
	}
}

func TestIndexExpression(t *testing.T) {
	input := `arr[0]`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	indexExpr, ok := expStmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expr is not *ast.IndexExpression. got=%T", expStmt.Expression)
	}

	if !testIdentifier(t, indexExpr.Left, "arr") {
		return
	}

	if !testIntegerLiteral(t, indexExpr.Index, 0) {
		return
	}
}

// Helper functions

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %q. got=%q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %q. got=%q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intlit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if intlit.Value != value {
		t.Errorf("intlit.Value not %d. got=%d", value, intlit.Value)
		return false
	}

	if intlit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intlit.TokenLiteral() not %q. got=%q", fmt.Sprintf("%d", value), intlit.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %q. got=%q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %q. got=%q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %v. got=%v", value, bo.Value)
		return false
	}

	return true
}
