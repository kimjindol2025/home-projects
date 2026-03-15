package compiler

import (
	"fmt"

	"github.com/freelang-ai/gofree/internal/ast"
)

// Compiler converts AST to IR (Intermediate Representation)
type Compiler struct {
	program    *IRProgram
	currentFn  *IRFunction
	scopes     []*CompileScope
	breakLabel int
	continueLabel int
	labelCounter int
	errors     []CompileError
}

// CompileError represents a compilation error
type CompileError struct {
	Line    int
	Column  int
	Message string
}

// CompileScope represents a scope during compilation
type CompileScope struct {
	variables map[string]int // Variable name to index
	parent    *CompileScope
}

// NewCompiler creates a new compiler
func NewCompiler() *Compiler {
	return &Compiler{
		program:    NewIRProgram(),
		scopes:     []*CompileScope{NewCompileScope(nil)},
		errors:     []CompileError{},
		labelCounter: 0,
	}
}

// NewCompileScope creates a new compile scope
func NewCompileScope(parent *CompileScope) *CompileScope {
	return &CompileScope{
		variables: make(map[string]int),
		parent:    parent,
	}
}

// Compile compiles an AST module to IR
func (c *Compiler) Compile(module *ast.Module) (*IRProgram, error) {
	if module == nil {
		return nil, fmt.Errorf("module is nil")
	}

	// Compile all statements
	for _, stmt := range module.Statements {
		c.compileStatement(stmt)
	}

	if len(c.errors) > 0 {
		return nil, fmt.Errorf("%d compilation errors", len(c.errors))
	}

	return c.program, nil
}

// compileStatement compiles a single statement
func (c *Compiler) compileStatement(stmt ast.Statement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		c.compileVariableDeclaration(s)
	case *ast.FunctionStatement:
		c.compileFunctionStatement(s)
	case *ast.IfStatement:
		c.compileIfStatement(s)
	case *ast.ForStatement:
		c.compileForStatement(s)
	case *ast.ForOfStatement:
		c.compileForOfStatement(s)
	case *ast.WhileStatement:
		c.compileWhileStatement(s)
	case *ast.ReturnStatement:
		c.compileReturnStatement(s)
	case *ast.ExpressionStatement:
		c.compileExpressionStatement(s)
	case *ast.BlockStatement:
		c.compileBlockStatement(s)
	case *ast.ThrowStatement:
		c.compileThrowStatement(s)
	case *ast.BreakStatement:
		c.emit(OpJump, c.breakLabel)
	case *ast.ContinueStatement:
		c.emit(OpJump, c.continueLabel)
	case *ast.TryStatement:
		c.compileTryStatement(s)
	default:
		// Unknown statement type
	}
}

// compileVariableDeclaration compiles variable declaration
func (c *Compiler) compileVariableDeclaration(decl *ast.VariableDeclaration) {
	// Compile the value expression
	if decl.Value != nil {
		c.compileExpression(decl.Value)
	} else {
		// Default value (null/0/false depending on type)
		c.emit(OpLoadConst, nil)
	}

	// Store in variable
	scope := c.currentScope()
	if _, exists := scope.variables[decl.Name]; exists {
		c.addError(0, 0, fmt.Sprintf("variable %q already defined", decl.Name))
	}

	// Assign index and store
	index := len(scope.variables)
	scope.variables[decl.Name] = index
	c.emit(OpStoreVar, index)
}

// compileFunctionStatement compiles function definition
func (c *Compiler) compileFunctionStatement(fn *ast.FunctionStatement) {
	// Create IR function
	irFn := NewIRFunction(fn.Name)
	irFn.ParamCount = len(fn.Parameters)

	// Register function
	c.program.Functions[fn.Name] = irFn

	// Push new scope for function
	c.scopes = append(c.scopes, NewCompileScope(c.currentScope()))

	// Define parameters
	scope := c.currentScope()
	for i, param := range fn.Parameters {
		scope.variables[param.Name] = i
	}

	// Set current function for compilation
	prevFn := c.currentFn
	c.currentFn = irFn

	// Compile function body
	if fn.Body != nil {
		c.compileBlockStatement(fn.Body)
	}

	// Add implicit return
	c.currentFn.Instructions = append(c.currentFn.Instructions, NewInstruction(OpLoadConst, nil))
	c.currentFn.Instructions = append(c.currentFn.Instructions, NewInstruction(OpReturn))

	// Restore previous function
	c.currentFn = prevFn
	c.scopes = c.scopes[:len(c.scopes)-1]
}

// compileIfStatement compiles if statement
func (c *Compiler) compileIfStatement(stmt *ast.IfStatement) {
	// Compile condition
	c.compileExpression(stmt.Condition)

	// Jump to else if false
	elseLabel := c.newLabel()
	c.emit(OpJumpIfFalse, elseLabel)

	// Compile then branch
	if stmt.ThenBranch != nil {
		c.compileStatement(stmt.ThenBranch)
	}

	if stmt.ElseBranch != nil {
		endLabel := c.newLabel()
		c.emit(OpJump, endLabel)
		c.patchLabel(elseLabel)
		c.compileStatement(stmt.ElseBranch)
		c.patchLabel(endLabel)
	} else {
		c.patchLabel(elseLabel)
	}
}

// compileForStatement compiles for loop
func (c *Compiler) compileForStatement(stmt *ast.ForStatement) {
	// Create new scope for loop
	c.pushScope()

	// Compile init
	if stmt.Init != nil {
		c.compileStatement(stmt.Init)
	}

	var loopStart int
	if c.currentFn != nil {
		loopStart = len(c.currentFn.Instructions)
	} else {
		loopStart = len(c.program.Instructions)
	}
	breakLabel := c.newLabel()
	continueLabel := c.newLabel()

	prevBreak := c.breakLabel
	prevContinue := c.continueLabel
	c.breakLabel = breakLabel
	c.continueLabel = continueLabel

	// Compile condition
	if stmt.Condition != nil {
		c.compileExpression(stmt.Condition)
		c.emit(OpJumpIfFalse, breakLabel)
	}

	// Compile body
	if stmt.Body != nil {
		c.compileStatement(stmt.Body)
	}

	// Patch continue label
	c.patchLabel(continueLabel)

	// Compile update
	if stmt.Update != nil {
		c.compileExpression(stmt.Update)
	}

	// Jump back to start
	c.emit(OpJump, loopStart)

	// Patch break label
	c.patchLabel(breakLabel)

	c.breakLabel = prevBreak
	c.continueLabel = prevContinue
	c.popScope()
}

// compileForOfStatement compiles for...of loop
func (c *Compiler) compileForOfStatement(stmt *ast.ForOfStatement) {
	// Compile iterable
	c.compileExpression(stmt.Iterable)

	// Create new scope
	c.pushScope()

	// Store loop variable
	scope := c.currentScope()
	scope.variables[stmt.Variable] = 0

	// Compile body
	if stmt.Body != nil {
		c.compileStatement(stmt.Body)
	}

	c.popScope()
}

// compileWhileStatement compiles while loop
func (c *Compiler) compileWhileStatement(stmt *ast.WhileStatement) {
	c.pushScope()

	var loopStart int
	if c.currentFn != nil {
		loopStart = len(c.currentFn.Instructions)
	} else {
		loopStart = len(c.program.Instructions)
	}
	breakLabel := c.newLabel()
	continueLabel := loopStart

	prevBreak := c.breakLabel
	prevContinue := c.continueLabel
	c.breakLabel = breakLabel
	c.continueLabel = continueLabel

	// Compile condition
	if stmt.Condition != nil {
		c.compileExpression(stmt.Condition)
		c.emit(OpJumpIfFalse, breakLabel)
	}

	// Compile body
	if stmt.Body != nil {
		c.compileStatement(stmt.Body)
	}

	// Jump back
	c.emit(OpJump, loopStart)

	// Patch break label
	c.patchLabel(breakLabel)

	c.breakLabel = prevBreak
	c.continueLabel = prevContinue
	c.popScope()
}

// compileReturnStatement compiles return statement
func (c *Compiler) compileReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.Value != nil {
		c.compileExpression(stmt.Value)
	} else {
		c.emit(OpLoadConst, nil)
	}
	c.emit(OpReturn)
}

// compileExpressionStatement compiles expression statement
func (c *Compiler) compileExpressionStatement(stmt *ast.ExpressionStatement) {
	if stmt.Expression != nil {
		c.compileExpression(stmt.Expression)
		c.emit(OpPop) // Discard result
	}
}

// compileBlockStatement compiles block statement
func (c *Compiler) compileBlockStatement(stmt *ast.BlockStatement) {
	c.pushScope()
	for _, s := range stmt.Statements {
		c.compileStatement(s)
	}
	c.popScope()
}

// compileTryStatement compiles try/catch/finally
func (c *Compiler) compileTryStatement(stmt *ast.TryStatement) {
	c.emit(OpTryEnter)
	if stmt.Try != nil {
		c.compileBlockStatement(stmt.Try)
	}
	c.emit(OpTryExit)
	// TODO: Handle catch and finally blocks
}

// compileThrowStatement compiles throw statement
func (c *Compiler) compileThrowStatement(stmt *ast.ThrowStatement) {
	if stmt.Argument != nil {
		c.compileExpression(stmt.Argument)
	}
	c.emit(OpThrow)
}

// compileExpression compiles an expression
func (c *Compiler) compileExpression(expr ast.Expression) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.LiteralExpression:
		c.compileLiteralExpression(e)
	case *ast.IdentifierExpression:
		c.compileIdentifierExpression(e)
	case *ast.BinaryOpExpression:
		c.compileBinaryOpExpression(e)
	case *ast.CallExpression:
		c.compileCallExpression(e)
	case *ast.ArrayExpression:
		c.compileArrayExpression(e)
	case *ast.MemberExpression:
		c.compileMemberExpression(e)
	case *ast.MatchExpression:
		c.compileMatchExpression(e)
	default:
		// Unknown expression type
	}
}

// compileLiteralExpression compiles literal
func (c *Compiler) compileLiteralExpression(expr *ast.LiteralExpression) {
	c.emit(OpLoadConst, expr.Value)
}

// compileIdentifierExpression compiles identifier
func (c *Compiler) compileIdentifierExpression(expr *ast.IdentifierExpression) {
	scope := c.currentScope()
	if index, exists := scope.resolve(expr.Name); exists {
		c.emit(OpLoadVar, index)
	} else {
		c.addError(0, 0, fmt.Sprintf("undefined variable %q", expr.Name))
	}
}

// compileBinaryOpExpression compiles binary operation
func (c *Compiler) compileBinaryOpExpression(expr *ast.BinaryOpExpression) {
	c.compileExpression(expr.Left)
	c.compileExpression(expr.Right)

	switch expr.Operator {
	case "+":
		c.emit(OpAdd)
	case "-":
		c.emit(OpSub)
	case "*":
		c.emit(OpMul)
	case "/":
		c.emit(OpDiv)
	case "%":
		c.emit(OpMod)
	case "==":
		c.emit(OpEqual)
	case "!=":
		c.emit(OpNotEqual)
	case "<":
		c.emit(OpLess)
	case "<=":
		c.emit(OpLessOrEqual)
	case ">":
		c.emit(OpGreater)
	case ">=":
		c.emit(OpGreaterOrEqual)
	case "&&":
		c.emit(OpAnd)
	case "||":
		c.emit(OpOr)
	default:
		c.addError(0, 0, fmt.Sprintf("unknown operator %q", expr.Operator))
	}
}

// compileCallExpression compiles function call
func (c *Compiler) compileCallExpression(expr *ast.CallExpression) {
	// Compile arguments
	for _, arg := range expr.Arguments {
		c.compileExpression(arg)
	}
	// Emit call
	c.emit(OpCall, expr.Callee, len(expr.Arguments))
}

// compileArrayExpression compiles array literal
func (c *Compiler) compileArrayExpression(expr *ast.ArrayExpression) {
	for _, elem := range expr.Elements {
		c.compileExpression(elem)
	}
	c.emit(OpArrayCreate, len(expr.Elements))
}

// compileMemberExpression compiles member access
func (c *Compiler) compileMemberExpression(expr *ast.MemberExpression) {
	c.compileExpression(expr.Object)
	if expr.IsIndex {
		c.emit(OpArrayAccess, expr.Property)
	} else {
		c.emit(OpObjectAccess, expr.Property)
	}
}

// compileMatchExpression compiles pattern matching
func (c *Compiler) compileMatchExpression(expr *ast.MatchExpression) {
	// TODO: Implement pattern matching
	c.compileExpression(expr.Expr)
}

// Helper methods

func (c *Compiler) emit(opcode Opcode, args ...interface{}) {
	ins := NewInstruction(opcode, args...)
	if c.currentFn != nil {
		c.currentFn.Instructions = append(c.currentFn.Instructions, ins)
	} else {
		c.program.Instructions = append(c.program.Instructions, ins)
	}
}

func (c *Compiler) newLabel() int {
	label := c.labelCounter
	c.labelCounter++
	return label
}

func (c *Compiler) patchLabel(label int) {
	// Find instruction references to label and patch them
	// This is a simplified version
}

func (c *Compiler) currentScope() *CompileScope {
	if len(c.scopes) == 0 {
		return nil
	}
	return c.scopes[len(c.scopes)-1]
}

func (c *Compiler) pushScope() {
	parent := c.currentScope()
	c.scopes = append(c.scopes, NewCompileScope(parent))
}

func (c *Compiler) popScope() {
	if len(c.scopes) > 1 {
		c.scopes = c.scopes[:len(c.scopes)-1]
	}
}

func (s *CompileScope) resolve(name string) (int, bool) {
	if index, ok := s.variables[name]; ok {
		return index, true
	}
	if s.parent != nil {
		return s.parent.resolve(name)
	}
	return -1, false
}

func (c *Compiler) addError(line, column int, message string) {
	c.errors = append(c.errors, CompileError{
		Line:    line,
		Column:  column,
		Message: message,
	})
}

// Errors returns compilation errors
func (c *Compiler) Errors() []CompileError {
	return c.errors
}
