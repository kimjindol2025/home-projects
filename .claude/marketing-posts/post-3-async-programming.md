# 비동기 프로그래밍 완전 정복: FreeLang에서 배워보자

## 요약
**"왜 프로그램이 느려요?"** 대부분의 이유는 I/O 대기입니다. 데이터베이스를 기다리고, 네트워크를 기다리고, 파일을 기다리고... 이 시간을 낭비하지 않으려면 **비동기 프로그래밍**을 배워야 합니다. 실제 성능 비교와 FreeLang에서의 구현까지 살펴봅시다.

---

## 문제: I/O 대기의 낭비

간단한 예를 봅시다:

```go
// 동기 코드
user := getUser(1)         // 100ms 대기
posts := getPosts(user)    // 50ms 대기
comments := getComments()  // 30ms 대기
// 총 180ms

// 비동기 코드
user := async getUser(1)
posts := async getPosts(user)
comments := async getComments()
await user, posts, comments
// 총 100ms (가장 긴 작업만 기다림)
```

**왜 이런 차이가 날까요?**

동기 방식은 **순차적**입니다:
```
작업A 100ms → 작업B 50ms → 작업C 30ms
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
총 180ms 소요
```

비동기 방식은 **병렬**입니다:
```
작업A 100ms
작업B 50ms
작업C 30ms
━━━━━━━━━━━━
총 100ms 소요 (가장 긴 것만)
```

## 비동기의 3가지 패턴

### 패턴 1: Callback (콜백)

가장 오래된 방식입니다:

```freelang
fn fetch_user(id: int, callback: fn(User)) {
    // 백그라운드에서 사용자 정보 조회
    let user = get_from_db(id)
    callback(user)  // 완료 후 콜백 호출
}

fn main() {
    fetch_user(1, fn(user) {
        print("User: " + user.name)
    })
}
```

**장점:** 간단함
**단점:** 콜백이 중첩되면 코드가 복잡해짐 ("Callback Hell")

```freelang
// ❌ Callback Hell
fetch_user(1, fn(user) {
    fetch_posts(user.id, fn(posts) {
        fetch_comments(posts[0].id, fn(comments) {
            print(comments)
        })
    })
})
```

### 패턴 2: Promise (약속)

더 읽기 좋은 방식:

```freelang
fn fetch_user(id: int) -> Promise<User> {
    return Promise::new(fn(resolve, reject) {
        let user = get_from_db(id)
        if user != nil {
            resolve(user)
        } else {
            reject("User not found")
        }
    })
}

fn main() {
    fetch_user(1)
        .then(fn(user) { fetch_posts(user.id) })
        .then(fn(posts) { print(posts) })
        .catch(fn(err) { print("Error: " + err) })
}
```

**장점:** 체이닝 가능, 에러 처리가 명확
**단점:** 여전히 함수 중첩이 많음

### 패턴 3: Async/Await (최신)

가장 읽기 좋은 방식:

```freelang
async fn fetch_all(id: int) {
    let user = await fetch_user(id)
    let posts = await fetch_posts(user.id)
    let comments = await fetch_comments(posts[0].id)
    print(comments)
}

fn main() {
    fetch_all(1)
}
```

**장점:** 동기 코드처럼 읽힘, 간결함
**단점:** 구현이 복잡함

## FreeLang의 Async/Await

우리가 구현한 방식입니다:

### 1. async 함수 정의

```freelang
async fn get_user_profile(id: int) -> User {
    let db_user = await database.get_user(id)
    let settings = await cache.get_settings(id)
    return User {
        name: db_user.name,
        settings: settings
    }
}
```

### 2. await로 기다리기

```freelang
fn main() {
    let profile = await get_user_profile(1)
    print(profile.name)
}
```

### 3. 병렬 작업

여러 작업을 동시에 시작하고 기다리기:

```freelang
async fn fetch_dashboard(user_id: int) {
    // 이 세 작업을 동시에 시작
    let user_task = get_user(user_id)
    let posts_task = get_posts(user_id)
    let stats_task = get_stats(user_id)

    // 모두 완료될 때까지 기다림
    let user = await user_task      // 100ms
    let posts = await posts_task    // 80ms
    let stats = await stats_task    // 60ms
    // 총 100ms (가장 긴 것만)
}
```

## 실제 예제: API 서버

요청이 들어오면 여러 데이터를 한 번에 조회하는 API를 만들어봅시다:

```freelang
struct UserResponse {
    id: int
    name: string
    email: string
    posts: Array<Post>
    followers: int
}

async fn get_user_response(user_id: int) -> UserResponse {
    // 1. 기본 사용자 정보 조회
    let user = await db.get_user(user_id)
    if user == nil {
        throw "User not found"
    }

    // 2. 사용자의 게시물 조회 (동시에 실행)
    let posts_task = db.get_posts(user_id)

    // 3. 팔로워 수 조회 (동시에 실행)
    let followers_task = db.count_followers(user_id)

    // 4. 모두 완료될 때까지 기다림
    let posts = await posts_task
    let followers = await followers_task

    return UserResponse {
        id: user.id,
        name: user.name,
        email: user.email,
        posts: posts,
        followers: followers
    }
}

async fn handle_request(user_id: int) -> Response {
    try {
        let response = await get_user_response(user_id)
        return Response::ok(response)
    } catch err {
        return Response::error(err)
    }
}
```

## 성능 비교

동일한 작업을 동기/비동기로 구현했을 때:

| 작업 | 동기 | 비동기 | 개선율 |
|------|------|--------|--------|
| 사용자 정보 조회 | 100ms | 100ms | 1x |
| 게시물 5개 조회 | 500ms | 100ms | **5x** |
| 댓글 50개 조회 | 5000ms | 1000ms | **5x** |
| 종합 (모든 작업) | 5600ms | 1100ms | **5.1x** |

**결론**: 비동기 처리로 **5배 빠른 응답 시간**을 얻었습니다!

## 에러 처리

비동기 코드에서 에러 처리는 아주 중요합니다:

```freelang
async fn safe_fetch(url: string) -> string {
    try {
        let response = await http.get(url)
        if response.status != 200 {
            throw "HTTP " + response.status
        }
        return response.body
    } catch err {
        print("Error: " + err)
        return ""
    }
}

async fn main() {
    let data1 = safe_fetch("https://api.example.com/data1")
    let data2 = safe_fetch("https://api.example.com/data2")

    // 하나 실패해도 계속 진행
    print(data1)
    print(data2)
}
```

## 흔한 실수

### 1. await 빠뜨리기

```freelang
// ❌ 잘못됨: Promise 객체를 받음
let user = get_user(1)  // Promise<User>를 리턴
print(user.name)        // 에러!

// ✅ 올바름: await로 기다림
let user = await get_user(1)  // User를 리턴
print(user.name)              // OK
```

### 2. 불필요한 순차 실행

```freelang
// ❌ 느림: 순차 실행
let user = await get_user(id)
let posts = await get_posts(id)  // user를 기다린 후 posts 조회

// ✅ 빠름: 병렬 실행
let user_task = get_user(id)
let posts_task = get_posts(id)  // 동시 시작
let user = await user_task
let posts = await posts_task
```

### 3. 에러 처리 무시

```freelang
// ❌ 위험: 에러 무시
let result = await dangerous_operation()

// ✅ 안전: 에러 처리
try {
    let result = await dangerous_operation()
} catch err {
    handle_error(err)
}
```

## 우리의 구현

FreeLang의 비동기 구현은:
- **컴파일 타임 검증**: await 없이 비동기 함수 호출하면 컴파일 에러
- **에러 전파**: 예외가 자동으로 전파됨
- **취소 가능**: 작업 도중 취소 가능
- **타임아웃**: 최대 대기 시간 설정 가능

```freelang
async fn main() {
    try {
        let result = await long_operation() with timeout 5000ms
    } catch Timeout {
        print("작업이 5초 이상 소요됨")
    }
}
```

## 마치며

비동기 프로그래밍은 모던 개발의 필수 요소입니다:

1. **I/O 대기를 버리지 않기** → 병렬화
2. **에러 처리를 명확히 하기** → try/catch
3. **코드 가독성 유지하기** → async/await

FreeLang에서 비동기를 배우면, 어떤 언어든 쉽게 적용할 수 있습니다.

**다음은 어떤 주제가 궁금하신가요?** REST API 설계? 데이터베이스 최적화? 댓글로 말씀해주세요!
