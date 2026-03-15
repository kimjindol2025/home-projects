// 키보드 드라이버
// Phase 4: PS/2 키보드 입력 처리

use crate::io::IOPort;
use alloc::vec::Vec;
use spin::Mutex;
use lazy_static::lazy_static;

extern crate alloc;

/// 키보드 포트 (PS/2)
const KEYBOARD_DATA_PORT: u16 = 0x60;
const KEYBOARD_STATUS_PORT: u16 = 0x64;

/// 키보드 상태 플래그
const KB_OUTPUT_FULL: u8 = 0x01;  // 출력 버퍼에 데이터 있음
const KB_INPUT_FULL: u8 = 0x02;   // 입력 버퍼가 가득 참

/// PS/2 키보드 스캔 코드 → ASCII 변환
pub struct KeyboardDriver {
    data_port: IOPort,
    status_port: IOPort,
    shift_pressed: bool,
    ctrl_pressed: bool,
    alt_pressed: bool,
    input_buffer: Vec<u8>,
}

impl KeyboardDriver {
    /// 새 키보드 드라이버 생성
    pub fn new() -> Self {
        KeyboardDriver {
            data_port: IOPort::new(KEYBOARD_DATA_PORT),
            status_port: IOPort::new(KEYBOARD_STATUS_PORT),
            shift_pressed: false,
            ctrl_pressed: false,
            alt_pressed: false,
            input_buffer: Vec::new(),
        }
    }

    /// 키보드 초기화
    pub fn init(&mut self) {
        // PS/2 컨트롤러 활성화 (기본 설정)
        crate::println!("🔧 PS/2 Keyboard initialized");
    }

    /// 키보드 인터럽트 처리
    pub fn handle_interrupt(&mut self) -> Option<u8> {
        // 출력 버퍼에 데이터가 있는지 확인
        if (self.status_port.read_u8() & KB_OUTPUT_FULL) == 0 {
            return None;
        }

        let scancode = self.data_port.read_u8();

        // 스캔 코드 처리
        match scancode {
            // Shift 키
            0x2A => {
                self.shift_pressed = true;
                None
            }
            0xAA => {
                self.shift_pressed = false;
                None
            }
            // Ctrl 키
            0x1D => {
                self.ctrl_pressed = true;
                None
            }
            0x9D => {
                self.ctrl_pressed = false;
                None
            }
            // Alt 키
            0x38 => {
                self.alt_pressed = true;
                None
            }
            0xB8 => {
                self.alt_pressed = false;
                None
            }
            // 일반 키 (Make code만, Break code는 무시)
            code if code < 0x80 => Some(self.scancode_to_ascii(code)),
            // Break code (키를 뗄 때)
            _code => None,
        }
    }

    /// 스캔 코드를 ASCII 문자로 변환
    fn scancode_to_ascii(&self, scancode: u8) -> u8 {
        let base_char = match scancode {
            0x02 => b'1',
            0x03 => b'2',
            0x04 => b'3',
            0x05 => b'4',
            0x06 => b'5',
            0x07 => b'6',
            0x08 => b'7',
            0x09 => b'8',
            0x0A => b'9',
            0x0B => b'0',
            0x0C => b'-',
            0x0D => b'=',
            0x0E => 0x08, // Backspace
            0x0F => b'\t', // Tab
            0x10 => b'q',
            0x11 => b'w',
            0x12 => b'e',
            0x13 => b'r',
            0x14 => b't',
            0x15 => b'y',
            0x16 => b'u',
            0x17 => b'i',
            0x18 => b'o',
            0x19 => b'p',
            0x1A => b'[',
            0x1B => b']',
            0x1C => b'\n', // Enter
            0x1E => b'a',
            0x1F => b's',
            0x20 => b'd',
            0x21 => b'f',
            0x22 => b'g',
            0x23 => b'h',
            0x24 => b'j',
            0x25 => b'k',
            0x26 => b'l',
            0x27 => b';',
            0x28 => b'\'',
            0x29 => b'`',
            0x2B => b'\\',
            0x2C => b'z',
            0x2D => b'x',
            0x2E => b'c',
            0x2F => b'v',
            0x30 => b'b',
            0x31 => b'n',
            0x32 => b'm',
            0x33 => b',',
            0x34 => b'.',
            0x35 => b'/',
            0x39 => b' ', // Space
            _ => return 0, // 알 수 없는 코드
        };

        // Shift + 문자는 대문자
        if self.shift_pressed && base_char >= b'a' && base_char <= b'z' {
            base_char - 32 // 대문자로 변환
        } else if self.shift_pressed {
            // Shift + 숫자는 특수 문자
            match scancode {
                0x02 => b'!',
                0x03 => b'@',
                0x04 => b'#',
                0x05 => b'$',
                0x06 => b'%',
                0x07 => b'^',
                0x08 => b'&',
                0x09 => b'*',
                0x0A => b'(',
                0x0B => b')',
                _ => base_char,
            }
        } else {
            base_char
        }
    }

    /// 입력 버퍼에 문자 추가
    pub fn push_char(&mut self, ch: u8) {
        self.input_buffer.push(ch);
    }

    /// 입력 버퍼에서 한 줄 읽기
    pub fn read_line(&mut self) -> Option<alloc::string::String> {
        use alloc::string::String;

        // 개행 문자 찾기
        if let Some(pos) = self.input_buffer.iter().position(|&c| c == b'\n') {
            let line_bytes: Vec<u8> = self.input_buffer.drain(0..=pos).collect();
            // 개행 제거 후 String으로 변환
            let line = String::from_utf8_lossy(&line_bytes[..line_bytes.len() - 1]);
            return Some(line.into_owned());
        }

        None
    }

    /// 입력 버퍼 상태 출력
    pub fn print_status(&self) {
        crate::println!("\n🎹 Keyboard Status:");
        crate::println!("   Shift: {}", if self.shift_pressed { "pressed" } else { "released" });
        crate::println!("   Ctrl: {}", if self.ctrl_pressed { "pressed" } else { "released" });
        crate::println!("   Alt: {}", if self.alt_pressed { "pressed" } else { "released" });
        crate::println!("   Input buffer: {} bytes", self.input_buffer.len());
    }
}

lazy_static! {
    /// 글로벌 키보드 드라이버
    pub static ref KEYBOARD: Mutex<KeyboardDriver> = {
        let mut kb = KeyboardDriver::new();
        kb.init();
        Mutex::new(kb)
    };
}

/// 키보드 인터럽트 핸들러 (interrupts.rs에서 호출)
pub fn handle_keyboard_interrupt() {
    let mut kb = KEYBOARD.lock();
    if let Some(ascii_code) = kb.handle_interrupt() {
        kb.push_char(ascii_code);

        // 특수 키 처리
        match ascii_code {
            b'\n' => {
                // 개행은 입력 완료
                crate::println!("[KBD] Line received");
            }
            0x08 => {
                // Backspace
                crate::print!("\x08 \x08"); // VGA 버퍼의 마지막 문자 제거
            }
            c if c >= 32 && c < 127 => {
                // 출력 가능한 문자
                crate::print!("{}", c as char);
            }
            _ => {}
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_keyboard_creation() {
        let kb = KeyboardDriver::new();
        assert!(!kb.shift_pressed);
        assert!(!kb.ctrl_pressed);
        assert!(!kb.alt_pressed);
    }

    #[test]
    fn test_scancode_to_ascii() {
        let kb = KeyboardDriver::new();

        // 'a' 키
        assert_eq!(kb.scancode_to_ascii(0x1E), b'a');

        // 숫자 키
        assert_eq!(kb.scancode_to_ascii(0x02), b'1');
        assert_eq!(kb.scancode_to_ascii(0x03), b'2');

        // 공백
        assert_eq!(kb.scancode_to_ascii(0x39), b' ');
    }

    #[test]
    fn test_shift_modifier() {
        let mut kb = KeyboardDriver::new();
        kb.shift_pressed = true;

        // Shift + 'a'는 'A'
        assert_eq!(kb.scancode_to_ascii(0x1E), b'A');

        // Shift + '1'은 '!'
        assert_eq!(kb.scancode_to_ascii(0x02), b'!');
    }

    #[test]
    fn test_input_buffer() {
        let mut kb = KeyboardDriver::new();
        kb.push_char(b'h');
        kb.push_char(b'i');
        kb.push_char(b'\n');

        let line = kb.read_line();
        assert!(line.is_some());
        assert_eq!(line.unwrap(), "hi");
    }
}
