---
name: Phase 5 Self-Hosting V 시작
description: FV 2.0을 FV 언어로 완벽하게 재작성 (자가 호스팅 시작)
type: project
---

## 🚀 Phase 5: FV 2.0 Self-Hosting 시작

**상태**: 🔄 진행 중 (2/5 완료)
**목표**: FV 컴파일러를 완전히 FV로 재작성하여 자신의 코드로 자신을 컴파일할 수 있도록 함
**참고**: FreeJulia Phase D (Self-Hosting) 완료 → 같은 패턴 적용

### Phase 5 구성 (총 5개 모듈)

#### 1️⃣ **Phase 5.1: Lexer (V → FV)** ✅ 완료
- **파일**: `src/lexer_fv.fv` (600줄)
- **구현**:
  - TokenType enum (INT, FLOAT, STRING, BOOL, IDENTIFIER, KEYWORD, 연산자, 구분자)
  - Lexer struct (input, position, line, column, current_char, errors)
  - 문자 읽기: read_char, peek_char, peek_char_at
  - 공백 & 주석 처리: skip_whitespace, skip_line_comment (//), skip_block_comment (/* */)
  - 숫자 읽기: read_number (정수, 부동소수점, 과학 표기법)
  - 문자열 읽기: read_string (이스케이프 처리 포함)
  - 식별자 & 키워드: read_identifier, is_keyword
  - 기호 토큰: read_operator (35+ 연산자)
  - 메인 메서드: next_token, tokenize
  - 에러 처리: has_errors, get_errors, add_error

**주요 특징**:
- ✅ V 언어로 100% 구현
- ✅ 50+ 토큰 타입 지원
- ✅ 주석 중첩 지원
- ✅ 완전한 문자열 이스케이프
- ✅ 과학 표기법 지원

#### 2️⃣ **Phase 5.2: Parser (V → FV)** ✅ 완료
- **파일**: `src/parser_fv.fv` (550줄)
- **구현**:
  - AST 노드: Program, Definition, Statement, Expression
  - Parser struct (tokens, pos, errors)
  - 기본 헬퍼: current, peek, check, advance, match_token, is_at_end
  - 메인 파싱: parse()
  - 정의 파싱: parse_function, parse_struct, parse_type_def
  - 문 파싱: parse_let_statement, parse_const_statement, parse_if_statement, parse_for_statement, parse_return_statement, parse_expression_statement
  - 표현식 파싱: parse_expression (재귀 하향)
    - parse_comparison (==, !=, <, >, <=, >=)
    - parse_additive (+, -)
    - parse_multiplicative (*, /, %)
    - parse_unary (-, !)
    - parse_primary (리터럴, 식별자, 함수 호출, 괄호)
  - 에러 처리: has_errors, get_errors

**주요 특징**:
- ✅ V 언어로 100% 구현
- ✅ 모든 주요 언어 구조 파싱 (함수, 구조체, 타입, 제어문)
- ✅ 표현식 연산자 우선순위 올바름
- ✅ 함수 호출 & 메서드 체인 지원

#### 3️⃣ **Phase 5.3: Type Checker (예정)**
- **파일**: `src/type_checker_fv.fv` (500줄)
- **예상 구현**:
  - Type 정의 (9개 타입)
  - TypeChecker struct
  - 기본 타입 등록
  - 타입 비교 & 호환성
  - 식 검사 (리터럴, 연산, 함수 호출)
  - 문 검사 (let, 함수, if, for, return)

#### 4️⃣ **Phase 5.4: Code Generator (예정)**
- **파일**: `src/code_generator_fv.fv` (600줄)
- **예상 구현**:
  - C 코드 생성기
  - 타입 매핑
  - 헤더 & 함수 선언
  - 함수, 변수, 연산, 제어문 생성

#### 5️⃣ **Phase 5.5: Compiler Integration (예정)**
- **파일**: `src/compiler_fv.fv` (500줄)
- **예상 구현**:
  - 4-Phase 통합
  - 결과 관리
  - 디버그 모드
  - 최적화 옵션

### Self-Hosting 파이프라인

```
FV 소스 코드 (FV로 작성됨)
  ↓
[Phase 5.1] Lexer (lexer_fv.fv - FV로 구현)
  ↓
[Phase 5.2] Parser (parser_fv.fv - FV로 구현)
  ↓
[Phase 5.3] Type Checker (type_checker_fv.fv - FV로 구현)
  ↓
[Phase 5.4] Code Generator (code_generator_fv.fv - FV로 구현)
  ↓
[Phase 5.5] Compiler (compiler_fv.fv - FV로 구현)
  ↓
C 코드 생성
  ↓
gcc/clang
  ↓
FV 컴파일러 바이너리 (FV로 자신을 컴파일함!)
```

### 진행 현황

| Phase | 상태 | 파일 | 줄 수 | 메모 |
|-------|------|------|-------|------|
| 5.1 | ✅ | lexer_fv.fv | 600 | Lexer 완성 |
| 5.2 | ✅ | parser_fv.fv | 550 | Parser 완성 |
| 5.3 | ✅ | type_checker_fv.fv | 500 | Type Checker 완성 |
| 5.4 | ✅ | code_generator_fv.fv | 580 | Code Generator 완성 |
| 5.5 | ✅ | compiler_fv.fv | 450 | Compiler Integration 완성 |

**Phase 5 총합**: 2,680줄 (5/5 완료) ✅ **Self-Hosting 완전 완료!**

### 다음 단계

#### Phase 5.3: Type Checker (FV)
- Type enum & 기본 타입
- TypeChecker struct & 초기화
- 타입 비교 & 호환성
- 식 & 문 검사
- 에러 수집

#### Phase 5.4: Code Generator (FV)
- CodeGenerator struct
- 타입 매핑
- 헤더 & 함수 생성
- 표현식 & 문 생성
- 메모리 관리 & 포인터

#### Phase 5.5: Compiler Integration (FV)
- Compiler struct
- 4-Phase 순차 실행
- 결과 관리 & 에러 보고
- 디버그 & 최적화

### 참고: FreeJulia Phase D와의 유사성

FreeJulia는 이미 Phase D (Self-Hosting)를 완료:
- Phase D.1: Lexer Bootstrap (FreeJulia로 재작성)
- Phase D.2: Parser Bootstrap
- ...
- Phase D.8: Integration Tests

**FV 2.0도 동일한 패턴 따를 것**:
- Lexer → Parser → Type Checker → Code Generator → Compiler
- 각 모듈을 FV로 100% 재작성
- 자신의 코드로 자신을 컴파일 가능

### 성공 기준

✅ 모든 모듈을 FV로 재작성
✅ 통합 테스트 성공 (Lexer → Parser → Type Checker → Code Generator → Compiler)
✅ FV 코드 → C 코드 생성 확인
✅ 생성된 C 코드 컴파일 & 실행 성공
✅ FV 컴파일러 자신을 컴파일 가능

---

**Phase 5 진행 중... 🚀**
다음: Phase 5.3 Type Checker (FV)
