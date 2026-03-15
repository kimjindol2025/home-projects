// ATA/SATA 디스크 드라이버
// Phase 4: 디스크 읽기/쓰기 기본 구현

use crate::io::IOPort;
use alloc::vec::Vec;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// ATA 포트 주소
const ATA_PRIMARY_DATA: u16 = 0x1F0;
const ATA_PRIMARY_ERROR: u16 = 0x1F1;
const ATA_PRIMARY_SECTOR_COUNT: u16 = 0x1F2;
const ATA_PRIMARY_SECTOR_NUMBER: u16 = 0x1F3;
const ATA_PRIMARY_CYLINDER_LOW: u16 = 0x1F4;
const ATA_PRIMARY_CYLINDER_HIGH: u16 = 0x1F5;
const ATA_PRIMARY_DRIVE_HEAD: u16 = 0x1F6;
const ATA_PRIMARY_STATUS: u16 = 0x1F7;
const ATA_PRIMARY_COMMAND: u16 = 0x1F7;

/// ATA 상태 플래그
const ATA_STATUS_BUSY: u8 = 0x80;
const ATA_STATUS_READY: u8 = 0x40;
const ATA_STATUS_ERROR: u8 = 0x01;

/// ATA 명령어
const ATA_CMD_READ_PIO: u8 = 0x20;
const ATA_CMD_WRITE_PIO: u8 = 0x30;
const ATA_CMD_IDENTIFY: u8 = 0xEC;

/// 디스크 섹터 크기 (바이트)
pub const DISK_SECTOR_SIZE: usize = 512;

/// ATA 디스크 드라이버
pub struct DiskDriver {
    data_port: IOPort,
    error_port: IOPort,
    sector_count_port: IOPort,
    sector_number_port: IOPort,
    cylinder_low_port: IOPort,
    cylinder_high_port: IOPort,
    drive_head_port: IOPort,
    status_port: IOPort,
    command_port: IOPort,

    initialized: bool,
    total_sectors: u32,
    sectors_read: u64,
    sectors_written: u64,
}

impl DiskDriver {
    /// 새 ATA 디스크 드라이버 생성
    pub fn new() -> Self {
        DiskDriver {
            data_port: IOPort::new(ATA_PRIMARY_DATA),
            error_port: IOPort::new(ATA_PRIMARY_ERROR),
            sector_count_port: IOPort::new(ATA_PRIMARY_SECTOR_COUNT),
            sector_number_port: IOPort::new(ATA_PRIMARY_SECTOR_NUMBER),
            cylinder_low_port: IOPort::new(ATA_PRIMARY_CYLINDER_LOW),
            cylinder_high_port: IOPort::new(ATA_PRIMARY_CYLINDER_HIGH),
            drive_head_port: IOPort::new(ATA_PRIMARY_DRIVE_HEAD),
            status_port: IOPort::new(ATA_PRIMARY_STATUS),
            command_port: IOPort::new(ATA_PRIMARY_COMMAND),

            initialized: false,
            total_sectors: 0,
            sectors_read: 0,
            sectors_written: 0,
        }
    }

    /// 디스크 초기화 (IDENTIFY 명령어 실행)
    pub fn init(&mut self) -> Result<(), &'static str> {
        // IDENTIFY 명령어 전송
        self.select_drive(0)?;
        self.send_command(ATA_CMD_IDENTIFY)?;

        // 응답 대기
        self.wait_ready()?;

        // 현재는 시뮬레이션: 기본값 설정
        self.total_sectors = 1024 * 1024; // 512MB (기본값)
        self.initialized = true;

        crate::println!("💾 ATA Disk initialized ({}MB)", self.total_sectors * DISK_SECTOR_SIZE as u32 / (1024 * 1024));

        Ok(())
    }

    /// 드라이브 선택
    fn select_drive(&self, drive: u8) -> Result<(), &'static str> {
        // 0xA0 = Master drive, 0xB0 = Slave drive
        let drive_byte = if drive == 0 { 0xA0 } else { 0xB0 };
        self.drive_head_port.write_u8(drive_byte);
        Ok(())
    }

    /// 명령어 전송
    fn send_command(&self, cmd: u8) -> Result<(), &'static str> {
        self.command_port.write_u8(cmd);
        Ok(())
    }

    /// 상태 확인 및 대기
    fn wait_ready(&self) -> Result<(), &'static str> {
        // 최대 30회 반복 (약 3초 대기)
        for _ in 0..30 {
            let status = self.status_port.read_u8();

            // 에러 확인
            if (status & ATA_STATUS_ERROR) != 0 {
                return Err("ATA device error");
            }

            // READY 상태 확인
            if (status & ATA_STATUS_READY) != 0 && (status & ATA_STATUS_BUSY) == 0 {
                return Ok(());
            }

            // 짧은 대기
            for _ in 0..1000 {
                unsafe { core::arch::asm!("nop"); }
            }
        }

        Err("ATA device timeout")
    }

    /// LBA 모드 설정 (LBA 주소로 섹터 지정)
    fn set_lba_mode(&self, lba: u32, sector_count: u8) {
        // LBA 모드: bit 6 = 1
        let drive_byte = 0xE0 | ((lba >> 24) & 0x0F) as u8;
        self.drive_head_port.write_u8(drive_byte);

        // 섹터 주소 설정
        self.sector_number_port.write_u8((lba & 0xFF) as u8);
        self.cylinder_low_port.write_u8(((lba >> 8) & 0xFF) as u8);
        self.cylinder_high_port.write_u8(((lba >> 16) & 0xFF) as u8);

        // 읽을 섹터 수
        self.sector_count_port.write_u8(sector_count);
    }

    /// 섹터 읽기 (LBA 주소)
    pub fn read_sector(&mut self, lba: u32) -> Result<[u8; DISK_SECTOR_SIZE], &'static str> {
        if !self.initialized {
            return Err("Disk not initialized");
        }

        if lba >= self.total_sectors {
            return Err("LBA out of range");
        }

        // LBA 모드 설정
        self.set_lba_mode(lba, 1);

        // READ PIO 명령어 전송
        self.send_command(ATA_CMD_READ_PIO)?;

        // 데이터 준비 대기
        self.wait_ready()?;

        // 512 바이트 읽기
        let mut buffer = [0u8; DISK_SECTOR_SIZE];
        for i in 0..DISK_SECTOR_SIZE / 2 {
            let word = self.data_port.read_u16();
            buffer[i * 2] = (word & 0xFF) as u8;
            buffer[i * 2 + 1] = ((word >> 8) & 0xFF) as u8;
        }

        self.sectors_read += 1;
        Ok(buffer)
    }

    /// 섹터 쓰기 (LBA 주소)
    pub fn write_sector(&mut self, lba: u32, data: &[u8; DISK_SECTOR_SIZE]) -> Result<(), &'static str> {
        if !self.initialized {
            return Err("Disk not initialized");
        }

        if lba >= self.total_sectors {
            return Err("LBA out of range");
        }

        // LBA 모드 설정
        self.set_lba_mode(lba, 1);

        // WRITE PIO 명령어 전송
        self.send_command(ATA_CMD_WRITE_PIO)?;

        // 512 바이트 쓰기
        for i in 0..DISK_SECTOR_SIZE / 2 {
            let word = (data[i * 2 + 1] as u16) << 8 | (data[i * 2] as u16);
            self.data_port.write_u16(word);
        }

        // 쓰기 완료 대기
        self.wait_ready()?;

        self.sectors_written += 1;
        Ok(())
    }

    /// 여러 섹터 읽기
    pub fn read_sectors(&mut self, lba: u32, count: u32) -> Result<Vec<u8>, &'static str> {
        let mut buffer = Vec::with_capacity((count as usize) * DISK_SECTOR_SIZE);

        for i in 0..count {
            let sector = self.read_sector(lba + i)?;
            buffer.extend_from_slice(&sector);
        }

        Ok(buffer)
    }

    /// 디스크 정보 출력
    pub fn print_status(&self) {
        crate::println!("\n💾 ATA Disk Driver Status:");
        crate::println!("   Initialized: {}", if self.initialized { "yes" } else { "no" });
        crate::println!("   Total size: {} MB", self.total_sectors / 2048);
        crate::println!("   Sectors read: {}", self.sectors_read);
        crate::println!("   Sectors written: {}", self.sectors_written);
    }
}

lazy_static! {
    /// 글로벌 디스크 드라이버
    pub static ref DISK: Mutex<DiskDriver> = {
        let mut disk = DiskDriver::new();
        let _ = disk.init();
        Mutex::new(disk)
    };
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_disk_creation() {
        let disk = DiskDriver::new();
        assert!(!disk.initialized);
        assert_eq!(disk.total_sectors, 0);
        assert_eq!(disk.sectors_read, 0);
    }

    #[test]
    fn test_disk_sector_size() {
        assert_eq!(DISK_SECTOR_SIZE, 512);
    }

    #[test]
    fn test_io_port_initialization() {
        let disk = DiskDriver::new();
        assert_eq!(disk.data_port.port(), ATA_PRIMARY_DATA);
        assert_eq!(disk.status_port.port(), ATA_PRIMARY_STATUS);
    }

    #[test]
    fn test_ata_constants() {
        assert_eq!(ATA_CMD_READ_PIO, 0x20);
        assert_eq!(ATA_CMD_WRITE_PIO, 0x30);
        assert_eq!(ATA_STATUS_BUSY, 0x80);
    }
}
