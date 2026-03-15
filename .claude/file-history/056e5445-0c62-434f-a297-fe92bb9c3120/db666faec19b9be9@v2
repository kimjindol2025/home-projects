# 🦎 VSCode에서 FreeLiulia 사용하기

FreeLiulia 문법 강조를 VSCode에서 설정하는 방법입니다.

## 설치 방법

### 방법 1: 수동 설정 (권장)

1. **문법 파일 복사**
   ```bash
   cp docs/freeliulia.tmLanguage.json ~/.config/Code/User/
   ```

2. **VSCode 사용자 설정에 추가**
   `Ctrl+Shift+P` → "Open User Settings (JSON)" → 아래 내용 추가

   ```json
   {
     "files.associations": {
       "*.fl": "freeliulia"
     },
     "[freeliulia]": {
       "editor.formatOnSave": true,
       "editor.tabSize": 2
     }
   }
   ```

### 방법 2: 구성 파일 (macOS/Linux)

FreeLiulia 프로젝트 디렉토리에 `.vscode/settings.json` 생성:

```json
{
  "files.associations": {
    "*.fl": "freeliulia"
  },
  "[freeliulia]": {
    "editor.formatOnSave": true,
    "editor.tabSize": 2
  }
}
```

### 방법 3: VS Code Extension 설치 (향후)

현재는 공식 확장이 없지만, 다음과 같이 로컬 확장으로 등록 가능:

```bash
cd ~/.vscode/extensions
mkdir freeliulia-syntax
cd freeliulia-syntax
# package.json 및 문법 파일 복사
```

---

## 기능

### 문법 강조 (Syntax Highlighting)

✅ **한글 키워드**
- 제어흐름: `함수`, `변수`, `만약`, `아니면`, `반복`, `동안`, `반환`, `끝`
- 타입: `정수`, `실수`, `논리값`, `문자열`, `배열`, `없음`
- 값: `참`, `거짓`

✅ **내장 함수**
- 문자열: `길이`, `포함`, `시작`, `끝`, `찾기`, `대문자`, `소문자`, `연결`, `자르기`, `바꾸기`
- 고급 문자열: `반복`, `채우기`, `나누기_앞`, `나누기_뒷`
- 수학: `절대값`, `제곱`, `제곱근`, `올림`, `내림`, `최대값`, `최소값`
- 배열: `배열길이`, `배열합계`, `배열평균`, `배열최대`
- 파일: `파일열기`, `파일닫기`, `줄읽기`, `줄쓰기`, `파일존재`
- 타입: `정수인가`, `실수인가`, `논리값인가`, `문자열인가`

✅ **구문 요소**
- 주석: `#` 이후 줄 끝까지
- 문자열: `"..."` 또는 `'...'`
- 숫자: 정수 및 실수
- 연산자: `+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `=`
- 그리고/또는/아니

---

## 사용 예제

### 기본 코드 (문법 강조 예)

```freeliulia
# 주석: 이 줄은 회색으로 표시됩니다

함수 계산(수1: 정수, 수2: 정수) -> 정수
    변수 합 = 수1 + 수2
    반환 합
끝

# 메인 코드
변수 입력 = "hello world"
변수 결과 = 계산(10, 20)

만약 결과 > 25
    변수 메시지 = "결과가 25보다 큽니다"
아니면
    변수 메시지 = "결과가 25 이하입니다"
끝
```

**색상 맵**:
- 🟦 **파란색**: 키워드 (`함수`, `변수`, `만약`)
- 🟩 **초록색**: 문자열 (`"..."`)
- 🟨 **노란색**: 숫자 (`42`, `3.14`)
- ⚫ **회색**: 주석 (`#`)
- 🟪 **보라색**: 내장 함수 (`길이`, `포함`)

---

## 테마 커스터마이징

### 테마별 색상 조정

**Dark+ (기본 어두운 테마)**:
```json
"editor.tokenColorCustomizations": {
  "[Dark+]": {
    "keywords": "#569CD6",
    "strings": "#CE9178",
    "numbers": "#B5CEA8",
    "comments": "#6A9955"
  }
}
```

**Light+ (밝은 테마)**:
```json
"editor.tokenColorCustomizations": {
  "[Light+]": {
    "keywords": "#0000FF",
    "strings": "#A31515",
    "numbers": "#098658",
    "comments": "#008000"
  }
}
```

---

## 문제 해결

### Q: `.fl` 파일이 인식되지 않습니다
**A**: VSCode를 재시작하고, `Ctrl+K Ctrl+T`로 테마를 선택한 후 다시 시도하세요.

### Q: 문법 강조가 작동하지 않습니다
**A**:
1. 파일 확장자가 `.fl`인지 확인
2. VSCode 상단 우측 "언어 모드" 드롭다운에서 "FreeLiulia" 선택
3. VSCode 재시작

### Q: 한글이 제대로 표시되지 않습니다
**A**: VSCode 설정에서 폰트를 한글 지원하는 것으로 변경:
```json
"editor.fontFamily": "'Noto Sans Mono CJK KR', Menlo, Monaco, 'Courier New'"
```

---

## 확장 기능 (향후)

현재 구현:
- ✅ 문법 강조

향후 구현 계획:
- ⏳ 자동 완성 (IntelliSense)
- ⏳ 코드 포매팅
- ⏳ 디버거 통합
- ⏳ 린트 (Linter)
- ⏳ 테스트 실행기

---

## 피드백

이 설정 파일을 개선하고 싶으신 부분이 있으면:

1. GitHub Issue 생성
2. Pull Request 제출
3. 토론 게시판에 제안

---

**마지막 업데이트**: 2026-03-11
**버전**: 1.0
**작성자**: FreeLiulia Team
