#!/bin/bash

echo "📊 프로젝트 자동 분석"
echo "===================="
echo ""

for category in core modules experiments tools archived; do
  echo "[$category] 분석 중..."
  
  for project_dir in "$category"/*; do
    if [ -d "$project_dir" ] && [ -f "$project_dir/CLAUDE.md" ]; then
      project=$(basename "$project_dir")
      
      # README 줄 수 확인
      if [ -f "$project_dir/README.md" ]; then
        lines=$(wc -l < "$project_dir/README.md")
      else
        lines=0
      fi
      
      # CLAUDE.md에서 Status 추출
      status=$(grep "Status:" "$project_dir/CLAUDE.md" | head -1 | sed 's/.*Status[: ]*//')
      
      # 우선순위 판정
      if [ "$lines" -lt 50 ]; then
        priority="🔴 HIGH"
      elif [ "$lines" -lt 200 ]; then
        priority="🟡 MED"
      else
        priority="🟢 OK"
      fi
      
      echo "  $priority | $project ($lines줄)"
    fi
  done
  echo ""
done

