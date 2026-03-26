# ❓ fv2-lang-go가 뭔데? FreeLang? Go? 뭐니?

**작성일**: 2026-03-26
**목적**: fv2-lang-go 정체 명확히 하기

---

## 🎯 정답

### fv2-lang-go는?

```
📦 fv2-lang-go

= FreeLang 프로그래밍 언어의 Go 버전

구성:
├─ "fv2" = FreeLang V2 의미
├─ "lang" = 프로그래밍 언어
└─ "go" = Go 언어로 구현됨
```

---

## 📚 이해하기

### 1단계: FreeLang이 뭔가?

**FreeLang**: 새로운 프로그래밍 언어
```
Python처럼 배우기 쉽고
Go처럼 빠른 언어

예시 코드:
fn add(a, b) {
    return a + b
}

fn main() {
    let result = add(5, 3)
    println(result)  // 8
}
```

**특징**:
- ✅ 간단한 문법
- ✅ 타입 안전
- ✅ 비동기 지원
- ✅ 100+ 내장 함수

---

### 2단계: V2는 뭔가?

**FreeLang V2**: 버전 2

```
Timeline:
V1 (초기)
  ↓
V2 (개선)  ← 우리가 지금 보는 것
  ↓
V3 (최신?)
```

**V2의 종류**:
1. freelang-v2 (큰 프로젝트, JavaScript/Go/Python 섞임)
2. fv2-lang (Rust 실험)
3. fv2-lang-go (Go 구현, 메인) ⭐

---

### 3단계: fv2-lang-go는 뭔가?

**fv2-lang-go**: FreeLang V2를 Go 언어로 구현한 것

```
┌─────────────────────────────┐
│  FreeLang 프로그래밍 언어    │
│  (개념, 문법, 규칙)          │
└──────────────┬──────────────┘
               │
        ┌──────┴──────┐
        │             │
        ↓             ↓
   JavaScript     Go (Go 언어)
   구현 버전      구현 버전
                   ↑
            fv2-lang-go
```

---

## 🔍 fv2-lang-go 상세 분석

### 뭘 하는가?

**FreeLang 코드를 실행하는 컴파일러**

```
입력:
hello.fl 파일
│
│ fn main() {
│     println("Hello World")
│ }
│

처리:
└─ fv2-lang-go 컴파일러 (Go로 만들어짐)
   ├─ Lexer: 코드 → 토큰
   ├─ Parser: 토큰 → AST
   ├─ Compiler: AST → 바이트코드
   └─ Runtime: 바이트코드 실행

출력:
실행 결과: "Hello World"
```

---

### 뭐로 만들어졌나?

**Go 언어**

```
fv2-lang-go/
├── internal/
│   ├── lexer/      (Go 코드)
│   ├── parser/     (Go 코드)
│   ├── codegen/    (Go 코드)
│   ├── stdlib/     (Go 코드)
│   └── typechecker/ (Go 코드)
└── cmd/fv2/main.go (Go 코드)

모든 파일: *.go (Go 소스 파일)
```

**의존성**:
```
외부:
├─ go-sqlite3 (데이터베이스)
└─ testify (테스트용)

표준:
├─ fmt, strings, time
├─ crypto/aes, crypto/sha256
├─ net/http, encoding/json
└─ ... (Go 표준 라이브러리)
```

---

### 역할은?

**FreeLang의 공식 구현**

| 항목 | 설명 |
|------|------|
| **언어** | FreeLang (새 언어) |
| **구현** | Go (Go 언어로 작성) |
| **역할** | 컴파일러 + 런타임 |
| **상태** | 100% 완성 ✅ |
| **품질** | 프로덕션 준비 완료 |

---

## 🗺️ V2 프로젝트들 비교

```
목표: FreeLang 언어를 여러 방식으로 구현

┌─────────────────────────────────────────┐
│         FreeLang 프로그래밍 언어         │
│              (개념/규칙)                 │
└────────────┬────────────┬───────────────┘
             │            │
             ↓            ↓
        ┌────────┐   ┌──────────┐
        │JavaScript│   │Go       │
        │구현      │   │구현     │
        │(freelang │   │(fv2-lang│
        │-v2)      │   │-go)     │
        │85%완성   │   │100%완성 │
        └────────┘   │✅메인   │
                     └──────────┘

추가:
├─ Rust 실험 (fv2-lang, 75% 완성)
└─ AI 연구 (v2-freelang-ai, 50% 완성)
```

---

## 💡 간단히 말하면

### 비유

```
음악 작곡가가 곡을 만들었다
(= FreeLang 언어 설계)

그 곡을 여러 악기로 연주했다
(= 여러 언어로 구현)

├─ 피아노 연주 (= JavaScript 구현)
├─ 바이올린 연주 (= Go 구현) ⭐ fv2-lang-go
└─ 기타 연주 (= Rust 구현)

가장 좋은 연주는?
→ 바이올린 (Go)! 🎻
  가장 깔끔하고 빠르다
```

---

## 📊 정리

| 항목 | 내용 |
|------|------|
| **fv2-lang-go는** | Go 언어로 구현한 FreeLang 컴파일러 |
| **FreeLang은** | 새로운 프로그래밍 언어 (개념) |
| **Go는** | fv2-lang-go를 만드는 데 사용된 언어 |
| **완성도** | 100% ✅ |
| **의존성** | 최소 (1개 + 표준) |
| **상태** | 프로덕션 준비 완료 |

---

## 🎯 최종 답변

### Q: fv2-lang-go가 뭔데?

### A:
```
FreeLang이라는 새 프로그래밍 언어를
Go 언어로 구현한 컴파일러/런타임

= FreeLang을 실행하는 엔진 (Go로 만들어짐)

예시:
사용자가 쓸 코드: hello.fl (FreeLang)
처리 프로그램: fv2-lang-go (Go로 만들어짐)
실행: FreeLang 코드 실행
```

### Q: 왜 Go로 만들었나?

### A:
```
특징:
1. 빠르다 (컴파일 언어)
2. 간단하다 (배우기 쉬움)
3. 포팅 쉽다 (Windows, Mac, Linux)
4. 의존성 적다 (표준 라이브러리 충실)

다른 구현:
- JavaScript (freelang-v2): 느리지만 웹 브라우저에서 실행
- Rust (fv2-lang): 빠르지만 구현 안 함
- Go (fv2-lang-go): 빠르고 완성도 높음 ⭐
```

---

**정리일**: 2026-03-26
**상태**: ✅ 명확함
