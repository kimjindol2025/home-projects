# 🌍 글로드 클로드 (Global Claude)

멀티프로세스 기반의 독립적인 Claude AI 에이전트 협업 시스템

## 개요

여러 개의 Claude AI 에이전트가 **독립 프로세스**로 실행되면서 **중앙 메시지 큐**를 통해 자유롭게 양방향 대화를 나누는 시스템입니다.

- ✅ **완전한 독립성**: 각 에이전트가 별도 프로세스에서 실행
- ✅ **IPC 기반 통신**: 프로세스 간 큐를 통한 메시지 라우팅
- ✅ **병렬 처리**: 진정한 GIL 회피 병렬 실행
- ✅ **성격 기반 대화**: 각 에이전트의 개성이 살아있는 응답 생성
- ✅ **확장성**: 에이전트 추가/제거 용이

## 핵심 개념

### 메시지 기반 아키텍처

```
┌─────────────────────────────────────────┐
│      Global Message Queue (mp.Queue)    │
└─────────────────────────────────────────┘
     ▲                    ▲                ▲
     │                    │                │
  [Process 1]        [Process 2]      [Process 3]
   Optimist           Realist          Critic
```

### 에이전트 생명주기

1. **초기화**: 에이전트가 독립 프로세스로 시작
2. **구독**: 전역 큐에서 자신 이외의 메시지 수신
3. **처리**: 성격에 맞는 응답 생성 (1초 사고 시간)
4. **발행**: 응답을 큐에 돌려보냄
5. **종료**: STOP 신호 수신 시 정상 종료

## 빠른 시작

### 설치 및 실행

```bash
# 저장소 클론
git clone https://gogs.dclub.kr/kim/global-claude.git
cd global-claude

# 기본 예제 실행 (권장)
python examples/01-simple-multiagent.py

# 실행 결과 예시
# ════════════════════════════════════════
# 💬 멀티프로세스 Claude 에이전트
# ════════════════════════════════════════
#
# [Step 1] 에이전트 프로세스 생성
# ✅ [optimist] 프로세스 시작 (PID: 12345)
# ✅ [realist] 프로세스 시작 (PID: 12346)
# ✅ [critic] 프로세스 시작 (PID: 12347)
#
# [Step 2] 대화 시작
# 📤 [user] → [optimist] 메시지 전송
#    "새로운 AI 제품을 개발해야 할까요?"
#
# [12:58:01] 💬 [optimist] 📨 [user]로부터 메시지 수신
# [12:58:02] 💬 [optimist] ✅ 응답 생성
# [12:58:02] 💬 [optimist] 📤 [realist]로 응답 전송
#
# ... (30초 동안 자유로운 대화)
#
# ✨ 대화 완료!
```

## 구현 방식별 예제

### 1. 간단한 멀티에이전트 (권장)
```bash
python examples/01-simple-multiagent.py
```
- **특징**: 최소한의 복잡도, 가장 안정적
- **규모**: 260줄
- **사용**: 전역 mp.Queue 하나로 모든 에이전트 통신
- **성능**: 30초 내 30+ 메시지 처리
- **추천**: 처음 사용자, 프로토타입

### 2. 대화형 에이전트
```bash
python examples/02-conversational-agent.py
```
- **특징**: Message 클래스, MessageBroker 패턴
- **규모**: 450줄
- **사용**: 메시지 타입 분류 (chat, question, answer, broadcast)
- **추가 기능**: 메시지 히스토리 추적
- **추천**: 복잡한 대화 흐름 필요할 때

### 3. 실시간 라우팅
```bash
python examples/03-interactive-realtime.py
```
- **특징**: 스레드 기반 실시간 메시지 라우팅
- **규모**: 380줄
- **사용**: 분리된 inbox/outbox 큐
- **추가 기능**: Router 스레드가 메시지 자동 라우팅
- **추천**: 높은 처리량 필요할 때

### 4. 오케스트레이터 패턴
```bash
python examples/04-orchestrator-pattern.py
```
- **특징**: 중앙 집중식 오케스트레이터
- **규모**: 280줄
- **사용**: 에이전트 생명주기 완전 제어
- **추가 기능**: 프로세스 상태 모니터링, 메시지 로깅
- **추천**: 프로덕션 환경, 상세 추적 필요할 때

## 아키텍처 비교

| 특징 | 간단 | 대화형 | 실시간 | 오케스트레이터 |
|------|------|--------|---------|-------------|
| 복잡도 | ⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| 메모리 사용 | 낮음 | 중간 | 높음 | 중간 |
| 확장성 | 좋음 | 좋음 | 최고 | 좋음 |
| 디버깅 | 쉬움 | 중간 | 어려움 | 쉬움 |
| 프로덕션 준비 | ✅ | ✅ | ⚠️ | ✅ |

## 에이전트 성격 시스템

### 기본 성격 (3가지)

**Optimist (낙관주의자)**
```python
"좋은 질문입니다! 제 생각:
✅ 긍정적 접근이 필요합니다
💡 시도해볼 가치가 충분합니다
🎯 함께 성공할 수 있을 겁니다!"
```

**Realist (현실주의자)**
```python
"현실적으로 분석하면:
📊 장단점을 균형있게 봐야 합니다
⚠️ 리스크 관리가 중요합니다
✓ 실행 계획이 명확해야 합니다"
```

**Critic (비판주의자)**
```python
"비판적 검토:
🔍 몇 가지 우려사항이 있습니다
⚡ 더 견고한 근거가 필요합니다
💭 이런 점들을 개선해야 합니다"
```

### 커스텀 성격 추가

```python
def add_custom_personality():
    """새로운 성격 추가 예시"""
    class CustomAgent(ChatAgent):
        def _respond(self, content):
            # 커스텀 로직
            return "당신의 아이디어에 대해..."
```

## 메시지 형식

### 메시지 구조
```python
{
    "from": "optimist",        # 발신자
    "to": "realist",           # 수신자
    "content": "...",          # 본문
    "type": "response",        # 메시지 타입
    "timestamp": "2026-03-27T12:58:01.234567"
}
```

### 제어 메시지
```python
# 에이전트 종료
{"type": "STOP"}

# 종료 확인
if msg.get("type") == "STOP":
    break
```

## 성능 특성

### 벤치마크
- **에이전트 수**: 3개
- **메시지 처리량**: 30초 내 30+ 메시지
- **응답 지연**: 1초/메시지 (생각 시간 포함)
- **메모리 사용**: ~50MB (전체 프로세스)
- **CPU 사용**: 저 (대부분 I/O 대기)

### 확장성
- **최대 에이전트 수**: 테스트됨 (3개), 이론상 100+ 가능
- **메시지 큐 크기**: 무제한 (시스템 메모리 한계까지)
- **지연 시간**: 선형 증가

## 문제 해결

### 에이전트가 응답하지 않음
```python
# 체크리스트
1. STOP 신호가 강제로 발송되지 않았는지 확인
2. 메시지 큐가 가득 차지 않았는지 확인
3. 프로세스가 살아있는지 확인: ps aux | grep python
```

### 메시지 손실
```python
# 원인: Queue 크기 제한 또는 처리 속도 초과
# 해결: mp.Queue() 크기 늘리기
message_queue = mp.Queue(maxsize=1000)
```

### 데드락
```python
# 예방: STOP 신호 여러 번 전송
for _ in range(len(processes)):
    global_queue.put({"type": "STOP"})
```

## 고급 사용

### 커스텀 에이전트 생성
```python
class SpecialistAgent(ChatAgent):
    def __init__(self, name, specialty, message_queue):
        super().__init__(name, "specialist", message_queue)
        self.specialty = specialty

    def _respond(self, content):
        if self.specialty == "math":
            return f"수학적 분석: {content}"
        elif self.specialty == "code":
            return f"코드 리뷰: {content}"
```

### 메시지 필터링
```python
def run(self):
    while True:
        if not self.queue.empty():
            msg = self.queue.get(timeout=0.5)

            # 필터: 자신의 메시지 스킵
            if msg.get("from") == self.name:
                continue

            # 필터: 특정 타입만 처리
            if msg.get("type") != "question":
                continue
```

### 에이전트 간 상태 공유
```python
# manager.dict() 사용 예시
from multiprocessing import Manager
manager = Manager()
shared_knowledge = manager.dict()

# 에이전트 내부에서
shared_knowledge["last_decision"] = response
```

## API 레퍼런스

### ChatAgent
```python
class ChatAgent:
    def __init__(self, name: str, personality: str, message_queue: mp.Queue)
    def run(self) -> None
    def _respond(self, content: str) -> str
    def _log(self, msg: str) -> None
```

### Process Control
```python
# 프로세스 시작
process = ctx.Process(target=agent.run)
process.start()

# 프로세스 종료
queue.put({"type": "STOP"})
process.join(timeout=3)
if process.is_alive():
    process.terminate()
```

## 라이선스

MIT

## 참고 자료

- [Python multiprocessing 문서](https://docs.python.org/3/library/multiprocessing.html)
- [메시지 큐 패턴](https://en.wikipedia.org/wiki/Message_queue)
- [마이크로서비스 아키텍처](https://martinfowler.com/articles/microservices.html)

## 문의 & 기여

이슈, 질문, 개선사항은 GOGS 저장소에서 관리합니다.

---

**Made with ❤️ by FreeLang Marketing Team**
