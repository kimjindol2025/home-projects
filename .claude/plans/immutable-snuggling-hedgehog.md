# v7.2 마이크로서비스 아키텍처(MSA)와 gRPC - 구현 계획

## Context

v7.1 메시지 큐 (느슨한 결합, 비동기 통신) 완성 후, 석사 과정 두 번째 연구.
목표: gRPC 기반 전자상거래 MSA 구현으로 서비스 간 직접 고성능 통신을 이해.

**사용자 선택**
- 도메인: 전자상거래 (주문-결제-배송)
- Protobuf 전략: .proto 파일로 완전 자동 코드 생성 (grpcio-tools)
- 고급 기능: 양방향 스트리밍, 서비스 디스커버리, Load Balancing, 인터셉터

---

## 파일 구조

```
src/v7_2_grpc/
├── __init__.py
├── protos/
│   ├── order.proto            # 4가지 RPC 패턴 정의
│   ├── payment.proto
│   └── delivery.proto
├── generated/                 # grpcio-tools 자동 생성
│   ├── __init__.py
│   ├── order_pb2.py           (자동 생성)
│   ├── order_pb2_grpc.py      (자동 생성)
│   ├── payment_pb2.py         (자동 생성)
│   ├── payment_pb2_grpc.py    (자동 생성)
│   ├── delivery_pb2.py        (자동 생성)
│   └── delivery_pb2_grpc.py   (자동 생성)
├── services/
│   ├── order_service.py       # OrderBusinessLogic + OrderServicer
│   ├── payment_service.py     # PaymentBusinessLogic + PaymentServicer
│   └── delivery_service.py    # DeliveryBusinessLogic + DeliveryServicer
├── interceptors/
│   └── grpc_interceptors.py   # Logging, Auth, Performance
├── discovery/
│   └── service_registry.py    # GogsServiceRegistry + ServiceInstance
├── load_balancer/
│   └── load_balancer.py       # GogsLoadBalancer (RoundRobin/Weighted/LeastConn)
├── client/
│   └── ecommerce_client.py    # 통합 클라이언트 (grpc/local 모드)
└── main.py                    # 10개 교육 섹션

tests/v7_2/
├── __init__.py
├── test_order_service.py      # OrderBusinessLogic 테스트 (grpc 미사용)
├── test_payment_service.py
├── test_delivery_service.py
├── test_interceptors.py       # LoggingInterceptor, Auth, Perf 테스트
└── test_integration.py        # EcommerceClient(local 모드) 통합 테스트

docs/V7_2_GRPC_MSA_GUIDE.md
```

---

## 구현 순서

### Phase 1: Proto 파일 + 자동 코드 생성

**Step 1** `src/v7_2_grpc/__init__.py`
**Step 2-4** `.proto` 파일 3개 (각각 Unary, ServerStream, ClientStream, Bidirectional 4가지 RPC 포함)
**Step 5** `pip install grpcio grpcio-tools` 실행
**Step 6** `python -m grpc_tools.protoc` → `generated/` 자동 생성
**Step 7** `src/v7_2_grpc/generated/__init__.py`

### Phase 2: 서비스 구현 (2계층 패턴)

각 서비스는 계층 분리:
- `BusinessLogic` 클래스: grpc import 없음, 순수 Python, 테스트 용이
- `Servicer` 클래스: BusinessLogic 위임, grpc 어댑터 (`GRPC_AVAILABLE` 플래그)

**Step 8** `services/order_service.py`
```python
class OrderBusinessLogic:
    def create_order(order_id, product, qty, price) -> dict
    def track_order_stream(order_id) -> Generator[dict]       # Server Streaming
    def process_batch_orders(orders: list) -> dict            # Client Streaming
    def process_orders_bidirectional(requests) -> Generator   # Bidirectional

class OrderServicer(order_pb2_grpc.Servicer if GRPC_AVAILABLE else object):
    pass  # BusinessLogic 위임
```

**Step 9** `services/payment_service.py` - 동일 패턴
**Step 10** `services/delivery_service.py` - 동일 패턴

### Phase 3: 인터셉터

**Step 11** `interceptors/grpc_interceptors.py`
```python
class LoggingInterceptor:    # call_count, call_log, log_call()
class AuthInterceptor:        # VALID_KEYS = {gogs-key-2026, ...}, validate_api_key()
class PerformanceInterceptor: # measurements dict, measure(), get_stats()
```

### Phase 4: 서비스 디스커버리 + 로드 밸런서

**Step 12** `discovery/service_registry.py`
```python
@dataclass class ServiceInstance:   # host, port, status, weight, is_alive(), heartbeat()
class GogsServiceRegistry:           # register, deregister, discover, cleanup_stale
```

**Step 13** `load_balancer/load_balancer.py`
```python
class BalancingStrategy(Enum): ROUND_ROBIN, WEIGHTED, LEAST_CONNECTIONS, RANDOM
class GogsLoadBalancer:        # get_instance, _round_robin, _weighted, _least_connections
```

### Phase 5: 통합 클라이언트

**Step 14** `client/ecommerce_client.py`
- `"grpc"` 모드: 실제 gRPC 채널
- `"local"` 모드: BusinessLogic 직접 호출 (테스트/개발용)
- 메서드: `create_order`, `track_order_stream`, `process_payment`, `create_delivery`
- `get_stats()`: 인터셉터 + LB 통계 합산

### Phase 6: 교육 모듈 main.py (10개 섹션)

**Step 15** `main.py`
1. gRPC란? (REST vs gRPC 비교 다이어그램)
2. Protocol Buffers (Protobuf 문법 설명)
3. 4가지 RPC 패턴 (구조 다이어그램)
4. 서비스 구현 패턴 (2계층 분리)
5. 인터셉터 체인 (파이프라인)
6. 서비스 디스커버리 (TTL, 헬스체크)
7. 로드 밸런싱 전략 비교
8. 전자상거래 워크플로우 (주문→결제→배송)
9. 성능 분석 (gRPC vs REST)
10. 실전 MSA 아키텍처

### Phase 7: 테스트

**Step 16** `tests/v7_2/__init__.py`
**Step 17** `test_order_service.py` (Layer 1, 6 tests)
**Step 18** `test_payment_service.py` (Layer 1, 5 tests)
**Step 19** `test_delivery_service.py` (Layer 1, 5 tests)
**Step 20** `test_interceptors.py` (Layer 2, 8 tests)
**Step 21** `test_integration.py` (Layer 3, 7 tests)

### Phase 8: 문서

**Step 22** `docs/V7_2_GRPC_MSA_GUIDE.md` (400+ 줄)

---

## 핵심 설계 원칙

### 1. grpcio 설치 실패 대비
```python
try:
    import grpc
    from ..generated import order_pb2, order_pb2_grpc
    GRPC_AVAILABLE = True
except ImportError:
    GRPC_AVAILABLE = False
```
모든 31개+ 테스트가 grpc 없어도 통과 가능.

### 2. .proto 자동 생성 워크플로우
```bash
python -m grpc_tools.protoc \
    -I src/v7_2_grpc/protos \
    --python_out=src/v7_2_grpc/generated \
    --grpc_python_out=src/v7_2_grpc/generated \
    src/v7_2_grpc/protos/order.proto \
    src/v7_2_grpc/protos/payment.proto \
    src/v7_2_grpc/protos/delivery.proto
```

### 3. 테스트 3계층
- Layer 1: BusinessLogic 단독 (grpc 없음)
- Layer 2: 인터셉터/레지스트리/LB 단독
- Layer 3: EcommerceClient(mode="local") 통합 (서버 없음)

---

## 코드 통계 목표

| 항목 | 줄 수 |
|------|-------|
| protos/ (3개) | 150 |
| generated/ (6개, 자동) | ~600 |
| services/ (3개) | 600 |
| interceptors/ | 200 |
| discovery/ | 180 |
| load_balancer/ | 150 |
| client/ | 200 |
| main.py | 800 |
| tests/ (5개) | 600 |
| docs/ | 400 |
| **총합** | **~3,880줄** |

테스트: **31개 이상** (100% PASS)

---

## 검증 방법

```bash
pip install grpcio grpcio-tools
python -m grpc_tools.protoc -I src/v7_2_grpc/protos \
  --python_out=src/v7_2_grpc/generated \
  --grpc_python_out=src/v7_2_grpc/generated \
  src/v7_2_grpc/protos/*.proto
python -m pytest tests/v7_2/ -v   # 31개+ PASS
python -m src.v7_2_grpc.main       # 10개 섹션 출력
git commit + push to gogs
```

저장 필수: 너는 기록이 증명이다 gogs 👑
