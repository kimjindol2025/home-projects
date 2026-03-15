---
name: GOGS Repository Rules Index
description: GOGS 레포 신규 생성 규칙 전체 인덱스
type: reference
version: 1.0
---

# 📦 GOGS 레포지토리 신규 생성 규칙

**작성**: 2026-03-15
**버전**: 1.0
**상태**: ✅ 활성

---

## 🎯 규칙 개요

GOGS에 새로운 레포지토리를 생성할 때 준수해야 할 모든 규칙을 정의합니다.

---

## 📚 규칙 문서

### 1. 📖 전체 가이드
**파일**: `~/.prompts/workflows/gogs-repo-creation.md`
**내용**:
- 레포 이름 규칙 (명명 규칙, 카테고리 프리픽스)
- 5단계 생성 프로세스
- 필수 파일 템플릿 (.gitignore, README.md, CLAUDE.md)
- Git 워크플로우
- 체크리스트 & 주의사항

### 2. ⚡ 빠른 시작
**파일**: `~/.prompts/templates/gogs-quick-start.md`
**내용**:
- 5분 내 GOGS 레포 생성
- 체크리스트
- 명령어 요약
- 필수 커밋 형식

### 3. 🔧 자동화 스크립트
**파일**: `~/.prompts/tasks/coding/gogs-repo-auto-setup.sh`
**내용**:
- 프로젝트 디렉토리 자동 생성
- 필수 파일 자동 생성
- Git 초기화 & 커밋 자동화
- 한 줄 명령으로 완료

---

## 🎯 핵심 규칙 (한눈에)

### 레포 이름 규칙
```
{category}-{project-name}[-v{version}]

예시:
✅ freelang-v4               (core project)
✅ module-compiler           (module)
✅ tool-gogs-api             (tool)
✅ experiment-neural-network (experiment)
✅ archive-v2-backup         (archive)
```

### 5단계 프로세스
```
1️⃣  GOGS 웹에서 레포 생성 (2분)
2️⃣  로컬에 클론 또는 remote 추가 (1분)
3️⃣  필수 파일 확인 (1분)
4️⃣  초기 커밋 & 푸시 (1분)
5️⃣  Claude 메모리 초기화 (1분)
```

### 필수 파일
```
✓ README.md          GOGS 생성 또는 수동
✓ CLAUDE.md          AI 작업 가이드
✓ .gitignore         표준 탈제
✓ package.json       프로젝트 메타
✓ MEMORY.md          진행 상황 기록
```

### 커밋 규칙
```
feat:    새 기능
fix:     버그 수정
docs:    문서
refactor: 코드 정리
test:    테스트

[필수]
Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
```

---

## 🚀 사용 방법

### 옵션 1: 자동화 (권장)
```bash
# 모든 것이 자동으로 생성됨
bash ~/.prompts/tasks/coding/gogs-repo-auto-setup.sh freelang-v4 core freelang

# 그 후:
cd ~/.projects/core/freelang-v4
git push -u origin master
```

### 옵션 2: 빠른 시작 (5분)
```bash
cat ~/.prompts/templates/gogs-quick-start.md
# → 체크리스트 따라하기
```

### 옵션 3: 전체 가이드 (상세)
```bash
cat ~/.prompts/workflows/gogs-repo-creation.md
# → 모든 규칙 상세히 읽기
```

---

## 📋 체크리스트 (신규 레포 생성 시)

### GOGS 설정
- [ ] 레포 이름 규칙 확인 (소문자, 하이픈)
- [ ] 설명 입력
- [ ] 라이선스: MIT
- [ ] .gitignore 추가
- [ ] README 생성

### 로컬 설정
- [ ] 디렉토리: ~/.projects/[category]/[name]/
- [ ] CLAUDE.md 작성
- [ ] Git remote 설정
- [ ] 필수 파일 존재 확인

### 첫 커밋
- [ ] 모든 파일 git add
- [ ] 표준 커밋 메시지
- [ ] Co-Authored-By 포함
- [ ] git push -u origin master

### Claude 메모리
- [ ] MEMORY.md 생성
- [ ] GOGS URL 기록
- [ ] 프로젝트 상태 입력
- [ ] 다음 액션 작성

---

## 🔗 카테고리별 프리픽스

```
freelang-*        → FreeLang 핵심 프로젝트
                    (freelang-v4, freelang-final 등)

module-*          → 언어 모듈 & 컴파일러
                    (module-compiler, module-vm 등)

tool-*            → 도구 & 유틸
                    (tool-gogs-api, tool-benchmark 등)

experiment-*      → 실험 & 학습 프로젝트
                    (experiment-neural-network 등)

archive-*         → 완료된 프로젝트
                    (archive-v2, archive-phase5 등)

integration-*     → 통합 프로젝트
                    (integration-ci-cd 등)
```

---

## 🚨 주의사항

### ❌ 하지 말 것
- 커밋 메시지에 Co-Authored-By 빠뜨리기
- 대문자 사용 (Freelang-V4 ❌)
- 공백 사용 (freelang v4 ❌)
- README.md 없이 푸시
- .env 파일 커밋
- MEMORY.md 미업데이트

### ✅ 해야 할 것
- 정기적 pull/push (보안 & 동기화)
- MEMORY.md 매 세션 업데이트
- 표준 커밋 형식 준수
- 테스트 후 커밋
- 브랜치 전략 유지 (master, develop, feature/*)

---

## 📞 자주 묻는 질문

### Q1: GOGS 레포는 몇 개까지 생성 가능?
A: 제한 없음. 카테고리별로 정렬하여 관리.

### Q2: 기존 프로젝트를 GOGS에 추가하려면?
A: `git remote add origin [gogs-url]` 후 `git push -u origin master`

### Q3: MEMORY.md는 어디에 저장?
A: `.claude/projects/[project-name]/memory/MEMORY.md`

### Q4: 커밋 서명이 왜 필수?
A: AI 기여 추적 & Claude 메모리 시스템 통합

### Q5: 브랜치는 어떻게 관리?
A: master (상용), develop (개발), feature/* (기능)

---

## 📊 예시 구조

```
GOGS 레포 (http://gogs.local/freelang/freelang-v4)
        ↓
로컬 저장소 (~/.projects/core/freelang-v4)
        ├─ .git/              (GOGS 연동)
        ├─ src/               (소스 코드)
        ├─ tests/             (테스트)
        ├─ docs/              (문서)
        ├─ .claude/
        │  └─ projects/freelang-v4/
        │     └─ memory/
        │        └─ MEMORY.md
        ├─ README.md
        ├─ CLAUDE.md
        ├─ .gitignore
        └─ package.json
```

---

## 🔗 관련 문서

- **프로젝트 정리**: ~/.projects/PROJECT-INDEX.md
- **시스템 보고서**: ~/.SYSTEM-COMPLETE-REPORT.md
- **Claude 메모리**: ~/.claude/projects/*/memory/MEMORY.md
- **프롬프트 저장소**: ~/.prompts/

---

## 📝 버전 히스토리

| 버전 | 날짜 | 변경사항 |
|------|------|---------|
| 1.0 | 2026-03-15 | 초기 버전 |

---

**Status**: ✅ Active & Ready
**Next**: GOGS 레포 생성 시 이 규칙 확인

