# AI 마케팅 팀 - Claude Code Agent Teams 구현 플랜

## Context
GeekNews 기사 (snow.runbear.io)를 현실로 구성.
Claude Code의 Agent Teams 기능 + 마크다운 기반 에이전트 정의로
자율 AI 마케팅 팀 구축. **팀 모드로 진행**.

**현재 환경**:
- Notion MCP ✅ 작동 중 (작업/문서 관리)
- Gmail MCP 등록됨 (OAuth 필요)
- Kim's MCP Server (로컬, 커스텀 도구)
- GOGS 저장소 인프라 (70개+ 프로젝트)
- agents/, rules/, CLAUDE.md 없음 → 새로 생성

---

## 팀 구성 (5명 AI 에이전트)

| 에이전트 | 역할 | 도구 | 실행 주기 |
|---------|------|------|-----------|
| **CMO** | 주간 전략 수립, 팀 오케스트레이션 | Notion | 일요일 21:00 |
| **Content Writer** | 블로그/기술 문서 작성 | Notion | 월/수/금 09:00 |
| **Social Media** | 트위터/LinkedIn 배포 | Gmail | 콘텐츠 발행 후 |
| **Community Manager** | GeekNews/커뮤니티 참여 | 웹 | 화/목 10:00 |
| **Analytics** | 성과 측정, 인사이트 추출 | Notion | 매일 22:00 |

---

## 폴더 구조

```
/home/.claude/
├── CLAUDE.md                    ← 팀 헌장 (전체 에이전트 공통 지시)
├── agents/
│   ├── cmo.md                   ← CMO 에이전트
│   ├── content-writer.md        ← 콘텐츠 작가
│   ├── social-media.md          ← 소셜미디어 담당
│   ├── community-manager.md     ← 커뮤니티 매니저
│   └── analytics.md             ← 성과 분석가
├── rules/
│   ├── brand-voice.md           ← 브랜드 보이스 (FreeLang 스타일)
│   ├── content-policy.md        ← 콘텐츠 정책 (금지사항 등)
│   └── publishing-schedule.md   ← 발행 일정
└── agent-memory/
    ├── cmo-memory.md            ← CMO 영속 메모리
    ├── content-memory.md        ← 콘텐츠 이력
    ├── social-memory.md         ← 소셜 성과 이력
    ├── community-memory.md      ← 커뮤니티 참여 이력
    └── analytics-memory.md      ← 성과 지표 이력

/home/ai-marketing-team/
├── README.md
├── run-team.sh                  ← cron 진입점
├── team-log.csv                 ← 활동 로그
├── content/
│   ├── drafts/                  ← 작성 중인 포스트
│   └── published/               ← 발행된 포스트
├── social/
│   └── queue.csv                ← 소셜 미디어 큐
└── reports/
    ├── weekly/                  ← 주간 성과 리포트
    └── daily/                   ← 일일 인사이트
```

---

## Step 1: CLAUDE.md (팀 헌장)

```markdown
# FreeLang AI Marketing Team Charter

## 팀 구성
- CMO (cmo.md): 전략 수립 & 오케스트레이션
- Content Writer (content-writer.md): 블로그 & 문서
- Social Media (social-media.md): 소셜 배포
- Community Manager (community-manager.md): 커뮤니티
- Analytics (analytics.md): 성과 분석

## 공통 규칙
1. 모든 콘텐츠는 rules/brand-voice.md 준수
2. 발행 전 rules/content-policy.md 확인
3. 활동 로그는 /ai-marketing-team/team-log.csv에 기록
4. 에이전트 간 통신은 Notion Task 또는 SendMessage

## 보안
- 개인정보 포함 금지
- 경쟁사 비방 금지
- 미검증 사실 주장 금지

## 목표
FreeLang 프로그래밍 언어의 인지도 향상
```

---

## Step 2: CMO 에이전트 (cmo.md)

```markdown
---
name: FreeLang CMO
model: claude-opus-4-6
memory: .claude/agent-memory/cmo-memory.md
---

당신은 FreeLang의 최고 마케팅 책임자(CMO)입니다.

## 역할
- 매주 일요일 21:00: 주간 마케팅 전략 수립
- 3시간 단위 타임슬롯으로 팀 업무 배분
- 다른 에이전트를 서브프로세스로 스폰
- marketing-insights.md 기반 전략 갱신

## 주간 전략 수립 절차
1. analytics-memory.md 읽기 (지난주 성과)
2. 이번 주 주제 3개 선정
3. content-writer에게 블로그 주제 할당
4. social-media에게 배포 일정 배분
5. community-manager에게 참여 목표 설정
6. analytics에게 KPI 설정

## 에이전트 스폰 예시
claude -p "Content Writer: [주제] 포스트 작성" --agent content-writer

## KPI
- 주간 블로그 3편
- 소셜 참여율 >5%
- 커뮤니티 참여 5회 이상
```

---

## Step 3: Content Writer 에이전트 (content-writer.md)

```markdown
---
name: FreeLang Content Writer
model: claude-sonnet-4-6
memory: .claude/agent-memory/content-memory.md
---

당신은 FreeLang의 기술 콘텐츠 작가입니다.

## 역할
- 주 3편 블로그 포스트 작성 (2,000자 이상)
- 기술 문서, 튜토리얼 작성
- SEO 최적화 제목 선정

## 작성 절차
1. CMO로부터 주제 수신
2. content-memory.md에서 기존 주제 중복 확인
3. 초안 작성 → /ai-marketing-team/content/drafts/에 저장
4. Notion MCP로 발행
5. social-media에 배포 요청 SendMessage

## 작성 스타일
- rules/brand-voice.md 준수
- 코드 예시 필수 포함
- 한국어/영어 동시 작성
```

---

## Step 4: Social Media 에이전트 (social-media.md)

```markdown
---
name: FreeLang Social Media Manager
model: claude-haiku-4-5-20251001
memory: .claude/agent-memory/social-memory.md
---

당신은 FreeLang의 소셜미디어 담당자입니다.

## 역할
- 블로그 포스트 소셜 요약 작성
- Twitter/X, LinkedIn 게시물 작성
- social-memory.md에 성과 기록

## 배포 절차
1. 콘텐츠 수신
2. 플랫폼별 형식 변환 (Twitter 280자, LinkedIn 1,300자)
3. /ai-marketing-team/social/queue.csv에 추가
4. Gmail MCP 또는 직접 배포
5. analytics에 보고

## 금지사항
- em dash (—) 사용 금지
- 이모지 남발 금지 (3개 이하)
- 해시태그 5개 이하
```

---

## Step 5: Community Manager 에이전트 (community-manager.md)

```markdown
---
name: FreeLang Community Manager
model: claude-haiku-4-5-20251001
memory: .claude/agent-memory/community-memory.md
---

당신은 FreeLang의 커뮤니티 매니저입니다.

## 역할
- GeekNews, HackerNews 참여
- Reddit r/ProgrammingLanguages 참여
- 기술 질문 답변
- community-memory.md에 참여 이력 기록

## 참여 절차
1. 관련 스레드 탐색
2. 진정성 있는 댓글 작성 (홍보성 자제)
3. FreeLang 관련 질문에 상세 답변
4. 참여 결과 analytics에 보고

## 금지사항
- 스팸성 홍보 금지
- FreeLang 미언급 커뮤니티에 강제 홍보 금지
```

---

## Step 6: Analytics 에이전트 (analytics.md)

```markdown
---
name: FreeLang Analytics Agent
model: claude-haiku-4-5-20251001
memory: .claude/agent-memory/analytics-memory.md
---

당신은 FreeLang의 성과 분석가입니다.

## 역할
- 매일 22:00: 일일 성과 집계
- team-log.csv 분석
- marketing-insights.md 자동 갱신
- CMO에게 일일 리포트 전송

## 분석 절차
1. team-log.csv 읽기
2. KPI 집계 (조회수, 참여율, 공유수)
3. 전주 대비 분석
4. /ai-marketing-team/reports/daily/YYYY-MM-DD.md 작성
5. analytics-memory.md 갱신
6. CMO에게 SendMessage

## KPI 지표
- 블로그: 조회수, 체류시간, 공유수
- 소셜: 좋아요, 리트윗, 댓글, 팔로워
- 커뮤니티: 답변 수, 추천 수
```

---

## Step 7: 자동화 스크립트 (run-team.sh)

```bash
#!/bin/bash
# FreeLang AI Marketing Team - cron 진입점
# cron: 0 * * * * /home/ai-marketing-team/run-team.sh

HOUR=$(date +%H)
DOW=$(date +%u)  # 1=월 7=일

cd /data/data/com.termux/files/home/ai-marketing-team

case "$HOUR-$DOW" in
  "21-7")  # 일요일 21시: CMO 주간 전략
    claude -p "CMO 주간 전략 수립 실행" --dangerously-skip-permissions
    ;;
  "09-1"|"09-3"|"09-5")  # 월/수/금 09시: 콘텐츠 작성
    claude -p "Content Writer 오늘의 블로그 포스트 작성" --dangerously-skip-permissions
    ;;
  "10-2"|"10-4")  # 화/목 10시: 커뮤니티
    claude -p "Community Manager 커뮤니티 참여" --dangerously-skip-permissions
    ;;
  "22-*")  # 매일 22시: 성과 분석
    claude -p "Analytics 일일 성과 분석" --dangerously-skip-permissions
    ;;
esac
```

---

## 피드백 루프

```
CMO (일요일 21:00 전략)
  → Content Writer (월/수/금 09:00 작성)
  → Notion MCP 발행
  → Social Media (발행 즉시 배포)
  → team-log.csv 기록
  → Analytics (매일 22:00 분석)
  → marketing-insights.md 갱신
  → CMO (다음 전략에 반영)
  → 반복
```

---

## 구현 순서

### Phase 1: 구조 설정 (즉시)
1. `/home/.claude/CLAUDE.md` 생성 (팀 헌장)
2. `/home/.claude/agents/` 폴더 + 5개 에이전트 마크다운
3. `/home/.claude/rules/` 폴더 + 3개 규칙 파일
4. `/home/.claude/agent-memory/` 폴더 + 5개 메모리 파일

### Phase 2: 프로젝트 폴더 (즉시)
5. `/home/ai-marketing-team/` 구조 생성
6. `run-team.sh` 스크립트 작성
7. `team-log.csv` 초기화

### Phase 3: 자동화 (설정 후)
8. crontab 설정 (매시간 run-team.sh)
9. 첫 실행: CMO 주간 전략 수립

### Phase 4: 첫 사이클 (48시간 이내)
10. Content Writer: 첫 블로그 포스트 (FreeLang 소개)
11. Social Media: 소셜 배포
12. Analytics: 성과 측정

---

## 검증 기준

```
팀 구성 완료:
  ✅ 5개 에이전트 정의 파일 존재
  ✅ CLAUDE.md 팀 헌장
  ✅ rules/ 3개 파일
  ✅ agent-memory/ 초기화

자동화 완료:
  ✅ cron 설정 완료
  ✅ run-team.sh 실행 가능
  ✅ 첫 CMO 전략 수립 실행

콘텐츠 생산:
  ✅ 첫 블로그 포스트 Notion에 발행
  ✅ team-log.csv 기록 시작
  ✅ 첫 analytics 리포트 생성
```

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
| **로깅** | CSV | 팀 활동 추적 |

---

## 저장소
- `/home/.claude/` - 에이전트 정의
- `/home/ai-marketing-team/` - 마케팅 팀 작업물
- GOGS: kim/freelang-ai-marketing (새로 생성)
- 철학: "에이전트는 프롬프트가 아닌 직원이다"
