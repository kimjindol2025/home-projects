# 📋 TODO: FreeLang 자체호스팅 완성

## 🎯 Stage 2: 부트스트랩 체인 (6주)

### Week 1-2: 실행 환경 준비

- [ ] QEMU x86_64 설치
- [ ] x86-64 시뮬레이션 환경 구성
- [ ] ELF 바이너리 실행 테스트
- [ ] v1.elf 실행 가능 확인

**예상 시간**: 1-2주

---

### Week 3-4: FreeLang Runtime 구현

**필수 구현**:
- [ ] malloc / free 함수
- [ ] syscall 인터페이스
- [ ] 문자열 함수
- [ ] 배열 연산

**테스트**:
- [ ] Runtime 단위 테스트 (10개 이상)
- [ ] 통합 테스트

**예상 시간**: 2주

---

### Week 5-6: Parser 확장 & 통합

**struct 문법 추가**:
- [ ] struct 키워드 인식
- [ ] 멤버 변수 파싱
- [ ] 인스턴스 생성
- [ ] 멤버 접근 (.dot notation)

**컴파일러 통합**:
- [ ] Lexer 업데이트
- [ ] Parser 업데이트
- [ ] Code Generator 업데이트
- [ ] 테스트 (5개 이상)

**예상 시간**: 2주

---

## 🚀 Stage 3: 자체호스팅 증명 (4주)

### Week 1-2: v2 생성

```bash
# v1.elf 실행
./v1.elf freelang-compiler-tiny.fl > v2.elf

# MD5 검증
md5sum v1.elf v2.elf
# 결과: 동일해야 함 ✅
```

**체크리스트**:
- [ ] v1.elf 성공적으로 실행
- [ ] v2.elf 생성됨
- [ ] MD5(v1) == MD5(v2) 확인

---

### Week 3-4: v3 생성 & 최종 검증

```bash
# v2.elf 실행
./v2.elf freelang-compiler-tiny.fl > v3.elf

# MD5 검증
md5sum v2.elf v3.elf
# 결과: 동일해야 함 ✅
```

**최종 검증**:
- [ ] v2.elf 성공적으로 실행
- [ ] v3.elf 생성됨
- [ ] MD5(v2) == MD5(v3) 확인
- [ ] MD5(v1) == MD5(v2) == MD5(v3)
- [ ] **자체호스팅 완전 증명 ✅**

---

## 📊 진행도 추적

| 단계 | 상태 | 완료율 |
|------|------|--------|
| Stage 1: 결정론적 컴파일 | ✅ 완료 | 100% |
| Stage 2: 환경 준비 | ⏳ 대기 | 0% |
| Stage 2: Runtime | ⏳ 대기 | 0% |
| Stage 2: Parser | ⏳ 대기 | 0% |
| Stage 3: v2 생성 | ⏳ 대기 | 0% |
| Stage 3: v3 생성 | ⏳ 대기 | 0% |
| **전체** | ⏳ 진행중 | **14%** |

---

## 🔍 품질 관리

**코드 검증**:
- [ ] 모든 코드는 테스트 가능해야 함
- [ ] 벤치마크 기록 (성능 퇴화 방지)
- [ ] 거짓 주장 금지 (실행 증거 필수)

**문서화**:
- [ ] 각 단계마다 PHASE_X_COMPLETE.md 작성
- [ ] 성능 벤치마크 기록
- [ ] 블로커 발생 시 즉시 기록

**git 관리**:
- [ ] 매일 커밋
- [ ] 명확한 커밋 메시지
- [ ] GOGS 푸시 (https 또는 SSH)

---

## 📞 참고 문서

| 문서 | 내용 |
|------|------|
| [MEMORY.md](./projects/-data-data-com-termux-files-home/memory/MEMORY.md) | 핵심 메모리 (1000자) |
| [PHASE_STATUS.md](./PHASE_STATUS.md) | Phase 진행 상태 |
| [CHECKLIST.md](./CHECKLIST.md) | 2중 검증 체크리스트 |
| [SHORTCUTS.md](./SHORTCUTS.md) | 접속정보 & 명령어 |
| [TODO.md](./TODO.md) | 이 파일 |

---

**마지막 업데이트**: 2026-03-08
**예상 완료**: 2026-04-19 (10주)
