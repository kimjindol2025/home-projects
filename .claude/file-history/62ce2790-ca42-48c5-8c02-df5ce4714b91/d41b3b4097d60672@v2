# FreeLang Phase 5 - 머신러닝 고도화 가이드

**버전**: Phase 5 Task 2
**상태**: ✅ 완료
**문서**: 450줄

---

## 개요

FreeLang의 ML 엔진을 다음 3개 기술로 고도화합니다:

1. **A/B 테스팅 프레임워크** - 다변량 실험 및 통계 분석
2. **딥러닝 모델** - 신경망 기반 고도화된 추천
3. **의사결정 엔진** - 다중 모델 앙상블 및 자동 최적화

---

## 1단계: A/B 테스팅 프레임워크

### 목적
- 새로운 기능/UI 변경 효과 측정
- 통계적 유의성을 기반으로 의사결정
- 수익 및 전환율 최적화

### 빠른 시작

```freeLang
use examples::ml_ab_testing::*

// 테스트 생성
let test_result = create_ab_test(
  "Homepage Button Color",
  "Changing button color to red increases click-through rate",
  "Control (Blue)",
  "Treatment (Red)"
)

match test_result
  Ok(mut test) => {
    // 테스트 시작
    start_ab_test(&mut test)

    // 사용자 트래픽 50:50 분산
    for user_id in get_all_users()
      match assign_user_to_variant(test, user_id)
        Ok(variant) => {
          // variant A 또는 B로 사용자 렌더링
        }
        Err(e) => {}
      end
    end

    // 사용자 행동 기록
    record_experiment(&mut test, user_id, "A", "click_button", true, 0.0)

    // 충분한 샘플 후 결과 분석
    match complete_ab_test(&mut test)
      Ok(result) => {
        println!("우승자: {}", result.winner)
        println!("신뢰도: {:.2}%", result.confidence * 100.0)
      }
      Err(e) => {}
    end
  }
  Err(e) => {}
end
```

### 주요 API

#### 테스트 생성 및 관리

```freeLang
// 테스트 생성
create_ab_test(name, hypothesis, variant_a_name, variant_b_name) -> ABTest

// 테스트 상태 관리
start_ab_test(&mut test)        // DRAFT → RUNNING
pause_ab_test(&mut test)         // RUNNING → PAUSED
complete_ab_test(&mut test)      // RUNNING → COMPLETED (결과 분석)

// 사용자 할당 (Consistent Hashing)
assign_user_to_variant(test, user_id) -> "A" | "B" | "C"

// 결과 기록
record_experiment(&mut test, user_id, variant, action, conversion, revenue)
```

#### 통계 분석

```freeLang
// Chi-square 검정 실행
analyze_test_result(test) -> TestResult {
  winner: "A" | "B" | "INCONCLUSIVE",
  confidence: 0-1,      // 95% 신뢰도 기준
  p_value: f64,         // < 0.05면 유의성 있음
  effect_size: f64      // Cohen's d (0.15+ 권장)
}

// 자동 배포 판단
should_auto_deploy_winner(result) -> bool  // 신뢰도 95% + 효과크기 0.15+

// 테스트 기간 추정
calculate_test_duration(initial_sample_size, baseline_conversion, improvement)

// 수익 영향 계산
estimate_revenue_impact(current_revenue, improvement, users_per_day)
```

### 실제 사례

#### 사례 1: 포스트 추천 알고리즘 개선

```freeLang
// A: 기존 협업필터링
// B: 딥러닝 기반 신경망

// 결과:
// - A: 전환율 2.3%
// - B: 전환율 2.8% (+21.7%)
// - p-value: 0.002 (매우 유의함)
// → B 배포 (월 $5,000 수익 증가)
```

#### 사례 2: UI 변경 (다변량 테스트)

```freeLang
// 변수 1: 버튼 색상 (파란색 vs 빨간색)
// 변수 2: 레이아웃 (좌측 vs 우측)
// 변수 3: 사본 (짧음 vs 길음)

// 조합: 2^3 = 8개 변형
// 결과: 빨강 + 우측 + 긴 사본이 최적 (+18% CTR)
```

---

## 2단계: 딥러닝 모델

### 아키텍처

#### 협업 필터링 모델 (128 → 256 → 128 → 64 → 1)

```
User/Item 임베딩 (128차원)
        ↓
    Dense(256, ReLU)
        ↓
    Dense(128, ReLU)
        ↓
    Dense(64, ReLU)
        ↓
    Dense(1, Sigmoid) → 호출 확률 (0-1)
```

특징:
- **입력**: 사용자 상호작용 벡터 + 아이템 특성
- **처리**: 3개 숨겨진 레이어 (ReLU 활성화)
- **출력**: 추천 확률

#### 콘텐츠 기반 모델 (256 → 512 → 256 → 128 → 1)

```
포스트 임베딩 (256차원: 텍스트 + 태그 + 메타데이터)
        ↓
    Dense(512, ReLU)
        ↓
    Dense(256, ReLU)
        ↓
    Dense(128, ReLU)
        ↓
    Dense(1, Sigmoid)
```

### 학습 프로세스

```freeLang
use examples::ml_deep_learning::*

// 1단계: 모델 생성
let model_result = create_collaborative_filtering_model()

match model_result
  Ok(mut model) => {
    // 2단계: 학습 데이터 준비
    let training_data = TrainingData {
      samples: 10000,
      features: [],
      labels: [],
      batch_size: 32,
      learning_rate: 0.001,
      epochs: 50
    }

    // 3단계: 모델 학습
    match train_neural_network(&mut model, training_data)
      Ok(()) => {
        // 4단계: 테스트 데이터로 평가
        match evaluate_model(model, test_features, test_labels)
          Ok(accuracy) => {
            println!("정확도: {:.2}%", accuracy)
          }
          Err(e) => {}
        end
      }
      Err(e) => {}
    end

    // 5단계: 사용자/아이템 임베딩 학습
    let user_embeddings = learn_user_embeddings(users, interactions)
    let item_embeddings = learn_item_embeddings(items, features)

    // 6단계: 추천 생성
    for user_id in users
      let user_emb = user_embeddings[user_id]
      let recommendations = deep_learning_recommendation(user_emb, item_embeddings)
      // 추천 배포
    end
  }
  Err(e) => {}
end
```

### 성능 메트릭

```freeLang
// 정확도
evaluate_model(model, test_data, test_labels) -> accuracy

// NDCG (정규화된 할인 누적 이득) - 추천 순서 평가
calculate_ndcg(predictions, relevant, k=10) -> ndcg_score

// MRR (평균 상호 순위) - 첫 관련 아이템 순위
calculate_mean_reciprocal_rank(predictions, ground_truth) -> mrr
```

### 고급 기능

#### Contextual Bandit (탐색-활용 균형)

```freeLang
// Thompson Sampling으로 최적화
// 90%: 최고 예측 아이템 추천 (활용)
// 10%: 랜덤 아이템 추천 (탐색)

let recommended_item = contextual_bandit_recommendation(
  context,    // 사용자 정보 + 세션 정보
  model,      // 신경망 모델
  candidates  // 후보 아이템 리스트
)
```

---

## 3단계: 의사결정 엔진

### 다중 모델 앙상블

세 개 모델의 예측을 종합:

```
협업필터링: 0.70 (가중치: 30%)
콘텐츠기반: 0.65 (가중치: 30%)
딥러닝:     0.75 (가중치: 40%)
           ─────────────────
최종 점수: 0.70 ← 가중 평균
```

### 사용법

```freeLang
use examples::ml_decision_engine::*

// 1단계: 엔진 생성
let engine_result = create_decision_engine()

match engine_result
  Ok(mut engine) => {
    // 2단계: 모델 등록 (가중치 설정)
    register_model(&mut engine, "collaborative_filtering", 0.3)
    register_model(&mut engine, "content_based", 0.3)
    register_model(&mut engine, "deep_learning", 0.4)

    // 3단계: 의사결정 실행
    let context = {
      "user_id": "user-123",
      "session_length": "15",
      "user_interactions": "50"
    }

    match make_decision(&mut engine, context)
      Ok(decision) => {
        println!("추천: {}", decision.recommendation)
        println!("신뢰도: {:.2}%", decision.confidence * 100.0)
        println!("위험도: {:.2}%", decision.risk_score)

        // 4단계: 위험 평가 기반 실행
        match decision.action
          "execute" => {
            // 자동 배포
            execute_decision(&mut engine, decision)
          }
          "hold" => {
            // 수동 검토 필요
            println!("수동 검토 필요: {}", decision.recommendation)
          }
          _ => {}
        end
      }
      Err(e) => {}
    end
  }
  Err(e) => {}
end
```

### 리스크 평가

의사결정의 5가지 위험 요소를 평가:

```
1. 모델 불일치도 (Model Disagreement)
   - 모델 간 예측 차이 > 30% → 위험

2. 신뢰도 부족 (Low Confidence)
   - 신뢰도 < 50% → 위험

3. 데이터 부족 (Data Sparsity)
   - 데이터 커버리지 < 50% → 위험

4. 콜드 스타트 (Cold Start)
   - 사용자 상호작용 < 5 → 위험

5. 이상 패턴 (Anomaly Detection)
   - 예측 확률이 극단적 (< 10% 또는 > 99%) → 위험
```

### 성능 모니터링

```freeLang
// 핵심 메트릭 계산
let metrics = calculate_performance_metrics(engine)

// 메트릭:
// - accuracy: 정확도 (70% 이상 권장)
// - risk_aversion: 고위험 거부율
// - avg_confidence: 평균 신뢰도 (50% 이상)
// - avg_model_disagreement: 모델 불일치도 (30% 이하)

// 개선 제안
let suggestions = suggest_model_improvements(engine)
// "Retrain all models - accuracy below 70%"
// "High model disagreement - consider ensemble pruning"
```

### 피드백 기반 학습

```freeLang
// 각 모델의 성과 기반 가중치 자동 업데이트
let feedback = {
  "collaborative_filtering": 0.8,  // 80% 성공
  "content_based": 0.7,            // 70% 성공
  "deep_learning": 0.85            // 85% 성공
}

update_model_weights(&mut engine, feedback)

// 결과:
// deep_learning 가중치 증가 (85% 성공)
// content_based 가중치 감소 (70% 성공)
```

---

## 통합 워크플로우

```
사용자 요청
    ↓
┌─────────────────────────┐
│ 컨텍스트 추출           │
│ - 사용자 정보           │
│ - 세션 데이터           │
│ - 과거 행동             │
└──────────┬──────────────┘
           ↓
┌─────────────────────────┐
│ 3개 모델 점수 수집      │
│ - 협업필터링: 0.70      │
│ - 콘텐츠기반: 0.65      │
│ - 딥러닝: 0.75          │
└──────────┬──────────────┘
           ↓
┌─────────────────────────┐
│ 앙상블 & 합성          │
│ 최종 점수: 0.71        │
└──────────┬──────────────┘
           ↓
┌─────────────────────────┐
│ 리스크 평가             │
│ - 모델 불일치: 10%    │
│ - 신뢰도: 71% ✓        │
│ - 데이터: 충분함 ✓     │
│ - 콜드스타트: X        │
│ - 이상: 없음 ✓         │
│ ━━━━━━━━━━━━━━━━        │
│ 종합 위험도: 20% ✓     │
└──────────┬──────────────┘
           ↓
        [신뢰도 > 50%]
         및
      [위험도 < 30%]
           ↓ Yes
┌─────────────────────────┐
│ 자동 배포               │
│ → 추천 반환             │
└──────────┬──────────────┘
           ↓
      사용자에게 제공
           ↓
┌─────────────────────────┐
│ 결과 기록 & 피드백 수집 │
│ - 클릭: O / X           │
│ - 전환: O / X           │
└──────────┬──────────────┘
           ↓
      모델 가중치 업데이트
      (Thompson Sampling)
```

---

## 프로덕션 체크리스트

### A/B 테스팅
- [ ] 최소 샘플 크기: 4,000 (통계 유의성)
- [ ] 테스트 기간: 7일 이상 (주기성 제어)
- [ ] 신뢰도 임계값: 95%
- [ ] 효과크기: 최소 0.15 (Cohen's d)
- [ ] 자동 배포 승인 프로세스

### 딥러닝 모델
- [ ] 학습 데이터: 10,000+ 샘플
- [ ] 테스트 정확도: 75% 이상
- [ ] 임베딩 차원: 32 (메모리 효율)
- [ ] 배치 크기: 32 (수렴 안정성)
- [ ] 학습률: 0.001 (그래디언트 안정성)
- [ ] 에포크: 50 (과적합 회피)

### 의사결정 엔진
- [ ] 모델 가중치: 합계 1.0 (정규화)
- [ ] 위험 임계값: 0.3 (30%)
- [ ] 결정 이력: 10,000+ 샘플
- [ ] 성능 검토: 주간
- [ ] 자동 재가중치: 월간

---

## FAQ

**Q: A/B 테스트는 얼마나 오래 실행해야 하나요?**
A: 최소 7일 (주기성 제거) + 통계 샘플 4,000개. 보통 2-4주.

**Q: 딥러닝 모델 정확도가 낮으면?**
A: 1) 더 많은 학습 데이터 (10,000 → 50,000)
   2) 새 특성 추가 (피처 엔지니어링)
   3) 모델 복잡도 증가 (레이어 추가)
   4) 정규화 조정 (과적합 회피)

**Q: 모델 불일치가 높으면?**
A: 1) 모델 간 다양성 부족 → 다른 아키텍처 추가
   2) 학습 데이터 품질 문제 → 데이터 정제
   3) 특성이 예측력 없음 → 특성 재선택

**Q: 신뢰도 50% 미만이면 권장사항?**
A: 자동 배포 중단, 수동 검토
   - 콜드 스타트: 인기도 기반 폴백
   - 데이터 부족: 임시 A/B 테스트
   - 모델 오류: 모델 재학습

---

**완료**: 2026-03-13 ✅

**다음**: Phase 5 Task 3 (마이크로서비스 아키텍처) / 또는 프로덕션 배포
