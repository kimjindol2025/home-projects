# FreeLang OS Kernel - Phase 1 구현

**상태**: ✅ Phase 1 완료 (Multiboot2 부트로더 + 실제 x86-64 OS)

---

## 📁 프로젝트 구조

```
freelang-os-kernel/
├── src/
│   ├── main.rs              # 커널 메인 진입점
│   ├── boot.asm             # Multiboot2 부트로더 (어셈블리)
│   ├── serial.rs            # QEMU 시리얼 포트 출력
│   ├── vga_buffer.rs        # VGA 텍스트 모드 (80x25)
│   ├── memory.rs            # 물리 메모리 관리자
│   ├── gdt.rs               # Global Descriptor Table
│   └── interrupts.rs        # IDT 및 인터럽트 핸들러
│
├── .cargo/
│   └── config.toml          # Cargo 설정 (bare-metal)
│
├── x86_64-unknown-none.json # x86_64 타겟 스펙
├── linker.ld                # 링커 스크립트
├── Cargo.toml               # 프로젝트 메타데이터
├── build.sh                 # 빌드 스크립트 (자동 컴파일)
└── README_PHASE1.md         # 이 파일
```

---

## 🔧 빌드 전 요구사항

### 1. Rust 설정
```bash
# Rust 및 Cargo 설치
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# x86_64 bare-metal 타겟 추가
rustup target add x86_64-unknown-none

# LLD 링커 설치 (권장)
rustup component add llvm-tools
```

### 2. 시스템 도구
```bash
# NASM 어셈블러
sudo apt update
sudo apt install nasm

# QEMU (실행용)
sudo apt install qemu-system-x86

# GRUB (ISO 생성용, 선택사항)
sudo apt install grub-pc xorriso

# 링커 (선택사항)
sudo apt install binutils-x86-64-linux-gnu
```

### 3. macOS 사용자
```bash
# Homebrew 이용
brew install nasm qemu

# GRUB (ISO 생성)
brew install grub

# 또는 멀티부트 이미지 생성 대신 QEMU 커널 직접 로드
```

---

## 🏗️ 빌드 방법

### 방법 1: 자동 빌드 (권장)

```bash
# Release 버전 (최적화)
./build.sh release

# Debug 버전 (더 느림, 디버깅 정보 포함)
./build.sh debug
```

**결과**:
- `target/kernel.bin` - 커널 바이너리
- `target/freelang.iso` - 부팅 가능한 ISO 이미지

### 방법 2: 수동 빌드

```bash
# 1. 어셈블리 컴파일
mkdir -p target/objects
nasm -f elf64 src/boot.asm -o target/objects/boot.o

# 2. Rust 컴파일
cargo build --release --target x86_64-unknown-none

# 3. 링킹
ld.lld -T linker.ld \
    target/objects/boot.o \
    target/x86_64-unknown-none/release/kernel \
    -o target/kernel.bin

# 4. ISO 생성 (선택사항)
mkdir -p target/iso/boot/grub
cp target/kernel.bin target/iso/boot/
cat > target/iso/boot/grub/grub.cfg << 'EOF'
set default=0
set timeout=0
menuentry "FreeLang OS" {
    multiboot2 /boot/kernel.bin
    boot
}
EOF
grub-mkrescue -o target/freelang.iso target/iso
```

---

## 🚀 실행 방법

### QEMU에서 부팅

```bash
# ISO 이미지로 부팅
qemu-system-x86_64 \
    -cdrom target/freelang.iso \
    -m 512 \
    -serial stdio \
    -no-reboot

# 또는 커널 바이너리 직접 로드 (더 빠름)
qemu-system-x86_64 \
    -kernel target/kernel.bin \
    -m 512 \
    -serial stdio \
    -nographic
```

### 고급 QEMU 옵션

```bash
# 하드웨어 가속 활성화 (Linux)
qemu-system-x86_64 -cdrom target/freelang.iso -m 512 -enable-kvm -serial stdio

# 멀티코어 (2개 CPU)
qemu-system-x86_64 -cdrom target/freelang.iso -m 512 -smp 2 -serial stdio

# GDB 디버깅 활성화
qemu-system-x86_64 -cdrom target/freelang.iso -m 512 -s -S -serial stdio
# 다른 터미널에서: gdb target/x86_64-unknown-none/release/kernel
#   (gdb) target remote :1234
#   (gdb) c  # 계속 실행
```

---

## 📊 실행 결과 예상

QEMU 부팅 시:

```
╔════════════════════════════════════════════════════╗
║      FreeLang OS Kernel - Phase G Bare-Metal       ║
║          실제 x86-64 부팅 및 실행                  ║
╚════════════════════════════════════════════════════╝

[SERIAL] Serial port initialized
✓ GDT initialized
✓ IDT initialized
🔧 Initializing physical memory manager...
📊 Physical Memory Manager initialized
   Total pages: 131072 (512 MB)
   Bitmap size: 16384 bytes
✓ Physical memory manager initialized
✓ Interrupts enabled

=== 커널 부팅 완료 ===
타이머 인터럽트: 4ms마다 발생
프로세스 관리: 준비 중

⏱️ Uptime: 1s
⏱️ Uptime: 2s
⏱️ Uptime: 3s
...
```

---

## 🔍 구현 내용

### 1. Multiboot2 부트로더 (boot.asm)
- ✅ Multiboot2 헤더 (GRUB 호환)
- ✅ CPUID 확인 (Long mode 지원)
- ✅ 페이지 테이블 설정 (4KB 페이지, 2MB 매핑)
- ✅ 페이징 활성화 (CR3, CR4, EFER, CR0)
- ✅ GDT 로드 및 64-bit 모드 진입
- ✅ 커널 메인 함수 호출

### 2. Rust 커널 (main.rs)
- ✅ VGA 텍스트 출력 (초록색)
- ✅ 시리얼 포트 출력 (QEMU -serial stdio)
- ✅ Panic 핸들러
- ✅ 커널 메인 루프
- ✅ 타이머 기반 업타임 표시

### 3. GDT 초기화 (gdt.rs)
- ✅ x86_64 크레이트 사용
- ✅ 코드 세그먼트 정의
- ✅ TSS (Task State Segment) 설정
- ✅ 더블 폴트 스택 설정

### 4. IDT 및 인터럽트 (interrupts.rs)
- ✅ 예외 핸들러 (Breakpoint, Double Fault, Page Fault, GPF, Divide Error)
- ✅ PIC (Programmable Interrupt Controller) 초기화
- ✅ 타이머 인터럽트 (IRQ0, 4ms)
- ✅ 키보드 인터럽트 (IRQ1)
- ✅ EOI (End of Interrupt) 신호

### 5. 메모리 관리 (memory.rs)
- ✅ 물리 메모리 비트맵
- ✅ 페이지 할당/해제
- ✅ Fragmentation 추적
- ✅ 메모리 상태 출력

---

## 🐛 트러블슈팅

### QEMU에 grub-mkrescue 없음
```bash
# 해결책 1: GRUB 설치
sudo apt install grub-pc

# 해결책 2: 커널 바이너리 직접 로드 (더 간단)
qemu-system-x86_64 -kernel target/kernel.bin -m 512 -serial stdio
```

### 컴파일 오류: "cannot find native static library"
```bash
# LLD 설치 또는 업데이트
rustup component add llvm-tools

# Rust 버전 확인
rustup update stable
```

### QEMU 검은 화면
1. 시리얼 포트 출력 활성화: `qemu-system-x86_64 ... -serial stdio`
2. VGA 텍스트 모드 문제? `-nographic` 옵션 제거

### 64-bit 모드 전환 실패
```bash
# 부트로더 에러 코드 확인 (VGA 우측 상단)
# C = CPUID 미지원
# L = Long mode 미지원
# E = 기타 에러

# CPU 기능 확인
grep -o 'lm' /proc/cpuinfo | head -1  # "lm" 있으면 OK
```

---

## 📚 다음 단계

**Phase 2**: 메모리 관리 고도화
- [ ] Demand paging (PageFault 시 자동 할당)
- [ ] 페이지 교체 알고리즘 (메모리 부족 시)
- [ ] 힙 할당자 (malloc/free)

**Phase 3**: 프로세스 관리
- [ ] 사용자 프로세스 생성
- [ ] Context switching (실제 어셈블리)
- [ ] 프로세스 격리 (페이지 테이블)

**Phase 4+**: 고급 기능
- [ ] 파일 시스템
- [ ] 네트워크 드라이버
- [ ] 쉘 및 명령어

---

## 📖 참고 자료

- **OSDev.org**: https://wiki.osdev.org/Main_Page
- **Intel x86-64 Manual**: https://www.intel.com/content/www/en/en/developer/articles/technical/intel-sdm.html
- **x86_64 crate**: https://docs.rs/x86_64/
- **Multiboot2 Specification**: https://www.gnu.org/software/grub/manual/multiboot2/

---

**상태**: ✅ Phase 1 (Multiboot2 부트로더) 완료
**다음**: Phase 2 (메모리 관리 고도화) 진행 시 알려주세요!
