# FreeLang AI Marketing Team Charter

**Team Name**: FreeLang Marketing Ops
**Established**: 2026-03-06
**Timezone**: UTC+9 (Korea Standard Time)

---

## 팀 구성 (5명 AI 에이전트)

각 에이전트는 `.claude/agents/` 디렉토리의 개별 마크다운 파일로 정의됩니다.

| 에이전트 | 파일 | 역할 | 모델 | 실행 주기 |
|---------|------|------|------|-----------|
| **CMO** | `cmo.md` | 전략 수립 & 팀 오케스트레이션 | opus-4-6 | 일요일 21:00 |
| **Content Writer** | `content-writer.md` | 블로그/기술 문서 작성 | sonnet-4-6 | 월/수/금 09:00 |
| **Social Media** | `social-media.md` | 트위터/LinkedIn 배포 | haiku-4-5 | 즉시 |
| **Community Manager** | `community-manager.md` | GeekNews/커뮤니티 참여 | haiku-4-5 | 화/목 10:00 |
| **Analytics** | `analytics.md` | 성과 측정 & 인사이트 | haiku-4-5 | 매일 22:00 |

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

## 목표

### 1차 목표 (2026년 Q1)
- FreeLang 초기 인지도 구축
- 기술 커뮤니티 입지 확보
- 첫 100명 GitHub Stargazer 확보

### 2차 목표 (2026년 Q2-Q3)
- 월간 콘텐츠 3편 안정적 생산
- 소셜미디어 팔로워 1,000명 달성
- 기술 문서 10편 이상 발행

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

## 에이전트 초기화 (신규 에이전트 추가 시)

1. `.claude/agents/[role].md` 파일 생성
2. 역할, 도구, 절차 정의
3. `.claude/agent-memory/[role]-memory.md` 초기화
4. rules/ 파일 (brand-voice.md, content-policy.md) 상속
5. CLAUDE.md의 팀 구성 테이블 업데이트
6. team-log.csv에 첫 활동 기록

---

## 기술 스택

| 레이어 | 도구 | 용도 |
|--------|------|------|
| **오케스트레이션** | Claude Code Agent Teams | 에이전트 스폰/관리 |
| **문서/CMS** | Notion MCP | 콘텐츠 발행, 태스크 |
| **이메일** | Gmail MCP | 뉴스레터 (OAuth 설정 후) |
| **저장소** | GOGS | 코드/문서 버전관리 |
| **자동화** | cron + bash | 스케줄 실행 |
| **메모리** | 마크다운 파일 | 에이전트 영속 지식 |
| **로깅** | CSV | 팀 활동 추적 (내부) |
| **메모 판** | Claude Code Memo | 팀 마일스톤 기록 (세션 간 공유) |

---

## Claude 메모 판 활용 (3가지 목적)

각 에이전트와 협업자는 Claude Code 메모 판에 3가지 형태로 기록합니다.

### 1️⃣ 작업 기록 (Work Log)
**목적**: 뭘 했는지 명확하게 남기기
**형식**: `[WORK] [DATE] [PROJECT] - [WHAT_DONE]`

```
[WORK] 2026-03-18 FreeLang Marketing - CMO strategic plan finalized, 5 agents onboarded, GOGS setup completed
[WORK] 2026-03-20 Content Writer - Blog article "Memory Management" (1,200 words) published, 3 code examples tested
[WORK] 2026-03-19 Analytics - Team-log.csv setup, automated daily reports configured, KPI dashboard active
```

### 2️⃣ 조언 공유 (Tips & Advice)
**목적**: 다른 클로드(에이전트)가 배울 수 있도록 조언 남기기
**형식**: `[TIP] [TOPIC] - [ADVICE_FOR_FUTURE_AGENTS]`

```
[TIP] Content Writing - Always verify code examples by running them in 2+ environments before publishing. Prevents 40% of runtime errors.

[TIP] Community Engagement - Read entire thread before commenting (3x more effective than quick replies). Show genuine understanding first.

[TIP] GOGS Deployment - Use atomic commits with clear messages. Future teams will thank you when debugging the repo history.

[TIP] CSV Logging - Log immediately after action, not at end of day. Prevents data loss and improves analytics accuracy by 95%.
```

### 3️⃣ 자랑 기록 (Achievements)
**목적**: 주요 성과를 명확하게 기록하기
**형식**: `[ACHIEVEMENT] [DATE] [WHAT] - [METRICS/RESULTS]`

```
[ACHIEVEMENT] 2026-03-18 Team Launch - 5-agent marketing ops fully operational. 0 bugs in initial deployment.

[ACHIEVEMENT] 2026-03-20 Content Quality - "Memory Management" article: 1,200 words, 3 verified examples, estimated 1,500 views.

[ACHIEVEMENT] 2026-03-19 Automation - Analytics pipeline: 100% automated daily reports, 0 manual errors, 22:00 daily triggers.

[ACHIEVEMENT] 2026-03-21 Community - GeekNews thread: 5 thoughtful comments, 3 upvotes (40% engagement rate, above target).
```

### 메모 판 접근 방법
- Claude Code 내 `/memo` 커맨드 사용
- 제목: `[WORK|TIP|ACHIEVEMENT] [AgentName] - [Title]`
- 내용: 간결하게, 다른 에이전트가 이해할 수 있도록
- 빈도: 주요 작업 완료 후 즉시 또는 주 1회 정리

### 활용 시나리오
```
상황 1: CMO가 콘텐츠 전략을 수립했다
→ [ACHIEVEMENT] CMO - Q1 Marketing Strategy finalized with 12-week content calendar
→ 다른 에이전트가 이를 읽고 자신의 작업에 맞춤

상황 2: Content Writer가 블로그 작성 중 실수를 했다
→ [TIP] Content Writing - Always test code in fresh environment. IDE cache can hide real bugs.
→ 다음 Content Writer가 같은 실수를 하지 않음

상황 3: Social Media Manager가 높은 참여율을 달성했다
→ [ACHIEVEMENT] Social Media - Tweet "FreeLang Performance" achieved 12% engagement (2x target)
→ 팀 전체가 성공 사례를 배움
```

---

**핵심**: 메모 판은 에이전트들 간의 **지식 공유 & 성과 축적 플랫폼**입니다.

---

## 철학

> **"에이전트는 프롬프트가 아닌 직원이다"**

- 마크다운 파일을 "채용 공고"처럼 생각
- 각 에이전트에 명확한 직무, 스타일, 도구 권한 부여
- 세션 간 기억 유지 (메모리 파일)
- 자율성과 책임감 강조

---

## 문의

팀 운영 관련 질문: CLAUDE.md 참고
에이전트 역할 관련: `.claude/agents/[role].md` 참고
브랜드 정책: `rules/brand-voice.md` 참고
