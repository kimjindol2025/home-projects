package parser

import (
	"fmt"

	"github.com/freelang-ai/gofree/internal/ast"
	"github.com/freelang-ai/gofree/internal/lexer"
)

// Parser parses tokens into an AST
type Parser struct {
	tokens       []lexer.Token
	position     int // current token index
	current      lexer.Token
	next         lexer.Token

	// Performance optimization: operator precedence cache
	precedenceCache map[string]int
}

// NewParser creates a new parser
func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens:          tokens,
		position:        0,
		precedenceCache: make(map[string]int),
	}
	p.initPrecedenceCache()
	p.advance() // load first token
	p.advanceNext() // preload next token
	return p
}

// initPrecedenceCache initializes operator precedence
func (p *Parser) initPrecedenceCache() {
	// Logical operators
	p.precedenceCache["||"] = 1
	p.precedenceCache["&&"] = 2

	// Bitwise operators
	p.precedenceCache["|"] = 3
	p.precedenceCache["^"] = 4
	p.precedenceCache["&"] = 5

	// Comparison operators
	p.precedenceCache["=="] = 6
	p.precedenceCache["!="] = 6
	p.precedenceCache["<"] = 7
	p.precedenceCache[">"] = 7
	p.precedenceCache["<="] = 7
	p.precedenceCache[">="] = 7

	// Shift operators
	p.precedenceCache["<<"] = 8
	p.precedenceCache[">>"] = 8

	// Additive
	p.precedenceCache["+"] = 9
	p.precedenceCache["-"] = 9

	// Multiplicative
	p.precedenceCache["*"] = 10
	p.precedenceCache["/"] = 10
	p.precedenceCache["%"] = 10
}

// advance moves to the next token
func (p *Parser) advance() {
	if p.position < len(p.tokens) {
		p.current = p.tokens[p.position]
		p.position++
	}
}

// advanceNext loads the next token
func (p *Parser) advanceNext() {
	if p.position < len(p.tokens) {
		p.next = p.tokens[p.position]
	} else {
		p.next = lexer.Token{Type: lexer.EOF}
	}
}

// peek returns the next token without consuming
func (p *Parser) peek() lexer.Token {
	return p.next
}

// expect consumes a token of the expected type
func (p *Parser) expect(tokenType lexer.TokenType) (lexer.Token, error) {
	if p.current.Type != tokenType {
		return lexer.Token{}, ast.NewParseError(
			p.current.Line,
			p.current.Column,
			fmt.Sprintf("expected %s, got %s", tokenType, p.current.Type),
		)
	}
	token := p.current
	p.advance()
	p.advanceNext()
	return token, nil
}

// match checks if current token matches any of the types
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.current.Type == t {
			return true
		}
	}
	return false
}

// consume consumes current token if it matches, returns true if matched
func (p *Parser) consume(tokenType lexer.TokenType) bool {
	if p.current.Type == tokenType {
		p.advance()
		p.advanceNext()
		return true
	}
	return false
}

// Parse parses tokens into a Module AST
func (p *Parser) Parse() (*ast.Module, error) {
	module := &ast.Module{
		Imports:    []*ast.ImportStatement{},
		Exports:    []*ast.ExportStatement{},
		Statements: []ast.Statement{},
	}

	// Parse imports first
	for p.match(lexer.IMPORT) {
		impStmt, err := p.parseImportStatement()
		if err != nil {
			return nil, err
		}
		module.Imports = append(module.Imports, impStmt)
	}

	// Parse statements
	for p.current.Type != lexer.EOF {
		// Skip newlines
		if p.consume(lexer.NEWLINE) {
			continue
		}

		// Parse export
		if p.match(lexer.EXPORT) {
			expStmt, err := p.parseExportStatement()
			if err != nil {
				return nil, err
			}
			module.Exports = append(module.Exports, expStmt)
			continue
		}

		// Parse statement
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			module.Statements = append(module.Statements, stmt)
		}
	}

	return module, nil
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.current.Type {
	case lexer.FN:
		return p.parseFunctionStatement()
	case lexer.LET:
		return p.parseVariableDeclaration(false)
	case lexer.CONST:
		return p.parseVariableDeclaration(true)
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	case lexer.TRY:
		return p.parseTryStatement()
	case lexer.THROW:
		return p.parseThrowStatement()
	case lexer.STRUCT:
		return p.parseStructDeclaration()
	case lexer.ENUM:
		return p.parseEnumDeclaration()
	case lexer.LBRACE:
		return p.parseBlockStatement()
	default:
		// Try expression statement
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		// Consume optional semicolon or newline
		p.consume(lexer.SEMICOLON)
		p.consume(lexer.NEWLINE)
		return &ast.ExpressionStatement{Expression: expr}, nil
	}
}

// parseFunctionStatement parses a function declaration
func (p *Parser) parseFunctionStatement() (*ast.FunctionStatement, error) {
	p.expect(lexer.FN)

	nameToken, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	// Parse type parameters if present
	var typeParams []string
	if p.consume(lexer.LT) {
		// TODO: Parse type parameters
		p.expect(lexer.GT)
	}

	// Parse parameters
	params, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	// Parse return type if present
	returnType := ""
	if p.consume(lexer.ARROW) {
		typeToken, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		returnType = typeToken.Value
	}

	// Parse body
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStatement{
		Name:       nameToken.Value,
		TypeParams: typeParams,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
		IsAsync:    false,
	}, nil
}

// parseVariableDeclaration parses variable declarations
func (p *Parser) parseVariableDeclaration(isConst bool) (*ast.VariableDeclaration, error) {
	kind := "let"
	if isConst {
		kind = "const"
	}

	p.expect(lexer.LET)
	if isConst {
		p.expect(lexer.CONST)
	}

	nameToken, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	// Parse type if present
	typeStr := ""
	if p.consume(lexer.COLON) {
		typeToken, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		typeStr = typeToken.Value
	}

	// Parse value if present
	var value ast.Expression
	if p.consume(lexer.ASSIGN) {
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// Consume optional semicolon or newline
	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)

	return &ast.VariableDeclaration{
		Kind:      kind,
		Name:      nameToken.Value,
		Type:      typeStr,
		Value:     value,
		IsMutable: !isConst,
	}, nil
}

// parseIfStatement parses if/else statements
func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	p.expect(lexer.IF)

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Parse then branch
	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Parse optional else branch
	var elseBranch ast.Statement
	if p.consume(lexer.ELSE) {
		var err error
		elseBranch, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

// parseForStatement parses for/for-of loops
func (p *Parser) parseForStatement() (ast.Statement, error) {
	p.expect(lexer.FOR)

	// Try to parse as for-of
	if p.match(lexer.IDENT) {
		ident := p.current.Value
		p.advance()
		p.advanceNext()

		if p.consume(lexer.OF) {
			// for...of loop
			iterable, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			body, err := p.parseStatement()
			if err != nil {
				return nil, err
			}

			return &ast.ForOfStatement{
				Variable: ident,
				Iterable: iterable,
				Body:     body,
			}, nil
		}

		// Reset if not for-of
		p.position--
		p.current = p.tokens[p.position-1]
	}

	// Parse regular for loop
	var init ast.Statement
	if !p.match(lexer.SEMICOLON) {
		var err error
		init, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	} else {
		p.expect(lexer.SEMICOLON)
	}

	var condition ast.Expression
	if !p.match(lexer.SEMICOLON) {
		var err error
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	p.expect(lexer.SEMICOLON)

	var update ast.Expression
	if !p.match(lexer.LBRACE) {
		var err error
		update, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ForStatement{
		Init:      init,
		Condition: condition,
		Update:    update,
		Body:      body,
	}, nil
}

// parseWhileStatement parses while loops
func (p *Parser) parseWhileStatement() (*ast.WhileStatement, error) {
	p.expect(lexer.WHILE)

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	p.expect(lexer.RETURN)

	var value ast.Expression
	if !p.match(lexer.SEMICOLON, lexer.NEWLINE, lexer.EOF, lexer.RBRACE) {
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)

	return &ast.ReturnStatement{Value: value}, nil
}

// parseBreakStatement parses break statements
func (p *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	p.expect(lexer.BREAK)
	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)
	return &ast.BreakStatement{}, nil
}

// parseContinueStatement parses continue statements
func (p *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	p.expect(lexer.CONTINUE)
	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)
	return &ast.ContinueStatement{}, nil
}

// parseBlockStatement parses block statements
func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	p.expect(lexer.LBRACE)

	block := &ast.BlockStatement{Statements: []ast.Statement{}}

	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		if p.consume(lexer.NEWLINE) {
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	p.expect(lexer.RBRACE)

	return block, nil
}

// parseTryStatement parses try/catch/finally blocks
func (p *Parser) parseTryStatement() (*ast.TryStatement, error) {
	p.expect(lexer.TRY)

	tryBlock, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	stmt := &ast.TryStatement{Try: tryBlock}

	// Parse catch clauses
	for p.match(lexer.CATCH) {
		p.expect(lexer.CATCH)

		// Parse parameter
		paramToken, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}

		// Parse optional type
		typeStr := ""
		if p.consume(lexer.COLON) {
			typeToken, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			typeStr = typeToken.Value
		}

		body, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}

		stmt.Catches = append(stmt.Catches, &ast.CatchClause{
			Parameter: paramToken.Value,
			Type:      typeStr,
			Body:      body,
		})
	}

	// Parse optional finally
	if p.consume(lexer.FINALLY) {
		finallyBlock, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		stmt.Finally = finallyBlock
	}

	return stmt, nil
}

// parseThrowStatement parses throw statements
func (p *Parser) parseThrowStatement() (*ast.ThrowStatement, error) {
	p.expect(lexer.THROW)

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)

	return &ast.ThrowStatement{Argument: expr}, nil
}

// parseStructDeclaration parses struct declarations
func (p *Parser) parseStructDeclaration() (*ast.StructDeclaration, error) {
	p.expect(lexer.STRUCT)

	nameToken, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	p.expect(lexer.LBRACE)

	fields := make(map[string]string)
	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		if p.consume(lexer.NEWLINE) {
			continue
		}

		fieldName, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}

		p.expect(lexer.COLON)

		fieldType, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}

		fields[fieldName.Value] = fieldType.Value

		p.consume(lexer.SEMICOLON)
		p.consume(lexer.NEWLINE)
	}

	p.expect(lexer.RBRACE)

	return &ast.StructDeclaration{
		Name:   nameToken.Value,
		Fields: fields,
	}, nil
}

// parseEnumDeclaration parses enum declarations
func (p *Parser) parseEnumDeclaration() (*ast.EnumDeclaration, error) {
	p.expect(lexer.ENUM)

	nameToken, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}

	p.expect(lexer.LBRACE)

	var values []string
	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		if p.consume(lexer.NEWLINE) {
			continue
		}

		valueToken, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}

		values = append(values, valueToken.Value)

		p.consume(lexer.COMMA)
		p.consume(lexer.SEMICOLON)
		p.consume(lexer.NEWLINE)
	}

	p.expect(lexer.RBRACE)

	return &ast.EnumDeclaration{
		Name:   nameToken.Value,
		Values: values,
	}, nil
}

// parseImportStatement parses import statements
func (p *Parser) parseImportStatement() (*ast.ImportStatement, error) {
	p.expect(lexer.IMPORT)

	// TODO: Parse import specifiers

	p.expect(lexer.FROM)

	// Parse module path
	pathToken, err := p.expect(lexer.STRING)
	if err != nil {
		return nil, err
	}

	p.consume(lexer.SEMICOLON)
	p.consume(lexer.NEWLINE)

	return &ast.ImportStatement{
		From: pathToken.Value,
	}, nil
}

// parseExportStatement parses export statements
func (p *Parser) parseExportStatement() (*ast.ExportStatement, error) {
	p.expect(lexer.EXPORT)

	// Parse declaration (function or variable)
	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ExportStatement{Declaration: stmt}, nil
}

// parseParameterList parses function parameters
func (p *Parser) parseParameterList() ([]*ast.Parameter, error) {
	p.expect(lexer.LPAREN)

	var params []*ast.Parameter

	for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
		nameToken, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}

		// Parse optional type
		typeStr := ""
		if p.consume(lexer.COLON) {
			typeToken, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			typeStr = typeToken.Value
		}

		params = append(params, &ast.Parameter{
			Name: nameToken.Value,
			Type: typeStr,
		})

		if !p.consume(lexer.COMMA) {
			break
		}
	}

	p.expect(lexer.RPAREN)

	return params, nil
}

// parseExpression parses expressions
func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseBinaryOp(0)
}

// parseBinaryOp parses binary operations with precedence
func (p *Parser) parseBinaryOp(minPrec int) (ast.Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		op := p.current.Value
		prec, ok := p.precedenceCache[op]
		if !ok || prec < minPrec {
			break
		}

		if !p.isOperator() {
			break
		}

		p.advance()
		p.advanceNext()

		right, err := p.parseBinaryOp(prec + 1)
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOpExpression{
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}

	return left, nil
}

// parsePrimary parses primary expressions
func (p *Parser) parsePrimary() (ast.Expression, error) {
	switch p.current.Type {
	case lexer.NUMBER:
		token := p.current
		p.advance()
		p.advanceNext()
		return &ast.LiteralExpression{
			Value:    token.Value,
			DataType: "number",
		}, nil

	case lexer.STRING:
		token := p.current
		p.advance()
		p.advanceNext()
		return &ast.LiteralExpression{
			Value:    token.Value,
			DataType: "string",
		}, nil

	case lexer.TRUE:
		p.advance()
		p.advanceNext()
		return &ast.LiteralExpression{
			Value:    true,
			DataType: "bool",
		}, nil

	case lexer.FALSE:
		p.advance()
		p.advanceNext()
		return &ast.LiteralExpression{
			Value:    false,
			DataType: "bool",
		}, nil

	case lexer.IDENT:
		return p.parseIdentifierOrCall()

	case lexer.LPAREN:
		p.advance()
		p.advanceNext()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.expect(lexer.RPAREN)
		return expr, nil

	case lexer.LBRACKET:
		return p.parseArrayLiteral()

	default:
		return nil, ast.NewParseError(
			p.current.Line,
			p.current.Column,
			fmt.Sprintf("unexpected token: %s", p.current.Type),
		)
	}
}

// parseIdentifierOrCall parses identifiers and function calls
func (p *Parser) parseIdentifierOrCall() (ast.Expression, error) {
	nameToken := p.current
	p.advance()
	p.advanceNext()

	// Check for function call
	if p.consume(lexer.LPAREN) {
		var args []ast.Expression

		for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if !p.consume(lexer.COMMA) {
				break
			}
		}

		p.expect(lexer.RPAREN)

		return &ast.CallExpression{
			Callee:    nameToken.Value,
			Arguments: args,
		}, nil
	}

	return &ast.IdentifierExpression{Name: nameToken.Value}, nil
}

// parseArrayLiteral parses array literals
func (p *Parser) parseArrayLiteral() (*ast.ArrayExpression, error) {
	p.expect(lexer.LBRACKET)

	var elements []ast.Expression

	for p.current.Type != lexer.RBRACKET && p.current.Type != lexer.EOF {
		elem, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		if !p.consume(lexer.COMMA) {
			break
		}
	}

	p.expect(lexer.RBRACKET)

	return &ast.ArrayExpression{Elements: elements}, nil
}

// isOperator checks if current token is an operator
func (p *Parser) isOperator() bool {
	_, ok := p.precedenceCache[p.current.Value]
	return ok
}
