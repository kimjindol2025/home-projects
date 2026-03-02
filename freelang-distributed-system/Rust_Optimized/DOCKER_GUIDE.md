# 🐳 Docker 배포 가이드

**Phase E: TLS/보안 + Docker 컨테이너화**

---

## 📚 목차

1. [빠른 시작](#빠른-시작)
2. [사전 요구사항](#사전-요구사항)
3. [인증서 설정](#인증서-설정)
4. [Docker 빌드](#docker-빌드)
5. [docker-compose 사용](#docker-compose-사용)
6. [모니터링](#모니터링)
7. [문제 해결](#문제-해결)
8. [프로덕션 체크리스트](#프로덕션-체크리스트)

---

## 빠른 시작

### 3단계로 실행하기

```bash
# 1. 인증서 생성
bash scripts/generate-certs.sh

# 2. Docker 이미지 빌드
docker-compose build

# 3. 클러스터 실행
docker-compose up -d
```

### 상태 확인

```bash
# 로그 보기
docker-compose logs -f

# 컨테이너 상태 확인
docker-compose ps

# 특정 서비스 로그
docker-compose logs raft-leader
```

### 종료

```bash
# 컨테이너 중지
docker-compose down

# 데이터 포함 전체 삭제
docker-compose down -v
```

---

## 사전 요구사항

### 필수 소프트웨어

| 도구 | 버전 | 용도 |
|------|------|------|
| Docker | 20.10+ | 컨테이너 실행 |
| Docker Compose | 2.0+ | 다중 컨테이너 관리 |
| OpenSSL | 1.1+ | TLS 인증서 생성 |
| Rust | 1.75+ | (빌드 시에만) |

### 설치 확인

```bash
# Docker 확인
docker --version
docker run hello-world

# Docker Compose 확인
docker-compose --version

# OpenSSL 확인
openssl version
```

### 시스템 요구사항

- **디스크**: 최소 5GB (이미지 + 데이터)
- **메모리**: 최소 2GB RAM
- **CPU**: 최소 2 코어
- **포트**: 8080, 9000, 9100, 9200 사용 가능

---

## 인증서 설정

### 자체 서명 인증서 생성

```bash
# 기본값 (365일 유효)
bash scripts/generate-certs.sh

# 커스텀 유효 기간 (730일)
bash scripts/generate-certs.sh ./certs 730
```

### 생성되는 파일

```
certs/
├── server.crt              # TLS 인증서
├── server.key              # 개인키
└── README                  # 생성 정보
```

### 인증서 검증

```bash
# 인증서 정보 확인
openssl x509 -in certs/server.crt -text -noout

# 개인키 확인
openssl rsa -in certs/server.key -text -noout

# 만료일 확인
openssl x509 -enddate -nodate -in certs/server.crt

# 인증서와 개인키 일치 확인
openssl x509 -noout -modulus -in certs/server.crt | openssl md5
openssl rsa -noout -modulus -in certs/server.key | openssl md5
```

### 프로덕션 인증서

```bash
# Let's Encrypt 인증서 사용 (추천)
sudo certbot certonly --standalone -d your-domain.com

# 생성된 인증서 복사
cp /etc/letsencrypt/live/your-domain.com/fullchain.pem certs/server.crt
cp /etc/letsencrypt/live/your-domain.com/privkey.pem certs/server.key
```

---

## Docker 빌드

### 이미지 빌드

```bash
# docker-compose를 사용한 빌드
docker-compose build

# 캐시 무시하고 다시 빌드
docker-compose build --no-cache

# 특정 서비스만 빌드
docker-compose build raft-leader
```

### 빌드 출력 확인

```bash
# 빌드 상세 로그
docker-compose build --verbose

# 최종 이미지 크기 확인
docker images | grep distributed_bank
```

### 예상 이미지 크기

| 단계 | 크기 |
|------|------|
| Builder 이미지 | ~2.5GB |
| Runtime 이미지 | ~80-100MB |
| 최종 압축 | ~30-40MB |

---

## docker-compose 사용

### 서비스 구조

```
┌─────────────────────────────────────┐
│     docker-compose (Orchestrator)   │
├─────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  │
│  │   Leader    │  │ Follower-1  │  │
│  │  (Node 0)   │  │ (Node 1)    │  │
│  │   :8080     │  │  :9100      │  │
│  │   :9000     │  │             │  │
│  └─────────────┘  └─────────────┘  │
│  ┌─────────────┐                    │
│  │ Follower-2  │                    │
│  │  (Node 2)   │                    │
│  │   :9200     │                    │
│  └─────────────┘                    │
└─────────────────────────────────────┘
       bank_network (172.20.0.0/16)
```

### 기본 명령어

```bash
# 모든 서비스 시작
docker-compose up -d

# 특정 서비스만 시작
docker-compose up -d raft-leader

# 모든 서비스 로그 보기
docker-compose logs -f

# 특정 서비스 로그 (최근 100줄)
docker-compose logs -f --tail=100 raft-leader

# 리소스 사용량 확인
docker-compose stats

# 실행 중인 컨테이너 확인
docker-compose ps

# 컨테이너 상태 상세 확인
docker-compose ps -a

# 서비스 재시작
docker-compose restart raft-leader

# 모든 서비스 중지
docker-compose stop

# 모든 서비스 삭제
docker-compose down

# 데이터 포함 전체 삭제
docker-compose down -v
```

### 환경변수 설정

`docker-compose.yml`에서 환경변수 수정:

```yaml
environment:
  - RUST_LOG=info          # 로그 레벨: trace, debug, info, warn, error
  - NODE_ID=0              # 노드 ID
  - NODE_TYPE=leader       # 노드 타입: leader, follower
  - TLS_ENABLED=false      # TLS 활성화
  - ENCRYPTION_ENABLED=true # 암호화 활성화
```

### 볼륨 설정

```yaml
volumes:
  - ./certs:/app/certs:ro           # 인증서 (읽기 전용)
  - bank_data_0:/app/data           # 영속 저장소
```

---

## 모니터링

### 헬스 체크

```bash
# 리더 노드 헬스 체크
curl http://localhost:8080/health

# 특정 포드 상태
docker-compose exec raft-leader sh -c 'curl http://localhost:8080/health'
```

### 로그 모니터링

```bash
# 실시간 로그 (모든 서비스)
docker-compose logs -f

# 리더 노드 로그 필터링
docker-compose logs -f raft-leader | grep "Transfer\|✅"

# 에러만 표시
docker-compose logs -f | grep "ERROR\|❌"
```

### 메트릭 수집

```bash
# 컨테이너 CPU/메모리 사용량
docker stats

# 특정 컨테이너 모니터링
docker stats distributed_bank_leader

# 디스크 사용량
docker volume ls
docker volume inspect bank_data_0
```

### 성능 벤치마크

```bash
# 리더에 요청 100개 전송 (성능 테스트)
for i in {1..100}; do
  curl http://localhost:8080/transfer &
done
wait

# 응답 시간 측정
time curl http://localhost:8080/transfer
```

---

## 문제 해결

### 빌드 실패

```bash
# 캐시 초기화
docker-compose down -v
docker system prune -a
docker-compose build --no-cache

# 빌드 로그 확인
docker-compose build --verbose 2>&1 | tail -100
```

### 포트 충돌

```bash
# 사용 중인 포트 확인
lsof -i :8080
netstat -tlnp | grep 8080

# docker-compose.yml에서 포트 변경
# ports:
#   - "8081:8080"  # 외부 포트 변경
```

### 네트워크 문제

```bash
# 네트워크 상태 확인
docker network ls
docker network inspect bank_network

# 컨테이너 간 통신 테스트
docker-compose exec raft-leader ping raft-follower-1

# DNS 확인
docker-compose exec raft-leader nslookup raft-leader
```

### 데이터 문제

```bash
# 데이터 볼륨 확인
docker volume ls
docker volume inspect bank_data_0

# 데이터 초기화
docker-compose down -v  # 주의: 모든 데이터 삭제

# 백업
docker cp distributed_bank_leader:/app/data ./backup-data
```

### 권한 문제

```bash
# Docker 소켓 권한
sudo usermod -aG docker $USER
newgrp docker

# 파일 권한
chmod 600 certs/server.key
chmod 644 certs/server.crt
```

---

## 프로덕션 체크리스트

### 배포 전 확인사항

- [ ] **보안**
  - [ ] 프로덕션 인증서 설치 (Let's Encrypt)
  - [ ] TLS 1.3 확인 (--tlsv1.3)
  - [ ] JWT 시크릿 변경 (강력한 비밀번호)
  - [ ] 방화벽 설정 (필요 포트만 개방)
  - [ ] Docker 업데이트 설치

- [ ] **성능**
  - [ ] 메모리 제한 설정 (예: 512MB)
  - [ ] CPU 제한 설정 (예: 1 코어)
  - [ ] 로그 로테이션 설정
  - [ ] 디스크 여유 공간 확인 (최소 10GB)

- [ ] **모니터링**
  - [ ] 헬스 체크 활성화
  - [ ] 로깅 시스템 구성
  - [ ] 알림 설정 (예: Slack)
  - [ ] 메트릭 수집 (Prometheus)

- [ ] **백업 & 복구**
  - [ ] 데이터 백업 계획 수립
  - [ ] 자동 백업 스크립트 작성
  - [ ] 복구 테스트 실행
  - [ ] 재해 복구 계획 (RTO/RPO)

- [ ] **문서화**
  - [ ] 배포 절차 문서화
  - [ ] 운영 매뉴얼 작성
  - [ ] 트러블슈팅 가이드
  - [ ] 연락처 정보 기록

### 프로덕션 docker-compose.yml

```yaml
version: '3.9'

services:
  raft-leader:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: distributed_bank_leader
    ports:
      - "8080:8080"
      - "9000:9000"
    environment:
      - RUST_LOG=warn          # 프로덕션은 warn 이상
      - NODE_ID=0
      - TLS_ENABLED=true       # HTTPS 필수
      - ENCRYPTION_ENABLED=true
    volumes:
      - ./certs:/app/certs:ro
      - bank_data_0:/app/data
    restart: always            # 자동 재시작
    healthcheck:
      test: ["CMD", "curl", "-f", "https://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '1'
        reservations:
          cpus: '0.5'
          memory: 512M
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "10"
```

---

## 고급 설정

### 커스텀 네트워크

```bash
# 기존 네트워크 사용
docker network create bank_network_custom
# docker-compose.yml에서:
# networks:
#   bank_network:
#     external: true
```

### 레지스트리에 푸시

```bash
# Docker Hub에 푸시
docker tag distributed_bank:latest yourusername/distributed_bank:1.0.0
docker push yourusername/distributed_bank:1.0.0

# 프라이빗 레지스트리
docker tag distributed_bank:latest registry.mycompany.com/banking/system:1.0.0
docker push registry.mycompany.com/banking/system:1.0.0
```

### Kubernetes 배포

```bash
# Docker 이미지에서 Kubernetes manifest 생성
docker-compose config > k8s-manifest.yaml

# 또는 kompose 사용
kompose convert -f docker-compose.yml -o k8s-deployment.yaml
kubectl apply -f k8s-deployment.yaml
```

---

## 성능 최적화

### 빌드 최적화

```bash
# 의존성 캐싱 활용
docker-compose build --no-cache=false

# 병렬 빌드
docker-compose build --parallel
```

### 런타임 최적화

```yaml
# docker-compose.yml
services:
  raft-leader:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

### 이미지 최적화

```dockerfile
# Dockerfile
FROM debian:bookworm-slim  # 이미 최적화됨

# 레이어 최소화
RUN apt-get update && \
    apt-get install -y package1 package2 && \
    rm -rf /var/lib/apt/lists/*
```

---

## 참고 자료

- [Docker 공식 문서](https://docs.docker.com/)
- [Docker Compose 참고](https://docs.docker.com/compose/compose-file/)
- [Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Security](https://docs.docker.com/engine/security/)

---

**작성자**: Claude Code AI
**작성일**: 2026-03-02
**버전**: 1.0.0
