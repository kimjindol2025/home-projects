# 🔍 freelang-v2 완전 독립성 분석

**작성일**: 2026-03-26
**목표**: koffi FFI 의존성 제거 가능성 조사

---

## 📊 현황

### koffi 사용 현황

**사용 파일** (3개):
1. `src/ffi/c-function-caller.ts` - **핵심** (FFI 엔진)
2. `src/stdlib/serial.ts` - 시리얼 통신
3. `rt/entropy_core.v2.ts` - 엔트로피 관리

**koffi 관련 코드**: ~854줄

### koffi의 역할

```
koffi = FFI (Foreign Function Interface) 라이브러리
  ↓
C 라이브러리와의 상호작용
  ├─ 동적 라이브러리 로드
  ├─ C 함수 호출
  ├─ 메모리/타입 마샬링
  └─ 네이티브 코드 실행
```

---

## 🔎 상세 분석

### 1. c-function-caller.ts (FFI 엔진)

**목적**: C 함수 동적 호출 엔진

**koffi 사용 방식**:
```typescript
import { load as loadLibrary } from 'koffi';

// C 라이브러리 로드
const lib = loadLibrary('/usr/local/lib/libstream.so');

// C 함수 호출
const result = koffiFunc(...args);
```

**로드 대상 라이브러리**:
```
- /usr/local/lib/libstream.so   (스트림)
- /usr/local/lib/libws.so        (WebSocket)
- /usr/local/lib/libhttp.so      (HTTP)
- /usr/local/lib/libhttp2.so     (HTTP/2)
- /usr/local/lib/libevent_loop.so (이벤트 루프)
- /usr/local/lib/libtimer.so     (타이머)
```

**평가**:
- ❌ 제거 불가능 (FFI가 핵심 기능)
- koffi 없으면 C 라이브러리 호출 불가

---

### 2. serial.ts (시리얼 통신)

**용도**: 직렬 포트 통신

**koffi 사용**:
```typescript
try {
  const koffi = require('koffi');
  // 시리얼 포트 접근
} catch {
  // koffi 없을 시 폴백
}
```

**특징**:
- ✅ Try-catch로 처리 (선택적)
- ⚠️ koffi 없으면 시리얼 통신 불가
- 일부 기능은 fallback 있음

**평가**:
- 🟡 부분적으로 제거 가능
- 시리얼 통신만 비활성화

---

### 3. entropy_core.v2.ts (엔트로피)

**용도**: 난수 생성 (보안 관련)

**koffi 사용**: 네이티브 난수 생성기

**평가**:
- 🟡 선택적 (fallback 있을 수 있음)

---

## ❓ 완전 독립 가능성

### 현실적 평가

| 기능 | koffi 필수? | 제거 비용 | 대체 방안 |
|------|-----------|---------|---------|
| **FFI 엔진** | ✅ 필수 | 높음 | Pure JS 재구현 (비현실적) |
| **시리얼 통신** | 부분 | 중간 | Node.js `serialport` 패키지 |
| **HTTP/WebSocket** | ❓ 불명 | 확인 필요 | Node.js `http/ws` 표준 |
| **이벤트 루프** | ❓ 불명 | 확인 필요 | Node.js 표준 Event Loop |
| **난수 생성** | 부분 | 낮음 | `crypto` 모듈 |

---

## 🛠️ 완전 독립을 위한 시나리오

### 방안 1: koffi 완전 제거 ❌ (불가능)

**문제**:
- FFI 엔진이 C 라이브러리 호출 담당
- 이를 제거하면 `libstream.so`, `libws.so` 등 사용 불가
- 이들 라이브러리를 JS로 재구현해야 함
- 비용: **매우 높음** (수주 이상)

### 방안 2: koffi 선택적 비활성화 ⚠️ (부분 가능)

**가능한 부분**:
1. **시리얼 통신** → Node.js `serialport` 사용
2. **난수 생성** → Node.js `crypto` 사용
3. **기본 HTTP/WS** → Node.js 표준 라이브러리

**불가능한 부분**:
- C 네이티브 라이브러리 호출 필요 시

**비용**: 중간 (1-2주)

**명령어**:
```bash
# 1. serialport 패키지 추가
npm install serialport

# 2. koffi 제거
npm uninstall koffi

# 3. 코드 수정
# - serial.ts: koffi → serialport
# - entropy_core.v2.ts: native → crypto
# - c-function-caller.ts: stub 구현
```

### 방안 3: koffi를 Optional 의존성으로 변경 ✅ (권장)

**개념**:
```json
{
  "dependencies": {
    "serialport": "^10.0.0",  // 선택적 (시리얼 통신)
    "chalk": "^4.1.2"
    // koffi 제거
  },
  "optionalDependencies": {
    "koffi": "^2.15.1"  // 있으면 사용, 없으면 스킵
  }
}
```

**코드**:
```typescript
let ffiEngine: any = null;

try {
  // koffi가 설치되어 있으면 사용
  const { load } = require('koffi');
  ffiEngine = { load };
} catch {
  // koffi 없을 시 stub
  ffiEngine = {
    load: () => {
      throw new Error('FFI support disabled: install koffi');
    }
  };
}
```

**비용**: 낮음 (1-2일)
**장점**:
- koffi 없어도 기본 기능 동작
- C 라이브러리가 필요한 경우만 설치

---

## 💡 최종 권장사항

### 현재 상태
```
freelang-v2:
├─ 필수 의존성: 6개 (chalk, express, better-sqlite3 등)
├─ koffi: C FFI 엔진 (완전 제거 불가)
└─ 결론: 완전 독립 불가능
```

### 3가지 옵션

**옵션 1: 현 상태 유지** ✅ 간단
- 의존성 7개 유지
- 모든 기능 동작
- 비용: 0

**옵션 2: koffi 선택적 처리** ⭐ 권장
- serialport로 시리얼 통신
- crypto로 난수 생성
- koffi 제거 또는 optional
- 비용: 1-2일
- 결과: 6개 의존성 (koffi 제거)

**옵션 3: 완전 JS 재구현** ❌ 비현실적
- 모든 C 라이브러리 JS로 변환
- 비용: 수주 이상
- 성능 저하

---

## 📋 실행 계획 (옵션 2 선택 시)

### Phase 1: 분석 (1일)
1. 각 C 라이브러리의 정확한 역할 파악
2. JS 대체 방안 검증
3. 호환성 확인

### Phase 2: 구현 (1-2주)
1. serialport 통합
2. crypto 모듈 사용
3. c-function-caller stub 구현
4. 테스트

### Phase 3: 검증 (2-3일)
1. 모든 기능 테스트
2. 성능 비교
3. 문서화

---

## 🎯 결론

**Q: freelang-v2는 완전 독립 가능한가?**

**A: 아니오. 하지만 개선 가능합니다.**

- ❌ koffi 없이 완전 독립 불가능 (FFI가 핵심)
- ✅ 부분적으로 의존성 제거 가능 (6개 → 5개 또는 optional)
- ⭐ **권장**: koffi를 optional 의존성으로 변경

**대안**:
- **지금**: fv2-lang-go 사용 (완전히 독립적, Go 기본 라이브러리만 사용)
- **나중**: freelang-v2는 선택적 기능으로 유지

---

**조사일**: 2026-03-26
**상태**: ✅ 완료
