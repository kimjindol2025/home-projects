package formatter

import (
	"fmt"
	"strings"

	"github.com/freelang-ai/gofree/internal/ast"
)

// Formatter formats FreeLang source code
type Formatter struct {
	source string
	indent int
}

// NewFormatter creates a new formatter
func NewFormatter(source string) *Formatter {
	return &Formatter{source: source, indent: 0}
}

// Format formats the source code
func (f *Formatter) Format(module *ast.Module) (string, error) {
	if module == nil {
		return "", fmt.Errorf("module is nil")
	}

	var result strings.Builder

	// Format statements
	for _, stmt := range module.Statements {
		f.formatStatement(stmt, &result)
	}

	return result.String(), nil
}

// formatStatement formats a statement
func (f *Formatter) formatStatement(stmt ast.Statement, result *strings.Builder) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		f.formatVariableDeclaration(s, result)
	case *ast.FunctionStatement:
		f.formatFunctionStatement(s, result)
	case *ast.IfStatement:
		f.formatIfStatement(s, result)
	case *ast.ForStatement:
		f.formatForStatement(s, result)
	case *ast.WhileStatement:
		f.formatWhileStatement(s, result)
	case *ast.BlockStatement:
		f.formatBlockStatement(s, result)
	case *ast.ExpressionStatement:
		f.formatExpressionStatement(s, result)
	case *ast.ReturnStatement:
		f.formatReturnStatement(s, result)
	}
}

func (f *Formatter) formatVariableDeclaration(decl *ast.VariableDeclaration, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString(fmt.Sprintf("%s %s", decl.Kind, decl.Name))
	if decl.Type != "" {
		result.WriteString(fmt.Sprintf(": %s", decl.Type))
	}
	if decl.Value != nil {
		result.WriteString(" = ")
		f.formatExpression(decl.Value, result)
	}
	result.WriteString("\n")
}

func (f *Formatter) formatFunctionStatement(fn *ast.FunctionStatement, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString(fmt.Sprintf("fn %s(", fn.Name))

	for i, param := range fn.Parameters {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.Name)
		if param.Type != "" {
			result.WriteString(fmt.Sprintf(": %s", param.Type))
		}
	}

	result.WriteString(")")
	if fn.ReturnType != "" {
		result.WriteString(fmt.Sprintf(" -> %s", fn.ReturnType))
	}
	result.WriteString(" ")

	if fn.Body != nil {
		f.formatBlockStatement(fn.Body, result)
	}
	result.WriteString("\n")
}

func (f *Formatter) formatIfStatement(stmt *ast.IfStatement, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString("if (")
	f.formatExpression(stmt.Condition, result)
	result.WriteString(") ")

	if stmt.ThenBranch != nil {
		f.formatStatement(stmt.ThenBranch, result)
	}

	if stmt.ElseBranch != nil {
		f.writeIndent(result)
		result.WriteString("else ")
		f.formatStatement(stmt.ElseBranch, result)
	}
}

func (f *Formatter) formatForStatement(stmt *ast.ForStatement, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString("for (")

	if stmt.Init != nil {
		f.formatStatement(stmt.Init, result)
	}
	result.WriteString("; ")

	if stmt.Condition != nil {
		f.formatExpression(stmt.Condition, result)
	}
	result.WriteString("; ")

	if stmt.Update != nil {
		f.formatExpression(stmt.Update, result)
	}
	result.WriteString(") ")

	if stmt.Body != nil {
		f.formatStatement(stmt.Body, result)
	}
	result.WriteString("\n")
}

func (f *Formatter) formatWhileStatement(stmt *ast.WhileStatement, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString("while (")
	f.formatExpression(stmt.Condition, result)
	result.WriteString(") ")

	if stmt.Body != nil {
		f.formatStatement(stmt.Body, result)
	}
	result.WriteString("\n")
}

func (f *Formatter) formatBlockStatement(stmt *ast.BlockStatement, result *strings.Builder) {
	result.WriteString("{\n")
	f.indent++

	for _, s := range stmt.Statements {
		f.formatStatement(s, result)
	}

	f.indent--
	f.writeIndent(result)
	result.WriteString("}\n")
}

func (f *Formatter) formatExpressionStatement(stmt *ast.ExpressionStatement, result *strings.Builder) {
	f.writeIndent(result)
	f.formatExpression(stmt.Expression, result)
	result.WriteString("\n")
}

func (f *Formatter) formatReturnStatement(stmt *ast.ReturnStatement, result *strings.Builder) {
	f.writeIndent(result)
	result.WriteString("return")
	if stmt.Value != nil {
		result.WriteString(" ")
		f.formatExpression(stmt.Value, result)
	}
	result.WriteString("\n")
}

func (f *Formatter) formatExpression(expr ast.Expression, result *strings.Builder) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.LiteralExpression:
		switch v := e.Value.(type) {
		case string:
			result.WriteString(fmt.Sprintf("\"%s\"", v))
		default:
			result.WriteString(fmt.Sprintf("%v", v))
		}

	case *ast.IdentifierExpression:
		result.WriteString(e.Name)

	case *ast.BinaryOpExpression:
		f.formatExpression(e.Left, result)
		result.WriteString(fmt.Sprintf(" %s ", e.Operator))
		f.formatExpression(e.Right, result)

	case *ast.CallExpression:
		result.WriteString(e.Callee)
		result.WriteString("(")
		for i, arg := range e.Arguments {
			if i > 0 {
				result.WriteString(", ")
			}
			f.formatExpression(arg, result)
		}
		result.WriteString(")")

	case *ast.ArrayExpression:
		result.WriteString("[")
		for i, elem := range e.Elements {
			if i > 0 {
				result.WriteString(", ")
			}
			f.formatExpression(elem, result)
		}
		result.WriteString("]")
	}
}

func (f *Formatter) writeIndent(result *strings.Builder) {
	for i := 0; i < f.indent; i++ {
		result.WriteString("    ")
	}
}

// FormatSimple formats code with simple rules
func (f *Formatter) FormatSimple() string {
	// Simple formatting: normalize whitespace
	lines := strings.Split(f.source, "\n")
	var result []string

	for _, line := range lines {
		// Trim trailing whitespace
		trimmed := strings.TrimRight(line, " \t")
		// Skip empty lines (or keep them)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}
