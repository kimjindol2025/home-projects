#!/bin/bash

# 🚀 Gogs API 테스트 스크립트

# 설정
API_KEY="dclub-api-key-2025-secure"
GOGS_TOKEN="826b3705d8a0602cf89a02327dcee25e991dd630"
API_URL="http://192.168.45.73:50400"
GOGS_URL="https://gogs.dclub.kr"
REPO="kim/freelang-v4-ide"

# 색상
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Gogs API 테스트 시작${NC}\n"

# 테스트 1: 저장소 목록
echo -e "${YELLOW}📦 [TEST 1] 저장소 목록 조회${NC}"
REPOS=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/repos?limit=5")
if echo "$REPOS" | grep -q '"repos"'; then
  COUNT=$(echo "$REPOS" | grep -o '"repos"' | wc -l)
  echo -e "${GREEN}✅ 성공: 저장소 목록 조회 완료${NC}"
  echo "$REPOS" | head -50
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 2: 저장소 검색
echo -e "${YELLOW}🔍 [TEST 2] 저장소 검색 (q=freelang)${NC}"
SEARCH=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/search?q=freelang&limit=10")
if echo "$SEARCH" | grep -q '"repos"'; then
  echo -e "${GREEN}✅ 성공: 저장소 검색 완료${NC}"
  echo "$SEARCH" | jq '.repos[] | {name, description}' 2>/dev/null || echo "$SEARCH"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 3: 코드 검색
echo -e "${YELLOW}🔎 [TEST 3] 코드 검색 (q=express)${NC}"
CODE=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/freelang-v4-ide/search?q=express&limit=5")
if echo "$CODE" | grep -q '"results"'; then
  echo -e "${GREEN}✅ 성공: 코드 검색 완료${NC}"
  echo "$CODE" | jq '.results[] | {file, line, content}' 2>/dev/null || echo "$CODE"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 4: 커밋 목록
echo -e "${YELLOW}📊 [TEST 4] 커밋 목록 조회${NC}"
COMMITS=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/freelang-v4-ide/commits?limit=10")
if echo "$COMMITS" | grep -q '"commits"'; then
  echo -e "${GREEN}✅ 성공: 커밋 목록 조회 완료${NC}"
  echo "$COMMITS" | jq '.commits[] | {sha_short, message, author}' 2>/dev/null || echo "$COMMITS"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 5: 커밋 메시지 검색
echo -e "${YELLOW}📝 [TEST 5] 커밋 메시지 검색 (q=env)${NC}"
COMMIT_SEARCH=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/freelang-v4-ide/commits/search?q=env")
if echo "$COMMIT_SEARCH" | grep -q '"commits"'; then
  echo -e "${GREEN}✅ 성공: 커밋 메시지 검색 완료${NC}"
  echo "$COMMIT_SEARCH" | jq '.commits[] | {sha_short, message}' 2>/dev/null || echo "$COMMIT_SEARCH"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 6: 저장소 정보
echo -e "${YELLOW}📦 [TEST 6] 저장소 정보 조회${NC}"
REPO_INFO=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/gogs/$REPO")
if echo "$REPO_INFO" | grep -q '"name"'; then
  echo -e "${GREEN}✅ 성공: 저장소 정보 조회 완료${NC}"
  echo "$REPO_INFO" | jq '{name, stars_count, forks_count, open_issues_count}' 2>/dev/null || echo "$REPO_INFO"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 7: 시스템 정보
echo -e "${YELLOW}💻 [TEST 7] 시스템 통계${NC}"
SYSTEM=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/system/stats")
if echo "$SYSTEM" | grep -q '"cpu"'; then
  echo -e "${GREEN}✅ 성공: 시스템 정보 조회 완료${NC}"
  echo "$SYSTEM" | jq '.' 2>/dev/null || echo "$SYSTEM"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 8: CPU 정보
echo -e "${YELLOW}⚙️ [TEST 8] CPU 정보${NC}"
CPU=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/system/cpu")
if echo "$CPU" | grep -q '"usage"'; then
  echo -e "${GREEN}✅ 성공: CPU 정보 조회 완료${NC}"
  echo "$CPU" | jq '{usage: .usage, cores: .cores, frequency: .frequency}' 2>/dev/null || echo "$CPU"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 9: 메모리 정보
echo -e "${YELLOW}🧠 [TEST 9] 메모리 정보${NC}"
MEMORY=$(curl -s -H "x-api-key: $API_KEY" "$API_URL/api/system/memory")
if echo "$MEMORY" | grep -q '"total"'; then
  echo -e "${GREEN}✅ 성공: 메모리 정보 조회 완료${NC}"
  echo "$MEMORY" | jq '{total, used, available, percent}' 2>/dev/null || echo "$MEMORY"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

# 테스트 10: Gogs 원본 API
echo -e "${YELLOW}🔗 [TEST 10] Gogs 원본 API (토큰 인증)${NC}"
GOGS_API=$(curl -s -H "Authorization: token $GOGS_TOKEN" "$GOGS_URL/api/v1/repos/$REPO")
if echo "$GOGS_API" | grep -q '"name"'; then
  echo -e "${GREEN}✅ 성공: Gogs 원본 API 테스트 완료${NC}"
  echo "$GOGS_API" | jq '{name, stars_count, open_issues_count, description}' 2>/dev/null || echo "$GOGS_API"
else
  echo -e "${RED}❌ 실패${NC}"
fi
echo ""

echo -e "${BLUE}🎉 테스트 완료!${NC}"
