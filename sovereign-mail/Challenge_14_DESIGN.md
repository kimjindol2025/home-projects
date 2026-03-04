# Challenge 14: L0-Mail-Core (종단간 암호화 엔진)
## Zero-Knowledge Encrypted Mail Message Storage

**Status**: 🚀 DESIGN PHASE
**Target Date**: 2026-03-11 (1주)
**Scope**: 4개 모듈, 2,300줄 코드, 30개 테스트, 6개 무관용 규칙

---

## 📋 Challenge Objectives

### Primary Goals
1. **메일이 메모리에 올라오는 순간 즉시 암호화**
   - 평문이 메모리에 노출되는 시간 최소화
   - Sovereign-FS CAS 해시로 자동 변환
   - 누구도 복호화 불가능한 검은 상자 생성

2. **PGP/OpenPGP 표준 준수**
   - 다른 메일 클라이언트와 호환
   - RFC 4880 표준 구현
   - 기존 GPG 키 지원

3. **하드웨어 암호화 가속**
   - L0-Crypto 모듈로 <5ms 암호화
   - CPU 오버헤드 최소화
   - AES-NI, SHA 명령어 활용

4. **키 관리 자동화**
   - 키 생성/저장/순환 자동화
   - Sovereign-FS에 안전하게 저장
   - 키 손상 감지 및 재생성

### Innovation Points
- **메모리 단계 암호화**: 입력 직후 자동 변환
- **CAS 통합**: 암호화된 메시지 → CAS 해시
- **오프라인 안정성**: 네트워크 없이 암호화/복호화
- **재발급 불가능 설계**: 한번 암호화되면 복구 불가능

### Success Criteria
- ✅ 암호화 지연 <5ms (AES-256, 1MB 메시지)
- ✅ 복호화 실패율 0% (정상 키로)
- ✅ 키 교환 시간 <50ms (RSA-4096)
- ✅ 메모리 캐시 <1MB (키 관리용)
- ✅ 암호화 강도: 256-bit minimum
- ✅ 오프라인 메시지 저장 100%

---

## 🏗️ Architecture

### 4-Module L0-Mail-Core Stack

#### Module 1: PGP-Engine (800 lines)
**Purpose**: OpenPGP 표준 메일 암호화/복호화

```
메일 메시지
  ↓ [PGP Engine]
  ├─ 해더 파싱 (From, To, Subject, Date)
  ├─ 본문 추출 (multipart MIME 분석)
  ├─ 첨부파일 처리 (바이너리 encoding)
  └─ 암호화 준비
    ↓
  ├─ 발신자의 개인키로 서명
  ├─ 수신자의 공개키로 암호화
  ├─ PGP 아머 형식 생성
  └─ Sovereign-FS CAS 해시 반환
```

**Components**:
- `MessageParser`: MIME/SMTP 형식 파싱
- `MailComposer`: RFC 5322 메일 작성
- `PGPMessage`: OpenPGP 구조 표현
- `ArmorEncoder`: ASCII armor 형식

**Techniques**:
- RFC 4880 (OpenPGP) 준수
- MIME multipart/signed
- S/MIME compatibility (선택사항)
- Base64 encoding

**Unforgiving Rules (Module 1)**:
- **Rule 1**: 암호화 <5ms (AES-256)
- **Rule 2**: 복호화 실패 0%

---

#### Module 2: Crypto-Primitives (600 lines)
**Purpose**: 저수준 암호화 원시함수 (AES, RSA, SHA)

```
저수준 암호화
  ↓ [Crypto Primitives]
  ├─ 대칭키 암호화
  │   ├─ AES-256-GCM (메시지 본문)
  │   └─ ChaCha20-Poly1305 (대체)
  ├─ 비대칭키 암호화
  │   ├─ RSA-4096 (키 교환)
  │   └─ ECDH (효율적 대안)
  ├─ 해시 함수
  │   ├─ SHA-256 (메시지 다이제스트)
  │   └─ RIPEMD-160 (키 ID 생성)
  └─ MAC (메시지 무결성)
      └─ HMAC-SHA256
```

**Components**:
- `AESCrypto`: AES-256-GCM 구현
- `RSACrypto`: RSA-4096 key agreement
- `HashAlgo`: SHA-256/RIPEMD-160
- `MACGenerator`: HMAC-SHA256

**Techniques**:
- AES-NI (hardware acceleration)
- Constant-time operations (side-channel resistance)
- Key derivation (PBKDF2)
- Random IV/nonce generation

**Unforgiving Rules (Module 2)**:
- **Rule 3**: 키 교환 <50ms
- **Rule 5**: 암호화 강도 256-bit minimum

---

#### Module 3: Key-Manager (500 lines)
**Purpose**: 키 생성, 저장, 순환, 복구

```
키 관리 시스템
  ↓ [Key Manager]
  ├─ 키 생성
  │   ├─ RSA-4096 마스터 키 생성
  │   ├─ AES-256 메시지 키 파생
  │   └─ ECDH subkey 생성
  ├─ 키 저장 (Sovereign-FS)
  │   ├─ 개인키: encrypted at rest
  │   ├─ 공개키: DHT 등록
  │   └─ 키 ID: 발급자 정보
  ├─ 키 순환
  │   ├─ 연간 회전 정책
  │   ├─ 이전 키 아카이브
  │   └─ 새 키 발급
  └─ 복구 (키 손상 시)
      ├─ 백업 키 사용
      └─ 재발급 알림
```

**Components**:
- `KeyGenerator`: RSA/ECDH/AES 키 생성
- `KeyStore`: Sovereign-FS 기반 키 저장소
- `KeyRotation`: 주기적 키 업데이트
- `KeyRecovery`: 손상된 키 복구

**Techniques**:
- PBKDF2 with salt (개인키 암호화)
- Key fingerprinting (SHA-256)
- Expiration policy (annual rotation)
- Backup key slots (3 generations)

**Unforgiving Rules (Module 3)**:
- **Rule 4**: 메모리 캐시 <1MB
- **Rule 6**: 오프라인 메시지 저장 100%

---

#### Module 4: Message-Codec (400 lines)
**Purpose**: 메시지 인코딩/디코딩 (MIME, armor)

```
메시지 변환
  ↓ [Message Codec]
  ├─ MIME 파싱
  │   ├─ multipart/signed 분해
  │   ├─ 본문/첨부 분리
  │   └─ 인코딩 감지 (UTF-8, etc)
  ├─ ASCII Armor
  │   ├─ 바이너리 → Base64
  │   ├─ CRC24 체크섬 추가
  │   └─ PGP 헤더/푸터 삽입
  ├─ 메타데이터 보존
  │   ├─ 발신자 (공개키 ID)
  │   ├─ 수신자 (공개키 ID)
  │   └─ 타임스탐프 (암호화된)
  └─ CAS 변환
      ├─ SHA-256 해시 계산
      └─ Sovereign-FS 블록 ID 반환
```

**Components**:
- `MIMEParser`: RFC 2045 MIME 파싱
- `ArmorCodec`: OpenPGP ASCII armor
- `MetadataPreserver`: 메타데이터 암호화 저장
- `CASConverter`: CAS 해시 변환

**Techniques**:
- Streaming parsing (메모리 효율)
- Base64 in-place encoding
- CRC-24 validation
- Content-Type detection

---

## 🧪 Test Plan (30+ tests)

### Group A: PGP-Engine (10 tests)
```
✓ test_mime_parsing                  - MIME 헤더 파싱
✓ test_multipart_extraction          - multipart/signed 분해
✓ test_signature_creation            - 메시지 서명
✓ test_signature_verification        - 서명 검증
✓ test_encryption                    - 메시지 암호화
✓ test_decryption                    - 메시지 복호화
✓ test_pgp_armor_format              - ASCII armor 형식
✓ test_attachment_handling           - 첨부파일 처리
✓ test_unicode_support               - UTF-8/UTF-16 지원
✓ test_large_message                 - 대용량 메시지 (>100MB)
```

### Group B: Crypto-Primitives (10 tests)
```
✓ test_aes256_gcm_encrypt            - AES-256-GCM 암호화
✓ test_aes256_gcm_decrypt            - AES-256-GCM 복호화
✓ test_aes256_invalid_tag            - 위조된 태그 감지
✓ test_rsa4096_encryption            - RSA-4096 암호화
✓ test_rsa4096_decryption            - RSA-4096 복호화
✓ test_sha256_hash                   - SHA-256 해시
✓ test_hmac_sha256                   - HMAC-SHA256
✓ test_key_derivation                - PBKDF2 키 파생
✓ test_random_nonce                  - 난수 생성 (unique)
✓ test_constant_time                 - Side-channel resistance
```

### Group C: Key-Manager (5 tests)
```
✓ test_key_generation                - RSA-4096 키 생성
✓ test_key_storage                   - Sovereign-FS 저장
✓ test_key_retrieval                 - 키 검색/로드
✓ test_key_rotation                  - 주기적 키 회전
✓ test_key_recovery                  - 손상 키 복구
```

### Group D: Message-Codec (5 tests)
```
✓ test_ascii_armor_encoding          - Base64 armor
✓ test_ascii_armor_decoding          - Armor 파싱
✓ test_crc24_validation              - CRC-24 체크
✓ test_cas_hash_generation           - CAS 해시 생성
✓ test_streaming_encode              - 스트리밍 인코딩
```

---

## 📊 Unforgiving Rules (6 total)

| Rule | Target | Implementation |
|------|--------|-----------------|
| **R1** | Encryption <5ms | PGPEngine::encrypt_message() |
| **R2** | Decryption failures 0% | CryptoPrimitives::aes_decrypt() |
| **R3** | Key exchange <50ms | KeyManager::derive_shared_key() |
| **R4** | Memory cache <1MB | KeyManager::get_memory_usage() |
| **R5** | Crypto strength 256-bit | CryptoPrimitives validate at init |
| **R6** | Offline storage 100% | MessageCodec::cas_convert() |

---

## 📁 File Structure

```
Challenge_14/
├── src/
│   ├── pgp_engine.rs                (800 lines)
│   │   ├── MessageParser
│   │   ├── MailComposer
│   │   ├── PGPMessage
│   │   └── [10 test functions]
│   │
│   ├── crypto_primitives.rs         (600 lines)
│   │   ├── AESCrypto
│   │   ├── RSACrypto
│   │   ├── HashAlgo
│   │   └── [10 test functions]
│   │
│   ├── key_manager.rs               (500 lines)
│   │   ├── KeyGenerator
│   │   ├── KeyStore
│   │   ├── KeyRotation
│   │   └── [5 test functions]
│   │
│   ├── message_codec.rs             (400 lines)
│   │   ├── MIMEParser
│   │   ├── ArmorCodec
│   │   ├── MetadataPreserver
│   │   └── [5 test functions]
│   │
│   └── lib.rs                       (updated)
│
├── Cargo.toml
├── DESIGN.md                        (this file)
└── COMPLETION_REPORT.md             (to be written)
```

---

## 🎯 Implementation Strategy

### Phase A: Crypto Foundations (Days 1-2)
1. AES-256-GCM encryption/decryption
2. RSA-4096 key agreement
3. SHA-256 hashing, HMAC-SHA256
4. Test accuracy & performance

### Phase B: Key Management (Day 3)
1. Key generation (RSA, ECDH, AES)
2. Key storage in Sovereign-FS
3. Key rotation policy
4. Recovery mechanisms

### Phase C: PGP Protocol (Day 4)
1. OpenPGP message format
2. MIME parsing
3. Signature creation/verification
4. Encryption/decryption

### Phase D: Codec & CAS Integration (Day 5)
1. ASCII armor encoding
2. CAS hash generation
3. Metadata preservation
4. Integration testing

### Phase E: Performance Optimization (Day 6)
1. <5ms encryption target
2. Memory efficiency (<1MB)
3. Hardware acceleration (AES-NI)
4. Final testing

---

## 📈 Privacy Guarantees

**Zero-Knowledge Properties**:
- ✅메일이 평문으로 메모리에 존재하지 않음
- ✅ 암호화된 메시지만 저장 (CAS 해시로)
- ✅ 발신자: 수신자의 공개키로 서명 후 암호화
- ✅ 수신자: 개인키로만 복호화 가능
- ✅ 중간자 공격 방지 (서명 검증)
- ✅ 재생 공격 방지 (타임스탐프)

**Example Privacy Flow**:
```
Kim이 Alice에게 메일 작성
  ↓
메모리에 올라옴 (평문)
  ↓ [L0-Crypto <5ms]
Kim의 개인키로 서명
Alice의 공개키로 암호화
  ↓
암호화된 메시지 (평문 없음)
  ↓
Sovereign-FS CAS 해시로 저장
  ↓
누구도 복호화 불가능 (Alice만 가능)
  ↓ [DHT를 통해 Alice에게 전송]
Alice가 수신
  ↓
Alice의 개인키로만 복호화 가능
  ↓
Kim의 서명 검증 (Kim이 보낸 것 확인)
```

---

## 🚀 Expected Outcomes

**Deliverables**:
- 2,300 lines of secure encryption code
- 30 tests (100% coverage)
- 6 unforgiving rules (all satisfied)
- Full E2E encryption pipeline
- CAS integration (no plaintext storage)
- RFC 4880 OpenPGP compliance

**Performance**:
- Encryption: <5ms (AES-256)
- Decryption: <5ms
- Key exchange: <50ms (RSA-4096)
- Memory overhead: <1MB

**Security**:
- 256-bit symmetric encryption
- 4096-bit asymmetric encryption
- Authenticated encryption (AEAD)
- Constant-time operations (no timing attacks)
- Zero plaintext in storage

---

**Next Step**: Implement L0-Mail-Core → Challenge 15 (Sovereign-Naming) → Challenge 16 (L0NN-Mail-Sentry)

**Status**: Design approved, ready for implementation 🔧

**Philosophy**: "메일이 메모리에 올라오는 순간, 검은 상자가 된다."
