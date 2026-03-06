# Week 2 에이전트 배치 명령서

**시간**: 2026-03-14 09:00 KST
**목표**: 40,000줄 신규 코드 + 모든 에이전트 동시 실행

---

## 🚀 동시 실행 스케줄

모든 8개 에이전트는 **병렬로** 실행됩니다. (의존성 없음)

```
시작 시간: 2026-03-14 09:00
└─ [09:00] 모든 에이전트 동시 스폰
   ├─ Agent 1: v4 데이터 계층 (15,000줄) - 자세한 지침: agent-1-week2.md
   ├─ Agent 2: Sovereign-DNS (2,500줄)
   ├─ Agent 3: Sovereign-Mail (2,500줄)
   ├─ Agent 4: Phone & Quantum (4,000줄)
   ├─ Agent 5: Nano-Kernel + JIT (4,000줄)
   ├─ Agent 6: Monitoring & Security (3,500줄)
   ├─ Agent 7: Communications & Ledger (4,000줄)
   └─ Agent 8: Learning Lessons (2,500줄)

예상 완료: 2026-03-20 18:00 (약 80시간 소요)
```

---

## 📝 에이전트별 상세 지침 (요약)

### Agent 2: Sovereign-DNS - Challenge 17
**목표**: 2,500줄 (DNSSEC, DoT, 다중 마스터)

#### 작업 분배
```
Task 2.1: DNS 프로토콜 확장 (1,500줄)
├─ DNSSEC 구현 (Zone Signing, RRSIG, DS) - 800줄
├─ EDNS 지원 - 400줄
└─ DNS over TLS (DoT) - 300줄

Task 2.2: 분산 DNS 서버 (1,000줄)
├─ 다중 마스터 동기화 - 500줄
├─ 캐싱 개선 - 300줄
└─ 모니터링 - 200줄

테스트: 18개 (프로토콜 준수도 + 성능)
```

#### 구현 전략
1. **DNSSEC**:
   - Zone file에 RSA-2048 서명 추가
   - RRSIG 검증 로직
   - DS 레코드 연쇄 검증

2. **DoT**:
   - TLS 1.3 기반 DNS over TCP
   - 암호화된 쿼리/응답

3. **다중 마스터**:
   - Zone transfer (AXFR) 동기화
   - 충돌 감지 & 해결 (Last-Write-Wins)

---

### Agent 3: Sovereign-Mail - Challenge 15
**목표**: 2,500줄 (SMTP 확장, POP3, IMAP)

#### 작업 분배
```
Task 3.1: SMTP 확장 (1,200줄)
├─ STARTTLS - 300줄
├─ AUTH (PLAIN, CRAM-MD5) - 400줄
├─ 큐잉 & 재시도 - 300줄
└─ SPF/DKIM/DMARC 검증 - 200줄

Task 3.2: POP3/IMAP (1,300줄)
├─ POP3 프로토콜 구현 - 400줄
├─ IMAP4 프로토콜 구현 - 600줄
├─ 폴더 관리 - 200줄
└─ IDLE 지원 - 100줄

테스트: 22개
```

#### 핵심 프로토콜
- **SMTP AUTH**: 클라이언트 인증
- **POP3**: 간단한 메일 다운로드
- **IMAP4**: 원격 메일박스 관리
- **DKIM/DMARC**: 이메일 인증

---

### Agent 4: Sovereign-Phone - Phase 11 (Quantum Crypto)
**목표**: 4,000줄 (Post-Quantum 암호화)

#### 작업 분배
```
Task 4.1: PQC 구현 (2,500줄)
├─ CRYSTALS-Kyber (Lattice-based) - 700줄
├─ XMSS (Hash-based) - 600줄
├─ McEliece (Code-based) - 700줄
└─ NTRU (Polynomial-based) - 500줄

Task 4.2: 혼합 키 교환 (1,500줄)
├─ Classical + PQC 핸드셰이크 - 700줄
├─ 마이그레이션 전략 - 500줄
└─ 성능 벤치마크 - 300줄

테스트: 25개
```

#### 암호화 알고리즘
```
기존: ECDH (256-bit) = 128-bit 양자 안전
추가:
- Kyber512 (768 bytes)
- XMSS-SHA2_10_256
- McEliece 6960119
```

---

### Agent 5: Low-level Systems - Nano-Kernel + JIT
**목표**: 4,000줄 (Nano-Kernel 완성 + JIT 스켈레톤)

#### 작업 분배
```
Task 5.1: Nano-Kernel (2,000줄)
├─ 프로세스 스케줄러 (선점형) - 600줄
├─ 가상 메모리 & 페이징 - 600줄
├─ 인터럽트 핸들러 - 400줄
└─ 타이머 관리 - 400줄

Task 5.2: JIT 컴파일러 (2,000줄)
├─ 핫스팟 감지 - 500줄
├─ x86-64 기계어 생성 - 700줄
├─ 루프 언롤링 & 최적화 - 500줄
└─ 인라인 캐싱 - 300줄

테스트: 22개
```

#### 성능 목표
- 프로세스 스위칭: <100μs
- 페이지 폴트: <1ms
- JIT 컴파일: <10ms

---

### Agent 6: Monitoring & Security - Integrity Engine
**목표**: 3,500줄 (무결성 엔진 + Sentry)

#### 작업 분배
```
Task 6.1: 무결성 검증 (2,000줄)
├─ 파일 체크섬 (SHA256, BLAKE3) - 500줄
├─ 코드 무결성 확인 - 400줄
├─ 실시간 파일 감시 (fswatch) - 600줄
└─ 위조 감지 - 500줄

Task 6.2: Sentry (침입 감지) (1,500줄)
├─ 비정상 행동 탐지 - 600줄
├─ ML 기반 분류 - 500줄
├─ 알림 및 응답 - 250줄
└─ 인시던트 로깅 - 150줄

테스트: 25개
```

#### 감지 시나리오
- 파일 변조 감지
- 프로세스 비정상 행동
- 네트워크 이상 탐지
- 권한 에스컬레이션

---

### Agent 7: Communications & Data - Atomic Ledger
**목표**: 4,000줄 (분산 원장 + 스트리밍)

#### 작업 분배
```
Task 7.1: Atomic Ledger (2,000줄)
├─ 블록 구조 & 해싱 - 300줄
├─ 트랜잭션 검증 - 400줄
├─ BFT 합의 알고리즘 - 700줄
└─ 스마트 컨트랙트 기초 - 600줄

Task 7.2: 스트리밍 처리 (2,000줄)
├─ 이벤트 스트림 - 500줄
├─ 윈도우 집계 - 600줄
├─ 백프레셔 처리 - 500줄
└─ 정확히 한 번 처리 - 400줄

테스트: 24개
```

#### 분산 합의
- PBFT (Practical Byzantine Fault Tolerance)
- Raft 대체 검토
- 최대 f < n/3 노드 장애 허용

---

### Agent 8: Integration & Learning - Lessons 11-20
**목표**: 2,500줄 (중급 과정 10개 강의)

#### 작업 분배
```
Lessons 11-20 (1,500줄)
├─ Lesson 11: 모듈 시스템 - 150줄
├─ Lesson 12: 제너릭 & 타입 - 150줄
├─ Lesson 13: 에러 핸들링 - 150줄
├─ Lesson 14: 함수형 프로그래밍 - 200줄
├─ Lesson 15: 동시성 - 200줄
├─ Lesson 16: 비동기 처리 - 200줄
├─ Lesson 17: 메타프로그래밍 - 150줄
├─ Lesson 18: 성능 최적화 - 150줄
├─ Lesson 19: 테스팅 - 150줄
└─ Lesson 20: 프로젝트 구조 - 150줄

연습문제 (1,000줄)
├─ 총 30개 (3개 × 10개 강의)
├─ 상세 솔루션 코드
└─ 학습자 피드백 템플릿

테스트: 평가 문제 자동 채점
```

#### 학습 경로
1. 기초 (Lessons 1-10) ✅ Week 1 완료
2. 중급 (Lessons 11-20) 🚀 Week 2 진행 중
3. 고급 (Lessons 21-30) 📅 Week 3 예정

---

## 📊 리소스 할당

| 에이전트 | CPU | 메모리 | 스토리지 | 담당자 |
|---------|-----|--------|----------|--------|
| Agent 1 | 4x | 16GB | 500GB | Core Team |
| Agent 2 | 2x | 8GB | 200GB | Network Team |
| Agent 3 | 2x | 8GB | 200GB | Mail Team |
| Agent 4 | 2x | 8GB | 200GB | Security Team |
| Agent 5 | 2x | 8GB | 200GB | Systems Team |
| Agent 6 | 2x | 8GB | 200GB | Monitoring Team |
| Agent 7 | 2x | 8GB | 200GB | Data Team |
| Agent 8 | 1x | 4GB | 100GB | Education Team |

**총**: 17x CPU, 68GB RAM, 1.6TB SSD

---

## 🎯 성공 지표

### 각 에이전트별
- [x] 목표 코드량 80% 이상 달성
- [x] 무관용 테스트 15+ 개 통과
- [x] GOGS 저장소 커밋
- [x] 문서화 포함

### 전체
- [x] 총 38,000줄 신규 코드
- [x] 외부 의존도 0%
- [x] 누적 ~90,000줄 (Week 1 + 2)
- [x] 모든 에이전트 동시 완료

---

## 📅 일정

```
2026-03-14 09:00 → 2026-03-20 18:00 (약 80시간)

Daily Sync:
- 매일 09:00 KST: 진도 체크
- 매일 18:00 KST: 상태 리포팅
- 병목 현상 즉시 처리
```

---

## ⚠️ 주의사항

1. **병렬 실행**: 모든 에이전트는 독립적이므로 동시 실행 가능
2. **메모리 관리**: 개별 제약 준수
3. **GOGS 커밋**: 하루 1회 이상
4. **테스트 우선**: 코드 작성 전 테스트 작성
5. **문서화**: 매 500줄마다 주석 추가

---

**준비 완료** ✅
**명령**: 에이전트 스폰 대기 🚀

