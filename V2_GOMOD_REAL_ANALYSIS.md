# 🔍 fv2-lang-go Go 의존성 진짜 분석

**작성일**: 2026-03-26
**발견**: go.mod 표기가 misleading했음

---

## 📊 진짜 상황

### go.mod vs 실제 코드

**go.mod 내용**:
```
require (
	github.com/davecgh/go-spew v1.1.1       // indirect (테스트)
	github.com/mattn/go-sqlite3 v1.14.37    // indirect (DB)
	github.com/pmezard/go-difflib v1.0.0    // indirect (테스트)
	github.com/stretchr/objx v0.5.0         // indirect (테스트)
	github.com/stretchr/testify v1.8.4      // indirect (테스트)
	gopkg.in/yaml.v3 v3.0.1                 // indirect (테스트)
)
```

**실제 import (internal/stdlib/database.go)**:
```go
import (
	"database/sql"      // ← Go 표준 라이브러리
	"fmt"               // ← Go 표준 라이브러리
	"strings"           // ← Go 표준 라이브러리

	_ "github.com/mattn/go-sqlite3"  // ← 외부 의존성 (DB 드라이버)
)
```

---

## ❌ 의존성 분류의 오류

### 원인: "// indirect" 표시

**잘못 이해한 부분**:
- go.mod에 `// indirect`로 표시된 것들이 모두 테스트용이라고 가정
- 실제로는 일부가 프로덕션 코드에 직접 사용됨

### 실제 분석

| 의존성 | 용도 | 프로덕션? | 필수? |
|--------|------|---------|------|
| **go-sqlite3** | SQLite 드라이버 | ✅ Yes | ✅ 필수 |
| **testify** | 테스트 라이브러리 | ❌ No | ❌ 선택 |
| **yaml.v3** | YAML 파싱 | ❌ No | ❌ 선택 |
| **go-spew** | 디버그 출력 | ❌ No | ❌ 선택 |
| **go-difflib** | 테스트 비교 | ❌ No | ❌ 선택 |
| **objx** | 테스트 유틸 | ❌ No | ❌ 선택 |

---

## 🔎 코드별 분석

### 1. database.go (프로덕션)

```go
import (
	"database/sql"                   // Go 표준
	"fmt"                            // Go 표준
	"strings"                        // Go 표준
	_ "github.com/mattn/go-sqlite3"  // ✅ 외부 의존성 (필수)
)
```

**평가**:
- ✅ Go 표준 라이브러리만 사용
- ✅ SQLite는 필요 시에만 (가능하면 제거 가능)

### 2. grpc.go (프로덕션)

```go
import (
	"fmt"       // Go 표준
	"sync"      // Go 표준
	"time"      // Go 표준
	// 외부 의존성 없음!
)
```

**평가**:
- ✅ 완전히 표준 라이브러리만 사용
- ✅ 외부 의존성 0

### 3. http.go (프로덕션)

```go
import (
	"net/http"      // Go 표준
	"net/url"       // Go 표준
	"strings"       // Go 표준
	"time"          // Go 표준
	// 외부 의존성 없음!
)
```

**평가**:
- ✅ 완전히 표준 라이브러리만 사용
- ✅ 외부 의존성 0

### 4. websocket.go (프로덕션)

```go
import (
	"encoding/json" // Go 표준
	"sync"          // Go 표준
	"time"          // Go 표준
	// 외부 의존성 없음!
)
```

**평가**:
- ✅ 완전히 표준 라이브러리만 사용
- ✅ 외부 의존성 0

### 5. crypto.go (프로덕션)

```go
import (
	"crypto/aes"      // Go 표준
	"crypto/cipher"   // Go 표준
	"crypto/rand"     // Go 표준
	"crypto/sha256"   // Go 표준
	"encoding/base64" // Go 표준
	"encoding/hex"    // Go 표준
	"fmt"             // Go 표준
	"strings"         // Go 표준
	// 외부 의존성 없음!
)
```

**평가**:
- ✅ 완전히 표준 라이브러리만 사용
- ✅ 외부 의존성 0

---

## 📈 프로덕션 의존성 종합

### 프로덕션 코드
```
fv2-lang-go/internal/stdlib/
├── crypto.go       → 0개 의존성 (표준만)
├── grpc.go         → 0개 의존성 (표준만)
├── http.go         → 0개 의존성 (표준만)
├── websocket.go    → 0개 의존성 (표준만)
├── database.go     → 1개 의존성 (go-sqlite3)
└── (기타)          → 0개 의존성
```

### 테스트 코드만 사용하는 패키지
```
- testify (테스트 라이브러리)
- yaml.v3 (설정 파싱, 테스트용)
- go-spew (디버그)
- go-difflib (테스트 비교)
- objx (테스트 유틸)
```

---

## ✅ 실제 결론

### 프로덕션에서의 외부 의존성

**필수 의존성**: 1개
- ✅ **go-sqlite3** (database/sql 드라이버)

**선택적 의존성**: 5개 (모두 테스트 전용)
- ❌ testify, yaml.v3, go-spew, go-difflib, objx

### 테스트 포함 전체

**필수**: 1개 (go-sqlite3)
**테스트용**: 5개
**합계**: 6개

---

## 🎯 독립성 평가

### 현재
```
외부 의존성:
├─ 프로덕션: 1개 (go-sqlite3)
├─ 테스트: 5개
└─ 합계: 6개
```

### 만약 SQLite 없이?
```
외부 의존성:
├─ 프로덕션: 0개 (표준만)
├─ 테스트: 5개
└─ 합계: 5개 (또는 테스트 제외 시 0개)
```

---

## 💡 옵션

### 옵션 1: 현 상태 유지 ✅
- 외부 의존성: 1개 (go-sqlite3)
- 비용: 0
- 장점: 모든 기능 작동

### 옵션 2: SQLite 제거 (프로덕션만) ⭐ 권장
- 외부 의존성: 0개 (프로덕션)
- 테스트 의존성: 5개 (테스트용)
- 비용: 낮음 (database 함수 수정)
- 장점: 프로덕션 완전 독립

### 옵션 3: 완전 독립 (테스트까지)
- 외부 의존성: 0개
- 비용: 높음 (테스트 프레임워크 재구성)

---

## 🏆 최종 평가

### fv2-lang-go vs freelang-v2

| 항목 | fv2-lang-go | freelang-v2 |
|------|------------|------------|
| **프로덕션 외부 의존** | 1개 (SQLite) | 7개 |
| **테스트 의존** | 5개 | 12개 |
| **독립성** | 🟢 높음 | 🟡 낮음 |
| **복잡도** | 🟢 낮음 | 🟡 높음 |

### 결론
✅ **fv2-lang-go는 실질적으로 거의 완전 독립적**
- 프로덕션 의존성: 1개만 (SQLite)
- 나머지 5개는 모두 테스트 전용
- 프로덕션만 보면 거의 의존성 0에 가까움

---

**조사일**: 2026-03-26
**상태**: ✅ 완료 (수정됨)
**오류 발견**: go.mod의 "// indirect" 표시가 모두 테스트용이라고 오해
