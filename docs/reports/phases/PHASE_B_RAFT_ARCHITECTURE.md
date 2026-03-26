# 🏗️ Phase B: Raft 분산 합의 엔진 - 상세 아키텍처 설계

**프로젝트**: FreeLang Distributed System (Phase B)
**목표**: Raft 논문 완벽 재현 (FreeLang으로)
**기간**: 4주
**코드**: 2,200줄 (구현) + 1,200줄 (테스트) + 1,500줄 (문서)

---

## 📊 아키텍처 개요

```
┌─────────────────────────────────────────────────────────┐
│                 10-Node Raft Cluster                     │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────┐ │
│  │   Node 1       │  │   Node 2       │  │   Node 3   │ │
│  │   (Leader)     │  │  (Follower)    │  │(Follower)  │ │
│  └────────────────┘  └────────────────┘  └────────────┘ │
│         │                   │                    │       │
│         └───────────────────┼────────────────────┘       │
│              Raft RPC Protocol (TCP)                     │
│                                                           │
├─────────────────────────────────────────────────────────┤
│  각 노드 내부 구조:                                       │
│  ┌─────────────────────────────────────────────────┐   │
│  │ RaftNode                                        │   │
│  │ ├─ Persistent State (디스크)                    │   │
│  │ │  ├─ currentTerm                               │   │
│  │ │  ├─ votedFor                                  │   │
│  │ │  └─ log[]                                     │   │
│  │ ├─ Volatile State (메모리)                      │   │
│  │ │  ├─ commitIndex                               │   │
│  │ │  ├─ lastApplied                               │   │
│  │ │  └─ state (Leader/Follower/Candidate)        │   │
│  │ ├─ Leader's Volatile State                      │   │
│  │ │  ├─ nextIndex[]                               │   │
│  │ │  └─ matchIndex[]                              │   │
│  │ ├─ StateMachine (상태)                          │   │
│  │ │  └─ 은행 시스템 데이터                         │   │
│  │ └─ Storage (디스크 I/O)                         │   │
│  │    ├─ Log 저장                                  │   │
│  │    └─ Snapshot 저장                             │   │
│  └─────────────────────────────────────────────────┘   │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

---

## 🔄 상태 전환도 (State Machine)

```
┌─────────────┐
│  Follower   │ (초기 상태)
└──────┬──────┘
       │
       │ election timeout
       │ (150-300ms 동안 메시지 없음)
       │
       ▼
┌─────────────────┐
│   Candidate     │
│                 │
│ • Term 증가     │ ◄─────────┐
│ • Vote 요청     │           │
│ • 투표 기다림   │           │
└────────┬────────┘           │
         │                    │
         │ 과반수 투표 획득    │
         │                    │
         ▼                    │
┌────────────────┐           │
│    Leader      │           │
│                │           │
│ • heartbeat    │           │
│   전송         │           │
│ • 로그 복제    │           │
└────────┬───────┘           │
         │                   │
         │ higher term 수신  │
         │ 또는              │
         │ leader 발견       │
         │                   │
         └───────────────────┘
              ▼
        ┌──────────────┐
        │  Follower    │
        └──────────────┘
```

---

## 📝 Raft 알고리즘 - 3가지 핵심

### 1️⃣ **Leader Election** (리더 선출)

**시간흐름**:
```
Time: 0 ─────── 50ms ──────── 100ms ────── 150ms ────── 200ms
      │         │             │           │           │
      │ Node1   │ timeout     │           │ START     │
      │ Follower│ → Candidate │           │ ELECTION  │
      │         │             │           │           │
      │─────────┤ send votes  │ votes     │ got       │ become
      │         │ to Node2,3  │ return    │ 2 votes   │ leader
      │         │             │           │           │
      └─────────┴─────────────┴───────────┴───────────┘
```

**FreeLang 구현 구조**:
```
fn start_election():
  1. currentTerm 증가
  2. votedFor = 자신
  3. resetElectionTimer() (새로운 150-300ms)
  4. 모든 피어에게 RequestVote RPC 전송
  5. 응답 대기:
     - 과반수 이상 찬성 → 리더 선출
     - 더 높은 term 수신 → 팔로워로 복귀
     - timeout → 새로운 선거 시작

fn handle_request_vote(candidateTerm, candidateId, lastLogIndex, lastLogTerm):
  1. candidateTerm < currentTerm → 거부
  2. votedFor != null && votedFor != candidateId → 거부
  3. 후보의 로그가 더 최신인가? (Log up-to-date check)
     - 같은 term: 더 긴 로그가 최신
     - 다른 term: 더 큰 term이 최신
  4. 모두 통과 → 투표 YES
```

### 2️⃣ **Log Replication** (로그 복제)

**흐름**:
```
Client Request
      │
      ▼
┌───────────────────────────┐
│  Leader 수신 및 로그 추가  │
│  log[5] = {term:5, cmd}   │
└───────┬───────────────────┘
        │
        ├─────────────────────────┬──────────────────────┐
        ▼                         ▼                      ▼
  ┌─────────────┐        ┌─────────────┐        ┌─────────────┐
  │ Follower-A  │        │ Follower-B  │        │ Follower-C  │
  │ (index: 4)  │        │ (index: 4)  │        │ (index: 4)  │
  └─────────────┘        └─────────────┘        └─────────────┘
        │ AppendEntries        │                      │
        │ with log[5]          │                      │
        │                      │                      │
        ├─────────────────────────────────────────────┤
        │         Leader 대기: 과반수(2) 복제됨       │
        ▼
   log[5] 적용 (commitIndex 증가)

   상태 기계에 적용
   (은행 시스템 거래 실행)
```

**FreeLang 구현**:
```
fn append_entries(leaderTerm, leaderId, prevLogIndex, prevLogTerm, entries[], leaderCommit):
  1. leaderTerm < currentTerm → 거부
  2. log[prevLogIndex].term != prevLogTerm → 로그 불일치
  3. 기존 로그 충돌 해결:
     - 새 로그로 덮어쓰기
  4. 새 로그 추가
  5. commitIndex 업데이트
  6. 상태 기계에 적용

fn replicate_to_follower(followerId):
  • nextIndex[followerId] 확인
  • AppendEntries RPC 전송
  • 실패 시:
    - nextIndex 감소
    - 재시도 (지수 백오프)
  • 성공 시:
    - matchIndex[followerId] 증가
    - commitIndex 갱신 (과반수 조건)
```

### 3️⃣ **Safety & Correctness** (안정성 증명)

**핵심 성질**:
```
┌────────────────────────────────────────────────────┐
│ 1. Election Safety                                 │
│    • 한 term에 최대 하나의 리더만 선출 가능        │
│    검증: votedFor는 최대 1개                       │
│                                                    │
│ 2. Leader Append-Only                             │
│    • 리더는 로그 항목을 삭제/수정하지 않음         │
│    검증: 로그는 append-only 구조                   │
│                                                    │
│ 3. State Machine Safety                           │
│    • 모든 노드가 같은 log index에 같은 명령 적용   │
│    검증: 과반수 복제 + 상태 기계 적용             │
│                                                    │
│ 4. Log Matching Property                          │
│    • 같은 index, term → 같은 명령                 │
│    검증: prevLogTerm 체크                         │
│                                                    │
│ 5. Leader Completeness                           │
│    • 리더는 모든 committed 항목을 포함함          │
│    검증: 높은 term의 로그만 선택                  │
└────────────────────────────────────────────────────┘
```

---

## 📂 파일 구조 및 모듈 설계

### **src/ 디렉토리 구조**

```
freelang-distributed-system/src/
│
├─ raft_node.fl (500줄) ⭐ 핵심
│  ├─ struct RaftNode
│  │  ├─ id: int (노드 ID)
│  │  ├─ currentTerm: int (현재 term)
│  │  ├─ votedFor: int? (투표한 후보, null이면 미투표)
│  │  ├─ log: []LogEntry (로그 배열)
│  │  ├─ commitIndex: int (적용된 로그 인덱스)
│  │  ├─ lastApplied: int (마지막 적용 인덱스)
│  │  ├─ state: enum (Leader/Follower/Candidate)
│  │  ├─ electionTimer: Timer (election timeout)
│  │  ├─ heartbeatTimer: Timer (150ms)
│  │  └─ peers: []int (다른 노드 ID들)
│  │
│  ├─ fn create_node(id, peers)
│  ├─ fn step() (main loop - 100ms마다 호출)
│  └─ fn handle_rpc(rpc_type, args)
│
├─ leader_election.fl (400줄)
│  ├─ fn start_election()
│  ├─ fn handle_request_vote_request(args) → response
│  ├─ fn handle_request_vote_response(response)
│  ├─ fn reset_election_timer()
│  ├─ fn send_vote_request_to_all()
│  └─ fn count_votes()
│
├─ log_replication.fl (400줄)
│  ├─ fn append_to_log(command) → LogEntry
│  ├─ fn handle_append_entries_request(args) → response
│  ├─ fn handle_append_entries_response(response)
│  ├─ fn replicate_to_follower(followerId)
│  ├─ fn advance_commit_index()
│  └─ fn check_log_conflict(prevLogIndex, prevLogTerm)
│
├─ state_machine.fl (300줄)
│  ├─ struct StateMachine
│  │  ├─ lastApplied: int
│  │  └─ state: BankSystem (은행 시스템)
│  │
│  ├─ fn apply_log_entry(entry)
│  │  └─ 상태 기계에 명령 적용
│  │
│  ├─ fn take_snapshot(index) → Snapshot
│  └─ fn load_snapshot(snapshot)
│
├─ rpc_server.fl (350줄)
│  ├─ fn start_rpc_server(node, port)
│  ├─ fn handle_request_vote_rpc(request)
│  ├─ fn handle_append_entries_rpc(request)
│  ├─ fn send_rpc_to_peer(peer, method, args)
│  └─ fn serialize_rpc(rpc) → json string
│
├─ storage.fl (250줄)
│  ├─ fn write_log_entry(entry) (디스크에 저장)
│  ├─ fn read_log(fromIndex, toIndex) (디스크에서 읽음)
│  ├─ fn write_metadata(currentTerm, votedFor)
│  ├─ fn read_metadata() → (term, votedFor)
│  ├─ fn create_snapshot(index, lastIncludedTerm)
│  └─ fn load_snapshot() → (lastIncludedIndex, data)
│
├─ types.fl (150줄)
│  ├─ struct LogEntry
│  │  ├─ term: int
│  │  ├─ command: string
│  │  └─ index: int
│  │
│  ├─ struct RequestVoteRequest
│  │  ├─ term: int
│  │  ├─ candidateId: int
│  │  ├─ lastLogIndex: int
│  │  └─ lastLogTerm: int
│  │
│  ├─ struct RequestVoteResponse
│  │  ├─ term: int
│  │  └─ voteGranted: bool
│  │
│  ├─ struct AppendEntriesRequest
│  │  ├─ term: int
│  │  ├─ leaderId: int
│  │  ├─ prevLogIndex: int
│  │  ├─ prevLogTerm: int
│  │  ├─ entries: []LogEntry
│  │  └─ leaderCommit: int
│  │
│  └─ struct AppendEntriesResponse
│     ├─ term: int
│     └─ success: bool
│
└─ cluster.fl (200줄)
   ├─ fn create_cluster(num_nodes, num_failures) → [RaftNode]
   ├─ fn run_cluster_step() (모든 노드 step())
   ├─ fn inject_failure(nodeId, failure_type)
   └─ fn check_cluster_consistency() → bool
```

---

## ⏱️ 타임라인 및 주간 계획

### **Week 1: 기초 구현 (500줄)**

```
Day 1-2: 데이터 구조 정의
├─ RaftNode struct
├─ LogEntry, RPC types
├─ Storage interface
└─ 기본 테스트 1개

Day 3-4: RPC 서버 구현
├─ TCP 리스너 (포트 9000 + nodeId*100)
├─ 직렬화/역직렬화 (JSON)
├─ 기본 요청/응답 처리
└─ 테스트 3개

Day 5: 로그 저장소
├─ WAL (Write-Ahead Logging)
├─ 메타데이터 (currentTerm, votedFor)
└─ 테스트 2개

완성: RPC + Storage 기본 구현 (테스트 6개)
```

### **Week 2: Leader Election (400줄)**

```
Day 1-2: Election 로직
├─ start_election()
├─ request_vote() 요청/응답
├─ timer 관리 (150-300ms random)
└─ 테스트 8개 (선거 시나리오)

Day 3-4: 동시성 처리
├─ 여러 선거 동시 진행
├─ term 업데이트 경쟁
├─ vote counting
└─ 테스트 7개 (복잡한 시나리오)

Day 5: 정합성 검증
├─ leader 선출 유일성 (election safety)
├─ 과반수 검증
└─ 테스트 5개

완성: Election 완벽 작동 (테스트 20개)
```

### **Week 3: Log Replication (400줄)**

```
Day 1-2: append_entries 기본 구현
├─ 로그 항목 추가
├─ replicate_to_followers()
├─ conflict resolution
└─ 테스트 8개

Day 3-4: 상태 기계 적용
├─ commitIndex 관리
├─ apply_log_entries()
├─ 은행 시스템과 연동
└─ 테스트 10개

Day 5: Snapshot & Log Compaction
├─ 스냅샷 생성
├─ 로그 압축
└─ 테스트 7개

완성: Replication 완벽 작동 (테스트 25개)
```

### **Week 4: 검증 & 최적화 (300줄)**

```
Day 1-2: 네트워크 실패 시뮬레이션
├─ 패킷 손실 (10-50%)
├─ 지연 추가 (50-500ms)
├─ 노드 충돌 시뮬레이션
└─ 테스트 10개 (chaos engineering)

Day 3-4: 최종 검증
├─ 10개 노드, 5개 장애
├─ 1000개 로그 항목
├─ 일관성 검증
└─ 테스트 10개

Day 5: 성능 최적화
├─ 메모리 프로파일링
├─ CPU 효율성
├─ 병목 제거
└─ 완료 보고서 작성

완성: 논문급 Raft 구현 (테스트 50개, 100% 통과)
```

---

## 🧪 테스트 시나리오 (50개)

### **Phase 1: Basic Functionality (15개)**
```
1. test_node_creation
2. test_initial_state
3. test_log_append
4. test_election_start
5. test_single_election
6. test_multi_election_simultaneous
7. test_log_replication_basic
8. test_log_replication_multiple
9. test_commit_advance
10. test_state_apply
11. test_rpc_serialization
12. test_rpc_deserialization
13. test_storage_write_read
14. test_storage_metadata
15. test_cluster_creation
```

### **Phase 2: Correctness (20개)**
```
16. test_election_safety (한 term 한 leader)
17. test_leader_append_only (로그 수정 없음)
18. test_state_machine_consistency (모든 노드 같은 상태)
19. test_log_matching_property (같은 index → 같은 명령)
20. test_leader_completeness (리더가 모든 commit 포함)
21-25. test_candidate_scenarios_1_to_5
26-30. test_log_conflict_resolution_1_to_5
31-35. test_partial_replication_failures_1_to_5
```

### **Phase 3: Failure Handling (15개)**
```
36. test_network_partition (분할된 네트워크)
37. test_leader_crash (리더 충돌)
38. test_follower_crash (팔로워 충돌)
39. test_packet_loss (패킷 손실 10%)
40. test_packet_loss_50_percent (50%)
41. test_network_delay (지연 100ms)
42. test_network_delay_500ms (500ms)
43. test_cascading_failures (연쇄 장애)
44. test_recovery_after_partition (분할 복구)
45. test_log_divergence_resolution (로그 발산)
46. test_snapshot_recovery (스냅샷 복구)
47. test_rapid_leader_changes (빠른 리더 변경)
48. test_concurrent_client_requests (동시 요청)
49. test_stale_term_handling (오래된 term)
50. test_final_consistency_check (최종 일관성)
```

---

## 📈 성능 목표

| 지표 | 목표 | 측정 방법 |
|------|------|---------|
| **Election Time** | < 300ms | timeout 기준 |
| **Log Replication** | < 50ms | append_entries 왕복 |
| **Commit Latency** | < 100ms | 클라이언트 요청 → 적용 |
| **Throughput** | 1,000 ops/sec | 동시 요청 1000개 |
| **Memory Per Node** | < 10MB | 1000개 로그 항목 |
| **Network Utilization** | < 10Mbps | 10 노드 클러스터 |

---

## 🔗 다음 단계 (Phase C)

Phase B 완료 후:
- Phase C: 고성능 로드 밸런서 & 프록시 설계
- Phase D: 전체 통합 및 검증

---

**작성**: Claude Code AI
**날짜**: 2026-03-02
**상태**: 아키텍처 설계 완료 (구현 준비 완료)
