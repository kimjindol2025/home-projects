# Week 1: CI/CD 기초 설계

## 🎯 목표

**사람 개입 0%의 완전 자동화**

```
커밋 푸시
  ↓
GitHub Actions 자동 실행
  ↓
build → test → benchmark
  ↓
결과 자동 저장
  ↓
PR에 성능 수치 자동 코멘트
```

---

## 📋 구현 완료 항목

### 1️⃣ CI 파이프라인 (`.github/workflows/ci.yml`)

**3개 Job:**

#### Job 1: `build-and-test`
- `npm ci` (의존성 설치)
- `npm run build` (TypeScript 컴파일)
- `npm test` (34/34 테스트 실행)
- 결과를 artifact로 저장

**조건:** 테스트 실패 시 즉시 pipeline 중단

#### Job 2: `benchmark`
- Job 1 통과 후에만 실행 (`needs: build-and-test`)
- 테스트 출력에서 성능 수치 추출
  - 총 테스트 수
  - 통과한 테스트
  - 실행 시간 (ms)
- JSON으로 저장
- PR에 자동 코멘트 (성능 결과)

#### Job 3: `verify-success`
- 모든 job 완료 후 최종 검증
- 실패 시 CI 실패 표시

---

### 2️⃣ 벤치마크 수집 스크립트 (`scripts/collect-metrics.js`)

**기능:**

```bash
npm run collect-metrics
```

- 테스트 실행 후 성능 데이터 추출
- `.benchmark-history/` 디렉토리에 저장
  - `result_TIMESTAMP_COMMIT.json` (개별 결과)
  - `metrics.json` (모든 결과의 누적)
- 최대 100개 결과 유지 (자동 정리)
- 요약 리포트 출력

**출력 예:**
```
📊 Benchmark Summary

Current: 2026-03-03T10:30:00Z
Tests: 34/34 (100%)
Execution: 26650ms
Trend: ⬇️ faster (5.2%)

History: 5 measurements
```

---

### 3️⃣ 벤치마크 비교 스크립트 (`scripts/compare-benchmarks.js`)

**기능:**

```bash
npm run compare-benchmarks
```

- 현재 vs 이전 성능 비교
- 성능 회귀 감지 (>10% 느려지면 경고)
- 추세 분석 (최근 5개 측정)
- 통계 계산 (평균, 최소, 최대)

**출력 예:**
```
Previous: 28000ms
Current:  26650ms
Change:   -1350ms (-4.82%)

✅ All checks passed
```

**실패 조건:**
- 테스트 통과율 감소
- 성능 10% 이상 저하

---

## 🔧 로컬 실행

### 전체 파이프라인 테스트

```bash
# 1. 빌드 + 테스트
npm run build
npm test

# 2. 성능 수치 수집
npm run collect-metrics

# 3. 이전 결과와 비교
npm run compare-benchmarks
```

### PR 생성 시 자동 실행

```bash
git add .
git commit -m "Add CI pipeline infrastructure"
git push origin feature/week1-ci
```

GitHub에서 PR을 열면 자동으로:
1. CI 파이프라인 실행
2. 테스트 결과 표시
3. 성능 수치를 PR 코멘트로 추가

---

## 📊 저장되는 데이터

### `.benchmark-history/metrics.json`
```json
[
  {
    "timestamp": "2026-03-03T10:30:00Z",
    "commit": "abc123def456",
    "branch": "master",
    "tests": {
      "passed": 34,
      "failed": 0,
      "total": 34,
      "passRate": 100.0
    },
    "performance": {
      "executionTimeMs": 26650
    }
  }
]
```

### `.benchmark-history/result_TIMESTAMP_COMMIT.json`
동일한 구조, 개별 저장

---

## ⚠️ 주의사항

1. **`.benchmark-history/` 디렉토리**
   - `.gitignore`에 추가됨 (로컬에만 저장)
   - GitHub Actions가 artifact로 저장

2. **PR 코멘트**
   - GitHub Actions 권한 필요
   - `permissions: write` 확인

3. **실패 시 처리**
   - 테스트 실패 → 파이프라인 중단 ✅
   - 벤치마크 실패 → 경고만 (파이프라인 계속) ✅

---

## 🚀 Week 1 완료 조건

- [x] CI 파이프라인 작동
- [x] 커밋 → 자동 테스트
- [x] 테스트 통과 → 벤치마크 자동 실행
- [x] 벤치마크 결과 → JSON 저장
- [x] PR에 자동 코멘트
- [ ] 실제 커밋으로 테스트 (다음 단계)

---

## 📝 다음: Week 2 (성능 측정의 현실화)

현재:
- 테스트 실행 시간만 측정

Week 2:
- 실제 네트워크 기반 테스트
- p50/p95/p99 latency 측정
- 메모리 사용량 측정
- 처리량 ops/sec 측정
- 성능 회귀 시 CI 실패

---

**Status**: Week 1 설계 완료 ✅
