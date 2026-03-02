// ============================================
// Raft Consensus Benchmarks
// ============================================
// 성능 측정: 리더 선출, 로그 복제, 노드 관리

use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};
use std::time::Duration;

// 간단한 Raft 시뮬레이션 (실제 모듈 import 시 수정)
#[derive(Clone)]
struct SimpleRaftCluster {
    nodes: usize,
}

impl SimpleRaftCluster {
    fn new(nodes: usize) -> Self {
        SimpleRaftCluster { nodes }
    }

    fn elect_leader(&mut self) -> u64 {
        // 리더 선출 시뮬레이션
        std::thread::sleep(Duration::from_micros(100));
        1
    }

    fn append_entry(&mut self, _entry: String) -> u64 {
        // 로그 항목 추가 시뮬레이션
        std::thread::sleep(Duration::from_micros(50));
        1
    }

    fn replicate_logs(&mut self, count: usize) -> usize {
        // 로그 복제 시뮬레이션
        std::thread::sleep(Duration::from_micros((count as u64) * 5));
        count
    }
}

fn raft_leader_election(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft_leader_election");
    group.sample_size(100);
    group.measurement_time(Duration::from_secs(10));
    group.warm_up_time(Duration::from_secs(2));

    for nodes in [3, 5, 7].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(nodes),
            nodes,
            |b, &nodes| {
                b.iter(|| {
                    let mut cluster = black_box(SimpleRaftCluster::new(nodes));
                    cluster.elect_leader()
                });
            },
        );
    }

    group.finish();
}

fn raft_log_append(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft_log_append");
    group.sample_size(1000);
    group.measurement_time(Duration::from_secs(10));

    group.bench_function("single_entry", |b| {
        b.iter(|| {
            let mut cluster = black_box(SimpleRaftCluster::new(5));
            cluster.append_entry(black_box("test_entry".to_string()))
        });
    });

    group.finish();
}

fn raft_log_replication(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft_log_replication");
    group.sample_size(50);
    group.measurement_time(Duration::from_secs(15));

    for entries in [10, 100, 1000].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(entries),
            entries,
            |b, &entries| {
                b.iter(|| {
                    let mut cluster = black_box(SimpleRaftCluster::new(5));
                    cluster.replicate_logs(black_box(entries))
                });
            },
        );
    }

    group.finish();
}

fn raft_throughput(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft_throughput");
    group.sample_size(100);

    group.bench_function("log_entries_per_second", |b| {
        let mut cluster = SimpleRaftCluster::new(5);
        let start = std::time::Instant::now();

        b.iter(|| {
            cluster.append_entry("entry".to_string())
        });

        let elapsed = start.elapsed();
        let ops = 100.0 / elapsed.as_secs_f64();
        println!("Throughput: {:.0} entries/sec", ops);
    });

    group.finish();
}

fn raft_node_scaling(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft_node_scaling");
    group.sample_size(20);

    for nodes in [3, 5, 7, 9].iter() {
        group.bench_with_input(
            BenchmarkId::new("nodes", nodes),
            nodes,
            |b, &nodes| {
                b.iter(|| {
                    let mut cluster = black_box(SimpleRaftCluster::new(nodes));
                    let _ = cluster.elect_leader();
                    cluster.replicate_logs(10)
                });
            },
        );
    }

    group.finish();
}

criterion_group!(
    benches,
    raft_leader_election,
    raft_log_append,
    raft_log_replication,
    raft_throughput,
    raft_node_scaling
);
criterion_main!(benches);
