#!/bin/bash

# 📦 GOGS 레포 자동 생성 스크립트
# 사용: bash gogs-repo-auto-setup.sh <project-name> [category] [user]

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 변수
PROJECT_NAME=${1:-}
CATEGORY=${2:-core}
GOGS_USER=${3:-freelang}
GOGS_HOST=${GOGS_HOST:-gogs.local}
GOGS_PROTOCOL=${GOGS_PROTOCOL:-http}

# 입력 검증
if [ -z "$PROJECT_NAME" ]; then
  echo -e "${RED}❌ 사용법: bash gogs-repo-auto-setup.sh <project-name> [category] [user]${NC}"
  echo "예시: bash gogs-repo-auto-setup.sh freelang-v4 core freelang"
  exit 1
fi

# 전역 경로
PROJECTS_DIR="$HOME/.projects"
PROJECT_PATH="$PROJECTS_DIR/$CATEGORY/$PROJECT_NAME"

echo -e "${BLUE}🚀 GOGS 레포 자동 생성 시작${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "프로젝트: $PROJECT_NAME"
echo "카테고리: $CATEGORY"
echo "GOGS 사용자: $GOGS_USER"
echo "GOGS 주소: $GOGS_PROTOCOL://$GOGS_HOST"
echo "로컬 경로: $PROJECT_PATH"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo

# 1단계: 로컬 디렉토리 생성
echo -e "${BLUE}1️⃣ 로컬 디렉토리 생성...${NC}"
mkdir -p "$PROJECT_PATH"/{src,tests,docs,.claude/projects/$PROJECT_NAME/memory}
cd "$PROJECT_PATH"
echo -e "${GREEN}✓ 완료${NC}"

# 2단계: 필수 파일 생성
echo -e "${BLUE}2️⃣ 필수 파일 생성...${NC}"

# .gitignore
cat > .gitignore << 'GITIGNORE'
node_modules/
dist/
build/
target/
.env
.env.local
.env.*.local
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store
Thumbs.db
*.log
npm-debug.log*
.cache/
.eslintcache
GITIGNORE

# README.md
cat > README.md << EOF
# $PROJECT_NAME

**GOGS 레포**: $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME
**카테고리**: $CATEGORY
**상태**: 🆕 신규

## 프로젝트 개요
프로젝트 설명을 작성하세요.

## 설치
\`\`\`bash
git clone $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME.git
cd $PROJECT_NAME
npm install
\`\`\`

## 폴더 구조
- \`src/\`: 소스 코드
- \`tests/\`: 테스트
- \`docs/\`: 문서

## 라이선스
MIT
EOF

# CLAUDE.md
cat > CLAUDE.md << EOF
# $PROJECT_NAME - Claude AI 작업 가이드

## GOGS 정보
- **레포**: $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME
- **브랜치**: master (기본)
- **클론**: git clone $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME.git

## 작업 규칙

### 커밋 메시지 형식
\`\`\`
feat:    새 기능
fix:     버그 수정
docs:    문서
refactor: 코드 정리
test:    테스트

[마지막 줄]
Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
\`\`\`

### 폴더 구조
- src/: 수정 권장
- tests/: 확장 권장
- docs/: 읽기 권장

## 메모리 파일
- .claude/projects/$PROJECT_NAME/memory/MEMORY.md (필수 유지)
EOF

# package.json
cat > package.json << EOF
{
  "name": "$PROJECT_NAME",
  "version": "1.0.0",
  "description": "FreeLang Project",
  "main": "src/index.ts",
  "scripts": {
    "build": "tsc",
    "test": "jest",
    "dev": "ts-node src/index.ts"
  },
  "keywords": ["freelang"],
  "license": "MIT"
}
EOF

# MEMORY.md
cat > .claude/projects/$PROJECT_NAME/memory/MEMORY.md << EOF
# $PROJECT_NAME - 프로젝트 메모리

**생성**: $(date +%Y-%m-%d)
**상태**: 🆕 신규
**GOGS**: $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME

## 프로젝트 개요
프로젝트 설명 입력

## 진행 상황
- [ ] 단계 1: 기초 설정
- [ ] 단계 2: 핵심 기능 개발
- [ ] 단계 3: 테스트 & 최적화

## 다음 액션
1. README.md 작성
2. 구조 확정
3. 첫 테스트 작성
EOF

echo -e "${GREEN}✓ 완료${NC}"

# 3단계: Git 초기화
echo -e "${BLUE}3️⃣ Git 초기화...${NC}"
git init
git remote add origin "$GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME.git"
echo -e "${GREEN}✓ 완료${NC}"

# 4단계: 초기 커밋
echo -e "${BLUE}4️⃣ 초기 커밋...${NC}"
git add -A
git commit -m "Initial commit: Project setup

- 프로젝트 구조 초기화
- 표준 파일 추가
- CLAUDE.md 작업 규칙 정의

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"
echo -e "${GREEN}✓ 완료${NC}"

echo
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ GOGS 레포 설정 완료!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo
echo "📚 다음 단계:"
echo "1. GOGS에서 레포 생성 (웹 인터페이스)"
echo "2. 다음 명령 실행:"
echo "   git push -u origin master"
echo
echo "🔗 GOGS 레포: $GOGS_PROTOCOL://$GOGS_HOST/$GOGS_USER/$PROJECT_NAME"
echo "📁 로컬 경로: $PROJECT_PATH"
echo

# 권한 설정
chmod +x "$0"

