package compiler

import (
	"fmt"

	"github.com/kimjindol2025/freelang-go/pkg/ast"
	"github.com/kimjindol2025/freelang-go/pkg/checker"
)

// Compiler generates bytecode from AST
type Compiler struct {
	constants  []interface{}
	bytecode   []byte
	errors     []string
	symtab     *checker.SymbolTable
	scopeLevel int
	locals     map[string]int // local variable offsets
}

// New creates a new compiler
func New(symtab *checker.SymbolTable) *Compiler {
	return &Compiler{
		constants:  []interface{}{},
		bytecode:   []byte{},
		errors:     []string{},
		symtab:     symtab,
		scopeLevel: 0,
		locals:     make(map[string]int),
	}
}

// Errors returns compilation errors
func (c *Compiler) Errors() []string {
	return c.errors
}

// Bytecode returns the compiled bytecode
func (c *Compiler) Bytecode() []byte {
	return c.bytecode
}

// Constants returns the constants pool
func (c *Compiler) Constants() []interface{} {
	return c.constants
}

// Compile compiles a program to bytecode
func (c *Compiler) Compile(program *ast.Program) {
	for _, stmt := range program.Statements {
		c.compileStatement(stmt)
	}
}

// compileStatement compiles a statement
func (c *Compiler) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		c.compileLetStatement(s)
	case *ast.ReturnStatement:
		c.compileReturnStatement(s)
	case *ast.ExpressionStatement:
		c.compileExpressionStatement(s)
	case *ast.IfStatement:
		c.compileIfStatement(s)
	case *ast.ForStatement:
		c.compileForStatement(s)
	case *ast.BlockStatement:
		c.compileBlockStatement(s)
	default:
		c.addError("unknown statement type: %T", stmt)
	}
}

// compileLetStatement compiles a let statement
func (c *Compiler) compileLetStatement(stmt *ast.LetStatement) {
	// Compile the value expression
	c.compileExpression(stmt.Value)

	// Store in variable
	localIdx := len(c.locals)
	c.locals[stmt.Name.Value] = localIdx

	// Emit SET_LOCAL
	c.emitInstruction(OpSetLocal, localIdx)
}

// compileReturnStatement compiles a return statement
func (c *Compiler) compileReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.ReturnValue == nil {
		c.emitInstruction(OpReturn)
	} else {
		c.compileExpression(stmt.ReturnValue)
		c.emitInstruction(OpReturnValue)
	}
}

// compileExpressionStatement compiles an expression statement
func (c *Compiler) compileExpressionStatement(stmt *ast.ExpressionStatement) {
	c.compileExpression(stmt.Expression)
	c.emitInstruction(OpPop)
}

// compileIfStatement compiles an if statement
func (c *Compiler) compileIfStatement(stmt *ast.IfStatement) {
	// Compile condition
	c.compileExpression(stmt.Condition)

	// Emit jump instruction (placeholder)
	jumpPos := len(c.bytecode)
	c.emitInstruction(OpJumpNotTruthy, 0) // Placeholder for jump address

	// Compile consequence
	c.compileBlockStatement(stmt.Consequence)

	// Patch jump address
	c.patchJump(jumpPos)

	// If alternative exists, compile it
	if stmt.Alternative != nil {
		c.compileBlockStatement(stmt.Alternative)
	}
}

// compileForStatement compiles a for statement
func (c *Compiler) compileForStatement(stmt *ast.ForStatement) {
	// Enter new scope for loop
	c.scopeLevel++
	loopLocals := make(map[string]int)

	// Define loop variable if present
	if stmt.LoopVar != nil {
		loopLocals[stmt.LoopVar.Value] = len(loopLocals)
	}

	// Save current locals and switch to loop locals
	savedLocals := c.locals
	c.locals = loopLocals

	// Compile loop body (simplified - would need proper loop implementation)
	c.compileBlockStatement(stmt.Body)

	// Restore locals and exit scope
	c.locals = savedLocals
	c.scopeLevel--
}

// compileBlockStatement compiles a block statement
func (c *Compiler) compileBlockStatement(stmt *ast.BlockStatement) {
	for _, s := range stmt.Statements {
		c.compileStatement(s)
	}
}

// compileExpression compiles an expression
func (c *Compiler) compileExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		c.compileIntegerLiteral(e)
	case *ast.StringLiteral:
		c.compileStringLiteral(e)
	case *ast.BooleanLiteral:
		c.compileBooleanLiteral(e)
	case *ast.Identifier:
		c.compileIdentifier(e)
	case *ast.ArrayLiteral:
		c.compileArrayLiteral(e)
	case *ast.HashLiteral:
		c.compileHashLiteral(e)
	case *ast.PrefixExpression:
		c.compilePrefixExpression(e)
	case *ast.InfixExpression:
		c.compileInfixExpression(e)
	case *ast.CallExpression:
		c.compileCallExpression(e)
	case *ast.IndexExpression:
		c.compileIndexExpression(e)
	case *ast.FunctionLiteral:
		c.compileFunctionLiteral(e)
	case *ast.AssignmentExpression:
		c.compileAssignmentExpression(e)
	default:
		c.addError("unknown expression type: %T", expr)
	}
}

// compileIntegerLiteral compiles an integer literal
func (c *Compiler) compileIntegerLiteral(expr *ast.IntegerLiteral) {
	idx := c.addConstant(expr.Value)
	c.emitInstruction(OpConstant, idx)
}

// compileStringLiteral compiles a string literal
func (c *Compiler) compileStringLiteral(expr *ast.StringLiteral) {
	idx := c.addConstant(expr.Value)
	c.emitInstruction(OpConstant, idx)
}

// compileBooleanLiteral compiles a boolean literal
func (c *Compiler) compileBooleanLiteral(expr *ast.BooleanLiteral) {
	if expr.Value {
		c.emitInstruction(OpTrue)
	} else {
		c.emitInstruction(OpFalse)
	}
}

// compileIdentifier compiles an identifier
func (c *Compiler) compileIdentifier(expr *ast.Identifier) {
	if idx, ok := c.locals[expr.Value]; ok {
		c.emitInstruction(OpGetLocal, idx)
	} else {
		// Try global
		c.emitInstruction(OpGetGlobal, 0) // Simplified
	}
}

// compileArrayLiteral compiles an array literal
func (c *Compiler) compileArrayLiteral(expr *ast.ArrayLiteral) {
	// Compile each element
	for _, elem := range expr.Elements {
		c.compileExpression(elem)
	}

	// Emit array instruction with count
	c.emitInstruction(OpArray, len(expr.Elements))
}

// compileHashLiteral compiles a hash literal
func (c *Compiler) compileHashLiteral(expr *ast.HashLiteral) {
	// Compile each key-value pair
	for k, v := range expr.Pairs {
		// If key is a bare identifier, treat it as a string constant (object property key)
		if ident, ok := k.(*ast.Identifier); ok {
			idx := c.addConstant(ident.Value)
			c.emitInstruction(OpConstant, idx)
		} else {
			c.compileExpression(k)
		}
		c.compileExpression(v)
	}

	// Emit hash instruction with count
	c.emitInstruction(OpHash, len(expr.Pairs))
}

// compilePrefixExpression compiles a prefix expression
func (c *Compiler) compilePrefixExpression(expr *ast.PrefixExpression) {
	c.compileExpression(expr.Right)

	switch expr.Operator {
	case "!":
		c.emitInstruction(OpNot)
	case "-":
		c.emitInstruction(OpNegate)
	case "+":
		// No-op
	case "~":
		c.emitInstruction(OpBitNot)
	default:
		c.addError("unknown prefix operator: %s", expr.Operator)
	}
}

// compileInfixExpression compiles an infix expression
func (c *Compiler) compileInfixExpression(expr *ast.InfixExpression) {
	c.compileExpression(expr.Left)
	c.compileExpression(expr.Right)

	switch expr.Operator {
	case "+":
		c.emitInstruction(OpAdd)
	case "-":
		c.emitInstruction(OpSubtract)
	case "*":
		c.emitInstruction(OpMultiply)
	case "/":
		c.emitInstruction(OpDivide)
	case "%":
		c.emitInstruction(OpModulo)
	case "**":
		c.emitInstruction(OpPower)
	case "==":
		c.emitInstruction(OpEqual)
	case "!=":
		c.emitInstruction(OpNotEqual)
	case "<":
		c.emitInstruction(OpLessThan)
	case "<=":
		c.emitInstruction(OpLessThanOrEqual)
	case ">":
		c.emitInstruction(OpGreaterThan)
	case ">=":
		c.emitInstruction(OpGreaterThanOrEqual)
	case "&&":
		c.emitInstruction(OpAnd)
	case "||":
		c.emitInstruction(OpOr)
	case "&":
		c.emitInstruction(OpBitAnd)
	case "|":
		c.emitInstruction(OpBitOr)
	case "^":
		c.emitInstruction(OpBitXor)
	case "<<":
		c.emitInstruction(OpLeftShift)
	case ">>":
		c.emitInstruction(OpRightShift)
	default:
		c.addError("unknown infix operator: %s", expr.Operator)
	}
}

// compileCallExpression compiles a function call
func (c *Compiler) compileCallExpression(expr *ast.CallExpression) {
	c.compileExpression(expr.Function)

	// Compile arguments
	for _, arg := range expr.Arguments {
		c.compileExpression(arg)
	}

	// Emit call instruction
	c.emitInstruction(OpCall, len(expr.Arguments))
}

// compileIndexExpression compiles an index expression
func (c *Compiler) compileIndexExpression(expr *ast.IndexExpression) {
	c.compileExpression(expr.Left)
	c.compileExpression(expr.Index)
	c.emitInstruction(OpIndex)
}

// compileFunctionLiteral compiles a function literal
func (c *Compiler) compileFunctionLiteral(expr *ast.FunctionLiteral) {
	// For now, emit a placeholder
	idx := c.addConstant(fmt.Sprintf("func(%d)", len(expr.Parameters)))
	c.emitInstruction(OpFunction, idx)
}

// compileAssignmentExpression compiles an assignment expression
func (c *Compiler) compileAssignmentExpression(expr *ast.AssignmentExpression) {
	c.compileExpression(expr.Value)

	if ident, ok := expr.Left.(*ast.Identifier); ok {
		if idx, ok := c.locals[ident.Value]; ok {
			c.emitInstruction(OpSetLocal, idx)
		} else {
			c.emitInstruction(OpSetGlobal, 0) // Simplified
		}
	}
}

// Helper methods

// addConstant adds a constant to the constants pool
func (c *Compiler) addConstant(obj interface{}) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emitInstruction emits an instruction to bytecode
func (c *Compiler) emitInstruction(op OpCode, operands ...int) int {
	instructions := EncodeInstructions(op, operands...)
	startPos := len(c.bytecode)
	c.bytecode = append(c.bytecode, instructions...)
	return startPos
}

// patchJump patches a jump instruction at the given position
func (c *Compiler) patchJump(pos int) {
	// This is simplified - in a real implementation, we would properly patch addresses
	jumpAddress := len(c.bytecode)

	// Encode jump address
	c.bytecode[pos+1] = byte((jumpAddress >> 24) & 0xFF)
	c.bytecode[pos+2] = byte((jumpAddress >> 16) & 0xFF)
	c.bytecode[pos+3] = byte((jumpAddress >> 8) & 0xFF)
	c.bytecode[pos+4] = byte(jumpAddress & 0xFF)
}

// addError adds a compilation error
func (c *Compiler) addError(format string, args ...interface{}) {
	c.errors = append(c.errors, fmt.Sprintf(format, args...))
}
