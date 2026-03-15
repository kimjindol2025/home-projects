# OAuth 2.0 Setup Guide

## 🔐 OAuth 2.0 소셜 로그인 설정

FreeLang에서 Google, GitHub, Naver 소셜 로그인을 사용하기 위한 완전한 가이드입니다.

---

## 📋 목차

1. [Google OAuth 설정](#google-oauth-설정)
2. [GitHub OAuth 설정](#github-oauth-설정)
3. [Naver OAuth 설정](#naver-oauth-설정)
4. [로컬 테스트](#로컬-테스트)
5. [프로덕션 배포](#프로덕션-배포)

---

## Google OAuth 설정

### 1단계: Google Cloud Console 접속

1. https://console.cloud.google.com 접속
2. 새 프로젝트 생성 (Project Name: "FreeLang")

### 2단계: OAuth 2.0 동의 화면 설정

1. **API 및 서비스** → **OAuth 동의 화면**
2. 사용자 유형: **외부** 선택
3. 다음 정보 입력:
   ```
   - 앱 이름: FreeLang
   - 사용자 지원 이메일: your-email@gmail.com
   - 개발자 연락처: your-email@gmail.com
   ```

### 3단계: OAuth 2.0 자격증명 생성

1. **자격증명** → **자격증명 만들기** → **OAuth 2.0 클라이언트 ID**
2. 애플리케이션 유형: **웹 애플리케이션**
3. 이름: "FreeLang Web Client"
4. 승인된 리다이렉션 URI 추가:
   ```
   http://localhost:3001/auth/google/callback
   https://253.dclub.kr/auth/google/callback
   ```
5. **만들기** 클릭 → 클라이언트 ID/Secret 복사

### 4단계: .env 파일 업데이트

```bash
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
```

---

## GitHub OAuth 설정

### 1단계: GitHub Settings 접속

1. https://github.com/settings/developers 접속
2. **OAuth Apps** → **New OAuth App**

### 2단계: 애플리케이션 등록

```
Application name: FreeLang
Homepage URL: https://253.dclub.kr
Authorization callback URL: https://253.dclub.kr/auth/github/callback
```

### 3단계: 자격증명 복사

- **Client ID** 복사
- **Generate Client Secret** 클릭 후 Secret 복사

### 4단계: .env 파일 업데이트

```bash
GITHUB_CLIENT_ID=your-client-id
GITHUB_CLIENT_SECRET=your-client-secret
```

---

## Naver OAuth 설정

### 1단계: Naver Developers 접속

1. https://developers.naver.com 접속
2. **로그인** → **애플리케이션 등록**

### 2단계: 애플리케이션 정보 입력

```
애플리케이션 이름: FreeLang
사용 API: Naver Login
환경: 웹 애플리케이션
```

### 3단계: 로그인 오픈 API 권한 설정

- **설정** 탭
- **Naver Login** 체크
- **Callback URL** 등록:
  ```
  https://253.dclub.kr/auth/naver/callback
  ```

### 4단계: 자격증명 복사

- **Client ID** 복사
- **Client Secret** 복사

### 5단계: .env 파일 업데이트

```bash
NAVER_CLIENT_ID=your-client-id
NAVER_CLIENT_SECRET=your-client-secret
```

---

## 로컬 테스트

### 1단계: 환경 설정

```bash
# .env 파일 생성
cp .env.example .env

# 필요한 OAuth 자격증명 입력
nano .env
```

### 2단계: 의존성 설치

```bash
npm install
```

### 3단계: 서버 실행

```bash
npm run dev
```

### 4단계: 로그인 테스트

```
http://localhost:3001
```

로그인 페이지에서:
- ✅ Google 로그인 버튼 클릭
- ✅ GitHub 로그인 버튼 클릭
- ✅ Naver 로그인 버튼 클릭

---

## 프로덕션 배포

### 1단계: GitHub Secrets 설정

```bash
# GitHub Repository Settings → Secrets and variables
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
NAVER_CLIENT_ID=xxx
NAVER_CLIENT_SECRET=xxx
JWT_SECRET=your-strong-random-secret
```

### 2단계: 배포

```bash
# 마스터 브랜치에 푸시 (자동 배포)
git push origin master
```

### 3단계: 검증

```bash
# 배포 상태 확인
curl https://253.dclub.kr/api/health

# 로그인 테스트
https://253.dclub.kr
```

---

## 🔒 보안 체크리스트

- ✅ JWT_SECRET은 강력한 랜덤 키 사용
- ✅ 모든 자격증명을 .env에 저장 (git에 커밋 금지)
- ✅ HTTPS 사용 (프로덕션)
- ✅ CORS 출처 제한
- ✅ 토큰 만료 시간 설정 (기본 7일)
- ✅ 정기적인 로그 모니터링

---

## 🐛 문제 해결

### "Invalid redirect URI" 오류

**해결**: OAuth 설정의 Callback URL이 코드와 일치하는지 확인
```javascript
// 코드에서 확인
callbackURL: `${baseURL}/auth/google/callback`
```

### "Missing client ID" 경고

**해결**: .env 파일에 모든 자격증명 입력
```bash
GOOGLE_CLIENT_ID=xxx
GITHUB_CLIENT_ID=xxx
NAVER_CLIENT_ID=xxx
```

### 로그인 후 토큰 저장 안 됨

**해결**: localStorage 활성화 확인
```javascript
localStorage.setItem('authToken', token);
```

---

## 📚 참고 자료

- [Google OAuth 공식 문서](https://developers.google.com/identity/protocols/oauth2)
- [GitHub OAuth 공식 문서](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [Naver OAuth 공식 문서](https://developers.naver.com/docs/login/overview)

---

## ✅ 완료 체크리스트

- [ ] Google OAuth 설정 완료
- [ ] GitHub OAuth 설정 완료
- [ ] Naver OAuth 설정 완료
- [ ] 로컬 테스트 완료
- [ ] .env 파일 보안 확인
- [ ] GitHub Secrets 설정 완료
- [ ] 프로덕션 배포 완료
