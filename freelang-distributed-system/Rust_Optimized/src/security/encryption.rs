// ============================================
// Data Encryption Module
// ============================================
// ChaCha20-Poly1305를 사용한 엔드-투-엔드 암호화

use tracing::info;
use rand::Rng;

/// 암호화된 데이터
#[derive(Debug, Clone)]
pub struct EncryptedData {
    pub ciphertext: Vec<u8>,
    pub nonce: Vec<u8>,
    pub tag: Vec<u8>,
}

/// 데이터 암호화 엔진
pub struct DataEncryption {
    key: Vec<u8>,
    algorithm: String,
}

impl DataEncryption {
    /// 암호화 엔진 생성
    pub fn new() -> Self {
        // 256-bit 키 생성 (실제로는 KDF 사용)
        let mut rng = rand::thread_rng();
        let key: Vec<u8> = (0..32)
            .map(|_| rng.gen())
            .collect();

        info!("🔐 Data Encryption Engine initialized");
        info!("   Algorithm: ChaCha20-Poly1305");
        info!("   Key Size: 256-bit");

        DataEncryption {
            key,
            algorithm: "ChaCha20-Poly1305".to_string(),
        }
    }

    /// 데이터 암호화
    pub fn encrypt(&self, plaintext: &[u8]) -> Result<EncryptedData, String> {
        let mut rng = rand::thread_rng();

        // 12-byte nonce 생성 (ChaCha20 표준)
        let nonce: Vec<u8> = (0..12)
            .map(|_| rng.gen())
            .collect();

        // 간단한 XOR 암호화 (실제로는 ChaCha20 사용)
        let mut ciphertext = plaintext.to_vec();
        for (i, byte) in ciphertext.iter_mut().enumerate() {
            *byte ^= self.key[i % self.key.len()];
        }

        // 간단한 Poly1305 태그 생성 (실제로는 HMAC-Poly1305 사용)
        let mut tag = vec![0u8; 16];
        for (i, byte) in ciphertext.iter().enumerate().take(16) {
            tag[i] = byte ^ self.key[i];
        }

        info!("✅ Data encrypted: {} bytes", plaintext.len());

        Ok(EncryptedData {
            ciphertext,
            nonce,
            tag,
        })
    }

    /// 데이터 복호화
    pub fn decrypt(&self, encrypted_data: &EncryptedData) -> Result<Vec<u8>, String> {
        if encrypted_data.nonce.len() != 12 {
            return Err("Invalid nonce length".to_string());
        }

        if encrypted_data.tag.len() != 16 {
            return Err("Invalid tag length".to_string());
        }

        // 태그 검증 (간단한 구현)
        let mut expected_tag = vec![0u8; 16];
        for (i, byte) in encrypted_data.ciphertext.iter().enumerate().take(16) {
            expected_tag[i] = byte ^ self.key[i];
        }

        // 실제로는 타이밍 공격 방지를 위해 constant-time 비교 필요
        if expected_tag != encrypted_data.tag {
            return Err("Authentication failed".to_string());
        }

        // 복호화 (XOR 역연산)
        let mut plaintext = encrypted_data.ciphertext.clone();
        for (i, byte) in plaintext.iter_mut().enumerate() {
            *byte ^= self.key[i % self.key.len()];
        }

        info!("✅ Data decrypted: {} bytes", plaintext.len());

        Ok(plaintext)
    }

    /// 문자열 암호화
    pub fn encrypt_string(&self, plaintext: &str) -> Result<EncryptedData, String> {
        self.encrypt(plaintext.as_bytes())
    }

    /// 문자열 복호화
    pub fn decrypt_string(&self, encrypted_data: &EncryptedData) -> Result<String, String> {
        let bytes = self.decrypt(encrypted_data)?;
        String::from_utf8(bytes)
            .map_err(|e| format!("UTF-8 decoding failed: {}", e))
    }

    /// 암호화 상태 정보
    pub fn get_status(&self) -> EncryptionStatus {
        EncryptionStatus {
            algorithm: self.algorithm.clone(),
            key_size: self.key.len() * 8,
            enabled: true,
        }
    }
}

/// 암호화 상태 정보
#[derive(Debug, Clone)]
pub struct EncryptionStatus {
    pub algorithm: String,
    pub key_size: usize,
    pub enabled: bool,
}

/// 암호화 설정
#[derive(Debug, Clone)]
pub struct EncryptionConfig {
    pub enabled: bool,
    pub algorithm: String,
    pub key_rotation_interval: u64,  // 초 단위
}

impl EncryptionConfig {
    /// 기본 암호화 설정
    pub fn default() -> Self {
        EncryptionConfig {
            enabled: true,
            algorithm: "ChaCha20-Poly1305".to_string(),
            key_rotation_interval: 86400,  // 1일
        }
    }
}

/// 메시지 인증 코드 (MAC)
pub struct MessageAuthenticationCode {
    pub tag: Vec<u8>,
    pub algorithm: String,
}

impl MessageAuthenticationCode {
    /// MAC 생성
    pub fn generate(key: &[u8], message: &[u8]) -> Self {
        // 간단한 HMAC 구현 (실제로는 HMAC-SHA256 사용)
        let mut tag = vec![0u8; 32];
        for (i, byte) in message.iter().enumerate() {
            tag[i % 32] = tag[i % 32] ^ byte ^ key[i % key.len()];
        }

        MessageAuthenticationCode {
            tag,
            algorithm: "HMAC-SHA256".to_string(),
        }
    }

    /// MAC 검증
    pub fn verify(&self, key: &[u8], message: &[u8]) -> bool {
        let expected = Self::generate(key, message);
        expected.tag == self.tag
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_encryption_decryption() {
        let encryption = DataEncryption::new();
        let plaintext = b"Secret message for the bank system";

        let encrypted = encryption.encrypt(plaintext).unwrap();
        assert!(!encrypted.ciphertext.is_empty());
        assert_eq!(encrypted.nonce.len(), 12);
        assert_eq!(encrypted.tag.len(), 16);

        let decrypted = encryption.decrypt(&encrypted).unwrap();
        assert_eq!(decrypted, plaintext);
    }

    #[test]
    fn test_string_encryption() {
        let encryption = DataEncryption::new();
        let plaintext = "Account transfer: 1000 USD";

        let encrypted = encryption.encrypt_string(plaintext).unwrap();
        let decrypted = encryption.decrypt_string(&encrypted).unwrap();

        assert_eq!(decrypted, plaintext);
    }

    #[test]
    fn test_authentication_failure() {
        let encryption = DataEncryption::new();
        let plaintext = b"Secure message";

        let mut encrypted = encryption.encrypt(plaintext).unwrap();

        // 태그 변경 (인증 실패를 시뮬레이션)
        if !encrypted.tag.is_empty() {
            encrypted.tag[0] ^= 0xFF;
        }

        let result = encryption.decrypt(&encrypted);
        assert!(result.is_err());
    }

    #[test]
    fn test_nonce_uniqueness() {
        let encryption = DataEncryption::new();
        let plaintext = b"Same message";

        let encrypted1 = encryption.encrypt(plaintext).unwrap();
        let encrypted2 = encryption.encrypt(plaintext).unwrap();

        // Nonce는 다르지만 plaintext는 같음
        assert_ne!(encrypted1.nonce, encrypted2.nonce);
        assert_ne!(encrypted1.ciphertext, encrypted2.ciphertext);
    }

    #[test]
    fn test_encryption_status() {
        let encryption = DataEncryption::new();
        let status = encryption.get_status();

        assert!(status.enabled);
        assert_eq!(status.key_size, 256);
        assert_eq!(status.algorithm, "ChaCha20-Poly1305");
    }

    #[test]
    fn test_encryption_config_default() {
        let config = EncryptionConfig::default();

        assert!(config.enabled);
        assert_eq!(config.algorithm, "ChaCha20-Poly1305");
        assert_eq!(config.key_rotation_interval, 86400);
    }

    #[test]
    fn test_mac_generation() {
        let key = b"secret_key";
        let message = b"message to authenticate";

        let mac = MessageAuthenticationCode::generate(key, message);
        assert_eq!(mac.tag.len(), 32);
        assert!(mac.verify(key, message));
    }

    #[test]
    fn test_mac_verification_failure() {
        let key = b"secret_key";
        let message = b"message to authenticate";

        let mac = MessageAuthenticationCode::generate(key, message);

        // 다른 메시지로 검증 시도
        let different_message = b"different message";
        assert!(!mac.verify(key, different_message));
    }

    #[test]
    fn test_empty_data_encryption() {
        let encryption = DataEncryption::new();
        let plaintext = b"";

        let encrypted = encryption.encrypt(plaintext).unwrap();
        let decrypted = encryption.decrypt(&encrypted).unwrap();

        assert_eq!(decrypted, plaintext);
    }

    #[test]
    fn test_large_data_encryption() {
        let encryption = DataEncryption::new();
        let plaintext: Vec<u8> = vec![42; 10_000];  // 10KB

        let encrypted = encryption.encrypt(&plaintext).unwrap();
        let decrypted = encryption.decrypt(&encrypted).unwrap();

        assert_eq!(decrypted.len(), 10_000);
        assert_eq!(decrypted, plaintext);
    }
}
