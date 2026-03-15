// 페이지 테이블 관리 및 가상 메모리 매핑
// Phase 2: Demand Paging 기반 구현

use core::ptr::NonNull;
use x86_64::structures::paging::PageTableFlags;

const PAGE_SIZE: usize = 4096;
const ENTRY_COUNT: usize = 512;

/// 페이지 테이블 항목 (64-bit)
#[repr(transparent)]
#[derive(Clone, Copy)]
pub struct PageTableEntry(u64);

impl PageTableEntry {
    /// 새 항목 생성
    pub fn new(physical_addr: u64, flags: PageTableFlags) -> Self {
        PageTableEntry(physical_addr | flags.bits())
    }

    /// 물리 주소 추출
    pub fn addr(&self) -> u64 {
        self.0 & 0x000FFFFFFFFFF000
    }

    /// Present 플래그 확인
    pub fn is_present(&self) -> bool {
        (self.0 & PageTableFlags::PRESENT.bits()) != 0
    }

    /// 쓰기 가능 확인
    pub fn is_writable(&self) -> bool {
        (self.0 & PageTableFlags::WRITABLE.bits()) != 0
    }

    /// 사용자 접근 가능 확인
    pub fn is_user_accessible(&self) -> bool {
        (self.0 & PageTableFlags::USER_ACCESSIBLE.bits()) != 0
    }

    /// Accessed 플래그 확인 (최근 접근 여부)
    pub fn is_accessed(&self) -> bool {
        (self.0 & (1 << 5)) != 0
    }

    /// Dirty 플래그 확인 (수정 여부)
    pub fn is_dirty(&self) -> bool {
        (self.0 & (1 << 6)) != 0
    }

    /// Present 플래그 설정
    pub fn set_present(&mut self) {
        self.0 |= PageTableFlags::PRESENT.bits();
    }

    /// Present 플래그 해제
    pub fn clear_present(&mut self) {
        self.0 &= !PageTableFlags::PRESENT.bits();
    }
}

/// 페이지 테이블 (512개 항목)
#[repr(align(4096))]
pub struct PageTable {
    pub entries: [PageTableEntry; ENTRY_COUNT],
}

impl PageTable {
    /// 새 페이지 테이블 생성
    pub fn new() -> Self {
        PageTable {
            entries: [PageTableEntry(0); ENTRY_COUNT],
        }
    }

    /// 페이지 맵핑 (가상 → 물리)
    pub fn map_page(&mut self, virtual_addr: u64, physical_addr: u64, flags: PageTableFlags) {
        let vaddr = virtual_addr as usize;
        let index = (vaddr >> 12) & 0x1FF;  // 페이지 오프셋 추출

        if index < ENTRY_COUNT {
            self.entries[index] = PageTableEntry::new(physical_addr, flags);
        }
    }

    /// 페이지 언맵 (맵핑 제거)
    pub fn unmap_page(&mut self, virtual_addr: u64) {
        let vaddr = virtual_addr as usize;
        let index = (vaddr >> 12) & 0x1FF;

        if index < ENTRY_COUNT {
            self.entries[index] = PageTableEntry(0);
        }
    }

    /// 가상 주소를 물리 주소로 변환
    pub fn translate(&self, virtual_addr: u64) -> Option<u64> {
        let vaddr = virtual_addr as usize;
        let index = (vaddr >> 12) & 0x1FF;
        let offset = vaddr & 0xFFF;  // 페이지 내 오프셋

        if index < ENTRY_COUNT {
            let entry = self.entries[index];
            if entry.is_present() {
                return Some(entry.addr() + offset as u64);
            }
        }
        None
    }

    /// 메모리 상태 출력
    pub fn print_status(&self) {
        let mut present = 0;
        let mut writable = 0;
        let mut user = 0;

        for entry in &self.entries {
            if entry.is_present() {
                present += 1;
                if entry.is_writable() {
                    writable += 1;
                }
                if entry.is_user_accessible() {
                    user += 1;
                }
            }
        }

        crate::println!("\n📊 Page Table Status:");
        crate::println!("   Present entries: {}", present);
        crate::println!("   Writable: {}", writable);
        crate::println!("   User-accessible: {}", user);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_page_table_entry() {
        let flags = PageTableFlags::PRESENT | PageTableFlags::WRITABLE;
        let entry = PageTableEntry::new(0x1000, flags);

        assert!(entry.is_present());
        assert!(entry.is_writable());
        assert_eq!(entry.addr(), 0x1000);
    }

    #[test]
    fn test_page_table_mapping() {
        let mut table = PageTable::new();
        let flags = PageTableFlags::PRESENT | PageTableFlags::WRITABLE;

        table.map_page(0x2000, 0x3000, flags);
        assert_eq!(table.translate(0x2000), Some(0x3000));
    }
}
