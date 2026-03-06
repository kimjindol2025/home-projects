# Phase 2: LLVM Backend - 35 Unforgiving Tests

**목표**: 5배 성능 향상 (900µs → 180µs) 달성 검증
**철학**: "숫자가 증명이다" — 정량 지표만 인정

---

## 📊 테스트 그룹 구성 (35개)

| 그룹 | 테스트 | 목표 | 지표 |
|------|--------|------|------|
| **A. LLVM Codegen** | 12개 | IR 생성 정확성 | Instruction count, Type correctness, Optimization |
| **B. Assembly Parser** | 12개 | asm! 파싱 및 검증 | Constraint parsing, SIMD validation |
| **C. Semantic Validator** | 11개 | 최적화 안정성 | CFG preservation, Data flow integrity |

**총 합계**: 35개 무관용 테스트 (각 100% 통과 필수)

---

## 🎯 그룹 A: LLVM Codegen Tests (12개)

### A1: Type System Correctness
**목표**: 모든 FreeLang 타입이 올바른 LLVM 타입으로 변환되는가?

```
Input:  u8, u16, u32, u64, u128, f32, f64
Output: i8, i16, i32, i64, i128, float, double
Rule: 100% 정확도 (1개 오류 = 실패)
```

**검증 방식**:
```rust
test_llvm_type_i8() -> assert_eq!(LLVMType::I8.to_llvm_string(), "i8")
test_llvm_type_i64() -> assert_eq!(LLVMType::I64.to_llvm_string(), "i64")
test_llvm_type_pointer() -> assert_eq!(LLVMType::Ptr.to_llvm_string(), "i8*")
test_llvm_type_array() -> assert!(LLVMType::Array(...).to_llvm_string().contains("["))
test_llvm_type_struct() -> assert!(LLVMType::Struct(...).to_llvm_string().contains("{"))
```

**Unforgiving Rule 1**: 타입 변환 오류 = 0 (허용불가)

---

### A2: Instruction Generation Accuracy
**목표**: 25개 LLVM 명령어 모두 올바른 IR 문법 생성하는가?

```
Categories:
- Arithmetic: add, sub, mul, div (4)
- Bitwise: and, or, xor, shl, shr (5)
- Memory: load, store, getelementptr (3)
- Control: br, ret, call (3)
- Comparison: icmp, fcmp (2)
- Conversion: bitcast, zext, sext (3)
- Volatile: volatile load/store (2)
- Barrier: fence (1)
Total: 23 instructions
```

**Unforgiving Rule 2**: 명령어 문법 오류 = 0 (허용불가)

**검증**:
```
test_inst_arithmetic() -> Check add, sub, mul generate valid syntax
test_inst_bitwise() -> Check and, or, xor, shl, shr
test_inst_memory() -> Check load, store, getelementptr
test_inst_control() -> Check br, ret, call
test_inst_volatile() -> Volatile ops include memory barriers
test_inst_barrier() -> Fence instructions present
```

---

### A3: Basic Block Grouping
**목표**: 명령어가 기본 블록으로 올바르게 그룹화되는가?

```
Rule: Each block has entry, instructions[], successors
- Block ID must be unique
- Instructions must be sequential
- Successors must be valid blocks
```

**Unforgiving Rule 3**: 기본 블록 구조 위반 = 0 (허용불가)

---

### A4: Function Signature Generation
**목표**: 함수 서명이 LLVM 형식으로 올바르게 생성되는가?

```
Format: define <return_type> @<name>(<param_list>)
Example: define i64 @safe_read(i8* %ptr, i64 %offset)

Rule: Return type + parameters 정확한가?
```

**Unforgiving Rule 4**: 함수 서명 오류 = 0

---

### A5: Module Structure Validity
**목표**: 전체 LLVM 모듈이 유효한 IR 문법인가?

```
Structure:
- target datalayout = "..."
- target triple = "..."
- declare void @llvm.mfence()
- define ... @function() { ... }
```

**Unforgiving Rule 5**: 모듈 문법 오류 = 0

---

### A6: Optimization Level Mapping
**목표**: 5가지 최적화 레벨(O0/O1/O2/O3/Os)이 올바른 LLVM 패스로 매핑되는가?

```
O0 -> -O0 (no optimization)
O1 -> -O1 (minimal)
O2 -> -O2 (moderate, recommended)
O3 -> -O3 (aggressive)
Os -> -Os (size optimize)
```

**검증**:
```
test_opt_o0() -> assert_eq!(OptimizationLevel::O0.to_llc_flag(), "-O0")
test_opt_o2() -> assert_eq!(OptimizationLevel::O2.to_llc_flag(), "-O2")
test_opt_estimated_speedup() -> estimated_speedup(O2) > estimated_speedup(O0)
```

**Unforgiving Rule 6**: 최적화 레벨 매핑 오류 = 0

---

### A7: Code Generation Example (add function)
**목표**: `generate_add_function()` 생성 코드가 유효한가?

```
Expected:
- Function name: @add_func
- Parameters: i64 %a, i64 %b
- Operation: add i64 %a, %b
- Return: ret i64 %result
```

**검증**: 생성된 코드가 유효한 LLVM IR 문법 확인

---

### A8: Code Generation Example (safe_read)
**목표**: `generate_safe_read_function()` 생성 코드가 bounds checking을 포함하는가?

```
Expected Operations:
1. Load from pointer
2. Type conversion (if needed)
3. Return value
4. Include bounds validation logic
```

---

### A9: Volatile Read Function
**목표**: `generate_volatile_read_function()` 메모리 배리어를 포함하는가?

```
Expected:
- fence seq_cst (before)
- volatile load
- fence seq_cst (after)
```

**Unforgiving Rule 7**: volatile 읽기에서 배리어 누락 = 실패

---

### A10: Fibonacci Function Generation
**목표**: 재귀 함수 생성이 올바른가?

```
test_fib_structure() -> assert has basic blocks for recursion
test_fib_termination() -> assert base case (n <= 1) exists
test_fib_recursion() -> assert recursive calls present
```

---

### A11: Type Casting Accuracy
**목표**: 타입 변환(zext, sext, bitcast) 명령어가 올바른가?

```
test_zext() -> zero extend i32 to i64
test_sext() -> sign extend i32 to i64
test_bitcast() -> reinterpret bits without change
```

---

### A12: Memory Ordering Preservation
**목표**: 메모리 접근 순서가 유지되는가?

```
Pattern:
  load i64* %ptr1
  load i64* %ptr2
  store i64 %val1, i64* %ptr3

Rule: Sequential order must match source
```

---

## 🔧 그룹 B: Assembly Parser Tests (12개)

### B1: Lexer Token Recognition
**목표**: 모든 토큰이 올바르게 인식되는가?

```
Tokens: asm!, {, }, :, ;, "string", constraint, register, modifier, number
Rule: 각 토큰 유형 100% 정확
```

---

### B2: String Template Parsing
**목표**: asm! { "template string" } 파싱 정확한가?

```
test_template_parsing() -> "mov $0, $1" 정확하게 파싱
test_template_with_escapes() -> "\\t", "\\n" 처리
```

---

### B3: Constraint Parser - Input
**목표**: "r" 입력 제약 조건 파싱 정확한가?

```
Input: "r", "=r", "m", "i"
Expected: ConstraintType::Input, ConstraintType::Output, etc.
Rule: 100% 정확도
```

---

### B4: Constraint Parser - Output
**목표**: "=r" 출력 제약 조건 파싱 정확한가?

```
test_output_reg() -> assert_eq!(constraint_type, Output)
test_output_memory() -> assert_eq!(constraint_type, Output)
```

---

### B5: SIMD Register Constraints
**목표**: "x" (XMM), "y" (YMM), "z" (ZMM) 제약 파싱 정확한가?

```
test_simd_xmm() -> "x" → register_class = "x"
test_simd_ymm() -> "y" → register_class = "y"
test_simd_zmm() -> "z" → register_class = "z"
```

**Unforgiving Rule 8**: SIMD 제약 오류 = 0

---

### B6: Early Clobber Detection
**목표**: "&" 수정자가 올바르게 파싱되는가?

```
Input: "&=r"
Expected: is_early_clobber = true, constraint_type = Output
```

---

### B7: LLVM Constraint Translation
**목표**: 파싱된 제약이 LLVM 형식으로 올바르게 변환되는가?

```
FreeLang: "=r" → LLVM: "=r"
FreeLang: "=x" → LLVM: "=x"
FreeLang: "&=m" → LLVM: "&=m"
```

---

### B8: Complete AST Parsing
**목표**: 전체 asm! 블록이 InlineAsmBlock으로 파싱되는가?

```
Input: asm! { "mov $0, $1" : "=r"(result) : "r"(input) : : }
Output: InlineAsmBlock { template, outputs, inputs, clobbers }
```

---

### B9: SSE Instruction Validation
**목표**: SSE 명령어(movdqa, paddb, etc.) 검증 정확한가?

```
test_movdqa() -> assert!(validate("movdqa xmm0, xmm1").is_ok())
test_paddb() -> assert!(validate("paddb xmm0, xmm1, xmm2").is_ok())
test_invalid_sse() -> assert!(validate("invalid_sse").is_err())
```

**Unforgiving Rule 9**: SSE 검증 오류 = 0

---

### B10: AVX Instruction Validation
**목표**: AVX 명령어(vmovdqa, vpaddb, etc.) 검증 정확한가?

```
test_vmovdqa() -> assert!(validate("vmovdqa ymm0, ymm1").is_ok())
test_vpaddb() -> assert!(validate("vpaddb ymm0, ymm1, ymm2").is_ok())
```

---

### B11: AVX-512 Instruction Validation
**목표**: AVX-512 명령어 검증 정확한가?

```
test_vpmullq() -> assert!(validate("vpmullq zmm0, zmm1, zmm2").is_ok())
```

---

### B12: Invalid Instruction Rejection
**목표**: 잘못된 명령어 거부하는가?

```
test_typo_detection() -> assert!(validate("mov rax, ").is_err())
test_unknown_mnemonic() -> assert!(validate("invalid_op").is_err())
```

**Unforgiving Rule 10**: 잘못된 명령어 허용 = 실패

---

## 🛡️ 그룹 C: Semantic Validator Tests (11개)

### C1: Data Flow Analysis
**목표**: 데이터 흐름 분석이 변수 정의/사용 추적하는가?

```
test_variable_definition() -> add_definition("x", 0) works
test_variable_use() -> add_use("x", 1) works
test_use_before_define() -> detect error if use before define
```

---

### C2: Control Flow Graph Construction
**목표**: CFG가 올바르게 구성되는가?

```
test_cfg_nodes() -> blocks added correctly
test_cfg_edges() -> edges connect valid blocks
test_cfg_entry() -> entry block identified
```

---

### C3: Reachability Analysis
**목표**: 도달 가능성 분석이 정확한가?

```
test_reachability_from_entry() -> compute_reachability() finds all reachable blocks
test_dead_code_detection() -> find_dead_code() identifies unreachable blocks
Rule: Unreachable blocks = 0 (except designated dead code)
```

**Unforgiving Rule 11**: 도달 불가능 블록 누락 = 실패

---

### C4: Use-Before-Define Detection
**목표**: use-before-define 오류 감지하는가?

```
test_detect_ubd() -> add_use("x", 0); add_definition("x", 1) → Err
test_valid_sequence() -> add_definition("x", 0); add_use("x", 1) → Ok
```

---

### C5: Memory Access Order Preservation
**목표**: 메모리 접근 순서 변경 감지하는가?

```
Before: load ptr[0], load ptr[1], store ptr[2]
After:  load ptr[1], load ptr[0], store ptr[2]  ← INVALID

test_memory_reordering() -> detect this violation
```

---

### C6: Volatile Operation Preservation
**목표**: volatile 연산 제거 감지하는가?

```
Before: volatile_read(&ptr, 0)
After:  (operation removed)

test_volatile_elimination() -> assert_eq!(before.side_effects.len(), after.side_effects.len())
```

---

### C7: Barrier Preservation Check
**목표**: 메모리 배리어 제거 감지하는가?

```
Before: fence seq_cst
After:  (fence removed)

test_barrier_removal() -> check_invariants() returns Err("barriers removed")
```

**Unforgiving Rule 12**: 배리어 제거 허용 = 실패

---

### C8: Instruction Count Sanity Check
**목표**: 명령어 폭증 감지하는가?

```
Before: 10 instructions
After:  25 instructions (2.5x increase)

Rule: After ≤ Before × 2 (합리적 컴파일러 한계)
test_instruction_explosion() -> after > before * 2 → Err
```

**Unforgiving Rule 13**: 지나친 명령어 증가 = 실패

---

### C9: Register Allocation Validation
**목표**: 레지스터 할당 폭증 감지하는가?

```
Before: 5 registers used
After:  20 registers used (unrealistic)

Rule: After ≤ Before + 10
test_register_explosion() -> detect this and fail
```

---

### C10: Control Flow Preservation
**목표**: 제어 흐름 그래프 변경 감지하는가?

```
Before CFG: entry → {block1, block2} → exit
After CFG:  entry → block1 → exit (block2 removed!)

test_cfg_preservation() -> check_equivalence() detects this
```

---

### C11: Semantic Equivalence Summary
**목표**: 모든 의미론적 검증 통과하는가?

```
Requirements:
✓ Data flow preserved (variables exist)
✓ Memory access order preserved
✓ Side effects not reordered
✓ Control flow intact
✓ Barriers present
✓ Instructions reasonable count
✓ Registers reasonable count

test_full_equivalence_check() -> all checks pass
```

---

## 📈 8가지 Unforgiving Rules 요약

| # | 규칙 | 대상 | 실패 기준 |
|---|------|------|----------|
| 1 | 타입 변환 정확도 | LLVM Codegen | 오류 > 0 |
| 2 | 명령어 문법 정확도 | LLVM Codegen | 문법 오류 > 0 |
| 3 | 기본 블록 구조 | LLVM Codegen | 구조 위반 > 0 |
| 4 | 함수 서명 정확도 | LLVM Codegen | 서명 오류 > 0 |
| 5 | 모듈 구조 유효성 | LLVM Codegen | 문법 오류 > 0 |
| 6 | 최적화 레벨 매핑 | LLVM Codegen | 매핑 오류 > 0 |
| 7 | volatile 배리어 | LLVM Codegen | 배리어 누락 > 0 |
| 8 | SIMD 제약 검증 | Assembly Parser | 제약 오류 > 0 |
| 9 | SSE/AVX 검증 | Assembly Parser | 검증 오류 > 0 |
| 10 | 명령어 거부 | Assembly Parser | 잘못된 명령어 허용 > 0 |
| 11 | 도달 가능성 | Semantic Validator | 누락된 블록 > 0 |
| 12 | 배리어 보존 | Semantic Validator | 배리어 제거 > 0 |
| 13 | 명령어 폭증 | Semantic Validator | 증가율 > 200% |

---

## ✅ 검증 기준

**모든 35개 테스트 = 100% 통과 (13가지 Unforgiving Rule 모두 만족)**

```
Phase 2 Success = (Passed Tests / Total Tests == 100%) AND (All Unforgiving Rules == Satisfied)
```

---

**목표**:
- 2,200줄 코드 (3개 모듈)
- 35개 무관용 테스트 (13가지 규칙)
- 5배 성능 향상 검증

**철학**: "기록이 증명이다" (Your Record is Your Proof)
