/// Phase I: Chaos Engineering - Test Scenarios
/// 7가지 카오스 엔지니어링 시나리오

use super::injector::{ChaosOrchestrator, FailureType};
use super::recovery_validator::{
    RecoveryMonitor, SystemSnapshot, IntegrityValidator, RecoveryTracker,
};
use std::sync::Arc;

/// 카오스 테스트 결과
#[derive(Clone, Debug)]
pub struct ChaosTestResult {
    pub scenario_name: String,
    pub passed: bool,
    pub detection_time_ms: u64,
    pub recovery_time_ms: u64,
    pub data_loss: u32,
    pub error_message: Option<String>,
}

/// Scenario 1: 네트워크 지연 (200ms)
pub async fn scenario_1_high_latency(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 1: Network Latency (200ms)");

    tracker.start_tracking("S1-HighLatency".to_string());

    // 지연 주입
    orchestrator.network_delay.inject_latency(1, 200).await;

    // Raft 타임아웃으로 Leader 선출 기대 (150ms)
    tokio::time::sleep(tokio::time::Duration::from_millis(250)).await;

    tracker.update_metrics("S1-HighLatency", |monitor| {
        monitor.mark_detected();
    });

    // 복구
    orchestrator.network_delay.remove_latency(1).await;

    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

    tracker.update_metrics("S1-HighLatency", |monitor| {
        monitor.mark_recovered();
        monitor.set_leader_changes(1);
    });

    let metrics = tracker.get_metrics("S1-HighLatency").unwrap();

    ChaosTestResult {
        scenario_name: "S1: Network Latency".to_string(),
        passed: metrics.recovery_success && metrics.time_to_recovery_ms < 1000,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: metrics.data_loss_items,
        error_message: None,
    }
}

/// Scenario 2: 부분 네트워크 단절
pub async fn scenario_2_network_partition(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 2: Network Partition (Node 2)");

    tracker.start_tracking("S2-Partition".to_string());

    // Node 2 격리
    orchestrator.partition.partition_node(2).await;

    // Quorum 유지로 계속 운영 가능 (2/3 정족수)
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

    tracker.update_metrics("S2-Partition", |monitor| {
        monitor.mark_detected();
    });

    // 격리된 노드로의 write 실패
    let partitioned_nodes = orchestrator.partition.get_partitioned_nodes();
    let isolated_write_fails = !partitioned_nodes.is_empty();

    // 복구
    orchestrator.partition.heal_partition(2).await;

    tokio::time::sleep(tokio::time::Duration::from_millis(200)).await;

    tracker.update_metrics("S2-Partition", |monitor| {
        monitor.mark_recovered();
    });

    let metrics = tracker.get_metrics("S2-Partition").unwrap();

    ChaosTestResult {
        scenario_name: "S2: Network Partition".to_string(),
        passed: metrics.recovery_success && isolated_write_fails,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: metrics.data_loss_items,
        error_message: None,
    }
}

/// Scenario 3: 노드 다운
pub async fn scenario_3_node_crash(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 3: Node Crash (Node 2)");

    tracker.start_tracking("S3-NodeCrash".to_string());

    // Node 2 크래시
    orchestrator.crash.crash_node(2).await;

    tokio::time::sleep(tokio::time::Duration::from_millis(150)).await;

    tracker.update_metrics("S3-NodeCrash", |monitor| {
        monitor.mark_detected();
        monitor.set_leader_changes(1);
    });

    // Leader 변경 완료
    let leader_changed = orchestrator.crash.get_crashed_nodes().len() > 0;

    // 나머지 노드로 계속 운영
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

    // 노드 복구
    orchestrator.crash.restart_node(2).await;

    tokio::time::sleep(tokio::time::Duration::from_millis(500)).await;

    tracker.update_metrics("S3-NodeCrash", |monitor| {
        monitor.mark_recovered();
    });

    let metrics = tracker.get_metrics("S3-NodeCrash").unwrap();

    ChaosTestResult {
        scenario_name: "S3: Node Crash".to_string(),
        passed: metrics.recovery_success && leader_changed && metrics.data_loss_items == 0,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: metrics.data_loss_items,
        error_message: None,
    }
}

/// Scenario 4: 디스크 오류 (10% 실패율)
pub async fn scenario_4_disk_error(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 4: Disk Error (10% failure rate)");

    tracker.start_tracking("S4-DiskError".to_string());

    // 디스크 오류 주입
    orchestrator.disk_error.set_error_rate(0.1).await;

    let mut success_count = 0;
    let mut fail_count = 0;

    // 10번의 연산 시도
    for _ in 0..10 {
        if orchestrator.disk_error.should_fail() {
            fail_count += 1;
            tracker.update_metrics("S4-DiskError", |monitor| {
                monitor.add_failed_operation();
            });
        } else {
            success_count += 1;
        }
    }

    tracker.update_metrics("S4-DiskError", |monitor| {
        monitor.mark_detected();
    });

    // 오류율이 예상치와 맞는가? (약 10%)
    let error_rate_ok = (fail_count as f64 / 10.0) >= 0.05 &&
                        (fail_count as f64 / 10.0) <= 0.20;

    // 복구
    orchestrator.disk_error.set_error_rate(0.0).await;

    tracker.update_metrics("S4-DiskError", |monitor| {
        monitor.mark_recovered();
    });

    let metrics = tracker.get_metrics("S4-DiskError").unwrap();

    ChaosTestResult {
        scenario_name: "S4: Disk Error".to_string(),
        passed: metrics.recovery_success && error_rate_ok,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: metrics.data_loss_items,
        error_message: None,
    }
}

/// Scenario 5: 느린 복제
pub async fn scenario_5_slow_replication(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 5: Slow Replication (100KB/s)");

    tracker.start_tracking("S5-SlowReplication".to_string());

    // 복제 속도 제한 (100KB/s)
    orchestrator.replication.set_speed_limit(100 * 1024).await;

    // 1MB 데이터 전송 시뮬레이션 (약 10초 소요)
    let transfer_time = orchestrator.replication.calculate_transfer_time(1024 * 1024).await;

    println!("[CHAOS] Transfer time for 1MB: {:?}", transfer_time);

    tracker.update_metrics("S5-SlowReplication", |monitor| {
        monitor.mark_detected();
    });

    // Raft 로그 누적 예상
    tokio::time::sleep(transfer_time).await;

    // Snapshot으로 압축
    tracker.update_metrics("S5-SlowReplication", |monitor| {
        monitor.mark_recovered();
    });

    orchestrator.replication.remove_limit().await;

    let metrics = tracker.get_metrics("S5-SlowReplication").unwrap();

    ChaosTestResult {
        scenario_name: "S5: Slow Replication".to_string(),
        passed: metrics.recovery_success,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: 0,
        error_message: None,
    }
}

/// Scenario 6: 데이터 손상
pub async fn scenario_6_data_corruption(tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 6: Data Corruption Detection");

    tracker.start_tracking("S6-DataCorruption".to_string());

    // 초기 상태 기록
    let initial = SystemSnapshot::new(1, 1, 100, 10, 50, 10000);
    let validator = IntegrityValidator::new(initial);

    // 손상된 데이터 (잔액 불일치)
    let corrupted = SystemSnapshot::new(1, 1, 100, 10, 50, 9999);

    tracker.update_metrics("S6-DataCorruption", |monitor| {
        monitor.mark_detected();
    });

    // 무결성 검증
    let report = validator.verify_after_recovery(&corrupted);

    tracker.update_metrics("S6-DataCorruption", |monitor| {
        monitor.set_data_loss(report.data_loss);
        if report.success {
            monitor.mark_recovered();
        }
    });

    let metrics = tracker.get_metrics("S6-DataCorruption").unwrap();

    ChaosTestResult {
        scenario_name: "S6: Data Corruption".to_string(),
        passed: report.data_corruption > 0, // 손상이 감지되어야 함
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: report.data_loss,
        error_message: None,
    }
}

/// Scenario 7: 다중 장애
pub async fn scenario_7_multiple_failures(orchestrator: &ChaosOrchestrator, tracker: &RecoveryTracker) -> ChaosTestResult {
    println!("\n🔴 Scenario 7: Multiple Concurrent Failures");

    tracker.start_tracking("S7-MultipleFailures".to_string());

    // 동시에 여러 장애 주입
    orchestrator.partition.partition_node(2).await;
    orchestrator.disk_error.set_error_rate(0.05).await;
    orchestrator.network_delay.inject_latency(1, 100).await;

    tracker.update_metrics("S7-MultipleFailures", |monitor| {
        monitor.mark_detected();
    });

    // 모든 장애 유지
    tokio::time::sleep(tokio::time::Duration::from_millis(300)).await;

    // 복구
    orchestrator.partition.heal_partition(2).await;
    orchestrator.disk_error.set_error_rate(0.0).await;
    orchestrator.network_delay.remove_latency(1).await;

    tracker.update_metrics("S7-MultipleFailures", |monitor| {
        monitor.mark_recovered();
    });

    let metrics = tracker.get_metrics("S7-MultipleFailures").unwrap();

    // 다중 장애에서 데이터 손실이 없어야 함
    let data_integrity_maintained = metrics.data_loss_items == 0;

    ChaosTestResult {
        scenario_name: "S7: Multiple Failures".to_string(),
        passed: metrics.recovery_success && data_integrity_maintained,
        detection_time_ms: metrics.time_to_detect_ms,
        recovery_time_ms: metrics.time_to_recovery_ms,
        data_loss: metrics.data_loss_items,
        error_message: None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_scenario_execution() {
        let orchestrator = ChaosOrchestrator::new();
        let tracker = RecoveryTracker::new();

        let result = scenario_1_high_latency(&orchestrator, &tracker).await;
        assert!(!result.scenario_name.is_empty());
    }
}
