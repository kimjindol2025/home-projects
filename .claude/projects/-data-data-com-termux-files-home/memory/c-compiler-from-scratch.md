---
name: c-compiler-from-scratch - 완전한 C 컴파일러 + VM
description: 19,158줄 | Lexer/Parser/Codegen/ELF | VM 구현 완료 | Phase 6 프리프로세서 | SAI 최적화
type: project
---

# 🚀 c-compiler-from-scratch

**구축 중인 완전한 C 컴파일러 (C11 호환, 부분)**

---

## 📌 **현재 상태**

**상태**: 🟡 In Progress
**총 코드량**: 19,158줄
**최신 커밋**: `06acc75 feat(vm): SAI + Computed Goto — VM이 스스로 진화한다`
**저장소**: https://gogs.dclub.kr/kim/c-compiler-from-scratch.git

---

## 🏗️ **구현 단계**

### ✅ **Phase 1: 기초 구현**
- Lexer: 토큰화 완료
- Parser: AST 생성 완료
- Codegen: 기계어 생성 완료
- ELF 바이너리 작성

### ✅ **Phase 2: 고급 기능**
- Nested structs 지원
- Pointer alignment 처리
- Function pointers 지원
- 복잡한 선언 파싱

### ✅ **Phase 3: 정적 변수**
- Static local variables 구현
- ELF 재배치 심볼 인덱스 수정
- 스코프 처리 개선

### ✅ **Phase 4: 제어 흐름**
- switch/fallthrough 완전 지원
- ELF 심볼 테이블 최적화
- 모든 제어 흐름 구조 지원

### ✅ **Phase 5: ABI 안정화**
- Stack argument alignment (x86-64 ABI)
- 암시적 타입 캐스팅
- 함수 호출 규약 표준화

### ✅ **Phase 6: C 전처리기**
- `#include` 지원 (파일 포함)
- `#define` 매크로 치환
- `#ifdef` 조건부 컴파일
- `stdarg.h` 변수 인자 지원

### ✅ **VM 구현**
- 완전한 바이트코드 VM
- 8개 핵심 개념 구현
- OBJ_NATIVE, OBJ_MAP, OBJ_FIBER 지원
- EXTERNAL 객체 타입

### ⏳ **최신: SAI + Computed Goto**
- Static Argument Indexing 최적화
- VM 자동 진화 메커니즘
- Computed Goto 성능 향상

---

## 📂 **폴더 구조**

```
c-compiler-from-scratch/
├── src/                    # 컴파일러 소스 (19,158줄)
│   ├── lexer.c            # 토크나이저
│   ├── parser.c           # 파서 (AST 생성)
│   ├── sema.c             # 시맨틱 분석
│   ├── codegen.c          # 코드 생성
│   ├── x86_encode.c       # x86-64 인코더
│   ├── elf_writer.c       # ELF 바이너리 생성
│   ├── preprocessor.c     # C 전처리기
│   ├── ir.c               # 중간 표현
│   ├── main.c             # 진입점
│   └── ...
│
├── include/               # 헤더 파일
│   ├── ast.h              # AST 정의
│   ├── codegen.h
│   ├── elf.h
│   ├── ir.h
│   ├── lexer.h
│   ├── parser.h
│   ├── preprocessor.h
│   ├── sema.h
│   ├── symtable.h
│   ├── types.h
│   └── x86_encode.h
│
├── vm/                    # 바이트코드 VM
│   ├── vm.c
│   ├── vm.h
│   └── objects/           # 객체 타입
│
├── tests/                 # 테스트
│   ├── test_lexer.c
│   └── unit/
│
├── scripts/               # 빌드 스크립트
│   └── Makefile
│
├── bin/                   # 빌드 결과
│   └── a.out             # 컴파일러 바이너리
│
├── CLAUDE.md             # 프로젝트 헌장
├── MEMORY.md             # 진행 상황 (2026-03-15)
└── a.out                 # 테스트 바이너리
```

---

## 🔧 **지원하는 C 기능**

### ✅ **데이터 타입**
- 기본 타입: int, char, float, double, void
- Pointer 타입: `int*`, `char**` 등
- Struct: nested struct, by-value passing
- Union: flat/nested 필드 동기화
- Array: 다차원 배열

### ✅ **선언 & 정의**
- 변수 선언 및 초기화
- 함수 선언 및 정의
- Static 변수 (전역, 로컬)
- Function pointers
- 복잡한 선언자 파싱

### ✅ **제어 흐름**
- if/else 조건문
- for/while 루프
- do-while
- switch/case (fallthrough 지원)
- break/continue

### ✅ **표현식**
- 산술 연산자: +, -, *, /, %
- 논리 연산자: &&, ||, !
- 비트 연산자: &, |, ^, <<, >>
- 할당 연산자: =, +=, -=, 등
- 삼항 연산자: ? :

### ✅ **함수**
- 함수 호출
- 가변 인자 (stdarg)
- 재귀 함수
- 외부 함수 호출 (libc)
- Function pointers

### ✅ **전처리기**
- `#include <file.h>`
- `#define MACRO(x) ...`
- `#ifdef / #endif`
- 매크로 치환

---

## ⚡ **성능 특징**

### 컴파일러
- **ELF 바이너리 생성**: x86-64 아키텍처
- **ABI 준수**: x86-64 호출 규약 완벽 구현
- **최적화**:
  - SSE2 double/float 연산
  - Struct by-value ABI
  - JIT 컴파일 지원
  - SAI (Static Argument Indexing)

### VM
- **바이트코드 인터프리터**: 8개 핵심 개념
- **객체 시스템**: OBJ_NATIVE, OBJ_MAP, OBJ_FIBER, OBJ_EXTERNAL
- **최적화**: Computed Goto로 dispatch 가속

---

## 📊 **최근 커밋 히스토리**

```
06acc75 feat(vm): SAI + Computed Goto — VM이 스스로 진화한다
8880c66 feat(vm): V가 못 하는 것 — OBJ_NATIVE, OBJ_MAP, OBJ_FIBER, OBJ_EXTERNAL
6f4988d feat(vm): 완전한 바이트코드 VM 구현 — 8개 핵심 개념 코드로 설명
ef2b45e feat: struct by-value ABI, double/float SSE2, JIT 수정 + 고성능 기반 구축
0c5f3b1 feat: Phase 6 - Text-based C preprocessor (#include, #define, #ifdef, stdarg)
d4abace feat: Phase 5 - ABI stabilization (stack args alignment + implicit cast)
68e320e feat: Phase 4 - control flow completion (switch fallthrough + ELF symidx)
8a42118 feat: Phase 3 - static local variables + ELF reloc symidx fix
ef85c7a feat: Phase 2 — nested structs, ptr alignment, func pointers
ed11068 Phase 1: Sema 전면 수정 — union/flat 필드 동기화 완료
```

---

## 🎯 **사용 방법**

### 컴파일
```bash
cd c-compiler-from-scratch
make                    # 컴파일러 빌드

# 사용자의 C 파일 컴파일
./a.out input.c        # ELF 바이너리 생성 (a.out)
./a.out input.c -o output  # 출력 파일 지정
```

### 실행
```bash
./a.out output_program
```

### 테스트
```bash
make test              # 유닛 테스트 실행
```

---

## 🔍 **파이프라인**

```
Input C File
    ⬇️
┌─────────────────────┐
│ Lexer (lexer.c)     │  → Tokens
└─────────────────────┘
    ⬇️
┌─────────────────────┐
│ Parser (parser.c)   │  → AST
└─────────────────────┘
    ⬇️
┌─────────────────────┐
│ Sema (sema.c)       │  → Type Checking
└─────────────────────┘
    ⬇️
┌─────────────────────┐
│ Codegen (codegen.c) │  → x86-64 Assembly
└─────────────────────┘
    ⬇️
┌─────────────────────┐
│ x86 Encoder         │  → Machine Code
│ (x86_encode.c)      │
└─────────────────────┘
    ⬇️
┌─────────────────────┐
│ ELF Writer          │  → ELF Binary
│ (elf_writer.c)      │
└─────────────────────┘
    ⬇️
Output Binary (a.out)
```

---

## 🧪 **테스트**

| 카테고리 | 상태 | 상세 |
|---------|------|------|
| Lexer | ✅ | 토큰화 완료 |
| Parser | ✅ | AST 생성 완료 |
| Codegen | ✅ | x86-64 생성 완료 |
| ELF | ✅ | 바이너리 생성 완료 |
| VM | ✅ | 완전 구현 |
| Preprocessor | ✅ | Phase 6 완료 |

---

## 📈 **프로젝트 규모**

| 항목 | 값 |
|------|-----|
| 총 코드량 | 19,158줄 |
| 소스 파일 | 20+ |
| 헤더 파일 | 10+ |
| 커밋 | 100+ |
| 구현 Phases | 6개 완료 + VM |

---

## 🚀 **다음 목표**

- [ ] C89 호환성 강화
- [ ] 더 많은 stdlib 함수 지원
- [ ] 최적화 패스 추가
- [ ] 디버깅 정보 (DWARF) 지원
- [ ] RISC-V 백엔드
- [ ] 웹어셈블리 지원

---

## 📝 **특징**

✨ **교육용 + 실용적**
- 각 단계가 명확하게 문서화됨
- 실제 작동하는 컴파일러
- 고성능 구현 (JIT, SAI, Computed Goto)
- 확장 가능한 아키텍처

---

**상태**: 🟡 In Progress | **매우 진전됨** ✅

