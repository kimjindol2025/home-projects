#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

bool 소수인가(int64_t n) {
  int64_t t1 = n;
  int64_t t2 = 2;
  bool t3 = (t1 < t2);
  if (t3) { goto then_1; } else { goto else_1; }

  then_1:
  return false;

  else_1:

  label_loop_1:
  int64_t t4 = 2;
  int64_t t5 = n;
  int64_t i = t4;

  label_cond_1:
  int64_t t6 = i;
  int64_t t7 = n;
  bool t8 = (t6 <= t7);
  if (t8) { goto body_1; } else { goto end_1; }

  body_1:
  // 만약 i * i > n -> 반환 참
  int64_t t9 = i;
  int64_t t10 = i;
  int64_t t11 = t9 * t10;
  int64_t t12 = n;
  bool t13 = (t11 > t12);
  if (t13) { goto then_2; } else { goto else_2; }

  then_2:
  return true;

  else_2:
  // 만약 n % i == 0 -> 반환 거짓
  int64_t t14 = n;
  int64_t t15 = i;
  int64_t t16 = t14 % t15;
  int64_t t17 = 0;
  bool t18 = (t16 == t17);
  if (t18) { goto then_3; } else { goto else_3; }

  then_3:
  return false;

  else_3:

  int64_t t19 = i;
  int64_t t20 = 1;
  int64_t t21 = t19 + t20;
  i = t21;
  goto label_cond_1;

  end_1:
  return true;
}

int main() {
  int64_t num = 17;
  bool is_prime = 소수인가(num);
  if (is_prime) {
    printf("17은 소수입니다!\n");
  } else {
    printf("17은 소수가 아닙니다.\n");
  }
  return 0;
}
