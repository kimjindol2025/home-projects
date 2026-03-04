# Challenge 13: Sovereign-Pay (독자적 결제 게이트웨이)
## Zero-Knowledge Proof 기반 오프라인 P2P 결제 프로토콜

**Challenge Status**: 🚀 STARTING
**Target Date**: 2026-03-25
**Scope**: FL-Signal 확장 + NFC/UWB P2P 결제 + ZKP 프라이버시
**Expected Deliverables**: 3,400-4,000 lines | 38+ tests | 8 unforgiving rules

---

## 📋 Challenge Objectives

### Primary Goals
1. **Zero-Knowledge Proof**: 잔액 공개 없이 지불 능력만 증명
2. **NFC/UWB 근거리 통신**: 오프라인 즉시 결제
3. **Double-Spend Prevention**: Offline에서 중복 지불 방지
4. **Privacy by Design**: 상점에 개인정보 노출 안함
5. **Distributed Settlement**: 나중의 배치 정산

### Innovation Points
- **Pedersen Commitment**: 잔액을 숨긴 상태로 증명
- **Range Proof**: "잔액 ≥ 지불액"만 증명 (actual amount 불공개)
- **Offline Ledger**: NFC 거래 로컬 저장 + 나중에 동기화
- **UWB Distance**: 물리적 거리 기반 보안 (MITM 방지)
- **Merkle Tree**: 거래 무결성 증명

### Success Criteria
- ✅ ZKP 증명 시간 <1 second
- ✅ NFC/UWB 거래 완료 <3 seconds
- ✅ Offline 거래 저장 성공률 ≥99.5%
- ✅ Settlement 충돌률 <0.1%
- ✅ Privacy: 상점이 수신인 신원 추론 불가능
- ✅ Range proof 정확도 ≥99%

---

## 🏗️ Architecture

### 5-Module Sovereign-Pay Stack

#### Module 1: ZKProofEngine (800 lines)
**Purpose**: Zero-Knowledge Proof로 잔액 증명

```
Commitment Scheme (Pedersen):
  Balance b ∈ [0, 2^32]
  Random r (blinding factor)
  C = g^b · h^r (commitment - amount hidden)
    ↓
  Prover: "C represents ≥ payment_amount" (without revealing b)
    ↓
  Range Proof:
    - Prove b ∈ [0, 2^32] (valid balance range)
    - Prove b ≥ payment_amount (sufficient funds)
    - Zero-knowledge: payment_amount NOT leaked
    ↓
  Verifier: Accept only if proof valid
    (learns: payment is feasible, not actual amount)
```

**Techniques**:
- **Commitment**: Pedersen commitment (secure, homomorphic)
- **Range proof**: Bulletproof-style (efficient, <1 second)
- **Hash-based**: SHA256 for challenges in Fiat-Shamir transform
- **Verification**: Batch verification for multiple proofs

**Unforgiving Rules (Module 1)**:
- **Rule 1**: Proof generation <1 second
- **Rule 2**: False positive rate <0.01% (incorrect balance accepted)

#### Module 2: TransactionProtocol (700 lines)
**Purpose**: Offline P2P transaction creation & validation

```
Transaction Structure:
  Payer: Device A identity (anonymized)
  Payee: Device B identity
  Amount: payment_amount
  Timestamp: creation time
  Nonce: prevent replay
  Signature: payer's digital signature
  ZKProof: range proof of balance
    ↓
  Double-Spend Prevention (Offline):
    1. Local check: transaction not in local ledger
    2. Timestamp ordering: prevent temporal conflicts
    3. Nonce uniqueness: each transaction unique
    ↓
  Validation:
    - Signature valid?
    - ZKProof valid?
    - No duplicate nonce?
    ↓
  Result: Valid/Invalid transaction
```

**Unforgiving Rules (Module 2)**:
- **Rule 3**: Signature verification ≥99.9% accuracy
- **Rule 4**: No offline double-spend recorded

#### Module 3: NFCUWBProtocol (600 lines)
**Purpose**: Secure NFC/UWB communication

```
NFC/UWB Transaction Flow:
  Device A (Payer) [NFC range: ~4cm]
    ↓ [NFC: establish secure channel]
  Device B (Payee)
    ├─ Handshake (ECDH key exchange)
    ├─ UWB distance check: <5m (prevents remote attacks)
    ├─ Mutual authentication
    └─ Payment confirmation
    ↓
  Transaction details exchange (encrypted)
    ↓
  Both devices sign & store locally
    ↓
  Completion <3 seconds
```

**Security**:
- **NFC encryption**: AES-128 payload encryption
- **UWB distance**: Physical proximity verification (MITM resistant)
- **Authentication**: ECDH + signature verification
- **Replay protection**: Timestamp + nonce

**Unforgiving Rules (Module 3)**:
- **Rule 5**: Transaction completion <3 seconds
- **Rule 6**: UWB MITM detection ≥99%

#### Module 4: DistributedLedger (700 lines)
**Purpose**: Offline transaction log with integrity

```
Local Ledger Structure:
  Transaction 1: [payer, payee, amount, proof, sig, timestamp, hash]
  Transaction 2: ...
  Transaction N: ...
    ↓
  Merkle Tree:
    H(T1) → ┐
    H(T2) → ├─ H(H(T1)+H(T2)) → ┐
    H(T3) → ┐                      ├─ Merkle Root
    H(T4) → ├─ H(H(T3)+H(T4)) → ┘
    ...
    ↓
  Root Hash: cryptographic commitment to all transactions
    ↓
  Storage:
    - Persistent: SD card / local storage
    - Synchronization: batch upload when online
    - Conflict resolution: first-write-wins + merkle proof
```

**Unforgiving Rules (Module 4)**:
- **Rule 7**: Ledger write success ≥99.5%

#### Module 5: SettlementEngine (500 lines)
**Purpose**: Batch settlement & double-spend detection

```
Settlement Pipeline (when online):
  1. Upload local ledger to server
  2. Server aggregates all devices' transactions
  3. Detect conflicts:
     - Same nonce from different device? → Fraud
     - Same payer, >1 conflicting transaction? → Double-spend
  4. Resolution:
     - Accept: first-seen-wins (timestamp)
     - Reject: second transaction marked invalid
  5. Account reconciliation:
     - Verify payer balance = initial - (all valid outflows) + (all inflows)
  6. Settlement: update canonical ledger
    ↓
  Merkle proof of settlement sent to all participants
```

**Unforgiving Rules (Module 5)**:
- **Rule 8**: Settlement conflict detection ≥99.5%

---

## 🧪 Test Plan (38+ tests)

### Group A: ZKProofEngine (10 tests)
```
✓ test_pedersen_commitment
✓ test_range_proof_generation
✓ test_proof_verification
✓ test_proof_generation_time (<1s)
✓ test_false_positive_rate (<0.01%)
✓ test_balance_hiding
✓ test_proof_completeness
✓ test_batch_verification
✓ test_invalid_range_rejection
✓ test_commitment_homomorphism
```

### Group B: TransactionProtocol (8 tests)
```
✓ test_transaction_creation
✓ test_signature_verification (≥99.9%)
✓ test_nonce_uniqueness
✓ test_double_spend_detection
✓ test_transaction_validation
✓ test_offline_double_spend_prevention
✓ test_replay_attack_prevention
✓ test_transaction_integrity
```

### Group C: NFCUWBProtocol (6 tests)
```
✓ test_nfc_handshake
✓ test_uwb_distance_verification
✓ test_encryption
✓ test_transaction_completion_time (<3s)
✓ test_mitm_detection (≥99%)
✓ test_concurrent_payments
```

### Group D: DistributedLedger (8 tests)
```
✓ test_ledger_creation
✓ test_transaction_append
✓ test_merkle_tree_construction
✓ test_merkle_root_verification
✓ test_ledger_persistence
✓ test_write_success_rate (≥99.5%)
✓ test_conflict_resolution
✓ test_ledger_synchronization
```

### Group E: SettlementEngine (6+ tests)
```
✓ test_settlement_creation
✓ test_double_spend_detection (≥99.5%)
✓ test_batch_settlement
✓ test_balance_reconciliation
✓ test_merkle_proof_verification
✓ test_settlement_finality
```

---

## 📊 Unforgiving Rules (8 total)

| Rule | Target | Implementation |
|------|--------|-----------------|
| **R1** | ZKP generation <1s | ZKProofEngine::generate_proof() |
| **R2** | False positive <0.01% | ZKProofEngine::verify_proof() |
| **R3** | Signature accuracy ≥99.9% | TransactionProtocol::verify_signature() |
| **R4** | No offline double-spend | TransactionProtocol::check_double_spend() |
| **R5** | NFC completion <3s | NFCUWBProtocol::execute_transaction() |
| **R6** | MITM detection ≥99% | NFCUWBProtocol::verify_distance() |
| **R7** | Ledger write ≥99.5% | DistributedLedger::append_transaction() |
| **R8** | Settlement conflict <0.1% | SettlementEngine::detect_conflicts() |

---

## 📁 File Structure

```
src/
├── zk_proof_engine.rs          (800 lines)
│   ├── PedersenCommitment
│   ├── RangeProof
│   ├── ZKProofEngine
│   └── [10 test functions]
│
├── transaction_protocol.rs      (700 lines)
│   ├── Transaction
│   ├── Signature
│   ├── TransactionProtocol
│   └── [8 test functions]
│
├── nfc_uwb_protocol.rs          (600 lines)
│   ├── NFCChannel
│   ├── UWBDistance
│   ├── NFCUWBProtocol
│   └── [6 test functions]
│
├── distributed_ledger.rs        (700 lines)
│   ├── LedgerEntry
│   ├── MerkleTree
│   ├── DistributedLedger
│   └── [8 test functions]
│
├── settlement_engine.rs         (500 lines)
│   ├── SettlementBatch
│   ├── ConflictResolution
│   ├── SettlementEngine
│   └── [6+ test functions]
│
└── lib.rs                       (updated)
```

---

## 🎯 Implementation Strategy

### Phase A: ZKP Foundation (Days 1-2)
1. Pedersen commitment implementation
2. Range proof generation & verification
3. Test accuracy & performance
4. Batch verification optimization

### Phase B: Transaction Protocol (Days 3-4)
1. Transaction structure & serialization
2. Digital signature implementation
3. Offline double-spend prevention
4. Validation logic

### Phase C: NFC/UWB Communication (Day 5)
1. NFC handshake & encryption
2. UWB distance verification (MITM prevention)
3. Concurrent payment handling
4. Performance <3 seconds

### Phase D: Ledger & Settlement (Day 6)
1. Merkle tree for transaction integrity
2. Offline ledger persistence
3. Batch settlement & conflict detection
4. Account reconciliation

---

## 📈 Privacy Guarantees

**Zero-Knowledge Properties**:
- ✅ Merchant learns: Transaction valid, amount ≤ balance
- ✅ Merchant does NOT learn: Actual balance, transaction history
- ✅ Server learns: Only settlement-time amounts (aggregated)
- ✅ No deanonymization: Payer identity cryptographically hidden

**Example**:
```
Scenario: Alice pays Bob ₩100,000

Alice's device shows Bob's device:
  - Commitment C = g^b · h^r (b = balance, r = random)
  - RangeProof: π_range (proves b ≥ ₩100,000)

Bob learns:
  ✓ Alice can pay ₩100,000
  ✗ Alice's actual balance (could be ₩100k or ₩1M)
  ✗ Alice's identity (anonymous device key only)

Even server (on settlement) learns minimal:
  ✓ Settlement timestamp & amounts
  ✗ Individual transaction history (merged in batch)
```

---

**Next Step**: Implement ZKProofEngine → TransactionProtocol → NFCUWBProtocol → DistributedLedger → SettlementEngine

**Status**: Design approved, ready for implementation 🔧
