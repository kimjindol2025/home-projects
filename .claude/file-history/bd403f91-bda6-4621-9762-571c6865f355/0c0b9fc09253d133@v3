// 💾 Memory Profiler - Week 2
// 메모리 사용 추적 및 누수 감지

/**
 * MemorySample: 메모리 샘플
 */
export interface MemorySample {
  timestamp: bigint;
  rss: number;      // Resident Set Size (MB)
  heapTotal: number;
  heapUsed: number;
  external: number;
}

/**
 * MemoryTrend: 추이 분석
 */
export interface MemoryTrend {
  startRss: number;
  endRss: number;
  peakRss: number;
  minRss: number;
  delta: number;
  isLeaking: boolean;
}

/**
 * MemoryProfiler: 메모리 프로파일링
 *
 * 특징:
 * - 연속 샘플링
 * - 누수 감지 (선형 증가)
 * - 피크 추적
 */
export class MemoryProfiler {
  private samples: MemorySample[] = [];
  private sampleInterval: NodeJS.Timeout | null = null;
  private startTime: bigint;

  constructor(private intervalMs: number = 100) {
    this.startTime = BigInt(Date.now()) * 1_000_000n;
  }

  /**
   * 샘플링 시작
   */
  public start(): void {
    if (this.sampleInterval) return;

    this.sample(); // 초기 샘플

    this.sampleInterval = setInterval(() => {
      this.sample();
    }, this.intervalMs);
  }

  /**
   * 샘플링 중지
   */
  public stop(): void {
    if (this.sampleInterval) {
      clearInterval(this.sampleInterval);
      this.sampleInterval = null;
    }
  }

  /**
   * 메모리 샘플 수집
   */
  private sample(): void {
    const usage = process.memoryUsage();
    const sample: MemorySample = {
      timestamp: BigInt(Date.now()) * 1_000_000n,
      rss: usage.rss / 1024 / 1024,
      heapTotal: usage.heapTotal / 1024 / 1024,
      heapUsed: usage.heapUsed / 1024 / 1024,
      external: usage.external / 1024 / 1024
    };
    this.samples.push(sample);
  }

  /**
   * 추이 분석
   */
  public analyzeTrend(): MemoryTrend {
    if (this.samples.length === 0) {
      return {
        startRss: 0,
        endRss: 0,
        peakRss: 0,
        minRss: 0,
        delta: 0,
        isLeaking: false
      };
    }

    const rssValues = this.samples.map(s => s.rss);
    const startRss = rssValues[0];
    const endRss = rssValues[rssValues.length - 1];
    const peakRss = Math.max(...rssValues);
    const minRss = Math.min(...rssValues);
    const delta = endRss - startRss;

    // 누수 감지: 지속적인 증가 (95% 이상의 샘플이 증가)
    let increasingCount = 0;
    for (let i = 1; i < rssValues.length; i++) {
      if (rssValues[i] >= rssValues[i - 1]) {
        increasingCount++;
      }
    }
    const isLeaking = increasingCount / rssValues.length > 0.95;

    return {
      startRss,
      endRss,
      peakRss,
      minRss,
      delta,
      isLeaking
    };
  }

  /**
   * 샘플 개수
   */
  public getSampleCount(): number {
    return this.samples.length;
  }

  /**
   * 현재 메모리
   */
  public getCurrent(): MemorySample | null {
    return this.samples.length > 0 ? this.samples[this.samples.length - 1] : null;
  }

  /**
   * 리포트 출력
   */
  public printReport(): void {
    const trend = this.analyzeTrend();
    const current = this.getCurrent();

    console.log('\n💾 Memory Profile\n');
    console.log('═'.repeat(60));
    console.log(`Samples:       ${this.getSampleCount()}`);
    console.log(`Start RSS:     ${trend.startRss.toFixed(2)} MB`);
    console.log(`End RSS:       ${trend.endRss.toFixed(2)} MB`);
    console.log(`Peak RSS:      ${trend.peakRss.toFixed(2)} MB`);
    console.log(`Min RSS:       ${trend.minRss.toFixed(2)} MB`);
    console.log(`Delta:         ${trend.delta.toFixed(2)} MB ${trend.delta > 0 ? '⬆️' : '⬇️'}`);
    console.log(`Leak Status:   ${trend.isLeaking ? '⚠️ LEAKING' : '✅ OK'}`);

    if (current) {
      console.log(`\nCurrent:`);
      console.log(`  RSS:         ${current.rss.toFixed(2)} MB`);
      console.log(`  Heap Used:   ${current.heapUsed.toFixed(2)} MB`);
      console.log(`  Heap Total:  ${current.heapTotal.toFixed(2)} MB`);
      console.log(`  External:    ${current.external.toFixed(2)} MB`);
    }

    console.log('\n' + '═'.repeat(60));
  }

  /**
   * 데이터 내보내기
   */
  public toJSON() {
    return {
      samples: this.samples,
      trend: this.analyzeTrend(),
      sampleCount: this.samples.length
    };
  }
}

/**
 * 전역 프로파일러
 */
let globalProfiler: MemoryProfiler | null = null;

export function getGlobalProfiler(): MemoryProfiler {
  if (!globalProfiler) {
    globalProfiler = new MemoryProfiler();
  }
  return globalProfiler;
}
