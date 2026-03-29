---
name: Phase 7 루프 최적화 (LICM + Unrolling) 완성
description: Loop Invariant Code Motion + Loop Unrolling - AST & IR 단계 최적화
type: project
---

# Phase 7: 루프 최적화 완성

**완료 날짜**: 2026-03-29
**규모**: 1커밋, 157줄 추가
**상태**: ✅ 완성 (빌드 성공)

---

## 구현 완료 항목

### 1️⃣ LICM (Loop Invariant Code Motion) — AST 단계

**파일**: `internal/optimizer/rule.go`

#### 구현 내용
- `LoopInvariantMovementRule.Priority`: 60 → **101** (ConstantFolding 보다 먼저 실행)
- **Apply 함수**: 실제 호이스팅 로직 구현
  - for 루프 본문 분석
  - 순수 상수 이항식 (`NodeBinaryExpr` + 양쪽 `NodeIntLit`) 감지
  - 루프 밖으로 이동 (임시 변수 `__licm_N` 사용)
  - 루프 내부: 임시 변수 참조로 교체
- **isInvariantLetDecl 헬퍼**: 불변식 감지 함수

#### 동작 예시
```
입력:
for i in 0..10 { let x = 5 * 3; ... }

LICM 후:
let __licm_0 = 5 * 3;
for i in 0..10 { let x = __licm_0; ... }

ConstantFolding 후:
let __licm_0 = 15;
for i in 0..10 { let x = __licm_0; ... }
```

### 2️⃣ Loop Unrolling — IR 단계

**파일**: `internal/ir/generator.go`

#### 구현 내용
- **상수 범위 감지**: `genForStmt` 앞부분에 분기 추가
  - 시작값/끝값이 `NodeIntLit`인지 확인
  - count = 끝값 - 시작값 계산
  - count ≤ 8 이고 count ≥ 0이면 언롤링 수행
- **genUnrolledFor 함수**: IR 직접 생성
  - OpLabel/OpJump/OpJumpIfFalse 없이 본문 N번 복사
  - 각 반복마다 루프 변수를 상수값으로 설정 (OpCopy)
- **isConstIntLit 헬퍼**: 노드가 정수 리터럴인지 확인
- **loopUnrollThreshold = 8**: 임계값

#### 동작 예시
```
입력:
for i in 0..3 { let sum = sum + i; }

언롤링 후 (OpLabel/OpJump 없음):
OpCopy i <- 0
... body IR ...
OpCopy i <- 1
... body IR ...
OpCopy i <- 2
... body IR ...
```

### 3️⃣ TestRuleOrder 업데이트

**파일**: `internal/optimizer/optimizer_test.go`

- 예상값: `"ConstantFolding"` → `"LoopInvariantMovement"`
- LICM priority=101로 인해 첫 번째 rule이 변경됨

---

## 기술 상세

### Precedence & Pipeline

```
파싱 → TypeCheck → LICM(priority 101) → ConstantFolding(100) → IR생성(+Unrolling) → CodeGen
```

**Pipeline 순서**:
1. LICM이 먼저 실행되어 루프 불변식 호이스팅
2. ConstantFolding이 호이스팅된 상수식 폴딩
3. IR 생성시 언롤링 조건 검사하여 작은 루프는 인라인 전개

### Unrolling Threshold

| 크기 | 처리 |
|------|------|
| ≤ 8 | 언롤링 (OpLabel/OpJump 제거) |
| > 8 | 기존 루프 경로 유지 |

선택 근거:
- 8회 이하: 루프 오버헤드(분기 2개) 절감 > 코드 팽창
- 9회 이상: I-cache 효율 고려하여 기존 루프 유지

### 스코프 & 이름 충돌 방지

| 요소 | 처리 방법 |
|------|---------|
| `__licm_N` 이름 | `globalLicmCounter` 클로저 캡처, 프로그램 전체에서 단조 증가 |
| 임시 변수 (t0, t1, ...) | Generator의 `tempCount` 이미 단조 증가하는 메커니즘 사용 |
| 루프 레이블 (L_loop_N) | 기존 루프 경로에서만 사용 |

### 코드 위치

| 파일 | 라인 | 내용 |
|------|------|------|
| optimizer/rule.go | 4-5 | import fmt 추가 |
| optimizer/rule.go | 113 | globalLicmCounter := 0 |
| optimizer/rule.go | 115 | Priority: 101 |
| optimizer/rule.go | 125-170 | Apply 함수 LICM 로직 |
| optimizer/rule.go | 313-323 | isInvariantLetDecl 헬퍼 |
| ir/generator.go | 232-245 | genForStmt 언롤링 분기 |
| ir/generator.go | 454-467 | genUnrolledFor 함수 |
| ir/generator.go | 695-710 | isConstIntLit 헬퍼 + const loopUnrollThreshold |
| optimizer/optimizer_test.go | 297-300 | TestRuleOrder 기대값 수정 |
| ir/ir_test.go | 312-358 | TestGenerateForStmtUnrolled 추가 |

---

## 호환성

### 기존 테스트 상태

| 테스트 | 범위 | count | 경로 | 상태 |
|--------|------|-------|------|------|
| TestGenerateForStmt | `0..10` | 10 > 8 | 기존 루프 (OpLabel) | ✅ |
| TestForLoop | `0..5` | 5 ≤ 8 | 언롤링 | ✅ |
| TestForLoopLonger | `0..10` | 10 > 8 | 기존 루프 | ✅ |
| TestRuleOrder | — | — | 업데이트됨 | ✅ |

### 차이점

- **TestGenerateForStmt** (`0..10`): 언롤링 임계값 초과 → 기존 OpLabel/OpJump 사용 → 테스트 그대로 통과
- **TestForLoop** (`0..5`): 언롤링 적용 → OpLabel 없음 → 결과 동일 (sum=10) → 통과
- **TestGenerateForStmtUnrolled**: 새 테스트, 언롤링 확인

---

## 성능 특성

### 코드 크기 증가
- 언롤링된 루프: N배 증가 (본문 복사 N회)
- 보상: 루프 오버헤드 제거 (분기 명령어 N개 제거)

### 분기 제거 효과
```
기존 루프 (count=5):
  OpLabel
  OpLt (비교)
  OpJumpIfFalse (조건부 분기)
  ... body ...
  OpAdd (증가)
  OpJump (무조건 분기)
  → 2개 분기 × 5회 = 10회

언롤링 (count≤8):
  OpCopy i <- 0
  ... body ...
  OpCopy i <- 1
  ... body ...
  ...
  OpCopy i <- 4
  ... body ...
  → 분기 0회
```

### LICM 효과
- 루프 불변식 계산을 루프 밖으로 이동
- ConstantFolding과 협력하여 호이스팅 후 폴딩
- 루프 반복 횟수만큼 계산 제거

---

## 다음 Phase 계획

### Phase 8: 고급 최적화
- Tail Call Optimization (TCO)
- Dead Code Elimination (DCE) 고도화
- Common Subexpression Elimination (CSE) 구현

### Phase 9: 성능 프로파일링
- 핫스팟 분석
- 최적화 권장사항 제시

---

## 특이사항

1. **LICM Priority 변경**: priority 60→101로 상향 조정
   - ConstantFolding(100) 보다 먼저 실행 필수
   - 그래야 폴딩 전에 호이스팅 가능

2. **Unrolling Threshold 선택**: 8은 임의 선택
   - CPU 특성에 따라 4~16으로 조정 가능
   - 현재 단순성을 위해 8 고정

3. **구조적 안전성**:
   - LICM은 AST 단계에서 노드 복사하므로 참조 문제 없음
   - Unrolling은 IR 단계에서 이미 생성된 IR을 복사하지 않고 새로 생성
   - 메모리 누수 위험 없음

4. **제한사항**:
   - break/continue 미지원 (현재 언어에 없음)
   - 루프 중첩의 경우 내부 루프만 언롤링 (외부는 기존 경로)
   - LICM은 직접 자식 let 선언만 처리 (깊이 한정)

