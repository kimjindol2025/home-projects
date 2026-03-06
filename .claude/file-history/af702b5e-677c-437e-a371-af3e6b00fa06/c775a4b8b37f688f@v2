# Challenge 15: Sovereign-Naming (분산 DNS 주소 체계)
## Decentralized Email Addressing without ICANN Control

**Status**: 🚀 DESIGN PHASE
**Target Date**: 2026-03-18 (1주)
**Scope**: 4개 모듈, 2,400줄 코드, 30개 테스트, 6개 무관용 규칙

---

## 📋 Challenge Objectives

### Primary Goals
1. **중앙화된 DNS 제거**
   - kim@gmail.com → kim@sovereign
   - 구글/ICANN이 도메인을 차단할 수 없음
   - 누구도 당신의 주소를 빼앗을 수 없음

2. **분산 해시 테이블(DHT) 기반**
   - Kademlia 알고리즘 (IPFS 호환)
   - 피어-투-피어 네트워크 (중앙 서버 없음)
   - 자동 복제 및 동기화

3. **공개키 중심 주소**
   - 이메일 주소 = 공개키 지문(fingerprint)
   - kim@sovereign → Kim's public key (RSA-4096)
   - 서버 부재 시에도 주소는 영구 유효

4. **평판 시스템**
   - 발신자의 신용도를 DHT에 저장
   - 스팸/악의 행위자 차단
   - 오픈 소스 신뢰 메커니즘

### Innovation Points
- **DNS 탈중앙화**: ICANN 의존도 0%
- **불변 주소**: 서버 사라져도 공개키는 영원히 유효
- **자기치유 네트워크**: 노드 손실 시 자동 복구
- **평판 검증**: 비잔틴 내성 (Byzantine-resistant) DHT

### Success Criteria
- ✅ DNS 조회 <100ms (DHT lookup)
- ✅ 노드 장애 시 복구 <500ms (quorum-based)
- ✅ 위조 응답 탐지 99.9% (signature verification)
- ✅ 주소 일관성 100% (gossip protocol)
- ✅ 개인키 노출 0건 (secure storage)
- ✅ 대역폭 <10KB per lookup (compression)

---

## 🏗️ Architecture

### 4-Module Sovereign-Naming Stack

#### Module 1: DHT-Node (900 lines)
**Purpose**: Kademlia DHT 노드 구현 (분산 저장소)

```
분산 주소 저장소
  ↓ [DHT Node]
  ├─ DHT 참여
  │   ├─ 부트스트랩 노드 연결 (initial peer discovery)
  │   ├─ Kademlia 라우팅 테이블 유지
  │   └─ k-buckets (20 bucket, 각 20개 피어)
  ├─ 데이터 저장
  │   ├─ put(key, value) - 주소 등록
  │   ├─ get(key) - 주소 조회
  │   └─ replicate() - 자동 복제 (3개 이상 노드)
  ├─ RPC 핸들러
  │   ├─ PING (노드 살아있나)
  │   ├─ FIND_NODE (가까운 노드 찾기)
  │   ├─ FIND_VALUE (데이터 찾기)
  │   └─ STORE (데이터 저장)
  └─ 네트워크 복구
      ├─ 죽은 피어 감지 (timeout)
      ├─ 새 피어 탐색 (refresh)
      └─ 데이터 재배치 (rebalance)
```

**Components**:
- `KademliaNode`: Kademlia DHT 구현
- `RoutingTable`: k-bucket 관리
- `RPCHandler`: PING/FIND_NODE/FIND_VALUE/STORE
- `NetworkResilience`: 장애 복구

**Techniques**:
- Kademlia distance metric (XOR distance)
- k-bucket refresh (replicate)
- Data expiration (time-to-live)
- Parallel lookups (α=3 concurrent)

**Unforgiving Rules (Module 1)**:
- **Rule 1**: DNS 조회 <100ms
- **Rule 2**: 장애 복구 <500ms

---

#### Module 2: Address-Resolver (600 lines)
**Purpose**: kim@sovereign 주소 해석 (공개키 조회)

```
주소 해석
  ↓ [Address Resolver]
  ├─ 주소 파싱
  │   ├─ kim@sovereign 분해
  │   ├─ 로컬 주소 vs 원격 주소 판별
  │   └─ 유효성 검증
  ├─ DHT에서 공개키 조회
  │   ├─ kim의 공개키 ID 계산 (SHA-256)
  │   ├─ DHT get(key) 실행
  │   └─ 캐시 저장 (10분)
  ├─ 공개키 검증
  │   ├─ 서명 확인 (발신자가 정말 kim인가)
  │   ├─ 타임스탐프 확인 (만료되지 않았나)
  │   └─ 신용도 확인 (DHT 평판)
  └─ 결과 반환
      ├─ 유효 공개키 → 메일 송신 가능
      ├─ 만료 공개키 → 경고 (키 순환 확인)
      └─ 없는 주소 → 에러
```

**Components**:
- `AddressParser`: kim@sovereign 분해
- `PublicKeyFetcher`: DHT에서 공개키 조회
- `KeyValidator`: 공개키 검증
- `CacheManager`: 조회 결과 캐시

**Techniques**:
- DNS-like caching (10 minute TTL)
- Fallback mechanisms (broadcast lookup)
- Key fingerprinting (SHA-256 -> 16 hex chars)
- Validity checking (expiration, signature)

**Unforgiving Rules (Module 2)**:
- **Rule 3**: 위조 탐지 99.9%
- **Rule 6**: 대역폭 <10KB per lookup

---

#### Module 3: Reputation-Manager (500 lines)
**Purpose**: 발신자 신용도 관리 (스팸/악의 행위자 차단)

```
신용도 시스템
  ↓ [Reputation Manager]
  ├─ 신용도 점수
  │   ├─ 초기값: 100점
  │   ├─ 스팸 신고당 -10점
  │   ├─ 악성 메일 감지 -5점
  │   ├─ 정상 메일 수신 +1점
  │   └─ 임계: 0점 이하 = 블랙리스트
  ├─ 신용도 저장 (DHT)
  │   ├─ kim's reputation
  │   │   ├─ Score: 95
  │   │   ├─ Reports: ["spam-1", "phishing-2"]
  │   │   ├─ Last update: 2026-03-05
  │   │   └─ TTL: 30 days
  ├─ 신고 메커니즘
  │   ├─ 사용자가 스팸 신고
  │   ├─ 신고는 (발신자, 이유, 증거) 튜플
  │   ├─ DHT에 저장 (다중 서명)
  │   └─ 다른 노드에 복제
  └─ 폭발 반경 차단
      ├─ 발신자 IP 블랙리스트
      ├─ ASN (자율 시스템 번호) 제한
      └─ 기관 도메인 차단 (선택사항)
```

**Components**:
- `ReputationScore`: 신용도 점수 계산
- `ReportHandler`: 스팸 신고 처리
- `ReputationStore`: DHT 저장소
- `BlocklistManager`: IP/ASN 블랙리스트

**Techniques**:
- Weighted scoring (스팸 -10, 악성 -5, 정상 +1)
- Decaying scores (오래된 신고는 점수 회복)
- Multi-signature reports (신뢰성)
- Distributed consensus (quorum voting)

**Unforgiving Rules (Module 3)**:
- **Rule 4**: 주소 일관성 100%
- **Rule 5**: 개인키 노출 0건

---

#### Module 4: Consensus (400 lines)
**Purpose**: DHT 노드 간 동기화 (CRDT 기반)

```
분산 동기화
  ↓ [Consensus Engine]
  ├─ Gossip Protocol
  │   ├─ 자신의 데이터를 이웃에 전파
  │   ├─ 이웃의 데이터를 받아 병합
  │   └─ 주기적 동기화 (30초마다)
  ├─ CRDT (Conflict-free Replicated Data Type)
  │   ├─ Last-Write-Wins (신용도)
  │   ├─ Add-Wins Set (블랙리스트)
  │   └─ Counter (신고 수)
  ├─ 충돌 해결
  │   ├─ 타임스탐프 기반 (최신 우선)
  │   ├─ 서명 기반 (신뢰할 수 있는 발신자)
  │   └─ 쿼럼 투표 (다수결)
  └─ 버전 관리
      ├─ 벡터 시계 (인과관계 추적)
      ├─ 체인 해시 (역사 검증)
      └─ 스냅샷 (체크포인트)
```

**Components**:
- `GossipProtocol`: 피어-투-피어 전파
- `CRDTResolver`: 충돌 해결
- `VectorClock`: 인과관계 추적
- `Snapshot`: 상태 저장/복구

**Techniques**:
- Epidemic gossip (exponential spread)
- Vector clocks (causality tracking)
- CRDT merge functions (deterministic)
- Merkle tree verification (consistency check)

---

## 🧪 Test Plan (30+ tests)

### Group A: DHT-Node (10 tests)
```
✓ test_dht_node_initialization       - 노드 초기화
✓ test_kademlia_distance             - XOR distance 계산
✓ test_k_bucket_management           - k-bucket 관리
✓ test_dht_put_value                 - 데이터 저장 (put)
✓ test_dht_get_value                 - 데이터 조회 (get)
✓ test_dht_value_replication         - 자동 복제 (3+ nodes)
✓ test_dht_node_lookup               - 노드 발견 (FIND_NODE)
✓ test_dht_rpc_handlers              - RPC 핸들러 (PING, STORE)
✓ test_dht_network_resilience        - 죽은 노드 감지 및 복구
✓ test_dht_data_expiration           - TTL 만료 처리
```

### Group B: Address-Resolver (8 tests)
```
✓ test_address_parsing               - kim@sovereign 파싱
✓ test_public_key_fetch              - DHT에서 공개키 조회
✓ test_key_validation                - 공개키 검증
✓ test_key_signature                 - 서명 확인
✓ test_key_expiration                - 만료 처리
✓ test_cache_management              - 조회 캐싱
✓ test_fallback_lookup               - 브로드캐스트 조회
✓ test_invalid_address               - 유효하지 않은 주소 감지
```

### Group C: Reputation-Manager (6 tests)
```
✓ test_reputation_scoring            - 신용도 점수 계산
✓ test_spam_report_handling          - 스팸 신고 처리
✓ test_reputation_store              - DHT 저장
✓ test_blacklist_enforcement         - 블랙리스트 적용
✓ test_reputation_decay              - 신용도 회복
✓ test_ipv4_asn_blocking             - IP/ASN 차단
```

### Group D: Consensus (6 tests)
```
✓ test_gossip_protocol               - 정보 전파
✓ test_crdt_merge                    - 충돌 해결
✓ test_vector_clock                  - 인과관계 추적
✓ test_merkle_consistency             - 일관성 검증
✓ test_snapshot_save                 - 상태 저장
✓ test_snapshot_restore              - 상태 복구
```

---

## 📊 Unforgiving Rules (6 total)

| Rule | Target | Implementation |
|------|--------|-----------------|
| **R1** | DNS lookup <100ms | DHTNode::get_value() |
| **R2** | Recovery <500ms | NetworkResilience::recover() |
| **R3** | Forgery detection 99.9% | AddressResolver::validate_key() |
| **R4** | Address consistency 100% | Consensus::merge_CRDT() |
| **R5** | Key exposure 0 events | ReputationManager::secure_store() |
| **R6** | Bandwidth <10KB | DHTNode::compress_response() |

---

## 📁 File Structure

```
Challenge_15/
├── src/
│   ├── dht_node.rs                  (900 lines)
│   │   ├── KademliaNode
│   │   ├── RoutingTable
│   │   ├── RPCHandler
│   │   └── [10 test functions]
│   │
│   ├── address_resolver.rs          (600 lines)
│   │   ├── AddressParser
│   │   ├── PublicKeyFetcher
│   │   ├── KeyValidator
│   │   └── [8 test functions]
│   │
│   ├── reputation_manager.rs        (500 lines)
│   │   ├── ReputationScore
│   │   ├── ReportHandler
│   │   ├── BlocklistManager
│   │   └── [6 test functions]
│   │
│   ├── consensus.rs                 (400 lines)
│   │   ├── GossipProtocol
│   │   ├── CRDTResolver
│   │   ├── VectorClock
│   │   └── [6 test functions]
│   │
│   └── lib.rs                       (updated)
│
├── Cargo.toml
├── DESIGN.md                        (this file)
└── COMPLETION_REPORT.md             (to be written)
```

---

## 🎯 Implementation Strategy

### Phase A: DHT Foundation (Days 1-2)
1. Kademlia 노드 구현
2. k-bucket 관리
3. RPC 핸들러 (PING, FIND_NODE, FIND_VALUE, STORE)
4. 네트워크 복구 메커니즘

### Phase B: Address Resolution (Day 3)
1. kim@sovereign 주소 파싱
2. DHT에서 공개키 조회
3. 공개키 검증 (서명, 타임스탐프)
4. 조회 캐싱

### Phase C: Reputation System (Day 4)
1. 신용도 점수 계산
2. 스팸 신고 처리
3. 블랙리스트 관리
4. IP/ASN 차단

### Phase D: Consensus (Day 5)
1. Gossip protocol 구현
2. CRDT 병합 함수
3. 벡터 시계
4. 스냅샷 저장/복구

### Phase E: Testing & Optimization (Day 6)
1. <100ms 조회 성능 검증
2. <500ms 복구 성능 검증
3. 메모리 효율성 확인
4. 분산 일관성 검증

---

## 📈 Privacy Guarantees

**Decentralization Properties**:
- ✅ ICANN 의존도 0% (DHT 기반)
- ✅ 중앙 서버 없음 (피어-투-피어)
- ✅ 도메인 차단 불가능 (공개키로 주소화)
- ✅ 주소 영구성 (서버 부재 시에도 유효)
- ✅ 신뢰성 검증 (서명 기반)

**Example Decentralization Flow**:
```
Kim이 주소 등록
  kim@sovereign → RSA-4096 public key fingerprint

DHT에 저장 (복제)
  kim의 공개키 → Node A, B, C (3개 이상)

5년 후, Node A가 사라짐
  → Node B, C가 여전히 kim의 공개키 보유
  → 누군가 kim이 주소를 차단할 수 없음
  → 네트워크가 존재하면 kim@sovereign은 영구 유효

정부/회사가 차단 시도
  → DHT는 분산됨 (중앙 DNS 없음)
  → kim@sovereign은 계속 작동
  → "이제 kim은 자유롭다"
```

---

## 🚀 Expected Outcomes

**Deliverables**:
- 2,400 lines of decentralized naming code
- 30 tests (100% coverage)
- 6 unforgiving rules (all satisfied)
- Full DHT implementation (Kademlia)
- Distributed reputation system
- Byzantine-resistant consensus

**Performance**:
- DNS lookup: <100ms
- Recovery: <500ms
- Address consistency: 100%
- Forgery detection: 99.9%
- Bandwidth: <10KB per lookup

**Security**:
- Decentralized (0% ICANN)
- Permanent addresses (even if server gone)
- Reputation-based trust
- Byzantine-resistant

---

**Next Step**: Implement Sovereign-Naming → Challenge 16 (L0NN-Mail-Sentry)

**Status**: Design approved, ready for implementation 🔧

**Philosophy**: "당신의 주소는 더 이상 구글의 것이 아니다. 공개키 그 자체다."
