# 🚀 Gogs API 가이드 2025

## 📋 기본 설정

### 서버 정보
- **API URL**: `http://192.168.45.73:50400`
- **Gogs URL**: `https://gogs.dclub.kr`
- **API Key**: `dclub-api-key-2025-secure`
- **Gogs Token**: `826b3705d8a0602cf89a02327dcee25e991dd630`

### 사용자 정보
- **Username**: `kim`
- **Repository**: `freelang-v4-ide`
- **Full Repo**: `kim/freelang-v4-ide`

---

## 🔍 Gogs 검색 API (6개)

| 엔드포인트 | 메서드 | 설명 | 예제 |
|-----------|--------|------|------|
| `/api/gogs/repos` | GET | 저장소 목록 (9,645개) | `/api/gogs/repos?page=1&limit=20` |
| `/api/gogs/search` | GET | 저장소 검색 (이름 기반) | `/api/gogs/search?q=freelang` |
| `/api/gogs/:repo/search` | GET | 코드 검색 (git grep) | `/api/gogs/freelang-v4-ide/search?q=express` |
| `/api/gogs/:repo/commits` | GET | 커밋 목록 | `/api/gogs/freelang-v4-ide/commits?limit=10` |
| `/api/gogs/:repo/commits/search` | GET | 커밋 메시지 검색 | `/api/gogs/freelang-v4-ide/commits/search?q=fix` |
| `/api/gogs/:owner/:repo` | GET | 저장소 상세 정보 | `/api/gogs/kim/freelang-v4-ide` |

---

## 🖥️ NexusSSH Pro API (25개)

### Core API (5개)
- **POST** `/api/ssh/connect` - SSH 서버 연결
- **GET** `/api/ssh/status` - 연결 상태 확인
- **POST** `/api/ssh/execute` - 명령 실행
- **POST** `/api/ssh/upload` - 파일 업로드
- **POST** `/api/ssh/download` - 파일 다운로드

### System Resources (8개)
- **GET** `/api/system/cpu` - CPU 정보
- **GET** `/api/system/memory` - 메모리 정보
- **GET** `/api/system/disk` - 디스크 정보
- **GET** `/api/system/network` - 네트워크 정보
- **GET** `/api/system/uptime` - 서버 운영 시간
- **GET** `/api/system/temperature` - 온도 정보
- **GET** `/api/system/load` - 로드 평균
- **GET** `/api/system/stats` - 전체 통계

### PM2 Management (2개)
- **GET** `/api/pm2/processes` - 프로세스 목록
- **POST** `/api/pm2/control` - 프로세스 제어

### Service Management (2개)
- **GET** `/api/services/list` - 서비스 목록
- **POST** `/api/services/control` - 서비스 제어

### Process Management (2개)
- **GET** `/api/process/list` - 프로세스 목록
- **POST** `/api/process/kill` - 프로세스 종료

### Log Management (2개)
- **GET** `/api/logs/list` - 로그 목록
- **GET** `/api/logs/search` - 로그 검색

### File Operations (3개)
- **GET** `/api/files/list` - 파일 목록
- **GET** `/api/files/read` - 파일 읽기
- **POST** `/api/files/write` - 파일 쓰기

**총 31개 API 엔드포인트**

---

## 💻 사용 예시

### Bash / cURL

```bash
#!/bin/bash

API_KEY="dclub-api-key-2025-secure"
GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
API_URL="http://192.168.45.73:50400"
GOGS_URL="https://gogs.dclub.kr"

# 1️⃣ 저장소 검색
echo "=== 저장소 검색 ==="
curl -s -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/search?q=freelang" | jq '.repos[] | .name'

# 2️⃣ 코드 검색
echo "=== 코드 검색 (express 검색) ==="
curl -s -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/freelang-v4-ide/search?q=express&limit=5" | \
  jq -r '.results[] | "\(.file):\(.line) - \(.content)"'

# 3️⃣ 커밋 목록 조회
echo "=== 최근 커밋 ==="
curl -s -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/freelang-v4-ide/commits?limit=5" | \
  jq -r '.commits[] | "\(.sha_short) - \(.message)"'

# 4️⃣ 저장소 상세 정보
echo "=== 저장소 정보 ==="
curl -s -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/kim/freelang-v4-ide" | \
  jq '{name: .name, stars: .stars_count, forks: .forks_count}'

# 5️⃣ 커밋 메시지 검색
echo "=== 특정 커밋 조회 ==="
curl -s -H "x-api-key: $API_KEY" \
  "$API_URL/api/gogs/kim/freelang-v4-ide/commits/search?q=env" | \
  jq '.commits'

# 6️⃣ Gogs 원본 API (토큰 인증)
echo "=== Gogs 원본 API ==="
curl -s -H "Authorization: token $GOGS_TOKEN" \
  "$GOGS_URL/api/v1/repos/kim/freelang-v4-ide" | jq '{name, stars_count}'
```

### JavaScript (Node.js)

```javascript
const axios = require('axios');

const apiClient = axios.create({
  baseURL: 'http://192.168.45.73:50400',
  headers: { 'x-api-key': 'dclub-api-key-2025-secure' }
});

const gogsClient = axios.create({
  baseURL: 'https://gogs.dclub.kr/api/v1',
  headers: { 'Authorization': 'token 826b3705d8a0602cf89a02327dcee25e991dd630' }
});

// 저장소 검색
const searchRepos = async (query) => {
  try {
    const { data } = await apiClient.get('/api/gogs/search', {
      params: { q: query }
    });
    return data.repos;
  } catch (error) {
    console.error('Error searching repos:', error.message);
  }
};

// 코드 검색
const searchCode = async (repo, query, limit = 10) => {
  try {
    const { data } = await apiClient.get(`/api/gogs/${repo}/search`, {
      params: { q: query, limit }
    });
    return data.results;
  } catch (error) {
    console.error('Error searching code:', error.message);
  }
};

// 커밋 목록
const getCommits = async (repo, limit = 20) => {
  try {
    const { data } = await apiClient.get(`/api/gogs/${repo}/commits`, {
      params: { limit }
    });
    return data.commits;
  } catch (error) {
    console.error('Error fetching commits:', error.message);
  }
};

// 저장소 정보
const getRepoInfo = async (owner, repo) => {
  try {
    const { data } = await gogsClient.get(`/repos/${owner}/${repo}`);
    return data;
  } catch (error) {
    console.error('Error fetching repo info:', error.message);
  }
};

// 사용 예
(async () => {
  console.log('🔍 저장소 검색...');
  const repos = await searchRepos('freelang');
  console.log(`✅ 찾은 저장소: ${repos.length}개`);

  console.log('\n📝 코드 검색...');
  const code = await searchCode('freelang-v4-ide', 'express', 5);
  console.log(`✅ 검색 결과: ${code.length}개`);

  console.log('\n📊 커밋 조회...');
  const commits = await getCommits('freelang-v4-ide', 10);
  console.log(`✅ 커밋: ${commits.length}개`);

  console.log('\n📦 저장소 정보...');
  const repo = await getRepoInfo('kim', 'freelang-v4-ide');
  console.log(`✅ 저장소: ${repo.name} (⭐ ${repo.stars_count})`);
})();
```

### Python

```python
import requests
import json

class GogsAPI:
    def __init__(self, api_url='http://192.168.45.73:50400', api_key='dclub-api-key-2025-secure'):
        self.api_url = api_url
        self.headers = {'x-api-key': api_key}

    def search_repos(self, query):
        """저장소 검색"""
        response = requests.get(
            f'{self.api_url}/api/gogs/search',
            headers=self.headers,
            params={'q': query}
        )
        return response.json().get('repos', [])

    def search_code(self, repo, query, limit=20):
        """코드 검색"""
        response = requests.get(
            f'{self.api_url}/api/gogs/{repo}/search',
            headers=self.headers,
            params={'q': query, 'limit': limit}
        )
        return response.json().get('results', [])

    def get_commits(self, repo, limit=20):
        """커밋 목록"""
        response = requests.get(
            f'{self.api_url}/api/gogs/{repo}/commits',
            headers=self.headers,
            params={'limit': limit}
        )
        return response.json().get('commits', [])

    def search_commits(self, repo, query):
        """커밋 메시지 검색"""
        response = requests.get(
            f'{self.api_url}/api/gogs/{repo}/commits/search',
            headers=self.headers,
            params={'q': query}
        )
        return response.json().get('commits', [])

class NexusAPI:
    def __init__(self, api_url='http://192.168.45.73:50400', api_key='dclub-api-key-2025-secure'):
        self.api_url = api_url
        self.headers = {'x-api-key': api_key}

    def get_system_stats(self):
        """시스템 통계"""
        response = requests.get(
            f'{self.api_url}/api/system/stats',
            headers=self.headers
        )
        return response.json()

    def get_cpu_info(self):
        """CPU 정보"""
        response = requests.get(
            f'{self.api_url}/api/system/cpu',
            headers=self.headers
        )
        return response.json()

    def execute_command(self, command):
        """명령 실행"""
        response = requests.post(
            f'{self.api_url}/api/ssh/execute',
            headers=self.headers,
            json={'command': command}
        )
        return response.json()

# 사용 예
if __name__ == '__main__':
    gogs = GogsAPI()
    nexus = NexusAPI()

    # Gogs API 사용
    print("🔍 Gogs API 테스트")
    repos = gogs.search_repos('freelang')
    print(f"✅ 찾은 저장소: {len(repos)}개")

    code = gogs.search_code('freelang-v4-ide', 'express', 5)
    print(f"✅ 코드 검색: {len(code)}개")

    commits = gogs.get_commits('freelang-v4-ide', 10)
    print(f"✅ 커밋: {len(commits)}개")

    # Nexus API 사용
    print("\n💻 Nexus API 테스트")
    stats = nexus.get_system_stats()
    print(f"✅ 시스템 통계: {json.dumps(stats, indent=2)}")
```

---

## 🔐 보안 주의사항

1. **API Key 보안**
   - API Key를 코드에 노출시키지 마세요
   - 환경 변수로 관리하세요: `.env` 파일 사용

2. **Token 관리**
   - Gogs Token을 안전하게 보관하세요
   - 정기적으로 토큰을 갱신하세요

3. **HTTPS 사용**
   - 프로덕션 환경에서는 반드시 HTTPS를 사용하세요

4. **Rate Limiting**
   - API 호출에 제한이 있을 수 있습니다
   - 적절한 딜레이를 추가하세요

---

## 📦 .env 파일 설정

```env
# Gogs API 설정
GOGS_URL=https://gogs.dclub.kr
GOGS_API_TOKEN=826b3705d8a0602cf89a02327dcee25e991dd630
GOGS_REPO_OWNER=kim
GOGS_REPO_NAME=freelang-v4-ide

# NexusSSH Pro API 설정
NEXUS_API_URL=http://192.168.45.73:50400
NEXUS_API_KEY=dclub-api-key-2025-secure

# 기타 설정
NODE_ENV=production
LOG_LEVEL=info
```

---

## 🎯 빠른 시작

### 1. 저장소 목록 조회
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  http://192.168.45.73:50400/api/gogs/repos
```

### 2. 특정 저장소 검색
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  "http://192.168.45.73:50400/api/gogs/search?q=freelang"
```

### 3. 코드 검색
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  "http://192.168.45.73:50400/api/gogs/freelang-v4-ide/search?q=env"
```

### 4. 커밋 조회
```bash
curl -H "x-api-key: dclub-api-key-2025-secure" \
  "http://192.168.45.73:50400/api/gogs/freelang-v4-ide/commits"
```

---

## 📞 지원 정보

- **Gogs 서버**: https://gogs.dclub.kr
- **API 서버**: http://192.168.45.73:50400
- **문서**: 이 가이드 참고

---

**마지막 업데이트**: 2025년 2월 20일
