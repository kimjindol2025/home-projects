/// Phase H: Observability - Distributed Tracing System
/// 분산 시스템의 모든 요청에 TraceID를 부여하고 경로별 지연 측정

pub mod core;
pub mod logging;
pub mod raft_tracing;
pub mod transaction_tracing;
pub mod analytics;

pub use core::{
    TraceId, Span, SpanRecord, SpanStatus, TracingContext, TraceReport,
    set_trace_context, get_trace_context, get_trace_id, clear_trace_context,
};

pub use logging::{
    TracedLogger, PerformanceLogger, PerformanceThresholds,
    ErrorLogger, NetworkLogger,
};

pub use raft_tracing::{RaftTracer, RaftTracingMetrics};
pub use transaction_tracing::{TransactionTracer, TransactionTracingMetrics};
pub use analytics::{TraceAnalytics, BottleneckReport};
