---
name: FV 2.0 Phase 1 - 통합 지점 설계
description: V 언어를 FreeLang 컴파일러에 통합하기 위한 구체적인 설계 및 구현 전략
type: project
---

# FV 2.0 Phase 1: 통합 지점 설계

**작성일**: 2026-03-19
**프로젝트**: FV 2.0 (V Language + FreeLang Integration)
**상태**: 🟢 설계 완료

---

## 핵심 통합 아키텍처

### 현재 (FreeLang)
```
FreeLang 소스 (.fl)
    ↓
FreeLang Lexer
    ↓
FreeLang Parser → AST
    ↓
Type Checker
    ↓
Code Generator
    ↓
C 코드 / 바이너리
```

### FV 2.0 (목표)
```
FV 2.0 소스 (V 문법, .fv)
    ↓
V-compatible Lexer
    ↓
V-compatible Parser → V AST
    ↓
AST Adapter (V AST → FreeLang AST)
    ↓
Type Checker (기존 재사용)
    ↓
Code Generator (기존 재사용)
    ↓
C 코드 / 바이너리 (동일)
```

### 핵심 개선
- **Green Path**: Lexer + Parser (20% 신규 작성)
- **Blue Path**: AST Adapter (기존 구조 매핑)
- **전체**: Type Checker ~ Code Generator (100% 재사용)

---

## 문법 매핑 테이블

### 1. 기본 선언

| V | FreeLang | FV 2.0 |
|---|----------|--------|
| `fn name() { }` | `fn name() { }` | `fn name() { }` |
| `let x = 5` | `let x = 5` | `let x = 5` |
| `mut x := 5` | `let mut x = 5` | `mut x := 5` |
| `const X = 5` | `const X = 5` | `const X = 5` |

**적응**: 거의 동일. Lexer에서 `let mut` → `mut` 변환

### 2. 타입 선언

| V | FreeLang | FV 2.0 (매핑) |
|---|----------|--------------|
| `struct User { }` | `type User = { }` | `type User = { }` |
| `interface Reader { }` | `trait Reader { }` | `trait Reader { }` |
| `type MyInt = int` | `type MyInt = i64` | `type MyInt = i64` |

**적응**: V `struct` → FreeLang `type` (키워드 변환)

### 3. 타입 시스템

#### 정수형
| V | FreeLang | 매핑 |
|---|----------|------|
| `i8, i16, i32, i64` | `i8, i16, i32, i64` | 동일 |
| `u8, u16, u32, u64` | `u8, u16, u32, u64` | 동일 |
| `int` | `i64` | `int` → `i64` (기본) |
| `uint` | `u64` | `uint` → `u64` (기본) |

#### 컬렉션
| V | FreeLang | 매핑 |
|---|----------|------|
| `[]T` | `Vec(T)` | 동적 배열 |
| `[N]T` | `[T; N]` | 고정 배열 |
| `map[K]V` | `HashMap(K, V)` | 해시맵 |
| `?T` | `Option(T)` | NULL 안전 |
| `!E` | `Result(T, E)` | 에러 처리 |

### 4. 에러 처리

#### V 스타일
```v
fn read_file(path string) ?string {
  return none  // 에러
}

value := read_file('file.txt')?  // 에러 전파
```

#### FreeLang 스타일
```fv
fn read_file(path: String) -> Option(String) {
  return None
}

value := read_file("file.txt")?
```

**매핑**: V `?T` = FreeLang `Option(T)` (동일)

### 5. 함수 정의

#### V
```v
fn add(a int, b int) int {
  a + b
}
```

#### FreeLang
```fv
fn add(a: i64, b: i64) -> i64 {
  a + b
}
```

#### 변환 규칙
1. `fn` 키워드는 유지
2. 매개변수: `(name type)` → `(name: type)`
3. 반환형: `-> type` 유지
4. 본체는 동일

**Parser 규칙**:
```
// V 문법
fn_def = "fn" IDENT "(" params ")" (IDENT)? "{" body "}"

// 변환
params: IDENT IDENT → IDENT: IDENT (타입 추가)
```

### 6. HTTP 라우팅

#### V 에상 문법 (RESTful)
```v
fn GET /api/users -> []User {
  db.all_users()?
}

fn GET /api/users/:id -> User {
  db.find(id)?
}

fn POST /api/users -> User {
  user := User { ... }
  db.save(user)?
  user
}
```

#### FreeLang 현재 문법
```fv
fn route GET /api/users -> JSON {
  users := db.all_users()?
  json(users)
}
```

#### 변환 전략
```
FV 2.0:
fn GET /api/users -> []User {
  ...
}

변환 →

FreeLang:
fn route GET /api/users -> JSON {
  ...
  json(...)
}
```

**특수 처리**:
- `fn METHOD /path` 감지
- 자동으로 `fn route METHOD /path -> JSON`으로 변환
- 반환값을 `json()` 호출로 래핑

---

## AST 어댑터 설계

### V AST 구조 (예상)
```rust
enum VExpression {
  IntLiteral(i64),
  StringLiteral(String),
  Identifier(String),
  FunctionCall {
    name: String,
    args: Vec<Expression>,
  },
  If {
    condition: Box<Expression>,
    then_body: Vec<Statement>,
    else_body: Option<Vec<Statement>>,
  },
  // ...
}
```

### FreeLang AST 구조 (기존)
```rust
enum FLExpression {
  Int(i64),
  String(String),
  Var(String),
  Call {
    func: String,
    args: Vec<Expression>,
  },
  Cond {
    cond: Box<Expression>,
    then_e: Box<Expression>,
    else_e: Option<Box<Expression>>,
  },
  // ...
}
```

### 어댑터 함수
```rust
fn adapt_v_expr(v_expr: VExpression) -> FLExpression {
  match v_expr {
    VExpression::IntLiteral(n) => FLExpression::Int(n),
    VExpression::StringLiteral(s) => FLExpression::String(s),
    VExpression::Identifier(id) => FLExpression::Var(id),
    VExpression::FunctionCall { name, args } => {
      FLExpression::Call {
        func: name,
        args: args.into_iter().map(adapt_v_expr).collect(),
      }
    },
    // ...
  }
}
```

---

## 호환성 레벨 정의

### Level 1: 기본 타입 & 함수 (Week 1)
```
목표 호환율: 70%

✅ 기본 타입 (int, string, bool)
✅ 함수 정의
✅ 변수 선언
✅ 산술/논리 연산
✅ if/else
❌ 제네릭 (미루기)
❌ 고급 에러 처리 (미루기)
```

### Level 2: 구조체 & 메서드 (Week 2)
```
목표 호환율: 85%

✅ 구조체 정의
✅ 메서드
✅ 배열 & 맵
✅ for 루프
✅ match 문
❌ 인터페이스 (Week 3)
```

### Level 3: 완전 호환 (Week 3)
```
목표 호환율: 95%

✅ 인터페이스
✅ 제네릭
✅ 고급 에러 처리
✅ HTTP 라우팅
✅ 데이터베이스 쿼리
```

---

## 구현 체크리스트 (Phase 2)

### Task 2.1: Lexer 수정 (1주)

- [ ] V 키워드 추가
  - `fn`, `mut`, `let`, `const`
  - `struct`, `type`, `trait`
  - `match`, `if`, `else`, `for`
  - `interface`, `enum`

- [ ] V 연산자 추가
  - `?` (에러 전파)
  - `:=` (변수 선언)
  - `..` (범위)

- [ ] 테스트
  - [ ] 토큰화 테스트 (20개 예제)
  - [ ] 에러 처리 테스트

### Task 2.2: Parser 수정 (1주)

- [ ] V 문법 규칙 추가
  - [ ] 함수 정의
  - [ ] 구조체 정의
  - [ ] 타입 정의
  - [ ] 제어문 (if, for, match)

- [ ] AST 생성
  - [ ] V AST 타입 정의
  - [ ] 어댑터 함수 구현

- [ ] 테스트
  - [ ] 파싱 테스트 (30개 프로그램)
  - [ ] AST 어댑터 테스트

### Task 2.3: 호환성 테스트 (3-4일)

- [ ] V 예제 코드 수집 (30-50개)
- [ ] 컴파일 테스트
- [ ] 호환율 측정 (목표: 95%)
- [ ] 실패 사례 분석

---

## 성공 기준

### 최종 목표 (Phase 1 완료)

| 항목 | 기준 | 측정 방법 |
|------|------|----------|
| V 문법 호환율 | 95% | V 예제 50개 중 47개 이상 컴파일 |
| 파서 정확도 | 100% | 테스트 30개 모두 통과 |
| 타입 매핑 | 100% | 모든 V 타입 → FreeLang 타입 매핑 |
| 성능 | <1초 | 평균 컴파일 시간 |
| 코드 재사용 | 90%+ | 기존 Type Checker, Code Generator 재사용 |

---

## Phase 2-4 로드맵 요약

### Phase 2 (Week 2-3): V 문법 채택
- **완료**: V 파서, AST 어댑터, 호환성 테스트
- **결과**: V 코드 95% 컴파일 가능

### Phase 3 (Week 4-7): 라이브러리 통합
- **완료**: HTTP, Database, WebSocket, gRPC를 V 문법으로 노출
- **결과**: 실제 프로덕션 백엔드 작성 가능

### Phase 4 (Week 8-9): 마케팅 & 배포
- **완료**: 문서, 예제, 130개 프로젝트 마이그레이션
- **결과**: FV 2.0 공식 출시

---

## 위험 요소 & 완화 전략

### 위험 1: 문법 충돌
**문제**: V와 FreeLang의 문법이 호환되지 않을 수 있음
**완화**: 주기적인 호환성 테스트 (주 2회)

### 위험 2: 성능 저하
**문제**: V 파서가 느릴 수 있음
**완화**: 기존 FreeLang 파서 최적화 활용

### 위험 3: 라이브러리 누락
**문제**: V 특화 라이브러리가 FreeLang에 없을 수 있음
**완화**: 가장 중요한 라이브러리부터 구현

---

**상태**: 🟢 Phase 1 분석 & 설계 완료

**다음**: Phase 2 (V 문법 채택) 시작 준비
