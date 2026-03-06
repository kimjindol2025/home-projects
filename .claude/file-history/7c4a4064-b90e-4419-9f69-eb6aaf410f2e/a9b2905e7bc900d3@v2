#!/bin/bash
# Write 작업 검증 훅: 파일이 실제로 작성되었는지 확인

read -r event

FILE=$(echo "$event" | jq -r '.tool_input.file_path' 2>/dev/null)
CONTENT=$(echo "$event" | jq -r '.tool_input.content' 2>/dev/null)

if [ -z "$FILE" ]; then
  exit 0
fi

# 파일이 실제로 존재하는지 확인
if [ -f "$FILE" ]; then
  ACTUAL_SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
  CLAIMED_SIZE=${#CONTENT}

  if [ "$ACTUAL_SIZE" -gt 0 ]; then
    echo "✅ Write 검증 통과: $FILE (크기: $ACTUAL_SIZE bytes)"
    echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"success\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
  else
    echo "⚠️ Write 경고: 파일이 비어있음 - $FILE"
    echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"empty\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
  fi
else
  echo "❌ Write 검증 실패: 파일이 없음 - $FILE"
  echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"not_found\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
  exit 1
fi

exit 0
