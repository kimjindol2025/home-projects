# GoFree 🚀

**Go-based FreeLang Compiler**

TypeScript 구현 → Go 이식 프로젝트

---

## 📊 프로젝트 개요

| 항목 | 설명 |
|------|------|
| **원본** | [v2-freelang-ai](https://github.com/freelang-ai/v2-freelang-ai) (TypeScript) |
| **목표** | Go로 완전 이식 (성능 + 배포 개선) |
| **상태** | 🔄 Phase 1: Lexer 이식 중 |
| **예상** | 5-6주 (4인 팀 기준) |

---

## 🏗️ 구조

```
gofree/
├── internal/
│   ├── lexer/          # Phase 1: 토큰화 ✅
│   ├── parser/         # Phase 2: AST 생성 (예정)
│   ├── compiler/       # Phase 3: IR 생성 (예정)
│   ├── runtime/        # Phase 4: VM 실행 (예정)
│   ├── analyzer/       # Phase 5: 의미분석 (예정)
│   └── ...
├── cmd/
│   └── freelang/       # CLI 엔트리 포인트
├── stdlib/             # 표준 라이브러리
├── examples/           # 예제 코드
├── tests/              # 테스트 스위트
├── go.mod
├── go.sum
└── README.md
```

---

## 🔄 상세 로드맵

**기반**: FreeLang TypeScript v2.8.0 → Go 1.21+
**총 기간**: 3~6개월 (11 Phase)
**최종 목표**: Self-Hosting (compiler.free를 자기 자신으로 컴파일)

### 📋 Phase 요약

| Phase | 모듈 | 기간 | 난이도 | 병렬화 | 상태 |
|-------|------|------|--------|--------|------|
| **P0** | 초기화 | 1일 | ⭐ | - | ✅ 완료 |
| **P1** | Lexer | 3-5일 | ⭐⭐ | P2와 병렬 | ✅ 완료 |
| **P2** | AST | 2-3일 | ⭐⭐ | P1과 병렬 | ✅ 완료 |
| **P3** | Parser | 7-10일 | ⭐⭐⭐⭐ | P4와 병렬 | 🔄 진행 중 |
| **P4** | Analyzer | 5-7일 | ⭐⭐⭐ | P3과 병렬 | ⏳ 대기 |
| **P5** | IR Compiler | 7-10일 | ⭐⭐⭐⭐ | P6과 병렬 | ⏳ 대기 |
| **P6** | Runtime VM | 10-14일 | ⭐⭐⭐⭐⭐ | P5와 병렬 | ⏳ 대기 |
| **P7** | Codegen | 7-10일 | ⭐⭐⭐⭐ | P5/P6 후 | ⏳ 대기 |
| **P8** | Formatter/Linter | 5-7일 | ⭐⭐⭐ | P6 후 | ⏳ 대기 |
| **P9** | GC & Memory | 10-14일 | ⭐⭐⭐⭐⭐ | 병렬 불가 | ⏳ 대기 |
| **P10** | Self-Hosting | 3-4주 | ⭐⭐⭐⭐⭐ | 마지막 | ⏳ 대기 |

### 🎯 마일스톤 & 성공 지표

**M1: 첫 실행 (Phase 6 완료 후)**
```
✅ go run cmd/freelang/main.go examples/hello.fl
   → "Hello, FreeLang!" 출력
✅ 테스트: 5개 예제 파일 모두 실행 성공
✅ 벤치마크: hello.fl 파싱 + 실행 < 100ms
```

**M2: 기본 언어 (Phase 8 완료 후)**
```
✅ 변수/함수/제어흐름 모두 동작
✅ 100줄 이상의 .free 파일 실행 가능
✅ for/while/if/match 모두 지원
✅ 테스트: stdlib 함수 20+ 지원
```

**M3: 완전 언어 (Phase 9 완료 후)**
```
✅ async/await, try/catch 동작
✅ 메모리 누수 0 (valgrind 수준)
✅ 5000줄 .free 파일 파싱 < 1s
✅ 테스트: 1GB 힙 메모리 관리
```

**M4: 도구 완성 (Phase 8 완료 후)**
```
✅ freelang fmt --check: 코드 포맷 검사
✅ freelang lint: ESLint 완전 대체
✅ freelang test: 모든 test{} 블록 실행
✅ freelang build --release: 최적화 빌드
```

**M5: 셀프호스팅 (Phase 10 완료 후)**
```
✅ compiler.free → Go-FreeLang으로 컴파일
✅ 컴파일된 바이너리: 원본과 동일 동작
✅ .flir 공통 IR: TS/Go 컴파일러 모두 호환
✅ @git_hook 어노테이션: Git 자동 검증
```

### 📊 성능 목표 (Benchmarks)

**Phase 3 (Parser)**:
- 파싱 속도: < 10ms (hello.free)
- 메모리 사용: < 1MB
- AST 노드 개수: 정확성 100%

**Phase 5 (IR Compiler)**:
- 컴파일 속도: < 20ms (hello.free)
- IR 크기: < 2KB
- 최적화: 상수 폴딩 + 데드 코드 제거

**Phase 6 (Runtime VM)**:
- 실행 속도: < 10ms (hello.free)
- 메모리 피크: < 10MB (1000 변수)
- 함수 호출 오버헤드: < 1μs

**Phase 9 (GC)**:
- GC 일시 정지 시간: < 10ms
- 메모리 해제율: > 95%
- 단편화율: < 10%

### 🔄 병렬 작업 계획

**레벨 1: P1 + P2** (완료)
```
Day 1-5: Lexer + AST 동시 진행 ✅
```

**레벨 2: P3 + P4** (진행 중)
```
Week 1-2: P3 (Parser) 핵심 구현
Week 2: P4 (Analyzer) 시작 (스코프 관리)
Week 3: P3 완료 → P4 타입 검사 병렬
```

**레벨 3: P5 + P6** (예정)
```
Week 4-5: P5 (IR Compiler) + P6 (VM) 동시
Week 6: P5 IR → P6 VM 통합 테스트
```

### 🐍 Python 통합 전략

**시점 1: Phase 3 (Parser) 완료 후**
- cgo를 통한 C 바인딩
- Python 타입 → Go 타입 자동 변환
- 양방향: Python → Go Lexer/Parser → JSON AST

**시점 2: Phase 6 (Runtime) 완료 후**
- GoFree VM을 Python 임베드 가능하게
- CFFI/ctypes로 Go 함수 노출
- 고성능 데이터 교환 (protobuf/msgpack)

**시점 3: Phase 7 (Codegen) 완료 후**
- Python 패키지로 배포 (`pip install gofree`)
- Go 바이너리 자동 다운로드
- Python 래퍼 제공

---

## 🛠️ 개발 환경

**요구사항**:
- Go 1.21+
- git
- Linux/macOS/Windows (WSL2)

**빌드**:
```bash
go build -o freelang cmd/freelang/main.go
```

**실행**:
```bash
./freelang examples/hello.fl
```

**테스트**:
```bash
go test ./...
```

---

## 📝 라이선스

MIT License (FreeLang 원본 따름)

---

## 🤝 기여

PR/Issue는 언제든 환영합니다!

- Slack: `#freelang-go-dev`
- Email: `kim@freelang-ai.dev`

---

**Last Updated**: 2026-03-10
