#!/bin/bash

# BookCommunity API 测试脚本
# 使用方法: chmod +x test_api.sh && ./test_api.sh

set -e

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/douyin"

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    local headers=$4

    echo -e "${BLUE}=== Testing: $name ===${NC}"
    echo "URL: $method $url"

    if [ -z "$headers" ]; then
        response=$(curl -s -X $method "$url")
    else
        response=$(curl -s -X $method -H "$headers" "$url")
    fi

    echo "Response:"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
    echo ""
}

echo -e "${GREEN}BookCommunity API 测试${NC}"
echo "================================"
echo ""

# 1. 健康检查
test_endpoint "Health Check" "$BASE_URL/health"

# 2. 用户注册
echo -e "${BLUE}=== 用户注册 ===${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/user/register/?username=testuser_$(date +%s)&password=password123")
echo "$REGISTER_RESPONSE" | python3 -m json.tool
TOKEN=$(echo "$REGISTER_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))")
USER_ID=$(echo "$REGISTER_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('user_id', ''))")
echo -e "${GREEN}Token: $TOKEN${NC}"
echo -e "${GREEN}User ID: $USER_ID${NC}"
echo ""

# 3. 获取用户信息
if [ -n "$TOKEN" ]; then
    test_endpoint "Get User Info" "$API_URL/user/?user_id=$USER_ID" "GET" "Authorization: Bearer $TOKEN"
fi

# 4. 搜索图书（不需要认证）
test_endpoint "Search Books" "$API_URL/search?q=golang&top_k=3"

# 5. 获取推荐（需要认证）
if [ -n "$TOKEN" ]; then
    test_endpoint "Get Recommendations" "$API_URL/recommend?top_k=5" "GET" "Authorization: Bearer $TOKEN"
fi

# 6. 获取图书详情
test_endpoint "Get Book Detail" "$API_URL/book/9787111544937"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo ""
echo "查看 Swagger 文档: $BASE_URL/swagger/index.html"
