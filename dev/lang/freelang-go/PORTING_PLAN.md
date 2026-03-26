# FreeLang Go 포팅 계획서

**기준 버전**: FreeLang v4 (TypeScript)
**목표**: 동일한 213개 테스트 통과
**타겟 기간**: 2-3주
**언어**: Go (1.21+)

---

## 📊 자원 분석

### v4 TypeScript 구현
- **크기**: 77MB
- **파일**: 6,005개
- **테스트**: 213개 (모두 PASS ✅)
- **최종 수정**: 2026-03-26

### v4 핵심 모듈
```
src/
├── lexer.ts          → 토큰화 (Lexer)
├── parser.ts         → AST 생성 (Parser)
├── checker.ts        → 타입 검사 (Type Checker)
├── compiler.ts       → 바이트코드 생성 (Compiler)
├── vm.ts             → 가상 머신 (Runtime)
├── ast.ts            → AST 정의
└── db-*.ts           → DB 벤치마크 (선택사항)
```

### 테스트 구성
```
19개 테스트 파일:
├── lexer.test.ts           (토큰화)
├── parser.test.ts          (파싱)
├── checker.test.ts         (타입 검사)
├── compiler.test.ts        (컴파일)
├── vm.test.ts              (실행)
├── struct.test.ts          (구조체)
├── function-literal.test.ts (함수 리터럴)
├── while-loop.test.ts      (루프)
├── for-of.test.ts          (반복)
└── ... 10개 더 (총 213개 테스트)
```

---

## 🎯 Go 포팅 로드맵

### Phase 1: 핵심 구조 구현 (3-4일)

#### 1.1 Lexer (토큰화)
```go
// cmd/lexer/main.go
type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           byte
}

// 메서드:
// - NextToken() Token
// - readNumber() string
// - readIdentifier() string
// - skipWhitespace()
```

**TypeScript 원본 분석**:
- 정규표현식 기반 토큰 분류
- 44개 토큰 타입
- 들여쓰기 기반 INDENT/DEDENT

**Go 구현 목표**:
- 상태 머신 기반 (정규식 불필요)
- 동일한 44개 토큰 타입 정의
- INDENT/DEDENT 처리

#### 1.2 AST (추상 문법 트리)
```go
// pkg/ast/ast.go
type Node interface {
    node()
}

type Expression interface {
    Node
    exprNode()
}

type Statement interface {
    Node
    stmtNode()
}

// 구체적 노드들:
type Identifier struct { Value string }
type IntegerLiteral struct { Value int64 }
type StringLiteral struct { Value string }
type ArrayLiteral struct { Elements []Expression }
type HashLiteral struct { Pairs map[string]Expression }
type FunctionLiteral struct { Parameters []*Identifier; Body *BlockStatement }
type BlockStatement struct { Statements []Statement }
type IfExpression struct { Condition Expression; Consequence *BlockStatement; Alternative *BlockStatement }
// ... 더 많음
```

**Go 구현 목표**:
- interface{} 사용 최소화
- 타입 안전성 보장
- 리플렉션 활용 (선택사항)

#### 1.3 Parser (파싱)
```go
// pkg/parser/parser.go
type Parser struct {
    l              *lexer.Lexer
    curToken       token.Token
    peekToken      token.Token
    errors         []string
    infixParseFns  map[token.TokenType]infixParseFn
    prefixParseFns map[token.TokenType]prefixParseFn
}

// 메서드:
// - ParseProgram() *ast.Program
// - parseStatement() ast.Statement
// - parseExpression() ast.Expression
// - nextToken()
// - curTokenIs() bool
// - peekTokenIs() bool
```

**Go 구현 목표**:
- Pratt Parser 구현 (우선순위 파싱)
- 재귀 하강 파싱
- 에러 복구 메커니즘

### Phase 2: 타입 시스템 (3-4일)

#### 2.1 Type Checker
```go
// pkg/checker/checker.go
type Checker struct {
    symbolTable *SymbolTable
    errors      []string
}

type SymbolTable struct {
    store          map[string]*Symbol
    outer          *SymbolTable
    builtins       map[string]*Symbol
}

type Symbol struct {
    Name  string
    Type  Type
    Scope ScopeType
}

// 메서드:
// - Check(program *ast.Program) *CheckResult
// - checkStatement(stmt ast.Statement) Type
// - checkExpression(expr ast.Expression) Type
// - Define(name string, typ Type)
// - Resolve(name string) (*Symbol, bool)
```

**Go 구현 목표**:
- TypeScript의 타입 체커 로직 이식
- 심볼 테이블 관리
- 타입 추론 (type inference)

#### 2.2 Type 정의
```go
// pkg/types/types.go
type Type interface {
    TypeName() string
}

type IntType struct{}
type StringType struct{}
type BoolType struct{}
type ArrayType struct{ ElementType Type }
type HashType struct{ KeyType, ValueType Type }
type FunctionType struct{ ParameterTypes []Type; ReturnType Type }
type StructType struct{ Fields map[string]Type }
// ... 더 많음
```

### Phase 3: 컴파일러 (3-4일)

#### 3.1 Bytecode 생성
```go
// pkg/compiler/compiler.go
type Compiler struct {
    instructions    code.Instructions
    lastInstruction EmittedInstruction
    previousInstruction EmittedInstruction
    symbolTable     *symbol.SymbolTable
    constants       []interface{}
}

type OpCode byte

const (
    OpConstant OpCode = iota
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpMod
    // ... 더 많음
)

// 메서드:
// - Compile(program *ast.Program) error
// - compileStatement(stmt ast.Statement) error
// - compileExpression(expr ast.Expression) error
// - Emit(op OpCode, operands ...int) int
// - addConstant(obj interface{}) int
```

**Go 구현 목표**:
- TypeScript 바이트코드 ISA 완전 구현
- 상수 폴딩
- 데드 코드 제거 (선택사항)

### Phase 4: 런타임 (4-5일)

#### 4.1 Virtual Machine
```go
// pkg/vm/vm.go
type VM struct {
    constants   []interface{}
    stack       []interface{}
    sp          int // Stack Pointer
    globals     []interface{}
}

type Object interface {
    Type() ObjectType
    Inspect() string
}

type (
    Integer struct{ Value int64 }
    String struct{ Value string }
    Boolean struct{ Value bool }
    Array struct{ Elements []Object }
    Hash struct{ Pairs map[string]Object }
    Function struct{ ... }
)

// 메서드:
// - Run() error
// - execute(bytecode []byte) error
// - push(obj Object) error
// - pop() Object
// - executeBinaryOperation(op code.OpCode, left, right Object) Object
```

**Go 구현 목표**:
- 스택 기반 VM
- 메모리 관리
- 가비지 컬렉션 (Go 자동)

### Phase 5: 테스트 (2-3일)

#### 5.1 테스트 구조
```go
// tests/
├── lexer_test.go      (19개 테스트 → 1개 Go 파일)
├── parser_test.go
├── checker_test.go
├── compiler_test.go
├── vm_test.go
├── struct_test.go
├── function_test.go
└── integration_test.go (end-to-end)
```

**Go 테스팅**:
- testing.T 표준 라이브러리
- Table-driven tests
- Subtests (`t.Run()`)

---

## 📁 Go 프로젝트 표준 구조

```
~/dev/lang/freelang-go/
├── cmd/                      # 실행 가능 프로그램
│   ├── freelang-compile/
│   │   └── main.go          # 컴파일러 CLI
│   └── freelang-run/
│       └── main.go          # 런타임 CLI
├── pkg/                      # 라이브러리 패키지
│   ├── lexer/
│   │   ├── lexer.go
│   │   └── lexer_test.go
│   ├── token/
│   │   └── token.go
│   ├── ast/
│   │   └── ast.go
│   ├── parser/
│   │   ├── parser.go
│   │   └── parser_test.go
│   ├── types/
│   │   └── types.go
│   ├── checker/
│   │   ├── checker.go
│   │   └── checker_test.go
│   ├── code/
│   │   └── code.go          # Bytecode ISA
│   ├── compiler/
│   │   ├── compiler.go
│   │   └── compiler_test.go
│   ├── vm/
│   │   ├── vm.go
│   │   └── vm_test.go
│   └── object/
│       └── object.go        # Runtime objects
├── tests/
│   ├── integration_test.go
│   └── testdata/
│       ├── valid/          # 통과해야 할 코드
│       └── invalid/        # 에러 검증
├── go.mod
├── go.sum
├── Makefile                 # 빌드 자동화
├── README.md
└── PORTING_PLAN.md (이 파일)
```

---

## 🔄 포팅 전략

### 1단계: 직역 (Direct Translation)
```
TypeScript 코드 → Go 코드 (1:1 매핑)
- 변수명, 함수명 유지
- 로직 그대로 이식
```

### 2단계: 개선 (Idiomatic Go)
```
Go best practices 적용
- interface{} → 제네릭 (Go 1.18+)
- 에러 처리 개선
- 성능 최적화
```

### 3단계: 검증 (Testing)
```
213개 테스트 모두 통과
- 단위 테스트
- 통합 테스트
- 호환성 테스트 (TypeScript와 비교)
```

---

## ⚡ 즉시 시작 가능한 작업

### Step 1: Lexer 구현 (1일)
```bash
cd ~/dev/lang/freelang-go
mkdir -p pkg/{lexer,token,ast} cmd/freelang-run
# lexer/lexer.go 작성 시작
```

### Step 2: Parser 구현 (1-2일)
```bash
# parser/parser.go 작성
# Pratt parser 패턴 구현
```

### Step 3: TypeScript 테스트케이스 분석 (0.5일)
```bash
cd ~/dev/lang/freelang-v4
npm test -- --verbose  # 213개 테스트 상세 분석
```

### Step 4: Go 테스트 작성 (1-2일)
```bash
cd ~/dev/lang/freelang-go
go test ./...
# 각 패키지별 테스트 추가
```

---

## 📈 성공 지표

| 항목 | 목표 | 진행도 |
|------|------|--------|
| Lexer | 44개 토큰 타입 완전 구현 | ⏳ |
| Parser | 모든 문법 규칙 파싱 가능 | ⏳ |
| Checker | 타입 검사 100% 일치 | ⏳ |
| Compiler | 바이트코드 생성 정확 | ⏳ |
| VM | 213개 테스트 PASS | ⏳ |
| **전체** | **Go 포팅 완료** | ⏳ |

---

## 🚀 Go 포팅 시작!

**다음 명령어로 즉시 시작:**

```bash
cd ~/dev/lang/freelang-go

# 1. 기본 구조 생성
mkdir -p pkg/{lexer,token,ast,parser,types,checker,code,compiler,vm,object}
mkdir -p cmd/{freelang-compile,freelang-run}
mkdir -p tests/testdata/{valid,invalid}

# 2. Makefile 작성
touch Makefile

# 3. 첫 Lexer 구현 시작
touch pkg/token/token.go
touch pkg/lexer/lexer.go
touch pkg/lexer/lexer_test.go

# 4. 테스트 실행
go test ./...
```

---

**생성일**: 2026-03-26
**상태**: 📍 포팅 준비 완료
**다음**: Lexer 구현 시작
