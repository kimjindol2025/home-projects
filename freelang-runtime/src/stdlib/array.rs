// Array Functions (12)
// Array manipulation and operations

use crate::core::Value;

/// Get array length
pub fn length(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("length requires 1 argument".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            Ok(Value::Number(arr.borrow().len() as f64))
        }
        _ => Err("Argument must be an array".to_string())
    }
}

/// Push element to array
pub fn push(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("push requires 2 arguments: array, value".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            arr.borrow_mut().push(args[1].clone());
            Ok(Value::Number(arr.borrow().len() as f64))
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Pop element from array
pub fn pop(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("pop requires 1 argument".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            match arr.borrow_mut().pop() {
                Some(value) => Ok(value),
                None => Ok(Value::Null)
            }
        }
        _ => Err("Argument must be an array".to_string())
    }
}

/// Get element at index
pub fn get(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("get requires 2 arguments: array, index".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let index = args[1].to_number() as usize;
            let arr_ref = arr.borrow();
            if index < arr_ref.len() {
                Ok(arr_ref[index].clone())
            } else {
                Ok(Value::Null)
            }
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Set element at index
pub fn set(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 3 {
        return Err("set requires 3 arguments: array, index, value".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let index = args[1].to_number() as usize;
            let mut arr_mut = arr.borrow_mut();
            if index < arr_mut.len() {
                arr_mut[index] = args[2].clone();
                Ok(Value::Null)
            } else {
                Err("Index out of bounds".to_string())
            }
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Slice array from start to end
pub fn slice(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 3 {
        return Err("slice requires 3 arguments: array, start, end".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let start = args[1].to_number() as usize;
            let end = args[2].to_number() as usize;
            let arr_ref = arr.borrow();

            if start >= arr_ref.len() || end > arr_ref.len() || start > end {
                return Err("Invalid slice range".to_string());
            }

            let sliced: Vec<Value> = arr_ref[start..end].to_vec();
            Ok(Value::array(sliced))
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Find index of element
pub fn index_of(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("index_of requires 2 arguments: array, value".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let arr_ref = arr.borrow();
            for (i, item) in arr_ref.iter().enumerate() {
                if item.equals(&args[1]) {
                    return Ok(Value::Number(i as f64));
                }
            }
            Ok(Value::Number(-1.0))
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Check if array contains value
pub fn contains(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("contains requires 2 arguments: array, value".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let arr_ref = arr.borrow();
            for item in arr_ref.iter() {
                if item.equals(&args[1]) {
                    return Ok(Value::Bool(true));
                }
            }
            Ok(Value::Bool(false))
        }
        _ => Err("First argument must be an array".to_string())
    }
}

/// Reverse array in place
pub fn reverse(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("reverse requires 1 argument".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            arr.borrow_mut().reverse();
            Ok(Value::Null)
        }
        _ => Err("Argument must be an array".to_string())
    }
}

/// Sort array (numeric/lexical)
pub fn sort(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("sort requires 1 argument".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            arr.borrow_mut().sort_by(|a, b| {
                a.to_number()
                    .partial_cmp(&b.to_number())
                    .unwrap_or(std::cmp::Ordering::Equal)
            });
            Ok(Value::Null)
        }
        _ => Err("Argument must be an array".to_string())
    }
}

/// Join array elements into string
pub fn join(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("join requires 2 arguments: array, delimiter".to_string());
    }

    match &args[0] {
        Value::Array(arr) => {
            let delimiter = args[1].to_string();
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_length() {
        let arr = Value::array(vec![
            Value::Number(1.0),
            Value::Number(2.0),
            Value::Number(3.0),
        ]);
        let result = length(vec![arr]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 3.0);
    }

    #[test]
    fn test_push() {
        let arr = Value::array(vec![Value::Number(1.0)]);
        let result = push(vec![arr.clone(), Value::Number(2.0)]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_pop() {
        let arr = Value::array(vec![Value::Number(1.0), Value::Number(2.0)]);
        let result = pop(vec![arr.clone()]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_get() {
        let arr = Value::array(vec![
            Value::Number(10.0),
            Value::Number(20.0),
            Value::Number(30.0),
        ]);
        let result = get(vec![arr, Value::Number(1.0)]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 20.0);
    }

    #[test]
    fn test_set() {
        let arr = Value::array(vec![Value::Number(1.0), Value::Number(2.0)]);
        let result = set(vec![arr.clone(), Value::Number(0.0), Value::Number(10.0)]);
        assert!(result.is_ok());

        // Verify the change
        let val = get(vec![arr, Value::Number(0.0)]);
        assert_eq!(val.unwrap().to_number(), 10.0);
    }

    #[test]
    fn test_index_of() {
        let arr = Value::array(vec![
            Value::Number(10.0),
            Value::Number(20.0),
            Value::Number(30.0),
        ]);
        let result = index_of(vec![arr, Value::Number(20.0)]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 1.0);
    }

    #[test]
    fn test_contains() {
        let arr = Value::array(vec![Value::Number(1.0), Value::Number(2.0)]);
        let result = contains(vec![arr, Value::Number(1.0)]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_reverse() {
        let arr = Value::array(vec![
            Value::Number(1.0),
            Value::Number(2.0),
            Value::Number(3.0),
        ]);
        let result = reverse(vec![arr.clone()]);
        assert!(result.is_ok());

        // Verify reversal
        let val = get(vec![arr.clone(), Value::Number(0.0)]);
        assert_eq!(val.unwrap().to_number(), 3.0);
    }

    #[test]
    fn test_sort() {
        let arr = Value::array(vec![
            Value::Number(3.0),
            Value::Number(1.0),
            Value::Number(2.0),
        ]);
        let result = sort(vec![arr.clone()]);
        assert!(result.is_ok());

        // Verify sorting
        let val0 = get(vec![arr.clone(), Value::Number(0.0)]);
        let val1 = get(vec![arr.clone(), Value::Number(1.0)]);
        let val2 = get(vec![arr, Value::Number(2.0)]);

        assert_eq!(val0.unwrap().to_number(), 1.0);
        assert_eq!(val1.unwrap().to_number(), 2.0);
        assert_eq!(val2.unwrap().to_number(), 3.0);
    }

    #[test]
    fn test_join() {
        let arr = Value::array(vec![
            Value::String("a".to_string()),
            Value::String("b".to_string()),
            Value::String("c".to_string()),
        ]);
        let result = join(vec![arr, Value::String(",".to_string())]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_string(), "a,b,c");
    }
}
