// ============================================
// Bank System Benchmarks
// ============================================
// 성능 측정: 계좌 관리, 송금, 동시 거래

use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};
use std::collections::HashMap;
use std::sync::Arc;
use parking_lot::RwLock;
use std::time::Duration;

// 간단한 Bank 시뮬레이션
#[derive(Clone)]
struct SimpleBank {
    accounts: Arc<RwLock<HashMap<String, i64>>>,
}

impl SimpleBank {
    fn new() -> Self {
        SimpleBank {
            accounts: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    fn create_account(&self, id: &str, balance: i64) {
        let mut accounts = self.accounts.write();
        accounts.insert(id.to_string(), balance);
    }

    fn get_balance(&self, id: &str) -> Option<i64> {
        let accounts = self.accounts.read();
        accounts.get(id).copied()
    }

    fn transfer(&self, from: &str, to: &str, amount: i64) -> bool {
        let mut accounts = self.accounts.write();

        if let Some(from_balance) = accounts.get_mut(from) {
            if *from_balance >= amount {
                *from_balance -= amount;
                if let Some(to_balance) = accounts.get_mut(to) {
                    *to_balance += amount;
                    return true;
                }
            }
        }
        false
    }

    fn deposit(&self, id: &str, amount: i64) -> bool {
        let mut accounts = self.accounts.write();
        if let Some(balance) = accounts.get_mut(id) {
            *balance += amount;
            return true;
        }
        false
    }

    fn withdraw(&self, id: &str, amount: i64) -> bool {
        let mut accounts = self.accounts.write();
        if let Some(balance) = accounts.get_mut(id) {
            if *balance >= amount {
                *balance -= amount;
                return true;
            }
        }
        false
    }
}

fn bank_account_creation(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_account_creation");
    group.sample_size(10000);

    group.bench_function("create_single_account", |b| {
        let bank = black_box(SimpleBank::new());
        bank.create_account("alice", 1000_00);

        b.iter(|| {
            bank.create_account(
                &format!("account_{}", b.iter().count()),
                100_00
            )
        });
    });

    group.finish();
}

fn bank_balance_query(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_balance_query");
    group.sample_size(50000);

    group.bench_function("get_balance", |b| {
        let bank = black_box(SimpleBank::new());
        bank.create_account("alice", 1000_00);

        b.iter(|| {
            bank.get_balance(black_box("alice"))
        });
    });

    group.finish();
}

fn bank_transfer(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_transfer");
    group.sample_size(1000);
    group.measurement_time(Duration::from_secs(15));

    group.bench_function("simple_transfer", |b| {
        let bank = black_box(SimpleBank::new());
        bank.create_account("alice", 1000_00);
        bank.create_account("bob", 500_00);

        b.iter(|| {
            bank.transfer(
                black_box("alice"),
                black_box("bob"),
                black_box(100_00)
            )
        });
    });

    group.finish();
}

fn bank_deposit(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_deposit");
    group.sample_size(10000);

    group.bench_function("deposit", |b| {
        let bank = black_box(SimpleBank::new());
        bank.create_account("alice", 0);

        b.iter(|| {
            bank.deposit(black_box("alice"), black_box(500_00))
        });
    });

    group.finish();
}

fn bank_withdraw(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_withdraw");
    group.sample_size(10000);

    group.bench_function("withdraw", |b| {
        let bank = black_box(SimpleBank::new());
        bank.create_account("alice", 10000_00);

        b.iter(|| {
            bank.withdraw(black_box("alice"), black_box(100_00))
        });
    });

    group.finish();
}

fn bank_concurrent_transfers(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_concurrent_transfers");
    group.sample_size(50);
    group.measurement_time(Duration::from_secs(20));

    for num_transfers in [10, 50, 100].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(num_transfers),
            num_transfers,
            |b, &num_transfers| {
                b.iter(|| {
                    let bank = black_box(SimpleBank::new());

                    // 계좌 초기화
                    for i in 0..num_transfers {
                        bank.create_account(&format!("account_{}", i), 1000_00);
                    }

                    // 동시 송금 시뮬레이션
                    let mut success = 0;
                    for i in 0..num_transfers {
                        let from = format!("account_{}", i);
                        let to = format!("account_{}", (i + 1) % num_transfers);
                        if bank.transfer(&from, &to, 100_00) {
                            success += 1;
                        }
                    }

                    black_box(success)
                });
            },
        );
    }

    group.finish();
}

fn bank_mixed_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_mixed_operations");
    group.sample_size(100);

    group.bench_function("realistic_workload", |b| {
        let bank = black_box(SimpleBank::new());

        // 초기화: 10개 계좌
        for i in 0..10 {
            bank.create_account(&format!("account_{}", i), 5000_00);
        }

        b.iter(|| {
            // 30% 조회
            let _ = bank.get_balance("account_0");

            // 20% 입금
            let _ = bank.deposit("account_1", 500_00);

            // 20% 출금
            let _ = bank.withdraw("account_2", 200_00);

            // 30% 송금
            let _ = bank.transfer("account_3", "account_4", 300_00);
        });
    });

    group.finish();
}

fn bank_throughput(c: &mut Criterion) {
    let mut group = c.benchmark_group("bank_throughput");
    group.sample_size(100);

    group.bench_function("operations_per_second", |b| {
        let bank = SimpleBank::new();

        // 초기화
        for i in 0..5 {
            bank.create_account(&format!("account_{}", i), 10000_00);
        }

        let start = std::time::Instant::now();
        let mut ops = 0;

        b.iter(|| {
            let _ = bank.transfer("account_0", "account_1", 100_00);
            ops += 1;
        });

        let elapsed = start.elapsed();
        let throughput = ops as f64 / elapsed.as_secs_f64();
        println!("Throughput: {:.0} operations/sec", throughput);
    });

    group.finish();
}

criterion_group!(
    benches,
    bank_account_creation,
    bank_balance_query,
    bank_transfer,
    bank_deposit,
    bank_withdraw,
    bank_concurrent_transfers,
    bank_mixed_operations,
    bank_throughput
);
criterion_main!(benches);
