---
name: kim-project-cli API 엔드포인트 참조
description: 라이브 배포된 프로젝트 관리 API의 모든 엔드포인트 정보
type: reference
---

# 📍 kim-project-cli API 엔드포인트 참조

## 기본 정보

### 프로젝트 API (기존)
**도메인**: https://project-cli.dclub.kr
**배포 위치**: 73 서버
**포트**: 기본 (HTTPS)
**상태**: ✅ 라이브 운영 (2026-03-16 확인)

### Claude 메모 API (신규)
**도메인**: localhost 또는 내부 네트워크
**포트**: 40012
**상태**: ⏳ 준비됨 (배포 대기)

---

## 엔드포인트 목록

### 1️⃣ 헬스 체크
```
GET /health
```

**응답** (200 OK):
```json
{
  "status": "ok",
  "service": "gogs-lookup-api"
}
```

**용도**: API 서버 가용성 확인

**응답 시간**: ~50ms

---

### 2️⃣ 전체 프로젝트 조회
```
GET /api/gogs/projects
```

**응답** (200 OK):
```json
{
  "total": 27,
  "projects": [
    {
      "repo_name": "FREELANG_Projects",
      "repo_path": "/kim/FREELANG_Projects",
      "repo_url": "https://gogs.dclub.kr/kim/FREELANG_Projects.git",
      "folder_path": "/home/kimjin/kim/Desktop/kim/01_Active_Projects/FREELANG_Projects",
      "status": "활성중",
      "last_work": "2026-03-15",
      "lookup_url": "/api/gogs/lookup?repo=FREELANG_Projects"
    },
    // ... 26개 더
  ]
}
```

**파라미터**: 없음

**용도**: 모든 추적 프로젝트 목록 조회

**응답 시간**: ~100ms

---

### 3️⃣ 개별 저장소 상세 조회
```
GET /api/gogs/lookup?repo=<repo-name>
```

**파라미터**:
- `repo`: 저장소 이름 (필수)
  - 예: `gogs-architect`, `dns-manager`, `FREELANG_Projects`

**응답** (200 OK):
```json
{
  "success": true,
  "gogs_repo": {
    "name": "gogs-architect",
    "url": "https://gogs.dclub.kr/kim/gogs-architect.git",
    "path": "/kim/gogs-architect"
  },
  "project": {
    "id": 7,
    "name": "gogs-architect",
    "status": "활성중",
    "category": ["도구"],
    "priority": 3,
    "progress": 0,
    "description": "gogs-architect - 2026-03-15 마이그레이션",
    "tech_stack": ["Other"],
    "last_work_date": "2026-03-15",
    "status_detail": "마이그레이션됨",
    "created_at": "2026-03-15T13:53:41.646Z",
    "updated_at": "2026-03-15T13:53:41.646Z"
  },
  "folder": {
    "path": "/home/kimjin/kim/Desktop/kim/01_Active_Projects/gogs-architect",
    "exists": false,
    "size": "계산 불가"
  },
  "files": {
    "changes_md": {
      "exists": false,
      "path": "/home/kimjin/kim/Desktop/kim/01_Active_Projects/gogs-architect/.claude/changes.md",
      "content": null,
      "size": 0
    },
    "metadata_json": {
      "exists": false,
      "path": "...",
      "content": null
    },
    "readme_md": {
      "exists": false,
      "path": "...",
      "preview": null
    }
  }
}
```

**용도**: 특정 저장소의 모든 메타데이터 조회

**응답 시간**: ~80ms

---

### 4️⃣ 폴더 경로 조회
```
GET /api/gogs/resolve-path?repo=<repo-name>
```

**파라미터**:
- `repo`: 저장소 이름 (필수)

**응답** (200 OK):
```json
{
  "repo": "dns-manager",
  "folder_path": "/home/kimjin/kim/Desktop/kim/01_Active_Projects/dns-manager",
  "exists": false
}
```

**용도**: 스크립트에서 폴더 경로만 빠르게 추출

**응답 시간**: ~60ms

---

## 추적 프로젝트 목록 (27개)

모든 저장소는 `/kim/` 경로로 시작하며, 로컬 폴더는 `/home/kimjin/kim/Desktop/kim/01_Active_Projects/` 아래에 위치

1. FREELANG_Projects
2. aion-backup
3. claude-push-api
4. dns-manager
5. fv-lang
6. gemini-bridge
7. gogs-architect
8. gogs-github-webhook
9. gogs_project
10. grie-engine
11. home-cleanup-tool
12. kim
13. kim-ai-os
14. kimdb
15. kimsearch
16. linux-file-manager
17. moss-state-core
18. projects
19. pyfree
20. robot-ai-project
21. secure-hasher
22. src
23. ssh-hub
24. ssh-sftp-server
25. synaptic-pilot
26. team-agents-guide
27. xpress

---

## 사용 예시

### Bash/Shell
```bash
# 1️⃣ 헬스 체크
curl https://project-cli.dclub.kr/health

# 2️⃣ 모든 프로젝트 조회
curl https://project-cli.dclub.kr/api/gogs/projects

# 3️⃣ 특정 저장소 조회
curl "https://project-cli.dclub.kr/api/gogs/lookup?repo=gogs-architect"

# 4️⃣ 경로만 추출해서 cd
FOLDER=$(curl -s "https://project-cli.dclub.kr/api/gogs/resolve-path?repo=gogs-architect" | grep -o '"folder_path":"[^"]*' | cut -d'"' -f4)
cd $FOLDER
```

### JavaScript/Node.js
```javascript
// 전체 프로젝트 조회
fetch('https://project-cli.dclub.kr/api/gogs/projects')
  .then(r => r.json())
  .then(data => console.log(`${data.total}개 프로젝트`));

// 특정 저장소 조회
fetch('https://project-cli.dclub.kr/api/gogs/lookup?repo=gogs-architect')
  .then(r => r.json())
  .then(data => console.log(data.folder.path));
```

### Python
```python
import requests

# 경로 조회
resp = requests.get('https://project-cli.dclub.kr/api/gogs/resolve-path?repo=dns-manager')
data = resp.json()
print(data['folder_path'])
```

---

## 에러 처리

### 저장소를 찾을 수 없는 경우
```json
{
  "error": "프로젝트 없음",
  "repo": "존재하지-않는-이름",
  "hint": "DB에 등록되지 않은 저장소입니다"
}
```

**상태 코드**: 404

---

## 성능 메트릭

| 엔드포인트 | 응답 시간 | 특징 |
|-----------|---------|------|
| `/health` | ~50ms | 가장 빠름 |
| `/api/gogs/resolve-path` | ~60ms | 경로만 필요할 때 |
| `/api/gogs/lookup` | ~80ms | 모든 메타데이터 |
| `/api/gogs/projects` | ~100ms | 대량 데이터 |

---

## 관련 문서

- **GOGS 저장소**: https://gogs.dclub.kr/kim/kim-project-cli.git
- **로컬 클론**: ~/kim-project-cli/
- **로컬 데이터베이스**: ~/kim/kim-projects.db
- **상세 가이드**: [kim-project-cli.md](./kim-project-cli.md)

