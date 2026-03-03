# FreeLang → Z-Lang 트랜스파일러 구현 계획

## Context

FreeLang v4 소스 코드(.fl)를 Z-Lang 소스 코드(.z)로 변환하는 트랜스파일러를 구현한다.
변환된 .z 파일은 Z-Lang 컴파일러(C++/LLVM)로 컴파일되어 네이티브 실행 파일이 된다.

두 언어의 핵심 문법이 90% 동일하므로 트랜스파일러는 비교적 간단하게 구현 가능하다.

---

## 구현 방식

**TypeScript + FreeLang 파서 재사용** 방식.
FreeLang v4의 기존 파서(`freelang-v4/src/parser.ts`)를 그대로 재사용해서 AST를 생성한 후,
AST를 Z-Lang 소스 코드 문자열로 변환(code generation)한다.

```
fl 소스
  ↓  (기존 FreeLang Lexer + Parser 재사용)
FreeLang AST
  ↓  (새로 구현: ZLangCodeGen)
Z-Lang 소스 (.z)
  ↓  (기존 Z-Lang 컴파일러)
실행 파일
```

---

## 문법 변환 매핑표

| FreeLang | Z-Lang | 비고 |
|----------|--------|------|
| `var x: i32 = 10` | `let x: i32 = 10;` | var→let, 세미콜론 추가 |
| `let x: i32 = 10` | `let x: i32 = 10;` | 세미콜론만 추가 |
| `const x: i32 = 10` | `let x: i32 = 10;` | const→let, 세미콜론 추가 |
| `fn f(): i32 {}` | `fn f() -> i32 {}` | `:` → `->` |
| `fn f() -> i32 {}` | `fn f() -> i32 {}` | 동일 |
| `for i in range(1,5){}` | `let i:i64=1;\nwhile i<=5 { ... i=i+1; }` | for-in → while 변환 |
| `for x in arr {}` | (제한적 지원) | 배열 루프는 Phase 2 |
| `match expr { 1 => ... }` | (미지원) | Phase 2 |
| `spawn { ... }` | (미지원) | Phase 2 (무시 처리) |
| `x * y` (암묵적 반환) | `return x * y;` | 마지막 식 → return |
| `println(x)` | `print(x); print("\n");` | stdlib 매핑 |
| `str(x)` | `toString(x)` | 함수명 매핑 |

---

## 구현 파일 구조

```
freelang-to-zlang/           (새 프로젝트 디렉토리)
├── src/
│   ├── transpiler.ts        (핵심: AST → Z-Lang 변환)
│   ├── stdlib_map.ts        (FreeLang stdlib → Z-Lang stdlib 매핑)
│   └── index.ts             (CLI 엔트리)
├── tests/
│   ├── test_basic.ts        (기본 변환 테스트)
│   └── test_e2e.ts          (E2E: .fl → .z → 실행 테스트)
├── examples/
│   ├── factorial.fl         (테스트 입력)
│   ├── fizzbuzz.fl          (테스트 입력)
│   └── hello.fl             (테스트 입력)
└── package.json
```

---

## 핵심 구현: transpiler.ts

FreeLang AST를 순회하면서 Z-Lang 코드 문자열 생성.

```typescript
class ZLangCodeGen {
  generate(ast: Program): string;

  // 선언문
  genFunction(node: FnDecl): string;
  genVarDecl(node: VarDecl): string;       // var/let/const → let + ;

  // 제어흐름
  genIf(node: IfStmt): string;
  genWhile(node: WhileStmt): string;
  genForIn(node: ForInStmt): string;       // for-in → while 변환
  genReturn(node: ReturnStmt): string;     // return + ;

  // 표현식
  genBinaryOp(node: BinaryOp): string;
  genCall(node: CallExpr): string;         // println → print + \n
  genLiteral(node: Literal): string;
  genIdentifier(node: Identifier): string;
  genAssignment(node: Assignment): string;

  // 타입 변환
  convertType(freeLangType: string): string;   // i32→i32, f64→f64 (대부분 동일)
  convertReturnType(node: FnDecl): string;     // `: type` → `-> type`
}
```

---

## 핵심 변환 로직

### 1. var/const → let (세미콜론 추가)
```typescript
genVarDecl(node: VarDecl): string {
  return `let ${node.name}: ${this.convertType(node.type)} = ${this.genExpr(node.init)};\n`;
}
```

### 2. 함수 반환타입 변환
```typescript
genFunction(node: FnDecl): string {
  const params = node.params.map(p => `${p.name}: ${this.convertType(p.type)}`).join(", ");
  const retType = node.returnType ? ` -> ${this.convertType(node.returnType)}` : "";
  const body = this.genBlock(node.body, node.returnType);
  return `fn ${node.name}(${params})${retType} {\n${body}}\n`;
}
```

### 3. for-in range → while 변환 (핵심)
```typescript
genForIn(node: ForInStmt): string {
  if (isRangeCall(node.iter)) {
    const [start, end] = extractRangeArgs(node.iter);
    return `
let ${node.var}: i64 = ${start};
while ${node.var} < ${end} {
${this.genBlock(node.body)}
  ${node.var} = ${node.var} + 1;
}`;
  }
  // 배열 루프: 추후 지원
  return `// TODO: for-in array not yet supported\n`;
}
```

### 4. 암묵적 반환 처리
```typescript
genBlock(block: Block, returnType?: string): string {
  const stmts = block.stmts;
  if (returnType && returnType !== "void" && stmts.length > 0) {
    const last = stmts[stmts.length - 1];
    if (isExpression(last) && !isReturn(last)) {
      // 마지막 표현식을 return으로 변환
      return stmts.slice(0, -1).map(s => this.genStmt(s)).join("") +
             `return ${this.genExpr(last)};\n`;
    }
  }
  return stmts.map(s => this.genStmt(s)).join("");
}
```

### 5. println → print + \n 변환
```typescript
const STDLIB_MAP: Record<string, string> = {
  "println": "__fl_println",   // 특수 처리: print(x); print("\n");
  "str":     "toString",
  "print":   "print",
  "input":   "input",
};

genCall(node: CallExpr): string {
  if (node.name === "println") {
    return `print(${this.genExpr(node.args[0])});\nprint("\\n");\n`;
  }
  const mapped = STDLIB_MAP[node.name] ?? node.name;
  return `${mapped}(${node.args.map(a => this.genExpr(a)).join(", ")});\n`;
}
```

---

## 단계별 구현 순서

### Phase 1: 기본 트랜스파일러 (Day 1)
- [ ] 프로젝트 셋업 (package.json, tsconfig.json)
- [ ] FreeLang 파서 의존성 연결
- [ ] transpiler.ts 기본 구조
- [ ] VarDecl, Function, Return, BinaryOp, Literal 변환
- [ ] 3개 예제 파일 변환 테스트 (hello.fl, factorial.fl, fizzbuzz.fl)

### Phase 2: for-in + 고급 기능 (Day 2)
- [ ] for-in range → while 변환
- [ ] println/str stdlib 매핑
- [ ] if/else, while 변환
- [ ] 암묵적 반환 처리
- [ ] 10개 예제 변환 테스트

### Phase 3: 검증 + CLI (Day 3)
- [ ] E2E 테스트 (변환 → 컴파일 → 실행)
- [ ] CLI 도구 (fl2z 명령)
- [ ] 변환 보고서 생성
- [ ] GOGS 저장

---

## 주요 파일 경로

### 재사용 (기존)
- `/data/data/com.termux/files/home/freelang-v4/src/lexer.ts`
- `/data/data/com.termux/files/home/freelang-v4/src/parser.ts`
- `/data/data/com.termux/files/home/freelang-v4/src/ast.ts`

### 신규 생성
- `/data/data/com.termux/files/home/freelang-to-zlang/src/transpiler.ts`
- `/data/data/com.termux/files/home/freelang-to-zlang/src/stdlib_map.ts`
- `/data/data/com.termux/files/home/freelang-to-zlang/src/index.ts`
- `/data/data/com.termux/files/home/freelang-to-zlang/tests/test_basic.ts`
- `/data/data/com.termux/files/home/freelang-to-zlang/tests/test_e2e.ts`

---

## 검증 방법

### 기능 테스트
```bash
# hello.fl 변환 테스트
npx ts-node src/index.ts examples/hello.fl -o /tmp/hello.z
cat /tmp/hello.z    # Z-Lang 코드 확인

# 팩토리얼 E2E 테스트
npx ts-node src/index.ts examples/factorial.fl -o /tmp/factorial.z
./zlang_compiler /tmp/factorial.z -o /tmp/factorial
/tmp/factorial
```

### 예상 변환 결과
```
입력 (FreeLang):
fn factorial(n: i32): i32 {
  if n <= 1 { return 1 }
  return n * factorial(n + -1)
}

출력 (Z-Lang):
fn factorial(n: i32) -> i32 {
  if n <= 1 {
    return 1;
  }
  return n * factorial(n + -1);
}
```

---

## 예상 통계
- 총 코드: 600줄 (transpiler 400 + tests 200)
- 변환 커버리지: 85% (FreeLang 핵심 기능)
- E2E 테스트: 5개 예제 파일
- 소요시간: 3일
