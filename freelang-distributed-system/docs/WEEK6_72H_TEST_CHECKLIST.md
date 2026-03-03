# ⏰ Week 6: 72시간 Long-Running Test 체크리스트

**기간**: 2026-03-04 09:00 (화) ~ 2026-03-07 09:00 (금)
**정확 기간**: 259,200초 = 72시간
**목표**: 시스템 안정성 실증 (메모리 누수 0, 에러율 < 0.1%)

---

## 🎯 **72시간의 의미**

### 왜 정확히 72시간인가?

```
1시간    = 초기 안정화 확인
1일 (24h)= 일일 사이클 패턴 감지
3일 (72h)= 주간 패턴 + 장기 안정성 증명

메모리 누수를 감지하려면:
├─ 처음 1시간: 정상
├─ 첫 날: 문제 없음 (짧은 기간은 누수가 안 보임)
├─ 3일: 선형 누수도 명확 (100MB 누수 = 1.4MB/hour × 72h)
└─ 누수가 있으면 반드시 3일 내에 드러남

Performance degradation:
├─ 처음 10시간: 캐시 워밍
├─ 24-48시간: 패턴 변화 감지 가능
├─ 72시간: 진정한 안정 여부 판정
```

---

## 📊 **Pre-Test 체크리스트 (Mon 08:00)**

### 시스템 상태 확인

```
□ 5-node 클러스터 상태
  └─ All nodes alive? (5/5 up)
  └─ Raft quorum 가능? (≥3 up)
  └─ Network connectivity? (ping all)

□ 메모리 상태
  └─ Total RAM: 2GB/node (확인)
  └─ Available: > 1.5GB (test 시작 전)
  └─ Swap disabled? (메모리 정확성 위해)

□ Disk 상태
  └─ Available space: > 50GB
  └─ I/O speed: fio 테스트 (baseline)
  └─ fsync latency < 5ms

□ 네트워크 상태
  └─ RTT: < 1ms (loopback)
  └─ Cross-node: < 5ms
  └─ Packet loss: 0%

□ 소프트웨어 상태
  └─ FreeLang binary 최신 버전
  └─ Configuration 로드 확인
  └─ Logging 준비 (모든 노드)
  └─ Metrics collector 시작

□ 모니터링 준비
  └─ Prometheus 또는 custom collector
  └─ System metrics (top, iostat)
  └─ Application logs 중앙화
  └─ Time sync (NTP) - 모든 노드
  └─ Clock skew < 100ms
```

---

## 🔄 **Test 구조: 3가지 Tier**

### **Tier 1: Baseline (처음 1시간)**

```
목표: 정상 상태 확립

Mon 09:00 - 10:00:
├─ Warm-up phase (0-10분)
│  ├─ 1K vectors 삽입
│  ├─ 100 search queries
│  ├─ GC 실행 (초기 메모리 할당)
│  └─ Latency baseline 측정

├─ Measurement phase (10분 - 1시간)
│  ├─ 10K ops/sec 시작
│  ├─ 각 operation:
│  │  ├─ Latency 측정
│  │  ├─ Memory usage 기록
│  │  └─ CPU usage 기록
│  │
│  └─ Percentiles 계산:
│     ├─ P50 (중앙값)
│     ├─ P99 (99th percentile)
│     ├─ P99.9 (99.9th percentile)
│     └─ Max

┌────────────────────────┐
│ Baseline Metrics (1h)  │
├────────────────────────┤
│ Throughput:    10K/sec │
│ Latency P99:   ~8ms    │
│ Memory:        ~1500MB │
│ Error rate:    0%      │
│ CPU:           ~60%    │
└────────────────────────┘
```

---

### **Tier 2: Sustained Load (26시간)**

```
목표: 패턴 안정성 확인

Mon 10:00 - Tue 12:00 (26시간):
├─ Continuous operation
│  ├─ 매초 10,000 ops 유지
│  ├─ INSERT 50%, SEARCH 50% 비율
│  ├─ 1,000 concurrent connections
│  └─ NO interruption

├─ Hourly measurements:
│  ├─ Throughput (ops/sec, current vs baseline)
│  ├─ Latency (P50, P99, P99.9)
│  ├─ Memory usage (RSS)
│  ├─ CPU usage (%)
│  ├─ Error count
│  └─ GC frequency (if applicable)

├─ Anomaly detection:
│  └─ "대각선" 패턴
│     ├─ Latency 선형 증가? → 메모리 부족
│     ├─ Throughput 감소? → 병목 발생
│     ├─ Error 급증? → 클러스터 불안정
│     └─ CPU 선형 증가? → 누수 패턴

監視 데이터 포맷:
┌─────────┬───────┬──────┬──────┬──────┬──────┐
│ Hour    │ Tps   │ P99  │ Mem  │ CPU  │ Err  │
├─────────┼───────┼──────┼──────┼──────┼──────┤
│ 1:00    │ 10000 │ 8.2  │ 1500 │ 62%  │ 0    │
│ 2:00    │ 10000 │ 8.3  │ 1510 │ 61%  │ 0    │
│ 3:00    │ 10000 │ 8.1  │ 1512 │ 60%  │ 0    │
│ ...     │ ...   │ ...  │ ...  │ ...  │ ...  │
│ 24:00   │ 9980  │ 9.5  │ 1850 │ 58%  │ 0    │
│ 26:00   │ 9950  │ 10.2 │ 1920 │ 55%  │ 0    │
└─────────┴───────┴──────┴──────┴──────┴──────┘

패턴 분석:
- Throughput: 10000 → 9950 (0.5% 감소)
- Latency: 8.2 → 10.2ms (24% 증가)
- Memory: 1500 → 1920MB (420MB 증가)
- Memory leak rate: 420MB / 26h = 16MB/h ❌ (목표 < 1.5MB/h)
```

---

### **Tier 3: Long-term Stability (45시간)**

```
목표: 최종 안정성 판정

Tue 12:00 - Fri 09:00 (45시간):
├─ 추가 부하 패턴 도입
│  ├─ 12시간마다 높은 부하 (20K ops/sec)
│  ├─ 장시간 검색 쿼리 (top-1000 검색)
│  └─ 메모리 압박 구간
│
├─ 12시간 주기 체크 (3회):
│  ├─ Tue 12:00 (26h)
│  ├─ Wed 12:00 (50h)
│  └─ Thu 12:00 (74h) ← 이미 72h 초과 주의!
│
├─ 최종 일관성 검증:
│  └─ Fri 08:00
│     ├─ 전체 데이터 스캔
│     ├─ 모든 노드 간 data consistency check
│     ├─ Replication lag 확인
│     └─ 최종 메트릭 수집

패턴:
┌─────────┬───────┬──────┬──────┬──────┬──────┐
│ Hour    │ Tps   │ P99  │ Mem  │ CPU  │ Err  │
├─────────┼───────┼──────┼──────┼──────┼──────┤
│ 26:00   │ 10000 │ 10.2 │ 1920 │ 55%  │ 0    │
│ 27:00   │ 10000 │ 10.1 │ 1928 │ 54%  │ 0    │
│ 30:00   │ 20000 │ 18.5 │ 1950 │ 75%  │ 0    │ ← 높은 부하
│ 36:00   │ 10000 │ 10.3 │ 2050 │ 52%  │ 0    │
│ ...     │ ...   │ ...  │ ...  │ ...  │ ...  │
│ 50:00   │ 10000 │ 10.5 │ 2100 │ 48%  │ 0    │
│ 72:00   │ 10000 │ 10.7 │ 2150 │ 45%  │ 0    │
└─────────┴───────┴──────┴──────┴──────┴──────┘

패턴 분석:
- Throughput: 일관성 유지 (±1%)
- Latency: 점진적 증가 (8.2 → 10.7ms = +30%)
- Memory: 1500 → 2150MB (650MB 증가 / 72h = 9MB/h)
- 평가: ❌ 메모리 누수 (목표 < 1.5MB/h)
```

---

## 📈 **메트릭 수집 기준**

### **매 1분 (Granular)**

```
timestamp: 2026-03-04T09:01:00Z

current_metrics = {
  "throughput": {
    "ops_sec": 10000,
    "inserts": 5000,
    "searches": 5000,
    "errors": 0
  },
  "latency": {
    "p50_ms": 3.2,
    "p99_ms": 8.5,
    "p99_9_ms": 12.3,
    "max_ms": 45.8
  },
  "resources": {
    "memory_mb": 1512,
    "memory_percent": 75.6,
    "cpu_percent": 62.3,
    "network_in_mbps": 24.5,
    "network_out_mbps": 18.3
  },
  "errors": {
    "count": 0,
    "rate": 0.0
  }
}

저장: metrics_<minute>.json
```

### **매 1시간 (Hourly Summary)**

```
hourly_summary = {
  "hour": 1,
  "time_start": "2026-03-04T09:00:00Z",
  "time_end": "2026-03-04T10:00:00Z",

  "throughput": {
    "avg": 10000,
    "min": 9950,
    "max": 10050
  },
  "latency": {
    "p50_avg": 3.1,
    "p99_avg": 8.4,
    "p99_avg_delta": 0.2  // vs baseline
  },
  "memory": {
    "start": 1500,
    "end": 1510,
    "delta": 10,  // MB
    "rate": 10  // MB/hour
  },
  "cpu": {
    "avg": 62,
    "min": 58,
    "max": 68
  },
  "errors": {
    "total": 0,
    "rate": 0.0
  },

  "anomalies": [
    // 탐지된 이상 현상
  ]
}

저장: hourly_<hour>.json
```

### **매 12시간 (Major Checkpoint)**

```
checkpoint_12h = {
  "checkpoint_id": 1,
  "elapsed_hours": 12,
  "cumulative_stats": {
    "total_ops": 432000000,  // 10K ops/s × 12h
    "total_inserts": 216000000,
    "total_searches": 216000000,
    "total_errors": 0,
    "error_rate": 0.0
  },
  "memory_trend": {
    "initial": 1500,
    "current": 1680,
    "delta": 180,
    "leak_rate_mb_per_hour": 15.0,  // 목표: < 1.5
    "projection_72h": 2580  // 예상 최종 메모리
  },
  "latency_trend": {
    "p99_baseline": 8.2,
    "p99_current": 8.8,
    "p99_increase_percent": 7.3
  },
  "health_score": {
    "memory_health": "WARN",  // 누수 속도 높음
    "latency_health": "OK",
    "error_health": "OK",
    "overall": "WARN"
  },
  "decisions": {
    "continue": true,
    "reason": "메모리 누수 있지만 72h 내 OOM 안 될 것 같음",
    "next_action": "계속 모니터링, 72h 후 분석"
  }
}

저장: checkpoint_12h_<n>.json
```

---

## 🔍 **실시간 모니터링 (Streaming Dashboard)**

### **자동 경고 (Auto Alert)**

```
메모리 누수 감시:
├─ IF memory_rate > 5 MB/h:
│  └─ 🟠 WARNING: "High memory leak rate"
│
├─ IF memory_rate > 10 MB/h:
│  └─ 🔴 CRITICAL: "Severe memory leak"
│     └─ Action: Start GC analysis
│
└─ IF projected_memory > 2200MB at 72h:
   └─ 🔴 CRITICAL: "OOM risk within 72h"
      └─ Action: Consider early termination

Latency 감시:
├─ IF p99_latency > baseline × 1.5:
│  └─ 🟠 WARNING: "Latency degradation"
│
└─ IF p99_latency > baseline × 2.0:
   └─ 🔴 CRITICAL: "Severe latency"
      └─ Action: Investigate bottleneck

Error 감시:
├─ IF error_rate > 0.01% (1 per 10K):
│  └─ 🔴 CRITICAL: "Errors detected"
│     └─ Action: Immediate investigation
│
└─ IF any hard errors (OOM, crash, network):
   └─ 🔴 STOP: Immediate termination
```

---

## 📋 **Hourly Checklist (자동화)**

```python
# 매 시간 실행
def hourly_checklist():
    # 1. 메트릭 수집
    metrics = collect_metrics()

    # 2. 경고 확인
    for alert in check_alerts(metrics):
        if alert.severity == CRITICAL:
            notify_critical(alert)

    # 3. 트렌드 분석
    memory_trend = analyze_memory_trend(metrics)
    if memory_trend.leak_rate > 5:  # MB/h
        log_warning(f"Memory leak detected: {memory_trend.leak_rate} MB/h")

    # 4. 데이터 무결성 체크 (12시간마다)
    if hour % 12 == 0:
        verify_consistency()

    # 5. 로그 기록
    save_hourly_summary()

    # 6. 계속할지 판정
    if should_terminate():
        stop_test()
        collect_final_data()
    else:
        schedule_next_hour()
```

---

## ⚠️ **중단 조건 (Stop Criteria)**

### **Immediate Termination**

```
ANY 발생 시 즉시 중단:

1. Hard Error
   □ Process crash
   □ OOM exception
   □ Network partition
   □ Data inconsistency

2. Performance Cliff
   □ P99 latency > 50ms (갑작스러운 증가)
   □ Throughput < 5K ops/s (50% 이상 감소)
   □ Error rate > 0.1%

3. Memory Crisis
   □ Memory > 2000MB (안전 한계)
   □ Projected OOM within 12h

4. System Hung
   □ No response > 5분
   □ CPU 0% (deadlock?)
```

### **Graceful Termination**

```
72시간 완료 후:
└─ 마지막 메트릭 수집
└─ 최종 데이터 일관성 검증
└─ 클러스터 상태 스냅샷
└─ 로그 수집 및 분석
```

---

## 📊 **최종 분석 (Fri 09:00)**

### **통과 기준 (All Must Pass)**

```
✅ Error Rate:
   └─ Actual: 0.0%
   └─ Target: < 0.1%
   └─ Result: PASS ✓

✅ Memory Leak:
   └─ Actual: 9 MB/h (650 / 72)
   └─ Target: < 1.5 MB/h
   └─ Result: FAIL ✗ (누수 있음, 심각)

✅ Latency Degradation:
   └─ Initial P99: 8.2ms
   └─ Final P99: 10.7ms
   └─ Increase: 30%
   └─ Target: < 10%
   └─ Result: FAIL ✗ (성능 저하)

✅ Throughput Maintained:
   └─ Actual: 99.5% (10000 → 9950 avg)
   └─ Target: > 95%
   └─ Result: PASS ✓

✅ Consistency:
   └─ Data corruption: 0
   └─ Replication lag: 0
   └─ Result: PASS ✓
```

### **최종 판정**

```
Tier 1: ✅ Baseline OK
  └─ 정상 상태 확립

Tier 2: ⚠️ Sustained Load WARN
  └─ 메모리 누수 감지됨

Tier 3: ❌ Long-term FAIL
  └─ 72시간 내 메모리 누수 명확
  └─ Latency 30% 증가

OVERALL VERDICT:
├─ Functional: ✅ YES (시스템 작동)
├─ Stable: ⚠️ PARTIAL (누수 있음)
├─ Production Ready: ❌ NO (메모리 누수)
│
└─ NEXT STEP:
   ├─ 메모리 프로파일링 (누수 원인)
   ├─ GC 튜닝
   ├─ Latency 최적화
   └─ 재테스트 (개선 후)
```

---

## 🔧 **Fallback Plan (실패 시)**

### **메모리 누수 발견**

```
원인 분석:
├─ 1. Heap dump 수집 (Fri 09:00)
├─ 2. Object allocation 추적
├─ 3. Reference cycle 확인
└─ 4. GC behavior 분석

수정:
├─ 1. 누수 원인 코드 식별
├─ 2. 참조 제거 또는 weak reference 사용
├─ 3. GC 호출 명시
└─ 4. Unit test 추가

재테스트:
└─ 72시간 다시 실행 (개선 후)
```

### **Latency 저하**

```
원인 분석:
├─ 1. CPU profile (flame graph)
├─ 2. I/O latency 측정
├─ 3. Lock contention 확인
└─ 4. Cache hit rate 분석

최적화:
├─ 1. 핫스팟 코드 최적화
├─ 2. I/O batching
├─ 3. Lock-free 구조 검토
└─ 4. 캐시 워밍

재테스트:
└─ 성능 복구 확인 후 재실행
```

---

## 📌 **72시간 일정 정리**

```
Week 6 Timeline:

Mon 09:00 - 10:00
└─ Tier 1: Baseline (1시간)
   └─ Warm-up: latency baseline 측정

Mon 10:00 - Tue 12:00
└─ Tier 2: Sustained Load (26시간)
   └─ 10K ops/sec 유지
   └─ 패턴 안정성 확인

Tue 12:00 - Fri 09:00
└─ Tier 3: Long-term (45시간)
   └─ 부하 변화 도입
   └─ 최종 안정성 판정

Fri 09:00 - 11:00
└─ Final Analysis
   └─ 모든 메트릭 분석
   └─ Production Ready 판정
```

---

## ✅ **Pre-Test 확인 사항 (Mon 08:00)**

```
Hardware:
□ 5 nodes, 2GB RAM each
□ SSD storage > 50GB each
□ Network 1Gbps
□ NTP sync (skew < 100ms)

Software:
□ FreeLang binary latest
□ Configuration loaded
□ Logging prepared
□ Metrics collector ready

Monitoring:
□ Prometheus up
□ Time sync verified
□ Alerting configured
□ Storage prepared

Baseline:
□ Memory baseline recorded
□ Latency baseline measured
□ Network baseline tested
```

---

**생성일**: 2026-03-03
**적용**: 2026-03-04 09:00 (정확히 72시간)
**목표**: Production Ready 최종 판정
**상태**: 🔨 Week 6 72h Long-Running 준비 완료

