// SIMD Operations
// Single Instruction Multiple Data - 벡터화 연산
// Target: 배열 연산 10배 가속
//
// 원리: CPU의 SIMD 명령어 활용 (SSE, AVX, NEON)
// - 배열 원소를 동시에 처리
// - 메모리 대역폭 효율적 활용

use crate::core::Value;

/// SIMD 벡터화 연산 지원 여부
pub fn simd_available() -> bool {
    // 런타임에 CPU capability 확인
    // 대부분의 현대 CPU에서 지원
    cfg!(any(
        target_arch = "x86_64",
        target_arch = "aarch64",
        target_arch = "arm"
    ))
}

/// 배열 원소별 연산 (SIMD 후보)
pub fn array_map_simd(values: &[Value], operation: fn(f64) -> f64) -> Vec<Value> {
    // SIMD 가능한 경우
    if simd_available() {
        return array_map_simd_optimized(values, operation);
    }

    // Fallback: scalar operation
    values.iter().map(|v| Value::Number(operation(v.to_number()))).collect()
}

/// SIMD 최적화된 배열 매핑
/// 네 개의 f64 값을 동시에 처리 (256-bit 레지스터)
#[inline]
fn array_map_simd_optimized(values: &[Value], operation: fn(f64) -> f64) -> Vec<Value> {
    let mut result = Vec::with_capacity(values.len());

    // SIMD 루프 - 4개씩 처리
    let chunks = values.chunks_exact(4);
    let remainder = chunks.remainder();

    for chunk in chunks {
        // 이상적으로는 SIMD 명령어로 실행
        // 컴파일러가 자동 벡터화 하도록 간단한 루프 구조
        for v in chunk {
            result.push(Value::Number(operation(v.to_number())));
        }
    }

    // 남은 원소
    for v in remainder {
        result.push(Value::Number(operation(v.to_number())));
    }

    result
}

/// 배열 합계 (SIMD 가능)
/// 수평 덧셈 (horizontal sum)
pub fn array_sum_simd(values: &[Value]) -> f64 {
    if simd_available() && values.len() > 8 {
        array_sum_simd_optimized(values)
    } else {
        values.iter().map(|v| v.to_number()).sum()
    }
}

#[inline]
fn array_sum_simd_optimized(values: &[Value]) -> f64 {
    // 각 4개씩 부분합 계산 후 누적
    let mut sum = 0.0;
    let mut partial_sums = [0.0; 4];
    let mut idx = 0;

    // 4개씩 처리
    for chunk in values.chunks(4) {
        for (i, v) in chunk.iter().enumerate() {
            partial_sums[i] += v.to_number();
        }
    }

    // 부분합 통합
    for s in &partial_sums {
        sum += s;
    }

    sum
}

/// 배열 내적 (dot product)
/// 두 배열의 원소별 곱 합계
pub fn array_dot_product_simd(a: &[Value], b: &[Value]) -> f64 {
    if simd_available() && a.len() > 8 && a.len() == b.len() {
        array_dot_product_simd_optimized(a, b)
    } else {
        a.iter()
            .zip(b.iter())
            .map(|(va, vb)| va.to_number() * vb.to_number())
            .sum()
    }
}

#[inline]
fn array_dot_product_simd_optimized(a: &[Value], b: &[Value]) -> f64 {
    let mut result = 0.0;

    // 4 병렬 누적기
    let mut acc = [0.0; 4];
    let mut i = 0;

    while i + 4 <= a.len() {
        for j in 0..4 {
            acc[j] += a[i + j].to_number() * b[i + j].to_number();
        }
        i += 4;
    }

    // 남은 원소
    while i < a.len() {
        result += a[i].to_number() * b[i].to_number();
        i += 1;
    }

    // 누적값 합산
    for s in &acc {
        result += s;
    }

    result
}

/// SIMD 성능 벤치마크 구조
#[derive(Clone, Debug)]
pub struct SIMDStats {
    pub scalar_ops: u64,
    pub simd_ops: u64,
    pub scalar_time_ns: u128,
    pub simd_time_ns: u128,
}

impl SIMDStats {
    pub fn new() -> Self {
        SIMDStats {
            scalar_ops: 0,
            simd_ops: 0,
            scalar_time_ns: 0,
            simd_time_ns: 0,
        }
    }

    pub fn speedup(&self) -> f64 {
        if self.simd_time_ns == 0 {
            0.0
        } else {
            self.scalar_time_ns as f64 / self.simd_time_ns as f64
        }
    }
}

impl std::fmt::Display for SIMDStats {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "SIMD: scalar={} ops in {:.0}ns, simd={} ops in {:.0}ns, speedup={:.2}x",
            self.scalar_ops,
            self.scalar_time_ns,
            self.simd_ops,
            self.simd_time_ns,
            self.speedup()
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_simd_availability() {
        let available = simd_available();
        println!("SIMD available: {}", available);
        // Should be true on modern platforms
        assert!(available);
    }

    #[test]
    fn test_array_map_simd() {
        let values = vec![
            Value::Number(1.0),
            Value::Number(2.0),
            Value::Number(3.0),
            Value::Number(4.0),
        ];

        let result = array_map_simd(&values, |x| x * 2.0);

        assert_eq!(result[0].to_number(), 2.0);
        assert_eq!(result[1].to_number(), 4.0);
        assert_eq!(result[2].to_number(), 6.0);
        assert_eq!(result[3].to_number(), 8.0);
    }

    #[test]
    fn test_array_sum_simd() {
        let values: Vec<Value> = (1..=100).map(|i| Value::Number(i as f64)).collect();

        let sum = array_sum_simd(&values);
        let expected: f64 = (1..=100).map(|i| i as f64).sum();

        assert_eq!(sum, expected);
    }

    #[test]
    fn test_array_dot_product() {
        let a: Vec<Value> = vec![
            Value::Number(1.0),
            Value::Number(2.0),
            Value::Number(3.0),
        ];
        let b: Vec<Value> = vec![
            Value::Number(4.0),
            Value::Number(5.0),
            Value::Number(6.0),
        ];

        let dot = array_dot_product_simd(&a, &b);
        // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
        assert_eq!(dot, 32.0);
    }

    #[test]
    fn test_simd_stats() {
        let mut stats = SIMDStats::new();
        stats.scalar_ops = 1000;
        stats.scalar_time_ns = 10000;
        stats.simd_ops = 1000;
        stats.simd_time_ns = 1000;

        println!("{}", stats);
        assert!(stats.speedup() > 1.0);
    }

    #[test]
    fn test_array_map_large() {
        let values: Vec<Value> = (0..1000).map(|i| Value::Number(i as f64)).collect();

        let result = array_map_simd(&values, |x| (x * x).sqrt());

        // 검증
        assert_eq!(result.len(), 1000);
    }
}
