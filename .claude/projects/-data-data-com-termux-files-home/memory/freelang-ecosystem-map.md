---
name: FreeLang 생태계 전체 지도
description: FreeLang 프로젝트 전체 구조, 런타임 현황, 실행 가능 여부 — 매 대화 재탐색 방지용 영구 참조
type: reference
updated: 2026-03-23
---

# FreeLang 생태계 전체 지도

> 조사일: 2026-03-23 | 242개 프로젝트, 347개 .fl/.free 소스 파일

---

## 핵심 3줄 요약

1. **v2-freelang-ai** (TypeScript 2.8.0) — 지금 당장 실행 가능한 유일한 완전한 FreeLang 런타임 (176/176 PASS, Self-hosting 가능)
2. **fir** (Free IR 1.0, `.projects/core/freelang-v7/`) — 미래 표준. IR 인터프리터+codegen 완성, Phase 3 진행 중
3. **gofree** — Go 이식 중 (Phase 3 Parser), Phase 6 완료 시 첫 `.fl` 실행 가능

---

## 런타임 실행 준비도 순위

| 순위 | 프로젝트 | 경로 | 상태 | 실행 명령 |
|-----|---------|------|------|---------|
| 1 | v2-freelang-ai | `~/v2-freelang-ai/` | 즉시 가능 | `npm run build && npx freelang hello.free` |
| 2 | fir 인터프리터 | `~/.projects/core/freelang-v7/` | 제한적 (Phase 2 기능) | interp.fl 통해 실행 |
| 3 | fv2-lang-go | `~/projects/fv2-lang-go/` | .fv→C 컴파일 가능 | `./bin/fv2 hello.fv` |
| 4 | gofree | `~/gofree/` | Parser 미완 | Phase 6 완료 후 |

---

## 프로젝트 계열별 위치

### 런타임/컴파일러
| 프로젝트 | 경로 | 언어 | 상태 |
|---------|------|------|------|
| fir (Free IR 1.0) | `~/.projects/core/freelang-v7/` | FreeLang+Go | Phase 3 진행 |
| gofree | `~/gofree/` | Go | Phase 3 Parser |
| freelang-runtime | `~/.projects/modules/freelang-runtime/` | Rust | Production |
| freelang-vm | `~/.projects/modules/freelang-vm/` | Rust | Active |
| freelang-jit-compiler | `~/.projects/modules/freelang-jit-compiler/` | Rust | Production |
| freelang-aot-compiler | `~/.projects/modules/freelang-aot-compiler/` | Rust | Active |

### 언어 버전
| 프로젝트 | 경로 | 언어 | 상태 |
|---------|------|------|------|
| v2-freelang-ai | `~/v2-freelang-ai/` | TypeScript | 2.8.0 Production |
| freelang-v2 | `~/projects/freelang-v2/` | TypeScript | 보관됨 |
| freelang-hybrid | `~/.projects/core/freelang-hybrid/` | TypeScript | Core |
| freelang-light | `~/.projects/core/freelang-light/` | TypeScript | Core |
| fv-lang | `~/projects/fv-lang/` | Rust | 개발 중 |
| fv2-lang-go | `~/projects/fv2-lang-go/` | Go | B+ 완성 |
| fv-lang-wasm | `~/projects/fv-lang-wasm/` | Rust | WIP |

### 서버/HTTP (.fl 구현)
| 프로젝트 | 경로 | 구현 파일 | 상태 |
|---------|------|---------|------|
| freelang-http-engine | `~/.projects/core/freelang-http-engine/src/` | server.fl, tcp_socket.fl, http_parser.fl, http_handler.fl, mod.fl (총 35KB) | 코드 완성, 실행 미검증 |
| freelang-v4-http | `~/.projects/core/freelang-v4-http/` | TypeScript | Production |
| freelang-rest-api | `~/.projects/core/freelang-rest-api/` | - | Active |

### 응용
| 프로젝트 | 경로 | 언어 | 상태 |
|---------|------|------|------|
| freelang-gpt | `~/projects/freelang-gpt/` | Go | Phase H 완료, 26 API |
| freelang-nexus | `~/projects/freelang-nexus/` | Go | Phase N3 완료 |
| freelang-to-c | `~/projects/freelang-to-c/` | Rust | Active |
| faf | `~/projects/faf/` | - | - |

---

## fir 파이프라인 (현재 표준)

```
FreeLang 소스 (.fl)
    ↓ ir_gen.fl (파서+변환기)
Free IR 1.0 (SSA, FROZEN)
    ├─ interp.fl → 인터프리터 실행 (정확성 기준)
    └─ codegen → x86-64 ELF64 네이티브
```

**fir IR 1.0 불변 규칙**:
- 타입: i8~i64, f32/f64, bool, void, ptr<T>, array<T,N>, slice<T>
- SSA: 각 레지스터 정확히 한 번 정의
- ABI: C ABI x86-64 System V
- 모듈: DAG (순환 의존 없음)

---

## 중요 질문 답변

### Q: `.fl` 파일 지금 실행 가능한가?
- **v2-freelang-ai로**: YES (`npm run build` 후 즉시)
- **fir 인터프리터로**: 제한적 YES (Phase 2 기능 범위 내)
- **gofree로**: NO (Phase 3 미완)

### Q: HTTP 서버를 FreeLang으로?
- **freelang-http-engine/src/server.fl** — 35KB 완전 구현 존재
- 실행 경로: v2-freelang-ai 또는 fir 인터프리터 통해 가능 (미검증)

### Q: 가장 완성도 높은 컴파일러?
- **프로덕션**: v2-freelang-ai (2.8.0, Self-hosting 가능)
- **아키텍처**: fir (IR 철학, 수학적 증명, 미래 표준)

---

## 구현 언어 분포
- TypeScript/Node.js: 60+ 프로젝트
- Rust: 18 프로젝트
- Go: 11 프로젝트
- FreeLang 자체: 5+ 프로젝트
