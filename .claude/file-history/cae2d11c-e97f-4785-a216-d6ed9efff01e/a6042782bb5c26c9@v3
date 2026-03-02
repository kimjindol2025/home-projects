// Execution Context
// Manages scope, variables, and call stack

use std::collections::HashMap;
use crate::core::Value;
use std::time::Instant;

/// Call frame for function calls
#[derive(Clone, Debug)]
pub struct CallFrame {
    pub name: String,
    pub locals: HashMap<String, Value>,
    pub created_at: Instant,
}

impl CallFrame {
    /// Create new call frame
    pub fn new(name: String) -> Self {
        CallFrame {
            name,
            locals: HashMap::new(),
            created_at: Instant::now(),
        }
    }

    /// Get variable
    pub fn get(&self, name: &str) -> Option<Value> {
        self.locals.get(name).cloned()
    }

    /// Set variable
    pub fn set(&mut self, name: String, value: Value) {
        self.locals.insert(name, value);
    }
}

/// Execution context for FreeLang programs
pub struct ExecutionContext {
    globals: HashMap<String, Value>,
    call_stack: Vec<CallFrame>,
    variable_count: u64,
}

impl ExecutionContext {
    /// Create new context
    pub fn new() -> Self {
        ExecutionContext {
            globals: HashMap::new(),
            call_stack: vec![CallFrame::new("global".to_string())],
            variable_count: 0,
        }
    }

    /// Set global variable
    pub fn set_global(&mut self, name: String, value: Value) {
        self.globals.insert(name, value);
        self.variable_count += 1;
    }

    /// Get global variable
    pub fn get_global(&self, name: &str) -> Option<Value> {
        self.globals.get(name).cloned()
    }

    /// Push new call frame
    pub fn push_frame(&mut self, name: String) {
        self.call_stack.push(CallFrame::new(name));
    }

    /// Pop call frame
    pub fn pop_frame(&mut self) -> Option<CallFrame> {
        if self.call_stack.len() > 1 {
            self.call_stack.pop()
        } else {
            None
        }
    }

    /// Get current frame
    pub fn current_frame(&mut self) -> Option<&mut CallFrame> {
        self.call_stack.last_mut()
    }

    /// Get current frame (immutable)
    pub fn current_frame_ref(&self) -> Option<&CallFrame> {
        self.call_stack.last()
    }

    /// Set local variable in current scope
    pub fn set_local(&mut self, name: String, value: Value) {
        if let Some(frame) = self.call_stack.last_mut() {
            frame.set(name, value);
            self.variable_count += 1;
        }
    }

    /// Get local variable from current scope
    pub fn get_local(&self, name: &str) -> Option<Value> {
        if let Some(frame) = self.call_stack.last() {
            frame.get(name)
        } else {
            None
        }
    }

    /// Get variable (local first, then global)
    pub fn get(&self, name: &str) -> Option<Value> {
        self.get_local(name).or_else(|| self.get_global(name))
    }

    /// Set variable (in current scope or global if not found locally)
    pub fn set(&mut self, name: String, value: Value) {
        if let Some(frame) = self.call_stack.last_mut() {
            frame.set(name, value);
            self.variable_count += 1;
        }
    }

    /// Get call stack depth
    pub fn depth(&self) -> usize {
        self.call_stack.len()
    }

    /// Get frame names (for debugging)
    pub fn frame_names(&self) -> Vec<String> {
        self.call_stack.iter().map(|f| f.name.clone()).collect()
    }

    /// Get total variable count
    pub fn variable_count(&self) -> u64 {
        self.variable_count
    }

    /// Clear context
    pub fn clear(&mut self) {
        self.globals.clear();
        self.call_stack = vec![CallFrame::new("global".to_string())];
        self.variable_count = 0;
    }

    /// Get global variable count
    pub fn global_count(&self) -> usize {
        self.globals.len()
    }

    /// Get local variable count (current frame)
    pub fn local_count(&self) -> usize {
        self.call_stack.last().map(|f| f.locals.len()).unwrap_or(0)
    }
}

impl Default for ExecutionContext {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_global_vars() {
        let mut ctx = ExecutionContext::new();
        ctx.set_global("x".to_string(), Value::Number(42.0));
        assert_eq!(ctx.get_global("x").unwrap().to_number(), 42.0);
    }

    #[test]
    fn test_local_vars() {
        let mut ctx = ExecutionContext::new();
        ctx.set_local("x".to_string(), Value::Number(42.0));
        assert_eq!(ctx.get_local("x").unwrap().to_number(), 42.0);
    }

    #[test]
    fn test_call_frame() {
        let mut ctx = ExecutionContext::new();
        assert_eq!(ctx.depth(), 1);

        ctx.push_frame("foo".to_string());
        assert_eq!(ctx.depth(), 2);

        ctx.pop_frame();
        assert_eq!(ctx.depth(), 1);
    }

    #[test]
    fn test_frame_names() {
        let mut ctx = ExecutionContext::new();
        ctx.push_frame("foo".to_string());
        ctx.push_frame("bar".to_string());

        let names = ctx.frame_names();
        assert_eq!(names.len(), 3);
        assert_eq!(names[0], "global");
        assert_eq!(names[1], "foo");
        assert_eq!(names[2], "bar");
    }

    #[test]
    fn test_scoped_variables() {
        let mut ctx = ExecutionContext::new();
        ctx.set_local("x".to_string(), Value::Number(1.0));

        ctx.push_frame("foo".to_string());
        ctx.set_local("x".to_string(), Value::Number(2.0));
        assert_eq!(ctx.get_local("x").unwrap().to_number(), 2.0);

        ctx.pop_frame();
        assert_eq!(ctx.get_local("x").unwrap().to_number(), 1.0);
    }

    #[test]
    fn test_variable_shadowing() {
        let mut ctx = ExecutionContext::new();
        ctx.set_global("x".to_string(), Value::Number(1.0));

        ctx.push_frame("func".to_string());
        ctx.set_local("x".to_string(), Value::Number(2.0));

        // Local should shadow global
        assert_eq!(ctx.get("x").unwrap().to_number(), 2.0);

        ctx.pop_frame();
        // Back to global
        assert_eq!(ctx.get("x").unwrap().to_number(), 1.0);
    }

    #[test]
    fn test_variable_count() {
        let mut ctx = ExecutionContext::new();
        ctx.set_global("a".to_string(), Value::Number(1.0));
        ctx.set_global("b".to_string(), Value::Number(2.0));

        assert_eq!(ctx.global_count(), 2);

        ctx.set_local("c".to_string(), Value::Number(3.0));
        assert_eq!(ctx.local_count(), 1);
    }

    #[test]
    fn test_clear() {
        let mut ctx = ExecutionContext::new();
        ctx.set_global("x".to_string(), Value::Number(1.0));
        ctx.push_frame("func".to_string());
        ctx.set_local("y".to_string(), Value::Number(2.0));

        ctx.clear();

        assert_eq!(ctx.depth(), 1);
        assert_eq!(ctx.global_count(), 0);
        assert_eq!(ctx.local_count(), 0);
    }

    #[test]
    fn test_nested_frames() {
        let mut ctx = ExecutionContext::new();
        for i in 0..5 {
            ctx.push_frame(format!("frame_{}", i));
            ctx.set_local(
                format!("var_{}", i),
                Value::Number(i as f64),
            );
        }

        assert_eq!(ctx.depth(), 6);

        for _ in 0..5 {
            ctx.pop_frame();
        }

        assert_eq!(ctx.depth(), 1);
    }
}
