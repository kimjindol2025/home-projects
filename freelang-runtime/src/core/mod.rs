// Core Types Module
// Defines fundamental types used throughout the runtime

pub mod value;
pub mod object;
pub mod traits;
pub mod function;
pub mod compact_value;

pub use value::Value;
pub use object::Object;
pub use traits::{Callable, Displayable};
pub use function::{FunctionKind, UserFunction, NativeFunction};
pub use compact_value::ValueSize;
