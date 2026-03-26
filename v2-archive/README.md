# 📦 V2 Archive - 참고 코드 및 레거시

**위치**: `/home/v2-archive/`
**목적**: FreeLang V2 및 FV2 관련 참고 코드 보관
**상태**: 📦 보관 중 (필요시 언제든 접근 가능)

---

## 📁 폴더 구조

```
v2-archive/
├─ freelang-v2/           FreeLang V2 (Go 구현, 레거시)
├─ fv2-lang/              FV2 Rust 구현
└─ v2-freelang-ai/        V2 AI 통합 프로젝트
```

---

## 📝 각 폴더 설명

### 1. freelang-v2/
- **언어**: Go
- **크기**: ~7.2GB
- **상태**: 레거시 (아카이브됨)
- **용도**: 참고용 (FreeLang의 이전 버전)
- **마지막 수정**: 2026-03-21

### 2. fv2-lang/
- **언어**: Rust
- **크기**: ~15MB
- **상태**: 보관 (Rust 구현)
- **용도**: Rust 버전 참고
- **마지막 수정**: 2026-03-22

### 3. v2-freelang-ai/
- **언어**: 혼합
- **크기**: ~25MB
- **상태**: 보관
- **용도**: V2 AI 통합 프로젝트
- **마지막 수정**: 루트에서 이동됨

---

## 🔗 관련 링크

**활성 프로젝트**:
- `projects/fv2-lang-go/` - 🌟 최신 Go 구현 (100% 완성도)
- `.claude/projects/fv2-lang-go/` - 개발 중

**백업**:
- `dev/archived/v2-backup/` - 통합 백업

---

## ⚡ 사용 방법

```bash
# 참고 코드 확인
cd v2-archive/freelang-v2/
git log --oneline

# Rust 버전 참고
cd v2-archive/fv2-lang/
cat src/main.rs
```

---

## ✨ 정리 내역

**2026-03-26**:
- ✅ `projects/freelang-v2` 이동
- ✅ `projects/fv2-lang` 이동
- ✅ 루트 `v2-freelang-ai` 이동
- ✅ v2-archive 통합 완료

**결과**:
- projects/ 폴더 간소화 (핵심 프로젝트만)
- 참고 코드 한 곳에 집중 (v2-archive)
- 루트 정리 완료

---

## 📌 주의사항

- 이 폴더는 **읽기 전용** 권장
- 수정 필요시 백업 후 진행
- 대용량 폴더 (freelang-v2: 7.2GB)

---

**상태**: ✅ 정리 완료 (2026-03-26)
**관리자**: V2 Cleanup Initiative
