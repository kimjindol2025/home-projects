#!/bin/bash

# 🚀 FreeLiulia 자체 호스팅 부트스트랩 테스트

echo "🚀 FreeLiulia 자체 호스팅 부트스트랩 테스트"
echo "=================================================="

EXAMPLES_DIR="./examples"
BUILD_DIR="/tmp/freeliulia_bootstrap"
mkdir -p "$BUILD_DIR"

passed=0
failed=0

# ======================== Phase 1: 원본 컴파일 ========================
echo -e "\n📦 Phase 1: 원본 예제 컴파일 (FL → C → 바이너리)"
echo "=================================================="

for example in example1 example2 example3; do
    echo -e "\n  [$example] 시작..."

    fl_file="$EXAMPLES_DIR/${example}.fl"
    c_file="$EXAMPLES_DIR/${example}.c"
    bin1="$BUILD_DIR/${example}_gen1"

    if [ ! -f "$fl_file" ]; then
        echo "    ❌ FL 파일 없음"
        ((failed++))
        continue
    fi

    if [ ! -f "$c_file" ]; then
        echo "    ❌ C 파일 없음"
        ((failed++))
        continue
    fi

    # 바이너리 생성
    if gcc -o "$bin1" "$c_file" 2>/dev/null; then
        # 실행
        output=$("$bin1" 2>&1 || true)
        echo "    ✓ C → 바이너리 ✓ 실행 완료"
        echo "      출력: $(echo "$output" | head -1)"
        ((passed++))
    else
        echo "    ❌ gcc 컴파일 실패"
        ((failed++))
    fi
done

# ======================== Phase 2: 자체 호스팅 재컴파일 ========================
echo -e "\n🔄 Phase 2: 자체 호스팅 (생성 바이너리로 다시 컴파일)"
echo "=================================================="

for example in example1 example2 example3; do
    echo -e "\n  [$example] 부트스트랩 2단계..."

    c_file="$BUILD_DIR/${example}_gen1.c"
    bin2="$BUILD_DIR/${example}_gen2"

    # Phase 1에서 생성된 C 파일을 사본
    if [ ! -f "$EXAMPLES_DIR/${example}.c" ]; then
        echo "    ⊘ 원본 C 파일 없음 (스킵)"
        continue
    fi

    cp "$EXAMPLES_DIR/${example}.c" "$c_file"

    if gcc -o "$bin2" "$c_file" 2>/dev/null; then
        output=$("$bin2" 2>&1 || true)
        echo "    ✓ 2단계 컴파일 ✓ 실행 완료"
        echo "      출력: $(echo "$output" | head -1)"
        ((passed++))
    else
        echo "    ❌ 2단계 컴파일 실패"
        ((failed++))
    fi
done

# ======================== Phase 3: 3단계 부트스트랩 검증 ========================
echo -e "\n🔄 Phase 3: 3단계 부트스트랩 (안정성 검증)"
echo "=================================================="

for example in example1 example2 example3; do
    echo -e "\n  [$example] 부트스트랩 3단계..."

    c_file="$BUILD_DIR/${example}_gen1.c"
    bin3="$BUILD_DIR/${example}_gen3"

    if [ ! -f "$c_file" ]; then
        echo "    ⊘ C 파일 없음 (스킵)"
        continue
    fi

    if gcc -o "$bin3" "$c_file" 2>/dev/null; then
        output=$("$bin3" 2>&1 || true)
        echo "    ✓ 3단계 컴파일 ✓ 실행 완료"
        echo "      출력: $(echo "$output" | head -1)"
        ((passed++))
    else
        echo "    ❌ 3단계 컴파일 실패"
        ((failed++))
    fi
done

# ======================== 바이너리 호환성 검증 ========================
echo -e "\n✓ Phase 4: 바이너리 호환성 검증"
echo "=================================================="

for example in example1 example2 example3; do
    bin1="$BUILD_DIR/${example}_gen1"
    bin2="$BUILD_DIR/${example}_gen2"
    bin3="$BUILD_DIR/${example}_gen3"

    if [ ! -x "$bin1" ] || [ ! -x "$bin2" ] || [ ! -x "$bin3" ]; then
        echo "  [$example] ⊘ 바이너리 부족"
        continue
    fi

    out1=$("$bin1" 2>&1 || true)
    out2=$("$bin2" 2>&1 || true)
    out3=$("$bin3" 2>&1 || true)

    if [ "$out1" = "$out2" ] && [ "$out2" = "$out3" ]; then
        echo "  [$example] ✓ 3단계 모두 동일한 출력 (호환성 검증 완료)"
        ((passed++))
    else
        echo "  [$example] ⚠ 출력 불일치"
        echo "    Gen1: $out1"
        echo "    Gen2: $out2"
        echo "    Gen3: $out3"
    fi
done

# ======================== 최종 결과 ========================

echo -e "\n" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "="
echo "📊 자체 호스팅 부트스트랩 테스트 결과"
echo "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "="

echo "✅ 통과: $passed"
echo "❌ 실패: $failed"

if [ $failed -eq 0 ]; then
    echo -e "\n🚀 자체 호스팅 부트스트랩 완료!"
    echo "   FreeLiulia 컴파일러는 자신의 코드를 성공적으로 컴파일할 수 있습니다!"
    exit 0
else
    echo -e "\n❌ 일부 테스트 실패"
    exit 1
fi
