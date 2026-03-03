# Phase 11A: AI-Driven Auto-Remediation
## 자율 치유 시스템 설계 문서

**상태**: 구현 계획
**기간**: 집중 구현 (4-5시간)
**산출물**: 4개 모듈 + 20개 테스트 + 완료 보고서

---

## 🎯 목표

**이전 Phase 9A+10A**에서 구축한 시스템의 **수동 개입을 최소화**하고 **자동 복구**를 실현하는 것이 목표입니다.

핵심 철학: **"기록이 증명" → "기록이 학습"**
- Phase 9A의 사건 기록
- Phase 10A의 성능 기록
- **이를 학습하여 다음 사건을 예방**

---

## 📦 4개 모듈 상세 설계

### 1️⃣ Incident Pattern Analysis (500줄)
**파일**: `src/ai/incident_pattern_analyzer.fl`

**개념**:
- 과거 사건 데이터로부터 **패턴 학습**
- K-means 클러스터링으로 유사한 사건들 분류
- 각 클러스터의 **근본 원인 식별**
- 새로운 사건과의 **유사도 점수 계산**

**구조체**:
```
struct IncidentPattern
  patternId: string
  clusterCentroid: array<f64>    // 패턴 특성 벡터
  incidentCount: number
  avgResolutionTimeMs: number
  rootCauseKeywords: array<string>
  frequency: f64

struct PatternMatch
  patternId: string
  matchScore: f64                // 0.0-1.0, >0.7 = HIGH confidence
  suggestedRemediations: array<string>
  timeToResolveMs: number

struct IncidentCluster
  clusterId: string
  center: array<f64>
  members: array<string>         // 사건 ID들
  characteristics: map<string, any>
```

**주요 함수**:
```
fn newIncidentPatternAnalyzer() -> map<string, any>
  초기 데이터 구조 생성

fn learnFromIncidents(analyzer, incidents[]) -> Result<array<IncidentPattern>, string>
  - K-means 클러스터링 (k=5)
  - 각 클러스터 특성 계산
  - 근본 원인 키워드 추출

fn matchIncidentPattern(analyzer, currentIncident) -> Result<PatternMatch, string>
  - 유사도 점수 계산
  - 가장 유사한 패턴 찾기
  - 추천 복구 방법 반환

fn calculateIncidentVector(incident) -> array<f64>
  사건을 특성 벡터로 변환:
  [severity(0-1), duration(분), affectedServices(개), errorRate(%), ...]

fn identifyRootCauseKeywords(pattern) -> array<string>
  패턴 내 사건들의 공통 키워드 추출

fn getPatternStats(analyzer) -> string
  패턴 통계 레포트
```

**예시**:
```
사건 1: DB 연결 풀 부족 (10분 해결)
사건 2: 메모리 누수 (15분 해결)
사건 3: DB 연결 풀 부족 (12분 해결)

패턴 분석:
- Cluster 1 (사건 1, 3): "Database Connection Pool"
  rootCause: ["database", "connection", "timeout"]
  avgResolutionTime: 11분

- Cluster 2 (사건 2): "Memory Leak"
  rootCause: ["memory", "heap", "gc"]
  avgResolutionTime: 15분
```

---

### 2️⃣ Auto-Remediation Engine (500줄)
**파일**: `src/ai/auto_remediation_engine.fl`

**개념**:
- 패턴 매칭 결과를 **자동 복구 액션**으로 변환
- 복구 전 **안전성 검증**
- 복구 실행 및 **효과 측정**
- 실패 시 **자동 escalation**

**구조체**:
```
struct RemediationAction
  actionId: string
  name: string                   // "restartDatabasePool", "clearCache", etc
  parameters: map<string, any>
  preconditions: array<string>   // 실행 전 확인사항
  expectedEffectMs: number       // 복구 예상 시간
  riskLevel: string              // "SAFE" / "MODERATE" / "HIGH"

struct RemediationExecution
  executionId: string
  actionId: string
  timestamp: number
  status: string                 // "PENDING" / "EXECUTING" / "SUCCESS" / "FAILED"
  durationMs: number
  effectMeasured: f64            // 0-100%, SLI 개선도
  logs: array<string>

struct RemediationPolicy
  policyId: string
  patternId: string
  actionSequence: array<RemediationAction>
  autoExecute: bool              // 자동 실행 여부
  maxRetries: number
  escalationThreshold: number    // 실패 후 escalate 기준
```

**주요 함수**:
```
fn newAutoRemediationEngine(analyzer) -> map<string, any>
  Pattern Analyzer와 연결하여 초기화

fn registerRemediationPolicy(engine, policy) -> Result<bool, string>
  복구 정책 등록
  예: DB 풀 부족 → ["increasePoolSize", "restartPool", "loadBalance"]

fn matchPatternToActions(engine, patternMatch) -> Result<array<RemediationAction>, string>
  패턴과 매칭되는 복구 액션들 반환

fn validatePreconditions(engine, action) -> Result<bool, string>
  복구 액션 실행 전 안전성 검증
  예: "restartDatabase" 전에 "activeConnections < 100" 확인

fn executeRemediationAction(engine, action) -> Result<RemediationExecution, string>
  복구 액션 실행
  - 로깅
  - 타이밍 측정
  - 성공/실패 기록

fn measureRemediationEffect(engine, execution) -> Result<f64, string>
  복구 후 SLI 개선도 측정 (0-100%)
  예: 복구 전 SLI=95% → 복구 후=99.5% → effect=4.5%

fn rollbackRemediationIfNeeded(engine, execution) -> Result<bool, string>
  복구 실패 시 이전 상태로 롤백

fn escalateToHuman(engine, execution, reason) -> Result<bool, string>
  자동 복구 실패 시 → 운영자에게 alert 발송
```

**복구 액션 예시**:
```
DB 연결 풀 부족:
  Action 1: increasePoolSize(from=20, to=50)
    Precondition: "currentLoad < 80%"
    Risk: SAFE
    ExpectedTime: 5초

  Action 2: restartPoolConnection()
    Precondition: "activeConnections < 10"
    Risk: MODERATE
    ExpectedTime: 15초

캐시 히트율 저하:
  Action 1: warmupCache(keys=["popular_vectors"])
    Risk: SAFE
    ExpectedTime: 2초

  Action 2: increaseMemoryAllocation(from=512MB, to=1GB)
    Risk: MODERATE
    ExpectedTime: 30초
```

---

### 3️⃣ Predictive Scaling (500줄)
**파일**: `src/ai/predictive_scaling.ml`

**개념**:
- **시계열 예측**으로 트래픽 스파이크를 **미리 감지**
- 자동으로 용량 확장 (Pod 추가, 캐시 확대)
- **비용 최적화** (필요할 때만 확장)
- 트래프 패턴 학습 (요일별, 시간대별)

**구조체**:
```
struct TimeSeriesMetric
  timestamp: number
  value: f64
  dayOfWeek: number              // 0=Monday, 6=Sunday
  hourOfDay: number              // 0-23
  isHoliday: bool

struct PredictionModel
  modelId: string
  metricName: string             // "requestsPerSecond", "latencyP99", etc
  trainingDataPoints: number
  accuracy: f64                  // RMSE 또는 MAPE
  lastTrainedAt: number

struct ScalingAction
  actionId: string
  metric: string                 // 트리거 메트릭
  threshold: f64
  scaleUpCount: number           // Pod 추가 개수
  scaleDownCount: number         // Pod 제거 개수
  cooldownMs: number             // 최소 대기 시간
```

**주요 함수**:
```
fn newPredictiveScaler() -> map<string, any>
  초기화

fn trainPredictionModel(scaler, historicalMetrics[]) -> Result<PredictionModel, string>
  과거 메트릭으로 시계열 모델 학습
  - ARIMA(Auto-Regressive Integrated Moving Average)
  - 또는 지수 평활
  - 계절성(요일, 시간대) 고려

fn predictFutureLoad(scaler, hoursAhead: number) -> Result<f64, string>
  n시간 뒤의 예상 트래픽/지연시간 예측
  예: predictFutureLoad(scaler, 1) → 예상 RPS 2000

fn detectAnomalyInForecast(prediction, current) -> Result<bool, string>
  예측과 실제가 크게 벗어나는지 확인
  이상치 감지 = 예측 모델 재학습 필요

fn recommendScalingAction(scaler, prediction) -> Result<ScalingAction, string>
  예측 기반 용량 확장 권장
  예: "CPU 사용률 예상 85% → Pod 2개 추가 추천"

fn executeScalingAction(scaler, action) -> Result<bool, string>
  실제 확장 실행
  - K8s API 호출 (Pod 생성)
  - 또는 클라우드 API (인스턴스 추가)

fn learnSeasonalPatterns(scaler) -> Result<map<string,f64>, string>
  트래픽 패턴 학습
  결과: {
    "monday": 1.2,      // 월요일 1.2배
    "friday": 0.8,      // 금요일 0.8배
    "9am": 1.5,         // 9시 1.5배
    "2am": 0.3          // 2시 0.3배
  }
```

**예측 모델 수식**:
```
지수 평활:
  Ft+1 = α × At + (1-α) × Ft

  At = 실제값
  Ft = 예측값
  α = 평활 계수 (0.2-0.3)

계절 조정:
  Ft_adj = Ft × seasonal_factor[dayOfWeek][hourOfDay]
```

---

### 4️⃣ Intelligent Alert Tuning (500줄)
**파일**: `src/ai/intelligent_alert_tuning.fl`

**개념**:
- 알림 임계값을 **자동으로 조정**하여 거짓 양성(False Positive) 감소
- **Context 기반** 중요도 재계산
- **연관된 알림 그룹화** (연쇄 반응 방지)
- 알림 **소음 감소** (알람 피로 방지)

**구조체**:
```
struct AlertBaseline
  alertType: string              // "HighLatency", "LowCacheHitRate", etc
  normalMin: f64
  normalMax: f64
  stdDev: f64                    // 표준편차
  lastUpdated: number

struct AlertContext
  timestamp: number
  incidentHappening: bool        // 다른 사건 진행 중?
  scalingInProgress: bool        // 확장 중?
  deploymentInProgress: bool     // 배포 중?
  timeOfDay: number
  dayOfWeek: number

struct TuningPolicy
  policyId: string
  baseThreshold: f64
  contextAdjustments: map<string, f64>  // context별 조정값
  falsePositiveThreshold: f64    // 거짓 양성 비율 > 이면 조정
```

**주요 함수**:
```
fn newIntelligentAlertTuner() -> map<string, any>
  초기화

fn calculateDynamicBaseline(tuner, historicalAlerts[]) -> Result<AlertBaseline, string>
  과거 알림 데이터로 정상 범위 학습
  - 평균값 계산
  - 표준편차 계산
  - 이상치 제외

fn adjustAlertThreshold(tuner, context: AlertContext) -> Result<f64, string>
  Context 기반 임계값 동적 조정

  예:
  - 기본 Latency 임계값: 50ms
  - 배포 중: +20ms (100ms)
  - 확장 중: +15ms (65ms)
  - 밤 2시: -5ms (45ms)

fn filterFalsePositives(tuner, alert) -> Result<bool, string>
  알림이 거짓 양성인지 판단

  거짓 양성 조건:
  - 알림 지속시간 < 10초
  - 값이 정상 범위 ±3σ 내
  - 다른 관련 알림 없음

fn groupRelatedAlerts(tuner, alerts[]) -> Result<array<array<string>>, string>
  관련된 알림들을 그룹화

  예:
  - [HighLatency, HighCPU, LowCacheHitRate] → 1개 그룹
    (메모리 부족 사건의 증상들)
  - [DatabaseDown] → 1개 그룹
    (독립적 사건)

fn learnFalsePositivePattern(tuner) -> Result<map<string,f64>, string>
  거짓 양성 비율 계산
  결과: {
    "HighLatency": 0.35,     // 35% 거짓 양성
    "LowCacheHit": 0.02,     // 2% 거짓 양성
  }

  > 0.2 (20%)이면 → 임계값 조정

fn recommendAlertPolicy(tuner, alertType) -> Result<TuningPolicy, string>
  알림 타입별 권장 정책 생성
```

---

## 🧪 Test Suite (20개 테스트)

### Group A: Pattern Analysis (5 tests)
- `testIncidentVectorization` - 사건 → 벡터 변환
- `testKMeansClustering` - 클러스터링
- `testRootCauseIdentification` - 근본 원인 식별
- `testPatternMatching` - 유사도 점수
- `testPatternAccuracy` - 패턴 정확도 (>80%)

### Group B: Auto-Remediation (5 tests)
- `testRemediationPolicyRegistration` - 정책 등록
- `testPreconditionValidation` - 안전성 검증
- `testRemediationExecution` - 복구 실행
- `testEffectMeasurement` - SLI 개선도 측정
- `testEscalationOnFailure` - 실패 시 escalation

### Group C: Predictive Scaling (5 tests)
- `testTimeSeriesModelTraining` - 모델 학습
- `testFutureLoadPrediction` - 트래픽 예측
- `testSeasonalPatternDetection` - 계절 패턴
- `testScalingActionRecommendation` - 확장 추천
- `testAnomalyDetection` - 예측 이상치 감지

### Group D: Alert Tuning (5 tests)
- `testBaselineCalculation` - 정상 범위 계산
- `testDynamicThresholdAdjustment` - 임계값 조정
- `testFalsePositiveDetection` - 거짓 양성 판단
- `testAlertGrouping` - 알림 그룹화
- `testFalsePositiveReduction` - 거짓 양성 비율 감소

---

## 🏗️ 아키텍처 통합

```
Phase 9A (운영 기반)
  ├─ Incident Recording
  └─ SLO/SLI Tracking
     ↓
Phase 11A (AI 자동 복구) ← NEW
  ├─ Pattern Analysis ← Phase 9A 사건 데이터 학습
  ├─ Auto-Remediation ← Pattern 기반 자동 복구
  ├─ Predictive Scaling ← Phase 10A 성능 데이터 학습
  └─ Alert Tuning ← 거짓 양성 감소
     ↓
결과: 자율 치유 시스템
  - 사건 자동 감지 + 복구
  - 미래 트래픽 예상 + 선제 확장
  - 알림 노이즈 제거
```

---

## 📊 성능 목표

| 메트릭 | 목표 |
|--------|------|
| **패턴 인식 정확도** | >80% |
| **자동 복구 성공률** | >85% |
| **복구 시간** | <2분 |
| **예측 정확도 (MAPE)** | <15% |
| **거짓 양성 감소** | >50% |
| **알림 응답 시간** | <30초 |

---

## 📝 구현 순서

1. **incident_pattern_analyzer.fl** (의존성 없음, 기반)
2. **auto_remediation_engine.fl** (Pattern Analyzer 활용)
3. **predictive_scaling.fl** (독립적, 병렬 가능)
4. **intelligent_alert_tuning.fl** (Phase 9A 데이터 활용)
5. **phase_11a_integration_test.fl** (20개 테스트)
6. **PHASE_11A_COMPLETION_REPORT.md**

---

## 🚀 성공 기준

✅ 4개 모듈 완성 (2,000줄)
✅ 20개 통합 테스트 (100% 통과)
✅ 패턴 정확도 >80%
✅ 자동 복구 성공률 >85%
✅ 거짓 양성 50% 감소
✅ GOGS 커밋 3개
✅ 완료 보고서

---

## 🎓 철학

**"기록이 학습한다"**:
- Phase 9A: 사건을 **기록**
- Phase 10A: 성능을 **추적**
- Phase 11A: 이를 **학습**하여 다음 사건 **예방**

결과: 운영자 개입 없이 **자동으로 치유되는 시스템**
