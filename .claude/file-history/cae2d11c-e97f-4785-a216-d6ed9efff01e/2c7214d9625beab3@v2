// Performance measurement utilities
// Macros and utilities for profiling and benchmarking the runtime

use std::time::{Instant, Duration};

/// Performance counter for measuring execution time
#[derive(Clone, Debug)]
pub struct PerfCounter {
    pub name: String,
    pub count: u64,
    pub total_time: Duration,
    pub min_time: Duration,
    pub max_time: Duration,
}

impl PerfCounter {
    pub fn new(name: &str) -> Self {
        PerfCounter {
            name: name.to_string(),
            count: 0,
            total_time: Duration::ZERO,
            min_time: Duration::MAX,
            max_time: Duration::ZERO,
        }
    }

    pub fn record(&mut self, duration: Duration) {
        self.count += 1;
        self.total_time += duration;
        if duration < self.min_time {
            self.min_time = duration;
        }
        if duration > self.max_time {
            self.max_time = duration;
        }
    }

    pub fn avg_time_us(&self) -> f64 {
        if self.count == 0 {
            0.0
        } else {
            self.total_time.as_micros() as f64 / self.count as f64
        }
    }

    pub fn avg_time_ns(&self) -> f64 {
        if self.count == 0 {
            0.0
        } else {
            self.total_time.as_nanos() as f64 / self.count as f64
        }
    }

    pub fn total_time_ms(&self) -> f64 {
        self.total_time.as_secs_f64() * 1000.0
    }

    pub fn throughput(&self) -> f64 {
        if self.total_time.as_secs_f64() == 0.0 {
            0.0
        } else {
            self.count as f64 / self.total_time.as_secs_f64()
        }
    }
}

impl std::fmt::Display for PerfCounter {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}: {} calls, avg {:.2}μs ({:.0}ns), min {:.2}μs, max {:.2}μs, total {:.2}ms, throughput {:.0} ops/sec",
            self.name,
            self.count,
            self.avg_time_us(),
            self.avg_time_ns(),
            self.min_time.as_micros() as f64,
            self.max_time.as_micros() as f64,
            self.total_time_ms(),
            self.throughput()
        )
    }
}

/// Macro for timing a block of code
#[macro_export]
macro_rules! measure {
    ($counter:expr, $block:block) => {{
        let start = std::time::Instant::now();
        let result = $block;
        let elapsed = start.elapsed();
        $counter.record(elapsed);
        result
    }};
}

/// Macro for quick timing
#[macro_export]
macro_rules! time_it {
    ($name:expr, $block:block) => {{
        let start = std::time::Instant::now();
        let result = $block;
        let elapsed = start.elapsed();
        println!("{}: {:.2}μs ({:.0}ns)", $name, elapsed.as_micros() as f64, elapsed.as_nanos() as f64);
        result
    }};
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_perf_counter_creation() {
        let counter = PerfCounter::new("test");
        assert_eq!(counter.count, 0);
        assert_eq!(counter.name, "test");
    }

    #[test]
    fn test_perf_counter_recording() {
        let mut counter = PerfCounter::new("test");
        counter.record(Duration::from_micros(100));
        counter.record(Duration::from_micros(200));

        assert_eq!(counter.count, 2);
        assert_eq!(counter.min_time, Duration::from_micros(100));
        assert_eq!(counter.max_time, Duration::from_micros(200));
        assert!((counter.avg_time_us() - 150.0).abs() < 0.1);
    }

    #[test]
    fn test_perf_counter_throughput() {
        let mut counter = PerfCounter::new("test");
        for _ in 0..1000 {
            counter.record(Duration::from_micros(1));
        }
        // 1000 ops in 1000μs = 1 million ops/sec
        assert!((counter.throughput() - 1_000_000.0).abs() < 100_000.0);
    }
}
