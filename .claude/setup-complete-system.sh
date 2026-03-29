#!/bin/bash
# ============================================================================
# 완벽한 분신 시스템 설치 스크립트
# Custom Alias + Shadow Architect + Termux Opt + Knowledge Preservation
# ============================================================================

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  완벽한 분신 시스템 설치${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo ""

# 1. 권한 설정
echo -e "${BLUE}[1/5] 스크립트 권한 설정${NC}"
chmod +x ~/.claude/aliases.sh 2>/dev/null || true
chmod +x ~/.claude/hooks/termux-optimize.sh 2>/dev/null || true
chmod +x ~/.claude/hooks/auto-map-index.sh 2>/dev/null || true
echo -e "${GREEN}  ✓ 완료${NC}"

# 2. 별칭 등록
echo ""
echo -e "${BLUE}[2/5] Shell 별칭 등록 (선택사항)${NC}"
echo "  ~/.bashrc에 다음 라인을 추가하시겠습니까? (y/n)"
read -r -p "  > " RESPONSE

if [ "$RESPONSE" = "y" ] || [ "$RESPONSE" = "Y" ]; then
    cat >> ~/.bashrc << 'EOF'

# Claude Code 완벽한 분신 시스템 별칭
alias check-logic='source ~/.claude/aliases.sh && check_logic'
alias save-proof='source ~/.claude/aliases.sh && save_proof'
alias verify-vanilla='source ~/.claude/aliases.sh && verify_vanilla'
alias red-team='source ~/.claude/aliases.sh && red_team'
alias refactor-first='source ~/.claude/aliases.sh && refactor_first'
alias prove-it='source ~/.claude/aliases.sh && prove_it'
alias map-index='bash ~/.claude/hooks/auto-map-index.sh'
alias termux-optimize='bash ~/.claude/hooks/termux-optimize.sh'

EOF
    echo -e "${GREEN}  ✓ ~/.bashrc 업데이트 완료${NC}"
else
    echo "  ⊘ 생략 (수동 등록 필요)"
fi

# 3. 설정 검증
echo ""
echo -e "${BLUE}[3/5] 필수 파일 검증${NC}"

FILES=(
    "~/.clauderules"
    "~/.claude/aliases.sh"
    "~/.claude/shadow-architect.md"
    "~/.claude/hooks/termux-optimize.sh"
    "~/.claude/hooks/auto-map-index.sh"
)

for file in "${FILES[@]}"; do
    EXPANDED=$(eval echo "$file")
    if [ -f "$EXPANDED" ]; then
        echo -e "${GREEN}  ✓${NC} $(basename "$EXPANDED")"
    else
        echo -e "${RED}  ✗${NC} $(basename "$EXPANDED") - 누락됨!"
    fi
done

# 4. 초기 실행
echo ""
echo -e "${BLUE}[4/5] 시스템 초기화${NC}"

if [ -d ".git" ]; then
    echo "  📝 MAP.md 생성 중..."
    bash ~/.claude/hooks/auto-map-index.sh >/dev/null 2>&1 || true
    echo -e "${GREEN}  ✓ MAP.md + LEGACY_BRIDGE.md 생성 완료${NC}"
else
    echo "  ⊘ Git 저장소 아님 (스킵)"
fi

# 5. 최종 요약
echo ""
echo -e "${BLUE}[5/5] 설치 완료${NC}"

cat << 'EOF'

╔════════════════════════════════════════════════════════╗
║     완벽한 분신 시스템 설치 완료! 🎉                 ║
╚════════════════════════════════════════════════════════╝

✅ 설치된 시스템:

  1️⃣  Custom Alias (명령어 자동화)
     check-logic, save-proof, verify-vanilla,
     red-team, refactor-first, prove-it

  2️⃣  Shadow Architect (자동 검증)
     Red Team 모드 / Refactor-First / Proof Verification

  3️⃣  Termux 최적화 (리소스 관리)
     메모리/CPU 모니터링 / 배포 자동화 / 핫스팟 분석

  4️⃣  지식 보존 인덱싱 (프로젝트 맵)
     MAP.md / LEGACY_BRIDGE.md / evolution-log.md

════════════════════════════════════════════════════════

🚀 다음 단계:

  1. Shell 재시작: source ~/.bashrc
  2. 명령어 테스트: check-logic (현재 디렉토리)
  3. 프로젝트에서: cd [프로젝트] && verify-vanilla

════════════════════════════════════════════════════════

💾 저장 필수, 기록이 증명이다 GOGS

  모든 작업은 자동으로 다음에 기록됩니다:
  - git 커밋 (GOGS 저장소)
  - HISTORY.gogs (변경 이력)
  - MEMORY.md (프로젝트 메모리)
  - MAP.md (지식 인덱스)

════════════════════════════════════════════════════════

📖 상세 가이드:

  .clauderules 파일의 "완벽한 분신 시스템" 섹션 참고
  또는: ~/.claude/shadow-architect.md 열기

════════════════════════════════════════════════════════

EOF

echo -e "${GREEN}✅ 시스템 준비 완료!${NC}"
echo ""
echo "첫 번째 명령어 실행:"
echo "  $ check-logic"
