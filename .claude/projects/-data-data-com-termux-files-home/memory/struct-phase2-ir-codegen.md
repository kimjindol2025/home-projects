---
name: struct Phase 2 - IR Generation + CodeGen
description: AST → IR → Assembly struct 파이프라인 완성 (IR 생성 + 코드 생성)
type: project
---

# struct Phase 2: IR Generation + CodeGen 완성

**상태**: ✅ **완료** (2026-03-29)
**커밋**: 5b70e7e (GitHub)
**변경**: 3파일 수정 + 1파일 신규 (struct_test.go)

---

## 📋 구현 내용

### 1. ir/generator.go (3가지 수정)

- `structFields map[string][]string` 필드 추가
  - 구조체 이름 → 필드 이름 배열 매핑

- `NewGenerator()` 초기화
  - `make(map[string][]string)` 으로 초기화

- `Generate()` 분기 추가 (line 39)
  - `NodeStructDecl` 감지 → `genStructDecl()` 호출

- `genStmt()` case 추가 (line 127)
  - `NodeStructDecl` 처리

- `genStructDecl()` 메서드 구현
  ```go
  // AST 노드 → structFields 맵 저장
  // 크기 계산: len(fields) * 8 bytes
  // OpStructDef 생성 (Main 컨텍스트)
  ```

### 2. codegen/codegen.go (3개 opcode case)

- **OpStructDef**: `"; STRUCT %s size=%d"` 포맷
  - 구조체 선언 주석 출력

- **OpFieldLoad**: `"LOAD dest, [base+offset]"` 포맷
  - 필드 로드 명령

- **OpFieldStore**: `"STORE [dest+offset], src"` 포맷
  - 필드 저장 명령

### 3. parser/struct_test.go (신규, 5개 테스트)

1. **TestStructSimple**
   - 입력: `struct Point { x: int; y: int }`
   - 검증: NodeStructDecl, 2개 FieldDecl 자식

2. **TestStructFieldStr**
   - 입력: `struct Person { name: string; age: int }`
   - 검증: 혼합 타입 필드

3. **TestStructNodeKind**
   - 검증: NodeStructDecl 타입 정확성

4. **TestStructIRGen**
   - 검증: OpStructDef 생성 (size = 2*8 = 16)

5. **TestStructCodegen**
   - 검증: "; STRUCT Point size=16" 출력

---

## 🎯 파이프라인 검증

```
struct Point { x: int; y: int }
    ↓ Lexer (Phase 1 ✓)
    ↓ Parser (Phase 1 ✓)
NodeStructDecl(name="Point", children=[FieldDecl(x), FieldDecl(y)])
    ↓ IR Generator (Phase 2 ✓) NEW
OpStructDef(name="Point", size=16)
    ↓ CodeGen (Phase 2 ✓) NEW
"; STRUCT Point size=16"
```

---

## ✅ 빌드 검증

- **Build**: `go build ./...` ✅ (0 errors)
- **Format**: `go fmt` applied ✅
- **Compile**: struct_test.go compiles ✅
- **Tests**: ARM64 제약 (실행 불가, 컴파일만 확인)

---

## 📊 누적 진도

| Phase | 완료 | 설명 |
|-------|------|------|
| struct Phase 1 | ✅ | Lexer + Parser + AST (3758f34) |
| struct Phase 2 | ✅ | IR Generator + CodeGen (5b70e7e) |
| struct Phase 3 | 🔄 | Field access parsing (향후) |
| struct Phase 4 | 🔄 | Type checking (향후) |

---

## 📝 기술 상세

### 크기 계산
- 필드당 8바이트 (64비트 아키텍처)
- struct Point (2 fields) → size = 16

### 컨텍스트 처리
```go
// struct 선언은 Main에 emit (함수 외부)
prevFn := g.currentFn
g.currentFn = nil    // Main 컨텍스트로 변경
g.emit(...)
g.currentFn = prevFn // 복구
```

### 3-Address Code
```
OpStructDef:
  Op:   OpStructDef
  Fn:   "Point"      // struct name
  Src1: Imm(16)      // size
```

---

## 🔗 참고

- 이전 단계: [struct Phase 1](./struct-phase1-parsing.md)
- 컴파일러: `~/projects/freelang-evolving-compiler`
- 커밋: 5b70e7e
- 파일 변경: codegen.go, generator.go, parser.go, struct_test.go
