# 글로드 클로드 아키텍처

## 시스템 설계

### 고수준 다이어그램

```
┌────────────────────────────────────────────────────────┐
│            Main Process (오케스트레이터)               │
│  - 프로세스 생성/종료                                  │
│  - 초기 메시지 전송                                    │
│  - 30초 대기                                           │
└────────────────────────────────────────────────────────┘
           ▲                      ▲                    ▲
           │                      │                    │
           │ Global Message       │ Global Message     │ Global Message
           │ Queue                │ Queue              │ Queue
           │                      │                    │
      ┌────▼────┐         ┌──────▼──┐         ┌───────▼────┐
      │Optimist │         │ Realist │         │   Critic   │
      │Process  │         │Process  │         │  Process   │
      │(PID:X)  │         │(PID:Y)  │         │ (PID:Z)    │
      └────┬────┘         └────┬────┘         └────┬───────┘
           │                   │                    │
           └───────────────────┼────────────────────┘
                    📨 Message Flow
```

### 핵심 컴포넌트

#### 1. Global Message Queue (전역 메시지 큐)

```python
message_queue = mp.Queue()  # 모든 프로세스가 공유
```

**역할**:
- 중앙 메시지 허브
- 에이전트 간 양방향 통신
- FIFO 순서 보장
- 스레드/프로세스 안전

**특징**:
- 무제한 크기 (시스템 메모리 한계)
- 블로킹 가능 (timeout 설정)
- 자동 직렬화/역직렬화

#### 2. ChatAgent (에이전트 프로세스)

```python
class ChatAgent:
    def __init__(self, name, personality, message_queue):
        self.name = name
        self.personality = personality
        self.queue = message_queue
        self.message_count = 0

    def run(self):  # 메인 이벤트 루프
        while True:
            if not self.queue.empty():
                msg = self.queue.get(timeout=0.5)
                response = self._respond(msg["content"])
                self.queue.put({...})  # 응답 전송
```

**상태 머신**:
```
┌─────────┐
│ Running │
└────┬────┘
     │ msg.type == "STOP"
     ▼
┌─────────┐
│Stopped  │
└─────────┘
```

**메시지 처리 플로우**:
```
1. queue.get() 에서 메시지 수신 (타임아웃 0.5초)
2. 자신이 보낸 메시지면 스킵 (무한 루프 방지)
3. _respond() 로 성격에 맞는 응답 생성
4. 응답 메시지를 queue.put() 으로 전송
5. 타임아웃 경과 또는 STOP 신호 수신 시 종료
```

#### 3. Response Generation (응답 생성)

```python
def _respond(self, content):
    if self.personality == "optimist":
        return """긍정적 응답..."""
    elif self.personality == "realist":
        return """현실적 분석..."""
    elif self.personality == "critic":
        return """비판적 검토..."""
```

**특징**:
- 결정적 응답 (LLM 없이)
- 성격 기반 다양성
- 즉시 반환 (지연 없음)
- 시뮬레이션용 1초 sleep

## 메시지 프로토콜

### 메시지 형식

```json
{
    "from": "string",           // 발신자 이름
    "to": "string",             // 수신자 이름
    "content": "string",        // 메시지 본문
    "type": "string",           // message, response, STOP
    "timestamp": "ISO 8601"     // 메시지 시간
}
```

### 메시지 플로우 예시

```
시간    Optimist        Queue          Realist      Critic
─────────────────────────────────────────────────────────
T0      발송 대기
T1                    [User Question]
T1                    → Optimist.queue
T2      수신 ✓
        생각 중 ⏳
T3      응답 생성 ✓
                    [Optimist → Realist]
T3                    → Realist.queue
T4                                  수신 ✓
T4                    [Optimist → Critic]
T4                    → Critic.queue
T5                                             수신 ✓
```

## 병렬 처리 모델

### 진정한 병렬 실행

Python의 `multiprocessing` 은 별도 **프로세스**를 생성하므로:

```
프로세스 1 (Optimist):    [작업] [대기] [작업] [대기]
프로세스 2 (Realist):     [대기] [작업] [대기] [작업]
프로세스 3 (Critic):      [대기] [대기] [작업] [대기]
───────────────────────────────────────────────────
실제 CPU 시간                 ←→ 모두 동시 실행
```

### GIL 회피

```python
# ❌ 스레드 (threading): GIL로 인해 직렬화됨
import threading
t1 = threading.Thread(target=agent1.run)
t2 = threading.Thread(target=agent2.run)
# → 한 번에 하나만 실행

# ✅ 프로세스 (multiprocessing): 진정한 병렬 실행
import multiprocessing
p1 = multiprocessing.Process(target=agent1.run)
p2 = multiprocessing.Process(target=agent2.run)
# → 여러 CPU 코어에서 동시 실행
```

## 메모리 모델

### 메모리 격리

각 프로세스는 **완전히 독립적인 메모리 공간**을 가짐:

```
┌─────────────────────────┐
│ Main Process (메인)      │
│  ├─ globals             │
│  ├─ message_queue (공유)│
│  └─ processes (리스트)  │
└─────────────────────────┘
     ↑ spawn
     │
     ├─→ ┌──────────────┐
     │   │ Optimist     │
     │   │ (독립 메모리)│
     │   │ queue(참조)  │
     │   └──────────────┘
     │
     ├─→ ┌──────────────┐
     │   │ Realist      │
     │   │ (독립 메모리)│
     │   │ queue(참조)  │
     │   └──────────────┘
     │
     └─→ ┌──────────────┐
         │ Critic       │
         │ (독립 메모리)│
         │ queue(참조)  │
         └──────────────┘
```

### 공유 상태

오직 **메시지 큐**만 공유:

```python
# ✓ 안전 (Queue는 스레드/프로세스 안전)
message_queue = mp.Queue()

# ✗ 위험 (일반 dict는 불안전)
shared_dict = {}

# ✓ 안전 (manager.dict 사용 시)
manager = mp.Manager()
shared_dict = manager.dict()
```

## 생명주기 관리

### 프로세스 시작

```python
# 1. 컨텍스트 선택 (spawn)
ctx = mp.get_context('spawn')

# 2. 에이전트 인스턴스 생성
agent = ChatAgent("optimist", "optimist", global_queue)

# 3. 프로세스 시작
process = ctx.Process(target=agent.run)
process.start()

# 4. PID 확인
print(f"Started with PID: {process.pid}")
```

### 프로세스 종료

```python
# 방법 1: STOP 신호 (안전한 종료)
queue.put({"type": "STOP"})
process.join(timeout=3)  # 최대 3초 대기

# 방법 2: 강제 종료
if process.is_alive():
    process.terminate()

# 방법 3: 강력한 강제 종료
process.kill()
```

### Spawn vs Fork

```python
# spawn (안전, 권장)
ctx = mp.get_context('spawn')
# 각 프로세스가 처음부터 시작

# fork (빠름, Unix/Linux만)
ctx = mp.get_context('fork')
# 부모 프로세스 복제 (메모리 대량 복사)
```

## 동시성 제어

### 메시지 큐 (내장 동시성)

```python
# mp.Queue 는 자동으로 스레드/프로세스 안전성 보장
queue.put(msg)      # 원자적 (atomic)
msg = queue.get()   # 원자적 (atomic)
```

### 더 복잡한 동기화 필요 시

```python
# Lock 사용
lock = mp.Lock()
with lock:
    # 임계 영역
    shared_resource = value

# Semaphore 사용
sem = mp.Semaphore(3)  # 최대 3개 프로세스만 진입
with sem:
    critical_section()
```

## 성능 특성

### 타이밍 분석

```
에이전트 시작:          ~100ms (프로세스 생성)
메시지 전송/수신:       <1ms (IPC)
응답 생성:              1000ms (time.sleep)
메시지 왕복:            ~3000ms (3개 에이전트)
```

### 확장성 특성

| 에이전트 수 | 메모리 | 시작 시간 | 지연 |
|-----------|--------|----------|------|
| 1개       | 20MB   | 50ms     | 1s   |
| 3개       | 60MB   | 150ms    | 3s   |
| 10개      | 200MB  | 500ms    | 10s  |
| 100개     | 2GB    | 5s       | 100s |

### 병목

1. **프로세스 생성 오버헤드**: fork/spawn 비용
2. **IPC 지연**: 메시지 직렬화/역직렬화
3. **응답 시간**: 각 에이전트의 처리 시간
4. **메모리 오버헤드**: 프로세스당 ~20MB 기본

## 장단점

### 장점
- ✅ 진정한 병렬 처리 (GIL 회피)
- ✅ 완전한 메모리 격리 (안정성)
- ✅ 프로세스 간 독립적 실행 (장애 격리)
- ✅ 확장 가능 (분산 시스템으로 확장 가능)
- ✅ 디버깅 용이 (독립적 로그)

### 단점
- ❌ 높은 메모리 오버헤드 (프로세스당 ~20MB)
- ❌ IPC 지연 (스레드보다 느림)
- ❌ 공유 상태 관리 복잡
- ❌ 프로세스 생성 비용 (첫 시작 느림)

## 향후 개선

### 1단계: 네트워크 확장
```python
# TCP/IP를 통한 원격 에이전트
agent = RemoteAgent("agent1", "tcp://10.0.0.2:5000")
```

### 2단계: 에이전트 풀
```python
# 동적 에이전트 생성/제거
pool = AgentPool(size=5)
pool.submit_task(task)
```

### 3단계: 상태 관리
```python
# Redis/etcd 를 통한 분산 상태
shared = DistributedDict(backend="redis")
```

### 4단계: 장애 복구
```python
# 프로세스 감시 및 자동 재시작
supervisor = ProcessSupervisor()
supervisor.watch_and_restart()
```

---

이 아키텍처는 **간단하면서도 확장 가능한** 멀티프로세스 협업 시스템의 기초를 제공합니다.
