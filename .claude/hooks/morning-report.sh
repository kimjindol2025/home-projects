#!/bin/bash
# Morning Report - SessionStart 훅
# 매 세션 시작 시 어제의 학습 내용 + 현재 상태 + 오늘의 추천 작업을 출력

HOME_DIR="${HOME}"
CLAUDE_DIR="${HOME_DIR}/.claude"
HOOKS_LOG="${CLAUDE_DIR}/.hook-logs/self-improve.log"
AGENT_MEMORY="${CLAUDE_DIR}/agent-memory"
MEMORY_FILE="${CLAUDE_DIR}/projects/-data-data-com-termux-files-home/memory/MEMORY.md"
TODO_FILE="${CLAUDE_DIR}/TODO.md"

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}🌅 Morning Report - $(date '+%Y-%m-%d %H:%M:%S')${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 1. 어제 학습 내용 (Stop 훅이 남긴 로그)
if [ -f "$HOOKS_LOG" ]; then
    echo -e "${GREEN}📚 어제 학습한 규칙:${NC}"
    tail -n 3 "$HOOKS_LOG" | sed 's/^/  /' | head -n 1
    echo ""
else
    echo -e "${YELLOW}⚠️  어제 학습 로그 없음${NC}"
    echo ""
fi

# 2. 최근 수정된 에이전트 메모리
if [ -d "$AGENT_MEMORY" ]; then
    LATEST_AGENT=$(ls -t "$AGENT_MEMORY"/*.md 2>/dev/null | head -n 1)
    if [ -n "$LATEST_AGENT" ]; then
        echo -e "${GREEN}🧠 최근 에이전트 메모리:${NC}"
        echo "  $(basename "$LATEST_AGENT") ($(stat -c %y "$LATEST_AGENT" 2>/dev/null | cut -d' ' -f1,2 || stat -f '%Sm' "$LATEST_AGENT" 2>/dev/null))"
        echo ""
    fi
fi

# 3. 현재 진행 프로젝트 (MEMORY.md 최근 2줄)
if [ -f "$MEMORY_FILE" ]; then
    echo -e "${GREEN}📊 진행 중인 프로젝트:${NC}"
    grep -E "^#|완료|진행|상태" "$MEMORY_FILE" | head -n 4 | sed 's/^/  /'
    echo ""
fi

# 4. 오늘의 추천 작업 (TODO.md 미완료 항목 1개)
if [ -f "$TODO_FILE" ]; then
    echo -e "${GREEN}✅ 오늘의 추천 작업:${NC}"
    grep -E "^\s*-\s*\[" "$TODO_FILE" | grep "\[ \]" | head -n 1 | sed 's/^/  /'
    echo ""
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}💪 좋은 아침입니다! 오늘도 화이팅!${NC}"
echo ""
