# FreeLang v1.0 - 공식 언어 명세서

**작성일**: 2026-03-02
**버전**: v1.0.0 (Phase 1 기본)
**상태**: 🟢 공식 사양
**목표**: AI가 쉽게 쓸 수 있는 안전한 프로그래밍 언어

---

## 📖 목차

1. [핵심 철학](#1-핵심-철학)
2. [언어 특징](#2-언어-특징)
3. [문법 (Syntax)](#3-문법-syntax)
4. [타입 시스템](#4-타입-시스템)
5. [의미론 (Semantics)](#5-의미론-semantics)
6. [컴파일러 아키텍처](#6-컴파일러-아키텍처)
7. [예제 모음](#7-예제-모음)
8. [표준 라이브러리](#8-표준-라이브러리)
9. [Phase 1 기능](#9-phase-1-기능)
10. [호환성 및 마이그레이션](#10-호환성-및-마이그레이션)

---

## 1. 핵심 철학

### 1.1 설계 목표

```
FreeLang = "AI를 위한 프로그래밍 언어"
```

**3가지 기본 원칙**:

1. **🎯 AI 친화성 (AI-Friendly)**
   - AI가 쉽게 읽고 생성할 수 있는 문법
   - 명확한 Intent 표현
   - 자동 완성 및 자동 생성 지원

2. **🔒 안전성 (Safety)**
   - 컴파일 타입 안전 (compile-time type safety)
   - 메모리 안전 (memory safety)
   - Null 안전 (null safety)

3. **📝 간결성 (Simplicity)**
   - 필수 요소만 작성 (intent, input, output)
   - 나머지는 자동 추론
   - 명확한 의도 표현

### 1.2 AI-First 철학

```
AI 코딩의 흐름:

[Natural Language] → [Intent] → [Function Signature]
      ↓                    ↓              ↓
  "합계 계산"         calculate_sum    array<number> → number
                                            ↓
                                    [AI Generates Code]
                                            ↓
                                    fn calculate_sum
                                      input: array<number>
                                      output: number
                                      intent: "배열의 합계"
                                      do
                                        result = 0
                                        for i in 0..input.len()
                                          result = result + input[i]
                                        return result
```

---

## 2. 언어 특징

### 2.1 핵심 특징 (Phase 1)

| 특징 | 설명 | 예시 |
|------|------|------|
| **Intent-Based** | 함수의 의도를 자연어로 선언 | `intent: "배열의 합계 계산"` |
| **Type-Safe** | 모든 값에 타입 정보 | `input: array<number>` |
| **Immutable-First** | 기본값은 불변 | `let x = 5` (변경 불가) |
| **Pattern Matching** | 패턴 기반 값 분해 | `match result { Ok(v) => ..., Err(e) => ... }` |
| **Error Handling** | Result 기반 오류 처리 | `fn foo() -> Result<number, Error>` |
| **Trait-Based** | 트레이트 기반 다형성 | `trait Iterable<T> { ... }` |
| **Generic** | 제너릭 타입 지원 | `fn process<T>(input: array<T>) -> array<T>` |

### 2.2 제거된 기능 (Phase 1 범위 외)

다음 기능들은 Phase 2+ 에서 추가됩니다:

- ❌ Concurrency (async/await, threads)
- ❌ Macros (매크로 시스템)
- ❌ Reflection (동적 타입 정보)
- ❌ FFI (외부 함수 인터페이스)
- ❌ Generics Specialization

---

## 3. 문법 (Syntax)

### 3.1 함수 선언

#### 형식 1: 최소 형식 (Minimal) ⭐ 권장

```freelang
fn functionName
  intent: "함수의 의도"
  input: paramName: paramType
  output: returnType
  do
    // 함수 본체
```

**예시**:
```freelang
fn sum
  intent: "배열의 합계 계산"
  input: numbers: array<number>
  output: number
  do
    result = 0
    for i in 0..numbers.len()
      result = result + numbers[i]
    return result
```

**특징**:
- 들여쓰기로 블록 표현 (중괄호 제외)
- 세미콜론 제외
- Intent로 함수의 의도를 명확히 함

#### 형식 2: 헤더 계약 형식 (Header)

```freelang
fn functionName: inputType → outputType ~ "intent"
```

**예시**:
```freelang
fn sum: array<number> → number ~ "배열의 합계"
  do
    result = 0
    for i in 0..input.len()
      result = result + input[i]
    return result
```

**특징**:
- 한 줄로 함수 시그니처 정의
- `→` (arrow)로 반환 타입 표현
- `~` (tilde)로 Intent 표현

#### 형식 3: 자유 형식 (Free - AI-Generated)

```freelang
fn sum(numbers: array<number>) -> number
  "배열의 합계를 계산한다"
  {
    result = 0
    for i in numbers
      result += i
    result
  }
```

**특징**:
- 일반적인 프로그래밍 언어 스타일
- AI가 생성하는 형식
- 괄호, 중괄호 사용 가능

### 3.2 변수 선언

```freelang
// 불변 변수 (기본)
let x = 5
let name = "Alice"
let numbers = [1, 2, 3]

// 가변 변수 (mut)
mut counter = 0
counter = counter + 1

// 타입 명시
let x: number = 5
let name: string = "Alice"
let items: array<string> = ["a", "b"]

// 타입 추론
let result = sum([1, 2, 3])  // 타입: number (추론됨)
```

### 3.3 제어문

#### if/else 문

```freelang
if x > 0
  print("양수")
else if x < 0
  print("음수")
else
  print("0")
```

#### match 문 (패턴 매칭)

```freelang
match result
  Ok(value) =>
    print(value)
  Err(error) =>
    print(error)

match number
  1 => print("하나")
  2 => print("둘")
  _ => print("기타")
```

#### for 루프

```freelang
// Range 기반
for i in 0..10
  print(i)

// 배열 반복
for item in array
  print(item)

// 인덱스 포함
for (i, item) in array.enumerate()
  print(f"{i}: {item}")
```

#### while 루프

```freelang
mut i = 0
while i < 10
  print(i)
  i = i + 1
```

### 3.4 구조체 선언

```freelang
struct Person
  name: string
  age: number
  email: string

// 인스턴스 생성
let person = Person { name: "Alice", age: 30, email: "alice@example.com" }

// 필드 접근
print(person.name)
```

### 3.5 Enum 선언

```freelang
enum Result<T, E>
  Ok(T)
  Err(E)

enum Color
  Red
  Green
  Blue

// 사용
let ok_result = Result.Ok(42)
let color = Color.Red
```

### 3.6 Trait 선언

```freelang
trait Drawable
  fn draw(self)

trait Comparable<T>
  fn compare(self, other: T) -> number

// Trait 구현
impl Drawable for Circle
  fn draw(self)
    print(f"Drawing circle with radius {self.radius}")
```

---

## 4. 타입 시스템

### 4.1 기본 타입

| 타입 | 설명 | 예시 |
|------|------|------|
| `number` | 64-bit 부동소수점 | `3.14`, `42` |
| `string` | UTF-8 문자열 | `"Hello"`, `'World'` |
| `bool` | 불린 값 | `true`, `false` |
| `()` | Unit (없음) | 반환값 없는 함수 |

### 4.2 복합 타입

```
// 배열
array<T>           // 배열
array<number>      // 숫자 배열
array<string>      // 문자 배열

// 맵
map<K, V>          // 키-값 맵
map<string, number> // 문자-숫자 맵

// 함수 타입
(number, number) -> number   // 입력 2개, 숫자 반환
() -> string                 // 입력 없음, 문자 반환

// 옵션/결과
Option<T>          // 값 또는 None
Result<T, E>       // 성공 또는 실패
```

### 4.3 타입 제약

```freelang
// 제너릭 타입 파라미터
fn first<T>(array: array<T>) -> Option<T>
  if array.len() > 0
    Some(array[0])
  else
    None

// Trait 제약
fn process<T: Drawable>(item: T)
  item.draw()

// 다중 제약
fn compare<T: Comparable<T> + Drawable>(a: T, b: T)
  // ...
```

---

## 5. 의미론 (Semantics)

### 5.1 타입 추론

FreeLang은 **Intent-Based Type Inference (IBTI)**를 사용합니다.

```freelang
fn calculate
  intent: "두 수를 더한 후 2배"
  input: a, b           // 타입 미지정
  output: number        // 반환 타입 명시
  do
    sum = a + b         // a, b는 number로 추론
    return sum * 2      // Result: number
```

### 5.2 불변성 (Immutability)

기본값은 불변입니다:

```freelang
let x = 5
x = 10          // ❌ 컴파일 에러

mut x = 5
x = 10          // ✅ OK

// 컬렉션도 불변
let arr = [1, 2, 3]
arr.push(4)     // ❌ 컴파일 에러

mut arr = [1, 2, 3]
arr.push(4)     // ✅ OK
```

### 5.3 오류 처리

모든 실패 가능한 작업은 Result를 반환합니다:

```freelang
enum Result<T, E>
  Ok(T)         // 성공
  Err(E)        // 실패

// 사용
fn divide(a: number, b: number) -> Result<number, string>
  if b == 0
    Err("Division by zero")
  else
    Ok(a / b)

// 호출
match divide(10, 2)
  Ok(result) =>
    print(result)        // 5
  Err(error) =>
    print(error)
```

### 5.4 생명주기 (Lifetime) - Phase 2+

생명주기는 Phase 2 이상에서 지원됩니다.

```freelang
// Phase 1: 문제 없음
fn first(arr: array<number>) -> number
  arr[0]

// Phase 2+: 생명주기 명시 가능
fn first<'a>(arr: &'a array<number>) -> &'a number
  &arr[0]
```

---

## 6. 컴파일러 아키텍처

### 6.1 컴파일 파이프라인

```
┌──────────────┐
│ Source Code  │  .fl 파일
└──────┬───────┘
       │
       v
┌──────────────┐
│   Lexer      │  토큰화
└──────┬───────┘
       │
       v
┌──────────────┐
│   Parser     │  구문 분석 → AST
└──────┬───────┘
       │
       v
┌──────────────────────┐
│ Type Checker         │  타입 검증
│ (Intent-Based)       │
└──────┬───────────────┘
       │
       v
┌──────────────┐
│ IR Generator │  중간 표현 생성
└──────┬───────┘
       │
       v
┌──────────────────────┐
│ Code Generator       │  최종 코드 생성
│ (TypeScript/C/VM)    │
└──────┬───────────────┘
       │
       v
┌──────────────┐
│ Output Code  │  .ts / .c / .bytecode
└──────────────┘
```

### 6.2 각 단계 상세

#### Lexer (토큰화)
- 소스 코드를 토큰 스트림으로 변환
- 공백, 주석 제거
- 토큰 타입: KEYWORD, IDENTIFIER, NUMBER, STRING, OPERATOR, ...

#### Parser (구문 분석)
- 토큰 스트림을 추상 구문 트리(AST)로 변환
- 문법 규칙 검증
- 오류 복구

#### Type Checker (타입 검증)
- Intent 기반 타입 추론
- 타입 일치성 검증
- 안전성 검증 (null, bounds, ...)

#### IR Generator (중간 표현)
- AST를 SSA 기반 IR로 변환
- 최적화 기회 식별
- 메타데이터 추가

#### Code Generator (코드 생성)
- IR을 최종 코드로 변환
- 타겟별 최적화 (TypeScript, C, VM bytecode)

---

## 7. 예제 모음

### 7.1 기본 예제

#### 예제 1: 합계 계산

```freelang
fn sum
  intent: "배열의 모든 요소를 더하기"
  input: numbers: array<number>
  output: number
  do
    result = 0
    for num in numbers
      result = result + num
    return result

// 사용
sum([1, 2, 3, 4, 5])  // 15
```

#### 예제 2: 최대값 찾기

```freelang
fn max
  intent: "배열의 최대값 찾기"
  input: numbers: array<number>
  output: Option<number>
  do
    if numbers.len() == 0
      None
    else
      mut largest = numbers[0]
      for i in 1..numbers.len()
        if numbers[i] > largest
          largest = numbers[i]
      Some(largest)
```

#### 예제 3: 필터링

```freelang
fn filter_greater_than
  intent: "threshold보다 큰 값만 필터링"
  input: numbers: array<number>, threshold: number
  output: array<number>
  do
    result = []
    for num in numbers
      if num > threshold
        result.push(num)
    return result

// 사용
filter_greater_than([1, 5, 3, 8, 2], 4)  // [5, 8]
```

### 7.2 중급 예제

#### 예제 4: 패턴 매칭

```freelang
enum HTTPStatus
  Ok
  NotFound
  ServerError

fn handle_status(status: HTTPStatus)
  match status
    HTTPStatus.Ok =>
      print("요청 성공")
    HTTPStatus.NotFound =>
      print("리소스를 찾을 수 없습니다")
    HTTPStatus.ServerError =>
      print("서버 오류 발생")
```

#### 예제 5: 오류 처리

```freelang
fn safe_divide
  intent: "안전한 나눗셈 (0으로 나누기 처리)"
  input: dividend: number, divisor: number
  output: Result<number, string>
  do
    if divisor == 0
      Err("Division by zero is not allowed")
    else
      Ok(dividend / divisor)

// 사용
match safe_divide(10, 2)
  Ok(result) =>
    print(f"Result: {result}")    // Result: 5
  Err(error) =>
    print(f"Error: {error}")
```

#### 예제 6: 구조체 사용

```freelang
struct Point
  x: number
  y: number

fn distance_from_origin
  intent: "원점으로부터의 거리 계산"
  input: point: Point
  output: number
  do
    sqrt(point.x * point.x + point.y * point.y)

// 사용
point = Point { x: 3, y: 4 }
distance_from_origin(point)  // 5
```

#### 예제 7: 제너릭 함수

```freelang
fn first<T>
  intent: "배열의 첫 번째 요소 반환"
  input: array: array<T>
  output: Option<T>
  do
    if array.len() > 0
      Some(array[0])
    else
      None

// 사용
first([1, 2, 3])           // Some(1)
first(["a", "b"])          // Some("a")
first([])                  // None
```

### 7.3 실전 예제

#### 예제 8: 은행 계좌 (Banking)

```freelang
struct Account
  id: string
  balance: number
  transactions: array<Transaction>

struct Transaction
  amount: number
  timestamp: number
  type: string  // "deposit", "withdrawal"

fn deposit
  intent: "계좌에 금액을 입금"
  input: account: Account, amount: number
  output: Result<Account, string>
  do
    if amount <= 0
      Err("Amount must be positive")
    else
      account.balance = account.balance + amount
      transaction = Transaction {
        amount: amount,
        timestamp: now(),
        type: "deposit"
      }
      account.transactions.push(transaction)
      Ok(account)

fn withdraw
  intent: "계좌에서 금액을 출금"
  input: account: Account, amount: number
  output: Result<Account, string>
  do
    if amount <= 0
      Err("Amount must be positive")
    else if account.balance < amount
      Err("Insufficient balance")
    else
      account.balance = account.balance - amount
      transaction = Transaction {
        amount: amount,
        timestamp: now(),
        type: "withdrawal"
      }
      account.transactions.push(transaction)
      Ok(account)
```

---

## 8. 표준 라이브러리

### 8.1 기본 함수

| 함수 | 입력 | 출력 | 설명 |
|------|------|------|------|
| `print(value)` | any | () | 값 출력 |
| `println(value)` | any | () | 값 출력 후 줄바꿈 |
| `len(collection)` | array\|string\|map | number | 길이 반환 |
| `push(array, item)` | array<T>, T | () | 배열에 항목 추가 |
| `pop(array)` | array<T> | Option<T> | 배열에서 항목 제거 |
| `contains(collection, item)` | array\|map, any | bool | 포함 여부 확인 |

### 8.2 문자열 함수

| 함수 | 설명 |
|------|------|
| `len(str)` | 문자열 길이 |
| `substring(str, start, end)` | 부분 문자열 |
| `concat(str1, str2)` | 문자열 연결 |
| `split(str, separator)` | 문자열 분할 |
| `trim(str)` | 공백 제거 |
| `to_lowercase(str)` | 소문자 변환 |
| `to_uppercase(str)` | 대문자 변환 |

### 8.3 배열 함수

| 함수 | 설명 |
|------|------|
| `map(array, fn)` | 함수 적용하여 변환 |
| `filter(array, fn)` | 조건에 맞는 요소만 필터 |
| `reduce(array, fn, initial)` | 배열을 단일 값으로 축약 |
| `sort(array)` | 배열 정렬 |
| `reverse(array)` | 배열 역순 |

### 8.4 수학 함수

| 함수 | 입력 | 출력 |
|------|------|------|
| `abs(n)` | number | number |
| `sqrt(n)` | number | number |
| `pow(n, p)` | number, number | number |
| `sin(n)` | number | number |
| `cos(n)` | number | number |
| `floor(n)` | number | number |
| `ceil(n)` | number | number |

---

## 9. Phase 1 기능

### 9.1 포함된 기능 (Phase 1)

✅ **완벽하게 지원**:
- 변수 선언 (let, mut)
- 기본 타입 (number, string, bool)
- 함수 선언 및 호출
- 제어문 (if/else, match, for, while)
- 배열과 맵
- 구조체
- Enum과 패턴 매칭
- 기본 트레이트
- Result 기반 오류 처리
- 제너릭 (기본)

### 9.2 제외된 기능 (Phase 2+)

❌ **Phase 2+에서 추가**:
- async/await (비동기 프로그래밍)
- 스레드 및 동시성
- 매크로 시스템
- Reflection (동적 타입 정보)
- FFI (외부 함수 인터페이스)
- 고급 제너릭 (Specialization)
- Lifetime (명시적 생명주기)

---

## 10. 호환성 및 마이그레이션

### 10.1 v2-freelang-ai와의 호환성

| 기능 | v2-freelang-ai | FreeLang v1.0 | 상태 |
|------|-----------------|---------------|----|
| 함수 선언 | ✅ | ✅ | 호환 |
| Intent | ✅ | ✅ | 호환 |
| 기본 타입 | ✅ | ✅ | 호환 |
| 패턴 매칭 | ✅ | ✅ | 호환 |
| 제너릭 | ⚠️ (부분) | ✅ | 향상 |
| 에러 처리 | ✅ | ✅ | 호환 |

### 10.2 Breaking Changes

**v2-freelang-ai에서 마이그레이션할 때 주의**:

1. **타입 명시**: 일부 타입 추론이 더 엄격해짐
2. **Null 안전**: `null`이 명시적으로 처리되어야 함
3. **불변성**: 기본값이 불변으로 변경됨

### 10.3 마이그레이션 가이드

#### Step 1: 함수 시그니처 확인

**Before (v2)**:
```freelang
fn sum input array output number intent "합계"
```

**After (v1.0)**:
```freelang
fn sum
  intent: "합계"
  input: array: array<number>
  output: number
```

#### Step 2: 타입 명시

모든 입력과 출력에 타입을 명시:

```freelang
// 기존 (부분 생략 가능)
fn process input intent "처리"

// v1.0 (완전히 명시)
fn process
  intent: "처리"
  input: data: array<number>
  output: array<number>
```

#### Step 3: 오류 처리

Result 타입 사용:

```freelang
// Before
fn divide a b
  if b == 0: return None
  return a / b

// After
fn divide
  input: a: number, b: number
  output: Result<number, string>
  do
    if b == 0
      Err("Cannot divide by zero")
    else
      Ok(a / b)
```

---

## 부록

### A. EBNF 문법

```ebnf
program          = { item }
item             = function_decl | struct_decl | enum_decl | trait_decl | impl_decl

function_decl    = "fn" identifier signature ["do" block]
signature        = ["intent" ":" string]
                   ["input" ":" parameters]
                   ["output" ":" type]

parameters       = parameter { "," parameter }
parameter        = identifier [":" type]

struct_decl      = "struct" identifier "{" { field } "}"
field            = identifier ":" type

enum_decl        = "enum" identifier "{" { variant } "}"
variant          = identifier ["(" types ")"]

trait_decl       = "trait" identifier "{" { function_decl } "}"
impl_decl        = "impl" identifier "for" identifier "{" { function_decl } "}"

type             = base_type [ "<" type_params ">" ]
base_type        = "number" | "string" | "bool" | "(" ")"
type_params      = type { "," type }

block            = { statement }
statement        = expr_stmt | if_stmt | match_stmt | loop_stmt | return_stmt
expr_stmt        = expression [";"]
if_stmt          = "if" expression block ["else" block]
match_stmt       = "match" expression "{" { pattern "=>" block } "}"
loop_stmt        = "for" identifier "in" expression block
                 | "while" expression block
return_stmt      = "return" expression

expression       = assignment
assignment       = or_expr [ "=" expression ]
or_expr          = and_expr { "||" and_expr }
and_expr         = equality { "&&" equality }
equality         = comparison { ("==" | "!=") comparison }
comparison       = addition { ("<" | ">" | "<=" | ">=") addition }
addition         = multiplication { ("+" | "-") multiplication }
multiplication   = unary { ("*" | "/") unary }
unary            = ("!" | "-") unary | postfix
postfix          = primary { "." identifier | "[" expression "]" | "(" args ")" }
primary          = literal | identifier | "(" expression ")"

literal          = number | string | bool | array_lit | map_lit
array_lit        = "[" { expression "," } "]"
map_lit          = "{" { identifier ":" expression "," } "}"
```

### B. 키워드 목록

```
fn        - 함수 선언
let       - 불변 변수
mut       - 가변 변수
struct    - 구조체
enum      - 열거형
trait     - 트레이트
impl      - 구현
if        - 조건
else      - 다른 조건
match     - 패턴 매칭
for       - for 루프
while     - while 루프
return    - 반환
intent    - 함수 의도
input     - 입력 파라미터
output    - 출력 타입
do        - 함수 본체 시작
Some      - 옵션 Some 값
None      - 옵션 None 값
Ok        - Result Ok 값
Err       - Result Err 값
true      - 참
false     - 거짓
```

### C. 연산자 목록

| 연산자 | 이름 | 예시 |
|--------|------|------|
| `+` | 덧셈 | `a + b` |
| `-` | 뺄셈 | `a - b` |
| `*` | 곱셈 | `a * b` |
| `/` | 나눗셈 | `a / b` |
| `%` | 나머지 | `a % b` |
| `==` | 같음 | `a == b` |
| `!=` | 다름 | `a != b` |
| `<` | 작음 | `a < b` |
| `>` | 큼 | `a > b` |
| `<=` | 작거나 같음 | `a <= b` |
| `>=` | 크거나 같음 | `a >= b` |
| `&&` | 논리 AND | `a && b` |
| `\|\|` | 논리 OR | `a \|\| b` |
| `!` | 논리 NOT | `!a` |
| `=` | 할당 | `a = b` |

---

## 최종 선언

**FreeLang v1.0은 다음을 선언합니다**:

1. ✅ **AI 친화적** - AI가 읽고 생성하기 쉬운 명확한 문법
2. ✅ **안전성** - 컴파일 타입 안전 + 메모리 안전
3. ✅ **간결성** - 필수 요소만 작성, 나머지는 자동 추론
4. ✅ **확장성** - 트레이트와 제너릭으로 확장 가능
5. ✅ **실용성** - 실제 프로젝트에 사용 가능

**라이센스**: Apache 2.0
**저장소**: https://github.com/freelang/freelang (예정)
**버전 공개**: 2026-03-02

---

**작성자**: Claude (AI)
**검토자**: FreeLang 개발 팀
**상태**: 🟢 공식 명세서
**라이센스**: CC-BY-SA 4.0
