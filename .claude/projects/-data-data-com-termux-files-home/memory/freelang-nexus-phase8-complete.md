---
name: FreeLang Nexus Phase 8 - Stdlib 추가 확장 완료
description: Phase 8 완전 완료 - if/else 문장 + int_cast 함수 지원
type: project
---

## 상태
✅ **완전 완료** (2026-03-20)

## 목표
Phase 7에서 기본 4개 내장 함수 구현.
Phase 8: 조건문(if/else) + 실용적인 타입 변환(int_cast) 추가.

## 구현 내용

### 1. if/else 문장 지원

**파일 3개 수정**:

1. **ast.ts** (+15줄): `IfStatement` 노드 추가
   ```typescript
   export interface IfStatement extends ASTNode {
     type: 'IfStatement';
     condition: Expression;
     thenBranch: Statement[];
     elseBranch?: Statement[];
   }
   ```

2. **nexus-parser.ts** (+35줄): `parseIfStatement` 메서드
   ```typescript
   if (token.type === TokenType.IF) {
     return this.parseIfStatement();
   }
   ```
   - 조건 파싱
   - then 블록 파싱
   - else 블록 (선택사항) 파싱

3. **nexus-codegen.ts** (+20줄): `genVStatement`에서 IfStatement 처리
   ```typescript
   } else if (s.type === 'IfStatement') {
     const ifStmt = s as AST.IfStatement;
     const cond = this.genVExpression(ifStmt.condition);
     this.writeC(`${this.indent()}if (${cond}) {`);
     // then/else 블록 생성
   }
   ```

### 2. int_cast 내장 함수

**파일**: `nexus-codegen.ts`

```typescript
// V_BUILTINS에 추가
private readonly V_BUILTINS = new Set([
  'println', 'print', 'len', 'to_string', 'int_cast'
]);

// genVBuiltinCall에서 처리
if (name === 'int_cast') {
  const firstArg = args[0] as any;
  if (firstArg.type === 'String') {
    const str = firstArg.value.slice(1, -1);
    return `atoi("${str}")`;
  }
  const arg = this.genVExpression(firstArg);
  return `atoi(${arg})`;
}
```

**매핑**:
- `int_cast("42")` → `atoi("42")`
- `int_cast(s)` → `atoi(s)`

**헤더**: `#include <stdlib.h>` 추가

### 3. 예제 파일 (`examples/if_demo.fl`)

```fl
@mode(v)
fn check(x: i64) -> i64 {
  if x > 0 {
    println("positive")
    return 1
  } else {
    println("non-positive")
    return 0
  }
}
```

## 생성 결과

**Input**:
```fl
fn check(x: i64) -> i64 {
  if x > 0 {
    println("positive")
    return 1
  } else {
    println("non-positive")
    return 0
  }
}
```

**Output C 코드**:
```c
long long check(long long x) {
    if (x > 0) {
        printf("positive\n");
        return 1;
    }
    else {
        printf("non-positive\n");
        return 0;
    }
}
```

## 메트릭

| 항목 | 값 |
|------|-----|
| AST 노드 추가 | +15줄 |
| Parser 메서드 추가 | +35줄 |
| Codegen if 처리 | +20줄 |
| Codegen int_cast | +10줄 |
| 테스트 | 7개 |
| 예제 파일 | 1개 |
| 전체 테스트 | 59/59 통과 (100%) |
| 커밋 | e6abaad |

## 완성도

- **Parser**: return, assign, expr, **if/else** (4개 문장 타입)
- **Codegen**: println, print, len, to_string, **int_cast** (5개 내장 함수)
- **조건식**: 비교연산자(>, <, >=, <=), 논리연산자(&&, ||) 모두 지원

---

**완성도**: 85-90% → 90-95%
**다음 옵션**: while 루프 / for 루프 / array 기본 지원

