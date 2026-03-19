# GOGS 푸시 가이드 📚

**작성일**: 2026-03-19
**대상**: fv2-lang-go 프로젝트
**상태**: ✅ 준비 완료

---

## 🔑 인증 정보

```
기관: GOGS (Git On Git Service)
서버: gogs.dclub.kr
계정: kim
토큰: 929b266ef2e69d9dcdc4f1b6cbf3141a068258b3
저장위치: ~/.netrc
```

---

## 🚀 빠른 시작 (3단계)

### Step 1: 변경사항 스테이징
```bash
cd ~/projects/fv2-lang-go
git add -A
```

### Step 2: 커밋
```bash
git commit -m "✨ 기능 설명"
```

### Step 3: GOGS에 푸시
```bash
git push gogs master
```

**끝!** ✅

---

## 📋 상세 워크플로우

### 1️⃣ 작업 디렉토리 이동
```bash
cd ~/projects/fv2-lang-go
```

### 2️⃣ 상태 확인
```bash
git status
```

**출력 예시**:
```
On branch master
Changes not staged for commit:
  modified:   internal/typechecker/checker.go

Untracked files:
  internal/typechecker/types.go
```

### 3️⃣ 변경사항 스테이징

**옵션 A: 모든 파일 추가**
```bash
git add -A
```

**옵션 B: 특정 파일만 추가**
```bash
git add internal/typechecker/checker.go
git add internal/typechecker/types.go
```

**옵션 C: 대화형 추가**
```bash
git add -i
```

### 4️⃣ 커밋

**기본 형식**:
```bash
git commit -m "메시지"
```

**좋은 커밋 메시지 예시**:
```bash
git commit -m "✨ Type Checker 구현 완료"
git commit -m "🐛 파서 오류 수정"
git commit -m "📚 README 업데이트"
```

**상세 메시지**:
```bash
git commit -m "✨ Phase 3.1: Type Checker 구현

- 16개 테스트 통과
- 20+ 타입 검사 규칙
- CLI 통합 완료"
```

### 5️⃣ GOGS로 푸시

**전용 저장소**:
```bash
git push gogs master
```

**메인 저장소**:
```bash
git push origin master
```

**양쪽 모두**:
```bash
git push gogs master && git push origin master
```

---

## 🌳 저장소 설정

### 현재 Remote 확인
```bash
git remote -v
```

### Remote 설명

| 이름 | 용도 | URL |
|------|------|-----|
| **gogs** | 👍 추천 (전용) | https://gogs.dclub.kr/kim/fv2-lang-go.git |
| **origin** | 메인 저장소 | https://gogs.dclub.kr/kim/projects.git |

---

## 📖 실제 사용 예시

### 예시 1: Phase 3.1 완료 후 푸시

```bash
cd ~/projects/fv2-lang-go
git status
git add -A
git commit -m "✨ Phase 3.1: Type Checker 구현 완료

- types.go: 280줄 (타입 시스템)
- checker.go: 430줄 (검사 엔진)
- 16개 테스트 100% 통과
- CLI 통합"

git push gogs master
git log --oneline -1
```

### 예시 2: 버그 수정 후 푸시

```bash
cd ~/projects/fv2-lang-go
git add internal/typechecker/checker.go
git commit -m "🐛 Type Checker: MatchExpression 타입 추론 수정"
git push gogs master
```

---

## 🔍 푸시 전 확인

### 커밋 로그 보기
```bash
git log --oneline -5
git log -1 --stat
git show HEAD
```

### 푸시될 커밋 확인
```bash
git log origin/master..HEAD
```

---

## 🚨 문제 해결

### ❌ 푸시 실패: "fatal: Authentication failed"

**해결**:
```bash
cat ~/.netrc                    # 파일 확인
ls -la ~/.netrc                 # 권한 확인
chmod 600 ~/.netrc              # 권한 수정
git push gogs master            # 다시 시도
```

### ❌ 커밋 실패: "no changes added"

**해결**:
```bash
git status
git add -A
git commit -m "메시지"
```

### ❌ 푸시 실패: "updates were rejected"

**해결**:
```bash
git pull gogs master            # 원격 먼저 가져오기
# 충돌이 있으면 해결
git push gogs master            # 다시 푸시
```

---

## 📊 유용한 명령어

### 상태 관련
```bash
git status              # 현재 상태
git status -s           # 간단한 표시
git diff                # 변경 내용 보기
```

### 스테이징 관련
```bash
git add .               # 현재 디렉토리 모든 파일
git add -A              # 전체 변경사항
git reset               # 스테이징 취소
```

### 커밋 관련
```bash
git commit -m "메시지"  # 커밋
git log --oneline       # 커밋 로그
git log -p              # 상세 로그
```

### 푸시 관련
```bash
git push gogs master                # 푸시
git push origin master              # 메인 저장소 푸시
git push gogs master && \
git push origin master              # 양쪽 모두
```

---

## 💾 체크리스트

**푸시하기 전에**:
- [ ] `cd ~/projects/fv2-lang-go` 이동
- [ ] `git status` 확인
- [ ] `git add -A` 스테이징
- [ ] `git commit -m "메시지"` 커밋
- [ ] `git log --oneline -1` 최종 확인

**푸시 후**:
- [ ] `git push gogs master` 실행
- [ ] GOGS 웹 확인: https://gogs.dclub.kr/kim/fv2-lang-go

---

## 🎯 Best Practices

### ✅ 좋은 습관

```bash
git commit -m "✨ 새 기능"       # 명확한 메시지
git commit -m "🐛 버그 수정"     # 이모지 사용
git push gogs master            # 자주 푸시 (매일 1-2회)
git add -A && git commit -m "..." && \
git push gogs master            # 파이프라인으로 연결
```

### ❌ 피해야 할 습관

```bash
git commit -m "수정"             # 애매한 메시지
git commit -m "작업"             # 무의미한 메시지
# 푸시하지 않고 여러 커밋 쌓기
git push --force                # 경고 없이 강제 푸시
```

---

## 🔗 유용한 링크

| 항목 | URL |
|------|-----|
| **GOGS 서버** | https://gogs.dclub.kr |
| **FV2 저장소** | https://gogs.dclub.kr/kim/fv2-lang-go |
| **Projects** | https://gogs.dclub.kr/kim/projects |

---

## 📝 커밋 메시지 규칙

```
타입(범위): 제목

본문

---

타입:
  ✨ feat     - 새 기능
  🐛 fix      - 버그 수정
  📚 docs     - 문서
  🎨 style    - 코드 스타일
  🔧 refactor - 리팩토링
  ⚡ perf     - 성능 개선
  ✅ test     - 테스트
  🚀 chore    - 기타

예시:
  ✨ feat(typechecker): Type Checker 구현
  🐛 fix(parser): IfExpression 파싱 오류 수정
  📚 docs: README 업데이트
```

---

**핵심**: 현재 모든 인증이 완료되어 있으므로, `git push gogs master`만 실행하면 됩니다! 🚀
