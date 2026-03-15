# GoFree: 상세 이식 로드맵

**기반**: FreeLang TypeScript v2.8.0 → Go 1.21+
**총 기간**: 3~6개월 (11 Phase)
**최종 목표**: Self-Hosting (compiler.free를 자기 자신으로 컴파일)

---

## 📋 Phase 요약표

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

---

## 🎯 마일스톤 & 성공 지표

### M1: 첫 실행 (Phase 6 완료 후)
```
성공 기준:
✅ go run cmd/freelang/main.go examples/hello.fl
   → "Hello, FreeLang!" 출력
✅ 테스트: 5개 예제 파일 모두 실행 성공
✅ 벤치마크: hello.fl 파싱 + 실행 < 100ms
```

### M2: 기본 언어 (Phase 8 완료 후)
```
성공 기준:
✅ 변수/함수/제어흐름 모두 동작
✅ 100줄 이상의 .free 파일 실행 가능
✅ for/while/if/match 모두 지원
✅ 테스트: stdlib 함수 20+ 지원
```

### M3: 완전 언어 (Phase 9 완료 후)
```
성공 기준:
✅ async/await, try/catch 동작
✅ 메모리 누수 0 (valgrind 수준)
✅ 5000줄 .free 파일 파싱 < 1s
✅ 테스트: 1GB 힙 메모리 관리
```

### M4: 도구 완성 (Phase 8 완료 후)
```
성공 기준:
✅ freelang fmt --check: 코드 포맷 검사
✅ freelang lint: ESLint 완전 대체
✅ freelang test: 모든 test{} 블록 실행
✅ freelang build --release: 최적화 빌드
```

### M5: 셀프호스팅 (Phase 10 완료 후)
```
성공 기준:
✅ compiler.free → Go-FreeLang으로 컴파일
✅ 컴파일된 바이너리: 원본과 동일 동작
✅ .flir 공통 IR: TS/Go 컴파일러 모두 호환
✅ @git_hook 어노테이션: Git 자동 검증
```

---

## 📊 성능 목표 (Benchmarks)

### Phase 3 (Parser)
```
벤치마크 (hello.free 기준):
- 파싱 속도: < 10ms
- 메모리 사용: < 1MB
- AST 노드 개수: 정확성 100%

기준:
$ go test -bench=BenchmarkParser ./internal/parser
BenchmarkParser-8  100000  10023 ns/op  ≈ 10ms
```

### Phase 5 (IR Compiler)
```
벤치마크:
- 컴파일 속도: < 20ms (hello.free)
- IR 크기: < 2KB
- 최적화: 상수 폴딩 + 데드 코드 제거

목표: 1000줄 .free 파일 < 100ms
```

### Phase 6 (Runtime)
```
벤치마크:
- 실행 속도: < 10ms (hello.free)
- 메모리 피크: < 10MB (1000 변수)
- 함수 호출 오버헤드: < 1μs

목표: fibonacci(50) < 5s
```

### Phase 9 (GC)
```
벤치마크:
- GC 일시 정지 시간: < 10ms
- 메모리 해제율: > 95%
- 단편화율: < 10%

목표: 1시간 지속 실행 중 메모리 누수 0
```

---

## 🔄 병렬 작업 계획

### 레벨 1: P1 + P2 (병렬 가능)
```
Timeline:
Day 1-3:  P1 (Lexer) + P2 (AST) 동시 진행
Day 4-5:  P1 완료, P2 완료
Day 6+:   P3 (Parser) 시작
```

### 레벨 2: P3 + P4 (부분 병렬)
```
Timeline:
Week 1-2:  P3 (Parser) 핵심 구현
Week 2:    P4 (Analyzer) 시작 (스코프 관리)
Week 3:    P3 완료 → P4 타입 검사 병렬
```

### 레벨 3: P5 + P6 (완전 병렬)
```
Timeline:
Week 4-5:  P5 (IR Compiler) + P6 (VM) 동시
           - 팀원 A: P5 (Opcode 정의)
           - 팀원 B: P6 (VM 루프)
Week 6:    P5 IR → P6 VM 통합 테스트
```

### 레벨 4: P7 병렬 (Optional)
```
Timeline:
Week 5-6:  P7 (Codegen) 시작 (선택사항)
           - C 코드 생성 구현
           - WASM 생성 구현
```

**병렬화 불가능**:
- P9 (GC): P6 완료 후 메모리 구조 확정 필요
- P10 (Self-Hosting): P1~P9 모두 완료 필수

---

## 🐍 Python 통합 전략

### 시점 1: Phase 3 (Parser) 완료 후
```
목표: Python에서 GoFree 파서 호출 가능
구현:
- cgo를 통한 C 바인딩
- Python 타입 → Go 타입 자동 변환
- 양방향: Python → Go Lexer/Parser → JSON AST

사용 사례:
from gofree import Lexer
tokens = Lexer("let x = 10").tokenize()
```

### 시점 2: Phase 6 (Runtime) 완료 후
```
목표: Python에서 FreeLang 코드 직접 실행
구현:
- GoFree VM을 Python 임베드 가능하게
- CFFI/ctypes로 Go 함수 노출
- 고성능 데이터 교환 (protobuf/msgpack)

사용 사례:
result = gofree.run("let x = 10 + 5")
```

### 시점 3: Phase 7 (Codegen) 완료 후
```
목표: Python 패키지로 배포
구현:
- pip install gofree
- Go 바이너리 자동 다운로드
- Python 래퍼 제공

사용 사례:
$ pip install gofree
$ python -c "from gofree import compile; compile('hello.free')"
```

---

## ✅ 각 Phase별 통합 테스트

### P3 Parser 테스트
```go
// internal/parser/parser_test.go
func TestHelloFree(t *testing.T) {
    src := `fn main() { println("Hello, FreeLang!") }`
    ast, err := Parse(Lex(src))
    assert.NoError(t, err)
    assert.Equal(t, 1, len(ast.Functions))
    assert.Equal(t, "main", ast.Functions[0].Name)
}

func TestComplexExpression(t *testing.T) {
    // 100줄 복잡한 표현식 테스트
}

func BenchmarkParser(b *testing.B) {
    // 성능 벤치마크
}
```

### P6 Runtime 통합 테스트
```bash
# 모든 예제 파일 자동 테스트
for file in examples/*.free; do
    go run cmd/freelang/main.go "$file" > /tmp/out.txt
    diff -u "$file.expected" /tmp/out.txt || exit 1
done

# 성능 테스트
time go run cmd/freelang/main.go examples/fibonacci.free
# 목표: fibonacci(50) < 5s
```

### P8 Linter 통합 테스트
```bash
# Lint 규칙 검증
freelang lint examples/unused_var.free
# 예상: [lint] unused_var.free:2:5 ✘ [no_unused] Variable 'x' is never used

# 포맷 일관성 검증
freelang fmt examples/messy.free
diff -u examples/messy.free examples/messy.free.formatted
```

---

## ⚠️ 위험 요소 & 완화 전략

### 위험 1: Parser 복잡도
```
위험도: 높음
영향: Phase 3 2-3주 지연 가능
완화:
- TS 파서 코드 직접 참고 (변환)
- 단위 테스트: 각 구문별 10+ 케이스
- 점진적 구현: 기본 문법 → 고급 문법
```

### 위험 2: 메모리 관리
```
위험도: 높음
영향: Phase 6 실패 가능
완화:
- Phase 9 GC 조기 설계 (Phase 5에서)
- Go 런타임 GC 활용 (초기)
- 메모리 프로파일링: pprof 사용
```

### 위험 3: Python 통합
```
위험도: 중간
영향: Python 사용자 차단
완화:
- cgo 대신 REST API 사용 (폴백)
- gRPC로 언어 독립적 통신
- Docker 컨테이너화
```

### 위험 4: Self-Hosting 복잡도
```
위험도: 높음
영향: Phase 10 2-4주 지연 가능
완화:
- Phase 10 초기부터 lexer.free 작성
- 단계별 부트스트랩 (Stage 1-4)
- .flir 공통 IR로 검증
```

---

## 📈 진행 추적 (Git/GOGS)

각 Phase 완료마다 태그 생성:
```bash
# Phase 3 완료 후
git tag -a v0.3.0 -m "Phase 3: Parser ✅"
git tag -a m1-ready -m "M1 마일스톤 준비"
git push origin --tags

# 각 Phase 브랜치
git checkout -b phase/3-analyzer
```

GOGS 저장소:
```
https://gogs.dclub.kr/kim/gofree
├── master (주 개발)
├── phase/3-analyzer (현재)
├── phase/4-ir-compiler (예정)
└── tags/v0.3.0, v0.6.0 등 (마일스톤)
```

---

## 🗂️ 파일 구조 (최종)

```
gofree/
├── internal/
│   ├── lexer/        (P1)  ✅
│   ├── parser/       (P3)  🔄
│   ├── ast/          (P2)  ✅
│   ├── analyzer/     (P4)  ⏳
│   ├── compiler/     (P5)  ⏳
│   ├── runtime/      (P6)  ⏳
│   ├── codegen/      (P7)  ⏳
│   ├── formatter/    (P8)  ⏳
│   ├── linter/       (P8)  ⏳
│   ├── gc/           (P9)  ⏳
│   └── stdlib/       (P6+) ⏳
├── cmd/freelang/     (CLI)  ✅
├── examples/         (테스트) ✅
├── tests/           (통합테스트) ⏳
├── bench/           (벤치마크) ⏳
├── docs/            (설명서) ⏳
├── ROADMAP.html     (일반 로드맵)
├── ROADMAP_DETAILED.md  (이 파일)
├── go.mod           ✅
└── Makefile         ✅
```

---

## 🚀 시작하기

### 즉시 (이번 주)
```bash
# Phase 3 브랜치 생성
git checkout -b phase/3-analyzer

# Analyzer 뼈대 생성
mkdir -p internal/analyzer
touch internal/analyzer/analyzer.go
touch internal/analyzer/analyzer_test.go

# 첫 커밋
git commit -m "🌱 Phase 3: Analyzer 시작"
```

### 조직 (매주)
```
월요일: Phase 목표 정의 & 테스트 작성
화~목요일: 핵심 구현
금요일: 통합 테스트 & 마일스톤 확인
```

---

**이 세분화 로드맵으로 시작할 준비 됐어요?** 💪