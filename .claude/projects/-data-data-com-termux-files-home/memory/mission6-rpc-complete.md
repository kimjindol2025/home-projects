---
name: Mission 6 - Custom RPC Framework Complete
description: Mission 6 완성 (1,300줄, 18/18 테스트 PASS) - Binary wire protocol, zero-copy, Client/Server, KV over RPC
type: project
---

# ✅ **Mission 6: Custom RPC Framework - 100% 완성**

**완료일**: 2026-03-26
**규모**: ~1,300줄 (코드 850 + 테스트 450)
**테스트**: 18/18 PASS ✅

---

## 구현 완료 항목

### ✅ 파일 구조 (9개)
1. **errors.go** (12줄) - 6개 sentinel errors (Timeout, MethodNotFound, etc)
2. **protocol.go** (40줄) - 20-byte 고정 헤더 포맷, MessageType 상수
3. **message.go** (90줄) - Message Encode/Decode, zero-copy 원칙
4. **codec.go** (195줄) - Codec interface + BinaryCodec (6 type tags: string, bytes, bool, int64, nil, error)
5. **transport.go** (45줄) - Transport interface + net.Pipe() InProcess
6. **server.go** (115줄) - Server, Register, Serve, ServeAsync, error logging
7. **client.go** (180줄) - Client, Call, CallWithTimeout, background readLoop, ID matching
8. **kvstore_handler.go** (115줄) - KV.Set, KV.Get, KV.Delete, KV.Stats over RPC
9. **rpc_test.go** (450줄) - 18개 테이블 기반 테스트

---

## 핵심 기술

### 1️⃣ Binary Wire Protocol (20-byte 고정 헤더)
```
[Version:1][MsgType:1][Flags:2][ID:8][MethodLen:2][ErrorLen:2][PayloadLen:4]
                                                    ↓
                        [Method|Error|Body] (zero-copy sub-slices)
```

**MessageType**: MsgRequest(0x01), MsgResponse(0x02), MsgError(0x03), MsgPing(0x04), MsgPong(0x05)

### 2️⃣ Zero-Copy Deserialization
```
- headerBuf 1회 할당 (루프 밖) 재사용
- payloadBuf: 요청당 1회 할당
- Method/Error/Body: payloadBuf의 sub-slice (복사 X)
```

### 3️⃣ BinaryCodec (6 type tags)
```
0x01 string:  [type][len:4][UTF-8]
0x02 []byte:  [type][len:4][bytes]
0x03 bool:    [type][0x00/0x01]
0x04 int64:   [type][uint64:8 BE]
0x05 nil:     [type]
0x06 error:   [type][len:4][msg]
```

### 4️⃣ Client-Server Interaction
```
동시성 모델:
- Client.Call: nextID++, inflight[id]=pending, writeMu.Lock, Write
- readLoop: 백그라운드 고루틴, 응답 match by ID, pending.ch <- msg
- Timeout: time.After() with select on pending.ch (buffered 1)
```

### 5️⃣ KVStore Integration (4개 메서드)
```
KV.Set(key, value) → Cluster.Set()
KV.Get(key) → Cluster.Get()
KV.Delete(key) → Cluster.Delete()
KV.Stats() → Cluster.Stats()
```

---

## 테스트 커버리지 (18개)

| 그룹 | 테스트 | 항목 |
|------|--------|------|
| Message(4) | Encode/Decode roundtrip, Request layout, InvalidVersion, PayloadTooLarge | ✅ |
| Codec(5) | String, Bytes, Bool, Int64, Nil roundtrip | ✅ |
| Server(3) | DuplicateRegister, KnownMethod, UnknownMethod | ✅ |
| Client(5) | Basic call, Concurrent(10 goroutines), Timeout, ConnectionLost, PingPong | ✅ |
| Integration(1) | KVStore Stats | ✅ |

---

## 설계 결정

### 20-byte 고정 헤더 vs varint
- ✅ **고정**: 수신자가 항상 HeaderSize를 먼저 읽음, 루프 구현 단순
- ❌ varint: 2-3바이트 절약, 루프 복잡

### Synchronous Dispatch vs per-request goroutine
- ✅ **동기**: 호출 스택 명확, 메모리 효율, 테스트 단순
- ❌ per-goroutine: 확장성 높지만 복잡도 증가

### net.Pipe() for InProcess
- ✅ **net.Pipe()**: 실제 net.Conn과 동일한 코드 경로, 테스트-프로덕션 일관성
- ❌ 인메모리 채널: 불필요한 커스텀 로직

### BinaryCodec (hand-rolled) vs encoding/gob
- ✅ **hand-rolled**: 교육적, transparent, 타입 태그 명확
- ❌ gob: 복잡, 디버깅 어려움

---

## 모든 테스트 통과

```bash
go test ./pkg/rpc/... -v
# 18/18 PASS ✅

go test ./... -v
# 모든 프로젝트 테스트 PASS ✅
```

---

## 다음: Mission 7

**Mission 7: Security Gateway** 예정
- Formal verification을 통한 보안 검증
- 예상 규모: ~1,500줄

---

## 커밋 정보
```
pkg/rpc/
├── errors.go
├── protocol.go
├── message.go
├── codec.go
├── transport.go
├── server.go
├── client.go
├── kvstore_handler.go
└── rpc_test.go

전체: ~1,300줄, 18/18 테스트
```
