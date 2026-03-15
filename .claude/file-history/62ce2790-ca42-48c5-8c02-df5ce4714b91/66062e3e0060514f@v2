# 🚀 FreeLang WebSocket 실시간 API 설정 가이드

**버전**: Phase 4-2 (Real-time Features)
**상태**: ✅ 구현 완료
**성능**: <50ms 메시지 지연 (로컬 네트워크)

---

## 목차

1. [WebSocket 개요](#-websocket-개요)
2. [설정 및 시작](#-설정-및-시작)
3. [클라이언트 연결](#-클라이언트-연결)
4. [실시간 기능](#-실시간-기능)
5. [채널 & 구독](#-채널--구독)
6. [성능 최적화](#-성능-최적화)

---

## 📡 WebSocket 개요

FreeLang WebSocket API는 RFC 6455 표준을 준수하는 실시간 양방향 통신을 제공합니다.

### 주요 특징

| 특징 | 설명 |
|------|------|
| **낮은 지연** | <50ms 평균 메시지 지연 |
| **채널 기반** | 여러 채널 동시 구독 |
| **자동 하트비트** | 연결 유지 위한 PING/PONG |
| **재연결 지원** | 자동 재연결 메커니즘 |
| **인증** | JWT 토큰 기반 보안 |

---

## 🚀 설정 및 시작

### 1단계: WebSocket 엔드포인트 활성화

```bash
# docker-compose.yml에서 WebSocket 포트 활성화
docker-compose up -d

# WebSocket 엔드포인트
# ws://localhost:5021/ws (개발)
# wss://253.dclub.kr/ws (프로덕션)
```

### 2단계: 핸드셰이크 및 인증

```bash
# curl을 사용한 WebSocket 핸드셰이크 테스트
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:5021/ws
```

### 3단계: 메시지 형식 확인

```
← 연결 성공 응답
→ 첫 번째 메시지 전송
← 서버 응답
```

---

## 💻 클라이언트 연결

### JavaScript (Browser)

```javascript
// WebSocket 연결
const token = localStorage.getItem('auth_token');
const ws = new WebSocket(
  'ws://localhost:5021/ws',
  token ? ['Authorization', token] : []
);

// 연결 성공
ws.onopen = (event) => {
  console.log('WebSocket 연결 성공');

  // 채널 구독
  ws.send(JSON.stringify({
    type: 'subscribe',
    channel: 'post:123'
  }));
};

// 메시지 수신
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('실시간 메시지:', message);

  switch (message.event_type) {
    case 'comment_added':
      handleNewComment(message.data);
      break;
    case 'post_updated':
      handlePostUpdate(message.data);
      break;
    case 'user_typing':
      handleUserTyping(message.data);
      break;
  }
};

// 연결 종료
ws.onclose = (event) => {
  console.log('WebSocket 연결 종료');
  // 재연결 로직
};

// 에러 처리
ws.onerror = (error) => {
  console.error('WebSocket 에러:', error);
};
```

### Python (AsyncIO)

```python
import asyncio
import websockets
import json

async def connect_websocket():
    token = "YOUR_JWT_TOKEN"
    uri = "ws://localhost:5021/ws"

    async with websockets.connect(
        uri,
        subprotocols=['Authorization', token]
    ) as websocket:
        # 채널 구독
        await websocket.send(json.dumps({
            'type': 'subscribe',
            'channel': 'post:123'
        }))

        # 메시지 수신
        async for message in websocket:
            data = json.loads(message)
            print(f"실시간 메시지: {data}")

# 실행
asyncio.run(connect_websocket())
```

### Go (Gorilla WebSocket)

```go
import (
    "github.com/gorilla/websocket"
)

func connectWebSocket() {
    token := "YOUR_JWT_TOKEN"

    dialer := websocket.Dialer{
        HandshakeTimeout: time.Second * 5,
    }

    header := http.Header{}
    header.Add("Authorization", "Bearer "+token)

    conn, _, err := dialer.Dial("ws://localhost:5021/ws", header)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 채널 구독
    msg := map[string]string{
        "type": "subscribe",
        "channel": "post:123",
    }
    conn.WriteJSON(msg)

    // 메시지 수신
    for {
        var data map[string]interface{}
        err := conn.ReadJSON(&data)
        if err != nil {
            log.Println(err)
            break
        }
        fmt.Printf("실시간 메시지: %v\n", data)
    }
}
```

---

## 📊 실시간 기능

### 1. 댓글 실시간 알림

#### 클라이언트 구독

```javascript
// 특정 포스트의 댓글 구독
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'post:123'
}));
```

#### 서버 이벤트

```
← {
  "event_type": "comment_added",
  "channel": "post:123",
  "data": {
    "comment_id": "c456",
    "post_id": "123",
    "content": "좋은 글입니다!",
    "author": "user1",
    "author_id": "u789",
    "created_at": "2026-03-13T12:00:00Z"
  },
  "timestamp": 1234567890,
  "sender_id": "u789"
}
```

#### 클라이언트 처리

```javascript
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.event_type === 'comment_added') {
    const comment = msg.data;

    // UI 업데이트
    const commentEl = document.createElement('div');
    commentEl.className = 'comment';
    commentEl.innerHTML = `
      <strong>${comment.author}</strong>
      <p>${comment.content}</p>
      <small>${comment.created_at}</small>
    `;

    document.getElementById('comments').appendChild(commentEl);

    // 알림 표시
    showNotification(`${comment.author}님이 댓글을 남겼습니다`);
  }
};
```

### 2. 포스트 좋아요 실시간 업데이트

```javascript
// 포스트 좋아요 버튼
function toggleLike(postId) {
  fetch(`/api/posts/${postId}/like`, { method: 'POST' })
    .then(() => {
      // 서버가 WebSocket으로 모든 구독자에게 전송
      // ← {
      //   "event_type": "post_like",
      //   "channel": "post:123",
      //   "data": {
      //     "post_id": "123",
      //     "user_id": "u456",
      //     "action": "liked"
      //   }
      // }
    });
}

// 좋아요 수 실시간 업데이트
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.event_type === 'post_like') {
    const likeCount = document.getElementById('like-count');
    const newCount = parseInt(likeCount.textContent) +
                     (msg.data.action === 'liked' ? 1 : -1);
    likeCount.textContent = newCount;
  }
};
```

### 3. 사용자 타이핑 표시기

#### 타이핑 감지

```javascript
let typingTimeout;

function onCommentInputChange(text) {
  if (!typingTimeout) {
    // 타이핑 시작
    ws.send(JSON.stringify({
      type: 'typing',
      post_id: '123',
      user_id: currentUser.id
    }));
  }

  // 타이머 리셋 (2초 후 타이핑 중지 신호)
  clearTimeout(typingTimeout);
  typingTimeout = setTimeout(() => {
    ws.send(JSON.stringify({
      type: 'stopped_typing',
      post_id: '123',
      user_id: currentUser.id
    }));
    typingTimeout = null;
  }, 2000);
}

// 입력창 변화 감지
document.getElementById('comment-input')
  .addEventListener('input', (e) => onCommentInputChange(e.target.value));
```

#### 타이핑 표시 UI

```javascript
const typingUsers = {};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.event_type === 'user_typing') {
    const { user_id, username } = msg.data;
    typingUsers[user_id] = username;
    updateTypingIndicator();
  }

  if (msg.event_type === 'user_stopped_typing') {
    delete typingUsers[msg.data.user_id];
    updateTypingIndicator();
  }
};

function updateTypingIndicator() {
  const names = Object.values(typingUsers).join(', ');
  const indicator = document.getElementById('typing-indicator');

  if (names) {
    indicator.textContent = `${names}님이 댓글을 작성 중입니다...`;
    indicator.style.display = 'block';
  } else {
    indicator.style.display = 'none';
  }
}
```

### 4. 활동 피드 (Activity Feed)

```javascript
// 팔로우한 사용자의 활동 실시간 수신
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'user:current_user_id'
}));

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.event_type === 'activity_feed_update') {
    const activity = msg.data;

    let feedItem = '';
    switch (activity.activity_type) {
      case 'post_published':
        feedItem = `${activity.actor_name}이 새 글을 작성했습니다`;
        break;
      case 'comment_added':
        feedItem = `${activity.actor_name}이 댓글을 남겼습니다`;
        break;
      case 'follow':
        feedItem = `${activity.actor_name}이 당신을 팔로우했습니다`;
        break;
    }

    const feedEl = document.createElement('li');
    feedEl.textContent = feedItem;
    document.getElementById('activity-feed').prepend(feedEl);
  }
};
```

### 5. 사용자 온라인 상태

```javascript
// 사용자 온라인 상태 추적
const onlineUsers = new Set();

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.event_type === 'user_online') {
    onlineUsers.add(msg.data.user_id);
    updateOnlineIndicators();
  }

  if (msg.event_type === 'user_offline') {
    onlineUsers.delete(msg.data.user_id);
    updateOnlineIndicators();
  }
};

function updateOnlineIndicators() {
  document.querySelectorAll('[data-user-id]').forEach(el => {
    const userId = el.dataset.userId;
    if (onlineUsers.has(userId)) {
      el.classList.add('online');
    } else {
      el.classList.remove('online');
    }
  });
}
```

---

## 📡 채널 & 구독

### 채널 구조

| 채널 | 용도 | 예시 |
|------|------|------|
| `global` | 전역 이벤트 | 새 포스트 발행 |
| `post:ID` | 특정 포스트 | `post:123` |
| `user:ID` | 특정 사용자 | `user:456` |
| `doc:ID` | 협업 문서 | `doc:789` |
| `direct:CONN_ID` | 직접 메시지 | (개인용) |

### 구독 관리

```javascript
// 구독
function subscribe(channel) {
  ws.send(JSON.stringify({
    type: 'subscribe',
    channel: channel
  }));
}

// 구독 해제
function unsubscribe(channel) {
  ws.send(JSON.stringify({
    type: 'unsubscribe',
    channel: channel
  }));
}

// 사용 예
subscribe('post:123');
subscribe('user:456');

// 포스트 변경 시 구독 변경
function switchPost(newPostId) {
  unsubscribe('post:123');
  subscribe(`post:${newPostId}`);
}
```

---

## ⚡ 성능 최적화

### 1. 메시지 배치

```javascript
// 여러 메시지를 한 번에 처리
const messageBuffer = [];
let flushTimer;

function queueMessage(message) {
  messageBuffer.push(message);

  if (messageBuffer.length >= 10 || !flushTimer) {
    flushMessages();
  }

  // 100ms 후 자동 flush
  if (!flushTimer) {
    flushTimer = setTimeout(flushMessages, 100);
  }
}

function flushMessages() {
  messageBuffer.forEach(msg => processMessage(msg));
  messageBuffer = [];
  flushTimer = null;
}
```

### 2. 자동 재연결

```javascript
class ReconnectingWebSocket {
  constructor(url, options = {}) {
    this.url = url;
    this.retryCount = 0;
    this.maxRetries = options.maxRetries || 5;
    this.retryDelay = options.retryDelay || 1000;

    this.connect();
  }

  connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      this.retryCount = 0;
      this.onopen?.();
    };

    this.ws.onerror = () => {
      this.reconnect();
    };

    this.ws.onmessage = this.onmessage;
    this.ws.onclose = () => {
      if (this.retryCount < this.maxRetries) {
        this.reconnect();
      }
    };
  }

  reconnect() {
    this.retryCount++;
    const delay = this.retryDelay * Math.pow(2, this.retryCount - 1);
    console.log(`${delay}ms 후 재연결 시도...`);
    setTimeout(() => this.connect(), delay);
  }

  send(message) {
    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message);
    }
  }
}

// 사용
const ws = new ReconnectingWebSocket('ws://localhost:5021/ws');
```

### 3. 메모리 관리

```javascript
// 오래된 메시지 자동 제거
const MAX_MESSAGES = 1000;
const messages = [];

function addMessage(message) {
  messages.push(message);

  if (messages.length > MAX_MESSAGES) {
    messages.shift();
  }
}

// 주기적 리소스 정리
setInterval(() => {
  // 정크 수집
  if (messages.length > 500) {
    messages.splice(0, 250);
  }
}, 60000);
```

---

## 🔐 보안

### 인증 토큰 전달

```javascript
const token = localStorage.getItem('auth_token');

const ws = new WebSocket('ws://localhost:5021/ws');

// 연결 후 인증
ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'authenticate',
    token: token
  }));
};
```

### CORS & Same-Origin

```javascript
// 프로덕션에서는 TLS 필수
const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const ws = new WebSocket(
  `${wsProtocol}//${window.location.host}/ws`
);
```

---

## 📊 모니터링

### 연결 상태 추적

```javascript
class WebSocketMonitor {
  constructor(ws) {
    this.ws = ws;
    this.stats = {
      connected: false,
      messagesReceived: 0,
      messagesSent: 0,
      errors: 0,
      lastMessageTime: null
    };

    this.setupMonitoring();
  }

  setupMonitoring() {
    this.ws.onopen = () => {
      this.stats.connected = true;
      console.log('📡 WebSocket 연결됨');
    };

    this.ws.onmessage = () => {
      this.stats.messagesReceived++;
      this.stats.lastMessageTime = Date.now();
    };

    this.ws.onerror = () => {
      this.stats.errors++;
    };

    this.ws.onclose = () => {
      this.stats.connected = false;
      console.log('📡 WebSocket 연결 해제');
    };
  }

  getStats() {
    return {
      ...this.stats,
      uptime: Date.now() - this.startTime,
      messageRate: this.stats.messagesReceived / (Date.now() - this.startTime)
    };
  }
}
```

---

## 🚀 다음 단계

### Phase 4-3: GraphQL Subscriptions

```graphql
subscription OnCommentAdded {
  commentAdded(postId: "123") {
    id
    content
    author
    createdAt
  }
}
```

### Phase 4-4: 메시지 암호화

```javascript
// 엔드-투-엔드 암호화 (E2EE)
import crypto from 'crypto';

function encryptMessage(message, publicKey) {
  const encrypted = crypto.publicEncrypt(
    publicKey,
    Buffer.from(message)
  );
  return encrypted.toString('base64');
}
```

---

**마지막 업데이트**: 2026-03-13
**WebSocket 버전**: RFC 6455
**FreeLang 버전**: 2.8.0
