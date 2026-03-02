/// Phase J: Supply Chain Security - Vulnerability Scanner
/// 알려진 보안 취약점 스캔 및 경고

use std::collections::HashMap;

/// 취약점 정보
#[derive(Clone, Debug)]
pub struct Vulnerability {
    pub id: String,           // CVE-XXXX-XXXXX
    pub crate_name: String,
    pub affected_versions: Vec<String>,
    pub severity: VulnerabilitySeverity,
    pub published_date: String,
    pub description: String,
    pub cwe_ids: Vec<String>,
}

/// 취약점 심각도
#[derive(Clone, Debug, PartialEq, Eq, PartialOrd, Ord)]
pub enum VulnerabilitySeverity {
    Low,
    Medium,
    High,
    Critical,
}

impl ToString for VulnerabilitySeverity {
    fn to_string(&self) -> String {
        match self {
            VulnerabilitySeverity::Low => "Low".to_string(),
            VulnerabilitySeverity::Medium => "Medium".to_string(),
            VulnerabilitySeverity::High => "High".to_string(),
            VulnerabilitySeverity::Critical => "Critical".to_string(),
        }
    }
}

/// 취약점 스캔 결과
#[derive(Clone, Debug)]
pub struct VulnerabilityScanResult {
    pub scanned_crates: usize,
    pub vulnerable_crates: usize,
    pub total_vulnerabilities: usize,
    pub vulnerabilities_by_severity: HashMap<String, usize>,
    pub vulnerabilities: Vec<Vulnerability>,
}

/// 취약점 스캐너
pub struct VulnerabilityScanner {
    known_vulnerabilities: HashMap<String, Vec<Vulnerability>>,
}

impl VulnerabilityScanner {
    /// 새로운 스캐너 생성
    pub fn new() -> Self {
        let mut scanner = VulnerabilityScanner {
            known_vulnerabilities: HashMap::new(),
        };

        scanner.load_vulnerabilities();

        scanner
    }

    /// 알려진 취약점 로드
    fn load_vulnerabilities(&mut self) {
        // 실제 환경에서는 RustSec, NVD 같은 데이터베이스에서 로드
        // 여기서는 시뮬레이션 데이터

        // tokio 취약점
        self.known_vulnerabilities
            .entry("tokio".to_string())
            .or_insert_with(Vec::new)
            .push(Vulnerability {
                id: "RUSTSEC-2023-0001".to_string(),
                crate_name: "tokio".to_string(),
                affected_versions: vec!["1.0.0".to_string(), "1.34.0".to_string()],
                severity: VulnerabilitySeverity::High,
                published_date: "2023-01-15".to_string(),
                description: "Potential panic in tokio runtime under specific conditions".to_string(),
                cwe_ids: vec!["CWE-248".to_string()],
            });

        // hyper 취약점
        self.known_vulnerabilities
            .entry("hyper".to_string())
            .or_insert_with(Vec::new)
            .push(Vulnerability {
                id: "RUSTSEC-2023-0002".to_string(),
                crate_name: "hyper".to_string(),
                affected_versions: vec!["0.14.0".to_string(), "0.14.26".to_string()],
                severity: VulnerabilitySeverity::Medium,
                published_date: "2023-03-20".to_string(),
                description: "Improper handling of HTTP header injection".to_string(),
                cwe_ids: vec!["CWE-94".to_string()],
            });

        // rustls 취약점
        self.known_vulnerabilities
            .entry("rustls".to_string())
            .or_insert_with(Vec::new)
            .push(Vulnerability {
                id: "RUSTSEC-2024-0001".to_string(),
                crate_name: "rustls".to_string(),
                affected_versions: vec!["0.21.0".to_string()],
                severity: VulnerabilitySeverity::Critical,
                published_date: "2024-01-10".to_string(),
                description: "Certificate validation bypass in TLS handshake".to_string(),
                cwe_ids: vec!["CWE-295".to_string()],
            });

        // chrono 취약점
        self.known_vulnerabilities
            .entry("chrono".to_string())
            .or_insert_with(Vec::new)
            .push(Vulnerability {
                id: "RUSTSEC-2023-0003".to_string(),
                crate_name: "chrono".to_string(),
                affected_versions: vec!["0.4.0".to_string(), "0.4.23".to_string()],
                severity: VulnerabilitySeverity::High,
                published_date: "2023-11-01".to_string(),
                description: "Potential panic in timezone handling".to_string(),
                cwe_ids: vec!["CWE-248".to_string()],
            });
    }

    /// 크레이트 버전에서 취약점 스캔
    pub fn scan_crate(&self, crate_name: &str, version: &str) -> Vec<Vulnerability> {
        let mut results = Vec::new();

        if let Some(vulnerabilities) = self.known_vulnerabilities.get(crate_name) {
            for vuln in vulnerabilities {
                if self.version_matches(&vuln.affected_versions, version) {
                    results.push(vuln.clone());
                }
            }
        }

        results
    }

    /// 모든 크레이트 스캔
    pub fn scan_all_crates(&self, crates: &[(String, String)]) -> VulnerabilityScanResult {
        let mut vulnerabilities = Vec::new();
        let mut vulnerable_crates = 0;
        let mut vulnerabilities_by_severity: HashMap<String, usize> = HashMap::new();

        for (crate_name, version) in crates {
            let crate_vulns = self.scan_crate(crate_name, version);

            if !crate_vulns.is_empty() {
                vulnerable_crates += 1;
            }

            for vuln in crate_vulns {
                let severity_str = vuln.severity.to_string();
                *vulnerabilities_by_severity.entry(severity_str).or_insert(0) += 1;
                vulnerabilities.push(vuln);
            }
        }

        let total_vulnerabilities = vulnerabilities.len();

        VulnerabilityScanResult {
            scanned_crates: crates.len(),
            vulnerable_crates,
            total_vulnerabilities,
            vulnerabilities_by_severity,
            vulnerabilities,
        }
    }

    /// 버전이 영향받는 버전 목록과 일치하는지 확인
    fn version_matches(&self, affected_versions: &[String], current_version: &str) -> bool {
        // 간단한 버전 범위 매칭
        // 실제로는 semver 크레이트 사용
        for affected in affected_versions {
            if affected == current_version {
                return true;
            }

            // 0.x 버전의 경우 마이너 버전까지 일치하면 포함
            if affected.starts_with("0.") && current_version.starts_with("0.") {
                let affected_parts: Vec<&str> = affected.split('.').collect();
                let current_parts: Vec<&str> = current_version.split('.').collect();

                if affected_parts.len() >= 2 && current_parts.len() >= 2 {
                    if affected_parts[0] == current_parts[0] && affected_parts[1] == current_parts[1] {
                        return true;
                    }
                }
            }
        }

        false
    }

    /// 취약점 보고서 생성
    pub fn generate_report(&self, result: &VulnerabilityScanResult) -> String {
        let mut report = String::new();

        report.push_str("\n╔════════════════════════════════════════════════════════╗\n");
        report.push_str("║       Vulnerability Scan Report                      ║\n");
        report.push_str("╠════════════════════════════════════════════════════════╣\n");

        report.push_str(&format!(
            "║ Scanned Crates: {:<37} ║\n",
            result.scanned_crates
        ));
        report.push_str(&format!(
            "║ Vulnerable Crates: {:<32} ║\n",
            result.vulnerable_crates
        ));
        report.push_str(&format!(
            "║ Total Vulnerabilities: {:<28} ║\n",
            result.total_vulnerabilities
        ));

        report.push_str("╠════════════════════════════════════════════════════════╣\n");
        report.push_str("║ Severity Breakdown:                                    ║\n");

        for severity in &["Critical", "High", "Medium", "Low"] {
            let count = result.vulnerabilities_by_severity.get(*severity).copied().unwrap_or(0);
            report.push_str(&format!("║ {}: {:<48} ║\n", severity, count));
        }

        report.push_str("╚════════════════════════════════════════════════════════╝\n");

        if !result.vulnerabilities.is_empty() {
            report.push_str("\nDetailed Vulnerabilities:\n");
            report.push_str("─────────────────────────────────────────────────────────\n");

            for vuln in &result.vulnerabilities {
                report.push_str(&format!("[{}] {} ({})\n", vuln.id, vuln.crate_name, vuln.severity.to_string()));
                report.push_str(&format!("  Description: {}\n", vuln.description));
                report.push_str(&format!("  Affected Versions: {}\n", vuln.affected_versions.join(", ")));
                report.push_str(&format!("  CWEs: {}\n\n", vuln.cwe_ids.join(", ")));
            }
        } else {
            report.push_str("\n✅ No known vulnerabilities detected!\n");
        }

        report
    }

    /// 심각도별 취약점 필터링
    pub fn filter_by_severity(
        &self,
        result: &VulnerabilityScanResult,
        severity: &VulnerabilitySeverity,
    ) -> Vec<Vulnerability> {
        result
            .vulnerabilities
            .iter()
            .filter(|v| v.severity == *severity)
            .cloned()
            .collect()
    }

    /// 특정 크레이트의 취약점 조회
    pub fn get_vulnerabilities_for_crate(&self, crate_name: &str) -> Vec<Vulnerability> {
        self.known_vulnerabilities
            .get(crate_name)
            .cloned()
            .unwrap_or_default()
    }

    /// 취약점 데이터베이스 크기
    pub fn vulnerability_db_size(&self) -> usize {
        self.known_vulnerabilities.values().map(|v| v.len()).sum()
    }
}

impl Default for VulnerabilityScanner {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_scanner_creation() {
        let scanner = VulnerabilityScanner::new();
        assert!(scanner.vulnerability_db_size() > 0);
    }

    #[test]
    fn test_scan_vulnerable_crate() {
        let scanner = VulnerabilityScanner::new();
        let vulns = scanner.scan_crate("tokio", "1.0.0");
        assert!(!vulns.is_empty());
    }

    #[test]
    fn test_scan_safe_crate() {
        let scanner = VulnerabilityScanner::new();
        let vulns = scanner.scan_crate("tokio", "2.0.0");
        assert!(vulns.is_empty());
    }

    #[test]
    fn test_scan_all_crates() {
        let scanner = VulnerabilityScanner::new();
        let crates = vec![
            ("tokio".to_string(), "1.0.0".to_string()),
            ("serde".to_string(), "1.0.0".to_string()),
        ];

        let result = scanner.scan_all_crates(&crates);
        assert_eq!(result.scanned_crates, 2);
        assert!(result.total_vulnerabilities > 0);
    }

    #[test]
    fn test_severity_filter() {
        let scanner = VulnerabilityScanner::new();
        let crates = vec![("tokio".to_string(), "1.0.0".to_string())];

        let result = scanner.scan_all_crates(&crates);
        let critical = scanner.filter_by_severity(&result, &VulnerabilitySeverity::Critical);

        // tokio 1.0.0은 Critical이 아님
        assert!(critical.is_empty());
    }

    #[test]
    fn test_version_matching() {
        let scanner = VulnerabilityScanner::new();

        assert!(scanner.version_matches(
            &vec!["1.0.0".to_string()],
            "1.0.0"
        ));
        assert!(!scanner.version_matches(
            &vec!["1.0.0".to_string()],
            "1.0.1"
        ));
    }
}
