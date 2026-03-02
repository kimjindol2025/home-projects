// ============================================
// Proxy Server Benchmarks
// ============================================
// 성능 측정: 요청 처리, 로드 밸런싱, 백엔드 선택

use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::time::Duration;

// 간단한 Proxy 시뮬레이션
struct SimpleProxy {
    backends: usize,
    current_index: Arc<AtomicUsize>,
}

impl SimpleProxy {
    fn new(backends: usize) -> Self {
        SimpleProxy {
            backends,
            current_index: Arc::new(AtomicUsize::new(0)),
        }
    }

    fn select_backend_round_robin(&self) -> usize {
        let index = self.current_index.fetch_add(1, Ordering::Relaxed);
        index % self.backends
    }

    fn select_backend_least_connections(&self, loads: &[usize]) -> usize {
        loads.iter()
            .enumerate()
            .min_by_key(|(_, &load)| load)
            .map(|(idx, _)| idx)
            .unwrap_or(0)
    }

    fn handle_request(&self, request_id: usize) -> u64 {
        // 요청 처리 시뮬레이션 (1-10ms)
        let delay = 1 + (request_id % 10) as u64;
        std::thread::sleep(Duration::from_millis(delay));
        delay
    }

    fn process_request_with_load_balancing(&self, loads: &[usize]) -> usize {
        let backend = self.select_backend_least_connections(loads);
        self.handle_request(backend);
        backend
    }
}

fn proxy_backend_selection_round_robin(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_backend_selection");
    group.sample_size(10000);

    group.bench_function("round_robin", |b| {
        let proxy = black_box(SimpleProxy::new(3));
        b.iter(|| {
            proxy.select_backend_round_robin()
        });
    });

    group.finish();
}

fn proxy_backend_selection_least_connections(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_backend_selection_lc");
    group.sample_size(1000);

    for backends in [3, 5, 10].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(backends),
            backends,
            |b, &backends| {
                b.iter(|| {
                    let proxy = black_box(SimpleProxy::new(backends));
                    let loads = vec![0; backends];
                    proxy.select_backend_least_connections(&black_box(loads))
                });
            },
        );
    }

    group.finish();
}

fn proxy_request_handling(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_request_handling");
    group.sample_size(100);
    group.measurement_time(Duration::from_secs(20));

    group.bench_function("single_request", |b| {
        let proxy = black_box(SimpleProxy::new(3));
        b.iter(|| {
            proxy.handle_request(black_box(1))
        });
    });

    group.finish();
}

fn proxy_concurrent_handling(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_concurrent");
    group.sample_size(50);
    group.measurement_time(Duration::from_secs(30));

    for requests in [10, 50, 100].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(requests),
            requests,
            |b, &requests| {
                b.iter(|| {
                    let proxy = black_box(SimpleProxy::new(3));
                    let mut total = 0u64;
                    for i in 0..requests {
                        total += proxy.handle_request(i);
                    }
                    black_box(total)
                });
            },
        );
    }

    group.finish();
}

fn proxy_weighted_distribution(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_weighted");
    group.sample_size(1000);

    group.bench_function("distribution", |b| {
        let proxy = black_box(SimpleProxy::new(3));
        let loads = black_box(vec![5, 3, 2]); // 가중치: 5:3:2
        b.iter(|| {
            proxy.select_backend_least_connections(&loads)
        });
    });

    group.finish();
}

fn proxy_throughput_analysis(c: &mut Criterion) {
    let mut group = c.benchmark_group("proxy_throughput");
    group.sample_size(100);

    group.bench_function("requests_per_second", |b| {
        let proxy = SimpleProxy::new(3);
        let start = std::time::Instant::now();
        let mut request_count = 0;

        b.iter(|| {
            request_count += 1;
            let _ = proxy.select_backend_round_robin();
        });

        let elapsed = start.elapsed();
        let throughput = request_count as f64 / elapsed.as_secs_f64();
        println!("Throughput: {:.0} requests/sec", throughput);
    });

    group.finish();
}

criterion_group!(
    benches,
    proxy_backend_selection_round_robin,
    proxy_backend_selection_least_connections,
    proxy_request_handling,
    proxy_concurrent_handling,
    proxy_weighted_distribution,
    proxy_throughput_analysis
);
criterion_main!(benches);
