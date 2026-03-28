---
name: FreeLang 생태계 통합 계획
description: 3개 TOP 프로젝트 최종 통합 및 배포 전략
type: project
---

# 🌍 FreeLang 생태계 통합 & 배포

**현황**: 2026-03-25
**상태**: 🔄 **최종 통합 진행 중**

---

## 📊 프로젝트 현황

| 프로젝트 | 상태 | 코드량 | GOGS |
|---------|------|--------|------|
| **Bank System** | ✅ 100% | 5,383줄 | ✅ 푸시 |
| **Playground** | ✅ 100% | 2,094줄 | ✅ 푸시 (Phase 3) |
| **Website** | 🟡 80% | 407줄 | ✅ 푸시 (Phase 1) |
| **합계** | **90%** | **7,884줄** | - |

---

## 🎯 통합 목표

### 1단계: 웹 기반 접근
모든 서비스에 **localhost 통일 포트**:
```
localhost:3000  → Bank System Dashboard
localhost:3001  → Playground IDE
localhost:4000  → Website (Docusaurus)
```

### 2단계: Nginx 리버스 프록시
단일 진입점:
```
/api/           → Bank System API (8080)
/playground     → Playground IDE (3000)
/docs/          → Website (4000)
/               → Landing Page
```

### 3단계: 통합 docker-compose.yml
모든 서비스를 한 명령어로 배포:
```bash
docker-compose up -d
```

---

## 🏗️ 통합 아키텍처

```
┌─────────────────────────────────────┐
│   Nginx 리버스 프록시 (포트 80)      │
│  - hostname: localhost (또는 도메인)  │
├─────────────────────────────────────┤
│ GET /          → Landing Page       │
│ GET /docs/*    → Website            │
│ GET /playground → Playground        │
│ GET /api/*     → Bank API           │
│ POST /api/*    → Bank API           │
└─────────────────────────────────────┘
         ↓         ↓           ↓
    Website    Playground   Bank System
    (4000)     (3000)       (8080)
```

---

## 📝 필요한 작업

### Website (Phase 2)
- **현황**: Docusaurus 기본 설정 + 9개 마크다운 문서
- **필요**:
  1. `package.json` 업데이트 (빌드 스크립트)
  2. `docusaurus.config.js` 수정 (URL 설정)
  3. Dockerfile 추가 (빌드 & 배포)
  4. 예상: +80줄

### 통합 docker-compose.yml
- 3개 서비스 조정:
  - freelang-bank-system
  - freelang-playground
  - freelang-website
- Nginx 리버스 프록시 추가
- 네트워크 통합 (ecosystem-network)
- 볼륨 및 의존성 정의
- 예상: ~200줄

### Landing Page
- 간단한 HTML 페이지
- 3개 서비스 링크
- 프로젝트 소개
- 예상: ~100줄

---

## 🔧 배포 체크리스트

- [ ] Website Phase 2 완료
  - [ ] Dockerfile 추가
  - [ ] npm 설치
  - [ ] GOGS 푸시

- [ ] 통합 docker-compose.yml 생성
  - [ ] 3개 서비스 정의
  - [ ] Nginx 리버스 프록시
  - [ ] 네트워크 격리
  - [ ] 헬스체크 체인

- [ ] nginx.conf (리버스 프록시)
  - [ ] /api 라우팅
  - [ ] /playground 라우팅
  - [ ] /docs 라우팅
  - [ ] CORS 설정

- [ ] 최종 검증
  - [ ] docker-compose up -d 성공
  - [ ] curl localhost/api/health
  - [ ] curl localhost/playground
  - [ ] curl localhost/docs

---

## 🚀 완성 예상 일정

| 단계 | 작업 | 예상 시간 |
|------|------|----------|
| 1 | Website Phase 2 | 30분 |
| 2 | 통합 docker-compose | 20분 |
| 3 | Nginx 설정 | 15분 |
| 4 | 최종 테스트 & GOGS 푸시 | 15분 |
| **합계** | | **80분** |

---

## 📈 완성도 전망

```
현재: 90% (7,884줄)
↓
최종: 100% (8,150줄+)

- Bank System: 100% ✅
- Playground: 100% ✅
- Website: 100% → Phase 2
- Ecosystem: 100% → 통합 완료
```

