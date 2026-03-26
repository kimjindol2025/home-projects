---
name: fir 프로젝트 기본기 및 도전 원칙
description: fir (Free IR Platform) 프로젝트의 핵심 원칙 — 어떤 작업을 하든 이 기본기는 고수한다
type: project
---

## fir 프로젝트 — 기본기 (절대 흔들리지 말 것)

**Why:** 이 프로젝트는 메타 플랫폼을 목표로 한다. 기반이 흔들리면 위에 쌓는 모든 것이 무너진다.
**How to apply:** 어떤 작업 요청이 오든 이 원칙에 위배되면 먼저 지적하고 진행한다.

---

### 1. IR이 본질이다

- 문법 변경, 기능 추가, 서버 기능은 IR 안정화 전에 하지 않는다
- IR_SPEC_1.0.md는 FROZEN — 변경은 IR 2.0 버전업으로만
- 구현이 문서를 앞서면 안 된다

### 2. 순서를 지킨다

```
interp 검증 (test 10/10) → codegen (IR→x86) → native==interp diff → 선언
```

- 선언은 증명 후에 한다
- 증명 없는 선언은 빈 말이다

### 3. 테스트가 완료 기준이다

- "될 것 같다"는 완료가 아니다
- 테스트 PASS가 완료다
- Day 기준: test_interp.fl 10/10 PASS = Day 1 완료

### 4. 3일 계획 중 하지 말 것

- 최적화
- 속도 개선
- 구조 변경
- IR 수정
- "이거 조금 고치면 더 깔끔한데?" → 그 순간 3일이 3주 된다

### 5. 모델/도구가 아니라 집중이 필요하다

- 모델 변경 불필요 — Sonnet으로 충분
- 필요한 것: 집중, 순서, 검증

### 6. 메타 플랫폼 선언 조건

다음이 모두 완성된 후에만 선언한다:
- [ ] interp 10/10 PASS
- [ ] codegen.fl (IR → x86-64)
- [ ] native 결과 == interp 결과
- 그때 그림 그리고 선언

---

## 저장소

- GOGS: https://gogs.dclub.kr/kim/fir.git
- 로컬: /data/data/com.termux/files/home/.projects/core/freelang-v7/
- 브랜치: master

## 현재 위치 (2026-03-22)

```
✅ IR_SPEC_1.0.md  FROZEN
✅ ir_node.fl      IR 노드 전체
✅ ir_gen.fl       AST → IR
✅ interp.fl       reference interpreter (Day 1 버그 3개 수정)
✅ test_interp.fl  10개 테스트 작성
⏳ 실행 검증      FreeLang 런타임으로 test_interp.fl 실행 필요
❌ codegen.fl      미착수
```

## FreeLang과의 관계

프리랭은 버리지 않는다.
fir의 레퍼런스 프론트엔드다.
FreeLang syntax → Free IR → 실행.
