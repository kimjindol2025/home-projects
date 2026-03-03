# FreeLang Phase H Week 2: Design by Failure
## Observability → Control Authority → Adaptive System

**기간**: 2026-03-04 ~ 2026-03-10 (Days 8-14)
**모드**: Design by Failure (역설계: 실패에서 정책으로)
**목표**: 단순 관찰에서 능동적 제어 시스템으로 진화

---

## 🔴 1. Failure Class 정의 (3개만)

### **Failure Class A: Latency Explosion**

**신호**:
- P99 latency > SLO target (연속 3회)
- Scheduler queue depth > 50%
- Context switch rate spike

**원인 후보**:
- GC pause
- Lock contention
- Cache miss

**대응 정책**:
```
IF P99 > threshold 3회 연속:
  → GC aggressive mode ON
  → Thread affinity 재조정
  → Priority boost for critical paths
```

**측정 기준**:
- P99 recovery time < 10초
- Queue depth 정상화

---

### **Failure Class B: Memory Collapse**

**신호**:
- Allocator fragmentation > 0.7
- GC pause > 100ms
- Memory pressure > 80%

**원인 후보**:
- Long-lived objects
- Fragmentation pattern
- Allocation size mismatch

**대응 정책**:
```
IF memory_pressure > 0.80:
  → Prealloc pool size +256MB
  → GC mode: incremental (from stop-the-world)
  → Compaction trigger threshold lower
```

**측정 기준**:
- Fragmentation ratio < 0.5
- GC pause < 50ms

---

### **Failure Class C: I/O Backpressure**

**신호**:
- Queue depth > 64
- I/O syscall latency > 1000µs
- Cache hit rate drop > 20%

**원인 후보**:
- Disk saturation
- Sequential vs random access mismatch
- Prefetch miss

**대응 정책**:
```
IF queue_depth > 64:
  → I/O scheduler: switch to deadline mode
  → Cache: predictive prefetch intensity up
  → Task scheduling: deprioritize I/O intensive tasks
```

**측정 기준**:
- Queue depth < 32
- I/O latency < 500µs

---

## 🔗 2. Causality Graph 생성 (핵심 구조)

### **문제점: 현재 Week 1**
```
Timeline (단순 나열):
  T=0ms:   start_trace
  T=10ms:  GC begin
  T=120ms: GC end
  T=121ms: scheduler backlog detected
  T=340ms: P99 spike
  T=500ms: SLO violation
```

### **목표: Week 2 Causality**
```
Causality Graph (원인 연결):

GC pause (120ms)
  │
  ├─ causes: Scheduler queue accumulation
  │   │
  │   └─ causes: Context switch spike
  │       │
  │       └─ causes: P99 latency explosion (340ms)
  │           │
  │           └─ causes: SLO violation
  │
  └─ impact_score: 0.72 (confidence)

I/O contention (secondary)
  │
  └─ impact_score: 0.21 (confidence)
```

### **생성 알고리즘**

```
1. Collect Events (last 60 seconds)
   → [GC, Scheduler, I/O, Latency, SLO]

2. Event Clustering
   → Group by temporal proximity
   → Weight by impact (latency contribution)

3. Causality Inference
   → T(cause) < T(effect)
   → Correlation analysis
   → Impact quantification

4. Confidence Scoring
   → How strong is the causal link?
   → Multiple competing hypotheses?
   → Evidence strength calculation

5. Output: Root Cause Ranked List
   → [Cause 1: 0.72, Cause 2: 0.21, ...]
```

---

## 🎮 3. Metric → Policy 매핑 (Control Authority)

### **Rule Engine 구조**

```rust
// failure_driven_policy.fl

pub struct PolicyRule {
    name: String,
    trigger: Condition,        // When to apply
    action: PolicyChange,      // What to do
    impact: f64,              // Expected improvement (0-1)
    side_effects: Vec<String>, // What else changes
}

pub enum Condition {
    P99Exceeded { count: u32, threshold_ms: u64 },
    MemoryPressure { ratio: f64 },
    IOQueueDepth { depth: u32 },
    Combined { op: LogicOp, conditions: Vec<Condition> },
}

pub enum PolicyChange {
    GCMode(String),            // "aggressive" / "incremental"
    PreallocSize(u64),         // +MB
    IOSchedulerPolicy(String), // "deadline" / "cfq"
    ThreadAffinity(String),    // "rebalance" / "lock"
}

pub struct ControlDecision {
    trigger_cause: String,
    applied_rules: Vec<PolicyRule>,
    expected_impact: f64,
    applied_at: u64,
    logged_reason: String,     // "왜 바꿨는지"
}
```

### **Soft Control (Option B)**

**허용되는 변경**:
- ✅ GC aggressive mode ON/OFF
- ✅ Prealloc pool size ±X
- ✅ Scheduler weight adjustment
- ✅ I/O scheduling policy switch

**금지되는 변경**:
- ❌ 코어 알고리즘 교체
- ❌ Hardware config 변경
- ❌ Permanent data structure 수정

**모든 변경은**:
- 🔄 Reversible (원상복구 가능)
- 📝 Logged (이유 기록)
- ⏱️ Timed (60초 후 재평가)

---

## 🧠 4. Root Cause Engine v0.1

### **구조**

```
SLO Violation 감지
  ↓
[Root Cause Engine]
  ├─ Phase 1: Event Collection
  │   └─ Last 60s trace + metrics
  │
  ├─ Phase 2: Causality Analysis
  │   ├─ Temporal ordering
  │   ├─ Correlation strength
  │   └─ Impact quantification
  │
  ├─ Phase 3: Hypothesis Ranking
  │   ├─ Most probable cause
  │   ├─ Secondary causes
  │   └─ Confidence scores
  │
  └─ Phase 4: Recommendation
      ├─ "Why this happened"
      ├─ "Policy we applied"
      └─ "Expected improvement"
  ↓
Auto-generated PostMortem
```

### **출력 포맷**

```
╔════════════════════════════════════════════╗
║         SLO Violation Report               ║
║         2026-03-04 15:23:45 UTC            ║
╚════════════════════════════════════════════╝

VIOLATION: P99 latency 340ms > target 100ms
SEVERITY: HIGH (SLO 5th consecutive breach)

ROOT CAUSE ANALYSIS
═══════════════════════════════════════════════

Primary Cause (confidence: 0.72)
  ├─ Event: GC pause 120ms
  ├─ Timing: started T-200ms relative to violation
  ├─ Impact: 340ms P99 latency
  └─ Mechanism: Scheduler queue blocked during pause

Secondary Cause (confidence: 0.21)
  ├─ Event: I/O queue depth spike
  └─ Impact: 20% P99 contribution

POLICY DECISIONS APPLIED
═══════════════════════════════════════════════

Decision 1: GC Aggressive Mode
  ├─ Trigger: P99 > threshold 3 consecutive times
  ├─ Applied at: T=0ms
  ├─ Expected impact: 30-40% P99 reduction
  └─ Logs: "GC_MODE_CHANGE" rule fired

Decision 2: Prealloc +256MB
  ├─ Trigger: Memory pressure > 80%
  ├─ Applied at: T=50ms
  ├─ Expected impact: Reduce GC frequency 50%
  └─ Logs: "MEMORY_PREALLOC" rule fired

OUTCOME PREDICTION
═══════════════════════════════════════════════

If above policies successful:
  └─ Expected P99 recovery: < 10 seconds
  └─ Probability: 0.68

NEXT STEPS
═══════════════════════════════════════════════

1. Monitor next 300 seconds for confirmation
2. If P99 still > 80ms: escalate to manual review
3. If recovered: commit policy changes for 24 hours
4. Log this incident for pattern learning

═══════════════════════════════════════════════
Generated by: Root Cause Engine v0.1
Timestamp: 2026-03-04 15:23:45.123 UTC
```

---

## 📊 5. Week 2 모듈 구성 (새로운 설계)

### **기존 계획 → Design by Failure 변경**

```
Before (Old Plan):
  - sre_operations.fl (350 lines)
  - chaos_real_injection.fl (350 lines)
  - postmortem_analyzer.fl (300 lines)

After (Design by Failure):
  - failure_definitions.fl (250 lines)
  - causality_graph.fl (350 lines)
  - policy_controller.fl (350 lines)
  - root_cause_engine.fl (400 lines)
  - adaptive_loop.fl (250 lines)
```

### **각 모듈의 역할**

#### **failure_definitions.fl** (Days 8-9 전반부)
```
목적: Failure의 명확한 정의
  ├─ FailureClass enum (Latency, Memory, IO)
  ├─ Signal thresholds (정확한 수치)
  ├─ Detection logic (어떻게 감지할 것인가)
  └─ Expected recovery time

결과: 모호함 제거, 자동 감지 가능
```

#### **causality_graph.fl** (Days 9-10)
```
목적: Trace → Causality 변환
  ├─ Event collection (60s 윈도우)
  ├─ Temporal ordering
  ├─ Correlation analysis
  ├─ Impact scoring
  └─ Graph output

결과: "왜 실패했는지"를 자동으로 파악
```

#### **policy_controller.fl** (Days 10-11)
```
목적: Metric 기반 자동 정책 변경
  ├─ PolicyRule 정의 (3개 failure class)
  ├─ Rule matching (조건 판단)
  ├─ Action execution (정책 변경)
  ├─ Side effect tracking
  └─ Reversibility guarantee

결과: 자동 대응 (Control Authority)
```

#### **root_cause_engine.fl** (Days 11-13)
```
목적: 자동 진단 및 보고서 생성
  ├─ Causality graph 입력
  ├─ Confidence scoring
  ├─ Hypothesis ranking
  ├─ Recommendation generation
  └─ PostMortem auto-report

결과: "이게 맞는 원인이고, 이렇게 대응했으며, 이 정도 효과"
```

#### **adaptive_loop.fl** (Days 13-14)
```
목적: 5개 모듈 통합 (Feedback Loop)
  ├─ Failure detection
  ├─ Causality analysis
  ├─ Policy decision
  ├─ Action execution
  ├─ Outcome measurement
  └─ Learning feedback

결과: 폐쇄 루프 (Adaptive System)
```

---

## ✅ 6. 성공 기준 (엄격)

### **통과 필수 조건 (4가지, 모두 자동화)**

#### **1️⃣ 실제 커널 이벤트 Trace 캡처**

```
Evidence Required:
  □ GC pause detection in trace
  □ Scheduler queue depth in metrics
  □ syscall latency measurement
  □ Context switch counting

Acceptance: Week 1 Phase G와의 연결 코드
```

#### **2️⃣ Metric 기반 정책 변경 실행**

```
Evidence Required:
  □ P99 > threshold → GC mode change
  □ Memory pressure > 80% → prealloc trigger
  □ I/O queue depth > 64 → scheduler policy switch
  □ All changes logged with reason

Acceptance: ControlDecision 구조에 적어도 3개 이상의 실행된 규칙
```

#### **3️⃣ 정책 변경 로그 기록**

```
Evidence Required:
  □ "Why this policy was triggered"
  □ "What condition was met"
  □ "What action was taken"
  □ "When was it applied (timestamp)"
  □ "Expected impact"

Acceptance: PostMortem에 "Policy decision log" 섹션 존재
```

#### **4️⃣ SLO 위반 시 Root Cause 자동 출력**

```
Evidence Required:
  □ Causality graph generated (confidence > 0.5)
  □ Primary cause identified
  □ Impact quantified
  □ Policy decisions explained
  □ Outcome predicted

Acceptance: 실제 SLO violation 시뮬레이션 후
            자동으로 정확한 root cause 출력
```

### **실패 기준**

```
❌ 수동 개입 필요
❌ 로그 없는 정책 변경
❌ Trace와 Metric의 단절
❌ Policy 역할이 불명확
❌ PostMortem이 템플릿 수준
```

---

## 🎯 7. 아키텍처 흐름도

```
┌─────────────────────────────────────────────────────────┐
│         Phase H Week 2: Adaptive Loop                   │
└─────────────────────────────────────────────────────────┘

Metrics Collection (Week 1)
  ↓
Failure Detection (failure_definitions.fl)
  "P99 > threshold 3회" → FailureClass::LatencyExplosion
  ↓
Causality Analysis (causality_graph.fl)
  Events: [GC_pause_120ms, SchedulerQueue_spike, ...]
  Graph: GC → Scheduler → Latency → SLO_violation
  Confidence: 0.72
  ↓
Policy Decision (policy_controller.fl)
  Rule Match: "P99 exceeded → apply GC aggressive"
  Side Effects: [scheduler_weight_change, prealloc_trigger]
  Decision: ControlDecision { rules: [...], logged_reason: "..." }
  ↓
Action Execution
  Phase 7 receives: "GC_MODE = aggressive"
  Phase 7 receives: "PREALLOC_SIZE += 256MB"
  (Both logged with timestamp)
  ↓
Outcome Measurement
  P99 trend: 340ms → 280ms → 150ms → 85ms (recovery!)
  Memory pressure: 85% → 72% → 58%
  ↓
Root Cause Output (root_cause_engine.fl)
  AutoReport:
    └─ Primary Cause: GC pause (0.72)
    └─ Applied Policies: [GC_aggressive, PREALLOC_256MB]
    └─ Expected Impact: 30-40% P99 reduction
    └─ Measured Impact: 75% (EXCEEDED expectations!)
  ↓
Learning Feedback
  "GC_aggressive rule worked well → increase confidence"
  "Prealloc threshold from 80% → 75% for faster trigger"
  ↓
[Loop continues...]
```

---

## ⚠️ 8. 위험 요소 및 대응

### **위험 1: 정책 충돌**

```
Problem:
  GC aggressive ON → more CPU usage
  I/O scheduler: deadline → latency increase
  Conflict!

Mitigation:
  □ Mutual exclusion rules
  □ Priority ordering (GC > I/O > others)
  □ Rollback on failure detection
```

### **위험 2: Thrashing (진동)**

```
Problem:
  GC on → P99 drops → off → P99 rises → on (반복)

Mitigation:
  □ Hysteresis (P99 < 80ms to turn off, > 100ms to turn on)
  □ Min duration: 60 seconds per policy
  □ Frequency limit: max 1 change per 10 seconds
```

### **위험 3: Metric Corruption**

```
Problem:
  Policy A changes metrics
  Causality engine mistakes it for real improvement
  Applies wrong Policy B

Mitigation:
  □ Tag metrics with "policy_induced" flag
  □ Filter out policy-caused signals
  □ Validation: is the improvement from policy or external?
```

---

## 🎬 9. 예상 실행 흐름 (실제 시나리오)

### **Scenario: Memory Leak Detection**

```
T=0s:   Normal operation
        GC pause: 10ms
        Memory: 512MB
        P99: 45ms

T=30s:  GC pause: 25ms (detected)
        Memory: 562MB (trend up)
        P99: 58ms

T=60s:  GC pause: 55ms (spike)
        Memory: 620MB
        P99: 82ms (approaching SLO 100ms)

T=90s:  GC pause: 120ms ← THRESHOLD
        Memory: 700MB
        P99: 340ms ← SLO VIOLATION

        [System Action]
        1. Failure detected: MemoryCollapse
        2. Causality: GC spike → Scheduler queue → P99
        3. Policy applied: PREALLOC +256MB, GC_MODE aggressive
        4. Phase 7 receives commands: [GC_aggressive, prealloc_256]

T=120s: GC pause: 45ms ← policy working
        Memory: 456MB (drops due to better GC)
        P99: 95ms ← below SLO

        [PostMortem Generated]
        "Primary Cause: Memory fragmentation (0.82 confidence)
         Applied: GC_aggressive + Prealloc +256MB
         Result: P99 recovered within 30 seconds
         Recommendation: Keep increased GC frequency for 24h"
```

---

## 📋 10. Week 2 구현 일정

```
Day 8-9 (Mon-Tue):
  ├─ failure_definitions.fl (250 lines, 6 tests)
  ├─ Integration: Phase H Week 1과 연결
  └─ Checkpoint: Failure detection working

Day 9-10 (Tue-Wed):
  ├─ causality_graph.fl (350 lines, 8 tests)
  ├─ Event clustering algorithm
  └─ Checkpoint: Causality graph generation validated

Day 10-11 (Wed-Thu):
  ├─ policy_controller.fl (350 lines, 8 tests)
  ├─ Rule engine implementation
  └─ Checkpoint: Policy changes apply correctly

Day 11-13 (Thu-Fri-Mon):
  ├─ root_cause_engine.fl (400 lines, 10 tests)
  ├─ PostMortem auto-generation
  └─ Checkpoint: Report quality acceptable

Day 13-14 (Mon-Tue):
  ├─ adaptive_loop.fl (250 lines, 8 tests)
  ├─ Integration: all 5 modules connected
  ├─ End-to-end testing
  └─ Final validation (4 criteria)

Total Week 2:
  └─ ~1,600 lines + ~40 tests + complete adaptive system
```

---

## 🏆 최종 평가

### **Week 1 (현재 상태)**:
```
Observability Infrastructure ✅
  └─ "눈은 있지만 손이 없다"
```

### **Week 2 목표**:
```
Adaptive Control System ✅
  └─ "눈과 손과 뇌가 통합된 시스템"
```

### **Week 2 성공 시 달성**:
```
┌─ Autonomous Failure Detection
├─ Automatic Root Cause Analysis
├─ Self-Healing Policies
└─ Production-Hardened Operation

= Truly Adaptive OS Kernel
```

---

**Design by Failure 시작 준비 완료.**

다음: 코드 레벨 구현 구조 제시

기록이 증명이다. Your record is your proof. 🎉
