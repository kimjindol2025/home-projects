---
name: FV-Lang이란 무엇인가
description: FV-Lang의 정의, 특징, 현재 상태, 그리고 우리의 전략
type: project
---

# FV-Lang: 완전 이해 가이드

**작성일**: 2026-03-19
**현재 상태**: Phase 1-6 완료, 자체호스팅 준비 중

---

## 🎯 FV-Lang의 정의

### 공식 설명

**"순수 함수형 프로그래밍을 위한 현대적 언어"**

- **개발언어**: Rust (2,197줄)
- **런타임**: Rust 기반 인터프리터 + 런타임
- **코드생성**: FV → C (트랜스파일)
- **목표**: Haskell + Lisp + ML의 영향, 순수함수 중심

---

## 🏗️ FV-Lang의 구조

### 계층 구조

```
Application Layer
    ↓
FV-Lang (우리의 언어)
    ↓
Rust 런타임 + 컴파일러
    ↓
C 코드 생성 (또는 직접 해석)
    ↓
바이너리 실행
```

### 파이프라인

```
FV 소스 코드
    ↓ (Lexer: src/lexer.rs)
토큰 스트림
    ↓ (Parser: src/parser.rs)
AST (Abstract Syntax Tree)
    ↓ (Type Checker: src/type_checker.rs)
타입 검증 완료 AST
    ↓ (Code Generator: src/codegen.rs)
C 코드 또는 IR
    ↓ (gcc 또는 해석)
실행 결과
```

---

## ✨ 핵심 특징

### 1. 순수 함수형 (Pure Functional)

```fv
fn add(a: i64, b: i64) -> i64 {
  a + b
}

// Side effect 없음
// 같은 입력 = 항상 같은 출력
```

**의미**: 함수가 외부 상태를 변경하지 않음

### 2. 불변성 (Immutability)

```fv
let x = 5;
// x = 10;  // ❌ 불가능 (기본 불변)

let mut y = 5;
// y = 10;  // ✅ 가능 (명시적 mut)
```

### 3. 패턴 매칭 (Pattern Matching)

```fv
type Option =
  | Some(i64)
  | None

fn unwrap(opt: Option) -> i64 {
  match opt {
    Some(val) => val,
    None => 0,
  }
}
```

### 4. 고차 함수 (Higher-Order Functions)

```fv
fn apply(f: fn(i64) -> i64, x: i64) -> i64 {
  f(x)
}

fn double(x: i64) -> i64 {
  x * 2
}

let result = apply(double, 5);  // 10
```

### 5. 타입 안전성 (Type Safety)

```fv
fn safe_divide(a: i64, b: i64) -> Option(i64) {
  if b == 0 {
    None
  } else {
    Some(a / b)
  }
}
```

---

## 📊 현재 구현 상태

### Phase 1-6: 완료됨

| Phase | 내용 | 상태 | 줄수 |
|-------|------|------|------|
| 1 | 기본 구조 (CLI, Makefile) | ✅ | 200 |
| 2 | Lexer (50개 토큰) | ✅ | 500 |
| 3 | Parser (재귀 하강) | ✅ | 700 |
| 4 | Code Generator (FV→C) | ✅ | 500 |
| 5 | Runtime (인터프리터) | ✅ | 400 |
| 6 | Stdlib (I/O, String, Array) | ✅ | 1,200+ |

**총**: 2,197줄 Rust

### Stdlib 모듈

```
src/stdlib/
├── io.rs          - 입출력 (print, read)
├── string.rs      - 문자열 처리
├── array.rs       - 배열 연산
├── math.rs        - 수학 함수
├── collections.rs - 맵, 셋
├── bytes.rs       - 바이트 처리
├── fs.rs          - 파일시스템
├── sys.rs         - 시스템 함수
├── elf.rs         - ELF 포맷
├── advanced.rs    - 고급 함수
└── types.rs       - 타입 유틸
```

---

## 💡 우리가 지금 하는 것

### Singularity Phase 1: Bootstrap (자체호스팅)

**목표**: FV 컴파일러를 FV로 작성하여 증명

#### 현재 상황

```
기존 (Rust): fvc 컴파일러 (2,197줄)
    ↓ 포팅
우리가 작성: FV 컴파일러 (2,287줄 FV)
    ↓ 컴파일 (Rust fvc 사용)
C 코드
    ↓ gcc
FV 컴파일러 바이너리
    ↓ 자신으로 재컴파일
C 코드 (동일한가?)
    ↓ (결정론적 증명)
✅ 자체호스팅 성공!
```

#### 우리가 작성한 5개 모듈 (FV)

1. **lexer.fv** (430줄)
   - 토큰화
   - 50개 토큰 타입 정의

2. **ast.fv** (175줄)
   - 구문 트리 구조

3. **parser.fv** (520줄)
   - 재귀 하강 파싱

4. **type_checker.fv** (280줄)
   - 타입 검사

5. **codegen.fv** (390줄)
   - FV → C 코드 생성

**합계**: 2,287줄 FV

---

## 🎨 FV 문법 예제

### Hello World

```fv
fn main() -> i64 {
  print("Hello, FV!");
  0
}
```

### 함수

```fv
fn factorial(n: i64) -> i64 {
  if n <= 1 {
    1
  } else {
    n * factorial(n - 1)
  }
}
```

### 타입 정의

```fv
type Person = {
  name: String,
  age: i64,
}

type Result =
  | Ok(String)
  | Error(String)
```

### 패턴 매칭

```fv
fn describe(result: Result) -> String {
  match result {
    Ok(msg) => string_concat("Success: ", msg),
    Error(err) => string_concat("Failed: ", err),
  }
}
```

### 고차 함수

```fv
fn map(f: fn(i64) -> i64, arr: Vec(i64)) -> Vec(i64) {
  // TODO: 구현
}

let doubled = map(fn(x: i64) -> i64 { x * 2 }, [1, 2, 3]);
```

---

## 🚀 FV의 강점

### 1. 작은 핵심
- 2,197줄 Rust로 완전한 컴파일러
- 복잡하지 않은 아키텍처
- 이해하기 쉬운 구조

### 2. 결정론적 (Deterministic)
- 같은 입력 → 항상 같은 출력
- 난수 없음
- 부동소수점 신중히 다룸

### 3. 이식성
- Rust로 작성 (크로스플랫폼)
- C로 생성 (모든 시스템)
- 단일 바이너리

### 4. 함수형 패러다임
- 순수 함수 중심
- 부작용 최소화
- 추론하기 쉬운 코드

---

## ⚠️ 현재 제한사항

### 1. 미완성 기능

```fv
// TODO: 아직 구현 안 됨
- Lazy evaluation
- 고급 패턴 매칭 (리스트, 중첩)
- Generic type parameters
- Module system
- 문자열 보간
```

### 2. 성능

- **런타임**: 인터프리터 (느림)
- **개선책**: JIT 컴파일 (TODO)
- **현재**: 작은 프로그램은 빠름, 큰 프로그램은 느림

### 3. 에러 처리

- Result/Option 타입 있음
- 런타임 에러 스택 트레이스 부족

---

## 🔮 우리의 전략

### 왜 자체호스팅인가?

**이유 1: 완성도 검증**
```
FV로 컴파일러 작성 가능
= FV는 "real" 언어
```

**이유 2: 자기개선의 기초**
```
컴파일러가 자신의 코드를 분석
→ 최적화 제안 생성
→ 자동 리팩토링
→ Singularity
```

**이유 3: 결정론적 증명**
```
Round 1: Rust fvc로 FV 컴파일러 컴파일
Round 2: FV 컴파일러로 자신 재컴파일
Round 3: 3회 반복해서 SHA256 확인

모두 동일? → 완벽한 증명
```

### Phase 1 로드맵 (2주)

```
Week 1 (지금):
  ✅ 5개 모듈 작성 (2,287줄)
  ⏳ Stdlib 최소화

Week 2:
  ⏳ Main module (compiler.fv)
  ⏳ 자체호스팅 테스트
  ⏳ 결정론적 증명
```

---

## 📈 "미니 언어"인가?

### 사용자 질문: "FV는 미니 언어야?"

**답: 아니다. 하지만 아직 작고 젊다.**

#### 비교

| 언어 | 복잡도 | 용도 | 성숙도 |
|------|--------|------|--------|
| Python | 높음 | 범용 | 매우 높음 |
| Rust | 높음 | 시스템 | 높음 |
| **FV-Lang** | **중간** | **함수형** | **초기** |
| Lisp | 낮음 | 함수형 | 높음 |

#### 특징

- **미니**: 코드 2K줄, 초기 Phase
- **완전**: 완전한 컴파일러 파이프라인 (Lexer→Parser→TypeChecker→CodeGen)
- **실용**: stdlib, 런타임, I/O 지원
- **학습용**: 컴파일러 이해에 좋음

### 우리가 증명할 것

```
"작은 크기가 완성도를 결정하지 않는다"

FV는:
✓ Haskell 같은 패턴 매칭
✓ ML 같은 타입 시스템
✓ Lisp 같은 고차 함수
✓ + 결정론적 실행
✓ + 자체호스팅 가능

= "작은 FV가 자신을 컴파일할 수 있다" 증명
```

---

## 🎯 최종 평가

**FV-Lang은:**

1. **완전한 언어** - 컴파일러, 런타임, stdlib 포함
2. **함수형** - 순수함수, 패턴 매칭, 고차함수
3. **작지만 강력** - 2K줄로 완성도 높음
4. **자체호스팅 가능** - 우리가 증명하는 중
5. **Singularity의 기초** - 자기개선 시스템의 첫 단계

**다음**: 자체호스팅 완료 → 자기진화 시작

---

**상태**: 🟡 완전히 이해됨, Phase 1 계속 진행 중
