---
name: freelang-to-c - FreeLang to C 컴파일러 (신규)
description: FreeLang 소스 → C 코드 생성 → c-compiler-from-scratch → 네이티브 ELF. Phase 1-5 로드맵. 초기 프로젝트
type: project
---

# 🚀 freelang-to-c - FreeLang to C Transpiler

**FreeLang 소스를 C 코드로 변환하는 컴파일러**

---

## 📌 **현재 상태**

**상태**: 🟢 Just Started (2026-03-16)
**저장소**: https://gogs.dclub.kr/kim/freelang-to-c.git
**커밋**: `b799533 feat: freelang-to-c 프로젝트 초기화`
**브랜치**: main

---

## 🎯 **프로젝트 목표**

```
FreeLang 소스 코드
    ⬇️ Lexer/Parser
FreeLang IR (중간 표현)
    ⬇️ C 코드 생성
C 코드 (.c)
    ⬇️ c-compiler-from-scratch
ELF 바이너리 (네이티브)
    ⬇️ 실행
실행 결과
```

---

## 💡 **아이디어**

### 왜 C로 생성?
1. **성능**: 네이티브 기계어 생성
2. **호환성**: C 라이브러리 직접 사용 가능
3. **디버깅**: 생성된 C 코드 검사 가능
4. **이식성**: x86-64, ARM 등 다양한 플랫폼

### 세 언어의 결합
```
FreeLang  (표현력)
    ↓
C 코드  (성능)
    ↓
네이티브 바이너리 (속도)
```

---

## 📂 **프로젝트 구조**

```
freelang-to-c/
├── src/
│   ├── main.fl                  # 진입점
│   ├── codegen/                 # C 코드 생성
│   │   ├── c_codegen.fl        # IR → C 변환 엔진
│   │   ├── c_types.fl          # 타입 매핑
│   │   ├── c_expr.fl           # 표현식 변환
│   │   └── c_stmt.fl           # 문장 변환
│   │
│   └── runtime/                 # FreeLang 런타임 (C)
│       ├── freelang_runtime.c  # 핵심 런타임
│       ├── freelang.h
│       └── builtins.c          # 내장 함수
│
├── tests/
│   ├── test_simple.fl          # 간단한 예제
│   ├── test_func.fl            # 함수 테스트
│   ├── test_struct.fl          # 구조체 테스트
│   └── expected/
│       └── test_simple.c       # 예상 C 코드
│
├── examples/
│   ├── hello.fl                # "Hello, World!"
│   ├── fibonacci.fl            # 피보나치
│   ├── factorial.fl            # 팩토리얼
│   └── sort.fl                 # 정렬
│
├── scripts/
│   └── Makefile                # 빌드 스크립트
│
├── CLAUDE.md                   # 프로젝트 헌장
├── MEMORY.md                   # 진행 상황
└── README.md                   # 사용 가이드
```

---

## 🔧 **변환 규칙 (FreeLang → C)**

### 타입 매핑

| FreeLang | C |
|----------|---|
| `int` | `int` |
| `string` | `char*` |
| `[T]` (배열) | `T[]` |
| `*T` (포인터) | `T*` |
| `{a:T, b:U}` (구조체) | `struct { T a; U b; }` |
| `fn() -> T` (함수) | `T ()(...)` |

### 함수 변환

```
FreeLang:
fn foo(x: int) -> int {
  x + 1
}

↓ 생성 C:
int foo(int x) {
  return x + 1;
}
```

### 제어 흐름

```
FreeLang if/else → C if/else
FreeLang for/while → C for/while
FreeLang match → C switch
```

---

## 📋 **구현 로드맵**

### ✅ **Phase 0: 프로젝트 초기화** (완료)
- [x] GOGS 저장소 생성
- [x] 프로젝트 구조 설계
- [x] CLAUDE.md, MEMORY.md 작성

### ⏳ **Phase 1: 기본 타입 변환** (예정)
**목표**: 기본 데이터 타입과 구조체의 C 매핑
- [ ] int, char, float 타입 변환
- [ ] 배열 처리
- [ ] 구조체 필드 동기화
- [ ] 포인터 타입 처리
- **테스트**: test_simple.fl → test_simple.c

### ⏳ **Phase 2: 표현식 & 함수** (예정)
**목표**: 산술/논리 연산, 함수 호출
- [ ] 산술 연산 변환
- [ ] 논리 연산 변환
- [ ] 함수 호출 변환
- [ ] 메모리 접근 (포인터 역참조)
- **테스트**: test_func.fl

### ⏳ **Phase 3: 제어 흐름** (예정)
**목표**: if/while/for/switch 변환
- [ ] if/else 생성
- [ ] for 루프 생성
- [ ] while 루프 생성
- [ ] switch/case 생성
- **테스트**: 제어 흐름 예제

### ⏳ **Phase 4: 모듈 시스템** (예정)
**목표**: 함수 조합, 외부 함수 호출
- [ ] 함수 선언 및 정의
- [ ] 함수 포인터
- [ ] 외부 C 함수 호출 (libc)
- [ ] 헤더 파일 생성

### ⏳ **Phase 5: 통합 & 최적화** (예정)
**목표**: 완전한 파이프라인
- [ ] FreeLang → C 파이프라인 통합
- [ ] c-compiler-from-scratch 연동
- [ ] C 코드 최적화
- [ ] 성능 벤치마크

---

## 🧪 **테스트 전략**

### 단위 테스트
```
test_simple.fl  → c_codegen → test_simple.c
                → c-compiler-from-scratch
                → a.out (실행)
                → 결과 비교
```

### 예제 프로그램
```
examples/hello.fl          → 컴파일 → hello.out → "Hello, World!"
examples/fibonacci.fl      → 컴파일 → fib.out → fib(10) = 55
examples/sort.fl           → 컴파일 → sort.out → 배열 정렬
```

### 성능 비교
```
FreeLang 직접 실행 vs
FreeLang → C → 네이티브 실행
(속도, 메모리)
```

---

## 🔗 **의존성**

### 필수
- **c-compiler-from-scratch**: C → ELF 컴파일
- **FreeLang**: 소스 파싱 및 IR 생성

### 선택
- **clang/gcc**: C 코드 검증 (테스트용)

---

## 📊 **현재 통계**

| 항목 | 값 |
|------|-----|
| 파일 수 | 3 (CLAUDE.md, MEMORY.md, README.md) |
| 커밋 | 1 |
| 코드량 | ~178줄 (문서 포함) |
| 상태 | 초기 프로젝트 |

---

## 🎯 **다음 액션**

1. **examples/hello.fl 작성** (가장 간단한 예제)
   ```
   fn main() {
     println("Hello, World!")
   }
   ```

2. **예상 C 코드 작성**
   ```c
   #include <stdio.h>

   int main() {
     printf("Hello, World!\n");
     return 0;
   }
   ```

3. **Phase 1: 기본 타입 변환 시작**
   - IR → C 타입 매핑 규칙 구현
   - 기본 테스트 케이스

4. **c-compiler-from-scratch 통합**
   - 생성된 C 코드 자동 컴파일
   - ELF 바이너리 생성

---

## 📝 **노트**

**설계 철학**:
- 간단한 것부터 시작 (hello.fl)
- 각 Phase마다 완전한 예제 제공
- 생성된 C 코드의 읽기 가능성 우선
- 성능 최적화는 나중에

**도전 과제**:
- FreeLang의 동적 특성 → C의 정적 타입 시스템
- 메모리 관리 (GC vs 수동)
- 문자열 처리
- 에러 처리

---

**상태**: 🟡 Phase 1 완료 | Phase 2 설계 중

---

## 📊 **Phase 1 완성 (2026-03-16)**

✅ **완료된 것**:
- 타입 매핑 시스템 (FlType → C type) 구현
- 17개 단위 테스트 통과 (100%)
- 코드젠 인프라 (변수, 함수, 제어문)
- 생성 C 코드 컴파일 및 실행 확인
- 아키텍처: 유니온 기반 타입 설계 (c-compiler 패턴 활용)

**코드량**: ~1,000줄 (테스트 포함)
**Commit**: 109f7ae

---

## 🔄 **FreeLang + C 시너지**

### 1. **타입 시스템 우위**
- FreeLang: 유연한 타입 (동적 + 정적)
- C: 고성능, 예측 가능한 메모리
- **시너지**: FreeLang의 고수준 추상화 → C의 저수준 제어
  ```
  FL: fn process(data: [i32]) -> bool { ... }
  → C: bool process(int *data, int len) { ... }
  ```

### 2. **메모리 관리**
- FreeLang GC 제거 → C의 수동 메모리 관리
- 스택 할당으로 성능 향상 (GC 오버헤드 제거)
- **시너지**: 소유권 추적 → C 라이프타임

### 3. **동시성**
- FreeLang async/await → C 스레드 풀
- 비동기 → 네이티브 멀티스레드
- **시너지**: 고수준 언어의 편의성 + 저수준 성능

### 4. **라이브러리 통합**
- FreeLang stdlib → C 바인딩 자동 생성
- 기존 C 라이브러리 직접 활용
- **시너지**: 에코시스템 무제한 확장

### 5. **최적화 기회**
- FreeLang IR → C → 재최적화 (c-compiler SAI 등)
- 2단계 최적화 (FreeLang + C 컴파일러)
- **시너지**: 이중 최적화로 성능 극대화

