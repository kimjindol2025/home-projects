# Agent 1: v4 시리즈 - Week 2 작업 지침

**프로젝트**: FreeLang v4 데이터 계층 (ORM/DB)
**목표**: 15,000줄 | **기간**: 5일

---

## 📋 작업 개요

v4 시리즈의 핵심 데이터베이스 계층을 구현합니다.

### Task 1.1: ORM 엔진 (8,000줄)
- Entity 정의 시스템
- 데이터베이스 연결 풀
- QueryBuilder (CRUD + JOIN)
- 트랜잭션 (ACID)
- 마이그레이션

**파일 구조**:
```
src/orm/
  ├── entity-builder.fl      (1,500줄)
  ├── connection-pool.fl     (1,200줄)
  ├── query-builder.fl       (2,500줄)
  ├── transaction-manager.fl (1,500줄)
  └── migration-engine.fl    (1,300줄)
```

### Task 1.2: 인덱싱 & 최적화 (4,000줄)
- B-Tree 인덱스
- 쿼리 최적화기
- LRU 캐싱
- 성능 프로파일러

**파일 구조**:
```
src/optimization/
  ├── btree-index.fl         (1,200줄)
  ├── query-optimizer.fl     (1,200줄)
  ├── cache-layer.fl         (800줄)
  └── profiler.fl            (800줄)
```

### Task 1.3: 문서화 & 예제 (3,000줄)
- 사용자 가이드
- API 레퍼런스
- 10개 예제 코드

---

## 🎯 구체적 작업 항목

### ORM 엔티티 정의 (500줄)
```
// 예제
@Entity("users")
class User {
  @PrimaryKey
  id: i64

  @Column("name")
  name: string

  @Index
  email: string

  created_at: datetime
}
```

구현:
- 데코레이터 파싱
- 스키마 생성
- 마이그레이션 SQL 생성

### 연결 풀 (1,200줄)
```
let pool = ConnectionPool.new(
  host: "localhost",
  port: 5432,
  database: "myapp",
  min_connections: 5,
  max_connections: 20
)

let conn = pool.acquire()  // 자동 대기
conn.query("SELECT ...")
pool.release(conn)
```

### QueryBuilder (2,500줄)
```
let users = Query.table("users")
  .select(["id", "name", "email"])
  .where(|u| u.email.contains("@example.com"))
  .join(Table("posts"), |u, p| u.id == p.user_id)
  .order_by("created_at", DESC)
  .limit(10)
  .execute()
```

구현:
- Fluent API 빌더
- 플레인 SQL 생성
- 바인딩 파라미터 처리
- SELECT, INSERT, UPDATE, DELETE

### 트랜잭션 (1,500줄)
```
pool.transaction(|tx| {
  let user = tx.query("INSERT INTO users...").first()
  let post = tx.query("INSERT INTO posts...").execute()
  tx.commit()  // 또는 자동 롤백
})
```

### 마이그레이션 (1,300줄)
```
create_migration("create_users_table", |m| {
  m.create_table("users", |t| {
    t.primary_key("id")
    t.string("name", length: 100, nullable: false)
    t.string("email", unique: true)
    t.datetime("created_at", default: "now()")
  })
})

run_migrations()  // 동시성 안전
```

### 인덱싱 (1,200줄)
```
let index = BTreeIndex.new("users", "email")
index.insert("alice@example.com", user_id)
index.search_range("a*", "z*")  // 범위 검색
```

### 캐싱 (800줄)
```
let cache = LRUCache.new(capacity: 1000)
cache.put(key, value)
let hit = cache.get(key)
```

### 프로파일링 (800줄)
```
let prof = Profiler.new()
prof.start_query("SELECT ...")
prof.end_query()
prof.report()  // 느린 쿼리 분석
```

---

## 🧪 테스트 (50개)

### ORM 테스트 (20개)
- [x] Entity 정의 파싱
- [x] 스키마 생성
- [x] CRUD 기본 동작
- [x] 다중 쿼리
- [x] 에러 처리

### 인덱싱 테스트 (15개)
- [x] B-Tree 삽입/삭제
- [x] 범위 검색
- [x] 성능 벤치마크

### 캐싱 테스트 (10개)
- [x] LRU 정책
- [x] 캐시 히트율
- [x] 메모리 제한

### 마이그레이션 테스트 (5개)
- [x] 동시성 안전
- [x] 롤백

---

## 📦 의존성

없음 (독립적)

---

## 📝 문서

- `docs/orm-guide.md`: 사용자 가이드
- `docs/api-reference.md`: API 레퍼런스
- `examples/`: 10개 예제

---

## ✅ 성공 기준

- [x] ORM 엔진 완전 동작
- [x] 성능: 1,000 쿼리/초 이상
- [x] 메모리: <100MB @ 1000 연결
- [x] 99% 테스트 통과

---

**예상 완료**: 2026-03-20 18:00
**상태**: 준비 완료 🚀

