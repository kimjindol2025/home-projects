#!/data/data/com.termux/files/usr/bin/bash
# self-healing-test.sh - 실패 테스트 자동 분석 및 수정 제안
# 실행 방식: PostToolUse/Bash hook에서 테스트 명령 실패 감지 시
# 입력: $1 = 테스트 실행 디렉토리, $2 = exit code

REPO_DIR="${1:-$(pwd)}"
EXIT_CODE="${2:-1}"
LOG="$HOME/.claude/self-healing-log.txt"
NOW=$(date '+%Y-%m-%d %H:%M:%S')
REPORT="$HOME/.claude/debug/test-failure-$(date '+%Y%m%d-%H%M%S').md"

# 정상 종료면 아무것도 하지 않음
if [ "$EXIT_CODE" -eq 0 ]; then
    exit 0
fi

echo "[$NOW] SELF-HEALING 시작: exit_code=$EXIT_CODE, dir=$REPO_DIR" >> "$LOG"
echo "=== Self-Healing Test Analysis ==="

# debug 디렉토리 확인 (이미 존재)
mkdir -p "$HOME/.claude/debug" 2>/dev/null || true

# 분석 리포트 생성
cat > "$REPORT" << 'HEREDOC_MARKER'
# 테스트 실패 자동 분석 리포트
HEREDOC_MARKER

echo "**일시**: $NOW" >> "$REPORT"
echo "**디렉토리**: $REPO_DIR" >> "$REPORT"
echo "**Exit Code**: $EXIT_CODE" >> "$REPORT"
echo "" >> "$REPORT"
echo "## 최근 변경 파일 (git diff)" >> "$REPORT"

# git diff 정보 추가
git -C "$REPO_DIR" diff --name-only HEAD 2>/dev/null >> "$REPORT" || echo "(git 정보 없음)" >> "$REPORT"

echo "" >> "$REPORT"
echo "## 프로젝트 타입별 분석" >> "$REPORT"

# 1. Go 테스트 분석
if [ -f "$REPO_DIR/go.mod" ]; then
    echo "" >> "$REPORT"
    echo "### Go 프로젝트" >> "$REPORT"
    echo "실패한 테스트 찾기: \`go test ./... 2>&1 | grep FAIL\`" >> "$REPORT"

    # 일반적인 Go 테스트 실패 패턴 검사
    FAILING_FILES=$(find "$REPO_DIR" -name "*_test.go" 2>/dev/null | head -5)
    for tf in $FAILING_FILES; do
        echo "- 테스트 파일: $tf" >> "$REPORT"
    done
fi

# 2. FreeLang 테스트 분석
FL_TESTS=$(find "$REPO_DIR" -name "*-tests.fl" -o -name "*_test.fl" 2>/dev/null | head -5)
if [ -n "$FL_TESTS" ]; then
    echo "" >> "$REPORT"
    echo "### FreeLang 테스트" >> "$REPORT"
    for ft in $FL_TESTS; do
        echo "- $ft" >> "$REPORT"
    done
fi

# 3. 자가 수복 제안 생성
echo "" >> "$REPORT"
echo "## 자동 수복 제안" >> "$REPORT"
echo "" >> "$REPORT"
echo "다음 명령으로 상세 분석 가능:" >> "$REPORT"
echo "\`\`\`bash" >> "$REPORT"
echo "# Go 테스트 상세 출력" >> "$REPORT"
echo "cd $REPO_DIR && go test -v ./... 2>&1 | tail -50" >> "$REPORT"
echo "" >> "$REPORT"
echo "# 최근 변경사항 확인" >> "$REPORT"
echo "git -C $REPO_DIR log --oneline -5" >> "$REPORT"
echo "git -C $REPO_DIR diff HEAD~1 --stat" >> "$REPORT"
echo "\`\`\`" >> "$REPORT"

# 결과 출력
echo ""
echo "분석 리포트 저장: $REPORT"
echo "[$NOW] REPORT: $REPORT" >> "$LOG"

# Claude에게 리포트 내용 출력 (hook 출력으로 Claude가 읽음)
echo ""
echo "=== 자동 분석 결과 ==="
cat "$REPORT" 2>/dev/null || echo "(리포트 생성 오류)"
