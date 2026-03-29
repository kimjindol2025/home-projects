# 🔬 완벽한 분신 시스템 - 검증 프로토콜

**목적**: 3가지 훅(Stop/pre-commit/post-commit)의 실제 동작 확인
**기간**: 2026-03-28 ~ 2026-04-04 (1주)
**담당**: Bigwash-nim

---

## 📋 검증 체크리스트

### Phase 1️⃣: Stop 훅 검증 (이번 세션)

#### ✅ Step 1: 로깅 추가 (완료)
- [x] `self_improve.sh`에 훅 실행 로깅 추가
- [x] 로그 파일: `~/.claude/.hook-logs/self-improve.log`
- [x] 비동기 모드에서도 추적 가능

#### ⏳ Step 2: 세션 종료 시 Stop 훅 실행
**다음 행동**: 이번 세션 종료 시 `exit` 입력

```bash
exit  # 또는 Ctrl+C 여러 번
```

**예상 결과**:
- [ ] 터미널 종료 전 메시지: "🧠 Self-Improving: Analyzing session patterns..."
- [ ] 세션 종료 후 로그 파일 생성

#### ⏳ Step 3: 로그 확인 (다음 세션)
```bash
cat ~/.claude/.hook-logs/self-improve.log
```

**체크 항목**:
- [ ] 타임스탬프: 세션 종료 시간과 일치?
- [ ] Status: "COMPLETED"?
- [ ] 오류 메시지: "⚠ analyze failed" 또는 "⚠ update-rules failed" 없음?

#### ⏳ Step 4: .clauderules 변경 확인 (다음 세션)
```bash
tail -20 ~/.claude/.clauderules  # 또는 ~/.claude/CLAUDE.md
```

**체크 항목**:
- [ ] 새로운 규칙이 추가되었나?
- [ ] 추가된 규칙에 타임스탬프가 있나? (예: `<!-- 2026-03-28 22:XX -->`)
- [ ] 중복 규칙은 없나?

---

### Phase 2️⃣: Pre-commit 훅 검증 (다음 세션, 2026-03-29)

#### ⏳ Step 1: 테스트 파일 생성
```bash
cd ~/projects/freelang-evolving-compiler

# 1-1. 외부 라이브러리를 import하는 파일 생성
cat > test_external.go << 'EOF'
package main

import "github.com/some/external"

func main() {
    external.DoSomething()
}
EOF

# 1-2. 스테이징
git add test_external.go
```

**예상 결과**:
- [ ] `git commit` 명령 실행
- [ ] pre-commit 훅 실행: "sandbox.sh" 시작
- [ ] ❌ 커밋 차단 메시지: "External dependency detected"
- [ ] ViolationReport 생성: `~/.claude/missions/mission2-zero-dep-sandbox/violation_report.txt`

#### ⏳ Step 2: 통과 테스트
```bash
# 외부 라이브러리 제거
git reset HEAD test_external.go
rm test_external.go

# 정상 파일로 다시 커밋
echo "// Clean code" >> README.md
git add README.md
git commit -m "test: verify pre-commit pass"
```

**예상 결과**:
- [ ] pre-commit 훅 통과
- [ ] 커밋 성공
- [ ] post-commit 훅 자동 실행 (다음 Step 참고)

---

### Phase 3️⃣: Post-commit 훅 검증 (다음 세션, 2026-03-29)

#### ⏳ Step 1: GOGS wiki 확인
```bash
# 로컬에서
cat ~/.claude/gogs-pulse-history.json | tail -1
```

**체크 항목**:
- [ ] 새 커밋 레코드 추가?
- [ ] FNV-1a 해시 생성?
- [ ] 타임스탬프 기록?

#### ⏳ Step 2: GOGS 저장소 확인 (웹 브라우저)
```
https://gogs.dclub.kr/kim/freelang-missions/wiki
```

**체크 항목**:
- [ ] Wiki 페이지 생성됨?
- [ ] 커밋 메시지가 wiki에 반영됨?
- [ ] 날짜/시간 메타데이터 포함?

---

## 📊 검증 결과 기록 양식

각 단계 완료 후 아래 형식으로 기록:

```markdown
### [SESSION-DATE] Phase [N] Step [N]

**상태**: ✅ PASS / ⚠️ PARTIAL / ❌ FAIL

**실행 시간**: [HH:MM]
**환경**: [env details]

**결과**:
- [Check 1]: ✅
- [Check 2]: ✅
- ...

**로그**:
[relevant log output]

**다음 단계**: [Action item]
```

---

## 🎯 성공 기준

### Phase 1 (Stop 훅) 성공
```
✅ 로그 파일 생성됨
✅ 타임스탬프 기록됨
✅ .clauderules 변경됨
```

### Phase 2 (Pre-commit 훅) 성공
```
✅ 외부 라이브러리 감지 → 커밋 차단
✅ 정상 코드 → 커밋 통과
✅ ViolationReport 생성
```

### Phase 3 (Post-commit 훅) 성공
```
✅ GOGS wiki 페이지 자동 생성
✅ gogs-pulse-history.json 기록 누적
✅ 메타데이터 저장
```

---

## 💾 검증 로그 위치

| 항목 | 위치 |
|------|------|
| Stop 훅 로그 | `~/.claude/.hook-logs/self-improve.log` |
| Pre-commit 리포트 | `~/.claude/missions/mission2-zero-dep-sandbox/violation_report.txt` |
| Post-commit 기록 | `~/.claude/gogs-pulse-history.json` |
| GOGS wiki | `https://gogs.dclub.kr/kim/freelang-missions/wiki` |

---

## ⚠️ 예상 문제 & 대응

### 문제 1: Stop 훅이 실행 안 됨
**원인**: `async: true` 설정으로 로그 확인 어려움
**대응**: 로그 파일 직접 확인 (`~/.claude/.hook-logs/self-improve.log`)

### 문제 2: Pre-commit 훅이 작동 안 함
**원인**: Symlink 경로 오류 또는 권한 문제
**대응**:
```bash
ls -la ~/projects/freelang-evolving-compiler/.git/hooks/pre-commit
chmod +x ~/.claude/missions/mission2-zero-dep-sandbox/sandbox.sh
```

### 문제 3: Post-commit 훅 실행 안 됨
**원인**: GOGS API 인증 토큰 만료
**대응**: `~/.git-credentials` 다시 확인

### 문제 4: .clauderules 변경 안 됨
**원인**: `history.jsonl`에 패턴이 없음
**대응**: 다음 세션에서 실제 오류 발생 후 재시도

---

## 📅 검증 일정

| 날짜 | 작업 | 상태 |
|------|------|------|
| 2026-03-28 | Stop 훅 로깅 추가 | ✅ |
| 2026-03-29 | Phase 1 확인 (로그 검증) | ⏳ |
| 2026-03-29 | Phase 2 테스트 (pre-commit) | ⏳ |
| 2026-03-29 | Phase 3 테스트 (post-commit) | ⏳ |
| 2026-04-04 | 1주 후 성능 평가 | ⏳ |

---

## 🎓 검증 완료 후 다음 단계

✅ 모든 훅 정상 작동 확인 후:

1. **Morning Report 구현** (1시간)
   - SessionStart 훅 추가
   - 어제 학습 내용 브리핑

2. **프로덕션 배포** (선택)
   - 다른 프로젝트에도 훅 설치

3. **성능 최적화** (선택)
   - 훅 실행 시간 단축
   - 메모리 사용량 최소화

---

**검증 시작: 2026-03-28 22:25 KST**
**목표: 2026-04-04 완료**

"기록이 증명이다" — 이 검증 프로토콜 자체가 증명 과정입니다.
