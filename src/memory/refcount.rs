// Reference Counting
// Automatic reference counting for immediate memory cleanup

use std::collections::HashMap;

/// Reference counter
pub struct ReferenceCounter {
    counts: HashMap<u64, usize>,
}

impl ReferenceCounter {
    /// Create new reference counter
    pub fn new() -> Self {
        ReferenceCounter {
            counts: HashMap::new(),
        }
    }

    /// Increment reference count
    pub fn increment(&mut self, id: u64) {
        *self.counts.entry(id).or_insert(0) += 1;
    }

    /// Decrement reference count
    pub fn decrement(&mut self, id: u64) -> usize {
        let count = self.counts.entry(id).or_insert(0);
        if *count > 0 {
            *count -= 1;
        }
        *count
    }

    /// Get reference count
    pub fn count(&self, id: u64) -> usize {
        self.counts.get(&id).copied().unwrap_or(0)
    }

    /// Check if no references
    pub fn is_unreferenced(&self, id: u64) -> bool {
        self.count(id) == 0
    }

    /// Remove reference count
    pub fn remove(&mut self, id: u64) -> bool {
        self.counts.remove(&id).is_some()
    }
}

impl Default for ReferenceCounter {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_refcount() {
        let mut rc = ReferenceCounter::new();
        rc.increment(1);
        assert_eq!(rc.count(1), 1);
        rc.increment(1);
        assert_eq!(rc.count(1), 2);
        rc.decrement(1);
        assert_eq!(rc.count(1), 1);
    }
}
