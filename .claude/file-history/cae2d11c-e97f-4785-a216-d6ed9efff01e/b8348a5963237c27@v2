#!/bin/bash
# ============================================
# TLS 인증서 생성 스크립트
# ============================================
# 자체 서명 인증서 생성 (테스트/개발용)

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 설정
CERT_DIR="${1:-.}/certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"
CSR_FILE="$CERT_DIR/server.csr"
DAYS="${2:-365}"

echo -e "${GREEN}🔐 TLS 인증서 생성${NC}"
echo "=================================="

# 디렉토리 생성
mkdir -p "$CERT_DIR"
echo -e "${YELLOW}✓${NC} 디렉토리 생성: $CERT_DIR"

# 이미 존재하는 파일 확인
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo -e "${YELLOW}⚠${NC}  기존 인증서가 있습니다."
    read -p "덮어쓰시겠습니까? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}작업 취소됨${NC}"
        exit 0
    fi
    rm -f "$CERT_FILE" "$KEY_FILE" "$CSR_FILE"
fi

# 1. 개인키 생성 (2048-bit RSA)
echo -e "${YELLOW}➊${NC} 개인키 생성 중... (2048-bit RSA)"
openssl genrsa -out "$KEY_FILE" 2048 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} 개인키 생성 완료"
else
    echo -e "${RED}✗${NC} 개인키 생성 실패"
    exit 1
fi

# 2. 인증서 서명 요청 (CSR) 생성
echo -e "${YELLOW}➋${NC} CSR (Certificate Signing Request) 생성 중..."
openssl req -new \
    -key "$KEY_FILE" \
    -out "$CSR_FILE" \
    -subj "/C=KR/ST=Seoul/L=Seoul/O=FreeLang/CN=localhost" \
    2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} CSR 생성 완료"
else
    echo -e "${RED}✗${NC} CSR 생성 실패"
    exit 1
fi

# 3. 자체 서명 인증서 생성
echo -e "${YELLOW}➌${NC} 자체 서명 인증서 생성 중... (${DAYS}일 유효)"
openssl x509 -req \
    -days "$DAYS" \
    -in "$CSR_FILE" \
    -signkey "$KEY_FILE" \
    -out "$CERT_FILE" \
    2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} 인증서 생성 완료"
else
    echo -e "${RED}✗${NC} 인증서 생성 실패"
    exit 1
fi

# 4. 권한 설정
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"
echo -e "${GREEN}✓${NC} 파일 권한 설정 완료"

# 5. CSR 파일 정리 (불필요)
rm -f "$CSR_FILE"

# 인증서 정보 표시
echo ""
echo -e "${GREEN}=================================="
echo "✅ TLS 인증서 생성 완료${NC}"
echo "=================================="
echo ""
echo "📄 인증서 정보:"
echo "   경로: $(cd "$CERT_DIR" && pwd)"
echo "   인증서: $(basename "$CERT_FILE")"
echo "   개인키: $(basename "$KEY_FILE")"
echo ""

# 인증서 상세 정보
echo "🔍 인증서 상세 정보:"
openssl x509 -in "$CERT_FILE" -text -noout 2>/dev/null | grep -E "Issuer:|Subject:|Not Before|Not After" | \
    sed 's/^/   /'

echo ""
echo -e "${YELLOW}⚠  주의사항:${NC}"
echo "   • 이 인증서는 테스트/개발용입니다"
echo "   • 프로덕션 환경에서는 CA에서 발급한 인증서를 사용하세요"
echo "   • HTTPS 접속 시 보안 경고가 나타날 수 있습니다"
echo ""

# 인증서 갱신 안내
EXPIRY_DATE=$(openssl x509 -enddate -nodate -in "$CERT_FILE" 2>/dev/null)
echo -e "${YELLOW}📅 갱신 일정:${NC}"
echo "   인증서 만료일: $EXPIRY_DATE"
echo "   갱신 명령: $0 \"$CERT_DIR\" <새_유효일>"
echo ""

echo -e "${GREEN}설정 완료! Docker 실행 준비가 되었습니다.${NC}"
echo "   \$ docker-compose up -d"
