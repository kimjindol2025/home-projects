# FreeLang 팀모드 시작 + v2.6.0 구현 계획

## Context

**OPTION A-E 모두 완료됨** (2026-03-06):
- OPTION A: test_runner.sh + 논리 검증 완료
- OPTION B: Challenge 15 Sovereign-Naming
- OPTION C: Challenge 14 L0-Mail-Core (1,546줄)
- OPTION D: Sovereign Backend Node.js (40/40 테스트, 8개 규칙 100%)
- OPTION E: V2_6_0_DESIGN.md (394줄, 5개 신기능 설계)

**현재 요청**:
1. **팀모드 시작**: FreeLang AI Marketing Team (5명 에이전트)
2. **v2.6.0 구현 (Week 1)**: ? 연산자 + f-string

---

## OPTION D: Sovereign Backend Node.js 구현

### 저장소: `/data/data/com.termux/files/home/freelang-sovereign-backend/`
기존 FreeLang 설계 파일 1,862줄 참조하여 Node.js로 실제 구현

### 구현 파일

**1. `server.js`** (~300줄)
- Node.js http 모듈, 포트 8080
- Circuit Breaker (CLOSED/OPEN/HALF_OPEN)
- Rate Limiter (Token Bucket)
- Graceful Shutdown

**2. `routes.js`** (~200줄)
- GET  /health
- GET  /metrics
- GET  /api/data
- POST /api/data
- GET  /api/status

**3. `middleware.js`** (~150줄)
- 요청 로깅 + 타이밍
- CORS, Keep-Alive
- 에러 핸들링 (400/404/500)

**4. `metrics.js`** (~100줄)
- P50/P95/P99 지연 시간
- 요청/에러 카운터
- CB 상태 추적

**5. `test_all.sh`** (~150줄)
- curl 테스트 (8가지 시나리오)
- > 100 req/s 성능 검증
- 8개 무관용 규칙 검증

### 무관용 규칙 (8개)
- R1: 시작 < 5초
- R2: GET/POST < 100ms
- R3: 에러 코드 100% 정확
- R4: Circuit Breaker < 100µs
- R5: P95 < 50ms
- R6: > 100 req/s
- R7: 메트릭 100% 정확
- R8: 종료 < 30초

---

## OPTION E: v2.6.0 설계

### 저장소: `/data/data/com.termux/files/home/freelang-final/`
### 파일: `V2_6_0_DESIGN.md` (~200줄)

**5개 신기능 후보**:
1. 모듈 시스템 (`import { fn } from "mod"`)
2. Union 타입 (`type Result = Ok | Err`)
3. try/catch 에러 처리
4. `?` 연산자 (Result 전파)
5. f-string (`f"Hello {name}!"`)

**구현 우선순위**: 타입 강화 → 에러처리 → 모듈 개선

---

---

## 작업 1: 팀모드 시작 (FreeLang AI Marketing Team)

### 배경 (CLAUDE.md 기반)
5명 AI 에이전트 팀이 이미 정의되어 있음. 에이전트 파일 및 지원 구조 생성 필요.

### 생성할 파일 목록

**에이전트 파일** (`.claude/agents/`):
1. `cmo.md` - CMO, opus-4-6, 일요일 21:00
2. `content-writer.md` - Content Writer, sonnet-4-6
3. `social-media.md` - Social Media, haiku-4-5
4. `community-manager.md` - Community Manager, haiku-4-5
5. `analytics.md` - Analytics, haiku-4-5

**에이전트 메모리** (`.claude/agent-memory/`):
- `cmo-memory.md`
- `content-writer-memory.md`
- `social-media-memory.md`
- `community-manager-memory.md`
- `analytics-memory.md`

**규칙 파일** (`rules/`):
- `brand-voice.md` - FreeLang 브랜드 보이스 (기술적, 혁신적, 한국어 우선)
- `content-policy.md` - 콘텐츠 정책 (금지사항 포함)

**로그 파일**:
- `team-log.csv` - 팀 활동 기록 CSV

### 에이전트 파일 내용 (각 파일 구조)
```
# [에이전트명] Agent

## 역할
## 모델
## 실행 주기
## 도구 및 MCP
## 절차 (SOP)
## 메모리 파일 경로
## 기억해야 할 규칙
```

---

## 작업 2: v2.6.0 Week 1 구현

### 대상 저장소
`/data/data/com.termux/files/home/freelang-final/`

### 구현 대상 (Week 1)

#### 1. `?` 연산자 (Result 전파)

**Lexer** (`src/lexer.js`):
- `QUESTION: 'QUESTION'` TokenType 추가
- `?` 문자 처리 (readOperator() 또는 tokenize())

**Parser** (`src/parser.js`):
- `postfix()` 메서드에 `?` 처리 추가
- `QuestionOperator` 또는 `ErrorPropagation` AST 노드 추가

**Evaluator** (`src/evaluator.js`):
- Result 타입 언래핑 처리
- Err이면 즉시 함수 반환

**테스트** (`v2_6_tests.fl` 신규 생성):
- 기본 `?` 연산자 (6개 테스트)

#### 2. f-string

**Lexer** (`src/lexer.js`):
- `f"..."` 시작 감지
- 보간 `{expr}` 토큰화

**Parser** (`src/parser.js`):
- `FStringLiteral` AST 노드
- 내부 표현식 파싱

**Evaluator** (`src/evaluator.js`):
- 각 보간 표현식 평가
- 문자열 조합

**테스트** (`v2_6_tests.fl`):
- 기본 보간 (6개 테스트)

### 핵심 파일 (기존 패턴 활용)
- Lexer: `src/lexer.js` (494줄) - TokenType + KEYWORDS 맵 패턴
- Parser: `src/parser.js` (791줄) - ASTNode + statement() 분기 패턴
- Evaluator: `src/evaluator.js` (466줄) - eval() if-else 체인 패턴

---

## 구현 순서

### Phase 1: 팀모드 (30분)
1. `.claude/agents/` 5개 파일 생성
2. `.claude/agent-memory/` 5개 메모리 파일 생성
3. `rules/brand-voice.md`, `rules/content-policy.md` 생성
4. `team-log.csv` 초기화
5. 팀 첫 활동 실행 (Content Writer 블로그 초안)

### Phase 2: v2.6.0 Week 1 (1-2시간)
1. Lexer에 `?` 토큰 + f-string 토큰 추가
2. Parser에 새 AST 노드 + 파싱 메서드 추가
3. Evaluator에 실행 로직 추가
4. `v2_6_tests.fl` 테스트 파일 생성 (12개 테스트)
5. 테스트 실행 확인
6. GOGS 커밋

## 예상 결과
- 5명 에이전트 팀 활성화 (팀모드)
- v2.6.0 Week 1: ~300줄 추가, 12개 테스트 통과
- GOGS 커밋 2회
