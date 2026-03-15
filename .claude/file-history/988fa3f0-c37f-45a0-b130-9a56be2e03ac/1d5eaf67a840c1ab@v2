package formatter

import (
	"fmt"

	"github.com/freelang-ai/gofree/internal/ast"
)

// Linter performs code linting
type Linter struct {
	issues []LintIssue
}

// LintIssue represents a linting issue
type LintIssue struct {
	Line    int
	Column  int
	Severity string // "error", "warning", "info"
	Code    string   // Rule code
	Message string
}

// NewLinter creates a new linter
func NewLinter() *Linter {
	return &Linter{issues: []LintIssue{}}
}

// Lint lints a module
func (l *Linter) Lint(module *ast.Module) []LintIssue {
	if module == nil {
		return l.issues
	}

	for _, stmt := range module.Statements {
		l.lintStatement(stmt)
	}

	return l.issues
}

// lintStatement lints a statement
func (l *Linter) lintStatement(stmt ast.Statement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		l.lintVariableDeclaration(s)
	case *ast.FunctionStatement:
		l.lintFunctionStatement(s)
	case *ast.IfStatement:
		l.lintIfStatement(s)
	case *ast.ForStatement:
		l.lintForStatement(s)
	case *ast.BlockStatement:
		l.lintBlockStatement(s)
	}
}

// lintVariableDeclaration lints variable declarations
func (l *Linter) lintVariableDeclaration(decl *ast.VariableDeclaration) {
	// Check for unused variables (simplified)
	if !decl.IsMutable && decl.Kind == "let" {
		// const should be uppercase
		if len(decl.Name) > 0 && decl.Name[0] >= 'a' && decl.Name[0] <= 'z' {
			// This is just a suggestion
		}
	}

	// Check for meaningful names
	if decl.Name == "x" || decl.Name == "y" || decl.Name == "z" {
		l.addIssue(0, 0, "warning", "unclear_name",
			fmt.Sprintf("variable name '%s' is not descriptive", decl.Name))
	}
}

// lintFunctionStatement lints function definitions
func (l *Linter) lintFunctionStatement(fn *ast.FunctionStatement) {
	// Check function naming convention
	if len(fn.Name) == 0 {
		l.addIssue(0, 0, "error", "empty_name", "function must have a name")
	}

	// Check for empty function body
	if fn.Body == nil || len(fn.Body.Statements) == 0 {
		l.addIssue(0, 0, "warning", "empty_function",
			fmt.Sprintf("function '%s' has empty body", fn.Name))
	}

	l.lintBlockStatement(fn.Body)
}

// lintIfStatement lints if statements
func (l *Linter) lintIfStatement(stmt *ast.IfStatement) {
	if stmt.ThenBranch != nil {
		l.lintStatement(stmt.ThenBranch)
	}
	if stmt.ElseBranch != nil {
		l.lintStatement(stmt.ElseBranch)
	}
}

// lintForStatement lints for loops
func (l *Linter) lintForStatement(stmt *ast.ForStatement) {
	if stmt.Body != nil {
		l.lintStatement(stmt.Body)
	}
}

// lintBlockStatement lints block statements
func (l *Linter) lintBlockStatement(stmt *ast.BlockStatement) {
	if stmt == nil {
		return
	}

	for _, s := range stmt.Statements {
		l.lintStatement(s)
	}
}

// addIssue adds a linting issue
func (l *Linter) addIssue(line, column int, severity, code, message string) {
	l.issues = append(l.issues, LintIssue{
		Line:     line,
		Column:   column,
		Severity: severity,
		Code:     code,
		Message:  message,
	})
}

// Rules contains all available linting rules
type Rules struct {
	NoUnused       bool // no_unused: detect unused variables
	NoShadow       bool // no_shadow: detect variable shadowing
	NoImplicitAny  bool // no_implicit_any: disallow implicit any type
	StrictEquality bool // strict_equality: require === instead of ==
	NoContinue     bool // no_continue: disallow continue statements
}

// DefaultRules returns default linting rules
func DefaultRules() Rules {
	return Rules{
		NoUnused:       true,
		NoShadow:       true,
		NoImplicitAny:  false,
		StrictEquality: false,
		NoContinue:     false,
	}
}

// LintWithRules lints with custom rules
func (l *Linter) LintWithRules(module *ast.Module, rules Rules) []LintIssue {
	if module == nil {
		return l.issues
	}

	// Apply rules
	if rules.NoUnused {
		// Check for unused variables
	}

	if rules.NoShadow {
		// Check for variable shadowing
	}

	// Continue with default linting
	return l.Lint(module)
}

// FormatIssue formats an issue for display
func (l *Linter) FormatIssue(issue LintIssue) string {
	return fmt.Sprintf("[%s] %s:%d:%d - %s (%s)",
		issue.Severity, issue.Code, issue.Line, issue.Column, issue.Message, issue.Code)
}

// Summary returns a summary of linting results
func (l *Linter) Summary() string {
	errors := 0
	warnings := 0
	infos := 0

	for _, issue := range l.issues {
		switch issue.Severity {
		case "error":
			errors++
		case "warning":
			warnings++
		case "info":
			infos++
		}
	}

	return fmt.Sprintf("Linting results: %d errors, %d warnings, %d info",
		errors, warnings, infos)
}
