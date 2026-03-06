# 🌐 FreeLang HTTP Engine

**상태**: 🚀 Phase B 구현 중
**버전**: 0.1.0
**목표**: TCP 소켓 + HTTP/1.1 파서 완성 (1,000줄)

---

## 📋 개요

HTTP/1.1 웹 서버 핵심 엔진입니다.

**특징**:
- ✅ 순수 FreeLang (100%)
- ✅ TCP 소켓 직접 제어
- ✅ HTTP/1.1 파싱
- ✅ Keep-Alive 지원
- ✅ 멀티플렉싱

**의존성**: 0 (커널만)

---

## 📁 구조

```
freelang-http-engine/
├── src/
│   ├── tcp_socket.fl      (250줄) - TCP 관리
│   ├── http_parser.fl     (350줄) - 요청 파싱
│   ├── http_handler.fl    (200줄) - 응답 생성
│   ├── server.fl          (200줄) - Server loop
│   └── mod.fl             (50줄)  - 공개 API
├── tests/
│   ├── test_tcp_socket.fl (150줄)
│   ├── test_http_parser.fl (250줄)
│   ├── test_server.fl     (150줄)
│   └── benchmark.fl       (100줄)
├── docs/
│   └── HTTP-ENGINE-DESIGN.md
└── README.md (이 파일)
```

---

## 🚀 빠른 시작

### 컴파일

```bash
cd freelang-http-engine
freelang compile src/server.fl -o http-server
```

### 실행

```bash
./http-server
# 🚀 Server listening on 127.0.0.1:8080
```

### 테스트

```bash
curl -X GET http://localhost:8080/health
# {"status":"ok"}

curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test"}'
# {"id":"1","status":"created"}
```

---

## 📊 구현 현황

| 모듈 | 상태 | 줄 수 |
|------|------|-------|
| TCP Socket Manager | ✅ 완료 | 250줄 |
| HTTP Parser | ✅ 완료 | 350줄 |
| HTTP Handler | ✅ 완료 | 200줄 |
| HTTP Server | ✅ 완료 | 200줄 |
| mod.fl | ✅ 완료 | 50줄 |
| **합계** | **✅ 완료** | **1,050줄** |

---

## 🧪 테스트

### 실행

```bash
freelang test tests/test_tcp_socket.fl
freelang test tests/test_http_parser.fl
freelang test tests/test_server.fl
freelang run tests/benchmark.fl
```

### 테스트 커버리지

| Suite | 테스트 수 | 상태 |
|-------|----------|------|
| TCP Socket | 5개 | ⏳ 준비 |
| HTTP Parser | 12개 | ⏳ 준비 |
| Server | 5개 | ⏳ 준비 |
| Benchmark | 3개 | ⏳ 준비 |
| **합계** | **25개** | **⏳ 준비** |

---

## 📈 성능 목표

| 지표 | 목표 | 측정 방법 |
|------|------|----------|
| Throughput | > 10K req/sec | wrk |
| Latency (p99) | < 50ms | histogram |
| Keep-Alive overhead | < 1ms | 연속 요청 |
| Parser | < 100µs | 평균 시간 |
| Memory (1K conn) | < 10MB | RSS |

---

## 🔗 Integration

### Input
- freelang-rest-api (라우팅)

### Output
- Layer 2: Raft DB (저장소)

---

## 📝 주요 함수

### TCP Socket

```freelang
fn socket_create() -> Result<i32, string>
fn socket_bind(fd, addr, port) -> Result<void, string>
fn socket_listen(fd, backlog) -> Result<void, string>
fn socket_accept(fd) -> Result<i32, string>
fn socket_read(fd, max_len) -> Result<array<u8>, string>
fn socket_write(fd, data) -> Result<i32, string>
```

### HTTP Parser

```freelang
fn parse_request(raw: array<u8>) -> Result<HttpRequest, ParseError>
fn parse_request_line(line) -> Result<(string, string, string), ParseError>
fn parse_headers_section(lines, start) -> Result<(map, i32), ParseError>
```

### HTTP Handler

```freelang
fn build_response(status, headers, body) -> HttpResponse
fn serialize_response(resp) -> array<u8>
fn build_error_response(status, message) -> HttpResponse
fn build_success_response(data) -> HttpResponse
```

### HTTP Server

```freelang
fn http_server_main(addr, port, handler) -> Result<void, string>
fn handle_client_connection(fd, handler) -> void
fn read_http_request(fd) -> Result<HttpRequest, string>
fn send_http_response(fd, response) -> Result<void, string>
```

---

## 🎯 다음 단계

- [ ] TCP Socket 테스트 작성 및 통과
- [ ] HTTP Parser 테스트 작성 및 통과
- [ ] Server 통합 테스트
- [ ] 성능 벤치마크
- [ ] GOGS 커밋

---

## 🔐 설계 원칙

1. **Zero-copy** (가능 범위)
2. **µs 단위 성능**
3. **커널 제어권**
4. **100% FreeLang**

---

**목표**: Layer 1 (TCP/HTTP) 완성 → Phase C (Production) → Phase A (통합)
