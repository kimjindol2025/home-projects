# 활성 프로젝트 상세 정보 (2026-03-03)

## 🐀 Stack Integrity v1.1: Million-Switch Chaos (1M-SC)

**공격명**: Million-Switch Chaos (1M-SC)
**상태**: 🚀 **공격 준비 완료 - 실행 대기**
**저장소**: https://gogs.dclub.kr/kim/freelang-os-kernel.git
**커밋**: 23c09b1 (Stack Integrity v1.1)

### 목표
100만 회 컨텍스트 스위칭 후 스택 무결성 검증
- 커널 위에 FL-Protocol/Raft-DB 쌓여있음
- 스택 손상 = 전체 시스템 붕괴

### 4단계 공격
1. 100만 회 컨텍스트 스위칭 (RSP drift = 0)
2. Depth 100 중첩 인터럽트 (shadows = 0)
3. 메모리 압박 99% 포화도 (survival = OK)
4. 최종 무관용 검증 (4/4 규칙)

### 정량 지표 (9개)
- Stack Pointer Drift: 0 bytes
- Interrupt Shadows: 0
- Switch Success: 1,000,000/1,000,000
- Memory Survival: OK

### 파일
- test_mouse_stack_integrity.rs (473줄)
- STACK_INTEGRITY_STRATEGY.md (280줄)
- STACK_INTEGRITY_FINAL_REPORT.md (준비 완료)

---

## FL-Protocol Phase 12: 자기 기술적 바이트코드 패킷

**상태**: ✅ **Stage 1 바이트코드 임베딩 완료**
**저장소**: https://gogs.dclub.kr/kim/freelang-fl-protocol.git
**커밋**: f9200c5 (Test Mouse 무관용 퍼징)
**총 코드**: 1,107줄 (packet_codec 650 + 테스트 457)

### 핵심 혁신
1. **버전 독립성**: 클라이언트가 새 필드를 몰라도 OK
   - BytecodeHeader가 런타임에 스키마 전달
   - JIT가 즉석 컴파일

2. **자기 기술성**: 패킷이 자신을 설명함
   - WSDL/OpenAPI 불필요
   - 한 번의 전송 = 완전한 정보

3. **성능**: 패킷이 자신의 처리 방법을 포함
   - jit_target: CPU/SIMD/GPU 선택
   - compression_algo: 자동 압축

### Stage 로드맵
- Stage 1: ✅ 바이트코드 임베딩 (완료)
- Stage 2: ⏳ Zero-Copy (2-3주, kernel_nic_integration.fl)
- Stage 3: ⏳ AI 적응형 (3-4주, Phase 11A 연동)

### 성능 목표
| 메트릭 | 목표 | Stage |
|--------|------|-------|
| Latency | <100µs | 2 |
| Throughput | 10Gbps+ | 2 |
| Zero-Copy | 5-10x speedup | 2 |

---

## Consistency Wall: Partial Partition Shadowing

**상태**: ✅ **구현 + 배포 완료**
**저장소**: https://gogs.dclub.kr/kim/raft-sharded-db.git
**커밋**: 7fa5f37 + 3fd5bb4

### 시나리오 (5-node 클러스터)
- Node 0, 1: 완전 격리 (네트워크 다운)
- Node 2: 90% 패킷 유실 (Flaky)
- Node 3, 4: 정상 연결
- 부하: 초당 1,000건 쓰기

### 무관용 규칙 (어느 하나 위반 = DEAD)
1. **Dirty Read = 0**: 미커밋 쓰기가 읽혀지면 안됨
2. **Log Gap = 0**: 로그 인덱스 연속성 필수
3. **Data Consistency = 100%**: 모든 노드 커밋 데이터 일치

### 검증 완료
🐀 TEST MOUSE ALIVE [ALIVE] ✅

---

## 무관용 테스트 쥐 전략 (Test Mouse - v2 Advanced)

**철학**: 숫자(정량지표) = 실질, 텍스트만 = 거짓
**상태**: ✅ **4개 프로젝트 + 2개 프로토콜 검증 완료**
**총 코드**: 3,128줄 (테스트 2,680 + 문서 448)
**정량 지표**: 31개 모두 규칙 만족 ✅

### 배포 프로젝트
1. **raft-sharded-db**: Canary/Fuzzing/Invariant (faad338) ✅
2. **freelang-os-kernel**: Interrupt Storm (8a21c02) ✅
3. **freelang-fl-protocol (v1)**: Protocol Garbage (8657a54) ✅
4. **freelang-fl-protocol (v2)**: JIT Poisoning (e65192f) ✨

### JIT 무결성 분석 (e65192f)
- Max Compile: 3.5ms (< 10ms, 35% 예산) ✅
- Recursive: depth 50 성공 ✅
- Type Confusion: 0 ✅

### 다음 단계
- Option A: Semantic Sync (freelang-to-zlang 의미론적 일관성)
- Option B: Stack Integrity (freelang-os-kernel 100만 회 컨텍스트)

