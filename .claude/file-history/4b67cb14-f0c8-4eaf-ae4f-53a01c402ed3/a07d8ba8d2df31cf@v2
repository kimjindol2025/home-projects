// 가상 파일 시스템 (VFS)
// Phase 5: 파일 시스템 추상화 계층

use alloc::vec::Vec;
use alloc::string::String;
use alloc::collections::BTreeMap;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// 파일 시스템 타입
#[derive(Clone, Copy, Debug, PartialEq)]
pub enum FileSystemType {
    FAT32,
    EXT2,
    Ramfs,
}

/// 파일 권한
#[repr(u16)]
#[derive(Clone, Copy, Debug)]
pub struct FileMode {
    pub bits: u16,
}

impl FileMode {
    pub const OWNER_READ: u16 = 0o400;
    pub const OWNER_WRITE: u16 = 0o200;
    pub const OWNER_EXEC: u16 = 0o100;
    pub const GROUP_READ: u16 = 0o40;
    pub const GROUP_WRITE: u16 = 0o20;
    pub const GROUP_EXEC: u16 = 0o10;
    pub const OTHER_READ: u16 = 0o4;
    pub const OTHER_WRITE: u16 = 0o2;
    pub const OTHER_EXEC: u16 = 0o1;
    pub const DIR: u16 = 0o40000;
    pub const FILE: u16 = 0o100000;

    pub const fn new(bits: u16) -> Self {
        FileMode { bits }
    }

    pub fn is_directory(&self) -> bool {
        (self.bits & Self::DIR) != 0
    }

    pub fn is_file(&self) -> bool {
        (self.bits & Self::FILE) != 0
    }
}

/// Inode (파일 메타데이터)
#[derive(Clone, Debug)]
pub struct Inode {
    pub inode_number: u32,
    pub mode: FileMode,
    pub size: u64,
    pub owner_uid: u32,
    pub owner_gid: u32,
    pub created_time: u32,
    pub modified_time: u32,
    pub accessed_time: u32,
    pub link_count: u32,
    pub block_pointers: Vec<u32>,  // 데이터 블록 포인터
}

impl Inode {
    /// 새 Inode 생성
    pub fn new(inode_number: u32, mode: FileMode, size: u64) -> Self {
        Inode {
            inode_number,
            mode,
            size,
            owner_uid: 0,
            owner_gid: 0,
            created_time: 0,
            modified_time: 0,
            accessed_time: 0,
            link_count: 1,
            block_pointers: Vec::new(),
        }
    }

    pub fn is_directory(&self) -> bool {
        self.mode.is_directory()
    }

    pub fn is_file(&self) -> bool {
        self.mode.is_file()
    }
}

/// 디렉토리 항목
#[derive(Clone, Debug)]
pub struct DirEntry {
    pub name: String,
    pub inode_number: u32,
    pub file_type: u8,  // 0=unknown, 1=file, 2=directory
}

/// 가상 파일 시스템
pub struct VirtualFileSystem {
    fs_type: FileSystemType,
    inode_table: BTreeMap<u32, Inode>,
    next_inode: u32,
    mounted_at: String,
    readonly: bool,
}

impl VirtualFileSystem {
    /// 새 VFS 생성
    pub fn new(fs_type: FileSystemType) -> Self {
        VirtualFileSystem {
            fs_type,
            inode_table: BTreeMap::new(),
            next_inode: 1,
            mounted_at: String::from("/"),
            readonly: false,
        }
    }

    /// 루트 디렉토리 초기화
    pub fn init(&mut self) -> Result<(), &'static str> {
        // 루트 디렉토리 Inode 생성
        let root = Inode::new(
            0,
            FileMode::new(FileMode::DIR | FileMode::OWNER_READ | FileMode::OWNER_WRITE),
            4096,
        );

        self.inode_table.insert(0, root);
        self.next_inode = 1;

        crate::println!("📁 VFS initialized ({})", match self.fs_type {
            FileSystemType::FAT32 => "FAT32",
            FileSystemType::EXT2 => "EXT2",
            FileSystemType::Ramfs => "RAMFS",
        });

        Ok(())
    }

    /// 파일 생성
    pub fn create(&mut self, path: &str, mode: FileMode) -> Result<u32, &'static str> {
        if self.readonly {
            return Err("File system is read-only");
        }

        let inode_num = self.next_inode;
        let inode = Inode::new(inode_num, mode, 0);

        self.inode_table.insert(inode_num, inode);
        self.next_inode += 1;

        crate::println!("✏️ File created: {} (inode {})", path, inode_num);
        Ok(inode_num)
    }

    /// 디렉토리 생성
    pub fn mkdir(&mut self, path: &str) -> Result<u32, &'static str> {
        if self.readonly {
            return Err("File system is read-only");
        }

        let inode_num = self.next_inode;
        let inode = Inode::new(
            inode_num,
            FileMode::new(FileMode::DIR | FileMode::OWNER_READ | FileMode::OWNER_WRITE),
            4096,
        );

        self.inode_table.insert(inode_num, inode);
        self.next_inode += 1;

        crate::println!("📁 Directory created: {} (inode {})", path, inode_num);
        Ok(inode_num)
    }

    /// Inode 조회
    pub fn stat(&self, inode_number: u32) -> Result<Inode, &'static str> {
        self.inode_table
            .get(&inode_number)
            .cloned()
            .ok_or("Inode not found")
    }

    /// 파일 삭제
    pub fn unlink(&mut self, inode_number: u32) -> Result<(), &'static str> {
        if self.readonly {
            return Err("File system is read-only");
        }

        self.inode_table.remove(&inode_number);
        Ok(())
    }

    /// 읽기 전용 마운트
    pub fn mount_readonly(&mut self) {
        self.readonly = true;
        crate::println!("🔒 File system mounted as read-only");
    }

    /// 읽기-쓰기 마운트
    pub fn mount_readwrite(&mut self) {
        self.readonly = false;
        crate::println!("🔓 File system mounted as read-write");
    }

    /// 파일 시스템 상태 출력
    pub fn print_status(&self) {
        crate::println!("\n📁 Virtual File System Status:");
        crate::println!("   Type: {:?}", self.fs_type);
        crate::println!("   Mounted at: {}", self.mounted_at);
        crate::println!("   Read-only: {}", self.readonly);
        crate::println!("   Total inodes: {}", self.inode_table.len());
        crate::println!("   Next inode: {}", self.next_inode);

        if !self.inode_table.is_empty() {
            crate::println!("   Inode list:");
            for (num, inode) in self.inode_table.iter() {
                let type_str = if inode.is_directory() { "[DIR]" } else { "[FILE]" };
                crate::println!("    {} #{} - {} bytes", type_str, num, inode.size);
            }
        }
    }
}

lazy_static! {
    /// 글로벌 VFS (FAT32)
    pub static ref VFS: Mutex<VirtualFileSystem> = {
        let mut vfs = VirtualFileSystem::new(FileSystemType::FAT32);
        let _ = vfs.init();
        Mutex::new(vfs)
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_vfs_creation() {
        let vfs = VirtualFileSystem::new(FileSystemType::FAT32);
        assert_eq!(vfs.fs_type, FileSystemType::FAT32);
    }

    #[test]
    fn test_file_mode_flags() {
        let mode = FileMode::new(FileMode::FILE | FileMode::OWNER_READ | FileMode::OWNER_WRITE);
        assert!(mode.is_file());
        assert!(!mode.is_directory());
    }

    #[test]
    fn test_inode_creation() {
        let mode = FileMode::new(FileMode::FILE);
        let inode = Inode::new(1, mode, 1024);
        assert_eq!(inode.size, 1024);
        assert!(inode.is_file());
    }

    #[test]
    fn test_directory_mode() {
        let mode = FileMode::new(FileMode::DIR);
        let inode = Inode::new(1, mode, 4096);
        assert!(inode.is_directory());
        assert!(!inode.is_file());
    }

    #[test]
    fn test_vfs_stat() {
        let mut vfs = VirtualFileSystem::new(FileSystemType::FAT32);
        let _ = vfs.init();

        let stat = vfs.stat(0);
        assert!(stat.is_ok());
        let inode = stat.unwrap();
        assert!(inode.is_directory());
    }
}
