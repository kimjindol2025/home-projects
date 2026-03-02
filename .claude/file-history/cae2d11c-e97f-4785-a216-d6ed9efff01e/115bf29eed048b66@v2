// Memory Allocator
// Manages heap allocation and deallocation

use std::collections::HashMap;
use crate::core::Value;

/// Memory allocator with GC support
pub struct MemoryAllocator {
    heap: HashMap<u64, Value>,
    next_id: u64,
    total_allocated: usize,
    gc_threshold: usize,  // Trigger GC at 5MB
}

impl MemoryAllocator {
    /// Create new allocator
    pub fn new() -> Self {
        MemoryAllocator {
            heap: HashMap::new(),
            next_id: 1,
            total_allocated: 0,
            gc_threshold: 5 * 1024 * 1024,  // 5MB
        }
    }

    /// Allocate value on heap
    pub fn allocate(&mut self, value: Value) -> u64 {
        let id = self.next_id;
        self.next_id += 1;

        // Estimate size
        let size = self.estimate_size(&value);
        self.total_allocated += size;

        self.heap.insert(id, value);

        id
    }

    /// Free allocated value
    pub fn free(&mut self, id: u64) -> Option<Value> {
        self.heap.remove(&id)
    }

    /// Get value from heap
    pub fn get(&self, id: u64) -> Option<&Value> {
        self.heap.get(&id)
    }

    /// Check if GC should run
    pub fn should_gc(&self) -> bool {
        self.total_allocated > self.gc_threshold
    }

    /// Get memory usage
    pub fn memory_usage(&self) -> usize {
        self.total_allocated
    }

    /// Estimate value size in bytes
    fn estimate_size(&self, value: &Value) -> usize {
        match value {
            Value::Null => 8,
            Value::Bool(_) => 8,
            Value::Number(_) => 8,
            Value::String(s) => 8 + s.len(),
            Value::Array(_) => 24,  // Approximate
            Value::Object(_) => 24, // Approximate
            Value::Function(_) => 16,
            Value::Error(e) => 8 + e.len(),
        }
    }

    /// Reset allocator
    pub fn reset(&mut self) {
        self.heap.clear();
        self.next_id = 1;
        self.total_allocated = 0;
    }
}

impl Default for MemoryAllocator {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_allocate() {
        let mut alloc = MemoryAllocator::new();
        let id = alloc.allocate(Value::Number(42.0));
        assert!(alloc.get(id).is_some());
    }

    #[test]
    fn test_free() {
        let mut alloc = MemoryAllocator::new();
        let id = alloc.allocate(Value::String("test".to_string()));
        let freed = alloc.free(id);
        assert!(freed.is_some());
        assert!(alloc.get(id).is_none());
    }
}
