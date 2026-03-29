#!/bin/bash
# Gogs-Pulse: git post-commit 훅
# 커밋 후 자동으로 GOGS wiki에 동기화

set -e

HOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 색상 정의
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${YELLOW}📡 Gogs-Pulse: Syncing commit to GOGS wiki...${NC}"

# FreeLang 또는 Go로 실행
if command -v freelang &>/dev/null; then
    freelang run "$HOOK_DIR/main.fl" sync
else
    # Go 바이너리 사용
    "$HOOK_DIR/main" sync 2>/dev/null || true
fi

echo -e "${GREEN}✅ Sync complete${NC}"
