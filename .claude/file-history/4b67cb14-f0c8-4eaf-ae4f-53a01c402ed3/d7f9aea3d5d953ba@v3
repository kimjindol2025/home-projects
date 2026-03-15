// IDT (Interrupt Descriptor Table) 및 인터럽트 핸들러

use x86_64::structures::idt::{InterruptDescriptorTable, InterruptStackFrame, PageFaultErrorCode};
use x86_64::registers::control::Cr2;
use lazy_static::lazy_static;

lazy_static! {
    /// IDT (Interrupt Descriptor Table)
    static ref IDT: InterruptDescriptorTable = {
        let mut idt = InterruptDescriptorTable::new();

        // 예외 핸들러
        idt.breakpoint.set_handler_fn(breakpoint_handler);
        idt.double_fault.set_handler_fn(double_fault_handler);
        idt.page_fault.set_handler_fn(page_fault_handler);
        idt.general_protection_fault.set_handler_fn(gp_fault_handler);
        idt.divide_error.set_handler_fn(divide_error_handler);

        // 하드웨어 인터럽트 (PIC 설정 필요)
        idt[32].set_handler_fn(timer_interrupt_handler);
        idt[33].set_handler_fn(keyboard_interrupt_handler);

        idt
    };
}

/// IDT 초기화
pub fn init() {
    IDT.load();
    init_pic();

    crate::println!("✓ IDT loaded and initialized");
}

/// PIC (Programmable Interrupt Controller) 초기화
fn init_pic() {
    use x86_64::instructions::port::Port;

    // PIC 초기화 (8259A 칩)
    // ICW1 (Initialization Command Word 1)
    unsafe {
        // master PIC: 0x20
        let mut port = Port::new(0x20);
        port.write(0x11u8);

        // slave PIC: 0xA0
        let mut port = Port::new(0xA0);
        port.write(0x11u8);

        // ICW2: 인터럽트 벡터 주소
        let mut port = Port::new(0x21);
        port.write(0x20u8);  // master: 32-39

        let mut port = Port::new(0xA1);
        port.write(0x28u8);  // slave: 40-47

        // ICW3: master/slave 설정
        let mut port = Port::new(0x21);
        port.write(0x04u8);

        let mut port = Port::new(0xA1);
        port.write(0x02u8);

        // ICW4: 8086 모드
        let mut port = Port::new(0x21);
        port.write(0x01u8);

        let mut port = Port::new(0xA1);
        port.write(0x01u8);

        // OCW1: 인터럽트 마스크
        let mut port = Port::new(0x21);
        port.write(0x00u8);  // 모든 인터럽트 활성화

        let mut port = Port::new(0xA1);
        port.write(0x00u8);
    }

    crate::println!("✓ PIC initialized");
}

// ==================== 예외 핸들러 ====================

extern "x86-interrupt" fn breakpoint_handler(stack_frame: InterruptStackFrame) {
    crate::println!("🔴 EXCEPTION: Breakpoint");
    crate::println!("{:#?}", stack_frame);
}

extern "x86-interrupt" fn double_fault_handler(
    stack_frame: InterruptStackFrame,
    _error_code: u64,
) -> ! {
    crate::println!("🔴 EXCEPTION: Double Fault!");
    crate::println!("{:#?}", stack_frame);

    loop {
        unsafe { asm!("hlt"); }
    }
}

extern "x86-interrupt" fn page_fault_handler(
    stack_frame: InterruptStackFrame,
    error_code: PageFaultErrorCode,
) {
    let virtual_addr = Cr2::read();

    crate::println!("📄 PAGE FAULT");
    crate::println!("   Virtual Address: 0x{:x}", virtual_addr);
    crate::println!("   Error Code: {:?}", error_code);
    crate::println!("   RIP: 0x{:x}", stack_frame.instruction_pointer);

    // Phase 2: Demand Paging 구현
    // TODO: 실제 페이지 할당 및 페이지 테이블 업데이트

    // 현재는 커널 패닉
    crate::println!("\n❌ Demand paging not yet implemented");
    crate::println!("   Phase 2에서 구현 예정");

    loop {
        unsafe { asm!("hlt"); }
    }
}

extern "x86-interrupt" fn gp_fault_handler(
    stack_frame: InterruptStackFrame,
    error_code: u64,
) {
    crate::println!("⚠️  EXCEPTION: General Protection Fault");
    crate::println!("   Error Code: 0x{:x}", error_code);
    crate::println!("{:#?}", stack_frame);

    loop {
        unsafe { asm!("hlt"); }
    }
}

extern "x86-interrupt" fn divide_error_handler(stack_frame: InterruptStackFrame) {
    crate::println!("💥 EXCEPTION: Divide by Zero");
    crate::println!("{:#?}", stack_frame);

    loop {
        unsafe { asm!("hlt"); }
    }
}

// ==================== 하드웨어 인터럽트 ====================

static mut TIMER_TICKS: u64 = 0;

extern "x86-interrupt" fn timer_interrupt_handler(_stack_frame: InterruptStackFrame) {
    unsafe {
        TIMER_TICKS += 1;

        // 250 ticks = 1 second (4ms per tick)
        if TIMER_TICKS % 250 == 0 {
            crate::print!("⏱️ ");
        }
    }

    // EOI (End of Interrupt) 신호
    send_eoi(0x20);
}

extern "x86-interrupt" fn keyboard_interrupt_handler(_stack_frame: InterruptStackFrame) {
    // Phase 4: 키보드 드라이버 호출
    crate::keyboard::handle_keyboard_interrupt();

    send_eoi(0x21);
}

/// EOI (End of Interrupt) 신호 전송
fn send_eoi(irq: u8) {
    use x86_64::instructions::port::Port;

    unsafe {
        if irq >= 8 {
            let mut port = Port::new(0xA0u16);
            port.write(0x20u8);  // slave EOI
        }

        let mut port = Port::new(0x20u16);
        port.write(0x20u8);  // master EOI
    }
}

use core::arch::asm;
