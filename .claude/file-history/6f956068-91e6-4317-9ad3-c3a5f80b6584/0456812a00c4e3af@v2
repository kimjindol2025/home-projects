# FreeLang v1~v6 버전별 코드 예제 비교

**같은 문제를 v4, v5, v6으로 어떻게 푸는가?**

---

## 예제 1: Fibonacci 수열

### v4 - 기초 버전

```freelang
// v4: 함수형, 단순 제어
fn fib(n: i32): i32 {
  if n <= 1 { return n }
  return fib(n - 1) + fib(n - 2)
}

// 실행
for i in range(0, 10) {
  println(str(i) + " : " + str(fib(i)))
}
```

**특징**:
- ❌ while/loop 없음 → for...in만 사용
- ❌ 메모이제이션 직접 구현 어려움
- ✅ 간단한 재귀

---

### v5 - 확장 버전

```freelang
// v5: 모듈화 + 메모이제이션
module fib_module {
  let memo: {i32: i32} = {}

  fn fib_memo(n: i32): i32 {
    if let Some(result) = memo.get(n) {
      return result
    }

    let result = if n <= 1 {
      n
    } else {
      fib_memo(n - 1) + fib_memo(n - 2)
    }

    memo.set(n, result)
    result
  }
}

// 클로저로 map
let nums = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
let fibs = map(nums, |n| fib_memo(n))
```

**특징**:
- ✅ 모듈 지원
- ✅ HashMap으로 메모이제이션
- ✅ 클로저 (고차 함수)
- ⚠️ HashMap은 아직 dict 문법 필요 (미구현)

---

### v6 - 최적화 버전

```freelang
import { memoize } from "std:function"

// v6: 데코레이터 + stdlib 활용
@memoize
fn fib(n: i32): i32 {
  if n <= 1 { return n }
  return fib(n - 1) + fib(n - 2)
}

// 한 줄로 끝
println(map([0..10], |n| str(fib(n))).join(", "))
```

**특징**:
- ✓ @memoize 데코레이터 (v6)
- ✓ Range literal `[0..10]` (v6)
- ✓ 메서드 체이닝 (v6)
- ✓ stdlib:function 모듈

---

## 예제 2: HTTP 요청 + JSON 파싱

### v4 - 불가능

```freelang
// v4에는 HTTP 기능이 없음
// 이 코드는 작동하지 않음
❌ HTTP_CLIENT 없음
❌ JSON 파싱은 있으나 외부 데이터 불가능
```

---

### v5 - 기본 지원

```freelang
import { http } from "std:http"
import { json } from "std:json"

// v5: 모듈로 http 지원 (메서드 X)
let url = "https://jsonplaceholder.typicode.com/users/1"

// http::get 함수 호출
let response = http::get(url)
let data = json::parse(response.body)

// 데이터 접근
match data {
  Some(user) => {
    println("User: " + user["name"])
    println("Email: " + user["email"])
  },
  None => println("Failed to parse")
}
```

**특징**:
- ✓ http::get() 지원
- ✓ json::parse() 지원
- ✗ 메서드 체이닝 없음
- ✗ Response 구조체 문법 미완성

---

### v6 - 전문 수준

```freelang
import { http, json } from "std:http"

// v6: HTTP 메서드 + json 직접 파싱
let response = http::get("https://jsonplaceholder.typicode.com/users/1")
  .timeout(5000)
  .header("User-Agent", "FreeLang/v6")

let user = json::parse(response.body)?

println("Name: " + user.name)
println("Email: " + user.email)
println("Company: " + user.company.name)

// 배열 처리
let users = http::get("https://jsonplaceholder.typicode.com/users")
let user_list = json::parse(users.body)?
let names = filter(user_list, |u| u.id > 5)
  .map(|u| u.name)
  .join(", ")
println(names)
```

**특징**:
- ✓ 메서드 체이닝 (v6)
- ✓ ?연산자 (error propagation)
- ✓ 구조체 필드 접근 (v5+)
- ✓ 메서드 (v5+)

---

## 예제 3: 파일 처리 + 정규표현식

### v4 - 기본 지원

```freelang
// v4: 파일 읽기/쓰기만 가능
let content = read_file("data.txt")
let lines = split(content, "\n")

// 수동 필터링
let var filtered = []
for line in lines {
  if !contains(line, "#") {  // 주석 제거
    push(filtered, line)
  }
}

write_file("output.txt", join(filtered, "\n"))
```

**특징**:
- ✓ read_file, write_file 있음
- ❌ 정규표현식 없음
- ❌ 배열 필터링은 수동

---

### v5 - 정규표현식 추가

```freelang
import { regex } from "std:regex"

// v5: 정규표현식 지원 (기본)
let content = read_file("data.txt")
let lines = split(content, "\n")

let pattern = regex("^[a-z]+:\\s*(.*)$")
let var results = []

for line in lines {
  if let Some(match) = pattern.exec(line) {
    push(results, match.group(1))
  }
}

write_file("output.txt", join(results, "\n"))
```

**특징**:
- ✓ regex() 함수
- ✓ pattern.exec() 메서드
- ✓ group() 캡처 그룹
- ⚠️ 아직 리터럴 문법 아님

---

### v6 - 정규표현식 리터럴 + 함수형

```freelang
import { regex } from "std:regex"

// v6: 리터럴 문법 + 메서드 체이닝
let pattern = /^[a-z]+:\s*(.*)$/m  // 리터럴

let results = read_file("data.txt")
  .split("\n")
  .filter(|line| pattern.test(line))
  .map(|line| pattern.exec(line).group(1))

write_file("output.txt", results.join("\n"))
```

**특징**:
- ✓ /regex/ 리터럴 (v6)
- ✓ 메서드 체이닝 (v6)
- ✓ .test(), .exec() 메서드 (v6)
- ✓ 함수형 스타일

---

## 예제 4: 구조체 + 메서드

### v4 - 불가능

```freelang
// v4: 구조체 선언 불가능
❌ struct 키워드 없음
❌ 메서드 정의 불가능
```

---

### v5 - struct + impl

```freelang
// v5: 구조체 정의 및 메서드
struct Rectangle {
  width: i32,
  height: i32,
}

impl Rectangle {
  fn area(self): i32 {
    self.width * self.height
  }

  fn perimeter(self): i32 {
    (self.width + self.height) * 2
  }

  fn resize(mut self, w: i32, h: i32) {
    self.width = w
    self.height = h
  }
}

let rect = Rectangle { width: 10, height: 20 }
println("Area: " + str(rect.area()))
println("Perimeter: " + str(rect.perimeter()))
```

**특징**:
- ✓ struct 선언 (v5)
- ✓ impl 블록 (v5)
- ✓ self / &self / &mut self (v5)
- ✓ 메서드 호출 (v5)

---

### v6 - 데코레이터 + 고급 기능

```freelang
// v6: 데코레이터로 메타정보 추가
@derive(Debug, PartialEq)
struct Rectangle {
  width: i32,
  height: i32,
}

impl Rectangle {
  fn new(width: i32, height: i32): Self {
    Rectangle { width, height }
  }

  @inline
  fn area(self): i32 {
    self.width * self.height
  }

  fn scale(mut self, factor: i32) {
    self.width = self.width * factor
    self.height = self.height * factor
  }

  fn to_string(self): string {
    "Rectangle(" + str(self.width) + "x" + str(self.height) + ")"
  }
}

let rect = Rectangle::new(10, 20)
println(rect.to_string())
```

**특징**:
- ✓ @derive 데코레이터 (v6)
- ✓ @inline 데코레이터 (v6)
- ✓ Self 타입 (v6)
- ✓ 관례적 메서드 이름 (v6)

---

## 예제 5: 동시성 (Actor)

### v4 - 기본 Actor

```freelang
// v4: Actor + Channel
let ch = channel()

spawn {
  for i in range(0, 5) {
    send(ch, i * 2)
  }
}

for i in range(0, 5) {
  let val = recv(ch)
  println(str(val))
}
```

**특징**:
- ✓ spawn 블록
- ✓ channel 생성
- ✓ send/recv 메서드
- ❌ select (다중 채널) 없음

---

### v5 - Actor + 워크플로우

```freelang
// v5: Action으로 비즈니스 로직 캡슐화
action process_data {
  input: [i32],
  output: Result<[i32], string>,

  on_validate {
    if length(input) == 0 {
      return Err("Empty input")
    }
  }

  on_execute {
    return Ok(map(input, |x| x * 2))
  }

  on_success {
    println("Processing complete")
  }

  on_failure {
    println("Processing failed")
  }
}

let result = process_data([1, 2, 3, 4, 5])
match result {
  Ok(data) => println("Result: " + str(data)),
  Err(msg) => println("Error: " + msg)
}
```

**특징**:
- ✓ action 선언 (v5)
- ✓ 선언형 워크플로우 (v5)
- ✓ on_validate, on_execute, on_success, on_failure (v5)
- ✓ Result 타입 처리 (v5)

---

### v6 - Actor + Observability

```freelang
// v6: 메트릭 추적 + 워크플로우
@observe(metrics: ["duration", "error_rate"])
action process_data {
  input: [i32],
  output: Result<[i32], string>,

  on_execute {
    return Ok(
      map(input, |x| x * 2)
        .filter(|x| x > 2)
        .sort()
    )
  }
}

let result = process_data([1, 2, 3, 4, 5])

// 메트릭 확인
let duration = metrics::get("duration", process_data)
let error_rate = metrics::get("error_rate", process_data)
println("Duration: " + str(duration) + "ms")
println("Errors: " + str(error_rate * 100) + "%")
```

**특징**:
- ✓ @observe 데코레이터 (v6)
- ✓ metrics API (v6)
- ✓ 함수형 체이닝 (v6)
- ✓ 자동 관찰성 (v6)

---

## 예제 6: 모듈과 타입 안전

### v4 - 단일 파일

```freelang
// v4: 모든 코드가 한 파일에
// main.fl

fn sum(arr: [i32]): i32 {
  var total = 0
  for x in arr {
    total = total + x
  }
  return total
}

fn average(arr: [i32]): f64 {
  return i64(sum(arr)) / i64(length(arr))
}

println(str(sum([1, 2, 3, 4, 5])))
println(str(average([1, 2, 3, 4, 5])))
```

**특징**:
- ❌ 모듈 없음
- ❌ import/export 없음
- ✓ 함수 정의는 명확

---

### v5 - 모듈 분리

```freelang
// math.fl - 모듈
export fn sum(arr: [i32]): i32 {
  var total = 0
  for x in arr {
    total = total + x
  }
  return total
}

export fn average(arr: [i32]): f64 {
  return i64(sum(arr)) / i64(length(arr))
}

// main.fl - 메인
import { sum, average } from "./math.fl"

println(str(sum([1, 2, 3, 4, 5])))
println(str(average([1, 2, 3, 4, 5])))
```

**특징**:
- ✓ export로 공개 (v5)
- ✓ import로 사용 (v5)
- ✓ 모듈 분리 (v5)
- ✓ 타입 안전 (v5)

---

### v6 - 표준 라이브러리

```freelang
// main.fl
import { sum, average } from "std:math"
import { print_table } from "std:io"

let data = [1, 2, 3, 4, 5]
let result = {
  sum: sum(data),
  avg: average(data),
  count: length(data),
  min: min(data),
  max: max(data)
}

print_table(result)
```

**특징**:
- ✓ 표준 라이브러리 "std:*" (v6)
- ✓ 풍부한 내장 함수 (v6)
- ✓ 구조체 리터럴 (v5+)
- ✓ 타입 안전 (v4+)

---

## 예제 7: 에러 처리

### v4 - Result<T,E>

```freelang
// v4: Result 타입으로 에러 처리
fn divide(a: i32, b: i32): Result<i32, string> {
  if b == 0 {
    return Err("Division by zero")
  }
  return Ok(a / b)
}

match divide(10, 2) {
  Ok(result) => println("Result: " + str(result)),
  Err(msg) => println("Error: " + msg)
}

// ?연산자로 전파
fn safe_divide(a: i32, b: i32, c: i32): Result<i32, string> {
  let x = divide(a, b)?  // Err이면 즉시 반환
  let y = divide(x, c)?
  return Ok(y)
}
```

**특징**:
- ✓ Result<T,E> 타입 (v4)
- ✓ match로 처리 (v4)
- ✓ ? 연산자 (v4)

---

### v5 - try-catch 추가

```freelang
// v5: try-catch 문법도 지원
try {
  let result = divide(10, 0)?
  println("Result: " + str(result))
} catch err {
  println("Caught error: " + err.message)
} finally {
  println("Cleanup")
}

// action에서 자동 처리
action safe_compute {
  input: {a: i32, b: i32, c: i32},
  output: Result<i32, string>,

  on_validate {
    if input.b == 0 || input.c == 0 {
      return Err("Invalid input")
    }
  }

  on_execute {
    return Ok(
      divide(input.a, input.b)? / input.c
    )
  }

  on_failure {
    println("Computation failed")
  }
}
```

**특징**:
- ✓ try-catch-finally (v5)
- ✓ Result 타입 (v4+)
- ✓ ? 연산자 (v4+)
- ✓ action으로 자동 처리 (v5)

---

### v6 - 데코레이터 기반 에러 처리

```freelang
// v6: @error_handler 데코레이터
@error_handler(retries: 3, timeout: 5000)
fn risky_operation(url: string): Result<string, string> {
  return http::get(url)
    .expect("HTTP request failed")
    .text()
}

// 자동으로 3번 재시도, 5초 타임아웃
let result = risky_operation("https://api.example.com")

// 또는 명시적
match result {
  Ok(data) => println("Success: " + data),
  Err(e) => println("Final error: " + e)
}
```

**특징**:
- ✓ @error_handler 데코레이터 (v6)
- ✓ 자동 재시도 (v6)
- ✓ 타임아웃 처리 (v6)

---

## 버전별 기능 적용 행렬

```
기능            v4    v5      v6
─────────────────────────────────
기본 연산        ✓     ✓       ✓
함수            ✓     ✓       ✓
for...in        ✓     ✓       ✓
while/loop      ❌    ✓       ✓
배열            ✓     ✓       ✓
구조체          ❌    ✓       ✓
메서드          ❌    ✓       ✓
클로저          ❌    ✓       ✓
모듈            ❌    ✓       ✓
Action          ❌    ✓       ✓
Observability   ❌    ✓       ✓
정규표현식      ❌    △       ✓
HTTP            ❌    △       ✓
데이터베이스    ❌    ❌       ✓
데코레이터      ❌    △       ✓
메트릭          ❌    ✓       ✓
```

---

## 마이그레이션 경로

### v4 → v5

```freelang
// v4 코드
fn process(arr: [i32]): i32 {
  var total = 0
  for x in arr {
    total = total + x
  }
  return total
}

// v5에서도 그대로 작동 (호환성 100%)
// 하지만 다음과 같이 개선 가능:

// 방법 1: 모듈로 분리
// math.fl
export fn process(arr: [i32]): i32 { ... }

// main.fl
import { process } from "./math.fl"

// 방법 2: 클로저 활용
let result = fold(arr, 0, |acc, x| acc + x)

// 방법 3: Action으로 캡슐화
action sum_array {
  input: [i32],
  output: i32,
  on_execute { return fold(input, 0, |a, x| a + x) }
}
```

### v5 → v6

```freelang
// v5 코드
let pattern = regex("^\\d{3}$")
if pattern.test(phone) {
  println("Valid")
}

// v6에서 개선:
if /^\d{3}$/.test(phone) {  // 리터럴
  println("Valid")
}

// v5 코드
@observe(metrics: ["duration"])
fn expensive() { ... }

// v6에서는 메트릭이 자동 수집되고 더 많은 통계 제공
```

---

**이 예제들은 FreeLang의 진화를 직접 보여줍니다:**
- v4: 기초 (간결하고 명확)
- v5: 기능 (유연하고 강력)
- v6: 완성 (최적화되고 실용적)
