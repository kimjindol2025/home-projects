# GOGS API 학습 가이드

**작성일**: 2026-03-04
**버전**: 1.0
**목적**: GOGS 저장소 관리 및 자동화 가이드

---

## 📋 목차

1. [GOGS 기본 개념](#gogs-기본-개념)
2. [저장소 생성](#저장소-생성)
3. [Git 기본 작업](#git-기본-작업)
4. [GOGS API 활용](#gogs-api-활용)
5. [실습 사례](#실습-사례)
6. [팁 & 트러블슈팅](#팁--트러블슈팅)

---

## GOGS 기본 개념

### GOGS란?
**GOGS** (Go Git Service)는 자체 호스팅 Git 서비스입니다.
- 가볍고 빠른 Git 저장소 관리
- Go 언어로 작성됨
- GitHub 같은 인터페이스 제공

### 주요 특징
- ✅ 저장소 생성 & 관리
- ✅ 사용자 계정 관리
- ✅ 사용 권한 제어
- ✅ Webhook 지원
- ✅ REST API 제공

### GOGS 서버 정보
```
서버: https://gogs.dclub.kr
사용자명: kim
저장소 경로: /kim/<repository-name>.git
```

---

## 저장소 생성

### 1️⃣ GOGS 웹 인터페이스에서 생성

**절차:**
```
1. GOGS 서버 접속 → https://gogs.dclub.kr
2. 로그인 (사용자명: kim)
3. 우측 상단 "+" 버튼 → "New Repository"
4. 저장소 이름 입력 (예: freelang-module-system)
5. 설명(Description) 입력 (선택사항)
6. "Create Repository" 클릭
```

**생성 후 URL:**
```
https://gogs.dclub.kr/kim/<repository-name>.git
```

### 2️⃣ 저장소명 명명 규칙

본 프로젝트에서 사용하는 명명 규칙:
```
패턴: freelang-<system-name>

예시:
- freelang-async-system
- freelang-closure-system
- freelang-module-system
- freelang-iterator-system
- freelang-distributed-system
```

---

## Git 기본 작업

### 저장소 초기화

```bash
# 새 디렉토리 생성
mkdir -p ~/freelang-<system-name>
cd ~/freelang-<system-name>

# git 초기화
git init

# git 사용자 설정 (첫 사용시만)
git config user.email "claude@anthropic.com"
git config user.name "Claude Haiku 4.5"
```

### 파일 추가 & 커밋

```bash
# 모든 파일 추가
git add -A

# 상태 확인
git status

# 커밋 생성
git commit -m "feat: 기능 설명"

# 로그 확인
git log --oneline
```

### 원격 저장소 설정 & 푸시

```bash
# 원격 저장소 추가
git remote add origin https://gogs.dclub.kr/kim/<repository-name>.git

# 원격 저장소 확인
git remote -v

# 푸시 (첫 번째는 -u 옵션 권장)
git push -u origin master
```

### 전체 워크플로우 예시

```bash
cd ~/freelang-module-system

# 1. 모든 파일 준비 완료
git add -A
git status

# 2. 커밋 생성
git commit -m "feat: 🎯 Module System Phase 9 Option E - 1,600줄"

# 3. 원격 설정
git remote add origin https://gogs.dclub.kr/kim/freelang-module-system.git

# 4. 푸시
git push -u origin master

# 5. 확인
git log --oneline
git remote -v
```

---

## GOGS API 활용

### REST API 기본 정보

**기본 URL:**
```
https://gogs.dclub.kr/api/v1
```

**인증:**
```
- Token-based: Authorization: token <access_token>
- Basic Auth: -u username:password
```

### 주요 API 엔드포인트

#### 1️⃣ 저장소 목록 조회

```bash
curl -H "Authorization: token <TOKEN>" \
  https://gogs.dclub.kr/api/v1/user/repos
```

**응답 예시:**
```json
[
  {
    "id": 1,
    "name": "freelang-async-system",
    "full_name": "kim/freelang-async-system",
    "description": "Async/Await System",
    "url": "https://gogs.dclub.kr/kim/freelang-async-system"
  }
]
```

#### 2️⃣ 저장소 생성 (API)

```bash
curl -X POST \
  -H "Authorization: token <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "freelang-new-system",
    "description": "New System Description",
    "private": false
  }' \
  https://gogs.dclub.kr/api/v1/user/repos
```

#### 3️⃣ 저장소 정보 조회

```bash
curl -H "Authorization: token <TOKEN>" \
  https://gogs.dclub.kr/api/v1/repos/kim/<repo-name>
```

#### 4️⃣ 저장소 삭제

```bash
curl -X DELETE \
  -H "Authorization: token <TOKEN>" \
  https://gogs.dclub.kr/api/v1/repos/kim/<repo-name>
```

### Access Token 얻기

**GOGS 웹 인터페이스에서:**
```
1. 로그인
2. Settings (우측 상단 아바타 → Settings)
3. Applications
4. Generate New Token
5. Token 이름 입력 (예: "API Token")
6. 생성된 Token 복사 (다시 보기 불가)
```

---

## 실습 사례

### 사례 1: Phase 9 Async System 푸시

**상황:**
- 로컬에서 async_definition.fl, await_executor.fl 등 구현 완료
- GOGS에 freelang-async-system 저장소 생성됨

**작업:**
```bash
cd ~/freelang-async-system

# 모든 파일 스테이징
git add -A

# 커밋 생성
git commit -m "feat: 🚀 Async/Await System Phase 9 Option D - 1,600줄"

# 원격 설정
git remote add origin https://gogs.dclub.kr/kim/freelang-async-system.git

# 푸시
git push -u origin master

# 확인
git log --oneline
git remote -v
```

**결과:**
```
✅ 커밋: d0dd471
✅ GOGS 저장소에 1,735 insertions
✅ 5개 파일 업로드 (src/, docs/)
```

### 사례 2: 여러 저장소 관리

**상황:**
- 5개의 독립적인 프로젝트가 각각 GOGS 저장소에 있음
- Option A, B, C, D, E 모두 완료

**저장소 목록:**
```
1. freelang-lifetime-analyzer     ✅ GOGS 푸시 완료
2. freelang-iterator-system       ✅ GOGS 푸시 완료
3. freelang-closure-system        ⏳ 로컬 준비 완료
4. freelang-async-system          ✅ GOGS 푸시 완료
5. freelang-module-system         ⏳ 로컬 준비 완료
```

**배치 작업 스크립트:**
```bash
#!/bin/bash

REPOS=(
  "freelang-lifetime-analyzer"
  "freelang-iterator-system"
  "freelang-closure-system"
  "freelang-async-system"
  "freelang-module-system"
)

for repo in "${REPOS[@]}"; do
  if [ -d "~/$repo" ]; then
    cd ~/$repo
    git status
    echo "---"
  fi
done
```

---

## 팁 & 트러블슈팅

### ✅ 유용한 팁

#### 1️⃣ .gitignore 설정
```bash
# 불필요한 파일 제외
echo "node_modules/" > .gitignore
echo ".DS_Store" >> .gitignore
git add .gitignore
```

#### 2️⃣ 커밋 메시지 규칙
```
feat: 새 기능
fix: 버그 수정
docs: 문서 업데이트
refactor: 코드 리팩토링
test: 테스트 추가
chore: 기타 작업

예: feat: Add async/await support
```

#### 3️⃣ 여러 파일 한 번에 커밋
```bash
git add -A              # 모든 변경 추가
git commit -m "message" # 커밋
git push origin master  # 푸시
```

#### 4️⃣ 이전 커밋 수정
```bash
# 마지막 커밋 메시지 수정
git commit --amend -m "new message"

# 마지막 커밋에 파일 추가
git add forgotten_file.fl
git commit --amend --no-edit
```

### ⚠️ 트러블슈팅

#### 문제 1: "repository not found"
```
원인: GOGS에 저장소가 없음
해결: GOGS 웹에서 저장소를 먼저 생성
```

#### 문제 2: "Permission denied (publickey)"
```
원인: SSH 키 설정 미흡
해결: HTTPS URL 사용 (https://gogs.dclub.kr/...)
```

#### 문제 3: "fatal: remote origin already exists"
```
원인: git remote add 중복
해결: git remote set-url origin <new-url>
```

#### 문제 4: 로컬 변경사항 손실
```
원인: force push로 히스토리 덮어씀
해결: git reflog로 이전 커밋 복구
```

### 🔍 유용한 명령어

```bash
# 원격 저장소 상태 확인
git remote -v
git remote show origin

# 커밋 히스토리 보기
git log --oneline
git log --oneline --graph --all

# 특정 파일 변경 추적
git log -p -- src/async_definition.fl

# 변경 사항 비교
git diff                    # 커밋되지 않은 변경
git diff --cached           # 스테이징된 변경
git diff HEAD~1 HEAD        # 마지막 커밋 vs 이전

# 현재 상태 확인
git status
git status -s               # 짧은 형식

# 브랜치 관리
git branch                  # 로컬 브랜치 목록
git branch -r               # 원격 브랜치 목록
```

---

## 📚 참고 자료

### 공식 문서
- GOGS 공식 사이트: https://gogs.io
- GOGS API 문서: https://gogs.io/docs/advanced/api
- Git 공식 문서: https://git-scm.com/doc

### 자주 사용하는 명령어 모음

**초기 설정:**
```bash
git config user.email "claude@anthropic.com"
git config user.name "Claude Haiku 4.5"
```

**기본 푸시:**
```bash
git add -A
git commit -m "message"
git push origin master
```

**저장소 클론:**
```bash
git clone https://gogs.dclub.kr/kim/<repo-name>.git
```

**원격 업데이트:**
```bash
git pull origin master
```

---

## 🎯 Phase 9 GOGS 푸시 상황판

| 옵션 | 프로젝트 | 상태 | GOGS URL |
|------|---------|------|----------|
| A | Lifetime Analyzer | ✅ 푸시됨 | freelang-lifetime-analyzer |
| B | Iterator System | ✅ 푸시됨 | freelang-iterator-system |
| C | Closure/Lambda | ⏳ 로컬 준비 | freelang-closure-system |
| D | Async/Await | ✅ 푸시됨 | freelang-async-system |
| E | Module System | ⏳ 로컬 준비 | freelang-module-system |

---

**이 가이드는 실무 경험을 바탕으로 작성되었습니다.**
**마지막 업데이트: 2026-03-04**

