# Phase 6: ML Online Learning Design

## Context

Phase 5의 **분산 시스템** 위에서 Neural-Kernel-Sentinel의 신경망을 **실시간으로 학습**시킵니다.

**목표**:
- Stochastic Gradient Descent (SGD) 구현
- 새로운 공격 패턴 자동 학습 (<1일)
- Model drift detection (KL divergence)
- Adaptive learning rate scheduling

---

## Architecture Overview

```
                    ┌──────────────────────────────┐
                    │  Distributed NKS (Phase 5)   │
                    │  (3-5 nodes, 99.99% SLA)    │
                    └──────────────┬───────────────┘
                                   │
                    ┌──────────────▼───────────────┐
                    │   ML Online Learning (Phase 6)│
                    └──────────────┬───────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
    ┌───▼────┐              ┌─────▼─────┐            ┌──────▼──┐
    │Gradient │              │   Online  │            │  Model  │
    │ Descent │              │ Learning  │            │Evaluation│
    │ (SGD)   │              │ Pipeline  │            │(Drift)   │
    └────┬────┘              └─────┬─────┘            └──────┬──┘
         │                         │                        │
         ├─────────────────────────┼────────────────────────┤
         │                         │                        │
         │  Mini-batch SGD         │  Incremental Updates   │
         │  • Forward Pass         │  • Weight Updates      │
         │  • Backward Pass        │  • Batch Norm          │
         │  • Weight Update        │  • Regularization      │
         │                         │                        │
         │  Learning Rate          │  Drift Detection       │
         │  • Fixed               │  • KL Divergence       │
         │  • Decay               │  • Population Stats     │
         │  • Adaptive            │  • Threshold Alerts     │
         │                         │                        │
         └─────────────────────────┴────────────────────────┘
                             │
                   ┌─────────▼──────────┐
                   │ All 8 Threat Types │
                   │ Update in Parallel │
                   └────────────────────┘
```

---

## 4개 모듈 설계

### 1. `gradient_descent.fl` (Day 1-2, ~350줄)

**목표**: Stochastic Gradient Descent 구현

**핵심 구조체**:
```rust
pub struct SGDOptimizer {
    pub learning_rate: f64,
    pub momentum: f64,  // 0.9 (관성)
    pub decay: f64,     // 0.001 (L2 정규화)
    pub batch_size: u32,
    pub iteration: u64,
}

pub struct NeuralNetworkWeights {
    pub weights_l1: Vec<f64>,  // 16×32 = 512
    pub bias_l1: Vec<f64>,     // 32
    pub weights_l2: Vec<f64>,  // 32×16 = 512
    pub bias_l2: Vec<f64>,     // 16
    pub weights_l3: Vec<f64>,  // 16×8 = 128
    pub bias_l3: Vec<f64>,     // 8
    pub gradients: Gradients,
}

pub struct Gradients {
    pub dw_l1: Vec<f64>,
    pub db_l1: Vec<f64>,
    pub dw_l2: Vec<f64>,
    pub db_l2: Vec<f64>,
    pub dw_l3: Vec<f64>,
    pub db_l3: Vec<f64>,
    pub loss: f64,
}

pub struct TrainingBatch {
    pub inputs: Vec<Vec<f64>>,     // 배치 크기 × 16
    pub targets: Vec<Vec<f64>>,    // 배치 크기 × 8
}
```

**주요 함수**:

1. **`forward_pass(inputs: &[f64]) -> Vec<f64>`**
   - Layer 1: `z1 = input @ W1 + b1`, `a1 = ReLU(z1)`
   - Layer 2: `z2 = a1 @ W2 + b2`, `a2 = ReLU(z2)`
   - Layer 3: `z3 = a2 @ W3 + b3`, `a3 = Softmax(z3)`
   - 반환: 8차원 확률 벡터

2. **`backward_pass(output: &[f64], target: &[f64])`**
   - Loss 계산: `L = -Σ(target * log(output))`
   - Gradient 계산: 역전파 (backpropagation)
   - `dL/dW = (1/m) * Σ(dL/dA @ A.T) + λW` (L2 정규화)

3. **`update_weights(gradients: &Gradients, lr: f64)`**
   - Momentum: `v = momentum × v_prev - lr × gradient`
   - Weight update: `W = W + v`
   - Decay: `W = W × (1 - decay × lr)`

4. **`compute_learning_rate(iteration: u64) -> f64`**
   - Fixed LR: `lr = 0.001`
   - Decay schedule: `lr = 0.001 / (1 + 0.0001 × iteration)`
   - Adaptive: `lr = base_lr × sqrt(epsilon / (m_squared + epsilon))`

5. **`train_batch(batch: &TrainingBatch) -> f64`**
   - Mini-batch SGD (크기 32)
   - 10번 iteration (수렴 가속)
   - 최종 loss 반환

**테스트 A1-A6**:
- A1: Forward pass (정방향 계산 정확성)
- A2: Backward pass (gradient 계산)
- A3: Weight update (수렴 확인)
- A4: Learning rate scheduling
- A5: Mini-batch processing
- A6: Loss 감소 > 10% per batch

---

### 2. `online_learning_pipeline.fl` (Day 3-4, ~350줄)

**목표**: 실시간 가중치 업데이트 및 배치 관리

**핵심 구조체**:
```rust
pub struct OnlineLearningPipeline {
    pub sgd: SGDOptimizer,
    pub weights: NeuralNetworkWeights,
    pub batch_buffer: VecDeque<TrainingBatch>,
    pub batch_size: u32,
    pub update_interval_ms: u32,  // 100ms
    pub samples_seen: u64,
    pub batches_processed: u64,
}

pub struct TrainingSample {
    pub features: Vec<f64>,        // 16차원 (Phase 4 output)
    pub label: u32,                // 0-7 (8가지 위협 타입)
    pub confidence: f64,           // 0.0-1.0
    pub timestamp_ns: u64,
}

pub struct TrainingMetrics {
    pub loss_history: VecDeque<f64>,
    pub accuracy_history: VecDeque<f64>,
    pub throughput: u64,           // samples/sec
    pub update_latency_us: u32,
}
```

**주요 함수**:

1. **`add_training_sample(sample: TrainingSample)`**
   - 새로운 위협 감지 → 즉시 학습 데이터 추가
   - 배치 버퍼에 저장 (FiFO)
   - 정규화: `(features - mean) / std`

2. **`create_mini_batch() -> TrainingBatch`**
   - 32개 샘플로 배치 생성
   - 배치 정규화: `(x - batch_mean) / sqrt(batch_var + eps)`
   - Label one-hot encoding: `[0,0,1,0,0,0,0,0]` (threat_type=2)

3. **`update_weights_online()`** (100ms마다)
   - 대기 중인 배치 처리
   - SGD 학습: 10 iterations
   - 가중치 업데이트 (<100ms)

4. **`get_updated_model() -> NeuralNetworkWeights`**
   - 최신 가중치 반환
   - 모든 노드에 브로드캐스트 (Phase 5)

5. **`compute_batch_statistics(batch: &TrainingBatch)`**
   - Mean/variance 계산
   - Batch normalization 파라미터

**테스트 B1-B6**:
- B1: 샘플 추가 및 버퍼 관리
- B2: Mini-batch 생성 (배치 정규화)
- B3: 온라인 가중치 업데이트 (<100ms)
- B4: 멀티 배치 처리 (연속 학습)
- B5: 정규화 (L2 정규화 효과)
- B6: 처리량 > 1000 samples/sec

---

### 3. `model_evaluation.fl` (Day 5-6, ~300줄)

**목표**: Model drift detection 및 성능 모니터링

**핵심 구조체**:
```rust
pub struct ModelEvaluator {
    pub validation_set: Vec<TrainingSample>,
    pub test_set: Vec<TrainingSample>,
    pub baseline_distribution: Vec<f64>,  // 각 위협 타입의 원래 확률
    pub current_distribution: Vec<f64>,   // 현재 확률
    pub kl_divergence_threshold: f64,     // 0.1 (10% 변화)
}

pub struct PerformanceMetrics {
    pub precision: f64,           // TP / (TP + FP)
    pub recall: f64,              // TP / (TP + FN)
    pub f1_score: f64,            // 2 × (precision × recall) / (precision + recall)
    pub accuracy: f64,            // (TP + TN) / (TP + TN + FP + FN)
    pub auc_roc: f64,            // Area Under ROC Curve
}

pub struct DriftDetectionResult {
    pub is_drifting: bool,
    pub kl_divergence: f64,
    pub affected_classes: Vec<u32>,
    pub severity: f64,            // 0.0 ~ 1.0
    pub recommendation: String,
}
```

**주요 함수**:

1. **`evaluate_on_validation_set() -> PerformanceMetrics`**
   - Validation set (전체의 20%)에서 평가
   - 정확도, 정밀도, 재현율, F1 점수 계산
   - 목표: 정확도 > 99% 유지

2. **`compute_kl_divergence() -> f64`**
   - `KL(P||Q) = Σ P(i) × log(P(i) / Q(i))`
   - P = baseline_distribution (원래)
   - Q = current_distribution (현재)
   - 기준: KL < 0.1 (정상), KL > 0.2 (드리프트 경고)

3. **`detect_drift() -> DriftDetectionResult`**
   - KL divergence 계산
   - 클래스별 확률 변화 분석
   - 심각도: `severity = (kl_divergence - 0.1) / 0.1` (clamped)
   - 권장사항: "Retrain with new data" 또는 "Use conservative policy"

4. **`compare_model_versions(old: &NN, new: &NN) -> f64`**
   - 가중치 변화량 측정
   - `weight_change = Σ |new_w - old_w| / Σ |old_w|`
   - 목표: 변화 < 5% (과도한 변화 방지)

5. **`get_class_distribution(predictions: &[Vec<f64>]) -> Vec<f64>`**
   - 각 위협 타입의 평균 확률
   - `distribution[i] = mean(predictions[*][i])`

**테스트 C1-C6**:
- C1: Validation set 평가
- C2: KL divergence 계산 정확도
- C3: Drift detection (임계값)
- C4: 모델 버전 비교
- C5: 클래스 분포 추적
- C6: Drift 발생 시 경고 발생

---

### 4. `mod.fl` (Day 7, ~250줄)

**목표**: ML 온라인 학습 통합 API

**핵심 구조체**:
```rust
pub struct MLOnlineLearningSystem {
    pub gradient_descent: SGDOptimizer,
    pub pipeline: OnlineLearningPipeline,
    pub evaluator: ModelEvaluator,
    pub learning_enabled: bool,
    pub total_samples_trained: u64,
    pub training_history: Vec<(u64, f64, f64)>,  // timestamp, loss, accuracy
}

pub struct LearningConfig {
    pub learning_rate: f64,
    pub batch_size: u32,
    pub update_interval_ms: u32,
    pub drift_threshold: f64,
}

pub struct TrainingSummary {
    pub samples_trained: u64,
    pub avg_loss: f64,
    pub accuracy: f64,
    pub drift_detected: bool,
    pub time_elapsed_minutes: u64,
    pub new_patterns_learned: Vec<String>,
}
```

**주요 함수**:

1. **`learn_from_threat(threat: ThreatSample)`**
   - Phase 5의 탐지된 위협 → 학습 데이터로 변환
   - 특성 벡터 추출: NKS의 16차원 벡터
   - 라벨 할당: 위협 타입 (0-7)
   - 배치 추가

2. **`update_model_if_ready()`** (100ms마다)
   - 배치 충분히 모였는지 확인
   - SGD 학습 수행
   - 가중치 업데이트

3. **`check_and_handle_drift()`**
   - Model drift 감지
   - Drift 발생 시 보수적 정책 전환 (Phase 5)
   - 추가 학습 강화

4. **`get_learning_progress() -> TrainingSummary`**
   - 학습 진행 상황 반환
   - 학습된 새로운 패턴 목록

5. **`broadcast_updated_weights()`** (Phase 5 연동)
   - 모든 분산 노드에 새로운 가중치 전송
   - 동기화 완료 확인

**테스트 E1-E6**:
- E1: 위협 샘플 학습
- E2: 온라인 모델 업데이트 (<100ms)
- E3: Drift 감지 및 대응
- E4: 학습 진행도 조회
- E5: 가중치 브로드캐스트 (Phase 5)
- E6: 신 패턴 학습 (<1일)

---

## 타이밍 예산 (온라인 학습 < 100ms/batch)

| 단계 | 예산 | 설명 |
|------|------|------|
| Mini-batch 생성 | 10ms | 정규화 + one-hot encoding |
| Forward pass | 20ms | 3-layer NN 추론 |
| Backward pass | 30ms | Gradient 계산 (역전파) |
| Weight update | 20ms | Momentum + decay 적용 |
| Validation (옵션) | 15ms | Accuracy 계산 |
| **합계** | **95ms** | **< 100ms 달성** |

---

## 8개 무관용 규칙

### Rule 1: 학습 가능성 > 95%
```
= 배치 처리 성공률
목표: 95% 이상
```

### Rule 2: SGD 수렴 속도
```
= Loss 감소율 > 10% per batch
목표: 첫 배치부터 수렴
```

### Rule 3: 온라인 업데이트 < 100ms
```
= Mini-batch 처리 시간
목표: 100ms 내 완료
```

### Rule 4: Model drift detection
```
= KL divergence < 0.2 (정상)
목표: 자동 감지, 경고 발생
```

### Rule 5: Validation accuracy ≥ 99%
```
= 모든 위협 타입 정확도
목표: 99% 이상 유지
```

### Rule 6: 신 패턴 학습 < 1일
```
= 새로운 공격 탐지 능력 획득
목표: 하루 이내 자동 학습
```

### Rule 7: 가중치 변화 < 5%
```
= Model stability 측정
목표: 과도한 변화 방지
```

### Rule 8: 병렬 학습 동기화
```
= 모든 분산 노드 가중치 동일
목표: 100ms 내 동기화 (Phase 5)
```

---

## 6가지 실제 환경 시나리오

### Scenario 1: Dirty Pipe 신 변종
```
기존: pipe→write→splice→mmap→execve (Phase 4.5 탐지)
신 변종: pipe→write→splice→mmap→mprotect→exec→syscall chaining
  ↓
모델이 본 적 없는 패턴
  ↓
1. 초기 탐지: confidence = 0.5
2. 학습: 10 batches
3. 다음 변종: confidence = 0.85 ✅
결과: 24시간 내 자동 학습
```

### Scenario 2: Privilege Escalation 새로운 기법
```
기존: ptrace→setuid→chroot→execve
신기법: ptrace→open(/proc/self/cwd)→chroot→setgid→execve
  ↓
학습 전: FN (false negative)
  ↓
1-10 샘플: 학습 시작
11-50 샘플: 정확도 70%
51-100 샘플: 정확도 95% ✅
결과: ~1시간 내 적응
```

### Scenario 3: Zero-day 공격 패턴
```
공격: 예측 불가능한 새로운 기법
  ↓
1. 초기 탐지 (anomaly score 높음): confidence = 0.6
2. 이상 징후 수집 (10-20 샘플)
3. 패턴 학습 (SGD)
4. 신뢰도 향상: confidence = 0.9 ✅
결과: 자동 적응, 수동 개입 불필요
```

### Scenario 4: Model drift
```
공격 환경 변화 → 분포 이동
  ↓
1. KL divergence 모니터링
2. KL > 0.15: 경고
3. 추가 학습 시작
4. KL < 0.1: 정상화 ✅
결과: 자동 drift 복구
```

### Scenario 5: 분산 노드 간 불일치
```
Node 1, 2, 3이 다른 학습 데이터 수집
  ↓
가중치 차이 발생
  ↓
1. 정기적 동기화 (100ms)
2. 평균 가중치 계산
3. 모든 노드에 브로드캐스트
결과: 100% 일관성 유지 ✅
```

### Scenario 6: Continuous learning
```
1주일 운영
  ↓
샘플 누적: 1M+ threat detections
  ↓
학습: 100K batches
  ↓
모델 진화:
  - 초기 정확도: 99.0%
  - 최종 정확도: 99.5% ✅
  - 새로운 위협 유형 10개 추가 학습
```

---

## 구현 일정

- **Day 1-2**: `gradient_descent.fl` 구현 + 테스트 A1-A6
- **Day 3-4**: `online_learning_pipeline.fl` 구현 + 테스트 B1-B6
- **Day 5-6**: `model_evaluation.fl` 구현 + 테스트 C1-C6
- **Day 7**: `mod.fl` 통합 + 6가지 시나리오 + 무관용 규칙 검증

---

## 성공 지표

| 지표 | 목표 | 달성 기준 |
|------|------|----------|
| **학습 가능성** | 95% | 배치 처리 성공률 ≥ 95% |
| **수렴 속도** | >10% loss 감소 | 첫 배치부터 수렴 |
| **온라인 업데이트** | <100ms | Mini-batch 처리 완료 |
| **Drift detection** | KL < 0.2 | 자동 감지 + 경고 |
| **정확도** | ≥ 99% | 모든 위협 타입 유지 |
| **신 패턴 학습** | < 1일 | 자동 학습 완료 |
| **가중치 안정성** | < 5% 변화 | 과도한 변화 방지 |
| **동기화** | 100% | 모든 노드 동일 |

---

## 다음 단계 (Phase 7)

Phase 7부터는 **프로덕션 배포**를 구현합니다:
- Kubernetes 컨테이너화
- 자동 배포 파이프라인 (CI/CD)
- 모니터링 & 알림 시스템
- 성능 프로파일링

목표: **1초 내 자동 배포**, **99.99% SLA 유지**
