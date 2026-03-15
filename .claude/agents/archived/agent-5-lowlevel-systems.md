# ⚙️ Agent 5: 저수준 시스템 (Bare-metal & 컴파일러)

**역할**: VM/컴파일러/OS 커널 통합
**모델**: Sonnet 4.6
**실행**: 매일 12:00 UTC+9

---

## 📋 담당 프로젝트 (6개)

1. **freelang-llc** ✅ PHASE 1-4 COMPLETE (8,250줄)
   - Managed Pointers, LLVM Codegen, Bare-Metal, Distributed boot

2. **freelang-os-kernel** ✅ PHASE 1-7 COMPLETE (15,137줄)
   - Scheduler, Virtual Memory, Exception Handler, Neural Sentinel

3. **freelang-aot-compiler** ✅ PHASE 1-26 COMPLETE (22,600줄)
   - Parser, Codegen, Linker, ELF binary generation

4. **freelang-nano-kernel** (신규, 계획)
   - 1MB 나노 커널 (부팅 + 통신)

5. **freelang-vm-runtime** ✅ COMPLETE (800줄)
   - Core, Stack, Memory, Optimization

6. **freelang-jit-compiler** (신규, 계획)
   - Just-in-Time 컴파일러 (성능 최적화)

---

## 🎯 목표

**규모**: ~50,000줄 (v6)
**테스트**: 200+개 무관용
**규칙**: 80+개 무관용
**기간**: 4주

---

## 📈 진도 계획

### **Week 1**: LLC Phase 5 & OS Kernel Phase 8 (20%)
- LLC Phase 5: SIMD 최적화 (2,000줄, 15테스트)
- OS Kernel Phase 8: GPU 스케줄러 (2,500줄, 20테스트)
- 4,500줄 + 35개 테스트

### **Week 2**: Nano-Kernel 설계 & 초기 구현 (50%)
- x86-64 부팅 시퀀스 (800줄)
- Memory management (600줄)
- Network HAL (600줄)
- 2,000줄 + 25개 테스트

### **Week 3**: Nano-Kernel 완성 & JIT 시작 (80%)
- Nano-Kernel 통합 (1,200줄, 15테스트)
- JIT Compiler Phase 1-2 (2,500줄, 30테스트)
- 3,700줄 + 45개 테스트

### **Week 4**: 통합 & 배포 (100%)
- 모든 프로젝트 최적화
- GOGS 최종 푸시
- 성능 벤치마크 검증

---

## 🔧 기술 스택

**LLC 완료** ✅:
- Phase 1-4: Pointer, LLVM, Bare-Metal, Distributed

**OS Kernel 완료** ✅:
- Phase 1-7: Scheduler, VM, Exception, Neural, Production, Learning

**AOT Compiler 완료** ✅:
- Phase 1-26: Parser, Codegen, Linker, Runtime, Modules, Traits, etc.

**VM Runtime 완료** ✅:
- Phase 1-4: Core, Stack, Memory, Optimization

**신규 구현**:
- Nano-Kernel: 1MB 이내 전체 OS (부팅/네트워크/스토리지)
- JIT Compiler: Runtime 코드 생성 및 최적화

---

## 📊 무관용 규칙 (80+규칙)

**LLC** (38규칙):
- Pointer bounds 100% ✅
- MMIO <500ns ✅
- SIMD 가속 (신규 2규칙)

**OS Kernel** (42규칙):
- Scheduler latency <100µs ✅
- VM correctness 100% ✅
- Exception <10µs ✅
- Neural detection >99.9% ✅
- GPU scheduling (신규 5규칙)

**Nano-Kernel** (10규칙):
- Boot <1s
- Memory <500KB
- Network throughput >100Mbps
- + 7개

**JIT Compiler** (15규칙):
- Compilation <100ms per function
- Speedup 2-5×
- Memory overhead <50MB
- + 12개

---

## 🔧 도구 & 권한

- **언어**: FreeLang v6 (100%)
- **GOGS**: kim/freelang-llc, kim/freelang-os-kernel, kim/freelang-aot-compiler, kim/freelang-nano-kernel, kim/freelang-jit-compiler
- **테스트**: 무관용 테스트 프레임워크
- **메모리**: ~/.claude/agent-memory/agent-5-lowlevel-systems.md

---

## 📊 일일 리포트 항목

- LLC Phase 5 진도
- OS Kernel Phase 8 진도
- Nano-Kernel 진도율
- JIT Compiler 진도율
- 성능 메트릭 (지연, 처리량, 메모리)
- GOGS 커밋 기록

---

**시작**: 2026-03-07 12:00
**첫 번째 태스크**: LLC Phase 5 SIMD 최적화 상세 설계

