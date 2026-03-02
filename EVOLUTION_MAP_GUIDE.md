# 📋 Evolution Map 사용 가이드

**41개 프로젝트, 362,000+ LOC를 4대 축으로 통합 시각화**

---

## 🚀 빠른 시작

### **3개 파일이 생성되었습니다:**

1. **evolution_map.json** (데이터 소스)
   - 41개 프로젝트 메타데이터
   - 개념 연결 그래프
   - 축별 통계

2. **EVOLUTION_MAP.md** (마크다운 문서)
   - 전체 구조 설명
   - 트리맵 형태 시각화
   - 상세 프로젝트 정보

3. **evolution_map_interactive.html** (인터랙티브 웹)
   - D3.js 기반 네트워크 시각화
   - 필터링 / 상세 보기 / 드래그 조작
   - 브라우저에서 열어서 사용

---

## 📊 파일 설명

### **1. evolution_map.json**

```json
{
  "metadata": { ... },
  "core_concepts": { ... },
  "projects": [
    {
      "id": 1,
      "name": "project-name",
      "repo": "https://gogs.dclub.kr/kim/...",
      "status": "production",
      "axes": ["Compiler", "Performance"],
      "concepts": ["Dialect Design", "GPU Acceleration"],
      "loc": 20170,
      "tests": 200,
      "performance": "8.2x improvement"
    }
  ],
  "concept_connections": [ ... ],
  "axis_statistics": { ... },
  "evolution_stages": { ... }
}
```

**용도:**
- 프로그래밍으로 데이터 접근 (JavaScript, Python, etc.)
- 새로운 프로젝트 추가
- 통계 분석

---

### **2. EVOLUTION_MAP.md**

**구성:**
```
📊 Overview (통계)
├─ 4 Major Axes (축 설명)
├─ 🌳 Project Tree (프로젝트 분류)
├─ 🔗 Concept Connection Graph (개념 연결)
├─ 📈 Evolution Stages (진화 단계)
├─ 🎓 Key Statistics (통계표)
├─ 🚀 Top Performance Gains (성능 개선)
├─ 💡 Key Innovations (혁신)
├─ 🔄 Cross-Project Dependencies (의존성)
├─ 📚 Learning Path (학습 경로)
└─ 🎯 Future Connections (향후 계획)
```

**용도:**
- 문서로 읽기 (GitHub, 마크다운 리더)
- 프로젝트 개요 파악
- 축별 프로젝트 탐색

---

### **3. evolution_map_interactive.html**

**기능:**
- ✅ D3.js 네트워크 시각화
- ✅ 축별 필터링 (4개 축 선택)
- ✅ 프로젝트 상세 정보 패널
- ✅ 드래그 가능한 노드
- ✅ 줌 / 패닝 조작
- ✅ 실시간 통계

**사용법:**
```bash
# 1. 파일을 브라우저로 열기
open evolution_map_interactive.html

# 또는

# 2. 간단한 서버로 열기
python -m http.server 8000
# http://localhost:8000/evolution_map_interactive.html
```

---

## 🎯 4대 축 설명

### **Axis 1: Language / DSL Experiment** 🔴
**목표:** 새로운 언어 설계 및 구현

| 항목 | 수치 |
|------|------|
| 프로젝트 수 | 15개 |
| 총 LOC | 85,000 |
| 총 테스트 | 650개 |
| 핵심 프로젝트 | FreeLang v4.5 (45K) |

**주요 개념:**
- SSA Form, Type System, Dialect Design
- Concurrency Models, Memory Management
- Standard Library, Security Framework

**학습 경로:**
1. golang_study (9K) → Basics
2. zig-study (9.5K) → Systems
3. zig-lang (4K) → Professional
4. FreeLang v4 (45K) → Production

---

### **Axis 2: Compiler / Optimization** 🟣
**목표:** 컴파일러 인프라 및 최적화

| 항목 | 수치 |
|------|------|
| 프로젝트 수 | 14개 |
| 총 LOC | 92,000 |
| 총 테스트 | 720개 |
| 성능 개선 | 8.2배 (ARIA) |

**주요 개념:**
- Polyhedral Model, Loop Tiling
- SMT Encoding, Formal Verification
- Pass Framework, Code Generation

**성능 순위:**
1. ARIA: 8.2x (ResNet-50)
2. AIAccel: 4-5x (일반)
3. Query Optimizer: 100-1000x (복합 쿼리)

---

### **Axis 3: Knowledge Engine / Analysis** 🔵
**목표:** 다중 저장소 분석 및 설계 조언

| 항목 | 수치 |
|------|------|
| 프로젝트 수 | 6개 |
| 총 LOC | 42,000 |
| 총 테스트 | 350개 |
| 핵심 프로젝트 | GOGS v8.0 (21.3K) |

**주요 개념:**
- Concept Transfer Detection
- Ecosystem Health Analysis
- Evolution Prediction
- Production Hardening

**GOGS Knowledge Engine 버전:**
- v1.2-2.0: Search + Hybrid
- v3.0: Evolution Inference
- v4.0: Design Intent + Architecture
- v5.0: Design Cognition
- v6.0: Multi-Repo + Ecosystem
- v7.0: Active Design Engine
- v8.0: Production Hardening ⭐

---

### **Axis 4: Performance / Distributed Systems** 🟢
**목표:** 분산 시스템 및 성능 최적화

| 항목 | 수치 |
|------|------|
| 프로젝트 수 | 16개 |
| 총 LOC | 78,000 |
| 총 테스트 | 580개 |
| 핵심 프로젝트 | GRIE (8.5K) |

**주요 개념:**
- Raft Consensus, Circuit Breaker
- GPU/TPU Acceleration
- Transaction Management
- Stream Processing

---

## 🔗 개념 연결 맵

```
┌─────────────────────────────────────────┐
│ 핵심 개념 통합 네트워크                  │
└─────────────────────────────────────────┘

[SSA Form]
    ↓
[MLIR Syntax] ──> [Dialect Design] ──> [Pass Optimization]
    ↓                                      ↓
[Operation]                        [Code Generation]
    ↓                                      ↓
[Polyhedral Model] ──> [Loop Tiling] ──> [GPU Target]
    ↓                       ↓
[Memory Hierarchy]    [Memory Arena]
    ↓                       ↓
[Bufferization] ────> [Performance]
                            ↓
                    [Benchmark Suite]
                            ↓
[PoW/UTXO/Wallet] ──────────┘
    ↓
[Transaction Pool]
    ↓
[Concurrency]
    ↓
[Work Stealing]
    ↓
[Distributed Systems]
    ↓
[Raft Consensus] ──> [Byzantine Resilience]
    ↓
[Active Design Engine]
    ↓
[Multi-Repo Analyzer] ──> [Concept Transfer] ──> [Prediction]
    ↓
[GOGS Knowledge Engine]
```

---

## 📈 진화 단계

### **Stage 0: Research** 🔬
- mlir-postdoc-nextgen (Formal Semantics)
- phd-v16-byzantine-fl (Byzantine FL)

**특징:**
- 아이디어 검증 중
- 논문 작성 단계
- 24개월 장기 프로젝트

### **Stage 1: Learning** 📚
- mlir-study, golang_study, zig-study
- mlir-cpp-learning

**특징:**
- 교육 자료
- 개념 학습
- 문제 풀이

### **Stage 2: Prototype** 🔄
- stochastic-code-optimizer
- freelang-v6 (next-gen design)
- zig-lang (professional course)

**특징:**
- 초기 구현
- 아키텍처 설계
- 성능 측정

### **Stage 3: Development** 🚀
- AIAccel (variant research)
- mlir-adaptive-pipeline (Phase 2.x)

**특징:**
- 적극적 개발
- 기능 추가
- 최적화 진행

### **Stage 4: Production** ✅
- freelang-v4-integrated (v4.5)
- gogs-knowledge-engine (v8.0)
- raft-consensus-engine (Phase 3)
- kimsearch-mini (3 versions)

**특징:**
- 배포 준비 완료
- 운영 가능
- 완전한 문서

### **Stage 5: Mature** ⭐
- freelang-v4 (all 20+ modules)
- mlir-postdoc (ARIA, PhD)
- zlang (LLVM, Phase 3)
- AIAccel (Post-Doc)

**특징:**
- 완전 안정화
- 장기 운영 중
- 확장 가능

---

## 🎓 권장 학습 경로

### **초급 (1-2주)**
```
Week 1: golang_study (9K LOC)
  → Basics: variables, functions, structs, interfaces

Week 2: zig-study (9.5K LOC)
  → Systems programming: memory, comptime, unsafe
```

### **중급 (2-3주)**
```
Week 3: mlir-cpp-learning (12K LOC)
  → MLIR API, Dialect, Pass framework

Week 4: stochastic-code-optimizer (3.5K LOC)
  → Optimization techniques, evolutionary search

Week 5: zig-lang (4K LOC)
  → Professional-level systems programming
```

### **상급 (4-6주)**
```
Week 6-7: ARIA/AIAccel (25K LOC)
  → Production compiler architecture
  → GPU/TPU code generation
  → Performance optimization

Week 8-10: FreeLang v4 (45K LOC)
  → Complete language implementation
  → Type system, runtime, standard library

Week 11-12: GOGS Knowledge Engine (21K LOC)
  → Multi-repo analysis
  → ML + architecture
  → Design intent extraction
```

### **연구 (8주+)**
```
Quarter 1-2: mlir-adaptive-pipeline (45.9K LOC)
  → Cutting-edge optimization
  → Formal verification
  → Polyhedral model

Quarter 3-4: mlir-postdoc-nextgen (2.1K LOC starter)
  → Formal semantics
  → SMT solver integration
  → Research paper writing
```

---

## 🚀 새 프로젝트 추가하기

### **Step 1: JSON에 추가**

```json
{
  "id": 42,
  "name": "my-new-project",
  "repo": "https://gogs.dclub.kr/kim/my-new-project.git",
  "status": "prototype",
  "axes": ["Compiler", "Performance"],
  "concepts": ["Custom Concept", "New Idea"],
  "loc": 5000,
  "tests": 50,
  "performance": "Initial measurements",
  "stage": "Development"
}
```

### **Step 2: 개념 연결**

```json
{
  "from": "Dialect Design",
  "to": "My New Concept",
  "type": "extension",
  "description": "Extends dialectdesign..."
}
```

### **Step 3: 재생성**

```bash
# HTML 시각화 업데이트 (브라우저 새로고침)
# Markdown 문서 재구성
```

---

## 💾 데이터 업데이트

### **프로젝트 상태 업데이트**

```bash
# 1. evolution_map.json 수정
vi evolution_map.json

# 2. 프로젝트 정보 변경
# - loc (줄 수)
# - tests (테스트 수)
# - status (상태)
# - stage (단계)

# 3. 브라우저 새로고침
# evolution_map_interactive.html
```

### **개념 추가**

```json
"core_concepts": {
  "Language": [..., "New Concept"],
  "Compiler": [..., "New Idea"],
  "Knowledge": [...],
  "Performance": [...]
}
```

---

## 📊 통계 분석

### **축별 분포**

```
Axis 1: Language (15/41 = 37%)
Axis 2: Compiler (14/41 = 34%)
Axis 3: Knowledge (6/41 = 15%)
Axis 4: Performance (16/41 = 39%)
```

*참고: 다중 축 프로젝트는 중복 계산됨*

### **LOC 분포**

```
Language: 85,000 (23%)
Compiler: 92,000 (25%)
Knowledge: 42,000 (12%)
Performance: 78,000 (22%)
Other: 65,000 (18%)
────────────────────
Total: 362,000 LOC
```

### **테스트 커버리지**

```
Average: 85 tests/project
Median: 48 tests
Range: 0-500 tests
Total: 3,500+ tests
```

---

## 🔍 검색 및 필터링

### **Interactive HTML에서:**
- 왼쪽 사이드바에서 축 선택
- 프로젝트명 클릭 → 상세 정보
- 드래그로 노드 이동
- 줌 / 패닝으로 조작

### **Markdown에서:**
- `Ctrl+F` (또는 `Cmd+F`) → 프로젝트명 검색
- 각 축별 섹션으로 탐색

### **JSON으로:**
```python
import json

with open('evolution_map.json') as f:
    data = json.load(f)

# 특정 축 필터링
compiler_projects = [p for p in data['projects']
                     if 'Compiler' in p['axes']]

# LOC 합계
total_loc = sum(p['loc'] for p in data['projects'])
```

---

## 🔐 버전 관리

### **Current Version: 1.0** (2026-02-27)

**포함된 프로젝트:**
- 41개 프로젝트
- 4개 축
- 8개 개념 연결

**다음 버전 (계획):**
- v1.1: 50개 프로젝트 추가
- v1.2: 시각화 개선 (3D 네트워크)
- v2.0: 자동 업데이트 스크립트

---

## 🆘 문제 해결

### **HTML이 로드되지 않음**

```bash
# CORS 문제 해결
python -m http.server 8000
# http://localhost:8000/evolution_map_interactive.html
```

### **JSON 데이터가 오래됨**

```bash
# 최신 메타데이터로 재생성
# evolution_map.json 수동 편집
```

### **새로운 프로젝트가 보이지 않음**

```bash
# 1. JSON에 항목 추가 확인
# 2. 브라우저 캐시 삭제 (Ctrl+Shift+Delete)
# 3. HTML 새로고침 (Ctrl+F5)
```

---

## 📚 참고 자료

- **메모리**: `/data/.../memory/MEMORY.md`
- **프로젝트 상세**: 각 `https://gogs.dclub.kr/kim/` 저장소
- **논문/보고서**: 각 프로젝트의 `README.md` / `THESIS.md`

---

## 🎯 다음 단계

1. **interactive.html 열기** (브라우저)
   → 네트워크 시각화 탐색

2. **EVOLUTION_MAP.md 읽기** (마크다운)
   → 전체 구조 이해

3. **evolution_map.json 분석** (프로그래밍)
   → 자동화 / 커스터마이징

4. **새 프로젝트 추가**
   → 자신의 아이디어 등록

---

**마지막 업데이트**: 2026-02-27 ✨

**질문이 있으신가요?**
- 🔍 검색: JSON에서 프로젝트 찾기
- 📖 문서: Markdown 상세 설명
- 🎨 시각화: HTML 네트워크 탐색
