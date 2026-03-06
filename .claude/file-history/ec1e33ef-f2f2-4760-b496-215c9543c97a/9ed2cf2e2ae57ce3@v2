# 👥 FreeLang 100% 3명 팀 역할 분배

**팀 구성**: 3명 (전담)
**기간**: 10개월 (2026-03-08 ~ 2026-12-31)
**목표**: 100% 완성도 달성
**관리 방식**: Sprint 기반 (2주 단위), 주 1회 동기화

---

## 👤 **팀원 1: Core Language & Compiler 리더**

**직책**: CTO / Language Architect
**경험**: 컴파일러 또는 인터프리터 개발 경험 필수
**급여**: 월 $35K

### **담당 영역**

```
1. Self-Compiling Compiler (3000줄)
   ├─ AST 정의 및 변환
   ├─ Type checking engine
   ├─ Bytecode code generation
   ├─ Optimization passes
   └─ Bootstrap mechanism

2. Type System 완성 (1500줄)
   ├─ Generics 구현
   ├─ Trait system
   ├─ Type inference
   └─ Constraint solving

3. Error Handling (800줄)
   ├─ Custom exceptions
   ├─ Stack trace generation
   ├─ Error recovery
   └─ Diagnostic messages

4. Runtime & VM (1200줄)
   ├─ Bytecode VM 최적화
   ├─ Memory layout optimization
   ├─ GC tuning
   └─ Native interface

총 책임: 6500줄
```

### **월별 마일스톤**

```
Month 1 (Mar 8-Apr 7): Compiler 설계 및 기초
├─ Architecture document 작성
├─ AST 정의 완료
├─ Parser integration 설계
├─ Weekly review with team
└─ Deliverable: Compiler design doc (50페이지)

Month 2 (Apr 8-May 7): Compiler 구현 (50%)
├─ Type checking engine
├─ Code generation (bytecode)
├─ Basic optimization
├─ First compilation 성공
└─ Deliverable: FL → bytecode compiler 작동

Month 3 (May 8-Jun 7): Compiler 완성 (100%)
├─ Optimization passes 완성
├─ Bootstrap 메커니즘 구현
├─ 버그 수정 및 테스트
├─ TypeScript compiler와 비교 검증
└─ Deliverable: Self-compiling compiler 완성

Month 4 (Jun 8-Jul 7): Type System & Generics
├─ Generic constraint checking
├─ Trait resolution
├─ Type inference engine
├─ Complex type 지원 (tuple, union)
└─ Deliverable: Complete type system

Month 5 (Jul 8-Aug 7): Error Handling & Runtime
├─ Exception hierarchy 설계
├─ Stack trace collection
├─ Runtime error handling
├─ Memory optimization
└─ Deliverable: Production-grade error system

Month 6-10: Maintenance & Optimization
├─ Bug fixes
├─ Performance tuning
├─ Documentation review
├─ Community feedback 대응
└─ Deliverable: Stable 1.0 release
```

### **코드 예시**

```fl
// src/compiler.fl
class Compiler {
  fn compile(source: string) -> bytes {
    // 1. Tokenize
    let tokens = self.lexer.tokenize(source)

    // 2. Parse
    let ast = self.parser.parse(tokens)

    // 3. Type check
    self.type_checker.check(ast)

    // 4. Optimize
    ast = self.optimizer.optimize(ast)

    // 5. Code gen
    self.codegen.generate(ast)
  }

  fn optimize(ast: AST) -> AST {
    // Dead code elimination
    // Constant folding
    // Inline optimization
    // Instruction scheduling
  }
}

// src/bootstrap.fl
fn main(args: string[]) {
  // Step 1: Load TypeScript compiler
  let ts_compiler = load_typescript_compiler()

  // Step 2: Compile FL compiler with TS
  let fl_compiler_v1 = ts_compiler.compile("src/compiler.fl")

  // Step 3: Compile FL compiler with FL v1
  let fl_compiler_v2 = fl_compiler_v1.compile("src/compiler.fl")

  // Step 4: Verify identical
  if sha256(fl_compiler_v2) == sha256(fl_compiler_v1) {
    println("✅ Bootstrap successful!")
    return fl_compiler_v2
  } else {
    println("❌ Bootstrap failed: bytecode differs")
    return error
  }
}
```

### **협력 지점**

```
With Team 2 (Stdlib):
├─ Week 2, 4, 6, 8: Stdlib 함수 시그니처 확정
├─ Type system 변경 시 즉시 공유
└─ Performance profiling 데이터 공유

With Team 3 (Tools):
├─ Week 1: LSP/IDE 요구사항 수집
├─ Month 5-6: Documentation 작성 검토
└─ Month 8-10: Performance optimization 협력
```

### **위험과 대비**

```
위험 1: Bootstrap 실패 (순환 의존성)
├─ 확률: 30%
├─ 임팩트: 매우 높음 (일정 연장)
├─ 대비: 2-stage bootstrap, dynamic linking
└─ 재계획: +2개월 추가 시간 예약

위험 2: Type system 복잡도 증가
├─ 확률: 50%
├─ 임팩트: 중간
├─ 대비: Simplified generics로 시작
└─ 기한: Month 4 완료 (Month 5로 연장 가능)

위험 3: Performance 목표 미달
├─ 확률: 40%
├─ 임팩트: 낮음 (v1.0에는 필수 아님)
├─ 대비: Performance는 Month 6-10에 집중
└─ 기한: v1.0 이후 최적화 가능
```

---

## 👤 **팀원 2: Standard Library & Libraries**

**직책**: Senior Rust/Go Engineer (Systems background)
**경험**: 시스템 라이브러리 개발 경험
**급여**: 월 $32K

### **담당 영역**

```
1. Collections (1500줄)
   ├─ HashMap (hash table + chaining)
   ├─ LinkedList (doubly-linked)
   ├─ TreeMap (balanced BST)
   ├─ PriorityQueue (min-heap)
   ├─ HashSet, TreeSet
   └─ Queue, Stack, Deque

2. String & Text (1000줄)
   ├─ Regex engine (basic)
   ├─ String formatting
   ├─ Unicode support (UTF-8)
   ├─ String split/join/replace
   └─ Regular expressions

3. Data Formats (1500줄)
   ├─ JSON parser & serializer
   ├─ CSV reader/writer
   ├─ YAML parser (basic)
   ├─ msgpack encoder/decoder
   └─ Protocol buffers (basic)

4. System APIs (800줄)
   ├─ File I/O (complete)
   ├─ Directory operations
   ├─ Process management
   ├─ Environment variables
   └─ Command execution

5. Network (1200줄)
   ├─ HTTP client (minimal)
   ├─ TCP socket
   ├─ UDP socket
   ├─ DNS resolver
   └─ TLS/SSL support

총 책임: 6000줄
```

### **월별 마일스톤**

```
Month 1 (Mar 8-Apr 7): Collections 기초
├─ HashMap 구현
├─ LinkedList 구현
├─ HashSet 구현
└─ Deliverable: 3가지 핵심 collection

Month 2 (Apr 8-May 7): Collections 완성
├─ TreeMap, PriorityQueue
├─ Queue, Stack, Deque
├─ 성능 최적화
└─ Deliverable: 8가지 collection 완성

Month 3 (May 8-Jun 7): String & Text + JSON
├─ Regex engine (기본)
├─ String utilities
├─ JSON parser & serializer
└─ Deliverable: JSON I/O 작동

Month 4 (Jun 8-Jul 7): Data Formats 완성
├─ CSV reader/writer
├─ YAML parser
├─ msgpack support
└─ Deliverable: 4가지 format 지원

Month 5 (Jul 8-Aug 7): System & File I/O
├─ File I/O 완성
├─ Directory operations
├─ Process management
└─ Deliverable: Complete filesystem API

Month 6 (Aug 8-Sep 7): Network
├─ HTTP client
├─ Socket programming
├─ DNS resolver
└─ Deliverable: Basic networking

Month 7-10: Optimization & Testing
├─ Performance tuning
├─ Edge case 처리
├─ Comprehensive tests
├─ Documentation
└─ Deliverable: Production-grade stdlib
```

### **코드 예시**

```fl
// lib/collections.fl
class HashMap<K, V> {
  var buckets: array<LinkedList<Pair<K, V>>>
  var size: number

  fn get(key: K) -> V? {
    let hash = key.hash() % buckets.length
    for pair in buckets[hash] {
      if pair.key == key {
        return pair.value
      }
    }
    return null
  }

  fn put(key: K, value: V) -> V? {
    let hash = key.hash() % buckets.length
    for i, pair in enumerate(buckets[hash]) {
      if pair.key == key {
        let old = pair.value
        pair.value = value
        return old
      }
    }
    buckets[hash].push(Pair(key, value))
    size += 1
    if size > buckets.length * 0.75 {
      self.resize()
    }
    return null
  }
}

// lib/json.fl
fn parse_json(input: string) -> JSONValue {
  let parser = JSONParser(input)
  return parser.parse_value()
}

fn stringify_json(value: JSONValue) -> string {
  match value {
    JSONNull => "null"
    JSONBool(b) => b.to_string()
    JSONNumber(n) => n.to_string()
    JSONString(s) => "\"" + escape(s) + "\""
    JSONArray(arr) => "[" + arr.map(stringify_json).join(",") + "]"
    JSONObject(obj) => "{" + ...
  }
}
```

### **협력 지점**

```
With Team 1 (Compiler):
├─ Week 1: Stdlib 함수 시그니처 확정
├─ Monthly: Type system 변경 대응
└─ Month 3+: Stdlib 함수 최적화

With Team 3 (Tools):
├─ Week 4: 문서 형식 정의
├─ Month 6: Performance profiling
└─ Month 8+: Test suite 작성
```

### **위험과 대비**

```
위험 1: JSON parser 복잡도
├─ 확률: 40%
├─ 임팩트: 낮음
├─ 대비: Simple JSON만 먼저 구현
└─ 기한: Month 3에서 Month 4로 연장 가능

위험 2: Performance 달성 (for collections)
├─ 확률: 50%
├─ 임팩트: 중간
├─ 대비: Simple 구현 먼저, 나중 최적화
└─ 기한: Month 7-10에 시간 예약

위험 3: Network 구현 복잡도
├─ 확률: 60%
├─ 임팩트: 낮음 (optional)
├─ 대비: HTTP client만 지원, full TLS 나중
└─ 기한: Month 6에서 Month 7로 연장 가능
```

---

## 👤 **팀원 3: Tools, Ecosystem & DevOps**

**직책**: DevOps / Tools Engineer
**경험**: IDE plugin, package manager, CI/CD 경험
**급여**: 월 $30K

### **담당 영역**

```
1. Package Manager (2500줄)
   ├─ CLI (freelang-pkg)
   ├─ Package registry client
   ├─ Dependency resolver
   ├─ Local cache management
   ├─ Version management
   └─ Publishing tools

2. Language Server Protocol (800줄)
   ├─ LSP server implementation
   ├─ Type hints
   ├─ Go-to-definition
   ├─ Diagnostics
   └─ Code completion (기본)

3. VSCode Extension (1000줄, TypeScript)
   ├─ Syntax highlighting
   ├─ Snippets
   ├─ Debugger integration
   ├─ Testing support
   └─ Theme support

4. CLI Tools (1500줄)
   ├─ freelang fmt (formatter)
   ├─ freelang lint (linter)
   ├─ freelang doc (documentation generator)
   ├─ freelang test runner
   └─ freelang profiler (basic)

5. Registry Server (3000줄, Go/Rust)
   ├─ Package storage
   ├─ User management
   ├─ REST API
   ├─ Search functionality
   └─ Version tracking

6. CI/CD & Infrastructure
   ├─ GitHub Actions workflows
   ├─ Automated testing
   ├─ Release automation
   ├─ Docker images
   └─ Documentation hosting

총 책임: ~9300줄 (+ 3000줄 Go/Rust)
```

### **월별 마일스톤**

```
Month 1 (Mar 8-Apr 7): CLI Tools 기초
├─ freelang fmt (formatter)
├─ Basic linter
├─ CLI structure
└─ Deliverable: fmt + lint 작동

Month 2 (Apr 8-May 7): Package Manager 시작
├─ CLI 인터페이스
├─ Dependency resolver 알고리즘
├─ Local cache
└─ Deliverable: freelang install <pkg> 기초

Month 3 (May 8-Jun 7): Registry Server
├─ Go/Rust로 registry 구현
├─ User management
├─ Package upload API
└─ Deliverable: Central registry 운영 가능

Month 4 (Jun 8-Jul 7): Package Manager 완성
├─ Version management
├─ Dependency locking
├─ Update/remove 명령
├─ 초기 official packages 배포
└─ Deliverable: Production-ready pkg manager

Month 5 (Jul 8-Aug 7): VSCode + LSP
├─ LSP server 구현
├─ VSCode extension
├─ Basic IDE features
└─ Deliverable: VSCode에서 FL 개발 가능

Month 6 (Aug 8-Sep 7): Advanced Tools
├─ freelang doc (doc generator)
├─ freelang test (test runner)
├─ Profiler (기본)
└─ Deliverable: Complete toolchain

Month 7-10: CI/CD, Docs, Polish
├─ GitHub Actions setup
├─ Release automation
├─ Documentation 완성
├─ Performance 최적화
└─ Deliverable: Production-ready 1.0
```

### **코드 예시**

```fl
// bin/freelang-pkg.fl
struct PackageManager {
  var registry_url: string
  var cache_dir: string
  var manifest: Manifest
}

impl PackageManager {
  fn install(pkg_name: string, version: string?) -> Result<(), string> {
    // 1. Resolve version
    let version = version.unwrap_or(self.resolve_latest(pkg_name)?)

    // 2. Check cache
    if self.is_cached(pkg_name, version) {
      self.load_from_cache(pkg_name, version)?
      return Ok(())
    }

    // 3. Download from registry
    let package = self.registry.download(pkg_name, version)?

    // 4. Extract and save to cache
    self.save_to_cache(pkg_name, version, package)?

    // 5. Update manifest
    self.manifest.add_dependency(pkg_name, version)
    self.manifest.save()?

    Ok(())
  }

  fn resolve_dependencies(pkg: Package) -> Result<Graph<Package>, string> {
    let graph = Graph::new()
    let mut queue = Queue::new()
    queue.push(pkg)

    while let Some(current) = queue.pop() {
      for dep in current.dependencies {
        let resolved = self.resolve(dep)?
        graph.add_edge(current, resolved)
        queue.push(resolved)
      }
    }

    Ok(graph)
  }
}

// src/lsp.fl
class LanguageServer {
  fn on_hover(file: string, line: number, col: number) -> Hover? {
    let symbol = self.find_symbol_at(file, line, col)?
    let type_info = self.type_checker.get_type(symbol)?
    return Hover {
      contents: type_info.to_string()
    }
  }

  fn on_completion(file: string, line: number, col: number) -> Completion[] {
    let prefix = self.get_prefix(file, line, col)
    let candidates = self.collect_candidates(prefix)
    return candidates.map(|c| Completion {
      label: c.name,
      kind: c.kind,
      detail: c.type_info
    })
  }
}
```

### **협력 지점**

```
With Team 1 (Compiler):
├─ Month 3+: Registry에 official packages 배포
├─ Month 5: LSP type info 연동
└─ Month 8+: Performance profiling

With Team 2 (Stdlib):
├─ Month 1: Stdlib docs 자동화
├─ Month 4+: Official packages (json, regex, etc)
└─ Month 8+: Test coverage 측정
```

### **위험과 대비**

```
위험 1: Registry server 복잡도
├─ 확률: 40%
├─ 임팩트: 중간
├─ 대비: Simple version 먼저 (in-memory), DB 나중
└─ 기한: Month 3에서 Month 4로 연장 가능

위험 2: VSCode extension 호환성
├─ 확률: 50%
├─ 임팩트: 낮음
├─ 대비: 기본 기능만 먼저, 고급 기능 나중
└─ 기한: Month 5에서 Month 6으로 연장 가능

위험 3: LSP performance
├─ 확률: 30%
├─ 임팩트: 낮음 (optional optimization)
├─ 대비: Simple LSP 먼저, 나중 최적화
└─ 기한: Month 7-10에 예약
```

---

## 📊 **3명 팀 전체 협력 일정**

### **주간 동기화 (Weekly Sync)**

```
월요일 10:00 UTC: Team standup (30분)
├─ 각자 진행 상황 공유
├─ 블로커 식별
└─ 이주 계획 수립

수요일 15:00 UTC: Cross-team sync (1시간)
├─ 협력 지점 확인 (APIs, interfaces)
├─ 의존성 관리
└─ 이슈 해결

금요일 16:00 UTC: Demo & Planning (1시간)
├─ 완성된 기능 데모
├─ 다음주 계획 리뷰
└─ 위험 식별
```

### **월간 동기화 (Monthly)**

```
Month start (1st week): Planning
├─ 월 목표 설정
├─ 일정 조정
└─ 리소스 할당

Month mid (2nd-3rd week): Progress check
├─ 50% 마일스톤 확인
├─ 위험 재평가
└─ 필요한 조정

Month end (4th week): Review & retrospective
├─ 월간 deliverable 검증
├─ 팀 건강도 체크
├─ 다음 월 계획 확정
└─ Lessons learned
```

### **의존성 매트릭스**

```
                Team 1          Team 2          Team 3
                (Compiler)      (Stdlib)        (Tools)
────────────────────────────────────────────────
Team 1          -               Type sigs       LSP types

Team 2          Type system     -               Doc format

Team 3          Compiler API    Stdlib APIs     -
```

---

## 💼 **리소스 & 타임라인**

### **월별 비용**

```
Team 1: $35K/month
Team 2: $32K/month
Team 3: $30K/month
─────────────────
Total: $97K/month

10개월 × $97K = $970K (인력비)
+ Infrastructure $5K/month = $50K
+ 기타 (tools, cloud) = $50K
─────────────────────────────
Total: ~$1.07M
```

### **병렬 작업 타임라인**

```
Month 1: 3팀 모두 기초 구현 시작
├─ Team 1: Compiler 설계
├─ Team 2: Collections 구현
└─ Team 3: CLI 기본 구조

Month 2-3: 병렬 진행, 주간 sync
├─ 의존성 최소화 (interfaces 먼저)
├─ 각 팀 독립적 진행
└─ Weekly cross-team sync

Month 3: Bootstrap (Team 1 중심)
├─ Team 2, 3은 각자 진행
├─ Team 1 critical path
└─ 1주 intensive work

Month 4+: 통합 단계
├─ 3팀 협력 증가
├─ API 연동 시작
└─ 통합 테스트
```

---

## 🎯 **팀 역할 요약**

| 역할 | 이름 | 월급 | 책임 | 줄 수 |
|------|------|------|------|-------|
| **CTO** | Team 1 | $35K | Compiler, Type system | 6,500 |
| **Stdlib Lead** | Team 2 | $32K | Collections, Stdlib | 6,000 |
| **DevOps/Tools** | Team 3 | $30K | Package manager, IDE | 9,300 |
| **Total** | 3명 | $97K | 100% 프로젝트 | 21,800 |

---

## 🚨 **팀 건강도 체크리스트**

### **매월 평가 항목**

```
□ 일정 준수율 (목표: > 90%)
□ 블로커 해결 시간 (목표: < 2일)
□ 코드 리뷰 품질 (동료 검증)
□ 팀 분위기 (1-10 점수)
□ 기술 채무 (누적 코드 리팩 시간)
□ 버그율 (per 1000 LOC)
```

### **리스크 모니터링**

```
Team 1: Compiler bootstrap 위험 높음
├─ Mitigation: Backup plan (2-stage bootstrap)
└─ Contingency: +2개월 시간 예약

Team 2: Stdlib 복잡도
├─ Mitigation: Simple 구현 먼저
└─ Contingency: 기능 우선순위 조정

Team 3: Registry/IDE 미지 영역
├─ Mitigation: POC 먼저 (2주)
└─ Contingency: Simple version으로 start
```

---

## 🎬 **팀 온보딩**

### **Week 1 (Onboarding week)**

```
Day 1: Team 미팅 + 환경 설정
Day 2: Codebase 둘러보기 + git workflow
Day 3-4: 각 팀원이 작은 task 해결
Day 5: 첫 PR review + team lunch

Deliverable: 모든 팀원 PR 1개씩 merge
```

### **Week 2-4 (Ramp-up)**

```
Daily 15분 standup
Weekly cross-team sync (1시간)
Pair programming (필요 시)
Code review with detailed feedback
```

---

**팀 구성**: 3명 (CTO, Stdlib Lead, DevOps/Tools)
**총 투입**: $1.07M (10개월)
**목표**: 2026년 12월 31일 v1.0 완성
**성공 확률**: 60-70% (조직이 잘되면)

