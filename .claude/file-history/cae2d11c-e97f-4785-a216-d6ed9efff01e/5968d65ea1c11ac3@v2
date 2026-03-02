// Function Registry
// Maps function names to stdlib implementations

use crate::core::Value;
use super::{io, string, array, math, system, crypto, json};
use std::collections::HashMap;

pub type FunctionPointer = fn(Vec<Value>) -> Result<Value, String>;

/// Global function registry
pub struct FunctionRegistry {
    functions: HashMap<String, FunctionPointer>,
}

impl FunctionRegistry {
    /// Create new function registry
    pub fn new() -> Self {
        let mut registry = FunctionRegistry {
            functions: HashMap::new(),
        };
        registry.register_all();
        registry
    }

    /// Register all stdlib functions
    fn register_all(&mut self) {
        // I/O functions (15)
        self.register("print", io::print);
        self.register("println", io::println);
        self.register("input", io::input);
        self.register("open", io::open);
        self.register("close", io::close);
        self.register("read_file", io::read_file);
        self.register("write_file", io::write_file);
        self.register("append_file", io::append_file);
        self.register("file_exists", io::file_exists);
        self.register("delete_file", io::delete_file);
        self.register("file_size", io::file_size);
        self.register("list_dir", io::list_dir);
        self.register("mkdir", io::mkdir);
        self.register("eprint", io::eprint);
        self.register("eprintln", io::eprintln);

        // String functions (18)
        self.register("concat", string::concat);
        self.register("uppercase", string::uppercase);
        self.register("lowercase", string::lowercase);
        self.register("length", string::length);
        self.register("char_at", string::char_at);
        self.register("substring", string::substring);
        self.register("split", string::split);
        self.register("join", string::join);
        self.register("trim", string::trim);
        self.register("trim_left", string::trim_left);
        self.register("trim_right", string::trim_right);
        self.register("replace", string::replace);
        self.register("index_of", string::index_of);
        self.register("contains", string::contains);
        self.register("starts_with", string::starts_with);
        self.register("ends_with", string::ends_with);
        self.register("repeat", string::repeat);
        self.register("reverse", string::reverse);

        // Array functions (12)
        self.register("array_length", array::length);
        self.register("push", array::push);
        self.register("pop", array::pop);
        self.register("get", array::get);
        self.register("set", array::set);
        self.register("slice", array::slice);
        self.register("array_index_of", array::index_of);
        self.register("array_contains", array::contains);
        self.register("array_reverse", array::reverse);
        self.register("sort", array::sort);
        self.register("array_join", array::join);

        // Math functions (15)
        self.register("abs", math::abs);
        self.register("sqrt", math::sqrt);
        self.register("pow", math::pow);
        self.register("ln", math::ln);
        self.register("log10", math::log10);
        self.register("log2", math::log2);
        self.register("sin", math::sin);
        self.register("cos", math::cos);
        self.register("tan", math::tan);
        self.register("min", math::min);
        self.register("max", math::max);
        self.register("floor", math::floor);
        self.register("ceil", math::ceil);
        self.register("round", math::round);
        self.register("radians", math::radians);

        // System functions (8)
        self.register("sleep", system::sleep);
        self.register("exit", system::exit);
        self.register("getenv", system::getenv);
        self.register("setenv", system::setenv);
        self.register("time_ms", system::time_ms);
        self.register("time", system::time);
        self.register("random", system::random);
        self.register("cpu_count", system::cpu_count);

        // Crypto functions (8)
        self.register("hash", crypto::hash);
        self.register("md5", crypto::md5);
        self.register("sha256", crypto::sha256);
        self.register("hmac", crypto::hmac);
        self.register("verify_hash", crypto::verify_hash);
        self.register("base64_encode", crypto::base64_encode);
        self.register("base64_decode", crypto::base64_decode);
        self.register("hex_encode", crypto::hex_encode);

        // JSON functions (5)
        self.register("json_parse", json::parse);
        self.register("json_stringify", json::stringify);
        self.register("json_dumps", json::dumps);
        self.register("json_load", json::load);
        self.register("json_validate", json::validate);
    }

    /// Register a function
    fn register(&mut self, name: &str, func: FunctionPointer) {
        self.functions.insert(name.to_string(), func);
    }

    /// Call a function by name
    pub fn call(&self, name: &str, args: Vec<Value>) -> Result<Value, String> {
        match self.functions.get(name) {
            Some(func) => func(args),
            None => Err(format!("Function not found: {}", name))
        }
    }

    /// Check if function exists
    pub fn exists(&self, name: &str) -> bool {
        self.functions.contains_key(name)
    }

    /// Get list of all function names
    pub fn list_functions(&self) -> Vec<&str> {
        self.functions.keys().map(|s| s.as_str()).collect()
    }

    /// Get function count
    pub fn function_count(&self) -> usize {
        self.functions.len()
    }
}

impl Default for FunctionRegistry {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_registry_creation() {
        let registry = FunctionRegistry::new();
        assert!(registry.function_count() > 80); // Should have 91 functions
    }

    #[test]
    fn test_function_exists() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("print"));
        assert!(registry.exists("println"));
        assert!(registry.exists("strlen")); // Some function that doesn't exist
    }

    #[test]
    fn test_call_print() {
        let registry = FunctionRegistry::new();
        let result = registry.call("print", vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_call_math_abs() {
        let registry = FunctionRegistry::new();
        let result = registry.call("abs", vec![Value::Number(-5.0)]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_nonexistent_function() {
        let registry = FunctionRegistry::new();
        let result = registry.call("nonexistent", vec![]);
        assert!(result.is_err());
    }

    #[test]
    fn test_list_functions() {
        let registry = FunctionRegistry::new();
        let funcs = registry.list_functions();
        assert!(funcs.len() > 80);
        assert!(funcs.contains(&"print"));
        assert!(funcs.contains(&"abs"));
    }

    #[test]
    fn test_io_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("read_file"));
        assert!(registry.exists("write_file"));
        assert!(registry.exists("file_exists"));
    }

    #[test]
    fn test_string_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("uppercase"));
        assert!(registry.exists("lowercase"));
        assert!(registry.exists("concat"));
    }

    #[test]
    fn test_array_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("push"));
        assert!(registry.exists("pop"));
        assert!(registry.exists("sort"));
    }

    #[test]
    fn test_math_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("abs"));
        assert!(registry.exists("sqrt"));
        assert!(registry.exists("pow"));
    }

    #[test]
    fn test_system_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("sleep"));
        assert!(registry.exists("getenv"));
        assert!(registry.exists("time"));
    }

    #[test]
    fn test_crypto_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("hash"));
        assert!(registry.exists("md5"));
        assert!(registry.exists("sha256"));
    }

    #[test]
    fn test_json_functions() {
        let registry = FunctionRegistry::new();
        assert!(registry.exists("json_parse"));
        assert!(registry.exists("json_stringify"));
    }
}
