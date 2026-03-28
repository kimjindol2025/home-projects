---
name: Sovereign Self-Evolving Code Factory v2 - Team Setup
description: 2026년 최고 난이도 프로젝트 팀 모드 준비 완료 (5 AI Agents)
type: project
---

# 🚀 Sovereign Self-Evolving Code Factory v2 - Team Setup

**상태**: ✅ **Team Setup Complete** (2026-03-18)
**위치**: `.projects/core/sovereign-self-evolving-factory/`
**로컬 Git**: 초기화됨 (2커밋)
**GOGS**: 배포 대기

---

## 📋 프로젝트 개요

**난이도**: ★★★★★ (2026년 AI 연구 최전선)
**비전**: 인간 개입 없이 **코드가 스스로 진화하는 완전 폐쇄 루프**

### 핵심 아이디어

```
입력: "간단한 TODO 앱 만들어줘"
  ↓
5명 AI 에이전트 자동 협력
  ↓
세대 1: 기본 (ProofScore: 0.60)
세대 2: 개선 (ProofScore: 0.75)
세대 3: 최적화 (ProofScore: 0.85)
세대 4: 고도화 (ProofScore: 0.92)
세대 5: 완성 (ProofScore: 0.98) ✅
  ↓
배포 완료
```

---

## 🤖 Team Structure (5 Agents)

| # | 에이전트 | 역할 | 모델 |
|---|---------|------|------|
| 1 | **Intent Architect** | 자연어 → 구조화된 스펙 | Opus 4.6 |
| 2 | **Graph Orchestrator** | 스펙 → FreeWire 그래프 | Sonnet 4.6 |
| 3 | **Code Generator** | 그래프 → fv-lang 코드 | Sonnet 4.6 |
| 4 | **Healing Surgeon** | 테스트 & 버그 자동 수정 | Opus 4.6 |
| 5 | **Evolution Tracker** | 성과 기록 & 다음 세대 피드백 | Haiku 4.5 |

**조율**: Human Oversight (당신)

---

## 📁 프로젝트 구조

```
sovereign-self-evolving-factory/
├── .git/                          # 로컬 Git 저장소
├── .claude/
│   ├── CLAUDE.md                  # 마스터 설계 문서
│   └── agents/
│       ├── intent-architect.md        (100줄)
│       ├── graph-orchestrator.md      (100줄)
│       ├── code-generator.md          (100줄)
│       ├── healing-surgeon.md         (100줄)
│       └── evolution-tracker.md       (100줄)
├── README.md                      # 프로젝트 개요
├── SETUP.md                       # GOGS 배포 가이드
├── MEMORY.md                      # (예정)
├── src/                           # (Phase 0에서 생성)
├── tests/                         # (Phase 0에서 생성)
└── Cargo.toml                     # (Phase 0에서 생성)
```

---

## 📊 로컬 Git 상태

### 커밋 히스토리
```
* 19dc458 (HEAD -> main)
  docs: GOGS 배포 가이드 추가

* cbd4749
  chore: Team setup complete - Phase 0 planning
  - 5 AI Agents 설계 완료
  - 프로젝트 구조 설정
  - Phase 0 상세 계획
  - 팀 모드 준비 완료
```

### 파일 현황
- ✅ CLAUDE.md (마스터 설계, 500+ 줄)
- ✅ README.md (프로젝트 개요)
- ✅ SETUP.md (배포 가이드)
- ✅ 5개 Agent 설계 문서 (각 100줄)
- ⏳ Cargo.toml (Phase 0)
- ⏳ src/ (Phase 0)

---

## 🎯 Phase 0 계획 (1주일)

### Day 1-2: Intent Parser
```
도구: Claude API + Intent Parser
입력: "간단한 TODO 앱 만들어줘"
출력: {
  "name": "todo_app",
  "components": ["List", "Add", "Delete"],
  "features": ["CRUD", "persistence"]
}
```

### Day 3: Graph Builder
```
도구: FreeWire DSL
입력: spec JSON
출력: 시각 그래프 (ASCII/SVG)
```

### Day 4: Code Generator
```
도구: fv-lang compiler
입력: 그래프
출력: fv-lang 코드 (구조체 + 함수)
```

### Day 5: Healing Loop
```
도구: Self-healing agent + ProofScore
입력: 코드
출력: 수정된 코드 (버그 없음)
```

### Day 6: Integration
```
모든 단계 연결
E2E 테스트 (1회 루프 완전 동작)
```

### Day 7: Polish
```
문서화
GOGS PR 자동 생성 준비
배포
```

---

## ✅ Success Criteria (MVP)

| 기준 | 목표 | 현재 |
|------|------|------|
| **Intent 정확도** | 90%+ 파싱 성공 | ⏳ |
| **Code Gen** | fv-lang 컴파일 성공 | ⏳ |
| **Self-Healing** | 자동 버그 수정 | ⏳ |
| **E2E Loop** | 입력 → 배포까지 1회 완성 | ⏳ |
| **GOGS Integration** | PR 자동 생성 | ⏳ |
| **Deterministic** | 같은 입력 → 같은 출력 | ⏳ |

---

## 🔧 기술 스택

| 레이어 | 도구 | 용도 |
|--------|------|------|
| **오케스트레이션** | Claude API | Intent parsing + Agent control |
| **시각화** | FreeWire | Graph generation |
| **코드 생성** | fv-lang | Compiler + codegen |
| **테스트 & 수정** | Self-healing | Bug detection + fix |
| **검증** | ProofScore | Confidence scoring |
| **빌드** | Rust + Cargo | Main implementation |
| **배포** | GOGS | Version control + PR |

---

## 🚀 GOGS 배포 (다음 단계)

### 1단계: 원격 저장소 생성
```
https://gogs.dclub.kr
→ 새 저장소
→ sovereign-self-evolving-factory
→ Public
```

### 2단계: Push
```bash
cd .projects/core/sovereign-self-evolving-factory
git remote add origin https://gogs.dclub.kr/kim/sovereign-self-evolving-factory.git
git push -u origin main
```

### 3단계: 확인
```
https://gogs.dclub.kr/kim/sovereign-self-evolving-factory
```

---

## 💡 핵심 기술 통찰

### 1. Intent Parsing (자연어 이해)
- Claude의 장점: 복잡한 요구사항을 정확히 분석
- 도전: 모호한 입력을 구조화된 스펙으로

### 2. Graph Orchestration (시각화)
- FreeWire의 강점: 시각적 프로그래밍
- 도전: 자동으로 최적 레이아웃 생성

### 3. Code Generation (코드 생성)
- fv-lang의 강점: 안전한 코드 생성 + 컴파일
- 도전: 그래프 → 타입-안전 코드 변환

### 4. Self-Healing (자동 수정)
- 기존 시스템 강점: ProofScore + agent loop
- 도전: 무한 루프 방지 + catastrophic forgetting 방지

### 5. Evolution (학습 & 진화)
- 비전: 세대를 거치면서 자동 개선
- 도전: 진화 방향 설정 + 수렴성 보장

---

## ⚠️ 리소스 제약 (Termux)

- **메모리**: 1GB 이상 사용 금지
- **CPU**: 최대 5개 프로세스 동시
- **시간**: 48시간 연속 실행 최대
- **스토리지**: 100MB 이상 생성 금지

---

## 📈 2026년 AI 트렌드 적합성

✅ **Self-Improving Systems** - MIT/Stanford 2025 최고 난제
✅ **Agentic AI** - OpenAI/Anthropic 최고 우선순위
✅ **Multi-Agent Orchestration** - 업계 표준으로 부상
✅ **Continual Learning** - AGI 연구 핵심
✅ **Code as First-Class Citizen** - Copilot/Claude Code 트렌드

---

## 🎊 최종 상태

```
✅ 설계 완료 (785줄 문서)
✅ Team 구성 완료 (5 Agents)
✅ 로컬 Git 준비 완료 (2커밋)
✅ Phase 0 계획 완료 (1주일)
⏳ GOGS 배포 (다음)
⏳ Phase 0 구현 (Day 1부터)
⏳ 무한 진화 루프 (Phase 2+)
```

---

**생성일**: 2026-03-18
**상태**: ✅ **Team Setup Complete - Ready for Phase 0**
**다음**: GOGS 배포 → Phase 0 구현 (Day 1)

> "AI가 스스로 진화한다" - 역사가 기록할 순간
