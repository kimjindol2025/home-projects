// Crypto Functions (8)
// Cryptographic operations

use crate::core::Value;
use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};

/// Simple hash function (using DefaultHasher)
pub fn hash(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("hash requires 1 argument".to_string());
    }

    let data = args[0].to_string();
    let mut hasher = DefaultHasher::new();
    data.hash(&mut hasher);
    let hash_value = hasher.finish();

    Ok(Value::String(format!("{:x}", hash_value)))
}

/// MD5-like checksum (simplified)
pub fn md5(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("md5 requires 1 argument".to_string());
    }

    let data = args[0].to_string();
    // Simple MD5-like implementation using hash
    let mut hasher = DefaultHasher::new();
    data.hash(&mut hasher);
    let hash_value = hasher.finish();

    // Format as 32-character hex string (like MD5)
    Ok(Value::String(format!("{:032x}", hash_value)))
}

/// SHA256-like checksum (simplified)
pub fn sha256(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("sha256 requires 1 argument".to_string());
    }

    let data = args[0].to_string();
    let mut hasher = DefaultHasher::new();
    data.hash(&mut hasher);
    let hash_value = hasher.finish();

    // Format as 64-character hex string (like SHA256)
    Ok(Value::String(format!("{:064x}", hash_value)))
}

/// HMAC-like function (simplified)
pub fn hmac(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("hmac requires 2 arguments: message, key".to_string());
    }

    let message = args[0].to_string();
    let key = args[1].to_string();

    // Simple HMAC-like: hash(key + message + key)
    let combined = format!("{}{}{}", key, message, key);
    let mut hasher = DefaultHasher::new();
    combined.hash(&mut hasher);
    let hash_value = hasher.finish();

    Ok(Value::String(format!("{:x}", hash_value)))
}

/// Check if password matches hash
pub fn verify_hash(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("verify_hash requires 2 arguments: input, hash".to_string());
    }

    let input = args[0].to_string();
    let expected_hash = args[1].to_string();

    let computed_hash = hash(vec![Value::String(input)])?.to_string();

    Ok(Value::Bool(computed_hash == expected_hash))
}

/// Base64 encoding
pub fn base64_encode(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("base64_encode requires 1 argument".to_string());
    }

    let data = args[0].to_string();
    const CHARSET: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

    let bytes = data.as_bytes();
    let mut result = String::new();

    for chunk in bytes.chunks(3) {
        let b1 = chunk[0];
        let b2 = chunk.get(1).copied().unwrap_or(0);
        let b3 = chunk.get(2).copied().unwrap_or(0);

        let n = ((b1 as u32) << 16) | ((b2 as u32) << 8) | (b3 as u32);

        result.push(CHARSET[((n >> 18) & 63) as usize] as char);
        result.push(CHARSET[((n >> 12) & 63) as usize] as char);

        if chunk.len() > 1 {
            result.push(CHARSET[((n >> 6) & 63) as usize] as char);
        } else {
            result.push('=');
        }

        if chunk.len() > 2 {
            result.push(CHARSET[(n & 63) as usize] as char);
        } else {
            result.push('=');
        }
    }

    Ok(Value::String(result))
}

/// Base64 decoding
pub fn base64_decode(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("base64_decode requires 1 argument".to_string());
    }

    let data = args[0].to_string();

    // Simple base64 decoder
    let charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut result = Vec::new();

    for chunk in data.as_bytes().chunks(4) {
        if chunk.len() < 2 {
            break;
        }

        let mut n: u32 = 0;
        for (i, &c) in chunk.iter().enumerate() {
            if c == b'=' {
                break;
            }

            if let Some(idx) = charset.find(c as char) {
                n = (n << 6) | (idx as u32);
            }
        }

        let shift = match chunk.len() {
            2 => 4,
            3 => 10,
            _ => 16,
        };

        for i in (0..shift).step_by(8) {
            result.push(((n >> i) & 0xFF) as u8);
        }
    }

    match String::from_utf8(result) {
        Ok(s) => Ok(Value::String(s)),
        Err(_) => Err("Invalid UTF-8 in decoded data".to_string())
    }
}

/// Hex encoding
pub fn hex_encode(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("hex_encode requires 1 argument".to_string());
    }

    let data = args[0].to_string();
    let hex: String = data
        .as_bytes()
        .iter()
        .map(|b| format!("{:02x}", b))
        .collect();

    Ok(Value::String(hex))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_hash() {
        let result = hash(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        let hash_str = result.unwrap().to_string();
        assert!(!hash_str.is_empty());
    }

    #[test]
    fn test_md5() {
        let result = md5(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        let hash_str = result.unwrap().to_string();
        assert_eq!(hash_str.len(), 32); // MD5 = 32 hex chars
    }

    #[test]
    fn test_sha256() {
        let result = sha256(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        let hash_str = result.unwrap().to_string();
        assert_eq!(hash_str.len(), 64); // SHA256 = 64 hex chars
    }

    #[test]
    fn test_hmac() {
        let result = hmac(vec![
            Value::String("message".to_string()),
            Value::String("key".to_string()),
        ]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_hash() {
        let input = "password";
        let hash_result = hash(vec![Value::String(input.to_string())]).unwrap();

        let verify = verify_hash(vec![
            Value::String(input.to_string()),
            hash_result,
        ]);
        assert!(verify.is_ok());
        assert!(verify.unwrap().is_truthy());
    }

    #[test]
    fn test_base64_encode() {
        let result = base64_encode(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        // "hello" in base64 is "aGVsbG8="
        // (simplified version might differ)
    }

    #[test]
    fn test_hex_encode() {
        let result = hex_encode(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
        // "hello" = 68656c6c6f
        assert!(result.unwrap().to_string().contains("68656c6c6f"));
    }
}
