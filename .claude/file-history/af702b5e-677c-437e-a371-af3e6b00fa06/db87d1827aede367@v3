# 계획: KimGraf - FreeLang Native Monitoring Dashboard (Grafana 대체)

## Context

Phase 7 (Byzantine Resilience) + Sovereign-Mesh (5-Layer Network)가 253 서버에 배포된 상태에서,
이 시스템들을 모니터링할 Grafana 대체 시스템을 **FreeLang v2.2.0**으로 완전히 구현.

**목표**: "Grafana 없이 FreeLang으로 Production 모니터링 가능한 시스템"

---

## 프로젝트 구조

```
freelang-kimgraf/
├── src/kimgraf/
│   ├── core/
│   │   ├── tsdb.fl              (700줄) - Time-Series DB (ring buffer + file)
│   │   ├── query_engine.fl      (500줄) - 메트릭 쿼리 언어 (KimQL)
│   │   └── storage.fl           (400줄) - 파일 기반 영구 저장
│   ├── sources/
│   │   ├── prometheus.fl        (400줄) - Prometheus /metrics 스크래핑
│   │   ├── phase7_source.fl     (300줄) - Phase 7 Byzantine/Raft 메트릭
│   │   └── mesh_source.fl       (300줄) - Sovereign-Mesh 5-Layer 메트릭
│   ├── ui/
│   │   ├── dashboard.fl         (500줄) - 대시보드 레이아웃 시스템
│   │   ├── chart_engine.fl      (600줄) - Line/Bar/Gauge/Heatmap 차트
│   │   └── panel.fl             (400줄) - 패널 컴포넌트
│   ├── alerts/
│   │   ├── rule_engine.fl       (400줄) - 임계값 기반 알림 규칙
│   │   └── notifier.fl          (300줄) - Webhook/Email/Slack 알림
│   └── server.fl                (500줄) - HTTP + WebSocket 서버
│
├── tests/
│   └── kimgraf_tests.fl         (600줄) - 60개 무관용 테스트
│
├── Dockerfile                   (50줄)  - Docker 배포
└── KIMGRAF_COMPLETION.md        (200줄) - 완료 보고서
```

**총 목표**: ~5,000줄 FreeLang v2.2.0

---

## 8개 무관용 규칙

| 규칙 | 설명 | 목표 |
|-----|------|------|
| K1 | 메트릭 수집 주기 | < 1초 (1Hz) |
| K2 | 쿼리 응답 시간 | < 100ms |
| K3 | 대시보드 로딩 | < 500ms |
| K4 | WebSocket 지연 | < 50ms |
| K5 | 데이터 보존 | 30일 (86,400 샘플/메트릭) |
| K6 | 알림 지연 | < 5초 (임계값 감지 후) |
| K7 | 메모리 사용 | < 200MB |
| K8 | Prometheus 호환 | 100% (PromQL 서브셋) |

---

## 핵심 모듈 상세

### 1. core/tsdb.fl (700줄)
- TimeSeries 구조체 (ring buffer 기반)
- DataPoint: timestamp_ms + value
- insert(), query_range(), downsample()
- flush_to_disk() → K5 30일 보존

### 2. core/query_engine.fl (500줄)
- KimQL: Prometheus PromQL 서브셋 호환 (K8)
- parse_query(), execute(), rate(), avg_over_time()
- percentile() → P99, P95 계산

### 3. sources/phase7_source.fl (300줄)
- Phase 7 메트릭: byzantine_detected, election_latency, replication_rate
- fault_detection_ms, recovery_ms, production_latency_ms
- memory_usage_mb, circuit_breaker_trips

### 4. sources/mesh_source.fl (300줄)
- Sovereign-Mesh 메트릭: olsr_overhead, neural_relay_ms
- packet_loss_rate, anonymity_score, battery_usage, active_nodes

### 5. ui/chart_engine.fl (600줄)
- SVG 직접 생성 (D3.js 의존 없음)
- render_line_chart(), render_bar_chart()
- render_gauge(), render_heatmap(), render_stat_panel()

### 6. server.fl (500줄)
- HTTP 서버 (포트 9000) + WebSocket (포트 9001)
- /api/query, /api/metrics, /dashboard 엔드포인트
- Prometheus 호환 /metrics 엔드포인트 (K8)

### 7. alerts/rule_engine.fl (400줄)
- AlertRule: 조건 + 임계값 + 심각도
- evaluate_rules(), fire_alert() → K6 < 5s

---

## 기본 대시보드 레이아웃

```
┌─────────────────────────────────────────────────────────┐
│  KimGraf v1.0 - FreeLang Native Monitoring              │
├────────────────┬────────────────┬───────────────────────┤
│  Byzantine     │  Election      │  Log Replication       │
│  Tolerance     │  Latency       │  Rate                  │
│  [GAUGE: OK]   │  [LINE: 150ms] │  [GAUGE: 100%]         │
├────────────────┼────────────────┼───────────────────────┤
│  Production    │  Memory        │  Circuit Breaker       │
│  Latency P99   │  Usage         │  Status                │
│  [LINE: 45ms]  │  [BAR: 35MB]   │  [STATUS: CLOSED]      │
├────────────────┴────────────────┴───────────────────────┤
│  Sovereign-Mesh: OLSR Overhead / Anonymity / Battery    │
├────────────────────────────────────────────────────────-┤
│  Alerts: ✅ All systems operational                     │
└─────────────────────────────────────────────────────────┘
```

---

## 테스트 구조 (60개)

```
Group A: Core TSDB (10개)      - K2, K5
Group B: Query Engine (10개)   - K2, K8
Group C: Data Sources (10개)   - K1, K8
Group D: UI & Charts (10개)    - K3
Group E: Alerts (10개)         - K6
Group F: Server & E2E (10개)   - K4, K7
```

---

## 구현 순서

1. GOGS에 freelang-kimgraf 저장소 생성
2. src/kimgraf/core/ 3개 파일 구현
3. src/kimgraf/sources/ 3개 파일 구현
4. src/kimgraf/ui/ 3개 파일 구현
5. src/kimgraf/alerts/ 2개 파일 구현
6. src/kimgraf/server.fl 구현
7. tests/kimgraf_tests.fl (60개 테스트)
8. Dockerfile + KIMGRAF_COMPLETION.md
9. GOGS 커밋 + 253 서버 Docker 배포

---

## 기술 스택

- **언어**: FreeLang v2.2.0 (100%)
- **외부 의존도**: 0%
- **차트**: SVG 직접 생성 (D3.js 없음)
- **DB**: 파일 기반 TSDB (Prometheus DB 없음)
- **배포**: Dockerfile → 253 서버
- **포트**: 9000 (HTTP 대시보드), 9001 (WebSocket)
