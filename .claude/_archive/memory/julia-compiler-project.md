# 🔨 Julia 컴파일러 구현 프로젝트 (JuliaCC)

**프로젝트 시작**: 2026-03-11
**상태**: 📋 계획 수립 중
**목표**: Julia 호환 고성능 컴파일러 완전 구현 (Go 기반)

---

## 📊 프로젝트 개요

| 항목 | 내용 |
|------|------|
| **프로젝트명** | JuliaCC (Julia Compiler Complete) |
| **구현 언어** | Go (GoFree 경험 활용) |
| **기반 기술** | Julia 공식 컴파일러 파이프라인 |
| **총 Phase** | 10 Phase (40+ 상세 단계) |
| **타겟 완성도** | Self-Hosting (컴파일러가 자기 자신을 컴파일) |
| **예상 코드량** | 30,000+ 라인 (GoFree 기준 확대) |

---

## 🗺️ 10 Phase 로드맵

### Phase 0: 프로젝트 초기화 ⚙️
**목표**: 프로젝트 기초 구축
**상세 작업**:
- [ ] GitHub 저장소 생성
- [ ] 프로젝트 구조 설계
- [ ] Go 모듈 초기화 (go.mod)
- [ ] 의존성 정의 (LLVM 바인딩, 테스트 프레임워크)
- [ ] CI/CD 파이프라인 설정
- [ ] 문서 템플릿 작성

**예상 산출물**:
- 저장소 구조: `cmd/`, `internal/`, `pkg/`, `test/` 등
- 초기 테스트 프레임워크
- 개발 가이드

---

### Phase 1: Lexer (토큰화) 🔤
**목표**: Julia 소스코드를 토큰으로 분해
**상세 작업**:
- [ ] 키워드 정의 (function, if, for, ...)
- [ ] 연산자 정의 (+, -, *, /, ^, ...)
- [ ] 식별자 파싱
- [ ] 리터럴 파싱 (정수, 부동소수, 문자열, 특수문자)
- [ ] 주석 처리 (#, #=...=#)
- [ ] 줄 바꿈/공백 처리
- [ ] 에러 위치 추적 (LineInfo)

**테스트 케이스**:
- `x = 5` → `[IDENT(x), ASSIGN, INT(5)]`
- `sin(x) + 1.0` → `[IDENT(sin), LPAREN, IDENT(x), RPAREN, PLUS, FLOAT(1.0)]`
- `# comment` → `[COMMENT]` (버림)

**예상 코드**: 500-800줄

---

### Phase 2: Parser (구문 분석) 🌳
**목표**: 토큰 → AST (추상 구문 트리)
**상세 작업**:
- [ ] AST 노드 정의 (Expr, Stmt, Type, ...)
- [ ] 식(Expression) 파싱 (이항 연산, 단항 연산)
- [ ] 연산자 우선순위 처리
- [ ] 문(Statement) 파싱 (assignment, if/else, loops)
- [ ] 함수 정의 파싱
- [ ] 함수 호출 파싱
- [ ] 타입 주석 파싱 (::Type)
- [ ] 오류 복구 (에러 보고)

**AST 구조 예**:
```go
type FuncDef struct {
    Name     string
    Params   []Param
    RetType  Type
    Body     []Stmt
}

type BinOp struct {
    Left  Expr
    Op    string
    Right Expr
}
```

**예상 코드**: 1,200-1,500줄

---

### Phase 3: 타입 시스템 🎯
**목표**: Julia의 다중 디스패치 타입 시스템 구현
**상세 작업**:
- [ ] 기본 타입 정의 (Int64, Float64, String, Bool, ...)
- [ ] 복합 타입 (Vector{T}, Matrix{T}, ...)
- [ ] 추상 타입 (Number, Real, ...)
- [ ] 함수 서명 (Method Signature)
- [ ] 타입 계층 구조 (subtyping rules)
- [ ] 타입 매개변수화 (Parametric Types)
- [ ] Union 타입
- [ ] 타입 변환 규칙

**다중 디스패치 테이블**:
```
Function: multiply
- Method 1: multiply(::Int64, ::Int64) -> Int64
- Method 2: multiply(::Float64, ::Float64) -> Float64
- Method 3: multiply(::String, ::Int64) -> String
```

**예상 코드**: 1,000-1,300줄

---

### Phase 4: Semantic Analysis 🔍
**목표**: 의미 검증 및 심볼 테이블 관리
**상세 작업**:
- [ ] 심볼 테이블 (변수, 함수, 타입 등록)
- [ ] 스코프 관리 (global, local, nested)
- [ ] 변수 참조 검증 (undefined variable 감지)
- [ ] 함수 호출 검증
- [ ] 타입 호환성 검사
- [ ] 미사용 변수 경고
- [ ] 재정의 감지

**데이터 구조**:
```go
type Scope struct {
    Parent    *Scope
    Symbols   map[string]*Symbol
    ScopeType ScopeType // Global, Local, Function
}

type Symbol struct {
    Name    string
    Type    Type
    Kind    SymbolKind // Var, Func, Type, Const
    Defined bool
}
```

**예상 코드**: 800-1,000줄

---

### Phase 5: IR 생성 (Lowering) 📍
**목표**: AST → Intermediate Representation (SSA IR)
**상세 작업**:
- [ ] SSA IR 포맷 정의
- [ ] 기본 블록 생성 (BasicBlock)
- [ ] 제어흐름 그래프 (CFG) 구축
- [ ] 로워링: AST → goto 기반 IR
- [ ] Phi 노드 생성 (제어흐름 합류)
- [ ] 임시 변수 생성 (%1, %2, ...)
- [ ] 초기 최적화 (상수 폴딩, 데드 코드 제거)

**SSA IR 예**:
```
Function main(x: Float64):
  %1 = call sin(%0)         ; sin(x)
  %2 = fcmp gt %0, 5.0      ; x > 5.0
  br i1 %2, label %then, label %else
then:
  %3 = call cos(%0)         ; cos(x)
  %4 = fadd %1, %3          ; y + cos(x)
  br label %end
else:
  br label %end
end:
  %5 = phi [%4, %then], [%1, %else]  ; Phi 노드
  ret %5
```

**예상 코드**: 1,500-2,000줄

---

### Phase 6: 타입 추론 엔진 🧠
**목표**: 모든 변수/표현식의 타입 결정
**상세 작업**:
- [ ] 타입 추론 알고리즘 (Type Inference Algorithm)
- [ ] 제약 수집 (Constraint Collection)
- [ ] 제약 해결 (Unification)
- [ ] 타입 환경 관리
- [ ] 제너릭 함수 특화 (Specialization)
- [ ] 타입 오류 보고

**추론 규칙**:
```
Given: x::Int64, y::Int64
Find:  typeof(x + y)

Rule: Int64 + Int64 → Int64
Result: typeof(x + y) = Int64
```

**예상 코드**: 1,800-2,200줄

---

### Phase 7: LLVM 코드 생성 ⚡
**목표**: Typed IR → LLVM IR
**상세 작업**:
- [ ] LLVM Go 바인딩 설정 (llvm.org/bindings/go)
- [ ] LLVM IR 생성
- [ ] 함수 정의 생성
- [ ] 기본 블록 생성
- [ ] 명령어 생성 (add, mul, call, br 등)
- [ ] 디버그 정보 추가
- [ ] 모듈 검증 (verify)
- [ ] 최적화 패스 적용

**LLVM IR 생성 예**:
```llvm
define double @julia_sin_x(double %x) {
  %1 = call double @sin(double %x)
  ret double %1
}
```

**예상 코드**: 2,000-2,500줄

---

### Phase 8: 런타임 & VM 🖥️
**목표**: 기본 실행 엔진 구현
**상세 작업**:
- [ ] 기본 VM 구현 (또는 JIT 실행)
- [ ] 메모리 관리 (heap, stack)
- [ ] 가비지 컬렉션 (GC) — Mark & Sweep
- [ ] 함수 호출 스택
- [ ] 내장 함수 구현 (sin, cos, sqrt, ...)
- [ ] 배열 런타임 지원
- [ ] 문자열 런타임 지원

**메모리 구조**:
```
Stack (함수 호출, 로컬 변수)
↓
Heap (객체, 배열 할당)
├─ GC 루트들
└─ Mark-Sweep 추적
```

**예상 코드**: 2,500-3,000줄

---

### Phase 9: 최적화 엔진 🔥
**목표**: 고성능 코드 생성
**상세 작업**:
- [ ] 인라이닝 (Inlining) — 작은 함수 직접 삽입
- [ ] 루프 최적화 (Loop Optimization) — 루프 언롤링
- [ ] 스칼라화 (Scalarization) — 벡터화 준비
- [ ] 조건부 제거 (Branch Elimination)
- [ ] SIMD 생성 (SIMD Instructions)
- [ ] 메모리 접근 최적화
- [ ] 레지스터 할당

**최적화 전후 비교**:
```
Before: for i in 1:1000000 s += v[i]  (루프)
After:  %1 = SIMD_ADD %v             (SIMD 가속)
```

**예상 코드**: 2,000-2,500줄

---

### Phase 10: Self-Hosting & 통합 🎯
**목표**: 컴파일러로 컴파일러를 컴파일하기
**상세 작업**:
- [ ] 컴파일러를 Go에서 Julia로 포팅 (선택)
- [ ] Julia 코드로 컴파일러 재구현
- [ ] 컴파일러로 자신의 소스 컴파일
- [ ] 부트스트랩 프로세스 (bootstrap sequence)
- [ ] 성능 벤치마킹
- [ ] 회귀 테스트
- [ ] 문서 완성
- [ ] 배포 패키징

**Self-Hosting 성공 기준**:
```
jcc.go (Go 컴파일러)
    ↓ go build
jcc (실행파일)
    ↓ ./jcc compiler.jl (Julia로 된 컴파일러 컴파일)
jcc-julia (Julia로 만든 컴파일러)
    ↓ 자신의 소스 컴파일 가능
✅ Self-hosting 달성!
```

**예상 코드**: 통합 및 통일화

---

## 📈 상세 구현 계획

### 주간 목표 (예상 일정)

| 주차 | Phase | 목표 | 산출물 |
|------|-------|------|--------|
| 1 | Phase 0 | 프로젝트 기초 설정 | 저장소 + 구조 |
| 2-3 | Phase 1 | Lexer 완성 | 토큰화 엔진 |
| 4-5 | Phase 2 | Parser 완성 | AST 생성기 |
| 6 | Phase 3 | 타입 시스템 | Multiple Dispatch |
| 7 | Phase 4 | Semantic Analysis | 심볼 테이블 |
| 8-9 | Phase 5 | IR 생성 | SSA IR 엔진 |
| 10 | Phase 6 | 타입 추론 | Type Inference |
| 11-12 | Phase 7 | LLVM Codegen | LLVM IR 생성 |
| 13 | Phase 8 | 런타임 | VM 구현 |
| 14 | Phase 9 | 최적화 | 최적화 패스 |
| 15-16 | Phase 10 | 통합 & Self-Hosting | 완성 |

---

## 🛠️ 기술 스택

| 계층 | 기술 |
|------|------|
| **언어** | Go 1.21+ |
| **중간표현** | Julia SSA IR + LLVM IR |
| **LLVM** | llvm-c binding (v14+) |
| **런타임** | 커스텀 VM + GC |
| **테스트** | Go testing + Julia compatibility |
| **문서** | Markdown + code comments |

---

## 📦 저장소 구조 (초안)

```
julia-compiler/
├── cmd/
│   └── jcc/                  # 메인 컴파일러 바이너리
│       └── main.go
├── internal/
│   ├── lexer/                # 토큰화
│   │   ├── lexer.go
│   │   └── token.go
│   ├── parser/               # 파싱
│   │   ├── parser.go
│   │   └── ast.go
│   ├── types/                # 타입 시스템
│   │   ├── types.go
│   │   └── dispatch.go
│   ├── semantic/             # 의미 분석
│   │   ├── analyzer.go
│   │   └── symbols.go
│   ├── ir/                   # IR 생성
│   │   ├── ir.go
│   │   ├── ssa.go
│   │   └── optimizer.go
│   ├── codegen/              # LLVM Codegen
│   │   ├── llvm.go
│   │   └── codegen.go
│   ├── runtime/              # 런타임
│   │   ├── vm.go
│   │   └── gc.go
│   └── optimizer/            # 최적화
│       ├── inline.go
│       ├── loop.go
│       └── simd.go
├── pkg/
│   └── stdlib/               # 표준 라이브러리
│       ├── math.go
│       ├── array.go
│       └── string.go
├── test/
│   ├── lexer_test.go
│   ├── parser_test.go
│   ├── types_test.go
│   ├── codegen_test.go
│   └── integration_test.go
├── examples/
│   ├── hello.jl
│   ├── fibonacci.jl
│   └── matrix_mul.jl
├── docs/
│   ├── ARCHITECTURE.md       # 아키텍처
│   ├── PHASES.md             # Phase별 가이드
│   ├── TESTING.md            # 테스트 전략
│   └── OPTIMIZATION.md       # 최적화 가이드
├── Makefile
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## 🧪 테스트 전략

### 단계별 테스트

| Phase | 테스트 유형 | 검증 항목 |
|-------|-----------|---------|
| 1 | 유닛 테스트 | 토큰 정확성 |
| 2 | 파서 테스트 | AST 구조 |
| 3 | 타입 테스트 | 타입 추론 정확성 |
| 4 | 의미 테스트 | 심볼 검증 |
| 5 | IR 테스트 | SSA 형식 |
| 6 | 타입 추론 테스트 | 모든 타입 추론 |
| 7 | LLVM 테스트 | 생성된 IR 검증 |
| 8 | 런타임 테스트 | 실행 결과 검증 |
| 9 | 성능 테스트 | 최적화 효율 |
| 10 | 통합 테스트 | Self-hosting 검증 |

### 테스트 케이스 예

```go
// Phase 1: Lexer Test
func TestLexer_BasicTokens(t *testing.T) {
    input := "x = 5 + 3.14"
    lexer := NewLexer(input)
    tokens := lexer.Tokenize()
    assert.Equal(t, tokens[0].Type, IDENT)
    assert.Equal(t, tokens[0].Value, "x")
    // ...
}

// Phase 7: Codegen Test
func TestCodegen_SimpleFunctionCall(t *testing.T) {
    ast := parseCode("sin(1.0)")
    llvmIR := codegenAST(ast)
    assert.Contains(t, llvmIR, "call double @sin")
}
```

---

## 🎯 성공 기준

| 기준 | 목표 | 검증 |
|------|------|------|
| **정확성** | 모든 Julia 기본 문법 지원 | 호환성 테스트 |
| **성능** | Python 대비 50배 이상 | 벤치마크 비교 |
| **완성도** | Self-hosting 달성 | 컴파일러 자체 컴파일 |
| **코드 품질** | 80%+ 테스트 커버리지 | go test -cover |
| **문서화** | 모든 함수 주석 + 가이드 | godoc |

---

## 💾 커밋 전략

각 Phase 완료 시 다음 형식으로 커밋:

```
[Phase X] 제목 (설명)

- 구현한 기능 1
- 구현한 기능 2
- ...

테스트: X개 PASS
코드: XXX줄 추가
```

예:
```
[Phase 1] Lexer 구현 완료

- 키워드/연산자 토큰화
- 숫자/문자열 파싱
- 주석 처리
- 에러 위치 추적

테스트: 15개 PASS
코드: 650줄 추가
```

---

## 🚀 다음 단계

1. **저장소 생성**: GitHub에 `julia-compiler` 생성
2. **Phase 0 시작**: 프로젝트 기초 설정
3. **커뮤니티**: Julia 공식 포럼에 진행 공유

---

**상태**: 📋 계획 완료 → 🚀 Phase 0 시작 준비
**마지막 업데이트**: 2026-03-11
