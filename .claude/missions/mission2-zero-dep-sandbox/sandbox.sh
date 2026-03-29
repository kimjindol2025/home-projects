#!/bin/bash
# Zero-Dep-Sandbox: git pre-commit 훅
# 외부 의존성 검증 후 커밋 차단

set -e

HOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(git rev-parse --show-toplevel)"

# FreeLang 런타임이 없으면 Go로 구성된 버전 사용
if command -v freelang &>/dev/null; then
    SANDBOX="freelang run $HOOK_DIR/main.fl"
else
    SANDBOX="$HOOK_DIR/main.go"
fi

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🔒 Zero-Dep-Sandbox: Checking external dependencies...${NC}"

# 검증 실행
if $SANDBOX check "$PROJECT_ROOT"; then
    echo -e "${GREEN}✅ PASS: No external dependencies detected${NC}"
    exit 0
else
    echo -e "${RED}❌ FAIL: External dependencies found (commit blocked)${NC}"
    echo ""
    echo -e "${YELLOW}To bypass (NOT RECOMMENDED):${NC}"
    echo "  git commit --no-verify"
    echo ""
    echo -e "${YELLOW}To fix:${NC}"
    echo "  1. Review the external imports"
    echo "  2. Move them to 'internal/' or 'vendor/'"
    echo "  3. Try committing again"
    echo ""
    exit 1
fi
