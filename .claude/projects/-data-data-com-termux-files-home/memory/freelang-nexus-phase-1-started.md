---
name: FreeLang Nexus Phase 1 - 시작! (2026-03-21)
description: 완전체 언어 설계 시작. V + Python 혼합 Lexer 개발 (목표: 1주)
type: project
---

# 🚀 FreeLang Nexus Phase 1 시작!

**결정 날짜**: 2026-03-21
**프로젝트 이름**: **FreeLang Nexus** (V와 Python의 결합점)
**상태**: 🟡 Phase 1 진행 중 (Lexer 통합)
**소요 예상**: 1주 (토큰화)

---

## 🎯 **Phase 1 목표**

한 파일에서 V와 Python 스타일을 섞어도 Lexer가 정확히 구분

```fl
@mode(v)
fn fast(x: i64) -> i64 { return x * x }

@mode(python)
def slow(items):
  return [i*2 for i in items]

fn mixed():
  return [fast(5), slow([1,2,3])]
```

이 코드의 토큰 스트림을 정확히 분리

---

## 📋 **구현 계획**

### **Task 1: Token 타입 정의**
- V 키워드: fn, let, struct, : (타입)
- Python 키워드: def, class, indent, dedent
- 혼합: @mode(v), @mode(python)

### **Task 2: Mode Detector**
- @mode(v) / @mode(python) 인식
- 자동 모드 감지

### **Task 3: NexusLexer 작성**
- 토큰화 with 모드 추적
- Python 들여쓰기 처리 (INDENT/DEDENT)

### **Task 4: 테스트 작성**
- V 키워드 인식 테스트
- Python 키워드 인식 테스트
- 혼합 코드 테스트
- 타입 표기 테스트
- 들여쓰기 토큰 테스트

---

## 🎨 **설계 핵심**

### Zero-copy Architecture (Phase 4 목표)
```
V 메모리 (정적 struct)
  ↓ (주소 직접 공유)
Python VM (동적 object)

= 복사 없이 주소만 전달
```

### Mode 전환 메커니즘
```
@mode(v)
  ↓
TokenType.MODE_MARKER (값: @mode(v))
  ↓
currentMode = 'v'
  ↓
이후 토큰: mode='v' 태그

@mode(python)로 전환되면:
  ↓
currentMode = 'python'
  ↓
이후 토큰: mode='python' 태그
```

---

## 📁 **파일 구조**

```
projects/freelang-nexus/
├── src/nexus/lexer/
│   ├── token.ts
│   ├── nexus-lexer.ts
│   └── mode-detector.ts
├── tests/
│   └── nexus-lexer.test.ts
└── docs/
    └── NEXUS_DESIGN.md
```

---

## ✅ **완료 조건**

- [ ] Token 타입 30개 이상 정의
- [ ] NexusLexer 클래스 구현
- [ ] 5개 이상 테스트 케이스 통과
- [ ] nexus_basic.fl 예제 실행 성공
- [ ] GOGS 커밋

---

## 💡 **핵심 학습 포인트**

1. **Lexer의 상태 관리**
   - 모드 전환이 어떻게 토큰 분류를 바꾸는가?

2. **Python의 들여쓰기 처리**
   - INDENT/DEDENT 토큰이 왜 필요한가?

3. **Type annotation의 토큰화**
   - V의 `: i64` 같은 타입 표기를 어떻게 처리?

---

## 📚 **참고: 기존 코드 활용**

- PyFree Lexer (Python 들여쓰기) → 참고
- FV Lexer (V 키워드) → 참고
- 통합 아이디어: 둘을 한 Lexer에서 관리

---

**시작 시간**: 2026-03-21 (지금!)
**목표 완료**: 2026-03-28 (1주)
**다음**: Phase 2 (Parser 통합)

---

> "저장 필수, 기록이 증명이다 GOGS!"
