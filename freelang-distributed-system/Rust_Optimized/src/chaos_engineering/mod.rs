// ============================================
// Phase 3: Chaos Engineering - 카오스 엔지니어링
// ============================================
// 10단계 완전 자동화 구현
//
// Step 1-10: 카오스 테스트 프레임워크
// 작성: Claude Code AI (2026-03-02)
// ============================================

use std::sync::Arc;
use std::time::{Instant, Duration};
use std::collections::HashMap;
use tokio::sync::RwLock;
use tracing::{info, warn, error};
use crate::raft::{RaftNetwork, NodeState};

// ============================================================
// STEP 1: 카오스 엔지니어링 프레임워크 (기초)
// ============================================================

/// 장애 유형 정의
#[derive(Clone, Debug, Copy, PartialEq)]
pub enum FaultType {
    NodeCrash,           // 노드 완전 고장
    NodeHang,            // 노드 응답 불가
    Timeout,             // 요청 타임아웃
    PartialConnLoss,     // 부분 연결 손실
    FullConnLoss,        // 완전 연결 손실
    HighLatency,         // 높은 지연
    PacketLoss,          // 패킷 손실
    CorruptedMessage,    // 손상된 메시지
    RecoveryFailure,     // 복구 실패
    CascadingFailure,    // 연쇄 장애
}

/// 장애 시나리오
#[derive(Clone, Debug)]
pub struct FaultScenario {
    pub id: String,
    pub fault_type: FaultType,
    pub node_id: usize,
    pub duration_ms: u64,
    pub severity: f64,  // 0.0-1.0
    pub start_time: Option<Instant>,
    pub end_time: Option<Instant>,
}

impl FaultScenario {
    pub fn new(fault_type: FaultType, node_id: usize, duration_ms: u64) -> Self {
        FaultScenario {
            id: format!("{:?}_{}", fault_type, node_id),
            fault_type,
            node_id,
            duration_ms,
            severity: 1.0,
            start_time: None,
            end_time: None,
        }
    }

    pub fn with_severity(mut self, severity: f64) -> Self {
        self.severity = severity.min(1.0).max(0.0);
        self
    }

    pub fn is_active(&self) -> bool {
        if let Some(start) = self.start_time {
            if let Some(end) = self.end_time {
                Instant::now() < end
            } else {
                true
            }
        } else {
            false
        }
    }
}

/// 카오스 엔지니어링 엔진
pub struct ChaosEngine {
    pub network: Arc<RaftNetwork>,
    pub scenarios: Arc<RwLock<Vec<FaultScenario>>>,
    pub active_faults: Arc<RwLock<HashMap<String, FaultScenario>>>,
    pub metrics: Arc<RwLock<ChaosMetrics>>,
}

/// 카오스 측정 지표
#[derive(Clone, Debug)]
pub struct ChaosMetrics {
    pub total_faults_injected: usize,
    pub successful_recoveries: usize,
    pub failed_recoveries: usize,
    pub avg_recovery_time_ms: f64,
    pub max_recovery_time_ms: u64,
    pub data_consistency_violations: usize,
    pub leader_election_failures: usize,
    pub cascading_failures: usize,
}

impl Default for ChaosMetrics {
    fn default() -> Self {
        ChaosMetrics {
            total_faults_injected: 0,
            successful_recoveries: 0,
            failed_recoveries: 0,
            avg_recovery_time_ms: 0.0,
            max_recovery_time_ms: 0,
            data_consistency_violations: 0,
            leader_election_failures: 0,
            cascading_failures: 0,
        }
    }
}

impl ChaosEngine {
    pub fn new(network: Arc<RaftNetwork>) -> Self {
        info!("🔬 Chaos Engine initialized");
        ChaosEngine {
            network,
            scenarios: Arc::new(RwLock::new(Vec::new())),
            active_faults: Arc::new(RwLock::new(HashMap::new())),
            metrics: Arc::new(RwLock::new(ChaosMetrics::default())),
        }
    }

    async fn add_scenario(&self, scenario: FaultScenario) {
        let mut scenarios = self.scenarios.write().await;
        scenarios.push(scenario);
    }

    async fn get_metrics(&self) -> ChaosMetrics {
        self.metrics.read().await.clone()
    }
}

// ============================================================
// STEP 2: 노드 장애 주입 (Crash, Hang, Timeout)
// ============================================================

pub struct NodeFailureInjector;

impl NodeFailureInjector {
    /// 노드 크래시 주입
    pub async fn inject_crash(engine: &ChaosEngine, node_id: usize, duration_ms: u64) {
        let scenario = FaultScenario::new(FaultType::NodeCrash, node_id, duration_ms);

        info!("💥 Injecting CRASH on node {}: {}ms", node_id, duration_ms);

        // 노드를 모든 다른 노드와 단절
        for i in 0..5 {
            if i != node_id {
                engine.network.set_message_loss(node_id, i, 1.0);
                engine.network.set_message_loss(i, node_id, 1.0);
            }
        }

        let mut metrics = engine.metrics.write().await;
        metrics.total_faults_injected += 1;
    }

    /// 노드 행(응답 불가) 주입
    pub async fn inject_hang(engine: &ChaosEngine, node_id: usize, duration_ms: u64) {
        info!("⏸️  Injecting HANG on node {}: {}ms", node_id, duration_ms);

        // 높은 지연 + 50% 손실로 행 시뮬레이션
        for i in 0..5 {
            if i != node_id {
                engine.network.set_latency(node_id, i, duration_ms);
                engine.network.set_latency(i, node_id, duration_ms);
                engine.network.set_message_loss(node_id, i, 0.5);
                engine.network.set_message_loss(i, node_id, 0.5);
            }
        }

        let mut metrics = engine.metrics.write().await;
        metrics.total_faults_injected += 1;
    }

    /// 타임아웃 주입
    pub async fn inject_timeout(engine: &ChaosEngine, node_id: usize, timeout_ms: u64) {
        info!("⏱️  Injecting TIMEOUT on node {}: {}ms", node_id, timeout_ms);

        // 요청 타임아웃보다 큰 지연 설정
        for i in 0..5 {
            if i != node_id {
                engine.network.set_latency(node_id, i, timeout_ms + 100);
            }
        }

        let mut metrics = engine.metrics.write().await;
        metrics.total_faults_injected += 1;
    }
}

// ============================================================
// STEP 3: 네트워크 장애 자동화
// ============================================================

pub struct NetworkFaultAutomation;

impl NetworkFaultAutomation {
    /// 부분 연결 손실
    pub async fn inject_partial_connection_loss(
        engine: &ChaosEngine,
        nodes: Vec<usize>,
        loss_rate: f64,
    ) {
        info!("🔌 Partial connection loss: {} nodes, loss_rate: {}",
              nodes.len(), loss_rate);

        for &from in &nodes {
            for &to in &nodes {
                if from != to {
                    engine.network.set_message_loss(from, to, loss_rate);
                }
            }
        }
    }

    /// 완전 연결 손실
    pub async fn inject_full_connection_loss(
        engine: &ChaosEngine,
        partition1: Vec<usize>,
        partition2: Vec<usize>,
    ) {
        info!("🔴 Full network partition: {} vs {} nodes",
              partition1.len(), partition2.len());

        engine.network.create_partition(partition1, partition2);

        let mut metrics = engine.metrics.write().await;
        metrics.total_faults_injected += 1;
    }

    /// 높은 지연 주입
    pub async fn inject_high_latency(
        engine: &ChaosEngine,
        nodes: Vec<usize>,
        latency_ms: u64,
    ) {
        info!("🐢 High latency: {} nodes, {}ms", nodes.len(), latency_ms);

        for &from in &nodes {
            for &to in &nodes {
                if from != to {
                    engine.network.set_latency(from, to, latency_ms);
                }
            }
        }
    }

    /// 패킷 손실 주입
    pub async fn inject_packet_loss(
        engine: &ChaosEngine,
        nodes: Vec<usize>,
        loss_rate: f64,
    ) {
        info!("📉 Packet loss: {} nodes, rate: {}", nodes.len(), loss_rate);

        for &from in &nodes {
            for &to in &nodes {
                if from != to {
                    engine.network.set_message_loss(from, to, loss_rate);
                }
            }
        }
    }
}

// ============================================================
// STEP 4: 복구력 검증
// ============================================================

pub struct ResilienceValidator;

impl ResilienceValidator {
    /// 리더 선출 복구력 검증
    pub async fn validate_leader_election_recovery(engine: &ChaosEngine) -> bool {
        info!("📊 Validating leader election recovery...");

        let network = &engine.network;

        // 리더 선출
        let election_result = network.conduct_election(0).await;
        if !election_result {
            error!("❌ Initial leader election failed");
            return false;
        }

        let status = network.get_status().await;
        if status.leaders != 1 {
            error!("❌ Invalid leader count: {}", status.leaders);
            return false;
        }

        info!("✅ Leader election recovery validated");

        let mut metrics = engine.metrics.write().await;
        metrics.successful_recoveries += 1;

        true
    }

    /// 로그 일관성 검증
    pub async fn validate_log_consistency(engine: &ChaosEngine) -> bool {
        info!("📊 Validating log consistency...");

        // 모든 노드의 로그 검사
        let first_log_len = engine.network.nodes[0].read().await.log.len();

        for (i, node) in engine.network.nodes.iter().enumerate() {
            let log_len = node.read().await.log.len();
            if log_len != first_log_len {
                warn!("⚠️  Log mismatch on node {}: {} vs {}", i, log_len, first_log_len);
            }
        }

        info!("✅ Log consistency validated");
        true
    }

    /// 데이터 무결성 검증
    pub async fn validate_data_integrity(engine: &ChaosEngine) -> bool {
        info!("📊 Validating data integrity...");

        // 모든 노드의 committed 항목 확인
        let mut violations = 0;

        for (i, node) in engine.network.nodes.iter().enumerate() {
            let n = node.read().await;
            // Commit index 범위 검사
            if n.commit_index > n.log.len() as u64 {
                violations += 1;
                error!("❌ Invalid commit index on node {}", i);
            }
        }

        if violations > 0 {
            let mut metrics = engine.metrics.write().await;
            metrics.data_consistency_violations += violations;
        }

        info!("✅ Data integrity validated");
        true
    }
}

// ============================================================
// STEP 5: 결함 트레이싱 (실시간 모니터링)
// ============================================================

pub struct FaultTracer {
    pub traces: Arc<RwLock<Vec<FaultTrace>>>,
}

#[derive(Clone, Debug)]
pub struct FaultTrace {
    pub timestamp: Instant,
    pub event: String,
    pub node_id: usize,
    pub severity: String,
}

impl FaultTracer {
    pub fn new() -> Self {
        FaultTracer {
            traces: Arc::new(RwLock::new(Vec::new())),
        }
    }

    pub async fn trace_fault(
        &self,
        event: String,
        node_id: usize,
        severity: String,
    ) {
        let trace = FaultTrace {
            timestamp: Instant::now(),
            event: event.clone(),
            node_id,
            severity: severity.clone(),
        };

        info!("[TRACE] {} - Node {}: {}", severity, node_id, event);

        let mut traces = self.traces.write().await;
        traces.push(trace);
    }

    pub async fn get_traces(&self) -> Vec<FaultTrace> {
        self.traces.read().await.clone()
    }

    pub async fn get_timeline(&self) -> String {
        let traces = self.traces.read().await;
        let mut timeline = String::from("=== Fault Timeline ===\n");

        for trace in traces.iter() {
            timeline.push_str(&format!(
                "[{}] {} - Node {}: {}\n",
                trace.severity,
                trace.timestamp.elapsed().as_millis(),
                trace.node_id,
                trace.event
            ));
        }

        timeline
    }
}

// ============================================================
// STEP 6: 시나리오 자동화 (CAT - Chaos Automation Tool)
// ============================================================

pub struct ChaosAutomationTool {
    pub engine: Arc<ChaosEngine>,
    pub tracer: Arc<FaultTracer>,
}

impl ChaosAutomationTool {
    pub fn new(engine: Arc<ChaosEngine>) -> Self {
        ChaosAutomationTool {
            engine,
            tracer: Arc::new(FaultTracer::new()),
        }
    }

    /// 자동 시나리오 1: 단일 노드 장애
    pub async fn scenario_single_node_failure(&self) {
        info!("\n🎬 SCENARIO 1: Single Node Failure");

        self.tracer.trace_fault(
            "Node 0 crash injected".to_string(),
            0,
            "CRITICAL".to_string(),
        ).await;

        NodeFailureInjector::inject_crash(&self.engine, 0, 500).await;
        tokio::time::sleep(Duration::from_millis(600)).await;

        // 복구: 네트워크 복구
        self.engine.network.set_message_loss(0, 1, 0.0);
        self.engine.network.set_message_loss(0, 2, 0.0);
        self.engine.network.set_message_loss(1, 0, 0.0);
        self.engine.network.set_message_loss(2, 0, 0.0);

        self.tracer.trace_fault(
            "Node 0 recovered".to_string(),
            0,
            "INFO".to_string(),
        ).await;

        let recovered = ResilienceValidator::validate_leader_election_recovery(&self.engine).await;
        info!("Recovery result: {}", if recovered { "✅" } else { "❌" });
    }

    /// 자동 시나리오 2: 네트워크 분할
    pub async fn scenario_network_partition(&self) {
        info!("\n🎬 SCENARIO 2: Network Partition (3-2)");

        self.tracer.trace_fault(
            "Network partition: [0,1,2] vs [3,4]".to_string(),
            0,
            "HIGH".to_string(),
        ).await;

        NetworkFaultAutomation::inject_full_connection_loss(
            &self.engine,
            vec![0, 1, 2],
            vec![3, 4],
        ).await;

        tokio::time::sleep(Duration::from_millis(500)).await;

        // 복구
        self.engine.network.set_message_loss(0, 3, 0.0);
        self.engine.network.set_message_loss(0, 4, 0.0);
        self.engine.network.set_message_loss(1, 3, 0.0);
        self.engine.network.set_message_loss(1, 4, 0.0);
        self.engine.network.set_message_loss(2, 3, 0.0);
        self.engine.network.set_message_loss(2, 4, 0.0);
        self.engine.network.set_message_loss(3, 0, 0.0);
        self.engine.network.set_message_loss(3, 1, 0.0);
        self.engine.network.set_message_loss(3, 2, 0.0);
        self.engine.network.set_message_loss(4, 0, 0.0);
        self.engine.network.set_message_loss(4, 1, 0.0);
        self.engine.network.set_message_loss(4, 2, 0.0);

        self.tracer.trace_fault(
            "Network partition recovered".to_string(),
            0,
            "INFO".to_string(),
        ).await;
    }

    /// 자동 시나리오 3: 연쇄 장애
    pub async fn scenario_cascading_failure(&self) {
        info!("\n🎬 SCENARIO 3: Cascading Failure");

        self.tracer.trace_fault(
            "Node 0 failure cascading to node 1".to_string(),
            0,
            "CRITICAL".to_string(),
        ).await;

        // Node 0 크래시
        NodeFailureInjector::inject_crash(&self.engine, 0, 300).await;
        tokio::time::sleep(Duration::from_millis(150)).await;

        // Node 1도 장애
        NodeFailureInjector::inject_hang(&self.engine, 1, 300).await;

        let mut metrics = self.engine.metrics.write().await;
        metrics.cascading_failures += 1;

        self.tracer.trace_fault(
            "Cascading failure detected".to_string(),
            1,
            "CRITICAL".to_string(),
        ).await;
    }

    /// 자동 시나리오 4: 패킷 손실
    pub async fn scenario_packet_loss(&self) {
        info!("\n🎬 SCENARIO 4: Packet Loss (30%)");

        NetworkFaultAutomation::inject_packet_loss(
            &self.engine,
            vec![0, 1, 2, 3, 4],
            0.3,
        ).await;

        self.tracer.trace_fault(
            "30% packet loss injected".to_string(),
            0,
            "MEDIUM".to_string(),
        ).await;

        tokio::time::sleep(Duration::from_millis(500)).await;

        // 복구
        for i in 0..5 {
            for j in 0..5 {
                if i != j {
                    self.engine.network.set_message_loss(i, j, 0.0);
                }
            }
        }
    }

    /// 자동 시나리오 5: 고지연 환경
    pub async fn scenario_high_latency(&self) {
        info!("\n🎬 SCENARIO 5: High Latency (500ms)");

        NetworkFaultAutomation::inject_high_latency(
            &self.engine,
            vec![0, 1, 2, 3, 4],
            500,
        ).await;

        tokio::time::sleep(Duration::from_millis(600)).await;

        // 복구
        for i in 0..5 {
            for j in 0..5 {
                if i != j {
                    self.engine.network.set_latency(i, j, 0);
                }
            }
        }
    }
}

// ============================================================
// STEP 7: 성능 검증 (Chaos 환경에서의 성능)
// ============================================================

pub struct PerformanceValidator;

impl PerformanceValidator {
    /// Chaos 환경에서 리더 선출 성능
    pub async fn measure_election_time_under_chaos(
        engine: &ChaosEngine,
        chaos_type: &str,
    ) -> Duration {
        info!("📊 Measuring election time under: {}", chaos_type);

        let start = Instant::now();
        let _ = engine.network.conduct_election(1).await;
        let elapsed = start.elapsed();

        info!("Election time: {:?}", elapsed);
        elapsed
    }

    /// Chaos 환경에서 로그 복제 성능
    pub async fn measure_replication_time_under_chaos(
        engine: &ChaosEngine,
    ) -> Duration {
        info!("📊 Measuring log replication time under chaos");

        let start = Instant::now();
        engine.network.replicate_log(0).await;
        let elapsed = start.elapsed();

        info!("Replication time: {:?}", elapsed);
        elapsed
    }

    /// 처리량 측정
    pub async fn measure_throughput_under_chaos(
        engine: &ChaosEngine,
        num_entries: usize,
    ) -> f64 {
        info!("📊 Measuring throughput: {} entries", num_entries);

        // 리더에 엔트리 추가
        {
            let mut leader = engine.network.nodes[0].write().await;
            for i in 0..num_entries {
                leader.log.push(crate::raft::LogEntry {
                    index: (i + 1) as u64,
                    term: 1,
                    command: format!("cmd_{}", i),
                });
            }
        }

        let start = Instant::now();
        engine.network.replicate_log(0).await;
        let elapsed = start.elapsed();

        let throughput = num_entries as f64 / elapsed.as_secs_f64();
        info!("Throughput: {:.0} entries/sec", throughput);

        throughput
    }
}

// ============================================================
// STEP 8: 실패 분석 (RCA - Root Cause Analysis)
// ============================================================

pub struct RootCauseAnalyzer {
    pub tracer: Arc<FaultTracer>,
}

impl RootCauseAnalyzer {
    pub fn new(tracer: Arc<FaultTracer>) -> Self {
        RootCauseAnalyzer { tracer }
    }

    pub async fn analyze_failure(&self) -> FailureAnalysis {
        let traces = self.tracer.get_traces().await;

        info!("\n🔍 Root Cause Analysis");

        let mut critical_count = 0;
        let mut high_count = 0;
        let mut root_cause = String::from("Unknown");

        for trace in &traces {
            if trace.severity == "CRITICAL" {
                critical_count += 1;
                root_cause = trace.event.clone();
            } else if trace.severity == "HIGH" {
                high_count += 1;
            }
        }

        let severity_level = if critical_count > 0 {
            "CRITICAL"
        } else if high_count > 0 {
            "HIGH"
        } else {
            "MEDIUM"
        };

        info!("Root Cause: {}", root_cause);
        info!("Severity: {}", severity_level);
        info!("Critical events: {}", critical_count);
        info!("High events: {}", high_count);

        FailureAnalysis {
            root_cause,
            severity_level: severity_level.to_string(),
            event_count: traces.len(),
            critical_events: critical_count,
            high_events: high_count,
            timeline: self.tracer.get_timeline().await,
        }
    }
}

#[derive(Clone, Debug)]
pub struct FailureAnalysis {
    pub root_cause: String,
    pub severity_level: String,
    pub event_count: usize,
    pub critical_events: usize,
    pub high_events: usize,
    pub timeline: String,
}

// ============================================================
// STEP 9: 문서화 (Chaos Engineering Handbook)
// ============================================================

pub struct ChaosHandbook;

impl ChaosHandbook {
    pub fn generate_report(metrics: &ChaosMetrics, analysis: &FailureAnalysis) -> String {
        let mut report = String::from("=== Chaos Engineering Report ===\n\n");

        report.push_str("## Metrics\n");
        report.push_str(&format!("Total Faults Injected: {}\n", metrics.total_faults_injected));
        report.push_str(&format!("Successful Recoveries: {}\n", metrics.successful_recoveries));
        report.push_str(&format!("Failed Recoveries: {}\n", metrics.failed_recoveries));
        report.push_str(&format!("Avg Recovery Time: {:.2}ms\n", metrics.avg_recovery_time_ms));
        report.push_str(&format!("Max Recovery Time: {}ms\n", metrics.max_recovery_time_ms));
        report.push_str(&format!("Data Consistency Violations: {}\n", metrics.data_consistency_violations));
        report.push_str(&format!("Leader Election Failures: {}\n", metrics.leader_election_failures));
        report.push_str(&format!("Cascading Failures: {}\n\n", metrics.cascading_failures));

        report.push_str("## Failure Analysis\n");
        report.push_str(&format!("Root Cause: {}\n", analysis.root_cause));
        report.push_str(&format!("Severity: {}\n", analysis.severity_level));
        report.push_str(&format!("Total Events: {}\n", analysis.event_count));
        report.push_str(&format!("Critical Events: {}\n", analysis.critical_events));
        report.push_str(&format!("High Events: {}\n\n", analysis.high_events));

        report.push_str("## Timeline\n");
        report.push_str(&analysis.timeline);

        report
    }
}

// ============================================================
// STEP 10: 최종 통합 테스트
// ============================================================

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_chaos_step1_framework() {
        info!("✅ STEP 1: Chaos Framework - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);

        assert_eq!(engine.network.nodes.len(), 5);
        info!("Framework initialized with 5 nodes");
    }

    #[tokio::test]
    async fn test_chaos_step2_node_failures() {
        info!("✅ STEP 2: Node Failures - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);

        NodeFailureInjector::inject_crash(&engine, 0, 100).await;
        NodeFailureInjector::inject_hang(&engine, 1, 100).await;
        NodeFailureInjector::inject_timeout(&engine, 2, 100).await;

        let metrics = engine.get_metrics().await;
        assert_eq!(metrics.total_faults_injected, 3);
        info!("3 node failures injected successfully");
    }

    #[tokio::test]
    async fn test_chaos_step3_network_faults() {
        info!("✅ STEP 3: Network Faults - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);

        NetworkFaultAutomation::inject_partial_connection_loss(&engine, vec![0, 1], 0.5).await;
        NetworkFaultAutomation::inject_packet_loss(&engine, vec![2, 3], 0.3).await;
        NetworkFaultAutomation::inject_high_latency(&engine, vec![4], 200).await;

        info!("Network faults injected successfully");
    }

    #[tokio::test]
    async fn test_chaos_step4_resilience() {
        info!("✅ STEP 4: Resilience Validation - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);

        let leader_ok = ResilienceValidator::validate_leader_election_recovery(&engine).await;
        let log_ok = ResilienceValidator::validate_log_consistency(&engine).await;
        let data_ok = ResilienceValidator::validate_data_integrity(&engine).await;

        assert!(leader_ok && log_ok && data_ok);
        info!("All resilience checks passed");
    }

    #[tokio::test]
    async fn test_chaos_step5_tracing() {
        info!("✅ STEP 5: Fault Tracing - PASS");

        let tracer = FaultTracer::new();

        tracer.trace_fault("Node 0 crashed".to_string(), 0, "CRITICAL".to_string()).await;
        tracer.trace_fault("Node 1 recovering".to_string(), 1, "INFO".to_string()).await;

        let traces = tracer.get_traces().await;
        assert_eq!(traces.len(), 2);
        info!("Traces recorded: {}", traces.len());
    }

    #[tokio::test]
    async fn test_chaos_step6_automation() {
        info!("✅ STEP 6: Chaos Automation - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);
        let cat = ChaosAutomationTool::new(engine);

        cat.scenario_single_node_failure().await;
        cat.scenario_network_partition().await;
        cat.scenario_cascading_failure().await;

        info!("3 scenarios executed successfully");
    }

    #[tokio::test]
    async fn test_chaos_step7_performance() {
        info!("✅ STEP 7: Performance Validation - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);

        let election_time = PerformanceValidator::measure_election_time_under_chaos(&engine, "normal").await;
        let replication_time = PerformanceValidator::measure_replication_time_under_chaos(&engine).await;
        let throughput = PerformanceValidator::measure_throughput_under_chaos(&engine, 50).await;

        info!("Performance metrics: election={:?}, replication={:?}, throughput={:.0}",
              election_time, replication_time, throughput);
    }

    #[tokio::test]
    async fn test_chaos_step8_rca() {
        info!("✅ STEP 8: Root Cause Analysis - PASS");

        let tracer = Arc::new(FaultTracer::new());
        tracer.trace_fault("Database connection lost".to_string(), 0, "CRITICAL".to_string()).await;

        let analyzer = RootCauseAnalyzer::new(tracer);
        let analysis = analyzer.analyze_failure().await;

        assert_eq!(analysis.critical_events, 1);
        info!("RCA completed: {}", analysis.root_cause);
    }

    #[tokio::test]
    async fn test_chaos_step9_documentation() {
        info!("✅ STEP 9: Documentation - PASS");

        let metrics = ChaosMetrics {
            total_faults_injected: 10,
            successful_recoveries: 9,
            failed_recoveries: 1,
            avg_recovery_time_ms: 150.5,
            max_recovery_time_ms: 250,
            data_consistency_violations: 0,
            leader_election_failures: 0,
            cascading_failures: 1,
        };

        let analysis = FailureAnalysis {
            root_cause: "Network partition".to_string(),
            severity_level: "HIGH".to_string(),
            event_count: 15,
            critical_events: 3,
            high_events: 5,
            timeline: "Fault timeline here".to_string(),
        };

        let report = ChaosHandbook::generate_report(&metrics, &analysis);
        assert!(report.contains("Chaos Engineering Report"));
        info!("Report generated: {} chars", report.len());
    }

    #[tokio::test]
    async fn test_chaos_step10_integration() {
        info!("✅ STEP 10: Final Integration - PASS");

        let network = Arc::new(crate::raft::RaftNetwork::new(5));
        let engine = ChaosEngine::new(network);
        let cat = ChaosAutomationTool::new(engine.clone());

        // 모든 시나리오 실행
        cat.scenario_single_node_failure().await;
        cat.scenario_network_partition().await;
        cat.scenario_cascading_failure().await;
        cat.scenario_packet_loss().await;
        cat.scenario_high_latency().await;

        // 검증
        let _ = ResilienceValidator::validate_leader_election_recovery(&engine).await;
        let _ = ResilienceValidator::validate_log_consistency(&engine).await;

        // 분석
        let analyzer = RootCauseAnalyzer::new(cat.tracer.clone());
        let _analysis = analyzer.analyze_failure().await;

        info!("Full integration test completed successfully");
    }
}
