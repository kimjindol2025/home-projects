# FreeLang 전체 버전 자원 맵

**업데이트**: 2026-03-29
**커버**: v2, v3, v4, v5, v6
**상태**: 전체 자원 인덱싱 완료

---

## 📍 버전별 위치

| 버전 | 경로 | 언어 | 상태 | 용도 |
|------|------|------|------|------|
| **v2** | `v2-archive/freelang-v2/` | Go | 📦 보관 | 참고/레거시 |
| **v2.5** | `dev/lang/versions/freelang-v2-5/` | - | 📦 보관 | 버전 참고 |
| **v3** | `dev/archived/v3-freelang-ai/` | - | 📦 보관 | AI 통합 참고 |
| **v4** | `dev/lang/freelang-v4/` | Go | 📦 보관 | 구현 참고 |
| **v5** | `.projects/tools/freelang-review/freelang-v5/` | - | 📦 보관 | 리뷰/검증 |
| **v6** | `dev/lang/versions/freelang-v6-source/` | - | 📦 보관 | 최신 구현 |
| **v6** | `dev/archived/freelang-v6/` | - | 📦 보관 | 아카이브 |

---

## 📚 버전별 특징

### v2 (Go 구현)
**위치**:
- `v2-archive/freelang-v2/` (7.2GB, 완전한 git 히스토리)
- `dev/archived/freelang-v2/` (백업)

**특징**:
- 기본 언어 구조 정의
- 렉서, 파서, 컴파일러, 런타임 포함
- 40+ 문서 (API, 레퍼런스 등)

**주요 파일**:
```
├─ lexer/              # 토큰화
├─ parser/             # 문법 분석
├─ compiler/           # 컴파일
├─ vm/                 # 가상머신
├─ runtime/            # 런타임
├─ builtin/            # 기본 함수
└─ API.md              # API 문서
```

**검색**:
```bash
cd v2-archive/freelang-v2
git log --oneline
grep -r "func.*String" runtime/
```

---

### v2.5 (마이너 버전)
**위치**: `dev/lang/versions/freelang-v2-5/`

**특징**:
- v2의 개선 버전
- 버전별 변경사항 추적 가능

---

### v3 (AI 통합)
**위치**: `dev/archived/v3-freelang-ai/`

**특징**:
- AI 기능 실험
- v2 기반 확장

---

### v4 (확장 버전)
**위치**: `dev/lang/freelang-v4/`

**특징**:
- 여러 변형 구현
  - `freelang-v4-core` (핵심)
  - `freelang-v4-stdlib` (표준 라이브러리)
  - `freelang-v4-sqlite-integration` (DB)
  - `freelang-v4-http` (네트워크)
  - `freelang-v4-concurrency` (동시성)
  - `freelang-v4-crypto` (암호화)
  - `freelang-v4-jit` (JIT 컴파일)
  - `freelang-v4-orm` (ORM)
  - 기타 15+ 모듈

**주요 모듈**:
```bash
dev/archived/
├─ freelang-v4-core/
├─ freelang-v4-stdlib/
├─ freelang-v4-sqlite-integration/
├─ freelang-v4-http/
├─ freelang-v4-concurrency/
├─ freelang-v4-crypto/
├─ freelang-v4-jit/
├─ freelang-v4-orm/
├─ freelang-v4-compiler-optimizer/
└─ [더 많은 모듈...]
```

---

### v5 (검증/리뷰)
**위치**:
- `.projects/tools/freelang-review/freelang-v5/`
- `.projects/freelang-v5-ai/`

**특징**:
- v4 검증 및 개선
- 리뷰 중심

---

### v6 (최신)
**위치**:
- `dev/lang/versions/freelang-v6-source/` (소스)
- `dev/archived/freelang-v6/` (아카이브)
- `dev/archived/freelang-v6-ai-sovereign/` (AI 버전)

**특징**:
- 가장 최신 구현
- AI 기능 강화
- 성능 최적화

---

## 🔍 검색 방법

### 1. 기능별 검색
```bash
# 채널 구현 찾기
grep -r "Channel" dev/archived/freelang-v4-concurrency/

# 타입 시스템
grep -r "TypeKind\|TypeInfo" v2-archive/freelang-v2/

# HTTP 기능
grep -r "HandleHTTP\|Request" dev/archived/freelang-v4-http/
```

### 2. 버전 비교
```bash
# v2 vs v4 변경사항
diff -r v2-archive/freelang-v2/compiler/ dev/archived/freelang-v4-core/
```

### 3. 최신 코드 찾기
```bash
# v6 최신 구현
ls -lh dev/lang/versions/freelang-v6-source/
```

---

## 🛠️ 기능별 추천 위치

| 기능 | 버전 | 위치 |
|------|------|------|
| 기본 문법 | v2 | `v2-archive/freelang-v2/parser/` |
| 런타임 | v2 | `v2-archive/freelang-v2/runtime/` |
| 타입 시스템 | v4 | `dev/archived/freelang-v4-core/` |
| 동시성 | v4 | `dev/archived/freelang-v4-concurrency/` |
| 데이터베이스 | v4 | `dev/archived/freelang-v4-sqlite-integration/` |
| 네트워킹 | v4 | `dev/archived/freelang-v4-http/` |
| JIT 컴파일 | v4 | `dev/archived/freelang-v4-jit/` |
| 암호화 | v4 | `dev/archived/freelang-v4-crypto/` |
| 최신 구현 | v6 | `dev/lang/versions/freelang-v6-source/` |

---

## 📋 부족한 기능 보고 프로세스

### 1단계: 검색
```
"이 기능이 어느 버전에 있나?"
→ 위 테이블에서 버전 확인
→ 해당 위치 검색
```

### 2단계: 분석
```bash
# 예: 채널 기능 검색
grep -r "Channel" v2-archive/freelang-v2/
grep -r "Channel" dev/archived/freelang-v4-concurrency/

# 결과: v2에는 없고, v4에 있음 → v4 참고
```

### 3단계: 보고 (부족할 경우)
```markdown
## 기능명: [기능]

**부재 위치**: v2, v3, v4, v5
**발견 위치**: v6 (일부) 또는 없음

**설명**: [무엇을 하는가]

**권장 구현**: FreeLang v2 기반 (기본) 또는 v4 참고
```

---

## 🎯 추천 학습 순서

1. **v2** - 기본 구조 이해
   ```bash
   cd v2-archive/freelang-v2
   cat API.md
   ```

2. **v4** - 확장 기능 학습
   ```bash
   ls -la dev/archived/freelang-v4-*/
   ```

3. **v6** - 최신 구현 확인
   ```bash
   cd dev/lang/versions/freelang-v6-source/
   ```

---

## 📊 버전 크기 비교

```
v2-archive/freelang-v2/      : 7.2GB (완전한 git 히스토리)
dev/lang/freelang-v4/        : ~500MB
dev/lang/versions/freelang-v6/ : ~300MB
```

---

## ⚠️ 주의사항

- **읽기 전용** 권장 (참고용)
- 수정 필요시 백업 후 진행
- 대용량 파일 주의 (v2: 7.2GB)

---

**마지막 업데이트**: 2026-03-29
**목적**: FreeLang 전체 버전 통일된 자원 접근
