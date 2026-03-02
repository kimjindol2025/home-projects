/// Phase H: Observability - Logging Integration
/// 모든 로그에 자동으로 TraceID 추가

use super::core::{get_trace_id, Span};
use tracing::{info, warn, error};
use std::fmt;
use tracing::{info, warn, error};
use std::time::Instant;
use tracing::{info, warn, error};

/// TracedLogger: TraceID를 포함한 로거
pub struct TracedLogger;

impl TracedLogger {
    /// Info 레벨 로그 (TraceID 자동 포함)
    pub fn info(msg: &str) {
        let trace_id = get_trace_id();
        info!("[INFO] [TRACE:{}] {}", trace_id, msg);
    }

    /// Warning 레벨 로그
    pub fn warn(msg: &str) {
        let trace_id = get_trace_id();
        info!("[WARN] [TRACE:{}] {}", trace_id, msg);
    }

    /// Error 레벨 로그
    pub fn error(msg: &str) {
        let trace_id = get_trace_id();
        info!("[ERROR] [TRACE:{}] {}", trace_id, msg);
    }

    /// Debug 레벨 로그
    pub fn debug(msg: &str) {
        let trace_id = get_trace_id();
        info!("[DEBUG] [TRACE:{}] {}", trace_id, msg);
    }

    /// Span 진입 로그
    pub fn debug_enter(operation: &str) {
        let trace_id = get_trace_id();
        info!("[TRACE:{}] >>> {}", trace_id, operation);
    }

    /// Span 종료 로그
    pub fn debug_exit(operation: &str, duration_ms: u64) {
        let trace_id = get_trace_id();
        info!("[TRACE:{}] <<< {} ({}ms)", trace_id, operation, duration_ms);
    }

    /// 구조화된 로그 (JSON 형식)
    pub fn structured(level: &str, operation: &str, data: &[(String, String)]) {
        let trace_id = get_trace_id();
        let mut json = format!(r#"{{"level":"{}","trace_id":"{}","operation":"{}""#, level, trace_id, operation);

        for (key, value) in data {
            json.push_str(&format!(r#","{}":"{}""#, key, value));
        }
        json.push('}');

        info!("{}", json);
    }
}

/// 자동 시간 측정 매크로
/// 사용: measure_span!("operation_name", { ... code ... })
#[macro_export]
macro_rules! measure_span {
    ($operation:expr, $service:expr, $code:block) => {{
        let start = Instant::now();
        let trace_id = $crate::tracing::core::get_trace_id();

        $crate::tracing::logging::TracedLogger::debug_enter($operation);

        let result = $code;

        let duration = start.elapsed().as_millis() as u64;
        $crate::tracing::logging::TracedLogger::debug_exit($operation, duration);

        result
    }};
}

/// 자동 시간 측정 매크로 (간단한 버전)
#[macro_export]
macro_rules! measure {
    ($operation:expr, $code:block) => {{
        let start = std::time::Instant::now();
        let result = $code;
        let duration = start.elapsed().as_millis() as u64;

        if duration > 100 {
            einfo!("[SLOW] {} took {}ms", $operation, duration);
        }

        result
    }};
}

/// 성능 메트릭 기록기
pub struct PerformanceLogger {
    thresholds: PerformanceThresholds,
}

pub struct PerformanceThresholds {
    pub warn_ms: u64,    // 이 이상이면 경고
    pub error_ms: u64,   // 이 이상이면 에러
}

impl Default for PerformanceThresholds {
    fn default() -> Self {
        PerformanceThresholds {
            warn_ms: 50,   // 50ms 이상 경고
            error_ms: 200, // 200ms 이상 에러
        }
    }
}

impl PerformanceLogger {
    pub fn new() -> Self {
        PerformanceLogger {
            thresholds: PerformanceThresholds::default(),
        }
    }

    pub fn with_thresholds(warn_ms: u64, error_ms: u64) -> Self {
        PerformanceLogger {
            thresholds: PerformanceThresholds { warn_ms, error_ms },
        }
    }

    /// 작업 실행 및 성능 로깅
    pub fn log_operation<F, R>(&self, operation: &str, f: F) -> R
    where
        F: FnOnce() -> R,
    {
        let trace_id = get_trace_id();
        let start = Instant::now();

        let result = f();

        let duration_ms = start.elapsed().as_millis() as u64;

        if duration_ms >= self.thresholds.error_ms {
            info!(
                "[ERROR] [TRACE:{}] {} took {}ms (> {}ms threshold)",
                trace_id, operation, duration_ms, self.thresholds.error_ms
            );
        } else if duration_ms >= self.thresholds.warn_ms {
            info!(
                "[WARN] [TRACE:{}] {} took {}ms (> {}ms threshold)",
                trace_id, operation, duration_ms, self.thresholds.warn_ms
            );
        } else {
            info!(
                "[DEBUG] [TRACE:{}] {} took {}ms",
                trace_id, operation, duration_ms
            );
        }

        result
    }
}

/// 에러 로깅 유틸
pub struct ErrorLogger;

impl ErrorLogger {
    /// 에러를 TraceID와 함께 로깅
    pub fn log(error: impl fmt::Display, context: &str) {
        let trace_id = get_trace_id();
        info!(
            "[ERROR] [TRACE:{}] {}: {}",
            trace_id, context, error
        );
    }

    /// 복구 가능한 에러
    pub fn recoverable(error: impl fmt::Display, context: &str) {
        let trace_id = get_trace_id();
        info!(
            "[WARN] [TRACE:{}] [RECOVERABLE] {}: {}",
            trace_id, context, error
        );
    }

    /// 치명적인 에러
    pub fn fatal(error: impl fmt::Display, context: &str) {
        let trace_id = get_trace_id();
        info!(
            "[FATAL] [TRACE:{}] {}: {}",
            trace_id, context, error
        );
    }
}

/// 네트워크 요청 로깅
pub struct NetworkLogger;

impl NetworkLogger {
    /// RPC 요청 로깅
    pub fn log_rpc_request(service: &str, method: &str, target_node: u32) {
        let trace_id = get_trace_id();
        info!(
            "[RPC] [TRACE:{}] {} → {} ({})",
            trace_id, service, method, target_node
        );
    }

    /// RPC 응답 로깅
    pub fn log_rpc_response(service: &str, method: &str, duration_ms: u64, status: &str) {
        let trace_id = get_trace_id();
        info!(
            "[RPC] [TRACE:{}] {} ← {} ({}ms) [{}]",
            trace_id, service, method, duration_ms, status
        );
    }

    /// 네트워크 지연 감지
    pub fn log_latency_warning(node: u32, latency_ms: u64, threshold_ms: u64) {
        let trace_id = get_trace_id();
        info!(
            "[WARN] [TRACE:{}] High latency to Node {}: {}ms (threshold: {}ms)",
            trace_id, node, latency_ms, threshold_ms
        );
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tracing::core::{set_trace_context, TracingContext, TraceId};

    #[test]
    fn test_traced_logger_info() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        // Should not panic
        TracedLogger::info("Test info message");
    }

    #[test]
    fn test_traced_logger_error() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        // Should not panic
        TracedLogger::error("Test error message");
    }

    #[test]
    fn test_performance_logger() {
        let logger = PerformanceLogger::new();
        let result = logger.log_operation("test_op", || {
            std::thread::sleep(std::time::Duration::from_millis(10));
            42
        });
        assert_eq!(result, 42);
    }

    #[test]
    fn test_error_logger() {
        // Should not panic
        ErrorLogger::log("Test error", "test_context");
        ErrorLogger::recoverable("Test recoverable", "test_context");
        ErrorLogger::fatal("Test fatal", "test_context");
    }

    #[test]
    fn test_network_logger() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        NetworkLogger::log_rpc_request("raft", "append_entries", 1);
        NetworkLogger::log_rpc_response("raft", "append_entries", 10, "ok");
        NetworkLogger::log_latency_warning(2, 150, 100);
    }

    #[test]
    fn test_structured_logging() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        let data = vec![
            ("status".to_string(), "ok".to_string()),
            ("duration_ms".to_string(), "42".to_string()),
        ];
        TracedLogger::structured("info", "test_op", &data);
    }
}
