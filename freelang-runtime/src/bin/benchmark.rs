// Performance Benchmarking Tool
// Measures current performance baseline

use freelang_runtime::{RuntimeEngine, Value, PerfCounter};
use std::time::Instant;

fn main() {
    println!("═══════════════════════════════════════════════════════════════");
    println!("   FreeLang Runtime - Performance Baseline Measurement");
    println!("═══════════════════════════════════════════════════════════════\n");

    // Test 1: Single Function Call
    benchmark_single_call();

    // Test 2: Function Lookup Performance
    benchmark_function_lookup();

    // Test 3: 1000 Function Calls
    benchmark_1000_calls();

    // Test 4: Mixed Operations
    benchmark_mixed_operations();

    // Test 5: Function Chaining
    benchmark_chaining();

    // Test 6: String Operations
    benchmark_strings();

    // Test 7: Array Operations
    benchmark_arrays();

    // Test 8: Memory Usage
    benchmark_memory();

    println!("\n═══════════════════════════════════════════════════════════════");
    println!("   Baseline measurement complete!");
    println!("═══════════════════════════════════════════════════════════════");
}

fn benchmark_single_call() {
    println!("\n📊 Test 1: Single Function Call");
    println!("────────────────────────────────────────────────────────────");

    let mut counter = PerfCounter::new("abs()");
    let mut engine = RuntimeEngine::new();

    for _ in 0..10000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("✓ Single call: {:.2}μs\n", counter.avg_time_us());
}

fn benchmark_function_lookup() {
    println!("\n📊 Test 2: Function Lookup (Registry)");
    println!("────────────────────────────────────────────────────────────");

    let engine = RuntimeEngine::new();
    let functions = vec!["abs", "sqrt", "pow", "sin", "cos", "floor", "ceil"];

    let start = Instant::now();
    for _ in 0..100000 {
        for func in &functions {
            let _ = engine.function_exists(func);
        }
    }
    let elapsed = start.elapsed();

    let total_lookups = 100000 * functions.len();
    let avg_ns = elapsed.as_nanos() as f64 / total_lookups as f64;

    println!("Lookups: {} in {:.2}ms", total_lookups, elapsed.as_secs_f64() * 1000.0);
    println!("Average: {:.0}ns per lookup", avg_ns);
    println!("Throughput: {:.0}M lookups/sec\n", total_lookups as f64 / elapsed.as_secs_f64() / 1_000_000.0);
}

fn benchmark_1000_calls() {
    println!("\n📊 Test 3: 1000 Function Calls");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("1000_calls");

    for _ in 0..100 {
        let start = Instant::now();
        for _ in 0..1000 {
            let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
        }
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("✓ Per-call average: {:.2}μs", counter.avg_time_us() / 1000.0);
    println!("✓ Per 1000-calls: {:.2}ms\n", counter.avg_time_us() / 1000.0);
}

fn benchmark_mixed_operations() {
    println!("\n📊 Test 4: Mixed Operations (1000 calls)");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("mixed_1000");

    for _ in 0..100 {
        let start = Instant::now();
        for i in 0..1000 {
            let _ = match i % 4 {
                0 => engine.call_stdlib("abs", vec![Value::Number(-5.0)]),
                1 => engine.call_stdlib("sqrt", vec![Value::Number(16.0)]),
                2 => engine.call_stdlib("pow", vec![Value::Number(2.0), Value::Number(2.0)]),
                _ => engine.call_stdlib("floor", vec![Value::Number(3.7)]),
            };
        }
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("✓ Throughput: {:.0}M ops/sec\n", counter.throughput() / 1_000_000.0);
}

fn benchmark_chaining() {
    println!("\n📊 Test 5: Function Chaining");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("chaining");

    for _ in 0..1000 {
        let start = Instant::now();
        let step1 = engine.call_stdlib("abs", vec![Value::Number(-5.0)]).unwrap();
        let _ = engine.call_stdlib("sqrt", vec![step1]);
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("✓ 2-step chain: {:.2}μs\n", counter.avg_time_us());
}

fn benchmark_strings() {
    println!("\n📊 Test 6: String Operations");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();

    // String length
    let mut len_counter = PerfCounter::new("length");
    for _ in 0..1000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("length", vec![Value::String("hello world".to_string())]);
        len_counter.record(start.elapsed());
    }
    println!("{}", len_counter);

    // String uppercase
    let mut upper_counter = PerfCounter::new("uppercase");
    for _ in 0..1000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("uppercase", vec![Value::String("hello world".to_string())]);
        upper_counter.record(start.elapsed());
    }
    println!("{}", upper_counter);

    // String concat
    let mut concat_counter = PerfCounter::new("concat");
    for _ in 0..1000 {
        let start = Instant::now();
        let _ = engine.call_stdlib(
            "concat",
            vec![Value::String("hello".to_string()), Value::String("world".to_string())],
        );
        concat_counter.record(start.elapsed());
    }
    println!("{}\n", concat_counter);
}

fn benchmark_arrays() {
    println!("\n📊 Test 7: Array Operations");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();
    let arr = Value::array(vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(2.0),
    ]);

    // Array length
    let mut len_counter = PerfCounter::new("array_length");
    for _ in 0..1000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("array_length", vec![arr.clone()]);
        len_counter.record(start.elapsed());
    }
    println!("{}", len_counter);

    // Array get
    let mut get_counter = PerfCounter::new("array_get");
    for _ in 0..1000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("array_get", vec![arr.clone(), Value::Number(1.0)]);
        get_counter.record(start.elapsed());
    }
    println!("{}\n", get_counter);
}

fn benchmark_memory() {
    println!("\n📊 Test 8: Memory Pressure (100 string creations)");
    println!("────────────────────────────────────────────────────────────");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("100_strings");

    for iteration in 0..10 {
        let start = Instant::now();
        for i in 0..100 {
            let _ = engine.call_stdlib(
                "concat",
                vec![
                    Value::String(format!("str{}", i)),
                    Value::String("suffix".to_string()),
                ],
            );
        }
        counter.record(start.elapsed());

        if iteration == 9 {
            println!("Iteration {}: {:.2}ms", iteration + 1, start.elapsed().as_secs_f64() * 1000.0);
        }
    }

    println!("{}", counter);
    println!("✓ Per-string: {:.2}μs\n", counter.avg_time_us() / 100.0);
}
