// 🌐 Global Synapse Engine - Core Types
// Layer: Shared across all 4 layers

import * as crypto from 'crypto';

/**
 * RemoteMemoryRef: 원격 노드의 메모리 위치를 안전하게 참조
 * Zero-Copy 보장을 위해 RDMA 직접 접근 가능
 */
export interface RemoteMemoryRef {
  nodeId: bigint;         // 원격 노드 ID
  address: bigint;        // 메모리 주소 (0부터 1GB 범위)
  size: number;           // 데이터 크기 (바이트)
  readOnly: boolean;      // 읽기 전용 플래그
}

/**
 * DirectMemoryRef: Zero-Copy를 위한 메모리 직접 참조
 * RDMA 하드웨어가 직접 접근 가능하도록 설계
 */
export interface DirectMemoryRef {
  pointer: bigint;        // 실제 메모리 포인터 (CPU 주소)
  offset: number;         // 오프셋
  length: number;         // 데이터 길이
  accessMask: number;     // 접근 권한 (READ=1, WRITE=2, ATOMIC=4)
}

/**
 * RDMAOperation: RDMA 하드웨어 작업 단위
 * Zero-Copy를 보장하며 메모리 카피 없음
 */
export interface RDMAOperation {
  id: bigint;             // 작업 ID
  type: 'read' | 'write' | 'atomic' | 'compare_swap';
  sourceRef: DirectMemoryRef;
  destRef: DirectMemoryRef;
  timestamp: bigint;      // 작업 시각 (나노초)
  signature: string;      // HMAC 서명 (무결성 보증)
}

/**
 * StateSnapshot: Layer 2 Semantic Sync에서 사용
 * 노드 상태를 시점별로 저장하여 동등성 검증
 */
export interface StateSnapshot {
  nodeId: bigint;
  executionId: bigint;    // 실행 ID (Increment)
  timestamp: bigint;      // 스냅샷 시각
  memoryChecksum: string; // SHA256 체크섬
  instructionCounter: number;  // 실행된 명령어 수
  variableState: Map<string, string>;  // 변수=값 맵
  hashChainLink: string;  // 이전 스냅샷 해시
}

/**
 * SemanticEquivalenceProof: 두 노드 간의 의미 동등성 증명
 * 의미 동등성 = 동일 입력 → 동일 출력 (절대 성능 무관)
 */
export interface SemanticEquivalenceProof {
  node1Id: bigint;
  node2Id: bigint;
  inputHash: string;      // 입력 데이터 SHA256
  outputHash: string;     // 출력 데이터 SHA256
  snapshot1: StateSnapshot;
  snapshot2: StateSnapshot;
  isEquivalent: boolean;
  proofTime: bigint;      // 증명 생성 시간 (나노초)
}

/**
 * FaultEvent: 장애 감지 이벤트
 * Layer 4 Self-Healing에서 사용
 */
export interface FaultEvent {
  detectTime: bigint;
  nodeId: bigint;
  faultType:
    | 'RDMATimeout'
    | 'SemanticDesync'
    | 'MemoryChecksum'
    | 'HeartbeatMissed'
    | 'DataLoss'
    | 'NetworkPartition';
  severity: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW';
  details: string;
}

/**
 * HashChainLink: Test Mouse 검증용 해시 체인
 * 모든 중요 작업을 기록하여 개조 불가능하도록 보증
 */
export interface HashChainLink {
  index: number;
  timestamp: bigint;
  previousHash: string;   // SHA256
  content: {
    type: string;
    operationId: bigint;
    signature: string;
  };
  hash: string;           // SHA256(previousHash + content)
}

/**
 * PerformanceMetrics: 성능 추적용 메트릭
 * Unforgiving Rule 검증에 사용
 */
export interface PerformanceMetrics {
  operationName: string;
  minLatency: number;     // 나노초
  maxLatency: number;
  avgLatency: number;
  p99Latency: number;
  throughput: number;     // ops/sec
  memoryUsed: number;     // 바이트
  errorRate: number;      // %
}

/**
 * 상수 정의
 */
export const CONSTANTS = {
  MAX_MEMORY_SIZE: 1_073_741_824,  // 1 GB
  MAX_NODES: 10_000,
  MAX_OPERATIONS_PER_SECOND: 1_000_000,

  // Unforgiving Rules 임계값
  MAX_INTER_NODE_LATENCY_US: 10,        // 10 마이크로초
  MAX_RECOVERY_TIME_MS: 1,              // 1 밀리초
  MAX_FAULT_DETECTION_MS: 100,          // 100 밀리초
  MIN_SEMANTIC_EQUIVALENCE: 1.0,        // 100% 정확성
  MAX_DATA_LOSS_RATE: 0.0,              // 0% 손실
  MIN_RECOVERY_SUCCESS_RATE: 1.0,       // 100% 성공

  // 메모리 관련
  MEMORY_PAGE_SIZE: 4096,               // 4KB
  RDMA_BATCH_SIZE: 256,                 // 배치 크기

  // 시간 관련
  HEARTBEAT_INTERVAL_MS: 100,
  HEARTBEAT_TIMEOUT_MS: 500,
  NODE_RECOVERY_TIMEOUT_MS: 5000,
};

/**
 * 성능 측정 함수
 * @param name 작업 이름
 * @param fn 실행할 함수
 * @returns 작업 결과와 나노초 단위 경과 시간
 */
export async function measurePerformance<T>(
  name: string,
  fn: () => Promise<T>
): Promise<{ result: T; elapsedNs: bigint }> {
  const startTime = process.hrtime.bigint();
  const result = await fn();
  const endTime = process.hrtime.bigint();
  const elapsedNs = endTime - startTime;

  console.log(`[PERF] ${name}: ${Number(elapsedNs) / 1000}μs`);

  return { result, elapsedNs };
}

/**
 * SHA256 해시 계산
 */
export function sha256(data: Buffer | string): string {
  const buffer = typeof data === 'string' ? Buffer.from(data) : data;
  return crypto.createHash('sha256').update(buffer).digest('hex');
}

/**
 * HMAC 서명 생성 (RDMA 무결성)
 */
export function createSignature(data: Buffer, key: Buffer): string {
  return crypto
    .createHmac('sha256', key)
    .update(data)
    .digest('hex');
}

/**
 * 서명 검증
 */
export function verifySignature(data: Buffer, signature: string, key: Buffer): boolean {
  const expected = createSignature(data, key);
  return expected === signature;
}

export class NodeState {
  nodeId: bigint;
  epoch: number = 0;
  isAlive: boolean = true;
  lastHeartbeat: bigint = 0n;

  constructor(nodeId: bigint) {
    this.nodeId = nodeId;
    this.lastHeartbeat = BigInt(Date.now()) * 1_000_000n;
  }

  updateHeartbeat(): void {
    this.lastHeartbeat = BigInt(Date.now()) * 1_000_000n;
  }

  isStale(currentTime: bigint, timeoutMs: number): boolean {
    const timeoutNs = BigInt(timeoutMs) * 1_000_000n;
    return currentTime - this.lastHeartbeat > timeoutNs;
  }
}
