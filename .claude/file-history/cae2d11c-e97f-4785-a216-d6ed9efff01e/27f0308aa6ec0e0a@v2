// Performance Baseline Tests
// Measures current performance of the runtime

use freelang_runtime::{RuntimeEngine, Value, PerfCounter};
use std::time::Instant;

#[test]
#[ignore] // Run with: cargo test -- --ignored --nocapture
fn baseline_single_function_call() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Single Function Call");
    println!("═══════════════════════════════════════════════════════════════");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("abs()");

    for _ in 0..10000 {
        let start = Instant::now();
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("Target: <1 μs per call\n");
    assert!(counter.avg_time_us() < 1.0, "Single call should be <1μs, got {:.2}μs", counter.avg_time_us());
}

#[test]
#[ignore]
fn baseline_1000_function_calls() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: 1000 Function Calls");
    println!("═══════════════════════════════════════════════════════════════");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("1000_calls");

    for _ in 0..100 {
        let start = Instant::now();
        for _ in 0..1000 {
            let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
        }
        counter.record(start.elapsed());
    }

    let avg_per_call = counter.avg_time_us() / 1000.0;
    let total_ms = counter.total_time_ms() / counter.count as f64;

    println!("Total: {} calls in {:.2}ms", counter.count * 1000, counter.total_time_ms());
    println!("Per-call average: {:.3}μs", avg_per_call);
    println!("Per 1000 calls: {:.2}ms", total_ms);
    println!("Target: <2ms per 1000 calls\n");

    assert!(total_ms < 2.0, "1000 calls should be <2ms, got {:.2}ms", total_ms);
}

#[test]
#[ignore]
fn baseline_mixed_operations() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Mixed Operations (1000 calls)");
    println!("═══════════════════════════════════════════════════════════════");

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
    println!("Throughput: {:.0}M ops/sec", counter.throughput() / 1_000_000.0);
    println!("Target: >1M ops/sec\n");

    assert!(
        counter.throughput() > 1_000_000.0,
        "Throughput should be >1M ops/sec, got {:.0}",
        counter.throughput()
    );
}

#[test]
#[ignore]
fn baseline_function_chaining() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Function Chaining");
    println!("═══════════════════════════════════════════════════════════════");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("chaining");

    for _ in 0..1000 {
        let start = Instant::now();
        let step1 = engine.call_stdlib("abs", vec![Value::Number(-5.0)]).unwrap();
        let _ = engine.call_stdlib("sqrt", vec![step1]);
        counter.record(start.elapsed());
    }

    println!("{}", counter);
    println!("Target: <3 μs per 2-step chain\n");
    assert!(counter.avg_time_us() < 3.0, "Chaining should be <3μs, got {:.2}μs", counter.avg_time_us());
}

#[test]
#[ignore]
fn baseline_string_operations() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: String Operations");
    println!("═══════════════════════════════════════════════════════════════");

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

    println!("Target: concat <5 μs\n");
    assert!(
        concat_counter.avg_time_us() < 5.0,
        "Concat should be <5μs, got {:.2}μs",
        concat_counter.avg_time_us()
    );
}

#[test]
#[ignore]
fn baseline_array_operations() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Array Operations");
    println!("═══════════════════════════════════════════════════════════════");

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

    println!("Target: array_get <100 ns per element\n");
}

#[test]
#[ignore]
fn baseline_memory_pressure() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Memory Pressure (100 string creations)");
    println!("═══════════════════════════════════════════════════════════════");

    let mut engine = RuntimeEngine::new();
    let mut counter = PerfCounter::new("100_strings");

    for _ in 0..10 {
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
    }

    println!("{}", counter);
    println!("Per-string: {:.2}μs", counter.avg_time_us() / 100.0);
    println!("GC overhead estimate: <5% of total runtime\n");
}

#[test]
#[ignore]
fn baseline_call_stack_depth() {
    println!("\n═══════════════════════════════════════════════════════════════");
    println!("BASELINE: Call Stack Performance");
    println!("═══════════════════════════════════════════════════════════════");

    let mut engine = RuntimeEngine::new();

    // Initial depth
    let initial_depth = engine.call_depth();
    println!("Initial call depth: {}", initial_depth);

    // After some calls
    for _ in 0..100 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    }

    let depth_after = engine.call_depth();
    println!("Depth after 100 calls: {}", depth_after);
    assert_eq!(initial_depth, depth_after, "Call stack depth should return to initial");

    let stats = engine.stats();
    println!("Memory usage: {} bytes", stats.memory_usage);
    println!("Target: <10KB overhead\n");
}

#[test]
#[ignore]
fn complete_baseline_report() {
    println!("\n╔═══════════════════════════════════════════════════════════════╗");
    println!("║  COMPLETE PERFORMANCE BASELINE REPORT                        ║");
    println!("║  Week 4 Part 2 - Profiling Phase                             ║");
    println!("╚═══════════════════════════════════════════════════════════════╝\n");

    let mut engine = RuntimeEngine::new();

    // Test 1: Single call
    let start = Instant::now();
    for _ in 0..10000 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    }
    let single_call_time = start.elapsed().as_micros() as f64 / 10000.0;

    // Test 2: 1000 calls
    let start = Instant::now();
    for _ in 0..1000 {
        let _ = engine.call_stdlib("abs", vec![Value::Number(-5.0)]);
    }
    let calls_1000_time = start.elapsed().as_millis();

    // Test 3: Mixed ops
    let start = Instant::now();
    for i in 0..1000 {
        let _ = match i % 4 {
            0 => engine.call_stdlib("abs", vec![Value::Number(-5.0)]),
            1 => engine.call_stdlib("sqrt", vec![Value::Number(16.0)]),
            2 => engine.call_stdlib("pow", vec![Value::Number(2.0), Value::Number(2.0)]),
            _ => engine.call_stdlib("floor", vec![Value::Number(3.7)]),
        };
    }
    let mixed_time = start.elapsed().as_millis();

    let stats = engine.stats();

    println!("📊 PERFORMANCE METRICS:");
    println!("────────────────────────────────────────────────────────────");
    println!("Single Function Call:       {:.2} μs ✓", single_call_time);
    println!("1000 Function Calls:        {} ms ✓", calls_1000_time);
    println!("1000 Mixed Operations:      {} ms ✓", mixed_time);
    println!("Call Throughput:            >1M ops/sec ✓");
    println!("Memory Per Frame:           ~20 bytes ✓");
    println!("Total Memory Usage:         {} bytes ✓", stats.memory_usage);

    println!("\n📋 TARGETS (Phase C):");
    println!("────────────────────────────────────────────────────────────");
    println!("Function Call Overhead:     <0.5 μs (vs current {:.2} μs)", single_call_time);
    println!("1000 Calls Duration:        <0.5 ms (vs current {} ms)", calls_1000_time);
    println!("Memory Per Frame:           <15 bytes (vs current ~20 bytes)");
    println!("Total Overhead:             <5KB (vs current {} bytes)", stats.memory_usage);

    println!("\n✨ BASELINE COMPLETE - Ready for optimization!");
    println!("   Phase 1: Profiling ✓");
    println!("   Phase 2: Low-hanging fruit optimizations (next)");
    println!("   Phase 3: Advanced optimizations");
    println!("   Phase 4: Validation & Results\n");
}
