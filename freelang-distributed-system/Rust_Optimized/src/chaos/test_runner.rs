/// Phase I: Chaos Engineering - Test Runner
/// 모든 카오스 엔지니어링 시나리오를 자동으로 실행하고 결과를 집계

use super::injector::ChaosOrchestrator;
use tracing::{info, warn, error};
use super::recovery_validator::RecoveryTracker;
use tracing::{info, warn, error};
use super::scenarios::{
use tracing::{info, warn, error};
    scenario_1_high_latency, scenario_2_network_partition, scenario_3_node_crash,
    scenario_4_disk_error, scenario_5_slow_replication, scenario_6_data_corruption,
    scenario_7_multiple_failures, ChaosTestResult,
};
use std::time::Instant;
use tracing::{info, warn, error};

/// 카오스 테스트 스위트
pub struct ChaosTestSuite {
    orchestrator: ChaosOrchestrator,
    tracker: RecoveryTracker,
    results: Vec<ChaosTestResult>,
}

impl ChaosTestSuite {
    /// 새로운 테스트 스위트 생성
    pub fn new() -> Self {
        ChaosTestSuite {
            orchestrator: ChaosOrchestrator::new(),
            tracker: RecoveryTracker::new(),
            results: Vec::new(),
        }
    }

    /// 모든 시나리오 실행
    pub async fn run_all_scenarios(&mut self) -> ChaosTestSuiteResult {
        info!("\n╔════════════════════════════════════════════════════════════╗");
        info!("║      Chaos Engineering Test Suite (Phase I)              ║");
        info!("║      7 Scenarios | Full Recovery Validation              ║");
        info!("╚════════════════════════════════════════════════════════════╝");

        let suite_start = Instant::now();

        // 시나리오 1: 네트워크 지연
        let result1 = scenario_1_high_latency(&self.orchestrator, &self.tracker).await;
        self.results.push(result1.clone());
        self.print_result(&result1);

        // 시나리오 2: 네트워크 단절
        let result2 = scenario_2_network_partition(&self.orchestrator, &self.tracker).await;
        self.results.push(result2.clone());
        self.print_result(&result2);

        // 시나리오 3: 노드 크래시
        let result3 = scenario_3_node_crash(&self.orchestrator, &self.tracker).await;
        self.results.push(result3.clone());
        self.print_result(&result3);

        // 시나리오 4: 디스크 오류
        let result4 = scenario_4_disk_error(&self.orchestrator, &self.tracker).await;
        self.results.push(result4.clone());
        self.print_result(&result4);

        // 시나리오 5: 느린 복제
        let result5 = scenario_5_slow_replication(&self.orchestrator, &self.tracker).await;
        self.results.push(result5.clone());
        self.print_result(&result5);

        // 시나리오 6: 데이터 손상
        let result6 = scenario_6_data_corruption(&self.tracker).await;
        self.results.push(result6.clone());
        self.print_result(&result6);

        // 시나리오 7: 다중 장애
        let result7 = scenario_7_multiple_failures(&self.orchestrator, &self.tracker).await;
        self.results.push(result7.clone());
        self.print_result(&result7);

        let suite_duration_ms = suite_start.elapsed().as_millis() as u64;

        // 결과 집계
        self.aggregate_results(suite_duration_ms)
    }

    /// 개별 시나리오 실행 (지정된 시나리오만)
    pub async fn run_scenario(&mut self, scenario_id: u32) -> Option<ChaosTestResult> {
        match scenario_id {
            1 => {
                let result = scenario_1_high_latency(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            2 => {
                let result = scenario_2_network_partition(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            3 => {
                let result = scenario_3_node_crash(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            4 => {
                let result = scenario_4_disk_error(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            5 => {
                let result = scenario_5_slow_replication(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            6 => {
                let result = scenario_6_data_corruption(&self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            7 => {
                let result = scenario_7_multiple_failures(&self.orchestrator, &self.tracker).await;
                self.results.push(result.clone());
                Some(result)
            }
            _ => None,
        }
    }

    /// 결과 집계
    fn aggregate_results(&self, total_duration_ms: u64) -> ChaosTestSuiteResult {
        let total_tests = self.results.len() as u32;
        let passed_tests = self.results.iter().filter(|r| r.passed).count() as u32;
        let failed_tests = total_tests - passed_tests;

        let avg_detection_time = if !self.results.is_empty() {
            let sum: u64 = self.results.iter().map(|r| r.detection_time_ms).sum();
            sum / self.results.len() as u64
        } else {
            0
        };

        let avg_recovery_time = if !self.results.is_empty() {
            let sum: u64 = self.results.iter().map(|r| r.recovery_time_ms).sum();
            sum / self.results.len() as u64
        } else {
            0
        };

        let total_data_loss: u32 = self.results.iter().map(|r| r.data_loss).sum();

        let max_detection_time = self.results.iter().map(|r| r.detection_time_ms).max().unwrap_or(0);
        let max_recovery_time = self.results.iter().map(|r| r.recovery_time_ms).max().unwrap_or(0);

        let success_rate = if total_tests > 0 {
            (passed_tests as f64 / total_tests as f64) * 100.0
        } else {
            0.0
        };

        ChaosTestSuiteResult {
            total_tests,
            passed_tests,
            failed_tests,
            success_rate,
            avg_detection_time_ms: avg_detection_time,
            avg_recovery_time_ms: avg_recovery_time,
            max_detection_time_ms: max_detection_time,
            max_recovery_time_ms: max_recovery_time,
            total_data_loss: total_data_loss,
            total_duration_ms,
        }
    }

    /// 결과 프린트
    fn print_result(&self, result: &ChaosTestResult) {
        let status = if result.passed { "✅ PASS" } else { "❌ FAIL" };
        info!(
            "\n{} {} | Detection: {}ms | Recovery: {}ms | Data Loss: {}",
            status,
            result.scenario_name,
            result.detection_time_ms,
            result.recovery_time_ms,
            result.data_loss
        );

        if let Some(err) = &result.error_message {
            info!("   Error: {}", err);
        }
    }

    /// 최종 리포트 생성
    pub fn generate_final_report(&self, suite_result: &ChaosTestSuiteResult) -> String {
        let mut report = String::new();

        report.push_str("\n╔════════════════════════════════════════════════════════════╗\n");
        report.push_str("║         Chaos Engineering Test Suite - Final Report       ║\n");
        report.push_str("╠════════════════════════════════════════════════════════════╣\n");

        report.push_str(&format!(
            "║ Total Tests: {:<50} ║\n",
            suite_result.total_tests
        ));
        report.push_str(&format!(
            "║ Passed: {:<54} ║\n",
            suite_result.passed_tests
        ));
        report.push_str(&format!(
            "║ Failed: {:<54} ║\n",
            suite_result.failed_tests
        ));
        report.push_str(&format!(
            "║ Success Rate: {:.1}%{:<42} ║\n",
            suite_result.success_rate, ""
        ));

        report.push_str("╠════════════════════════════════════════════════════════════╣\n");

        report.push_str(&format!(
            "║ Avg Detection Time: {}ms{:<38} ║\n",
            suite_result.avg_detection_time_ms, ""
        ));
        report.push_str(&format!(
            "║ Max Detection Time: {}ms{:<38} ║\n",
            suite_result.max_detection_time_ms, ""
        ));
        report.push_str(&format!(
            "║ Avg Recovery Time: {}ms{:<39} ║\n",
            suite_result.avg_recovery_time_ms, ""
        ));
        report.push_str(&format!(
            "║ Max Recovery Time: {}ms{:<39} ║\n",
            suite_result.max_recovery_time_ms, ""
        ));

        report.push_str("╠════════════════════════════════════════════════════════════╣\n");

        report.push_str(&format!(
            "║ Total Data Loss: {:<45} ║\n",
            suite_result.total_data_loss
        ));
        report.push_str(&format!(
            "║ Total Duration: {}ms{:<39} ║\n",
            suite_result.total_duration_ms, ""
        ));

        report.push_str("╚════════════════════════════════════════════════════════════╝\n");

        // 상세 시나리오별 결과
        report.push_str("\nDetailed Results by Scenario:\n");
        report.push_str("─────────────────────────────────────────────────────────────\n");

        for (idx, result) in self.results.iter().enumerate() {
            let status = if result.passed { "✅" } else { "❌" };
            report.push_str(&format!(
                "{} S{}: {} | Det: {}ms | Rec: {}ms | Loss: {}\n",
                status,
                idx + 1,
                result.scenario_name,
                result.detection_time_ms,
                result.recovery_time_ms,
                result.data_loss
            ));
        }

        report.push_str("─────────────────────────────────────────────────────────────\n");

        // 권장사항
        report.push_str("\nRecommendations:\n");
        if suite_result.success_rate == 100.0 {
            report.push_str("✅ All chaos engineering tests passed!\n");
            report.push_str("✅ System is production-ready for deployment.\n");
        } else {
            report.push_str(&format!(
                "⚠️  {} scenarios failed. Review and fix before production deployment.\n",
                suite_result.failed_tests
            ));
        }

        if suite_result.avg_recovery_time_ms > 5000 {
            report.push_str("⚠️  Average recovery time exceeds 5 seconds. Consider optimization.\n");
        }

        if suite_result.total_data_loss > 0 {
            report.push_str(&format!(
                "❌ Data loss detected in {} items. Data integrity must be verified.\n",
                suite_result.total_data_loss
            ));
        }

        report.push_str("\n");
        report
    }

    /// 결과 조회
    pub fn get_results(&self) -> &[ChaosTestResult] {
        &self.results
    }

    /// 결과 JSON 형식 내보내기
    pub fn export_json(&self) -> String {
        let mut json = String::from("{\"test_results\": [");

        for (idx, result) in self.results.iter().enumerate() {
            if idx > 0 {
                json.push(',');
            }
            json.push_str(&format!(
                "{{\"scenario\":\"{}\",\"passed\":{},\"detection_ms\":{},\"recovery_ms\":{},\"data_loss\":{}}}",
                result.scenario_name,
                result.passed,
                result.detection_time_ms,
                result.recovery_time_ms,
                result.data_loss
            ));
        }

        json.push_str("]}");
        json
    }
}

impl Default for ChaosTestSuite {
    fn default() -> Self {
        Self::new()
    }
}

/// 테스트 스위트 결과 집계
#[derive(Clone, Debug)]
pub struct ChaosTestSuiteResult {
    pub total_tests: u32,
    pub passed_tests: u32,
    pub failed_tests: u32,
    pub success_rate: f64,
    pub avg_detection_time_ms: u64,
    pub avg_recovery_time_ms: u64,
    pub max_detection_time_ms: u64,
    pub max_recovery_time_ms: u64,
    pub total_data_loss: u32,
    pub total_duration_ms: u64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_suite_creation() {
        let suite = ChaosTestSuite::new();
        assert_eq!(suite.results.len(), 0);
    }

    #[tokio::test]
    async fn test_single_scenario_execution() {
        let mut suite = ChaosTestSuite::new();
        let result = suite.run_scenario(1).await;
        assert!(result.is_some());
        assert_eq!(suite.results.len(), 1);
    }

    #[test]
    fn test_result_aggregation() {
        let suite = ChaosTestSuite::new();
        let result = suite.aggregate_results(5000);
        assert_eq!(result.total_tests, 0);
        assert_eq!(result.success_rate, 0.0);
    }

    #[test]
    fn test_json_export() {
        let suite = ChaosTestSuite::new();
        let json = suite.export_json();
        assert!(json.contains("test_results"));
    }
}
