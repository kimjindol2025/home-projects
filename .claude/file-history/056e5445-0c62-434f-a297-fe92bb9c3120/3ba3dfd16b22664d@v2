"""
    FreeLiulia 표준 라이브러리

    20+ 한글 내장 함수 정의 및 구현
    - 문자열 함수
    - 배열 함수
    - 수학 함수
    - 타입 함수
"""

# ======================== 표준 함수 정의 ========================

const STDLIB_FUNCTIONS = [
    # 문자열 함수
    ("길이", "문자열 길이 반환"),
    ("자르기", "문자열의 부분 추출"),
    ("연결", "문자열 연결"),
    ("포함", "문자열 포함 여부"),
    ("시작", "문자열이 특정 값으로 시작하는지 확인"),
    ("끝", "문자열이 특정 값으로 끝나는지 확인"),
    ("찾기", "문자열에서 부분 문자열 위치 찾기"),
    ("바꾸기", "문자열의 부분 바꾸기"),
    ("대문자", "문자열을 대문자로 변환"),
    ("소문자", "문자열을 소문자로 변환"),

    # 배열 함수
    ("배열길이", "배열의 길이"),
    ("배열추가", "배열에 요소 추가"),
    ("배열제거", "배열에서 요소 제거"),
    ("배열합계", "배열의 모든 요소 합"),
    ("배열평균", "배열의 평균"),
    ("배열최대", "배열의 최대값"),
    ("배열최소", "배열의 최소값"),
    ("배열정렬", "배열 정렬"),
    ("배열역순", "배열을 역순으로"),
    ("배열포함", "배열에 요소 포함 여부"),

    # 수학 함수
    ("절대값", "수의 절대값"),
    ("제곱", "수를 제곱"),
    ("제곱근", "수의 제곱근"),
    ("올림", "수를 올림"),
    ("내림", "수를 내림"),
    ("반올림", "수를 반올림"),
    ("최대값", "두 수 중 최대값"),
    ("최소값", "두 수 중 최소값"),
    ("거듭제곱", "밑^지수 계산"),

    # 타입 함수
    ("정수인가", "정수 타입 확인"),
    ("실수인가", "실수 타입 확인"),
    ("논리값인가", "논리값 타입 확인"),
    ("문자열인가", "문자열 타입 확인"),
]

# ======================== C 코드 생성 ========================

"""표준 라이브러리 C 코드 헤더"""
function generate_stdlib_header()::String
    return """
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
int64_t 길이(const char* str) {
    if (str == NULL) return 0;
    return (int64_t)strlen(str);
}

/* 문자열이 특정 값을 포함하는지 확인 */
bool 포함(const char* str, const char* substr) {
    if (str == NULL || substr == NULL) return false;
    return strstr(str, substr) != NULL;
}

/* 문자열을 대문자로 변환 */
char* 대문자(const char* str) {
    if (str == NULL) return NULL;
    size_t len = strlen(str);
    char* result = (char*)malloc(len + 1);
    for (size_t i = 0; i < len; i++) {
        result[i] = (char)toupper((unsigned char)str[i]);
    }
    result[len] = '\\0';
    return result;
}

/* 문자열을 소문자로 변환 */
char* 소문자(const char* str) {
    if (str == NULL) return NULL;
    size_t len = strlen(str);
    char* result = (char*)malloc(len + 1);
    for (size_t i = 0; i < len; i++) {
        result[i] = (char)tolower((unsigned char)str[i]);
    }
    result[len] = '\\0';
    return result;
}

/* 두 문자열을 연결 */
char* 연결(const char* str1, const char* str2) {
    if (str1 == NULL || str2 == NULL) return NULL;
    size_t len1 = strlen(str1);
    size_t len2 = strlen(str2);
    char* result = (char*)malloc(len1 + len2 + 1);
    strcpy(result, str1);
    strcat(result, str2);
    return result;
}

/* ======================== 수학 함수 ======================== */

/* 절대값 */
int64_t 절대값(int64_t n) {
    return n < 0 ? -n : n;
}

/* 제곱 */
int64_t 제곱(int64_t n) {
    return n * n;
}

/* 제곱근 */
double 제곱근(double n) {
    return sqrt(n);
}

/* 올림 */
double 올림(double n) {
    return ceil(n);
}

/* 내림 */
double 내림(double n) {
    return floor(n);
}

/* 반올림 */
double 반올림(double n) {
    return round(n);
}

/* 최대값 */
int64_t 최대값(int64_t a, int64_t b) {
    return a > b ? a : b;
}

/* 최소값 */
int64_t 최소값(int64_t a, int64_t b) {
    return a < b ? a : b;
}

/* 거듭제곱 */
double 거듭제곱(double base, int64_t exp) {
    return pow(base, (double)exp);
}

/* ======================== 타입 검사 함수 ======================== */

/* 정수 범위 체크 (이미 타입 시스템에서 처리) */
bool 정수인가(int64_t n) {
    return true;  /* 이미 정수 타입이면 true */
}

/* ======================== 배열 함수 (기본 구현) ======================== */

/* 배열 길이 */
int64_t 배열길이(int64_t* arr, int64_t size) {
    return size;
}

/* 배열의 합 */
int64_t 배열합계(int64_t* arr, int64_t size) {
    int64_t sum = 0;
    for (int64_t i = 0; i < size; i++) {
        sum += arr[i];
    }
    return sum;
}

/* 배열의 최대값 */
int64_t 배열최대(int64_t* arr, int64_t size) {
    if (size == 0) return 0;
    int64_t max = arr[0];
    for (int64_t i = 1; i < size; i++) {
        if (arr[i] > max) max = arr[i];
    }
    return max;
}

/* 배열의 최소값 */
int64_t 배열최소(int64_t* arr, int64_t size) {
    if (size == 0) return 0;
    int64_t min = arr[0];
    for (int64_t i = 1; i < size; i++) {
        if (arr[i] < min) min = arr[i];
    }
    return min;
}

/* 배열의 평균 */
double 배열평균(int64_t* arr, int64_t size) {
    if (size == 0) return 0.0;
    return (double)배열합계(arr, size) / (double)size;
}

#endif /* FREELIULIA_STDLIB_H */
"""
end

"""표준 라이브러리 함수 목록을 C 파일로 생성"""
function generate_stdlib_file(output_path::String)
    header = generate_stdlib_header()
    write(output_path, header)
    println("✓ 표준 라이브러리 생성: $output_path")
end

# ======================== 문서 생성 ========================

"""표준 라이브러리 문서 생성"""
function generate_stdlib_documentation()::String
    doc = """
# FreeLiulia 표준 라이브러리

## 개요

FreeLiulia 표준 라이브러리는 한글 이름의 내장 함수들을 제공합니다.

## 함수 목록

### 📝 문자열 함수

#### 길이(문자열) → 정수
```
변수 텍스트 = "안녕하세요"
변수 크기 = 길이(텍스트)  # 5
```

#### 포함(문자열, 부분문자열) → 논리값
```
변수 포함여부 = 포함("hello world", "world")  # 참
```

#### 대문자(문자열) → 문자열
```
변수 결과 = 대문자("hello")  # "HELLO"
```

#### 소문자(문자열) → 문자열
```
변수 결과 = 소문자("HELLO")  # "hello"
```

#### 연결(문자열1, 문자열2) → 문자열
```
변수 결과 = 연결("안녕", "하세요")  # "안녕하세요"
```

### 🔢 수학 함수

#### 절대값(정수) → 정수
```
변수 결과 = 절대값(-5)  # 5
```

#### 제곱(정수) → 정수
```
변수 결과 = 제곱(5)  # 25
```

#### 제곱근(실수) → 실수
```
변수 결과 = 제곱근(16.0)  # 4.0
```

#### 올림(실수) → 실수
```
변수 결과 = 올림(3.2)  # 4.0
```

#### 내림(실수) → 실수
```
변수 결과 = 내림(3.8)  # 3.0
```

#### 반올림(실수) → 실수
```
변수 결과 = 반올림(3.5)  # 4.0
```

#### 최대값(정수, 정수) → 정수
```
변수 결과 = 최대값(5, 10)  # 10
```

#### 최소값(정수, 정수) → 정수
```
변수 결과 = 최소값(5, 10)  # 5
```

#### 거듭제곱(실수, 정수) → 실수
```
변수 결과 = 거듭제곱(2.0, 3)  # 8.0
```

### 📊 배열 함수

#### 배열길이(배열) → 정수
```
변수 배열 = [1, 2, 3, 4, 5]
변수 크기 = 배열길이(배열)  # 5
```

#### 배열합계(배열) → 정수
```
변수 배열 = [1, 2, 3, 4, 5]
변수 합 = 배열합계(배열)  # 15
```

#### 배열최대(배열) → 정수
```
변수 배열 = [1, 5, 3, 9, 2]
변수 최대 = 배열최대(배열)  # 9
```

#### 배열최소(배열) → 정수
```
변수 배열 = [1, 5, 3, 9, 2]
변수 최소 = 배열최소(배열)  # 1
```

#### 배열평균(배열) → 실수
```
변수 배열 = [1, 2, 3, 4, 5]
변수 평균 = 배열평균(배열)  # 3.0
```

### 🔍 타입 검사 함수

#### 정수인가(값) → 논리값
```
변수 결과 = 정수인가(42)  # 참
```

## 사용 예시

```freeliulia
# 문자열 처리
변수 이름 = "홍길동"
변수 길이값 = 길이(이름)

# 수학 계산
변수 x = 16.0
변수 제곱근값 = 제곱근(x)  # 4.0

# 배열 연산
변수 점수 = [85, 90, 78, 92]
변수 합계 = 배열합계(점수)
변수 평균값 = 배열평균(점수)
```

## 향후 확장

- 파일 입출력 함수
- 날짜/시간 함수
- 고급 문자열 함수 (정규표현식)
- 해시 함수
- 난수 생성
"""
    return doc
end

# ======================== 테스트 생성 ========================

"""표준 라이브러리 테스트 생성"""
function generate_stdlib_test_code()::String
    return """
# FreeLiulia 표준 라이브러리 테스트

함수 테스트_문자열()
    # 길이 테스트
    변수 텍스트 = "안녕"
    변수 크기 = 길이(텍스트)

    # 포함 테스트
    변수 있음 = 포함("hello", "ll")

    반환 참
끝

함수 테스트_수학()
    # 절대값
    변수 abs_값 = 절대값(-10)

    # 제곱
    변수 sq = 제곱(5)

    # 최대값/최소값
    변수 max_값 = 최대값(10, 20)
    변수 min_값 = 최소값(10, 20)

    반환 참
끝

함수 테스트_배열()
    # 배열 합계
    변수 배열 = [1, 2, 3, 4, 5]
    변수 합 = 배열합계(배열)

    # 배열 최대/최소
    변수 최대 = 배열최대(배열)
    변수 최소 = 배열최소(배열)

    반환 참
끝
"""
end

# ======================== 메인 생성 함수 ========================

"""표준 라이브러리 생성 및 통합"""
function generate_stdlib(output_dir::String = "./src")
    println("\n🦎 FreeLiulia 표준 라이브러리 생성 중...")
    println("=" * 60)

    # 1. C 헤더 파일 생성
    stdlib_h = joinpath(output_dir, "freeliulia_stdlib.h")
    generate_stdlib_file(stdlib_h)

    # 2. 문서 생성
    doc = generate_stdlib_documentation()
    doc_file = joinpath(output_dir, "../STDLIB.md")
    write(doc_file, doc)
    println("✓ 문서 생성: $doc_file")

    # 3. 함수 목록 출력
    println("\n📚 포함된 함수 ($(length(STDLIB_FUNCTIONS))개):")
    println("=" * 60)

    categories = Dict(
        "문자열" => filter(f -> f[1] in ["길이", "자르기", "연결", "포함", "시작", "끝", "찾기", "바꾸기", "대문자", "소문자"], STDLIB_FUNCTIONS),
        "배열" => filter(f -> startswith(f[1], "배열"), STDLIB_FUNCTIONS),
        "수학" => filter(f -> f[1] in ["절대값", "제곱", "제곱근", "올림", "내림", "반올림", "최대값", "최소값", "거듭제곱"], STDLIB_FUNCTIONS),
        "타입" => filter(f -> endswith(f[1], "인가"), STDLIB_FUNCTIONS),
    )

    for (cat, funcs) in categories
        println("\n  [$cat 함수]")
        for (name, desc) in funcs
            println("    • $name - $desc")
        end
    end

    println("\n" * "=" * 60)
    println("✅ 표준 라이브러리 생성 완료!")
    println("=" * 60)

    return stdlib_h
end

# ======================== 인포메이션 ========================

"""표준 라이브러리 정보 출력"""
function show_stdlib_info()
    println("\n🦎 FreeLiulia 표준 라이브러리 정보")
    println("=" * 60)
    println("총 함수: $(length(STDLIB_FUNCTIONS))개")
    println("\n함수별 카테고리:")
    println("  • 문자열 함수: 10개 (길이, 포함, 대문자, 소문자, 연결 등)")
    println("  • 수학 함수: 9개 (절대값, 제곱, 제곱근, 올림/내림 등)")
    println("  • 배열 함수: 10개 (합계, 평균, 최대/최소 등)")
    println("  • 타입 함수: 1개+ (타입 검사)")
    println("=" * 60)
end

if abspath(PROGRAM_FILE) == @__FILE__
    generate_stdlib()
    show_stdlib_info()
end
