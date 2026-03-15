package lexer

import "fmt"

// TokenType represents the type of a lexical token
type TokenType string

const (
	// Keywords (50개)
	FN       TokenType = "FN"
	LET      TokenType = "LET"
	CONST    TokenType = "CONST"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	MATCH    TokenType = "MATCH"
	FOR      TokenType = "FOR"
	WHILE    TokenType = "WHILE"
	LOOP     TokenType = "LOOP"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	RETURN   TokenType = "RETURN"
	ASYNC    TokenType = "ASYNC"
	AWAIT    TokenType = "AWAIT"
	IMPORT   TokenType = "IMPORT"
	EXPORT   TokenType = "EXPORT"
	FROM     TokenType = "FROM"
	STRUCT   TokenType = "STRUCT"
	ENUM     TokenType = "ENUM"
	TRAIT    TokenType = "TRAIT"
	TYPE     TokenType = "TYPE"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	NULL     TokenType = "NULL"
	IN       TokenType = "IN"
	OF       TokenType = "OF"
	AS       TokenType = "AS"
	IS       TokenType = "IS"
	PUB      TokenType = "PUB"
	MUT      TokenType = "MUT"
	SELF     TokenType = "SELF"
	SUPER    TokenType = "SUPER"
	IMPL     TokenType = "IMPL"

	// Exception Handling
	TRY     TokenType = "TRY"
	CATCH   TokenType = "CATCH"
	THROW   TokenType = "THROW"
	FINALLY TokenType = "FINALLY"

	// Phase 5: Minimal .free Format
	INPUT  TokenType = "INPUT"
	OUTPUT TokenType = "OUTPUT"
	INTENT TokenType = "INTENT"

	// Security & Style
	SECRET TokenType = "SECRET"
	STYLE  TokenType = "STYLE"

	// Testing & GraphQL
	TEST           TokenType = "TEST"
	EXPECT         TokenType = "EXPECT"
	SCHEMA         TokenType = "SCHEMA"
	RESOLVER       TokenType = "RESOLVER"
	QUERY          TokenType = "QUERY"
	MUTATION       TokenType = "MUTATION"
	FORMAT_POLICY  TokenType = "FORMAT_POLICY"

	// Identifiers & Literals
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"
	CHAR   TokenType = "CHAR"
	REGEX  TokenType = "REGEX"

	// Operators (Arithmetic)
	PLUS   TokenType = "PLUS"   // +
	MINUS  TokenType = "MINUS"  // -
	STAR   TokenType = "STAR"   // *
	SLASH  TokenType = "SLASH"  // /
	PERCENT TokenType = "PERCENT" // %
	POWER  TokenType = "POWER"  // **

	// Comparison
	EQ TokenType = "EQ" // ==
	NE TokenType = "NE" // !=
	LT TokenType = "LT" // <
	GT TokenType = "GT" // >
	LE TokenType = "LE" // <=
	GE TokenType = "GE" // >=

	// Logical
	AND TokenType = "AND" // &&
	OR  TokenType = "OR"  // ||
	NOT TokenType = "NOT" // !

	// Bitwise
	BIT_AND TokenType = "BIT_AND" // &
	BIT_OR  TokenType = "BIT_OR"  // |
	BIT_XOR TokenType = "BIT_XOR" // ^
	BIT_NOT TokenType = "BIT_NOT" // ~
	SHL     TokenType = "SHL"     // <<
	SHR     TokenType = "SHR"     // >>

	// Assignment
	ASSIGN        TokenType = "ASSIGN"        // =
	PLUS_ASSIGN   TokenType = "PLUS_ASSIGN"   // +=
	MINUS_ASSIGN  TokenType = "MINUS_ASSIGN"  // -=
	STAR_ASSIGN   TokenType = "STAR_ASSIGN"   // *=
	SLASH_ASSIGN  TokenType = "SLASH_ASSIGN"  // /=
	PERCENT_ASSIGN TokenType = "PERCENT_ASSIGN" // %=

	// Range & Special
	RANGE     TokenType = "RANGE"     // ..
	RANGE_INC TokenType = "RANGE_INC" // ..=

	// Delimiters
	DOT         TokenType = "DOT"         // .
	COLON_COLON TokenType = "COLON_COLON" // ::
	QUESTION    TokenType = "QUESTION"   // ?
	PIPE_GT     TokenType = "PIPE_GT"    // |>

	LPAREN    TokenType = "LPAREN"    // (
	RPAREN    TokenType = "RPAREN"    // )
	LBRACKET  TokenType = "LBRACKET"  // [
	RBRACKET  TokenType = "RBRACKET"  // ]
	LBRACE    TokenType = "LBRACE"    // {
	RBRACE    TokenType = "RBRACE"    // }
	COMMA     TokenType = "COMMA"     // ,
	SEMICOLON TokenType = "SEMICOLON" // ;
	COLON     TokenType = "COLON"     // :
	ARROW     TokenType = "ARROW"     // ->
	FAT_ARROW TokenType = "FAT_ARROW" // =>
	AT        TokenType = "AT"        // @
	HASH      TokenType = "HASH"      // #

	// Special
	EOF     TokenType = "EOF"
	NEWLINE TokenType = "NEWLINE"
	COMMENT TokenType = "COMMENT"
	ILLEGAL TokenType = "ILLEGAL"
)

// Token represents a lexical token
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Value: %q, Line: %d, Column: %d}",
		t.Type, t.Value, t.Line, t.Column)
}

// KeywordMap maps string keywords to token types
var KeywordMap = map[string]TokenType{
	// v1 기본 (33개)
	"fn":       FN,
	"let":      LET,
	"const":    CONST,
	"if":       IF,
	"else":     ELSE,
	"match":    MATCH,
	"for":      FOR,
	"while":    WHILE,
	"loop":     LOOP,
	"break":    BREAK,
	"continue": CONTINUE,
	"return":   RETURN,
	"async":    ASYNC,
	"await":    AWAIT,
	"import":   IMPORT,
	"export":   EXPORT,
	"from":     FROM,
	"struct":   STRUCT,
	"enum":     ENUM,
	"trait":    TRAIT,
	"type":     TYPE,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"in":       IN,
	"of":       OF,
	"as":       AS,
	"is":       IS,
	"pub":      PUB,
	"mut":      MUT,
	"self":     SELF,
	"super":    SUPER,
	"impl":     IMPL,

	// Exception Handling (4개)
	"try":     TRY,
	"catch":   CATCH,
	"throw":   THROW,
	"finally": FINALLY,

	// Phase 5 (3개)
	"input":  INPUT,
	"output": OUTPUT,
	"intent": INTENT,

	// Security & Style
	"secret": SECRET,
	"style":  STYLE,

	// Testing & GraphQL
	"test":           TEST,
	"expect":         EXPECT,
	"schema":         SCHEMA,
	"resolver":       RESOLVER,
	"query":          QUERY,
	"mutation":       MUTATION,
	"format_policy":  FORMAT_POLICY,
}

// IsKeyword checks if a string is a keyword
func IsKeyword(s string) bool {
	_, ok := KeywordMap[s]
	return ok
}

// GetKeywordType returns the token type for a keyword, or IDENT if not a keyword
func GetKeywordType(s string) TokenType {
	if tt, ok := KeywordMap[s]; ok {
		return tt
	}
	return IDENT
}
