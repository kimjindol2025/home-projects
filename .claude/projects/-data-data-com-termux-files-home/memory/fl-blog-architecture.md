---
name: FreeLang 블로그 아키텍처 (최종)
description: GitHub 프론트 + 253 서버 백엔드 연결 구조
type: project
---

# FreeLang 블로그 최종 아키텍처

**업데이트**: 2026-03-29
**상태**: ✅ 운영 중 (GitHub + 253 서버)

---

## 🎯 지시사항 (이것이 규칙)

1. **GitHub = 프론트엔드**
   - `index.html` (517줄)
   - BACKEND_URL = 'http://192.0.0.2:8253'
   - 저장소: github.com/kimjindol2025/projects/fl-blog

2. **253 서버 = 백엔드 (메인)**
   - `server-comments.js` (526줄)
   - 포트: 8253
   - 포스트: 89개
   - 댓글: JSON 저장소 (comments.json)
   - 모든 API 엔드포인트 제공

3. **블로그 유지**
   - 지속적인 개발/유지보수
   - 포스트 추가/수정/삭제
   - 댓글 기능 관리
   - API 개선/최적화

---

## 🔗 연결 구조

```
GitHub Pages (프론트엔드)
        ↓ Fetch API
        ↓
253 서버 (백엔드)
        ↓
comments.json (댓글) + posts/ (89개)
```

---

## 📊 현재 상태

**GitHub 저장소**
- ✅ index.html (BACKEND_URL 설정됨)
- ✅ server-comments.js (253 서버 코드)
- ✅ posts/ (89개 마크다운)
- ✅ comments.json (댓글 저장소)
- ✅ package.json (Node.js 설정)

**253 서버 (실행 중)**
- ✅ 포트: 8253
- ✅ 상태: 정상 (PID: 21078)
- ✅ 모든 API: 작동 중 (HTTP 200)
- ✅ 포스트: 89개 로드됨
- ✅ 댓글: 저장소 정상

---

## 🔄 개발 워크플로우

1. **코드 수정**
   - server-comments.js (백엔드)
   - index.html (프론트엔드)
   - posts/ (포스트)

2. **253 서버 테스트**
   - http://localhost:8253 확인
   - API 테스트

3. **GitHub 커밋**
   - git add .
   - git commit -m "..."
   - git push origin master

4. **서버 재시작** (필요시)
   - pkill -f "node.*server"
   - node server-comments.js 8253

---

## 📍 접속 정보

- **로컬**: http://localhost:8253
- **네트워크**: http://192.0.0.2:8253
- **API**: /api/posts, /api/comments, /api/search 등

---

## ⚠️ 주의사항

- 253 서버는 항상 실행 상태로 유지
- GitHub는 프론트엔드만 담당 (index.html)
- 백엔드 코드 변경 후 서버 재시작 필수
- comments.json은 자동 생성/관리됨

---

**Why**: GitHub와 로컬 서버 분리로 안정적이고 빠른 개발 가능

**How to apply**: 모든 기능 개발/유지는 이 구조를 기반으로 진행
