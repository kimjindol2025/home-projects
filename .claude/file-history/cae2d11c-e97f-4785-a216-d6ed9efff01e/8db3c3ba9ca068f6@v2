# FreeLang IR v1.0 - Intermediate Representation Specification

## 📋 개요

**IR v1.0** (Intermediate Representation v1.0)은 FreeLang의 **Kernel Task 변환 계층**으로, 고수준 FreeLang 프로그램을 저수준 Kernel Task로 변환하는 중간 표현입니다.

### 역할
```
FreeLang Source Code
    ↓
IR Compiler (이 명세서)
    ↓
Kernel Task (interpreter_v2.fl)
    ↓
Process Scheduler & Neural Execution
```

---

## 🎯 핵심 구조

### IR_Instruction

모든 IR의 기본 단위는 **IR_Instruction**입니다:

```freeling
type IR_Instruction {
    opcode: string            // 명령어 종류
    operands: list            // 명령어 피연산자
    metadata: map             // 추가 정보
}
```

### Metadata 구조

```freeling
// 기본 메타데이터
{
    "id": "IR_001",           // 고유 ID
    "lineNumber": 42,         // 소스 행번호
    "priority": 1,            // 우선순위 (1=highest, 10=lowest)
    "timeout": 1000,          // 실행 타임아웃 (ms)
    "retries": 3,             // 재시도 횟수
    "neuronGroup": "compute", // 신경망 그룹
}
```

---

## 📝 Opcode 정의

### Category 1: Data Movement (데이터 이동)

#### LOAD (값 로드)
```
opcode: "LOAD"
operands: [address, value_type, initial_value]
metadata: { "scope": "local|global", "size": bytes }
```
**예제**:
```
LOAD <pid=1> <type=int> <value=42>
```

#### STORE (값 저장)
```
opcode: "STORE"
operands: [address, source_value]
metadata: { "scope": "local|global" }
```
**예제**:
```
STORE <pid=1> <value=100>
```

#### MOVE (데이터 이동)
```
opcode: "MOVE"
operands: [source, destination]
metadata: { "copySize": bytes }
```

---

### Category 2: Arithmetic (산술 연산)

#### ADD (덧셈)
```
opcode: "ADD"
operands: [result_addr, operand1, operand2]
metadata: { "type": "int|float" }
```

#### MUL (곱셈)
```
opcode: "MUL"
operands: [result_addr, operand1, operand2]
metadata: { "type": "int|float" }
```

#### SUB (뺄셈)
```
opcode: "SUB"
operands: [result_addr, operand1, operand2]
metadata: { "type": "int|float" }
```

#### DIV (나눗셈)
```
opcode: "DIV"
operands: [result_addr, operand1, operand2]
metadata: { "type": "int|float" }
```

---

### Category 3: Control Flow (제어 흐름)

#### JUMP (무조건 점프)
```
opcode: "JUMP"
operands: [target_instruction_id]
metadata: { "label": "loop_start" }
```

#### BRANCH (조건부 분기)
```
opcode: "BRANCH"
operands: [condition, true_target, false_target]
metadata: { "compareOp": "==|!=|<|>|<=|>=" }
```

#### CALL (함수 호출)
```
opcode: "CALL"
operands: [function_id, arg1, arg2, ...]
metadata: { "argCount": N, "returnAddress": id }
```

#### RETURN (함수 반환)
```
opcode: "RETURN"
operands: [return_value]
metadata: {}
```

---

### Category 4: Process Management (프로세스 관리)

#### CREATE_PROCESS (프로세스 생성)
```
opcode: "CREATE_PROCESS"
operands: [pid, process_name, priority]
metadata: { "initialContext": {...}, "cpuAffinity": cores }
```

#### SUSPEND_PROCESS (프로세스 일시 중지)
```
opcode: "SUSPEND_PROCESS"
operands: [pid, reason]
metadata: { "timeout": 1000 }
```

#### RESUME_PROCESS (프로세스 재개)
```
opcode: "RESUME_PROCESS"
operands: [pid]
metadata: {}
```

#### TERMINATE_PROCESS (프로세스 종료)
```
opcode: "TERMINATE_PROCESS"
operands: [pid, exit_code]
metadata: { "cleanup": true }
```

---

### Category 5: Neural Integration (신경망 통합)

#### REGISTER_NEURON (뉴런 등록)
```
opcode: "REGISTER_NEURON"
operands: [neuron_id, neuron_type, connections]
metadata: { "weights": [0.1, 0.2, ...], "activation": "sigmoid|relu|linear" }
```

#### FIRE_NEURON (뉴런 발화)
```
opcode: "FIRE_NEURON"
operands: [neuron_id, signal_data, strength]
metadata: { "propagationMode": "sync|async" }
```

#### ATTACH_PROCESS_TO_NEURON (프로세스-뉴런 연결)
```
opcode: "ATTACH_PROCESS_TO_NEURON"
operands: [pid, neuron_id]
metadata: { "direction": "input|output|bidirectional" }
```

---

### Category 6: Memory Operations (메모리 연산)

#### ALLOCATE_MEMORY (메모리 할당)
```
opcode: "ALLOCATE_MEMORY"
operands: [address, size_bytes]
metadata: { "alignment": 16, "zeroClear": true }
```

#### DEALLOCATE_MEMORY (메모리 해제)
```
opcode: "DEALLOCATE_MEMORY"
operands: [address, size_bytes]
metadata: {}
```

#### COPY_MEMORY (메모리 복사)
```
opcode: "COPY_MEMORY"
operands: [source_addr, dest_addr, size_bytes]
metadata: { "async": false }
```

---

### Category 7: Synchronization (동기화)

#### WAIT (대기)
```
opcode: "WAIT"
operands: [event_id, timeout_ms]
metadata: { "eventType": "neuron_fired|process_ready|custom" }
```

#### SIGNAL (신호 발생)
```
opcode: "SIGNAL"
operands: [event_id, signal_data]
metadata: { "broadcast": false }
```

#### LOCK (뮤텍스 잠금)
```
opcode: "LOCK"
operands: [lock_id]
metadata: { "timeout": 5000 }
```

#### UNLOCK (뮤텍스 해제)
```
opcode: "UNLOCK"
operands: [lock_id]
metadata: {}
```

---

### Category 8: I/O Operations (입출력)

#### READ_INPUT (입력 읽기)
```
opcode: "READ_INPUT"
operands: [destination_addr, input_source]
metadata: { "format": "text|binary|json" }
```

#### WRITE_OUTPUT (출력 쓰기)
```
opcode: "WRITE_OUTPUT"
operands: [source_addr, output_target]
metadata: { "format": "text|binary|json" }
```

#### READ_FILE (파일 읽기)
```
opcode: "READ_FILE"
operands: [file_path, destination_addr]
metadata: { "offset": 0, "size": -1 }
```

#### WRITE_FILE (파일 쓰기)
```
opcode: "WRITE_FILE"
operands: [source_addr, file_path]
metadata: { "append": false, "permissions": "rw-r--r--" }
```

---

### Category 9: Debug & Monitoring (디버그 및 모니터링)

#### DEBUG_PRINT (디버그 출력)
```
opcode: "DEBUG_PRINT"
operands: [message, debug_level]
metadata: { "level": "TRACE|DEBUG|INFO|WARN|ERROR" }
```

#### COLLECT_METRICS (메트릭 수집)
```
opcode: "COLLECT_METRICS"
operands: [metric_scope]
metadata: { "scope": "process|neuron|system", "outputAddr": address }
```

#### SET_BREAKPOINT (중단점 설정)
```
opcode: "SET_BREAKPOINT"
operands: [breakpoint_id, condition]
metadata: { "condition": "value > 100" }
```

---

### Category 10: Utility (유틸리티)

#### NO_OP (아무것도 하지 않음)
```
opcode: "NO_OP"
operands: []
metadata: {}
```

#### HALT (정지)
```
opcode: "HALT"
operands: [exit_code]
metadata: { "reason": "normal|error|emergency" }
```

#### TIME_MEASURE (시간 측정)
```
opcode: "TIME_MEASURE"
operands: [label, operation]
metadata: { "unit": "ms|us|ns" }
```

---

## 🔄 Execution Model

### IR Instruction 실행 순서

```
1. Instruction Fetch (IR에서 명령어 가져오기)
2. Operand Resolution (피연산자 해석)
3. Metadata Apply (메타데이터 적용)
4. Kernel Task Creation (커널 태스크 생성)
5. Task Execution (프로세스 스케줄러에서 실행)
6. Result Collection (결과 수집)
```

### State Machine

```
[NEW] ──fetch──> [FETCHED]
                    ↓
              [RESOLVED] ──create──> [TASK_CREATED]
                    ↑                      ↓
              [ERROR] <────────── [EXECUTION]
                                       ↓
                                  [COMPLETED]
```

---

## 📊 타입 시스템

### Supported Data Types

```freeling
// 기본 타입
type BasicType = "int" | "float" | "bool" | "string"

// 복합 타입
type ComplexType = "list" | "map" | "struct"

// 메모리 타입
type MemoryType = "stack" | "heap" | "shared"
```

### Type Conversion Table

| From | To | Cost | Safe |
|------|-----|------|------|
| int | float | low | ✅ |
| float | int | medium | ⚠️ (loss) |
| int | string | low | ✅ |
| string | int | medium | ⚠️ |
| bool | int | low | ✅ |
| any | any | - | ❌ |

---

## 🎛️ Configuration

### IR Execution Configuration

```freeling
type IRConfig {
    maxInstructions: int = 100000,      // 최대 명령어 수
    executionTimeout: int = 30000,      // 실행 타임아웃 (ms)
    memoryLimit: int = 1073741824,      // 메모리 제한 (1GB)
    stackSize: int = 8388608,           // 스택 크기 (8MB)
    enableProfiling: bool = false,      // 프로파일링 활성화
    debugLevel: string = "INFO"         // 디버그 레벨
}
```

---

## 📚 Complete Opcode Reference

| Category | Opcode | Operands | Metadata | Status |
|----------|--------|----------|----------|--------|
| **Data** | LOAD, STORE, MOVE | 2-3 | scope | ✅ |
| **Arithmetic** | ADD, SUB, MUL, DIV | 3 | type | ✅ |
| **Control** | JUMP, BRANCH, CALL, RETURN | 1-3 | label | ✅ |
| **Process** | CREATE/SUSPEND/RESUME/TERMINATE | 1-3 | context | ✅ |
| **Neural** | REGISTER/FIRE/ATTACH | 2-3 | weights | ✅ |
| **Memory** | ALLOCATE/DEALLOCATE/COPY | 2-3 | align | ✅ |
| **Sync** | WAIT, SIGNAL, LOCK, UNLOCK | 1-2 | event | ✅ |
| **I/O** | READ/WRITE (INPUT/FILE) | 2 | format | ✅ |
| **Debug** | DEBUG_PRINT, METRICS, BREAKPOINT | 1-2 | level | ✅ |
| **Utility** | NO_OP, HALT, TIME_MEASURE | 0-2 | reason | ✅ |

---

## 🔗 Kernel Integration

### IR to kernel_core.Task Mapping

```
IR_Instruction
    ↓
Metadata Extraction
    ↓
Opcode Classification
    ↓
Task Parameters Assignment
    ↓
kernel_core.Task Creation
    ↓
Scheduler Enqueue
```

---

## ✅ Version History

| Version | Date | Status | Changes |
|---------|------|--------|---------|
| v1.0 | 2026-03-03 | Released | Initial specification with 30+ opcodes |
| v1.1 | TBD | Planning | Vector operations, Custom instructions |
| v2.0 | TBD | Planning | Async/await, Distributed IR |

---

## 🎓 Examples

### Example 1: Simple Arithmetic

```
IR_001: LOAD <pid=1> <int> <42>
IR_002: LOAD <pid=1> <int> <8>
IR_003: ADD <pid=1> <42> <8>
IR_004: STORE <pid=1> <50>
IR_005: DEBUG_PRINT "Result: 50" INFO
IR_006: HALT 0
```

### Example 2: Process Creation & Execution

```
IR_001: CREATE_PROCESS <pid=2> "Worker" <priority=5>
IR_002: LOAD <pid=2> <int> <100>
IR_003: REGISTER_NEURON <nid=1> <INPUT> <connections=[2]>
IR_004: FIRE_NEURON <nid=1> <{"value":100}> <0.9>
IR_005: WAIT <neuron_fired> <5000>
IR_006: COLLECT_METRICS <system>
IR_007: HALT 0
```

### Example 3: Memory Operations

```
IR_001: ALLOCATE_MEMORY <0x1000> <1024>
IR_002: LOAD <0x1000> <int> <42>
IR_003: COPY_MEMORY <0x1000> <0x2000> <1024>
IR_004: DEALLOCATE_MEMORY <0x1000> <1024>
IR_005: HALT 0
```

---

## 📖 Design Principles

1. **Simplicity**: 최소 30개의 필수 Opcode만 정의
2. **Composability**: 모든 복잡한 연산은 기본 Opcode의 조합
3. **Traceability**: 모든 Instruction에 ID와 메타데이터 포함
4. **Extensibility**: 새로운 Opcode 추가 용이
5. **Type Safety**: 모든 피연산자 타입 검증
6. **Performance**: 최소 오버헤드로 실행

---

**작성**: Claude Haiku 4.5
**버전**: IR v1.0 (Phase C Week 2)
**상태**: Production Specification ✅
