// Challenge 13: Sovereign-Pay
// Zero-Knowledge Proof based Offline P2P Payment Protocol

pub mod zk_proof_engine;
pub mod transaction_protocol;
pub mod nfc_uwb_protocol;
pub mod distributed_ledger;
pub mod settlement_engine;

pub use zk_proof_engine::{
    PedersenParams, PedersenCommitment, RangeProof, ZKProofEngine,
};

pub use transaction_protocol::{
    Signature, Transaction, TransactionProtocol,
};

pub use nfc_uwb_protocol::{
    NFCChannel, UWBDistance, NFCUWBProtocol,
};

pub use distributed_ledger::{
    LedgerEntry, MerkleNode, MerkleTree, DistributedLedger,
};

pub use settlement_engine::{
    SettlementBatch, DoubleSpendResult, ConflictResolution, SettlementEngine,
};
