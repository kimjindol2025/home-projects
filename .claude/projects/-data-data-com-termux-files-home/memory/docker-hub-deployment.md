---
name: Docker Hub 공용 배포 (Sovereign Workspace v1.0.0)
description: Docker Hub에서 freelang/sovereign-workspace를 공용 배포하기 위한 계획 및 진행
type: project
---

# Docker Hub 공용 배포 진행

**시작일**: 2026-03-18
**상태**: 🚀 배포 준비 완료 (Step 1-5 체크리스트 작성)
**목표**: freelang/sovereign-workspace를 Docker Hub 공용 레지스트리에 배포

## 완료 사항

### Phase 1-12 구현 완료
- 20,990줄 FV-Lang 코드
- 413개 테스트 (100% 통과)
- HTTP API + WebSocket + SQLite + gRPC
- 멀티테넌트 격리 (API 키 기반)

### 배포 파일 완성
1. **Dockerfile** (95줄)
   - Ubuntu 22.04 + FV-Lang 컴파일러
   - HTTP 8080 + gRPC 50051 포트
   - /data 볼륨 (SQLite)
   - healthcheck 30초 간격

2. **docker-compose.yml** (72줄)
   - sovereign-workspace 서비스
   - dashboard (Nginx) 서비스
   - sovereign-net 네트워크
   - 자동 재시작 정책

3. **DOCKER_HUB.md** (475줄)
   - 이미지 정보 (Tags: latest, v1.0.0, stable)
   - 3가지 빠른 시작 옵션
   - API 엔드포인트 (curl 예제 포함)
   - 환경 변수 & 포트 매핑
   - 보안 설정 (API 키, localhost, HTTPS)
   - 모니터링 & 트러블슈팅

4. **DEPLOYMENT_FVLANG.md** (450줄)
   - 단계별 설치 가이드 (6 steps)
   - FV-Lang 컴파일러 설치
   - docker-compose.fvlang.yml 사용법
   - 포트/성능 최적화
   - 트러블슈팅 (5가지 문제 해결)

5. **.dockerignore** (53줄)
   - Git, 빌드, IDE, 로그, 데이터, 테스트 제외
   - 이미지 크기 최적화

### 문서화
- **DOCKER_HUB_DEPLOYMENT.md** (200줄 이상)
  - 배포 체크리스트
  - 다음 단계 가이드 (Step 1-5)
  - Docker Hub 이미지 정보
  - 보안 고려사항
  - 예상 타임라인

## 다음 단계 (필수)

### Step 1: Docker Hub 계정 및 로그인
```bash
docker login
# 또는 토큰 사용
export DOCKER_TOKEN=<token>
docker login -u freelang --password-stdin <<< "$DOCKER_TOKEN"
```

### Step 2: 로컬 이미지 빌드
```bash
docker build -t freelang/sovereign-workspace:latest .
docker build -t freelang/sovereign-workspace:v1.0.0 .
docker build -t freelang/sovereign-workspace:stable .
```

### Step 3: Docker Hub에 푸시
```bash
docker push freelang/sovereign-workspace:latest
docker push freelang/sovereign-workspace:v1.0.0
docker push freelang/sovereign-workspace:stable
```

### Step 4: 공용 GOGS 저장소 생성
```bash
# https://gogs.dclub.kr/repo/create에서 새 저장소 생성
# 또는 API로 생성

git remote add public https://gogs.dclub.kr/freelang/sovereign-workspace.git
git push -u public master
```

### Step 5: 문서 작성
- [ ] Docker Hub README
- [ ] GOGS README
- [ ] LICENSE (MIT)
- [ ] CONTRIBUTING.md

## Demo 테넌트 정보

```
🔐 Demo Tenant 1: key_123456
🔐 Demo Tenant 2: key_789012

프로덕션: openssl rand -hex 32로 생성
```

## 기본 사용법

```bash
# 다운로드 & 실행
docker pull freelang/sovereign-workspace:latest
docker run -d \
  -p 8080:8080 \
  -p 50051:50051 \
  -v sovereign-data:/data \
  freelang/sovereign-workspace:latest

# 헬스 체크
curl -H "X-API-Key: key_123456" http://localhost:8080/health
```

## 브랜치 정보

### 현재 리포지토리
- **주소**: https://gogs.dclub.kr/kim/sovereign-workspace
- **브랜치**: master (Phase 1-12 완성)
- **커밋**: ce2f418 (Docker Hub 배포 구성)

### 새 공용 리포지토리 (예정)
- **주소**: https://gogs.dclub.kr/freelang/sovereign-workspace
- **상태**: 생성 대기
- **목적**: 커뮤니티 공용 배포

## 예상 이미지 크기

- **Base**: Ubuntu 22.04 (~77MB)
- **FV-Lang 컴파일러**: ~150MB
- **소스 코드**: ~100MB
- **의존성**: ~50MB
- **총합**: ~500MB

## 지원 플랫폼

- ✅ linux/amd64 (x86-64)
- ✅ linux/arm64 (ARM64, M1/M2/M3 Mac)
- ✅ WSL2 (Windows)
- ✅ Termux (Android - docker-compose.fvlang.yml)

## 라이선스

MIT License - 자유로운 사용, 수정, 배포

## 성공 기준

- [x] Dockerfile 작성 및 검증
- [x] docker-compose 작성
- [x] 배포 가이드 문서 완성
- [ ] Docker Hub 이미지 빌드 & 푸시
- [ ] 공용 GOGS 저장소 생성 & 코드 푸시
- [ ] Docker Hub README 작성
- [ ] 커뮤니티 공개 (Issues/Discussions 활성화)

