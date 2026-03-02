# Phase H: Observability 심화 - 분산 트레이싱 시스템

**목표**: 분산 시스템의 모든 요청이 "어디서 얼마나 지연되는지" 실시간으로 추적

**기간**: 1주
**산출물**: 1,500줄 (Rust)
**핵심 문제**: "원인 모를 지연"을 제거한다

---

## 🎯 설계 목표

### 1. Trace ID 전역 추적

```
Request 1
├─ TraceID: abc-def-123
├─ Service: API
│  ├─ Span 1: "auth" (5ms)
│  ├─ Span 2: "validate" (2ms)
│  └─ Span 3: "raft_consensus" (50ms)
│     ├─ Node 1: "log_append" (15ms)
│     ├─ Node 2: "log_append" (12ms)
│     └─ Node 3: "log_append" (10ms)
└─ Total: 57ms
```

### 2. 4가지 추적 대상

**A. HTTP 요청 전체 경로**
```
HTTP Request
→ Auth (JWT 검증) [5ms]
→ Input Validation [2ms]
→ Business Logic [10ms]
→ Raft Consensus [50ms]
→ HTTP Response [1ms]
---
Total: 68ms
```

**B. Raft 로그 복제**
```
Leader.append_entry()
├─ Local WAL write [3ms]
├─ RPC to Node 1 [10ms]
│  └─ Node 1.append_entry() [8ms]
├─ RPC to Node 2 [12ms]
│  └─ Node 2.append_entry() [9ms]
└─ Commit [5ms]
---
Total: 47ms
```

**C. 2PC 트랜잭션**
```
2PC.begin()
├─ Phase 1: Prepare [20ms]
│  ├─ Node 1 [5ms]
│  ├─ Node 2 [8ms]
│  └─ Node 3 [7ms]
└─ Phase 2: Commit [15ms]
   ├─ Node 1 [4ms]
   ├─ Node 2 [5ms]
   └─ Node 3 [6ms]
---
Total: 35ms (병렬)
```

**D. 은행 연산**
```
transfer(from, to, amount)
├─ From.lock() [1ms]
├─ To.lock() [0.5ms]
├─ From.debit() [0.5ms]
├─ To.credit() [0.5ms]
├─ Raft.replicate() [50ms]
└─ Unlock [0.5ms]
---
Total: 53ms
```

---

## 📐 아키텍처

### Core Components

```rust
// 1. TraceID 생성 (UUID 기반)
pub struct TraceId {
    id: String,          // "abc-def-123"
    start_time: u64,     // Unix timestamp (us)
    root_span: String,   // "http_request"
}

// 2. Span 기록 (각 작업의 시간 측정)
pub struct Span {
    trace_id: String,
    span_id: String,
    parent_span_id: Option<String>,
    operation: String,   // "raft_append", "jwt_verify", etc
    service: String,     // "api", "raft_node_1", etc
    start_time: u64,
    end_time: u64,
    tags: HashMap<String, String>,
    status: SpanStatus,  // OK, ERROR, SLOW
}

// 3. 추적 수집기
pub struct TracingCollector {
    spans: Arc<DashMap<String, Vec<Span>>>,
    slow_threshold_ms: u64,  // 병목 판단 기준
}
```

### 통합 포인트

```rust
// HTTP Request 시작
http_handler()
├─ trace_id = TraceId::new()
├─ ctx = TracingContext::new(trace_id.clone())
└─ handle_with_tracing(req, ctx)

// Raft Consensus 통합
raft_consensus()
├─ parent_span = ctx.current_span()
├─ new_span = ctx.start_span("raft_consensus")
├─ replicate_log()
│  ├─ new_span = ctx.start_span("replicate_to_node_1")
│  ├─ measure_time!(node1_rpc)
│  └─ ctx.end_span()
└─ ctx.end_span()

// 응답 전송
response()
├─ trace = ctx.finalize()
└─ response.headers["X-Trace-ID"] = trace_id
```

---

## 📊 구현 계획 (5단계)

### Stage 1: 기본 추적 인프라 (300줄)

**파일**: `src/tracing/core.rs`

```rust
// 1. TraceId 생성 및 전파
pub fn generate_trace_id() -> String {
    format!("{}-{}-{}",
        Uuid::new_v4(),
        SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_micros(),
        rand::random::<u16>()
    )
}

// 2. 스레드 로컬 Trace Context
thread_local! {
    static TRACE_CONTEXT: RefCell<Option<TracingContext>> = RefCell::new(None);
}

pub fn set_trace_context(ctx: TracingContext) {
    TRACE_CONTEXT.with(|tc| {
        *tc.borrow_mut() = Some(ctx);
    });
}

// 3. Span 시작/종료
pub struct Span {
    trace_id: String,
    span_id: String,
    operation: String,
    start_time: Instant,
    tags: DashMap<String, String>,
}

impl Span {
    pub fn new(trace_id: String, operation: String) -> Self {
        Span {
            trace_id,
            span_id: Uuid::new_v4().to_string(),
            operation,
            start_time: Instant::now(),
            tags: DashMap::new(),
        }
    }

    pub fn end(self) -> SpanRecord {
        let duration_ms = self.start_time.elapsed().as_millis() as u64;
        SpanRecord {
            trace_id: self.trace_id,
            span_id: self.span_id,
            operation: self.operation,
            duration_ms,
            tags: self.tags.into_iter().collect(),
        }
    }
}
```

### Stage 2: 로그 통합 (300줄)

**파일**: `src/tracing/logging.rs`

```rust
// 모든 로그에 자동으로 Trace ID 추가
pub struct TracedLogger;

impl TracedLogger {
    pub fn info(msg: &str) {
        let trace_id = get_trace_id();
        println!("[TRACE: {}] {}", trace_id, msg);
    }

    pub fn error(msg: &str) {
        let trace_id = get_trace_id();
        println!("[ERROR] [TRACE: {}] {}", trace_id, msg);
    }

    pub fn debug_span(operation: &str) {
        let trace_id = get_trace_id();
        println!("[TRACE: {}] >>> {}", trace_id, operation);
    }
}

// 매크로: 자동 시간 측정
#[macro_export]
macro_rules! measure_span {
    ($operation:expr, $code:block) => {{
        let span = Span::new(get_trace_id(), $operation.to_string());
        let result = $code;
        let record = span.end();
        record.log();
        result
    }};
}
```

### Stage 3: Raft 추적 (400줄)

**파일**: `src/tracing/raft_tracing.rs`

```rust
// Raft 로그 복제 전체 과정 추적
pub struct RaftTracer {
    collector: Arc<TracingCollector>,
}

impl RaftTracer {
    pub async fn trace_replicate_log(&self, entries: Vec<LogEntry>) {
        let trace_id = get_trace_id();
        let parent_span = Span::new(trace_id.clone(), "raft_replicate".to_string());

        // 로컬 쓰기
        {
            let local_span = Span::new(trace_id.clone(), "local_wal_write".to_string());
            self.local_wal_write(&entries).await;
            let record = local_span.end();
            record.log(); // 5ms 정도
        }

        // 원격 노드 복제 (병렬)
        {
            let node1_span = Span::new(trace_id.clone(), "replicate_to_node1".to_string());
            let node1_rpc = self.replicate_to_node(1, &entries);

            let node2_span = Span::new(trace_id.clone(), "replicate_to_node2".to_string());
            let node2_rpc = self.replicate_to_node(2, &entries);

            tokio::try_join!(node1_rpc, node2_rpc).ok();
        }

        // 커밋
        {
            let commit_span = Span::new(trace_id.clone(), "commit".to_string());
            self.commit(&entries).await;
        }

        parent_span.end().log();
    }

    async fn replicate_to_node(&self, node_id: usize, entries: &[LogEntry]) {
        // RPC 호출 + 원격 서버 처리 시간 측정
        let start = Instant::now();
        let response = self.send_rpc(node_id, entries).await;
        let remote_time = response.processing_time_us;
        let network_time = start.elapsed().as_micros() as u64 - remote_time;

        println!(
            "[TRACE] Node {} - Network: {}us, Remote: {}us",
            node_id, network_time, remote_time
        );
    }
}
```

### Stage 4: 2PC 추적 (300줄)

**파일**: `src/tracing/transaction_tracing.rs`

```rust
// 2PC 트랜잭션의 각 단계 추적
pub struct TransactionTracer {
    collector: Arc<TracingCollector>,
}

impl TransactionTracer {
    pub async fn trace_transaction(&self, tx: &Transaction) {
        let trace_id = get_trace_id();
        let parent_span = Span::new(trace_id.clone(), "2pc_transaction".to_string());

        // Phase 1: Prepare (병렬)
        {
            let prepare_span = Span::new(trace_id.clone(), "phase1_prepare".to_string());
            let prepare_tasks: Vec<_> = (1..=3)
                .map(|node_id| {
                    let span = Span::new(trace_id.clone(), format!("prepare_node_{}", node_id));
                    async move {
                        self.send_prepare(node_id, tx).await;
                        span.end().log();
                    }
                })
                .collect();

            futures::future::join_all(prepare_tasks).await;
            prepare_span.end().log();
        }

        // Phase 2: Commit (병렬)
        {
            let commit_span = Span::new(trace_id.clone(), "phase2_commit".to_string());
            let commit_tasks: Vec<_> = (1..=3)
                .map(|node_id| {
                    let span = Span::new(trace_id.clone(), format!("commit_node_{}", node_id));
                    async move {
                        self.send_commit(node_id, tx).await;
                        span.end().log();
                    }
                })
                .collect();

            futures::future::join_all(commit_tasks).await;
            commit_span.end().log();
        }

        parent_span.end().log();
    }
}
```

### Stage 5: 분석 & 대시보드 (200줄)

**파일**: `src/tracing/analytics.rs`

```rust
// 수집된 Trace 분석
pub struct TraceAnalytics {
    spans: Arc<DashMap<String, Vec<Span>>>,
}

impl TraceAnalytics {
    // 1. 병목 지점 자동 감지
    pub fn find_bottlenecks(&self, threshold_ms: u64) -> Vec<BottleneckReport> {
        let mut bottlenecks = Vec::new();

        for span_vec in self.spans.iter() {
            for span in span_vec.value() {
                if span.duration_ms > threshold_ms {
                    bottlenecks.push(BottleneckReport {
                        operation: span.operation.clone(),
                        duration_ms: span.duration_ms,
                        severity: if span.duration_ms > threshold_ms * 2 {
                            "CRITICAL"
                        } else {
                            "WARNING"
                        },
                    });
                }
            }
        }

        bottlenecks.sort_by(|a, b| b.duration_ms.cmp(&a.duration_ms));
        bottlenecks
    }

    // 2. 평균 지연 계산
    pub fn average_latency(&self, operation: &str) -> u64 {
        let mut total = 0u64;
        let mut count = 0u64;

        for span_vec in self.spans.iter() {
            for span in span_vec.value() {
                if span.operation == operation {
                    total += span.duration_ms;
                    count += 1;
                }
            }
        }

        if count == 0 { 0 } else { total / count }
    }

    // 3. P50/P95/P99 지연
    pub fn percentile_latency(&self, operation: &str, percentile: u8) -> u64 {
        let mut latencies: Vec<u64> = Vec::new();

        for span_vec in self.spans.iter() {
            for span in span_vec.value() {
                if span.operation == operation {
                    latencies.push(span.duration_ms);
                }
            }
        }

        latencies.sort();
        let index = (latencies.len() * percentile as usize) / 100;
        latencies.get(index).copied().unwrap_or(0)
    }
}

// REST API 엔드포인트
pub async fn trace_api_handler(
    axum::extract::Query(params): Query<TraceQuery>,
) -> String {
    let analytics = get_analytics();

    match params.action.as_str() {
        "bottlenecks" => {
            let bottlenecks = analytics.find_bottlenecks(params.threshold);
            format!("{:#?}", bottlenecks)
        }
        "avg_latency" => {
            let avg = analytics.average_latency(&params.operation);
            format!("Average latency for {}: {}ms", params.operation, avg)
        }
        "percentile" => {
            let p99 = analytics.percentile_latency(&params.operation, 99);
            format!("P99 latency for {}: {}ms", params.operation, p99)
        }
        _ => "Unknown action".to_string(),
    }
}
```

---

## 📈 기대 효과

### 1. 병목 지점 실시간 감지

```
[TRACE] GET /transfer HTTP/1.1
├─ auth: 5ms ✓
├─ validate: 2ms ✓
├─ raft_consensus: 150ms ❌ CRITICAL
│  ├─ local_wal: 3ms ✓
│  ├─ replicate_node1: 45ms ⚠️ (network 35ms + remote 10ms)
│  ├─ replicate_node2: 52ms ❌ (network 40ms + remote 12ms)
│  └─ commit: 50ms ⚠️
└─ response: 1ms ✓

RECOMMENDATION: Node 2 network latency 40ms - check connectivity
```

### 2. 성능 회귀 감지

```
Before: avg raft_consensus = 50ms
After: avg raft_consensus = 150ms (+200%)
Action: Automatically alert if latency increases > 20%
```

### 3. 자동 병목 보고서

```
/metrics/traces?action=bottlenecks&threshold=50
→ 병목이 되는 작업들을 자동으로 리스트업
→ P99 지연이 높은 연산 자동 감지
→ 원격 노드 네트워크 지연 분석
```

---

## ✅ 검증 방법

### 테스트 시나리오

```bash
# 1. 정상 요청
$ curl -H "X-Trace-ID: abc-123" http://localhost:8080/transfer
→ Trace 출력: 모든 span 순차 기록

# 2. 느린 네트워크 시뮬레이션
$ tc qdisc add dev eth0 root netem delay 100ms
→ Trace 출력: network latency 자동 감지

# 3. 노드 다운 시뮬레이션
$ docker stop node2
→ Trace 출력: node2_replicate timeout 감지
```

---

## 🎯 최종 구조

```
Tracing System
├── Core (core.rs) - TraceId, Span, TracingContext
├── Logging (logging.rs) - 로그 통합, measure! 매크로
├── Raft Tracing (raft_tracing.rs) - Raft 로그 복제 추적
├── Transaction Tracing (transaction_tracing.rs) - 2PC 추적
├── Analytics (analytics.rs) - 병목 감지, P99 계산
└── API (api.rs) - REST 엔드포인트

→ 모든 요청이 TraceID로 추적 가능
→ 병목 지점 자동 감지
→ 성능 회귀 자동 알람
```

---

**다음**: Phase I (Chaos Engineering) - 실제 장애에서의 복구 검증
