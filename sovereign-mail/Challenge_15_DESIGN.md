# Challenge 15: Sovereign-Naming (분산 DNS)

**Status**: 🚀 DESIGN PHASE
**Target Date**: 2026-03-11 (1주)
**Scope**: 4개 모듈, 2,400줄 코드, 30개 테스트, 6개 무관용 규칙

---

## 📋 Challenge Objectives

### Primary Goals

1. **@sovereign 주소 체계**
   - 중앙화 제거 (ICANN 독립)
   - 누구도 차단 불가능
   - 영구 불멸 (서버 없이 작동)

2. **DHT 기반 분산 저장**
   - Kademlia 프로토콜 (IPFS 호환)
   - 피어-투-피어 네트워크
   - 자동 복제 (최소 3개 노드)

3. **공개키 주소화**
   - kim@sovereign → kim의 공개키
   - 직접 암호화 통신
   - 중간자 공격 방지

4. **네트워크 일관성**
   - CRDT 기반 동기화
   - Gossip 프로토콜
   - 최종 일관성 보장

### Success Criteria
- ✅ DNS 조회 지연 <100ms (Kademlia DHT)
- ✅ 노드 장애 복구 <500ms (quorum-based)
- ✅ 위조 응답 탐지율 99.9% (zero-knowledge proof)
- ✅ 주소 일관성 100% (gossip protocol)
- ✅ 개인키 노출 0건 (encrypted-at-rest)
- ✅ 대역폭 <10KB per lookup (optimized messages)

---

## 🏗️ Architecture

### 4-Module Sovereign-Naming Stack

**Module 1: DHT-Node (900 lines)**
- Kademlia 노드 구현
- 피어 발견/관리
- 분산 저장소

**Module 2: Address-Resolver (600 lines)**
- @sovereign 주소 해석
- 공개키 조회
- DNS-like 쿼리

**Module 3: Reputation-Manager (500 lines)**
- 발신자 신용도 추적
- 블랙리스트/화이트리스트
- 평판 점수 계산

**Module 4: Consensus (400 lines)**
- CRDT (Conflict-free Replicated Data Type)
- 네트워크 동기화
- 일관성 검증

---

## 📊 Unforgiving Rules (6 total)

| Rule | Target | Implementation |
|------|--------|-----------------|
| **R1** | DNS lookup <100ms | AddressResolver::resolve() |
| **R2** | Failover <500ms | DHTNode::find_alternative_peer() |
| **R3** | Forgery detection 99.9% | ZKProof::verify() |
| **R4** | Consistency 100% | Consensus::gossip_sync() |
| **R5** | Key exposure 0 events | encrypted_at_rest() |
| **R6** | Bandwidth <10KB/lookup | optimize_dht_message() |

---

## 🧪 Test Plan (30+ tests)

### Group A: DHT-Node (10 tests)
- ✓ Kademlia distance calculation
- ✓ k-bucket management
- ✓ Peer discovery
- ✓ Iterative lookup
- ✓ Data storage/retrieval
- ✓ Automatic replication (3x)
- ✓ Node join
- ✓ Node leave
- ✓ Network partition recovery
- ✓ Large network simulation (1000+ nodes)

### Group B: Address-Resolver (10 tests)
- ✓ @sovereign address parsing
- ✓ Address resolution
- ✓ Public key lookup
- ✓ Zero-knowledge proof verification
- ✓ Resolution caching (TTL)
- ✓ Fallback resolution
- ✓ Malformed address handling
- ✓ Unicode address support
- ✓ Resolution latency (<100ms)
- ✓ Batch resolution

### Group C: Reputation-Manager (5 tests)
- ✓ Reputation scoring
- ✓ Blacklist management
- ✓ Feedback loop
- ✓ Reputation decay
- ✓ Distributed reputation (DHT sync)

### Group D: Consensus (5 tests)
- ✓ CRDT merge
- ✓ Gossip protocol
- ✓ Vector clock
- ✓ Eventual consistency
- ✓ Network partition recovery

---

## 📈 Performance Goals

- Address lookup: <100ms
- Failover recovery: <500ms
- Bandwidth: <10KB per lookup
- Consistency: 100% (eventual)
- Forgery detection: 99.9%

---

## 🎯 Implementation Strategy

### Phase A: DHT Foundation (Days 1-2)
- Kademlia distance metric (XOR)
- k-bucket management
- Iterative lookup

### Phase B: Address Resolution (Day 3)
- @sovereign 주소 파싱
- Zero-knowledge proof
- Public key lookup

### Phase C: Reputation System (Day 4)
- Reputation scoring
- Blacklist/Whitelist
- Feedback loop

### Phase D: Consensus & Sync (Day 5)
- CRDT implementation
- Gossip protocol
- Vector clocks

### Phase E: Integration (Day 6)
- Module integration
- End-to-end tests
- Performance optimization

---

**Philosophy**: "누구도 당신의 주소를 차단할 수 없다"
(No one can block your address)

**Next Step**: Challenge 16 (L0NN-Mail-Sentry - AI 방어)

**Status**: Design approved, ready for implementation 🔧
