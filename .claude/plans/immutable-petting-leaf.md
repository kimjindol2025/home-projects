# GRIE Phase 3: 운영 준비 구현 계획

## Context
Phase 1-4 완료 및 Phase 2 신뢰성 검증(260M ops, 0% 에러, 1시간 장기 테스트 통과) 이후,
GRIE를 프로덕션에 배포하기 위한 운영 기반 강화 작업이다.
하드코딩된 timeout, 부재한 Circuit Breaker, 불완전한 Graceful Shutdown, 모니터링 부재
문제를 표준 라이브러리만으로 해결한다.

**제약사항**: 외부 의존성 추가 없음 (표준 라이브러리만), 기존 API 유지 확장

---

## 현황 분석

| 항목 | 현재 상태 | 문제점 |
|------|-----------|--------|
| Timeout | 부분 구현 | `worker.go:processJob()` 5초 하드코딩 |
| Circuit Breaker | 없음 | 전혀 미구현 |
| 에러 복구/재시도 | 없음 | 에러 분류, 재시도 로직 없음 |
| Graceful Shutdown | 일부 완성 | `cmd/writer`: `select{}` 무한 블록 |
| 모니터링 | GetStats()만 | HTTP 엔드포인트 없음 |

---

## 구현 단계 (8단계)

### 단계 1: 기반 패키지 (병렬 구현)

**신규: `internal/config/timeout.go`**
```go
type TimeoutConfig struct {
    JobProcessTimeout time.Duration // 5초 하드코딩 대체
    SHMWriteTimeout   time.Duration // 1초
    SHMReadTimeout    time.Duration // 100ms
    KernelWaitTimeout time.Duration // 5초
    ShutdownTimeout   time.Duration // 30초
}
func DefaultTimeoutConfig() *TimeoutConfig
func FastTimeoutConfig() *TimeoutConfig  // 테스트용 빠른 값
```

**신규: `internal/retry/backoff.go`**
```go
type ErrorType int  // ErrorTransient | ErrorPermanent
type RetryableError struct{ Err error; Type ErrorType }
type RetryConfig struct {
    MaxAttempts int
    BaseDelay, MaxDelay time.Duration
    Multiplier, JitterFactor float64
}
func DefaultRetryConfig() RetryConfig  // 3회, 10ms 기반, 2배 배율
func Do(ctx context.Context, cfg RetryConfig, fn func() error) error
```

### 단계 2: Circuit Breaker

**신규: `internal/circuitbreaker/breaker.go`**
```go
type State int32  // StateClosed=0 | StateOpen=1 | StateHalfOpen=2
type BreakerConfig struct {
    FailureThreshold int           // 연속 실패 5회 → Open
    RecoveryTimeout  time.Duration // 30초 후 Half-Open
    SuccessThreshold int           // 연속 성공 2회 → Closed
    HalfOpenMaxCalls int           // Half-Open 동시 허용 1개
}
type CircuitBreaker struct { /* 모든 필드 atomic */ }
func (cb *CircuitBreaker) Allow() error         // 차단 시 ErrCircuitOpen
func (cb *CircuitBreaker) RecordSuccess()
func (cb *CircuitBreaker) RecordFailure()
func (cb *CircuitBreaker) CurrentState() State
func (cb *CircuitBreaker) Stats() map[string]interface{}
```

상태 전이:
- Closed → Open: `consecutiveFail >= FailureThreshold` 시
- Open → Half-Open: `time.Since(openedAt) >= RecoveryTimeout` 시
- Half-Open → Closed: `consecutiveSucc >= SuccessThreshold` 시
- Half-Open → Open: Half-Open 중 실패 즉시

### 단계 3: worker.go 수정

**파일**: `internal/dispatcher/worker.go`

```go
// Worker 구조체에 3개 필드 추가
type Worker struct {
    // ... 기존 ...
    timeoutCfg     *config.TimeoutConfig
    circuitBreaker *circuitbreaker.CircuitBreaker
    retryConfig    retry.RetryConfig
}

// processJob 개선
func (w *Worker) processJob(ctx context.Context, job *Job) ([]byte, error) {
    if err := w.circuitBreaker.Allow(); err != nil {  // CB 체크
        return nil, &retry.RetryableError{Err: err, Type: retry.ErrorPermanent}
    }
    err := retry.Do(ctx, w.retryConfig, func() error {
        deadline := time.Now().Add(w.timeoutCfg.JobProcessTimeout) // 하드코딩 제거
        // ... 기존 로직 ...
    })
    if err != nil { w.circuitBreaker.RecordFailure(); return nil, err }
    w.circuitBreaker.RecordSuccess()
    return result, nil
}
```

### 단계 4: dispatcher.go 수정

**파일**: `internal/dispatcher/dispatcher.go`

```go
// wg 추가, assignJobToWorker에 wg.Add(1)/Done() 추가
// Close()에 wg.Wait() + ShutdownTimeout 추가
func (d *Dispatcher) Close() error {
    for _, w := range d.workers { w.SetState("SHUTDOWN") }
    done := make(chan struct{})
    go func() { d.wg.Wait(); close(done) }()
    select {
    case <-done: return nil
    case <-time.After(30 * time.Second):
        return fmt.Errorf("shutdown timeout: jobs still in-flight")
    }
}
```

### 단계 5: Graceful Shutdown (cmd 수정)

**파일**: `cmd/writer/main.go`
- `select{}` 제거 → `signal.Notify` + `<-sigChan` 대기
- SIGINT/SIGTERM 수신 시 mgr.Close() 실행 보장

**파일**: `cmd/reader/main.go`
- 동일 패턴 적용

### 단계 6: 모니터링 패키지 (독립, 병렬 가능)

**신규: `internal/monitoring/metrics.go`**
```go
type Counter struct{ name, help string; value uint64 }  // atomic
type Gauge   struct{ name, help string; value int64  }  // atomic
type Registry struct {
    JobsReceived, JobsDropped, JobsProcessed, JobsFailed *Counter
    QueueSize, ActiveWorkers *Gauge
    CBAllowed, CBRejected *Counter
    SHMWrites, SHMReads, SHMErrors *Counter
    RetryAttempts, RetrySuccess, RetryExhausted *Counter
    startTime time.Time
}
var Default = NewRegistry()
func (r *Registry) WritePrometheusFormat(w io.Writer)  // Prometheus text format
```

**신규: `internal/monitoring/server.go`**
```go
type ServerConfig struct {
    Addr            string        // ":9090"
    ReadTimeout     time.Duration // 5초
    WriteTimeout    time.Duration // 10초
    ShutdownTimeout time.Duration // 5초
}
type Server struct{ cfg, registry, srv }
func (s *Server) Start() error          // 백그라운드 goroutine
func (s *Server) Stop(ctx) error        // graceful shutdown
```

엔드포인트:
| 경로 | 응답 | 용도 |
|------|------|------|
| `GET /metrics` | Prometheus text | 전체 메트릭 |
| `GET /health` | JSON 200 | Liveness probe |
| `GET /ready` | 200/503 | Readiness probe |
| `GET /debug/pprof/` | pprof UI | 런타임 프로파일링 |

### 단계 7: 헬스체크 CLI

**신규: `cmd/health_check/main.go`**
```
사용법: health_check -addr localhost:9090 [-check health|ready] [-json]
Exit code: 0=정상, 1=비정상, 2=연결 오류
```

### 단계 8: 운영 문서

**신규: `OPERATIONS_GUIDE.md`**
- 시작/종료 절차 (SOP)
- 헬스체크 방법, 메트릭 해석
- Circuit Breaker 상태 모니터링 및 복구 절차
- 트러블슈팅 가이드

---

## 파일 목록

### 신규 파일 (10개)
```
internal/config/timeout.go
internal/retry/backoff.go
internal/retry/backoff_test.go
internal/circuitbreaker/breaker.go
internal/circuitbreaker/breaker_test.go
internal/monitoring/metrics.go
internal/monitoring/server.go
internal/monitoring/metrics_test.go
cmd/health_check/main.go
OPERATIONS_GUIDE.md
```

### 수정 파일 (4개)
```
internal/dispatcher/worker.go       (timeout 하드코딩 제거, CB/retry 주입)
internal/dispatcher/dispatcher.go   (wg.Wait 추가로 Graceful Shutdown 완성)
cmd/writer/main.go                  (select{} → signal 핸들링)
cmd/reader/main.go                  (signal 핸들링 추가)
```

---

## 구현 의존성 순서

```
단계 1 (병렬): config/timeout.go  ←→  retry/backoff.go
단계 2:        circuitbreaker/breaker.go
단계 3:        dispatcher/worker.go       (1,2 완료 후)
단계 4:        dispatcher/dispatcher.go   (3 완료 후)
단계 5 (병렬): cmd/writer, cmd/reader     (4 완료 후)
단계 6 (독립): monitoring/               (언제든 가능)
단계 7:        cmd/health_check/          (6 완료 후)
단계 8:        OPERATIONS_GUIDE.md        (전체 완료 후)
```

---

## 검증 계획

```bash
# 전체 테스트 (race detector)
go test -race ./... -v

# Circuit Breaker 단독
go test -race ./internal/circuitbreaker/... -v

# Graceful Shutdown 수동 검증
go run cmd/writer/main.go &
sleep 2 && kill -SIGTERM $!
# 기대: "Received signal: terminated. Shutting down..." 출력

# 모니터링 서버 수동 검증
go run cmd/writer/main.go &
curl localhost:9090/health
curl localhost:9090/ready
curl localhost:9090/metrics

# 헬스체크 CLI
./bin/health_check -addr localhost:9090 -check ready -json
```

---

## 우선순위

**필수 (Phase 3 핵심)**:
1. Circuit Breaker
2. Graceful Shutdown 완성
3. 모니터링 서버

**선택 (품질 개선)**:
4. Timeout 표준화
5. 재시도 로직
6. 헬스체크 CLI
7. 운영 문서
