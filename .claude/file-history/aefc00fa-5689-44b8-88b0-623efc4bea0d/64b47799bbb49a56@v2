#!/bin/bash
# GRIE Kernel Build & Test Automation Script
# Supports multi-target compilation and verification
#
# Usage:
#   ./build.sh                    # Build native (ReleaseFast)
#   ./build.sh native             # Same as above
#   ./build.sh android            # Cross-compile Android ARM64
#   ./build.sh linux-x86          # Cross-compile x86_64 Linux
#   ./build.sh linux-arm          # Cross-compile ARM64 Linux
#   ./build.sh all                # Build all targets
#   ./build.sh test               # Run tests
#   ./build.sh clean              # Clean artifacts
#   ./build.sh help               # Show help

set -e  # Exit on error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ZIG_MIN_VERSION="0.14.0"
PROJECT_ROOT="../../.."
OPT="${OPT:-ReleaseFast}"

# Functions
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

check_zig() {
    if ! command -v zig &> /dev/null; then
        print_error "Zig not found. Please install Zig 0.14.0+"
        echo "Download: https://ziglang.org/download/"
        exit 1
    fi

    local zig_version=$(zig version)
    print_success "Zig found: v$zig_version"
}

build_native() {
    print_header "Building Native (${OPT})"
    local out="${PROJECT_ROOT}/bin/native"
    mkdir -p "$out"

    zig build \
        --prefix="$out" \
        -Doptimize="$OPT" \
        2>&1 | tee "$out/build.log"

    print_success "Native build complete: $out/lib"
}

build_android() {
    print_header "Cross-compiling for ARM64 Android (${OPT})"
    local out="${PROJECT_ROOT}/bin/android/aarch64"
    mkdir -p "$out"

    zig build \
        --prefix="$out" \
        -Dtarget=aarch64-linux-android \
        -Doptimize="$OPT" \
        2>&1 | tee "$out/build.log"

    print_success "Android build complete: $out/lib"
}

build_linux_x86() {
    print_header "Cross-compiling for x86_64 Linux (${OPT})"
    local out="${PROJECT_ROOT}/bin/linux/x86_64"
    mkdir -p "$out"

    zig build \
        --prefix="$out" \
        -Dtarget=x86_64-linux-gnu \
        -Doptimize="$OPT" \
        2>&1 | tee "$out/build.log"

    print_success "x86_64 Linux build complete: $out/lib"
}

build_linux_arm() {
    print_header "Cross-compiling for ARM64 Linux (${OPT})"
    local out="${PROJECT_ROOT}/bin/linux/aarch64"
    mkdir -p "$out"

    zig build \
        --prefix="$out" \
        -Dtarget=aarch64-linux-gnu \
        -Doptimize="$OPT" \
        2>&1 | tee "$out/build.log"

    print_success "ARM64 Linux build complete: $out/lib"
}

build_all() {
    print_header "Building All Targets (${OPT})"
    build_native
    build_android
    build_linux_x86
    build_linux_arm

    echo ""
    print_success "All builds complete!"
    echo ""
    echo "Build artifacts:"
    echo "  📦 Native:        ${PROJECT_ROOT}/bin/native/lib"
    echo "  📦 Android ARM64: ${PROJECT_ROOT}/bin/android/aarch64/lib"
    echo "  📦 Linux x86_64:  ${PROJECT_ROOT}/bin/linux/x86_64/lib"
    echo "  📦 Linux ARM64:   ${PROJECT_ROOT}/bin/linux/aarch64/lib"
}

run_tests() {
    print_header "Running Unit Tests (${OPT})"
    zig build test -Doptimize="$OPT"
    print_success "All tests passed"
}

clean_build() {
    print_header "Cleaning Build Artifacts"
    rm -rf "${PROJECT_ROOT}/bin"
    rm -rf zig-cache
    print_success "Clean complete"
}

show_help() {
    cat << EOF
GRIE Kernel Build Script

Usage: ./build.sh [COMMAND] [OPTIONS]

Commands:
  native              Build for native platform (default)
  android             Cross-compile for ARM64 Android
  linux-x86           Cross-compile for x86_64 Linux
  linux-arm           Cross-compile for ARM64 Linux
  all                 Build all targets
  test                Run unit tests
  clean               Clean build artifacts
  help                Show this help message

Environment Variables:
  OPT                 Optimization level (default: ReleaseFast)
                      Options: ReleaseFast, ReleaseSafe, ReleaseSmall, Debug

Examples:
  ./build.sh                          # Native build (ReleaseFast)
  ./build.sh all                      # Build all targets
  OPT=ReleaseSafe ./build.sh all      # Build all with ReleaseSafe
  ./build.sh test                     # Run tests
  ./build.sh clean && ./build.sh all  # Clean and rebuild

Optimization Levels:
  ReleaseFast         Maximum performance (default)
  ReleaseSafe         Performance + safety checks
  ReleaseSmall        Minimum binary size
  Debug               Debug symbols + no optimizations

For more information, see README.md

EOF
}

# Main
main() {
    local cmd="${1:-native}"

    check_zig

    case "$cmd" in
        native)
            build_native
            ;;
        android)
            build_android
            ;;
        linux-x86|linux-x86_64)
            build_linux_x86
            ;;
        linux-arm|linux-aarch64)
            build_linux_arm
            ;;
        all)
            build_all
            ;;
        test)
            run_tests
            ;;
        clean)
            clean_build
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $cmd"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main
main "$@"
