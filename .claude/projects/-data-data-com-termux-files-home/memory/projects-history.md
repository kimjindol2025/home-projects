---
name: projects-history
description: 완료된 프로젝트 히스토리 상세 목록 (MEMORY.md 인덱스에서 분리)
type: project
---

# 완료 프로젝트 히스토리

## FVX — Formal Verification Tool (2026-03-24) ✅ 완전 완료
- **저장소**: https://gogs.dclub.kr/kim/fvx.git
- **로컬**: `/data/data/com.termux/files/home/.projects/core/fvx/`
- **규모**: 32개 파일 / 6,996줄 / GOGS 커밋 9개
- **최종 테스트**: 20/20 PASS (Phase K 17/17 + Deep Stress 3/3)
- **핵심 성과**:
  - 5개 공리(A1~A5) 기반 형식 검증 엔진
  - Tarski Fixed Point Solver + AC-3 전파 + μ 함수 정지 증명
  - FreeJulia IR → FVX IR 브리지 (7개 Bytecode 타입 1:1 변환)
  - `fvx run/report/ir/check/test-all/repl` CLI
  - 수학적 증명 예제 3종: range_check / sort_verify / fib_safety
  - Deep Stress: 공리충돌(TOP 단락) / N=2000 O(N) / √2 수렴(ε=1e-10)
- **증명된 이론**: Tarski / Inversion Pair φ / Well-Founded Relation μ(k)
- **성능**: 2,000노드 체인 → 2회 반복 / 3.90ms (O(N) 선형)
- **다음 연동 후보**: Mandate-DB 트랜잭션 검증 / FreeJulia 정적 분석기

## fir — Free IR Platform (2026-03-22)
- **저장소**: https://gogs.dclub.kr/kim/fir.git
- **로컬**: `/data/data/com.termux/files/home/.projects/core/freelang-v7/`
- **완료**: Interpreter 17/17 + Codegen 24/24 + IRGen 16/16 PASS (91/91)
- **현재**: Phase 3 — ConstFold → DCE → struct/GEP → mem2reg → RegAlloc
- **원칙**: [fir-project-principles.md](./fir-project-principles.md)
- **핵심**: fir는 언어가 아니다. 실행 의미를 정의하는 IR 플랫폼. FreeLang is a Free IR frontend.

## FreeLang GPT (2026-03-21)
- **위치**: projects/freelang-gpt
- Priority 3 Adam Optimizer: 하이브리드 (FL 285줄 + Go 376줄), 7/7 PASS
- Phase H Checkpoint: 2,241줄, 26개 REST API, 10/10 PASS
- Phase G REST API: 7/8 PASS, 1,388줄, Docker/Nginx 배포
- 누적 코드: ~12,228줄 Go (Phase 1-6 + A-H)

## FV 2.0 Go (2026-03-21)
- **완성도**: 100% (파서 2개 + 코드젠 6개 갭 해결)
- extern fn void, else if 체인, 배열선언, MethodCall/Struct codegen
- Phase 7 최종 등급: B+ (9,895줄 Go + 9,116줄 테스트)
- 5개 언어 혼합 통합 (FV + C FFI + PyFree CLI)

## FreeLang Nexus (2026-03-20)
- Phase 1-7 완료: 52/52 테스트
- 기능: run/compile/check/repl CLI

## FreeJulia (2026-03-20)
- Phase A-H 완료: 92% 완성도, 21,036줄, 451+ 테스트
- Phase D Self-Hosting: 4,241줄 + 121개 테스트
- Phase G VFS+Collections+Integration: 14,351줄, 398+ 테스트

## FV-Julia (2026-03-20)
- Phase E 완료: 11/11 컴파일 테스트, 93.3% 기능 커버리지
- Phase 1 CodeGen: 1,422줄, 50개 테스트

## Multi-Language PoC (2026-03-21)
- Phase 1: 859줄, 14/14 테스트 PASS
- 6개 파일, 4개 언어 동시 실행, 12개 타입 매핑

## FV-Lang Go Phase 1-5 (기초)
- 3,650줄 Go + 1,020줄 테스트, 58/58 통과
- GOGS: 커밋 1cf1804

## Sovereign Workspace v1.0.0
- 20,370줄 FV-Lang + 393개 테스트
- GOGS: https://gogs.dclub.kr/kim/sovereign-workspace
- Phase 1-11 완성, Docker/Termux/K8s 배포 옵션

## Julia Compiler v0.2.0
- Code Quality: 3/10 → 9/10, 커밋 79723259

## Genspark Clone v3.0
- 4개 버그 수정, 36/36 테스트, 커밋 2230abf
