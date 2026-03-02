// ============================================
// Security Module - TLS, Auth, Encryption
// ============================================
// 분산 은행 시스템의 보안 계층

pub mod tls;
pub mod auth;
pub mod encryption;

pub use self::tls::TlsConfig;
pub use self::auth::{TokenValidator, AuthToken};
pub use self::encryption::DataEncryption;

use std::sync::Arc;
use tracing::{info, warn, error};

/// 보안 설정
#[derive(Clone, Debug)]
pub struct SecurityConfig {
    pub tls_enabled: bool,
    pub cert_path: String,
    pub key_path: String,
    pub jwt_secret: String,
    pub encryption_enabled: bool,
}

impl SecurityConfig {
    /// 보안 설정 생성 (기본값)
    pub fn default() -> Self {
        SecurityConfig {
            tls_enabled: false,
            cert_path: "./certs/server.crt".to_string(),
            key_path: "./certs/server.key".to_string(),
            jwt_secret: "your-secret-key-change-in-production".to_string(),
            encryption_enabled: true,
        }
    }

    /// 보안 설정 생성 (커스텀)
    pub fn new(
        tls_enabled: bool,
        cert_path: &str,
        key_path: &str,
        jwt_secret: &str,
        encryption_enabled: bool,
    ) -> Self {
        SecurityConfig {
            tls_enabled,
            cert_path: cert_path.to_string(),
            key_path: key_path.to_string(),
            jwt_secret: jwt_secret.to_string(),
            encryption_enabled,
        }
    }
}

/// 보안 관리자
pub struct SecurityManager {
    config: Arc<SecurityConfig>,
    tls: Option<Arc<TlsConfig>>,
    auth: Arc<TokenValidator>,
    encryption: Option<Arc<DataEncryption>>,
}

impl SecurityManager {
    /// 보안 관리자 초기화
    pub fn new(config: SecurityConfig) -> Result<Self, String> {
        let tls = if config.tls_enabled {
            let tls_config = TlsConfig::new(&config.cert_path, &config.key_path)?;
            Some(Arc::new(tls_config))
        } else {
            None
        };

        let auth = Arc::new(TokenValidator::new(&config.jwt_secret));

        let encryption = if config.encryption_enabled {
            let enc = DataEncryption::new();
            Some(Arc::new(enc))
        } else {
            None
        };

        info!("🔐 Security Manager initialized");
        info!("   TLS: {}", if config.tls_enabled { "✅ Enabled" } else { "❌ Disabled" });
        info!("   Encryption: {}", if config.encryption_enabled { "✅ Enabled" } else { "❌ Disabled" });

        Ok(SecurityManager {
            config: Arc::new(config),
            tls,
            auth,
            encryption,
        })
    }

    /// TLS 설정 반환
    pub fn tls(&self) -> Option<&Arc<TlsConfig>> {
        self.tls.as_ref()
    }

    /// 토큰 검증기 반환
    pub fn auth(&self) -> &Arc<TokenValidator> {
        &self.auth
    }

    /// 암호화 엔진 반환
    pub fn encryption(&self) -> Option<&Arc<DataEncryption>> {
        self.encryption.as_ref()
    }

    /// 설정 반환
    pub fn config(&self) -> &Arc<SecurityConfig> {
        &self.config
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_security_config_default() {
        let config = SecurityConfig::default();
        assert!(!config.tls_enabled);
        assert!(config.encryption_enabled);
    }

    #[test]
    fn test_security_manager_creation() {
        let config = SecurityConfig::default();
        let manager = SecurityManager::new(config);
        assert!(manager.is_ok());
    }

    #[test]
    fn test_token_generation() {
        let config = SecurityConfig::default();
        let manager = SecurityManager::new(config).unwrap();

        let token = manager.auth().generate("test_user", "admin");
        assert!(!token.is_empty());
    }
}
