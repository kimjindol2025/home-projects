// Challenge 13: Sovereign-Pay
// Module 4: Distributed Ledger
// Offline transaction log with Merkle tree integrity

use sha2::{Sha256, Digest};
use std::collections::HashMap;

/// Ledger entry
#[derive(Debug, Clone)]
pub struct LedgerEntry {
    pub payer: String,
    pub payee: String,
    pub amount: u64,
    pub proof_hash: String,
    pub signature_hash: String,
    pub timestamp: u64,
    pub entry_hash: String,
}

impl LedgerEntry {
    pub fn new(
        payer: &str,
        payee: &str,
        amount: u64,
        proof_hash: &str,
        signature_hash: &str,
        timestamp: u64,
    ) -> Self {
        // Compute entry hash
        let mut hasher = Sha256::new();
        hasher.update(payer.as_bytes());
        hasher.update(payee.as_bytes());
        hasher.update(amount.to_le_bytes());
        hasher.update(proof_hash.as_bytes());
        hasher.update(signature_hash.as_bytes());
        hasher.update(timestamp.to_le_bytes());
        let entry_hash = format!("{:x}", hasher.finalize());

        LedgerEntry {
            payer: payer.to_string(),
            payee: payee.to_string(),
            amount,
            proof_hash: proof_hash.to_string(),
            signature_hash: signature_hash.to_string(),
            timestamp,
            entry_hash,
        }
    }
}

/// Merkle Tree node
#[derive(Debug, Clone)]
pub struct MerkleNode {
    pub hash: String,
    pub left: Option<Box<MerkleNode>>,
    pub right: Option<Box<MerkleNode>>,
}

impl MerkleNode {
    pub fn new(hash: &str) -> Self {
        MerkleNode {
            hash: hash.to_string(),
            left: None,
            right: None,
        }
    }

    pub fn with_children(hash: &str, left: MerkleNode, right: MerkleNode) -> Self {
        MerkleNode {
            hash: hash.to_string(),
            left: Some(Box::new(left)),
            right: Some(Box::new(right)),
        }
    }
}

/// Merkle Tree for transaction integrity
pub struct MerkleTree {
    root: Option<MerkleNode>,
    leaf_count: u32,
}

impl MerkleTree {
    pub fn new() -> Self {
        MerkleTree {
            root: None,
            leaf_count: 0,
        }
    }

    /// Build Merkle tree from transaction hashes
    pub fn build_from_hashes(&mut self, entry_hashes: &[String]) -> String {
        if entry_hashes.is_empty() {
            return String::from("empty_tree");
        }

        self.leaf_count = entry_hashes.len() as u32;

        // Create leaf nodes
        let mut nodes: Vec<MerkleNode> = entry_hashes
            .iter()
            .map(|hash| MerkleNode::new(hash))
            .collect();

        // Build tree bottom-up
        while nodes.len() > 1 {
            let mut new_level = Vec::new();

            // Process pairs
            for i in (0..nodes.len()).step_by(2) {
                let left = nodes[i].clone();
                let right = if i + 1 < nodes.len() {
                    nodes[i + 1].clone()
                } else {
                    left.clone()  // Duplicate if odd number
                };

                // Compute parent hash
                let mut hasher = Sha256::new();
                hasher.update(&left.hash);
                hasher.update(&right.hash);
                let parent_hash = format!("{:x}", hasher.finalize());

                new_level.push(MerkleNode::with_children(&parent_hash, left, right));
            }

            nodes = new_level;
        }

        if let Some(root) = nodes.into_iter().next() {
            let root_hash = root.hash.clone();
            self.root = Some(root);
            root_hash
        } else {
            String::from("error")
        }
    }

    /// Get Merkle root hash
    pub fn get_root_hash(&self) -> String {
        if let Some(root) = &self.root {
            root.hash.clone()
        } else {
            String::new()
        }
    }

    /// Verify inclusion of entry in tree
    pub fn verify_inclusion(&self, entry_hash: &str) -> bool {
        if let Some(root) = &self.root {
            self.verify_in_subtree(root, entry_hash)
        } else {
            false
        }
    }

    fn verify_in_subtree(&self, node: &MerkleNode, target: &str) -> bool {
        if node.hash == target {
            return true;
        }

        if let Some(left) = &node.left {
            if self.verify_in_subtree(left, target) {
                return true;
            }
        }

        if let Some(right) = &node.right {
            if self.verify_in_subtree(right, target) {
                return true;
            }
        }

        false
    }

    pub fn get_leaf_count(&self) -> u32 {
        self.leaf_count
    }
}

/// Distributed Ledger
pub struct DistributedLedger {
    entries: Vec<LedgerEntry>,
    merkle_tree: MerkleTree,
    root_hash: String,
    write_success_count: u32,
    write_failure_count: u32,
    conflict_count: u32,
}

impl DistributedLedger {
    pub fn new() -> Self {
        DistributedLedger {
            entries: Vec::new(),
            merkle_tree: MerkleTree::new(),
            root_hash: String::new(),
            write_success_count: 0,
            write_failure_count: 0,
            conflict_count: 0,
        }
    }

    /// Append transaction to ledger
    pub fn append_transaction(&mut self, entry: LedgerEntry) -> bool {
        // Simulate persistent write (Rule 7: ≥99.5% success)
        use std::time::{SystemTime, UNIX_EPOCH};
        let nanos = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .subsec_nanos();

        let failure_chance = (nanos % 2000) as f32 / 2000.0;  // ~0.05% failure rate
        if failure_chance > 0.995 {
            self.write_failure_count += 1;
            return false;
        }

        self.entries.push(entry);
        self.write_success_count += 1;

        // Rebuild Merkle tree
        self.rebuild_merkle_tree();
        true
    }

    /// Rebuild Merkle tree with current entries
    fn rebuild_merkle_tree(&mut self) {
        let hashes: Vec<String> = self.entries.iter().map(|e| e.entry_hash.clone()).collect();
        self.root_hash = self.merkle_tree.build_from_hashes(&hashes);
    }

    /// Get Merkle root for verification
    pub fn get_merkle_root(&self) -> String {
        self.root_hash.clone()
    }

    /// Verify transaction is in ledger
    pub fn verify_transaction_in_ledger(&self, entry_hash: &str) -> bool {
        self.merkle_tree.verify_inclusion(entry_hash)
    }

    /// Persist ledger to disk (simplified)
    pub fn persist_to_disk(&self) -> bool {
        // Simulate persistence
        !self.entries.is_empty()
    }

    /// Synchronize with server (batch upload)
    pub fn synchronize_with_server(&mut self) -> bool {
        if self.entries.is_empty() {
            return true;
        }

        // Simulate conflict resolution
        for entry in &self.entries {
            // Check for conflicts based on payer + timestamp
            let conflicts = self
                .entries
                .iter()
                .filter(|e| e.payer == entry.payer && e.timestamp == entry.timestamp)
                .count();

            if conflicts > 1 {
                self.conflict_count += 1;
            }
        }

        true
    }

    /// Resolve conflicts (first-write-wins)
    pub fn resolve_conflict(&mut self, entry1: &LedgerEntry, entry2: &LedgerEntry) -> Option<LedgerEntry> {
        if entry1.timestamp <= entry2.timestamp {
            Some(entry1.clone())
        } else {
            Some(entry2.clone())
        }
    }

    pub fn get_ledger_size(&self) -> usize {
        self.entries.len()
    }

    pub fn get_write_success_rate(&self) -> f32 {
        let total = self.write_success_count + self.write_failure_count;
        if total == 0 {
            100.0
        } else {
            (self.write_success_count as f32 / total as f32) * 100.0
        }
    }

    pub fn get_conflict_count(&self) -> u32 {
        self.conflict_count
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ledger_creation() {
        let ledger = DistributedLedger::new();
        assert_eq!(ledger.get_ledger_size(), 0);
    }

    #[test]
    fn test_transaction_append() {
        let mut ledger = DistributedLedger::new();
        let entry = LedgerEntry::new(
            "payer1",
            "payee1",
            50000,
            "proof_hash_1",
            "sig_hash_1",
            1000,
        );

        assert!(ledger.append_transaction(entry));
        assert_eq!(ledger.get_ledger_size(), 1);
    }

    #[test]
    fn test_merkle_tree_construction() {
        let mut ledger = DistributedLedger::new();

        let entries = vec![
            LedgerEntry::new("payer1", "payee1", 50000, "proof1", "sig1", 1000),
            LedgerEntry::new("payer2", "payee2", 30000, "proof2", "sig2", 1001),
            LedgerEntry::new("payer3", "payee3", 20000, "proof3", "sig3", 1002),
        ];

        for entry in entries {
            ledger.append_transaction(entry);
        }

        assert!(!ledger.get_merkle_root().is_empty());
    }

    #[test]
    fn test_merkle_root_verification() {
        let mut merkle_tree = MerkleTree::new();
        let hashes = vec![
            "hash1".to_string(),
            "hash2".to_string(),
            "hash3".to_string(),
        ];

        let _root = merkle_tree.build_from_hashes(&hashes);

        assert!(!merkle_tree.get_root_hash().is_empty());
        assert_eq!(merkle_tree.get_leaf_count(), 3);
    }

    #[test]
    fn test_ledger_persistence() {
        let mut ledger = DistributedLedger::new();
        let entry = LedgerEntry::new(
            "payer1",
            "payee1",
            50000,
            "proof_hash",
            "sig_hash",
            1000,
        );

        ledger.append_transaction(entry);

        assert!(ledger.persist_to_disk());
    }

    #[test]
    fn test_write_success_rate() {
        let mut ledger = DistributedLedger::new();

        for i in 0..20 {
            let entry = LedgerEntry::new(
                "payer",
                "payee",
                50000,
                &format!("proof_{}", i),
                &format!("sig_{}", i),
                1000 + i as u64,
            );
            ledger.append_transaction(entry);
        }

        let success_rate = ledger.get_write_success_rate();
        // Should be ≥99.5% (Rule 7)
        assert!(success_rate >= 99.5);
    }

    #[test]
    fn test_conflict_resolution() {
        let mut ledger = DistributedLedger::new();

        let entry1 = LedgerEntry::new("payer1", "payee1", 50000, "proof1", "sig1", 1000);
        let entry2 = LedgerEntry::new("payer1", "payee2", 30000, "proof2", "sig2", 1001);

        let resolved = ledger.resolve_conflict(&entry1, &entry2);

        assert!(resolved.is_some());
        assert_eq!(resolved.unwrap().payer, "payer1");
    }

    #[test]
    fn test_ledger_synchronization() {
        let mut ledger = DistributedLedger::new();

        for i in 0..5 {
            let entry = LedgerEntry::new(
                &format!("payer_{}", i % 2),
                &format!("payee_{}", i),
                50000,
                &format!("proof_{}", i),
                &format!("sig_{}", i),
                1000 + i as u64,
            );
            ledger.append_transaction(entry);
        }

        assert!(ledger.synchronize_with_server());
    }
}
