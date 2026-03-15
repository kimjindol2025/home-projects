package main

import (
	"fmt"
	"os"

	"github.com/freelang-ai/gofree/internal/compiler"
	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/parser"
	"github.com/freelang-ai/gofree/internal/runtime"
)

// Bootstrap implements self-hosting compilation
// This program compiles FreeLang source code using the Go-based compiler
// and can be the foundation for a pure-FreeLang bootstrap

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bootstrap <input-file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Read source file
	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Compile FreeLang code
	result, err := compile(string(source))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

// compile compiles FreeLang source to IR and executes it
func compile(source string) (string, error) {
	// Phase 1: Lexical analysis
	lexer := lexer.NewLexer(source)
	tokens := lexer.Tokenize()

	// Phase 2: Parsing
	parser := parser.NewParser(tokens)
	module, err := parser.Parse()
	if err != nil {
		return "", fmt.Errorf("parse error: %v", err)
	}

	// Phase 3: Compilation to IR
	compiler := compiler.NewCompiler()
	program, err := compiler.Compile(module)
	if err != nil {
		return "", fmt.Errorf("compilation error: %v", err)
	}

	// Phase 4: Execution
	vm := runtime.NewVM(program)
	err = vm.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}

	// Phase 5: Return result
	return fmt.Sprintf("✅ Compilation successful\n" +
		"Compiled instructions: %d\n" +
		"Functions: %d\n" +
		"Status: Ready for execution",
		len(program.Instructions), len(program.Functions)), nil
}
