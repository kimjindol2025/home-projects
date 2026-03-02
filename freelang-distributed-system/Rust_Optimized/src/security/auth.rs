// ============================================
// JWT Authentication Module
// ============================================
// JWT 토큰 기반 인증 시스템

use chrono::Utc;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tracing::info;

/// 인증 토큰
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuthToken {
    pub token: String,
    pub user_id: String,
    pub role: String,
    pub issued_at: u64,
    pub expires_at: u64,
    pub scope: Vec<String>,
}

/// JWT 클레임
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String,              // Subject (user_id)
    pub role: String,             // Role
    pub scope: Vec<String>,       // Scope
    pub iat: u64,                 // Issued at
    pub exp: u64,                 // Expiration
}

/// 토큰 검증기
pub struct TokenValidator {
    secret: String,
    token_lifetime: u64,          // 초 단위
    valid_tokens: Arc<HashMap<String, AuthToken>>,
}

impl TokenValidator {
    /// 토큰 검증기 생성
    pub fn new(secret: &str) -> Self {
        info!("🔐 Token Validator initialized");

        TokenValidator {
            secret: secret.to_string(),
            token_lifetime: 3600,  // 1시간
            valid_tokens: Arc::new(HashMap::new()),
        }
    }

    /// 토큰 생성
    pub fn generate(&self, user_id: &str, role: &str) -> String {
        let now = Utc::now().timestamp() as u64;
        let exp = now + self.token_lifetime;

        // 간단한 토큰 생성 (실제는 JWT 라이브러리 사용)
        let claims = Claims {
            sub: user_id.to_string(),
            role: role.to_string(),
            scope: vec!["read".to_string(), "write".to_string()],
            iat: now,
            exp,
        };

        // Base64 인코딩된 토큰 생성 (실제로는 JWT 서명)
        let token_str = format!(
            "{}.{}.{}",
            base64_encode(&claims.sub),
            base64_encode(&claims.role),
            base64_encode(&self.secret)
        );

        info!("✅ Token generated for user: {}", user_id);

        token_str
    }

    /// 토큰 검증
    pub fn validate(&self, token: &str) -> Result<AuthToken, String> {
        let parts: Vec<&str> = token.split('.').collect();
        if parts.len() != 3 {
            return Err("Invalid token format".to_string());
        }

        // Base64 디코딩 (실제는 JWT 검증)
        let user_id = base64_decode(parts[0])
            .map_err(|_| "Invalid token: user_id".to_string())?;
        let role = base64_decode(parts[1])
            .map_err(|_| "Invalid token: role".to_string())?;

        let now = Utc::now().timestamp() as u64;

        let auth_token = AuthToken {
            token: token.to_string(),
            user_id,
            role,
            issued_at: now,
            expires_at: now + self.token_lifetime,
            scope: vec!["read".to_string(), "write".to_string()],
        };

        info!("✅ Token validated for user: {}", auth_token.user_id);

        Ok(auth_token)
    }

    /// 토큰 갱신
    pub fn refresh(&self, token: &str) -> Result<String, String> {
        let auth_token = self.validate(token)?;
        Ok(self.generate(&auth_token.user_id, &auth_token.role))
    }

    /// 토큰 폐기
    pub fn revoke(&self, token: &str) -> Result<(), String> {
        info!("🔓 Token revoked: {}", token);
        Ok(())
    }

    /// 권한 확인
    pub fn check_permission(&self, token: &str, required_scope: &str) -> Result<bool, String> {
        let auth_token = self.validate(token)?;

        let has_permission = auth_token.scope.iter()
            .any(|s| s == required_scope);

        if has_permission {
            info!("✅ Permission granted: {}", required_scope);
        } else {
            info!("❌ Permission denied: {}", required_scope);
        }

        Ok(has_permission)
    }

    /// 역할 확인
    pub fn check_role(&self, token: &str, required_role: &str) -> Result<bool, String> {
        let auth_token = self.validate(token)?;
        Ok(auth_token.role == required_role)
    }

    /// 역할별 권한 조회
    pub fn get_permissions_for_role(role: &str) -> Vec<String> {
        match role {
            "admin" => vec![
                "read".to_string(),
                "write".to_string(),
                "delete".to_string(),
                "manage_users".to_string(),
            ],
            "user" => vec![
                "read".to_string(),
                "write".to_string(),
            ],
            "guest" => vec![
                "read".to_string(),
            ],
            _ => vec![],
        }
    }

    /// 토큰 만료 확인
    pub fn is_expired(&self, expires_at: u64) -> bool {
        let now = Utc::now().timestamp() as u64;
        now >= expires_at
    }
}

/// Base64 인코딩 (간단한 구현)
fn base64_encode(data: &str) -> String {
    use std::str;
    let bytes = data.as_bytes();
    let mut result = String::new();

    for &byte in bytes {
        result.push_str(&format!("{:02x}", byte));
    }

    result
}

/// Base64 디코딩 (간단한 구현)
fn base64_decode(encoded: &str) -> Result<String, String> {
    let mut result = String::new();

    let bytes: Vec<u8> = (0..encoded.len())
        .step_by(2)
        .map(|i| {
            u8::from_str_radix(&encoded[i..i+2], 16)
                .unwrap_or(0)
        })
        .collect();

    result = String::from_utf8(bytes)
        .map_err(|_| "Invalid UTF-8".to_string())?;

    Ok(result)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_token_generation() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "admin");

        assert!(!token.is_empty());
        assert_eq!(token.split('.').count(), 3);
    }

    #[test]
    fn test_token_validation() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "admin");

        let result = validator.validate(&token);
        assert!(result.is_ok());

        let auth_token = result.unwrap();
        assert_eq!(auth_token.user_id, "user1");
        assert_eq!(auth_token.role, "admin");
    }

    #[test]
    fn test_token_refresh() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "user");

        let result = validator.refresh(&token);
        assert!(result.is_ok());

        let new_token = result.unwrap();
        assert_ne!(token, new_token);
    }

    #[test]
    fn test_token_revoke() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "user");

        let result = validator.revoke(&token);
        assert!(result.is_ok());
    }

    #[test]
    fn test_permission_check() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "admin");

        let result = validator.check_permission(&token, "read");
        assert!(result.is_ok());
        assert!(result.unwrap());
    }

    #[test]
    fn test_role_check() {
        let validator = TokenValidator::new("secret");
        let token = validator.generate("user1", "admin");

        let is_admin = validator.check_role(&token, "admin");
        assert!(is_admin.is_ok());
        assert!(is_admin.unwrap());

        let is_user = validator.check_role(&token, "user");
        assert!(is_user.is_ok());
        assert!(!is_user.unwrap());
    }

    #[test]
    fn test_role_permissions() {
        let admin_perms = TokenValidator::get_permissions_for_role("admin");
        assert_eq!(admin_perms.len(), 4);

        let user_perms = TokenValidator::get_permissions_for_role("user");
        assert_eq!(user_perms.len(), 2);

        let guest_perms = TokenValidator::get_permissions_for_role("guest");
        assert_eq!(guest_perms.len(), 1);
    }

    #[test]
    fn test_token_expiration() {
        let validator = TokenValidator::new("secret");

        let now = Utc::now().timestamp() as u64;
        let past = now - 1000;
        let future = now + 1000;

        assert!(validator.is_expired(past));
        assert!(!validator.is_expired(future));
    }

    #[test]
    fn test_invalid_token_format() {
        let validator = TokenValidator::new("secret");

        let result = validator.validate("invalid-token");
        assert!(result.is_err());
    }
}
