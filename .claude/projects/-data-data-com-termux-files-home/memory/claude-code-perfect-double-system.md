---
name: Claude Code 완벽한 분신 시스템 (4단계 실전 가이드)
description: 기록이 증명이다 철학 기반, 커스텀 명령·Red Team·리소스 자동화·지식 인덱싱 4단계 시스템
type: project
---

# 🤖 Claude Code를 "완벽한 분신"으로 만드는 4단계 실전 시스템

**철학**: "기록이 증명이다" + "외부 의존성 제로" + "자동화된 검증"
**대상**: FreeLang v6.1.0 이상
**시간**: 세션당 30초 setup으로 자동 활성화

---

## 🎯 시스템 1: Custom Alias (전용 명령)

### 구현: `.clauderules` 확장

```bash
# ~/.clauderules (사용자 글로벌) 또는 .claude/rules.md (프로젝트 로컬)

## 📋 커스텀 명령 (Alias)

### /check-logic
목적: 현재 코드의 타입 안전성 및 의존성 검사
실행: 모든 함수의 입출력 타입 → 의존성 선언 확인 → 외부 라이브러리 호출 전수 조사

명령문:
```
현재 작성된 모든 함수의 입출력 타입을 검사하세요.
의존성 선언이 없는 부분을 찾으세요.
외부 라이브러리가 호출되는지 확인하세요.
FORMAT: 함수명 | 입력(타입) | 출력(타입) | 의존성 | 외부호출 여부
```

### /save-proof
목적: 작업 내용 → GOGS 커밋 메시지로 자동 정리
실행: 작업 요약 → 기록 증명 포맷 → PLAN.md 자동 업데이트

명령문:
```
지금까지의 작업 내용을 요약해주세요.
GOGS 커밋 메시지 규격에 맞춰 정리하세요:
  [타입] 제목 (50자 이내)

  상세:
  - 변경사항 (불릿 포인트)
  - 테스트 결과
  - 증빙 데이터 (벤치마크, 타입 검사)

PLAN.md를 업데이트하세요.
마지막으로 '저장 필수, 기록이 증명이다 gogs'를 출력하세요.
```

### /verify-vanilla
목적: 외부 라이브러리 의존성 전수 조사
실행: 임포트 전수 조사 → 의존성 검사 → 순수 코드 대체안 제시

명령문:
```
외부 라이브러리 호출을 전수 조사하세요.
각 외부 호출마다:
  1. 라이브러리명
  2. 호출 위치 (파일:라인)
  3. 순수 코드로 대체 가능 여부
  4. 대체 방안 (있으면)

표 형식으로 정리하고,
'외부 의존성 제로' 달성 로드맵을 제시하세요.
```

### /red-team [aspect]
목적: 설계 약점 찾기 (반대 의견 제시)
실행: 성능/메모리/확장성 관점 3가지 취약점 도출

명령문:
```
방금 설계한 이 구조에서 가장 먼저 터질 것 같은 취약점을 분석하세요.
범주: [성능 | 메모리누수 | 확장성 | 보안]
분석 형식:
  1. 취약점: ~
     원인: ~
     재현 조건: ~
     영향도: 심각(1~5)
  2. ...

각 취약점마다 해결 방안을 제시하세요.
```

### /refactor-first [기능명]
목적: 기능 추가 전, 구조 개선 검토
실행: 기존 코드 성능 10% 향상 + 신규 기능 통합 방안 모색

명령문:
```
[기능명] 기능을 구현하기 전에,
기존 FreeLang 엔진 성능을 10% 향상시키면서
이 기능을 넣을 수 있는 구조적 대안을 제시해주세요.

분석 기준:
  - 메모리 레이아웃 (SoA vs AoS)
  - 캐시 지역성 (L1/L2/L3)
  - 호출 경로 (함수 호출 오버헤드)
  - 병렬화 가능성

결과 형식:
  현재: [메트릭]
  개선안: [메트릭]
  성능 향상: X%
```

---

## 🎯 시스템 2: Shadow Architect (검증 자동화)

### 구현: `.claude/prompts/architect-review.md`

```markdown
# Shadow Architect Review Protocol

매 코드 작성 완료 시 자동 실행:

## 1단계: Red Team Analysis
현재 설계에서:
- 가장 큰 성능 병목은?
- 메모리 누수 가능성은?
- 확장할 때 깨질 부분은?

## 2단계: Refactor-First Check
이 기능을 추가하기 전:
- 기존 코드에서 10% 성능 개선할 부분이 있는가?
- 구조 변경으로 신규 기능을 더 우아하게 넣을 수 있는가?
- 함수 호출 체인을 단축할 수 있는가?

## 3단계: Vanilla Check
의존성 검사:
- 외부 라이브러리 호출이 있는가?
- 있다면 순수 코드로 대체 가능한가?

## 4단계: Proof Record
기록 생성:
- 변경사항을 GOGS 커밋 형식으로 요약
- 벤치마크/타입 검사 증빙
- HISTORY.gogs에 누적
```

---

## 🎯 시스템 3: Termux 리소스 자동화

### 구현: `~/.claude/scripts/termux-monitor.sh`

```bash
#!/bin/bash
# Termux 환경에서 Claude Code 리소스 자동 관리

# 1. Node.js 메모리 모니터링
check_memory() {
  local max_mem=500  # MB (모바일 환경)
  local current=$(ps aux | grep node | grep -v grep | awk '{sum+=$6} END {print sum}')

  if [ $current -gt $max_mem ]; then
    echo "⚠️  메모리 경고: ${current}MB > ${max_mem}MB"
    echo "작업 중단 및 사용자 보고 실행"
    pkill -STOP node
    return 1
  fi
  return 0
}

# 2. 빌드 및 테스트 자동 실행
auto_test() {
  local project=$1
  cd "$project"

  echo "🔨 빌드 시작..."
  if go build ./... 2>&1 | tee /tmp/build.log; then
    echo "✅ 빌드 성공"
  else
    echo "❌ 빌드 실패"
    return 1
  fi

  echo "🧪 테스트 시작..."
  if go test ./... -v 2>&1 | tee /tmp/test.log; then
    echo "✅ 테스트 성공"
  else
    echo "❌ 테스트 실패"
    return 1
  fi

  # 한 줄 요약
  echo ""
  echo "📊 결과 요약:"
  local build_lines=$(wc -l < /tmp/build.log)
  local test_lines=$(wc -l < /tmp/test.log)
  echo "- 빌드: $build_lines 줄, 테스트: $test_lines 줄"
}

# 3. 배포 자동화
auto_deploy() {
  local message=$1

  echo "📦 배포 시작..."
  git add -A
  git commit -m "$message"
  git push origin master

  echo "✅ 배포 완료"
}

# Main
case "$1" in
  monitor) check_memory ;;
  test) auto_test "$2" ;;
  deploy) auto_deploy "$2" ;;
  *) echo "Usage: $0 {monitor|test|deploy} [args]" ;;
esac
```

---

## 🎯 시스템 4: 지식 보존 인덱싱

### 구현: `.claude/scripts/auto-map.py`

```python
#!/usr/bin/env python3
import os
import re
from pathlib import Path

def generate_map(project_root):
    """파일 구조 + 함수 호출 관계 지도 자동 생성"""

    map_data = {}

    # 1. 파일 구조 수집
    for root, dirs, files in os.walk(project_root):
        for file in files:
            if file.endswith(('.go', '.py', '.ts', '.rs')):
                filepath = os.path.join(root, file)
                map_data[filepath] = {
                    'functions': [],
                    'imports': [],
                    'calls': []
                }

    # 2. 함수 정의 및 호출 관계 분석
    for filepath in map_data.keys():
        with open(filepath, 'r') as f:
            content = f.read()

        # 함수 정의 찾기
        if filepath.endswith('.go'):
            functions = re.findall(r'func\s+\(?\w+\)?\s+(\w+)\s*\(', content)
            imports = re.findall(r'import\s+"([^"]+)"', content)
        elif filepath.endswith('.py'):
            functions = re.findall(r'def\s+(\w+)\s*\(', content)
            imports = re.findall(r'from\s+(\w+)\s+import|import\s+(\w+)', content)

        map_data[filepath]['functions'] = functions
        map_data[filepath]['imports'] = imports

    # 3. MAP.md 생성
    with open(os.path.join(project_root, 'MAP.md'), 'w') as f:
        f.write('# 프로젝트 맵\n\n')
        f.write('## 파일 구조\n\n```\n')

        for filepath in sorted(map_data.keys()):
            indent = '  ' * (filepath.count(os.sep) - project_root.count(os.sep))
            f.write(f'{indent}{os.path.basename(filepath)}\n')

        f.write('```\n\n')

        f.write('## 함수 및 호출 관계\n\n')
        for filepath, data in map_data.items():
            f.write(f'### {filepath}\n\n')
            for func in data['functions']:
                f.write(f'- `{func}`\n')
            if data['imports']:
                f.write(f'  - 의존: {", ".join(data["imports"])}\n')
            f.write('\n')

    print('✅ MAP.md 생성 완료')

if __name__ == '__main__':
    import sys
    project_root = sys.argv[1] if len(sys.argv) > 1 else '.'
    generate_map(project_root)
```

---

## 🎯 통합: `.claude/setup-perfect-double.sh`

### 세션 시작 시 자동 실행 (30초)

```bash
#!/bin/bash

echo "🤖 Claude Code 완벽한 분신 시스템 초기화..."

# 1. Custom Alias 등록
if [ ! -f ".claude/rules.md" ]; then
  cat > .claude/rules.md << 'EOF'
# 커스텀 명령 (Custom Alias)

/check-logic   → 타입 안전성 + 의존성 검사
/save-proof    → GOGS 커밋 자동 정리
/verify-vanilla → 외부 의존성 제거
/red-team      → 설계 약점 분석
/refactor-first → 성능 개선 + 기능 통합

---

# Shadow Architect 설정
자동으로 코드 완료 시 검증 실행:
- Red Team Analysis
- Refactor-First Check
- Vanilla Check
- Proof Record

---

저장 필수, 기록이 증명이다 gogs
EOF
fi

# 2. Termux 스크립트 설치
chmod +x ~/.claude/scripts/termux-monitor.sh 2>/dev/null || true

# 3. MAP.md 생성
python3 ~/.claude/scripts/auto-map.py . 2>/dev/null || true

echo "✅ 초기화 완료"
echo "📋 사용 가능한 명령:"
echo "   /check-logic       → 타입 검사"
echo "   /save-proof        → 커밋 정리"
echo "   /verify-vanilla    → 의존성 검사"
echo "   /red-team          → 약점 분석"
echo "   /refactor-first    → 성능 개선"
echo ""
echo "저장 필수, 기록이 증명이다 gogs"
```

---

## 📊 사용 흐름 (실제 예시)

### 예시 1: 새 기능 구현

```
1. /refactor-first "병렬 쿼리 엔진"
   → 기존 성능 10% 개선 + 신규 기능 통합 방안

2. (코딩)

3. /check-logic
   → 모든 함수 타입 검사
   → 외부 의존성 확인

4. /verify-vanilla
   → 외부 라이브러리 호출 전수 조사
   → 순수 코드 대체안

5. /red-team 성능
   → 가장 먼저 터질 취약점 3가지 분석
   → 해결책 제시

6. /save-proof
   → GOGS 커밋 메시지 자동 정리
   → PLAN.md 업데이트
   → "저장 필수, 기록이 증명이다 gogs" 출력
```

### 예시 2: 버그 수정 + 자동 배포

```
1. (버그 수정)

2. Termux: ./termux-monitor.sh test .
   → 빌드 + 테스트 자동 실행
   → 한 줄 요약 출력

3. /save-proof
   → 커밋 메시지 생성

4. Termux: ./termux-monitor.sh deploy "메시지"
   → git add/commit/push 자동 실행
```

---

## 🎯 설정 파일 체크리스트

```
✅ ~/.clauderules          (글로벌 규칙)
✅ .claude/rules.md        (프로젝트 로컬)
✅ .claude/prompts/architect-review.md
✅ ~/.claude/scripts/termux-monitor.sh
✅ ~/.claude/scripts/auto-map.py
✅ setup-perfect-double.sh (초기화 스크립트)
```

---

## 💡 핵심 슬로건

> **"저장 필수, 기록이 증명이다 gogs"**
>
> 매 작업마다:
> - ✅ 코드 작성 → /check-logic
> - ✅ 검증 완료 → /red-team
> - ✅ 최적화 → /refactor-first
> - ✅ 의존성 제거 → /verify-vanilla
> - ✅ 기록 저장 → /save-proof → GOGS 푸시

---

**시스템 준비 완료!** 🚀
