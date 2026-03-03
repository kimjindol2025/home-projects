# FreeLang Phase H: Observability & SRE Operations 설계 문서

**기간**: 2026-03-03 ~ (2주 집중)
**목표**: 약한 축(실제 프로덕션 운영)을 강화
**철학**: "기록이 증명이다" → "운영이 증명이다"

---

## 🎯 Phase H 핵심 목표

Phase 7의 **적응형 리소스 할당**을 실제 프로덕션 환경에서 **운영 가능하게** 만드는 것.

### 목표 달성 기준

✅ **분산 트레이싱**: 모든 요청의 전체 경로 가시화
✅ **실시간 모니터링**: CPU/Memory/I/O 메트릭 자동 수집
✅ **SRE 운영**: Error budget, SLO/SLI 자동 측정
✅ **Chaos Engineering**: 실제 장애 주입 후 자동 복구 검증
✅ **성능 측정**: perf/valgrind/eBPF로 실제 성능 측정
✅ **Post-Mortem**: 장애 분석 자동 생성

---

## 🏗️ Phase H 구조 (2주, ~2,000줄)

```
Phase H: Observability & SRE Operations (2,000줄 목표)
═══════════════════════════════════════════════════════

Week 1: 분산 트레이싱 & 모니터링 (1,000줄)
├─ distributed_tracer.fl     (350줄) - Trace ID, span, context
├─ metrics_collector.fl      (350줄) - CPU/Memory/I/O 메트릭
└─ monitoring_dashboard.fl   (300줄) - 실시간 가시화

Week 2: SRE 운영 & Chaos 검증 (1,000줄)
├─ sre_operations.fl         (350줄) - Error budget, SLO/SLI
├─ chaos_real_injection.fl   (350줄) - 실제 장애 주입
└─ postmortem_analyzer.fl    (300줄) - 자동 분석 & 리포팅
```

---

## 📊 Week 1: 분산 트레이싱 & 모니터링

### 1️⃣ distributed_tracer.fl (350줄, 8 tests)

**핵심 개념**:
```
Request Flow (Phase 7 적응 결정 포함):

[Client]
  ↓ (Trace ID: 550e8400-e29b-41d4-a716-446655440000)
[HTTP Listener]
  ├─ span: http_request (10ms)
  ├─ ctx: trace_id, request_id, user_id
  ↓
[Phase 7: CPU Analyzer]
  ├─ span: cpu_analysis (2ms)
  ├─ span: cpu_detect_pattern (1ms) [MemoryLeak 감지]
  ↓
[Phase 7: Memory Predictor]
  ├─ span: memory_predict (1.5ms)
  ├─ span: prealloc_check (0.5ms)
  ↓
[Phase 7: I/O Scheduler]
  ├─ span: io_dispatch (0.8ms)
  ├─ span: cache_lookup (0.2ms) [Hit]
  ↓
[Response] (total: 17ms)
```

**구현 내용**:
```rust
pub struct TraceContext {
    trace_id: String,        // 550e8400-e29b-41d4-a716-446655440000
    span_id: String,
    parent_span_id: Option<String>,
    start_time_us: u64,
    end_time_us: u64,
    operation_name: String,
    tags: HashMap<String, String>,  // {phase: "Phase7", pattern: "MemoryLeak"}
    logs: Vec<(u64, String)>,       // 타임스탬프, 이벤트
}

pub struct DistributedTracer {
    traces: HashMap<String, Vec<TraceContext>>,  // trace_id → spans
    root_traces: Vec<String>,
    sampling_rate: f64,  // 1.0 = 100% (프로덕션에서는 0.1 = 10%)
}

impl DistributedTracer {
    pub fn start_trace(&mut self, trace_id: String) → TraceContext
    pub fn start_span(&mut self, trace_id: String, parent_span_id: String, operation: String) → TraceContext
    pub fn end_span(&mut self, trace_id: String, span_id: String) → Option<u64>  // latency_us
    pub fn add_tag(&mut self, trace_id: String, span_id: String, key: String, value: String)
    pub fn add_log(&mut self, trace_id: String, span_id: String, message: String)
    pub fn get_trace_latency(&self, trace_id: String) → Option<u64>  // 전체 지연
    pub fn export_jaeger_format(&self) → String  // JSON export
}
```

**테스트** (8개):
- test_trace_context_creation
- test_span_hierarchy
- test_latency_calculation
- test_tag_addition
- test_log_tracking
- test_phase7_integration (Phase 7 span 추적)
- test_jaeger_export
- test_sampling_rate

---

### 2️⃣ metrics_collector.fl (350줄, 8 tests)

**핵심 메트릭**:
```
┌─────────────────────────────────────────┐
│  Real-time Metrics Collection           │
├─────────────────────────────────────────┤
│ CPU:      utilization, context_switches │
│ Memory:   usage, gc_pause, fragmentation│
│ I/O:      latency, throughput, queue    │
│ Network:  latency, packet_loss          │
│ Phase 7:  adaptation_count, decisions   │
└─────────────────────────────────────────┘
```

**구현 내용**:
```rust
pub struct Metric {
    name: String,              // "cpu.utilization" / "memory.gc_pause_ms"
    timestamp_us: u64,
    value: f64,
    tags: HashMap<String, String>,  // {phase: "Phase7", pattern: "MemoryLeak"}
    exemplar: Option<String>,  // trace_id 연결
}

pub struct MetricsCollector {
    metrics: Vec<Metric>,
    max_history: usize,  // 1시간 = 3600개 (1초 interval)
    flush_interval_ms: u32,
    phase7_adapter: AdaptiveResourceManager,  // Phase 7 통합
}

impl MetricsCollector {
    pub fn record_cpu_metric(&mut self, util: f64, ctx_switches: u32, trace_id: Option<String>)
    pub fn record_memory_metric(&mut self, usage_mb: u64, gc_pause_ms: u32, trace_id: Option<String>)
    pub fn record_io_metric(&mut self, latency_us: u64, throughput_mb: f64, queue_depth: u32)
    pub fn record_adaptation_metric(&mut self, phase: &str, decision: &str, trace_id: String)

    pub fn get_percentile(&self, metric_name: &str, percentile: u32) → Option<f64>  // P99
    pub fn get_rate_of_change(&self, metric_name: &str, window_sec: u32) → Option<f64>  // 변화율
    pub fn detect_anomaly(&self, metric_name: &str) → bool  // σ > 3 detection
    pub fn export_prometheus(&self) → String  // Prometheus format
}
```

**메트릭 목록**:
```
CPU:
  - cpu.utilization                  (0-100%)
  - cpu.context_switches_per_sec
  - cpu.thermal_level                (°C)

Memory:
  - memory.usage_mb
  - memory.gc_pause_ms
  - memory.fragmentation_ratio       (0-1)
  - memory.pressure_level            (Normal/Warning/Critical/Emergency)

I/O:
  - io.read_latency_us
  - io.write_latency_us
  - io.queue_depth
  - io.cache_hit_rate                (0-1)
  - io.bandwidth_utilization         (0-1)

Phase 7 Adaptation:
  - phase7.cpu_decisions_per_sec
  - phase7.memory_prealloc_triggered
  - phase7.io_pattern_switches
  - phase7.adaptation_latency_us
```

**테스트** (8개):
- test_metric_recording
- test_percentile_calculation
- test_rate_of_change
- test_anomaly_detection
- test_phase7_metric_integration
- test_prometheus_export
- test_correlation (CPU ↔ Memory ↔ I/O)
- test_metric_retention

---

### 3️⃣ monitoring_dashboard.fl (300줄, 8 tests)

**대시보드 구조**:
```
╔═══════════════════════════════════════════════════════════╗
║              Real-time Operations Dashboard              ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  [CPU Util]  [Memory]  [I/O Queue]  [GC Pause]  [SLO]   ║
║   ▄▄▄▄▄▄▄▄    ████      ▃▃▃▃▃     20ms↓        99.5%   ║
║   75%→↑80%   680/1024   depth:12   (↓ -5%)      ✓       ║
║                                                           ║
║  Phase 7 Decisions (last 5 minutes):                     ║
║  ├─ CPU priority boost: 3 times [MemoryLeak]            ║
║  ├─ Memory prealloc: 512MB [x2]                         ║
║  ├─ I/O cache switch: LRU→Predictive [IOBottleneck]    ║
║  └─ Average adaptation latency: 1.2ms                   ║
║                                                           ║
║  Active Traces: 247 | Sampled: 25 (10%) | Errors: 0    ║
║                                                           ║
║  Anomalies Detected: 1                                   ║
║  └─ Memory GC pause spike at 10:15:32 UTC              ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
```

**구현 내용**:
```rust
pub enum DashboardView {
    RealTime,       // 1초 업데이트
    Historical,     // 1시간 그래프
    HeatMap,        // CPU×Memory×I/O 3D
    TraceDetail,    // 특정 trace 분석
}

pub struct MonitoringDashboard {
    metrics: Vec<Metric>,
    traces: HashMap<String, Vec<TraceContext>>,
    alerts: Vec<Alert>,
    slo_status: (f64, f64, f64),  // (current, target, error_budget)
}

impl MonitoringDashboard {
    pub fn render_realtime(&self) → String
    pub fn render_historical(&self, minutes: u32) → String  // 그래프
    pub fn get_alerts(&self) → Vec<Alert>
    pub fn get_top_slow_traces(&self, count: usize) → Vec<(String, u64)>  // trace_id, latency
    pub fn get_adaptation_timeline(&self) → Vec<(u64, String, String)>  // (timestamp, phase, decision)
    pub fn export_html(&self) → String  // 웹 대시보드
}
```

**테스트** (8개):
- test_dashboard_rendering
- test_metric_aggregation
- test_alert_generation
- test_trace_correlation
- test_adaptation_timeline
- test_anomaly_highlighting
- test_slo_visualization
- test_html_export

---

## 📊 Week 2: SRE 운영 & Chaos 검증

### 4️⃣ sre_operations.fl (350줄, 8 tests)

**SRE 핵심 개념**:
```
SLO (Service Level Objective):
  - Availability: 99.9% (43분/월 downtime)
  - Latency P99: <100ms
  - Error rate: <0.1%
  - Cache hit rate: >80%

SLI (Service Level Indicator) - 측정:
  - (success_requests / total_requests) × 100%
  - (requests_under_100ms / total_requests) × 100%
  - ...

Error Budget:
  - 월간 허용 downtime: 43.2분
  - 현재 소비: 12분 (27.8%)
  - 남은 budget: 31.2분 (72.2%) ← 남은 "운영 여유"

Decision:
  - budget > 50%: 적극적 배포 가능
  - budget 20-50%: 신중한 배포만
  - budget < 20%: 긴급 대응 전용
```

**구현 내용**:
```rust
pub struct SLO {
    name: String,
    objective: f64,         // 0.999 = 99.9%
    window_days: u32,       // 30일
    metric_type: String,    // "availability" / "latency" / "error_rate"
}

pub struct SLI {
    name: String,
    current_value: f64,
    historical: Vec<(u64, f64)>,  // timestamp, value
}

pub struct ErrorBudget {
    total_allowed_failures: u32,  // 43분 = 2,580초
    consumed: u32,
    remaining: u32,
    can_deploy: bool,  // remaining > total * 0.2
}

pub struct SREOperations {
    slos: Vec<SLO>,
    slis: HashMap<String, SLI>,
    error_budgets: HashMap<String, ErrorBudget>,
    incident_log: Vec<Incident>,  // 모든 장애 기록
}

pub struct Incident {
    start_time: u64,
    end_time: u64,
    duration_sec: u32,
    root_cause: String,
    impact: String,  // "API latency +200ms"
    mitigation: String,
    prevention: String,  // 앞으로 어떻게?
}

impl SREOperations {
    pub fn measure_sli(&mut self, slo_name: &str, success: bool, latency_ms: u32)
    pub fn get_slo_status(&self, slo_name: &str) → (f64, f64)  // (current, target)
    pub fn get_error_budget(&self, slo_name: &str) → ErrorBudget
    pub fn can_deploy(&self) → bool  // 모든 SLO error_budget이 충분한가?
    pub fn record_incident(&mut self, incident: Incident)
    pub fn generate_monthly_report(&self) → String
    pub fn auto_rollback_decision(&self) → bool  // error_budget 기반 결정
}
```

**테스트** (8개):
- test_slo_measurement
- test_sli_calculation
- test_error_budget_consumption
- test_deployment_decision
- test_incident_recording
- test_monthly_report
- test_phase7_slo_integration
- test_auto_rollback

---

### 5️⃣ chaos_real_injection.fl (350줄, 8 tests)

**실제 장애 주입 시나리오**:

```
Scenario 1: CPU 스파이크 (CFS 우회)
  → CPU Analyzer 반응 확인
  → Thread affinity 재설정 확인
  → Memory 압박 연쇄 반응 확인
  → I/O 큐 증가 감지 확인
  ✓ 자동 복구 (priority boost)

Scenario 2: 메모리 누수 (계단식)
  → Memory Predictor 누수 감지
  → GC 긴급 실행
  → Prealloc 풀 소비
  → Memory pressure handler 단계 상향
  ✓ 자동 복구 (compaction + pageout)

Scenario 3: I/O 마감시간 위반
  → I/O Scheduler deadline 추적
  → Cache 적중률 급락 감지
  → LFU → Predictive 정책 전환
  ✓ 자동 복구 (prefetch 활성화)

Scenario 4: 네트워크 분할 (대기 중)
  → 분산 합의 실패
  → Leader election 재시작
  → Phase 6 예측 일시 중단
  → Phase 7 컨저버 모드 (보수적 정책)
  ✓ 자동 복구 (partition heal)

Scenario 5: Cascading 장애
  → CPU spike → Memory pressure → I/O backlog
  → 동시 3개 Phase 압박
  → 우선순위 충돌 해결
  ✓ 자동 복구 (차등 감속)
```

**구현 내용**:
```rust
pub enum ChaosInjectionType {
    CPUSpike(duration_sec: u32, target_utilization: f64),
    MemoryLeak(growth_rate_mb_per_sec: f64, duration_sec: u32),
    IOLatencyIncrease(multiplier: f64, duration_sec: u32),
    NetworkPartition(duration_sec: u32),
    ProcessKill(process_type: String),
}

pub struct ChaosInjector {
    active_injections: Vec<(u64, ChaosInjectionType)>,
    recovery_timeline: Vec<(u64, String, bool)>,  // (time, action, success)
    metrics_before: MetricsSnapshot,
    metrics_after: MetricsSnapshot,
}

pub struct RecoveryMetrics {
    time_to_detect_sec: u32,  // 장애 감지까지
    time_to_recover_sec: u32, // 정상 복구까지
    data_consistency: bool,   // 데이터 손상 없음?
    slo_violations: u32,      // 복구 중 SLO 위반 횟수
    automatic_recovery: bool, // 자동 복구 성공?
}

impl ChaosInjector {
    pub fn inject_chaos(&mut self, injection: ChaosInjectionType) → String  // execution_id
    pub fn wait_for_recovery(&mut self, max_wait_sec: u32) → RecoveryMetrics
    pub fn verify_consistency(&self) → bool
    pub fn export_chaos_report(&self) → String
}
```

**테스트** (8개):
- test_cpu_spike_injection
- test_memory_leak_detection
- test_io_backlog_recovery
- test_network_partition_heal
- test_cascading_failure_isolation
- test_recovery_metrics
- test_data_consistency
- test_phase7_chaos_integration

---

### 6️⃣ postmortem_analyzer.fl (300줄, 8 tests)

**자동 Post-Mortem 분석**:

```
Incident Report (자동 생성):

Title: Memory GC Pause Spike on 2026-03-03 10:15:32 UTC
Duration: 45 seconds
Severity: HIGH

Timeline (자동 추출):
  10:15:30  MemoryLeak pattern detected (confidence: 0.95)
  10:15:31  GC triggered by MemoryPredictor
  10:15:32  GC pause: 450ms (SLO target: <100ms)
  10:15:33  CPU priority boost applied [+5sec]
  10:15:34  Cache switch: LRU → Predictive
  10:15:45  GC pause normalized
  10:15:46  System recovered

Root Cause (자동 분석):
  - Primary: Memory fragmentation > 0.6
  - Trigger: MemoryLeak pattern × 0.95 confidence
  - Propagation: Memory pressure → CPU → I/O chain reaction

Impact:
  - API latency P99: 580ms (SLO: 100ms) → +480% ✗
  - Error rate: 0.2% (SLO: 0.1%) ✗
  - Cache hit rate: 62% (SLO: 80%) ✗
  - Error budget consumed: 3 minutes (5%)

Lessons Learned (자동):
  1. Fragmentation threshold should be < 0.5 (current: 0.6)
  2. GC urgency scoring needs Phase 6 confidence × 1.5
  3. Cascading prevention: I/O queue should throttle earlier

Prevention (자동 권고):
  - Increase memory compaction frequency (every 5min → every 2min)
  - Lower memory pressure Critical threshold (80% → 75%)
  - Add early warning: fragmentation > 0.4 → preemptive compact
  - Monitor Phase 6 confidence score < 0.8 → conservative policies

Follow-up Actions:
  □ Test increased compaction frequency (Week of 03-10)
  □ Implement early fragmentation detection (PATCH)
  □ Add metric: "phase6_confidence_weight_factor"
  □ Review similar incidents in last 30 days
```

**구현 내용**:
```rust
pub struct PostMortemAnalysis {
    incident: Incident,
    timeline: Vec<(u64, String, MetricsSnapshot)>,
    root_causes: Vec<(String, f64)>,  // (cause, confidence)
    impact_metrics: HashMap<String, (f64, f64)>,  // (actual, slo)
    lessons_learned: Vec<String>,
    prevention_measures: Vec<String>,
}

impl PostMortemAnalysis {
    pub fn extract_timeline(&self) → Vec<TimestepEntry>
    pub fn find_root_causes(&self) → Vec<(String, f64)>
    pub fn calculate_impact(&self) → HashMap<String, f64>
    pub fn generate_lessons(&self) → Vec<String>
    pub fn recommend_prevention(&self) → Vec<String>
    pub fn export_markdown(&self) → String
    pub fn suggest_config_changes(&self) → HashMap<String, String>  // 자동 설정 제안
}
```

**테스트** (8개):
- test_timeline_extraction
- test_root_cause_analysis
- test_impact_calculation
- test_lesson_generation
- test_prevention_recommendation
- test_config_suggestion
- test_markdown_export
- test_phase7_postmortem_integration

---

## 🔗 Phase H의 Phase 7 통합

```
Phase 7의 적응 결정 → Phase H의 추적

1. CPU Analyzer decision
   → distributed_tracer가 span 기록
   → metrics_collector가 decision_latency 측정
   → dashboard에 실시간 표시
   → SRO metrics에 포함 (decision_count, decision_latency)

2. Memory Predictor decision
   → tracer가 "prealloc_triggered" 로그
   → metrics에 "memory.prealloc_size_mb" 기록
   → dashboard에 "Memory decisions (last hour)" 표시
   → incident 발생 시 postmortem에 "Phase 7 decision" 포함

3. I/O Scheduler decision
   → tracer가 "cache_policy_switch" 로그
   → metrics에 "io.cache_hit_rate_change" 기록
   → SLO 측정에 포함 (cache_hit_rate SLO)
   → chaos injection 시 "I/O policy effectiveness" 검증
```

---

## 🎯 Success Criteria

| 항목 | 목표 | 증명 |
|------|------|------|
| **Tracing** | 100% 요청 추적 | trace_id가 없는 요청 0개 |
| **Metrics** | P99 latency < 100ms | distribution 그래프로 확인 |
| **Anomaly Detection** | false positive < 5% | historical data 30일 |
| **Chaos Recovery** | 자동 복구 95% | 50개 scenario, 47~50 성공 |
| **SLO 달성** | 99.9% uptime | 30일 연속 측정 |
| **Post-Mortem** | 자동화율 100% | 모든 incident의 automated report |
| **Deployment Safety** | error_budget 기반 결정 | 배포 결정이 100% SRE policy 준수 |

---

## 📅 구현 타임라인

### Week 1 (Days 1-7): 분산 트레이싱 & 모니터링
- Day 1-2: distributed_tracer.fl 설계 & 구현 (350줄)
- Day 3-4: metrics_collector.fl 구현 (350줄)
- Day 5-6: monitoring_dashboard.fl 구현 (300줄)
- Day 7: 통합 & 테스트 (모든 24개 tests PASS)

### Week 2 (Days 8-14): SRE 운영 & Chaos 검증
- Day 8-9: sre_operations.fl 구현 (350줄)
- Day 10-11: chaos_real_injection.fl 구현 (350줄)
- Day 12-13: postmortem_analyzer.fl 구현 (300줄)
- Day 14: 통합 & 최종 검증 (모든 24개 tests PASS)

### Deliverables
- 6개 모듈, 2,000줄 코드
- 48개 테스트 (100% PASS)
- Phase 7과의 완전한 통합
- 운영 가능한 대시보드
- 자동 post-mortem 생성 시스템

---

## 💡 핵심 철학 전환

**Phase 7까지**: "내 코드가 완벽하다"
**Phase H부터**: "내 코드가 실제 프로덕션에서 무중단으로 운영된다"

기록이 증명이다 (Your record is your proof).
운영이 증명이다 (Your operations prove it). 🎯

---

**다음 명령**: Week 1 Day 1 시작 ✅
