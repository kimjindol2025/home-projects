# 🔥 Mojo & 🎉 Julia 언어 학습 프로젝트

## 📊 Mojo 학습 (완료) ✅

**저장소**: `gogs.dclub.kr/kim/mojo-learning`
**상태**: ✅ **Phase 1-10 완료 + 컴파일러 구현 중**

### 구조
```
step01-setup/       → 환경 설정, Hello World (4 파일)
step02-basics/      → 변수, 타입, 제어문 (6 파일)
step03-types/       → 타입 시스템 심화 (5 파일)
step04-collections/ → 배열, 딕셔너리, SIMD (5 파일)
step05-ownership/   → 소유권, 참조 차용 (3 파일)
step06-functions/   → 고급 함수, 재귀 (2 파일)
step07-structs/     → 구조체, 메서드 (2 파일)
step08-traits/      → 다형성, Trait (1 파일)
step09-performance/ → SIMD, 병렬 처리 (1 파일)
step10-aiml/        → 신경망, 행렬 연산 (3 파일)
compiler-impl/      → Phase 2: 들여쓰기 파서 (25+ 파일)
```

### 주요 성과
- 총 28개 .mojo 파일 + 10개 NOTES.md
- 4,500+ 줄 코드
- 2일 학습 완료 (2026-03-10~11)
- AI/ML 신경망 완전 구현

### 핵심 배운 점
- **fn vs def**: fn은 컴파일 타임 최적화
- **SIMD**: 벡터 연산으로 4배 성능 향상
- **Ownership**: Rust처럼 소유권 관리
- **Multiple Dispatch** 없음 (오버로딩으로 대체)

---

## 🎉 Julia 학습 (준비 중) 🔄

**목표**: Mojo와 동일한 10단계 구조로 Julia 완전 학습

### 차별점: Julia vs Mojo

| 항목 | Mojo | Julia |
|------|------|-------|
| **기반** | Python 문법 | 수학 기반 (MATLAB 유사) |
| **성능** | C++ 수준 | Fortran 수준 |
| **학습곡선** | 쉬움 | 중간 |
| **Multiple Dispatch** | ❌ | ✅ (핵심 기능) |
| **병렬화** | SIMD | @simd, @threads |
| **용도** | AI/ML | 과학계산 |
| **패키지** | 적음 | 많음 (Plots, Flux 등) |

### 준비할 내용

**10단계 로드맵**:
1. 환경 설정 (REPL, Jupyter, IDE)
2. 기본 문법 (변수, 함수, 타입)
3. Multiple Dispatch (Julia의 핵심)
4. 배열과 선형대수
5. 제어 흐름과 고차 함수
6. 모듈과 패키지 시스템
7. 메타프로그래밍과 매크로
8. 성능 최적화 기법
9. 병렬 처리 (@threads, pmap)
10. 과학계산 프로젝트 (선형대수, 미분)

**예상 결과**:
- 32개 .jl 파일 (Mojo보다 10% 많음)
- 10개 NOTES.md
- 5,000+ 줄 코드
- 3일 학습 예상

### 시작 전 준비
- [ ] Julia 다운로드 및 설치 가이드
- [ ] Jupyter + VSCode 환경 설정
- [ ] 패키지 관리 (Pkg.jl) 이해
- [ ] REPL 단축키 학습

---

## 🔗 연관 프로젝트

### GoFree (Go-based FreeLang Compiler)
- **Phase 0-12 완료** (고급 최적화)
- ~4,750줄 코드
- 115개 테스트 (모두 통과)

### FreeLang v2 (Self-Parse)
- **Phase 3 준비 완료**
- 실행 환경만 필요

### 학습 패턴
> Mojo → Julia → 언어 설계로 확대 (프로그래밍 언어 이해도 깊이화)

---

**마지막 업데이트**: 2026-03-11 08:15 UTC+9
**상태**: Mojo 완료, Julia 준비 단계
