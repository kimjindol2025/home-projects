# 🌐 FreeLang HTTP Engine - 상세 설계 문서

**프로젝트**: freelang-http-engine
**버전**: 1.0
**상태**: 🚀 구현 중 (Phase B)
**목표**: TCP 소켓 + HTTP/1.1 파서 완성 (1,000줄)

---

## 📋 개요

**목적**: Sovereign Backend의 Layer 1 (TCP/HTTP) 완성

**구현 범위**:
- TCP 소켓 관리
- HTTP/1.1 요청 파싱
- HTTP/1.1 응답 생성
- Keep-Alive 지원
- 멀티플렉싱 (다중 동시 연결)

**언어**: FreeLang 100%
**의존성**: 0 (커널만 의존)

---

## 🏗️ 아키텍처

```
┌─────────────────────────────────────────┐
│         HTTP Server (main)              │
│  (라우터와 연동, accept loop)            │
├─────────────────────────────────────────┤
│      HTTP Parser & Handler              │
│  (요청 파싱, 응답 생성)                  │
├─────────────────────────────────────────┤
│      TCP Socket Manager                 │
│  (소켓 생성, read/write, timeout)       │
├─────────────────────────────────────────┤
│           Kernel                        │
│  (socket(), bind(), listen() 등)       │
└─────────────────────────────────────────┘
```

---

## 📂 파일 구조

```
freelang-http-engine/
├── src/
│   ├── tcp_socket.fl          (250줄) - TCP 소켓 관리
│   ├── http_parser.fl         (350줄) - HTTP 요청 파싱
│   ├── http_handler.fl        (200줄) - 응답 생성
│   ├── server.fl              (200줄) - Server loop
│   └── mod.fl                 (50줄)  - 모듈 공개 API
├── tests/
│   ├── test_tcp_socket.fl     (150줄) - TCP 테스트
│   ├── test_http_parser.fl    (250줄) - HTTP 파서 테스트
│   ├── test_server.fl         (150줄) - 통합 테스트
│   └── benchmark.fl           (100줄) - 성능 벤치마크
├── docs/
│   └── HTTP-ENGINE-DESIGN.md  (이 파일)
└── README.md                  (50줄)
```

---

## 📐 Module 1: TCP Socket Manager (250줄)

### 목적
저수준 소켓 작업 추상화

### 구조
```freelang
struct TcpSocket {
  fd: i32,              # 파일 디스크립터
  addr: string,         # 바인드 주소
  port: i32,            # 포트
  is_listening: bool,   # listen() 상태
  timeout_ms: i32,      # 타임아웃
}

struct SocketBuffer {
  data: array<u8>,      # 읽은 데이터
  offset: i32,          # 현재 위치
  size: i32,            # 버퍼 크기
}
```

### 주요 함수

| 함수 | 설명 | 반환 |
|------|------|------|
| `socket_create()` | 소켓 생성 | `Result<TcpSocket>` |
| `socket_bind(sock, addr, port)` | bind() | `Result<void>` |
| `socket_listen(sock, backlog)` | listen() | `Result<void>` |
| `socket_accept(sock)` | accept() (블로킹) | `Result<TcpSocket>` |
| `socket_read(sock, len)` | read() | `Result<array<u8>>` |
| `socket_write(sock, data)` | write() | `Result<i32>` |
| `socket_close(sock)` | close() | `Result<void>` |

### 무관용 규칙

**R1**: 연결당 throughput > 10K req/sec
- 측정: 1000개 요청, 단일 연결
- 예상 시간: 100ms 이상

**R2**: Keep-Alive 지연 < 1ms
- 측정: 연속 2개 요청, 동일 연결
- 포함: TCP 재사용 오버헤드

**R3**: HTTP 파싱 < 100µs
- 측정: 평균 파싱 시간
- 대상: 1000 바이트 요청

**R4**: 메모리 < 10MB (1000 동시 연결)
- 측정: RSS (상주 메모리)
- 계산: 연결당 ~10KB

**R5**: Zero-copy 전송 (가능 범위)
- 측정: sendfile() 사용률
- 목표: > 50% 정적 파일

---

## 📐 Module 2: HTTP Parser (350줄)

### 목적
HTTP/1.1 요청 파싱

### 구조
```freelang
struct HttpRequest {
  method: string,           # GET, POST, PUT, DELETE
  path: string,             # /todos/123
  version: string,          # HTTP/1.1
  headers: map<string, string>,  # 헤더
  body: array<u8>,          # 요청 본문
  body_len: i32,            # Content-Length
}

struct HttpResponse {
  version: string,          # HTTP/1.1
  status: i32,              # 200, 404, 500
  reason: string,           # OK, Not Found
  headers: map<string, string>,
  body: array<u8>,
}

enum ParseState {
  RequestLine,
  Headers,
  Body,
  Complete,
}
```

### 파싱 프로세스

```
1. Request Line 파싱
   "GET /todos HTTP/1.1\r\n"
   → method=GET, path=/todos, version=HTTP/1.1

2. Headers 파싱
   "Content-Type: application/json\r\n"
   "Content-Length: 42\r\n"
   "\r\n"
   → headers map + body_len = 42

3. Body 읽기
   body_len 바이트만큼 읽음
   → body = [...]

4. 상태: Complete
```

### 주요 함수

| 함수 | 설명 |
|------|------|
| `parse_request_line(line)` | Request line 파싱 |
| `parse_headers(lines)` | 헤더 파싱 |
| `parse_body(buffer, len)` | 본문 읽기 |
| `validate_request(req)` | 유효성 검증 |

### 무관용 규칙 (위 R3 참고)

**예외 처리**:
- 악의적 요청: 400 Bad Request
- 대용량 요청: 413 Payload Too Large (limit: 1MB)
- 느린 클라이언트: 408 Request Timeout (30s)

---

## 📐 Module 3: HTTP Handler (200줄)

### 목적
응답 생성 + 상태 코드 변환

### 구조
```freelang
fn build_response(
  status: i32,
  headers: map<string, string>,
  body: array<u8>
) -> HttpResponse
  # HTTP/1.1 응답 생성

fn serialize_response(resp: HttpResponse) -> array<u8>
  # 응답을 바이트로 변환
  # "HTTP/1.1 200 OK\r\n..."

fn build_status_line(status: i32) -> string
  # 상태 코드 → 이유 구문
  # 200 → "200 OK"
  # 404 → "404 Not Found"
```

### 상태 코드 지원

| 코드 | 이유 | 용도 |
|------|------|------|
| 200 | OK | 성공 |
| 201 | Created | 생성됨 |
| 204 | No Content | 삭제됨 |
| 400 | Bad Request | 잘못된 요청 |
| 404 | Not Found | 찾을 수 없음 |
| 500 | Internal Server Error | 서버 에러 |

---

## 📐 Module 4: HTTP Server (200줄)

### 목적
TCP 서버 + 요청 디스패칭

### Server Loop

```freelang
fn server_main(addr: string, port: i32) -> void
  1. server_socket = socket_create()
  2. socket_bind(server_socket, addr, port)
  3. socket_listen(server_socket, backlog=128)

  4. while true:
     a. client_socket = socket_accept(server_socket)
     b. handle_client(client_socket)
     c. socket_close(client_socket)
```

### 요청 처리

```freelang
fn handle_client(client_socket: TcpSocket) -> void
  1. buffer = socket_read(client_socket, 4096)
  2. request = parse_request(buffer)
  3. response = route_request(request)  # Layer 5 REST API 연동
  4. response_bytes = serialize_response(response)
  5. socket_write(client_socket, response_bytes)
```

### Keep-Alive 지원

```
요청 1: GET /todos
응답 1: 200 OK (Connection: keep-alive)
[동일 TCP 연결 유지]
요청 2: POST /todos
응답 2: 201 Created (Connection: keep-alive)
[연결 닫기 또는 타임아웃]
```

---

## 🧪 테스트 전략

### Test Suite 1: TCP Socket (150줄, 5개 테스트)

```
T1: socket_create() - 소켓 생성
T2: socket_bind/listen() - 포트 바인딩
T3: socket_accept() - 수락 (클라이언트 시뮬레이션)
T4: socket_read/write() - 데이터 송수신
T5: socket_close() - 정상 종료
```

### Test Suite 2: HTTP Parser (250줄, 12개 테스트)

```
T1-T3: Request Line 파싱
  - T1: "GET / HTTP/1.1"
  - T2: "POST /todos HTTP/1.1"
  - T3: "DELETE /todos/123 HTTP/1.1"

T4-T7: Headers 파싱
  - T4: Content-Type
  - T5: Content-Length
  - T6: Host
  - T7: Authorization

T8-T10: Body 처리
  - T8: JSON 본문
  - T9: Form 데이터
  - T10: 빈 본문

T11-T12: 에러 처리
  - T11: 잘못된 Request Line
  - T12: 헤더 형식 오류
```

### Test Suite 3: Server (150줄, 5개 테스트)

```
T1: 단일 요청 (GET /todos)
T2: 다중 요청 (Keep-Alive)
T3: 동시 연결 (N=10)
T4: 대용량 요청 (1MB)
T5: 타임아웃
```

### Benchmark (100줄, 성능 측정)

```
B1: Throughput (req/sec)
B2: Latency (p50/p95/p99)
B3: 메모리 사용 (RSS)
B4: CPU 사용률
```

---

## 📊 구현 체크리스트

### TCP Socket Manager (250줄)

- [ ] `socket_create()` - 소켓 생성
- [ ] `socket_bind()` - bind()
- [ ] `socket_listen()` - listen()
- [ ] `socket_accept()` - accept()
- [ ] `socket_read()` - read() with timeout
- [ ] `socket_write()` - write()
- [ ] `socket_close()` - close()
- [ ] 에러 처리 (errno 매핑)
- [ ] Buffer 관리
- [ ] Timeout 처리

### HTTP Parser (350줄)

- [ ] Request line 파싱
- [ ] Headers 파싱 (map)
- [ ] Body 처리 (Content-Length)
- [ ] Chunked 인코딩 (선택)
- [ ] 유효성 검증
- [ ] 에러 응답 생성
- [ ] Request 구조체
- [ ] Parser 상태 기계

### HTTP Handler (200줄)

- [ ] 상태 코드 → 이유 구문
- [ ] 응답 헤더 생성
- [ ] 응답 직렬화
- [ ] Keep-Alive 헤더
- [ ] CRLF 정확성

### HTTP Server (200줄)

- [ ] Socket 생성/바인드
- [ ] Accept loop
- [ ] 요청 읽기
- [ ] 라우팅 (Layer 5 연동)
- [ ] 응답 전송
- [ ] Keep-Alive 관리
- [ ] 타임아웃 처리
- [ ] Graceful shutdown

### 모듈 통합 (50줄)

- [ ] mod.fl 공개 API
- [ ] 타입 정의 내보내기
- [ ] 함수 내보내기

---

## 📈 성능 목표

| 지표 | 목표 | 측정 방법 |
|------|------|----------|
| Throughput | > 10K req/sec | ab or wrk |
| Latency (p99) | < 50ms | latency histogram |
| Keep-Alive overhead | < 1ms | 연속 요청 비교 |
| Parser overhead | < 100µs | 평균 파싱 시간 |
| Memory (1K conn) | < 10MB | RSS |
| CPU (idle) | < 5% | top |

---

## 🔗 Integration Points

### Input
- 기존 REST API (freelang-rest-api)
- Layer 5 Router

### Output
- HTTP 응답 (클라이언트)
- Layer 2 (Raft DB) 연결

### Dependencies
- Kernel syscalls (socket, bind, listen, accept, read, write, close)
- FreeLang stdlib (기본만)

---

## 📝 구현 노트

**왜 프레임워크 안 쓰나?**
- 저수준 제어 필요 (µs 단위)
- Zero-copy 최적화 필요
- 커널 제어권 필수

**왜 HTTP/1.1만?**
- HTTP/2는 복잡도 ↑ (200줄 → 800줄)
- 대부분 사용 사례에 충분
- 성능은 충분 (10K req/sec)

**멀티플렉싱?**
- Keep-Alive (TCP 재사용)
- Non-blocking I/O (선택사항)
- 다중 동시 연결 (128 backlog)

---

## 🎯 완성 조건

✅ 모든 체크리스트 항목 완성
✅ 모든 테스트 통과 (20개)
✅ 성능 목표 달성 (5/5)
✅ 무관용 규칙 준수 (R1-R5)
✅ GOGS에 커밋

---

**다음**: TCP Socket Manager 구현 시작 (src/tcp_socket.fl)
