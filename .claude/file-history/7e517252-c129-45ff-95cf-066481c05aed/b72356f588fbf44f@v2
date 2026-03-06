# Z-Lang Transpiler - Phase 3: Optimization Plan

**프로젝트**: FreeLang to Z-Lang Transpiler
**Phase**: 3 (최적화)
**기간**: 2026-03-06 ~ 2026-03-09
**규모**: 1,000줄 + 20개 테스트

---

## 📋 Phase 3 목표

### 목표 1: Loop Optimization (300줄)
- [ ] **loop_optimizer.fl** (300줄)
  - LICM (Loop-Invariant Code Motion)
  - Induction Variable 제거
  - Loop Unrolling 기본 구현
  - Nested loop 처리

### 목표 2: Constant Folding (250줄)
- [ ] **constant_folder.fl** (250줄)
  - 컴파일 타임 상수 계산
  - 타입 안전성 유지
  - 부분 평가 (Partial Evaluation)
  - 문자열 상수 최적화

### 목표 3: Dead Code Elimination (200줄)
- [ ] **dead_code_eliminator.fl** (200줄)
  - 도달 불가능한 코드 제거
  - 부작용 없는 명령어 제거
  - 사용되지 않는 변수 추적
  - Control Flow Graph 분석

### 목표 4: Register Allocation (200줄)
- [ ] **register_allocator.fl** (200줄)
  - Linear Scan Register Allocation
  - Liveness Analysis
  - Spilling 전략
  - Coalescing 기본 구현

### 목표 5: Integration & Benchmarking (150줄)
- [ ] **optimizer_integration.fl** (150줄)
  - 최적화 파이프라인
  - IR 변환 검증
  - 성능 메트릭 수집

---

## 🧪 Test Plan (20개 테스트)

### Group P1: Loop Optimization (5개)
- [ ] **P1-T1**: LICM - 불변 코드 이동
- [ ] **P1-T2**: IV 제거 - Induction variable 최소화
- [ ] **P1-T3**: Nested loops - 중첩 루프 처리
- [ ] **P1-T4**: Loop unrolling - 루프 전개
- [ ] **P1-T5**: Side effects - 부작용 있는 코드 보존

### Group P2: Constant Folding (5개)
- [ ] **P2-T1**: 정수 연산 - 1 + 2 = 3
- [ ] **P2-T2**: 실수 연산 - 2.5 * 2.0 = 5.0
- [ ] **P2-T3**: 문자열 - concat 최적화
- [ ] **P2-T4**: Mixed types - 타입 변환 최적화
- [ ] **P2-T5**: 중첩 표현식 - (1 + 2) * (3 + 4) = 21

### Group P3: Dead Code Elimination (5개)
- [ ] **P3-T1**: 도달 불가능한 코드 - if false { ... }
- [ ] **P3-T2**: 사용되지 않는 변수 - let x = 5; (unused)
- [ ] **P3-T3**: 도달 불가능한 분기 - return 후 코드
- [ ] **P3-T4**: 부작용 없는 명령어 - x + 1; (unused result)
- [ ] **P3-T5**: 복잡한 CFG - 다중 분기 정리

### Group P4: Register Allocation (5개)
- [ ] **P4-T1**: Liveness analysis - 변수 생존 범위
- [ ] **P4-T2**: Register coloring - 최소 색상 수
- [ ] **P4-T3**: Spilling - 메모리 사용
- [ ] **P4-T4**: Coalescing - 복사 제거
- [ ] **P4-T5**: 성능 메트릭 - 할당 효율성

---

## 🎯 무관용 규칙 (Unforgiving Rules)

### R1: Optimization Correctness
- **정의**: 최적화 후 의미론(Semantics) 변경 금지
- **검증**: 최적화 전후 실행 결과 동일
- **테스트**: 모든 테스트에서 검증

### R2: Performance Improvement
- **정의**: 최적화 > 10% 성능 향상 또는 코드 크기 감소
- **검증**: IR 크기 또는 실행 시간 비교
- **목표**: 평균 20% 개선

### R3: Type Safety Preservation
- **정의**: 최적화 후 타입 안전성 보존
- **검증**: 타입 체커 통과
- **테스트**: P2-T4 (Mixed types)

### R4: Side Effect Awareness
- **정의**: 부작용 있는 코드 절대 제거 불가
- **검증**: 부작용 분석 완전성
- **테스트**: P1-T5, P3-T4

---

## 📊 파일 구조

```
src/
  ├─ loop_optimizer.fl        (300줄) ← NEW
  ├─ constant_folder.fl       (250줄) ← NEW
  ├─ dead_code_eliminator.fl  (200줄) ← NEW
  ├─ register_allocator.fl    (200줄) ← NEW
  ├─ optimizer_integration.fl (150줄) ← NEW
  ├─ mod.fl                   (업데이트)
  ├─ parser.fl                (Phase 2)
  ├─ ir.fl                    (Phase 2)
  └─ codegen.fl               (Phase 2)

tests/
  ├─ phase3_tests.fl          (300줄, 20개 테스트) ← NEW
  ├─ phase3_unforgiving.fl    (100줄, 4개 규칙)   ← NEW
  └─ phase2_tests.fl          (Phase 2)

docs/
  └─ PHASE_3_OPTIMIZATION_REPORT.md (최종 보고서)
```

---

## 🚀 구현 전략

### 1단계: Loop Optimizer (Day 1-2)
- CFG (Control Flow Graph) 구축
- Dominance Frontier 계산
- LICM 알고리즘 구현
- IV 제거 (Strength Reduction)

### 2단계: Constant Folder (Day 2-3)
- 상수 표현식 인식
- 타입별 연산 구현 (int, float, string)
- 부분 평가 (Partial Evaluation)
- 메모이제이션

### 3단계: Dead Code Eliminator (Day 3)
- 도달 가능성 분석
- 부작용 추적
- Liveness 분석
- 반복 제거

### 4단계: Register Allocator (Day 3-4)
- Linear Scan 알고리즘
- Spill 구간 결정
- Coalescing 규칙
- 할당 검증

### 5단계: Integration & Test (Day 4)
- 파이프라인 통합
- 성능 벤치마크
- 테스트 작성
- GOGS 푸시

---

## 📈 성능 목표

| 최적화 | 목표 | 측정 방법 |
|--------|------|----------|
| **Loop** | 루프 15% 가속 | 루프 반복 실행 시간 |
| **Constant** | 10줄 감소 | IR 명령어 수 |
| **DCE** | 20% 코드 제거 | 도달 불가능 코드 |
| **Register** | 메모리 30% 감소 | 스필 횟수 |
| **Overall** | 20%+ 개선 | 종합 성능 |

---

## 🎯 검증 기준

### ✅ 완료 조건
- [ ] 20개 테스트 100% 통과
- [ ] 4개 규칙 100% 검증
- [ ] 성능 목표 달성 (> 10%)
- [ ] 타입 안전성 보존
- [ ] 부작용 보존
- [ ] GOGS 푸시 완료

---

**시작일**: 2026-03-06
**예상 완료**: 2026-03-09
**상태**: 📋 계획 수립 완료
