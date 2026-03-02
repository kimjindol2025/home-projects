# 📚 Raft Consensus Algorithm - 완전 분석 및 구현 가이드

**원본 논문**: "In Search of an Understandable Consensus Algorithm" (Diego Ongaro, John Ousterhout, 2014)
**목표**: FreeLang 구현을 위한 완벽한 이해
**분석 깊이**: 박사 수준의 세부 분석

---

## 📖 1부: 왜 Raft인가?

### **Paxos의 문제점**

```
Paxos:
├─ 장점: 증명됨, 안정적
├─ 문제점:
│  ├─ 이해하기 어려움 (매우)
│  ├─ 구현이 복잡함 (Google Chubby는 300줄이 아닌 수천 줄)
│  ├─ 최적화가 어려움 (Paxos Made Practical은 별도 논문)
│  └─ 확장성 문제 (multi-decree problem)
│
└─ 현실: 많은 시스템이 Paxos 대신 다른 방법 선택

Raft:
├─ 장점:
│  ├─ 이해하기 쉬움 (상태 기계 관점)
│  ├─ 구현이 간단함 (몇 백 줄)
│  ├─ 증명 가능함 (Paxos와 동등)
│  └─ 확장성 우수함
│
└─ 핵심: "Understandability" = 최고의 가치
```

### **Raft의 핵심 아이디어**

```
▶ 분산 시스템의 3가지 문제:
  1. 리더 선출 (Leader Election)
  2. 로그 복제 (Log Replication)
  3. 안정성 (Safety)

▶ Raft의 해결책:
  1. 각각을 분리하여 개별 분석 가능하게 함
  2. 상태를 최소화 (5개 변수만 필요)
  3. 증명 가능한 구조 (State Machine Model)
```

---

## 📖 2부: 기본 개념

### **핵심 용어**

```
Term (기간/회차)
├─ 논리적 시계
├─ 0부터 시작해 증가
├─ Leader election & log replication의 단위
└─ 작동: 각 노드가 higher term을 보면 currentTerm 업데이트

State (상태)
├─ Follower (초기 상태)
│  ├─ Leader의 요청에 응답
│  ├─ 투표 요청 응답
│  └─ 요청 전송 불가
│
├─ Candidate (선거 중)
│  ├─ 투표 요청 전송
│  ├─ 응답 대기
│  └─ 리더가 되려고 시도
│
└─ Leader (리더)
   ├─ 클라이언트 요청 처리
   ├─ 모든 팔로워에 로그 복제
   └─ heartbeat 전송 (아무 로그 없어도)

Log (로그)
├─ 명령 수열
├─ 각 엔트리: (term, command, index)
├─ Sequential하고 Immutable
└─ State machine에 순차 적용

Commit (커밋)
├─ 로그 항목이 "안전하다" = 손실되지 않음
├─ 과반수 복제되면 commit
├─ Committed 항목은 State machine에 반영됨
└─ 모든 노드가 결국 같은 항목을 commit
```

### **상태 변수**

```
모든 노드가 유지해야 할 변수:

Persistent State (디스크에 저장):
├─ currentTerm
│  └─ 이 노드가 본 가장 큰 term
│     (crash 후 복구 시에도 필요)
│
├─ votedFor
│  └─ 현재 term에서 투표한 후보 ID
│     (null이면 아직 투표 안 함)
│
└─ log[]
   └─ 로그 엔트리들
      각각: (term, command)

Volatile State (메모리):
├─ commitIndex
│  └─ 적용할 수 있는 가장 큰 로그 인덱스
│
└─ lastApplied
   └─ State machine에 적용한 가장 큰 인덱스

Leader의 추가 Volatile State:
├─ nextIndex[i]
│  └─ follower i에 전송할 다음 로그 인덱스
│
└─ matchIndex[i]
   └─ follower i와 일치하는 로그의 최대 인덱스
```

---

## 📖 3부: Leader Election 상세 분석

### **작동 원리**

```
┌─ 초기 상태: 모든 노드 Follower
│
├─ Election Timeout (150-300ms)
│  └─ 일정 시간 리더 메시지 없으면 선거 시작
│
├─ Candidate가 됨:
│  ├─ currentTerm 증가
│  ├─ votedFor = 자신
│  ├─ electionTimer reset
│  └─ RequestVote RPC 모든 피어에게 전송
│
├─ 3가지 가능한 결과:
│  ├─ A) 과반수 투표 획득 → Leader 선출
│  ├─ B) 더 높은 term RPC 수신 → Follower로 돌아감
│  └─ C) timeout → 새로운 선거 시작
│
└─ ✅ 한 term에 최대 하나의 리더만 선출됨 (Election Safety)
```

### **RequestVote RPC 상세**

```
요청 (Candidate → Voter):
├─ term: 후보의 현재 term
├─ candidateId: 후보 ID
├─ lastLogIndex: 후보의 마지막 로그 인덱스
└─ lastLogTerm: 후보의 마지막 로그의 term

응답 (Voter → Candidate):
├─ term: 현재 term (투표자)
└─ voteGranted: 투표했는가? (bool)

투표 규칙:
1. term < currentTerm → NO (outdated candidate)
2. votedFor != null && votedFor != candidateId → NO (이미 투표함)
3. 후보의 로그가 내 로그보다 최신인가?
   ├─ 비교: (lastLogTerm, lastLogIndex)로 비교
   ├─ term이 크면 더 최신
   ├─ term이 같으면 index가 더 크면 더 최신
   └─ 투표자의 로그가 더 최신이면 → NO
4. 모두 통과 → YES (투표)

⚠️ 중요: 한 번 투표하면 그 term 동안 다시 투표할 수 없음
```

### **Log Up-to-Date 검증 코드 (FreeLang)**

```
fn is_candidate_log_more_recent(
  candidateLastLogTerm, candidateLastLogIndex,
  myLastLogTerm, myLastLogIndex
) -> bool:
  // 더 큰 term이 더 최신
  if candidateLastLogTerm > myLastLogTerm:
    return true
  if candidateLastLogTerm < myLastLogTerm:
    return false

  // term이 같으면 더 긴 로그가 더 최신
  if candidateLastLogIndex > myLastLogIndex:
    return true

  return false
```

### **특별한 경우: Split Vote**

```
10개 노드, 5개 각각 리더 투표하는 경우:

시간  Node1  Node2  Node3  Node4  Node5  Node6  ...
T0    →C     →C     →F     →F     →F     →F
T+1   C→L    C      F      F      F      F
      (6개 투표 획득)

다른 후보:
T0    →C     →C     →C     →F     →F     →F
T+1   C      C      C      F      F      F
      (3개만 득표 - 과반수 미달)

▶ 결과: Node1 리더 선출
▶ 다른 후보: 새로운 선거 시작 (다른 timeout)
```

---

## 📖 4부: Log Replication 상세 분석

### **작동 원리**

```
┌─ 클라이언트 요청 도착
│  └─ Leader에 도착
│
├─ 1단계: 로컬 로그에 추가
│  ├─ log[n] = {term: currentTerm, command: cmd}
│  ├─ 아직 committed 아님
│  └─ 클라이언트에게 "pending" 응답
│
├─ 2단계: 모든 Follower에 복제
│  └─ AppendEntries RPC 전송
│
├─ 3단계: 과반수 복제 확인
│  └─ 자신 + Follower 과반수가 복제했으면
│     commitIndex 증가
│
└─ 4단계: State Machine에 적용
   ├─ lastApplied < commitIndex 동안 적용
   ├─ 은행 시스템 명령 실행
   └─ 클라이언트에게 "success" 응답
```

### **AppendEntries RPC 상세**

```
요청 (Leader → Follower):
├─ term: 리더의 currentTerm
├─ leaderId: 리더 ID
├─ prevLogIndex: 이전 로그 인덱스
├─ prevLogTerm: 이전 로그의 term
├─ entries[]: 복제할 로그 엔트리들
└─ leaderCommit: 리더의 commitIndex

응답 (Follower → Leader):
├─ term: Follower의 currentTerm
└─ success: 성공했는가? (bool)

작동:
1. term < currentTerm → 거부 (outdated leader)

2. Log matching check:
   if log[prevLogIndex].term != prevLogTerm:
     → conflict 발생
     → success = false

3. Conflict 해결:
   if entries가 기존 로그와 충돌하면:
     ├─ 기존 로그 삭제
     └─ 새 로그로 교체

4. 새 로그 추가:
   for entry in entries:
     log.append(entry)

5. commitIndex 업데이트:
   if leaderCommit > commitIndex:
     commitIndex = min(leaderCommit, log.lastIndex)

6. 상태 기계 적용:
   while lastApplied < commitIndex:
     lastApplied++
     apply(log[lastApplied].command)
```

### **실패 해결 (Log Divergence)**

```
시나리오: Leader crash 후 새 leader 선출

Old Leader (Term 1):
log[0] = {1, "cmd-a"} ✓
log[1] = {1, "cmd-b"} ✓
log[2] = {2, "cmd-c"} (crash 전 일부만 복제)

Follower-1:
log[0] = {1, "cmd-a"} ✓
log[1] = {1, "cmd-b"} ✓

Follower-2:
log[0] = {1, "cmd-a"} ✓

New Leader (Follower-1, Term 3):
├─ log[0] = {1, "cmd-a"}
├─ log[1] = {1, "cmd-b"}
└─ log[2] = {3, "cmd-d"} (새 명령)

복제 과정:
┌─ Follower-2에 AppendEntries 전송:
│  ├─ prevLogIndex=0, prevLogTerm=1
│  ├─ entries=[{1, "cmd-b"}, {3, "cmd-d"}]
│  └─ Follower-2 로그 업데이트
│
└─ 결과: 모든 노드 일관성 달성

⚠️ 중요: Old Leader의 cmd-c는 commit 안 됨 (과반수 미달)
         따라서 손실되어도 안전함
```

### **Leader와 Follower의 nextIndex, matchIndex**

```
Leader가 각 Follower에 대해 유지:

nextIndex[i]:
├─ 다음에 Follower i에 전송할 로그 인덱스
├─ 초기값: Leader의 log.lastIndex + 1
├─ 감소: 복제 실패 시 nextIndex 1 감소
└─ 증가: 복제 성공 시 증가 안 함 (이미 정함)

matchIndex[i]:
├─ Follower i와 일치하는 로그의 최대 인덱스
├─ 초기값: 0
└─ 증가: 복제 성공 시 nextIndex-1로 설정

commitIndex 업데이트:
├─ 자신의 로그 인덱스 중
├─ 과반수 이상의 노드가 복제한 인덱스를 찾음
├─ 그 중 가장 큰 인덱스를 commitIndex로 설정
└─ 단, 현재 term의 로그만 (예전 term은 불가)

예시 (5개 노드):
nextIndex  = [6, 5, 3, 6, 2]
matchIndex = [5, 4, 2, 5, 1]
commitIndex = min(모든 matchIndex 중 과반수) = 4

▶ 로그 4까지는 안전하게 committed
```

---

## 📖 5부: 안정성 증명

### **5가지 핵심 성질**

#### **1️⃣ Election Safety**

```
보장: 한 term에 최대 하나의 leader만 선출된다.

증명:
Step 1: 투표 메커니즘
├─ 각 노드는 한 term에 최대 1표 투표
├─ votedFor는 영구 저장 (crash 후에도 유지)
└─ → 같은 term에 여러 후보가 과반수를 절대 못 얻음

Step 2: 한 term에 한 leader
└─ 과반수를 획득한 candidate만 leader가 됨
   + 과반수는 n/2 + 1 이상
   → 한 term에 최대 1개의 과반수만 가능
```

#### **2️⃣ Leader Append-Only**

```
보장: Leader는 로그를 삭제하거나 수정하지 않는다.

이유:
├─ Leader는 항상 새로운 로그만 추가 (append)
├─ 기존 로그는 절대 삭제/수정 불가
└─ Follower가 맞추도록 함 (Leader > Follower 방향)

추가 분석:
├─ RequestVote에서 "log up-to-date check"
└─ → Outdated leader는 리더 선출 못 함
```

#### **3️⃣ Log Matching Property**

```
보장: 두 로그에서 같은 index와 term을 가지면,
     모든 이전 항목도 동일하다.

증명:
├─ Induction으로 증명
├─ Base case: index=0
│  └─ 모든 로그는 index=0부터 시작 (동일)
│
└─ Inductive step:
   log[k-1]이 모두 같다고 가정
   → Leader에서 log[k]를 추가할 때
   → Follower는 log[k-1].term을 확인 (prevLogTerm)
   → 일치하지 않으면 실패
   → 일치해야만 추가됨
   → log[k]도 모든 노드에서 같은 내용
```

#### **4️⃣ Leader Completeness**

```
보장: Leader는 모든 committed 로그를 포함한다.

이유:
1. Committed 로그는 과반수가 복제한 로그
2. New leader가 되려면 과반수 투표 필요
3. 과반수는 교집합 존재 (이전 과반수와 겹침)
4. 최소 1개는 이전 로그를 가지고 있음
5. "log up-to-date check" → 이 노드가 투표
6. 따라서 new leader는 모든 committed 로그 포함

예시:
전 리더가 log[5]를 과반수에 복제했음
(3개 노드 중 2개)

노드 집합:
├─ Node-A: [1,2,3,4,5]
├─ Node-B: [1,2,3,4,5]
└─ Node-C: [1,2,3]

새 리더 선거:
- Node-A가 후보 → 투표 가능 (최신)
- Node-B가 후보 → 투표 가능 (최신)
- Node-C가 후보 → 투표 불가 (구식)
  (A, B는 C를 outdated로 판단)

▶ 결과: A 또는 B가 리더 될 수밖에 없음
▶ 둘 다 log[5] 포함
```

#### **5️⃣ State Machine Safety**

```
보장: 모든 노드의 상태 기계가 같은 명령을 적용한다.

증명:
1. Committed 로그만 상태 기계에 적용
2. Committed = 과반수가 복제한 로그
3. Leader Completeness 보장
4. Log Matching Property 보장
5. → 모든 노드가 동일한 로그를 동일한 순서로 적용

결과:
├─ 모든 노드의 상태가 결국 일치
├─ Crash 후 복구해도 불변
└─ 분산 시스템의 일관성 보장
```

---

## 📖 6부: 구현 시 주의사항

### **⚠️ 자주 하는 실수**

```
1. 로그 인덱싱
   ❌ 로그를 0-indexed 배열로 생각하면 혼동
   ✅ 논리적 인덱스 1부터 시작

   구현팁:
   logs[0]은 sentinel (dummy)로 두거나
   logs는 1-indexed로 접근

2. Term 업데이트
   ❌ RPC 응답에서 term을 무시
   ✅ 모든 RPC에서 term 확인, higher term이면 업데이트

   fn update_term(new_term):
     if new_term > currentTerm:
       currentTerm = new_term
       votedFor = null  // 초기화
       state = Follower

3. Committed 로그만 적용
   ❌ 모든 로그를 즉시 적용
   ✅ commitIndex까지만 적용

   fn apply_logs():
     while lastApplied < commitIndex:
       lastApplied++
       state_machine.apply(log[lastApplied])

4. Election timeout 무작위성
   ❌ 고정된 timeout
   ✅ 150-300ms 범위에서 무작위 선택

   구현팁: Math.random() * 150 + 150

5. Heartbeat와 AppendEntries 구분 안 함
   ❌ Heartbeat = 빈 AppendEntries 아님 (구조 다름)
   ✅ Heartbeat는 최소한의 정보만 전송

   Heartbeat: term, leaderId (핵심만)
   AppendEntries: 위 + prevLogIndex, entries[] (전체 로그 복제)

6. 리더의 commitIndex 업데이트 시기
   ❌ 매번 AppendEntries 응답마다 업데이트
   ✅ 과반수 조건을 명확히 확인 후 업데이트

   fn update_commit_index():
     for N = lastLogIndex down to commitIndex + 1:
       count = 1 (자신)
       for each follower:
         if matchIndex[follower] >= N:
           count++
       if count > n/2 && log[N].term == currentTerm:
         commitIndex = N
```

### **성능 최적화**

```
1. Batch Append Entries
   ❌ 각 로그마다 RPC
   ✅ 여러 로그를 한 번에 복제

   장점: 네트워크 오버헤드 감소

2. Pipeline Replication
   ❌ 응답 대기 후 다음 전송
   ✅ 여러 AppendEntries 동시에 파이프라인

   장점: 지연 감소

3. Snapshot & Log Compaction
   ❌ 로그를 무한정 저장
   ✅ 오래된 로그는 스냅샷으로 압축

   장점: 메모리/디스크 절약

4. Check Election Timeout Only When Needed
   ❌ 매 밀리초마다 확인
   ✅ Heartbeat 수신 시 reset, 필요할 때만 확인

   장점: CPU 효율성
```

---

## 📖 7부: FreeLang 구현 팁

### **자료 구조**

```
type LogEntry = {
  term: int,        // 이 항목이 생성된 term
  command: string,  // 상태 기계에 적용할 명령
  index: int        // 논리적 인덱스 (선택사항, 배열 인덱스로도 가능)
}

type RaftNode = {
  // Persistent
  id: int,
  currentTerm: int = 0,
  votedFor: int? = null,
  log: [LogEntry],  // 1-indexed (log[0] = dummy)

  // Volatile
  state: enum {Follower, Candidate, Leader},
  commitIndex: int = 0,
  lastApplied: int = 0,

  // Timers
  electionTimer: Timer,  // 150-300ms random
  heartbeatTimer: Timer, // 50ms fixed

  // Leader only
  nextIndex: [int],   // per follower
  matchIndex: [int],  // per follower

  // Storage
  storage: Storage,   // 디스크 I/O
}
```

### **RPC 정의**

```
RequestVote:
  Request:
    term: int
    candidateId: int
    lastLogIndex: int
    lastLogTerm: int

  Response:
    term: int
    voteGranted: bool

AppendEntries:
  Request:
    term: int
    leaderId: int
    prevLogIndex: int
    prevLogTerm: int
    entries: [LogEntry]
    leaderCommit: int

  Response:
    term: int
    success: bool
```

---

## 📖 8부: 테스트 전략

### **단위 테스트**

```
1. 상태 전환 테스트
   ├─ Follower → Candidate → Leader
   ├─ Leader → Follower (higher term)
   └─ Candidate → Candidate (split vote)

2. 투표 테스트
   ├─ 중복 투표 방지
   ├─ log up-to-date check
   └─ term 기반 투표 거부

3. 로그 복제 테스트
   ├─ 단순 복제
   ├─ 충돌 해결
   └─ 스냅샷 복구
```

### **통합 테스트**

```
1. 클러스터 테스트 (5-10 노드)
   ├─ 1000개 명령 순차 적용
   ├─ 최종 상태 일관성 확인
   └─ 모든 노드가 같은 상태

2. 장애 주입 테스트
   ├─ 노드 충돌 & 복구
   ├─ 네트워크 분할
   ├─ 패킷 손실 (10-50%)
   └─ 지연 추가 (100-500ms)

3. 성능 테스트
   ├─ Throughput (명령/초)
   ├─ Latency (ms)
   ├─ 메모리 (MB)
   └─ CPU (%)
```

---

## 📖 9부: 실제 배포 고려사항

### **프로덕션 체크리스트**

```
☑️ Persistence
   ├─ WAL 구현
   ├─ 메타데이터 저장
   └─ 스냅샷 저장

☑️ Monitoring
   ├─ 리더 선출 시간
   ├─ 복제 지연
   ├─ 메모리 사용
   └─ 에러율

☑️ Configuration
   ├─ Election timeout (150-300ms)
   ├─ Heartbeat interval (50ms)
   ├─ RPC timeout (50-100ms)
   └─ Batch size (최대 로그 항목 수)

☑️ Debugging
   ├─ 상세한 로깅
   ├─ 상태 덤프
   └─ 네트워크 추적
```

---

## 🎯 결론

**Raft의 강점**:
1. ✅ 이해하기 쉬움
2. ✅ 증명 가능함
3. ✅ 구현하기 간단함
4. ✅ 확장 가능함

**FreeLang 구현 난이도**: ⭐⭐⭐⭐ (4/5, 고급)

**예상 구현 시간**: 4주 (주 40시간 기준)

**최종 코드량**: 2,500줄 (구현) + 1,500줄 (테스트)

---

**작성**: Claude Code AI
**날짜**: 2026-03-02
**상태**: 완전 분석 완료
