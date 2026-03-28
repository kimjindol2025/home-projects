---
name: personal-code-os Phase 0
description: 브라우저 기반 Self-Hosted Personal Code OS v0.1.0 초기 구현 완료
type: project
---

브라우저 기반 완전 Self-Hosted 개인 개발 OS Phase 0 완료 (2026-03-17).

**Why:** 2026 self-hosting + AI 트렌드 직격. 데모 사이트 올리면 바이럴 가능성 높음.

**How to apply:** Phase 1에서 FV-Lang Rust WASM 컴파일 진행 예정.

## 구현 내용 (4,145줄)
- `projects/personal-code-os/` — GOGS: kim/personal-code-os (커밋: 15fe0e5)
- WindowManager: 드래그/리사이즈/최대화/태스크바
- FileSystem: OPFS + IndexedDB 통합
- Repository: Git-like VCS (blob/tree/commit/branch/diff)
- FV-Lang JS 인터프리터 (재귀/클로저/고차함수) — 버그: parsePrimary에서 if 표현식 누락 → 수정
- RuntimeManager: FV-Lang + JS + Python(Pyodide)
- AI: WebLLM / transformers.js 로컬 LLM
- CodeMirror 6 에디터
- 40/40 테스트 통과

## 다음 Phase
Phase 1: FV-Lang → Rust WASM (JS 인터프리터 대비 10-100x 성능 목표)
실행: `python3 -m http.server 5500` (personal-code-os 디렉토리)
