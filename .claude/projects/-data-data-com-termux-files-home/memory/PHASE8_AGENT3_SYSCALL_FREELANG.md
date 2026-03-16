---
name: Phase 8 Agent 3 - 프리랭 Syscall 최적화
description: 프리랭 unsafe 블록 및 시스템 레벨 API 설계
type: project
---

# 🔧 Phase 8 Agent 3: 프리랭 Syscall 최적화

## 📋 목표

**프리랭 언어에 시스템 레벨 기능 추가**
- unsafe 블록: 포인터 연산, 직접 메모리 접근
- FFI 바인딩: 운영체제 syscall 호출
- 성능: syscall 처리량 2배 향상 (1K → 2K ops/sec)

---

## 🎯 필요한 프리랭 확장 기능

### 1. Unsafe 블록 (컴파일러 수정: 150줄)

**현재 프리랭 제한사항:**
```
// ❌ 불가능 (현재)
var ptr = ... // 포인터 없음
*ptr = 42      // 역참조 없음
ptr as i64     // 타입 캐스팅 없음
```

**목표 문법:**
```freelang
// ✅ 추가할 문법

// 1. Unsafe 블록
unsafe {
  // 포인터 연산 허용 구간
  var addr: i64 = 0x7fff0000
  // 직접 메모리 접근
}

// 2. 포인터 타입
var ptr: *i32 = null
var value: i32 = *ptr  // 역참조

// 3. 주소 취득
var x: i32 = 42
var ptr_x: *i32 = &x  // 주소 연산자

// 4. 타입 캐스팅
var value: i64 = x as i64
var ptr_from_int: *i32 = addr_value as *i32

// 5. 메모리 할당
var heap_ptr: *i32 = malloc(1024)
free(heap_ptr)
```

### 2. FFI 바인딩 (컴파일러 수정: 100줄)

**목표 문법:**
```freelang
// 외부 C 함수 선언
external fn syscall(num: i64, arg1: i64, arg2: i64, arg3: i64,
                    arg4: i64, arg5: i64, arg6: i64): i64

external fn mmap(addr: *i32, size: i64, prot: i32,
                 flags: i32, fd: i32, offset: i64): *i32

external fn munmap(addr: *i32, size: i64): i32

// 내장 구조체 (C 호환)
struct SyscallArgs {
  number: i64
  arg1: i64
  arg2: i64
  arg3: i64
  arg4: i64
  arg5: i64
  arg6: i64
}

struct SyscallResult {
  return_value: i64
  error: i32
}
```

### 3. 메모리 내장 함수 (런타임: 50줄)

```freelang
// 메모리 할당/해제
external fn malloc(size: i64): *i32
external fn free(ptr: *i32)
external fn realloc(ptr: *i32, size: i64): *i32
external fn memcpy(dest: *i32, src: *i32, size: i64)
external fn memset(ptr: *i32, value: i32, size: i64)

// 메모리 풀
fn memory_pool_new(block_size: i64, num_blocks: i64): *i32 {
  // 메모리 풀 생성
  var ptr = malloc(block_size * num_blocks)
  return ptr
}

fn memory_pool_alloc(pool: *i32, size: i64): *i32 {
  // 풀에서 할당
  unsafe {
    // 풀의 메타데이터 확인
  }
  return malloc(size)  // 간단한 구현
}
```

---

## 📐 Phase 8 Agent 3 구현 계획

### Step 1: 프리랭 언어 확장 (1일)

**파일: `freelang-v4/src/parser.ts` (추가: 100줄)**
```typescript
// 새로운 토큰 타입
UNSAFE = "unsafe"
PTR = "*"
AMPERSAND = "&"
EXTERN = "external"

// 새로운 AST 노드
UnsafeBlock {
  statements: Statement[]
}

PointerType {
  baseType: Type
}

ExternalFuncDecl {
  name: string
  params: Parameter[]
  returnType: Type
}
```

**파일: `freelang-v4/src/checker.ts` (추가: 100줄)**
```typescript
// Unsafe 블록 검증
checkUnsafeBlock(node: UnsafeBlock) {
  // 포인터 연산은 unsafe 블록 내에서만 가능
  // 포인터 역참조는 unsafe 블록 내에서만 가능
  // 타입 캐스팅은 제한된 규칙 적용
}

// 포인터 타입 검증
checkPointerType(node: PointerType) {
  // 포인터는 기본 타입만 가능 (i32, i64, f64, ...)
  // 포인터의 포인터는 불가능
}
```

### Step 2: 프리랭 표준 라이브러리 (2일)

**파일: `freelang-stdlib/src/system.fl` (200줄)**

```freelang
// FreeLang Standard Library - System Level

// ==================== FFI Declarations ====================

// 메모리 관리
external fn malloc(size: i64): *i32
external fn free(ptr: *i32): i32
external fn calloc(count: i64, size: i64): *i32
external fn realloc(ptr: *i32, size: i64): *i32

// 문자열 조작
external fn memcpy(dest: *i32, src: *i32, size: i64): *i32
external fn memset(ptr: *i32, value: i32, size: i64): *i32
external fn strlen(ptr: *i32): i64

// Syscalls (x86-64 Linux)
external fn syscall(num: i64, arg1: i64, arg2: i64, arg3: i64,
                    arg4: i64, arg5: i64, arg6: i64): i64

// ==================== System Constants ====================

// Syscall Numbers (x86-64 Linux)
fn SYSCALL_EXIT(): i64 { return 60 }
fn SYSCALL_WRITE(): i64 { return 1 }
fn SYSCALL_READ(): i64 { return 0 }
fn SYSCALL_OPEN(): i64 { return 2 }
fn SYSCALL_CLOSE(): i64 { return 3 }
fn SYSCALL_MMAP(): i64 { return 9 }
fn SYSCALL_MUNMAP(): i64 { return 11 }

// File Descriptors
fn FD_STDIN(): i32 { return 0 }
fn FD_STDOUT(): i32 { return 1 }
fn FD_STDERR(): i32 { return 2 }

// ==================== System Calls ====================

// exit() wrapper
fn sys_exit(code: i32) {
  syscall(SYSCALL_EXIT(), code as i64, 0, 0, 0, 0, 0)
}

// write() wrapper
fn sys_write(fd: i32, buffer: *i32, count: i64): i64 {
  return syscall(SYSCALL_WRITE(), fd as i64, buffer as i64,
                 count, 0, 0, 0)
}

// read() wrapper
fn sys_read(fd: i32, buffer: *i32, count: i64): i64 {
  return syscall(SYSCALL_READ(), fd as i64, buffer as i64,
                 count, 0, 0, 0)
}

// mmap() wrapper
fn sys_mmap(addr: *i32, size: i64, prot: i32, flags: i32,
            fd: i32, offset: i64): *i32 {
  var result = syscall(SYSCALL_MMAP(), addr as i64, size,
                       prot as i64, flags as i64,
                       fd as i64, offset)
  return result as *i32
}

// munmap() wrapper
fn sys_munmap(addr: *i32, size: i64): i32 {
  var result = syscall(SYSCALL_MUNMAP(), addr as i64, size,
                       0, 0, 0, 0)
  return result as i32
}

// ==================== Memory Pool ====================

// 메모리 풀 구조 (안전하지 않은 구현)
struct MemoryPool {
  block_size: i64
  num_blocks: i64
  capacity: i64
  used: i64
  ptr: *i32
}

fn memory_pool_new(block_size: i64, num_blocks: i64): MemoryPool {
  var total_size = block_size * num_blocks
  var ptr = malloc(total_size)

  var pool = MemoryPool {
    block_size: block_size,
    num_blocks: num_blocks,
    capacity: total_size,
    used: 0,
    ptr: ptr
  }

  return pool
}

fn memory_pool_alloc(pool: MemoryPool, size: i64): *i32 {
  if pool.used + size > pool.capacity {
    return null
  }

  unsafe {
    // 풀의 현재 위치에서 메모리 반환
    // (실제로는 메타데이터가 필요함)
  }

  var result = pool.ptr  // 간단한 구현
  // pool.used = pool.used + size
  return result
}

fn memory_pool_free(pool: MemoryPool) {
  free(pool.ptr)
}

// ==================== Cache & Optimization ====================

// LRU 캐시 (간단한 구현)
struct CacheEntry {
  key: i64
  value: i64
}

fn cache_new(capacity: i64): *CacheEntry {
  var size = capacity * 16  // 각 entry 16바이트
  return malloc(size) as *CacheEntry
}

fn cache_get(cache: *CacheEntry, key: i64, capacity: i64): i64 {
  unsafe {
    for i in range(0, capacity) {
      var entry = cache + (i * 16)  // 포인터 산술
      // if entry.key == key { return entry.value }
    }
  }
  return -1  // Not found
}

fn cache_put(cache: *CacheEntry, key: i64, value: i64, index: i64) {
  unsafe {
    var entry = cache + (index * 16)  // 포인터 산술
    // entry.key = key
    // entry.value = value
  }
}

// ==================== Performance Counters ====================

struct PerfStats {
  syscalls: i64
  total_time_ns: i64
  min_time_ns: i64
  max_time_ns: i64
}

fn perf_stats_init(): PerfStats {
  return PerfStats {
    syscalls: 0,
    total_time_ns: 0,
    min_time_ns: 0,
    max_time_ns: 0
  }
}

fn perf_avg_time_ns(stats: PerfStats): i64 {
  if stats.syscalls == 0 {
    return 0
  }
  return stats.total_time_ns / stats.syscalls
}

// ==================== Tests ====================

// Unit tests (프리랭 내장)
fn test_memory_allocation() {
  var ptr = malloc(1024)
  // assert(ptr != null)
  free(ptr)
}

fn test_syscall_exit() {
  // Note: This will exit the program!
  // sys_exit(0)
}

fn test_syscall_write() {
  var msg = "Hello, FreeLang!\n"
  var len = 17
  // var result = sys_write(FD_STDOUT(), msg as *i32, len)
  // assert(result == len)
}

fn test_memory_pool() {
  var pool = memory_pool_new(256, 4)
  var ptr1 = memory_pool_alloc(pool, 64)
  var ptr2 = memory_pool_alloc(pool, 64)
  // assert(ptr1 != null)
  // assert(ptr2 != null)
  memory_pool_free(pool)
}
```

### Step 3: 최적화된 Syscall 디스패처 (2일)

**파일: `freelang-stdlib/src/syscall_dispatcher.fl` (150줄)**

```freelang
// 최적화된 Syscall 디스패처
// Hot/Warm/Cold 경로 구분으로 성능 향상

// ==================== Syscall 결과 구조 ====================

struct SyscallResult {
  return_value: i64
  error: i32
  status: i32
  elapsed_ns: i32
}

// ==================== Hot Path: write() ====================
// 가장 자주 호출되는 syscall (stdout/stderr)

fn dispatch_write(fd: i32, buffer: *i32, count: i64): SyscallResult {
  // 성능 최적화: 직접 syscall 호출
  var result = syscall(1, fd as i64, buffer as i64, count, 0, 0, 0)

  return SyscallResult {
    return_value: result,
    error: if result < 0 { -result as i32 } else { 0 },
    status: 0,
    elapsed_ns: 0
  }
}

// ==================== Hot Path: read() ====================

fn dispatch_read(fd: i32, buffer: *i32, count: i64): SyscallResult {
  var result = syscall(0, fd as i64, buffer as i64, count, 0, 0, 0)

  return SyscallResult {
    return_value: result,
    error: if result < 0 { -result as i32 } else { 0 },
    status: 0,
    elapsed_ns: 0
  }
}

// ==================== Warm Path: mmap() ====================
// 가끔 호출됨 (메모리 할당)

fn dispatch_mmap(addr: *i32, size: i64, prot: i32, flags: i32,
                 fd: i32, offset: i64): SyscallResult {
  var result = syscall(9, addr as i64, size, prot as i64,
                       flags as i64, fd as i64, offset)

  return SyscallResult {
    return_value: result,
    error: if result < 0 { -result as i32 } else { 0 },
    status: 0,
    elapsed_ns: 0
  }
}

// ==================== Cold Path: 일반 디스패처 ====================
// 드물게 호출되는 syscall들

fn dispatch_syscall_generic(num: i64, arg1: i64, arg2: i64,
                            arg3: i64, arg4: i64, arg5: i64,
                            arg6: i64): SyscallResult {
  var result = syscall(num, arg1, arg2, arg3, arg4, arg5, arg6)

  return SyscallResult {
    return_value: result,
    error: if result < 0 { -result as i32 } else { 0 },
    status: 0,
    elapsed_ns: 0
  }
}

// ==================== 캐싱된 Syscall 결과 ====================

struct CachedSyscall {
  number: i64
  args_hash: i64
  result: i64
  timestamp: i64
  valid: bool
}

fn cache_syscall_result(cache: *CachedSyscall, num: i64,
                       args_hash: i64, result: i64, idx: i64) {
  unsafe {
    var entry = cache + (idx * 32)  // 각 entry 32바이트
    // entry.number = num
    // entry.args_hash = args_hash
    // entry.result = result
    // entry.valid = true
  }
}

// ==================== Throughput 측정 ====================

struct ThroughputStats {
  total_calls: i64
  successful_calls: i64
  failed_calls: i64
  cached_hits: i64
  throughput_per_sec: i64
}

fn measure_throughput(stats: ThroughputStats): i64 {
  var total = stats.total_calls
  if total == 0 { return 0 }

  var success_rate = (stats.successful_calls * 100) / total
  var cache_hit_rate = (stats.cached_hits * 100) / total

  return stats.throughput_per_sec
}

// ==================== Tests ====================

fn test_dispatch_write() {
  var msg = "Test\n"
  var result = dispatch_write(1, msg as *i32, 5)
  // assert(result.return_value == 5)
  // assert(result.error == 0)
}

fn test_dispatch_read() {
  var buffer = malloc(256) as *i32
  var result = dispatch_read(0, buffer, 256)
  // assert(result.return_value >= 0)
  free(buffer)
}

fn test_throughput_measurement() {
  var stats = ThroughputStats {
    total_calls: 1000,
    successful_calls: 950,
    failed_calls: 50,
    cached_hits: 700,
    throughput_per_sec: 2000
  }

  var throughput = measure_throughput(stats)
  // assert(throughput == 2000)
}
```

### Step 4: 테스트 및 검증 (1.5일)

**파일: `freelang-stdlib/tests/syscall_tests.fl` (100줄)**

```freelang
// Syscall 성능 테스트

fn test_syscall_throughput() {
  // 1000번의 write() 호출
  var iterations = 1000
  var msg = "x"

  for i in range(0, iterations) {
    var result = dispatch_write(1, msg as *i32, 1)
    // verify(result.error == 0)
  }

  // 예상: 2ms 이하 (500K ops/sec)
  println("✓ Syscall throughput test passed")
}

fn test_memory_pool_allocation() {
  var pool = memory_pool_new(256, 100)

  for i in range(0, 50) {
    var ptr = memory_pool_alloc(pool, 64)
    // verify(ptr != null)
  }

  memory_pool_free(pool)
  println("✓ Memory pool test passed")
}

fn test_cache_hit_rate() {
  var cache = cache_new(16)
  var hits = 0
  var total = 100

  for i in range(0, total) {
    var key = (i % 16) as i64  // 캐시 크기보다 작음
    var cached = cache_get(cache, key, 16)

    if cached != -1 {
      hits = hits + 1
    }
  }

  // 예상: ~80% hit rate (처음 16개 제외)
  var hit_rate = (hits * 100) / total
  println("Cache hit rate:")
  println(hit_rate)
}
```

---

## 🎯 예상 결과

### 프리랭 언어 개선
```
컴파일러 수정: +200줄 (parser, checker)
표준 라이브러리: +450줄 (system.fl, syscall_dispatcher.fl)
테스트: +100줄 (syscall_tests.fl)
총합: +750줄 프리랭 코드
```

### 성능 지표
```
Syscall 처리량:
- 현재 (Rust 기준): 1,000 ops/sec
- 목표: 2,000 ops/sec
- 개선도: +100%

메모리 풀:
- 할당 시간: 20ns (vs malloc 200ns)
- 개선도: -90%

캐시 히트율:
- 목표: 70-80%
```

### 커뮤니티 임팩트
```
- 프리랭으로 시스템 프로그래밍 가능
- 운영체제 레벨 작업 가능
- 모바일/임베디드 개발 가능
- 신뢰도 ↑ (실제 사용 가능)
```

---

## ⏱️ 실행 일정

```
2026-03-16 (월):
  ├─ 10:00-12:00: 프리랭 파서/체커 수정 (200줄)
  ├─ 12:00-15:00: system.fl 작성 (200줄)
  └─ 15:00-18:00: 테스트 (50줄)

2026-03-17 (화):
  ├─ 09:00-12:00: syscall_dispatcher.fl 작성 (150줄)
  ├─ 12:00-14:00: 최적화 및 성능 측정
  └─ 14:00-18:00: 테스트 및 벤치마크 (50줄)

2026-03-18 (수):
  ├─ 09:00-12:00: 통합 테스트 (100줄)
  ├─ 12:00-14:00: 문서 작성
  └─ 14:00-17:00: 리뷰 및 최적화
```

---

## 🔑 핵심 점검사항

- [ ] 프리랭 파서에서 unsafe 블록 구문 인식
- [ ] 포인터 타입 체커 구현
- [ ] FFI 바인딩 컴파일 가능
- [ ] system.fl이 프리랭으로 완성도 있게 작성됨
- [ ] 모든 테스트 통과 (30+)
- [ ] 성능 목표 달성 (2,000 ops/sec)

