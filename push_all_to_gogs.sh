#!/bin/bash

# GOGS 푸시 자동화 스크립트
# 모든 프로젝트 폴더를 GOGS에 푸시

set -e

GOGS_URL="https://gogs.dclub.kr"
GOGS_USER="kim"
GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
HOME_DIR="/data/data/com.termux/files/home"

# 색상
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 통계
TOTAL=0
SUCCESS=0
FAILED=0
SKIPPED=0

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}GOGS 전체 프로젝트 푸시 시작${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 프로젝트 폴더 순회
cd "$HOME_DIR"

for folder in */; do
    folder_name="${folder%/}"

    # 제외 폴더 (숨김, 설정 등)
    if [[ "$folder_name" =~ ^(\.|[a-z0-9]*cache|tmp|tests|storage|node_modules)$ ]]; then
        continue
    fi

    TOTAL=$((TOTAL + 1))

    echo -e "${YELLOW}[${TOTAL}] 처리 중: $folder_name${NC}"

    cd "$HOME_DIR/$folder_name" 2>/dev/null || {
        echo -e "${RED}   ✗ 폴더 접근 실패${NC}"
        FAILED=$((FAILED + 1))
        continue
    }

    # Git 저장소 초기화 확인
    if [ ! -d ".git" ]; then
        echo -e "${YELLOW}   → Git 저장소 초기화 중...${NC}"
        git init > /dev/null 2>&1
        git config user.name "Kim Developer" > /dev/null 2>&1
        git config user.email "kim@freelang.dev" > /dev/null 2>&1
        git add . > /dev/null 2>&1
        git commit -m "Initial commit: $folder_name" > /dev/null 2>&1 || true
    fi

    # 원격 저장소 추가/업데이트
    REMOTE_URL="https://${GOGS_USER}:${GOGS_TOKEN}@gogs.dclub.kr/${GOGS_USER}/${folder_name}.git"

    if git remote | grep -q origin; then
        git remote set-url origin "$REMOTE_URL" > /dev/null 2>&1
    else
        git remote add origin "$REMOTE_URL" > /dev/null 2>&1
    fi

    # 현재 브랜치 확인
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "master")

    # 푸시 시도
    if git push -f -u origin "$CURRENT_BRANCH" > /dev/null 2>&1; then
        echo -e "${GREEN}   ✓ 성공 (${CURRENT_BRANCH})${NC}"
        SUCCESS=$((SUCCESS + 1))
    else
        echo -e "${RED}   ✗ 실패${NC}"
        FAILED=$((FAILED + 1))
    fi

done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}푸시 완료!${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "총 폴더: ${YELLOW}${TOTAL}${NC}"
echo -e "성공: ${GREEN}${SUCCESS}${NC}"
echo -e "실패: ${RED}${FAILED}${NC}"
echo ""
echo -e "${GREEN}GOGS 서버: https://gogs.dclub.kr/kim${NC}"
echo ""
