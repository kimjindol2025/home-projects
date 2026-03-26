package checker

import (
	"fmt"

	"github.com/kimjindol2025/freelang-go/pkg/ast"
)

// Checker performs semantic analysis and type checking
type Checker struct {
	symtab *SymbolTable
	errors []string
}

// New creates a new type checker
func New() *Checker {
	return &Checker{
		symtab: NewSymbolTable(),
		errors: []string{},
	}
}

// Errors returns all type checking errors
func (c *Checker) Errors() []string {
	return c.errors
}

// GlobalScope returns the global symbol table
func (c *Checker) GlobalScope() *SymbolTable {
	return c.symtab
}

// Check performs type checking on the entire program
func (c *Checker) Check(program *ast.Program) {
	for _, stmt := range program.Statements {
		c.checkStatement(stmt)
	}
}

// checkStatement checks a statement
func (c *Checker) checkStatement(stmt ast.Statement) Type {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return c.checkLetStatement(s)
	case *ast.ReturnStatement:
		return c.checkReturnStatement(s)
	case *ast.ExpressionStatement:
		return c.checkExpressionStatement(s)
	case *ast.IfStatement:
		return c.checkIfStatement(s)
	case *ast.ForStatement:
		return c.checkForStatement(s)
	case *ast.BlockStatement:
		return c.checkBlockStatement(s)
	default:
		c.addError("unknown statement type: %T", stmt)
		return TypeVoid
	}
}

// checkLetStatement checks a let statement
func (c *Checker) checkLetStatement(stmt *ast.LetStatement) Type {
	// Check if variable already exists in current scope
	if _, ok := c.symtab.ResolveLocal(stmt.Name.Value); ok {
		c.addError("variable %q already defined in this scope", stmt.Name.Value)
		return TypeVoid
	}

	// Infer type from value
	valueType := c.checkExpression(stmt.Value)
	if valueType == nil {
		valueType = TypeAny
	}

	// Define the variable
	c.symtab.Define(
		stmt.Name.Value,
		SymbolVariable,
		valueType,
		true, // mutable
	)

	return TypeVoid
}

// checkReturnStatement checks a return statement
func (c *Checker) checkReturnStatement(stmt *ast.ReturnStatement) Type {
	if stmt.ReturnValue == nil {
		return TypeVoid
	}

	return c.checkExpression(stmt.ReturnValue)
}

// checkExpressionStatement checks an expression statement
func (c *Checker) checkExpressionStatement(stmt *ast.ExpressionStatement) Type {
	return c.checkExpression(stmt.Expression)
}

// checkIfStatement checks an if statement
func (c *Checker) checkIfStatement(stmt *ast.IfStatement) Type {
	// Check condition
	condType := c.checkExpression(stmt.Condition)
	if condType != nil && !TypesEqual(condType, TypeBool) && !TypesEqual(condType, TypeAny) {
		c.addError("if condition must be boolean, got %s", condType.String())
	}

	// Check consequence
	c.checkBlockStatement(stmt.Consequence)

	// Check alternative
	if stmt.Alternative != nil {
		c.checkBlockStatement(stmt.Alternative)
	}

	return TypeVoid
}

// checkForStatement checks a for statement
func (c *Checker) checkForStatement(stmt *ast.ForStatement) Type {
	// Check iterable
	iterableType := c.checkExpression(stmt.Iterable)

	// Push new scope for loop variable
	c.symtab.PushScope()

	// Define loop variable
	if stmt.LoopVar != nil {
		var elementType Type = TypeAny

		// If iterable is array, element type is array's element type
		if arrayType, ok := iterableType.(*ArrayType); ok {
			elementType = arrayType.ElementType
		}

		c.symtab.Define(
			stmt.LoopVar.Value,
			SymbolVariable,
			elementType,
			true,
		)
	}

	// Check loop body
	c.checkBlockStatement(stmt.Body)

	// Pop loop scope
	c.symtab.PopScope()

	return TypeVoid
}

// checkBlockStatement checks a block statement
func (c *Checker) checkBlockStatement(stmt *ast.BlockStatement) Type {
	var lastType Type = TypeVoid

	for _, s := range stmt.Statements {
		lastType = c.checkStatement(s)
	}

	return lastType
}

// checkExpression checks an expression and returns its type
func (c *Checker) checkExpression(expr ast.Expression) Type {
	if expr == nil {
		return TypeAny
	}

	switch e := expr.(type) {
	case *ast.Identifier:
		return c.checkIdentifier(e)
	case *ast.IntegerLiteral:
		return TypeInt
	case *ast.StringLiteral:
		return TypeString
	case *ast.BooleanLiteral:
		return TypeBool
	case *ast.ArrayLiteral:
		return c.checkArrayLiteral(e)
	case *ast.HashLiteral:
		return c.checkHashLiteral(e)
	case *ast.FunctionLiteral:
		return c.checkFunctionLiteral(e)
	case *ast.PrefixExpression:
		return c.checkPrefixExpression(e)
	case *ast.InfixExpression:
		return c.checkInfixExpression(e)
	case *ast.CallExpression:
		return c.checkCallExpression(e)
	case *ast.IndexExpression:
		return c.checkIndexExpression(e)
	case *ast.AssignmentExpression:
		return c.checkAssignmentExpression(e)
	default:
		c.addError("unknown expression type: %T", expr)
		return TypeAny
	}
}

// checkIdentifier checks an identifier
func (c *Checker) checkIdentifier(expr *ast.Identifier) Type {
	symbol, ok := c.symtab.Resolve(expr.Value)
	if !ok {
		c.addError("undefined variable: %q", expr.Value)
		return TypeAny
	}

	if symbol.Type == nil {
		return TypeAny
	}

	return symbol.Type
}

// checkArrayLiteral checks an array literal
func (c *Checker) checkArrayLiteral(expr *ast.ArrayLiteral) Type {
	if len(expr.Elements) == 0 {
		return NewArrayType(TypeAny)
	}

	// Infer element type from first element
	elementType := c.checkExpression(expr.Elements[0])

	// Check that all elements have compatible type
	for i, elem := range expr.Elements[1:] {
		elemType := c.checkExpression(elem)
		if !IsCompatible(elemType, elementType) {
			c.addError("element %d has incompatible type: expected %s, got %s",
				i+1, elementType.String(), elemType.String())
		}
	}

	return NewArrayType(elementType)
}

// checkHashLiteral checks a hash literal
func (c *Checker) checkHashLiteral(expr *ast.HashLiteral) Type {
	if len(expr.Pairs) == 0 {
		return NewHashType(TypeAny, TypeAny)
	}

	var keyType, valueType Type

	for k, v := range expr.Pairs {
		// For identifiers as keys, treat them as strings (object property keys)
		var kt Type
		if _, isIdent := k.(*ast.Identifier); isIdent {
			kt = TypeString
		} else {
			kt = c.checkExpression(k)
		}

		vt := c.checkExpression(v)

		if keyType == nil {
			keyType = kt
		} else if !IsCompatible(kt, keyType) {
			c.addError("hash key has incompatible type: expected %s, got %s",
				keyType.String(), kt.String())
		}

		if valueType == nil {
			valueType = vt
		} else if !IsCompatible(vt, valueType) {
			c.addError("hash value has incompatible type: expected %s, got %s",
				valueType.String(), vt.String())
		}
	}

	if keyType == nil {
		keyType = TypeAny
	}
	if valueType == nil {
		valueType = TypeAny
	}

	return NewHashType(keyType, valueType)
}

// checkFunctionLiteral checks a function literal
func (c *Checker) checkFunctionLiteral(expr *ast.FunctionLiteral) Type {
	// Push new scope for function
	c.symtab.PushScope()

	// Define parameters
	paramTypes := make([]Type, len(expr.Parameters))
	for i, param := range expr.Parameters {
		paramTypes[i] = TypeAny // Parameters are inferred as Any by default
		c.symtab.Define(param.Value, SymbolParameter, TypeAny, false)
	}

	// Check function body
	returnType := c.checkBlockStatement(expr.Body)
	if returnType == nil {
		returnType = TypeVoid
	}

	// Pop function scope
	c.symtab.PopScope()

	return NewFunctionType(paramTypes, returnType)
}

// checkPrefixExpression checks a prefix expression
func (c *Checker) checkPrefixExpression(expr *ast.PrefixExpression) Type {
	rightType := c.checkExpression(expr.Right)

	switch expr.Operator {
	case "!":
		if !TypesEqual(rightType, TypeBool) && !TypesEqual(rightType, TypeAny) {
			c.addError("! operator requires boolean, got %s", rightType.String())
		}
		return TypeBool
	case "-":
		if !TypesEqual(rightType, TypeInt) && !TypesEqual(rightType, TypeFloat) && !TypesEqual(rightType, TypeAny) {
			c.addError("- operator requires number, got %s", rightType.String())
		}
		return rightType
	case "+":
		if !TypesEqual(rightType, TypeInt) && !TypesEqual(rightType, TypeFloat) && !TypesEqual(rightType, TypeAny) {
			c.addError("+ operator requires number, got %s", rightType.String())
		}
		return rightType
	case "~":
		if !TypesEqual(rightType, TypeInt) && !TypesEqual(rightType, TypeAny) {
			c.addError("~ operator requires int, got %s", rightType.String())
		}
		return TypeInt
	default:
		c.addError("unknown prefix operator: %s", expr.Operator)
		return TypeAny
	}
}

// checkInfixExpression checks an infix expression
func (c *Checker) checkInfixExpression(expr *ast.InfixExpression) Type {
	leftType := c.checkExpression(expr.Left)
	rightType := c.checkExpression(expr.Right)

	switch expr.Operator {
	case "+", "-", "*", "/", "%":
		// Arithmetic operators require numbers
		if !isNumber(leftType) || !isNumber(rightType) {
			c.addError("%s operator requires numbers, got %s and %s",
				expr.Operator, leftType.String(), rightType.String())
		}
		return leftType
	case "==", "!=":
		// Equality operators work on all types
		return TypeBool
	case "<", ">", "<=", ">=":
		// Comparison operators require numbers
		if !isNumber(leftType) || !isNumber(rightType) {
			c.addError("%s operator requires numbers, got %s and %s",
				expr.Operator, leftType.String(), rightType.String())
		}
		return TypeBool
	case "&&", "||":
		// Logical operators require booleans
		if !TypesEqual(leftType, TypeBool) || !TypesEqual(rightType, TypeBool) {
			c.addError("%s operator requires booleans, got %s and %s",
				expr.Operator, leftType.String(), rightType.String())
		}
		return TypeBool
	default:
		c.addError("unknown infix operator: %s", expr.Operator)
		return TypeAny
	}
}

// checkCallExpression checks a function call
func (c *Checker) checkCallExpression(expr *ast.CallExpression) Type {
	funcType := c.checkExpression(expr.Function)

	if funcType == nil {
		return TypeAny
	}

	// Function type check
	if funcType, ok := funcType.(*FunctionType); ok {
		// Check argument count
		if len(expr.Arguments) != len(funcType.ParameterTypes) {
			c.addError("function call expects %d arguments, got %d",
				len(funcType.ParameterTypes), len(expr.Arguments))
		}

		// Check argument types
		for i, arg := range expr.Arguments {
			argType := c.checkExpression(arg)
			if i < len(funcType.ParameterTypes) {
				if !IsCompatible(argType, funcType.ParameterTypes[i]) {
					c.addError("argument %d has wrong type: expected %s, got %s",
						i, funcType.ParameterTypes[i].String(), argType.String())
				}
			}
		}

		return funcType.ReturnType
	}

	c.addError("calling non-function type: %s", funcType.String())
	return TypeAny
}

// checkIndexExpression checks an index expression
func (c *Checker) checkIndexExpression(expr *ast.IndexExpression) Type {
	leftType := c.checkExpression(expr.Left)
	indexType := c.checkExpression(expr.Index)

	switch left := leftType.(type) {
	case *ArrayType:
		if !TypesEqual(indexType, TypeInt) && !TypesEqual(indexType, TypeAny) {
			c.addError("array index must be int, got %s", indexType.String())
		}
		return left.ElementType

	case *HashType:
		if !IsCompatible(indexType, left.KeyType) {
			c.addError("hash key has wrong type: expected %s, got %s",
				left.KeyType.String(), indexType.String())
		}
		return left.ValueType

	default:
		c.addError("cannot index type %s", leftType.String())
		return TypeAny
	}
}

// checkAssignmentExpression checks an assignment expression
func (c *Checker) checkAssignmentExpression(expr *ast.AssignmentExpression) Type {
	// Get target type
	var targetType Type

	if ident, ok := expr.Left.(*ast.Identifier); ok {
		symbol, ok := c.symtab.Resolve(ident.Value)
		if !ok {
			c.addError("undefined variable: %q", ident.Value)
			return TypeAny
		}

		if !symbol.Mutable {
			c.addError("cannot assign to immutable variable: %q", ident.Value)
			return TypeAny
		}

		targetType = symbol.Type
	} else {
		// For complex assignments (array[i] = value, obj.prop = value), infer from value
		targetType = c.checkExpression(expr.Left)
	}

	// Check value type compatibility
	valueType := c.checkExpression(expr.Value)
	if !IsCompatible(valueType, targetType) {
		c.addError("cannot assign %s to %s", valueType.String(), targetType.String())
	}

	return targetType
}

// Helper functions

func isNumber(t Type) bool {
	return TypesEqual(t, TypeInt) || TypesEqual(t, TypeFloat) || TypesEqual(t, TypeAny)
}

func (c *Checker) addError(format string, args ...interface{}) {
	c.errors = append(c.errors, fmt.Sprintf(format, args...))
}
