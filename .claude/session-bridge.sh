#!/data/data/com.termux/files/usr/bin/bash
# session-bridge.sh - 세션 종료 시 컨텍스트 저장 및 다음 세션으로 전달
# SessionStop hook에서 실행됨

CLAUDE_DIR="$HOME/.claude"
JOURNAL="$CLAUDE_DIR/JOURNAL.md"
SESSION_LOG="$CLAUDE_DIR/session-log.txt"
TODAY=$(date '+%Y-%m-%d')
NOW=$(date '+%Y-%m-%d %H:%M:%S')
TIME_ONLY=$(date '+%H:%M')

# 1. 세션 종료 기록
echo "[$NOW] SESSION_STOP" >> "$SESSION_LOG"

# 2. 저널에 세션 종료 항목 추가
# 저널이 없으면 생성
if [ ! -f "$JOURNAL" ]; then
    echo "# Claude Code 작업 저널" > "$JOURNAL"
    echo "" >> "$JOURNAL"
fi

# 오늘 날짜 섹션이 없으면 추가
if ! grep -q "## $TODAY" "$JOURNAL" 2>/dev/null; then
    # 파일 앞에 새 섹션 삽입 (최신 항목이 위에 오도록)
    python3 -c "
import sys
journal = '$JOURNAL'
today = '$TODAY'
time_only = '$TIME_ONLY'

try:
    with open(journal, 'r') as f:
        content = f.read()
except:
    content = '# Claude Code 작업 저널\n\n'

# 헤더와 내용 분리
lines = content.split('\n')
header_end = 0
for i, line in enumerate(lines):
    if line.startswith('---') and i > 0:
        header_end = i + 1
        break

new_entry = f'''
## {today} {time_only} - [세션 자동 기록]

**작업**: (이 세션에서 수행한 작업을 기록하세요)
**결과**:
**다음**:

---
'''

# 헤더 다음에 삽입
if header_end > 0:
    new_content = '\n'.join(lines[:header_end]) + new_entry + '\n'.join(lines[header_end:])
else:
    new_content = content + new_entry

with open(journal, 'w') as f:
    f.write(new_content)
" 2>/dev/null
fi

# 3. 작업 로그에 세션 통계 기록
WORK_LOG="$CLAUDE_DIR/work-log.txt"
echo "[$NOW] SESSION_END - 로그 저장 완료" >> "$WORK_LOG"

# 4. 에이전트 메모리 동기화 체크
MEMORY_DIR="$CLAUDE_DIR/agent-memory"
if [ -d "$MEMORY_DIR" ]; then
    echo "[$NOW] 에이전트 메모리 파일 (세션 종료):" >> "$SESSION_LOG"
    ls "$MEMORY_DIR" 2>/dev/null | while read f; do
        echo "  - $f" >> "$SESSION_LOG"
    done
fi

# stdout으로 완료 알림 (Claude가 읽음)
echo "세션 브릿지: 컨텍스트 저장 완료"
