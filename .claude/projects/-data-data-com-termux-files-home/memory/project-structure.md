---
name: 프로젝트 폴더 구조 (2026-03-16)
description: ~/projects 중앙 관리. 6개 활성 프로젝트. freelang-to-c, c-compiler-from-scratch 추가
type: reference
---

# 📁 프로젝트 폴더 구조

**최종 업데이트**: 2026-03-16

---

## 🗂️ **홈 디렉토리 구조**

```
~/
├── projects/                          # 🆕 프로젝트 중앙 관리 폴더
│   ├── c-compiler-from-scratch/      # 19,158줄 C 컴파일러 + VM
│   ├── freelang-to-c/                # 🆕 FreeLang → C 트랜스파일러
│   ├── freelang-library-extraction/  # FreeLang 라이브러리 추출
│   ├── golang_study/                 # Go 학습 프로젝트
│   ├── julia-compiler/               # Julia 컴파일러
│   ├── pulseai-lang/                 # PulseAI 언어
│   └── README.md
│
├── kim-project-cli/                   # 프로젝트 관리 시스템 (315개 프로젝트)
│   ├── .git/
│   ├── bin/
│   ├── memo/
│   ├── memo-config.json
│   ├── package.json
│   └── README.md
│
└── .projects/                         # 로컬 프로젝트 폴더들
    ├── archived/      (12개)
    ├── core/         (90개)
    ├── experiments/  (10개)
    └── modules/      (20개)
```

---

## 📊 **프로젝트 현황**

### **활성 프로젝트 (~/projects/)**

| 프로젝트 | 경로 | 상태 | 역할 |
|---------|------|------|------|
| **c-compiler-from-scratch** | `~/projects/c-compiler-from-scratch` | 🟡 진행중 | C 컴파일러 + VM |
| **freelang-to-c** | `~/projects/freelang-to-c` | 🟢 신규 | FreeLang → C 변환 |
| **freelang-library-extraction** | `~/projects/freelang-library-extraction` | 🟡 진행중 | 라이브러리 추출 |
| **golang_study** | `~/projects/golang_study` | ⏳ 학습용 | Go 언어 학습 |
| **julia-compiler** | `~/projects/julia-compiler` | 🟡 진행중 | Julia 컴파일러 |
| **pulseai-lang** | `~/projects/pulseai-lang` | 🟡 진행중 | PulseAI 언어 |

---

## 🔗 **GOGS 저장소 매핑**

| 프로젝트 | GOGS URL | 로컬 경로 |
|---------|----------|---------|
| c-compiler-from-scratch | https://gogs.dclub.kr/kim/c-compiler-from-scratch.git | `~/projects/c-compiler-from-scratch` |
| freelang-to-c | https://gogs.dclub.kr/kim/freelang-to-c.git | `~/projects/freelang-to-c` |
| kim-project-cli | https://gogs.dclub.kr/kim/kim-project-cli.git | `~/kim-project-cli` |

---

## 📈 **프로젝트 규모**

| 프로젝트 | 코드량 | 커밋 | 상태 |
|---------|--------|------|------|
| c-compiler-from-scratch | 19,158줄 | 100+ | ✅ 작동 |
| freelang-to-c | 178줄 (문서) | 1 | 🟢 시작 |
| kim-project-cli | 3,563줄 | 6 | ✅ 운영 중 |
| .projects 하위 | 132개 프로젝트 | - | 📦 추적 |

---

## 🎯 **작업 흐름**

### 개발 사이클

```
~/projects/freelang-to-c/
    ⬇️ (작업)
코드 수정 & 테스트
    ⬇️ (커밋)
git commit
    ⬇️ (푸시)
git push origin main
    ⬇️ (GOGS에 반영)
https://gogs.dclub.kr/kim/freelang-to-c.git
```

### 빌드 파이프라인

```
~/projects/freelang-to-c/
    (FreeLang 소스)
    ⬇️
~/projects/c-compiler-from-scratch/
    (C 컴파일러로 변환)
    ⬇️
ELF 바이너리
    (실행)
```

---

## 💾 **중요 위치**

### 메인 프로젝트 (지금 작업 중)
```
~/projects/freelang-to-c/          # ← 현재 활성 프로젝트
  ├── CLAUDE.md                     # 프로젝트 헌장
  ├── MEMORY.md                     # 진행 상황
  ├── src/codegen/                  # C 코드 생성 (구현 예정)
  ├── examples/                      # 테스트 예제
  └── tests/                         # 단위 테스트
```

### 의존 프로젝트
```
~/projects/c-compiler-from-scratch/
  ├── src/                           # C 컴파일러 소스
  ├── vm/                            # VM 구현
  └── a.out                          # 컴파일러 바이너리
```

### 프로젝트 관리 도구
```
~/kim-project-cli/
  ├── memo/                          # Claude 메모 저장소
  ├── gogs-lookup-api.js            # REST API
  └── bin/kim-lookup                # CLI 도구
```

---

## 🔄 **이동 기록 (2026-03-16)**

### 변경사항
```
Before:
  ~/freelang-to-c                   ❌ 삭제됨
  ~/kim-project-cli/c-compiler...   ❌ 옮겨짐

After:
  ~/projects/freelang-to-c/         ✅ 추가됨
  ~/projects/c-compiler-from-scratch/ ✅ 추가됨
  ~/kim-project-cli/                ✅ 유지됨
```

---

## 📝 **커맨드 빠른 참조**

```bash
# freelang-to-c 작업 디렉토리
cd ~/projects/freelang-to-c

# c-compiler-from-scratch 참조
cd ~/projects/c-compiler-from-scratch

# 프로젝트 관리 도구 사용
kim-lookup <repo>

# 프로젝트 목록 확인
ls -la ~/projects/
```

---

## 🚀 **다음 작업**

1. **freelang-to-c Phase 1 시작**
   ```bash
   cd ~/projects/freelang-to-c
   # 기본 타입 변환 구현
   ```

2. **c-compiler-from-scratch 참고**
   ```bash
   cd ~/projects/c-compiler-from-scratch
   # 생성된 C 코드 컴파일 테스트
   ```

3. **통합 파이프라인 구축**
   ```
   freelang-to-c → (생성) → C 코드
   C 코드 → c-compiler-from-scratch
   → ELF 바이너리 → 실행
   ```

---

**상태**: 📁 구조 정리 완료 | **준비 완료** ✅

