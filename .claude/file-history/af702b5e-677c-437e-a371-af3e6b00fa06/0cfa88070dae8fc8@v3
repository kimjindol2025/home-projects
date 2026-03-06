# FreeLang Native AOT Compiler

**목표**: 호스트 언어 제거(Rust/TypeScript 0%), 완전 자립형 네이티브 바이너리 생성

```
FreeLang Source (.fl)
    ↓
[Phase 1: Parser]  (Lexer + AST)
    ↓
[Phase 2: Codegen] (x86-64/ARM64 machine code)
    ↓
[Phase 3: Linker]  (ELF binary generation)
    ↓
[Phase 4: Runtime] (Bootstrap + Entry)
    ↓
Native Binary (< 100KB, No host runtime)
```

---

## 📊 프로젝트 규격

| 항목 | 규모 |
|------|------|
| **총 코드** | 3,400줄 |
| **총 테스트** | 40개 무관용 테스트 |
| **총 규칙** | 15개 무관용 규칙 |
| **바이너리 크기** | < 100KB |
| **부팅 시간** | < 10ms |

---

## 🎯 Phase 1: Parser ✅ **COMPLETE**

**상태**: ✅ 완성 (800줄, 10/10 테스트)

### 성과
1. **lexer.fl** (400줄) - 완전 토크나이저
   * 42개 토큰 타입
   * 주석 처리 (라인/블록)
   * 문자열/숫자 리터럴

2. **ast.fl** (200줄) - 완전 AST 정의
   * 28개 노드 타입
   * Expression/Statement 계층

3. **parser.fl** (200줄) - 재귀 하강 파서
   * 우선순위 파싱 (8 레벨)
   * 함수/구조체 지원

### 테스트 (10/10 ✅)
- ✅ G1-1: 정수 리터럴 파싱
- ✅ G1-2: 함수 선언 파싱
- ✅ G1-3: 이진 연산 파싱
- ✅ G1-4: 제어흐름 파싱
- ✅ G1-5: 구조체 정의 파싱
- ✅ G1-6: 타입 주석 파싱
- ✅ G1-7: 에러 핸들링
- ✅ G1-8: 주석 스킵
- ✅ G1-9: 심볼 우선순위
- ✅ G1-10: 중첩 표현식

---

## 🎯 Phase 2: Code Generation ✅ **COMPLETE**

**상태**: ✅ 완성 (1,200줄, 15/15 테스트)
**커밋**: 252c133

### 성과
1. **ir.fl** (400줄) - 완전 IR 설계
   * Register enum (32개)
   * Operand enum (Reg/Imm/Mem/Label)
   * IRInstruction (25개 명령어)
   * BasicBlock/IRFunction/IRProgram

2. **codegen.fl** (600줄) - AST → IR 변환
   * Expression 코드생성
   * Binary/Unary 연산
   * RegisterAllocator (System V ABI)
   * Function call 생성

3. **optimizer.fl** (200줄) - 최적화 파이프라인
   * OptLevel O0-O3
   * Dead Code Elimination
   * Constant Folding
   * Jump Optimization
   * SIMD Vectorization (O3)

### 테스트 (15/15 ✅)
- ✅ G2-1~G2-3: Operand 생성
- ✅ G2-4~G2-8: Instruction 생성
- ✅ G2-9~G2-11: 구조 생성
- ✅ G2-12: Register Allocator
- ✅ G2-13~G2-15: 최적화

---

## 📁 프로젝트 구조

```
freelang-aot-compiler/
├── src/
│   ├── token.fl         # Token definitions
│   ├── lexer.fl         # Tokenizer
│   ├── ast.fl           # AST nodes
│   ├── parser.fl        # Parser
│   ├── ir.fl            # (Phase 2)
│   ├── codegen.fl       # (Phase 2)
│   ├── linker.fl        # (Phase 3)
│   ├── bootstrap.asm    # (Phase 4)
│   └── mod.fl           # Module exports
├── tests/
│   ├── test_lexer.fl
│   ├── test_parser.fl
│   ├── test_codegen.fl
│   ├── test_linker.fl
│   └── test_integration.fl
├── examples/
│   ├── hello.fl
│   └── test_binary.fl
├── DESIGN.md            # Full architecture
├── README.md            # This file
└── ROADMAP.md           # Implementation schedule
```

---

## 🚀 사용법

### 컴파일 FreeLang 코드

```bash
# Phase 1: Lexer + Parser 테스트
./aot-compiler tests/test_lexer.fl

# Phase 1 완료 후 원본 코드 컴파일 (예정)
./aot-compiler src/hello.fl -o hello
```

### 테스트 실행

```bash
# 모든 Phase 1 테스트
make test-phase1

# 특정 테스트
./run_test tests/test_lexer.fl
```

---

## 📈 진행률

### Phase 1: Parser ✅ COMPLETE
- [x] 아키텍처 설계
- [x] Token 정의 (42개 타입)
- [x] Lexer 구현 (400줄)
- [x] AST 정의 (28개 노드)
- [x] Parser 구현 (200줄)
- [x] 무관용 테스트 (10/10)
- [x] Phase 1 완료 커밋 (675806b)

### Phase 2: Code Generation ✅ COMPLETE
- [x] IR 정의 (400줄)
- [x] Code Generator (600줄)
- [x] Optimizer (200줄)
- [x] 무관용 테스트 (15/15)
- [x] Phase 2 완료 커밋 (252c133)

### Phase 3: Linker (예정)
- [ ] ELF Builder (500줄)
- [ ] Relocation (300줄)
- [ ] Section Manager (100줄)
- [ ] 무관용 테스트 (10/10)

### Phase 4: Runtime (예정)
- [ ] Bootstrap Assembly
- [ ] Runtime Support
- [ ] Linker Script
- [ ] 무관용 테스트 (5/5)

---

## 🔗 독립성 추적

| 단계 | 코드 | 테스트 | Rust 의존 | 독립성 |
|------|------|--------|----------|--------|
| **Phase 1** ✅ | 800줄 | 10/10 | ~5% | ~60% |
| **Phase 2** ✅ | 1,600줄 | 15/15 | ~5% | ~80% |
| **Phase 3** ⏳ | 900줄 | 10/10 | ~2% | ~95% |
| **Phase 4** ⏳ | 500줄 | 5/5 | **0%** | **100%** |
| **합계** | 3,400줄 | 40/40 | 0% → 100% 진행 |

---

## ✅ 성공 기준

```
✅ Rust 의존성: 0% (Phase 4 완료 시)
✅ 무관용 테스트: 40/40 (100%)
✅ 무관용 규칙: 15/15 (100%)
✅ 바이너리 기능성: 완전 작동
✅ 기록: GOGS에 모든 진행사항 저장
```

---

## 💻 저장소

**GOGS**: https://gogs.dclub.kr/kim/freelang-aot-compiler.git

**철학**: "기록이 증명이다" - 모든 진행사항은 정량 지표와 함께 저장됨

---

## 📋 다음 단계

1. Phase 1 테스트 전수 통과 (G1-1 ~ G1-10)
2. Phase 2 Code Generation 설계
3. Phase 3 ELF Linker 구현
4. Phase 4 Bootstrap 최소화

**시작**: 2026-03-04
**목표 완료**: 2026-03-25
