# 🚀 실제 환경 테스트 실행 계획 (Week 5-6)

**시작**: 2026-03-04 09:00 (정확히 이 시간)
**종료**: 2026-03-07 09:00 (72시간 후)
**환경**: 5-node 클러스터 (실제 환경)
**목표**: Production Ready 최종 판정

---

## 🎯 **테스트 구성 해석 (2-5-3)**

### **의미: "2개 노드 DOWN → 5-node 클러스터에서 → 3개만 ALIVE"**

```
기본 구성:
├─ Node 1 (leader)
├─ Node 2
├─ Node 3
├─ Node 4
└─ Node 5

2-5-3 구성 (추가 extreme test):
├─ Node 1: DOWN ❌
├─ Node 2: DOWN ❌
├─ Node 3: ALIVE ✅ (quorum OK: 3/5)
├─ Node 4: ALIVE ✅
└─ Node 5: ALIVE ✅

특징:
└─ 정확히 quorum 경계 (3/5 = barely alive)
└─ 4/5는 안전, 2/5는 불가능 → 3/5는 위험한 경계
```

---

## 📅 **Week 5-6 상세 일정 및 체크리스트**

### **Mon 2026-03-04**

#### **08:00 - 09:00: Pre-Test Setup (1시간)**

```
□ 시스템 헬스 체크
  ├─ 5개 노드 모두 alive? (ping)
  ├─ Raft quorum OK? (5/5)
  ├─ Network latency < 5ms?
  ├─ Memory > 1.5GB available?
  └─ Disk I/O speed baseline (fio)

□ 모니터링 준비
  ├─ Prometheus start (if used)
  ├─ Log aggregation start
  ├─ Metrics directory 생성
  ├─ Timestamp sync (NTP all nodes)
  └─ Network packet capture (optional)

□ Baseline 측정
  ├─ Normal latency (10초 × 1000 ops)
  ├─ Memory baseline (각 노드)
  ├─ CPU baseline (idle)
  ├─ Network RTT (node-to-node)
  └─ 결과 저장: baseline_metrics.json

Status: ✅ Ready to Start Test 1
```

#### **09:00 - 21:00: Test 1 - Raft Leadership Kill (12시간)**

```
타겟: MTTR < 1000ms, Data Loss = 0

Procedure (100회 반복):
├─ Iteration 1-20 (09:00 - 11:00):
│  ├─ 08:00: Current leader 확인
│  ├─ 08:01: kill -9 <pid>
│  ├─ 08:02: 재선출 대기 (MTTR 측정)
│  ├─ 08:05: 복구 노드 재시작
│  ├─ 08:06: Cluster rejoins
│  ├─ 08:07: Consistency check
│  ├─ 08:12: Wait 5분
│  └─ Metrics 기록: mttr_1.json
│
├─ Iteration 21-50 (11:00 - 16:00):
│  └─ 같은 절차 × 30회
│
├─ Iteration 51-100 (16:00 - 21:00):
│  └─ 같은 절차 × 50회

Hourly Checkpoints (매시간 11:00, 12:00, ..., 21:00):
├─ MTTR 통계 (min, max, avg, std)
├─ Data loss count (cumulative)
├─ 데이터 일관성 검증
└─ 이상 감지

Critical Stops:
├─ MTTR > 2000ms (3회 연속) → STOP
├─ Data loss 감지 → IMMEDIATE STOP
├─ Node crash (leader 아님) → STOP
└─ Network error → STOP

Result File: test1_leadership_kill.json
├─ mttr_statistics
├─ data_loss_count
├─ consistency_violations
└─ overall_pass (Y/N)
```

#### **21:00 - 22:00: Day 1 분석 (1시간)**

```
□ 로그 분석
  ├─ 비정상 종료 여부 확인
  ├─ MTTR 통계 계산
  ├─ 데이터 불일치 확인
  └─ Root cause 분석

□ 의사결정
  ├─ MTTR < 1000ms? → OK, 계속
  ├─ Data loss? → STOP, 분석
  └─ Anomalies? → 기록, 계속

Status: ✅ Complete or ❌ STOP
```

---

### **Tue 2026-03-05**

#### **09:00 - 12:00: Test 1 계속 (추가 50회) + 분석**

```
Iteration 101-150 (50회 추가):
├─ 같은 절차 (Mon과 동일)
├─ Cumulative statistics 업데이트
└─ Trend analysis (MTTR 변화)

12:00: Test 1 최종 판정
├─ 100회 × 10초 × baseline + quorum vote
├─ 평균 MTTR 계산
├─ 표준편차 계산
├─ 최종 PASS/FAIL
└─ 결과 저장: test1_final_report.json

PASS Criteria:
├─ MTTR average < 1000ms ✓
├─ MTTR max < 2000ms ✓
├─ Data loss = 0 ✓
└─ Consistency violations = 0 ✓
```

#### **12:00 - 13:00: Test 2 준비 (Network Partition)**

```
□ Network 도구 설정
  ├─ tc (traffic control) 확인
  ├─ iptables 규칙 준비
  ├─ 30% packet loss 스크립트
  ├─ 200ms delay 스크립트
  └─ 모든 명령 테스트

□ Baseline 재측정
  ├─ Normal latency (정상 상태)
  ├─ P50, P99, P99.9 기록
  ├─ Timeout count baseline
  └─ 1000 vectors 사전 삽입

Status: ✅ Ready for Test 2
```

#### **13:00 - 14:30: Test 2 - Network Partition (1.5시간)**

```
Timeline:
├─ 13:00-13:10: Baseline 측정 (정상 상태)
│  ├─ 100 INSERT ops
│  ├─ Latency 수집
│  └─ 통계 계산 (P50, P99)
│
├─ 13:10-13:40: Fault 적용 (30% loss + 200ms delay)
│  ├─ tc 명령 실행
│  ├─ Parallel INSERT/SEARCH (30초)
│  ├─ Latency 실시간 수집
│  ├─ Timeout count 기록
│  └─ Consistency check (매 10초)
│
├─ 13:40-13:50: Fault 해제
│  ├─ tc 규칙 제거
│  ├─ Recovery latency 측정 (10초)
│  └─ 최종 consistency 검증
│
└─ 13:50-14:30: 분석 및 결과

Result File: test2_network_partition.json
├─ baseline_latency
├─ during_partition_latency
├─ consistency_violations
├─ timeout_rate
└─ overall_pass (Y/N)

PASS Criteria:
├─ Consistency violations = 0 ✓
├─ P99 latency (normal) < 15ms ✓
├─ Timeout rate < 1% ✓
├─ Quorum writes > 99% ✓
```

#### **14:30 - 15:30: Test 3 준비 (Shard Rebalancing)**

```
□ 10K concurrent connections 준비
  ├─ Apache Bench or wrk 확인
  ├─ 스크립트 준비 (10K parallel requests)
  ├─ Connection pool 설정
  └─ Load generator 테스트

□ Partition 3 대상 선택
  ├─ 현재 위치 확인 (node1?)
  ├─ 마이그레이션 대상 선택 (node2?)
  └─ Rebalance 시간 예측

Status: ✅ Ready for Test 3
```

#### **15:30 - 16:30: Test 3 - Shard Rebalancing (1시간)**

```
Timeline:
├─ 15:30-15:40: Baseline 측정
│  ├─ 10K concurrent 시작
│  ├─ Latency baseline (10초)
│  └─ P99 기록
│
├─ 15:40-16:40: Rebalancing 진행
│  ├─ 15:40: rebalance partition 3 trigger
│  ├─ 15:40-16:40: 실시간 모니터링
│  │  ├─ P99 latency (매초)
│  │  ├─ Throughput (ops/sec)
│  │  ├─ Timeout count
│  │  └─ CPU/Memory
│  ├─ 16:30: rebalance 완료 예상
│  └─ 16:40: Post-rebalance 측정
│
└─ 결과 저장

Result File: test3_rebalancing.json
├─ baseline_p99
├─ during_p99
├─ post_p99
├─ rebalance_duration
├─ throughput_maintained
└─ overall_pass (Y/N)

PASS Criteria:
├─ During P99 < 100ms ✓
├─ Timeout rate < 1% ✓
├─ Throughput > 80% ✓
├─ Rebalance < 120s ✓
```

#### **16:30 - 17:30: Day 2 분석**

```
□ Test 1-3 결과 종합
├─ 모든 PASS?
├─ 이상 감지?
├─ 다음 테스트 가능?

Decision:
├─ ALL PASS → Test 4 진행
├─ PARTIAL → 분석 후 결정
└─ FAIL → STOP, root cause
```

---

### **Wed 2026-03-05 (계속)**

#### **17:30 - 18:30: Test 4 준비 (Resource Stress)**

```
□ cgroup 메모리 제한 설정
  ├─ Memory limit: 1.9GB (95% of 2GB)
  ├─ Test: 프로세스 재시작 시 확인
  └─ OOM killer 비활성화

□ fio disk slowdown 준비
  ├─ fio 설치 확인
  ├─ fsync 10배 느려짐 스크립트
  ├─ 테스트 실행 (5초)
  └─ 확인 후 중지

Status: ✅ Ready for Test 4
```

#### **18:30 - 20:00: Test 4a - Memory Stress (1.5시간)**

```
Timeline:
├─ 18:30: cgroup 95% 설정
├─ 18:31-19:30: INSERT 시도
│  ├─ 1000개 벡터 삽입 시도
│  ├─ Memory usage 매초 기록
│  ├─ OOM exception 감지
│  └─ Success count
│
├─ 19:30: Memory stress 해제
├─ 19:31-19:45: Recovery 측정
│  ├─ Memory 사용량 정상화
│  └─ 시간 기록
│
└─ 19:45-20:00: 분석

Result File: test4a_memory_stress.json
├─ oom_exceptions
├─ insert_success_rate
├─ recovery_time
└─ overall_pass (Y/N)

PASS Criteria:
├─ OOM exceptions = 0 ✓
├─ System crashes = 0 ✓
├─ Recovery time < 30s ✓
├─ Insert success > 99% ✓
```

#### **20:00 - 21:30: Test 4b - Disk Slowdown (1.5시간)**

```
Timeline:
├─ 20:00: fio 시작 (10배 느려짐)
├─ 20:01-21:00: INSERT 시도
│  ├─ 1000개 벡터 삽입 시도
│  ├─ fsync latency 기록
│  ├─ Insert latency 기록
│  └─ Timeout count
│
├─ 21:00: fio 중지
├─ 21:01-21:15: Recovery 측정
│  └─ Latency 정상화 확인
│
└─ 21:15-21:30: 분석

Result File: test4b_disk_slowdown.json
├─ insert_success_rate
├─ fsync_latency
├─ recovery_time
└─ overall_pass (Y/N)

PASS Criteria:
├─ Insert success > 99% ✓
├─ Recovery < 30s ✓
├─ No crashes = 0 ✓
```

#### **21:30 - 22:00: Test 4 최종 분석**

```
□ Memory + Disk 결과 종합
├─ 모두 PASS?
├─ 어느 것이 더 영향?
└─ Next action

Decision:
├─ ALL PASS → Week 6 72h 진행
├─ FAIL → root cause 분석
```

---

### **Thu 2026-03-06**

#### **09:00 - 10:00: Test 5 - Real Data Validation (1시간)**

```
Preparation:
├─ OpenAI embedding 1M 데이터 준비
│  ├─ 1,536 dimensions
│  ├─ float32 format
│  └─ 총 ~6GB (raw), ~2.2GB (compressed)

Timeline:
├─ 09:00: Insert 1M vectors 시작
├─ 09:XX: Insert 진행 (시간 기록)
│  ├─ 매 100K마다 progress 기록
│  └─ Memory usage 추적
│
├─ 10:00: Insert 완료 (또는 timeout)
├─ 10:01: Search performance 측정
│  ├─ 1K random queries
│  ├─ Latency 수집
│  ├─ P99 계산
│  └─ 정확도 검증
│
└─ 10:15: 결과 저장

Result File: test5_real_data.json
├─ insert_success
├─ insert_duration
├─ search_latency_p99
├─ memory_accuracy
├─ compression_ratio
└─ overall_pass (Y/N)

PASS Criteria:
├─ Insert success = 100% ✓
├─ Search P99 < 20ms ✓
├─ Memory accuracy > 80% ✓
├─ Compression ≈ 0.35 ✓
```

#### **10:00 - 11:00: Test 5 분석 및 최종 결정**

```
□ 모든 Test 1-5 결과 검토
├─ Test 1: Leadership ✓/✗
├─ Test 2: Network ✓/✗
├─ Test 3: Rebalancing ✓/✗
├─ Test 4: Resource ✓/✗
├─ Test 5: Real Data ✓/✗

Decision:
├─ ALL PASS (5/5) → Week 6 72h 진행
├─ 4/5 PASS → 실패한 것 분석 후 결정
├─ < 4 PASS → 심각한 문제, 코드 수정 필요

Next:
└─ 11:00: Week 6 72h Long-Running 시작 준비
```

---

### **Mon 2026-03-04 09:00 ~ Fri 2026-03-07 09:00: 72시간 Long-Running**

```
정확 타이밍: Mon 09:00 시작 → Fri 09:00 종료 (259,200초)

Tier 1: Warm-up (1시간)
├─ 09:00-10:00: Baseline 측정
└─ P99 latency ~8ms 기록

Tier 2: Sustained (26시간)
├─ 10:00 Mon ~ 12:00 Tue
├─ 10K ops/sec 지속
└─ 메모리 누수 초기 감지

Tier 3: Long-term (45시간)
├─ 12:00 Tue ~ 09:00 Fri
├─ 부하 변화 도입 (12시간마다 20K ops/sec)
└─ 최종 안정성 판정

Final: 09:00 Fri
└─ 모든 메트릭 수집 완료
```

---

### **Fri 2026-03-07 09:00 ~ 12:00: 최종 분석 및 판정**

```
09:00-10:00: 데이터 수집
├─ 72시간 로그 모두 수집
├─ 메트릭 데이터베이스 dump
├─ 시스템 상태 스냅샷
└─ 모든 노드 일관성 확인

10:00-11:00: 분석
├─ Memory leak rate 계산
├─ Latency degradation 분석
├─ Throughput stability 검증
├─ Error pattern 분석

11:00-12:00: 최종 판정
├─ Production Ready Criteria 적용
├─ 5가지 Category 평가
├─ 최종 PASS/CONDITIONAL/FAIL
└─ 보고서 생성

Output:
├─ FINAL_EVALUATION_REPORT.md
├─ metrics_analysis.json
├─ production_ready_verdict.txt
└─ GOGS 커밋
```

---

## 📊 **실시간 모니터링 명령어**

### **각 테스트 중 모니터링 (Terminal 1)**

```bash
# Memory & CPU
watch -n 1 'free -h && echo && ps aux | grep freelang'

# Network
iftop -n

# Disk I/O
iostat -x 1

# Process detail
htop

# Network packet
tcpdump -i eth0 -n
```

### **로그 수집 (Terminal 2)**

```bash
# Real-time log
tail -f /path/to/logs/system.log

# Error grep
tail -f /path/to/logs/error.log | grep -i "error\|exception\|fail"

# Metrics dump
watch -n 10 'cat current_metrics.json | jq .'
```

---

## 🔍 **실패 시 대응 (Quick Reference)**

```
Scenario 1: MTTR > 2000ms
├─ Raft election timeout 확인
├─ Network latency 확인
├─ Leader는 정상인가?
└─ 로그에서 heartbeat 메시지 추적

Scenario 2: Data loss 발생
├─ Quorum write 메커니즘 확인
├─ Replication lag 분석
├─ WAL 검토
└─ Snapshot consistency 확인

Scenario 3: Memory leak
├─ Heap dump 생성
├─ Object allocation 추적
├─ GC behavior 분석
└─ Memory profile 생성

Scenario 4: System hang
├─ Process status 확인
├─ Lock trace 수집
├─ Deadlock 분석
└─ Kill -9 + restart

Scenario 5: Network error
├─ Network packet capture
├─ Latency 측정
├─ Packet loss 확인
└─ Interface 상태 확인
```

---

## 📋 **최종 체크리스트**

```
Pre-Test (Mon 08:00):
□ 5개 노드 모두 alive
□ Clock sync (skew < 100ms)
□ Monitoring tools ready
□ Baseline metrics 측정
□ Logging infrastructure 준비

During Tests (Mon 09:00 - Fri 09:00):
□ 매시간 metrics 수집
□ 이상 감지 시 로그 기록
□ Critical error 시 즉시 보고

Post-Test (Fri 09:00 - 12:00):
□ 모든 로그 수집
□ 최종 분석 수행
□ Production Ready 판정
□ GOGS 커밋

Success Criteria:
□ Test 1-5: 모두 PASS
□ 72h: Memory leak < 1.5 MB/h
□ 72h: Latency increase < 10%
□ 72h: Throughput > 95%
□ 72h: Error rate < 0.1%
```

---

**시작**: 2026-03-04 09:00 (정확히)
**종료**: 2026-03-07 09:00 (72시간 후)
**목표**: Production Ready 최종 판정
**상태**: 🚀 실행 준비 완료

