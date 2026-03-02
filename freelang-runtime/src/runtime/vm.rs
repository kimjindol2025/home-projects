// Runtime Virtual Machine
// Main execution engine

use crate::core::Value;
use crate::memory::{MemoryAllocator, GarbageCollector, ReferenceCounter, MemoryStats};
use crate::runtime::context::ExecutionContext;
use crate::runtime::executor::Executor;
use crate::stdlib::FunctionRegistry;
use std::time::Instant;

/// Runtime engine statistics
#[derive(Clone, Debug)]
pub struct RuntimeStats {
    pub uptime_ms: u128,
    pub memory: MemoryStats,
    pub memory_usage: usize,
}

/// FreeLang runtime engine
pub struct RuntimeEngine {
    memory: MemoryAllocator,
    gc: GarbageCollector,
    refcount: ReferenceCounter,
    context: ExecutionContext,
    executor: Executor,
    function_registry: FunctionRegistry,
    started_at: Instant,
    execution_count: u64,
}

impl RuntimeEngine {
    /// Create new runtime
    pub fn new() -> Self {
        RuntimeEngine {
            memory: MemoryAllocator::new(),
            gc: GarbageCollector::new(),
            refcount: ReferenceCounter::new(),
            context: ExecutionContext::new(),
            executor: Executor::new(),
            function_registry: FunctionRegistry::new(),
            started_at: Instant::now(),
            execution_count: 0,
        }
    }

    /// Execute a value
    pub fn execute(&mut self, value: Value) -> Result<Value, String> {
        self.execution_count += 1;

        // Check if GC should run
        if self.memory.should_gc() {
            self.run_gc();
        }

        Ok(value)
    }

    /// Execute binary operation
    pub fn execute_binop(&mut self, op: &str, left: Value, right: Value) -> Result<Value, String> {
        self.executor.binary_op(op, left, right)
    }

    /// Execute unary operation
    pub fn execute_unop(&mut self, op: &str, operand: Value) -> Result<Value, String> {
        self.executor.unary_op(op, operand)
    }

    /// Call a function
    pub fn call_function(&mut self, func: Value, args: Vec<Value>) -> Result<Value, String> {
        self.executor.call_function(func, args)
    }

    /// Call a stdlib function by name
    pub fn call_stdlib(&self, name: &str, args: Vec<Value>) -> Result<Value, String> {
        self.function_registry.call(name, args)
    }

    /// Check if function exists
    pub fn function_exists(&self, name: &str) -> bool {
        self.function_registry.exists(name)
    }

    /// List all available functions
    pub fn list_functions(&self) -> Vec<&str> {
        self.function_registry.list_functions()
    }

    /// Get function count
    pub fn function_count(&self) -> usize {
        self.function_registry.function_count()
    }

    /// Run garbage collection
    fn run_gc(&mut self) {
        // Mark phase
        self.gc.mark(&[]);

        // Sweep phase
        self.gc.sweep();
    }

    /// Get memory statistics
    pub fn memory_stats(&self) -> usize {
        self.memory.memory_usage()
    }

    /// Get runtime statistics
    pub fn stats(&self) -> RuntimeStats {
        RuntimeStats {
            uptime_ms: self.started_at.elapsed().as_millis(),
            memory: MemoryStats {
                allocated: 0,
                freed: 0,
                active: self.memory.object_count(),
                gc_runs: 0,
                gc_collected: 0,
            },
            memory_usage: self.memory.memory_usage(),
        }
    }

    /// Get execution count
    pub fn execution_count(&self) -> u64 {
        self.execution_count
    }

    /// Get variable count
    pub fn variable_count(&self) -> u64 {
        self.context.variable_count()
    }

    /// Reset runtime
    pub fn reset(&mut self) {
        self.context.clear();
        self.memory.reset();
        self.gc.reset();
        self.refcount.clear();
        self.executor.reset_stats();
        self.execution_count = 0;
        self.started_at = Instant::now();
    }
}

impl Default for RuntimeEngine {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_runtime_creation() {
        let engine = RuntimeEngine::new();
        assert_eq!(engine.memory_stats(), 0);
        assert_eq!(engine.execution_count(), 0);
    }

    #[test]
    fn test_execute() {
        let mut engine = RuntimeEngine::new();
        let result = engine.execute(Value::Number(42.0));
        assert!(result.is_ok());
        assert_eq!(engine.execution_count(), 1);
    }

    #[test]
    fn test_execute_binop() {
        let mut engine = RuntimeEngine::new();
        let result = engine.execute_binop("+", Value::Number(5.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 8.0);
    }

    #[test]
    fn test_execute_unop() {
        let mut engine = RuntimeEngine::new();
        let result = engine.execute_unop("-", Value::Number(5.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), -5.0);
    }

    #[test]
    fn test_variable_tracking() {
        let mut engine = RuntimeEngine::new();
        engine.context.set_global("x".to_string(), Value::Number(42.0));
        engine.context.set_local("y".to_string(), Value::Number(3.14));

        assert_eq!(engine.variable_count(), 2);
    }

    #[test]
    fn test_stats() {
        let mut engine = RuntimeEngine::new();
        engine.execute(Value::Number(1.0)).ok();
        engine.execute(Value::Number(2.0)).ok();

        let stats = engine.stats();
        assert!(stats.uptime_ms > 0);
    }

    #[test]
    fn test_reset() {
        let mut engine = RuntimeEngine::new();
        engine.context.set_global("x".to_string(), Value::Number(42.0));
        engine.execute(Value::Number(1.0)).ok();

        assert_eq!(engine.variable_count(), 1);
        assert_eq!(engine.execution_count(), 1);

        engine.reset();

        assert_eq!(engine.variable_count(), 0);
        assert_eq!(engine.execution_count(), 0);
    }
}
