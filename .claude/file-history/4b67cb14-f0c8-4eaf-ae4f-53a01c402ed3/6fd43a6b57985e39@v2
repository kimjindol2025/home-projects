// 힙 할당자 (개선)
// Phase 2: First-Fit + Best-Fit 알고리즘

use alloc::vec::Vec;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// 메모리 블록
#[derive(Clone, Debug)]
pub struct MemoryBlock {
    pub addr: u64,
    pub size: u64,
    pub allocated: bool,
}

/// 힙 할당자
pub struct HeapAllocator {
    blocks: Vec<MemoryBlock>,
    total_size: u64,
    allocated_size: u64,
}

impl HeapAllocator {
    /// 새 할당자 생성
    pub fn new(start: u64, size: u64) -> Self {
        let mut blocks = Vec::new();
        blocks.push(MemoryBlock {
            addr: start,
            size,
            allocated: false,
        });

        HeapAllocator {
            blocks,
            total_size: size,
            allocated_size: 0,
        }
    }

    /// First-Fit 할당 (처음 맞는 블록 사용)
    pub fn allocate_first_fit(&mut self, size: u64) -> Result<u64, &'static str> {
        for block in &mut self.blocks {
            if !block.allocated && block.size >= size {
                let addr = block.addr;

                // 정확한 크기면 그대로 할당
                if block.size == size {
                    block.allocated = true;
                    self.allocated_size += size;
                    return Ok(addr);
                }

                // 남은 공간 분할
                let remaining = block.size - size;
                block.size = size;
                block.allocated = true;
                self.allocated_size += size;

                self.blocks.push(MemoryBlock {
                    addr: addr + size,
                    size: remaining,
                    allocated: false,
                });

                return Ok(addr);
            }
        }

        Err("Out of memory: First-Fit failed")
    }

    /// Best-Fit 할당 (가장 적합한 블록 사용)
    pub fn allocate_best_fit(&mut self, size: u64) -> Result<u64, &'static str> {
        let mut best_idx = None;
        let mut best_waste = u64::MAX;

        // 가장 적합한 블록 찾기
        for (idx, block) in self.blocks.iter().enumerate() {
            if !block.allocated && block.size >= size {
                let waste = block.size - size;
                if waste < best_waste {
                    best_waste = waste;
                    best_idx = Some(idx);
                    if waste == 0 {
                        break;  // 완벽한 맞춤 찾음
                    }
                }
            }
        }

        if let Some(idx) = best_idx {
            let block = &mut self.blocks[idx];
            let addr = block.addr;

            if block.size == size {
                block.allocated = true;
            } else {
                let remaining = block.size - size;
                block.size = size;
                block.allocated = true;

                self.blocks.push(MemoryBlock {
                    addr: addr + size,
                    size: remaining,
                    allocated: false,
                });
            }

            self.allocated_size += size;
            Ok(addr)
        } else {
            Err("Out of memory: Best-Fit failed")
        }
    }

    /// 메모리 해제
    pub fn deallocate(&mut self, addr: u64) -> Result<(), &'static str> {
        for block in &mut self.blocks {
            if block.addr == addr && block.allocated {
                block.allocated = false;
                self.allocated_size = self.allocated_size.saturating_sub(block.size);

                // 인접한 free 블록 병합 (Coalescing)
                self.coalesce_blocks();
                return Ok(());
            }
        }

        Err("Block not found or already free")
    }

    /// 인접한 free 블록 병합
    fn coalesce_blocks(&mut self) {
        let mut i = 0;
        while i < self.blocks.len() - 1 {
            if !self.blocks[i].allocated && !self.blocks[i + 1].allocated {
                let next_addr = self.blocks[i + 1].addr;
                let next_size = self.blocks[i + 1].size;

                if self.blocks[i].addr + self.blocks[i].size == next_addr {
                    self.blocks[i].size += next_size;
                    self.blocks.remove(i + 1);
                    continue;
                }
            }
            i += 1;
        }
    }

    /// Fragmentation 비율
    pub fn fragmentation_ratio(&self) -> f64 {
        let free_blocks = self.blocks.iter().filter(|b| !b.allocated).count();
        if self.blocks.is_empty() {
            return 0.0;
        }
        (free_blocks as f64) / (self.blocks.len() as f64)
    }

    /// 사용 가능한 메모리
    pub fn available_memory(&self) -> u64 {
        self.blocks
            .iter()
            .filter(|b| !b.allocated)
            .map(|b| b.size)
            .sum()
    }

    /// 메모리 상태 출력
    pub fn print_status(&self) {
        let used_percent = (self.allocated_size as f64 / self.total_size as f64) * 100.0;

        crate::println!("\n📊 Heap Allocator Status:");
        crate::println!("   Total: {} KB", self.total_size / 1024);
        crate::println!("   Used: {} KB ({:.2}%)",
            self.allocated_size / 1024, used_percent);
        crate::println!("   Free: {} KB", self.available_memory() / 1024);
        crate::println!("   Blocks: {}", self.blocks.len());
        crate::println!("   Fragmentation: {:.2}%", self.fragmentation_ratio() * 100.0);
    }
}

lazy_static! {
    /// 글로벌 힙 할당자
    pub static ref HEAP_ALLOCATOR: Mutex<HeapAllocator> = {
        // 0x200000 ~ 0x20000000 (510MB)
        Mutex::new(HeapAllocator::new(0x200000, 510 * 1024 * 1024))
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_first_fit_allocation() {
        let mut alloc = HeapAllocator::new(0x1000, 0x10000);
        let addr = alloc.allocate_first_fit(0x1000).unwrap();
        assert_eq!(addr, 0x1000);
    }

    #[test]
    fn test_best_fit_allocation() {
        let mut alloc = HeapAllocator::new(0x1000, 0x10000);
        alloc.allocate_first_fit(0x1000).unwrap();
        alloc.allocate_first_fit(0x2000).unwrap();

        let addr = alloc.allocate_best_fit(0x1000).unwrap();
        assert_eq!(addr, 0x3000);
    }

    #[test]
    fn test_deallocation() {
        let mut alloc = HeapAllocator::new(0x1000, 0x10000);
        let addr = alloc.allocate_first_fit(0x1000).unwrap();
        assert!(alloc.deallocate(addr).is_ok());
    }

    #[test]
    fn test_fragmentation() {
        let mut alloc = HeapAllocator::new(0x1000, 0x10000);
        alloc.allocate_first_fit(0x1000).unwrap();
        alloc.allocate_first_fit(0x1000).unwrap();
        alloc.allocate_first_fit(0x1000).unwrap();

        assert!(alloc.fragmentation_ratio() > 0.0);
    }
}
