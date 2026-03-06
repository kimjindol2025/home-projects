// 🔄 Retry Strategy - 지수 백오프 + 지터 (Week 4)

import { RetryConfig } from './types';

/**
 * 재시도 결과
 */
export interface RetryResult<T> {
  success: boolean;
  result?: T;
  error?: Error;
  attempts: number;
  totalDelayMs: number;
}

/**
 * RetryStrategy: 지수 백오프 + 지터를 사용한 재시도 관리
 *
 * 알고리즘:
 * 1. 첫 시도 (attempt 0)
 * 2. 실패 시 지수 백오프: delay = baseDelay * (backoffMultiplier ^ attempt)
 * 3. 지터 추가: delay = delay * (0.5 + random(0.5)) → [50%, 150%]
 * 4. 최대 지연 시간 제한
 * 5. 재시도 불가능한 오류는 즉시 반환
 *
 * 예시:
 * Attempt 0: 즉시 시도
 * Attempt 1: delay = 100ms * 2^1 * jitter ≈ 100-300ms
 * Attempt 2: delay = 100ms * 2^2 * jitter ≈ 200-600ms
 * Attempt 3: delay = 100ms * 2^3 * jitter ≈ 400-1200ms (제한 시간 적용)
 */
export class RetryStrategy {
  private readonly maxAttempts: number;
  private readonly baseDelayMs: number;
  private readonly maxDelayMs: number;
  private readonly jitterFactor: number;
  private readonly backoffMultiplier: number;
  private readonly retryableErrors: Set<string>;
  private totalRetries: number = 0;
  private totalSuccesses: number = 0;

  constructor(config: Partial<RetryConfig> = {}) {
    this.maxAttempts = config.maxAttempts ?? 3;
    this.baseDelayMs = config.baseDelayMs ?? 100;
    this.maxDelayMs = config.maxDelayMs ?? 5000;
    this.jitterFactor = config.jitterFactor ?? 0.5;
    this.backoffMultiplier = config.backoffMultiplier ?? 2;
    this.retryableErrors = new Set(
      config.retryableErrors ?? [
        'ETIMEDOUT',
        'ECONNREFUSED',
        'ECONNRESET',
        'EHOSTUNREACH',
        'ENETUNREACH'
      ]
    );
  }

  /**
   * 비동기 작업 재시도
   * @param operation 실행할 비동기 작업
   * @param operationName 작업 이름 (로깅용)
   * @returns 재시도 결과
   */
  public async execute<T>(
    operation: () => Promise<T>,
    operationName: string = 'Unknown'
  ): Promise<RetryResult<T>> {
    let lastError: Error | undefined;
    let totalDelayMs = 0;

    for (let attempt = 0; attempt < this.maxAttempts; attempt++) {
      try {
        const startTime = Date.now();
        const result = await operation();
        const elapsed = Date.now() - startTime;

        this.totalSuccesses++;
        return {
          success: true,
          result,
          attempts: attempt + 1,
          totalDelayMs
        };
      } catch (error) {
        lastError = error instanceof Error ? error : new Error(String(error));

        // 재시도 불가능한 에러 확인
        if (!this.isRetryable(lastError)) {
          this.totalRetries++;
          return {
            success: false,
            error: lastError,
            attempts: attempt + 1,
            totalDelayMs
          };
        }

        // 마지막 시도가 아니면 대기 후 재시도
        if (attempt < this.maxAttempts - 1) {
          const delayMs = this.calculateBackoffDelay(attempt);
          totalDelayMs += delayMs;
          console.log(
            `⏳ Retry ${operationName}: Attempt ${attempt + 1}/${this.maxAttempts} ` +
              `failed (${lastError.message}), waiting ${delayMs}ms...`
          );
          await this.delay(delayMs);
        }
      }
    }

    this.totalRetries++;
    return {
      success: false,
      error: lastError,
      attempts: this.maxAttempts,
      totalDelayMs
    };
  }

  /**
   * 동기 작업 재시도 (타임아웃 포함)
   */
  public async executeWithTimeout<T>(
    operation: () => Promise<T>,
    timeoutMs: number,
    operationName: string = 'Unknown'
  ): Promise<RetryResult<T>> {
    const executeWithTimeout = async () => {
      return Promise.race([
        operation(),
        new Promise<T>((_, reject) =>
          setTimeout(
            () => reject(new Error('ETIMEDOUT')),
            timeoutMs
          )
        )
      ]);
    };

    return this.execute(executeWithTimeout, operationName);
  }

  /**
   * 지수 백오프 지연 시간 계산
   * @param attempt 재시도 횟수 (0부터 시작)
   * @returns 지연 시간 (ms)
   */
  private calculateBackoffDelay(attempt: number): number {
    // 지수 백오프: delay = baseDelay * (multiplier ^ attempt)
    const exponentialDelay = this.baseDelayMs * Math.pow(this.backoffMultiplier, attempt);

    // 최대 지연 시간 제한
    const cappedDelay = Math.min(exponentialDelay, this.maxDelayMs);

    // 지터 추가: delay * (0.5 + random(0.5)) → [50%, 150%]
    const jitter = 0.5 + Math.random() * this.jitterFactor;
    const finalDelay = Math.ceil(cappedDelay * jitter);

    return finalDelay;
  }

  /**
   * 재시도 가능한 에러인지 확인
   */
  private isRetryable(error: Error): boolean {
    // 에러 메시지에서 에러 코드 추출
    const errorCode = error.message.split(':')[0].trim();
    return this.retryableErrors.has(errorCode) || this.retryableErrors.has(error.message);
  }

  /**
   * 지연 실행 (Promise 기반)
   */
  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * 통계 조회
   */
  public getStats() {
    return {
      totalRetries: this.totalRetries,
      totalSuccesses: this.totalSuccesses,
      retryRate:
        this.totalRetries + this.totalSuccesses > 0
          ? ((this.totalRetries / (this.totalRetries + this.totalSuccesses)) * 100)
              .toFixed(2) + '%'
          : 'N/A'
    };
  }

  /**
   * 통계 리포트 출력
   */
  public printReport(): void {
    console.log('\n🔄 Retry Strategy Report\n');
    console.log('═'.repeat(70));
    console.log(`Max Attempts:        ${this.maxAttempts}`);
    console.log(`Base Delay:          ${this.baseDelayMs}ms`);
    console.log(`Max Delay:           ${this.maxDelayMs}ms`);
    console.log(`Backoff Multiplier:  ${this.backoffMultiplier}x`);
    console.log(`Jitter Factor:       ${this.jitterFactor * 100}%`);
    console.log(`─`.repeat(70));
    console.log(`Total Retries:       ${this.totalRetries}`);
    console.log(`Total Successes:     ${this.totalSuccesses}`);
    console.log(
      `Retry Rate:          ${this.getStats().retryRate}`
    );
    console.log(`Retryable Errors:    ${Array.from(this.retryableErrors).join(', ')}`);
    console.log('═'.repeat(70));
  }

  /**
   * 통계 초기화 (테스트용)
   */
  public resetStats(): void {
    this.totalRetries = 0;
    this.totalSuccesses = 0;
  }
}

// 전역 RetryStrategy 인스턴스
let globalRetryStrategy: RetryStrategy | null = null;

export function getGlobalRetryStrategy(): RetryStrategy {
  if (!globalRetryStrategy) {
    globalRetryStrategy = new RetryStrategy();
  }
  return globalRetryStrategy;
}
