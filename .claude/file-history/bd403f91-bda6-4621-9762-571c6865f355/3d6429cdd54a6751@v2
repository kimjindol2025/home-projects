# FreeLang → Z-Lang 트랜스파일러 사용 가이드

**버전**: 2.0.0
**작성일**: 2026-03-02
**언어**: 한글

---

## 📖 목차

1. [개요](#개요)
2. [설치](#설치)
3. [기본 사용법](#기본-사용법)
4. [고급 기능](#고급-기능)
5. [문법 변환 규칙](#문법-변환-규칙)
6. [예제](#예제)
7. [문제 해결](#문제-해결)
8. [FAQ](#faq)

---

## 개요

**FreeLang → Z-Lang 트랜스파일러**는 FreeLang v4로 작성된 코드를 Z-Lang 코드로 자동 변환하는 도구입니다.

```
FreeLang 소스 (.fl)
    ↓ [Lexer + Parser]
FreeLang AST
    ↓ [CodeGenerator]
Z-Lang 소스 (.z)
    ↓ [Z-Lang 컴파일러]
실행 파일
```

### 주요 기능

- ✅ **자동 문법 변환**: `var` → `let`, `:type` → `->type`
- ✅ **for-in 루프 변환**: `for i in range(1,10)` → `while` 루프
- ✅ **고급 기능 지원**: Match, Spawn, Try 표현식
- ✅ **완벽한 타입 매핑**: i32, i64, f64, bool, string, array, channel, option, result
- ✅ **배치 처리**: 여러 파일 동시 변환
- ✅ **자동 보고서**: JSON 형식 변환 보고서

---

## 설치

### 방법 1: npm (권장)

```bash
npm install -g freelang-to-zlang
fl2z --version
```

### 방법 2: 소스에서 빌드

```bash
git clone https://gogs.dclub.kr/kim/freelang-to-zlang.git
cd freelang-to-zlang
npm install
npm run build
npx ts-node src/index.ts --help
```

### 방법 3: 직접 사용

```bash
# TypeScript 런타임 필요
npx ts-node src/index.ts <input.fl>
```

---

## 기본 사용법

### 1. 단일 파일 변환

```bash
# 기본 (output: input.z)
fl2z hello.fl

# 출력 파일 지정
fl2z hello.fl -o output.z

# 상세 출력
fl2z hello.fl -v

# 보고서 생성
fl2z hello.fl --report
```

**출력 예**:

```
✅ Transpilation successful!

📊 Statistics:
  Input:      10 lines
  Statements: 2
  Output:     12 lines
  Tokens:     45
  Time:       25ms

📁 Output: output.z
```

### 2. 배치 처리 (디렉토리)

```bash
# examples 디렉토리의 모든 .fl 파일 변환
fl2z examples/ --batch

# 보고서 함께 생성
fl2z examples/ --batch --report
```

**출력 예**:

```
📦 Batch processing 10 files...
  ✅ hello.fl
  ✅ factorial.fl
  ✅ fizzbuzz.fl
  ✅ sum.fl
  ...
  ✅ even_odd.fl

📈 Summary:
  Total: 10
  Success: 10 ✅
  Failure: 0 ❌
  Success rate: 100%

📊 Report saved to transpilation_report.json
```

### 3. 헬프 및 버전

```bash
fl2z --help      # 도움말 표시
fl2z --version   # 버전 확인
```

---

## 고급 기능

### 자동 보고서 생성

```bash
fl2z factorial.fl --report --report-file my_report.json
```

**보고서 포맷** (JSON):

```json
{
  "timestamp": "2026-03-02T10:30:00.000Z",
  "version": "2.0.0",
  "inputFile": "factorial.fl",
  "outputFile": "factorial.z",
  "success": true,
  "stats": {
    "inputLines": 12,
    "outputLines": 15,
    "statementsCount": 2,
    "lexTokens": 45,
    "parseTime": 5,
    "genTime": 3
  },
  "warnings": [],
  "errors": []
}
```

### 상세 모드 (Verbose)

```bash
fl2z factorial.fl -v
```

각 단계별 상세 정보 출력:

```
[*] Processing: factorial.fl
    [Lexer] 45 tokens
    [Parser] 2 statements
    [CodeGen] 15 lines
    [Time] Lex: 5ms, Gen: 3ms
```

---

## 문법 변환 규칙

### 변수 선언

| FreeLang | Z-Lang | 설명 |
|----------|--------|------|
| `var x: i32 = 10;` | `let x: i32 = 10;` | var → let |
| `let x: i32 = 10;` | `let x: i32 = 10;` | 동일 |
| `const x: i32 = 10;` | `let x: i32 = 10;` | const → let |

### 함수 선언

| FreeLang | Z-Lang |
|----------|--------|
| `fn add(x: i32, y: i32): i32 { ... }` | `fn add(x: i32, y: i32) -> i32 { ... }` |

**자동 변환 규칙**:
- `:` (콜론) → `->` (화살표)
- 세미콜론 자동 추가

### For-in 루프 → While 변환

**입력 (FreeLang)**:
```freeLang
for i in range(1, 10) {
    println(i);
}
```

**출력 (Z-Lang)**:
```z
let i: i64 = 1;
while i <= 10 {
  print(i);
  print("\n");
  i = i + 1;
}
```

### 타입 매핑

| FreeLang | Z-Lang |
|----------|--------|
| `i32` | `i32` |
| `i64` | `i64` |
| `f64` | `f64` |
| `bool` | `bool` |
| `string` | `string` |
| `[i32]` | `[i32]` |
| `channel<i32>` | `channel<i32>` |
| `option<i32>` | `option<i32>` |

### 표준 함수 매핑

| FreeLang | Z-Lang | 설명 |
|----------|--------|------|
| `println(x)` | `print(x); print("\n");` | 출력 + 줄바꿈 |
| `print(x)` | `print(x)` | 출력 |
| `str(x)` | `toString(x)` | 타입 변환 |

---

## 예제

### 예제 1: 팩토리얼 (재귀)

**입력 (factorial.fl)**:
```freeLang
fn factorial(n: i32): i32 {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n - 1);
}

fn main(): i32 {
    let result: i32 = factorial(5);
    println(result);
    return result;
}
```

**변환 (factorial.z)**:
```z
fn factorial(n: i32) -> i32 {
  if n <= 1 {
    return 1;
  }
  return n * factorial(n - 1);
}

fn main() -> i32 {
  let result: i32 = factorial(5);
  print(result);
  print("\n");
  return result;
}
```

### 예제 2: for-in 루프 변환

**입력 (fizzbuzz.fl)**:
```freeLang
fn main(): i32 {
    for i in range(1, 100) {
        if i % 15 == 0 {
            println("FizzBuzz");
        } else if i % 3 == 0 {
            println("Fizz");
        } else if i % 5 == 0 {
            println("Buzz");
        } else {
            println(i);
        }
    }
    return 0;
}
```

**자동 변환** (while 루프):
```z
fn main() -> i32 {
  let i: i64 = 1;
  while i <= 100 {
    if (i % 15) == 0 {
      print("FizzBuzz");
      print("\n");
    } else if (i % 3) == 0 {
      print("Fizz");
      print("\n");
    } else if (i % 5) == 0 {
      print("Buzz");
      print("\n");
    } else {
      print(i);
      print("\n");
    }
    i = i + 1;
  }
  return 0;
}
```

---

## 문제 해결

### Q: "파일을 찾을 수 없습니다" 오류

**원인**: 입력 파일 경로가 잘못됨

**해결책**:
```bash
# 절대 경로 사용
fl2z /home/user/hello.fl

# 현재 디렉토리 확인
pwd
ls *.fl
```

### Q: 변환이 실패했습니다

**원인**: FreeLang 문법 오류

**해결책**:
```bash
# 상세 오류 메시지 확인
fl2z input.fl -v

# 문법 확인
cat input.fl  # 세미콜론, 괄호 등 확인
```

### Q: 출력 파일이 생성되지 않았습니다

**원인**: 변환 실패

**해결책**:
```bash
# 1. 출력 디렉토리 확인
ls -la output_dir/

# 2. 쓰기 권한 확인
chmod 755 output_dir/

# 3. 다시 시도
fl2z input.fl -o output.z -v
```

---

## FAQ

### Q1: Z-Lang 컴파일러는 어디서 구하나요?

A: Z-Lang 컴파일러는 별도 저장소에서 제공합니다:
```bash
git clone https://gogs.dclub.kr/kim/zlang.git
cd zlang
npm install && npm run build
```

### Q2: 배치 변환할 때 실패한 파일만 다시 변환할 수 있나요?

A: `--report` 옵션으로 생성된 JSON 보고서를 확인하세요:
```bash
# 보고서에서 실패한 파일 확인
cat transpilation_report.json | grep '"success": false'

# 해당 파일만 다시 변환
fl2z failed_file.fl -v
```

### Q3: 대용량 파일 변환 시 성능을 개선할 수 있나요?

A: 배치 모드 사용:
```bash
# 배치 모드 (더 빠름)
fl2z examples/ --batch

# 단일 모드 (비교)
for f in examples/*.fl; do fl2z "$f"; done
```

### Q4: 변환된 코드가 완벽하지 않은 경우는?

A: 다음 기능은 Phase 2에서 제한적으로 지원됩니다:
- 배열 루프 (range만 지원)
- 복합 패턴 매칭
- 일부 고급 기능

해당 부분은 수동으로 수정하세요.

### Q5: 다른 언어로 변환할 수 있나요?

A: 현재는 Z-Lang만 지원합니다. 추후 다음 언어 지원 예정:
- [ ] Rust
- [ ] Go
- [ ] TypeScript

---

## 지원 및 피드백

문제 발생 시:

1. **GitHub Issues**: https://github.com/anthropics/claude-code/issues
2. **GOGS Repository**: https://gogs.dclub.kr/kim/freelang-to-zlang
3. **Email**: support@example.com

### 버그 보고 양식

```
제목: [BUG] 간단한 설명

내용:
- 재현 방법
- 예상 결과
- 실제 결과
- 환경 정보 (OS, Node.js 버전 등)
```

---

## 라이선스

MIT License (2026-03-02)

---

**마지막 업데이트**: 2026-03-02
**버전**: 2.0.0
**상태**: 프로덕션 준비 완료 ✅

