# REST API 설계 완벽 가이드: 실전 팁들

## 요약
**"좋은 API란 뭘까요?"** 빠른 건 기본이고, 사용하기 쉬워야 하고, 이해하기 쉬워야 합니다. 우리가 FreeLang 백엔드 프로젝트에서 배운 REST API 설계 원칙들을 공개합니다. 50K req/sec를 처리하는 실전 API를 어떻게 만드는지 살펴봅시다.

---

## REST는 뭘까요?

**REST** = REpresentational State Transfer

쉽게 말해, HTTP 메서드(GET, POST, PUT, DELETE)로 **자원(Resource)**을 다루는 방식입니다.

```
자원: User, Post, Comment 등
메서드:
  GET    - 조회
  POST   - 생성
  PUT    - 수정
  DELETE - 삭제
```

### 좋은 API의 특징
1. **직관적**: URL만 봐도 뭔지 알 수 있음
2. **일관성**: 모든 엔드포인트가 같은 패턴
3. **버전 관리**: 변경사항을 관리할 수 있음
4. **에러 처리**: 명확한 에러 메시지

## 설계 원칙

### 1. 명확한 URL 구조

```
좋은 예:
GET /api/v1/users               - 모든 사용자 조회
GET /api/v1/users/123           - 특정 사용자 조회
POST /api/v1/users              - 새 사용자 생성
PUT /api/v1/users/123           - 사용자 수정
DELETE /api/v1/users/123        - 사용자 삭제

GET /api/v1/users/123/posts     - 사용자의 게시물
POST /api/v1/users/123/posts    - 사용자에게 게시물 추가
```

**피해야 할 것들:**
```
❌ GET /getUser?id=123
❌ GET /user/getAll
❌ POST /users/create
❌ GET /deleteUser?id=123 (GET으로 삭제 금지!)
```

### 2. 적절한 HTTP 상태 코드

```
2xx - 성공
  200 OK              - 요청 성공
  201 Created         - 자원 생성 성공
  204 No Content      - 요청 성공, 응답 본문 없음

4xx - 클라이언트 오류
  400 Bad Request     - 잘못된 요청
  401 Unauthorized    - 인증 필요
  403 Forbidden       - 권한 없음
  404 Not Found       - 자원 없음

5xx - 서버 오류
  500 Internal Server Error     - 서버 오류
  503 Service Unavailable       - 서비스 이용 불가
```

**올바른 사용:**
```go
// 사용자 조회
GET /api/v1/users/123
→ 200 OK + 사용자 정보

// 잘못된 ID
GET /api/v1/users/999
→ 404 Not Found

// 미인증
GET /api/v1/users (인증 필요)
→ 401 Unauthorized

// 생성
POST /api/v1/users
→ 201 Created + 새 사용자 정보
```

### 3. 일관된 응답 형식

모든 응답을 같은 구조로:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  "error": null
}
```

에러 응답도 같은 형식:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "사용자를 찾을 수 없습니다",
    "details": {
      "user_id": 999
    }
  }
}
```

## 실전 예제: User API

### 1. 사용자 조회

```go
// 요청
GET /api/v1/users/123
Authorization: Bearer {token}

// 응답 (200 OK)
{
  "success": true,
  "data": {
    "id": 123,
    "name": "Bob",
    "email": "bob@example.com",
    "created_at": "2026-01-15T10:30:00Z"
  }
}
```

### 2. 사용자 목록 (페이징)

```go
// 요청
GET /api/v1/users?page=1&limit=20&sort=created_at&order=desc
Authorization: Bearer {token}

// 응답 (200 OK)
{
  "success": true,
  "data": {
    "items": [
      { "id": 1, "name": "Alice", ... },
      { "id": 2, "name": "Bob", ... }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1050,
      "pages": 53
    }
  }
}
```

### 3. 새 사용자 생성

```go
// 요청
POST /api/v1/users
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Charlie",
  "email": "charlie@example.com",
  "password": "secure_password"
}

// 응답 (201 Created)
{
  "success": true,
  "data": {
    "id": 124,
    "name": "Charlie",
    "email": "charlie@example.com"
  }
}
```

### 4. 사용자 수정

```go
// 요청
PUT /api/v1/users/123
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "Bob Smith"
}

// 응답 (200 OK)
{
  "success": true,
  "data": {
    "id": 123,
    "name": "Bob Smith",
    "email": "bob@example.com"
  }
}
```

### 5. 사용자 삭제

```go
// 요청
DELETE /api/v1/users/123
Authorization: Bearer {token}

// 응답 (204 No Content)
// 응답 본문 없음
```

## 고급 기능

### 1. 필터링 & 검색

```
GET /api/v1/users?role=admin&status=active
GET /api/v1/users/search?q=alice
GET /api/v1/posts?author_id=123&tags=golang,database
```

### 2. 정렬

```
GET /api/v1/users?sort=created_at&order=desc
GET /api/v1/posts?sort=-created_at  (음수 = 역순)
```

### 3. 부분 응답 (필요한 필드만)

```
GET /api/v1/users/123?fields=id,name,email
→ {
    "id": 123,
    "name": "Bob",
    "email": "bob@example.com"
  }
```

### 4. 포함/확장 (관련 자원)

```
GET /api/v1/users/123?include=posts,followers
→ {
    "id": 123,
    "name": "Bob",
    "posts": [ ... ],
    "followers": [ ... ]
  }
```

## 성능 최적화

### 1. 캐싱 헤더

```
GET /api/v1/users/123

응답 헤더:
Cache-Control: public, max-age=3600
ETag: "abc123"
Last-Modified: 2026-01-15T10:30:00Z
```

클라이언트는 1시간 동안 캐시를 사용할 수 있습니다.

### 2. 압축

```go
// 요청에 gzip 지원을 알려줌
GET /api/v1/users
Accept-Encoding: gzip, deflate

// 응답은 압축됨
Response-Encoding: gzip
```

### 3. 페이징 (대량 데이터)

```
❌ 나쁜 예: 10,000개 항목을 한 번에
GET /api/v1/posts → 10,000개 모두

✅ 좋은 예: 페이징으로 나눔
GET /api/v1/posts?page=1&limit=20 → 20개
GET /api/v1/posts?page=2&limit=20 → 다음 20개
```

## 인증 & 권한

### 1. Bearer Token (JWT)

```
GET /api/v1/users/123
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

백엔드에서 토큰 검증:
- 서명 확인
- 만료 시간 확인
- 권한 확인
```

### 2. 권한 확인

```
GET /api/v1/users/123/private-data
(자신의 데이터만 조회 가능)

토큰에 user_id: 456이면:
- user_id = 123 → 404 (다른 사용자)
- user_id = 456 → 200 OK (자신)
```

## 에러 처리 (실제 예제)

```go
// 상황 1: 잘못된 요청
POST /api/v1/users
{ "email": "invalid-email" }

↓
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "입력값이 유효하지 않습니다",
    "details": {
      "email": "유효한 이메일이 아닙니다"
    }
  }
}

// 상황 2: 자원 없음
GET /api/v1/users/999

↓
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "사용자를 찾을 수 없습니다"
  }
}

// 상황 3: 권한 없음
DELETE /api/v1/users/456  (다른 사용자)

↓
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "이 작업을 수행할 권한이 없습니다"
  }
}
```

## 버전 관리

API는 계속 변합니다. 하지만 기존 사용자를 망가뜨리면 안 됩니다:

```
v1: GET /api/v1/users → 100 req/sec
v2: GET /api/v2/users → 더 빠른 응답

클라이언트가 v1을 계속 쓸 수 있음
```

**언제 새 버전을 만들까요?**
- 응답 형식 변경
- 새 필수 필드
- 엔드포인트 제거

## 성능 벤치마크

FreeLang 백엔드의 실제 성능:

| 엔드포인트 | 처리량 | 평균 응답시간 |
|-----------|--------|-------------|
| GET /users | 50K req/sec | 20ms |
| GET /users/{id} | 45K req/sec | 22ms |
| POST /users | 30K req/sec | 33ms |
| PUT /users/{id} | 28K req/sec | 35ms |
| DELETE /users/{id} | 42K req/sec | 23ms |

**비결:**
- 캐싱 (Redis)
- 연결 풀링 (500 커넥션)
- 쿼리 최적화

## 체크리스트

좋은 API 설계 확인사항:

```
☐ URL이 직관적인가?
☐ HTTP 메서드를 올바르게 사용했는가?
☐ 상태 코드가 적절한가?
☐ 응답 형식이 일관되는가?
☐ 에러 메시지가 명확한가?
☐ 인증/권한이 구현되었는가?
☐ 버전 관리 계획이 있는가?
☐ 성능이 충분한가?
☐ 문서가 충분한가?
```

## 마치며

좋은 REST API는:
1. **직관적** - URL만 봐도 뭔지 알 수 있음
2. **일관적** - 모든 엔드포인트가 같은 패턴
3. **빠른** - 캐싱, 압축, 최적화
4. **안전한** - 인증, 권한, 검증

이 원칙들을 따르면 누구나 사용하기 좋은 API를 만들 수 있습니다!

**궁금한 점:** "내 API도 50K req/sec를 처리할 수 있을까?" 라는 질문이 있으신가요? 댓글로 물어봐주세요. GraphQL, gRPC 같은 다른 프로토콜도 설명해드릴 수 있습니다!
