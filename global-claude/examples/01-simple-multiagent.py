#!/usr/bin/env python3
"""
간단한 멀티프로세스 Claude 에이전트 대화 시스템
- 각 에이전트가 독립 프로세스로 실행
- 중앙 메시지 큐로 모든 메시지 관리
- 자유로운 양방향 대화
"""

import multiprocessing as mp
import time
from datetime import datetime
import sys

# ============= 글로벌 메시지 큐 =============
message_queue = mp.Queue()  # 모든 에이전트가 사용

# ============= Claude 에이전트 =============
class ChatAgent:
    """대화형 에이전트"""

    def __init__(self, name: str, personality: str, message_queue: mp.Queue):
        self.name = name
        self.personality = personality
        self.queue = message_queue
        self.message_count = 0

    def _log(self, msg: str):
        """로그"""
        t = datetime.now().strftime("%H:%M:%S")
        print(f"[{t}] 💬 [{self.name}] {msg}")

    def _respond(self, content: str) -> str:
        """성격에 따른 응답"""

        if self.personality == "optimist":
            return f"""좋은 질문입니다! 제 생각:
✅ 긍정적 접근이 필요합니다
💡 시도해볼 가치가 충분합니다
🎯 함께 성공할 수 있을 겁니다!"""

        elif self.personality == "realist":
            return f"""현실적으로 분석하면:
📊 장단점을 균형있게 봐야 합니다
⚠️ 리스크 관리가 중요합니다
✓ 실행 계획이 명확해야 합니다"""

        elif self.personality == "critic":
            return f"""비판적 검토:
🔍 몇 가지 우려사항이 있습니다
⚡ 더 견고한 근거가 필요합니다
💭 이런 점들을 개선해야 합니다"""

        else:
            return f"""흥미로운 주제네요!
🚀 더 깊이 논의해봅시다
💬 당신의 의견은 어떤가요?"""

    def run(self):
        """에이전트 실행 루프"""

        self._log(f"🚀 시작 (성격: {self.personality})")

        message_buffer = []
        conversation_round = 0

        while True:
            try:
                # 메시지 수신 (타임아웃 1초)
                if not self.queue.empty():
                    try:
                        msg = self.queue.get(timeout=0.5)

                        # 종료 신호
                        if msg.get("type") == "STOP":
                            self._log("⛔ 종료 신호 수신")
                            break

                        # 자신의 메시지면 스킵
                        if msg.get("from") == self.name:
                            continue

                        # 메시지 처리
                        sender = msg.get("from", "unknown")
                        content = msg.get("content", "")

                        self._log(f"📨 [{sender}]로부터 메시지 수신")
                        self._log(f"   내용: {content[:50]}...")

                        # 응답 생성
                        time.sleep(1)  # 생각하는 시간
                        response = self._respond(content)

                        self._log(f"✅ 응답 생성")

                        # 응답 전송
                        response_msg = {
                            "from": self.name,
                            "to": sender,
                            "content": response,
                            "type": "response",
                            "timestamp": datetime.now().isoformat()
                        }

                        self.queue.put(response_msg)
                        self._log(f"📤 [{sender}]로 응답 전송\n")

                        self.message_count += 1
                        conversation_round += 1

                        # 3라운드 후 종료
                        if conversation_round >= 3:
                            self._log("🎯 충분한 대화 후 대기 중...")

                    except Exception as e:
                        continue

                time.sleep(0.2)

            except KeyboardInterrupt:
                break
            except Exception as e:
                self._log(f"❌ 에러: {str(e)}")
                break

        self._log(f"🔴 종료 (총 {self.message_count}개 메시지 처리)")

# ============= 글로벌 메시지 전송 함수 =============
def send_message(queue: mp.Queue, from_name: str, to_name: str, content: str):
    """메시지 전송"""

    msg = {
        "from": from_name,
        "to": to_name,
        "content": content,
        "type": "message",
        "timestamp": datetime.now().isoformat()
    }

    queue.put(msg)
    t = datetime.now().strftime("%H:%M:%S")
    print(f"[{t}] 📤 [user] → [{from_name}] 메시지 전송")
    print(f"    {content[:60]}...\n")

# ============= 메인 =============
def main():
    """메인 실행"""

    print(f"\n{'='*70}")
    print(f"💬 멀티프로세스 Claude 에이전트 - 자유로운 대화")
    print(f"{'='*70}\n")

    # 글로벌 큐 생성 (모든 프로세스가 공유)
    ctx = mp.get_context('spawn')
    global_queue = ctx.Queue()

    # Step 1: 에이전트 생성
    print(f"{'─'*70}")
    print(f"[Step 1] 에이전트 프로세스 생성\n")

    agents_config = [
        ("optimist", "optimist"),
        ("realist", "realist"),
        ("critic", "critic")
    ]

    processes = []

    for name, personality in agents_config:
        agent = ChatAgent(name, personality, global_queue)
        p = ctx.Process(target=agent.run)
        p.start()
        processes.append(p)
        print(f"✅ [{name}] 프로세스 시작 (PID: {p.pid})")
        time.sleep(0.3)

    print()

    # Step 2: 초기 메시지 전송
    print(f"{'─'*70}")
    print(f"[Step 2] 대화 시작\n")

    topic = "새로운 AI 제품을 개발해야 할까요?"
    send_message(global_queue, "user", "optimist", topic)

    # Step 3: 대화 진행
    print(f"{'─'*70}")
    print(f"[Step 3] 대화 진행 중... (30초)\n")

    time.sleep(30)

    # Step 4: 모든 에이전트 종료
    print(f"\n{'─'*70}")
    print(f"[Step 4] 모든 에이전트 종료\n")

    for _ in range(len(processes)):
        global_queue.put({"type": "STOP"})

    for name, p in zip([cfg[0] for cfg in agents_config], processes):
        p.join(timeout=3)
        if p.is_alive():
            p.terminate()
        print(f"✅ [{name}] 종료됨")

    print(f"\n{'='*70}")
    print(f"✨ 대화 완료!")
    print(f"{'='*70}\n")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n⛔ 중단됨")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ 에러: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
