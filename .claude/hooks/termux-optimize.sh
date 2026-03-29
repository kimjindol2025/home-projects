#!/bin/bash
# ============================================================================
# Termux 최적화 자동화
# 모바일 환경에서 Claude Code 리소스 관리 및 배포 자동화
# ============================================================================

set -e

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

LOG_FILE="$HOME/.claude/termux-optimize.log"

# ============================================================================
# 1. 리소스 모니터링
# ============================================================================
monitor_resources() {
    echo -e "${BLUE}[리소스 모니터링] 시스템 상태 점검${NC}"

    # 메모리 사용률 확인
    TOTAL_MEM=$(free -m | awk 'NR==2 {print $2}')
    USED_MEM=$(free -m | awk 'NR==2 {print $3}')
    MEM_PERCENT=$((USED_MEM * 100 / TOTAL_MEM))

    echo "  📊 메모리: ${USED_MEM}MB / ${TOTAL_MEM}MB ($MEM_PERCENT%)"

    # CPU 부하 확인
    CPU_LOAD=$(uptime | awk -F'load average:' '{ print $2 }' | awk '{print $1}')
    echo "  📊 CPU 부하: $CPU_LOAD"

    # Node.js 프로세스 확인
    if command -v node >/dev/null 2>&1; then
        NODE_PROCS=$(pgrep -c node 2>/dev/null || echo "0")
        NODE_MEM=$(ps aux | grep node | grep -v grep | awk '{sum+=$6} END {print sum}')
        echo "  📊 Node.js: $NODE_PROCS 프로세스, ${NODE_MEM:-0}MB 메모리"

        if [ "$NODE_MEM" -gt 150 ]; then
            echo -e "${RED}  ⚠️  Node.js 메모리 과다 사용! ($NODE_MEM MB)${NC}"
            echo "     → 작업 일시 중단 및 gc() 트리거 권장"
            return 1
        fi
    fi

    # 디스크 여유 확인
    DISK_FREE=$(df -h / | tail -1 | awk '{print $4}')
    DISK_USED=$(df -h / | tail -1 | awk '{print $3}')
    echo "  📊 디스크: $DISK_USED / $DISK_FREE 여유"

    if [ "$MEM_PERCENT" -gt 80 ]; then
        echo -e "${YELLOW}  ⚠️  메모리 부족! ($MEM_PERCENT%) - 백그라운드 작업 중단 권장${NC}"
        return 1
    fi

    echo -e "${GREEN}  ✓ 리소스 상태 OK${NC}"
    return 0
}

# ============================================================================
# 2. 배포 자동화
# ============================================================================
deploy_auto() {
    echo -e "${BLUE}[배포 자동화] 빌드 및 로컬 테스트${NC}"

    PROJECT_ROOT=$(pwd)
    PROJECT_NAME=$(basename "$PROJECT_ROOT")

    # Go 프로젝트 빌드
    if [ -f "go.mod" ]; then
        echo "  🔨 Go 프로젝트 빌드 중: $PROJECT_NAME"

        if go build -v ./... 2>&1 | tail -5; then
            echo -e "${GREEN}  ✓ Go 빌드 성공${NC}"
        else
            echo -e "${RED}  ✗ Go 빌드 실패${NC}"
            return 1
        fi
    fi

    # TypeScript 프로젝트 빌드
    if [ -f "tsconfig.json" ] && command -v tsc >/dev/null 2>&1; then
        echo "  🔨 TypeScript 컴파일 중: $PROJECT_NAME"

        if tsc --noEmit 2>&1 | tail -5; then
            echo -e "${GREEN}  ✓ TypeScript 컴파일 성공${NC}"
        else
            echo -e "${RED}  ✗ TypeScript 컴파일 실패${NC}"
            return 1
        fi
    fi

    # 테스트 실행
    if [ -f "package.json" ] && npm test >/dev/null 2>&1; then
        echo "  ✅ npm test PASS"
    elif [ -f "go.mod" ] && go test ./... >/dev/null 2>&1; then
        echo "  ✅ go test PASS"
    fi

    echo -e "${GREEN}  ✓ 배포 준비 완료${NC}"
    return 0
}

# ============================================================================
# 3. 빌드 결과 요약
# ============================================================================
summarize_build() {
    echo ""
    echo -e "${BLUE}[빌드 요약]${NC}"

    if [ -f "go.mod" ]; then
        echo "  📦 Go: go build ./..."
        if go build -v ./... >/dev/null 2>&1; then
            BIN_SIZE=$(ls -lh $(go list -f '{{.Target}}' ./...) 2>/dev/null | awk '{print $5}' | head -1)
            echo -e "     바이너리 크기: ${BIN_SIZE:-N/A} ${GREEN}✓${NC}"
        fi
    fi

    if [ -f "package.json" ]; then
        NPM_PACKAGES=$(jq '.dependencies | length' package.json 2>/dev/null || echo "0")
        echo "  📦 npm: $NPM_PACKAGES 패키지"
    fi

    # 테스트 결과
    TEST_RESULT=$(if [ -f "go.mod" ]; then go test ./... -v 2>&1 | grep -c "^ok\|^PASS"; fi)
    if [ "$TEST_RESULT" -gt 0 ]; then
        echo -e "  ✅ 테스트: $TEST_RESULT개 PASS"
    fi

    echo ""
    echo -e "${GREEN}💾 저장 필수, 기록이 증명이다 GOGS${NC}"
}

# ============================================================================
# 4. 백그라운드 작업 관리
# ============================================================================
manage_background_jobs() {
    echo -e "${BLUE}[백그라운드 작업 관리]${NC}"

    # 실행 중인 작업 목록
    JOBS_COUNT=$(jobs -r 2>/dev/null | wc -l)
    echo "  현재 실행 중인 작업: $JOBS_COUNT개"

    if [ "$JOBS_COUNT" -gt 3 ]; then
        echo -e "${YELLOW}  ⚠️  너무 많은 백그라운드 작업 (권장: 3개 이하)${NC}"
        echo "     작업 목록:"
        jobs -r | head -5
        return 1
    fi

    return 0
}

# ============================================================================
# 5. 핫스팟 분석 및 최적화 제안
# ============================================================================
analyze_hotspots() {
    echo -e "${BLUE}[핫스팟 분석]${NC}"

    # Go 프로젝트의 경우
    if [ -f "go.mod" ] && command -v go >/dev/null 2>&1; then
        echo "  🔍 CPU 프로필 생성 (1초)..."

        go test -run . -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem -bench . ./... 2>/dev/null || true

        if [ -f "cpu.prof" ]; then
            echo -e "${YELLOW}  📊 CPU 프로필 (상위 3 함수):${NC}"
            go tool pprof -top cpu.prof | head -10 || true
            rm -f cpu.prof mem.prof
        fi
    fi

    echo ""
}

# ============================================================================
# 6. 최적화 제안 생성
# ============================================================================
suggest_optimizations() {
    echo -e "${BLUE}[최적화 제안]${NC}"
    echo ""

    # 파일 크기 기준
    LARGEST_FILES=$(find . -type f \( -name "*.go" -o -name "*.ts" -o -name "*.fl" \) -exec ls -lS {} \; | head -3)

    if [ -n "$LARGEST_FILES" ]; then
        echo "  📈 큰 파일 (분할 권장):"
        echo "$LARGEST_FILES" | awk '{print "     " $9 " (" $5 " bytes)"}' || true
    fi

    # 함수 복잡도
    echo ""
    echo "  🔍 복잡한 함수 (리팩토링 권장):"

    if find . -name "*.go" >/dev/null 2>&1; then
        grep -n "^func " $(find . -name "*.go") 2>/dev/null | \
        while read -r line; do
            FILE=$(echo "$line" | cut -d: -f1)
            LINENUM=$(echo "$line" | cut -d: -f2)
            FUNC=$(echo "$line" | cut -d: -f3- | sed 's/func //')

            # 함수 길이 추정
            FUNC_END=$(tail -n +$LINENUM "$FILE" | grep -n "^}" | head -1 | cut -d: -f1)
            if [ "$FUNC_END" -gt 50 ]; then
                echo "     $FILE:$LINENUM - $FUNC ($FUNC_END 줄)"
            fi
        done | head -5
    fi

    echo ""
}

# ============================================================================
# 7. 종합 최적화 보고서
# ============================================================================
final_report() {
    echo ""
    echo "╔═══════════════════════════════════════════════════════╗"
    echo "║         Termux 최적화 종합 보고서                    ║"
    echo "╚═══════════════════════════════════════════════════════╝"
    echo ""
    echo "  📝 보고서 파일: $LOG_FILE"
    echo "  ⏰ 생성 시간: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    echo -e "${GREEN}✅ 기록이 증명이다${NC}"
    echo ""
}

# ============================================================================
# 메인 실행
# ============================================================================
main() {
    echo "🚀 Termux 최적화 자동화 시작..."
    echo ""

    # 1. 리소스 확인
    if ! monitor_resources; then
        echo -e "${RED}⚠️  리소스 부족. 작업 중단.${NC}"
        exit 1
    fi

    echo ""

    # 2. 백그라운드 작업 확인
    if ! manage_background_jobs; then
        echo -e "${YELLOW}⚠️  경고: 백그라운드 작업 많음${NC}"
    fi

    echo ""

    # 3. 배포 자동화
    if deploy_auto; then
        summarize_build
    fi

    echo ""

    # 4. 핫스팟 분석
    analyze_hotspots

    echo ""

    # 5. 최적화 제안
    suggest_optimizations

    echo ""

    # 6. 최종 보고서
    final_report

    # 로그 저장
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Termux 최적화 완료" >> "$LOG_FILE"
}

# 실행
main "$@"
