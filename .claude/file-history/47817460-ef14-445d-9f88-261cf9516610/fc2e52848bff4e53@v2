# Phase I: Chaos Engineering - Fault Injection & Validation
**FreeLang Bare-Metal OS Kernel - Production Resilience Verification**

**기간**: 3주 (Days 1-21)
**목표**: Phase H Design by Failure를 실제 chaos 환경에서 검증
**철학**: "관찰과 제어가 실제 장애 상황에서도 동작하는가?"

---

## 🎯 Phase I 목표

### Primary Goal
Phase H에서 구현한 5개 모듈이 **실제 시스템 장애** 속에서도:
1. ✅ 장애를 정확히 감지하는가?
2. ✅ 원인을 올바르게 분석하는가?
3. ✅ 적절한 정책을 선택하는가?
4. ✅ 성공적으로 복구하는가?

### 4가지 Success Criteria 검증
1. **Real Kernel Event Trace**: 실제 fault injection에서 events 수집
2. **Automatic Policy Execution**: 정책이 자동으로 적용되고 개선되는가?
3. **Full Reasoning Logging**: 모든 결정이 명확한 이유와 함께 로깅되는가?
4. **Automatic Root Cause Output**: SLO violation 시 root cause가 자동으로 생성되는가?

---

## 📋 3주 구현 계획

### Week 1: Fault Injection Framework (Days 1-7)

**목표**: Phase H 모듈들이 inject되는 장애를 감지할 수 있는 기반 구축

#### Day 1-2: Core Fault Types (300줄)
```
pub enum FaultType {
    LatencySpike(u64),           // 5ms ~ 500ms 레이턴시 추가
    MemoryLeak(u64),             // 1MB ~ 100MB 메모리 소비
    GCPause(u32),                // 50ms ~ 200ms GC pause
    IOSaturation(u32),           // 10 ~ 100개 pending I/O ops
    NetworkPartition,            // 완전 네트워크 단절
    CPUThrottling(f64),          // 50% ~ 100% CPU 제한
    ContextSwitchSpike(u32),     // normal × 2 ~ 5
}

pub struct FaultInjectionEngine {
    active_faults: Vec<(FaultType, u64, Duration)>, // (fault, start_us, duration)
    fault_timeline: Vec<(u64, FaultType, String)>,  // 감사 추적
    fault_effects: HashMap<String, f64>,            // metric → degradation
}
```

#### Day 3-4: Metric Degradation (250줄)
```
struct MetricDegradation {
    p99_latency_increase: f64,      // LatencySpike 적용 시 증가율
    memory_usage_increase: f64,     // MemoryLeak 적용 시
    gc_pause_increase: f64,         // GCPause 적용 시
    io_queue_depth_increase: f64,   // IOSaturation 적용 시
    cache_hit_drop: f64,            // NetworkPartition 시
}

// 각 fault type마다 realistic degradation 모델
```

#### Day 5-7: Fault Injection API (200줄)
```
pub struct FaultInjector {
    engine: FaultInjectionEngine,
    observability_stack: &mut ObservabilityStack,  // Phase H와 통합
}

impl FaultInjector {
    pub fn inject_fault(&mut self, fault_type: FaultType, duration_ms: u32) → FaultId
    pub fn check_fault_active(&self, metric_name: &str) → Option<f64> // degradation
    pub fn get_injected_metrics(&self) → HashMap<String, f64>
    pub fn generate_fault_report(&self) → String
}
```

**테스트**: 7개 (각 fault type당 1개 + unified test)

---

### Week 2: Realistic Chaos Scenarios (Days 8-14)

**목표**: Phase H 모듈들이 실제 운영 환경에서 발생하는 장애를 처리할 수 있는가?

#### Scenario 1: "Tail Latency Under Load" (Days 8-9)
```
Timeline:
  T-5s:   Load increases (normal)
  T-0s:   LatencySpike(+150ms) injected
  T+1s:   P99 jumps from 100ms to 250ms → FailureDetected
  T+2s:   Causality: LatencySpike → Scheduler queue increase
  T+3s:   Policy: GCMode(aggressive) selected
  T+10s:  P99 drops to 120ms (52% improvement) → Policy success
  T+70s:  Return to idle (improvement stable)

Validation:
  ✅ FailureDetector catches P99 > 100ms
  ✅ CausalityGraph links LatencySpike → P99
  ✅ PolicyController selects correct rule
  ✅ AdaptiveLoop measures improvement
```

#### Scenario 2: "Memory Degradation" (Days 9-10)
```
Timeline:
  T-10s:  MemoryLeak(1MB/sec) starts
  T-0s:   Memory pressure reaches 0.82 → MemoryCollapse detected
  T+1s:   Causality: MemoryLeak + GC pause spike
  T+2s:   Policy: PreallocSize(256MB) applied (priority 5)
  T+15s:  Allocation failures stop → Improvement measured
  T+80s:  System stable

Expected:
  - Allocation failures go to 0 after prealloc
  - GC pause time drops 30-40%
  - No OOM crash
```

#### Scenario 3: "I/O Bottleneck" (Days 10-11)
```
Timeline:
  T-5s:   IOSaturation(50 pending) injected
  T-2s:   Queue depth spikes, cache hit rate drops
  T-0s:   IOBackpressure detected (confidence 0.85)
  T+1s:   Causality: IOQueue → latency spike
  T+2s:   Policy: IOSchedulerPolicy(deadline) selected
  T+20s:  Queue depth normalized
  T+90s:  Rollback check (improvement > 30%?)

Expected:
  - I/O latency decreases by 30-35%
  - No thrashing (single policy application)
```

#### Scenario 4: "Cascading Failures" (Days 11-12)
```
Timeline:
  T-5s:   GCPause(150ms) injected
  T-2s:   LatencyExplosion detected
  T-0s:   MemoryLeak(2MB/sec) injected (compounding)
  T+1s:   TWO failures detected simultaneously
  T+3s:   Policy priority: MEMORY_PREALLOC (5) > LATENCY_GC (4)
  T+10s:  Memory stabilizes first
  T+25s:  Latency stabilizes second
  T+100s: Both metrics normal

Validation:
  ✅ Multiple failures handled correctly
  ✅ Priority-based policy selection works
  ✅ No race conditions or oscillations
  ✅ Learning engine updates for both failures
```

#### Scenario 5: "Network Partition & Recovery" (Days 12-14)
```
Timeline:
  T-5s:   NetworkPartition injected
  T-2s:   Cache hit rate drops 0.9 → 0.3
  T-0s:   IOBackpressure detected (but root cause is network, not disk)
  T+1s:   Causality inference challenges (no direct IO queue depth)
  T+3s:   Policy applied anyway (pessimistic but safe)
  T+30s:  NetworkPartition removed (recovery)
  T+35s:  Cache hit rate recovers to 0.85
  T+90s:  System returns to normal

Validation:
  ✅ Graceful degradation (policy helps even with wrong diagnosis)
  ✅ Recovery is automatic
  ✅ No policy thrashing during recovery
  ✅ Learning engine notes this scenario for next time
```

**테스트**: 5개 scenarios × 3 iterations = 15개

---

### Week 3: Validation & Integration (Days 15-21)

**목표**: All 4 success criteria를 chaos 환경에서 검증

#### Day 15-17: Success Criteria Validation
```
TEST 1: Real Kernel Event Trace Capture
  Input: All 5 chaos scenarios running simultaneously
  Expect: CausalityGraph captures ALL relevant events
  Verify:
    - No missed events (100% event capture)
    - Correct temporal ordering
    - Confidence scores consistent with chaos severity

TEST 2: Automatic Policy Execution
  Input: Fault injection + ObservabilityStack detect
  Expect: PolicyController applies correct policy automatically
  Verify:
    - No manual intervention needed
    - Policy matches failure type
    - Timing: detect to execute < 2s

TEST 3: Full Reasoning Logging
  Input: Each scenario's failure → root cause → policy
  Expect: AutoReport generated with complete reasoning
  Verify:
    - Evidence field populated
    - Confidence scores justified
    - Timeline reconstruction accurate
    - Markdown report valid

TEST 4: Automatic Root Cause Output
  Input: SLO violation detection
  Expect: Root cause identified within 1 second
  Verify:
    - Primary cause confidence > 0.7
    - Secondary causes ranked correctly
    - Expected improvement matches actual
```

#### Day 18-19: Integration Testing
```
INTEGRATION TEST 1: Week 1 → Week 2 Flow
  Observe (Week 1: Tracer + Metrics) →
  Detect (failure_definitions) →
  Analyze (causality_graph) →
  Decide (policy_controller) →
  Execute (via Phase 7 adapter) →
  Measure (collect improvement metrics) →
  Learn (update LearningEngine)

INTEGRATION TEST 2: All 5 Chaos Scenarios Simultaneously
  - Run all 5 scenarios in parallel
  - Verify no race conditions
  - Check adaptive loop state machine integrity
  - Measure overall system recovery time

INTEGRATION TEST 3: Rollback Correctness
  - 3 policies that should NOT improve (intentional)
  - Verify rollback happens at 90s or 20% threshold
  - Confirm system returns to pre-policy state
  - Check learning engine penalizes failed policies
```

#### Day 20-21: Final Report & Metrics
```
PHASE_I_FINAL_REPORT.md (2,000줄):
  1. Executive Summary (Success/Failure)
  2. Chaos Scenarios Results (5 scenarios × 3 iterations = 15 tests)
  3. Success Criteria Verification
  4. Performance Metrics
  5. Root Cause Analysis: Why some failures?
  6. Recommendations for Phase II
  7. Lessons Learned

Key Metrics:
  - Policy success rate: target > 75%
  - False positive rate: target = 0%
  - Recovery time: target < 90s
  - Learning effectiveness: target > 60% (policies repeat correctly)
```

---

## 🔧 파일 구조

```
freelang-os-kernel/
├── src/chaos_engineering/
│   ├── fault_injection.fl         (300줄, Days 1-2)
│   ├── metric_degradation.fl      (250줄, Days 3-4)
│   ├── fault_injector_api.fl      (200줄, Days 5-7)
│   ├── scenario_generator.fl      (400줄, Days 8-14)
│   ├── chaos_validator.fl         (350줄, Days 15-17)
│   ├── integration_harness.fl     (300줄, Days 18-19)
│   └── mod.fl                     (150줄)
│
├── tests/
│   ├── chaos_scenario_tests.fl    (800줄, 15 scenarios)
│   ├── success_criteria_tests.fl  (600줄, 4 criteria × 3 runs)
│   └── integration_tests.fl       (500줄, 3 integration scenarios)
│
└── docs/
    ├── PHASE_I_CHAOS_ENGINEERING_DESIGN.md (이 파일)
    ├── PHASE_I_WEEK1_REPORT.md
    ├── PHASE_I_WEEK2_REPORT.md
    └── PHASE_I_FINAL_REPORT.md
```

---

## 📊 성과 기대치

| 항목 | Week 1 | Week 2 | Week 3 | 합계 |
|------|--------|--------|--------|------|
| **파일** | 4 | 1 | 2 | **7** |
| **줄수** | 750 | 400 | 1,900 | **3,050** |
| **테스트** | 7 | 15 | 10 | **32** |

### Final Phase I Statistics
- **Total Chaos Scenarios**: 5 × 3 iterations = 15
- **Total Integration Tests**: 3
- **Success Criteria Validations**: 4 × 3 runs = 12
- **Total Test Cases**: 32
- **Expected Pass Rate**: > 95% (allow 1-2 failures for debugging)

---

## 🎯 Phase I 성공 기준

### Overall Success
✅ **PASS**: All 4 success criteria verified in real chaos, > 75% scenario success rate
⚠️ **CONDITIONAL**: > 70% scenario success rate, but specific criteria issues identified
❌ **FAIL**: < 70% success rate OR critical race conditions detected

### Critical Requirements
- **Zero OOM crashes**: All memory collapse scenarios must NOT crash system
- **Zero deadlocks**: No priority inversion or thread deadlock
- **Zero data corruption**: All improvements are safe (no lost data)
- **Automatic recovery**: Every scenario must recover without manual intervention

---

## 🚀 다음 단계 (Phase II)

### Phase II: Production Hardening (4주)
If Phase I succeeds:
1. **Observability at Scale**: 1000+ events/sec handling
2. **Multi-node Coordination**: Phase H across distributed cluster
3. **Persistent State**: Failure history + policy effectiveness log
4. **Real Kernel Integration**: Linux/RTOS kernel hooks

---

## 🔍 Phase I 철학

### "실패는 성공의 어머니다"
- 이 phase는 **실패하는 시나리오를 찾는 것**이 목표
- 각 실패는 Phase H 설계의 약점을 드러냄
- 약점을 보강하여 production-ready 시스템을 만듦

### "관찰만으로는 부족하다"
- Phase H의 관찰력이 진정한 능력인지 증명
- Chaos 속에서 정확한 원인 분석이 가능한가?
- 자동 복구가 정말 작동하는가?

### "기록이 증명이다"
- 각 scenario의 결과를 상세히 로깅
- 모든 결정과 그 이유를 기록
- 나중에 재현 가능하도록 (deterministic chaos)

---

## 📝 결론

Phase I은 Phase H의 "이론적 검증"을 "실제 검증"으로 전환하는 중요한 단계입니다.

**이 phase를 통해 증명하는 것**:
1. Phase H의 Design by Failure가 실제로 작동하는가?
2. Automatic policy selection이 올바른가?
3. Learning engine이 실제로 개선하는가?
4. Production 환경에서 안전한가?

**성공 시 얻는 것**:
- FreeLang OS가 **자율적으로 자신의 장애를 진단하고 복구**할 수 있는 시스템
- SRE 팀의 부담을 극적으로 감소시킬 수 있는 infrastructure
- Cloud-native 환경에서 요구하는 **resilience와 observability를 동시에 만족**

---

**생성**: 2026-03-03
**다음**: Day 1-2 (Fault Injection Framework) 구현 시작
