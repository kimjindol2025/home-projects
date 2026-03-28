---
name: FreeLang에 Julia 컴파일러 이식 계획
description: Julia 컴파일러(Go)를 FreeLang으로 재구현하기 위한 전략 및 실행 계획
type: project
---

# 🚀 FreeLang에 Julia 컴파일러 이식 계획

**시작 날짜**: 2026-03-19
**목표 완료**: 2026-04-30 (6주)
**상태**: 🟢 계획 수립 중

---

## 📊 **현황 분석**

### Julia Compiler (Go) 현재 상태
| 항목 | 상태 |
|------|------|
| **언어** | Go 1.18+ |
| **코드량** | ~3,000줄 |
| **테스트** | 30+ (E2E + 벤치마크) |
| **아키텍처** | 8단계 파이프라인 |
| **버전** | v0.2.0 (Production Ready) |
| **커밋** | 79723259 |

**8단계 파이프라인**:
1. Lexing (토크나이제이션)
2. Parsing (AST 생성)
3. Semantic Analysis (의미 분석)
4. IR Generation (중간 표현)
5. Type Inference (타입 추론)
6. Optimization (최적화)
7. Code Generation (바이트코드 생성)
8. VM Execution (실행)

---

### FreeLang 현재 상태
| 항목 | 상태 |
|------|------|
| **언어** | FreeLang (자체 언어) |
| **표준 라이브러리** | Collections, String, IO, Math |
| **컴파일 대상** | C 언어 (프로덕션) |
| **자체호스팅** | ✅ 증명됨 (freelang-to-c) |
| **프로젝트** | freelang-to-c, freelang-library-extraction |

---

## 🎯 **이식 전략**

### Phase 1: 설계 & 분석 (Week 1)
**목표**: 이식 가능성 판단 및 설계 문서 작성

#### Task 1.1: Julia 컴파일러 구조 분석
```
목표: Go 코드의 모듈 구조를 FreeLang으로 매핑
산출물:
  - module-mapping.md (Go 모듈 → FreeLang 모듈)
  - interface-design.md (공개 인터페이스 설계)
  - type-system-mapping.md (Go 타입 → FreeLang 타입)

예시:
  Go: type Lexer struct { source string }
  FL: record Lexer { source: String }
```

#### Task 1.2: FreeLang 언어 기능 평가
```
목표: Julia 이식에 필요한 FreeLang 기능 확인
검증 항목:
  ✅ 배열/맵 (Go slice/map → FreeLang Array/Map)
  ✅ 구조체 (Go struct → FreeLang record)
  ✅ 에러 처리 (Go error → FreeLang Result/Option)
  ✅ 패턴 매칭 (Go switch → FreeLang match)
  ✅ 함수형 기능 (map/filter/fold 지원?)

가능성: 90% (대부분 FreeLang 기능으로 구현 가능)
```

#### Task 1.3: 이식 로드맵 작성
```
산출물: PORTING_ROADMAP.md
구성:
  - Phase별 마일스톤
  - 모듈별 우선순위
  - 리스크 및 대응 방안
  - 테스트 전략
```

---

### Phase 2: 핵심 모듈 이식 (Week 2-3)
**목표**: 가장 의존성이 낮은 모듈부터 이식

#### Task 2.1: Lexer 이식
```
Go 파일: internal/lexer/lexer.go (350줄)
FL 파일: freelang-julia/src/lexer.fl (420줄 예상)

구현:
  1. Token 타입 정의 (record)
  2. Lexer 구조체 정의 (record)
  3. Tokenize() 메서드 구현
  4. 도우미 함수 (isAlpha, isDigit 등)

테스트:
  - 50개 토큰 타입 검증
  - 18개 테스트 케이스 이식
```

#### Task 2.2: Parser 이식
```
Go 파일: internal/parser/parser.go (450줄)
FL 파일: freelang-julia/src/parser.fl (550줄 예상)

구현:
  1. AST 노드 정의 (record 타입)
  2. Parser 구조체 정의
  3. 재귀 하강 파서 구현
  4. Precedence climbing (연산자 우선순위)

의존성: Lexer 완료 필요
```

#### Task 2.3: Type System 이식
```
Go 파일: internal/types/ (200줄)
FL 파일: freelang-julia/src/types.fl (280줄 예상)

구현:
  1. 기본 타입 정의 (i32, i64, f64, string, bool)
  2. 복합 타입 (Array, Function, Union)
  3. 타입 비교 & 호환성 검사

테스트: 14개 타입 검증 케이스
```

---

### Phase 3: 분석 & 생성 (Week 4-5)
**목표**: 중간 단계 모듈 이식

#### Task 3.1: Semantic Analyzer 이식
```
Go 파일: internal/sema/sema.go (500줄)
FL 파일: freelang-julia/src/sema.fl (620줄 예상)

구현:
  1. Scope & Symbol Table
  2. Type Checking
  3. Error Reporting

테스트: 20개 의미 분석 케이스
```

#### Task 3.2: IR Builder 이식
```
Go 파일: internal/ir/builder.go (250줄)
FL 파일: freelang-julia/src/ir_builder.fl (320줄 예상)

구현:
  1. IR 노드 정의
  2. AST → IR 변환
  3. BasicBlock 생성 및 연결

테스트: IR 생성 검증 (15개 케이스)
```

#### Task 3.3: Code Generator 이식
```
Go 파일: internal/codegen/codegen.go (400줄)
FL 파일: freelang-julia/src/codegen.fl (500줄 예상)

구현:
  1. IR → Bytecode 변환
  2. Instruction 생성
  3. 바이트코드 시리얼화

테스트: 코드 생성 검증 (12개 케이스)
```

---

### Phase 4: 최적화 & 통합 (Week 6)
**목표**: 전체 파이프라인 완성 및 테스트

#### Task 4.1: 통합 테스트
```
목표: E2E 파이프라인 검증

테스트 케이스:
  1. simple_literal: "42" → 42
  2. arithmetic: "2 + 3" → 5
  3. nested: "2 * 3 + 4" → 10
  4. complex: 함수 호출, 루프 등

각 단계 검증:
  - Lexer: 토큰 생성 확인
  - Parser: AST 구조 확인
  - Sema: 타입 검사 통과
  - IR: 중간 표현 생성
  - Codegen: 바이트코드 생성
  - VM: 실행 결과 검증
```

#### Task 4.2: 성능 최적화
```
목표: Go 버전 수준의 성능 달성

최적화 대상:
  1. Lexer: 토큰 버퍼링
  2. Parser: AST 캐싱
  3. IR: 불변성 활용
  4. Codegen: 바이트코드 최적화

벤치마크:
  - BenchmarkLexer: <100μs
  - BenchmarkParser: <200μs
  - BenchmarkFullPipeline: <1ms
```

#### Task 4.3: 문서화
```
산출물:
  - IMPLEMENTATION.md (구현 상세)
  - BUILD.md (빌드 가이드)
  - EXAMPLES.md (사용 예제)
  - API.md (공개 API)
```

---

## 📋 **구현 세부사항**

### 모듈 매핑 (Go → FreeLang)

| Go 모듈 | FreeLang 모듈 | 줄 수 | 의존성 |
|--------|--------------|------|--------|
| lexer | lexer.fl | 420 | 없음 |
| parser | parser.fl | 550 | lexer |
| types | types.fl | 280 | 없음 |
| sema | sema.fl | 620 | types, parser |
| ir | ir_builder.fl | 320 | parser |
| ir | ir.fl | 200 | 없음 |
| codegen | codegen.fl | 500 | ir |
| optimizer | optimizer.fl | 300 | ir |
| runtime | vm.fl | 400 | codegen |

**총 예상 코드량**: 3,590줄 (Go 버전 3,000줄과 유사)

---

## 🔄 **테스트 전략**

### Unit Tests (각 모듈별)
```
lexer_test.fl       (18 tests)
parser_test.fl      (14 tests)
types_test.fl       (12 tests)
sema_test.fl        (20 tests)
ir_test.fl          (15 tests)
codegen_test.fl     (12 tests)

총 91개 테스트
```

### Integration Tests (E2E)
```
e2e_test.fl (3 기본 + 5 고급)
  - simple_literal
  - arithmetic_operation
  - nested_arithmetic
  - function_call
  - loop_execution
  - array_operation
  - map_operation
  - string_manipulation
```

### Performance Tests
```
benchmark.fl
  - BenchmarkLexer
  - BenchmarkParser
  - BenchmarkFullPipeline
  - BenchmarkVM
```

---

## 🎯 **성공 기준**

| 지표 | 목표 |
|------|------|
| **기능성** | 8단계 파이프라인 100% 이식 |
| **호환성** | Go 버전과 동일한 출력 |
| **성능** | Go 버전 대비 95% 이상 |
| **테스트** | 90+ 테스트 모두 통과 |
| **코드품질** | Code Review 피드백 통합 |
| **문서화** | BUILD.md, API.md 완성 |

---

## 🚧 **리스크 및 대응**

| 리스크 | 심각도 | 대응 방안 |
|--------|--------|----------|
| FreeLang 기능 부족 | 🟠 중 | fallback to Go 구현 (hybrid mode) |
| 성능 저하 | 🟠 중 | 프로파일링 & 최적화 (loop unrolling, caching) |
| 복잡한 타입 변환 | 🟡 낮음 | 명시적 변환 함수 제공 |
| 메모리 효율 | 🟡 낮음 | 구조체 최적화, 불변성 활용 |

---

## 📅 **타임라인**

```
2026-03-19 (Week 1)
├─ Task 1.1: Julia 구조 분석 ✅
├─ Task 1.2: FreeLang 기능 평가
└─ Task 1.3: 로드맵 작성

2026-03-26 (Week 2)
├─ Task 2.1: Lexer 이식 (50% 완료)
├─ Task 2.2: Parser 이식 (시작)
└─ Task 2.3: Type System 이식 (준비)

2026-04-02 (Week 3)
├─ Task 2.1: Lexer 이식 (완료 + 테스트)
├─ Task 2.2: Parser 이식 (완료 + 테스트)
└─ Task 2.3: Type System 이식 (완료 + 테스트)

2026-04-09 (Week 4)
├─ Task 3.1: Semantic Analyzer 이식 (완료 + 테스트)
└─ Task 3.2: IR Builder 이식 (완료 + 테스트)

2026-04-16 (Week 5)
├─ Task 3.3: Code Generator 이식 (완료 + 테스트)
└─ Task 4.1: 통합 테스트 (준비)

2026-04-23 (Week 6)
├─ Task 4.1: 통합 테스트 (완료)
├─ Task 4.2: 성능 최적화 (완료)
└─ Task 4.3: 문서화 (완료)

2026-04-30
└─ 🎉 v1.0.0 Release (FreeLang Julia Compiler)
```

---

## 💡 **추가 고려사항**

### 1. 라이선스
- Go Julia Compiler: MIT
- FreeLang Julia Compiler: MIT (동일)

### 2. 커뮤니티
- Julia 커뮤니티에 공개
- FreeLang 커뮤니티 통합

### 3. 향후 확장
- LLVM 백엔드 추가
- JIT 컴파일 고도화
- 병렬 처리 지원

---

**버전**: 1.0 (계획)
**상태**: 🟢 계획 수립 완료
**다음 단계**: Task 1.1 시작 (Julia 구조 분석)
