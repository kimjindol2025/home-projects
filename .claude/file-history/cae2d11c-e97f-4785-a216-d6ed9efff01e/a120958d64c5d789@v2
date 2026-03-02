// System Functions (8)
// System-level operations

use crate::core::Value;
use std::time::{SystemTime, UNIX_EPOCH};
use std::env;
use std::process;

/// Sleep for milliseconds
pub fn sleep(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("sleep requires 1 argument: milliseconds".to_string());
    }

    let ms = args[0].to_number() as u64;
    std::thread::sleep(std::time::Duration::from_millis(ms));
    Ok(Value::Null)
}

/// Exit program with exit code
pub fn exit(args: Vec<Value>) -> Result<Value, String> {
    let code = if args.is_empty() {
        0
    } else {
        args[0].to_number() as i32
    };
    process::exit(code);
}

/// Get environment variable
pub fn getenv(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("getenv requires 1 argument: variable name".to_string());
    }

    let var_name = args[0].to_string();
    match env::var(&var_name) {
        Ok(value) => Ok(Value::String(value)),
        Err(_) => Ok(Value::Null)
    }
}

/// Set environment variable
pub fn setenv(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("setenv requires 2 arguments: name, value".to_string());
    }

    let var_name = args[0].to_string();
    let var_value = args[1].to_string();
    env::set_var(&var_name, &var_value);
    Ok(Value::Null)
}

/// Get current Unix timestamp (milliseconds)
pub fn time_ms(args: Vec<Value>) -> Result<Value, String> {
    let _ = args;
    match SystemTime::now().duration_since(UNIX_EPOCH) {
        Ok(duration) => Ok(Value::Number(duration.as_millis() as f64)),
        Err(_) => Err("System time error".to_string())
    }
}

/// Get current Unix timestamp (seconds)
pub fn time(args: Vec<Value>) -> Result<Value, String> {
    let _ = args;
    match SystemTime::now().duration_since(UNIX_EPOCH) {
        Ok(duration) => Ok(Value::Number(duration.as_secs() as f64)),
        Err(_) => Err("System time error".to_string())
    }
}

/// Generate random number between 0 and 1
pub fn random(_args: Vec<Value>) -> Result<Value, String> {
    use std::time::{SystemTime};

    // Simple pseudo-random using system time
    let nanos = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.subsec_nanos())
        .unwrap_or(0);

    let random_val = ((nanos as f64) % 1000000.0) / 1000000.0;
    Ok(Value::Number(random_val))
}

/// Get CPU count
pub fn cpu_count(_args: Vec<Value>) -> Result<Value, String> {
    let count = num_cpus::get();
    Ok(Value::Number(count as f64))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_time() {
        let result = time(vec![]);
        assert!(result.is_ok());
        let timestamp = result.unwrap().to_number();
        assert!(timestamp > 0.0);
    }

    #[test]
    fn test_time_ms() {
        let result = time_ms(vec![]);
        assert!(result.is_ok());
        let timestamp = result.unwrap().to_number();
        assert!(timestamp > 0.0);
    }

    #[test]
    fn test_random() {
        let result = random(vec![]);
        assert!(result.is_ok());
        let rand_val = result.unwrap().to_number();
        assert!(rand_val >= 0.0 && rand_val <= 1.0);
    }

    #[test]
    fn test_getenv() {
        // Set a variable first
        setenv(vec![
            Value::String("TEST_VAR".to_string()),
            Value::String("test_value".to_string()),
        ]).ok();

        // Get it back
        let result = getenv(vec![Value::String("TEST_VAR".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "test_value");
    }

    #[test]
    fn test_setenv() {
        let result = setenv(vec![
            Value::String("MY_VAR".to_string()),
            Value::String("my_value".to_string()),
        ]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_cpu_count() {
        let result = cpu_count(vec![]);
        assert!(result.is_ok());
        let count = result.unwrap().to_number();
        assert!(count > 0.0);
    }

    #[test]
    fn test_sleep() {
        use std::time::Instant;
        let start = Instant::now();
        let result = sleep(vec![Value::Number(10.0)]); // Sleep 10ms
        let elapsed = start.elapsed().as_millis();

        assert!(result.is_ok());
        assert!(elapsed >= 10); // Should take at least 10ms
    }
}
