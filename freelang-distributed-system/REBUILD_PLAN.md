# FreeLang 재구축 계획 (Rebuild Plan)
## 진정한 프로덕션 시스템 만들기

**시작**: 2026-03-02
**목표**: 정직하고 검증된 분산 시스템
**원칙**: "기록이 증명이다" - 실제 동작을 증명하자

---

## 📋 현황 분석

### 현재 문제점
```
1. Raft: 하드코드된 리더 선출 (노드 0만)
2. 카오스: 가짜 시뮬레이션
3. 코드: unwrap 45개, println 112개
4. 품질: 컴파일 불가능
5. 테스트: 의미 없는 자동 통과
6. 보안: 미검증
7. 성능: 로컬에서만 측정
```

### 실제 코드량
```
총 7,453줄 (23,015줄이 아님)
- supply_chain: 2,403줄
- tracing: 1,583줄
- chaos: 1,426줄
- security: 1,027줄
- bank/proxy/raft: 701줄
- main: 313줄
```

---

## 🎯 재구축 전략

### Phase 1: 기초 정비 (1주일)
**목표**: 컴파일 가능한 상태, 코드 품질 개선

#### 1-1. Raft 정리
- [ ] **현재 상태**: 하드코드된 리더 선출만 있음
- [ ] **할 일**:
  ```
  ❌ 제거: 하드코드된 노드 0
  ✅ 추가: 실제 election timeout 메커니즘
  ✅ 추가: 투표 기반 리더 선출
  ✅ 추가: Term 관리
  ✅ 추가: 최소 3개 노드 필수 검증
  ```
- [ ] **수정할 파일**: `src/raft/mod.rs`
- [ ] **예상 줄수**: 185줄 → 400줄
- [ ] **테스트**: 리더 변경 검증 (node 0이 아닌 다른 노드도)

#### 1-2. 코드 품질 개선
- [ ] **unwrap() 제거** (45개)
  ```rust
  // Before
  let account = bank.get_account("alice").await.unwrap();

  // After
  let account = match bank.get_account("alice").await {
      Ok(acc) => acc,
      Err(e) => {
          error!("Failed to get account: {}", e);
          return Err(e);
      }
  };
  ```

- [ ] **println! 제거** (112개)
  ```rust
  // Before
  println!("[CHAOS] Injected {}ms latency to node {}", delay_ms, node);

  // After
  tracing::info!("[CHAOS] Injected {}ms latency to node {}", delay_ms, node);
  ```

- [ ] **하드코드 제거** (47개)
  ```rust
  // Before
  const LEADER_TIMEOUT: u64 = 1000;
  let port = 8080;

  // After
  let config = Config::from_env()?;
  let leader_timeout = config.leader_timeout_ms;
  let port = config.server_port;
  ```

#### 1-3. 컴파일 가능화
- [ ] Cargo check 성공
- [ ] Cargo build --release 성공
- [ ] 모든 경고 제거

**소요 시간**: 3-4일
**테스트**: 컴파일, clippy linting

---

### Phase 2: Raft 정확한 구현 (1주일)
**목표**: RFC 5740 준수하는 Raft 구현

#### 2-1. Raft 핵심 구현
```
필수 구현:
✅ Leader Election
  - Randomized election timeout
  - Vote counting
  - Term management
  - Multi-node scenario

✅ Log Replication
  - AppendEntries RPC
  - Log matching property
  - Commit index tracking
  - Snapshot mechanism

✅ Safety
  - Vote restriction (only voted_for)
  - Log matching verification
  - Commit only current term entries

✅ Failure Recovery
  - Leader failure detection
  - Automatic failover
  - Log consistency restoration
```

#### 2-2. 실제 네트워크 시뮬레이션 추가
```rust
// 이전: DashMap에 지연 값만 저장
// 지금: 실제로 네트워크 요청 지연

pub struct RaftNetwork {
    messages: Arc<RwLock<VecDeque<Message>>>,
    latencies: Arc<DashMap<(u32, u32), u64>>,  // (from, to) -> delay_ms
}

// 메시지 전송 시
// 1. 지연 조회: latency_ms?
// 2. sleep: Duration::from_millis(latency_ms)
// 3. 메시지 전송
// 4. 응답 수신
```

#### 2-3. 테스트 케이스
```
✅ Normal operation (all healthy)
✅ Single leader failure → election
✅ Follower failure → continues
✅ Network partition → quorum determines leader
✅ Log consistency verification
✅ Snapshot & recovery
```

**소요 시간**: 4-5일
**참고**: etcd, Consul 코드 참고

---

### Phase 3: 카오스 엔지니어링 재구현 (5일)
**목표**: 실제 장애 시뮬레이션 및 검증

#### 3-1. 실제 카오스 주입
```
Before (가짜):
✗ 지연만 설정 (tokio::sleep)
✗ 자동으로 "감지됨" 표시
✗ 자동으로 "복구됨" 표시
✗ 결과 미리 정해짐

After (진짜):
✅ 네트워크 메시지 드롭
✅ 지연 추가 (실제 RPC 호출 시)
✅ 노드 프로세스 시뮬레이트 중단
✅ WAL 손상 시뮬레이션
✅ 실제 동작 검증 (읽기/쓰기 성공 여부)
```

#### 3-2. 신뢰할 수 있는 검증
```rust
pub async fn scenario_1_network_latency() -> ChaosTestResult {
    // Step 1: 초기 상태 스냅샷
    let initial_state = system.capture_state().await;
    assert_no_errors(&initial_state);

    // Step 2: 장애 주입
    let mut chaos = ChaosInjector::new();
    chaos.add_latency(node_1, node_2, 500).await;

    // Step 3: 실제 작업 수행 (자동 표시 아님)
    let write_result = system.write(
        "key1", "value1",
        Duration::from_secs(5)  // timeout
    ).await;

    // Step 4: 결과 검증 (실제 동작 확인)
    match write_result {
        Ok(_) => {
            println!("✅ Write succeeded despite latency");
            // 성공한 이유 분석: 다른 두 노드에서 quorum 달성?
        }
        Err(e) => {
            println!("❌ Write failed: {}", e);
            // 실패한 이유 분석: 타임아웃? 리더 없음?
        }
    }

    // Step 5: 데이터 일관성 검증
    let final_state = system.capture_state().await;
    let data_loss = verify_consistency(&initial_state, &final_state);

    ChaosTestResult {
        scenario: "Network Latency",
        passed: write_result.is_ok() && data_loss == 0,
        actual_behavior: format!("{:?}", write_result),  // 실제 동작
        metrics: ChaosMetrics {
            detection_time: /* 측정함 */,
            recovery_time: /* 측정함 */,
            data_loss: data_loss,
        }
    }
}
```

#### 3-3. 시나리오 개선
```
원래 (7개, 모두 가짜):
- Scenario 1-7: 모두 자동 통과

새로운 (5개, 진짜):
✅ S1: Network Latency
   - 검증: Write 성공 여부 실제 확인
   - 동작: Timeout 기반

✅ S2: Network Partition
   - 검증: Quorum 유지 확인
   - 동작: 소수 side는 읽기만 가능

✅ S3: Node Crash
   - 검증: 데이터 손실 없음
   - 동작: 다른 노드에서 복구

✅ S4: Multiple Failures (Concurrent)
   - 검증: 2개 노드 장애에서 생존
   - 동작: 3/5 quorum

✅ S5: Recovery & Consistency
   - 검증: 복구 후 모든 데이터 일치
   - 동작: 로그 동기화
```

**소요 시간**: 3-4일

---

### Phase 4: 보안 검증 (3일)
**목표**: 실제 동작 확인

#### 4-1. TLS 검증
```
Current: 코드만 있음
Todo:
✅ 실제 인증서 생성 (self-signed OK)
✅ TLS handshake 성공 확인
✅ 암호화된 통신 검증 (패킷 분석)
✅ 인증서 검증 로직 테스트
```

#### 4-2. JWT 완성
```
Current: 30% 구현
Todo:
✅ Token expiration 추가
✅ Token refresh mechanism
✅ Token revocation list (blacklist)
✅ Signature verification
✅ Claim validation
```

#### 4-3. 암호화 검증
```
Current: ChaCha20-Poly1305 코드만 있음
Todo:
✅ Key generation 및 저장
✅ Key rotation 메커니즘
✅ Plaintext 암호화 → 복호화 일치 확인
✅ AEAD 인증 실패 테스트
```

**소요 시간**: 2-3일

---

### Phase 5: 성능 벤치마킹 (3일)
**목표**: 실제 환경에서의 성능 측정

#### 5-1. 로컬 벤치마크 (기준)
```
Current: 1000 tx/sec (로컬)
Todo:
✅ Baseline 측정
✅ CPU 프로파일링
✅ 메모리 사용량 추적
✅ 지연 분포 (P50/P95/P99)
```

#### 5-2. 네트워크 환경 벤치마크
```
Setup:
- 3개 Docker 컨테이너 (별도 네트워크)
- loopback이 아닌 실제 TCP/IP

측정:
✅ 처리량 (TPS)
✅ 지연 (ms)
✅ 네트워크 대역폭
✅ 메모리 누수 (1시간 이상)
```

#### 5-3. 스트레스 테스트
```
Tests:
✅ 10K 동시 연결
✅ 100K TPS (또는 한계까지)
✅ 24시간 안정성 (soak test)
✅ 메모리 누수 확인
```

**소요 시간**: 2-3일

---

## 📅 전체 일정

```
Phase 1 (기초 정비):        3-4일   ✅ 컴파일 가능화
Phase 2 (Raft 정확화):      4-5일   ✅ RFC 5740 준수
Phase 3 (카오스 재구현):    3-4일   ✅ 실제 검증
Phase 4 (보안 검증):        2-3일   ✅ 동작 확인
Phase 5 (성능 벤치마크):    2-3일   ✅ 실제 측정
────────────────────────────────
총소요시간:                2-3주
```

---

## 🎯 성공 기준

### Phase 1 완료
```
✅ cargo build --release 성공
✅ cargo test 모두 통과
✅ Clippy warnings 0
✅ unwrap() 5개 이하 (테스트만)
✅ println! 10개 이하 (로깅만)
```

### Phase 2 완료
```
✅ 3개 이상 노드에서 리더 선출
✅ Node 0이 아닌 다른 노드도 리더 가능
✅ 리더 장애 → 새 리더 선출
✅ 로그 일관성 검증
```

### Phase 3 완료
```
✅ 각 시나리오에서 실제 동작 관찰
✅ 데이터 손실 = 0
✅ 복구 시간 < 5초
✅ 가짜 자동 통과 없음
```

### Phase 4 완료
```
✅ TLS handshake 성공
✅ JWT token 검증 성공
✅ 암호화/복호화 일치
```

### Phase 5 완료
```
✅ 로컬: 기준선 수립
✅ 네트워크: 실제 성능 측정
✅ 스트레스: 24시간 안정성 확인
```

---

## 📊 정직한 평가 기준

```
"프로덕션 준비됨"이라고 주장하기 전에:

✅ 코드가 컴파일되는가? (당연)
✅ 테스트가 통과하는가? (자동 통과 말고)
✅ 실제 작업을 한다고 검증했는가? (실제 RPC 호출)
✅ 장애에서 복구되는가? (가짜 아닌)
✅ 성능을 실제 환경에서 측정했는가? (loopback 아님)
✅ 24시간 이상 안정적인가? (soak test)
✅ 보안 감사를 받았는가? (적어도 내부)
```

---

## 🚫 더 이상 하지 않을 것

```
❌ "100% 통과"라고 말하기 (실제 의미 불명)
❌ "프로덕션 준비"라고 주장하기 (검증 전에)
❌ 자동 성공하는 테스트
❌ 과장된 수치
❌ 측정 없는 성능 주장
❌ 미검증 보안
❌ Unwrap과 println 혼용
```

---

## ✅ 새로운 원칙

```
"기록이 증명이다"

의미:
✅ 코드가 실제로 컴파일된다 (증명: 빌드 로그)
✅ 테스트가 실제로 통과한다 (증명: 테스트 실행)
✅ 기능이 실제로 동작한다 (증명: 작동 시나리오)
✅ 성능을 실제로 측정했다 (증명: 벤치마크 데이터)
✅ 장애에서 복구된다 (증명: 카오스 테스트 결과)

"말하지 말고, 보여주라"
Don't talk, show the results.
```

---

## 📌 다음 액션

### 즉시 (오늘)
1. [ ] 이 계획 검토 및 승인
2. [ ] GitHub Issue 생성 (각 Phase별)
3. [ ] 개발 시작

### Phase 1 시작 (내일)
```
우선순위:
1. Cargo.toml 확인 및 의존성 정리
2. src/raft/mod.rs 리더 선출 수정
3. 모든 unwrap() 제거
4. 모든 println! 제거 (tracing으로 교체)
5. cargo build 성공
6. cargo test 성공
```

---

## 📝 커밋 원칙

```
각 Phase 완료마다:

git commit -m "
chore(phase-X): 완료 - 구체적인 성과

실제 완료 항목:
✅ Feature A 구현
✅ Feature B 테스트
✅ Bug C 수정

검증 결과:
✅ 컴파일 성공
✅ 테스트 통과
✅ 동작 확인

성과 수치:
- 코드 추가: XXX줄
- 코드 수정: XXX줄
- 테스트 추가: X개
- 버그 수정: X개

"기록이 증명이다"
"Your record is your proof"

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
"
```

---

## 🎯 최종 목표

```
이 계획이 완료되면:

✅ 정직한 코드 (과장 없음)
✅ 검증된 기능 (자동 통과 없음)
✅ 측정된 성능 (로컬 환경만 아님)
✅ 실제 보안 (미검증 아님)
✅ 프로덕션 가능한 품질

그때야 비로소:
"프로덕션 준비됨"이라고 말할 수 있다.
```

---

**시작 준비 완료**
2026-03-02 | 다시 시작합시다
