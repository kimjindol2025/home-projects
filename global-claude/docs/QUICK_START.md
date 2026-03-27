# 빠른 시작 가이드

## 5분 안에 시작하기

### 1단계: 저장소 클론 (1분)

```bash
git clone https://gogs.dclub.kr/kim/global-claude.git
cd global-claude
```

### 2단계: 파이썬 버전 확인 (30초)

```bash
python3 --version
# Python 3.7 이상 필요
```

### 3단계: 예제 실행 (3분 30초)

```bash
# 가장 간단한 예제 (권장)
python3 examples/01-simple-multiagent.py
```

### 예상 출력

```
======================================================================
💬 멀티프로세스 Claude 에이전트 - 자유로운 대화
======================================================================

──────────────────────────────────────────────────────────────────────
[Step 1] 에이전트 프로세스 생성

✅ [optimist] 프로세스 시작 (PID: 12345)
✅ [realist] 프로세스 시작 (PID: 12346)
✅ [critic] 프로세스 시작 (PID: 12347)

──────────────────────────────────────────────────────────────────────
[Step 2] 대화 시작

📤 [user] → [optimist] 메시지 전송
    새로운 AI 제품을 개발해야 할까요?...

──────────────────────────────────────────────────────────────────────
[Step 3] 대화 진행 중... (30초)

[12:58:01] 💬 [optimist] 📨 [user]로부터 메시지 수신
[12:58:01] 💬 [optimist]    내용: 새로운 AI 제품을 개발해야 할까요?...
[12:58:02] 💬 [optimist] ✅ 응답 생성
[12:58:02] 💬 [optimist] 📤 [realist]로 응답 전송

[12:58:03] 💬 [realist] 📨 [optimist]로부터 메시지 수신
[12:58:03] 💬 [realist]    내용: 좋은 질문입니다! 제 생각:...
[12:58:04] 💬 [realist] ✅ 응답 생성
[12:58:04] 💬 [realist] 📤 [critic]로 응답 전송

[12:58:05] 💬 [critic] 📨 [realist]로부터 메시지 수신
[12:58:05] 💬 [critic]    내용: 현실적으로 분석하면:...
[12:58:06] 💬 [critic] ✅ 응답 생성
[12:58:06] 💬 [critic] 📤 [optimist]로 응답 전송

... (30초 동안 계속 대화)

──────────────────────────────────────────────────────────────────────
[Step 4] 모든 에이전트 종료

✅ [optimist] 종료됨
✅ [realist] 종료됨
✅ [critic] 종료됨

======================================================================
✨ 대화 완료!
======================================================================
```

## 코드 이해하기

### 최소 예제 (30줄)

```python
import multiprocessing as mp
import time
from datetime import datetime

class SimpleAgent:
    def __init__(self, name, queue):
        self.name = name
        self.queue = queue

    def run(self):
        print(f"🚀 [{self.name}] 시작")

        # 10번 메시지 처리
        for i in range(10):
            if not self.queue.empty():
                msg = self.queue.get()

                # STOP 신호 확인
                if msg.get("type") == "STOP":
                    break

                # 메시지 처리
                response = f"Response from {self.name}"
                self.queue.put({
                    "from": self.name,
                    "content": response
                })

            time.sleep(0.1)

        print(f"🔴 [{self.name}] 종료")

if __name__ == "__main__":
    # 공유 큐 생성
    queue = mp.Queue()

    # 에이전트 시작
    agent = SimpleAgent("Worker", queue)
    process = mp.Process(target=agent.run)
    process.start()

    # 초기 메시지 전송
    queue.put({"from": "user", "content": "Hello!"})

    # 30초 대기
    time.sleep(30)

    # 종료
    queue.put({"type": "STOP"})
    process.join()
```

### 실행 결과

```
🚀 [Worker] 시작
🔴 [Worker] 종료
```

## 다양한 시나리오

### 시나리오 1: 세 에이전트가 대화하기

```bash
python3 examples/01-simple-multiagent.py
```

**실행 흐름:**
1. 3개 에이전트 시작 (optimist, realist, critic)
2. User가 optimist에 질문 전송
3. 에이전트들이 자동으로 대화
4. 30초 후 모든 에이전트 종료

### 시나리오 2: 메시지 히스토리 추적

```bash
python3 examples/02-conversational-agent.py
```

**추가 기능:**
- 모든 메시지 기록
- 메시지 타입 분류
- 대화 로그 출력

### 시나리오 3: 실시간 라우팅

```bash
python3 examples/03-interactive-realtime.py
```

**추가 기능:**
- 별도 Router 스레드
- 자동 메시지 라우팅
- 높은 처리량

### 시나리오 4: 프로덕션 오케스트레이션

```bash
python3 examples/04-orchestrator-pattern.py
```

**추가 기능:**
- 중앙 오케스트레이터
- 상태 모니터링
- 상세 로깅

## 커스터마이제이션

### 새 에이전트 추가

```python
# 1. ChatAgent 상속
class EngineerAgent(ChatAgent):
    def __init__(self, name, message_queue):
        super().__init__(name, "engineer", message_queue)

    def _respond(self, content):
        return f"기술적 관점: {content}"

# 2. 에이전트 시작
agent = EngineerAgent("engineer", queue)
p = ctx.Process(target=agent.run)
p.start()
```

### 성격 커스터마이징

```python
class MyAgent(ChatAgent):
    def _respond(self, content):
        if self.personality == "engineer":
            return f"""
            엔지니어 의견:
            📋 기술적 검토
            🔧 구현 방안
            ⚡ 성능 최적화
            """
        return "기본 응답"
```

### 메시지 필터링

```python
def run(self):
    while True:
        if not self.queue.empty():
            msg = self.queue.get(timeout=0.5)

            # 자신의 메시지 스킵
            if msg.get("from") == self.name:
                continue

            # 특정 타입만 처리
            if msg.get("type") == "question":
                response = self._respond(msg["content"])
```

## 문제 해결

### 에러: "ModuleNotFoundError: No module named 'multiprocessing'"

```bash
# Python 3 사용 확인
python3 --version

# 올바른 python3 사용
python3 examples/01-simple-multiagent.py
```

### 에러: "ProcessError: spawn not available"

```bash
# macOS/Windows에서 발생 가능
# 해결: 스크립트 마지막에 추가
if __name__ == "__main__":
    main()
```

### 에러: "Queue is empty" (Queue.get() 에러)

```python
# ❌ 잘못된 사용
msg = queue.get()  # 메시지 없으면 무한 대기

# ✅ 올바른 사용
msg = queue.get(timeout=1)  # 1초 후 타임아웃

# 또는
if not queue.empty():
    msg = queue.get()
```

### 성능 저하 (응답 느림)

```python
# 원인 1: 에이전트가 STOP 신호 받지 못함
# 해결: STOP 신호 여러 번 전송
for _ in range(num_agents):
    queue.put({"type": "STOP"})

# 원인 2: 메시지 큐 가득 참
# 해결: 큐 크기 늘리기
queue = mp.Queue(maxsize=1000)

# 원인 3: 프로세스 리소스 부족
# 해결: 에이전트 수 줄이기
```

## 다음 단계

### 기본 학습
1. ✅ [01-simple-multiagent.py](../examples/01-simple-multiagent.py) 실행
2. 📖 [ARCHITECTURE.md](./ARCHITECTURE.md) 읽기
3. 🔧 코드 수정하여 자신의 에이전트 추가

### 심화 학습
1. 📚 [Python multiprocessing 공식 문서](https://docs.python.org/3/library/multiprocessing.html)
2. 🎓 [분산 시스템 개념](https://en.wikipedia.org/wiki/Distributed_computing)
3. 💡 다른 예제 살펴보기 (02, 03, 04)

### 프로덕션 준비
1. 에러 핸들링 추가
2. 로깅 시스템 구축
3. 모니터링 대시보드 연결
4. 네트워크 확장 (원격 에이전트)

## 유용한 팁

### Tip 1: 디버깅
```bash
# 프로세스 상태 확인
ps aux | grep python

# 메모리 사용량 확인
ps -o pid,rss,comm
```

### Tip 2: 로깅
```python
import logging

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

logger.debug(f"[{self.name}] 메시지 수신")
```

### Tip 3: 프로파일링
```python
import cProfile

cProfile.run('main()')
```

### Tip 4: 타임아웃 설정
```python
# 프로세스 종료 타임아웃
process.join(timeout=5)

# 큐 접근 타임아웃
msg = queue.get(timeout=1)
```

## 자주 묻는 질문 (FAQ)

**Q: 네트워크를 통해 다른 머신의 에이전트와 통신할 수 있나요?**
A: 현재는 로컬 IPC만 지원. 네트워크 확장은 Road Map에 있습니다.

**Q: 에이전트 수를 늘리면 어떻게 되나요?**
A: 메모리와 성능에 영향. 10개까지는 무방하지만, 100개 이상은 네트워크 기반 아키텍처 권장.

**Q: 실제 LLM을 사용할 수 있나요?**
A: 네! `_respond()` 메서드를 Anthropic API 호출로 변경하면 됩니다.

**Q: 메시지 손실이 발생할 수 있나요?**
A: mp.Queue는 신뢰성이 높으나, 시스템 크래시 시 손실 가능. 중요한 메시지는 디스크에 기록하세요.

---

**More questions?** → [ARCHITECTURE.md](./ARCHITECTURE.md) 참고 또는 GOGS 이슈 생성
