// 사용자 모드 (Ring 3) 구현
// Phase 6: 권한 분리 및 사용자 프로세스

use x86_64::structures::gdt::SegmentSelector;
use x86_64::registers::segmentation::Segment;
use spin::Mutex;
use lazy_static::lazy_static;

/// 권한 수준
#[derive(Clone, Copy, Debug, PartialEq)]
pub enum PrivilegeLevel {
    Ring0,  // 커널 모드
    Ring1,  // 시스템 소프트웨어
    Ring2,  // 시스템 소프트웨어
    Ring3,  // 사용자 모드
}

impl PrivilegeLevel {
    pub fn as_u8(&self) -> u8 {
        match self {
            PrivilegeLevel::Ring0 => 0,
            PrivilegeLevel::Ring1 => 1,
            PrivilegeLevel::Ring2 => 2,
            PrivilegeLevel::Ring3 => 3,
        }
    }
}

/// 사용자 프로세스 컨텍스트
#[repr(C)]
#[derive(Clone, Copy, Debug)]
pub struct UserContext {
    pub rip: u64,           // 사용자 코드의 실행 포인터
    pub rsp: u64,           // 사용자 스택 포인터
    pub rbp: u64,           // 사용자 베이스 포인터
    pub rax: u64,
    pub rbx: u64,
    pub rcx: u64,
    pub rdx: u64,
    pub rsi: u64,
    pub rdi: u64,
    pub r8: u64,
    pub r9: u64,
    pub r10: u64,
    pub r11: u64,
    pub r12: u64,
    pub r13: u64,
    pub r14: u64,
    pub r15: u64,
}

impl UserContext {
    /// 새 사용자 컨텍스트 생성
    pub fn new(entry_point: u64, stack_top: u64) -> Self {
        UserContext {
            rip: entry_point,
            rsp: stack_top,
            rbp: stack_top,
            rax: 0,
            rbx: 0,
            rcx: 0,
            rdx: 0,
            rsi: 0,
            rdi: 0,
            r8: 0,
            r9: 0,
            r10: 0,
            r11: 0,
            r12: 0,
            r13: 0,
            r14: 0,
            r15: 0,
        }
    }

    /// 레지스터 값 설정
    pub fn set_register(&mut self, reg: u8, value: u64) {
        match reg {
            0 => self.rax = value,
            1 => self.rbx = value,
            2 => self.rcx = value,
            3 => self.rdx = value,
            4 => self.rsi = value,
            5 => self.rdi = value,
            6 => self.rbp = value,
            7 => self.rsp = value,
            _ => {} // 무시
        }
    }

    /// 시스템 호출 준비
    pub fn prepare_syscall(&mut self, syscall_number: u64) {
        // RAX에 시스템 호출 번호 저장
        self.rax = syscall_number;
    }
}

/// 사용자 모드 관리자
pub struct UserModeManager {
    user_segments: UserSegments,
    active_user_processes: usize,
}

/// 사용자 세그먼트 (GDT에서)
#[derive(Clone, Copy, Debug)]
pub struct UserSegments {
    pub code_selector: SegmentSelector,  // 사용자 코드 세그먼트
    pub data_selector: SegmentSelector,  // 사용자 데이터 세그먼트
}

impl UserModeManager {
    /// 새 사용자 모드 관리자 생성
    pub fn new() -> Self {
        UserModeManager {
            user_segments: UserSegments {
                code_selector: SegmentSelector::new(3, x86_64::structures::gdt::ring::Ring3),
                data_selector: SegmentSelector::new(4, x86_64::structures::gdt::ring::Ring3),
            },
            active_user_processes: 0,
        }
    }

    /// 사용자 모드로 전환 (iretq)
    #[inline]
    pub fn switch_to_user_mode(&self, context: &UserContext) -> ! {
        unsafe {
            // 사용자 모드로 전환하기 위해 iretq 사용
            // 스택에 다음 정보를 쌓음:
            // SS (스택 세그먼트)
            // RSP (스택 포인터)
            // RFLAGS (플래그)
            // CS (코드 세그먼트)
            // RIP (명령어 포인터)

            use core::arch::asm;

            asm!(
                "push {data}",      // SS
                "push {rsp}",       // RSP
                "pushf",            // RFLAGS
                "push {code}",      // CS
                "push {rip}",       // RIP
                "iretq",            // 사용자 모드로 전환
                data = in(reg) self.user_segments.data_selector.0,
                rsp = in(reg) context.rsp,
                code = in(reg) self.user_segments.code_selector.0,
                rip = in(reg) context.rip,
                options(noreturn)
            );
        }
    }

    /// 사용자 프로세스 생성
    pub fn create_user_process(
        &mut self,
        entry_point: u64,
        stack_start: u64,
        stack_size: u64,
    ) -> Result<UserContext, &'static str> {
        if stack_start == 0 || stack_size == 0 {
            return Err("Invalid stack parameters");
        }

        let stack_top = stack_start + stack_size - 1;
        let context = UserContext::new(entry_point, stack_top);

        self.active_user_processes += 1;
        crate::println!("👤 User process created (entry: 0x{:x}, stack: 0x{:x}-0x{:x})",
            entry_point, stack_start, stack_top);

        Ok(context)
    }

    /// 사용자 프로세스 종료
    pub fn exit_user_process(&mut self) {
        if self.active_user_processes > 0 {
            self.active_user_processes -= 1;
        }
        crate::println!("👤 User process exited");
    }

    /// 세그먼트 권한 확인
    pub fn check_privilege(&self, level: PrivilegeLevel) -> Result<(), &'static str> {
        match level {
            PrivilegeLevel::Ring0 => Ok(()),   // 항상 허용
            PrivilegeLevel::Ring3 => Ok(()),   // 사용자는 사용자 모드만 사용
            _ => Err("Invalid privilege level"),
        }
    }

    /// 상태 출력
    pub fn print_status(&self) {
        crate::println!("\n👤 User Mode Manager Status:");
        crate::println!("   User code selector: {:?}", self.user_segments.code_selector);
        crate::println!("   User data selector: {:?}", self.user_segments.data_selector);
        crate::println!("   Active user processes: {}", self.active_user_processes);
    }
}

/// 시스템 호출 번호
pub mod syscalls {
    pub const SYS_EXIT: u64 = 0;
    pub const SYS_WRITE: u64 = 1;
    pub const SYS_READ: u64 = 2;
    pub const SYS_OPEN: u64 = 3;
    pub const SYS_CLOSE: u64 = 4;
    pub const SYS_FORK: u64 = 5;
    pub const SYS_EXEC: u64 = 6;
    pub const SYS_WAIT: u64 = 7;
    pub const SYS_GETPID: u64 = 8;
}

lazy_static! {
    /// 글로벌 사용자 모드 관리자
    pub static ref USER_MODE: Mutex<UserModeManager> = {
        Mutex::new(UserModeManager::new())
    };
}

/// 시스템 호출 핸들러
pub fn handle_syscall(number: u64, arg1: u64, arg2: u64, arg3: u64) -> u64 {
    match number {
        syscalls::SYS_EXIT => {
            crate::println!("📤 [syscall] exit({})", arg1);
            0
        }
        syscalls::SYS_WRITE => {
            crate::println!("📝 [syscall] write(fd={}, addr=0x{:x}, size={})", arg1, arg2, arg3);
            arg3 as u64  // 쓴 바이트 수 반환
        }
        syscalls::SYS_READ => {
            crate::println!("📖 [syscall] read(fd={}, addr=0x{:x}, size={})", arg1, arg2, arg3);
            0  // 읽은 바이트 수 반환
        }
        syscalls::SYS_OPEN => {
            crate::println!("📂 [syscall] open(path=0x{:x}, flags={})", arg1, arg2);
            1  // 파일 디스크립터 반환
        }
        syscalls::SYS_CLOSE => {
            crate::println!("🔒 [syscall] close(fd={})", arg1);
            0
        }
        syscalls::SYS_GETPID => {
            crate::println!("🆔 [syscall] getpid()");
            1  // PID 반환
        }
        _ => {
            crate::println!("❌ [syscall] Unknown syscall: {}", number);
            -1i64 as u64
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_privilege_level() {
        let level = PrivilegeLevel::Ring3;
        assert_eq!(level.as_u8(), 3);
    }

    #[test]
    fn test_user_context_creation() {
        let ctx = UserContext::new(0x400000, 0x7FFF0000);
        assert_eq!(ctx.rip, 0x400000);
        assert_eq!(ctx.rsp, 0x7FFF0000);
        assert_eq!(ctx.rax, 0);
    }

    #[test]
    fn test_user_context_register_set() {
        let mut ctx = UserContext::new(0x400000, 0x7FFF0000);
        ctx.set_register(0, 42);  // RAX = 42
        assert_eq!(ctx.rax, 42);
    }

    #[test]
    fn test_user_mode_manager_creation() {
        let mgr = UserModeManager::new();
        assert_eq!(mgr.active_user_processes, 0);
    }

    #[test]
    fn test_syscall_numbers() {
        assert_eq!(syscalls::SYS_EXIT, 0);
        assert_eq!(syscalls::SYS_WRITE, 1);
        assert_eq!(syscalls::SYS_READ, 2);
        assert_eq!(syscalls::SYS_GETPID, 8);
    }
}
