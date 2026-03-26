---
title: "REST API 설계 완벽 가이드 - 은행 시스템으로 배우기"
date: 2026-03-26
author: Content Writer
category: Technical
tags:
  - REST API
  - API Design
  - Banking System
  - HTTP Methods
  - Database Design
  - Best Practices
---

# REST API 설계 완벽 가이드
## 은행 시스템으로 배우는 실전 API 설계

**글을 읽으면 얻을 수 있는 것:**
1. REST API의 핵심 원칙 이해하기
2. 은행 시스템의 API 설계 전략 배우기
3. HTTP 메서드와 상태 코드 올바르게 사용하기
4. 실무에서 즉시 적용할 수 있는 패턴들

---

## 배경: API는 왜 중요한가?

현대의 애플리케이션은 **분리된 시스템들의 조합**입니다.

```
모바일 앱 ──┐
           ├─→ REST API (서버) ←─→ 데이터베이스
웹 앱      ├─→
           │
데스크톱 앱 ┘
```

**API가 없다면?**
- 각 클라이언트마다 다른 코드 작성 필요
- 비즈니스 로직 중복
- 유지보수 악몽

**REST API가 있다면?**
- 하나의 인터페이스로 모든 클라이언트 지원 ✅
- 비즈니스 로직 중앙화 ✅
- 확장과 유지보수 용이 ✅

---

## 문제: 나쁜 API 설계의 위험성

### 1. 일관성 없는 엔드포인트

```
❌ Bad API 설계:
POST   /create-user
GET    /get_user
PUT    /updateUser
DELETE /user_delete

👉 각기 다른 규칙 = 사용자 혼란
```

### 2. 잘못된 HTTP 메서드 사용

```
❌ Bad:
GET /user/delete/123              (읽기 작업에 GET 사용)
POST /user/get/123               (데이터 조회에 POST 사용)

✅ Good:
DELETE /users/123                (삭제는 DELETE)
GET /users/123                   (조회는 GET)
```

### 3. 명확하지 않은 응답

```
❌ Bad:
{
  "result": "ok"
  "code": 1
  "data": { ... }
}

사용자: "code 1이 뭐지? 성공인가 실패인가?"

✅ Good:
HTTP 200 OK
{
  "status": "success",
  "data": { ... },
  "timestamp": "2026-03-26T10:30:00Z"
}

명확함!
```

### 4. 보안 고려 부족

```
❌ Bad:
GET /api/user/balance?user_id=123

사용자: "아, 그럼 user_id=124는?"
(권한 없이 다른 사용자 데이터 접근 가능)

✅ Good:
GET /api/users/me/balance
(로그인한 사용자 본인 데이터만)
```

---

## 해결책: RESTful API 설계 원칙

### REST란?

**REST (Representational State Transfer)**
- **Representational**: 자원을 표현
- **State**: 상태를 전송
- **Transfer**: 이동/변경

### 6가지 핵심 원칙

```
1️⃣ 자원 중심 설계 (Resource-Oriented)
   /users, /accounts, /transactions

2️⃣ 표준 HTTP 메서드 사용
   GET (조회), POST (생성), PUT (수정), DELETE (삭제)

3️⃣ 상태 비저장 (Stateless)
   각 요청이 독립적 (세션 기반 아님)

4️⃣ 표준 상태 코드 사용
   200 (성공), 201 (생성), 400 (잘못됨), 404 (없음)

5️⃣ JSON 형식으로 통일
   일관된 데이터 포맷

6️⃣ API 버전 관리
   /api/v1, /api/v2 등
```

---

## 실전: 은행 시스템 API 설계

### 데이터 모델

```
User (사용자)
├─ id: INT
├─ name: STRING
├─ email: STRING
└─ password: STRING (암호화)

Account (계좌)
├─ id: INT
├─ user_id: INT (누구의 계좌?)
├─ account_number: STRING
├─ balance: FLOAT
└─ created_at: TIMESTAMP

Transaction (거래)
├─ id: INT
├─ from_account_id: INT
├─ to_account_id: INT
├─ amount: FLOAT
├─ type: STRING (deposit, withdraw, transfer)
└─ created_at: TIMESTAMP
```

### 초보자라면?

기본적인 CRUD 엔드포인트를 먼저 설계합시다.

#### 1. 사용자 관리 API

**엔드포인트 목록:**

```
POST   /api/v1/auth/register      사용자 회원가입
POST   /api/v1/auth/login         사용자 로그인
GET    /api/v1/users/me           현재 사용자 정보
PUT    /api/v1/users/me           사용자 정보 수정
DELETE /api/v1/users/me           계정 삭제
```

**예제 1: 회원가입**

```
요청:
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "Kim Jindol",
  "email": "kim@freelang.dev",
  "password": "securePassword123"
}

응답 (성공 - 201 Created):
{
  "id": 1,
  "name": "Kim Jindol",
  "email": "kim@freelang.dev",
  "created_at": "2026-03-26T10:30:00Z"
}

응답 (실패 - 400 Bad Request):
{
  "error": "Email already exists",
  "code": "DUPLICATE_EMAIL"
}
```

**예제 2: 로그인**

```
요청:
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "kim@freelang.dev",
  "password": "securePassword123"
}

응답 (성공 - 200 OK):
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 604800,
  "user": {
    "id": 1,
    "name": "Kim Jindol",
    "email": "kim@freelang.dev"
  }
}

응답 (실패 - 401 Unauthorized):
{
  "error": "Invalid credentials",
  "code": "INVALID_PASSWORD"
}
```

#### 2. 계좌 관리 API

```
POST   /api/v1/accounts           새 계좌 생성
GET    /api/v1/accounts           내 계좌 목록
GET    /api/v1/accounts/:id       특정 계좌 조회
PUT    /api/v1/accounts/:id       계좌 정보 수정
DELETE /api/v1/accounts/:id       계좌 삭제
```

**예제: 계좌 생성**

```
요청:
POST /api/v1/accounts
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_type": "savings"  // savings, checking
}

응답 (201 Created):
{
  "id": 101,
  "user_id": 1,
  "account_number": "1234567890",
  "balance": 0.0,
  "account_type": "savings",
  "created_at": "2026-03-26T10:30:00Z"
}
```

#### 3. 거래 API

```
POST   /api/v1/transactions/deposit    입금
POST   /api/v1/transactions/withdraw   출금
POST   /api/v1/transactions/transfer   송금
GET    /api/v1/transactions            거래 내역 조회
GET    /api/v1/transactions/:id        특정 거래 조회
```

**예제: 입금**

```
요청:
POST /api/v1/transactions/deposit
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_id": 101,
  "amount": 50000
}

응답 (201 Created):
{
  "id": 1001,
  "type": "deposit",
  "account_id": 101,
  "amount": 50000,
  "previous_balance": 0,
  "new_balance": 50000,
  "created_at": "2026-03-26T10:31:00Z"
}

응답 (실패 - 400 Bad Request):
{
  "error": "Invalid amount",
  "code": "INVALID_AMOUNT",
  "details": "Amount must be positive"
}
```

**예제: 송금**

```
요청:
POST /api/v1/transactions/transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "from_account_id": 101,
  "to_account_id": 102,
  "amount": 10000
}

응답 (201 Created):
{
  "id": 1002,
  "type": "transfer",
  "from_account_id": 101,
  "to_account_id": 102,
  "amount": 10000,
  "status": "completed",
  "created_at": "2026-03-26T10:32:00Z"
}

응답 (실패 - 402 Payment Required):
{
  "error": "Insufficient balance",
  "code": "INSUFFICIENT_FUNDS",
  "current_balance": 30000,
  "required": 40000
}
```

### 전문가라면?

복잡한 비즈니스 로직과 성능 최적화를 다루봅시다.

#### 고급 기능: 페이지네이션

```
요청:
GET /api/v1/accounts/101/transactions?page=2&limit=20&sort=created_at:desc

응답 (200 OK):
{
  "data": [
    { "id": 1050, "type": "transfer", "amount": 5000, ... },
    { "id": 1049, "type": "withdraw", "amount": 2000, ... },
    ...
  ],
  "pagination": {
    "page": 2,
    "limit": 20,
    "total": 150,
    "pages": 8
  }
}
```

#### 고급 기능: 필터링

```
요청:
GET /api/v1/transactions?
  from_date=2026-03-01&
  to_date=2026-03-31&
  type=transfer&
  min_amount=1000

응답:
{
  "data": [
    { "id": 1001, "type": "transfer", "amount": 5000, ... },
    { "id": 1005, "type": "transfer", "amount": 2000, ... }
  ],
  "count": 2
}
```

#### 고급 기능: 트랜잭션 무결성

```
코드: 송금이 원자적(Atomic)으로 처리됨

TRANSACTION START {
  1. from_account의 잔액 확인
  2. 잔액 부족하면 ROLLBACK
  3. from_account에서 금액 차감
  4. to_account에 금액 추가
  5. 거래 기록 저장
} COMMIT

중간에 오류 발생 → 모든 변경사항 취소
```

---

## HTTP 상태 코드 올바르게 사용하기

### 성공 코드 (2xx)

```
200 OK              - 요청 성공, 데이터 반환
201 Created         - 리소스 생성 성공
202 Accepted        - 요청 수락 (처리 중)
204 No Content      - 성공하지만 반환 데이터 없음
```

**은행 예제:**
```
POST /api/v1/accounts → 201 Created (계좌 생성됨)
GET /api/v1/accounts → 200 OK (계좌 목록 반환)
DELETE /api/v1/accounts/101 → 204 No Content (삭제 완료)
```

### 클라이언트 오류 (4xx)

```
400 Bad Request     - 요청이 잘못됨
401 Unauthorized    - 인증 필요
403 Forbidden       - 권한 없음
404 Not Found       - 리소스 없음
409 Conflict        - 상태 충돌
```

**은행 예제:**
```
POST /api/v1/auth/login (잘못된 비번) → 401 Unauthorized
GET /api/v1/accounts/999 (없는 계좌) → 404 Not Found
POST /api/v1/transactions/transfer (잔액 부족) → 409 Conflict
```

### 서버 오류 (5xx)

```
500 Internal Server Error  - 예상 밖의 오류
503 Service Unavailable    - 서비스 이용 불가 (DB 다운 등)
```

---

## 보안: API 인증과 권한

### JWT 토큰 기반 인증

```
1️⃣ 로그인
  POST /auth/login
  → 토큰 발급: "eyJhbGciOiJIUzI1NiIs..."

2️⃣ API 호출 시 토큰 첨부
  GET /api/v1/accounts
  Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

3️⃣ 서버에서 토큰 검증
  ✅ 유효한 토큰 → 요청 처리
  ❌ 만료된 토큰 → 401 Unauthorized
  ❌ 잘못된 토큰 → 401 Unauthorized
```

### 권한 검증

```
사용자 A가 사용자 B의 계좌를 조회 시도:
GET /api/v1/accounts/102 (사용자 B의 계좌)
Authorization: Bearer <사용자 A의 토큰>

응답:
403 Forbidden
{
  "error": "You don't have permission to access this account"
}
```

---

## API 문서화

### OpenAPI/Swagger 형식

```yaml
openapi: 3.0.0
info:
  title: Banking API
  version: 1.0.0

paths:
  /api/v1/auth/login:
    post:
      summary: 사용자 로그인
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: 로그인 성공
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                  expires_in:
                    type: integer
```

---

## Best Practices

### ✅ 하는 것

```
1️⃣ 자원 중심 명명
   GET    /api/v1/users
   POST   /api/v1/accounts
   DELETE /api/v1/transactions/123

2️⃣ 버전 관리
   /api/v1/        (현재)
   /api/v2/        (미래)

3️⃣ 일관된 응답 형식
   {
     "status": "success",
     "data": { ... },
     "timestamp": "2026-03-26T10:30:00Z"
   }

4️⃣ 에러 메시지 상세함
   {
     "error": "Insufficient balance",
     "code": "INSUFFICIENT_FUNDS",
     "current_balance": 1000,
     "required": 5000
   }

5️⃣ 타임스탬프 포함
   created_at, updated_at, timestamp
```

### ❌ 하지 말 것

```
❌ 동사를 URL에 포함
   /api/users/get/123
   /api/accounts/delete/456

✅ 대신:
   GET    /api/users/123
   DELETE /api/accounts/456

❌ 깊은 경로
   /api/users/123/accounts/456/transactions/789

✅ 대신:
   /api/users/123       (사용자 조회)
   /api/accounts/456    (계좌 조회)
   GET /api/accounts/456/transactions (거래 목록)

❌ 불명확한 상태 코드
   모든 응답에 200만 사용

✅ 대신:
   201 Created (리소스 생성)
   204 No Content (삭제)
   400 Bad Request (잘못된 입력)
   401 Unauthorized (인증 필요)
```

---

## 성능 최적화

### 1. 캐싱

```
GET /api/v1/accounts/101

응답 헤더:
Cache-Control: max-age=3600    (1시간 캐싱)
ETag: "abc123"

클라이언트는 1시간 동안 캐시된 데이터 사용
→ 서버 부하 감소 ⚡
```

### 2. 압축

```
요청:
GET /api/v1/transactions?limit=1000

응답 헤더:
Content-Encoding: gzip

크기: 1MB → 100KB (10배 압축!)
```

### 3. 배치 요청

```
❌ Bad: 여러 번 요청
GET /api/v1/accounts/101
GET /api/v1/accounts/102
GET /api/v1/accounts/103

✅ Good: 한 번에 요청
POST /api/v1/batch
{
  "requests": [
    { "method": "GET", "path": "/accounts/101" },
    { "method": "GET", "path": "/accounts/102" },
    { "method": "GET", "path": "/accounts/103" }
  ]
}
```

---

## 다음 단계

REST API 설계를 완벽히 이해했으니, 더 고급 주제를 배울 수 있습니다.

**다음에 읽을 거리:**
1. [GraphQL - REST의 대안](#)
2. [API 속도 최적화 - 성능 벤치마크](#)
3. [마이크로서비스 API 설계](#)

**직접 해보고 싶으신가요?**

```bash
# FreeLang 은행 API 실행
git clone https://github.com/freelang/freelang
cd freelang/examples/banking-api
freelang main.fl

# 테스트
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Kim","email":"kim@freelang.dev","password":"123"}'
```

---

**참고 자료:**
- [REST API Best Practices - Microsoft](https://docs.microsoft.com/en-us/azure/architecture/best-practices/api-design)
- [JSON:API - A specification for building APIs](https://jsonapi.org/)
- [OpenAPI Specification](https://spec.openapis.org/)
- FreeLang Banking API Docs: https://freelang-docs.example.com/banking-api

---

이 글이 도움이 되었다면? 👍

- 다른 개발자와 공유해주세요
- "이런 API 패턴을 배우고 싶어요" 댓글로 제안해주세요
- 자신이 설계한 API에 대한 피드백을 공유해주세요!

**Happy API Designing! 🚀**
