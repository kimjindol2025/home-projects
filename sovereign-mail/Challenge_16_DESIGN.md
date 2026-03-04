# Challenge 16: L0NN-Mail-Sentry (AI 지능형 방어막)
## Neural Network-Based Real-Time Spam/Phishing Detection

**Status**: 🚀 DESIGN PHASE
**Target Date**: 2026-03-25 (1주)
**Scope**: 4개 모듈, 2,300줄 코드, 30개 테스트, 6개 무관용 규칙

---

## 📋 Challenge Objectives

### Primary Goals
1. **99.9% 정확도 스팸/피싱 필터**
   - Phase 8의 L0NN 확장 (신경망)
   - 기존 블랙리스트 방식 제거
   - 실시간 위협 탐지 (<10µs)

2. **8가지 특성 추출**
   - 발신자 신용도 (DHT)
   - 헤더 이상 탐지 (SPF/DKIM/DMARC)
   - 콘텐츠 패턴 분석
   - URL 피싱 탐지
   - 첨부파일 악성코드 스캔
   - 메시지 요청 빈도 (brute-force)
   - 언어 모델 분석 (GPT 생성 탐지)
   - 암호화 서명 검증

3. **네트워크 계층 통합**
   - 메일 라우팅 단계에서 필터
   - 악의 메일 차단 (조기)
   - 폭발 반경 차단 (발신자 IP)

4. **분산 학습**
   - 사용자 스팸 신고 피드백
   - 신경망 온라인 학습
   - DHT 평판 업데이트

### Innovation Points
- **신경망 가속**: <10µs 추론 (L0NN)
- **8가지 복합 특성**: 단순한 키워드 필터 아님
- **실시간 탐지**: 메일 입장 순간 판정
- **분산 학습**: 사용자 신고로 신경망 개선

### Success Criteria
- ✅ 오탐률 (False Positive) <0.01%
- ✅ 미탐률 (False Negative) <0.5%
- ✅ 분석 속도 <10µs (per message)
- ✅ 신경망 정확도 ≥99.9%
- ✅ 평판 업데이트 <1초
- ✅ 메모리 오버헤드 <50MB

---

## 🏗️ Architecture

### 4-Module L0NN-Mail-Sentry Stack

#### Module 1: Mail-Classifier (800 lines)
**Purpose**: 메일 분류 (정상/스팸/피싱/악성)

```
메일 입장
  ↓ [Mail Classifier]
  ├─ 분류 유형
  │   ├─ LEGITIMATE (정상 메일)
  │   ├─ SPAM (스팸, 광고)
  │   ├─ PHISHING (피싱, 사기)
  │   └─ MALWARE (악성코드, 감염)
  ├─ 분류 경로
  │   ├─ 특성 추출 (8가지)
  │   ├─ 신경망 입력
  │   ├─ 신경망 추론 (<10µs)
  │   └─ 점수 계산 (0-100)
  ├─ 판정 기준
  │   ├─ Score ≥ 80 → ACCEPT
  │   ├─ Score 40-80 → QUARANTINE
  │   └─ Score < 40 → REJECT
  └─ 실행 액션
      ├─ ACCEPT: 메일함에 저장
      ├─ QUARANTINE: 스팸 폴더
      └─ REJECT: 폐기 + 신고
```

**Components**:
- `MailClassifier`: 메일 분류 엔진
- `ScoringEngine`: 위험도 점수 계산
- `ActionExecutor`: 수락/격리/거부
- `FeedbackCollector`: 사용자 신고

**Techniques**:
- Probabilistic classification (Bayesian)
- Confidence scoring
- Multi-class classification (4개 클래스)
- Real-time execution

**Unforgiving Rules (Module 1)**:
- **Rule 1**: 오탐률 <0.01%
- **Rule 2**: 미탐률 <0.5%

---

#### Module 2: Feature-Extractor (600 lines)
**Purpose**: 메일에서 8가지 특성 추출

```
메일 분석
  ↓ [Feature Extractor]
  ├─ Feature 1: 발신자 신용도 (DHT)
  │   ├─ kim@sovereign의 DHT 평판 조회
  │   ├─ 신용도 점수 (0-100)
  │   └─ 블랙리스트 여부 확인
  ├─ Feature 2: 헤더 이상 탐지
  │   ├─ SPF 레코드 검증 (발신자 도메인)
  │   ├─ DKIM 서명 검증
  │   └─ DMARC 정책 확인
  ├─ Feature 3: 콘텐츠 패턴
  │   ├─ 의심 키워드 (결제 요청, 긴급, 클릭)
  │   ├─ HTML 비율 (피싱 = 높음)
  │   └─ 언어 통계 (자연언어 vs 생성)
  ├─ Feature 4: URL 피싱 탐지
  │   ├─ URL 숫자 (피싱 = 많음)
  │   ├─ 단축 URL 사용 여부
  │   └─ 도메인 신용도 (VirusTotal 등)
  ├─ Feature 5: 첨부파일 스캔
  │   ├─ 파일 형식 (exe, zip 등)
  │   ├─ 파일 크기 (비정상적 = 의심)
  │   └─ 바이러스 스캔 (clamav 호환)
  ├─ Feature 6: 메시지 빈도
  │   ├─ 발신자가 1시간에 몇 개 보냈나
  │   ├─ 같은 시간대 수신한 메시지 수
  │   └─ 점수: 높음(의심) = 낮은 점수
  ├─ Feature 7: 언어 모델 (GPT 탐지)
  │   ├─ 텍스트 엔트로피 (생성 = 낮음)
  │   ├─ 반복 패턴 분석
  │   └─ 자연스러움 점수
  └─ Feature 8: 암호화 검증
      ├─ 메일 서명 검증 (DKIM, S/MIME)
      ├─ 발신자 인증 확인
      └─ 전송 암호화 (TLS) 확인
```

**Components**:
- `ReputationFeature`: DHT 조회
- `HeaderAnalyzer`: SPF/DKIM/DMARC
- `ContentAnalyzer`: 텍스트 분석
- `URLAnalyzer`: 링크 위험도
- `AttachmentScanner`: 파일 분석
- `FrequencyAnalyzer`: 빈도 분석
- `LanguageModelAnalyzer`: 자연언어 분석
- `CryptoVerifier`: 서명 검증

**Techniques**:
- TF-IDF (텍스트 특성화)
- Entropy calculation (언어 분석)
- Rate limiting (빈도 분석)
- Cryptographic verification

**Unforgiving Rules (Module 2)**:
- **Rule 3**: 분석 속도 <10µs
- **Rule 6**: 메모리 <50MB

---

#### Module 3: Neural-Network (700 lines)
**Purpose**: 신경망 모델 (Phase 8 L0NN 확장)

```
신경망 구조
  ↓ [Neural Network]
  ├─ 입력층 (8개 특성)
  │   ├─ Feature 1: 발신자 신용도 (정규화)
  │   ├─ Feature 2: 헤더 이상 (0/1)
  │   ├─ Feature 3: 콘텐츠 위험도 (정규화)
  │   ├─ Feature 4: URL 위험도 (정규화)
  │   ├─ Feature 5: 첨부파일 위험도 (정규화)
  │   ├─ Feature 6: 메시지 빈도 (정규화)
  │   ├─ Feature 7: 생성 메일 여부 (0/1)
  │   └─ Feature 8: 암호화 검증 (0/1)
  ├─ 은닉층 1 (64개 노드)
  │   ├─ ReLU 활성화
  │   ├─ 배치 정규화
  │   └─ 드롭아웃 0.3
  ├─ 은닉층 2 (32개 노드)
  │   ├─ ReLU 활성화
  │   ├─ 배치 정규화
  │   └─ 드롭아웃 0.3
  ├─ 출력층 (4개 클래스)
  │   ├─ LEGITIMATE (정상)
  │   ├─ SPAM (스팸)
  │   ├─ PHISHING (피싱)
  │   └─ MALWARE (악성코드)
  └─ 추론
      ├─ Forward pass (<10µs)
      ├─ Softmax 확률화
      ├─ 최대 확률 클래스 선택
      └─ 신뢰도 점수 반환
```

**Components**:
- `Layer`: 신경망 계층
- `NeuralNetwork`: 전체 신경망
- `ForwardPass`: 추론 실행
- `Optimizer`: SGD/Adam (학습)

**Techniques**:
- 3-layer architecture (Phase 8 L0NN)
- Batch normalization
- Dropout regularization
- Softmax output

**Unforgiving Rules (Module 3)**:
- **Rule 4**: 정확도 ≥99.9%
- **Rule 5**: 업데이트 <1초

---

#### Module 4: Reputation-Tracker (500 lines)
**Purpose**: 발신자 평판 추적 및 신경망 학습

```
평판 추적
  ↓ [Reputation Tracker]
  ├─ 사용자 신고
  │   ├─ "이 메일은 스팸이다" → -20점
  │   ├─ "이 메일은 피싱이다" → -30점
  │   ├─ "이 메일은 정상이다" → +5점
  │   └─ 신고는 서명되어 DHT에 저장
  ├─ 신경망 학습
  │   ├─ 매시간 학습 (온라인 학습)
  │   ├─ 사용자 신고를 라벨로 사용
  │   ├─ SGD로 가중치 업데이트
  │   └─ 모델 정확도 향상
  ├─ 피드백 루프
  │   ├─ 사용자 신고 수집
  │   ├─ 신경망 재학습
  │   ├─ 모델 배포
  │   └─ 정확도 모니터링
  └─ 폭발 반경 차단
      ├─ 스팸 발신자 IP 기록
      ├─ ASN (자율 시스템) 누적
      ├─ 같은 IP 대역의 메일 감시
      └─ 높은 위험도 자동 할당
```

**Components**:
- `UserFeedback`: 사용자 신고 수집
- `ModelTrainer`: 온라인 학습
- `ReputationUpdate`: DHT 평판 업데이트
- `BlastRadiusDetector`: IP/ASN 분석

**Techniques**:
- Online learning (incremental)
- Exponential decay (오래된 신고는 영향 줄임)
- Collaborative filtering (다른 사용자 신고)
- CIDR block analysis (IP 대역)

---

## 🧪 Test Plan (30+ tests)

### Group A: Mail-Classifier (10 tests)
```
✓ test_mail_classification           - 메일 분류
✓ test_legitimate_mail               - 정상 메일 탐지
✓ test_spam_detection                - 스팸 탐지
✓ test_phishing_detection            - 피싱 탐지
✓ test_malware_detection             - 악성 메일 탐지
✓ test_score_calculation             - 위험도 점수
✓ test_false_positive_rate           - 오탐률 <0.01%
✓ test_false_negative_rate           - 미탐률 <0.5%
✓ test_action_execution              - 수락/격리/거부 실행
✓ test_feedback_collection           - 신고 수집
```

### Group B: Feature-Extractor (8 tests)
```
✓ test_reputation_feature            - 신용도 조회 (DHT)
✓ test_header_analysis               - SPF/DKIM/DMARC 검증
✓ test_content_analysis              - 콘텐츠 패턴 분석
✓ test_url_analysis                  - URL 피싱 탐지
✓ test_attachment_scan               - 첨부파일 스캔
✓ test_frequency_analysis            - 빈도 분석
✓ test_language_model_analysis       - 생성 메일 탐지
✓ test_crypto_verification           - 서명 검증
```

### Group C: Neural-Network (6 tests)
```
✓ test_forward_pass                  - 신경망 추론
✓ test_inference_speed               - <10µs 추론 (Rule 4)
✓ test_output_softmax                - 확률 정규화
✓ test_model_accuracy                - ≥99.9% 정확도
✓ test_gradient_computation          - 역전파 계산
✓ test_model_update                  - 가중치 업데이트
```

### Group D: Reputation-Tracker (6 tests)
```
✓ test_user_feedback_handling        - 신고 수집
✓ test_model_retraining              - 온라인 학습
✓ test_reputation_update             - DHT 평판 업데이트
✓ test_update_latency                - <1초 업데이트
✓ test_blast_radius_detection        - IP/ASN 분석
✓ test_collaborative_filtering       - 협력 필터링
```

---

## 📊 Unforgiving Rules (6 total)

| Rule | Target | Implementation |
|------|--------|-----------------|
| **R1** | False Positive <0.01% | MailClassifier::classify() |
| **R2** | False Negative <0.5% | MailClassifier::verify_accuracy() |
| **R3** | Analysis speed <10µs | FeatureExtractor::extract_all() |
| **R4** | Accuracy ≥99.9% | NeuralNetwork::forward_pass() |
| **R5** | Update latency <1s | ReputationTracker::update_model() |
| **R6** | Memory <50MB | SystemMonitor::get_memory_usage() |

---

## 📁 File Structure

```
Challenge_16/
├── src/
│   ├── mail_classifier.rs           (800 lines)
│   │   ├── MailClassifier
│   │   ├── ScoringEngine
│   │   ├── ActionExecutor
│   │   └── [10 test functions]
│   │
│   ├── feature_extractor.rs         (600 lines)
│   │   ├── ReputationFeature
│   │   ├── HeaderAnalyzer
│   │   ├── ContentAnalyzer
│   │   ├── URLAnalyzer
│   │   ├── AttachmentScanner
│   │   ├── FrequencyAnalyzer
│   │   └── [8 test functions]
│   │
│   ├── neural_network.rs            (700 lines)
│   │   ├── Layer
│   │   ├── NeuralNetwork
│   │   ├── ForwardPass
│   │   ├── Optimizer
│   │   └── [6 test functions]
│   │
│   ├── reputation_tracker.rs        (500 lines)
│   │   ├── UserFeedback
│   │   ├── ModelTrainer
│   │   ├── ReputationUpdate
│   │   ├── BlastRadiusDetector
│   │   └── [6 test functions]
│   │
│   └── lib.rs                       (updated)
│
├── Cargo.toml
├── DESIGN.md                        (this file)
└── COMPLETION_REPORT.md             (to be written)
```

---

## 🎯 Implementation Strategy

### Phase A: Mail Classification (Days 1-2)
1. 분류 엔진 구현
2. 위험도 점수 계산
3. 수락/격리/거부 액션
4. 신고 수집

### Phase B: Feature Extraction (Day 3)
1. 8가지 특성 추출
2. DHT 평판 조회
3. SPF/DKIM/DMARC 검증
4. 파일 및 URL 분석

### Phase C: Neural Network (Day 4)
1. 3-layer 신경망 구현
2. Forward pass (<10µs)
3. Batch normalization
4. Softmax output

### Phase D: Reputation Tracking (Day 5)
1. 사용자 신고 수집
2. 온라인 학습 (<1초)
3. DHT 평판 업데이트
4. IP/ASN 블래스트 반경

### Phase E: Integration & Testing (Day 6)
1. 모든 모듈 통합
2. <10µs 분석 속도 검증
3. 99.9% 정확도 검증
4. 메모리 <50MB 확인

---

## 📈 Expected Performance

**Accuracy**:
- False Positive: <0.01% (99.99% specificity)
- False Negative: <0.5% (99.5% sensitivity)
- Overall Accuracy: ≥99.9%

**Speed**:
- Feature extraction: <5µs
- Neural network inference: <3µs
- Total per-message: <10µs

**Resource**:
- Memory: <50MB (model + cache)
- Network: <10KB (DHT reputation lookup)

**Scalability**:
- Handles 1M+ emails/day
- Online learning updates real-time
- Distributed reputation system

---

## 🏆 Final Outcome

```
Project Sovereign-Mail (완성도 예상)

Challenge 14: L0-Mail-Core
├─ 메모리에서 암호화 (<5ms)
├─ CAS 자동 변환 (검은 상자)
└─ RFC 4880 OpenPGP 준수

Challenge 15: Sovereign-Naming
├─ DHT 기반 분산 주소 (<100ms)
├─ 공개키 주소화 (영구 불멸)
└─ 평판 시스템 (스팸 차단)

Challenge 16: L0NN-Mail-Sentry
├─ 99.9% 정확도 스팸 필터
├─ <10µs 신경망 추론
└─ 분산 학습 (사용자 신고)

= Project Sovereign-Mail ✅
  진정한 개인 메일 서버 완성
  구글/MS 클라우드 의존 종료
  "메일은 이제 내 것이다"
```

---

**Next Step**: Implement L0NN-Mail-Sentry → All 3 Challenges complete → Final integration

**Status**: Design approved, ready for implementation 🔧

**Philosophy**: "99.9% 정확도로 악의 메일을 차단하되, 단 1개의 정상 메일도 놓치지 않는다."
