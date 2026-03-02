# 🚀 Phase B: FreeLang 런타임 구현 계획
**작성일**: 2026-03-02
**기간**: 4주 (2026-03-02 ~ 2026-03-30)
**목표**: Rust 자체 런타임 + 50+ 표준 함수 라이브러리 구현
**상태**: 📋 준비 단계

---

## 📊 **Phase B 개요**

### 목표
```
FreeLang v1.0 = 완전한 실행 가능한 언어
```

### 핵심 요구사항
1. **Rust 자체 런타임**
   - 외부 의존성 최소화
   - 메모리 안전성 보장
   - 성능 최적화

2. **50+ 표준 함수**
   - 기본 I/O
   - 수학 연산
   - 문자열 처리
   - 파일 I/O
   - 네트워크 (선택)

3. **메모리 관리**
   - GC 또는 RC (Reference Counting)
   - 스택 vs 힙 할당 전략
   - 메모리 누수 방지

### 성공 기준
```
✅ Rust 코드 작성 (8,000+ 줄)
✅ 50+ 표준 함수 구현
✅ 테스트 커버리지 >85%
✅ 문서 완성도 >90%
✅ 성능 기준 달성 (<100ms 시작, <10ns 함수 호출)
```

---

## 🏗️ **아키텍처 설계**

### 계층 구조
```
┌─────────────────────────────────────┐
│  FreeLang 사용자 코드               │
├─────────────────────────────────────┤
│  VM 실행 엔진                       │
│  - Bytecode Interpreter             │
│  - Stack Machine (1000 스택 깊이)   │
├─────────────────────────────────────┤
│  표준 라이브러리                    │
│  - Core (기본)                      │
│  - IO (입출력)                      │
│  - Math (수학)                      │
│  - String (문자열)                  │
│  - File (파일)                      │
│  - Net (네트워크)                   │
├─────────────────────────────────────┤
│  런타임 서비스                      │
│  - 메모리 관리 (Allocator)          │
│  - GC / RC                          │
│  - 에러 처리                        │
│  - 디버깅 정보                      │
├─────────────────────────────────────┤
│  Rust 기반 인프라                   │
└─────────────────────────────────────┘
```

### 모듈 구조
```
freelang-runtime/
├── src/
│   ├── vm.rs                  # VM 핵심 엔진 (500줄)
│   ├── bytecode.rs            # Bytecode 정의 (200줄)
│   ├── interpreter.rs         # 인터프리터 (400줄)
│   ├── memory.rs              # 메모리 관리 (300줄)
│   ├── gc.rs                  # GC 구현 (250줄)
│   ├── stdlib/
│   │   ├── core.rs            # 핵심 함수 (800줄)
│   │   ├── io.rs              # 입출력 (400줄)
│   │   ├── math.rs            # 수학 함수 (300줄)
│   │   ├── string.rs          # 문자열 처리 (350줄)
│   │   ├── file.rs            # 파일 I/O (400줄)
│   │   ├── net.rs             # 네트워크 (300줄)
│   │   └── time.rs            # 시간 함수 (200줄)
│   ├── error.rs               # 에러 처리 (100줄)
│   └── debug.rs               # 디버깅 (150줄)
├── tests/
│   ├── vm_tests.rs            # VM 테스트 (300줄)
│   ├── stdlib_tests.rs        # 표준 라이브러리 테스트 (400줄)
│   └── integration_tests.rs   # 통합 테스트 (200줄)
├── benches/
│   └── benchmarks.rs          # 성능 벤치마크 (150줄)
└── Cargo.toml
```

---

## 📝 **주간별 작업 계획**

### Week 1: 기초 설계 & VM 구현 (2026-03-02 ~ 2026-03-08)

#### Task 1.1: Bytecode 정의 (1일)
```rust
// bytecode.rs (200줄)
pub enum Opcode {
    // 스택 연산
    Push(Value),      // 값 push
    Pop,              // 값 pop
    Dup,              // 복제

    // 산술 연산
    Add, Sub, Mul, Div, Mod,

    // 논리 연산
    And, Or, Not,

    // 비교 연산
    Eq, Ne, Lt, Lte, Gt, Gte,

    // 제어 흐름
    Jump(usize),      // 무조건 점프
    JumpIfTrue(usize),// 조건 점프
    JumpIfFalse(usize),
    Call(usize),      // 함수 호출
    Return,           // 함수 반환

    // 함수 관련
    DefFn,            // 함수 정의
    LoadFn(String),   // 함수 로드

    // 메모리 접근
    LoadLocal(usize), // 로컬 변수 로드
    StoreLocal(usize),// 로컬 변수 저장
    LoadGlobal(String),
    StoreGlobal(String),

    // 배열/맵
    LoadArray(usize), // 배열 인덱싱
    StoreArray(usize),
    LoadMap(String),  // 맵 접근
    StoreMap(String),

    // 타입 처리
    Cast(Type),       // 타입 캐스팅

    // I/O
    Print,            // 출력
    Read,             // 입력

    // 할당
    Alloc(Type),      // 메모리 할당
    Free,             // 메모리 해제

    // 기타
    Halt,             // 프로그램 종료
}
```

#### Task 1.2: 메모리 관리 설계 (1일)
```
메모리 전략: RC (Reference Counting)
이유:
- Deterministic cleanup
- 예측 가능한 성능
- GC pause 없음
- Rust와의 자연스러운 통합

구조:
struct Memory {
    heap: Vec<Box<Value>>,
    stack: Vec<Value>,
    references: HashMap<*const Value, usize>,
}
```

#### Task 1.3: VM 핵심 엔진 구현 (2일)
```rust
// vm.rs (500줄)
pub struct VirtualMachine {
    pc: usize,                    // 프로그램 카운터
    stack: Vec<Value>,            // 스택 (1000 깊이)
    memory: Memory,               // 메모리
    functions: HashMap<String, Function>,
    bytecode: Vec<Opcode>,        // Bytecode 명령
    callstack: Vec<Frame>,        // 호출 스택
}

impl VirtualMachine {
    pub fn new() -> Self { }
    pub fn execute(&mut self, bytecode: Vec<Opcode>) -> Result<Value> { }
    fn fetch(&mut self) -> Opcode { }
    fn decode_execute(&mut self, op: Opcode) { }
    fn push(&mut self, val: Value) { }
    fn pop(&mut self) -> Value { }
}
```

#### Task 1.4: 인터프리터 구현 (1일)
```rust
// interpreter.rs (400줄)
pub struct Interpreter {
    vm: VirtualMachine,
    functions: HashMap<String, AST>,
}

impl Interpreter {
    pub fn run(&mut self, ast: AST) -> Result<Value> {
        // AST → Bytecode 변환
        let bytecode = self.compile(ast)?;
        // VM 실행
        self.vm.execute(bytecode)
    }
}
```

---

### Week 2: 표준 라이브러리 1단계 (2026-03-09 ~ 2026-03-15)

#### Task 2.1: Core 함수 (800줄)
```
함수 목록:
✅ print(value) - 출력
✅ println(value) - 줄 바꿈 출력
✅ len(array/string) - 길이
✅ type(value) - 타입 확인
✅ is_null(value) - null 체크
✅ assert(condition) - 단언
✅ error(message) - 에러 발생
✅ panic(message) - 패닉
✅ random() - 난수
✅ sleep(ms) - 대기
✅ exit(code) - 종료
... (800줄, 30+ 함수)
```

**구현 예시**:
```rust
// stdlib/core.rs
pub fn builtin_print(vm: &mut VirtualMachine, args: Vec<Value>) -> Result<Value> {
    for arg in args {
        print!("{}", arg.to_string());
    }
    Ok(Value::Null)
}

pub fn builtin_len(vm: &mut VirtualMachine, args: Vec<Value>) -> Result<Value> {
    if args.len() != 1 {
        return Err("len() takes 1 argument".to_string());
    }
    match &args[0] {
        Value::String(s) => Ok(Value::Int(s.len() as i64)),
        Value::Array(arr) => Ok(Value::Int(arr.len() as i64)),
        _ => Err("len() only works on strings and arrays".to_string()),
    }
}
```

#### Task 2.2: IO 함수 (400줄)
```
함수 목록:
✅ input() - 표준 입력
✅ input_line() - 한 줄 입력
✅ read_file(path) - 파일 읽기
✅ write_file(path, content) - 파일 쓰기
✅ append_file(path, content) - 파일 추가
✅ exists(path) - 파일 존재 확인
✅ delete_file(path) - 파일 삭제
✅ list_dir(path) - 디렉토리 목록
✅ create_dir(path) - 디렉토리 생성
... (400줄, 15+ 함수)
```

#### Task 2.3: Math 함수 (300줄)
```
함수 목록:
✅ abs(n) - 절댓값
✅ sqrt(n) - 제곱근
✅ pow(base, exp) - 거듭제곱
✅ sin/cos/tan(n) - 삼각함수
✅ log(n) - 로그
✅ floor(n) - 내림
✅ ceil(n) - 올림
✅ round(n) - 반올림
✅ min(a, b) - 최솟값
✅ max(a, b) - 최댓값
... (300줄, 20+ 함수)
```

---

### Week 3: 표준 라이브러리 2단계 (2026-03-16 ~ 2026-03-22)

#### Task 3.1: String 함수 (350줄)
```
함수 목록:
✅ str_len(s) - 문자열 길이
✅ str_upper(s) - 대문자
✅ str_lower(s) - 소문자
✅ str_trim(s) - 공백 제거
✅ str_split(s, delim) - 분할
✅ str_join(arr, delim) - 결합
✅ str_substring(s, start, len) - 부분 문자열
✅ str_index_of(s, substr) - 위치 찾기
✅ str_replace(s, old, new) - 치환
✅ str_contains(s, substr) - 포함 여부
✅ str_starts_with(s, prefix) - 시작 확인
✅ str_ends_with(s, suffix) - 종료 확인
... (350줄, 18+ 함수)
```

#### Task 3.2: File 함수 (400줄)
```
구현:
- FileHandle 구조체 (파일 핸들)
- 파일 읽기/쓰기 모드
- 바이너리/텍스트 모드
- 라인 단위 읽기
- 전체 파일 읽기
- 스트리밍 처리
```

#### Task 3.3: Time 함수 (200줄)
```
함수 목록:
✅ now() - 현재 시간 (밀리초)
✅ now_sec() - 현재 시간 (초)
✅ sleep(ms) - 대기
✅ format_time(ts) - 시간 포맷
... (200줄, 10+ 함수)
```

---

### Week 4: 테스트 & 최적화 (2026-03-23 ~ 2026-03-30)

#### Task 4.1: 종합 테스트 (900줄)
```
테스트 커버리지:
✅ VM 기본 연산 (100줄)
✅ 스택 연산 (100줄)
✅ 함수 호출 (150줄)
✅ 메모리 할당/해제 (150줄)
✅ Core 함수 (150줄)
✅ IO 함수 (150줄)
✅ Math 함수 (100줄)
✅ String 함수 (100줄)
```

#### Task 4.2: 성능 최적화
```
목표:
✅ VM 시작: <100ms
✅ 함수 호출: <10ns
✅ 메모리 오버헤드: <1MB 베이스
✅ 캐시 친화적 설계

최적화 기법:
- inline 함수
- 범위 기반 루프
- SIMD 가능성 (미래)
- 메모리 풀 (선택)
```

#### Task 4.3: 문서 작성 (500줄)
```
산출물:
✅ RUNTIME_SPECIFICATION.md (1,000줄)
  - Bytecode 명세
  - API 문서
  - 내부 구조

✅ STDLIB_REFERENCE.md (500줄)
  - 모든 함수 레퍼런스
  - 예제 코드
  - 성능 특성

✅ RUNTIME_IMPLEMENTATION.md (300줄)
  - 구현 가이드
  - 설계 결정 사항
  - 알려진 제한사항
```

---

## 📦 **산출물**

### 코드
```
런타임 구현:        3,700줄 (Rust)
표준 라이브러리:    2,750줄 (50+ 함수)
테스트:             900줄
총계:               7,350줄
```

### 문서
```
설명서:             1,800줄
API 레퍼런스:       500줄
예제:               300줄
총계:               2,600줄
```

### 성능 목표
```
✅ VM 시작: <100ms
✅ 함수 호출: <10ns
✅ 메모리: <1MB 베이스
✅ 테스트: >90초 (전체 스위트)
```

---

## 🎯 **핵심 체크포인트**

### Week 1 완료 (2026-03-08)
```
✅ Bytecode 정의 완료
✅ VM 엔진 작동
✅ 기본 연산 (Add, Sub 등) 통과
✅ 함수 호출 틀 완성
```

### Week 2 완료 (2026-03-15)
```
✅ Core 함수 30+ 개 구현
✅ IO 함수 15+ 개 구현
✅ Math 함수 20+ 개 구현
✅ 기본 테스트 50+ 개 통과
```

### Week 3 완료 (2026-03-22)
```
✅ String 함수 18+ 개 구현
✅ File 함수 10+ 개 구현
✅ Time 함수 10+ 개 구현
✅ 통합 테스트 75+ 개 통과
```

### Week 4 완료 (2026-03-30) - Phase B 종료
```
✅ 총 함수 50+ 개 구현 완료
✅ 900줄 테스트, >90% 커버리지
✅ 문서 2,600줄 완성
✅ 성능 목표 달성
✅ 모든 체크포인트 통과
```

---

## 🚀 **성공 기준**

### 기술적 기준
```
✅ cargo build 성공 (0 경고)
✅ cargo test 100% 통과
✅ 함수 50+ 개 구현
✅ 커버리지 >90%
✅ 성능 목표 달성
```

### 품질 기준
```
✅ 코드 품질: clippy 통과
✅ 문서: rustdoc 완성
✅ 에러 처리: 모든 경로 커버
✅ 안전성: unsafe 블록 최소화
```

### 통합 기준
```
✅ v2-freelang-ai와 통합 가능
✅ 은행 시스템 예제 실행
✅ 추가 예제 3+ 개 동작
```

---

## 📌 **다음 단계 (Phase C)**

Phase B 완료 후 Phase C로 진행:

### Phase C: 배포 준비 (4주)
```
주제: 패키지 관리자 + 공식 릴리즈
기한: 2026-04-27
목표: v1.0.0 공식 출시

작업:
1. 패키지 빌드 (cargo build --release)
2. 설치 스크립트 작성
3. 바이너리 배포 (Linux, macOS, Windows)
4. 문서 정리 및 영어 번역
5. 첫 공식 릴리즈 (v1.0.0)
```

---

## 💡 **설계 원칙**

### 1. 단순성
- 복잡한 최적화 피하기
- 명확한 인터페이스
- 문서화된 구조

### 2. 성능
- 프로파일링 기반 최적화
- 메모리 효율성
- 예측 가능한 성능

### 3. 안전성
- Rust의 메모리 안전성 활용
- 철저한 에러 처리
- Panic 최소화

### 4. 확장성
- 플러그인 아키텍처 (미래)
- 표준 함수 추가 용이
- 외부 라이브러리 통합 가능

---

## 📋 **일일 진행 보고서**

매일 10:00에 작업 진행 상황 보고:
```
일자: 2026-03-XX
- 완료: [작업명] (XXX줄)
- 진행중: [작업명] (XX%)
- 이슈: [문제점]
- 다음: [내일 계획]
```

---

**상태**: 📋 준비 완료
**시작 예정**: 2026-03-02 (오늘)
**완료 목표**: 2026-03-30
**담당**: Claude (자동 구현)

🚀 **Phase B 시작 준비 완료!**
