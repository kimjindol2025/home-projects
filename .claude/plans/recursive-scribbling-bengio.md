# ClaudeScript Phase 5: Runtime Interpreter 계획 (보류)

## Context
Phase 5는 ClaudeScript VT 바이트코드를 실제 실행하는 인터프리터 구현 단계입니다.
현재 v4 GOGS 저장소 작업을 먼저 처리 후 재개 예정.

## 조사 결과 요약
- VT 바이트코드 형식: S-expression (let, defn, if, while, match, try-catch 등)
- freelang-vm (Rust): 스택 기반 VM, GC 포함 재활용 가능
- freelang-final (JS): evaluator.js 스코프 체인, runtime.js 150개+ 내장함수 재활용 가능

## 구현 예정 (Phase 5)
1. src/interpreter.ts - VT S-expression 파서 + 실행기
2. src/runtime.ts - 내장 함수 (println, print, length, push, pop 등)
3. src/scope.ts - Environment 스코프 체인
4. tests/interpreter.test.ts - 실행 테스트 20개
