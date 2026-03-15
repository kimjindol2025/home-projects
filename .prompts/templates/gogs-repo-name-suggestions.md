---
name: GOGS Repository Name Suggestions
description: FreeLang 프로젝트에 추천되는 GOGS 저장소 이름 모음
type: template
---

# 💡 GOGS 저장소 이름 추천

**기준**: 규칙에 맞고, 명확하며, 검색하기 쉬운 이름

---

## 🎯 카테고리별 추천명

### 1️⃣ Core Projects (핵심 프로젝트)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **freelang-v5** | FreeLang v5 메인 (차기 버전) | ⭐⭐⭐⭐⭐ |
| **freelang-core** | 핵심 라이브러리 | ⭐⭐⭐⭐ |
| **freelang-sdk** | 공식 SDK | ⭐⭐⭐⭐ |
| **freelang-stdlib** | 표준 라이브러리 | ⭐⭐⭐⭐ |

### 2️⃣ Modules (언어 모듈)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **module-compiler** | 컴파일러 | ⭐⭐⭐⭐⭐ |
| **module-vm** | 가상머신 | ⭐⭐⭐⭐⭐ |
| **module-runtime** | 런타임 | ⭐⭐⭐⭐ |
| **module-optimizer** | 최적화 패스 | ⭐⭐⭐⭐ |
| **module-allocator** | 메모리 관리 | ⭐⭐⭐ |

### 3️⃣ Tools (도구)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **tool-cli** | 명령행 도구 | ⭐⭐⭐⭐⭐ |
| **tool-package-manager** | 패키지 관리자 | ⭐⭐⭐⭐ |
| **tool-debugger** | 디버거 | ⭐⭐⭐⭐ |
| **tool-profiler** | 프로파일러 | ⭐⭐⭐⭐ |
| **tool-formatter** | 코드 포매터 | ⭐⭐⭐ |

### 4️⃣ Infrastructure (인프라)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **infra-gogs-api** | GOGS 통합 API | ⭐⭐⭐⭐ |
| **infra-ci-cd** | CI/CD 파이프라인 | ⭐⭐⭐⭐ |
| **infra-docker** | Docker 환경 | ⭐⭐⭐ |
| **infra-kubernetes** | K8s 배포 | ⭐⭐⭐ |

### 5️⃣ Libraries (라이브러리)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **lib-json** | JSON 처리 | ⭐⭐⭐ |
| **lib-crypto** | 암호화 | ⭐⭐⭐ |
| **lib-async** | 비동기 처리 | ⭐⭐⭐ |
| **lib-http** | HTTP 클라이언트 | ⭐⭐⭐ |

### 6️⃣ Experiments (실험 프로젝트)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **exp-ai-optimization** | AI 기반 최적화 | ⭐⭐⭐ |
| **exp-quantum-compute** | 양자 컴퓨팅 | ⭐⭐⭐ |
| **exp-parallel-runtime** | 병렬 런타임 | ⭐⭐⭐ |
| **exp-llvm-backend** | LLVM 백엔드 | ⭐⭐⭐ |

### 7️⃣ Archive (아카이브)

| 추천명 | 설명 | 선택도 |
|--------|------|--------|
| **archive-v4-final** | v4 최종 버전 | ⭐⭐⭐ |
| **archive-v3-stable** | v3 안정 버전 | ⭐⭐ |
| **archive-legacy** | 레거시 코드 | ⭐⭐ |

---

## 🎯 추천 우선순위 (지금 시작할 프로젝트)

### 🥇 최고 우선순위
```
1. freelang-v5           (차기 메인 프로젝트)
2. module-compiler       (컴파일러 개선)
3. tool-cli              (사용자 도구)
```

### 🥈 2순위
```
4. module-vm             (런타임 개선)
5. infra-ci-cd           (빌드 자동화)
6. freelang-sdk          (개발자 키트)
```

### 🥉 3순위
```
7. tool-package-manager  (패키지 관리)
8. module-optimizer      (성능 최적화)
9. exp-ai-optimization   (실험)
```

---

## 📋 선택 기준

### ✅ 좋은 이름
```
✓ 짧고 명확 (freelang-v5)
✓ 목적이 분명 (module-compiler)
✓ 검색 용이 (tool-cli)
✓ 카테고리 명시 (freelang-*, module-*)
```

### ❌ 피해야 할 이름
```
✗ 너무 긺 (freelang-v5-main-compiler-final)
✗ 모호함 (project-something)
✗ 대문자 (Freelang-V5)
✗ 공백 (freelang v5)
✗ 언더스코어 (freelang_v5)
```

---

## 🚀 추천 선택

### 상황별 추천

**새 프로젝트 시작**
```
freelang-v5            (차기 메인)
module-compiler        (개선 모듈)
tool-cli               (사용자 도구)
```

**기존 프로젝트 마이그레이션**
```
archive-v4-final       (v4 보관)
archive-v3-stable      (v3 보관)
```

**실험/학습**
```
exp-ai-optimization    (AI 실험)
exp-quantum-compute    (양자 실험)
```

---

## 💾 생성 명령어 (자동화)

```bash
# 추천명으로 레포 자동 생성
bash ~/.prompts/tasks/coding/gogs-repo-auto-setup.sh freelang-v5 core freelang
cd ~/.projects/core/freelang-v5
git push -u origin master
```

---

## 🎖️ 최종 추천

### 🏆 TOP 3 추천명
1. **freelang-v5** ⭐⭐⭐⭐⭐ (메인 프로젝트)
2. **module-compiler** ⭐⭐⭐⭐⭐ (컴파일러)
3. **tool-cli** ⭐⭐⭐⭐⭐ (사용자 도구)

---

**Generated**: 2026-03-15
**Style**: 규칙 준수, 명확함, 검색성

