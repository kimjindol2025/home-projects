---
name: Mission 3 - LSM 완성 (Phase 1-4)
description: Log-Structured Merge Tree 완전 구현, 1,670줄 코드 + 54/54 테스트
type: project
---

# 🎉 Mission 3: Log-Structured Merge Tree (LSM) 완성!!

**상태**: ✅ 100% 완료
**규모**: ~1,670줄 코드 + 630줄 테스트 = 2,300줄
**테스트**: 54/54 PASS ✅
**GOGS**: https://gogs.dclub.kr/kim/freelang-lsm.git
**완료일**: 2026-03-27

---

## Phase 별 완성 내용

### Phase 1: SkipList (380줄 + 180줄 테스트)
```go
type SkipList struct {
    header *Node
    level int
    mu sync.RWMutex
}

Key Features:
- 레벨 16 (확률 기반 O(log N))
- Put, Get, Delete, Range, Rank, RangeByRank
- 메모리 사용량 추적
- 레벨 분포 최적화
```

**9/9 테스트 PASS**:
- 기본 연산 (Put/Get/Delete)
- Range 쿼리
- Concurrency 안전성
- 메모리 추정

---

### Phase 2: WAL + MemTable (230줄 + 180줄 + 180줄 테스트)

#### WAL (Write-Ahead Log)
```go
type WAL struct {
    file *os.File
    buffer []byte
    maxSize int64
}

Features:
- [1B: type][4B: key_len][4B: value_len][key][value] 인코딩
- 버퍼링 + 파일 로테이션
- 복구 가능한 형식
```

#### MemTable
```go
type MemTable struct {
    sl *SkipList
    wal *WAL
    maxSize int64
    isFrozen bool
}

Features:
- SkipList 기반 고속 읽기
- WAL 기반 내구성
- Freeze 상태로 플러시 신호
```

**24/24 테스트 PASS**:
- Put/Get/Delete
- Range 쿼리
- Freeze 상태 관리
- 메모리 추정
- WAL 동기화

---

### Phase 3: SSTable (580줄 + 150줄 테스트)

```go
type SSTable struct {
    path string
    file *os.File
    indexOffset int64
    entries []*IndexEntry
}

File Format:
[Data Blocks][Index Block][Footer]
  - Data Block: [4B key_len][4B value_len][key][value]...
  - Index Block: [4B count][key_len][key][8B offset]...
  - Footer: [8B index_offset]
```

**Key Features**:
- 블록 기반 구조 (4KB 기본)
- 이진 검색 인덱스 (O(log N) 조회)
- Range 쿼리 지원
- 블록 단위 읽기 (I/O 최적화)

**10/10 테스트 PASS**:
- 생성, 열기, Get, Range
- 다중 블록 처리
- 큰 파일 (1000+ 엔트리)
- 엣지 케이스 (빈 리스트, 없는 키)

---

### Phase 4: Compactor (300줄 + 380줄 테스트)

```go
type Compactor struct {
    ssTables []*SSTable
    config *CompactorConfig
    stopCh chan struct{}
    isRunning bool
    lastCompact time.Time
}

type CompactorConfig struct {
    CompactInterval time.Duration
    MaxLevelSize int64
    MaxSSTables int
    BlockSize int64
}
```

**Key Functions**:
- `Start/Stop`: 백그라운드 서비스
- `tryCompact`: 타이머 기반 자동 트리거
- `compact`: 다중 SSTable 병합 + 중복 제거
- `deduplicateEntries`: 마지막 값 유지
- `GetStats`: 상태 조회

**Trigger Conditions** (shouldCompact):
```
1. len(ssTables) >= MaxSSTables
2. totalSize > MaxLevelSize
3. lastCompact + 2*CompactInterval 경과 && len > 1
```

**11/11 테스트 PASS**:
- SSTable 추가, 크기 계산
- 트리거 조건 검증
- 실제 컴팩션 + 중복 제거
- 백그라운드 실행
- 통계 조회
- 다중 SSTable 병합 (5개)
- 엣지 케이스 (빈 리스트)

---

## 누적 통계

| 항목 | 수치 |
|------|------|
| **총 코드** | 1,670줄 |
| **총 테스트** | 630줄 |
| **테스트 케이스** | 54/54 PASS ✅ |
| **컴포넌트** | 4개 (SkipList, WAL, MemTable, SSTable, Compactor) |
| **Git 커밋** | 4개 |

---

## 핵심 구현 포인트

### 1. 동시성 안전성
```go
// 모든 구조체가 sync.RWMutex 사용
mu sync.RWMutex
defer mu.Unlock()  // 항상 defer로 보호
```

### 2. 이진 형식 (Binary Wire Protocol)
```go
// 일관된 인코딩: 길이 + 데이터
binary.BigEndian.PutUint32(buf, uint32(len(key)))
binary.BigEndian.PutUint64(buf, uint64(offset))
```

### 3. 메모리 효율성
```go
// SkipList: 확률 기반 레벨 (평균 1.44 포인터/노드)
// SSTable: 인덱스만 메모리 (데이터는 디스크)
// Compactor: 배치 처리 (스트리밍 없음)
```

### 4. 에러 처리
```go
// 파일 I/O는 모두 에러 체크
if err := st.file.WriteAt(data, offset); err != nil {
    return fmt.Errorf("failed to write: %w", err)
}
```

---

## LSM 데이터 흐름

```
Write 요청
    ↓
MemTable.Put() [SkipList + WAL]
    ↓
MemTable.IsFull()?
    ├─ No  → 계속 쓰기
    └─ Yes → Flush to SSTable
         ↓
      [SSTable 생성]
         ↓
Compactor.tryCompact()
    ↓
SSTables >= MaxSSTables?
    ├─ No  → 대기
    └─ Yes → Merge + Deduplicate
         ↓
      [새로운 SSTable 생성]

Read 요청
    ↓
MemTable.Get() [SkipList 검색] → Found? Return
    ↓ Not Found
SSTables.Get() [이진 검색] → Found? Return
    ↓
Not Found
```

---

## GOGS 배포

```bash
git remote: https://gogs.dclub.kr/kim/freelang-lsm.git
Last Commit: a79e4b8 (Phase 4 컴팩션)
Status: ✅ Pushed
```

---

## 다음 단계 옵션

### Option A: Phase 5 (DB Engine)
- LSM을 완전한 데이터베이스로 확장
- Get, Put, Delete, Scan, Compact 통합
- 예상: 1,000줄 + 테스트

### Option B: Mission 4 (IaC)
- Infrastructure as Code Engine
- Terraform 스타일 선언형 구문
- 예상: 2,000줄 + 테스트

---

## 검증 명령어

```bash
# 모든 테스트 실행
go test ./internal/lsm -v

# 특정 테스트만
go test ./internal/lsm -run TestCompactor -v

# 벤치마크 (Option)
go test ./internal/lsm -bench=. -benchmem
```

---

## 핵심 학습점

1. **프로빙**: SkipList의 확률 기반 높이가 O(log N)을 보장
2. **인덱싱**: SSTable의 이진 검색 인덱스로 O(log N) 조회
3. **배치**: Compactor의 배치 처리로 쓰기 증폭 최소화
4. **동시성**: RWMutex의 정확한 사용으로 안전성 보장
5. **내구성**: WAL로 크래시 복구 가능

---

**상태**: ✅ 완료 (54/54 테스트 PASS, GOGS 배포 완료)

