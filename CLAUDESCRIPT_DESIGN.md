# ClaudeScript (CS) 언어 설계 명세서

**상태**: 설계 단계 (검증 기반)
**버전**: 0.1.0 (초안)
**기반**: FreeLang v2.6.0 인프라
**목표**: Claude가 쓰기 편하고, 컴파일 가능하고, AI 실행 가능한 언어

---

## 1. 언어 철학

### 핵심 원칙

```
1️⃣ 거짓 없이
   - 선언한 기능은 반드시 동작
   - 제한사항은 명확히 표시
   - 테스트로 검증

2️⃣ Claude 친화성
   - 간단한 문법 (학습 비용 낮음)
   - JSON 기반 (파싱 용이)
   - 타입 안전성 (예측 가능)

3️⃣ 안전성 우선
   - null 참조 불가
   - 암묵적 타입 변환 금지
   - 범위 검사 필수
   - 컴파일 타임 검증

4️⃣ 실용성
   - FreeLang 바이트코드로 컴파일
   - 네이티브 성능 (C 바이딩 가능)
   - 모듈 시스템 지원
```

---

## 2. 문법 설계

### 2.1 기본 타입

```typescript
// 기본 타입 (모두 선택적)
let name: string = "Claude";     // 문자열
let age: i32 = 30;               // 32비트 정수
let score: f64 = 95.5;           // 64비트 실수
let active: bool = true;         // 불린
let items: Array<i32> = [];      // 배열
let data: Object<string> = {};   // 객체

// 타입 추론 (컴파일러가 결정)
let x = 42;                      // 타입: i32
let msg = "hello";               // 타입: string

// Optional (null 허용)
let maybe: Option<i32> = None;   // 아직 값 없음
let maybe2: Option<string> = Some("value");
```

### 2.2 함수 정의

```typescript
// 기본 함수
fn add(a: i32, b: i32) -> i32 {
  a + b
}

// 선택적 반환값 (Option 사용)
fn divide(a: f64, b: f64) -> Option<f64> {
  if b == 0.0 {
    None
  } else {
    Some(a / b)
  }
}

// 제너릭 함수
fn first<T>(items: Array<T>) -> Option<T> {
  if items.length() == 0 {
    None
  } else {
    Some(items[0])
  }
}
```

### 2.3 제어 구문

```typescript
// if/else
if x > 0 {
  println("양수");
} else {
  println("음수");
}

// match (패턴 매칭)
match value {
  None => println("없음"),
  Some(x) => println("값: " + to_string(x)),
}

// for 루프
for i in 0..10 {
  println(to_string(i));
}

// while 루프
while count < 10 {
  count = count + 1;
}

// try/catch (에러 처리)
try {
  let result = risky_operation();
  println(result);
} catch err {
  println("에러: " + err);
}
```

### 2.4 배열과 객체

```typescript
// 배열 (타입 안전)
let nums: Array<i32> = [1, 2, 3, 4, 5];
let first = nums[0];              // ✓ 유효
let invalid = nums[100];          // ✗ 컴파일 에러 또는 None
let safe = nums.get(100);         // Option<i32> 반환

// 객체 (타입 안전)
let person: Object<string> = {
  "name": "Alice",
  "city": "Seoul"
};
let name = person["name"];        // Some("Alice")
let unknown = person["age"];      // None (없는 키)

// Map (제너릭)
let scores: Map<string, i32> = {};
scores.set("Alice", 95);
scores.set("Bob", 87);
let alice_score = scores.get("Alice");  // Some(95)
```

---

## 3. 타입 시스템 (안전성 보장)

### 3.1 Optional 타입

모든 참조 타입은 기본적으로 `Option<T>`:

```typescript
// 선언
let maybe_string: Option<string> = Some("hello");
let nothing: Option<i32> = None;

// 사용 (패턴 매칭 강제)
match value {
  Some(x) => println(x),        // x는 string 타입 보장
  None => println("없음"),
}

// Option 메서드
value.is_some()          // bool
value.is_none()          // bool
value.unwrap()           // T (None이면 panic!)
value.unwrap_or(default) // T (기본값 제공)
value.map(f)             // Option<U> (함수 적용)
```

### 3.2 타입 검증 (컴파일 타임)

```typescript
// ✗ 컴파일 에러 (타입 불일치)
let x: i32 = "hello";

// ✗ 컴파일 에러 (암묵적 변환 금지)
let result = 5 + "3";

// ✓ 명시적 변환 필요
let result = 5 + to_i32("3");
let result = to_string(5) + "3";

// ✗ 컴파일 에러 (null 참조)
let s: string = null;

// ✓ Option 사용
let maybe_s: Option<string> = None;
```

### 3.3 범위 검사 (런타임)

```typescript
let arr = [1, 2, 3];

// ✗ 런타임 에러 (범위 벗어남)
let val = arr[-1];
let val = arr[999];

// ✓ 안전한 접근
let val = arr.get(0);         // Option<i32> = Some(1)
let val = arr.get(999);       // Option<i32> = None
let val = arr.get_or(999, 0); // i32 = 0 (기본값)
```

---

## 4. 컴파일 파이프라인

### 4.1 클로드 코드 → ClaudeScript AST

```json
{
  "type": "program",
  "definitions": [
    {
      "type": "function",
      "name": "add",
      "params": [
        {"name": "a", "type": {"base": "i32"}},
        {"name": "b", "type": {"base": "i32"}}
      ],
      "return_type": {"base": "i32"},
      "body": [
        {
          "type": "return",
          "value": {
            "type": "binary_op",
            "op": "+",
            "left": {"type": "ref", "name": "a"},
            "right": {"type": "ref", "name": "b"}
          }
        }
      ]
    }
  ],
  "instructions": [
    {
      "type": "call",
      "function": "add",
      "args": [
        {"type": "literal", "value_type": "i32", "value": 5},
        {"type": "literal", "value_type": "i32", "value": 3}
      ]
    }
  ]
}
```

### 4.2 ClaudeScript AST → FreeLang 바이트코드

```
(program
  (define-function add (a:i32 b:i32) -> i32
    (+ a b)
  )
  (call add 5 3)
)
```

### 4.3 FreeLang 바이트코드 → 네이티브 코드 (기존 FreeLang 컴파일러)

```c
int add(int a, int b) {
  return a + b;
}

void main() {
  int result = add(5, 3);
  printf("%d\n", result);  // 출력: 8
}
```

---

## 5. 구현 계획

### Phase 1: 기초 (1주일)
- [ ] ClaudeScript JSON AST 정의
- [ ] Lexer 구현 (문법 토큰화)
- [ ] Parser 구현 (AST 생성)
- [ ] Type Checker 구현 (안전성 검증)
- [ ] 기본 테스트 20개 작성

### Phase 2: 코드 생성 (1주일)
- [ ] ClaudeScript → FreeLang 변환
- [ ] Option 타입 런타임 지원
- [ ] 범위 검사 런타임 구현
- [ ] 통합 테스트 30개 작성

### Phase 3: 표준 라이브러리 (1주일)
- [ ] 기본 함수 (to_string, to_i32, 등)
- [ ] 배열 함수 (get, push, pop, etc.)
- [ ] 문자열 함수 (substring, length, etc.)
- [ ] 파일 I/O 함수
- [ ] 라이브러리 테스트 40개 작성

### Phase 4: 고급 기능 (1주일)
- [ ] 모듈 시스템 (import/export)
- [ ] 제너릭 타입 지원
- [ ] 비동기 함수 (async/await)
- [ ] 고급 테스트 50개 작성

---

## 6. 안전성 보장

### 6.1 타입 안전성

| 문제 | CLAUDELang v6.0 | ClaudeScript |
|------|---|---|
| 타입 불일치 | ⚠️ 암묵적 변환 | ✅ 컴파일 에러 |
| Null 참조 | ⚠️ undefined 반환 | ✅ Option<T> 강제 |
| 배열 범위 | ⚠️ 범위 검사 없음 | ✅ None 또는 panic |
| 함수 정의 | ⚠️ 50개 미만 실제 구현 | ✅ 100% 구현 보증 |
| 테스트 | ⚠️ 기본 케이스만 | ✅ 모든 엣지 케이스 |

### 6.2 검증 프로세스

```
ClaudeScript 소스코드
    ↓
① 렉싱 (Tokenization)
    ↓
② 파싱 (AST 생성)
    ↓
③ 타입 검사 (Type Checking)
   ✓ 모든 참조 타입 검증
   ✓ 암묵적 변환 금지
   ✓ 범위 검사 태그
    ↓
④ 코드 생성 (Code Generation)
    ↓
⑤ 런타임 검증
   ✓ 범위 검사 실행
   ✓ Null 참조 보호
   ✓ 함수 존재 검증
    ↓
✅ 정상 실행 또는 명확한 에러
```

---

## 7. 사용 예시

### 7.1 간단한 프로그램

```typescript
// calculate.cs
fn factorial(n: i32) -> i32 {
  if n <= 1 {
    1
  } else {
    n * factorial(n - 1)
  }
}

fn main() {
  let result = factorial(5);
  println("5! = " + to_string(result));
}
```

**실행**:
```bash
$ claudescript compile calculate.cs
✅ 컴파일 완료

$ claudescript run calculate.cs
5! = 120
```

### 7.2 배열 처리

```typescript
// array_operations.cs
fn sum_array(items: Array<i32>) -> i32 {
  let total: i32 = 0;
  for i in 0..items.length() {
    total = total + items[i];
  }
  total
}

fn main() {
  let numbers = [1, 2, 3, 4, 5];
  let result = sum_array(numbers);
  println("합: " + to_string(result));
}
```

### 7.3 에러 처리

```typescript
// error_handling.cs
fn safe_divide(a: f64, b: f64) -> Option<f64> {
  if b == 0.0 {
    None
  } else {
    Some(a / b)
  }
}

fn main() {
  let result1 = safe_divide(10.0, 2.0);
  match result1 {
    Some(val) => println("결과: " + to_string(val)),
    None => println("0으로 나눌 수 없음"),
  }

  let result2 = safe_divide(10.0, 0.0);
  match result2 {
    Some(val) => println("결과: " + to_string(val)),
    None => println("0으로 나눌 수 없음"),
  }
}
```

---

## 8. Claude 생성 코드 예시

Claude가 ClaudeScript 코드를 직접 생성할 수 있도록 설계:

```
🤖 Claude:
"Fibonacci 수열을 계산하는 함수를 만들어 줄게요:"

fn fibonacci(n: i32) -> i32 {
  if n <= 1 {
    n
  } else {
    fibonacci(n - 1) + fibonacci(n - 2)
  }
}

"이 함수는:
- 타입 안전함 (i32 → i32)
- 재귀적으로 안전
- 범위 검사 자동
- 컴파일 보장됨"
```

---

## 9. 검증 및 테스트

### 9.1 단위 테스트

```typescript
#[test]
fn test_basic_arithmetic() {
  assert_eq(2 + 2, 4);
  assert_eq(10 - 3, 7);
  assert_eq(5 * 3, 15);
}

#[test]
fn test_option_type() {
  let some_val: Option<i32> = Some(42);
  let none_val: Option<i32> = None;

  assert_true(some_val.is_some());
  assert_true(none_val.is_none());
}

#[test]
fn test_array_bounds() {
  let arr = [1, 2, 3];
  let val = arr.get(100);  // Option<i32> = None
  assert_true(val.is_none());
}

#[test]
fn test_type_checking() {
  // 컴파일 타임에 검사되므로
  // 런타임 타입 불일치는 불가능
}
```

### 9.2 회귀 테스트

```bash
# 매주 검증 (진정성 보장)
$ claudescript test --all
Running tests...
✅ 단위 테스트: 100/100 통과
✅ 통합 테스트: 50/50 통과
✅ 타입 안전성: 모든 케이스 통과
✅ 런타임 안전성: 패닉 없음

총 150개 테스트 통과 (0 실패)
```

---

## 10. 거짓 없는 약속

### 10.1 할 것 ✅

```
✅ 타입 안전성 (컴파일 타임)
✅ Null 안전성 (Option 강제)
✅ 범위 검사 (런타임)
✅ 함수 정의 (100% 구현)
✅ 테스트 검증 (모든 엣지 케이스)
✅ 에러 메시지 (명확한 원인)
✅ FreeLang 통합 (완전한 컴파일)
✅ Claude 친화성 (간단한 문법)
```

### 10.2 하지 않을 것 ❌

```
❌ 암묵적 타입 변환
❌ Null 참조 허용
❌ 범위 검사 없는 배열 접근
❌ 선언만 하고 구현 안 한 함수
❌ 테스트 없는 기능 배포
❌ 모호한 에러 메시지
❌ 클로드라 불완전한 컴파일
❌ 과장된 성능 주장
```

### 10.3 진정성 검증

```
매주 금요일 22:00 자동 검증:

1. 코드 실행 테스트
   - 모든 100+ 테스트 실행
   - 패닉 또는 undefined 있으면 빌드 실패

2. 타입 안전성 검증
   - 타입 불일치 샘플 10개 컴파일 시도
   - 모두 에러 나야 함

3. 성능 벤치마크
   - 실제 성능 측정 (과장 금지)
   - 경쟁사 비교 공정한지 확인

4. 문서와 코드 일치 검증
   - 문서에 있는 기능 모두 동작하는지 확인
   - 구현되지 않은 기능은 문서에서 제거
```

---

## 11. 다음 단계

1. **ClaudeScript 문법** → `.g4` Antlr 문법 파일
2. **TypeScript Compiler** → TypeScript 구현
3. **테스트 작성** → 150개 테스트 케이스
4. **FreeLang 통합** → 바이트코드 생성 및 실행
5. **CI/CD 자동화** → 매 커밋 검증
6. **문서화** → 사용자 가이드 및 API 문서

---

## 12. 최종 선언

```
"ClaudeScript는 거짓 없는 언어다.

선언한 기능은 반드시 동작한다.
제한사항은 명확히 표시된다.
모든 기능은 테스트로 검증된다.

이것이 우리의 약속이다."
```

---

**작성**: 2026-03-07
**상태**: 설계 승인 대기
**다음 검토**: 아키텍처 상세 설계

