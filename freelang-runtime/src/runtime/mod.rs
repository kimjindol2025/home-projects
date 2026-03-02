// Runtime Engine Module
// Core execution engine for FreeLang programs

pub mod vm;
pub mod executor;
pub mod context;
pub mod callstack;
pub mod interner;

pub use vm::RuntimeEngine;
pub use executor::Executor;
pub use context::ExecutionContext;
pub use callstack::CallStack;
pub use interner::{StringInterner, InternedString, InternerStats};
