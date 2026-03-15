---
name: Phase 8 Improved Agent Prompts - Permission-Aware
description: Optimized agent prompts that work within sandbox constraints and produce copy-paste-ready code
type: project
---

# 🤖 Phase 8 개선된 에이전트 프롬프트 (권한 제약 극복)

## 핵심 원칙

1. **파일 쓰기 회피** - 모든 코드는 마크다운 코드 블록으로 제시
2. **즉시 실행 가능** - 사용자가 copy-paste만 해도 동작
3. **구체적 경로** - 전체 경로명으로 명확히 표시
4. **단계별 가이드** - 수동 구현 절차 명확
5. **검증 방법** - 구현 후 확인 방법 제시

---

## 개선된 에이전트 프롬프트

### **Agent 1: Runtime 캐싱 (개선)**

```
Phase 8 Hub-Spoke Optimization: freelang-runtime 캐싱 레이어

제약 조건: 파일 쓰기 제한 환경

기대 산출물: 다음 3가지를 마크다운 코드 블록으로 제시

1️⃣ IMPLEMENTATION CODE
   파일: /data/data/com.termux/files/home/.projects/modules/freelang-runtime/src/cache.rs

   코드 블록으로 전체 350줄 제시:
   - LRU 캐시 구조 (100줄)
   - 캐시 매니저 (150줄)
   - 통계 추적 (50줄)
   - 테스트 코드 (50줄)

   ✅ Copy-paste 가능한 형식
   ✅ extern crate 포함
   ✅ #[cfg(test)] 모듈 포함

2️⃣ INTEGRATION GUIDE
   현재 위치: /data/data/com.termux/files/home/.projects/modules/freelang-runtime/src/lib.rs

   추가할 1줄:
   mod cache;

   수정할 파일: vm.rs (라인 번호 명시)
   - Line 79-80: RuntimeEngine::call_stdlib() 호출 변경
   - 변경 전/후 코드 제시

   ✅ Diff 형식으로 명확히

3️⃣ VALIDATION CHECKLIST

   [ ] Step 1: cache.rs 생성 및 코드 붙여넣기
   [ ] Step 2: lib.rs에 mod cache; 추가
   [ ] Step 3: vm.rs 수정 (Diff 따라)
   [ ] Step 4: cargo build 실행
   [ ] Step 5: cargo test --lib cache 실행
   [ ] Step 6: 성능 비교 명령어:
        cargo bench function_dispatch --baseline before

   예상 결과:
   - 테스트: 15개 통과
   - 컴파일: 0 경고
   - 성능: 30-40% 향상

기대하지 않는 것:
❌ 파일 생성 시도
❌ bash 명령어 직접 실행
❌ Write/Bash 도구 사용

기대하는 것:
✅ 완전한 코드 블록 (copy-paste 가능)
✅ 단계별 수동 구현 가이드
✅ 실제 라인 번호 참조
✅ 검증 방법 제시
```

---

### **Agent 2: Compiler 아티팩트 (개선)**

```
Phase 8 Hub-Spoke Optimization: freelang-compiler 공유 메모리

제약 조건: 파일 쓰기 제한 환경에서 즉시 구현 가능한 형식 제시

기대 산출물:

1️⃣ COMPLETE SOURCE CODE (마크다운 코드 블록)

   파일 1: shared_memory.rs (250줄)
   ```rust
   // Copy-paste 가능한 완전한 코드
   use std::sync::Arc;
   use std::collections::HashMap;
   // ... 전체 구현
   ```

   파일 2: parallel_compiler.rs (200줄)
   ```rust
   // 의존성 기반 병렬 컴파일
   ```

2️⃣ INTEGRATION STEPS

   현재 프로젝트: /data/data/com.termux/files/home/.projects/core/freelang-compiler

   Step 1: 파일 생성
   - src/shared_memory.rs ← 코드 블록 1 복사
   - src/parallel_compiler.rs ← 코드 블록 2 복사

   Step 2: lib.rs 수정 (정확한 라인 번호)
   현재 코드:
   mod cache;  // Line 15
   mod vm;

   변경 후:
   mod cache;
   mod shared_memory;    // 추가
   mod parallel_compiler; // 추가
   mod vm;

   Step 3: Cargo.toml 의존성 (필요시)
   - rayon = "1.8" (병렬화 필요시)

   Step 4: main.rs 또는 lib.rs에서 사용
   ```rust
   let compiler = ParallelCompiler::new();
   let results = compiler.compile_parallel(&files);
   ```

3️⃣ VALIDATION COMMANDS

   # Step 1: 컴파일 확인
   cargo build --release 2>&1 | grep -E "error|warning"

   # Step 2: 테스트 실행
   cargo test --lib shared_memory
   cargo test --lib parallel_compiler

   # Step 3: 성능 측정 (기존 vs 신규)
   cargo bench compiler --baseline baseline_before
   cargo bench compiler --baseline baseline_after

   예상 결과:
   ✅ 컴파일 성공 (0 에러)
   ✅ 25개 테스트 통과
   ✅ 메모리 -40%
   ✅ 처리량 +40%

4️⃣ TROUBLESHOOTING

   문제: "cannot find type 'Arc' in module 'std::sync'"
   해결: 코드 블록 상단의 use 문 확인

   문제: 테스트 실패
   해결: 코드 블록 내 #[cfg(test)] 섹션 확인

기대하지 않는 것:
❌ 파일 시스템 접근
❌ 직접 파일 생성

기대하는 것:
✅ 마크다운 코드 블록으로 완전한 구현 코드
✅ 라인 번호 정확히 지정
✅ 복사-붙여넣기만으로 동작
✅ 검증 명령어 제시
```

---

### **Agent 3: OS Kernel (개선)**

```
Phase 8 Hub-Spoke Optimization: freelang-os-kernel 시스템콜 표준화

제약 조건: 권한 없음 → 설계문서 + 코드 제시 형식

기대 산출물:

1️⃣ ARCHITECTURE DOCUMENT (마크다운 + 코드 블록)

   ## 현재 상태
   - 파일: /data/data/com.termux/files/home/.projects/core/freelang-os-kernel/src/syscall.rs
   - 라인 수: 142줄
   - 현재 인터페이스:
     ```rust
     pub fn syscall(num: u64, arg1: u64, ..., arg6: u64) -> i64
     ```
   - 문제: 타입 불안전, 캐싱 불가

   ## 변경 후 상태

   새 구조:
   ```rust
   pub struct SyscallArgs {
       pub number: u64,
       pub args: [u64; 6],
   }

   pub struct SyscallResult {
       pub value: i64,
       pub error: i32,
   }
   ```

2️⃣ COMPLETE NEW FILES (마크다운 코드 블록)

   파일 1: syscall_types.rs (100줄)
   ```rust
   // 완전한 타입 정의
   ```

   파일 2: syscall_dispatcher.rs (150줄)
   ```rust
   // 최적화된 디스패처
   ```

3️⃣ MODIFICATION GUIDE

   수정 대상: syscall.rs (기존 파일)

   # 변경 1: 함수 시그니처 변경
   라인 50-60:

   현재:
   ```rust
   pub fn syscall_handler(num: u64, ...) -> i64 {
       match num {
           0 => sys_exit(...),
   ```

   변경 후:
   ```rust
   pub fn syscall_handler(args: &SyscallArgs) -> SyscallResult {
       match args.number {
           0 => sys_exit(&args.args[0]),
   ```

   # 변경 2: RDTSC 추가 (성능 측정)
   라인 15에 추가:
   ```rust
   #[inline]
   fn rdtsc() -> u64 {
       unsafe { std::arch::x86_64::_rdtsc() }
   }
   ```

4️⃣ TESTING CODE (마크다운 코드 블록)

   파일: tests/syscall_test.rs (새 파일)
   ```rust
   #[test]
   fn test_syscall_throughput() {
       // 30개 테스트
   }
   ```

5️⃣ VALIDATION

   컴파일 확인:
   cargo build -p freelang-os-kernel 2>&1

   테스트 실행:
   cargo test --test syscall_test

   성능 확인:
   time cargo run --bin kernel_bench --release

   예상:
   ✅ 컴파일 성공
   ✅ 30개 테스트 통과
   ✅ syscalls/sec: 1,000 → 2,000+

기대하는 것:
✅ 현재 코드 구체적 분석
✅ 정확한 라인 번호 지정
✅ 전체 코드 블록 제시
✅ 수정 전/후 명확히 표시
```

---

## 개선 체크리스트

### ✅ 에이전트가 해야 할 것

```
1. 파일 위치 분석
   ☑ 절대 경로 명시 (/data/data/com.termux/files/home/.projects/...)
   ☑ 현재 라인 수 확인
   ☑ 현재 구조 설명

2. 완전한 코드 제시
   ☑ 마크다운 코드 블록 (```rust, ```python 등)
   ☑ Copy-paste 가능한 형식
   ☑ 모든 import/use 문 포함
   ☑ 테스트 코드 포함

3. 통합 가이드
   ☑ 수정할 파일 명시
   ☑ 정확한 라인 번호
   ☑ 변경 전/후 코드 제시
   ☑ 단계별 수동 절차

4. 검증 방법
   ☑ 실행 가능한 명령어
   ☑ 예상 출력 명시
   ☑ 문제 해결 팁
   ☑ 성능 측정 방법
```

### ❌ 에이전트가 하지 말아야 할 것

```
☑ Write/Bash 도구로 파일 생성 시도
☑ 권한 부여 요청 (대신 해결책 제시)
☑ 애매한 코드 (명시적으로)
☑ "권한이 없습니다" 만 반복 (대신 대안 제시)
☑ 불완전한 코드 블록
```

---

## 개선된 에이전트 통합 프롬프트 (템플릿)

```
# Phase 8 최적화 작업 - [프로젝트 이름]

## 제약 조건
- 파일 쓰기 권한 없음
- Bash 직접 실행 불가
- 모든 산출물은 마크다운 텍스트 형식

## 기대 산출물 (이 순서로)

### 1️⃣ 현재 상태 분석 (구체적)
- [ ] 프로젝트 위치: [절대 경로]
- [ ] 현재 코드 라인 수: [숫자]
- [ ] 주요 파일: [파일명]
- [ ] 현재 병목점: [구체적 설명]

### 2️⃣ 완전한 구현 코드 (마크다운)
- [ ] 파일 1: [이름] ([줄수]줄)
  ```rust/python/fl
  [완전한 코드 - copy-paste 가능]
  ```

- [ ] 파일 2: [이름] ([줄수]줄)
  ```
  [완전한 코드]
  ```

### 3️⃣ 수정 가이드 (라인 번호 명시)
- [ ] 수정할 파일: [경로]
  - 변경 위치: Line [번호]-[번호]
  - 현재 코드: [코드 블록]
  - 변경 후: [코드 블록]

- [ ] 새 파일 생성:
  - 경로: [절대 경로]
  - 내용: [위의 코드 블록 복사]

### 4️⃣ 검증 방법 (명령어)
```bash
# Step 1: 컴파일 확인
[구체적 명령어]

# Step 2: 테스트 실행
[구체적 명령어]

# Step 3: 성능 측정
[구체적 명령어]
```

### 5️⃣ 예상 결과
- ✅ 컴파일: [예상 결과]
- ✅ 테스트: [몇 개 통과]
- ✅ 성능: [구체적 수치]

## 핵심 원칙
1. 모든 코드는 마크다운 코드 블록으로
2. 절대 경로 명시
3. Copy-paste 만으로 동작
4. 수동 구현 가능해야 함
5. 검증 방법 명확히
```

---

## 사용 방법

다음 7개 에이전트를 **개선된 프롬프트**로 재실행:

```bash
# Agent 1-7 모두 위 프롬프트 템플릿 사용
# "제약 조건" 섹션 포함 → 권한 문제 해결

Agent 1: Runtime 캐싱
Agent 2: Compiler 아티팩트
Agent 3: OS Kernel 인터페이스
Agent 4: Async/Distributed
Agent 5: 파이프라인 최적화
Agent 6: 클러스터 배포
Agent 7: 지표 표준화
```

---

## 결과물 형식

**모든 에이전트가 생성할 결과**:

```
# Phase 8: [작업 이름]

## 1️⃣ 현재 상태
[분석]

## 2️⃣ 구현 코드

### shared_memory.rs (250줄)
```rust
[완전한 코드]
```

### parallel_compiler.rs (200줄)
```rust
[완전한 코드]
```

## 3️⃣ 통합 가이드

Step 1: shared_memory.rs 생성
- 위 코드 블록 전체 복사
- 저장 위치: /path/to/src/shared_memory.rs

Step 2: lib.rs 수정
- Line 15: `mod shared_memory;` 추가
- [변경 전/후]

## 4️⃣ 검증
[명령어]

## 5️⃣ 예상 결과
[예상값]
```

---

**이 템플릿으로 Agent 재실행 시 100% 실행 가능한 결과물 생성됩니다** ✅
