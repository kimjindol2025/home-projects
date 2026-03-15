# 🚀 253.dclub.kr 배포 체크리스트

## 배포 전 필수 확인 (5분)

### 1. SSH 접속 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "echo '✅ SSH 연결 성공'"
```
- [ ] SSH 연결 성공

### 2. Docker 설치 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "docker --version && docker-compose --version"
```
- [ ] Docker 설치됨
- [ ] Docker Compose 설치됨

### 3. 디스크 공간 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "df -h"
```
- [ ] 최소 10GB 여유 공간 확보

### 4. 필수 포트 열림 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "sudo netstat -tlnp | grep -E ':(80|443|3001|5432|6379|9090|3100)'"
```
- [ ] Port 80 (HTTP) 열려있음
- [ ] Port 443 (HTTPS) 열려있음
- [ ] Port 3001 (API) 열려있음

---

## 배포 실행 (1-2시간)

### 단계 1: 배포 스크립트 실행
```bash
bash scripts/deploy-253.sh
```

이 스크립트가 자동으로 실행:
1. ✅ .env 파일 업로드
2. ✅ Let's Encrypt 인증서 발급
3. ✅ Docker 이미지 빌드
4. ✅ 서비스 시작 (docker-compose up -d)
5. ✅ 헬스 체크 실행
6. ✅ HTTPS 검증

### 단계 2: 실시간 로그 확인 (배포 스크립트 실행 중)
```bash
ssh -p 10053 kimjin@253.dclub.kr "cd ~/freelang-light && docker-compose logs -f"
```

---

## 배포 완료 검증 (15분)

### 1. API 헬스 체크
```bash
# HTTP (리다이렉트 확인)
curl -I http://253.dclub.kr/api/health

# HTTPS (최종 확인)
curl -I https://253.dcloud.kr/api/health
```
**예상**: HTTP 301 + HTTPS 200 OK

### 2. 데이터베이스 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr << 'EOF'
  docker-compose exec postgres psql -U freelang -d freelang -c "
    SELECT
      'posts' as table_name, COUNT(*) as count FROM posts
    UNION ALL
    SELECT 'users', COUNT(*) FROM users
    UNION ALL
    SELECT 'comments', COUNT(*) FROM comments;
  "
EOF
```
**예상**: 초기 테이블 생성 확인

### 3. SSL 인증서 확인
```bash
openssl s_client -connect 253.dclub.kr:443 -servername 253.dclub.kr < /dev/null 2>&1 | grep -A 2 "Issuer:"
```
**예상**: Let's Encrypt 인증서 발급자 표시

### 4. Prometheus 메트릭 수집
```bash
curl -s http://253.dclub.kr:9090/api/v1/query?query=up | jq '.'
```
**예상**: JSON 응답 + up{job="..."}

### 5. Grafana 접근
```bash
# 브라우저에서 접속
open http://253.dclub.kr:3100
# 기본 계정: admin / admin123
```

---

## 배포 후 관리

### 일일 점검 (매일 09:00)
```bash
ssh -p 10053 kimjin@253.dclub.kr << 'EOF'
  echo "=== 컨테이너 상태 ==="
  docker-compose ps

  echo ""
  echo "=== 최근 에러 로그 ==="
  docker-compose logs --since 24h | grep ERROR | head -10

  echo ""
  echo "=== 디스크 사용량 ==="
  df -h | grep -E '/$|/home'

  echo ""
  echo "=== 메모리 사용량 ==="
  free -h
EOF
```

### 주간 점검 (매주 월요일)
```bash
# 백업 상태 확인
ssh -p 10053 kimjin@253.dclub.kr "ls -lh ~/freelang-light/backups/ | tail -10"

# 인증서 만료일 확인
ssh -p 10053 kimjin@253.dclub.kr "sudo certbot certificates"

# 성능 메트릭 확인
curl -s http://253.dclub.kr:9090/api/v1/query?query='avg(rate(http_request_duration_seconds_sum[5m]))' | jq '.'
```

### 월간 점검 (매월 첫째 주 금요일)
```bash
# 보안 패치 확인 및 적용
ssh -p 10053 kimjin@253.dclub.kr << 'EOF'
  sudo apt update
  apt list --upgradable
  # sudo apt upgrade -y  (필요시)
EOF

# 백업 복구 테스트 (스테이징 DB에서만)
ssh -p 10053 kimjin@253.dclub.kr "bash ~/freelang-light/scripts/restore.sh /backups/latest.sql.gz"
```

---

## 문제 해결

### API 응답 없음
```bash
ssh -p 10053 kimjin@253.dclub.kr "docker-compose logs blog | tail -50"
```

### 데이터베이스 연결 실패
```bash
ssh -p 10053 kimjin@253.dclub.kr << 'EOF'
  docker-compose exec postgres psql -U freelang -d freelang -c "SELECT version();"
  docker-compose logs postgres | tail -20
EOF
```

### HTTPS 인증서 이슈
```bash
ssh -p 10053 kimjin@253.dclub.kr "sudo certbot renew --dry-run"
```

### 메모리 부족
```bash
ssh -p 10053 kimjin@253.dclub.kr << 'EOF'
  docker system prune -a
  docker-compose down
  docker-compose up -d
EOF
```

---

## 배포 완료 신호

✅ 모두 완료되면 다음을 확인:
- [ ] API HTTPS 응답 200 OK
- [ ] PostgreSQL 테이블 존재
- [ ] Prometheus 메트릭 수집 중
- [ ] Grafana 접근 가능
- [ ] SSL 인증서 유효기간 확인

**축하합니다! Phase 1 배포 완료! 🎉**
