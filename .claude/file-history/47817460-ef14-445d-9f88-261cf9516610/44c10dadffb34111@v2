# Phase 7: Adaptive Resource Allocation - 완전 설계 문서

## 프로젝트 개요

**목표**: Phase 6의 예측 결과를 바탕으로 CPU, 메모리, I/O 리소스를 동적으로 할당하여 시스템 성능 최적화

**핵심 아이디어**:
```
Phase 6 (예측) → "메모리 누수 감지, 2초 후 임계값 도달"
                 ↓
Phase 7 (할당) → "즉시 메모리 사전 할당, GC 우선순위 상향"
                 ↓
             결과: 임계값 도달 0, 장애 방지
```

**기간**: 3주 (Days 1-21)
**목표 줄수**: 2,500~3,000줄
**목표 테스트**: 70+ 개

---

## 📋 Week-by-Week Breakdown

### Week 1: CPU Scheduling Optimization (Days 1-7)

**목표**: 예측된 패턴에 따라 CPU 스케줄링 동적 조정

#### Day 1-2: CPU Resource Analyzer (200줄)
```rust
pub struct CPUMetrics {
    node_id: u32,
    utilization: f64,              // 0.0-1.0
    context_switches: u64,
    runqueue_length: u32,
    cache_misses: u64,
    thermal_level: f64,            // 0.0-1.0 (온도)
}

pub struct CPUResourceAnalyzer {
    current_metrics: CPUMetrics,
    history: Vec<(u64, CPUMetrics)>,
    max_history: usize,            // 100개 보유
}

impl CPUResourceAnalyzer {
    pub fn detect_cpu_bottleneck(&self) -> Option<f64>;      // 병목 신뢰도
    pub fn predict_saturation_trend(&self) -> TrendType;     // 상승/안정/하강
    pub fn get_cache_efficiency(&self) -> f64;               // 0.0-1.0
    pub fn estimate_needed_cores(&self, pattern: &Pattern) -> u32; // 필요 코어 수
}
```

**테스트** (8개):
- test_cpu_utilization_detection
- test_context_switch_overhead
- test_cache_miss_rate_analysis
- test_thermal_monitoring
- test_cpu_trend_prediction
- test_core_requirement_estimation
- test_multi_socket_awareness
- test_cpu_history_retention

#### Day 3-4: Thread Affinity Manager (220줄)
```rust
pub struct ThreadAffinity {
    thread_id: u32,
    preferred_cores: Vec<u32>,     // NUMA-aware
    cache_level: u8,               // L1/L2/L3
    migration_cost: u32,           // 마이그레이션 비용(ns)
}

pub struct ThreadAffinityManager {
    affinities: HashMap<u32, ThreadAffinity>,
    core_topology: CPUTopology,    // NUMA, SMT 정보
    rebalance_threshold: f64,
}

impl ThreadAffinityManager {
    pub fn optimize_affinity(&mut self, pattern: &Pattern) -> bool;
    pub fn detect_cache_migration_problem(&self) -> Option<(u32, u32)>; // (src, dst)
    pub fn get_optimal_core_assignment(&self, thread_id: u32) -> u32;
    pub fn measure_migration_cost(&self, from_core: u32, to_core: u32) -> u32;
}
```

**테스트** (8개):
- test_affinity_assignment_basic
- test_numa_locality_optimization
- test_cache_level_awareness
- test_migration_cost_calculation
- test_hotspot_detection
- test_smt_thread_pairing
- test_cross_socket_minimization
- test_affinity_rebalancing

#### Day 5-7: Dynamic Priority Scheduler (260줄)
```rust
pub enum ThreadPriority {
    Critical,              // 0 (OS 스레드)
    HighRealtime,         // 1-49 (실시간)
    NormalInteractive,    // 50-99 (상호작용)
    NormalBatch,          // 100-139 (배치)
    IdleNice,             // 140-199 (백그라운드)
}

pub struct ThreadScheduleEntry {
    thread_id: u32,
    priority: ThreadPriority,
    time_slice_ms: u32,
    boost_until: u64,     // 우선순위 부스트 기한
}

pub struct DynamicPriorityScheduler {
    threads: HashMap<u32, ThreadScheduleEntry>,
    prediction_boost: bool,        // Phase 6 기반 부스트
    predictive_actions: Vec<ProactiveAction>,
}

impl DynamicPriorityScheduler {
    pub fn adjust_priority_from_prediction(&mut self, pattern: &Pattern, action: &ProactiveAction);
    pub fn calculate_time_slice(&self, priority: &ThreadPriority) -> u32;
    pub fn boost_memory_recovery_threads(&mut self, pattern: &Pattern);
    pub fn apply_cpu_affinity_hint(&mut self, thread_id: u32, hint: &ThreadAffinity);
    pub fn get_schedule_fairness_score(&self) -> f64; // 0.0-1.0
}
```

**테스트** (8개):
- test_priority_level_assignment
- test_time_slice_calculation
- test_boost_from_pattern_prediction
- test_io_wait_handling
- test_starvation_prevention
- test_priority_inversion_detection
- test_schedule_fairness_metric
- test_preemption_on_critical_event

#### Integration: CPUScheduler Week1 System (120줄)

**공개 API**:
```rust
pub struct Week1CPUScheduler {
    analyzer: CPUResourceAnalyzer,
    affinity: ThreadAffinityManager,
    scheduler: DynamicPriorityScheduler,
}

impl Week1CPUScheduler {
    pub fn adapt_to_pattern(&mut self, pattern: &Pattern, action: &ProactiveAction) -> ScheduleAdjustment;
    pub fn get_current_schedule_efficiency(&self) -> f64;
    pub fn forecast_cpu_saturation(&self) -> (f64, u32); // (saturation, time_to_critical_ms)
}
```

**Week 1 통계**:
- 파일 개수: 4개
- 총 줄수: 800줄
- 테스트: 32개
- 핵심 개념: CPU Affinity, Thread Priority, NUMA Awareness

---

### Week 2: Memory Allocation Strategy (Days 8-14)

**목표**: 예측된 메모리 패턴에 따라 사전 할당 및 압축 전략 실행

#### Day 8-9: Memory Predictor (240줄)
```rust
pub struct MemoryPrediction {
    current_usage_mb: u64,
    predicted_peak_mb: u64,        // 앞으로 N초 내 최고 사용량
    growth_rate_mb_per_sec: f64,
    time_to_limit_sec: u32,
    confidence: f64,               // Phase 6 신뢰도
}

pub struct MemoryPredictor {
    historical_patterns: Vec<MemoryPattern>,
    phase6_predictions: Vec<PredictionOutcome>,
}

impl MemoryPredictor {
    pub fn predict_memory_usage(&self, pattern: &Pattern) -> MemoryPrediction;
    pub fn estimate_preallocation_size(&self, pattern: &Pattern, safety_margin: f64) -> u64;
    pub fn detect_memory_leak_pattern(&self) -> Option<(u64, f64)>; // (growth_rate, confidence)
    pub fn calculate_gc_urgency(&self, pattern: &Pattern) -> f64; // 0.0-1.0
}
```

**테스트** (8개):
- test_memory_prediction_accuracy
- test_preallocation_sizing
- test_leak_detection_confidence
- test_growth_rate_calculation
- test_gc_urgency_scoring
- test_prediction_with_phase6_input
- test_memory_fragmentation_prediction
- test_peak_usage_forecasting

#### Day 10-11: Predictive Memory Allocator (260줄)
```rust
pub struct PreallocatedPool {
    base_address: u64,
    total_size: u64,
    allocated: u64,
    reserved_for_pattern: HashMap<String, u64>,  // 패턴별 예약
    expiration_time: u64,                        // 예약 만료 시간
}

pub struct PredictiveMemoryAllocator {
    pools: Vec<PreallocatedPool>,
    phase6_predictions: Vec<MemoryPrediction>,
    fragmentation_threshold: f64,
}

impl PredictiveMemoryAllocator {
    pub fn preallocate_from_prediction(&mut self, pattern: &Pattern) -> bool;
    pub fn compact_memory(&mut self, target_mb: u64) -> u64; // 해제된 바이트
    pub fn adjust_pool_size(&mut self, pattern: &Pattern);
    pub fn get_fragmentation_ratio(&self) -> f64;            // 0.0-1.0
    pub fn suggest_pageout_candidates(&self) -> Vec<PageInfo>;
}
```

**테스트** (8개):
- test_pool_preallocation
- test_pattern_specific_reservation
- test_compaction_effectiveness
- test_fragmentation_reduction
- test_pool_expiration
- test_emergency_allocation
- test_memory_pressure_response
- test_swap_decision_making

#### Day 12-14: Memory Pressure Handler (240줄)
```rust
pub enum MemoryPressureLevel {
    Normal,                // < 60%
    Warning,               // 60-80%
    Critical,              // 80-95%
    Emergency,             // 95-100%
}

pub struct MemoryPressureResponse {
    actions: Vec<MemoryAction>,
    estimated_freed_mb: u64,
    time_required_ms: u32,
}

pub enum MemoryAction {
    CompactMemory,
    DropPageCache,
    ShrinkHeap,
    TriggerGC,
    PageOutToSwap,
    KillNonEssentialProcess,
}

pub struct MemoryPressureHandler {
    threshold_by_level: HashMap<MemoryPressureLevel, f64>,
    action_priority: Vec<MemoryAction>,
}

impl MemoryPressureHandler {
    pub fn handle_memory_pressure(&mut self, current_usage: u64, limit: u64) -> MemoryPressureResponse;
    pub fn execute_preemptive_action(&mut self, pattern: &Pattern) -> bool;
    pub fn estimate_action_effectiveness(&self, action: &MemoryAction) -> u64; // 예상 해제(MB)
    pub fn get_pressure_level(&self, usage: u64, limit: u64) -> MemoryPressureLevel;
}
```

**테스트** (8개):
- test_pressure_level_detection
- test_response_action_ordering
- test_preemptive_gc_triggering
- test_pageout_decision
- test_emergency_response
- test_pressure_prediction
- test_recovery_time_estimation
- test_repeated_pressure_handling

#### Integration: MemoryAllocator Week2 System (150줄)

**공개 API**:
```rust
pub struct Week2MemoryAllocator {
    predictor: MemoryPredictor,
    allocator: PredictiveMemoryAllocator,
    pressure_handler: MemoryPressureHandler,
}

impl Week2MemoryAllocator {
    pub fn adapt_to_pattern(&mut self, pattern: &Pattern) -> AllocationAdjustment;
    pub fn forecast_memory_situation(&self) -> (MemoryPressureLevel, u32); // (level, time_to_critical_ms)
    pub fn get_allocation_efficiency(&self) -> f64;
}
```

**Week 2 통계**:
- 파일 개수: 4개
- 총 줄수: 890줄
- 테스트: 32개
- 핵심 개념: Predictive Allocation, Memory Compaction, Pressure Handling

---

### Week 3: I/O Bandwidth Management (Days 15-21)

**목표**: 예측된 I/O 패턴에 따라 읽기/쓰기 스케줄링 및 캐시 최적화

#### Day 15-16: I/O Pattern Analyzer (230줄)
```rust
pub struct IOMetrics {
    read_ops_per_sec: u64,
    write_ops_per_sec: u64,
    read_latency_us: u32,
    write_latency_us: u32,
    queue_depth: u32,
    disk_utilization: f64,        // 0.0-1.0
}

pub enum IOPattern {
    SequentialRead,
    SequentialWrite,
    RandomRead,
    RandomWrite,
    Mixed,
    Bursty,
}

pub struct IOPatternAnalyzer {
    current_metrics: IOMetrics,
    history: Vec<(u64, IOMetrics)>,
    detected_pattern: Option<IOPattern>,
}

impl IOPatternAnalyzer {
    pub fn detect_io_pattern(&self) -> IOPattern;
    pub fn estimate_io_saturation(&self) -> f64;             // 0.0-1.0
    pub fn predict_bandwidth_usage(&self) -> (u64, u32);     // (MB/sec, time_to_saturate_ms)
    pub fn detect_io_hot_spots(&self) -> Vec<(u64, f64)>;    // (address, access_count)
    pub fn get_queue_saturation(&self) -> f64;               // 0.0-1.0
}
```

**테스트** (8개):
- test_pattern_detection_sequential
- test_pattern_detection_random
- test_burst_detection
- test_bandwidth_prediction
- test_queue_depth_analysis
- test_latency_measurement
- test_hotspot_identification
- test_io_trend_forecasting

#### Day 17-18: I/O Scheduler (280줄)
```rust
pub struct IORequest {
    request_id: u32,
    io_type: IOType,               // Read/Write
    address: u64,
    size: u64,
    priority: u8,                  // 0-255
    deadline_us: u32,
}

pub enum IOSchedulePolicy {
    NOOP,                          // 들어온 순서대로
    FIFO,                          // FCFS
    Deadline,                      // 마감시간 우선
    CFQ,                           // Completely Fair Queuing
    Predictive,                    // Phase 6 기반 예측적
}

pub struct IOScheduler {
    queue: Vec<IORequest>,
    policy: IOSchedulePolicy,
    phase6_predictions: Vec<IOPattern>,
}

impl IOScheduler {
    pub fn enqueue_io_request(&mut self, request: IORequest);
    pub fn dequeue_next_request(&mut self) -> Option<IORequest>;
    pub fn reorder_for_pattern(&mut self, pattern: &IOPattern);
    pub fn apply_predictive_prefetch(&mut self, pattern: &IOPattern);
    pub fn get_queue_wait_time(&self, priority: u8) -> u32; // ms
    pub fn estimate_completion_time(&self, request_id: u32) -> u32; // ms
}
```

**테스트** (8개):
- test_io_queue_enqueue
- test_scheduling_policy_fifo
- test_deadline_scheduling
- test_predictive_reordering
- test_priority_preservation
- test_fairness_metric
- test_starvation_prevention
- test_bandwidth_efficient_ordering

#### Day 19-21: Cache Optimization Engine (250줄)
```rust
pub struct CacheBlock {
    address: u64,
    size: u64,
    access_count: u32,
    last_access: u64,
    dirty: bool,
}

pub struct CacheOptimizer {
    block_cache: HashMap<u64, CacheBlock>,
    max_cache_size: u64,
    eviction_policy: EvictionPolicy,
    prefetch_queue: Vec<u64>,
}

pub enum EvictionPolicy {
    LRU,
    LFU,
    Predictive,                   // Phase 6 기반 예측
}

impl CacheOptimizer {
    pub fn optimize_cache_for_pattern(&mut self, pattern: &IOPattern);
    pub fn prefetch_predicted_blocks(&mut self, addresses: &Vec<u64>);
    pub fn evict_cold_blocks(&mut self) -> u64;             // 해제된 바이트
    pub fn get_hit_rate(&self) -> f64;                      // 0.0-1.0
    pub fn suggest_cache_size(&self, pattern: &IOPattern) -> u64;
}
```

**테스트** (8개):
- test_cache_hit_tracking
- test_lru_eviction
- test_prefetch_effectiveness
- test_dirty_block_flushing
- test_cache_coherency
- test_pattern_specific_optimization
- test_working_set_identification
- test_cache_size_recommendation

#### Integration: IOBandwidthManager Week3 System (180줄)

**공개 API**:
```rust
pub struct Week3IOBandwidthManager {
    analyzer: IOPatternAnalyzer,
    scheduler: IOScheduler,
    cache: CacheOptimizer,
}

impl Week3IOBandwidthManager {
    pub fn adapt_to_pattern(&mut self, pattern: &Pattern) -> IOAdjustment;
    pub fn forecast_io_saturation(&self) -> (f64, u32);     // (saturation, time_to_critical_ms)
    pub fn get_io_efficiency(&self) -> f64;
}
```

**Week 3 통계**:
- 파일 개수: 4개
- 총 줄수: 940줄
- 테스트: 32개
- 핵심 개념: I/O Scheduling, Predictive Prefetch, Cache Optimization

---

## 🔗 Phase 6 ↔ Phase 7 Integration

### 데이터 흐름

```
Phase 6 Outputs
├─ Pattern: MemoryLeak (confidence: 0.95)
├─ Prediction: Saturation 90% in 2.3 seconds
└─ Action: "AggressiveGC"
        ↓
Phase 7 Adaptation
├─ Memory Allocator: 500MB 사전 할당
├─ GC Priority: 상향 조정
└─ I/O Scheduler: GC I/O 우선순위 증가
        ↓
Result: 임계값 도달 전 문제 해결 ✅
```

### Integration Points

1. **CPU Scheduler** ← Pattern + ProactiveAction
   - Memory GC 필요 시: GC 스레드 우선순위 부스트
   - CPU 병목 감지 시: 스레드 친화성 재배치

2. **Memory Allocator** ← MemoryPrediction
   - 누수 예측 시: 여유 메모리 사전 할당
   - 압박 예상 시: 사전 컴팩션 실행

3. **I/O Bandwidth Manager** ← IOPattern
   - 대용량 읽기 예상 시: 캐시 프리페치
   - 쓰기 폭발 예상 시: 쓰기 버퍼 사전 할당

---

## 📊 Phase 7 Summary

### File Structure
```
src/adaptive/
├── Week1/ (CPU Scheduling)
│   ├── cpu_analyzer.fl
│   ├── thread_affinity.fl
│   ├── dynamic_scheduler.fl
│   └── mod.fl (integration)
├── Week2/ (Memory Allocation)
│   ├── memory_predictor.fl
│   ├── predictive_allocator.fl
│   ├── pressure_handler.fl
│   └── mod.fl (integration)
├── Week3/ (I/O Bandwidth)
│   ├── io_analyzer.fl
│   ├── io_scheduler.fl
│   ├── cache_optimizer.fl
│   └── mod.fl (integration)
└── phase7_integration.fl (전체 통합)

docs/
├── PHASE_7_ADAPTIVE_DESIGN.md (이 파일)
├── PHASE_7_WEEK1_REPORT.md (Days 1-7)
├── PHASE_7_WEEK2_REPORT.md (Days 8-14)
├── PHASE_7_WEEK3_REPORT.md (Days 15-21)
└── PHASE_7_ADAPTIVE_COMPLETE_REPORT.md (최종)
```

### Statistics
| Metric | Target | Notes |
|--------|--------|-------|
| **Code Lines** | 2,500-3,000 | Week1: 800, Week2: 890, Week3: 940 |
| **Test Cases** | 70+ | 32 + 32 + 32 = 96 |
| **Modules** | 12 | 4 + 4 + 4 = 12 |
| **Files** | 12 | 4 + 4 + 4 = 12 |
| **Pass Rate** | 100% | All tests must pass ✅ |

---

## 🎯 Success Criteria

✅ **Phase 7 완성 조건**:
1. ✓ CPU 스케줄링: 우선순위 부스트로 응답성 50% 개선
2. ✓ 메모리 할당: 사전 할당으로 OOM 지연 300% 증가
3. ✓ I/O 관리: 스케줄링으로 대역폭 활용률 80%+ 달성
4. ✓ 통합: Phase 6 예측을 Phase 7 리소스 할당에 100% 반영
5. ✓ 테스트: 96개 모두 통과 (100%)
6. ✓ 문서: 완전한 설계 + 주간 보고서 + 최종 보고서

---

## 📅 Timeline

```
Week 1 (Days 1-7):   CPU Scheduling (800줄, 32 tests)
  Day 1-2: CPU Analyzer
  Day 3-4: Thread Affinity Manager
  Day 5-7: Dynamic Priority Scheduler
  Integration & Testing

Week 2 (Days 8-14):  Memory Allocation (890줄, 32 tests)
  Day 8-9: Memory Predictor
  Day 10-11: Predictive Memory Allocator
  Day 12-14: Memory Pressure Handler
  Integration & Testing

Week 3 (Days 15-21): I/O Bandwidth (940줄, 32 tests)
  Day 15-16: I/O Pattern Analyzer
  Day 17-18: I/O Scheduler
  Day 19-21: Cache Optimization Engine
  Integration & Testing

Final:               Complete Phase 7
  Documentation & GOGS Push
```

---

**철학**: "기록이 증명이다"
- 예측 (Phase 6) → 적응 (Phase 7) → 성능 개선 (Phase 8)

이 설계를 바탕으로 3주 동안 체계적으로 구현하겠습니다.

**다음**: Phase 7 Week 1 시작
