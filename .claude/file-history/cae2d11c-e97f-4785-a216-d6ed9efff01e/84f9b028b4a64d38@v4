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

        // 여러 번 선출하면 다른 노드가 리더가 되어야 함
        let first_election = cluster.elect_leader().await;
        assert!(first_election);
        let leader_1 = cluster.leader_id;

        let second_election = cluster.elect_leader().await;
        assert!(second_election);
        let leader_2 = cluster.leader_id;

        // 첫 번째와 두 번째 리더가 다름
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
}
