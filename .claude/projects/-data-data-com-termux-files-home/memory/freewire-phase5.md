---
name: FreeWire Phase 5-6
description: FreeWire Phase 5 (AI 노드 자동 생성) + Phase 6 (브라우저 Live Preview) 완료. 162/162 테스트.
type: project
---

FreeWire Phase 1-6 완료 (2026-03-17). GOGS: kim/FreeWire (커밋: 692a360)
로컬 경로: `projects/FreeWire/`

**Why:** Phase 5: 자연어 → 노드 그래프. Phase 6: 노드 → 실시간 HTML 미리보기 (서버 불필요).

**How to apply:** FreeWire 기능 완성. personal-code-os 통합 가능.

## 누적 테스트
- Phase 1: 32/32 (노드 JSON 스펙)
- Phase 2: 16/16 (Graph→FreeLang 컴파일러)
- Phase 3: 24/24 (API 서버 + FreeExecutor)
- Phase 4: 30/30 (노드 에디터 UI)
- Phase 5: 35/35 (AI 노드 자동 생성기)
- Phase 6: 25/25 (브라우저 Live Preview)
- **합계: 162/162**

## Phase 6 구현
- `www/index.html` 분할 뷰 (에디터 좌 + 미리보기 패널 우 360px)
- `browserCompile()`: 노드 그래프 → HTML/CSS/JS/FreeLang (서버 불필요)
  - generateFreeLang(): component/state/markup 블록
  - generateHTML(): 와이어 기반 자식 배치, data-bind 연동
  - generateCSS(): 컴포넌트 스타일 (버튼/카드/입력/컨테이너)
  - generateJS(): state IIFE + DOMContentLoaded 바인딩
- "▶ 실행" 버튼 → iframe Live Preview
- 미리보기 탭 ↔ FreeLang 코드 탭 전환
- `src/test/test-phase6.js` — 25개 테스트
- 실행: `python3 -m http.server 5500 --directory www`
