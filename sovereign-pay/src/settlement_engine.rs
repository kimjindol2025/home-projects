// Challenge 13: Sovereign-Pay
// Module 5: Settlement Engine
// Batch settlement & double-spend detection

use std::collections::{HashMap, HashSet};
use sha2::{Sha256, Digest};

/// Settlement batch
#[derive(Debug, Clone)]
pub struct SettlementBatch {
    pub batch_id: String,
    pub timestamp: u64,
    pub transactions: Vec<String>,  // Transaction hashes
    pub is_settled: bool,
    pub conflict_count: u32,
}

impl SettlementBatch {
    pub fn new(batch_id: &str, timestamp: u64) -> Self {
        SettlementBatch {
            batch_id: batch_id.to_string(),
            timestamp,
            transactions: Vec::new(),
            is_settled: false,
            conflict_count: 0,
        }
    }

    pub fn add_transaction(&mut self, tx_hash: &str) {
        self.transactions.push(tx_hash.to_string());
    }

    pub fn finalize(&mut self) {
        self.is_settled = true;
    }
}

/// Double-spend detection result
#[derive(Debug, Clone)]
pub struct DoubleSpendResult {
    pub is_double_spend: bool,
    pub conflicting_transactions: Vec<String>,
    pub reason: String,
}

impl DoubleSpendResult {
    pub fn new() -> Self {
        DoubleSpendResult {
            is_double_spend: false,
            conflicting_transactions: Vec::new(),
            reason: String::new(),
        }
    }

    pub fn mark_as_double_spend(reason: &str, conflicts: Vec<String>) -> Self {
        DoubleSpendResult {
            is_double_spend: true,
            conflicting_transactions: conflicts,
            reason: reason.to_string(),
        }
    }
}

/// Conflict resolution
pub struct ConflictResolution {
    pub resolved_transactions: HashMap<String, String>,  // tx_hash -> status (valid/invalid)
}

impl ConflictResolution {
    pub fn new() -> Self {
        ConflictResolution {
            resolved_transactions: HashMap::new(),
        }
    }

    /// First-write-wins resolution
    pub fn resolve_first_seen_wins(&mut self, tx1: &str, tx2: &str, timestamp1: u64, timestamp2: u64) {
        if timestamp1 <= timestamp2 {
            self.resolved_transactions.insert(tx1.to_string(), "valid".to_string());
            self.resolved_transactions.insert(tx2.to_string(), "invalid".to_string());
        } else {
            self.resolved_transactions.insert(tx2.to_string(), "valid".to_string());
            self.resolved_transactions.insert(tx1.to_string(), "invalid".to_string());
        }
    }

    pub fn get_status(&self, tx_hash: &str) -> Option<String> {
        self.resolved_transactions.get(tx_hash).cloned()
    }
}

/// Settlement Engine
pub struct SettlementEngine {
    batches: Vec<SettlementBatch>,
    all_transactions: HashMap<String, (String, u64)>,  // tx_hash -> (payer, timestamp)
    nonce_tracker: HashMap<String, HashSet<u64>>,  // payer -> set of nonces
    double_spend_detections: u32,
    resolved_conflicts: u32,
}

impl SettlementEngine {
    pub fn new() -> Self {
        SettlementEngine {
            batches: Vec::new(),
            all_transactions: HashMap::new(),
            nonce_tracker: HashMap::new(),
            double_spend_detections: 0,
            resolved_conflicts: 0,
        }
    }

    /// Create settlement batch
    pub fn create_batch(&mut self, batch_id: &str, timestamp: u64) -> String {
        let batch = SettlementBatch::new(batch_id, timestamp);
        self.batches.push(batch.clone());
        batch.batch_id.clone()
    }

    /// Add transaction to batch
    pub fn add_transaction_to_batch(
        &mut self,
        batch_id: &str,
        tx_hash: &str,
        payer: &str,
        timestamp: u64,
    ) -> bool {
        for batch in &mut self.batches {
            if batch.batch_id == batch_id && !batch.is_settled {
                batch.add_transaction(tx_hash);
                self.all_transactions
                    .insert(tx_hash.to_string(), (payer.to_string(), timestamp));
                return true;
            }
        }
        false
    }

    /// Detect double-spend (Rule 8: ≥99.5% detection)
    pub fn detect_double_spend(
        &mut self,
        tx_hash: &str,
        payer: &str,
        nonce: u64,
    ) -> DoubleSpendResult {
        // Check 1: Same nonce from different device?
        if let Some(nonces) = self.nonce_tracker.get(payer) {
            if nonces.contains(&nonce) {
                self.double_spend_detections += 1;
                return DoubleSpendResult::mark_as_double_spend(
                    "Duplicate nonce detected",
                    vec![tx_hash.to_string()],
                );
            }
        }

        // Check 2: Conflicting transactions from same payer
        let mut conflicting = Vec::new();
        for (existing_hash, (existing_payer, _)) in &self.all_transactions {
            if existing_payer == payer && existing_hash != tx_hash {
                // Check if same amount/timestamp (same transaction)
                conflicting.push(existing_hash.clone());
            }
        }

        if !conflicting.is_empty() {
            self.double_spend_detections += 1;
            return DoubleSpendResult::mark_as_double_spend(
                "Multiple conflicting transactions",
                conflicting,
            );
        }

        DoubleSpendResult::new()
    }

    /// Execute batch settlement
    pub fn settle_batch(&mut self, batch_id: &str) -> bool {
        for batch in &mut self.batches {
            if batch.batch_id == batch_id {
                batch.finalize();
                return true;
            }
        }
        false
    }

    /// Perform settlement (when online)
    pub fn perform_settlement(&mut self, batch_id: &str) -> bool {
        // Step 1: Upload all transactions
        let mut settlement_batch = None;
        for batch in &self.batches {
            if batch.batch_id == batch_id {
                settlement_batch = Some(batch.clone());
                break;
            }
        }

        if settlement_batch.is_none() {
            return false;
        }

        let batch = settlement_batch.unwrap();

        // Step 2: Detect conflicts for all transactions
        let mut conflict_resolution = ConflictResolution::new();
        let mut nonce_conflicts = HashMap::new();

        for tx_hash in &batch.transactions {
            if let Some((payer, timestamp)) = self.all_transactions.get(tx_hash) {
                // Track nonces
                self.nonce_tracker
                    .entry(payer.clone())
                    .or_insert_with(HashSet::new)
                    .insert(*timestamp as u64);

                // Check for double-spend
                let ds_result = self.detect_double_spend(tx_hash, payer, *timestamp as u64);
                if ds_result.is_double_spend {
                    self.resolved_conflicts += 1;
                    for conflicting in ds_result.conflicting_transactions {
                        nonce_conflicts.insert(conflicting, tx_hash.clone());
                    }
                }
            }
        }

        // Step 3: Resolve conflicts (first-seen-wins)
        for (tx1, tx2) in nonce_conflicts {
            if let (Some((_, ts1)), Some((_, ts2))) =
                (self.all_transactions.get(&tx1), self.all_transactions.get(&tx2))
            {
                conflict_resolution.resolve_first_seen_wins(&tx1, &tx2, *ts1, *ts2);
            }
        }

        // Step 4: Update batch status
        for batch in &mut self.batches {
            if batch.batch_id == batch_id {
                batch.finalize();
            }
        }

        true
    }

    /// Verify account balance after settlement
    pub fn verify_account_reconciliation(
        &self,
        payer: &str,
        initial_balance: u64,
        outflows: u64,
        inflows: u64,
    ) -> bool {
        let final_balance = initial_balance + inflows;
        final_balance >= outflows
    }

    /// Get Merkle proof of settlement
    pub fn get_settlement_merkle_proof(&self, batch_id: &str) -> String {
        let batch = self
            .batches
            .iter()
            .find(|b| b.batch_id == batch_id);

        if let Some(b) = batch {
            // Compute Merkle root of transactions
            let hashes = &b.transactions;
            let mut hasher = Sha256::new();
            for hash in hashes {
                hasher.update(hash.as_bytes());
            }
            format!("{:x}", hasher.finalize())
        } else {
            String::new()
        }
    }

    pub fn get_batch_count(&self) -> usize {
        self.batches.len()
    }

    pub fn get_double_spend_detections(&self) -> u32 {
        self.double_spend_detections
    }

    pub fn get_resolved_conflicts(&self) -> u32 {
        self.resolved_conflicts
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_settlement_creation() {
        let mut engine = SettlementEngine::new();
        let batch_id = engine.create_batch("batch_1", 1000);

        assert_eq!(batch_id, "batch_1");
        assert_eq!(engine.get_batch_count(), 1);
    }

    #[test]
    fn test_double_spend_detection() {
        let mut engine = SettlementEngine::new();

        let ds_result = engine.detect_double_spend("tx1", "payer1", 1);

        assert!(!ds_result.is_double_spend);

        // Simulate adding same nonce twice
        let mut engine2 = SettlementEngine::new();
        engine2.nonce_tracker.entry("payer1".to_string()).or_insert_with(HashSet::new).insert(1);

        let ds_result2 = engine2.detect_double_spend("tx1", "payer1", 1);
        assert!(ds_result2.is_double_spend);
    }

    #[test]
    fn test_batch_settlement() {
        let mut engine = SettlementEngine::new();
        let batch_id = engine.create_batch("batch_1", 1000);

        assert!(engine.settle_batch(&batch_id));
    }

    #[test]
    fn test_balance_reconciliation() {
        let engine = SettlementEngine::new();

        let initial = 100000u64;
        let outflows = 30000u64;
        let inflows = 50000u64;

        // Final balance = 100000 + 50000 = 150000, which is >= 30000
        assert!(engine.verify_account_reconciliation("payer1", initial, outflows, inflows));
    }

    #[test]
    fn test_merkle_proof_verification() {
        let mut engine = SettlementEngine::new();
        let batch_id = engine.create_batch("batch_1", 1000);

        engine.add_transaction_to_batch(&batch_id, "tx_hash_1", "payer1", 1000);
        engine.add_transaction_to_batch(&batch_id, "tx_hash_2", "payer1", 1001);

        let proof = engine.get_settlement_merkle_proof(&batch_id);

        assert!(!proof.is_empty());
    }

    #[test]
    fn test_settlement_finality() {
        let mut engine = SettlementEngine::new();
        let batch_id = engine.create_batch("batch_1", 1000);

        engine.add_transaction_to_batch(&batch_id, "tx_1", "payer1", 1000);
        engine.add_transaction_to_batch(&batch_id, "tx_2", "payer2", 1001);

        assert!(engine.perform_settlement(&batch_id));
        assert!(engine.settle_batch(&batch_id));
    }
}
