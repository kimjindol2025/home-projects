---
name: Genspark Clone - 웹 검색 + AI 합산 엔진
description: 실시간 웹 검색 + Claude AI 분석 + Sparkpage 자동 생성 (913줄)
type: project
---

# Genspark Clone 프로젝트

## 프로젝트 개요

**이름**: Genspark Clone (웹 검색 + AI 합산 엔진)
**상태**: ✅ Alpha v1.0.0 완료
**완료일**: 2026-03-18
**코드**: 913줄 (Python)
**위치**: `/data/data/com.termux/files/home/projects/genspark-clone/`

## 핵심 기능

5단계 파이프라인으로 사용자 질문을 자동으로 분석하고 Sparkpage 생성

```
질문 → 분해 → 검색 → 추출 → AI 합산 → Sparkpage (HTML/MD)
```

## 구현 완료 사항

### 1️⃣ Query Analyzer (145줄)
- 사용자 질문 분석 (Claude haiku)
- 서브쿼리 분해 (2~5개)
- 예상 섹션 추천
- requests 기반 API 호출
- 파싱 실패 시 fallback

### 2️⃣ Web Searcher (182줄)
- DuckDuckGo HTML 파싱
- 단일/다중 검색 지원
- 1초 딜레이로 봇 차단 방지
- API 키 불필요

### 3️⃣ Content Fetcher (218줄)
- URL별 콘텐츠 페칭
- ThreadPoolExecutor 병렬화 (MAX_WORKERS=3)
- BeautifulSoup로 본문 추출
- 타임아웃 처리 (8초)
- 상세 에러 분류 (ok/timeout/error/blocked)

### 4️⃣ Claude Synthesizer (228줄)
- 멀티소스 콘텐츠 분석 (Claude sonnet)
- 섹션 구조화 (3~5개)
- 핵심 사실 추출
- 신뢰도 점수 계산
- 컨텍스트 크기 제한 (60K자)

### 5️⃣ Sparkpage Generator (272줄)
- 마크다운 생성
- HTML 생성 (외부 라이브러리 없음)
- Markdown → HTML 변환
- 반응형 CSS 포함
- 타임스탬프 기반 파일명

### 6️⃣ GensparkAgent (198줄)
- 5단계 파이프라인 오케스트레이션
- 설정 기반 실행
- 상세 로깅

### 테스트 & 문서
- `test_basic.py`: 기본 검증 (API 키 불필요)
- `test_integration.py`: 통합 테스트 (API 키 필수)
- `README.md`: 사용 설명서
- `ARCHITECTURE.md`: 상세 설계
- `QUICK_START.md`: 빠른 시작
- `COMPLETION_SUMMARY.md`: 구현 현황

## 주요 특징

### ✅ Termux 최적화
- 메모리 제약 (1GB) 대응
- MAX_WORKERS=3 병렬 제한
- 콘텐츠 크기 제한 (3K자)
- 컨텍스트 제한 (60K자)

### ✅ 의존성 최소화
- requests + beautifulsoup4만 필수
- stdlib만으로 완성 가능

### ✅ 에러 처리
- 계층별 fallback
- 타임아웃 안전성
- 봇 차단 회피

### ✅ 성능
- 전체 소요 시간: ~48초 (목표: 60초)
- API 호출: 2회 (haiku + sonnet)
- 메모리: < 100MB

## 테스트 결과

```
✅ QueryAnalyzer fallback OK
✅ DuckDuckGo 검색 (네트워크 상태)
✅ ContentFetcher OK (Python.org 790단어 추출)
✅ SparkpageGenerator OK
✅ 모든 기본 기능 동작
```

## 사용 방법

### CLI
```bash
python main.py "파이썬이란"
```

### 프로그래밍
```python
from src.genspark_agent import GensparkAgent, AgentConfig
import os

config = AgentConfig(
    anthropic_api_key=os.environ["ANTHROPIC_API_KEY"]
)
agent = GensparkAgent(config)
result = agent.run("REST API란")

print(f"HTML: {result.html_path}")
```

## 파일 구조

```
genspark-clone/
├── src/
│   ├── query_analyzer.py           # 질문 분석
│   ├── web_searcher.py             # 웹 검색
│   ├── content_fetcher.py          # 콘텐츠 추출
│   ├── claude_synthesizer.py       # AI 합산
│   ├── sparkpage_generator.py      # HTML/MD 생성
│   └── genspark_agent.py           # 오케스트레이션
├── main.py                         # CLI
├── test_basic.py                   # 기본 테스트
├── test_integration.py             # 통합 테스트
├── requirements.txt                # 의존성
└── 문서 (README, ARCHITECTURE 등)
```

## 기술 스택

- **Language**: Python 3.8+
- **Libraries**: requests, beautifulsoup4
- **API**: Claude API (haiku + sonnet)
- **Architecture**: 모듈식, 병렬 처리

## 성능 지표

| 지표 | 값 |
|------|-----|
| 총 코드 | 913줄 |
| 테스트 | 178줄 |
| 문서 | ~2,000줄 |
| 메모리 | < 100MB |
| 속도 | ~48초 |
| API 호출 | 2회 |

## 다음 단계

### 단기 (v1.1)
- 검색 결과 캐싱
- 이미지 추출
- 한글 폰트 최적화

### 중기 (v2.0)
- REST API 서버
- 데이터베이스
- 웹 UI

### 장기 (v3.0)
- 다중 검색 엔진
- 다국어 지원
- AI 모델 파인튜닝

## 특징 요약

🔍 **검색**: DuckDuckGo (API 키 불필요)
🤖 **분석**: Claude haiku + sonnet (최적화된 비용)
⚡ **병렬**: ThreadPoolExecutor (MAX_WORKERS=3)
📄 **생성**: HTML + Markdown (반응형 레이아웃)
🎯 **성능**: 48초 내 완성 (Termux 최적화)

## 관련 파일

- 소스: `src/` (6개 모듈)
- 테스트: `test_*.py`
- 문서: `README.md`, `ARCHITECTURE.md`, `QUICK_START.md`

## 활용 시나리오

1. 기술 질문 자동 답변
2. 뉴스 요약 생성
3. 학습 자료 자동 수집
4. 리서치 초안 작성

---

## 최종 배포 상태

**완료일**: 2026-03-18
**버전**: v1.0.0 Alpha
**GOGS 커밋**: 628b63d
**상태**: ✅ 완전 배포

### README 개선 사항
- 웹 UI 사용법 추가 (포트 커스터마이징, API 엔드포인트)
- 실제 Genspark와의 차이점 명확히 (표 형식)
- 프로덕션 배포 가이드 (Gunicorn, Systemd)
- 관련 문서 링크 통합

### 구현 완료 현황
- **Core**: 913줄 Python (6개 모듈)
- **Web UI**: 350줄 (Flask)
- **Docs**: 2,000줄 이상
- **Tests**: 2개 파일 (기본 + 통합)
- **GOGS**: 배포 완료 (628b63d)

### 핵심 특징
✅ 5단계 파이프라인 (분석→검색→추출→합산→생성)
✅ CLI + 웹 UI 지원
✅ Termux 메모리 최적화
✅ Claude haiku + sonnet 이중 모델
✅ 완전 문서화

### 미구현 기능 (v2.0+)
❌ Multi-Agent 교차 검증
❌ Consensus Engine (정보 충돌 해결)
❌ 동적 Sparkpage (위젯 기반)
❌ 캐싱 시스템

---

## 기술 스택 최종 정리

| 계층 | 도구 |
|------|------|
| API | Claude haiku + sonnet (requests 직접 사용) |
| 검색 | DuckDuckGo (HTML 파싱) |
| 파싱 | BeautifulSoup4 |
| 병렬화 | ThreadPoolExecutor (MAX_WORKERS=3) |
| 웹 | Flask + Vanilla JS |
| 배포 | Git + GOGS + Gunicorn (선택) |
