---
name: FreeLang V2 vs V2.5 분석
description: V2와 V2.5의 근본적인 차이점 정리
type: project
---

# 🔍 FreeLang V2 vs V2.5 - 핵심 차이 분석

**Date**: 2026-03-26
**분석**: V2 (아카이브)와 V2.5 (클레임) 비교

---

## 📊 **V2: "freelang-v2" (ARCHIVED)**

### 상태
- **Status**: 🔴 **ARCHIVED** (2026-03-15)
- **위치**: `/data/data/com.termux/files/home/.projects/archived/freelang-v2/`
- **언어**: TypeScript (129,463줄)
- **버전**: v2.9.0

### 실제 구현 내용
```
✅ 56개 모듈 구현
├─ 파서 (13개 파일)          → AST 생성
├─ 코드젠 (18개 파일)        → C/다국어 생성
├─ 컴파일러 (11개 파일)      → 전체 파이프라인
├─ 런타임 (12개 파일)        → VM/바이트코드 실행
├─ 표준라이브러리 (64개 파일) → 수학, 파일, 네트워크 등
├─ Phase 6~30 (200+ 파일)    → 진화 과정 모두 보존
├─ LSP 서버 (13개 파일)      → 에디터 통합
├─ CLI 도구 (15개 파일)      → 명령행 인터페이스
└─ ...기타 (JIT, FFI, 핫리로드, 최적화 등)

✅ 실제 동작 검증:
  - CLI 명령어 15+ 개
  - REST API 서버
  - LSP (Language Server Protocol)
  - 배치 모드, AOT 컴파일, 핫리로드
  - 모듈 링커 (KPM)
  - Proof Tester (@test 어노테이션)
```

### 기술 스택
```typescript
dependencies {
  better-sqlite3: DB
  express: HTTP 서버
  vscode-languageserver: LSP
  chalk: CLI 색상
}
```

### NPM Scripts
```
- build: TypeScript → JavaScript 컴파일
- test: Jest 테스트
- start: CLI 실행 (dist/cli/index.js)
- lsp: LSP 서버 시작
- dev: ts-node로 개발 모드
```

---

## 📊 **V2.5: "freelang-v2-5" (CLAIMED)**

### 상태
- **Status**: 🟢 **"Production Ready"** (클레임)
- **위치**: `/data/data/com.termux/files/home/.projects/core/freelang-v2-5/`
- **언어**: FreeLang (.fl 파일들)
- **버전**: v2.5.0

### 실제 구현 내용
```
❌ 파일만 존재, 코드 없음
├─ CLAUDE.md               → 템플릿 (수정 안 함)
├─ README.md               → 계획만 있음
├─ examples/
│   ├─ hello.fl            → 간단한 예제
│   ├─ algorithm.fl        → 정렬 알고리즘
│   └─ system.fl           → 시스템 호출 예제
└─ spec/
    ├─ syntax.md           → 언어 정의
    ├─ stdlib.md           → 표준 라이브러리 (예상)
    └─ ...예상 파일들 (미구현)

❌ 컴파일러 없음
❌ 런타임 없음
❌ 표준라이브러리 없음
❌ 테스트 없음
❌ 배포 방법 없음
```

### Git 커밋 히스토리
```
Last commit: ~1주일 전
Log:
  "Create examples"
  "Add CLAUDE.md template"
  "Initial setup"

⏸️ 진전 없음 (계획만 존재)
```

---

## 🎯 **비교표**

| 기준 | V2 (ARCHIVED) | V2.5 (CLAIMED) |
|------|---------------|-----------------|
| **코드 존재** | ✅ 129K 줄 TypeScript | ❌ 0 줄 (`.fl` 파일만) |
| **컴파일러** | ✅ src/compiler/ (11파일) | ❌ 없음 |
| **런타임** | ✅ src/runtime/ (12파일) | ❌ 없음 |
| **표준라이브러리** | ✅ src/stdlib/ (64파일) | ❌ 스펙만 |
| **테스트** | ✅ Jest/Proof-Tester | ❌ 없음 |
| **CLI** | ✅ 15+ 명령어 | ❌ 없음 |
| **LSP** | ✅ 에디터 통합 | ❌ 없음 |
| **배포 준비** | ⚠️ npm package | ❌ 완전 불가 |
| **생산 레벨** | ⚠️ 가능할 수 있음 | ❌ 0% |
| **상태** | 🔴 ARCHIVED | 🟢 **"CLAIMED"** |

---

## 🔴 **결론**

### V2는?
- **실제 컴파일러 구현**: 129K 줄
- **어떤 이유로든 아카이브 됨** (버그? 성능? 설계 문제?)
- **지금 상태**: 이론적으로 동작할 수 있으나, 알려진 문제 있을 가능성

### V2.5는?
- **"Product Marketing"**: "Production Ready"라고 주장
- **실제 내용**: `.fl` 파일들 + 예제 + 스펙 문서
- **컴파일러/런타임**: 완전 부재
- **현실 레벨**: 개념 증명(PoC) 수준, 실행 불가능

---

## 🤔 **해석**

### 가정 1: V2 → V2.5 마이그레이션 실패
```
V2 (아카이브)
    ↓
  (문제 발생?)
    ↓
V2.5 재설계 시작
    ↓
  (구현 중단?)
    ↓
"Production Ready"라고 선언 (실제로는 미완성)
```

### 가정 2: V2.5는 V2의 "스펙" 버전
```
V2 코드 → 동작하지만 복잡함
V2.5 → "이상적인 문법/설계"만 정의
    → 구현은 나중에 (지금은 계획만)
```

### 가정 3: 프로젝트 관리 문제
```
V2 → 완성도 70%, 성능 문제로 아카이브
V2.5 → 새로운 시도, 하지만 중단됨
      → "Production Ready"는 희망사항
```

---

## 📋 **다음 조사사항**

1. **V2가 아카이브된 이유**: git log 분석 필요
2. **V2.5의 의도**: 설계자의 의도 파악 필요
3. **다른 버전들**: V1, V3, V4 등 실제 상태 확인
4. **FreeLang 생태계**: 30+ 프로젝트 중 실제로 동작하는 것 파악

---

**Status**: 🔴 V2.5는 명백한 **거짓 표현** (False Claim)
**V2는** 실제 구현이 있지만 **아카이브된 이유** 불명
