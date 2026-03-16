#!/bin/bash

################################################################################
# 🚀 FreeLang Blog Auto-Deploy Script
# 폴더 단위 업로드 + 자동 버전 관리
################################################################################

set -e

PROJECT_NAME="freelang-blog"
COMUPLOAD_URL="https://comupload.ai-empire.kr/api/q"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFEST="$SCRIPT_DIR/MANIFEST.json"

# 색상 정의
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

################################################################################
# 함수: 현재 버전 조회
################################################################################
get_current_version() {
  if [ -f "$MANIFEST" ]; then
    grep -oP '"current_version":\s*\K[0-9]+' "$MANIFEST" | head -1
  else
    echo "0"
  fi
}

################################################################################
# 함수: 아카이브 생성
################################################################################
create_archive() {
  local version=$1
  local archive_name="${PROJECT_NAME}-v${version}.tar"

  echo -e "${BLUE}📦 Creating archive: $archive_name${NC}"
  cd "$SCRIPT_DIR"
  tar -cf "$archive_name" _posts/ README.md *.json 2>/dev/null || true

  if [ -f "$archive_name" ]; then
    echo -e "${GREEN}✅ Archive created: $archive_name${NC}"
    echo "$archive_name"
  else
    echo -e "${YELLOW}⚠️ Archive creation failed${NC}"
    exit 1
  fi
}

################################################################################
# 함수: comupload에 업로드
################################################################################
upload_to_comupload() {
  local archive=$1
  local url="$COMUPLOAD_URL/$PROJECT_NAME"

  echo -e "${BLUE}📤 Uploading to comupload: $url${NC}"

  response=$(curl -s -F "f=@$archive" "$url")

  if echo "$response" | grep -q '"ok":true'; then
    echo -e "${GREEN}✅ Upload successful${NC}"
    echo "$response"
  else
    echo -e "${YELLOW}⚠️ Upload response: $response${NC}"
    echo "$response"
  fi
}

################################################################################
# 함수: MANIFEST 업데이트
################################################################################
update_manifest() {
  local version=$1
  local archive=$2
  local response=$3

  local comupload_id=$(echo "$response" | grep -oP '"id":\s*\K[0-9]+')
  local hash=$(echo "$response" | grep -oP '"hash":"?\K[^"]+')
  local file_count=$(find "$SCRIPT_DIR/_posts" -type f | wc -l)
  local post_count=$(find "$SCRIPT_DIR/_posts" -name "*.md" -type f | wc -l)
  local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  echo -e "${BLUE}📝 Updating MANIFEST.json${NC}"

  # 간단한 JSON 업데이트 (jq 없이)
  # 실제로는 Python이나 jq를 사용하는 것이 좋음
  cat > "$MANIFEST" << EOF
{
  "project": "$PROJECT_NAME",
  "current_version": $version,
  "last_deploy": "$timestamp",
  "status": "active",
  "latest": {
    "version": $version,
    "archive": "$archive",
    "comupload_id": $comupload_id,
    "hash": "$hash",
    "files": $((file_count + 1)),
    "posts": $post_count,
    "date": "$timestamp",
    "gogs_synced": true
  }
}
EOF

  echo -e "${GREEN}✅ MANIFEST updated${NC}"
}

################################################################################
# 함수: 배포 로그 기록
################################################################################
log_deployment() {
  local version=$1
  local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  local log_file="$SCRIPT_DIR/deploy.log"

  echo "$timestamp | Version $version | freelang-blog deployed" >> "$log_file"
  echo -e "${GREEN}✅ Logged to $log_file${NC}"
}

################################################################################
# 메인 배포 로직
################################################################################
main() {
  echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
  echo -e "${BLUE}║ 🚀 FreeLang Blog Auto-Deploy          ║${NC}"
  echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"

  # 현재 버전 확인
  current_version=$(get_current_version)
  next_version=$((current_version + 1))

  echo -e "${YELLOW}Current version: $current_version${NC}"
  echo -e "${YELLOW}Next version: $next_version${NC}\n"

  # 아카이브 생성
  archive=$(create_archive "$next_version")
  echo ""

  # comupload에 업로드
  response=$(upload_to_comupload "$archive")
  echo ""

  # MANIFEST 업데이트
  update_manifest "$next_version" "$archive" "$response"
  echo ""

  # 로그 기록
  log_deployment "$next_version"
  echo ""

  echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║ ✅ Deployment Complete               ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
  echo -e "\nResponse: $response"
}

################################################################################
# 실행
################################################################################
main "$@"
