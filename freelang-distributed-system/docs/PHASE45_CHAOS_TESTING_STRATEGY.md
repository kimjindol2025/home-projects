# Phase 4.5: Chaos Testing & Validation Strategy

**목표**: Phase 4 API 레이어를 실제 장애 환경에서 검증
**철학**: "부숴보기 전까지는 신뢰하지 않는다"
**기간**: Week 5-6 (10일)

---

## 📋 핵심 가설 검증

### 가설 1: Raft Consensus는 Leader 장애를 1초 이내에 복구한다

**현재 상태**: 이론 (설계 문서)
**검증 방법**: 
```
- Leader를 5초마다 강제 종료
- 100회 반복
- 각 회차마다 MTTR 기록
```

**통과 기준**:
- 평균 MTTR < 1초
- 표준편차 < 0.3초
- 0.1% 이상의 재선출 실패 없음

**실패 시 의미**:
- Raft 구현에 로깅 지연 있음
- 또는 파이프라인이 느림

---

### 가설 2: 네트워크 불안정 상황에서 데이터 일관성을 유지한다

**현재 상태**: Quorum 기반 설계 (이론)
**검증 방법**:
```
- 30% 패킷 손실 주입 (30초 지속)
- 200ms 왕복 지연 추가
- 이 동안 INSERT/SEARCH 계속
- 파티션 해제 후 데이터 스캔
```

**통과 기준**:
- Dirty read 0건
- Write 손실 0건
- Consistency violation 0건

**실패 시 의미**:
- Quorum write 로직에 구멍
- 또는 recovery 프로토콜 버그

---

### 가설 3: Shard Rebalancing 중에도 latency SLA를 유지한다

**현재 상태**: 설계 (구현 미완료)
**검증 방법**:
```
- 동시 10K 연결 유지
- 초당 10K ops 생성
- 파티션 마이그레이션 시작
- Latency 분포 기록
```

**통과 기준**:
- P99 latency < 100ms
- Timeout rate < 1%
- Throughput drop < 20%

**실패 시 의미**:
- Rebalancing이 lock을 잡음
- 또는 네트워크 대역폭 포화

---

### 가설 4: 메모리 압박 상황에서 Graceful degradation이 작동한다

**현재 상태**: 미설계
**검증 방법**:
```
- 메모리를 95% 포화시킴
- 동시에 INSERT 계속
- OOM 발생 시 동작 관찰
- 디스크 I/O stall도 추가 (fsync 10배 느려짐)
```

**통과 기준**:
- OOM crash 0건
- Out-of-memory exception 처리
- 자동 복구 < 30초

**실패 시 의미**:
- 메모리 관리 정책 필요
- 또는 circuit breaker 패턴 추가

---

### 가설 5: 실제 embedding 데이터에서도 성능이 유지된다

**현재 상태**: 랜덤 벡터로 테스트됨
**검증 방법**:
```
- OpenAI embedding 1M 데이터 로드
- 같은 쿼리 패턴 실행
- 성능/메모리 측정
- 랜덤 벡터 결과와 비교
```

**통과 기준**:
- 성능 ±10% 범위
- 메모리 사용량 예측 정확도 > 80%
- 압축률 가정 재검증

**실패 시 의미**:
- 벡터 분포에 따라 성능 변동
- 양자화 가능성 탐색 필요

---

## 🔧 Chaos Test Framework 설계

### 기본 구조

```
ChaosTestFramework
├─ ChaosEngine (장애 주입)
│  ├─ injectPacketLoss(percent, duration)
│  ├─ injectLatency(ms, variance)
│  ├─ killRandomNode()
│  └─ exhaustMemory(percent)
│
├─ MetricsCollector (측정)
│  ├─ recordLatency(operation, ms)
│  ├─ recordThroughput(ops/sec)
│  ├─ recordConsistency(pass/fail)
│  └─ recordRecoveryTime(ms)
│
├─ ValidationEngine (검증)
│  ├─ validateDataConsistency()
│  ├─ validateLatencySLA()
│  ├─ validateThroughput()
│  └─ checkMemoryLeaks()
│
└─ ReportGenerator (보고)
   ├─ generateTestReport()
   ├─ generateLatencyHistogram()
   ├─ generateFailureAnalysis()
   └─ compareToPreviousRun()
```

### 실행 패턴

```
for each chaos_scenario:
  for iteration = 1 to max_iterations:
    1. Setup (clean state)
    2. Warm up (5s)
    3. Baseline measurement (5s)
    4. Inject chaos (apply fault)
    5. Measure under fault (30s)
    6. Fault removal
    7. Measure recovery (10s)
    8. Validate results
    9. Record metrics
    10. Teardown
```

---

## 📊 Test Case 상세 정의

### Test 1: Raft Leadership Chaos

**File**: `tests/chaos_raft_leadership.fl` (200줄)

```
fn chaosRaftLeadershipKill()
  // Setup: 5-node Raft cluster
  let cluster = newRaftIndexCluster(nodeCount: 5)
  startCluster(cluster)
  
  // Baseline: normal ops
  let baseline = measureThroughput(cluster, duration: 5s)
  
  // Chaos loop
  var results = []
  for iteration = 0 to 100
    // Get current leader
    let leader = cluster.getLeader()
    let startTime = getCurrentTime()
    
    // Kill leader
    cluster.killNode(leader)
    
    // Measure time to new leader elected
    let newLeader = cluster.waitForLeaderElection(timeout: 5s)
    let mttr = getCurrentTime() - startTime
    
    // Record
    results.append({
      "iteration": iteration,
      "mttr_ms": mttr,
      "leader_before": leader,
      "leader_after": newLeader,
      "lost_requests": countLostOps(during: mttr)
    })
    
    // Wait 5s before next kill
    sleep(5000)
    
    // Revive killed node for next iteration
    cluster.reviveNode(leader)
    sleep(1000)  // Wait for catchup
  
  // Analysis
  let avgMttr = average(results.map(r => r.mttr_ms))
  let maxMttr = max(results.map(r => r.mttr_ms))
  let stdDev = stdDeviation(results.map(r => r.mttr_ms))
  
  println("Leadership Chaos Results:")
  println("  Avg MTTR: " + avgMttr + "ms (target < 1000ms)")
  println("  Max MTTR: " + maxMttr + "ms")
  println("  Std Dev: " + stdDev + "ms")
  println("  Pass: " + (avgMttr < 1000))
```

**측정 항목**:
- MTTR (재선출 시간)
- Lost operations
- Consistency violations
- 각 iteration별 상세 기록

**통과 기준**:
```
avgMttr < 1000ms &&
maxMttr < 2000ms &&
stdDev < 300ms &&
dataInconsistency == 0
```

### Test 2: Network Partition

**File**: `tests/chaos_network_partition.fl` (200줄)

```
fn chaosNetworkPartition()
  let cluster = newRaftIndexCluster(nodeCount: 5)
  startCluster(cluster)
  
  // Pre-write baseline data
  for i = 0 to 1000
    cluster.insert("vec_" + i, randomVector())
  
  // Baseline metrics
  let baselineLatency = measureLatency(cluster, ops: 1000)  // [p50, p99, p99.9]
  let baselineThroughput = measureThroughput(cluster, duration: 5s)
  
  // Introduce 30% packet loss + 200ms latency
  cluster.networkManager.injectPacketLoss(30)
  cluster.networkManager.injectLatency(200)
  
  // Run ops for 30s under fault
  var faultyMetrics = []
  for duration = 0 to 30000 step 1000
    let latency = measureLatency(cluster, ops: 100)
    let throughput = measureThroughput(cluster, duration: 1s)
    faultyMetrics.append({
      "time_ms": duration,
      "latency_p99": latency.p99,
      "throughput": throughput
    })
  
  // Remove fault
  cluster.networkManager.removePacketLoss()
  cluster.networkManager.removeLatency()
  
  // Measure recovery (consistency check)
  sleep(5000)  // Wait for catch-up
  
  // Validate consistency
  var inconsistencies = 0
  for replica = 0 to 4
    let data1 = cluster.getNode(0).scanAll()
    let data2 = cluster.getNode(replica).scanAll()
    if data1 != data2
      inconsistencies = inconsistencies + 1
  
  println("Network Partition Results:")
  println("  Packet Loss: 30%, Latency: 200ms")
  println("  Latency during fault:")
  println("    P99: " + avgOfP99s + "ms")
  println("    P99.9: " + avgOfP99_9s + "ms")
  println("  Data Inconsistencies: " + inconsistencies + " (target: 0)")
  println("  Pass: " + (inconsistencies == 0))
```

**측정 항목**:
- Latency 분포 (P50, P99, P99.9)
- Throughput 변화
- 일관성 위반 (scan 비교)

**통과 기준**:
```
consistencyViolations == 0 &&
p99Latency < 100ms (under 30% loss)
```

### Test 3: Shard Rebalancing

**File**: `tests/chaos_shard_rebalance.fl` (150줄)

```
fn chaosShardbRebalancing()
  let cluster = newRaftIndexCluster(nodeCount: 5, partitions: 5)
  startCluster(cluster)
  
  // Warm up: 10K concurrent connections
  let connectionPool = newConnectionPool(10000)
  let workloadGenerator = startWorkload(
    connections: connectionPool,
    opsPerSecond: 10000,
    operationType: ["insert", "search"]
  )
  
  // Baseline (30s)
  sleep(30000)
  let baselineLatencies = collectLatencies(duration: 30s)
  let baselineP99 = baselineLatencies.p99  // Assume ~12ms
  
  // Trigger rebalancing
  let rebalanceStart = getCurrentTime()
  cluster.triggerRebalancing()  // Move partition 3 from node1 → node2
  
  // Collect metrics during rebalancing (60s)
  var duringRebalanceLatencies = []
  var timeouts = 0
  
  for second = 0 to 60
    let latencies = collectLatencies(duration: 1s)
    duringRebalanceLatencies.append(latencies)
    
    if latencies.timeouts > 0.01 * operationCount
      timeouts = timeouts + latencies.timeouts
    
    sleep(1000)
  
  let rebalanceDuration = getCurrentTime() - rebalanceStart
  
  // Post-rebalance validation (30s)
  sleep(30000)
  let postLatencies = collectLatencies(duration: 30s)
  
  // Analysis
  let duringP99 = percentile(duringRebalanceLatencies, 99)
  let latencyIncrease = (duringP99 - baselineP99) / baselineP99 * 100
  
  println("Shard Rebalancing Chaos Results:")
  println("  Baseline P99: " + baselineP99 + "ms")
  println("  During P99: " + duringP99 + "ms")
  println("  Increase: " + latencyIncrease + "%")
  println("  Timeouts: " + timeouts)
  println("  Rebalance duration: " + rebalanceDuration + "ms")
  println("  Pass: " + (latencyIncrease < 10 && timeouts == 0))
```

**측정 항목**:
- Latency P99 변화
- Timeout 비율
- Rebalancing 시간
- Throughput 유지율

**통과 기준**:
```
latencyIncreasePercent < 10 &&
timeoutRate < 1% &&
rebalanceDuration < 120000  // 2 minutes
```

### Test 4: Resource Stress

**File**: `tests/chaos_resource_stress.fl` (150줄)

```
fn chaosResourceStress()
  let cluster = newRaftIndexCluster(nodeCount: 5)
  startCluster(cluster)
  
  // Stress memory to 95%
  let memoryManager = cluster.getMemoryManager()
  let totalMem = memoryManager.getTotalMemory()
  let targetMem = totalMem * 0.95
  
  let stressAlloc = memoryManager.allocate(targetMem)
  
  // Also simulate slow disk
  let ioManager = cluster.getIOManager()
  ioManager.slowDownFsync(10)  // 10x slower
  
  // Try to insert while under stress
  var insertResults = []
  for i = 0 to 1000
    let result = cluster.insert("vec_" + i, randomVector())
    insertResults.append({
      "index": i,
      "success": result.ok,
      "error": result.err if err else "none"
    })
  
  // Check for OOM
  let oomCount = insertResults.filter(r => r.error == "OutOfMemory").length()
  
  // Release stress
  memoryManager.deallocate(stressAlloc)
  ioManager.normalizeIoSpeed()
  
  // Measure recovery time
  let recoveryStart = getCurrentTime()
  var recovered = false
  while !recovered && getCurrentTime() - recoveryStart < 30000
    let memUsage = memoryManager.getUsagePercent()
    if memUsage < 70
      recovered = true
    sleep(1000)
  
  let recoveryTime = getCurrentTime() - recoveryStart
  
  println("Resource Stress Chaos Results:")
  println("  Memory stress: 95%")
  println("  Disk slowdown: 10x")
  println("  OOM exceptions: " + oomCount)
  println("  Recovery time: " + recoveryTime + "ms")
  println("  Pass: " + (recoveryTime < 30000 && oomCount == 0))
```

**측정 항목**:
- OOM 발생 여부
- Exception 처리
- 복구 시간
- 메모리 안정화

**통과 기준**:
```
oomCrashes == 0 &&
recoveryTime < 30000 &&
systemStable == true
```

### Test 5: Real Data Validation

**File**: `tests/validate_real_data.fl` (150줄)

```
fn validateRealDataPerformance()
  // Load OpenAI embeddings (1M vectors, 1536 dimensions)
  println("Loading OpenAI embedding dataset...")
  let dataset = loadOpenAIEmbeddings(count: 1000000)
  
  // Insert into cluster
  let cluster = newRaftIndexCluster(nodeCount: 5)
  startCluster(cluster)
  
  println("Inserting 1M vectors...")
  let insertStart = getCurrentTime()
  var insertedCount = 0
  
  for vector in dataset
    let result = cluster.insert(vector.id, vector.embedding)
    if result.ok
      insertedCount = insertedCount + 1
  
  let insertDuration = getCurrentTime() - insertStart
  
  // Search performance
  println("Benchmarking search performance...")
  var searchLatencies = []
  
  for i = 0 to 10000  // 10K random searches
    let queryVector = randomFromDataset(dataset)
    let searchStart = getCurrentTime()
    let results = cluster.search(queryVector, topK: 10)
    let latency = getCurrentTime() - searchStart
    
    searchLatencies.append(latency)
  
  // Compare with random vector results
  let randomVectorLatencies = compareWithRandomVectors(10000)
  
  // Memory usage
  let memoryUsage = cluster.getMemoryManager().getUsageBytes()
  let estimatedMemory = calculateExpectedMemory(1000000)
  let memoryAccuracy = memoryUsage / estimatedMemory
  
  // Compression ratio (actual)
  let actualCompressionRatio = cluster.getCompressionStats().actualRatio
  
  // Results
  println("Real Data Validation Results:")
  println("  Inserted: " + insertedCount + " / 1000000")
  println("  Insert time: " + insertDuration + "ms")
  println("  Search latency (real data):")
  println("    P50: " + percentile(searchLatencies, 50) + "ms")
  println("    P99: " + percentile(searchLatencies, 99) + "ms")
  println("  Search latency (random): P99=" + percentile(randomVectorLatencies, 99) + "ms")
  println("  Memory usage: " + (memoryUsage / 1024 / 1024) + "MB")
  println("  Estimated: " + (estimatedMemory / 1024 / 1024) + "MB")
  println("  Accuracy: " + (memoryAccuracy * 100) + "%")
  println("  Compression ratio: " + actualCompressionRatio + " (expected: 0.35)")
  println("  Pass: " + (
    insertedCount == 1000000 &&
    abs(memoryAccuracy - 1.0) < 0.2 &&
    abs(actualCompressionRatio - 0.35) < 0.05
  ))
```

**측정 항목**:
- Insert 성능 (1M 벡터)
- Search latency (실제 데이터)
- 메모리 사용량 vs 예측값
- 압축률 (실제)

**통과 기준**:
```
insertCount == 1000000 &&
memoryAccuracy > 80% &&
compressionRatio ≈ 0.35 (±0.05)
```

### Test 6: Long-Running Stability

**File**: `tests/chaos_long_running.fl` (100줄)

```
fn chaosLongRunningStability()
  let cluster = newRaftIndexCluster(nodeCount: 5)
  startCluster(cluster)
  
  println("Starting 72-hour long-running test...")
  
  var metrics = {
    "startTime": getCurrentTime(),
    "iterations": 0,
    "errors": 0,
    "memorySnapshots": [],
    "latencyBuckets": []
  }
  
  // Run for 72 hours (or until failure)
  let testDuration = 72 * 60 * 60 * 1000  // 72 hours in ms
  let startTime = getCurrentTime()
  
  while getCurrentTime() - startTime < testDuration
    // Every second: INSERT + SEARCH
    for i = 0 to 100
      let vectorId = "iter_" + metrics.iterations + "_vec_" + i
      let result = cluster.insert(vectorId, randomVector())
      
      if !result.ok
        metrics.errors = metrics.errors + 1
      
      let searchResult = cluster.search(randomVector(), topK: 10)
      if !searchResult.ok
        metrics.errors = metrics.errors + 1
    
    metrics.iterations = metrics.iterations + 1
    
    // Every hour: snapshot memory
    if metrics.iterations % 3600 == 0
      let memUsage = cluster.getMemoryManager().getUsageBytes()
      metrics.memorySnapshots.append({
        "hour": metrics.iterations / 3600,
        "memoryMB": memUsage / 1024 / 1024
      })
      
      println("Hour " + (metrics.iterations / 3600) + 
              ": " + (memUsage / 1024 / 1024) + "MB, " +
              "Errors: " + metrics.errors)
    
    sleep(1000)
  
  // Analysis
  let errorRate = metrics.errors / (metrics.iterations * 200)  // 200 ops per iteration
  let memoryTrend = detectTrend(metrics.memorySnapshots)  // 0 = stable, >0 = leak
  
  println("Long-Running Test Results:")
  println("  Iterations: " + metrics.iterations)
  println("  Error rate: " + (errorRate * 100) + "%")
  println("  Memory trend: " + memoryTrend + " (0 = stable)")
  println("  Pass: " + (errorRate < 0.01 && memoryTrend < 0.05))
```

**측정 항목**:
- 누적 operation count
- Error rate
- 메모리 추이 (누수 감시)
- 성능 degradation

**통과 기준**:
```
errorRate < 1% &&
memoryLeakRate < 5% (per day) &&
noUnplannedCrashes == true
```

---

## 📈 결과 해석 가이드

### 통과 시: ✅ PASS

```
모든 6개 test pass
→ "이 시스템은 chaos를 견딘다"
→ Production deployment 가능
```

### 부분 통과: ⚠️ CONDITIONAL PASS

```
예: Test 1-3 pass, Test 4 fail
→ "메모리 관리에 문제가 있다"
→ Phase 4 코드 수정 필요
→ Test 4만 재실행
```

### 실패: ❌ FAIL

```
예: Test 2 fail (일관성 위반 발생)
→ "Quorum 구현에 구멍이 있다"
→ 근본 원인 분석 필요
→ 구조 재설계 가능성
→ 모든 test 재실행
```

---

## 🎯 최종 선언 기준

```
모든 test 통과 + Long-running 72h 안정 확인

↓

"Production Ready"

실제 선언:
"이 시스템은 장애 환경에서 테스트되었고,
의도적으로 부숴지지 않았다"
```

---

**다음 페이즈**: 실제로 chaos test 구현 및 실행
