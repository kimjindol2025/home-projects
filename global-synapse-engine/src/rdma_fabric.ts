// 📡 Layer 1: Inter-Node Fabric - RDMA Zero-Copy Architecture
// Unforgiving Rule: Zero-Copy Guarantee = 100% (memory copies = 0)

import {
  RemoteMemoryRef,
  DirectMemoryRef,
  RDMAOperation,
  sha256,
  createSignature,
  verifySignature,
  CONSTANTS,
  measurePerformance,
} from './types';
import { HashChain } from './hash_chain';

/**
 * RemoteNodeHandle: 원격 노드와의 통신 핸들
 * Zero-Copy RDMA 작업을 위한 저수준 인터페이스
 */
class RemoteNodeHandle {
  nodeId: bigint;
  private isConnected: boolean = false;
  private signatureKey: Buffer;

  constructor(nodeId: bigint) {
    this.nodeId = nodeId;
    // 실제 구현에서는 TLS 인증서로부터 파생
    this.signatureKey = Buffer.from(`node-${nodeId.toString()}`.repeat(4).slice(0, 32));
  }

  /**
   * 원격 노드와 연결 (모의 구현)
   */
  async connect(): Promise<void> {
    // 실제로는 RDMA 디바이스와 통신
    await new Promise((resolve) => setTimeout(resolve, 1));
    this.isConnected = true;
  }

  /**
   * 연결 상태
   */
  isAlive(): boolean {
    return this.isConnected;
  }

  /**
   * 서명 키 획득 (내부용)
   */
  getSignatureKey(): Buffer {
    return this.signatureKey;
  }

  /**
   * 연결 종료
   */
  async disconnect(): Promise<void> {
    this.isConnected = false;
  }
}

/**
 * RDMAQueue: RDMA 작업 큐
 * Zero-Copy 작업을 배치로 관리하여 메모리 오버헤드 최소화
 */
class RDMAQueue {
  private queue: RDMAOperation[] = [];
  private inFlight: number = 0;
  private completedOperations: number = 0;

  /**
   * 작업을 큐에 추가 (배치)
   */
  enqueue(operation: RDMAOperation): void {
    this.queue.push(operation);
  }

  /**
   * 배치 처리 (실제 RDMA 하드웨어에 전송)
   * @param batchSize 배치 크기
   */
  async flushBatch(batchSize: number = CONSTANTS.RDMA_BATCH_SIZE): Promise<number> {
    const toProcess = Math.min(batchSize, this.queue.length);
    const batch = this.queue.splice(0, toProcess);

    if (batch.length === 0) return 0;

    this.inFlight += batch.length;

    // 실제 RDMA 송수신 (비동기)
    // 여기서는 시뮬레이션으로 1μs 레이턴시를 가정
    await new Promise((resolve) => setTimeout(resolve, 0.001));

    this.inFlight -= batch.length;
    this.completedOperations += batch.length;

    return batch.length;
  }

  /**
   * 모든 배치 처리 완료 대기
   */
  async waitAll(): Promise<void> {
    while (this.queue.length > 0 || this.inFlight > 0) {
      await this.flushBatch();
    }
  }

  /**
   * 큐 상태
   */
  getStats(): { queuedOps: number; inFlightOps: number; completedOps: number } {
    return {
      queuedOps: this.queue.length,
      inFlightOps: this.inFlight,
      completedOps: this.completedOperations,
    };
  }
}

/**
 * RDMAFabric: Layer 1 - Inter-Node Fabric
 *
 * Zero-Copy 보장:
 *   1. 메모리 복사 없음 (RDMA 하드웨어가 직접 접근)
 *   2. DirectMemoryRef를 통해 실제 메모리 포인터 사용
 *   3. 무결성은 HMAC 서명으로 검증 (데이터 복사 아님)
 *
 * Unforgiving Rules:
 *   - Zero-Copy Guarantee: 메모리 복사 = 0
 *   - Latency: <10μs (Inter-node)
 *   - Throughput: 1M ops/sec
 *   - Reliability: 0% data loss
 */
export class RDMAFabric {
  readonly nodeId: bigint;
  private localMemory: Buffer;
  private remoteNodes: Map<bigint, RemoteNodeHandle>;
  private rdmaQueue: RDMAQueue;
  private integrityLog: HashChain;
  private operationId: bigint = 1n;

  constructor(nodeId: bigint, memorySizeMB: number = 1024) {
    this.nodeId = nodeId;
    // 로컬 메모리 할당 (1GB)
    this.localMemory = Buffer.alloc(memorySizeMB * 1024 * 1024);
    this.remoteNodes = new Map();
    this.rdmaQueue = new RDMAQueue();
    this.integrityLog = new HashChain();

    console.log(`[LAYER1] RDMAFabric initialized: node ${nodeId}, memory ${memorySizeMB}MB`);
  }

  /**
   * 원격 노드 등록
   */
  async registerRemoteNode(remoteNodeId: bigint): Promise<void> {
    if (!this.remoteNodes.has(remoteNodeId)) {
      const handle = new RemoteNodeHandle(remoteNodeId);
      await handle.connect();
      this.remoteNodes.set(remoteNodeId, handle);
    }
  }

  /**
   * RDMA Read: 원격 메모리에서 로컬로 읽기 (Zero-Copy)
   *
   * Zero-Copy 메커니즘:
   *   1. 원격 주소를 직접 지정
   *   2. RDMA 하드웨어가 네트워크 전송 (CPU 개입 X)
   *   3. DirectMemoryRef로 로컬 메모리 위치 지정
   *   4. 메모리 복사 연산 없음 (RDMA가 처리)
   */
  async read(remoteRef: RemoteMemoryRef, localOffset: number): Promise<Buffer> {
    const { result: buffer, elapsedNs } = await measurePerformance(
      `RDMA.read from node ${remoteRef.nodeId}`,
      async () => {
        // RDMA 작업 생성
        const operation: RDMAOperation = {
          id: this.operationId++,
          type: 'read',
          sourceRef: {
            pointer: BigInt(remoteRef.address),
            offset: 0,
            length: remoteRef.size,
            accessMask: 1, // READ
          },
          destRef: {
            pointer: BigInt(this.localMemory.buffer.byteLength + localOffset),
            offset: localOffset,
            length: remoteRef.size,
            accessMask: 2, // WRITE
          },
          timestamp: BigInt(Date.now()) * 1_000_000n,
          signature: '', // 아래에서 계산
        };

        // 서명 계산 (HMAC)
        const remoteHandle = this.remoteNodes.get(remoteRef.nodeId);
        if (!remoteHandle) {
          throw new Error(`Remote node ${remoteRef.nodeId} not registered`);
        }

        const sigKey = remoteHandle.getSignatureKey();
        operation.signature = createSignature(
          Buffer.from(
            JSON.stringify(operation, (key, value) =>
              typeof value === 'bigint' ? value.toString() : value
            )
          ),
          sigKey
        );

        // 무결성 로그에 기록
        this.integrityLog.addLink(operation.id, 'read', operation.signature);

        // 큐에 추가
        this.rdmaQueue.enqueue(operation);
        await this.rdmaQueue.flushBatch(1);

        // 결과 데이터 반환 (실제로는 RDMA 완료 대기)
        return this.localMemory.slice(localOffset, localOffset + remoteRef.size);
      }
    );

    // Latency 검증
    const latencyUs = Number(elapsedNs) / 1000;
    if (latencyUs > CONSTANTS.MAX_INTER_NODE_LATENCY_US) {
      console.warn(`[WARN] RDMA.read latency exceeded: ${latencyUs}μs > ${CONSTANTS.MAX_INTER_NODE_LATENCY_US}μs`);
    }

    return buffer;
  }

  /**
   * RDMA Write: 로컬에서 원격 메모리로 쓰기 (Zero-Copy)
   */
  async write(localBuffer: Buffer, remoteRef: RemoteMemoryRef): Promise<void> {
    await measurePerformance(`RDMA.write to node ${remoteRef.nodeId}`, async () => {
      if (remoteRef.readOnly) {
        throw new Error(`Cannot write to read-only remote ref`);
      }

      const operation: RDMAOperation = {
        id: this.operationId++,
        type: 'write',
        sourceRef: {
          pointer: BigInt(localBuffer.buffer.byteLength),
          offset: 0,
          length: localBuffer.length,
          accessMask: 1, // READ
        },
        destRef: {
          pointer: BigInt(remoteRef.address),
          offset: 0,
          length: localBuffer.length,
          accessMask: 2, // WRITE
        },
        timestamp: BigInt(Date.now()) * 1_000_000n,
        signature: '',
      };

      const remoteHandle = this.remoteNodes.get(remoteRef.nodeId);
      if (!remoteHandle) {
        throw new Error(`Remote node ${remoteRef.nodeId} not registered`);
      }

      const sigKey = remoteHandle.getSignatureKey();
      operation.signature = createSignature(
        Buffer.from(
          JSON.stringify(operation, (key, value) =>
            typeof value === 'bigint' ? value.toString() : value
          )
        ),
        sigKey
      );

      this.integrityLog.addLink(operation.id, 'write', operation.signature);
      this.rdmaQueue.enqueue(operation);
      await this.rdmaQueue.flushBatch(1);
    });
  }

  /**
   * RDMA Compare-and-Swap: 원자적 수정
   * Zero-Copy 보장: 실제 메모리 위치에서 직접 연산
   */
  async compareAndSwap(
    remoteRef: RemoteMemoryRef,
    expectedValue: bigint,
    newValue: bigint
  ): Promise<boolean> {
    const { result: success, elapsedNs } = await measurePerformance(`RDMA.CAS on node ${remoteRef.nodeId}`, async () => {
      const operation: RDMAOperation = {
        id: this.operationId++,
        type: 'compare_swap',
        sourceRef: {
          pointer: BigInt(expectedValue),
          offset: 0,
          length: 8,
          accessMask: 4, // ATOMIC
        },
        destRef: {
          pointer: BigInt(newValue),
          offset: 0,
          length: 8,
          accessMask: 4, // ATOMIC
        },
        timestamp: BigInt(Date.now()) * 1_000_000n,
        signature: '',
      };

      const remoteHandle = this.remoteNodes.get(remoteRef.nodeId);
      if (!remoteHandle) {
        throw new Error(`Remote node ${remoteRef.nodeId} not registered`);
      }

      const sigKey = remoteHandle.getSignatureKey();
      operation.signature = createSignature(
        Buffer.from(
          JSON.stringify(operation, (key, value) =>
            typeof value === 'bigint' ? value.toString() : value
          )
        ),
        sigKey
      );

      this.integrityLog.addLink(operation.id, 'compare_swap', operation.signature);
      this.rdmaQueue.enqueue(operation);
      await this.rdmaQueue.flushBatch(1);

      // 모의 결과: 50% 성공률
      return Math.random() > 0.5;
    });

    return success;
  }

  /**
   * 모든 RDMA 작업 완료 대기
   */
  async waitAllOperations(): Promise<void> {
    await this.rdmaQueue.waitAll();
  }

  /**
   * 무결성 로그 검증
   */
  public verifyIntegrity(): boolean {
    const verification = this.integrityLog.verify();
    return verification.isValid;
  }

  /**
   * Layer 1 통계
   */
  public getStats(): {
    nodeId: bigint;
    operationsCompleted: number;
    integrityValid: boolean;
    remoteNodesCount: number;
    memoryUsedBytes: number;
  } {
    const queueStats = this.rdmaQueue.getStats();
    return {
      nodeId: this.nodeId,
      operationsCompleted: queueStats.completedOps,
      integrityValid: this.verifyIntegrity(),
      remoteNodesCount: this.remoteNodes.size,
      memoryUsedBytes: this.localMemory.length,
    };
  }

  /**
   * Layer 1 무결성 로그 획득
   */
  public getIntegrityLog(): HashChain {
    return this.integrityLog;
  }
}
