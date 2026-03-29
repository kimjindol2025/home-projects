#!/usr/bin/env python3
"""
서버 내 멀티 Claude 에이전트 통신 시스템
같은 서버에서 여러 Claude 인스턴스가 메시지 교환
"""

import asyncio
import anthropic
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
        print(f"   메시지: {content[:80]}...")

        return msg

    def get_inbox(self, agent_id: str):
        """에이전트의 받은 메시지 반환"""
        messages = self.inbox[agent_id]
        self.inbox[agent_id] = []
        return messages

    def get_log(self):
        """모든 메시지 로그 반환"""
        return self.message_log

# ============= Claude 에이전트 =============
class ClaudeAgent:
    def __init__(self, agent_id: str, role: str, bus: MessageBus):
        self.id = agent_id
        self.role = role
        self.bus = bus
        import os
        self.client = anthropic.Client(api_key=os.environ.get("ANTHROPIC_API_KEY"))
        self.message_history = []

        self.system_prompts = {
            "researcher": """당신은 연구원입니다. 주어진 주제에 대해 깊이 있게 분석합니다.
다른 에이전트의 피드백을 수집하고 종합합니다.""",
            "critic": """당신은 비평가입니다. 제시된 아이디어를 비판적으로 검토합니다.
논리적 결함과 개선점을 지적합니다.""",
            "synthesizer": """당신은 종합가입니다. 여러 의견을 통합하여 최종 결론을 도출합니다.
모든 관점을 존중하면서 균형잡힌 결론을 만듭니다."""
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

        system_prompt = self.system_prompts.get(
            self.role, self.system_prompts["researcher"]
        )

        response = self.client.messages.create(
            model="claude-opus-4-6",
            max_tokens=400,
            system=system_prompt,
            messages=[
                {
                    "role": "user",
                    "content": f'다른 에이전트({message.from_agent})로부터 다음 메시지를 받았습니다:\n\n"{message.content}"\n\n이에 대해 당신의 관점을 간단히 (2-3문장) 제시해주세요.'
                }
            ]
        )

        reply = response.content[0].text
        self.message_history.append({
            "received_from": message.from_agent,
            "received_content": message.content,
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
    print(f"\n📌 주제: {topic}\n")

    # Step 1: 연구원의 초기 관점
    print(f"\n{'─'*70}")
    print(f"[Step 1] 연구원의 초기 관점")
    print(f"{'─'*70}")

    initial_view = """AI 기술은 빠르게 발전하고 있습니다. 특히 생성형 AI의 발전으로 많은 산업에서
혁신이 일어나고 있습니다. 다만 규제와 윤리 문제도 동시에 고려해야 합니다."""

    print(f"\n🔵 Researcher:\n{initial_view}\n")
    researcher.message_history.append({
        "received_from": None,
        "received_content": None,
        "response": initial_view
    })

    # Step 2: 비평가의 검토
    print(f"\n{'─'*70}")
    print(f"[Step 2] 비평가의 검토")
    print(f"{'─'*70}")

    await researcher.send_message("critic", initial_view)
    await asyncio.sleep(1)

    critic_response = await critic.process_inbox()
    print(f"\n🔴 Critic:\n{critic_response}\n")

    # Step 3: 종합가의 최종 결론
    print(f"\n{'─'*70}")
    print(f"[Step 3] 종합가의 최종 결론")
    print(f"{'─'*70}")

    combined_input = f"""
초기 의견:
"{initial_view}"

비평:
"{critic_response}"

이 두 관점을 종합하여 균형잡힌 결론을 제시해주세요."""

    await researcher.send_message("synthesizer", combined_input)
    await asyncio.sleep(1)

    final_conclusion = await synthesizer.process_inbox()
    print(f"\n🟢 Synthesizer:\n{final_conclusion}\n")

    # Step 4: 결과 정리
    print(f"\n{'='*70}")
    print(f"📊 협업 결과 요약")
    print(f"{'='*70}\n")

    print(f"✅ 전송된 메시지: {len(bus.get_log())}개")
    print(f"✅ 활성 에이전트: 3개")
    print(f"✅ 완료된 상호작용: 2개\n")

    print(f"메시지 로그:")
    for idx, log in enumerate(bus.get_log(), 1):
        print(f"  {idx}. [{log.from_agent}] → [{log.to_agent}]")

# ============= 실시간 토론 =============
async def realtime_debate():
    print(f"\n{'='*70}")
    print(f"💬 실시간 멀티 Claude 토론")
    print(f"{'='*70}\n")

    bus = MessageBus()
    agent_a = ClaudeAgent("agentA", "researcher", bus)
    agent_b = ClaudeAgent("agentB", "critic", bus)

    topic = "원격 근무는 생산성을 높이는가?"

    # Agent A의 주장
    print(f"\n[Round 1] Agent A의 주장")
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
    await asyncio.sleep(1)

    print(f"\n[Round 2] Agent B의 반박")
    rebuttal = await agent_b.process_inbox()
    print(f"🟥 Agent B:\n{rebuttal}\n")

    # Agent A의 재반박
    await agent_b.send_message("agentA", rebuttal)
    await asyncio.sleep(1)

    print(f"\n[Round 3] Agent A의 재반박")
    counter = await agent_a.process_inbox()
    print(f"🟦 Agent A:\n{counter}\n")

    print(f"\n✅ 토론 완료! (총 3라운드)\n")

# ============= 메인 =============
async def main():
    try:
        # 협업 워크플로우
        await collaborative_workflow()

        # 실시간 토론
        await realtime_debate()

        print(f"\n{'='*70}")
        print(f"✨ 모든 에이전트 상호작용 완료")
        print(f"{'='*70}\n")

    except Exception as e:
        print(f"❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    asyncio.run(main())
