// 🌐 Global Synapse Engine - Main Export

export { RDMAFabric } from './rdma_fabric';
export { SemanticSyncProtocol, SemanticSyncValidator } from './semantic_sync';
export { HashChain, HashChainVerifier } from './hash_chain';
export * from './types';

// Week 4: Auto Recovery Mechanisms
export {
  CircuitBreaker,
  CircuitBreakerState,
  CircuitBreakerConfig,
  CircuitBreakerMetrics,
  RetryStrategy,
  RetryPolicy,
  TimeoutManager,
  TimeoutPolicy,
  AutoRecoveryOrchestrator,
} from './auto_recovery';

export { ChaosMonkey, ChaosEvent, ChaosStats } from './chaos-monkey';

console.log('Global Synapse Engine - Phase 1 + Week 4 Auto-Recovery Loaded');
