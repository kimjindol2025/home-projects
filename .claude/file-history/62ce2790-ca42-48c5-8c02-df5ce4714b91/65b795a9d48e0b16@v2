# 🚀 FreeLang Light - 프로덕션 배포 가이드

**버전**: Production Deployment (Phase A)
**대상 서버**: 253.dclub.kr (IP: 123.212.111.26:10053)
**상태**: 배포 준비 완료
**배포 일정**: 2026-03-13 (즉시 가능)

---

## 📋 배포 전 체크리스트

### 1단계: 환경 변수 설정 (30분)

```bash
# .env 파일 생성
cd /tmp/freelang-light
cat > .env << 'EOF'
# 서버 설정
SERVER_HOST=0.0.0.0
SERVER_PORT=5021
ENVIRONMENT=production
LOG_LEVEL=info

# PostgreSQL
DATABASE_URL=postgresql://freelang:secure_password@postgres:5432/freelang
DB_BACKUP_ENABLED=true
DB_BACKUP_INTERVAL=86400  # 1일

# Redis
REDIS_URL=redis://redis:6379
REDIS_TTL=3600
CACHE_ENABLED=true

# SSL/TLS
SSL_ENABLED=true
SSL_CERT_PATH=/etc/letsencrypt/live/253.dclub.kr/cert.pem
SSL_KEY_PATH=/etc/letsencrypt/live/253.dclub.kr/key.pem

# JWT 인증
JWT_SECRET=your-secure-jwt-secret-key-here
JWT_EXPIRY=3600
REFRESH_TOKEN_EXPIRY=2592000  # 30일

# OAuth (Google, GitHub, Naver)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
NAVER_CLIENT_ID=your-naver-client-id
NAVER_CLIENT_SECRET=your-naver-client-secret

# Monitoring
PROMETHEUS_ENABLED=true
GRAFANA_ADMIN_PASSWORD=admin123
JAEGER_AGENT_HOST=localhost
JAEGER_AGENT_PORT=6831

# Email (백업 알림)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@freelang.kr
SMTP_PASSWORD=your-app-password

# API Gateway
GATEWAY_PORT=8000
RATE_LIMIT_RPS=100
RATE_LIMIT_HOUR=100000

# API 문서
API_DOCS_ENABLED=true
SWAGGER_UI_PATH=/api/docs

# Slack 알림
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
SLACK_CHANNEL=#alerts

EOF
```

### 2단계: GitHub Secrets 설정 (15분)

```bash
# GitHub 저장소 Settings → Secrets and variables → Actions

# 필수 Secrets:
DEPLOY_KEY=your-ssh-private-key
DEPLOY_HOST=253.dclub.kr
DEPLOY_PORT=10053
DEPLOY_USER=kimjin
DEPLOY_PATH=/home/kimjin/freelang-light

# 추가 Secrets:
DB_PASSWORD=secure_password
JWT_SECRET=your-jwt-secret
GOOGLE_CLIENT_SECRET=xxx
```

### 3단계: SSL 인증서 설정 (20분)

```bash
# 서버에서 Let's Encrypt 인증서 발급
ssh -p 10053 kimjin@253.dclub.kr

# Certbot 설치
sudo apt-get install certbot python3-certbot-nginx

# 인증서 발급
sudo certbot certonly --standalone \
  -d 253.dclub.kr \
  --email noreply@freelang.kr \
  --agree-tos \
  --non-interactive

# 자동 갱신 설정
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer

# 확인
sudo ls -la /etc/letsencrypt/live/253.dclub.kr/
```

### 4단계: 데이터베이스 백업 확인 (15분)

```bash
# 로컬에서 최신 백업 생성
bash scripts/backup.sh

# 확인
ls -lh /backups/*.sql.gz | tail -5

# S3 업로드 (선택)
aws s3 cp /backups/ s3://freelang-backups/ --recursive
```

### 5단계: 모니터링 알림 테스트 (20분)

```bash
# Prometheus 데이터 소스 추가
curl -X POST http://localhost:9090/api/v1/admin/tsdb/delete_series \
  -H "Content-Type: application/json" \
  -d '{"match": ["test_metric"]}'

# Grafana 대시보드 확인
open http://localhost:3100
# 기본 계정: admin / admin123

# Slack 알림 테스트
curl -X POST $SLACK_WEBHOOK_URL \
  -H 'Content-Type: application/json' \
  -d '{"text": "Test alert from FreeLang"}'
```

---

## 🚀 배포 절차

### Phase A-1: 스테이징 환경 테스트 (1시간)

```bash
# 1. 로컬에서 Docker Compose 테스트
docker-compose -f docker-compose.yml build
docker-compose -f docker-compose.yml up -d

# 2. 서비스 헬스 확인
curl http://localhost:5021/api/health
curl http://localhost:8000/api/gateway/health

# 3. API 엔드포인트 테스트
curl -X GET http://localhost:5021/api/posts
curl -X GET http://localhost:8000/api/v1/posts

# 4. 인증 테스트
TOKEN=$(curl -X POST http://localhost:5021/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"test","password":"test"}' | jq -r '.token')

curl -X GET http://localhost:5021/api/posts \
  -H "Authorization: Bearer $TOKEN"

# 5. 모니터링 대시보드 확인
open http://localhost:9090  # Prometheus
open http://localhost:3100  # Grafana
open http://localhost:16686 # Jaeger

# 6. 성능 테스트
ab -n 1000 -c 10 http://localhost:5021/api/posts
```

### Phase A-2: 서버 배포 (1시간)

```bash
# 1. 서버 접속
ssh -p 10053 kimjin@253.dclub.kr

# 2. 저장소 클론
cd ~
git clone https://github.com/your-user/freelang-light.git
cd freelang-light

# 3. 환경 변수 복사
# (사전에 준비한 .env 파일 업로드)
scp -P 10053 .env kimjin@253.dclub.kr:~/freelang-light/

# 4. Docker 이미지 빌드
docker-compose build --no-cache

# 5. 마이그레이션 실행
docker-compose run --rm blog \
  freec run scripts/migrate.fl

# 6. 서비스 시작
docker-compose -f docker-compose.yml up -d

# 7. 시작 로그 확인
docker-compose logs -f blog
# 또는 services 상태 확인
docker-compose ps
```

### Phase A-3: 검증 및 모니터링 (30분)

```bash
# 1. 서비스 가용성 확인
curl https://253.dclub.kr/api/health

# 2. SSL 인증서 확인
openssl s_client -connect 253.dclub.kr:443 \
  -servername 253.dclub.kr < /dev/null

# 3. 데이터베이스 연결 확인
docker-compose exec postgres psql -U freelang -d freelang -c "SELECT COUNT(*) FROM posts;"

# 4. 캐시 작동 확인
docker-compose exec redis redis-cli PING

# 5. 메트릭 수집 확인
curl http://localhost:9090/api/v1/query?query=up
# 또는 Prometheus UI: http://253.dclub.kr:9090

# 6. 로그 수집 확인
docker-compose logs blog | tail -50

# 7. 모니터링 대시보드 접근
# (nginx를 통해 443 포트로 프록시)
open https://253.dclub.kr/monitoring/prometheus
open https://253.dclub.kr/monitoring/grafana
```

---

## 📊 배포 후 모니터링

### 1단계: 실시간 메트릭 (첫 1시간)

```bash
# Prometheus 대시보드 열기
curl https://253.dclub.kr:9090

# 모니터할 주요 메트릭:
http_requests_total[5m]       # 요청률
http_request_duration_seconds # 응답시간
up{job="blog"}                # 서비스 가용성
postgres_connections         # DB 연결
redis_connected_clients       # 캐시 연결
```

### 2단계: 알림 규칙 검증 (1시간 후)

```yaml
# Prometheus 알림 규칙 (/etc/prometheus/alert.rules.yml)

# 규칙 1: 서비스 다운
- alert: ServiceDown
  expr: up{job="blog"} == 0
  for: 1m
  annotations:
    summary: "Blog service down"

# 규칙 2: 높은 에러율
- alert: HighErrorRate
  expr: (rate(http_requests_total{status=~"5.."}[5m])) > 0.05
  for: 5m
  annotations:
    summary: "Error rate > 5%"

# 규칙 3: DB 연결 풀 소진
- alert: DBPoolExhausted
  expr: postgres_connections > 90
  for: 2m
```

### 3단계: 비즈니스 메트릭 (일일)

```
대시보드: https://253.dclub.kr/monitoring/grafana

모니터할 지표:
├─ 성능
│  ├─ 평균 응답시간 (<200ms)
│  ├─ P95 레이턴시 (<500ms)
│  └─ RPS (요청/초)
│
├─ 안정성
│  ├─ 가용성 (>99.95%)
│  ├─ 에러율 (<0.1%)
│  └─ 에러 유형별 분포
│
├─ 리소스
│  ├─ CPU 사용률 (<70%)
│  ├─ 메모리 사용률 (<80%)
│  ├─ 디스크 사용률 (<70%)
│  └─ 네트워크 대역폭
│
└─ 비즈니스
   ├─ 사용자 수
   ├─ 포스트 수
   ├─ 활성 세션
   └─ API 호출 수
```

---

## 🔄 배포 후 관리

### 일일 업무

```bash
# 1. 09:00 - 일일 로그 확인
docker-compose logs --since 1h blog | grep ERROR

# 2. 11:00 - 성능 메트릭 확인
curl http://localhost:9090/api/v1/query?query='avg(rate(http_request_duration_seconds_sum[5m]))'

# 3. 14:00 - 백업 상태 확인
ls -lh /backups/freelang-backup-*.sql.gz | tail -3

# 4. 17:00 - 보안 패치 확인
docker system prune -a --filter 'until=72h'
apt update && apt list --upgradable
```

### 주간 업무

```bash
# 월요일: 성능 분석
# - 평균 응답시간 추이
# - 에러율 분석
# - 용량 예측

# 수요일: 백업 테스트
bash scripts/restore.sh /backups/latest.sql.gz
# → 테스트 DB에서 검증
# → 복구 시간 측정

# 금요일: 보안 검토
# - 로그 분석 (401/403 에러)
# - 의존성 업데이트
# - 취약점 스캔
```

### 월간 업무

```bash
# 보안 패치 적용
docker-compose down
git pull origin main
docker-compose build --no-cache
docker-compose up -d

# 성능 최적화
# - 느린 쿼리 분석
# - 인덱스 최적화
# - 캐시 전략 검토

# 비용 분석
# - 대역폭 사용량
# - 스토리지 용량
# - 리소스 할당
```

---

## 🆘 장애 대응

### Scenario 1: 서비스 다운

```bash
# 1. 즉시 대응
docker-compose ps
docker-compose logs blog | tail -100

# 2. 자동 재시작 시도
docker-compose down
docker-compose up -d

# 3. 헬스 체크
curl http://localhost:5021/api/health

# 4. 데이터 무결성 확인
docker-compose exec postgres pg_dump -U freelang freelang | wc -l

# 5. 모니터링 재개
curl http://localhost:9090/api/v1/query?query=up
```

### Scenario 2: 메모리 부족

```bash
# 1. 현재 사용량 확인
free -h
docker stats

# 2. 캐시 정리
docker-compose exec redis redis-cli FLUSHDB

# 3. 컨테이너 재시작
docker-compose restart blog redis

# 4. 모니터 설정 조정
# - Redis maxmemory 증가
# - DB 연결 풀 크기 조정
```

### Scenario 3: 높은 에러율

```bash
# 1. 에러 로그 분석
docker-compose logs --since 30m blog | grep ERROR | head -20

# 2. 외부 서비스 확인
# - OAuth 제공자 상태
# - CDN 상태
# - 메일 서비스

# 3. 데이터베이스 상태 확인
docker-compose exec postgres \
  psql -U freelang -d freelang -c "SELECT pg_database.datname, \
    CAST(pg_database_size(pg_database.datname) AS bigint) \
    FROM pg_database WHERE datname = 'freelang';"

# 4. 연결 풀 상태
# PostgreSQL connection count
# Redis connected clients
```

---

## 📈 성능 최적화

### 주간 최적화 작업

```bash
# 1. 느린 쿼리 분석
docker-compose exec postgres \
  psql -U freelang -d freelang -c \
  "SELECT query, calls, mean_exec_time, max_exec_time FROM pg_stat_statements \
   ORDER BY mean_exec_time DESC LIMIT 10;"

# 2. 인덱스 추가
docker-compose exec postgres \
  psql -U freelang -d freelang -c \
  "CREATE INDEX idx_posts_author_created ON posts(author, created_at DESC);"

# 3. 캐시 효율성 확인
docker-compose exec redis redis-cli INFO stats | grep hits
# Hit ratio = Hits / (Hits + Misses)

# 4. 연결 풀 최적화
# 현재: max_connections = 100
# 추천: (RAM_GB * 2) 또는 CPU * 4
```

---

## ✅ 최종 체크리스트

배포 전:
- [ ] .env 파일 설정
- [ ] GitHub Secrets 설정
- [ ] SSL 인증서 발급
- [ ] 데이터베이스 백업
- [ ] 모니터링 알림 테스트

배포 중:
- [ ] 로컬 테스트 완료
- [ ] 서버 배포 완료
- [ ] 헬스 체크 확인
- [ ] 모니터링 시작

배포 후:
- [ ] 실시간 메트릭 모니터링 (1시간)
- [ ] 알림 규칙 검증 (1시간 후)
- [ ] 일일 체크 완료
- [ ] 주간 성능 분석 예약

---

## 📞 긴급 연락처

```
DevOps 담당자: [연락처]
DBA 담당자: [연락처]
보안 담당자: [연락처]

모니터링 대시보드:
- Prometheus: https://253.dclub.kr:9090
- Grafana: https://253.dclub.kr:3100 (admin/admin123)
- Jaeger: https://253.dclub.kr:16686

로그 수집:
- Docker logs: docker-compose logs -f
- Application logs: /var/log/freelang/app.log
- Nginx logs: /var/log/nginx/access.log
```

---

**배포 준비 상태**: ✅ 완료
**예상 배포 시간**: 2-3시간
**예상 다운타임**: <5분 (무중단 배포 권장)

**배포 담당자**: CloudOps Team
**배포 일시**: 2026-03-13 (협의 후 결정)

🚀 **프로덕션 배포 준비 완료**
