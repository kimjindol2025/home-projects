# Raft Consensus DB - Group A 구현 계획

**프로젝트**: Raft Consensus based Sharded Database
**기간**: 2026-03-05 ~ 2026-03-06
**규모**: 1,000줄 + 25개 테스트

## 📋 Group A 목표

### 목표 1: Core Consensus Engine (400-500줄)
- [ ] **raft_state_machine.fl** (200줄)
  - State: Follower, Candidate, Leader
  - Term, VotedFor, CommitIndex, LastApplied
  - 상태 전환 로직

- [ ] **raft_log_entry.fl** (150줄)
  - LogEntry struct (Term, Index, Command)
  - LogStore with append/get/range operations
  - PrevLogIndex/PrevLogTerm 검증

- [ ] **leader_election.fl** (150줄)
  - RequestVote RPC protocol
  - Election timeout randomization
  - 과반수 투표 로직

### 목표 2: Replication Engine (300-400줄)
- [ ] **log_replication.fl** (250줄)
  - AppendEntries RPC protocol
  - Heartbeat mechanism (<100ms)
  - NextIndex, MatchIndex 추적
  - Committed entries 반영

- [ ] **snapshot_manager.fl** (150줄)
  - InstallSnapshot RPC
  - Snapshot persistence
  - Snapshot offset calculation

### 목표 3: Utilities & Integration (200-300줄)
- [ ] **raft_utils.fl** (150줄)
  - RaftConfig struct
  - Timer management (Election, Heartbeat)
  - Message serialization helpers

- [ ] **raft_integration.fl** (150줄)
  - RaftNode struct (assembling all components)
  - Process RequestVote
  - Process AppendEntries
  - Apply committed entries

## 🧪 Test Plan (25개 테스트)

### Group A1: State Machine (5개)
- [ ] **A1-T1**: Follower initial state
- [ ] **A1-T2**: Candidate → Leader transition
- [ ] **A1-T3**: Leader → Follower on higher term
- [ ] **A1-T4**: Term increment logic
- [ ] **A1-T5**: Multiple state transitions (state machine)

### Group A2: Log Entry (5개)
- [ ] **A2-T1**: LogEntry creation & storage
- [ ] **A2-T2**: Log append operations
- [ ] **A2-T3**: PrevLogIndex/PrevLogTerm validation
- [ ] **A2-T4**: Log range queries
- [ ] **A2-T5**: Log consistency check

### Group A3: Leader Election (5개)
- [ ] **A3-T1**: RequestVote with higher term
- [ ] **A3-T2**: Vote granted logic (term, log)
- [ ] **A3-T3**: Election timeout trigger
- [ ] **A3-T4**: Split vote scenario
- [ ] **A3-T5**: Majority voting

### Group A4: Log Replication (5개)
- [ ] **A4-T1**: AppendEntries heartbeat
- [ ] **A4-T2**: Log replication with entries
- [ ] **A4-T3**: NextIndex/MatchIndex update
- [ ] **A4-T4**: Committed index advancement
- [ ] **A4-T5**: Follower log consistency

### Group A5: Integration E2E (5개)
- [ ] **A5-T1**: Single node → Candidate → Leader
- [ ] **A5-T2**: Leader heartbeat to followers
- [ ] **A5-T3**: Log entry replication (1 entry)
- [ ] **A5-T4**: Committed entry application
- [ ] **A5-T5**: Snapshot installation

## 🎯 무관용 규칙 (Unforgiving Rules)

### R1: State Safety
- **정의**: 같은 index에서 다른 Term을 가진 log가 없어야 함
- **테스트**: A5-T4 (committed entries)
- **검증**: 모든 committed log는 불변

### R2: Leader Completeness
- **정의**: Leader는 모든 이전 committed entries를 포함
- **테스트**: A3-T5, A4-T5 (majority check)
- **검증**: Leader election 시 log completeness 확인

### R3: Election Liveness
- **정의**: Follower timeout 없이 Leader 선출 불가
- **테스트**: A3-T3, A3-T4
- **검증**: Timeout (100-300ms) 범위 내

### R4: Replication Safety
- **정의**: Replication 실패 시 entry는 committed되지 않음
- **테스트**: A4-T3, A4-T5
- **검증**: MatchIndex >= CommitIndex (majority)

### R5: Consistency (No Divergence)
- **정의**: 같은 index의 모든 logs는 같은 term을 가짐
- **테스트**: A2-T3, A2-T5
- **검증**: PrevLogIndex/PrevLogTerm 일치

## 📊 파일 구조

```
src/
  ├─ raft_state_machine.fl    (200줄)
  ├─ raft_log_entry.fl        (150줄)
  ├─ leader_election.fl       (150줄)
  ├─ log_replication.fl       (250줄)
  ├─ snapshot_manager.fl      (150줄)
  ├─ raft_utils.fl            (150줄)
  ├─ raft_integration.fl      (150줄)
  └─ mod.fl                   (업데이트)

tests/
  ├─ group_a_tests.fl         (350줄, 25개 테스트)
  └─ group_a_unforgiving.fl   (150줄, 5개 규칙 검증)

docs/
  └─ GROUP_A_SUMMARY.md       (이 파일)
```

## ✅ 구현 체크리스트

- [ ] raft_state_machine.fl 작성 및 테스트
- [ ] raft_log_entry.fl 작성 및 테스트
- [ ] leader_election.fl 작성 및 테스트
- [ ] log_replication.fl 작성 및 테스트
- [ ] snapshot_manager.fl 작성 및 테스트
- [ ] raft_utils.fl + raft_integration.fl 작성
- [ ] group_a_tests.fl 작성 (25개 테스트)
- [ ] group_a_unforgiving.fl 작성 (5개 규칙)
- [ ] 로컬 테스트 100% 통과
- [ ] GOGS 푸시 (커밋)
- [ ] GROUP_A_SUMMARY.md 작성

## 🚀 시작

**시작일**: 2026-03-05
**예상 완료**: 2026-03-06
**상태**: 준비 완료 ✅

---

**다음 단계**: raft_state_machine.fl 작성 시작
