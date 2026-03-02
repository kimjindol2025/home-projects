// String Interning Tests
// Tests for StringInterner integration with CallStack
// Target: 80% deduplication for function names, 15% performance improvement

use freelang_runtime::runtime::{CallStack, StringInterner};
use freelang_runtime::Value;

#[test]
fn test_string_interner_basic() {
    let mut interner = StringInterner::new();

    let s1 = interner.intern("global");
    let s2 = interner.intern("main");
    let s3 = interner.intern("global");

    assert_eq!(s1.as_str(), "global");
    assert_eq!(s2.as_str(), "main");
    assert_eq!(s3.as_str(), "global");
    assert_eq!(interner.pool_size(), 2);
}

#[test]
fn test_string_interner_deduplication() {
    let mut interner = StringInterner::new();

    // Intern same string 100 times
    for _ in 0..100 {
        interner.intern("function_name");
    }

    let stats = interner.stats();
    assert_eq!(stats.new_strings, 1);
    assert_eq!(stats.total_requests, 100);
    assert_eq!(stats.cache_hits, 99);
    assert!(stats.hit_rate() > 95.0);
}

#[test]
fn test_string_interner_function_names() {
    let mut interner = StringInterner::new();

    // Simulate CallStack with repeated function names
    let names = vec!["global", "main", "process", "helper", "utility"];

    for _ in 0..50 {
        for name in &names {
            interner.intern(name);
        }
    }

    let stats = interner.stats();
    println!("Function names stats: {}", stats);

    // Should have only 5 unique strings
    assert_eq!(stats.new_strings, 5);
    assert!(stats.hit_rate() > 95.0);

    // Memory savings calculation
    let total_interned = stats.new_strings + stats.cache_hits;
    let bytes_saved = stats.deduplicated_bytes;
    println!("Total: {}, Saved: {}", total_interned, bytes_saved);
    assert!(bytes_saved > 0);
}

#[test]
fn test_callstack_with_interning() {
    let mut stack = CallStack::new();

    // Simulate multiple function calls
    for i in 0..10 {
        let func_name = format!("function_{}", i % 3); // Only 3 unique names
        stack.push(func_name).ok();
        stack.pop();
    }

    let stats = stack.interner_stats();
    println!("CallStack interner stats: {}", stats);

    // Should have good hit rate for repeated function names
    assert!(stats.hit_rate() > 0.0);
}

#[test]
fn test_callstack_deep_nesting_interning() {
    let mut stack = CallStack::new();

    // Create deep nesting with repeated names
    for i in 0..20 {
        let func_name = if i % 2 == 0 {
            "process".to_string()
        } else {
            "helper".to_string()
        };
        stack.push(func_name).ok();
    }

    let stats = stack.interner_stats();
    println!("Deep nesting stats: {}", stats);

    // Should deduplicate "process" and "helper"
    assert!(stats.cache_hits > 10);
    assert!(stats.hit_rate() > 50.0);

    // Unwind
    for _ in 0..20 {
        stack.pop();
    }

    assert!(stack.is_global());
}

#[test]
fn test_string_interner_memory_savings() {
    let mut interner = StringInterner::new();

    // Intern typical CallStack function names
    let names = vec![
        "global", "main", "process", "handler", "callback", "utility",
        "helper", "worker", "executor", "evaluator",
    ];

    for _ in 0..100 {
        for name in &names {
            interner.intern(name);
        }
    }

    let stats = interner.stats();

    // Memory savings should be significant
    let bytes_without_interning = names.iter().map(|n| n.len()).sum::<usize>() * 100;
    let bytes_with_interning = names.iter().map(|n| n.len()).sum::<usize>();

    println!("Without interning: {} bytes", bytes_without_interning);
    println!("With interning: {} bytes", bytes_with_interning);
    println!("Savings: {} bytes ({:.1}% reduction)",
        bytes_without_interning - bytes_with_interning,
        ((bytes_without_interning - bytes_with_interning) as f64 / bytes_without_interning as f64) * 100.0
    );

    assert!(stats.hit_rate() > 90.0);
}

#[test]
fn test_interner_with_callstack_operations() {
    let mut stack = CallStack::new();

    // Simulate realistic CallStack usage
    stack.push("main".to_string()).unwrap();
    stack.set_local("x".to_string(), Value::Number(1.0)).ok();

    stack.push("process".to_string()).unwrap();
    stack.set_local("y".to_string(), Value::Number(2.0)).ok();

    stack.push("helper".to_string()).unwrap();
    stack.set_local("z".to_string(), Value::Number(3.0)).ok();

    // Pop and check interner stats
    stack.pop();
    stack.pop();
    stack.pop();

    assert!(stack.is_global());

    let stats = stack.interner_stats();
    println!("Final stats: {}", stats);
}

#[test]
fn test_interner_bytes_saved() {
    let mut interner = StringInterner::new();

    // Intern 5 unique strings 20 times each
    for _ in 0..20 {
        interner.intern("global");
        interner.intern("main");
        interner.intern("process");
        interner.intern("helper");
        interner.intern("utility");
    }

    let bytes_saved = interner.bytes_saved();
    println!("Total bytes saved: {}", bytes_saved);

    // Each word has been saved multiple times
    // "global" (6) + "main" (4) + "process" (7) + "helper" (6) + "utility" (7) = 30 bytes
    // Saved 20 times each minus 1 (initial) = 19 times
    // Expected: 30 * 19 = 570 bytes
    let expected = (6 + 4 + 7 + 6 + 7) * 19;
    println!("Expected bytes saved: {}", expected);
    assert_eq!(bytes_saved, expected);
}

#[test]
#[ignore] // Performance test
fn benchmark_interner_vs_string() {
    use std::time::Instant;

    // Benchmark: 100,000 calls with 10 unique function names
    let unique_names = vec!["func_0", "func_1", "func_2", "func_3", "func_4",
                            "func_5", "func_6", "func_7", "func_8", "func_9"];

    // Without interning
    let start = Instant::now();
    let mut strings = Vec::new();
    for i in 0..100_000 {
        strings.push(unique_names[i % 10].to_string());
    }
    let without = start.elapsed();

    // With interning
    let start = Instant::now();
    let mut interner = StringInterner::new();
    for i in 0..100_000 {
        interner.intern(unique_names[i % 10]);
    }
    let with = start.elapsed();

    println!("Without interning: {:?}", without);
    println!("With interning: {:?}", with);
    println!("Overhead: {:.2}%", (with.as_nanos() as f64 / without.as_nanos() as f64 - 1.0) * 100.0);

    let stats = interner.stats();
    println!("{}", stats);
}

#[test]
fn test_interner_equality() {
    let mut interner = StringInterner::new();

    let s1 = interner.intern("test");
    let s2 = interner.intern("test");

    assert_eq!(s1, s2);
    assert_eq!(s1.as_str(), s2.as_str());
}

#[test]
fn test_interner_display() {
    let mut interner = StringInterner::new();

    let s = interner.intern("hello world");
    assert_eq!(format!("{}", s), "hello world");
    assert_eq!(format!("{:?}", s), "InternedString(hello world)");
}

#[test]
fn test_interner_clear() {
    let mut interner = StringInterner::new();

    interner.intern("test1");
    interner.intern("test2");

    assert_eq!(interner.pool_size(), 2);

    interner.clear();

    assert_eq!(interner.pool_size(), 0);
    let stats = interner.stats();
    assert_eq!(stats.total_requests, 0);
}

#[test]
fn test_interning_memory_efficiency() {
    let mut stack = CallStack::new();

    // Simulate 1000 function calls with only 5 unique names
    for i in 0..1000 {
        let name = match i % 5 {
            0 => "global",
            1 => "main",
            2 => "process",
            3 => "handler",
            _ => "helper",
        };
        stack.push(name.to_string()).ok();
        stack.pop();
    }

    let stats = stack.interner_stats();

    println!("Interning efficiency after 1000 calls:");
    println!("  Pool size: {}", stats.new_strings);
    println!("  Cache hits: {}", stats.cache_hits);
    println!("  Hit rate: {:.1}%", stats.hit_rate());

    // Should have high deduplication
    assert_eq!(stats.new_strings, 5);
    assert!(stats.hit_rate() > 95.0);
    assert!(stats.cache_hits > 950);
}
