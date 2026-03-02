# Phase 2: Raft RFC 5740 정확한 구현

**시작**: 2026-03-02 (즉시)
**목표**: RFC 5740 준수하는 Raft 구현으로 분산 합의 신뢰도 확보
**기한**: 4-5일 (2026-03-06 ~ 2026-03-07)
**성공 기준**: 다중 노드 리더 선출, 로그 복제, 장애 복구 모두 검증

---

## 🎯 Phase 2 핵심 목표

### 현재 문제점 (Phase 1에서 발견)
```
❌ 리더 선출이 라운드로빈 (결정론적이지만 현실적이지 않음)
❌ 실제 투표 메커니즘 없음
❌ 로그 복제 미구현
❌ 네트워크 지연 시뮬레이션 미흡
❌ 장애 복구 검증 불완전
```

### Phase 2로 해결할 것
```
✅ 투표 기반 리더 선출 (Random Timeout)
✅ 투표 카운팅 (Quorum 기반)
✅ Log Replication (AppendEntries RPC)
✅ 네트워크 시뮬레이션 (실제 지연)
✅ 자동 장애 복구 (Leader Failure → Re-election)
```

---

## 📋 Phase 2 구현 단계

### Step 1: Raft 핵심 구조 정의 (1일)

#### 1-1. RaftNode 확장
```rust
pub struct RaftNode {
    // 기존
    pub node_id: usize,
    pub state: NodeState,
    pub current_term: u64,
    pub voted_for: Option<usize>,
    pub log: Vec<LogEntry>,
    pub commit_index: u64,
    pub last_applied: u64,

    // 새로 추가
    // Leader 전용
    pub next_index: Vec<u64>,          // 각 노드에 전송할 다음 로그 인덱스
    pub match_index: Vec<u64>,         // 각 노드에 복제된 로그 인덱스

    // 모든 노드
    pub election_timeout_ms: u64,      // 선출 타임아웃 (150-300ms)
    pub heartbeat_interval_ms: u64,    // 하트비트 간격 (50ms)
    pub last_heartbeat: std::time::Instant,  // 마지막 하트비트 시간
    pub last_election_time: std::time::Instant, // 마지막 선출 시도 시간
}
```

#### 1-2. RPC 메시지 정의
```rust
// RequestVote RPC
pub struct RequestVoteRequest {
    pub term: u64,
    pub candidate_id: usize,
    pub last_log_index: u64,
    pub last_log_term: u64,
}

pub struct RequestVoteResponse {
    pub term: u64,
    pub vote_granted: bool,
}

// AppendEntries RPC (로그 복제)
pub struct AppendEntriesRequest {
    pub term: u64,
    pub leader_id: usize,
    pub prev_log_index: u64,
    pub prev_log_term: u64,
    pub entries: Vec<LogEntry>,
    pub leader_commit: u64,
}

pub struct AppendEntriesResponse {
    pub term: u64,
    pub success: bool,
    pub match_index: u64,
}
```

#### 1-3. 네트워크 시뮬레이터
```rust
pub struct RaftNetwork {
    nodes: Vec<Arc<RwLock<RaftNode>>>,
    pending_messages: Arc<DashMap<(usize, usize), VecDeque<Message>>>,
    latencies: Arc<DashMap<(usize, usize), u64>>,  // (from, to) → delay_ms
    message_loss: Arc<DashMap<(usize, usize), f64>>, // (from, to) → loss_rate
}
```

---

### Step 2: 리더 선출 구현 (1-2일)

#### 2-1. RequestVote Handler
```rust
pub async fn request_vote(&mut self, req: RequestVoteRequest) -> RequestVoteResponse {
    // 1. Term 업데이트
    if req.term > self.current_term {
        self.current_term = req.term;
        self.voted_for = None;
        self.state = NodeState::Follower;
    }

    // 2. 투표 거부 조건
    if req.term < self.current_term {
        return RequestVoteResponse {
            term: self.current_term,
            vote_granted: false,
        };
    }

    // 3. 이미 투표했거나 로그가 구식이면 거부
    if let Some(voted) = self.voted_for {
        if voted != req.candidate_id {
            return RequestVoteResponse {
                term: self.current_term,
                vote_granted: false,
            };
        }
    }

    // 4. Log Matching: 후보자의 로그가 우리 로그보다 최신인지 확인
    if req.last_log_term < self.last_log_term() {
        return RequestVoteResponse {
            term: self.current_term,
            vote_granted: false,
        };
    }

    if req.last_log_term == self.last_log_term() && req.last_log_index < self.log.len() as u64 {
        return RequestVoteResponse {
            term: self.current_term,
            vote_granted: false,
        };
    }

    // 5. 투표 부여
    self.voted_for = Some(req.candidate_id);
    RequestVoteResponse {
        term: self.current_term,
        vote_granted: true,
    }
}
```

#### 2-2. Election Timeout 처리
```rust
pub async fn check_election_timeout(&mut self) -> bool {
    let elapsed = self.last_heartbeat.elapsed().as_millis() as u64;

    // 리더는 선출 타임아웃 무시
    if matches!(self.state, NodeState::Leader) {
        return false;
    }

    // 팔로워가 타임아웃되면 후보자가 됨
    if elapsed > self.election_timeout_ms {
        self.state = NodeState::Candidate;
        self.current_term += 1;
        self.voted_for = Some(self.node_id);  // 자신에게 투표
        self.last_election_time = std::time::Instant::now();

        tracing::info!("Node {}: Election timeout, starting new election (term: {})",
            self.node_id, self.current_term);

        return true;
    }

    false
}
```

#### 2-3. 투표 수집 및 리더 선출
```rust
pub async fn conduct_election(network: &RaftNetwork, node_id: usize) -> bool {
    let node = &network.nodes[node_id];
    let mut node_lock = node.write().await;

    let my_term = node_lock.current_term;
    let my_last_index = node_lock.log.len() as u64;
    let my_last_term = node_lock.last_log_term();

    drop(node_lock);  // 락 해제

    // 모든 다른 노드에게 RequestVote 전송
    let mut votes = 1;  // 자신에게 투표
    let quorum = network.nodes.len() / 2 + 1;

    for i in 0..network.nodes.len() {
        if i == node_id {
            continue;
        }

        let req = RequestVoteRequest {
            term: my_term,
            candidate_id: node_id,
            last_log_index: my_last_index,
            last_log_term: my_last_term,
        };

        // 네트워크를 통해 RPC 전송 (지연 포함)
        let response = network.send_request_vote(i, req).await;

        if response.vote_granted {
            votes += 1;
        }
    }

    // Quorum 달성 여부 확인
    if votes >= quorum {
        let mut node_lock = node.write().await;
        if node_lock.current_term == my_term && matches!(node_lock.state, NodeState::Candidate) {
            node_lock.state = NodeState::Leader;
            node_lock.next_index = vec![node_lock.log.len() as u64; network.nodes.len()];
            node_lock.match_index = vec![0; network.nodes.len()];
            tracing::info!("Node {}: Won election with {} votes (term: {})",
                node_id, votes, my_term);
            return true;
        }
    }

    false
}
```

---

### Step 3: Log Replication 구현 (1-2일)

#### 3-1. AppendEntries Handler
```rust
pub async fn append_entries(&mut self, req: AppendEntriesRequest) -> AppendEntriesResponse {
    // 1. Term 업데이트
    if req.term > self.current_term {
        self.current_term = req.term;
        self.voted_for = None;
        self.state = NodeState::Follower;
    }

    // 2. Leader의 term이 낮으면 거부
    if req.term < self.current_term {
        return AppendEntriesResponse {
            term: self.current_term,
            success: false,
            match_index: 0,
        };
    }

    self.last_heartbeat = std::time::Instant::now();

    // 3. Log Matching: prev_log_index의 항목이 존재하고 term이 맞는지 확인
    if req.prev_log_index > 0 {
        if (req.prev_log_index as usize) > self.log.len() {
            // 로그가 부족함
            return AppendEntriesResponse {
                term: self.current_term,
                success: false,
                match_index: self.log.len() as u64,
            };
        }

        let prev_entry = &self.log[(req.prev_log_index - 1) as usize];
        if prev_entry.term != req.prev_log_term {
            // Term 불일치
            return AppendEntriesResponse {
                term: self.current_term,
                success: false,
                match_index: (req.prev_log_index - 1) as u64,
            };
        }
    }

    // 4. 새로운 엔트리 추가
    let mut match_index = req.prev_log_index;
    for (i, entry) in req.entries.iter().enumerate() {
        let log_index = req.prev_log_index + i as u64 + 1;

        if log_index as usize > self.log.len() {
            self.log.push(entry.clone());
        } else if self.log[(log_index - 1) as usize].term != entry.term {
            // 기존 엔트리와 충돌: 삭제 후 추가
            self.log.truncate((log_index - 1) as usize);
            self.log.push(entry.clone());
        }

        match_index = log_index;
    }

    // 5. Commit Index 업데이트
    if req.leader_commit > self.commit_index {
        self.commit_index = std::cmp::min(req.leader_commit, self.log.len() as u64);
    }

    AppendEntriesResponse {
        term: self.current_term,
        success: true,
        match_index,
    }
}
```

#### 3-2. Leader의 Log 복제
```rust
pub async fn replicate_log(network: &RaftNetwork, leader_id: usize) {
    let leader = &network.nodes[leader_id];
    let mut leader_lock = leader.write().await;

    if !matches!(leader_lock.state, NodeState::Leader) {
        return;
    }

    let my_term = leader_lock.current_term;
    let commit_index = leader_lock.commit_index;
    let log = leader_lock.log.clone();
    let next_index = leader_lock.next_index.clone();

    drop(leader_lock);

    // 각 팔로워에게 로그 복제
    for i in 0..network.nodes.len() {
        if i == leader_id {
            continue;
        }

        let prev_index = next_index[i] - 1;
        let prev_term = if prev_index > 0 {
            log[(prev_index - 1) as usize].term
        } else {
            0
        };

        let entries: Vec<LogEntry> = log[(next_index[i] as usize)..]
            .iter()
            .cloned()
            .collect();

        let req = AppendEntriesRequest {
            term: my_term,
            leader_id,
            prev_log_index: prev_index,
            prev_log_term: prev_term,
            entries,
            leader_commit: commit_index,
        };

        // 네트워크를 통해 전송 (비동기, 지연 포함)
        let network_clone = network.clone();
        let next_index_clone = next_index.clone();

        tokio::spawn(async move {
            if let Some(response) = network_clone.send_append_entries(i, req).await {
                // 응답 처리
                let mut leader_lock = network_clone.nodes[leader_id].write().await;

                if response.term > leader_lock.current_term {
                    leader_lock.current_term = response.term;
                    leader_lock.state = NodeState::Follower;
                } else if response.success {
                    leader_lock.match_index[i] = response.match_index;
                    leader_lock.next_index[i] = response.match_index + 1;
                } else {
                    // 복제 실패: next_index 감소
                    leader_lock.next_index[i] = std::cmp::max(1, next_index_clone[i] - 1);
                }
            }
        });
    }
}
```

---

### Step 4: 네트워크 시뮬레이션 (1일)

#### 4-1. 네트워크 지연 추가
```rust
pub async fn send_request_vote(&self, to: usize, req: RequestVoteRequest) -> RequestVoteResponse {
    // 지연 시뮬레이션
    let latency = self.latencies.get(&(req.candidate_id, to))
        .map(|r| r.clone())
        .unwrap_or(0);

    if latency > 0 {
        tokio::time::sleep(Duration::from_millis(latency)).await;
    }

    // 메시지 손실 시뮬레이션
    let loss_rate = self.message_loss.get(&(req.candidate_id, to))
        .map(|r| r.clone())
        .unwrap_or(0.0);

    if rand::random::<f64>() < loss_rate {
        // 메시지 손실됨 → 타임아웃 발생
        return RequestVoteResponse {
            term: req.term,
            vote_granted: false,
        };
    }

    // 실제 처리
    let node = &self.nodes[to];
    let mut node_lock = node.write().await;
    node_lock.request_vote(req).await
}
```

#### 4-2. 장애 주입
```rust
pub fn set_network_latency(&self, from: usize, to: usize, latency_ms: u64) {
    self.latencies.insert((from, to), latency_ms);
}

pub fn set_message_loss(&self, from: usize, to: usize, loss_rate: f64) {
    self.message_loss.insert((from, to), loss_rate);
}

pub fn create_partition(&self, partition1: Vec<usize>, partition2: Vec<usize>) {
    // 두 그룹 간 네트워크 단절
    for i in partition1.iter() {
        for j in partition2.iter() {
            self.set_message_loss(*i, *j, 1.0);  // 100% 손실
            self.set_message_loss(*j, *i, 1.0);
        }
    }
}
```

---

### Step 5: 테스트 케이스 (1일)

#### 5-1. 정상 작동 테스트
```rust
#[tokio::test]
async fn test_normal_election() {
    let network = RaftNetwork::new(5);

    // 모든 노드 시작
    network.start().await;

    // 첫 번째 선출: 임의의 노드가 리더가 됨
    tokio::time::sleep(Duration::from_secs(1)).await;

    let status = network.get_status().await;
    assert_eq!(status.leaders, 1, "정확히 1개 리더 필요");
    assert!(status.leader_id.is_some());
}
```

#### 5-2. 리더 장애 테스트
```rust
#[tokio::test]
async fn test_leader_failure() {
    let network = RaftNetwork::new(5);
    network.start().await;

    tokio::time::sleep(Duration::from_secs(1)).await;
    let old_leader = network.get_status().await.leader_id;

    // 리더 시뮬레이션 중단
    network.stop_node(old_leader.unwrap()).await;

    // 새로운 리더 선출
    tokio::time::sleep(Duration::from_secs(2)).await;

    let status = network.get_status().await;
    assert_ne!(status.leader_id, old_leader, "새로운 리더 선출됨");
    assert_eq!(status.leaders, 1, "정확히 1개 리더");
}
```

#### 5-3. 네트워크 분할 테스트
```rust
#[tokio::test]
async fn test_network_partition() {
    let network = RaftNetwork::new(5);
    network.start().await;

    // 3-2 분할: [0,1,2] vs [3,4]
    network.create_partition(vec![0, 1, 2], vec![3, 4]).await;

    tokio::time::sleep(Duration::from_secs(2)).await;

    let status = network.get_status().await;
    // 큰 파티션(3개)에서만 리더 선출 가능 (quorum = 3)
    assert_eq!(status.leaders, 1, "큰 파티션에서 리더 1개");
}
```

#### 5-4. Log Replication 테스트
```rust
#[tokio::test]
async fn test_log_replication() {
    let network = RaftNetwork::new(5);
    network.start().await;

    tokio::time::sleep(Duration::from_secs(1)).await;
    let leader = network.get_status().await.leader_id.unwrap();

    // 리더에 로그 추가
    network.append_entry(leader, "entry1".to_string()).await;

    // 복제 대기
    tokio::time::sleep(Duration::from_millis(500)).await;

    // 모든 노드에 로그가 복제되었는지 확인
    for i in 0..5 {
        let node = &network.nodes[i];
        let node_lock = node.read().await;
        assert!(node_lock.log.len() > 0, "노드 {}에 로그 복제됨", i);
    }
}
```

---

## 📊 Phase 2 검증 기준

| 항목 | 검증 방법 | 성공 기준 |
|------|---------|---------|
| **리더 선출** | test_normal_election | 1개 리더, 다른 노드에서도 선출 |
| **장애 복구** | test_leader_failure | 새로운 리더 < 1초 |
| **네트워크 분할** | test_network_partition | quorum 유지하는 파티션만 선출 |
| **로그 복제** | test_log_replication | 모든 노드에 복제 |
| **Quorum 유지** | 통합 테스트 | 2개 노드 장애에서도 생존 |

---

## 📝 성공 기준

### Phase 2 완료 조건 (모두 필수)
```
✅ cargo build --release 성공
✅ cargo test: 모든 테스트 통과
✅ 3개 이상 노드에서 리더 선출 가능
✅ 노드 0 이외의 노드도 리더 됨
✅ 리더 장애 → 새 리더 선출 < 1초
✅ 로그 일관성 보장
✅ Quorum 기반 안전성
✅ 문서: RFC 5740 검증 보고서
```

---

## 🚀 구현 로드맵

```
Day 1: Raft 구조 정의 + RequestVote
       └─ RaftNode 확장, RPC 메시지, 네트워크 정의

Day 2: 리더 선출 + 타임아웃
       └─ 투표 수집, Election Timeout, Quorum 로직

Day 3: Log Replication
       └─ AppendEntries RPC, 로그 복제, Commit

Day 4: 네트워크 시뮬레이션 + 테스트
       └─ 지연, 손실, 네트워크 분할 테스트

Day 5: 최종 검증 + 문서화
       └─ 모든 테스트 통과, 보고서 작성
```

---

## 📚 참고 자료

- **RFC 5740**: https://raft.github.io/
- **Raft Visualization**: https://raft.github.io/raftscope/index.html
- **etcd Raft**: https://github.com/etcd-io/etcd/tree/main/raft
- **Consul**: https://github.com/hashicorp/raft

---

## 🎯 최종 목표

```
"기록이 증명이다"

Phase 2 완료 후:

✅ Raft 구현: RFC 5740 준수
✅ 검증: "정말 동작한다" 증명
✅ 커밋: 5-6개 (체계적 기록)
✅ 테스트: 모든 시나리오 통과
✅ 신뢰도: 85% → 95%

"10년 무중단" 준비 완료
```

---

**시작**: 2026-03-02 (지금부터)
**기한**: 2026-03-07 (5일)
**다음**: Phase 3 (카오스 엔지니어링)
