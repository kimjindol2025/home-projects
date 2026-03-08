# FreeLang 셀프호스팅 완성 계획 (C Transpiler 방식)

## Context

**목표**: FreeLang이 자기 자신을 컴파일하는 셀프호스팅 언어를 완성한다.
**전략**: FreeLang → C 코드 생성 → clang → ARM64 네이티브 바이너리
**기준**: 시간보다 정확성. 증거 없으면 "완료" 선언 불가.

### 환경
- OS: ARM64 Android (Termux)
- 컴파일러: clang 21.1.8 (aarch64-unknown-linux-android30)
- Node.js: TypeScript (v6 현재 컴파일러)
- 저장소: `freelang-c-final/`

---

## 셀프호스팅 체인

```
Stage 0 (지금):
  TypeScript 컴파일러 (v6)
  → FreeLang 소스를 C 코드로 변환
  → clang으로 네이티브 바이너리 생성

Stage 1 (Phase 3):
  FreeLang으로 작성한 FreeLang 컴파일러 (.fl 파일들)
  → Stage 0이 이것을 C 코드로 변환
  → clang → fl-compiler 바이너리

Stage 2 (Phase 4 — 진짜 셀프호스팅):
  Stage 1 바이너리 (fl-compiler)가
  자신의 소스 (.fl 파일들)를 컴파일
  → 동일한 fl-compiler 바이너리 생성

검증:
  md5(Stage1 바이너리) == md5(Stage2 바이너리) ← 이게 증명
```

---

## Phase별 작업 계획

---

### Phase 1: 언어 기반 완성

**목적**: 제네릭, Enum, Result<T,E>가 없으면 컴파일러를 FreeLang으로 쓸 수 없다.

**작업 저장소**: `freelang-v6/` (기존 v6를 확장)

#### 1-1. 타입 어노테이션 AST 보존

**문제**: 현재 `parser.ts`의 `parseParams()`가 타입 어노테이션을 소비 후 버림.

**수정 파일**: `freelang-v6/src/ast.ts`, `freelang-v6/src/parser.ts`

```typescript
// ast.ts 수정: fn 노드의 params 타입 변경
// Before:
//   params: string[]
// After:
type Param = { name: string; typeAnnotation?: string };

// Stmt.fn 수정:
{ kind: "fn"; name: string; params: Param[]; returnType?: string; body: Stmt[] }

// Stmt.struct 수정: 이미 type: string이 있으나 제네릭 대응 필요
{ kind: "struct"; name: string; typeParams?: string[]; fields: Array<{name: string; type: string}> }
```

**수정 파일**: `parser.ts`의 `parseParams()` (라인 669-688)
- 타입 어노테이션을 버리지 않고 `Param` 객체에 저장

#### 1-2. 사용자 정의 Enum (Tagged Union)

**수정 파일**: `freelang-v6/src/token.ts`, `ast.ts`, `parser.ts`, `compiler.ts`, `vm.ts`

```typescript
// token.ts: Enum 키워드 추가
Enum = "enum",

// ast.ts: enum 노드 추가
type EnumVariant = { name: string; fields?: string[] }; // fields = 데이터 필드
{ kind: "enum"; name: string; variants: EnumVariant[] }

// 생성 예시:
// enum Token { Number(f64), Ident(string), Plus, EOF }
// enum Result { Ok(value), Err(message) }
```

```typescript
// compiler.ts: Value 타입 확장
| { tag: "enum_val"; enumName: string; variant: string; data?: Value }

// Opcode 추가:
EnumNew,    // 열거형 변체 생성
EnumMatch,  // 열거형 패턴 매칭
EnumGetData // 열거형 내부 데이터 추출
```

#### 1-3. 제네릭 타입 시스템

**수정 파일**: `token.ts`, `ast.ts`, `parser.ts`, `type-inference.ts`

```typescript
// 목표: 이런 코드가 파싱되고 타입 체크돼야 함
fn map<T, U>(arr: [T], f: fn(T) -> U) -> [U] { ... }
struct HashMap<K, V> { ... }
let result: Result<i32, string> = Ok(42)
```

```typescript
// ast.ts: 제네릭 파라미터 추가
type TypeExpr =
  | { kind: "simple"; name: string }         // i32, string
  | { kind: "generic"; name: string; args: TypeExpr[] }  // Result<T, E>
  | { kind: "array"; elem: TypeExpr }        // [T]
  | { kind: "fn"; params: TypeExpr[]; ret: TypeExpr }  // fn(T) -> U
  | { kind: "typevar"; name: string };       // T (타입 변수)

// Stmt.fn 추가:
typeParams?: string[]; // ["T", "U"]

// Stmt.struct 추가:
typeParams?: string[]; // ["K", "V"]
```

**구현 전략**: 1단계에서는 **타입 소거(type erasure)** 방식으로 구현.
- 제네릭 함수는 컴파일 시 void* 또는 FreeLang 동적 타입으로 처리
- 정확한 특수화는 Phase 2 C 생성기에서 처리

#### 1-4. 내장 Result<T, E> 언어 레벨 타입

```typescript
// FreeLang 문법:
let r: Result<i32, string> = Ok(42)
let e: Result<i32, string> = Err("fail")

match r {
  Ok(v) => print(v),
  Err(msg) => print(msg),
}
```

- Enum의 특수 케이스로 구현: `enum Result<T, E> { Ok(T), Err(E) }`
- `?` 연산자: `let v = r?` (Err이면 조기 return)

**검증**: 각 기능 추가 후 테스트 파일 실행으로 확인
- `tests/enum_basic.fl`
- `tests/generic_basic.fl`
- `tests/result_basic.fl`

---

### Phase 2: C 코드 생성 백엔드

**목적**: FreeLang AST → C 소스코드 → clang → ARM64 바이너리

**작업 저장소**: `freelang-c-final/` 및 `freelang-v6/src/codegen/`

#### 2-1. 런타임 라이브러리 (C)

**신규 파일**: `freelang-c-final/runtime/fl_runtime.h`, `fl_runtime.c`

```c
// fl_runtime.h - FreeLang 값 표현의 핵심

// 1. 동적 값 태그드 유니온 (Phase 2 초기: 모든 값을 이 타입으로)
typedef enum {
    FL_NULL, FL_BOOL, FL_INT, FL_FLOAT,
    FL_STRING, FL_ARRAY, FL_MAP,
    FL_CLOSURE, FL_ENUM,
} fl_tag;

typedef struct fl_value {
    fl_tag tag;
    union {
        int64_t  ival;
        double   fval;
        int      bval;
        char*    sval;
        struct fl_array*   arrval;
        struct fl_map*     mapval;
        struct fl_closure* fnval;
        struct fl_enum_val* enumval;
    };
} fl_value;

// 2. 배열
typedef struct fl_array {
    fl_value** data;
    size_t     length;
    size_t     capacity;
} fl_array;

// 3. 해시맵
typedef struct fl_map_entry { char* key; fl_value* val; struct fl_map_entry* next; } fl_map_entry;
typedef struct fl_map { fl_map_entry** buckets; size_t bucket_count; size_t size; } fl_map;

// 4. 클로저
typedef struct fl_closure {
    fl_value* (*fn_ptr)(void*, fl_value**, int);
    void*  env;     // 캡처된 변수들
    int    arity;
} fl_closure;

// 5. Enum 변체
typedef struct fl_enum_val {
    const char* enum_name;
    const char* variant;
    fl_value**  data;    // 데이터 필드들
    int         data_count;
} fl_enum_val;

// 6. 기본 연산 함수들
fl_value* fl_int(int64_t v);
fl_value* fl_float(double v);
fl_value* fl_bool(int v);
fl_value* fl_string(const char* s);
fl_value* fl_null(void);
fl_value* fl_add(fl_value* a, fl_value* b);
fl_value* fl_sub(fl_value* a, fl_value* b);
fl_value* fl_eq(fl_value* a, fl_value* b);
void      fl_print(fl_value* v);
fl_value* fl_array_new(void);
void      fl_array_push(fl_array* arr, fl_value* v);
fl_value* fl_enum_new(const char* ename, const char* variant, fl_value** data, int n);

// 7. 메모리 관리 (초기: arena allocator)
void* fl_alloc(size_t size);
void  fl_gc_run(void);

// 8. 예외 처리
#include <setjmp.h>
extern jmp_buf  fl_exception_buf;
extern fl_value* fl_last_exception;
#define FL_TRY    if (setjmp(fl_exception_buf) == 0)
#define FL_CATCH  else
void fl_throw(fl_value* err);
```

#### 2-2. C 코드 생성기 (TypeScript)

**신규 파일**: `freelang-v6/src/codegen/c-codegen.ts`

```typescript
// AST 노드 → C 소스코드 문자열 생성
export class CCodegen {
  private lines: string[] = [];
  private indent = 0;

  generate(stmts: Stmt[]): string {
    this.emitLine('#include "fl_runtime.h"');
    for (const s of stmts) this.genStmt(s);
    this.emitLine('int main(int argc, char** argv) { return (int)fl_int(0)->ival; }');
    return this.lines.join('\n');
  }

  genStmt(s: Stmt): void { /* switch(s.kind) ... */ }
  genExpr(e: Expr): string { /* switch(e.kind) ... */ }
}
```

**핵심 변환 규칙**:

| FreeLang 구문 | 생성되는 C 코드 |
|-------------|---------------|
| `fn foo(a: i32) -> i32 { }` | `fl_value* fl_foo(fl_value* a) { }` |
| `struct Point { x: i32, y: i32 }` | `fl_value* fl_Point_new(fl_value* x, fl_value* y);` |
| `let x = 42` | `fl_value* x = fl_int(42);` |
| `x + y` | `fl_add(x, y)` |
| `if cond { } else { }` | `if (fl_is_truthy(cond)) { } else { }` |
| `for item in arr { }` | `for (size_t _i=0; _i<arr->arrval->length; _i++) { fl_value* item = arr->arrval->data[_i]; }` |
| `match x { Ok(v) => ... }` | `if (x->enumval && strcmp(x->enumval->variant,"Ok")==0) { fl_value* v = x->enumval->data[0]; ... }` |
| `print(x)` | `fl_print(x);` |
| `fn(a) { return a }` (클로저) | `fl_closure* _cl_N` + 환경 구조체 + 함수 포인터 |

#### 2-3. CLI 통합 (`freelang-v6/src/main.ts`)

```bash
# 현재:
node dist/main.js input.fl          # 실행

# 추가:
node dist/main.js --emit-c input.fl  # C 코드 출력
node dist/main.js --build input.fl   # C 생성 + clang 자동 빌드
```

**검증**: `tests/` 디렉토리의 모든 .fl 파일을 `--build`로 컴파일 후 실행 결과 일치 확인

---

### Phase 3: FreeLang 컴파일러 자작 (.fl)

**목적**: FreeLang 언어로 FreeLang 컴파일러를 작성한다. (Stage 1)

**작업 저장소**: `freelang-c-final/compiler/`

#### 3-1. 파일 구조

```
freelang-c-final/
├── runtime/
│   ├── fl_runtime.h
│   └── fl_runtime.c
├── compiler/
│   ├── lexer.fl         ← FreeLang 렉서
│   ├── token.fl         ← 토큰 타입 (Enum 활용)
│   ├── parser.fl        ← FreeLang 파서
│   ├── ast.fl           ← AST 타입 (struct + Enum)
│   ├── types.fl         ← 타입 시스템
│   ├── codegen_c.fl     ← C 코드 생성기
│   └── main.fl          ← 진입점
├── tests/
│   ├── hello_world.fl
│   ├── fibonacci.fl
│   ├── enum_test.fl
│   └── generic_test.fl
└── SELFHOSTING_CHECKLIST.md
```

#### 3-2. lexer.fl

```fl
// 토큰 타입을 Enum으로 정의
enum TokenKind {
  Number(f64),
  StringLit(string),
  Ident(string),
  Fn, Let, If, Else, While, For, In,
  Return, Struct, Enum, Match,
  Plus, Minus, Star, Slash,
  Eq, Neq, Lt, Gt, Lte, Gte,
  Assign, Arrow, FatArrow,
  LParen, RParen, LBrace, RBrace,
  Comma, Colon, Semicolon, Dot,
  EOF,
}

struct Token {
  kind: TokenKind,
  line: i32,
  col: i32,
}

fn tokenize(source: string) -> [Token] {
  // ...
}
```

#### 3-3. ast.fl

```fl
// AST 노드를 Enum으로 정의
enum Expr {
  Number(f64),
  StringLit(string),
  Bool(bool),
  Null,
  Ident(string),
  Binary(string, Expr, Expr),    // op, left, right
  Call(Expr, [Expr]),             // callee, args
  Member(Expr, string),           // obj.field
  Match(Expr, [MatchArm]),
}

enum Stmt {
  Let(string, bool, Expr),        // name, mutable, init
  Fn(string, [Param], string, [Stmt]),  // name, params, rettype, body
  Struct(string, [Field]),
  EnumDecl(string, [Variant]),
  If(Expr, [Stmt], [Stmt]),       // cond, then, else
  While(Expr, [Stmt]),
  For(string, Expr, [Stmt]),      // var, iterable, body
  Return(Expr),
  Print(Expr),
}

struct Param { name: string, type_annotation: string }
struct Field { name: string, field_type: string }
struct Variant { name: string, fields: [string] }
struct MatchArm { pattern: Expr, body: [Stmt] }
```

#### 3-4. codegen_c.fl (C 코드 생성기)

```fl
// FreeLang으로 작성한 C 코드 생성기

fn codegen_stmt(stmt: Stmt) -> string {
  match stmt {
    Stmt.Let(name, _, init) => {
      let val = codegen_expr(init)
      return "fl_value* " + name + " = " + val + ";\n"
    },
    Stmt.Fn(name, params, ret_type, body) => {
      // ... C 함수 생성
    },
    Stmt.If(cond, then_body, else_body) => {
      // ...
    },
    // ...
  }
}

fn codegen_expr(expr: Expr) -> string {
  match expr {
    Expr.Number(n)      => "fl_int(" + to_string(n) + ")",
    Expr.StringLit(s)   => "fl_string(\"" + s + "\")",
    Expr.Binary(op, l, r) => {
      let lc = codegen_expr(l)
      let rc = codegen_expr(r)
      match op {
        "+" => "fl_add(" + lc + ", " + rc + ")",
        "-" => "fl_sub(" + lc + ", " + rc + ")",
        // ...
      }
    },
    // ...
  }
}
```

#### 3-5. 빌드 절차 (Stage 1)

```bash
# Stage 1 빌드: TypeScript 컴파일러(v6)가 .fl 파일들 → C → 바이너리
node dist/main.js --build compiler/main.fl -o fl-compiler-stage1

# 검증: Stage 1 컴파일러가 hello.fl을 컴파일 가능한지
./fl-compiler-stage1 --build tests/hello_world.fl -o hello
./hello  # → "Hello, World!" 출력
```

---

### Phase 4: 부트스트랩 검증

**목적**: Stage 2 빌드 후 바이너리 동일성 검증 (진짜 셀프호스팅 증명)

#### 4-1. Stage 2 빌드

```bash
# Stage 1 바이너리가 자신의 소스를 컴파일
./fl-compiler-stage1 --build compiler/main.fl -o fl-compiler-stage2

# 검증: 두 바이너리가 동일한가?
md5sum fl-compiler-stage1 fl-compiler-stage2
# 출력:
# a3f8...  fl-compiler-stage1
# a3f8...  fl-compiler-stage2   ← 동일해야 함!
```

#### 4-2. Stage 3 (안정성 확인)

```bash
# Stage 2가 또 다시 컴파일
./fl-compiler-stage2 --build compiler/main.fl -o fl-compiler-stage3

md5sum fl-compiler-stage2 fl-compiler-stage3
# 동일해야 함 → deterministic build 증명
```

#### 4-3. TypeScript 컴파일러 제거 테스트

```bash
# Node.js 없이 FreeLang 소스 컴파일 가능한가?
./fl-compiler-stage2 --build tests/fibonacci.fl -o fib
./fib  # 실행 성공 = 완전한 셀프호스팅
```

---

## 구체적 파일 수정 목록

### Phase 1 수정 파일

| 파일 | 변경 내용 |
|------|----------|
| `freelang-v6/src/token.ts` | `Enum`, `Lt`(제네릭 구분) 토큰 추가 |
| `freelang-v6/src/ast.ts` | `Param` 타입 변경, `enum` 노드, 제네릭 `TypeExpr` 추가 |
| `freelang-v6/src/parser.ts` | `parseParams()` 타입 보존, `parseEnum()` 추가, 제네릭 파싱 |
| `freelang-v6/src/compiler.ts` | `enum` 컴파일, `Value` 확장 |
| `freelang-v6/src/vm.ts` | enum 관련 Op 처리 |

### Phase 2 신규 파일

| 파일 | 내용 |
|------|------|
| `freelang-c-final/runtime/fl_runtime.h` | 모든 FreeLang C 타입 정의 |
| `freelang-c-final/runtime/fl_runtime.c` | 런타임 함수 구현 |
| `freelang-v6/src/codegen/c-codegen.ts` | AST → C 코드 생성기 |
| `freelang-v6/src/codegen/c-types.ts` | 타입 매핑 유틸리티 |

### Phase 3 신규 파일 (.fl)

| 파일 | 내용 |
|------|------|
| `freelang-c-final/compiler/token.fl` | TokenKind enum |
| `freelang-c-final/compiler/lexer.fl` | 렉서 (tokenize 함수) |
| `freelang-c-final/compiler/ast.fl` | AST enum/struct 정의 |
| `freelang-c-final/compiler/parser.fl` | 파서 |
| `freelang-c-final/compiler/types.fl` | 타입 시스템 |
| `freelang-c-final/compiler/codegen_c.fl` | C 코드 생성기 |
| `freelang-c-final/compiler/main.fl` | 진입점 + CLI |

---

## 검증 기준 (단계별)

### Phase 1 완료 조건

```
✅ enum Token { ... } 선언 후 match로 패턴 분해 가능
✅ fn foo<T>(x: T) -> T { x } 선언 및 호출 가능
✅ let r: Result<i32, string> = Ok(42) 동작
✅ r? 연산자 동작
✅ 기존 77개 테스트 전부 통과 (회귀 없음)
```

### Phase 2 완료 조건

```
✅ freelang-v6$ node dist/main.js --build tests/hello_world.fl -o hello
✅ ./hello 실행 → "Hello, World!" 출력
✅ --build로 컴파일된 fibonacci, struct, enum, closure 테스트 통과
✅ 생성된 .c 파일이 clang으로 경고 없이 컴파일
```

### Phase 3 완료 조건

```
✅ Stage 1 바이너리 생성됨 (fl-compiler-stage1 실행 가능)
✅ fl-compiler-stage1이 간단한 .fl 파일 컴파일 가능
✅ fl-compiler-stage1이 tests/ 전체 컴파일 성공
```

### Phase 4 완료 조건 (셀프호스팅 증명)

```
✅ md5(fl-compiler-stage1) == md5(fl-compiler-stage2)
✅ md5(fl-compiler-stage2) == md5(fl-compiler-stage3)
✅ TypeScript 없이 fl-compiler-stage2만으로 .fl 컴파일 가능
✅ 3회 연속 재빌드 동일 결과
```

---

## 정직한 현황 기록

### 현재 상태 (2026-03-08 기준)

| 항목 | 상태 |
|------|------|
| Stage 0 (TypeScript 컴파일러) | ✅ 동작 중 |
| Phase 1 (Enum/제네릭) | ❌ 미시작 |
| Phase 2 (C 코드 생성) | ❌ 미시작 |
| Phase 3 (FreeLang 자작 컴파일러) | ❌ 미시작 |
| Phase 4 (부트스트랩 검증) | ❌ 미시작 |

### 진행하면서 업데이트할 것

각 Phase 완료 시 `freelang-c-final/SELFHOSTING_CHECKLIST.md` 업데이트.
거짓 보고 금지 — 실행 증거(md5, 바이너리, 테스트 로그)만 기록.

---

## 시작 순서

**첫 번째 작업**: `freelang-v6/src/`에서 Phase 1-1 시작
→ `ast.ts`의 `Param` 타입 수정 + `parser.ts`의 `parseParams()` 수정
→ 타입 어노테이션이 AST에 보존되는지 테스트

**두 번째 작업**: `token.ts` + `ast.ts`에 `enum` 노드 추가
→ `parser.ts`에 `parseEnum()` 구현
→ `compiler.ts`에 enum 컴파일 추가

**세 번째 작업**: 제네릭 타입 파싱 (타입 소거 방식)

순서대로 하나씩, 검증하면서 진행.
