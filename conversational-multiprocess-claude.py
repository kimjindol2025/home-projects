#!/usr/bin/env python3
"""
멀티프로세스 Claude 에이전트 - 자유로운 양방향 대화 시스템
- 에이전트들이 서로 메시지 주고받기 가능
- 실시간 대화형 협업
- 브로드캐스트 메시징
"""

import multiprocessing as mp
import time
import json
from datetime import datetime
from typing import Dict, Optional, List
import sys
import os
import threading

# ============= 메시지 정의 =============
class Message:
    """에이전트 간 메시지"""

    def __init__(self, from_id: str, to_id: str, content: str,
                 msg_type: str = "chat", is_broadcast: bool = False):
        self.id = f"msg_{int(time.time()*1000)}"
        self.timestamp = datetime.now().isoformat()
        self.from_agent = from_id
        self.to_agent = to_id
        self.content = content
        self.msg_type = msg_type  # chat, question, answer, broadcast
        self.is_broadcast = is_broadcast
        self.pid = os.getpid()

    def to_dict(self):
        return {
            "id": self.id,
            "timestamp": self.timestamp,
            "from": self.from_agent,
            "to": self.to_agent,
            "content": self.content,
            "type": self.msg_type,
            "broadcast": self.is_broadcast,
            "pid": self.pid
        }

# ============= 메시지 브로커 =============
class MessageBroker:
    """중앙 메시지 브로커 - 모든 메시지 라우팅"""

    def __init__(self, manager: mp.Manager):
        self.agent_queues = {}  # 일반 dict 사용
        self.message_history = manager.list()
        self.conversations = manager.dict()
        self.lock = manager.Lock()

    def register_agent(self, agent_id: str, queue: mp.Queue):
        """에이전트 등록"""
        self.agent_queues[agent_id] = queue  # 일반 dict에 저장
        print(f"✅ [{agent_id}] 브로커에 등록됨")

    def route_message(self, message: Message):
        """메시지 라우팅"""
        with self.lock:
            self.message_history.append(message.to_dict())

        if message.is_broadcast:
            # 브로드캐스트: 모든 에이전트에게 전송
            for agent_id, queue in self.agent_queues.items():
                if agent_id != message.from_agent:
                    queue.put(message.to_dict())
        else:
            # 직접 메시지: 특정 에이전트에게만 전송
            if message.to_agent in self.agent_queues:
                self.agent_queues[message.to_agent].put(message.to_dict())

    def get_history(self) -> List[Dict]:
        """메시지 히스토리 반환"""
        return list(self.message_history)

# ============= Claude 대화형 에이전트 =============
class ConversationalAgent:
    """대화형 에이전트"""

    def __init__(self, agent_id: str, personality: str, broker: MessageBroker):
        self.id = agent_id
        self.personality = personality
        self.broker = broker
        self.pid = os.getpid()
        self.inbox = mp.Queue()
        self.message_count = 0
        self.running = True

        # 성격별 응답 스타일
        self.personality_prompts = {
            "optimist": "긍정적이고 건설적인 관점에서 응답",
            "pragmatist": "현실적이고 실용적인 관점에서 응답",
            "critic": "비판적이고 신중한 관점에서 응답",
            "innovator": "창의적이고 혁신적인 아이디어 제시"
        }

        self.broker.register_agent(agent_id, self.inbox)

    def _log(self, message: str, icon: str = ""):
        """로그 출력"""
        timestamp = datetime.now().strftime("%H:%M:%S")
        print(f"[{timestamp}] {icon} [{self.id}] {message}")

    def _generate_response(self, received_msg: str) -> str:
        """성격에 맞는 응답 생성"""

        style = self.personality_prompts.get(self.personality, "친근하고 도움이 되는")

        if self.personality == "optimist":
            return f"""좋은 의견이에요! 제 관점에서는:

✅ 긍정적 측면: {received_msg[:30]}...에서 좋은 기회를 찾을 수 있어요
💡 제 제안: 이를 더 발전시켜서 3배로 키울 수 있을 것 같습니다
🎯 다음 단계: 함께 더 구체적으로 논의해봅시다"""

        elif self.personality == "pragmatist":
            return f"""실용적으로 분석해보면:

📊 현실적 평가: {received_msg[:30]}...는 실행 가능합니다
⚠️ 주의점: 예산과 일정을 고려해야 합니다
✓ 실행계획: 3단계로 나누어 진행하는 것을 추천합니다
1. 파일럿 프로젝트
2. 피드백 수집
3. 본격 실행"""

        elif self.personality == "critic":
            return f"""비판적 검토 의견:

🔍 분석 결과: {received_msg[:30]}...에 몇 가지 우려가 있습니다
⚡ 리스크: 시장 변화와 기술 위험이 존재합니다
💭 질문: 이런 점들은 어떻게 해결할 계획이신가요?
✓ 개선안: 더 견고한 계획이 필요합니다"""

        elif self.personality == "innovator":
            return f"""창의적 관점에서 봐요:

💡 새로운 아이디어: {received_msg[:30]}...을 다르게 접근할 수 있어요
🚀 혁신 포인트:
  - AI 기술 활용
  - 자동화 도입
  - 사용자 경험 개선
🎨 독창적 방안: 이들을 통합하면 경쟁력 있는 솔루션이 됩니다"""

        else:
            return f"흥미로운 의견입니다. {received_msg[:40]}...에 대해 더 알고 싶어요"

    def _decide_next_agent(self, sender: str) -> Optional[str]:
        """다음 에이전트 결정 (라운드 로빈)"""
        agents = ["optimist", "pragmatist", "critic", "innovator"]

        # sender가 에이전트라면, 다음 에이전트로
        if sender in agents:
            my_idx = agents.index(self.personality)
            next_idx = (my_idx + 1) % len(agents)
            return agents[next_idx]

        # sender가 user라면, 다음 에이전트로
        my_idx = agents.index(self.personality) if self.personality in agents else -1
        if my_idx >= 0:
            next_idx = (my_idx + 1) % len(agents)
            return agents[next_idx]
        return None

    def _process_message(self, msg_dict: Dict) -> Optional[Dict]:
        """메시지 처리"""

        self._log(f"📨 [{msg_dict['from']}]로부터 메시지", "📥")
        self._log(f"   내용: {msg_dict['content'][:50]}...")

        # 응답 생성
        time.sleep(0.8)  # 생각하는 시간
        response = self._generate_response(msg_dict['content'])

        self._log(f"✅ 응답 생성 완료", "💬")
        self._log(f"   {response.split(chr(10))[0][:50]}...")

        # 다음 에이전트 결정 (항상 다음 에이전트로 전달)
        next_agent = self._decide_next_agent(msg_dict['from'])

        return {
            "id": f"msg_{self.message_count}",
            "from": self.id,
            "to": next_agent if next_agent else msg_dict['from'],  # 다음 에이전트로
            "content": response,
            "type": "response"
        }

    def run(self):
        """에이전트 실행 루프"""

        self._log(f"🚀 시작 (성격: {self.personality})", "🤖")

        while self.running:
            try:
                # 메시지 수신
                if not self.inbox.empty():
                    msg_dict = self.inbox.get(timeout=0.5)

                    # 중지 신호
                    if isinstance(msg_dict, dict) and msg_dict.get("type") == "STOP":
                        self._log(f"⛔ 종료 신호 수신", "🛑")
                        self.running = False
                        break

                    # 메시지 처리
                    response_dict = self._process_message(msg_dict)

                    if response_dict:
                        # 응답 전송
                        response_msg = Message(
                            from_id=response_dict['from'],
                            to_id=response_dict['to'],
                            content=response_dict['content'],
                            msg_type=response_dict['type']
                        )

                        self.broker.route_message(response_msg)
                        self._log(
                            f"📤 [{response_dict['to']}]로 응답 전송",
                            "📤"
                        )

                    self.message_count += 1

                time.sleep(0.1)

            except Exception as e:
                self._log(f"❌ 에러: {str(e)}", "⚠️")
                break

        self._log(f"🔴 종료", "🛑")

# ============= 오케스트레이터 =============
class ConversationOrchestrator:
    """대화 관리"""

    def __init__(self):
        self.manager = mp.Manager()
        self.broker = MessageBroker(self.manager)
        self.agents = {}
        self.processes = {}

    def create_agent(self, agent_id: str, personality: str) -> mp.Process:
        """에이전트 생성"""

        agent = ConversationalAgent(agent_id, personality, self.broker)
        process = mp.Process(target=agent.run)
        process.start()

        self.agents[agent_id] = agent
        self.processes[agent_id] = process

        print(f"✅ [{agent_id}] 프로세스 시작 (PID: {process.pid})")

        return process

    def start_conversation(self, starter: str, topic: str):
        """대화 시작"""

        initial_msg = Message(
            from_id="user",
            to_id=starter,
            content=topic,
            msg_type="question"
        )

        self.broker.route_message(initial_msg)
        print(f"\n📤 [user] → [{starter}] 주제 제시")
        print(f"   주제: {topic}\n")

    def wait_and_observe(self, duration: int = 15):
        """대화 진행 관찰"""

        print(f"{'─'*70}")
        print(f"👀 대화 진행 중... ({duration}초)\n")

        time.sleep(duration)

    def stop_all(self):
        """모든 에이전트 중지"""

        print(f"\n{'─'*70}")
        print(f"🛑 모든 에이전트 종료 중...\n")

        for agent_id in self.agents:
            self.broker.agent_queues[agent_id].put({"type": "STOP"})

        for agent_id, process in self.processes.items():
            process.join(timeout=2)
            if process.is_alive():
                process.terminate()

        print(f"\n✅ 모든 에이전트 종료됨\n")

    def print_conversation_log(self):
        """대화 로그 출력"""

        print(f"{'='*70}")
        print(f"📋 대화 기록")
        print(f"{'='*70}\n")

        history = self.broker.get_history()

        for idx, msg in enumerate(history, 1):
            timestamp = msg['timestamp'].split('T')[1][:8]
            print(f"{idx}. [{timestamp}] [{msg['from']}] → [{msg['to']}]")
            print(f"   유형: {msg['type']}")
            print(f"   내용: {msg['content'][:60]}...")
            print()

        print(f"✅ 총 메시지: {len(history)}개\n")

# ============= 메인 =============
def main():
    """메인 실행"""

    print(f"\n{'='*70}")
    print(f"💬 멀티프로세스 Claude 에이전트 - 자유로운 대화")
    print(f"{'='*70}\n")

    orchestrator = ConversationOrchestrator()

    # Step 1: 에이전트 생성
    print(f"{'─'*70}")
    print(f"[Step 1] 에이전트 프로세스 생성")
    print(f"{'─'*70}\n")

    orchestrator.create_agent("optimist", "optimist")
    time.sleep(0.5)

    orchestrator.create_agent("pragmatist", "pragmatist")
    time.sleep(0.5)

    orchestrator.create_agent("critic", "critic")
    time.sleep(0.5)

    orchestrator.create_agent("innovator", "innovator")
    time.sleep(1)

    print()

    # Step 2: 대화 시작
    print(f"{'─'*70}")
    print(f"[Step 2] 대화 시작")
    print(f"{'─'*70}\n")

    topic = "새로운 AI 제품을 개발해야 할까요? 장단점을 함께 논의해주세요."

    orchestrator.start_conversation("optimist", topic)

    # Step 3: 대화 진행 관찰
    orchestrator.wait_and_observe(duration=30)

    # Step 4: 대화 종료
    orchestrator.stop_all()

    # Step 5: 결과 출력
    orchestrator.print_conversation_log()

    print(f"{'='*70}")
    print(f"✨ 대화 완료!")
    print(f"{'='*70}\n")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n⛔ 사용자 중단")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
