/// Phase J: Supply Chain Security - Dependency Checker
/// Cargo.toml 파싱 및 의존성 검증

use std::collections::HashMap;
use std::fs;

/// 의존성 정보
#[derive(Clone, Debug)]
pub struct Dependency {
    pub name: String,
    pub version: String,
    pub features: Vec<String>,
    pub is_dev: bool,
    pub risk_level: RiskLevel,
    pub validation_status: ValidationStatus,
}

/// 위험 수준
#[derive(Clone, Debug, PartialEq, Eq)]
pub enum RiskLevel {
    Low,      // 안정적인 크레이트 (1000+ downloads/day)
    Medium,   // 중간 정도 (100-1000 downloads/day)
    High,     // 낮은 다운로드 (10-100 downloads/day)
    Critical, // 매우 낮은 다운로드 또는 unmaintained
}

/// 검증 상태
#[derive(Clone, Debug, PartialEq, Eq)]
pub enum ValidationStatus {
    Valid,              // 검증됨
    VersionMismatch,    // 버전 범위 문제
    Unmaintained,       // 유지보수 중단
    SecurityIssue,      // 보안 문제
    Deprecated,         // 더 이상 사용 안 함
    Unknown,            // 검증 불가
}

/// 의존성 검증 결과
#[derive(Clone, Debug)]
pub struct DependencyCheckResult {
    pub total_dependencies: usize,
    pub valid_count: usize,
    pub risk_count: HashMap<String, usize>,
    pub issues: Vec<DependencyIssue>,
    pub overall_risk: RiskLevel,
}

/// 의존성 문제
#[derive(Clone, Debug)]
pub struct DependencyIssue {
    pub dependency_name: String,
    pub issue_type: String,
    pub severity: String,
    pub recommendation: String,
}

/// 의존성 검사기
pub struct DependencyChecker {
    dependencies: HashMap<String, Dependency>,
    known_vulnerabilities: HashMap<String, Vec<String>>,
    unmaintained_crates: Vec<String>,
}

impl DependencyChecker {
    /// 새로운 검사기 생성
    pub fn new() -> Self {
        let mut checker = DependencyChecker {
            dependencies: HashMap::new(),
            known_vulnerabilities: HashMap::new(),
            unmaintained_crates: vec![],
        };

        // 알려진 취약점 초기화
        checker.initialize_vulnerabilities();

        // 유지보수 중단 크레이트 초기화
        checker.initialize_unmaintained();

        checker
    }

    /// 취약점 데이터 초기화
    fn initialize_vulnerabilities(&mut self) {
        // 실제 환경에서는 CVE 데이터베이스에서 로드
        self.known_vulnerabilities.insert(
            "tokio".to_string(),
            vec!["CVE-2023-1234".to_string()], // 예시
        );
    }

    /// 유지보수 중단 크레이트 초기화
    fn initialize_unmaintained(&mut self) {
        // 실제 환경에서는 crates.io API에서 로드
        self.unmaintained_crates = vec![
            // 예시: 실제로는 더 많음
        ];
    }

    /// Cargo.toml에서 의존성 로드
    pub fn load_from_cargo_toml(&mut self, cargo_path: &str) -> Result<(), String> {
        let content = fs::read_to_string(cargo_path)
            .map_err(|e| format!("Failed to read Cargo.toml: {}", e))?;

        // 간단한 파싱 (실제로는 toml 크레이트 사용)
        self.parse_cargo_content(&content)?;

        Ok(())
    }

    /// Cargo 콘텐츠 파싱
    fn parse_cargo_content(&mut self, content: &str) -> Result<(), String> {
        let mut in_dependencies = false;
        let mut in_dev_dependencies = false;

        for line in content.lines() {
            let trimmed = line.trim();

            if trimmed == "[dependencies]" {
                in_dependencies = true;
                in_dev_dependencies = false;
                continue;
            }

            if trimmed == "[dev-dependencies]" {
                in_dependencies = false;
                in_dev_dependencies = true;
                continue;
            }

            if trimmed.starts_with('[') {
                in_dependencies = false;
                in_dev_dependencies = false;
                continue;
            }

            if (in_dependencies || in_dev_dependencies) && !trimmed.is_empty() && !trimmed.starts_with('#') {
                if let Some(parsed) = self.parse_dependency_line(trimmed, in_dev_dependencies) {
                    self.dependencies.insert(parsed.name.clone(), parsed);
                }
            }
        }

        Ok(())
    }

    /// 의존성 라인 파싱
    fn parse_dependency_line(&self, line: &str, is_dev: bool) -> Option<Dependency> {
        // 형식: name = "version" 또는 name = { version = "...", features = [...] }
        if let Some(eq_pos) = line.find('=') {
            let name = line[..eq_pos].trim().to_string();

            let rest = line[eq_pos + 1..].trim();

            let version = if rest.starts_with('{') {
                // 복잡한 형식
                rest.split("version = \"")
                    .nth(1)
                    .and_then(|v| v.split('"').next())
                    .unwrap_or("unknown")
                    .to_string()
            } else if rest.starts_with('"') {
                // 간단한 형식
                rest.trim_matches('"').to_string()
            } else {
                return None;
            };

            // 피처 추출
            let features = self.extract_features(line);

            // 위험 수준 결정
            let risk_level = self.determine_risk_level(&name);

            return Some(Dependency {
                name,
                version,
                features,
                is_dev,
                risk_level,
                validation_status: ValidationStatus::Unknown,
            });
        }

        None
    }

    /// 피처 추출
    fn extract_features(&self, line: &str) -> Vec<String> {
        if let Some(start) = line.find("features = [") {
            if let Some(end) = line[start..].find(']') {
                let features_str = &line[start + 12..start + end];
                return features_str
                    .split(',')
                    .map(|f| f.trim().trim_matches('"').to_string())
                    .filter(|f| !f.is_empty())
                    .collect();
            }
        }

        Vec::new()
    }

    /// 위험 수준 결정
    fn determine_risk_level(&self, name: &str) -> RiskLevel {
        // 메이저 크레이트는 Low risk
        let major_crates = vec![
            "tokio", "serde", "hyper", "tracing", "uuid", "chrono",
            "parking_lot", "dashmap", "futures", "rustls", "jsonwebtoken"
        ];

        if major_crates.contains(&name) {
            RiskLevel::Low
        } else if self.unmaintained_crates.contains(&name.to_string()) {
            RiskLevel::Critical
        } else if self.known_vulnerabilities.contains_key(name) {
            RiskLevel::High
        } else {
            RiskLevel::Medium
        }
    }

    /// 모든 의존성 검증
    pub fn validate_all(&mut self) -> DependencyCheckResult {
        let mut valid_count = 0;
        let mut risk_count: HashMap<String, usize> = HashMap::new();
        let mut issues = Vec::new();

        for dep in self.dependencies.values_mut() {
            // 버전 범위 검증
            if !self.validate_version(&dep.version) {
                issues.push(DependencyIssue {
                    dependency_name: dep.name.clone(),
                    issue_type: "Invalid version range".to_string(),
                    severity: "HIGH".to_string(),
                    recommendation: "Fix version specification in Cargo.toml".to_string(),
                });
                dep.validation_status = ValidationStatus::VersionMismatch;
            }

            // 취약점 검사
            if self.known_vulnerabilities.contains_key(&dep.name) {
                issues.push(DependencyIssue {
                    dependency_name: dep.name.clone(),
                    issue_type: "Known vulnerability".to_string(),
                    severity: "CRITICAL".to_string(),
                    recommendation: "Update to patched version".to_string(),
                });
                dep.validation_status = ValidationStatus::SecurityIssue;
            }

            // 유지보수 상태 검사
            if self.unmaintained_crates.contains(&dep.name) {
                issues.push(DependencyIssue {
                    dependency_name: dep.name.clone(),
                    issue_type: "Unmaintained crate".to_string(),
                    severity: "MEDIUM".to_string(),
                    recommendation: "Consider migrating to actively maintained alternative".to_string(),
                });
                dep.validation_status = ValidationStatus::Unmaintained;
            }

            if dep.validation_status == ValidationStatus::Unknown {
                dep.validation_status = ValidationStatus::Valid;
                valid_count += 1;
            }

            // 위험 수준별 카운트
            let risk_str = format!("{:?}", dep.risk_level);
            *risk_count.entry(risk_str).or_insert(0) += 1;
        }

        // 전체 위험 수준 결정
        let overall_risk = if !issues.is_empty() {
            RiskLevel::Critical
        } else if risk_count.get("High").copied().unwrap_or(0) > 0 {
            RiskLevel::High
        } else if risk_count.get("Medium").copied().unwrap_or(0) > 0 {
            RiskLevel::Medium
        } else {
            RiskLevel::Low
        };

        DependencyCheckResult {
            total_dependencies: self.dependencies.len(),
            valid_count,
            risk_count,
            issues,
            overall_risk,
        }
    }

    /// 버전 범위 검증
    fn validate_version(&self, version: &str) -> bool {
        // 유효한 버전 형식: 1.0, 1.0.0, 1.0+, ^1.0, ~1.0, >=1.0, etc.
        if version.is_empty() {
            return false;
        }

        // 간단한 검증: 최소한 숫자와 점을 포함해야 함
        version
            .chars()
            .any(|c| c.is_ascii_digit())
    }

    /// 특정 의존성 조회
    pub fn get_dependency(&self, name: &str) -> Option<&Dependency> {
        self.dependencies.get(name)
    }

    /// 모든 의존성 조회
    pub fn get_all_dependencies(&self) -> Vec<&Dependency> {
        self.dependencies.values().collect()
    }

    /// 고위험 의존성만 조회
    pub fn get_high_risk_dependencies(&self) -> Vec<&Dependency> {
        self.dependencies
            .values()
            .filter(|d| d.risk_level == RiskLevel::High || d.risk_level == RiskLevel::Critical)
            .collect()
    }

    /// 의존성 수 조회
    pub fn dependency_count(&self) -> usize {
        self.dependencies.len()
    }

    /// 저장 의존성 수 조회
    pub fn dev_dependency_count(&self) -> usize {
        self.dependencies.values().filter(|d| d.is_dev).count()
    }
}

impl Default for DependencyChecker {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_checker_creation() {
        let checker = DependencyChecker::new();
        assert_eq!(checker.dependency_count(), 0);
    }

    #[test]
    fn test_version_validation() {
        let checker = DependencyChecker::new();
        assert!(checker.validate_version("1.0"));
        assert!(checker.validate_version("1.35"));
        assert!(checker.validate_version("0.7"));
        assert!(!checker.validate_version(""));
    }

    #[test]
    fn test_risk_level_determination() {
        let checker = DependencyChecker::new();
        assert_eq!(checker.determine_risk_level("tokio"), RiskLevel::Low);
        assert_eq!(checker.determine_risk_level("serde"), RiskLevel::Low);
        assert_eq!(checker.determine_risk_level("unknown_crate"), RiskLevel::Medium);
    }

    #[test]
    fn test_parse_dependency_line() {
        let checker = DependencyChecker::new();

        let dep1 = checker.parse_dependency_line("tokio = \"1.35\"", false);
        assert!(dep1.is_some());
        assert_eq!(dep1.unwrap().name, "tokio");

        let dep2 = checker.parse_dependency_line("serde = { version = \"1.0\", features = [\"derive\"] }", false);
        assert!(dep2.is_some());
        let dep2_unwrap = dep2.unwrap();
        assert_eq!(dep2_unwrap.name, "serde");
    }
}
