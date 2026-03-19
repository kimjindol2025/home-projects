---
name: FV 2.0 Phase 1 - FreeLang 현재 상태 분석
description: FreeLang (Sovereign Workspace)의 아키텍처, 표준 라이브러리, 컴파일 파이프라인 분석
type: project
---

# FV 2.0 Phase 1: FreeLang 현재 상태 분석

**작성일**: 2026-03-19
**프로젝트**: FV 2.0 (V Language + FreeLang Integration)
**상태**: 🟡 분석 중

---

## FreeLang 개요

### 정의
FreeLang = 프로덕션 준비 완료된 백엔드 언어
- **개발언어**: FV-Lang (자체호스팅)
- **규모**: 20,370줄 핵심 + 393개 테스트 + 130개 프로젝트 생태계
- **상태**: Phase 1-11 완료, 프로덕션 배포 가능
- **저장소**: GOGS (gogs.dclub.kr/kim/sovereign-workspace)

### 핵심 특징

#### 1. 프로덕션 준비
- ✅ HTTP engine 내장
- ✅ WebSocket 지원 (RFC 6455)
- ✅ SQLite ORM 내장
- ✅ JWT/OAuth2 지원
- ✅ Redis 드라이버
- ✅ gRPC 지원
- ✅ Docker 배포 가능
- ✅ Kubernetes 지원
- ✅ 멀티테넌트 지원
- ✅ 성능 최적화됨

---

## FreeLang 아키텍처

### 컴파일 파이프라인

```
FreeLang 소스 (.fl)
    ↓ [Lexer]
토큰 스트림
    ↓ [Parser]
AST (Abstract Syntax Tree)
    ↓ [Type Checker]
타입 검증된 AST
    ↓ [Optimizer]
최적화된 IR
    ↓ [Code Generator]
C 코드 또는 바이너리
    ↓ [Runtime / JIT]
실행
```

### 주요 컴포넌트

#### Lexer (토큰화)
- **역할**: FreeLang 소스를 토큰으로 변환
- **토큰 수**: 60+ 타입
- **특징**: 명령형 + 함수형 하이브리드 문법 지원

#### Parser (파싱)
- **역할**: 토큰 → AST 변환
- **방식**: 재귀 하강 파싱
- **특징**:
  - 함수, 구조체, 트레이트 정의
  - 라우팅 문법 (`fn route GET /path`)
  - 에러 처리 (? 연산자, or 블록)

#### Type Checker (타입 검사)
- **역할**: 타입 추론 및 검증
- **특징**:
  - Hindley-Milner 타입 추론
  - 제네릭 타입 지원
  - 트레이트 바운드

#### Code Generator (코드 생성)
- **역할**: AST → C 코드 생성
- **출력**: 표준 C 코드 (gcc/clang으로 컴파일 가능)
- **특징**: 결정론적 생성 (같은 입력 = 같은 출력)

#### Runtime (런타임)
- **메모리 관리**: GC (Mark & Sweep)
- **실행 모드**: 인터프리터 또는 JIT 컴파일
- **성능**: 중간 크기 프로젝트 (<10K 줄) 최적화

---

## FreeLang 표준 라이브러리 (20,370줄)

### 모듈 구조

```
stdlib/
├── http/          - HTTP 엔진 (Web Framework)
│   ├── server.fl  - 서버 구현
│   ├── router.fl  - 라우팅
│   ├── middleware - 미들웨어
│   └── response   - 응답 처리
│
├── database/      - 데이터베이스 ORM
│   ├── sqlite.fl  - SQLite
│   ├── query.fl   - 쿼리 빌더
│   ├── migration  - 마이그레이션
│   └── transaction- 트랜잭션
│
├── websocket/     - 실시간 통신
│   ├── server.fl  - WebSocket 서버
│   ├── client.fl  - 클라이언트
│   └── protocol   - RFC 6455
│
├── grpc/          - 마이크로서비스
│   ├── server.fl  - gRPC 서버
│   ├── codegen    - Protocol Buffers
│   └── reflection - 서비스 검사
│
├── auth/          - 인증 & 권한
│   ├── jwt.fl     - JWT
│   ├── oauth2.fl  - OAuth2
│   └── session    - 세션 관리
│
├── crypto/        - 암호화
│   ├── aes.fl     - AES
│   ├── rsa.fl     - RSA
│   ├── sha.fl     - SHA
│   └── hash       - 해시
│
├── json/          - JSON 처리
│   ├── encode.fl  - 인코딩
│   └── decode.fl  - 디코딩
│
├── redis/         - 캐싱
│   ├── client.fl  - Redis 클라이언트
│   ├── pub_sub    - 발행/구독
│   └── cache      - 캐시 유틸
│
├── fs/            - 파일시스템
│   ├── file.fl    - 파일 I/O
│   ├── dir.fl     - 디렉토리
│   └── path       - 경로 처리
│
└── sys/           - 시스템
    ├── os.fl      - 운영체제
    ├── env.fl     - 환경 변수
    └── process    - 프로세스
```

---

## FreeLang 언어 기능

### 기본 문법

#### 변수와 상수
```fv
let x = 10          // 불변
let mut y = 20      // 가변
const PI = 3.14159  // 상수
```

#### 함수
```fv
fn add(a: i64, b: i64) -> i64 {
  a + b
}

// 마지막 식이 반환값
fn multiply(a: i64, b: i64) -> i64 {
  a * b
}
```

#### 구조체
```fv
type User = {
  name: String,
  age: i64,
  mut email: String,
}

impl User {
  fn greet(self) -> String {
    string_concat("Hello, ", self.name)
  }
}
```

#### 에러 처리
```fv
// Option (NULL 안전)
fn divide(a: i64, b: i64) -> Option(i64) {
  if b == 0 {
    None
  } else {
    Some(a / b)
  }
}

// Result (에러 메시지)
fn parse(s: String) -> Result(i64, String) {
  // ... 파싱
  if error {
    Error("Invalid number")
  } else {
    Ok(42)
  }
}

// 에러 전파 (? 연산자)
value := divide(10, 2)?
```

### HTTP 라우팅 문법

#### 기본 라우팅
```fv
fn route GET /api/users -> JSON {
  users := db.all_users()?
  json(users)
}

fn route GET /api/users/:id -> JSON {
  user := db.find(id)?
  json(user)
}

fn route POST /api/users -> JSON {
  user := User {
    name: request.body.name,
    email: request.body.email,
  }
  db.save(user)?
  json(user)
}
```

#### 미들웨어
```fv
fn middleware auth_required(req: Request) -> Result(Request, String) {
  token := req.header("Authorization")?
  if is_valid_token(token) {
    Ok(req)
  } else {
    Error("Unauthorized")
  }
}

// 라우트에 적용
fn route GET /admin @auth_required -> String {
  "Admin only"
}
```

### 데이터베이스 ORM

#### 모델 정의
```fv
type User = {
  id: i64,
  name: String,
  email: String,
  created_at: DateTime,
}

// 자동으로 users 테이블 매핑
```

#### 쿼리
```fv
// 모두 조회
users := db.all::<User>()?

// 하나 조회
user := db.find::<User>(id)?

// 조건 조회
active_users := db.where::<User>("active = true")?

// 저장
db.save(user)?

// 수정
db.update(user)?

// 삭제
db.delete::<User>(id)?
```

#### 마이그레이션
```fv
fn migration_001_create_users() {
  db.exec("""
    CREATE TABLE users (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      email TEXT UNIQUE,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  """)?
}
```

---

## FreeLang vs V vs 기타 언어

### 비교표

| 항목 | V | FreeLang | FV 2.0 (목표) | Python | Go |
|------|---|----------|--------------|--------|-----|
| **문법** | 간단 | 중간 | 간단 | 느슨함 | 간단 |
| **컴파일** | <1초 | 1-5초 | <1초 | 해석 | 빠름 |
| **HTTP 라이브러리** | 미성숙 | ✅ | ✅ | Django/FastAPI | gin |
| **Database ORM** | ❌ | ✅ | ✅ | SQLAlchemy | gorm |
| **WebSocket** | ❌ | ✅ | ✅ | ❌ | ❌ |
| **gRPC** | ❌ | ✅ | ✅ | ❌ | protobuf |
| **JWT/OAuth2** | ❌ | ✅ | ✅ | ❌ | ❌ |
| **Docker** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **K8s** | ❌ | ✅ | ✅ | ✅ | ✅ |
| **성숙도** | 초기 | 프로덕션 | 계획 | 매우 높음 | 높음 |

---

## FreeLang → FV 2.0 마이그레이션 전략

### 무엇을 유지할 것인가

#### 1. IR (Intermediate Representation)
- FreeLang의 IR은 V-호환 AST로부터도 생성 가능
- **이득**: 기존 최적화, JIT 컴파일, 코드 생성 재사용

#### 2. 표준 라이브러리
- HTTP, Database, WebSocket, gRPC 모두 유지
- **이득**: 0에서 시작하는 대신 기존 구현 재사용

#### 3. 타입 시스템
- FreeLang의 타입 검사 로직 재사용
- V 타입 → FreeLang 타입 매핑만 추가

#### 4. 런타임
- GC, 메모리 관리, 인터프리터 재사용

### 무엇을 변경할 것인가

#### 1. Lexer (토큰화)
- **변경**: V 키워드 및 토큰 추가
  - V: `fn`, `mut`, `:=`, `?`, `or`
  - FreeLang: 같은 토큰 (호환)
- **노력**: 낮음 (기존 토큰 90%+ 재사용)

#### 2. Parser (파싱)
- **변경**: V 문법 규칙 추가
  - V: `fn name(params) ReturnType { body }`
  - FreeLang: 비슷하지만 약간 다름
- **노력**: 중간 (새 파싱 규칙 20-30%만 추가)

#### 3. 라우팅 문법
- **변경**: `fn route GET /path` → V 스타일 래핑
  ```fv
  fn handle_users_get(req: Request) -> Response { ... }
  // 또는 V 메타프로그래밍으로 래핑
  ```
- **노력**: 낮음 (호환성 레이어만 추가)

---

## 통합 이점

### 1. 코드 재사용
- 20,370줄 기존 라이브러리 그대로 사용
- 393개 테스트 자동 검증

### 2. 빠른 개발
- V 문법만 추가하면 됨
- 내부는 FreeLang과 동일

### 3. 커뮤니티
- V 개발자는 쉬운 문법 사용
- FreeLang 팬은 풍부한 라이브러리 사용
- Win-win

### 4. 프로덕션 준비
- 즉시 프로덕션 배포 가능
- Docker, K8s 지원
- 130개 프로젝트 생태계

---

## Phase 2 준비 사항

### Lexer 어댑터
V 키워드를 FreeLang 토큰으로 매핑:
- `fn` → `fn` (동일)
- `:=` → `:=` (동일)
- `?` → `?` (동일)
- `mut` → `mut` (동일)

### Parser 어댑터
V 문법 규칙을 FreeLang AST로 변환:
- V 함수 정의 → FreeLang 함수 정의
- V 구조체 → FreeLang 타입
- V 인터페이스 → FreeLang 트레이트

### 호환성 테스트
V 예제 코드를 수집하여 컴파일 테스트

---

**상태**: 🟡 분석 완료 대기 → Task 1.3 시작 예정
