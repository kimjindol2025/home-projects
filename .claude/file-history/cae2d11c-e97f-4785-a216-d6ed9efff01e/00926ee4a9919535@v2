/// Phase J: Supply Chain Security - Integration Engine
/// 모든 공급망 보안 컴포넌트를 통합하는 상위 엔진

use super::dependency_checker::{DependencyChecker, DependencyCheckResult};
use super::sbom_generator::{SBOMGenerator, SBOM};
use super::vulnerability_scanner::{VulnerabilityScanner, VulnerabilityScanResult};
use super::regression_tester::{RegressionTester, RegressionTestSuite};
use super::audit_logger::AuditLogger;

/// 공급망 보안 보고서
#[derive(Clone, Debug)]
pub struct SupplyChainSecurityReport {
    pub timestamp: String,
    pub project_name: String,
    pub project_version: String,
    pub overall_risk_level: String,
    pub dependency_check: DependencyCheckResult,
    pub vulnerability_scan: VulnerabilityScanResult,
    pub issues_found: usize,
    pub recommendations: Vec<String>,
}

/// 공급망 보안 엔진
pub struct SupplyChainSecurityEngine {
    project_name: String,
    project_version: String,
    cargo_path: String,
    dependency_checker: DependencyChecker,
    sbom_generator: SBOMGenerator,
    vulnerability_scanner: VulnerabilityScanner,
    regression_tester: RegressionTester,
    audit_logger: AuditLogger,
}

impl SupplyChainSecurityEngine {
    /// 새로운 공급망 보안 엔진 생성
    pub fn new(
        project_name: &str,
        project_version: &str,
        cargo_path: &str,
    ) -> Self {
        SupplyChainSecurityEngine {
            project_name: project_name.to_string(),
            project_version: project_version.to_string(),
            cargo_path: cargo_path.to_string(),
            dependency_checker: DependencyChecker::new(),
            sbom_generator: SBOMGenerator::new(project_name, project_version),
            vulnerability_scanner: VulnerabilityScanner::new(),
            regression_tester: RegressionTester::new(),
            audit_logger: AuditLogger::new(),
        }
    }

    /// 전체 공급망 보안 검사 실행
    pub async fn run_full_security_check(&mut self) -> SupplyChainSecurityReport {
        println!("\n╔════════════════════════════════════════════════════════╗");
        println!("║    Phase J: Supply Chain Security Assessment (Full)   ║");
        println!("╚════════════════════════════════════════════════════════╝\n");

        // 1단계: 의존성 로드 및 검증
        println!("📋 Stage 1: Loading and validating dependencies...");
        self.dependency_checker
            .load_from_cargo_toml(&self.cargo_path)
            .unwrap_or_else(|e| println!("Warning: {}", e));

        let dependency_check = self.dependency_checker.validate_all();
        println!(
            "  ✅ Found {} dependencies ({} high-risk)",
            dependency_check.total_dependencies,
            dependency_check.issues.len()
        );

        // 2단계: SBOM 생성
        println!("\n📄 Stage 2: Generating SBOM...");
        let sbom = self.sbom_generator.generate_from_checker(&self.dependency_checker);
        self.audit_logger.log_sbom_generated(sbom.components.len());
        println!("  ✅ SBOM generated with {} components", sbom.components.len());

        // 3단계: 취약점 스캔
        println!("\n🔍 Stage 3: Running vulnerability scanner...");
        let crates: Vec<(String, String)> = self
            .dependency_checker
            .get_all_dependencies()
            .iter()
            .map(|d| (d.name.clone(), d.version.clone()))
            .collect();

        let vulnerability_scan = self.vulnerability_scanner.scan_all_crates(&crates);
        self.audit_logger.log_security_scan(
            vulnerability_scan.scanned_crates,
            vulnerability_scan.total_vulnerabilities,
        );

        for vuln in &vulnerability_scan.vulnerabilities {
            self.audit_logger.log_vulnerability_detected(
                &vuln.crate_name,
                &vuln.id,
                &vuln.severity.to_string(),
            );
        }

        println!(
            "  ✅ Scanned {} crates, found {} vulnerabilities",
            vulnerability_scan.scanned_crates, vulnerability_scan.total_vulnerabilities
        );

        // 4단계: 회귀 테스트 (선택사항)
        println!("\n🧪 Stage 4: Running regression tests...");
        self.run_regression_tests().await;

        // 5단계: 컴플라이언스 검사
        println!("\n✅ Stage 5: Compliance check...");
        let issues_found = dependency_check.issues.len() + vulnerability_scan.total_vulnerabilities;
        self.audit_logger.log_compliance_check(issues_found == 0, issues_found);

        // 보고서 생성
        let mut recommendations = Vec::new();

        if !dependency_check.issues.is_empty() {
            for issue in &dependency_check.issues {
                recommendations.push(issue.recommendation.clone());
            }
        }

        if vulnerability_scan.total_vulnerabilities > 0 {
            recommendations.push("Update vulnerable dependencies immediately".to_string());
        }

        if recommendations.is_empty() {
            recommendations.push("✅ No issues found - supply chain is secure".to_string());
        }

        let overall_risk = format!("{:?}", dependency_check.overall_risk);

        SupplyChainSecurityReport {
            timestamp: chrono::Utc::now().to_rfc3339(),
            project_name: self.project_name.clone(),
            project_version: self.project_version.clone(),
            overall_risk_level: overall_risk,
            dependency_check,
            vulnerability_scan,
            issues_found,
            recommendations,
        }
    }

    /// 회귀 테스트 실행
    async fn run_regression_tests(&mut self) {
        let test_cases = RegressionTester::create_core_compatibility_tests();
        let mut tester = RegressionTester::new();

        for test_case in test_cases {
            tester.add_test_case(test_case);
        }

        // 토쿠 업그레이드 테스트 시뮬레이션
        let suite = tester
            .run_regression_test("tokio", "1.34.0", "1.35.0")
            .await;

        let check = tester.check_compatibility(&suite);
        self.audit_logger.log_regression_test(
            &suite.dependency,
            &suite.old_version,
            &suite.new_version,
            check.is_compatible,
            suite.test_results.len(),
        );
    }

    /// SBOM 내보내기 (JSON)
    pub fn export_sbom_json(&self) -> String {
        let sbom = self.sbom_generator.generate_from_checker(&self.dependency_checker);
        sbom.to_json()
    }

    /// SBOM 내보내기 (CSV)
    pub fn export_sbom_csv(&self) -> String {
        let sbom = self.sbom_generator.generate_from_checker(&self.dependency_checker);
        sbom.to_csv()
    }

    /// SBOM 내보내기 (Markdown)
    pub fn export_sbom_markdown(&self) -> String {
        let sbom = self.sbom_generator.generate_from_checker(&self.dependency_checker);
        sbom.to_markdown()
    }

    /// 감사 보고서 생성
    pub fn generate_audit_report(&self) -> String {
        self.audit_logger.generate_report()
    }

    /// 전체 보안 보고서 생성
    pub fn generate_security_report(&self, report: &SupplyChainSecurityReport) -> String {
        let mut output = String::new();

        output.push_str("\n╔════════════════════════════════════════════════════════╗\n");
        output.push_str("║      Supply Chain Security Assessment Report           ║\n");
        output.push_str("╠════════════════════════════════════════════════════════╣\n");

        output.push_str(&format!(
            "║ Project: {} v{:<37} ║\n",
            report.project_name, report.project_version
        ));
        output.push_str(&format!(
            "║ Overall Risk: {:<41} ║\n",
            report.overall_risk_level
        ));
        output.push_str(&format!(
            "║ Total Issues: {:<42} ║\n",
            report.issues_found
        ));
        output.push_str(&format!(
            "║ Timestamp: {:<40} ║\n",
            &report.timestamp[..10]
        ));

        output.push_str("╠════════════════════════════════════════════════════════╣\n");
        output.push_str("║ Dependency Analysis:                                   ║\n");
        output.push_str(&format!(
            "║ Total: {} | Valid: {} | High-Risk: {}                   ║\n",
            report.dependency_check.total_dependencies,
            report.dependency_check.valid_count,
            report.dependency_check.issues.len()
        ));

        output.push_str("╠════════════════════════════════════════════════════════╣\n");
        output.push_str("║ Vulnerability Scan:                                    ║\n");
        output.push_str(&format!(
            "║ Scanned: {} | Vulnerable: {} | Critical: {}            ║\n",
            report.vulnerability_scan.scanned_crates,
            report.vulnerability_scan.vulnerable_crates,
            report
                .vulnerability_scan
                .vulnerabilities_by_severity
                .get("Critical")
                .copied()
                .unwrap_or(0)
        ));

        output.push_str("╠════════════════════════════════════════════════════════╣\n");
        output.push_str("║ Recommendations:                                       ║\n");

        for (idx, rec) in report.recommendations.iter().enumerate() {
            if idx < 5 {
                // 최대 5개 권장사항만 표시
                output.push_str(&format!("║ {}. {:<45} ║\n", idx + 1, rec));
            }
        }

        if report.recommendations.len() > 5 {
            output.push_str(&format!(
                "║ ... and {} more recommendations                      ║\n",
                report.recommendations.len() - 5
            ));
        }

        output.push_str("╚════════════════════════════════════════════════════════╝\n");

        output
    }

    /// 종합 평가 결과
    pub fn get_security_assessment(&self) -> String {
        format!(
            "Supply Chain Security Engine Ready\nProject: {} v{}\nCargo Path: {}",
            self.project_name, self.project_version, self.cargo_path
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_engine_creation() {
        let engine = SupplyChainSecurityEngine::new(
            "test_project",
            "1.0.0",
            "Cargo.toml",
        );

        assert_eq!(engine.project_name, "test_project");
        assert_eq!(engine.project_version, "1.0.0");
    }

    #[test]
    fn test_assessment_message() {
        let engine = SupplyChainSecurityEngine::new(
            "test",
            "1.0.0",
            "Cargo.toml",
        );

        let msg = engine.get_security_assessment();
        assert!(msg.contains("test"));
        assert!(msg.contains("1.0.0"));
    }

    #[tokio::test]
    async fn test_exports() {
        let engine = SupplyChainSecurityEngine::new(
            "test",
            "1.0.0",
            "Cargo.toml",
        );

        let json = engine.export_sbom_json();
        assert!(json.contains("\"components\""));

        let csv = engine.export_sbom_csv();
        assert!(csv.contains("Name,Version"));

        let md = engine.export_sbom_markdown();
        assert!(md.contains("# Software Bill of Materials"));
    }
}
