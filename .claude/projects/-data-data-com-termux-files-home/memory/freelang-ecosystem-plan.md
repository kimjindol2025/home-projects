---
name: FreeLang 언어생태계 종합 조성계획 v1.0 (2026)
description: FLIR 기반 8개 컴파일러 통합, Q1~Q4 로드맵, 4분기 KPI 정의
type: project
---

# FreeLang 언어생태계 종합 조성계획 v1.0

**상태**: ✅ APPROVED & IN EFFECT (2026-03-21)
**문서**: `FREELANG_ECOSYSTEM_PLAN.md` (1,200+ 줄)
**저장소**: https://gogs.dclub.kr/kim/freelang-nexus

## 핵심 전략

### 문제 정의
- 현황: 8개 프로젝트 (Nexus, FV2-Go, FreeJulia, PyFree, FV-Lang, GPT, Multi-Lang, To-C)가 각자 독립 발전
- 사일로화: 각 프로젝트가 stdlib, 패키지, 문서를 모두 별도 구현
- 공통 기반 부재: FFI, 공통 타입 시스템, 패키지 매니저 없음

### 해결책: FLIR (FreeLang IR)
- 기반: `fv-lang/src/ir.rs`의 IRNode enum (이미 좋은 구조)
- 역할: 모든 컴파일러의 공통 언어
- 포맷: JSON 스키마 (프로젝트 간 호환성)
- 효과: Parser(FL) → FLIR JSON → Backend(C/Go/WASM) 표준화

### 기대 효과
- Nexus: 공식 레퍼런스 (FLIR 생성)
- FV2-Go: 네이티브 백엔드 (FLIR 처리)
- FreeJulia: FLIR 인터프리터 (FL 언어로)
- 통합: 어떤 방언으로도 같은 stdlib/패키지 사용 가능

## 4분기 로드맵

### Q1 (지금 ~ 3월): 기반 다지기

**주요 작업** (6주):
1. **FLIR v1.0 스펙** (2주)
   - FV-Lang ir.rs 분석
   - Loop, For, Match, Struct, FieldAccess, Closure, TypeCast, ExternCall 추가
   - JSON Schema + 설명서 작성

2. **컴파일러 연결** (2주)
   - Nexus FLIR emitter: AST → FLIR JSON
   - FV2-Go FLIR reader: FLIR JSON → 내부 AST
   - Roundtrip 테스트: 양방향 동일성 검증

3. **stdlib 코어 30개** (2주)
   - 명세: `core.json` (io, string, math, array, type)
   - C 구현: `c/core.h`
   - Go 구현: `go/core.go`

4. **공통 타입 시스템** (1주)
   - 기반: `fv2-lang-go/internal/typechecker/types.go`
   - 공통 타입: i8~i64, u8~u64, f32, f64, bool, string, void, never
   - 각 컴파일러 매핑표

**목표**: FLIR로 2개 컴파일러 연결, stdlib core 30개 함수 명세

### Q2 (4~6월): 도구체인

**주요 작업** (13주):
1. **flpkg 패키지 매니저** (6주)
   - CLI: init, add, build, run, test, publish
   - 설정: `flang.toml` (TOML 형식)
   - 레지스트리: Gogs API 연동
   - 컴파일러 자동 선택

2. **VSCode 확장** (3주)
   - 문법 강조: nexus-lexer.ts → TextMate grammar
   - 오류 표시, flpkg 단축키, REPL 패널

3. **온라인 플레이그라운드** (4주)
   - Monaco Editor + FV-WASM
   - 예제 공유 (단축 URL)

**목표**: flpkg 설치 동작, VSCode 확장 배포, 플레이그라운드 런칭

### Q3 (7~9월): 커뮤니티 & 문서

**주요 작업** (13주):
1. **공식 문서 사이트** (6주)
   - freelang.dev
   - 페이지: getting-started, language-guide, stdlib, tutorials, playground, packages
   - stdlib API 자동 생성 (core.json → Markdown)

2. **커뮤니티 채널**
   - Discord (국문), Reddit r/freelang (영문)
   - Gogs Issues/Discussions 포럼화

3. **튜토리얼 4편**
   - "시스템 프로그래밍" (Nexus)
   - "함수형 프로그래밍" (FV-Lang)
   - "FreeLang GPT: Transformer 이해" (Killer app)
   - "한글로 코딩하기" (PyFree, 교육)

**목표**: 문서 30+ 페이지, GitHub Star 100, Discord 활성화

### Q4 (10~12월): 생태계 성숙

**주요 작업** (13주):
1. **공식 패키지 10개**
   - freelang-json, http, regex, crypto, uuid, datetime, path, env, log, test

2. **기여자 프로그램**
   - Good First Issue (입문/중급/고급 3단계)
   - 인센티브: 사이트 등재, Discord 역할, 인터뷰

3. **컨퍼런스 발표**
   - 국내: PyCon, NAVER DevTalk, 카카오 (1회 이상)
   - 주제: "사일로에서 생태계로: FreeLang의 1년"

**목표**: 패키지 10개, 활성 기여자 15명, Star 300

## KPI 트래킹 (측정 가능)

```
분기별 목표:

개발자 지표
  GitHub Stars:      10 → 30 → 100 → 300
  활성 기여자:       1  → 3  → 8   → 15
  패키지 수:         0  → 5  → 15  → 50

기술 지표
  FLIR 채택 컴파일러: 2  → 4  → 6   → 8
  stdlib 함수:       30 → 80 → 150 → 250
  테스트 커버리지:    70% → 80% → 85% → 90%

커뮤니티 지표
  문서 페이지:       10 → 25 → 50 → 80
  Discord 멤버:      10 → 50 → 200 → 500
  블로그 포스트:     3  → 8  → 15 → 25
```

## 포지셔닝

**슬로건**: "한국에서 만든, 세계 수준의 언어 생태계"

**2개 축**:
1. 🇰🇷 한국어 기반 (PyFree 활용, 교육 시장)
2. ⚡ C 성능 + Python 편의

**채널**:
- 국내: GeekNews (월 1회), Okky (주 2회), 카카오 오픈채팅
- 영문: Reddit r/ProgrammingLanguages, Hacker News (Show HN)

## 각 프로젝트의 새로운 역할

| 프로젝트 | 현재 | 생태계 역할 |
|---------|------|----------|
| Nexus | V 모드 | **공식 레퍼런스** |
| FV2-Go | FV 컴파일러 | **네이티브 백엔드** |
| FreeJulia | 자기호스팅 | **FLIR VM** |
| PyFree | 한글 | **교육 진입점** |
| FV-Lang | 함수형 | **최적화 패스** |
| FreeLang-to-C | 트랜스파일러 | **C 백엔드** |
| GPT | Transformer | **stdlib 검증** |
| Multi-Lang PoC | 성능 비교 | **공식 벤치마크** |
| **Ecosystem** | **미시작** | **flpkg 허브** |

## 핵심 파일 (Q1에서 생성될 것)

| 파일 | 예상 크기 | 소유자 |
|------|---------|--------|
| `freelang-ecosystem/specs/flir-v1.0.json` | 500 라인 | FV-Lang |
| `freelang-ecosystem/specs/flir-v1.0.md` | 300 라인 | 생태계 |
| `freelang-nexus/src/nexus/ir/flir-emitter.ts` | 200 라인 | Nexus |
| `fv2-lang-go/internal/flir/reader.go` | 200 라인 | FV2-Go |
| `freelang-ecosystem/stdlib/core.json` | 300 라인 | 생태계 |
| `freelang-ecosystem/stdlib/c/core.h` | 400 라인 | 생태계 |
| `freelang-ecosystem/stdlib/go/core.go` | 400 라인 | 생태계 |

## 중요 포인트

### FLIR 설계의 핵심
- IR은 언어-독립적이어야 함 (어떤 구문도 표현 가능)
- JSON 직렬화로 프로젝트 간 호환성 보장
- 각 백엔드는 자신의 최적화 패스 적용 가능

### stdlib 표준화 방법
- 명세 (core.json) → 레퍼런스 (C/Go/JS) → 각 컴파일러 채택
- 모든 컴파일러가 같은 함수 시그니처로 stdlib 제공

### 패키지 생태계의 출발점
- flpkg (도구) + FLIR (표준) + stdlib (기반) = 완성 시스템
- Q2부터 실제 패키지 개발 가능

## 다음 즉시 액션 (Q1 시작)

1. 📅 FLIR 스펙 작성 회의 (이번 주)
2. 👥 Q1 팀 구성 및 역할 분담
3. 📝 stdlib core.json 함수 목록 최종 확정
4. 🔧 Nexus FLIR emitter 착수

---

**상태**: 계획 완료, Q1 시작 대기
**마지막 갱신**: 2026-03-21
**검토 예정**: Q1 완료 (2026-03-31)
