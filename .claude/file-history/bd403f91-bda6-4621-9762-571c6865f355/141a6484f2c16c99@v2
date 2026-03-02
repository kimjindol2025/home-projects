# 🚀 Phase B: FreeLang 런타임 구현 (TypeScript 버전)
**작성일**: 2026-03-02
**언어**: TypeScript (Node.js)
**기간**: 4주 (2026-03-02 ~ 2026-03-30)
**목표**: 고성능 VM + 50+ 표준 함수
**상태**: 🚀 시작 중

---

## 📊 **계획 변경 사항**

### Why TypeScript?
```
✅ 빠른 개발 (Rust 대비 2배 빠름)
✅ Node.js 생태계 활용
✅ 테스트 환경 완비 (Jest)
✅ 디버깅 용이
✅ v2-freelang-ai와 완벽 통합
✅ 성능도 충분함 (<100ms 시작 가능)
```

### 구현 방식
```
freelang-v2/src/runtime/
├── vm.ts                 # VM 핵심 엔진
├── bytecode.ts           # Bytecode 정의
├── interpreter.ts        # 인터프리터
├── memory.ts             # 메모리 관리
├── stdlib/               # 표준 라이브러리
│   ├── core.ts
│   ├── io.ts
│   ├── math.ts
│   ├── string.ts
│   ├── file.ts
│   └── time.ts
└── __tests__/            # Jest 테스트
```

---

## 🏗️ **아키텍처**

### Value 타입 정의
```typescript
// value.ts
export type Value =
  | { type: 'null' }
  | { type: 'bool'; value: boolean }
  | { type: 'int'; value: number }
  | { type: 'float'; value: number }
  | { type: 'string'; value: string }
  | { type: 'array'; value: Value[] }
  | { type: 'object'; value: Map<string, Value> }
  | { type: 'function'; params: string[]; body: Opcode[] }
  | { type: 'reference'; ptr: number };

export class ValueHandler {
  toString(v: Value): string { }
  toNumber(v: Value): number { }
  toBoolean(v: Value): boolean { }
  typeOf(v: Value): string { }
  equals(a: Value, b: Value): boolean { }
}
```

### Bytecode 정의
```typescript
// bytecode.ts
export enum OpcodeType {
  // 스택 연산
  PUSH,
  POP,
  DUP,

  // 산술
  ADD, SUB, MUL, DIV, MOD,

  // 논리
  AND, OR, NOT,

  // 비교
  EQ, NE, LT, LTE, GT, GTE,

  // 제어 흐름
  JUMP,
  JUMP_IF_TRUE,
  JUMP_IF_FALSE,
  CALL,
  RETURN,

  // 변수
  LOAD_LOCAL,
  STORE_LOCAL,
  LOAD_GLOBAL,
  STORE_GLOBAL,

  // 배열/맵
  LOAD_ARRAY,
  STORE_ARRAY,
  LOAD_MAP,
  STORE_MAP,

  // 기타
  PRINT,
  HALT,
}

export interface Instruction {
  op: OpcodeType;
  arg?: any;
  line: number;
  column: number;
}
```

### VM 엔진
```typescript
// vm.ts
export class VirtualMachine {
  private pc: number = 0;                    // 프로그램 카운터
  private stack: Value[] = [];               // 스택
  private memory: Map<number, Value> = new Map(); // 힙
  private variables: Map<string, Value> = new Map(); // 전역 변수
  private functions: Map<string, Function> = new Map(); // 함수
  private callStack: Frame[] = [];           // 호출 스택

  public execute(bytecode: Instruction[]): Value {
    while (this.pc < bytecode.length) {
      const instr = bytecode[this.pc];
      this.executeInstruction(instr);
      this.pc++;
    }
    return this.stack.pop() || { type: 'null' };
  }

  private executeInstruction(instr: Instruction) {
    switch (instr.op) {
      case OpcodeType.PUSH:
        this.stack.push(instr.arg);
        break;
      case OpcodeType.POP:
        this.stack.pop();
        break;
      case OpcodeType.ADD:
        this.binaryOp((a, b) => a + b);
        break;
      // ... 더 많은 연산
    }
  }

  private binaryOp(op: (a: number, b: number) => number) {
    const b = this.popNumber();
    const a = this.popNumber();
    this.stack.push({ type: 'int', value: op(a, b) });
  }

  private popNumber(): number {
    const v = this.stack.pop();
    if (!v || v.type === 'null') throw new Error('Type error');
    if (v.type === 'int' || v.type === 'float') return v.value;
    throw new Error('Expected number');
  }
}
```

---

## 📅 **주간 작업 계획**

### Week 1: 기초 구현 (2026-03-02 ~ 2026-03-08)

#### Task 1.1: Value 타입 & 기본 VM (2일)
```typescript
// src/runtime/value.ts (150줄)
// src/runtime/bytecode.ts (100줄)
// src/runtime/vm.ts (250줄)

산출: 500줄
테스트: PUSH, POP, 기본 산술 연산
```

#### Task 1.2: 인터프리터 & 컴파일러 (2일)
```typescript
// src/runtime/interpreter.ts (300줄)
// src/compiler/astToOpcodes.ts (200줄)

산출: 500줄
기능:
- AST → Bytecode 변환
- 타입 체크
- 에러 보고
```

#### Task 1.3: 메모리 & 변수 관리 (2일)
```typescript
// src/runtime/memory.ts (150줄)
// src/runtime/context.ts (150줄)

산출: 300줄
기능:
- 스택 & 힙 관리
- 지역/전역 변수
- 함수 호출 스택
```

**Week 1 산출**: 1,300줄 (코드 + 테스트)

---

### Week 2: 표준 함수 1단계 (2026-03-09 ~ 2026-03-15)

#### Task 2.1: Core 함수 (500줄)
```typescript
// src/stdlib/core.ts

함수 목록 (30개):
✅ print(value) - 출력
✅ println(value) - 줄바꿈 출력
✅ input() - 사용자 입력
✅ len(arr|str) - 길이
✅ type(value) - 타입
✅ is_null(value) - null 체크
✅ assert(cond, msg) - 단언
✅ error(msg) - 에러 발생
✅ panic(msg) - 패닉
✅ random() - 난수
✅ sleep(ms) - 대기
✅ exit(code) - 종료
✅ toInt(value) - 정수 변환
✅ toFloat(value) - 실수 변환
✅ toString(value) - 문자열 변환
✅ toBool(value) - 불 변환
✅ min(a, b) - 최소값
✅ max(a, b) - 최대값
✅ keys(object) - 키 목록
✅ values(object) - 값 목록
... 10개 더
```

**구현 예시**:
```typescript
export const coreLibrary = {
  print: (vm: VirtualMachine, args: Value[]) => {
    process.stdout.write(
      args.map(v => ValueHandler.toString(v)).join('')
    );
    return { type: 'null' };
  },

  println: (vm: VirtualMachine, args: Value[]) => {
    console.log(args.map(v => ValueHandler.toString(v)).join(''));
    return { type: 'null' };
  },

  len: (vm: VirtualMachine, args: Value[]) => {
    if (args.length !== 1) throw new Error('len() requires 1 argument');
    const v = args[0];
    if (v.type === 'string') {
      return { type: 'int', value: v.value.length };
    } else if (v.type === 'array') {
      return { type: 'int', value: v.value.length };
    }
    throw new Error('len() only works on strings and arrays');
  },

  random: (vm: VirtualMachine, args: Value[]) => {
    return { type: 'float', value: Math.random() };
  },
};
```

#### Task 2.2: Math 함수 (400줄)
```typescript
// src/stdlib/math.ts

함수 목록 (25개):
✅ abs(n) - 절댓값
✅ sqrt(n) - 제곱근
✅ pow(base, exp) - 거듭제곱
✅ sin/cos/tan(n) - 삼각함수
✅ asin/acos/atan(n) - 역삼각함수
✅ log(n) - 자연로그
✅ log10(n) - 10진 로그
✅ exp(n) - e^n
✅ floor(n) - 내림
✅ ceil(n) - 올림
✅ round(n) - 반올림
✅ PI - 원주율
✅ E - 자연상수
... 더 많음
```

#### Task 2.3: IO 함수 (300줄)
```typescript
// src/stdlib/io.ts

함수 목록 (15개):
✅ print(value) - [core에서]
✅ input() - [core에서]
✅ read_stdin() - 표준 입력 전체
✅ write_stdout(str) - 표준 출력
✅ write_stderr(str) - 에러 출력
```

**Week 2 산출**: 1,200줄

---

### Week 3: 표준 함수 2단계 (2026-03-16 ~ 2026-03-22)

#### Task 3.1: String 함수 (400줄)
```typescript
// src/stdlib/string.ts

함수 목록 (20개):
✅ str_len(s) - 길이 (len과 동일)
✅ str_upper(s) - 대문자
✅ str_lower(s) - 소문자
✅ str_trim(s) - 공백 제거
✅ str_ltrim(s) - 좌측 공백 제거
✅ str_rtrim(s) - 우측 공백 제거
✅ str_split(s, delim) - 분할
✅ str_join(arr, delim) - 결합
✅ str_substring(s, start, len) - 부분 추출
✅ str_index_of(s, substr) - 위치 찾기
✅ str_last_index_of(s, substr) - 마지막 위치
✅ str_replace(s, old, new) - 치환
✅ str_replace_all(s, old, new) - 전체 치환
✅ str_contains(s, substr) - 포함 여부
✅ str_starts_with(s, prefix) - 시작 확인
✅ str_ends_with(s, suffix) - 종료 확인
✅ str_pad_left(s, len, char) - 좌측 패딩
✅ str_pad_right(s, len, char) - 우측 패딩
✅ str_repeat(s, count) - 반복
```

#### Task 3.2: File 함수 (350줄)
```typescript
// src/stdlib/file.ts

함수 목록 (12개):
✅ file_exists(path) - 파일 존재 확인
✅ file_read(path) - 전체 읽기
✅ file_read_lines(path) - 라인 단위 읽기
✅ file_write(path, content) - 전체 쓰기
✅ file_append(path, content) - 추가 쓰기
✅ file_delete(path) - 파일 삭제
✅ file_size(path) - 파일 크기
✅ file_created_time(path) - 생성 시간
✅ file_modified_time(path) - 수정 시간
✅ dir_exists(path) - 디렉토리 존재
✅ dir_list(path) - 목록 조회
✅ dir_create(path) - 디렉토리 생성
```

#### Task 3.3: Time 함수 (200줄)
```typescript
// src/stdlib/time.ts

함수 목록 (10개):
✅ now() - 현재 시간 (밀리초)
✅ now_sec() - 현재 시간 (초)
✅ now_ns() - 현재 시간 (나노초)
✅ sleep(ms) - 대기
✅ format_time(ts, fmt) - 포맷팅
✅ parse_time(str, fmt) - 파싱
```

**Week 3 산출**: 950줄

---

### Week 4: 통합 테스트 & 문서 (2026-03-23 ~ 2026-03-30)

#### Task 4.1: 종합 테스트 (800줄)
```typescript
// src/runtime/__tests__/

테스트 스위트:
✅ VM 테스트 (200줄)
   - Push/Pop/Dup
   - 산술 연산 (+, -, *, /, %)
   - 논리 연산 (&&, ||, !)
   - 비교 연산 (==, !=, <, >, <=, >=)

✅ 함수 호출 (150줄)
   - 함수 정의
   - 함수 호출
   - 반환 값
   - 재귀 함수

✅ 변수 관리 (100줄)
   - 지역 변수
   - 전역 변수
   - 변수 스코프

✅ Core 함수 (150줄)
   - print, println
   - len, type
   - type 변환

✅ 표준 함수 (200줄)
   - Math 함수
   - String 함수
   - File 함수

✅ 통합 테스트 (100줄)
   - 은행 시스템 예제
   - 복잡한 프로그램
```

#### Task 4.2: 성능 최적화
```typescript
// src/runtime/optimizations/

최적화 기법:
✅ 인라인 캐싱 (Inline Caching)
✅ 메모리 풀 (Object Pool)
✅ 함수 메모이제이션
✅ 상수 폴딩 (Constant Folding)
```

#### Task 4.3: 문서 작성 (500줄)
```
산출물:
✅ RUNTIME_API.md (400줄)
   - API 레퍼런스
   - 모든 함수 설명
   - 사용 예제

✅ STDLIB_REFERENCE.md (300줄)
   - 표준 함수 목록
   - 상세 설명
   - 예제

✅ IMPLEMENTATION_GUIDE.md (200줄)
   - 아키텍처 설명
   - 확장 방법
   - 성능 팁
```

**Week 4 산출**: 1,500줄

---

## 📦 **최종 산출물**

### 코드
```
Week 1:  1,300줄 (VM 기초)
Week 2:  1,200줄 (함수 1단계)
Week 3:    950줄 (함수 2단계)
Week 4:  1,500줄 (테스트 + 최적화)
────────────────
총합:    4,950줄
```

### 문서
```
API 문서:    400줄
레퍼런스:    300줄
구현 가이드:  200줄
────────────────
총합:        900줄
```

### 테스트
```
단위 테스트:  600줄
통합 테스트:  200줄
────────────────
총합:        800줄
```

### 전체
```
코드:        4,950줄
테스트:      800줄
문서:        900줄
────────────────
전체:        6,650줄
```

---

## 🎯 **주간 체크포인트**

### Week 1 완료 (2026-03-08) ✅
```
목표: VM 기본 작동
기준:
✅ npm test 통과
✅ PUSH, POP, ADD, SUB 작동
✅ 함수 호출 기능
✅ 변수 저장 및 로드
```

### Week 2 완료 (2026-03-15) ✅
```
목표: 기본 함수 30개 구현
기준:
✅ Core 함수 (print, len, type 등)
✅ Math 함수 (abs, sqrt, sin 등)
✅ IO 함수 (input, read 등)
✅ 50개+ 테스트 통과
```

### Week 3 완료 (2026-03-22) ✅
```
목표: 고급 함수 40개 추가
기준:
✅ String 함수 (split, join, upper 등)
✅ File 함수 (read, write, delete 등)
✅ Time 함수 (now, sleep, format 등)
✅ 100개+ 테스트 통과
```

### Week 4 완료 (2026-03-30) ✅
```
목표: Phase B 완료
기준:
✅ 총 함수 50+ 개 완성
✅ 테스트 커버리지 >90%
✅ 문서 완성
✅ 성능 기준 달성
✅ npm test 100% 통과
```

---

## 🚀 **빠른 시작**

### 프로젝트 구조
```bash
freelang-v2/
├── src/
│   ├── runtime/
│   │   ├── value.ts          # Value 타입
│   │   ├── bytecode.ts        # Bytecode 명세
│   │   ├── vm.ts              # VM 엔진
│   │   ├── interpreter.ts     # 인터프리터
│   │   ├── memory.ts          # 메모리 관리
│   │   ├── stdlib/
│   │   │   ├── core.ts        # 기본 함수
│   │   │   ├── math.ts        # 수학 함수
│   │   │   ├── io.ts          # 입출력
│   │   │   ├── string.ts      # 문자열
│   │   │   ├── file.ts        # 파일
│   │   │   └── time.ts        # 시간
│   │   └── __tests__/
│   │       ├── vm.test.ts
│   │       ├── stdlib.test.ts
│   │       └── integration.test.ts
│   └── compiler/
│       └── astToOpcodes.ts    # AST → Bytecode
└── package.json
```

### 초기화 명령
```bash
# v2-freelang-ai로 이동
cd freelang-v2

# 필요한 패키지 설치 (이미 있을 것)
npm install

# 런타임 테스트 실행
npm test -- runtime

# 빌드
npm run build
```

---

## 💡 **설계 원칙**

### 1. 성능
```
✅ <100ms VM 시작
✅ <10ns 함수 호출
✅ <1MB 메모리 오버헤드
```

### 2. 안정성
```
✅ 모든 에러 처리
✅ 타입 체크
✅ 100% 테스트 커버리지
```

### 3. 확장성
```
✅ 함수 추가 용이
✅ 플러그인 아키텍처 (미래)
✅ 명확한 인터페이스
```

### 4. 품질
```
✅ TypeScript strict mode
✅ ESLint + Prettier
✅ 상세한 문서
```

---

## 📊 **성공 기준**

### 기술적
```
✅ npm build: 0 경고
✅ npm test: 100% 통과
✅ 함수: 50+ 개 구현
✅ 커버리지: >90%
```

### 성능
```
✅ VM 시작: <100ms
✅ 함수 호출: <10ns
✅ 메모리: <1MB 베이스
✅ 테스트: <60초
```

### 품질
```
✅ TypeScript 완벽
✅ 문서 완성
✅ 예제 5+ 개
✅ 실행 가능
```

---

## 📌 **다음 단계 (Phase C)**

Phase B 완료 후:

### Phase C: 배포 준비 (4주)
```
기한: 2026-04-27
주제: 패키지 빌드 + 공식 릴리즈
목표: v1.0.0 공식 출시

작업:
1. 번들링 (webpack/esbuild)
2. CLI 도구 (freelang 명령)
3. 문서 영어 번역
4. npm 패키지 등록
5. 첫 공식 릴리즈
```

---

**상태**: 🚀 준비 완료
**시작**: 2026-03-02 (지금)
**완료**: 2026-03-30
**언어**: TypeScript + Node.js

🎯 **Week 1 Task 1.1 시작하시겠습니까?**
