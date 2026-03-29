# FreeLang V3: 컴파일러 파이프라인을 처음부터 만들어봤어요

## 요약
프로그래밍 언어의 가장 복잡한 부분이 뭘까요? 바로 **컴파일러**입니다. 우리가 FreeLang V3에서 처음부터 만든 완전한 컴파일러 파이프라인을 공개합니다. Lexer, Parser, Type System, Code Generation까지 3,000줄의 코드로 어떻게 동작하는지 살펴봐요.

---

## 들어가기

코드를 작성하면, 그것이 어떻게 실행 가능한 프로그램이 될까요?

```
Source Code → Lexer → Tokens → Parser → AST → Type Checker → Code Gen → Binary
```

이 과정이 **컴파일 파이프라인**입니다. 겉보기는 간단해 보이지만, 각 단계마다 정말 중요한 결정들이 일어나요.

## Phase 1: Lexer (토큰화)

첫 번째는 **Lexer**입니다. 이건 소스 코드를 읽고 의미 있는 단위(토큰)로 쪼개는 역할을 합니다.

```go
// 입력: let x = 42;
// 출력: [KEYWORD(let), IDENT(x), ASSIGN(=), NUMBER(42), SEMICOLON(;)]
```

**왜 이 단계가 필요할까요?**
- 공백, 주석 제거
- 키워드와 식별자 구분
- 문법 오류 조기 감지

우리는 30개의 토큰 타입을 정의하고, 위치 정보(line, col)까지 추적합니다. 나중에 에러 메시지를 정확하게 보여줄 수 있게요.

## Phase 2: Parser (문법 분석)

다음은 **Parser**입니다. 토큰들을 규칙에 따라 조합하는 거죠.

예를 들어, 변수 선언은 이렇게 됩니다:

```
let x: int = 42;
    ↓ ↓  ↓   ↓ ↓
    1 2  3   4 5

1: KEYWORD(let)
2: IDENTIFIER(x)
3: TYPE_ANNOTATION(int)
4: VALUE(42)
5: SEMICOLON
```

우리는 **Pratt Parser**를 사용했습니다. 이건 연산자의 우선순위를 깔끔하게 처리할 수 있는 방식이에요.

```go
// 1 + 2 * 3 을 올바르게 파싱
// AST: Add(1, Mul(2, 3)) ✅
// (1 + 2) * 3 아닙니다! ❌
```

## Phase 3: Type System (타입 검증)

이제 AST(Abstract Syntax Tree)가 생겼습니다. 다음은 타입이 맞는지 확인하는 거예요.

```go
let x: int = "hello";  // ❌ Error: string을 int로 할당 불가
let y: int = 42;       // ✅ OK
```

우리는 이 과정에서:
- 변수 타입 추론
- 함수 반환값 타입 확인
- 제네릭 타입 치환
- 타입 캐스팅 검증

을 합니다.

## Phase 4: Code Generation (최종 생성)

마지막으로 **Code Gen** 단계입니다. 검증된 AST를 실행 가능한 형태로 변환합니다.

우리는 두 가지 방식을 지원합니다:

### 1. Go 코드로 변환
```go
// FreeLang 원본
let sum = fn(a, b) { a + b }
sum(3, 4)

// 생성된 Go 코드
func sum(a int, b int) int {
    return a + b
}
sum(3, 4)
```

### 2. 바이트코드로 변환
```
LOAD_VAR a
LOAD_VAR b
ADD
RETURN
```

**어느 게 더 빠를까요?**
- Go 변환: 컴파일 오버헤드 있지만, 최적화 가능
- 바이트코드: 빠른 시작, 하지만 해석 오버헤드

우리는 상황에 맞춰 선택할 수 있게 했습니다.

## 실제 예제

간단한 계산기를 만들어봅시다:

```freelang
fn factorial(n: int) -> int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

let result = factorial(5)
print(result)  // 120
```

이 코드가 어떻게 변환되는지 보면:

```
1. Lexer: [FN, IDENT, LPAREN, ... ]
2. Parser:
   FunctionDef {
     name: "factorial",
     params: [{ name: "n", type: "int" }],
     returnType: "int",
     body: { ... }
   }
3. TypeChecker: ✅ 모든 타입 검증 완료
4. CodeGen (Go):
   func factorial(n int) int {
     if n <= 1 {
       return 1
     }
     return n * factorial(n - 1)
   }
```

## 성능 테스트

우리는 다양한 크기의 프로그램으로 테스트했습니다:

| 프로그램 | 코드 | Lexer | Parser | TypeCheck | CodeGen | 총합 |
|---------|------|-------|--------|-----------|---------|------|
| 간단 (50줄) | 50줄 | 0.1ms | 0.2ms | 0.05ms | 0.1ms | 0.45ms |
| 중간 (500줄) | 500줄 | 1ms | 2ms | 0.5ms | 1ms | 4.5ms |
| 복잡 (5000줄) | 5,000줄 | 10ms | 20ms | 5ms | 10ms | 45ms |

**결론**: 대부분의 시간은 Parser에서 소비됩니다. 다음 버전에서 최적화할 계획입니다.

## 배운 점

1. **설계가 중요합니다**
   - Lexer, Parser, TypeChecker의 책임을 명확히 나누세요
   - 각 단계의 출력이 다음 단계의 입력이 되게

2. **에러 메시지는 최우선**
   - 위치 정보를 정확하게 기록하세요
   - 개발자가 쉽게 문제를 찾을 수 있어야 합니다

3. **테스트는 매 단계마다**
   - Lexer 단위 테스트
   - Parser 기능 테스트
   - 통합 테스트

   우리는 100+ 테스트로 안정성을 보장합니다.

## 다음은?

지금 이 컴파일러는:
- ✅ 정적 타입 검증
- ✅ 함수 정의 및 호출
- ✅ 제어 흐름 (if, while, for)
- ⏳ 클래스/구조체 (개발 중)
- ⏳ 모듈 시스템 (다음 포스트!)

---

## 마치며

컴파일러를 처음부터 만드는 건 정말 어렵습니다. 하지만 이 과정에서 우리가 매일 사용하는 프로그래밍 언어가 얼마나 정교하게 설계되었는지 배울 수 있어요.

FreeLang의 컴파일러 코드는 GitHub에서 완전히 공개되어 있습니다. 관심 있으신 분들은 살펴봐주세요!

**궁금한 점이 있으신가요?** 댓글로 물어봐주세요. 어떤 단계를 더 자세히 설명해달라는 요청도 환영합니다.
