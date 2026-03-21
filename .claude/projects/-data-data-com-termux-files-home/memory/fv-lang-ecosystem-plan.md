---
name: FV-Lang 스탠드얼론 생태계 조성계획 v1.0
description: Q1-Q4 4분기 로드맵, fvpkg 매니저, 문서 사이트, 커뮤니티 구축 계획
type: project
---

# FV-Lang 스탠드얼론 생태계 조성계획 v1.0

**상태**: ✅ APPROVED & IN EFFECT (2026-03-21)
**파일**: `projects/fv-lang/FV_ECOSYSTEM_PLAN.md` (350줄)

## 핵심 전략

FV-Lang (Rust, 57% 완성, 7,703줄)을 FreeLang Nexus와 별도로 **완전한 독립 생태계**로 구축.

### 3가지 핵심 구축 항목

1. **fvpkg 패키지 매니저** (Rust 구현)
   - `fvpkg init/add/build/run/publish` 명령어
   - `fv.toml` 설정 파일 (TOML 형식)
   - Gogs API 연동 레지스트리

2. **공식 문서 사이트** (mdBook)
   - 15페이지 (getting-started ~ stdlib reference)
   - GitHub Pages 자동 배포
   - 튜토리얼 4편 (입문/함수/패턴매칭/fvpkg)

3. **커뮤니티** (Discord + Reddit)
   - Discord 서버 (한글 메인)
   - Reddit r/fvlang (영어)
   - Good First Issue 3단계 기여자 프로그램

## 4분기 KPI

```
Q1 (3월):    기반 구축
  ├─ fvpkg 초안 (5개 명령어)
  ├─ stdlib 4개 패키지 (25개 함수)
  └─ 문서 6페이지
  → 목표: GitHub Stars 5-10, 활성 기여자 1

Q2 (4-6월):  도구체인 완성
  ├─ fvpkg 완성 (8개 명령어)
  ├─ 문서 사이트 (15페이지)
  ├─ 튜토리얼 4편
  └─ crates.io v0.1.0 배포
  → 목표: GitHub Stars 20, crates.io 100+ DL

Q3 (7-9월):  커뮤니티 구축
  ├─ Discord 활성화 (50명)
  ├─ stdlib 6개 추가 (50개 함수)
  ├─ Good First Issue 프로그램
  └─ 튜토리얼 심화 (7편)
  → 목표: GitHub Stars 60, Discord 50명

Q4 (10-12월): 성숙
  ├─ 패키지 10개 완성 (150개 함수)
  ├─ 언어 95% 완성
  ├─ 컨퍼런스 발표 3회
  └─ 활성 기여자 12명
  → 목표: GitHub Stars 150, crates.io 5,000+ DL
```

## 시작 체크리스트

### 즉시 (이번 주)
- [ ] `projects/fv-lang/FV_ECOSYSTEM_PLAN.md` ✅ 생성 완료
- [ ] `packages/` 디렉토리 구조 생성 준비

### Q1 실행 (3월)
- [ ] fvpkg 초안 (`fvpkg/src/main.rs` + `Cargo.toml`)
- [ ] stdlib 4개 패키지 (각 50줄 FV 코드)
- [ ] docs/ 기초 (6개 Markdown)
- [ ] GitHub 설정 (Pages 활성화)

## 주요 파일 참조

| 파일 | 역할 |
|-----|------|
| `projects/fv-lang/FV_ECOSYSTEM_PLAN.md` | 메인 계획 문서 |
| `projects/fv-lang/src/stdlib.rs` | stdlib 기반 (591줄) |
| `projects/fv-lang/Cargo.toml` | 패키지 메타데이터 |
| `projects/fv-lang/examples/` | 예제 코드 |

## 차별점

1. **FreeLang과 독립** - Nexus와 별개의 생태계
2. **순수 함수형** - 부작용 없는 언어 설계
3. **성능 중심** - Rust 수준의 성능 보장
4. **한글 우선** - 문서/커뮤니티 한글 중심 (영어 이중언어)

---

**상태**: APPROVED & IN EFFECT (2026-03-21)
**다음 검토**: Q1 완료 (2026-03-31)
