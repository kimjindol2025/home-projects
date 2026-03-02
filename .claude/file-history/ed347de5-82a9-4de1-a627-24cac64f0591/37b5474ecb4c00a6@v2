# FreeLang 시작 가이드

**작성일**: 2026-03-02
**버전**: v1.0.0
**대상**: 신규 사용자 & 개발자

---

## 📚 목차

1. [설치 방법](#1-설치-방법)
2. [첫 번째 프로그램](#2-첫-번째-프로그램)
3. [개발 환경 설정](#3-개발-환경-설정)
4. [기본 예제](#4-기본-예제)
5. [자주 묻는 질문 (FAQ)](#5-자주-묻는-질문-faq)
6. [다음 단계](#6-다음-단계)

---

## 1. 설치 방법

### 1.1 시스템 요구사항

- **Node.js**: v16.0 이상
- **npm**: v7.0 이상
- **운영체제**: Linux, macOS, Windows (WSL2)

### 1.2 Option A: npm으로 설치 (권장)

```bash
# FreeLang 설치
npm install -g @freelang/cli

# 버전 확인
freelang --version
```

### 1.3 Option B: 소스에서 빌드

```bash
# 저장소 클론
git clone https://github.com/freelang/freelang.git
cd freelang

# 의존성 설치
npm install

# 빌드
npm run build

# 전역 등록
npm link

# 설치 확인
freelang --version
```

### 1.4 Option C: Docker로 실행

```bash
# Docker 이미지 가져오기
docker pull freelang:latest

# 컨테이너 실행
docker run -it freelang:latest

# FreeLang 실행
freelang hello.fl
```

---

## 2. 첫 번째 프로그램

### 2.1 프로젝트 생성

```bash
# 새 프로젝트 생성
freelang new my_project
cd my_project

# 프로젝트 구조
my_project/
├── src/
│   └── main.fl           # 주 프로그램
├── tests/
├── freelang.toml         # 프로젝트 설정
└── README.md
```

### 2.2 Hello World

**파일**: `src/main.fl`

```freelang
fn main
  intent: "프로그램 진입점"
  output: ()
  do
    println("Hello, FreeLang!")
```

### 2.3 실행

```bash
# 방법 1: 직접 실행 (해석)
freelang run src/main.fl

# 방법 2: 컴파일 후 실행
freelang compile src/main.fl
freelang src/main.compiled

# 방법 3: npm 스크립트
npm run start
```

**출력**:
```
Hello, FreeLang!
```

---

## 3. 개발 환경 설정

### 3.1 에디터 설정

#### VS Code 설정

1. **확장 프로그램 설치**:
   ```
   Command Palette (Ctrl+Shift+P) → "Extensions: Install Extensions"
   검색: "FreeLang"
   "FreeLang Language Support" 설치
   ```

2. **문법 강조**:
   - `.fl` 파일이 자동으로 FreeLang으로 인식됨
   - 키워드 강조, 자동 완성 지원

3. **디버거 설정**:
   ```json
   // .vscode/launch.json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "FreeLang",
         "type": "freelang",
         "request": "launch",
         "program": "${workspaceFolder}/src/main.fl",
         "console": "integratedTerminal"
       }
     ]
   }
   ```

#### Vim/Neovim 설정

1. **문법 파일 설치**:
   ```bash
   mkdir -p ~/.vim/syntax
   cp freelang.vim ~/.vim/syntax/
   ```

2. **vim 설정** (`~/.vimrc`):
   ```vim
   au BufRead,BufNewFile *.fl set filetype=freelang
   ```

### 3.2 IDE 통합

#### IntelliJ IDEA / WebStorm

1. **플러그인 설치**:
   - Plugins → Browse Repositories
   - "FreeLang" 검색 및 설치

2. **프로젝트 설정**:
   - File → New → Project
   - FreeLang 선택

### 3.3 명령어 라인 도구

```bash
# 도움말 보기
freelang help

# 버전 확인
freelang --version

# 프로젝트 생성
freelang new <project_name>

# 실행
freelang run <file.fl>

# 컴파일
freelang compile <file.fl>

# 테스트 실행
freelang test

# 형식 확인
freelang fmt

# 린트
freelang lint
```

---

## 4. 기본 예제

### 4.1 변수와 기본 타입

**파일**: `src/01_variables.fl`

```freelang
fn variables_example
  intent: "변수 선언 및 사용"
  output: ()
  do
    // 불변 변수
    let name = "Alice"
    let age = 30
    let score = 95.5
    let active = true

    // 변수 출력
    println(f"Name: {name}")
    println(f"Age: {age}")
    println(f"Score: {score}")
    println(f"Active: {active}")

    // 가변 변수
    mut counter = 0
    counter = counter + 1
    println(f"Counter: {counter}")
```

**실행**:
```bash
freelang run src/01_variables.fl
```

**출력**:
```
Name: Alice
Age: 30
Score: 95.5
Active: true
Counter: 1
```

### 4.2 배열과 반복

**파일**: `src/02_arrays.fl`

```freelang
fn arrays_example
  intent: "배열 선언 및 반복"
  output: ()
  do
    // 배열 선언
    let numbers = [1, 2, 3, 4, 5]
    let fruits = ["apple", "banana", "orange"]

    // 길이 출력
    println(f"Numbers: {numbers}")
    println(f"Length: {numbers.len()}")

    // 반복 (for-in)
    println("Fruits:")
    for fruit in fruits
      println(f"  - {fruit}")

    // 인덱스 반복
    println("Numbers with index:")
    for (i, num) in numbers.enumerate()
      println(f"  [{i}] = {num}")
```

### 4.3 함수와 제어문

**파일**: `src/03_functions.fl`

```freelang
fn sum
  intent: "배열의 합계 계산"
  input: numbers: array<number>
  output: number
  do
    result = 0
    for num in numbers
      result = result + num
    return result

fn max_value
  intent: "배열의 최대값 찾기"
  input: numbers: array<number>
  output: Option<number>
  do
    if numbers.len() == 0
      None
    else
      mut largest = numbers[0]
      for i in 1..numbers.len()
        if numbers[i] > largest
          largest = numbers[i]
      Some(largest)

fn examples
  intent: "함수 사용 예제"
  output: ()
  do
    numbers = [10, 20, 30, 40, 50]

    total = sum(numbers)
    println(f"Sum: {total}")          // Sum: 150

    match max_value(numbers)
      Some(max) =>
        println(f"Max: {max}")        // Max: 50
      None =>
        println("Array is empty")
```

### 4.4 구조체와 패턴 매칭

**파일**: `src/04_structs.fl`

```freelang
struct Person
  name: string
  age: number
  email: string

fn create_person
  intent: "새로운 Person 생성"
  input: name: string, age: number
  output: Person
  do
    return Person {
      name: name,
      age: age,
      email: f"{name.to_lowercase()}@example.com"
    }

fn is_adult
  intent: "성인 여부 판단"
  input: person: Person
  output: bool
  do
    return person.age >= 18

fn example
  intent: "구조체 사용 예제"
  output: ()
  do
    person = create_person("Alice", 30)

    println(f"Name: {person.name}")
    println(f"Age: {person.age}")
    println(f"Email: {person.email}")

    if is_adult(person)
      println("Adult")
    else
      println("Minor")
```

### 4.5 오류 처리

**파일**: `src/05_error_handling.fl`

```freelang
fn safe_divide
  intent: "안전한 나눗셈"
  input: a: number, b: number
  output: Result<number, string>
  do
    if b == 0
      Err("Division by zero")
    else
      Ok(a / b)

fn parse_number
  intent: "문자열을 숫자로 변환"
  input: str: string
  output: Result<number, string>
  do
    // 문자열이 숫자인지 확인
    match str
      "0" => Ok(0)
      "1" => Ok(1)
      "2" => Ok(2)
      _ => Err(f"Cannot parse '{str}' as number")

fn example
  intent: "오류 처리 예제"
  output: ()
  do
    // divide
    match safe_divide(10, 2)
      Ok(result) =>
        println(f"10 / 2 = {result}")
      Err(error) =>
        println(f"Error: {error}")

    match safe_divide(10, 0)
      Ok(result) =>
        println(f"Result: {result}")
      Err(error) =>
        println(f"Error: {error}")    // Error: Division by zero

    // parse
    match parse_number("42")
      Ok(num) =>
        println(f"Parsed: {num}")
      Err(error) =>
        println(f"Error: {error}")
```

---

## 5. 자주 묻는 질문 (FAQ)

### Q1. FreeLang은 인터프리터인가요, 컴파일러인가요?

**A**: FreeLang은 **컴파일러** 언어입니다.
- `.fl` 파일 → **컴파일** → TypeScript/C/bytecode
- 실행 속도가 빠릅니다
- 타입 안전성을 제공합니다

하지만 개발 편의성을 위해 **직접 실행** 모드도 지원합니다:
```bash
freelang run script.fl      # 해석 실행 (편함)
freelang compile script.fl  # 컴파일 후 실행 (빠름)
```

### Q2. Python/JavaScript와 무엇이 다른가요?

| 특징 | Python | JavaScript | FreeLang |
|------|--------|------------|----------|
| 타입 | 동적 | 동적 | 정적 |
| 문법 | 쉬움 | 중간 | 쉬움 |
| AI 친화 | ❌ | ❌ | ✅ |
| Intent 기반 | ❌ | ❌ | ✅ |
| 컴파일 | 없음 | 없음 | ✅ |

### Q3. 어떤 프로젝트에 사용할 수 있나요?

**✅ 적합한 프로젝트**:
- 데이터 처리
- 알고리즘 구현
- 비즈니스 로직
- AI 생성 코드

**❌ 부적합한 프로젝트** (Phase 1):
- 네트워크 애플리케이션
- 웹 서버
- 실시간 시스템
- 게임 개발

### Q4. 생산 환경에서 사용 가능한가요?

**A**: Phase 1은 **프로토타입/학습** 단계입니다.
- ✅ 알고리즘 검증
- ✅ 프로토타입 개발
- ⚠️ 중규모 프로젝트
- ❌ 프로덕션 고성능 시스템

Phase 2+ (Rust 런타임)부터 프로덕션 준비가 됩니다.

### Q5. 커뮤니티 지원이 있나요?

**A**: 예, 다양한 채널이 있습니다:
- 🌐 **웹사이트**: https://freelang.dev
- 💬 **Discord**: https://discord.gg/freelang
- 📧 **메일**: support@freelang.dev
- 🐛 **GitHub Issues**: https://github.com/freelang/freelang/issues

### Q6. 학습 곡선은 어떻게 되나요?

```
시간 투자    학습 곡선
   |
   |     /
   |    /
   |   /
   |__/______
     |
     기초    중급    고급
   (1-2주)  (2-4주) (1개월+)

기초: 변수, 함수, 배열, if/for
중급: 구조체, enum, 오류 처리
고급: 트레이트, 제너릭, Phase 2 기능
```

---

## 6. 다음 단계

### 6.1 학습 경로

**초급** (1-2주):
1. ✅ 설치 및 Hello World
2. ✅ 변수와 기본 타입
3. ✅ 배열과 반복
4. ✅ 함수 선언
5. ✅ 조건문 (if/else)

**중급** (2-4주):
6. 구조체 및 Enum
7. 오류 처리 (Result)
8. 패턴 매칭
9. 기본 트레이트
10. 컬렉션 조작

**고급** (1개월+):
11. 제너릭 프로그래밍
12. 고급 트레이트
13. 성능 최적화
14. Phase 2 기능

### 6.2 추천 리소스

| 리소스 | 링크 | 타입 |
|--------|------|------|
| 공식 튜토리얼 | https://docs.freelang.dev | 📚 문서 |
| 예제 저장소 | https://github.com/freelang/examples | 💻 코드 |
| 언어 사양 v1.0 | FREELANG-v1.0-SPEC.md | 📖 명세 |
| API 문서 | https://api.freelang.dev | 📋 레퍼런스 |

### 6.3 커뮤니티 참여

1. **GitHub Star** ⭐
   ```bash
   git clone https://github.com/freelang/freelang.git
   ```

2. **Discord 가입** 💬
   https://discord.gg/freelang

3. **이슈 보고** 🐛
   https://github.com/freelang/freelang/issues

4. **PR 제출** 🤝
   기여 가이드: CONTRIBUTION.md

---

## 빠른 참고

### 명령어 요약

```bash
# 설치
npm install -g @freelang/cli

# 프로젝트 생성
freelang new my_project

# 파일 실행
freelang run main.fl

# 컴파일
freelang compile main.fl

# 테스트
freelang test

# 도움말
freelang help
```

### 기본 문법 요약

```freelang
// 함수
fn function_name
  intent: "설명"
  input: param: type
  output: type
  do
    // 본체

// 변수
let x = 5          // 불변
mut y = 10         // 가변

// 제어문
if condition
  // ...
else
  // ...

for item in collection
  // ...

match value
  Pattern => // ...
  _ => // ...

// 반환
return result
```

---

**다음**: [공식 언어 사양](FREELANG-v1.0-SPEC.md)으로 진행하기

**문의**: support@freelang.dev

---

**작성자**: FreeLang 개발팀
**라이센스**: CC-BY-SA 4.0
