// 🔄 Layer 2: Semantic Sync Protocol - Deterministic Execution
// Unforgiving Rule: Semantic Equivalence 1.0 (identical output for identical input)

import {
  StateSnapshot,
  SemanticEquivalenceProof,
  sha256,
  CONSTANTS,
  measurePerformance,
} from './types';
import { RDMAFabric } from './rdma_fabric';
import { HashChain } from './hash_chain';

/**
 * ExecutionState: 단일 노드의 실행 상태
 * Deterministic하게 추적하여 다른 노드와 비교 가능
 */
class ExecutionState {
  executionId: bigint;
  nodeId: bigint;
  memoryImage: Buffer;
  instructionCounter: number = 0;
  variables: Map<string, string> = new Map();
  outputs: string[] = [];

  constructor(executionId: bigint, nodeId: bigint) {
    this.executionId = executionId;
    this.nodeId = nodeId;
    this.memoryImage = Buffer.alloc(CONSTANTS.MAX_MEMORY_SIZE);
  }

  /**
   * 변수 설정 (Deterministic)
   */
  setVariable(name: string, value: string): void {
    this.variables.set(name, value);
    this.instructionCounter++;
  }

  /**
   * 변수 조회
   */
  getVariable(name: string): string | undefined {
    return this.variables.get(name);
  }

  /**
   * 메모리 쓰기
   */
  writeMemory(offset: number, data: Buffer): void {
    if (offset + data.length > this.memoryImage.length) {
      throw new Error(`Memory write out of bounds: ${offset}+${data.length}`);
    }
    data.copy(this.memoryImage, offset);
    this.instructionCounter++;
  }

  /**
   * 메모리 읽기
   */
  readMemory(offset: number, size: number): Buffer {
    return this.memoryImage.slice(offset, offset + size);
  }

  /**
   * 출력 기록
   */
  output(value: string): void {
    this.outputs.push(value);
    this.instructionCounter++;
  }

  /**
   * 상태 스냅샷 생성
   */
  snapshot(): StateSnapshot {
    const variableState = new Map<string, string>();
    for (const [key, value] of this.variables.entries()) {
      variableState.set(key, value);
    }

    return {
      nodeId: this.nodeId,
      executionId: this.executionId,
      timestamp: BigInt(Date.now()) * 1_000_000n,
      memoryChecksum: sha256(this.memoryImage),
      instructionCounter: this.instructionCounter,
      variableState,
      hashChainLink: '', // Layer 3에서 설정
    };
  }

  /**
   * 체크섬 계산
   */
  calculateChecksum(): string {
    const stateJson = JSON.stringify({
      variables: Array.from(this.variables.entries()),
      outputs: this.outputs,
      instructionCounter: this.instructionCounter,
    });
    return sha256(stateJson);
  }
}

/**
 * SemanticSyncProtocol: Layer 2
 *
 * 의미 동등성 보장:
 *   1. 모든 노드가 동일한 코드 실행
 *   2. 동일한 입력 → 동일한 출력
 *   3. 성능 차이는 무관 (절대 정확성만 중요)
 *   4. 10,000개 노드 전체 동기화
 *
 * Unforgiving Rules:
 *   - Semantic Equivalence: 1.0 (100%)
 *   - All outputs identical
 *   - All memory checksums identical
 *   - Instruction counter identical
 *   - 0 semantic desync events
 */
export class SemanticSyncProtocol {
  readonly nodeId: bigint;
  private executionId: bigint = 1n;
  private fabric: RDMAFabric;
  private currentState: ExecutionState | null = null;
  private stateSnapshots: StateSnapshot[] = [];
  private syncLog: HashChain;
  private otherNodeSnapshots: Map<bigint, StateSnapshot> = new Map();

  constructor(nodeId: bigint, fabric: RDMAFabric) {
    this.nodeId = nodeId;
    this.fabric = fabric;
    this.syncLog = new HashChain();

    console.log(`[LAYER2] SemanticSyncProtocol initialized on node ${nodeId}`);
  }

  /**
   * 새로운 실행 시작
   */
  public async startExecution(
    code: string,
    input: Map<string, string>
  ): Promise<ExecutionState> {
    const { result: state } = await measurePerformance(`SemanticSync.startExecution on node ${this.nodeId}`, async () => {
      const newState = new ExecutionState(this.executionId++, this.nodeId);

      // 입력 변수 초기화
      for (const [key, value] of input.entries()) {
        newState.setVariable(key, value);
      }

      // 코드 실행 (모의 구현)
      await this.executeCode(newState, code);

      this.currentState = newState;
      return newState;
    });

    return state;
  }

  /**
   * 코드 실행 (모의 구현)
   * 실제 구현에서는 FreeLang JIT을 사용
   */
  private async executeCode(state: ExecutionState, code: string): Promise<void> {
    // 간단한 모의 실행
    // 실제: FreeLang parser + JIT compiler + execution

    // 예: code = "x = input; y = x + 1; print(y)"
    const lines = code.split(';').map((s) => s.trim());

    for (const line of lines) {
      if (line.startsWith('x =')) {
        const value = state.getVariable('input') || '0';
        state.setVariable('x', value);
      } else if (line.startsWith('y =')) {
        const x = parseInt(state.getVariable('x') || '0', 10);
        state.setVariable('y', String(x + 1));
      } else if (line.startsWith('print')) {
        const value = state.getVariable('y') || '';
        state.output(value);
      }
    }

    // Deterministic delay (고정)
    await new Promise((resolve) => setTimeout(resolve, 1));
  }

  /**
   * 현재 상태 스냅샷 생성
   */
  public async createSnapshot(): Promise<StateSnapshot> {
    const { result: snapshot } = await measurePerformance(`SemanticSync.createSnapshot on node ${this.nodeId}`, async () => {
      if (!this.currentState) {
        throw new Error('No execution in progress');
      }

      const snap = this.currentState.snapshot();

      // 동기화 로그에 기록
      this.syncLog.addLink(snap.executionId, 'state_snapshot', snap.memoryChecksum);

      this.stateSnapshots.push(snap);

      return snap;
    });

    return snapshot;
  }

  /**
   * 다른 노드와의 동등성 검증
   * @param otherNodeId 비교 대상 노드
   * @param otherSnapshot 비교 대상 스냅샷
   * @returns 동등성 증명
   */
  public async verifyEquivalence(
    otherNodeId: bigint,
    otherSnapshot: StateSnapshot
  ): Promise<SemanticEquivalenceProof> {
    const { result: proof } = await measurePerformance(`SemanticSync.verifyEquivalence node ${this.nodeId} vs ${otherNodeId}`, async () => {
      if (!this.currentState) {
        throw new Error('No execution in progress');
      }

      const proofObj: SemanticEquivalenceProof = {
        node1Id: this.nodeId,
        node2Id: otherNodeId,
        inputHash: 'todo', // 실제로는 입력 데이터 해시
        outputHash: 'todo',
        snapshot1: this.currentState.snapshot(),
        snapshot2: otherSnapshot,
        isEquivalent: false,
        proofTime: BigInt(Date.now()) * 1_000_000n,
      };

      // 동등성 검증
      const equiv =
        proofObj.snapshot1.memoryChecksum === proofObj.snapshot2.memoryChecksum &&
        proofObj.snapshot1.instructionCounter === proofObj.snapshot2.instructionCounter;

      proofObj.isEquivalent = equiv;

      if (equiv) {
        // 동기화 로그에 기록
        this.syncLog.addLink(
          this.executionId,
          'equivalence_verified',
          sha256(
            JSON.stringify(proofObj, (key, value) =>
              typeof value === 'bigint' ? value.toString() : value
            )
          )
        );
      }

      return proofObj;
    });

    return proof;
  }

  /**
   * 모든 노드와의 동등성 동시 검증
   * @param nodeSnapshots 다른 노드들의 스냅샷
   */
  public async verifyAllEquivalence(
    nodeSnapshots: Map<bigint, StateSnapshot>
  ): Promise<{ allEquivalent: boolean; proofs: SemanticEquivalenceProof[] }> {
    const { result, elapsedNs } = await measurePerformance('SemanticSync.verifyAllEquivalence', async () => {
      const proofs: SemanticEquivalenceProof[] = [];

      for (const [nodeId, snapshot] of nodeSnapshots.entries()) {
        if (nodeId === this.nodeId) continue;

        const proof = await this.verifyEquivalence(nodeId, snapshot);
        proofs.push(proof);
      }

      const allEquivalent = proofs.every((p) => p.isEquivalent);

      return { allEquivalent, proofs };
    });

    const latencyMs = Number(elapsedNs) / 1_000_000;
    console.log(`[LAYER2] All equivalence verified in ${latencyMs.toFixed(2)}ms`);

    return result;
  }

  /**
   * 동기화 로그 검증
   */
  public async verifySyncLog(): Promise<boolean> {
    const verification = await this.syncLog.verify();
    return verification.isValid;
  }

  /**
   * 크로스-노드 의미 동등성 검증 (완전)
   *
   * 모든 노드의 스냅샷을 수집하여 동일성 검증
   * Unforgiving Rule: 100% 동등성 필요
   */
  public async verifyGlobalSemanticConsistency(
    allNodeSnapshots: Map<bigint, StateSnapshot>
  ): Promise<{
    isConsistent: boolean;
    equivalentGroups: StateSnapshot[][];
    inconsistentNodes: bigint[];
    proofTime: bigint;
  }> {
    const { result } = await measurePerformance(`SemanticSync.verifyGlobalSemanticConsistency ${allNodeSnapshots.size} nodes`, async () => {
      const startTime = BigInt(Date.now()) * 1_000_000n;

      // 체크섬으로 그룹핑
      const groups = new Map<string, StateSnapshot[]>();

      for (const [nodeId, snapshot] of allNodeSnapshots.entries()) {
        const key = snapshot.memoryChecksum;
        if (!groups.has(key)) {
          groups.set(key, []);
        }
        groups.get(key)!.push(snapshot);
      }

      // 동등성 검증: 모든 노드가 동일 그룹에 속해야 함
      const equivalentGroups = Array.from(groups.values());
      const isConsistent = equivalentGroups.length === 1;

      // 불일치 노드 식별
      const inconsistentNodes: bigint[] = [];
      if (!isConsistent) {
        const majorityGroup = equivalentGroups[0];
        const majorityHashes = new Set(majorityGroup.map((s) => s.memoryChecksum));

        for (const [nodeId, snapshot] of allNodeSnapshots.entries()) {
          if (!majorityHashes.has(snapshot.memoryChecksum)) {
            inconsistentNodes.push(nodeId);
          }
        }
      }

      const proofTime = BigInt(Date.now()) * 1_000_000n - startTime;

      return {
        isConsistent,
        equivalentGroups,
        inconsistentNodes,
        proofTime,
      };
    });

    return result;
  }

  /**
   * Layer 2 통계
   */
  public async getStats(): Promise<{
    nodeId: bigint;
    executionsCompleted: number;
    snapshotsTaken: number;
    syncLogValid: boolean;
    currentStateChecksum: string;
  }> {
    const syncLogValid = await this.verifySyncLog();
    return {
      nodeId: this.nodeId,
      executionsCompleted: Number(this.executionId - 1n),
      snapshotsTaken: this.stateSnapshots.length,
      syncLogValid,
      currentStateChecksum: this.currentState?.calculateChecksum() || 'N/A',
    };
  }

  /**
   * 동기화 로그 획득
   */
  public getSyncLog(): HashChain {
    return this.syncLog;
  }

  /**
   * 현재 상태 획득
   */
  public getCurrentState(): ExecutionState | null {
    return this.currentState;
  }

  /**
   * 스냅샷 히스토리 획득
   */
  public getSnapshotHistory(): ReadonlyArray<StateSnapshot> {
    return Object.freeze([...this.stateSnapshots]);
  }
}

/**
 * SemanticSyncValidator: 외부 검증용 유틸리티
 */
export class SemanticSyncValidator {
  /**
   * 모든 노드의 의미 동등성 검증
   */
  public static async validateAllNodes(
    protocols: Map<bigint, SemanticSyncProtocol>
  ): Promise<boolean> {
    if (protocols.size < 2) {
      return true; // 노드가 1개 이하면 자명하게 동등
    }

    // 첫 번째 노드의 스냅샷을 기준으로 사용
    const firstEntry = protocols.entries().next();
    if (firstEntry.done || !firstEntry.value) {
      return false;
    }

    const [firstNodeId, firstProtocol] = firstEntry.value;

    for (const [nodeId, protocol] of protocols.entries()) {
      if (nodeId === firstNodeId) continue;

      // 실제 비교 로직
      const firstState = firstProtocol.getCurrentState();
      const otherState = protocol.getCurrentState();

      if (!firstState || !otherState) {
        return false;
      }

      if (firstState.calculateChecksum() !== otherState.calculateChecksum()) {
        return false;
      }
    }

    return true;
  }
}
