/// Phase H: Observability - Raft Consensus Tracing
/// Raft 로그 복제 경로의 전체 가시화

use super::core::{get_trace_id, Span, SpanRecord, TracingContext, get_trace_context};
use super::logging::{NetworkLogger, PerformanceLogger, TracedLogger};
use std::sync::Arc;
use std::time::Instant;
use dashmap::DashMap;

/// Raft 추적 메트릭
#[derive(Clone, Debug)]
pub struct RaftTracingMetrics {
    pub local_wal_time_us: u64,      // Local WAL 쓰기 시간
    pub node_latencies: Vec<(u32, u64)>, // 노드별 지연 (node_id, latency_us)
    pub total_replication_time_us: u64,
    pub commit_time_us: u64,
}

impl Default for RaftTracingMetrics {
    fn default() -> Self {
        RaftTracingMetrics {
            local_wal_time_us: 0,
            node_latencies: Vec::new(),
            total_replication_time_us: 0,
            commit_time_us: 0,
        }
    }
}

/// Raft 추적기
pub struct RaftTracer {
    metrics: Arc<DashMap<String, RaftTracingMetrics>>,
    perf_logger: PerformanceLogger,
}

impl RaftTracer {
    pub fn new() -> Self {
        RaftTracer {
            metrics: Arc::new(DashMap::new()),
            perf_logger: PerformanceLogger::with_thresholds(50, 200),
        }
    }

    /// Raft 로그 복제 추적 (전체 과정)
    /// 단계:
    /// 1. Local WAL 쓰기
    /// 2. 원격 노드로 복제 (병렬)
    /// 3. 커밋
    pub async fn trace_replicate_log(&self, entries: &[(u32, Vec<u8>)]) {
        let trace_id = get_trace_id();
        let parent_span = Span::new(
            trace_id.clone(),
            "raft_replicate_log".to_string(),
            "raft".to_string(),
        );

        let mut metrics = RaftTracingMetrics::default();
        let parent_start = Instant::now();

        // Stage 1: Local WAL 쓰기
        {
            let span = Span::new(
                trace_id.clone(),
                "local_wal_write".to_string(),
                "raft".to_string(),
            ).with_threshold(10);

            TracedLogger::debug_enter("local_wal_write");
            let local_start = Instant::now();

            // 시뮬레이션: WAL에 엔트리 쓰기
            let entry_count = entries.len();
            let total_size: usize = entries.iter().map(|(_, data)| data.len()).sum();

            // I/O 시뮬레이션 (실제로는 디스크 쓰기)
            tokio::time::sleep(tokio::time::Duration::from_millis(3)).await;

            let local_record = span.end();
            metrics.local_wal_time_us = (local_record.duration_ms * 1000) as u64;
            local_record.log();

            TracedLogger::debug_exit("local_wal_write", local_record.duration_ms);
            TracedLogger::structured(
                "info",
                "wal_write_complete",
                &[
                    ("entries".to_string(), entry_count.to_string()),
                    ("total_bytes".to_string(), total_size.to_string()),
                    ("time_ms".to_string(), local_record.duration_ms.to_string()),
                ],
            );
        }

        // Stage 2: 원격 노드로 복제 (병렬)
        {
            let span = Span::new(
                trace_id.clone(),
                "replicate_to_followers".to_string(),
                "raft".to_string(),
            ).with_threshold(100);

            TracedLogger::debug_enter("replicate_to_followers");

            // 3개 노드로 병렬 복제
            let node_ids = vec![1, 2, 3];
            let mut node_tasks = Vec::new();

            for node_id in node_ids {
                let trace_id = trace_id.clone();
                let entries_copy = entries.to_vec();

                let task = async move {
                    let node_span = Span::new(
                        trace_id.clone(),
                        format!("replicate_to_node{}", node_id),
                        "raft".to_string(),
                    ).with_threshold(50);

                    let rpc_start = Instant::now();

                    // RPC 요청 로깅
                    NetworkLogger::log_rpc_request("raft", "append_entries", node_id);

                    // RPC 시뮬레이션
                    // 네트워크 지연: 10-40ms
                    let network_delay = 10 + (node_id as u64 * 5);
                    tokio::time::sleep(tokio::time::Duration::from_millis(network_delay)).await;

                    // 원격 서버 처리 시간: 5-10ms
                    tokio::time::sleep(tokio::time::Duration::from_millis(5)).await;

                    let node_record = node_span.end();
                    node_record.log();

                    // RPC 응답 로깅
                    NetworkLogger::log_rpc_response(
                        "raft",
                        "append_entries",
                        node_record.duration_ms,
                        "ok",
                    );

                    (node_id, node_record.duration_ms)
                };

                node_tasks.push(task);
            }

            // 모든 노드에 대한 복제를 병렬로 실행
            let results: Vec<_> = futures::future::join_all(node_tasks).await;

            // 지연 기록
            for (node_id, latency_ms) in results {
                metrics.node_latencies.push((node_id, latency_ms * 1000));
                if latency_ms > 50 {
                    NetworkLogger::log_latency_warning(node_id, latency_ms, 50);
                }
            }

            let replication_record = span.end();
            metrics.total_replication_time_us = (replication_record.duration_ms * 1000) as u64;
            replication_record.log();

            TracedLogger::debug_exit("replicate_to_followers", replication_record.duration_ms);
        }

        // Stage 3: Commit
        {
            let span = Span::new(
                trace_id.clone(),
                "commit".to_string(),
                "raft".to_string(),
            ).with_threshold(20);

            TracedLogger::debug_enter("commit");

            // Commit 시뮬레이션
            tokio::time::sleep(tokio::time::Duration::from_millis(5)).await;

            let commit_record = span.end();
            metrics.commit_time_us = (commit_record.duration_ms * 1000) as u64;
            commit_record.log();

            TracedLogger::debug_exit("commit", commit_record.duration_ms);
        }

        // 최종 리포트
        let parent_record = parent_span.end();
        parent_record.log();

        // 메트릭 저장
        self.metrics.insert(trace_id.clone(), metrics.clone());

        // 병목 분석
        self.analyze_replication(&trace_id, &metrics);
    }

    /// 복제 성능 분석
    fn analyze_replication(&self, trace_id: &str, metrics: &RaftTracingMetrics) {
        let slowest_node = metrics
            .node_latencies
            .iter()
            .max_by_key(|(_, latency)| latency);

        if let Some((node_id, latency_us)) = slowest_node {
            let latency_ms = latency_us / 1000;
            if latency_ms > 50 {
                TracedLogger::warn(&format!(
                    "[TRACE:{}] Slow replication to node {}: {}ms (> 50ms threshold)",
                    trace_id, node_id, latency_ms
                ));
            }
        }

        let total_time_ms = (metrics.total_replication_time_us / 1000) as u64;
        let local_time_ms = (metrics.local_wal_time_us / 1000) as u64;
        let commit_time_ms = (metrics.commit_time_us / 1000) as u64;

        TracedLogger::structured(
            "info",
            "replication_complete",
            &[
                ("total_ms".to_string(), total_time_ms.to_string()),
                ("local_wal_ms".to_string(), local_time_ms.to_string()),
                ("commit_ms".to_string(), commit_time_ms.to_string()),
                ("node_count".to_string(), metrics.node_latencies.len().to_string()),
            ],
        );
    }

    /// 특정 노드의 복제 추적
    pub async fn trace_replicate_to_node(&self, node_id: u32, entries: &[(u32, Vec<u8>)]) {
        let trace_id = get_trace_id();
        let span = Span::new(
            trace_id.clone(),
            format!("replicate_to_node{}", node_id),
            "raft".to_string(),
        ).with_threshold(100);

        let rpc_start = Instant::now();

        // RPC 요청
        NetworkLogger::log_rpc_request("raft", "append_entries", node_id);

        // 네트워크 지연 시뮬레이션
        let network_delay = 10 + (node_id as u64 * 5);
        tokio::time::sleep(tokio::time::Duration::from_millis(network_delay)).await;

        // 원격 처리 시간
        tokio::time::sleep(tokio::time::Duration::from_millis(5)).await;

        let record = span.end();

        // RPC 응답
        NetworkLogger::log_rpc_response(
            "raft",
            "append_entries",
            record.duration_ms,
            "ok",
        );

        record.log();
    }

    /// Leader 선출 추적
    pub async fn trace_leader_election(&self) {
        let trace_id = get_trace_id();
        let span = Span::new(
            trace_id.clone(),
            "leader_election".to_string(),
            "raft".to_string(),
        ).with_threshold(500);

        TracedLogger::debug_enter("leader_election");

        // 단계 1: 타임아웃 감지
        {
            let timeout_span = Span::new(
                trace_id.clone(),
                "election_timeout".to_string(),
                "raft".to_string(),
            );
            tokio::time::sleep(tokio::time::Duration::from_millis(150)).await;
            let record = timeout_span.end();
            record.log();
        }

        // 단계 2: Vote 수집 (병렬)
        {
            let vote_span = Span::new(
                trace_id.clone(),
                "collect_votes".to_string(),
                "raft".to_string(),
            );

            for node in 1..=3 {
                let node_vote_span = Span::new(
                    trace_id.clone(),
                    format!("vote_from_node{}", node),
                    "raft".to_string(),
                );

                let vote_delay = 20 + (node as u64 * 10);
                tokio::time::sleep(tokio::time::Duration::from_millis(vote_delay)).await;

                let record = node_vote_span.end();
                record.log();
            }

            let record = vote_span.end();
            record.log();
        }

        // 단계 3: Leader 선출
        {
            let leader_span = Span::new(
                trace_id.clone(),
                "become_leader".to_string(),
                "raft".to_string(),
            );
            tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
            let record = leader_span.end();
            record.log();
        }

        let parent_record = span.end();
        parent_record.log();
        TracedLogger::debug_exit("leader_election", parent_record.duration_ms);

        TracedLogger::info(&format!(
            "[TRACE:{}] New leader elected ({}ms)",
            trace_id, parent_record.duration_ms
        ));
    }

    /// 메트릭 조회
    pub fn get_metrics(&self, trace_id: &str) -> Option<RaftTracingMetrics> {
        self.metrics.get(trace_id).map(|r| r.clone())
    }

    /// 모든 메트릭 조회
    pub fn get_all_metrics(&self) -> Vec<(String, RaftTracingMetrics)> {
        self.metrics
            .iter()
            .map(|r| (r.key().clone(), r.value().clone()))
            .collect()
    }
}

impl Default for RaftTracer {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tracing::core::{set_trace_context, TracingContext, TraceId};

    #[tokio::test]
    async fn test_raft_tracer_creation() {
        let tracer = RaftTracer::new();
        assert_eq!(tracer.get_all_metrics().len(), 0);
    }

    #[tokio::test]
    async fn test_replicate_log_tracing() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        let tracer = RaftTracer::new();
        let entries = vec![(0u32, vec![1, 2, 3, 4, 5])];

        tracer.trace_replicate_log(&entries).await;

        let trace_id_str = get_trace_id();
        let metrics = tracer.get_metrics(&trace_id_str);
        assert!(metrics.is_some());
    }
}
