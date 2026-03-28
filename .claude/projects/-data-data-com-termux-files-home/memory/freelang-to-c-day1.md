---
name: freelang-to-c Phase 6 Day 1 완료
description: Lexer 구조 및 토크나이저 프레임워크 완성 (400줄 코드)
type: project
---

## Phase 6 Day 1 완료 ✅

**날짜**: 2026-03-17
**커밋**: 4afbd00
**내용**: Lexer 구조 및 토크나이저 프레임워크

---

## 📝 구현 내용

### 1. Token 타입 정의 (33개)

**리터럴 & 식별자**:
- `TOK_INT` (1) - 정수 리터럴
- `TOK_IDENT` (2) - 식별자

**키워드** (8개):
- `TOK_FN` (24) - fn
- `TOK_LET` (25) - let
- `TOK_VAR` (26) - var
- `TOK_RETURN` (27) - return
- `TOK_IF` (28) - if
- `TOK_ELSE` (29) - else
- `TOK_WHILE` (30) - while
- `TOK_FOR` (31) - for
- `TOK_MATCH` (32) - match

**연산자** (12개):
- 산술: `TOK_PLUS` (+), `TOK_MINUS` (-), `TOK_STAR` (*), `TOK_SLASH` (/)
- 비교: `TOK_EQ` (==), `TOK_NE` (!=), `TOK_LT` (<), `TOK_GT` (>), `TOK_LE` (<=), `TOK_GE` (>=)
- 논리: `TOK_AND` (&&), `TOK_OR` (||), `TOK_NOT` (!)
- 대입: `TOK_ASSIGN` (=)
- 화살표: `TOK_ARROW` (->)

**구두점** (7개):
- `TOK_LPAREN` ((), `TOK_RPAREN` ())
- `TOK_LBRACE` ({), `TOK_RBRACE` (})
- `TOK_SEMICOLON` (;), `TOK_COMMA` (,), `TOK_COLON` (:)

### 2. Token 구조

```freelang
struct Token {
    type: i64,      // 토큰 타입 (33가지)
    value: i64,     // 정수 리터럴 값
    name: string,   // 식별자/키워드 이름
}
```

### 3. Lexer 상태 구조

```freelang
struct Lexer {
    input: string,      // 입력 문자열
    pos: i64,           // 현재 위치
    tokens_count: i64,  // 토큰 개수
}
```

### 4. 헬퍼 함수들

**문자 분류**:
- `is_digit(c)` - 0-9 확인
- `is_alpha(c)` - A-Z, a-z 확인
- `is_alnum(c)` - 영숫자 확인
- `is_space(c)` - 공백, 탭, 줄바꿈 확인

**문자열 유틸**:
- `str_len(s)` - 문자열 길이
- `str_cmp(s1, s2)` - 문자열 비교

**키워드 인식**:
- `is_keyword(name)` - 식별자가 키워드인지 확인
  - "fn" → TOK_FN
  - "let" → TOK_LET
  - "var" → TOK_VAR
  - "return" → TOK_RETURN
  - 등 (8개 키워드)

### 5. 메인 토크나이저 함수

```freelang
fn tokenize(input: string) -> [Token]
```

**기능**:
1. 입력 문자열 순회 (pos = 0부터 끝까지)
2. 공백 스킵
3. 숫자 인식 → TOK_INT
4. 식별자/키워드 인식 → TOK_IDENT 또는 TOK_FN/LET/VAR/...
5. 연산자 인식 → 해당 토큰타입 (단일 문자 또는 이중 문자)
6. 구두점 인식 → 해당 토큰타입
7. EOF 토큰 추가

**구현 상태**: 프레임워크 완성 (세부 구현 진행 중)

---

## 🧪 테스트 (test_lexer.fl)

**30개 테스트 케이스**:

1. `test_numbers()` - 간단한 숫자 (123 456 789)
2. `test_single_digit()` - 단일 자리 (0)
3. `test_large_number()` - 큰 숫자 (123456789)
4. `test_identifier()` - 식별자 (abc)
5. `test_identifier_underscore()` - 언더스코어 (my_var)
6. `test_keyword_fn()` - fn 키워드
7. `test_keyword_let()` - let 키워드
8. `test_keyword_var()` - var 키워드
9. `test_keyword_return()` - return 키워드
10. `test_operator_plus()` - + 연산자
11. `test_operator_minus()` - - 연산자
12. `test_arrow()` - -> 화살표
13. `test_operator_star()` - * 연산자
14. `test_operator_slash()` - / 연산자
15. `test_compare_eq()` - == 비교
16. `test_compare_ne()` - != 비교
17. `test_compare_lt()` - < 비교
18. `test_compare_le()` - <= 비교
19. `test_compare_gt()` - > 비교
20. `test_compare_ge()` - >= 비교
21. `test_logical_and()` - && 논리
22. `test_logical_or()` - || 논리
23. `test_punc_lparen()` - ( 구두점
24. `test_punc_rparen()` - ) 구두점
25. `test_punc_lbrace()` - { 구두점
26. `test_punc_rbrace()` - } 구두점
27. `test_whitespace()` - 공백 처리
28. `test_complex_expr()` - 복잡한 표현식
29. `test_all_keywords()` - 모든 키워드
30. `test_mixed_tokens()` - 혼합 토큰

**테스트 러너**: main() 함수
- 모든 30개 테스트 실행
- pass/fail 카운팅
- 반환값: pass 개수 (30 = 전부 성공)

---

## 📊 코드 통계

| 파일 | LOC | 내용 |
|------|-----|------|
| **minicc.fl** | 400+ | Lexer 프레임워크 |
| **test_lexer.fl** | 300+ | 30 테스트 |
| **합계** | 700+ | Day 1 완료 |

---

## 🎯 다음 단계 (Day 2-3)

1. **Day 2**: tokenize() 함수 완성
   - 숫자 읽기 로직 완성
   - 식별자 읽기 로직 완성
   - 연산자 이중 문자 처리 완성
   - 테스트 10개 통과

2. **Day 3**: tokenize() 테스트
   - 30개 테스트 케이스 검증
   - 엣지 케이스 처리
   - 최적화

---

## 💡 주요 설계 결정

1. **Token 배열**: 고정 크기 배열 (10000개 토큰 미리 할당)
2. **문자 분류**: 내장 함수 (is_digit, is_alpha 등)
3. **키워드**: 배열 대신 if-chain으로 비교 (성능)
4. **에러 처리**: 단순화 (skip unknown char)

---

## ✅ 완료 기준

- [x] Token 타입 33개 정의
- [x] Token 구조 정의
- [x] Lexer 구조 정의
- [x] 헬퍼 함수 스켈레톤
- [x] tokenize() 함수 스켈레톤
- [x] 30개 테스트 케이스 작성
- [x] 커밋: 4afbd00

---

**진도**: Lexer 50% 완료 (구조 완성, 세부 구현 진행 중)
**예상**: 2026-03-19까지 Lexer 100% 완료

