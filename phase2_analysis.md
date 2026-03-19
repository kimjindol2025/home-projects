# Phase 2: 공통 모듈 추출 & Standard Core 설계

생성일: 2026-03-19T12:18:07.360937

## 📊 분석 개요

- 분석된 FreeLang 파일: 3194개
- 식별된 함수: 22568개
- 식별된 모듈: 423개
- 디자인 패턴: 18개

## 🔥 상위 함수 (재사용도 높음)

- `new`: 1395개 프로젝트에서 사용
- `main`: 994개 프로젝트에서 사용
- `clone`: 396개 프로젝트에서 사용
- `add`: 331개 프로젝트에서 사용
- `tokenize`: 138개 프로젝트에서 사용
- `stats`: 128개 프로젝트에서 사용
- `isDigit`: 127개 프로젝트에서 사용
- `isAlpha`: 126개 프로젝트에서 사용
- `describe`: 112개 프로젝트에서 사용
- `isAlphaNumeric`: 110개 프로젝트에서 사용
- `assert`: 97개 프로젝트에서 사용
- `peek`: 92개 프로젝트에서 사용
- `parsePrimary`: 88개 프로젝트에서 사용
- `test`: 80개 프로젝트에서 사용
- `advance`: 79개 프로젝트에서 사용
- `validate`: 79개 프로젝트에서 사용
- `isWhitespace`: 71개 프로젝트에서 사용
- `parseBlock`: 71개 프로젝트에서 사용
- `parseProgram`: 69개 프로젝트에서 사용
- `to_string`: 67개 프로젝트에서 사용
- `get`: 67개 프로젝트에서 사용
- `panic`: 66개 프로젝트에서 사용
- `getCurrentTime`: 66개 프로젝트에서 사용
- `printIR`: 62개 프로젝트에서 사용
- `parseFnDecl`: 61개 프로젝트에서 사용
- `factorial`: 55개 프로젝트에서 사용
- `close`: 55개 프로젝트에서 사용
- `parseStatement`: 53개 프로젝트에서 사용
- `parseExpression`: 53개 프로젝트에서 사용
- `push`: 53개 프로젝트에서 사용

## 🏗️ 상위 모듈 (반복되는 구조)

- `tests`: 288개 파일에서 사용
- `unforgiving_rules`: 14개 파일에서 사용
- `n`: 7개 파일에서 사용
- `q`: 6개 파일에서 사용
- `p`: 5개 파일에서 사용
- `predictor`: 4개 파일에서 사용
- `ttf_parser`: 4개 파일에서 사용
- `consensus`: 3개 파일에서 사용
- `neural_network`: 3개 파일에서 사용
- `parser`: 3개 파일에서 사용
- `blog_db`: 3개 파일에서 사용
- `blog_api`: 3개 파일에서 사용
- `blog_server`: 3개 파일에서 사용
- `http_engine`: 3개 파일에서 사용
- `simd`: 3개 파일에서 사용
- `core`: 3개 파일에서 사용
- `challenge_16_tests`: 3개 파일에서 사용
- `module_security`: 3개 파일에서 사용
- `module`: 3개 파일에서 사용
- `distributed_tracer`: 2개 파일에서 사용

## 🎨 디자인 패턴 분포

- `struct`: 1285회
- `while_loop`: 1118회
- `for_loop`: 800회
- `impl_block`: 542회
- `public_function`: 506회
- `public_struct`: 461회
- `unit_test`: 365회
- `enum`: 325회
- `result_pattern`: 309회
- `pattern_matching`: 280회
- `option_pattern`: 275회
- `derive_macro`: 149회
- `await`: 97회
- `trait`: 87회
- `loop`: 42회

## 🎯 Standard Core 라이브러리 설계

### 핵심 모듈 (12개)

#### ASYNC
- 설명: 비동기
- 함수: async, await, Future, Executor, spawn...
- 실제 사용 빈도: 2회
- 예상 코드량: 400줄

#### COLLECTIONS
- 설명: 기본 자료구조
- 함수: Vec, HashMap, HashSet, LinkedList, BTreeMap...
- 실제 사용 빈도: 0회
- 예상 코드량: 400줄

#### CONCURRENCY
- 설명: 동시성
- 함수: Mutex, RwLock, Channel, sync, thread...
- 실제 사용 빈도: 0회
- 예상 코드량: 350줄

#### CRYPTO
- 설명: 암호화
- 함수: sha256, md5, encrypt, decrypt, hash...
- 실제 사용 빈도: 17회
- 예상 코드량: 450줄

#### IO
- 설명: 입출력
- 함수: File, read, write, println, BufReader...
- 실제 사용 빈도: 41회
- 예상 코드량: 250줄

#### ITER
- 설명: 반복자
- 함수: Iterator, map, filter, fold, zip...
- 실제 사용 빈도: 15회
- 예상 코드량: 300줄

#### MATH
- 설명: 수학 연산
- 함수: sqrt, sin, cos, pow, abs...
- 실제 사용 빈도: 100회
- 예상 코드량: 150줄

#### NETWORK
- 설명: 네트워크
- 함수: TcpStream, TcpListener, http, request...
- 실제 사용 빈도: 0회
- 예상 코드량: 500줄

#### OPTION
- 설명: 선택적 값
- 함수: Option, Some, None, map, unwrap...
- 실제 사용 빈도: 20회
- 예상 코드량: 120줄

#### RESULT
- 설명: 에러 처리
- 함수: Result, Ok, Err, map, unwrap...
- 실제 사용 빈도: 28회
- 예상 코드량: 120줄

#### SERIALIZATION
- 설명: 직렬화
- 함수: serialize, deserialize, json, serde...
- 실제 사용 빈도: 25회
- 예상 코드량: 400줄

#### STRING
- 설명: 문자열 처리
- 함수: String, str, trim, split, join...
- 실제 사용 빈도: 67회
- 예상 코드량: 300줄

### 총 예상 Standard Core 규모

- **코드량**: 3740줄
- **모듈 수**: 12개
- **개발 기간**: 2-3주 (병렬 개발 가능)
- **테스트**: 각 모듈당 50-100개 테스트 케이스 (총 800-1200개)

## 🔗 모듈 간 의존성


```
result ──────┐
option ──────┤
             ├──> collections (Vec, HashMap, ...)
iter ────────┤       ↓
             ├──> string (문자열 처리)
string ──────┤       ↓
math ────────┤──> io (파일, 콘솔)
             ↓
             ├──> serialization (JSON)
             ├──> crypto (암호화)
             ├──> network (HTTP, TCP)
             └──> async/concurrency
```

## 📈 마이그레이션 로드맵


### Phase 2-A: Standard Core 설계 (1주)
- [x] 함수/모듈 패턴 분석 (완료)
- [ ] API 설계 (공개 인터페이스)
- [ ] 문서화 및 예제 작성

### Phase 2-B: 핵심 모듈 구현 (2주)
**우선순위 1** (기초):
- collections (Vec, HashMap)
- result/option (에러 처리)
- string (문자열)

**우선순위 2** (활용):
- io (파일 입출력)
- iter (반복자)
- math (수학)

**우선순위 3** (고급):
- async/concurrency (병렬)
- serialization (직렬화)
- crypto (암호화)
- network (통신)

### Phase 2-C: 130개 프로젝트 마이그레이션 (3주)
- 자동 코드 변환 스크립트
- 호환성 테스트
- 점진적 롤아웃

### Phase 3: Z3 SMT 검증 (진행 예정)
