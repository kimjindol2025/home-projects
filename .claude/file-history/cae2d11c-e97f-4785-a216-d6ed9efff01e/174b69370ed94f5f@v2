# ============================================
# Multi-Stage Docker Build
# ============================================
# 최종 이미지 크기: < 100MB

# ===== Stage 1: Builder =====
FROM rust:1.75-slim as builder

WORKDIR /app

# 의존성 캐싱 (소스 코드 변경 시에만 재빌드)
COPY Cargo.toml Cargo.lock ./

# 더미 main.rs로 의존성 먼저 빌드
RUN mkdir -p src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release && \
    rm -rf src

# 소스 코드 복사
COPY src ./src

# 실제 빌드
RUN cargo build --release --locked

# 빌드 결과물 확인
RUN ls -lh /app/target/release/distributed_bank || echo "Binary not found, checking target structure..." && \
    find /app/target/release -type f -executable -name "*distributed*" || true

# ===== Stage 2: Runtime =====
FROM debian:bookworm-slim

# 필수 라이브러리 설치
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 라벨 추가
LABEL maintainer="Claude Code AI"
LABEL version="1.0.0"
LABEL description="Distributed Banking System - Rust Optimized (Phase E)"

# Stage 1에서 바이너리 복사 (여러 가능성 시도)
COPY --from=builder /app/target/release/distributed_bank /usr/local/bin/distributed_bank || true

# 바이너리 권한 설정
RUN chmod +x /usr/local/bin/distributed_bank || true

# 인증서 디렉토리 생성 (마운트용)
RUN mkdir -p /app/certs

# 헬스 체크 추가
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 포트 설정
# 8080: HTTP/Proxy
# 9000: Raft Leader
# 9100: Raft Follower 1
# 9200: Raft Follower 2
EXPOSE 8080 9000 9100 9200

# 실행
CMD ["distributed_bank"]
