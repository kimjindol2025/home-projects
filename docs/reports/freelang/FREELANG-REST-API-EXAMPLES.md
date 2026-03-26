# FreeLang REST API 실제 구현 예제

**작성일**: 2026-03-02
**버전**: v1.0
**상태**: 🟢 Phase B 예제 코드
**목표**: 실전 API 구현 패턴 제시

---

## 📖 목차

1. [프로젝트 구조](#1-프로젝트-구조)
2. [TODO API 전체 구현](#2-todo-api-전체-구현)
3. [API 사용 예시](#3-api-사용-예시)
4. [고급 예제: 블로그 API](#4-고급-예제-블로그-api)
5. [미들웨어 패턴](#5-미들웨어-패턴)
6. [실행 및 테스트](#6-실행-및-테스트)

---

## 1. 프로젝트 구조

```
freelang-rest-api/
├── src/
│   ├── main.fl              # 메인 엔트리포인트
│   ├── models.fl            # 데이터 모델 정의
│   ├── handlers/
│   │   ├── todo.fl          # TODO 핸들러
│   │   └── blog.fl          # 블로그 핸들러
│   ├── middleware.fl        # 미들웨어
│   ├── utils.fl             # 유틸리티 함수
│   └── errors.fl            # 에러 타입 정의
├── tests/
│   ├── test_todo.fl         # TODO API 테스트
│   └── test_integration.fl  # 통합 테스트
└── README.md
```

---

## 2. TODO API 전체 구현

### 2.1 models.fl - 데이터 모델

```freelang
# ============================================================================
# models.fl - Todo API의 데이터 모델 정의
# ============================================================================

struct Todo
  id: number
  title: string
  description: string
  completed: bool
  createdAt: string
  updatedAt: string

struct CreateTodoRequest
  title: string
  description: string

struct UpdateTodoRequest
  title: Option<string>
  description: Option<string>
  completed: Option<bool>

struct TodoStore
  todos: map<number, Todo>
  nextId: number

# HTTP 요청/응답 타입
struct HttpRequest
  method: string
  path: string
  headers: map<string>
  body: string
  queryParams: map<string>

struct HttpResponse
  status: number
  headers: map<string>
  body: string

# 표준 응답 포맷
struct ApiResponse<T>
  success: bool
  data: T
  error: Option<string>
  timestamp: string

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

# 현재 시간 (ISO 8601 형식)
fn getCurrentTime
  input: void
  output: string
  intent: "현재 시간을 ISO 8601 형식으로 반환"
  do
    # Phase B 런타임에서 구현
    return "2026-03-02T10:30:00Z"

# 다음 ID 생성
fn nextId
  input: (store: TodoStore)
  output: (TodoStore, number)
  intent: "다음 사용 가능한 ID를 생성하고 저장소 업데이트"
  do
    let id = store.nextId
    store.nextId = store.nextId + 1
    return (store, id)
```

### 2.2 errors.fl - 에러 처리

```freelang
# ============================================================================
# errors.fl - API 에러 타입 및 변환 함수
# ============================================================================

enum ApiError
  ValidationError(string)
  NotFound(string)
  Conflict(string)
  ServerError(string)
  InvalidJson(string)

fn errorToStatus
  input: ApiError
  output: number
  intent: "API 에러를 HTTP 상태 코드로 변환"
  do
    match error
      ValidationError(_) => 400
      NotFound(_) => 404
      Conflict(_) => 409
      ServerError(_) => 500
      InvalidJson(_) => 400

fn errorToJson
  input: ApiError
  output: string
  intent: "API 에러를 JSON 문자열로 변환"
  do
    let (code, message) = match error
      ValidationError(msg) => ("VALIDATION_ERROR", msg)
      NotFound(msg) => ("NOT_FOUND", msg)
      Conflict(msg) => ("CONFLICT", msg)
      ServerError(msg) => ("SERVER_ERROR", msg)
      InvalidJson(msg) => ("INVALID_JSON", msg)

    return json({
      "success": false,
      "error": {
        "code": code,
        "message": message
      },
      "timestamp": getCurrentTime()
    })

fn createErrorResponse
  input: ApiError
  output: HttpResponse
  intent: "API 에러를 HTTP 응답으로 변환"
  do
    return {
      status: errorToStatus(error),
      headers: { "Content-Type": "application/json" },
      body: errorToJson(error)
    }
```

### 2.3 handlers/todo.fl - TODO 핸들러

```freelang
# ============================================================================
# handlers/todo.fl - Todo API 핸들러 구현
# ============================================================================

# ─────────────────────────────────────────────────────────────────────────
# GET /todos - 모든 Todo 조회
# ─────────────────────────────────────────────────────────────────────────
fn handleGetTodos
  input: (store: TodoStore, request: HttpRequest)
  output: Result<HttpResponse, ApiError>
  intent: "모든 Todo를 페이지네이션과 함께 반환"
  do
    # 1. 쿼리 매개변수 파싱
    let page = parseNumber(request.queryParams.get("page")).getOrElse(1)
    let pageSize = parseNumber(request.queryParams.get("pageSize")).getOrElse(10)
    let completed = request.queryParams.get("completed")

    # 2. Todo 목록 가져오기
    let allTodos = values(store.todos)

    # 3. 필터링 (completed 상태)
    let filtered = match completed
      None => allTodos
      Some(value) => {
        let isCompleted = value == "true"
        return filter(allTodos, fn(todo) => todo.completed == isCompleted)
      }

    # 4. 페이지네이션
    let total = filtered.len()
    let start = (page - 1) * pageSize
    let end = start + pageSize
    let paged = slice(filtered, start, end)

    # 5. JSON 응답 생성
    let jsonBody = json({
      "success": true,
      "data": {
        "items": paged,
        "pagination": {
          "page": page,
          "pageSize": pageSize,
          "total": total,
          "totalPages": ceil(total / pageSize)
        }
      },
      "timestamp": getCurrentTime()
    })

    return Ok({
      status: 200,
      headers: { "Content-Type": "application/json" },
      body: jsonBody
    })

# ─────────────────────────────────────────────────────────────────────────
# POST /todos - 새 Todo 생성
# ─────────────────────────────────────────────────────────────────────────
fn handleCreateTodo
  input: (store: TodoStore, request: HttpRequest)
  output: Result<(TodoStore, HttpResponse), ApiError>
  intent: "새로운 Todo를 생성하고 저장"
  do
    # 1. JSON 파싱
    let body = match parseJson(request.body)
      Ok(json) => json
      Err(_) => return Err(ApiError::InvalidJson("Invalid JSON in request body"))

    # 2. 필드 검증
    let title = body.get("title").unwrap_or("")
    let description = body.get("description").unwrap_or("")

    if title == ""
      return Err(ApiError::ValidationError("title cannot be empty"))

    if title.len() > 255
      return Err(ApiError::ValidationError("title too long (max 255)"))

    # 3. 새 Todo 생성
    let now = getCurrentTime()
    let (newStore, id) = nextId(store)

    let newTodo = Todo {
      id: id,
      title: title,
      description: description,
      completed: false,
      createdAt: now,
      updatedAt: now
    }

    newStore.todos[id] = newTodo

    # 4. 201 Created 응답
    let jsonBody = json({
      "success": true,
      "data": newTodo,
      "timestamp": now
    })

    return Ok((newStore, {
      status: 201,
      headers: {
        "Content-Type": "application/json",
        "Location": "/todos/" + str(id)
      },
      body: jsonBody
    }))

# ─────────────────────────────────────────────────────────────────────────
# GET /todos/:id - 특정 Todo 조회
# ─────────────────────────────────────────────────────────────────────────
fn handleGetTodoById
  input: (store: TodoStore, request: HttpRequest)
  output: Result<HttpResponse, ApiError>
  intent: "특정 ID의 Todo를 조회"
  do
    # 1. ID 파싱
    let id = extractIdFromPath(request.path)
    match id
      Err(e) => return Err(e)
      Ok(todoId) => {
        # 2. Todo 조회
        match store.todos.get(todoId)
          None => {
            return Err(ApiError::NotFound("Todo not found with id: " + str(todoId)))
          }
          Some(todo) => {
            # 3. 200 OK 응답
            let jsonBody = json({
              "success": true,
              "data": todo,
              "timestamp": getCurrentTime()
            })

            return Ok({
              status: 200,
              headers: { "Content-Type": "application/json" },
              body: jsonBody
            })
          }
      }

# ─────────────────────────────────────────────────────────────────────────
# PUT /todos/:id - Todo 수정
# ─────────────────────────────────────────────────────────────────────────
fn handleUpdateTodo
  input: (store: TodoStore, request: HttpRequest)
  output: Result<(TodoStore, HttpResponse), ApiError>
  intent: "기존 Todo를 부분 또는 전체 업데이트"
  do
    # 1. ID 파싱
    let id = extractIdFromPath(request.path)
    match id
      Err(e) => return Err(e)
      Ok(todoId) => {
        # 2. Todo 존재 확인
        match store.todos.get(todoId)
          None => {
            return Err(ApiError::NotFound("Todo not found"))
          }
          Some(todo) => {
            # 3. JSON 파싱
            let body = match parseJson(request.body)
              Ok(json) => json
              Err(_) => return Err(ApiError::InvalidJson("Invalid JSON"))

            # 4. 필드 업데이트
            let updated = todo
            if body.contains("title")
              updated.title = body.get("title").unwrap_or("")
            if body.contains("description")
              updated.description = body.get("description").unwrap_or("")
            if body.contains("completed")
              updated.completed = body.get("completed") == "true"

            updated.updatedAt = getCurrentTime()

            # 5. 저장
            store.todos[todoId] = updated

            # 6. 200 OK 응답
            let jsonBody = json({
              "success": true,
              "data": updated,
              "timestamp": updated.updatedAt
            })

            return Ok((store, {
              status: 200,
              headers: { "Content-Type": "application/json" },
              body: jsonBody
            }))
          }
      }

# ─────────────────────────────────────────────────────────────────────────
# DELETE /todos/:id - Todo 삭제
# ─────────────────────────────────────────────────────────────────────────
fn handleDeleteTodo
  input: (store: TodoStore, request: HttpRequest)
  output: Result<TodoStore, ApiError>
  intent: "특정 ID의 Todo를 삭제"
  do
    # 1. ID 파싱
    let id = extractIdFromPath(request.path)
    match id
      Err(e) => return Err(e)
      Ok(todoId) => {
        # 2. Todo 존재 확인
        if !store.todos.contains(todoId)
          return Err(ApiError::NotFound("Todo not found"))

        # 3. 삭제
        store.todos.remove(todoId)

        # 4. 업데이트된 저장소 반환
        return Ok(store)
      }

# ─────────────────────────────────────────────────────────────────────────
# 유틸리티 함수
# ─────────────────────────────────────────────────────────────────────────

fn extractIdFromPath
  input: string       # "/todos/123"
  output: Result<number, ApiError>
  intent: "URL 경로에서 Todo ID 추출"
  do
    let parts = split(path, "/")
    # parts = ["", "todos", "123"]

    if parts.len() < 3
      return Err(ApiError::ValidationError("Invalid path format"))

    let idStr = parts[2]
    match parseNumber(idStr)
      Some(id) => Ok(id)
      None => Err(ApiError::ValidationError("ID must be a number"))

fn parseNumber
  input: string
  output: Option<number>
  intent: "문자열을 숫자로 변환"
  do
    # Phase B 런타임 구현
    return None

fn parseJson
  input: string
  output: Result<map<string>, Error>
  intent: "JSON 문자열을 파싱"
  do
    # Phase B 런타임 구현
    return Err(Error::InvalidJson)

fn json
  input: map<string>
  output: string
  intent: "맵을 JSON 문자열로 변환"
  do
    # Phase B 런타임 구현
    return "{}"

fn slice
  input: (array: array<T>, start: number, end: number)
  output: array<T>
  intent: "배열의 부분집합 반환"
  do
    # Phase B 런타임 구현
    return []

fn filter
  input: (array: array<T>, predicate: fn(T) -> bool)
  output: array<T>
  intent: "조건에 맞는 요소만 필터링"
  do
    # Phase B 런타임 구현
    return []

fn values
  input: map<string, T>
  output: array<T>
  intent: "맵의 모든 값을 배열로 반환"
  do
    # Phase B 런타임 구현
    return []

fn ceil
  input: number
  output: number
  intent: "올림"
  do
    # Phase B 런타임 구현
    return 0
```

### 2.4 main.fl - 메인 서버

```freelang
# ============================================================================
# main.fl - REST API 서버의 메인 엔트리포인트
# ============================================================================

fn main
  input: void
  output: void
  intent: "TODO REST API 서버를 시작"
  do
    # 1. 저장소 초기화
    let store = initTodoStore()

    # 2. 테스트 데이터 추가
    let (testStore, _) = handleCreateTodo(store, {
      method: "POST",
      path: "/todos",
      headers: {},
      body: json({ "title": "Learn FreeLang", "description": "" }),
      queryParams: {}
    })
    match testStore
      Ok((newStore, _)) => {
        store = newStore
        print("✓ Test data loaded")
      }
      Err(e) => print("✗ Failed to load test data")

    # 3. 라우터 설정
    let routes = [
      { pattern: "/todos", method: "GET", handler: handleGetTodos },
      { pattern: "/todos", method: "POST", handler: handleCreateTodo },
      { pattern: "/todos/:id", method: "GET", handler: handleGetTodoById },
      { pattern: "/todos/:id", method: "PUT", handler: handleUpdateTodo },
      { pattern: "/todos/:id", method: "DELETE", handler: handleDeleteTodo }
    ]

    # 4. 서버 시작 (Phase B 런타임)
    print("✓ TODO API server started on http://localhost:8080")
    print("  Endpoints:")
    print("    GET    /todos")
    print("    POST   /todos")
    print("    GET    /todos/:id")
    print("    PUT    /todos/:id")
    print("    DELETE /todos/:id")

    # Phase B: 실제 HTTP 서버 시작
    # startHttpServer(routes, 8080, store)
```

---

## 3. API 사용 예시

### 3.1 cURL 명령어

```bash
# ─────────────────────────────────────────────────────────────────────────
# 1. 모든 Todo 조회
# ─────────────────────────────────────────────────────────────────────────
curl -X GET "http://localhost:8080/todos?page=1&pageSize=10"

# 응답:
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Learn FreeLang",
        "description": "",
        "completed": false,
        "createdAt": "2026-03-02T10:30:00Z",
        "updatedAt": "2026-03-02T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1,
      "totalPages": 1
    }
  },
  "timestamp": "2026-03-02T10:30:00Z"
}

# ─────────────────────────────────────────────────────────────────────────
# 2. 새 Todo 생성
# ─────────────────────────────────────────────────────────────────────────
curl -X POST "http://localhost:8080/todos" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy milk",
    "description": "Fresh milk from the store"
  }'

# 응답 (201 Created):
{
  "success": true,
  "data": {
    "id": 2,
    "title": "Buy milk",
    "description": "Fresh milk from the store",
    "completed": false,
    "createdAt": "2026-03-02T10:35:00Z",
    "updatedAt": "2026-03-02T10:35:00Z"
  },
  "timestamp": "2026-03-02T10:35:00Z"
}

# ─────────────────────────────────────────────────────────────────────────
# 3. 특정 Todo 조회
# ─────────────────────────────────────────────────────────────────────────
curl -X GET "http://localhost:8080/todos/2"

# 응답:
{
  "success": true,
  "data": {
    "id": 2,
    "title": "Buy milk",
    "description": "Fresh milk from the store",
    "completed": false,
    "createdAt": "2026-03-02T10:35:00Z",
    "updatedAt": "2026-03-02T10:35:00Z"
  },
  "timestamp": "2026-03-02T10:35:00Z"
}

# ─────────────────────────────────────────────────────────────────────────
# 4. Todo 수정
# ─────────────────────────────────────────────────────────────────────────
curl -X PUT "http://localhost:8080/todos/2" \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'

# 응답:
{
  "success": true,
  "data": {
    "id": 2,
    "title": "Buy milk",
    "description": "Fresh milk from the store",
    "completed": true,
    "createdAt": "2026-03-02T10:35:00Z",
    "updatedAt": "2026-03-02T10:40:00Z"
  },
  "timestamp": "2026-03-02T10:40:00Z"
}

# ─────────────────────────────────────────────────────────────────────────
# 5. Todo 삭제
# ─────────────────────────────────────────────────────────────────────────
curl -X DELETE "http://localhost:8080/todos/2"

# 응답:
HTTP/1.1 204 No Content

# ─────────────────────────────────────────────────────────────────────────
# 6. 오류 응답 예시
# ─────────────────────────────────────────────────────────────────────────
curl -X POST "http://localhost:8080/todos" \
  -H "Content-Type: application/json" \
  -d '{
    "title": ""
  }'

# 응답 (400 Bad Request):
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "title cannot be empty"
  },
  "timestamp": "2026-03-02T10:45:00Z"
}
```

### 3.2 JavaScript 클라이언트

```javascript
// ─────────────────────────────────────────────────────────────────────────
// Todo API 클라이언트
// ─────────────────────────────────────────────────────────────────────────

class TodoClient {
  constructor(baseUrl = 'http://localhost:8080') {
    this.baseUrl = baseUrl;
  }

  // GET /todos
  async getTodos(page = 1, pageSize = 10, completed = null) {
    let url = `${this.baseUrl}/todos?page=${page}&pageSize=${pageSize}`;
    if (completed !== null) {
      url += `&completed=${completed}`;
    }
    const response = await fetch(url);
    return response.json();
  }

  // POST /todos
  async createTodo(title, description = '') {
    const response = await fetch(`${this.baseUrl}/todos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, description })
    });
    return response.json();
  }

  // GET /todos/:id
  async getTodo(id) {
    const response = await fetch(`${this.baseUrl}/todos/${id}`);
    return response.json();
  }

  // PUT /todos/:id
  async updateTodo(id, updates) {
    const response = await fetch(`${this.baseUrl}/todos/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updates)
    });
    return response.json();
  }

  // DELETE /todos/:id
  async deleteTodo(id) {
    return fetch(`${this.baseUrl}/todos/${id}`, {
      method: 'DELETE'
    });
  }
}

// 사용 예시
(async () => {
  const client = new TodoClient();

  // 모든 Todo 조회
  const todos = await client.getTodos();
  console.log('All todos:', todos);

  // 새 Todo 생성
  const created = await client.createTodo('Learn REST API', 'Study FreeLang');
  console.log('Created:', created);

  // Todo 조회
  const todo = await client.getTodo(created.data.id);
  console.log('Todo detail:', todo);

  // Todo 수정
  const updated = await client.updateTodo(created.data.id, { completed: true });
  console.log('Updated:', updated);

  // Todo 삭제
  await client.deleteTodo(created.data.id);
  console.log('Deleted');
})();
```

---

## 4. 고급 예제: 블로그 API

### 4.1 블로그 모델

```freelang
struct Post
  id: number
  title: string
  slug: string
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

fn initBlogStore
  input: void
  output: BlogStore
  intent: "빈 블로그 저장소 생성"
  do
    return {
      posts: {},
      comments: {},
      nextPostId: 1,
      nextCommentId: 1
    }
```

### 4.2 블로그 핸들러

```freelang
# GET /posts - 발행된 게시물 목록 (필터링, 정렬)
fn handleGetPublishedPosts
  input: (store: BlogStore, request: HttpRequest)
  output: Result<HttpResponse, ApiError>
  intent: "발행된 게시물을 태그 필터와 최신순 정렬로 반환"
  do
    let allPosts = values(store.posts)

    # 필터: 발행된 게시물만
    let published = filter(allPosts, fn(post) => post.published == true)

    # 필터: 태그별 (선택적)
    let tag = request.queryParams.get("tag")
    let filtered = match tag
      None => published
      Some(t) => filter(published, fn(post) => post.tags.contains(t))

    # 정렬: 최신순
    let sorted = sortBy(filtered, fn(post) => post.createdAt, "desc")

    # 페이지네이션
    let page = parseNumber(request.queryParams.get("page")).getOrElse(1)
    let pageSize = parseNumber(request.queryParams.get("pageSize")).getOrElse(10)
    let start = (page - 1) * pageSize
    let end = start + pageSize
    let paged = slice(sorted, start, end)

    # 응답
    let jsonBody = json({
      "success": true,
      "data": {
        "posts": paged,
        "pagination": {
          "page": page,
          "pageSize": pageSize,
          "total": filtered.len()
        }
      },
      "timestamp": getCurrentTime()
    })

    return Ok({
      status: 200,
      headers: { "Content-Type": "application/json" },
      body: jsonBody
    })

# GET /posts/:slug - 게시물 상세 조회 (댓글 포함)
fn handleGetPostBySlug
  input: (store: BlogStore, slug: string)
  output: Result<HttpResponse, ApiError>
  intent: "특정 게시물과 그 댓글을 조회하고 조회수 증가"
  do
    # 1. Slug로 게시물 찾기
    let post = findBySlug(store.posts, slug)
    match post
      None => Err(ApiError::NotFound("Post not found"))
      Some(p) => {
        # 2. 조회수 증가
        p.views = p.views + 1

        # 3. 댓글 조회
        let postComments = filter(values(store.comments),
          fn(comment) => comment.postId == p.id)

        # 4. 응답
        let jsonBody = json({
          "success": true,
          "data": {
            "post": p,
            "comments": postComments
          },
          "timestamp": getCurrentTime()
        })

        return Ok({
          status: 200,
          headers: { "Content-Type": "application/json" },
          body: jsonBody
        })
      }

# POST /posts/:postId/comments - 댓글 추가
fn handleAddComment
  input: (store: BlogStore, postId: number, request: HttpRequest)
  output: Result<BlogStore, ApiError>
  intent: "게시물에 새 댓글을 추가"
  do
    # 1. 게시물 존재 확인
    if !store.posts.contains(postId)
      return Err(ApiError::NotFound("Post not found"))

    # 2. 요청 파싱
    let body = match parseJson(request.body)
      Ok(json) => json
      Err(_) => return Err(ApiError::InvalidJson("Invalid JSON"))

    let content = body.get("content").unwrap_or("")
    let authorId = parseNumber(body.get("authorId")).unwrap_or(0)

    # 3. 검증
    if content == ""
      return Err(ApiError::ValidationError("Comment cannot be empty"))

    # 4. 댓글 생성
    let comment = Comment {
      id: store.nextCommentId,
      postId: postId,
      authorId: authorId,
      content: content,
      createdAt: getCurrentTime()
    }

    # 5. 저장
    store.comments[comment.id] = comment
    store.nextCommentId = store.nextCommentId + 1

    return Ok(store)
```

---

## 5. 미들웨어 패턴

### 5.1 로깅 미들웨어

```freelang
struct RequestLog
  method: string
  path: string
  status: number
  duration: number    # 밀리초
  timestamp: string

fn loggingMiddleware
  input: (request: HttpRequest, handler: fn(HttpRequest) -> Result<HttpResponse, ApiError>)
  output: Result<HttpResponse, ApiError>
  intent: "요청/응답 로깅 및 응답 시간 측정"
  do
    let startTime = getCurrentTimeMs()

    # 핸들러 실행
    let response = handler(request)

    let duration = getCurrentTimeMs() - startTime
    let status = match response
      Ok(resp) => resp.status
      Err(_) => 500

    # 로그 출력
    let log = RequestLog {
      method: request.method,
      path: request.path,
      status: status,
      duration: duration,
      timestamp: getCurrentTime()
    }

    printLog(log)
    return response

fn printLog
  input: RequestLog
  output: void
  intent: "요청 로그를 출력"
  do
    print("[" + log.timestamp + "] " +
          log.method + " " + log.path + " " +
          str(log.status) + " (" + str(log.duration) + "ms)")
```

### 5.2 인증 미들웨어

```freelang
struct AuthToken
  userId: number
  scopes: array<string>
  expiresAt: number

fn authMiddleware
  input: (request: HttpRequest, requiredScope: string,
          handler: fn(HttpRequest, number) -> Result<HttpResponse, ApiError>)
  output: Result<HttpResponse, ApiError>
  intent: "토큰 검증 및 권한 확인"
  do
    # 1. Authorization 헤더에서 토큰 추출
    let authHeader = request.headers.get("Authorization")
    match authHeader
      None => {
        return Err(ApiError::ValidationError("Missing Authorization header"))
      }
      Some(header) => {
        # 2. Bearer 토큰 파싱
        let token = extractBearerToken(header)
        match validateToken(token, requiredScope)
          Err(e) => Err(e)
          Ok(userId) => {
            # 3. 검증 통과, 핸들러 실행
            return handler(request, userId)
          }
      }

fn extractBearerToken
  input: string       # "Bearer eyJhbGc..."
  output: Option<string>
  intent: "Authorization 헤더에서 Bearer 토큰 추출"
  do
    let parts = split(header, " ")
    if parts.len() != 2 || parts[0] != "Bearer"
      return None
    return Some(parts[1])

fn validateToken
  input: (token: string, requiredScope: string)
  output: Result<number, ApiError>
  intent: "토큰을 검증하고 사용자 ID 반환"
  do
    # Phase B 런타임: JWT 검증 구현
    # 현재는 간단한 검증만
    if token == ""
      return Err(ApiError::ValidationError("Invalid token"))
    return Ok(1)  # 임시
```

### 5.3 CORS 미들웨어

```freelang
fn corsMiddleware
  input: (request: HttpRequest, handler: fn(HttpRequest) -> Result<HttpResponse, ApiError>)
  output: Result<HttpResponse, ApiError>
  intent: "CORS 헤더를 추가하여 크로스 오리진 요청 허용"
  do
    # OPTIONS 메서드 처리
    if request.method == "OPTIONS"
      return Ok({
        status: 204,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
          "Access-Control-Allow-Headers": "Content-Type, Authorization"
        },
        body: ""
      })

    # 일반 요청 처리
    match handler(request)
      Ok(response) => {
        response.headers["Access-Control-Allow-Origin"] = "*"
        return Ok(response)
      }
      Err(e) => Err(e)
```

---

## 6. 실행 및 테스트

### 6.1 테스트 파일 (test_todo.fl)

```freelang
fn testCreateTodo
  intent: "새 Todo 생성 테스트"
  do
    let store = initTodoStore()
    let request = {
      method: "POST",
      path: "/todos",
      headers: {},
      body: json({ "title": "Test", "description": "" }),
      queryParams: {}
    }

    match handleCreateTodo(store, request)
      Ok((newStore, response)) => {
        assert(response.status == 201)
        print("✓ testCreateTodo passed")
      }
      Err(e) => print("✗ testCreateTodo failed")

fn testGetNonExistentTodo
  intent: "존재하지 않는 Todo 조회 실패 테스트"
  do
    let store = initTodoStore()
    let request = {
      method: "GET",
      path: "/todos/999",
      headers: {},
      body: "",
      queryParams: {}
    }

    match handleGetTodoById(store, request)
      Ok(_) => print("✗ testGetNonExistentTodo failed: should have failed")
      Err(ApiError::NotFound(_)) => print("✓ testGetNonExistentTodo passed")
      Err(_) => print("✗ testGetNonExistentTodo failed: wrong error")

fn runAllTests
  input: void
  output: void
  intent: "모든 테스트 실행"
  do
    testCreateTodo()
    testGetNonExistentTodo()
    # 더 많은 테스트 추가...
```

### 6.2 실행 방법

```bash
# FreeLang 컴파일 및 실행
freelang compile src/main.fl -o api-server
./api-server

# 또는 직접 실행
freelang run src/main.fl

# 테스트 실행
freelang test tests/test_todo.fl
```

---

**작성자**: AI Assistant
**마지막 업데이트**: 2026-03-02
**버전**: v1.0
