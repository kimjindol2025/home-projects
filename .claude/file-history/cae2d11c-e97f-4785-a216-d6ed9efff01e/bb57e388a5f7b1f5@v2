// JSON Functions (5)
// JSON parsing and serialization

use crate::core::Value;

/// Parse JSON string to Value
pub fn parse(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("parse requires 1 argument".to_string());
    }

    let json_str = args[0].to_string();
    parse_json(&json_str)
}

/// Convert Value to JSON string
pub fn stringify(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("stringify requires 1 argument".to_string());
    }

    let json_str = value_to_json(&args[0]);
    Ok(Value::String(json_str))
}

/// Pretty-print JSON (formatted)
pub fn dumps(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("dumps requires 1 argument".to_string());
    }

    let json_str = value_to_json_pretty(&args[0], 0);
    Ok(Value::String(json_str))
}

/// Load JSON from string (alias for parse)
pub fn load(args: Vec<Value>) -> Result<Value, String> {
    parse(args)
}

/// Validate JSON string
pub fn validate(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("validate requires 1 argument".to_string());
    }

    let json_str = args[0].to_string();
    match parse_json(&json_str) {
        Ok(_) => Ok(Value::Bool(true)),
        Err(_) => Ok(Value::Bool(false))
    }
}

// Helper functions

fn parse_json(s: &str) -> Result<Value, String> {
    let trimmed = s.trim();

    if trimmed.is_empty() {
        return Err("Empty JSON string".to_string());
    }

    let (value, _) = parse_value(trimmed)?;
    Ok(value)
}

fn parse_value(s: &str) -> Result<(Value, &str), String> {
    let trimmed = skip_whitespace(s);

    match trimmed.chars().next() {
        Some('"') => parse_string(trimmed),
        Some('{') => parse_object(trimmed),
        Some('[') => parse_array(trimmed),
        Some('t') | Some('f') => parse_bool(trimmed),
        Some('n') => parse_null(trimmed),
        Some(c) if c == '-' || c.is_ascii_digit() => parse_number(trimmed),
        _ => Err("Invalid JSON value".to_string()),
    }
}

fn parse_string(s: &str) -> Result<(Value, &str), String> {
    if !s.starts_with('"') {
        return Err("Expected string".to_string());
    }

    let mut result = String::new();
    let mut chars = s[1..].chars();
    let mut remaining = &s[1..];

    loop {
        match chars.next() {
            Some('"') => {
                let rest = remaining.split('"').nth(1).unwrap_or("");
                return Ok((Value::String(result), rest));
            }
            Some('\\') => {
                if let Some(c) = chars.next() {
                    match c {
                        'n' => result.push('\n'),
                        't' => result.push('\t'),
                        'r' => result.push('\r'),
                        '\\' => result.push('\\'),
                        '"' => result.push('"'),
                        _ => {
                            result.push('\\');
                            result.push(c);
                        }
                    }
                }
            }
            Some(c) => result.push(c),
            None => return Err("Unterminated string".to_string()),
        }
    }
}

fn parse_number(s: &str) -> Result<(Value, &str), String> {
    let mut end = 0;
    let mut has_dot = false;

    for (i, c) in s.chars().enumerate() {
        match c {
            '-' if i == 0 => {}
            '0'..='9' => {}
            '.' if !has_dot => has_dot = true,
            'e' | 'E' => {
                end = i + 1;
                if i + 1 < s.len() && (s.chars().nth(i + 1) == Some('+') || s.chars().nth(i + 1) == Some('-')) {
                    end += 1;
                }
                continue;
            }
            _ => {
                end = i;
                break;
            }
        }
        end = i + 1;
    }

    let num_str = &s[..end];
    match num_str.parse::<f64>() {
        Ok(n) => Ok((Value::Number(n), &s[end..])),
        Err(_) => Err("Invalid number".to_string()),
    }
}

fn parse_bool(s: &str) -> Result<(Value, &str), String> {
    if s.starts_with("true") {
        Ok((Value::Bool(true), &s[4..]))
    } else if s.starts_with("false") {
        Ok((Value::Bool(false), &s[5..]))
    } else {
        Err("Invalid boolean".to_string())
    }
}

fn parse_null(s: &str) -> Result<(Value, &str), String> {
    if s.starts_with("null") {
        Ok((Value::Null, &s[4..]))
    } else {
        Err("Invalid null".to_string())
    }
}

fn parse_array(s: &str) -> Result<(Value, &str), String> {
    if !s.starts_with('[') {
        return Err("Expected array".to_string());
    }

    let mut elements = Vec::new();
    let mut remaining = skip_whitespace(&s[1..]);

    if remaining.starts_with(']') {
        return Ok((Value::array(elements), &remaining[1..]));
    }

    loop {
        let (value, rest) = parse_value(remaining)?;
        elements.push(value);
        remaining = skip_whitespace(rest);

        if remaining.starts_with(',') {
            remaining = skip_whitespace(&remaining[1..]);
        } else if remaining.starts_with(']') {
            return Ok((Value::array(elements), &remaining[1..]));
        } else {
            return Err("Invalid array format".to_string());
        }
    }
}

fn parse_object(s: &str) -> Result<(Value, &str), String> {
    if !s.starts_with('{') {
        return Err("Expected object".to_string());
    }

    let obj = Value::empty_object();
    let mut remaining = skip_whitespace(&s[1..]);

    if remaining.starts_with('}') {
        return Ok((obj, &remaining[1..]));
    }

    loop {
        remaining = skip_whitespace(remaining);

        // Parse key
        if !remaining.starts_with('"') {
            return Err("Expected string key in object".to_string());
        }

        let (Value::String(key), rest) = parse_string(remaining)? else {
            return Err("Invalid object key".to_string());
        };

        remaining = skip_whitespace(rest);
        if !remaining.starts_with(':') {
            return Err("Expected ':' in object".to_string());
        }

        remaining = skip_whitespace(&remaining[1..]);

        // Parse value
        let (value, rest) = parse_value(remaining)?;
        obj.set_property(key, value)
            .map_err(|e| format!("Object set error: {}", e))?;

        remaining = skip_whitespace(rest);

        if remaining.starts_with(',') {
            remaining = skip_whitespace(&remaining[1..]);
        } else if remaining.starts_with('}') {
            return Ok((obj, &remaining[1..]));
        } else {
            return Err("Invalid object format".to_string());
        }
    }
}

fn skip_whitespace(s: &str) -> &str {
    s.trim_start()
}

fn value_to_json(value: &Value) -> String {
    match value {
        Value::Null => "null".to_string(),
        Value::Bool(b) => b.to_string(),
        Value::Number(n) => {
            if n.fract() == 0.0 {
                format!("{:.0}", n)
            } else {
                n.to_string()
            }
        }
        Value::String(s) => format!("\"{}\"", s.replace("\"", "\\\"").replace("\n", "\\n")),
        Value::Array(arr) => {
            let elements: Vec<String> = arr.borrow()
                .iter()
                .map(value_to_json)
                .collect();
            format!("[{}]", elements.join(","))
        }
        Value::Object(obj) => {
            let obj_ref = obj.borrow();
            let mut pairs: Vec<String> = Vec::new();

            for (key, val) in obj_ref.iter() {
                let json_key = format!("\"{}\"", key.replace("\"", "\\\""));
                let json_val = value_to_json(val);
                pairs.push(format!("{}:{}", json_key, json_val));
            }

            format!("{{{}}}", pairs.join(","))
        }
        _ => "null".to_string(),
    }
}

fn value_to_json_pretty(value: &Value, indent: usize) -> String {
    let indent_str = " ".repeat(indent);
    let next_indent_str = " ".repeat(indent + 2);

    match value {
        Value::Null => "null".to_string(),
        Value::Bool(b) => b.to_string(),
        Value::Number(n) => {
            if n.fract() == 0.0 {
                format!("{:.0}", n)
            } else {
                n.to_string()
            }
        }
        Value::String(s) => format!("\"{}\"", s.replace("\"", "\\\"").replace("\n", "\\n")),
        Value::Array(arr) => {
            let arr_ref = arr.borrow();
            if arr_ref.is_empty() {
                "[]".to_string()
            } else {
                let elements: Vec<String> = arr_ref
                    .iter()
                    .map(|v| format!("{}{}", next_indent_str, value_to_json_pretty(v, indent + 2)))
                    .collect();
                format!("[\n{}\n{}]", elements.join(",\n"), indent_str)
            }
        }
        Value::Object(obj) => {
            let obj_ref = obj.borrow();
            if obj_ref.is_empty() {
                "{}".to_string()
            } else {
                let mut pairs: Vec<String> = Vec::new();
                for (key, val) in obj_ref.iter() {
                    let json_key = format!("\"{}\"", key);
                    let json_val = value_to_json_pretty(val, indent + 2);
                    pairs.push(format!("{}{}: {}", next_indent_str, json_key, json_val));
                }
                format!("{{\n{}\n{}}}", pairs.join(",\n"), indent_str)
            }
        }
        _ => "null".to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_null() {
        let result = parse(vec![Value::String("null".to_string())]);
        assert!(result.is_ok());
        assert!(matches!(result.unwrap(), Value::Null));
    }

    #[test]
    fn test_parse_bool() {
        let result = parse(vec![Value::String("true".to_string())]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_parse_number() {
        let result = parse(vec![Value::String("42".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 42.0);
    }

    #[test]
    fn test_parse_string() {
        let result = parse(vec![Value::String("\"hello\"".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "hello");
    }

    #[test]
    fn test_stringify() {
        let val = Value::Number(42.0);
        let result = stringify(vec![val]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "42");
    }

    #[test]
    fn test_validate() {
        let result = validate(vec![Value::String("{}".to_string())]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());

        let result = validate(vec![Value::String("invalid".to_string())]);
        assert!(result.is_ok());
        assert!(!result.unwrap().is_truthy());
    }
}
