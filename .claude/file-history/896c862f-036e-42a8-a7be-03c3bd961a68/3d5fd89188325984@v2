# FreeLang Light Architecture & Security

## 🏗️ System Architecture

### OAuth 2.0 Complete Flow

```mermaid
sequenceDiagram
    participant User as 사용자 (브라우저)
    participant Frontend as Frontend (React)
    participant AuthServer as OAuth Provider<br/>(Google/GitHub/Naver)
    participant Backend as Backend Server<br/>(Node.js + Express)
    participant Database as Session Store<br/>(Memory/Redis)

    User->>Frontend: 1. "Google로 로그인" 클릭
    Frontend->>Backend: 2. GET /auth/google
    Backend->>AuthServer: 3. 리다이렉트 to OAuth 로그인 페이지
    AuthServer->>User: 4. 로그인 & 권한 동의 페이지
    User->>AuthServer: 5. 권한 동의
    AuthServer->>Backend: 6. 리다이렉트 /auth/google/callback?code=AUTH_CODE
    Backend->>AuthServer: 7. POST 토큰 교환 (code → access_token)
    AuthServer->>Backend: 8. access_token 반환
    Backend->>Backend: 9. 사용자 정보 조회 & 토큰 페어 생성
    Backend->>Database: 10. 단기 인증 코드 저장 (60초 유효)
    Backend->>Frontend: 11. 리다이렉트 /auth/callback?code=AUTH_CODE
    Frontend->>Backend: 12. POST /auth/exchange?code=AUTH_CODE
    Backend->>Backend: 13. 인증 코드 검증 & 제거
    Backend->>Frontend: 14. accessToken + refreshToken (HttpOnly Cookie)
    Frontend->>Frontend: 15. localStorage에 accessToken 저장
    User->>User: ✅ 로그인 완료
```

### Token Refresh Flow

```mermaid
sequenceDiagram
    participant User as 사용자
    participant Frontend as Frontend (React)
    participant Backend as Backend Server

    User->>Frontend: 1. 보호된 리소스 요청
    Frontend->>Backend: 2. GET /api/admin/stats<br/>Header: Authorization: Bearer {accessToken}

    alt accessToken이 유효한 경우
        Backend->>Backend: 3. 토큰 검증 ✅
        Backend->>Frontend: 4. 데이터 반환 (200 OK)
    else accessToken이 만료된 경우
        Backend->>Frontend: 3. 만료 에러 (401 Unauthorized)
        Frontend->>Backend: 4. POST /auth/refresh<br/>Cookie: refreshToken
        Backend->>Backend: 5. refreshToken 검증

        alt refreshToken이 유효한 경우
            Backend->>Backend: 6. 새로운 accessToken 생성
            Backend->>Frontend: 7. 새 accessToken 반환 (200 OK)
            Frontend->>Backend: 8. 원래 요청 재시도<br/>(새 accessToken 사용)
            Backend->>Frontend: 9. 데이터 반환
        else refreshToken이 만료된 경우
            Backend->>Frontend: 6. 갱신 실패 (401 Unauthorized)
            Frontend->>User: 7. 재로그인 필요 안내
        end
    end
```

### Security Layer Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      External Requests                          │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 1: Security Headers (helmet)                             │
│  ├─ HSTS (HTTP Strict Transport Security)                       │
│  ├─ X-Frame-Options (Clickjacking Protection)                   │
│  ├─ X-Content-Type-Options (MIME Sniffing Prevention)           │
│  └─ Content-Security-Policy (XSS Prevention)                    │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 2: CORS & Rate Limiting                                  │
│  ├─ CORS: credentials + origin validation                       │
│  ├─ Auth Rate Limit: 10 requests/15min                          │
│  └─ API Rate Limit: 100 requests/min                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 3: Input Validation (Joi)                                │
│  ├─ Query Parameters: code (alphanum, 10-500 chars)             │
│  ├─ Body Parameters: title, content validation                  │
│  └─ Auto-cleanup: stripUnknown=true                             │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 4: Authentication (JWT)                                  │
│  ├─ Access Token: 15 minutes (Authorization header)             │
│  ├─ Refresh Token: 30 days (HttpOnly cookie)                    │
│  ├─ Token Type: access/refresh (type field)                     │
│  └─ Secret Validation: Throws error if not set                  │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 5: Business Logic                                        │
│  ├─ Single-use Auth Code (60sec expiry)                         │
│  ├─ Session Management (7 days)                                 │
│  └─ Route Handlers                                              │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│  Layer 6: Response Security                                     │
│  ├─ No sensitive data in responses                              │
│  ├─ No JWT in URL query parameters                              │
│  ├─ Timestamp for replay attack prevention                      │
│  └─ Error messages without stack traces                         │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│                    Database / Storage                            │
│              (PostgreSQL / Redis / Memory)                       │
└─────────────────────────────────────────────────────────────────┘
```

## 🔐 Security Improvements (Phase 1-3)

### Phase 1: Security Foundation
- ✅ Helmet middleware for security headers
- ✅ Express rate limiting (Auth: 10/15min, API: 100/min)
- ✅ Cookie parser for HttpOnly cookie handling
- ✅ JWT_SECRET validation (throws error if not set)

### Phase 2: JWT Enhancement
- ✅ Token pair generation (Access + Refresh)
- ✅ Access token: 15 minutes (for sensitive operations)
- ✅ Refresh token: 30 days (stored in HttpOnly cookie)
- ✅ Single-use auth code (60 seconds)
- ✅ Removed JWT exposure from URL (was: `?token=...&user=...`)
- ✅ New endpoints: `/auth/exchange` and `/auth/refresh`

### Phase 3: Input Validation
- ✅ Joi schema validation for all inputs
- ✅ OAuth callback code validation (alphanum, 10-500)
- ✅ Post creation validation (title, content constraints)
- ✅ Auto-cleanup of unknown fields (stripUnknown)

## 📊 Testing Infrastructure (Phase 4)

### Jest Configuration
```typescript
// jest.config.ts
roots: ['<rootDir>/tests', '<rootDir>/src']
coverageThreshold: {
  global: {
    branches: 60,
    functions: 70,
    lines: 70,
    statements: 70
  }
}
```

### Test Scripts
```bash
# Run all tests
npm test

# Run auth-specific tests
npm run test:auth

# Generate coverage report
npm run coverage

# CI coverage upload format
npm run coverage:ci
```

### OAuth Token Tests
- ✅ Token generation with valid JWT structure
- ✅ Token verification success/failure cases
- ✅ Token expiration handling
- ✅ Token pair generation (access + refresh)
- ✅ Refresh token validation
- ✅ Access token rejection when used as refresh token

## 📚 API Documentation (Phase 5)

### Swagger UI
- Endpoint: `/api/docs`
- Spec Format: OpenAPI 3.0.0
- Auto-refresh: Reloads from `docs/api.yaml`

### Documented Endpoints
- POST `/auth/exchange` — Code to token exchange
- POST `/auth/refresh` — Refresh token to new access token
- GET `/api/health` — Health check
- GET `/api/admin/stats` — Dashboard statistics
- POST `/api/admin/posts` — Create post (with Joi validation)

## 🛠️ Environment Variables

```bash
# Security
JWT_SECRET=your-strong-random-secret-change-this
JWT_REFRESH_SECRET=your-strong-random-refresh-secret
SESSION_SECRET=your-session-secret

# Server
NODE_ENV=production
SERVER_PORT=3001
BASE_URL=https://253.dclub.kr
FRONTEND_URL=https://253.dclub.kr

# OAuth Providers
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GITHUB_CLIENT_ID=...
GITHUB_CLIENT_SECRET=...
NAVER_CLIENT_ID=...
NAVER_CLIENT_SECRET=...
```

## 🔍 Validation Rules

### OAuth Code Validation
```
- Type: Alphanumeric only
- Min length: 10 characters
- Max length: 500 characters
- Auto-cleanup: Unknown fields stripped
```

### Post Creation Validation
```
- Title:
  - Min length: 1 character
  - Max length: 200 characters
  - Required: true
- Content:
  - Min length: 1 character
  - Max length: 50,000 characters
  - Required: true
```

## ⚡ Performance & Limits

| Component | Limit | Duration |
|-----------|-------|----------|
| Auth Rate Limit | 10 requests | 15 minutes |
| API Rate Limit | 100 requests | 1 minute |
| Access Token | - | 15 minutes |
| Refresh Token | - | 30 days |
| Auth Code | - | 60 seconds |
| Session | - | 7 days |

## 📋 Deployment Checklist

- [ ] All environment variables set (.env file)
- [ ] JWT_SECRET and JWT_REFRESH_SECRET are strong random values
- [ ] SSL/HTTPS enabled in production
- [ ] CORS origin restricted to frontend domain
- [ ] Rate limiting active
- [ ] Helmet security headers enabled
- [ ] Swagger UI accessible at /api/docs
- [ ] Health check endpoint working
- [ ] Tests passing with >70% coverage
- [ ] OAuth providers configured with correct callback URLs

---

**Last Updated**: 2026-03-14
**Phase**: 1-5 Complete (Security → JWT → Validation → Testing → Documentation)
