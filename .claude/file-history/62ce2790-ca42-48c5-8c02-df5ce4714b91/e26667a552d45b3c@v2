#!/bin/bash
# ============================================================================
# FreeLang Light - 253.dclub.kr 배포 스크립트
# ============================================================================
# 사용법: bash scripts/deploy-253.sh
# 선행요구: SSH 접속 가능, docker-compose 설치, sudo 권한

set -e

SERVER_HOST="253.dclub.kr"
SERVER_USER="kimjin"
SERVER_PORT="10053"
SERVER_IP="123.212.111.26"
DEPLOY_PATH="/home/kimjin/freelang-light"

echo "🚀 FreeLang Light 프로덕션 배포 시작 (253.dclub.kr)"
echo "=================================================="

# ============================================================================
# Phase 1-1: 환경변수 설정 (30분)
# ============================================================================

echo ""
echo "📋 [Step 1/5] .env 파일 업로드"
echo "----------------------------------------"

scp -P $SERVER_PORT .env.production $SERVER_USER@$SERVER_IP:$DEPLOY_PATH/.env

ssh -p $SERVER_PORT $SERVER_USER@$SERVER_IP << 'REMOTE_STEP1'
  echo "✅ .env 파일 확인:"
  head -10 ~/freelang-light/.env
REMOTE_STEP1

# ============================================================================
# Phase 1-2: Let's Encrypt SSL 인증서 발급 (20분)
# ============================================================================

echo ""
echo "🔒 [Step 2/5] Let's Encrypt 인증서 발급"
echo "----------------------------------------"

ssh -p $SERVER_PORT $SERVER_USER@$SERVER_IP << 'REMOTE_STEP2'
  echo "certbot 설치 확인..."
  if ! command -v certbot &> /dev/null; then
    echo "certbot 설치 중..."
    sudo apt-get update
    sudo apt-get install -y certbot python3-certbot-nginx
  fi

  echo ""
  echo "⏳ Let's Encrypt 인증서 발급..."
  sudo certbot certonly --standalone \
    -d 253.dclub.kr \
    --email noreply@freelang.kr \
    --agree-tos \
    --non-interactive \
    --expand || echo "⚠️  인증서가 이미 존재하거나 발급 실패. 수동 확인 필요."

  echo ""
  echo "✅ 인증서 확인:"
  sudo ls -la /etc/letsencrypt/live/253.dclub.kr/
REMOTE_STEP2

# ============================================================================
# Phase 1-3: Docker 이미지 빌드 및 서비스 시작 (30분)
# ============================================================================

echo ""
echo "🐳 [Step 3/5] Docker Compose 빌드 & 시작"
echo "----------------------------------------"

ssh -p $SERVER_PORT $SERVER_USER@$SERVER_IP << 'REMOTE_STEP3'
  cd ~/freelang-light

  echo "Docker daemon 상태 확인..."
  sudo systemctl is-active docker || sudo systemctl start docker

  echo ""
  echo "⏳ Docker 이미지 빌드 중 (2-3분 소요)..."
  docker-compose build --no-cache

  echo ""
  echo "⏳ 서비스 시작 중..."
  docker-compose up -d

  echo ""
  echo "✅ 컨테이너 상태:"
  docker-compose ps
REMOTE_STEP3

# ============================================================================
# Phase 1-4: 서비스 헬스 체크 (10분)
# ============================================================================

echo ""
echo "🏥 [Step 4/5] 서비스 헬스 체크"
echo "----------------------------------------"

echo "⏳ 서비스 시작 대기 중... (30초)"
sleep 30

ssh -p $SERVER_PORT $SERVER_USER@$SERVER_IP << 'REMOTE_STEP4'
  echo "API 헬스 체크:"
  curl -s http://localhost:3001/api/health | jq . || echo "⚠️  API 응답 없음"

  echo ""
  echo "PostgreSQL 연결 확인:"
  docker-compose exec -T postgres psql -U freelang -d freelang -c "SELECT COUNT(*) FROM posts;" || echo "⚠️  DB 연결 실패"

  echo ""
  echo "Redis 상태 확인:"
  docker-compose exec -T redis redis-cli PING || echo "⚠️  Redis 연결 실패"

  echo ""
  echo "Prometheus 메트릭 수집:"
  curl -s http://localhost:9090/api/v1/query?query=up | jq '.data.result | length' || echo "⚠️  Prometheus 응답 없음"
REMOTE_STEP4

# ============================================================================
# Phase 1-5: HTTPS 검증 (5분)
# ============================================================================

echo ""
echo "🔐 [Step 5/5] HTTPS 검증"
echo "----------------------------------------"

sleep 10

echo "SSL 인증서 확인:"
openssl s_client -connect 253.dclub.kr:443 -servername 253.dclub.kr < /dev/null 2>/dev/null | grep -A 2 "Issuer:" || echo "⚠️  HTTPS 응답 없음"

echo ""
echo "HTTPS 연결 테스트:"
curl -s -I https://253.dclub.kr/api/health | head -5

# ============================================================================
# 배포 완료
# ============================================================================

echo ""
echo "=================================================="
echo "✅ Phase 1 배포 완료!"
echo "=================================================="
echo ""
echo "📊 접근 가능한 서비스:"
echo "  • API:          https://253.dclub.kr/api/health"
echo "  • Prometheus:   http://253.dclub.kr:9090"
echo "  • Grafana:      http://253.dclub.kr:3100 (admin/admin123)"
echo ""
echo "📝 다음 단계:"
echo "  1. 로그 확인: ssh -p 10053 kimjin@253.dclub.kr 'docker-compose logs -f blog'"
echo "  2. 백업 테스트: ssh -p 10053 kimjin@253.dclub.kr 'bash scripts/backup.sh'"
echo "  3. 모니터링: https://253.dclub.kr:9090/graph"
echo ""
