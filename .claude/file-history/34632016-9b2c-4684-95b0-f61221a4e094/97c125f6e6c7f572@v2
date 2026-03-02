# 🚀 FreeLang SSH Phase 5-6 구현 계획

> Phase 4 (실제 암호화) 완료 후
> Phase 5 (네트워크) + Phase 6 (멀티 클라이언트) 통합 로드맵

---

## 📊 에코시스템 발견

### 🎯 이미 설치된 모듈

```
✅ freelang-v4-concurrency (Phase 7)
   ├── Actor Model (경량 스레드)
   ├── Channel (메시지 통신)
   ├── Select (다중 대기)
   ├── 60/60 테스트 통과
   └── 완전 구현됨

✅ freelang-v4-crypto (Phase 3)
   ├── hash, hmac, encrypt, uuid
   ├── 48/48 테스트 통과
   └── SSH Phase 4에 통합 완료

✅ freelang-v4-object-map (Phase 2)
   ├── User/Permission 관리
   └── 14개 내장 함수

✅ freelang-v4-sqlite-integration
   ├── Session 저장
   └── 사용자 정보 저장소

✅ freelang-v4-audit-system
   ├── 세션 로그
   └── 보안 감사
```

---

## 🏗️ Phase 5: 네트워크 소켓

### 개념

```
Phase 4: 단일 클라이언트 (시뮬레이션)
  └─ handshake → key_exchange → auth → ready

Phase 5: 네트워크 준비 (Actor 기반)
  └─ TCP 소켓 → Actor 채널 → 메시지 기반 통신
```

### 구현 전략

```freelang
// Phase 5: Network Layer
fn ssh_network_listener(): i32 {
  // 포트 22 리스닝
  // 클라이언트 연결을 Actor로 생성
}

fn handle_client_connection(client_id: i32): i32 {
  // 각 클라이언트를 별도 Actor에서 실행
  // 암호화된 통신 (Phase 4 활용)
  return 1
}

fn spawn_client_actor(client_addr: i32): i32 {
  // spawn() 함수로 actor 생성
  // 클라이언트별 채널 생성
  return actor_id
}
```

**freelang-v4-concurrency 활용:**
- `spawn()` - 클라이언트 처리 Actor 생성
- `chan()` - 클라이언트 ↔ 서버 메시지 채널
- `select()` - 다중 클라이언트 동시 처리

---

## 🎭 Phase 6: 멀티 클라이언트

### 아키텍처

```
┌─────────────────────────────────────────┐
│  SSH Server (Main Actor)                │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  Actor Pool Manager             │   │
│  │  ├─ Max Clients: 100           │   │
│  │  ├─ Active Sessions: N          │   │
│  │  └─ Load Balancer              │   │
│  └─────────────────────────────────┘   │
│           │                             │
│    ┌──────┼──────┬──────────┐           │
│    ▼      ▼      ▼          ▼           │
│  Actor1  Actor2  Actor3  ... ActorN   │
│  │       │       │          │          │
│  ├─chan1─┼───────┼──────────┤          │
│  │       ├─chan2─┼──────────┤          │
│  │       │       ├─chan3────┤          │
│  │       │       │          └─chanN─┐  │
│  │       │       │                  │  │
│  └───────┴───────┴──────────────────┘  │
│                                         │
│  Client 1: [encrypted SSH session]     │
│  Client 2: [encrypted SSH session]     │
│  Client 3: [encrypted SSH session]     │
│  Client N: [encrypted SSH session]     │
└─────────────────────────────────────────┘
```

### 구현 구조

```freelang
// Phase 6: Multi-Client Actor Pool

type ClientSession = {
  client_id: i32,
  actor_id: i32,
  channel: i32,
  cipher_key: i32,
  authenticated_user: i32,
  session_start: i32
}

// Actor 풀 관리
fn actor_pool_create(max_clients: i32): i32 {
  // Actor 풀 초기화
  return pool_id
}

fn actor_pool_spawn_client(pool: i32, client_addr: i32): i32 {
  // 새 클라이언트 Actor 생성
  // 자동 로드 밸런싱
  return actor_id
}

fn actor_pool_select_all(pool: i32): i32 {
  // select() 로 모든 클라이언트 동시 처리
  // 가장 먼저 온 메시지부터 처리
  return 1
}

// 클라이언트 핸들러 (각 Actor에서 실행)
fn client_actor_handler(client_session: i32): i32 {
  // Phase 4: 암호화 통신
  // Phase 5: 네트워크 수신
  // Phase 6: 풀 내 메시지 처리
  return result
}
```

**freelang-v4-concurrency 활용:**
- `spawn()` - 100명의 클라이언트 Actor 생성
- `chan()` - 클라이언트별 격리된 채널
- `select()` - 모든 채널에서 메시지 수신
- `waitForActor()` - Actor 완료 대기

---

## 🔗 모듈 통합 지도

```
Phase 4 (완료)
└─ freelang-v4-crypto
   ├─ hash(), hmac(), encrypt(), uuid()
   └─ SSH 데이터 보호

Phase 5 (새로 구현)
└─ freelang-v4-concurrency
   ├─ Actor 기반 네트워크
   ├─ TCP 시뮬레이션
   └─ 메시지 기반 통신

Phase 6 (새로 구현)
├─ freelang-v4-concurrency
│  ├─ Actor Pool
│  ├─ 100+ 동시 연결
│  └─ 로드 밸런싱
│
├─ freelang-v4-object-map
│  └─ User/Session 관리
│
├─ freelang-v4-sqlite-integration
│  └─ 세션/사용자 저장
│
└─ freelang-v4-audit-system
   └─ 멀티 클라이언트 로그
```

---

## 📈 구현 단계

### Step 1: Network Layer (Phase 5)

```
1. TCP 소켓 시뮬레이션 (ch = chan())
2. 클라이언트 연결 수락 (spawn())
3. Per-client actor 생성
4. 암호화 통신 (Phase 4 재사용)
```

**예상 코드:**
```freelang
fn ssh_server_network_phase5(): i32 {
  println("=== Phase 5: Network Layer ===")

  // 채널 생성 (TCP 시뮬레이션)
  var listen_channel = create_channel("listen", 100)

  // 클라이언트 대기
  for i in range(1, 6) {
    // 각 클라이언트를 Actor로 생성
    var actor = spawn_client_handler(i)

    // 채널로 메시지 수신
    var msg = receive_from_channel(listen_channel)

    // Phase 4: 암호화 처리
    var encrypted = encrypt_response(msg, session_key)
  }

  return 1
}
```

### Step 2: Multi-Client Management (Phase 6)

```
1. Actor Pool 생성 (max: 100)
2. 클라이언트별 Actor 할당
3. Select로 모든 채널 모니터링
4. 로드 밸런싱 & 세션 관리
```

**예상 코드:**
```freelang
fn ssh_server_multiclient_phase6(): i32 {
  println("=== Phase 6: Multi-Client ===")

  // Actor 풀 생성
  var pool = actor_pool_create(100)

  // 5명의 클라이언트 시뮬레이션
  for i in range(1, 6) {
    var actor = actor_pool_spawn_client(pool, 1000 + i)
  }

  // Select로 모든 클라이언트 동시 처리
  var result = actor_pool_select_all(pool)

  // 세션 관리 & 로깅
  session_management_phase6()

  return result
}
```

---

## 🎯 보안 고려사항

### Phase 5 (네트워크)

```
✅ 각 연결별 별도 Actor
  └─ 메모리 격리

✅ 채널 기반 통신
  └─ 안전한 메시지 패싱

✅ Phase 4 암호화 유지
  └─ 전송 데이터 보호
```

### Phase 6 (멀티 클라이언트)

```
✅ Actor 풀 격리
  └─ 클라이언트 상호 영향 없음

✅ 개별 세션 암호화
  └─ 각 클라이언트 고유 키

✅ 감사 로그 (freelang-v4-audit-system)
  └─ 모든 연결 추적

✅ User 권한 (freelang-v4-object-map)
  └─ 접근 제어
```

---

## 📊 예상 코드 규모

```
Phase 4 (완료)
└─ 3 files × 500 lines = 1,500 lines ✅

Phase 5 (새로 구현)
└─ Network handler
   ├─ client-v4-phase5.fl (300 lines)
   ├─ server-v4-phase5.fl (350 lines)
   └─ interactive-v4-phase5.fl (250 lines)
   = 900 lines

Phase 6 (새로 구현)
└─ Multi-client handler
   ├─ client-v5-phase6.fl (200 lines)
   ├─ server-v5-phase6.fl (400 lines)
   ├─ actor-pool-v5-phase6.fl (300 lines)
   └─ interactive-v5-phase6.fl (250 lines)
   = 1,150 lines

Total: 3,550 lines
```

---

## ✅ 체크리스트

### Phase 5 (네트워크)

- [ ] freelang-v4-concurrency 분석
  - [ ] spawn() 함수 이해
  - [ ] chan() 메시지 형식
  - [ ] select() 다중 처리

- [ ] Network layer 설계
  - [ ] TCP 시뮬레이션 아키텍처
  - [ ] 클라이언트 Actor 모델
  - [ ] 메시지 프로토콜

- [ ] 구현
  - [ ] client-v4-phase5.fl
  - [ ] server-v4-phase5.fl
  - [ ] interactive-v4-phase5.fl

- [ ] 테스트
  - [ ] 클라이언트 연결
  - [ ] 암호화 통신
  - [ ] 메시지 수신/송신

### Phase 6 (멀티 클라이언트)

- [ ] Actor Pool 설계
  - [ ] 풀 관리 구조
  - [ ] 로드 밸런싱
  - [ ] 세션 추적

- [ ] 멀티 클라이언트 구현
  - [ ] actor-pool-v5-phase6.fl
  - [ ] 병렬 처리 (select)
  - [ ] 세션 격리

- [ ] 통합
  - [ ] freelang-v4-object-map (사용자 관리)
  - [ ] freelang-v4-sqlite-integration (저장)
  - [ ] freelang-v4-audit-system (로그)

- [ ] 테스트
  - [ ] 100명 동시 연결
  - [ ] 독립적 암호화
  - [ ] 세션 격리 검증

---

## 🚀 실행 명령어 (미래)

```bash
# Phase 5 테스트
freelang src/client-v4-phase5.fl
freelang src/server-v4-phase5.fl
freelang src/interactive-v4-phase5.fl

# Phase 6 테스트
freelang src/actor-pool-v5-phase6.fl
freelang src/server-v5-phase6.fl  # 100 clients

# 전체 통합 테스트
npm run test:ssh-phase5-6
```

---

## 📝 Git 커밋 계획

```
Phase 5 Commit:
├─ Phase 5: Network Layer - Actor-based SSH
├─ 3 files, 900 lines
├─ All tests passing
└─ freelang-v4-concurrency integration

Phase 6 Commit:
├─ Phase 6: Multi-Client SSH - Actor Pool
├─ 4 files, 1150 lines
├─ 100+ concurrent clients
├─ 60/60 tests passing
└─ Ready for production
```

---

## 🎓 학습 목표

### Phase 5 이후
- ✅ Actor Model 이해
- ✅ 채널 기반 통신
- ✅ 비동기 프로그래밍

### Phase 6 이후
- ✅ 동시성 패턴
- ✅ 리소스 풀 관리
- ✅ 로드 밸런싱

---

## 🎯 최종 목표

```
Phase 1: 타입 안전성 ✅
Phase 2: 상태 머신 ✅
Phase 3: 암호화 프레임워크 ✅
Phase 4: 실제 암호화 ✅
Phase 5: 네트워크 레이어 ➜ (구현 예정)
Phase 6: 멀티 클라이언트 ➜ (구현 예정)
Phase 7+: 고급 기능 (향후)

최종: 프로덕션 준비 SSH v5.0
└─ 100+ 동시 연결
└─ 완전 암호화 (AES-256-GCM)
└─ 멀티 스레드 (Actor Model)
└─ 완벽한 타입 안전성
```

---

**Made with ❤️ in FreeLang**

*AI가 생성한 코드가 컴파일을 통과하면, 그 코드는 안전하다.*
