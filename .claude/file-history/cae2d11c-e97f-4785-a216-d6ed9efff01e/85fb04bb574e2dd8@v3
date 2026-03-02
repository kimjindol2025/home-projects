/// Phase I: Chaos Engineering - Fault Injection Engine
/// 네트워크 지연, 단절, 노드 크래시 등을 주입하여 시스템 복원력 테스트

use std::sync::Arc;
use tracing::{info, warn, error};
use std::time::Duration;
use tracing::{info, warn, error};
use dashmap::DashMap;
use tracing::{info, warn, error};
use tokio::time::sleep;
use tracing::{info, warn, error};

/// 장애 유형
#[derive(Clone, Debug, Copy, PartialEq, Eq)]
pub enum FailureType {
    NetworkLatency,     // 네트워크 지연
    NetworkPartition,   // 네트워크 단절
    NodeCrash,         // 노드 다운
    DiskError,         // 디스크 오류
    SlowReplication,   // 느린 복제
    DataCorruption,    // 데이터 손상
    MultipleFailures,  // 다중 장애
}

/// 장애 대상
#[derive(Clone, Debug)]
pub struct FailureTarget {
    pub failure_type: FailureType,
    pub target_node: u32,
    pub duration_ms: u64,
    pub severity: f64, // 0.0 ~ 1.0 (0%~100%)
}

/// 네트워크 지연 주입기
pub struct NetworkDelayInjector {
    delays: Arc<DashMap<u32, u64>>, // node_id -> delay_ms
}

impl NetworkDelayInjector {
    pub fn new() -> Self {
        NetworkDelayInjector {
            delays: Arc::new(DashMap::new()),
        }
    }

    /// 특정 노드에 지연 추가
    pub async fn inject_latency(&self, node: u32, delay_ms: u64) {
        self.delays.insert(node, delay_ms);
        info!("[CHAOS] Injected {}ms latency to node {}", delay_ms, node);
    }

    /// 지연 해제
    pub async fn remove_latency(&self, node: u32) {
        self.delays.remove(&node);
        info!("[CHAOS] Removed latency from node {}", node);
    }

    /// 특정 노드의 현재 지연 조회
    pub fn get_latency(&self, node: u32) -> u64 {
        self.delays
            .get(&node)
            .map(|v| *v)
            .unwrap_or(0)
    }

    /// 지연 시뮬레이션
    pub async fn apply_delay(&self, node: u32) {
        if let Some(delay) = self.get_latency(node) {
            if delay > 0 {
                sleep(Duration::from_millis(delay)).await;
            }
        }
    }
}

impl Default for NetworkDelayInjector {
    fn default() -> Self {
        Self::new()
    }
}

/// 네트워크 단절 주입기
pub struct NetworkPartitionInjector {
    partitions: Arc<DashMap<u32, bool>>, // node_id -> is_partitioned
}

impl NetworkPartitionInjector {
    pub fn new() -> Self {
        NetworkPartitionInjector {
            partitions: Arc::new(DashMap::new()),
        }
    }

    /// 특정 노드 격리
    pub async fn partition_node(&self, node: u32) {
        self.partitions.insert(node, true);
        info!("[CHAOS] Network partition on node {} (ISOLATED)", node);
    }

    /// 노드 복구
    pub async fn heal_partition(&self, node: u32) {
        self.partitions.insert(node, false);
        info!("[CHAOS] Network partition healed on node {}", node);
    }

    /// 도달 가능성 확인
    pub fn is_reachable(&self, node: u32) -> bool {
        !self.partitions
            .get(&node)
            .map(|v| *v)
            .unwrap_or(false)
    }

    /// 모든 격리된 노드 조회
    pub fn get_partitioned_nodes(&self) -> Vec<u32> {
        self.partitions
            .iter()
            .filter(|r| *r.value())
            .map(|r| *r.key())
            .collect()
    }
}

impl Default for NetworkPartitionInjector {
    fn default() -> Self {
        Self::new()
    }
}

/// 노드 크래시 시뮬레이터
pub struct NodeCrashSimulator {
    crashed_nodes: Arc<DashMap<u32, bool>>, // node_id -> is_crashed
}

impl NodeCrashSimulator {
    pub fn new() -> Self {
        NodeCrashSimulator {
            crashed_nodes: Arc::new(DashMap::new()),
        }
    }

    /// 노드 크래시
    pub async fn crash_node(&self, node: u32) {
        self.crashed_nodes.insert(node, true);
        info!("[CHAOS] Node {} CRASHED (all operations failed)", node);
    }

    /// 노드 복구
    pub async fn restart_node(&self, node: u32) {
        self.crashed_nodes.insert(node, false);
        info!("[CHAOS] Node {} restarted (recovering state)", node);
    }

    /// 노드 실행 상태 확인
    pub fn is_running(&self, node: u32) -> bool {
        !self.crashed_nodes
            .get(&node)
            .map(|v| *v)
            .unwrap_or(false)
    }

    /// 크래시된 노드 목록
    pub fn get_crashed_nodes(&self) -> Vec<u32> {
        self.crashed_nodes
            .iter()
            .filter(|r| *r.value())
            .map(|r| *r.key())
            .collect()
    }
}

impl Default for NodeCrashSimulator {
    fn default() -> Self {
        Self::new()
    }
}

/// 디스크 오류 시뮬레이터
pub struct DiskErrorInjector {
    error_rate: Arc<parking_lot::Mutex<f64>>, // 0.0 ~ 1.0
}

impl DiskErrorInjector {
    pub fn new() -> Self {
        DiskErrorInjector {
            error_rate: Arc::new(parking_lot::Mutex::new(0.0)),
        }
    }

    /// 에러율 설정 (0.0 ~ 1.0, 예: 0.1 = 10%)
    pub async fn set_error_rate(&self, rate: f64) {
        let rate = rate.clamp(0.0, 1.0);
        *self.error_rate.lock() = rate;
        info!("[CHAOS] Disk error rate set to {:.1}%", rate * 100.0);
    }

    /// 오류 발생 시뮬레이션
    pub fn should_fail(&self) -> bool {
        let rate = *self.error_rate.lock();
        if rate <= 0.0 {
            return false;
        }
        rand::random::<f64>() < rate
    }

    /// 현재 에러율
    pub fn get_error_rate(&self) -> f64 {
        *self.error_rate.lock()
    }
}

impl Default for DiskErrorInjector {
    fn default() -> Self {
        Self::new()
    }
}

/// 복제 속도 제한기
pub struct ReplicationThrottler {
    bytes_per_second: Arc<parking_lot::Mutex<u64>>,
}

impl ReplicationThrottler {
    pub fn new() -> Self {
        ReplicationThrottler {
            bytes_per_second: Arc::new(parking_lot::Mutex::new(u64::MAX)),
        }
    }

    /// 복제 속도 제한 (bytes/sec)
    pub async fn set_speed_limit(&self, bps: u64) {
        *self.bytes_per_second.lock() = bps;
        info!("[CHAOS] Replication speed limited to {} bytes/sec", bps);
    }

    /// 속도 제한 제거
    pub async fn remove_limit(&self) {
        *self.bytes_per_second.lock() = u64::MAX;
    }

    /// 주어진 데이터 전송 시간 계산
    pub async fn calculate_transfer_time(&self, bytes: u64) -> Duration {
        let bps = *self.bytes_per_second.lock();
        if bps == u64::MAX {
            Duration::from_millis(0)
        } else {
            let ms = (bytes as f64 / bps as f64) * 1000.0;
            Duration::from_millis(ms as u64)
        }
    }
}

impl Default for ReplicationThrottler {
    fn default() -> Self {
        Self::new()
    }
}

/// 통합 카오스 오케스트레이터
pub struct ChaosOrchestrator {
    pub network_delay: Arc<NetworkDelayInjector>,
    pub partition: Arc<NetworkPartitionInjector>,
    pub crash: Arc<NodeCrashSimulator>,
    pub disk_error: Arc<DiskErrorInjector>,
    pub replication: Arc<ReplicationThrottler>,
}

impl ChaosOrchestrator {
    pub fn new() -> Self {
        ChaosOrchestrator {
            network_delay: Arc::new(NetworkDelayInjector::new()),
            partition: Arc::new(NetworkPartitionInjector::new()),
            crash: Arc::new(NodeCrashSimulator::new()),
            disk_error: Arc::new(DiskErrorInjector::new()),
            replication: Arc::new(ReplicationThrottler::new()),
        }
    }

    /// 카오스 시나리오 시작
    pub async fn inject_chaos(&self, scenario: FailureType, node: u32, duration_ms: u64, severity: f64) {
        info!(
            "\n🔴 [CHAOS] Injecting {:?} on node {} for {}ms (severity: {:.1}%)\n",
            scenario, node, duration_ms, severity * 100.0
        );

        match scenario {
            FailureType::NetworkLatency => {
                let delay = (duration_ms as f64 * severity) as u64;
                self.network_delay.inject_latency(node, delay).await;
            }
            FailureType::NetworkPartition => {
                self.partition.partition_node(node).await;
            }
            FailureType::NodeCrash => {
                self.crash.crash_node(node).await;
            }
            FailureType::DiskError => {
                self.disk_error.set_error_rate(severity).await;
            }
            FailureType::SlowReplication => {
                let limit = (10_000_000.0 * severity) as u64;
                self.replication.set_speed_limit(limit).await;
            }
            FailureType::DataCorruption => {
                info!("[CHAOS] Data corruption rate: {:.1}%", severity * 100.0);
            }
            FailureType::MultipleFailures => {
                info!("[CHAOS] Multiple concurrent failures injected");
            }
        }

        // 지정된 시간 동안 장애 유지
        sleep(Duration::from_millis(duration_ms)).await;

        // 복구
        self.recover_from_chaos(scenario, node).await;
    }

    /// 카오스 복구
    pub async fn recover_from_chaos(&self, scenario: FailureType, node: u32) {
        info!("[CHAOS] Recovering from {:?}...", scenario);

        match scenario {
            FailureType::NetworkLatency => {
                self.network_delay.remove_latency(node).await;
            }
            FailureType::NetworkPartition => {
                self.partition.heal_partition(node).await;
            }
            FailureType::NodeCrash => {
                self.crash.restart_node(node).await;
            }
            FailureType::DiskError => {
                self.disk_error.set_error_rate(0.0).await;
            }
            FailureType::SlowReplication => {
                self.replication.remove_limit().await;
            }
            _ => {}
        }

        info!("[CHAOS] ✅ Recovery complete\n");
    }
}

impl Default for ChaosOrchestrator {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_network_delay_injection() {
        let injector = NetworkDelayInjector::new();
        injector.inject_latency(1, 100).await;
        assert_eq!(injector.get_latency(1), 100);
    }

    #[tokio::test]
    async fn test_network_partition() {
        let partitioner = NetworkPartitionInjector::new();
        partitioner.partition_node(2).await;
        assert!(!partitioner.is_reachable(2));

        partitioner.heal_partition(2).await;
        assert!(partitioner.is_reachable(2));
    }

    #[tokio::test]
    async fn test_node_crash() {
        let simulator = NodeCrashSimulator::new();
        simulator.crash_node(3).await;
        assert!(!simulator.is_running(3));

        simulator.restart_node(3).await;
        assert!(simulator.is_running(3));
    }

    #[tokio::test]
    async fn test_disk_error_injection() {
        let injector = DiskErrorInjector::new();
        injector.set_error_rate(0.5).await;
        assert_eq!(injector.get_error_rate(), 0.5);
    }

    #[tokio::test]
    async fn test_chaos_orchestrator() {
        let orchestrator = ChaosOrchestrator::new();
        assert_eq!(orchestrator.crash.get_crashed_nodes().len(), 0);
    }
}
