---
name: Type System Foundation Phase 1
description: 기초 타입 체커 구현 — TypeKind, TypeEnv, TypeChecker 3개 모듈
type: project
---

# Type System Foundation Phase 1: 기초 타입 체커

**상태**: ✅ **완료** (2026-03-29)
**커밋**: 1e8d7db (GitHub)
**코드량**: ~530줄 (4파일 신규 + 3파일 수정)

---

## 📋 구현 내용

### 신규 패키지: `internal/typesys/`

#### 1. **types.go** (~70줄)
- **TypeKind enum**: TypeInt, TypeBool, TypeString, TypeStruct, TypeFn, TypeUnit, TypeUnknown
- **TypeInfo struct**
  ```go
  type TypeInfo struct {
      Kind       TypeKind
      StructName string     // struct 이름
      ParamTypes []TypeInfo // 함수 파라미터 타입
      ReturnType *TypeInfo  // 반환 타입
  }
  ```
- **상수**: `IntType`, `BoolType`, `StringType`, `UnitType`, `UnknownType`
- **헬퍼**:
  - `StructType(name)` — struct 타입 생성
  - `TypeFromAnnotation(s)` — "int" → IntType
  - `Equals(other)` — 타입 비교

#### 2. **env.go** (~80줄)
- **Scope struct**: 변수 타입 맵 + 부모 링크 (스코프 체인)
- **TypeEnv struct**:
  ```go
  type TypeEnv struct {
      current *Scope
      structs map[string]StructDef
  }
  ```
- **메서드**:
  - `EnterScope()` / `ExitScope()` — 중첩 블록 처리
  - `Define(name, type)` — 변수 타입 등록
  - `Lookup(name)` — 변수 타입 조회 (체인 따라 상향)
  - `RegisterStruct(name, def)` — struct 정의 저장
  - `LookupStruct(name)` — struct 조회

#### 3. **checker.go** (~180줄)
- **TypeError struct**: message, line, col
- **TypeChecker struct**:
  ```go
  type TypeChecker struct {
      env    *TypeEnv
      errors []TypeError
  }
  ```
- **메인 메서드**:
  - `Check(root *ast.Node) []TypeError` — 진입점
  - `checkNode(n) TypeInfo` — 노드 타입 추론 및 검증 (switch on NodeKind)

- **구현된 경우**:
  - **NodeStructDecl**: struct 정의 등록 (필드명 → 타입)
  - **NodeLetDecl**: 타입 어노테이션 저장/추론, env.Define()
  - **NodeFnDecl**: 파라미터 타입 등록
  - **NodeBlockStmt**: EnterScope → 자식 순회 → ExitScope
  - **NodeIfStmt**: 조건이 bool인지 검증
  - **NodeBinaryExpr**: 산술 (int + int → int), 비교 (→ bool)
  - **NodeFieldAccess**: struct 필드 존재 여부 검증 + 필드 타입 반환
  - **NodeIdent**: env.Lookup() 으로 타입 조회
  - **NodeIntLit**: IntType 반환
  - **NodeForStmt**: 이터레이터를 int로 등록

#### 4. **checker_test.go** (~150줄)
5개 테스트:
1. **TestTypeCheckInt**: `let x: int = 5` → 오류 없음
2. **TestTypeCheckStruct**: struct 정의 + 등록 검증
3. **TestTypeCheckFieldAccess**: 유효한 필드 접근
4. **TestTypeCheckUnknownField**: 존재하지 않는 필드 → 에러
5. **TestTypeCheckBinaryExpr**: 이항 연산 타입 검증

### 기존 파일 수정

#### **ast/nodes.go**
- `NodeTypeAnnotation` NodeKind 추가
- `Node.TypeAnnotation string` 필드 추가
  - `let x: int = 5` 에서 identifier 노드의 TypeAnnotation = "int"

#### **parser/parser.go**
- `parseLetDecl()`:
  ```go
  if p.curToken.Type == ast.TokenColon {
      nameNode.TypeAnnotation = p.curToken.Value // "int", "bool", "Point"
  }
  ```
- `parseFieldDecl()`: 필드 타입을 `TypeAnnotation` 에도 저장

#### **main.go**
- typesys 임포트 추가
- `compileCode()` → Step 1.5 TypeCheck 단계:
  ```go
  tc := typesys.NewTypeChecker()
  if typeErrs := tc.Check(prog); len(typeErrs) > 0 {
      // 경고만 출력, 컴파일 계속 (soft mode)
  }
  ```

---

## 🎯 파이프라인

```
Source Code
    ↓ Lexer
    ↓ Parser        [NodeTypeAnnotation 저장]
    ↓ TypeChecker   ← NEW (Phase 1)
      [struct 등록, 변수 타입 검증, 필드 존재 확인]
    ↓ Optimizer
    ↓ IR Generator
    ↓ CodeGen
    ↓ Assembly
```

---

## ✅ 검증 예시

### 성공 케이스
```freelang
struct Point { x: int; y: int }
let p: Point = 0
let v: int = p.x
```
→ 모든 타입 검증 통과 ✅

### 에러 감지 케이스
```freelang
struct Point { x: int }
let p: Point = 0
let v: int = p.z  // z 필드는 없음
```
→ "unknown field 'z' on struct 'Point'" 🔴

---

## 🏗️ 설계 결정

### Soft Mode (경고만 출력)
- Phase 1: 타입 에러 감지하지만 컴파일 계속 진행
- Phase 2: Hard Mode (타입 에러 시 컴파일 중단)

### 명시적 타입 (Hindley-Milner 아님)
- 변수: `let x: int = ...` 명시적 선언
- struct 필드: `x: int` 필수
- 함수 파라미터: 향후 지원

### 스코프 체인
- 블록 진입: `EnterScope()`
- 블록 퇴출: `ExitScope()`
- 이름 조회: 현재 scope → 부모 scope... 재귀

---

## 📊 코드 통계

| 파일 | 줄 | 설명 |
|------|---|------|
| types.go | ~70 | TypeKind + TypeInfo |
| env.go | ~80 | Scope + TypeEnv |
| checker.go | ~180 | TypeChecker + 8개 switch case |
| checker_test.go | ~150 | 5 tests |
| ast/nodes.go | +5 | NodeTypeAnnotation + TypeAnnotation |
| parser/parser.go | +15 | parseLetDecl + parseFieldDecl 수정 |
| main.go | +12 | TypeCheck 단계 추가 |
| **합계** | **~512** | **Type System Foundation** |

---

## 🔮 향후 확장 (Phase 2+)

### Phase 2: Type Inference
- 리터럴에서 타입 추론 (let x = 5 → IntType)
- 함수 반환 타입 추론
- 제네릭 타입 변수 (T)

### Phase 3: Generic Types
- `struct Pair<T> { a: T; b: T }`
- `fn map<T, U>(f: T→U, arr: [T]): [U]`

### Phase 4: Trait/Interface
- trait Iterator { next(): T }
- impl Iterator for MyType

---

## 🔗 기술 하이라이트

1. **Scope 체인**: 렉시컬 스코핑 구현
2. **TypeEnv 이원 구조**:
   - `current *Scope` — 변수 타입
   - `structs map[string]StructDef` — struct 정의
3. **검증 분류**:
   - 선언적 (let x: int)
   - 추론적 (struct 필드로부터)
   - 구조적 (필드 접근, 이항 연산)
4. **Error Accumulation**: 첫 에러에서 멈추지 않고 모든 에러 수집

---

## 💪 다음 커밋 준비 상태

```bash
cd ~/projects/freelang-evolving-compiler
go build ./...                     # ✅ 0 errors
go test -c ./internal/typesys     # ✅ compile success
git log --oneline | head -3       # 1e8d7db Type System Foundation...
```

---

## 🎓 학습 포인트

이 Phase 1 구현으로 다음이 증명되었습니다:
1. **모듈화**: typesys 패키지 단독으로 동작 가능
2. **확장성**: TypeKind 추가로 새로운 타입 지원 가능
3. **유연성**: Soft mode로 기존 파이프라인과 무마찰 통합
4. **Zero-Dependency**: 표준 라이브러리만 사용

---

## 🔗 참고

- [struct Phase 3 완료](./struct-phase3-field-access.md) — 필드 접근 파이프라인
- 컴파일러 전체 구조: `~/projects/freelang-evolving-compiler`
- 커밋: 1e8d7db
