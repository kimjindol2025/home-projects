# 📚 프로젝트 메모리 인덱스

> 최신 업데이트: 2026-03-29 | 활성 프로젝트 3개 | 완료: 15개 아카이브 이동

---

## 🚀 **활성 프로젝트**

### freelang-evolving-compiler
- [Phase 7: 루프 최적화](./phase7-loop-optimization.md) — LICM + Loop Unrolling (157줄)
- [Phase 6: 논리 연산자](./phase6-logical-operators.md) — !, &&, || 완전 구현 (256줄, 6파일)
- [Phase 5: Array 지원](./phase5-array-support.md) — 배열 리터럴/인덱싱/범위검사 (187줄, 8파일)
- [Phase 4: String/Bool REPL](./phase4-string-bool-repl.md) — 문자열/불리언 + 대화식 REPL (383줄)
- [Type System Phase 2](./type-system-phase2-complete.md) — Type Inference + Hard Mode (572줄)
- [Type System Phase 1](./type-system-foundation-phase1.md) — 기초 타입 체커 (530줄)
- [struct Phase 3: Field Access](./struct-phase3-field-access.md) — obj.field 문법 (251줄)
- [struct Phase 2: IR + CodeGen](./struct-phase2-ir-codegen.md) — AST→IR→Assembly (256줄)

### Zero-Copy-DB
- [Phase 12: 성능 최적화](./zdb-phase12-complete.md) — 인덱스 + 통계 캐싱 (852줄, 6파일, **20,392줄 누적**)
- [Phase 11: 통합 쿼리 파이프라인](./zdb-phase11-complete.md) — SQL 파싱 + 플래너 + 실행 (1,358줄)
- [Phase 10: 쿼리 실행 엔진](./zdb-phase10-complete.md) — 물리계획 → ResultSet (1,177줄)
- [Phase 9: 물리 계획 + CodeGen](./zdb-phase9-complete.md) — 논리→물리계획 변환 (2,092줄)

### Bigwash Native Shell (BNS)
- [Phase 3: 동적 데이터](./bns-phase3-dynamic-data.md) — MEMORY/git 파싱 (307줄 신규)
- [Phase 2: Webhook + SSE](./bns-phase2-webhook-sse.md) — Gogs 실시간 피드 (428줄)
- [Phase 1: MVP](./bns-phase1-mvp.md) — 100% FreeLang HTTP 서버 (730줄, 0 외부의존성)

---

## 📦 **완료됨 (아카이브)**

다음 15개 항목이 [`MEMORY_ARCHIVE.md`](./MEMORY_ARCHIVE.md)로 이동됨:
- [COMPLETE] struct Phase 1-2
- [COMPLETE] Zero-Copy-DB Phase 1-8
- [COMPLETE] GitHub 블로그 저장소
- [COMPLETE] 완벽한 분신 시스템 Step 3
- 외 11개 마일스톤

**아카이브 접근**: 과거 학습이 필요하면 [`MEMORY_ARCHIVE.md`](./MEMORY_ARCHIVE.md) 참조

---

## 🔑 **핵심 통계**

| 프로젝트 | 상태 | 규모 | 최신 |
|---------|------|------|------|
| freelang-evolving-compiler | 🟢 활성 | 7개 phase | Phase 7 |
| Zero-Copy-DB | 🟢 활성 | 4개 phase (12종) | Phase 12 |
| BNS | 🟢 활성 | 3개 phase | Phase 3 |

---

## 💡 **다음 단계**

- 각 `.md` 파일에 **세부 아키텍처** 및 **코드 위치** 기록됨
- 기억할 특정 내용이 있으면 이 인덱스 상단에 추가
- 완료된 항목은 자동으로 MEMORY_ARCHIVE.md로 정리
