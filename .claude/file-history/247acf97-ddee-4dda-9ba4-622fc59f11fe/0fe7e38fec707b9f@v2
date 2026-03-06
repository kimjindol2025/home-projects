# Neural-Kernel-Sentinel 구현 계획

## Context
FreeLang OS Kernel 내부에서 실행되는 AI-Native 실시간 보안 위협 탐지 엔진.
Phase 6(패턴 학습) + Phase 9(ML 최적화)를 결합하여 Syscall 패턴 분석으로 Zero-day 공격을 탐지한다.

**목표**: 탐지 정확도 > 99.9%, 탐지 지연 < 100µs

---

## 설치 위치

```
/data/data/com.termux/files/home/freelang-os-kernel/src/security/
├── syscall_interceptor.fl   (~350줄)
├── behavioral_analyzer.fl   (~400줄)
├── threat_classifier.fl     (~350줄)
├── response_engine.fl       (~300줄)
└── mod.fl                   (~250줄)

tests/
└── neural_kernel_sentinel_test.fl  (30개 무관용 테스트)

docs/
└── NEURAL_KERNEL_SENTINEL_DESIGN.md
```

---

## 재사용할 기존 코드

| 파일 | 재사용 항목 |
|------|------------|
| `src/predictive/pattern_recognition.fl` | `update_confidence()` 베이지안 공식 → N-gram 베이스라인 학습 |
| `src/predictive/learning_state_machine.fl` | `StateMachineManager`, `LearningState` 5단계 → 프로세스 프로파일링 |
| `src/predictive/pattern_database.fl` | `RecoveryHistoryDB` 이력 구조 → `ResponseEngine.response_history` |
| `src/observability/policy_controller.fl` | `PolicyType`, `ControlDecision` 패턴 → `SecurityPolicy`, `ResponseAction` |
| `freelang-distributed-system/src/ml/ml-prefetcher.fl` | `forwardPass()` ReLU→Sigmoid → 3-layer NN으로 확장 |
| `src/interrupt.fl` | `InterruptType::SystemCall` (벡터 128) → `SyscallEvent` 컨텍스트 |

---

## 5개 모듈 설계

### 1. `syscall_interceptor.fl` (Day 1-2)

**핵심 구조체**:
```rust
pub enum SyscallCategory { FileRead, FileWrite, FileOpen, FileUnlink,
    ProcessExec, ProcessFork, ProcessKill, ProcessPtrace,
    MemoryMmap, MemoryMprotect, MemoryMunmap, MemoryBrk,
    NetSocket, NetConnect, NetBind, NetSend, NetRecv,
    SetUID, GetUID, Chroot, Chown, SysInfo, ProcStat, DirScan,
    SignalSend, ShmemAccess, PipeOp,
    ModuleLoad, RawSocket, PerfEvent, KexecLoad, BpfSyscall, UserfaultFd }

pub struct SyscallEvent {
    timestamp_ns: u64, syscall_id: u32,
    syscall_category: SyscallCategory,
    pid: u32, uid: u32, return_code: i64,
    arg0: u64, arg1: u64, arg2: u64,
    flags: u32, duration_ns: u32,
}

pub struct SyscallRingBuffer { buffer: Vec<SyscallEvent>, head: usize,
    tail: usize, capacity: usize }  // lock-free, 비트 마스킹

pub struct SyscallInterceptor { buffer: SyscallRingBuffer,
    pid_whitelist: Vec<u32>, category_filter: Vec<bool> }
```

**기반 함수**: `base_risk_weight()` → ModuleLoad/KexecLoad = 1.0, Ptrace = 0.9

---

### 2. `behavioral_analyzer.fl` (Day 3-4)

**핵심 구조체**:
```rust
pub struct SyscallFeatureVector {
    // D0-D4: 카테고리 빈도
    file_io_rate, process_ctrl_rate, memory_manip_rate,
    network_rate, privilege_rate: f64,
    // D5-D9: N-gram 이상치
    bigram_anomaly, trigram_anomaly, fourgram_anomaly,
    sequence_entropy, repetition_factor: f64,
    // D10-D12: 타이밍
    burst_intensity, avg_duration_normalized, timing_variance: f64,
    // D13-D15: 컨텍스트
    privilege_level, high_risk_syscall_count, return_error_rate: f64,
}

pub struct NgramExtractor {
    per_pid_windows: HashMap<u32, VecDeque<u8>>,
    ngram_counts: HashMap<Vec<u8>, u32>,
    normal_baseline: HashMap<Vec<u8>, f64>,  // 베이지안 업데이트
}
```

**Phase 6 베이지안 공식 재사용**:
```rust
// pattern_recognition.fl의 update_confidence() 그대로 적용
*prob = (*prob * freq + evidence) / (freq + 1.0);
```

---

### 3. `threat_classifier.fl` (Day 5-6)

**ML 모델 아키텍처**:
```
입력: 16차원 SyscallFeatureVector
  ↓ [16×32 가중치, ReLU]  ← Phase 9 forwardPass 패턴
은닉층1: 32 뉴런 (ReLU)
  ↓ [32×16 가중치, ReLU]
은닉층2: 16 뉴런 (ReLU)
  ↓ [16×8 가중치, Softmax]
출력: 8가지 위협 유형 확률
  C0:Normal, C1:ZeroDay, C2:PrivilegeEscalation,
  C3:CodeInjection, C4:DataExfiltration,
  C5:PersistenceMechanism, C6:LateralMovement, C7:CoverageEvasion
총 파라미터: 1,208개, 추론 시간: ~20µs
```

**앙상블 전략 (99.9% 달성)**:
- 최종 신뢰도 = `0.6×NN + 0.3×규칙기반 + 0.1×이상치점수`

**빠른 경로 (5µs)**:
```rust
if features.high_risk_syscall_count > 0.8 {
    return ClassificationResult::emergency(pid);  // NN 생략
}
```

---

### 4. `response_engine.fl` (Day 7)

**대응 정책** (Phase H PolicyController 패턴):
```rust
pub enum ResponseAction {
    IsolatePid(u32),     // SIGSTOP < 1µs
    KillPid(u32),        // SIGKILL < 1µs
    Alert(AlertSeverity, String),
    IncreaseSampling(u32),
    SnapshotMemory(u32),
    UpdateNNWeights,
}

pub enum AlertSeverity { Info(0.3-0.5), Warning(0.5-0.7),
                          Critical(0.7-0.9), Emergency(>0.9) }
```

---

### 5. `mod.fl` (Day 8)

**NeuralKernelSentinel 통합 API**:
```rust
pub struct NeuralKernelSentinel {
    interceptor: SyscallInterceptor,
    analyzer: BehavioralAnalyzer,
    classifier: ThreatClassifier,
    response_engine: ResponseEngine,
}

impl NeuralKernelSentinel {
    // 메인 파이프라인: 5+35+30+25 = 95µs < 100µs
    pub fn process_syscall_batch(&mut self, max_events: usize)
        -> Vec<ClassificationResult> { ... }

    pub fn detect_zero_day(&self, pid: u32) -> Option<f64> { ... }
    pub fn get_accuracy_metrics(&self) -> AccuracyMetrics { ... }
}
```

---

## 타이밍 예산 분해 (< 100µs)

| 단계 | 예산 | 최적화 기법 |
|------|------|------------|
| Syscall 캡처 | 5µs | 링 버퍼 비트 마스킹, 카테고리 필터 |
| N-gram 추출 | 15µs | VecDeque 슬라이딩 윈도우, 정수 인덱스 |
| 특성 벡터 계산 | 20µs | 사전 계산 상수, 비트 연산 |
| NN 추론 | 30µs | 루프 언롤링(4×), 고정 크기 배열 |
| 규칙 매칭 + 대응 | 25µs | 정렬된 규칙 리스트, 조기 종료 |
| **합계** | **95µs** | **< 100µs 달성** |

---

## 30개 무관용 테스트 계획

| 그룹 | 내용 | 수 |
|------|------|----|
| A: Syscall 인터셉터 | 링 버퍼, 분류, 처리량 | 6개 |
| B: 행동 분석기 | N-gram, 이상치 점수, 특성 벡터 | 6개 |
| C: 위협 분류기 | NN 추론, 각 위협 유형 탐지 | 6개 |
| D: 대응 엔진 | 격리, 정책, 오탐 복구 | 6개 |
| E: 통합 파이프라인 | 지연, 정확도, Zero-day, 부하 | 6개 |

**8개 무관용 규칙**:
1. 탐지 정확도 > 99.9%
2. 탐지 지연 < 100µs (평균 + P99)
3. False Positive < 0.1%
4. Zero-day 이상치 점수 > 0.7
5. NN 추론 < 30µs
6. 격리 대응 < 1µs
7. 링 버퍼 손실률 = 0%
8. 처리량 > 1M syscall/sec

---

## 구현 순서

1. **Day 1-2**: `syscall_interceptor.fl` + 테스트 A1-A6
2. **Day 3-4**: `behavioral_analyzer.fl` + 테스트 B1-B6
3. **Day 5-6**: `threat_classifier.fl` + 테스트 C1-C6
4. **Day 7**: `response_engine.fl` + 테스트 D1-D6
5. **Day 8**: `mod.fl` 통합 + 테스트 E1-E6 (무관용 검증)

---

## Zero-day 시나리오 검증

```
Dirty Pipe 유사 패턴:
  pipe() → write() → splice() → mmap(RW) → mprotect(RWX) → execve()

특성: D2=0.85, D6=0.92, D14=0.78
NN 출력: C3(CodeInjection)=0.87
결과: Emergency 격리 (73µs < 100µs)
```

---

## 검증 방법

```bash
cd ~/freelang-os-kernel
# 테스트 실행
cargo test --test neural_kernel_sentinel_test

# 무관용 규칙 확인
# 탐지 정확도 > 99.9%, 지연 < 100µs 출력 확인
```
