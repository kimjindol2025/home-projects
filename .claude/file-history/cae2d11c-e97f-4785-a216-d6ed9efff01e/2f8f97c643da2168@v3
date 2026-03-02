// ============================================
// TLS Configuration Module
// ============================================
// HTTPS/TLS 1.3 설정 및 인증서 관리

use std::fs;
use std::path::Path;
use tracing::info;

/// TLS 인증서
#[derive(Clone, Debug)]
pub struct Certificate {
    pub data: Vec<u8>,
    pub path: String,
}

/// TLS 개인키
#[derive(Clone, Debug)]
pub struct PrivateKey {
    pub data: Vec<u8>,
    pub path: String,
}

/// TLS 설정
pub struct TlsConfig {
    pub cert: Certificate,
    pub key: PrivateKey,
    pub enabled: bool,
    pub min_version: TlsVersion,
    pub cipher_suites: Vec<String>,
}

/// TLS 버전
#[derive(Clone, Debug, PartialEq, Eq)]
pub enum TlsVersion {
    V12,
    V13,
}

impl TlsConfig {
    /// TLS 설정 생성
    pub fn new(cert_path: &str, key_path: &str) -> Result<Self, String> {
        let cert = Certificate::load(cert_path)?;
        let key = PrivateKey::load(key_path)?;

        info!("🔐 TLS Configuration");
        info!("   Certificate: {}", cert_path);
        info!("   Private Key: {}", key_path);
        info!("   Min Version: TLS 1.3");

        Ok(TlsConfig {
            cert,
            key,
            enabled: true,
            min_version: TlsVersion::V13,
            cipher_suites: vec![
                "TLS_AES_256_GCM_SHA384".to_string(),
                "TLS_CHACHA20_POLY1305_SHA256".to_string(),
                "TLS_AES_128_GCM_SHA256".to_string(),
            ],
        })
    }

    /// 기본 TLS 설정 (테스트용)
    pub fn test_default() -> Result<Self, String> {
        // 테스트용 자체 서명 인증서
        let cert_data = vec![
            48, 130, 2, 90, 48, 130, 1, 146, // ... (PEM 형식 더미)
        ];
        let key_data = vec![
            48, 130, 4, 168, 2, 1, 0, 2, 129, // ... (PEM 형식 더미)
        ];

        Ok(TlsConfig {
            cert: Certificate {
                data: cert_data,
                path: "<test_cert>".to_string(),
            },
            key: PrivateKey {
                data: key_data,
                path: "<test_key>".to_string(),
            },
            enabled: true,
            min_version: TlsVersion::V13,
            cipher_suites: vec![
                "TLS_AES_256_GCM_SHA384".to_string(),
            ],
        })
    }

    /// TLS 활성화 상태
    pub fn is_enabled(&self) -> bool {
        self.enabled
    }

    /// 암호 스위트 추가
    pub fn add_cipher_suite(&mut self, cipher: &str) {
        if !self.cipher_suites.contains(&cipher.to_string()) {
            self.cipher_suites.push(cipher.to_string());
        }
    }

    /// 인증서 정보 조회
    pub fn cert_info(&self) -> CertInfo {
        CertInfo {
            path: self.cert.path.clone(),
            size: self.cert.data.len(),
            loaded: true,
        }
    }

    /// 개인키 정보 조회
    pub fn key_info(&self) -> KeyInfo {
        KeyInfo {
            path: self.key.path.clone(),
            size: self.key.data.len(),
            loaded: true,
        }
    }

    /// TLS 설정 상태 출력
    pub fn print_status(&self) {
        info!("\n🔐 TLS Status:");
        info!("├─ Enabled: {}", self.enabled);
        info!("├─ Min Version: {:?}", self.min_version);
        info!("├─ Cipher Suites: {}", self.cipher_suites.len());
        info!("├─ Certificate: {}", self.cert.path);
        info!("└─ Private Key: {}", self.key.path);
    }
}

/// 인증서 정보
#[derive(Debug, Clone)]
pub struct CertInfo {
    pub path: String,
    pub size: usize,
    pub loaded: bool,
}

/// 개인키 정보
#[derive(Debug, Clone)]
pub struct KeyInfo {
    pub path: String,
    pub size: usize,
    pub loaded: bool,
}

impl Certificate {
    /// 파일에서 인증서 로드
    pub fn load(path: &str) -> Result<Self, String> {
        if !Path::new(path).exists() {
            return Err(format!("Certificate not found: {}", path));
        }

        let data = fs::read(path)
            .map_err(|e| format!("Failed to read certificate: {}", e))?;

        if data.is_empty() {
            return Err("Certificate is empty".to_string());
        }

        Ok(Certificate {
            data,
            path: path.to_string(),
        })
    }

    /// PEM 형식인지 확인
    pub fn is_pem(&self) -> bool {
        let pem_header = b"-----BEGIN CERTIFICATE-----";
        self.data.windows(pem_header.len())
            .any(|w| w == pem_header)
    }

    /// 인증서 크기 반환 (바이트)
    pub fn size(&self) -> usize {
        self.data.len()
    }
}

impl PrivateKey {
    /// 파일에서 개인키 로드
    pub fn load(path: &str) -> Result<Self, String> {
        if !Path::new(path).exists() {
            return Err(format!("Private key not found: {}", path));
        }

        let data = fs::read(path)
            .map_err(|e| format!("Failed to read private key: {}", e))?;

        if data.is_empty() {
            return Err("Private key is empty".to_string());
        }

        Ok(PrivateKey {
            data,
            path: path.to_string(),
        })
    }

    /// PEM 형식인지 확인
    pub fn is_pem(&self) -> bool {
        let pem_headers = [
            b"-----BEGIN PRIVATE KEY-----",
            b"-----BEGIN RSA PRIVATE KEY-----",
        ];

        pem_headers.iter().any(|header| {
            self.data.windows(header.len())
                .any(|w| w == *header)
        })
    }

    /// 개인키 크기 반환 (바이트)
    pub fn size(&self) -> usize {
        self.data.len()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tls_config_default() {
        let config = TlsConfig::test_default().unwrap();
        assert!(config.is_enabled());
        assert_eq!(config.min_version, TlsVersion::V13);
    }

    #[test]
    fn test_tls_cipher_suites() {
        let mut config = TlsConfig::test_default().unwrap();
        let initial_count = config.cipher_suites.len();

        config.add_cipher_suite("TLS_AES_128_CCM_SHA256");
        assert_eq!(config.cipher_suites.len(), initial_count + 1);

        // 중복 추가 확인
        config.add_cipher_suite("TLS_AES_128_CCM_SHA256");
        assert_eq!(config.cipher_suites.len(), initial_count + 1);
    }

    #[test]
    fn test_cert_info() {
        let config = TlsConfig::test_default().unwrap();
        let info = config.cert_info();

        assert!(info.loaded);
        assert!(!info.path.is_empty());
        assert!(info.size > 0);
    }

    #[test]
    fn test_key_info() {
        let config = TlsConfig::test_default().unwrap();
        let info = config.key_info();

        assert!(info.loaded);
        assert!(!info.path.is_empty());
        assert!(info.size > 0);
    }

    #[test]
    fn test_certificate_size() {
        let cert = Certificate {
            data: vec![1, 2, 3, 4, 5],
            path: "test.crt".to_string(),
        };
        assert_eq!(cert.size(), 5);
    }

    #[test]
    fn test_private_key_size() {
        let key = PrivateKey {
            data: vec![1, 2, 3],
            path: "test.key".to_string(),
        };
        assert_eq!(key.size(), 3);
    }
}
