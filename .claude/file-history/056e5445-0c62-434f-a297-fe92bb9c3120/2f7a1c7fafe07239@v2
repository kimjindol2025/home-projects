# Julia 컴파일러 아키텍처 분석

## 1. 컴파일러 파이프라인 개요

Julia 컴파일러는 **다단계 JIT 컴파일러**입니다:

```
소스코드 (*.jl)
    ↓
┌─────────────────────────────────────┐
│ 1. LEXER (토크나이저)               │
│  - 문자 스트림 → 토큰 스트림         │
│  - 89개 토큰 타입 인식               │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 2. PARSER (파서)                    │
│  - 토큰 스트림 → AST                │
│  - Recursive Descent Parser         │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 3. LOWERING (로워링)                │
│  - AST → Untyped IR                 │
│  - goto 기반 제어흐름                │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 4. TYPE INFERENCE (타입 추론)       │
│  - Untyped IR → Typed SSA IR        │
│  - 모든 변수 타입 확정               │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 5. OPTIMIZATION (최적화)             │
│  - 인라이닝                          │
│  - 상수 전파 (Constant Propagation) │
│  - 불필요 코드 제거                  │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 6. CODEGEN (코드 생성)              │
│  - Typed IR → LLVM IR               │
│  - 플랫폼 독립적 중간 표현            │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 7. LLVM 백엔드                      │
│  - LLVM IR → 네이티브 어셈블리       │
│  - 아키텍처별 최적화                  │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ 8. 링킹 & 실행                      │
│  - RuntimeDyld 링킹                 │
│  - 메모리 로드 → CPU 실행            │
└─────────────────────────────────────┘

기계어 (x86-64, ARM64 등)
```

## 2. 주요 IR (중간 표현) 형식

### 2.1 Untyped IR (로워링 후)
- 제어흐름: `goto`, `label`, `conditional_branch`
- 연산: `call`, `load`, `store`, `binary_op`
- 특징: 타입 정보 없음, 정규화된 제어흐름

### 2.2 Typed SSA IR (타입 추론 후)
- **SSA (Static Single Assignment)**: 각 변수는 정확히 한 번만 할당
- **Phi 노드**: 제어흐름 합류점에서 여러 값 합치기
- 모든 변수에 구체적 타입 정보 포함

예:
```
# Julia 코드
if condition
    x = 1.0
else
    x = 2.0
end
y = x + 1.0

# SSA 형식
label_true:
  x₁ = 1.0
  goto merge
label_false:
  x₂ = 2.0
  goto merge
merge:
  x₃ = φ(x₁, x₂)          # Phi 노드
  y = x₃ + 1.0
```

## 3. 다중 디스패치 처리

Julia의 핵심 특성: 동일 함수명의 여러 구현

```julia
function process(x::Int)
    return x * 2
end

function process(x::Float64)
    return x * 2.5
end

function process(x::String)
    return x * 3
end
```

컴파일러는:
1. **호출 시점의 인자 타입 분석**
2. **매칭되는 메서드 선택**
3. **해당 메서드로 특화된(specialized) 코드 생성**
4. **결과 캐싱** (다시 호출할 때 재사용)

## 4. Julia 타입 시스템

```
Any
├── Type{T}
├── Union{A, B}
├── Tuple{T...}
├── Function
├── Number
│   ├── Integer
│   │   ├── Int8, Int16, Int32, Int64
│   │   ├── UInt8, UInt16, UInt32, UInt64
│   │   └── BigInt
│   ├── AbstractFloat
│   │   ├── Float16, Float32, Float64
│   │   └── BigFloat
│   └── Complex{T}
├── AbstractArray{T, N}
│   ├── Array{T, N}
│   ├── Matrix{T}  (== Array{T, 2})
│   └── Vector{T}  (== Array{T, 1})
├── String
├── Symbol
└── ...
```

## 5. 토큰 타입 (89개)

### 키워드 (35개)
- `function`, `return`, `if`, `else`, `elseif`, `end`
- `for`, `while`, `break`, `continue`
- `begin`, `let`, `do`
- `const`, `global`, `local`
- `try`, `catch`, `finally`
- `import`, `using`, `module`
- `export`, `include`
- `abstract`, `struct`, `mutable`
- `macro`, `quote`
- `true`, `false`, `nothing`

### 연산자 & 구분자 (30개)
- 산술: `+`, `-`, `*`, `/`, `%`, `//`, `^`
- 비교: `==`, `!=`, `<`, `<=`, `>`, `>=`
- 논리: `&&`, `||`, `!`
- 비트: `&`, `|`, `~`, `<<`, `>>`
- 할당: `=`, `+=`, `-=`, `*=`, `/=`
- 기타: `(`, `)`, `[`, `]`, `{`, `}`, `,`, `;`, `:`, `::`

### 리터럴 (6개)
- `INTEGER` (정수)
- `FLOAT` (실수)
- `STRING` (문자열)
- `CHAR` (문자)
- `SYMBOL` (`:symbol`)
- `IDENTIFIER` (변수명)

## 6. AST 노드 타입 (40+)

```
Expression
├── Literal (Int, Float, String, Symbol)
├── Identifier
├── BinaryOp (op, left, right)
├── UnaryOp (op, operand)
├── Call (func, args)
├── IndexAccess (object, indices)
├── MemberAccess (object, field)
├── TypeAnnotation (value, type)
├── Assignment (target, value)
├── FunctionDef (name, params, return_type, body)
├── If (condition, then_body, else_body)
├── While (condition, body)
├── For (variable, iterator, body)
├── Try (body, catch_clause, finally_clause)
├── Block (statements)
├── Return (value)
├── ArrayLiteral (elements)
├── TupleLiteral (elements)
└── ...
```

## 7. 자체 호스팅 부트스트랩

자체 호스팅 컴파일러의 최종 목표:

```
Julia 컴파일러 (Julia로 작성)
  ↓
Julia 코드 입력
  ↓
[위 8단계 파이프라인]
  ↓
Julia 코드 출력 (또는 기계어)
```

**부트스트랩 과정**:
1. 처음: Julia 컴파일러 v1.10 사용 (외부 컴파일러)
2. 우리의 Julia 컴파일러 구현 → Julia v1.10으로 컴파일
3. 결과물: 우리의 컴파일러 실행 파일
4. 최종: 우리 컴파일러로 자기 자신을 컴파일

---

## 참고

- [Julia Docs: Devdocs - JIT Compiler](https://docs.julialang.org/en/v1/devdocs/jit/)
- [Julia Source: src/julia.h](https://github.com/JuliaLang/julia/blob/master/src/julia.h)
- [Julia IR Docs](https://docs.julialang.org/en/v1/devdocs/ssair/)
