package ast

import "fmt"

// Node is the base interface for all AST nodes
type Node interface {
	node()
}

// Expression is the base interface for expressions
type Expression interface {
	Node
	expr()
}

// Statement is the base interface for statements
type Statement interface {
	Node
	stmt()
}

// ============================================================================
// Expressions
// ============================================================================

// LiteralExpression represents a literal value
type LiteralExpression struct {
	Value    interface{} // string, number, bool
	DataType string      // "number", "string", "bool"
}

func (e *LiteralExpression) node() {}
func (e *LiteralExpression) expr() {}

// IdentifierExpression represents a variable/identifier
type IdentifierExpression struct {
	Name string
}

func (e *IdentifierExpression) node() {}
func (e *IdentifierExpression) expr() {}

// BinaryOpExpression represents binary operations
type BinaryOpExpression struct {
	Operator string      // "+", "-", "*", "/", "%", "==", "!=", ">", "<", ">=", "<=", etc.
	Left     Expression
	Right    Expression
}

func (e *BinaryOpExpression) node() {}
func (e *BinaryOpExpression) expr() {}

// CallExpression represents function calls
type CallExpression struct {
	Callee    string       // function name
	Arguments []Expression
}

func (e *CallExpression) node() {}
func (e *CallExpression) expr() {}

// ArrayExpression represents array literals
type ArrayExpression struct {
	Elements []Expression
}

func (e *ArrayExpression) node() {}
func (e *ArrayExpression) expr() {}

// MemberExpression represents member access (obj.prop or obj[index])
type MemberExpression struct {
	Object   Expression
	Property string
	IsIndex  bool // true if obj[index], false if obj.prop
}

func (e *MemberExpression) node() {}
func (e *MemberExpression) expr() {}

// LambdaExpression represents anonymous functions
type LambdaExpression struct {
	Params     []*Parameter
	ParamTypes []string        // optional type annotations
	Body       interface{}     // Expression or *BlockStatement
	ReturnType string          // optional return type
	CapturedVs []string        // captured variables
}

func (e *LambdaExpression) node() {}
func (e *LambdaExpression) expr() {}

// AwaitExpression represents await in async functions
type AwaitExpression struct {
	Argument Expression
}

func (e *AwaitExpression) node() {}
func (e *AwaitExpression) expr() {}

// MatchExpression represents pattern matching
type MatchExpression struct {
	Expr Expression
	Arms []*MatchArm
}

func (e *MatchExpression) node() {}
func (e *MatchExpression) expr() {}

// MatchArm represents a single arm in a match expression
type MatchArm struct {
	Pattern Pattern
	Guard   Expression // optional
	Body    Expression
}

// ============================================================================
// Patterns
// ============================================================================

// Pattern is the base interface for patterns
type Pattern interface {
	Node
	pattern()
}

// LiteralPattern matches a literal value
type LiteralPattern struct {
	Value interface{}
}

func (p *LiteralPattern) node()    {}
func (p *LiteralPattern) pattern() {}

// VariablePattern binds a variable
type VariablePattern struct {
	Name string
}

func (p *VariablePattern) node()    {}
func (p *VariablePattern) pattern() {}

// WildcardPattern matches anything
type WildcardPattern struct{}

func (p *WildcardPattern) node()    {}
func (p *WildcardPattern) pattern() {}

// ArrayPattern matches array patterns
type ArrayPattern struct {
	Patterns []Pattern
}

func (p *ArrayPattern) node()    {}
func (p *ArrayPattern) pattern() {}

// StructPattern matches struct patterns
type StructPattern struct {
	Fields map[string]Pattern
}

func (p *StructPattern) node()    {}
func (p *StructPattern) pattern() {}

// ============================================================================
// Statements
// ============================================================================

// BlockStatement represents a block of statements
type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) node() {}
func (s *BlockStatement) stmt() {}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Expression Expression
}

func (s *ExpressionStatement) node() {}
func (s *ExpressionStatement) stmt() {}

// VariableDeclaration represents let/const declarations
type VariableDeclaration struct {
	Kind      string       // "let" or "const"
	Name      string
	Type      string       // optional type annotation
	Value     Expression   // optional initial value
	IsMutable bool
}

func (s *VariableDeclaration) node() {}
func (s *VariableDeclaration) stmt() {}

// FunctionStatement represents function declarations
type FunctionStatement struct {
	Name       string
	TypeParams []string     // generic type parameters
	Parameters []*Parameter
	ReturnType string       // optional return type
	Body       *BlockStatement
	IsAsync    bool
}

func (s *FunctionStatement) node() {}
func (s *FunctionStatement) stmt() {}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string // optional type annotation
}

// IfStatement represents if/else statements
type IfStatement struct {
	Condition   Expression
	ThenBranch  Statement
	ElseBranch  Statement // optional
}

func (s *IfStatement) node() {}
func (s *IfStatement) stmt() {}

// ForStatement represents for loops
type ForStatement struct {
	Init      Statement   // optional
	Condition Expression  // optional
	Update    Expression  // optional
	Body      Statement
}

func (s *ForStatement) node() {}
func (s *ForStatement) stmt() {}

// ForOfStatement represents for...of loops
type ForOfStatement struct {
	Variable string     // loop variable
	Iterable Expression
	Body     Statement
}

func (s *ForOfStatement) node() {}
func (s *ForOfStatement) stmt() {}

// WhileStatement represents while loops
type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (s *WhileStatement) node() {}
func (s *WhileStatement) stmt() {}

// BreakStatement represents break
type BreakStatement struct {
	Label string // optional label
}

func (s *BreakStatement) node() {}
func (s *BreakStatement) stmt() {}

// ContinueStatement represents continue
type ContinueStatement struct {
	Label string // optional label
}

func (s *ContinueStatement) node() {}
func (s *ContinueStatement) stmt() {}

// ReturnStatement represents return statements
type ReturnStatement struct {
	Value Expression // optional return value
}

func (s *ReturnStatement) node() {}
func (s *ReturnStatement) stmt() {}

// TryStatement represents try/catch/finally blocks
type TryStatement struct {
	Try     *BlockStatement
	Catches []*CatchClause
	Finally *BlockStatement // optional
}

func (s *TryStatement) node() {}
func (s *TryStatement) stmt() {}

// CatchClause represents a catch block
type CatchClause struct {
	Parameter string          // error variable name
	Type      string          // optional error type
	Body      *BlockStatement
}

// ThrowStatement represents throw statements
type ThrowStatement struct {
	Argument Expression
}

func (s *ThrowStatement) node() {}
func (s *ThrowStatement) stmt() {}

// StructDeclaration represents struct definitions
type StructDeclaration struct {
	Name   string
	Fields map[string]string // field name -> type
}

func (s *StructDeclaration) node() {}
func (s *StructDeclaration) stmt() {}

// EnumDeclaration represents enum definitions
type EnumDeclaration struct {
	Name   string
	Values []string
}

func (s *EnumDeclaration) node() {}
func (s *EnumDeclaration) stmt() {}

// ImportStatement represents import statements
type ImportStatement struct {
	Imports       []*ImportSpecifier
	From          string
	IsNamespace   bool
	Namespace     string
}

func (s *ImportStatement) node() {}
func (s *ImportStatement) stmt() {}

// ImportSpecifier represents a single import
type ImportSpecifier struct {
	Name  string
	Alias string // optional
}

// ExportStatement represents export statements
type ExportStatement struct {
	Declaration Statement // FunctionStatement or VariableDeclaration
}

func (s *ExportStatement) node() {}
func (s *ExportStatement) stmt() {}

// SecretDeclaration represents secret variable declarations
type SecretDeclaration struct {
	Name  string
	Value Expression
}

func (s *SecretDeclaration) node() {}
func (s *SecretDeclaration) stmt() {}

// StyleDeclaration represents style declarations
type StyleDeclaration struct {
	Name       string
	Properties []*StyleProperty
}

func (s *StyleDeclaration) node() {}
func (s *StyleDeclaration) stmt() {}

// StyleProperty represents a style property
type StyleProperty struct {
	Key   string
	Value string
}

// TestBlock represents test blocks
type TestBlock struct {
	Name       string
	Assertions []*AssertStatement
}

func (s *TestBlock) node() {}
func (s *TestBlock) stmt() {}

// AssertStatement represents assertions in tests
type AssertStatement struct {
	Condition Expression
	Message   string // optional error message
}

func (s *AssertStatement) node() {}
func (s *AssertStatement) stmt() {}

// ============================================================================
// Module Level
// ============================================================================

// Module represents a complete source file
type Module struct {
	Path              string
	Imports           []*ImportStatement
	Exports           []*ExportStatement
	Statements        []Statement
	LintConfig        *LintConfig
	AllowOrigins      []string
	CSPPolicy         string
	ValidateSchemas   []map[string]interface{}
	LocalVault        *LocalVaultConfig
}

// MinimalFunctionAST represents a function in .free format
type MinimalFunctionAST struct {
	Decorator  string // optional "minimal"
	FnName     string
	TypeParams []string
	InputType  string
	OutputType string
	Intent     string
	Reason     string
	Body       string
	Source     map[string]int
}

// ============================================================================
// Configuration
// ============================================================================

// LintConfig represents @lint(...) configuration
type LintConfig struct {
	NoUnused       string // "error", "warn", "off"
	ShadowingCheck string // "error", "warn", "off"
	StrictPointers bool
	Line           int
	Column         int
}

// LocalVaultConfig represents @local_vault(...) configuration
type LocalVaultConfig struct {
	Path     string
	Autosave bool
	Line     int
	Column   int
}

// ============================================================================
// Error Handling
// ============================================================================

// ParseError represents a parsing error
type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", e.Line, e.Column, e.Message)
}

// NewParseError creates a new parse error
func NewParseError(line, column int, message string) *ParseError {
	return &ParseError{
		Line:    line,
		Column:  column,
		Message: message,
	}
}
