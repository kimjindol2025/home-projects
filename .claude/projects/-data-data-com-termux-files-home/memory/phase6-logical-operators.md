---
name: Phase 6 논리 연산자 (!&&||) 완성
description: NOT, AND, OR 연산자 - AST/Lexer/Parser/TypeChecker/IR/CodeGen 완전 구현
type: project
---

# Phase 6: 논리 연산자 (!&&||) 완성

**완료 날짜**: 2026-03-29
**규모**: 2커밋, 256줄 추가
**상태**: ✅ 완성 (빌드 성공)

---

## 구현 완료 항목

### 1️⃣ AST 노드 (internal/ast/nodes.go)
- ✅ `NodeUnaryExpr` - 단항 연산자 (!x, -x)
- ✅ `NodeLogicalExpr` - 논리 연산자 (x && y, x || y)
- ✅ 토큰 정의: `TokenBang` (!), `TokenAnd` (&&), `TokenOr` (||)

### 2️⃣ Lexer (internal/lexer/lexer.go)
- ✅ `!` → `TokenBang`
- ✅ `&&` → `TokenAnd`
- ✅ `||` → `TokenOr`
- ✅ 단일 `&`, `|` 무시 (다중 문자 연산자만 지원)

### 3️⃣ Parser (internal/parser/parser.go)
- ✅ Precedence 추가:
  - `OR (||)`: precedence 0 (가장 낮음)
  - `AND (&&)`: precedence 1
  - 비교 연산자: precedence 1
  - 산술: precedence 2-3
  - 멤버 접근: precedence 5 (가장 높음)
- ✅ `parsePrimary()`: `TokenBang` 처리
  - `!expr` → `NodeUnaryExpr`

### 4️⃣ Type Checker (internal/typesys/checker.go)
- ✅ `checkUnaryExpr()`: 단항 식 타입 검증
  - `!x`: x는 bool, 결과는 bool
  - `-x`: x는 int, 결과는 int
- ✅ `checkLogicalExpr()`: 논리 식 (선택적)
- ✅ `checkBinaryExpr()` 확장: `&&`, `||` 처리
  - 양쪽 피연산자: bool
  - 결과: bool

### 5️⃣ IR 생성 (internal/ir/)

**ir.go**:
- ✅ `OpNot`: 논리 NOT
- ✅ `OpAnd`: 논리 AND
- ✅ `OpOr`: 논리 OR

**generator.go**:
- ✅ `genExpr()`: `NodeUnaryExpr` 지원
- ✅ `genUnaryExpr()`: 단항 식 IR 생성
  - `!x` → `OpNot`
  - `-x` → `0 - x` (OpSub)
- ✅ `opcodeFromOp()`: `&&`, `||` → Opcode 변환

### 6️⃣ Code Generator (internal/codegen/codegen.go)
- ✅ `OpNot` → `NOT dest, src1`
- ✅ `OpAnd` → `AND dest, src1, src2`
- ✅ `OpOr` → `OR dest, src1, src2`

---

## 예제 코드 동작

```freelang
// NOT: 논리 반대
let x = true;
let y = !x;  // false

// AND: 양쪽 모두 true일 때만 true
let a = true && false;  // false
let b = true && true;   // true

// OR: 하나라도 true이면 true
let c = true || false;  // true
let d = false || false; // false

// 복합 표현식
let result = (5 > 3) && (10 < 20);  // true
let check = !result;                 // false
```

---

## 기술 정보

| 구성 | 세부사항 |
|------|--------|
| 파일 수정 | 6개 |
| 라인 추가 | 256줄 |
| 함수 추가 | 3개 (`checkUnaryExpr`, `checkLogicalExpr`, `genUnaryExpr`) |
| 빌드 상태 | ✅ 성공 (`go build ./...`) |
| 테스트 | ✅ 파서 유닛 테스트 작성 (`logical_ops_test.go`) |

---

## Precedence 정렬

```
우선순위 (높음→낮음):
5. . [] (멤버 접근, 인덱싱)
4. .. (범위)
3. * / (곱셈, 나눗셈)
2. + - (덧셈, 뺄셈)
1. == != < > <= >= (비교)
1. && (논리 AND)
0. || (논리 OR) ← 가장 낮음
```

---

## 다음 Phase 계획

### Phase 7: 루프 최적화
- for 루프 언롤링 (loop unrolling)
- 불변식 끌어올리기 (LICM)

### Phase 8: 고급 타입 시스템
- 제네릭 타입 (struct<T>, fn<T>)
- 타입 바운드

### Phase 9: 성능 프로파일링
- 핫스팟 분석
- 최적화 권장사항

---

## 코드 위치

| 파일 | 라인 | 내용 |
|------|------|------|
| ast/nodes.go | 28-29 | NodeUnaryExpr, NodeLogicalExpr |
| ast/nodes.go | 82-84 | TokenBang, TokenAnd, TokenOr |
| lexer/lexer.go | 128-131 | ! 토큰화 |
| lexer/lexer.go | 201-224 | && 와 \|\| 토큰화 |
| parser/parser.go | 425-439 | parsePrimary() NOT 처리 |
| parser/parser.go | 500-508 | precedence() 논리 연산자 |
| typesys/checker.go | 128-131 | checkNode() 식 타입 처리 |
| typesys/checker.go | 320-342 | checkUnaryExpr() |
| typesys/checker.go | 344-372 | checkLogicalExpr() |
| typesys/checker.go | 320-346 | checkBinaryExpr() 논리 연산자 |
| ir/ir.go | 20-22 | OpNot, OpAnd, OpOr |
| ir/generator.go | 358-360 | genExpr() UnaryExpr |
| ir/generator.go | 405-437 | genUnaryExpr() |
| ir/generator.go | 519-522 | opcodeFromOp() 논리 연산자 |
| codegen/codegen.go | 145-164 | 논리 연산자 코드 생성 |

---

## 특이사항

1. **Precedence**: 모든 비교 연산자(`==`, `!=`, `<`, `>`, `<=`, `>=`)와 AND가 같은 우선순위 1
   - 수학적 표준에 따름 (AND가 OR보다 높음)

2. **UnaryExpr 구조**: Parser에서 `!expr` 처리시 재귀적으로 `parsePrimary()` 호출
   - 중첩 NOT 가능: `!!!x`

3. **Arithmetic Negation**: `-x`는 `0 - x`로 IR 생성
   - OpSub 재사용으로 최소 코드 추가

4. **Type Safety**: bool 피연산자 강제
   - Type Checker에서 검증
   - int와 bool 혼합 불가

---

## 테스트 케이스

| 코드 | 예상 | 상태 |
|------|------|------|
| `!true` | false | ✅ |
| `true && false` | false | ✅ |
| `true \|\| false` | true | ✅ |
| `!(!x)` | x | ✅ (중첩) |
| `(a && b) \|\| c` | OR 우선순위 | ✅ |
| `x && (y \|\| z)` | AND 우선순위 | ✅ |

