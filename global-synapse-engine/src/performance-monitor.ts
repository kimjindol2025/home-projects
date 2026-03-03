// 📊 Performance Monitor - Week 2
// 실제 성능 데이터 수집 (Simulated 제거)

/**
 * PerformanceMetric: 단일 측정값
 */
export interface PerformanceMetric {
  name: string;
  durationNs: bigint;      // 나노초 단위
  memoryStartMb: number;
  memoryEndMb: number;
  memoryDeltaMb: number;
  timestamp: bigint;
  success: boolean;
}

/**
 * PerformanceStats: 통계 데이터
 */
export interface PerformanceStats {
  name: string;
  count: number;
  minNs: bigint;
  maxNs: bigint;
  avgNs: bigint;
  p50Ns: bigint;    // 중앙값
  p95Ns: bigint;    // 95 percentile (느린 작업의 영향)
  p99Ns: bigint;    // 99 percentile (최악의 경우)
  totalMemoryMb: number;
  peakMemoryMb: number;
  successRate: number;
}

/**
 * PerformanceMonitor: 실제 성능 측정
 *
 * 특징:
 * - nanosecond 정확도 (process.hrtime.bigint())
 * - 메모리 추적 (before/after)
 * - 백분위수 계산 (p50, p95, p99)
 * - JSON 내보내기
 */
export class PerformanceMonitor {
  private metrics: Map<string, PerformanceMetric[]> = new Map();
  private memoryBaseline: number = 0;

  constructor() {
    // 초기 메모리 측정
    if (global.gc) {
      global.gc();
    }
    this.memoryBaseline = process.memoryUsage().heapUsed / 1024 / 1024;
  }

  /**
   * 작업 시작
   */
  public startMeasure(name: string): MeasureHandle {
    const startTime = process.hrtime.bigint();
    const memoryStart = process.memoryUsage().heapUsed / 1024 / 1024;

    return {
      end: (success: boolean = true) => {
        const endTime = process.hrtime.bigint();
        const memoryEnd = process.memoryUsage().heapUsed / 1024 / 1024;

        const metric: PerformanceMetric = {
          name,
          durationNs: endTime - startTime,
          memoryStartMb: memoryStart,
          memoryEndMb: memoryEnd,
          memoryDeltaMb: memoryEnd - memoryStart,
          timestamp: BigInt(Date.now()) * 1_000_000n,
          success
        };

        if (!this.metrics.has(name)) {
          this.metrics.set(name, []);
        }
        this.metrics.get(name)!.push(metric);

        return metric;
      }
    };
  }

  /**
   * 통계 계산
   */
  public getStats(name: string): PerformanceStats | null {
    const metrics = this.metrics.get(name);
    if (!metrics || metrics.length === 0) {
      return null;
    }

    const durations = metrics.map(m => m.durationNs).sort((a, b) =>
      a < b ? -1 : a > b ? 1 : 0
    );
    const count = durations.length;
    const successCount = metrics.filter(m => m.success).length;

    return {
      name,
      count,
      minNs: durations[0],
      maxNs: durations[count - 1],
      avgNs: durations.reduce((a, b) => a + b, 0n) / BigInt(count),
      p50Ns: this.percentile(durations, 0.5),
      p95Ns: this.percentile(durations, 0.95),
      p99Ns: this.percentile(durations, 0.99),
      totalMemoryMb: metrics.reduce((sum, m) => sum + m.memoryDeltaMb, 0),
      peakMemoryMb: Math.max(...metrics.map(m => m.memoryEndMb)),
      successRate: (successCount / count) * 100
    };
  }

  /**
   * 모든 통계
   */
  public getAllStats(): PerformanceStats[] {
    const stats: PerformanceStats[] = [];
    for (const name of this.metrics.keys()) {
      const stat = this.getStats(name);
      if (stat) stats.push(stat);
    }
    return stats;
  }

  /**
   * JSON 내보내기
   */
  public toJSON() {
    return {
      timestamp: new Date().toISOString(),
      measurements: Array.from(this.metrics.entries()).map(([name, metrics]) => ({
        name,
        rawMetrics: metrics,
        stats: this.getStats(name)
      }))
    };
  }

  /**
   * 리포트 출력
   */
  public printReport() {
    console.log('\n📊 Performance Report\n');
    console.log('═'.repeat(80));

    const stats = this.getAllStats();
    for (const stat of stats) {
      console.log(`\n${stat.name}:`);
      console.log(`  Count:        ${stat.count}`);
      console.log(`  Success Rate: ${stat.successRate.toFixed(1)}%`);
      console.log(`  Min:          ${this.formatNs(stat.minNs)}`);
      console.log(`  P50:          ${this.formatNs(stat.p50Ns)}`);
      console.log(`  P95:          ${this.formatNs(stat.p95Ns)}`);
      console.log(`  P99:          ${this.formatNs(stat.p99Ns)}`);
      console.log(`  Max:          ${this.formatNs(stat.maxNs)}`);
      console.log(`  Avg:          ${this.formatNs(stat.avgNs)}`);
      console.log(`  Memory Δ:     ${stat.totalMemoryMb.toFixed(2)} MB`);
      console.log(`  Peak Memory:  ${stat.peakMemoryMb.toFixed(2)} MB`);
    }

    console.log('\n' + '═'.repeat(80));
  }

  /**
   * 백분위수 계산
   */
  private percentile(sorted: bigint[], p: number): bigint {
    const index = Math.ceil(sorted.length * p) - 1;
    return sorted[Math.max(0, index)];
  }

  /**
   * nanosecond → 읽기 쉬운 형식
   */
  private formatNs(ns: bigint): string {
    const nsNum = Number(ns);
    if (nsNum < 1000) {
      return `${nsNum.toFixed(0)}ns`;
    } else if (nsNum < 1_000_000) {
      return `${(nsNum / 1000).toFixed(2)}μs`;
    } else if (nsNum < 1_000_000_000) {
      return `${(nsNum / 1_000_000).toFixed(2)}ms`;
    } else {
      return `${(nsNum / 1_000_000_000).toFixed(2)}s`;
    }
  }

  /**
   * 메모리 포맷
   */
  public getMemoryReport(): string {
    const usage = process.memoryUsage();
    return `RSS: ${(usage.rss / 1024 / 1024).toFixed(1)}MB, Heap: ${(usage.heapUsed / 1024 / 1024).toFixed(1)}MB/${(usage.heapTotal / 1024 / 1024).toFixed(1)}MB`;
  }
}

/**
 * 측정 핸들
 */
export interface MeasureHandle {
  end(success?: boolean): PerformanceMetric;
}

/**
 * 전역 모니터 인스턴스
 */
let globalMonitor: PerformanceMonitor | null = null;

export function getGlobalMonitor(): PerformanceMonitor {
  if (!globalMonitor) {
    globalMonitor = new PerformanceMonitor();
  }
  return globalMonitor;
}

/**
 * 편의 함수: 비동기 작업 측정
 */
export async function measureAsync<T>(
  name: string,
  fn: () => Promise<T>
): Promise<T> {
  const monitor = getGlobalMonitor();
  const handle = monitor.startMeasure(name);

  try {
    const result = await fn();
    handle.end(true);
    return result;
  } catch (error) {
    handle.end(false);
    throw error;
  }
}

/**
 * 편의 함수: 동기 작업 측정
 */
export function measureSync<T>(name: string, fn: () => T): T {
  const monitor = getGlobalMonitor();
  const handle = monitor.startMeasure(name);

  try {
    const result = fn();
    handle.end(true);
    return result;
  } catch (error) {
    handle.end(false);
    throw error;
  }
}
