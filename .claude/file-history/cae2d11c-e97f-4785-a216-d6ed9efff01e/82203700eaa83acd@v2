# FreeLang Phase B - Rust 런타임 설계

**작성일**: 2026-03-02
**Phase**: B (런타임 구현)
**기간**: 4주 (2026-03-02 ~ 2026-03-30)
**목표**: 완전히 독립적인 Rust 기반 런타임 + 50+ 표준 함수

---

## 🎯 **Phase B 목표**

### 핵심 목표
```
✅ Rust 자체 런타임 구현 (외부 의존성 최소)
✅ 메모리 관리 시스템 (GC + Reference Counting)
✅ 표준 함수 라이브러리 (50+ 함수)
✅ 성능 최적화 (<100ms 시작, <10ms 실행)
```

### 완성도 목표
```
런타임 구현:    100% (Full Rust)
표준 라이브러리: 80% (Core functions)
성능 테스트:     100% (모든 케이스)
문서:           100% (API + 예제)
```

---

## 📐 **Rust 런타임 아키텍처**

### 1. **계층 구조** (4-Layer)

```
┌─────────────────────────────────────────┐
│  Layer 4: Standard Library (50+ funcs)  │
│  (I/O, String, Array, Math, System)     │
├─────────────────────────────────────────┤
│  Layer 3: Runtime Engine                │
│  (Value System, Execution, Call Stack)  │
├─────────────────────────────────────────┤
│  Layer 2: Memory Management             │
│  (Allocation, GC, Reference Counting)   │
├─────────────────────────────────────────┤
│  Layer 1: Core Types & Traits           │
│  (Value, Object, Array, String)         │
└─────────────────────────────────────────┘
```

### 2. **핵심 모듈**

```
freelang-runtime (Rust)
├── src/
│   ├── core/
│   │   ├── value.rs          (타입 시스템)
│   │   ├── object.rs         (객체 표현)
│   │   └── traits.rs         (핵심 trait)
│   │
│   ├── memory/
│   │   ├── allocator.rs      (메모리 할당)
│   │   ├── gc.rs             (가비지 컬렉션)
│   │   └── refcount.rs       (참조 카운팅)
│   │
│   ├── runtime/
│   │   ├── vm.rs             (가상 머신)
│   │   ├── executor.rs       (실행 엔진)
│   │   ├── stack.rs          (콜 스택)
│   │   └── builtin.rs        (내장 함수)
│   │
│   ├── stdlib/
│   │   ├── io.rs             (입출력 - 15 funcs)
│   │   ├── string.rs         (문자열 - 18 funcs)
│   │   ├── array.rs          (배열 - 12 funcs)
│   │   ├── math.rs           (수학 - 15 funcs)
│   │   ├── system.rs         (시스템 - 8 funcs)
│   │   ├── crypto.rs         (암호 - 8 funcs)
│   │   └── json.rs           (JSON - 5 funcs)
│   │
│   └── main.rs               (진입점)
│
├── Cargo.toml
├── benches/
│   ├── startup.rs            (시작 성능)
│   ├── execution.rs          (실행 성능)
│   └── stdlib.rs             (라이브러리 성능)
│
└── tests/
    ├── unit/                 (단위 테스트)
    ├── integration/          (통합 테스트)
    └── benchmarks/           (성능 테스트)
```

---

## 🔧 **Core Types 설계**

### 1. **Value 타입** (8가지)

```rust
// src/core/value.rs
#[derive(Clone)]
pub enum Value {
    Null,                    // null
    Bool(bool),             // true/false
    Number(f64),            // 정수 & 부동소수점
    String(String),         // 문자열
    Array(Rc<RefCell<Vec<Value>>>),     // 배열
    Object(Rc<RefCell<HashMap<String, Value>>>),  // 객체
    Function(Box<dyn Fn(Vec<Value>) -> Value>),   // 함수
    Error(String),          // 에러
}

// 핵심 메서드
impl Value {
    pub fn is_truthy(&self) -> bool;
    pub fn to_number(&self) -> f64;
    pub fn to_string(&self) -> String;
    pub fn type_of(&self) -> &str;
    pub fn equals(&self, other: &Value) -> bool;
}
```

### 2. **Object 표현**

```rust
// src/core/object.rs
pub struct FLObject {
    id: u64,                           // 고유 ID
    class: String,                     // 클래스명
    fields: HashMap<String, Value>,    // 필드
    methods: HashMap<String, Box<dyn Fn(Vec<Value>) -> Value>>,
    created_at: Instant,
    last_accessed: Instant,
}
```

### 3. **메모리 관리 전략**

```rust
// src/memory/allocator.rs
pub struct MemoryAllocator {
    heap: HashMap<u64, Value>,
    gc_threshold: usize,    // GC 트리거 (5MB)
    gc_interval: Duration,  // 최소 GC 간격 (100ms)
}

// GC Strategy: Mark & Sweep + Reference Counting
// - Reference Counting for immediate cleanup
// - Mark & Sweep for circular references (10초마다)
```

---

## 📚 **표준 함수 라이브러리** (50+ 함수)

### 1. **I/O 함수** (15개)

```
print(msg)          - 콘솔 출력
println(msg)        - 개행 포함 출력
input()             - 콘솔 입력
read_file(path)     - 파일 읽기
write_file(path, content)  - 파일 쓰기
append_file(path, content) - 파일 추가
file_exists(path)   - 파일 존재 확인
delete_file(path)   - 파일 삭제
list_files(dir)     - 디렉토리 목록
mkdir(path)         - 디렉토리 생성
rmdir(path)         - 디렉토리 삭제
get_cwd()           - 현재 디렉토리
cd(path)            - 디렉토리 변경
env(name)           - 환경변수 읽기
system(cmd)         - 시스템 명령 실행
```

### 2. **String 함수** (18개)

```
length(str)         - 문자열 길이
upper(str)          - 대문자 변환
lower(str)          - 소문자 변환
trim(str)           - 공백 제거
split(str, sep)     - 분할
join(arr, sep)      - 결합
substring(str, start, end)  - 부분 추출
contains(str, substr)  - 포함 확인
startswith(str, prefix)  - 시작 확인
endswith(str, suffix)    - 종료 확인
replace(str, old, new)   - 교체
find(str, substr)        - 위치 찾기
reverse(str)             - 역순
repeat(str, count)       - 반복
pad_left(str, len, char) - 좌측 패딩
pad_right(str, len, char)- 우측 패딩
to_upper_case(str)  - 대문자
to_lower_case(str)  - 소문자
```

### 3. **Array 함수** (12개)

```
length(arr)         - 배열 길이
push(arr, val)      - 요소 추가
pop(arr)            - 요소 제거
shift(arr)          - 첫 요소 제거
unshift(arr, val)   - 첫 요소 추가
slice(arr, start, end)  - 부분 추출
concat(arr1, arr2)  - 배열 연결
reverse(arr)        - 역순
sort(arr)           - 정렬
unique(arr)         - 중복 제거
contains(arr, val)  - 포함 확인
index_of(arr, val)  - 인덱스 찾기
```

### 4. **Math 함수** (15개)

```
abs(num)            - 절댓값
round(num)          - 반올림
floor(num)          - 내림
ceil(num)           - 올림
sqrt(num)           - 제곱근
pow(base, exp)      - 거듭제곱
sin(angle)          - 사인
cos(angle)          - 코사인
tan(angle)          - 탄젠트
log(num)            - 자연로그
log10(num)          - 상용로그
exp(num)            - 지수
random()            - 난수 (0-1)
min(a, b)           - 최솟값
max(a, b)           - 최댓값
```

### 5. **System 함수** (8개)

```
now()               - 현재 시간 (ms)
sleep(ms)           - 대기
exit(code)          - 종료
typeof(val)         - 타입 확인
is_null(val)        - null 확인
is_array(val)       - 배열 확인
is_object(val)      - 객체 확인
memory_usage()      - 메모리 사용량
```

### 6. **Crypto 함수** (8개)

```
hash_md5(str)       - MD5 해시
hash_sha256(str)    - SHA256 해시
random_bytes(len)   - 난수 바이트
base64_encode(str)  - Base64 인코딩
base64_decode(str)  - Base64 디코딩
hex_encode(str)     - 16진수 인코딩
hex_decode(str)     - 16진수 디코딩
md5(str)            - MD5 (별칭)
```

### 7. **JSON 함수** (5개)

```
json_parse(str)     - JSON 파싱
json_stringify(obj) - JSON 직렬화
json_pretty(obj, indent)  - JSON 포매팅
json_validate(str)  - JSON 검증
json_minify(str)    - JSON 축소
```

---

## ⚙️ **런타임 실행 프로세스**

### 1. **시작 프로세스** (<100ms)

```
[1] Runtime 초기화 (20ms)
    - 메모리 할당자 생성
    - GC 초기화
    - 전역 심볼 테이블 생성

[2] 표준 라이브러리 로딩 (30ms)
    - 50+ 함수 등록
    - 표준 객체 생성
    - 기본 프로토타입 설정

[3] IR 로드 (20ms)
    - 컴파일된 IR 읽기
    - 코드 객체 생성
    - 함수 등록

[4] 준비 완료
    - 글로벌 컨텍스트 설정
    - 실행 스택 초기화
```

### 2. **실행 프로세스** (<10ms per function)

```
[1] 함수 호출 (0.5ms)
    - 파라미터 평가
    - 새 스택 프레임 생성

[2] 함수 본문 실행 (8ms)
    - IR 명령어 해석
    - 값 계산
    - 부작용 수행

[3] 반환 (0.5ms)
    - 반환값 준비
    - 스택 프레임 정리
    - 제어 흐름 복귀

[4] 메모리 정리 (1ms)
    - 지역 변수 정리
    - 참조 카운팅 업데이트
    - 필요시 GC 트리거
```

---

## 📊 **성능 요구사항**

### 1. **시작 시간**
```
목표: <100ms
- 런타임 초기화: <20ms
- 표준 라이브러리: <30ms
- IR 로드: <20ms
- 준비: <30ms
```

### 2. **실행 시간**
```
목표: <10ms per function
- 함수 호출: <1ms
- 본문 실행: <8ms
- 반환 + 정리: <1ms
```

### 3. **메모리 사용**
```
목표: <10MB base
- 런타임: 2MB
- 표준 라이브러리: 3MB
- 할당 여유: 5MB
```

### 4. **처리량**
```
목표: 10,000+ function calls/sec
- Simple functions: 100,000+ calls/sec
- Complex functions: 1,000+ calls/sec
```

---

## 🗂️ **주간 계획**

### **Week 1: Core + Memory (3/2-3/8)**
- [ ] Rust 프로젝트 초기화
- [ ] Value, Object 타입 정의
- [ ] MemoryAllocator 구현
- [ ] Reference Counting 구현
- [ ] Mark & Sweep GC 구현
- [ ] 단위 테스트 (50개)

**목표**: Core 타입 + 메모리 관리 100% 완성

### **Week 2: Runtime Engine (3/9-3/15)**
- [ ] ExecutionContext 구현
- [ ] CallStack 구현
- [ ] IR 해석기 (Interpreter)
- [ ] 함수 호출 메커니즘
- [ ] 표준 라이브러리 기본 구조
- [ ] 통합 테스트 (30개)

**목표**: 기본 실행 엔진 100% 완성

### **Week 3: Standard Library (3/16-3/22)**
- [ ] I/O 함수 (15개)
- [ ] String 함수 (18개)
- [ ] Array 함수 (12개)
- [ ] Math 함수 (15개)
- [ ] System 함수 (8개)
- [ ] 함수 테스트 (68개)

**목표**: 표준 함수 80% 완성

### **Week 4: Crypto + JSON + Optimization (3/23-3/30)**
- [ ] Crypto 함수 (8개)
- [ ] JSON 함수 (5개)
- [ ] 성능 벤치마크
- [ ] 최적화
- [ ] 문서 작성
- [ ] 통합 테스트 (20개)

**목표**: 전체 런타임 100% 완성, 성능 최적화

---

## 📈 **성공 기준**

### 1. **기술 기준**
```
✅ 모든 50+ 함수 구현 완료
✅ GC 안정성 (no memory leak)
✅ 성능 요구사항 달성
✅ 단위 테스트 100+ 개
✅ 통합 테스트 30+ 개
```

### 2. **코드 품질**
```
✅ 코드 커버리지: >90%
✅ 린트 에러: 0개
✅ 타입 에러: 0개
✅ Clippy 경고: 0개
```

### 3. **문서**
```
✅ API 문서: 완전
✅ 예제: 50+ 개
✅ 성능 리포트: 상세
✅ 아키텍처: 설명
```

---

## 🔍 **의존성 전략** (최소화)

```toml
[dependencies]
# Core only (필수)
parking_lot = "0.12"    # Mutex (std보다 빠름)

# Optional
serde = { version = "1.0", optional = true }
serde_json = { version = "1.0", optional = true }

[dev-dependencies]
criterion = "0.5"       # 성능 벤치마크
proptest = "1.0"        # 속성 기반 테스트
```

---

## 🎯 **다음 단계**

### 즉시 (3/2-3/3)
1. [ ] Rust 프로젝트 초기화
2. [ ] 기본 파일 구조 생성
3. [ ] Cargo.toml 설정
4. [ ] Value 타입 정의 시작

### 단기 (3/4-3/8)
1. [ ] Core 타입 100% 완성
2. [ ] Memory allocator 기본 구현
3. [ ] 첫 테스트 실행

### 중기 (3/9-3/15)
1. [ ] 런타임 엔진 기본 완성
2. [ ] IR 해석 가능

### 장기 (3/16-3/30)
1. [ ] 표준 라이브러리 완성
2. [ ] 성능 최적화
3. [ ] Phase B 완료

---

## 📝 **추적 항목**

- [ ] Rust 프로젝트 생성
- [ ] GitHub/GOGS 저장소
- [ ] CI/CD 파이프라인
- [ ] 성능 벤치마크
- [ ] 커버리지 리포트
- [ ] 릴리즈 준비

---

**상태**: 🟢 설계 완료, 구현 준비됨
**분위기**: 🚀 전진 상태
**목표**: 📈 명확함
**기한**: ⏰ 4주

완벽한 독립 런타임을 만들자! 🎯
