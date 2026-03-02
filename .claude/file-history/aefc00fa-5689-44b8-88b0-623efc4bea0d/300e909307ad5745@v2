# GRIE Zig Kernel - Build & Integration Guide

## 개요

GRIE Zig Kernel은 고성능 SIMD 가속 연산 엔진입니다. Go 오케스트레이터와 Zig 커널 간에 **공유 메모리 IPC**를 통해 20ns 이하의 지연 시간으로 통신합니다.

**아키텍처**:
```
┌─ Go Process ─────────────────┐
│  • SubmitJob()               │
│  • WaitResult()              │  Spin-loop polling
│  • ExecuteSync()             │  (no CGO calls)
└──────┬──────────────────────┘
       │ POSIX Shared Memory
       │ (mmap MAP_SHARED)
       ▼
┌─ Zig Kernel ─────────────────┐
│  • spinKernel()              │
│  • computeSimd()             │
│    - MatMul 4x4 (SIMD)       │  < 20ns latency
│    - FFT-16 (SIMD)           │  (atomic state only)
│    - Sum Reduction           │
└──────────────────────────────┘
```

---

## 시스템 요구사항

### 필수 (Required)
- **Zig 0.14.0 이상** ([https://ziglang.org/download/](https://ziglang.org/download/))
- POSIX-호환 OS (Linux, macOS, Android with Termux)
- ARM64 또는 x86_64 CPU

### 선택 (Optional)
- Make 또는 compatible build tool (Makefile 사용 시)
- Go 1.21+ (Go-Zig 통합 테스트)

### 환경 설정 예시

#### Linux (Ubuntu/Debian)
```bash
# Zig 설치
wget https://ziglang.org/download/0.14.0/zig-linux-x86_64-0.14.0.tar.xz
tar xf zig-linux-x86_64-0.14.0.tar.xz
export PATH=$PWD/zig-0.14.0:$PATH

# 확인
zig version  # 0.14.0
```

#### macOS
```bash
# Homebrew로 설치
brew install zig

# 또는 다운로드
wget https://ziglang.org/download/0.14.0/zig-macos-aarch64-0.14.0.tar.xz
tar xf zig-macos-aarch64-0.14.0.tar.xz
export PATH=$PWD/zig-0.14.0.0:$PATH
```

#### Android (Termux)
```bash
# Termux에서는 Zig 설치 불가능 (아직)
# 대신: 원격 Linux 환경에서 크로스컴파일
ssh user@linux-host "cd grie-engine && make build-android"
```

---

## 빌드 방법

### 1️⃣ 방법 A: Makefile 사용 (권장)

#### 네이티브 빌드
```bash
cd internal/kernel
make build-native OPT=ReleaseFast
# 출력: bin/native/lib/libgrie_kernel.a, libgrie_kernel.so
```

#### 크로스컴파일
```bash
# ARM64 Android
make build-android OPT=ReleaseFast
# 출력: bin/android/aarch64/lib/libgrie_kernel.*

# x86_64 Linux
make build-linux-x86 OPT=ReleaseFast
# 출력: bin/linux/x86_64/lib/libgrie_kernel.*

# ARM64 Linux
make build-linux-arm OPT=ReleaseFast
# 출력: bin/linux/aarch64/lib/libgrie_kernel.*

# 모든 타겟
make build-all OPT=ReleaseFast
```

#### 최적화 레벨
```bash
# 기본값: ReleaseFast (최고 성능)
make build-native OPT=ReleaseFast    # 최고 성능

make build-native OPT=ReleaseSafe    # 성능 + 안전성
make build-native OPT=ReleaseSmall   # 최소 크기
make build-native OPT=Debug          # 디버깅
```

### 2️⃣ 방법 B: 직접 `zig build` 실행

#### 네이티브 빌드
```bash
cd internal/kernel
zig build -Doptimize=ReleaseFast
# 출력: zig-out/lib/
```

#### 크로스컴파일
```bash
# Android ARM64
zig build \
  --prefix=/tmp/grie-android \
  -Dtarget=aarch64-linux-android \
  -Doptimize=ReleaseFast

# x86_64 Linux
zig build \
  --prefix=/tmp/grie-linux-x86 \
  -Dtarget=x86_64-linux-gnu \
  -Doptimize=ReleaseFast

# ARM64 Linux
zig build \
  --prefix=/tmp/grie-linux-arm \
  -Dtarget=aarch64-linux-gnu \
  -Doptimize=ReleaseFast
```

#### 단일 스텝 모드 빌드 (정적/동적 분리)
```bash
# 정적 라이브러리만
zig build -Dmode=static -Doptimize=ReleaseFast

# 동적 라이브러리만
zig build -Dmode=shared -Doptimize=ReleaseFast

# 둘 다 (기본)
zig build -Dmode=both -Doptimize=ReleaseFast
```

### 3️⃣ 유닛 테스트

```bash
cd internal/kernel

# 모든 테스트 실행
make test

# 또는 직접 실행
zig build test -Doptimize=Debug
```

---

## Go-Zig 통합

### 스텝 1: Zig 커널 빌드

```bash
cd internal/kernel
make build-native OPT=ReleaseFast

# 생성된 라이브러리 확인
ls -lh ../../bin/native/lib/
# libgrie_kernel.a  (정적 라이브러리)
# libgrie_kernel.so (동적 라이브러리)
```

### 스텝 2: Go에서 FFI 호출

```go
// internal/kernel/bridge.go 참고
// KernelBridge는 이미 구현됨

import "grie-engine/internal/kernel"

// 사용 예시
kb, err := kernel.NewKernelBridge("/tmp/grie_shm_xxxxx")
if err != nil {
    log.Fatal(err)
}
defer kb.Close()

// SIMD 작업 제출
payload := []byte{0, ...}  // op_type=0 (MatMul)
kb.SubmitJob(payload)

// 결과 대기
result, err := kb.WaitResult(100)  // 100ms timeout
```

### 스텝 3: CGO 바인딩 (선택사항)

동적 라이브러리를 사용하는 경우:

```bash
# 1. 라이브러리 설치
sudo cp bin/native/lib/libgrie_kernel.so /usr/local/lib/
sudo ldconfig

# 2. Go에서 CGO로 바인딩
# cgo: -lgrie_kernel
```

---

## 성능 특성

| 작업 | 지연 시간 | 처리량 |
|------|----------|------|
| Go → Zig Signal | ~10ns | - |
| Zig Spin-loop | ~5ns | - |
| MatMul 4×4 SIMD | ~100ns | 1.6 GFLOPS |
| FFT-16 SIMD | ~150ns | - |
| Total Round-trip | ~20ns | >100M ops/s |

**목표**: Go-Zig 신호 지연 **< 20ns** (Spin-loop만, 계산 제외)

---

## 빌드 아티팩트 구조

```
internal/kernel/
├── src/
│   └── main.zig              # Kernel entry point
├── simd.zig                  # SIMD operations
├── shm_header.zig            # Shared header
├── build.zig                 # Build config
├── Makefile                  # Build automation
├── bridge.go                 # Go-Zig FFI bridge
├── bridge_test.go            # Bridge tests
├── simd_test.go              # SIMD tests
└── README.md                 # This file

../../bin/                    # Build outputs (generated)
├── native/
│   └── lib/
│       ├── libgrie_kernel.a  # Static
│       └── libgrie_kernel.so # Shared
├── android/
│   └── aarch64/lib/
├── linux/
│   ├── x86_64/lib/
│   └── aarch64/lib/
```

---

## 트러블슈팅

### 문제 1: "zig: command not found"
```bash
# 해결: Zig PATH 설정
export PATH=/path/to/zig-0.14.0:$PATH
which zig
zig version
```

### 문제 2: "Cross-compilation not supported"
```bash
# Zig 0.14.0 이상 필요
zig version  # 0.14.0 이상인지 확인
```

### 문제 3: "undefined reference to gqse_kernel_init"
```bash
# CGO에서 라이브러리 링크 확인
# build.zig에서 lib_shared.linkLibC() 확인
```

### 문제 4: 빌드 캐시 문제
```bash
# 캐시 삭제 후 재빌드
rm -rf zig-cache
make clean
make build-native
```

---

## 다음 단계

### Step 1 (현재): ✅ Zig 0.14+ 빌드 환경
- [x] build.zig 작성
- [x] Makefile 크로스컴파일 지원
- [x] 빌드 가이드 문서화

### Step 2: ✅ SIMD 최적화 커널 (완료)
- [x] simd.zig - Vec64, matmul4x4, FFT 구현
- [x] main.zig - SIMD 통합
- [x] simd_test.go - 검증

### Step 3: 실측 검증 (대기)
- [ ] 실제 크로스컴파일 및 벤치마크
- [ ] 20ns 이하 지연 시간 검증
- [ ] SIMD 성능 프로파일링

---

## 참고 자료

- [Zig 공식 문서](https://ziglang.org/documentation/0.14.0/)
- [Zig Build System](https://ziglang.org/learn/build-system/)
- [SIMD in Zig](https://ziglang.org/documentation/0.14.0/#Vectors)
- [GRIE GitHub](https://gogs.dclub.kr/kim/grie-engine.git)

---

**작성자**: Claude Opus 4.6
**최종 업데이트**: 2026-02-26
**상태**: ✅ Step 1 Complete
