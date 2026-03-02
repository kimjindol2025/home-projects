// FreeLang Rust Runtime - Main Entry Point
// Phase B: Independent execution engine with 91 standard functions

mod core;
mod memory;
mod runtime;
mod stdlib;

use runtime::RuntimeEngine;
use core::Value;

fn main() {
    println!("🚀 FreeLang Rust Runtime v0.2.0");
    println!("═══════════════════════════════════════════════════════");
    println!();

    // Initialize runtime
    let engine = RuntimeEngine::new();

    println!("✅ Runtime initialized");
    println!("📚 Standard library loaded: {} functions", engine.function_count());
    println!();

    // Demonstrate function registry
    println!("📋 Available function categories:");
    println!("   • I/O Functions (15)");
    println!("   • String Functions (18)");
    println!("   • Array Functions (12)");
    println!("   • Math Functions (15)");
    println!("   • System Functions (8)");
    println!("   • Crypto Functions (8)");
    println!("   • JSON Functions (5)");
    println!();

    // Demo: Math operations
    println!("🧮 Math Operations Demo:");

    // abs
    if let Ok(result) = engine.call_stdlib("abs", vec![Value::Number(-42.0)]) {
        println!("   abs(-42) = {}", result.to_number());
    }

    // sqrt
    if let Ok(result) = engine.call_stdlib("sqrt", vec![Value::Number(16.0)]) {
        println!("   sqrt(16) = {}", result.to_number());
    }

    // pow
    if let Ok(result) = engine.call_stdlib("pow", vec![
        Value::Number(2.0),
        Value::Number(8.0),
    ]) {
        println!("   pow(2, 8) = {}", result.to_number());
    }

    // min/max
    if let Ok(result) = engine.call_stdlib("min", vec![
        Value::Number(5.0),
        Value::Number(2.0),
        Value::Number(8.0),
    ]) {
        println!("   min(5, 2, 8) = {}", result.to_number());
    }

    println!();

    // Demo: String operations
    println!("📝 String Operations Demo:");

    // concat
    if let Ok(result) = engine.call_stdlib("concat", vec![
        Value::String("Hello".to_string()),
        Value::String(" ".to_string()),
        Value::String("FreeLang".to_string()),
    ]) {
        println!("   concat(\"Hello\", \" \", \"FreeLang\") = \"{}\"", result.to_string());
    }

    // uppercase
    if let Ok(result) = engine.call_stdlib("uppercase", vec![
        Value::String("hello world".to_string()),
    ]) {
        println!("   uppercase(\"hello world\") = \"{}\"", result.to_string());
    }

    // length
    if let Ok(result) = engine.call_stdlib("length", vec![
        Value::String("freelang".to_string()),
    ]) {
        println!("   length(\"freelang\") = {}", result.to_number());
    }

    println!();

    // Demo: Crypto operations
    println!("🔐 Crypto Operations Demo:");

    // hash
    if let Ok(result) = engine.call_stdlib("hash", vec![
        Value::String("freelang".to_string()),
    ]) {
        println!("   hash(\"freelang\") = \"{}\"", result.to_string());
    }

    // hex_encode
    if let Ok(result) = engine.call_stdlib("hex_encode", vec![
        Value::String("hello".to_string()),
    ]) {
        println!("   hex_encode(\"hello\") = \"{}\"", result.to_string());
    }

    println!();

    // Demo: System functions
    println!("⏱️  System Functions Demo:");

    // time
    if let Ok(result) = engine.call_stdlib("time", vec![]) {
        println!("   current_time() = {} seconds", result.to_number() as u64);
    }

    // random
    if let Ok(result) = engine.call_stdlib("random", vec![]) {
        println!("   random() = {:.4}", result.to_number());
    }

    // cpu_count
    if let Ok(result) = engine.call_stdlib("cpu_count", vec![]) {
        println!("   cpu_count() = {} cores", result.to_number() as usize);
    }

    println!();
    println!("═══════════════════════════════════════════════════════");
    println!("✅ Runtime ready for execution");
    println!("🎯 Ready to run FreeLang programs!");
}
