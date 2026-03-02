# 🔐 Phase E: TLS/암호화 + Docker 컨테이너화 계획

**작성일**: 2026-03-02
**목표 기한**: 1주일 (2026-03-09)
**상태**: 🟢 **시작**

---

## 📋 개요

FreeLang 분산 은행 시스템을 프로덕션 환경에 배포하기 위해 두 가지 핵심 작업 수행:

1. **TLS/암호화 추가** (보안)
2. **Docker 컨테이너화** (배포)

---

## 🔐 Part 1: TLS/암호화 추가

### 1.1 의존성 추가

**Cargo.toml 업데이트**:
```toml
# 기존 의존성
tokio = { version = "1.35", features = ["full"] }
hyper = { version = "1.0", features = ["full"] }

# 신규 추가
rustls = "0.21"                 # TLS 라이브러리
rustls-pemfile = "2.0"          # PEM 인증서 파싱
tokio-rustls = "0.24"           # tokio 통합
hyper-rustls = "0.24"           # hyper 통합

# 자체 서명 인증서 생성용 (선택)
rcgen = "0.11"                  # 테스트 인증서 생성
```

### 1.2 모듈 구조

```
src/
├── main.rs                      # 메인 (기존)
├── raft/mod.rs                  # Raft (기존)
├── proxy/mod.rs                 # Proxy (기존)
├── bank/mod.rs                  # Bank (기존)
└── security/mod.rs              # 신규: 보안
    ├── tls.rs                   # TLS 설정
    ├── auth.rs                  # 인증
    └── encryption.rs            # 암호화
```

### 1.3 구현 계획

**security/tls.rs** (300줄):
```rust
pub struct TlsConfig {
    pub cert_path: String,
    pub key_path: String,
    pub enabled: bool,
}

impl TlsConfig {
    pub fn new(cert: &str, key: &str) -> Self { ... }
    pub fn load_cert(&self) -> Result<Certificate> { ... }
    pub fn load_key(&self) -> Result<PrivateKey> { ... }
}

pub struct SecureProxy {
    inner: ProxyServer,
    tls: Arc<TlsConfig>,
}

impl SecureProxy {
    pub async fn handle_secure_request(&self, ...) -> Result<()> { ... }
}
```

**security/auth.rs** (250줄):
```rust
pub struct AuthToken {
    pub token: String,
    pub issued_at: u64,
    pub expires_at: u64,
}

pub struct TokenValidator {
    secret: String,
}

impl TokenValidator {
    pub fn validate(&self, token: &str) -> Result<AuthToken> { ... }
    pub fn generate(&self, user: &str) -> String { ... }
}
```

**security/encryption.rs** (200줄):
```rust
use chacha20poly1305::ChaCha20Poly1305;

pub struct DataEncryption {
    cipher: ChaCha20Poly1305,
}

impl DataEncryption {
    pub fn encrypt(&self, data: &[u8]) -> Result<Vec<u8>> { ... }
    pub fn decrypt(&self, data: &[u8]) -> Result<Vec<u8>> { ... }
}
```

### 1.4 테스트

**tests/security_tests.rs** (400줄):
```rust
#[tokio::test]
async fn test_tls_handshake() { ... }

#[tokio::test]
async fn test_token_generation() { ... }

#[tokio::test]
async fn test_encryption_decryption() { ... }

#[tokio::test]
async fn test_invalid_token_rejection() { ... }

#[tokio::test]
async fn test_secure_transfer() { ... }
```

---

## 🐳 Part 2: Docker 컨테이너화

### 2.1 Dockerfile (Multi-Stage Build)

**파일**: `Dockerfile`

```dockerfile
# ===== Stage 1: Builder =====
FROM rust:1.75-slim as builder

WORKDIR /app

# 의존성 캐싱
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release && \
    rm -rf src

# 소스 코드
COPY src ./src
RUN cargo build --release --locked

# ===== Stage 2: Runtime =====
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Stage 1에서 바이너리 복사
COPY --from=builder /app/target/release/distributed_bank /usr/local/bin/

# 포트 설정
EXPOSE 8080 9000 9100 9200

# 실행
CMD ["distributed_bank"]
```

**최적화 전략**:
- ✅ Multi-stage build (최종 크기 < 100MB)
- ✅ 의존성 캐싱 (빌드 시간 단축)
- ✅ 최소 런타임 이미지 (debian slim)
- ✅ 불필요한 파일 제거

### 2.2 docker-compose.yml

**파일**: `docker-compose.yml`

```yaml
version: '3.9'

services:
  # Raft 리더 (Node 0)
  raft-leader:
    build: .
    container_name: distributed_bank_leader
    ports:
      - "8080:8080"     # Proxy
      - "9000:9000"     # Raft
    environment:
      - RUST_LOG=info
      - NODE_ID=0
      - TLS_ENABLED=true
    volumes:
      - ./certs:/app/certs:ro
    networks:
      - bank_network
    healthcheck:
      test: ["CMD", "curl", "-f", "https://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  # Raft 팔로워 (Node 1-2)
  raft-follower-1:
    build: .
    container_name: distributed_bank_follower_1
    ports:
      - "9100:9100"
    environment:
      - RUST_LOG=info
      - NODE_ID=1
    networks:
      - bank_network

  raft-follower-2:
    build: .
    container_name: distributed_bank_follower_2
    ports:
      - "9200:9200"
    environment:
      - RUST_LOG=info
      - NODE_ID=2
    networks:
      - bank_network

networks:
  bank_network:
    driver: bridge
```

### 2.3 인증서 생성 스크립트

**파일**: `scripts/generate-certs.sh`

```bash
#!/bin/bash
# TLS 인증서 생성

set -e

CERT_DIR="./certs"
mkdir -p "$CERT_DIR"

# 개인키 생성
openssl genrsa -out "$CERT_DIR/server.key" 2048

# 인증서 서명 요청 (CSR)
openssl req -new \
    -key "$CERT_DIR/server.key" \
    -out "$CERT_DIR/server.csr" \
    -subj "/C=KR/ST=Seoul/L=Seoul/O=FreeLang/CN=localhost"

# 자체 서명 인증서 생성 (365일)
openssl x509 -req \
    -days 365 \
    -in "$CERT_DIR/server.csr" \
    -signkey "$CERT_DIR/server.key" \
    -out "$CERT_DIR/server.crt"

# 권한 설정
chmod 600 "$CERT_DIR/server.key"
chmod 644 "$CERT_DIR/server.crt"

echo "✅ TLS 인증서 생성 완료"
echo "  - 인증서: $CERT_DIR/server.crt"
echo "  - 개인키: $CERT_DIR/server.key"
```

---

## 📊 구현 로드맵

### Week 1 (2026-03-02 ~ 2026-03-09)

| 날짜 | 작업 | 상태 |
|------|------|------|
| 3/2 | Cargo.toml 업데이트, 계획 수립 | 🟢 진행중 |
| 3/3 | security/tls.rs 구현 (300줄) | ⏳ 예정 |
| 3/4 | security/auth.rs 구현 (250줄) | ⏳ 예정 |
| 3/5 | security/encryption.rs 구현 (200줄) | ⏳ 예정 |
| 3/6 | 보안 테스트 작성 (400줄) | ⏳ 예정 |
| 3/7 | Dockerfile + docker-compose 작성 | ⏳ 예정 |
| 3/8 | 인증서 생성 스크립트 + 문서 | ⏳ 예정 |
| 3/9 | 최종 검증 + GOGS 커밋 | ⏳ 예정 |

---

## 🎯 구현 세부사항

### Phase E-1: TLS/보안 모듈 (750줄)

**목표**: 모든 통신을 암호화

**파일 구성**:
- `src/security/mod.rs` (50줄): 모듈 정의
- `src/security/tls.rs` (300줄): TLS 설정
- `src/security/auth.rs` (250줄): 인증 토큰
- `src/security/encryption.rs` (150줄): 데이터 암호화
- `tests/security_tests.rs` (400줄): 5+ 보안 테스트

**기능**:
1. ✅ HTTPS 지원 (TLS 1.3)
2. ✅ JWT 토큰 인증
3. ✅ ChaCha20-Poly1305 암호화
4. ✅ 인증서 검증
5. ✅ 보안 로깅

### Phase E-2: Docker 컨테이너화 (500줄)

**파일 구성**:
- `Dockerfile` (50줄): Multi-stage build
- `docker-compose.yml` (120줄): 클러스터 설정
- `scripts/generate-certs.sh` (40줄): 인증서 생성
- `scripts/docker-build.sh` (30줄): 빌드 자동화
- `.dockerignore` (20줄): 무시 파일
- `docs/DOCKER_GUIDE.md` (240줄): 사용 가이드

**기능**:
1. ✅ Multi-stage 빌드 (<100MB)
2. ✅ Docker Compose (3노드)
3. ✅ Health check
4. ✅ 볼륨 마운팅
5. ✅ 환경변수 설정

---

## 📈 최종 성과 예상

### 코드 규모
- **TLS/보안**: 750줄 (모듈 + 테스트)
- **Docker**: 500줄 (설정 + 스크립트)
- **총 Phase E**: **1,250줄** 신규

### 전체 프로젝트 규모
- FreeLang: 9,191줄 (기존)
- Rust: 1,185줄 (기존)
- Phase E: 1,250줄 (신규)
- **총합**: **11,626줄**

### 기능 확대
- ✅ HTTPS/TLS 1.3
- ✅ JWT 인증
- ✅ 엔드-투-엔드 암호화
- ✅ Docker 배포
- ✅ 컨테이너 오케스트레이션
- ✅ Health monitoring

---

## 🔍 기술 상세

### TLS 구현 패턴

```rust
// 보안 프록시 생성
let tls_config = TlsConfig::new(
    "./certs/server.crt",
    "./certs/server.key"
)?;

let secure_proxy = SecureProxy::new(proxy, tls_config);

// HTTPS 수신
secure_proxy.listen("0.0.0.0:8443").await?;
```

### Docker 이미지 최적화

```
빌드 전: ~2GB (컴파일러 + 소스 + 의존성)
빌드 후: ~80MB (런타임만)
압축률: 96% 감소
```

### 배포 흐름

```
소스코드
    ↓
Docker build (Dockerfile)
    ↓
Multi-stage (Stage 1 제거)
    ↓
최종 이미지 (~80MB)
    ↓
docker-compose up
    ↓
3노드 클러스터 실행
    ↓
HTTPS + TLS 활성화
```

---

## ✅ 완료 기준

1. ✅ TLS/암호화 모듈 구현
2. ✅ 보안 테스트 (5+ 개)
3. ✅ Docker 이미지 생성
4. ✅ docker-compose 배포 테스트
5. ✅ 인증서 생성 자동화
6. ✅ 완전한 문서 작성
7. ✅ GOGS 커밋 및 저장

---

## 📚 산출물

1. **src/security/** (3 파일, 750줄)
   - TLS 설정
   - JWT 인증
   - 데이터 암호화

2. **Docker 파일** (5 파일, 250줄)
   - Dockerfile
   - docker-compose.yml
   - .dockerignore
   - 빌드 스크립트

3. **인증서 관리** (2 파일, 70줄)
   - 생성 스크립트
   - 갱신 가이드

4. **문서** (500줄)
   - Docker 사용 가이드
   - TLS 설정 가이드
   - 배포 체크리스트

---

**시작 날짜**: 2026-03-02
**예상 완료**: 2026-03-09
**예상 산출**: **1,250줄** + 완전한 문서

---

이제 구현을 시작합니다! 🚀
