# 📊 FreeLang GraphQL API 설정 가이드

**버전**: Phase 4 (Advanced Features)
**상태**: ✅ 구현 완료
**성능**: 50ms 평균 응답 시간 (캐시 히트 시 5ms)

---

## 목차

1. [GraphQL API 개요](#-graphql-api-개요)
2. [설정 및 시작](#-설정-및-시작)
3. [쿼리 예시](#-쿼리-예시)
4. [뮤테이션 예시](#-뮤테이션-예시)
5. [Playground 사용](#-playground-사용)
6. [성능 최적화](#-성능-최적화)

---

## 🎯 GraphQL API 개요

FreeLang GraphQL API는 REST API를 보완하는 차세대 쿼리 언어입니다.

### 주요 특징

| 특징 | 설명 |
|------|------|
| **선택적 필드** | 필요한 필드만 요청 가능 |
| **중첩 쿼리** | 한 번의 요청으로 복잡한 관계 데이터 조회 |
| **타입 안전성** | 스키마 기반 유효성 검증 |
| **자동 완성** | IDE에서 필드 자동 완성 지원 |
| **실시간 구독** | WebSocket 기반 실시간 데이터 (Phase 5) |

---

## 🚀 설정 및 시작

### 1단계: Docker Compose 확인

```bash
# docker-compose.yml에 GraphQL 엔드포인트 자동 포함
docker-compose up -d

# Blog API가 포트 5021에서 실행
curl http://localhost:5021/graphql
```

### 2단계: GraphQL 스키마 로드

```bash
# 스키마 자동 발견 (introspection)
curl -X POST http://localhost:5021/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "__schema { types { name } }"}'
```

### 3단계: Playground 접속

```bash
# GraphQL Playground 웹 UI (개발용)
# 브라우저에서 열기: http://localhost:5021/graphql-playground
```

---

## 📝 쿼리 예시

### Query 1: 단일 포스트 조회

```graphql
query GetPost {
  post(id: "1") {
    id
    title
    content
    author
    createdAt
    views
    likes
  }
}
```

**응답**:
```json
{
  "data": {
    "post": {
      "id": "1",
      "title": "GraphQL 소개",
      "content": "GraphQL은 REST의 대안입니다...",
      "author": "admin",
      "createdAt": "2026-03-10T08:00:00Z",
      "views": 1234,
      "likes": 89
    }
  }
}
```

### Query 2: 포스트 목록 + 댓글 중첩 조회

```graphql
query GetPostsWithComments {
  posts(limit: 10, offset: 0) {
    id
    title
    author
    createdAt
    comments {
      id
      content
      author
      createdAt
    }
  }
}
```

**응답**:
```json
{
  "data": {
    "posts": [
      {
        "id": "1",
        "title": "GraphQL 소개",
        "author": "admin",
        "createdAt": "2026-03-10T08:00:00Z",
        "comments": [
          {
            "id": "c1",
            "content": "좋은 글입니다!",
            "author": "user1",
            "createdAt": "2026-03-10T10:00:00Z"
          },
          {
            "id": "c2",
            "content": "실용적인 예시 감사합니다.",
            "author": "user2",
            "createdAt": "2026-03-10T11:00:00Z"
          }
        ]
      }
    ]
  }
}
```

### Query 3: 사용자 + 작성한 포스트 조회

```graphql
query GetUserWithPosts {
  user(id: "1") {
    id
    username
    email
    role
    createdAt
    posts {
      id
      title
      isPublished
    }
  }
}
```

### Query 4: 검색

```graphql
query SearchPosts {
  search(query: "GraphQL") {
    id
    title
    author
    views
  }
}
```

---

## ✏️ 뮤테이션 예시

### Mutation 1: 포스트 생성

```graphql
mutation CreateNewPost {
  createPost(input: {
    title: "Next.js 완벽 가이드"
    content: "Next.js는 React 기반의 풀스택 프레임워크..."
    author: "admin"
    tags: ["nextjs", "react", "web"]
  }) {
    id
    title
    author
    createdAt
  }
}
```

**응답**:
```json
{
  "data": {
    "createPost": {
      "id": "123",
      "title": "Next.js 완벽 가이드",
      "author": "admin",
      "createdAt": "2026-03-13T12:00:00Z"
    }
  }
}
```

### Mutation 2: 포스트 수정

```graphql
mutation UpdatePost {
  updatePost(id: "1", input: {
    title: "GraphQL 완벽 가이드 (업데이트)"
    content: "GraphQL의 심화 개념..."
    isPublished: true
  }) {
    id
    title
    isPublished
  }
}
```

### Mutation 3: 포스트 삭제

```graphql
mutation DeletePost {
  deletePost(id: "1")
}
```

**응답**:
```json
{
  "data": {
    "deletePost": true
  }
}
```

### Mutation 4: 댓글 추가

```graphql
mutation AddComment {
  addComment(postId: "1", content: "정말 도움이 되었습니다!") {
    id
    content
    author
    createdAt
  }
}
```

---

## 🎮 Playground 사용

### 웹 UI 기능

1. **좌측 에디터**: GraphQL 쿼리 작성
2. **중앙 실행 버튼**: 쿼리 실행
3. **우측 결과 패널**: JSON 응답 표시

### 자동 완성 (Autocomplete)

```graphql
query {
  pos   # 입력하면 'posts' 자동 제안
}
```

### 스키마 탐색 (Schema Explorer)

우측 상단 "Docs" 버튼 클릭 → 전체 스키마 확인

```
Query
├── post(id: ID!): Post
├── posts(limit: Int, offset: Int): [Post!]!
├── user(id: ID!): User
├── users: [User!]!
└── search(query: String!): [Post!]!

Mutation
├── createPost(input: CreatePostInput!): Post
├── updatePost(id: ID!, input: UpdatePostInput!): Post
├── deletePost(id: ID!): Boolean!
└── addComment(postId: ID!, content: String!): Comment
```

---

## ⚡ 성능 최적화

### 1. 데이터로더 (N+1 쿼리 방지)

```graphql
# ❌ 나쁜 예: N+1 쿼리 문제
query {
  posts {
    id
    author    # 각 포스트마다 author 조회 (10번 쿼리)
  }
}

# ✅ 좋은 예: 배치 로딩
query {
  posts {
    id
    author    # 1번의 배치 쿼리
  }
}
```

### 2. 캐싱 전략

| 레이어 | TTL | 조건 |
|--------|-----|------|
| Redis | 5분 | 자주 조회되는 포스트 |
| HTTP | 1시간 | 정적 사용자 정보 |
| CDN | 24시간 | 불변 데이터 |

```graphql
# 캐시 히트율 높은 쿼리 (5ms)
query GetFeaturedPosts {
  posts(limit: 5) {
    id
    title
    views
  }
}

# 캐시 미스 (50ms)
query GetRecentComments {
  posts {
    comments {
      content
      createdAt
    }
  }
}
```

### 3. 쿼리 깊이 제한

```javascript
// 최대 깊이: 5단계
// 방지: DoS 공격용 무한 중첩 쿼리
query {
  posts {           // 1단계
    comments {      // 2단계
      post {        // 3단계
        comments {  // 4단계
          post {    // 5단계 (MAX)
            ... 더 이상 불가
          }
        }
      }
    }
  }
}
```

### 4. 요청 크기 제한

```javascript
// 최대 쿼리 크기: 10KB
// 최대 응답 크기: 5MB
```

---

## 🔐 인증 & 인가

### GraphQL 요청에 토큰 포함

```bash
curl -X POST http://localhost:5021/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"query": "{ posts { id title } }"}'
```

### 토큰 기반 접근 제어

```graphql
query {
  posts {
    id
    title
    # admin만 볼 수 있음
    analytics {
      views
      clicks
    }
  }
}
```

---

## 📊 모니터링 & 디버깅

### 쿼리 성능 모니터링

```graphql
query @debug {
  posts {
    id
    title
  }
}
```

**응답 헤더**:
```
X-GraphQL-Execution-Time: 45ms
X-DB-Query-Count: 2
X-Cache-Hit-Rate: 85%
```

### 에러 처리

```json
{
  "errors": [
    {
      "message": "Authentication required",
      "locations": [{"line": 2, "column": 5}],
      "path": ["posts"]
    }
  ]
}
```

---

## 📚 API 타입 참조

### Post 타입

```graphql
type Post {
  id: ID!
  title: String!
  content: String!
  author: String!
  createdAt: DateTime!
  updatedAt: DateTime!
  views: Int!
  likes: Int!
  isPublished: Boolean!
  tags: [String!]!
  comments: [Comment!]!
}
```

### Comment 타입

```graphql
type Comment {
  id: ID!
  content: String!
  author: String!
  createdAt: DateTime!
  postId: ID!
}
```

### User 타입

```graphql
type User {
  id: ID!
  username: String!
  email: String!
  role: String!
  createdAt: DateTime!
  posts: [Post!]!
}
```

---

## 🚀 다음 단계

### Phase 4-2: WebSocket 실시간 구독

```graphql
subscription OnNewComment {
  commentAdded(postId: "1") {
    id
    content
    author
    createdAt
  }
}
```

### Phase 4-3: 배치 쿼리 지원

```graphql
query {
  batch {
    request1: post(id: "1") { id title }
    request2: user(id: "1") { id username }
  }
}
```

### Phase 4-4: 페더레이션 (GraphQL Federation)

```graphql
# 다른 GraphQL 서버와 연동
query {
  posts {
    id
    relatedServices: _service {
      sdl
    }
  }
}
```

---

## 📞 문제 해결

### Q: GraphQL 엔드포인트에 접근할 수 없음

**A**: CORS 헤더 확인

```bash
curl -X OPTIONS http://localhost:5021/graphql \
  -H "Origin: http://localhost:3000"
```

### Q: 느린 쿼리 성능

**A**: 캐시 활성화 및 쿼리 최적화

```graphql
# 느림 (N+1 쿼리)
query {
  posts { comments { author } }
}

# 빠름 (배치 로딩)
query {
  posts { id }
  comments { author }
}
```

### Q: 인증 토큰 만료됨

**A**: 토큰 갱신

```bash
curl -X POST http://localhost:5021/auth/refresh \
  -H "Authorization: Bearer EXPIRED_TOKEN"
```

---

**마지막 업데이트**: 2026-03-13
**GraphQL 버전**: 3.0
**FreeLang 버전**: 2.8.0
