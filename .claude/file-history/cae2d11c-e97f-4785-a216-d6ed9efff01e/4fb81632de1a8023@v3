// ============================================
// Load Balancer & Proxy - Rust Implementation
// ============================================
// 고성능 로드 밸런싱 및 프록시

use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use dashmap::DashMap;
use tracing::info;

#[derive(Debug, Clone)]
pub enum LoadBalancingStrategy {
    RoundRobin,
    LeastConnections,
    Weighted,
}

pub struct Backend {
    pub id: usize,
    pub host: String,
    pub port: u16,
    pub weight: u32,
    pub active: bool,
    pub connection_count: Arc<AtomicUsize>,
}

pub struct ProxyServer {
    listen_port: u16,
    backends: Vec<Arc<Backend>>,
    current_index: Arc<AtomicUsize>,
    strategy: LoadBalancingStrategy,
    metrics: Arc<ProxyMetrics>,
}

pub struct ProxyMetrics {
    total_requests: Arc<AtomicUsize>,
    successful_requests: Arc<AtomicUsize>,
    failed_requests: Arc<AtomicUsize>,
    total_latency: Arc<AtomicUsize>,
}

impl ProxyServer {
    /// 새 프록시 서버 생성
    pub fn new(port: u16) -> Self {
        let backends = vec![
            Arc::new(Backend {
                id: 0,
                host: "localhost".to_string(),
                port: 9000,
                weight: 1,
                active: true,
                connection_count: Arc::new(AtomicUsize::new(0)),
            }),
            Arc::new(Backend {
                id: 1,
                host: "localhost".to_string(),
                port: 9100,
                weight: 2,
                active: true,
                connection_count: Arc::new(AtomicUsize::new(0)),
            }),
            Arc::new(Backend {
                id: 2,
                host: "localhost".to_string(),
                port: 9200,
                weight: 1,
                active: true,
                connection_count: Arc::new(AtomicUsize::new(0)),
            }),
        ];

        info!("🌐 Proxy Server initialized on port {}", port);
        info!("   Backends: {}", backends.len());
        info!("   Strategy: Round-Robin");

        ProxyServer {
            listen_port: port,
            backends,
            current_index: Arc::new(AtomicUsize::new(0)),
            strategy: LoadBalancingStrategy::RoundRobin,
            metrics: Arc::new(ProxyMetrics {
                total_requests: Arc::new(AtomicUsize::new(0)),
                successful_requests: Arc::new(AtomicUsize::new(0)),
                failed_requests: Arc::new(AtomicUsize::new(0)),
                total_latency: Arc::new(AtomicUsize::new(0)),
            }),
        }
    }

    /// 백엔드 선택 (Round-Robin)
    pub fn select_backend(&self) -> Option<Arc<Backend>> {
        if self.backends.is_empty() {
            return None;
        }

        let index = self.current_index.fetch_add(1, Ordering::Relaxed);
        let selected = index % self.backends.len();

        Some(self.backends[selected].clone())
    }

    /// 최소 연결 백엔드 선택
    pub fn select_backend_least_connections(&self) -> Option<Arc<Backend>> {
        self.backends
            .iter()
            .min_by_key(|b| b.connection_count.load(Ordering::Relaxed))
            .cloned()
    }

    /// 요청 처리
    pub async fn handle_request(&self, request_id: usize) -> Result<(), String> {
        // 요청 카운트 증가
        self.metrics.total_requests.fetch_add(1, Ordering::Relaxed);

        // 백엔드 선택
        let backend = self.select_backend()
            .ok_or_else(|| "No backends available".to_string())?;

        // 연결 증가
        backend.connection_count.fetch_add(1, Ordering::Relaxed);

        // 요청 처리 시뮬레이션 (10-50ms)
        let latency = 10 + (request_id % 40) as usize;
        tokio::time::sleep(tokio::time::Duration::from_millis(latency as u64)).await;

        // 연결 감소
        backend.connection_count.fetch_sub(1, Ordering::Relaxed);

        // 성공 카운트 증가 (95% 성공률)
        if request_id % 100 < 95 {
            self.metrics.successful_requests.fetch_add(1, Ordering::Relaxed);
        } else {
            self.metrics.failed_requests.fetch_add(1, Ordering::Relaxed);
            return Err("Request failed".to_string());
        }

        // 지연시간 기록
        self.metrics.total_latency.fetch_add(latency, Ordering::Relaxed);

        Ok(())
    }

    /// 메트릭 조회
    pub fn get_metrics(&self) -> ProxyMetricsSnapshot {
        let total = self.metrics.total_requests.load(Ordering::Relaxed);
        let successful = self.metrics.successful_requests.load(Ordering::Relaxed);
        let failed = self.metrics.failed_requests.load(Ordering::Relaxed);
        let total_latency = self.metrics.total_latency.load(Ordering::Relaxed);

        let avg_latency = if total > 0 {
            total_latency / total
        } else {
            0
        };

        ProxyMetricsSnapshot {
            total_requests: total,
            successful_requests: successful,
            failed_requests: failed,
            average_latency: avg_latency as u64,
        }
    }

    /// 상태 보고
    pub fn print_status(&self) {
        let metrics = self.get_metrics();

        info!("📊 Proxy Metrics:");
        info!("├─ Total Requests: {}", metrics.total_requests);
        info!("├─ Successful: {}", metrics.successful_requests);
        info!("├─ Failed: {}", metrics.failed_requests);
        info!("└─ Avg Latency: {}ms", metrics.average_latency);

        info!("📈 Backend Status:");
        for backend in &self.backends {
            info!(
                "  Node {}: {} connections (weight: {})",
                backend.id,
                backend.connection_count.load(Ordering::Relaxed),
                backend.weight
            );
        }
    }
}

impl Clone for ProxyServer {
    fn clone(&self) -> Self {
        ProxyServer {
            listen_port: self.listen_port,
            backends: self.backends.clone(),
            current_index: self.current_index.clone(),
            strategy: self.strategy.clone(),
            metrics: self.metrics.clone(),
        }
    }
}

impl Clone for LoadBalancingStrategy {
    fn clone(&self) -> Self {
        match self {
            LoadBalancingStrategy::RoundRobin => LoadBalancingStrategy::RoundRobin,
            LoadBalancingStrategy::LeastConnections => LoadBalancingStrategy::LeastConnections,
            LoadBalancingStrategy::Weighted => LoadBalancingStrategy::Weighted,
        }
    }
}

#[derive(Debug, Clone)]
pub struct ProxyMetricsSnapshot {
    pub total_requests: usize,
    pub successful_requests: usize,
    pub failed_requests: usize,
    pub average_latency: u64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_proxy_creation() {
        let proxy = ProxyServer::new(8080);
        assert_eq!(proxy.listen_port, 8080);
        assert_eq!(proxy.backends.len(), 3);
    }

    #[test]
    fn test_backend_selection() {
        let proxy = ProxyServer::new(8080);
        let b1 = proxy.select_backend();
        let b2 = proxy.select_backend();

        assert!(b1.is_some());
        assert!(b2.is_some());
        assert_ne!(b1.unwrap().id, b2.unwrap().id);
    }

    #[tokio::test]
    async fn test_request_handling() {
        let proxy = ProxyServer::new(8080);
        let result = proxy.handle_request(0).await;
        assert!(result.is_ok());

        let metrics = proxy.get_metrics();
        assert_eq!(metrics.total_requests, 1);
    }
}
