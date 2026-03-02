# C11 라이브러리 저장소 분석 템플릿

이 템플릿을 각 저장소마다 복사하여 사용하세요.

---

## 저장소명: `${REPO_NAME}`

**분석 날짜**: `${ANALYSIS_DATE}`
**분석자**: `${ANALYST_NAME}`

---

## 📌 1. 기본 정보

### 1.1 저장소 메타데이터

| 항목 | 값 |
|------|-----|
| **저장소명** | `${REPO_NAME}` |
| **Git URL** | `https://gogs.dclub.kr/kim/${REPO_NAME}` |
| **설명** | `${DESCRIPTION}` |
| **생성 날짜** | `${CREATED_AT}` |
| **마지막 푸시** | `${PUSHED_AT}` |
| **저장소 크기** | `${SIZE}` MB |
| **Stars** | `${STARS}` |
| **Watchers** | `${WATCHERS}` |
| **공개 여부** | PUBLIC / PRIVATE |

### 1.2 분류

- **그룹**: Tier 1 / Tier 2 / Tier 3
- **카테고리**: 유틸리티 / 고급기능 / 프레임워크
- **타입**: Pure C11 / C11 + 의존성

---

## 📂 2. 파일 구조 분석

### 2.1 디렉토리 맵

```
${REPO_NAME}/
├── README.md              ☐ ✅ ☐ ❌
├── LICENSE                ☐ ✅ ☐ ❌  (라이센스: ${LICENSE_TYPE})
├── Makefile               ☐ ✅ ☐ ❌
├── CMakeLists.txt         ☐ ✅ ☐ ❌
├── .gitignore             ☐ ✅ ☐ ❌
├── src/
│   ├── (파일 목록)
│   └── 총 라인: ${SRC_LINES}
├── include/
│   ├── (파일 목록)
│   └── 총 라인: ${INCLUDE_LINES}
├── tests/
│   ├── (파일 목록)
│   └── 총 라인: ${TEST_LINES}
├── examples/              ☐ ✅ ☐ ❌
│   └── (파일 개수: ${EXAMPLES_COUNT})
├── docs/                  ☐ ✅ ☐ ❌
│   └── (문서 파일: ${DOCS_FILES})
└── CI/CD 구성            ☐ ✅ ☐ ❌
    └── (.github/workflows, .gitea 등)
```

### 2.2 파일 구조 평가

```
파일 구조 표준 준수도:

☐ 완벽 (100%): 모든 표준 디렉토리 포함
☐ 우수 (80%): 대부분의 표준 디렉토리 포함
☐ 양호 (60%): 기본 디렉토리만 포함
☐ 부족 (40%): 일부 디렉토리만 포함
☐ 매우부족 (20%): 최소한의 구조만
```

### 2.3 파일 목록 (상세)

#### src/ 디렉토리

| 파일명 | 라인 | 함수 | 설명 |
|--------|------|------|------|
| `${FILE1}` | `${LINES1}` | `${FUNCS1}` | ${DESC1} |
| `${FILE2}` | `${LINES2}` | `${FUNCS2}` | ${DESC2} |
| ... | ... | ... | ... |

#### include/ 디렉토리

| 파일명 | 라인 | 구조체 | 설명 |
|--------|------|--------|------|
| `${HEADER1}` | `${LINES1}` | `${STRUCTS1}` | ${DESC1} |
| `${HEADER2}` | `${LINES2}` | `${STRUCTS2}` | ${DESC2} |
| ... | ... | ... | ... |

#### tests/ 디렉토리

| 파일명 | 라인 | 테스트 수 | 커버리지 |
|--------|------|----------|---------|
| `${TEST1}` | `${LINES1}` | `${COUNT1}` | ${COVERAGE1}% |
| `${TEST2}` | `${LINES2}` | `${COUNT2}` | ${COVERAGE2}% |
| ... | ... | ... | ... |

---

## 📊 3. 코드 메트릭

### 3.1 종합 통계

| 메트릭 | 값 | 비고 |
|--------|-----|------|
| **C 소스 파일** | `${C_FILES}` 개 | src/ 디렉토리 |
| **헤더 파일** | `${H_FILES}` 개 | include/ 디렉토리 |
| **테스트 파일** | `${TEST_FILES}` 개 | tests/ 디렉토리 |
| **총 파일 수** | `${TOTAL_FILES}` 개 | src/ + include/ + tests/ |
| **소스 코드량** | `${SRC_LINES}` 라인 | 주석 제외 |
| **헤더 코드량** | `${INCLUDE_LINES}` 라인 | 주석 제외 |
| **테스트 코드량** | `${TEST_LINES}` 라인 | 주석 제외 |
| **총 코드량** | `${TOTAL_LINES}` 라인 | 전체 합계 |
| **테스트/소스 비율** | `${TEST_RATIO}`% | 이상적: 30-50% |

### 3.2 라이브러리 크기 분류

```
☐ 마이크로 라이브러리    (<500 라인)
☐ 소형 라이브러리       (500-2,000 라인)
☐ 중형 라이브러리       (2,000-5,000 라인)
☐ 대형 라이브러리       (5,000-10,000 라인)
☐ 대규모 라이브러리     (>10,000 라인)

현재 크기: ${TOTAL_LINES} 라인 = [분류]
```

### 3.3 코드 구성 분석

#### 함수 분석

| 메트릭 | 값 | 평가 |
|--------|-----|------|
| **전체 함수 개수** | `${FUNC_COUNT}` | ☐ 많음 ☐ 적당 ☐ 적음 |
| **평균 함수 크기** | `${AVG_FUNC_SIZE}` 라인 | ☐ 우수(<50) ☐ 양호(50-100) ☐ 부족(>100) |
| **최대 함수 크기** | `${MAX_FUNC_SIZE}` 라인 | ☐ 관리가능 ☐ 분해필요 |
| **최소 함수 크기** | `${MIN_FUNC_SIZE}` 라인 | - |

#### 타입 정의

| 타입 | 개수 | 설명 |
|------|------|------|
| **struct** | `${STRUCT_COUNT}` | 데이터 구조 정의 |
| **typedef** | `${TYPEDEF_COUNT}` | 타입 별명 |
| **enum** | `${ENUM_COUNT}` | 열거형 정의 |
| **union** | `${UNION_COUNT}` | 공용체 정의 |

#### 메모리 관리

```
메모리 할당:
  ☐ malloc 사용: ${MALLOC_COUNT} 개 호출
  ☐ calloc 사용: ${CALLOC_COUNT} 개 호출
  ☐ realloc 사용: ${REALLOC_COUNT} 개 호출

메모리 해제:
  ☐ free 사용: ${FREE_COUNT} 개 호출
  ☐ 할당-해제 짝: ${MATCH_RATE}%

메모리 누수 가능성:
  ☐ LOW (<5% 위험)
  ☐ MEDIUM (5-20% 위험)
  ☐ HIGH (>20% 위험)

주요 우려사항:
  1. ${CONCERN1}
  2. ${CONCERN2}
  3. ${CONCERN3}
```

### 3.4 복잡도 분석

| 항목 | 값 | 평가 |
|------|-----|------|
| **순환 복잡도 평균** | `${CYCLOMATIC}` | ☐ 낮음 ☐ 중간 ☐ 높음 |
| **가장 복잡한 함수** | `${COMPLEX_FUNC}` | ${COMPLEX_VALUE} |
| **복합 함수 비율** | `${COMPLEX_RATIO}`% | ☐ 양호 ☐ 개선필요 |

---

## 📚 4. 문서화 평가

### 4.1 README.md 분석

**존재 여부**: ☐ YES ☐ NO

**README 주요 내용**:

```markdown
${README_SUMMARY}
```

**항목별 존재 여부**:

| 항목 | 존재 | 품질 |
|------|------|------|
| 프로젝트 설명 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 주요 기능 목록 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 설치 방법 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 기본 사용 예제 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| API 문서 링크 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 라이센스 정보 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 기여 가이드 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 변경로그 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |
| 저자/라이센스 | ☐ ✅ ☐ ❌ | ☐ 우수 ☐ 양호 ☐ 부족 |

**README 품질 평가**:
- 명확성: ☐ 우수 ☐ 양호 ☐ 개선필요
- 완전성: ☐ 우수 ☐ 양호 ☐ 개선필요
- 실용성: ☐ 우수 ☐ 양호 ☐ 개선필요

**README 점수**: `${README_SCORE}` / 5

### 4.2 API 문서

**문서 위치**: ${API_DOC_LOCATION}

**API 문서 형식**:
- ☐ 헤더 주석 (*.h 파일의 Doxygen 스타일)
- ☐ 별도 API 문서 (docs/API.md)
- ☐ 함수별 상세 설명 (Doxygen, JavaDoc 스타일)
- ☐ 온라인 API 참고서

**API 문서 범위**:

| 항목 | 상태 | 설명 |
|------|------|------|
| 전체 API 함수 설명 | ☐ 100% ☐ 80% ☐ 50% ☐ <50% | ${API_COVERAGE}개 / ${TOTAL_FUNC}개 |
| 함수 매개변수 설명 | ☐ 100% ☐ 80% ☐ 50% ☐ <50% | 각 함수마다 |
| 반환값 설명 | ☐ 100% ☐ 80% ☐ 50% ☐ <50% | 반환값 명확성 |
| 에러 코드 설명 | ☐ 100% ☐ 80% ☐ 50% ☐ <50% | 오류 처리 명확성 |
| 타입 정의 설명 | ☐ 100% ☐ 80% ☐ 50% ☐ <50% | 구조체, typedef 등 |

**API 문서 품질 평가**:
```
☐ 완벽함 (✅): 모든 함수에 상세 설명, 매개변수-반환값-오류 명확
☐ 부분적 (⚠️): 주요 함수만 설명, 일부 누락
☐ 없음 (❌): API 문서 부재 또는 최소한의 설명
```

**API 문서 점수**: `${API_SCORE}` / 5

### 4.3 사용 예제

**예제 위치**: ${EXAMPLES_LOCATION}
**예제 파일 개수**: `${EXAMPLES_COUNT}` 개

**예제 유형**:

| 유형 | 존재 | 파일 | 설명 |
|------|------|------|------|
| **기본 사용 예제** | ☐ YES ☐ NO | ${BASIC_EXAMPLE} | 간단한 시작 예제 |
| **고급 기능 예제** | ☐ YES ☐ NO | ${ADVANCED_EXAMPLE} | 고급 기능 활용 |
| **에러 처리 예제** | ☐ YES ☐ NO | ${ERROR_EXAMPLE} | 에러 처리 방법 |
| **성능 예제** | ☐ YES ☐ NO | ${PERF_EXAMPLE} | 성능 최적화 팁 |
| **통합 테스트 예제** | ☐ YES ☐ NO | ${INTEGRATION_EXAMPLE} | 실제 사용 시나리오 |

**예제 품질**:
- 실행 가능: ☐ YES ☐ NO
- 상세 주석: ☐ YES ☐ NO
- 다양한 use case: ☐ YES ☐ NO
- 현재성: ☐ 최신 ☐ 일부 구식 ☐ 구식

**예제 평가**:
```
☐ 완벽함 (✅): 다양한 예제 + 상세 설명 + 실행 가능
☐ 부분적 (⚠️): 기본 예제만 있거나 일부 오래됨
☐ 없음 (❌): 예제 부재
```

**예제 점수**: `${EXAMPLES_SCORE}` / 5

### 4.4 설치 가이드

**설치 가이드 위치**: ${INSTALL_LOCATION}

**설명 항목**:

| 항목 | 존재 | 상세도 |
|------|------|--------|
| 빌드 방법 | ☐ YES ☐ NO | ${BUILD_DETAIL} |
| 의존성 설명 | ☐ YES ☐ NO | ${DEPS_DETAIL} |
| 단계별 설치 | ☐ YES ☐ NO | ${STEP_DETAIL} |
| 플랫폼별 지침 | ☐ YES ☐ NO | ${PLATFORM_DETAIL} |
| 문제 해결 | ☐ YES ☐ NO | ${TROUBLESHOOT_DETAIL} |

**빌드 시스템**:
- ☐ Makefile: ${MAKEFILE_DETAIL}
- ☐ CMake: ${CMAKE_DETAIL}
- ☐ Autotools: ${AUTOTOOLS_DETAIL}
- ☐ 기타: ${OTHER_BUILD}

**설치 가이드 평가**:
```
☐ 완벽함 (✅): 명확한 단계별 설명, 트러블슈팅 포함
☐ 부분적 (⚠️): 기본 설명만 있음
☐ 없음 (❌): 설치 가이드 부재
```

**설치 가이드 점수**: `${INSTALL_SCORE}` / 5

### 4.5 인라인 주석 품질

```
헤더 주석 품질: ☐ 우수 ☐ 양호 ☐부족
코드 주석: ☐ 우수 ☐ 양호 ☐부족
함수 문서: ☐ 우수 ☐ 양호 ☐부족

주석/코드 비율: ${COMMENT_RATIO}%
  ☐ 과다주석 (>30%)
  ☐ 적절 (15-30%)
  ☐ 부족 (<15%)
```

---

## 🔧 5. 빌드 및 테스트

### 5.1 빌드 설정

**빌드 파일 확인**:

| 파일 | 존재 | 상태 |
|------|------|------|
| Makefile | ☐ YES ☐ NO | ☐ 완성 ☐ 부분 ☐ 오류 |
| CMakeLists.txt | ☐ YES ☐ NO | ☐ 완성 ☐ 부분 ☐ 오류 |
| meson.build | ☐ YES ☐ NO | ☐ 완성 ☐ 부분 ☐ 오류 |
| build.sh | ☐ YES ☐ NO | ☐ 완성 ☐ 부분 ☐ 오류 |

**빌드 명령어**:

```bash
${BUILD_COMMAND_1}
${BUILD_COMMAND_2}
${BUILD_COMMAND_3}
```

**빌드 옵션**:
```
☐ DEBUG 빌드: ${DEBUG_FLAG}
☐ RELEASE 빌드: ${RELEASE_FLAG}
☐ 최적화 옵션: ${OPT_FLAG}
☐ 경고 플래그: ${WARN_FLAG}
☐ 추가 옵션: ${EXTRA_FLAGS}
```

**빌드 테스트 결과**:
```
☐ 성공: 오류 없음
☐ 경고 있음: ${WARNING_COUNT}개
  - ${WARNING1}
  - ${WARNING2}
  - ${WARNING3}
☐ 실패: ${ERROR_DESC}
```

**빌드 시간**: ${BUILD_TIME} 초

**빌드 평가**: ☐ 우수 ☐ 양호 ☐ 개선필요

### 5.2 테스트

**테스트 파일 위치**: ${TEST_LOCATION}
**테스트 파일 개수**: `${TEST_FILES}` 개
**테스트 라인 수**: `${TEST_LINES}` 라인

**테스트 범위**:

| 항목 | 상태 | 커버리지 |
|------|------|---------|
| **전체 함수 테스트** | ☐ YES ☐ PARTIAL ☐ NO | ${FUNC_COVERAGE}% |
| **주요 기능 테스트** | ☐ YES ☐ PARTIAL ☐ NO | ${FEATURE_COVERAGE}% |
| **에러 처리 테스트** | ☐ YES ☐ PARTIAL ☐ NO | ${ERROR_COVERAGE}% |
| **엣지 케이스 테스트** | ☐ YES ☐ PARTIAL ☐ NO | ${EDGE_COVERAGE}% |
| **통합 테스트** | ☐ YES ☐ PARTIAL ☐ NO | ${INTEGRATION_COVERAGE}% |

**테스트 실행 방법**:

```bash
${TEST_COMMAND_1}
${TEST_COMMAND_2}
${TEST_COMMAND_3}
```

**테스트 프레임워크**:
- ☐ CUnit
- ☐ Check
- ☐ Minunit
- ☐ 수동 테스트
- ☐ 기타: ${OTHER_FRAMEWORK}

**테스트 결과**:

| 항목 | 값 | 상태 |
|------|-----|------|
| **총 테스트 케이스** | `${TOTAL_TESTS}` | - |
| **통과** | `${PASSED_TESTS}` | ✅ |
| **실패** | `${FAILED_TESTS}` | ❌ |
| **건너뜀** | `${SKIPPED_TESTS}` | ⏭️ |
| **성공률** | `${SUCCESS_RATE}`% | ${SUCCESS_EVAL} |
| **커버리지** | `${COVERAGE}`% | ${COVERAGE_EVAL} |

**테스트 평가**:
```
☐ 우수: 광범위한 커버리지 (>80%) + 높은 성공률
☐ 양호: 기본적인 테스트 존재 (>50%)
☐ 개선필요: 제한된 테스트 (>20%)
☐ 부실: 테스트 없음 또는 매우 제한적
```

**테스트 점수**: `${TEST_SCORE}` / 5

### 5.3 CI/CD

**CI/CD 설정 위치**: ${CI_LOCATION}

**CI/CD 도구**:
- ☐ GitHub Actions (.github/workflows/)
- ☐ GitLab CI (.gitlab-ci.yml)
- ☐ Gitea Actions (.gitea/workflows/)
- ☐ Jenkins (Jenkinsfile)
- ☐ Travis CI (.travis.yml)
- ☐ 기타: ${OTHER_CI}

**CI/CD 파이프라인**:

```yaml
${CI_CONFIG}
```

**자동화 항목**:
- ☐ 자동 빌드
- ☐ 자동 테스트
- ☐ 코드 품질 검사
- ☐ 문서 생성
- ☐ 배포 자동화

---

## 📜 6. 마지막 커밋 정보

### 6.1 최신 커밋

| 항목 | 값 |
|------|-----|
| **Commit Hash** | `${COMMIT_HASH}` |
| **Short Hash** | `${SHORT_HASH}` |
| **메시지** | `${COMMIT_MESSAGE}` |
| **저자** | `${AUTHOR_NAME}` |
| **이메일** | `${AUTHOR_EMAIL}` |
| **날짜** | `${COMMIT_DATE}` |
| **시간** | `${COMMIT_TIME}` |
| **표준시간대** | `${TIMEZONE}` |

### 6.2 커밋 통계

| 항목 | 값 |
|------|-----|
| **총 커밋 수** | `${TOTAL_COMMITS}` |
| **첫 커밋 날짜** | `${FIRST_COMMIT_DATE}` |
| **프로젝트 기간** | `${PROJECT_DAYS}` 일 |
| **평균 커밋 간격** | `${AVG_INTERVAL}` 일 |

### 6.3 최근 활동

| 기간 | 커밋 수 | 평가 |
|------|---------|------|
| 최근 7일 | `${COMMITS_7D}` | ☐ 활발 ☐ 정상 ☐ 저조 |
| 최근 30일 | `${COMMITS_30D}` | ☐ 활발 ☐ 정상 ☐ 저조 |
| 최근 90일 | `${COMMITS_90D}` | ☐ 활발 ☐ 정상 ☐ 저조 |
| 최근 1년 | `${COMMITS_1Y}` | ☐ 활발 ☐ 정상 ☐ 저조 |

### 6.4 유지보수 평가

```
최근 활동 기준:

☐ 활발한 개발 (주 1회 이상 커밋)
  → 높은 유지보수성, 신기술 적용, 빠른 버그 수정

☐ 정기적 유지보수 (월 1회 이상)
  → 안정적 유지보수, 버그 수정 꾸준함

☐ 불규칙한 업데이트 (분기별 또는 그 이상)
  → 느린 반응, 버그 수정 지연 가능

☐ 미유지보수 (1년 이상 커밋 없음)
  → 프로젝트 정체, 보안 위험 가능성
```

**유지보수 점수**: `${MAINTENANCE_SCORE}` / 5

---

## 📦 7. 의존성 분석

### 7.1 외부 라이브러리 의존성

**총 의존성 수**: `${TOTAL_DEPS}` 개

| 순서 | 라이브러리 | 버전 | 용도 | 필수도 |
|------|----------|------|------|--------|
| 1 | `${DEP1_NAME}` | `${DEP1_VER}` | ${DEP1_USE} | ☐ 필수 ☐ 선택 |
| 2 | `${DEP2_NAME}` | `${DEP2_VER}` | ${DEP2_USE} | ☐ 필수 ☐ 선택 |
| 3 | `${DEP3_NAME}` | `${DEP3_VER}` | ${DEP3_USE} | ☐ 필수 ☐ 선택 |

### 7.2 시스템 라이브러리

```
기본 라이브러리:
  ☐ libc: ${LIBC_USE}
  ☐ libm (math): ${LIBM_USE}
  ☐ libpthread: ${PTHREAD_USE}

기타:
  ☐ ${OTHER_LIB1}
  ☐ ${OTHER_LIB2}
  ☐ ${OTHER_LIB3}
```

### 7.3 의존성 관리

**의존성 파일**:
- ☐ package.json
- ☐ requirements.txt
- ☐ vcpkg.json
- ☐ conanfile.txt
- ☐ CMake (찾기 모듈)
- ☐ 관리 안 함 (수동)

### 7.4 의존성 평가

```
의존성 수준:

☐ 순수 C11 라이브러리 (0개 외부 의존성)
  → 최고의 이식성, 높은 자유도

☐ 최소 의존성 (1-2개)
  → 우수한 이식성, 관리 용이

☐ 적당한 의존성 (3-5개)
  → 기능성과 의존성 균형

☐ 높은 의존성 (>5개)
  → 설치 복잡도 증가, 호환성 이슈 가능성
```

**의존성 점수**: `${DEPS_SCORE}` / 5

---

## 📄 8. 라이센스

### 8.1 라이센스 정보

| 항목 | 값 |
|------|-----|
| **라이센스 파일** | ${LICENSE_FILE} |
| **라이센스 종류** | ${LICENSE_TYPE} |
| **라이센스 URL** | ${LICENSE_URL} |

### 8.2 라이센스 타입

```
☐ MIT: 매우 관대한 라이센스, 상업 사용 가능
☐ Apache 2.0: 관대, 특허권 보호
☐ GPL 2.0: 코드공개 의무
☐ GPL 3.0: GPL 2.0의 강화 버전
☐ BSD (2-Clause/3-Clause): 관대
☐ ISC: MIT와 유사
☐ Unlicense: 공개 도메인
☐ 기타: ${OTHER_LICENSE_TYPE}
```

### 8.3 라이센스 명확성

```
☐ 명확함: LICENSE 파일 명시적 존재
☐ 부분적: README에만 언급
☐ 불명확함: 라이센스 정보 불분명
```

---

## ⭐ 9. 품질 점수 계산

### 9.1 점수 항목

#### 1. 코드 품질 (1-5점)

```
평가 기준:
- 모듈화 정도
- 함수 크기 (평균 <50줄 권장)
- 메모리 관리 (누수 위험도)
- 에러 처리 (반환값 검사율)
- 코딩 스타일 일관성

점수:
5점: 프로덕션급 - 모듈화 우수, 메모리 안전, 에러 처리 완벽
4점: 매우 양호 - 구조 명확, 대부분 에러 처리
3점: 양호 - 기본 구조 준수, 일부 에러 처리
2점: 기본 - 구조 미흡, 에러 처리 부족
1점: 부실 - 구조 없음, 에러 처리 미흡

계산: 코드 품질 점수 = ${CODE_QUALITY_SCORE} / 5
```

#### 2. 문서화 (1-5점)

```
평가 기준:
- README 품질
- API 문서 완성도
- 사용 예제 충실도
- 설치 가이드 명확도
- 인라인 주석 품질

점수:
5점: 포괄적 - API + 예제 + 설치 가이드 + 상세 설명
4점: 상세 - API + 예제 + 설치 가이드
3점: 기본 - README + 기본 예제
2점: 최소 - README만 있음
1점: 없음 - README 부재

계산: 문서화 점수 = ${DOCUMENTATION_SCORE} / 5
```

#### 3. 테스트 (1-5점)

```
평가 기준:
- 테스트 파일 존재
- 테스트 라인 수 (소스:테스트 비율 30-50% 권장)
- 테스트 커버리지
- 자동화 테스트

점수:
5점: 완벽 - 단위+통합 테스트, 커버리지 >80%
4점: 광범위 - 주요 기능 테스트, 커버리지 >60%
3점: 기본 - 기본 기능 테스트 존재
2점: 제한적 - 일부 테스트만 존재
1점: 없음 - 테스트 없음

계산: 테스트 점수 = ${TESTING_SCORE} / 5
```

#### 4. 유지보수 (1-5점)

```
평가 기준:
- 최근 커밋 날짜
- 커밋 빈도
- 버그 수정 커밋율
- 기능 추가 커밋율

점수:
5점: 활발 - 주 1회 이상 업데이트
4점: 정기 - 월 1회 이상 업데이트
3점: 불규칙 - 분기별 업데이트
2점: 간헐적 - 연 1-2회 업데이트
1점: 미유지 - 1년 이상 업데이트 없음

계산: 유지보수 점수 = ${MAINTENANCE_SCORE} / 5
```

#### 5. 성능 (1-5점)

```
평가 기준:
- 성능 테스트 여부
- 벤치마크 결과 문서화
- 최적화 노력
- 복잡도 고려

점수:
5점: 벤치마크됨 - 성능 테스트 + 결과 문서화
4점: 최적화 - 성능 고려한 구현
3점: 기본 - 성능 고려
2점: 미흡 - 성능 미고려
1점: 미정의 - 성능 정보 없음

계산: 성능 점수 = ${PERFORMANCE_SCORE} / 5
```

### 9.2 종합 점수

| 항목 | 점수 | 가중치 | 기여도 |
|------|------|--------|--------|
| 코드 품질 | `${CODE_SCORE}` | 25% | ${CODE_WEIGHTED} |
| 문서화 | `${DOC_SCORE}` | 25% | ${DOC_WEIGHTED} |
| 테스트 | `${TEST_SCORE}` | 20% | ${TEST_WEIGHTED} |
| 유지보수 | `${MAINT_SCORE}` | 20% | ${MAINT_WEIGHTED} |
| 성능 | `${PERF_SCORE}` | 10% | ${PERF_WEIGHTED} |
| **평균** | **`${AVG_SCORE}`** | **100%** | **`${TOTAL_SCORE}`** |

### 9.3 등급 평가

```
평균 점수: ${AVG_SCORE} / 5.0

등급:
  ★★★★★ 우수 (Excellent): 4.5점 이상
  ★★★★☆ 양호 (Good): 3.5 - 4.4점
  ★★★☆☆ 개선필요 (Needs Improvement): 2.5 - 3.4점
  ★★☆☆☆ 부실 (Poor): 2.5점 미만

최종 등급: [SELECT]
  ☐ 우수 (Excellent) - 4.5점 이상
  ☐ 양호 (Good) - 3.5점 - 4.4점
  ☐ 개선필요 (Needs Improvement) - 2.5점 - 3.4점
  ☐ 부실 (Poor) - 2.5점 미만
```

---

## 🎯 10. 핵심 특징

### 10.1 주요 강점

1. **${STRENGTH1}**
   - 설명: ${STRENGTH1_DESC}
   - 영향: ${STRENGTH1_IMPACT}

2. **${STRENGTH2}**
   - 설명: ${STRENGTH2_DESC}
   - 영향: ${STRENGTH2_IMPACT}

3. **${STRENGTH3}**
   - 설명: ${STRENGTH3_DESC}
   - 영향: ${STRENGTH3_IMPACT}

### 10.2 주요 약점

1. **${WEAKNESS1}**
   - 설명: ${WEAKNESS1_DESC}
   - 영향: ${WEAKNESS1_IMPACT}

2. **${WEAKNESS2}**
   - 설명: ${WEAKNESS2_DESC}
   - 영향: ${WEAKNESS2_IMPACT}

3. **${WEAKNESS3}**
   - 설명: ${WEAKNESS3_DESC}
   - 영향: ${WEAKNESS3_IMPACT}

---

## 💡 11. 개선 권장사항

### 우선순위 1 (긴급)

1. **${IMPROVEMENT1}**
   - 현황: ${CURRENT_STATE1}
   - 개선 방안: ${ACTION_PLAN1}
   - 기대 효과: ${EXPECTED_BENEFIT1}
   - 예상 시간: ${TIME_ESTIMATE1}

2. **${IMPROVEMENT2}**
   - 현황: ${CURRENT_STATE2}
   - 개선 방안: ${ACTION_PLAN2}
   - 기대 효과: ${EXPECTED_BENEFIT2}
   - 예상 시간: ${TIME_ESTIMATE2}

### 우선순위 2 (중요)

1. **${IMPROVEMENT3}**
   - 현황: ${CURRENT_STATE3}
   - 개선 방안: ${ACTION_PLAN3}
   - 기대 효과: ${EXPECTED_BENEFIT3}

2. **${IMPROVEMENT4}**
   - 현황: ${CURRENT_STATE4}
   - 개선 방안: ${ACTION_PLAN4}
   - 기대 효과: ${EXPECTED_BENEFIT4}

### 우선순위 3 (선택)

1. **${IMPROVEMENT5}**
   - 설명: ${IMPROVEMENT5_DESC}

---

## 🏆 12. 최종 평가

### 종합 평가

```
저장소명: ${REPO_NAME}
분석 날짜: ${ANALYSIS_DATE}
최종 점수: ${FINAL_SCORE} / 5.0
최종 등급: [최종 등급]
```

### 평가 요약

${EVALUATION_SUMMARY}

### 추천 여부

```
☐ 강하게 추천 (우수)
  → 프로덕션 사용 권장, 높은 품질

☐ 추천 (양호)
  → 일반적 사용 가능, 몇 가지 개선 필요

☐ 유의해서 사용 (개선필요)
  → 사용 가능하나 주의 필요, 개선 권장

☐ 권장하지 않음 (부실)
  → 사용 부적절, 주요 개선 필요
```

### 사용 적합성

| 용도 | 적합도 | 비고 |
|------|--------|------|
| **프로덕션** | ☐ 적합 ☐ 조건부 ☐ 부적합 | ${PROD_NOTE} |
| **개발** | ☐ 적합 ☐ 조건부 ☐ 부적합 | ${DEV_NOTE} |
| **교육** | ☐ 적합 ☐ 조건부 ☐ 부적합 | ${EDU_NOTE} |
| **참고 코드** | ☐ 적합 ☐ 조건부 ☐ 부적합 | ${REF_NOTE} |

---

## 📋 부록

### A. 상세 파일 목록

(소스 파일, 헤더 파일, 테스트 파일 상세 목록)

### B. 코드 샘플

(주요 함수, 구조체 정의 등)

### C. 테스트 결과

(테스트 실행 로그, 결과 요약)

### D. 참고 자료

- [공식 저장소](https://gogs.dclub.kr/kim/${REPO_NAME})
- [라이센스](${LICENSE_URL})
- [기여 가이드](${CONTRIBUTING_URL})

---

**분석 완료**
작성일: ${ANALYSIS_DATE}
분석자: ${ANALYST_NAME}

