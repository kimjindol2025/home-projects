// Function Call Mechanism Tests
// Tests for user-defined functions, closures, and call stack management

use freelang_runtime::{RuntimeEngine, Value};
use freelang_runtime::core::{UserFunction, FunctionKind};

#[test]
fn test_user_function_creation() {
    let func = UserFunction::new(
        "add".to_string(),
        vec!["x".to_string(), "y".to_string()],
        vec![],
    );

    assert_eq!(func.get_name(), "add");
    assert_eq!(func.param_count(), 2);
    assert_eq!(func.get_params(), &["x", "y"]);
}

#[test]
fn test_user_function_display() {
    let func = UserFunction::new(
        "multiply".to_string(),
        vec!["a".to_string(), "b".to_string(), "c".to_string()],
        vec![],
    );

    let display = format!("{}", func);
    assert!(display.contains("multiply"));
    assert!(display.contains("3")); // Parameter count
}

#[test]
fn test_closure_variables() {
    let func = UserFunction::new(
        "calculator".to_string(),
        vec!["x".to_string()],
        vec![],
    );

    // Add closure variables
    func.set_closure_var("base".to_string(), Value::Number(10.0));
    func.set_closure_var("factor".to_string(), Value::Number(2.0));

    // Retrieve them
    assert_eq!(func.get_closure_var("base").unwrap().to_number(), 10.0);
    assert_eq!(func.get_closure_var("factor").unwrap().to_number(), 2.0);
    assert!(func.get_closure_var("nonexistent").is_none());
}

#[test]
fn test_function_kind_native() {
    fn my_func(args: Vec<Value>) -> Result<Value, String> {
        if args.is_empty() {
            return Err("requires arguments".to_string());
        }
        Ok(Value::Number(args[0].to_number() * 2.0))
    }

    let kind = FunctionKind::Native(my_func);

    assert!(kind.is_native());
    assert!(!kind.is_user());
    assert!(!kind.is_closure());

    let result = kind.call(vec![Value::Number(5.0)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 10.0);
}

#[test]
fn test_function_kind_user() {
    let user_func = UserFunction::new(
        "process".to_string(),
        vec!["input".to_string()],
        vec![],
    );

    let kind = FunctionKind::User(user_func);

    assert!(!kind.is_native());
    assert!(kind.is_user());
    assert!(!kind.is_closure());

    assert_eq!(kind.name(), Some("process"));
    assert_eq!(kind.param_count(), Some(1));
}

#[test]
fn test_function_kind_closure() {
    let func = UserFunction::new(
        "closure".to_string(),
        vec!["x".to_string()],
        vec![],
    );

    let captured = std::collections::HashMap::new();
    let kind = FunctionKind::Closure {
        func,
        captured: std::rc::Rc::new(std::cell::RefCell::new(captured)),
    };

    assert!(!kind.is_native());
    assert!(!kind.is_user());
    assert!(kind.is_closure());

    assert_eq!(kind.name(), Some("closure"));
    assert_eq!(kind.param_count(), Some(1));
}

#[test]
fn test_callstack_basics() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();

    assert_eq!(stack.depth(), 1); // Global frame
    assert!(stack.is_global());

    // Push frames
    stack.push("function1".to_string()).unwrap();
    assert_eq!(stack.depth(), 2);
    assert!(!stack.is_global());

    stack.push("function2".to_string()).unwrap();
    assert_eq!(stack.depth(), 3);

    // Pop frames
    stack.pop();
    assert_eq!(stack.depth(), 2);

    stack.pop();
    assert_eq!(stack.depth(), 1);
    assert!(stack.is_global());
}

#[test]
fn test_callstack_variables() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();

    // Global variables
    stack.set_local("global_x".to_string(), Value::Number(1.0)).ok();

    // Function scope 1
    stack.push("func1".to_string()).unwrap();
    stack.set_local("x".to_string(), Value::Number(10.0)).ok();
    stack.set_local("y".to_string(), Value::Number(20.0)).ok();

    // Verify we can get local variables
    assert_eq!(stack.get_local("x").unwrap().to_number(), 10.0);
    assert_eq!(stack.get_local("y").unwrap().to_number(), 20.0);

    // Can search up the stack
    assert_eq!(stack.get_variable("global_x").unwrap().to_number(), 1.0);

    // Total variable count
    assert_eq!(stack.total_variables(), 3); // 1 global + 2 local
}

#[test]
fn test_callstack_variable_shadowing() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();

    // Global
    stack.set_local("x".to_string(), Value::Number(1.0)).ok();

    // Function 1 - shadows global
    stack.push("func1".to_string()).unwrap();
    stack.set_local("x".to_string(), Value::Number(10.0)).ok();

    // In func1, we see the shadowed variable
    assert_eq!(stack.get_local("x").unwrap().to_number(), 10.0);

    // Function 2 - shadows further
    stack.push("func2".to_string()).unwrap();
    stack.set_local("x".to_string(), Value::Number(100.0)).ok();

    // In func2, we see the most recent shadowing
    assert_eq!(stack.get_local("x").unwrap().to_number(), 100.0);

    // Pop back to func1
    stack.pop();
    assert_eq!(stack.get_local("x").unwrap().to_number(), 10.0);
}

#[test]
fn test_callstack_depth_limit() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();
    stack.max_depth = 5;

    // Can push up to max_depth - 1 (global frame counts as 1)
    for i in 0..4 {
        assert!(stack.push(format!("func{}", i)).is_ok());
    }

    assert_eq!(stack.depth(), 5);

    // Next push should exceed limit
    assert!(stack.push("overflow".to_string()).is_err());
}

#[test]
fn test_callstack_frame_names() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();
    stack.push("main".to_string()).unwrap();
    stack.push("helper".to_string()).unwrap();
    stack.push("utility".to_string()).unwrap();

    let names = stack.frame_names();
    assert_eq!(names, vec!["global", "main", "helper", "utility"]);
}

#[test]
fn test_callstack_trace() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();
    stack.push("main".to_string()).unwrap();
    stack.push("process".to_string()).unwrap();

    let trace = stack.get_trace();
    assert!(trace.contains("Call Stack Trace"));
    assert!(trace.contains("global"));
    assert!(trace.contains("main"));
    assert!(trace.contains("process"));
}

#[test]
fn test_callstack_clear() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();
    stack.push("func1".to_string()).unwrap();
    stack.push("func2".to_string()).unwrap();
    stack.set_local("x".to_string(), Value::Number(42.0)).ok();

    assert_eq!(stack.depth(), 3);
    assert_eq!(stack.total_variables(), 1);

    stack.clear();

    assert_eq!(stack.depth(), 1);
    assert_eq!(stack.total_variables(), 0);
    assert!(stack.is_global());
}

#[test]
fn test_callstack_call_count() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();
    assert_eq!(stack.call_count(), 0);

    for i in 0..10 {
        stack.push(format!("func{}", i)).ok();
        assert_eq!(stack.call_count(), (i + 1) as u64);
    }

    // Clearing resets call count
    stack.clear();
    assert_eq!(stack.call_count(), 0);
}

#[test]
fn test_function_integration_with_runtime() {
    let engine = RuntimeEngine::new();

    // Test that we can create a user function and check its properties
    let func = UserFunction::new(
        "calculate".to_string(),
        vec!["a".to_string(), "b".to_string(), "c".to_string()],
        vec![],
    );

    assert_eq!(func.get_name(), "calculate");
    assert_eq!(func.param_count(), 3);

    // Set some closure variables
    func.set_closure_var("multiplier".to_string(), Value::Number(10.0));

    // Verify closure
    assert_eq!(
        func.get_closure_var("multiplier").unwrap().to_number(),
        10.0
    );
}

#[test]
fn test_function_call_sequence() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();

    // Simulate: main() -> process() -> calculate()
    stack.push("main".to_string()).unwrap();
    assert_eq!(stack.depth(), 2);

    stack.push("process".to_string()).unwrap();
    assert_eq!(stack.depth(), 3);

    stack.set_local("result".to_string(), Value::Number(0.0)).ok();

    stack.push("calculate".to_string()).unwrap();
    assert_eq!(stack.depth(), 4);

    stack.set_local("temp".to_string(), Value::Number(42.0)).ok();

    // Can access variables through the stack
    assert_eq!(stack.get_variable("result").unwrap().to_number(), 0.0);
    assert_eq!(stack.get_variable("temp").unwrap().to_number(), 42.0);

    // Return from calculate
    stack.pop();
    assert_eq!(stack.depth(), 3);

    // Return from process
    stack.pop();
    assert_eq!(stack.depth(), 2);

    // Return from main
    stack.pop();
    assert_eq!(stack.depth(), 1);
}

#[test]
fn test_nested_function_calls() {
    use freelang_runtime::runtime::CallStack;

    let mut stack = CallStack::new();

    // Create a deeply nested call
    for i in 0..50 {
        stack.push(format!("func_{}", i)).ok();
    }

    assert_eq!(stack.depth(), 51); // 1 global + 50 functions

    // Unwind the stack
    for _ in 0..50 {
        stack.pop();
    }

    assert_eq!(stack.depth(), 1);
    assert!(stack.is_global());
}

#[test]
fn test_function_native_signature() {
    fn add_numbers(args: Vec<Value>) -> Result<Value, String> {
        if args.len() < 2 {
            return Err("add_numbers requires 2 arguments".to_string());
        }

        let sum = args[0].to_number() + args[1].to_number();
        Ok(Value::Number(sum))
    }

    let kind = FunctionKind::Native(add_numbers);

    // Test calling with correct args
    let result = kind.call(vec![Value::Number(5.0), Value::Number(3.0)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 8.0);

    // Test calling with insufficient args
    let result = kind.call(vec![Value::Number(5.0)]);
    assert!(result.is_err());
}

#[test]
fn test_multiple_user_functions() {
    let func1 = UserFunction::new(
        "function1".to_string(),
        vec!["x".to_string()],
        vec![],
    );

    let func2 = UserFunction::new(
        "function2".to_string(),
        vec!["a".to_string(), "b".to_string()],
        vec![],
    );

    let func3 = UserFunction::new(
        "function3".to_string(),
        vec!["p".to_string(), "q".to_string(), "r".to_string()],
        vec![],
    );

    assert_eq!(func1.param_count(), 1);
    assert_eq!(func2.param_count(), 2);
    assert_eq!(func3.param_count(), 3);

    assert_eq!(func1.get_name(), "function1");
    assert_eq!(func2.get_name(), "function2");
    assert_eq!(func3.get_name(), "function3");
}
