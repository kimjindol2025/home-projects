#!/usr/bin/env python3
"""
API 키 없이 서버 내 에이전트들이 서로 대화하는 시스템
- 로컬 메모리 기반 상태 관리
- 규칙 기반 응답 생성
- 에이전트 간 지식 공유
"""

import asyncio
import json
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, List, Dict
from enum import Enum

# ============= 상태 정의 =============
class AgentRole(Enum):
    """에이전트 역할"""
    RESEARCHER = "researcher"    # 연구자
    CRITIC = "critic"           # 비평가
    COORDINATOR = "coordinator" # 조정자
    EXECUTOR = "executor"       # 실행자

# ============= 메시지 정의 =============
@dataclass
class Message:
    """에이전트 간 메시지"""
    id: str
    timestamp: str
    from_agent: str
    to_agent: str
    content: str
    message_type: str = "text"  # text, query, response, command
    metadata: Dict = field(default_factory=dict)

    def to_dict(self):
        return {
            "id": self.id,
            "timestamp": self.timestamp,
            "from": self.from_agent,
            "to": self.to_agent,
            "type": self.message_type,
            "content": self.content,
            "metadata": self.metadata
        }

# ============= 지식 저장소 =============
class KnowledgeStore:
    """모든 에이전트가 공유하는 지식 저장소"""

    def __init__(self):
        self.facts = {}        # 사실들
        self.opinions = {}     # 의견들
        self.conclusions = {}  # 결론들
        self.context = {}      # 현재 컨텍스트

    def store_fact(self, key: str, value: str, agent_id: str):
        """사실 저장"""
        self.facts[key] = {
            "value": value,
            "source": agent_id,
            "timestamp": datetime.now().isoformat()
        }
        print(f"💾 [{agent_id}] 사실 저장: {key}")

    def store_opinion(self, key: str, value: str, agent_id: str):
        """의견 저장"""
        self.opinions[key] = {
            "value": value,
            "source": agent_id,
            "timestamp": datetime.now().isoformat()
        }
        print(f"💭 [{agent_id}] 의견 저장: {key}")

    def store_conclusion(self, key: str, value: str, agent_id: str):
        """결론 저장"""
        self.conclusions[key] = {
            "value": value,
            "source": agent_id,
            "timestamp": datetime.now().isoformat()
        }
        print(f"✅ [{agent_id}] 결론 저장: {key}")

    def get_facts(self) -> Dict:
        """모든 사실 조회"""
        return self.facts

    def get_opinions(self) -> Dict:
        """모든 의견 조회"""
        return self.opinions

    def get_conclusions(self) -> Dict:
        """모든 결론 조회"""
        return self.conclusions

    def set_context(self, key: str, value: str):
        """컨텍스트 설정"""
        self.context[key] = value

    def get_context(self) -> Dict:
        """컨텍스트 조회"""
        return self.context

# ============= 메시지 버스 =============
class MessageBus:
    """에이전트 간 메시지 교환 버스"""

    def __init__(self, knowledge_store: KnowledgeStore):
        self.agents = {}
        self.message_log = []
        self.inbox = {}
        self.knowledge_store = knowledge_store

    def register(self, agent_id: str, agent):
        """에이전트 등록"""
        self.agents[agent_id] = agent
        self.inbox[agent_id] = []
        print(f"✅ 에이전트 등록: {agent_id}")

    async def send(self, from_id: str, to_id: str, content: str,
                   msg_type: str = "text", metadata: Dict = None) -> Message:
        """메시지 전송"""
        msg = Message(
            id=f"msg_{len(self.message_log)}",
            timestamp=datetime.now().isoformat(),
            from_agent=from_id,
            to_agent=to_id,
            content=content,
            message_type=msg_type,
            metadata=metadata or {}
        )

        self.message_log.append(msg)
        self.inbox[to_id].append(msg)

        print(f"\n📨 [{from_id}] → [{to_id}] ({msg_type})")
        print(f"   {content[:80]}...")
        await asyncio.sleep(0.3)

        return msg

    def get_inbox(self, agent_id: str) -> List[Message]:
        """받은 메시지 반환"""
        messages = self.inbox[agent_id]
        self.inbox[agent_id] = []
        return messages

    def get_log(self) -> List[Message]:
        """메시지 로그 반환"""
        return self.message_log

# ============= 로컬 Claude 에이전트 =============
class LocalClaudeAgent:
    """API 없이 로컬에서 작동하는 Claude 에이전트"""

    def __init__(self, agent_id: str, role: AgentRole,
                 bus: MessageBus, knowledge_store: KnowledgeStore):
        self.id = agent_id
        self.role = role
        self.bus = bus
        self.knowledge_store = knowledge_store
        self.message_history = []
        self.decision_log = []

        self.bus.register(agent_id, self)

    async def send_message(self, target_id: str, content: str,
                          msg_type: str = "text", metadata: Dict = None):
        """메시지 전송"""
        await self.bus.send(self.id, target_id, content, msg_type, metadata)

    async def process_inbox(self) -> Optional[str]:
        """받은 메시지 처리"""
        messages = self.bus.get_inbox(self.id)

        if not messages:
            return None

        message = messages[0]
        print(f"\n🔄 [{self.id}] 메시지 처리 중...")
        await asyncio.sleep(0.5)

        # 역할에 따른 응답 생성
        response = await self._generate_response(message)

        self.message_history.append({
            "received_from": message.from_agent,
            "received_content": message.content[:100],
            "generated_response": response
        })

        return response

    async def _generate_response(self, message: Message) -> str:
        """역할 기반 응답 생성"""

        if self.role == AgentRole.RESEARCHER:
            # 연구자: 분석 및 사실 제시
            response = self._researcher_response(message.content)

        elif self.role == AgentRole.CRITIC:
            # 비평가: 비판적 검토
            response = self._critic_response(message.content)

        elif self.role == AgentRole.COORDINATOR:
            # 조정자: 의견 종합
            response = self._coordinator_response(message.content)

        elif self.role == AgentRole.EXECUTOR:
            # 실행자: 액션 수행
            response = self._executor_response(message.content)

        else:
            response = "알겠습니다."

        # 지식 저장소에 학습
        if message.message_type == "text":
            self.knowledge_store.store_opinion(
                f"{self.id}_response_{len(self.message_history)}",
                response,
                self.id
            )

        self.decision_log.append({
            "timestamp": datetime.now().isoformat(),
            "input": message.content[:50],
            "output": response[:50],
            "role": self.role.value
        })

        return response

    def _researcher_response(self, content: str) -> str:
        """연구자 응답"""
        keywords = ["AI", "발전", "혁신", "기술"]

        if any(k in content for k in keywords):
            return f"""흥미로운 질문이네요. 제 분석으로는:

1. 기술 발전: {content[:40]}...에 대해 더 깊이 있게 살펴볼 필요가 있습니다.
2. 증거: 이를 뒷받침하는 여러 연구 논문들이 있습니다.
3. 시사점: 이는 향후 정책 수립에 영향을 미칠 것입니다."""
        else:
            return f"""주제에 대해 분석해보니, 다음과 같은 관점을 제시합니다:

- 긍정적 측면: 효율성 증대
- 부정적 측면: 관리 복잡성
- 권장사항: 균형잡힌 접근"""

    def _critic_response(self, content: str) -> str:
        """비평가 응답"""
        return f"""좋은 지적이지만, 다음 점들을 고려해야 합니다:

1. 논리적 결함: 일반화의 오류가 있을 수 있습니다.
2. 반례: 실제로는 이런 경우들이 존재합니다.
3. 개선안: 더 엄밀한 정의와 범위 설정이 필요합니다."""

    def _coordinator_response(self, content: str) -> str:
        """조정자 응답"""
        # 지식 저장소에서 기존 의견들 확인
        opinions = self.knowledge_store.get_opinions()
        opinion_count = len(opinions)

        return f"""여러 관점을 종합한 결과:

1. 수집된 의견: {opinion_count}개
2. 공통점: 변화의 필요성
3. 차이점: 방법론과 속도에 대한 의견 차이
4. 최종 합의: 단계적이고 신중한 접근이 필요합니다.

추천 방향: {content[:30]}...을 기반으로 한 하이브리드 전략"""

    def _executor_response(self, content: str) -> str:
        """실행자 응답"""
        return f"""지시사항을 이해했습니다. 다음을 실행하겠습니다:

1. 준비 단계: 환경 설정 중...
2. 실행 단계: {content[:40]}...을(를) 처리 중...
3. 검증 단계: 결과 확인 중...
4. 완료: 작업이 성공적으로 완료되었습니다."""

    def get_decision_log(self) -> List[Dict]:
        """의사결정 로그 반환"""
        return self.decision_log

# ============= 협업 시나리오 =============
async def scenario_product_development():
    """시나리오: 새로운 제품 개발 협업"""

    print(f"\n{'='*70}")
    print(f"🏢 시나리오: 새로운 AI 제품 개발")
    print(f"{'='*70}\n")

    # 초기화
    knowledge_store = KnowledgeStore()
    bus = MessageBus(knowledge_store)
    knowledge_store.set_context("project", "AI-기반 고객 서비스 챗봇")

    # 에이전트 생성
    researcher = LocalClaudeAgent("researcher", AgentRole.RESEARCHER, bus, knowledge_store)
    critic = LocalClaudeAgent("critic", AgentRole.CRITIC, bus, knowledge_store)
    coordinator = LocalClaudeAgent("coordinator", AgentRole.COORDINATOR, bus, knowledge_store)

    project = knowledge_store.get_context()["project"]
    print(f"📌 프로젝트: {project}\n")

    # Phase 1: 연구자의 제안
    print(f"{'─'*70}")
    print(f"[Phase 1] 연구자의 기술 제안")
    print(f"{'─'*70}\n")

    proposal = f"""우리의 {project} 개발 제안:

1. 아키텍처: 클라우드 기반 마이크로서비스
2. 기술: LLM + 규칙 엔진 하이브리드
3. 확장성: 초당 1000개 요청 처리 가능
4. 비용: 연 500만원 예상"""

    print(f"🔵 Researcher:\n{proposal}\n")
    researcher.message_history.append({
        "received_from": None,
        "generated_response": proposal
    })

    knowledge_store.store_fact("technical_proposal", proposal, "researcher")

    # Phase 2: 비평가의 검토
    print(f"{'─'*70}")
    print(f"[Phase 2] 비평가의 위험 평가")
    print(f"{'─'*70}\n")

    await researcher.send_message("critic", proposal, "query")
    critic_review = await critic.process_inbox()
    print(f"🔴 Critic:\n{critic_review}\n")

    knowledge_store.store_opinion("risk_assessment", critic_review, "critic")

    # Phase 3: 조정자의 최종 결정
    print(f"{'─'*70}")
    print(f"[Phase 3] 조정자의 통합 결정")
    print(f"{'─'*70}\n")

    combined_input = f"제안: {proposal[:100]}...\n검토: {critic_review[:100]}..."

    await critic.send_message("coordinator", combined_input, "synthesis")
    final_decision = await coordinator.process_inbox()
    print(f"🟢 Coordinator:\n{final_decision}\n")

    knowledge_store.store_conclusion("project_decision", final_decision, "coordinator")

    # 결과 정리
    print(f"\n{'='*70}")
    print(f"📊 협업 결과")
    print(f"{'='*70}\n")

    print(f"✅ 전송된 메시지: {len(bus.get_log())}개")
    print(f"✅ 저장된 사실: {len(knowledge_store.get_facts())}개")
    print(f"✅ 저장된 의견: {len(knowledge_store.get_opinions())}개")
    print(f"✅ 저장된 결론: {len(knowledge_store.get_conclusions())}개\n")

    print("📋 메시지 로그:")
    for idx, msg in enumerate(bus.get_log(), 1):
        print(f"  {idx}. [{msg.from_agent}] → [{msg.to_agent}] ({msg.message_type})")

    return bus, knowledge_store

# ============= 실시간 협의 시나리오 =============
async def scenario_meeting():
    """시나리오: 실시간 팀 회의"""

    print(f"\n\n{'='*70}")
    print(f"💼 시나리오: 실시간 팀 회의")
    print(f"{'='*70}\n")

    knowledge_store = KnowledgeStore()
    bus = MessageBus(knowledge_store)

    # 에이전트들
    pm = LocalClaudeAgent("pm", AgentRole.COORDINATOR, bus, knowledge_store)
    engineer = LocalClaudeAgent("engineer", AgentRole.RESEARCHER, bus, knowledge_store)
    qa = LocalClaudeAgent("qa", AgentRole.CRITIC, bus, knowledge_store)

    topic = "마감일 단축 가능성 검토"
    print(f"📌 회의 주제: {topic}\n")

    # PM의 제안
    print(f"[PM 발언]")
    pm_proposal = """2주 단축이 가능할까요? 클라이언트가 요청했습니다."""
    print(f"💼 PM: {pm_proposal}\n")
    pm.message_history.append({
        "received_from": None,
        "generated_response": pm_proposal
    })

    # Engineer의 기술적 검토
    await pm.send_message("engineer", pm_proposal, "question")
    await asyncio.sleep(0.5)
    print(f"[Engineer 발언]")
    engineer_response = await engineer.process_inbox()
    print(f"🔧 Engineer:\n{engineer_response}\n")

    # QA의 품질 우려
    await engineer.send_message("qa", engineer_response, "response")
    await asyncio.sleep(0.5)
    print(f"[QA 발언]")
    qa_response = await qa.process_inbox()
    print(f"🧪 QA:\n{qa_response}\n")

    # PM의 최종 결정
    await qa.send_message("pm", qa_response, "concern")
    await asyncio.sleep(0.5)
    print(f"[PM 최종 결정]")
    final_decision = await pm.process_inbox()
    print(f"💼 PM:\n{final_decision}\n")

    print(f"✅ 회의 완료!\n")

    return bus, knowledge_store

# ============= 메인 =============
async def main():
    try:
        # 시나리오 1: 제품 개발
        bus1, ks1 = await scenario_product_development()

        # 시나리오 2: 실시간 회의
        bus2, ks2 = await scenario_meeting()

        # 최종 요약
        print(f"\n{'='*70}")
        print(f"✨ 모든 협업 완료")
        print(f"{'='*70}\n")

        print("🎯 시스템 특징:")
        print("  ✓ API 키 불필요 (로컬 작동)")
        print("  ✓ 에이전트 간 지식 공유 (KnowledgeStore)")
        print("  ✓ 역할 기반 응답 생성")
        print("  ✓ 완전한 메시지 추적")
        print("  ✓ 비동기 협업\n")

        print("🔧 확장 가능 영역:")
        print("  - 실제 LLM 통합 (ollama, llama.cpp)")
        print("  - 외부 시스템 연동 (API, DB)")
        print("  - 의사결정 엔진 추가")
        print("  - 영구 저장소 (파일/DB)")
        print("  - REST API 서버화\n")

    except Exception as e:
        print(f"❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    asyncio.run(main())
