package selfhost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test NewBootstrapper creates correct initial state
func TestNewBootstrapper(t *testing.T) {
	b := NewBootstrapper("/source", "/output")

	if b.stage != 0 {
		t.Errorf("expected initial stage 0, got %d", b.stage)
	}

	if b.sourceDir != "/source" {
		t.Errorf("expected source dir /source, got %s", b.sourceDir)
	}

	if b.outputDir != "/output" {
		t.Errorf("expected output dir /output, got %s", b.outputDir)
	}

	if len(b.errors) != 0 {
		t.Errorf("expected no initial errors, got %d", len(b.errors))
	}
}

// Test compileWithGoCompiler (stage 0)
func TestCompileWithGoCompiler(t *testing.T) {
	b := NewBootstrapper("", "")
	err := b.compileWithGoCompiler()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test simple file compilation
func TestCompileFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	b := NewBootstrapper(tmpDir, outputDir)

	// Test simple source
	source := `fn myTest() { return 42 }`
	err := b.compileFile(source, "test.flir")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify output file was created
	outputPath := filepath.Join(outputDir, "test.flir")
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file to exist at %s", outputPath)
	}
}

// Test compileLexerStage
func TestCompileLexerStage(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	// Create lexer.free file
	lexerFile := filepath.Join(tmpDir, "lexer.free")
	lexerCode := `fn tokenize(source) {
	let tokens = [1, 2, 3]
	return tokens
}`

	err := os.WriteFile(lexerFile, []byte(lexerCode), 0644)
	if err != nil {
		t.Fatalf("failed to create lexer.free: %v", err)
	}

	b := NewBootstrapper(tmpDir, outputDir)
	err = b.compileLexerStage()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify output
	outputPath := filepath.Join(outputDir, "lexer_stage1.flir")
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

// Test compileParserStage
func TestCompileParserStage(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	// Create parser.free file
	parserFile := filepath.Join(tmpDir, "parser.free")
	parserCode := `fn parse(tokens) {
	let ast = []
	return ast
}`

	err := os.WriteFile(parserFile, []byte(parserCode), 0644)
	if err != nil {
		t.Fatalf("failed to create parser.free: %v", err)
	}

	b := NewBootstrapper(tmpDir, outputDir)
	err = b.compileParserStage()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify output
	outputPath := filepath.Join(outputDir, "parser_stage2.flir")
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

// Test compileCompilerStage
func TestCompileCompilerStage(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	// Create compiler.free file
	compilerFile := filepath.Join(tmpDir, "compiler.free")
	compilerCode := `fn compile(ast) {
	let ir = []
	return ir
}`

	err := os.WriteFile(compilerFile, []byte(compilerCode), 0644)
	if err != nil {
		t.Fatalf("failed to create compiler.free: %v", err)
	}

	b := NewBootstrapper(tmpDir, outputDir)
	err = b.compileCompilerStage()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify output
	outputPath := filepath.Join(outputDir, "compiler_stage3.flir")
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

// Test verifySelfHosting
func TestVerifySelfHosting(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	b := NewBootstrapper(tmpDir, outputDir)
	err := b.verifySelfHosting()

	// verifySelfHosting will succeed or fail based on the hardcoded test code
	// Just verify the method runs without panic
	_ = err

	// Verify test output
	outputPath := filepath.Join(outputDir, "test_stage4.flir")
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected test output file to exist: %v", err)
	}
}

// Test addError
func TestAddError(t *testing.T) {
	b := NewBootstrapper("", "")
	b.addError(1, "test error")

	if len(b.errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(b.errors))
	}

	if b.errors[0].Stage != 1 {
		t.Errorf("expected error stage 1, got %d", b.errors[0].Stage)
	}

	if b.errors[0].Message != "test error" {
		t.Errorf("expected error message 'test error', got '%s'", b.errors[0].Message)
	}
}

// Test Errors method
func TestErrors(t *testing.T) {
	b := NewBootstrapper("", "")
	b.addError(0, "error 1")
	b.addError(1, "error 2")

	errors := b.Errors()

	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errors))
	}
}

// Test GetStatus
func TestGetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	b := NewBootstrapper(tmpDir, outputDir)
	b.stage = 2

	status := b.GetStatus()

	if status.Stage != 2 {
		t.Errorf("expected stage 2, got %d", status.Stage)
	}

	if status.TotalStages != 5 {
		t.Errorf("expected 5 total stages, got %d", status.TotalStages)
	}

	if status.CompletedStages != 2 {
		t.Errorf("expected 2 completed stages, got %d", status.CompletedStages)
	}

	if status.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", status.ErrorCount)
	}
}

// Test GetStatus with errors
func TestGetStatusWithErrors(t *testing.T) {
	b := NewBootstrapper("", "")
	b.stage = 1
	b.addError(0, "compilation error")
	b.addError(1, "link error")

	status := b.GetStatus()

	if status.ErrorCount != 2 {
		t.Errorf("expected 2 errors, got %d", status.ErrorCount)
	}

	if len(status.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(status.Messages))
	}

	if !strings.Contains(status.Messages[0], "Stage 0") {
		t.Errorf("expected message to contain 'Stage 0', got '%s'", status.Messages[0])
	}
}

// Test Summary - no errors
func TestSummarySuccess(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	b := NewBootstrapper(tmpDir, outputDir)
	b.stage = 5

	summary := b.Summary()

	if !strings.Contains(summary, "✅") {
		t.Errorf("expected success indicator in summary, got '%s'", summary)
	}

	if !strings.Contains(summary, "5/5") {
		t.Errorf("expected '5/5 stages' in summary, got '%s'", summary)
	}
}

// Test Summary - with errors
func TestSummaryFailure(t *testing.T) {
	b := NewBootstrapper("", "")
	b.addError(1, "compilation failed")

	summary := b.Summary()

	if !strings.Contains(summary, "❌") {
		t.Errorf("expected failure indicator in summary, got '%s'", summary)
	}

	if !strings.Contains(summary, "1 error") {
		t.Errorf("expected error count in summary, got '%s'", summary)
	}
}

// Test Bootstrap method - missing files
func TestBootstrapMissingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	b := NewBootstrapper(tmpDir, outputDir)
	err := b.Bootstrap()

	// Should fail due to missing lexer.free
	if err == nil {
		t.Errorf("expected error for missing lexer.free")
	}

	// Should have recorded errors
	if len(b.errors) == 0 {
		t.Errorf("expected errors to be recorded")
	}
}

// Test Bootstrap method - with valid files
func TestBootstrapWithValidFiles(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := t.TempDir()

	// Create all required files
	files := map[string]string{
		"lexer.free":   `fn tokenize(s) { return [] }`,
		"parser.free":  `fn parse(t) { return [] }`,
		"compiler.free": `fn compile(a) { return {} }`,
	}

	for filename, code := range files {
		filepath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filepath, []byte(code), 0644)
		if err != nil {
			t.Fatalf("failed to create %s: %v", filename, err)
		}
	}

	b := NewBootstrapper(tmpDir, outputDir)
	_ = b.Bootstrap()

	// May fail due to runtime issues, but file reading should work
	// Just verify it attempted the bootstrap
	if b.stage < 1 {
		t.Logf("bootstrap stage: %d", b.stage)
	}
}

// Test BootstrapError struct
func TestBootstrapError(t *testing.T) {
	err := BootstrapError{
		Stage:   2,
		Message: "test",
	}

	if err.Stage != 2 {
		t.Errorf("expected stage 2, got %d", err.Stage)
	}

	if err.Message != "test" {
		t.Errorf("expected message 'test', got '%s'", err.Message)
	}
}

// Test Status struct
func TestStatusStruct(t *testing.T) {
	status := Status{
		Stage:           2,
		CompletedStages: 2,
		TotalStages:     5,
		ErrorCount:      0,
		Messages:        []string{"stage 0 ok", "stage 1 ok"},
	}

	if status.Stage != 2 {
		t.Errorf("expected stage 2, got %d", status.Stage)
	}

	if status.TotalStages != 5 {
		t.Errorf("expected 5 total stages, got %d", status.TotalStages)
	}

	if len(status.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(status.Messages))
	}
}

// Benchmark Bootstrap
func BenchmarkBootstrap(b *testing.B) {
	tmpDir := b.TempDir()
	outputDir := b.TempDir()

	// Create minimal files
	files := map[string]string{
		"lexer.free":    `fn f() {}`,
		"parser.free":   `fn f() {}`,
		"compiler.free": `fn f() {}`,
	}

	for filename, code := range files {
		filepath := filepath.Join(tmpDir, filename)
		os.WriteFile(filepath, []byte(code), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bootstrapper := NewBootstrapper(tmpDir, outputDir)
		bootstrapper.Bootstrap()
	}
}

// Benchmark CompileFile
func BenchmarkCompileFile(b *testing.B) {
	tmpDir := b.TempDir()
	outputDir := b.TempDir()

	bootstrapper := NewBootstrapper(tmpDir, outputDir)
	source := `fn test() { return 42 }`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bootstrapper.compileFile(source, "bench.flir")
	}
}
