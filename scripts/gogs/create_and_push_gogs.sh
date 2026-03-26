#!/bin/bash

# GOGS 저장소 생성 및 푸시 자동화 스크립트
# 1. GOGS API로 저장소 생성
# 2. 각 폴더를 git 저장소로 초기화
# 3. GOGS에 푸시

set -e

GOGS_URL="https://gogs.dclub.kr"
GOGS_USER="kim"
GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
HOME_DIR="/data/data/com.termux/files/home"
API_URL="${GOGS_URL}/api/v1"

# 색상
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 통계
TOTAL=0
CREATED=0
PUSHED=0
FAILED=0

echo -e "${BLUE}======================================== ${NC}"
echo -e "${CYAN}GOGS 저장소 생성 및 푸시 시작${NC}"
echo -e "${BLUE}======================================== ${NC}"
echo ""

# 제외 폴더 목록
is_excluded() {
    local folder="$1"
    if [[ "$folder" =~ ^(\.|.*cache.*|tmp|tests|storage|node_modules|bbbb|cccc|go|projects|src|kim|gogs_.*|python-.*)$ ]]; then
        return 0
    fi
    return 1
}

# 프로젝트 폴더 순회
cd "$HOME_DIR"

for folder in */; do
    folder_name="${folder%/}"

    # 제외 폴더 스킵
    if is_excluded "$folder_name"; then
        continue
    fi

    TOTAL=$((TOTAL + 1))

    echo -e "${YELLOW}[${TOTAL}] ${folder_name}${NC}"

    # 1단계: GOGS에 저장소 생성
    echo -ne "   → 저장소 생성 중... "

    CREATE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: token ${GOGS_TOKEN}" \
        -H "Content-Type: application/json" \
        -d "{\"name\": \"${folder_name}\", \"description\": \"${folder_name} repository\"}" \
        "${API_URL}/user/repos")

    # 저장소 생성 확인
    if echo "$CREATE_RESPONSE" | grep -q '"id"'; then
        echo -e "${GREEN}✓${NC}"
        CREATED=$((CREATED + 1))
    else
        # 이미 존재하는 경우
        if echo "$CREATE_RESPONSE" | grep -q '"message"'; then
            echo -e "${YELLOW}(이미 존재)${NC}"
        else
            echo -e "${RED}✗${NC}"
            echo "   응답: $CREATE_RESPONSE"
            FAILED=$((FAILED + 1))
            continue
        fi
    fi

    # 2단계: 폴더로 이동 및 git 초기화
    cd "$HOME_DIR/$folder_name" 2>/dev/null || {
        echo -e "${RED}   ✗ 폴더 접근 실패${NC}"
        FAILED=$((FAILED + 1))
        continue
    }

    echo -ne "   → 저장소 초기화 중... "

    # Git 저장소 초기화
    if [ ! -d ".git" ]; then
        git init > /dev/null 2>&1
        git config user.name "Kim Developer" > /dev/null 2>&1
        git config user.email "kim@freelang.dev" > /dev/null 2>&1
        git add . > /dev/null 2>&1
        git commit -m "Initial commit: $folder_name" > /dev/null 2>&1 || true
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${CYAN}(기존 저장소)${NC}"
    fi

    # 3단계: GOGS에 푸시
    echo -ne "   → 푸시 중... "

    REMOTE_URL="https://${GOGS_USER}:${GOGS_TOKEN}@gogs.dclub.kr/${GOGS_USER}/${folder_name}.git"

    if git remote | grep -q origin; then
        git remote set-url origin "$REMOTE_URL" > /dev/null 2>&1
    else
        git remote add origin "$REMOTE_URL" > /dev/null 2>&1
    fi

    # 현재 브랜치 확인
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "master")

    # 푸시
    if git push -f -u origin "$CURRENT_BRANCH" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        echo -e "   ${CYAN}→ https://gogs.dclub.kr/${GOGS_USER}/${folder_name}${NC}"
        PUSHED=$((PUSHED + 1))
    else
        echo -e "${RED}✗${NC}"
        FAILED=$((FAILED + 1))
    fi

done

echo ""
echo -e "${BLUE}======================================== ${NC}"
echo -e "${CYAN}작업 완료!${NC}"
echo -e "${BLUE}======================================== ${NC}"
echo -e "총 폴더: ${YELLOW}${TOTAL}${NC}"
echo -e "저장소 생성: ${GREEN}${CREATED}${NC}"
echo -e "푸시 성공: ${GREEN}${PUSHED}${NC}"
echo -e "실패: ${RED}${FAILED}${NC}"
echo ""
echo -e "${CYAN}GOGS 대시보드:${NC}"
echo -e "${GREEN}https://gogs.dclub.kr/${GOGS_USER}${NC}"
echo ""
