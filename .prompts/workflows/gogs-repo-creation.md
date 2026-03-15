---
name: GOGS Repository Creation Workflow
description: GOGS 신규 레포지토리 생성 표준 프로세스
type: workflow
version: 1.0
author: Claude Code
created: 2026-03-15
---

# 📦 GOGS 신규 레포 생성 규칙

## 🎯 목적
GOGS에 새로운 Git 레포지토리를 생성할 때 준수해야 할 표준 규칙입니다.

---

## 📋 레포 이름 규칙

### 명명 규칙
```
{project-type}-{project-name}[-v{version}]

예시:
✅ freelang-v4                  (코어 프로젝트)
✅ module-compiler              (모듈)
✅ tool-gogs-api                (도구)
✅ experiment-neural-network    (실험)
✅ archive-v2-backup            (아카이브)
```

### 카테고리 프리픽스
```
freelang-*        → FreeLang 핵심 프로젝트
module-*          → 언어 모듈
tool-*            → 도구/유틸
experiment-*      → 실험 프로젝트
archive-*         → 아카이브
integration-*     → 통합 프로젝트
```

---

## 📁 레포 생성 순서 (단계별)

### Step 1: GOGS 웹 인터페이스에서 생성

```
1. GOGS 로그인 (http://gogs.local 또는 설정된 주소)
2. "+" → "New Repository"
3. 정보 입력:
   - Repository Name: {project-type}-{project-name}
   - Description: 프로젝트 설명
   - Private: 필요시 체크
   - Initialize with: README
   - License: MIT
```

### Step 2: 로컬에 클론

```bash
# 디렉토리 이동
cd ~/.projects/[category]/[project-name]

# GOGS 레포 초기화
git clone http://gogs.local/[user]/freelang-v4.git .
# 또는 기존 프로젝트에 원격 추가
git remote add origin http://gogs.local/[user]/freelang-v4.git
```

### Step 3: 필수 파일 생성

```bash
# 이미 존재하는 것 확인
ls -la
# 필수 파일:
# ✓ README.md
# ✓ CLAUDE.md
# ✓ .gitignore
# ✓ package.json
```

### Step 4: 초기 커밋

```bash
git add -A
git commit -m "Initial commit: Project setup

- 프로젝트 구조 초기화
- 표준 파일 추가
- CLAUDE.md 작업 규칙 정의

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"

git push -u origin master
```

### Step 5: Claude 메모리 초기화

```bash
# 메모리 파일 생성
mkdir -p .claude/projects/[project-name]/memory

# MEMORY.md 초기화
cat > .claude/projects/[project-name]/memory/MEMORY.md << 'MEMORY'
# 프로젝트 메모리 (2026-03-15 생성)
- 상태: 🆕 신규
- GOGS 레포: http://gogs.local/[user]/[repo-name]
MEMORY
```

---

## 📝 필수 파일 템플릿

### .gitignore (표준)
```gitignore
# Dependencies
node_modules/
dist/
build/
target/

# Environment
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log
npm-debug.log*

# Cache
.cache/
.eslintcache
```

### README.md (템플릿)
```markdown
# {{project-name}}

**GOGS 레포**: http://gogs.local/[user]/{{project-name}}
**카테고리**: {{category}}
**상태**: {{status}}

## 프로젝트 개요
[프로젝트 설명]

## 설치
\`\`\`bash
git clone http://gogs.local/[user]/{{project-name}}.git
cd {{project-name}}
npm install
\`\`\`

## 사용법
[사용 방법]

## 구조
```
src/       - 소스 코드
tests/     - 테스트
docs/      - 문서
\`\`\`

## 기여 가이드
[기여 방법]

## 라이선스
MIT
\`\`\`

### CLAUDE.md (템플릿)
```markdown
# {{project-name}} - Claude AI 작업 가이드

## GOGS 정보
- **레포**: http://gogs.local/[user]/{{project-name}}
- **브랜치**: master (기본)
- **클론**: git clone http://gogs.local/[user]/{{project-name}}.git

## 작업 규칙

### 브랜치 전략
- master: 프로덕션 코드
- develop: 개발 브랜치
- feature/*: 기능 개발
- fix/*: 버그 수정

### 커밋 메시지 규칙
\`\`\`
feat:    새 기능
fix:     버그 수정
docs:    문서
refactor: 코드 정리
test:    테스트
\`\`\`

### 필수 서명
모든 커밋은 다음으로 끝남:
\`\`\`
Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
\`\`\`

## 폴더 구조
- src/: 수정 권장
- tests/: 확장 권장
- docs/: 읽기 권장

## 메모리 파일
- .claude/projects/{{project}}/memory/MEMORY.md (필수 유지)
```

---

## 🔄 Git 워크플로우

### 기본 작업 흐름

```bash
# 1. 최신 코드 가져오기
git pull origin master

# 2. 기능 브랜치 생성 (선택)
git checkout -b feature/new-feature

# 3. 코드 작성
# → src/ 폴더에 작성

# 4. 테스트
npm test

# 5. 스테이징
git add src/
git add tests/

# 6. 커밋
git commit -m "feat: 새 기능 추가

상세 설명...

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"

# 7. 푸시
git push origin [branch-name]

# 8. 메모리 업데이트
cat .claude/projects/{{project}}/memory/MEMORY.md
# → 완료된 작업 기록
```

### Merge 프로세스 (Master)

```bash
# 1. Master로 이동
git checkout master

# 2. 최신 업데이트
git pull origin master

# 3. Feature 브랜치 병합
git merge feature/new-feature

# 4. 푸시
git push origin master

# 5. 브랜치 정리
git branch -d feature/new-feature
git push origin --delete feature/new-feature
```

---

## ✅ 체크리스트 (레포 생성 시)

### GOGS 생성
- [ ] 레포 이름 규칙 준수
- [ ] 설명 입력
- [ ] 라이선스 선택 (기본: MIT)
- [ ] .gitignore 자동 추가
- [ ] README 자동 생성

### 로컬 설정
- [ ] 클론 또는 remote 추가
- [ ] 필수 파일 확인 (.gitignore, README, package.json)
- [ ] CLAUDE.md 작성
- [ ] .claude/projects/[name]/memory/ 생성

### 초기 커밋
- [ ] 모든 파일 스테이징
- [ ] 초기 커밋 메시지 작성
- [ ] Co-Authored-By 포함
- [ ] Master에 푸시

### Claude 메모리
- [ ] MEMORY.md 생성
- [ ] GOGS URL 기록
- [ ] 상태 설정 (🆕 신규)
- [ ] 다음 액션 입력

---

## 🔗 GOGS 및 로컬 동기화

### 푸시 (로컬 → GOGS)
```bash
git push origin master
```

### 풀 (GOGS → 로컬)
```bash
git pull origin master
```

### 상태 확인
```bash
git status
git log --oneline (최근 커밋 확인)
```

---

## 📊 프로젝트 구조 (예시)

```
~/.projects/core/freelang-v4/
├─ .git/                      (GOGS 연동)
├─ .claude/
│  └─ projects/freelang-v4/
│     └─ memory/
│        └─ MEMORY.md         (세션 기록)
│
├─ src/
│  ├─ main.ts
│  └─ utils/
│
├─ tests/
│  └─ test.ts
│
├─ docs/
│  └─ README.md
│
├─ .gitignore
├─ CLAUDE.md                 (이 파일)
├─ README.md
├─ package.json
└─ .env.example
```

---

## 🚀 자동화 스크립트

### 신규 레포 자동 생성

```bash
#!/bin/bash
# gogs-create-repo.sh

PROJECT_NAME=$1
CATEGORY=${2:-core}
GOGS_USER=${3:-default}
GOGS_HOST="gogs.local"

# 1. GOGS 레포 생성 (API)
curl -X POST \
  -H "Authorization: token $GOGS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$PROJECT_NAME\",\"description\":\"FreeLang Project\"}" \
  http://$GOGS_HOST/api/v1/user/repos

# 2. 로컬 디렉토리 생성
mkdir -p ~/.projects/$CATEGORY/$PROJECT_NAME/{src,tests,docs}
cd ~/.projects/$CATEGORY/$PROJECT_NAME

# 3. CLAUDE.md 생성
cat > CLAUDE.md << EOF
# $PROJECT_NAME - Claude AI 작업 가이드

## GOGS
- 레포: http://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME
EOF

# 4. Git 초기화
git init
git remote add origin http://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME.git

# 5. 초기 커밋
git add -A
git commit -m "Initial commit: Project setup

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"

git push -u origin master

echo "✅ 레포 생성 완료: $PROJECT_NAME"
```

---

## 📌 주의사항

❌ **하지 말 것**:
- 커밋 메시지에 Co-Authored-By 빼먹기
- README.md 없이 푸시
- .env 파일 커밋
- MEMORY.md 미업데이트

✅ **해야 할 것**:
- 정기적으로 pull/push
- MEMORY.md 유지
- 표준 커밋 형식 준수
- 테스트 후 커밋

---

## 🔗 참고 링크

- GOGS 공식: https://gogs.io/
- FreeLang 프로젝트: ~/.projects/
- Claude 메모리: ~/.claude/projects/*/memory/
- 프롬프트 저장소: ~/.prompts/

---

**Version**: 1.0
**Last Updated**: 2026-03-15
**Author**: Claude Code

