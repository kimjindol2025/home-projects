# FreeLang Native AOT Compiler - 아키텍처 설계

**목표**: 호스트 언어 제로(Rust/TypeScript 완전 제거), 완전 자립형 네이티브 바이너리 생성

---

## 📊 프로젝트 개요

```
┌─────────────────────────────────────────┐
│   FreeLang Source Code (.fl files)      │
└──────────────────┬──────────────────────┘
                   ↓
        ┌──────────────────────┐
        │  PHASE 1: Parser     │ (800줄)
        │  - Lexer             │
        │  - AST Builder       │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │  PHASE 2: Codegen    │ (1,200줄)
        │  - IR Generation     │
        │  - x86-64/ARM64      │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │  PHASE 3: Linker     │ (900줄)
        │  - ELF Sections      │
        │  - Relocation        │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │  PHASE 4: Runtime    │ (500줄)
        │  - Bootstrap         │
        │  - Entry Point       │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │  Native Binary       │
        │  (ELF/PE, <100KB)    │
        │  No Host Runtime!    │
        └──────────────────────┘
```

---

## 🎯 4-Phase 구현 계획

### **Phase 1: Parser (800줄, 10개 무관용 테스트)**

#### 목표
- FreeLang 소스 코드를 추상 구문 트리(AST)로 변환
- 호스트 언어 의존성 없이 순수 FreeLang으로 작성

#### 구성
1. **lexer.fl** (400줄)
   - 토크나이저: `fn`, `let`, `const`, `struct`, `fn`, `return` 등
   - 심볼 인식: `+`, `-`, `*`, `/`, `(`, `)`, `{`, `}` 등
   - 문자열/숫자 리터럴 처리
   - 에러 핸들링 (라인 번호, 위치 정보)

2. **ast.fl** (200줄)
   - AST 노드 정의:
     ```
     enum Token
     enum ASTNode
     struct Program
     struct Function
     struct Statement
     struct Expression
     ```

3. **parser.fl** (200줄)
   - 재귀 하강 파서 (RDP)
   - 우선순위 파싱 (precedence climbing)
   - 에러 복구

#### 무관용 테스트 (10개)
- G1-1: 정수 리터럴 파싱
- G1-2: 함수 선언 파싱
- G1-3: 이진 연산 파싱
- G1-4: 제어흐름 (if/for) 파싱
- G1-5: 구조체 정의 파싱
- G1-6: 타입 주석 파싱
- G1-7: 에러 핸들링 (잘못된 구문)
- G1-8: 주석 스킵
- G1-9: 심볼 우선순위
- G1-10: 중첩 표현식 파싱

#### 성공 지표
```
TokenizeSuccess: 100%
ParseSuccess: ≥95%
ErrorRecovery: ≥90%
AST_Correctness: 100%
```

---

### **Phase 2: Code Generation (1,200줄, 15개 무관용 테스트)**

#### 목표
- AST를 x86-64 또는 ARM64 머신 코드로 변환
- 최적화 패스 (O0-O3)

#### 구성
1. **ir.fl** (400줄)
   - Intermediate Representation 정의
   - Basic blocks, SSA form
   - 명령어 형식

2. **codegen_x86_64.fl** (600줄)
   - x86-64 명령어 생성
   - 레지스터 할당
   - 호출 규약 (System V AMD64 ABI)

3. **optimizer.fl** (200줄)
   - 상수 폴딩
   - 데드 코드 제거
   - 루프 언롤링

#### 무관용 테스트 (15개)
- G2-1 ~ G2-7: 기본 연산 (add, sub, mul, div, mod)
- G2-8 ~ G2-11: 함수 호출 & 반환
- G2-12 ~ G2-15: 최적화 검증

#### 성공 지표
```
CodeGenSuccess: 100%
Instruction_Count_Efficiency: ≥90%
Register_Allocation: Optimal
Binary_Size_x86_64: <50KB
```

---

### **Phase 3: Linker (900줄, 10개 무관용 테스트)**

#### 목표
- 목적 파일(.o)을 실행 파일(ELF)로 링크
- 재배치(relocation) 처리
- 섹션 배치

#### 구성
1. **elf_builder.fl** (500줄)
   - ELF 헤더 생성
   - 섹션 관리 (.text, .data, .rodata, .bss)
   - 심볼 테이블

2. **relocation.fl** (300줄)
   - 상대 주소 계산
   - 심볼 해석
   - 재배치 적용

3. **section_manager.fl** (100줄)
   - 메모리 레이아웃 관리

#### 무관용 테스트 (10개)
- G3-1 ~ G3-5: ELF 구조 검증
- G3-6 ~ G3-8: 재배치 정확성
- G3-9 ~ G3-10: 링크 결과 검증

#### 성공 지표
```
ELF_Validity: 100%
Relocation_Accuracy: 100%
Entry_Point_Correct: 100%
Binary_Size: <100KB
```

---

### **Phase 4: Runtime (500줄, 5개 무관용 테스트)**

#### 목표
- 최소 부트스트랩 코드
- 진입점 관리
- 호출 스택 초기화

#### 구성
1. **bootstrap.asm** (200줄)
   - _start 진입점
   - 스택 초기화
   - main() 호출

2. **runtime_support.fl** (200줄)
   - 메모리 할당 (brk)
   - 시스템 콜 래퍼
   - 종료 처리

3. **linker_script.ld** (100줄)
   - GNU 링커 스크립트
   - 메모리 레이아웃

#### 무관용 테스트 (5개)
- G4-1: 부트스트랩 실행
- G4-2: main() 호출 확인
- G4-3: 반환값 처리
- G4-4: 스택 정렬
- G4-5: 종료 코드 전달

---

## 📈 전체 통계

| 항목 | 규모 |
|------|------|
| **총 코드** | 3,400줄 |
| **총 테스트** | 40개 무관용 |
| **총 규칙** | 15개 |
| **바이너리 크기** | <100KB |
| **부팅 시간** | <10ms |

---

## 🔗 의존성 제거 전략

| 호스트 언어 | 제거 방법 | Phase |
|-----------|---------|-------|
| Rust RuntimeCompiler | → FreeLang Lexer/Parser | 1 |
| LLVM Backend | → 네이티브 Codegen | 2 |
| Gold/ld Linker | → ELF Builder | 3 |
| libc | → syscall 직접 호출 | 4 |

**최종 상태**: 완전 자립형, 호스트 언어 제로

---

## ✅ 성공 기준

```
✅ 호스트 언어 의존성: 0% (Rust/TypeScript 제거)
✅ 언어 독립성: 100% (모든 코드 FreeLang)
✅ 무관용 테스트: 40/40 (100%)
✅ 무관용 규칙: 15/15 (100%)
✅ 바이너리 기능성: 완전 작동
```

---

## 📋 커밋 전략

```
Phase 1 완료 → Commit: "feat: Lexer + AST 구현 (10/10 테스트 통과)"
Phase 2 완료 → Commit: "feat: x86-64 Codegen + 최적화 (15/15 테스트)"
Phase 3 완료 → Commit: "feat: ELF Linker + 재배치 (10/10 테스트)"
Phase 4 완료 → Commit: "feat: 부트스트랩 + 런타임 (5/5 테스트)"
```

---

**"기록이 증명이다" - 모든 진행사항은 GOGS에 저장됨**
