# Project Sovereign-DNS Phase 4: Security Layer 구현

## Context

Phase 3에서 DHT + Registry + Resolver 3계층 완성 (48개 테스트, 6/8 규칙 달성).
**미달성 규칙 2개**:
- **Rule 3**: ZKP 검증 < 1ms — proof_validator.fl 골격만 존재 (더미 구현)
- **Rule 5**: 하이재킹 방어 100% — detect_hijack_attempt() 항상 false 반환

목표: Security Layer 완성으로 **8/8 무관용 규칙 달성** + 54개 무관용 테스트 완성.

---

## 현재 상태

### 파일 상태
| 파일 | 현재 상태 | 필요 작업 |
|------|----------|----------|
| `src/security/proof_validator.fl` | 121줄 골격 | 400줄로 강화 (ZKP 수학 구현) |
| `src/security/domain_session.fl` | 106줄 기본 | 350줄로 강화 (HMAC + nonce 캐시) |
| `src/security/mod.fl` | 완성 (6줄) | 변경 없음 |
| `tests/security_tests.fl` | 없음 | 600줄 신규 생성 (S1-S6) |

### 미구현 핵심 함수
- `proof_validator.fl::verify_domain_ownership()` → 단순 non-zero 체크만
- `proof_validator.fl::create_domain_proof()` → 하드코딩 더미 값
- `proof_validator.fl::detect_hijack_attempt()` → 항상 false
- `domain_session.fl::verify_session()` → HMAC 검증 없음

---

## 재사용 자산 (freelang-signal)

| 자산 | 파일 경로 | 재사용 패턴 |
|------|----------|-----------|
| `constant_time_eq()` | `freelang-signal/src/zkp/schnorr_proof.fl` | XOR 누적, 타이밍 공격 방지 |
| Fiat-Shamir 패턴 | `freelang-signal/src/zkp/schnorr_proof.fl` | e = H(R || message), s = k - secret*e |
| `sha3_256()` 패턴 | `freelang-signal/src/crypto/sha3_hash.fl` | 24라운드 Keccak 시뮬레이션 |
| `hmac_sha256()` 패턴 | `freelang-signal/src/crypto/sha3_hash.fl` | ipad(0x36)/opad(0x5C) XOR |
| `ReplayGuard` 패턴 | `freelang-signal/src/security/replay_guard.fl` | 1024 슬롯 링 버퍼, TTL 기반 |

---

## 구현 계획

### 1. proof_validator.fl 강화 (121줄 → 400줄)

추가 구현 내용:
- `constant_time_eq()`: XOR 누적, 타이밍 공격 방지
- `sha3_256_dns()`: 24라운드 Keccak 시뮬레이션
- `verify_domain_ownership()` 강화:
  - Step 1: commitment != [0; 32] 체크
  - Step 2: challenge 재계산 = H(commitment || domain_hash || nonce)
  - Step 3: constant_time_eq 로 response 검증
  - Step 4: < 1ms 카운터 (fast_verifications++)
- `create_domain_proof()` 강화:
  - commitment = sha3_256(domain || nonce)
  - challenge = sha3_256(commitment || nonce)
  - response = pseudo-Schnorr 계산
- `detect_hijack_attempt()` 구현:
  - N-gram bigram 엔트로피 분석
  - 엔트로피 > 4.5이면 하이재킹 판정

### 2. domain_session.fl 강화 (106줄 → 350줄)

추가 구현 내용:
- `hmac_sha256_dns()`: ipad(0x36)/opad(0x5C) XOR 패턴
- `detect_replay()` 강화:
  - 8192 nonce 캐시 슬라이딩 윈도우
  - 해시 기반 인덱싱 (nonce % 8192)
  - 링 버퍼 자동 재활용
- `verify_session()` 강화:
  - HMAC 서명 검증
  - 만료 체크 (expires_ns)
  - nonce 재사용 거부 (100% 차단)
- `detect_hijack_attempt()` 구현:
  - N-gram 도메인 이상 탐지
  - 100 시도 카운터 (Rule 5)
- 통계: hijack_attempts / total_verifications / replay_blocks

### 3. tests/security_tests.fl 신규 생성 (600줄)

S1-S6 무관용 테스트:

| 테스트 | 검증 내용 | 무관용 규칙 |
|--------|---------|-----------|
| **S1** | ZKP 증명 생성 + 검증 (< 1ms, 100회) | Rule 3 |
| **S2** | 하이재킹 탐지 (100번 시도 → 100% 감지) | Rule 5 |
| **S3** | nonce 재생 공격 방지 (1000번 → 100% 차단) | Rule 5 |
| **S4** | 세션 생성/만료/검증 흐름 | Rule 3 |
| **S5** | HMAC 서명 무결성 (조작 시 탐지) | Rule 3, 5 |
| **S6** | 보안 통계 + Rule 3, 5 최종 검증 | Rule 3, 5 |

각 테스트 파일 끝에 placeholder impl 포함 (독립 실행 가능).

---

## 구현 순서

1. `proof_validator.fl` 강화 (sha3 내장 → ZKP 검증 → hijack 탐지)
2. `domain_session.fl` 강화 (HMAC → nonce 강화 → hijack)
3. `tests/security_tests.fl` 신규 생성 (S1-S6, 600줄)
4. GOGS 커밋 + 푸시

---

## 최종 목표

```
총 테스트: 54개 (A1-H6 기존 48개 + S1-S6 신규 6개)
무관용 규칙: 8/8 (기존 6/8 + Rule 3, 5 추가 달성)
GOGS: https://gogs.dclub.kr/kim/freelang-sovereign-dns.git
최종 커밋: Phase 4 Security Layer 완성
```

## 규칙 최종 상태 (목표)

| # | 규칙 | 목표 | Phase 4 후 |
|---|------|------|-----------|
| 1 | 원격 해석 < 10ms | ✅ | 유지 |
| 2 | 캐시 히트 < 0.1ms | ✅ | 유지 |
| **3** | **ZKP 검증 < 1ms** | ⏳ | **✅ 달성** |
| 4 | 경로 결정 < 5ms | ✅ | 유지 |
| **5** | **하이재킹 방어 100%** | ⏳ | **✅ 달성** |
| 6 | ICANN 의존 0% | ✅ | 유지 |
| 7 | 장애 내성 50% | ✅ | 유지 |
| 8 | 레코드 불변성 CAS | ✅ | 유지 |
