// FAT32 파일 시스템 구현
// Phase 5: 파일 및 디렉토리 관리

use alloc::vec::Vec;
use alloc::string::String;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// FAT32 부트 섹터 (512바이트)
#[repr(C)]
pub struct BootSector {
    pub jmp_boot: [u8; 3],           // 0x00: 부팅 코드 점프
    pub oem_name: [u8; 8],           // 0x03: OEM 이름
    pub bytes_per_sector: u16,       // 0x0B: 섹터당 바이트 (보통 512)
    pub sectors_per_cluster: u8,     // 0x0D: 클러스터당 섹터
    pub reserved_sectors: u16,       // 0x0E: 예약 섹터 수
    pub fat_count: u8,               // 0x10: FAT 테이블 개수 (보통 2)
    pub root_entries: u16,           // 0x11: 루트 디렉토리 항목 (FAT32는 0)
    pub total_sectors_16: u16,       // 0x13: 전체 섹터 (FAT32는 0)
    pub media_descriptor: u8,        // 0x15: 미디어 유형
    pub sectors_per_fat_16: u16,     // 0x16: FAT당 섹터 (FAT32는 0)
    pub sectors_per_track: u16,      // 0x18: 트랙당 섹터
    pub heads: u16,                  // 0x1A: 헤드 수
    pub hidden_sectors: u32,         // 0x1C: 숨겨진 섹터
    pub total_sectors_32: u32,       // 0x20: 전체 섹터 (32비트)

    // FAT32 확장 영역
    pub sectors_per_fat_32: u32,     // 0x24: FAT당 섹터 (32비트)
    pub ext_flags: u16,              // 0x28: 확장 플래그
    pub fs_version: u16,             // 0x2A: 파일 시스템 버전
    pub root_cluster: u32,           // 0x2C: 루트 디렉토리 클러스터
    pub fsinfo_sector: u16,          // 0x30: FSINFO 섹터
    pub backup_boot_sector: u16,     // 0x32: 백업 부트 섹터
    pub reserved: [u8; 12],          // 0x34: 예약
    pub drive_number: u8,            // 0x40: BIOS 드라이브 번호
    pub reserved2: u8,               // 0x41: 예약
    pub boot_signature: u8,          // 0x42: 부트 시그니처 (0x29)
    pub volume_serial: u32,          // 0x43: 볼륨 일련번호
    pub volume_label: [u8; 11],      // 0x47: 볼륨 레이블
    pub fs_type: [u8; 8],            // 0x52: 파일 시스템 타입 ("FAT32   ")
}

/// FAT32 디렉토리 항목 (32바이트)
#[repr(C)]
#[derive(Clone, Copy, Debug)]
pub struct DirectoryEntry {
    pub name: [u8; 8],               // 0x00: 파일 이름
    pub extension: [u8; 3],          // 0x08: 파일 확장자
    pub attributes: u8,              // 0x0B: 파일 속성
    pub reserved: u8,                // 0x0C: 예약
    pub create_time_tenths: u8,      // 0x0D: 생성 시간 (밀리초)
    pub create_time: u16,            // 0x0E: 생성 시간 (시:분:초)
    pub create_date: u16,            // 0x10: 생성 날짜
    pub access_date: u16,            // 0x12: 액세스 날짜
    pub cluster_high: u16,           // 0x14: 시작 클러스터 (상위 16비트)
    pub write_time: u16,             // 0x16: 쓰기 시간
    pub write_date: u16,             // 0x18: 쓰기 날짜
    pub cluster_low: u16,            // 0x1A: 시작 클러스터 (하위 16비트)
    pub file_size: u32,              // 0x1C: 파일 크기
}

/// 파일 속성 플래그
pub mod attributes {
    pub const READ_ONLY: u8 = 0x01;
    pub const HIDDEN: u8 = 0x02;
    pub const SYSTEM: u8 = 0x04;
    pub const VOLUME_ID: u8 = 0x08;
    pub const DIRECTORY: u8 = 0x10;
    pub const ARCHIVE: u8 = 0x20;
    pub const LONG_NAME: u8 = 0x0F;  // LFN (Long File Name)
}

/// 파일 핸들
#[derive(Clone, Debug)]
pub struct FileHandle {
    pub name: String,
    pub start_cluster: u32,
    pub file_size: u32,
    pub current_position: u32,
    pub attributes: u8,
}

impl FileHandle {
    /// 새 파일 핸들 생성
    pub fn new(name: String, cluster: u32, size: u32, attr: u8) -> Self {
        FileHandle {
            name,
            start_cluster: cluster,
            file_size: size,
            current_position: 0,
            attributes: attr,
        }
    }

    /// 파일이 디렉토리인지 확인
    pub fn is_directory(&self) -> bool {
        (self.attributes & attributes::DIRECTORY) != 0
    }

    /// 파일이 읽기 전용인지 확인
    pub fn is_read_only(&self) -> bool {
        (self.attributes & attributes::READ_ONLY) != 0
    }

    /// 파일 포인터 리셋
    pub fn seek_to_start(&mut self) {
        self.current_position = 0;
    }

    /// 파일 끝인지 확인
    pub fn is_eof(&self) -> bool {
        self.current_position >= self.file_size
    }
}

/// FAT32 파일 시스템 관리자
pub struct FileSystem {
    boot_sector: Option<BootSector>,
    open_files: Vec<FileHandle>,
    fat_table: Vec<u32>,  // FAT 테이블 (클러스터 링크)
}

impl FileSystem {
    /// 새 파일 시스템 생성
    pub fn new() -> Self {
        FileSystem {
            boot_sector: None,
            open_files: Vec::new(),
            fat_table: Vec::new(),
        }
    }

    /// 파일 시스템 초기화 (디스크에서 읽기)
    pub fn init(&mut self) -> Result<(), &'static str> {
        // TODO: 실제로는 디스크에서 부트 섹터 읽기
        // 현재는 시뮬레이션
        crate::println!("🗂️ FAT32 File System initialized");
        Ok(())
    }

    /// 파일 열기 (읽기)
    pub fn open(&mut self, filename: &str) -> Result<FileHandle, &'static str> {
        // TODO: 실제로는 디렉토리에서 파일 검색
        // 현재는 시뮬레이션

        // 파일을 찾았다고 가정
        let handle = FileHandle::new(
            String::from(filename),
            1,  // 시작 클러스터
            512, // 파일 크기
            0,   // 속성
        );

        self.open_files.push(handle.clone());
        Ok(handle)
    }

    /// 파일 닫기
    pub fn close(&mut self, filename: &str) -> Result<(), &'static str> {
        self.open_files.retain(|f| f.name.as_str() != filename);
        Ok(())
    }

    /// 파일 읽기
    pub fn read(&mut self, handle: &mut FileHandle, size: usize) -> Result<Vec<u8>, &'static str> {
        if handle.current_position >= handle.file_size {
            return Ok(Vec::new());  // 파일 끝
        }

        let remaining = handle.file_size - handle.current_position;
        let read_size = if size as u32 > remaining {
            remaining as usize
        } else {
            size
        };

        // TODO: 실제로는 디스크에서 데이터 읽기
        let mut buffer = Vec::with_capacity(read_size);
        buffer.resize(read_size, 0);

        handle.current_position += read_size as u32;
        Ok(buffer)
    }

    /// 파일 쓰기
    pub fn write(&mut self, handle: &mut FileHandle, data: &[u8]) -> Result<usize, &'static str> {
        if handle.is_read_only() {
            return Err("File is read-only");
        }

        // TODO: 실제로는 디스크에 데이터 쓰기
        let written = data.len();
        handle.current_position += written as u32;

        if handle.current_position > handle.file_size {
            handle.file_size = handle.current_position;
        }

        Ok(written)
    }

    /// 디렉토리 나열
    pub fn list_directory(&self, dirname: &str) -> Result<Vec<FileHandle>, &'static str> {
        // TODO: 실제로는 디렉토리 항목 읽기
        let mut entries = Vec::new();

        // 시뮬레이션: 기본 항목
        entries.push(FileHandle::new(String::from("file1.txt"), 1, 512, 0));
        entries.push(FileHandle::new(String::from("file2.bin"), 2, 1024, 0));

        Ok(entries)
    }

    /// 파일 삭제
    pub fn delete(&mut self, filename: &str) -> Result<(), &'static str> {
        // 열려있는 파일이면 닫기
        self.open_files.retain(|f| f.name.as_str() != filename);

        // TODO: 실제로는 FAT 테이블에서 클러스터 해제
        crate::println!("🗑️ File deleted: {}", filename);
        Ok(())
    }

    /// 파일 이름 변경
    pub fn rename(&mut self, old_name: &str, new_name: &str) -> Result<(), &'static str> {
        for file in &mut self.open_files {
            if file.name.as_str() == old_name {
                file.name = String::from(new_name);
                return Ok(());
            }
        }

        // TODO: 실제로는 디렉토리 항목 업데이트
        crate::println!("📝 File renamed: {} → {}", old_name, new_name);
        Ok(())
    }

    /// 열려있는 파일 목록
    pub fn print_status(&self) {
        crate::println!("\n🗂️ File System Status:");
        crate::println!("   Open files: {}", self.open_files.len());

        if !self.open_files.is_empty() {
            crate::println!("   Files:");
            for file in &self.open_files {
                let type_str = if file.is_directory() { "[DIR]" } else { "[FILE]" };
                crate::println!("    {} {} - {} bytes @ cluster {}",
                    type_str, file.name, file.file_size, file.start_cluster);
            }
        }

        crate::println!("   FAT table entries: {}", self.fat_table.len());
    }
}

lazy_static! {
    /// 글로벌 파일 시스템
    pub static ref FILESYSTEM: Mutex<FileSystem> = {
        let mut fs = FileSystem::new();
        let _ = fs.init();
        Mutex::new(fs)
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_filesystem_creation() {
        let fs = FileSystem::new();
        assert_eq!(fs.open_files.len(), 0);
    }

    #[test]
    fn test_file_handle_creation() {
        let handle = FileHandle::new(String::from("test.txt"), 1, 512, 0);
        assert_eq!(handle.file_size, 512);
        assert!(!handle.is_directory());
        assert!(!handle.is_eof());
    }

    #[test]
    fn test_file_attributes() {
        let mut handle = FileHandle::new(String::from("test.txt"), 1, 512, 0);
        assert!(!handle.is_directory());

        handle.attributes = attributes::DIRECTORY;
        assert!(handle.is_directory());

        handle.attributes = attributes::READ_ONLY;
        assert!(handle.is_read_only());
    }

    #[test]
    fn test_file_seek() {
        let mut handle = FileHandle::new(String::from("test.txt"), 1, 512, 0);
        handle.current_position = 256;
        handle.seek_to_start();
        assert_eq!(handle.current_position, 0);
    }

    #[test]
    fn test_boot_sector_size() {
        assert_eq!(core::mem::size_of::<BootSector>(), 512);
    }

    #[test]
    fn test_directory_entry_size() {
        assert_eq!(core::mem::size_of::<DirectoryEntry>(), 32);
    }

    #[test]
    fn test_file_handle_eof() {
        let mut handle = FileHandle::new(String::from("test.txt"), 1, 512, 0);
        handle.current_position = 512;
        assert!(handle.is_eof());
    }
}
