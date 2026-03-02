// Call Stack
// Manages function call contexts and variable scoping

use crate::core::Value;
use std::collections::HashMap;

/// A single frame in the call stack
#[derive(Clone, Debug)]
pub struct StackFrame {
    pub function_name: String,
    pub locals: HashMap<String, Value>,
    pub depth: usize,
}

impl StackFrame {
    /// Create new stack frame
    pub fn new(function_name: String, depth: usize) -> Self {
        StackFrame {
            function_name,
            locals: HashMap::new(),
            depth,
        }
    }

    /// Set local variable
    pub fn set_local(&mut self, name: String, value: Value) {
        self.locals.insert(name, value);
    }

    /// Get local variable
    pub fn get_local(&self, name: &str) -> Option<Value> {
        self.locals.get(name).cloned()
    }

    /// Get all local variables
    pub fn locals(&self) -> &HashMap<String, Value> {
        &self.locals
    }

    /// Variable count in frame
    pub fn var_count(&self) -> usize {
        self.locals.len()
    }
}

/// Call stack for managing function execution
pub struct CallStack {
    frames: Vec<StackFrame>,
    max_depth: usize,
    call_count: u64,
}

impl CallStack {
    /// Create new call stack
    pub fn new() -> Self {
        let mut stack = CallStack {
            frames: Vec::new(),
            max_depth: 1000, // Prevent stack overflow
            call_count: 0,
        };
        // Push global frame
        stack.frames.push(StackFrame::new("global".to_string(), 0));
        stack
    }

    /// Push a new frame
    pub fn push(&mut self, function_name: String) -> Result<(), String> {
        let depth = self.frames.len();

        if depth >= self.max_depth {
            return Err(format!("Maximum call stack depth ({}) exceeded", self.max_depth));
        }

        self.frames.push(StackFrame::new(function_name, depth));
        self.call_count += 1;
        Ok(())
    }

    /// Pop a frame
    pub fn pop(&mut self) -> Option<StackFrame> {
        if self.frames.len() > 1 {
            self.frames.pop()
        } else {
            None
        }
    }

    /// Get current frame (mutable)
    pub fn current_mut(&mut self) -> Option<&mut StackFrame> {
        self.frames.last_mut()
    }

    /// Get current frame (immutable)
    pub fn current(&self) -> Option<&StackFrame> {
        self.frames.last()
    }

    /// Set variable in current frame
    pub fn set_local(&mut self, name: String, value: Value) -> Result<(), String> {
        if let Some(frame) = self.frames.last_mut() {
            frame.set_local(name, value);
            Ok(())
        } else {
            Err("No frame in call stack".to_string())
        }
    }

    /// Get variable from current frame
    pub fn get_local(&self, name: &str) -> Option<Value> {
        if let Some(frame) = self.frames.last() {
            frame.get_local(name)
        } else {
            None
        }
    }

    /// Get variable from any frame in the stack (search from current to global)
    pub fn get_variable(&self, name: &str) -> Option<Value> {
        for frame in self.frames.iter().rev() {
            if let Some(val) = frame.get_local(name) {
                return Some(val);
            }
        }
        None
    }

    /// Get call stack depth
    pub fn depth(&self) -> usize {
        self.frames.len()
    }

    /// Get frame names
    pub fn frame_names(&self) -> Vec<&str> {
        self.frames.iter().map(|f| f.function_name.as_str()).collect()
    }

    /// Get total variables in stack
    pub fn total_variables(&self) -> usize {
        self.frames.iter().map(|f| f.var_count()).sum()
    }

    /// Get call count
    pub fn call_count(&self) -> u64 {
        self.call_count
    }

    /// Clear the stack (reset to global frame)
    pub fn clear(&mut self) {
        self.frames.clear();
        self.frames.push(StackFrame::new("global".to_string(), 0));
        self.call_count = 0;
    }

    /// Get stack trace
    pub fn get_trace(&self) -> String {
        let mut trace = String::from("Call Stack Trace:\n");
        for (idx, frame) in self.frames.iter().enumerate() {
            trace.push_str(&format!("  [{}] {}\n", idx, frame.function_name));
        }
        trace
    }

    /// Check if we're in global scope
    pub fn is_global(&self) -> bool {
        self.frames.len() == 1
    }
}

impl Default for CallStack {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_callstack_creation() {
        let stack = CallStack::new();
        assert_eq!(stack.depth(), 1); // Global frame
        assert_eq!(stack.is_global(), true);
    }

    #[test]
    fn test_push_pop_frames() {
        let mut stack = CallStack::new();
        assert_eq!(stack.depth(), 1);

        stack.push("func1".to_string()).unwrap();
        assert_eq!(stack.depth(), 2);
        assert!(!stack.is_global());

        stack.push("func2".to_string()).unwrap();
        assert_eq!(stack.depth(), 3);

        stack.pop();
        assert_eq!(stack.depth(), 2);

        stack.pop();
        assert_eq!(stack.depth(), 1);
    }

    #[test]
    fn test_set_get_local() {
        let mut stack = CallStack::new();
        stack.push("test".to_string()).unwrap();

        stack.set_local("x".to_string(), Value::Number(42.0)).ok();
        let val = stack.get_local("x");
        assert!(val.is_some());
        assert_eq!(val.unwrap().to_number(), 42.0);
    }

    #[test]
    fn test_variable_scoping() {
        let mut stack = CallStack::new();

        // Global scope
        stack.set_local("global_var".to_string(), Value::Number(1.0)).ok();

        // Function scope 1
        stack.push("func1".to_string()).unwrap();
        stack.set_local("func1_var".to_string(), Value::Number(2.0)).ok();

        // Function scope 2
        stack.push("func2".to_string()).unwrap();
        stack.set_local("func2_var".to_string(), Value::Number(3.0)).ok();

        // Should see func2_var
        assert_eq!(stack.get_local("func2_var").unwrap().to_number(), 3.0);

        // Search up the stack for func1_var
        assert_eq!(stack.get_variable("func1_var").unwrap().to_number(), 2.0);

        // Pop and check we're back in func1
        stack.pop();
        assert_eq!(stack.get_local("func1_var").unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_total_variables() {
        let mut stack = CallStack::new();

        stack.set_local("a".to_string(), Value::Number(1.0)).ok();
        assert_eq!(stack.total_variables(), 1);

        stack.push("func".to_string()).unwrap();
        stack.set_local("b".to_string(), Value::Number(2.0)).ok();
        stack.set_local("c".to_string(), Value::Number(3.0)).ok();

        assert_eq!(stack.total_variables(), 3);

        stack.pop();
        assert_eq!(stack.total_variables(), 1);
    }

    #[test]
    fn test_frame_names() {
        let mut stack = CallStack::new();
        stack.push("func1".to_string()).unwrap();
        stack.push("func2".to_string()).unwrap();

        let names = stack.frame_names();
        assert_eq!(names.len(), 3);
        assert_eq!(names[0], "global");
        assert_eq!(names[1], "func1");
        assert_eq!(names[2], "func2");
    }

    #[test]
    fn test_max_depth_exceeded() {
        let mut stack = CallStack::new();
        stack.max_depth = 5;

        for i in 0..4 {
            assert!(stack.push(format!("func{}", i)).is_ok());
        }

        // Next push should fail
        assert!(stack.push("func_overflow".to_string()).is_err());
    }

    #[test]
    fn test_call_count() {
        let mut stack = CallStack::new();
        assert_eq!(stack.call_count(), 0);

        stack.push("func1".to_string()).unwrap();
        assert_eq!(stack.call_count(), 1);

        stack.push("func2".to_string()).unwrap();
        assert_eq!(stack.call_count(), 2);

        stack.pop();
        assert_eq!(stack.call_count(), 2); // Call count doesn't decrease
    }

    #[test]
    fn test_stack_trace() {
        let mut stack = CallStack::new();
        stack.push("main".to_string()).unwrap();
        stack.push("helper".to_string()).unwrap();

        let trace = stack.get_trace();
        assert!(trace.contains("global"));
        assert!(trace.contains("main"));
        assert!(trace.contains("helper"));
    }

    #[test]
    fn test_clear() {
        let mut stack = CallStack::new();
        stack.push("func".to_string()).unwrap();
        stack.set_local("x".to_string(), Value::Number(42.0)).ok();

        assert_eq!(stack.depth(), 2);
        assert_eq!(stack.total_variables(), 1);

        stack.clear();

        assert_eq!(stack.depth(), 1);
        assert_eq!(stack.total_variables(), 0);
        assert_eq!(stack.call_count(), 0);
    }
}
