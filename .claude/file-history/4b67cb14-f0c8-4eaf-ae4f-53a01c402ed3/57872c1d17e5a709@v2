// Demand Paging 구현
// Phase 2: PageFault 시 자동으로 물리 페이지 할당

use crate::memory::PhysicalMemoryManager;
use crate::paging::PageTable;
use x86_64::structures::paging::PageTableFlags;
use spin::Mutex;
use lazy_static::lazy_static;

/// Demand Paging 매니저
pub struct DemandPagingManager {
    /// 현재 프로세스의 페이지 테이블
    page_table: Option<&'static mut PageTable>,

    /// 물리 메모리 관리자 참조
    /// (실제로는 Arc<Mutex<>> 사용)

    /// 할당된 페이지 수
    allocated_pages: usize,

    /// PageFault 발생 횟수
    fault_count: usize,
}

lazy_static! {
    /// 글로벌 Demand Paging 매니저
    pub static ref DEMAND_PAGING: Mutex<DemandPagingManager> = {
        Mutex::new(DemandPagingManager {
            page_table: None,
            allocated_pages: 0,
            fault_count: 0,
        })
    };
}

impl DemandPagingManager {
    /// 새 Demand Paging 매니저 생성
    pub fn new() -> Self {
        DemandPagingManager {
            page_table: None,
            allocated_pages: 0,
            fault_count: 0,
        }
    }

    /// PageFault 처리 (Demand Paging)
    ///
    /// 동작:
    /// 1. 가상 주소의 유효성 확인
    /// 2. 물리 페이지 할당
    /// 3. 페이지 테이블 업데이트
    /// 4. 명령어 재실행
    pub fn handle_page_fault(&mut self, virtual_addr: u64) -> Result<(), &'static str> {
        // 가상 주소 정렬 (페이지 경계)
        let aligned_vaddr = virtual_addr & !0xFFF;

        crate::println!("  [Demand Paging] Allocating page for 0x{:x}", aligned_vaddr);

        // TODO: 실제 물리 메모리 할당
        // 현재는 시뮬레이션만 수행
        self.allocated_pages += 1;
        self.fault_count += 1;

        crate::println!("  ✓ Page allocated (total: {})", self.allocated_pages);
        Ok(())
    }

    /// 페이지 해제 (메모리 부족 시)
    ///
    /// 알고리즘: LRU (Least Recently Used)
    /// 1. Accessed 플래그가 가장 오래된 페이지 찾기
    /// 2. Dirty 플래그 확인 (디스크 쓰기 필요 여부)
    /// 3. 페이지 테이블에서 제거
    pub fn evict_page(&mut self) -> Result<(), &'static str> {
        if self.allocated_pages == 0 {
            return Err("No pages to evict");
        }

        // TODO: LRU 페이지 선택 및 제거
        self.allocated_pages -= 1;

        crate::println!("  [Paging] Evicted 1 page (LRU)");
        Ok(())
    }

    /// 메모리 상태 출력
    pub fn print_status(&self) {
        crate::println!("\n📊 Demand Paging Status:");
        crate::println!("   Allocated Pages: {}", self.allocated_pages);
        crate::println!("   Page Faults: {}", self.fault_count);
        crate::println!("   Fault Rate: {:.2} faults/page",
            if self.allocated_pages > 0 {
                self.fault_count as f64 / self.allocated_pages as f64
            } else {
                0.0
            }
        );
    }
}

/// PageFault 핸들러 진입점
pub fn handle_page_fault(virtual_addr: u64) {
    let mut manager = DEMAND_PAGING.lock();

    match manager.handle_page_fault(virtual_addr) {
        Ok(_) => {
            crate::println!("✓ Page fault resolved");
        }
        Err(e) => {
            crate::println!("❌ Page fault unrecoverable: {}", e);
            loop {
                unsafe { asm!("hlt"); }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_demand_paging_allocation() {
        let mut mgr = DemandPagingManager::new();
        let result = mgr.handle_page_fault(0x2000);

        assert!(result.is_ok());
        assert_eq!(mgr.allocated_pages, 1);
    }

    #[test]
    fn test_page_fault_count() {
        let mut mgr = DemandPagingManager::new();
        mgr.handle_page_fault(0x2000).unwrap();
        mgr.handle_page_fault(0x3000).unwrap();
        mgr.handle_page_fault(0x4000).unwrap();

        assert_eq!(mgr.fault_count, 3);
        assert_eq!(mgr.allocated_pages, 3);
    }
}

use core::arch::asm;
