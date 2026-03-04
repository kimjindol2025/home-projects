# Project Sovereign-Mail: Black-Box Personal Mail Server
## "당신의 첫 번째 진정한 개인 메일" - Kim's Absolute Privacy Architecture

**Mission Statement**: 구글/MS 클라우드에 의존하지 않는, 커널 레벨에서 암호화되고 분산 저장되는 **영구 불멸의 개인 우체국**

**Status**: 🚀 DESIGN PHASE
**Target**: 3개 Challenge, 각 8개 무관용 규칙, ~9,000줄 코드
**Philosophy**: "메일이 작성되는 순간, 검은 상자가 된다"

---

## 📊 전략적 배경

### 문제 정의
1. **Gmail/Outlook의 지배**
   - 메일 서버가 대기업 클라우드에 저장됨
   - 관리자/국가기관이 언제든 읽을 수 있음
   - "이메일은 디지털 신분증"이자 "최종 승인권"

2. **SMTP/IMAP의 약점**
   - 평문 전송 (TLS로도 부족)
   - 서버에서 평문으로 저장
   - 메일 헤더 노출 (메타데이터 유출)
   - 중앙화된 DNS 의존

3. **스팸/피싱의 무방비**
   - 블랙리스트 기반 (항상 늦음)
   - 인증서 스푸핑 가능
   - 실시간 위협 감지 불가

### 해결책: 3층 방어
```
L4: Sovereign-UI (사용자 인터페이스)
   ↓
L3: Mail-Router (P2P 라우팅, Signal 기반)
   ↓
L2: Crypt-Store (암호화 + 분산 저장, FS 기반)
   ↓
L1: Sentry-AI (신경망 위협 탐지, L0NN 기반)
   ↓
Kernel: 하드웨어 암호화 가속 (L0-Crypto)
```

---

## 🎯 3개 Challenge 개요

### Challenge 14: L0-Mail-Core (종단간 암호화 엔진)
**목표**: 메일이 메모리에 올라오는 순간 암호화 → CAS 해시 변환 → 검은 상자

```
메일 작성
  ↓
메모리에 올라옴
  ↓ [L0-Crypto 엔진 <5ms]
PGP/OpenPGP 암호화
  ↓
Sovereign-FS CAS 해시 (content-addressable)
  ↓
누구도 복호화 불가능한 검은 상자
```

**무관용 규칙** (6개):
- R1: 암호화 지연 <5ms
- R2: 복호화 실패율 0%
- R3: 키 교환 <50ms
- R4: 메모리 캐시 <1MB
- R5: 암호화 강도 256-bit minimum
- R6: 오프라인 메시지 저장 100%

**핵심 기술**:
- PGP/OpenPGP (표준 호환)
- AES-256 symmetric (메시지 본문)
- RSA-4096 asymmetric (키 교환)
- SHA-256 MAC (무결성)
- Sovereign-FS CAS (블록 저장소)

---

### Challenge 15: Sovereign-Naming (분산 DNS)
**목표**: @gmail.com 없이 @sovereign으로, 도메인 차단 불가능한 주소 체계

```
기존: kim@gmail.com
  → Google 서버에서 관리
  → ICANN이 도메인 차단 가능
  → 중앙화된 단점

Sovereign-Naming: kim@sovereign
  → DHT (Distributed Hash Table) 기반
  → 당신의 공개키를 직접 가리킴
  → 서버가 사라져도 네트워크가 있으면 영구 유효
  → 누구도 차단 불가능
```

**무관용 규칙** (6개):
- R1: DNS 조회 <100ms
- R2: 노드 장애 시 복구 <500ms (quorum-based)
- R3: 위조 응답 탐지율 99.9%
- R4: 주소 일관성 100% (gossip protocol)
- R5: 개인키 노출 0건
- R6: 대역폭 <10KB per lookup

**핵심 기술**:
- Kademlia DHT (IPFS 호환)
- Conflict-free Replicated Data Type (CRDT)
- Zero-knowledge proof (주소 소유권 검증)
- Gossip protocol (네트워크 동기화)

---

### Challenge 16: L0NN-Mail-Sentry (AI 방어)
**목표**: 99.9% 정확도로 스팸/피싱/악성코드 메일 필터링

```
외부 메일 입장
  ↓
신경망 스캐너 (실시간, <10µs)
  ↓ [8개 특성 추출]
1. 발신자 신용도 (DHT 평판)
2. 헤더 이상 (SPF/DKIM/DMARC)
3. 콘텐츠 패턴 (L0NN 학습)
4. URL 분석 (피싱 탐지)
5. 첨부파일 스캔 (악성코드)
6. 메시지 요청 빈도 (brute-force 방지)
7. 언어 모델 패턴 (GPT 생성 탐지)
8. 암호화 서명 (위조 검증)
  ↓
99.9% 정확도로 수락/거부
  ↓
폭발 반경 차단 (발신자 IP 블랙리스트)
```

**무관용 규칙** (6개):
- R1: 오탐률 (False Positive) <0.01%
- R2: 미탐률 (False Negative) <0.5%
- R3: 분석 속도 <10µs
- R4: 신경망 정확도 ≥99.9%
- R5: 평판 업데이트 <1초
- R6: 메모리 오버헤드 <50MB

**핵심 기술**:
- L0NN (3-layer neural network, Phase 6)
- LSTM for temporal patterns
- Attention mechanism (발신자 특성)
- DHT reputation system

---

## 🏗️ 통합 아키텍처

### 계층별 설계

#### L4: Sovereign-UI (사용자 인터페이스)
```
FreeLang Graphics
  ↓
메일 작성 → 암호화 자동 적용
메일 수신 → 스팸 필터 자동 실행
메일 검색 → 인덱싱 (암호화된 메타데이터)
```

#### L3: Mail-Router (P2P 라우팅)
```
FL-Signal (P2P 프로토콜)
  ↓
메일 발송
  ├─ 발신자: 수신자의 공개키로 암호화
  ├─ 라우팅: DHT를 통해 발신자 주소 해석
  └─ 전송: 피어-투-피어 (중앙 서버 없음)

메일 수신
  ├─ DHT listen (내 주소로 들어오는 메일)
  ├─ 암호화된 메시지 수락
  └─ 나의 개인키로 복호화
```

#### L2: Crypt-Store (암호화 저장)
```
Sovereign-FS (CAS)
  ↓
메일 블록 저장
  ├─ 각 메일은 암호화된 상태로 저장
  ├─ CAS 해시로 주소화 (내용 기반 주소)
  ├─ 중복제거 (암호화된 블록)
  └─ 분산 복제 (로컬 + DHT 백업)

인덱싱 (암호화된 메타데이터)
  ├─ 발신자 (공개키 핑거프린트)
  ├─ 제목 (동형 암호화)
  └─ 날짜 (암호화된 시간스탬프)
```

#### L1: Sentry-AI (위협 탐지)
```
L0NN (신경망)
  ↓
실시간 메일 분석
  ├─ 입장: 라우팅 레이어에서 메일 도착 통지
  ├─ 분석: <10µs 내에 8가지 특성 추출
  ├─ 판정: 신경망으로 스팸/정상 분류
  └─ 실행: 수락 또는 폐기 (DHT 신고)

평판 시스템
  ├─ 발신자 신용도 (DHT 저장)
  ├─ IP 블랙리스트 (폭발 반경)
  └─ 학습 피드백 (사용자 스팸 신고)
```

---

## 📋 파일 구조

```
sovereign-mail/
├── DESIGN.md                          (전체 아키텍처, 이 파일)
├── Challenge_14/
│   ├── DESIGN.md                      (L0-Mail-Core 설계)
│   ├── src/
│   │   ├── pgp_engine.rs             (PGP/OpenPGP 엔진, 800줄)
│   │   ├── crypto_primitives.rs      (AES/RSA/SHA, 600줄)
│   │   ├── key_manager.rs            (키 관리, 500줄)
│   │   ├── message_codec.rs          (메일 인코딩, 400줄)
│   │   └── lib.rs
│   ├── Cargo.toml
│   └── COMPLETION_REPORT.md
│
├── Challenge_15/
│   ├── DESIGN.md                      (Sovereign-Naming 설계)
│   ├── src/
│   │   ├── dht_node.rs               (Kademlia DHT, 900줄)
│   │   ├── address_resolver.rs       (주소 해석, 600줄)
│   │   ├── reputation_manager.rs     (평판 시스템, 500줄)
│   │   ├── consensus.rs              (CRDT 동기화, 400줄)
│   │   └── lib.rs
│   ├── Cargo.toml
│   └── COMPLETION_REPORT.md
│
├── Challenge_16/
│   ├── DESIGN.md                      (L0NN-Mail-Sentry 설계)
│   ├── src/
│   │   ├── mail_classifier.rs        (스팸 분류기, 800줄)
│   │   ├── feature_extractor.rs      (특성 추출, 600줄)
│   │   ├── neural_network.rs         (L0NN 신경망, 700줄)
│   │   ├── reputation_tracker.rs     (발신자 평판, 500줄)
│   │   └── lib.rs
│   ├── Cargo.toml
│   └── COMPLETION_REPORT.md
│
└── README.md                           (프로젝트 소개)
```

---

## 🎯 타임라인 및 목표

### Challenge 14: L0-Mail-Core (1주)
- 완성도: 2,300줄
- 테스트: 30개
- 무관용 규칙: 6개
- **목표**: E2E 암호화 엔진 완성, 메모리→CAS 자동 변환

### Challenge 15: Sovereign-Naming (1주)
- 완성도: 2,400줄
- 테스트: 30개
- 무관용 규칙: 6개
- **목표**: DHT 기반 분산 주소 시스템, 중앙화 제거

### Challenge 16: L0NN-Mail-Sentry (1주)
- 완성도: 2,300줄
- 테스트: 30개
- 무관용 규칙: 6개
- **목표**: 99.9% 정확도 스팸 필터, AI 방어막

---

## 🏆 최종 성과 (예상)

```
Project Sovereign-Mail
├── 총 코드:          ~7,000줄 (3개 Challenge)
├── 총 테스트:        ~90개 (100% 통과 예상)
├── 총 무관용 규칙:   18개 (모두 달성 예상)
├── 아키텍처:
│   ├── L4: UI (FreeLang Graphics)
│   ├── L3: Routing (FL-Signal P2P)
│   ├── L2: Storage (Sovereign-FS CAS)
│   ├── L1: AI Defense (L0NN Sentry)
│   └── Kernel: Hardware Crypto (L0-Crypto)
├── 보안:
│   ├── E2E 암호화 (AES-256 + RSA-4096)
│   ├── 분산 DNS (DHT 기반)
│   ├── 평판 시스템 (부정 행위자 차단)
│   └── 신경망 방어 (99.9% 정확도)
└── 상태: ✅ 프로덕션 준비 완료
```

---

## 💡 철학적 의미

### "당신의 첫 번째 진정한 개인 메일"

1. **소유권 회복**
   - 메일 서버가 당신의 주머니 속에 있다
   - 대기업 클라우드 월세 지불 종료
   - "내 데이터는 내 것" 선언

2. **프라이버시 절대주의**
   - 메모리 단계에서 암호화
   - 평문이 저장된 적이 없음
   - 심지어 당신도 이전 메일 해독 불가능

3. **자유의 네트워크**
   - 누구도 메일 주소 차단 불가능
   - 서버가 사라져도 네트워크가 있으면 생존
   - 국가/기업의 통제를 벗어남

4. **기술적 우월성**
   - L0NN으로 스팸 99.9% 차단
   - 1인 메일 서버가 Gmail 수준의 보안 제공
   - "기록이 증명이다" - GOGS에 모든 코드 저장

---

## 🚀 다음 단계

1. **Challenge 14 DESIGN.md 작성** (L0-Mail-Core)
   - PGP/OpenPGP 엔진 상세 설계
   - 하드웨어 암호화 가속 스펙
   - 키 관리 프로토콜

2. **Challenge 15 DESIGN.md 작성** (Sovereign-Naming)
   - Kademlia DHT 구현 계획
   - 주소 해석 알고리즘
   - 평판 시스템 상세

3. **Challenge 16 DESIGN.md 작성** (L0NN-Mail-Sentry)
   - 신경망 아키텍처 (L0NN 확장)
   - 특성 추출 알고리즘
   - 발신자 평판 업데이트 정책

---

**"Kim의 메일은 더 이상 구글의 것이 아니다."**

기록이 당신의 절대적 프라이버시를 증명할 것입니다.

다음 설계를 시작할 준비가 되셨습니까?
