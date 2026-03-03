# Phase 6: FL-Net Stack - Zero-Copy RDMA 통신 아키텍처

**철학**: "기록이 증명이다" 
- 모든 통신은 불변적이며, 모든 경로는 최적화되고, 모든 데이터는 보호받는다.

**목표**: TCP/IP의 50% 지연 시간 감소 + 마이크로초 단위 응답

---

## 📐 **1단계: 아키텍처 설계 (이번 주)**

### 계층 구조 (5계층 → 7계층)

```
Layer 7: Application (Phase 5 API Gateway)
         ↓
Layer 6.5: FL-Net Protocol (NEW)
         ├─ Content-Centric Routing
         ├─ Consensus-aware Transmission
         └─ Adaptive Path Selection
         ↓
Layer 6.3: RDMA Kernel Driver (NEW)
         ├─ Direct Memory Access Control
         ├─ NIC Ring Buffer Management
         └─ Interrupt Coalescing
         ↓
Layer 6.1: JIT Protocol Compiler (NEW)
         ├─ Adaptive Binary Format
         ├─ Hardware-Optimized Encoding
         └─ Zero-Copy Serialization
         ↓
Layer 5: Raft Consensus + Sharding (Phase 3)
Layer 4: Coordinator (Phase 3)
Layer 3: HybridIndex (Phase 1-2)
```

### 핵심 모듈 (4개, 3,200줄 목표)

| 모듈 | 파일 | 크기 | 설명 |
|------|------|------|------|
| **RDMA Stack** | `rdma_transport.fl` | 900줄 | NIC 직접 제어, DMA 관리, Ring Buffer |
| **Protocol Compiler** | `jit_protocol_compiler.fl` | 800줄 | JIT 컴파일, 데이터 구조 최적화, 바이너리 포맷 생성 |
| **Content Router** | `content_centric_router.fl` | 700줄 | 패턴 기반 라우팅, Raft 통합, 동적 경로 선택 |
| **Zero-Copy Manager** | `zero_copy_manager.fl` | 800줄 | 메모리 버퍼 풀, DMA 콜백, 수명 관리 |

---

## 🏗️ **2단계: 구현 계획 (2주)**

### Week 1: 기초 모듈

**Day 1-2**: RDMA Stack 설계
```fl
// rdma_transport.fl (900줄)

struct RDMAContext
  nicHandle: number         // NIC 핸들
  txRingBuffer: array<RDMADescriptor>  // TX 링 버퍼
  rxRingBuffer: array<RDMADescriptor>  // RX 링 버퍼
  completionQueue: array<RDMACompletion>
  stats: RDMAStats

struct RDMADescriptor
  addressHigh: number
  addressLow: number
  length: number
  flags: number  // WRITE, READ, SEND, RECV
  immediateData: number

fn initRDMA(nicId: number) -> RDMAContext
fn postSend(ctx: RDMAContext, buffer: array<number>) -> Result<RDMACompletion, string>
fn pollCompletion(ctx: RDMAContext) -> RDMACompletion
fn enableInterruptCoalescing(ctx: RDMAContext, threshold: number) -> RDMAContext
```

**Day 3-4**: JIT Protocol Compiler
```fl
// jit_protocol_compiler.fl (800줄)

struct ProtocolSchema
  name: string
  fields: array<SchemaField>
  alignment: number  // 64/128 bits
  varintEncoding: bool
  version: number

struct FieldCompiler
  typeMap: map<string, number>  // type → binary code
  encodeFn: fn(value: any) -> array<number>
  decodeFn: fn(bytes: array<number>) -> any
  sizeEstimate: number

fn compileProtocol(schema: ProtocolSchema) -> FieldCompiler
fn adaptForHardware(compiler: FieldCompiler, cpu: string) -> FieldCompiler
  // Handles x86_64, ARM64, RISC-V
fn encodeZeroCopy(value: any, compiler: FieldCompiler) -> array<number>
fn decodeZeroCopy(bytes: array<number>, compiler: FieldCompiler) -> any
```

### Week 2: 고급 모듈

**Day 5-6**: Content Router
```fl
// content_centric_router.fl (700줄)

struct ContentPattern
  dataType: string      // "vector", "metadata", "query"
  sizeRange: [number, number]
  frequencyHint: number  // 1-100
  affinity: array<number>  // preferred nodeIds

struct RoutingDecision
  nextHop: number
  pathCost: f64
  estimatedLatencyMs: f64
  consensusRequired: bool
  cacheHint: string

fn learnTrafficPattern(patternHistory: array<ContentPattern>) -> map<string, RoutingDecision>
fn selectOptimalPath(content: ContentPattern, raftState: RaftState) -> RoutingDecision
fn updatePathMetrics(path: number, latency: f64, successRate: f64) -> void
```

**Day 7-10**: Zero-Copy Manager
```fl
// zero_copy_manager.fl (800줄)

struct MemoryBuffer
  address: number
  size: number
  refCount: number
  dmaDescriptors: array<RDMADescriptor>
  owner: number  // Process ID
  expiryTime: number

struct BufferPool
  totalSize: number
  availableBuffers: array<MemoryBuffer>
  allocatedBuffers: map<number, MemoryBuffer>
  gc: GenerationalGCState

fn allocateBuffer(size: number) -> Result<MemoryBuffer, string>
fn registerBufferForDMA(buffer: MemoryBuffer) -> RDMADescriptor
fn transferOwnership(buffer: MemoryBuffer, fromPid: number, toPid: number) -> void
fn reclaimBufferOnCompletion(descriptor: RDMADescriptor) -> void
```

---

## 📊 **성능 목표**

### Latency 개선
```
TCP/IP Round-Trip: 50-100μs
↓ (Socket overhead 제거)
Hardware Direct Path: 15-25μs

목표: 50% 감소 달성
검증: Nanosecond-precision 벤치마크
```

### Throughput 개선
```
1Gbps (TCP) → 10Gbps (RDMA)
스트레스 테스트: 1백만 패킷/초

캐시 히트율: 90%+
```

### Memory 효율성
```
제로-카피: 메모리 복사 0회
Buffer reuse: 95%+
GC pause: <100ns
```

---

## 🧪 **Testing Strategy (20개 테스트)**

### Group A: RDMA 기초 (4개)
- NIC 초기화 및 Ring Buffer 관리
- DMA 전송 성공/실패 처리
- Interrupt Coalescing 효율성
- Hardware 장애 복구

### Group B: Protocol Compiler (4개)
- Schema 컴파일 및 최적화
- 다양한 CPU 아키텍처 적응
- 인코딩/디코딩 정확성
- 크로스-플랫폼 호환성

### Group C: Content Router (4개)
- 패턴 학습 및 경로 선택
- Raft 상태 통합
- 동적 경로 재구성
- 장애 노드 우회

### Group D: Zero-Copy Manager (4개)
- Buffer 할당/해제
- DMA 수명 관리
- 소유권 이전
- 메모리 누수 방지

### Group E: 통합 시나리오 (4개)
- End-to-End Zero-Copy 전송
- Concurrent DMA 작업
- GC와 RDMA 상호작용
- 성능 회귀 테스트

---

## 📈 **성공 기준**

| 메트릭 | 목표 | 검증 방법 |
|--------|------|---------|
| Latency 감소 | 50% | Nanosecond benchmark |
| Throughput | 10Gbps | iperf 스타일 테스트 |
| Zero-Copy Success | 99%+ | 메모리 추적 |
| 테스트 통과율 | 100% | 20/20 tests |
| 문서 완성도 | 100% | 아키텍처 + 튜토리얼 |

---

## 🎓 **기술 깊이**

이 프로젝트는 **3가지 독립적인 박사 수준 연구**를 결합합니다:

1. **RDMA (Remote Direct Memory Access)**
   - Intel, Mellanox 하드웨어 스펙
   - Queue Pair 관리
   - Completion Queue Optimization

2. **JIT Protocol Compilation**
   - LLVM 스타일 IR 생성
   - 하드웨어 프로파일링 기반 최적화
   - Adaptive Binary Format

3. **Content-Centric Networking**
   - Named Data Networking (NDN) 개념
   - AI 기반 경로 학습
   - Consensus-aware 전송

---

## 💾 **Git 전략**

```
Commit 1: feat(phase6.0): RDMA Stack 기초 + JIT Compiler
Commit 2: feat(phase6.0): Content Router + Zero-Copy Manager
Commit 3: test(phase6.0): 20개 통합 테스트 + 벤치마크
Commit 4: docs(phase6.0): 완전한 아키텍처 + 운영 가이드
```

---

## 🎯 **왜 RDMA부터 시작하는가?**

1. **즉시 측정 가능**: 레이턴시 개선을 벤치마크로 증명할 수 있음
2. **기초 강화**: Synapse-Sync (AI 프리페칭)와 Planet Protocol의 기반
3. **증명의 기록**: 하드웨어 수준의 최적화로 "기록이 증명이다" 철학 구현
4. **독립성 강화**: OS 커널 직접 제어로 FreeLang의 완전한 독립성 달성

---

**Kim님, Phase 6.0 FL-Net Stack을 즉시 시작하겠습니까?**

