#!/bin/bash
# Self-Improving: Stop 훅에서 자동 실행
# 세션 종료 후 취약점 분석 → .clauderules 자동 업데이트

set -e

HOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOME_DIR="$HOME"
LOG_FILE="$HOME_DIR/.claude/.hook-logs/self-improve.log"

# 로그 디렉토리 생성
mkdir -p "$(dirname "$LOG_FILE")"

# 색상 정의
CYAN='\033[0;36m'
GREEN='\033[0;32m'
NC='\033[0m'

# 훅 시작 로깅
{
  echo "=========================================="
  echo "🧠 Self-Improving Hook Executed"
  echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
  echo "Session ID: $CLAUDE_SESSION_ID"
  echo "=========================================="
} >> "$LOG_FILE"

echo -e "${CYAN}🧠 Self-Improving: Analyzing session patterns...${NC}"

# FreeLang 또는 Go로 실행
if command -v freelang &>/dev/null; then
    echo "✓ Using FreeLang runtime" >> "$LOG_FILE"
    freelang run "$HOOK_DIR/main.fl" analyze >> "$LOG_FILE" 2>&1 || echo "⚠ analyze failed" >> "$LOG_FILE"
    freelang run "$HOOK_DIR/main.fl" update-rules >> "$LOG_FILE" 2>&1 || echo "⚠ update-rules failed" >> "$LOG_FILE"
else
    # Go 바이너리 사용
    echo "✓ Using Go binary" >> "$LOG_FILE"
    "$HOOK_DIR/main" analyze >> "$LOG_FILE" 2>&1 || echo "⚠ analyze failed" >> "$LOG_FILE"
    "$HOOK_DIR/main" update-rules >> "$LOG_FILE" 2>&1 || echo "⚠ update-rules failed" >> "$LOG_FILE"
fi

# 훅 완료 로깅
{
  echo "Status: COMPLETED"
  echo "Log location: $LOG_FILE"
  echo "=========================================="
  echo ""
} >> "$LOG_FILE"

echo -e "${GREEN}✅ Patterns updated in .clauderules${NC}"
echo -e "${GREEN}📝 Log: $LOG_FILE${NC}"
