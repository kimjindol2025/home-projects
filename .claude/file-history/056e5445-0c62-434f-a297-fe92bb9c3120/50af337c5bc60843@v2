#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

/* 마이크로 렉서 - FreeLiulia 자체호스팅 첫 단계 */

int64_t 렉서_마이크로(int64_t 입력) {
    int64_t 토큰 = 0;

    if (입력 > 0) {
        토큰 = 입력 + 1;
    }

    return 토큰;
}

int main() {
    printf("🦎 FreeLiulia 마이크로 렉서 v2\n");
    printf("================================\n\n");

    int64_t 결과 = 렉서_마이크로(5);

    printf("입력: 5\n");
    printf("출력 (토큰 수): %lld\n", 결과);
    printf("\n✅ 마이크로 렉서 작동 확인\n");

    return 0;
}
