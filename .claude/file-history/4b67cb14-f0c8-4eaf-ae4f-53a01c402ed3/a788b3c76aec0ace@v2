// FreeLang OS Kernel - Phase G: Bare-Metal 구현
// 실제 x86-64 하드웨어에서 부팅되는 OS

#![no_std]
#![no_main]
#![feature(asm_const)]

use core::fmt::Write;
use core::panic::PanicInfo;

mod serial;
mod vga_buffer;
mod memory;
mod gdt;
mod interrupts;

use vga_buffer::WRITER;
use memory::PhysicalMemoryManager;

/// Multiboot2 정보 구조체 (간단 버전)
#[repr(C)]
pub struct MultibootInfo {
    size: u32,
    reserved: u32,
}

/// 커널 메인 함수
#[no_mangle]
pub extern "C" fn kernel_main(_multiboot_info: u64) -> ! {
    // VGA 버퍼 초기화
    writeln!(WRITER.lock(), "╔════════════════════════════════════════════════════╗").unwrap();
    writeln!(WRITER.lock(), "║      FreeLang OS Kernel - Phase G Bare-Metal       ║").unwrap();
    writeln!(WRITER.lock(), "║          실제 x86-64 부팅 및 실행                  ║").unwrap();
    writeln!(WRITER.lock(), "╚════════════════════════════════════════════════════╝").unwrap();

    // 시리얼 포트 초기화
    serial::init();
    println!("[SERIAL] Serial port initialized");

    // GDT 초기화
    gdt::init();
    println!("✓ GDT initialized");

    // IDT 초기화
    interrupts::init();
    println!("✓ IDT initialized");

    // 메모리 관리자 초기화
    println!("🔧 Initializing physical memory manager...");
    let _phys_mem = unsafe { PhysicalMemoryManager::new(0x1000) };
    println!("✓ Physical memory manager initialized");

    // 인터럽트 활성화
    unsafe { asm!("sti"); }
    println!("✓ Interrupts enabled");

    // 메인 루프
    println!("\n=== 커널 부팅 완료 ===");
    println!("타이머 인터럽트: 4ms마다 발생");
    println!("프로세스 관리: 준비 중");

    kernel_loop();
}

/// 커널 메인 루프
fn kernel_loop() -> ! {
    let mut tick_count = 0;

    loop {
        tick_count += 1;

        // 1초마다 메시지 출력 (timer_ticks = 250 = 1000ms)
        if tick_count % 250 == 0 {
            let seconds = tick_count / 250;
            println!("⏱️  Uptime: {}s", seconds);
        }

        // 프로세스 스케줄링 (여기서 context switching 발생)
        // scheduler_tick();

        // CPU 대기
        unsafe { asm!("hlt"); }
    }
}

/// Panic 핸들러
#[panic_handler]
fn panic(info: &PanicInfo) -> ! {
    writeln!(
        WRITER.lock(),
        "🔴 KERNEL PANIC: {}",
        info
    ).unwrap();

    println!("PANIC: {}", info);

    loop {
        unsafe { asm!("hlt"); }
    }
}

// 어셈블리 인라인
use core::arch::asm;
