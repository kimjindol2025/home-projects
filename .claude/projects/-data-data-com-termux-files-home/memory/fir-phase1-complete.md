---
name: fir Phase 1+2 Complete
description: fir Phase 1+2 완료 — IR 인터프리터 + x86-64 코드젠 + FL 프론트엔드 검증
type: project
---

# fir Phase 1 COMPLETE (2026-03-22)

**Why:** IR이 실행 의미를 정의하는 플랫폼임을 native 실행으로 증명.
**How to apply:** Phase 2 진입 시 이 기반 위에서 작업. IR 1.0 spec은 FROZEN.

## 검증 결과

| 항목 | 결과 |
|------|------|
| Interpreter | 17/17 PASS |
| Codegen (native) | 9/9 PASS |
| interp == native | ✅ 완전 일치 |

## 핵심 버그 2개 해결

### 1. RSP 정렬 (System V ABI)
- **증상**: 함수 호출 후 반환값 틀림
- **원인**: `sub rsp, N` 에서 N이 16배수 → CALL 전 RSP 미정렬
- **수정**: `(8 + N) % 16 == 0` → `N = ceil(size/16)*16 - 8` (8, 24, 40...)
- **위치**: `codegen.go:stackSize()`

### 2. 함수 진입 레이블 위치
- **증상**: `call id` 가 prologue를 건너뛰고 함수 바디 중간으로 점프
- **원인**: `MarkLabel(fn__entry)`가 prologue emit **이후** 블록 루프에서 호출됨
- **수정**: prologue **이전**에 `MarkLabel(fn__entry)` 등록, 블록 루프에서 entry 블록은 스킵
- **위치**: `codegen.go:compileFuncInto()`

## Phase 1 코드 구조

```
src/
  ir/
    types.go     — IR 타입 시스템 (i8~i64, f32/f64, ptr, array...)
    nodes.go     — IR 노드 (Value, Inst, Block, Function, Module)
  interp/
    interp.go    — 레퍼런스 인터프리터 (정확성 100% 목표)
  codegen/
    codegen.go   — x86-64 코드 생성기 (SSA → 기계어)
    elf.go       — ELF64 바이너리 생성기
  cmd/
    main.go          — 인터프리터 테스트 10개
    codegen_verify.go — 코드젠 검증 (interp == native)
```

## 실행 파이프라인

```
IR (SSA) → CompileModule → ELF64 → qemu-x86_64 → exit code
IR (SSA) → ExecModule   → RVal  → AsInt()
            양쪽 결과 비교 → 일치 확인
```

## 다음: Phase 2 우선순위

1. **FreeLang 파서 → Free IR 낙하** (ref frontend)
2. **phi 노드 codegen** (현재 interp-only)
3. **global 변수 codegen**
4. **struct/GEP codegen**
