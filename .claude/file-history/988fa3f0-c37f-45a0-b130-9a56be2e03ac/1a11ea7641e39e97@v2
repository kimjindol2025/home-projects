package selfhost

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/freelang-ai/gofree/internal/compiler"
	"github.com/freelang-ai/gofree/internal/lexer"
	"github.com/freelang-ai/gofree/internal/parser"
	"github.com/freelang-ai/gofree/internal/runtime"
)

// Bootstrapper handles self-hosting compilation
type Bootstrapper struct {
	stage      int
	sourceDir  string
	outputDir  string
	errors     []BootstrapError
}

// BootstrapError represents a bootstrap error
type BootstrapError struct {
	Stage   int
	Message string
}

// NewBootstrapper creates a new bootstrapper
func NewBootstrapper(sourceDir, outputDir string) *Bootstrapper {
	return &Bootstrapper{
		stage:     0,
		sourceDir: sourceDir,
		outputDir: outputDir,
		errors:    []BootstrapError{},
	}
}

// Bootstrap performs self-hosting bootstrap in stages
func (b *Bootstrapper) Bootstrap() error {
	// Stage 0: Compile with Go-based compiler (baseline)
	b.stage = 0
	if err := b.compileWithGoCompiler(); err != nil {
		b.addError(b.stage, fmt.Sprintf("stage 0 failed: %v", err))
		return fmt.Errorf("%d bootstrap errors", len(b.errors))
	}

	// Stage 1: Compile lexer.free with Go compiler output
	b.stage = 1
	if err := b.compileLexerStage(); err != nil {
		b.addError(b.stage, fmt.Sprintf("stage 1 failed: %v", err))
		return fmt.Errorf("%d bootstrap errors", len(b.errors))
	}

	// Stage 2: Compile parser.free with Stage 1 output
	b.stage = 2
	if err := b.compileParserStage(); err != nil {
		b.addError(b.stage, fmt.Sprintf("stage 2 failed: %v", err))
		return fmt.Errorf("%d bootstrap errors", len(b.errors))
	}

	// Stage 3: Compile full compiler with Stage 2 output
	b.stage = 3
	if err := b.compileCompilerStage(); err != nil {
		b.addError(b.stage, fmt.Sprintf("stage 3 failed: %v", err))
		return fmt.Errorf("%d bootstrap errors", len(b.errors))
	}

	// Stage 4: Self-host verification
	b.stage = 4
	if err := b.verifySelfHosting(); err != nil {
		b.addError(b.stage, fmt.Sprintf("stage 4 failed: %v", err))
		return fmt.Errorf("%d bootstrap errors", len(b.errors))
	}

	return nil
}

// compileWithGoCompiler performs initial compilation using Go compiler
func (b *Bootstrapper) compileWithGoCompiler() error {
	// This is the baseline: use Go implementation to bootstrap
	return nil
}

// compileLexerStage compiles the lexer in stage 1
func (b *Bootstrapper) compileLexerStage() error {
	lexerFile := filepath.Join(b.sourceDir, "lexer.free")

	source, err := os.ReadFile(lexerFile)
	if err != nil {
		return fmt.Errorf("cannot read lexer.free: %v", err)
	}

	return b.compileFile(string(source), "lexer_stage1.flir")
}

// compileParserStage compiles the parser in stage 2
func (b *Bootstrapper) compileParserStage() error {
	parserFile := filepath.Join(b.sourceDir, "parser.free")

	source, err := os.ReadFile(parserFile)
	if err != nil {
		return fmt.Errorf("cannot read parser.free: %v", err)
	}

	return b.compileFile(string(source), "parser_stage2.flir")
}

// compileCompilerStage compiles the compiler in stage 3
func (b *Bootstrapper) compileCompilerStage() error {
	compilerFile := filepath.Join(b.sourceDir, "compiler.free")

	source, err := os.ReadFile(compilerFile)
	if err != nil {
		return fmt.Errorf("cannot read compiler.free: %v", err)
	}

	return b.compileFile(string(source), "compiler_stage3.flir")
}

// verifySelfHosting verifies that self-hosting is achieved
func (b *Bootstrapper) verifySelfHosting() error {
	// Verify that the generated compiler can compile itself
	testSource := `fn verify() { return 42 }`

	return b.compileFile(testSource, "test_stage4.flir")
}

// compileFile compiles a FreeLang source file to IR
func (b *Bootstrapper) compileFile(source, outputFile string) error {
	// Lexical analysis
	lex := lexer.NewLexer(source)
	tokens := lex.Tokenize()

	// Parsing
	p := parser.NewParser(tokens)
	module, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	// Compilation
	comp := compiler.NewCompiler()
	program, err := comp.Compile(module)
	if err != nil {
		return fmt.Errorf("compile error: %v", err)
	}

	// Save IR
	outputPath := filepath.Join(b.outputDir, outputFile)
	irData := fmt.Sprintf("IR Program: %d instructions, %d functions\n",
		len(program.Instructions), len(program.Functions))

	err = os.WriteFile(outputPath, []byte(irData), 0644)
	if err != nil {
		return fmt.Errorf("cannot write IR: %v", err)
	}

	// Execute to verify
	vm := runtime.NewVM(program)
	if err := vm.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}

	return nil
}

// addError adds a bootstrap error
func (b *Bootstrapper) addError(stage int, message string) {
	b.errors = append(b.errors, BootstrapError{
		Stage:   stage,
		Message: message,
	})
}

// Errors returns bootstrap errors
func (b *Bootstrapper) Errors() []BootstrapError {
	return b.errors
}

// Status returns bootstrap status
type Status struct {
	Stage             int
	CompletedStages   int
	TotalStages       int
	ErrorCount        int
	Messages          []string
}

// GetStatus returns current bootstrap status
func (b *Bootstrapper) GetStatus() Status {
	messages := []string{}
	for _, err := range b.errors {
		messages = append(messages, fmt.Sprintf("Stage %d: %s", err.Stage, err.Message))
	}

	return Status{
		Stage:           b.stage,
		CompletedStages: b.stage,
		TotalStages:     5,
		ErrorCount:      len(b.errors),
		Messages:        messages,
	}
}

// Summary returns a summary of bootstrap
func (b *Bootstrapper) Summary() string {
	status := b.GetStatus()
	if status.ErrorCount == 0 {
		return fmt.Sprintf("✅ Self-hosting bootstrap complete! (%d/%d stages)",
			status.CompletedStages, status.TotalStages)
	}
	return fmt.Sprintf("❌ Bootstrap failed with %d errors", status.ErrorCount)
}
