---
name: fir Phase 3 코드 품질 리뷰 결과
description: 2026-03-24 코드 품질 검사 결과 — 72/100점, 핵심 문제 3가지
type: project
---

# fir Phase 3 코드 품질 리뷰 (2026-03-24)

**종합 점수**: 72/100

**Why:** 수동 코드 품질 체크 결과. 수정 전 기준점으로 보존.
**How to apply:** fir 작업 재개 시 아래 3가지 문제를 우선 해결 대상으로 삼을 것.

---

## 항목별 점수

| 영역 | 점수 | 비고 |
|------|------|------|
| 정확성 (테스트) | 16/20 | 97/98 PASS. RA-3b FAIL 1개 |
| 설계 일관성 | 15/20 | Pass 인터페이스 깔끔, 레이어 분리 명확 |
| 코드 중복 | 10/20 | codegen.go + regalloc_codegen.go 거의 동일 |
| 에러 처리 | 12/20 | compileInst 에러 리턴 있지만 실제론 nil만 반환 |
| 명확성/주석 | 19/20 | 알고리즘 출처 명시(Poletto, Cytron, Cooper) 우수 |

---

## 핵심 문제 3가지 (우선순위 순)

### 1. 코드 중복 (가장 심각)
- `src/codegen/codegen.go` 763줄 — `FuncCompiler`
- `src/codegen/regalloc_codegen.go` 458줄 — `FuncCompilerRA`
- `compileInst` 내용이 거의 동일하게 중복
- **개선 방향**: 공통 베이스 추출 또는 인터페이스 추상화

### 2. 버그: phi lowering 미구현
- 위치: `src/codegen/regalloc_codegen.go:351` (InstPhi → no-op)
- 증상: RA-3b FAIL — if/else false 경로에서 native=0, expected=2
- 원인: `br`/`brcond` 직전에 phi incoming 값을 슬롯/레지스터에 쓰는 "phi copy" 코드 없음
- predecessor 블록이 phi 슬롯에 아무것도 쓰지 않고 분기함
- 체크포인트 태그: `checkpoint/phase3-ra-3b-bugfix-before`

### 3. 사문 코드
- 위치: `src/opt/regalloc.go:137`
- 내용: `_ = u` — 컴파일러 경고 회피 흔적, 의미 없음
