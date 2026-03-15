package lexer

import (
	"fmt"
)

// Lexer tokenizes FreeLang source code
type Lexer struct {
	input          string
	position       int       // current character index
	line           int       // current line number
	column         int       // current column number
	current        rune      // current character
	lastTokenType  TokenType // track last token for regex disambiguation
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	lex := &Lexer{
		input:         input,
		position:      0,
		line:          1,
		column:        0,
		current:       '\x00',
		lastTokenType: EOF,
	}
	lex.readChar()
	return lex
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.current = '\x00'
	} else {
		l.current = rune(l.input[l.position])
	}
	l.position++
	l.column++
}

// peekChar returns the next character without consuming it
func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return '\x00'
	}
	return rune(l.input[l.position])
}

// skipWhitespace skips spaces, tabs, and carriage returns
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\r' {
		l.readChar()
	}
}

// skipComment skips a single-line comment (// ...)
func (l *Lexer) skipComment() {
	// Skip //
	l.readChar()
	l.readChar()

	// Read until newline
	for l.current != '\n' && l.current != '\x00' {
		l.readChar()
	}
}

// skipMultiLineComment skips a multi-line comment (/* ... */)
func (l *Lexer) skipMultiLineComment() {
	// Skip /*
	l.readChar()
	l.readChar()

	for l.current != '\x00' {
		if l.current == '\n' {
			l.line++
			l.column = 0
		}

		if l.current == '*' && l.peekChar() == '/' {
			l.readChar() // *
			l.readChar() // /
			break
		}

		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	start := l.position - 1
	for l.isIdentifierChar(l.current) {
		l.readChar()
	}
	return l.input[start : l.position-1]
}

// readNumber reads a number (integer, float, or scientific notation)
func (l *Lexer) readNumber() string {
	start := l.position - 1

	// Integer part
	for l.isDigit(l.current) {
		l.readChar()
	}

	// Decimal part
	if l.current == '.' && l.isDigit(l.peekChar()) {
		l.readChar() // .
		for l.isDigit(l.current) {
			l.readChar()
		}
	}

	// Exponent (e or E)
	if l.current == 'e' || l.current == 'E' {
		l.readChar()
		if l.current == '+' || l.current == '-' {
			l.readChar()
		}
		for l.isDigit(l.current) {
			l.readChar()
		}
	}

	return l.input[start : l.position-1]
}

// readString reads a string literal with escape sequences
func (l *Lexer) readString() string {
	quote := l.current
	l.readChar() // skip opening quote

	var result string
	for l.current != quote && l.current != '\x00' {
		if l.current == '\\' {
			l.readChar()
			// Handle escape sequences
			switch l.current {
			case 'n':
				result += "\n"
			case 't':
				result += "\t"
			case 'r':
				result += "\r"
			case '\\':
				result += "\\"
			case '"':
				result += "\""
			case '\'':
				result += "'"
			case '0':
				result += "\x00"
			default:
				result += string(l.current)
			}
		} else {
			result += string(l.current)
		}
		l.readChar()
	}

	if l.current == quote {
		l.readChar() // skip closing quote
	}

	return result
}

// readCharLiteral reads a character literal
func (l *Lexer) readCharLiteral() string {
	l.readChar() // skip opening quote

	var result string
	if l.current == '\\' {
		l.readChar()
		// Handle escape sequences
		switch l.current {
		case 'n':
			result = "\n"
		case 't':
			result = "\t"
		case 'r':
			result = "\r"
		case '\\':
			result = "\\"
		case '\'':
			result = "'"
		case '0':
			result = "\x00"
		default:
			result = string(l.current)
		}
		l.readChar()
	} else {
		result = string(l.current)
		l.readChar()
	}

	if l.current == '\'' {
		l.readChar() // skip closing quote
	}

	return result
}

// readRegex reads a regex literal (/pattern/flags)
func (l *Lexer) readRegex() string {
	l.readChar() // skip opening /

	var pattern string
	for l.current != '/' && l.current != '\x00' {
		if l.current == '\\' {
			pattern += string(l.current)
			l.readChar()
			if l.position < len(l.input) {
				pattern += string(l.current)
				l.readChar()
			}
		} else {
			pattern += string(l.current)
			l.readChar()
		}
	}

	if l.current == '/' {
		l.readChar() // skip closing /
	}

	// Read flags (g, i, m, s, u, y, etc.)
	var flags string
	for l.isLetter(l.current) {
		flags += string(l.current)
		l.readChar()
	}

	if flags != "" {
		return pattern + "/" + flags
	}
	return pattern
}

// isLetter checks if a character is a letter or underscore
func (l *Lexer) isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isDigit checks if a character is a digit
func (l *Lexer) isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isIdentifierChar checks if a character can be part of an identifier
func (l *Lexer) isIdentifierChar(ch rune) bool {
	return l.isLetter(ch) || l.isDigit(ch)
}

// makeToken creates a token and tracks the last token type
func (l *Lexer) makeToken(tokenType TokenType, value string) Token {
	l.lastTokenType = tokenType
	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	}
}

// isRegexContext checks if a regex literal can appear in current context
// Regex can follow: =, (, [, {, return, throw, :, ,, ;, &&, ||, !, ?, etc.
func (l *Lexer) isRegexContext() bool {
	switch l.lastTokenType {
	// After assignment/operators
	case ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN, STAR_ASSIGN, SLASH_ASSIGN, PERCENT_ASSIGN:
		// After delimiters
		fallthrough
	case LPAREN, LBRACKET, LBRACE, COMMA, SEMICOLON, COLON:
		// After keywords
		fallthrough
	case RETURN, THROW, IF, WHILE, FOR:
		// After logical operators
		fallthrough
	case AND, OR, NOT:
		// After comparison operators
		fallthrough
	case EQ, NE, LT, GT, LE, GE:
		// At start
		fallthrough
	case EOF, LET, CONST:
		return true
	default:
		return false
	}
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	// Handle newline
	if l.current == '\n' {
		token := l.makeToken(NEWLINE, "\\n")
		l.readChar()
		l.line++
		l.column = 1
		return token
	}

	// Handle comments
	if l.current == '/' && l.peekChar() == '/' {
		l.skipComment()
		return l.NextToken()
	}

	if l.current == '/' && l.peekChar() == '*' {
		l.skipMultiLineComment()
		return l.NextToken()
	}

	// Handle regex literals: /pattern/flags
	if l.current == '/' && l.isRegexContext() {
		value := l.readRegex()
		token := l.makeToken(REGEX, value)
		l.lastTokenType = REGEX
		return token
	}

	// EOF
	if l.current == '\x00' {
		return l.makeToken(EOF, "")
	}

	// String (double quotes)
	if l.current == '"' {
		value := l.readString()
		return l.makeToken(STRING, value)
	}

	// Char literal (single quotes)
	if l.current == '\'' {
		value := l.readCharLiteral()
		return l.makeToken(CHAR, value)
	}

	// Number
	if l.isDigit(l.current) {
		value := l.readNumber()
		return l.makeToken(NUMBER, value)
	}

	// Identifier or Keyword
	if l.isLetter(l.current) {
		value := l.readIdentifier()
		tokenType := GetKeywordType(value)
		return l.makeToken(tokenType, value)
	}

	// Two-character operators
	twoChar := string(l.current) + string(l.peekChar())
	switch twoChar {
	case "==":
		l.readChar()
		l.readChar()
		return l.makeToken(EQ, "==")
	case "!=":
		l.readChar()
		l.readChar()
		return l.makeToken(NE, "!=")
	case "<=":
		l.readChar()
		l.readChar()
		return l.makeToken(LE, "<=")
	case ">=":
		l.readChar()
		l.readChar()
		return l.makeToken(GE, ">=")
	case "&&":
		l.readChar()
		l.readChar()
		return l.makeToken(AND, "&&")
	case "||":
		l.readChar()
		l.readChar()
		return l.makeToken(OR, "||")
	case "<<":
		l.readChar()
		l.readChar()
		return l.makeToken(SHL, "<<")
	case ">>":
		l.readChar()
		l.readChar()
		return l.makeToken(SHR, ">>")
	case "+=":
		l.readChar()
		l.readChar()
		return l.makeToken(PLUS_ASSIGN, "+=")
	case "-=":
		l.readChar()
		l.readChar()
		return l.makeToken(MINUS_ASSIGN, "-=")
	case "*=":
		l.readChar()
		l.readChar()
		return l.makeToken(STAR_ASSIGN, "*=")
	case "/=":
		l.readChar()
		l.readChar()
		return l.makeToken(SLASH_ASSIGN, "/=")
	case "%=":
		l.readChar()
		l.readChar()
		return l.makeToken(PERCENT_ASSIGN, "%=")
	case "->":
		l.readChar()
		l.readChar()
		return l.makeToken(ARROW, "->")
	case "=>":
		l.readChar()
		l.readChar()
		return l.makeToken(FAT_ARROW, "=>")
	case "::":
		l.readChar()
		l.readChar()
		return l.makeToken(COLON_COLON, "::")
	case "**":
		l.readChar()
		l.readChar()
		return l.makeToken(POWER, "**")
	case "..":
		l.readChar()
		l.readChar()
		if l.current == '=' {
			l.readChar()
			return l.makeToken(RANGE_INC, "..=")
		}
		return l.makeToken(RANGE, "..")
	case "|>":
		l.readChar()
		l.readChar()
		return l.makeToken(PIPE_GT, "|>")
	}

	// Single-character tokens
	ch := l.current
	l.readChar()

	switch ch {
	case '+':
		return l.makeToken(PLUS, "+")
	case '-':
		return l.makeToken(MINUS, "-")
	case '*':
		return l.makeToken(STAR, "*")
	case '/':
		return l.makeToken(SLASH, "/")
	case '%':
		return l.makeToken(PERCENT, "%")
	case '<':
		return l.makeToken(LT, "<")
	case '>':
		return l.makeToken(GT, ">")
	case '!':
		return l.makeToken(NOT, "!")
	case '&':
		return l.makeToken(BIT_AND, "&")
	case '|':
		return l.makeToken(BIT_OR, "|")
	case '^':
		return l.makeToken(BIT_XOR, "^")
	case '~':
		return l.makeToken(BIT_NOT, "~")
	case '=':
		return l.makeToken(ASSIGN, "=")
	case '.':
		return l.makeToken(DOT, ".")
	case '?':
		return l.makeToken(QUESTION, "?")
	case '(':
		return l.makeToken(LPAREN, "(")
	case ')':
		return l.makeToken(RPAREN, ")")
	case '[':
		return l.makeToken(LBRACKET, "[")
	case ']':
		return l.makeToken(RBRACKET, "]")
	case '{':
		return l.makeToken(LBRACE, "{")
	case '}':
		return l.makeToken(RBRACE, "}")
	case ',':
		return l.makeToken(COMMA, ",")
	case ';':
		return l.makeToken(SEMICOLON, ";")
	case ':':
		return l.makeToken(COLON, ":")
	case '@':
		return l.makeToken(AT, "@")
	case '#':
		return l.makeToken(HASH, "#")
	default:
		return l.makeToken(ILLEGAL, string(ch))
	}
}

// Tokenize returns all tokens (excluding newlines)
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		token := l.NextToken()
		if token.Type != NEWLINE {
			tokens = append(tokens, token)
		}
		if token.Type == EOF {
			break
		}
	}
	return tokens
}

// TokenizeWithNewlines returns all tokens including newlines
func (l *Lexer) TokenizeWithNewlines() []Token {
	var tokens []Token
	for {
		token := l.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}
	return tokens
}

// DebugPrint prints all tokens for debugging
func (l *Lexer) DebugPrint() {
	tokens := l.Tokenize()
	for i, token := range tokens {
		fmt.Printf("[%d] %s\n", i, token.String())
	}
}
