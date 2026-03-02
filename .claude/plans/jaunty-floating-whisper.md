# Phase E: 런타임 구현 계획

## Context

Phase A~D를 통해 FreeLang의 핵심 시스템(동시성, 표준 라이브러리, FPM, 문서화)을 완성했습니다.
Phase E 목표: **기존 Phase B(sync) + Phase C(stdlib) 코드를 실제로 실행 가능**한 런타임 엔진 구축.
기존 코드 재활용을 극대화하고, REPL(대화형 실행 환경)까지 완성하는 것이 목표입니다.

---

## 현재 상태

```
freelang-backend-system/
├── src/sync/    (16개 .fl - Phase B 동시성 명세 완료)
├── src/stdlib/  (11개 .fl - Phase C 표준 라이브러리 105+함수 명세 완료)
├── fpm/         (5개  .fl - Phase D 패키지 관리자 명세 완료)
└── runtime/     (없음 → Phase E에서 신규 생성)
```

핵심 기반:
- `FlValue` enum (Int/Float/Bool/Str/Array/Object/Null) → mod.fl에 정의
- 105+ stdlib 함수 명세 (types, io, math, string, array, time, system, json, object, error)
- 동시성 모듈 명세 (Arc, Mutex, Channel, ThreadPool, RefCell, CondVar)

---

## Phase E 파일 구조 (신규 생성, 6개)

```
Week 1 (Lexer + Parser):
└── runtime/lexer.fl             (500줄, 30+ 토큰, 8개 함수, 10개 테스트)

Week 2 (Evaluator + stdlib 연결):
├── runtime/evaluator.fl         (400줄, FlValue 평가기, 8개 함수, 10개 테스트)
└── runtime/stdlib_bindings.fl   (200줄, 30+개 내장 함수 등록, 10개 테스트)

Week 3 (sync 모듈 런타임 연결):
└── runtime/sync_bindings.fl     (600줄, Arc/Mutex/Channel/ThreadPool 런타임, 10개 테스트)

Week 4 (REPL + 에러 처리):
├── runtime/repl.fl              (250줄, Read-Eval-Print Loop, 5개 테스트)
└── runtime/error.fl             (150줄, 런타임 에러 타입, 5개 테스트)
```

---

## 주차별 구현 계획

### Week 1: Lexer + Parser (500줄)
**파일**: `runtime/lexer.fl`

**Token 종류 (30+개)**:
```rust
pub enum Token {
    // 리터럴
    Int(i64), Float(f64), Bool(bool), Str(String), Null,

    // 키워드
    Let, Fn, If, Else, While, For, Return, Use, Match,

    // 연산자
    Plus, Minus, Star, Slash, Percent,
    Eq, Neq, Lt, Gt, Le, Ge,
    And, Or, Not, Assign, Arrow,

    // 구분자
    LParen, RParen, LBrace, RBrace, LBracket, RBracket,
    Comma, Semicolon, Colon, Dot,

    // 식별자 / 특수
    Ident(String), EOF, Illegal(char),
}
```

**AST 노드 (20+개)**:
```rust
pub enum Expr {
    Int(i64), Float(f64), Bool(bool), Str(String), Null,
    Ident(String),
    BinOp { left: Box<Expr>, op: BinOpKind, right: Box<Expr> },
    UnaryOp { op: UnaryOpKind, operand: Box<Expr> },
    Call { func: String, args: Vec<Expr> },
    Index { obj: Box<Expr>, index: Box<Expr> },
    Array(Vec<Expr>),
    Object(Vec<(String, Expr)>),
    If { cond: Box<Expr>, then_block: Vec<Stmt>, else_block: Option<Vec<Stmt>> },
}

pub enum Stmt {
    Let { name: String, value: Expr },
    Return(Expr),
    Expr(Expr),
    While { cond: Expr, body: Vec<Stmt> },
}

pub struct Program { pub stmts: Vec<Stmt> }
```

**핵심 함수 (8개)**:
- `tokenize(src: &str) -> Vec<Token>` — 소스 → 토큰 목록
- `parse(tokens: &[Token]) -> Result<Program, ParseError>` — 토큰 → AST
- `parse_expr(...)`, `parse_stmt(...)`, `parse_fn_call(...)`, `parse_block(...)`

**테스트 (10개)**:
- 리터럴 토크나이징 (정수, 실수, 문자열, bool)
- 연산자 파싱 (+, -, *, /, ==, !=)
- 함수 호출 파싱 (`println("hello")`)
- let 문 파싱, if-else 파싱, 배열/객체 리터럴 파싱

---

### Week 2: Evaluator + stdlib 내장 함수 (600줄)

#### evaluator.fl (400줄)
```rust
pub struct Environment {
    vars: HashMap<String, FlValue>,
    parent: Option<Box<Environment>>,
}

pub struct Evaluator {
    env: Environment,
    builtins: HashMap<String, BuiltinFn>,
}

type BuiltinFn = fn(Vec<FlValue>) -> Result<FlValue, RuntimeError>;
```

핵심 함수 (8개):
- `eval(prog, env)` — Program 전체 실행
- `eval_expr(expr, env)` — 표현식 평가
- `eval_stmt(stmt, env)` — 문장 실행
- `eval_binop(op, left, right)` — 이항 연산
- `eval_call(func, args, env)` — 함수 호출

#### stdlib_bindings.fl (200줄)

30+개 내장 함수 등록:

| 모듈 | 함수 |
|------|------|
| IO | `println`, `print`, `read_file`, `write_file` |
| String | `str_length`, `str_concat`, `str_split`, `str_format`, `str_upper`, `str_lower` |
| Array | `array_push`, `array_length`, `array_map`, `array_filter`, `array_reduce`, `array_sort` |
| Math | `sqrt`, `abs`, `pow`, `sin`, `cos`, `floor`, `ceil` |
| Types | `type_of`, `to_int`, `to_float`, `to_string`, `to_bool`, `is_null` |
| JSON | `json_parse`, `json_stringify` |
| Time | `now_ms`, `sleep` |

**테스트 (10개)**:
- `println("hello")` 실행
- `str_length("hello") == 5`
- `array_map([1,2,3], fn(x) x*2) == [2,4,6]`
- `json_parse("{\"a\":1}").a == 1`

---

### Week 3: sync 모듈 런타임 연결 (600줄)
**파일**: `runtime/sync_bindings.fl`

| 모듈 | 내장 함수 |
|------|---------|
| Arc | `arc_new`, `arc_clone`, `arc_get` |
| Mutex | `mutex_new`, `mutex_lock`, `mutex_set` |
| Channel | `channel_new`, `channel_send`, `channel_recv` |
| ThreadPool | `threadpool_new`, `threadpool_execute` |
| Thread | `thread_spawn`, `thread_sleep` |
| CondVar | `condvar_new`, `condvar_wait`, `condvar_notify` |

각 함수: `FlValue` 인수 받아 `Result<FlValue, RuntimeError>` 반환.
`thread_spawn`은 클로저(함수 이름 + 인수)를 받아 새 스레드에서 평가.

**테스트 (10개)**:
- `arc_new(42)`, `arc_get(arc)` 검증
- `mutex_new(0)`, `mutex_lock`, `mutex_set` 시나리오
- `channel_new`, `channel_send(tx, 42)`, `channel_recv(rx)` 파이프라인
- `threadpool_new(4)`, `threadpool_execute` 작업 분산

---

### Week 4: REPL + 에러 처리 (400줄)

#### repl.fl (250줄)
```rust
pub struct Repl {
    evaluator: Evaluator,
    history: Vec<String>,
}
```

**REPL 특수 명령**:
- `.help` — 명령 목록
- `.history` — 명령 히스토리
- `.env` — 현재 변수 목록
- `.load <file>` — .fl 파일 실행
- `exit` / `quit` — 종료

#### error.fl (150줄)
```rust
pub enum RuntimeError {
    TypeError { expected: &'static str, got: &'static str },
    UndefinedVariable(String),
    IndexOutOfBounds { index: i64, length: usize },
    CallError { name: String, reason: String },
    ParseError { line: usize, col: usize, msg: String },
    LockError(String),
    ChannelClosed(String),
}
```

**테스트 (10개)**:
- REPL 루프 시작/종료
- `.help` / `.env` 명령
- 에러 발생 후 REPL 계속 실행
- TypeError / UndefinedVariable 메시지 검증

---

## 최종 목표

| 항목 | Week 1 | Week 2 | Week 3 | Week 4 | 합계 |
|------|--------|--------|--------|--------|------|
| 파일 | 1 | 2 | 1 | 2 | **6개** |
| 줄수 | 500 | 600 | 600 | 400 | **2,100줄** |
| 테스트 | 10 | 20 | 10 | 10 | **50개** |

**달성 시 가능한 것**:
```bash
$ freelang --repl         # REPL 실행
FreeLang REPL v1.0
> let x = 42
> println("x = " + to_string(x))
x = 42
> array_map([1,2,3], fn(n) n * 2)
[2, 4, 6]
> let m = mutex_new(0)
> mutex_lock(m)
```

---

## 완료 기준

- ✅ 50개 테스트 통과
- ✅ REPL에서 FreeLang 코드 실행 가능
- ✅ 30+개 stdlib 내장 함수 등록 및 동작
- ✅ 6개 sync 모듈 런타임 연결
- ✅ GOGS 커밋 4회 (Week별)
