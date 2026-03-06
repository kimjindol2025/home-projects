# FL-Ghost-Writer: 독립형 텍스트 렌더링 엔진

**프로젝트명**: FL-Ghost-Writer (유령 작가)
**난이도**: 7 (심화 수준의 '휴식')
**목표**: 외부 폰트 라이브러리 없이 순수 FreeLang으로 구현한 SDF 기반 텍스트 렌더링
**철학**: "폰트를 남의 손에 맡기는 순간, 우리는 또 다시 갇힌다."

---

## 📋 프로젝트 개요

### 목표 명세

| 항목 | 요구사항 |
|------|---------|
| **의존도** | Rust/C/FreeType 0% |
| **메모리** | < 64KB (SDF 버퍼 포함) |
| **성능** | 단일 glyph 렌더 < 100µs |
| **해상도** | 해상도 독립적 (벡터 기반) |
| **부팅 화면** | "FreeLang Sovereign OS v1.0" 16µs 내 출력 |

### 핵심 기술

1. **Font Parser**: TTF 바이너리 파싱 (의존도 0%)
2. **Bezier 곡선**: Quadratic Bezier 거리 계산 (Newton-Raphson)
3. **SDF 생성**: 32-bit float 거리 필드
4. **DMA 가속**: CPU 점유율 < 1% 렌더링
5. **수학적 검증**: 버퍼 오버플로우 형식 검증

---

## 🏗️ 4개 Phase 구조

### Phase 1: Font Parser (Week 1-2)
**목표**: TTF 파일을 Bezier 곡선으로 변환
**규모**: ~500줄

**구현 항목**:
1. `ttf_parser.fl` (300줄)
   - TTF 헤더 파싱 (offset table)
   - 'glyf' 테이블 읽기
   - 윤곽선 데이터 추출
   - Quadratic Bezier 곡선 재구성

2. `glyph_data.fl` (200줄)
   - Glyph 메타데이터 구조
   - Bezier 곡선 저장소
   - 변환 행렬 (scale, translate)

**무관용 규칙**:
- [ ] TTF 파일 유효성 검증 (magic 0x00010000)
- [ ] 256 glyph 로드 완료 < 100ms
- [ ] 메모리 누수 = 0
- [ ] 버퍼 오버플로우 = 0

**테스트**:
- P1-T1: TTF 헤더 파싱 (magic, table count)
- P1-T2: 'glyf' 테이블 위치 계산
- P1-T3: Quadratic outline 데이터 추출
- P1-T4: 좌표 변환 및 스케일링

---

### Phase 2: SDF Generator (Week 3-4)
**목표**: Bezier 곡선을 수학적 거리 필드로 변환
**규모**: ~400줄

**구현 항목**:
1. `bezier_math.fl` (250줄)
   - Quadratic Bezier 계산: B(t)
   - 1차/2차 도함수: B'(t), B''(t)
   - Newton-Raphson 점-곡선 거리 계산
   - 고정밀 부동소수점 연산

2. `sdf_generator.fl` (150줄)
   - Scanline 기반 처리
   - SDF 버퍼 관리
   - 거리 필드 생성 (32-bit float)

**무관용 규칙**:
- [ ] 점-Bezier 거리 < 1µs
- [ ] SDF 생성 < 100µs per glyph
- [ ] Newton-Raphson 수렴성 보증 (convex 함수)
- [ ] 정밀도 > 0.01 픽셀

**테스트**:
- P2-T1: Bezier 곡선 계산 (B(t), B'(t), B''(t))
- P2-T2: 점-곡선 거리 (Newton-Raphson)
- P2-T3: SDF 버퍼 생성 (해상도 64×64)
- P2-T4: 부호 결정 (inside/outside)

---

### Phase 3: DMA & Rendering (Week 5)
**목표**: GPU/DMA를 통한 고속 렌더링
**규모**: ~300줄

**구현 항목**:
1. `dma_controller.fl` (150줄)
   - DMA 채널 제어
   - 타일 기반 메모리 전송
   - CPU 점유율 0%에 가까운 렌더링

2. `alpha_blending.fl` (100줄)
   - SDF → Alpha 변환
   - RGBA 합성
   - 안티에일리어싱

3. `framebuffer_compositor.fl` (50줄)
   - 다층 합성 (text + background)
   - 색상 관리

**무관용 규칙**:
- [ ] DMA 전송 대역폭 > 10GB/s
- [ ] Alpha 계산 < 5µs per pixel
- [ ] CPU 점유율 < 1% (측정)
- [ ] 품질 > 95% (시각 검사)

**테스트**:
- P3-T1: DMA 타일 전송 (64×64)
- P3-T2: SDF → Alpha 변환
- P3-T3: 색상 합성 (RGBA blend)
- P3-T4: 다중 glyph 렌더링

---

### Phase 4: Sovereign Console (Week 6)
**목표**: 부팅 화면 렌더링 (16µs 목표)
**규모**: ~200줄

**구현 항목**:
1. `boot_screen.fl` (150줄)
   - "FreeLang Sovereign OS v1.0" 텍스트 배치
   - 배경 색상 (프로그래시브 블루)
   - 애니메이션 (선택사항)

2. `performance_meter.fl` (50줄)
   - 렌더 시간 측정
   - 성능 카운터

**무관용 규칙**:
- [ ] 부팅 화면 < 16µs (측정)
- [ ] 텍스트 선명도 > 98% (95점이상)
- [ ] 메모리 사용 < 100KB
- [ ] CPU 코어 1개 점유율 < 10%

**테스트**:
- P4-T1: 텍스트 메타데이터 (글자 위치, 크기)
- P4-T2: 단일 glyph 렌더 (< 16µs)
- P4-T3: 전체 문구 렌더링
- P4-T4: 부팅 시간 측정

---

## 📊 구현 통계

| Phase | 파일 | 줄 수 | 주요 작업 | 기간 |
|-------|------|-------|---------|------|
| 1 | ttf_parser.fl, glyph_data.fl | 500 | TTF 파싱, Bezier 추출 | Week 1-2 |
| 2 | bezier_math.fl, sdf_generator.fl | 400 | SDF 생성, 거리 계산 | Week 3-4 |
| 3 | dma_controller.fl, compositor.fl | 300 | DMA 가속, 렌더링 | Week 5 |
| 4 | boot_screen.fl, meter.fl | 200 | 부팅 화면, 성능 측정 | Week 6 |
| Tests | test_suite.fl | 600 | 20개 무관용 테스트 | 전체 |
| **합계** | **8개** | **~2,000줄** | **완전 독립형** | **6주** |

---

## 🔬 무관용 규칙 (Unforgiving Rules)

### Tier 1: 기능 검증

| # | 규칙 | 목표 | 검증 방법 |
|----|------|------|---------|
| R1 | TTF 파싱 정확도 | 100% 정상 문자 | 표준 폰트 테스트 |
| R2 | Bezier 재구성 | 오차 < 0.01px | 좌표 비교 |
| R3 | SDF 생성 | 정밀도 > 98% | 거리 필드 검증 |
| R4 | 렌더링 품질 | 시각 점수 > 95 | 화질 검사 |
| R5 | 독립성 | Rust/C = 0% | 의존도 분석 |

### Tier 2: 성능 검증

| # | 규칙 | 목표 | 측정 |
|----|------|------|------|
| R6 | 단일 glyph | < 100µs | 시간 프로파일링 |
| R7 | 부팅 화면 | < 16µs | 정밀 측정 |
| R8 | DMA 점유율 | < 1% | CPU 모니터링 |
| R9 | 메모리 | < 64KB | 버퍼 추적 |
| R10 | 오버플로우 | = 0 | 형식 검증 |

---

## 📂 디렉토리 구조

```
~/freelang-ghost-writer/
├── src/
│   ├── fonts/
│   │   ├── ttf_parser.fl       (Phase 1)
│   │   ├── glyph_data.fl       (Phase 1)
│   │   └── mod.fl
│   ├── sdf/
│   │   ├── bezier_math.fl      (Phase 2)
│   │   ├── sdf_generator.fl    (Phase 2)
│   │   └── mod.fl
│   ├── rendering/
│   │   ├── dma_controller.fl   (Phase 3)
│   │   ├── alpha_blending.fl   (Phase 3)
│   │   ├── compositor.fl       (Phase 3)
│   │   └── mod.fl
│   ├── dma/
│   │   ├── boot_screen.fl      (Phase 4)
│   │   ├── perf_meter.fl       (Phase 4)
│   │   └── mod.fl
│   └── lib.fl                   (통합)
├── tests/
│   ├── phase1_parser_tests.fl   (10개 테스트)
│   ├── phase2_sdf_tests.fl      (10개 테스트)
│   ├── phase3_render_tests.fl   (5개 테스트)
│   └── phase4_boot_tests.fl     (5개 테스트)
├── docs/
│   ├── MATH_SPEC.md             (수학 명세)
│   ├── TTF_SPEC.md              (TTF 형식)
│   └── PERFORMANCE.md           (성능 분석)
├── assets/
│   └── font_samples.ttf         (테스트 폰트)
└── PROJECT_PLAN.md              (이 파일)
```

---

## 🎯 주요 도전 과제

### 1. TTF 파싱 (Phase 1)
**문제**: 이진 형식의 복잡한 TTF 구조
**해결**: 순차적 파싱, 오프셋 테이블 추적
**검증**: 각 단계마다 체크섬 확인

### 2. Bezier 거리 계산 (Phase 2)
**문제**: Newton-Raphson의 수렴성 보증
**해결**: Convexity 수학적 증명
**검증**: 100000회 반복 거리 계산

### 3. DMA 가속 (Phase 3)
**문제**: CPU-GPU 동기화, 메모리 일관성
**해결**: 명시적 배리어, 타일 단위 처리
**검증**: CPU 점유율 실시간 모니터링

### 4. 16µs 부팅 (Phase 4)
**문제**: 극도로 짧은 시간 제약
**해결**: SDF 사전 캐싱, DMA 병렬화
**검증**: 나노초 단위 시간 측정

---

## 📝 구현 시작 지점

### 우선순위

1. **즉시 시작**: Phase 1 (TTF Parser)
   - 가장 중요한 기초
   - 다른 Phase에 의존하지 않음
   - 명확한 명세

2. **병렬 진행**: Bezier Math 라이브러리 (Phase 2)
   - Phase 1과 독립적
   - 단위 테스트 가능

3. **의존성 후**: Phase 3, 4
   - Phase 1, 2의 결과 필요

---

## 🔐 안전 검증

### 형식 검증

```freelang
// TTF 매직 넘버
magic_number ∈ {0x00010000, 'OTTO', 'true'}

// Glyph 수
num_glyphs ∈ [1, 65535]

// Offset 유효성
offset < file_size
```

### 버퍼 오버플로우 방지

```freelang
// 최대 glyph 크기 (TTF 스펙)
max_glyph_width = 4096px
max_points_per_glyph = 1000

// SDF 버퍼
sdf_size = 64 × 64 × 4 bytes = 16KB per glyph
256 glyphs = 4MB (충분)
```

---

## 결론

**FL-Ghost-Writer**는 FreeLang의 완전한 독립성을 입증하는 프로젝트입니다.

> "기록이 증명이다."

이 프로젝트가 완성되면:
- ✅ 외부 폰트 라이브러리 의존도 0%
- ✅ 순수 FreeLang 2,000줄로 전문가급 텍스트 렌더링
- ✅ 10배 빠른 부팅 화면 (16µs)
- ✅ Nano-Kernel의 "시각적 영혼" 완성

**다음**: Phase 1 (TTF Parser) 즉시 시작

---

**생성**: 2026-03-04
**상태**: 🚀 **READY TO START**
