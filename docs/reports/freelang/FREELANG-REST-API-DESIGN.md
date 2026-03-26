# FreeLang REST API 설계 가이드

**작성일**: 2026-03-02
**버전**: v1.0
**상태**: 🟢 Phase B 설계
**목표**: AI-친화적 REST API 구축 패턴

---

## 📖 목차

1. [핵심 설계 원칙](#1-핵심-설계-원칙)
2. [HTTP 핸들링 구조](#2-http-핸들링-구조)
3. [라우팅 패턴](#3-라우팅-패턴)
4. [요청/응답 설계](#4-요청응답-설계)
5. [예제: 간단한 TODO API](#5-예제-간단한-todo-api)
6. [예제: 블로그 API](#6-예제-블로그-api)
7. [에러 처리 전략](#7-에러-처리-전략)
8. [테스트 패턴](#8-테스트-패턴)

---

## 1. 핵심 설계 원칙

### 1.1 REST의 4가지 기본 원칙

```
Resource-Based URLs:
  - GET    /todos            → 모든 todo 조회
  - POST   /todos            → 새 todo 생성
  - GET    /todos/{id}       → 특정 todo 조회
  - PUT    /todos/{id}       → todo 수정
  - DELETE /todos/{id}       → todo 삭제

Standard Status Codes:
  - 200: OK (요청 성공)
  - 201: Created (자원 생성 성공)
  - 400: Bad Request (잘못된 요청)
  - 404: Not Found (자원 없음)
  - 500: Server Error (서버 오류)

Stateless Design:
  - 각 요청은 독립적
  - 이전 요청 정보 불필요
  - 서버 상태 최소화

JSON Format:
  - 표준 JSON 메시지
  - 일관된 응답 스키마
  - 타입 정보 명시
```

### 1.2 FreeLang REST API의 특징

| 특징 | 설명 |
|------|------|
| **Intent-Clear** | 각 엔드포인트의 의도가 명확 |
| **Type-Safe** | 요청/응답 타입이 정의됨 |
| **Error-Explicit** | 모든 오류가 Result로 표현 |
| **Composable** | 작은 핸들러를 조합 |

---

## 2. HTTP 핸들링 구조

### 2.1 기본 아키텍처

```
┌─────────────────────────────────────────┐
│         HTTP Request                    │
│  GET /todos/123 HTTP/1.1                │
│  Host: api.example.com                  │
│  Content-Type: application/json         │
└─────────────────────────────────────────┘
              ↓
        ┌─────────────┐
        │   Router    │
        │   /todos/:id│
        └─────────────┘
              ↓
    ┌──────────────────────┐
    │  Handler (Function)  │
    │  getTodoById(id)     │
    └──────────────────────┘
              ↓
    ┌──────────────────────┐
    │  Business Logic      │
    │  Query Database      │
    └──────────────────────┘
              ↓
    ┌──────────────────────┐
    │  Response Handler    │
    │  ToJSON Response     │
    └──────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│      HTTP Response                      │
│  HTTP/1.1 200 OK                        │
│  Content-Type: application/json         │
│  { "id": 123, "title": "..." }          │
└─────────────────────────────────────────┘
```

### 2.2 Phase B 런타임 기반 HTTP 핸들러

```freelang
# HTTP 요청을 처리하는 핸들러 기본 구조
struct HttpRequest
  method: string           # "GET", "POST", "PUT", "DELETE"
  path: string            # "/todos/123"
  headers: map<string>    # { "Content-Type": "application/json" }
  body: string            # 요청 본문 (JSON 등)

struct HttpResponse
  status: number          # 200, 201, 400, 404, 500
  headers: map<string>    # { "Content-Type": "application/json" }
  body: string            # 응답 본문

# 핸들러 함수의 기본 형태
fn handleRequest
  input: HttpRequest
  output: Result<HttpResponse, Error>
  intent: "HTTP 요청을 처리하고 응답 생성"
  do
    match request.method
      "GET" => handleGet(request)
      "POST" => handlePost(request)
      "PUT" => handlePut(request)
      "DELETE" => handleDelete(request)
      _ => Err(Error::InvalidMethod)
```

---

## 3. 라우팅 패턴

### 3.1 단순 라우터 패턴

```freelang
# 라우터: URL을 핸들러에 매핑
struct Route
  pattern: string         # "/todos/:id" → /todos/123 매칭
  method: string          # "GET", "POST", etc.
  handler: fn             # 실제 처리 함수

fn matchRoute
  input: HttpRequest
  routes: array<Route>
  output: Option<Route>
  intent: "요청에 맞는 라우트 찾기"
  do
    for route in routes
      if route.method == request.method && matchPattern(route.pattern, request.path)
        return Some(route)
    return None

fn matchPattern
  input: (pattern: string, path: string)
  output: bool
  intent: "패턴과 경로 매칭"
  do
    # /todos/:id와 /todos/123 매칭
    # 구현 예시
    let parts = split(path, "/")
    let patterns = split(pattern, "/")
    if parts.len() != patterns.len()
      return false

    for i in 0..parts.len()
      if patterns[i].startsWith(":")
        continue  # 매개변수, 모두 매칭
      if parts[i] != patterns[i]
        return false
    return true
```

### 3.2 경로 매개변수 추출

```freelang
fn extractParams
  input: (pattern: string, path: string)
  output: map<string>
  intent: "URL 패턴에서 매개변수 추출"
  do
    let result = {}
    let parts = split(path, "/")
    let patterns = split(pattern, "/")

    for i in 0..parts.len()
      if patterns[i].startsWith(":")
        let paramName = patterns[i].substring(1)  # ":" 제거
        result[paramName] = parts[i]

    return result

# 사용 예시:
# extractParams("/todos/:id", "/todos/123")
# → { "id": "123" }
```

---

## 4. 요청/응답 설계

### 4.1 표준 JSON 응답 포맷

```freelang
# 성공 응답 (200, 201)
struct SuccessResponse<T>
  success: bool       # true
  data: T            # 실제 데이터
  meta: {
    timestamp: string # ISO 8601
    version: string   # API 버전
  }

# 예시:
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Buy milk",
    "completed": false
  },
  "meta": {
    "timestamp": "2026-03-02T10:30:00Z",
    "version": "v1"
  }
}

# 오류 응답 (400, 404, 500)
struct ErrorResponse
  success: bool       # false
  error: {
    code: string      # "NOT_FOUND", "INVALID_INPUT"
    message: string   # 사용자 친화적 메시지
    details: string   # 상세 정보
  }
  meta: {
    timestamp: string
    version: string
  }

# 예시:
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Todo not found",
    "details": "Todo with id 999 does not exist"
  },
  "meta": {
    "timestamp": "2026-03-02T10:30:00Z",
    "version": "v1"
  }
}
```

### 4.2 페이지네이션

```freelang
struct PagedResponse<T>
  success: bool
  data: array<T>
  pagination: {
    page: number       # 현재 페이지 (1부터 시작)
    pageSize: number   # 페이지당 항목 수
    total: number      # 전체 항목 수
    totalPages: number # 전체 페이지 수
  }

# 예시:
# GET /todos?page=2&pageSize=10
# {
#   "success": true,
#   "data": [ ... ],
#   "pagination": {
#     "page": 2,
#     "pageSize": 10,
#     "total": 145,
#     "totalPages": 15
#   }
# }
```

### 4.3 필터링과 정렬

```freelang
# 쿼리 매개변수를 파싱하는 구조
struct QueryParams
  filter: map<string>   # { "status": "completed" }
  sort: string          # "title:asc,date:desc"
  page: number
  pageSize: number

fn parseQueryParams
  input: string         # "?status=completed&sort=title:asc"
  output: QueryParams
  intent: "쿼리 매개변수 파싱"
  do
    # 구현: ? 이후를 파싱하여 QueryParams 반환
    result = {
      filter: {},
      sort: "",
      page: 1,
      pageSize: 10
    }
    # 파싱 로직...
    return result
```

---

## 5. 예제: 간단한 TODO API

### 5.1 데이터 모델

```freelang
# Todo 항목의 데이터 구조
struct Todo
  id: number
  title: string
  description: string
  completed: bool
  createdAt: string       # ISO 8601
  updatedAt: string

# 데이터베이스 (인메모리)
struct TodoStore
  todos: map<number, Todo>
  nextId: number

# Todo 저장소 초기화
fn initTodoStore
  input: void
  output: TodoStore
  intent: "빈 Todo 저장소 생성"
  do
    return {
      todos: {},
      nextId: 1
    }
```

### 5.2 GET /todos - 모든 Todo 조회

```freelang
fn handleGetTodos
  input: (store: TodoStore, queryParams: QueryParams)
  output: Result<SuccessResponse<array<Todo>>, Error>
  intent: "모든 Todo 항목을 조회하고 반환"
  do
    # 1. 모든 todo 가져오기
    let allTodos = values(store.todos)  # map의 모든 값

    # 2. 필터링 적용
    let filtered = allTodos
    if queryParams.filter.contains("completed")
      let completedValue = queryParams.filter["completed"] == "true"
      filtered = filter(filtered, fn(todo) => todo.completed == completedValue)

    # 3. 정렬 적용
    if queryParams.sort != ""
      filtered = applySort(filtered, queryParams.sort)

    # 4. 페이지네이션 적용
    let total = filtered.len()
    let start = (queryParams.page - 1) * queryParams.pageSize
    let end = start + queryParams.pageSize
    let paged = slice(filtered, start, end)

    # 5. 응답 생성
    return Ok({
      success: true,
      data: paged,
      meta: {
        timestamp: getCurrentTime(),
        version: "v1"
      }
    })
```

### 5.3 POST /todos - 새 Todo 생성

```freelang
struct CreateTodoRequest
  title: string
  description: string

fn handleCreateTodo
  input: (store: TodoStore, request: CreateTodoRequest)
  output: Result<(TodoStore, SuccessResponse<Todo>), Error>
  intent: "새로운 Todo 항목을 생성"
  do
    # 1. 입력 검증
    if request.title == ""
      return Err(Error::InvalidInput("title cannot be empty"))

    # 2. 새 Todo 생성
    let now = getCurrentTime()
    let newTodo = {
      id: store.nextId,
      title: request.title,
      description: request.description,
      completed: false,
      createdAt: now,
      updatedAt: now
    }

    # 3. 저장소에 추가
    let updatedStore = {
      todos: store.todos,
      nextId: store.nextId + 1
    }
    updatedStore.todos[newTodo.id] = newTodo

    # 4. 응답 생성 (201 Created 상태)
    return Ok((updatedStore, {
      success: true,
      data: newTodo,
      meta: {
        timestamp: now,
        version: "v1"
      }
    }))
```

### 5.4 GET /todos/{id} - 특정 Todo 조회

```freelang
fn handleGetTodoById
  input: (store: TodoStore, id: number)
  output: Result<SuccessResponse<Todo>, Error>
  intent: "ID로 특정 Todo 항목을 조회"
  do
    # 1. Todo 찾기
    match store.todos.get(id)
      Some(todo) => {
        # 2. 응답 생성 (200 OK)
        return Ok({
          success: true,
          data: todo,
          meta: {
            timestamp: getCurrentTime(),
            version: "v1"
          }
        })
      }
      None => {
        # 404 Not Found 오류
        return Err(Error::NotFound("Todo not found with id: " + str(id)))
      }
```

### 5.5 PUT /todos/{id} - Todo 수정

```freelang
struct UpdateTodoRequest
  title: Option<string>
  description: Option<string>
  completed: Option<bool>

fn handleUpdateTodo
  input: (store: TodoStore, id: number, request: UpdateTodoRequest)
  output: Result<(TodoStore, SuccessResponse<Todo>), Error>
  intent: "기존 Todo 항목을 부분 업데이트"
  do
    # 1. Todo 찾기
    match store.todos.get(id)
      None => return Err(Error::NotFound("Todo not found"))
      Some(todo) => {
        # 2. 필드 업데이트
        let updated = todo
        if request.title.isSome()
          updated.title = request.title.unwrap()
        if request.description.isSome()
          updated.description = request.description.unwrap()
        if request.completed.isSome()
          updated.completed = request.completed.unwrap()

        # 3. updatedAt 갱신
        updated.updatedAt = getCurrentTime()

        # 4. 저장소에 저장
        let updatedStore = store
        updatedStore.todos[id] = updated

        # 5. 응답 생성 (200 OK)
        return Ok((updatedStore, {
          success: true,
          data: updated,
          meta: {
            timestamp: updated.updatedAt,
            version: "v1"
          }
        }))
      }
```

### 5.6 DELETE /todos/{id} - Todo 삭제

```freelang
fn handleDeleteTodo
  input: (store: TodoStore, id: number)
  output: Result<TodoStore, Error>
  intent: "Todo 항목을 삭제"
  do
    # 1. Todo 존재 확인
    if !store.todos.contains(id)
      return Err(Error::NotFound("Todo not found with id: " + str(id)))

    # 2. 삭제
    let updatedStore = store
    updatedStore.todos.remove(id)

    # 3. 업데이트된 저장소 반환 (응답 본문 없음)
    return Ok(updatedStore)
```

### 5.7 라우터 통합

```freelang
fn createTodoRouter
  input: void
  output: array<Route>
  intent: "Todo API의 모든 라우트 정의"
  do
    return [
      Route {
        pattern: "/todos",
        method: "GET",
        handler: handleGetTodos
      },
      Route {
        pattern: "/todos",
        method: "POST",
        handler: handleCreateTodo
      },
      Route {
        pattern: "/todos/:id",
        method: "GET",
        handler: handleGetTodoById
      },
      Route {
        pattern: "/todos/:id",
        method: "PUT",
        handler: handleUpdateTodo
      },
      Route {
        pattern: "/todos/:id",
        method: "DELETE",
        handler: handleDeleteTodo
      }
    ]
```

---

## 6. 예제: 블로그 API

### 6.1 데이터 모델

```freelang
struct Post
  id: number
  title: string
  slug: string           # URL-friendly: "my-first-post"
  content: string
  authorId: number
  published: bool
  tags: array<string>
  views: number
  createdAt: string
  updatedAt: string

struct Comment
  id: number
  postId: number
  authorId: number
  content: string
  createdAt: string

struct BlogStore
  posts: map<number, Post>
  comments: map<number, Comment>
  nextPostId: number
  nextCommentId: number
```

### 6.2 블로그 API 핸들러 예제

```freelang
# 게시물 목록 조회 (필터링, 정렬, 페이지네이션)
fn handleGetPosts
  input: (store: BlogStore, queryParams: QueryParams)
  output: Result<SuccessResponse<array<Post>>, Error>
  intent: "블로그 게시물 목록을 조회"
  do
    let allPosts = values(store.posts)

    # 필터: 발행된 게시물만
    let filtered = filter(allPosts, fn(post) => post.published == true)

    # 필터: 특정 태그로 검색
    if queryParams.filter.contains("tag")
      let tag = queryParams.filter["tag"]
      filtered = filter(filtered, fn(post) => post.tags.contains(tag))

    # 정렬: 최신 순
    filtered = sortBy(filtered, fn(post) => post.createdAt, "desc")

    # 페이지네이션
    let paginated = paginate(filtered, queryParams.page, queryParams.pageSize)

    return Ok({
      success: true,
      data: paginated.items,
      meta: {
        timestamp: getCurrentTime(),
        version: "v1"
      }
    })

# 특정 게시물과 댓글 조회
fn handleGetPostDetail
  input: (store: BlogStore, postId: number)
  output: Result<SuccessResponse<(Post, array<Comment>)>, Error>
  intent: "특정 게시물과 그 댓글을 조회"
  do
    match store.posts.get(postId)
      None => Err(Error::NotFound("Post not found"))
      Some(post) => {
        # 조회수 증가
        post.views = post.views + 1

        # 댓글 조회
        let postComments = filter(values(store.comments),
          fn(comment) => comment.postId == postId)

        return Ok({
          success: true,
          data: (post, postComments),
          meta: {
            timestamp: getCurrentTime(),
            version: "v1"
          }
        })
      }

# 새 댓글 추가
fn handleAddComment
  input: (store: BlogStore, postId: number, authorId: number, content: string)
  output: Result<BlogStore, Error>
  intent: "게시물에 댓글을 추가"
  do
    # 게시물 존재 확인
    if !store.posts.contains(postId)
      return Err(Error::NotFound("Post not found"))

    # 댓글 생성
    let newComment = Comment {
      id: store.nextCommentId,
      postId: postId,
      authorId: authorId,
      content: content,
      createdAt: getCurrentTime()
    }

    # 저장소에 추가
    store.comments[newComment.id] = newComment
    store.nextCommentId = store.nextCommentId + 1

    return Ok(store)
```

---

## 7. 에러 처리 전략

### 7.1 표준 에러 타입

```freelang
enum ApiError
  ValidationError(string)    # 입력 검증 실패
  NotFound(string)          # 자원 없음
  Unauthorized(string)      # 인증 실패
  Forbidden(string)         # 권한 없음
  Conflict(string)          # 충돌 (예: 중복)
  ServerError(string)       # 서버 오류
  DatabaseError(string)     # 데이터베이스 오류
  ExternalServiceError(string) # 외부 서비스 오류

fn errorToHttpStatus
  input: ApiError
  output: number
  intent: "API 오류를 HTTP 상태 코드로 변환"
  do
    match error
      ValidationError(_) => 400
      NotFound(_) => 404
      Unauthorized(_) => 401
      Forbidden(_) => 403
      Conflict(_) => 409
      ServerError(_) => 500
      DatabaseError(_) => 503
      ExternalServiceError(_) => 502
```

### 7.2 에러 응답 생성

```freelang
fn createErrorResponse
  input: (status: number, error: ApiError)
  output: ErrorResponse
  intent: "API 에러를 표준 응답 형식으로 변환"
  do
    let (code, message) = match error
      ValidationError(msg) => ("VALIDATION_ERROR", msg)
      NotFound(msg) => ("NOT_FOUND", msg)
      Unauthorized(msg) => ("UNAUTHORIZED", msg)
      Forbidden(msg) => ("FORBIDDEN", msg)
      Conflict(msg) => ("CONFLICT", msg)
      ServerError(msg) => ("SERVER_ERROR", msg)
      DatabaseError(msg) => ("DATABASE_ERROR", msg)
      ExternalServiceError(msg) => ("SERVICE_ERROR", msg)

    return {
      success: false,
      error: {
        code: code,
        message: message,
        details: ""
      },
      meta: {
        timestamp: getCurrentTime(),
        version: "v1"
      }
    }
```

### 7.3 트랜잭션 안전성

```freelang
# 여러 작업을 원자적으로 실행
fn executeTransaction<T>
  input: (store: TodoStore, operation: fn(TodoStore) -> Result<T, Error>)
  output: Result<(TodoStore, T), Error>
  intent: "트랜잭션 형태로 저장소 업데이트"
  do
    # 롤백을 위해 이전 상태 저장
    let backup = cloneStore(store)

    match operation(store)
      Ok(result) => Ok((store, result))
      Err(error) => {
        # 롤백
        return Err(error)
      }
```

---

## 8. 테스트 패턴

### 8.1 단위 테스트

```freelang
fn testCreateTodo
  intent: "새로운 Todo 생성 테스트"
  do
    let store = initTodoStore()
    let request = CreateTodoRequest {
      title: "Test Todo",
      description: "A test"
    }

    match handleCreateTodo(store, request)
      Ok((newStore, response)) => {
        assert(response.success == true)
        assert(response.data.id == 1)
        assert(response.data.title == "Test Todo")
        assert(newStore.nextId == 2)
        print("✓ testCreateTodo passed")
      }
      Err(error) => {
        print("✗ testCreateTodo failed: " + str(error))
      }

fn testValidationError
  intent: "빈 제목으로 Todo 생성 실패 테스트"
  do
    let store = initTodoStore()
    let request = CreateTodoRequest {
      title: "",
      description: ""
    }

    match handleCreateTodo(store, request)
      Ok(_) => {
        print("✗ testValidationError failed: should have failed")
      }
      Err(error) => {
        assert(error == Error::InvalidInput)
        print("✓ testValidationError passed")
      }
```

### 8.2 통합 테스트

```freelang
fn testTodoWorkflow
  intent: "전체 Todo 워크플로우 테스트"
  do
    let store = initTodoStore()

    # 1. Todo 생성
    let createReq = CreateTodoRequest {
      title: "Learn FreeLang",
      description: "Complete REST API guide"
    }
    match handleCreateTodo(store, createReq)
      Ok((store1, resp1)) => {
        assert(resp1.success == true)
        let todoId = resp1.data.id

        # 2. Todo 조회
        match handleGetTodoById(store1, todoId)
          Ok(resp2) => {
            assert(resp2.data.title == "Learn FreeLang")

            # 3. Todo 수정
            let updateReq = UpdateTodoRequest {
              completed: Some(true)
            }
            match handleUpdateTodo(store1, todoId, updateReq)
              Ok((store2, resp3)) => {
                assert(resp3.data.completed == true)
                print("✓ testTodoWorkflow passed")
              }
              Err(e) => print("✗ Update failed")
          }
          Err(e) => print("✗ Get failed")
      }
      Err(e) => print("✗ Create failed")
```

---

## 9. 성능 고려사항

### 9.1 응답 시간 목표

```
GET /todos          → < 50ms   (페이지네이션)
GET /todos/{id}     → < 10ms   (ID 조회)
POST /todos         → < 20ms   (생성)
PUT /todos/{id}     → < 20ms   (수정)
DELETE /todos/{id}  → < 15ms   (삭제)
```

### 9.2 최적화 전략

```freelang
# 인덱싱 (복합 쿼리 최적화)
struct OptimizedStore
  todosById: map<number, Todo>           # ID 조회용
  todosByStatus: map<bool, array<Todo>>  # 상태별 조회용
  todosByCreatedDate: array<Todo>        # 시간순 정렬용

# 캐싱 (자주 접근하는 데이터)
struct CachedTodoStore
  store: TodoStore
  cache: map<string, string>   # 캐시 (key -> JSON)
  cacheExpiry: map<string, number>

# 배치 처리
fn batchCreateTodos
  input: (store: TodoStore, requests: array<CreateTodoRequest>)
  output: Result<TodoStore, Error>
  intent: "여러 Todo를 한 번에 생성"
  do
    let result = store
    for request in requests
      match handleCreateTodo(result, request)
        Ok((newStore, _)) => result = newStore
        Err(e) => return Err(e)
    return Ok(result)
```

---

## 10. 보안 고려사항

### 10.1 입력 검증

```freelang
fn validateTodoInput
  input: CreateTodoRequest
  output: Result<void, Error>
  intent: "Todo 입력 데이터 검증"
  do
    # 필수 필드 확인
    if request.title == ""
      return Err(Error::ValidationError("title is required"))

    if request.title.len() > 255
      return Err(Error::ValidationError("title too long"))

    # SQL Injection 방지 (쿼리에 직접 사용할 경우)
    if containsSuspiciousChars(request.title)
      return Err(Error::ValidationError("invalid characters"))

    return Ok(void)

fn containsSuspiciousChars
  input: string
  output: bool
  intent: "위험한 문자 패턴 감지"
  do
    # ; ' " -- % 등 SQL 키워드 감지
    # 실제 구현은 더 복잡할 수 있음
    let dangerous = ["';", "\"", "--", "%"]
    for char in dangerous
      if string.contains(char)
        return true
    return false
```

### 10.2 속도 제한 (Rate Limiting)

```freelang
struct RateLimiter
  requestCounts: map<string, number>  # IP -> 요청 수
  windowResetTime: map<string, number> # IP -> 리셋 시간

fn checkRateLimit
  input: (limiter: RateLimiter, ip: string, limit: number, windowSeconds: number)
  output: Result<void, Error>
  intent: "IP 기반 속도 제한 확인"
  do
    let now = getCurrentTime()
    let resetTime = limiter.windowResetTime.get(ip).getOrElse(0)

    # 윈도우 리셋 확인
    if now > resetTime
      limiter.requestCounts[ip] = 1
      limiter.windowResetTime[ip] = now + windowSeconds
      return Ok(void)

    # 제한 확인
    let count = limiter.requestCounts[ip].getOrElse(0)
    if count >= limit
      return Err(Error::Forbidden("Rate limit exceeded"))

    limiter.requestCounts[ip] = count + 1
    return Ok(void)
```

---

## 11. 배포 체크리스트

- [ ] 모든 엔드포인트 구현 완료
- [ ] 핸들러 단위 테스트 작성 (100% 커버리지)
- [ ] 통합 테스트 작성 (주요 시나리오)
- [ ] 에러 처리 검증
- [ ] 입력 검증 검증
- [ ] 성능 테스트 (응답 시간)
- [ ] 보안 검토 (SQL Injection, XSS 등)
- [ ] API 문서 작성
- [ ] 로깅 설정
- [ ] 모니터링 설정

---

## 12. 다음 단계 (Phase B+)

### 12.1 인증/인가

```freelang
struct AuthToken
  userId: number
  expiresAt: number
  scope: array<string>  # ["read:todos", "write:todos"]

fn validateToken
  input: (token: AuthToken, scope: string)
  output: Result<number, Error>
  intent: "토큰 검증 및 사용자 ID 반환"
  do
    if token.expiresAt < getCurrentTime()
      return Err(Error::Unauthorized("Token expired"))
    if !token.scope.contains(scope)
      return Err(Error::Forbidden("Insufficient permissions"))
    return Ok(token.userId)
```

### 12.2 데이터베이스 통합

Phase B 런타임 완성 후:

```freelang
# 데이터베이스 쿼리
fn queryTodosFromDB
  input: (db: Database, filter: map<string>, page: number)
  output: Result<array<Todo>, Error>
  intent: "데이터베이스에서 Todo 조회"
  do
    let query = "SELECT * FROM todos"
    # ... 쿼리 빌더 로직
    return db.query(query)
```

### 12.3 메시지 큐 통합

비동기 작업 처리:

```freelang
struct QueueMessage
  type: string
  payload: string
  createdAt: number

fn publishTodoCreatedEvent
  input: (queue: MessageQueue, todo: Todo)
  output: Result<void, Error>
  intent: "Todo 생성 이벤트를 큐에 발행"
  do
    return queue.publish(QueueMessage {
      type: "todo:created",
      payload: serializeToJSON(todo),
      createdAt: getCurrentTime()
    })
```

---

**작성자**: AI Assistant
**마지막 업데이트**: 2026-03-02
**버전**: v1.0
