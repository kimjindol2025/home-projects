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

## 🔄 로드맵

### Phase 1: Lexer ✅ 진행 중
- [x] TokenType 정의 (50+ 토큰)
- [x] Lexer 구현 (렉서 로직)
- [ ] 단위 테스트 작성
- [ ] 성능 최적화

### Phase 2: Parser (예정)
- [ ] AST 노드 정의
- [ ] Recursive Descent Parser 구현
- [ ] 표현식 파싱
- [ ] 문장 파싱

### Phase 3: Compiler (예정)
- [ ] IR (중간언어) 설계
- [ ] 코드 생성
- [ ] 최적화

### Phase 4: Runtime VM (예정)
- [ ] 가상머신 구현
- [ ] 메모리 관리
- [ ] 가비지 컬렉션

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
