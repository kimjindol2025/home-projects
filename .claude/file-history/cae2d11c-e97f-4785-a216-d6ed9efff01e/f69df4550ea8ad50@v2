// Value Type System
// Represents all possible values in FreeLang runtime

use std::collections::HashMap;
use std::rc::Rc;
use std::cell::RefCell;
use std::fmt;

/// Main Value type - represents all FreeLang values
#[derive(Clone)]
pub enum Value {
    /// Null value
    Null,

    /// Boolean (true/false)
    Bool(bool),

    /// Number (f64 - supports both integers and floats)
    Number(f64),

    /// String (immutable UTF-8)
    String(String),

    /// Array (mutable, reference-counted)
    Array(Rc<RefCell<Vec<Value>>>),

    /// Object (mutable HashMap, reference-counted)
    Object(Rc<RefCell<HashMap<String, Value>>>),

    /// Function (callable)
    Function(Rc<Box<dyn Fn(Vec<Value>) -> Result<Value, String>>>),

    /// Error value
    Error(String),
}

impl Value {
    /// Get type name as string
    pub fn type_of(&self) -> &str {
        match self {
            Value::Null => "null",
            Value::Bool(_) => "boolean",
            Value::Number(_) => "number",
            Value::String(_) => "string",
            Value::Array(_) => "array",
            Value::Object(_) => "object",
            Value::Function(_) => "function",
            Value::Error(_) => "error",
        }
    }

    /// Check if value is truthy
    pub fn is_truthy(&self) -> bool {
        match self {
            Value::Null => false,
            Value::Bool(b) => *b,
            Value::Number(n) => *n != 0.0 && !n.is_nan(),
            Value::String(s) => !s.is_empty(),
            Value::Array(_) => true,
            Value::Object(_) => true,
            Value::Function(_) => true,
            Value::Error(_) => false,
        }
    }

    /// Convert to number
    pub fn to_number(&self) -> f64 {
        match self {
            Value::Null => 0.0,
            Value::Bool(b) => if *b { 1.0 } else { 0.0 },
            Value::Number(n) => *n,
            Value::String(s) => s.parse().unwrap_or(0.0),
            Value::Array(_) => f64::NAN,
            Value::Object(_) => f64::NAN,
            Value::Function(_) => f64::NAN,
            Value::Error(_) => f64::NAN,
        }
    }

    /// Convert to string
    pub fn to_string(&self) -> String {
        match self {
            Value::Null => "null".to_string(),
            Value::Bool(b) => b.to_string(),
            Value::Number(n) => {
                if n.fract() == 0.0 && !n.is_infinite() {
                    format!("{:.0}", n)
                } else {
                    n.to_string()
                }
            }
            Value::String(s) => s.clone(),
            Value::Array(arr) => {
                let arr_ref = arr.borrow();
                let elements: Vec<String> = arr_ref
                    .iter()
                    .map(|v| v.to_string())
                    .collect();
                format!("[{}]", elements.join(", "))
            }
            Value::Object(obj) => {
                let obj_ref = obj.borrow();
                let entries: Vec<String> = obj_ref
                    .iter()
                    .map(|(k, v)| format!("{}: {}", k, v.to_string()))
                    .collect();
                format!("{{{}}}", entries.join(", "))
            }
            Value::Function(_) => "[Function]".to_string(),
            Value::Error(e) => format!("Error: {}", e),
        }
    }

    /// Check equality
    pub fn equals(&self, other: &Value) -> bool {
        match (self, other) {
            (Value::Null, Value::Null) => true,
            (Value::Bool(a), Value::Bool(b)) => a == b,
            (Value::Number(a), Value::Number(b)) => a == b,
            (Value::String(a), Value::String(b)) => a == b,
            (Value::Array(a), Value::Array(b)) => {
                Rc::ptr_eq(a, b) || {
                    let a_ref = a.borrow();
                    let b_ref = b.borrow();
                    a_ref.len() == b_ref.len()
                        && a_ref
                            .iter()
                            .zip(b_ref.iter())
                            .all(|(x, y)| x.equals(y))
                }
            }
            (Value::Object(a), Value::Object(b)) => Rc::ptr_eq(a, b),
            _ => false,
        }
    }

    /// Create array from vec
    pub fn array(values: Vec<Value>) -> Self {
        Value::Array(Rc::new(RefCell::new(values)))
    }

    /// Create object from HashMap
    pub fn object(map: HashMap<String, Value>) -> Self {
        Value::Object(Rc::new(RefCell::new(map)))
    }

    /// Create empty object
    pub fn empty_object() -> Self {
        Value::Object(Rc::new(RefCell::new(HashMap::new())))
    }

    /// Create empty array
    pub fn empty_array() -> Self {
        Value::Array(Rc::new(RefCell::new(Vec::new())))
    }
}

impl fmt::Display for Value {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.to_string())
    }
}

impl fmt::Debug for Value {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Value::Null => write!(f, "Null"),
            Value::Bool(b) => write!(f, "Bool({})", b),
            Value::Number(n) => write!(f, "Number({})", n),
            Value::String(s) => write!(f, "String(\"{}\")", s),
            Value::Array(_) => write!(f, "Array(...)"),
            Value::Object(_) => write!(f, "Object(...)"),
            Value::Function(_) => write!(f, "Function"),
            Value::Error(e) => write!(f, "Error(\"{}\")", e),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_type_of() {
        assert_eq!(Value::Null.type_of(), "null");
        assert_eq!(Value::Bool(true).type_of(), "boolean");
        assert_eq!(Value::Number(42.0).type_of(), "number");
        assert_eq!(Value::String("hello".to_string()).type_of(), "string");
    }

    #[test]
    fn test_is_truthy() {
        assert!(!Value::Null.is_truthy());
        assert!(!Value::Bool(false).is_truthy());
        assert!(Value::Bool(true).is_truthy());
        assert!(Value::Number(42.0).is_truthy());
        assert!(!Value::Number(0.0).is_truthy());
    }

    #[test]
    fn test_to_number() {
        assert_eq!(Value::Null.to_number(), 0.0);
        assert_eq!(Value::Bool(true).to_number(), 1.0);
        assert_eq!(Value::Bool(false).to_number(), 0.0);
        assert_eq!(Value::Number(42.0).to_number(), 42.0);
    }

    #[test]
    fn test_to_string() {
        assert_eq!(Value::Null.to_string(), "null");
        assert_eq!(Value::Bool(true).to_string(), "true");
        assert_eq!(Value::Number(42.0).to_string(), "42");
        assert_eq!(Value::String("hello".to_string()).to_string(), "hello");
    }
}
