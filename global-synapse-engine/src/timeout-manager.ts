// ⏱️ Timeout Manager - 타임아웃 표준화 및 동적 조정 (Week 4)

import { TimeoutConfig } from './types';

/**
 * TimeoutManager: 작업 타입별 표준 타임아웃 관리
 *
 * 기능:
 * - 작업 타입별 기본 타임아웃 정의
 * - 평균 응답 시간 기반 동적 조정
 * - 타임아웃 초과 감지
 * - 적응형 타임아웃 (moving average 기반)
 */
export class TimeoutManager {
  private baseConfig: TimeoutConfig;
  private adaptiveTimeouts: Map<string, number> = new Map();
  private operationLatencies: Map<string, number[]> = new Map();
  private readonly maxHistorySize = 100;

  constructor(baseConfig: Partial<TimeoutConfig> = {}) {
    this.baseConfig = {
      rdmaReadMs: baseConfig.rdmaReadMs ?? 100,
      rdmaWriteMs: baseConfig.rdmaWriteMs ?? 150,
      semanticSyncMs: baseConfig.semanticSyncMs ?? 500,
      heartbeatMs: baseConfig.heartbeatMs ?? 500,
      globalMs: baseConfig.globalMs ?? 1000
    };

    // 적응형 타임아웃 초기화
    this.adaptiveTimeouts.set('rdmaRead', this.baseConfig.rdmaReadMs);
    this.adaptiveTimeouts.set('rdmaWrite', this.baseConfig.rdmaWriteMs);
    this.adaptiveTimeouts.set('semanticSync', this.baseConfig.semanticSyncMs);
    this.adaptiveTimeouts.set('heartbeat', this.baseConfig.heartbeatMs);
    this.adaptiveTimeouts.set('global', this.baseConfig.globalMs);
  }

  /**
   * 작업 시작 (타이머 생성)
   * @returns 타임아웃 ID
   */
  public startTimeout(
    operationType: string,
    onTimeout: () => void
  ): NodeJS.Timeout {
    const timeoutMs = this.getTimeout(operationType);
    return setTimeout(onTimeout, timeoutMs);
  }

  /**
   * 타임아웃 취소
   */
  public clearTimeout(timeoutHandle: NodeJS.Timeout): void {
    clearTimeout(timeoutHandle);
  }

  /**
   * 작업 완료 시 레이턴시 기록 (동적 조정용)
   */
  public recordLatency(operationType: string, latencyMs: number): void {
    if (!this.operationLatencies.has(operationType)) {
      this.operationLatencies.set(operationType, []);
    }

    const latencies = this.operationLatencies.get(operationType)!;
    latencies.push(latencyMs);

    // 히스토리 크기 제한
    if (latencies.length > this.maxHistorySize) {
      latencies.shift();
    }

    // 적응형 타임아웃 계산 (p99 레이턴시 + 마진)
    this.updateAdaptiveTimeout(operationType);
  }

  /**
   * 현재 타임아웃 값 조회 (적응형)
   */
  public getTimeout(operationType: string): number {
    return this.adaptiveTimeouts.get(operationType) ?? this.baseConfig.globalMs;
  }

  /**
   * 기본 타임아웃 값 조회
   */
  public getBaseTimeout(operationType: string): number {
    switch (operationType) {
      case 'rdmaRead':
        return this.baseConfig.rdmaReadMs;
      case 'rdmaWrite':
        return this.baseConfig.rdmaWriteMs;
      case 'semanticSync':
        return this.baseConfig.semanticSyncMs;
      case 'heartbeat':
        return this.baseConfig.heartbeatMs;
      default:
        return this.baseConfig.globalMs;
    }
  }

  /**
   * 평균 레이턴시 조회
   */
  public getAverageLatency(operationType: string): number {
    const latencies = this.operationLatencies.get(operationType);
    if (!latencies || latencies.length === 0) {
      return -1;
    }

    const sum = latencies.reduce((a, b) => a + b, 0);
    return sum / latencies.length;
  }

  /**
   * p99 레이턴시 조회 (99 백분위수)
   */
  public getP99Latency(operationType: string): number {
    const latencies = this.operationLatencies.get(operationType);
    if (!latencies || latencies.length === 0) {
      return -1;
    }

    const sorted = [...latencies].sort((a, b) => a - b);
    const p99Index = Math.ceil(sorted.length * 0.99) - 1;
    return sorted[Math.max(0, p99Index)];
  }

  /**
   * 적응형 타임아웃 업데이트
   * 적응 공식: max(baseTimeout, p99 * 1.5) (50% 마진)
   */
  private updateAdaptiveTimeout(operationType: string): void {
    const p99 = this.getP99Latency(operationType);
    if (p99 < 0) return;

    const baseTimeout = this.getBaseTimeout(operationType);
    const adaptiveTimeout = Math.max(baseTimeout, Math.ceil(p99 * 1.5));

    this.adaptiveTimeouts.set(operationType, adaptiveTimeout);
  }

  /**
   * 타임아웃 초과 여부 확인
   * @param startTimeMs 작업 시작 시각
   * @param operationType 작업 타입
   * @returns true면 타임아웃 초과
   */
  public isTimedOut(startTimeMs: number, operationType: string): boolean {
    const elapsed = Date.now() - startTimeMs;
    const timeout = this.getTimeout(operationType);
    return elapsed > timeout;
  }

  /**
   * 통계 리포트 출력
   */
  public printReport(): void {
    console.log('\n⏱️  Timeout Manager Report\n');
    console.log('═'.repeat(80));

    const operationTypes = Array.from(this.operationLatencies.keys());
    operationTypes.forEach(opType => {
      const avgLatency = this.getAverageLatency(opType);
      const p99Latency = this.getP99Latency(opType);
      const baseTimeout = this.getBaseTimeout(opType);
      const adaptiveTimeout = this.getTimeout(opType);
      const historySize = this.operationLatencies.get(opType)?.length ?? 0;

      console.log(`${opType}:`);
      console.log(`  Base Timeout:      ${baseTimeout}ms`);
      console.log(`  Adaptive Timeout:  ${adaptiveTimeout}ms`);
      console.log(`  Avg Latency:       ${avgLatency.toFixed(2)}ms`);
      console.log(`  P99 Latency:       ${p99Latency.toFixed(2)}ms`);
      console.log(`  History Size:      ${historySize}/${this.maxHistorySize}`);
      console.log();
    });

    console.log('═'.repeat(80));
    console.log('\n📊 적응형 타임아웃 공식: max(baseTimeout, p99 * 1.5)');
  }

  /**
   * 히스토리 초기화 (테스트용)
   */
  public resetHistory(): void {
    this.operationLatencies.clear();
    this.adaptiveTimeouts.clear();

    this.adaptiveTimeouts.set('rdmaRead', this.baseConfig.rdmaReadMs);
    this.adaptiveTimeouts.set('rdmaWrite', this.baseConfig.rdmaWriteMs);
    this.adaptiveTimeouts.set('semanticSync', this.baseConfig.semanticSyncMs);
    this.adaptiveTimeouts.set('heartbeat', this.baseConfig.heartbeatMs);
    this.adaptiveTimeouts.set('global', this.baseConfig.globalMs);
  }
}

// 전역 TimeoutManager 인스턴스
let globalTimeoutManager: TimeoutManager | null = null;

export function getGlobalTimeoutManager(): TimeoutManager {
  if (!globalTimeoutManager) {
    globalTimeoutManager = new TimeoutManager();
  }
  return globalTimeoutManager;
}
