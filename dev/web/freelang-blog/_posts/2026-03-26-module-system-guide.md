---
title: "FreeLang의 모듈 시스템 - 코드 재사용 완벽 가이드"
date: 2026-03-26
author: Content Writer
category: Technical
tags:
  - FreeLang
  - Module System
  - Import/Export
  - Code Organization
  - Best Practices
---

# FreeLang의 모듈 시스템
## 코드 재사용을 위한 완벽 가이드

**글을 읽으면 얻을 수 있는 것:**
1. FreeLang 모듈 시스템의 핵심 개념 이해하기
2. import/export로 코드 구조화하기
3. 실제 프로젝트에서 모듈화하는 방법 배우기

---

## 배경: 왜 모듈이 필요한가?

작은 프로젝트에선 모든 코드를 한 파일에 써도 괜찮습니다. 하지만 프로젝트가 커지면?

```freelang
// 1000줄이 넘는 main.fl
// ❌ 뭐가 어디에 있는지 찾기 힘듦
// ❌ 다른 프로젝트에서 재사용 불가
// ❌ 팀원과 협업 어려움
```

이 문제를 해결하는 게 **모듈 시스템**입니다.

---

## 문제: 거대한 단일 파일의 한계

### 1. 코드 찾기 어려움
```freelang
// main.fl - 2000줄짜리 파일
// ...
fn calculate_salary(base: INT, bonus: INT) -> INT { ... }  // 700줄 쯤?
fn validate_email(email: STRING) -> BOOL { ... }          // 1500줄 쯤?
fn process_payment(amount: FLOAT) -> BOOL { ... }         // 1800줄 쯤?
// ...

// 어디가 어디지? 😱
```

### 2. 코드 재사용 불가능
```freelang
// project-a/main.fl
fn calculate_salary(base: INT, bonus: INT) -> INT {
  RETURN base + bonus
}

// project-b/main.fl
// 같은 함수를 또 작성해야 함... 😢 코드 중복
```

### 3. 협업 어려움
```freelang
// 팀원 A: salary 함수 수정 중
// 팀원 B: email 함수 수정 중
// 🔥 같은 파일에서 작업하다가 conflict!
```

---

## 해결책: FreeLang의 모듈 시스템

FreeLang은 **모듈(Module)**을 통해 코드를 논리적으로 나누고 재사용할 수 있게 합니다.

### 핵심 개념

```
모듈 = 관련된 함수/타입들의 모음
└─ 어떤 부분을 공개할지(export) 결정
└─ 다른 모듈의 코드를 가져올 수 있음(import)
```

### 기본 구조

```freelang
// salary.fl - 급여 관련 모듈
fn calculate_salary(base: INT, bonus: INT) -> INT {
  RETURN base + bonus
}

fn tax_deduction(salary: INT) -> INT {
  RETURN salary / 10
}

EXPORT calculate_salary
EXPORT tax_deduction
```

```freelang
// main.fl - 메인 프로그램
IMPORT calculate_salary FROM "salary"
IMPORT tax_deduction FROM "salary"

LET base = 1000
LET bonus = 200
LET gross = calculate_salary(base, bonus)
LET tax = tax_deduction(gross)

PRINTLN("급여: " + gross)
PRINTLN("세금: " + tax)
PRINTLN("실수령액: " + (gross - tax))
```

---

## 실전: 프로젝트 구조화하기

### 초보자라면?

간단한 계산기 프로젝트를 모듈로 나눠봅시다.

**프로젝트 구조:**
```
calculator/
├── math.fl        (수학 연산)
├── validation.fl  (입력 검증)
└── main.fl        (메인 프로그램)
```

**Step 1: math.fl - 수학 연산 모듈**

```freelang
// math.fl
fn add(a: INT, b: INT) -> INT {
  RETURN a + b
}

fn subtract(a: INT, b: INT) -> INT {
  RETURN a - b
}

fn multiply(a: INT, b: INT) -> INT {
  RETURN a * b
}

fn divide(a: INT, b: INT) -> INT {
  IF b = 0:
    RETURN 0
  ELSE:
    RETURN a / b
}

EXPORT add
EXPORT subtract
EXPORT multiply
EXPORT divide
```

**Step 2: validation.fl - 검증 모듈**

```freelang
// validation.fl
fn is_valid_number(n: STRING) -> BOOL {
  // 문자열이 숫자인지 확인
  RETURN TRUE  // 간단히 구현
}

fn is_positive(n: INT) -> BOOL {
  RETURN n > 0
}

EXPORT is_valid_number
EXPORT is_positive
```

**Step 3: main.fl - 메인 프로그램**

```freelang
// main.fl
IMPORT add FROM "math"
IMPORT subtract FROM "math"
IMPORT multiply FROM "math"
IMPORT divide FROM "math"
IMPORT is_valid_number FROM "validation"
IMPORT is_positive FROM "validation"

FN run_calculator() {
  PRINTLN("계산기에 오신 걸 환영합니다!")

  LET result = add(10, 5)
  PRINTLN("10 + 5 = " + result)

  LET result = multiply(4, 3)
  PRINTLN("4 * 3 = " + result)

  IF is_positive(result):
    PRINTLN("결과는 양수입니다")
}

run_calculator()
```

**실행 결과:**
```
계산기에 오신 걸 환영합니다!
10 + 5 = 15
4 * 3 = 12
결과는 양수입니다
```

### 전문가라면?

복잡한 의존성을 가진 모듈 시스템을 설계합시다.

**프로젝트 구조:**
```
bank-system/
├── core/
│   ├── types.fl      (타입 정의)
│   └── errors.fl     (에러 처리)
├── domain/
│   ├── account.fl    (계좌 로직)
│   ├── transaction.fl (거래 로직)
│   └── user.fl       (사용자 로직)
├── utils/
│   ├── validation.fl (검증)
│   ├── crypto.fl     (암호화)
│   └── logger.fl     (로깅)
└── main.fl           (진입점)
```

**core/types.fl - 핵심 타입**

```freelang
// core/types.fl
RECORD User {
  id: INT
  name: STRING
  email: STRING
}

RECORD Account {
  id: INT
  user_id: INT
  balance: FLOAT
}

RECORD Transaction {
  id: INT
  from_account: INT
  to_account: INT
  amount: FLOAT
  timestamp: STRING
}

EXPORT User
EXPORT Account
EXPORT Transaction
```

**domain/account.fl - 계좌 도메인**

```freelang
// domain/account.fl
IMPORT Account FROM "core/types"
IMPORT validate_amount FROM "utils/validation"

FN create_account(user_id: INT) -> Account {
  RETURN {
    id: 1,
    user_id: user_id,
    balance: 0.0
  }
}

FN deposit(account: Account, amount: FLOAT) -> Account {
  IF NOT validate_amount(amount):
    RETURN account

  LET new_balance = account.balance + amount
  RETURN {
    id: account.id,
    user_id: account.user_id,
    balance: new_balance
  }
}

FN withdraw(account: Account, amount: FLOAT) -> Account {
  IF NOT validate_amount(amount):
    RETURN account

  IF amount > account.balance:
    RETURN account  // 잔액 부족

  LET new_balance = account.balance - amount
  RETURN {
    id: account.id,
    user_id: account.user_id,
    balance: new_balance
  }
}

EXPORT create_account
EXPORT deposit
EXPORT withdraw
```

**domain/transaction.fl - 거래 도메인**

```freelang
// domain/transaction.fl
IMPORT Transaction FROM "core/types"
IMPORT Account FROM "core/types"
IMPORT deposit FROM "domain/account"
IMPORT withdraw FROM "domain/account"

FN transfer(
  from_account: Account,
  to_account: Account,
  amount: FLOAT
) -> BOOL {
  // 출금 시도
  LET from_after = withdraw(from_account, amount)
  IF from_after.balance = from_account.balance:
    RETURN FALSE  // 출금 실패

  // 입금 처리
  LET to_after = deposit(to_account, amount)
  RETURN TRUE
}

FN create_transaction(
  from_id: INT,
  to_id: INT,
  amount: FLOAT
) -> Transaction {
  RETURN {
    id: 1,
    from_account: from_id,
    to_account: to_id,
    amount: amount,
    timestamp: "2026-03-26"
  }
}

EXPORT transfer
EXPORT create_transaction
```

**main.fl - 통합**

```freelang
// main.fl
IMPORT User FROM "core/types"
IMPORT Account FROM "core/types"
IMPORT create_account FROM "domain/account"
IMPORT deposit FROM "domain/account"
IMPORT withdraw FROM "domain/account"
IMPORT transfer FROM "domain/transaction"

FN main() {
  // 사용자 생성
  LET user = { id: 1, name: "Kim", email: "kim@freelang.dev" }

  // 계좌 생성 및 입금
  LET account1 = create_account(user.id)
  LET account1 = deposit(account1, 1000.0)

  LET account2 = create_account(2)

  // 송금
  LET success = transfer(account1, account2, 500.0)

  IF success:
    PRINTLN("송금 성공!")
    PRINTLN("계좌1 잔액: " + account1.balance)
    PRINTLN("계좌2 잔액: " + account2.balance)
}

main()
```

---

## 모듈 시스템의 고급 기능

### 1. 네임스페이스 (Namespace)

같은 이름의 함수가 여러 모듈에 있을 때:

```freelang
IMPORT add FROM "math"        // math.add
IMPORT add FROM "bigint"      // bigint.add (충돌!)

// 해결책: 별칭(alias) 사용
IMPORT add AS math_add FROM "math"
IMPORT add AS bigint_add FROM "bigint"

LET result1 = math_add(10, 5)
LET result2 = bigint_add(10000000, 5000000)
```

### 2. 선택적 Export (부분 공개)

```freelang
// internal.fl
FN public_function() { ... }
FN internal_function() { ... }

EXPORT public_function
// internal_function은 export 안 함 = 비공개
```

### 3. 재export (중계)

```freelang
// utils.fl
IMPORT add FROM "math"
IMPORT subtract FROM "math"

EXPORT add
EXPORT subtract
// utils 사용자는 utils에서 바로 import 가능
```

```freelang
// main.fl
IMPORT add FROM "utils"  // math를 직접 import하지 않아도 됨
```

---

## 모듈 사용 시 Best Practices

### ✅ 하는 것

```freelang
// ✅ Good: 명확한 책임
// user.fl
EXPORT create_user
EXPORT update_user
EXPORT delete_user

// ✅ Good: 의존성 명시
// order.fl
IMPORT create_user FROM "user"
IMPORT validate_email FROM "validation"
```

### ❌ 하지 말 것

```freelang
// ❌ Bad: 모든 걸 export
EXPORT *

// ❌ Bad: 순환 의존성
// a.fl imports b.fl
// b.fl imports a.fl (무한 루프!)

// ❌ Bad: 깊은 중첩
IMPORT create_user FROM "a/b/c/d/e/f"
```

---

## 성능과 조직

### 모듈 로딩

FreeLang의 모듈 시스템은 **지연 로딩(Lazy Loading)**을 지원합니다:

```freelang
// heavy_module.fl (큰 파일)
fn expensive_calculation() { ... }

// main.fl
// heavy_module을 사용하지 않으면 로드되지 않음
// 필요할 때만 로드 = 빠른 시작 시간 ⚡
```

### 순환 의존성 피하기

```
┌─────────┐
│ user.fl │
└────┬────┘
     │ imports
     ↓
┌──────────────┐
│ account.fl   │
└────┬─────────┘
     │ imports
     ↓
┌──────────────┐
│ transaction  │
└──────────────┘

// 단방향 의존성 = 깔끔한 구조 ✅
```

---

## 실제 사례: FreeLang의 표준 라이브러리

FreeLang 자체도 모듈로 구성되어 있습니다:

```freelang
// stdlib 구조
freelang/
├── core/
│   ├── array.fl     (배열 함수)
│   ├── string.fl    (문자열 함수)
│   └── math.fl      (수학 함수)
├── io/
│   ├── file.fl      (파일 I/O)
│   └── network.fl   (네트워크)
└── utils/
    ├── crypto.fl    (암호화)
    └── time.fl      (시간)
```

**사용 예:**

```freelang
IMPORT len FROM "freelang/core/array"
IMPORT read_file FROM "freelang/io/file"
IMPORT hash FROM "freelang/utils/crypto"

LET arr = [1, 2, 3, 4, 5]
LET size = len(arr)  // 5

LET content = read_file("data.txt")
LET checksum = hash(content)
```

---

## 다음 단계

모듈 시스템을 마스터했으니, 더 고급 주제를 배울 수 있습니다.

**다음에 읽을 거리:**
1. [FreeLang 비동기 모듈 - async 함수를 모듈화하기](#)
2. [패키지 관리 시스템 - 모듈을 배포하고 공유하기](#)
3. [의존성 주입(DI) 패턴 - 느슨한 결합 유지하기](#)

**직접 해보고 싶으신가요?**

```bash
# FreeLang 모듈 예제 실행
git clone https://github.com/freelang/freelang
cd freelang/examples/modules
freelang main.fl
```

---

**참고 자료:**
- [Module Pattern in Software Design - Martin Fowler](https://martinfowler.com/)
- [Clean Architecture - Robert C. Martin](https://www.oreilly.com/library/view/clean-architecture/9780134494272/)
- FreeLang 공식 문서: Module System https://freelang-docs.example.com/modules

---

이 글이 도움이 되었다면? 👍

- 다른 개발자와 공유해주세요
- 모듈화 경험담을 댓글로 나눠주세요
- "이런 모듈 패턴을 배우고 싶어요" 제안 환영합니다!

**Happy Modularizing! 🚀**
