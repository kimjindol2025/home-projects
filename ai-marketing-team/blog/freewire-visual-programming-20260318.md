# FreeWire — 노드 기반 시각 프로그래밍으로 UI를 만들다

> **요약**: FreeWire는 노드를 연결하는 것만으로 FreeLang UI 코드를 자동 생성하는 시각 프로그래밍 환경입니다. AI 자동 생성, Intent 계약 검사, Undo/Redo까지 — 코드를 "그리는" 새로운 방식을 소개합니다.

**작성일**: 2026-03-18 | **카테고리**: 도구 소개, 시각 프로그래밍 | **읽기 시간**: 약 8분

---

## 프로그래밍을 "그린다"는 것

코드를 작성할 때 우리는 보통 텍스트 에디터를 엽니다. 함수를 정의하고, 변수를 선언하고, 로직을 한 줄씩 씁니다. 이 방식은 수십 년간 검증된 훌륭한 방법입니다.

하지만 UI 컴포넌트를 만들 때는 조금 다른 고민이 생겨요. "이 버튼이 클릭되면 저 상태가 바뀌고, 그 상태가 저 텍스트에 반영되고..." — 이런 데이터 흐름과 이벤트 연결을 머릿속에서 그리면서 코드로 옮기는 과정이 꽤 번거롭습니다.

**FreeWire**는 이 "머릿속 그림"을 그대로 화면에 그릴 수 있게 해줍니다.

언리얼 엔진의 Blueprints, Node-RED의 플로우 방식처럼 — 노드를 캔버스에 배치하고 연결하면 FreeLang `.free` 코드가 자동으로 완성됩니다. 단순한 시각화 도구가 아닙니다. 노드가 "왜 존재하는지"를 이해하고, 실수를 기억하며, 스스로 고칩니다.

---

## FreeWire의 세 가지 노드 세계

FreeWire의 그래프는 세 종류의 노드로 구성됩니다. 이 구분은 UI 개발의 세 가지 관심사 — **표현**, **데이터**, **로직** — 를 그대로 반영한 것입니다.

### component 노드: 화면에 보이는 것들

```json
{
  "id": "btn-1",
  "type": "component",
  "subtype": "button",
  "props": {
    "label": "카운트 증가",
    "variant": "primary"
  }
}
```

버튼, 텍스트, 입력창, 카드, 컨테이너, 이미지, 배지 — 총 7가지 UI 요소를 노드로 표현합니다. `variant`처럼 세부 스타일도 props로 지정할 수 있어요.

### data 노드: 흐르는 정보들

```json
{
  "id": "state-count",
  "type": "data",
  "subtype": "state",
  "props": {
    "name": "count",
    "initial": "0"
  }
}
```

상태(state), 상수(const), 배열(array), 외부 props까지 4가지 데이터 노드가 있습니다. `state` 노드는 초기값을 설정하면, 나중에 `setState` 로직 노드와 연결해 값을 업데이트할 수 있습니다.

### logic 노드: 일어나는 일들

```json
{
  "id": "handler-click",
  "type": "logic",
  "subtype": "handler",
  "props": {
    "name": "handleClick"
  }
}
```

이벤트 핸들러(handler), 상태 업데이트(setState), 조건 분기(ifelse) — 세 가지 로직 노드가 흐름을 제어합니다. handler 노드를 button 노드의 `toPort: 1`에 연결하면 클릭 이벤트가 연결됩니다.

### 와이어: 노드를 잇는 의도

노드들은 **와이어(Wire)**로 연결됩니다. 와이어 하나가 데이터 흐름 또는 이벤트 바인딩 하나를 의미합니다.

```json
[
  { "from": "state-count", "fromPort": 0, "to": "text-display", "toPort": 0 },
  { "from": "handler-click", "fromPort": 0, "to": "btn-1", "toPort": 1 }
]
```

state → text 연결은 "상태가 텍스트에 표시된다"는 의미, handler → button 연결은 "핸들러가 버튼 클릭에 응답한다"는 의미입니다. 그래프 자체가 UI의 설계 의도를 담습니다.

---

## 그래프가 코드가 된다

노드와 와이어를 연결하고 나면 **GraphCompiler**가 FreeLang Light `.free` 코드를 자동 생성합니다. 위 예시의 카운터 그래프는 이런 코드로 변환됩니다.

```
component Counter {

  script {
    let count = 0;
    function set_count(v) { count = v; }
    function handleClick() { set_count(count + 1); }
  }

  markup {
    text(content: "카운트", tag: "h2")
    button(label: "카운트 증가", variant: "primary")
  }
}
```

`script` 블록에는 상태와 핸들러가, `markup` 블록에는 UI 구조가 들어갑니다. 컴파일된 코드는 바로 **FreeExecutor**에 전달되어 HTML/CSS/JS로 변환되고, 브라우저에서 Live Preview로 확인할 수 있습니다.

서버 없이 브라우저 단독으로 실행되는 것도 장점 중 하나입니다.

```bash
# 브라우저 에디터 시작 (서버 불필요)
python3 -m http.server 5500 --directory www
# http://localhost:5500 접속
```

---

## "설명만 하면 그래프가 완성된다" — AI 자동 생성

노드를 손으로 하나씩 배치하는 것도 충분히 직관적이지만, FreeWire는 한 단계 더 나아갑니다.

**자연어로 UI를 설명하면 노드 그래프가 자동 생성됩니다.**

API 서버를 실행하고 `/api/ai-generate` 엔드포인트를 호출해 보세요.

```bash
curl -X POST http://localhost:4010/api/ai-generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "이메일과 비밀번호 입력칸, 로그인 버튼이 있는 로그인 폼"}'
```

응답으로는 버튼, 입력창, 컨테이너, 핸들러, 상태 노드들이 포함된 완성된 그래프 JSON이 돌아옵니다. AI가 "키워드 매칭"이 아니라 의미를 이해해서 노드 구조를 설계하기 때문에, 복잡한 UI 설명도 비교적 정확하게 변환됩니다.

15개 노드 타입(component 7개, data 4개, logic 3개 + 조합)을 지원하며, 생성된 노드들은 자동으로 캔버스에 배치됩니다.

AI 생성 기능은 Claude API 키 없이도 목(mock) 데이터로 테스트할 수 있어요.

```bash
# API 키 없이도 테스트 가능
node src/test/test-phase5.js   # 35/35 테스트 통과
```

---

## Undo/Redo: 되돌리는 것도 프로그래밍이다

이번 Phase 9에서 새롭게 추가된 기능이 있습니다 — **GraphHistory Undo/Redo 시스템**입니다.

그래프 에디터를 쓰다 보면 "방금 전 연결을 취소하고 싶다"거나 "아까 배치가 더 나았는데" 하는 순간이 반드시 옵니다. FreeWire의 GraphHistory는 이 문제를 스냅샷 기반으로 해결합니다.

### 작동 원리

GraphHistory는 그래프 상태(JSON)를 두 개의 스택으로 관리합니다.

- **past 스택**: 저장된 이전 상태들 (최대 50개)
- **future 스택**: undo로 되돌린 상태들 (redo에 사용)

```javascript
const { GraphHistory } = require('./src/core/graph-history');

const history = new GraphHistory(50);  // 최대 50단계

// 노드 추가 후 저장
history.push(graph.toJSON());

// 실수했다면 되돌리기
const prev = history.undo();   // past.pop() → future.push()
graph.loadFromJSON(prev);

// 역시 원래 게 더 나았다면 다시 앞으로
const next = history.redo();   // future.pop() → past.push()
graph.loadFromJSON(next);

// 현재 상태 확인
console.log(history.canUndo);  // true / false
console.log(history.canRedo);  // true / false
console.log(history.size);     // 현재 저장된 스냅샷 수
```

### 설계 원칙

`push()`를 호출하면 future 스택이 비워집니다. "새 길을 선택하면 이전에 되돌아갔던 기록은 사라진다"는 직관적인 규칙입니다. 대부분의 텍스트 에디터 Undo/Redo가 이 방식을 씁니다.

스냅샷은 JSON 문자열로 직렬화해서 저장하기 때문에, 그래프가 복잡해져도 메모리 사용이 예측 가능합니다. 50단계 한도를 초과하면 가장 오래된 항목이 자동으로 제거됩니다.

작지만 의미 있는 기능입니다. 실험적으로 노드를 연결해보고, 마음에 안 들면 되돌리고 — 이 자유로움이 시각 프로그래밍의 진입 장벽을 낮춰줍니다.

---

## Intent 시스템: 노드의 "왜"를 안다

FreeWire가 단순한 노드 에디터와 다른 점은 **Intent 시스템**에 있습니다.

각 노드는 `requires / ensures / produces / warns`로 구성된 **계약(Contract)**을 가집니다.

```javascript
// button 노드의 Contract
{
  produces: ['ui.button.clickable'],
  warns:    ['ui.button.connected'],  // 핸들러 미연결 시 경고
}

// handler 노드의 Contract
{
  requires: ['action.triggered'],
  produces: ['action.handler.present', 'ui.button.connected'],
}
```

IntentEvaluator는 그래프 전체를 보면서 계약 충족도를 **0-100점**으로 평가합니다. 버튼 노드에 핸들러가 연결되지 않으면 경고(`ui.button.connected` 미충족)를 발생시키고, RepairEngine이 자동으로 패치를 생성해 적용합니다.

더 나아가, 수정 이력이 localStorage에 쌓이면 PatternAdvisor가 "이 패턴에서 이런 실수가 자주 발생한다"고 **사전 경고**를 제공합니다. 학습하는 에디터입니다.

---

## 테스트로 검증된 235개의 약속

FreeWire의 각 Phase는 독립적인 테스트로 검증됩니다.

| Phase | 내용 | 테스트 |
|-------|------|--------|
| 1 | 노드 JSON 스펙 | 32개 |
| 2 | GraphCompiler | 16개 |
| 3 | API 서버 + FreeExecutor | 24개 |
| 4 | 브라우저 에디터 UI | 30개 |
| 5 | AI 자동 생성 | 35개 |
| 6 | Live Preview | 25개 |
| 7 | Intent 시스템 | 35개 |
| 8 | Adaptive Learning | 38개 |
| **합계** | | **235/235** |

모든 테스트가 통과됩니다. GraphHistory Undo/Redo(Phase 9)는 현재 추가 중입니다.

테스트를 직접 실행해보고 싶다면:

```bash
git clone https://gogs.dclub.kr/kim/FreeWire
cd FreeWire
npm install
node src/test/test-phase1.js
```

---

## 다음 단계: Phase 6 WASM

현재 FreeWire는 Node.js 서버와 브라우저 에디터로 동작합니다. 다음 목표는 **WebAssembly(WASM) 빌드**입니다.

서버 없이 브라우저 단독으로 전체 파이프라인(노드 → FreeLang 컴파일 → HTML 렌더링)이 돌아가면, FreeWire는 설치 없이 URL 하나로 공유할 수 있는 도구가 됩니다.

---

## 마무리 & 피드백 요청

FreeWire는 "코드 작성 방식을 바꿀 수 있을까?"라는 질문에서 시작했습니다.

아직 초기 단계이고, 텍스트 코딩을 완전히 대체하려는 것도 아닙니다. 오히려 "UI 구조를 시각적으로 설계하고, 그 설계에서 코드가 나오는" 워크플로우를 탐구하는 실험적인 프로젝트입니다.

한계도 솔직하게 말씀드리면:
- 현재 지원 노드 타입이 15종으로 제한적입니다
- 복잡한 비즈니스 로직은 여전히 직접 코드 작성이 낫습니다
- WASM 빌드가 완성되기 전까지는 로컬 서버가 필요합니다

그럼에도 "버튼 클릭 → 상태 업데이트 → 텍스트 반영" 같은 흔한 UI 패턴은 노드 몇 개를 연결하는 것만으로 완성됩니다. 직접 써보시고 피드백을 주시면 큰 도움이 됩니다.

저장소: [gogs.dclub.kr/kim/FreeWire](https://gogs.dclub.kr/kim/FreeWire)

다음 글: **"FreeLang Light — 브라우저에서 바로 실행되는 프론트엔드 언어"**

---

*이 글은 FreeLang Marketing Ops Content Writer가 작성했습니다. 코드 예시는 실제 FreeWire 저장소 기준입니다.*
