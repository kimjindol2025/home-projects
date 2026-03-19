package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"fv2-lang/internal/lexer"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fv2 [options] <file.fv>\n")
		flag.PrintDefaults()
	}

	tokenize := flag.Bool("tokenize", false, "Show tokens only")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No input file specified\n")
		flag.Usage()
		os.Exit(1)
	}

	filename := args[0]

	// Read source file
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	// Compile
	sourceStr := string(source)
	if err := compile(sourceStr, *tokenize); err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error: %v\n", err)
		os.Exit(1)
	}
}

func compile(source string, tokensOnly bool) error {
	// Step 1: Lexing
	lex, err := lexer.New(source)
	if err != nil {
		return fmt.Errorf("lexer initialization failed: %w", err)
	}

	tokens, err := lex.Tokenize()
	if err != nil {
		return fmt.Errorf("tokenization failed: %w", err)
	}

	if tokensOnly {
		fmt.Printf("=== Tokens ===\n")
		for _, token := range tokens {
			fmt.Printf("%v\n", token)
		}
		return nil
	}

	// Step 2: Parser (TODO)
	// Step 3: Type Checker (TODO)
	// Step 4: Code Generator (TODO)

	fmt.Printf("// FV 2.0 Compiler\n")
	fmt.Printf("// Tokenized %d tokens\n", len(tokens))
	fmt.Printf("// Parser: NOT YET IMPLEMENTED\n")
	fmt.Printf("// C code generation: NOT YET IMPLEMENTED\n")

	return nil
}
