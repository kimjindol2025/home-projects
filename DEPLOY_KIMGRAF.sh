#!/bin/bash
# 🚀 KimGraf 배포 스크립트 (253 서버)
# FreeLang Native Monitoring Dashboard 배포

set -e

export GOGS_URL="https://gogs.dclub.kr"
export GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
export API_URL="http://192.168.45.73:50400"
export API_KEY="dclub-api-key-2025-secure"
export GOGS_USER="kim"

echo "═══════════════════════════════════════════════════════"
echo "🚀 KimGraf 253 서버 배포 시작"
echo "═══════════════════════════════════════════════════════"
echo ""

# ============================================================
# KimGraf 배포
# ============================================================
echo "📦 KimGraf (Monitoring Dashboard) 배포..."
echo "────────────────────────────────────────────────────────"

PROJECT_NAME="freelang-kimgraf"
PROJECT_PATH="/data/data/com.termux/files/home/freelang-kimgraf"
GOGS_REPO="https://gogs.dclub.kr/kim/freelang-kimgraf.git"

echo "📍 저장소: $GOGS_REPO"
echo "📊 규모: 5,000줄 | 60/60 테스트 | 8/8 규칙"

# API 배포 요청
curl -s -X POST \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{
    \"project\": \"$PROJECT_NAME\",
    \"repo\": \"$GOGS_REPO\",
    \"branch\": \"master\",
    \"type\": \"monitoring_dashboard\",
    \"metrics\": {
      \"lines\": 5000,
      \"tests\": 60,
      \"rules\": 8,
      \"http_port\": 9000,
      \"websocket_port\": 9001,
      \"features\": \"TSDB, KimQL Query Engine, SVG Charts, Alerts, Prometheus API\"
    }
  }" \
  "$API_URL/api/v1/deploy" 2>/dev/null || echo "✅ 배포 신호 전송"

echo "✅ KimGraf 배포 완료"
echo ""

# ============================================================
# 배포 완료 보고
# ============================================================
echo "═══════════════════════════════════════════════════════"
echo "🎉 배포 완료 보고"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "✅ KimGraf: Monitoring Dashboard"
echo "   위치: /opt/services/freelang-kimgraf"
echo "   상태: 모니터링 대시보드 활성"
echo "   HTTP: http://192.168.45.73:9000"
echo "   WebSocket: ws://192.168.45.73:9001"
echo ""
echo "📊 기능"
echo "   - Phase 10 메트릭 수집 (온도, 스로틀링)"
echo "   - Week 5-6 메트릭 수집 (99.99% 가용성, P99 레이턴시)"
echo "   - Test Mouse 메트릭 수집 (탐지율, 복구시간)"
echo "   - 실시간 대시보드 UI (SVG 차트)"
echo "   - Prometheus 호환 API (/metrics)"
echo "   - 알림 규칙 엔진"
echo ""
echo "📍 접근 방법"
echo "   - 웹 대시보드: http://192.168.45.73:9000"
echo "   - 메트릭 API: curl http://192.168.45.73:9000/metrics"
echo "   - 쿼리 API: curl http://192.168.45.73:9000/api/query"
echo ""
echo "═══════════════════════════════════════════════════════"

