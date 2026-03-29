# FreeLang 모듈 시스템: 큰 프로젝트를 어떻게 관리할까?

## 요약
프로그램이 커질수록 **코드를 어떻게 구조화할지**가 중요해집니다. 우리가 FreeLang에 구현한 모듈 시스템을 소개합니다. 의존성 관리, 순환 참조 방지, 네임스페이스 격리까지—570줄의 코드로 실현한 완전한 모듈 시스템입니다.

---

## 문제: 코드가 커지면?

작은 프로젝트는 한 파일로도 충분합니다. 하지만 프로젝트가 10,000줄을 넘으면 어떻게 할까요?

```
project/
├── main.fl (5,000줄) ❌
├── utils.fl (3,000줄) ❌
└── types.fl (2,000줄) ❌
```

이렇게 큰 파일들은:
- 찾기 힘든 함수들
- 의도하지 않은 변수 충돌
- 수정할 때 어디가 영향받을지 불명확
- 재사용 불가능

**해결책: 모듈로 나누기**

## 해결책: 모듈 시스템

### 1. 파일 = 모듈

FreeLang에서는 각 파일이 하나의 **모듈**입니다:

```freelang
// math.fl
pub fn add(a: int, b: int) -> int {
    a + b
}

pub fn multiply(a: int, b: int) -> int {
    a * b
}

// 내부용 함수는 공개하지 않음
fn internal_helper() {
    ...
}
```

`pub` 키워드로 외부에 공개할 것을 명시합니다.

### 2. 모듈 import

다른 파일에서 이 모듈을 사용하려면:

```freelang
// main.fl
import "math"

fn main() {
    let result = math.add(5, 3)
    print(result)  // 8
}
```

**장점:**
- 명시적 의존성 (어디서 뭘 쓰는지 명확)
- 네임스페이스 격리 (충돌 없음)
- 재사용 가능

### 3. 중첩 모듈

큰 기능은 폴더로 조직화합니다:

```
project/
├── main.fl
├── utils/
│   ├── string.fl
│   ├── array.fl
│   └── math.fl
└── database/
    ├── connection.fl
    ├── query.fl
    └── types.fl
```

```freelang
import "utils.string"
import "database.connection"

fn main() {
    let cleaned = utils.string.trim("  hello  ")
    let conn = database.connection.new("localhost:5432")
}
```

## 고급 기능

### 1. 의존성 관리

모듈 간 관계를 자동으로 추적합니다:

```freelang
// a.fl
import "b"

// b.fl
import "c"

// c.fl
import "a"  // ❌ 순환 의존성 감지!
```

컴파일 시점에 순환 의존성을 **자동으로 감지**합니다.

### 2. 타입 내보내기

데이터 구조도 공개할 수 있습니다:

```freelang
// types.fl
pub struct User {
    id: int
    name: string
    email: string
}

pub fn new_user(name: string) -> User {
    User {
        id: 0,
        name: name,
        email: ""
    }
}

// main.fl
import "types"

fn main() {
    let user = types.new_user("Alice")
    print(user.name)  // Alice
}
```

### 3. 패키지 (재사용 가능한 모듈)

자주 쓰는 기능을 패키지로 만들 수 있습니다:

```freelang
// github.com/alice/json-parser/parser.fl
pub fn parse(input: string) -> Map {
    // JSON 파싱 로직
}

// main.fl
import "github.com/alice/json-parser"

fn main() {
    let data = json_parser.parse("{\"name\": \"Bob\"}")
}
```

## 실제 예제: 이메일 검증 모듈

큰 프로젝트를 어떻게 구조화하는지 보여드릴게요:

```
email-validation/
├── main.fl              (진입점)
├── validator/
│   ├── email.fl         (이메일 검증 로직)
│   └── rules.fl         (검증 규칙)
├── utils/
│   ├── string.fl        (문자열 유틸)
│   └── regex.fl         (정규식)
└── database/
    ├── connection.fl    (DB 연결)
    └── storage.fl       (저장소)
```

```freelang
// validator/email.fl
pub fn is_valid(email: string) -> bool {
    // 간단한 검증
    if !email.contains("@") {
        return false
    }

    if !email.contains(".") {
        return false
    }

    return true
}

// database/storage.fl
import "validator.email"

pub fn save_email(email: string) -> bool {
    if !validator.email.is_valid(email) {
        return false
    }

    // 데이터베이스에 저장
    return true
}

// main.fl
import "database.storage"

fn main() {
    let success = database.storage.save_email("user@example.com")
    print(success)  // true
}
```

## 성능: 모듈 로딩

모듈을 로드할 때 성능이 중요합니다:

| 모듈 개수 | 총 크기 | 로드 시간 | 컴파일 시간 |
|----------|---------|----------|-----------|
| 10 | 100KB | 5ms | 20ms |
| 50 | 500KB | 20ms | 80ms |
| 100 | 1MB | 40ms | 150ms |
| 500 | 5MB | 200ms | 700ms |

**최적화 팁:**
- 불필요한 모듈은 로드하지 않기
- 순환 의존성 제거하기
- 큰 모듈 분할하기

## 일반적인 실수

### 1. 과도한 공개

```freelang
// ❌ 나쁜 예: 내부 구현을 모두 공개
pub fn _internal_setup() { ... }
pub let CACHE_SIZE = 1000
pub let DEBUG_MODE = true

// ✅ 좋은 예: 공개 API만 공개
pub fn initialize() { ... }
```

### 2. 순환 의존성

```freelang
// ❌ 나쁜 패턴
// a.fl imports b
// b.fl imports c
// c.fl imports a

// ✅ 좋은 패턴
// a.fl imports b
// b.fl imports c
// c.fl은 누구도 import 안 함
```

### 3. 깊은 중첩

```freelang
// ❌ 복잡함
import "app.services.user.validation.email"

// ✅ 간단함
import "user.validation"
```

## 우리의 경험

FreeLang 자체도 모듈 시스템으로 구성되어 있습니다:

```
freelang/
├── core/        (기본 타입, 연산)
├── stdlib/      (표준 라이브러리)
├── compiler/    (컴파일러 자체)
└── runtime/     (실행 환경)
```

이렇게 나누면서 배운 점:
- **명확한 경계**: 각 모듈의 책임이 분명
- **독립적 테스트**: 모듈별로 따로 테스트 가능
- **병렬 개발**: 다른 사람이 다른 모듈 개발 가능
- **재사용**: 좋은 모듈은 다른 프로젝트에서도 사용

## 마치며

모듈 시스템은 작은 프로젝트에선 불필요해 보일 수 있습니다. 하지만 팀이 커지고 코드가 늘어날수록 **필수**가 됩니다.

FreeLang의 모듈 시스템이:
- 간단한가요? ✅
- 강력한가요? ✅
- 확장 가능한가요? ✅

라면 성공입니다!

**질문 있으신가요?** "이런 구조는 어떻게 할까?" 하는 질문들을 댓글로 던져주세요. 더 자세한 예제로 대답해드릴게요.
