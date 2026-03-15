// FreeLang OS Kernel - Phase G: Bare-Metal 구현
// 실제 x86-64 하드웨어에서 부팅되는 OS

#![no_std]
#![no_main]
#![feature(asm_const)]
#![feature(allocator_api)]

extern crate alloc;

use core::fmt::Write;
use core::panic::PanicInfo;

mod serial;
mod vga_buffer;
mod memory;
mod gdt;
mod interrupts;
mod paging;
mod demand_paging;
mod allocator;
mod context;
mod scheduler;
mod io;
mod keyboard;
mod disk;
mod fat32;
mod vfs;
mod usermode;
mod syscall;

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

    // Phase 2: Demand Paging 초기화
    println!("\n🔧 Initializing demand paging system...");
    let demand_paging_mgr = demand_paging::DEMAND_PAGING.lock();
    println!("✓ Demand paging manager ready");
    drop(demand_paging_mgr);

    // 힙 할당자 초기화
    println!("🔧 Initializing heap allocator...");
    {
        let alloc = allocator::HEAP_ALLOCATOR.lock();
        println!("✓ Heap allocator initialized");
        println!("  Heap range: 0x200000 ~ 0x20000000 (510 MB)");
        drop(alloc);
    }

    // 인터럽트 활성화
    unsafe { asm!("sti"); }
    println!("✓ Interrupts enabled");

    // 메모리 상태 출력
    println!("\n📊 System Memory Status:");
    {
        let alloc = allocator::HEAP_ALLOCATOR.lock();
        alloc.print_status();
    }

    // Phase 3: 프로세스 스케줄러 초기화
    println!("\n🔧 Initializing process scheduler...");
    {
        let mut sched = scheduler::SCHEDULER.lock();
        // 테스트용 프로세스 생성 (Phase 4에서 실제 프로세스 로드로 대체)
        sched.add_process(1, 0x400000, 0x7FFF0000);
        sched.add_process(2, 0x410000, 0x7FFE0000);
        println!("✓ Process scheduler initialized with 2 test processes");
        drop(sched);
    }

    // Phase 4: I/O 드라이버 초기화
    println!("\n🔧 Initializing I/O drivers...");
    {
        let kb = keyboard::KEYBOARD.lock();
        println!("✓ PS/2 Keyboard driver ready");
        drop(kb);
    }
    {
        let disk = disk::DISK.lock();
        disk.print_status();
        drop(disk);
    }

    // Phase 5: 파일 시스템 초기화
    println!("\n🔧 Initializing file system...");
    {
        let vfs = vfs::VFS.lock();
        vfs.print_status();
        drop(vfs);
    }

    // Phase 6: 사용자 모드 초기화
    println!("\n🔧 Initializing user mode...");
    {
        let user_mode = usermode::USER_MODE.lock();
        user_mode.print_status();
        drop(user_mode);
    }

    // 메인 루프
    println!("\n╔════════════════════════════════════════════════════╗");
    println!("║           🚀 커널 부팅 완료 (Phase 6)              ║");
    println!("╠════════════════════════════════════════════════════╣");
    println!("║ ✓ Multiboot2 부트로더                             ║");
    println!("║ ✓ GDT/IDT 초기화                                  ║");
    println!("║ ✓ 타이머 & 키보드 인터럽트                        ║");
    println!("║ ✓ Demand Paging 시스템                           ║");
    println!("║ ✓ 힙 할당자 (First-Fit, Best-Fit)               ║");
    println!("║ ✓ Context Switching & Round-Robin 스케줄러       ║");
    println!("║ ✓ I/O Drivers (PS/2 Keyboard, ATA Disk)         ║");
    println!("║ ✓ FAT32 & 가상 파일 시스템 (VFS)                ║");
    println!("║ ✓ 사용자 모드 (Ring 3) & 시스템 호출            ║");
    println!("╚════════════════════════════════════════════════════╝");
    println!("\n프로세스 타임 슬라이스: 4ms");
    println!("Context Switching: 타이머 틱마다 활성");
    println!("Keyboard: 인터럽트 기반 입력 처리");
    println!("Disk: LBA 기반 섹터 읽기/쓰기");
    println!("File System: FAT32 및 Inode 기반 메타데이터");
    println!("User Mode: Ring 3 권한 분리 및 시스템 호출\n");

    kernel_loop();
}

/// 커널 메인 루프
fn kernel_loop() -> ! {
    let mut tick_count = 0;

    loop {
        tick_count += 1;

        // 1초마다 상태 출력 (timer_ticks = 250 = 1000ms)
        if tick_count % 250 == 0 {
            let seconds = tick_count / 250;
            println!("\n⏱️  Uptime: {}s", seconds);

            // 메모리 상태 출력
            {
                let alloc = allocator::HEAP_ALLOCATOR.lock();
                let available = alloc.available_memory();
                let frag = alloc.fragmentation_ratio();
                println!("  📊 Heap: {} KB free, {:.1}% fragmentation",
                    available / 1024, frag * 100.0);
            }

            // Demand Paging 상태
            {
                let dp = demand_paging::DEMAND_PAGING.lock();
                dp.print_status();
            }
        }

        // 5초마다 상세 정보 출력
        if tick_count % 1250 == 0 && tick_count > 0 {
            println!("\n📋 Detailed System Status:");
            {
                let alloc = allocator::HEAP_ALLOCATOR.lock();
                alloc.print_status();
            }
            {
                let sched = scheduler::SCHEDULER.lock();
                sched.print_status();
            }
            {
                let kb = keyboard::KEYBOARD.lock();
                kb.print_status();
            }
            {
                let disk = disk::DISK.lock();
                disk.print_status();
            }
            {
                let vfs = vfs::VFS.lock();
                vfs.print_status();
            }
            {
                let user_mode = usermode::USER_MODE.lock();
                user_mode.print_status();
            }
            println!();
        }

        // 프로세스 스케줄링 (타이머 틱마다 호출)
        {
            let mut sched = scheduler::SCHEDULER.lock();
            sched.tick();
        }

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
