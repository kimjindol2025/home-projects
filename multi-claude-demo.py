#!/usr/bin/env python3
"""
서버 내 멀티 Claude 에이전트 통신 시스템 (데모)
시각화된 IPC 메시지 버스 예제
"""

import asyncio
from dataclasses import dataclass
from datetime import datetime
from typing import Optional

# ============= 메시지 정의 =============
@dataclass
class Message:
    timestamp: str
    from_agent: str
    to_agent: str
    content: str
    message_id: str

# ============= 메시지 버스 (IPC) =============
class MessageBus:
    """에이전트 간 메시지 교환을 위한 버스"""

    def __init__(self):
        self.agents = {}
        self.message_log = []
        self.inbox = {}

    def register(self, agent_id: str, agent):
        """에이전트를 버스에 등록"""
        self.agents[agent_id] = agent
        self.inbox[agent_id] = []
        print(f"✅ 에이전트 등록: {agent_id}")

    async def send(self, from_id: str, to_id: str, content: str) -> Message:
        """메시지 전송"""
        msg = Message(
            timestamp=datetime.now().isoformat(),
            from_agent=from_id,
            to_agent=to_id,
            content=content,
            message_id=f"msg_{len(self.message_log)}"
        )

        self.message_log.append(msg)
        self.inbox[to_id].append(msg)

        print(f"\n📨 [{from_id}] → [{to_id}]")
        print(f"   메시지: {content[:100]}...")
        await asyncio.sleep(0.5)

        return msg

    def get_inbox(self, agent_id: str):
        """에이전트의 받은 메시지 반환"""
        messages = self.inbox[agent_id]
        self.inbox[agent_id] = []
        return messages

    def get_log(self):
        """모든 메시지 로그 반환"""
        return self.message_log

# ============= Claude 에이전트 (시뮬레이션) =============
class ClaudeAgent:
    """Claude API를 사용하는 에이전트"""

    def __init__(self, agent_id: str, role: str, bus: MessageBus):
        self.id = agent_id
        self.role = role
        self.bus = bus
        self.message_history = []

        # 역할별 응답 템플릿
        self.responses = {
            "researcher": {
                "initial": """AI 기술은 빠르게 발전하고 있습니다. 특히 생성형 AI의 발전으로 많은 산업에서
혁신이 일어나고 있습니다. 다만 규제와 윤리 문제도 동시에 고려해야 합니다.""",
                "to_critic": """감사합니다. 우려 사항을 고려하여 책임 있는 AI 개발에 집중해야 한다는 점에 동의합니다.
다만 혁신의 속도도 중요하다고 생각합니다.""",
                "debate": """원격 근무는 많은 이점이 있지만, 팀 협업과 회사 문화 형성에 어려움이 있을 수 있습니다.
하이브리드 모델이 최선의 접근일 수 있습니다."""
            },
            "critic": {
                "response": """흥미로운 관점이지만 몇 가지 우려사항이 있습니다. 생성형 AI의 환경 영향과
개인정보 보호 측면에서 더 깊이 있는 논의가 필요합니다.""",
                "debate": """당신의 의견도 일리가 있지만, 실제 상황은 훨씬 복잡합니다.
신뢰 구축과 온보딩 문제, 회사 간 소통 등이 문제가 될 수 있습니다."""
            },
            "synthesizer": {
                "conclusion": """종합하면, AI의 미래는 밝지만 신중한 접근이 필요합니다.
혁신과 책임성의 균형을 맞추는 것이 핵심입니다."""
            }
        }

        self.bus.register(agent_id, self)

    async def send_message(self, target_id: str, content: str):
        """메시지 전송"""
        await self.bus.send(self.id, target_id, content)

    async def process_inbox(self) -> Optional[str]:
        """받은 메시지 처리 및 응답 생성"""
        messages = self.bus.get_inbox(self.id)

        if not messages:
            print(f"[{self.id}] 받은 메시지 없음")
            return None

        message = messages[0]
        print(f"\n🔄 [{self.id}] 메시지 처리 중...")
        await asyncio.sleep(1)

        # 역할에 따른 응답 반환
        if self.role == "critic":
            reply = self.responses["critic"]["response"]
        elif self.role == "synthesizer":
            reply = self.responses["synthesizer"]["conclusion"]
        else:
            reply = self.responses.get(self.role, {}).get("to_critic", "알겠습니다.")

        self.message_history.append({
            "received_from": message.from_agent,
            "received_content": message.content[:100],
            "response": reply
        })

        return reply

    def get_history(self):
        """상호작용 히스토리 반환"""
        return self.message_history

# ============= 협업 워크플로우 =============
async def collaborative_workflow():
    print(f"\n{'='*70}")
    print(f"🤝 서버 내 멀티 Claude 협업 시스템")
    print(f"{'='*70}\n")

    bus = MessageBus()

    # 3개 에이전트 생성
    researcher = ClaudeAgent("researcher", "researcher", bus)
    critic = ClaudeAgent("critic", "critic", bus)
    synthesizer = ClaudeAgent("synthesizer", "synthesizer", bus)

    topic = "AI의 미래는 밝을까?"
    print(f"📌 주제: {topic}\n")

    # Step 1: 연구원의 초기 관점
    print(f"{'─'*70}")
    print(f"[Step 1] 연구원의 초기 관점")
    print(f"{'─'*70}\n")

    initial_view = """AI 기술은 빠르게 발전하고 있습니다. 특히 생성형 AI의 발전으로 많은 산업에서
혁신이 일어나고 있습니다. 다만 규제와 윤리 문제도 동시에 고려해야 합니다."""

    print(f"🔵 Researcher:\n{initial_view}\n")
    researcher.message_history.append({
        "received_from": None,
        "received_content": None,
        "response": initial_view
    })

    # Step 2: 비평가의 검토
    print(f"{'─'*70}")
    print(f"[Step 2] 비평가의 검토")
    print(f"{'─'*70}")

    await researcher.send_message("critic", initial_view)
    critic_response = await critic.process_inbox()
    print(f"\n🔴 Critic:\n{critic_response}\n")

    # Step 3: 종합가의 최종 결론
    print(f"{'─'*70}")
    print(f"[Step 3] 종합가의 최종 결론")
    print(f"{'─'*70}")

    combined_input = f"""
초기 의견:
"{initial_view}"

비평:
"{critic_response}"

이 두 관점을 종합하여 균형잡힌 결론을 제시해주세요."""

    await researcher.send_message("synthesizer", combined_input)
    final_conclusion = await synthesizer.process_inbox()
    print(f"\n🟢 Synthesizer:\n{final_conclusion}\n")

    # Step 4: 결과 정리
    print(f"{'='*70}")
    print(f"📊 협업 결과 요약")
    print(f"{'='*70}\n")

    print(f"✅ 전송된 메시지: {len(bus.get_log())}개")
    print(f"✅ 활성 에이전트: 3개")
    print(f"✅ 완료된 상호작용: 2개\n")

    print(f"메시지 로그:")
    for idx, log in enumerate(bus.get_log(), 1):
        print(f"  {idx}. [{log.from_agent}] → [{log.to_agent}]")

    return bus, [researcher, critic, synthesizer]

# ============= 실시간 토론 =============
async def realtime_debate():
    print(f"\n{'='*70}")
    print(f"💬 실시간 멀티 Claude 토론")
    print(f"{'='*70}\n")

    bus = MessageBus()
    agent_a = ClaudeAgent("agentA", "researcher", bus)
    agent_b = ClaudeAgent("agentB", "critic", bus)

    # Agent A의 주장
    print(f"[Round 1] Agent A의 주장\n")
    argument_a = """원격 근무는 통근 시간을 절약하고, 개발자들이 더 집중할 수 있는 환경을
제공합니다. 실제 많은 기업들이 원격 근무 도입 후 생산성 향상을 보고했습니다."""
    print(f"🟦 Agent A:\n{argument_a}\n")
    agent_a.message_history.append({
        "received_from": None,
        "received_content": None,
        "response": argument_a
    })

    # Agent B의 반박
    await agent_a.send_message("agentB", argument_a)
    print(f"[Round 2] Agent B의 반박\n")
    rebuttal = await agent_b.process_inbox()
    print(f"🟥 Agent B:\n{rebuttal}\n")

    # Agent A의 재반박
    await agent_b.send_message("agentA", rebuttal)
    print(f"[Round 3] Agent A의 재반박\n")
    counter = await agent_a.process_inbox()
    print(f"🟦 Agent A:\n{counter}\n")

    print(f"✅ 토론 완료! (총 3라운드)\n")

    return bus, [agent_a, agent_b]

# ============= 메인 =============
async def main():
    try:
        # 협업 워크플로우
        bus1, agents1 = await collaborative_workflow()

        # 실시간 토론
        bus2, agents2 = await realtime_debate()

        print(f"\n{'='*70}")
        print(f"✨ 모든 에이전트 상호작용 완료")
        print(f"{'='*70}\n")

        print("🎯 주요 기능:")
        print("  1. 메시지 버스 (MessageBus) - IPC 채널")
        print("  2. 3개 이상의 Claude 에이전트 동시 실행")
        print("  3. 비동기 메시지 교환 (async/await)")
        print("  4. 메시지 로그 추적")
        print("  5. 역할별 응답 생성\n")

        print("💡 다음 단계:")
        print("  - 실제 Claude API 연동")
        print("  - 웹소켓 기반 실시간 통신")
        print("  - 복잡한 멀티 에이전트 협업")
        print("  - 상태 저장소 추가\n")

    except Exception as e:
        print(f"❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    asyncio.run(main())
