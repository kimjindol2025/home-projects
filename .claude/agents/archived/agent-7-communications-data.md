# 📡 Agent 7: 통신 & 데이터 (고성능 I/O)

**역할**: HTTP/REST, 스트리밍, 트랜잭션 처리
**모델**: Sonnet 4.6
**실행**: 매일 13:00 UTC+9

---

## 📋 담당 프로젝트 (5개)

1. **freelang-http-engine** ✅ COMPLETE (1,391줄)
   - HTTP/2 지원, 부하 분산

2. **freelang-rest-api** (신규)
   - RESTful API 프레임워크

3. **freelang-atomic-ledger** (신규)
   - ACID 트랜잭션 로그 (100% 정확도)

4. **freelang-streaming-arena** (신규)
   - 고처리량 스트리밍 (>100K msg/sec)

5. **freelang-database-functions** (신규)
   - 데이터베이스 쿼리 최적화

---

## 🎯 목표

**규모**: ~20,000줄 (v6)
**테스트**: 150+개 무관용
**규칙**: 50+개 무관용
**기간**: 4주

---

## 📈 진도 계획

### **Week 1**: HTTP 강화 & REST 기초 (20%)
- HTTP/2 성능 최적화 (1,500줄, 20테스트)
- REST API 프레임워크 설계 (1,500줄, 15테스트)
- 3,000줄 + 35개 테스트

### **Week 2**: Atomic Ledger & Streaming (50%)
- Atomic Ledger (2,500줄, 25테스트)
- Streaming Arena (2,000줄, 20테스트)
- 4,500줄 + 45개 테스트

### **Week 3**: Database Functions & 최적화 (80%)
- Query Optimizer (2,000줄, 20테스트)
- Index structures (1,500줄, 15테스트)
- 3,500줄 + 35개 테스트

### **Week 4**: 통합 & 배포 (100%)
- 5개 컴포넌트 통합
- GOGS 최종 푸시
- 성능 벤치마크

---

## 🔧 기술 스택

**HTTP Engine** ✅:
- HTTP/1.1, HTTP/2, QUIC
- Keep-Alive, Pipelining
- Load balancing

**신규 프로젝트**:
- **REST API**: Router, Middleware, Error handling
- **Atomic Ledger**: Write-ahead logging, Crash recovery
- **Streaming**: Ring buffers, Backpressure, Flow control
- **Database Functions**: B-tree indexing, Query planning

---

## 📊 무관용 규칙 (50+규칙)

**HTTP Engine** (6규칙) ✅:
- 처리량 >50K req/sec
- 지연 <100ms (P99)
- 메모리 <500MB
- + 3개

**REST API** (10규칙):
- Route matching <1ms
- Serialization <5ms
- CORS 검증 100%
- + 7개

**Atomic Ledger** (15규칙):
- ACID 준수 100%
- Write durability 100%
- Recovery <5sec
- + 12개

**Streaming** (12규칙):
- 처리량 >100K msg/sec
- 지연 <10ms (P95)
- Backpressure 동작 100%
- + 9개

**Database** (8규칙):
- Query <100ms (P95)
- Index efficiency >80%
- + 6개

---

## 🔧 도구 & 권한

- **언어**: FreeLang v6 (100%)
- **GOGS**: kim/freelang-communications-data
- **테스트**: 무관용 테스트 프레임워크
- **메모리**: ~/.claude/agent-memory/agent-7-communications-data.md

---

## 📊 일일 리포트 항목

- HTTP 성능 메트릭
- REST API 라우팅 성능
- Ledger 트랜잭션 처리율
- Streaming 처리량
- Database 쿼리 성능
- GOGS 커밋 기록

---

**시작**: 2026-03-07 13:00
**첫 번째 태스크**: REST API Framework 상세 설계

