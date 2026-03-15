# JuliaCC - Julia Compiler Complete

Julia 언어의 완전한 컴파일러 구현 (Go 기반)

## 개요

JuliaCC는 Julia 공식 컴파일러의 파이프라인을 기반으로 Go 언어로 구현한 고성능 컴파일러입니다.
자체호스팅(Self-Hosting)을 목표로 하며, Julia 코드를 LLVM IR을 거쳐 네이티브 코드로 컴파일합니다.

```
Julia Source
    ↓
[Lexer]    - 토큰화
    ↓
[Parser]   - AST 생성
    ↓
[Analyzer] - 의미 분석
    ↓
[Lowering] - Untyped IR
    ↓
[Type Inference] - Type SSA IR
    ↓
[Optimizer] - Julia SSA 최적화
    ↓
[Codegen]  - LLVM IR 생성
    ↓
[LLVM]     - 최적화 + 어셈블리
    ↓
Native Code
```

## 빌드 및 실행

### 필수 요구사항
- Go 1.20+
- make
- (선택) LLVM 12+

### 빌드

```bash
make build
# 또는
go build -o bin/jcc cmd/jcc/main.go
```

### 테스트

```bash
make test          # 기본 테스트
make test-verbose  # 상세 테스트 (race detector 포함)
```

### 예제 실행

```bash
make run
# 또는
./bin/jcc examples/hello.jl
```

## 프로젝트 구조

```
julia-compiler/
├── cmd/
│   └── jcc/                 # 메인 CLI 도구
│       └── main.go
├── internal/
│   ├── lexer/              # 토큰화 (Phase 1)
│   ├── parser/             # 구문 분석 (Phase 2)
│   ├── ast/                # 추상 구문 트리
│   ├── types/              # 타입 시스템 (Phase 3)
│   ├── ir/                 # 중간 표현 (Phase 5-6)
│   ├── codegen/            # 코드 생성 (Phase 7)
│   └── runtime/            # 런타임 (Phase 8)
├── pkg/
│   └── stdlib/             # 표준 라이브러리
├── test/
│   ├── fixtures/           # 테스트 코드 샘플
│   └── benchmarks/         # 성능 벤치마크
├── examples/               # 예제 코드
├── docs/                   # 설명서
├── Makefile
├── go.mod
└── README.md
```

## 개발 로드맵

### Phase 0: 초기화 ✅ 진행 중
- [x] 프로젝트 구조 설정
- [x] Go 모듈 초기화
- [ ] CI/CD 파이프라인 (GitHub Actions 또는 GitLab CI)
- [ ] 개발 가이드 작성

### Phase 1: Lexer 🔤 (500-800줄)
- [ ] 키워드 정의
- [ ] 연산자 정의
- [ ] 리터럴 파싱
- [ ] 에러 위치 추적

### Phase 2: Parser 🌳 (1,200-1,500줄)
- [ ] AST 노드 정의
- [ ] 식 파싱 (이항/단항 연산)
- [ ] 문 파싱 (제어 흐름)
- [ ] 함수 정의 파싱

### Phase 3: 타입 시스템 🎯 (1,000-1,300줄)
- [ ] 기본 타입 (Int64, Float64, String, ...)
- [ ] 복합 타입 (Vector{T}, Matrix{T}, ...)
- [ ] 다중 디스패치 구현
- [ ] 타입 계층 구조

**...** (총 10 Phase)

자세한 로드맵은 [julia-compiler-project.md](./docs/julia-compiler-project.md) 참고

## 의존성

주요 Go 패키지:
- `stretchr/testify` - 테스트 프레임워크
- (추가 예정) LLVM 바인딩

## 성능 목표

- **컴파일 속도**: < 100ms (1000줄 기준)
- **실행 속도**: Python 대비 10-100배 빠름
- **메모리 효율**: 타입 안정성으로 최소화

## 기여

프로젝트 기여 안내:
1. Fork 저장소
2. 기능 브랜치 생성 (`git checkout -b feature/amazing-feature`)
3. 변경사항 커밋 (`git commit -m 'Add amazing feature'`)
4. 브랜치 푸시 (`git push origin feature/amazing-feature`)
5. Pull Request 생성

## 라이선스

MIT License - [LICENSE](./LICENSE) 참고

## 참고 자료

- [Julia 공식 문서](https://docs.julialang.org)
- [Julia Compiler 개발 문서](https://docs.julialang.org/en/v1/devdocs/)
- [LLVM Language Reference](https://llvm.org/docs/LangRef/)

## 연락처

프로젝트 관련 질문이나 제안:
- Issue Tracker: [GitHub Issues](https://github.com/yourusername/juliacc/issues)
- 토론: [GitHub Discussions](https://github.com/yourusername/juliacc/discussions)

---

**상태**: Phase 0 진행 중 (2026-03-11)
**버전**: 0.1.0-alpha
**최종 업데이트**: 2026-03-11
