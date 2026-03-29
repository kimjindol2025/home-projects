# FreeLang v2 자원 가이드

**업데이트**: 2026-03-29
**기준**: FreeLang v2
**상태**: 활성 참고 자료

---

## 📍 주요 위치

| 경로 | 설명 | 상태 |
|------|------|------|
| `v2-archive/freelang-v2/` | Go 구현 (레거시) | 📦 보관 |
| `v2-archive/fv2-lang/` | Rust 구현 | 📦 보관 |
| `projects/fv2-lang-go/` | 현재 활성 Go 프로젝트 | 🌟 활성 |

---

## 📚 자원 유형

### 1. 언어 사양 (Specification)
- 문서: `v2-archive/freelang-v2/API.md`
- 내용: 기본 데이터 타입, 함수, 채널 API
- 크기: ~22KB

### 2. 코드 구현
- **Go**: `v2-archive/freelang-v2/` (7.2GB)
  - 렉서, 파서, 컴파일러, 런타임
  - git 히스토리 포함

- **Rust**: `v2-archive/fv2-lang/` (15MB)
  - Rust 기반 구현 참고

### 3. 문서화
```
v2-archive/freelang-v2/
├─ API.md                          # API 레퍼런스
├─ API_REFERENCE.md               # 상세 레퍼런스
├─ CHANGELOG.md                   # 버전 변경사항
├─ CONTRIBUTING.md                # 기여 가이드
└─ [기타 40+ 문서]
```

### 4. 테스트 & 검증
- 단위 테스트: `v2-archive/freelang-v2/tests/`
- 통합 테스트: `v2-archive/freelang-v2/integration_tests/`
- 문서: `*_TEST_REPORT.md` (20+ 리포트)

---

## 🔍 검색 방법

### API 확인
```bash
# v2 API 문서 읽기
cat v2-archive/freelang-v2/API.md

# Go 코드 구현 검색
grep -r "func.*String" v2-archive/freelang-v2/runtime/
```

### 특정 기능 찾기
```bash
# 채널 관련
grep -r "Channel" v2-archive/freelang-v2/ --include="*.go"

# 타입 시스템
grep -r "TypeKind" v2-archive/freelang-v2/ --include="*.go"

# 함수 정의
grep -r "FuncDef" v2-archive/freelang-v2/ --include="*.go"
```

---

## 🛠️ 사용 흐름

### 1단계: 기능 부족 확인
```
"이 기능이 v2에 있나?"
→ FREELANG_V2_RESOURCES.md 확인
→ v2-archive/ 검색
```

### 2단계: 구현 보고
```
"부족한 기능: [기능명]"
→ 리포트 작성 및 제시
→ 승인 후 FreeLang v2로 구현
```

### 3단계: 코드 작성
```
v2-archive/freelang-v2/ 참고
→ FreeLang v2로 구현
→ 테스트 작성
```

---

## 📖 주요 문서

| 문서 | 용도 |
|------|------|
| `API.md` | API 명세 |
| `API_REFERENCE.md` | 함수/타입 레퍼런스 |
| `C-BINDING-INTEGRATION-GUIDE.md` | C 바인딩 |
| `CHANGELOG.md` | 버전 히스토리 |
| `*_REPORT.md` | 기능 검증 보고서 |

---

## ⚙️ 자주 쓰는 파일

```
v2-archive/freelang-v2/
├─ runtime/           # 런타임 구현
├─ lexer/            # 렉서
├─ parser/           # 파서
├─ compiler/         # 컴파일러
├─ vm/               # 가상 머신
└─ builtin/          # 기본 라이브러리
```

---

## 📝 부족한 기능 보고 템플릿

```markdown
## [기능명]

**현재 상태**: v2에 없음 / 부분 구현됨 / 완전 구현됨

**설명**:
- 목적: [무엇을 하는가]
- 예시: [사용 예시]

**v2 위치**: (있다면)
- 파일: `v2-archive/freelang-v2/[path]`
- 함수: `[함수명]`

**구현 우선순위**: High / Medium / Low
```

---

**참고**: 이 문서는 FreeLang v2 자원 접근을 위한 가이드입니다.
