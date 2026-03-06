#!/bin/bash
# 모든 작업 감사 로깅 (비동기)

read -r event

TOOL=$(echo "$event" | jq -r '.tool' 2>/dev/null)
STATUS=$(echo "$event" | jq -r '.tool_response.error // "success"' 2>/dev/null)
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "{\"tool\": \"$TOOL\", \"status\": \"$STATUS\", \"timestamp\": \"$TIMESTAMP\"}" >> ~/.claude/audit.log

exit 0
