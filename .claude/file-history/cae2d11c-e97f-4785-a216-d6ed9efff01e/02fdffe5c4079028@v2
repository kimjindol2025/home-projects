// Function Type System
// Represents different kinds of functions in FreeLang

use crate::core::Value;
use std::collections::HashMap;
use std::rc::Rc;
use std::cell::RefCell;

/// Function signature type
pub type NativeFunction = fn(Vec<Value>) -> Result<Value, String>;

/// Represents a user-defined function
#[derive(Clone)]
pub struct UserFunction {
    pub name: String,
    pub params: Vec<String>,
    pub body: Vec<Value>, // Instructions/expressions
    pub closure: Rc<RefCell<HashMap<String, Value>>>,
}

impl UserFunction {
    /// Create new user-defined function
    pub fn new(
        name: String,
        params: Vec<String>,
        body: Vec<Value>,
    ) -> Self {
        UserFunction {
            name,
            params,
            body,
            closure: Rc::new(RefCell::new(HashMap::new())),
        }
    }

    /// Get parameter count
    pub fn param_count(&self) -> usize {
        self.params.len()
    }

    /// Get function name
    pub fn get_name(&self) -> &str {
        &self.name
    }

    /// Get parameters
    pub fn get_params(&self) -> &[String] {
        &self.params
    }

    /// Set closure variable
    pub fn set_closure_var(&self, name: String, value: Value) {
        self.closure.borrow_mut().insert(name, value);
    }

    /// Get closure variable
    pub fn get_closure_var(&self, name: &str) -> Option<Value> {
        self.closure.borrow().get(name).cloned()
    }
}

impl std::fmt::Debug for UserFunction {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("UserFunction")
            .field("name", &self.name)
            .field("params", &self.params)
            .field("param_count", &self.param_count())
            .finish()
    }
}

impl std::fmt::Display for UserFunction {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "[Function: {}({})]", self.name, self.params.join(", "))
    }
}

/// Function kind enumeration
#[derive(Clone)]
pub enum FunctionKind {
    /// Native stdlib function
    Native(NativeFunction),
    /// User-defined function
    User(UserFunction),
    /// Closure (captures variables)
    Closure {
        func: UserFunction,
        captured: Rc<RefCell<HashMap<String, Value>>>,
    },
}

impl std::fmt::Debug for FunctionKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            FunctionKind::Native(_) => f.write_str("Native"),
            FunctionKind::User(uf) => write!(f, "User({})", uf.name),
            FunctionKind::Closure { func, .. } => write!(f, "Closure({})", func.name),
        }
    }
}

impl FunctionKind {
    /// Call the function
    pub fn call(&self, args: Vec<Value>) -> Result<Value, String> {
        match self {
            FunctionKind::Native(func) => func(args),
            FunctionKind::User(_) => Err("User-defined functions require a runtime context".to_string()),
            FunctionKind::Closure { .. } => Err("Closures require a runtime context".to_string()),
        }
    }

    /// Check if this is a native function
    pub fn is_native(&self) -> bool {
        matches!(self, FunctionKind::Native(_))
    }

    /// Check if this is a user-defined function
    pub fn is_user(&self) -> bool {
        matches!(self, FunctionKind::User(_))
    }

    /// Check if this is a closure
    pub fn is_closure(&self) -> bool {
        matches!(self, FunctionKind::Closure { .. })
    }

    /// Get function name if available
    pub fn name(&self) -> Option<&str> {
        match self {
            FunctionKind::User(uf) => Some(&uf.name),
            FunctionKind::Closure { func, .. } => Some(&func.name),
            FunctionKind::Native(_) => None,
        }
    }

    /// Get parameter count if available
    pub fn param_count(&self) -> Option<usize> {
        match self {
            FunctionKind::User(uf) => Some(uf.param_count()),
            FunctionKind::Closure { func, .. } => Some(func.param_count()),
            FunctionKind::Native(_) => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_user_function_creation() {
        let func = UserFunction::new(
            "add".to_string(),
            vec!["x".to_string(), "y".to_string()],
            vec![],
        );
        assert_eq!(func.get_name(), "add");
        assert_eq!(func.param_count(), 2);
    }

    #[test]
    fn test_closure_variables() {
        let func = UserFunction::new(
            "multiply".to_string(),
            vec!["x".to_string()],
            vec![],
        );
        func.set_closure_var("factor".to_string(), Value::Number(2.0));
        let var = func.get_closure_var("factor");
        assert!(var.is_some());
        assert_eq!(var.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_native_function() {
        fn add(args: Vec<Value>) -> Result<Value, String> {
            if args.len() != 2 {
                return Err("add requires 2 arguments".to_string());
            }
            Ok(Value::Number(args[0].to_number() + args[1].to_number()))
        }

        let func_kind = FunctionKind::Native(add);
        assert!(func_kind.is_native());
        assert!(!func_kind.is_user());

        let result = func_kind.call(vec![Value::Number(2.0), Value::Number(3.0)]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_function_kind_debug() {
        let user_func = UserFunction::new(
            "test".to_string(),
            vec!["x".to_string()],
            vec![],
        );
        let kind = FunctionKind::User(user_func);
        assert_eq!(format!("{:?}", kind), "User(test)");
    }

    #[test]
    fn test_function_kind_name() {
        let user_func = UserFunction::new(
            "multiply".to_string(),
            vec!["a".to_string(), "b".to_string()],
            vec![],
        );
        let kind = FunctionKind::User(user_func);
        assert_eq!(kind.name(), Some("multiply"));
        assert_eq!(kind.param_count(), Some(2));
    }
}
