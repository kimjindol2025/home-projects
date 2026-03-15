package analyzer

import (
	"fmt"

	"github.com/freelang-ai/gofree/internal/ast"
)

// Analyzer performs semantic analysis on an AST
type Analyzer struct {
	scopes      []*Scope          // Stack of scopes (global, function, block)
	currentFunc *ast.FunctionStatement
	errors      []AnalysisError
	warnings    []AnalysisWarning
	typeCache   map[ast.Expression]string // Type inference cache
}

// AnalysisError represents a semantic error
type AnalysisError struct {
	Line    int
	Column  int
	Message string
}

// AnalysisWarning represents a semantic warning
type AnalysisWarning struct {
	Line    int
	Column  int
	Message string
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		scopes:    []*Scope{NewScope(nil, "global")},
		errors:    []AnalysisError{},
		warnings:  []AnalysisWarning{},
		typeCache: make(map[ast.Expression]string),
	}
}

// Analyze performs analysis on a module
func (a *Analyzer) Analyze(module *ast.Module) error {
	if module == nil {
		return fmt.Errorf("module is nil")
	}

	// Process imports (future)
	// Process exports (future)

	// Analyze statements
	for _, stmt := range module.Statements {
		a.analyzeStatement(stmt)
	}

	// Report errors
	if len(a.errors) > 0 {
		return fmt.Errorf("%d semantic errors found", len(a.errors))
	}

	return nil
}

// analyzeStatement analyzes a single statement
func (a *Analyzer) analyzeStatement(stmt ast.Statement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		a.analyzeVariableDeclaration(s)
	case *ast.FunctionStatement:
		a.analyzeFunctionStatement(s)
	case *ast.IfStatement:
		a.analyzeIfStatement(s)
	case *ast.ForStatement:
		a.analyzeForStatement(s)
	case *ast.ForOfStatement:
		a.analyzeForOfStatement(s)
	case *ast.WhileStatement:
		a.analyzeWhileStatement(s)
	case *ast.ReturnStatement:
		a.analyzeReturnStatement(s)
	case *ast.ExpressionStatement:
		a.analyzeExpressionStatement(s)
	case *ast.BlockStatement:
		a.analyzeBlockStatement(s)
	case *ast.TryStatement:
		a.analyzeTryStatement(s)
	case *ast.ThrowStatement:
		a.analyzeThrowStatement(s)
	case *ast.StructDeclaration:
		a.analyzeStructDeclaration(s)
	case *ast.EnumDeclaration:
		a.analyzeEnumDeclaration(s)
	case *ast.ImportStatement:
		a.analyzeImportStatement(s)
	case *ast.ExportStatement:
		a.analyzeExportStatement(s)
	case *ast.TestBlock:
		a.analyzeTestBlock(s)
	default:
		// Unknown statement type
	}
}

// analyzeVariableDeclaration analyzes variable declarations
func (a *Analyzer) analyzeVariableDeclaration(decl *ast.VariableDeclaration) {
	// Analyze the value expression first
	if decl.Value != nil {
		valueType := a.inferType(decl.Value)
		if decl.Type != "" && decl.Type != valueType && valueType != "" {
			a.addError(0, 0, fmt.Sprintf("type mismatch: expected %s, got %s", decl.Type, valueType))
		}
	}

	// Define variable in current scope
	scope := a.currentScope()
	if scope.isDefined(decl.Name) {
		a.addWarning(0, 0, fmt.Sprintf("variable %q is already defined", decl.Name))
	}
	scope.define(decl.Name, &Symbol{
		Name:     decl.Name,
		Type:     decl.Type,
		IsMutable: decl.IsMutable,
		IsUsed:   false,
		IsGlobal: len(a.scopes) == 1,
	})
}

// analyzeFunctionStatement analyzes function definitions
func (a *Analyzer) analyzeFunctionStatement(fn *ast.FunctionStatement) {
	// Define function in current scope
	scope := a.currentScope()
	scope.define(fn.Name, &Symbol{
		Name:      fn.Name,
		Type:      "function",
		IsUsed:    false,
		IsGlobal:  len(a.scopes) == 1,
		IsMutable: false,
	})

	// Create new scope for function body
	prevFunc := a.currentFunc
	a.currentFunc = fn
	a.pushScope("function")

	// Define parameters
	for _, param := range fn.Parameters {
		a.currentScope().define(param.Name, &Symbol{
			Name:      param.Name,
			Type:      param.Type,
			IsUsed:    false,
			IsGlobal:  false,
			IsMutable: true,
		})
	}

	// Analyze function body
	if fn.Body != nil {
		a.analyzeBlockStatement(fn.Body)
	}

	a.popScope()
	a.currentFunc = prevFunc
}

// analyzeIfStatement analyzes if statements
func (a *Analyzer) analyzeIfStatement(stmt *ast.IfStatement) {
	if stmt.Condition != nil {
		a.inferType(stmt.Condition)
	}
	if stmt.ThenBranch != nil {
		a.analyzeStatement(stmt.ThenBranch)
	}
	if stmt.ElseBranch != nil {
		a.analyzeStatement(stmt.ElseBranch)
	}
}

// analyzeForStatement analyzes for loops
func (a *Analyzer) analyzeForStatement(stmt *ast.ForStatement) {
	a.pushScope("block")

	if stmt.Init != nil {
		a.analyzeStatement(stmt.Init)
	}
	if stmt.Condition != nil {
		a.inferType(stmt.Condition)
	}
	if stmt.Update != nil {
		a.inferType(stmt.Update)
	}
	if stmt.Body != nil {
		a.analyzeStatement(stmt.Body)
	}

	a.popScope()
}

// analyzeForOfStatement analyzes for...of loops
func (a *Analyzer) analyzeForOfStatement(stmt *ast.ForOfStatement) {
	a.pushScope("block")

	// Define loop variable
	a.currentScope().define(stmt.Variable, &Symbol{
		Name:       stmt.Variable,
		Type:       "",
		IsUsed:     false,
		IsGlobal:   false,
		IsMutable:  true,
	})

	if stmt.Iterable != nil {
		a.inferType(stmt.Iterable)
	}
	if stmt.Body != nil {
		a.analyzeStatement(stmt.Body)
	}

	a.popScope()
}

// analyzeWhileStatement analyzes while loops
func (a *Analyzer) analyzeWhileStatement(stmt *ast.WhileStatement) {
	a.pushScope("block")

	if stmt.Condition != nil {
		a.inferType(stmt.Condition)
	}
	if stmt.Body != nil {
		a.analyzeStatement(stmt.Body)
	}

	a.popScope()
}

// analyzeReturnStatement analyzes return statements
func (a *Analyzer) analyzeReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.Value != nil {
		a.inferType(stmt.Value)
	}
	if a.currentFunc == nil {
		a.addError(0, 0, "return statement outside of function")
	}
}

// analyzeExpressionStatement analyzes expression statements
func (a *Analyzer) analyzeExpressionStatement(stmt *ast.ExpressionStatement) {
	if stmt.Expression != nil {
		a.inferType(stmt.Expression)
	}
}

// analyzeBlockStatement analyzes block statements
func (a *Analyzer) analyzeBlockStatement(stmt *ast.BlockStatement) {
	a.pushScope("block")
	for _, s := range stmt.Statements {
		a.analyzeStatement(s)
	}
	a.popScope()
}

// analyzeTryStatement analyzes try/catch/finally blocks
func (a *Analyzer) analyzeTryStatement(stmt *ast.TryStatement) {
	if stmt.Try != nil {
		a.analyzeBlockStatement(stmt.Try)
	}
	for _, catch := range stmt.Catches {
		a.pushScope("block")
		a.currentScope().define(catch.Parameter, &Symbol{
			Name:      catch.Parameter,
			Type:      catch.Type,
			IsUsed:    false,
			IsGlobal:  false,
			IsMutable: false,
		})
		if catch.Body != nil {
			a.analyzeBlockStatement(catch.Body)
		}
		a.popScope()
	}
	if stmt.Finally != nil {
		a.analyzeBlockStatement(stmt.Finally)
	}
}

// analyzeThrowStatement analyzes throw statements
func (a *Analyzer) analyzeThrowStatement(stmt *ast.ThrowStatement) {
	if stmt.Argument != nil {
		a.inferType(stmt.Argument)
	}
}

// analyzeStructDeclaration analyzes struct definitions
func (a *Analyzer) analyzeStructDeclaration(decl *ast.StructDeclaration) {
	scope := a.currentScope()
	scope.define(decl.Name, &Symbol{
		Name:      decl.Name,
		Type:      "struct",
		IsUsed:    false,
		IsGlobal:  len(a.scopes) == 1,
		IsMutable: false,
	})
}

// analyzeEnumDeclaration analyzes enum definitions
func (a *Analyzer) analyzeEnumDeclaration(decl *ast.EnumDeclaration) {
	scope := a.currentScope()
	scope.define(decl.Name, &Symbol{
		Name:      decl.Name,
		Type:      "enum",
		IsUsed:    false,
		IsGlobal:  len(a.scopes) == 1,
		IsMutable: false,
	})
}

// analyzeImportStatement analyzes import statements
func (a *Analyzer) analyzeImportStatement(stmt *ast.ImportStatement) {
	// Import analysis (future)
}

// analyzeExportStatement analyzes export statements
func (a *Analyzer) analyzeExportStatement(stmt *ast.ExportStatement) {
	// Export analysis (future)
}

// analyzeTestBlock analyzes test blocks
func (a *Analyzer) analyzeTestBlock(block *ast.TestBlock) {
	a.pushScope("block")
	for _, assertion := range block.Assertions {
		a.inferType(assertion.Condition)
	}
	a.popScope()
}

// inferType infers the type of an expression
func (a *Analyzer) inferType(expr ast.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.LiteralExpression:
		return e.DataType
	case *ast.IdentifierExpression:
		return a.resolveIdentifierType(e.Name)
	case *ast.BinaryOpExpression:
		leftType := a.inferType(e.Left)
		rightType := a.inferType(e.Right)
		return a.infeBinaryOpType(e.Operator, leftType, rightType)
	case *ast.CallExpression:
		return a.inferCallExpressionType(e)
	case *ast.ArrayExpression:
		return "array"
	case *ast.LambdaExpression:
		return "function"
	case *ast.AwaitExpression:
		return a.inferType(e.Argument)
	case *ast.MatchExpression:
		return a.inferMatchExpressionType(e)
	default:
		return ""
	}
}

// resolveIdentifierType resolves the type of an identifier
func (a *Analyzer) resolveIdentifierType(name string) string {
	// Search from innermost to outermost scope
	for i := len(a.scopes) - 1; i >= 0; i-- {
		if symbol, ok := a.scopes[i].symbols[name]; ok {
			symbol.IsUsed = true
			return symbol.Type
		}
	}
	a.addError(0, 0, fmt.Sprintf("undefined variable %q", name))
	return ""
}

// infeBinaryOpType infers the result type of a binary operation
func (a *Analyzer) infeBinaryOpType(op, leftType, rightType string) string {
	switch op {
	case "+", "-", "*", "/", "%":
		if leftType == rightType && (leftType == "i64" || leftType == "f64") {
			return leftType
		}
		if leftType == "string" && op == "+" {
			return "string"
		}
		return ""
	case "==", "!=", "<", ">", "<=", ">=":
		return "bool"
	case "&&", "||":
		return "bool"
	default:
		return ""
	}
}

// inferCallExpressionType infers the return type of a function call
func (a *Analyzer) inferCallExpressionType(expr *ast.CallExpression) string {
	// Look up function in scopes
	for i := len(a.scopes) - 1; i >= 0; i-- {
		if symbol, ok := a.scopes[i].symbols[expr.Callee]; ok {
			symbol.IsUsed = true
			if symbol.Type == "function" {
				return "unknown" // TODO: track function return types
			}
		}
	}
	return ""
}

// inferMatchExpressionType infers the type of a match expression
func (a *Analyzer) inferMatchExpressionType(expr *ast.MatchExpression) string {
	// Infer type from match arms (future)
	if len(expr.Arms) > 0 && expr.Arms[0] != nil && expr.Arms[0].Body != nil {
		return a.inferType(expr.Arms[0].Body)
	}
	return ""
}

// Helper methods for scope management

func (a *Analyzer) currentScope() *Scope {
	if len(a.scopes) == 0 {
		return nil
	}
	return a.scopes[len(a.scopes)-1]
}

func (a *Analyzer) pushScope(name string) {
	parent := a.currentScope()
	a.scopes = append(a.scopes, NewScope(parent, name))
}

func (a *Analyzer) popScope() {
	if len(a.scopes) > 1 {
		a.scopes = a.scopes[:len(a.scopes)-1]
	}
}

// Helper methods for error/warning reporting

func (a *Analyzer) addError(line, column int, message string) {
	a.errors = append(a.errors, AnalysisError{
		Line:    line,
		Column:  column,
		Message: message,
	})
}

func (a *Analyzer) addWarning(line, column int, message string) {
	a.warnings = append(a.warnings, AnalysisWarning{
		Line:    line,
		Column:  column,
		Message: message,
	})
}

// Getters for errors and warnings

func (a *Analyzer) Errors() []AnalysisError {
	return a.errors
}

func (a *Analyzer) Warnings() []AnalysisWarning {
	return a.warnings
}
