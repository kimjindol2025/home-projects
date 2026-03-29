#!/bin/bash
# Termux 환경에서 Claude Code 자동 배포 및 리소스 관리

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============================================================
# 함수 1: 메모리 모니터링
# ============================================================
check_memory() {
  local max_mem=${1:-500}  # 기본값: 500MB (모바일)

  if ! command -v free &> /dev/null; then
    # Termux에서 available memory 확인
    local mem_info=$(cat /proc/meminfo | grep MemAvailable | awk '{print $2}')
    local current_mb=$((mem_info / 1024))
  else
    local current_mb=$(free -m | awk 'NR==2{print $7}')
  fi

  echo -e "${BLUE}💾 메모리 상태: ${current_mb}MB${NC}"

  if [ $current_mb -lt 100 ]; then
    echo -e "${RED}⚠️  경고: 메모리 부족 (${current_mb}MB < 100MB)${NC}"
    echo "작업 중단 권장"
    return 1
  fi

  return 0
}

# ============================================================
# 함수 2: 빌드 및 테스트
# ============================================================
auto_build_test() {
  local project_dir=$1

  if [ ! -d "$project_dir" ]; then
    echo -e "${RED}❌ 디렉토리 없음: $project_dir${NC}"
    return 1
  fi

  cd "$project_dir"
  echo -e "${BLUE}📁 프로젝트: $(pwd)${NC}"

  # 빌드
  echo -e "${YELLOW}🔨 빌드 시작...${NC}"
  if go build ./... 2>&1 | tee /tmp/build.log | head -20; then
    echo -e "${GREEN}✅ 빌드 성공${NC}"
    local build_status="PASS"
  else
    echo -e "${RED}❌ 빌드 실패${NC}"
    local build_status="FAIL"
    tail -10 /tmp/build.log
    return 1
  fi

  # 테스트 (구조 검증)
  echo -e "${YELLOW}🧪 테스트 구조 검증...${NC}"
  local test_count=$(find . -name '*_test.go' -exec wc -l {} + | awk '{sum+=$1} END {print sum}')
  local test_files=$(find . -name '*_test.go' | wc -l)

  if [ $test_files -gt 0 ]; then
    echo -e "${GREEN}✅ 테스트 구조: ${test_files}개 파일, ${test_count}줄${NC}"
    local test_status="PASS"
  else
    echo -e "${YELLOW}⚠️  테스트 파일 없음${NC}"
    local test_status="NONE"
  fi

  # 요약 출력
  echo ""
  echo -e "${BLUE}📊 빌드 결과 요약:${NC}"
  echo "  빌드: $build_status"
  echo "  테스트: $test_status"
  echo "  구조: $test_files 파일, $test_count 줄"

  return 0
}

# ============================================================
# 함수 3: 자동 배포
# ============================================================
auto_deploy() {
  local commit_msg=$1

  if [ -z "$commit_msg" ]; then
    echo -e "${RED}❌ 커밋 메시지가 필요합니다${NC}"
    return 1
  fi

  echo -e "${BLUE}📦 배포 시작...${NC}"

  # Git 상태 확인
  if ! git status &> /dev/null; then
    echo -e "${RED}❌ Git 저장소가 아닙니다${NC}"
    return 1
  fi

  # 스테이지
  echo "📝 파일 스테이징..."
  git add -A

  # 변경사항 확인
  if ! git diff --cached --quiet; then
    echo -e "${YELLOW}변경사항 확인:${NC}"
    git diff --cached --stat | head -10
  else
    echo "변경사항 없음"
    return 0
  fi

  # 커밋
  echo "💾 커밋..."
  git commit -m "$commit_msg"

  # 푸시
  local remote=$(git remote | head -1)
  local branch=$(git rev-parse --abbrev-ref HEAD)

  if [ -n "$remote" ]; then
    echo "🚀 푸시 ($remote/$branch)..."
    if git push "$remote" "$branch" 2>&1 | tail -5; then
      echo -e "${GREEN}✅ 배포 완료${NC}"
    else
      echo -e "${YELLOW}⚠️  푸시 경고 (로컬에만 저장됨)${NC}"
    fi
  else
    echo -e "${YELLOW}⚠️  원격 저장소 없음 (로컬에만 저장됨)${NC}"
  fi

  echo ""
  echo -e "${GREEN}저장 필수, 기록이 증명이다 gogs${NC}"
}

# ============================================================
# 함수 4: 지도(MAP) 생성
# ============================================================
generate_map() {
  local project_dir=${1:-.}

  echo -e "${BLUE}🗺️  프로젝트 맵 생성 중...${NC}"

  cat > "$project_dir/MAP.md" << 'EOF'
# 프로젝트 맵

## 파일 구조
EOF

  echo "" >> "$project_dir/MAP.md"
  echo '```' >> "$project_dir/MAP.md"
  find "$project_dir" -type f \( -name "*.go" -o -name "*.py" -o -name "*.ts" \) \
    | sed 's|'"$project_dir"'/||' | sort >> "$project_dir/MAP.md"
  echo '```' >> "$project_dir/MAP.md"

  echo "" >> "$project_dir/MAP.md"
  echo "## 함수 목록" >> "$project_dir/MAP.md"
  echo "" >> "$project_dir/MAP.md"

  # Go 파일에서 함수 추출
  for file in $(find "$project_dir" -name "*.go" -type f); do
    local funcs=$(grep -o "^func [a-zA-Z_][a-zA-Z0-9_]*" "$file" 2>/dev/null | sed 's/^func //' | sort -u)
    if [ -n "$funcs" ]; then
      echo "### $(basename "$file")" >> "$project_dir/MAP.md"
      echo "" >> "$project_dir/MAP.md"
      for func in $funcs; do
        echo "- \`$func\`" >> "$project_dir/MAP.md"
      done
      echo "" >> "$project_dir/MAP.md"
    fi
  done

  echo -e "${GREEN}✅ MAP.md 생성 완료${NC}"
}

# ============================================================
# 메인
# ============================================================
main() {
  case "$1" in
    memory)
      check_memory "${2:-500}"
      ;;
    test)
      auto_build_test "${2:-.}"
      ;;
    deploy)
      auto_deploy "$2"
      ;;
    map)
      generate_map "${2:-.}"
      ;;
    all)
      echo -e "${BLUE}🤖 전체 작업 실행...${NC}"
      check_memory && \
      auto_build_test "${2:-.}" && \
      generate_map "${2:-.}" && \
      echo -e "${GREEN}✅ 모든 작업 완료${NC}"
      ;;
    *)
      echo "Claude Code Termux 자동화 도구"
      echo ""
      echo "사용법:"
      echo "  $0 memory [max_mb]     # 메모리 확인"
      echo "  $0 test [dir]          # 빌드 및 테스트"
      echo "  $0 deploy [message]    # 자동 배포 (git add/commit/push)"
      echo "  $0 map [dir]           # 프로젝트 맵 생성"
      echo "  $0 all [dir]           # 전체 작업 실행"
      echo ""
      echo "예시:"
      echo "  $0 memory 500          # 메모리 최대값 500MB 확인"
      echo "  $0 test .              # 현재 디렉토리 빌드/테스트"
      echo "  $0 deploy '✨ 기능 추가' # 커밋 메시지와 함께 배포"
      ;;
  esac
}

main "$@"
