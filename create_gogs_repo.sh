#!/bin/bash

# GOGS 설정
GOGS_URL="http://localhost:3000"
USERNAME="kim"

# 저장소 이름 목록
REPOS=(
  "freelang-unified-build-pipeline"
  "freelang-standard-core"
  "freelang-z3-verification"
)

echo "🚀 GOGS 저장소 생성 스크립트"
echo "=================================="
echo ""

# 각 저장소 정보 출력
for repo in "${REPOS[@]}"; do
  echo "📦 $repo"
  echo "   URL: $GOGS_URL/$USERNAME/$repo.git"
  echo "   설명: Phase 1-3 분석 결과 저장소"
  echo ""
done

echo "=================================="
echo ""
echo "생성할 저장소 이름:"
echo ""
for i in "${!REPOS[@]}"; do
  echo "$((i+1)). ${REPOS[$i]}"
done
echo ""
echo "각 저장소는 다음을 포함합니다:"
echo "  • Phase 별 분석 보고서"
echo "  • 자동화 스크립트"
echo "  • 산출물 (CSV, JSON)"
echo "  • 구현 로드맵"
echo ""

