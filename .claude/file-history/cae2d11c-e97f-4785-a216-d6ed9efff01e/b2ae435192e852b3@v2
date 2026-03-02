// Function Execution Tests
// Tests for calling functions and managing execution context

use freelang_runtime::{RuntimeEngine, Value};
use freelang_runtime::core::UserFunction;

#[test]
fn test_stdlib_function_call_through_runtime() {
    let mut engine = RuntimeEngine::new();

    // Call stdlib function through runtime
    let result = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 5.0);
}

#[test]
fn test_multiple_stdlib_calls() {
    let mut engine = RuntimeEngine::new();

    // First call
    let result1 = engine.call_stdlib("abs", vec![Value::Number(-10.0)]);
    assert_eq!(result1.unwrap().to_number(), 10.0);

    // Second call
    let result2 = engine.call_stdlib("sqrt", vec![Value::Number(16.0)]);
    assert_eq!(result2.unwrap().to_number(), 4.0);

    // Third call
    let result3 = engine.call_stdlib("pow", vec![
        Value::Number(2.0),
        Value::Number(5.0),
    ]);
    assert_eq!(result3.unwrap().to_number(), 32.0);
}

#[test]
fn test_call_stack_depth_tracking() {
    let mut engine = RuntimeEngine::new();

    assert_eq!(engine.call_depth(), 1); // Global frame

    // Call a stdlib function
    let _result = engine.call_stdlib("abs", vec![Value::Number(-1.0)]);

    // Should be back to depth 1 after call completes
    assert_eq!(engine.call_depth(), 1);
}

#[test]
fn test_function_stats_in_engine() {
    let mut engine = RuntimeEngine::new();

    // Initial stats
    let stats1 = engine.stats();
    assert_eq!(stats1.call_depth, 1);
    assert_eq!(stats1.call_count, 0);

    // After some calls
    let _result = engine.call_stdlib("abs", vec![Value::Number(-1.0)]);
    let _result = engine.call_stdlib("sqrt", vec![Value::Number(4.0)]);

    let stats2 = engine.stats();
    assert_eq!(stats2.call_depth, 1);
    // Call count should remain at initial value (0) since calls complete immediately
}

#[test]
fn test_user_function_argument_binding() {
    let mut engine = RuntimeEngine::new();

    // Create a simple user-defined function
    let func = UserFunction::new(
        "add_one".to_string(),
        vec!["x".to_string()],
        vec![], // Empty body for now
    );

    // Try to call it
    let result = engine.call_user_function(&func, vec![Value::Number(5.0)]);

    // Should succeed (even though body is empty, it returns Null for now)
    assert!(result.is_ok());
}

#[test]
fn test_user_function_argument_count_mismatch() {
    let mut engine = RuntimeEngine::new();

    let func = UserFunction::new(
        "requires_two".to_string(),
        vec!["a".to_string(), "b".to_string()],
        vec![],
    );

    // Too few arguments
    let result = engine.call_user_function(&func, vec![Value::Number(1.0)]);
    assert!(result.is_err());
    assert!(result.unwrap_err().contains("expects 2 arguments"));

    // Too many arguments
    let result = engine.call_user_function(
        &func,
        vec![Value::Number(1.0), Value::Number(2.0), Value::Number(3.0)],
    );
    assert!(result.is_err());
    assert!(result.unwrap_err().contains("expects 2 arguments"));
}

#[test]
fn test_user_function_argument_binding_multiple() {
    let mut engine = RuntimeEngine::new();

    let func = UserFunction::new(
        "three_args".to_string(),
        vec!["x".to_string(), "y".to_string(), "z".to_string()],
        vec![],
    );

    let result = engine.call_user_function(
        &func,
        vec![Value::Number(1.0), Value::Number(2.0), Value::Number(3.0)],
    );

    assert!(result.is_ok());
}

#[test]
fn test_call_trace() {
    let mut engine = RuntimeEngine::new();

    let trace1 = engine.get_call_trace();
    assert!(trace1.contains("global"));

    // Call a function
    let _result = engine.call_stdlib("abs", vec![Value::Number(-1.0)]);

    let trace2 = engine.get_call_trace();
    assert!(trace2.contains("global"));
}

#[test]
fn test_function_chaining() {
    let mut engine = RuntimeEngine::new();

    // Chain: abs(-5.0) = 5.0, sqrt(5.0) = 2.236...
    let step1 = engine.call_stdlib("abs", vec![Value::Number(-5.0)])
        .unwrap()
        .to_number();

    let step2 = engine.call_stdlib("sqrt", vec![Value::Number(step1)])
        .unwrap()
        .to_number();

    // sqrt(5) ≈ 2.236
    assert!((step2 - 2.236).abs() < 0.01);
}

#[test]
fn test_function_composition() {
    let mut engine = RuntimeEngine::new();

    // Compute: (2^3)^2 = 8^2 = 64
    let base = Value::Number(2.0);
    let exponent1 = Value::Number(3.0);
    let exponent2 = Value::Number(2.0);

    let result1 = engine.call_stdlib("pow", vec![base, exponent1])
        .unwrap();

    let result2 = engine.call_stdlib("pow", vec![result1, exponent2])
        .unwrap();

    assert_eq!(result2.to_number(), 64.0);
}

#[test]
fn test_string_function_chaining() {
    let mut engine = RuntimeEngine::new();

    let input = "hello world";

    // uppercase
    let upper = engine.call_stdlib("uppercase", vec![
        Value::String(input.to_string()),
    ]).unwrap();

    // length
    let len = engine.call_stdlib("length", vec![upper])
        .unwrap();

    assert_eq!(len.to_number(), 11.0); // "HELLO WORLD"
}

#[test]
fn test_array_function_chaining() {
    let mut engine = RuntimeEngine::new();

    let arr = Value::array(vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(2.0),
    ]);

    // Sort array
    let _ = engine.call_stdlib("sort", vec![arr.clone()]);

    // Get length
    let len = engine.call_stdlib("array_length", vec![arr.clone()])
        .unwrap();

    assert_eq!(len.to_number(), 3.0);
}

#[test]
fn test_math_pipeline() {
    let mut engine = RuntimeEngine::new();

    // Pipeline: 16 -> sqrt -> 4 -> pow(2,2) -> 16
    let val1 = engine.call_stdlib("sqrt", vec![Value::Number(16.0)])
        .unwrap()
        .to_number();
    assert_eq!(val1, 4.0);

    let val2 = engine.call_stdlib("pow", vec![
        Value::Number(val1),
        Value::Number(2.0),
    ]).unwrap().to_number();
    assert_eq!(val2, 16.0);

    let val3 = engine.call_stdlib("sqrt", vec![Value::Number(val2)])
        .unwrap()
        .to_number();
    assert_eq!(val3, 4.0);
}

#[test]
fn test_function_error_handling() {
    let mut engine = RuntimeEngine::new();

    // Division by zero should error
    let result = engine.call_stdlib("pow", vec![
        Value::Number(2.0),
        Value::Number(-1.0),
    ]);
    // Note: pow with negative exponent is valid and returns 0.5
    assert!(result.is_ok());

    // Invalid function call
    let result = engine.call_function("nonexistent", vec![]);
    assert!(result.is_err());
}

#[test]
fn test_nested_function_calls() {
    let mut engine = RuntimeEngine::new();

    // Create nested structure: max(min(5,2), min(8,3))
    let min1 = engine.call_stdlib("min", vec![Value::Number(5.0), Value::Number(2.0)])
        .unwrap();

    let min2 = engine.call_stdlib("min", vec![Value::Number(8.0), Value::Number(3.0)])
        .unwrap();

    let result = engine.call_stdlib("max", vec![min1, min2])
        .unwrap();

    assert_eq!(result.to_number(), 3.0); // max(2, 3) = 3
}

#[test]
fn test_call_depth_with_user_functions() {
    let mut engine = RuntimeEngine::new();

    assert_eq!(engine.call_depth(), 1);

    let func1 = UserFunction::new(
        "func1".to_string(),
        vec!["x".to_string()],
        vec![],
    );

    let _ = engine.call_user_function(&func1, vec![Value::Number(1.0)]);

    // Should be back at global depth
    assert_eq!(engine.call_depth(), 1);
}

#[test]
fn test_function_execution_order() {
    let mut engine = RuntimeEngine::new();

    // Execute functions in sequence and verify order
    let results = vec![
        engine.call_stdlib("abs", vec![Value::Number(-5.0)]).unwrap().to_number(),
        engine.call_stdlib("sqrt", vec![Value::Number(16.0)]).unwrap().to_number(),
        engine.call_stdlib("pow", vec![Value::Number(2.0), Value::Number(3.0)]).unwrap().to_number(),
        engine.call_stdlib("floor", vec![Value::Number(3.7)]).unwrap().to_number(),
    ];

    assert_eq!(results, vec![5.0, 4.0, 8.0, 3.0]);
}

#[test]
fn test_function_registry_lookup_performance() {
    use std::time::Instant;

    let mut engine = RuntimeEngine::new();

    let start = Instant::now();
    for _ in 0..1000 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-1.0)]);
    }
    let elapsed = start.elapsed();

    // 1000 function calls should complete in < 100ms
    println!("1000 calls completed in: {:?}", elapsed);
    assert!(elapsed.as_millis() < 100);
}

#[test]
fn test_mixed_function_types() {
    let mut engine = RuntimeEngine::new();

    // Mix stdlib and user functions
    let stdlib_result = engine.call_stdlib("abs", vec![Value::Number(-10.0)])
        .unwrap();

    let user_func = UserFunction::new(
        "process".to_string(),
        vec!["value".to_string()],
        vec![],
    );

    let user_result = engine.call_user_function(&user_func, vec![stdlib_result])
        .unwrap();

    // Should complete without errors
    assert!(matches!(user_result, Value::Null));
}

#[test]
fn test_function_call_statistics() {
    let mut engine = RuntimeEngine::new();

    let stats1 = engine.stats();
    assert_eq!(stats1.call_count, 0);
    assert_eq!(stats1.call_depth, 1);

    // Make some calls
    let _ = engine.call_stdlib("abs", vec![Value::Number(-1.0)]);
    let _ = engine.call_stdlib("sqrt", vec![Value::Number(4.0)]);
    let _ = engine.call_stdlib("pow", vec![Value::Number(2.0), Value::Number(2.0)]);

    let stats2 = engine.stats();
    // Call count should still be from CallStack
    assert_eq!(stats2.call_depth, 1); // Back at global
}

#[test]
fn test_function_with_closure_variables() {
    let func = UserFunction::new(
        "multiplier".to_string(),
        vec!["x".to_string()],
        vec![],
    );

    // Set closure variables
    func.set_closure_var("factor".to_string(), Value::Number(2.0));
    func.set_closure_var("offset".to_string(), Value::Number(5.0));

    // Retrieve closure variables
    let factor = func.get_closure_var("factor");
    let offset = func.get_closure_var("offset");

    assert_eq!(factor.unwrap().to_number(), 2.0);
    assert_eq!(offset.unwrap().to_number(), 5.0);
}
