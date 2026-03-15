#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

int64_t 배열합계(int64_t n) {
  int64_t 합 = 0;

  label_1:
  int64_t t1 = 1;
  int64_t t2 = n;
  int64_t t3 = (t1 <= t2) ? 1 : 0;
  if (t3) { goto label_2; } else { goto label_3; }

  label_2:
  int64_t i = 1;

  label_4:
  int64_t t4 = i;
  int64_t t5 = n;
  int64_t t6 = (t4 <= t5) ? 1 : 0;
  if (t6) { goto label_5; } else { goto label_3; }

  label_5:
  int64_t t7 = 합;
  int64_t t8 = i;
  int64_t t9 = t7 + t8;
  합 = t9;

  int64_t t10 = i;
  int64_t t11 = 1;
  int64_t t12 = t10 + t11;
  i = t12;
  goto label_4;

  label_3:
  return 합;
}

int main() {
  int64_t 총합 = 배열합계(5);
  printf("1부터 5까지의 합: %lld\n", 총합);
  return 0;
}
