# Phase 6.5 — Runtime Optimization Layer

## 컨텍스트

**목표**: Phase 5.5.4 (제네릭 TypeVar) 완료 후, Phase 6.5에서 런타임 최적화 레이어를 구현한다.
현재 최적화 (컴파일 타임: 상수 폴딩, 중복 제거, Typed Opcode)를 보완하는 3가지 최적화:
1. **Peephole Optimizer** — 바이트코드 패턴 최적화 (컴파일 후 패스)
2. **DCE Optimizer** — Dead Code Elimination (도달 불가능한 코드 제거)
3. **Runtime Profiler** — 핫 경로 추적 (opt-in 프로파일링)

**기존 구현 (재사용 참고)**:
- `compiler.ts:5-29`: Op enum (55개 opcode, 인자 있는 opcode 목록)
- `compiler.ts:43-47`: Chunk = `{ code: number[], constants: Value[], lines: number[] }`
- `jump-table.ts`: OffsetCalculator → 점프 주소 계산 패턴 참고
- `index.ts:32-84`: `run()` 함수 — 두 경로 (빠른/모듈) 모두 수정 필요

---

## 수정 대상 파일

| 파일 | 작업 |
|------|------|
| `src/peephole-optimizer.ts` | **신규**: PeepholeOptimizer 클래스 |
| `src/dce-optimizer.ts` | **신규**: DCEOptimizer 클래스 |
| `src/runtime-profiler.ts` | **신규**: RuntimeProfiler 클래스 |
| `src/vm.ts` | **수정**: profiler opt-in 파라미터 추가 (3곳) |
| `src/index.ts` | **수정**: compile() 후 최적화 패스 추가 (Line 47-49, 81-82 두 경로) |
| `src/vm-runtime-optimizer-test.ts` | **신규**: 15개 테스트 |

---

## Step 1: peephole-optimizer.ts 신규 생성

```typescript
import { Op, Chunk } from "./compiler";

interface Instruction {
  op: Op;
  args: number[];
  pos: number;  // 원래 chunk.code 내 시작 위치 (점프 주소 재계산용)
}

export class PeepholeOptimizer {
  // 인자 있는 opcode → 인자 개수 (arity)
  private static readonly ARITIES = new Map<number, number>([
    [Op.Const, 1], [Op.Load, 1], [Op.Store, 1],
    [Op.LoadGlobal, 1], [Op.StoreGlobal, 1],
    [Op.Jump, 1], [Op.JumpIfFalse, 1], [Op.JumpIfTrue, 1],
    [Op.Call, 1], [Op.NewArray, 1], [Op.NewObject, 1],
    [Op.GetProp, 1], [Op.TryBegin, 1], [Op.NewClosure, 1],
    [Op.LoadFree, 1], [Op.StoreFree, 1], [Op.Builtin, 1],
  ]);

  private stats = { removedCount: 0, passCount: 0 };

  // Fix-point 반복 (변화가 없을 때까지, 최대 10회)
  optimize(chunk: Chunk): Chunk { ... }

  // chunk.code[] → Instruction[] 디코딩
  private decode(chunk: Chunk): Instruction[] { ... }

  // Instruction[] → chunk.code[] 인코딩 + 점프 타깃 재계산
  private encode(instructions: Instruction[], original: Chunk): Chunk {
    // old_pos → new_code_index 맵 구성
    // Jump/JumpIfFalse/JumpIfTrue 인자를 맵으로 변환
  }

  // 7가지 패턴 적용
  private applyPatterns(instructions: Instruction[], chunk: Chunk):
    { result: Instruction[]; changed: boolean } {
    // 패턴 1: Const(비함수) + Pop → 제거
    //   chunk.constants[idx].tag !== "fn" 체크 필수
    // 패턴 2: Not + Not → 제거 (이중 부정)
    // 패턴 3: LoadGlobal + Pop → 제거
    // 패턴 4: Jump(addr == next_pos) → 제거 (NOP 점프)
    // 패턴 5: Const(true) + JumpIfFalse → 제거 (항상 진행)
    // 패턴 6: Const(false) + JumpIfFalse → unconditional Jump
    // 패턴 7: JumpIfFalse(addr) + Jump(addr) → Pop + Jump(addr)
  }

  getStats() { return { ...this.stats }; }
}
```

**핵심 위험**: 점프 타깃 재계산 (encode 단계)
```
원본: [0]Const [2]Not [3]Not [4]JumpIfFalse→8 [6]Print [7]Jump→9 [9]Halt
Not(3) 제거 후: [0]Const [2]Not [3]JumpIfFalse→? [5]Print [6]Jump→? [8]Halt
→ old_pos→new_pos 맵: {8→7, 9→8} 로 점프 인자 재계산
```

---

## Step 2: dce-optimizer.ts 신규 생성

```typescript
import { Op, Chunk } from "./compiler";

export class DCEOptimizer {
  private stats = { removedCount: 0 };

  optimize(chunk: Chunk): Chunk {
    const reachable = this.computeReachable(chunk);
    // reachable에 없는 위치 제거 + 점프 타깃 재계산
  }

  // BFS로 도달 가능한 바이트코드 위치 집합 계산
  private computeReachable(chunk: Chunk): Set<number> {
    // worklist = [0]
    // Jump → target으로만 이동 (fall-through 없음)
    // JumpIfFalse/JumpIfTrue → fall-through + target 두 경로
    // Halt/Return → 종료 (이후 코드 도달 불가)
    // 이미 방문한 위치는 스킵 (무한 루프 방지)
  }

  getStats() { return { ...this.stats }; }
}
```

---

## Step 3: runtime-profiler.ts 신규 생성

```typescript
import { Op } from "./compiler";

export interface ProfileStats {
  totalOps: number;
  opFrequency: Map<Op, number>;
  hotFunctions: Array<{ name: string; callCount: number; totalOps: number }>;
}

export class RuntimeProfiler {
  private totalOps = 0;
  private opFrequency = new Map<Op, number>();
  private fnCallCounts = new Map<string, number>();
  private fnOpCounts = new Map<string, number>();
  private enabled = false;

  constructor(private options = { hotThreshold: 10 }) {}

  enable(): void { this.enabled = true; }
  disable(): void { this.enabled = false; }
  isEnabled(): boolean { return this.enabled; }

  recordOp(op: Op): void {
    if (!this.enabled) return;
    this.totalOps++;
    this.opFrequency.set(op, (this.opFrequency.get(op) ?? 0) + 1);
  }

  recordCall(fnName: string): void {
    if (!this.enabled) return;
    this.fnCallCounts.set(fnName, (this.fnCallCounts.get(fnName) ?? 0) + 1);
  }

  getStats(): ProfileStats {
    const hotFunctions = Array.from(this.fnCallCounts.entries())
      .filter(([, count]) => count >= this.options.hotThreshold)
      .map(([name, callCount]) => ({
        name, callCount, totalOps: this.fnOpCounts.get(name) ?? 0
      }))
      .sort((a, b) => b.callCount - a.callCount);
    return { totalOps: this.totalOps, opFrequency: new Map(this.opFrequency), hotFunctions };
  }

  reset(): void { ... }
}
```

---

## Step 4: vm.ts 수정 (3곳)

### 4-1: import + 생성자 파라미터 추가 (Line 63)
```typescript
import { RuntimeProfiler } from "./runtime-profiler";
// ...
constructor(private chunk: Chunk, initialGlobals?: Map<string, Value>, private profiler?: RuntimeProfiler) {
```

### 4-2: execOp() 첫 줄에 추가 (Line 340 직후)
```typescript
private execOp(op: number) {
  this.profiler?.recordOp(op as Op);  // Phase 6.5 추가
  // ... 기존 코드
```

### 4-3: Op.Call 처리의 fn 분기에 추가
```typescript
if (callee.tag === "fn") {
  this.profiler?.recordCall(callee.name ?? "anonymous");  // Phase 6.5 추가
  // ... 기존 프레임 push 코드
```

---

## Step 5: index.ts 수정 (두 경로 모두)

```typescript
// import 추가
import { PeepholeOptimizer } from "./peephole-optimizer";
import { DCEOptimizer } from "./dce-optimizer";

// 헬퍼 함수 (중복 제거)
function applyOptimizations(chunk: Chunk): Chunk {
  return new DCEOptimizer().optimize(new PeepholeOptimizer().optimize(chunk));
}

// 경로 1 (Line 47-49):
const chunk = applyOptimizations(compile(ast, ...));

// 경로 2 (Line 81-82):
const chunk = applyOptimizations(compile({ stmts: allStmts }, ...));
```

---

## Step 6: vm-runtime-optimizer-test.ts 신규 생성 (15개)

```
Suite A: Peephole Optimizer (5개)
  A-1: Not+Not 이중 부정 → 코드 길이 축소 확인
  A-2: Const+Pop 패턴 → 제거 + 실행 결과 동일
  A-3: LoadGlobal+Pop → 제거
  A-4: Const(true)+JumpIfFalse → 분기 제거
  A-5: run() end-to-end 최적화 검증

Suite B: DCE Optimizer (5개)
  B-1: Halt 이후 dead code 제거
  B-2: Return 이후 dead code 제거 (함수 내)
  B-3: 도달 불가능 코드 제거 확인
  B-4: DCE 통계 removedCount 기록
  B-5: Peephole + DCE 연쇄 최적화

Suite C: RuntimeProfiler + 통합 (5개)
  C-1: totalOps 카운팅
  C-2: opFrequency - Print 빈도 확인
  C-3: 핫 함수 탐지 (threshold=3, 4회 호출)
  C-4: Profiler disabled 시 통계 없음
  C-5: 전체 파이프라인 통합 (factorial(5)=120, hot 함수 탐지)
```

---

## 검증 순서

```bash
npm run build
node dist/vm-runtime-optimizer-test.js   # 신규 15개
bash /tmp/run_all_tests.sh               # 전체 192개 (177 + 15)
```

목표: **192/192 전체 통과**

---

## 구현 순서

```
1. peephole-optimizer.ts — 외부 의존 없음 (Op, Chunk만)
2. dce-optimizer.ts — 외부 의존 없음 (Op, Chunk만)
3. runtime-profiler.ts — 외부 의존 없음 (Op만)
4. vm.ts — profiler 파라미터 추가 (3곳, 선택적)
5. index.ts — 최적화 패스 삽입 (두 경로)
6. vm-runtime-optimizer-test.ts — 15개 테스트
7. npm run build + 전체 테스트
```

---

## 리스크 통제

| 리스크 | 대책 |
|--------|------|
| Const(fn)+Pop 제거 → 함수 정의 손실 | `chunk.constants[idx].tag === "fn"` 체크 |
| 점프 타깃 주소 재계산 실패 | old_pos → new_pos 맵으로 모든 Jump 인자 재계산 |
| vm.ts 생성자 변경 → 기존 테스트 깨짐 | profiler는 선택적(`?`) 파라미터 유지 |
| index.ts 경로 중 하나 누락 | 두 경로(Line 47-49, Line 81-82) 모두 적용 |
| 기존 177개 테스트 회귀 | run_all_tests.sh 전체 검증으로 확인 |
