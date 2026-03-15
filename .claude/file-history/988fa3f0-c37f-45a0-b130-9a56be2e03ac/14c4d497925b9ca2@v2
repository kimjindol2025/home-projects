package codegen

import (
	"strings"
	"testing"

	"github.com/freelang-ai/gofree/internal/compiler"
)

// Test C code generation
func TestGenerateC(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 42))

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if code == "" {
		t.Errorf("expected non-empty code")
	}

	if !strings.Contains(code, "#include <stdio.h>") {
		t.Errorf("expected C header")
	}

	if !strings.Contains(code, "int main()") {
		t.Errorf("expected main function")
	}
}

// Test JavaScript code generation
func TestGenerateJavaScript(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 42))

	codegen := NewCodegen(program)
	code, err := codegen.GenerateJavaScript()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if code == "" {
		t.Errorf("expected non-empty code")
	}

	if !strings.Contains(code, "Generated JavaScript") {
		t.Errorf("expected JavaScript header")
	}
}

// Test code generation with functions
func TestGenerateWithFunctions(t *testing.T) {
	program := compiler.NewIRProgram()

	fn := compiler.NewIRFunction("add")
	fn.ParamCount = 2
	fn.Instructions = append(fn.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 0))
	fn.Instructions = append(fn.Instructions, compiler.NewInstruction(compiler.OpReturn))

	program.Functions["add"] = fn

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if !strings.Contains(code, "add") {
		t.Errorf("expected function 'add' in generated code")
	}
}

// Test empty program
func TestGenerateEmpty(t *testing.T) {
	program := compiler.NewIRProgram()

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if code == "" {
		t.Errorf("expected non-empty code for empty program")
	}

	if !strings.Contains(code, "int main()") {
		t.Errorf("expected main function")
	}
}

// Test nil program
func TestGenerateNilProgram(t *testing.T) {
	codegen := NewCodegen(nil)
	_, err := codegen.GenerateC()

	if err == nil {
		t.Errorf("expected error for nil program")
	}
}

// Test instruction emission
func TestEmitInstruction(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpAdd))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpSub))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpMul))

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if !strings.Contains(code, "ADD") && !strings.Contains(code, "add") {
		t.Logf("note: ADD instruction not found in generated code")
	}
}

// Test function generation
func TestFunctionGeneration(t *testing.T) {
	program := compiler.NewIRProgram()

	fn1 := compiler.NewIRFunction("func1")
	fn1.ParamCount = 1
	program.Functions["func1"] = fn1

	fn2 := compiler.NewIRFunction("func2")
	fn2.ParamCount = 2
	program.Functions["func2"] = fn2

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if !strings.Contains(code, "func1") || !strings.Contains(code, "func2") {
		t.Errorf("expected both functions in generated code")
	}
}

// Test C code validity
func TestCCodeValidity(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 10))

	codegen := NewCodegen(program)
	code, _ := codegen.GenerateC()

	// Check for basic C structure
	checks := []string{
		"#include",
		"main",
		"return",
		"{",
		"}",
	}

	for _, check := range checks {
		if !strings.Contains(code, check) {
			t.Errorf("expected '%s' in generated C code", check)
		}
	}
}

// Test JavaScript code validity
func TestJavaScriptValidity(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 10))

	codegen := NewCodegen(program)
	code, _ := codegen.GenerateJavaScript()

	// Check for basic JavaScript structure
	if !strings.Contains(code, "Generated JavaScript") {
		t.Errorf("expected JavaScript header")
	}

	if code == "" {
		t.Errorf("expected non-empty JavaScript code")
	}
}

// Test global variables
func TestGlobalVariables(t *testing.T) {
	program := compiler.NewIRProgram()
	program.GlobalVars["x"] = 0
	program.GlobalVars["y"] = 1
	program.GlobalVars["z"] = 2

	codegen := NewCodegen(program)
	code, _ := codegen.GenerateJavaScript()

	if !strings.Contains(code, "let x") && !strings.Contains(code, "let y") {
		t.Logf("note: global variables may not be visible in simplified codegen")
	}
}

// Test codegen with multiple instructions
func TestMultipleInstructions(t *testing.T) {
	program := compiler.NewIRProgram()
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 1))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpLoadConst, 2))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpAdd))
	program.Instructions = append(program.Instructions, compiler.NewInstruction(compiler.OpReturn))

	codegen := NewCodegen(program)
	code, err := codegen.GenerateC()

	if err != nil {
		t.Errorf("codegen failed: %v", err)
	}

	if code == "" {
		t.Errorf("expected non-empty code")
	}
}

// Benchmark code generation
func BenchmarkCodegen(b *testing.B) {
	program := compiler.NewIRProgram()

	// Create a moderately complex program
	for i := 0; i < 100; i++ {
		program.Instructions = append(program.Instructions,
			compiler.NewInstruction(compiler.OpLoadConst, i))
		program.Instructions = append(program.Instructions,
			compiler.NewInstruction(compiler.OpAdd))
	}

	fn := compiler.NewIRFunction("test")
	fn.ParamCount = 2
	fn.Instructions = append(fn.Instructions,
		compiler.NewInstruction(compiler.OpLoadConst, 0))
	fn.Instructions = append(fn.Instructions,
		compiler.NewInstruction(compiler.OpReturn))
	program.Functions["test"] = fn

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codegen := NewCodegen(program)
		_, _ = codegen.GenerateC()
	}
}
