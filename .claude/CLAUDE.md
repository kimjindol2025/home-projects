# FreeLang AI Marketing Team Charter

**Team Name**: FreeLang Marketing Ops
**Established**: 2026-03-06
**Timezone**: UTC+9 (Korea Standard Time)

---


## 공통 규칙

### 1. 콘텐츠 가이드라인
- 모든 콘텐츠는 `rules/brand-voice.md` 준수 (필수)
- 발행 전 `rules/content-policy.md` 확인 (필수)
- 한국어 우선, 영어 이중언어 작성 권장

### 2. 에이전트 간 통신
- **Notion Task**: 장기 협업 작업 (블로그 → 소셜 배포)
- **SendMessage**: 긴급 알림, 일일 리포트
- **CSV 파일**: 활동 로그 (`/ai-marketing-team/team-log.csv`)
- **Claude 메모 판**: 팀 활동 기록 및 공유 (claude-code memo system)

### 3. 메모리 관리
- 각 에이전트는 `.claude/agent-memory/`의 마크다운 파일 유지
- 메모리 파일: 성과, 학습, 의사결정 기록
- 세션 간 자동 상속 (persistent knowledge)

### 4. 활동 로깅
- **필수**: 모든 에이전트는 활동을 `team-log.csv`에 기록
- 형식: `[시간],[에이전트],[활동],[결과],[KPI]`
- Analytics가 매일 22:00에 수집 & 분석
- **메모 판**: 중요 마일스톤은 Claude 메모 판에 기록
- 메모 형식: `[AGENT] [PROJECT] [STATUS] - [BRIEF_DESC] | [DATE]`
- 예시: `[CMO] FreeLang Marketing Q1 - 5-agent team initialized, GOGS deployment ready | 2026-03-18`

### 5. 플랜모드 기본 적용
- **규칙**: 파일 수정/생성 코드가 100줄 이상인 경우 자동으로 플랜모드 진입
- **목적**: 큰 변경사항에 대한 사전 검토 및 승인
- **절차**: 사용자가 승인할 때까지 코드 작성 미연기
- **예외**: 버그픽스나 간단한 수정은 플랜모드 스킵 가능

### 6. FreeLang 코딩 규칙
- **기준**: FreeLang v2 사용 ([전체 버전 자원](./FREELANG_ALL_VERSIONS.md))
- **원칙**: 모든 코딩은 FreeLang 적극 사용
- **절차**: 부족한 기능은 먼저 보고 후 작성
- **우선순위**: FreeLang으로 구현 가능한지 먼저 검토

---

## 보안 정책

### 금지사항
- ❌ 개인정보(주민번호, 전화번호 등) 포함 금지
- ❌ 경쟁사 비방 금지
- ❌ 미검증 사실 주장 금지 (근거 제시 필수)
- ❌ 스팸/광고성 콘텐츠 금지

### 필수 검증
- 모든 기술 주장은 코드 예시 또는 공식 문서로 뒷받침
- 성능 지표는 실제 벤치마크 제시
- 커뮤니티 참여는 진정성 있는 기여 (홍보 자제)

---

## 성공 지표 (KPI)

| 지표 | 목표 | 측정 주기 |
|------|------|----------|
| 블로그 조회수 | 월 1,000+ | 일일 |
| 소셜 참여율 | >5% | 일일 |
| 커뮤니티 댓글 | 주 5회 이상 | 주간 |
| 콘텐츠 생산 | 주 3편 | 주간 |

---

## 피드백 루프

```
┌─────────────────────────────────────────┐
│ CMO (일요일 21:00)                      │
│ 주간 전략 수립 & 작업 배분               │
└──────────────┬──────────────────────────┘
               │
        ┌──────▼──────┐
        │ Content     │ (월/수/금 09:00)
        │ Writer      │ 블로그 작성
        └──────┬──────┘
               │
        ┌──────▼──────────┐
        │ Notion MCP      │ 자동 발행
        └──────┬──────────┘
               │
        ┌──────▼──────┐      ┌────────────────┐
        │ Social      │      │ Community Mgr  │
        │ Media       │      │ (화/목 10:00)  │
        │ (즉시)      │      │ 커뮤니티 참여   │
        └──────┬──────┘      └────────┬───────┘
               │                     │
        ┌──────┴──────────────────┬──┘
        │ team-log.csv 기록       │
        │ (에이전트별)             │
        └──────┬──────────────────┘
               │
        ┌──────▼──────┐
        │ Analytics   │ (매일 22:00)
        │ 성과 분석    │
        └──────┬──────┘
               │
        ┌──────▼──────────────────┐
        │ marketing-insights.md   │
        │ CMO가 읽고 전략 갱신   │
        └─────────────────────────┘
```

---


## 철학

> **"에이전트는 프롬프트가 아닌 직원이다"**

- 마크다운 파일을 "채용 공고"처럼 생각
- 각 에이전트에 명확한 직무, 스타일, 도구 권한 부여
- 세션 간 기억 유지 (메모리 파일)
- 자율성과 책임감 강조

---

## 참고 자료

| 자료 | 설명 | 링크 |
|------|------|------|
| **FreeLang 검증 보고서** | 실제 코드 분석 (완성도 68%) | [FREELANG_VERIFICATION_REPORT.md](./FREELANG_VERIFICATION_REPORT.md) |
| **FreeLang 언어 분석** | 기능, 성능, 타입시스템 평가 | [FREELANG_LANGUAGE_ANALYSIS.md](./FREELANG_LANGUAGE_ANALYSIS.md) |
| **FreeLang 전체 버전** | v2~v6 자원 맵 | [FREELANG_ALL_VERSIONS.md](./FREELANG_ALL_VERSIONS.md) |
| **v2 자원** | v2 특화 가이드 | [FREELANG_V2_RESOURCES.md](./FREELANG_V2_RESOURCES.md) |

---

## 문의

| 주제 | 참고 자료 |
|------|----------|
| 팀 운영 | CLAUDE.md |
| FreeLang 분석 | [FREELANG_LANGUAGE_ANALYSIS.md](./FREELANG_LANGUAGE_ANALYSIS.md) |
| FreeLang 자원 | [FREELANG_ALL_VERSIONS.md](./FREELANG_ALL_VERSIONS.md) |
| 에이전트 역할 | `.claude/agents/[role].md` |
| 브랜드 정책 | `rules/brand-voice.md` |
