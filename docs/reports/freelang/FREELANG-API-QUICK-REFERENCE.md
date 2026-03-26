# FreeLang REST API 빠른 참조 가이드

**버전**: v1.0
**마지막 업데이트**: 2026-03-02

---

## 📚 관련 문서

| 문서 | 내용 | 목적 |
|------|------|------|
| **FREELANG-REST-API-DESIGN.md** | 설계 원칙, 패턴, 이론 | 아키텍처 이해 |
| **FREELANG-REST-API-EXAMPLES.md** | 실제 구현 코드 | 실전 코딩 |
| **FREELANG-API-QUICK-REFERENCE.md** | 빠른 참조 & 치트시트 | ← 현재 문서 |

---

## 🚀 5분 시작 가이드

### 1️⃣ 기본 구조

```freelang
# Step 1: 모델 정의
struct Todo
  id: number
  title: string
  completed: bool

# Step 2: 핸들러 작성
fn handleGetTodos
  input: HttpRequest
  output: Result<HttpResponse, Error>
  intent: "모든 Todo 반환"
  do
    # 구현...

# Step 3: 서버 시작
fn main
  input: void
  output: void
  do
    startServer(routes, 8080)
```

### 2️⃣ REST 메서드 매핑

```
┌─────────┬────────────┬──────────┬───────┐
│ Method  │ Path       │ 의미     │ Status│
├─────────┼────────────┼──────────┼───────┤
│ GET     │ /todos     │ 목록     │ 200   │
│ POST    │ /todos     │ 생성     │ 201   │
│ GET     │ /todos/:id │ 조회     │ 200   │
│ PUT     │ /todos/:id │ 수정     │ 200   │
│ DELETE  │ /todos/:id │ 삭제     │ 204   │
└─────────┴────────────┴──────────┴───────┘
```

### 3️⃣ 표준 응답 포맷

```json
// ✅ 성공 (200, 201)
{
  "success": true,
  "data": { /* 실제 데이터 */ },
  "timestamp": "2026-03-02T10:30:00Z"
}

// ❌ 오류 (400, 404, 500)
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found"
  },
  "timestamp": "2026-03-02T10:30:00Z"
}
```

---

## 📝 핸들러 템플릿

### GET (조회)

```freelang
fn handleGetTodos
  input: (store: TodoStore, request: HttpRequest)
  output: Result<HttpResponse, ApiError>
  intent: "모든 Todo를 페이지네이션과 함께 반환"
  do
    # 1. 쿼리 파라미터 파싱
    let page = parseNumber(request.queryParams.get("page")).getOrElse(1)

    # 2. 데이터 조회
    let items = values(store.todos)

    # 3. 페이지네이션
    let paged = slice(items, (page-1)*10, page*10)

    # 4. JSON 응답
    let body = json({
      "success": true,
      "data": paged,
      "timestamp": getCurrentTime()
    })

    return Ok({
      status: 200,
      headers: { "Content-Type": "application/json" },
      body: body
    })
```

### POST (생성)

```freelang
fn handleCreateTodo
  input: (store: TodoStore, request: HttpRequest)
  output: Result<(TodoStore, HttpResponse), ApiError>
  intent: "새로운 Todo를 생성"
  do
    # 1. 요청 본문 파싱
    let body = parseJson(request.body).ok_or(
      ApiError::InvalidJson("Invalid JSON")
    )

    # 2. 입력 검증
    let title = body.get("title").unwrap_or("")
    if title == ""
      return Err(ApiError::ValidationError("title required"))

    # 3. 데이터 생성
    let (newStore, id) = nextId(store)
    let todo = Todo {
      id: id,
      title: title,
      completed: false
    }
    newStore.todos[id] = todo

    # 4. 응답
    let jsonBody = json({
      "success": true,
      "data": todo,
      "timestamp": getCurrentTime()
    })

    return Ok((newStore, {
      status: 201,
      headers: {
        "Content-Type": "application/json",
        "Location": "/todos/" + str(id)
      },
      body: jsonBody
    }))
```

### PUT (수정)

```freelang
fn handleUpdateTodo
  input: (store: TodoStore, id: number, request: HttpRequest)
  output: Result<(TodoStore, HttpResponse), ApiError>
  intent: "Todo를 부분 업데이트"
  do
    # 1. 존재 확인
    match store.todos.get(id)
      None => Err(ApiError::NotFound("Not found"))
      Some(todo) => {
        # 2. 요청 파싱
        let body = parseJson(request.body).ok_or(
          ApiError::InvalidJson("Invalid JSON")
        )

        # 3. 필드 업데이트
        let updated = todo
        if body.contains("title")
          updated.title = body.get("title").unwrap_or("")
        if body.contains("completed")
          updated.completed = body.get("completed") == "true"

        store.todos[id] = updated

        # 4. 응답
        return Ok((store, {
          status: 200,
          headers: { "Content-Type": "application/json" },
          body: json({
            "success": true,
            "data": updated,
            "timestamp": getCurrentTime()
          })
        }))
      }
```

### DELETE (삭제)

```freelang
fn handleDeleteTodo
  input: (store: TodoStore, id: number, request: HttpRequest)
  output: Result<TodoStore, ApiError>
  intent: "Todo를 삭제"
  do
    # 1. 존재 확인
    if !store.todos.contains(id)
      return Err(ApiError::NotFound("Not found"))

    # 2. 삭제
    store.todos.remove(id)

    # 3. 응답 (본문 없음, 204)
    return Ok(store)
```

---

## 🔗 라우팅 패턴

### 경로 매칭

```freelang
fn matchPattern
  input: (pattern: string, path: string)
  output: bool
  intent: "URL 패턴 매칭"
  do
    # /todos/:id와 /todos/123 매칭
    let parts = split(path, "/")
    let patterns = split(pattern, "/")

    if parts.len() != patterns.len()
      return false

    for i in 0..parts.len()
      if patterns[i].startsWith(":")
        continue  # 파라미터
      if parts[i] != patterns[i]
        return false
    return true
```

### 라우터 정의

```freelang
fn createRouter
  input: void
  output: array<Route>
  intent: "모든 라우트 정의"
  do
    return [
      { pattern: "/todos", method: "GET", handler: handleGetTodos },
      { pattern: "/todos", method: "POST", handler: handleCreateTodo },
      { pattern: "/todos/:id", method: "GET", handler: handleGetTodoById },
      { pattern: "/todos/:id", method: "PUT", handler: handleUpdateTodo },
      { pattern: "/todos/:id", method: "DELETE", handler: handleDeleteTodo }
    ]
```

---

## ⚠️ 에러 처리

### 에러 타입

| 타입 | 상태 | 사용 경우 |
|------|------|---------|
| `ValidationError` | 400 | 입력 검증 실패 |
| `NotFound` | 404 | 자원 없음 |
| `Conflict` | 409 | 충돌 (중복 등) |
| `ServerError` | 500 | 서버 오류 |
| `InvalidJson` | 400 | JSON 파싱 실패 |

### 에러 응답

```freelang
fn createErrorResponse
  input: ApiError
  output: HttpResponse
  intent: "에러를 HTTP 응답으로 변환"
  do
    let (status, code, message) = match error
      ValidationError(msg) => (400, "VALIDATION_ERROR", msg)
      NotFound(msg) => (404, "NOT_FOUND", msg)
      Conflict(msg) => (409, "CONFLICT", msg)
      ServerError(msg) => (500, "SERVER_ERROR", msg)
      InvalidJson(msg) => (400, "INVALID_JSON", msg)

    return {
      status: status,
      headers: { "Content-Type": "application/json" },
      body: json({
        "success": false,
        "error": { "code": code, "message": message },
        "timestamp": getCurrentTime()
      })
    }
```

---

## 🧪 테스트 패턴

### 단위 테스트

```freelang
fn testCreateTodo
  intent: "Todo 생성 테스트"
  do
    let store = initTodoStore()
    let request = {
      method: "POST",
      path: "/todos",
      body: json({ "title": "Test", "description": "" })
    }

    match handleCreateTodo(store, request)
      Ok((newStore, response)) => {
        assert(response.status == 201)
        assert(newStore.nextId == 2)
        print("✓ PASS")
      }
      Err(e) => print("✗ FAIL")
```

### 통합 테스트

```freelang
fn testTodoWorkflow
  intent: "전체 워크플로우 테스트"
  do
    let store = initTodoStore()

    # 1. 생성
    match handleCreateTodo(store, createReq)
      Ok((store1, resp1)) => {
        assert(resp1.status == 201)

        # 2. 조회
        match handleGetTodoById(store1, resp1.data.id)
          Ok(resp2) => {
            assert(resp2.status == 200)

            # 3. 수정
            match handleUpdateTodo(store1, resp1.data.id, updateReq)
              Ok((store2, resp3)) => {
                assert(resp3.status == 200)
                print("✓ WORKFLOW PASS")
              }
```

---

## 📊 데이터 구조 치트시트

### 요청

```freelang
struct HttpRequest
  method: string           # "GET", "POST", "PUT", "DELETE"
  path: string            # "/todos/123"
  headers: map<string>    # { "Content-Type": "application/json" }
  body: string            # JSON 또는 텍스트
  queryParams: map<string># ?page=1&size=10
```

### 응답

```freelang
struct HttpResponse
  status: number          # 200, 201, 400, 404, 500
  headers: map<string>    # { "Content-Type": "application/json" }
  body: string            # JSON 또는 텍스트
```

### 저장소

```freelang
struct TodoStore
  todos: map<number, Todo>    # ID → Todo
  nextId: number              # 다음 사용 ID
```

---

## 🔧 유틸리티 함수

| 함수 | 용도 | 예시 |
|------|------|------|
| `parseNumber(string)` | 문자열→숫자 | `parseNumber("123")` → `Some(123)` |
| `parseJson(string)` | JSON 파싱 | `parseJson("{...}")` → `Ok(map)` |
| `json(map)` | 맵→JSON | `json({"a": "b"})` → `"{...}"` |
| `split(string, sep)` | 문자열 분할 | `split("a/b", "/")` → `["a", "b"]` |
| `slice(array, i, j)` | 배열 부분 | `slice([1,2,3], 1, 2)` → `[2]` |
| `filter(array, pred)` | 필터링 | `filter([1,2,3], x > 1)` → `[2,3]` |
| `values(map)` | 맵→배열 | `values({a:1, b:2})` → `[1,2]` |
| `getCurrentTime()` | 현재 시간 | → `"2026-03-02T10:30:00Z"` |

---

## 💡 Best Practices

### ✅ DO

```freelang
# ✓ 명확한 intent
fn handleGetTodos
  intent: "모든 Todo를 페이지네이션과 함께 반환"
  do ...

# ✓ 입력 검증
if title == ""
  return Err(ValidationError("title required"))

# ✓ 일관된 응답
return Ok({
  status: 200,
  headers: { "Content-Type": "application/json" },
  body: json(...)
})

# ✓ 명확한 에러 메시지
Err(NotFound("Todo with id " + str(id) + " not found"))
```

### ❌ DON'T

```freelang
# ✗ 모호한 intent
fn handleTodos
  intent: "Handle request"
  do ...

# ✗ 검증 없음
let title = body.get("title")  # 빈 문자열 허용?

# ✗ 일관성 없는 응답
return { status: 200, body: "OK" }

# ✗ 모호한 에러
Err(ServerError("Error"))
```

---

## 🌐 HTTP 상태 코드

| 코드 | 의미 | 사용 |
|------|------|------|
| **200** | OK | 성공적 조회/수정 |
| **201** | Created | 자원 생성 |
| **204** | No Content | 성공 (응답 본문 없음) |
| **400** | Bad Request | 입력 검증 실패 |
| **401** | Unauthorized | 인증 필요 |
| **403** | Forbidden | 권한 없음 |
| **404** | Not Found | 자원 없음 |
| **409** | Conflict | 충돌 |
| **500** | Server Error | 서버 오류 |
| **503** | Service Unavailable | 서비스 불가 |

---

## 🔐 보안 체크리스트

- [ ] **입력 검증**: 모든 사용자 입력 검증
- [ ] **XSS 방지**: HTML 이스케이프
- [ ] **SQL 주입 방지**: 파라미터화된 쿼리
- [ ] **인증**: 토큰 검증
- [ ] **인가**: 권한 확인
- [ ] **HTTPS**: SSL/TLS 사용
- [ ] **CORS**: 크로스 오리진 정책
- [ ] **속도 제한**: Rate limiting
- [ ] **로깅**: 감사 로그
- [ ] **에러 처리**: 민감 정보 노출 금지

---

## 📈 성능 가이드

```
권장 응답 시간:
├─ GET  /items           < 50ms  (페이지네이션)
├─ GET  /items/:id       < 10ms  (직접 조회)
├─ POST /items           < 20ms  (생성)
├─ PUT  /items/:id       < 20ms  (수정)
└─ DELETE /items/:id     < 15ms  (삭제)

최적화 전략:
├─ 인덱싱 (ID, 상태 등)
├─ 캐싱 (자주 접근 데이터)
├─ 배치 처리 (여러 항목 한 번에)
└─ 비동기 작업 (무거운 작업은 큐로)
```

---

## 🧩 마이크로서비스 아키텍처

### 서비스 분리

```
┌──────────────────────────────────────┐
│        API Gateway                   │
│  (라우팅, 인증, 로깅)                 │
└─────────────┬────────────────────────┘
              │
    ┌─────────┼─────────┬──────────┐
    │         │         │          │
┌───▼──┐ ┌───▼──┐ ┌───▼──┐ ┌───▼──┐
│ Todo │ │ User │ │Post  │ │Comment│
│API   │ │API   │ │API   │ │API   │
└──────┘ └──────┘ └──────┘ └──────┘
```

### 서비스 간 통신

```freelang
# Service A (Todo API)
fn callUserService
  input: userId: number
  output: Result<User, Error>
  intent: "User Service 호출"
  do
    let response = httpGet("http://user-service:8081/users/" + str(userId))
    return parseJson(response.body)
```

---

## 📚 학습 순서

1. **기초** (1일)
   - HTTP 메서드 이해
   - 표준 응답 포맷
   - CRUD 핸들러

2. **중급** (2-3일)
   - 라우팅
   - 에러 처리
   - 데이터 검증
   - 페이지네이션

3. **고급** (4-5일)
   - 인증/인가
   - 미들웨어
   - 캐싱 및 최적화
   - 통합 테스트

4. **전문가** (6-7일)
   - 마이크로서비스
   - API 버전 관리
   - OpenAPI/Swagger
   - 분산 트레이싱

---

## 🔗 관련 링크

- 📖 [Full Design Guide](./FREELANG-REST-API-DESIGN.md)
- 💻 [Implementation Examples](./FREELANG-REST-API-EXAMPLES.md)
- 📋 [FreeLang Spec](./FREELANG-v1.0-SPEC.md)
- 🚀 [Getting Started](./FREELANG-GETTING-STARTED.md)

---

**버전**: v1.0
**마지막 업데이트**: 2026-03-02
**상태**: ✅ Phase B 설계 완료
