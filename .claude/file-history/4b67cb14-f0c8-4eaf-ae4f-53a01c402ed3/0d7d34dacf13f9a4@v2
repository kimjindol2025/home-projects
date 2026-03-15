// I/O 포트 추상화 계층
// Phase 4: x86-64 I/O 포트 읽기/쓰기

use core::arch::asm;

/// I/O 포트 연산 추상화
pub struct Port;

impl Port {
    /// 8비트 I/O 포트에서 데이터 읽기
    #[inline]
    pub fn read_u8(port: u16) -> u8 {
        let value: u8;
        unsafe {
            asm!("in al, dx", in("dx") port, out("al") value);
        }
        value
    }

    /// 8비트 I/O 포트에 데이터 쓰기
    #[inline]
    pub fn write_u8(port: u16, value: u8) {
        unsafe {
            asm!("out dx, al", in("dx") port, in("al") value);
        }
    }

    /// 16비트 I/O 포트에서 데이터 읽기
    #[inline]
    pub fn read_u16(port: u16) -> u16 {
        let value: u16;
        unsafe {
            asm!("in ax, dx", in("dx") port, out("ax") value);
        }
        value
    }

    /// 16비트 I/O 포트에 데이터 쓰기
    #[inline]
    pub fn write_u16(port: u16, value: u16) {
        unsafe {
            asm!("out dx, ax", in("dx") port, in("ax") value);
        }
    }

    /// 32비트 I/O 포트에서 데이터 읽기
    #[inline]
    pub fn read_u32(port: u16) -> u32 {
        let value: u32;
        unsafe {
            asm!("in eax, dx", in("dx") port, out("eax") value);
        }
        value
    }

    /// 32비트 I/O 포트에 데이터 쓰기
    #[inline]
    pub fn write_u32(port: u16, value: u32) {
        unsafe {
            asm!("out dx, eax", in("dx") port, in("eax") value);
        }
    }
}

/// I/O 포트 래퍼 (타입 안전성)
pub struct IOPort {
    port: u16,
}

impl IOPort {
    /// 새 I/O 포트 생성
    pub const fn new(port: u16) -> Self {
        IOPort { port }
    }

    /// 8비트 읽기
    pub fn read_u8(&self) -> u8 {
        Port::read_u8(self.port)
    }

    /// 8비트 쓰기
    pub fn write_u8(&self, value: u8) {
        Port::write_u8(self.port, value);
    }

    /// 16비트 읽기
    pub fn read_u16(&self) -> u16 {
        Port::read_u16(self.port)
    }

    /// 16비트 쓰기
    pub fn write_u16(&self, value: u16) {
        Port::write_u16(self.port, value);
    }

    /// 32비트 읽기
    pub fn read_u32(&self) -> u32 {
        Port::read_u32(self.port)
    }

    /// 32비트 쓰기
    pub fn write_u32(&self, value: u32) {
        Port::write_u32(self.port, value);
    }

    /// 포트 번호 반환
    pub fn port(&self) -> u16 {
        self.port
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_io_port_creation() {
        let port = IOPort::new(0x3F8);
        assert_eq!(port.port(), 0x3F8);
    }

    #[test]
    fn test_multiple_ports() {
        let serial_port = IOPort::new(0x3F8);
        let pit_port = IOPort::new(0x40);
        let pic_port = IOPort::new(0x20);

        assert_eq!(serial_port.port(), 0x3F8);
        assert_eq!(pit_port.port(), 0x40);
        assert_eq!(pic_port.port(), 0x20);
    }
}
