# Phase 6 Track B: Analytics & Monitoring Implementation Plan

**프로젝트**: FreeLang Distributed Vector Index System - Phase 6
**Track**: B (Analytics & Monitoring)
**기간**: 2주 (병렬 진행)
**목표**: Performance Dashboard, Distributed Tracing, Real-time Alerts, Cost Analysis
**상태**: 🚀 시작

---

## 🎯 Track B 구성 (4개 모듈, ~2,700줄)

### 모듈 1: Performance Dashboard (500줄)

**목적**: 실시간 성능 메트릭 시각화 (Prometheus + Grafana)

#### 파일: `src/analytics/performance_dashboard.fl`

```
struct Metric
  name: string
  value: number
  unit: string                 # "ms", "MB", "count", "%"
  timestamp: number
  labels: map<string, string>  # {"method": "POST", "endpoint": "/insert"}

struct MetricTimeSeries
  metric: Metric
  samples: array<tuple<number, number>>  # (timestamp, value) pairs

struct PerformanceMetrics
  // API Performance
  requestLatency: MetricTimeSeries      # P50, P95, P99
  throughput: MetricTimeSeries          # req/sec
  errorRate: MetricTimeSeries           # %

  // System Performance
  cpuUsage: MetricTimeSeries            # %
  memoryUsage: MetricTimeSeries         # MB
  diskIO: MetricTimeSeries              # MB/sec

  // Distributed System
  raftLeaderElections: MetricTimeSeries
  logReplicationLatency: MetricTimeSeries
  quorumAchievementRate: MetricTimeSeries  # %

  // Vector Operations
  vectorInsertLatency: MetricTimeSeries
  vectorSearchLatency: MetricTimeSeries
  vectorSearchThroughput: MetricTimeSeries

  // Cache Performance
  cacheHitRatio: MetricTimeSeries       # %
  cacheEvictions: MetricTimeSeries      # count/sec

  timestamp: number

struct DashboardConfig
  promUrl: string              # "http://localhost:9090"
  grafanaUrl: string           # "http://localhost:3000"
  scrapeInterval: number       # 15 seconds
  retentionDays: number        # 30 days

struct PerformanceDashboard
  config: DashboardConfig
  metrics: PerformanceMetrics
  alertRules: array<AlertRule>
  customDashboards: array<map<string, any>>

fn newPerformanceDashboard(config: DashboardConfig) -> PerformanceDashboard
fn recordMetric(dashboard, metric: Metric) -> Result<bool, string>
fn getMetrics(dashboard, metricName, timeRange) -> Result<array<Metric>, string>
fn getLatencyPercentiles(dashboard, endpoint) -> map<string, number>  # P50, P95, P99
fn calculateAverageLatency(dashboard, startTime, endTime) -> number
fn getErrorRateByEndpoint(dashboard) -> map<string, number>
fn generateDashboardSnapshot(dashboard) -> string  # JSON for Grafana
fn exportMetricsToPrometheus(dashboard) -> Result<string, string>  # Prometheus format
fn getHealthScore(dashboard) -> number  # 0-100
```

**메트릭 정의**:

```
Request Latency:
  P50: 50th percentile (중간값)
  P95: 95th percentile (거의 모든 요청)
  P99: 99th percentile (최악의 경우)

Example: P50=5ms, P95=20ms, P99=100ms
→ 50% 요청은 5ms 이내, 99%는 100ms 이내

Throughput: 초당 처리 요청 수
  Example: 10,000 req/sec

Error Rate: 실패한 요청 비율
  Example: 0.5% (5000 중 25개 실패)
```

**대시보드 종류**:

```
1. System Overview (전체 시스템)
   - CPU, Memory, Disk
   - Network (In/Out)
   - Error Rate

2. API Performance (API별)
   - Latency (P50, P95, P99)
   - Throughput
   - Error Rate by endpoint

3. Distributed System Health
   - Raft status
   - Replica status
   - Shard distribution

4. Vector Operations
   - Insert latency
   - Search latency
   - Search throughput

5. User Activity
   - Active users
   - Requests per user
   - Cost per user
```

---

### 모듈 2: Distributed Tracing (500줄)

**목적**: 요청 전체 경로 추적 (End-to-end visibility)

#### 파일: `src/tracing/distributed_tracing.fl`

```
struct TraceContext
  traceId: string              # Unique ID for entire request
  spanId: string               # Unique ID for this operation
  parentSpanId: string         # Parent operation ID
  samplingRate: number         # 0.0 to 1.0

struct Span
  spanId: string
  traceId: string
  parentSpanId: string
  operationName: string        # "http.request", "db.query", "cache.get"
  serviceName: string          # "api-server", "vector-db", "cache"
  startTime: number
  endTime: number
  duration: number             # milliseconds
  tags: map<string, string>    # {"http.method": "POST", "http.path": "/api/insert"}
  logs: array<LogEntry>        # Debug logs within span
  status: string               # "OK", "ERROR", "TIMEOUT"

struct LogEntry
  timestamp: number
  message: string
  level: string                # "DEBUG", "INFO", "WARN", "ERROR"

struct Trace
  traceId: string
  userId: string
  startTime: number
  endTime: number
  duration: number
  spans: array<Span>
  tags: map<string, string>    # Request-level tags

struct DistributedTracingConfig
  jaegerUrl: string            # "http://localhost:6831"
  serviceName: string          # "freelang-backend"
  samplingRate: number         # 0.1 (10% of requests)
  maxSpansPerTrace: number     # 1000

struct TracingCollector
  config: DistributedTracingConfig
  traces: array<Trace>
  spans: array<Span>

fn newTracingCollector(config: DistributedTracingConfig) -> TracingCollector
fn startSpan(collector, traceContext, operationName, serviceName) -> Span
fn endSpan(collector, span, status) -> Result<bool, string>
fn setSpanTag(span, key, value) -> Span
fn logSpanEvent(span, message, level) -> Span
fn getTrace(collector, traceId) -> Result<Trace, string>
fn getTraceTree(collector, traceId) -> string  # ASCII tree
fn analyzeTraceCriticalPath(collector, traceId) -> array<Span>  # Slowest path
fn getLatencyByService(collector, traceId) -> map<string, number>
fn exportTraceToJaeger(collector, trace) -> Result<string, string>
```

**Trace 예시**:

```
Request: POST /api/vector/insert

Trace ID: trace-12345
├─ Span 1: API Handler (10ms)
│  ├─ [0-2ms] TLS Decrypt
│  ├─ [2-3ms] JWT Validate
│  └─ [3-10ms] Authorization
│
├─ Span 2: Vector Processing (30ms)
│  ├─ [0-5ms] Vector Normalization
│  └─ [5-30ms] LSH Index Lookup
│
├─ Span 3: Raft Consensus (50ms)
│  ├─ [0-10ms] Log Append
│  ├─ [10-40ms] Heartbeat + Replicate
│  └─ [40-50ms] Leader Commit
│
├─ Span 4: Cache Write (5ms)
│  └─ [0-5ms] Redis SET
│
└─ Span 5: Response (5ms)
   └─ [0-5ms] JSON Serialize

Total: ~100ms
Critical Path: Raft (50ms) > Vector (30ms) > API (10ms) > Cache (5ms) > Response (5ms)
→ Raft이 병목, 최적화 여지 있음
```

**서비스 간 추적**:

```
Client
  ↓
API Server (span: api-handler)
  ├─ Security (span: tls, jwt, rbac)
  ├─ Coordinator (span: route-request)
  │  └─ Raft Cluster (span: raft-consensus)
  │     ├─ Node 1 (span: replicate)
  │     ├─ Node 2 (span: replicate)
  │     └─ Node 3 (span: replicate)
  ├─ Cache (span: cache-set)
  └─ Response (span: serialize)

모든 Span이 traceId로 연결 → 전체 경로 가시화
```

---

### 모듈 3: Real-time Alerts (500줄)

**목적**: 이상 감지 + 알림 생성

#### 파일: `src/alerts/anomaly_detection.fl`

```
struct AlertRule
  ruleId: string
  name: string                 # "High Latency Alert"
  metric: string               # "request_latency_p99"
  condition: string            # "> 1000"  (threshold or expression)
  duration: number             # Sustained for N seconds
  severity: string             # "CRITICAL", "HIGH", "MEDIUM", "LOW"
  notificationChannels: array<string>  # ["slack", "email", "pagerduty"]

struct Anomaly
  anomalyId: string
  ruleId: string
  timestamp: number
  metric: string
  currentValue: number
  threshold: number
  deviation: number            # Standard deviations from mean
  explanation: string
  suggestedAction: string

struct AnomalyDetectionConfig
  statisticalMethod: string    # "zscore", "iqr", "isolation_forest"
  baselineWindow: number       # 24 hours (86400 seconds)
  detectionWindow: number      # 5 minutes
  sensitivity: number          # 0.5-2.0 (lower = more sensitive)

struct AnomalyDetector
  config: AnomalyDetectionConfig
  rules: array<AlertRule>
  anomalies: array<Anomaly>
  baselineMetrics: map<string, array<number>>
  notificationHistory: array<map<string, any>>

fn newAnomalyDetector(config: AnomalyDetectionConfig) -> AnomalyDetector
fn addAlertRule(detector, rule: AlertRule) -> Result<bool, string>
fn detectAnomalies(detector, currentMetrics: map<string, number>) -> array<Anomaly>
fn calculateZScore(detector, metric: string, value: number) -> number
fn getBaseline(detector, metric: string) -> map<string, number>  # mean, stddev, min, max
fn sendNotification(detector, anomaly: Anomaly) -> Result<bool, string>
fn acknowledgeAnomaly(detector, anomalyId: string) -> Result<bool, string>
fn getAlertHistory(detector, timeRange) -> array<Anomaly>
fn calculateFalsePositiveRate(detector) -> number
fn updateBaseline(detector, metric: string) -> Result<bool, string>
```

**이상 탐지 방법**:

```
1. Z-Score (Statistical)
   zscore = (value - mean) / stddev

   Example:
   - Mean latency: 50ms
   - StdDev: 10ms
   - Current: 100ms
   - Z-score: (100-50)/10 = 5.0

   Rule: Z-score > 3.0 (99.7% threshold)
   → Alert! Anomaly detected

2. IQR (Interquartile Range)
   Lower = Q1 - 1.5 × IQR
   Upper = Q3 + 1.5 × IQR

   If value < Lower or > Upper → Anomaly

3. Isolation Forest
   Detects outliers by isolation path length
   More robust to non-Gaussian distributions
```

**알림 규칙 예**:

```
Rule 1: High Latency
  metric: request_latency_p99
  condition: > 1000ms
  duration: 60 seconds (sustained)
  severity: HIGH
  channels: [slack, pagerduty]
  Action: "Check Raft replication status"

Rule 2: High Error Rate
  metric: error_rate
  condition: > 5%
  duration: 30 seconds
  severity: CRITICAL
  channels: [slack, email, pagerduty]
  Action: "Check application logs"

Rule 3: Cache Miss Spike
  metric: cache_hit_ratio
  condition: < 60% AND zscore > 2.0
  duration: 120 seconds
  severity: MEDIUM
  channels: [slack]
  Action: "Check cache eviction"
```

---

### 모듈 4: Cost Analysis (500줄)

**목적**: 사용량 기반 청구 계산

#### 파일: `src/billing/cost_analysis.fl`

```
struct PricingTier
  name: string                 # "Free", "Pro", "Enterprise"
  baseCost: number             # Monthly base cost
  vectorUnitPrice: number      # $ per 1000 vectors
  requestPrice: number         # $ per 1000 requests
  storagePrice: number         # $ per GB/month
  cachedRequestPrice: number   # $ per cached request (lower)

struct UsageMetrics
  userId: string
  billingPeriod: string        # "2026-03"
  vectorsCreated: number
  vectorsQueried: number
  totalRequests: number
  cachedRequests: number
  storageUsedGB: number
  peakConcurrentUsers: number

struct CostBreakdown
  userId: string
  billingPeriod: string
  baseCost: number
  vectorCost: number           # (vectors / 1000) × price
  requestCost: number          # ((requests - cached) / 1000) × price
  cachedRequestSavings: number # -(cached / 1000) × price
  storageCost: number
  totalCost: number
  monthlyRecurringRevenue: number  # For all users

struct CostAnalysisConfig
  pricingTiers: map<string, PricingTier>
  maxFreeTierVectors: number   # 1000
  billingDay: number           # 1-28
  currencyCode: string         # "USD"

struct CostAnalyzer
  config: CostAnalysisConfig
  usageMetrics: array<UsageMetrics>
  costBreakdowns: array<CostBreakdown>
  billingRecords: array<map<string, any>>

fn newCostAnalyzer(config: CostAnalysisConfig) -> CostAnalyzer
fn recordUsage(analyzer, userId, usage: UsageMetrics) -> Result<bool, string>
fn calculateCostBreakdown(analyzer, userId, billingPeriod) -> Result<CostBreakdown, string>
fn generateInvoice(analyzer, userId, billingPeriod) -> Result<string, string>  # HTML/PDF
fn estimateNextMonthCost(analyzer, userId) -> Result<number, string>
fn getTopCostUsers(analyzer, limit: number) -> array<tuple<string, number>>  # (userId, cost)
fn analyzeCostTrends(analyzer, userId) -> array<CostBreakdown>  # Last 12 months
fn optimizeForCost(analyzer, userId) -> array<string>  # Recommendations
fn calculateCumulativeRevenue(analyzer, startDate, endDate) -> number
fn getMetricsForBilling(analyzer, userId) -> UsageMetrics
```

**가격 모델**:

```
Free Tier:
  - Base: $0
  - 1,000 vectors included
  - 10,000 requests/month included
  - 1GB storage included

Pro Tier ($99/month):
  - Base: $99
  - $0.10 per 1,000 vectors (overage)
  - $0.01 per 1,000 requests (overage)
  - $0.10 per GB storage (after 100GB)
  - Cache hit saves 50% (cached = $0.005 per 1,000)

Enterprise:
  - Custom pricing (contact sales)
  - Volume discounts
  - Dedicated support
```

**청구서 예**:

```
User: company-xyz
Period: 2026-03
Tier: Pro

Usage:
  Vectors created: 50,000
  Requests: 500,000
  Cached requests: 200,000
  Storage used: 150GB

Breakdown:
  Base cost: $99.00
  Vector overage: (50,000 - 0) / 1,000 × $0.10 = $5.00
  Requests: (500,000 - 10,000 - 200,000) / 1,000 × $0.01 = $2.90
  Cached request savings: -(200,000 / 1,000 × $0.005) = -$1.00
  Storage: (150 - 100) × $0.10 = $5.00

Total: $99.00 + $5.00 + $2.90 - $1.00 + $5.00 = $110.90
```

---

## 📁 파일 구조

```
src/
├─ analytics/
│  └─ performance_dashboard.fl (500줄)
├─ tracing/
│  └─ distributed_tracing.fl (500줄)
├─ alerts/
│  └─ anomaly_detection.fl (500줄)
└─ billing/
   └─ cost_analysis.fl (500줄)

tests/
└─ track_b_integration_test.fl (600줄, 20+ 테스트)
   ├─ Performance Dashboard Tests (5)
   ├─ Distributed Tracing Tests (5)
   ├─ Anomaly Detection Tests (5)
   ├─ Cost Analysis Tests (5)
   └─ Integration Scenarios (3)

docs/
├─ PHASE_6_TRACK_B_PLAN.md (this file)
├─ PHASE_6_TRACK_B_GUIDE.md (500줄)
└─ TRACK_B_IMPLEMENTATION_EXAMPLES.md (400줄)
```

---

## 🧪 테스트 계획 (20+ 테스트)

### Performance Dashboard Tests (5개)
1. `testMetricCollection` - 메트릭 수집
2. `testLatencyPercentiles` - P50, P95, P99 계산
3. `testErrorRateCalculation` - 에러율 계산
4. `testDashboardSnapshot` - 스냅샷 생성
5. `testHealthScore` - 건강도 점수

### Distributed Tracing Tests (5개)
1. `testTraceCreation` - Trace 생성
2. `testSpanHierarchy` - Span 계층
3. `testTracePropagation` - 전체 경로 추적
4. `testCriticalPathAnalysis` - 병목 지점
5. `testLatencyByService` - 서비스별 지연

### Anomaly Detection Tests (5개)
1. `testZScoreDetection` - Z-score 이상 감지
2. `testIQRDetection` - IQR 방식
3. `testBaselineCalculation` - 기준값 계산
4. `testAlertRuleMatching` - 규칙 매칭
5. `testFalsePositiveRate` - 거짓 양성률

### Cost Analysis Tests (5개)
1. `testUsageRecording` - 사용량 기록
2. `testCostBreakdown` - 비용 분해
3. `testInvoiceGeneration` - 청구서 생성
4. `testCostOptimization` - 최적화 제안
5. `testMonthlyRecurringRevenue` - MRR 계산

### Integration Scenarios (3개)
1. `testDashboardWithMetrics` - 메트릭 + 대시보드
2. `testTracingWithAlerting` - 추적 + 알림
3. `testEnd2EndAnalytics` - 전체 통합

---

## 📊 성능 목표

| 작업 | 목표 | 달성도 |
|------|------|--------|
| Metric Collection | <100ms | Target |
| Trace Analysis | <500ms | Target |
| Alert Detection | <1sec | Target |
| Cost Calculation | <2sec | Target |
| Dashboard Rendering | <1sec | Target |

---

## 🔄 Track A & B 병렬 진행 타임라인

```
Week 1:
  Mon-Wed: Track A (Web3, Streaming, Multi-tenant, Cache)
  Mon-Wed: Track B (Dashboard, Tracing, Alerts, Cost Analysis)
  Thu: 각 Track 통합 테스트
  Fri: 성능 최적화

Weekend: 크로스 테스트 (Track A + B 통합)

Week 2:
  Mon-Tue: 최종 최적화 & 문서화
  Wed: 최종 통합 테스트
  Thu: 완료 보고서 & GOGS 커밋
```

---

## ✅ 완료 기준

- [ ] 4개 모듈 구현 (2,000줄)
- [ ] 20+ 테스트 작성 및 통과
- [ ] 단위 테스트: 100% 통과
- [ ] 통합 테스트: 5개 시나리오 통과
- [ ] 성능 목표 달성
- [ ] 문서 완성 (500줄)
- [ ] GOGS 커밋 & Push

---

## 🎯 Track A + B 통합 기대효과

```
Advanced Features (A):
  ✅ 블록체인 기반 감사 추적
  ✅ 실시간 스트리밍
  ✅ 다중 테넌트 지원
  ✅ 분산 캐싱

Analytics & Monitoring (B):
  ✅ 실시간 성능 대시보드
  ✅ 전체 경로 추적
  ✅ 이상 감지 및 알림
  ✅ 사용량 기반 청구

결과:
  🏆 엔터프라이즈급 기능 완성
  🏆 운영 가시성 극대화
  🏆 비용 최적화 지원
  🏆 SaaS 플랫폼 준비 완료
```

---

**Track B 시작 준비 완료!** 🚀

Track A와 동시에 진행됩니다.

---

**작성일**: 2026-03-02
**상태**: 📋 계획 완료, 🚀 구현 시작 준비
