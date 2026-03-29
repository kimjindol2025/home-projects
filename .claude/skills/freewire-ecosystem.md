---
name: freewire-ecosystem
description: FreeWire 프로젝트에서 FreeLang 생태계 조성 원칙 적용 — 새 기능/UI 추가 시 .free 코드 우선, GraphCompiler/FreeExecutor 확장 가이드
---

# FreeLang 생태계 조성 원칙 적용

사용자가 FreeWire 프로젝트에 새 기능을 추가하거나 UI 패턴을 구현할 때:

## 판단 흐름

1. **FreeLang(.free)으로 표현 가능한가?**
   - YES → `GraphCompiler._compileNode()`에 case 추가하여 .free 코드 생성
   - NO → 다음 단계로

2. **FreeExecutor에 런타임 추가 가능한가?**
   - YES → `free-executor.js` `_render*()` 또는 핸들러 추가
   - NO → 다음 단계로

3. **그래도 불가한 경우만** → 직접 HTML/JS 생성 (최후 수단)

## 새 노드 subtype 추가 체크리스트

```
[ ] src/core/node-types-index.js — 타입 정의 (inputs/outputs/props)
[ ] src/compiler/graph-compiler.js — _compileNode() case 추가
[ ] src/intent/intent-spec.js — NODE_INTENT_DEFAULTS 계약 추가
[ ] src/intent/evaluator.js — 필요 시 구조 검사 추가
[ ] src/test/test-phase-new.js — 테스트 작성
[ ] run-tests.js — 새 테스트 파일 등록
```

## .free 코드 생성 패턴

```
// 컴포넌트 노드
component <id> type="<subtype>" {
  props { label = "<value>" }
}

// 상태 노드
state <id> { initial = <value> }

// 핸들러 노드
handler <id> { action = "<type>" }

// 와이어 연결
wire <from> -> <to>
```

## 금지 패턴 (FreeLang 원칙 위반)

```javascript
// ❌ 직접 HTML 생성 (FreeExecutor 우회)
return `<button onclick="${handler}">${label}</button>`;

// ✅ GraphCompiler 통해 .free 코드로 출력
_compileComponent(node) {
  return `component ${node.id} type="button" {\n  props { label = "${label}" }\n}`;
}
```

## 실행 확인

```bash
# 원칙 준수 확인: 파이프라인 전체 테스트
node pipeline-test.js "새 UI 이름"

# .free 코드 출력 확인 (③ 컴파일 단계)
# 반드시 .free 코드에 새 노드가 포함되어야 함

# 테스트 유지
node src/test/run-tests.js  # 모두 PASS 확인
```
