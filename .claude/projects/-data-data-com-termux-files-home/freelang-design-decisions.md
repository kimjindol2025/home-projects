# FreeLang 핵심 설계 결정 분석

**언어 설계자 관점에서 본 각 버전의 기로와 선택**

---

## 1️⃣ 메모리 모델: "GC vs Move" 기로

### v1-v3: 혼란 (실패의 원인)

```
의사 결정:
┌─────────────────────────────────────┐
│  "메모리를 어떻게 관리할 것인가?"    │
├─────────────────────────────────────┤
│  A) Garbage Collection (Python 방식) │
│     장점: 구현 간단, 사용자 부담 적음 │
│     단점: 성능, 예측 불가능, GC pause │
│                                      │
│  B) Manual Management (C 방식)       │
│     장점: 성능, 완벽 제어           │
│     단점: 복잡, 메모리 누수 위험     │
│                                      │
│  C) Move Semantics (Rust 방식)       │
│     장점: 안전, 성능, 예측 가능     │
│     단점: 학습곡선 가파름           │
│                                      │
│  D) Reference Counting (Swift 방식)  │
│     장점: 비교적 간단               │
│     단점: 순환 참조 문제            │
│                                      │
│  E) Hybrid (GC + Manual)             │
│     장점: 유연함                    │
│     단점: 복잡함                    │
└─────────────────────────────────────┘

v1-v3: 명확한 결정 없이 구현 시도
→ Checker에서 Move 규칙 불명확
→ Type checker 복잡도 폭발 (2500+ LOC)
→ 실패
```

### v4: "Move Semantics 선택"

```
결정: C) Move Semantics 채택

근거:
1. 타겟 사용자: AI 에이전트
   → 생성된 코드도 안전해야 함
   → GC의 비결정성 불가
   → 따라서 Rust 방식 필요

2. "AI가 생성한 코드가 컴파일을 통과하면, 그 코드는 안전하다"
   → Compile-time 검증 필수
   → Move semantics가 가장 적합

3. 성능
   → JIT 없이도 성능 예측 가능
   → GC pause 없음

설계:
  ┌─────────────────────────────────┐
  │ Copy 타입 (스택)                │
  │ - i32, i64, f64, bool          │
  │ - 암묵적 복사                   │
  │                                 │
  │ Move 타입 (힙)                  │
  │ - string, [T], {fields}        │
  │ - 명시적 소유권 이전            │
  │                                 │
  │ 메모리 해제                     │
  │ - Scope 벗어날 때 자동 drop     │
  │ - C++의 RAII 패턴               │
  └─────────────────────────────────┘

구현:
  Type Checker → Move 추적 (801줄로 명확화)
  VM → Value ownership 추적

테스트:
  ✓ use-after-move: Compile error
  ✓ double-move: Compile error
  ✓ partial-move: Compile error (array slice 등)
```

---

## 2️⃣ 제어 흐름: "for...in vs while" 기로

### v4: "for...in만" 선택

```
의사 결정:
┌────────────────────────────────────┐
│  "루프를 어떻게 표현할 것인가?"     │
├────────────────────────────────────┤
│  A) C/Java 방식: for, while, do-while
│  B) Python 방식: for..in만
│  C) Functional: map, filter, fold만
│  D) Go 방식: for만 (유연함)
│  E) Hybrid: for..in + while
└────────────────────────────────────┘

v4 선택: B) for...in만

근거:
1. 무한 루프 방지
   - 프로그래머 의도 명확
   - while (true) 같은 실수 방지

2. 반복자 명시
   - for i in [1,2,3] → i는 어디서 오는가?
   - 명확한 범위

3. AI 친화성
   - 무한 루프 버그 감소
   - 생성된 코드의 안전성

한계:
  ❌ 무한 루프 필요시 어려움
     → range(0, INF)? 불가능
     → 작업 루프 (job queue) 패턴 강제

설계:
  for <var> in <iterable> {
    ...
    if <cond> { return / break }  // v4에는 break 없음!
  }
```

### v5: "for...in + while + loop" 추가

```
기로:
┌────────────────────────────────────┐
│  "loop를 어떻게 할 것인가?"         │
├────────────────────────────────────┤
│  A) 그대로 for...in만 유지          │
│     → AI 친화적이나 표현력 부족    │
│                                    │
│  B) while/loop 추가                │
│     → 표현력 증가하지만 무한루프 위험
│                                    │
│  C) Conditional break 추가          │
│     → for...in 유지하면서 유연성  │
│                                    │
│  선택: B+C) while + break + continue
└────────────────────────────────────┘

v5 추가:

// While loop
while <condition> {
  ...
  if <cond> { break }
  if <cond> { continue }
}

// Infinite loop
loop {
  ...
  if <cond> { break }
}

영향:
  ✓ 표현력 증가
  ✓ Job queue, event loop 패턴 가능
  ❌ 무한 루프 위험 증가
     → 문서화로 보상 (피해야 할 패턴)

패턴:
  // 권장
  while condition {
    do_work()
  }

  // 권장하지 않음 (비상용)
  loop {
    if should_exit { break }
  }
```

### v6: 최적화 & 정규화

```
v5의 문제점:
  - while, for...in, loop이 섞여 있음
  - 언제 뭘 써야 하는가?

v6 가이드:

┌─────────────────────────────────────────┐
│ 1. 범위를 알 때 → for...in              │
│    for i in [1, 2, 3] { ... }          │
│    for i in range(0, 10) { ... }       │
│    권장도: ⭐⭐⭐ (명확)               │
│                                        │
│ 2. 조건이 있을 때 → while              │
│    while x > 0 { x = x - 1 }           │
│    권장도: ⭐⭐ (조심스러움)           │
│                                        │
│ 3. 사건 기반 → loop + break            │
│    loop {                              │
│      let event = wait_event()          │
│      if event == "quit" { break }      │
│    }                                   │
│    권장도: ⭐ (극히 제한적)            │
└─────────────────────────────────────────┘
```

---

## 3️⃣ 타입 시스템: "동적 vs 정적" 기로

### v4: "정적 타입 채택"

```
의사 결정:
┌────────────────────────────────────┐
│  "타입 검증을 언제 할 것인가?"     │
├────────────────────────────────────┤
│  A) Runtime (Python)               │
│     - 유연함                       │
│     - 에러는 실행 시 발생          │
│                                    │
│  B) Compile-time (Java, Rust)      │
│     - 안전함                       │
│     - 에러는 컴파일 시 발견        │
│                                    │
│  C) Gradual (TypeScript)           │
│     - 유연함 + 안전함              │
│     - 복잡함                       │
└────────────────────────────────────┘

v4 선택: B) 정적 타입 + 타입 추론

설계:
  // 함수 시그니처는 명시 필수
  fn add(x: i32, y: i32): i32 { x + y }

  // 로컬 변수는 타입 추론
  let x = 42  // x: i32 자동 추론

근거:
1. AI 친화성
   → "이 코드는 안전한가?"를 컴파일 타임에 검증
   → Runtime error 최소화

2. 성능 최적화
   → 컴파일 타임에 타입 결정
   → 런타임 타입 체크 불필요

3. 명확성
   → 함수 인터페이스 명확
   → 타입 에러는 컴파일 시 발견

제약:
  ✗ Generic 초기 버전에서 미지원
  ✗ Any 타입 없음
  ✗ interface/protocol 미지원
```

### v5: "Generic + Trait 추가"

```
기로:
┌────────────────────────────────────┐
│  "코드 재사용을 어떻게 할 것인가?" │
├────────────────────────────────────┤
│  A) 복사 (C 방식)                   │
│     fn add_i32(x, y) { ... }        │
│     fn add_f64(x, y) { ... }        │
│                                    │
│  B) Generic (Rust, C++ 방식)       │
│     fn add<T>(x: T, y: T): T       │
│                                    │
│  C) Duck typing (Python)            │
│     fn add(x, y): return x + y    │
│                                    │
│  선택: B) Generic + Trait
└────────────────────────────────────┘

v5 추가:

// Generic 함수
fn first<T>(arr: [T]): Option<T> {
  if length(arr) > 0 { Some(arr[0]) }
  else { None }
}

// Generic struct
struct Pair<T> {
  left: T,
  right: T,
}

// Trait (기초)
trait Comparable {
  fn compare(self, other: Self): i32
}

impl Comparable for i32 {
  fn compare(self, other: i32): i32 {
    if self < other { -1 } else { 1 }
  }
}

영향:
  ✓ 코드 재사용성 증가
  ✓ Type safety 유지
  ❌ 컴파일 시간 증가 (monomorphization)
```

---

## 4️⃣ 모듈 시스템: "단일 파일 vs 모듈" 기로

### v4: "단일 파일"

```
철학: "Core language는 간단하게"

특징:
  - 모든 코드가 main.fl 같은 파일에
  - import/export 없음
  - 전역 네임스페이스

장점:
  ✓ 간단함 (명세 줄임)
  ✓ 구현 단순
  ✓ 학습 곡선 낮음

단점:
  ❌ 코드 재사용 어려움
  ❌ 대규모 프로젝트 불가능
  ❌ 네임스페이스 충돌 가능

v4의 입장:
  "모듈은 v5의 '살'이다"
  "v4는 '핵'만 증명하면 된다"
```

### v5: "파일 기반 모듈"

```
기로:
┌────────────────────────────────────┐
│  "모듈을 어떻게 할 것인가?"        │
├────────────────────────────────────┤
│  A) 파일 = 모듈 (Go 방식)          │
│     각 파일이 자동으로 모듈        │
│                                    │
│  B) 명시적 모듈 선언 (Python)      │
│     module { ... }                │
│                                    │
│  C) 패키지 중심 (Java)             │
│     package com.example.foo        │
│                                    │
│  선택: A) 파일 = 모듈
└────────────────────────────────────┘

v5 설계:

// math.fl
export fn sum(arr: [i32]): i32 { ... }

// main.fl
import { sum } from "./math.fl"
import { to_upper } from "std:string"

let result = sum([1, 2, 3])

특징:
  ✓ 간단함 (파일 = 모듈)
  ✓ Go처럼 직관적
  ✓ 순환 import 방지 가능 (DAG 체크)

구현 복잡도:
  - Lexer: +3% (import/export 토큰)
  - Parser: +5% (import 문법)
  - Checker: +15% (모듈 해석, 심볼 해석)
  - Compiler: +10% (모듈 코드 생성)
```

---

## 5️⃣ 함수형 vs 객체지향: "패러다임" 기로

### v4: "함수형 기초"

```
철학: "순수 함수와 데이터 변환 중심"

특징:
  - 함수는 1급 객체 아님
  - 변수는 mutable
  - 메서드 없음

코드 예:
  fn process(arr: [i32]): i32 {
    var sum = 0
    for x in arr {
      sum = sum + x
    }
    return sum
  }

한계:
  ❌ 캡슐화 어려움
  ❌ OOP 패턴 구현 번거로움
```

### v5: "OOP 지원 추가"

```
기로:
┌────────────────────────────────────┐
│  "OOP를 어느 정도로 할 것인가?"     │
├────────────────────────────────────┤
│  A) 없음 (함수형만)                │
│     → 표현력 부족                  │
│                                    │
│  B) 경량 OOP (struct + method)     │
│     → 캡슐화 + 간단함              │
│                                    │
│  C) 전체 OOP (상속, 다형성)        │
│     → 강력하지만 복잡              │
│                                    │
│  선택: B) 경량 OOP (struct + impl)
└────────────────────────────────────┘

v5 추가:

struct Point {
  x: i32,
  y: i32,
}

impl Point {
  fn distance(self): f64 {
    sqrt(self.x * self.x + self.y * self.y)
  }

  fn translate(mut self, dx: i32, dy: i32) {
    self.x = self.x + dx
    self.y = self.y + dy
  }
}

특징:
  ✓ 데이터와 메서드 함께 정의
  ✓ self / &self / &mut self (Rust 방식)
  ✓ 상속 없음 (Composition 강조)

비철학:
  - "상속은 조합보다 나쁘다" (FP 관점)
  - "OOP는 도구일 뿐" (균형잡힌 관점)
```

### v6: "고급 OOP 미지원, 대신 데코레이터"

```
v5의 OOP는 충분한가?

장점:
  ✓ Struct + impl로 기본 캡슐화 가능
  ✓ 메서드 호출 가능
  ✓ 간단함

부족한 점:
  ❌ 상속 없음 → 코드 중복
  ❌ 다형성 없음 → type assertion 필요
  ❌ 인터페이스 없음 → duck typing 필요

v6의 선택: "OOP 더하기 대신 메타프로그래밍"

// @derive로 자동 구현 (메타프로그래밍)
@derive(Debug, Eq, Hash)
struct Point {
  x: i32,
  y: i32,
}

// @inline, @memoize 등 최적화 지시자
@inline
fn hot_path() { ... }

철학:
  "OOP보다는 composition + 메타프로그래밍"
  (Python, Go, Rust와 비슷한 방향)
```

---

## 6️⃣ 에러 처리: "Exception vs Result" 기로

### v4: "Result<T,E>만"

```
의사 결정:
┌────────────────────────────────────┐
│  "에러를 어떻게 표현할 것인가?"     │
├────────────────────────────────────┤
│  A) Exception (Java)               │
│     try-catch-finally              │
│                                    │
│  B) Result/Option (Rust)           │
│     match Ok/Err/Some/None         │
│                                    │
│  C) Error code (C)                 │
│     if err { ... }                 │
│                                    │
│  선택: B) Result<T,E> + Option<T>
└────────────────────────────────────┘

v4 설계:

// 성공/실패 명시
fn divide(a: i32, b: i32): Result<i32, string> {
  if b == 0 {
    Err("Division by zero")
  } else {
    Ok(a / b)
  }
}

// 처리 강제
match divide(10, 2) {
  Ok(result) => println(str(result)),
  Err(msg) => println(msg),
}

// ?연산자로 간결하게
fn safe_divide(a: i32, b: i32): Result<i32, string> {
  let x = divide(a, b)?
  Ok(x * 2)
}

장점:
  ✓ 에러를 값으로 표현 (명시적)
  ✓ Control flow 안전 (컴파일 시 강제)
  ✓ Stack unwinding 없음 (성능)

단점:
  ❌ Try-catch 없음 (복잡해 보임)
  ❌ 학습 곡선 (Rust처럼)
```

### v5: "try-catch 추가"

```
기로:
┌────────────────────────────────────┐
│  "try-catch를 추가할 것인가?"       │
├────────────────────────────────────┤
│  A) Result만 유지 (Rust 방식)      │
│                                    │
│  B) try-catch 추가 (예외 처리 지원) │
│     Result와 함께 사용 가능         │
│                                    │
│  선택: B) Result + try-catch 병행
└────────────────────────────────────┘

v5 추가:

try {
  let result = risky_operation()?
  println("Success: " + str(result))
} catch err {
  println("Error: " + err.message)
} finally {
  cleanup()
}

그리고 action으로 선언형 에러 처리:

action safe_operation {
  input: string,
  output: Result<string, Error>,

  on_validate {
    if length(input) == 0 {
      return Err("Empty input")
    }
  }

  on_execute {
    return Ok(process(input))
  }

  on_failure {
    log("Operation failed")
  }
}

철학:
  "Result는 데이터 에러 (비즈니스 로직)"
  "try-catch는 예외 (프로그래밍 에러)"
  "action은 워크플로우 (선언형)"
```

---

## 7️⃣ Observability: "시스템 이해" 기로

### v4-v5: 부재

```
v4-v5는 관찰성이 없다:
  ❌ 함수 호출 추적 없음
  ❌ 성능 메트릭 없음
  ❌ 에러율 통계 없음
```

### v6: "@observe 데코레이터"

```
기로:
┌────────────────────────────────────┐
│  "시스템을 어떻게 이해할 것인가?"   │
├────────────────────────────────────┤
│  A) 완전 자동 (APM처럼)            │
│     모든 함수 자동 추적             │
│     → 오버헤드 큼                  │
│                                    │
│  B) 명시적 (데코레이터)            │
│     선택적으로 표시                │
│     → 오버헤드 작음                │
│                                    │
│  C) 수동 (명시적 로깅)             │
│     개발자가 직접 추적             │
│     → 번거로움                     │
│                                    │
│  선택: B) 데코레이터 기반 (선택적)
└────────────────────────────────────┘

v6 설계:

@observe(metrics: ["duration", "error_rate", "memory"])
fn expensive_operation(x: i32): Result<i32, string> {
  // 자동으로 추적됨:
  // 1. 실행 시간 (duration)
  // 2. 에러율 (error_rate)
  // 3. 메모리 사용 (memory)

  return Ok(x * 2)
}

// 메트릭 조회
let metrics = metrics::get_all(expensive_operation)
println("Duration: " + str(metrics.duration))
println("Errors: " + str(metrics.error_rate * 100) + "%")

특징:
  ✓ 선택적 (필요한 함수만)
  ✓ 오버헤드 최소
  ✓ 메트릭 자동 수집
  ✓ "설명 가능한 AI" 지원

메트릭:
  - duration: 함수 실행 시간 (ms)
  - error_rate: 에러율 (0.0~1.0)
  - memory: 메모리 증가량 (bytes)
  - call_count: 호출 횟수
  - avg_duration: 평균 실행 시간
```

---

## 8️⃣ 성능 최적화: "해석 vs JIT" 기로

### v4: "Stack VM (해석)"

```
설계:
  source.fl → Bytecode → Stack VM → Output

오버헤드:
  - Opcode 디스패치 (각 명령마다 switch)
  - 스택 조작 (push/pop)
  - 런타임 타입 체크 (아직은 minimal)

성능:
  ~100k ops/sec (추정)

철학:
  "정확성이 성능보다 중요"
  "나중에 JIT 추가 가능"
```

### v6: "Bytecode 최적화"

```
기로:
┌────────────────────────────────────┐
│  "어디서 성능을 얻을 것인가?"       │
├────────────────────────────────────┤
│  A) JIT 컴파일 (→ v7)             │
│     → 성능 10배 증가 (목표)        │
│     → 구현 복잡                    │
│                                    │
│  B) 바이트코드 최적화               │
│     → 성능 2-3배 증가              │
│     → 구현 간단                    │
│                                    │
│  C) Opcode 축소                    │
│     → 성능 20-30% 증가             │
│     → 코드 명확성 개선             │
│                                    │
│  선택: B+C) 바이트코드 + Opcode 축소
└────────────────────────────────────┘

v6 최적화:

1. Opcode 축소 (45 → 29개)
   // v4
   ADD_I32, ADD_F64 (분산)

   // v6
   ARITH { op: "+", type: "i32" }
   → 런타임 타입 체크로 보상

2. 상수 폴딩 (Compile time)
   // v4
   PUSH_I32(2)
   PUSH_I32(3)
   ADD_I32

   // v6
   PUSH_I32(5)  // 컴파일 시 계산

3. Dead code elimination
   fn unused_func() { ... }  // 컴파일 결과에서 제거

성능:
  ~150k ops/sec (추정, 50% 향상)

미래 (v7):
  JIT 컴파일 (hot path 감지)
  → ~1M ops/sec (목표)
```

---

## 종합: 설계 철학의 진화

```
v1-v3: 혼란
  ❓ 메모리 모델이 불명확
  ❓ 타입 시스템이 느슨함
  ❓ 설계 명세가 부족
  → 실패

v4: "기초 원칙"
  1. Safety First (Move Semantics)
     → 컴파일 타임에 메모리 안전 보증

  2. Clarity (for...in 제한, 명시적 타입)
     → 프로그래머 의도 명확화

  3. Provable (SPEC 기반 설계)
     → 형식 명세부터 시작

  철학: "AI가 생성한 코드도 안전하게"

v5: "기능 추가"
  1. Modularity (파일 기반 모듈)
  2. Object-Oriented (struct + impl)
  3. Functional (클로저, 고차 함수)
  4. Action System (선언형 워크플로우)

  철학: "표현력 증가, 안전성 유지"

v6: "완성과 최적화"
  1. Stability (버그 수정, 테스트 강화)
  2. Stdlib (HTTP, regex, DB 내장)
  3. Performance (Opcode 최적화)
  4. Observability (메트릭 추적)

  철학: "실무 사용 가능한 수준"
```

---

## 🎓 언어 설계자가 배운 교훈

### 1. "설계 우선"

```
❌ v1-v3: 구현부터 시작
  → 타입 시스템 혼란
  → Type checker 복잡도 폭발
  → 실패

✅ v4: SPEC부터 시작
  → 18개 명세 문서 (9,136 LOC)
  → 각 단계를 명확히 함
  → 성공
```

### 2. "점진적 확장"

```
v4: 기초만 (6,934 LOC)
  - Lexer, Parser, TypeChecker, Compiler, VM
  - 50개 내장 함수
  - 334 테스트

v5: 기능 추가 (7,000+ LOC)
  - Module System
  - Struct + impl
  - Closures
  - Action System

v6: 최적화 (4,515 LOC)
  - 코드 35% 감소
  - 기능 60% 증가
  - Stdlib 11개 모듈

특징: 각 버전이 이전 호환성 유지
```

### 3. "안전성 vs 표현력 균형"

```
Axis:
      표현력
        ↑
        │  Python
        │    (고)
        │
        │  v5-v6
        │  (중)
        │
        │  v4
        │  (저)
        │
        └──────────────→ 안전성 (높)
        Rust          v4    v5  Python
                            (저)
```

---

**결론: FreeLang은 "AI-first, safety-preserving, gradually-extending" 철학으로 설계되었다.**
