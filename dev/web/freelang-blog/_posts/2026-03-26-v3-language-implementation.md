---
title: "V3 언어는 어떻게 만들어질까? - Lexer에서 Runtime까지"
date: 2026-03-26
author: Content Writer
category: Technical
tags:
  - FreeLang
  - V3
  - Compiler
  - Programming Language
  - Lexer
  - Parser
  - Code Generation
---

# V3 언어는 어떻게 만들어질까?
## Lexer → Parser → Compiler → Runtime 파이프라인 분석

**글을 읽으면 얻을 수 있는 것:**
1. 프로그래밍 언어가 어떻게 작동하는지 이해하기
2. 컴파일러 4단계 파이프라인의 역할 파악하기
3. 실제 V3 코드 예제로 따라하기

---

## 배경: 왜 언어 구현을 배워야 할까?

프로그래밍을 배우다 보면 자연스러운 질문이 생깁니다.

> "Python이나 JavaScript는 어떻게 코드를 이해할까?"
> "내 코드가 어떻게 컴퓨터가 이해할 수 있는 형태가 될까?"

이 질문에 답하려면 **컴파일러(Compiler)**를 이해해야 합니다. FreeLang의 V3는 간단하면서도 완전한 프로그래밍 언어를 구현하는 방법을 보여줍니다.

---

## 문제: 텍스트는 어떻게 코드가 될까?

우리가 작성한 텍스트 코드는 그저 문자열입니다:

```
INT x = 10
```

컴퓨터는 이 문자열을 이해하지 못합니다. 컴퓨터가 이해하려면:

1. **문자 → 토큰** (이게 뭐지?)
2. **토큰 → 문법** (어떻게 배열되어 있지?)
3. **문법 → 명령어** (뭘 하라는 건지?)
4. **명령어 → 실행** (실제로 동작!)

이 과정이 바로 **컴파일러의 일**입니다.

---

## 해결책: V3 컴파일러 파이프라인

V3는 다음 4단계를 거쳐 코드를 실행합니다:

```
┌─────────┐     ┌────────┐     ┌──────────┐     ┌──────────┐
│ 소스코드 │ --> │ Lexer  │ --> │  Parser  │ --> │ Compiler │ --> 실행
└─────────┘     └────────┘     └──────────┘     └──────────┘
   "INT x=10"    [INT, x, =, 10]  AST 트리      바이트코드
```

각 단계를 자세히 살펴봅시다.

---

## 1단계: Lexer (어휘 분석)

### 역할: 문자 → 토큰

Lexer는 코드를 **토큰(Token)**이라는 작은 단위로 쪼갭니다.

```
입력: "INT x = 10"
출력: [
  { type: 'INT', value: 'INT' },
  { type: 'IDENTIFIER', value: 'x' },
  { type: 'ASSIGN', value: '=' },
  { type: 'NUMBER', value: '10' }
]
```

### 코드 예제

```typescript
// Lexer의 핵심 로직
interface Token {
  type: string;
  value: string;
  line: number;
}

class Lexer {
  private tokens: Token[] = [];
  private current = 0;

  tokenize(code: string): Token[] {
    const lines = code.split('\n');
    let lineNum = 1;

    for (const line of lines) {
      this.tokenizeLine(line, lineNum);
      lineNum++;
    }

    return this.tokens;
  }

  private tokenizeLine(line: string, lineNum: number) {
    let i = 0;

    while (i < line.length) {
      const char = line[i];

      // 숫자
      if (/[0-9]/.test(char)) {
        let num = '';
        while (i < line.length && /[0-9]/.test(line[i])) {
          num += line[i];
          i++;
        }
        this.tokens.push({ type: 'NUMBER', value: num, line: lineNum });
        continue;
      }

      // 식별자 또는 키워드
      if (/[a-zA-Z_]/.test(char)) {
        let word = '';
        while (i < line.length && /[a-zA-Z0-9_]/.test(line[i])) {
          word += line[i];
          i++;
        }

        const type = ['INT', 'ARR', 'IF', 'FOR', 'FUNC'].includes(word)
          ? 'KEYWORD'
          : 'IDENTIFIER';

        this.tokens.push({ type, value: word, line: lineNum });
        continue;
      }

      // 연산자
      if (['=', '>', '<', '+', '-'].includes(char)) {
        this.tokens.push({ type: 'OPERATOR', value: char, line: lineNum });
        i++;
        continue;
      }

      i++;
    }
  }
}
```

---

## 2단계: Parser (구문 분석)

### 역할: 토큰 → 추상 문법 트리(AST)

Parser는 토큰들을 **규칙**에 따라 조직합니다. 이를 **추상 문법 트리(Abstract Syntax Tree, AST)**라고 합니다.

```
토큰: [INT, x, =, 10]
↓
AST:
├─ Declaration
   ├─ type: 'INT'
   ├─ name: 'x'
   └─ value: 10 (숫자)
```

### 코드 예제

```typescript
// AST 노드 타입 정의
interface ASTNode {
  type: string;
}

interface Declaration extends ASTNode {
  type: 'Declaration';
  dataType: string;  // INT, ARR 등
  name: string;
  value: ASTNode;
}

interface IfStatement extends ASTNode {
  type: 'IfStatement';
  condition: ASTNode;
  body: ASTNode[];
  elseBody?: ASTNode[];
}

class Parser {
  private tokens: Token[] = [];
  private current = 0;

  parse(tokens: Token[]): ASTNode[] {
    this.tokens = tokens;
    this.current = 0;

    const ast: ASTNode[] = [];
    while (!this.isAtEnd()) {
      const stmt = this.parseStatement();
      if (stmt) ast.push(stmt);
    }
    return ast;
  }

  private parseStatement(): ASTNode | null {
    const token = this.peek();

    if (token.value === 'INT') {
      return this.parseDeclaration();
    }

    if (token.value === 'IF') {
      return this.parseIfStatement();
    }

    if (token.value === 'FUNC') {
      return this.parseFunctionDef();
    }

    this.advance();
    return null;
  }

  private parseDeclaration(): Declaration {
    const dataType = this.advance().value;  // INT
    const name = this.advance().value;       // x
    this.advance();                          // =
    const value = this.parseExpression();    // 10

    return {
      type: 'Declaration',
      dataType,
      name,
      value
    };
  }

  private parseExpression(): ASTNode {
    const token = this.peek();
    this.advance();

    if (token.type === 'NUMBER') {
      return {
        type: 'Literal',
        value: parseInt(token.value)
      };
    }

    return { type: 'Unknown' };
  }

  private parseIfStatement(): IfStatement {
    this.advance();  // IF
    const condition = this.parseExpression();

    const body: ASTNode[] = [];
    while (this.peek().value !== 'ELSE' && !this.isAtEnd()) {
      const stmt = this.parseStatement();
      if (stmt) body.push(stmt);
    }

    return {
      type: 'IfStatement',
      condition,
      body
    };
  }

  private peek(): Token {
    return this.tokens[this.current];
  }

  private advance(): Token {
    return this.tokens[this.current++];
  }

  private isAtEnd(): boolean {
    return this.current >= this.tokens.length;
  }
}
```

---

## 3단계: Compiler (코드 생성)

### 역할: AST → 바이트코드

Compiler는 AST를 **바이트코드**로 변환합니다. 바이트코드는 런타임이 이해할 수 있는 간단한 명령어 집합입니다.

```
AST:
├─ Declaration (INT x = 10)

바이트코드:
PUSH 10      // 스택에 10 넣기
STORE x      // x에 저장하기
```

### 코드 예제

```typescript
interface Instruction {
  op: string;  // PUSH, STORE, LOAD 등
  arg?: any;
}

class Compiler {
  private instructions: Instruction[] = [];
  private variables: Map<string, number> = new Map();

  compile(ast: ASTNode[]): Instruction[] {
    this.instructions = [];
    this.variables.clear();

    for (const node of ast) {
      this.compileNode(node);
    }

    return this.instructions;
  }

  private compileNode(node: ASTNode) {
    if (node.type === 'Declaration') {
      const decl = node as Declaration;

      // 값 먼저 컴파일
      this.compileExpression(decl.value);

      // STORE 명령어 추가
      const varIndex = this.variables.size;
      this.variables.set(decl.name, varIndex);
      this.instructions.push({
        op: 'STORE',
        arg: varIndex
      });
    }

    if (node.type === 'IfStatement') {
      const ifStmt = node as IfStatement;

      // 조건 계산
      this.compileExpression(ifStmt.condition);

      // 조건부 점프
      const jumpIdx = this.instructions.length;
      this.instructions.push({
        op: 'JUMP_IF_FALSE',
        arg: -1  // 나중에 수정
      });

      // Body 컴파일
      for (const stmt of ifStmt.body) {
        this.compileNode(stmt);
      }

      // 점프 주소 수정
      this.instructions[jumpIdx].arg = this.instructions.length;
    }
  }

  private compileExpression(node: ASTNode) {
    if (node.type === 'Literal') {
      const lit = node as any;
      this.instructions.push({
        op: 'PUSH',
        arg: lit.value
      });
    }
  }
}
```

---

## 4단계: Runtime (실행)

### 역할: 바이트코드 실행

Runtime은 바이트코드를 하나씩 읽어 실행합니다. **스택(Stack)**을 사용해 값들을 관리합니다.

```
바이트코드:
[PUSH 10, STORE x]

실행:
1. PUSH 10: 스택 = [10]
2. STORE x: 변수 x = 10, 스택 = []
```

### 코드 예제

```typescript
class Runtime {
  private stack: number[] = [];
  private variables: Map<string, number> = new Map();
  private instructions: Instruction[];

  execute(instructions: Instruction[]): void {
    this.instructions = instructions;
    this.stack = [];
    this.variables.clear();

    for (const instr of instructions) {
      this.executeInstruction(instr);
    }
  }

  private executeInstruction(instr: Instruction) {
    switch (instr.op) {
      case 'PUSH':
        this.stack.push(instr.arg);
        break;

      case 'STORE':
        const value = this.stack.pop();
        this.variables.set(`var_${instr.arg}`, value!);
        break;

      case 'LOAD':
        const val = this.variables.get(`var_${instr.arg}`);
        if (val !== undefined) {
          this.stack.push(val);
        }
        break;

      case 'JUMP_IF_FALSE':
        const condition = this.stack.pop();
        if (!condition) {
          // 점프 처리 (실제 구현에선 IP 수정)
        }
        break;
    }
  }

  getVariables(): Map<string, number> {
    return this.variables;
  }
}
```

---

## 실제 사용: 전체 파이프라인

### 초보자라면?

간단한 코드를 어떻게 실행되는지 따라가봅시다.

```v3
INT x = 10
println(x)
```

**단계별 실행:**

```
1️⃣ Lexer
입력: "INT x = 10"
출력: [INT, x, =, 10]

2️⃣ Parser
입력: [INT, x, =, 10]
출력: AST { Declaration, type: 'INT', name: 'x', value: 10 }

3️⃣ Compiler
입력: AST
출력: [PUSH 10, STORE 0]

4️⃣ Runtime
PUSH 10  → 스택: [10]
STORE 0  → x = 10, 스택: []
```

### 전문가라면?

복잡한 제어 흐름도 같은 파이프라인으로 처리됩니다.

```v3
INT x = 5
IF x > 3:
  println("big")
ELSE:
  println("small")
```

이 코드도 정확히 같은 4단계를 거칩니다. Parser가 IF 구조를 인식하고, Compiler가 조건부 점프 명령어를 생성하며, Runtime이 실행합니다.

---

## 균형잡힌 관점: 언어 구현의 장점과 한계

### 장점 ✅

| 측면 | 설명 |
|------|------|
| **이해도** | 프로그래밍의 본질을 깊이 있게 이해 가능 |
| **문제 해결** | 언어 기능을 직접 디자인할 수 있음 |
| **디버깅** | 예상 밖의 동작을 쉽게 추적 가능 |

### 한계 ⚠️

| 측면 | 설명 |
|------|------|
| **복잡성** | 완전한 언어는 매우 복잡함 |
| **성능** | 직접 구현한 런타임은 최적화 부족 |
| **호환성** | 기존 생태계와 연동 어려움 |

**개선 방안:**
- V3는 최소한의 기능으로 시작 (학습 목표)
- 단계적으로 기능 추가 (Array, Function 등)
- WASM으로 컴파일해 성능 향상 계획 중

---

## 다음 단계

이제 컴파일러 파이프라인을 이해했으니, 더 깊이 있는 주제들을 탐색할 수 있습니다.

**다음에 읽을 거리:**
1. [FreeLang 모듈 시스템 완벽 가이드](#) - import/export 메커니즘
2. [비동기 프로그래밍: async/await 구현](#) - 동시성 처리
3. [성능 최적화: AST vs 바이트코드 vs JIT](#) - 실행 속도 비교

**V3 직접 해보고 싶으신가요?**

```bash
# FreeLang 저장소에서 V3 실행
git clone https://github.com/freelang/freelang
cd freelang
npm install
npm run v3-demo
```

---

**참고 자료:**
- [Crafting Interpreters - Robert Nystrom](https://craftinginterpreters.com/)
- [Engineering a Compiler - Keith Cooper & Linda Torczon](https://www.elsevier.com/books/engineering-a-compiler/cooper/978-0-12-815412-0)
- FreeLang V3 공식 구현: https://github.com/freelang/v3-compiler

---

이 글이 도움이 되었다면? 👍

- 다른 개발자와 공유해주세요
- "다음엔 뭘 구현해볼까?" 댓글로 주제 제안해주세요
- V3에 시도해본 흥미로운 기능이 있으면 공유해주세요!

**Happy Coding! 🚀**
