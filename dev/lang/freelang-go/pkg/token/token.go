package token

import "fmt"

// TokenType represents the type of a token in FreeLang
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}", t.Type, t.Literal, t.Line, t.Column)
}

// Token types for FreeLang
const (
	// Special tokens
	ILLEGAL   TokenType = "ILLEGAL"
	EOF       TokenType = "EOF"
	NEWLINE   TokenType = "NEWLINE"
	INDENT    TokenType = "INDENT"
	DEDENT    TokenType = "DEDENT"

	// Literals
	IDENTIFIER TokenType = "IDENTIFIER"
	INT        TokenType = "INT"
	STRING     TokenType = "STRING"
	FLOAT      TokenType = "FLOAT"
	BOOL       TokenType = "BOOL" // true, false

	// Keywords
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	ELIF     TokenType = "ELIF"
	FOR      TokenType = "FOR"
	WHILE    TokenType = "WHILE"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	RETURN   TokenType = "RETURN"
	FUNC     TokenType = "FUNC"
	STRUCT   TokenType = "STRUCT"
	ENUM     TokenType = "ENUM"
	INTERFACE TokenType = "INTERFACE"
	IMPORT   TokenType = "IMPORT"
	EXPORT   TokenType = "EXPORT"
	CONST    TokenType = "CONST"
	VAR      TokenType = "VAR"
	IN       TokenType = "IN"
	OF       TokenType = "OF"
	ASYNC    TokenType = "ASYNC"
	AWAIT    TokenType = "AWAIT"
	LET      TokenType = "LET"
	TRY      TokenType = "TRY"
	CATCH    TokenType = "CATCH"
	FINALLY  TokenType = "FINALLY"
	THROW    TokenType = "THROW"
	CLASS    TokenType = "CLASS"
	EXTENDS  TokenType = "EXTENDS"
	IMPLEMENTS TokenType = "IMPLEMENTS"
	THIS     TokenType = "THIS"
	SUPER    TokenType = "SUPER"
	STATIC   TokenType = "STATIC"
	PUBLIC   TokenType = "PUBLIC"
	PRIVATE  TokenType = "PRIVATE"
	PROTECTED TokenType = "PROTECTED"
	READONLY TokenType = "READONLY"
	ABSTRACT TokenType = "ABSTRACT"
	TYPE     TokenType = "TYPE"
	AS       TokenType = "AS"
	NEW      TokenType = "NEW"
	DELETE   TokenType = "DELETE"
	INSTANCEOF TokenType = "INSTANCEOF"
	TYPEOF   TokenType = "TYPEOF"
	VOID     TokenType = "VOID"
	NULL     TokenType = "NULL"
	UNDEFINED TokenType = "UNDEFINED"

	// Operators
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	MULTIPLY TokenType = "MULTIPLY"
	DIVIDE   TokenType = "DIVIDE"
	MODULO   TokenType = "MODULO"
	POWER    TokenType = "POWER"

	// Assignment operators
	ASSIGN      TokenType = "ASSIGN"
	PLUS_ASSIGN TokenType = "PLUS_ASSIGN"
	MINUS_ASSIGN TokenType = "MINUS_ASSIGN"
	MUL_ASSIGN   TokenType = "MUL_ASSIGN"
	DIV_ASSIGN   TokenType = "DIV_ASSIGN"
	MOD_ASSIGN   TokenType = "MOD_ASSIGN"

	// Comparison operators
	EQ      TokenType = "EQ"       // ==
	NOT_EQ  TokenType = "NOT_EQ"   // !=
	LT      TokenType = "LT"       // <
	LE      TokenType = "LE"       // <=
	GT      TokenType = "GT"       // >
	GE      TokenType = "GE"       // >=

	// Logical operators
	AND TokenType = "AND" // &&
	OR  TokenType = "OR"  // ||
	NOT TokenType = "NOT" // !

	// Bitwise operators
	BIT_AND TokenType = "BIT_AND"      // &
	BIT_OR  TokenType = "BIT_OR"       // |
	BIT_XOR TokenType = "BIT_XOR"      // ^
	BIT_NOT TokenType = "BIT_NOT"      // ~
	LSHIFT  TokenType = "LSHIFT"       // <<
	RSHIFT  TokenType = "RSHIFT"       // >>

	// Delimiters
	LPAREN    TokenType = "LPAREN"     // (
	RPAREN    TokenType = "RPAREN"     // )
	LBRACE    TokenType = "LBRACE"     // {
	RBRACE    TokenType = "RBRACE"     // }
	LBRACKET  TokenType = "LBRACKET"   // [
	RBRACKET  TokenType = "RBRACKET"   // ]

	// Punctuation
	COMMA     TokenType = "COMMA"      // ,
	COLON     TokenType = "COLON"      // :
	SEMICOLON TokenType = "SEMICOLON"  // ;
	DOT       TokenType = "DOT"        // .
	ARROW     TokenType = "ARROW"      // =>
	LAMBDA    TokenType = "LAMBDA"     // ->
	QUESTION  TokenType = "QUESTION"   // ?
	SPREAD    TokenType = "SPREAD"     // ...
	PIPE      TokenType = "PIPE"       // |

	// Special operators
	INCREMENT TokenType = "INCREMENT"  // ++
	DECREMENT TokenType = "DECREMENT"  // --
	OPTIONAL  TokenType = "OPTIONAL"   // ?.
)

// Keywords map
var keywords = map[string]TokenType{
	"if":         IF,
	"else":       ELSE,
	"elif":       ELIF,
	"for":        FOR,
	"while":      WHILE,
	"break":      BREAK,
	"continue":   CONTINUE,
	"return":     RETURN,
	"func":       FUNC,
	"struct":     STRUCT,
	"enum":       ENUM,
	"interface":  INTERFACE,
	"import":     IMPORT,
	"export":     EXPORT,
	"const":      CONST,
	"var":        VAR,
	"let":        LET,
	"in":         IN,
	"of":         OF,
	"async":      ASYNC,
	"await":      AWAIT,
	"try":        TRY,
	"catch":      CATCH,
	"finally":    FINALLY,
	"throw":      THROW,
	"class":      CLASS,
	"extends":    EXTENDS,
	"implements": IMPLEMENTS,
	"this":       THIS,
	"super":      SUPER,
	"static":     STATIC,
	"public":     PUBLIC,
	"private":    PRIVATE,
	"protected":  PROTECTED,
	"readonly":   READONLY,
	"abstract":   ABSTRACT,
	"type":       TYPE,
	"as":         AS,
	"new":        NEW,
	"delete":     DELETE,
	"instanceof": INSTANCEOF,
	"typeof":     TYPEOF,
	"void":       VOID,
	"null":       NULL,
	"undefined":  UNDEFINED,
	"true":       BOOL,
	"false":      BOOL,
}

// LookupIdent returns the token type for a given identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
