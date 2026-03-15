#ifndef FREELIULIA_STDLIB_H
#define FREELIULIA_STDLIB_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <ctype.h>
#include <stdbool.h>
#include <stdint.h>

/* ======================== 문자열 함수 ======================== */

/* 문자열 길이 */
static inline int64_t 길이(const char* str) {
    if (str == NULL) return 0;
    return (int64_t)strlen(str);
}

/* 문자열이 특정 값을 포함하는지 확인 */
static inline bool 포함(const char* str, const char* substr) {
    if (str == NULL || substr == NULL) return false;
    return strstr(str, substr) != NULL;
}

/* 문자열이 특정 값으로 시작하는지 확인 */
static inline bool 시작(const char* str, const char* prefix) {
    if (str == NULL || prefix == NULL) return false;
    return strncmp(str, prefix, strlen(prefix)) == 0;
}

/* 문자열이 특정 값으로 끝나는지 확인 */
static inline bool 끝(const char* str, const char* suffix) {
    if (str == NULL || suffix == NULL) return false;
    size_t str_len = strlen(str);
    size_t suffix_len = strlen(suffix);
    if (suffix_len > str_len) return false;
    return strcmp(str + str_len - suffix_len, suffix) == 0;
}

/* 부분 문자열의 위치 찾기 (없으면 -1) */
static inline int64_t 찾기(const char* str, const char* substr) {
    if (str == NULL || substr == NULL) return -1;
    char* pos = strstr(str, substr);
    if (pos == NULL) return -1;
    return (int64_t)(pos - str);
}

/* 문자열을 대문자로 변환 */
static inline char* 대문자(const char* str) {
    if (str == NULL) return NULL;
    size_t len = strlen(str);
    char* result = (char*)malloc(len + 1);
    for (size_t i = 0; i < len; i++) {
        result[i] = (char)toupper((unsigned char)str[i]);
    }
    result[len] = '\0';
    return result;
}

/* 문자열을 소문자로 변환 */
static inline char* 소문자(const char* str) {
    if (str == NULL) return NULL;
    size_t len = strlen(str);
    char* result = (char*)malloc(len + 1);
    for (size_t i = 0; i < len; i++) {
        result[i] = (char)tolower((unsigned char)str[i]);
    }
    result[len] = '\0';
    return result;
}

/* 두 문자열을 연결 */
static inline char* 연결(const char* str1, const char* str2) {
    if (str1 == NULL || str2 == NULL) return NULL;
    size_t len1 = strlen(str1);
    size_t len2 = strlen(str2);
    char* result = (char*)malloc(len1 + len2 + 1);
    strcpy(result, str1);
    strcat(result, str2);
    return result;
}

/* 부분 문자열 추출 (start부터 length까지) */
static inline char* 자르기(const char* str, int64_t start, int64_t length) {
    if (str == NULL) return NULL;
    size_t str_len = strlen(str);
    if (start < 0 || start >= (int64_t)str_len) return NULL;

    char* result = (char*)malloc(length + 1);
    strncpy(result, str + start, length);
    result[length] = '\0';
    return result;
}

/* 문자열의 일부를 다른 것으로 바꾸기 */
static inline char* 바꾸기(const char* str, const char* old, const char* new) {
    if (str == NULL || old == NULL || new == NULL) return NULL;

    size_t old_len = strlen(old);
    size_t new_len = strlen(new);
    size_t str_len = strlen(str);

    /* 단순 구현: 첫 번째 발견만 바꾸기 */
    char* pos = strstr(str, old);
    if (pos == NULL) return (char*)str;

    size_t before_len = pos - str;
    char* result = (char*)malloc(before_len + new_len + strlen(pos + old_len) + 1);

    strncpy(result, str, before_len);
    strcpy(result + before_len, new);
    strcpy(result + before_len + new_len, pos + old_len);

    return result;
}

/* ======================== 수학 함수 ======================== */

/* 절대값 */
static inline int64_t 절대값(int64_t n) {
    return n < 0 ? -n : n;
}

/* 절대값 (실수) */
static inline double 절대값_실수(double n) {
    return n < 0.0 ? -n : n;
}

/* 제곱 */
static inline int64_t 제곱(int64_t n) {
    return n * n;
}

/* 제곱근 */
static inline double 제곱근(double n) {
    return sqrt(n);
}

/* 올림 */
static inline double 올림(double n) {
    return ceil(n);
}

/* 내림 */
static inline double 내림(double n) {
    return floor(n);
}

/* 반올림 */
static inline double 반올림(double n) {
    return round(n);
}

/* 최대값 */
static inline int64_t 최대값(int64_t a, int64_t b) {
    return a > b ? a : b;
}

/* 최소값 */
static inline int64_t 최소값(int64_t a, int64_t b) {
    return a < b ? a : b;
}

/* 거듭제곱 */
static inline double 거듭제곱(double base, int64_t exp) {
    return pow(base, (double)exp);
}

/* ======================== 배열 함수 (기본 구현) ======================== */

/* 배열 길이 */
static inline int64_t 배열길이(int64_t* arr, int64_t size) {
    return size;
}

/* 배열의 합 */
static inline int64_t 배열합계(int64_t* arr, int64_t size) {
    int64_t sum = 0;
    for (int64_t i = 0; i < size; i++) {
        sum += arr[i];
    }
    return sum;
}

/* 배열의 최대값 */
static inline int64_t 배열최대(int64_t* arr, int64_t size) {
    if (size == 0) return 0;
    int64_t max = arr[0];
    for (int64_t i = 1; i < size; i++) {
        if (arr[i] > max) max = arr[i];
    }
    return max;
}

/* 배열의 최소값 */
static inline int64_t 배열최소(int64_t* arr, int64_t size) {
    if (size == 0) return 0;
    int64_t min = arr[0];
    for (int64_t i = 1; i < size; i++) {
        if (arr[i] < min) min = arr[i];
    }
    return min;
}

/* 배열의 평균 */
static inline double 배열평균(int64_t* arr, int64_t size) {
    if (size == 0) return 0.0;
    return (double)배열합계(arr, size) / (double)size;
}

/* 배열에 요소 포함 여부 */
static inline bool 배열포함(int64_t* arr, int64_t size, int64_t value) {
    for (int64_t i = 0; i < size; i++) {
        if (arr[i] == value) return true;
    }
    return false;
}

/* ======================== 고급 문자열 함수 ======================== */

/* 문자열 반복 (N회 반복) */
static inline char* 반복(const char* str, int64_t count) {
    if (str == NULL || count <= 0) {
        char* empty = (char*)malloc(1);
        empty[0] = '\0';
        return empty;
    }

    size_t str_len = strlen(str);
    size_t total_len = str_len * count;
    char* result = (char*)malloc(total_len + 1);

    for (int64_t i = 0; i < count; i++) {
        strcpy(result + i * str_len, str);
    }
    result[total_len] = '\0';

    return result;
}

/* 문자열 채우기 (왼쪽 공백 패딩으로 지정 길이 도달) */
static inline char* 채우기(const char* str, int64_t 길이) {
    if (str == NULL) return NULL;

    size_t str_len = strlen(str);
    if ((int64_t)str_len >= 길이) {
        /* 이미 충분히 길면 원본 복사 */
        char* result = (char*)malloc(str_len + 1);
        strcpy(result, str);
        return result;
    }

    int64_t pad_count = 길이 - str_len;
    char* result = (char*)malloc(길이 + 1);

    /* 앞에 공백 추가 */
    for (int64_t i = 0; i < pad_count; i++) {
        result[i] = ' ';
    }

    /* 문자열 추가 */
    strcpy(result + pad_count, str);
    result[길이] = '\0';

    return result;
}

/* 문자열 나누기 (첫 번째 구분자로 분할, 앞부분만 반환) */
static inline char* 나누기_앞(const char* str, const char* 구분자) {
    if (str == NULL || 구분자 == NULL) return NULL;

    char* pos = strstr(str, 구분자);
    if (pos == NULL) {
        /* 구분자를 찾지 못하면 원본 전체 반환 */
        char* result = (char*)malloc(strlen(str) + 1);
        strcpy(result, str);
        return result;
    }

    size_t before_len = pos - str;
    char* result = (char*)malloc(before_len + 1);
    strncpy(result, str, before_len);
    result[before_len] = '\0';

    return result;
}

/* 문자열 나누기 (첫 번째 구분자로 분할, 뒷부분만 반환) */
static inline char* 나누기_뒷(const char* str, const char* 구분자) {
    if (str == NULL || 구분자 == NULL) return NULL;

    char* pos = strstr(str, 구분자);
    if (pos == NULL) {
        /* 구분자를 찾지 못하면 빈 문자열 반환 */
        char* result = (char*)malloc(1);
        result[0] = '\0';
        return result;
    }

    size_t delim_len = strlen(구분자);
    size_t after_len = strlen(pos + delim_len);
    char* result = (char*)malloc(after_len + 1);
    strcpy(result, pos + delim_len);

    return result;
}

/* ======================== 타입 검사 함수 ======================== */

/* 정수 타입 검사 */
static inline bool 정수인가(int64_t n) {
    return true;  /* 이미 정수 타입이면 true */
}

/* 실수 타입 검사 */
static inline bool 실수인가(double n) {
    return true;  /* 이미 실수 타입이면 true */
}

/* 논리값 타입 검사 */
static inline bool 논리값인가(bool b) {
    return true;  /* 이미 논리값 타입이면 true */
}

/* 문자열 포인터 체크 */
static inline bool 문자열인가(const char* str) {
    return str != NULL;
}

#endif /* FREELIULIA_STDLIB_H */
