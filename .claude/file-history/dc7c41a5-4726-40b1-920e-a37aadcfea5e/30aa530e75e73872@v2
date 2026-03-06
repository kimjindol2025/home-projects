# Phase 3: Bare-Metal OS - 10 Unforgiving Tests

**목표**: Rust 의존성 제거 (0%) + ARM64 bare-metal 검증
**철학**: "기록이 증명이다" — 정량 지표만 인정

---

## 📊 테스트 그룹 구성 (10개)

| 그룹 | 테스트 | 목표 | 지표 |
|------|--------|------|------|
| **A. Bootloader** | 4개 | ARM64 부트 순서 검증 | Stage progression, CPU state, Memory init |
| **B. MMU** | 3개 | 가상 메모리 관리 | Page table, VA translation, TLB |
| **C. Exception Handler** | 3개 | 예외/인터럽트 처리 | Vector table, Dispatch, GIC |

**총 합계**: 10개 무관용 테스트

---

## 🎯 그룹 A: Bootloader Tests (4개)

### A1: Boot Stage Progression
**목표**: Stage 1 → Stage 2 → Stage 3 정확한 진행?

```
Sequence:
Stage1_init() → memory_initialized=true
Stage2_transition() → el=1, gic_initialized=true
Stage3_kernel_init() → mmu_initialized=true, boot_complete()=true

Rule: 모든 스테이지 100% 순서대로 실행
```

**Unforgiving Rule 1**: 스테이지 순서 위반 = 실패

**검증**:
```rust
test_stage1() -> assert!(boot.stage1_init().is_ok())
test_stage2() -> boot.stage1_init(); assert_eq!(boot.cpu_state.el, 1)
test_stage3() -> assert!(boot.boot_complete())
```

---

### A2: CPU State Initialization
**목표**: CPU 레지스터(PC, SP, EL, etc.) 올바르게 초기화?

```
Requirements:
- pc = 0x80000000 (kernel load address)
- sp = 0x80200000 (stack)
- el = 3 (EL3 start)
- mpidr_el1 = valid CPU ID
- vbar_el1 = exception vector base

Rule: 모든 필드 정확한가?
```

**Unforgiving Rule 2**: CPU 상태 초기화 오류 = 실패

**검증**:
```rust
test_pc() -> assert_eq!(cpu.pc, 0x80000000)
test_stack() -> assert_eq!(cpu.sp, 0x80200000)
test_el_level() -> assert_eq!(cpu.el, 3)
test_vbar() -> assert_eq!(cpu.vbar_el1, 0x80100000)
```

---

### A3: Memory Layout Validation
**목표**: 메모리 레이아웃 충돌 없음?

```
Regions (순서대로):
- Kernel: 0x80000000-0x80200000 (2MB)
- Stack: 0x80200000-0x80300000 (1MB)
- Heap: 0x80300000-0xFF800000
- Device: 0xFF800000+ (MMIO)

Rule: 겹치거나 역순 없음?
```

**Unforgiving Rule 3**: 메모리 겹침 = 실패

**검증**:
```rust
test_layout_order() -> assert!(layout.kernel_base < layout.stack_base)
test_layout_validation() -> assert!(layout.validate().is_ok())
test_uart_address() -> assert_eq!(layout.uart_base, 0xFF800000)
```

---

### A4: GIC & Clock Initialization
**목표**: 인터럽트 컨트롤러와 클럭 초기화?

```
GIC:
- Distributor: 0xFF841000
- CPU interface: 0xFF842000
- gic_initialized = true

Clock:
- System clock: 800 MHz
- Initialization success

Rule: 모두 성공적으로 초기화?
```

**Unforgiving Rule 4**: GIC/Clock init 실패 = 실패

---

## 🛡️ 그룹 B: MMU Tests (3개)

### B1: Page Table Entry Conversion
**목표**: PTE를 ARM64 디스크립터로 올바르게 변환?

```
Conversion Rules:
- Valid bit (bit 0)
- Physical address (bits 47:12)
- Access permission (bits 9:8)
- Shareability (bits 8:7)
- Memory attribute (bits 5:2)

Rule: 모든 비트 정확한가?
```

**Unforgiving Rule 5**: PTE 변환 오류 = 실패

**검증**:
```rust
test_pte_valid_bit() -> assert_eq!(desc & 1, 1)
test_pte_addr_mask() -> assert_eq!(desc & 0xFFFFFFFFF000, physical_addr)
test_pte_roundtrip() -> pte == PTE::from_descriptor(pte.to_descriptor())
```

---

### B2: Virtual Address Translation
**목표**: VA를 PA로 올바르게 변환?

```
Page Table Walk:
L0 (bits [47:39]) → L1 (bits [38:30]) → L2 (bits [29:21]) → L3 (bits [20:12])
Offset: bits [11:0]

Rule: 모든 레벨 정확히 인덱싱?
```

**Unforgiving Rule 6**: VA 변환 오류 = 실패

**검증**:
```rust
test_va_l0_index() -> assert_eq!(va.get_l0_index(), (addr >> 39) & 0x1FF)
test_va_l3_index() -> assert_eq!(va.get_l3_index(), (addr >> 12) & 0x1FF)
test_canonical_va() -> VirtualAddress::new(0x0000000080000000).validate().is_ok()
test_non_canonical_va() -> VirtualAddress::new(0x5555555555555555).validate().is_err()
```

---

### B3: MMU Enable/Disable & Page Mapping
**목표**: MMU 제어와 페이지 매핑 작동?

```
Operations:
- map_page(va, pa, perm) → 페이지 매핑
- translate_va(va) → pa 반환
- unmap_page(va) → 매핑 제거

Rule: 모든 연산 성공 & 일관성?
```

**Unforgiving Rule 7**: 페이지 매핑 오류 = 실패

**검증**:
```rust
test_map_page() -> mmu.map_page(0x80001000, 0x80001000, ...).is_ok()
test_unmap_page() -> mmu.unmap_page(0x80001000).is_ok()
test_translate() -> translate_va(0x80001000) == 0x80001000 (identity mapping)
```

---

## 🚨 그룹 C: Exception Handler Tests (3개)

### C1: Exception Vector Table Construction
**목표**: 16개 벡터가 128바이트 정렬되고 올바르게 배치?

```
Vector Layout:
- Offset 0x000: Sync EL1
- Offset 0x080: IRQ EL1
- Offset 0x100: FIQ EL1
- Offset 0x180: SError EL1
- Offset 0x200-0x580: EL0 variants
- Offset 0x600-0x780: Reserved

Rule: 모두 128바이트 경계?
      주소 = base + offset?
```

**Unforgiving Rule 8**: 벡터 테이블 배치 오류 = 실패

**검증**:
```rust
test_vector_alignment() -> (vector.handler_address & 0x7F) == 0
test_vector_offsets() -> vector[0] = base + 0x000, vector[1] = base + 0x080
test_vector_count() -> vectors.len() == 16
test_table_validation() -> table.validate().is_ok()
```

---

### C2: Exception Dispatch
**목표**: 예외 타입별 올바른 핸들러 호출?

```
Exception Types:
- SynchronousAbort → handle_sync_abort()
- IRQ → handle_irq()
- FIQ → handle_fiq()
- SError → handle_serror()
- PageFault → handle_page_fault()

Rule: 모든 예외 올바르게 분류 & 디스패치?
```

**Unforgiving Rule 9**: 예외 디스패치 오류 = 실패

**검증**:
```rust
test_sync_abort() -> classify(ESR_0x20) == SynchronousAbort
test_irq_dispatch() -> handle_exception(irq_context).is_ok()
test_page_fault() -> handle_exception(pf_context) increments counter
test_exception_count() -> exceptions_handled increases
```

---

### C3: GIC Interrupt Management
**목표**: GIC 인터럽트 활성화/비활성화/우선순위 설정?

```
Operations:
- init() → irq_count = 32
- enable_irq(n) → IRQ n 활성화 (0 ≤ n < 32)
- disable_irq(n) → IRQ n 비활성화
- set_priority(n, p) → 우선순위 설정 (0-255)
- get_irq() → 펜딩 IRQ 반환

Rule: 범위 검사 & 모든 연산 성공?
```

**Unforgiving Rule 10**: GIC 연산 오류 = 실패

**검증**:
```rust
test_gic_init() -> gic.init().is_ok(); assert_eq!(irq_count, 32)
test_enable_valid() -> enable_irq(0).is_ok()
test_enable_invalid() -> enable_irq(100).is_err()
test_priority() -> set_priority(0, 128).is_ok()
test_bounds_check() -> all invalid ranges rejected
```

---

## 🔐 10가지 Unforgiving Rules 요약

| # | 규칙 | 대상 | 실패 기준 |
|---|------|------|----------|
| 1 | 스테이지 진행 | Bootloader | 순서 위반 > 0 |
| 2 | CPU 상태 초기화 | Bootloader | 초기화 오류 > 0 |
| 3 | 메모리 레이아웃 | Bootloader | 겹침 > 0 |
| 4 | GIC/Clock init | Bootloader | init 실패 > 0 |
| 5 | PTE 변환 | MMU | 변환 오류 > 0 |
| 6 | VA 변환 | MMU | 변환 오류 > 0 |
| 7 | 페이지 매핑 | MMU | 매핑 오류 > 0 |
| 8 | 벡터 테이블 | Exception | 배치 오류 > 0 |
| 9 | 예외 디스패치 | Exception | 디스패치 오류 > 0 |
| 10 | GIC 관리 | Exception | 관리 오류 > 0 |

---

## ✅ 검증 기준

**모든 10개 테스트 = 100% 통과 (10가지 Unforgiving Rule 모두 만족)**

```
Phase 3 Success = (Passed Tests / Total Tests == 100%) AND (All Unforgiving Rules == Satisfied)
```

---

**목표**:
- 2,200줄 코드 (3개 모듈)
- 10개 무관용 테스트 (10가지 규칙)
- Rust 의존성 0% 달성

**철학**: "기록이 증명이다" (Your Record is Your Proof)

---

## 📁 Phase 3 파일 구조

```
freelang-llc/src/
├── bootloader.fl (900줄)
│   ├── ARM64 CPU state & modes
│   ├── Boot stages (1-3)
│   ├── Assembly intrinsics wrapper
│   ├── Memory layout initialization
│   └── Device init (UART, Clock)
├── mmu.fl (650줄)
│   ├── Page table entries
│   ├── Page table levels (L0-L3)
│   ├── Virtual address translation
│   ├── MMU management
│   └── 8개 테스트
├── exception_handler.fl (650줄)
│   ├── Exception types & context
│   ├── Exception vector table
│   ├── Handler dispatcher
│   ├── GIC (interrupt controller)
│   └── 9개 테스트
```

**총 코드**: 2,200줄
**총 테스트**: 10개 무관용 + 내부 단위테스트 30개

---

## 🎯 최종 목표

**FreeLang-LLC 완성 (Phase 1+2+3)**:
- Phase 1: 1,500줄 (ManagedPointer)
- Phase 2: 2,200줄 (LLVM Backend)
- Phase 3: 2,200줄 (Bare-Metal OS)
- **총합: 5,900줄**

**테스트**:
- Phase 1: 25개
- Phase 2: 35개
- Phase 3: 10개
- **총합: 70개 무관용 테스트**

**언어독립성**: 32.2% → 85%+ 달성
**Rust 의존도**: 15.4% → 0%

---

**Version**: 3.0
**Status**: Phase 3 Design Complete
**Philosophy**: "기록이 증명이다"
