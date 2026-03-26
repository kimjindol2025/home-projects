---
title: "FreeLang의 비동기 프로그래밍 - async/await 완전 정복"
date: 2026-03-26
author: Content Writer
category: Technical
tags:
  - FreeLang
  - Async Programming
  - async/await
  - Future
  - Concurrency
  - Non-blocking
---

# FreeLang의 비동기 프로그래밍
## async/await로 동시성 완전히 이해하기

**글을 읽으면 얻을 수 있는 것:**
1. 비동기 프로그래밍의 핵심 개념 이해하기
2. async/await 문법과 작동 원리 파악하기
3. 실제 프로젝트에서 비동기 코드 작성하기

---

## 배경: 왜 비동기가 필요한가?

일반적인 프로그램의 실행을 생각해봅시다.

```freelang
// 동기 코드 (한 줄씩 순서대로)
PRINTLN("작업 1 시작")     // 시간: 0초
do_work_1()               // 시간: 3초 (총 3초)
PRINTLN("작업 1 완료")

PRINTLN("작업 2 시작")     // 시간: 3초
do_work_2()               // 시간: 2초 (총 5초)
PRINTLN("작업 2 완료")

// 총 소요 시간: 5초 ⏱️
```

만약 작업 1, 2가 **독립적**이라면?

```
현재 (동기):
작업1 ▓▓▓ (3초) → 작업2 ▓▓ (2초) = 총 5초

이상적 (비동기):
작업1 ▓▓▓ (3초) \
               → 동시 실행 = 총 3초 ⏱️
작업2 ▓▓ (2초) /
```

**50% 빨라집니다!** 🚀

---

## 문제: 동기 코드의 한계

### 1. 대기 시간 낭비

```freelang
// 네트워크 요청 예제
FN get_user_data(user_id: INT) -> STRING {
  PRINTLN("데이터 요청 중...")
  let response = http_get("api.example.com/user/" + user_id)  // 2초 대기
  PRINTLN("데이터 받음!")
  RETURN response
}

// 메인 프로그램
PRINTLN("시작")                 // 0초
LET user1 = get_user_data(1)   // 0-2초 (대기)
LET user2 = get_user_data(2)   // 2-4초 (대기)
PRINTLN("완료")                // 4초

// 총 4초 걸림 😢
// 2초는 그냥 "대기"하는 시간
```

### 2. UI가 먹통

```freelang
// 버튼 클릭 → 느린 작업
FN on_button_click() {
  let result = slow_calculation()  // 5초 소요
  update_ui(result)
}

// 사용자 관점:
// ❌ 5초 동안 UI가 응답 없음
// ❌ 버튼 클릭 안 됨
// ❌ 나쁜 사용자 경험 👎
```

### 3. 확장성 문제

```freelang
// 1000명의 사용자 데이터 가져오기
FOR i IN 1..1000 {
  let user = get_user_data(i)  // 각 2초
  process_user(user)
}

// 총 소요 시간: 1000 * 2초 = 2000초 (33분!) 😱
```

---

## 해결책: FreeLang의 async/await

### 핵심 개념

**Future**: 미래에 완료될 작업을 나타내는 객체

```freelang
// 지금 실행되지 않고, 나중에 실행될 작업
ASYNC FN get_user_data(user_id: INT) -> STRING {
  // 이 함수는 즉시 실행되지 않음
  // 대신 Future<STRING>을 반환
  let response = await http_get("api.example.com/user/" + user_id)
  RETURN response
}
```

**await**: Future가 완료될 때까지 기다리기

```freelang
// async 함수의 결과를 기다림
let user1 = await get_user_data(1)  // 1번 사용자 데이터
```

### 기본 구조

```freelang
// async 함수 정의
ASYNC FN fetch_data() -> STRING {
  // 오래 걸리는 작업...
  let result = await slow_operation()
  RETURN result
}

// async 함수 호출
// 1. Future 반환 (즉시)
let future = fetch_data()

// 2. await로 기다림
let result = await future
```

---

## 실전: 비동기 코드 작성하기

### 초보자라면?

간단한 타이머 예제로 시작합시다.

```freelang
// timer.fl
ASYNC FN wait_seconds(seconds: INT) -> STRING {
  // seconds초 기다린 후 메시지 반환
  PRINTLN("⏳ " + seconds + "초 대기 중...")
  await sleep(seconds)  // 기본 제공 함수
  PRINTLN("✅ " + seconds + "초 경과!")
  RETURN "완료!"
}

FN main() {
  PRINTLN("프로그램 시작")

  // 동기 방식 (순차 실행)
  // TOTAL: 3초 + 2초 = 5초
  let result1 = await wait_seconds(3)
  let result2 = await wait_seconds(2)

  PRINTLN("모두 완료!")
}
```

**실행 결과:**
```
프로그램 시작
⏳ 3초 대기 중...
✅ 3초 경과!
⏳ 2초 대기 중...
✅ 2초 경과!
모두 완료!
(총 5초)
```

#### 비동기로 병렬 실행하기

```freelang
FN main() {
  PRINTLN("프로그램 시작")

  // 비동기 방식 (동시 실행)
  let future1 = wait_seconds(3)  // Future 생성 (await 없음!)
  let future2 = wait_seconds(2)  // Future 생성

  // 둘 다 완료될 때까지 대기
  let result1 = await future1
  let result2 = await future2

  PRINTLN("모두 완료!")
}
```

**실행 결과:**
```
프로그램 시작
⏳ 3초 대기 중...
⏳ 2초 대기 중...
✅ 2초 경과!   (2초 후)
✅ 3초 경과!   (3초 후)
모두 완료!
(총 3초) ⚡ 훨씬 빠름!
```

### 전문가라면?

복잡한 비동기 흐름을 다루봅시다.

**웹 크롤러 예제:**

```freelang
RECORD PageData {
  url: STRING
  title: STRING
  content: STRING
  status: INT
}

// HTTP 요청 (비동기)
ASYNC FN fetch_url(url: STRING) -> STRING {
  PRINTLN("🌐 " + url + " 요청 중...")
  let response = await http_get(url)
  PRINTLN("✅ " + url + " 수신 완료")
  RETURN response
}

// HTML 파싱 (비동기)
ASYNC FN parse_html(html: STRING) -> RECORD {
  PRINTLN("📄 HTML 파싱 중...")
  await sleep(1)  // 파싱 시뮬레이션
  RETURN {
    url: "",
    title: "Page Title",
    content: html,
    status: 200
  }
}

// 여러 페이지 동시에 크롤링
ASYNC FN crawl_pages(urls: ARRAY<STRING>) -> ARRAY<PageData> {
  PRINTLN("🕷️ 크롤링 시작 (" + len(urls) + "개 페이지)")

  // Step 1: 모든 URL을 동시에 요청
  let futures = []
  FOR url IN urls {
    let future = fetch_url(url)
    futures = append(futures, future)
  }

  // Step 2: 모든 응답을 기다림
  let responses = []
  FOR future IN futures {
    let html = await future
    responses = append(responses, html)
  }

  // Step 3: 각 HTML 파싱
  let parse_futures = []
  FOR html IN responses {
    let future = parse_html(html)
    parse_futures = append(parse_futures, future)
  }

  // Step 4: 모든 파싱 완료 대기
  let results = []
  FOR future IN parse_futures {
    let data = await future
    results = append(results, data)
  }

  PRINTLN("✅ 크롤링 완료!")
  RETURN results
}

// 메인 프로그램
ASYNC FN main() {
  let urls = [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3",
    "https://example.com/page4"
  ]

  let pages = await crawl_pages(urls)

  PRINTLN("수집된 페이지:")
  FOR page IN pages {
    PRINTLN("- " + page.title)
  }
}

await main()
```

**실행 예상 시간:**
```
동기: 4개 페이지 × 2초(요청) + 1초(파싱) = 12초
비동기: max(2초 요청 + 1초 파싱) = 3초

⚡ 4배 빨라짐!
```

---

## async/await의 고급 패턴

### 1. Error Handling (에러 처리)

```freelang
ASYNC FN fetch_with_retry(url: STRING, max_retries: INT) -> STRING {
  LET attempt = 0

  LOOP {
    IF attempt >= max_retries {
      RETURN "실패"
    }

    TRY {
      let response = await http_get(url)
      RETURN response
    } CATCH error {
      attempt = attempt + 1
      PRINTLN("⚠️ 재시도 " + attempt)
      await sleep(1)  // 1초 대기 후 재시도
    }
  }
}
```

### 2. Timeout (시간 제한)

```freelang
ASYNC FN fetch_with_timeout(url: STRING, timeout_seconds: INT) -> STRING {
  let fetch_future = http_get(url)
  let timeout_future = sleep(timeout_seconds)

  // 둘 중 먼저 완료되는 쪽
  TRY {
    let response = await race(fetch_future, timeout_future)
    RETURN response
  } CATCH {
    RETURN "타임아웃!"
  }
}
```

### 3. 병렬 처리 (Promise.all 패턴)

```freelang
ASYNC FN process_all(items: ARRAY<STRING>) -> ARRAY<INT> {
  // 모든 작업의 Future 생성
  let futures = []
  FOR item IN items {
    futures = append(futures, process_item(item))
  }

  // 모두 완료될 때까지 대기
  let results = []
  FOR future IN futures {
    results = append(results, await future)
  }

  RETURN results
}
```

### 4. Race (경쟁)

```freelang
ASYNC FN first_to_respond(urls: ARRAY<STRING>) -> STRING {
  let futures = []
  FOR url IN urls {
    futures = append(futures, fetch_url(url))
  }

  // 가장 먼저 응답하는 서버의 데이터 반환
  let winner = await race(futures)
  RETURN winner
}
```

---

## 비동기 vs 멀티스레딩

### 비동기 (async/await)
```
한 스레드에서 여러 작업을 번갈아가며 처리
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

시간: 0초    1초    2초    3초    4초
      │      │      │      │      │
Task1 ▓▓▓━━━━│      │      │      │
Task2 ┃    ▓▓▓━━━━│      │      │
Task3 ┃    ┃    ▓▓▓━━━━│      │
Task4 ┃    ┃    ┃    ▓▓▓━━━━│
```

**장점:**
- 가벼움 (스레드 생성 비용 없음)
- 컨텍스트 스위칭 적음
- 메모리 효율적

**단점:**
- CPU 집약적 작업에 부적합
- 하나의 작업이 오래 걸리면 다른 작업 블로킹

### 멀티스레딩

```
여러 스레드가 동시에 작업
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

스레드1 ▓▓▓▓▓▓▓▓ (병렬 처리)
스레드2 ▓▓▓▓▓▓▓▓
스레드3 ▓▓▓▓▓▓▓▓
```

**장점:**
- 진정한 병렬 처리
- CPU 집약적 작업에 적합

**단점:**
- 무거움 (스레드 생성 비용)
- 동기화 문제 (데이터 경쟁)

### 선택 기준

```
네트워크 I/O (API 호출) → async/await ⭐
파일 읽기/쓰기 → async/await ⭐
데이터베이스 쿼리 → async/await ⭐
CPU 집약적 계산 → 멀티스레딩 ⭐
```

---

## 성능 비교

**1000개의 API 호출 예제:**

```
동기:
for i in 1..1000:
  result = api_call(i)  // 각 100ms
총: 1000 × 100ms = 100초 😱

비동기 (10개 동시):
for i in 1..100:
  wait for 10 api_calls in parallel
총: 100 × 100ms = 10초 ⚡

10배 빨라짐!
```

---

## Best Practices

### ✅ 하는 것

```freelang
// ✅ Good: 병렬 처리
let f1 = operation1()
let f2 = operation2()
let r1 = await f1
let r2 = await f2

// ✅ Good: 에러 처리
TRY {
  let result = await risky_operation()
} CATCH {
  handle_error()
}

// ✅ Good: 타임아웃
let result = await with_timeout(operation(), 5000)
```

### ❌ 하지 말 것

```freelang
// ❌ Bad: 불필요한 순차 await
let r1 = await op1()
let r2 = await op2()  // op1 완료를 기다림

// ✅ 개선:
let f1 = op1()
let f2 = op2()
let r1 = await f1
let r2 = await f2

// ❌ Bad: 에러 무시
let result = await maybe_error()  // 어쩌면 실패

// ✓ 개선:
TRY {
  let result = await maybe_error()
} CATCH e {
  PRINTLN("에러: " + e)
}
```

---

## 다음 단계

비동기 프로그래밍을 마스터했으니, 더 고급 주제를 배울 수 있습니다.

**다음에 읽을 거리:**
1. [Reactive Programming - Streams와 Observables](#)
2. [async/await 디버깅 - 어디서 막히는가](#)
3. [성능 프로파일링 - 비동기 코드 최적화](#)

**직접 해보고 싶으신가요?**

```bash
# FreeLang async 예제 실행
git clone https://github.com/freelang/freelang
cd freelang/examples/async
freelang crawler.fl
```

---

**참고 자료:**
- [Async Patterns - Rust async book](https://rust-lang.github.io/async-book/)
- [JavaScript Promises - MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Using_promises)
- FreeLang 공식 문서: Async Programming https://freelang-docs.example.com/async

---

이 글이 도움이 되었다면? 👍

- 다른 개발자와 공유해주세요
- "이런 비동기 패턴을 배우고 싶어요" 댓글로 제안해주세요
- 자신의 비동기 코드 최적화 경험담을 공유해주세요!

**Happy Async Coding! 🚀**
