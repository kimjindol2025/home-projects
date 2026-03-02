/// Phase J: Supply Chain Security - Audit Logger
/// 모든 의존성 변경 및 보안 검사 기록

use chrono::Utc;
use std::collections::HashMap;

/// 감사 로그 항목
#[derive(Clone, Debug)]
pub struct AuditLog {
    pub timestamp: String,
    pub event_type: AuditEventType,
    pub dependency_name: String,
    pub details: HashMap<String, String>,
    pub actor: String,
    pub result: AuditResult,
}

/// 감사 이벤트 유형
#[derive(Clone, Debug, PartialEq, Eq)]
pub enum AuditEventType {
    DependencyAdded,
    DependencyRemoved,
    DependencyUpgraded,
    DependencyDowngraded,
    SecurityScan,
    VulnerabilityDetected,
    RegressionTest,
    SBOMGenerated,
    ComplianceCheck,
}

impl ToString for AuditEventType {
    fn to_string(&self) -> String {
        match self {
            AuditEventType::DependencyAdded => "DEPENDENCY_ADDED".to_string(),
            AuditEventType::DependencyRemoved => "DEPENDENCY_REMOVED".to_string(),
            AuditEventType::DependencyUpgraded => "DEPENDENCY_UPGRADED".to_string(),
            AuditEventType::DependencyDowngraded => "DEPENDENCY_DOWNGRADED".to_string(),
            AuditEventType::SecurityScan => "SECURITY_SCAN".to_string(),
            AuditEventType::VulnerabilityDetected => "VULNERABILITY_DETECTED".to_string(),
            AuditEventType::RegressionTest => "REGRESSION_TEST".to_string(),
            AuditEventType::SBOMGenerated => "SBOM_GENERATED".to_string(),
            AuditEventType::ComplianceCheck => "COMPLIANCE_CHECK".to_string(),
        }
    }
}

/// 감사 결과
#[derive(Clone, Debug, PartialEq, Eq)]
pub enum AuditResult {
    Success,
    Warning,
    Failure,
}

impl ToString for AuditResult {
    fn to_string(&self) -> String {
        match self {
            AuditResult::Success => "SUCCESS".to_string(),
            AuditResult::Warning => "WARNING".to_string(),
            AuditResult::Failure => "FAILURE".to_string(),
        }
    }
}

/// 감사 요약
#[derive(Clone, Debug)]
pub struct AuditSummary {
    pub total_events: usize,
    pub events_by_type: HashMap<String, usize>,
    pub events_by_result: HashMap<String, usize>,
    pub critical_events: usize,
    pub timeline_start: String,
    pub timeline_end: String,
}

/// 감사 로거
pub struct AuditLogger {
    logs: Vec<AuditLog>,
    max_logs: usize,
}

impl AuditLogger {
    /// 새로운 감사 로거 생성
    pub fn new() -> Self {
        AuditLogger {
            logs: Vec::new(),
            max_logs: 10000, // 최대 로그 수
        }
    }

    /// 감사 로그 기록
    pub fn log(&mut self, log: AuditLog) {
        // 순환 로깅 (오래된 로그 삭제)
        if self.logs.len() >= self.max_logs {
            self.logs.remove(0);
        }

        self.logs.push(log);
    }

    /// 의존성 추가 기록
    pub fn log_dependency_added(
        &mut self,
        dependency_name: &str,
        version: &str,
        reason: &str,
    ) {
        let mut details = HashMap::new();
        details.insert("version".to_string(), version.to_string());
        details.insert("reason".to_string(), reason.to_string());

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::DependencyAdded,
            dependency_name: dependency_name.to_string(),
            details,
            actor: "system".to_string(),
            result: AuditResult::Success,
        };

        self.log(log);
        println!(
            "[AUDIT] ✅ Dependency added: {} v{}",
            dependency_name, version
        );
    }

    /// 의존성 제거 기록
    pub fn log_dependency_removed(
        &mut self,
        dependency_name: &str,
        old_version: &str,
        reason: &str,
    ) {
        let mut details = HashMap::new();
        details.insert("old_version".to_string(), old_version.to_string());
        details.insert("reason".to_string(), reason.to_string());

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::DependencyRemoved,
            dependency_name: dependency_name.to_string(),
            details,
            actor: "system".to_string(),
            result: AuditResult::Success,
        };

        self.log(log);
        println!("[AUDIT] 🗑️  Dependency removed: {}", dependency_name);
    }

    /// 의존성 업그레이드 기록
    pub fn log_dependency_upgraded(
        &mut self,
        dependency_name: &str,
        old_version: &str,
        new_version: &str,
        reason: &str,
    ) {
        let mut details = HashMap::new();
        details.insert("old_version".to_string(), old_version.to_string());
        details.insert("new_version".to_string(), new_version.to_string());
        details.insert("reason".to_string(), reason.to_string());

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::DependencyUpgraded,
            dependency_name: dependency_name.to_string(),
            details,
            actor: "system".to_string(),
            result: AuditResult::Success,
        };

        self.log(log);
        println!(
            "[AUDIT] ⬆️  Dependency upgraded: {} {} -> {}",
            dependency_name, old_version, new_version
        );
    }

    /// 보안 스캔 기록
    pub fn log_security_scan(
        &mut self,
        scanned_crates: usize,
        vulnerabilities_found: usize,
    ) {
        let mut details = HashMap::new();
        details.insert("scanned_crates".to_string(), scanned_crates.to_string());
        details.insert(
            "vulnerabilities_found".to_string(),
            vulnerabilities_found.to_string(),
        );

        let result = if vulnerabilities_found > 0 {
            AuditResult::Warning
        } else {
            AuditResult::Success
        };

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::SecurityScan,
            dependency_name: "all".to_string(),
            details,
            actor: "system".to_string(),
            result,
        };

        self.log(log);

        let status = if vulnerabilities_found > 0 { "⚠️ " } else { "✅" };
        println!(
            "[AUDIT] {}Security scan completed: {} crates, {} vulnerabilities",
            status, scanned_crates, vulnerabilities_found
        );
    }

    /// 취약점 감지 기록
    pub fn log_vulnerability_detected(
        &mut self,
        crate_name: &str,
        vulnerability_id: &str,
        severity: &str,
    ) {
        let mut details = HashMap::new();
        details.insert("vulnerability_id".to_string(), vulnerability_id.to_string());
        details.insert("severity".to_string(), severity.to_string());

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::VulnerabilityDetected,
            dependency_name: crate_name.to_string(),
            details,
            actor: "system".to_string(),
            result: AuditResult::Failure,
        };

        self.log(log);
        println!(
            "[AUDIT] 🚨 Vulnerability detected: {} - {} [{}]",
            crate_name, vulnerability_id, severity
        );
    }

    /// 회귀 테스트 결과 기록
    pub fn log_regression_test(
        &mut self,
        dependency: &str,
        old_version: &str,
        new_version: &str,
        passed: bool,
        test_count: usize,
    ) {
        let mut details = HashMap::new();
        details.insert("old_version".to_string(), old_version.to_string());
        details.insert("new_version".to_string(), new_version.to_string());
        details.insert("test_count".to_string(), test_count.to_string());

        let result = if passed {
            AuditResult::Success
        } else {
            AuditResult::Failure
        };

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::RegressionTest,
            dependency_name: dependency.to_string(),
            details,
            actor: "system".to_string(),
            result,
        };

        self.log(log);

        let status = if passed { "✅" } else { "❌" };
        println!(
            "[AUDIT] {} Regression test: {} {} tests",
            status, dependency, test_count
        );
    }

    /// SBOM 생성 기록
    pub fn log_sbom_generated(&mut self, component_count: usize) {
        let mut details = HashMap::new();
        details.insert("component_count".to_string(), component_count.to_string());

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::SBOMGenerated,
            dependency_name: "all".to_string(),
            details,
            actor: "system".to_string(),
            result: AuditResult::Success,
        };

        self.log(log);
        println!(
            "[AUDIT] 📋 SBOM generated: {} components",
            component_count
        );
    }

    /// 컴플라이언스 체크 기록
    pub fn log_compliance_check(
        &mut self,
        passed: bool,
        issues_found: usize,
    ) {
        let mut details = HashMap::new();
        details.insert("issues_found".to_string(), issues_found.to_string());

        let result = if issues_found == 0 {
            AuditResult::Success
        } else {
            AuditResult::Warning
        };

        let log = AuditLog {
            timestamp: Utc::now().to_rfc3339(),
            event_type: AuditEventType::ComplianceCheck,
            dependency_name: "all".to_string(),
            details,
            actor: "system".to_string(),
            result,
        };

        self.log(log);

        let status = if passed { "✅" } else { "⚠️" };
        println!("[AUDIT] {} Compliance check: {} issues", status, issues_found);
    }

    /// 감사 로그 조회
    pub fn get_logs(&self) -> &[AuditLog] {
        &self.logs
    }

    /// 특정 의존성의 로그만 조회
    pub fn get_dependency_logs(&self, dependency_name: &str) -> Vec<&AuditLog> {
        self.logs
            .iter()
            .filter(|log| log.dependency_name == dependency_name)
            .collect()
    }

    /// 특정 이벤트 타입의 로그만 조회
    pub fn get_logs_by_type(&self, event_type: &AuditEventType) -> Vec<&AuditLog> {
        self.logs
            .iter()
            .filter(|log| log.event_type == *event_type)
            .collect()
    }

    /// 감사 요약 생성
    pub fn generate_summary(&self) -> AuditSummary {
        let mut events_by_type: HashMap<String, usize> = HashMap::new();
        let mut events_by_result: HashMap<String, usize> = HashMap::new();
        let mut critical_events = 0;

        for log in &self.logs {
            *events_by_type
                .entry(log.event_type.to_string())
                .or_insert(0) += 1;

            *events_by_result
                .entry(log.result.to_string())
                .or_insert(0) += 1;

            if log.event_type == AuditEventType::VulnerabilityDetected
                || log.result == AuditResult::Failure
            {
                critical_events += 1;
            }
        }

        let timeline_start = self.logs.first().map(|l| l.timestamp.clone()).unwrap_or_default();
        let timeline_end = self.logs.last().map(|l| l.timestamp.clone()).unwrap_or_default();

        AuditSummary {
            total_events: self.logs.len(),
            events_by_type,
            events_by_result,
            critical_events,
            timeline_start,
            timeline_end,
        }
    }

    /// 감사 리포트 생성
    pub fn generate_report(&self) -> String {
        let summary = self.generate_summary();

        let mut report = String::new();

        report.push_str("\n╔════════════════════════════════════════════════════════╗\n");
        report.push_str("║           Supply Chain Audit Report                    ║\n");
        report.push_str("╠════════════════════════════════════════════════════════╣\n");

        report.push_str(&format!(
            "║ Total Events: {:<42} ║\n",
            summary.total_events
        ));
        report.push_str(&format!(
            "║ Critical Events: {:<38} ║\n",
            summary.critical_events
        ));

        report.push_str("╠════════════════════════════════════════════════════════╣\n");
        report.push_str("║ Events by Type:                                        ║\n");

        for (event_type, count) in &summary.events_by_type {
            report.push_str(&format!("║ {}: {:<43} ║\n", event_type, count));
        }

        report.push_str("╠════════════════════════════════════════════════════════╣\n");
        report.push_str("║ Events by Result:                                      ║\n");

        for (result, count) in &summary.events_by_result {
            report.push_str(&format!("║ {}: {:<47} ║\n", result, count));
        }

        report.push_str("╠════════════════════════════════════════════════════════╣\n");
        report.push_str(&format!(
            "║ Timeline: {} to {}  ║\n",
            &summary.timeline_start[..10], &summary.timeline_end[..10]
        ));

        report.push_str("╚════════════════════════════════════════════════════════╝\n");

        report
    }

    /// 감사 로그 CSV 내보내기
    pub fn export_csv(&self) -> String {
        let mut csv = String::from("Timestamp,EventType,DependencyName,Result,Details\n");

        for log in &self.logs {
            let details = log
                .details
                .iter()
                .map(|(k, v)| format!("{}={}", k, v))
                .collect::<Vec<_>>()
                .join("|");

            csv.push_str(&format!(
                "\"{}\",\"{}\",\"{}\",\"{}\",\"{}\"\n",
                log.timestamp,
                log.event_type.to_string(),
                log.dependency_name,
                log.result.to_string(),
                details
            ));
        }

        csv
    }
}

impl Default for AuditLogger {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_audit_logger_creation() {
        let logger = AuditLogger::new();
        assert_eq!(logger.logs.len(), 0);
    }

    #[test]
    fn test_log_dependency_added() {
        let mut logger = AuditLogger::new();
        logger.log_dependency_added("test-crate", "1.0.0", "feature request");

        assert_eq!(logger.logs.len(), 1);
        assert_eq!(logger.logs[0].event_type, AuditEventType::DependencyAdded);
    }

    #[test]
    fn test_log_security_scan() {
        let mut logger = AuditLogger::new();
        logger.log_security_scan(20, 0);

        assert_eq!(logger.logs.len(), 1);
        assert_eq!(logger.logs[0].event_type, AuditEventType::SecurityScan);
    }

    #[test]
    fn test_get_dependency_logs() {
        let mut logger = AuditLogger::new();
        logger.log_dependency_added("tokio", "1.0.0", "required");
        logger.log_dependency_added("serde", "1.0.0", "required");

        let tokio_logs = logger.get_dependency_logs("tokio");
        assert_eq!(tokio_logs.len(), 1);
    }

    #[test]
    fn test_generate_summary() {
        let mut logger = AuditLogger::new();
        logger.log_dependency_added("test", "1.0.0", "test");
        logger.log_security_scan(10, 0);

        let summary = logger.generate_summary();
        assert_eq!(summary.total_events, 2);
    }

    #[test]
    fn test_export_csv() {
        let mut logger = AuditLogger::new();
        logger.log_dependency_added("test", "1.0.0", "test");

        let csv = logger.export_csv();
        assert!(csv.contains("DEPENDENCY_ADDED"));
        assert!(csv.contains("test"));
    }
}
