/// Phase J: Supply Chain Security Module
/// 공급망 보안 검증, SBOM 생성, 취약점 스캔, 회귀 테스트

pub mod dependency_checker;
pub mod sbom_generator;
pub mod vulnerability_scanner;
pub mod regression_tester;
pub mod audit_logger;
pub mod supply_chain;

pub use dependency_checker::{DependencyChecker, Dependency, RiskLevel, ValidationStatus};
pub use sbom_generator::{SBOMGenerator, SBOM, SBOMComponent};
pub use vulnerability_scanner::{VulnerabilityScanner, Vulnerability, VulnerabilitySeverity};
pub use regression_tester::{RegressionTester, RegressionTestSuite, CompatibilityCheckResult};
pub use audit_logger::{AuditLogger, AuditLog, AuditEventType};
pub use supply_chain::{SupplyChainSecurityEngine, SupplyChainSecurityReport};
