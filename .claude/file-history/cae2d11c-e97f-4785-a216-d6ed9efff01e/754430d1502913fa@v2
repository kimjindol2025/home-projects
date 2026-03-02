// String Functions (18)
// String manipulation and analysis

use crate::core::Value;

/// Concatenate strings
pub fn concat(args: Vec<Value>) -> Result<Value, String> {
    let mut result = String::new();
    for arg in args {
        result.push_str(&arg.to_string());
    }
    Ok(Value::String(result))
}

/// Convert to uppercase
pub fn uppercase(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("uppercase requires 1 argument".to_string());
    }
    Ok(Value::String(args[0].to_string().to_uppercase()))
}

/// Convert to lowercase
pub fn lowercase(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("lowercase requires 1 argument".to_string());
    }
    Ok(Value::String(args[0].to_string().to_lowercase()))
}

/// Get string length
pub fn length(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("length requires 1 argument".to_string());
    }
    let s = args[0].to_string();
    Ok(Value::Number(s.len() as f64))
}

/// Get character at index
pub fn char_at(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("char_at requires 2 arguments: string, index".to_string());
    }
    let s = args[0].to_string();
    let index = args[1].to_number() as usize;

    if index < s.len() {
        Ok(Value::String(s.chars().nth(index).unwrap().to_string()))
    } else {
        Err("Index out of bounds".to_string())
    }
}

/// Get substring
pub fn substring(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 3 {
        return Err("substring requires 3 arguments: string, start, end".to_string());
    }
    let s = args[0].to_string();
    let start = args[1].to_number() as usize;
    let end = args[2].to_number() as usize;

    let chars: Vec<char> = s.chars().collect();
    if start >= chars.len() || end > chars.len() || start > end {
        return Err("Invalid substring range".to_string());
    }

    let result: String = chars[start..end].iter().collect();
    Ok(Value::String(result))
}

/// Split string by delimiter
pub fn split(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("split requires 2 arguments: string, delimiter".to_string());
    }
    let s = args[0].to_string();
    let delimiter = args[1].to_string();

    let parts: Vec<Value> = s
        .split(&delimiter)
        .map(|part| Value::String(part.to_string()))
        .collect();

    Ok(Value::array(parts))
}

/// Join array elements with delimiter
pub fn join(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("join requires 2 arguments: array, delimiter".to_string());
    }

    let delimiter = args[1].to_string();

    match &args[0] {
        Value::Array(arr) => {
            let arr_ref = arr.borrow();
            let parts: Vec<String> = arr_ref
                .iter()
                .map(|v| v.to_string())
                .collect();
            Ok(Value::String(parts.join(&delimiter)))
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Trim whitespace from both ends
pub fn trim(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("trim requires 1 argument".to_string());
    }
    Ok(Value::String(args[0].to_string().trim().to_string()))
}

/// Trim left whitespace
pub fn trim_left(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("trim_left requires 1 argument".to_string());
    }
    Ok(Value::String(args[0].to_string().trim_start().to_string()))
}

/// Trim right whitespace
pub fn trim_right(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("trim_right requires 1 argument".to_string());
    }
    Ok(Value::String(args[0].to_string().trim_end().to_string()))
}

/// Replace substring occurrences
pub fn replace(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 3 {
        return Err("replace requires 3 arguments: string, from, to".to_string());
    }
    let s = args[0].to_string();
    let from = args[1].to_string();
    let to = args[2].to_string();

    Ok(Value::String(s.replace(&from, &to)))
}

/// Find index of substring
pub fn index_of(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("index_of requires 2 arguments: string, substring".to_string());
    }
    let s = args[0].to_string();
    let substring = args[1].to_string();

    match s.find(&substring) {
        Some(idx) => Ok(Value::Number(idx as f64)),
        None => Ok(Value::Number(-1.0))
    }
}

/// Check if string contains substring
pub fn contains(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("contains requires 2 arguments: string, substring".to_string());
    }
    let s = args[0].to_string();
    let substring = args[1].to_string();

    Ok(Value::Bool(s.contains(&substring)))
}

/// Check if string starts with substring
pub fn starts_with(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("starts_with requires 2 arguments: string, prefix".to_string());
    }
    let s = args[0].to_string();
    let prefix = args[1].to_string();

    Ok(Value::Bool(s.starts_with(&prefix)))
}

/// Check if string ends with substring
pub fn ends_with(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("ends_with requires 2 arguments: string, suffix".to_string());
    }
    let s = args[0].to_string();
    let suffix = args[1].to_string();

    Ok(Value::Bool(s.ends_with(&suffix)))
}

/// Repeat string n times
pub fn repeat(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("repeat requires 2 arguments: string, count".to_string());
    }
    let s = args[0].to_string();
    let count = args[1].to_number() as usize;

    Ok(Value::String(s.repeat(count)))
}

/// Reverse string
pub fn reverse(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("reverse requires 1 argument".to_string());
    }
    let s = args[0].to_string();
    Ok(Value::String(s.chars().rev().collect()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_concat() {
        let result = concat(vec![
            Value::String("hello".to_string()),
            Value::String(" ".to_string()),
            Value::String("world".to_string()),
        ]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "hello world");
    }

    #[test]
    fn test_uppercase() {
        let result = uppercase(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "HELLO");
    }

    #[test]
    fn test_lowercase() {
        let result = lowercase(vec![Value::String("HELLO".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "hello");
    }

    #[test]
    fn test_length() {
        let result = length(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_split() {
        let result = split(vec![
            Value::String("a,b,c".to_string()),
            Value::String(",".to_string()),
        ]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_trim() {
        let result = trim(vec![Value::String("  hello  ".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "hello");
    }

    #[test]
    fn test_replace() {
        let result = replace(vec![
            Value::String("hello world".to_string()),
            Value::String("world".to_string()),
            Value::String("friend".to_string()),
        ]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "hello friend");
    }

    #[test]
    fn test_index_of() {
        let result = index_of(vec![
            Value::String("hello".to_string()),
            Value::String("ll".to_string()),
        ]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_contains() {
        let result = contains(vec![
            Value::String("hello".to_string()),
            Value::String("ll".to_string()),
        ]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_starts_with() {
        let result = starts_with(vec![
            Value::String("hello".to_string()),
            Value::String("he".to_string()),
        ]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_ends_with() {
        let result = ends_with(vec![
            Value::String("hello".to_string()),
            Value::String("lo".to_string()),
        ]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_reverse() {
        let result = reverse(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "olleh");
    }
}
