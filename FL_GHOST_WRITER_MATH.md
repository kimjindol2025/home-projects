# FL-Ghost-Writer: Bezier → SDF 수학 명세

## 1. Bezier 곡선 기초

### 1.1 Quadratic Bezier (TrueType 표준)

```
B(t) = (1-t)²P₀ + 2(1-t)tP₁ + t²P₂,  t ∈ [0,1]
```

**프로퍼티**:
- P₀: 시작점 (anchor)
- P₁: 제어점 (off-curve)
- P₂: 끝점 (anchor)
- 전체 glyph = 여러 Quadratic Bezier 곡선의 연결

### 1.2 Cubic Bezier (OpenType 표준)

```
B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
```

---

## 2. SDF (Signed Distance Field) 개념

### 2.1 정의

SDF는 모든 픽셀 (x, y)에 대해:

```
SDF(x, y) = {
    +distance(x, y, curve)   if (x, y) is outside the glyph
    -distance(x, y, curve)   if (x, y) is inside the glyph
}
```

**의미**:
- **양수**: 픽셀이 glyph 경계 밖
- **음수**: 픽셀이 glyph 경계 안
- **절댓값**: 경계까지의 최단 거리

### 2.2 장점

| 특성 | 설명 |
|------|------|
| **벡터화** | 해상도 독립적 렌더링 |
| **안티에일리어싱** | 거리 값으로 smooth edge |
| **확장성** | 모든 해상도에서 고품질 |
| **메모리 효율** | 32-bit float로 충분 |

---

## 3. 점(Point)에서 Bezier 곡선까지의 거리

### 3.1 Quadratic Bezier 거리 계산

**입력**: 점 P, Bezier 곡선 B(t)
**출력**: min_distance(P, B(t))

#### 알고리즘 (Newton-Raphson)

```
1. 초기값 설정: t₀ = 0.5
2. 반복:
   d(t) = |P - B(t)|² (제곱거리, 미분 용이)
   t_new = t - d'(t) / d''(t)
   if |t_new - t| < ε: break
3. t ∈ [0, 1] 범위로 클립
4. min_dist = |P - B(t_final)|
```

#### FreeLang 구현 (pseudo-code)

```freelang
fn distance_point_to_quadratic_bezier(
    p: Vector2,
    p0: Vector2, p1: Vector2, p2: Vector2
) -> f32 {
    // 1단계: 초기 t 값
    let mut t = 0.5;

    for _ in 0..10 {  // Newton-Raphson 반복 (10회)
        // B(t) 계산
        let b_t = quadratic_bezier(t, p0, p1, p2);

        // d(t) = |P - B(t)|²
        let diff = p - b_t;
        let d_t = diff.dot(diff);

        // B'(t) 계산
        let b_prime = quadratic_bezier_derivative(t, p0, p1, p2);

        // d'(t) = 2 * (P - B(t)) · (-B'(t))
        let d_prime = 2.0 * diff.dot(-b_prime);

        // B''(t) 계산
        let b_double_prime = quadratic_bezier_second_derivative(p0, p1, p2);

        // d''(t) = 2 * B'(t)·B'(t) + 2 * (P-B(t))·B''(t)
        let d_double_prime =
            2.0 * b_prime.dot(b_prime) +
            2.0 * diff.dot(b_double_prime);

        // Newton-Raphson 업데이트
        if d_double_prime.abs() > 1e-6 {
            let t_new = t - d_prime / d_double_prime;
            t = clamp(t_new, 0.0, 1.0);
        } else {
            break;
        }
    }

    // 최종 거리
    let b_final = quadratic_bezier(t, p0, p1, p2);
    (p - b_final).magnitude()
}
```

### 3.2 필요한 보조 함수들

#### Quadratic Bezier 계산

```freelang
fn quadratic_bezier(
    t: f32,
    p0: Vector2, p1: Vector2, p2: Vector2
) -> Vector2 {
    let mt = 1.0 - t;
    let mt2 = mt * mt;
    let t2 = t * t;

    p0 * mt2 + p1 * 2.0 * mt * t + p2 * t2
}
```

#### 1차 도함수

```
B'(t) = 2(1-t)(P₁-P₀) + 2t(P₂-P₁)
      = 2[(P₁-P₀) + t(P₀-2P₁+P₂)]
```

```freelang
fn quadratic_bezier_derivative(
    t: f32,
    p0: Vector2, p1: Vector2, p2: Vector2
) -> Vector2 {
    let mt = 1.0 - t;
    (p1 - p0) * 2.0 * mt + (p2 - p1) * 2.0 * t
}
```

#### 2차 도함수

```
B''(t) = 2(P₀ - 2P₁ + P₂)
```

```freelang
fn quadratic_bezier_second_derivative(
    p0: Vector2, p1: Vector2, p2: Vector2
) -> Vector2 {
    (p0 - p1 * 2.0 + p2) * 2.0
}
```

---

## 4. Scanline 기반 고속 SDF 생성

### 4.1 개요

glyph의 모든 픽셀에 대해 거리를 계산하는 것은 O(w × h × curves)이므로 느립니다.

**최적화**: Scanline sweeping

```
for y in 0..height:
    for x in 0..width:
        distance_to_all_curves[x][y] = min(
            distance(point(x,y), curve_0),
            distance(point(x,y), curve_1),
            ...
        )
```

### 4.2 거리 필드 버퍼

```freelang
pub struct SDFBuffer {
    pub width: u32,
    pub height: u32,
    pub data: Vec<f32>,  // 32-bit float SDF 값
}

impl SDFBuffer {
    fn set_sdf(&mut self, x: u32, y: u32, distance: f32) {
        let index = (y * self.width + x) as usize;
        self.data[index] = distance;
    }

    fn get_sdf(&self, x: u32, y: u32) -> f32 {
        let index = (y * self.width + x) as usize;
        self.data[index]
    }
}
```

---

## 5. 렌더링: SDF → 화면

### 5.1 알파 계산

SDF 값을 화면상 픽셀의 알파(투명도)로 변환:

```
alpha(sdf) = {
    1.0                      if sdf ≤ -spread
    0.5 + sdf/(2*spread)     if -spread < sdf < +spread
    0.0                      if sdf ≥ +spread
}
```

**spread** = 품질 파라미터 (보통 0.5~2.0 픽셀)

```freelang
fn sdf_to_alpha(sdf: f32, spread: f32) -> f32 {
    let half_spread = spread * 0.5;
    if sdf <= -spread {
        1.0
    } else if sdf >= spread {
        0.0
    } else {
        0.5 + sdf / (2.0 * spread)
    }
}
```

### 5.2 DMA를 통한 고속 합성

```
[SDF Buffer (GPU)]
        ↓
   [ALPHA BLEND]
        ↓
[Framebuffer]
```

**DMA 타일링**: 너비 64픽셀 × 높이 64픽셀 타일 단위로 처리

```
타일당 계산:
- SDF 로드: ~4KB
- Alpha 계산: ~10µs
- Framebuffer 쓰기: ~10µs (DMA 병렬화)

총 시간: 16µs/glyph (64×64 타일 기준)
```

---

## 6. 무관용 규칙 (Unforgiving Rules)

| # | 규칙 | 목표 | 검증 |
|----|------|------|------|
| 1 | 폰트 렌더 < 100µs | single glyph | 정확한 시간 측정 |
| 2 | SDF 메모리 < 2MB | full ASCII set (256) | 버퍼 크기 검증 |
| 3 | CPU 점유율 < 1% | DMA 가속 | 성능 프로파일링 |
| 4 | 버퍼 오버플로우 = 0 | .ttf 파싱 | 형식 검증 |
| 5 | 품질 > 95% | anti-aliasing | 시각 검사 |

---

## 7. 구현 체크리스트

### Phase 1: Font Parser
- [ ] TTF 헤더 파싱
- [ ] 'glyf' 테이블 읽기
- [ ] Outline 좌표 추출
- [ ] Quadratic Bezier 곡선 재구성

### Phase 2: SDF Generator
- [ ] Newton-Raphson 거리 계산
- [ ] Scanline 고속 처리
- [ ] SDF 버퍼 관리
- [ ] 메모리 최적화

### Phase 3: DMA & Rendering
- [ ] DMA 컨트롤러 제어
- [ ] Alpha blending 구현
- [ ] 타일 기반 처리
- [ ] 성능 측정

### Phase 4: Sovereign Console
- [ ] 부팅 화면 설계
- [ ] 텍스트 배치
- [ ] 시간 측정 (16µs 목표)
- [ ] GOGS 푸시

---

## 8. 수학적 검증

### 8.1 거리 함수의 볼록성 (Convexity)

Quadratic Bezier 곡선까지의 거리 함수는 **convex**이므로 Newton-Raphson이 수렴 보장:

```
d''(t) = 2(|B'(t)|² + (P-B(t))·B''(t)) ≥ 0
```

### 8.2 오버플로우 방지

TTF 파일의 최대 glyph 크기:
```
max_width ≤ 4096 픽셀 (TTF 스펙)
SDF 버퍼: width × 32-bit = 16KB per glyph
256 glyphs = 256 × 16KB = 4MB (충분)
```

---

## 참고 자료

- **TrueType Spec**: https://docs.microsoft.com/en-us/typography/opentype/spec/
- **SDF Paper**: "Improper Contouring of Closed Curves"
- **Scanline**: "An Introduction to DIGITAL IMAGE PROCESSING" (Gonzalez & Woods)

---

**생성**: 2026-03-04
**목적**: FL-Ghost-Writer Phase 2 (SDF Generator) 구현 기초
