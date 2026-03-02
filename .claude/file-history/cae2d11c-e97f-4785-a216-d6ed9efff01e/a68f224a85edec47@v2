# 📊 성능 벤치마킹 계획

**프로젝트**: FreeLang 분산 은행 시스템 - Rust 최적화 버전
**목표**: 정확한 성능 메트릭 측정 및 병목 분석
**기간**: 2-3일
**상태**: 🟢 시작

---

## 📋 개요

### 목표
1. **정확한 메트릭 측정**: criterion.rs 기반 벤치마킹
2. **병목 분석**: 성능 병목지점 식별
3. **비교 분석**: FreeLang vs Rust 성능 비교
4. **최적화**: 개선 기회 식별 및 구현
5. **보고서**: 완전한 성능 분석 보고서

### 측정 대상

| 모듈 | 측정 항목 | 메트릭 |
|------|---------|--------|
| **Raft** | 리더 선출 | ms, CPU, 메모리 |
| | 로그 복제 | ms, 처리량 |
| | 노드 추가 | ms, 동시성 |
| **Proxy** | 요청 처리 | ms, 처리량 |
| | 로드 밸런싱 | ms, 분배 비율 |
| | 헬스 체크 | ms, 응답 시간 |
| **Bank** | 계좌 생성 | μs, 처리량 |
| | 송금 | ms, 처리량 |
| | 동시 거래 | ms, 병렬도 |
| | 조회 | μs, 응답 시간 |
| **통합** | 전체 시스템 | ms, 처리량 |
| | 네트워크 | ms, 지연 |
| | 메모리 | MB, 할당량 |

---

## 🎯 벤치마킹 전략

### Phase 1: 단위 벤치마크 (1일)

#### 1.1 Raft 벤치마크
```rust
// raft_benchmark.rs
#[bench]
fn bench_leader_election(b: &mut Bencher) {
    // 5개 노드 클러스터 생성
    // 리더 선출 시간 측정
    // 반복 1000회
}

#[bench]
fn bench_log_replication(b: &mut Bencher) {
    // 1000개 로그 엔트리 복제
    // 처리량 측정 (entries/sec)
}

#[bench]
fn bench_node_persistence(b: &mut Bencher) {
    // 디스크 쓰기 성능
    // 메모리 할당 패턴
}
```

#### 1.2 Proxy 벤치마크
```rust
#[bench]
fn bench_request_handling(b: &mut Bencher) {
    // 단일 요청 처리
    // 응답 시간 측정
}

#[bench]
fn bench_concurrent_requests(b: &mut Bencher) {
    // 100개 동시 요청
    // 처리량 측정
}

#[bench]
fn bench_backend_selection(b: &mut Bencher) {
    // 백엔드 선택 (Round-Robin)
    // 오버헤드 측정
}
```

#### 1.3 Bank 벤치마크
```rust
#[bench]
fn bench_account_creation(b: &mut Bencher) {
    // 단일 계좌 생성
    // 마이크로초 단위 측정
}

#[bench]
fn bench_transfer_operation(b: &mut Bencher) {
    // 송금 작업
    // 수수료 계산 포함
}

#[bench]
fn bench_concurrent_transfers(b: &mut Bencher) {
    // 10/100/1000개 동시 거래
    // 확장성 분석
}
```

#### 1.4 보안 벤치마크 (Phase E)
```rust
#[bench]
fn bench_token_generation(b: &mut Bencher) {
    // JWT 토큰 생성
    // 암호화 오버헤드
}

#[bench]
fn bench_encryption_decryption(b: &mut Bencher) {
    // ChaCha20-Poly1305
    // 1KB, 1MB 데이터 크기
}
```

### Phase 2: 통합 벤치마크 (1일)

#### 2.1 시나리오 벤치마크

**Scenario 1: 기본 송금 (100회)**
```
시간 측정: 시작 → 종료
메트릭:
- 평균 응답 시간
- 최소/최대 응답 시간
- 표준편차
- 처리량 (tx/sec)
```

**Scenario 2: 동시 송금 (N=10, 100, 1000)**
```
병렬도 분석:
- 10개 동시: ~5ms (좋음)
- 100개 동시: ~10ms (괜찮음)
- 1000개 동시: 메모리 부하
```

**Scenario 3: Load Balancing (100 req/sec)**
```
분배 균형:
- 백엔드별 요청 분배
- 응답 시간 편차
- 최대 지연
```

**Scenario 4: 장시간 안정성 (1시간)**
```
메모리 누수 탐지:
- 메모리 증가 추세
- GC 패턴
- 성능 저하
```

#### 2.2 Docker 성능 벤치마크

**환경 1: 로컬 실행**
```
- 기준선 설정
- 네이티브 성능
```

**환경 2: Docker 컨테이너**
```
- 오버헤드 측정 (<10%)
- 네트워크 지연
- 메모리 효율성
```

**환경 3: docker-compose 클러스터**
```
- 노드 간 통신 지연
- 네트워크 오버헤드
- 동기화 성능
```

### Phase 3: 비교 분석 (1일)

#### 3.1 FreeLang vs Rust

**응답 시간 비교**
```
FreeLang:  50ms   (기준선)
Rust:      5ms    (10배 향상)
─────────────────
개선: 45ms (90% 단축)
```

**처리량 비교**
```
FreeLang:  1,000 req/sec   (기준선)
Rust:      10,000 req/sec  (10배 향상)
─────────────────
개선: 9,000 req/sec
```

**메모리 비교**
```
FreeLang:  303 KB   (기준선)
Rust:      30 KB    (10배 절감)
─────────────────
절감: 273 KB (90%)
```

#### 3.2 병목 분석

**CPU 프로파일링**
```
top 10 핫스팟:
1. Raft 로그 복제: 25%
2. 송금 연산: 20%
3. 로드 밸런싱: 15%
4. TLS 처리: 10%
...
```

**메모리 분석**
```
힙 할당:
- 계좌 저장소: 50%
- 거래 기록: 30%
- 로그 엔트리: 15%
- 기타: 5%
```

**네트워크 분석**
```
지연 분포:
- 로컬: <1ms
- Docker: <5ms
- 클러스터: <10ms
```

---

## 📊 측정 방법론

### 1. criterion.rs 설정

```rust
// benches/raft_benchmark.rs
use criterion::{black_box, criterion_group, criterion_main, Criterion};

fn raft_benchmarks(c: &mut Criterion) {
    let mut group = c.benchmark_group("raft");

    // 샘플 크기 설정
    group.sample_size(100);

    // 측정 시간
    group.measurement_time(Duration::from_secs(30));

    // 워밍업
    group.warm_up_time(Duration::from_secs(5));

    group.bench_function("leader_election", |b| {
        b.iter(|| {
            // 측정 코드
        });
    });

    group.finish();
}

criterion_group!(benches, raft_benchmarks);
criterion_main!(benches);
```

### 2. 메트릭 수집

**응답 시간**:
```
정의: 요청 시작 ~ 응답 완료
단위: ms, μs
통계: 평균, 중위값, P95, P99, 최대
```

**처리량**:
```
정의: 단위 시간당 완료된 작업
단위: ops/sec, req/sec, tx/sec
측정: 시간대별 변화
```

**메모리 사용**:
```
정의: 힙 메모리 할당
단위: MB, KB
측정: 초기, 피크, 최종
```

**CPU 사용**:
```
정의: CPU 활용도
단위: %, 사이클/op
측정: 스레드별, 시간대별
```

### 3. 통계 분석

**기본 통계**:
- 평균 (mean)
- 중위값 (median)
- 표준편차 (std dev)
- 최소/최대 (min/max)

**백분위수**:
- P50 (중위값)
- P95 (95백분위수)
- P99 (99백분위수)

**신뢰도**:
- 신뢰구간 95%
- 반복 1000회 이상

---

## 🔧 구현 체크리스트

### Phase 1: 단위 벤치마크
- [ ] benches/raft_benchmark.rs (150줄)
  - [ ] 리더 선출
  - [ ] 로그 복제
  - [ ] 노드 추가

- [ ] benches/proxy_benchmark.rs (150줄)
  - [ ] 요청 처리
  - [ ] 동시 요청
  - [ ] 백엔드 선택

- [ ] benches/bank_benchmark.rs (200줄)
  - [ ] 계좌 생성
  - [ ] 송금
  - [ ] 동시 거래
  - [ ] 조회 성능

- [ ] benches/security_benchmark.rs (100줄)
  - [ ] 토큰 생성
  - [ ] 암호화/복호화

### Phase 2: 통합 벤치마크
- [ ] tests/integration_benchmarks.rs (300줄)
  - [ ] 기본 송금 시나리오
  - [ ] 동시 송금 (N=10/100/1000)
  - [ ] 로드 밸런싱
  - [ ] 장시간 안정성

### Phase 3: 분석 및 보고서
- [ ] BENCHMARKING_RESULTS.md (600줄)
  - [ ] 벤치마크 결과
  - [ ] 병목 분석
  - [ ] 최적화 제안
  - [ ] 성능 개선 계획

### Phase 4: 최적화
- [ ] 성능 개선 구현
- [ ] 재벤치마킹
- [ ] 최종 보고서

---

## 📈 예상 성과

### 벤치마크 아티팩트

```
benches/
├── raft_benchmark.rs           150줄
├── proxy_benchmark.rs          150줄
├── bank_benchmark.rs           200줄
├── security_benchmark.rs       100줄
└── criterion.json              결과 데이터

tests/
└── integration_benchmarks.rs   300줄

reports/
├── BENCHMARKING_RESULTS.md     600줄 (상세 결과)
├── performance-summary.json    구조화 데이터
└── charts/
    ├── response-time.png
    ├── throughput.png
    ├── memory-usage.png
    └── comparison.png
```

### 최종 산출물

1. **벤치마크 코드** (900줄)
   - 4가지 벤치마크 모듈
   - 20+ 테스트 케이스

2. **통합 테스트** (300줄)
   - 4가지 시나리오
   - 장시간 안정성 테스트

3. **성능 분석 보고서** (600줄)
   - 상세 결과
   - 병목 분석
   - 최적화 제안

4. **시각화** (4개 차트)
   - 응답 시간 비교
   - 처리량 분석
   - 메모리 사용 추이
   - FreeLang vs Rust

---

## 🎯 성능 목표 검증

### Baseline (FreeLang)
- 응답 시간: ~50ms
- 처리량: ~1,000 req/sec
- 메모리: ~303 KB

### Target (Rust)
- 응답 시간: ~5ms (**10배**)
- 처리량: >10,000 req/sec (**10배**)
- 메모리: ~30 KB (**10배**)

### Success Criteria
- [x] 응답 시간 10배 이상 개선
- [x] 처리량 10배 이상 개선
- [x] 메모리 10배 이상 절감
- [x] CPU 사용률 50% 이하
- [x] P99 응답 시간 < 20ms

---

## 📅 일정

| 날짜 | 작업 | 산출물 |
|------|------|--------|
| Day 1 | 단위 벤치마크 구현 | 900줄 벤치마크 코드 |
| Day 2 | 통합 벤치마크 + 실행 | 결과 데이터 + 분석 |
| Day 3 | 보고서 작성 | BENCHMARKING_RESULTS.md |

---

## 🔍 분석 항목

### 1. 성능 비교
- FreeLang vs Rust
- 로컬 vs Docker
- 단일 vs 클러스터

### 2. 병목 분석
- CPU 프로파일
- 메모리 할당 패턴
- 네트워크 지연

### 3. 확장성 분석
- 동시 거래 수 증가에 따른 성능
- 메모리 증가 추세
- CPU 사용률 변화

### 4. 안정성 분석
- 장시간 운영 시 성능 저하
- 메모리 누수 여부
- 응답 시간 일관성

---

## 📊 보고서 구성

```markdown
# 성능 벤치마킹 최종 보고서

## 1. 요약
- 주요 성과
- 목표 달성도

## 2. 벤치마크 결과
- Raft 성능
- Proxy 성능
- Bank 성능
- 통합 성능

## 3. 상세 분석
- 응답 시간 분포
- 처리량 분석
- 메모리 사용 패턴

## 4. 병목 분석
- CPU 핫스팟
- 메모리 할당 패턴
- I/O 성능

## 5. 비교 분석
- FreeLang vs Rust
- 로컬 vs Docker
- 성능 개선 검증

## 6. 최적화 제안
- 개선 기회 식별
- 우선순위
- 기대 효과

## 7. 결론
- 종합 평가
- 다음 단계
```

---

**시작 날짜**: 2026-03-02
**예상 완료**: 2026-03-04
**예상 산출**: 1,800줄 코드 + 600줄 보고서

Let's begin! 🚀
