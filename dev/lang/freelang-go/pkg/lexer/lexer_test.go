package lexer

import (
	"testing"

	"github.com/kimjindol2025/freelang-go/pkg/token"
)

func TestSimpleTokens(t *testing.T) {
	input := "let x = 5;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := "if else for while return func let const var"

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.IF},
		{token.ELSE},
		{token.FOR},
		{token.WHILE},
		{token.RETURN},
		{token.FUNC},
		{token.LET},
		{token.CONST},
		{token.VAR},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q (%q)", i, tt.expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := "+ - * / % == != < > <= >= && || !"

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.PLUS},
		{token.MINUS},
		{token.MULTIPLY},
		{token.DIVIDE},
		{token.MODULO},
		{token.EQ},
		{token.NOT_EQ},
		{token.LT},
		{token.GT},
		{token.LE},
		{token.GE},
		{token.AND},
		{token.OR},
		{token.NOT},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
	}
}

func TestStrings(t *testing.T) {
	input := `"hello" 'world'`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "hello"},
		{token.STRING, "world"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := "123 45.67 1e10"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.INT, "45.67"},
		{token.INT, "1e10"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestComplexCode(t *testing.T) {
	input := `func add(a, b) { return a + b; }`

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.FUNC},
		{token.IDENTIFIER},
		{token.LPAREN},
		{token.IDENTIFIER},
		{token.COMMA},
		{token.IDENTIFIER},
		{token.RPAREN},
		{token.LBRACE},
		{token.RETURN},
		{token.IDENTIFIER},
		{token.PLUS},
		{token.IDENTIFIER},
		{token.SEMICOLON},
		{token.RBRACE},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q (%q)", i, tt.expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	input := `x = 5; // comment
y = 10; /* block */`

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.IDENTIFIER},
		{token.ASSIGN},
		{token.INT},
		{token.SEMICOLON},
		{token.IDENTIFIER},
		{token.ASSIGN},
		{token.INT},
		{token.SEMICOLON},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
	}
}
