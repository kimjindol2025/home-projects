---
name: GOGS Repository Quick Start
description: GOGS 레포 신규 생성 빠른 시작 가이드
type: template
---

# ⚡ GOGS 레포 생성 빠른 시작 (5분)

## 📋 체크리스트 (순서대로 따르기)

### 1️⃣ GOGS에서 레포 생성 (2분)
```
http://gogs.local → + → New Repository
├─ Repository Name: freelang-v4
├─ Description: FreeLang v4 프로젝트
├─ Private: (선택)
├─ Initialize with: README
└─ License: MIT
```

### 2️⃣ 로컬 설정 (1분)
```bash
cd ~/.projects/core/freelang-v4
git remote add origin http://gogs.local/[user]/freelang-v4.git
```

### 3️⃣ 필수 파일 확인 (1분)
```
✓ README.md          (GOGS가 자동 생성)
✓ CLAUDE.md          (필수)
✓ .gitignore         (필수)
✓ package.json       (필수)
```

### 4️⃣ 초기 커밋 (1분)
```bash
git add -A
git commit -m "Initial commit: Project setup

Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>"
git push -u origin master
```

---

## 🔗 GOGS 레포 이름 규칙

```
✅ 올바른 예시:
- freelang-v4              (코어)
- module-compiler          (모듈)
- tool-gogs-api            (도구)
- experiment-ml            (실험)

❌ 잘못된 예시:
- Freelang_V4              (대문자)
- freelang v4              (공백)
- freelang-v4-backup-final (너무 길음)
```

---

## 📁 명령어 요약

```bash
# 1. 레포 클론
git clone http://gogs.local/[user]/freelang-v4.git ~/.projects/core/freelang-v4

# 2. 원격 추가 (기존 프로젝트)
git remote add origin http://gogs.local/[user]/freelang-v4.git

# 3. 푸시 (로컬 → GOGS)
git push -u origin master

# 4. 풀 (GOGS → 로컬)
git pull origin master

# 5. 상태 확인
git status
git log --oneline
```

---

## 💾 필수 커밋 형식

```
feat: 새 기능
fix: 버그 수정
docs: 문서
refactor: 코드 정리

[마지막 줄 필수]
Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
```

---

## 📚 추가 참고

📖 전체 가이드: `cat ~/.prompts/workflows/gogs-repo-creation.md`

