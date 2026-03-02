// Executor
// Executes individual instructions and functions

use crate::core::Value;

/// Execution statistics
#[derive(Clone, Debug)]
pub struct ExecutionStats {
    pub operations: u64,
    pub function_calls: u64,
    pub errors: u64,
}

/// Executor for FreeLang instructions
pub struct Executor {
    stats: ExecutionStats,
}

impl Executor {
    /// Create new executor
    pub fn new() -> Self {
        Executor {
            stats: ExecutionStats {
                operations: 0,
                function_calls: 0,
                errors: 0,
            },
        }
    }

    /// Execute a function call
    pub fn call_function(&mut self, func: Value, args: Vec<Value>) -> Result<Value, String> {
        self.stats.function_calls += 1;
        match func {
            Value::Function(_) => {
                // TODO: Implement function call
                Ok(Value::Null)
            }
            _ => {
                self.stats.errors += 1;
                Err(format!("Cannot call non-function: {}", func.type_of()))
            }
        }
    }

    /// Execute binary operation
    pub fn binary_op(&mut self, op: &str, left: Value, right: Value) -> Result<Value, String> {
        self.stats.operations += 1;
        match op {
            "+" => {
                let l = left.to_number();
                let r = right.to_number();
                Ok(Value::Number(l + r))
            }
            "-" => {
                let l = left.to_number();
                let r = right.to_number();
                Ok(Value::Number(l - r))
            }
            "*" => {
                let l = left.to_number();
                let r = right.to_number();
                Ok(Value::Number(l * r))
            }
            "/" => {
                let l = left.to_number();
                let r = right.to_number();
                if r == 0.0 {
                    self.stats.errors += 1;
                    Err("Division by zero".to_string())
                } else {
                    Ok(Value::Number(l / r))
                }
            }
            "%" => {
                let l = left.to_number();
                let r = right.to_number();
                if r == 0.0 {
                    self.stats.errors += 1;
                    Err("Modulo by zero".to_string())
                } else {
                    Ok(Value::Number(l % r))
                }
            }
            "==" => Ok(Value::Bool(left.equals(&right))),
            "!=" => Ok(Value::Bool(!left.equals(&right))),
            "<" => Ok(Value::Bool(left.to_number() < right.to_number())),
            ">" => Ok(Value::Bool(left.to_number() > right.to_number())),
            "<=" => Ok(Value::Bool(left.to_number() <= right.to_number())),
            ">=" => Ok(Value::Bool(left.to_number() >= right.to_number())),
            "&&" => Ok(Value::Bool(left.is_truthy() && right.is_truthy())),
            "||" => Ok(Value::Bool(left.is_truthy() || right.is_truthy())),
            _ => {
                self.stats.errors += 1;
                Err(format!("Unknown operator: {}", op))
            }
        }
    }

    /// Execute unary operation
    pub fn unary_op(&mut self, op: &str, operand: Value) -> Result<Value, String> {
        self.stats.operations += 1;
        match op {
            "-" => Ok(Value::Number(-operand.to_number())),
            "!" => Ok(Value::Bool(!operand.is_truthy())),
            "+" => Ok(Value::Number(operand.to_number())),
            _ => {
                self.stats.errors += 1;
                Err(format!("Unknown unary operator: {}", op))
            }
        }
    }

    /// Get execution statistics
    pub fn stats(&self) -> ExecutionStats {
        self.stats.clone()
    }

    /// Reset statistics
    pub fn reset_stats(&mut self) {
        self.stats = ExecutionStats {
            operations: 0,
            function_calls: 0,
            errors: 0,
        };
    }
}

impl Default for Executor {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_binary_add() {
        let mut exe = Executor::new();
        let result = exe.binary_op("+", Value::Number(5.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 8.0);
        assert_eq!(exe.stats().operations, 1);
    }

    #[test]
    fn test_binary_subtract() {
        let mut exe = Executor::new();
        let result = exe.binary_op("-", Value::Number(5.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 2.0);
    }

    #[test]
    fn test_binary_multiply() {
        let mut exe = Executor::new();
        let result = exe.binary_op("*", Value::Number(5.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 15.0);
    }

    #[test]
    fn test_binary_divide() {
        let mut exe = Executor::new();
        let result = exe.binary_op("/", Value::Number(15.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 5.0);
    }

    #[test]
    fn test_binary_div_by_zero() {
        let mut exe = Executor::new();
        let result = exe.binary_op("/", Value::Number(5.0), Value::Number(0.0));
        assert!(result.is_err());
        assert_eq!(exe.stats().errors, 1);
    }

    #[test]
    fn test_binary_modulo() {
        let mut exe = Executor::new();
        let result = exe.binary_op("%", Value::Number(10.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), 1.0);
    }

    #[test]
    fn test_comparison_equal() {
        let mut exe = Executor::new();
        let result = exe.binary_op("==", Value::Number(5.0), Value::Number(5.0));
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_comparison_not_equal() {
        let mut exe = Executor::new();
        let result = exe.binary_op("!=", Value::Number(5.0), Value::Number(3.0));
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_comparison_less() {
        let mut exe = Executor::new();
        let result = exe.binary_op("<", Value::Number(3.0), Value::Number(5.0));
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_logical_and() {
        let mut exe = Executor::new();
        let result = exe.binary_op("&&", Value::Bool(true), Value::Bool(true));
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());

        let result = exe.binary_op("&&", Value::Bool(true), Value::Bool(false));
        assert!(result.is_ok());
        assert!(!result.unwrap().is_truthy());
    }

    #[test]
    fn test_logical_or() {
        let mut exe = Executor::new();
        let result = exe.binary_op("||", Value::Bool(false), Value::Bool(true));
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());
    }

    #[test]
    fn test_unary_negate() {
        let mut exe = Executor::new();
        let result = exe.unary_op("-", Value::Number(5.0));
        assert!(result.is_ok());
        assert_eq!(result.unwrap().to_number(), -5.0);
    }

    #[test]
    fn test_unary_not() {
        let mut exe = Executor::new();
        let result = exe.unary_op("!", Value::Bool(true));
        assert!(result.is_ok());
        assert!(!result.unwrap().is_truthy());
    }

    #[test]
    fn test_stats_tracking() {
        let mut exe = Executor::new();
        exe.binary_op("+", Value::Number(1.0), Value::Number(2.0)).ok();
        exe.binary_op("-", Value::Number(1.0), Value::Number(2.0)).ok();
        exe.binary_op("/", Value::Number(1.0), Value::Number(0.0)).ok();

        let stats = exe.stats();
        assert_eq!(stats.operations, 3);
        assert_eq!(stats.errors, 1);
    }
}
