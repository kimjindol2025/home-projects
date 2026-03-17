---
name: FV-Lang Phase 5 WASM
description: FV-Lang을 WebAssembly로 브라우저에서 실행하는 REPL 프로젝트 진행 상황
type: project
---

# FV-Lang Phase 5 - WASM 브라우저 REPL

## 목표
- FV-Lang 코드를 브라우저에서 직접 입력 → 즉시 파싱/실행 → 결과 출력 + AST 시각화(SVG)
- freelang-light 스타일 컴포넌트 DSL을 FV-Lang으로 작성 → WASM eval → DOM 동적 렌더링
- 최종 데모: "브라우저에서 함수형 언어로 실시간 웹 UI 빌드"

## 로드맵
- Phase 0: 사전 준비 (1~2일) - wasm-pack, trunk 설치, fv-lang-wasm crate 생성
- Phase 1: fv-lang 코어 WASM 포팅 (5~7일)
- Phase 2: 브라우저 REPL 구현 (7~10일)
- Phase 3: freelang-light 통합 + AST 시각화 (5~7일)
- Phase 4: Polishing + 배포 (3~5일)

## 현재 상태
🟡 아직 시작 전 (Phase 0 준비 단계)

## Why
**Why:** Rust 기반 함수형 언어 + WASM = 2026년에도 드문 케이스. 포트폴리오 임팩트 큼.

## How to apply
**How to apply:** 다음 세션에서 "FV-Lang WASM" 언급 시 Phase 0부터 시작. fv-lang-wasm crate 신규 생성부터.
