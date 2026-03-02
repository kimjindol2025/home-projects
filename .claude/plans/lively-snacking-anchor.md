# Phase 4: 분산 벡터 인덱스 통합 API 레이어

## Context

Phase 3의 분산 벡터 인덱스 (Raft + Sharding + Replication)를
외부 클라이언트에게 노출하는 통합 레이어 구현.

기존 인프라 재사용:
- `src/runtime/websocket.fl` (439줄) → 실시간 스트리밍 API
- `src/grpc/server.fl` (525줄), `messages.fl` (195줄) → 고성능 RPC
- `src/serialization/protobuf.fl` (427줄) → 바이너리 직렬화
- `src/monitoring/metrics.fl` (367줄) → 분산 시스템 모니터링

## 목표

```
Layer 1: HybridIndexSystem (Phase 1-2)
Layer 2-4: 분산 클러스터 (Phase 3: Raft + Sharding + Replication)
Layer 5 (NEW): 통합 API 레이어 (Phase 4)
  ├─ WebSocket 실시간 스트리밍 (/ws/insert, /ws/search, /ws/cluster)
  ├─ gRPC 서비스 (VectorInsert, VectorSearch, ClusterStatus, BatchInsert, NodeHealth)
  ├─ Protocol Buffers 바이너리 직렬화 (벡터 도메인 전용)
  └─ 분산 시스템 통합 모니터링 (Raft + Replication + Sharding 메트릭)
```

## 구현 파일 (6개, ~2,700줄)

### 1. `src/distributed/websocket_stream.fl` (500줄)

**재사용 (src/runtime/websocket.fl)**:
- `WebSocketServer` struct - 연결 관리 (maxClients 10K)
- `WebSocketHandler` struct - 경로별 핸들러
- `WebSocketMessage` struct - 메시지 포맷
- `broadcastToPath(server, path, data)` - 경로별 브로드캐스트
- `sendToClient(server, clientId, data)` - 단일 클라이언트 전송

**새 구조체**:
```
struct VectorStreamServer
  wsServer: WebSocketServer
  coordinator: Coordinator       # Phase 3
  cluster: RaftIndexCluster      # Phase 3
  config: StreamServerConfig
  stats: StreamStats

struct StreamServerConfig
  maxConcurrentStreams: number   # 기본 1000
  searchBatchSize: number        # 기본 50
  heartbeatIntervalMs: number    # 기본 30000

struct StreamStats
  totalConnections: number
  activeStreams: number
  insertStreamed: number
  searchStreamed: number
```

**주요 함수**:
```
fn newVectorStreamServer(coordinator, cluster, config) -> VectorStreamServer
fn startVectorStreamServer(server) -> Result<VectorStreamServer, string>
  → 3개 경로 등록: /ws/insert, /ws/search, /ws/cluster

fn handleInsertStream(server, msg: WebSocketMessage) -> Result<WebSocketMessage, string>
  → JSON 파싱 (vectorId, vector[])
  → coordinator.routeInsertRequest(vectorId, vector)
  → 응답: { status: "ok", quorumAchieved, nodeIds[] }

fn handleSearchStream(server, msg: WebSocketMessage) -> Result<WebSocketMessage, string>
  → JSON 파싱 (query[], topK)
  → coordinator.routeSearchRequest(query, topK)
  → coordinator.aggregateSearchResults(partialResults, topK)
  → 응답: { results: ScoredVector[], totalMs }

fn handleClusterStream(server, msg) -> Result<WebSocketMessage, string>
  → coordinator.collectHealthReport()
  → 응답: ClusterHealthReport JSON

fn broadcastClusterUpdate(server, healthReport) -> Result<number, string>
fn getStreamServerStats(server) -> string
```

**WebSocket 경로**:
- `/ws/insert` → 실시간 벡터 삽입 (Quorum 결과 즉시 반환)
- `/ws/search` → 실시간 분산 검색 (글로벌 Top-K)
- `/ws/cluster` → 클러스터 상태 스트림

---

### 2. `src/distributed/grpc_vector_service.fl` (500줄)

**재사용 (src/grpc/)**:
- `GrpcServer` struct (server.fl:13) - 서버 인프라
- `GrpcMethodHandler` struct (server.fl:20) - 메서드 등록
- `GrpcCall` struct (server.fl:27) - 요청 컨텍스트
- `GrpcStatus` enum (messages.fl:101) - 17개 상태 코드
- `GrpcError` struct (messages.fl:107) - 에러 처리

**새 구조체**:
```
struct VectorGrpcServer
  grpcServer: GrpcServer         # 재사용
  coordinator: Coordinator        # Phase 3
  cluster: RaftIndexCluster       # Phase 3
  stats: VectorGrpcStats

struct VectorGrpcStats
  totalRpcCalls: number
  insertCalls: number
  searchCalls: number
  clusterStatusCalls: number
  batchInsertCalls: number
  failedCalls: number
```

**5개 RPC**:
```
fn rpcVectorInsert(server, call) -> Result<map<string,string>, GrpcError>
  → coordinator.routeInsertRequest(vectorId, vector)
  → Ok: { quorumAchieved: "true", nodeIds: "..." }
  → Err: GrpcStatus.INTERNAL / INVALID_ARGUMENT

fn rpcVectorSearch(server, call) -> Result<map<string,string>, GrpcError>
  → coordinator.routeSearchRequest(query, topK)
  → coordinator.aggregateSearchResults(...)
  → Ok: { results: JSON, totalMs: "..." }

fn rpcClusterStatus(server, call) -> Result<map<string,string>, GrpcError>
  → coordinator.collectHealthReport()
  → Ok: { activeNodes, deadNodes, overallHealth, issues }

fn rpcBatchInsert(server, call) -> Result<map<string,string>, GrpcError>
  → call.message["vectors"] (JSON 배열)
  → 각 벡터별 routeInsertRequest() 반복
  → Ok: { inserted, failed, quorumAchievements }

fn rpcNodeHealth(server, call) -> Result<map<string,string>, GrpcError>
  → coordinator.collectHealthReport()
  → replicationMgr.checkReplicationHealth()
  → Ok: 상세 노드별 상태

fn newVectorGrpcServer(coordinator, cluster, host, port) -> VectorGrpcServer
fn startVectorGrpcServer(server) -> Result<VectorGrpcServer, string>
fn getVectorGrpcStats(server) -> string
```

---

### 3. `src/distributed/proto_vector.fl` (400줄)

**재사용 (src/serialization/protobuf.fl)**:
- `encodeVarint(num)` (line 23) - 가변 정수 인코딩
- `encodeFieldTag(fieldNumber, wireType)` (line 73) - 필드 태그
- `encodeMessage(message, schema)` (line 153) - 메시지 직렬화
- `decodeMessage(bytes, schema)` (line 235) - 역직렬화
- `createProtoSchema(name, fields)` (line 332) - 스키마 생성
- `createProtoField(fieldNumber, name, fieldType, label)` (line 338) - 필드
- `compareSize(jsonSize, protoSize)` (line 380) - 압축률 비교

**6개 Proto 스키마**:
```
VectorInsertRequest:  vectorId(1), dimensions(2), values(3, repeated)
VectorInsertResponse: success(1), quorumAchieved(2), nodeCount(3)
VectorSearchRequest:  queryVector(1, repeated), topK(2), partitionId(3)
ScoredVectorProto:    vectorId(1), score(2)
VectorSearchResponse: results(1, repeated), totalSearchMs(2), partitionCount(3)
ClusterStatusResponse: activeNodes(1), deadNodes(2), overallHealth(3), topologyVersion(4)
```

**새 함수**:
```
fn buildVectorInsertSchema() -> ProtoSchema
fn buildVectorSearchSchema() -> ProtoSchema
fn buildClusterStatusSchema() -> ProtoSchema
fn serializeVectorInsertRequest(vectorId, vector[]) -> Result<array<number>, string>
fn deserializeVectorInsertRequest(bytes) -> Result<map<string,string>, string>
fn serializeVectorSearchResponse(results, totalMs) -> Result<array<number>, string>
fn deserializeVectorSearchResponse(bytes) -> Result<map<string,string>, string>
fn serializeClusterStatus(healthReport) -> Result<array<number>, string>
fn calcVectorCompressionRatio(originalVector) -> string   # compareSize 재사용
```

---

### 4. `src/distributed/distributed_monitor.fl` (400줄)

**재사용 (src/monitoring/metrics.fl)**:
- `Metrics` struct (line 10) - 요청 통계
- `RequestMetric` struct (line 18) - 개별 요청
- `PerformanceReport` struct (line 30) - 성능 보고서
- `recordRequest(metrics, metric)` (line 102) - 요청 기록
- `generatePerformanceReport(metrics)` (line 151) - 리포트
- `findBottlenecks(metrics)` (line 187) - 병목 감지

**새 구조체**:
```
struct DistributedMetrics
  apiMetrics: Metrics            # 재사용
  raftMetrics: RaftMetrics       # NEW
  replicationMetrics: ReplicationMetrics  # NEW
  shardingMetrics: ShardingMetrics  # NEW

struct RaftMetrics
  leaderElections: number
  logAppends: number
  snapshotsTaken: number
  avgCommitLatencyMs: number
  quorumSuccessRate: number

struct ReplicationMetrics
  quorumWrites: number
  quorumFailures: number
  failoversTriggered: number
  activeReplicas: number

struct ShardingMetrics
  rebalancesTriggered: number
  migrationsCompleted: number
  balanceScore: number
  avgPartitionLoad: number

struct DistributedHealthDashboard
  metrics: DistributedMetrics
  clusterHealth: ClusterHealthReport   # Phase 3 재사용
  replicationHealth: array<PartitionHealth>  # Phase 3 재사용
  timestamp: number
```

**새 함수**:
```
fn newDistributedMetrics() -> DistributedMetrics
fn recordRaftEvent(metrics, event: string) -> DistributedMetrics
  → event: "LEADER_ELECTION" | "LOG_APPEND" | "SNAPSHOT" | "COMMIT"
fn recordReplicationEvent(metrics, success: bool) -> DistributedMetrics
fn recordShardingEvent(metrics, event: string, balanceScore) -> DistributedMetrics
fn generateDistributedHealthDashboard(metrics, coordinator, cluster) -> DistributedHealthDashboard
fn reportDistributedMetrics(metrics) -> string
```

---

### 5. `tests/phase4_integration_test.fl` (500줄, 20개 테스트)

```
Group A: WebSocket 스트리밍 (4개)
  testWebSocketServerStartup     # 서버 시작
  testInsertStream               # WS → routeInsertRequest
  testSearchStream               # WS → routeSearchRequest → Top-K
  testClusterStatusBroadcast     # /ws/cluster 브로드캐스트

Group B: gRPC 서비스 (5개)
  testVectorInsertRpc            # VectorInsert RPC
  testVectorSearchRpc            # VectorSearch 결과 검증
  testClusterStatusRpc           # ClusterStatus 필드
  testBatchInsertRpc             # 10개 벡터 일괄
  testNodeHealthRpc              # 노드별 헬스

Group C: Protocol Buffers (4개)
  testVectorInsertSerialization  # 직렬화 → 역직렬화 → 동일성
  testSearchResponseSerialization
  testProtoCompressionRatio      # Proto < JSON × 0.5
  testSchemaValidation           # 잘못된 스키마 에러

Group D: 분산 모니터링 (4개)
  testRaftMetricsRecording       # LEADER_ELECTION 이벤트
  testReplicationMetricsTracking # Quorum 성공률 계산
  testShardingBalanceScore       # balanceScore 변화
  testHealthDashboardGeneration  # 전체 대시보드 생성

Group E: 통합 시나리오 (3개)
  testWebSocketGrpcInterop       # WS + gRPC 동일 Coordinator
  testProtoWebSocketStream       # Proto 직렬화 메시지 WS 전송
  testEnd2EndDistributedSearch   # Phase 3 → Phase 4 전체 파이프라인
```

---

### 6. `docs/PHASE4_GUIDE.md` (400줄)

- Phase 4 아키텍처 (5계층 ASCII 다이어그램)
- WebSocket API 레퍼런스 (경로, 메시지 포맷, 예시)
- gRPC 서비스 정의 (5개 RPC, 에러 코드)
- Protocol Buffers 스키마 (6개 메시지, Wire Type 설명)
- 모니터링 대시보드 (대시보드 해석, 알림 기준)
- 운영 가이드 (포트 설정, 클라이언트 연결 방법)
- 성능 벤치마크 (WebSocket vs REST vs gRPC 비교)

---

## 재사용 함수 맵

| Phase 4 파일 | 재사용 대상 | 위치 |
|-------------|-----------|------|
| websocket_stream.fl | WebSocketServer, WebSocketHandler | src/runtime/websocket.fl |
| websocket_stream.fl | broadcastToPath, sendToClient | websocket.fl:264 |
| grpc_vector_service.fl | GrpcServer, GrpcMethodHandler, GrpcCall | src/grpc/server.fl |
| grpc_vector_service.fl | GrpcStatus enum (17개), GrpcError | src/grpc/messages.fl:101 |
| proto_vector.fl | encodeVarint, encodeMessage, decodeMessage | src/serialization/protobuf.fl:23,153,235 |
| proto_vector.fl | createProtoSchema, createProtoField | protobuf.fl:332,338 |
| proto_vector.fl | compareSize | protobuf.fl:380 |
| distributed_monitor.fl | Metrics, recordRequest, findBottlenecks | src/monitoring/metrics.fl |
| ALL Phase 4 | coordinator, routeInsertRequest, routeSearchRequest | src/distributed/coordinator.fl |
| ALL Phase 4 | aggregateSearchResults, collectHealthReport | coordinator.fl:330,211 |
| ALL Phase 4 | checkReplicationHealth | replication.fl:374 |

---

## 구현 순서

```
1. proto_vector.fl         (의존성 없음, 스키마 먼저)
2. distributed_monitor.fl  (metrics.fl 패턴 따라 구현)
3. grpc_vector_service.fl  (proto_vector + grpc 재사용)
4. websocket_stream.fl     (coordinator + websocket 재사용)
5. phase4_integration_test.fl (20개 테스트)
6. PHASE4_GUIDE.md         (문서)
```

---

## 성능 목표

| 작업 | WebSocket | gRPC | REST (Phase 1) |
|------|-----------|------|----------------|
| 벡터 삽입 | 5-10ms | 5-10ms | - |
| 벡터 검색 | 0.2ms | 0.2ms | - |
| 직렬화 크기 | Proto (81% 압축) | Proto | JSON 100% |
| 동시 연결 | 10K | - | - |

## 검증 포인트

1. `testInsertStream` → WebSocket → coordinator.routeInsertRequest 연결
2. `testVectorSearchRpc` → gRPC → 분산 검색 글로벌 Top-K
3. `testProtoCompressionRatio` → Proto < JSON × 0.5
4. `testHealthDashboardGeneration` → Raft + Replication + Sharding 통합
5. `testEnd2EndDistributedSearch` → Phase 3 + Phase 4 전체 파이프라인
