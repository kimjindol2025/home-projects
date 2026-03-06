#!/bin/bash
# Bash 작업 검증 훅: 명령어가 실제로 성공했는지 확인

read -r event

COMMAND=$(echo "$event" | jq -r '.tool_input.command' 2>/dev/null)
EXIT_CODE=$(echo "$event" | jq -r '.tool_response.exit_code' 2>/dev/null)
OUTPUT=$(echo "$event" | jq -r '.tool_response.output' 2>/dev/null)

if [ -z "$COMMAND" ]; then
  exit 0
fi

# Claude의 주장 vs 실제 exit code 비교
if echo "$COMMAND" | grep -qiE "git push|git commit|npm test|docker"; then

  # 성공 주장하지만 실패한 경우 감지
  if [ "$EXIT_CODE" -ne 0 ]; then
    echo "🚨 거짓 감지: 명령어 실패"
    echo "  주장: 성공 완료"
    echo "  실제 Exit Code: $EXIT_CODE"
    echo "  명령어: $COMMAND"
    echo "{\"type\": \"bash_fraud\", \"command\": \"$COMMAND\", \"exit_code\": $EXIT_CODE, \"timestamp\": \"$(date -u +%s)\", \"severity\": \"high\"}" >> ~/.claude/verification-audit.log
    exit 1
  else
    echo "✅ Bash 검증 통과: 명령어 실행 성공 (exit $EXIT_CODE)"
    echo "{\"type\": \"bash_verify\", \"command\": \"${COMMAND:0:80}\", \"exit_code\": $EXIT_CODE, \"status\": \"success\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
  fi
fi

exit 0
