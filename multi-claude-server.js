#!/usr/bin/env node
/**
 * 서버 내 멀티 Claude 에이전트 통신 시스템
 * 같은 서버에서 여러 Claude 인스턴스가 메시지 교환
 */

const Anthropic = require("@anthropic-ai/sdk");
const EventEmitter = require("events");

// ============= 메시지 버스 (IPC) =============
class MessageBus extends EventEmitter {
  constructor() {
    super();
    this.agents = new Map();
    this.messageLog = [];
  }

  register(agentId, agent) {
    this.agents.set(agentId, agent);
    console.log(`✅ 에이전트 등록: ${agentId}`);
  }

  async send(fromId, toId, message) {
    const envelope = {
      timestamp: new Date().toISOString(),
      from: fromId,
      to: toId,
      message,
      id: `msg_${Date.now()}_${Math.random().toString(36).slice(2)}`
    };

    this.messageLog.push(envelope);
    console.log(`\n📨 [${fromId}] → [${toId}]`);
    console.log(`   메시지: ${message.substring(0, 100)}...`);

    // 수신자에게 메시지 전달
    const recipient = this.agents.get(toId);
    if (recipient) {
      await recipient.receiveMessage(envelope);
    } else {
      console.warn(`⚠️  수신자 ${toId}를 찾을 수 없습니다`);
    }

    return envelope;
  }

  getLog() {
    return this.messageLog;
  }
}

// ============= Claude 에이전트 =============
class ClaudeAgent {
  constructor(id, role, bus) {
    this.id = id;
    this.role = role;
    this.bus = bus;
    this.client = new Anthropic();
    this.inbox = [];
    this.messageHistory = [];

    // 역할별 시스템 프롬프트
    this.systemPrompts = {
      researcher: `당신은 연구원입니다. 주어진 주제에 대해 깊이 있게 분석합니다.
다른 에이전트의 피드백을 수집하고 종합합니다.`,
      critic: `당신은 비평가입니다. 제시된 아이디어를 비판적으로 검토합니다.
논리적 결함과 개선점을 지적합니다.`,
      synthesizer: `당신은 종합가입니다. 여러 의견을 통합하여 최종 결론을 도출합니다.
모든 관점을 존중하면서 균형잡힌 결론을 만듭니다.`
    };

    this.bus.register(id, this);
  }

  async receiveMessage(envelope) {
    this.inbox.push(envelope);
    console.log(`\n💬 [${this.id}] 메시지 수신 완료`);
  }

  async processInbox() {
    if (this.inbox.length === 0) {
      console.log(`[${this.id}] 받은 메시지 없음`);
      return null;
    }

    const message = this.inbox.shift();
    console.log(`\n🔄 [${this.id}] 메시지 처리 중...`);

    // Claude가 응답 생성
    const systemPrompt = this.systemPrompts[this.role] || this.systemPrompts.researcher;

    const response = await this.client.messages.create({
      model: "claude-opus-4-6",
      max_tokens: 500,
      system: systemPrompt,
      messages: [
        {
          role: "user",
          content: `다른 에이전트(${message.from})로부터 다음 메시지를 받았습니다:\n\n"${message.message}"\n\n이에 대해 당신의 관점을 간단히 (2-3문장) 제시해주세요.`
        }
      ]
    });

    const reply = response.content[0].text;
    this.messageHistory.push({
      received: message,
      response: reply
    });

    return reply;
  }

  async sendMessage(targetId, message) {
    await this.bus.send(this.id, targetId, message);
  }

  getHistory() {
    return this.messageHistory;
  }
}

// ============= 협업 워크플로우 =============
async function collaborativeWorkflow() {
  console.log(`\n${"=".repeat(70)}`);
  console.log(`🤝 서버 내 멀티 Claude 협업 시스템`);
  console.log(`${"=".repeat(70)}\n`);

  // 메시지 버스 생성
  const bus = new MessageBus();

  // 3개 에이전트 생성
  const researcher = new ClaudeAgent("researcher", "researcher", bus);
  const critic = new ClaudeAgent("critic", "critic", bus);
  const synthesizer = new ClaudeAgent("synthesizer", "synthesizer", bus);

  const topic = "AI의 미래는 밝을까?";
  console.log(`\n📌 주제: ${topic}\n`);

  // Step 1: 연구원이 초기 관점 제시
  console.log(`\n${"─".repeat(70)}`);
  console.log(`[Step 1] 연구원의 초기 관점`);
  console.log(`${"─".repeat(70)}`);

  const initialView = `AI 기술은 빠르게 발전하고 있습니다. 특히 생성형 AI의 발전으로 많은 산업에서
혁신이 일어나고 있습니다. 다만 규제와 윤리 문제도 동시에 고려해야 합니다.`;

  console.log(`\n🔵 Researcher:\n${initialView}\n`);
  researcher.messageHistory.push({
    received: null,
    response: initialView
  });

  // Step 2: 비평가가 비판 제시
  console.log(`\n${"─".repeat(70)}`);
  console.log(`[Step 2] 비평가의 검토`);
  console.log(`${"─".repeat(70)}`);

  await researcher.sendMessage("critic", initialView);
  await new Promise(resolve => setTimeout(resolve, 1000));

  const criticResponse = await critic.processInbox();
  console.log(`\n🔴 Critic:\n${criticResponse}\n`);

  // Step 3: 종합가가 통합 의견 제시
  console.log(`\n${"─".repeat(70)}`);
  console.log(`[Step 3] 종합가의 최종 결론`);
  console.log(`${"─".repeat(70)}`);

  const combinedInput = `
초기 의견:
"${initialView}"

비평:
"${criticResponse}"

이 두 관점을 종합하여 균형잡힌 결론을 제시해주세요.`;

  await researcher.sendMessage("synthesizer", combinedInput);
  await new Promise(resolve => setTimeout(resolve, 1000));

  const finalConclusion = await synthesizer.processInbox();
  console.log(`\n🟢 Synthesizer:\n${finalConclusion}\n`);

  // Step 4: 결과 정리
  console.log(`\n${"=".repeat(70)}`);
  console.log(`📊 협업 결과 요약`);
  console.log(`${"=".repeat(70)}\n`);

  console.log(`✅ 전송된 메시지: ${bus.messageLog.length}개`);
  console.log(`✅ 활성 에이전트: 3개`);
  console.log(`✅ 완료된 상호작용: 2개\n`);

  console.log(`메시지 로그:`);
  bus.getLog().forEach((log, idx) => {
    console.log(`  ${idx + 1}. [${log.from}] → [${log.to}]`);
  });
}

// ============= 실시간 토론 (더 복잡한 예제) =============
async function realtimeDebate() {
  console.log(`\n${"=".repeat(70)}`);
  console.log(`💬 실시간 멀티 Claude 토론`);
  console.log(`${"=".repeat(70)}\n`);

  const bus = new MessageBus();
  const agentA = new ClaudeAgent("agentA", "researcher", bus);
  const agentB = new ClaudeAgent("agentB", "critic", bus);

  const topic = "원격 근무는 생산성을 높이는가?";

  // Agent A의 주장
  console.log(`\n[Round 1] Agent A의 주장`);
  const argumentA = `원격 근무는 통근 시간을 절약하고, 개발자들이 더 집중할 수 있는 환경을
제공합니다. 실제 많은 기업들이 원격 근무 도입 후 생산성 향상을 보고했습니다.`;
  console.log(`🟦 Agent A:\n${argumentA}\n`);
  agentA.messageHistory.push({ received: null, response: argumentA });

  // Agent B의 반박
  await agentA.sendMessage("agentB", argumentA);
  await new Promise(resolve => setTimeout(resolve, 1000));

  console.log(`\n[Round 2] Agent B의 반박`);
  const rebuttal = await agentB.processInbox();
  console.log(`🟥 Agent B:\n${rebuttal}\n`);

  // Agent A의 재반박
  await agentB.sendMessage("agentA", rebuttal);
  await new Promise(resolve => setTimeout(resolve, 1000));

  console.log(`\n[Round 3] Agent A의 재반박`);
  const counter = await agentA.processInbox();
  console.log(`🟦 Agent A:\n${counter}\n`);

  console.log(`\n✅ 토론 완료! (총 3라운드)\n`);
}

// ============= 메인 =============
async function main() {
  try {
    // 협업 워크플로우
    await collaborativeWorkflow();

    // 실시간 토론
    await realtimeDebate();

    console.log(`\n${"=".repeat(70)}`);
    console.log(`✨ 모든 에이전트 상호작용 완료`);
    console.log(`${"=".repeat(70)}\n`);
  } catch (error) {
    console.error("❌ 에러:", error.message);
    process.exit(1);
  }
}

main();
