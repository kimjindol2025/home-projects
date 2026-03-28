#!/data/data/com.termux/files/usr/bin/bash
# context-injector.sh - SessionStart 시 이전 작업 컨텍스트 로드
# Claude Code의 SessionStart hook에서 실행됨
# 출력은 Claude가 세션 시작 시 컨텍스트로 읽음

CLAUDE_DIR="$HOME/.claude"
JOURNAL="$CLAUDE_DIR/JOURNAL.md"
MEMORY="$CLAUDE_DIR/MEMORY.md"
TODAY=$(date '+%Y-%m-%d')
LAST_RUN="$CLAUDE_DIR/.context-injector-last-run"

# 하루 1회 실행만 (같은 날에 반복 실행 방지)
if [ -f "$LAST_RUN" ] && grep -q "$TODAY" "$LAST_RUN" 2>/dev/null; then
    exit 0
fi

echo "=== SESSION CONTEXT INJECTION: $TODAY ==="
echo ""

# 1. 최근 저널 항목 (마지막 3개 작업)
if [ -f "$JOURNAL" ]; then
    echo "## 📋 최근 작업 이력 (JOURNAL 최신 3개)"
    echo ""
    # python3로 마지막 ## 섹션 3개 추출
    python3 -c "
import re, sys
try:
    content = open('$JOURNAL').read()
    sections = re.findall(r'(## \d{4}-\d{2}-\d{2}.*?)(?=## \d{4}-\d{2}-\d{2}|\Z)', content, re.DOTALL)
    for s in sections[-3:]:
        lines = s.strip().split('\n')[:10]  # 각 섹션 첫 10줄만
        print('\n'.join(lines))
        print()
except Exception as e:
    print(f'(저널 파싱 오류: {e})')
" 2>/dev/null || echo "(저널 파싱 오류)"
fi

# 2. 현재 MEMORY.md 상태 요약 (첫 50줄)
if [ -f "$MEMORY" ]; then
    echo "## 🧠 프로젝트 메모리 (현재 상태 스냅샷)"
    echo ""
    head -n 50 "$MEMORY" | head -n 40
    echo ""
    echo "... (더 보기: cat ~/.claude/MEMORY.md)"
fi

# 3. 오늘 날짜 컨텍스트
echo ""
echo "## 📅 시스템 정보"
echo "- 세션 시작: $(date '+%Y-%m-%d %H:%M:%S')'
echo "- 작업 디렉토리: $(pwd)"

# 4. 미완료 TODO 확인
TODO="$CLAUDE_DIR/TODO.md"
if [ -f "$TODO" ]; then
    PENDING=$(grep -c '^\- \[ \]' "$TODO" 2>/dev/null || echo "0")
    [ "$PENDING" != "0" ] && echo "- 미완료 TODO: ${PENDING}개"
fi

# 5. 세션 시작 로그 기록
echo ""
echo "[$(date '+%Y-%m-%d %H:%M:%S')] SESSION_START" >> "$CLAUDE_DIR/session-log.txt"

# 6. 마지막 실행 시간 기록 (하루 1회 체크용)
date '+%Y-%m-%d' > "$LAST_RUN"

echo "=== CONTEXT INJECTION 완료 ==="
