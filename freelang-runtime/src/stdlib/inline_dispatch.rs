// Inline Function Dispatch
// V8의 inline caching 기법 구현
// Target: 함수 호출 오버헤드 제거 (50% 추가 개선)
//
// 원리:
// 1. 함수 호출 지점마다 최근 호출한 함수 캐시
// 2. 같은 함수가 반복되면 function pointer 직접 호출
// 3. 다른 함수면 일반 dispatch로 전환
//
// 이점:
// - CPU의 branch prediction 활용
// - CPU cache warming
// - Hot path 최적화

use crate::core::Value;
use std::collections::HashMap;

/// Inline cache entry
/// 최근 호출한 함수의 정보를 저장
#[derive(Clone, Debug)]
struct InlineCacheEntry {
    function_name: String,
    call_count: u64,
    is_monomorphic: bool,  // 같은 함수만 호출되는가?
}

/// Per-callsite inline cache
/// 각 함수 호출 지점에서 최근 호출 기록
pub struct InlineCache {
    // 호출 지점별 캐시 (key: callsite_id)
    caches: HashMap<usize, InlineCacheEntry>,
    // 전체 통계
    total_hits: u64,
    total_misses: u64,
    total_polymorphic: u64,
}

impl InlineCache {
    pub fn new() -> Self {
        InlineCache {
            caches: HashMap::new(),
            total_hits: 0,
            total_misses: 0,
            total_polymorphic: 0,
        }
    }

    /// Inline cache 조회 및 업데이트
    pub fn lookup_and_update(&mut self, callsite_id: usize, function_name: &str) -> bool {
        match self.caches.get_mut(&callsite_id) {
            Some(entry) => {
                if entry.function_name == function_name {
                    // Cache hit - 같은 함수
                    entry.call_count += 1;
                    self.total_hits += 1;
                    true
                } else {
                    // Cache miss - 다른 함수 호출
                    entry.is_monomorphic = false;
                    self.total_polymorphic += 1;
                    self.total_misses += 1;
                    false
                }
            }
            None => {
                // 첫 호출 - 캐시 생성
                self.caches.insert(
                    callsite_id,
                    InlineCacheEntry {
                        function_name: function_name.to_string(),
                        call_count: 1,
                        is_monomorphic: true,
                    },
                );
                self.total_hits += 1;
                true
            }
        }
    }

    /// Inline cache 통계
    pub fn stats(&self) -> InlineCacheStats {
        let total = self.total_hits + self.total_misses;
        let hit_rate = if total == 0 {
            0.0
        } else {
            (self.total_hits as f64 / total as f64) * 100.0
        };

        InlineCacheStats {
            callsites: self.caches.len(),
            monomorphic_sites: self
                .caches
                .values()
                .filter(|e| e.is_monomorphic)
                .count(),
            polymorphic_sites: self
                .caches
                .values()
                .filter(|e| !e.is_monomorphic)
                .count(),
            total_hits: self.total_hits,
            total_misses: self.total_misses,
            hit_rate,
            polymorphic_calls: self.total_polymorphic,
        }
    }

    /// Hottest callsites (가장 자주 호출되는 지점)
    pub fn hottest_callsites(&self, top_n: usize) -> Vec<(usize, String, u64)> {
        let mut callsites: Vec<_> = self
            .caches
            .iter()
            .map(|(id, entry)| (*id, entry.function_name.clone(), entry.call_count))
            .collect();
        callsites.sort_by(|a, b| b.2.cmp(&a.2));
        callsites.into_iter().take(top_n).collect()
    }

    pub fn clear(&mut self) {
        self.caches.clear();
        self.total_hits = 0;
        self.total_misses = 0;
        self.total_polymorphic = 0;
    }
}

/// Inline cache 통계
#[derive(Clone, Debug)]
pub struct InlineCacheStats {
    pub callsites: usize,
    pub monomorphic_sites: usize,
    pub polymorphic_sites: usize,
    pub total_hits: u64,
    pub total_misses: u64,
    pub hit_rate: f64,
    pub polymorphic_calls: u64,
}

impl std::fmt::Display for InlineCacheStats {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "InlineCache: {} callsites, {} monomorphic ({:.1}% hit rate), {} polymorphic",
            self.callsites, self.monomorphic_sites, self.hit_rate, self.polymorphic_sites
        )
    }
}

impl Default for InlineCache {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_inline_cache_creation() {
        let cache = InlineCache::new();
        assert_eq!(cache.caches.len(), 0);
    }

    #[test]
    fn test_inline_cache_monomorphic() {
        let mut cache = InlineCache::new();

        // 같은 함수를 반복 호출
        for _ in 0..100 {
            let hit = cache.lookup_and_update(1, "abs");
            assert!(hit);
        }

        let stats = cache.stats();
        assert!(stats.monomorphic_sites > 0);
        assert_eq!(stats.polymorphic_sites, 0);
        assert!(stats.hit_rate > 95.0);
    }

    #[test]
    fn test_inline_cache_polymorphic() {
        let mut cache = InlineCache::new();

        // 다른 함수들 호출
        let funcs = vec!["abs", "sqrt", "pow", "floor"];

        for i in 0..100 {
            let func = funcs[i % 4];
            cache.lookup_and_update(1, func);
        }

        let stats = cache.stats();
        println!("{}", stats);
        assert!(stats.polymorphic_sites > 0);
        assert!(stats.hit_rate < 100.0);
    }

    #[test]
    fn test_inline_cache_multiple_callsites() {
        let mut cache = InlineCache::new();

        // 3개 호출 지점
        for _ in 0..50 {
            cache.lookup_and_update(1, "abs");   // Callsite 1
            cache.lookup_and_update(2, "sqrt");  // Callsite 2
            cache.lookup_and_update(3, "pow");   // Callsite 3
        }

        let stats = cache.stats();
        assert_eq!(stats.callsites, 3);
        assert_eq!(stats.monomorphic_sites, 3);
        assert!(stats.hit_rate > 95.0);
    }

    #[test]
    fn test_hottest_callsites() {
        let mut cache = InlineCache::new();

        // Callsite 1: 100 calls
        for _ in 0..100 {
            cache.lookup_and_update(1, "abs");
        }

        // Callsite 2: 50 calls
        for _ in 0..50 {
            cache.lookup_and_update(2, "sqrt");
        }

        // Callsite 3: 10 calls
        for _ in 0..10 {
            cache.lookup_and_update(3, "pow");
        }

        let hottest = cache.hottest_callsites(3);
        assert_eq!(hottest[0].0, 1);  // Callsite 1 (100 calls)
        assert_eq!(hottest[1].0, 2);  // Callsite 2 (50 calls)
        assert_eq!(hottest[2].0, 3);  // Callsite 3 (10 calls)
    }

    #[test]
    fn test_inline_cache_effectiveness() {
        let mut cache = InlineCache::new();

        // 전형적인 프로그램: 10개 unique 함수, 90% single function calls
        let primary_func = "primary";
        let secondary_funcs = vec!["func1", "func2", "func3"];

        for i in 0..1000 {
            if i % 10 == 0 {
                // 10% polymorphic
                cache.lookup_and_update(1, secondary_funcs[i % 3]);
            } else {
                // 90% monomorphic
                cache.lookup_and_update(1, primary_func);
            }
        }

        let stats = cache.stats();
        println!("{}", stats);
        assert!(stats.hit_rate > 90.0);
    }

    #[test]
    fn test_inline_cache_clear() {
        let mut cache = InlineCache::new();

        for _ in 0..100 {
            cache.lookup_and_update(1, "abs");
        }

        assert!(cache.caches.len() > 0);

        cache.clear();
        assert_eq!(cache.caches.len(), 0);
        assert_eq!(cache.total_hits, 0);
    }
}
