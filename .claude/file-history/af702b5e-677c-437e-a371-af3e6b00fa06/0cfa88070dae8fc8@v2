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

## 🎯 Phase 1: Parser (현재)

**상태**: 🚀 진행 중

### 목표
- FreeLang 소스 코드를 AST로 변환
- 100% 호스트 언어 의존성 제거

### 구성
1. **lexer.fl** (400줄) - 토크나이저
2. **ast.fl** (200줄) - AST 노드 정의
3. **parser.fl** (200줄) - 재귀 하강 파서

### 테스트 (10개 무관용)
- ✓ 정수 리터럴 파싱
- ✓ 함수 선언 파싱
- ✓ 이진 연산 파싱
- ✓ 제어흐름 파싱
- ✓ 구조체 정의 파싱
- ✓ 타입 주석 파싱
- ✓ 에러 핸들링
- ✓ 주석 스킵
- ✓ 심볼 우선순위
- ✓ 중첩 표현식

### 성공 지표
```
TokenizeSuccess: 100%
ParseSuccess: ≥95%
ErrorRecovery: ≥90%
AST_Correctness: 100%
```

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

### Phase 1: Parser
- [x] 아키텍처 설계 (DESIGN.md)
- [x] Token 정의 (token.fl)
- [x] Lexer 구현 (lexer.fl)
- [x] AST 정의 (ast.fl)
- [x] Parser 구현 (parser.fl)
- [x] 무관용 테스트 설계 (test_lexer.fl)
- [ ] 테스트 실행 & 디버깅
- [ ] Phase 1 완료 커밋

### Phase 2-4
- [ ] Code Generation (x86-64)
- [ ] ELF Linker
- [ ] Bootstrap + Runtime

---

## 🔗 독립성 추적

| 단계 | 호스트 언어 의존성 | 언어 독립성 |
|------|------------------|-----------|
| **Phase 1** | Rust 일부 | ~60% |
| **Phase 2** | Rust 최소화 | ~80% |
| **Phase 3** | 거의 없음 | ~95% |
| **Phase 4** | **0%** | **100%** |

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
