// 라운드-로빈 스케줄러
// Phase 3: 프로세스 스케줄링 & Context Switching

use crate::context::{ProcessContext, ContextSwitcher};
use alloc::vec::Vec;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// 프로세스 상태
#[derive(Clone, Copy, Debug, PartialEq)]
pub enum ProcessState {
    Ready,       // 실행 대기 중
    Running,     // 현재 실행 중
    Blocked,     // 대기 중 (I/O 등)
    Terminated,  // 종료됨
}

/// 프로세스 (PCB: Process Control Block)
#[repr(C)]
pub struct Process {
    pub pid: u32,
    pub state: ProcessState,
    pub context: ProcessContext,
    pub time_slice_ms: u32,      // 타임 슬라이스 (4ms)
    pub used_time_ms: u32,       // 현재 슬라이스에서 사용한 시간
    pub total_cpu_time_ms: u64,  // 누적 CPU 사용 시간
}

impl Process {
    /// 새 프로세스 생성
    pub fn new(pid: u32, entry_point: u64, stack_top: u64) -> Self {
        Process {
            pid,
            state: ProcessState::Ready,
            context: ProcessContext::new(pid, entry_point, stack_top),
            time_slice_ms: 4,
            used_time_ms: 0,
            total_cpu_time_ms: 0,
        }
    }

    /// 타임 슬라이스 시간 소비
    pub fn consume_time(&mut self, duration_ms: u32) {
        self.used_time_ms += duration_ms;
        self.total_cpu_time_ms += duration_ms as u64;
    }

    /// 타임 슬라이스 초기화
    pub fn reset_time_slice(&mut self) {
        self.used_time_ms = 0;
    }

    /// 타임 슬라이스 남은 시간
    pub fn time_remaining(&self) -> u32 {
        if self.used_time_ms < self.time_slice_ms {
            self.time_slice_ms - self.used_time_ms
        } else {
            0
        }
    }

    /// 상태 변경
    pub fn set_state(&mut self, state: ProcessState) {
        self.state = state;
    }

    /// 실행 가능 여부
    pub fn is_runnable(&self) -> bool {
        self.state == ProcessState::Ready || self.state == ProcessState::Running
    }
}

/// 라운드-로빈 스케줄러
pub struct RoundRobinScheduler {
    pub processes: Vec<Process>,
    pub current_index: Option<usize>,
    pub tick_count: u32,
    pub context_switches: u64,
}

impl RoundRobinScheduler {
    /// 새 스케줄러 생성
    pub fn new() -> Self {
        RoundRobinScheduler {
            processes: Vec::new(),
            current_index: None,
            tick_count: 0,
            context_switches: 0,
        }
    }

    /// 프로세스 추가
    pub fn add_process(&mut self, pid: u32, entry_point: u64, stack_top: u64) -> u32 {
        let process = Process::new(pid, entry_point, stack_top);
        self.processes.push(process);

        if self.current_index.is_none() {
            self.current_index = Some(0);
        }

        crate::println!("✨ Process {} created (entry: 0x{:x})", pid, entry_point);
        pid
    }

    /// 다음 프로세스로 스케줄
    pub fn schedule(&mut self) {
        if self.processes.is_empty() {
            return;
        }

        if let Some(curr_idx) = self.current_index {
            // 현재 프로세스 시간 초기화
            self.processes[curr_idx].reset_time_slice();

            // 다음 프로세스 선택 (라운드-로빈)
            let next_idx = (curr_idx + 1) % self.processes.len();

            // Context switching
            let prev_ctx = &mut self.processes[curr_idx].context;
            let next_ctx = &mut self.processes[next_idx].context;

            ContextSwitcher::switch(prev_ctx, next_ctx);

            // 페이지 테이블 변경 (프로세스 격리)
            if next_ctx.page_table_base != 0 {
                ContextSwitcher::set_page_table(next_ctx.page_table_base);
            }

            self.current_index = Some(next_idx);
            self.context_switches += 1;
        }
    }

    /// 타이머 틱 처리 (4ms)
    pub fn tick(&mut self) {
        self.tick_count += 1;

        if let Some(idx) = self.current_index {
            self.processes[idx].consume_time(4);

            // 타임 슬라이스 만료 시 다음 프로세스로
            if self.processes[idx].time_remaining() == 0 {
                self.schedule();
            }
        }
    }

    /// 현재 실행 중인 프로세스
    pub fn current_process(&self) -> Option<&Process> {
        self.current_index.and_then(|idx| self.processes.get(idx))
    }

    /// 현재 실행 중인 프로세스 (mut)
    pub fn current_process_mut(&mut self) -> Option<&mut Process> {
        let idx = self.current_index?;
        self.processes.get_mut(idx)
    }

    /// 스케줄러 상태 출력
    pub fn print_status(&self) {
        crate::println!("\n📊 Scheduler Status:");
        crate::println!("   Total processes: {}", self.processes.len());
        crate::println!("   Context switches: {}", self.context_switches);
        crate::println!("   Timer ticks: {}", self.tick_count);

        if !self.processes.is_empty() {
            crate::println!("\n   Process List:");
            for (idx, proc) in self.processes.iter().enumerate() {
                let marker = if self.current_index == Some(idx) { "▶️" } else { "⏸️" };
                crate::println!("    {} PID {} - {:?} - CPU: {}ms",
                    marker, proc.pid, proc.state, proc.total_cpu_time_ms);
            }
        }
    }
}

lazy_static! {
    /// 글로벌 스케줄러
    pub static ref SCHEDULER: Mutex<RoundRobinScheduler> = {
        Mutex::new(RoundRobinScheduler::new())
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process_creation() {
        let proc = Process::new(1, 0x400000, 0x7FFF0000);
        assert_eq!(proc.pid, 1);
        assert_eq!(proc.state, ProcessState::Ready);
    }

    #[test]
    fn test_scheduler_add_process() {
        let mut sched = RoundRobinScheduler::new();
        sched.add_process(1, 0x400000, 0x7FFF0000);
        sched.add_process(2, 0x410000, 0x7FFE0000);

        assert_eq!(sched.processes.len(), 2);
    }

    #[test]
    fn test_scheduler_tick() {
        let mut sched = RoundRobinScheduler::new();
        sched.add_process(1, 0x400000, 0x7FFF0000);

        sched.tick();
        assert_eq!(sched.tick_count, 1);
        assert_eq!(sched.processes[0].used_time_ms, 4);
    }
}
