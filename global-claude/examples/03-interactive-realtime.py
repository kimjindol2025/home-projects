#!/usr/bin/env python3
"""
실시간 멀티프로세스 Claude 에이전트 협업
- 각 에이전트가 독립 프로세스로 병렬 실행
- 실시간 메시지 처리 및 응답
- 상태 동기화
"""

import multiprocessing as mp
import time
import json
from datetime import datetime
from typing import Dict, Optional
import sys
import os

# ============= 공유 상태 =============
class SharedState:
    """모든 프로세스가 공유하는 상태"""

    def __init__(self, manager: mp.Manager):
        self.lock = manager.Lock()
        self.messages = manager.list()
        self.agent_status = manager.dict()
        self.knowledge = manager.dict()

# ============= Claude 에이전트 =============
class ClaudeAgent:
    """독립 프로세스로 실행되는 에이전트"""

    def __init__(self, agent_id: str, role: str, inbox: mp.Queue,
                 shared_state: SharedState, outbox: mp.Queue):
        self.id = agent_id
        self.role = role
        self.inbox = inbox
        self.shared_state = shared_state
        self.outbox = outbox
        self.pid = os.getpid()
        self.message_count = 0
        self.running = True

    def _log(self, message: str):
        """로그 출력"""
        timestamp = datetime.now().strftime("%H:%M:%S")
        print(f"[{timestamp}] [{self.id} (PID:{self.pid})] {message}")

    def _generate_response(self, content: str) -> str:
        """역할 기반 응답 생성"""

        if self.role == "analyst":
            return f"""분석 완료:
✓ 내용 분석: {content[:30]}...
✓ 패턴 식별: 3가지 주요 패턴
✓ 인사이트: 숨겨진 기회 2개 발견
→ 검증팀으로 이관"""

        elif self.role == "validator":
            return f"""검증 완료:
✓ 논리 검증: 통과
✓ 증거 확인: 충분
✓ 신뢰도: 95%
→ 최적화팀으로 이관"""

        elif self.role == "optimizer":
            return f"""최적화 완료:
✓ 효율성 개선: 40% ↑
✓ 비용 절감: 25% ↓
✓ 속도 향상: 3배
→ 최종 완료"""

        else:
            return f"처리 완료: {content[:30]}..."

    def _process_message(self, msg: Dict) -> Optional[Dict]:
        """메시지 처리"""

        self._log(f"📨 [{msg['from']}]로부터 메시지 수신")
        self._log(f"   내용: {msg['content'][:50]}...")

        # 응답 생성
        response = self._generate_response(msg['content'])

        time.sleep(1)  # 처리 시간 시뮬레이션

        self._log(f"✅ 응답 생성 완료")
        self._log(f"   {response[:50]}...")

        # 다음 에이전트 결정
        next_agent = None
        if self.role == "analyst":
            next_agent = "validator"
        elif self.role == "validator":
            next_agent = "optimizer"

        return {
            "id": f"msg_{self.message_count}",
            "from": self.id,
            "to": next_agent,
            "content": response,
            "timestamp": datetime.now().isoformat(),
            "process_id": self.pid
        }

    def run(self):
        """에이전트 실행 루프"""

        self._log(f"🚀 에이전트 시작 (역할: {self.role})")

        # 상태 등록
        with self.shared_state.lock:
            self.shared_state.agent_status[self.id] = "ready"

        while self.running:
            try:
                # 메시지 수신 (타임아웃 1초)
                if not self.inbox.empty():
                    msg = self.inbox.get(timeout=0.5)

                    # 중지 신호 확인
                    if isinstance(msg, dict) and msg.get("type") == "STOP":
                        self._log(f"⛔ 종료 신호 수신")
                        self.running = False
                        break

                    # 메시지 처리
                    response = self._process_message(msg)

                    # 응답 전송 (다음 에이전트가 있으면)
                    if response and response["to"]:
                        self.outbox.put(response)
                        self._log(f"📤 [{response['to']}]로 메시지 전송")

                    # 공유 상태 업데이트
                    with self.shared_state.lock:
                        self.shared_state.messages.append(response)
                        self.shared_state.knowledge[self.id] = response["content"]

                    self.message_count += 1

                time.sleep(0.1)

            except Exception as e:
                self._log(f"❌ 에러: {str(e)}")
                break

        with self.shared_state.lock:
            self.shared_state.agent_status[self.id] = "stopped"

        self._log(f"🔴 에이전트 종료")

# ============= 오케스트레이터 =============
class ProcessOrchestrator:
    """프로세스 관리 및 조율"""

    def __init__(self):
        self.manager = mp.Manager()
        self.shared_state = SharedState(self.manager)
        self.agents = {}
        self.processes = {}
        self.queues = {}

    def create_agent(self, agent_id: str, role: str) -> mp.Queue:
        """에이전트 프로세스 생성"""

        # 큐 생성
        inbox = mp.Queue()
        outbox = mp.Queue()

        self.queues[agent_id] = {
            "inbox": inbox,
            "outbox": outbox
        }

        # 에이전트 생성 및 프로세스 시작
        agent = ClaudeAgent(agent_id, role, inbox, self.shared_state, outbox)
        process = mp.Process(target=agent.run)
        process.start()

        self.agents[agent_id] = agent
        self.processes[agent_id] = process

        print(f"✅ [{agent_id}] 프로세스 시작됨 (PID: {process.pid})")

        return outbox

    def connect_agents(self, from_agent: str, to_agent: str):
        """에이전트 간 연결"""

        outbox = self.queues[from_agent]["outbox"]
        inbox = self.queues[to_agent]["inbox"]

        # 아웃박스의 메시지를 다음 에이전트의 인박스로 라우팅
        self.queues[from_agent]["outbox_target"] = to_agent
        self.queues[to_agent]["inbox_source"] = from_agent

        print(f"🔗 [{from_agent}] → [{to_agent}] 연결됨")

    def send_initial_message(self, agent_id: str, content: str):
        """초기 메시지 전송"""

        msg = {
            "id": "msg_0",
            "from": "user",
            "to": agent_id,
            "content": content,
            "timestamp": datetime.now().isoformat(),
            "type": "query"
        }

        self.queues[agent_id]["inbox"].put(msg)
        print(f"📤 [user] → [{agent_id}] 초기 메시지 전송")

    def route_messages(self):
        """메시지 라우팅 루프"""

        while True:
            try:
                for agent_id, queues in self.queues.items():
                    outbox = queues["outbox"]

                    while not outbox.empty():
                        msg = outbox.get(timeout=0.1)

                        if msg and msg.get("to"):
                            target = msg["to"]
                            if target in self.queues:
                                self.queues[target]["inbox"].put(msg)
                                timestamp = datetime.now().strftime("%H:%M:%S")
                                print(f"[{timestamp}] [Router] 메시지 라우팅: "
                                      f"[{msg['from']}] → [{target}]")

                time.sleep(0.1)

            except:
                break

    def wait_completion(self, timeout: int = 30):
        """모든 프로세스 완료 대기"""

        start = time.time()

        while time.time() - start < timeout:
            all_ready = all(
                self.processes[pid].is_alive()
                for pid in self.processes
            )

            if not all_ready and any(
                not self.processes[pid].is_alive()
                for pid in self.processes
            ):
                return True

            time.sleep(0.5)

        return False

    def stop_all(self):
        """모든 프로세스 중지"""

        print(f"\n{'─'*70}")
        print(f"🛑 모든 에이전트 종료 중...\n")

        for agent_id in self.agents:
            self.queues[agent_id]["inbox"].put({"type": "STOP"})

        for agent_id, process in self.processes.items():
            process.join(timeout=2)
            if process.is_alive():
                process.terminate()

        print(f"\n✅ 모든 에이전트 종료됨\n")

    def get_status(self) -> Dict:
        """상태 조회"""

        return {
            "timestamp": datetime.now().isoformat(),
            "agents": {
                agent_id: {
                    "pid": process.pid,
                    "alive": process.is_alive(),
                    "status": self.shared_state.agent_status.get(agent_id, "unknown")
                }
                for agent_id, process in self.processes.items()
            },
            "message_count": len(self.shared_state.messages),
            "knowledge": dict(self.shared_state.knowledge)
        }

# ============= 메인 워크플로우 =============
def main():
    """메인 워크플로우"""

    print(f"\n{'='*70}")
    print(f"🌍 멀티프로세스 Claude 에이전트 협업 (실시간)")
    print(f"{'='*70}\n")

    orchestrator = ProcessOrchestrator()

    # Step 1: 에이전트 생성
    print(f"{'─'*70}")
    print(f"[Step 1] 에이전트 프로세스 생성")
    print(f"{'─'*70}\n")

    orchestrator.create_agent("analyst", "analyst")
    time.sleep(0.5)

    orchestrator.create_agent("validator", "validator")
    time.sleep(0.5)

    orchestrator.create_agent("optimizer", "optimizer")
    time.sleep(1)

    # Step 2: 에이전트 연결
    print(f"\n{'─'*70}")
    print(f"[Step 2] 에이전트 간 연결")
    print(f"{'─'*70}\n")

    orchestrator.connect_agents("analyst", "validator")
    orchestrator.connect_agents("validator", "optimizer")
    time.sleep(0.5)

    # Step 3: 메시지 라우팅 스레드 시작 (별도 스레드)
    import threading

    router_thread = threading.Thread(
        target=orchestrator.route_messages,
        daemon=True
    )
    router_thread.start()
    time.sleep(0.5)

    # Step 4: 초기 메시지 전송
    print(f"\n{'─'*70}")
    print(f"[Step 3] 분석 요청")
    print(f"{'─'*70}\n")

    query = """새로운 AI 제품의 시장 진출 전략을 분석해주세요.
타겟 시장: B2B Enterprise
예상 수익: 연 500만 달러"""

    orchestrator.send_initial_message("analyst", query)

    # Step 5: 결과 수집
    print(f"\n{'─'*70}")
    print(f"[Step 4] 결과 수집 (20초 대기)")
    print(f"{'─'*70}\n")

    time.sleep(20)

    # Step 6: 최종 상태
    print(f"\n{'─'*70}")
    print(f"[Step 5] 최종 상태 리포트")
    print(f"{'─'*70}\n")

    status = orchestrator.get_status()

    print(f"✅ 총 메시지: {status['message_count']}개\n")

    print(f"📊 에이전트 상태:")
    for agent_id, info in status['agents'].items():
        alive = "🟢 실행 중" if info["alive"] else "🔴 종료됨"
        print(f"  {alive} [{agent_id}] (PID: {info['pid']}, Status: {info['status']})")

    print(f"\n📚 수집된 지식:")
    for agent_id, knowledge in status['knowledge'].items():
        print(f"  [{agent_id}]:")
        for line in knowledge.split('\n')[:3]:
            if line.strip():
                print(f"    {line}")

    # Step 7: 종료
    orchestrator.stop_all()

    print(f"{'='*70}")
    print(f"✨ 멀티프로세스 협업 완료!")
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
