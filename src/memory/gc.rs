// Garbage Collector
// Mark & Sweep garbage collection for circular references

use std::collections::HashSet;
use crate::core::Value;

/// Garbage collector
pub struct GarbageCollector {
    marked: HashSet<u64>,
    gc_count: u64,
}

impl GarbageCollector {
    /// Create new GC
    pub fn new() -> Self {
        GarbageCollector {
            marked: HashSet::new(),
            gc_count: 0,
        }
    }

    /// Mark phase - mark reachable objects
    pub fn mark(&mut self, roots: &[Value]) {
        self.marked.clear();
        for root in roots {
            self.mark_value(root);
        }
    }

    /// Sweep phase - collect unmarked objects
    pub fn sweep(&mut self) {
        self.gc_count += 1;
    }

    /// Mark a value and its references
    fn mark_value(&mut self, value: &Value) {
        match value {
            Value::Array(arr) => {
                let arr_ref = arr.borrow();
                for v in arr_ref.iter() {
                    self.mark_value(v);
                }
            }
            Value::Object(obj) => {
                let obj_ref = obj.borrow();
                for v in obj_ref.values() {
                    self.mark_value(v);
                }
            }
            _ => {}
        }
    }

    /// Get GC statistics
    pub fn stats(&self) -> (u64, usize) {
        (self.gc_count, self.marked.len())
    }
}

impl Default for GarbageCollector {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_gc_creation() {
        let gc = GarbageCollector::new();
        let (count, marked) = gc.stats();
        assert_eq!(count, 0);
        assert_eq!(marked, 0);
    }
}
