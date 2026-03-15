# FreeLang OS Kernel: 실제 동작 구현 로드맵

**목표**: 시뮬레이션 → 실제 x86-64 OS (QEMU에서 부팅 및 동작)

---

## 📋 Phase 0: 현재 상태 분석

### ✅ 이미 완성된 것
- bootloader.fl: 메모리 관리 알고리즘 (물리/가상/힙)
- kernel.fl: 프로세스 관리 구조
- scheduler.fl: Context switching 로직
- interrupt.fl: IDT 구조

### ❌ 부족한 것
- 실제 x86-64 부트로더 (GRUB/multiboot 호환)
- 실제 메모리 레이아웃 설정
- 실제 GDT/IDT 어셈블리
- 실제 인터럽트 처리 루틴 (어셈블리)
- Context switching 어셈블리 코드
- 실제 하드웨어 타이머 상호작용
- QEMU 시뮬레이터 설정

---

## 🏗️ Phase 1: Bootloader 재구성 (1-2주)

### 1.1 Multiboot2 Bootloader 작성

**목표**: QEMU에서 부팅 가능하도록 만들기

**파일**: `src/boot.s` (x86-64 어셈블리)

```asm
; src/boot.s - x86-64 Bootloader
[bits 32]  ; 32-bit 모드에서 시작

MULTIBOOT2_HEADER_MAGIC     equ 0xe85250d6
MULTIBOOT_ARCHITECTURE_I386 equ 0
MULTIBOOT2_HEADER_TAG_END   equ 0

section .multiboot_header
align 8
    dd MULTIBOOT2_HEADER_MAGIC
    dd MULTIBOOT_ARCHITECTURE_I386
    dd (header_end - header_start)
    dd -(MULTIBOOT2_HEADER_MAGIC + MULTIBOOT_ARCHITECTURE_I386 + (header_end - header_start))

    ; 태그 (선택사항)
    dw 0
    dw 0
    dd 8

header_start:
header_end:

section .text
extern long_mode_start

global _start
_start:
    mov esp, stack_top
    call check_multiboot
    call check_cpuid
    call check_long_mode
    call setup_page_tables
    call enable_paging

    ; 64-bit 모드 전환
    lgdt [gdt64.pointer]
    jmp gdt64.code:long_mode_start

check_multiboot:
    cmp eax, 0x36d76289
    jne .no_multiboot
    ret
.no_multiboot:
    mov al, "M"
    jmp error

check_cpuid:
    pushfd
    pop eax
    mov ecx, eax
    xor eax, 0x200000
    push eax
    popfd
    pushfd
    pop eax
    cmp eax, ecx
    je .no_cpuid
    ret
.no_cpuid:
    mov al, "C"
    jmp error

check_long_mode:
    mov eax, 0x80000000
    cpuid
    cmp eax, 0x80000001
    jb .no_long_mode

    mov eax, 0x80000001
    cpuid
    test edx, 1 << 29
    jz .no_long_mode
    ret
.no_long_mode:
    mov al, "L"
    jmp error

setup_page_tables:
    mov eax, p4_table
    or eax, 0b11
    mov [p4_table], eax

    mov eax, p3_table
    or eax, 0b11
    mov [p4_table + 0], eax

    mov eax, p2_table
    or eax, 0b11
    mov [p3_table + 0], eax

    ; 2MB 페이지 설정 (0x200000부터 시작)
    mov ecx, 0
.loop:
    mov eax, 0x200000
    add eax, ecx
    shl ecx, 21
    or eax, 0b10000011  ; 쓰기 가능, 존재, 큰 페이지
    mov [p2_table + ecx], eax

    inc ecx
    cmp ecx, 512
    jne .loop

    ret

enable_paging:
    mov eax, p4_table
    mov cr3, eax

    mov eax, cr4
    or eax, 1 << 5
    mov cr4, eax

    mov ecx, 0xC0000080
    rdmsr
    or eax, 1 << 8
    wrmsr

    mov eax, cr0
    or eax, 1 << 31
    or eax, 1 << 16
    mov cr0, eax

    ret

error:
    mov dword [0xb8000], 0x4f524f45
    mov dword [0xb8004], 0x4f3a4f52
    mov dword [0xb8008], 0x4f204f20
    mov byte [0xb800a], al
    hlt

section .bss
align 4096
p4_table: resq 512
p3_table: resq 512
p2_table: resq 512
stack_bottom:
    resb 4096 * 4
stack_top:

section .rodata
gdt64:
    dq 0
.code: equ $ - gdt64
    dq (1<<43) | (1<<47) | (1<<53)
.pointer:
    dw $ - gdt64 - 1
    dq gdt64
```

### 1.2 Long Mode 진입점

**파일**: `src/long_mode_init.asm`

```asm
; src/long_mode_init.asm
[bits 64]

extern kernel_main

global long_mode_start
long_mode_start:
    mov rax, kernel_main
    call rax
    hlt
```

### 1.3 Cargo.toml 설정

```toml
[package]
name = "freelang-kernel"
version = "0.1.0"
edition = "2021"

[dependencies]
x86_64 = "0.14"
volatile = "0.4"

[profile.dev]
panic = "abort"
lto = true

[[bin]]
name = "kernel"
path = "src/main.rs"

[build]
target = "x86_64-unknown-linux-gnu"
```

---

## 🔧 Phase 2: 실제 Memory Management 구현 (1주)

### 2.1 물리 메모리 관리자 재구현

**파일**: `src/memory/physical.rs`

```rust
// src/memory/physical.rs
use core::mem;

const PAGE_SIZE: usize = 4096;
const MEMORY_SIZE: usize = 512 * 1024 * 1024;  // 512MB
const NUM_PAGES: usize = MEMORY_SIZE / PAGE_SIZE;

pub struct PhysicalMemoryManager {
    bitmap: &'static mut [u8],  // 비트맵: 1=할당, 0=free
}

impl PhysicalMemoryManager {
    pub fn new(bitmap_addr: usize) -> Self {
        let bitmap_size = (NUM_PAGES + 7) / 8;
        let bitmap = unsafe {
            core::slice::from_raw_parts_mut(bitmap_addr as *mut u8, bitmap_size)
        };

        // 커널 메모리 영역 표시
        for i in 0..256 {  // 1MB 커널 영역
            bitmap[i / 8] |= 1 << (i % 8);
        }

        PhysicalMemoryManager { bitmap }
    }

    pub fn allocate_page(&mut self) -> Option<usize> {
        for byte_idx in 0..(NUM_PAGES / 8) {
            for bit in 0..8 {
                if (self.bitmap[byte_idx] & (1 << bit)) == 0 {
                    self.bitmap[byte_idx] |= 1 << bit;
                    return Some((byte_idx * 8 + bit) * PAGE_SIZE);
                }
            }
        }
        None
    }

    pub fn deallocate_page(&mut self, addr: usize) {
        let page_num = addr / PAGE_SIZE;
        let byte_idx = page_num / 8;
        let bit = page_num % 8;
        self.bitmap[byte_idx] &= !(1 << bit);
    }
}
```

### 2.2 가상 메모리 (페이지 테이블)

**파일**: `src/memory/virtual.rs`

```rust
// src/memory/virtual.rs
use x86_64::registers::control::Cr3;
use x86_64::structures::paging::{PageTable, OffsetPageTable, Mapper};

pub struct VirtualMemoryManager {
    mapper: OffsetPageTable<'static>,
}

impl VirtualMemoryManager {
    pub unsafe fn new(physical_memory_offset: u64) -> Self {
        let level_4_table = &mut *((Cr3::read().0.start_address().as_u64()
            + physical_memory_offset) as *mut PageTable);

        VirtualMemoryManager {
            mapper: OffsetPageTable::new(level_4_table, VirtAddr::new(physical_memory_offset)),
        }
    }
}
```

### 2.3 힙 할당자

**파일**: `src/memory/heap.rs`

```rust
// src/memory/heap.rs
use alloc::alloc::{GlobalAlloc, Layout};
use core::ptr::NonNull;

pub struct HeapAllocator;

unsafe impl GlobalAlloc for HeapAllocator {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        // First-Fit 구현
        // bootloader.fl의 HeapAllocator 로직 이용
        todo!()
    }

    unsafe fn dealloc(&self, ptr: *mut u8, layout: Layout) {
        todo!()
    }
}

#[global_allocator]
static HEAP: HeapAllocator = HeapAllocator;
```

---

## ⚡ Phase 3: 실제 인터럽트 처리 (1주)

### 3.1 GDT (Global Descriptor Table) 설정

**파일**: `src/gdt.rs`

```rust
// src/gdt.rs
use x86_64::structures::gdt::{GlobalDescriptorTable, Descriptor};
use x86_64::structures::tss::TaskStateSegment;

pub fn init_gdt() {
    let mut gdt = GlobalDescriptorTable::new();
    let code = gdt.add_entry(Descriptor::kernel_code_segment());
    let tss_selector = gdt.add_entry(Descriptor::tss_segment(&TSS));
    gdt.load();
}

lazy_static::lazy_static! {
    static ref TSS: TaskStateSegment = {
        let mut tss = TaskStateSegment::new();
        tss.interrupt_stack_table[DOUBLE_FAULT_IST_INDEX as usize] = {
            const STACK_SIZE: usize = 4096 * 5;
            static mut STACK: [u8; STACK_SIZE] = [0; STACK_SIZE];

            VirtAddr::from_ptr(unsafe { &STACK }) + STACK_SIZE
        };
        tss
    };
}
```

### 3.2 IDT (Interrupt Descriptor Table) 설정

**파일**: `src/interrupts.rs`

```rust
// src/interrupts.rs
use x86_64::structures::idt::{InterruptDescriptorTable, InterruptStackFrame};

lazy_static::lazy_static! {
    static ref IDT: InterruptDescriptorTable = {
        let mut idt = InterruptDescriptorTable::new();
        idt.breakpoint.set_handler_fn(breakpoint_handler);
        idt.double_fault.set_handler_fn(double_fault_handler);
        idt.page_fault.set_handler_fn(page_fault_handler);
        idt[InterruptIndex::Timer.as_u8() as usize]
            .set_handler_fn(timer_interrupt_handler);
        idt
    };
}

pub fn init() {
    IDT.load();
}

extern "x86-interrupt" fn breakpoint_handler(stack_frame: InterruptStackFrame) {
    println!("EXCEPTION: BREAKPOINT\n{:#?}", stack_frame);
}

extern "x86-interrupt" fn page_fault_handler(
    stack_frame: InterruptStackFrame,
    error_code: PageFaultErrorCode,
) {
    println!("EXCEPTION: PAGE FAULT");
    println!("Accessed Address: {:?}", Cr2::read());
    println!("Error Code: {:?}", error_code);
    println!("{:#?}", stack_frame);

    // 실제 메모리 할당 수행
    // demand_page_allocation(Cr2::read());

    loop {}
}

extern "x86-interrupt" fn timer_interrupt_handler(
    _stack_frame: InterruptStackFrame
) {
    // Context switching 트리거
    unsafe {
        PICS.lock()
            .notify_end_of_interrupt(InterruptIndex::Timer.as_u8());
    }

    // scheduler_tick();
}
```

---

## 🔄 Phase 4: Context Switching (1주)

### 4.1 Context 저장/복원 어셈블리

**파일**: `src/context_switch.asm`

```asm
; src/context_switch.asm
[bits 64]

global switch_context
; rdi = prev_context (ProcessContext*)
; rsi = next_context (ProcessContext*)

switch_context:
    ; 현재 레지스터 상태를 prev_context에 저장
    mov [rdi + 0], rax      ; offset 0
    mov [rdi + 8], rbx      ; offset 8
    mov [rdi + 16], rcx     ; offset 16
    mov [rdi + 24], rdx     ; offset 24
    mov [rdi + 32], rsi     ; offset 32
    mov [rdi + 40], rdi     ; offset 40 (RDI는 이미 스택에 있음)
    mov [rdi + 48], rbp     ; offset 48
    mov [rdi + 56], rsp     ; offset 56

    ; RIP 저장 (반환 주소)
    mov rax, [rsp]
    mov [rdi + 64], rax     ; offset 64: RIP

    ; next_context에서 레지스터 복원
    mov rax, [rsi + 0]      ; RAX
    mov rbx, [rsi + 8]      ; RBX
    mov rcx, [rsi + 16]     ; RCX
    mov rdx, [rsi + 24]     ; RDX
    mov r8, [rsi + 32]      ; RSI
    mov r9, [rsi + 40]      ; RDI
    mov rbp, [rsi + 48]     ; RBP
    mov rsp, [rsi + 56]     ; RSP

    ; 새 페이지 테이블로 전환 (CR3)
    ; rsi + 72 = page_table_base offset
    mov rax, [rsi + 72]
    mov cr3, rax

    ; 새 RIP로 점프
    mov rax, [rsi + 64]
    jmp rax
```

### 4.2 Rust에서 호출

**파일**: `src/scheduler.rs`

```rust
// src/scheduler.rs
extern "C" {
    fn switch_context(prev: *mut ProcessContext, next: *mut ProcessContext);
}

pub struct Scheduler {
    processes: Vec<Process>,
    current_idx: usize,
}

impl Scheduler {
    pub fn schedule(&mut self) {
        let next_idx = (self.current_idx + 1) % self.processes.len();

        unsafe {
            switch_context(
                &mut self.processes[self.current_idx].context as *mut _,
                &mut self.processes[next_idx].context as *mut _,
            );
        }

        self.current_idx = next_idx;
    }
}
```

---

## 📊 Phase 5: 통합 테스트 (1주)

### 5.1 main.rs 통합

**파일**: `src/main.rs`

```rust
#![no_std]
#![no_main]

extern crate alloc;

mod bootloader_info;
mod gdt;
mod interrupts;
mod memory;
mod scheduler;
mod uart;  // QEMU 시리얼 출력

use core::panic::PanicInfo;

#[no_mangle]
pub extern "C" fn kernel_main(boot_info: &'static BootloaderInput) -> ! {
    // 1. 시리얼 초기화 (QEMU에 출력)
    uart::init();
    println!("Freelang OS Kernel booting...");

    // 2. GDT 초기화
    gdt::init_gdt();
    println!("✓ GDT initialized");

    // 3. IDT 초기화
    interrupts::init();
    println!("✓ IDT initialized");

    // 4. 메모리 초기화
    let mut phys_mem = memory::PhysicalMemoryManager::new(0x1000);
    let mut virt_mem = unsafe { memory::VirtualMemoryManager::new(0xFFFF800000000000) };
    println!("✓ Memory initialized");

    // 5. 프로세스 생성
    let mut scheduler = scheduler::Scheduler::new();
    scheduler.create_process(0x400000);  // 커널 프로세스
    scheduler.create_process(0x500000);  // 사용자 프로세스 1
    scheduler.create_process(0x600000);  // 사용자 프로세스 2
    println!("✓ Processes created");

    // 6. 스케줄러 시작
    println!("Starting scheduler...");
    unsafe { asm!("sti"); }  // 인터럽트 활성화

    // 메인 루프
    loop {
        scheduler.schedule();
    }
}

#[panic_handler]
fn panic(info: &PanicInfo) -> ! {
    println!("PANIC: {}", info);
    loop {}
}
```

### 5.2 QEMU 빌드 및 실행

**파일**: `build.sh`

```bash
#!/bin/bash

# 1. x86_64 bare-metal 타겟 설정
rustup target add x86_64-unknown-none

# 2. 어셈블리 컴파일
nasm -f elf64 src/boot.s -o target/boot.o
nasm -f elf64 src/long_mode_init.asm -o target/long_mode.o
nasm -f elf64 src/context_switch.asm -o target/context.o

# 3. Rust 커널 빌드
cargo build --release --target x86_64-unknown-none

# 4. 링킹
x86_64-linux-gnu-ld \
    -n \
    -T linker.ld \
    target/boot.o \
    target/long_mode.o \
    target/context.o \
    target/x86_64-unknown-none/release/kernel \
    -o target/kernel.bin

# 5. ISO 생성 (GRUB)
mkdir -p target/iso/boot/grub
cp target/kernel.bin target/iso/boot/
cat > target/iso/boot/grub/grub.cfg << 'EOF'
set default=0
set timeout=0

menuentry "FreeLang OS" {
    multiboot2 /boot/kernel.bin
    boot
}
EOF

grub-mkrescue -o target/freelang.iso target/iso

# 6. QEMU에서 실행
qemu-system-x86_64 \
    -cdrom target/freelang.iso \
    -serial stdio \
    -m 512 \
    -no-reboot
```

---

## 🎯 Phase 6: 실제 기능 구현 (2-3주)

### 6.1 키보드 입력

```rust
// src/keyboard.rs
pub fn init_keyboard() {
    // PIC 설정 (Programmable Interrupt Controller)
    // 키보드: IRQ1
}
```

### 6.2 디스크 I/O

```rust
// src/disk.rs
pub fn read_sector(lba: u64, buffer: &mut [u8]) {
    // ATA/SATA 드라이버 구현
}
```

### 6.3 프로세스 전환 & 메모리 보호

```rust
// src/process.rs
pub fn create_user_process(binary: &[u8]) -> ProcessId {
    // 독립적 메모리 공간에 로드
    // 사용자 모드 설정
}
```

---

## 📈 Phase 7: 테스트 & 최적화 (1-2주)

### 7.1 단위 테스트

```bash
cargo test --target x86_64-unknown-none
```

### 7.2 통합 테스트 (QEMU)

```rust
#[cfg(test)]
mod tests {
    #[test]
    fn test_memory_allocation() {
        // 메모리 할당 테스트
    }

    #[test]
    fn test_context_switching() {
        // Context switching 테스트
    }

    #[test]
    fn test_page_fault() {
        // PageFault 처리 테스트
    }
}
```

---

## 📅 전체 타임라인

| Phase | 내용 | 기간 | 상태 |
|-------|------|------|------|
| **0** | 현재 상태 분석 | ✓ | 완료 |
| **1** | Multiboot2 부트로더 | 1-2주 | 📍 시작 |
| **2** | 메모리 관리 (실제) | 1주 | ⏳ 대기 |
| **3** | 인터럽트 처리 | 1주 | ⏳ 대기 |
| **4** | Context switching | 1주 | ⏳ 대기 |
| **5** | 통합 & 빌드 | 1주 | ⏳ 대기 |
| **6** | 실제 기능 구현 | 2-3주 | ⏳ 대기 |
| **7** | 테스트 & 최적화 | 1-2주 | ⏳ 대기 |
| **총** | | **8-10주** | |

---

## 🛠️ 필요한 도구

```bash
# 개발 환경
- Rust (stable + nightly)
- rustup target add x86_64-unknown-none
- NASM (어셈블리 어셈블러)
- x86_64-linux-gnu-binutils (gcc 도구)
- QEMU (x86_64 에뮬레이터)
- GRUB (부트로더)

# 설치
sudo apt install \
    nasm \
    binutils-x86-64-linux-gnu \
    qemu-system-x86 \
    grub-pc \
    xorriso
```

---

## 🚀 빠른 시작 (지금 바로)

### 옵션 A: 최소한의 커널 (1주)
1. Multiboot2 부트로더만 구현
2. 간단한 커널 진입점
3. QEMU에서 "Hello World" 출력
4. 타이머 인터럽트 처리

**결과**: 부팅 가능한 OS 이미지 생성

### 옵션 B: 완전한 OS (8-10주)
1. 모든 Phase 1-7 구현
2. 메모리 관리, Context switching, 인터럽트 모두 실제 동작
3. 사용자 프로세스 실행
4. 키보드, 디스크 I/O

**결과**: 실제 OS처럼 동작하는 시스템

---

## ✅ 체크리스트

- [ ] Multiboot2 부트로더 작성
- [ ] x86_64 long mode 진입
- [ ] GDT/IDT 초기화
- [ ] 페이지 테이블 설정
- [ ] 타이머 인터럽트 구현
- [ ] Context switching 구현
- [ ] 메모리 할당자 통합
- [ ] 프로세스 생성/관리
- [ ] QEMU에서 부팅 테스트
- [ ] 사용자 모드 프로세스 실행
- [ ] 키보드/디스크 입출력
- [ ] 성능 최적화

---

## 📚 추천 자료

- **OSDev.org**: https://wiki.osdev.org/Main_Page
- **Intel x86-64 Manual**: https://www.intel.com/content/www/en/en/developer/articles/technical/intel-sdm.html
- **OS in 30 Days**: 교육용 OS 개발 가이드
- **Philipp Oppermann's Blog**: https://os.phil-opp.com/

---

**다음 단계**: Phase 1 (Multiboot2 부트로더) 구현을 시작하겠습니까?
