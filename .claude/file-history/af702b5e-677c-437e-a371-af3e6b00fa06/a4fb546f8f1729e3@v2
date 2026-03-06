# Challenge 14: mail_encryptor.fl (L0 메일 암호화 엔진)
## Zero-Plaintext Email Encryption - Instant AES-256-GCM

**Mission**: 메일이 메모리에 올라오는 순간 즉시 암호화, 절대 평문 저장 금지
**Status**: 🚀 IMPLEMENTATION READY
**Language**: 100% FreeLang (no external crypto libraries)
**Target**: 3,000 lines + 6 unforgiving test cases

---

## 📋 Design Specification

### Core Philosophy
```
Raw Mail (Plaintext)
    ↓ [1ms 이내]
L0Crypto::encrypt()
    ↓
AES-256-GCM 암호화
    ↓
Sovereign-FS::push_chunk() (CAS)
    ↓
검은 상자 (누구도 복호화 불가능)

"메일은 메모리에 올라오는 순간 검은 상자가 된다"
```

---

## 🏗️ Data Structures

### 1. Cryptographic Primitives

```rust
// AES-256-GCM 암호화 키
struct AESKey {
    key: [u8; 32],           // 256-bit key
    nonce: [u8; 12],         // 96-bit nonce
    associated_data: Vec<u8>, // AAD (Additional Authenticated Data)
}

impl AESKey {
    fn new_from_master(master_key: &[u8; 32]) -> Self {
        // PBKDF2 파생
        let mut key = [0u8; 32];
        let mut nonce = [0u8; 12];
        Self::kdf(master_key, &mut key, &mut nonce);

        AESKey {
            key,
            nonce,
            associated_data: Vec::new(),
        }
    }

    fn kdf(master: &[u8; 32], out_key: &mut [u8; 32], out_nonce: &mut [u8; 12]) {
        // PBKDF2-SHA256 (2024 iterations)
        // 실제 구현: hmac_sha256 기반 KDF
    }
}

// AES-256-GCM 암호화 상태
struct AESGCMState {
    key: [u8; 32],
    nonce: [u8; 12],
    counter: u32,           // Block counter
    state: [u8; 16],        // Internal state
    tag: [u8; 16],          // Authentication tag
}

impl AESGCMState {
    fn encrypt_block(&mut self, plaintext: &[u8]) -> Vec<u8> {
        // AES-256 encrypt 1 block
        // GCM GHASH 업데이트
        // 누적 TAG 계산
        let mut ciphertext = vec![0u8; plaintext.len()];
        Self::aes256_encrypt_block(&self.key, &self.nonce, self.counter, plaintext, &mut ciphertext);
        self.counter += 1;
        ciphertext
    }

    fn finalize(&mut self) -> [u8; 16] {
        // GHASH 최종 계산
        // TAG 반환
        self.tag
    }

    fn aes256_encrypt_block(
        key: &[u8; 32],
        nonce: &[u8; 12],
        counter: u32,
        plaintext: &[u8],
        ciphertext: &mut [u8],
    ) {
        // AES S-box 기반 암호화
        // 10 rounds (AES-256)
    }
}
```

### 2. Mail Structure

```rust
// 원본 메일 (평문, 메모리에 최소시간만 존재)
struct RawMail {
    sender: String,
    recipient: String,
    subject: String,
    body: Vec<u8>,
    attachments: Vec<Attachment>,
    timestamp: u64,
}

// 암호화된 메일 (평문 제거)
struct EncryptedMail {
    sender_pubkey_id: [u8; 32],      // 발신자 공개키 해시
    recipient_pubkey_id: [u8; 32],   // 수신자 공개키 해시
    nonce: [u8; 12],                 // GCM nonce
    ciphertext: Vec<u8>,             // AES-256-GCM 암호화된 본문
    attachments_hashes: Vec<[u8; 32]>, // 첨부파일 CAS 해시
    authentication_tag: [u8; 16],    // GCM authentication tag
    timestamp: u64,                  // 타임스탐프
}

impl EncryptedMail {
    fn from_raw(raw: &RawMail, master_key: &[u8; 32]) -> Result<Self, Error> {
        // 1. 메모리에서 평문을 읽음 (1ms 이내)
        // 2. AES-256-GCM 암호화 (L0Crypto)
        // 3. 첨부파일 CAS 저장
        // 4. 원본 평문 메모리 제거
        // 5. 암호화된 메일 반환

        let aes_key = AESKey::new_from_master(master_key);
        let mut gcm = AESGCMState::new(&aes_key);

        // Body 암호화
        let ciphertext = gcm.encrypt(&raw.body)?;
        let tag = gcm.finalize();

        // Attachment CAS 처리
        let attachment_hashes: Vec<[u8; 32]> = raw.attachments
            .iter()
            .map(|att| SovereignFS::push_chunk(&att.data).unwrap())
            .collect();

        Ok(EncryptedMail {
            sender_pubkey_id: Self::hash_pubkey(&raw.sender),
            recipient_pubkey_id: Self::hash_pubkey(&raw.recipient),
            nonce: aes_key.nonce,
            ciphertext,
            attachments_hashes,
            authentication_tag: tag,
            timestamp: raw.timestamp,
        })
    }
}

// Vault 저장 구조
struct MailVault {
    owner_id: [u8; 32],              // 소유자 ID
    encrypted_mails: Vec<MailID>,    // CAS 메일 ID 목록
    master_key: [u8; 32],            // 마스터 키 (격리 저장)
    backup_keys: Vec<[u8; 32]>,      // 백업 키
}

impl MailVault {
    fn seal_and_store(&mut self, raw_mail: &RawMail) -> Result<MailID, Error> {
        // 1. 암호화
        let encrypted = EncryptedMail::from_raw(raw_mail, &self.master_key)?;

        // 2. 직렬화
        let serialized = bincode::serialize(&encrypted)?;

        // 3. Sovereign-FS에 CAS로 저장
        let mail_id = SovereignFS::push_chunk(&serialized)?;

        // 4. 인덱싱
        self.encrypted_mails.push(mail_id);

        Ok(mail_id)
    }

    fn retrieve_and_decrypt(&self, mail_id: MailID) -> Result<RawMail, Error> {
        // 1. CAS에서 검색
        let encrypted_bytes = SovereignFS::get_chunk(mail_id)?;

        // 2. 역직렬화
        let encrypted: EncryptedMail = bincode::deserialize(&encrypted_bytes)?;

        // 3. 복호화
        let mut gcm = AESGCMState::new_from_encrypted(&encrypted, &self.master_key)?;
        let plaintext = gcm.decrypt(&encrypted.ciphertext, &encrypted.authentication_tag)?;

        // 4. 원본 메일 재구성
        Ok(RawMail {
            body: plaintext,
            ..Default::default()
        })
    }
}
```

---

## 🧬 Cryptographic Functions

### AES-256 Encryption (S-box 기반)

```rust
mod AES256 {
    // AES S-box (256 values)
    const SBOX: [u8; 256] = [
        0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, // ...
        // 나머지 248개 값
    ];

    // AES Inverse S-box
    const ISBOX: [u8; 256] = [
        0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, // ...
    ];

    // AES Rcon (Round constants)
    const RCON: [u32; 10] = [
        0x01000000, 0x02000000, 0x04000000, // ...
    ];

    fn sub_bytes(state: &mut [[u8; 4]; 4]) {
        for i in 0..4 {
            for j in 0..4 {
                state[i][j] = SBOX[state[i][j] as usize];
            }
        }
    }

    fn shift_rows(state: &mut [[u8; 4]; 4]) {
        // Row 0: no shift
        // Row 1: shift left 1
        // Row 2: shift left 2
        // Row 3: shift left 3
        let temp = state[1][0];
        state[1][0] = state[1][1];
        state[1][1] = state[1][2];
        state[1][2] = state[1][3];
        state[1][3] = temp;
        // ... 나머지 행 처리
    }

    fn mix_columns(state: &mut [[u8; 4]; 4]) {
        // Galois Field multiplication
        // Matrix multiplication in GF(2^8)
    }

    fn add_round_key(state: &mut [[u8; 4]; 4], round_key: &[[u8; 4]; 4]) {
        for i in 0..4 {
            for j in 0..4 {
                state[i][j] ^= round_key[i][j];
            }
        }
    }

    fn key_expansion(key: &[u8; 32]) -> Vec<u32> {
        let mut w = vec![0u32; 60]; // 10 rounds + 1 = 11 round keys = 44 words

        // First 8 words come from the key
        for i in 0..8 {
            w[i] = u32::from_be_bytes([key[4*i], key[4*i+1], key[4*i+2], key[4*i+3]]);
        }

        // Remaining words
        for i in 8..60 {
            let mut temp = w[i - 1];
            if i % 8 == 0 {
                temp = sub_word(rot_word(temp)) ^ RCON[(i / 8) - 1];
            } else if i % 8 == 4 {
                temp = sub_word(temp);
            }
            w[i] = w[i - 8] ^ temp;
        }

        w
    }

    fn encrypt(plaintext: &[u8; 16], key: &[u8; 32]) -> [u8; 16] {
        let round_keys = key_expansion(key);
        let mut state = [[0u8; 4]; 4];

        // Copy plaintext to state
        for i in 0..16 {
            state[i % 4][i / 4] = plaintext[i];
        }

        // AddRoundKey (round 0)
        let mut rk_idx = 0;
        for i in 0..4 {
            for j in 0..4 {
                let rk = ((round_keys[rk_idx] >> (8 * (3 - j))) & 0xFF) as u8;
                state[i][j] ^= rk;
            }
            rk_idx += 1;
        }

        // 9 main rounds
        for round in 1..10 {
            sub_bytes(&mut state);
            shift_rows(&mut state);
            mix_columns(&mut state);

            // AddRoundKey
            let mut rk_idx = 0;
            for i in 0..4 {
                for j in 0..4 {
                    let rk = ((round_keys[rk_idx] >> (8 * (3 - j))) & 0xFF) as u8;
                    state[i][j] ^= rk;
                }
                rk_idx += 1;
            }
        }

        // Final round (no MixColumns)
        sub_bytes(&mut state);
        shift_rows(&mut state);

        // AddRoundKey
        let mut rk_idx = 0;
        for i in 0..4 {
            for j in 0..4 {
                let rk = ((round_keys[rk_idx] >> (8 * (3 - j))) & 0xFF) as u8;
                state[i][j] ^= rk;
            }
            rk_idx += 1;
        }

        // Copy state to ciphertext
        let mut ciphertext = [0u8; 16];
        for i in 0..16 {
            ciphertext[i] = state[i % 4][i / 4];
        }

        ciphertext
    }
}
```

### GCM (Galois/Counter Mode) GHASH

```rust
mod GHASH {
    // GCM Authentication
    fn ghash(h: &[u8; 16], aad: &[u8], ciphertext: &[u8], auth_tag: &mut [u8; 16]) {
        // 1. AAD 처리
        let aad_blocks = (aad.len() + 15) / 16;
        for i in 0..aad_blocks {
            let mut block = [0u8; 16];
            let len = std::cmp::min(16, aad.len() - i * 16);
            block[..len].copy_from_slice(&aad[i * 16..i * 16 + len]);

            // X_i = (X_{i-1} XOR block_i) * H
            gf128_multiply(auth_tag, &block, h);
            gf128_multiply(auth_tag, auth_tag, h);
        }

        // 2. Ciphertext 처리
        let ct_blocks = (ciphertext.len() + 15) / 16;
        for i in 0..ct_blocks {
            let mut block = [0u8; 16];
            let len = std::cmp::min(16, ciphertext.len() - i * 16);
            block[..len].copy_from_slice(&ciphertext[i * 16..i * 16 + len]);

            gf128_multiply(auth_tag, &block, h);
            gf128_multiply(auth_tag, auth_tag, h);
        }

        // 3. Length 처리 (AAD length || CT length)
        let mut len_block = [0u8; 16];
        let aad_bits = (aad.len() as u64) * 8;
        let ct_bits = (ciphertext.len() as u64) * 8;

        len_block[0..8].copy_from_slice(&aad_bits.to_be_bytes());
        len_block[8..16].copy_from_slice(&ct_bits.to_be_bytes());

        gf128_multiply(auth_tag, &len_block, h);
        gf128_multiply(auth_tag, auth_tag, h);
    }

    fn gf128_multiply(out: &mut [u8; 16], a: &[u8; 16], b: &[u8; 16]) {
        // GF(2^128) multiplication
        let mut v = *b;
        let mut z = [0u8; 16];

        for i in 0..128 {
            let byte = i / 8;
            let bit = 7 - (i % 8);

            if (a[byte] >> bit) & 1 == 1 {
                for j in 0..16 {
                    z[j] ^= v[j];
                }
            }

            let lsb = v[15] & 1;

            for j in (1..16).rev() {
                v[j] = (v[j] >> 1) | ((v[j - 1] & 1) << 7);
            }
            v[0] >>= 1;

            if lsb == 1 {
                v[0] ^= 0xE1; // R (GCM reduction polynomial)
            }
        }

        *out = z;
    }
}
```

---

## 🧪 Test Plan (A1-A6 Unforgiving Tests)

### A1: Basic Encryption/Decryption
```
Input: Raw mail (2KB)
Process:
  1. Encrypt with AES-256-GCM
  2. Decrypt with same key
Assertion:
  - Decrypted == Original (bit-perfect match)
  - Ciphertext != Plaintext
  - Time < 1ms
```

### A2: Plaintext Zero Time (Rule 1)
```
Input: Raw mail (100KB)
Process:
  1. Load plaintext to memory
  2. Start encryption
  3. Measure time plaintext exists in memory
Assertion:
  - Plaintext memory duration < 1ms
  - Metadata cleanup complete
  - No plaintext in heap after encryption
```

### A3: Authentication Tag Verification
```
Input: Encrypted mail with TAG
Process:
  1. Decrypt with correct TAG
  2. Decrypt with altered TAG
Assertion:
  - Correct TAG: Decryption success
  - Altered TAG: Decryption fails (authentication error)
  - No partial decryption leaks
```

### A4: CAS Integration (Deduplication)
```
Input: Two identical mails
Process:
  1. Encrypt mail A -> CAS ID_A
  2. Encrypt mail B (identical content) -> CAS ID_B
  3. Query CAS both IDs
Assertion:
  - ID_A == ID_B (same content hash)
  - Single copy in storage
  - Deduplication effective: 50% storage saved
```

### A5: Master Key Derivation
```
Input: Master key + 2 different nonces
Process:
  1. Derive round keys for mail A
  2. Derive round keys for mail B
  3. Encrypt same plaintext with both
Assertion:
  - Ciphertext A != Ciphertext B
  - Both decrypt correctly with respective keys
  - No key material leakage
```

### A6: Performance (Rule 1 - <1ms)
```
Input: Various mail sizes (1KB, 10KB, 100KB, 1MB)
Process:
  1. Encrypt each size
  2. Measure end-to-end time (plaintext load -> ciphertext storage)
Assertion:
  - 1KB: < 0.1ms
  - 10KB: < 0.5ms
  - 100KB: < 1ms
  - 1MB: < 5ms (optional, stretch goal)
  - Linear scaling (no exponential growth)
```

---

## 📊 Unforgiving Rules

| Rule | Target | Test |
|------|--------|------|
| **R1** | Zero Plaintext (<1ms) | A2 |
| **R2** | Perfect Decryption | A1 |
| **R3** | Auth Tag Integrity | A3 |
| **R4** | CAS Deduplication | A4 |
| **R5** | Key Derivation | A5 |
| **R6** | Performance (<1ms) | A6 |

---

## 📁 File Structure

```
freelang-mail-server/
├── src/
│   ├── core/
│   │   ├── mail_encryptor.fl       (1,200 lines)
│   │   │   ├── AES256Module
│   │   │   ├── GHASHModule
│   │   │   ├── MailEncryptor
│   │   │   └── [A1-A6 tests]
│   │   │
│   │   ├── crypto_primitives.fl    (800 lines)
│   │   │   ├── AESKey
│   │   │   ├── AESGCMState
│   │   │   └── PBKDF2
│   │   │
│   │   └── mail_structure.fl       (600 lines)
│   │       ├── RawMail
│   │       ├── EncryptedMail
│   │       ├── MailVault
│   │       └── [serialization]
│   │
│   └── lib.fl                       (400 lines)
│
├── Cargo.toml
├── MAIL_ENCRYPTOR_DESIGN.md        (this file)
└── TEST_SPECIFICATION.md            (A1-A6 details)
```

---

## 🎯 Implementation Checklist

- [ ] AES-256 S-box, Rcon, key expansion
- [ ] AES encryption (10 rounds)
- [ ] AES decryption (inverse operations)
- [ ] GCM GHASH (GF(2^128) multiplication)
- [ ] PBKDF2 key derivation
- [ ] MailVault seal/store operations
- [ ] CAS integration (Sovereign-FS)
- [ ] Test A1: Basic encrypt/decrypt
- [ ] Test A2: Plaintext duration <1ms
- [ ] Test A3: Authentication tag verification
- [ ] Test A4: CAS deduplication
- [ ] Test A5: Master key derivation
- [ ] Test A6: Performance benchmark

---

**Ready to build the L0 encryption fortress?**

"메일은 메모리에 올라오는 순간 검은 상자가 된다."
