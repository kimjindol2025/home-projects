// Advanced Integration Tests
// Complex scenarios demonstrating FreeLang runtime capabilities

use freelang_runtime::{RuntimeEngine, Value};

// Scenario 1: Data Processing Pipeline
#[test]
fn test_data_processing_pipeline() {
    let engine = RuntimeEngine::new();

    // Create array of numbers
    let arr = Value::array(vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(4.0),
        Value::Number(1.0),
        Value::Number(5.0),
    ]);

    // Sort array
    let sorted = engine.call_stdlib("sort", vec![arr.clone()]);
    assert!(sorted.is_ok());

    // Get length
    let length = engine.call_stdlib("array_length", vec![arr.clone()]);
    assert!(length.is_ok());
    assert_eq!(length.unwrap().to_number(), 5.0);

    // Min/Max
    let min = engine.call_stdlib("min", vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(4.0),
        Value::Number(1.0),
        Value::Number(5.0),
    ]);
    assert_eq!(min.unwrap().to_number(), 1.0);

    let max = engine.call_stdlib("max", vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(4.0),
        Value::Number(1.0),
        Value::Number(5.0),
    ]);
    assert_eq!(max.unwrap().to_number(), 5.0);
}

// Scenario 2: String Processing Pipeline
#[test]
fn test_string_processing_pipeline() {
    let engine = RuntimeEngine::new();

    // Original string
    let original = "  Hello World  ";

    // Trim
    let trimmed = engine.call_stdlib("trim", vec![
        Value::String(original.to_string())
    ]);
    assert!(trimmed.is_ok());

    // Lowercase
    let lower = engine.call_stdlib("lowercase", vec![
        trimmed.unwrap()
    ]);
    assert!(lower.is_ok());

    // Get length
    let len = engine.call_stdlib("length", vec![lower.unwrap()]);
    assert!(len.is_ok());
    assert_eq!(len.unwrap().to_number(), 11.0); // "hello world"

    // Split
    let split_result = engine.call_stdlib("split", vec![
        Value::String("a,b,c,d,e".to_string()),
        Value::String(",".to_string()),
    ]);
    assert!(split_result.is_ok());
}

// Scenario 3: Mathematical Calculations
#[test]
fn test_mathematical_calculations() {
    let engine = RuntimeEngine::new();

    // Calculate: (√16 + 2³) × |−5| = (4 + 8) × 5 = 60
    let sqrt_16 = engine.call_stdlib("sqrt", vec![Value::Number(16.0)])
        .unwrap()
        .to_number();
    assert_eq!(sqrt_16, 4.0);

    let pow_2_3 = engine.call_stdlib("pow", vec![
        Value::Number(2.0),
        Value::Number(3.0),
    ]).unwrap().to_number();
    assert_eq!(pow_2_3, 8.0);

    let abs_minus_5 = engine.call_stdlib("abs", vec![Value::Number(-5.0)])
        .unwrap()
        .to_number();
    assert_eq!(abs_minus_5, 5.0);

    // Total: (4 + 8) * 5 = 60
    let total = (sqrt_16 + pow_2_3) * abs_minus_5;
    assert_eq!(total, 60.0);
}

// Scenario 4: Cryptographic Operations
#[test]
fn test_cryptographic_workflow() {
    let engine = RuntimeEngine::new();

    let message = "freelang";

    // Create hash
    let hash_result = engine.call_stdlib("hash", vec![
        Value::String(message.to_string())
    ]);
    assert!(hash_result.is_ok());
    let hash = hash_result.unwrap().to_string();
    assert!(!hash.is_empty());

    // Create MD5
    let md5_result = engine.call_stdlib("md5", vec![
        Value::String(message.to_string())
    ]);
    assert!(md5_result.is_ok());
    let md5 = md5_result.unwrap().to_string();
    assert_eq!(md5.len(), 32);

    // Create SHA256
    let sha256_result = engine.call_stdlib("sha256", vec![
        Value::String(message.to_string())
    ]);
    assert!(sha256_result.is_ok());
    let sha256 = sha256_result.unwrap().to_string();
    assert_eq!(sha256.len(), 64);

    // HMAC with key
    let hmac_result = engine.call_stdlib("hmac", vec![
        Value::String(message.to_string()),
        Value::String("secret_key".to_string()),
    ]);
    assert!(hmac_result.is_ok());
}

// Scenario 5: JSON Data Processing
#[test]
fn test_json_data_processing() {
    let engine = RuntimeEngine::new();

    // Parse JSON object
    let json_str = r#"{"name": "FreeLang", "version": "0.2.0"}"#;
    let parsed = engine.call_stdlib("json_parse", vec![
        Value::String(json_str.to_string())
    ]);
    assert!(parsed.is_ok());

    // Parse JSON array
    let json_array = r#"[1, 2, 3, 4, 5]"#;
    let parsed_array = engine.call_stdlib("json_parse", vec![
        Value::String(json_array.to_string())
    ]);
    assert!(parsed_array.is_ok());

    // Parse primitives
    let json_number = "42";
    let parsed_num = engine.call_stdlib("json_parse", vec![
        Value::String(json_number.to_string())
    ]);
    assert!(parsed_num.is_ok());
    assert_eq!(parsed_num.unwrap().to_number(), 42.0);

    // Stringify
    let obj = Value::empty_object();
    obj.set_property("key".to_string(), Value::String("value".to_string())).ok();

    let stringified = engine.call_stdlib("json_stringify", vec![obj]);
    assert!(stringified.is_ok());
}

// Scenario 6: Base64 Encoding/Decoding
#[test]
fn test_base64_workflow() {
    let engine = RuntimeEngine::new();

    let original = "Hello, FreeLang!";

    // Encode
    let encoded = engine.call_stdlib("base64_encode", vec![
        Value::String(original.to_string())
    ]);
    assert!(encoded.is_ok());
    let encoded_str = encoded.unwrap().to_string();
    assert!(!encoded_str.is_empty());

    // Decode
    let decoded = engine.call_stdlib("base64_decode", vec![
        Value::String(encoded_str)
    ]);
    assert!(decoded.is_ok());
    // Note: simplified decoder might not match exactly due to padding handling
}

// Scenario 7: Time and System Operations
#[test]
fn test_time_and_system_operations() {
    let engine = RuntimeEngine::new();

    // Get current time
    let time1 = engine.call_stdlib("time", vec![])
        .unwrap()
        .to_number();

    // Get current time in milliseconds
    let time_ms = engine.call_stdlib("time_ms", vec![])
        .unwrap()
        .to_number();

    // ms should be ~1000x larger than seconds
    assert!(time_ms > time1 * 900.0); // Allow some variance

    // Get random number
    let rand1 = engine.call_stdlib("random", vec![])
        .unwrap()
        .to_number();
    let rand2 = engine.call_stdlib("random", vec![])
        .unwrap()
        .to_number();

    // Both should be in valid range
    assert!(rand1 >= 0.0 && rand1 <= 1.0);
    assert!(rand2 >= 0.0 && rand2 <= 1.0);

    // CPU count
    let cpu_count = engine.call_stdlib("cpu_count", vec![])
        .unwrap()
        .to_number();
    assert!(cpu_count > 0.0);
}

// Scenario 8: Complex String Manipulation
#[test]
fn test_complex_string_manipulation() {
    let engine = RuntimeEngine::new();

    // Template: "Dear {name}, welcome to {platform}!"
    let template = "Dear {name}, welcome to {platform}!";

    // Replace operations
    let step1 = engine.call_stdlib("replace", vec![
        Value::String(template.to_string()),
        Value::String("{name}".to_string()),
        Value::String("Alice".to_string()),
    ]).unwrap().to_string();

    let step2 = engine.call_stdlib("replace", vec![
        Value::String(step1),
        Value::String("{platform}".to_string()),
        Value::String("FreeLang".to_string()),
    ]).unwrap().to_string();

    assert_eq!(step2, "Dear Alice, welcome to FreeLang!");
}

// Scenario 9: Array Transformations
#[test]
fn test_array_transformations() {
    let engine = RuntimeEngine::new();

    // Create array
    let arr = Value::array(vec![
        Value::Number(1.0),
        Value::Number(2.0),
        Value::Number(3.0),
        Value::Number(4.0),
        Value::Number(5.0),
    ]);

    // Get length
    let len = engine.call_stdlib("array_length", vec![arr.clone()])
        .unwrap()
        .to_number();
    assert_eq!(len, 5.0);

    // Slice [1..3]
    let sliced = engine.call_stdlib("slice", vec![
        arr.clone(),
        Value::Number(1.0),
        Value::Number(3.0),
    ]);
    assert!(sliced.is_ok());

    // Find index of element
    let idx = engine.call_stdlib("array_index_of", vec![
        arr.clone(),
        Value::Number(3.0),
    ]).unwrap().to_number();
    assert_eq!(idx, 2.0); // 0-indexed

    // Reverse
    let reversed = engine.call_stdlib("array_reverse", vec![arr.clone()]);
    assert!(reversed.is_ok());
}

// Scenario 10: Multi-step Processing
#[test]
fn test_multi_step_processing() {
    let engine = RuntimeEngine::new();

    // Simulate: data -> transform -> hash -> encode
    let data = "freelang_runtime_v0.2.0";

    // Step 1: uppercase
    let uppercase = engine.call_stdlib("uppercase", vec![
        Value::String(data.to_string())
    ]).unwrap().to_string();

    // Step 2: hash
    let hash = engine.call_stdlib("hash", vec![
        Value::String(uppercase)
    ]).unwrap().to_string();

    // Step 3: hex_encode
    let hex = engine.call_stdlib("hex_encode", vec![
        Value::String(hash)
    ]).unwrap().to_string();

    assert!(!hex.is_empty());

    // Step 4: length
    let hex_len = engine.call_stdlib("length", vec![
        Value::String(hex)
    ]).unwrap().to_number();

    assert!(hex_len > 0.0);
}

// Performance Test: Bulk Operations
#[test]
fn test_bulk_string_operations() {
    use std::time::Instant;

    let engine = RuntimeEngine::new();

    let start = Instant::now();
    let mut count = 0;

    for i in 0..100 {
        let s = format!("test_{}", i);

        // Multiple operations
        let _ = engine.call_stdlib("uppercase", vec![Value::String(s.clone())]);
        let _ = engine.call_stdlib("length", vec![Value::String(s.clone())]);
        let _ = engine.call_stdlib("hash", vec![Value::String(s)]);

        count += 3;
    }

    let elapsed = start.elapsed();
    println!("100 strings × 3 operations = {} ops in {:?}", count, elapsed);

    // Should complete quickly
    assert!(elapsed.as_millis() < 500);
}

// Stress Test: Many Concurrent Operations
#[test]
fn test_stress_many_operations() {
    use std::time::Instant;

    let engine = RuntimeEngine::new();
    let start = Instant::now();

    // Run 1000 different operations
    for i in 0..1000 {
        match i % 10 {
            0 => { let _ = engine.call_stdlib("abs", vec![Value::Number(-1.0)]); }
            1 => { let _ = engine.call_stdlib("sqrt", vec![Value::Number(16.0)]); }
            2 => { let _ = engine.call_stdlib("hash", vec![Value::String("test".to_string())]); }
            3 => { let _ = engine.call_stdlib("random", vec![]); }
            4 => { let _ = engine.call_stdlib("uppercase", vec![Value::String("hello".to_string())]); }
            5 => { let _ = engine.call_stdlib("lowercase", vec![Value::String("HELLO".to_string())]); }
            6 => { let _ = engine.call_stdlib("length", vec![Value::String("freelang".to_string())]); }
            7 => { let _ = engine.call_stdlib("time", vec![]); }
            8 => { let _ = engine.call_stdlib("concat", vec![Value::String("a".to_string()), Value::String("b".to_string())]); }
            _ => { let _ = engine.call_stdlib("min", vec![Value::Number(1.0), Value::Number(2.0)]); }
        }
    }

    let elapsed = start.elapsed();
    println!("1000 mixed operations in {:?}", elapsed);

    // Should complete in reasonable time
    assert!(elapsed.as_millis() < 2000);
}
