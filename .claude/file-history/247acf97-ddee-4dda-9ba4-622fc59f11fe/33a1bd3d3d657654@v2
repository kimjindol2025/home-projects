# Phase 5: Distributed System Integration Design

## Context
확장된 Neural-Kernel-Sentinel을 **3-5개 노드의 분산 환경**에서 운영하기 위한 Raft 기반 합의 엔진과 Byzantine Fault Tolerance를 구현합니다.

**목표**:
- 99.99% 가용성 달성 (4개 9's)
- 노드 장애 자동 감지 및 대응
- 분산 메트릭 실시간 동기화
- 자동 네트워크 분할 복구 (Network Partition Healing)

---

## Architecture Overview

```
                    ┌─────────────────────────────────────────┐
                    │  NeuralKernelSentinel Cluster (3-5 nodes)│
                    └─────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
    ┌───▼───┐            ┌───▼───┐            ┌───▼───┐
    │ Node1 │            │ Node2 │            │ Node3 │
    │(Leader)│           │(Voter)│           │(Voter)│
    └───┬───┘            └───┬───┘            └───┬───┘
        │                    │                    │
        ├────────────────────┼────────────────────┤
        │   Raft Log & State Machine              │
        │   - Consensus History                   │
        │   - Threat Decisions                    │
        │   - Response Actions                    │
        │                                         │
        ├────────────────────┼────────────────────┤
        │   Real-time Metrics Sync                │
        │   - Detection Scores                    │
        │   - Latency Measurements                │
        │   - Resource Utilization                │
        │                                         │
        └────────────────────┼────────────────────┘
                             │
                   ┌─────────▼─────────┐
                   │  Byzantine Fault  │
                   │  Tolerance Layer  │
                   │ (Weighted Voting) │
                   └───────────────────┘
```

---

## 4개 모듈 설계

### 1. `node_coordinator.fl` (Day 1-2, ~400줄)

**목표**: 클러스터 내 노드 간 신뢰도 가중 합의

**핵심 구조체**:
```rust
pub enum NodeRole { Leader, Voter, Candidate, Follower }
pub enum NodeState { Healthy, Degraded, Unreachable, Dead }

pub struct NodeInfo {
    node_id: u32,
    role: NodeRole,
    state: NodeState,
    reliability_score: f64,  // 0.0 ~ 1.0
    last_heartbeat_ns: u64,
    metrics: ClusterMetrics,
}

pub struct ClusterMetrics {
    detection_score_avg: f64,
    latency_p99_us: u32,
    throughput_syscall_sec: u64,
    memory_usage_mb: u32,
    last_response_timestamp_ns: u64,
}

pub struct WeightedVote {
    node_votes: HashMap<u32, (f64, Decision)>,  // node_id -> (reliability_weight, decision)
    consensus_threshold: f64,  // 보통 0.66
    weighted_consensus: f64,
}

pub struct NodeCoordinator {
    nodes: HashMap<u32, NodeInfo>,
    local_node_id: u32,
    heartbeat_interval_ms: u32,  // 기본 100ms
    failure_timeout_ms: u32,     // 기본 300ms
}
```

**주요 함수**:

1. **`calculate_reliability_score(node: &NodeInfo) -> f64`**
   - 성공 비율 계산 (최근 100개 결정): `score = success_count / 100`
   - 응답 시간 페널티: `if latency > target { score *= 0.9 }`
   - 가용성 계산: `score = (uptime_seconds / total_seconds) * 100`
   - 최종: `reliability = 0.4×success + 0.3×latency + 0.3×availability`

2. **`perform_weighted_consensus(decision: &ThreatDecision) -> Decision`**
   - 각 노드의 신뢰도 점수를 가중치로 사용
   - `weighted_score = Σ(node_reliability × node_decision_confidence) / Σ(node_reliability)`
   - 합의 달성 조건: `weighted_score > consensus_threshold (기본 0.66)`
   - 결과: `ControlDecision::Isolated` (권장) or `ControlDecision::Alert` (보수)

3. **`detect_node_failure(node_id: u32) -> bool`**
   - 마지막 heartbeat 이후 경과 시간 확인
   - `if now_ns - last_heartbeat_ns > failure_timeout_ms * 1_000_000 { return true }`
   - 3번 연속 실패 시 노드 상태를 Dead로 변경

4. **`handle_network_partition() -> RecoveryStrategy`**
   - 파티션 양쪽이 각각 합의를 시도
   - minority partition: 새로운 결정 중지, 캐시된 결정만 사용
   - majority partition: 정상적으로 운영 계속
   - 재연결 시: 따라잡기 복제(catchup replication) 수행

**테스트 A1-A6**:
- A1: 신뢰도 점수 계산 (성공률, 지연, 가용성)
- A2: 가중 투표 합의 (다양한 신뢰도 조합)
- A3: 노드 장애 감지 (heartbeat timeout)
- A4: 네트워크 분할 감지 및 recovery
- A5: Leader 선출 검증
- A6: 합의 성공률 > 99%

---

### 2. `consensus.fl` (Day 3-4, ~350줄)

**목표**: Raft 기반 분산 합의 및 로그 복제

**핵심 구조체**:
```rust
pub enum RaftState { Follower, Candidate, Leader }

pub struct RaftLog {
    entries: Vec<LogEntry>,
    commit_index: u32,
    last_applied: u32,
}

pub struct LogEntry {
    index: u32,
    term: u32,
    command: DecisionCommand,  // 위협 결정, 응답 액션
    timestamp_ns: u64,
}

pub struct RaftEngine {
    current_term: u32,
    voted_for: Option<u32>,
    state: RaftState,
    log: RaftLog,
    leader_id: Option<u32>,
    election_timeout_ms: u32,  // 150~300ms
}

pub struct SnapshotManager {
    snapshots: Vec<Snapshot>,
    last_included_index: u32,
    last_included_term: u32,
}

pub struct Snapshot {
    index: u32,
    term: u32,
    state_machine_data: Vec<u8>,  // 직렬화된 클러스터 상태
    timestamp_ns: u64,
}
```

**주요 함수**:

1. **`append_entries_rpc(leader_id, prev_index, prev_term, entries) -> bool`**
   - Raft AppendEntries RPC 구현
   - 로그 일관성 확인: `if log[prev_index].term != prev_term { return false }`
   - 새 항목 추가 및 커밋 인덱스 업데이트
   - 응답: `{ success: true, term: current_term }`

2. **`request_vote_rpc(candidate_id, term, last_log_index, last_log_term) -> bool`**
   - Raft 리더 선출 RPC
   - 투표 권한 확인: `if term < current_term { return false }`
   - Candidate 로그가 더 최신인지 확인 (term 우선, 그 다음 index)
   - 투표 기록: `voted_for = candidate_id`

3. **`replicate_log_entry(entry: LogEntry) -> ReplicationStatus`**
   - 모든 follower에 로그 항목 복제
   - 재시도 로직: exponential backoff (100ms → 200ms → 400ms)
   - 성공 조건: majority (N/2+1) 노드가 복제 완료
   - 결과: `ReplicationStatus::Committed` or `ReplicationStatus::Failed`

4. **`create_snapshot() -> Snapshot`**
   - 현재 상태 머신 스냅샷 생성
   - 압축: 로그의 apply된 항목까지만 저장
   - 공간 절감: log 크기 > 10MB일 때 자동 스냅샷
   - 새 노드 가입 속도 향상

5. **`handle_leader_election() -> bool`**
   - Term 증가: `current_term += 1`
   - 자신에게 투표: `voted_for = self.node_id`
   - 모든 노드에 요청 투표 RPC 전송
   - 성공 조건: majority 투표 획득
   - 실패 시: 다음 election timeout까지 기다린 후 재시도

**테스트 B1-B6**:
- B1: AppendEntries RPC 정상 동작
- B2: RequestVote RPC 및 리더 선출
- B3: 로그 복제 (majority)
- B4: 스냅샷 생성 및 로드
- B5: 노드 추가/제거 (dynamic membership)
- B6: 합의율 > 99.99%

---

### 3. `failover.fl` (Day 5-6, ~300줄)

**목표**: 자동 노드 교체 및 트래픽 재분배

**핵심 구조체**:
```rust
pub enum FailoverStrategy { Active, Passive, GracefulDegradation }

pub struct FailoverController {
    leader: NodeInfo,
    failover_timeout_ms: u32,  // 기본 500ms
    strategy: FailoverStrategy,
    backup_nodes: Vec<NodeInfo>,
}

pub struct LoadBalancer {
    nodes: Vec<NodeInfo>,
    weights: HashMap<u32, f64>,  // node_id -> weight (신뢰도 기반)
    routing_table: HashMap<String, u32>,  // PID -> best_node
}

pub struct HealthChecker {
    check_interval_ms: u32,  // 기본 100ms
    failure_threshold: u32,  // 연속 실패 3회 이상
    health_metrics: HashMap<u32, HealthStatus>,
}

pub enum HealthStatus { Healthy, Degraded, Critical, Failed }

pub struct HealthMetrics {
    cpu_usage_percent: f32,
    memory_usage_percent: f32,
    response_time_p99_us: u32,
    error_rate_percent: f32,
    last_check_ns: u64,
}
```

**주요 함수**:

1. **`detect_leader_failure() -> Option<u32>`**
   - 리더로부터의 heartbeat 모니터링
   - `if now_ns - last_heartbeat_ns > election_timeout { trigger_election() }`
   - 새 리더 선출: Raft engine의 `handle_leader_election()`
   - 반환: 새로운 리더 ID

2. **`perform_failover(failed_node: &NodeInfo) -> bool`**
   - 장애 노드의 책임 다른 노드에 분산
   - 상태 머신 스냅샷 다운로드
   - 미복제 로그 항목 복제
   - 트래픽 재분배: `recalculate_weights_and_routes()`

3. **`recalculate_weights_and_routes()`**
   - 각 노드의 신뢰도 점수 재계산
   - 가중치 정규화: `weights[node] = reliability_score / Σ(reliability_scores)`
   - 라우팅 테이블 업데이트: PID별 최선의 노드 선택
   - 목표: 신뢰도 높은 노드로 트래픽 쏠림

4. **`health_check_all() -> HealthSummary`**
   - 모든 노드에 상태 확인 요청
   - 메트릭 수집:
     - CPU: `ps aux | awk '{print $3}' | tail -1`
     - Memory: `/proc/meminfo` 파싱
     - 응답시간: 마지막 100개 응답의 P99
     - 오류율: 최근 1000개 결정 중 실패 비율
   - 실패 감지: 연속 3회 이상 타임아웃

5. **`graceful_degradation(failed_count: u32) -> DegradationMode`**
   - `if failed_count == 1 { return Normal }`
   - `if failed_count == 2 { return Conservative (높은 임계값) }`
   - `if failed_count >= 3 { return ReadOnly (읽기만) }`
   - Alert: `send_alert(AlertSeverity::Critical)`

**테스트 C1-C6**:
- C1: 노드 장애 감지
- C2: Failover 수행 (로그 복제)
- C3: 트래픽 재분배
- C4: Health check 정확도
- C5: Graceful degradation
- C6: 복구 시간 < 500ms

---

### 4. `mod.fl` (Day 7, ~200줄)

**목표**: 분산 시스템 통합 API 및 E2E 파이프라인

**핵심 구조체**:
```rust
pub struct DistributedNeuralKernelSentinel {
    coordinator: NodeCoordinator,
    consensus: RaftEngine,
    failover: FailoverController,
    load_balancer: LoadBalancer,
    health_checker: HealthChecker,

    // 로컬 인스턴스
    local_sentinel: NeuralKernelSentinel,

    // 메트릭 수집
    cluster_metrics: ClusterMetrics,
    decision_log: Vec<Decision>,
}

pub struct Decision {
    timestamp_ns: u64,
    threat_id: String,
    local_score: f64,
    consensus_score: f64,
    final_action: ResponseAction,
    nodes_agreed: Vec<u32>,
    success: bool,
}
```

**주요 함수**:

1. **`process_distributed_syscall_batch(events: Vec<SyscallEvent>)`**
   - 로컬 처리: `local_sentinel.process_syscall_batch()`
   - 합의 필요 여부 판단: `if confidence < 0.7 { consensus_required = true }`
   - 합의 수행: `weighted_consensus(local_decision)`
   - 최종 결정: `execute_response_action(consensus_decision)`

2. **`heartbeat_and_sync()`** (백그라운드)
   - 100ms 간격으로 반복
   - 모든 노드로 heartbeat 전송
   - 메트릭 수집 및 동기화
   - 합의 상태 확인: `check_consensus_health()`

3. **`get_cluster_status() -> ClusterStatus`**
   ```
   {
     leader_id: Some(1),
     total_nodes: 5,
     healthy_nodes: 4,
     failed_nodes: 1,
     availability_percent: 99.5,
     avg_decision_latency_us: 125,
     consensus_agreements: 1523,
     consensus_failures: 2,
     last_update_ns: <timestamp>
   }
   ```

4. **`metrics_dashboard() -> MetricsDashboard`**
   - 실시간 메트릭 수집
   - 시각화 데이터 생성
   - Prometheus 형식 내보내기

---

## 타이밍 예산 (E2E 파이프라인 < 500ms)

| 단계 | 예산 | 설명 |
|------|------|------|
| 로컬 처리 | 95µs | 기존 NKS (syscall 캡처 + 분석 + NN + 응답) |
| 합의 여부 결정 | 10µs | confidence 체크 |
| Raft 합의 | 50-100ms | 네트워크 RTT + 로그 복제 |
| Failover (필요 시) | 300-500ms | 노드 감지 + 리더 선출 + 상태 복제 |
| **정상 경로** | **95-110µs** | 로컬만 처리 |
| **합의 필요** | **50-150ms** | 로컬 + Raft |
| **Failover** | **300-500ms** | 최악의 경우 |

---

## 6가지 Chaos 시나리오

### Scenario 1: Node Failure
```
Node1(Leader) → Crash
  ↓
Election timeout (150-300ms)
  ↓
Node2가 리더 당선
  ↓
데이터 일관성 검증: all log entries replicated ✅
```

### Scenario 2: Network Partition
```
Node1-2 vs Node3 (minority)
  ↓
Minority: 합의 중지, 캐시된 결정만 사용
  ↓
재연결 시 minority catch-up
  ↓
최종 일관성 달성: all states synchronized ✅
```

### Scenario 3: Cascading Failures
```
Node1 ↓ → Node2가 Leader
Node2 ↓ → Node3가 Leader (3/5 majority 유지)
  ↓
Graceful degradation (Conservative 모드)
```

### Scenario 4: Slowness
```
Node1이 100x 느려짐 (p99 = 10s)
  ↓
가중치 감소: weight = 0.3 (신뢰도 ↓)
  ↓
트래픽 재분배: Node2-5로 쏠림
  ↓
합의 성공률 = 99% (Node1 제외)
```

### Scenario 5: Byzantine (Incorrect Decision)
```
Node1이 잘못된 결정 (confidence 0.2)
  ↓
Weighted consensus: 0.6×1.0 + 0.3×0.8 + 0.1×0.2 = 0.78 ✅
  ↓
Node1 신뢰도 감소
```

### Scenario 6: Recovery
```
Node1 recovery
  ↓
Snapshot 다운로드 (state sync)
  ↓
미복제 로그 항목 복제
  ↓
정상 운영 재개
```

---

## 무관용 규칙 (Unforgiving Rules)

### Rule 1: 99.99% 가용성
```
가용성 = Σ(healthy_time) / Σ(total_time)
목표: 가용성 > 99.99%
즉, 월 2.16초 미만 다운타임만 허용
```

### Rule 2: 합의 성공률 > 99.5%
```
success_rate = consensus_agreements / (consensus_agreements + consensus_failures)
목표: > 99.5%
즉, 1000번 합의 중 최대 5번 실패만 허용
```

### Rule 3: Failover 시간 < 500ms
```
failover_time = new_leader_election_time + state_replication_time
목표: < 500ms
최악의 경우: 300ms (election) + 100ms (replication) = 400ms
```

### Rule 4: 네트워크 분할 복구 시간 < 1s
```
recovery_time = re_connection_detection + catch_up_replication
목표: < 1s
```

### Rule 5: 로그 손실 = 0%
```
lost_entries = entries_in_minority_partition
목표: 항상 0 (majority partition만 진행)
```

### Rule 6: 데이터 일관성 = 100%
```
consistency = nodes_with_identical_state / total_nodes
목표: = 100% (어느 순간이든 majority는 동일한 상태)
```

### Rule 7: 클러스터 간 메트릭 동기 < 100ms
```
sync_latency = max_time_to_sync_all_nodes
목표: < 100ms
```

### Rule 8: Quorum 가중치 편향 < 10%
```
max_weight_difference = max(weights) - min(weights)
목표: < 10%
즉, 신뢰도가 높은 노드의 가중치가 너무 크지 않아야 함
```

---

## 구현 일정

- **Day 1-2**: `node_coordinator.fl` 구현 + 테스트 A1-A6
- **Day 3-4**: `consensus.fl` 구현 + 테스트 B1-B6
- **Day 5-6**: `failover.fl` 구현 + 테스트 C1-C6
- **Day 7**: `mod.fl` 통합 + 6가지 Chaos 시나리오 + 무관용 규칙 검증

---

## 성공 지표

| 지표 | 목표 | 달성 기준 |
|------|------|----------|
| 가용성 | 99.99% | 월 2.16초 미만 다운타임 |
| 합의 성공률 | 99.5% | 1000번 중 최대 5번 실패 |
| Failover 시간 | <500ms | 자동 복구 빠름 |
| 데이터 일관성 | 100% | 어느 순간이든 모든 노드 동일 |
| 합의 지연 | 50-150ms | Raft 네트워크 RTT 포함 |
| 메트릭 동기 | <100ms | 모든 노드 실시간 동기화 |

---

## 다음 단계 (Phase 6)

Phase 6부터는 분산 시스템 위에서 **ML 온라인 학습**을 구현합니다:
- Real-time weight update (SGD)
- Model drift detection
- Adaptive learning rates
- 새로운 공격 패턴 자동 학습 (<1일)
