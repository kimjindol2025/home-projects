#!/bin/bash
# Claude의 주장과 실제 결과 비교 검증

read -r event

TOOL=$(echo "$event" | jq -r '.tool' 2>/dev/null)
MESSAGE=$(echo "$event" | jq -r '.tool_response.message // ""' 2>/dev/null)

case "$TOOL" in
  "Bash")
    COMMAND=$(echo "$event" | jq -r '.tool_input.command' 2>/dev/null)
    EXIT_CODE=$(echo "$event" | jq -r '.tool_response.exit_code' 2>/dev/null)

    # 패턴 1: 성공 주장하지만 실패
    if echo "$MESSAGE" | grep -qiE "success|complete|done|passed"; then
      if [ "$EXIT_CODE" -ne 0 ]; then
        echo "🚨 거짓 감지 (Type 1): 성공 주장하지만 실패"
        echo "{\"type\": \"claim_fraud_type1\", \"claim\": \"$MESSAGE\", \"actual_exit\": $EXIT_CODE, \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
      fi
    fi

    # 패턴 2: 파일 생성 주장하지만 파일 없음
    if echo "$COMMAND" | grep -q "touch\|create\|write"; then
      FILE=$(echo "$COMMAND" | awk '{print $NF}')
      if ! [ -f "$FILE" ] 2>/dev/null; then
        echo "🚨 거짓 감지 (Type 2): 파일 생성 주장하지만 실제 없음"
        echo "{\"type\": \"claim_fraud_type2\", \"file\": \"$FILE\", \"command\": \"$COMMAND\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
      fi
    fi
    ;;

  "Write")
    FILE=$(echo "$event" | jq -r '.tool_input.file_path' 2>/dev/null)

    # 패턴 3: 파일 작성 주장하지만 실제로 없음
    if ! [ -f "$FILE" ] 2>/dev/null; then
      echo "🚨 거짓 감지 (Type 3): 파일 작성 주장하지만 실제 없음"
      echo "{\"type\": \"claim_fraud_type3\", \"file\": \"$FILE\", \"timestamp\": \"$(date -u +%s)\"}" >> ~/.claude/verification-audit.log
    fi
    ;;
esac

exit 0
