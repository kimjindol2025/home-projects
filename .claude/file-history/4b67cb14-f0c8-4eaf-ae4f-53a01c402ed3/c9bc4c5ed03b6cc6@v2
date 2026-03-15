// 물리 메모리 관리자 (Phase G Week 1 bootloader 통합)

const PAGE_SIZE: usize = 4096;
const MEMORY_SIZE: usize = 512 * 1024 * 1024;  // 512MB
const NUM_PAGES: usize = MEMORY_SIZE / PAGE_SIZE;

/// 물리 메모리 페이지
#[repr(C)]
pub struct PhysicalPage {
    pub page_number: u64,
    pub page_size: u64,
    pub phys_addr: u64,
    pub allocated: bool,
}

/// 물리 메모리 관리자
pub struct PhysicalMemoryManager {
    bitmap: &'static mut [u8],  // 비트맵: 1=할당, 0=free
}

impl PhysicalMemoryManager {
    /// 새로운 물리 메모리 관리자 생성
    pub unsafe fn new(bitmap_addr: usize) -> Self {
        let bitmap_size = (NUM_PAGES + 7) / 8;
        let bitmap = core::slice::from_raw_parts_mut(bitmap_addr as *mut u8, bitmap_size);

        // 커널 메모리 영역 표시 (0 ~ 1MB)
        for i in 0..256 {
            bitmap[i / 8] |= 1 << (i % 8);
        }

        // 부트로더 영역 표시
        for i in 256..512 {
            bitmap[i / 8] |= 1 << (i % 8);
        }

        crate::println!("📊 Physical Memory Manager initialized");
        crate::println!("   Total pages: {} ({} MB)", NUM_PAGES, MEMORY_SIZE / (1024 * 1024));
        crate::println!("   Bitmap size: {} bytes", bitmap_size);

        PhysicalMemoryManager { bitmap }
    }

    /// 물리 페이지 할당
    pub fn allocate_page(&mut self) -> Option<u64> {
        for byte_idx in 0..(NUM_PAGES / 8) {
            for bit in 0..8 {
                if (self.bitmap[byte_idx] & (1 << bit)) == 0 {
                    self.bitmap[byte_idx] |= 1 << bit;
                    let page_num = (byte_idx * 8 + bit) as u64;
                    let phys_addr = page_num * PAGE_SIZE as u64;
                    return Some(phys_addr);
                }
            }
        }
        None
    }

    /// 물리 페이지 할당 해제
    pub fn deallocate_page(&mut self, addr: u64) {
        let page_num = (addr / PAGE_SIZE as u64) as usize;
        let byte_idx = page_num / 8;
        let bit = page_num % 8;
        if byte_idx < self.bitmap.len() {
            self.bitmap[byte_idx] &= !(1 << bit);
        }
    }

    /// 사용 가능한 페이지 수
    pub fn free_pages(&self) -> usize {
        let mut count = 0;
        for byte_idx in 0..(NUM_PAGES / 8) {
            for bit in 0..8 {
                if (self.bitmap[byte_idx] & (1 << bit)) == 0 {
                    count += 1;
                }
            }
        }
        count
    }

    /// Fragmentation 비율
    pub fn fragmentation_ratio(&self) -> f64 {
        let mut free_blocks = 0;
        let mut total_blocks = 0;

        for byte in self.bitmap.iter() {
            for bit in 0..8 {
                total_blocks += 1;
                if (byte & (1 << bit)) == 0 {
                    free_blocks += 1;
                }
            }
        }

        if total_blocks == 0 {
            return 0.0;
        }
        (free_blocks as f64) / (total_blocks as f64)
    }

    /// 메모리 상태 출력
    pub fn print_status(&self) {
        let used = NUM_PAGES - self.free_pages();
        let used_mb = (used * PAGE_SIZE) / (1024 * 1024);
        let total_mb = MEMORY_SIZE / (1024 * 1024);

        crate::println!("\n📊 Physical Memory Status:");
        crate::println!("   Used:  {} MB / {} MB", used_mb, total_mb);
        crate::println!("   Free:  {} pages", self.free_pages());
        crate::println!("   Usage: {:.2}%", (used as f64 / NUM_PAGES as f64) * 100.0);
        crate::println!("   Fragmentation: {:.2}%", self.fragmentation_ratio() * 100.0);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_allocate_page() {
        unsafe {
            let mut mgr = PhysicalMemoryManager::new(0x1000);
            let addr = mgr.allocate_page();
            assert!(addr.is_some());
        }
    }
}
