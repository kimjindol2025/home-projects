/// Phase I: Chaos Engineering - Recovery Validator
/// 장애 후 시스템 복구 상태 검증

use std::time::Instant;
use tracing::{info, warn, error};
use dashmap::DashMap;
use tracing::{info, warn, error};
use std::sync::Arc;
use tracing::{info, warn, error};

/// 시스템 상태 스냅샷
#[derive(Clone, Debug)]
pub struct SystemSnapshot {
    pub timestamp: u64,
    pub leader_id: u32,
    pub raft_term: u64,
    pub log_index: u64,
    pub account_count: u32,
    pub transaction_count: u32,
    pub total_balance: u64,
}

impl SystemSnapshot {
    pub fn new(
        leader_id: u32,
        raft_term: u64,
        log_index: u64,
        account_count: u32,
        transaction_count: u32,
        total_balance: u64,
    ) -> Self {
        SystemSnapshot {
            timestamp: std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .unwrap()
                .as_millis() as u64,
            leader_id,
            raft_term,
            log_index,
            account_count,
            transaction_count,
            total_balance,
        }
    }
}

/// 복구 메트릭
#[derive(Clone, Debug, Default)]
pub struct RecoveryMetrics {
    pub time_to_detect_ms: u64,      // 문제 감지까지의 시간
    pub time_to_recovery_ms: u64,    // 완전 복구까지의 시간
    pub leader_change_count: u32,    // Leader 변경 횟수
    pub data_loss_items: u32,        // 손실된 항목 수
    pub operations_failed: u32,      // 실패한 연산 수
    pub operations_retried: u32,     // 재시도된 연산 수
    pub recovery_success: bool,      // 복구 성공 여부
}

/// 복구 모니터
pub struct RecoveryMonitor {
    failure_start: Instant,
    detection_time: Option<Instant>,
    recovery_complete: Option<Instant>,
    metrics: RecoveryMetrics,
}

impl RecoveryMonitor {
    pub fn new() -> Self {
        RecoveryMonitor {
            failure_start: Instant::now(),
            detection_time: None,
            recovery_complete: None,
            metrics: RecoveryMetrics::default(),
        }
    }

    /// 문제 감지 표시
    pub fn mark_detected(&mut self) {
        self.detection_time = Some(Instant::now());
        if let Some(detect_time) = self.detection_time {
            let duration_ms = detect_time.duration_since(self.failure_start).as_millis() as u64;
            self.metrics.time_to_detect_ms = duration_ms;
            info!("[CHAOS] Problem detected in {}ms", duration_ms);
        }
    }

    /// 복구 완료 표시
    pub fn mark_recovered(&mut self) {
        self.recovery_complete = Some(Instant::now());
        if let Some(recover_time) = self.recovery_complete {
            let duration_ms = recover_time.duration_since(self.failure_start).as_millis() as u64;
            self.metrics.time_to_recovery_ms = duration_ms;
            self.metrics.recovery_success = true;
            info!("[CHAOS] System recovered in {}ms", duration_ms);
        }
    }

    /// 메트릭 업데이트
    pub fn add_failed_operation(&mut self) {
        self.metrics.operations_failed += 1;
    }

    pub fn add_retried_operation(&mut self) {
        self.metrics.operations_retried += 1;
    }

    pub fn set_data_loss(&mut self, count: u32) {
        self.metrics.data_loss_items = count;
    }

    pub fn set_leader_changes(&mut self, count: u32) {
        self.metrics.leader_change_count = count;
    }

    /// 현재 메트릭 조회
    pub fn metrics(&self) -> RecoveryMetrics {
        self.metrics.clone()
    }
}

impl Default for RecoveryMonitor {
    fn default() -> Self {
        Self::new()
    }
}

/// 무결성 검증기
pub struct IntegrityValidator {
    initial_snapshot: SystemSnapshot,
}

impl IntegrityValidator {
    pub fn new(snapshot: SystemSnapshot) -> Self {
        IntegrityValidator {
            initial_snapshot: snapshot,
        }
    }

    /// 복구 후 무결성 검증
    pub fn verify_after_recovery(&self, final_snapshot: &SystemSnapshot) -> IntegrityReport {
        let mut report = IntegrityReport {
            data_loss: 0,
            data_corruption: 0,
            inconsistencies: Vec::new(),
            success: true,
        };

        // 1. 계좌 수 검증
        if final_snapshot.account_count != self.initial_snapshot.account_count {
            report.data_loss += 1;
            report.inconsistencies.push(format!(
                "Account count mismatch: {} vs {}",
                self.initial_snapshot.account_count,
                final_snapshot.account_count
            ));
        }

        // 2. 트랜잭션 수 검증
        if final_snapshot.transaction_count < self.initial_snapshot.transaction_count {
            let lost = self.initial_snapshot.transaction_count - final_snapshot.transaction_count;
            report.data_loss += lost;
            report.inconsistencies.push(format!(
                "Lost {} transactions",
                lost
            ));
        }

        // 3. 잔액 검증
        if final_snapshot.total_balance != self.initial_snapshot.total_balance {
            report.data_corruption += 1;
            report.inconsistencies.push(format!(
                "Balance mismatch: {} vs {}",
                self.initial_snapshot.total_balance,
                final_snapshot.total_balance
            ));
        }

        // 4. Raft 상태 검증
        if final_snapshot.log_index < self.initial_snapshot.log_index {
            report.data_loss += 1;
            report.inconsistencies.push(format!(
                "Log index regressed: {} vs {}",
                self.initial_snapshot.log_index,
                final_snapshot.log_index
            ));
        }

        report.success = report.data_loss == 0 && report.data_corruption == 0;
        report
    }

    /// 데이터 일관성 확인
    pub fn check_consistency(&self, current: &SystemSnapshot) -> bool {
        current.account_count > 0
            && current.transaction_count >= 0
            && current.total_balance > 0
            && current.log_index >= self.initial_snapshot.log_index
    }
}

/// 무결성 리포트
#[derive(Clone, Debug)]
pub struct IntegrityReport {
    pub data_loss: u32,
    pub data_corruption: u32,
    pub inconsistencies: Vec<String>,
    pub success: bool,
}

impl IntegrityReport {
    pub fn print_summary(&self) {
        info!("\n╔════════════════════════════════════════╗");
        info!("║  Integrity Verification Report       ║");
        info!("╠════════════════════════════════════════╣");
        info!("║ Status: {:<31} ║", if self.success { "✅ PASS" } else { "❌ FAIL" });
        info!("║ Data Loss: {:<28} ║", self.data_loss);
        info!("║ Data Corruption: {:<22} ║", self.data_corruption);
        info!("╚════════════════════════════════════════╝");

        if !self.inconsistencies.is_empty() {
            info!("\nInconsistencies:");
            for (i, issue) in self.inconsistencies.iter().enumerate() {
                info!("  {}. {}", i + 1, issue);
            }
        }
    }
}

/// 복구 상태 추적기
pub struct RecoveryTracker {
    monitors: Arc<DashMap<String, RecoveryMonitor>>,
}

impl RecoveryTracker {
    pub fn new() -> Self {
        RecoveryTracker {
            monitors: Arc::new(DashMap::new()),
        }
    }

    /// 새로운 복구 프로세스 시작
    pub fn start_tracking(&self, test_id: String) {
        self.monitors.insert(test_id.clone(), RecoveryMonitor::new());
        info!("[CHAOS] Tracking recovery for test: {}", test_id);
    }

    /// 특정 테스트의 모니터 조회 및 업데이트
    pub fn update_metrics<F>(&self, test_id: &str, updater: F)
    where
        F: FnOnce(&mut RecoveryMonitor),
    {
        if let Some(mut monitor) = self.monitors.get_mut(test_id) {
            updater(&mut monitor);
        }
    }

    /// 복구 메트릭 조회
    pub fn get_metrics(&self, test_id: &str) -> Option<RecoveryMetrics> {
        self.monitors.get(test_id).map(|m| m.metrics())
    }

    /// 모든 테스트 결과 조회
    pub fn get_all_results(&self) -> Vec<(String, RecoveryMetrics)> {
        self.monitors
            .iter()
            .map(|r| (r.key().clone(), r.value().metrics()))
            .collect()
    }

    /// 평균 복구 시간
    pub fn average_recovery_time(&self) -> u64 {
        let all_results = self.get_all_results();
        if all_results.is_empty() {
            return 0;
        }

        let total: u64 = all_results.iter().map(|(_, m)| m.time_to_recovery_ms).sum();
        total / all_results.len() as u64
    }

    /// 복구 성공률
    pub fn recovery_success_rate(&self) -> f64 {
        let all_results = self.get_all_results();
        if all_results.is_empty() {
            return 0.0;
        }

        let success_count = all_results.iter().filter(|(_, m)| m.recovery_success).count();
        (success_count as f64 / all_results.len() as f64) * 100.0
    }
}

impl Default for RecoveryTracker {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_recovery_monitor() {
        let mut monitor = RecoveryMonitor::new();
        monitor.mark_detected();
        assert!(monitor.metrics().time_to_detect_ms > 0);

        monitor.mark_recovered();
        assert!(monitor.metrics().recovery_success);
    }

    #[test]
    fn test_integrity_validator() {
        let snapshot = SystemSnapshot::new(1, 1, 100, 10, 50, 1000);
        let validator = IntegrityValidator::new(snapshot.clone());

        let final_snapshot = SystemSnapshot::new(1, 1, 101, 10, 51, 1000);
        let report = validator.verify_after_recovery(&final_snapshot);
        assert!(report.success);
    }

    #[test]
    fn test_recovery_tracker() {
        let tracker = RecoveryTracker::new();
        tracker.start_tracking("test1".to_string());

        let metrics = tracker.get_metrics("test1");
        assert!(metrics.is_some());
    }
}
