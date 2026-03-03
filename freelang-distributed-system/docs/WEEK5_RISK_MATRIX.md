# 📊 Week 5 실제 환경 Chaos Test - 리스크 매트릭스

**작성일**: 2026-03-03
**적용 기간**: 2026-03-04 ~ 2026-03-08 (Week 5)
**목표**: 모든 리스크 식별 및 완화 전략 수립

---

## 🎯 매트릭스 개요

```
           Severity (심각도)
              ↑
         5   │  [위험영역]
             │  ・Leadership MTTR Fail
         4   │  ・Network Partition Violation
             │  ・Resource OOM
         3   │  ・Rebalance P99 > 100ms
             │  ・Memory Leak Detection
         2   │  ・Minor Timeout
             │  ・Latency Spike
         1   │  ・Log Message
             │
             └──────────────────────→ Probability
             1    2    3    4    5

빨강 (Critical):   Severity ≥ 4 AND Probability ≥ 3
주황 (High):       Severity ≥ 3 AND Probability ≥ 3
노랑 (Medium):     Severity ≥ 2 AND Probability ≥ 2
초록 (Low):        나머지
```

---

## 📋 Test별 상세 리스크 분석

### **Test 1: Raft Leadership Kill (Mon-Tue)**

목표: MTTR < 1000ms (평균), 데이터 손실 0건

#### Risk 1.1: MTTR 목표 초과

```
Scenario:
├─ Leader kill 후
├─ 재선출 대기 (network detection + election)
├─ 최악: 2000ms 초과
└─ 원인: 느린 heartbeat timeout (default 150ms)

Probability:   2/5 (학생 환경에서 네트워크 지연 가능)
Severity:      4/5 (SLA 위반)
Detection:     즉시 (MTTR 측정)
Detectability: 높음 (명확한 수치)

Mitigation:
├─ Pre-test: Heartbeat timeout 확인 (< 150ms)
├─ During: MTTR 실시간 수집
├─ Threshold: 1500ms 도달 시 즉시 중단
└─ Root Cause: 재선출 시간 분석

Fallback:
├─ Raft term 기록 분석
├─ Gossip 메시지 추적
├─ Candidate vote 지연 분석
└─ 네트워크 RTT 확인
```

**임계값 설정**:
```
MTTR 측정:
├─ Target:  847ms (시뮬레이션)
├─ OK:      < 1000ms ✓
├─ Warn:    1000-1500ms ⚠️
└─ Fail:    > 1500ms ❌

롤백 계획:
├─ 1500ms 초과 시 즉시 중단
├─ 재선출 로그 수집
├─ Raft configuration 재확인
└─ 네트워크 상태 진단
```

---

#### Risk 1.2: 데이터 손실 (Write Quorum 실패)

```
Scenario:
├─ Leader kill 후
├─ 일부 replicas에만 쓰여진 상태 (pre-commit)
├─ 새 leader에서 불일치 감지
└─ 최악: 100개 요청 손실

Probability:   1/5 (Raft quorum이 보호함)
Severity:      5/5 (데이터 무결성)
Detection:     명시적 검증 필요

Mitigation:
├─ Quorum size = 3 (5-node 클러스터)
├─ Pre-commit log 검증
├─ Leader election 후 committed index 확인
└─ 손실된 write 집계

Fallback:
├─ WAL (Write Ahead Log) 검토
├─ Replication leader → follower 메시지 확인
├─ Snapshot consistency 검증
└─ 데이터 복구 가능성 진단
```

**검증 체크리스트**:
```
각 Leadership Kill 후:
□ 새 Leader 선출 확인
□ Leader epoch 증가 확인
□ Pre-commit log 길이 기록
□ Replicas 간 log match 확인
□ Insert count 일관성 검증
  ├─ node1: X
  ├─ node2: X
  ├─ node3: X
  ├─ node4: X
  ├─ node5: X
  └─ 모두 동일 → OK ✓
□ Cumulative loss < 200 요청
```

---

#### Risk 1.3: 정상 노드 crash (연쇄 장애)

```
Scenario:
├─ Leader kill
├─ 재선출 중 follower가 과부하로 crash
├─ Quorum 상실 (2/5만 살아있음)
└─ 시스템 정지

Probability:   1/5 (과부하 테스트는 Fri인 Test 4)
Severity:      5/5 (서비스 마비)
Detection:     Health check

Mitigation:
├─ Memory/CPU monitoring (실시간)
├─ Election flood 제어 (max 3 candidates)
├─ Follower crash 조기 감지
└─ Cascading failure 예방

Fallback:
├─ Raft election flood 감지 시 중단
├─ OS process monitor (top, htop)
├─ Memory 사용량 > 80% 시 조사
└─ Core dump 생성
```

---

### **Test 2: Network Partition (Wed)**

목표: 일관성 위반 0건, P99 < 100ms, Timeout < 1%

#### Risk 2.1: 일관성 위반 (Split Brain)

```
Scenario:
├─ 30% packet loss + 200ms latency
├─ Partition 1 (2 nodes): 새로운 leader 선출
├─ Partition 2 (3 nodes): 기존 leader 유지
└─ 동일한 key에 다른 값 → 일관성 위반

Probability:   2/5 (네트워크 불안정)
Severity:      5/5 (데이터 무결성)
Detection:     명시적 검증 필요
Detectability: 매우 높음 (스캔 후 비교)

Mitigation:
├─ Quorum-based write (3/5 반드시 필요)
├─ Split brain이 발생해도 minority partition은 write 불가
├─ Heartbeat + election timeout으로 partition 감지
└─ Asymmetric partition 설계

Fallback:
├─ 모든 노드 데이터 scan (key/value count)
├─ Checksum 비교
├─ Vector의 vectorId 일관성 확인
├─ Replication log 분석
└─ 충돌 해결 전략 (last-write-wins vs custom)
```

**검증 절차**:
```
Network partition 적용 → 30초 운영 → 제거 → 복구 검증

매 5초마다 모든 노드에 INSERT 실행:
├─ node1 <- INSERT(key="test", value=X)
├─ node2 <- 같은 요청
├─ ...
└─ 모든 노드에서 최종 값이 동일한가?

최종 검증:
□ 모든 노드에서 데이터 scan
□ 동일 key에 서로 다른 값 존재? → FAIL
□ Insert count 일관성? → MUST match
□ Replication history? → Quorum majority 이겨야 함
```

---

#### Risk 2.2: Timeout 폭증

```
Scenario:
├─ 200ms latency
├─ Quorum write timeout (default 500ms)
├─ 50% 이상의 요청 timeout
└─ 사용자 체험 악화

Probability:   3/5 (200ms는 실제로 느림)
Severity:      3/5 (서비스는 되지만 느림)
Detection:     Timeout count 수집
Detectability: 즉시 (메트릭)

Mitigation:
├─ Timeout threshold 확인 (< 500ms)
├─ Retry logic 검증
├─ Exponential backoff 적용
└─ Client-side timeout 설정

Fallback:
├─ 실제 timeout rate 측정
├─ P99 latency 기록
├─ Timeout 원인 분석 (quorum wait? write? fsync?)
└─ Timeout threshold 동적 조정 고려
```

**모니터링 지표**:
```
Network Partition 중:

매 초:
├─ Request count: N
├─ Success count: N_success
├─ Timeout count: N_timeout
├─ Timeout rate: (N_timeout / N) × 100%

누적:
├─ P50 latency
├─ P99 latency
├─ P99.9 latency
└─ Max latency

임계값:
├─ Timeout rate < 1% ✓
├─ P99 < 100ms ✓
├─ Max < 500ms ✓
└─ 모두 통과 → OK
```

---

#### Risk 2.3: Rebalancing 자동 trigger (최악)

```
Scenario:
├─ Network partition 감지
├─ 클러스터 health check 실패
├─ 자동 rebalancing trigger
├─ 이미 느린 네트워크에서 rebalance + partition
└─ 이중 부하

Probability:   1/5 (rebalance trigger 로직 있어야 함)
Severity:      4/5 (성능 급락)
Detection:     Rebalance log 확인
Detectability: 높음

Mitigation:
├─ Partition 감지 중 rebalancing disable
├─ Manual trigger only during test
├─ Recovery 후에만 자동 rebalancing
└─ Configuration 사전 확인

Fallback:
├─ Rebalance 로그 수집
├─ Trigger 원인 분석
├─ Manual rebalancing 중단 가능성
└─ 시간 연장 고려
```

---

### **Test 3: Shard Rebalancing (Thu)**

목표: P99 < 100ms, Timeout < 1%, Throughput > 80%, Duration < 120s

#### Risk 3.1: P99 Latency 급증

```
Scenario:
├─ 10K 동시 연결
├─ Shard partition 3 마이그레이션 (node1 → node2)
├─ Lock contention 발생
├─ Latency spike: 12ms → 500ms+
└─ SLA 위반

Probability:   3/5 (lock-free 설계여도 risk)
Severity:      4/5 (SLA 위반)
Detection:     Latency percentile 수집
Detectability: 즉시 (메트릭)

Mitigation:
├─ Rebalance 중 partition 3 read 레플리카 우회
├─ Lock timeout 설정 (100ms)
├─ Partition 마이그레이션 batch size (100 vectors/batch)
├─ Rebalance throttling (max 10MB/s)
└─ P99 threshold 모니터링 (100ms)

Fallback:
├─ Rebalance 일시 중단 가능
├─ Hot partition 우선 제외
├─ 느린 시간대 우회 (node2가 과부하면 다른 노드로)
├─ Duration 연장 허용 (< 300s)
└─ Manual rebalance 필요 시 개입
```

**실시간 모니터링**:
```
Rebalance 진행 중 (60초):

매 1초:
├─ Request latency (P50, P99, P99.9)
├─ Rebalance progress (%)
├─ Partition 3 lock wait time
├─ Data transferred (MB)

기준:
├─ P99 < 100ms ← 초과 시 즉시 경고
├─ Lock wait < 50ms
├─ Transfer speed > 5MB/s (progress 확인)
└─ 비정상 시 pause 신호

Post-rebalance (10초):
├─ P99 복구 확인
├─ Consistency check (partition 3 데이터)
└─ Node 간 load balance
```

---

#### Risk 3.2: Rebalance Timeout (장시간 block)

```
Scenario:
├─ 데이터 양이 예상보다 큼
├─ 네트워크 속도 저하
├─ Rebalance > 120초 지속
└─ "시스템이 정지됨" 체감

Probability:   2/5 (시뮬레이션에서 67s, 실제는?)
Severity:      4/5 (운영 마비)
Detection:     Duration 측정
Detectability: 높음 (명확한 수치)

Mitigation:
├─ Rebalance 사전 크기 추정
├─ Partition 3 크기 확인 (테스트 전)
├─ Network bandwidth 확인 (iperf)
├─ Batch size 동적 조정
└─ Duration threshold 설정 (< 120s)

Fallback:
├─ Rebalance 진행률 매 30초 기록
├─ 평균 속도로 총 시간 예측
├─ > 150초 예상 시 중단 및 분석
├─ Partition 3 일부만 먼저 마이그레이션
└─ Multi-step rebalance (여러 batch)
```

---

#### Risk 3.3: Throughput 급락

```
Scenario:
├─ 10K concurrent requests
├─ Rebalance로 CPU 75% 사용
├─ Request processing CPU 25% 남음
├─ Throughput 50% → 80% 요구
└─ 실제: 60% 달성 → 실패

Probability:   2/5 (리소스 경쟁)
Severity:      3/5 (SLA 위반이지만 낮음)
Detection:     Throughput 측정
Detectability: 높음

Mitigation:
├─ CPU throttling 설정 (rebalance max 50%)
├─ Request queue 우선순위 (user request > rebalance)
├─ Batch processing 최적화
└─ Throughput threshold monitoring

Fallback:
├─ 실제 throughput 측정
├─ CPU usage 분석
├─ 병목 부분 식별 (disk I/O? memory? network?)
├─ Rebalance speed 조정
└─ 필요 시 duration 연장 협의
```

---

### **Test 4: Resource Stress (Fri)**

목표: OOM 0건, Crash 0건, Recovery < 30s, Insert Success > 99%

#### Risk 4.1: OOM Exception 발생

```
Scenario:
├─ 메모리 95% 압박
├─ Insert 요청 대기
├─ 새 vector allocation 시도
├─ malloc() 실패 → OOM
└─ Process exit (또는 exception)

Probability:   3/5 (95%는 매우 위험)
Severity:      5/5 (시스템 정지)
Detection:     OOM exception 또는 process exit
Detectability: 즉시

Mitigation:
├─ Cgroup memory limit 설정 (2GB, 95% = 1.9GB)
├─ Insert 전 메모리 체크
├─ Vector pool pre-allocation (GC 회피)
├─ OOM killer 비활성화 (graceful shutdown 위해)
└─ Out-of-memory handler (log + exit)

Fallback:
├─ Memory usage 매초 기록
├─ 메모리 누수 확인 (linear increase?)
├─ GC 효율 분석
├─ Vector compression 재확인
└─ 메모리 할당 패턴 분석
```

**상세 모니터링**:
```
Memory Stress 적용 (cgroup 95%로 제한):

매 1초:
├─ Memory used / total
├─ Memory used %
├─ Available memory
├─ Slab caches
├─ Page cache

Insert 시도:
├─ Insert 성공/실패 count
├─ 시도한 allocation size
├─ 실제 할당된 크기
├─ Allocation 시간

OOM 감지 시:
□ dmesg 수집 (kernel OOM message)
□ /proc/meminfo 기록
□ ps aux (메모리 사용 프로세스)
□ 정확한 OOM 시점 기록
□ Stack trace 수집
```

---

#### Risk 4.2: Disk I/O Stall (10배 느려짐)

```
Scenario:
├─ fio로 disk I/O 10배 느려짐 (fsync 1ms → 10ms)
├─ Raft log write block
├─ Replication latency 급증
├─ Timeout 발생
└─ 실제 시스템: 10배 느려짐

Probability:   3/5 (disk 느려짐 실제로 일어남)
Severity:      4/5 (모든 write 느려짐)
Detection:     Latency 급증
Detectability: 높음

Mitigation:
├─ fsync interval 조정 (batch)
├─ WAL (Write Ahead Log) batch size 증가
├─ Disk buffer 활용
├─ Timeout threshold 증가 (임시)
└─ Graceful degradation (write throughput 감소 허용)

Fallback:
├─ Actual disk latency 측정 (iostat)
├─ fsync duration 수집
├─ I/O wait time %
├─ Replication timeout 설정 동적 조정
├─ Insert queue 크기 증가
└─ 최악의 경우 write 불가 상태 detection
```

**I/O 성능 모니터링**:
```
Disk Slowdown 적용 (fio로 10배):

매 1초:
├─ fsync 호출 횟수
├─ fsync 평균 시간 (ms)
├─ fsync max 시간
├─ Queue depth
├─ I/O wait time (%)

Insert 영향:
├─ Insert latency (P99)
├─ Insert timeout
├─ WAL write latency
├─ Replication latency

임계값:
├─ fsync > 50ms 시 경고
├─ Insert P99 > 500ms 시 경고
├─ Insert timeout > 10% 시 중단
└─ Recovery 후 확인
```

---

#### Risk 4.3: Cascading Failure (메모리 + 디스크)

```
Scenario:
├─ 메모리 95% + 디스크 10배 느려짐 (동시)
├─ GC 강제 (메모리 회수)
├─ GC 중 fsync 대기 (디스크 느림)
├─ 전체 시스템 정지
└─ 최악: crash

Probability:   2/5 (동시 발생은 드물지만)
Severity:      5/5 (시스템 마비)
Detection:     System hang
Detectability: 매우 높음 (응답 없음)

Mitigation:
├─ 각 test 독립 실행 (메모리 + 디스크는 별도)
├─ Timeout 설정 (request max 1000ms)
├─ Watchdog timer (system hang detection)
├─ Force kill 메커니즘
└─ Recovery script 사전 준비

Fallback:
├─ Process kill + restart
├─ Memory dump 수집
├─ System log 분석
├─ Deadlock 분석 (lock trace)
└─ 개별 test로 다시 실행
```

---

## 🎬 실행 계획 (Week 5)

### Mon-Tue: Test 1 (Leadership Kill)

```
Mon 09:00
├─ Health check (모든 노드 alive?)
├─ Baseline latency 측정 (10초)
├─ MTTR 측정 준비 (timestamp 동기화)
└─ Log 파일 prepare

Mon 09:30 - 18:00
├─ Leadership kill 100회 실행
├─ 각 회차: kill → MTTR 측정 → 복구 → wait 5s
├─ MTTR 기록 (min, max, avg, p99)
├─ 데이터 일관성 검증 (각 50회마다)
└─ 비정상 발생 시 즉시 중단

Tue 09:00 - 12:00
├─ 추가 50회 반복
├─ 누적 통계 분석
├─ 실패 케이스 분석
└─ 롤백/복구 검증

Tue 14:00
├─ 최종 결과 정리
├─ 로그 분석
└─ MTTR 통계: min, max, avg, std, p99
```

**즉시 중단 조건**:
```
□ MTTR > 2000ms (3회 연속)
□ 데이터 불일치 감지 (1건이라도)
□ 노드 crash (leader 아님)
□ 네트워크 error (node 불가 도달)
□ Memory leak 발생 (RSS 50MB+ 증가)
```

---

### Wed: Test 2 (Network Partition)

```
Wed 09:00
├─ Baseline 측정 (normal latency)
├─ Network setup 확인
└─ tc/iptables 설정 검증

Wed 09:30
├─ Network 정상 상태 INSERT 1000개 (baseline)
├─ Latency 수집 (P50, P99, P99.9)

Wed 09:45
├─ Partition 적용: 30% loss + 200ms delay
├─ INSERT/SEARCH 병렬 실행 (30초)
├─ Latency 실시간 수집
├─ Consistency violation 감지

Wed 10:15
├─ Partition 해제
├─ Recovery latency 측정 (10초)
├─ 최종 데이터 consistency 검증

Wed 10:30
├─ 로그 분석
├─ Partition 동안의 timeout 분석
└─ Split brain 발생 여부 확인
```

**검증 기준**:
```
□ Consistency violations == 0
□ P99 latency (normal) < 15ms
□ Timeout rate < 1%
□ Quorum write success > 99%
□ 복구 후 latency 복원
```

---

### Thu: Test 3 (Shard Rebalancing)

```
Thu 09:00
├─ 10K concurrent connections 준비 (Apache Bench/wrk)
├─ Baseline latency 측정

Thu 09:30
├─ Rebalance trigger (partition 3: node1 → node2)
├─ 실시간 모니터링:
│  ├─ P99 latency
│  ├─ Throughput (ops/sec)
│  ├─ CPU usage
│  └─ Memory usage

Thu 10:30
├─ Rebalance 완료 확인 (duration 기록)
├─ Post-rebalance latency 측정
└─ Load balance 확인

Thu 11:00
├─ Consistency check
├─ Partition 3 데이터 검증
└─ 최종 결과 정리
```

**실패 시 롤백**:
```
□ P99 > 150ms 지속 (10s+) → 중단 + rollback
□ Timeout > 2% → 중단
□ Throughput < 50% → 중단
□ Rebalance > 180s → 강제 종료
```

---

### Fri: Test 4 (Resource Stress)

```
Fri 09:00
├─ Cgroup memory 95% 설정 (1.9GB / 2GB)
├─ fio disk slowdown 적용 안 함 (별도)
├─ Memory 할당 모니터링 준비

Fri 09:30 - 10:30
├─ Memory stress test만 (메모리 95%)
├─ 1000개 INSERT 시도
├─ Memory usage 실시간 수집
├─ OOM exception 감지

Fri 10:30
├─ Memory stress 해제
├─ Recovery time 측정 (< 30s)
├─ 메모리 누수 여부 확인

Fri 11:00
├─ Disk slowdown 적용 (fio 10배)
├─ 1000개 INSERT 시도
├─ fsync latency 수집

Fri 12:00
├─ Disk slowdown 해제
├─ 최종 데이터 검증
└─ 결과 정리
```

**OOM 감지 시**:
```
□ dmesg 수집
□ Process status (ps aux)
□ Memory map (pmap)
□ Core dump (if available)
□ Stack trace
□ 즉시 복구: kill + restart
```

---

## 📊 누적 모니터링 대시보드

### 매 Test 후 수집 항목

```
Test 1 결과:
├─ MTTR 통계 (min, max, avg, std, p50, p99)
├─ 손실된 요청 수
├─ 데이터 무결성 (일관성 위반 0/100)
└─ 리더 선출 성공률

Test 2 결과:
├─ Latency 통계 (normal vs partition)
├─ Timeout count
├─ Consistency violations
└─ Recovery 시간

Test 3 결과:
├─ Rebalance duration
├─ P99 latency (during, post)
├─ Throughput maintained (%)
├─ CPU/Memory peak
└─ Load balance score

Test 4 결과:
├─ OOM exceptions
├─ System crashes
├─ Recovery time
├─ Insert success rate
└─ Memory leak slope (MB/hour)
```

---

## ⚠️ 긴급 중단 기준 (Hard Stops)

```
ANY test 중 다음 발생 시 즉시 중단:

1. Data inconsistency (한 건이라도)
   └─ 동일 key 다른 값 또는 count 불일치

2. OOM exception
   └─ Process exit 또는 malloc() failure

3. System crash
   └─ Kernel panic 또는 process segfault

4. Data loss
   └─ Committed write가 사라짐

5. Network unreachable
   └─ 노드 간 통신 완전 단절

6. Hanging (> 5분)
   └─ 응답 완전 중지

발생 시:
├─ 즉시 log 수집
├─ 시스템 상태 기록 (memory, cpu, network)
├─ Process dump
├─ 영향 범위 분석
└─ 근본 원인 분석 대기
```

---

## 📋 Check List (Test 실행 전)

```
Environment Setup:
□ 5-node 클러스터 구성 완료
□ 각 노드 2GB RAM 확보
□ 네트워크 1Gbps 연결 확인
□ Storage SSD 100GB 확보
□ 모든 노드 alive 확인
□ Clock sync (NTP)
□ Logging 준비 (모든 노드)

Software Setup:
□ FreeLang 최신 버전 배포
□ Configuration 적용 (timeout, heartbeat 등)
□ Monitoring tool 준비 (metrics collection)
□ tc/iptables 설정 스크립트
□ fio disk slowdown 스크립트
□ Kill signal 스크립트
□ Consistency checker 준비
□ Memory/CPU monitor 준비

Baseline:
□ 정상 상태 latency 측정
□ Normal throughput 측정
□ Network RTT 측정
□ Disk I/O speed 측정
□ Memory baseline 기록

Logging:
□ All logs to shared storage
□ Timestamp sync (UTC)
□ Log rotation 설정
□ Error log aggregation
```

---

## 🎯 Success Criteria (Week 5 완료)

```
ALL 5가지 다음 만족:

1. Leadership Resilience
   └─ MTTR < 1000ms (avg), data loss = 0

2. Network Partition Consistency
   └─ Violations = 0, timeout rate < 1%

3. Rebalancing SLA
   └─ P99 < 100ms, throughput > 80%

4. Resource Recovery
   └─ OOM = 0, recovery < 30s

5. Real Data Validation
   └─ 1M vectors success, accuracy ±10%

모두 통과 → Week 6 72h Long-Running 진행
하나라도 실패 → 원인 분석 + 코드 수정 + 재테스트
```

---

**생성일**: 2026-03-03
**적용**: 2026-03-04 ~ 2026-03-08
**목표**: 모든 리스크 식별 및 완화 전략 수립
**상태**: 🔨 Week 5 실행 대기

