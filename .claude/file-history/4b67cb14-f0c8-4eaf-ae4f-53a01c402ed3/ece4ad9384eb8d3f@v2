// 시스템 호출 인터페이스
// Phase 6: 사용자 모드 ↔ 커널 모드 전환

use crate::usermode;

/// 시스템 호출 래퍼
pub struct SyscallInterface;

impl SyscallInterface {
    /// 시스템 호출 처리 (사용자 모드에서 호출)
    pub fn call(
        number: u64,
        arg1: u64,
        arg2: u64,
        arg3: u64,
        arg4: u64,
        arg5: u64,
    ) -> i64 {
        // 권한 검증 (현재는 스킵)
        match number {
            usermode::syscalls::SYS_EXIT => {
                Self::sys_exit(arg1 as i32);
                0
            }
            usermode::syscalls::SYS_WRITE => {
                Self::sys_write(arg1 as i32, arg2 as *const u8, arg3 as usize) as i64
            }
            usermode::syscalls::SYS_READ => {
                Self::sys_read(arg1 as i32, arg2 as *mut u8, arg3 as usize) as i64
            }
            usermode::syscalls::SYS_OPEN => {
                Self::sys_open(arg1 as *const u8, arg2 as i32) as i64
            }
            usermode::syscalls::SYS_CLOSE => {
                Self::sys_close(arg1 as i32) as i64
            }
            usermode::syscalls::SYS_GETPID => {
                Self::sys_getpid()
            }
            _ => {
                crate::println!("❌ Unknown syscall: {}", number);
                -1
            }
        }
    }

    /// exit(status) 시스템 호출
    fn sys_exit(status: i32) -> ! {
        crate::println!("👤 Process exit with status: {}", status);
        loop {
            unsafe { core::arch::asm!("hlt"); }
        }
    }

    /// write(fd, buf, count) 시스템 호출
    fn sys_write(fd: i32, buf: *const u8, count: usize) -> usize {
        if fd == 1 {
            // stdout
            unsafe {
                for i in 0..count {
                    let ch = *buf.add(i) as char;
                    crate::print!("{}", ch);
                }
            }
            count
        } else {
            // 다른 fd는 현재 미지원
            0
        }
    }

    /// read(fd, buf, count) 시스템 호출
    fn sys_read(fd: i32, buf: *mut u8, count: usize) -> usize {
        if fd == 0 {
            // stdin: 현재는 미지원
            0
        } else {
            0
        }
    }

    /// open(filename, flags) 시스템 호출
    fn sys_open(filename: *const u8, flags: i32) -> i32 {
        // 파일 시스템을 통해 파일 열기
        // 현재는 시뮬레이션
        unsafe {
            let mut name = [0u8; 256];
            let mut i = 0;
            while *filename.add(i) != 0 && i < 255 {
                name[i] = *filename.add(i);
                i += 1;
            }
            crate::println!("📂 open() called for file at 0x{:x}, flags={}", filename as u64, flags);
        }
        1  // 파일 디스크립터
    }

    /// close(fd) 시스템 호출
    fn sys_close(fd: i32) -> i32 {
        crate::println!("🔒 close(fd={})", fd);
        0
    }

    /// getpid() 시스템 호출
    fn sys_getpid() -> i64 {
        1  // 프로세스 ID (현재는 항상 1)
    }
}

/// SYSCALL 명령어 핸들러 (x86-64 빠른 시스템 호출)
pub fn handle_syscall_fast(rax: u64, rdi: u64, rsi: u64, rdx: u64, r10: u64, r8: u64) -> i64 {
    // x86-64 System V ABI 호출 규약:
    // RAX = syscall number
    // RDI = arg1
    // RSI = arg2
    // RDX = arg3
    // R10 = arg4 (RCX는 RIP 저장용)
    // R8 = arg5
    // R9 = arg6

    SyscallInterface::call(rax, rdi, rsi, rdx, r10, r8)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_syscall_interface() {
        // 실제로는 사용자 모드에서 호출되지만,
        // 테스트에서는 직접 호출 가능
        let _result = unsafe {
            // getpid() 시스템 호출
            let result = SyscallInterface::call(
                usermode::syscalls::SYS_GETPID,
                0, 0, 0, 0, 0
            );
            assert_eq!(result, 1);  // PID = 1
            result
        };
    }
}
