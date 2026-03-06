// Challenge 14: lib.rs - 모듈 통합
// L0-Mail-Core: Zero-Plaintext Email Encryption
// 2026-03-05, Kim

pub mod crypto_primitives;
pub mod mail_structure;
pub mod mail_encryptor;

pub use mail_encryptor::L0MailEncryptor;
pub use mail_structure::{RawMail, EncryptedMail, MailVault};

/// Challenge 14 초기화 함수
pub fn init_mail_system() -> L0MailEncryptor {
    L0MailEncryptor::new()
}

/// 버전 정보
pub const VERSION: &str = "1.0.0";
pub const PROJECT: &str = "L0-Mail-Core";
pub const CHALLENGE: u8 = 14;

#[cfg(test)]
mod integration_tests {
    use super::*;

    #[test]
    fn test_system_initialization() {
        let system = init_mail_system();
        assert_eq!(system.name, "L0-Mail-Core");
        assert_eq!(system.version, "1.0.0");
    }

    #[test]
    fn test_constants() {
        assert_eq!(CHALLENGE, 14);
        assert_eq!(PROJECT, "L0-Mail-Core");
    }
}
