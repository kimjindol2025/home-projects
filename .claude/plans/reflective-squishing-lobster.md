# v6.2 플랜: 단위 테스트 & TDD (Unit Test & TDD)
# 대학교 4학년 전공 심화

## Context
Python University v6.1 디자인 패턴 완성 후 다음 단계.
v6.1에서 만든 클래스들을 대상으로 체계적인 단위 테스트와 TDD를 실습한다.

- 대상 파일: `university_v6_1_DESIGN_PATTERNS.py`
- 새 파일: `university_v6_2_UNIT_TEST_TDD.py`
- 위치: `/data/data/com.termux/files/home/python-gogs/`
- 수정 파일: v6.1의 `ExperimentBuilder.set_duration()`에 음수 검증 추가 (TDD Green 단계)

핵심 학습:
- `unittest` 프레임워크 기반 TestCase 작성
- TDD 사이클: Red(실패) → Green(구현) → Refactor(정리)
- `MagicMock`으로 의존성 격리 (Mock & Stub)
- `setUp()`으로 싱글톤 초기화 → 테스트 독립성 보장

---

## 구현 구조 (5개 Part + 테스트 12개)

### Part 1: unittest 기초 — 싱글톤 검증 (TestGogsSystemLogger)
```python
class TestGogsSystemLogger(unittest.TestCase):
    def setUp(self):
        GogsSystemLogger._instance = None  # 싱글톤 초기화로 테스트 독립성 보장

    def test_singleton_identity(self):       # 1: assertIs(a, b)
    def test_log_recorded(self):             # 2: assertEqual(count, 1)
    def test_clear_resets_logs(self):        # 3: assertEqual(count, 0)
```

### Part 2: Mock & Stub — 팩토리 격리 (TestSensorFactory)
```python
class TestSensorFactory(unittest.TestCase):
    def test_get_sensor_returns_temperature(self):   # 4: assertIsInstance
    def test_get_sensor_invalid_raises(self):         # 5: assertRaises(ValueError)
    def test_register_custom_sensor(self):            # 6: 런타임 등록 검증
    def test_sensor_collect_with_mock(self):          # 7: MagicMock.assert_called_once()
```

### Part 3: TDD 사이클 — Red→Green→Refactor (TestExperimentBuilder)
```python
class TestExperimentBuilder(unittest.TestCase):
    def test_empty_name_raises(self):       # 8: build() without name → ValueError
    def test_negative_duration_raises(self): # 9: set_duration(-1) → ValueError (TDD!)
```
TDD Green 단계: `ExperimentBuilder.set_duration()`에 `if hours < 0: raise ValueError` 추가

### Part 4: 통합 테스트 (TestDatabaseConnectionPool + TestGogsResearchSystem)
```python
class TestDatabaseConnectionPool(unittest.TestCase):
    def setUp(self): DatabaseConnectionPool._instance = None

    def test_pool_has_max_connections(self):    # 10: stats['total'] == MAX
    def test_exhausted_pool_raises(self):        # 11: RuntimeError

class TestGogsResearchSystem(unittest.TestCase):
    def test_run_experiment_returns_results(self): # 12: len(results) == 센서 수
```

### Part 5: TDD 사이클 보고서 + 커버리지 시뮬레이션
```python
def print_tdd_report():
    """Red→Green→Refactor 각 단계 코드와 결과 출력"""

def print_coverage_report(suite):
    """테스트 수 / 커버된 클래스 수 등 간단한 커버리지 집계"""
```

---

## v6.1 파일 수정 (TDD Green 단계 실증)

파일: `university_v6_1_DESIGN_PATTERNS.py`
대상 메서드: `ExperimentBuilder.set_duration()` (line ~295)

```python
# 수정 전
def set_duration(self, hours: int) -> 'ExperimentBuilder':
    self._experiment['duration_hours'] = hours
    return self

# 수정 후 (Green)
def set_duration(self, hours: int) -> 'ExperimentBuilder':
    if hours < 0:
        raise ValueError("실험 지속 시간은 0 이상이어야 합니다")
    self._experiment['duration_hours'] = hours
    return self
```

---

## 테스트 계획 (12개)

| # | 대상 | 테스트 내용 | assert 메서드 |
|---|------|-----------|--------------|
| 1 | Singleton | 동일 인스턴스 보장 | assertIs |
| 2 | Singleton | 로그 기록 확인 | assertEqual |
| 3 | Singleton | clear() 초기화 | assertEqual |
| 4 | Factory | TemperatureSensor 생성 | assertIsInstance |
| 5 | Factory | 잘못된 타입 → ValueError | assertRaises |
| 6 | Factory | register() 런타임 등록 | assertIsInstance |
| 7 | Factory | Mock으로 collect() 검증 | assert_called_once |
| 8 | TDD(Red) | 이름 없음 → ValueError | assertRaises |
| 9 | TDD(Green)| 음수 시간 → ValueError | assertRaises |
| 10 | DB Pool | MAX 연결 수 검증 | assertEqual |
| 11 | DB Pool | 소진 시 RuntimeError | assertRaises |
| 12 | System | 통합 실험 실행 결과 | assertEqual |

---

## 파일 구조 예상

```
 80줄: Part 1 - unittest 기초 + 싱글톤 테스트
 90줄: Part 2 - Mock/Stub + 팩토리 테스트
 80줄: Part 3 - TDD Red→Green→Refactor 시각화
 90줄: Part 4 - DB Pool + System 통합 테스트
 80줄: Part 5 - TDD 보고서 + 커버리지
 80줄: Main, 결과 출력, 철학
---
500줄 예상
```

---

## 검증

```bash
# TDD Green 단계 먼저 v6.1 수정
# 이후 v6.2 실행
python3 university_v6_2_UNIT_TEST_TDD.py test
# → 12/12 PASS

git add university_v6_2_UNIT_TEST_TDD.py university_v6_1_DESIGN_PATTERNS.py
git commit -m "v6.2: Unit Test & TDD — 코드의 무결성 증명"
```

## 핵심 학습 목표

1. **unittest**: TestCase, setUp, assert* 메서드 완전 활용
2. **TDD 사이클**: Red → Green → Refactor 실제 코드로 증명
3. **Mock/MagicMock**: 의존성 격리로 순수한 단위 테스트
4. **테스트 독립성**: setUp()으로 싱글톤 상태 리셋
5. **커버리지 개념**: 어느 코드 경로까지 테스트했는가

철학: 기록이 증명이다 — 테스트가 코드의 진실을 증명한다 gogs
