# 📚 Gogs API 설정 완료

## 🎯 생성된 파일 목록

| 파일 | 설명 | 크기 |
|------|------|------|
| **GOGS_API_GUIDE.md** | 완전한 API 가이드 (31개 엔드포인트) | 11KB |
| **GOGS_CHEATSHEET.md** | 빠른 참고 치트시트 | 3.0KB |
| **test-gogs-api.sh** | API 테스트 스크립트 (10개 테스트) | 4.7KB |
| **README_GOGS.md** | 이 파일 | - |

---

## 🚀 빠른 시작

### 1️⃣ 테스트 스크립트 실행
```bash
./test-gogs-api.sh
```

### 2️⃣ 직접 API 호출
```bash
# 환경 변수 설정
export API_KEY="dclub-api-key-2025-secure"
export API_URL="http://192.168.45.73:50400"

# 저장소 검색
curl -H "x-api-key: $API_KEY" "$API_URL/api/gogs/search?q=freelang"
```

### 3️⃣ Python으로 사용
```python
import requests

headers = {'x-api-key': 'dclub-api-key-2025-secure'}
response = requests.get(
    'http://192.168.45.73:50400/api/gogs/search',
    headers=headers,
    params={'q': 'freelang'}
)
print(response.json())
```

---

## 📋 핵심 정보

### 서버 정보
```
API Server:     http://192.168.45.73:50400
Gogs Server:    https://gogs.dclub.kr
SSH Gateway:    gogs.dclub.kr (포트 22)
```

### 인증 정보
```
API Key:        dclub-api-key-2025-secure
Gogs Token:     826b3705d8a0602cf89a02327dcee25e991dd630
Username:       kim
Repository:     freelang-v4-ide
Full Path:      kim/freelang-v4-ide
```

### 저장소 정보
```
Git Clone:      git clone https://kim:TOKEN@gogs.dclub.kr/kim/freelang-v4-ide.git
SSH Clone:      ssh://kim@gogs.dclub.kr:22/kim/freelang-v4-ide.git
HTTP URL:       https://gogs.dclub.kr/kim/freelang-v4-ide
```

---

## 📊 API 엔드포인트 분류

### Gogs 검색 API (6개)
- `GET /api/gogs/repos` - 저장소 목록
- `GET /api/gogs/search` - 저장소 검색
- `GET /api/gogs/:repo/search` - 코드 검색
- `GET /api/gogs/:repo/commits` - 커밋 목록
- `GET /api/gogs/:repo/commits/search` - 커밋 메시지 검색
- `GET /api/gogs/:owner/:repo` - 저장소 정보

### 시스템 API (8개)
- `GET /api/system/cpu` - CPU 정보
- `GET /api/system/memory` - 메모리 정보
- `GET /api/system/disk` - 디스크 정보
- `GET /api/system/network` - 네트워크 정보
- `GET /api/system/uptime` - 운영 시간
- `GET /api/system/temperature` - 온도
- `GET /api/system/load` - 로드 평균
- `GET /api/system/stats` - 전체 통계

### SSH/파일 API (10개)
- `POST /api/ssh/execute` - 명령 실행
- `POST /api/ssh/upload` - 파일 업로드
- `POST /api/ssh/download` - 파일 다운로드
- `GET /api/pm2/processes` - PM2 프로세스
- `GET /api/services/list` - 서비스 목록
- `GET /api/process/list` - 프로세스 목록
- `POST /api/process/kill` - 프로세스 종료
- `GET /api/logs/list` - 로그 목록
- `GET /api/files/list` - 파일 목록
- `GET /api/files/read` - 파일 읽기

**총 31개 API 엔드포인트**

---

## 🔧 환경 변수 설정

### .env 파일 위치
```
~/.env (또는)
~/test-freelang-v4/.env
```

### 필수 환경 변수
```env
# Gogs 설정
GOGS_URL=https://gogs.dclub.kr
GOGS_API_TOKEN=826b3705d8a0602cf89a02327dcee25e991dd630
GOGS_REPO_OWNER=kim
GOGS_REPO_NAME=freelang-v4-ide

# API 설정
NEXUS_API_URL=http://192.168.45.73:50400
NEXUS_API_KEY=dclub-api-key-2025-secure

# 서버 설정
PORT=50101
NODE_ENV=development
```

---

## 💡 사용 예시

### cURL로 저장소 검색
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  "http://192.168.45.73:50400/api/gogs/search?q=freelang"
```

### cURL로 코드 검색
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  "http://192.168.45.73:50400/api/gogs/freelang-v4-ide/search?q=express"
```

### Node.js로 저장소 정보 조회
```javascript
const axios = require('axios');

const api = axios.create({
  baseURL: 'http://192.168.45.73:50400',
  headers: { 'x-api-key': 'dclub-api-key-2025-secure' }
});

const repo = await api.get('/api/gogs/kim/freelang-v4-ide');
console.log(repo.data);
```

### Python으로 시스템 정보 조회
```python
import requests

response = requests.get(
    'http://192.168.45.73:50400/api/system/stats',
    headers={'x-api-key': 'dclub-api-key-2025-secure'}
)
print(response.json())
```

---

## 🔐 보안 체크리스트

- [ ] API Key를 환경 변수로 관리
- [ ] Gogs Token을 .gitignore에 추가
- [ ] 프로덕션에서 HTTPS 사용
- [ ] API Rate Limiting 확인
- [ ] 토큰 만료 정책 설정
- [ ] 로깅 및 모니터링 활성화
- [ ] 정기적인 토큰 갱신

---

## 📖 더 알아보기

전체 상세 정보는 다음 파일을 참고하세요:

1. **GOGS_API_GUIDE.md** - 완전한 기술 문서
   - 모든 31개 API 엔드포인트 설명
   - Bash, JavaScript, Python 예제

2. **GOGS_CHEATSHEET.md** - 빠른 참고
   - 자주 사용하는 API
   - 한줄 코드 샘플

3. **test-gogs-api.sh** - 테스트 스크립트
   - 10개 자동 테스트
   - 실시간 API 검증

---

## 🎯 다음 단계

### 프로젝트 클론
```bash
cd ~/test-freelang-v4
npm install
npm run build
npm run dev
```

### API 테스트
```bash
./test-gogs-api.sh
```

### 환경 설정
```bash
# .env 파일 확인
cat ~/test-freelang-v4/.env
```

---

## 📞 트러블슈팅

### API 연결 실패
```bash
# 서버 상태 확인
curl -I http://192.168.45.73:50400

# DNS 확인
nslookup gogs.dclub.kr

# 네트워크 확인
ping 192.168.45.73
```

### 토큰 오류
- 토큰이 만료되었을 수 있습니다
- Gogs 대시보드에서 새 토큰 발급
- 환경 변수 업데이트

### 레이트 제한
- API 호출에 지연 추가
- 배치 요청 사용
- 서버 관리자에 문의

---

## ✅ 체크리스트

생성된 리소스:
- [x] .env 환경 변수 파일
- [x] config.ts (Node.js 설정)
- [x] GOGS_API_GUIDE.md (완전 가이드)
- [x] GOGS_CHEATSHEET.md (치트시트)
- [x] test-gogs-api.sh (테스트 스크립트)
- [x] README_GOGS.md (이 문서)

저장 위치: `~/` (홈 디렉토리)

---

## 📅 업데이트 정보

- **생성일**: 2025년 2월 20일
- **버전**: 1.0.0
- **API 엔드포인트**: 31개
- **문서 파일**: 4개

---

**모든 설정이 완료되었습니다! 🎉**

더 많은 정보는 `GOGS_API_GUIDE.md`를 참고하세요.
