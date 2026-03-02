// Generational Garbage Collector
// V8, JVM 같은 현대 언어 엔진의 GC 기법
// Target: 50% 더 적은 full collection pause time
//
// Theory: 대부분의 객체는 young generation에서 빠르게 소멸함
// (weak generational hypothesis)
//
// Strategy:
// 1. Young generation: 빈번하고 빠른 collection
// 2. Old generation: 드물게, 비용이 큼
// 3. Card marking: Old → Young 참조 추적

use std::collections::{HashMap, HashSet};
use crate::core::Value;

/// 객체의 나이 정보
#[derive(Clone, Copy, Debug, PartialEq, Eq, PartialOrd, Ord)]
pub enum Generation {
    Young,  // 새로 할당된 객체 (0-2 collections)
    Old,    // 여러 collections를 생존한 객체
}

/// Generational GC 통계
#[derive(Clone, Debug)]
pub struct GenerationalGCStats {
    pub young_collections: u64,
    pub old_collections: u64,
    pub young_pause_ms: f64,
    pub old_pause_ms: f64,
    pub young_collected: usize,
    pub old_collected: usize,
    pub promotion_count: u64,
}

impl GenerationalGCStats {
    pub fn new() -> Self {
        GenerationalGCStats {
            young_collections: 0,
            old_collections: 0,
            young_pause_ms: 0.0,
            old_pause_ms: 0.0,
            young_collected: 0,
            old_collected: 0,
            promotion_count: 0,
        }
    }

    /// Young collection은 old보다 훨씬 빠름 (10-100배)
    pub fn is_efficient(&self) -> bool {
        if self.old_collections == 0 {
            return true;
        }
        // Old collection은 최소 5배 느려야 함 (예상)
        let avg_old = self.old_pause_ms / self.old_collections as f64;
        let avg_young = if self.young_collections == 0 {
            0.0
        } else {
            self.young_pause_ms / self.young_collections as f64
        };

        avg_old > avg_young * 5.0 || avg_old > 50.0
    }
}

impl std::fmt::Display for GenerationalGCStats {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "GenerationalGC: {} young ({:.2}ms total), {} old ({:.2}ms total), {} promotions",
            self.young_collections,
            self.young_pause_ms,
            self.old_collections,
            self.old_pause_ms,
            self.promotion_count
        )
    }
}

/// Generational GC 구현
pub struct GenerationalGC {
    young: HashSet<usize>,  // Young generation object ids
    old: HashSet<usize>,    // Old generation object ids

    // Card marking: which old objects reference young objects
    // (Old → Young 참조를 빠르게 찾기 위함)
    card_marks: HashMap<usize, bool>,

    stats: GenerationalGCStats,

    // 임계값
    young_threshold: usize,  // Young collection 트리거
    promotion_age: u32,      // 이 나이가 되면 old로 승격
}

impl GenerationalGC {
    pub fn new() -> Self {
        GenerationalGC {
            young: HashSet::new(),
            old: HashSet::new(),
            card_marks: HashMap::new(),
            stats: GenerationalGCStats::new(),
            young_threshold: 100,    // 100개 young 객체 도달 시 collection
            promotion_age: 2,         // 2번 collection 생존 시 promotion
        }
    }

    /// 새 객체를 young generation에 할당
    pub fn allocate(&mut self, id: usize) {
        self.young.insert(id);
    }

    /// Young generation GC (빠름, 자주 실행)
    pub fn collect_young(&mut self, live_young: HashSet<usize>) -> usize {
        let start = std::time::Instant::now();

        // 죽은 young 객체 제거
        let before = self.young.len();
        self.young = live_young.clone();
        let collected = before - self.young.len();

        // 일부 객체를 old로 승격
        let to_promote: Vec<_> = self.young.iter().copied().collect();
        for obj_id in to_promote {
            if self.card_marks.get(&obj_id).copied().unwrap_or(false) {
                // Old → Young 참조가 없으면 promotion 가능
                self.young.remove(&obj_id);
                self.old.insert(obj_id);
                self.stats.promotion_count += 1;
            }
        }

        let elapsed = start.elapsed().as_secs_f64() * 1000.0;
        self.stats.young_collections += 1;
        self.stats.young_pause_ms += elapsed;
        self.stats.young_collected += collected;

        collected
    }

    /// Old generation GC (느림, 드물게 실행)
    pub fn collect_old(&mut self, live_old: HashSet<usize>) -> usize {
        let start = std::time::Instant::now();

        let before = self.old.len();
        self.old = live_old;
        let collected = before - self.old.len();

        // Card marks 정리
        self.card_marks.clear();

        let elapsed = start.elapsed().as_secs_f64() * 1000.0;
        self.stats.old_collections += 1;
        self.stats.old_pause_ms += elapsed;
        self.stats.old_collected += collected;

        collected
    }

    /// Old 객체가 young 객체를 참조할 때 호출
    pub fn mark_card(&mut self, old_id: usize, young_id: usize) {
        self.card_marks.insert(old_id, true);
    }

    /// Young collection 필요 여부
    pub fn should_collect_young(&self) -> bool {
        self.young.len() >= self.young_threshold
    }

    /// Old collection 필요 여부
    pub fn should_collect_old(&self) -> bool {
        self.old.len() > self.young_threshold * 5  // Young의 5배
    }

    /// 통계 반환
    pub fn stats(&self) -> &GenerationalGCStats {
        &self.stats
    }

    /// Young 객체 수
    pub fn young_count(&self) -> usize {
        self.young.len()
    }

    /// Old 객체 수
    pub fn old_count(&self) -> usize {
        self.old.len()
    }
}

impl Default for GenerationalGC {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_generational_gc_creation() {
        let gc = GenerationalGC::new();
        assert_eq!(gc.young_count(), 0);
        assert_eq!(gc.old_count(), 0);
    }

    #[test]
    fn test_allocate_in_young() {
        let mut gc = GenerationalGC::new();

        gc.allocate(1);
        gc.allocate(2);
        gc.allocate(3);

        assert_eq!(gc.young_count(), 3);
        assert_eq!(gc.old_count(), 0);
    }

    #[test]
    fn test_young_collection() {
        let mut gc = GenerationalGC::new();

        // Young에 5개 객체 할당
        for i in 0..5 {
            gc.allocate(i);
        }

        assert_eq!(gc.young_count(), 5);

        // 2개만 살아있음
        let live = vec![1, 3].into_iter().collect();
        let collected = gc.collect_young(live);

        assert_eq!(collected, 3);  // 3개 수집됨
        assert_eq!(gc.young_count(), 2);
    }

    #[test]
    fn test_promotion() {
        let mut gc = GenerationalGC::new();

        // Young에 할당
        for i in 0..10 {
            gc.allocate(i);
        }

        // Card mark - old → young 참조 없음
        let live = (0..10).collect();
        gc.collect_young(live);

        // Promotion이 발생했을 것
        assert!(gc.stats().promotion_count > 0);
    }

    #[test]
    fn test_should_collect_young() {
        let mut gc = GenerationalGC::new();

        for i in 0..100 {
            gc.allocate(i);
        }

        assert!(gc.should_collect_young());
    }

    #[test]
    fn test_generational_gc_pause_times() {
        let mut gc = GenerationalGC::new();

        // Young 할당 및 collection
        for i in 0..50 {
            gc.allocate(i);
        }

        let live = (0..50).into_iter().collect();
        gc.collect_young(live);

        let stats = gc.stats();
        println!("{}", stats);
        assert!(stats.young_collections > 0);
    }

    #[test]
    fn test_card_marking() {
        let mut gc = GenerationalGC::new();

        // Old → Young 참조 표시
        gc.mark_card(100, 1);
        gc.mark_card(101, 2);

        // Card marks가 있음
        assert!(gc.card_marks.contains_key(&100));
        assert!(gc.card_marks.contains_key(&101));
    }
}
