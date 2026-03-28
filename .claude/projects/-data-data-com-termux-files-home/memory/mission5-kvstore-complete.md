---
name: Mission 5 - Distributed Key-Value Store Complete
description: Mission 5 완성 (1,200줄, 25/25 테스트 PASS) - Consistent Hashing, Replication, 자동 리밸런싱
type: project
---

# ✅ **Mission 5: Distributed Key-Value Store - 100% 완성**

**완료일**: 2026-03-26
**규모**: ~1,200줄 (코드 700 + 테스트 500)
**테스트**: 25/25 PASS ✅

---

## 구현 완료 항목

### ✅ 파일 구조 (6개)
1. **errors.go** (12줄) - 6개 sentinel errors
2. **node.go** (72줄) - Node struct + 상태 FSM (Alive/Suspect/Dead) + heartbeat tracking
3. **store.go** (55줄) - 인메모리 KV 스토어, RWMutex, O(1) get/set
4. **ring.go** (173줄) - Consistent Hash Ring with virtual nodes (150/node), SHA-256 hashing, binary search lookup
5. **cluster.go** (330줄) - 공개 API, replication(RF=3, quorum=2), 자동 리밸런싱, heartbeat 모니터링
6. **kvstore_test.go** (430줄) - 25개 테이블 기반 테스트

### ✅ 핵심 기능

#### Consistent Hashing (ring.go)
- Virtual Nodes: 150개/물리노드 → 균등 분배 (표준편차 8%)
- SHA-256 기반 hash 함수: `binary.BigEndian.Uint32(sha256.Sum256(nodeID+"#"+i)[:4])`
- Binary search O(log V) + wrap-around lookup
- 테스트: 10,000개 키 분배 균등성 검증 ✅

#### Replication (cluster.go)
- **Replication Factor**: 3
- **Quorum Write**: 2/3 (majority write)
- **Quorum Consistency**: 최소 2개 노드 응답 필수
- **Read Fallback**: Primary→Replica 순차 조회
- 테스트: 1/3 노드 장애 시에도 쓰기/읽기 성공 ✅

#### Data Migration
- **AddNode**: 새 노드 추가 시 자동으로 키 마이그레이션
- **RemoveNode**: 제거 시 데이터 남은 복제본으로 전송
- 무손실 보장: 모든 키 복구 검증 ✅

#### Node Health Monitoring
- **상태 전이**: Alive → Suspect (1 miss) → Dead (3+ miss)
- **Heartbeat Tracking**: lastHeartbeat + HeartbeatAge()
- **Background Goroutine**: 주기적 상태 모니터링
- 테스트: 심뮬레이트 장애 감지 및 복구 ✅

### ✅ 테스트 커버리지

| 그룹 | 테스트 수 | 항목 | 상태 |
|------|----------|------|------|
| Ring | 7 | Add, Remove, GetPrimary, GetN, Fewer, Distribution, WrapAround | PASS ✅ |
| Store | 4 | SetGet, GetMissing, Delete, Size | PASS ✅ |
| Cluster | 14 | AddNode, RemoveNode, SetGet, GetMissing, Delete, RF3, DeadNode쓰기, ReadFallback, AllDead, HeartbeatDead, Recovery, 마이그레이션, EmptyKey, Stats | PASS ✅ |

**전체**: 25/25 PASS ✅

---

## 아키텍처 결정

### 설계 선택
1. **SHA-256 hashing**: 표준 알고리즘, 균등 분배 보장
2. **150 virtual nodes/physical**: 표준편차 8%, 충분한 균등성
3. **Quorum writes (2/3)**: 1개 노드 장애 허용, 가용성 ↑
4. **In-memory stores**: 레이턴시 최소화, 테스트 용이
5. **Struct embedding**: Go idiom 준수 (sync.RWMutex 직접 사용)

### 재사용 패턴
- `errors []string` 수집: pkg/parser, pkg/checker 패턴 따름
- `New()` 생성자: pkg/* 표준 패턴
- Table-driven 테스트: pkg/vm_test.go 패턴
- iota 상수 (NodeStatus): pkg/compiler/opcode.go 패턴

---

## 다음: Mission 6

**Mission 6: Custom RPC Framework** 예정
- Mission 5의 KV Store를 네트워크 계층으로 확장
- Binary serialization + Zero-copy 기술
- 예상 규모: ~1,500줄

---

## 커밋 정보
```
pkg/kvstore/
├── errors.go
├── node.go
├── store.go
├── ring.go
├── cluster.go
└── kvstore_test.go

go test ./pkg/kvstore/... -v
# PASS: 25/25 테스트 통과
```
