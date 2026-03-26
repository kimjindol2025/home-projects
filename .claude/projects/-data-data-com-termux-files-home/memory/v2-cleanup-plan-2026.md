---
name: V2 Projects Cleanup Plan 2026
description: 16개의 v2 폴더 정리 및 조직화 계획 (90MB 회수 가능)
type: project
---

# 📂 V2 프로젝트 폴더 정리 계획

**날짜**: 2026-03-26
**현황**: 16개의 v2 폴더 산재 상태
**목표**: 명확한 구조, 약 90MB 회수

---

## 🔍 폴더 분석 결과

### 다른 프로젝트 확인됨 ✅

```
1. projects/freelang-v2
   - FreeLang 버전 2 (Go 구현)
   - 크기: 7.2GB
   - 상태: 완전한 Git 저장소, AI 백업 포함
   - 결정: ✅ 보관 (명확히 레이블)

2. projects/fv2-lang
   - FV2 Rust 구현
   - 크기: ~15MB
   - 상태: 간단한 구조
   - 결정: ✅ 보관 (Rust 버전)

3. projects/fv2-lang-go
   - FV2 Go 구현 (최신)
   - 크기: ~50MB
   - 상태: 100% 완성도, 활성
   - 결정: ✅ 보관 (메인 프로젝트)
```

---

## 📊 크기 분석

### 활성 프로젝트 (~140MB)
- projects/fv2-lang-go: 50MB (보관)
- projects/freelang-v2: 7.2GB (보관, 별도)
- projects/fv2-lang: 15MB (보관)
- .claude/projects/fv2-lang-go: 40MB (개발 중)
- v2-freelang-ai (루트): 25MB (정리 필요)

### 아카이브 프로젝트 (~90MB) - 삭제 가능
- dev/archived/* : 70MB
- .projects/archived/* : 20MB

**💰 회수 가능 공간: ~90MB**

---

## 🎯 정리 계획 (3단계)

### 📍 1단계: 레이블링 (10분)

```
projects/
├─ freelang-v2/              (← ARCHIVE: Go 버전 2)
├─ fv2-lang/                 (← ARCHIVE: Rust 구현)
└─ fv2-lang-go/              (← ACTIVE: Go 구현)
```

**Action:**
- [ ] README-ARCHIVE.md 추가 (각 폴더에)
- [ ] 폴더명 앞에 [v2-archive] 또는 [ACTIVE] 표시

### 📍 2단계: 루트 정리 (10분)

```
❌ /data/data/com.termux/files/home/v2-freelang-ai/
✅ → projects/v2-freelang-ai/ 또는 .projects/core/v2-freelang-ai/
```

**Action:**
- [ ] v2-freelang-ai 내용 확인
- [ ] 적절한 위치로 이동
- [ ] 루트 정리

### 📍 3단계: 아카이브 통합 (20분)

```
dev/archived/
├─ v2-archive/              (← 새 폴더)
│  ├─ freelang-v2/          (복사)
│  ├─ fv2-lang/             (복사)
│  └─ fv2-lang-go/          (복사)
└─ (기존 파일들 삭제)
```

**Action:**
- [ ] dev/archived/v2-archive 폴더 생성
- [ ] 기존 v2 폴더들 이동
- [ ] 기존 아카이브 정리
- [ ] 백업 확인 후 삭제

### 📍 4단계: .projects 정리 (10분)

```
.projects/archived/
├─ v2-projects/            (← 통합)
│  ├─ freelang-v2/         (요약 후 보관)
│  └─ fv2-checks/          (통합)
└─ (불필요한 것 삭제)
```

**Action:**
- [ ] .projects/core/freelang-v2-* 파일 요약
- [ ] 통합 아카이브로 이동
- [ ] 중복 삭제

---

## 📋 체크리스트

### 우선순위 1️⃣ (긴급)
- [ ] projects 폴더 레이블링
- [ ] README 추가 (각 폴더)
- [ ] v2-freelang-ai (루트) 처리

### 우선순위 2️⃣ (이번 주)
- [ ] dev/archived 통합
- [ ] .projects 정리
- [ ] 백업 확인

### 우선순위 3️⃣ (다음 주)
- [ ] 불필요한 파일 안전 삭제
- [ ] 최종 검증

---

## 💾 백업 전략

**삭제 전 백업:**
```bash
# 전체 백업
tar -czf /backup/v2-projects-backup-$(date +%Y%m%d).tar.gz \
  dev/archived/* .projects/archived/*
```

**검증:**
- [ ] 백업 파일 크기 확인 (~90MB)
- [ ] 무결성 테스트 (tar -tzf)
- [ ] 별도 저장소에 보관

---

## 📊 예상 결과

### Before (현재)
```
v2 폴더: 16개 산재
구조: 혼란스러움
디스크: 낭비 (~90MB)
```

### After (정리 후)
```
projects/
  ├─ [v2-archive] freelang-v2/    (명확한 레이블)
  ├─ [v2-archive] fv2-lang/        (명확한 레이블)
  └─ [ACTIVE] fv2-lang-go/         (활성 프로젝트)

dev/archived/
  └─ v2-projects/                  (통합 아카이브)

.projects/
  └─ archived/                      (정리됨)

회수 가능 공간: ~90MB ✅
```

---

## 🔗 관련 파일

- MEMORY.md (이 메모리 파일)
- marketing-2026-03-26.md (마케팅 성과)

---

## 📝 노트

**Why:** v2 프로젝트들이 여러 위치에 산재되어 있어 관리가 어려움

**How:** 명확한 레이블링과 통합 아카이브를 통한 구조화

**Impact:**
- 명확한 프로젝트 구조 ✅
- 디스크 공간 회수 (90MB) ✅
- 유지보수 용이 ✅
- 새로운 멤버 온보딩 간편 ✅

---

**상태**: 계획 수립 완료 ✅
**다음 단계**: 1단계 시작 (레이블링)
**예상 완료**: 2026-03-27 또는 그 이후
