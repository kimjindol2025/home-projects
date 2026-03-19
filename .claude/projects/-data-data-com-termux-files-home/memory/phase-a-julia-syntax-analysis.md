---
name: Phase A - Julia 문법 상세 분석
description: Julia 핵심 문법을 FreeLang에서 구현 가능한 형태로 분석
type: project
---

# 📋 Phase A: Julia 문법 상세 분석

**완료 날짜**: 2026-03-19
**상태**: 🟢 Task A.1 진행 중
**목표**: Julia 기능 → FreeLang 매핑 & 구현 가능성 판단

---

## 1️⃣ **Julia 핵심 문법 분석**

### A. 기본 타입 & 변수

#### Julia 문법
```julia
# 정수
x::Int = 5
y::Int32 = 10

# 부동소수점
a::Float64 = 3.14
b::Float32 = 2.71f0

# 문자열
s::String = "hello"

# 불린
flag::Bool = true

# Nullable (Union)
maybe::Union{Int, Nothing} = nothing
maybe = 5
```

#### FreeLang 매핑
```freeling
// 정수
let x: Int = 5
let y: Int32 = 10

// 부동소수점
let a: Float64 = 3.14
let b: Float32 = 2.71

// 문자열
let s: String = "hello"

// 불린
let flag: Bool = true

// Optional (Union)
let maybe: Option[Int] = None
let maybe = Some(5)
```

**호환성**: ✅ 100% (FreeLang이 더 엄격함)

---

### B. 배열 & 인덱싱

#### Julia 문법
```julia
# 배열 생성
arr = [1, 2, 3, 4, 5]
arr = Int64[1, 2, 3]  # 타입 명시

# 인덱싱 (1-based)
arr[1]      # 첫 원소
arr[end]    # 마지막 원소
arr[1:3]    # 슬라이싱 [1, 2, 3]
arr[1:2:5]  # 스텝 [1, 3, 5]

# 다차원 배열
matrix = [1 2 3; 4 5 6; 7 8 9]
matrix[1, 2]     # 행 1, 열 2
matrix[2, :]     # 행 2 전체

# 배열 생성식
squares = [x^2 for x in 1:10]
evens = [x for x in 1:10 if x % 2 == 0]
```

#### FreeLang 구현
```freeling
// 배열 생성
let arr = [1, 2, 3, 4, 5]
let arr = [1, 2, 3]: Array[Int]

// 인덱싱 (0-based로 변환 필요)
arr[0]               // 첫 원소
arr[len(arr) - 1]   // 마지막 원소
arr.slice(0, 3)      // 슬라이싱 [1, 2, 3]
arr.stride(0, 5, 2)  // 스텝 [1, 3, 5]

// 다차원 배열
let matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
matrix[0][1]         // [0, 1]
matrix[1]            // [4, 5, 6] (행)

// 배열 생성식
let squares = [x^2 | x <- range(1, 10)]
let evens = [x | x <- range(1, 10), x % 2 == 0]
```

**호환성**: 🟡 80% (인덱싱 1-based → 0-based 변환 필요)

**변환 전략**:
```
Julia: arr[n]
→ FreeLang: arr[n-1] (자동 변환)

Julia: arr[1:3]
→ FreeLang: arr.slice(0, 3)

Julia: [x for x in 1:10]
→ FreeLang: [x | x <- range(1, 10)]
```

---

### C. 함수 & 다중 디스패치

#### Julia 문법
```julia
# 기본 함수
function add(x, y)
    return x + y
end

# 타입 명시
function add(x::Int, y::Int)::Int
    return x + y
end

# 단일줄 함수
square(x) = x^2

# 다중 디스패치 (Julia 핵심!)
function process(x::Int)
    return x + 1
end

function process(x::Float64)
    return x * 2
end

function process(x::String)
    return uppercase(x)
end

# 호출
process(5)          # Int 버전: 6
process(3.14)       # Float64 버전: 6.28
process("hello")    # String 버전: "HELLO"
```

#### FreeLang 구현
```freeling
// 기본 함수
function add(x, y) = x + y

// 타입 명시
function add(x: Int, y: Int): Int = x + y

// 다중 디스패치 (신규!)
// 방법 1: 명시적 dispatch table
let process_dispatch = {
  "Int" -> fn(x: Int) = x + 1,
  "Float" -> fn(x: Float) = x * 2,
  "String" -> fn(x: String) = uppercase(x)
}

function process(x: Dynamic) = {
  let type_name = typeof(x)
  let handler = lookup(process_dispatch, type_name)
  handler(x)
}

// 방법 2: pattern matching
function process(x: Dynamic) = match x {
  Int(v) -> v + 1,
  Float(v) -> v * 2,
  String(v) -> uppercase(v)
}

// 호출
process(5)          // 6
process(3.14)       // 6.28
process("hello")    // "HELLO"
```

**호환성**: 🟡 70% (다중 디스패치 구현 필요)

**신규 기능**:
- `Dynamic` 타입 추가
- `dispatch` 메커니즘 구현
- Pattern matching with types

---

### D. 연산자 오버로딩

#### Julia 문법
```julia
# 벡터 덧셈
v1 = [1, 2, 3]
v2 = [4, 5, 6]
v3 = v1 + v2  # [5, 7, 9]

# Broadcasting (점 연산자)
a = [1, 2, 3]
b = a .+ 5    # [6, 7, 8]
c = a .* 2    # [2, 4, 6]
d = sin.(a)   # [sin(1), sin(2), sin(3)]

# 사용자 정의 타입
struct Point
    x::Float64
    y::Float64
end

import Base.+
(p1::Point) + (p2::Point) = Point(p1.x + p2.x, p1.y + p2.y)

p1 = Point(1, 2)
p2 = Point(3, 4)
p3 = p1 + p2  # Point(4, 6)
```

#### FreeLang 구현
```freeling
// 벡터 덧셈
let v1 = [1, 2, 3]
let v2 = [4, 5, 6]
let v3 = zipWith(add, v1, v2)  // [5, 7, 9]

// Broadcasting (map 사용)
let a = [1, 2, 3]
let b = map(add(_, 5), a)     // [6, 7, 8]
let c = map(mul(_, 2), a)     // [2, 4, 6]
let d = map(sin, a)            // [sin(1), sin(2), sin(3)]

// 사용자 정의 타입
record Point {
  x: Float,
  y: Float
}

// 인스턴스 메서드
impl Point {
  function add(other: Point): Point = {
    Point { x: this.x + other.x, y: this.y + other.y }
  }
}

let p1 = Point { x: 1, y: 2 }
let p2 = Point { x: 3, y: 4 }
let p3 = p1.add(p2)  // Point { x: 4, y: 6 }
```

**호환성**: 🟡 75% (Broadcasting 매핑 필요)

---

### E. 제어 흐름

#### Julia 문법
```julia
# if-else
if x > 0
    println("positive")
elseif x < 0
    println("negative")
else
    println("zero")
end

# for 루프
for i in 1:10
    println(i)
end

# while 루프
while x > 0
    x = x - 1
end

# try-catch (예외처리)
try
    result = risky_operation()
catch e
    println("Error: $e")
finally
    cleanup()
end

# Comprehensions with multiple conditions
result = [x^2 for x in 1:10 if x % 2 == 0 for y in 1:5 if y > 2]
```

#### FreeLang 구현
```freeling
// if-else (호환)
if x > 0 {
    println("positive")
} else if x < 0 {
    println("negative")
} else {
    println("zero")
}

// for 루프
for i <- range(1, 10) {
    println(i)
}

// while 루프
while x > 0 {
    x = x - 1
}

// try-catch (호환)
try {
    let result = risky_operation()
} catch e {
    println(format("Error: {}", e))
} finally {
    cleanup()
}

// Nested comprehensions
let result = [x^2 | x <- range(1, 10), x % 2 == 0, y <- range(1, 5), y > 2]
```

**호환성**: ✅ 95% (문법만 약간 다름)

---

### F. 구조체 & 타입

#### Julia 문법
```julia
# 불변 구조체
struct Point
    x::Float64
    y::Float64
end

# 가변 구조체
mutable struct MutablePoint
    x::Float64
    y::Float64
end

# 추상 타입
abstract type Animal end

struct Dog <: Animal
    name::String
    age::Int
end

struct Cat <: Animal
    name::String
    age::Int
end

# 함수 오버로드
sound(d::Dog) = "Woof!"
sound(c::Cat) = "Meow!"
```

#### FreeLang 구현
```freeling
// 불변 record (기본)
record Point {
    x: Float,
    y: Float
}

// 가변 record
mutable record MutablePoint {
    x: Float,
    y: Float
}

// 프로토콜 (추상 타입 대체)
protocol Animal {
    method name() -> String
    method age() -> Int
    method sound() -> String
}

// 구현
record Dog {
    name: String,
    age: Int
}

impl Dog : Animal {
    function name() = this.name
    function age() = this.age
    function sound() = "Woof!"
}

record Cat {
    name: String,
    age: Int
}

impl Cat : Animal {
    function name() = this.name
    function age() = this.age
    function sound() = "Meow!"
}
```

**호환성**: 🟡 80% (프로토콜 사용, 약간 문법 다름)

---

### G. 모듈 & 패키지

#### Julia 문법
```julia
# 모듈 정의
module MyModule
    export public_func

    public_func(x) = x + 1
    private_func(x) = x - 1
end

# 모듈 사용
using MyModule
result = public_func(5)  # 6

# 패키지 추가
using DataFrames
using Plots

# 조건부 포함
if VERSION >= v"1.6"
    include("new_feature.jl")
end
```

#### FreeLang 구현
```freeling
// 모듈 정의
module MyModule {
    export public_func

    function public_func(x) = x + 1

    function private_func(x) = x - 1
}

// 모듈 사용
import MyModule
let result = MyModule::public_func(5)  // 6

// 라이브러리 import
import stdlib::arrays
import stdlib::collections

// 조건부 포함
if major_version >= 1 && minor_version >= 6 {
    import new_feature
}
```

**호환성**: ✅ 90% (FreeLang 모듈 시스템 사용)

---

## 2️⃣ **호환성 종합 평가**

### 호환성 매트릭스

| 기능 | 호환성 | 난이도 | 우선순위 |
|------|--------|--------|----------|
| 기본 타입 | ✅ 100% | 🟢 낮음 | P0 |
| 배열 & 인덱싱 | 🟡 80% | 🟡 중간 | P0 |
| 함수 | ✅ 100% | 🟢 낮음 | P0 |
| 다중 디스패치 | 🟡 70% | 🔴 높음 | P1 |
| 연산자 오버로딩 | 🟡 75% | 🟡 중간 | P1 |
| 제어 흐름 | ✅ 95% | 🟢 낮음 | P0 |
| 구조체 & 타입 | 🟡 80% | 🟡 중간 | P0 |
| 모듈 & 패키지 | ✅ 90% | 🟢 낮음 | P0 |
| Broadcasting | 🟡 75% | 🟡 중간 | P1 |
| 매크로 | 🟡 60% | 🔴 높음 | P2 |
| 메타프로그래밍 | 🟡 50% | 🔴 높음 | P2 |

**종합 호환성**: 🟡 **78%** (1단계: 기본 기능)

---

## 3️⃣ **Phase A Task 분해**

### Task A.1: Julia 문법 분석 ✅ (완료)
**산출물**:
- julia-syntax-mapping.md (이 문서)
- 호환성 평가 (78%)
- 구현 전략

### Task A.2: FreeLang 타입 시스템 확장
**목표**: Dynamic Type, Union, Type Parameters, Protocols 구현

**파일**: `src/types_extended.fl` (450줄)

**내용**:
```freeling
// 1. Dynamic Type
type Dynamic = Any | Int | Float | String | Bool | List | Map

// 2. Union Types
type Union[T, U] = T | U
type Option[T] = None | Some(T)
type Result[T, E] = Ok(T) | Err(E)

// 3. Type Parameters (Generics)
record Vector[T] {
    data: [T],
    len: Int
}

function Vector::new[T](cap: Int): Vector[T] = {
    Vector { data: Array::new(cap), len: 0 }
}

// 4. Protocols (Interfaces)
protocol Numeric {
    function add(other: Self) -> Self
    function sub(other: Self) -> Self
    function mul(other: Self) -> Self
}

impl Numeric for Int {
    function add(other: Int) = this + other
    function sub(other: Int) = this - other
    function mul(other: Int) = this * other
}
```

**테스트**: 20개 (types_test.fl)

### Task A.3: 동적 디스패치 엔진
**목표**: Multiple Dispatch 구현 (Julia 핵심!)

**파일**: `src/dispatch.fl` (300줄)

**구조**:
```
Method Registry
  ├─ Method 1: add(Int, Int) -> Int
  ├─ Method 2: add(Float, Float) -> Float
  ├─ Method 3: add(Int, Float) -> Float
  └─ Method 4: add(Array, Array) -> Array

Dispatch Resolution
  ├─ Type matching
  ├─ Specificity ranking
  └─ Method selection
```

**테스트**: 25개 (dispatch_test.fl)

---

## 4️⃣ **리스크 & 대응**

| 리스크 | 심각도 | 대응 |
|--------|--------|------|
| Dynamic Type 성능 | 🟠 중 | JIT 컴파일 & 캐싱 |
| Dispatch 복잡도 | 🟠 중 | 단계적 구현 (simple → complex) |
| 타입 안전성 | 🟡 낮음 | Type constraints 강제 |
| 메모리 오버헤드 | 🟡 낮음 | Struct optimization |

---

## 📅 **Phase A 스케줄**

```
2026-03-19 (Day 1)
└─ A.1: Julia 문법 분석 ✅

2026-03-20 (Day 2)
├─ A.2: 타입 시스템 확장 (50%)
└─ A.3: 디스패치 엔진 (준비)

2026-03-21 (Day 3)
├─ A.2: 완료 & 테스트 (18/20)
└─ A.3: 디스패치 엔진 (50%)

2026-03-22 (Day 4)
└─ A.3: 완료 & 테스트 (25/25)

2026-03-26 (Week 2)
├─ 최적화 & 성능 튜닝
└─ Phase B 준비
```

---

**버전**: Phase A v1.0
**상태**: 🟢 Task A.1 완료, A.2 시작 준비
**다음**: Task A.2 (FreeLang 타입 시스템 확장)
**예상 완료**: 2026-03-26
