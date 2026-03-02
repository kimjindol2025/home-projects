// Math Functions (15)
// Mathematical operations

use crate::core::Value;
use std::f64::consts;

/// Absolute value
pub fn abs(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("abs requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().abs()))
}

/// Square root
pub fn sqrt(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("sqrt requires 1 argument".to_string());
    }
    let n = args[0].to_number();
    if n < 0.0 {
        return Err("Cannot take square root of negative number".to_string());
    }
    Ok(Value::Number(n.sqrt()))
}

/// Power (base^exponent)
pub fn pow(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("pow requires 2 arguments: base, exponent".to_string());
    }
    let base = args[0].to_number();
    let exponent = args[1].to_number();
    Ok(Value::Number(base.powf(exponent)))
}

/// Natural logarithm
pub fn ln(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("ln requires 1 argument".to_string());
    }
    let n = args[0].to_number();
    if n <= 0.0 {
        return Err("Cannot take logarithm of non-positive number".to_string());
    }
    Ok(Value::Number(n.ln()))
}

/// Logarithm base 10
pub fn log10(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("log10 requires 1 argument".to_string());
    }
    let n = args[0].to_number();
    if n <= 0.0 {
        return Err("Cannot take logarithm of non-positive number".to_string());
    }
    Ok(Value::Number(n.log10()))
}

/// Logarithm base 2
pub fn log2(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("log2 requires 1 argument".to_string());
    }
    let n = args[0].to_number();
    if n <= 0.0 {
        return Err("Cannot take logarithm of non-positive number".to_string());
    }
    Ok(Value::Number(n.log2()))
}

/// Sine (radians)
pub fn sin(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("sin requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().sin()))
}

/// Cosine (radians)
pub fn cos(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("cos requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().cos()))
}

/// Tangent (radians)
pub fn tan(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("tan requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().tan()))
}

/// Minimum value
pub fn min(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("min requires at least 1 argument".to_string());
    }

    let mut min_val = args[0].to_number();
    for arg in &args[1..] {
        let val = arg.to_number();
        if val < min_val {
            min_val = val;
        }
    }
    Ok(Value::Number(min_val))
}

/// Maximum value
pub fn max(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("max requires at least 1 argument".to_string());
    }

    let mut max_val = args[0].to_number();
    for arg in &args[1..] {
        let val = arg.to_number();
        if val > max_val {
            max_val = val;
        }
    }
    Ok(Value::Number(max_val))
}

/// Floor (round down)
pub fn floor(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("floor requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().floor()))
}

/// Ceiling (round up)
pub fn ceil(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("ceil requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().ceil()))
}

/// Round to nearest integer
pub fn round(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("round requires 1 argument".to_string());
    }
    Ok(Value::Number(args[0].to_number().round()))
}

/// Convert degrees to radians
pub fn radians(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("radians requires 1 argument".to_string());
    }
    let degrees = args[0].to_number();
    Ok(Value::Number(degrees * consts::PI / 180.0))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_abs() {
        assert_eq!(abs(vec![Value::Number(-5.0)]).unwrap().to_number(), 5.0);
        assert_eq!(abs(vec![Value::Number(5.0)]).unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_sqrt() {
        assert_eq!(sqrt(vec![Value::Number(16.0)]).unwrap().to_number(), 4.0);
        assert_eq!(sqrt(vec![Value::Number(25.0)]).unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_pow() {
        assert_eq!(
            pow(vec![Value::Number(2.0), Value::Number(3.0)])
                .unwrap()
                .to_number(),
            8.0
        );
        assert_eq!(
            pow(vec![Value::Number(5.0), Value::Number(2.0)])
                .unwrap()
                .to_number(),
            25.0
        );
    }

    #[test]
    fn test_min() {
        let result = min(vec![
            Value::Number(5.0),
            Value::Number(2.0),
            Value::Number(8.0),
        ]);
        assert_eq!(result.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_max() {
        let result = max(vec![
            Value::Number(5.0),
            Value::Number(2.0),
            Value::Number(8.0),
        ]);
        assert_eq!(result.unwrap().to_number(), 8.0);
    }

    #[test]
    fn test_floor() {
        assert_eq!(floor(vec![Value::Number(5.7)]).unwrap().to_number(), 5.0);
        assert_eq!(floor(vec![Value::Number(5.2)]).unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_ceil() {
        assert_eq!(ceil(vec![Value::Number(5.2)]).unwrap().to_number(), 6.0);
        assert_eq!(ceil(vec![Value::Number(5.7)]).unwrap().to_number(), 6.0);
    }

    #[test]
    fn test_round() {
        assert_eq!(round(vec![Value::Number(5.4)]).unwrap().to_number(), 5.0);
        assert_eq!(round(vec![Value::Number(5.5)]).unwrap().to_number(), 6.0);
    }

    #[test]
    fn test_sin() {
        let result = sin(vec![Value::Number(0.0)]).unwrap().to_number();
        assert!((result - 0.0).abs() < 0.0001);
    }

    #[test]
    fn test_cos() {
        let result = cos(vec![Value::Number(0.0)]).unwrap().to_number();
        assert!((result - 1.0).abs() < 0.0001);
    }
}
