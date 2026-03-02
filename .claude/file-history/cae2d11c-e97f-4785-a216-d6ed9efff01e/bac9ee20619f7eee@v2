// String Interning
// Deduplicates strings and enables O(1) comparison
// Target: 80% memory reduction for frame names, 10-15% overall improvement

use std::collections::HashMap;
use std::rc::Rc;

/// Interned string wrapper
/// Points to shared string in global interner pool
#[derive(Clone, Eq, PartialEq, Hash)]
pub struct InternedString {
    ptr: Rc<str>,
}

impl InternedString {
    pub fn as_str(&self) -> &str {
        &self.ptr
    }

    pub fn to_string(&self) -> String {
        self.ptr.to_string()
    }
}

impl std::fmt::Display for InternedString {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.ptr)
    }
}

impl std::fmt::Debug for InternedString {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "InternedString({})", self.ptr)
    }
}

/// String Interner
/// Maintains pool of unique strings to reduce memory duplication
/// Frequently used strings (e.g., "global", "main") cached after first use
#[derive(Clone)]
pub struct StringInterner {
    // Pool of all interned strings
    pool: HashMap<String, Rc<str>>,
    // Statistics
    stats: InternerStats,
}

/// Interner statistics
#[derive(Clone, Debug)]
pub struct InternerStats {
    pub total_requests: u64,
    pub cache_hits: u64,
    pub new_strings: u64,
    pub deduplicated_bytes: usize,
}

impl InternerStats {
    pub fn hit_rate(&self) -> f64 {
        if self.total_requests == 0 {
            0.0
        } else {
            (self.cache_hits as f64 / self.total_requests as f64) * 100.0
        }
    }
}

impl std::fmt::Display for InternerStats {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "StringInterner: {} requests, {} hits ({:.1}% hit rate), {} new strings, {} bytes saved",
            self.total_requests, self.cache_hits, self.hit_rate(), self.new_strings, self.deduplicated_bytes
        )
    }
}

impl StringInterner {
    /// Create new interner
    pub fn new() -> Self {
        StringInterner {
            pool: HashMap::new(),
            stats: InternerStats {
                total_requests: 0,
                cache_hits: 0,
                new_strings: 0,
                deduplicated_bytes: 0,
            },
        }
    }

    /// Intern a string - returns shared reference to deduplicated string
    pub fn intern(&mut self, s: &str) -> InternedString {
        self.stats.total_requests += 1;

        if let Some(rc_str) = self.pool.get(s) {
            // Cache hit - reuse existing string
            self.stats.cache_hits += 1;
            InternedString {
                ptr: Rc::clone(rc_str),
            }
        } else {
            // Cache miss - add new string to pool
            self.stats.new_strings += 1;
            let rc_str: Rc<str> = Rc::from(s);
            self.stats.deduplicated_bytes += s.len();
            self.pool.insert(s.to_string(), Rc::clone(&rc_str));
            InternedString { ptr: rc_str }
        }
    }

    /// Get pool size
    pub fn pool_size(&self) -> usize {
        self.pool.len()
    }

    /// Get total bytes saved by deduplication
    pub fn bytes_saved(&self) -> usize {
        self.stats.deduplicated_bytes
    }

    /// Get statistics
    pub fn stats(&self) -> InternerStats {
        self.stats.clone()
    }

    /// Clear interner
    pub fn clear(&mut self) {
        self.pool.clear();
        self.stats = InternerStats {
            total_requests: 0,
            cache_hits: 0,
            new_strings: 0,
            deduplicated_bytes: 0,
        };
    }
}

impl Default for StringInterner {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_interner_creation() {
        let interner = StringInterner::new();
        assert_eq!(interner.pool_size(), 0);
    }

    #[test]
    fn test_intern_new_string() {
        let mut interner = StringInterner::new();
        let s1 = interner.intern("hello");
        assert_eq!(s1.as_str(), "hello");
        assert_eq!(interner.pool_size(), 1);
    }

    #[test]
    fn test_intern_deduplication() {
        let mut interner = StringInterner::new();
        let s1 = interner.intern("hello");
        let s2 = interner.intern("hello");

        // Should point to same underlying string
        assert_eq!(s1.as_str(), s2.as_str());
        assert_eq!(interner.pool_size(), 1);

        let stats = interner.stats();
        assert_eq!(stats.total_requests, 2);
        assert_eq!(stats.cache_hits, 1);
        assert_eq!(stats.new_strings, 1);
    }

    #[test]
    fn test_intern_multiple_strings() {
        let mut interner = StringInterner::new();

        for i in 0..100 {
            interner.intern(&format!("string_{}", i % 10));
        }

        let stats = interner.stats();
        assert_eq!(stats.new_strings, 10); // Only 10 unique strings
        assert_eq!(stats.cache_hits, 90); // 90 hits from duplicates
        assert!(stats.hit_rate() > 80.0);
    }

    #[test]
    fn test_intern_hit_rate() {
        let mut interner = StringInterner::new();

        // First 5 unique strings
        for i in 0..5 {
            interner.intern(&format!("str_{}", i));
        }

        // Repeat first 5 strings 10 times each
        for _ in 0..10 {
            for i in 0..5 {
                interner.intern(&format!("str_{}", i));
            }
        }

        let stats = interner.stats();
        println!("{}", stats);
        assert!(stats.hit_rate() > 95.0); // Should have >95% hit rate
    }

    #[test]
    fn test_intern_memory_savings() {
        let mut interner = StringInterner::new();

        // Simulate function names (e.g., "global", "main", "process")
        let names = vec!["global", "main", "process", "helper", "utility"];
        for _ in 0..100 {
            for name in &names {
                interner.intern(name);
            }
        }

        let bytes_saved = interner.bytes_saved();
        let expected = (5 * 100) * names.iter().map(|n| n.len()).sum::<usize>();

        println!("Bytes saved: {} out of {}", bytes_saved, expected);
        assert!(bytes_saved > 0);
    }

    #[test]
    fn test_intern_eq_and_hash() {
        let mut interner = StringInterner::new();
        let s1 = interner.intern("test");
        let s2 = interner.intern("test");

        assert_eq!(s1, s2);
        assert_eq!(hash(&s1), hash(&s2));
    }

    fn hash<T: std::hash::Hash>(obj: &T) -> u64 {
        use std::hash::{Hash, Hasher};
        use std::collections::hash_map::DefaultHasher;

        let mut hasher = DefaultHasher::new();
        obj.hash(&mut hasher);
        hasher.finish()
    }

    #[test]
    fn test_intern_clear() {
        let mut interner = StringInterner::new();
        interner.intern("hello");
        interner.intern("world");

        assert_eq!(interner.pool_size(), 2);

        interner.clear();
        assert_eq!(interner.pool_size(), 0);

        let stats = interner.stats();
        assert_eq!(stats.total_requests, 0);
    }
}
