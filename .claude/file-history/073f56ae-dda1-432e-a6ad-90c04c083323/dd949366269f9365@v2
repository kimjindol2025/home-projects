# GRIE Phase 3: 운영 가이드

GRIE (Raft Consensus Engine)의 프로덕션 운영을 위한 완벽한 가이드입니다.

## 📋 목차

1. [시작 및 종료 절차](#시작-및-종료-절차)
2. [헬스체크 및 모니터링](#헬스체크-및-모니터링)
3. [메트릭 해석](#메트릭-해석)
4. [Circuit Breaker 이해](#circuit-breaker-이해)
5. [트러블슈팅](#트러블슈팅)

---

## 시작 및 종료 절차

### 시작하기

#### 1. Writer 프로세스 시작 (공유 메모리 및 작업 관리)

```bash
go run cmd/writer/main.go -size 10485760 -count 5 -delay 100ms
```

**옵션:**
- `-size`: 공유 메모리 크기 (기본값: 10MB)
- `-count`: 초기 메시지 수 (기본값: 5)
- `-delay`: 메시지 간 간격 (기본값: 100ms)
- `-remove`: 종료 시 공유 메모리 파일 제거 (기본값: false)

**예상 출력:**
```
📝 GRIE Writer Started
📂 Shared Memory: /tmp/grie_shm_1234567890
💾 Size: 10485760 bytes
📤 Messages: 5

---

[1/5] Writing: "Message #1 - Timestamp: 2024-01-01T12:00:00.123456789Z"
  ✅ State: ready | SeqNum: 1 | DataLen: 57 | Writes: 1 | Reads: 0
  ⏳ Waiting 100ms...

...

⏱️  Waiting for readers...
📌 Path: /tmp/grie_shm_1234567890
💡 To run reader: go run cmd/reader/main.go /tmp/grie_shm_1234567890

Press Ctrl+C to exit and cleanup...
```

#### 2. Reader 프로세스 시작 (옵션, 데이터 검증용)

다른 터미널에서:

```bash
go run cmd/reader/main.go /tmp/grie_shm_1234567890 -timeout 5000 -count 0
```

**옵션:**
- `-timeout`: 읽기 타임아웃 (밀리초, 기본값: 5000)
- `-count`: 읽을 메시지 수 (0 = 무제한, 기본값: 0)

### 종료하기

#### 정상 종료 (Graceful Shutdown)

1. **Writer 프로세스 종료:**
   ```bash
   # 터미널에서 Ctrl+C 입력 또는 SIGTERM 신호 전송
   kill -SIGTERM <pid>
   ```

   **예상 출력:**
   ```
   🛑 Received signal: terminated. Shutting down...
   🗑️  Removing shared memory...
   ✅ Shutdown complete.
   ```

   **대기 시간:** 최대 30초 (구성 가능)
   - 진행 중인 작업이 완료될 때까지 대기
   - 타임아웃 후 강제 종료

2. **Reader 프로세스 종료:**
   ```bash
   # Ctrl+C 입력
   ```

---

## 헬스체크 및 모니터링

### 모니터링 서버 시작

```bash
go run cmd/monitor/main.go  # 미구현 - 수동 실행
```

기본 포트: `9090`

### 헬스체크 CLI 사용법

```bash
# 모든 검사 실행
go run cmd/health_check/main.go -addr localhost:9090 -check all

# 개별 검사
go run cmd/health_check/main.go -check health
go run cmd/health_check/main.go -check ready
go run cmd/health_check/main.go -check metrics

# JSON 형식 출력
go run cmd/health_check/main.go -check all -json
```

### HTTP 엔드포인트

#### `GET /health` - Liveness Probe

**목적:** 서비스가 살아있는지 확인

**응답:**
```json
{
  "status": "alive",
  "time": "2024-01-01T12:00:00Z"
}
```

**HTTP 상태 코드:** 200 (항상)

#### `GET /ready` - Readiness Probe

**목적:** 서비스가 트래픽을 처리할 준비가 되었는지 확인

**응답:**
```json
{
  "status": "ready",
  "time": "2024-01-01T12:00:00Z"
}
```

**HTTP 상태 코드:**
- 200: 준비 완료
- 503: 준비 안 됨

#### `GET /metrics` - Prometheus 메트릭

**목적:** 상세한 성능 메트릭 수집

**형식:** Prometheus text format (v0.0.4)

**예시:**
```
# HELP grie_jobs_received_total Total jobs received
# TYPE grie_jobs_received_total counter
grie_jobs_received_total 1000

# HELP grie_queue_size Current job queue size
# TYPE grie_queue_size gauge
grie_queue_size 42

# HELP grie_circuit_breaker_rejected_total Requests rejected by circuit breaker
# TYPE grie_circuit_breaker_rejected_total counter
grie_circuit_breaker_rejected_total 5
```

#### `GET /debug/pprof/` - 프로파일링

**목적:** 런타임 성능 분석

**사용 가능한 프로파일:**
- `/debug/pprof/heap` - 메모리 할당
- `/debug/pprof/goroutine` - 고루틴 정보
- `/debug/pprof/profile` - CPU 프로파일 (30초)
- `/debug/pprof/trace` - 실행 추적

---

## 메트릭 해석

### 작업 메트릭 (Job Metrics)

| 메트릭 | 의미 | 정상 범위 |
|--------|------|---------|
| `grie_jobs_received_total` | 총 수신 작업 수 | 계속 증가 |
| `grie_jobs_processed_total` | 성공적으로 처리된 작업 | >= dropped |
| `grie_jobs_failed_total` | 실패한 작업 | 낮음 |
| `grie_jobs_dropped_total` | 드롭된 작업 (백프레셔) | 0 또는 낮음 |
| `grie_jobs_retried_total` | 재시도된 작업 | 낮음 |

**주의:**
- 드롭이 증가하면 → 백프레셔 발생, 큐 크기 확인
- 실패가 많으면 → Circuit Breaker 상태 확인, 로그 분석

### 큐 메트릭 (Queue Metrics)

| 메트릭 | 의미 | 정상 범위 |
|--------|------|---------|
| `grie_queue_size` | 현재 큐에 있는 작업 | < 80% capacity |
| `grie_active_workers` | 활성 워커 수 | 구성된 워커 수 |

**주의:**
- 큐 크기 > 80% → 처리 속도 부족, 워커 추가 고려
- active_workers = 0 → 모든 워커 정지, Circuit Breaker 확인

### Circuit Breaker 메트릭

| 메트릭 | 의미 | 정상 범위 |
|--------|------|---------|
| `grie_circuit_breaker_allowed_total` | 허용된 요청 | 증가 추세 |
| `grie_circuit_breaker_rejected_total` | 거부된 요청 | 낮음 |

**상태 해석:**
- **CLOSED**: 정상 운영, 모든 요청 허용
- **OPEN**: 서비스 장애, 모든 요청 거부 (복구 대기)
- **HALF_OPEN**: 복구 테스트, 제한된 요청만 허용

---

## Circuit Breaker 이해

### 상태 전이 다이어그램

```
CLOSED ─(연속 5회 실패)→ OPEN
  ↑                        │
  │                        ↓
  └(연속 2회 성공)─ HALF_OPEN
           ↑
           │(30초 경과)
           │
         OPEN
```

### 설정 항목

```go
FailureThreshold:  5          // CLOSED → OPEN 전환 기준
RecoveryTimeout:   30s        // OPEN → HALF_OPEN 대기 시간
SuccessThreshold:  2          // HALF_OPEN → CLOSED 전환 기준
HalfOpenMaxCalls:  1          // HALF_OPEN에서 허용 요청 수
```

### 복구 절차

1. **Circuit Breaker가 OPEN이 되면:**
   ```
   시간: T=0s
   상태: OPEN
   거부: 모든 요청 → ErrCircuitOpen
   ```

2. **Recovery Timeout 경과 (30초):**
   ```
   시간: T=30s
   상태: HALF_OPEN 전환
   허용: 최대 1개 요청만
   ```

3. **HALF_OPEN 중 요청 결과:**
   - **성공**: 누적 성공 2회 → CLOSED 전환
   - **실패**: 즉시 OPEN으로 복귀 + 30초 재대기

### 수동 복구 (필요 시)

Circuit Breaker를 수동으로 리셋해야 하는 경우:

```go
// 코드에서 직접 리셋 (테스트/관리용)
circuitBreaker.Reset()
```

---

## 트러블슈팅

### 문제 1: "circuit breaker is open" 에러

**증상:**
```
circuit breaker is open
```

**원인:**
- 연속 5회 이상 요청 실패
- 서비스 또는 공유 메모리 장애

**해결 방법:**

1. **상태 확인:**
   ```bash
   curl localhost:9090/metrics | grep circuit_breaker
   ```

2. **로그 확인:**
   - 최근 에러 메시지 조사
   - SHM 또는 워커 상태 확인

3. **자동 복구 대기:**
   - 30초 후 자동으로 HALF_OPEN 상태로 전환
   - 이후 2회 연속 성공 시 CLOSED로 복구

4. **강제 복구 (응급):**
   - 프로세스 재시작

### 문제 2: 작업이 자꾸 드롭됨

**증상:**
```
grie_jobs_dropped_total 증가
```

**원인:**
- 큐가 80% 용량 초과 (백프레셔)
- 워커가 느림

**해결 방법:**

1. **큐 상태 확인:**
   ```bash
   curl localhost:9090/metrics | grep queue_size
   ```

2. **개선 방법:**
   - 큐 크기 증가: `NewDispatcher(workerCount, largercapacity)`
   - 워커 수 증가: `workerCount++`
   - 타임아웃 조정: `config.JobProcessTimeout` 증가

### 문제 3: 느린 응답 시간

**증상:**
```
요청 처리 시간 > 5초 (기본 타임아웃)
```

**원인:**
- 워커 부하 높음
- 공유 메모리 지연

**해결 방법:**

1. **메트릭 확인:**
   ```bash
   curl localhost:9090/metrics | grep -E "queue_size|active_workers"
   ```

2. **CPU/메모리 확인:**
   ```bash
   curl localhost:9090/debug/pprof/heap
   curl localhost:9090/debug/pprof/profile?seconds=30
   ```

3. **타임아웃 조정:**
   - `config.JobProcessTimeout` 증가 (보수적)
   - `config.SHMReadTimeout` 감소 (적극적)

### 문제 4: 정상 종료 실패

**증상:**
```
graceful shutdown timeout: jobs still in-flight after 30s
```

**원인:**
- 일부 작업이 응답하지 않음
- 워커가 멈춤

**해결 방법:**

1. **강제 종료 (마지막 수단):**
   ```bash
   kill -9 <pid>
   ```

2. **향후 예방:**
   - `config.ShutdownTimeout` 증가 (구성 파일에서)
   - `config.JobProcessTimeout` 감소 (더 빠른 실패)

### 문제 5: 메트릭 엔드포인트 접근 불가

**증상:**
```
curl: (7) Failed to connect to localhost port 9090
```

**원인:**
- 모니터링 서버 미시작
- 포트 충돌

**해결 방법:**

1. **프로세스 확인:**
   ```bash
   lsof -i :9090
   ```

2. **포트 변경:**
   ```bash
   ServerConfig{Addr: ":9091"}
   ```

3. **서비스 재시작**

---

## 성능 튜닝

### 추천 설정값

#### 고처리량 환경 (>1000 req/s)

```go
cfg := config.TimeoutConfig{
    JobProcessTimeout: 3 * time.Second,  // 단축
    SHMWriteTimeout:   500 * time.Millisecond,
    SHMReadTimeout:    50 * time.Millisecond,
    ShutdownTimeout:   60 * time.Second,  // 연장
}

breaker := circuitbreaker.BreakerConfig{
    FailureThreshold: 10,      // 좀 더 관대
    RecoveryTimeout:  60 * time.Second,
    SuccessThreshold: 5,       // 더 엄격한 복구
}
```

#### 저지연 환경 (<100ms 목표)

```go
cfg := config.TimeoutConfig{
    JobProcessTimeout: 100 * time.Millisecond,  // 짧음
    SHMWriteTimeout:   50 * time.Millisecond,
    SHMReadTimeout:    10 * time.Millisecond,
    ShutdownTimeout:   10 * time.Second,
}

breaker := circuitbreaker.BreakerConfig{
    FailureThreshold: 3,       // 엄격
    RecoveryTimeout:  10 * time.Second,
    SuccessThreshold: 1,       // 빠른 복구
}
```

---

## 체크리스트

### 배포 전
- [ ] 모든 테스트 통과 (`go test -race ./...`)
- [ ] 모니터링 서버 포트 확인
- [ ] 공유 메모리 크기 충분한지 확인
- [ ] 타임아웃 값 환경에 맞게 조정
- [ ] Circuit Breaker 설정 검토

### 배포 후 (1시간)
- [ ] `/health` 엔드포인트 정상
- [ ] `/ready` 엔드포인트 200 반환
- [ ] 메트릭 수집 확인
- [ ] Circuit Breaker CLOSED 상태
- [ ] 드롭된 작업 0 (정상 로드)
- [ ] 활성 워커 수 >= 1

### 일일 점검
- [ ] 에러 로그 없음
- [ ] Circuit Breaker 토글 없음
- [ ] 평균 응답시간 정상
- [ ] 메모리 누수 없음 (heap 프로파일)
- [ ] 업타임 연속

---

**문서 버전:** Phase 3 v1.0
**마지막 업데이트:** 2026-02-27
**다음 검토:** Phase 4 시작 시
