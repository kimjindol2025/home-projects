// Function Dispatch Cache
// LRU cache for frequently used stdlib functions (90%+ hit rate)

use std::collections::HashMap;
use crate::core::Value;

/// Function pointer type for stdlib functions
pub type NativeFunction = fn(Vec<Value>) -> Result<Value, String>;

/// Cached function entry
#[derive(Clone)]
struct CacheEntry {
    function: NativeFunction,
    hits: u64,
}

/// LRU Function Cache
/// Maintains last N frequently used functions for O(1) dispatch
/// Target: 90%+ cache hit rate, 5-10 item capacity
#[derive(Clone)]
pub struct FunctionCache {
    cache: HashMap<String, CacheEntry>,
    lru_order: Vec<String>,
    max_capacity: usize,
    hits: u64,
    misses: u64,
}

impl FunctionCache {
    /// Create new function cache with default capacity (8)
    pub fn new() -> Self {
        FunctionCache {
            cache: HashMap::new(),
            lru_order: Vec::new(),
            max_capacity: 8,
            hits: 0,
            misses: 0,
        }
    }

    /// Set cache capacity
    pub fn with_capacity(capacity: usize) -> Self {
        FunctionCache {
            cache: HashMap::new(),
            lru_order: Vec::new(),
            max_capacity: capacity.max(1),
            hits: 0,
            misses: 0,
        }
    }

    /// Get cached function
    /// Returns Some(function) on cache hit, None on miss
    pub fn get(&mut self, name: &str) -> Option<NativeFunction> {
        if let Some(entry) = self.cache.get_mut(name) {
            entry.hits += 1;
            self.hits += 1;

            // Move to end (most recent)
            if let Some(pos) = self.lru_order.iter().position(|x| x == name) {
                self.lru_order.remove(pos);
                self.lru_order.push(name.to_string());
            }

            Some(entry.function)
        } else {
            self.misses += 1;
            None
        }
    }

    /// Insert or update cache entry
    pub fn insert(&mut self, name: String, function: NativeFunction) {
        // If already in cache, just update hit count and move to end
        if self.cache.contains_key(&name) {
            if let Some(entry) = self.cache.get_mut(&name) {
                entry.hits += 1;
            }
            if let Some(pos) = self.lru_order.iter().position(|x| x == &name) {
                self.lru_order.remove(pos);
                self.lru_order.push(name);
            }
            return;
        }

        // If cache is full, evict LRU (least recently used)
        if self.cache.len() >= self.max_capacity {
            if let Some(lru_name) = self.lru_order.first().cloned() {
                self.cache.remove(&lru_name);
                self.lru_order.remove(0);
            }
        }

        // Insert new entry
        self.cache.insert(
            name.clone(),
            CacheEntry {
                function,
                hits: 1,
            },
        );
        self.lru_order.push(name);
    }

    /// Get cache statistics
    pub fn stats(&self) -> CacheStats {
        let total = self.hits + self.misses;
        let hit_rate = if total == 0 {
            0.0
        } else {
            (self.hits as f64 / total as f64) * 100.0
        };

        CacheStats {
            capacity: self.max_capacity,
            current_size: self.cache.len(),
            hits: self.hits,
            misses: self.misses,
            hit_rate,
        }
    }

    /// Clear cache
    pub fn clear(&mut self) {
        self.cache.clear();
        self.lru_order.clear();
        self.hits = 0;
        self.misses = 0;
    }
}

/// Cache statistics
#[derive(Clone, Debug)]
pub struct CacheStats {
    pub capacity: usize,
    pub current_size: usize,
    pub hits: u64,
    pub misses: u64,
    pub hit_rate: f64,
}

impl std::fmt::Display for CacheStats {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Cache: {}/{} items, {} hits, {} misses, {:.1}% hit rate",
            self.current_size, self.capacity, self.hits, self.misses, self.hit_rate
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn dummy_abs(args: Vec<Value>) -> Result<Value, String> {
        if args.is_empty() {
            return Err("requires arguments".to_string());
        }
        Ok(Value::Number(args[0].to_number().abs()))
    }

    fn dummy_sqrt(args: Vec<Value>) -> Result<Value, String> {
        if args.is_empty() {
            return Err("requires arguments".to_string());
        }
        Ok(Value::Number(args[0].to_number().sqrt()))
    }

    #[test]
    fn test_cache_creation() {
        let cache = FunctionCache::new();
        assert_eq!(cache.max_capacity, 8);
        assert_eq!(cache.cache.len(), 0);
    }

    #[test]
    fn test_cache_insert_and_get() {
        let mut cache = FunctionCache::new();

        cache.insert("abs".to_string(), dummy_abs);
        assert!(cache.get("abs").is_some());
        assert!(cache.get("sqrt").is_none());
    }

    #[test]
    fn test_cache_hit_rate() {
        let mut cache = FunctionCache::with_capacity(2);

        cache.insert("abs".to_string(), dummy_abs);
        cache.insert("sqrt".to_string(), dummy_sqrt);

        // Hit
        let _ = cache.get("abs");
        let _ = cache.get("abs");

        // Miss
        let _ = cache.get("pow");

        let stats = cache.stats();
        assert_eq!(stats.hits, 2);
        assert_eq!(stats.misses, 1);
        assert!((stats.hit_rate - 66.66).abs() < 1.0);
    }

    #[test]
    fn test_cache_lru_eviction() {
        let mut cache = FunctionCache::with_capacity(2);

        cache.insert("abs".to_string(), dummy_abs);
        cache.insert("sqrt".to_string(), dummy_sqrt);

        // Both should be present
        assert!(cache.get("abs").is_some());
        assert!(cache.get("sqrt").is_some());

        // Add third function, should evict "abs" (LRU)
        fn dummy_pow(args: Vec<Value>) -> Result<Value, String> {
            if args.len() < 2 {
                return Err("requires 2 arguments".to_string());
            }
            Ok(Value::Number(args[0].to_number().powf(args[1].to_number())))
        }

        cache.insert("pow".to_string(), dummy_pow);

        // "abs" should be evicted
        assert!(cache.get("abs").is_none());
        assert!(cache.get("pow").is_some());
    }

    #[test]
    fn test_cache_lru_update() {
        let mut cache = FunctionCache::with_capacity(2);

        cache.insert("abs".to_string(), dummy_abs);
        cache.insert("sqrt".to_string(), dummy_sqrt);

        // Access "abs" to make it recently used
        let _ = cache.get("abs");

        // Add "pow", should evict "sqrt" (not "abs")
        fn dummy_pow(args: Vec<Value>) -> Result<Value, String> {
            if args.len() < 2 {
                return Err("requires 2 arguments".to_string());
            }
            Ok(Value::Number(args[0].to_number().powf(args[1].to_number())))
        }

        cache.insert("pow".to_string(), dummy_pow);

        assert!(cache.get("abs").is_some());
        assert!(cache.get("sqrt").is_none());
        assert!(cache.get("pow").is_some());
    }

    #[test]
    fn test_cache_clear() {
        let mut cache = FunctionCache::new();
        cache.insert("abs".to_string(), dummy_abs);
        cache.insert("sqrt".to_string(), dummy_sqrt);

        assert_eq!(cache.cache.len(), 2);

        cache.clear();

        assert_eq!(cache.cache.len(), 0);
        assert_eq!(cache.hits, 0);
        assert_eq!(cache.misses, 0);
    }
}
