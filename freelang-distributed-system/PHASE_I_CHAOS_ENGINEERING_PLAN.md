# Phase I: Chaos Engineering - 장애 주입 테스트 프레임워크

**목표**: 실제 장애 상황(네트워크 지연, 노드 다운, 디스크 오류)에서 시스템이 자동 복구되는지 검증

**기간**: 1주
**산출물**: 2,000줄 (Rust)
**핵심 철학**: "장애를 예측하지 말고 유도하라"

---

## 🎯 설계 목표

### 현재 상태
- ✅ Raft 합의 알고리즘 (이론적으로 견고)
- ✅ 2PC 트랜잭션 (개념적으로 안전)
- ⚠️ **그러나 실제 네트워크/디스크 장애에서 정말 복구되는가?**

### 검증할 시나리오 (7가지)

```
S1: 네트워크 지연 (100ms ~ 1000ms)
    → Raft 타임아웃 동작, Leader 선출
    → 2PC Prepare 재시도

S2: 부분 네트워크 단절 (특정 노드만 끊김)
    → Quorum 유지 시 계속 진행
    → Quorum 손실 시 대기

S3: 노드 완전 다운
    → Leader 변경, 자동 rebalance
    → 데이터 손실 없음

S4: 디스크 오류 (Write 실패)
    → Automatic rollback
    → 트랜잭션 실패 → 재시도

S5: 느린 데이터 복제 (1MB 데이터 10초 소요)
    → Raft 로그 누적
    → Snapshot으로 압축

S6: 부분 데이터 손상 (Checksum 불일치)
    → 자동 감지
    → 원본 복제본에서 복구

S7: 동시 다중 장애 (네트워크 + 디스크)
    → 조율된 장애 복구
    → 데이터 무결성 보장
```

---

## 📐 아키텍처

### Core Components

```rust
// 1. 장애 주입 플랫폼
pub struct ChaosInjector {
    targets: Vec<FailureTarget>,
    duration_ms: u64,
    recovery_timeout_ms: u64,
}

pub enum FailureTarget {
    NetworkLatency { node: u32, delay_ms: u64 },
    NetworkPartition { nodes: Vec<u32> },
    NodeCrash { node: u32 },
    DiskError { error_rate: f64 },
    SlowReplication { bytes_per_sec: u64 },
    CorruptedData { probability: f64 },
}

// 2. 장애 복구 검증
pub struct RecoveryValidator {
    initial_state: SystemState,
    failure_period: (Instant, Instant),
    recovery_metrics: RecoveryMetrics,
}

pub struct RecoveryMetrics {
    time_to_detect_ms: u64,       // 문제 감지까지의 시간
    time_to_recovery_ms: u64,      // 복구까지의 시간
    data_loss_items: u32,          // 손실된 항목 수
    operations_failed: u32,        // 실패한 연산 수
    operations_retried: u32,       // 재시도된 연산 수
}

// 3. 데이터 무결성 검증
pub struct IntegrityValidator {
    checksums: HashMap<String, u64>,
    snapshots: Vec<SystemSnapshot>,
}
```

---

## 📊 구현 계획 (5단계)

### Stage 1: 장애 주입 엔진 (400줄)

**파일**: `src/chaos/injector.rs`

```rust
// 1. 네트워크 지연 주입
pub struct NetworkDelayInjector {
    delays: Arc<DashMap<u32, u64>>, // node_id -> delay_ms
}

impl NetworkDelayInjector {
    pub async fn inject_latency(&self, node: u32, delay_ms: u64) {
        self.delays.insert(node, delay_ms);
        tokio::time::sleep(Duration::from_millis(delay_ms)).await;
    }

    pub fn get_latency(&self, node: u32) -> u64 {
        self.delays.get(&node).map(|v| *v).unwrap_or(0)
    }
}

// 2. 네트워크 단절 주입
pub struct NetworkPartitionInjector {
    partitions: Arc<DashMap<u32, bool>>, // node_id -> is_partitioned
}

impl NetworkPartitionInjector {
    pub async fn partition_node(&self, node: u32) {
        self.partitions.insert(node, true);
    }

    pub async fn heal_partition(&self, node: u32) {
        self.partitions.insert(node, false);
    }

    pub fn is_reachable(&self, node: u32) -> bool {
        !self.partitions.get(&node).map(|v| *v).unwrap_or(false)
    }
}

// 3. 노드 크래시 시뮬레이션
pub struct NodeCrashSimulator {
    crashed_nodes: Arc<DashSet<u32>>,
}

impl NodeCrashSimulator {
    pub async fn crash_node(&self, node: u32) {
        self.crashed_nodes.insert(node);
        // 해당 노드의 모든 pending RPC 실패 처리
        self.fail_pending_rpcs(node).await;
    }

    pub async fn restart_node(&self, node: u32) {
        self.crashed_nodes.remove(&node);
        // 노드 복구 로직: Raft log catch-up, state sync
        self.recover_node_state(node).await;
    }

    pub fn is_running(&self, node: u32) -> bool {
        !self.crashed_nodes.contains(&node)
    }
}

// 4. 통합 주입 관리자
pub struct ChaosOrchestrator {
    network_delay: Arc<NetworkDelayInjector>,
    partition: Arc<NetworkPartitionInjector>,
    crash: Arc<NodeCrashSimulator>,
}

impl ChaosOrchestrator {
    pub async fn inject_scenario(&self, scenario: ChaosScenario) {
        match scenario {
            ChaosScenario::HighLatency => {
                for node in 1..=3 {
                    self.network_delay.inject_latency(node, 200).await;
                }
            }
            ChaosScenario::PartialPartition => {
                self.partition.partition_node(2).await;
            }
            ChaosScenario::NodeDown => {
                self.crash.crash_node(2).await;
            }
        }
    }
}
```

### Stage 2: 복구 검증 (400줄)

**파일**: `src/chaos/recovery_validator.rs`

```rust
// 1. 시스템 상태 스냅샷
pub struct SystemSnapshot {
    timestamp: u64,
    leader_id: u32,
    raft_term: u64,
    log_index: u64,
    accounts: HashMap<String, Account>,
    transactions: Vec<Transaction>,
}

impl SystemSnapshot {
    pub async fn capture() -> Self {
        SystemSnapshot {
            timestamp: now_us(),
            leader_id: raft_state.current_leader,
            raft_term: raft_state.current_term,
            log_index: raft_state.last_log_index,
            accounts: bank_state.accounts.clone(),
            transactions: bank_state.transactions.clone(),
        }
    }
}

// 2. 복구 과정 모니터링
pub struct RecoveryMonitor {
    failure_start: Instant,
    detection_time: Option<Instant>,
    recovery_complete: Option<Instant>,
}

impl RecoveryMonitor {
    pub fn mark_detected(&mut self) {
        self.detection_time = Some(Instant::now());
    }

    pub fn mark_recovered(&mut self) {
        self.recovery_complete = Some(Instant::now());
    }

    pub fn metrics(&self) -> RecoveryMetrics {
        let detection_ms = self.detection_time
            .map(|t| t.duration_since(self.failure_start).as_millis() as u64)
            .unwrap_or(0);

        let recovery_ms = self.recovery_complete
            .map(|t| t.duration_since(self.failure_start).as_millis() as u64)
            .unwrap_or(0);

        RecoveryMetrics {
            time_to_detect_ms: detection_ms,
            time_to_recovery_ms: recovery_ms,
            ..Default::default()
        }
    }
}

// 3. 데이터 무결성 검증
pub struct IntegrityChecker {
    initial_snapshot: SystemSnapshot,
}

impl IntegrityChecker {
    pub async fn verify_after_recovery(&self) -> IntegrityReport {
        let final_snapshot = SystemSnapshot::capture().await;

        let mut report = IntegrityReport {
            data_loss: 0,
            data_corruption: 0,
            inconsistencies: Vec::new(),
        };

        // 계좌 검증
        for (id, account) in &self.initial_snapshot.accounts {
            if let Some(final_account) = final_snapshot.accounts.get(id) {
                if account.balance != final_account.balance {
                    report.data_loss += 1;
                    report.inconsistencies.push(
                        format!("Account {} balance mismatch", id)
                    );
                }
            }
        }

        // 트랜잭션 검증
        let initial_tx_count = self.initial_snapshot.transactions.len();
        let final_tx_count = final_snapshot.transactions.len();
        if initial_tx_count > final_tx_count {
            report.data_loss += initial_tx_count - final_tx_count;
        }

        report
    }
}
```

### Stage 3: 7가지 카오스 시나리오 (500줄)

**파일**: `src/chaos/scenarios.rs`

```rust
pub enum ChaosScenario {
    HighLatency,
    NetworkPartition,
    NodeCrash,
    DiskError,
    SlowReplication,
    DataCorruption,
    MultipleFailures,
}

// S1: 네트워크 지연
pub async fn scenario_high_latency(
    orchestrator: &ChaosOrchestrator,
    monitor: &RecoveryMonitor,
) -> TestResult {
    println!("🔴 S1: Injecting 200ms latency...");

    for node in 1..=3 {
        orchestrator.network_delay.inject_latency(node, 200).await;
    }

    // Raft 타임아웃으로 인한 Leader 선출 기대
    let result = wait_for_leader_change(Duration::from_secs(5)).await;

    // 이후 요청 성공 여부 확인
    let transfer = send_transfer_request().await;

    TestResult {
        passed: transfer.is_ok() && result.is_ok(),
        metrics: monitor.metrics(),
    }
}

// S2: 부분 네트워크 단절
pub async fn scenario_partition(
    orchestrator: &ChaosOrchestrator,
) -> TestResult {
    println!("🔴 S2: Network partition on Node 2...");

    orchestrator.partition.partition_node(2).await;

    // 3-node 클러스터에서 2개는 정족수 유지
    let quorum_ok = wait_for_quorum_state(Duration::from_secs(3)).await;

    // 격리된 노드는 write 실패
    let isolated_write = send_write_to_node(2).await;

    TestResult {
        passed: quorum_ok && isolated_write.is_err(),
        metrics: Default::default(),
    }
}

// S3: 노드 크래시
pub async fn scenario_node_crash(
    orchestrator: &ChaosOrchestrator,
) -> TestResult {
    println!("🔴 S3: Node 2 crash...");

    orchestrator.crash.crash_node(2).await;

    // Leader가 변경되었는가?
    let leader_changed = wait_for_leader_change(Duration::from_secs(5)).await;

    // 나머지 노드로 계속 운영되는가?
    let transfer_ok = send_transfer_request().await;

    // 노드 복구 후 동기화되는가?
    orchestrator.crash.restart_node(2).await;
    let sync_ok = wait_for_sync(Duration::from_secs(10)).await;

    TestResult {
        passed: leader_changed.is_ok() && transfer_ok.is_ok() && sync_ok.is_ok(),
        metrics: Default::default(),
    }
}

// S4: 디스크 오류
pub async fn scenario_disk_error(
    orchestrator: &ChaosOrchestrator,
) -> TestResult {
    println!("🔴 S4: Disk write error (10% probability)...");

    orchestrator.disk_error.inject_errors(0.1).await;

    // 일부 write 실패
    let mut success_count = 0;
    let mut fail_count = 0;

    for _ in 0..10 {
        match send_transfer_request().await {
            Ok(_) => success_count += 1,
            Err(_) => fail_count += 1,
        }
    }

    // 오류율이 예상치와 맞는가? (약 10%)
    let error_rate_ok = (fail_count as f64 / 10.0) > 0.05 &&
                        (fail_count as f64 / 10.0) < 0.15;

    // 모든 실패가 감지되고 로깅되었는가?
    let errors_logged = check_error_logs().await;

    TestResult {
        passed: error_rate_ok && errors_logged,
        metrics: Default::default(),
    }
}

// S5: 느린 복제
pub async fn scenario_slow_replication(
    orchestrator: &ChaosOrchestrator,
) -> TestResult {
    println!("🔴 S5: Slow replication (1MB in 10s)...");

    orchestrator.replication.set_speed_limit(1024 * 1024 / 10).await; // 100KB/s

    // 대량 데이터 전송
    for i in 0..100 {
        send_transfer_request().await.ok();
    }

    // Raft 로그 크기 증가
    let log_size = raft_state.log_size();

    // Snapshot으로 압축되는가?
    let snapshot_triggered = wait_for_snapshot(Duration::from_secs(15)).await;

    TestResult {
        passed: snapshot_triggered.is_ok() && log_size > 1_000_000,
        metrics: Default::default(),
    }
}

// S6: 데이터 손상
pub async fn scenario_data_corruption(
    validator: &IntegrityChecker,
) -> TestResult {
    println!("🔴 S6: Data bit-flip (0.01% probability)...");

    let initial_snapshot = SystemSnapshot::capture().await;

    // 일부 메모리 데이터 손상
    corrupt_random_bits(0.0001).await;

    // Checksum 검증으로 감지되는가?
    let corruptions_detected = count_detected_corruptions().await;

    // 자동 복구되는가?
    let recovered_ok = recover_from_replicas().await;

    TestResult {
        passed: corruptions_detected > 0 && recovered_ok.is_ok(),
        metrics: Default::default(),
    }
}

// S7: 다중 장애
pub async fn scenario_multiple_failures(
    orchestrator: &ChaosOrchestrator,
) -> TestResult {
    println!("🔴 S7: Network partition + Disk error + High latency...");

    // 동시에 여러 장애 주입
    orchestrator.partition.partition_node(2).await;
    orchestrator.disk_error.inject_errors(0.05).await;
    orchestrator.network_delay.inject_latency(1, 100).await;

    // 어떤 시나리오도 데이터 손실을 초래하지 않아야 함
    let integrity = validator.verify_after_recovery().await;

    TestResult {
        passed: integrity.data_loss == 0,
        metrics: Default::default(),
    }
}
```

### Stage 4: 테스트 실행 프레임워크 (400줄)

**파일**: `src/chaos/test_runner.rs`

```rust
pub struct ChaosTestSuite {
    scenarios: Vec<(String, Box<dyn ChaosTest>)>,
    results: Arc<DashMap<String, TestResult>>,
}

pub trait ChaosTest: Send + Sync {
    async fn run(&self) -> TestResult;
    fn name(&self) -> &str;
}

impl ChaosTestSuite {
    pub async fn run_all(&self) {
        println!("\n╔════════════════════════════════════╗");
        println!("║  Chaos Engineering Test Suite     ║");
        println!("╚════════════════════════════════════╝\n");

        for (name, test) in &self.scenarios {
            println!("Running: {}", name);
            let result = test.run().await;
            self.results.insert(name.clone(), result.clone());

            if result.passed {
                println!("  ✅ PASSED");
            } else {
                println!("  ❌ FAILED");
            }
            println!("  Time to detect: {}ms", result.metrics.time_to_detect_ms);
            println!("  Time to recovery: {}ms", result.metrics.time_to_recovery_ms);
            println!();
        }

        self.print_summary();
    }

    fn print_summary(&self) {
        let total = self.results.len();
        let passed = self.results.iter()
            .filter(|r| r.value().passed)
            .count();

        println!("\n╔════════════════════════════════════╗");
        println!("║  Test Summary                     ║");
        println!("╠════════════════════════════════════╣");
        println!("║ Total: {} tests", total);
        println!("║ Passed: {} ✅", passed);
        println!("║ Failed: {} ❌", total - passed);
        println!("║ Success Rate: {:.1}%", (passed as f64 / total as f64) * 100.0);
        println!("╚════════════════════════════════════╝");
    }
}

// CLI 실행
pub async fn run_chaos_tests() {
    let suite = ChaosTestSuite::new();
    suite.add_scenario(
        "S1: High Latency",
        Box::new(HighLatencyTest),
    );
    suite.add_scenario(
        "S2: Network Partition",
        Box::new(PartitionTest),
    );
    suite.add_scenario(
        "S3: Node Crash",
        Box::new(NodeCrashTest),
    );
    suite.add_scenario(
        "S4: Disk Error",
        Box::new(DiskErrorTest),
    );
    suite.add_scenario(
        "S5: Slow Replication",
        Box::new(SlowReplicationTest),
    );
    suite.add_scenario(
        "S6: Data Corruption",
        Box::new(DataCorruptionTest),
    );
    suite.add_scenario(
        "S7: Multiple Failures",
        Box::new(MultipleFailuresTest),
    );

    suite.run_all().await;
}
```

### Stage 5: 대시보드 & 리포트 (300줄)

**파일**: `src/chaos/dashboard.rs`

```rust
pub struct ChaosDashboard {
    results: Arc<DashMap<String, TestResult>>,
}

impl ChaosDashboard {
    pub fn generate_html_report(&self) -> String {
        let mut html = String::new();

        html.push_str(r#"
<!DOCTYPE html>
<html>
<head>
    <title>Chaos Engineering Report</title>
    <style>
        body { font-family: Arial; margin: 20px; }
        .scenario { border: 1px solid #ccc; padding: 10px; margin: 10px 0; }
        .passed { background-color: #90EE90; }
        .failed { background-color: #FFB6C6; }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 10px; text-align: left; }
    </style>
</head>
<body>
    <h1>🔴 Chaos Engineering Test Report</h1>
        "#);

        let total = self.results.len();
        let passed = self.results.iter()
            .filter(|r| r.value().passed)
            .count();

        html.push_str(&format!(
            "<h2>Summary: {} / {} tests passed ({:.1}%)</h2>",
            passed, total, (passed as f64 / total as f64) * 100.0
        ));

        html.push_str("<table>");
        html.push_str("<tr><th>Scenario</th><th>Status</th><th>Detect (ms)</th><th>Recovery (ms)</th><th>Data Loss</th></tr>");

        for result in self.results.iter() {
            let status = if result.value().passed { "✅ PASS" } else { "❌ FAIL" };
            let class = if result.value().passed { "passed" } else { "failed" };
            html.push_str(&format!(
                "<tr class='{}'>
                    <td>{}</td>
                    <td>{}</td>
                    <td>{}</td>
                    <td>{}</td>
                    <td>{}</td>
                </tr>",
                class,
                result.key(),
                status,
                result.value().metrics.time_to_detect_ms,
                result.value().metrics.time_to_recovery_ms,
                result.value().metrics.data_loss_items,
            ));
        }

        html.push_str("</table>");
        html.push_str("</body></html>");

        html
    }

    pub fn generate_json_report(&self) -> String {
        let mut results = Vec::new();

        for result in self.results.iter() {
            results.push(json!({
                "scenario": result.key(),
                "passed": result.value().passed,
                "metrics": {
                    "time_to_detect_ms": result.value().metrics.time_to_detect_ms,
                    "time_to_recovery_ms": result.value().metrics.time_to_recovery_ms,
                    "data_loss": result.value().metrics.data_loss_items,
                }
            }));
        }

        serde_json::to_string_pretty(&results).unwrap()
    }
}

// REST API 엔드포인트
pub async fn chaos_report_handler() -> String {
    let dashboard = get_dashboard();
    dashboard.generate_html_report()
}

pub async fn chaos_json_handler() -> String {
    let dashboard = get_dashboard();
    dashboard.generate_json_report()
}
```

---

## 📊 기대 효과

### 1. 자동 장애 복구 검증

```
Scenario: Node 2 crashes
├─ Detection time: 150ms (Raft heartbeat timeout)
├─ Recovery time: 2.3s (Leader election + log replication)
├─ Data loss: 0 items ✅
└─ Status: PASSED

Recommendation: Recovery time is within SLA (3s)
```

### 2. 성능 저하 조기 감지

```
Scenario: Network latency 200ms
├─ Initial: avg latency 50ms
├─ With 200ms network latency: avg latency 250ms (+400%)
├─ Alert threshold: 20%
└─ Action: AUTO-ALERT triggered
```

### 3. 데이터 무결성 보증

```
Scenario: Simultaneous network partition + disk error
├─ Data consistency: MAINTAINED ✅
├─ No items lost: TRUE ✅
├─ No items corrupted: TRUE ✅
└─ Status: Production-ready
```

---

## ✅ 최종 구조

```
Chaos Engineering System
├── Injector (injector.rs) - 장애 주입 엔진
├── Recovery Validator (recovery_validator.rs) - 복구 검증
├── Scenarios (scenarios.rs) - 7가지 카오스 시나리오
├── Test Runner (test_runner.rs) - 자동 테스트 실행
└── Dashboard (dashboard.rs) - 리포트 생성

→ 7가지 장애 상황에서 자동 복구 검증
→ 데이터 무결성 보장 (0 loss)
→ 성능 저하 조기 감지
→ 프로덕션 준비도 평가
```

---

## 🎓 박사 수준의 안정성 입증

이 두 Phase (H + I)가 완료되면:

1. **Observability**: 모든 지연을 추적하고 병목 지점 파악 가능
2. **Resilience**: 7가지 장애에서 자동 복구 검증
3. **Durability**: 데이터 손실 없음 보증 (99.99999% uptime)

결과: **"10년 무중단 운영 가능한 시스템"** ✅

---

**다음**: Phase H, I 구현 시작
