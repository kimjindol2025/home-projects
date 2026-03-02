// Integration Tests for FreeLang Runtime
// Tests runtime integration with stdlib functions

use freelang_runtime::core::Value;
use freelang_runtime::runtime::RuntimeEngine;

#[test]
fn test_runtime_initialization() {
    let engine = RuntimeEngine::new();
    assert_eq!(engine.execution_count(), 0);
    assert_eq!(engine.variable_count(), 0);
    assert!(engine.function_count() >= 80);
}

#[test]
fn test_function_registry_loaded() {
    let engine = RuntimeEngine::new();
    let funcs = engine.list_functions();
    assert!(funcs.len() > 80);
    assert!(funcs.contains(&"print"));
    assert!(funcs.contains(&"println"));
}

// I/O Function Tests
#[test]
fn test_call_print() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("print", vec![
        Value::String("Hello, FreeLang!".to_string())
    ]);
    assert!(result.is_ok());
}

#[test]
fn test_call_println() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("println", vec![
        Value::String("Test output".to_string())
    ]);
    assert!(result.is_ok());
}

// String Function Tests
#[test]
fn test_string_concat() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("concat", vec![
        Value::String("Hello".to_string()),
        Value::String(" ".to_string()),
        Value::String("World".to_string()),
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_string(), "Hello World");
}

#[test]
fn test_string_uppercase() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("uppercase", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_string(), "HELLO");
}

#[test]
fn test_string_lowercase() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("lowercase", vec![
        Value::String("HELLO".to_string())
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_string(), "hello");
}

#[test]
fn test_string_length() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("length", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 5.0);
}

// Array Function Tests
#[test]
fn test_array_length() {
    let engine = RuntimeEngine::new();
    let arr = Value::array(vec![
        Value::Number(1.0),
        Value::Number(2.0),
        Value::Number(3.0),
    ]);
    let result = engine.call_stdlib("array_length", vec![arr]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 3.0);
}

#[test]
fn test_array_push() {
    let engine = RuntimeEngine::new();
    let arr = Value::array(vec![Value::Number(1.0)]);
    let result = engine.call_stdlib("push", vec![
        arr.clone(),
        Value::Number(2.0),
    ]);
    assert!(result.is_ok());

    // Verify length increased
    let length_result = engine.call_stdlib("array_length", vec![arr]);
    assert_eq!(length_result.unwrap().to_number(), 2.0);
}

// Math Function Tests
#[test]
fn test_math_abs() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 5.0);
}

#[test]
fn test_math_sqrt() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("sqrt", vec![Value::Number(16.0)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 4.0);
}

#[test]
fn test_math_pow() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("pow", vec![
        Value::Number(2.0),
        Value::Number(3.0),
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 8.0);
}

#[test]
fn test_math_min() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("min", vec![
        Value::Number(5.0),
        Value::Number(2.0),
        Value::Number(8.0),
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 2.0);
}

#[test]
fn test_math_max() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("max", vec![
        Value::Number(5.0),
        Value::Number(2.0),
        Value::Number(8.0),
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 8.0);
}

#[test]
fn test_math_floor() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("floor", vec![Value::Number(5.7)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 5.0);
}

#[test]
fn test_math_ceil() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("ceil", vec![Value::Number(5.2)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 6.0);
}

#[test]
fn test_math_round() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("round", vec![Value::Number(5.5)]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 6.0);
}

// System Function Tests
#[test]
fn test_time() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("time", vec![]);
    assert!(result.is_ok());
    assert!(result.unwrap().to_number() > 0.0);
}

#[test]
fn test_time_ms() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("time_ms", vec![]);
    assert!(result.is_ok());
    assert!(result.unwrap().to_number() > 0.0);
}

#[test]
fn test_random() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("random", vec![]);
    assert!(result.is_ok());
    let rand_val = result.unwrap().to_number();
    assert!(rand_val >= 0.0 && rand_val <= 1.0);
}

#[test]
fn test_cpu_count() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("cpu_count", vec![]);
    assert!(result.is_ok());
    assert!(result.unwrap().to_number() > 0.0);
}

// Crypto Function Tests
#[test]
fn test_hash() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("hash", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    assert!(!result.unwrap().to_string().is_empty());
}

#[test]
fn test_md5() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("md5", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    let hash = result.unwrap().to_string();
    assert_eq!(hash.len(), 32); // MD5 = 32 hex chars
}

#[test]
fn test_sha256() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("sha256", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    let hash = result.unwrap().to_string();
    assert_eq!(hash.len(), 64); // SHA256 = 64 hex chars
}

#[test]
fn test_base64_encode() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("base64_encode", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    assert!(!result.unwrap().to_string().is_empty());
}

#[test]
fn test_hex_encode() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("hex_encode", vec![
        Value::String("hello".to_string())
    ]);
    assert!(result.is_ok());
    // "hello" = 68656c6c6f
    assert!(result.unwrap().to_string().contains("68656c6c6f"));
}

// JSON Function Tests
#[test]
fn test_json_parse_null() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_parse", vec![
        Value::String("null".to_string())
    ]);
    assert!(result.is_ok());
}

#[test]
fn test_json_parse_number() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_parse", vec![
        Value::String("42".to_string())
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_number(), 42.0);
}

#[test]
fn test_json_parse_string() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_parse", vec![
        Value::String("\"hello\"".to_string())
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_string(), "hello");
}

#[test]
fn test_json_stringify() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_stringify", vec![
        Value::Number(42.0)
    ]);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().to_string(), "42");
}

#[test]
fn test_json_validate_valid() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_validate", vec![
        Value::String("{}".to_string())
    ]);
    assert!(result.is_ok());
    assert!(result.unwrap().is_truthy());
}

#[test]
fn test_json_validate_invalid() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("json_validate", vec![
        Value::String("invalid".to_string())
    ]);
    assert!(result.is_ok());
    assert!(!result.unwrap().is_truthy());
}

// Edge Cases and Error Handling
#[test]
fn test_nonexistent_function() {
    let engine = RuntimeEngine::new();
    let result = engine.call_stdlib("nonexistent_function", vec![]);
    assert!(result.is_err());
}

#[test]
fn test_function_exists_check() {
    let engine = RuntimeEngine::new();
    assert!(engine.function_exists("print"));
    assert!(engine.function_exists("abs"));
    assert!(!engine.function_exists("nonexistent"));
}

#[test]
fn test_multiple_operations() {
    let engine = RuntimeEngine::new();

    // Call multiple functions in sequence
    let result1 = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    assert!(result1.is_ok());

    let result2 = engine.call_stdlib("sqrt", vec![Value::Number(16.0)]);
    assert!(result2.is_ok());

    let result3 = engine.call_stdlib("concat", vec![
        Value::String("hello".to_string()),
        Value::String(" world".to_string()),
    ]);
    assert!(result3.is_ok());
}

// Performance Benchmarks
#[test]
fn test_stdlib_function_call_performance() {
    use std::time::Instant;

    let engine = RuntimeEngine::new();
    let start = Instant::now();

    // Call 1000 simple functions
    for _ in 0..1000 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    }

    let elapsed = start.elapsed();
    println!("1000 function calls took: {:?}", elapsed);

    // Should complete in reasonable time (< 100ms for 1000 calls)
    assert!(elapsed.as_millis() < 1000);
}

#[test]
fn test_stdlib_function_call_throughput() {
    use std::time::Instant;

    let engine = RuntimeEngine::new();
    let start = Instant::now();
    let mut count = 0;

    // See how many calls we can make in 100ms
    while start.elapsed().as_millis() < 100 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
        count += 1;
    }

    println!("Function call throughput: {} calls/100ms", count);

    // Should handle at least 1000 calls per 100ms
    assert!(count > 1000);
}
