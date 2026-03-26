package parser

import (
	"fmt"
	"strconv"

	"github.com/kimjindol2025/freelang-go/pkg/ast"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
	"github.com/kimjindol2025/freelang-go/pkg/token"
)

// Precedence constants for operator precedence parsing
const (
	LOWEST = iota
	ASSIGN
	OR
	AND
	EQUALS
	COMPARE
	TERM
	FACTOR
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN,
	token.PLUS_ASSIGN:  ASSIGN,
	token.MINUS_ASSIGN: ASSIGN,
	token.MUL_ASSIGN:   ASSIGN,
	token.DIV_ASSIGN:   ASSIGN,
	token.MOD_ASSIGN:   ASSIGN,
	token.OR:           OR,
	token.AND:          AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           COMPARE,
	token.LE:           COMPARE,
	token.GT:           COMPARE,
	token.GE:           COMPARE,
	token.PLUS:         TERM,
	token.MINUS:        TERM,
	token.MULTIPLY:     FACTOR,
	token.DIVIDE:       FACTOR,
	token.MODULO:       FACTOR,
	token.POWER:        FACTOR,
	token.LSHIFT:       COMPARE,
	token.RSHIFT:       COMPARE,
	token.BIT_AND:      EQUALS,
	token.BIT_OR:       EQUALS,
	token.BIT_XOR:      EQUALS,
	token.DOT:          CALL,
	token.LBRACKET:     INDEX,
	token.LPAREN:       CALL,
}

// Parser tokenizes and parses FreeLang code
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.STRING] = p.parseStringLiteral
	p.prefixParseFns[token.BOOL] = p.parseBooleanLiteral
	p.prefixParseFns[token.LBRACKET] = p.parseArrayLiteral
	p.prefixParseFns[token.LBRACE] = p.parseHashLiteral
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[token.FUNC] = p.parseFunctionLiteral
	p.prefixParseFns[token.NOT] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.PLUS] = p.parsePrefixExpression
	p.prefixParseFns[token.BIT_NOT] = p.parsePrefixExpression
	p.prefixParseFns[token.INCREMENT] = p.parsePrefixExpression
	p.prefixParseFns[token.DECREMENT] = p.parsePrefixExpression

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.MULTIPLY] = p.parseInfixExpression
	p.infixParseFns[token.DIVIDE] = p.parseInfixExpression
	p.infixParseFns[token.MODULO] = p.parseInfixExpression
	p.infixParseFns[token.POWER] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.LE] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.GE] = p.parseInfixExpression
	p.infixParseFns[token.AND] = p.parseInfixExpression
	p.infixParseFns[token.OR] = p.parseInfixExpression
	p.infixParseFns[token.BIT_AND] = p.parseInfixExpression
	p.infixParseFns[token.BIT_OR] = p.parseInfixExpression
	p.infixParseFns[token.BIT_XOR] = p.parseInfixExpression
	p.infixParseFns[token.LSHIFT] = p.parseInfixExpression
	p.infixParseFns[token.RSHIFT] = p.parseInfixExpression
	p.infixParseFns[token.ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.PLUS_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.MINUS_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.MUL_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.DIV_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.MOD_ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.LBRACKET] = p.parseIndexExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression
	p.infixParseFns[token.DOT] = p.parseDotExpression

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// CurToken returns current token (for debugging)
func (p *Parser) CurToken() token.Token {
	return p.curToken
}

// PeekToken returns next token (for debugging)
func (p *Parser) PeekToken() token.Token {
	return p.peekToken
}

// nextToken advances to next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses let x = 5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses return expression;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseIfStatement parses if (condition) { block } else { block }
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseForStatement parses for var in iterable { block }
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{
		Token: p.curToken,
	}

	p.nextToken()

	// Check if there's a loop variable (for x in ...)
	if p.curToken.Type == token.IDENTIFIER && p.peekToken.Type == token.IN {
		stmt.LoopVar = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		p.nextToken() // consume 'in'
		p.nextToken() // move to iterable
	}

	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseBlockStatement parses { statement; statement; }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpressionStatement parses an expression as a statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses expressions with operator precedence
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}

	leftExp := prefixFn()

	for precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infixFn(leftExp)
	}

	return leftExp
}

// parsePrefixExpression parses -5, !true, etc.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

// parseInfixExpression parses a + b, a == b, etc.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

// parseAssignmentExpression parses x = 5, x += 1, etc.
func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	expr := &ast.AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expr.Value = p.parseExpression(LOWEST)

	return expr
}

// parseIndexExpression parses array[index]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expr
}

// parseCallExpression parses func(arg1, arg2)
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: []ast.Expression{},
	}

	expr.Arguments = p.parseExpressionList(token.RPAREN)

	return expr
}

// parseDotExpression parses obj.property or obj.method()
func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	// For now, treat as a simple infix operator
	// Full implementation would handle member access
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: ".",
		Left:     left,
	}

	p.nextToken()
	expr.Right = p.parseExpression(CALL)

	return expr
}

// parseIdentifier parses variable name
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseIntegerLiteral parses 123
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral parses "hello"
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseBooleanLiteral parses true/false
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal == "true",
	}
}

// parseArrayLiteral parses [1, 2, 3]
func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{
		Token:    p.curToken,
		Elements: []ast.Expression{},
	}

	arr.Elements = p.parseExpressionList(token.RBRACKET)

	return arr
}

// parseHashLiteral parses {key: value, key: value}
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: p.curToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	if p.peekToken.Type == token.RBRACE {
		p.nextToken()
		return hash
	}

	for {
		p.nextToken()

		// Handle identifier or string keys
		var key ast.Expression
		switch p.curToken.Type {
		case token.IDENTIFIER:
			key = &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}
		case token.STRING:
			key = &ast.StringLiteral{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}
		case token.INT:
			val, _ := strconv.ParseInt(p.curToken.Literal, 0, 64)
			key = &ast.IntegerLiteral{
				Token: p.curToken,
				Value: val,
			}
		default:
			p.errors = append(p.errors, fmt.Sprintf("expected identifier, string, or number as hash key, got %s", p.curToken.Type))
			return nil
		}

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		hash.Pairs[key] = value

		if p.peekToken.Type != token.COMMA {
			break
		}
		p.nextToken() // consume comma
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

// parseFunctionLiteral parses func(a, b) { return a + b; }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// parseGroupedExpression parses (expression)
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}


// parseFunctionParameters parses parameter list
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseExpressionList parses comma-separated expressions until terminator
func (p *Parser) parseExpressionList(terminator token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekToken.Type == terminator {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(terminator) {
		return nil
	}

	return list
}

// expectPeek checks if next token is of expected type
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type != t {
		p.errors = append(p.errors, fmt.Sprintf("expected %s, got %s instead", t, p.peekToken.Type))
		return false
	}

	p.nextToken()
	return true
}

// curPrecedence returns precedence of current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// peekPrecedence returns precedence of next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}
