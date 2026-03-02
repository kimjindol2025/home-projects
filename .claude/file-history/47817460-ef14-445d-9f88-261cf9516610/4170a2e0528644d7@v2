# GC 2부: 자율적 메모리 해방 (Autonomous Memory Liberation)

**Kim님의 "언어 독립" 대망의 제2장**

---

## 🎯 미션

FreeLang이 외부 런타임(Node.js, Rust) 없이 **독립적으로 메모리를 관리**하는 시스템을 구축합니다.

- **1부 (Arc)**: 구조 = 건물의 기초 ✅ 완료
- **2부 (GC)**: 지능 = 건물의 신경계 🚀 NOW

---

## 📊 프로젝트 개요

| 항목 | 내용 |
|------|------|
| **총 코드** | 6,100줄 |
| **테스트** | 80개 |
| **기간** | 4주 |
| **목표** | "어른의 언어" = 메모리 자율 관리 |

---

## 🏗️ 4주 마일스톤

### Week 1: Tracing GC 기초 (1,500줄)
- `generational_gc.fl`: Young/Old Generation
- `mark_and_sweep.fl`: Tricolor Marking + 순환 참조 해결
- `test_memory_stability.fl`: 메모리 안정성 테스트

### Week 2: Compaction & Memory Layout (1,400줄)
- `memory_compactor.fl`: LISP2 Algorithm
- `heap_manager.fl`: Bump Pointer + Card Marking
- `test_compaction.fl`: 단편화 검증

### Week 3: Concurrent GC & Safepoint (1,600줄)
- `concurrent_gc.fl`: Tri-color + Write Barriers
- `safepoint_handler.fl`: Thread-safe 조정
- `test_concurrent.fl`: 동시성 검증

### Week 4: Optimization & Integration (1,600줄)
- `gc_optimizer.fl`: Adaptive Tuning
- `metrics_collector.fl`: 성능 수집
- `integration_tests.fl`: Raft DB + GC 협력

---

## 📌 핵심 개념

### Generational Hypothesis
> "대부분의 객체는 생성 후 곧 소멸한다"

**결과**: Young Generation을 자주 검사 → STW 시간 감소

### Tricolor Marking
```
White: 아직 탐색 안 함 (쓰레기 후보)
Gray: 탐색 중 (자식 검사 예정)
Black: 완전히 탐색함 (살아있음)
```

### Write Barrier
Old Generation → Young Generation 참조 추적으로 전체 힙 검사 회피

---

## 🧪 메모리 안정성 테스트 스위트

**Week 1 테스트 (400줄, 20개 테스트)**

1. **기본 할당/해제** (4개)
   - 단일 객체 할당/해제
   - 배열 할당/해제
   - 구조체 할당/해제
   - 중첩 구조 할당/해제

2. **순환 참조 해결** (4개)
   - 2-cycle (A→B→A)
   - 3-cycle (A→B→C→A)
   - Self-cycle (A→A)
   - 복합 cycle + 외부 참조

3. **Generation 전환** (4개)
   - Young → Old 승격
   - Old 내 참조 업데이트
   - Young 전체 회수
   - Old 전체 회수 (Full GC)

4. **메모리 단편화** (4개)
   - 할당/해제 패턴 1 (sequential)
   - 할당/해제 패턴 2 (interleaved)
   - 할당/해제 패턴 3 (fragmented)
   - Compaction 검증 (< 5% 단편화)

5. **동시성 기본** (4개)
   - 2 스레드 동시 할당
   - 2 스레드 동시 해제
   - 1 할당 + 1 GC 동시
   - Race condition 없음 검증

**성공 기준**:
- ✅ 모든 테스트 통과
- ✅ 메모리 누수 0%
- ✅ 순환 참조 해결 100%
- ✅ 단편화 < 5%

---

## 🚀 시작

```bash
# Week 1 구현
cd freelang-gc-part2/

# 메모리 구조 설계
cat src/generational_gc.fl

# 메모리 안정성 테스트
cat tests/test_memory_stability.fl

# 테스트 실행 (FreeLang Runtime 필요)
freelang tests/test_memory_stability.fl
```

---

## 📚 설계 문서

- `GC_DESIGN.md`: 상세 설계 (이론 + 코드)
- `MEMORY_SAFETY.md`: 메모리 안정성 검증 방법론
- `PERFORMANCE.md`: GC 성능 최적화 가이드

---

## 💡 "기록이 증명이다"

### 1부 (Arc) 증명
```
증명: Rc/Arc로 대부분의 메모리 관리 가능
한계: 순환 참조 남음
```

### 2부 (GC) 증명
```
증명: Mark-Sweep + Compaction으로 완전 자동화
성과: 순환 참조 해결 ✅ + 단편화 제거 ✅
```

### 최종 증명 (언어 독립)
```
1부 + 2부 = "FreeLang이 자체 메모리를 완전히 관리"
= 외부 런타임 없이 독립 실행
```

---

**철학**: "언어의 자아 = 메모리 정화의 지능"

**목표**: FreeLang을 "어른의 언어"로 진화시키기

---

**마지막 업데이트**: 2026-03-02
**상태**: Week 1 실행 중 🚀
