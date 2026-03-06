# Z-Lang Transpiler - Phase 4: Backend Code Generation Plan

**프로젝트**: FreeLang to Z-Lang Transpiler
**Phase**: 4 (백엔드 코드 생성)
**기간**: 2026-03-06 ~ 2026-03-10
**규모**: 1,400줄 + 25개 테스트

---

## 📋 Phase 4 목표

### 목표 1: Code Generation (400줄)
- [ ] **code_generator.fl** (400줄)
  - IR 명령어 → x86-64 어셈블리 변환
  - 함수 호출 규약 (System V AMD64 ABI)
  - 임시 변수 및 스택 할당
  - 기본 명령어 생성 (mov, add, sub, mul, div, cmp, jmp)

### 목표 2: Assembly Builder (300줄)
- [ ] **assembly_builder.fl** (300줄)
  - 어셈블리 명령어 생성 및 최적화
  - 레이블 및 섹션 관리
  - 인스턴션 인코딩
  - 피킹 및 알레이먼트

### 목표 3: Register Management (200줄)
- [ ] **register_management.fl** (200줄)
  - x86-64 레지스터 할당
  - Caller-saved/Callee-saved 규칙
  - 임시 변수 관리
  - 컨벤션 준수

### 목표 4: Stack Frame Management (150줄)
- [ ] **stack_frame.fl** (150줄)
  - 스택 프레임 구성
  - Prologue/Epilogue 생성
  - 로컬 변수 할당
  - 인자 전달 및 반환값

### 목표 5: Linker (150줄)
- [ ] **linker.fl** (150줄)
  - 객체 파일 링킹
  - 심볼 테이블 관리
  - 리로케이션 처리
  - 라이브러리 링킹

### 목표 6: Binary Writer (200줄)
- [ ] **binary_writer.fl** (200줄)
  - ELF 파일 형식 생성
  - 섹션 및 헤더 작성
  - 심볼 및 리로케이션 테이블
  - 실행 가능한 바이너리 생성

### 목표 7: Integration & Testing (150줄)
- [ ] **phase4_integration.fl** (150줄)
  - 전체 백엔드 파이프라인 통합
  - IR → Assembly → Binary 변환

---

## 🧪 Test Plan (25개 테스트)

### Group B1: Code Generation (5개)
- [ ] **B1-T1**: 기본 명령어 - mov, add, sub
- [ ] **B1-T2**: 함수 호출 - CALL, RET
- [ ] **B1-T3**: 제어흐름 - JMP, CMP, JNE, JE
- [ ] **B1-T4**: 메모리 접근 - MOV [mem], LOAD, STORE
- [ ] **B1-T5**: 레지스터 할당 - RAX, RBX, RCX 할당

### Group B2: Assembly Building (5개)
- [ ] **B2-T1**: 명령어 인코딩 - 정확한 바이트 코드
- [ ] **B2-T2**: 레이블 해결 - 주소 계산
- [ ] **B2-T3**: 섹션 관리 - .text, .data, .bss
- [ ] **B2-T4**: 얼라인먼트 - 16바이트 정렬
- [ ] **B2-T5**: 최적화 - NOP 제거, JMP 축약

### Group B3: Register Management (5개)
- [ ] **B3-T1**: Caller-saved - RAX, RCX, RDX, RSI, RDI, R8-R11
- [ ] **B3-T2**: Callee-saved - RBX, RSP, RBP, R12-R15
- [ ] **B3-T3**: 임시 변수 - 자동 할당
- [ ] **B3-T4**: 스필 처리 - 레지스터 부족 시 메모리 사용
- [ ] **B3-T5**: 컨벤션 준수 - System V AMD64 ABI

### Group B4: Stack Frame (5개)
- [ ] **B4-T1**: Prologue - PUSH RBP, MOV RSP, RBP, SUB RSP
- [ ] **B4-T2**: Epilogue - MOV RSP, RBP, POP RBP, RET
- [ ] **B4-T3**: 로컬 변수 - [RBP-8], [RBP-16] 할당
- [ ] **B4-T4**: 인자 전달 - RDI, RSI, RDX, RCX, R8, R9
- [ ] **B4-T5**: 반환값 - RAX, RDX (64비트 이상)

### Group B5: Integration E2E (5개)
- [ ] **B5-T1**: 단순 함수 - int add(int a, int b)
- [ ] **B5-T2**: 재귀 함수 - factorial(n)
- [ ] **B5-T3**: 여러 인자 - func(a, b, c, d, e, f)
- [ ] **B5-T4**: 루프 코드 - for 루프 생성
- [ ] **B5-T5**: ELF 바이너리 - 실행 가능한 바이너리 생성

---

## 🎯 무관용 규칙 (Unforgiving Rules)

### R1: Code Correctness
- **정의**: 생성된 코드의 의미론 정확성 보장
- **검증**: IR 실행 결과 = 어셈블리 실행 결과
- **테스트**: B1-T1, B5-T1, B5-T2

### R2: Calling Convention
- **정의**: System V AMD64 ABI 완벽 준수
- **검증**: 인자 전달, 반환값, 레지스터 보존
- **테스트**: B3-T1, B3-T2, B4-T4, B4-T5

### R3: Stack Safety
- **정의**: 스택 오버플로우, 언더플로우 없음
- **검증**: RSP 유효성 검사, 스택 얼라인먼트 16바이트
- **테스트**: B4-T1, B4-T2, B4-T3

### R4: Register Efficiency
- **정의**: 레지스터 할당 최적성 (> 80% 사용률)
- **검증**: 스필 횟수 최소화
- **테스트**: B3-T3, B3-T4, B3-T5

### R5: Binary Compatibility
- **정의**: 생성된 바이너리가 Linux x86-64에서 실행 가능
- **검증**: ELF 형식, 실행 권한, 링킹
- **테스트**: B5-T5

---

## 📊 파일 구조

```
src/
  ├─ code_generator.fl          (400줄) ← NEW
  ├─ assembly_builder.fl        (300줄) ← NEW
  ├─ register_management.fl     (200줄) ← NEW
  ├─ stack_frame.fl             (150줄) ← NEW
  ├─ linker.fl                  (150줄) ← NEW
  ├─ binary_writer.fl           (200줄) ← NEW
  ├─ phase4_integration.fl      (150줄) ← NEW
  ├─ mod.fl                     (업데이트)
  ├─ optimizer_integration.fl   (Phase 3)
  ├─ loop_optimizer.fl          (Phase 3)
  ├─ constant_folder.fl         (Phase 3)
  ├─ dead_code_eliminator.fl    (Phase 3)
  ├─ register_allocator.fl      (Phase 3)
  ├─ parser.fl                  (Phase 2)
  ├─ ir.fl                      (Phase 2)
  └─ codegen.fl                 (Phase 2)

tests/
  ├─ phase4_tests.fl            (400줄, 25개 테스트) ← NEW
  ├─ phase4_unforgiving.fl      (120줄, 5개 규칙)   ← NEW
  └─ phase3_tests.fl            (Phase 3)

docs/
  └─ PHASE_4_CODEGEN_REPORT.md  (최종 보고서)
```

---

## 🚀 구현 전략

### 1단계: Code Generator (Day 1)
- IR 명령어 타입 분류
- x86-64 기본 명령어 매핑
- 함수 호출 패턴 구현

### 2단계: Assembly Builder (Day 1-2)
- 어셈블리 명령어 생성
- 레이블 해결
- 섹션 관리

### 3단계: Register Management (Day 2)
- 레지스터 할당 정책
- Caller/Callee-saved 규칙
- 임시 변수 관리

### 4단계: Stack Frame (Day 2-3)
- Prologue/Epilogue 생성
- 로컬 변수 할당
- 인자 전달 처리

### 5단계: Linker & Binary (Day 3)
- ELF 파일 형식
- 심볼 테이블
- 리로케이션 처리

### 6단계: Integration & Test (Day 3-4)
- 파이프라인 통합
- 25개 테스트 작성
- 최종 검증

---

## 📈 성능 목표

| 목표 | 측정 방법 |
|------|---------|
| **Code Size** | 생성 어셈블리 크기 최소화 (< 1.5× IR 크기) |
| **Register Usage** | 스필 최소화 (> 80% 사용률) |
| **Execution Speed** | 생성 코드 실행 시간 (인터프리터 대비 50×+) |
| **Binary Size** | ELF 바이너리 크기 (< 10KB 스켈레톤) |
| **Link Time** | 링킹 시간 (< 100ms) |

---

## 🎯 검증 기준

### ✅ 완료 조건
- [ ] 1,400줄 코드 완성
- [ ] 25개 테스트 100% 통과
- [ ] 5개 규칙 100% 검증
- [ ] x86-64 ABI 준수
- [ ] 실행 가능한 바이너리 생성
- [ ] GOGS 푸시 완료

---

**시작일**: 2026-03-06
**예상 완료**: 2026-03-10
**상태**: 📋 계획 수립 완료

