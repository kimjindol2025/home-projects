/// Phase H: Observability - Trace Analytics
/// 수집된 Trace 데이터 분석 및 병목 지점 감지

use super::core::SpanRecord;
use std::collections::HashMap;

/// 병목 지점 리포트
#[derive(Clone, Debug)]
pub struct BottleneckReport {
    pub operation: String,
    pub duration_ms: u64,
    pub severity: String, // CRITICAL, WARNING, OK
    pub impact_percent: f64,
}

/// Trace 분석기
pub struct TraceAnalytics {
    spans: Vec<SpanRecord>,
    total_duration_ms: u64,
}

impl TraceAnalytics {
    /// 새로운 분석기 생성
    pub fn new(spans: Vec<SpanRecord>) -> Self {
        let total_duration_ms = spans
            .iter()
            .map(|s| s.duration_ms)
            .max()
            .unwrap_or(0);

        TraceAnalytics {
            spans,
            total_duration_ms,
        }
    }

    /// 병목 지점 감지 (threshold 이상의 작업)
    pub fn find_bottlenecks(&self, threshold_ms: u64) -> Vec<BottleneckReport> {
        let mut bottlenecks = Vec::new();

        for span in &self.spans {
            if span.duration_ms >= threshold_ms {
                let severity = if span.duration_ms >= threshold_ms * 2 {
                    "CRITICAL"
                } else if span.duration_ms >= threshold_ms {
                    "WARNING"
                } else {
                    "OK"
                };

                let impact_percent = if self.total_duration_ms > 0 {
                    (span.duration_ms as f64 / self.total_duration_ms as f64) * 100.0
                } else {
                    0.0
                };

                bottlenecks.push(BottleneckReport {
                    operation: span.operation.clone(),
                    duration_ms: span.duration_ms,
                    severity: severity.to_string(),
                    impact_percent,
                });
            }
        }

        // 지연 시간 순으로 정렬
        bottlenecks.sort_by(|a, b| b.duration_ms.cmp(&a.duration_ms));
        bottlenecks
    }

    /// 평균 지연 계산 (특정 작업)
    pub fn average_latency(&self, operation: &str) -> u64 {
        let matching: Vec<u64> = self.spans
            .iter()
            .filter(|s| s.operation == operation)
            .map(|s| s.duration_ms)
            .collect();

        if matching.is_empty() {
            0
        } else {
            matching.iter().sum::<u64>() / matching.len() as u64
        }
    }

    /// 백분위 지연 계산 (P50, P95, P99 등)
    pub fn percentile_latency(&self, operation: &str, percentile: u8) -> u64 {
        let mut latencies: Vec<u64> = self.spans
            .iter()
            .filter(|s| s.operation == operation)
            .map(|s| s.duration_ms)
            .collect();

        if latencies.is_empty() {
            return 0;
        }

        latencies.sort();
        let index = (latencies.len() * percentile as usize) / 100;
        latencies.get(index).copied().unwrap_or(0)
    }

    /// 최대 지연
    pub fn max_latency(&self, operation: &str) -> u64 {
        self.spans
            .iter()
            .filter(|s| s.operation == operation)
            .map(|s| s.duration_ms)
            .max()
            .unwrap_or(0)
    }

    /// 최소 지연
    pub fn min_latency(&self, operation: &str) -> u64 {
        self.spans
            .iter()
            .filter(|s| s.operation == operation)
            .map(|s| s.duration_ms)
            .min()
            .unwrap_or(0)
    }

    /// 작업별 통계
    pub fn operation_stats(&self) -> HashMap<String, OperationStats> {
        let mut stats_map: HashMap<String, Vec<u64>> = HashMap::new();

        for span in &self.spans {
            stats_map
                .entry(span.operation.clone())
                .or_insert_with(Vec::new)
                .push(span.duration_ms);
        }

        let mut result = HashMap::new();

        for (operation, latencies) in stats_map {
            let sum: u64 = latencies.iter().sum();
            let count = latencies.len();
            let avg = sum / count as u64;

            let mut sorted = latencies.clone();
            sorted.sort();

            let p50_idx = (sorted.len() * 50) / 100;
            let p95_idx = (sorted.len() * 95) / 100;
            let p99_idx = (sorted.len() * 99) / 100;

            result.insert(
                operation,
                OperationStats {
                    count: count as u32,
                    avg_ms: avg,
                    min_ms: *sorted.first().unwrap_or(&0),
                    max_ms: *sorted.last().unwrap_or(&0),
                    p50_ms: sorted.get(p50_idx).copied().unwrap_or(0),
                    p95_ms: sorted.get(p95_idx).copied().unwrap_or(0),
                    p99_ms: sorted.get(p99_idx).copied().unwrap_or(0),
                },
            );
        }

        result
    }

    /// 성능 회귀 감지
    pub fn detect_regression(&self, baseline_avg_ms: u64, threshold_percent: f64) -> Vec<String> {
        let current_avg = self.average_latency("raft_consensus");
        let regression_threshold = (baseline_avg_ms as f64 * (100.0 + threshold_percent) / 100.0) as u64;

        if current_avg > regression_threshold {
            vec![format!(
                "Performance regression detected: {}ms (baseline: {}ms, +{:.1}%)",
                current_avg,
                baseline_avg_ms,
                ((current_avg as f64 - baseline_avg_ms as f64) / baseline_avg_ms as f64) * 100.0
            )]
        } else {
            Vec::new()
        }
    }

    /// 전체 분석 리포트
    pub fn generate_report(&self) -> String {
        let mut report = String::new();

        report.push_str("═══════════════════════════════════════════\n");
        report.push_str("        Trace Analysis Report\n");
        report.push_str("═══════════════════════════════════════════\n\n");

        report.push_str(&format!("Total Spans: {}\n", self.spans.len()));
        report.push_str(&format!("Total Duration: {}ms\n\n", self.total_duration_ms));

        // 작업별 통계
        let stats = self.operation_stats();
        report.push_str("Operation Statistics:\n");
        report.push_str("─────────────────────────────────────────\n");

        for (op, stat) in stats {
            report.push_str(&format!("{}: ", op));
            report.push_str(&format!("count={} ", stat.count));
            report.push_str(&format!("avg={}ms ", stat.avg_ms));
            report.push_str(&format!("p99={}ms ", stat.p99_ms));
            report.push_str(&format!("max={}ms\n", stat.max_ms));
        }

        report.push_str("\n");

        // 병목 지점
        let bottlenecks = self.find_bottlenecks(50);
        if !bottlenecks.is_empty() {
            report.push_str("Bottlenecks (>50ms):\n");
            report.push_str("─────────────────────────────────────────\n");

            for bn in bottlenecks {
                report.push_str(&format!(
                    "  {} [{}] {}ms ({:.1}% of total)\n",
                    bn.operation, bn.severity, bn.duration_ms, bn.impact_percent
                ));
            }
        }

        report.push_str("\n═══════════════════════════════════════════\n");
        report
    }
}

/// 작업별 통계
#[derive(Clone, Debug)]
pub struct OperationStats {
    pub count: u32,
    pub avg_ms: u64,
    pub min_ms: u64,
    pub max_ms: u64,
    pub p50_ms: u64,
    pub p95_ms: u64,
    pub p99_ms: u64,
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tracing::core::{SpanRecord, SpanStatus};

    fn create_test_spans() -> Vec<SpanRecord> {
        vec![
            SpanRecord {
                trace_id: "test".to_string(),
                span_id: "1".to_string(),
                parent_span_id: None,
                operation: "auth".to_string(),
                service: "api".to_string(),
                start_time_us: 0,
                end_time_us: 5000,
                duration_ms: 5,
                tags: Default::default(),
                status: SpanStatus::Ok,
            },
            SpanRecord {
                trace_id: "test".to_string(),
                span_id: "2".to_string(),
                parent_span_id: None,
                operation: "raft_consensus".to_string(),
                service: "raft".to_string(),
                start_time_us: 5000,
                end_time_us: 55000,
                duration_ms: 50,
                tags: Default::default(),
                status: SpanStatus::Ok,
            },
        ]
    }

    #[test]
    fn test_analytics_creation() {
        let spans = create_test_spans();
        let analytics = TraceAnalytics::new(spans);
        assert_eq!(analytics.spans.len(), 2);
    }

    #[test]
    fn test_bottleneck_detection() {
        let spans = create_test_spans();
        let analytics = TraceAnalytics::new(spans);
        let bottlenecks = analytics.find_bottlenecks(40);
        assert_eq!(bottlenecks.len(), 1);
    }

    #[test]
    fn test_average_latency() {
        let spans = create_test_spans();
        let analytics = TraceAnalytics::new(spans);
        let avg = analytics.average_latency("auth");
        assert_eq!(avg, 5);
    }

    #[test]
    fn test_operation_stats() {
        let spans = create_test_spans();
        let analytics = TraceAnalytics::new(spans);
        let stats = analytics.operation_stats();
        assert_eq!(stats.len(), 2);
    }
}
