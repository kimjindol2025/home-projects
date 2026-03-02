use criterion::{black_box, criterion_group, criterion_main, Criterion};
use freelang_runtime::{RuntimeEngine, Value};

fn benchmark_stdlib_function_call(c: &mut Criterion) {
    let mut group = c.benchmark_group("function_dispatch");

    group.bench_function("single_abs_call", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib("abs", vec![black_box(Value::Number(-5.0))]);
        });
    });

    group.bench_function("single_sqrt_call", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib("sqrt", vec![black_box(Value::Number(16.0))]);
        });
    });

    group.bench_function("single_pow_call", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib(
                "pow",
                vec![black_box(Value::Number(2.0)), black_box(Value::Number(3.0))],
            );
        });
    });

    group.finish();
}

fn benchmark_mixed_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("mixed_operations");

    group.bench_function("1000_mixed_calls", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            for i in 0..1000 {
                let _ = match i % 4 {
                    0 => engine.call_stdlib("abs", vec![black_box(Value::Number(-5.0))]),
                    1 => engine.call_stdlib("sqrt", vec![black_box(Value::Number(16.0))]),
                    2 => engine.call_stdlib(
                        "pow",
                        vec![black_box(Value::Number(2.0)), black_box(Value::Number(2.0))],
                    ),
                    _ => engine.call_stdlib("floor", vec![black_box(Value::Number(3.7))]),
                };
            }
        });
    });

    group.finish();
}

fn benchmark_function_chaining(c: &mut Criterion) {
    let mut group = c.benchmark_group("function_chaining");

    group.bench_function("chained_abs_sqrt", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let step1 = engine
                .call_stdlib("abs", vec![black_box(Value::Number(-5.0))])
                .unwrap();
            let _ = engine.call_stdlib("sqrt", vec![black_box(step1)]);
        });
    });

    group.bench_function("chained_pow_sqrt", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let step1 = engine
                .call_stdlib(
                    "pow",
                    vec![black_box(Value::Number(2.0)), black_box(Value::Number(3.0))],
                )
                .unwrap();
            let _ = engine.call_stdlib("sqrt", vec![black_box(step1)]);
        });
    });

    group.finish();
}

fn benchmark_string_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("string_operations");

    group.bench_function("string_length", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib(
                "length",
                vec![black_box(Value::String("hello world".to_string()))],
            );
        });
    });

    group.bench_function("string_uppercase", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib(
                "uppercase",
                vec![black_box(Value::String("hello world".to_string()))],
            );
        });
    });

    group.bench_function("string_concat", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib(
                "concat",
                vec![
                    black_box(Value::String("hello".to_string())),
                    black_box(Value::String("world".to_string())),
                ],
            );
        });
    });

    group.finish();
}

fn benchmark_array_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("array_operations");

    let arr = Value::array(vec![
        Value::Number(3.0),
        Value::Number(1.0),
        Value::Number(2.0),
    ]);

    group.bench_function("array_length", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib("array_length", vec![black_box(arr.clone())]);
        });
    });

    group.bench_function("array_get", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            let _ = engine.call_stdlib(
                "array_get",
                vec![black_box(arr.clone()), black_box(Value::Number(1.0))],
            );
        });
    });

    group.finish();
}

fn benchmark_memory_pressure(c: &mut Criterion) {
    let mut group = c.benchmark_group("memory_pressure");

    group.bench_function("create_many_strings", |b| {
        let mut engine = RuntimeEngine::new();
        b.iter(|| {
            for i in 0..100 {
                let _ = engine.call_stdlib(
                    "concat",
                    vec![
                        black_box(Value::String(format!("str{}", i))),
                        black_box(Value::String("suffix".to_string())),
                    ],
                );
            }
        });
    });

    group.finish();
}

criterion_group!(
    benches,
    benchmark_stdlib_function_call,
    benchmark_mixed_operations,
    benchmark_function_chaining,
    benchmark_string_operations,
    benchmark_array_operations,
    benchmark_memory_pressure
);
criterion_main!(benches);
