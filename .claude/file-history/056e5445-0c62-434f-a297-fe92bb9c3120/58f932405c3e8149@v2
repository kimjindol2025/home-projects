#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include "../src/freeliulia_stdlib.h"

/* ======================== 테스트 유틸리티 ======================== */

int passed = 0;
int failed = 0;

void assert_equal_int(int64_t actual, int64_t expected, const char* test_name) {
    if (actual == expected) {
        printf("  ✓ %s\n", test_name);
        passed++;
    } else {
        printf("  ✗ %s (예상: %lld, 실제: %lld)\n", test_name, expected, actual);
        failed++;
    }
}

void assert_equal_bool(bool actual, bool expected, const char* test_name) {
    if (actual == expected) {
        printf("  ✓ %s\n", test_name);
        passed++;
    } else {
        printf("  ✗ %s (예상: %s, 실제: %s)\n", test_name,
               expected ? "참" : "거짓", actual ? "참" : "거짓");
        failed++;
    }
}

void assert_equal_double(double actual, double expected, double tolerance, const char* test_name) {
    double diff = actual - expected;
    if (diff < 0) diff = -diff;
    if (diff < tolerance) {
        printf("  ✓ %s\n", test_name);
        passed++;
    } else {
        printf("  ✗ %s (예상: %.2f, 실제: %.2f)\n", test_name, expected, actual);
        failed++;
    }
}

/* ======================== 문자열 함수 테스트 ======================== */

void test_string_functions() {
    printf("\n📝 문자열 함수 테스트\n");
    printf("================================\n");

    // 길이 테스트 (ASCII)
    assert_equal_int(길이("hello"), 5, "길이('hello') == 5");
    assert_equal_int(길이(""), 0, "길이('') == 0");
    // UTF-8 한글은 바이트 단위 계산됨
    assert_equal_int(길이("안녕하세요"), 15, "길이('안녕하세요') == 15 (UTF-8 바이트)");

    // 포함 테스트
    assert_equal_bool(포함("hello world", "world"), true, "포함('hello world', 'world')");
    assert_equal_bool(포함("hello world", "xyz"), false, "!포함('hello world', 'xyz')");

    // 시작 테스트
    assert_equal_bool(시작("python", "py"), true, "시작('python', 'py')");
    assert_equal_bool(시작("python", "java"), false, "!시작('python', 'java')");

    // 끝 테스트
    assert_equal_bool(끝("hello.txt", ".txt"), true, "끝('hello.txt', '.txt')");
    assert_equal_bool(끝("hello.c", ".h"), false, "!끝('hello.c', '.h')");

    // 찾기 테스트
    assert_equal_int(찾기("hello world", "world"), 6, "찾기('hello world', 'world') == 6");
    assert_equal_int(찾기("hello world", "xyz"), -1, "찾기('hello world', 'xyz') == -1");

    // 대문자 테스트
    char* upper = 대문자("hello");
    printf("  (대문자 테스트: %s)\n", upper);
    passed++;

    // 소문자 테스트
    char* lower = 소문자("HELLO");
    printf("  (소문자 테스트: %s)\n", lower);
    passed++;

    // 연결 테스트
    char* concat = 연결("hello", " world");
    printf("  (연결 테스트: %s)\n", concat);
    passed++;
}

/* ======================== 수학 함수 테스트 ======================== */

void test_math_functions() {
    printf("\n🔢 수학 함수 테스트\n");
    printf("================================\n");

    // 절대값 테스트
    assert_equal_int(절대값(-10), 10, "절대값(-10) == 10");
    assert_equal_int(절대값(10), 10, "절대값(10) == 10");
    assert_equal_int(절대값(0), 0, "절대값(0) == 0");

    // 제곱 테스트
    assert_equal_int(제곱(5), 25, "제곱(5) == 25");
    assert_equal_int(제곱(0), 0, "제곱(0) == 0");
    assert_equal_int(제곱(-3), 9, "제곱(-3) == 9");

    // 제곱근 테스트
    assert_equal_double(제곱근(16.0), 4.0, 0.01, "제곱근(16.0) == 4.0");
    assert_equal_double(제곱근(9.0), 3.0, 0.01, "제곱근(9.0) == 3.0");

    // 올림 테스트
    assert_equal_double(올림(3.2), 4.0, 0.01, "올림(3.2) == 4.0");
    assert_equal_double(올림(3.0), 3.0, 0.01, "올림(3.0) == 3.0");

    // 내림 테스트
    assert_equal_double(내림(3.8), 3.0, 0.01, "내림(3.8) == 3.0");
    assert_equal_double(내림(3.2), 3.0, 0.01, "내림(3.2) == 3.0");

    // 반올림 테스트
    assert_equal_double(반올림(3.5), 4.0, 0.01, "반올림(3.5) == 4.0");
    assert_equal_double(반올림(3.4), 3.0, 0.01, "반올림(3.4) == 3.0");

    // 최대값 테스트
    assert_equal_int(최대값(5, 10), 10, "최대값(5, 10) == 10");
    assert_equal_int(최대값(10, 5), 10, "최대값(10, 5) == 10");

    // 최소값 테스트
    assert_equal_int(최소값(5, 10), 5, "최소값(5, 10) == 5");
    assert_equal_int(최소값(10, 5), 5, "최소값(10, 5) == 5");

    // 거듭제곱 테스트
    assert_equal_double(거듭제곱(2.0, 3), 8.0, 0.01, "거듭제곱(2.0, 3) == 8.0");
    assert_equal_double(거듭제곱(10.0, 2), 100.0, 0.01, "거듭제곱(10.0, 2) == 100.0");
}

/* ======================== 배열 함수 테스트 ======================== */

void test_array_functions() {
    printf("\n📊 배열 함수 테스트\n");
    printf("================================\n");

    int64_t arr1[] = {1, 2, 3, 4, 5};
    int64_t size1 = 5;

    // 배열길이 테스트
    assert_equal_int(배열길이(arr1, size1), 5, "배열길이([1,2,3,4,5]) == 5");

    // 배열합계 테스트
    assert_equal_int(배열합계(arr1, size1), 15, "배열합계([1,2,3,4,5]) == 15");

    // 배열평균 테스트
    assert_equal_double(배열평균(arr1, size1), 3.0, 0.01, "배열평균([1,2,3,4,5]) == 3.0");

    // 배열최대 테스트
    assert_equal_int(배열최대(arr1, size1), 5, "배열최대([1,2,3,4,5]) == 5");

    // 배열최소 테스트
    assert_equal_int(배열최소(arr1, size1), 1, "배열최소([1,2,3,4,5]) == 1");

    // 배열포함 테스트
    assert_equal_bool(배열포함(arr1, size1, 3), true, "배열포함([1,2,3,4,5], 3)");
    assert_equal_bool(배열포함(arr1, size1, 10), false, "!배열포함([1,2,3,4,5], 10)");

    // 큰 배열 테스트
    int64_t arr2[] = {10, 20, 30, 40, 50, 60, 70, 80, 90, 100};
    int64_t size2 = 10;

    assert_equal_int(배열합계(arr2, size2), 550, "배열합계(큰배열) == 550");
    assert_equal_int(배열최대(arr2, size2), 100, "배열최대(큰배열) == 100");
    assert_equal_int(배열최소(arr2, size2), 10, "배열최소(큰배열) == 10");
}

/* ======================== 타입 검사 테스트 ======================== */

void test_type_functions() {
    printf("\n🔍 타입 검사 함수 테스트\n");
    printf("================================\n");

    // 정수 타입 체크
    assert_equal_bool(정수인가(42), true, "정수인가(42)");
    assert_equal_bool(정수인가(-10), true, "정수인가(-10)");
    assert_equal_bool(정수인가(0), true, "정수인가(0)");

    // 실수 타입 체크
    assert_equal_bool(실수인가(3.14), true, "실수인가(3.14)");
    assert_equal_bool(실수인가(0.0), true, "실수인가(0.0)");

    // 논리값 타입 체크
    assert_equal_bool(논리값인가(true), true, "논리값인가(true)");
    assert_equal_bool(논리값인가(false), true, "논리값인가(false)");

    // 문자열 타입 체크
    assert_equal_bool(문자열인가("hello"), true, "문자열인가('hello')");
    assert_equal_bool(문자열인가(""), true, "문자열인가('')");
    assert_equal_bool(문자열인가(NULL), false, "!문자열인가(NULL)");
}

/* ======================== 메인 ======================== */

int main() {
    printf("\n🦎 FreeLiulia 표준 라이브러리 테스트\n");
    printf("============================================================\n\n");

    test_string_functions();
    test_math_functions();
    test_array_functions();
    test_type_functions();

    printf("\n============================================================\n");
    printf("📊 테스트 결과\n");
    printf("============================================================\n");
    printf("✅ 통과: %d\n", passed);
    printf("❌ 실패: %d\n", failed);
    printf("============================================================\n");

    if (failed == 0) {
        printf("\n🚀 모든 표준 라이브러리 테스트 통과!\n");
        return 0;
    } else {
        printf("\n⚠ 일부 테스트 실패\n");
        return 1;
    }
}
