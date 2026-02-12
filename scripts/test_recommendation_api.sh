#!/bin/bash

# BookCommunity 推荐API测试脚本

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "  BookCommunity 推荐API测试"
echo "========================================"
echo ""

# 1. 健康检查
echo -e "${YELLOW}[1/5] 健康检查${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ 健康检查通过${NC}"
    echo "$body" | jq '.'
else
    echo -e "${RED}✗ 健康检查失败 (HTTP $http_code)${NC}"
    echo "$body"
    exit 1
fi
echo ""

# 2. 用户注册
echo -e "${YELLOW}[2/5] 注册测试用户${NC}"
username="testuser_$(date +%s)"
password="TestPassword123!"
email="${username}@example.com"

response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/douyin/user/register/" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"${username}\",
    \"password\": \"${password}\",
    \"email\": \"${email}\"
  }")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    token=$(echo "$body" | jq -r '.token')
    user_id=$(echo "$body" | jq -r '.user_id')
    echo -e "${GREEN}✓ 注册成功${NC}"
    echo "  用户名: $username"
    echo "  用户ID: $user_id"
    echo "  Token: ${token:0:30}..."
else
    echo -e "${RED}✗ 注册失败 (HTTP $http_code)${NC}"
    echo "$body" | jq '.'
    exit 1
fi
echo ""

# 3. 获取个性化推荐
echo -e "${YELLOW}[3/5] 获取个性化推荐 (mock数据)${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/douyin/recommend?token=${token}&top_k=5")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ 推荐获取成功${NC}"
    book_count=$(echo "$body" | jq '.books | length')
    echo "  返回图书数量: $book_count"
    echo ""
    echo "  推荐书籍列表:"
    echo "$body" | jq '.books[] | "  - \(.title) (\(.author)) - 评分: \(.rating)"'
    echo ""
    echo "  提示信息: $(echo "$body" | jq -r '.message')"
else
    echo -e "${RED}✗ 推荐获取失败 (HTTP $http_code)${NC}"
    echo "$body" | jq '.'
fi
echo ""

# 4. 搜索图书
echo -e "${YELLOW}[4/5] 搜索图书 (关键词: 计算机)${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/douyin/search?q=%E8%AE%A1%E7%AE%97%E6%9C%BA&top_k=3")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ 搜索成功${NC}"
    book_count=$(echo "$body" | jq '.books | length')
    echo "  返回图书数量: $book_count"
    echo ""
    echo "  搜索结果:"
    echo "$body" | jq '.books[] | "  - \(.title) (\(.isbn))"'
    echo ""
    echo "  提示信息: $(echo "$body" | jq -r '.message')"
else
    echo -e "${RED}✗ 搜索失败 (HTTP $http_code)${NC}"
    echo "$body" | jq '.'
fi
echo ""

# 5. 获取图书详情
echo -e "${YELLOW}[5/5] 获取图书详情 (ISBN: 9787111544937)${NC}"
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/douyin/book/9787111544937")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ 获取成功${NC}"
    echo "$body" | jq '.'
else
    echo -e "${RED}✗ 获取失败 (HTTP $http_code)${NC}"
    echo "$body" | jq '.'
fi
echo ""

echo "========================================"
echo -e "${GREEN}测试完成！${NC}"
echo "========================================"
echo ""
echo "说明："
echo "  - 推荐功能当前使用mock数据"
echo "  - 搜索功能当前使用mock数据"
echo "  - 未来可对接真实推荐系统"
echo ""
