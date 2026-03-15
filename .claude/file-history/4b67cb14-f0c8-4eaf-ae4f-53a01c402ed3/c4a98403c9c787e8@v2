// 시리얼 포트 (COM1: 0x3F8) 드라이버
// QEMU 출력용: -serial stdio

use core::fmt::{self, Write};

const COM1: u16 = 0x3F8;

/// 시리얼 포트 초기화
pub fn init() {
    // COM1 설정
    // 115200 baud, 8 bit, no parity, 1 stop bit

    unsafe {
        // IER (Interrupt Enable Register) 해제
        outb(COM1 + 1, 0x00);

        // LCR (Line Control Register)
        // bit 7: DLAB, bit 3: DLAB, bit 1-0: word length (11 = 8 bits)
        outb(COM1 + 3, 0x80);

        // 보율 설정 (115200 = 0x0001)
        outb(COM1 + 0, 0x01);  // DLL
        outb(COM1 + 1, 0x00);  // DLM

        // LCR 다시 설정 (DLAB 해제)
        outb(COM1 + 3, 0x03);

        // FCR (FIFO Control Register)
        outb(COM1 + 2, 0xC7);

        // MCR (Modem Control Register)
        outb(COM1 + 4, 0x0B);

        // IER 활성화
        outb(COM1 + 1, 0x01);
    }
}

/// 한 문자 전송
pub fn putchar(c: u8) {
    unsafe {
        // 전송 대기
        while (inb(COM1 + 5) & 0x20) == 0 {}
        outb(COM1, c);
    }
}

/// I/O 포트에서 읽기
unsafe fn inb(port: u16) -> u8 {
    let value: u8;
    asm!("in al, dx", out("al") value, in("dx") port);
    value
}

/// I/O 포트에 쓰기
unsafe fn outb(port: u16, value: u8) {
    asm!("out dx, al", in("dx") port, in("al") value);
}

/// Write trait 구현
pub struct SerialPort;

impl fmt::Write for SerialPort {
    fn write_str(&mut self, s: &str) -> fmt::Result {
        for byte in s.bytes() {
            putchar(byte);
        }
        Ok(())
    }
}

/// 전역 시리얼 락
use spin::Mutex;
use lazy_static::lazy_static;

lazy_static! {
    pub static ref SERIAL: Mutex<SerialPort> = Mutex::new(SerialPort);
}

/// println! 매크로
#[macro_export]
macro_rules! println {
    () => ($crate::serial::_print("\n"));
    ($($arg:tt)*) => ($crate::serial::_print(&format_args!("{}\n", format_args!($($arg)*))));
}

#[macro_export]
macro_rules! print {
    ($($arg:tt)*) => ($crate::serial::_print(&format_args!($($arg)*)));
}

pub fn _print(args: &core::fmt::Arguments) {
    use core::fmt::Write;
    SERIAL.lock().write_fmt(*args).unwrap();
}

use core::arch::asm;
