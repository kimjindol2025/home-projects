---
name: FreeWire Phase 4
description: FreeWire 노드 기반 시각 프로그래밍 언어 Phase 4 완료 (102/102 테스트)
type: project
---

FreeWire Phase 1-4 완료 (2026-03-17). GOGS: kim/FreeWire (커밋: a137527)
로컬 경로: `projects/FreeWire/`

**Why:** 노드 기반 시각 프로그래밍 → FreeLang 코드 생성. React Flow 없이 순수 Canvas/SVG.

**How to apply:** Phase 5 (AI 노드 자동 생성), Phase 6 (WASM 브라우저 REPL) 예정.

## 누적 테스트
- Phase 1: 32/32 (노드 JSON 스펙 — component/data/logic + Graph)
- Phase 2: 16/16 (Graph→FreeLang 컴파일러)
- Phase 3: 24/24 (API 서버 + FreeExecutor)
- Phase 4: 30/30 (노드 에디터 UI — 포트호환성/베지어/그래프모델/컴파일/팔레트)
- **합계: 102/102**

## Phase 4 구현
- `src/ui/node-editor.js` — 서버사이드 모듈
- `www/index.html` — 완전 브라우저 UI (인라인 JS, 의존성 없음)
- `src/test/test-phase4.js` — 30개 테스트
- 팔레트 3그룹, 4 시작 템플릿, JSON 내보내기/가져오기
- 실행: `python3 -m http.server 5500 --directory www`
- API 서버: `node src/api/server.js` (포트 4000)
