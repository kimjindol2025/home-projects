---
name: FreeLang 서버 기반 전수조사
description: 프리랭 생태계의 모든 서버/API 프로젝트 분석 보고서
type: project
---

# 🔍 FreeLang 서버 기반 전수조사 보고서
**작성일**: 2026-03-25 | **상태**: 완료 ✅

---

## 📊 전체 요약

| 범주 | 프로젝트 수 | 상태 | 언어 | 설명 |
|------|-----------|------|------|------|
| **REST API** | 1 | 🟡 설계중 | FreeLang | Todo REST API 프레임워크 |
| **Backend** | 1 | 🟡 개발중 | FreeLang | 다층 아키텍처 (Auth/Cache/DB) |
| **Mail Server** | 1 | 🟡 설계중 | FreeLang | 이메일 서버 구현 |
| **분산 시스템** | 1 | 🟡 설계중 | FreeLang | Raft 기반 분산DB |
| **Node.js** | 1 | 🟢 운영중 | JavaScript | FreeWire API (그래프 컴파일러) |
| **Go LLM** | 1 | 🟢 운영중 | Go | FreeLang GPT REST API |

**총 6개 서버 프로젝트** (FreeLang 5, JavaScript 1)

---

## 🟡 설계/개발중인 FreeLang 서버들

### 1️⃣ **freelang-rest-api** (REST API Framework)
- **위치**: `.projects/core/freelang-rest-api`
- **상태**: 🟡 Phase B 설계중
- **구현 언어**: FreeLang (.fl)
- **목표**: RESTful API 개발 프레임워크

**구조**:
```
src/
├── main.fl          — 라우터 + 디스패칭 엔진 (200줄+)
├── models.fl        — Todo 모델 + 저장소
├── handlers/
│   └── todo.fl      — 5개 엔드포인트
└── errors.fl        — 에러 타입 정의

tests/
└── test_todo.fl     — 단위 테스트
```

**구현 현황**:
- ✅ 라우트 정의 & 매칭 (패턴 기반)
- ✅ 요청 디스패칭
- ✅ Todo 모델 (CRUD)
- ✅ 에러 처리 (ApiError enum)
- ⏳ HTTP 서버 구현 (Phase B 런타임 대기)

**엔드포인트** (설계):
```
GET    /todos         — 모든 Todo 조회
POST   /todos         — Todo 생성
GET    /todos/:id     — ID로 조회
PUT    /todos/:id     — 업데이트
DELETE /todos/:id     — 삭제
```

**특징**:
- 함수형 프로그래밍 기반 라우팅
- 타입 안전성 (Result<T, E>)
- Immutable 저장소 패턴

---

### 2️⃣ **freelang-backend-system** (다층 아키텍처)
- **위치**: `.projects/core/freelang-backend-system`
- **상태**: 🟡 개발중
- **구현 언어**: FreeLang (.fl)
- **목표**: 엔터프라이즈급 백엔드 시스템

**6계층 아키텍처**:
```
1. 프레젠테이션 (Presentation)
   └─ src/layers/presentation.fl — HTTP 응답 변환

2. 비즈니스 로직 (Business)
   └─ src/layers/business.fl     — 핵심 로직

3. 인증/인가 (Auth)
   ├─ src/auth/jwt.fl            — JWT 토큰
   ├─ src/auth/user.fl           — 사용자 관리
   └─ src/auth/rbac.fl           — 역할 기반 접근

4. 캐싱 (Caching)
   └─ src/caching/cache.fl       — 메모리 캐시

5. 지속성 (Persistence)
   ├─ src/layers/persistence.fl  — DB 추상화
   ├─ src/database.fl            — 쿼리 엔진
   └─ src/optimization/indexing.fl — 인덱싱

6. 모니터링 (Monitoring)
   ├─ src/monitoring/metrics.fl  — 메트릭 수집
   └─ src/monitoring/profiler.fl — 성능 분석

부가:
├─ src/grpc/                    — gRPC 서비스
├─ src/runtime/websocket.fl     — WebSocket 지원
└─ src/serialization/protobuf.fl — Protobuf 직렬화
```

**구현 현황**:
- ✅ 모델 정의 (src/models.fl)
- ✅ 인증 시스템 (JWT + RBAC)
- ✅ 캐싱 전략
- ✅ gRPC 메시지/핸들러
- ⏳ HTTP 요청/응답 변환
- ⏳ 데이터베이스 연동

**특징**:
- 관심사의 분리 (각 계층 독립)
- 플러그인 아키텍처 (Middleware)
- 비동기 작업 큐
- 타입 안전 JSON 직렬화

---

### 3️⃣ **freelang-mail-server**
- **위치**: `.projects/core/freelang-mail-server`
- **상태**: 🟡 설계 단계
- **구현 언어**: FreeLang (.fl)
- **목표**: SMTP/IMAP 이메일 서버

**예상 구조**:
```
SMTP Handler      — 메일 수신
IMAP Handler      — 메일 조회/관리
Storage Engine    — 파일 기반 저장
Encryption        — TLS/SSL 암호화
Queue System      — 배달 재시도
```

**현재 상태**: 초기 설계 단계

---

### 4️⃣ **freelang-distributed-system** (Raft 기반 분산DB)
- **위치**: `.projects/core/freelang-distributed-system`
- **상태**: 🟡 개발중
- **구현 언어**: FreeLang (.fl)
- **목표**: Raft 합의 기반 분산 데이터베이스

**예상 기능**:
- Raft 합의 알고리즘
- 분산 트랜잭션
- 상태 복제 (State Machine Replication)
- 장애 조치 (Failover)

**현재 상태**: 핵심 알고리즘 구현중

---

## 🟢 운영중인 서버들

### 5️⃣ **FreeWire API** (Node.js)
- **위치**: `projects/FreeWire`
- **상태**: 🟢 운영중
- **포트**: 4010
- **언어**: JavaScript (Node.js)
- **목표**: 그래프 컴파일러 API

**엔드포인트**:
```
GET  /health                    — 상태 확인
GET  /api/nodes                 — 노드 타입 목록
POST /api/compile               — Graph → .free 코드
POST /api/execute               — Graph → HTML/CSS/JS
POST /api/preview               — 미리보기 생성 (캐시)
POST /api/ai-generate           — AI → Graph JSON
POST /api/evaluate              — Graph 의도 평가 (0-100점)
POST /api/ir                    — Graph → FIR IR
GET  /preview/:id               — 미리보기 렌더링
```

**주요 기능**:
- 📦 **미리보기 캐시**: LRU (100개, 1시간 TTL)
- 🔒 **CSP 보안**: nonce 기반 스크립트 샌드박스
- 🎨 **다중 포맷**: HTML/CSS/JS/IR 지원
- 🤖 **AI 통합**: Claude API → 그래프 생성

**구현 사항**:
```javascript
src/api/server.js           — HTTP 라우터 & 핸들러 (386줄)
src/api/ai-graph-generator.js — Claude API 호출
src/api/free-executor.js    — .free 코드 실행기
src/compiler/graph-compiler.js — Graph → .free 컴파일
src/intent/evaluator.js     — 의도 평가 엔진
```

**원칙**: 모든 UI는 FreeLang으로 생성 (직접 HTML 금지)

---

### 6️⃣ **FreeLang GPT** (Go REST API)
- **위치**: `projects/freelang-gpt/api`
- **상태**: 🟢 운영중
- **포트**: 8080
- **언어**: Go

**핵심 서버들**:

#### **server.go** (기본 REST API)
```
GET  /health         — 상태 확인
GET  /models         — 모델 정보
POST /api/generate   — 텍스트 생성
POST /api/encode     — 텍스트 → 토큰
POST /api/decode     — 토큰 → 텍스트
GET  /api/stats      — 모델 통계
GET  /docs           — API 문서 (HTML)
```

**응답 형식**:
```json
{
  "prompt": "FreeLang은",
  "generated": " 함수형 프로그래밍 언어입니다.",
  "full_text": "FreeLang은 함수형 프로그래밍 언어입니다.",
  "model": "simple",
  "tokens": 7,
  "time_ms": "0ms"
}
```

#### **fl8_gpt_server.go** (FL8 추론 서버)
```
GET  /api/fl8/status       — FL8 런타임 상태
POST /api/fl8/generate     — FL8 GPT로 텍스트 생성
POST /api/fl8/forward      — 단일 forward pass
GET  /api/fl8/capabilities — 지원 기능 목록
```

#### **checkpoint_server.go** (모델 체크포인트)
```
POST /api/checkpoint/save   — 체크포인트 저장
GET  /api/checkpoint/list   — 체크포인트 목록
GET  /api/checkpoint/:id    — 특정 체크포인트 조회
POST /api/checkpoint/load   — 체크포인트 로드
POST /api/checkpoint/best   — 최고 모델 로드
```

#### **training_server.go** (학습 감시)
```
GET  /api/train/status      — 학습 상태
POST /api/train/start       — 학습 시작
POST /api/train/stop        — 학습 중지
GET  /api/train/metrics     — 손실/정확도 메트릭
```

#### **generation_server.go** (텍스트 생성)
```
POST /api/gen/batch         — 배치 생성
POST /api/gen/stream        — 스트림 생성
GET  /api/gen/templates     — 프롬프트 템플릿
```

#### **phase_i_server.go** (Phase I 테스트)
```
테스트 관련 엔드포인트
벤치마크 실행
검증 도구
```

#### **evaluation_server.go** (모델 평가)
```
POST /api/eval/perplexity   — 난해도 계산
POST /api/eval/bleu         — BLEU 점수
POST /api/eval/custom       — 커스텀 평가
```

#### **monitoring.go** (시스템 모니터링)
```
GET  /api/metrics/system    — CPU/메모리
GET  /api/metrics/request   — 요청 통계
GET  /api/metrics/model     — 모델 성능
GET  /api/health/full       — 종합 헬스체크
```

#### **dashboard.go** (웹 대시보드)
```
GET  /dashboard             — 학습 대시보드 (HTML)
GET  /dashboard/metrics.js  — 메트릭 데이터 (JSON)
GET  /dashboard/charts.js   — 차트 라이브러리
```

**빌드 결과** (8MB 바이너리):
- `freelang-gpt-api` (기본)
- `freelang-gpt-api-v2` (고급)
- `freelang-gpt-api-v3` (최신)
- `freelang-gpt-api-eval` (평가용)

---

## 🏗️ 아키텍처 비교

### FreeLang 서버 (Phase B 설계)
```
Request
   ↓
Router (main.fl)
   ↓
Handler Selection
   ↓
Handler Execution
   ↓
Store Update (Immutable)
   ↓
Response
```

**특징**:
- 함수형 (Immutable 저장소)
- 타입 안전
- 런타임 미결정 (HTTP 서버 대기)

### Go 서버 (완전 구현)
```
HTTP Request
   ↓
Mux Router
   ↓
Handler (goroutine)
   ↓
Business Logic
   ↓
Database/Model
   ↓
JSON Response
```

**특징**:
- 명령형 (Stateful)
- 완전한 HTTP 지원
- 프로덕션 준비 완료

### Node.js 서버 (특수 목적)
```
Graph JSON
   ↓
Graph Parser
   ↓
Compiler (→ .free)
   ↓
Executor
   ↓
HTML/CSS/JS + Preview
```

**특징**:
- 도메인 특화 (그래프 처리)
- FreeLang 코드 생성
- CSP 보안 (샌드박스)

---

## 📈 코드 규모

| 프로젝트 | 파일 수 | 라인 수 | 상태 |
|---------|--------|--------|------|
| freelang-rest-api | 5 | 500+ | 설계중 |
| freelang-backend-system | 20 | 2,000+ | 개발중 |
| FreeWire API | 8 | 1,500+ | 운영 |
| FreeLang GPT | 10 | 4,000+ | 운영 |

**총**: 43개 파일, 8,000+ 라인

---

## 🔒 보안 분석

### FreeWire API
- ✅ **CSP 헤더**: `script-src 'nonce-{nonce}'`
- ✅ **HTML 살균**: 스크립트 태그 제거, on* 속성 제거
- ✅ **요청 크기 제한**: 1MB
- ✅ **타임아웃**: 5초
- ⚠️ **CORS**: `*` (너무 허용적)

### Go 서버
- ✅ **CORS 헤더** 설정
- ✅ **타입 검증**: JSON 파싱 에러 처리
- ✅ **기본값 설정**: 임의의 값 방지
- ⚠️ **입력 검증**: 최소화
- ⚠️ **속도 제한**: 없음

### FreeLang 설계
- ✅ **타입 시스템**: 컴파일 타입 안전성
- ✅ **에러 처리**: Result<T, E> 필수
- ⏳ **런타임 검증**: 미정

---

## 📋 TODO & 다음 단계

### 즉시 필요
- [ ] **HTTP 서버 런타임** (Phase B): FreeLang 서버 실행 가능하게
- [ ] **데이터베이스 바인딩**: freelang-backend-system
- [ ] **보안 검토**: 입력 검증, 속도 제한, 인증

### 중기 (2-4주)
- [ ] **FreeLang REST API** 런타임 테스트
- [ ] **Mail Server** 기본 SMTP 구현
- [ ] **분산 시스템** Raft 통합 테스트

### 장기 (1-3개월)
- [ ] **프로덕션 배포**: Docker 패키징
- [ ] **성능 최적화**: 벤치마크
- [ ] **생태계 문서**: API 가이드

---

## 🎯 결론

**FreeLang 서버 생태계는 다음 3단계로 진화중**:

1. **🟢 운영중** (Node.js + Go)
   - FreeWire: 그래프 컴파일러 API ✅
   - FreeLang GPT: LLM REST API ✅

2. **🟡 설계/개발중** (FreeLang)
   - REST API Framework (Phase B 대기)
   - Backend System (다층 아키텍처)
   - Mail Server (기본 설계)
   - Distributed System (Raft 구현중)

3. **미래**
   - 모든 서버를 FreeLang으로 통일
   - 자체 HTTP 런타임 구축
   - 프로덕션급 신뢰성

**현재**: Go/Node.js로 구현된 완전한 API 3개 + FreeLang 설계 4개

**목표**: FreeLang 생태계 완성 (Phase B-E)
