#!/bin/bash

TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
GOGS_URL="https://gogs.dclub.kr/kim"

# 모든 v4 모듈 목록
REPOS=(
  "freelang-v4-sqlite-integration"
  "freelang-streaming-arena"
  "freelang-v4-transaction-advanced"
  "freelang-v4-migration"
  "freelang-v4-orm"
  "freelang-v4-query-performance"
  "freelang-v4-audit-system"
  "freelang-v4-bytecode-cache"
  "freelang-v4-compiler-optimizer"
  "freelang-v4-crypto"
  "freelang-v4-input-validation"
  "freelang-v4-security"
  "freelang-v4-compliance"
  "freelang-v4-object-map"
  "freelang-regex-engine"
  "freelang-database-functions"
  "freelang-v4-concurrency"
)

echo "🚀 FreeLang v4 전체 에코시스템 클론 시작..."
echo ""

for repo in "${REPOS[@]}"; do
  if [ -d "$repo" ]; then
    echo "✅ $repo (이미 존재)"
  else
    echo "📥 $repo 클론 중..."
    git clone "$GOGS_URL/$repo.git" 2>/dev/null && echo "✅ $repo 완료" || echo "❌ $repo 실패"
  fi
done

echo ""
echo "✨ 완료!"
ls -d freelang-* | wc -l
echo "개의 freelang 저장소 설치됨"
