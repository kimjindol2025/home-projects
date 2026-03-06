// 🤖 Auto Recovery - 통합 자동 복구 메커니즘 (Week 4)
// Circuit Breaker + Timeout Manager + Retry Strategy

export { CircuitBreaker } from './circuit-breaker';
export { CircuitBreakerState } from './types';
export type { CircuitBreakerStatus } from './types';

export { TimeoutManager, getGlobalTimeoutManager } from './timeout-manager';
export type { TimeoutConfig } from './types';

export { RetryStrategy, getGlobalRetryStrategy } from './retry-strategy';
export type { RetryConfig } from './types';
export interface RetryResult<T> {
  success: boolean;
  result?: T;
  error?: Error;
  attempts: number;
  totalDelayMs: number;
}

export type { AutoRecoveryStats } from './types';

/**
 * AutoRecoveryOrchestrator: 3가지 메커니즘을 통합 관리
 * 작동 방식:
 * 1. 요청 전: Circuit Breaker가 노드 상태 확인 (차단 또는 통과)
 * 2. 실행 중: Timeout Manager가 타임아웃 감시
 * 3. 실패 시: Retry Strategy가 지수 백오프로 재시도
 * 4. 성공/실패: Circuit Breaker가 상태 업데이트
 */

import { CircuitBreaker } from './circuit-breaker';
import { CircuitBreakerState } from './types';
import { TimeoutManager, getGlobalTimeoutManager } from './timeout-manager';
import { RetryStrategy, getGlobalRetryStrategy } from './retry-strategy';
import { AutoRecoveryStats } from './types';

export class AutoRecoveryOrchestrator {
  private circuitBreaker: CircuitBreaker;
  private timeoutManager: TimeoutManager;
  private retryStrategy: RetryStrategy;
  private stats: AutoRecoveryStats = {
    totalAttempts: 0,
    successfulRecoveries: 0,
    failedRecoveries: 0,
    recoveryRate: 0,
    averageRecoveryTimeMs: 0,
    maxRecoveryTimeMs: 0,
    totalRetries: 0,
    circuitBreakerTrips: 0
  };
  private recoveryTimes: number[] = [];

  constructor(
    circuitBreaker?: CircuitBreaker,
    timeoutManager?: TimeoutManager,
    retryStrategy?: RetryStrategy
  ) {
    this.circuitBreaker = circuitBreaker ?? new CircuitBreaker();
    this.timeoutManager = timeoutManager ?? getGlobalTimeoutManager();
    this.retryStrategy = retryStrategy ?? getGlobalRetryStrategy();
  }

  /**
   * 자동 복구와 함께 작업 실행
   * @param operation 비동기 작업
   * @param nodeId 작업을 실행할 노드 ID
   * @param operationType 작업 타입 (타임아웃 설정용)
   */
  public async execute<T>(
    operation: () => Promise<T>,
    nodeId: bigint,
    operationType: string = 'global'
  ): Promise<{ result?: T; success: boolean; error?: Error }> {
    const startTime = Date.now();
    this.stats.totalAttempts++;

    // Step 1: Circuit Breaker 확인
    if (!this.circuitBreaker.canExecute(nodeId)) {
      this.stats.failedRecoveries++;
      this.stats.circuitBreakerTrips++;
      return {
        success: false,
        error: new Error(
          `Circuit Breaker OPEN for node ${nodeId} - request rejected`
        )
      };
    }

    // Step 2: Retry + Timeout 통합 실행
    try {
      const timeout = this.timeoutManager.getTimeout(operationType);
      const result = await this.retryStrategy.executeWithTimeout(
        operation,
        timeout,
        `Node${nodeId}:${operationType}`
      );

      if (result.success) {
        // Step 3: 성공 - Circuit Breaker 상태 업데이트
        this.circuitBreaker.recordSuccess(nodeId);
        this.timeoutManager.recordLatency(operationType, result.totalDelayMs);

        const recoveryTime = Date.now() - startTime;
        this.recoveryTimes.push(recoveryTime);
        this.stats.successfulRecoveries++;
        this.stats.totalRetries += result.attempts - 1;
        this.stats.maxRecoveryTimeMs = Math.max(
          this.stats.maxRecoveryTimeMs,
          recoveryTime
        );

        return { result: result.result, success: true };
      } else {
        // Step 3: 실패 - Circuit Breaker 상태 업데이트
        this.circuitBreaker.recordFailure(nodeId);
        this.stats.failedRecoveries++;
        this.stats.totalRetries += result.attempts;

        return { success: false, error: result.error };
      }
    } catch (error) {
      // 예상치 못한 에러
      this.circuitBreaker.recordFailure(nodeId);
      this.stats.failedRecoveries++;

      return {
        success: false,
        error: error instanceof Error ? error : new Error(String(error))
      };
    } finally {
      // 복구율 계산
      if (this.stats.totalAttempts > 0) {
        this.stats.recoveryRate =
          (this.stats.successfulRecoveries / this.stats.totalAttempts) * 100;
      }

      // 평균 복구 시간
      if (this.recoveryTimes.length > 0) {
        const sum = this.recoveryTimes.reduce((a, b) => a + b, 0);
        this.stats.averageRecoveryTimeMs = sum / this.recoveryTimes.length;
      }
    }
  }

  /**
   * 통계 조회
   */
  public getStats(): AutoRecoveryStats {
    return { ...this.stats };
  }

  /**
   * Circuit Breaker 조회
   */
  public getCircuitBreaker(): CircuitBreaker {
    return this.circuitBreaker;
  }

  /**
   * Timeout Manager 조회
   */
  public getTimeoutManager(): TimeoutManager {
    return this.timeoutManager;
  }

  /**
   * Retry Strategy 조회
   */
  public getRetryStrategy(): RetryStrategy {
    return this.retryStrategy;
  }

  /**
   * 통합 리포트 출력
   */
  public printReport(): void {
    console.log('\n🤖 Auto Recovery Orchestrator Report\n');
    console.log('═'.repeat(80));
    console.log(`Total Attempts:            ${this.stats.totalAttempts}`);
    console.log(`Successful Recoveries:     ${this.stats.successfulRecoveries}`);
    console.log(`Failed Recoveries:         ${this.stats.failedRecoveries}`);
    console.log(
      `Recovery Rate:             ${this.stats.recoveryRate.toFixed(2)}%`
    );
    console.log(
      `Average Recovery Time:     ${this.stats.averageRecoveryTimeMs.toFixed(2)}ms`
    );
    console.log(`Max Recovery Time:         ${this.stats.maxRecoveryTimeMs}ms`);
    console.log(`Total Retries:             ${this.stats.totalRetries}`);
    console.log(
      `Circuit Breaker Trips:     ${this.stats.circuitBreakerTrips}`
    );
    console.log('═'.repeat(80));

    console.log('\n📊 Sub-Components:\n');
    this.circuitBreaker.printReport();
    this.timeoutManager.printReport();
    this.retryStrategy.printReport();
  }

  /**
   * 통계 초기화
   */
  public resetStats(): void {
    this.stats = {
      totalAttempts: 0,
      successfulRecoveries: 0,
      failedRecoveries: 0,
      recoveryRate: 0,
      averageRecoveryTimeMs: 0,
      maxRecoveryTimeMs: 0,
      totalRetries: 0,
      circuitBreakerTrips: 0
    };
    this.recoveryTimes = [];
    this.circuitBreaker.reset(0n);
    this.timeoutManager.resetHistory();
    this.retryStrategy.resetStats();
  }
}

// 전역 Orchestrator 인스턴스
let globalOrchestrator: AutoRecoveryOrchestrator | null = null;

export function getGlobalAutoRecoveryOrchestrator(): AutoRecoveryOrchestrator {
  if (!globalOrchestrator) {
    globalOrchestrator = new AutoRecoveryOrchestrator();
  }
  return globalOrchestrator;
}
