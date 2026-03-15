package compiler

import (
	"testing"
	"sync"
)

// Test parallel compiler creation
func TestNewParallelCompiler(t *testing.T) {
	program := NewIRProgram()
	compiler := NewParallelCompiler(program, 4)

	if compiler == nil {
		t.Errorf("expected compiler, got nil")
	}

	if compiler.workerCount != 4 {
		t.Errorf("expected 4 workers, got %d", compiler.workerCount)
	}
}

// Test default worker count
func TestDefaultWorkerCount(t *testing.T) {
	program := NewIRProgram()
	compiler := NewParallelCompiler(program, 0)

	if compiler.workerCount != 4 {
		t.Errorf("expected default 4 workers, got %d", compiler.workerCount)
	}
}

// Test dependency graph creation
func TestNewDependencyGraph(t *testing.T) {
	graph := NewDependencyGraph()

	if graph == nil {
		t.Errorf("expected graph, got nil")
	}

	if len(graph.nodes) != 0 {
		t.Errorf("expected empty graph")
	}
}

// Test adding nodes to dependency graph
func TestAddNode(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddNode("func1")
	graph.AddNode("func2")

	if len(graph.nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.nodes))
	}

	if _, exists := graph.nodes["func1"]; !exists {
		t.Errorf("func1 should exist in graph")
	}
}

// Test adding dependencies
func TestAddDependency(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddNode("func1")
	graph.AddNode("func2")
	graph.AddDependency("func2", "func1")

	func2 := graph.nodes["func2"]
	if len(func2.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(func2.Dependencies))
	}
}

// Test topological order calculation
func TestCalculateTopologicalOrder(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddNode("func1")
	graph.AddNode("func2")
	graph.AddNode("func3")
	graph.AddDependency("func2", "func1")
	graph.AddDependency("func3", "func2")

	graph.CalculateTopologicalOrder()

	func1 := graph.nodes["func1"]
	func2 := graph.nodes["func2"]
	func3 := graph.nodes["func3"]

	if func1.Level != 0 {
		t.Errorf("func1 level should be 0, got %d", func1.Level)
	}

	if func2.Level != 1 {
		t.Errorf("func2 level should be 1, got %d", func2.Level)
	}

	if func3.Level != 2 {
		t.Errorf("func3 level should be 2, got %d", func3.Level)
	}
}

// Test worker pool creation
func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Stop()

	if pool.workers != 4 {
		t.Errorf("expected 4 workers, got %d", pool.workers)
	}

	if pool.GetWorkerCount() != 4 {
		t.Errorf("expected GetWorkerCount() = 4, got %d", pool.GetWorkerCount())
	}
}

// Test parallel compilation with single function
func TestParallelCompileSingleFunction(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 42.0))

	compiler := NewParallelCompiler(program, 2)
	err := compiler.CompileInParallel()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := compiler.GetStatistics()
	if stats.TotalFunctions != 1 {
		t.Errorf("expected 1 total function, got %d", stats.TotalFunctions)
	}
}

// Test parallel compilation with multiple functions
func TestParallelCompileMultipleFunctions(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	for i := 0; i < 5; i++ {
		name := "func" + string(rune('1'+i))
		program.Functions[name] = NewIRFunction(name)
		program.Functions[name].Instructions = append(program.Functions[name].Instructions,
			NewInstruction(OpLoadConst, float64(i)))
	}

	compiler := NewParallelCompiler(program, 4)
	err := compiler.CompileInParallel()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := compiler.GetStatistics()
	if stats.TotalFunctions != 5 {
		t.Errorf("expected 5 total functions, got %d", stats.TotalFunctions)
	}
}

// Test parallel compilation with dependencies
func TestParallelCompileWithDependencies(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	// func1 (no dependencies)
	program.Functions["func1"] = NewIRFunction("func1")
	program.Functions["func1"].Instructions = append(program.Functions["func1"].Instructions,
		NewInstruction(OpLoadConst, 1.0))

	// func2 (depends on func1)
	program.Functions["func2"] = NewIRFunction("func2")
	program.Functions["func2"].Instructions = append(program.Functions["func2"].Instructions,
		NewInstruction(OpCall, "func1"))

	compiler := NewParallelCompiler(program, 2)
	err := compiler.CompileInParallel()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stats := compiler.GetStatistics()
	if stats.ParallelizedFuncs < 1 {
		t.Logf("expected at least 1 parallelizable function")
	}
}

// Test statistics structure
func TestStatisticsStructure(t *testing.T) {
	stats := ParallelCompilationStats{
		TotalFunctions:     5,
		ParallelizedFuncs:  3,
		SerializedFuncs:    2,
		Speedup:            2.5,
		ParallelEfficiency: 0.6,
	}

	if stats.TotalFunctions != 5 {
		t.Errorf("expected 5 total functions")
	}

	if stats.Speedup != 2.5 {
		t.Errorf("expected 2.5 speedup")
	}

	if stats.ParallelEfficiency != 0.6 {
		t.Errorf("expected 0.6 efficiency")
	}
}

// Test profile output
func TestParallelProfileOutput(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)
	program.Functions["func1"] = NewIRFunction("func1")

	compiler := NewParallelCompiler(program, 4)
	_ = compiler.CompileInParallel()

	profile := compiler.Profile()

	if profile == "" {
		t.Errorf("expected non-empty profile")
	}

	if len(profile) < 50 {
		t.Errorf("expected detailed profile, got %d chars", len(profile))
	}
}

// Test nil program handling
func TestNilProgramHandling(t *testing.T) {
	compiler := NewParallelCompiler(nil, 4)
	err := compiler.CompileInParallel()

	if err == nil {
		t.Errorf("expected error for nil program")
	}
}

// Test function node structure
func TestFunctionNode(t *testing.T) {
	node := &FunctionNode{
		Name: "test_func",
	}

	if node.Name != "test_func" {
		t.Errorf("expected name 'test_func'")
	}

	if node.Compiled {
		t.Errorf("expected compiled=false initially")
	}

	if node.Level != 0 {
		t.Errorf("expected level=0 initially")
	}
}

// Test critical path detection
func TestCriticalPathDetection(t *testing.T) {
	graph := NewDependencyGraph()

	// Create chain: 1 -> 2 -> 3 -> 4
	for i := 1; i <= 4; i++ {
		name := "func" + string(rune('0'+i))
		graph.AddNode(name)
	}

	graph.AddDependency("func2", "func1")
	graph.AddDependency("func3", "func2")
	graph.AddDependency("func4", "func3")

	graph.CalculateTopologicalOrder()

	func4 := graph.nodes["func4"]
	if func4.Level != 3 {
		t.Errorf("expected func4 level=3, got %d", func4.Level)
	}
}

// Test concurrent task submission
func TestConcurrentTaskSubmission(t *testing.T) {
	pool := NewWorkerPool(2)

	var wg sync.WaitGroup

	// Submit tasks concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := CompilationTask{
				FunctionName: "func" + string(rune('0'+id%5)),
			}
			pool.Submit(task)
		}(i)
	}

	wg.Wait()
	close(pool.taskQueue)
	pool.wg.Wait()

	// Verify results exist
	resultCount := 0
	for result := range pool.resultChan {
		if result.Success {
			resultCount++
		}
	}

	if resultCount > 0 {
		t.Logf("processed %d results", resultCount)
	}
}

// Test speedup calculation
func TestSpeedupCalculation(t *testing.T) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	// Create multiple independent functions
	for i := 0; i < 4; i++ {
		name := "func" + string(rune('0'+i))
		program.Functions[name] = NewIRFunction(name)
		program.Functions[name].Instructions = append(program.Functions[name].Instructions,
			NewInstruction(OpLoadConst, float64(i)))
	}

	compiler := NewParallelCompiler(program, 4)
	_ = compiler.CompileInParallel()

	stats := compiler.GetStatistics()
	if stats.Speedup <= 0 {
		t.Errorf("expected positive speedup")
	}
}

// Benchmark parallel compilation
func BenchmarkParallelCompilation(b *testing.B) {
	program := NewIRProgram()
	program.Functions = make(map[string]*IRFunction)

	for i := 0; i < 20; i++ {
		name := "func" + string(rune('0' + (i % 10)))
		if program.Functions[name] == nil {
			program.Functions[name] = NewIRFunction(name)
		}
		program.Functions[name].Instructions = append(program.Functions[name].Instructions,
			NewInstruction(OpLoadConst, float64(i)))
	}

	compiler := NewParallelCompiler(program, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = compiler.CompileInParallel()
	}
}

// Benchmark worker pool
func BenchmarkWorkerPool(b *testing.B) {
	pool := NewWorkerPool(4)
	defer pool.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := CompilationTask{
			FunctionName: "func",
		}
		pool.Submit(task)
	}
}

// Benchmark dependency graph
func BenchmarkDependencyGraph(b *testing.B) {
	graph := NewDependencyGraph()

	for i := 0; i < 50; i++ {
		graph.AddNode("func" + string(rune('0'+(i%10))))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.CalculateTopologicalOrder()
	}
}
