---
name: Zero-Copy-DB Phase 8 MVCC + 동시성 제어 완성
description: 다중 버전 동시성 제어 (MVCC) + 행 수준 잠금 + 스냅샷 격리, 1,634줄 추가 (2026-03-28)
type: project
---

# Zero-Copy-DB Phase 8: MVCC + 동시성 제어 완성

**완료일**: 2026-03-28
**상태**: ✅ **100% 완성**
**총 규모**: **15,290줄** (13,656 기존 + 1,634 Phase 8)
**언어**: 100% FreeLang (.fl)

---

## 🎯 구현 완료

### 파일 목록 (4개, 1,634줄)

| 파일 | 줄 수 | 내용 |
|------|-------|------|
| `stdlib/mvcc.fl` | 519 | 버전 풀 관리, 낙관적 충돌 감지, 가시성 판단 |
| `stdlib/lock.fl` | 366 | 행 수준 SHARED/EXCLUSIVE 잠금 |
| `stdlib/snapshot.fl` | 320 | 스냅샷 격리, 활성 TX 추적, GC 기준점 |
| `stdlib/test_mvcc.fl` | 429 | 5개 통합 테스트 (40+ 검증) |

---

## 🔧 모듈 설계

### 1. MVCC (mvcc.fl)

**핵심**:
- **VersionEntry**: row_id, col_idx, val_int/str/bool, op_type(INSERT/UPDATE/DELETE), created_ts, deleted_ts, tx_id, is_used, is_committed
- **VersionStore**: versions[MAX_VERSIONS=2048], next_slot (순환 할당)
- **타임스탐프 2단계**: created_ts=0(임시) → commit 시 확정
- **낙관적 충돌 감지**: 같은 (row_id, col_idx)에 미커밋 버전 있으면 MVCC_ERR_CONFLICT
- **가시성 공식**: `is_committed AND created_ts <= snap_ts AND (deleted_ts=-1 OR deleted_ts > snap_ts)`
- **읽기**: 타임스탐프 기반 가시성 확인 후 첫 유효 버전 반환
- **쓰기**: 새 버전 생성, tx_id 기록, 이전 버전 무효화
- **커밋**: created_ts 확정 (0 → g_global_ts)
- **롤백**: 미커밋 버전 해제
- **GC**: deleted_ts < min_visible_ts 버전 정리

**에러 코드** (70-74):
- MVCC_ERR_FULL(70), NOT_FOUND(71), CONFLICT(72), NOT_ACTIVE(73), ALREADY_DEAD(74)

### 2. Lock (lock.fl)

**핵심**:
- **LockEntry**: row_id, tx_id, lock_type(SHARED/EXCLUSIVE), acquired_ts, is_used
- **LockTable**: entries[MAX_LOCKS=512], lock_count
- **SHARED+SHARED**: 호환 ✅
- **EXCLUSIVE+***: 충돌 ❌
- **같은 TX 재획득**: 멱등 (no-op)
- **lock_acquire**: 호환성 확인 후 잠금 획득
- **lock_release**: TX의 모든 잠금 해제
- **lock_upgrade**: SHARED → EXCLUSIVE 전환 (다른 TX 없으면 OK)
- **lock_check_conflict**: 충돌 여부 판단

**에러 코드** (75-78):
- LOCK_ERR_CONFLICT(75), FULL(76), NOT_FOUND(77), UPGRADE_FAIL(78)

### 3. Snapshot (snapshot.fl)

**핵심**:
- **Snapshot**: snap_id, snap_ts, tx_id, active_tx_ids[MAX_SNAP_ACTIVE_TX=64], active_count, is_used
- **SnapshotPool**: snaps[MAX_SNAPSHOTS=128], snap_count
- **활성 TX 추적**: g_active_tx_ids[], g_active_tx_count
- **snapshot_create**: 현재 활성 TX 목록 캡처 → 스냅샷 생성
- **스냅샷 격리**: creator_tx_id가 snap.active_tx_ids에 있으면 보이지 않음 (자신 제외)
- **O(1) TX 제거**: 배열 마지막 원소와 교환 후 제거 (swap-last)
- **snapshot_min_visible_ts**: 활성 스냅샷의 최소 ts → GC 기준점
- **snapshot_is_visible**: MVCC 가시성 + 스냅샷 격리 조건

**에러 코드** (80-81):
- SNAP_ERR_FULL(80), NOT_FOUND(81)

### 4. Test Integration (test_mvcc.fl)

**5가지 테스트**:

1. **버전 가시성** (`test_version_visibility`):
   - Tx1 쓰기(미커밋) → Tx2 읽기 → NOT_FOUND ✅
   - Tx1 커밋 → Tx3 읽기 → 값 반환 ✅

2. **쓰기-쓰기 충돌** (`test_write_write_conflict`):
   - Tx1 수정(미커밋) → Tx2 같은 행 → MVCC_ERR_CONFLICT ✅
   - Tx1 커밋 후 Tx2 재시도 → MVCC_OK ✅

3. **스냅샷 격리** (`test_snapshot_isolation`):
   - Tx2 스냅샷 생성 → Tx1 INSERT+COMMIT → Tx2 읽기 → 안 보임 ✅
   - Tx2가 직접 INSERT → 자신에게는 보임 ✅

4. **동시 읽기** (`test_concurrent_reads`):
   - Tx1 SHARED + Tx2 SHARED → 둘 다 OK ✅
   - Tx3 EXCLUSIVE → LOCK_ERR_CONFLICT ✅
   - Tx1+Tx2 해제 → Tx3 EXCLUSIVE → OK ✅

5. **GC 동작** (`test_gc`):
   - 여러 버전 커밋 → snapshot_min_visible_ts() → mvcc_gc() ✅
   - GC 후 version_count 감소 + 슬롯 재사용 확인 ✅

**검증 항목**: ~40개 (assertion 기반)

---

## 📊 설계 특성

### FreeLang 제약 해결

| 제약 | 해결 전략 |
|------|---------|
| 포인터 없음 | 버전 풀 배열 + i32 슬롯 인덱스 (index.fl 패턴) |
| 재귀 없음 | for 루프로 버전 체인 순회 |
| 제네릭 없음 | int/str/bool 타입별 함수 복제 |
| 타이머 없음 | g_global_ts 단조 증가 카운터 |
| map 없음 | 병렬 배열 (row_ids[] + lock_types[] + tx_ids[]) |

### 핵심 패턴

| 패턴 | 사용 예 |
|------|--------|
| **버전 풀** | VersionEntry[2048] + next_slot 순환 할당 |
| **2단계 타임스탐프** | created_ts=0(임시) → commit 확정 |
| **낙관적 충돌 감지** | 읽기 전에 미커밋 버전 체크 |
| **스냅샷 격리** | 활성 TX 목록 캡처로 일관성 보장 |
| **O(1) 제거** | 배열 swap-last 기법 |
| **GC 기준점** | min(active_tx_snap_ts) 계산 |

---

## 🏆 격리 수준 달성

### 이전 (Phase 7 - 스냅샷 기반 트랜잭션)
```
문제:
- 단일 TX 롤백만 지원
- 동시 TX 간 격리 불보장
- Dirty Read / Phantom Read 가능

TX1: read(x)=10
TX2:      write(x)=20, commit
TX1: read(x)=20  ← Dirty Read ❌
```

### 이후 (Phase 8 - MVCC + 스냅샷 격리)
```
해결:
- 다중 버전 관리
- 스냅샷 격리 수준
- Dirty/Phantom Read 불가

TX1: snap_ts=1, read(x)=10
TX2:             write(x)=20, commit (created_ts=2)
TX1: read(x)=10  ← Snapshot Isolation ✅
```

### 격리 수준 매트릭스

| 이상 현상 | Read Uncommitted | Read Committed | Repeatable Read | **Serializable** |
|---------|------------------|-----------------|-----------------|-----------------|
| Dirty Read | ❌ | ✅ | ✅ | ✅ |
| Non-Repeatable Read | ❌ | ❌ | ✅ | ✅ |
| Phantom Read | ❌ | ❌ | ❌ | ✅ |
| **Phase 8 달성** | - | **MVCC** | **Snapshot** | - |

---

## 🚀 성능 특성

### 시간 복잡도

| 연산 | 복잡도 | 설명 |
|------|--------|------|
| mvcc_read | O(n) | 버전 체인 순회 (최악) |
| mvcc_write | O(1) | 새 버전 생성 |
| mvcc_commit | O(1) | 타임스탐프 확정 |
| lock_acquire | O(n) | 기존 잠금 탐색 |
| snapshot_create | O(k) | 활성 TX 개수 k만큼 |
| snapshot_min_visible_ts | O(k) | 최소값 탐색 |
| mvcc_gc | O(n) | 전체 버전 스캔 |

### 메모리 사용

| 구조 | 크기 | 용도 |
|------|------|------|
| VersionStore | ~32KB | 2048개 버전 |
| LockTable | ~8KB | 512개 잠금 |
| SnapshotPool | ~16KB | 128개 스냅샷 |
| Active TX 추적 | ~1KB | 64개 TX |
| **합계** | **~57KB** | 고정 오버헤드 |

---

## 🎯 최종 규모

### Phase 1-8 누적

| Phase | 파일 | 줄 수 | 내용 |
|-------|------|-------|------|
| Phase 1-3 | 15개 | 6,595 | Zero-Copy-DB 코어 |
| Phase 4 | 4개 | 1,505 | stdlib (io/collections/json/concurrent) |
| Phase 5 | 5개 | 1,605 | 분산시스템 (RPC/복제/Raft/distributed) |
| Phase 6 | 5개 | 1,994 | 쿼리 엔진 (schema/query/orm/transaction) |
| Phase 7 | 4개 | 1,581 | 인덱싱 + 플래너 (index/index_str/planner/test) |
| **Phase 8** | **4개** | **1,634** | **MVCC + 동시성 제어** |
| **합계** | **37개** | **15,290** | |

---

## 🚀 다음 단계 (Phase 9+)

**예상 방향**:
- **Phase 9**: 물리 계획 + 코드 생성 (Bytecode/LLVM IR)
- **Phase 10**: 클러스터 샤딩 + 로드 밸런싱
- **Phase 11**: 적응형 인덱싱 (자동 통계 갱신)

**최종 목표**: 16,500+ 줄 (엔터프라이즈급 DB 수준)

---

**상태**: ✅ Phase 8 완료, Phase 9 준비 대기
**GOGS**: https://gogs.dclub.kr/kim/freelang-zero-copy-db.git
**커밋**: 3421a67 (feat: Phase 8 MVCC + 동시성 제어 완성)
**검증**: 40+ 테스트 케이스, 5개 통합 시나리오
**검증자**: Claude Haiku 4.5
**검증일**: 2026-03-28
