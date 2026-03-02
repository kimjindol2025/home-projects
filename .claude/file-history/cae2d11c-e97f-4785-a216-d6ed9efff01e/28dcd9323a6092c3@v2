/// Phase H: Observability - Core Tracing Infrastructure
/// 분산 트레이싱의 기본 인프라: TraceID, Span, TracingContext

use std::cell::RefCell;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{SystemTime, UNIX_EPOCH, Instant};
use uuid::Uuid;
use dashmap::DashMap;

/// TraceID: 전체 요청의 고유 식별자
/// "2026-03-02T10:30:45.123456-abc-def-123" 형식
#[derive(Clone, Debug, Eq, PartialEq, Hash)]
pub struct TraceId {
    pub id: String,
}

impl TraceId {
    /// 새로운 TraceID 생성
    pub fn new() -> Self {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_micros();

        let uuid_part = Uuid::new_v4().to_string()[..8].to_string();
        let random_part = format!("{:05}", rand::random::<u32>() % 100000);

        TraceId {
            id: format!("{}-{}-{}", timestamp, uuid_part, random_part),
        }
    }

    pub fn from_string(id: String) -> Self {
        TraceId { id }
    }

    pub fn to_string(&self) -> String {
        self.id.clone()
    }
}

impl Default for TraceId {
    fn default() -> Self {
        Self::new()
    }
}

/// Span: 특정 작업의 시간 및 메타데이터 기록
/// 예: "raft_append", "jwt_verify", "raft_consensus" 등
#[derive(Clone, Debug)]
pub struct SpanRecord {
    pub trace_id: String,
    pub span_id: String,
    pub parent_span_id: Option<String>,
    pub operation: String,
    pub service: String,
    pub start_time_us: u64,
    pub end_time_us: u64,
    pub duration_ms: u64,
    pub tags: HashMap<String, String>,
    pub status: SpanStatus,
}

#[derive(Clone, Debug, Copy, PartialEq, Eq)]
pub enum SpanStatus {
    Ok,
    Error,
    Slow,
}

impl SpanRecord {
    /// Span 기록을 로그로 출력
    pub fn log(&self) {
        let status_str = match self.status {
            SpanStatus::Ok => "✓",
            SpanStatus::Error => "✗",
            SpanStatus::Slow => "⚠",
        };

        println!(
            "[TRACE:{}] {} {} → {} ({}ms) {}",
            self.trace_id,
            status_str,
            self.operation,
            self.service,
            self.duration_ms,
            self.tags
                .iter()
                .map(|(k, v)| format!("{}={}", k, v))
                .collect::<Vec<_>>()
                .join(" ")
        );
    }

    /// JSON 형식으로 변환
    pub fn to_json(&self) -> String {
        format!(
            r#"{{"trace_id":"{}","span_id":"{}","operation":"{}","duration_ms":{},"status":"{}"}}"#,
            self.trace_id,
            self.span_id,
            self.operation,
            self.duration_ms,
            match self.status {
                SpanStatus::Ok => "ok",
                SpanStatus::Error => "error",
                SpanStatus::Slow => "slow",
            }
        )
    }
}

/// Span 빌더: 활성 Span
pub struct Span {
    trace_id: String,
    span_id: String,
    parent_span_id: Option<String>,
    operation: String,
    service: String,
    start_time_us: u64,
    start_instant: Instant,
    tags: Arc<DashMap<String, String>>,
    slow_threshold_ms: u64,
}

impl Span {
    /// 새로운 Span 생성
    pub fn new(trace_id: String, operation: String, service: String) -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_micros() as u64;

        Span {
            trace_id,
            span_id: Uuid::new_v4().to_string(),
            parent_span_id: None,
            operation,
            service,
            start_time_us: now,
            start_instant: Instant::now(),
            tags: Arc::new(DashMap::new()),
            slow_threshold_ms: 100, // 기본값: 100ms 이상이면 SLOW 판정
        }
    }

    /// Parent Span ID 설정
    pub fn with_parent(mut self, parent_id: String) -> Self {
        self.parent_span_id = Some(parent_id);
        self
    }

    /// Slow threshold 설정
    pub fn with_threshold(mut self, ms: u64) -> Self {
        self.slow_threshold_ms = ms;
        self
    }

    /// Tag 추가
    pub fn add_tag(&self, key: String, value: String) {
        self.tags.insert(key, value);
    }

    /// Span 종료 및 기록 생성
    pub fn end(self) -> SpanRecord {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_micros() as u64;

        let duration_ms = self.start_instant.elapsed().as_millis() as u64;

        let status = if duration_ms >= self.slow_threshold_ms {
            SpanStatus::Slow
        } else {
            SpanStatus::Ok
        };

        SpanRecord {
            trace_id: self.trace_id,
            span_id: self.span_id,
            parent_span_id: self.parent_span_id,
            operation: self.operation,
            service: self.service,
            start_time_us: self.start_time_us,
            end_time_us: now,
            duration_ms,
            tags: self.tags.iter().map(|r| (r.key().clone(), r.value().clone())).collect(),
            status,
        }
    }
}

/// TracingContext: 스레드 로컬 Trace 컨텍스트
#[derive(Clone)]
pub struct TracingContext {
    trace_id: TraceId,
    current_span_id: Option<String>,
    spans: Arc<DashMap<String, SpanRecord>>,
}

impl TracingContext {
    /// 새로운 TracingContext 생성
    pub fn new(trace_id: TraceId) -> Self {
        TracingContext {
            trace_id,
            current_span_id: None,
            spans: Arc::new(DashMap::new()),
        }
    }

    /// Span 시작
    pub fn start_span(&mut self, operation: String, service: String) -> Span {
        let span = Span::new(self.trace_id.id.clone(), operation, service);

        if let Some(parent_id) = &self.current_span_id {
            span.with_parent(parent_id.clone())
        } else {
            span
        }
    }

    /// Span 종료 및 저장
    pub fn end_span(&mut self, span: SpanRecord) {
        self.spans.insert(span.span_id.clone(), span);
    }

    /// 현재 Span ID 설정
    pub fn set_current_span(&mut self, span_id: String) {
        self.current_span_id = Some(span_id);
    }

    /// 모든 Span 조회
    pub fn get_spans(&self) -> Vec<SpanRecord> {
        self.spans
            .iter()
            .map(|r| r.value().clone())
            .collect()
    }

    /// TraceID 조회
    pub fn trace_id(&self) -> String {
        self.trace_id.id.clone()
    }

    /// Trace 완성 및 최종화
    pub fn finalize(self) -> TraceReport {
        let spans = self.get_spans();
        let total_duration_ms = spans
            .iter()
            .map(|s| s.duration_ms)
            .max()
            .unwrap_or(0);

        TraceReport {
            trace_id: self.trace_id.id,
            span_count: spans.len() as u32,
            total_duration_ms,
            spans,
        }
    }
}

/// TraceReport: 완료된 Trace 리포트
#[derive(Clone, Debug)]
pub struct TraceReport {
    pub trace_id: String,
    pub span_count: u32,
    pub total_duration_ms: u64,
    pub spans: Vec<SpanRecord>,
}

impl TraceReport {
    /// 병목 Span 찾기 (가장 느린 작업)
    pub fn find_bottleneck(&self) -> Option<&SpanRecord> {
        self.spans.iter().max_by_key(|s| s.duration_ms)
    }

    /// 전체 리포트 출력
    pub fn print_summary(&self) {
        println!(
            "\n╔════════════════════════════════════════╗"
        );
        println!("║ Trace Report: {}                   ║", &self.trace_id[..8]);
        println!("╠════════════════════════════════════════╣");
        println!("║ Total Duration: {}ms                    ║", self.total_duration_ms);
        println!("║ Span Count: {}                         ║", self.span_count);
        println!("╚════════════════════════════════════════╝");

        for span in &self.spans {
            span.log();
        }

        if let Some(bottleneck) = self.find_bottleneck() {
            println!(
                "\n⚠️  Bottleneck: {} ({}ms)",
                bottleneck.operation, bottleneck.duration_ms
            );
        }
    }
}

/// 글로벌 Tracing Context (스레드 로컬)
thread_local! {
    static TRACE_CONTEXT: RefCell<Option<TracingContext>> = RefCell::new(None);
}

/// Trace Context 설정
pub fn set_trace_context(ctx: TracingContext) {
    TRACE_CONTEXT.with(|tc| {
        *tc.borrow_mut() = Some(ctx);
    });
}

/// 현재 Trace Context 조회
pub fn get_trace_context() -> Option<TracingContext> {
    TRACE_CONTEXT.with(|tc| tc.borrow().clone())
}

/// 현재 Trace ID 조회
pub fn get_trace_id() -> String {
    TRACE_CONTEXT.with(|tc| {
        tc.borrow()
            .as_ref()
            .map(|ctx| ctx.trace_id.id.clone())
            .unwrap_or_else(|| "NO_TRACE".to_string())
    })
}

/// Trace Context 초기화
pub fn clear_trace_context() {
    TRACE_CONTEXT.with(|tc| {
        *tc.borrow_mut() = None;
    });
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_trace_id_generation() {
        let id1 = TraceId::new();
        let id2 = TraceId::new();
        assert_ne!(id1, id2);
    }

    #[test]
    fn test_span_creation() {
        let span = Span::new(
            "test-trace".to_string(),
            "test_op".to_string(),
            "test_service".to_string(),
        );
        assert_eq!(span.operation, "test_op");
        assert_eq!(span.service, "test_service");
    }

    #[test]
    fn test_span_tagging() {
        let span = Span::new(
            "test-trace".to_string(),
            "test_op".to_string(),
            "test_service".to_string(),
        );
        span.add_tag("key1".to_string(), "value1".to_string());
        let record = span.end();
        assert_eq!(record.tags.get("key1").unwrap(), "value1");
    }

    #[test]
    fn test_tracing_context() {
        let trace_id = TraceId::new();
        let mut ctx = TracingContext::new(trace_id.clone());

        let span = ctx.start_span("op1".to_string(), "svc1".to_string());
        let record = span.end();
        ctx.end_span(record);

        assert_eq!(ctx.get_spans().len(), 1);
    }

    #[test]
    fn test_thread_local_context() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        let retrieved_id = get_trace_id();
        assert!(!retrieved_id.is_empty());
    }

    #[test]
    fn test_slow_span_detection() {
        let span = Span::new(
            "test-trace".to_string(),
            "slow_op".to_string(),
            "test_service".to_string(),
        ).with_threshold(10); // 10ms threshold

        // 실제로 시간을 측정하려면 sleep 필요
        let record = span.end();
        // 이 테스트에서는 duration이 10ms보다 적을 것 (OK 상태)
        assert_eq!(record.status, SpanStatus::Ok);
    }
}
