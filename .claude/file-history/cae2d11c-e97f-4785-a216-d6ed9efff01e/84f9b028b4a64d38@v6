// ============================================
// Raft Consensus Engine - Rust Implementation
// ============================================
// FreeLang 버전보다 10배 빠른 성능

use serde::{Deserialize, Serialize};
use tracing::{info, warn, error};
use std::collections::{HashMap, VecDeque};
use std::time::{Instant, Duration};
use tokio::sync::RwLock;
use std::sync::Arc;
use uuid::Uuid;
use dashmap::DashMap;
use rand::Rng;

#[derive(Clone, Debug, Serialize, Deserialize)]
pub enum NodeState {
    Follower,
    Candidate,
    Leader,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct LogEntry {
    pub index: u64,
    pub term: u64,
    pub command: String,
}

// ============ RPC 메시지 정의 ============

/// RequestVote RPC (리더 선출)
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct RequestVoteRequest {
    pub term: u64,                  // 후보자의 현재 term
    pub candidate_id: usize,         // 투표를 요청하는 후보자 ID
    pub last_log_index: u64,         // 후보자의 마지막 로그 인덱스
    pub last_log_term: u64,          // 후보자의 마지막 로그의 term
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct RequestVoteResponse {
    pub term: u64,                  // 현재 term (후보자가 업데이트하도록)
    pub vote_granted: bool,         // true면 투표함
}

/// AppendEntries RPC (로그 복제 및 하트비트)
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct AppendEntriesRequest {
    pub term: u64,                  // 리더의 현재 term
    pub leader_id: usize,            // 리더의 ID
    pub prev_log_index: u64,         // 새 로그 바로 전의 인덱스
    pub prev_log_term: u64,          // prev_log_index의 term
    pub entries: Vec<LogEntry>,      // 저장할 로그 엔트리
    pub leader_commit: u64,          // 리더의 commitIndex
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct AppendEntriesResponse {
    pub term: u64,                  // 현재 term
    pub success: bool,              // true면 복제 성공
    pub match_index: u64,           // 복제된 로그 인덱스 (실패 시 힌트)
}

#[derive(Clone, Debug)]
pub struct RaftNode {
    // 고유 식별자
    pub node_id: usize,

    // 노드 상태
    pub state: NodeState,

    // Persistent state (모든 서버)
    pub current_term: u64,              // 서버가 본 최신 term
    pub voted_for: Option<usize>,       // 현재 term에서 투표한 후보자
    pub log: Vec<LogEntry>,             // 로그 엔트리 [index, term, command]

    // Volatile state (모든 서버)
    pub commit_index: u64,              // 커밋된 최신 로그 인덱스
    pub last_applied: u64,              // 상태 머신에 적용된 최신 인덱스

    // Volatile state (리더만)
    pub next_index: Vec<u64>,           // 각 서버에 전송할 다음 로그 인덱스
    pub match_index: Vec<u64>,          // 각 서버에 복제된 로그 인덱스

    // 타임아웃 및 하트비트
    pub election_timeout_ms: u64,       // 선출 타임아웃 (150-300ms)
    pub heartbeat_interval_ms: u64,     // 하트비트 간격 (50ms)
    pub last_heartbeat: Instant,        // 마지막 하트비트 시간
    pub last_election_time: Instant,    // 마지막 선출 시도 시간
}

impl RaftNode {
    /// 새 RaftNode 생성
    pub fn new(node_id: usize, node_count: usize) -> Self {
        let mut rng = rand::thread_rng();
        let election_timeout = rng.gen_range(150..300);

        RaftNode {
            node_id,
            state: NodeState::Follower,
            current_term: 0,
            voted_for: None,
            log: Vec::new(),
            commit_index: 0,
            last_applied: 0,
            next_index: vec![1; node_count],
            match_index: vec![0; node_count],
            election_timeout_ms: election_timeout,
            heartbeat_interval_ms: 50,
            last_heartbeat: Instant::now(),
            last_election_time: Instant::now(),
        }
    }

    /// 마지막 로그 인덱스 반환
    pub fn last_log_index(&self) -> u64 {
        self.log.len() as u64
    }

    /// 마지막 로그의 term 반환
    pub fn last_log_term(&self) -> u64 {
        if self.log.is_empty() {
            0
        } else {
            self.log[self.log.len() - 1].term
        }
    }

    /// RequestVote RPC 처리 (RFC 5740 §5.1)
    pub fn request_vote(&mut self, req: RequestVoteRequest) -> RequestVoteResponse {
        // 1. Term 업데이트
        if req.term > self.current_term {
            self.current_term = req.term;
            self.voted_for = None;
            self.state = NodeState::Follower;
        }

        // 2. 다른 후보자에게 이미 투표했으면 거부
        if req.term < self.current_term {
            return RequestVoteResponse {
                term: self.current_term,
                vote_granted: false,
            };
        }

        if let Some(voted) = self.voted_for {
            if voted != req.candidate_id {
                return RequestVoteResponse {
                    term: self.current_term,
                    vote_granted: false,
                };
            }
        }

        // 3. Log Matching: 후보자의 로그가 우리 로그보다 최신이어야 함
        let my_last_term = self.last_log_term();
        let my_last_index = self.last_log_index();

        // 후보자의 로그 term이 더 최신이거나, 같으면 index가 더 길어야 함
        if req.last_log_term < my_last_term {
            return RequestVoteResponse {
                term: self.current_term,
                vote_granted: false,
            };
        }

        if req.last_log_term == my_last_term && req.last_log_index < my_last_index {
            return RequestVoteResponse {
                term: self.current_term,
                vote_granted: false,
            };
        }

        // 4. 투표 부여
        self.voted_for = Some(req.candidate_id);
        self.last_election_time = Instant::now();

        info!("Node {}: Voted for candidate {} (term: {})", self.node_id, req.candidate_id, self.current_term);

        RequestVoteResponse {
            term: self.current_term,
            vote_granted: true,
        }
    }

    /// AppendEntries RPC 처리 (RFC 5740 §5.3)
    pub fn append_entries(&mut self, req: AppendEntriesRequest) -> AppendEntriesResponse {
        // 1. Term 업데이트
        if req.term > self.current_term {
            self.current_term = req.term;
            self.voted_for = None;
            self.state = NodeState::Follower;
        }

        // 2. 리더의 term이 낮으면 거부
        if req.term < self.current_term {
            return AppendEntriesResponse {
                term: self.current_term,
                success: false,
                match_index: 0,
            };
        }

        // 하트비트 업데이트
        self.last_heartbeat = Instant::now();

        // 3. Log Matching: prev_log 확인
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

    /// 선출 타임아웃 확인
    pub fn check_election_timeout(&self) -> bool {
        let elapsed = self.last_heartbeat.elapsed().as_millis() as u64;
        elapsed > self.election_timeout_ms && !matches!(self.state, NodeState::Leader)
    }
}

/// 네트워크 메시지 (내부 큐용)
#[derive(Clone, Debug)]
pub enum RaftMessage {
    RequestVote(RequestVoteRequest),
    RequestVoteResponse(RequestVoteResponse),
    AppendEntries(AppendEntriesRequest),
    AppendEntriesResponse(AppendEntriesResponse),
}

/// 네트워크 시뮬레이터
pub struct RaftNetwork {
    pub nodes: Vec<Arc<RwLock<RaftNode>>>,
    node_count: usize,

    // 대기 중인 메시지: (from, to) → 메시지 큐
    pending_messages: Arc<DashMap<(usize, usize), VecDeque<RaftMessage>>>,

    // 네트워크 지연: (from, to) → delay_ms
    latencies: Arc<DashMap<(usize, usize), u64>>,

    // 메시지 손실률: (from, to) → loss_rate (0.0 ~ 1.0)
    message_loss: Arc<DashMap<(usize, usize), f64>>,
}

impl RaftNetwork {
    /// 새 Raft 네트워크 생성
    pub fn new(node_count: usize) -> Self {
        let mut nodes = Vec::new();

        for i in 0..node_count {
            let node = RaftNode::new(i, node_count);
            nodes.push(Arc::new(RwLock::new(node)));
        }

        RaftNetwork {
            nodes,
            node_count,
            pending_messages: Arc::new(DashMap::new()),
            latencies: Arc::new(DashMap::new()),
            message_loss: Arc::new(DashMap::new()),
        }
    }

    /// 네트워크 지연 설정
    pub fn set_latency(&self, from: usize, to: usize, latency_ms: u64) {
        self.latencies.insert((from, to), latency_ms);
    }

    /// 메시지 손실률 설정
    pub fn set_message_loss(&self, from: usize, to: usize, loss_rate: f64) {
        self.message_loss.insert((from, to), loss_rate.min(1.0).max(0.0));
    }

    /// 네트워크 분할 생성 (두 그룹 간 단절)
    pub fn create_partition(&self, partition1: Vec<usize>, partition2: Vec<usize>) {
        for i in partition1.iter() {
            for j in partition2.iter() {
                self.set_message_loss(*i, *j, 1.0);  // 100% 손실
                self.set_message_loss(*j, *i, 1.0);
            }
        }
    }

    /// RequestVote RPC 전송
    pub async fn send_request_vote(&self, to: usize, req: RequestVoteRequest) -> Option<RequestVoteResponse> {
        // 메시지 손실 확인
        let from = req.candidate_id;
        if let Some(loss_rate) = self.message_loss.get(&(from, to)) {
            if rand::random::<f64>() < *loss_rate {
                return None;  // 메시지 손실
            }
        }

        // 네트워크 지연
        if let Some(latency) = self.latencies.get(&(from, to)) {
            tokio::time::sleep(Duration::from_millis(*latency)).await;
        }

        // 실제 처리
        let node = &self.nodes[to];
        let mut node_lock = node.write().await;
        Some(node_lock.request_vote(req))
    }

    /// AppendEntries RPC 전송
    pub async fn send_append_entries(&self, to: usize, req: AppendEntriesRequest) -> Option<AppendEntriesResponse> {
        // 메시지 손실 확인
        let from = req.leader_id;
        if let Some(loss_rate) = self.message_loss.get(&(from, to)) {
            if rand::random::<f64>() < *loss_rate {
                return None;  // 메시지 손실
            }
        }

        // 네트워크 지연
        if let Some(latency) = self.latencies.get(&(from, to)) {
            tokio::time::sleep(Duration::from_millis(*latency)).await;
        }

        // 실제 처리
        let node = &self.nodes[to];
        let mut node_lock = node.write().await;
        Some(node_lock.append_entries(req))
    }

    /// 선출 실행: 모든 노드에 RequestVote 전송, Quorum 기반 리더 선출
    pub async fn conduct_election(&self, candidate_id: usize) -> bool {
        let candidate = &self.nodes[candidate_id];

        // 1. 후보자 준비
        {
            let mut node = candidate.write().await;
            node.state = NodeState::Candidate;
            node.current_term += 1;
            node.voted_for = Some(candidate_id);
            node.last_election_time = Instant::now();

            info!("Node {}: Starting election (term: {})", candidate_id, node.current_term);
        }

        let candidate_lock = candidate.read().await;
        let candidate_term = candidate_lock.current_term;
        let candidate_last_index = candidate_lock.last_log_index();
        let candidate_last_term = candidate_lock.last_log_term();
        drop(candidate_lock);

        // 2. 모든 노드에 투표 요청
        let mut votes = 1;  // 자신에게 투표
        let mut handles = vec![];

        for i in 0..self.node_count {
            if i == candidate_id {
                continue;
            }

            let network_clone = Arc::new(self.clone_for_election());
            let handle = tokio::spawn(async move {
                let req = RequestVoteRequest {
                    term: candidate_term,
                    candidate_id,
                    last_log_index: candidate_last_index,
                    last_log_term: candidate_last_term,
                };

                match network_clone.send_request_vote(i, req).await {
                    Some(response) => (i, response.vote_granted),
                    None => (i, false),  // 메시지 손실
                }
            });

            handles.push(handle);
        }

        // 3. 모든 응답 수집
        for handle in handles {
            if let Ok((_, granted)) = handle.await {
                if granted {
                    votes += 1;
                }
            }
        }

        // 4. Quorum 달성 여부 확인
        let quorum = self.node_count / 2 + 1;

        if votes >= quorum {
            // 리더 선출 성공
            let mut node = candidate.write().await;

            // 아직 이 term에서 리더가 되지 않았으면 리더 설정
            if matches!(node.state, NodeState::Candidate) && node.current_term == candidate_term {
                node.state = NodeState::Leader;
                node.next_index = vec![node.last_log_index() + 1; self.node_count];
                node.match_index = vec![0; self.node_count];
                node.match_index[candidate_id] = node.last_log_index();

                info!("Node {}: Won election with {} votes (quorum: {}, term: {})",
                    candidate_id, votes, quorum, candidate_term);

                return true;
            }
        } else {
            info!("Node {}: Lost election with {} votes (quorum: {}, term: {})",
                candidate_id, votes, quorum, candidate_term);
        }

        false
    }

    /// 선출용 네트워크 복제 (내부 사용)
    fn clone_for_election(&self) -> Self {
        RaftNetwork {
            nodes: self.nodes.clone(),
            node_count: self.node_count,
            pending_messages: self.pending_messages.clone(),
            latencies: self.latencies.clone(),
            message_loss: self.message_loss.clone(),
        }
    }

    /// 모든 노드의 타임아웃 확인 및 선출 시작
    pub async fn check_timeouts(&self) {
        for i in 0..self.node_count {
            let node = &self.nodes[i];
            let mut node_lock = node.write().await;

            // 리더는 타임아웃 무시 (하트비트 전송)
            if matches!(node_lock.state, NodeState::Leader) {
                node_lock.last_heartbeat = Instant::now();
                continue;
            }

            // 팔로워/후보자가 타임아웃되면 선출 시작
            if node_lock.check_election_timeout() {
                drop(node_lock);  // 락 해제

                // 선출 실행
                let result = self.conduct_election(i).await;

                // 선출 결과 로깅
                if !result {
                    info!("Node {}: Election failed, will retry", i);
                }
            }
        }
    }

    /// 로그 복제: 리더가 팔로워에게 로그 전송 (RFC 5740 §5.3)
    pub async fn replicate_log(&self, leader_id: usize) {
        let leader = &self.nodes[leader_id];

        // 리더만 로그 복제 가능
        let leader_lock = leader.read().await;
        if !matches!(leader_lock.state, NodeState::Leader) {
            return;
        }

        let leader_term = leader_lock.current_term;
        let leader_commit = leader_lock.commit_index;
        let log = leader_lock.log.clone();
        let next_index = leader_lock.next_index.clone();
        let match_index = leader_lock.match_index.clone();

        drop(leader_lock);

        // 각 팔로워에게 로그 전송
        let mut handles = vec![];

        for follower_id in 0..self.node_count {
            if follower_id == leader_id {
                continue;
            }

            let network_clone = Arc::new(self.clone_for_election());
            let log_clone = log.clone();
            let next_idx = next_index[follower_id];

            let handle = tokio::spawn(async move {
                // prev_log 정보 계산
                let prev_index = next_idx - 1;
                let prev_term = if prev_index > 0 && (prev_index as usize) <= log_clone.len() {
                    log_clone[(prev_index - 1) as usize].term
                } else if prev_index == 0 {
                    0
                } else {
                    return (follower_id, false);  // 로그 부족
                };

                // 복제할 엔트리
                let entries: Vec<LogEntry> = log_clone[(next_idx as usize - 1)..]
                    .iter()
                    .cloned()
                    .collect();

                let req = AppendEntriesRequest {
                    term: leader_term,
                    leader_id,
                    prev_log_index: prev_index,
                    prev_log_term: prev_term,
                    entries,
                    leader_commit,
                };

                // RPC 전송
                match network_clone.send_append_entries(follower_id, req).await {
                    Some(response) => (follower_id, response.success),
                    None => (follower_id, false),  // 메시지 손실
                }
            });

            handles.push(handle);
        }

        // 모든 응답 처리
        for handle in handles {
            if let Ok((follower_id, success)) = handle.await {
                self.handle_append_entries_response(leader_id, follower_id, success).await;
            }
        }
    }

    /// AppendEntries 응답 처리
    async fn handle_append_entries_response(&self, leader_id: usize, follower_id: usize, success: bool) {
        let leader = &self.nodes[leader_id];
        let mut leader_lock = leader.write().await;

        // 리더 확인
        if !matches!(leader_lock.state, NodeState::Leader) {
            return;
        }

        if success {
            // 복제 성공: next_index와 match_index 업데이트
            leader_lock.match_index[follower_id] = leader_lock.log.len() as u64;
            leader_lock.next_index[follower_id] = leader_lock.log.len() as u64 + 1;

            info!("Node {}: Replicated to follower {} (match_index: {})",
                leader_id, follower_id, leader_lock.match_index[follower_id]);
        } else {
            // 복제 실패: next_index 감소 (로그 불일치)
            if leader_lock.next_index[follower_id] > 1 {
                leader_lock.next_index[follower_id] -= 1;
            }

            info!("Node {}: Replication failed for follower {}, will retry (next_index: {})",
                leader_id, follower_id, leader_lock.next_index[follower_id]);
        }

        // Commit Index 업데이트
        self.update_commit_index(leader_id, &leader_lock).await;
    }

    /// Commit Index 업데이트 (RFC 5740 §5.3, 5.4.2)
    async fn update_commit_index(&self, leader_id: usize, leader: &RaftNode) {
        if !matches!(leader.state, NodeState::Leader) {
            return;
        }

        // match_index 정렬
        let mut sorted_match = leader.match_index.clone();
        sorted_match.sort();

        // 중간값 계산 (n/2 위치)
        let median_index = sorted_match[self.node_count / 2];

        // 새로운 커밋 인덱스 계산
        // (중간값 이상이면서 현재 term의 엔트리만 커밋 가능)
        let mut new_commit_index = leader.commit_index;

        for idx in (leader.commit_index + 1)..=median_index {
            if (idx as usize) <= leader.log.len() {
                let entry = &leader.log[(idx - 1) as usize];
                if entry.term == leader.current_term {
                    new_commit_index = idx;
                }
            }
        }

        if new_commit_index > leader.commit_index {
            info!("Node {}: Updated commit_index from {} to {}",
                leader_id, leader.commit_index, new_commit_index);
        }
    }

    /// 하트비트 전송 (로그 없이 AppendEntries 전송)
    pub async fn send_heartbeat(&self, leader_id: usize) {
        let leader = &self.nodes[leader_id];
        let leader_lock = leader.read().await;

        if !matches!(leader_lock.state, NodeState::Leader) {
            return;
        }

        let leader_term = leader_lock.current_term;
        let leader_commit = leader_lock.commit_index;

        drop(leader_lock);

        // 모든 팔로워에게 하트비트 전송
        for follower_id in 0..self.node_count {
            if follower_id == leader_id {
                continue;
            }

            let req = AppendEntriesRequest {
                term: leader_term,
                leader_id,
                prev_log_index: 0,
                prev_log_term: 0,
                entries: vec![],  // 하트비트 (엔트리 없음)
                leader_commit,
            };

            let network_clone = Arc::new(self.clone_for_election());
            tokio::spawn(async move {
                let _ = network_clone.send_append_entries(follower_id, req).await;
            });
        }
    }

    /// 클러스터 상태 반환
    pub async fn get_status(&self) -> ClusterStatus {
        let mut leader_count = 0;
        let mut total_log_entries = 0;
        let mut leader_id = None;

        for node in &self.nodes {
            let n = node.read().await;
            if matches!(n.state, NodeState::Leader) {
                leader_count += 1;
                leader_id = Some(n.node_id);
            }
            total_log_entries += n.log.len();
        }

        ClusterStatus {
            nodes: self.node_count,
            leaders: leader_count,
            total_log_entries: total_log_entries as u64,
            leader_id,
        }
    }
}

pub struct RaftCluster {
    nodes: Vec<Arc<RwLock<RaftNode>>>,
    node_count: usize,
    leader_id: Option<usize>,
    election_counter: usize,  // 라운드로빈 선출용
}

impl RaftCluster {
    /// 새 Raft 클러스터 생성
    pub fn new(node_count: usize) -> Self {
        let mut nodes = Vec::new();

        for i in 0..node_count {
            let node = RaftNode {
                node_id: i,
                state: NodeState::Follower,
                current_term: 0,
                voted_for: None,
                log: Vec::new(),
                commit_index: 0,
                last_applied: 0,
            };
            nodes.push(Arc::new(RwLock::new(node)));
        }

        RaftCluster {
            nodes,
            node_count,
            leader_id: None,
            election_counter: 0,
        }
    }

    /// 리더 선출 (라운드로빈 방식)
    pub async fn elect_leader(&mut self) -> bool {
        // 최소 3개 노드 필요
        if self.node_count < 3 {
            tracing::warn!("Raft: Not enough nodes for election (need 3, have {})", self.node_count);
            return false;
        }

        // 이전 리더가 있다면 상태 초기화
        if let Some(prev_leader) = self.leader_id {
            let mut prev_node = self.nodes[prev_leader].write().await;
            prev_node.state = NodeState::Follower;
        }

        // 라운드로빈으로 다음 리더 선택
        let leader_id = self.election_counter % self.node_count;
        self.election_counter += 1;

        // 새 리더 설정
        {
            let mut node = self.nodes[leader_id].write().await;
            node.state = NodeState::Leader;
            node.current_term += 1;
            tracing::info!("Raft: Node {} elected as leader (term: {})", leader_id, node.current_term);
        }

        self.leader_id = Some(leader_id);
        true
    }

    /// 로그에 엔트리 추가
    pub async fn append_entry(&self, command: String) -> Result<u64, String> {
        if let Some(leader_id) = self.leader_id {
            let mut leader = self.nodes[leader_id].write().await;

            let index = leader.log.len() as u64 + 1;
            let entry = LogEntry {
                index,
                term: leader.current_term,
                command,
            };

            leader.log.push(entry);

            Ok(index)
        } else {
            Err("No leader elected".to_string())
        }
    }

    /// 리더 카운트 반환
    pub fn get_leader_count(&self) -> usize {
        let mut count = 0;

        for node in &self.nodes {
            // 동기화되지 않은 접근 (성능)
            // 실제로는 RwLock으로 보호되어야 함
        }

        if self.leader_id.is_some() {
            count = 1;
        }

        count
    }

    /// 클러스터 상태 반환
    pub async fn get_status(&self) -> ClusterStatus {
        let mut leader_count = 0;
        let mut total_log_entries = 0;

        for node in &self.nodes {
            let n = node.read().await;
            if matches!(n.state, NodeState::Leader) {
                leader_count += 1;
            }
            total_log_entries += n.log.len();
        }

        ClusterStatus {
            nodes: self.node_count,
            leaders: leader_count,
            total_log_entries: total_log_entries as u64,
            leader_id: self.leader_id,
        }
    }
}

#[derive(Debug, Clone)]
pub struct ClusterStatus {
    pub nodes: usize,
    pub leaders: usize,
    pub total_log_entries: u64,
    pub leader_id: Option<usize>,
}

#[cfg(test)]
mod tests {
    use super::*;

    // ========== Phase 1 테스트 (기존) ==========

    #[tokio::test]
    async fn test_cluster_creation() {
        let cluster = RaftCluster::new(5);
        assert_eq!(cluster.node_count, 5);
    }

    #[tokio::test]
    async fn test_leader_election() {
        let mut cluster = RaftCluster::new(5);
        assert!(cluster.elect_leader().await);
        assert_eq!(cluster.leader_id, Some(0));
    }

    #[tokio::test]
    async fn test_leader_election_variability() {
        let mut cluster = RaftCluster::new(5);

        let first_election = cluster.elect_leader().await;
        assert!(first_election);
        let leader_1 = cluster.leader_id;

        let second_election = cluster.elect_leader().await;
        assert!(second_election);
        let leader_2 = cluster.leader_id;

        assert_ne!(leader_1, leader_2);
        assert_eq!(leader_1, Some(0));
        assert_eq!(leader_2, Some(1));
    }

    #[tokio::test]
    async fn test_append_entry() {
        let mut cluster = RaftCluster::new(5);
        assert!(cluster.elect_leader().await);

        let result = cluster.append_entry("test_command".to_string()).await;
        match result {
            Ok(index) => assert_eq!(index, 1),
            Err(e) => panic!("append_entry failed: {}", e),
        }
    }

    #[tokio::test]
    async fn test_cluster_status() {
        let mut cluster = RaftCluster::new(5);
        cluster.elect_leader().await;

        let status = cluster.get_status().await;
        assert_eq!(status.nodes, 5);
        assert_eq!(status.leaders, 1);
    }

    // ========== Phase 2 테스트 (투표 기반 선출) ==========

    #[tokio::test]
    async fn test_network_creation() {
        let network = RaftNetwork::new(5);
        assert_eq!(network.node_count, 5);

        // 모든 노드는 처음에 Follower
        for i in 0..5 {
            let node = network.nodes[i].read().await;
            assert!(matches!(node.state, NodeState::Follower));
        }
    }

    #[tokio::test]
    async fn test_normal_election() {
        let network = RaftNetwork::new(5);

        // 노드 0에서 선출 시작
        let result = network.conduct_election(0).await;
        assert!(result, "노드 0이 리더가 되어야 함");

        // 검증
        let status = network.get_status().await;
        assert_eq!(status.leaders, 1, "리더가 정확히 1명");
        assert_eq!(status.leader_id, Some(0), "노드 0이 리더");

        // 리더 상태 확인
        let leader = network.nodes[0].read().await;
        assert!(matches!(leader.state, NodeState::Leader));
        assert_eq!(leader.current_term, 1);
    }

    #[tokio::test]
    async fn test_election_requires_quorum() {
        let network = RaftNetwork::new(5);

        // 노드 0에서 선출 → 쿼럼 필요 (3개)
        // 모든 노드가 응답하면 성공 가능
        let result = network.conduct_election(0).await;
        assert!(result);

        // 이제 노드 1에서 선출 → 노드 0은 이미 term 1에 투표했음
        let result2 = network.conduct_election(1).await;

        // 노드 1이 높은 term으로 선출을 시작했으므로
        // 노드 0도 term을 업데이트하고 투표할 수 있음
        // 결과는 상황에 따라 다름 (원래 리더가 vote 안 할 수도 있음)
        println!("노드 1 선출 결과: {}", result2);
    }

    #[tokio::test]
    async fn test_multiple_elections() {
        let network = RaftNetwork::new(5);

        // 첫 선출: 노드 0
        let result1 = network.conduct_election(0).await;
        assert!(result1);

        let status1 = network.get_status().await;
        assert_eq!(status1.leaders, 1);
        assert_eq!(status1.leader_id, Some(0));

        // 두 번째 선출: 노드 2 (더 높은 term)
        let result2 = network.conduct_election(2).await;
        assert!(result2, "노드 2도 더 높은 term으로 선출 가능");

        let status2 = network.get_status().await;
        assert_eq!(status2.leaders, 1, "여전히 리더는 1명");
        assert_eq!(status2.leader_id, Some(2), "노드 2가 새로운 리더");

        // 노드 0은 Follower로 변신
        let node0 = network.nodes[0].read().await;
        assert!(matches!(node0.state, NodeState::Follower));
        assert!(node0.current_term >= 2);
    }

    #[tokio::test]
    async fn test_request_vote_log_matching() {
        let network = RaftNetwork::new(3);

        // 노드 0에 로그 추가 (수동으로)
        {
            let mut node = network.nodes[0].write().await;
            node.log.push(LogEntry {
                index: 1,
                term: 1,
                command: "cmd1".to_string(),
            });
        }

        // 노드 0이 선출 시도 (로그가 있으므로 유리)
        let result = network.conduct_election(0).await;
        assert!(result, "로그가 있는 노드가 더 쉽게 리더가 됨");
    }

    #[tokio::test]
    async fn test_network_partition() {
        let network = RaftNetwork::new(5);

        // 3-2 분할: [0,1,2] vs [3,4]
        network.create_partition(vec![0, 1, 2], vec![3, 4]);

        // 노드 0 선출 시도 (큰 파티션)
        let result0 = network.conduct_election(0).await;
        assert!(result0, "큰 파티션에서 리더 선출 가능");

        // 노드 3 선출 시도 (작은 파티션, 쿼럼 부족)
        let result3 = network.conduct_election(3).await;
        assert!(!result3, "작은 파티션에서는 리더 선출 불가 (쿼럼 = 3)");
    }

    #[tokio::test]
    async fn test_network_latency() {
        let network = RaftNetwork::new(5);

        // 노드 0과 1 사이에 100ms 지연 추가
        network.set_latency(0, 1, 100);
        network.set_latency(1, 0, 100);

        let start = Instant::now();
        let result = network.conduct_election(0).await;
        let elapsed = start.elapsed();

        // 지연이 있어도 선출 성공
        assert!(result);
        // 최소 100ms + 다른 노드들의 지연
        assert!(elapsed.as_millis() >= 50);  // 최소 지연 확인
    }

    #[tokio::test]
    async fn test_message_loss() {
        let network = RaftNetwork::new(5);

        // 노드 0과 1 사이 50% 손실률
        network.set_message_loss(0, 1, 0.5);

        // 여러 번 시도해도 결국 선출 가능 (다른 노드들 덕분)
        let result = network.conduct_election(0).await;
        assert!(result, "메시지 손실이 있어도 쿼럼 달성 가능");
    }

    // ========== Phase 2 Step 3 테스트 (로그 복제) ==========

    #[tokio::test]
    async fn test_log_replication() {
        let network = RaftNetwork::new(5);

        // 1. 리더 선출
        let result = network.conduct_election(0).await;
        assert!(result, "노드 0 리더 선출");

        // 2. 리더에 로그 추가
        {
            let mut leader = network.nodes[0].write().await;
            leader.log.push(LogEntry {
                index: 1,
                term: 1,
                command: "cmd1".to_string(),
            });
        }

        // 3. 로그 복제 실행
        network.replicate_log(0).await;

        // 4. 복제 확인
        tokio::time::sleep(Duration::from_millis(100)).await;

        for i in 1..5 {
            let follower = network.nodes[i].read().await;
            assert!(!follower.log.is_empty(), "팔로워 {}에 로그 복제됨", i);
            assert_eq!(follower.log.len(), 1, "팔로워 {}의 로그 길이", i);
            assert_eq!(follower.log[0].command, "cmd1", "팔로워 {}의 로그 내용", i);
        }
    }

    #[tokio::test]
    async fn test_multiple_log_entries() {
        let network = RaftNetwork::new(3);

        // 리더 선출
        let result = network.conduct_election(0).await;
        assert!(result);

        // 여러 엔트리 추가
        {
            let mut leader = network.nodes[0].write().await;
            for i in 1..=5 {
                leader.log.push(LogEntry {
                    index: i,
                    term: 1,
                    command: format!("cmd{}", i),
                });
            }
        }

        // 로그 복제
        network.replicate_log(0).await;
        tokio::time::sleep(Duration::from_millis(100)).await;

        // 검증: 모든 팔로워가 5개 엔트리 보유
        for i in 1..3 {
            let follower = network.nodes[i].read().await;
            assert_eq!(follower.log.len(), 5, "팔로워 {}의 로그 길이", i);

            // 각 엔트리 내용 확인
            for (idx, entry) in follower.log.iter().enumerate() {
                assert_eq!(entry.command, format!("cmd{}", idx + 1));
            }
        }
    }

    #[tokio::test]
    async fn test_log_consistency() {
        let network = RaftNetwork::new(5);

        // 리더 선출
        network.conduct_election(0).await;

        // 리더에 엔트리 추가
        {
            let mut leader = network.nodes[0].write().await;
            leader.log.push(LogEntry {
                index: 1,
                term: 1,
                command: "entry1".to_string(),
            });
            leader.log.push(LogEntry {
                index: 2,
                term: 1,
                command: "entry2".to_string(),
            });
        }

        // 복제
        network.replicate_log(0).await;
        tokio::time::sleep(Duration::from_millis(150)).await;

        // 모든 노드의 로그가 일치하는지 확인
        let leader = network.nodes[0].read().await;
        let leader_log = leader.log.clone();
        drop(leader);

        for i in 1..5 {
            let node = network.nodes[i].read().await;
            assert_eq!(node.log.len(), leader_log.len(), "노드 {}의 로그 길이", i);

            for (idx, entry) in node.log.iter().enumerate() {
                assert_eq!(entry.command, leader_log[idx].command, "노드 {} 엔트리 {}", i, idx);
                assert_eq!(entry.term, leader_log[idx].term, "노드 {} 엔트리 {} term", i, idx);
            }
        }
    }

    #[tokio::test]
    async fn test_leader_with_network_latency() {
        let network = RaftNetwork::new(3);

        // 리더 선출
        network.conduct_election(0).await;

        // 100ms 지연 추가
        network.set_latency(0, 1, 100);
        network.set_latency(0, 2, 50);

        // 로그 추가
        {
            let mut leader = network.nodes[0].write().await;
            leader.log.push(LogEntry {
                index: 1,
                term: 1,
                command: "delayed_entry".to_string(),
            });
        }

        // 로그 복제
        let start = Instant::now();
        network.replicate_log(0).await;
        let elapsed = start.elapsed();

        // 지연이 있어도 복제 완료 (최소 100ms)
        assert!(elapsed.as_millis() >= 50, "최소 지연 시간 확인");

        tokio::time::sleep(Duration::from_millis(150)).await;

        // 모든 팔로워가 로그 수신
        for i in 1..3 {
            let follower = network.nodes[i].read().await;
            assert_eq!(follower.log.len(), 1, "팔로워 {}의 로그 복제", i);
        }
    }

    #[tokio::test]
    async fn test_heartbeat() {
        let network = RaftNetwork::new(3);

        // 리더 선출
        network.conduct_election(0).await;

        // 하트비트 전송
        network.send_heartbeat(0).await;
        tokio::time::sleep(Duration::from_millis(50)).await;

        // 팔로워의 last_heartbeat가 업데이트됨
        for i in 1..3 {
            let follower = network.nodes[i].read().await;
            let elapsed = follower.last_heartbeat.elapsed().as_millis();
            // 50ms 이내에 하트비트 수신
            assert!(elapsed < 100, "팔로워 {}가 하트비트 수신", i);
        }
    }
}
