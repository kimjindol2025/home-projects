/// Phase J: Supply Chain Security - SBOM Generator
/// Software Bill of Materials (SBOM) 생성 엔진

use super::dependency_checker::{Dependency, DependencyChecker};
use std::collections::HashMap;
use chrono::Utc;

/// SBOM (Software Bill of Materials) 항목
#[derive(Clone, Debug)]
pub struct SBOMComponent {
    pub component_type: String, // "library", "application", etc.
    pub name: String,
    pub version: String,
    pub supplier: String,
    pub license: String,
    pub download_location: String,
    pub hash: String,
    pub added_date: String,
    pub purpose: String,
}

/// SBOM 메타데이터
#[derive(Clone, Debug)]
pub struct SBOMMetadata {
    pub timestamp: String,
    pub version: String,
    pub tool_name: String,
    pub tool_version: String,
    pub spec_version: String, // CycloneDX or SPDX version
}

/// SBOM 문서
#[derive(Clone, Debug)]
pub struct SBOM {
    pub metadata: SBOMMetadata,
    pub components: Vec<SBOMComponent>,
    pub dependencies_graph: HashMap<String, Vec<String>>,
}

impl SBOM {
    /// JSON 형식으로 변환
    pub fn to_json(&self) -> String {
        let mut json = String::from("{\n");
        json.push_str("  \"version\": \"1.4\",\n");
        json.push_str("  \"format\": \"application/vnd.cyclonedx+json\",\n");
        json.push_str("  \"spec_version\": \"1.4\",\n");

        // 메타데이터
        json.push_str("  \"metadata\": {\n");
        json.push_str(&format!("    \"timestamp\": \"{}\",\n", self.metadata.timestamp));
        json.push_str(&format!("    \"tool\": {{\n"));
        json.push_str(&format!("      \"vendor\": \"Anthropic\",\n"));
        json.push_str(&format!("      \"name\": \"{}\",\n", self.metadata.tool_name));
        json.push_str(&format!("      \"version\": \"{}\"\n", self.metadata.tool_version));
        json.push_str("    }\n");
        json.push_str("  },\n");

        // 컴포넌트
        json.push_str("  \"components\": [\n");
        for (idx, comp) in self.components.iter().enumerate() {
            if idx > 0 {
                json.push_str(",\n");
            }
            json.push_str("    {\n");
            json.push_str(&format!("      \"type\": \"{}\",\n", comp.component_type));
            json.push_str(&format!("      \"name\": \"{}\",\n", comp.name));
            json.push_str(&format!("      \"version\": \"{}\",\n", comp.version));
            json.push_str(&format!("      \"supplier\": \"{}\",\n", comp.supplier));
            json.push_str(&format!("      \"license\": \"{}\",\n", comp.license));
            json.push_str(&format!("      \"downloadLocation\": \"{}\",\n", comp.download_location));
            json.push_str(&format!("      \"hashes\": [{{ \"alg\": \"SHA-256\", \"content\": \"{}\" }}],\n", comp.hash));
            json.push_str(&format!("      \"purpose\": \"{}\"\n", comp.purpose));
            json.push_str("    }");
        }
        json.push_str("\n  ]\n");
        json.push_str("}\n");

        json
    }

    /// CSV 형식으로 변환
    pub fn to_csv(&self) -> String {
        let mut csv = String::from("Name,Version,Type,License,Supplier,Purpose,Added Date\n");

        for comp in &self.components {
            csv.push_str(&format!(
                "\"{}\",\"{}\",\"{}\",\"{}\",\"{}\",\"{}\",\"{}\"\n",
                comp.name,
                comp.version,
                comp.component_type,
                comp.license,
                comp.supplier,
                comp.purpose,
                comp.added_date
            ));
        }

        csv
    }

    /// 의존성 그래프 마크다운으로 생성
    pub fn to_markdown(&self) -> String {
        let mut md = String::from("# Software Bill of Materials (SBOM)\n\n");

        md.push_str(&format!("**Generated**: {}\n\n", self.metadata.timestamp));
        md.push_str(&format!("**Tool**: {} v{}\n\n", self.metadata.tool_name, self.metadata.tool_version));
        md.push_str(&format!("**Total Components**: {}\n\n", self.components.len()));

        md.push_str("## Components\n\n");
        md.push_str("| Name | Version | Type | License | Purpose |\n");
        md.push_str("|------|---------|------|---------|----------|\n");

        for comp in &self.components {
            md.push_str(&format!(
                "| {} | {} | {} | {} | {} |\n",
                comp.name, comp.version, comp.component_type, comp.license, comp.purpose
            ));
        }

        md.push_str("\n## Dependency Tree\n\n");
        md.push_str("```\n");
        md.push_str("distributed_bank\n");
        for comp in &self.components {
            md.push_str(&format!("├── {} v{}\n", comp.name, comp.version));
        }
        md.push_str("```\n");

        md
    }
}

/// SBOM 생성기
pub struct SBOMGenerator {
    project_name: String,
    project_version: String,
}

impl SBOMGenerator {
    /// 새로운 생성기 생성
    pub fn new(project_name: &str, project_version: &str) -> Self {
        SBOMGenerator {
            project_name: project_name.to_string(),
            project_version: project_version.to_string(),
        }
    }

    /// 의존성 검사기로부터 SBOM 생성
    pub fn generate_from_checker(&self, checker: &DependencyChecker) -> SBOM {
        let mut components = Vec::new();

        // 프로젝트 자체를 첫 번째 컴포넌트로 추가
        components.push(SBOMComponent {
            component_type: "application".to_string(),
            name: self.project_name.clone(),
            version: self.project_version.clone(),
            supplier: "Anthropic".to_string(),
            license: "MIT".to_string(),
            download_location: "https://gogs.dclub.kr/kim/freelang-distributed-system".to_string(),
            hash: "unknown".to_string(),
            added_date: Utc::now().to_rfc3339(),
            purpose: "Distributed Banking System".to_string(),
        });

        // 모든 의존성을 컴포넌트로 추가
        for dep in checker.get_all_dependencies() {
            components.push(self.dependency_to_component(dep));
        }

        // 의존성 그래프 생성
        let dependencies_graph = self.build_dependency_graph(checker);

        SBOM {
            metadata: SBOMMetadata {
                timestamp: Utc::now().to_rfc3339(),
                version: "1.0".to_string(),
                tool_name: "Claude Supply Chain Security".to_string(),
                tool_version: "1.0".to_string(),
                spec_version: "CycloneDX 1.4".to_string(),
            },
            components,
            dependencies_graph,
        }
    }

    /// Dependency를 SBOMComponent로 변환
    fn dependency_to_component(&self, dep: &Dependency) -> SBOMComponent {
        let purpose = if dep.is_dev {
            "Development & Testing".to_string()
        } else {
            "Runtime".to_string()
        };

        let license = self.infer_license(&dep.name);

        SBOMComponent {
            component_type: "library".to_string(),
            name: dep.name.clone(),
            version: dep.version.clone(),
            supplier: "Rust Crates Community".to_string(),
            license,
            download_location: format!("https://crates.io/crates/{}", dep.name),
            hash: format!("sha256:{}", self.compute_component_hash(&dep.name, &dep.version)),
            added_date: Utc::now().to_rfc3339(),
            purpose,
        }
    }

    /// 라이선스 추론
    fn infer_license(&self, crate_name: &str) -> String {
        // 실제로는 crates.io API에서 조회
        match crate_name {
            "tokio" => "MIT".to_string(),
            "serde" => "MIT OR Apache-2.0".to_string(),
            "hyper" => "MIT".to_string(),
            "tracing" => "MIT".to_string(),
            "rustls" => "Apache-2.0".to_string(),
            "chacha20poly1305" => "Apache-2.0 OR MIT".to_string(),
            _ => "Unknown".to_string(),
        }
    }

    /// 컴포넌트 해시 계산 (시뮬레이션)
    fn compute_component_hash(&self, name: &str, version: &str) -> String {
        format!("{}_{}_sha256", name, version)
    }

    /// 의존성 그래프 구축
    fn build_dependency_graph(&self, _checker: &DependencyChecker) -> HashMap<String, Vec<String>> {
        let mut graph: HashMap<String, Vec<String>> = HashMap::new();

        // 실제 환경에서는 Cargo.lock에서 구성
        graph.insert(
            self.project_name.clone(),
            vec![
                "tokio".to_string(),
                "serde".to_string(),
                "hyper".to_string(),
                "tracing".to_string(),
                "rustls".to_string(),
            ],
        );

        // 2단계 의존성
        graph.insert(
            "tokio".to_string(),
            vec!["bytes".to_string(), "pin-project".to_string()],
        );

        graph
    }

    /// 라이선스 호환성 검사
    pub fn check_license_compatibility(components: &[SBOMComponent]) -> Vec<String> {
        let mut warnings = Vec::new();

        // GPL과 proprietary 라이선스 혼합 감지
        let has_gpl = components.iter().any(|c| c.license.contains("GPL"));
        let has_proprietary = components.iter().any(|c| c.license.contains("Proprietary"));

        if has_gpl && has_proprietary {
            warnings.push(
                "Warning: GPL and Proprietary licenses detected - compatibility may be affected"
                    .to_string(),
            );
        }

        warnings
    }

    /// 컴포넌트 수명주기 분석
    pub fn analyze_component_lifecycle(components: &[SBOMComponent]) -> HashMap<String, usize> {
        let mut analysis = HashMap::new();

        for comp in components {
            let lifecycle = if comp.added_date.contains("2024") || comp.added_date.contains("2025") || comp.added_date.contains("2026") {
                "Recent".to_string()
            } else if comp.added_date.contains("2023") {
                "Active".to_string()
            } else {
                "Legacy".to_string()
            };

            *analysis.entry(lifecycle).or_insert(0) += 1;
        }

        analysis
    }
}

impl Default for SBOMGenerator {
    fn default() -> Self {
        Self::new("distributed_bank", "1.0.0")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sbom_generator_creation() {
        let generator = SBOMGenerator::new("test_app", "1.0.0");
        assert_eq!(generator.project_name, "test_app");
    }

    #[test]
    fn test_license_inference() {
        let generator = SBOMGenerator::new("test_app", "1.0.0");
        assert_eq!(generator.infer_license("tokio"), "MIT");
        assert!(generator.infer_license("serde").contains("MIT"));
    }

    #[test]
    fn test_sbom_json_generation() {
        let sbom = SBOM {
            metadata: SBOMMetadata {
                timestamp: "2026-03-02T00:00:00Z".to_string(),
                version: "1.0".to_string(),
                tool_name: "Claude".to_string(),
                tool_version: "1.0".to_string(),
                spec_version: "CycloneDX 1.4".to_string(),
            },
            components: vec![SBOMComponent {
                component_type: "library".to_string(),
                name: "test".to_string(),
                version: "1.0".to_string(),
                supplier: "Test".to_string(),
                license: "MIT".to_string(),
                download_location: "https://example.com".to_string(),
                hash: "abc123".to_string(),
                added_date: "2026-03-02".to_string(),
                purpose: "Testing".to_string(),
            }],
            dependencies_graph: HashMap::new(),
        };

        let json = sbom.to_json();
        assert!(json.contains("\"components\""));
        assert!(json.contains("\"test\""));
    }

    #[test]
    fn test_sbom_csv_generation() {
        let sbom = SBOM {
            metadata: SBOMMetadata {
                timestamp: "2026-03-02T00:00:00Z".to_string(),
                version: "1.0".to_string(),
                tool_name: "Claude".to_string(),
                tool_version: "1.0".to_string(),
                spec_version: "CycloneDX 1.4".to_string(),
            },
            components: vec![SBOMComponent {
                component_type: "library".to_string(),
                name: "tokio".to_string(),
                version: "1.35".to_string(),
                supplier: "Rust Community".to_string(),
                license: "MIT".to_string(),
                download_location: "https://crates.io/crates/tokio".to_string(),
                hash: "abc123".to_string(),
                added_date: "2026-03-02".to_string(),
                purpose: "Runtime".to_string(),
            }],
            dependencies_graph: HashMap::new(),
        };

        let csv = sbom.to_csv();
        assert!(csv.contains("tokio"));
        assert!(csv.contains("1.35"));
    }
}
