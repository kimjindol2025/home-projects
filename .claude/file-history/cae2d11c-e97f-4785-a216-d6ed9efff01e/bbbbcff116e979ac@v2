/// Phase H: Observability - Transaction Tracing
/// 2PC (2-Phase Commit) 트랜잭션 추적

use super::core::{get_trace_id, Span};
use super::logging::TracedLogger;
use std::sync::Arc;
use dashmap::DashMap;

#[derive(Clone, Debug)]
pub struct TransactionTracingMetrics {
    pub tx_id: String,
    pub phase1_prepare_ms: u64,
    pub phase2_commit_ms: u64,
    pub total_ms: u64,
    pub node_times: Vec<(u32, u64)>,
}

pub struct TransactionTracer {
    metrics: Arc<DashMap<String, TransactionTracingMetrics>>,
}

impl TransactionTracer {
    pub fn new() -> Self {
        TransactionTracer {
            metrics: Arc::new(DashMap::new()),
        }
    }

    /// 2PC 트랜잭션 전체 추적
    pub async fn trace_2pc_transaction(&self, tx_id: &str) {
        let trace_id = get_trace_id();
        let parent_span = Span::new(
            trace_id.clone(),
            "2pc_transaction".to_string(),
            "bank".to_string(),
        );

        let mut metrics = TransactionTracingMetrics {
            tx_id: tx_id.to_string(),
            phase1_prepare_ms: 0,
            phase2_commit_ms: 0,
            total_ms: 0,
            node_times: Vec::new(),
        };

        // Phase 1: Prepare (병렬)
        {
            let phase1_span = Span::new(
                trace_id.clone(),
                "phase1_prepare".to_string(),
                "bank".to_string(),
            );

            TracedLogger::debug_enter("phase1_prepare");

            let mut phase1_tasks = Vec::new();

            for node_id in 1..=3 {
                let trace_id = trace_id.clone();
                let tx_id = tx_id.to_string();

                let task = async move {
                    let node_span = Span::new(
                        trace_id.clone(),
                        format!("prepare_node{}", node_id),
                        "bank".to_string(),
                    );

                    // Prepare 시뮬레이션: 5-10ms
                    tokio::time::sleep(
                        tokio::time::Duration::from_millis(5 + (node_id as u64 * 2))
                    ).await;

                    let record = node_span.end();
                    record.log();
                    (node_id, record.duration_ms)
                };

                phase1_tasks.push(task);
            }

            let results: Vec<_> = futures::future::join_all(phase1_tasks).await;
            for (node_id, duration) in results {
                metrics.node_times.push((node_id, duration * 1000));
            }

            let phase1_record = phase1_span.end();
            metrics.phase1_prepare_ms = phase1_record.duration_ms;
            phase1_record.log();
            TracedLogger::debug_exit("phase1_prepare", phase1_record.duration_ms);
        }

        // Phase 2: Commit (병렬)
        {
            let phase2_span = Span::new(
                trace_id.clone(),
                "phase2_commit".to_string(),
                "bank".to_string(),
            );

            TracedLogger::debug_enter("phase2_commit");

            let mut phase2_tasks = Vec::new();

            for node_id in 1..=3 {
                let trace_id = trace_id.clone();
                let tx_id = tx_id.to_string();

                let task = async move {
                    let node_span = Span::new(
                        trace_id.clone(),
                        format!("commit_node{}", node_id),
                        "bank".to_string(),
                    );

                    // Commit 시뮬레이션: 3-8ms
                    tokio::time::sleep(
                        tokio::time::Duration::from_millis(3 + (node_id as u64 * 1))
                    ).await;

                    let record = node_span.end();
                    record.log();
                    (node_id, record.duration_ms)
                };

                phase2_tasks.push(task);
            }

            let results: Vec<_> = futures::future::join_all(phase2_tasks).await;

            let phase2_record = phase2_span.end();
            metrics.phase2_commit_ms = phase2_record.duration_ms;
            phase2_record.log();
            TracedLogger::debug_exit("phase2_commit", phase2_record.duration_ms);
        }

        let parent_record = parent_span.end();
        metrics.total_ms = parent_record.duration_ms;
        parent_record.log();

        self.metrics.insert(trace_id.clone(), metrics.clone());

        TracedLogger::info(&format!(
            "[TRACE:{}] 2PC Transaction {} completed: {}ms (Prepare: {}ms, Commit: {}ms)",
            trace_id,
            tx_id,
            metrics.total_ms,
            metrics.phase1_prepare_ms,
            metrics.phase2_commit_ms
        ));
    }

    pub fn get_metrics(&self, trace_id: &str) -> Option<TransactionTracingMetrics> {
        self.metrics.get(trace_id).map(|r| r.clone())
    }
}

impl Default for TransactionTracer {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tracing::core::{set_trace_context, TracingContext, TraceId};

    #[tokio::test]
    async fn test_2pc_tracing() {
        let trace_id = TraceId::new();
        let ctx = TracingContext::new(trace_id);
        set_trace_context(ctx);

        let tracer = TransactionTracer::new();
        tracer.trace_2pc_transaction("TX001").await;
    }
}
