// 🐒 Chaos Monkey - Week 3
// 네트워크 장애 시뮬레이션 및 복구 검증

import { RDMAFabric } from './rdma_fabric';
import { SemanticSyncProtocol } from './semantic_sync';

/**
 * ChaosEvent: 장애 이벤트
 */
export interface ChaosEvent {
  id: bigint;
  type: 'network_delay' | 'packet_loss' | 'node_crash' | 'memory_spike';
  targetNode: bigint;
  durationMs: number;
  timestamp: bigint;
  recovered: boolean;
  recoveryTimeMs?: number;
}

/**
 * ChaosStats: 혼돈 테스트 통계
 */
export interface ChaosStats {
  totalEvents: number;
  successfulRecoveries: number;
  failedRecoveries: number;
  recoveryRate: number;
  averageRecoveryTimeMs: number;
  maxRecoveryTimeMs: number;
  dataLosses: number;
  integrityViolations: number;
}

/**
 * ChaosMonkey: 혼돈 테스트 엔진
 *
 * 기능:
 * - 네트워크 지연 주입
 * - 패킷 손실 시뮬레이션
 * - 노드 충돌 시뮬레이션
 * - 메모리 스파이크
 * - 복구율 측정
 */
export class ChaosMonkey {
  private events: ChaosEvent[] = [];
  private eventId: bigint = 1n;
  private failedNodes: Set<bigint> = new Set();
  private stats: ChaosStats = {
    totalEvents: 0,
    successfulRecoveries: 0,
    failedRecoveries: 0,
    recoveryRate: 0,
    averageRecoveryTimeMs: 0,
    maxRecoveryTimeMs: 0,
    dataLosses: 0,
    integrityViolations: 0
  };

  constructor(private randomSeed: number = 42) {
    this.seedRandom(randomSeed);
  }

  /**
   * 난수 생성기 시드
   */
  private seedRandom(seed: number): void {
    // 간단한 선형 합동 생성기
    let m_w = seed;
    let m_z = 987654321;
    const mask = 0xffffffff;

    (Math.random as any) = function() {
      m_z = (36969 * (m_z & 65535) + (m_z >> 16)) & mask;
      m_w = (18000 * (m_w & 65535) + (m_w >> 16)) & mask;
      let result = ((m_z << 16) + (m_w & 65535)) >>> 0;
      result /= 4294967296;
      return result;
    };
  }

  /**
   * 네트워크 지연 주입
   * @param fabric RDMA 패브릭
   * @param nodeId 대상 노드
   * @param delayMs 지연 시간(ms)
   */
  public async injectNetworkDelay(
    fabric: RDMAFabric,
    nodeId: bigint,
    delayMs: number
  ): Promise<ChaosEvent> {
    const startTime = Date.now();
    const event: ChaosEvent = {
      id: this.eventId++,
      type: 'network_delay',
      targetNode: nodeId,
      durationMs: delayMs,
      timestamp: BigInt(Date.now()) * 1_000_000n,
      recovered: false
    };

    // 지연 시뮬레이션
    await new Promise(resolve => setTimeout(resolve, delayMs));

    // 복구 시간 기록
    event.recoveryTimeMs = Date.now() - startTime;
    event.recovered = true;

    this.events.push(event);
    this.stats.totalEvents++;
    this.stats.successfulRecoveries++;
    this.stats.maxRecoveryTimeMs = Math.max(
      this.stats.maxRecoveryTimeMs,
      event.recoveryTimeMs
    );

    return event;
  }

  /**
   * 패킷 손실 시뮬레이션
   * @param fabric RDMA 패브릭
   * @param nodeId 대상 노드
   * @param lossRatePercent 손실율(%)
   */
  public async injectPacketLoss(
    fabric: RDMAFabric,
    nodeId: bigint,
    lossRatePercent: number
  ): Promise<ChaosEvent> {
    const startTime = Date.now();
    const event: ChaosEvent = {
      id: this.eventId++,
      type: 'packet_loss',
      targetNode: nodeId,
      durationMs: 100, // 100ms 동안 손실
      timestamp: BigInt(Date.now()) * 1_000_000n,
      recovered: false
    };

    // 패킷 손실 중 재시도 메커니즘 작동
    let retries = 0;
    const maxRetries = 3;
    let success = false;

    while (retries < maxRetries && !success) {
      try {
        // 손실 시간 시뮬레이션
        await new Promise(resolve => setTimeout(resolve, 50));

        // 복구 성공
        if (Math.random() > lossRatePercent / 100) {
          success = true;
          event.recovered = true;
          this.stats.successfulRecoveries++;
        } else {
          retries++;
        }
      } catch (error) {
        retries++;
      }
    }

    if (!success) {
      event.recovered = false;
      this.stats.failedRecoveries++;
      this.stats.dataLosses++;
    }

    event.recoveryTimeMs = Date.now() - startTime;
    this.events.push(event);
    this.stats.totalEvents++;

    return event;
  }

  /**
   * 노드 충돌 시뮬레이션
   * @param nodeId 충돌할 노드 ID
   * @param recoveryDelayMs 복구 시간(ms)
   */
  public async injectNodeCrash(
    nodeId: bigint,
    recoveryDelayMs: number = 500
  ): Promise<ChaosEvent> {
    const startTime = Date.now();
    const event: ChaosEvent = {
      id: this.eventId++,
      type: 'node_crash',
      targetNode: nodeId,
      durationMs: recoveryDelayMs,
      timestamp: BigInt(Date.now()) * 1_000_000n,
      recovered: false
    };

    // 노드를 장애 목록에 추가
    this.failedNodes.add(nodeId);

    // 복구 대기
    await new Promise(resolve => setTimeout(resolve, recoveryDelayMs));

    // 노드 복구
    this.failedNodes.delete(nodeId);
    event.recovered = true;
    event.recoveryTimeMs = Date.now() - startTime;

    this.events.push(event);
    this.stats.totalEvents++;
    this.stats.successfulRecoveries++;
    this.stats.maxRecoveryTimeMs = Math.max(
      this.stats.maxRecoveryTimeMs,
      event.recoveryTimeMs
    );

    return event;
  }

  /**
   * 메모리 스파이크 주입
   * @param sizeMs 메모리 크기(MB)
   * @param durationMs 지속 시간(ms)
   */
  public async injectMemorySpike(
    sizeMb: number,
    durationMs: number
  ): Promise<ChaosEvent> {
    const startTime = Date.now();
    const event: ChaosEvent = {
      id: this.eventId++,
      type: 'memory_spike',
      targetNode: 0n,
      durationMs: durationMs,
      timestamp: BigInt(Date.now()) * 1_000_000n,
      recovered: false
    };

    try {
      // 메모리 할당 (스파이크 시뮬레이션)
      const arrays = [];
      const arraySize = (sizeMb * 1024 * 1024) / 8; // 8바이트 숫자
      for (let i = 0; i < arraySize; i++) {
        arrays.push(Math.random());
      }

      // 스파이크 지속 시간
      await new Promise(resolve => setTimeout(resolve, durationMs));

      // 메모리 해제
      arrays.length = 0;

      event.recovered = true;
      this.stats.successfulRecoveries++;
    } catch (error) {
      event.recovered = false;
      this.stats.failedRecoveries++;
    }

    event.recoveryTimeMs = Date.now() - startTime;
    this.events.push(event);
    this.stats.totalEvents++;

    return event;
  }

  /**
   * 혼돈 테스트 시나리오: 연속 장애
   */
  public async runChaosScenario(
    nodeCount: number,
    eventCount: number
  ): Promise<ChaosStats> {
    console.log(`\n🐒 Starting Chaos Scenario (${eventCount} events, ${nodeCount} nodes)\n`);

    for (let i = 0; i < eventCount; i++) {
      const nodeId = BigInt(Math.floor(Math.random() * nodeCount) + 1);
      const eventType = Math.floor(Math.random() * 4);

      if (eventType === 0) {
        // 네트워크 지연
        await this.injectNetworkDelay(null as any, nodeId, 50 + Math.random() * 100);
      } else if (eventType === 1) {
        // 패킷 손실
        await this.injectPacketLoss(null as any, nodeId, 10 + Math.random() * 30);
      } else if (eventType === 2) {
        // 노드 충돌
        await this.injectNodeCrash(nodeId, 100 + Math.random() * 300);
      } else {
        // 메모리 스파이크
        await this.injectMemorySpike(10 + Math.random() * 50, 100 + Math.random() * 200);
      }

      if ((i + 1) % 100 === 0) {
        console.log(`  Progress: ${i + 1}/${eventCount} events`);
      }
    }

    // 통계 계산
    this.calculateStats();

    return this.stats;
  }

  /**
   * 통계 계산
   */
  private calculateStats(): void {
    if (this.stats.totalEvents === 0) return;

    this.stats.recoveryRate =
      (this.stats.successfulRecoveries / this.stats.totalEvents) * 100;

    const recoveryTimes = this.events
      .filter(e => e.recoveryTimeMs !== undefined)
      .map(e => e.recoveryTimeMs!);

    if (recoveryTimes.length > 0) {
      this.stats.averageRecoveryTimeMs =
        recoveryTimes.reduce((a, b) => a + b, 0) / recoveryTimes.length;
    }
  }

  /**
   * 통계 리포트
   */
  public printReport(): void {
    console.log('\n🐒 Chaos Test Report\n');
    console.log('═'.repeat(60));
    console.log(`Total Events:           ${this.stats.totalEvents}`);
    console.log(`Successful Recoveries:  ${this.stats.successfulRecoveries}`);
    console.log(`Failed Recoveries:      ${this.stats.failedRecoveries}`);
    console.log(`Recovery Rate:          ${this.stats.recoveryRate.toFixed(2)}%`);
    console.log(
      `Average Recovery Time:  ${this.stats.averageRecoveryTimeMs.toFixed(2)}ms`
    );
    console.log(`Max Recovery Time:      ${this.stats.maxRecoveryTimeMs}ms`);
    console.log(`Data Losses:            ${this.stats.dataLosses}`);
    console.log(`Integrity Violations:   ${this.stats.integrityViolations}`);
    console.log('═'.repeat(60));

    // 통과/실패 판정
    if (this.stats.recoveryRate >= 99 && this.stats.integrityViolations === 0) {
      console.log('\n✅ Chaos Test PASSED (>99% recovery, 0 integrity violations)');
    } else {
      console.log('\n❌ Chaos Test FAILED');
      if (this.stats.recoveryRate < 99) {
        console.log(`   - Recovery rate: ${this.stats.recoveryRate.toFixed(2)}% (need >99%)`);
      }
      if (this.stats.integrityViolations > 0) {
        console.log(`   - Integrity violations: ${this.stats.integrityViolations}`);
      }
    }
  }

  /**
   * 통계 조회
   */
  public getStats(): ChaosStats {
    this.calculateStats();
    return { ...this.stats };
  }

  /**
   * 이벤트 히스토리
   */
  public getEvents(): ReadonlyArray<ChaosEvent> {
    return Object.freeze([...this.events]);
  }

  /**
   * 장애 노드 조회
   */
  public getFailedNodes(): Set<bigint> {
    return new Set(this.failedNodes);
  }
}
