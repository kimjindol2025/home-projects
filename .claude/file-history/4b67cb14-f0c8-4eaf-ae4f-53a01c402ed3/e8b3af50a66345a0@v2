// Context Switching & 프로세스 관리
// Phase 3: CPU 레지스터 상태 저장/복원

use core::arch::asm;

/// x86-64 일반 목적 레지스터
#[repr(C)]
#[derive(Clone, Copy, Debug)]
pub struct RegisterState {
    pub rax: u64,
    pub rbx: u64,
    pub rcx: u64,
    pub rdx: u64,
    pub rsi: u64,
    pub rdi: u64,
    pub rbp: u64,
    pub rsp: u64,
    pub r8: u64,
    pub r9: u64,
    pub r10: u64,
    pub r11: u64,
    pub r12: u64,
    pub r13: u64,
    pub r14: u64,
    pub r15: u64,
}

impl RegisterState {
    /// 새 레지스터 상태 생성
    pub fn new() -> Self {
        RegisterState {
            rax: 0, rbx: 0, rcx: 0, rdx: 0,
            rsi: 0, rdi: 0, rbp: 0, rsp: 0,
            r8: 0, r9: 0, r10: 0, r11: 0,
            r12: 0, r13: 0, r14: 0, r15: 0,
        }
    }

    /// 스택 포인터 설정 (프로세스 시작)
    pub fn set_stack_pointer(mut self, rsp: u64) -> Self {
        self.rsp = rsp;
        self
    }

    /// 프로그램 카운터 설정 (진입점)
    pub fn set_program_counter(mut self, rip: u64) -> Self {
        self.rax = rip;  // RIP는 별도 저장 (아래 ProcessContext 참고)
        self
    }
}

/// 프로세스 컨텍스트 (CPU 상태 + 메타데이터)
#[repr(C)]
#[derive(Clone, Copy)]
pub struct ProcessContext {
    pub registers: RegisterState,
    pub instruction_pointer: u64,  // RIP
    pub flags: u64,                // RFLAGS
    pub page_table_base: u64,      // CR3 (페이지 테이블)
    pub context_id: u32,
    pub switch_count: u32,
}

impl ProcessContext {
    /// 새 프로세스 컨텍스트 생성
    pub fn new(context_id: u32, entry_point: u64, stack_top: u64) -> Self {
        ProcessContext {
            registers: RegisterState::new().set_stack_pointer(stack_top),
            instruction_pointer: entry_point,
            flags: 0x202,  // IF 플래그 (인터럽트 활성화)
            page_table_base: 0,
            context_id,
            switch_count: 0,
        }
    }

    /// 페이지 테이블 설정 (프로세스 격리)
    pub fn set_page_table(mut self, base: u64) -> Self {
        self.page_table_base = base;
        self
    }

    /// 컨텍스트 전환 횟수 증가
    pub fn record_switch(&mut self) {
        self.switch_count += 1;
    }
}

/// Context Switcher (CPU 상태 전환)
pub struct ContextSwitcher;

impl ContextSwitcher {
    /// Context switching 수행 (어셈블리 호출)
    ///
    /// 동작:
    /// 1. 이전 프로세스 레지스터 저장
    /// 2. 새 프로세스 레지스터 복원
    /// 3. CR3 (페이지 테이블) 변경
    /// 4. RIP로 점프
    pub fn switch(prev_ctx: &mut ProcessContext, next_ctx: &mut ProcessContext) {
        // 현재 프로세스의 상태 저장
        prev_ctx.record_switch();

        // 다음 프로세스 선택
        next_ctx.record_switch();

        // TODO: 실제 어셈블리 context switching
        // 현재는 시뮬레이션
        crate::println!("  🔄 Context switch: ctx{} → ctx{}",
            prev_ctx.context_id, next_ctx.context_id);
    }

    /// CR3 (페이지 테이블 레지스터) 설정
    pub fn set_page_table(page_table_base: u64) {
        unsafe {
            asm!("mov cr3, {}", in(reg) page_table_base);
        }
    }

    /// 현재 CR3 읽기
    pub fn read_page_table() -> u64 {
        let cr3: u64;
        unsafe {
            asm!("mov {}, cr3", out(reg) cr3);
        }
        cr3
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_register_state() {
        let regs = RegisterState::new()
            .set_stack_pointer(0x7FFF0000);
        assert_eq!(regs.rsp, 0x7FFF0000);
    }

    #[test]
    fn test_process_context() {
        let ctx = ProcessContext::new(1, 0x400000, 0x7FFF0000);
        assert_eq!(ctx.context_id, 1);
        assert_eq!(ctx.instruction_pointer, 0x400000);
        assert_eq!(ctx.registers.rsp, 0x7FFF0000);
    }

    #[test]
    fn test_context_switch_count() {
        let mut ctx = ProcessContext::new(1, 0x400000, 0x7FFF0000);
        assert_eq!(ctx.switch_count, 0);
        ctx.record_switch();
        assert_eq!(ctx.switch_count, 1);
    }
}
