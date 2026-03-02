// ============================================
// Security Benchmarks
// ============================================
// 성능 측정: JWT, 암호화, TLS

use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};
use std::time::Duration;

// 간단한 JWT 토큰 시뮬레이션
struct SimpleTokenValidator {
    secret: String,
}

impl SimpleTokenValidator {
    fn new(secret: &str) -> Self {
        SimpleTokenValidator {
            secret: secret.to_string(),
        }
    }

    fn generate_token(&self, user: &str) -> String {
        // 간단한 토큰 생성 시뮬레이션
        format!("{}.{}.{}", user, self.secret, chrono::Utc::now().timestamp())
    }

    fn validate_token(&self, token: &str) -> bool {
        // 간단한 토큰 검증 시뮬레이션
        token.contains(&self.secret)
    }
}

// 간단한 암호화 시뮬레이션
struct SimpleEncryption {
    key: Vec<u8>,
}

impl SimpleEncryption {
    fn new() -> Self {
        let key = vec![1, 2, 3, 4, 5];
        SimpleEncryption { key }
    }

    fn encrypt(&self, data: &[u8]) -> Vec<u8> {
        // 간단한 XOR 암호화
        data.iter()
            .enumerate()
            .map(|(i, &byte)| byte ^ self.key[i % self.key.len()])
            .collect()
    }

    fn decrypt(&self, data: &[u8]) -> Vec<u8> {
        // XOR은 자체 역함수
        self.encrypt(data)
    }
}

fn jwt_token_generation(c: &mut Criterion) {
    let mut group = c.benchmark_group("jwt_token_generation");
    group.sample_size(5000);

    group.bench_function("generate_token", |b| {
        let validator = black_box(SimpleTokenValidator::new("secret_key"));
        b.iter(|| {
            validator.generate_token(black_box("user@example.com"))
        });
    });

    group.finish();
}

fn jwt_token_validation(c: &mut Criterion) {
    let mut group = c.benchmark_group("jwt_token_validation");
    group.sample_size(5000);

    group.bench_function("validate_token", |b| {
        let validator = black_box(SimpleTokenValidator::new("secret_key"));
        let token = validator.generate_token("user@example.com");

        b.iter(|| {
            validator.validate_token(black_box(&token))
        });
    });

    group.finish();
}

fn encryption_speed(c: &mut Criterion) {
    let mut group = c.benchmark_group("encryption_speed");
    group.sample_size(1000);
    group.measurement_time(Duration::from_secs(15));

    for size in [64, 512, 4096].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(format!("{}_bytes", size)),
            size,
            |b, &size| {
                b.iter(|| {
                    let cipher = black_box(SimpleEncryption::new());
                    let data = vec![42u8; size];
                    cipher.encrypt(black_box(&data))
                });
            },
        );
    }

    group.finish();
}

fn decryption_speed(c: &mut Criterion) {
    let mut group = c.benchmark_group("decryption_speed");
    group.sample_size(1000);
    group.measurement_time(Duration::from_secs(15));

    for size in [64, 512, 4096].iter() {
        group.bench_with_input(
            BenchmarkId::from_parameter(format!("{}_bytes", size)),
            size,
            |b, &size| {
                let cipher = SimpleEncryption::new();
                let data = vec![42u8; size];
                let encrypted = cipher.encrypt(&data);

                b.iter(|| {
                    cipher.decrypt(black_box(&encrypted))
                });
            },
        );
    }

    group.finish();
}

fn encryption_throughput(c: &mut Criterion) {
    let mut group = c.benchmark_group("encryption_throughput");
    group.sample_size(100);

    group.bench_function("bytes_per_second", |b| {
        let cipher = SimpleEncryption::new();
        let data = vec![42u8; 4096];  // 4KB

        let start = std::time::Instant::now();
        let mut bytes = 0u64;

        b.iter(|| {
            let encrypted = cipher.encrypt(&data);
            bytes += encrypted.len() as u64;
            black_box(encrypted)
        });

        let elapsed = start.elapsed();
        let throughput_mbps = (bytes as f64 / (1024 * 1024) as f64) / elapsed.as_secs_f64();
        println!("Encryption Throughput: {:.2} MB/s", throughput_mbps);
    });

    group.finish();
}

fn token_refresh(c: &mut Criterion) {
    let mut group = c.benchmark_group("token_operations");
    group.sample_size(1000);

    group.bench_function("generate_and_validate", |b| {
        let validator = black_box(SimpleTokenValidator::new("secret_key"));
        b.iter(|| {
            let token = validator.generate_token("user@example.com");
            validator.validate_token(&token)
        });
    });

    group.finish();
}

fn security_mixed_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("security_mixed");
    group.sample_size(100);

    group.bench_function("realistic_security_workload", |b| {
        let token_validator = black_box(SimpleTokenValidator::new("secret_key"));
        let cipher = black_box(SimpleEncryption::new());

        b.iter(|| {
            // 50% 토큰 생성
            let _token = token_validator.generate_token("user");

            // 30% 토큰 검증
            let _valid = token_validator.validate_token("user.secret.1234567890");

            // 20% 암호화
            let data = vec![42u8; 256];
            let _encrypted = cipher.encrypt(&data);
        });
    });

    group.finish();
}

criterion_group!(
    benches,
    jwt_token_generation,
    jwt_token_validation,
    encryption_speed,
    decryption_speed,
    encryption_throughput,
    token_refresh,
    security_mixed_operations
);
criterion_main!(benches);
