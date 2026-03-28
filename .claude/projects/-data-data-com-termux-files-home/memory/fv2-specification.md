---
name: FV 2.0 프로젝트 사양서
description: FreeLang + V 언어 통합 - 새로운 백엔드 언어 개발
type: project
---

# FV 2.0: FreeLang V Language Integration

**프로젝트명**: FV 2.0 (FreeLang V-compatible)
**시작일**: 2026-03-19
**목표**: V 언어 문법 + FreeLang 기능성 통합
**기간**: 8-10주
**상태**: 🟢 기획 완료, 개발 준비 단계

---

## 🎯 핵심 비전

### 문제 인식
```
V 언어: 빠르고 간단하지만, 백엔드 라이브러리 부족
FreeLang: 강력한 백엔드지만, 문법이 복잡

해결책: FV 2.0
= V의 문법 + FreeLang의 기능성
= "가장 빠르고 안전한 백엔드 언어"
```

### 목표 상태
```
FV 2.0으로 쓸 수 있는 것:

fn main() {
  mut server = http.Server{
    addr: '0.0.0.0:8080',
    handler: handle_request,
  }

  db := sql.open('sqlite://app.db')?

  server.listen_and_serve() or {
    eprintln('Server error: $err')
  }
}

fn handle_request(req http.Request) http.Response {
  user := db.query('SELECT * FROM users WHERE id = ?', req.params['id'])?

  return http.Response{
    status_code: 200,
    body: user.to_json(),
  }
}

// 특징:
✅ V 같은 간단한 문법
✅ FreeLang의 HTTP, DB 기능
✅ 에러 처리 (? 와 or)
✅ 타입 안전성
✅ 빠른 컴파일 (<1초)
✅ Docker, K8s 배포 가능
```

---

## 📊 비교: V vs FreeLang vs FV 2.0

### V 언어
```
장점:
✅ 빠른 컴파일 (<1초)
✅ 간단한 문법 (Go 스타일)
✅ 안전성 (bounds checking, null safety)
✅ C 코드 생성 (이식성)
✅ 활발한 커뮤니티

단점:
❌ 백엔드 라이브러리 부족
❌ HTTP framework 미성숙
❌ Database ORM 없음
❌ 프로덕션 사용 사례 적음
```

### FreeLang (현재)
```
장점:
✅ 완전한 HTTP engine
✅ Database ORM
✅ WebSocket, gRPC
✅ JWT/OAuth2
✅ Docker, K8s 배포
✅ 130개 프로젝트 생태계
✅ 프로덕션 준비됨

단점:
❌ 문법이 복잡 (학습곡선)
❌ 커뮤니티 작음
❌ V 언어와 비호환
❌ 독립적인 표준화 필요
```

### FV 2.0 (목표)
```
V의 장점:
✅ 빠른 컴파일
✅ 간단한 문법
✅ 안전성

+ FreeLang의 장점:
✅ 완전한 백엔드
✅ 프로덕션 준비
✅ 130개 프로젝트

= FV 2.0
✅ "실무용 + 배우기 쉬운 언어"
```

---

## 🏗️ 아키텍처

### 현재 (FreeLang)
```
FreeLang 소스
    ↓
Rust 컴파일러 (FreeLang 자체)
    ↓
IR (Intermediate Representation)
    ↓
JIT 또는 C 코드 생성
    ↓
실행
```

### FV 2.0 (목표)
```
FV 2.0 소스 (V 문법)
    ↓
파서 (V 호환)
    ↓
FreeLang IR
    ↓
JIT 또는 C 코드 생성
    ↓
실행

주요 차이:
- 입력: V 문법
- 내부: FreeLang의 IR 사용
- 출력: 동일 (C 또는 바이너리)
```

---

## 📋 상세 로드맵

### Phase 1: 분석 & 설계 (1주)

#### Task 1.1: V 언어 분석
```
목표: V 언어의 구조 파악

내용:
- V 문법 문서 분석
- V 타입 시스템 이해
- V 라이브러리 생태계 조사
- V 컴파일러 아키텍처 분석

산출물:
- V 문법 요약 (20페이지)
- V 타입 시스템 매핑
- V 라이브러리 목록 분류
```

#### Task 1.2: FreeLang 분석
```
목표: FreeLang의 현재 상태 파악

내용:
- FreeLang 문법 정리
- AST 구조 분석
- 타입 시스템 문서화
- 라이브러리 분류

산출물:
- FreeLang 아키텍처 다이어그램
- 현재 문법 명세
- 라이브러리 매핑
```

#### Task 1.3: 통합 지점 설계
```
목표: 어떻게 결합할 것인가?

내용:
- V 문법 → FreeLang AST 매핑
- 호환성 모드 설계
- 마이그레이션 경로
- 문법 충돌 해결

산출물:
- 통합 설계 문서
- 문법 매핑 테이블
- 마이그레이션 전략
```

---

### Phase 2: V 문법 채택 (2주)

#### Task 2.1: 렉서 & 파서 수정
```
목표: V 문법으로 파싱 가능하게

내용:
- 렉서에서 V 키워드 지원
- 파서에서 V 문법 규칙 추가
- AST 호환성 유지
- 테스트 작성

예: V 스타일
  fn main() { ... }    // FreeLang: fn main() -> i64 { ... }
  mut x := 10          // FreeLang: let mut x = 10;
  db := sql.open(...)? // FreeLang: let db = sql::open(...)?
```

#### Task 2.2: 타입 시스템 정렬
```
목표: V와 FreeLang 타입 호환

내용:
- V 타입 ↔ FreeLang 타입 매핑
  - V int ↔ FV i32/i64
  - V string ↔ FV String
  - V []T ↔ FV Vec(T)
  - V ?T ↔ FV Option(T)

- 에러 처리 통일
  - V: result? (결과 전파)
  - FV: Result(T, E)
  - 매핑: ? → Result 언래핑
```

#### Task 2.3: 호환성 테스트
```
목표: V 문법 코드가 컴파일되나?

내용:
- V 코드 샘플 수집
- 각 샘플을 FV 2.0으로 컴파일
- 호환율 측정
- 부족한 부분 식별

성공 기준:
- 90% 이상의 V 코드가 컴파일
- 컴파일 오류 명확함
```

---

### Phase 3: 라이브러리 통합 (4주)

#### Task 3.1: HTTP 라이브러리
```
목표: V 스타일로 HTTP 사용

현재 (FreeLang):
  http.Server { port: 8080, ... }

목표 (FV 2.0, V 스타일):
  mut server := http.Server{
    addr: '0.0.0.0:8080',
    handler: fn(req http.Request) http.Response { ... },
  }
  server.listen_and_serve()?

내용:
- V http 라이브러리 문법 채택
- FreeLang HTTP engine과 통합
- 라우팅 (GET /users/:id)
- 미들웨어
- CORS, 압축 등
```

#### Task 3.2: Database ORM
```
목표: V 스타일로 DB 사용

현재 (FreeLang):
  db.find::<User>(id)?

목표 (FV 2.0, V 스타일):
  user := db.get<User>(id) or { return error }
  users := db.query('SELECT * FROM users')?

내용:
- V database 문법 채택
- FreeLang ORM 재구현
- 쿼리 빌더
- 마이그레이션
- 트랜잭션
```

#### Task 3.3: WebSocket & gRPC
```
목표: 실시간 통신

FV 2.0 코드:
  ws := http.WebSocket{
    on_message: fn(msg string) {
      broadcast(msg)
    },
  }

내용:
- WebSocket (RFC 6455)
- gRPC (Protocol Buffers)
- 양방향 통신
- 에러 처리
```

#### Task 3.4: 보안 & 암호화
```
목표: 프로덕션급 보안

내용:
- JWT/OAuth2 (V 스타일)
- SSL/TLS
- 입력 검증
- XSS, CSRF 방지
```

---

### Phase 4: 마케팅 & 배포 (2주)

#### Task 4.1: 문서화
```
목표: FV 2.0 사용 설명서

내용:
- FV 2.0 시작 가이드
- V 사용자를 위한 가이드
- API 문서
- 예제 모음
```

#### Task 4.2: 130개 프로젝트 마이그레이션
```
목표: FreeLang 프로젝트 → FV 2.0

내용:
- 자동 마이그레이션 도구
- 각 프로젝트 테스트
- 성능 벤치마크
```

#### Task 4.3: 커뮤니티 연동
```
목표: V 커뮤니티와 연결

내용:
- V 포럼/Reddit에 FV 2.0 소개
- V 개발자 대상 블로그 포스팅
- V 생태계에 기여
- 상호 홍보
```

---

## 📈 성공 지표

| 지표 | 목표 | 측정 시점 |
|------|------|----------|
| V 문법 호환율 | 95% | Week 2 끝 |
| 라이브러리 완성도 | 100% HTTP, 100% DB | Week 6 끝 |
| 테스트 커버리지 | 90% | Week 7 |
| 프로덕션 배포 | Docker, K8s 지원 | Week 8 |
| 커뮤니티 반응 | V 커뮤니티 긍정평 | Week 9 |

---

## 💡 핵심 통찰

### "왜 FV 2.0인가?"

```
기존 언어들의 문제:
- V: 빠르지만, 백엔드 약함
- Python: 느리지만, 라이브러리 많음
- Go: 빠르지만, 복잡함
- Rust: 너무 어려움

FV 2.0:
= V의 단순함 + FreeLang의 강점
= "쉽고 빠르고 완전한 백엔드"
```

### "프로덕션에 필요한가?"

```
당신의 의도:

FreeLang만으로도 가능하지만,
V 커뮤니티와 연결하면:

✅ V 개발자들이 백엔드 선택 가능
✅ V 생태계 확장
✅ FreeLang 채택률 증가
✅ 실제 프로덕션 사용 사례 증가

= 상호 win-win
```

---

## 🚀 다음 단계

### Week 1 목표
1. V 언어 완전 분석 완료
2. 통합 설계 문서 작성
3. Phase 1 완료

### 당장 해야 할 것
1. V 문서 읽기 (vlang.io)
2. V 예제 코드 수집
3. 통합 지점 정리

---

**상태**: 🟢 기획 완료, Week 1 시작 준비 완료

**예상 결과**: FV 2.0 (V-compatible, production-ready backend language)

**최종 비전**:
> "V 개발자도, FreeLang 팬도 모두 행복한 언어"
