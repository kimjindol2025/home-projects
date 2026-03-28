---
name: Singularity Playground - Phase 1 FV-Lang Bootstrap
description: FV-Lang 자체호스팅 구현 (FV로 작성된 컴파일러)
type: project
---

# Singularity Playground - Phase 1: FV-Lang Bootstrap

**시작일**: 2026-03-19
**목표**: FV-Lang 컴파일러를 FV 언어로 재작성하여 자체호스팅 증명
**규모**: 약 4,000줄 FV 코드 (기존 Rust 2,197줄 포팅)
**상태**: 🟡 Step 1 완료 - 핵심 모듈 작성 완료

---

## 🎯 최종 목표

```
fv-compiler.fv (자기 자신을 컴파일할 수 있는 컴파일러)
    ↓ (기존 Rust 컴파일러로 컴파일)
fv-compiler.c
    ↓ (gcc로 컴파일)
fv-compiler (바이너리)
    ↓ (자기 자신으로 재컴파일)
fv-compiler.c (동일한 결과)
    ↓ (deterministic proof)
✅ 자체호스팅 증명 완료
```

---

## 📋 Step 1: 핵심 모듈 작성 (COMPLETE)

### 파일 생성 완료

1. **lexer.fv** (430줄)
   - 50개 토큰 타입 정의
   - 상태 머신: 공백, 주석, 식별자, 숫자, 문자열, 연산자, 구분자
   - 보안: 입력 크기 검증, NULL 바이트 검사
   - 함수: `tokenize(input: String) -> Result(Vec(Token), String)`

2. **ast.fv** (175줄)
   - Program, Definition, Statement, Expression 구조
   - Literal 타입 (Integer, Float, String, Boolean)
   - Type 시스템 정의 (I64, F64, Str, Bool, Array, Tuple, Custom)
   - 헬퍼 함수: create_program, create_function, create_let_binding, 등

3. **parser.fv** (520줄)
   - 재귀 하강 파싱 (Recursive Descent)
   - 토큰 → AST 변환
   - 우선순위 기반 식 파싱
     - parse_equality (==, !=)
     - parse_comparison (<, <=, >, >=)
     - parse_addition (+, -)
     - parse_multiplication (*, /, %)
     - parse_unary (!, -)
     - parse_call (함수 호출, 배열 인덱싱)
     - parse_primary (리터럴, 변수, 괄호)
   - 구문 파싱: let, return, if, match (TODO)
   - 함수: `parse(parser: Parser) -> ParseResult(Program)`

4. **type_checker.fv** (280줄)
   - Type Environment 관리 (함수, 변수)
   - AST 검증 및 타입 추론
   - 정의 수집 (1차 pass)
   - 명령문 검사 (2차 pass)
   - 식 타입 검사
   - 함수: `check_program(prog: Program) -> CheckResult(Program)`

5. **codegen.fv** (390줄)
   - AST → C 코드 변환
   - 정의 생성 (함수 선언 및 본체)
   - 구문 생성 (let, return, if/else)
   - 식 생성 (리터럴, 변수, 이항/단항 연산, 함수 호출)
   - 들여쓰기 및 포맷팅
   - 함수: `generate_code(prog: Program) -> GenResult(String)`

### 통계
- **총 FV 코드**: 1,795줄 (5개 파일)
- **예상 최종 규모**: 4,000줄 (stdlib + 통합 로직 추가 시)

---

## 🔄 Step 2: 필요한 stdlib 함수 (IN PROGRESS)

FV 컴파일러가 동작하려면 다음 stdlib이 필요:

### 문자열 함수
```fv
fn string_length(s: String) -> i64
fn string_contains(s: String, sub: String) -> bool
fn char_at(s: String, pos: i64) -> String
fn string_concat(a: String, b: String) -> String
fn string_to_int(s: String) -> i64
fn string_to_float(s: String) -> f64
```

### 벡터 함수
```fv
fn vec_length(v: Vec(T)) -> i64
fn vec_get(v: Vec(T), i: i64) -> T
fn vec_append(v: Vec(T), e: T) -> Vec(T)
```

### 숫자 변환
```fv
fn i64_to_string(n: i64) -> String
fn f64_to_string(f: f64) -> String
```

### 맵 (Hash Map)
```fv
fn create_map() -> Map(K, V)
fn map_get(m: Map(K, V), k: K) -> Option(V)
fn map_insert(m: Map(K, V), k: K, v: V) -> Map(K, V)
```

**현황**: 기존 FV-Lang stdlib에서 부분 구현됨
**작업**: FV 컴파일러가 필요로 하는 최소 기능 확인 후 추가/보완

---

## 🛠️ Step 3: Main Compiler Module (TODO)

```fv
// compiler.fv: 통합 파이프라인
fn compile(source: String) -> Result(String, String) {
  // 1. Lexer
  let tokens = tokenize(source)?;

  // 2. Parser
  let program = parse(tokens)?;

  // 3. Type Checker
  let checked_program = check_program(program)?;

  // 4. Code Generator
  let c_code = generate_code(checked_program)?;

  Ok(c_code)
}
```

**작업 예상**: 100줄

---

## 🧪 Step 4: 자체호스팅 증명 (TODO)

### 테스트 계획

1. **기존 Rust 컴파일러로 FV 컴파일러 컴파일**
   ```bash
   fvc src/lexer.fv → lexer.c
   fvc src/parser.fv → parser.c
   fvc src/type_checker.fv → type_checker.c
   fvc src/codegen.fv → codegen.c
   fvc src/compiler.fv → compiler.c
   gcc -o fvc *.c
   ```

2. **FV 컴파일러로 자신 재컴파일**
   ```bash
   ./fvc src/lexer.fv → lexer.c (동일한가?)
   ```

3. **결정론적 증명**
   - 3회 반복 컴파일
   - 각 결과물 비교
   - 100% 동일성 검증

### 성공 기준
- ✅ FV → C 변환 동작
- ✅ 기존 컴파일러 결과와 동일
- ✅ 자체호스팅 완성

---

## 📊 진행 상황

| Step | 작업 | 상태 | 줄수 |
|------|------|------|------|
| 1 | Lexer 작성 | ✅ 완료 | 430 |
| 1 | AST 정의 | ✅ 완료 | 175 |
| 1 | Parser 작성 | ✅ 완료 | 520 |
| 1 | Type Checker 작성 | ✅ 완료 | 280 |
| 1 | Code Generator 작성 | ✅ 완료 | 390 |
| 2 | Stdlib 최소화/확인 | 🟡 진행 중 | 500+ |
| 3 | Main Compiler Module | ⏳ 예정 | 100 |
| 4 | 자체호스팅 테스트 | ⏳ 예정 | - |

**전체 진행률**: 40% (1,795/4,000줄 작성 완료)

---

## 🚨 블로커 & 의존성

### 현재 블로커
1. **FV 언어 완성도**
   - Pattern matching 미구현 (parser에서 TODO)
   - Option/Result 타입 미완성
   - 재귀 호출 테스트 필요

2. **Stdlib 부족**
   - 문자열 처리 최소화 필요
   - 벡터 연산 단순화 필요
   - Map 구현 필요 (타입 시스템용)

3. **FV 런타임**
   - 현재 인터프리터 기반 (속도 문제)
   - JIT 컴파일 필요 (자체호스팅 용)

### 해결 전략
1. **Stdlib 최소화**: "just enough" 접근 (과도 최적화 금지)
2. **단계적 테스트**: 각 모듈별 단독 테스트 후 통합
3. **Fallback 계획**: 필요시 Rust 헬퍼 함수 허용 (점진적 제거)

---

## 💡 핵심 아이디어

### 왜 FV로 재작성?
1. **자체호스팅 증명**: 언어의 완성도를 보여줌
2. **자기 자신을 진화시킬 수 있음**: Singularity 파운데이션
3. **Deterministic 속성**: 같은 입력 → 항상 같은 출력

### 왜 C로 코드 생성?
1. **이식성**: 모든 플랫폼 지원
2. **성능**: C 컴파일러 최적화 활용
3. **단순성**: 복잡한 백엔드 불필요

---

## 📚 관련 파일

- 기존 Rust 구현: `/data/data/com.termux/files/home/projects/fv-lang/src/`
- FV 컴파일러: `/data/data/com.termux/files/home/projects/fv-lang/src/*.fv`
- 테스트: `/data/data/com.termux/files/home/projects/fv-lang/tests/`

---

## 🎯 다음 마일스톤

- **Week 1 End**: Step 1-2 완료, Stdlib 정리
- **Week 2 Mid**: Main compiler module 완성
- **Week 2 End**: 자체호스팅 1차 성공 또는 블로커 식별

---

## 🧠 학습 포인트

1. **Bootstrap Problem**: 언어 설계와 구현의 닭과 계란 문제
   - 해결책: 먼저 다른 언어(Rust)로 구현 → FV로 포팅 → 자체호스팅

2. **Determinism**: 프로그램의 결정론적 실행
   - 중요: 부동소수점, 무작위, I/O 없이 항상 같은 결과

3. **Compiler Architecture**: 4단계 파이프라인
   - Lexer → Parser → Type Checker → Code Generator
   - 각 단계는 독립적이고 테스트 가능

---

**상태**: 🟡 Phase 1 진행 중 (Step 1 완료, Step 2 시작)
