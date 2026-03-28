---
name: FV 2.0 Phase 1 - V 언어 분석
description: FV 2.0 프로젝트에서 V 언어의 구조, 문법, 타입 시스템, 라이브러리 생태계 분석
type: project
---

# FV 2.0 Phase 1: V 언어 완전 분석

**작성일**: 2026-03-19
**프로젝트**: FV 2.0 (V Language + FreeLang Integration)
**상태**: 🟡 분석 시작

---

## 목표

V 언어의 완전한 이해를 통해:
1. V 문법을 FreeLang 파서에 적응시킬 수 있는 방법 찾기
2. V 타입 시스템과 FreeLang의 호환성 매핑
3. V 라이브러리를 FreeLang 표준 라이브러리로 래핑하는 전략 수립

---

## V 언어 개요

### 공식 정보
- **웹사이트**: https://vlang.io
- **저장소**: https://github.com/vlang/v
- **라이선스**: MIT
- **개발 시작**: 2019년
- **현재 상태**: 0.4.x (빠르게 성장 중)

### 핵심 특징

#### 1. 빠른 컴파일 (<1초)
```v
// V 코드
fn main() {
  println('Hello, V!')
}

// 컴파일
$ v hello.v
# 결과: 단일 바이너리 (의존성 없음)
```

**특징**:
- C 코드로 트랜스파일
- 컴파일 시간 <1초 (대부분의 프로젝트)
- 최소 의존성 (기본은 표준 라이브러리만)

#### 2. 간단한 문법 (Go 영향)

**변수 선언**:
```v
x := 10              // 타입 추론
mut y := 20          // 가변성 명시
const z = 30         // 상수

name: string = 'Kim' // 명시적 타입 (선택)
```

**함수**:
```v
fn add(a int, b int) int {
  return a + b
}

// 마지막 표현식이 반환값
fn multiply(a int, b int) int {
  a * b
}
```

#### 3. 메모리 안전성 (Rust와 비슷)

```v
// NULL 안전성
fn safe_divide(a int, b int) ?int {
  if b == 0 {
    return none
  }
  return a / b
}

// 사용
result := safe_divide(10, 2) or { panic('Division failed') }
// 또는
if value := safe_divide(10, 2) {
  println(value)
}
```

---

## V 문법 상세 분석

### 1. 기본 타입

#### 정수형
```v
i8, i16, i32, i64, int
u8, u16, u32, u64, uint
```

#### 부동소수점
```v
f32, f64
```

#### 문자열 & 문자
```v
string        // UTF-8
rune          // 단일 문자
byte          // u8 별칭
```

#### 불린
```v
bool
```

#### 컬렉션
```v
[]int         // 배열 (동적)
[5]int        // 고정 크기 배열
map[string]int // 맵
```

### 2. 제어문

#### if/else
```v
if x > 0 {
  println('positive')
} else if x < 0 {
  println('negative')
} else {
  println('zero')
}
```

#### for 루프
```v
// C 스타일
for i := 0; i < 10; i++ {
  println(i)
}

// While 스타일
for condition {
  // ...
}

// 범위
for i in 0..10 {
  println(i)
}

// 배열/맵 순회
for item in arr {
  println(item)
}
```

#### match (패턴 매칭)
```v
match value {
  1 {
    println('one')
  }
  2, 3 {
    println('two or three')
  }
  else {
    println('other')
  }
}
```

### 3. 구조체 & 메서드

```v
struct Person {
  name: string
  age: int
  mut email: string  // 가변 필드
}

fn (p Person) greet() string {
  return 'Hello, ${p.name}!'
}

fn (mut p Person) set_email(email string) {
  p.email = email
}

// 사용
mut person := Person{
  name: 'Kim',
  age: 30,
  email: 'kim@example.com',
}
person.set_email('kim2@example.com')
```

### 4. 인터페이스

```v
interface Reader {
  read(mut buf []u8) ?int
}

struct FileReader {
  path: string
}

fn (f FileReader) read(mut buf []u8) ?int {
  // 파일에서 읽기
}
```

### 5. 제네릭

```v
fn first<T>(arr []T) ?T {
  if arr.len == 0 {
    return none
  }
  return arr[0]
}

// 사용
first([1, 2, 3]) // 1 반환
first(['a', 'b']) // 'a' 반환
```

### 6. 에러 처리

#### Option (널 안전)
```v
fn find(arr []int, target int) ?int {
  for i, v in arr {
    if v == target {
      return i
    }
  }
  return none
}

// 사용
if index := find([1, 2, 3], 2) {
  println('Found at: $index')
} else {
  println('Not found')
}
```

#### Result (에러 메시지)
```v
fn parse_int(s string) !int {
  // ... 파싱
  if error {
    return error('Invalid number')
  }
  return 42
}

// 사용 (? 연산자)
value := parse_int('42') or {
  eprintln('Parse error: ${err}')
  return
}
```

---

## V 표준 라이브러리

### 모듈 구조
```
v/
├── vlib/              # 표준 라이브러리
│   ├── os/           # 운영체제 (파일, 환경)
│   ├── io/           # 입출력
│   ├── json/         # JSON 처리
│   ├── crypto/       # 암호화
│   ├── net/          # 네트워킹
│   ├── http/         # HTTP 클라이언트/서버
│   ├── sql/          # 데이터베이스
│   ├── time/         # 시간
│   ├── math/         # 수학
│   ├── strings/      # 문자열
│   └── ...           # 기타
```

### 주요 라이브러리

#### os (파일시스템)
```v
import os

// 파일 읽기
content := os.read_file('file.txt')?

// 파일 쓰기
os.write_file('file.txt', 'content')?

// 디렉토리 생성
os.mkdir('dir')?

// 환경 변수
env_var := os.getenv('HOME')
```

#### io (입출력)
```v
import io

reader := io.FileReader{
  path: 'data.txt'
}

data := reader.read(1024)?
```

#### json (JSON 처리)
```v
import json

data := json.decode(MyStruct, json_string)?
json_string := json.encode(obj)
```

#### http (HTTP)
```v
import net.http

// 서버
server := http.Server{
  port: 8080,
  handler: fn(req http.Request) http.Response {
    return http.Response{
      status_code: 200,
      body: 'Hello!',
    }
  },
}
server.listen_and_serve()

// 클라이언트
resp := http.get('https://example.com')?
println(resp.body)
```

#### sql (데이터베이스)
```v
import db.sqlite

db := sqlite.connect('app.db')?

// 쿼리
users := db.exec('SELECT * FROM users')?

// 매개변수화 쿼리
user := db.exec('SELECT * FROM users WHERE id = ?', [id])?
```

---

## V와 FreeLang의 타입 매핑

| V 타입 | FreeLang 타입 | 변환 전략 |
|--------|---------------|----------|
| `int` | `i64` | 기본 정수 |
| `i32` | `i32` | 명시적 |
| `u32` | `u32` | 부호 없음 |
| `f64` | `f64` | 부동소수점 |
| `string` | `String` | 동일 |
| `[]T` | `Vec(T)` | 동적 배열 |
| `[N]T` | `[T; N]` | 고정 배열 |
| `?T` | `Option(T)` | NULL 안전 |
| `!T` | `Result(T, String)` | 에러 처리 |
| `map[K]V` | `HashMap(K, V)` | 맵 |
| `struct` | `type Struct = {...}` | 구조체 |

---

## 다음 단계

### Task 1.2: FreeLang 현재 상태 분석
- FreeLang의 IR (Intermediate Representation) 구조
- FreeLang의 표준 라이브러리
- 기존 컴파일 파이프라인

### Task 1.3: 통합 지점 설계
- V 문법 파서 어댑터 설계
- 타입 매핑 테이블 작성
- 호환성 레벨 정의

---

**상태**: 🟡 V 언어 분석 완료 대기
