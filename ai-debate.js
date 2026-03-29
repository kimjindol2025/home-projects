#!/usr/bin/env node
/**
 * AI 멀티 에이전트 토론 시스템
 * Claude, Gemini, GPT-4가 한 주제로 서로 대화
 */

const Anthropic = require("@anthropic-ai/sdk");
const fetch = require("node-fetch");

const claude = new Anthropic();

// ============= AI 에이전트 정의 =============
const agents = {
  claude: {
    name: "Claude 🔵",
    emoji: "🔵",
    call: async (prompt) => {
      const response = await claude.messages.create({
        model: "claude-opus-4-6",
        max_tokens: 400,
        messages: [{ role: "user", content: prompt }]
      });
      return response.content[0].text;
    }
  },
  gemini: {
    name: "Gemini ✨",
    emoji: "✨",
    call: async (prompt) => {
      const response = await fetch(
        `https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${process.env.GEMINI_API_KEY}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            contents: [{ parts: [{ text: prompt }] }],
            generationConfig: { temperature: 0.8, maxOutputTokens: 400 }
          })
        }
      );
      const data = await response.json();
      if (!data.candidates?.[0]?.content?.parts?.[0]) {
        throw new Error("Gemini API error: " + JSON.stringify(data));
      }
      return data.candidates[0].content.parts[0].text;
    }
  },
  gpt: {
    name: "ChatGPT 🔴",
    emoji: "🔴",
    call: async (prompt) => {
      const response = await fetch("https://api.openai.com/v1/chat/completions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${process.env.OPENAI_API_KEY}`
        },
        body: JSON.stringify({
          model: "gpt-4-turbo",
          messages: [{ role: "user", content: prompt }],
          max_tokens: 400,
          temperature: 0.8
        })
      });
      const data = await response.json();
      if (!data.choices?.[0]?.message?.content) {
        throw new Error("OpenAI API error: " + JSON.stringify(data));
      }
      return data.choices[0].message.content;
    }
  }
};

// ============= 토론 시스템 =============
async function startDebate(topic, rounds = 3) {
  console.log(`\n${"=".repeat(70)}`);
  console.log(`📋 주제: ${topic}`);
  console.log(`🎬 라운드: ${rounds}라운드`);
  console.log(`${"=".repeat(70)}\n`);

  const agentNames = Object.keys(agents);
  let conversation = [];

  // Round 1: 각자 초기 의견
  console.log(`\n🔔 [라운드 1] 초기 의견 제시\n`);

  for (const agentName of agentNames) {
    try {
      const agent = agents[agentName];
      const response = await agent.call(
        `주제: "${topic}"\n\n이 주제에 대해 간단하고 명확하게 당신의 관점을 설명해줘. (2-3문장)`
      );

      console.log(`${agent.emoji} ${agent.name}:`);
      console.log(`${response}\n`);

      conversation.push({
        round: 1,
        agent: agentName,
        message: response
      });

      await sleep(1000); // API 레이트 리밋 대비
    } catch (error) {
      console.error(`❌ ${agentName} 오류: ${error.message}\n`);
    }
  }

  // Round 2-N: 서로 반박
  for (let round = 2; round <= rounds; round++) {
    console.log(`\n🔔 [라운드 ${round}] 상호 토론\n`);

    for (let i = 0; i < agentNames.length; i++) {
      const currentAgent = agentNames[i];
      const lastAgent = agentNames[(i - 1 + agentNames.length) % agentNames.length];

      try {
        const lastMessage = conversation
          .filter(m => m.agent === lastAgent)
          .pop();

        const agent = agents[currentAgent];
        const prompt = `
주제: "${topic}"

${agents[lastAgent].name}가 이렇게 말했어:
"${lastMessage?.message || ''}"

이것에 대해 동의하거나 다른 관점을 제시해줄래? (2-3문장)
`;

        const response = await agent.call(prompt);

        console.log(`${agent.emoji} ${agent.name}:`);
        console.log(`${response}\n`);

        conversation.push({
          round,
          agent: currentAgent,
          message: response
        });

        await sleep(1000);
      } catch (error) {
        console.error(`❌ ${currentAgent} 오류: ${error.message}\n`);
      }
    }
  }

  // 요약
  console.log(`\n${"=".repeat(70)}`);
  console.log(`📊 토론 요약`);
  console.log(`${"=".repeat(70)}\n`);

  for (const agentName of agentNames) {
    const agentMessages = conversation.filter(m => m.agent === agentName);
    console.log(`${agents[agentName].emoji} ${agents[agentName].name}: ${agentMessages.length}개 발언`);
  }
}

// ============= 토론 팀 모드 =============
async function teamDebate(topic) {
  console.log(`\n${"=".repeat(70)}`);
  console.log(`🏢 팀 모드: Pro vs Con`);
  console.log(`📋 주제: ${topic}`);
  console.log(`${"=".repeat(70)}\n`);

  const proTeam = ["claude", "gemini"];
  const conTeam = ["gpt"];

  // 입장 정리
  console.log(`🔵 찬성팀 (Claude, Gemini)`);
  for (const agent of proTeam) {
    try {
      const response = await agents[agent].call(
        `주제: "${topic}"\n\n이것이 좋은 이유를 설득력 있게 설명해줘. (2-3문장, 찬성 입장으로)`
      );
      console.log(`  ${agents[agent].emoji} ${agents[agent].name}: ${response}\n`);
      await sleep(1000);
    } catch (error) {
      console.error(`❌ ${agent} 오류\n`);
    }
  }

  console.log(`\n🔴 반대팀 (ChatGPT)`);
  try {
    const response = await agents.gpt.call(
      `주제: "${topic}"\n\n이것의 문제점을 설득력 있게 설명해줘. (2-3문장, 반대 입장으로)`
    );
    console.log(`  ${agents.gpt.emoji} ${agents.gpt.name}: ${response}\n`);
  } catch (error) {
    console.error(`❌ gpt 오류\n`);
  }

  // 재반박
  console.log(`\n🔵 찬성팀 최종 재반박`);
  try {
    const response = await agents.claude.call(
      `주제: "${topic}"\n\n반대팀의 우려에 대해 최종 재반박을 해줄래? (2-3문장)`
    );
    console.log(`  ${agents.claude.emoji} Claude: ${response}\n`);
  } catch (error) {
    console.error(`❌ claude 오류\n`);
  }
}

// ============= 유틸 =============
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// ============= 메인 =============
async function main() {
  // 환경변수 확인
  if (!process.env.GEMINI_API_KEY) {
    console.warn("⚠️  GEMINI_API_KEY 미설정 (Gemini 사용 불가)");
  }
  if (!process.env.OPENAI_API_KEY) {
    console.warn("⚠️  OPENAI_API_KEY 미설정 (GPT 사용 불가)");
  }

  // 토론 실행
  const topic = process.argv[2] || "AI가 인류에 미치는 영향";
  const mode = process.argv[3] || "debate"; // debate 또는 team

  try {
    if (mode === "team") {
      await teamDebate(topic);
    } else {
      await startDebate(topic, 2);
    }
  } catch (error) {
    console.error("❌ 에러:", error.message);
    process.exit(1);
  }
}

main();
