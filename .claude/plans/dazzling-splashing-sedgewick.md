# Pulse AI - 분산 에이전트 아키텍처 구현 계획

## Context

Phase 4까지 완성된 Pulse AI를 **다중 인스턴스(분산) 환경**으로 확장합니다.
현재는 단일 프로세스에 모든 에이전트 상태와 작업이 묶여 있어 확장이 불가능합니다.
목표: 워커 3개가 동시에 에이전트 작업을 처리하고, 장애 시 자동 Failover.

---

## 핵심 문제 (현재 단일 인스턴스)

| 위치 | 문제 |
|------|------|
| `AgentRegistry.agents` → 로컬 딕셔너리 | 다른 인스턴스 접근 불가 |
| `Agent.tasks` → `asyncio.Task` | 프로세스 종료 시 소멸 |
| `Agent.id` → uuid4 앞 8자리 | 인스턴스 간 충돌 가능 |
| `increment_cost` → GET+SET 분리 | race condition |

---

## 구현 순서 (의존성 순)

```
1. cache.py       수정 (가장 먼저 - 기반)
2. registry.py    수정 (cache에 의존)
3. agent.py       수정 (registry + cache에 의존)
4. distributed/worker.py  신규 (agent 의존)
5. distributed/cluster.py 신규 (worker + registry 의존)
6. gateway.py     수정 (cluster 의존)
7. docker-compose.yml 업데이트 (최종)
```

---

## 파일별 구현 명세

### 1. `src/cache.py` - increment_cost 원자화

```python
# 변경: GET+SET → INCRBYFLOAT (원자 연산)
async def increment_cost(self, amount: float) -> float:
    new_val = await self.redis.incrbyfloat("daily_cost", amount)
    await self.redis.expire("daily_cost", 86400)
    return new_val
```

### 2. `src/registry.py` - Redis Hash로 전환

Redis 키 설계:
```
pulse:agent:{uuid4}             → Hash  (에이전트 상태)
pulse:agents:index              → Set   (전체 ID 목록)
pulse:agents:by_status:{status} → Set   (상태별 인덱스)
```

주요 변경:
- `self.agents` 딕셔너리 제거
- 모든 메서드를 Redis Hash CRUD로 교체
- `update_status()` → async + Redis 갱신 + 인덱스 갱신

### 3. `src/agent.py` - 직렬화 + Redis 큐잉

```python
# Agent.id: 전체 uuid4 사용
self.id = str(uuid.uuid4())

# 직렬화 메서드 추가
def to_dict(self) -> Dict[str, Any]
@classmethod
def from_dict(cls, data: Dict, llm_manager=None) -> "Agent"

# asyncio.Task → Redis Stream 큐잉
await redis.xadd("pulse:jobs", {"data": json.dumps(job)})
```

AgentPool:
- `self.agents` 딕셔너리 제거
- `create_agent()`, `get_agent()`, `list_agents()` → Registry 경유

### 4. `src/distributed/worker.py` (신규)

```
Redis Streams 구조:
  Stream:          pulse:jobs
  Consumer Group:  pulse:workers
  Consumer ID:     worker:{hostname}:{pid}
```

주요 클래스 `DistributedWorker`:
- `start()` - Consumer Group 생성 + 작업 루프 시작
- `_job_loop()` - XREADGROUP으로 작업 수신 (5초 블로킹)
- `_execute_job(entry_id, job)` - 에이전트 복원 후 실행
- `_handle_job_failure()` - 3회 재시도 → Dead Letter Queue
- `_heartbeat_loop()` - 10초마다 생존 신호 (30초 TTL)
- `stop()` - 진행 중 작업 완료 후 종료

진입점:
```bash
python -m src.distributed.worker
```

### 5. `src/distributed/cluster.py` (신규)

Redis 키:
```
pulse:cluster:workers     → Hash  (worker_id → 상태 JSON, 30초 TTL)
pulse:cluster:assignments → Hash  (agent_id → worker_id)
pulse:cluster:load        → Sorted Set (worker_id, score=job_count)
```

주요 클래스 `ClusterManager`:
- `get_active_workers()` - 30초 이내 응답한 워커만 반환
- `select_worker(strategy)` - least_loaded / round_robin
- `assign_job(agent_id, job)` - 워커 선택 + xadd + 부하 카운터 증가
- `release_job(agent_id)` - 부하 카운터 감소
- `detect_dead_workers()` - 30초 이상 무응답 감지
- `failover(dead_worker_id)` - XCLAIM으로 작업 이전
- `start_failover_monitor()` - 15초마다 장애 감시
- `get_cluster_status()` - /health 엔드포인트용

### 6. `src/gateway.py` - 최소 수정

추가:
```python
from .distributed.cluster import ClusterManager

# startup: ClusterManager 초기화 + failover_monitor 백그라운드 시작
# /pulse: cluster.assign_job() 으로 배포
# /health: cluster_status 추가
# /cluster: 신규 엔드포인트 (워커 목록 + 부하 현황)
```

변경:
- `registry.update_status()` 호출에 `await` 추가

### 7. `docker-compose.yml` - 워커 서비스 추가

```yaml
pulse-worker:
  build: .
  command: python -m src.distributed.worker
  environment:
    - REDIS_URL=redis://redis:6379/0
    - WORKER_MAX_JOBS=10
  deploy:
    replicas: 3
    restart_policy:
      condition: on-failure
```

스케일링:
```bash
docker-compose up --scale pulse-worker=5  # 5개로 확장
```

---

## 신규 파일 목록

```
src/distributed/
├── __init__.py
├── worker.py    (DistributedWorker 클래스)
└── cluster.py   (ClusterManager 클래스)
```

---

## 수정 파일 목록

| 파일 | 변경 범위 |
|------|-----------|
| `src/cache.py` | `increment_cost()` 1개 메서드 |
| `src/registry.py` | 클래스 전체 재구현 (인터페이스 유지) |
| `src/agent.py` | `Agent.id`, `to_dict`, `from_dict`, `AgentPool` |
| `src/gateway.py` | startup, /pulse, /health, await 추가 |
| `docker-compose.yml` | pulse-worker 서비스 추가 |

---

## Redis 키 전체 구조

```
# 에이전트 상태
pulse:agent:{uuid4}                → Hash
pulse:agents:index                 → Set
pulse:agents:by_status:{status}    → Set

# 작업 큐
pulse:jobs                         → Stream
pulse:jobs:dead                    → Stream (Dead Letter)

# 클러스터
pulse:cluster:workers              → Hash (30초 TTL)
pulse:cluster:assignments          → Hash
pulse:cluster:load                 → Sorted Set

# 기존 유지
daily_cost, nlu:{hash}, llm:{hash} ...
```

---

## 검증 방법

```bash
# 1. 기본 동작
docker-compose up --scale pulse-worker=3
curl -X POST http://localhost:8000/pulse \
  -d '{"instruction": "슬랙 봇 만들어줘"}'

# 2. 클러스터 상태 확인
curl http://localhost:8000/health
# → "cluster": {"active_workers": 3, "pending_jobs": 0}

# 3. Failover 테스트
docker-compose stop pulse-worker  # 워커 1개 강제 종료
# → 15초 내에 작업이 다른 워커로 이전되는지 확인

# 4. 스케일 업
docker-compose up --scale pulse-worker=5
curl http://localhost:8000/health
# → "active_workers": 5

# 5. 부하 분산 확인
for i in {1..10}; do
  curl -X POST http://localhost:8000/pulse \
    -d '{"instruction": "봇 만들어줘"}'
done
# → 각 워커에 고르게 분배되는지 확인
```

---

## 마이그레이션 (무중단)

```
Phase 1: cache.py 수정 → 즉시 배포 가능 (하위 호환)
Phase 2: registry.py 수정 → Gateway 재시작 필요
Phase 3: agent.py + distributed/ 추가 → 워커 컨테이너 추가
Phase 4: docker-compose.yml → rolling restart
```

---

**상태**: 📋 계획 완료 - 구현 승인 대기
**예상 파일**: 7개 수정/추가
**추가 의존성**: 없음 (redis[asyncio] 이미 포함)
