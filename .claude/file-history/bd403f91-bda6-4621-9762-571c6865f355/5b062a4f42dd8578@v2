// ⚡ Circuit Breaker - 장애 격리 및 자동 복구 (Week 4)
// 상태: CLOSED → OPEN → HALF_OPEN → CLOSED

export { CircuitBreakerState } from './types';
import { CircuitBreakerState, CircuitBreakerStatus } from './types';

/**
 * Circuit Breaker: 장애 노드 격리 및 자동 복구 메커니즘
 *
 * 상태 전환 다이어그램:
 * CLOSED (정상)
 *   ↓ (failureThreshold 초과)
 * OPEN (격리)
 *   ↓ (timeoutMs 경과)
 * HALF_OPEN (복구 시도)
 *   ↓ (성공) → CLOSED
 *   ↓ (실패) → OPEN
 */
export class CircuitBreaker {
  private statuses: Map<bigint, CircuitBreakerStatus> = new Map();
  private readonly failureThreshold: number;      // CLOSED → OPEN 실패 횟수
  private readonly successThreshold: number;      // HALF_OPEN → CLOSED 성공 횟수
  private readonly timeout: number;               // OPEN → HALF_OPEN 대기 시간 (ms)

  constructor(
    failureThreshold: number = 5,
    successThreshold: number = 2,
    timeoutMs: number = 5000
  ) {
    this.failureThreshold = failureThreshold;
    this.successThreshold = successThreshold;
    this.timeout = timeoutMs;
  }

  /**
   * 작업 성공 처리
   */
  public recordSuccess(nodeId: bigint): void {
    let status = this.statuses.get(nodeId);
    if (!status) {
      status = this.createStatus(nodeId);
      this.statuses.set(nodeId, status);
    }

    status.totalSuccesses++;

    if (status.state === CircuitBreakerState.CLOSED) {
      // 정상 상태에서 성공: 실패 카운트 감소
      status.failureCount = Math.max(0, status.failureCount - 1);
    } else if (status.state === CircuitBreakerState.HALF_OPEN) {
      // 복구 시도 중 성공: 성공 카운트 증가
      status.successCount++;
      if (status.successCount >= this.successThreshold) {
        // HALF_OPEN → CLOSED
        this.transitionState(nodeId, CircuitBreakerState.CLOSED);
      }
    }
  }

  /**
   * 작업 실패 처리
   */
  public recordFailure(nodeId: bigint): void {
    let status = this.statuses.get(nodeId);
    if (!status) {
      status = this.createStatus(nodeId);
      this.statuses.set(nodeId, status);
    }

    status.totalFailures++;
    status.lastFailureTime = BigInt(Date.now()) * 1_000_000n;

    if (status.state === CircuitBreakerState.CLOSED) {
      // 실패 카운트 증가
      status.failureCount++;
      if (status.failureCount >= this.failureThreshold) {
        // CLOSED → OPEN
        this.transitionState(nodeId, CircuitBreakerState.OPEN);
      }
    } else if (status.state === CircuitBreakerState.HALF_OPEN) {
      // 복구 시도 실패: HALF_OPEN → OPEN
      this.transitionState(nodeId, CircuitBreakerState.OPEN);
    }
  }

  /**
   * 요청 허용 여부 확인
   * @returns true면 요청 통과, false면 차단
   */
  public canExecute(nodeId: bigint): boolean {
    let status = this.statuses.get(nodeId);
    if (!status) {
      status = this.createStatus(nodeId);
      this.statuses.set(nodeId, status);
    }

    if (status.state === CircuitBreakerState.CLOSED) {
      return true;
    }

    if (status.state === CircuitBreakerState.HALF_OPEN) {
      // HALF_OPEN: 제한된 요청만 허용 (1개씩)
      return true;
    }

    if (status.state === CircuitBreakerState.OPEN) {
      // OPEN: 타임아웃 경과 확인
      const currentTime = BigInt(Date.now()) * 1_000_000n;
      const timeoutNs = BigInt(this.timeout) * 1_000_000n;
      if (currentTime - status.lastTransitionTime >= timeoutNs) {
        // OPEN → HALF_OPEN
        this.transitionState(nodeId, CircuitBreakerState.HALF_OPEN);
        return true;
      }
      return false;
    }

    return false;
  }

  /**
   * 상태 조회
   */
  public getStatus(nodeId: bigint): CircuitBreakerStatus | undefined {
    return this.statuses.get(nodeId);
  }

  /**
   * 모든 상태 조회
   */
  public getAllStatuses(): Map<bigint, CircuitBreakerStatus> {
    return new Map(this.statuses);
  }

  /**
   * 상태 초기화
   */
  public reset(nodeId: bigint): void {
    const status = this.createStatus(nodeId);
    this.statuses.set(nodeId, status);
  }

  /**
   * 상태 전환
   */
  private transitionState(nodeId: bigint, newState: CircuitBreakerState): void {
    const status = this.statuses.get(nodeId);
    if (!status) return;

    const oldState = status.state;
    status.state = newState;
    status.lastTransitionTime = BigInt(Date.now()) * 1_000_000n;

    // 상태 전환 시 관련 카운트 초기화
    if (newState === CircuitBreakerState.CLOSED) {
      status.failureCount = 0;
      status.successCount = 0;
    } else if (newState === CircuitBreakerState.OPEN) {
      // (실패 카운트 유지)
    } else if (newState === CircuitBreakerState.HALF_OPEN) {
      status.successCount = 0;
    }

    console.log(
      `🔌 Circuit Breaker: Node ${nodeId} ${oldState} → ${newState}`
    );
  }

  /**
   * 상태 객체 생성
   */
  private createStatus(nodeId: bigint): CircuitBreakerStatus {
    return {
      nodeId,
      state: CircuitBreakerState.CLOSED,
      failureCount: 0,
      successCount: 0,
      lastFailureTime: 0n,
      lastTransitionTime: BigInt(Date.now()) * 1_000_000n,
      totalFailures: 0,
      totalSuccesses: 0
    };
  }

  /**
   * 통계 리포트 출력
   */
  public printReport(): void {
    console.log('\n⚡ Circuit Breaker Report\n');
    console.log('═'.repeat(70));

    let totalStatuses = 0;
    let closedCount = 0;
    let openCount = 0;
    let halfOpenCount = 0;

    for (const status of this.statuses.values()) {
      totalStatuses++;
      switch (status.state) {
        case CircuitBreakerState.CLOSED:
          closedCount++;
          break;
        case CircuitBreakerState.OPEN:
          openCount++;
          break;
        case CircuitBreakerState.HALF_OPEN:
          halfOpenCount++;
          break;
      }

      const successRate =
        status.totalFailures + status.totalSuccesses > 0
          ? ((status.totalSuccesses /
              (status.totalFailures + status.totalSuccesses)) *
              100)
              .toFixed(2)
          : 'N/A';

      console.log(`Node ${status.nodeId}:`);
      console.log(`  State:          ${status.state}`);
      console.log(`  Failures:       ${status.totalFailures}`);
      console.log(`  Successes:      ${status.totalSuccesses}`);
      console.log(`  Success Rate:   ${successRate}%`);
      console.log();
    }

    console.log('─'.repeat(70));
    console.log(`Total Nodes:    ${totalStatuses}`);
    console.log(`  CLOSED:       ${closedCount} (정상)`);
    console.log(`  OPEN:         ${openCount} (격리)`);
    console.log(`  HALF_OPEN:    ${halfOpenCount} (복구 중)`);
    console.log('═'.repeat(70));
  }
}
