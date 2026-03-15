# 🦎 Zig 미니 컴파일러 구현 프로젝트 (2026-03-12 23:00 UTC+9)

**상태**: ✅ Phase 6: LLVM IR → Native Binary 완료 (총 ~15,800줄 코드, 80% 완성도)

## 진행도

```
Phase 1: 준비 ........................... ✅ 완료
Phase 2.1: Tokenizer .................... ✅ 완료 (446줄)
Phase 2.2: Parser ....................... ✅ 완료 (972줄)
Phase 2.3: AstGen ..................... ✅ 완료 (581줄)
Phase 2.4: Sema (ZIR→AIR) ............ ✅ 완료 (571줄)
Phase 2.5-2.7: Codegen/Linker ........ ✅ 완료

Phase 3.1-3.8: Advanced Features ...... ✅ 완료
  └─ Generic, Union/Enum, Comptime, Overload
  └─ Arrays, Defer, Error, Pattern Matching

Phase 4: C Code Generation ............ ✅ 완료
  └─ Pattern Matching → switch/if-else if
  └─ Error Handling → alloca 기반 에러변수
  └─ codegen.zig 생성

Phase 5: LLVM IR Backend ............. ✅ 완료
  └─ Step 1: Basic Values & Functions (9 emit functions)
  └─ Step 2: Control Flow (10 emit functions)
  └─ Step 3: Advanced Features (8 emit functions)
  └─ llvm_ir.zig 생성 (1,070줄)
  └─ 39 tests: basic_ir, control_flow_ir, advanced_ir

Phase 6: LLVM IR → Native Binary .... ✅ 완료 🎉
  └─ Step 1: 버그 수정 + linker/main 통합
     • llvm_ir.zig: emitTryExpr() 버그 수정 (2줄)
     • linker.zig: LLVM 메서드 추가 (+60줄)
     • main.zig: CLI 플래그 + Step 5b/6b (+55줄)
  └─ Step 2: Phase 6 테스트 (13 tests)
     • Group A: 파일 I/O (3)
     • Group B: 헬퍼 함수 (3)
     • Group C: IR 생성+저장 (3)
     • Group D: 컴파일러 감지 (2)
     • Group E: 통합 (2)

누적: ~15,800줄 (src ~3,200줄 + test ~12,600줄)
테스트: 366개 (353개 기존 + 13개 Phase 6)
완성도: ~80%
```

## Phase 3.5 Step 1 (Parser) - 완료 내용

### 파일 변경사항
1. **tokenizer.zig** (+20줄)
   - `dot_dot` 토큰 추가
   - `..` 시퀀스 인식

2. **ast.zig** (+14줄)
   - `type_slice` 노드 태그
   - `slice_access` 노드 태그
   - `slice_type` 데이터 구조
   - `slice_access` 데이터 구조

3. **parser.zig** (+110줄)
   - `ParsingError` 타입 정의
   - `parseType()` 슬라이스/배열 구분 로직
   - `parsePostfix()` 슬라이싱 처리

4. **test/phase3_array_parser_test.zig** (신규 380줄)
   - 15개 테스트 (6개 그룹)

5. **build.zig** (+8줄)
   - `array_parser_tests` 타겟 추가

### 핵심 기능
- ✅ `[]T` 슬라이스 타입 인식
- ✅ `[]const T` 상수 슬라이스 인식
- ✅ `[N]T` 배열 타입 인식
- ✅ `arr[i]` 인덱싱 파싱
- ✅ `arr[i..j]` 슬라이싱 파싱 (선택적 시작/끝)

---

## Phase 3.5 Step 2 (AstGen) - 완료 내용

### 파일 변경사항
1. **zir.zig** (+50줄)
   - `array_type`: 배열/슬라이스 타입 ZIR
     - `element_type: Ref`
     - `size: ?u32` (null이면 슬라이스)
     - `is_const: bool`
   - `array_access`: 배열 인덱싱
   - `slice_access`: 슬라이싱 (선택적 start/end)
   - `array_length`: 배열 길이
   - `printZir()` 출력 로직 추가

2. **astgen.zig** (+105줄)
   - `genArrayType()`: [N]T 타입 생성
   - `genSliceType()`: []T 타입 생성
   - `genSliceAccess()`: arr[i..j] 슬라이싱
   - `genNode()` switch에 3개 케이스

3. **test/phase3_array_astgen_test.zig** (신규 360줄)
   - 12개 테스트 (5개 그룹)

4. **build.zig** (+10줄)
   - `array_astgen_tests` 타겟 추가

### 핵심 기능
- ✅ 배열 타입 ZIR 생성
- ✅ 슬라이스 타입 ZIR 생성
- ✅ 슬라이싱 ZIR 생성
- ✅ 선택적 인덱스 처리
- ✅ 중첩 배열 타입 지원

## Phase 6 완료 내용 (LLVM IR → Native Binary)

### 파일 변경사항
1. **src/llvm_ir.zig** (버그 수정, 2줄)
   - 라인 519: 문자열 `+` 연산 → writer().print() 수정
   - 라인 531: 문자열 `+` 연산 → writer().print() 수정
   - emitTryExpr()의 모든 Zig 비호환 문법 해결

2. **src/linker.zig** (+60줄)
   - `llvm_path: ?[]const u8` 필드 추가
   - `init()` 수정: findLLVMCompiler() 호출
   - `deinit()` 수정: llvm_path free 추가
   - `findLLVMCompiler()`: clang/cc 자동 감지
   - `writeIRToFile()`: IR 텍스트 → .ll 파일 저장
   - `compileFromLLVMIR()`: clang으로 .ll → 바이너리
   - `compileIRText()`: 통합 메서드

3. **src/main.zig** (+55줄)
   - `LLVMIRGen` 임포트 추가
   - CLI 플래그: `--emit-llvm`, `--llvm`
   - LLVM 경로 분기 (Step 5b/6b)
   - C 경로와 LLVM 경로 선택적 실행

4. **test/phase6_linker_test.zig** (신규 300줄, 13 tests)
   - Group A: 파일 I/O (3 tests)
   - Group B: 헬퍼 함수 (3 tests)
   - Group C: IR 생성+저장 (3 tests)
   - Group D: 컴파일러 감지 (2 tests)
   - Group E: 통합 (2 tests)

5. **build.zig** (+12줄)
   - phase6_linker_tests 타겟 추가
   - run_phase6_linker_tests 아티팩트 추가
   - test_step 의존성 추가

### 핵심 기능
✅ LLVM IR 텍스트를 .ll 파일로 저장
✅ clang을 통해 .ll → 네이티브 바이너리 컴파일
✅ --emit-llvm: IR 텍스트만 생성
✅ --llvm: LLVM IR 경유 바이너리 컴파일
✅ C 경로 유지 (기본값, --emit-c 플래그 지원)
✅ clang 자동 감지 (없으면 graceful skip)

## 다음 단계

### Phase 7: 예정 사항
- 리얼 바이너리 실행/테스트
- 최적화 패스 (DCE, 상수 folding)
- 크로스 컴파일 지원 (ARM, RISC-V)
- 고급 기능 (직렬화, 매크로 시스템)

## 주요 기술 결정

### 1. 토크나이저
- `dot_dot` 토큰: `..` 2-문자 시퀀스

### 2. Parser - 슬라이스 vs 배열 구분
```zig
// parseType()에서:
if (self.check(.rbracket) or (self.check(.const_) and peekNext == ']')) {
    // 슬라이스: []T, []const T
} else {
    // 배열: [N]T
}
```

### 3. Parser - 인덱싱 vs 슬라이싱 구분
```zig
// parsePostfix()에서:
if (self.match(.dot_dot)) {
    // 슬라이싱: arr[i..j]
} else {
    // 인덱싱: arr[i]
}
```

### 4. 선택적 인덱스
- `start: ?*Node` - 없으면 처음부터
- `end: ?*Node` - 없으면 끝까지

## 문서
- PHASE_3_5_PLAN.md - 전체 계획
- PHASE_3_5_STEP1_PROGRESS.md - Step 1 완료 보고서
