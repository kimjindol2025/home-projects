package compiler

import (
	"fmt"
	"sync"
)

// ParallelCompiler enables multi-threaded compilation
type ParallelCompiler struct {
	program         *IRProgram
	workerCount     int
	dependencyGraph *DependencyGraph
	stats           ParallelCompilationStats
	mu              sync.Mutex
}

// ParallelCompilationStats tracks parallel compilation metrics
type ParallelCompilationStats struct {
	TotalFunctions     int64
	ParallelizedFuncs  int64
	SerializedFuncs    int64
	CompilationTime    int64 // nanoseconds
	Speedup            float64
	ParallelEfficiency float64 // 0-1 scale
	CriticalPath       int64   // longest dependency chain
}

// DependencyGraph represents function call dependencies
type DependencyGraph struct {
	nodes map[string]*FunctionNode
	mu    sync.RWMutex
}

// FunctionNode represents a function in the dependency graph
type FunctionNode struct {
	Name         string
	Dependencies []*FunctionNode
	Dependents   []*FunctionNode
	Level        int // Topological level
	Compiled     bool
	mu           sync.Mutex
}

// WorkerPool manages compilation workers
type WorkerPool struct {
	workers    int
	taskQueue  chan CompilationTask
	resultChan chan CompilationResult
	wg         sync.WaitGroup
	done       chan struct{}
}

// CompilationTask represents a single compilation job
type CompilationTask struct {
	FunctionName string
	Instructions []*Instruction
	FunctionID   int
}

// CompilationResult represents completion of a compilation task
type CompilationResult struct {
	FunctionName string
	Success      bool
	Result       interface{}
	Error        error
}

// NewParallelCompiler creates a new parallel compiler
func NewParallelCompiler(program *IRProgram, workerCount int) *ParallelCompiler {
	if workerCount < 1 {
		workerCount = 4 // Default to 4 workers
	}

	return &ParallelCompiler{
		program:         program,
		workerCount:     workerCount,
		dependencyGraph: NewDependencyGraph(),
		stats:           ParallelCompilationStats{},
	}
}

// NewDependencyGraph creates an empty dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string]*FunctionNode),
	}
}

// CompileInParallel compiles functions in parallel respecting dependencies
func (pc *ParallelCompiler) CompileInParallel() error {
	if pc.program == nil {
		return fmt.Errorf("program is nil")
	}

	// Build dependency graph
	pc.buildDependencyGraph()

	// Calculate topological order
	pc.calculateTopologicalOrder()

	// Create worker pool
	pool := NewWorkerPool(pc.workerCount)
	defer pool.Stop()

	// Submit compilation tasks
	pc.submitCompilationTasks(pool)

	// Collect results
	pc.collectResults(pool)

	return nil
}

// buildDependencyGraph analyzes function dependencies
func (pc *ParallelCompiler) buildDependencyGraph() {
	pc.stats.TotalFunctions = int64(len(pc.program.Functions))

	// Create nodes for all functions
	for name := range pc.program.Functions {
		pc.dependencyGraph.AddNode(name)
	}

	// Analyze dependencies (simplified: scan for function calls)
	for name, fn := range pc.program.Functions {
		// Check instructions for call operations
		if fn.Instructions != nil {
			for _, ins := range fn.Instructions {
				if ins.Opcode == OpCall && len(ins.Args) > 0 {
					if callName, ok := ins.Args[0].(string); ok {
						pc.dependencyGraph.AddDependency(name, callName)
					}
				}
			}
		}
	}
}

// calculateTopologicalOrder orders functions for parallel execution
func (pc *ParallelCompiler) calculateTopologicalOrder() {
	pc.dependencyGraph.CalculateTopologicalOrder()

	// Calculate critical path length
	maxPath := 0
	for _, node := range pc.dependencyGraph.nodes {
		if node.Level > maxPath {
			maxPath = node.Level
		}
	}
	pc.stats.CriticalPath = int64(maxPath)

	// Identify parallelizable functions
	for _, node := range pc.dependencyGraph.nodes {
		if len(node.Dependencies) == 0 {
			pc.stats.ParallelizedFuncs++
		} else {
			pc.stats.SerializedFuncs++
		}
	}
}

// submitCompilationTasks sends functions to worker pool
func (pc *ParallelCompiler) submitCompilationTasks(pool *WorkerPool) {
	pc.dependencyGraph.mu.RLock()
	defer pc.dependencyGraph.mu.RUnlock()

	for name, fn := range pc.program.Functions {
		// Submit based on topological order
		if pc.canCompileNow(name) {
			task := CompilationTask{
				FunctionName: name,
				Instructions: fn.Instructions,
				FunctionID:   len(pc.program.Functions),
			}
			pool.Submit(task)
		}
	}
}

// canCompileNow checks if all dependencies are satisfied
func (pc *ParallelCompiler) canCompileNow(funcName string) bool {
	node, exists := pc.dependencyGraph.nodes[funcName]
	if !exists {
		return false
	}

	// Check if all dependencies are compiled
	for _, dep := range node.Dependencies {
		if !dep.Compiled {
			return false
		}
	}

	return true
}

// collectResults processes compilation results
func (pc *ParallelCompiler) collectResults(pool *WorkerPool) {
	for result := range pool.resultChan {
		pc.mu.Lock()
		if result.Success {
			// Mark function as compiled in dependency graph
			if node, exists := pc.dependencyGraph.nodes[result.FunctionName]; exists {
				node.mu.Lock()
				node.Compiled = true
				node.mu.Unlock()
			}
		}
		pc.mu.Unlock()

		// Check if more tasks can be submitted
		pc.submitPendingTasks(pool)
	}
}

// submitPendingTasks submits newly available tasks
func (pc *ParallelCompiler) submitPendingTasks(pool *WorkerPool) {
	pc.dependencyGraph.mu.RLock()
	defer pc.dependencyGraph.mu.RUnlock()

	for name, node := range pc.dependencyGraph.nodes {
		if !node.Compiled && pc.canCompileNow(name) {
			task := CompilationTask{
				FunctionName: name,
				Instructions: pc.program.Functions[name].Instructions,
				FunctionID:   len(pc.program.Functions),
			}
			pool.Submit(task)
		}
	}
}

// GetStatistics returns parallel compilation statistics
func (pc *ParallelCompiler) GetStatistics() ParallelCompilationStats {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Calculate efficiency
	if pc.stats.ParallelizedFuncs > 0 {
		pc.stats.ParallelEfficiency = float64(pc.stats.ParallelizedFuncs) /
			float64(pc.stats.TotalFunctions)
	}

	// Calculate speedup (theoretical)
	// Speedup = sequential time / parallel time
	// Simplified: workers / critical path
	if pc.stats.CriticalPath > 0 {
		pc.stats.Speedup = float64(pc.workerCount) / float64(pc.stats.CriticalPath)
	}

	return pc.stats
}

// Profile returns a profile of parallel compilation
func (pc *ParallelCompiler) Profile() string {
	stats := pc.GetStatistics()

	return fmt.Sprintf(`Parallel Compilation Profile:
  Total Functions: %d
  Parallelized: %d
  Serialized: %d
  Worker Threads: %d
  Speedup: %.2fx
  Efficiency: %.2f%%
  Critical Path: %d
`,
		stats.TotalFunctions,
		stats.ParallelizedFuncs,
		stats.SerializedFuncs,
		pc.workerCount,
		stats.Speedup,
		stats.ParallelEfficiency*100,
		stats.CriticalPath)
}

// DependencyGraph methods

// AddNode adds a function node to the graph
func (dg *DependencyGraph) AddNode(funcName string) {
	dg.mu.Lock()
	defer dg.mu.Unlock()

	if _, exists := dg.nodes[funcName]; !exists {
		dg.nodes[funcName] = &FunctionNode{
			Name:         funcName,
			Dependencies: make([]*FunctionNode, 0),
			Dependents:   make([]*FunctionNode, 0),
		}
	}
}

// AddDependency adds a dependency between functions
func (dg *DependencyGraph) AddDependency(dependent, dependency string) {
	dg.mu.Lock()
	defer dg.mu.Unlock()

	depNode, depExists := dg.nodes[dependent]
	indepNode, indepExists := dg.nodes[dependency]

	if depExists && indepExists {
		// Add edge: dependent depends on dependency
		depNode.Dependencies = append(depNode.Dependencies, indepNode)
		indepNode.Dependents = append(indepNode.Dependents, depNode)
	}
}

// CalculateTopologicalOrder assigns levels to nodes
func (dg *DependencyGraph) CalculateTopologicalOrder() {
	dg.mu.Lock()
	defer dg.mu.Unlock()

	// Simple topological sort using BFS
	visited := make(map[string]bool)
	var queue []*FunctionNode

	// Start with nodes that have no dependencies
	for _, node := range dg.nodes {
		if len(node.Dependencies) == 0 {
			queue = append(queue, node)
			node.Level = 0
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if visited[node.Name] {
			continue
		}
		visited[node.Name] = true

		// Process dependents
		for _, dependent := range node.Dependents {
			if !visited[dependent.Name] {
				if dependent.Level < node.Level+1 {
					dependent.Level = node.Level + 1
				}
				queue = append(queue, dependent)
			}
		}
	}
}

// WorkerPool methods

// NewWorkerPool creates a worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		workers:    workerCount,
		taskQueue:  make(chan CompilationTask, workerCount*2),
		resultChan: make(chan CompilationResult, workerCount),
		done:       make(chan struct{}),
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

// Submit sends a task to the worker pool
func (wp *WorkerPool) Submit(task CompilationTask) {
	wp.taskQueue <- task
}

// worker processes compilation tasks
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}
			// Simulate compilation
			result := CompilationResult{
				FunctionName: task.FunctionName,
				Success:      true,
				Result:       nil,
			}
			wp.resultChan <- result

		case <-wp.done:
			return
		}
	}
}

// Stop shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	wp.wg.Wait()
	close(wp.resultChan)
	close(wp.done)
}

// GetWorkerCount returns the number of workers
func (wp *WorkerPool) GetWorkerCount() int {
	return wp.workers
}
