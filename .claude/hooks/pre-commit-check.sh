#!/data/data/com.termux/files/usr/bin/bash
# pre-commit-check.sh - 외부 의존성 스캔 (Claude Code PreToolUse hook)
# 실행 방식: Claude hook이 이 스크립트를 bash로 호출
# 입력: $1 = 커밋 중인 git 저장소 디렉토리 (없으면 현재 디렉토리)

REPO_DIR="${1:-$(pwd)}"
LOG="$HOME/.claude/pre-commit-log.txt"
NOW=$(date '+%Y-%m-%d %H:%M:%S')
FAILED=0

echo "[$NOW] PRE-COMMIT CHECK 시작: $REPO_DIR" >> "$LOG"
echo "=== Pre-Commit 의존성 검사 ==="

# 1. go.mod 검사 (Go 프로젝트)
if [ -f "$REPO_DIR/go.mod" ]; then
    echo "[GO] go.mod 의존성 확인 중..."
    # 금지 패키지 검사 (프로젝트별 커스터마이즈)
    BANNED_GO="gorm.io github.com/gin-gonic"
    for pkg in $BANNED_GO; do
        if grep -q "$pkg" "$REPO_DIR/go.mod" 2>/dev/null; then
            echo "[WARN] 금지된 Go 패키지 감지: $pkg"
            echo "[$NOW] WARN: 금지 패키지 $pkg in $REPO_DIR/go.mod" >> "$LOG"
        fi
    done

    # go.sum 존재 확인
    if [ ! -f "$REPO_DIR/go.sum" ]; then
        echo "[WARN] go.sum 없음 - go mod tidy 실행 권고"
    fi
fi

# 2. package.json 검사 (Node.js 프로젝트)
if [ -f "$REPO_DIR/package.json" ]; then
    echo "[NODE] package.json 의존성 확인 중..."
    # python3로 JSON 파싱 (jq 대체)
    python3 -c "
import json, sys

BANNED_NPM = ['moment', 'lodash', 'request']
try:
    with open('$REPO_DIR/package.json') as f:
        pkg = json.load(f)

    all_deps = {}
    all_deps.update(pkg.get('dependencies', {}))
    all_deps.update(pkg.get('devDependencies', {}))

    for dep in BANNED_NPM:
        if dep in all_deps:
            print(f'[WARN] 금지된 npm 패키지: {dep} (경량 대안 권고)')
except Exception as e:
    print(f'[INFO] package.json 파싱 오류: {e}')
" 2>/dev/null
fi

# 3. .fl 파일에서 외부 import 검사 (FreeLang)
FL_FILES=$(find "$REPO_DIR" -name "*.fl" -not -path "*/.*" 2>/dev/null | head -20)
if [ -n "$FL_FILES" ]; then
    echo "[FL] FreeLang 파일 의존성 확인 중..."
    for fl_file in $FL_FILES; do
        # 외부 import 패턴 검사
        if grep -qE 'import\s+"[^.][^/]' "$fl_file" 2>/dev/null; then
            echo "[WARN] 외부 import 감지: $fl_file"
            grep -nE 'import\s+"[^.][^/]' "$fl_file" 2>/dev/null | head -3
        fi
    done
fi

# 4. 하드코딩된 시크릿 검사
echo "[SEC] 시크릿 패턴 검사 중..."
SECRET_PATTERNS='(password|passwd|secret|api_key|apikey|token)\s*=\s*"[^"]{8,}"'
STAGED_FILES=$(git -C "$REPO_DIR" diff --cached --name-only 2>/dev/null)
if [ -n "$STAGED_FILES" ]; then
    for f in $STAGED_FILES; do
        full_path="$REPO_DIR/$f"
        if [ -f "$full_path" ]; then
            if grep -qiE "$SECRET_PATTERNS" "$full_path" 2>/dev/null; then
                echo "[ERROR] 시크릿 패턴 감지: $f"
                echo "[$NOW] ERROR: 시크릿 감지 in $f" >> "$LOG"
                FAILED=1
            fi
        fi
    done
fi

# 5. 결과 요약
echo ""
if [ "$FAILED" -eq 1 ]; then
    echo "=== PRE-COMMIT 검사 실패 - 위 오류를 수정 후 커밋하세요 ==="
    echo "[$NOW] FAILED" >> "$LOG"
    exit 1
else
    echo "=== PRE-COMMIT 검사 통과 ==="
    echo "[$NOW] PASSED" >> "$LOG"
    exit 0
fi
