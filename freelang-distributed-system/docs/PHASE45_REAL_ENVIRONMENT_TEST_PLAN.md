# Phase 4.5: 실제 환경 테스트 계획

**시작일**: 2026-03-04
**기간**: 10일 (Week 5-6)
**목표**: Production Ready 최종 선언
**성공 기준**: 모든 Chaos Test + 72h Long-Running PASS

---

## 📋 Week 5: 실제 환경 Chaos Test (5가지)

### Monday-Tuesday: Raft Leadership Kill (실제)

**환경 설정**
```
- 5-node Raft cluster (실제 프로세스/컨테이너)
- 각 노드: 2GB RAM, 4 CPU cores
- 네트워크: 1Gbps 이더넷
- 스토리지: SSD 100GB
```

**테스트 실행**
```bash
# 실제 환경에서 100회 반복
for i in {1..100}
  1. Leader 선택
  2. Kill -9 (강제 종료)
  3. MTTR 측정 (재선출까지의 시간)
  4. 데이터 일관성 검증 (scan & compare)
  5. 5초 대기
  6. Kill된 노드 복구
  7. Catch-up 대기 (1초)
```

**측정 항목**
```
- MTTR (milliseconds)
  목표: 평균 < 1000ms, 최대 < 2000ms

- 데이터 손실
  목표: 0 bytes

- Lost operations
  목표: < 200 (전체 10,000 중)

- Consistency violations
  목표: 0
```

**통과 기준**
```
avgMTTR < 1000ms ✓
maxMTTR < 2000ms ✓
dataLoss == 0 ✓
consistencyViolations == 0 ✓
```

### Wednesday: Network Partition (실제)

**환경 설정**
```
- Traffic Control (tc) 사용
- 30% 패킷 손실 (iptables netem)
- 200ms 왕복 지연 (latency)
- 지속 시간: 30초
```

**테스트 실행**
```bash
# Pre-write: 1000개 벡터 삽입
# Baseline: 정상 상태 Latency 측정 (10초)
# Inject fault: 30% loss + 200ms delay (30초)
# Measure: Latency, throughput, consistency
# Remove fault: 10초
# Post-measure: 일관성 검증
```

**측정 항목**
```
- Latency (정상): P50, P99, P99.9
- Latency (장애중): P50, P99, P99.9
- Timeout 발생 건수 (< 1%)
- 일관성 위반 (0건)
- Quorum write 성공률 (> 99%)
```

**통과 기준**
```
P99 latency (normal) < 15ms ✓
Consistency violations == 0 ✓
Timeout rate < 1% ✓
Quorum writes > 99% ✓
```

### Thursday: Shard Rebalancing (실제)

**환경 설정**
```
- 10,000 동시 연결 (Apache Bench or wrk)
- 지속적인 INSERT/SEARCH (10K ops/sec)
- 파티션 3 마이그레이션: node1 → node2
```

**테스트 실행**
```bash
# Baseline: 30초 (정상 상태 latency)
# Trigger rebalancing
# During: 60초 (rebalancing 진행 중 측정)
# Post: 30초 (복구 확인)
# Total: ~2분
```

**측정 항목**
```
- Baseline P99: ~12ms
- During P99: < 100ms (목표)
- Post P99: ~12ms (복구)
- Timeout rate: < 1%
- Throughput maintained: > 80%
- Rebalance duration: < 120s
```

**통과 기준**
```
During P99 < 100ms ✓
Timeout rate < 1% ✓
Throughput > 80% ✓
Rebalance < 120s ✓
Post P99 회복 < 12ms ✓
```

### Friday: Resource Stress (실제)

**환경 설정**
```
- cgroup으로 메모리 제한: 95%
- fio로 디스크 I/O stall (fsync 10배 느려짐)
- 동시 INSERT 100개
```

**테스트 실행**
```bash
# Apply memory pressure (95%)
# Apply disk slowdown (10x)
# Try to INSERT 1000 vectors
# Observe behavior (OOM? crash? slow?)
# Release stress
# Measure recovery time
# Verify data integrity
```

**측정 항목**
```
- Insert success rate: > 99%
- OOM exceptions: 0
- System crashes: 0
- Recovery time: < 30s
- Graceful degradation: yes/no
```

**통과 기준**
```
OOM crashes == 0 ✓
System crashes == 0 ✓
Recovery time < 30s ✓
Insert success > 99% ✓
```

---

## 📊 Week 6: Long-Running + 분석

### Monday: Real Data Test (1M OpenAI embeddings)

**데이터 준비**
```
- OpenAI embedding 1M 벡터
- 각 벡터: 1536 dimensions (float32)
- 총 크기: ~6GB (압축 후 ~2.2GB)
```

**테스트 실행**
```bash
# Insert: 1M vectors
# Measure:
#   - Insert throughput
#   - Memory usage
#   - Search latency (1K queries)
#   - Compression ratio (실제)
# Compare with simulation
```

**측정 항목**
```
- Insert success: 100%
- Search Latency P99: ~15ms (simulation: 16ms)
- Memory actual vs estimated
- Compression ratio: 0.36 (simulation: 0.36)
- Performance accuracy: ±10%
```

**통과 기준**
```
Insert success == 100% ✓
Search P99 < 20ms ✓
Memory accuracy > 80% ✓
Compression ratio ≈ 0.35 (±0.05) ✓
```

### Tuesday-Thursday: 72시간 Long-Running Test

**설정**
```
시작: 2026-03-04 09:00 (화요일)
종료: 2026-03-07 09:00 (금요일)
기간: 정확히 72시간
```

**워크로드**
```
- 동시 1000 연결 (상수 유지)
- 초당 10,000 ops (INSERT 50%, SEARCH 50%)
- 벡터 크기: 512 dimensions
- 지속: 72시간
```

**자동 모니터링**
```
매 1분:
  - Current throughput (ops/sec)
  - Current latency (P50, P99, P99.9)
  - Memory usage (MB)
  - CPU usage (%)
  - Network I/O (Mbps)
  - Error count

매 1시간:
  - Rolling 1h statistics
  - Error rate (%)
  - Memory trend (증가하나?)
  - Latency trend (악화하나?)
  - Alerts (임계값 초과)

매 12시간:
  - Summary report
  - Cumulative stats
  - Any anomalies
```

**관찰 항목**
```
1. 메모리 누수
   - 시작: 1500MB (추정)
   - 72h 후: < 1500MB + 100MB (안정)
   - 누수 속도: < 1.5MB/hour (허용)

2. 성능 Degradation
   - 시작 latency: 8ms
   - 72h 후: < 9.5ms (±10% 이내)
   - P99: 12ms → < 13ms

3. 에러율
   - 시작: < 0.01%
   - 지속: < 0.01% 유지

4. Throughput
   - 유지: 정확히 10K ops/sec ±5%
```

### Friday: 분석 & 최종 판정

**데이터 수집**
```
- 72시간 로그 파일 (GB 단위)
- 메트릭 데이터베이스 (Prometheus)
- 에러 로그 (stderr)
- 시스템 로그 (journalctl)
```

**분석 작업**
```
1. 성능 분석
   - 시뮬레이션 vs 실제 비교
   - 예상치 vs 실제값
   - 편차 분석

2. 안정성 분석
   - 에러 패턴 분석
   - 메모리 누수 추이
   - Latency degradation

3. 리소스 사용
   - 피크 메모리
   - 피크 CPU
   - 네트워크 대역폭

4. 개선 사항
   - 문제점 식별
   - 최적화 제안
   - 알려진 제한사항
```

**최종 리포트**
```
생성: PHASE45_REAL_ENVIRONMENT_FINAL_REPORT.md

내용:
- Executive Summary
- Test Results (5가지 Chaos Test)
- Long-Running Results (72h)
- Performance Comparison (Simulation vs Real)
- Stability Assessment
- Recommendations
- Production Ready Assessment
```

---

## 🎯 최종 판정 로직

```python
def is_production_ready():
    # 모든 조건을 AND로 체크
    
    leadership_pass = (
        avg_mttr < 1000 and
        max_mttr < 2000 and
        data_loss == 0 and
        consistency_violations == 0
    )
    
    network_pass = (
        consistency_violations == 0 and
        timeout_rate < 1 and
        quorum_writes_success > 99
    )
    
    rebalancing_pass = (
        p99_latency < 100 and
        timeout_rate < 1 and
        throughput_maintained > 80 and
        rebalance_duration < 120
    )
    
    resource_stress_pass = (
        oom_crashes == 0 and
        system_crashes == 0 and
        recovery_time < 30 and
        insert_success > 99
    )
    
    real_data_pass = (
        insert_success == 100 and
        search_latency_p99 < 20 and
        memory_accuracy > 80 and
        abs(compression_ratio - 0.35) < 0.05
    )
    
    long_running_pass = (
        error_rate < 0.1 and
        memory_leak_rate < 1.5 and  # MB/hour
        latency_degradation < 10 and  # percent
        throughput_maintained > 95 and
        cumulative_errors == 0
    )
    
    return (
        leadership_pass and
        network_pass and
        rebalancing_pass and
        resource_stress_pass and
        real_data_pass and
        long_running_pass
    )
```

---

## 📌 의사결정 기준

### ✅ 선언 가능: "Production Ready"
```
if is_production_ready():
    "이 시스템은 실제 환경에서 테스트되었고
     72시간 안정적으로 운영되었다.
     장애 환경에서도 견디는 것을 확인했다."
```

### 🟡 조건부 승인: "Almost Ready"
```
일부 test는 PASS, 일부는 FAIL
→ 실패한 부분만 개선
→ 재테스트
→ 반복
```

### ❌ 미승인: "Not Ready"
```
if critical_failures > 0:
    근본 원인 분석
    Phase 4 코드 수정
    모든 test 재실행
    
    (반복 가능성 높음)
```

---

## 📊 성공 메트릭 요약

| 메트릭 | 시뮬레이션 | 목표 | 실제 (예상) |
|--------|-----------|------|-----------|
| Leadership MTTR | 847ms | <1000ms | ~850ms |
| Network P99 | 52ms | OK | ~50ms |
| Rebalance P99 | 68ms | <100ms | ~70ms |
| Resource Recovery | 15s | <30s | ~15s |
| Real Data P99 | 16ms | OK | ~16ms |
| 72h Error Rate | 0% | <0.1% | ~0.01% |

---

**목표**: 모든 항목에서 시뮬레이션 결과와 유사하게 실제 환경에서도 통과

**최종 선언**: 2026-03-07 (금요일) 예정

