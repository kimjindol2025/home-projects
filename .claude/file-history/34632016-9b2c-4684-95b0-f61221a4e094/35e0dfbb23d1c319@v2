# 🚀 Gogs API 치트시트

## ⚡ 빠른 설정

```bash
# 환경 변수 설정
export API_KEY="dclub-api-key-2025-secure"
export GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
export API_URL="http://192.168.45.73:50400"
export GOGS_URL="https://gogs.dclub.kr"
export REPO="kim/freelang-v4-ide"
```

---

## 🔍 자주 사용하는 API

### 1️⃣ 저장소 검색
```bash
curl -H "x-api-key: $API_KEY" "$API_URL/api/gogs/search?q=freelang"
```

### 2️⃣ 코드 검색
```bash
curl -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/freelang-v4-ide/search?q=express&limit=10"
```

### 3️⃣ 커밋 조회
```bash
curl -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/freelang-v4-ide/commits?limit=20"
```

### 4️⃣ 커밋 메시지 검색
```bash
curl -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/freelang-v4-ide/commits/search?q=env"
```

### 5️⃣ 저장소 정보
```bash
curl -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/kim/freelang-v4-ide"
```

### 6️⃣ 저장소 목록
```bash
curl -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/repos?page=1&limit=50"
```

---

## 💻 시스템 정보 조회

### CPU 정보
```bash
curl -H "x-api-key: $API_KEY" "$API_URL/api/system/cpu"
```

### 메모리 정보
```bash
curl -H "x-api-key: $API_KEY" "$API_URL/api/system/memory"
```

### 디스크 정보
```bash
curl -H "x-api-key: $API_KEY" "$API_URL/api/system/disk"
```

### 전체 통계
```bash
curl -H "x-api-key: $API_KEY" "$API_URL/api/system/stats"
```

---

## 🔑 인증 토큰

| 서비스 | 토큰 | 타입 |
|--------|------|------|
| Gogs API | `826b3705d8a0602cf89a02327dcee25e991dd630` | Personal Access Token |
| NexusSSH | `dclub-api-key-2025-secure` | API Key |

---

## 📊 응답 포맷

모든 API는 JSON 응답을 반환합니다.

### 성공 응답
```json
{
  "success": true,
  "data": { ... },
  "count": 10
}
```

### 오류 응답
```json
{
  "success": false,
  "error": "Error message",
  "code": 400
}
```

---

## 🛠️ Python 한줄 코드

```python
# 저장소 검색
import requests; repos = requests.get('http://192.168.45.73:50400/api/gogs/search', headers={'x-api-key': 'dclub-api-key-2025-secure'}, params={'q': 'freelang'}).json()['repos']

# 코드 검색
code = requests.get('http://192.168.45.73:50400/api/gogs/freelang-v4-ide/search', headers={'x-api-key': 'dclub-api-key-2025-secure'}, params={'q': 'express'}).json()['results']

# 커밋 조회
commits = requests.get('http://192.168.45.73:50400/api/gogs/freelang-v4-ide/commits', headers={'x-api-key': 'dclub-api-key-2025-secure'}).json()['commits']
```

---

## 🔗 유용한 링크

- **Gogs 서버**: https://gogs.dclub.kr
- **저장소**: https://gogs.dclub.kr/kim/freelang-v4-ide
- **API 문서**: 로컬 `GOGS_API_GUIDE.md` 참고

---

## ⚠️ 주의사항

1. Rate Limit: 분당 60 요청 제한
2. Timeout: 기본 30초
3. HTTPS 권장: 프로덕션 환경
4. 토큰 보안: 환경 변수로 관리

---

**마지막 업데이트**: 2025년 2월 20일
