# 🔗 V2 프로젝트 외부 의존성 조사

**작성일**: 2026-03-26
**목적**: V2 프로젝트 4개의 외부 의존성 분석

---

## 📊 의존성 요약

| 프로젝트 | 언어 | 의존성 개수 | 상태 | 평가 |
|---------|------|-----------|------|------|
| **fv2-lang-go** | Go 1.21 | 6개 (indirect) | ⚠️ 최소 | ⭐⭐⭐⭐⭐ |
| **freelang-v2** | Node.js | 12개 (prod) + 8개 (dev) | ⚠️ 중간 | ⭐⭐⭐⭐ |
| **fv2-lang** | Rust 2021 | 0개 | ✅ 없음 | ⭐⭐⭐⭐⭐ |
| **v2-freelang-ai** | - | - | - | - |

---

## 1️⃣ fv2-lang-go (Go 1.21)

### 외부 의존성: 6개 (모두 indirect)

```
require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.37 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

### 분석

**의존성 목록**:
1. **go-spew** (v1.1.1) - 디버깅/덤프 유틸
2. **go-sqlite3** (v1.14.37) - SQLite 드라이버
3. **go-difflib** (v1.0.0) - 문자열 비교
4. **testify** (v1.8.4) - 테스트 라이브러리
5. **objx** (v0.5.0) - 맵 조작
6. **yaml.v3** (v3.0.1) - YAML 파서

**평가**: ✅ **매우 안전**
- 모두 indirect (직접 사용 아님, 테스트용)
- 표준 라이브러리에 가까운 필수 도구들
- 버전이 낮음 (안정적)

**리스크**: 🟢 **낮음**
- SQLite 의존성 1개만 있음
- 나머지는 순수 Go 패키지

---

## 2️⃣ freelang-v2 (Node.js)

### 프로덕션 의존성: 7개

```json
{
  "dependencies": {
    "better-sqlite3": "^11.10.0",      // 고성능 SQLite
    "chalk": "^4.1.2",                  // 터미널 색상
    "express": "^5.2.1",                // 웹 프레임워크
    "koffi": "^2.15.1",                 // FFI (C 바인딩)
    "spdy": "^4.0.2",                   // HTTP/2 프로토콜
    "vscode-languageserver": "^8.1.0",  // LSP 서버
    "vscode-languageserver-textdocument": "^1.0.12"  // LSP 문서
  }
}
```

### 개발 의존성: 8개

```json
{
  "devDependencies": {
    "@types/jest": "^29.5.8",
    "@types/node": "^20.10.0",
    "@typescript-eslint/eslint-plugin": "^6.13.0",
    "@typescript-eslint/parser": "^6.13.0",
    "eslint": "^8.55.0",
    "jest": "^29.7.0",
    "ts-jest": "^29.1.1",
    "ts-node": "^10.9.2",
    "typescript": "^5.3.3"
  }
}
```

### 분석

**프로덕션 의존성 평가**:

| 패키지 | 용도 | 리스크 | 평가 |
|--------|------|--------|------|
| **better-sqlite3** | 데이터베이스 | ⚠️ 중간 | 필수, C 바인딩 |
| **chalk** | 터미널 UI | 🟢 낮음 | 안전 |
| **express** | 웹 서버 | ⚠️ 중간 | 표준, 버전 5.x (최신) |
| **koffi** | FFI 바인딩 | 🔴 높음 | C 코드 실행 (보안 주의) |
| **spdy** | HTTP/2 | 🟢 낮음 | 프로토콜 구현 |
| **vscode-languageserver** | LSP 지원 | 🟢 낮음 | VS Code 표준 |
| **vscode-languageserver-textdocument** | LSP 문서 | 🟢 낮음 | VS Code 표준 |

**리스크**: 🟡 **중간**
- **koffi**: FFI (C 바인딩) - 외부 C 코드 실행 위험
- **better-sqlite3**: 네이티브 바인딩 필요
- express v5.x: 최신 버전 (안정성 아직 검증 중)

**개발 의존성 평가**:
- 모두 표준 개발 도구 (TypeScript, Jest, ESLint)
- 🟢 **안전** - 프로덕션에 포함 안 됨

---

## 3️⃣ fv2-lang (Rust 2021 Edition)

### 외부 의존성: **0개** ✅

```toml
[dependencies]
# 빈 상태
```

### 분석

**평가**: ⭐⭐⭐⭐⭐ **최고 점수**
- 완전히 자체 구현
- 외부 의존성 0개
- 포팅/배포 용이

**장점**:
- ✅ 완전히 독립적
- ✅ 빌드 간단
- ✅ 버전 호환성 문제 없음
- ✅ 보안 위험 낮음

---

## 4️⃣ v2-freelang-ai

### 외부 의존성: **미정**
- 프로젝트 초기 단계
- 의존성 파일 없음
- README만 존재

---

## 📈 종합 평가

### 의존성 비교

```
fv2-lang-go:
├─ 의존성 적음 (6개, 모두 indirect)
├─ 언어: Go (자체 포함)
└─ 리스크: 🟢 낮음

freelang-v2:
├─ 의존성 중간 (7개 prod)
├─ 위험 요소: koffi (FFI)
└─ 리스크: 🟡 중간

fv2-lang:
├─ 의존성 없음 (0개)
└─ 리스크: 🟢 가장 낮음
```

### 보안 우려사항

**🔴 높음**:
- freelang-v2의 **koffi** (C 코드 직접 실행)
  - 임의 C 함수 호출 가능
  - 메모리 안전성 보장 없음

**🟡 중간**:
- better-sqlite3 (네이티브 바인딩)
- express v5.x (최신, 아직 검증 중)

**🟢 낮음**:
- fv2-lang-go: 모두 표준 도구
- fv2-lang: 의존성 0

---

## 💡 권장사항

### fv2-lang-go (현재 메인)
✅ **그대로 유지**
- 의존성 최소
- 안정적

### freelang-v2 (레거시)
⚠️ **검토 필요**
- koffi 사용 여부 확인
- 필수시에만 사용
- 메모리 안전성 점검

### fv2-lang (Rust 실험)
✅ **권장**
- 의존성 0
- 가장 안전한 선택
- 완성 시 최고 품질

### 프로덕션 배포
**순서**:
1. ⭐ fv2-lang-go (지금 사용 가능)
2. 🏛️ freelang-v2 (레거시, FFI 검토)
3. 🔬 fv2-lang (완성 후)

---

**조사일**: 2026-03-26
**상태**: ✅ 완료
