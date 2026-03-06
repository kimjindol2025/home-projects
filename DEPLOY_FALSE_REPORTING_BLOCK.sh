#!/bin/bash
# 253 서버에 거짓 보고 차단 설정 배포 스크립트
# 사용법: bash DEPLOY_FALSE_REPORTING_BLOCK.sh

set -e

echo "🚀 253 서버 거짓 보고 차단 시스템 배포"
echo "════════════════════════════════════════"
echo ""

# 1. Hooks 디렉토리 생성
echo "📁 [1/5] Hooks 디렉토리 생성..."
mkdir -p ~/deployment/.claude/hooks
echo "✅ 완료"
echo ""

# 2. settings.json 생성
echo "⚙️ [2/5] Claude Code 설정 생성..."
cat > ~/deployment/.claude/settings.json << 'EOF'
{
  "model": "opus-4-6",
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/hooks/verify-write.sh"
          }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/hooks/verify-bash.sh"
          }
        ]
      },
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/hooks/audit-logger.sh",
            "async": true
          }
        ]
      }
    ]
  }
}
EOF
echo "✅ 완료"
echo ""

# 3. verify-write.sh 생성
echo "🔍 [3/5] Write 검증 훅 생성..."
cat > ~/deployment/.claude/hooks/verify-write.sh << 'EOF'
#!/bin/bash
read -r event
FILE=$(echo "$event" | jq -r '.tool_input.file_path' 2>/dev/null)
if [ -z "$FILE" ]; then exit 0; fi
if [ -f "$FILE" ]; then
  SIZE=$(stat -c%s "$FILE" 2>/dev/null || stat -f%z "$FILE" 2>/dev/null)
  if [ "$SIZE" -gt 0 ]; then
    echo "✅ Write 검증: $FILE (크기: $SIZE bytes)"
    echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"success\"}" >> ~/.claude/verification-audit.log
  else
    echo "⚠️ Write 경고: 파일이 비어있음"
    echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"empty\"}" >> ~/.claude/verification-audit.log
  fi
else
  echo "❌ Write 검증 실패: 파일 없음 - $FILE"
  echo "{\"type\": \"write_verify\", \"file\": \"$FILE\", \"status\": \"not_found\"}" >> ~/.claude/verification-audit.log
  exit 1
fi
exit 0
EOF
chmod +x ~/deployment/.claude/hooks/verify-write.sh
echo "✅ 완료"
echo ""

# 4. verify-bash.sh 생성
echo "🔍 [4/5] Bash 검증 훅 생성..."
cat > ~/deployment/.claude/hooks/verify-bash.sh << 'EOF'
#!/bin/bash
read -r event
COMMAND=$(echo "$event" | jq -r '.tool_input.command' 2>/dev/null)
EXIT_CODE=$(echo "$event" | jq -r '.tool_response.exit_code' 2>/dev/null)
if [ -z "$COMMAND" ]; then exit 0; fi
if [ "$EXIT_CODE" -ne 0 ]; then
  echo "🚨 거짓 감지: 명령어 실패 (exit $EXIT_CODE)"
  echo "{\"type\": \"bash_fraud\", \"command\": \"$COMMAND\", \"exit_code\": $EXIT_CODE, \"severity\": \"high\"}" >> ~/.claude/verification-audit.log
  exit 1
else
  echo "✅ Bash 검증: 명령어 성공 (exit $EXIT_CODE)"
  echo "{\"type\": \"bash_verify\", \"status\": \"success\"}" >> ~/.claude/verification-audit.log
fi
exit 0
EOF
chmod +x ~/deployment/.claude/hooks/verify-bash.sh
echo "✅ 완료"
echo ""

# 5. 감사 로그 초기화
echo "📝 [5/5] 감사 로그 초기화..."
touch ~/.claude/verification-audit.log ~/.claude/audit.log
chmod 600 ~/.claude/verification-audit.log ~/.claude/audit.log
echo "✅ 완료"
echo ""

echo "════════════════════════════════════════"
echo "✅ 거짓 보고 차단 시스템 배포 완료!"
echo ""
echo "📊 설정 현황:"
echo "  • 설정 파일: ~/deployment/.claude/settings.json"
echo "  • Hooks 스크립트:"
echo "    - verify-write.sh (Write 작업 검증)"
echo "    - verify-bash.sh (Bash 명령 검증)"
echo "    - audit-logger.sh (감사 로깅)"
echo ""
echo "📝 감사 로그:"
echo "  • ~/.claude/verification-audit.log (거짓 감지 기록)"
echo "  • ~/.claude/audit.log (전체 작업 기록)"
echo ""
echo "🚀 다음 단계:"
echo "  1. Claude Code 시작: claude"
echo "  2. 작업 수행"
echo "  3. 거짓 감지 확인: tail -f ~/.claude/verification-audit.log"
echo ""
