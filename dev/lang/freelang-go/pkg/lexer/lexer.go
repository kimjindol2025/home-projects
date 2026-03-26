package lexer

import (
	"github.com/kimjindol2025/freelang-go/pkg/token"
)

// Lexer tokenizes FreeLang source code
type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position to read
	ch           byte // current char
	line         int
	column       int
}

// New creates a new lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at next char without consuming
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case 0:
		tok = l.newToken(token.EOF, "")
	case '(':
		tok = l.newToken(token.LPAREN, "(")
		l.readChar()
	case ')':
		tok = l.newToken(token.RPAREN, ")")
		l.readChar()
	case '{':
		tok = l.newToken(token.LBRACE, "{")
		l.readChar()
	case '}':
		tok = l.newToken(token.RBRACE, "}")
		l.readChar()
	case '[':
		tok = l.newToken(token.LBRACKET, "[")
		l.readChar()
	case ']':
		tok = l.newToken(token.RBRACKET, "]")
		l.readChar()
	case ',':
		tok = l.newToken(token.COMMA, ",")
		l.readChar()
	case ';':
		tok = l.newToken(token.SEMICOLON, ";")
		l.readChar()
	case ':':
		tok = l.newToken(token.COLON, ":")
		l.readChar()
	case '.':
		if l.peekChar() == '.' && l.position+2 < len(l.input) && l.input[l.position+2] == '.' {
			tok = l.newToken(token.SPREAD, "...")
			l.readChar()
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.DOT, ".")
			l.readChar()
		}
	case '?':
		if l.peekChar() == '.' {
			tok = l.newToken(token.OPTIONAL, "?.")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.QUESTION, "?")
			l.readChar()
		}
	case '+':
		if l.peekChar() == '+' {
			tok = l.newToken(token.INCREMENT, "++")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '=' {
			tok = l.newToken(token.PLUS_ASSIGN, "+=")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.PLUS, "+")
			l.readChar()
		}
	case '-':
		if l.peekChar() == '-' {
			tok = l.newToken(token.DECREMENT, "--")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '=' {
			tok = l.newToken(token.MINUS_ASSIGN, "-=")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '>' {
			tok = l.newToken(token.LAMBDA, "->")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.MINUS, "-")
			l.readChar()
		}
	case '*':
		if l.peekChar() == '*' {
			tok = l.newToken(token.POWER, "**")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '=' {
			tok = l.newToken(token.MUL_ASSIGN, "*=")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.MULTIPLY, "*")
			l.readChar()
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			tok = l.newToken(token.DIV_ASSIGN, "/=")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.DIVIDE, "/")
			l.readChar()
		}
	case '%':
		if l.peekChar() == '=' {
			tok = l.newToken(token.MOD_ASSIGN, "%=")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.MODULO, "%")
			l.readChar()
		}
	case '=':
		if l.peekChar() == '=' {
			tok = l.newToken(token.EQ, "==")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '>' {
			tok = l.newToken(token.ARROW, "=>")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.ASSIGN, "=")
			l.readChar()
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.newToken(token.NOT_EQ, "!=")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.NOT, "!")
			l.readChar()
		}
	case '<':
		if l.peekChar() == '=' {
			tok = l.newToken(token.LE, "<=")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '<' {
			tok = l.newToken(token.LSHIFT, "<<")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.LT, "<")
			l.readChar()
		}
	case '>':
		if l.peekChar() == '=' {
			tok = l.newToken(token.GE, ">=")
			l.readChar()
			l.readChar()
		} else if l.peekChar() == '>' {
			tok = l.newToken(token.RSHIFT, ">>")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.GT, ">")
			l.readChar()
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.newToken(token.AND, "&&")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.BIT_AND, "&")
			l.readChar()
		}
	case '|':
		if l.peekChar() == '|' {
			tok = l.newToken(token.OR, "||")
			l.readChar()
			l.readChar()
		} else {
			tok = l.newToken(token.BIT_OR, "|")
			l.readChar()
		}
	case '^':
		tok = l.newToken(token.BIT_XOR, "^")
		l.readChar()
	case '~':
		tok = l.newToken(token.BIT_NOT, "~")
		l.readChar()
	case '"', '\'', '`':
		quote := l.ch
		literal := l.readString(quote)
		tok = token.Token{
			Type:    token.STRING,
			Literal: literal,
			Line:    l.line,
			Column:  l.column,
		}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			typ := token.LookupIdent(literal)
			tok = token.Token{
				Type:    typ,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			tok = token.Token{
				Type:    token.INT,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) newToken(typ token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    typ,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '$' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	hasDot := false

	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}

	// Scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.position]
}

func (l *Lexer) readString(quote byte) string {
	l.readChar() // skip opening quote
	var result []byte

	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case byte(quote):
				result = append(result, quote)
			default:
				result = append(result, '\\', l.ch)
			}
			l.readChar()
		} else {
			result = append(result, l.ch)
			l.readChar()
		}
	}

	if l.ch == quote {
		l.readChar() // skip closing quote
	}

	return string(result)
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
