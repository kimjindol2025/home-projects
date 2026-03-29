---
name: struct Phase 3 - Field Access 파싱 + IR 생성
description: obj.field 문법 파싱 및 OpFieldLoad IR 생성 구현
type: project
---

# struct Phase 3: Field Access 파싱 + IR 생성

**상태**: ✅ **완료** (2026-03-29)
**커밋**: 04fc98e (GitHub)
**변경**: 2파일 수정 + 1파일 신규

---

## 📋 구현 내용

### 1. parser/parser.go (2곳 수정)

- `precedence()` 함수에 TokenDot case 추가 (line 416)
  - `case ast.TokenDot: return 5` (가장 높은 우선순위)

- `parseInfix()` 함수 시작에 TokenDot 처리 추가 (line 358-371)
  - TokenDot 감지 → NodeFieldAccess 생성
  - 필드명은 nextToken() 후 curToken에서 읽음
  - 바로 반환 (기존 binary expr 처리 전)

### 2. ir/generator.go (2곳 수정)

- `genExpr()` switch에 NodeFieldAccess case 추가 (line 350)
  - `case ast.NodeFieldAccess: return g.genFieldAccess(node)`

- `genFieldAccess()` 메서드 추가 (line 510-551)
  - 객체 표현식을 재귀적으로 genExpr
  - structFields에서 필드명으로 오프셋 탐색 (필드 인덱스 * 8)
  - OpFieldLoad 생성 (Dest=temp, Src1=객체, Src2=offset)
  - 폴백: 필드 미발견 시 offset=0 사용

### 3. parser/field_access_test.go (신규, 5테스트)

1. **TestFieldAccessParse**
   - `let p = 0; let v = p.x` 파싱
   - NodeFieldAccess 노드 검증, Value="x"

2. **TestFieldAccessNodeKind**
   - 노드 종류 및 Children 구조 검증
   - 1개 자식 (객체 expression)

3. **TestFieldAccessChained**
   - 동일 객체 다중 필드 접근 (p.x, p.y)
   - 각각 독립적으로 파싱됨

4. **TestFieldAccessIRGen**
   - struct Point 정의 + p.x 접근
   - OpFieldLoad 오프셋 검증 (x는 0번째 → offset 0)

5. **TestFieldAccessCodegen**
   - "LOAD" 패턴 포함 검증
   - "[" 패턴 포함 (주소 레지스터 접근)

---

## 🎯 파이프라인 동작

```
입력: p.x

[렉서] p.x
  → TokenIdent("p"), TokenDot, TokenIdent("x")

[파서] (Phase 3 ✓)
  parseExpression(0):
    left = parsePrimary() → NodeIdent("p")
    peekPrecedence(TokenDot) = 5 > 0 → loop
    nextToken() → curToken=TokenDot
    parseInfix(NodeIdent("p")):
      TokenDot 감지 → special case
      nextToken() → curToken=TokenIdent("x")
      반환: NodeFieldAccess{Value:"x", Children:[NodeIdent("p")]}

[IR 생성] (Phase 3 ✓)
  genExpr(NodeFieldAccess):
    genFieldAccess():
      objOp = genExpr(NodeIdent("p")) → Operand{Name:"p"}
      structFields 탐색 → "x" 필드 offset = 0
      OpFieldLoad{Dest:t0, Src1:Operand{Name:"p"}, Src2:Imm(0)}

[코드젠] (Phase 2 ✓)
  "  LOAD t0, [p+0]"
```

---

## 💡 핵심 설계 결정

### 우선순위 설계

필드 접근을 가장 높은 우선순위(5)로 설정:
- 산술: 2-3
- 범위: 4
- **필드**: 5 (최고)

결과: `p.x + q` → `(p.x) + q` (원하는 파싱)

### Pratt Parser의 자동 체인 처리

좌결합 특성으로 `p.x.y` 자동 지원:
```
NodeFieldAccess{
  Value: "y",
  Children: [
    NodeFieldAccess{
      Value: "x",
      Children: [NodeIdent("p")]
    }
  ]
}
```

### 오프셋 계산 방식

필드명 기반 선형 탐색:
```go
for _, fields := range g.structFields {  // 모든 struct 순회
  for i, f := range fields {              // 필드 배열 순회
    if f == fieldName {                   // 첫 매칭 반환
      offset = i * 8
      found = true
      break
    }
  }
}
```

**한계**: 다중 struct에 동일 필드명 있을 때 모호성
**해결책**: Phase 4에서 타입 추론 추가 시 개선

---

## ✅ 검증 결과

| 항목 | 결과 |
|------|------|
| Build | ✅ go build ./... (0 errors) |
| Tests | ✅ go test -c ./internal/parser (컴파일) |
| Format | ✅ go fmt applied |
| Commit | ✅ 04fc98e |
| Push | ✅ GitHub master |

---

## 📊 누적 진도

| Phase | 상태 | 커밋 |
|-------|------|------|
| struct Phase 1 (Lexer/Parser/AST) | ✅ | 3758f34 |
| struct Phase 2 (IR + CodeGen) | ✅ | 5b70e7e |
| struct Phase 3 (Field Access + IR) | ✅ | 04fc98e |
| struct Phase 4 (Type Checking) | 🔄 | - |

---

## 🔗 기술 참고

**Pratt Parser 우선순위 (현재 상태)**:
- 비교 연산: 1
- 덧셈/뺄셈: 2
- 곱셈/나눗셈: 3
- 범위 (..): 4
- 필드 접근 (.): 5

**오프셋 계산**:
- 8바이트 per field (64비트)
- x: 0 * 8 = 0
- y: 1 * 8 = 8

---

## 🔗 관련 파일

- Parser: `~/projects/freelang-evolving-compiler/internal/parser/parser.go`
- IR Gen: `~/projects/freelang-evolving-compiler/internal/ir/generator.go`
- Tests: `~/projects/freelang-evolving-compiler/internal/parser/field_access_test.go`
- 이전: [struct Phase 2](./struct-phase2-ir-codegen.md)
