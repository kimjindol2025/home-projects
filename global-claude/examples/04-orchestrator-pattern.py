#!/usr/bin/env python3
"""
서버 내 다중 Claude 프로세스 시스템
- 각 에이전트가 독립적인 프로세스로 실행
- 프로세스 간 통신 (IPC)
- 병렬 처리 및 협업
"""

import multiprocessing
import asyncio
import time
import json
from dataclasses import dataclass, asdict
from datetime import datetime
from typing import Optional, Dict, List
from enum import Enum
from queue import Queue
import os
import sys

# ============= 메시지 정의 =============
@dataclass
class Message:
    """프로세스 간 메시지"""
    id: str
    timestamp: str
    from_agent: str
    to_agent: str
    content: str
    process_id: int

    def to_dict(self):
        return asdict(self)

# ============= 에이전트 역할 =============
class AgentRole(Enum):
    ANALYST = "analyst"       # 분석가
    VALIDATOR = "validator"   # 검증자
    OPTIMIZER = "optimizer"   # 최적화자

# ============= Claude 에이전트 프로세스 =============
class ClaudeAgentProcess:
    """독립적인 프로세스로 실행되는 Claude 에이전트"""

    def __init__(self, agent_id: str, role: AgentRole,
                 inbox: multiprocessing.Queue,
                 outbox: multiprocessing.Queue,
                 broadcast: multiprocessing.Queue):
        self.agent_id = agent_id
        self.role = role
        self.inbox = inbox
        self.outbox = outbox
        self.broadcast = broadcast
        self.process_id = os.getpid()
        self.message_count = 0
        self.responses = []

        print(f"✅ [{self.agent_id}] 프로세스 시작 (PID: {self.process_id})")

    def run(self):
        """에이전트 실행 루프"""
        print(f"🔄 [{self.agent_id}] 이벤트 루프 시작...\n")

        while True:
            try:
                # 받은 메시지 처리
                if not self.inbox.empty():
                    message = self.inbox.get(timeout=1)

                    if isinstance(message, dict) and message.get("type") == "STOP":
                        print(f"⛔ [{self.agent_id}] 종료 신호 수신")
                        break

                    await_response = self._process_message(message)

                    # 응답 전송
                    if await_response:
                        response_msg = Message(
                            id=f"msg_{self.message_count}",
                            timestamp=datetime.now().isoformat(),
                            from_agent=self.agent_id,
                            to_agent=await_response["to"],
                            content=await_response["content"],
                            process_id=self.process_id
                        )

                        self.outbox.put(response_msg.to_dict())
                        self.broadcast.put({
                            "type": "message",
                            "data": response_msg.to_dict()
                        })

                        self.message_count += 1

                time.sleep(0.1)

            except Exception as e:
                print(f"❌ [{self.agent_id}] 에러: {str(e)}")
                break

        print(f"✅ [{self.agent_id}] 프로세스 종료\n")

    def _process_message(self, message: dict) -> Optional[Dict]:
        """메시지 처리 및 응답 생성"""

        if isinstance(message, Message):
            msg_content = message.content
            from_agent = message.from_agent
        else:
            msg_content = message.get("content", "")
            from_agent = message.get("from_agent", "unknown")

        print(f"📨 [{self.agent_id}] 메시지 수신 (from: {from_agent})")
        print(f"   내용: {msg_content[:60]}...\n")

        # 역할별 응답 생성
        if self.role == AgentRole.ANALYST:
            response = self._analyst_response(msg_content)
            next_agent = "validator"

        elif self.role == AgentRole.VALIDATOR:
            response = self._validator_response(msg_content)
            next_agent = "optimizer"

        elif self.role == AgentRole.OPTIMIZER:
            response = self._optimizer_response(msg_content)
            next_agent = None

        else:
            response = "알겠습니다"
            next_agent = None

        self.responses.append({
            "timestamp": datetime.now().isoformat(),
            "from": from_agent,
            "response": response
        })

        if next_agent:
            return {"to": next_agent, "content": response}
        else:
            print(f"✅ [{self.agent_id}] 최종 응답 완료\n")
            return None

    def _analyst_response(self, content: str) -> str:
        """분석가 응답"""
        return f"""📊 분석 결과:

입력: {content[:40]}...

분석:
1. 핵심 요소: 3가지 주요 포인트 식별
2. 패턴: 반복되는 패턴 발견
3. 인사이트: 숨겨진 기회 발견

다음 단계: 검증 필요"""

    def _validator_response(self, content: str) -> str:
        """검증자 응답"""
        return f"""✓ 검증 결과:

입력 분석: {content[:30]}...

검증 체크리스트:
- ✓ 논리적 일관성: 확인됨
- ✓ 증거 기반성: 충분함
- ✓ 예외 사항: 없음

평가: 양호 (95점)

다음 단계: 최적화"""

    def _optimizer_response(self, content: str) -> str:
        """최적화자 응답"""
        return f"""⚡ 최적화 제안:

분석 내용: {content[:30]}...

최적화 방안:
1. 효율성: 30% 개선 가능
2. 비용: 15% 절감 가능
3. 속도: 2배 향상 가능

권장 조치: 즉시 실행 가능

결론: 준비 완료"""

# ============= 글로벌 오케스트레이터 =============
class GlobalOrchestrator:
    """모든 에이전트 프로세스를 관리하는 오케스트레이터"""

    def __init__(self):
        self.agents = {}
        self.processes = {}
        self.message_log = []
        self.message_queues = {}

    def create_agent(self, agent_id: str, role: AgentRole):
        """새 에이전트 프로세스 생성"""

        # 각 에이전트별 큐 생성
        inbox = multiprocessing.Queue()
        outbox = multiprocessing.Queue()
        broadcast = multiprocessing.Queue()

        self.message_queues[agent_id] = {
            "inbox": inbox,
            "outbox": outbox,
            "broadcast": broadcast
        }

        # 에이전트 프로세스 생성
        agent = ClaudeAgentProcess(agent_id, role, inbox, outbox, broadcast)
        process = multiprocessing.Process(target=agent.run)
        process.start()

        self.agents[agent_id] = agent
        self.processes[agent_id] = process

        print(f"🚀 [{agent_id}] 프로세스 시작됨 (PID: {process.pid})\n")

        return process

    def send_message(self, from_agent: str, to_agent: str, content: str):
        """메시지 전송"""

        message = {
            "id": f"msg_{len(self.message_log)}",
            "timestamp": datetime.now().isoformat(),
            "from_agent": from_agent,
            "to_agent": to_agent,
            "content": content,
            "process_id": os.getpid()
        }

        self.message_log.append(message)

        if to_agent in self.message_queues:
            self.message_queues[to_agent]["inbox"].put(message)
            print(f"📤 [{from_agent}] → [{to_agent}]")
            print(f"   {content[:60]}...\n")

    def stop_all_agents(self):
        """모든 에이전트 중지"""
        print(f"\n🛑 모든 에이전트 종료 중...\n")

        for agent_id, process in self.processes.items():
            self.message_queues[agent_id]["inbox"].put({"type": "STOP"})
            process.join(timeout=2)
            if process.is_alive():
                process.terminate()

            print(f"✅ [{agent_id}] 종료됨")

    def wait_completion(self, timeout=10):
        """모든 프로세스 완료 대기"""
        start_time = time.time()

        while time.time() - start_time < timeout:
            all_done = all(not p.is_alive() for p in self.processes.values())
            if all_done:
                return True
            time.sleep(0.5)

        return False

    def get_message_log(self) -> List[Dict]:
        """메시지 로그 반환"""
        return self.message_log

    def get_process_info(self) -> Dict:
        """프로세스 정보 반환"""
        return {
            agent_id: {
                "pid": process.pid,
                "alive": process.is_alive()
            }
            for agent_id, process in self.processes.items()
        }

# ============= 워크플로우 =============
def workflow_product_analysis():
    """워크플로우: 제품 분석"""

    print(f"\n{'='*70}")
    print(f"🌍 멀티프로세스 Claude 에이전트 협업")
    print(f"{'='*70}\n")

    # 오케스트레이터 생성
    orchestrator = GlobalOrchestrator()

    # 3개 에이전트 프로세스 생성
    print(f"{'─'*70}")
    print(f"[Step 1] 에이전트 프로세스 시작")
    print(f"{'─'*70}\n")

    orchestrator.create_agent("analyst", AgentRole.ANALYST)
    time.sleep(0.5)

    orchestrator.create_agent("validator", AgentRole.VALIDATOR)
    time.sleep(0.5)

    orchestrator.create_agent("optimizer", AgentRole.OPTIMIZER)
    time.sleep(0.5)

    # 프로세스 정보 출력
    print(f"\n📊 실행 중인 프로세스:")
    for agent_id, info in orchestrator.get_process_info().items():
        status = "🟢 실행 중" if info["alive"] else "🔴 종료됨"
        print(f"  {status} [{agent_id}] (PID: {info['pid']})")

    # 분석 요청 시작
    print(f"\n{'─'*70}")
    print(f"[Step 2] 분석 시작")
    print(f"{'─'*70}\n")

    initial_query = """새로운 AI 제품의 시장 잠재력을 분석해주세요.
타겟 시장: B2B SaaS
예상 수익: 연 1000만 달러"""

    orchestrator.send_message("main", "analyst", initial_query)
    time.sleep(2)

    # 결과 수집
    print(f"\n{'─'*70}")
    print(f"[Step 3] 결과 수집")
    print(f"{'─'*70}\n")

    time.sleep(5)

    # 모든 에이전트 종료
    orchestrator.stop_all_agents()
    orchestrator.wait_completion()

    # 최종 보고
    print(f"\n{'='*70}")
    print(f"📋 최종 보고")
    print(f"{'='*70}\n")

    print(f"✅ 전송된 메시지: {len(orchestrator.get_message_log())}개\n")

    print(f"📝 메시지 로그:")
    for idx, msg in enumerate(orchestrator.get_message_log(), 1):
        print(f"  {idx}. [{msg['from_agent']}] → [{msg['to_agent']}]")

    print(f"\n")

# ============= 메인 =============
if __name__ == "__main__":
    try:
        workflow_product_analysis()

        print(f"{'='*70}")
        print(f"✨ 멀티프로세스 협업 완료!")
        print(f"{'='*70}\n")

        print("🎯 시스템 특징:")
        print("  ✓ 각 에이전트가 독립적 프로세스로 실행")
        print("  ✓ 진정한 병렬 처리 (GIL 회피)")
        print("  ✓ 프로세스 간 큐 기반 통신")
        print("  ✓ 완전히 독립적인 메모리 공간\n")

        print("💡 확장 가능 영역:")
        print("  - 다른 머신의 에이전트 원격 실행")
        print("  - 에이전트 생성/제거 동적 관리")
        print("  - 에이전트 풀(pool) 관리")
        print("  - 프로세스 상태 모니터링")
        print("  - 에러 복구 및 재시작\n")

    except KeyboardInterrupt:
        print("\n\n⛔ 사용자 중단")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
