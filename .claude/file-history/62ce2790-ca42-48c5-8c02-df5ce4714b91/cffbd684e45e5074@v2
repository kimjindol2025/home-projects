# 🚀 지금 바로 배포 시작

## 준비된 파일
```
✅ .env.production          - 프로덕션 환경변수
✅ scripts/deploy-253.sh    - 자동 배포 스크립트
✅ docker-compose.yml       - 서비스 정의
✅ nginx.conf               - HTTPS 역프록시
✅ postgres-init.sql        - 데이터베이스 스키마
```

---

## 🎯 배포 시작 (3줄만 실행)

### 1단계: 배포 스크립트 실행
```bash
cd /tmp/freelang-light
bash scripts/deploy-253.sh
```

**소요 시간**: 1-2시간
- Let's Encrypt 인증서: 5-10분
- Docker 빌드: 20-30분
- 서비스 시작: 5-10분
- 헬스 체크: 5분

### 2단계: 배포 중 실시간 로그 보기 (다른 터미널)
```bash
ssh -p 10053 kimjin@253.dclub.kr "cd ~/freelang-light && docker-compose logs -f"
```

### 3단계: 배포 완료 확인
```bash
curl -I https://253.dclub.kr/api/health
```

**예상 응답**:
```
HTTP/2 200
Strict-Transport-Security: max-age=31536000
```

---

## 📊 배포 후 확인할 것들

| 서비스 | URL | 예상값 |
|--------|-----|--------|
| API Health | https://253.dclub.kr/api/health | 200 OK |
| Prometheus | http://253.dclub.kr:9090 | Graph UI 표시 |
| Grafana | http://253.dclub.kr:3100 | 로그인 화면 (admin/admin123) |

---

## ⚠️ 배포 전 필수 사항

✅ SSH 접속 확인:
```bash
ssh -p 10053 kimjin@253.dclub.kr "uname -a"
```

✅ Docker 설치 확인:
```bash
ssh -p 10053 kimjin@253.dclub.kr "docker --version"
```

✅ 디스크 공간 확인 (최소 10GB):
```bash
ssh -p 10053 kimjin@253.dclub.kr "df -h"
```

---

## 🎬 전체 배포 흐름

```
$ bash scripts/deploy-253.sh

🚀 FreeLang Light 프로덕션 배포 시작 (253.dclub.kr)
==================================================

📋 [Step 1/5] .env 파일 업로드
----------------------------------------
✅ .env 파일 확인:
SERVER_HOST=0.0.0.0
SERVER_PORT=5021
...

🔒 [Step 2/5] Let's Encrypt 인증서 발급
----------------------------------------
⏳ Let's Encrypt 인증서 발급...
✅ 인증서 확인:
/etc/letsencrypt/live/253.dclub.kr/
├── cert.pem
├── chain.pem
├── fullchain.pem
└── privkey.pem

🐳 [Step 3/5] Docker Compose 빌드 & 시작
----------------------------------------
⏳ Docker 이미지 빌드 중 (2-3분 소요)...
[+] Building 45.2s
✅ 컨테이너 상태:
NAME                         STATUS
freelang-hybrid-api          Up 2s
freelang-hybrid-nginx        Up 1s
freelang-hybrid-db           Up 3s
freelang-hybrid-cache        Up 2s
freelang-hybrid-prometheus   Up 2s
freelang-hybrid-grafana      Up 1s

🏥 [Step 4/5] 서비스 헬스 체크
----------------------------------------
API 헬스 체크:
{"status": "ok", "uptime": 2.5}
✅ PostgreSQL 연결 확인:
1
✅ Redis 상태 확인:
PONG

🔐 [Step 5/5] HTTPS 검증
----------------------------------------
🔒 SSL 인증서 확인:
Issuer: C=US, O=Let's Encrypt, CN=Let's Encrypt Authority X3
✅ HTTPS 연결 테스트:
HTTP/2 200
Strict-Transport-Security: max-age=31536000

==================================================
✅ Phase 1 배포 완료!
==================================================

📊 접근 가능한 서비스:
  • API:          https://253.dclub.kr/api/health
  • Prometheus:   http://253.dclub.kr:9090
  • Grafana:      http://253.dclub.kr:3100 (admin/admin123)

📝 다음 단계:
  1. 로그 확인: docker-compose logs -f blog
  2. 백업 테스트: bash scripts/backup.sh
  3. 모니터링: https://253.dclub.kr:9090/graph
```

---

## 🔙 배포 취소/복구

만약 배포를 중단하려면:
```bash
ssh -p 10053 kimjin@253.dclub.kr "cd ~/freelang-light && docker-compose down"
```

이전 버전으로 복구:
```bash
ssh -p 10053 kimjin@253.dclub.kr "cd ~/freelang-light && docker-compose down && git checkout main && docker-compose up -d"
```

---

## 📞 배포 중 문제 발생 시

### 실시간 로그 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "docker-compose logs -f api"
```

### PostgreSQL 연결 확인
```bash
ssh -p 10053 kimjin@253.dclub.kr "docker-compose exec postgres psql -U freelang -d freelang -c 'SELECT COUNT(*) FROM posts;'"
```

### 인증서 발급 실패 시
```bash
ssh -p 10053 kimjin@253.dclub.kr "sudo certbot certonly -d 253.dclub.kr --email noreply@freelang.kr --agree-tos"
```

---

**준비 완료! 배포를 시작하세요 🚀**

```bash
bash scripts/deploy-253.sh
```
