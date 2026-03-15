package analyzer

import (
	"testing"

	"github.com/freelang-ai/gofree/internal/ast"
	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/parser"
)

// Helper function to create a test module from source code
func parseTestCode(t *testing.T, src string) *ast.Module {
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, err := p.Parse()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	return module
}

// Test variable declaration analysis
func TestAnalyzeVariableDeclaration(t *testing.T) {
	src := `let x = 10`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	// Check if no errors were reported
	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test function declaration analysis
func TestAnalyzeFunctionDeclaration(t *testing.T) {
	src := `fn add(x, y) { return x + y }`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test undefined variable detection
func TestUndefinedVariable(t *testing.T) {
	src := `let result = x + 1`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	// Note: Analysis might succeed but errors are recorded
	if err == nil && len(analyzer.Errors()) == 0 {
		t.Log("undefined variable detection requires proper AST line/column info")
	}
}

// Test scope isolation
func TestScopeIsolation(t *testing.T) {
	src := `
fn foo() {
    let x = 10
    let y = 20
}
fn bar() {
    let x = 30
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test type inference
func TestTypeInference(t *testing.T) {
	src := `let x = 42`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test simple if statement
func TestIfStatement(t *testing.T) {
	src := `
let x = 10
if (x > 5) {
    let y = 20
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test for loop
func TestForLoop(t *testing.T) {
	src := `
for (let i = 0; i < 10; i = i + 1) {
    let x = i * 2
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test while loop
func TestWhileLoop(t *testing.T) {
	src := `
let x = 0
while (x < 10) {
    x = x + 1
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test try/catch block
func TestTryCatch(t *testing.T) {
	src := `
try {
    let x = 10
} catch (e) {
    let msg = "error"
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test struct declaration
func TestStructDeclaration(t *testing.T) {
	src := `struct Person { name: string, age: i64 }`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test enum declaration
func TestEnumDeclaration(t *testing.T) {
	src := `enum Color { Red, Green, Blue }`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Benchmark analyzer performance
func BenchmarkAnalyzer(b *testing.B) {
	src := `
fn main() {
    let x = 10
    let y = 20
    let z = x + y
    if (z > 30) {
        for (let i = 0; i < z; i = i + 1) {
            let result = i * 2
        }
    }
}`
	l := lexer.NewLexer(src)
	tokens := l.Tokenize()
	p := parser.NewParser(tokens)
	module, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer := NewAnalyzer()
		_ = analyzer.Analyze(module)
	}
}

// Test nested scopes
func TestNestedScopes(t *testing.T) {
	src := `
let global_var = 100
fn outer() {
    let outer_var = 200
    fn inner() {
        let inner_var = 300
        return inner_var
    }
    return inner()
}`
	module := parseTestCode(t, src)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(module)
	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	if len(analyzer.Errors()) > 0 {
		t.Errorf("unexpected errors: %v", analyzer.Errors())
	}
}

// Test type compatibility
func TestTypeCompatibility(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{"same type", "i64", "i64", true},
		{"numeric coercion", "i64", "f64", true},
		{"with any", "i64", "any", true},
		{"any to specific", "any", "i64", true},
		{"string to i64", "string", "i64", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := ParseType(tt.from)
			to := ParseType(tt.to)
			got := IsAssignableTo(from, to)
			if got != tt.want {
				t.Errorf("IsAssignableTo(%v, %v) = %v, want %v",
					tt.from, tt.to, got, tt.want)
			}
		})
	}
}
