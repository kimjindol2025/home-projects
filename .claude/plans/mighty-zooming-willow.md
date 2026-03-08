# FreeLang 셀프호스팅 기술 문제 해결 계획

## Context

셀프호스팅 검증 결과 두 가지 핵심 기술 문제가 발견됨:
1. **클로저 런타임 미동작** — opcode 7개가 no-op/전파 오류
2. **Result<T,E> 타입 없음** — 메인 파이프라인에 미구현

두 문제 모두 컴파일러를 FreeLang으로 재작성할 때 필수 조건.
script-runner 파이프라인에 Result/Option 완전 구현체가 참조 코드로 존재.

---

## Phase 1: 클로저 런타임 구현

### 문제 요약 (7개 버그)

| # | 위치 | 문제 |
|---|------|------|
| B1 | `src/analyzer/type-checker.ts:913` | capturedVars 수집 후 AST에 쓰지 않음 |
| B2 | `src/vm.ts:1220-1227` | LAMBDA_NEW/CAPTURE/SET_BODY → pc++ (no-op) |
| B3 | `src/vm/instruction-dispatcher.ts:571-581` | Lambda 핸들러 3개 전부 no-op |
| B4 | `src/vm.ts:1376` | callClosure()가 `{name:string}[]` 기대, AST는 `string[]` |
| B5 | `src/codegen/ir-generator.ts:1545` | LAMBDA_SET_BODY 후 클로저 객체 스택 push 없음 |
| B6 | `src/parser/ast.ts:122` | body 타입 `Expression`만, `BlockStatement` 불가 |
| B7 | `src/vm.ts:52` | `currentScope?: LocalScope` 선언만, 사용 안 함 |

### 수정 순서

**Step 1: B1 수정 — type-checker.ts capturedVars 전파**

파일: `src/analyzer/type-checker.ts`

```typescript
// Line 913 근처: validateLambda() 내부
const capturedVars = this.collectClosureVariables(lambda.body, context, paramNames);
lambda.capturedVars = capturedVars; // ← 이 한 줄 추가
return { compatible: true, ..., capturedVars };
```

**Step 2: B4 수정 — callClosure() 타입 통일**

파일: `src/vm.ts`

```typescript
// callClosure() 내부: string[]로 통일
if (closure.capturedVars && Array.isArray(closure.capturedVars)) {
  for (const varName of closure.capturedVars) {
    // 기존: capturedVar.name (object) → 수정: varName (string)
    const name = typeof varName === 'string' ? varName : (varName as any).name;
    if (savedVars.has(name)) {
      this.vars.set(name, savedVars.get(name));
    }
  }
}
```

**Step 3: B2 수정 — LAMBDA opcode VM 실행**

파일: `src/vm.ts` (Line 1220-1227)

```typescript
case Op.LAMBDA_NEW:
  this.currentLambda = { type: 'lambda', capturedVars: [], params: [], body: null };
  this.pc++;
  break;

case Op.LAMBDA_CAPTURE: {
  const varName = inst.arg as string;
  if (this.currentLambda && this.vars.has(varName)) {
    this.currentLambda.capturedVars.push(varName);
    this.currentLambda[varName] = this.vars.get(varName); // 값도 스냅샷
  }
  this.pc++;
  break;
}

case Op.LAMBDA_SET_BODY:
  if (this.currentLambda) {
    this.currentLambda.params = inst.params || [];
    this.currentLambda.sub = inst.sub || [];
    this.stack.push(this.currentLambda); // ← B5도 해결: 스택에 push
    this.currentLambda = null;
  }
  this.pc++;
  break;
```

**Step 4: B6 수정 — LambdaExpression body 타입**

파일: `src/parser/ast.ts`

```typescript
export interface LambdaExpression {
  body: Expression | BlockStatement; // ← BlockStatement 추가
  capturedVars?: string[];
  // ...
}
```

**Step 5: 클로저 테스트 작성**

파일: `self-hosting/test_closure.fl` (새 파일)

```freelang
// 기본 클로저
let x = 10;
let add = fn(n) -> x + n;
print(add(5)); // 기대: 15

// 고차 함수 + 클로저
fn makeCounter() {
  let count = 0;
  return fn() -> {
    count = count + 1;
    count
  };
}
let c = makeCounter();
print(c()); // 1
print(c()); // 2
```

---

## Phase 2: Result<T,E> 타입 구현

### 전략: script-runner 참조 구현 → 메인 파이프라인 이식

script-runner에 완전 구현체 존재 → 그대로 메인 파이프라인에 추가.

**Step 1: Op 코드 추가**

파일: `src/types.ts`

```typescript
// Op enum에 추가 (0xB0 영역)
WRAP_OK   = 0xB0,
WRAP_ERR  = 0xB1,
WRAP_SOME = 0xB2,
WRAP_NONE = 0xB3,
IS_OK     = 0xB4,
IS_ERR    = 0xB5,
IS_SOME   = 0xB6,
IS_NONE   = 0xB7,
UNWRAP    = 0xB8,
UNWRAP_OR = 0xB9,
```

**Step 2: AST Pattern 타입 확장**

파일: `src/parser/ast.ts`

```typescript
// Pattern 유니온에 추가
export interface OkPattern  { type: 'ok_pattern';   inner: Pattern; }
export interface ErrPattern { type: 'err_pattern';  inner: Pattern; }
export interface SomePattern{ type: 'some_pattern'; inner: Pattern; }
export interface NonePattern{ type: 'none_pattern'; }

export type Pattern = LiteralPattern | VariablePattern | WildcardPattern
  | StructPattern | ArrayPattern
  | OkPattern | ErrPattern | SomePattern | NonePattern; // ← 추가
```

**Step 3: IR 생성기 — MatchExpression + Result 패턴**

파일: `src/codegen/ir-generator.ts`

```typescript
// traverse() switch에 추가
case 'MatchExpression':
case 'match':
  this.generateMatchIR(node, out);
  break;

// 새 메서드 추가
private generateMatchIR(node: any, out: Inst[]): void {
  // 1. scrutinee 평가
  this.traverse(node.scrutinee, out);
  // 2. 각 arm: IS_OK/IS_ERR/IS_SOME/IS_NONE 분기
  for (const arm of node.arms) {
    this.generatePatternTest(arm.pattern, out);
    out.push({ op: Op.JUMP_IF_FALSE, arg: `arm_end_${i}` });
    this.traverse(arm.body, out);
    out.push({ op: Op.JUMP, arg: `match_end` });
    out.push({ op: Op.LABEL, arg: `arm_end_${i}` });
  }
  out.push({ op: Op.LABEL, arg: 'match_end' });
}
```

**Step 4: VM — Result opcode 실행**

파일: `src/vm.ts`

```typescript
case Op.WRAP_OK:  { const v = this.stack.pop(); this.stack.push({tag:'ok',  val:v}); this.pc++; break; }
case Op.WRAP_ERR: { const v = this.stack.pop(); this.stack.push({tag:'err', val:v}); this.pc++; break; }
case Op.WRAP_SOME:{ const v = this.stack.pop(); this.stack.push({tag:'some',val:v}); this.pc++; break; }
case Op.WRAP_NONE: { this.stack.push({tag:'none'}); this.pc++; break; }
case Op.IS_OK:   { const v = this.stack.pop(); this.stack.push(v?.tag === 'ok');   this.pc++; break; }
case Op.IS_ERR:  { const v = this.stack.pop(); this.stack.push(v?.tag === 'err');  this.pc++; break; }
case Op.IS_SOME: { const v = this.stack.pop(); this.stack.push(v?.tag === 'some'); this.pc++; break; }
case Op.IS_NONE: { const v = this.stack.pop(); this.stack.push(v?.tag === 'none'); this.pc++; break; }
case Op.UNWRAP:  {
  const v = this.stack.pop();
  if (v?.tag === 'ok' || v?.tag === 'some') { this.stack.push(v.val); }
  else throw new Error(`panic: unwrap on ${v?.tag}`);
  this.pc++; break;
}
```

**Step 5: 빌트인 함수 등록**

파일: `src/stdlib-builtins.ts`

```typescript
// Ok(), Err(), Some(), None 함수 등록
registry.register('Ok',   (v) => ({ tag: 'ok',   val: v }));
registry.register('Err',  (v) => ({ tag: 'err',  val: v }));
registry.register('Some', (v) => ({ tag: 'some', val: v }));
registry.register('None', ()  => ({ tag: 'none' }));
registry.register('isOk', (v) => v?.tag === 'ok');
registry.register('isErr',(v) => v?.tag === 'err');
```

**Step 6: Result 테스트**

파일: `self-hosting/test_result.fl` (새 파일)

```freelang
fn divide(a: int, b: int) -> Result<int, string> {
  if b == 0 { return Err("division by zero"); }
  return Ok(a / b);
}

let r = divide(10, 2);
match r {
  Ok(v) => print(v),      // 기대: 5
  Err(e) => print(e),
}

let r2 = divide(10, 0);
match r2 {
  Ok(v) => print(v),
  Err(e) => print(e),     // 기대: division by zero
}
```

---

## 수정 파일 목록

| 파일 | Phase | 수정 내용 |
|------|-------|---------|
| `src/analyzer/type-checker.ts` | 1 | capturedVars AST 전파 (1줄) |
| `src/vm.ts` | 1, 2 | LAMBDA opcode 구현, callClosure() 수정, Result opcode |
| `src/vm/instruction-dispatcher.ts` | 1 | Lambda 핸들러 구현 |
| `src/codegen/ir-generator.ts` | 1, 2 | LAMBDA_SET_BODY push, MatchExpression 구현 |
| `src/parser/ast.ts` | 1, 2 | body 타입 수정, Pattern 타입 추가 |
| `src/types.ts` | 2 | Result opcode 추가 |
| `src/stdlib-builtins.ts` | 2 | Ok/Err/Some/None 함수 등록 |
| `self-hosting/test_closure.fl` | 1 | 클로저 테스트 |
| `self-hosting/test_result.fl` | 2 | Result 테스트 |

---

## 검증 방법

```bash
# Phase 1 검증
npm run build
node dist/cli/index.js self-hosting/test_closure.fl
# 기대: 15, 1, 2 출력

# Phase 2 검증
node dist/cli/index.js self-hosting/test_result.fl
# 기대: 5, division by zero 출력

# 전체 빌드 검증
npm run build && npm test
```

---

## 완료 기준 (거짓보고 금지)

- [ ] `test_closure.fl` 실행 로그: 15, 1, 2 출력 확인
- [ ] `test_result.fl` 실행 로그: 5, "division by zero" 출력 확인
- [ ] `npm run build` 에러 없음
- [ ] 기존 테스트 regression 없음 (기존 통과 테스트 유지)
