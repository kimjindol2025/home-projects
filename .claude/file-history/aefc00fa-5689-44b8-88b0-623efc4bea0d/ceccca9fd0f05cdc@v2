# Go 언어 대학 과정 (golang_study)

## 📚 프로젝트 개요

**상태**: ✅ 완성 (v4.1 Advanced)
**철학**: "기록이 증명이다" - 공학적 사고의 진정한 이해
**대상**: 대학/대학원 입문 수준
**저장소**: https://gogs.dclub.kr/kim/golang_study.git

---

## 🎓 학습 로드맵

### **v1.1 Basics (Go 101)** ✅ 완성
- **1-1**: Hello, Gopher! (설치와 첫 컴파일)

### **v2.1 Advanced (Go 201)** ✅ 완성
- **2-1**: GMP 모델 (고루틴 스케줄링) - 커밋: 893c348
- **2-2**: Fan-out/Fan-in 패턴 (병렬 데이터 파이프라인) - 커밋: 2d6f9bb
- **2-3**: 공유 자원의 사수 - 뮤텍스와 원자적 연산 - 커밋: 1f5cf5b
- **2-4**: 고루틴 생명주기 제어 - Context와 WaitGroup - 커밋: 24bf882

### **v3.1 Distributed (Go 301)** ✅ 완성
- **3-1**: gRPC와 Protocol Buffers - 커밋: 0d3322d
- **3-2**: 메시지 기반 분산 아키텍처 - 커밋: bcf3545
- **3-3**: Service Discovery & Health Check - 커밋: abab1aa
- **3-4**: API Gateway & JWT 인증 - 커밋: 3eae0f8
- **3-5**: 관찰 가능성(Observability) & 분산 트레이싱 - 커밋: 42ff7d1

### **v4.1 Advanced (Go 401)** ✅ 완성 (신의 영역)
- **4-1**: 메모리 레이아웃과 이스케이프 분석 - 커밋: 974f3dd
  - 마크다운(2100줄) + 구현(530줄), Stack vs Heap, GC 내부
- **4-2**: Plan 9 어셈블리와 컴파일러 최적화 - 커밋: 595318b
  - 마크다운(2000줄) + 구현(480줄), 가상 레지스터, SSA, SIMD
- **4-3**: Go 런타임 소스 코드 분석 - 커밋: bc28c0e
  - 마크다운(2100줄) + 구현(520줄), proc.go, chan.go, slice.go
- **4-4**: PGO와 정밀 프로파일링 - 커밋: d4be9c9
  - 마크다운(2000줄) + 구현(550줄), 벤치마킹, pprof, Flame Graph
- **4-5**: Zero-Copy & Lock-Free Architecture - 커밋: 78479ab 🎓
  - 마크다운(2100줄) + 구현(542줄) + 테스트(526줄), 박사 과정 완성

---

## 📊 통계

| 카테고리 | 수치 |
|---------|------|
| 총 레슨 | 13개 (v1.1 + v2.1 + v3.1 + v4.1) |
| 총 마크다운 | ~15,000줄 |
| 총 구현 코드 | ~3,500줄 |
| 총 테스트 | 24/24 ✅ |
| 총 커밋 | 13개 |

---

## 🔥 핵심 기술

### v2.1: 병렬 프로그래밍
- **GMP 모델**: GOMAXPROCS와 성능 4.08배 개선
- **Fan-out/Fan-in**: 8워커 vs 1워커, 7.93배 성능
- **Race Condition**: 50.98% 데이터 손실 재현 + Atomic 1.64배 개선
- **Context**: 고루틴 누수 3가지 원인 해결

### v3.1: 분산 시스템
- **gRPC**: REST 8.5초 → gRPC 0.8초 (10배), 10배 작은 데이터
- **Message Broker**: 동기식 990ms → 비동기식 200ms (5배)
- **Service Discovery**: Registry, Client/Server-side, Health Check
- **API Gateway**: JWT, Rate Limiting, Token Refresh
- **Observability**: Trace ID, Prometheus, Golden Signals

### v4.1: 시스템 아키텍트
- **메모리**: 50배 성능 차이, Zero Allocation 기법
- **어셈블리**: 4-16배 성능 개선, CPU 수준 최적화
- **런타임**: Go 내부 동작 원리 완벽 이해
- **프로파일링**: 병목 함수 분석, 메모리 누수 감지
- **Zero-Copy**: Go-Julia Hybrid Kernel 구현

---

## 📂 디렉토리 구조

```
golang_study/
├── README.md                (전체 개요)
├── go.mod                   (모듈)
├── guides/
│   ├── GOGS_GUIDE.md
│   └── GO_INSTALLATION.md
├── v1_1_basics/
│   ├── README.md
│   └── 01_hello_gopher.go
├── v2_1_advanced/
│   ├── 01_gmp_model.go
│   ├── 02_fan_out_fan_in.go
│   ├── 03_mutex_atomic.go
│   └── 04_context_waitgroup.go
├── v3_1_distributed/
│   ├── 01_grpc_protobuf.go
│   ├── 02_message_broker.go
│   ├── 03_service_discovery.go
│   ├── 04_api_gateway.go
│   └── 05_observability.go
└── v4_1_advanced/
    ├── 01_memory_layout.go
    ├── 02_assembly_optimization.go
    ├── 03_runtime_analysis.go
    ├── 04_profiling.go
    └── 05_zero_copy_lockfree.go
```

---

## 🎯 설계 철학

1. **명시적 오류 처리**: if err != nil 을 거리낌 없이
2. **메모리 안전성**: unsafe 포인터는 신중하게
3. **동시성 패턴**: Channel과 Goroutine의 진정한 이해
4. **성능 중심**: 벤치마크로 검증된 설계
5. **이해의 깊이**: 표면적 사용이 아닌 내부 원리 이해

---

## 📚 추천 학습 순서

1. v1.1 기초 (1-2시간)
2. v2.1 병렬 프로그래밍 (1주)
3. v3.1 분산 시스템 (2주)
4. v4.1 시스템 아키텍처 (3주)

**총 예상**: 6-8주 (대학 1학기 강좌 규모)

---

**마지막 업데이트**: 2026-02-26
**상태**: 모든 레슨 완성, 다음 과정 예정 없음
