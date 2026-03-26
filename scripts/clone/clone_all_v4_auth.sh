#!/bin/bash

TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
GOGS_BASE="https://kim:$TOKEN@gogs.dclub.kr/kim"

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

echo "🚀 FreeLang v4 전체 에코시스템 클론 (인증)..."
COUNT=0

for repo in "${REPOS[@]}"; do
  if [ -d "$repo" ]; then
    echo "✅ $repo (존재)"
  else
    echo "📥 클론: $repo"
    if git clone "$GOGS_BASE/$repo.git" 2>&1 | grep -q "fatal"; then
      echo "❌ $repo"
    else
      echo "✅ $repo"
      ((COUNT++))
    fi
  fi
done

echo ""
echo "✨ 완료! $(ls -d freelang-* 2>/dev/null | wc -l)개 저장소"
