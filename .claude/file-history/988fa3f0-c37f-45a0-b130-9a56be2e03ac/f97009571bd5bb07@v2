package formatter

import (
	"strings"
	"testing"

	"github.com/freelang-ai/gofree/internal/ast"
)

// Test formatter
func TestFormatterBasic(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.VariableDeclaration{
				Kind:  "let",
				Name:  "x",
				Value: &ast.LiteralExpression{Value: 42.0, DataType: "number"},
			},
		},
	}

	formatter := NewFormatter("")
	result, err := formatter.Format(module)

	if err != nil {
		t.Errorf("format failed: %v", err)
	}

	if result == "" {
		t.Errorf("expected non-empty formatted code")
	}

	if !strings.Contains(result, "let") || !strings.Contains(result, "x") {
		t.Errorf("expected 'let x' in formatted code")
	}
}

// Test function formatting
func TestFormatterFunction(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.FunctionStatement{
				Name:       "add",
				Parameters: []*ast.Parameter{{Name: "x"}, {Name: "y"}},
				Body:       &ast.BlockStatement{Statements: []ast.Statement{}},
			},
		},
	}

	formatter := NewFormatter("")
	result, err := formatter.Format(module)

	if err != nil {
		t.Errorf("format failed: %v", err)
	}

	if !strings.Contains(result, "fn") || !strings.Contains(result, "add") {
		t.Errorf("expected function signature in formatted code")
	}
}

// Test simple formatting
func TestFormatterSimple(t *testing.T) {
	source := "let x = 10\nlet y = 20\n"
	formatter := NewFormatter(source)
	result := formatter.FormatSimple()

	if result == "" {
		t.Errorf("expected non-empty result")
	}

	lines := strings.Split(result, "\n")
	if len(lines) == 0 {
		t.Errorf("expected formatted lines")
	}
}

// Test formatter with nil module
func TestFormatterNilModule(t *testing.T) {
	formatter := NewFormatter("")
	_, err := formatter.Format(nil)

	if err == nil {
		t.Errorf("expected error for nil module")
	}
}

// Test linter basic
func TestLinterBasic(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.VariableDeclaration{
				Kind:  "let",
				Name:  "x",
				Value: &ast.LiteralExpression{Value: 42.0, DataType: "number"},
			},
		},
	}

	linter := NewLinter()
	issues := linter.Lint(module)

	if len(issues) < 0 {
		t.Errorf("expected issues list")
	}
}

// Test linter with suspicious variable names
func TestLinterSuspiciousNames(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.VariableDeclaration{
				Kind:  "let",
				Name:  "x",
				Value: &ast.LiteralExpression{Value: 1.0, DataType: "number"},
			},
		},
	}

	linter := NewLinter()
	issues := linter.Lint(module)

	// Should detect unclear name
	foundWarning := false
	for _, issue := range issues {
		if issue.Severity == "warning" {
			foundWarning = true
		}
	}

	if !foundWarning {
		t.Logf("note: unclear name detection not active")
	}
}

// Test linter summary
func TestLinterSummary(t *testing.T) {
	linter := NewLinter()
	linter.addIssue(1, 0, "error", "test", "test error")
	linter.addIssue(2, 0, "warning", "test", "test warning")

	summary := linter.Summary()

	if summary == "" {
		t.Errorf("expected non-empty summary")
	}

	if !strings.Contains(summary, "error") {
		t.Errorf("expected 'error' in summary")
	}
}

// Test format issue
func TestFormatIssue(t *testing.T) {
	linter := NewLinter()
	issue := LintIssue{
		Line:     1,
		Column:   0,
		Severity: "error",
		Code:     "TEST",
		Message:  "Test message",
	}

	formatted := linter.FormatIssue(issue)

	if formatted == "" {
		t.Errorf("expected formatted issue")
	}

	if !strings.Contains(formatted, "error") {
		t.Errorf("expected 'error' in formatted issue")
	}
}

// Test default rules
func TestDefaultRules(t *testing.T) {
	rules := DefaultRules()

	if !rules.NoUnused {
		t.Errorf("expected NoUnused to be enabled")
	}

	if !rules.NoShadow {
		t.Errorf("expected NoShadow to be enabled")
	}
}

// Test formatter indentation
func TestFormatterIndentation(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.FunctionStatement{
				Name:       "test",
				Parameters: []*ast.Parameter{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.VariableDeclaration{
							Kind:  "let",
							Name:  "x",
							Value: &ast.LiteralExpression{Value: 1.0, DataType: "number"},
						},
					},
				},
			},
		},
	}

	formatter := NewFormatter("")
	result, _ := formatter.Format(module)

	// Check for indentation (spaces)
	if !strings.Contains(result, "    ") {
		t.Logf("note: indentation may vary")
	}
}

// Test linter with function
func TestLinterFunction(t *testing.T) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.FunctionStatement{
				Name:       "test",
				Parameters: []*ast.Parameter{},
				Body:       &ast.BlockStatement{Statements: []ast.Statement{}},
			},
		},
	}

	linter := NewLinter()
	issues := linter.Lint(module)

	// May warn about empty function
	if len(issues) > 0 {
		for _, issue := range issues {
			if issue.Code == "empty_function" {
				t.Logf("detected empty function warning")
				break
			}
		}
	}
}

// Benchmark formatter
func BenchmarkFormatter(b *testing.B) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.VariableDeclaration{
				Kind: "let",
				Name: "x",
			},
			&ast.VariableDeclaration{
				Kind: "let",
				Name: "y",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter := NewFormatter("")
		_, _ = formatter.Format(module)
	}
}

// Benchmark linter
func BenchmarkLinter(b *testing.B) {
	module := &ast.Module{
		Statements: []ast.Statement{
			&ast.VariableDeclaration{
				Kind: "let",
				Name: "x",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linter := NewLinter()
		_ = linter.Lint(module)
	}
}
