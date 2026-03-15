#!/bin/bash

# FreeLang OS Kernel 빌드 스크립트
# 용법: ./build.sh [debug|release]

set -e

BUILD_TYPE=${1:-release}
KERNEL_NAME="freelang-kernel"

echo "════════════════════════════════════════════════════════════"
echo "🚀 FreeLang OS Kernel Build Script"
echo "════════════════════════════════════════════════════════════"
echo ""

# 1. 디렉토리 준비
echo "[1/5] Preparing directories..."
mkdir -p target/iso/boot/grub
mkdir -p target/objects

# 2. 어셈블리 컴파일
echo "[2/5] Compiling assembly..."
if ! command -v nasm &> /dev/null; then
    echo "❌ NASM not found. Install with: sudo apt install nasm"
    exit 1
fi

nasm -f elf64 src/boot.asm -o target/objects/boot.o
nasm -f elf64 src/boot.asm -o target/objects/boot64.o
echo "   ✓ boot.asm -> boot.o"
echo "   ✓ boot.asm -> boot64.o"

# 3. Rust 커널 빌드
echo "[3/5] Compiling Rust kernel..."
if ! rustup target list | grep -q "x86_64-unknown-none (installed)"; then
    echo "Installing x86_64-unknown-none target..."
    rustup target add x86_64-unknown-none
fi

if [ "$BUILD_TYPE" = "debug" ]; then
    cargo build --target x86_64-unknown-none 2>&1 | tail -5
    KERNEL_BIN="target/x86_64-unknown-none/debug/kernel"
else
    cargo build --release --target x86_64-unknown-none 2>&1 | tail -5
    KERNEL_BIN="target/x86_64-unknown-none/release/kernel"
fi

if [ ! -f "$KERNEL_BIN" ]; then
    echo "❌ Kernel binary not found at $KERNEL_BIN"
    exit 1
fi

echo "   ✓ Rust kernel compiled: $KERNEL_BIN"

# 4. 링킹
echo "[4/5] Linking kernel..."

# x86_64-linux-gnu-ld 대신 ld.lld 사용 (더 안정적)
if command -v ld.lld &> /dev/null; then
    LINKER="ld.lld"
elif command -v x86_64-linux-gnu-ld &> /dev/null; then
    LINKER="x86_64-linux-gnu-ld"
else
    echo "❌ Linker not found"
    exit 1
fi

$LINKER \
    -n \
    --gc-sections \
    -T linker.ld \
    target/objects/boot.o \
    $KERNEL_BIN \
    -o target/kernel.bin

if [ ! -f target/kernel.bin ]; then
    echo "❌ Linking failed"
    exit 1
fi

echo "   ✓ Linked kernel: target/kernel.bin"
echo "   📊 Size: $(du -h target/kernel.bin | cut -f1)"

# 5. ISO 생성
echo "[5/5] Creating bootable ISO..."

# GRUB 설정
cat > target/iso/boot/grub/grub.cfg << 'EOF'
set default=0
set timeout=0
set gfxmode=1024x768

menuentry "FreeLang OS" {
    multiboot2 /boot/kernel.bin
    boot
}
EOF

cp target/kernel.bin target/iso/boot/kernel.bin

if command -v grub-mkrescue &> /dev/null; then
    grub-mkrescue -o target/freelang.iso target/iso 2>&1 | grep -v "Warning:" || true
    FINAL_ISO="target/freelang.iso"
else
    echo "⚠️  grub-mkrescue not found, trying mkisofs..."
    if command -v mkisofs &> /dev/null; then
        mkisofs -R -b boot/grub/stage2_eltorito \
            -no-emul-boot -boot-load-size 4 \
            -o target/freelang.iso target/iso
        FINAL_ISO="target/freelang.iso"
    else
        echo "⚠️  ISO creation skipped (install grub-pc or cdrtools)"
        FINAL_ISO=""
    fi
fi

echo ""
echo "════════════════════════════════════════════════════════════"
echo "✅ Build completed!"
echo "════════════════════════════════════════════════════════════"
echo ""

if [ -n "$FINAL_ISO" ]; then
    echo "📦 Bootable ISO: $FINAL_ISO"
    echo "   Size: $(du -h $FINAL_ISO | cut -f1)"
    echo ""
    echo "🚀 To run in QEMU:"
    echo "   qemu-system-x86_64 -cdrom $FINAL_ISO -m 512 -serial stdio"
    echo ""
    echo "📝 Other QEMU options:"
    echo "   -enable-kvm         # 하드웨어 가속 활성화"
    echo "   -smp 2              # 2개 CPU 코어"
    echo "   -s                  # GDB 디버깅 포트 1234 열기"
    echo "   -S                  # 부팅 시 멈춤"
fi
