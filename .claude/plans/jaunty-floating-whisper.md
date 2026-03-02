# Raft Consensus based Sharded DB - 구현 계획

## Context

Kim님이 스스로 비판했던 "Leader election 없는 Raft"를 정면 돌파하여, 실제 금융권 수준의 정합성을 가진 분산 DB를 구축하는 프로젝트.

**목표**: 3대 이상의 노드가 동적으로 리더를 선출하고, 데이터가 유실되지 않으며, 수평적 확장(Sharding)이 가능한 분산 DB 완성.

---

## 프로젝트 구조

```
freelang-raft-db/
├── src/
│   ├── raft_core.fl         (600줄, Week 1)  - 완전한 Raft 합의 알고리즘
│   ├── sharding.fl          (500줄, Week 2)  - Consistent Hashing + Sharding
│   ├── simulation.fl        (700줄, Week 3)  - Deterministic Simulation Testing
│   └── integration.fl       (400줄, Week 4)  - 완전한 통합 & 시나리오
├── RAFT_WEEK1_REPORT.md
├── RAFT_WEEK2_REPORT.md
├── RAFT_WEEK3_REPORT.md
├── RAFT_WEEK4_REPORT.md
└── README.md
```

---

## Week 1: Core Raft Algorithm (600줄)

**파일**: `src/raft_core.fl`

### 핵심 구조체

```rust
// 노드 상태 (3개)
pub enum RaftState { Follower, Candidate, Leader }

// 로그 엔트리
pub struct LogEntry {
    term: u64,
    index: u64,
    command: String,
}

// 투표 요청/응답
pub struct VoteRequest {
    term: u64,
    candidate_id: u32,
    last_log_index: u64,
    last_log_term: u64,
}

// 로그 복제 요청
pub struct AppendEntriesRequest {
    term: u64,
    leader_id: u32,
    prev_log_index: u64,
    prev_log_term: u64,
    entries: Vec<LogEntry>,
    leader_commit: u64,
}

// Raft 노드 핵심
pub struct RaftNode {
    id: u32,
    state: RaftState,
    current_term: u64,
    voted_for: Option<u32>,
    log: Vec<LogEntry>,
    commit_index: u64,
    last_applied: u64,
    election_timeout_ms: u32,    // 150-300ms 랜덤
    votes_received: u32,
    next_index: Vec<u64>,        // Leader 전용
    match_index: Vec<u64>,       // Leader 전용
}
```

### 핵심 메서드 (8개)

| 메서드 | 설명 |
|--------|------|
| `start_election()` | Candidate로 전환, term 증가, 투표 요청 |
| `handle_vote_request(req)` | 투표 처리 (Safety 조건 검증) |
| `become_leader()` | Leader 승격, next/match_index 초기화 |
| `replicate_log(entry)` | 과반수에게 AppendEntries 전송 |
| `handle_append_entries(req)` | 로그 복제 수신 처리 |
| `commit_entries()` | 과반수 확인 후 commit_index 업데이트 |
| `apply_committed()` | State Machine에 적용 |
| `reset_election_timeout()` | Heartbeat 수신 시 타이머 리셋 |

### 5가지 Safety 조건

1. **Election Safety**: 한 Term에 최대 1명의 Leader
2. **Leader Append-Only**: 기존 로그 덮어쓰지 않음
3. **Log Matching**: 같은 index+term이면 내용도 동일
4. **Leader Completeness**: 선출된 리더는 모든 committed 엔트리 보유
5. **State Machine Safety**: 모든 노드는 동일 순서로 command 적용

### 테스트 (20개)

- test_initial_follower_state
- test_term_increment_on_election
- test_candidate_vote_request
- test_vote_granted_conditions
- test_vote_rejected_on_old_term
- test_leader_election_majority
- test_split_vote_reelection
- test_leader_heartbeat
- test_log_replication_3nodes
- test_commit_requires_majority
- test_safety_no_two_leaders
- test_log_matching_property
- test_leader_completeness
- test_follower_catchup
- test_stale_leader_detection
- test_term_monotonic_increase
- test_apply_committed_entries
- test_network_partition_recovery
- test_leader_crash_reelection
- test_5node_cluster

---

## Week 2: Consistent Hashing + Sharding (500줄)

**파일**: `src/sharding.fl`

```rust
pub struct ConsistentHashRing {
    ring: BTreeMap<u64, NodeId>,
    virtual_nodes: u32,          // 100~150
    nodes: Vec<NodeId>,
}

pub struct Shard {
    id: u32,
    key_range: (String, String),
    primary_node: NodeId,
    replica_nodes: Vec<NodeId>,
    size_bytes: u64,
}

pub struct ShardedKV {
    shards: Vec<Shard>,
    ring: ConsistentHashRing,
    raft_nodes: Vec<RaftNode>,
}
```

### 핵심 기능

- `add_node(id)`: 노드 추가 → 리밸런싱 최소화
- `remove_node(id)`: 노드 제거 → 데이터 이동 최소화
- `get_shard_for_key(key)`: 샤드 결정
- `rebalance()`: 균등 분배

### 테스트 (15개)

- test_hash_ring_creation
- test_consistent_mapping
- test_add_node_minimal_movement
- test_remove_node_redistribution
- test_virtual_node_balance
- test_shard_key_lookup
- test_sharded_put_get
- test_replication_factor
- test_node_failure_redistribution
- test_hotspot_detection
- test_cross_shard_query
- test_range_scan
- test_shard_split
- test_shard_merge
- test_load_distribution

---

## Week 3: Deterministic Simulation Testing (700줄)

**파일**: `src/simulation.fl`

```rust
pub enum SimEvent {
    NodeCrash(NodeId),
    NodeRestart(NodeId),
    NetworkDelay { from: NodeId, to: NodeId, delay_ms: u32 },
    PacketLoss { from: NodeId, to: NodeId, rate: f32 },
    NetworkPartition { group_a: Vec<NodeId>, group_b: Vec<NodeId> },
    ClientRequest { key: String, value: String },
}

pub struct SimulationEngine {
    random_seed: u64,
    events: Vec<(u64, SimEvent)>,
    nodes: Vec<RaftNode>,
    network: NetworkSimulator,
    log: Vec<SimulationResult>,
}

pub struct NetworkSimulator {
    latency_ms: u32,
    packet_loss_rate: f32,
    partitioned: bool,
    message_queue: Vec<(u64, Message)>,
}
```

### 핵심 시나리오 (5개)

1. Leader Election: 리더 사망 → 재선출 검증
2. Network Partition: 분리 → 복구 검증
3. Log Divergence: 충돌 후 로그 정합성 검증
4. High Load: 1000 writes/sec 부하 테스트
5. Cascading Failure: 연쇄 노드 사망 복구

### 테스트 (20개)

- test_deterministic_replay_same_seed
- test_election_after_leader_crash
- test_no_data_loss_on_crash
- test_network_partition_minority_stops
- test_partition_heal_convergence
- test_split_brain_prevention
- test_log_divergence_resolution
- test_stale_leader_rejection
- test_packet_loss_resilience
- test_high_latency_election
- test_concurrent_elections
- test_5node_2crash_recovery
- test_replication_under_partition
- test_write_after_partition_heal
- test_read_consistency
- test_3node_election_speed
- test_cascading_failure
- test_log_compaction
- test_snapshot_install
- test_chaos_random_seed

---

## Week 4: 완전 통합 (400줄)

**파일**: `src/integration.fl`

```rust
pub struct RaftShardedDB {
    shards: Vec<ShardedKV>,
    ring: ConsistentHashRing,
    nodes: Vec<RaftNode>,
    sim: Option<SimulationEngine>,
}
```

### 최종 시나리오 (4개)

1. **3노드 클러스터**: Leader 선출 → Write → Commit → Read
2. **노드 사망 복구**: 사망 → 재선출 → 데이터 유실 0
3. **샤드 리밸런싱**: 노드 추가 → 최소 이동 → 균등 분배
4. **카오스 테스트**: 무작위 장애 → Safety 조건 유지

### 테스트 (15개)

- test_full_cluster_setup
- test_end_to_end_write_read
- test_sharded_db_operations
- test_leader_election_in_sharded_db
- test_cross_shard_consistency
- test_node_failure_no_data_loss
- test_scaling_add_node
- test_scaling_remove_node
- test_chaos_no_safety_violation
- test_performance_1000_writes
- test_concurrent_clients
- test_read_write_consistency
- test_snapshot_and_restore
- test_cluster_restart
- test_production_scenario

---

## 최종 목표

| 항목 | Week 1 | Week 2 | Week 3 | Week 4 | 합계 |
|------|--------|--------|--------|--------|------|
| 파일 | 1 | 1 | 1 | 1 | **4개** |
| 줄수 | 600 | 500 | 700 | 400 | **2,200줄** |
| 테스트 | 20 | 15 | 20 | 15 | **70개** |

**달성 시**:
- ✅ 3+ 노드 동적 Leader 선출 (150~300ms)
- ✅ 데이터 유실 0% (Log Replication + Majority Commit)
- ✅ 수평 확장 (Consistent Hashing + Virtual Nodes)
- ✅ 100% 재현 가능 테스트 (Deterministic Simulation)
- ✅ 5가지 Safety 조건 완벽 충족
- ✅ 금융권 수준 정합성

## GOGS 저장소

- 이름: `raft-sharded-db`
- URL: `https://gogs.dclub.kr/kim/raft-sharded-db.git`
- 커밋: Week별 1회 (총 4개)
