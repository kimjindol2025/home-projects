# Community Engagement — 2026-03-18

**작성자**: Community Manager
**날짜**: 2026-03-18 (수요일) 10:00 KST
**검토 기준**: brand-voice.md, content-policy.md 준수

---

## 1. GeekNews 본문 포스트

### 제목
**노드를 연결하면 FreeLang 코드가 나옵니다 — FreeWire 개발 후기**

### 본문

안녕하세요. 사이드 프로젝트를 공유합니다.

**FreeWire**는 브라우저에서 노드를 연결하면 FreeLang 컴포넌트 코드(`.free` 파일)가 자동으로 생성되는 시각 프로그래밍 도구입니다. 언리얼 엔진의 Blueprints나 Node-RED를 써보셨다면 "아, 그 방식"이라고 바로 이해하실 것 같아요.

---

**왜 만들었나요?**

FreeLang은 제가 만들고 있는 경량 컴포넌트 DSL인데, 코드를 직접 타이핑하기 전에 "일단 연결해보고 싶다"는 아이디어에서 시작했습니다. 처음엔 단순히 시각화 도구로 기획했는데, 개발하다 보니 질문이 생겼어요.

> "노드가 '왜' 거기 있어야 하는지를 도구가 알 수는 없을까?"

그래서 Intent 엔진을 붙이게 됐습니다. 각 노드에 `contract` (입출력 조건)와 `goal` (달성 목표)를 명세할 수 있고, 연결 후 자동으로 충족도를 0-100점으로 평가합니다. 점수가 낮으면 수정 패치를 제안하고, 승인하면 자동으로 적용됩니다.

---

**기술 스택**

- 노드 에디터: Canvas/SVG + 드래그·연결·줌/패닝 (외부 프레임워크 없음)
- Graph 컴파일러: 노드 JSON → FreeLang `.free` 코드
- FreeExecutor: `.free` → HTML/CSS/JS 즉시 변환 (서버 불필요)
- AI 자동 생성: 자연어 입력 → 노드 그래프 생성 (15개 노드 타입)
- Live Preview: 브라우저 안에서 iframe으로 즉시 실행

---

**개발하면서 어려웠던 점**

가장 고민했던 부분은 **자동 수정 루프**입니다.

계약 위반을 감지 → 패치 생성 → 적용 → 재평가 하는 폐쇄 루프인데, 이게 무한 루프가 될 수 있거든요. "고쳤더니 다른 곳이 깨지는" 상황. 해결책은 패치마다 고유 ID를 부여하고 동일 패턴이 N회 이상 반복되면 루프를 중단 + 경고를 띄우는 방식이었습니다.

또 하나는 **학습 기억**입니다. 수정 이력을 localStorage에 쌓아두고, 같은 노드 타입에서 비슷한 실수가 반복되면 사전에 경고해줍니다. 세션이 끊겨도 기억이 유지돼야 했기 때문에 직렬화 설계에 신경을 많이 썼어요.

---

**현재 상태**

Phase 1-6 완료, 162개 테스트 통과입니다.

- Phase 1: 노드 JSON 스펙 (32/32)
- Phase 2: Graph→FreeLang 컴파일러 (16/16)
- Phase 3: API 서버 + FreeExecutor (24/24)
- Phase 4: 노드 에디터 UI (30/30)
- Phase 5: AI 노드 자동 생성기 (35/35)
- Phase 6: 브라우저 Live Preview (25/25)

---

**실행해보기**

서버 없이 바로 실행됩니다.

```bash
git clone https://gogs.dclub.kr/kim/FreeWire
cd FreeWire
python3 -m http.server 5500 --directory www
# http://localhost:5500 접속
```

AI 자동 생성 기능은 API 서버가 필요합니다.

```bash
npm install
node src/api/server.js
```

---

**피드백 환영합니다**

- "이런 노드 타입이 있으면 좋겠다"
- "Intent 엔진 설계에서 다른 접근을 써봤는데..."
- "FreeLang 문법이 이해가 안 된다"

모두 반갑습니다. GOGS 저장소: https://gogs.dclub.kr/kim/FreeWire

아직 초기 단계라 부족한 부분이 많지만, "노드가 자기 자신의 목적을 안다"는 아이디어 자체에 대한 이야기를 나눠보고 싶어서 공유했습니다.

---

*태그: #시각프로그래밍 #FreeLang #브라우저툴 #사이드프로젝트 #노드에디터*

---

## 2. Reddit 댓글 초안

### 2-1. r/rust — "Visual programming tools in Rust/WASM" 관련 스레드

**타겟 스레드 유형**: "Has anyone built a visual programming editor with Rust + WASM?" 또는 "Showcasing my node editor built in Rust" 등

**댓글 초안**:

---

> Great timing on this thread — I've been down a similar rabbit hole recently.

I built a node-based visual editor called **FreeWire** for a DSL I'm working on (FreeLang). The editor itself runs pure in the browser (Canvas/SVG, no framework), but the interesting challenge was on the compilation side: the node graph needs to emit structured DSL code, not just visual output.

The part that took the most iteration was what I'm calling an "Intent Engine" — each node declares a `contract` (input/output preconditions) and a `goal` (what it should accomplish structurally). After each connection, the engine evaluates fulfillment 0-100 and proposes patches if a contract is violated. The tricky bit was preventing the auto-repair loop from becoming infinite — solved with patch ID tracking and a cycle-break on repeated patterns.

On the WASM side: I'm planning to port the FreeLang executor (currently Node.js) to Rust + WASM so the live preview can run fully compiled in the browser. Currently it compiles node graphs to HTML/CSS/JS via a JS executor, which works fine but a Rust-backed pipeline would be more consistent.

If you're doing something similar with the Rust + WASM approach, I'd be curious how you're handling the AST-to-DOM pipeline — especially state binding across recompiles.

Source: https://gogs.dclub.kr/kim/FreeWire (162 tests passing across 6 phases)

---

**글자 수**: 약 1,050자 (영문) / 진정성 있는 기술 교환 형태, FreeWire는 자연스럽게 언급

---

### 2-2. r/programming — "Visual / node-based programming: interesting or overhyped?" 스레드

**타겟 스레드 유형**: 시각 프로그래밍 가치 토론, "Why don't we use visual programming more?" 류의 스레드

**댓글 초안**:

---

> The "overhyped" camp usually points to scaling issues — and I think that's fair. But the nuance I'd add is that visual tools work well when the domain has a clear node semantics.

I've been building **FreeWire**, a node editor for a component DSL. The use case is UI component wiring — connecting a "Button" node to a "State" node to an "Output" node. In that domain, the graph *is* the mental model, so the visual representation actually reduces cognitive load rather than adding it.

Where I ran into the scaling problem was when the graph got large enough that "intent" became ambiguous — a node could be connected correctly by structure but still be semantically wrong. That pushed me toward adding what I call an Intent Engine: each node declares what it's *for* (not just what it *takes*), and the tool enforces that contract.

Whether that's the right abstraction I'm not sure, but it did make the tool more useful than a pure visual-syntax-sugar approach. The key insight was: visual programming helps when the connections *carry meaning*, not just when they carry data.

FreeWire is open on GOGS if you want to dig into how the contract/goal system is implemented: https://gogs.dclub.kr/kim/FreeWire

---

**글자 수**: 약 950자 (영문) / 토론 참여형, FreeWire는 근거 제시 형태로 언급

---

## 3. FAQ 초안 (별도 파일 참조)

→ `/data/data/com.termux/files/home/ai-marketing-team/community/faq-freewire.md`

---

## 콘텐츠 검토 (content-policy.md 기준)

| 항목 | 확인 | 비고 |
|------|------|------|
| 사실 검증 | ✅ | 162 테스트, Phase 1-6 실제 구현 수치 |
| 코드 동작 | ✅ | README 기반 실행 명령어 검증 |
| 경쟁사 존중 | ✅ | Blueprints/Node-RED 언급은 비방 아닌 분류 |
| 출처 명시 | ✅ | GOGS 링크 명시 |
| 개인정보 제외 | ✅ | 없음 |
| 톤 검증 | ✅ | 공유 중심, 홍보 자제, 진정성 유지 |
| 링크 유효성 | ⚠️ | GOGS 링크 게시 전 확인 필요 |

---

*작성 완료: 2026-03-18T10:45 KST*
