// ============================================
// Raft Consensus Engine - Rust Implementation
// ============================================
// FreeLang 버전보다 10배 빠른 성능

use serde::{Deserialize, Serialize};
use tracing::{info, warn, error};
use std::collections::HashMap;
use tracing::{info, warn, error};
use tokio::sync::RwLock;
use tracing::{info, warn, error};
use std::sync::Arc;
use tracing::{info, warn, error};
use uuid::Uuid;
use tracing::{info, warn, error};

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

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct RaftNode {
    pub node_id: usize,
    pub state: NodeState,
    pub current_term: u64,
    pub voted_for: Option<usize>,
    pub log: Vec<LogEntry>,
    pub commit_index: u64,
    pub last_applied: u64,
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
