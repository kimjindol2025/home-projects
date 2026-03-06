# 🔐 FreeLang Backend Production - 상세 설계

**프로젝트**: freelang-backend-production
**버전**: 1.0
**상태**: 🚀 구현 중 (Phase C)
**목표**: Production Hardening 8개 모듈 (2,600줄)

---

## 🎯 개요

**목적**: HTTP Engine + REST API를 운영 환경에 배포 가능하게 강화

**핵심 8개 모듈**:
1. Structured Logging (400줄) - JSON 로그 포맷
2. Distributed Tracing (500줄) - 분산 추적
3. Circuit Breaker (400줄) - 장애 자동 차단
4. Rate Limiting (300줄) - 요청 제한
5. Health Checks (250줄) - 상태 확인
6. Graceful Shutdown (200줄) - 우아한 종료
7. Metrics Export (300줄) - Prometheus 포맷
8. Configuration (250줄) - 설정 관리

---

## 📂 파일 구조

```
freelang-backend-production/
├── src/
│   ├── logger.fl              (400줄) - Structured Logging
│   ├── tracer.fl              (500줄) - Distributed Tracing
│   ├── circuit_breaker.fl     (400줄) - Circuit Breaker
│   ├── rate_limiter.fl        (300줄) - Rate Limiting
│   ├── health_checker.fl      (250줄) - Health Checks
│   ├── shutdown_handler.fl    (200줄) - Graceful Shutdown
│   ├── metrics_exporter.fl    (300줄) - Metrics Export
│   ├── config_manager.fl      (250줄) - Configuration
│   └── mod.fl                 (50줄)  - 공개 API
├── tests/
│   ├── test_logger.fl         (100줄)
│   ├── test_tracer.fl         (100줄)
│   ├── test_circuit_breaker.fl (150줄)
│   ├── test_rate_limiter.fl   (100줄)
│   ├── test_health_checker.fl (80줄)
│   ├── test_shutdown.fl       (80줄)
│   ├── test_metrics.fl        (80줄)
│   ├── test_config.fl         (80줄)
│   └── integration_test.fl    (150줄)
└── README.md
```

---

## 📐 Module 1: Structured Logging (400줄)

### 목적
JSON 형식의 구조화된 로그

### 구조
```freelang
struct LogEntry {
  timestamp: string,        // ISO 8601
  level: string,           // DEBUG, INFO, WARN, ERROR
  message: string,
  context: map<string, any>,
  trace_id: string,        // 분산 추적
  error: string,           // 에러 메시지
}

enum LogLevel {
  Debug,    // 개발용
  Info,     // 정보
  Warn,     // 경고
  Error,    // 에러
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `logger_init(level)` | 로거 초기화 |
| `logger_debug(msg)` | DEBUG 로그 |
| `logger_info(msg)` | INFO 로그 |
| `logger_warn(msg)` | WARN 로그 |
| `logger_error(msg)` | ERROR 로그 |
| `logger_with_context(entry, context)` | 컨텍스트 추가 |
| `logger_output(entry)` | 로그 출력 (파일/stdout) |

### 무관용 규칙

**R1**: 로그 처리 < 1ms
- 측정: 평균 로그 처리 시간
- 포함: JSON 직렬화, 쓰기

---

## 📐 Module 2: Distributed Tracing (500줄)

### 목적
분산 시스템 요청 추적

### 구조
```freelang
struct Trace {
  trace_id: string,         // 전역 유니크
  parent_span_id: string,   // 부모 span
  span_id: string,          // 현재 span
  operation: string,        // 작업명
  start_time: i64,          // 나노초
  end_time: i64,
  status: string,           // ok, error, timeout
  tags: map<string, string>,
  logs: array<string>,
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `tracer_start(operation)` | Span 시작 |
| `tracer_end(span)` | Span 종료 |
| `tracer_with_tag(span, key, value)` | 태그 추가 |
| `tracer_propagate(headers)` | 헤더로 전파 |
| `tracer_extract(headers)` | 헤더에서 추출 |
| `tracer_export(span)` | Jaeger 포맷 내보내기 |

### 무관용 규칙

**R2**: Trace 오버헤드 < 5%
- 측정: 전체 요청 시간 대비 trace 시간
- 목표: 성능 영향 최소화

---

## 📐 Module 3: Circuit Breaker (400줄)

### 목적
장애 자동 차단 및 복구

### 상태 기계

```
        ┌─────────────┐
        │   CLOSED    │ (정상)
        └─────────────┘
               │
        실패 N번 연속
               │
        ┌─────────────┐
        │     OPEN    │ (차단)
        └─────────────┘
               │
        timeout(30초)
               │
        ┌─────────────┐
        │  HALF-OPEN  │ (시험)
        └─────────────┘
               │
        테스트 요청
        ├─ 성공 → CLOSED
        └─ 실패 → OPEN
```

### 구조
```freelang
struct CircuitBreaker {
  state: string,            // CLOSED, OPEN, HALF_OPEN
  failure_count: i32,       // 연속 실패 수
  success_count: i32,
  failure_threshold: i32,   // N개 실패 시 차단
  last_failure_time: i64,
  timeout_ms: i32,          // 복구 시도 대기
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `circuit_breaker_new(threshold, timeout)` | CB 생성 |
| `circuit_breaker_call(cb, fn)` | 함수 호출 (CB 적용) |
| `circuit_breaker_record_success(cb)` | 성공 기록 |
| `circuit_breaker_record_failure(cb)` | 실패 기록 |
| `circuit_breaker_reset(cb)` | 리셋 |

### 무관용 규칙

**R3**: Circuit breaker < 100µs
- 측정: 상태 확인 + 전환 시간
- 목표: 오버헤드 무시할 수준

---

## 📐 Module 4: Rate Limiting (300줄)

### 목적
요청 속도 제한

### 알고리즘: Token Bucket

```
┌──────────────────┐
│  Token Bucket    │
│  capacity: 1000  │
│  tokens: 850     │
│  rate: 100/sec   │
└──────────────────┘
        │
   요청 → consume 1 token
   ├─ 토큰 있음 → 통과
   └─ 토큰 없음 → 거부 (429 Too Many Requests)
```

### 구조
```freelang
struct RateLimiter {
  bucket_capacity: i32,     // 최대 토큰
  tokens: i32,              // 현재 토큰
  refill_rate: i32,         // token/sec
  last_refill: i64,         // 마지막 충전 시간
  per_user: bool,           // 사용자별 제한
  per_ip: bool,             // IP별 제한
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `rate_limiter_new(capacity, rate)` | 제한기 생성 |
| `rate_limiter_allow(limiter, id)` | 요청 허용 여부 |
| `rate_limiter_reset(limiter)` | 리셋 |

### 무관용 규칙

**R4**: Rate limiting accuracy ≥ 99%
- 측정: 제한 정확도
- 목표: 정확한 rate limiting

---

## 📐 Module 5: Health Checks (250줄)

### 목적
서비스 상태 모니터링

### 2가지 Probe

```
Liveness Probe:   서비스 살아있나? (주기: 10초)
├─ 응답 있음 → OK
└─ 타임아웃 → DEAD

Readiness Probe:  트래픽 받을 준비됐나? (주기: 5초)
├─ DB 연결 OK, 캐시 OK → READY
└─ 뭔가 안 됨 → NOT_READY
```

### 구조
```freelang
struct HealthStatus {
  service: string,          // 서비스명
  status: string,           // HEALTHY, UNHEALTHY
  uptime: i64,              // 초 단위
  db_connected: bool,
  cache_connected: bool,
  memory_mb: i32,
  last_check: i64,
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `health_check_init()` | 헬스 체크 초기화 |
| `health_check_liveness()` | Liveness 체크 |
| `health_check_readiness()` | Readiness 체크 |
| `health_check_status()` | 전체 상태 조회 |

### 무관용 규칙

**R5**: Health check < 500ms
- 측정: 전체 체크 시간
- 목표: 빠른 응답

---

## 📐 Module 6: Graceful Shutdown (200줄)

### 목적
안전한 종료 및 연결 정리

### 프로세스

```
1. SIGTERM 신호 수신
2. "드레인" 시작: 새 연결 거부
3. 기존 요청 완료 대기 (최대 30초)
4. 타임아웃 또는 모든 요청 완료
5. 연결 종료
6. 종료
```

### 구조
```freelang
struct ShutdownManager {
  is_shutting_down: bool,
  active_connections: i32,
  drain_timeout_ms: i32,
  drain_start_time: i64,
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `shutdown_register_signal()` | SIGTERM 등록 |
| `shutdown_initiate()` | 종료 시작 |
| `shutdown_wait_connections()` | 연결 대기 |
| `shutdown_complete()` | 완료 |

### 무관용 규칙

**R6**: Graceful shutdown < 30s
- 측정: 신호 수신 → 완료
- 목표: 30초 내 종료

---

## 📐 Module 7: Metrics Export (300줄)

### 목적
Prometheus 호환 메트릭 내보내기

### 포맷

```
# HELP http_requests_total 총 요청 수
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/todos"} 1523

# HELP http_request_duration_seconds 요청 지연시간
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 1000
http_request_duration_seconds_bucket{le="0.5"} 1450
http_request_duration_seconds_bucket{le="1.0"} 1500

# HELP process_resident_memory_bytes 메모리 사용
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 52428800
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `metrics_counter(name, labels)` | Counter 생성 |
| `metrics_histogram(name, labels, buckets)` | Histogram 생성 |
| `metrics_gauge(name, value)` | Gauge 업데이트 |
| `metrics_export()` | Prometheus 포맷 내보내기 |

### 무관용 규칙

**R7**: Metrics 수집 < 10MB
- 측정: 메트릭 메모리 사용
- 목표: 저 메모리 오버헤드

---

## 📐 Module 8: Configuration (250줄)

### 목적
환경별 설정 관리

### 소스 우선순위

```
1. 환경 변수 (최고)
   BACKEND_PORT=8080
   BACKEND_LOG_LEVEL=DEBUG

2. 설정 파일
   config.toml 또는 config.json

3. 기본값 (최저)
   port: 8080
   log_level: "INFO"
```

### 구조
```freelang
struct Config {
  server_addr: string,
  server_port: i32,
  log_level: string,
  db_host: string,
  db_port: i32,
  cache_enabled: bool,
  metrics_enabled: bool,
  tracing_enabled: bool,
}
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `config_load_env()` | 환경 변수 로드 |
| `config_load_file(path)` | 파일 로드 |
| `config_get(key)` | 설정 값 조회 |
| `config_validate()` | 검증 |

### 무관용 규칙

**R8**: Config 재로드 < 100ms
- 측정: 재로드 시간
- 목표: 빠른 재설정

---

## 🧪 테스트 전략

**30개 테스트** (8개 모듈 × 3-4개 + 통합 2개)

| 모듈 | 테스트 | 상태 |
|------|--------|------|
| Logger | 4개 (정상, 에러, JSON, 성능) | ⏳ |
| Tracer | 4개 (span 생성, 전파, 추출, 내보내기) | ⏳ |
| Circuit Breaker | 6개 (전환, 복구, 타임아웃, 성공) | ⏳ |
| Rate Limiter | 4개 (허용, 거부, 리셋, 정확도) | ⏳ |
| Health Checker | 3개 (Liveness, Readiness, 상태) | ⏳ |
| Shutdown | 3개 (신호, 드레인, 완료) | ⏳ |
| Metrics | 3개 (Counter, Histogram, Gauge) | ⏳ |
| Config | 3개 (환경, 파일, 검증) | ⏳ |
| **Integration** | **2개** (전체 연동) | ⏳ |
| **합계** | **30개** | ⏳ |

---

## 📈 무관용 규칙 (10개)

| # | 규칙 | 목표 | 측정 |
|---|------|------|------|
| R1 | 로그 처리 | < 1ms | 평균 시간 |
| R2 | Trace 오버헤드 | < 5% | 성능 영향 |
| R3 | Circuit breaker | < 100µs | 상태 확인 |
| R4 | Rate limiting accuracy | ≥ 99% | 정확도 |
| R5 | Health check | < 500ms | 응답 시간 |
| R6 | Graceful shutdown | < 30s | 종료 시간 |
| R7 | Metrics memory | < 10MB | RSS |
| R8 | Config reload | < 100ms | 재로드 시간 |
| R9 | 에러 추적 | 100% | 모든 에러 기록 |
| R10 | 의존성 장애 검출 | 100% | DB/Cache 장애 감지 |

---

## 📊 구현 체크리스트

### Structured Logging (400줄)

- [ ] LogEntry 구조체
- [ ] LogLevel enum
- [ ] logger_init()
- [ ] 5가지 로그 함수 (debug/info/warn/error)
- [ ] JSON 직렬화
- [ ] 파일/stdout 출력
- [ ] 타임스탬프 (ISO 8601)
- [ ] 컨텍스트 추가

### Distributed Tracing (500줄)

- [ ] Trace/Span 구조체
- [ ] Trace ID 생성 (UUID)
- [ ] Span 생성/종료
- [ ] 태그/로그 추가
- [ ] 헤더 전파 (W3C TraceContext)
- [ ] 헤더 추출
- [ ] Jaeger 포맷 내보내기
- [ ] 성능: < 5% 오버헤드

### Circuit Breaker (400줄)

- [ ] 상태 기계 (CLOSED/OPEN/HALF_OPEN)
- [ ] 실패 카운트
- [ ] 상태 전환 로직
- [ ] 타임아웃 복구
- [ ] call() 래퍼 함수
- [ ] 성공/실패 기록
- [ ] 리셋

### Rate Limiter (300줄)

- [ ] Token Bucket 알고리즘
- [ ] 토큰 충전 (refill)
- [ ] allow() 판단
- [ ] 사용자별 제한 (per-user)
- [ ] IP별 제한 (per-IP)
- [ ] 정확도 ≥ 99%

### Health Checker (250줄)

- [ ] HealthStatus 구조체
- [ ] Liveness probe (10초)
- [ ] Readiness probe (5초)
- [ ] DB 연결 확인
- [ ] 캐시 연결 확인
- [ ] 메모리 사용량
- [ ] 응답 < 500ms

### Graceful Shutdown (200줄)

- [ ] SIGTERM 등록
- [ ] 드레인 시작 (새 연결 거부)
- [ ] 기존 연결 대기
- [ ] 타임아웃 (30초)
- [ ] 깔끔한 종료

### Metrics Export (300줄)

- [ ] Counter (누적)
- [ ] Histogram (분포)
- [ ] Gauge (현재값)
- [ ] Prometheus 포맷
- [ ] 메모리 < 10MB

### Configuration (250줄)

- [ ] 환경 변수 로드
- [ ] 파일 로드 (JSON/TOML)
- [ ] 기본값
- [ ] 검증
- [ ] 재로드 < 100ms

---

## 🎯 완성 조건

✅ 모든 체크리스트 항목 완성
✅ 모든 테스트 통과 (30개)
✅ 모든 무관용 규칙 준수 (R1-R10)
✅ GOGS에 커밋

---

**다음**: Structured Logging 구현 시작 (src/logger.fl)
