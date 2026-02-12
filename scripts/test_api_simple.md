# BookCommunity API 简单测试

## 前置条件

确保服务已启动：
```bash
go run main.go
# 或
./bookcommunity
```

---

## 1. 健康检查

```bash
curl http://localhost:8080/health
```

**预期响应：**
```json
{
  "service": "BookCommunity API",
  "status": "healthy"
}
```

---

## 2. 用户注册

```bash
curl -X POST http://localhost:8080/douyin/user/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "Password123!",
    "email": "alice@example.com"
  }'
```

**预期响应：**
```json
{
  "status_code": 0,
  "token": "encrypted_jwt_token...",
  "user_id": 1
}
```

**保存token供后续使用**

---

## 3. 用户登录

```bash
curl -X POST http://localhost:8080/douyin/user/login/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "Password123!"
  }'
```

---

## 4. 获取个性化推荐 (需要登录)

```bash
# 替换 YOUR_TOKEN 为实际token
curl "http://localhost:8080/douyin/recommend?token=YOUR_TOKEN&top_k=5"
```

**预期响应：**
```json
{
  "status_code": 0,
  "books": [
    {
      "isbn": "9787111544937",
      "title": "深入理解计算机系统（原书第3版）",
      "author": "Randal E. Bryant / David R. O'Hallaron",
      "cover_url": "https://img3.doubanio.com/view/subject/l/public/s29195878.jpg",
      "rating": 9.7,
      "reason": "基于你的阅读历史推荐",
      "publisher": "机械工业出版社",
      "pub_date": "2016-11",
      "summary": "从程序员的视角，看计算机系统！..."
    },
    ...
  ],
  "total": 5,
  "message": "当前为模拟推荐数据，可对接真实推荐系统"
}
```

---

## 5. 搜索图书 (无需登录)

```bash
curl "http://localhost:8080/douyin/search?q=计算机&top_k=3"
```

**预期响应：**
```json
{
  "status_code": 0,
  "books": [
    {
      "isbn": "9787111544937",
      "title": "深入理解计算机系统（原书第3版）",
      "author": "Randal E. Bryant / David R. O'Hallaron",
      ...
    },
    ...
  ],
  "total": 3,
  "message": "当前为模拟搜索结果，可对接RAG检索系统"
}
```

---

## 6. 获取图书详情 (无需登录)

```bash
curl http://localhost:8080/douyin/book/9787111544937
```

**预期响应：**
```json
{
  "status_code": 0,
  "message": "图书详情功能待实现",
  "book": {
    "isbn": "9787111544937",
    "title": "示例图书"
  }
}
```

---

## 测试所有Mock推荐书籍

当前mock数据包含以下经典计算机书籍：

1. **深入理解计算机系统（原书第3版）** - Bryant (评分9.7)
2. **Go语言圣经** - Donovan (评分9.5)
3. **编码：隐匿在计算机软硬件背后的语言** - Petzold (评分9.3)
4. **代码大全（第2版）** - McConnell (评分9.3)
5. **算法（第4版）** - Sedgewick (评分9.4)
6. **计算机程序的构造和解释** - Abelson (评分9.5)
7. **设计模式** - Gang of Four (评分9.1)
8. **数据结构与算法分析：C语言描述** - Weiss (评分9.0)
9. **Python编程：从入门到实践** - Matthes (评分9.1)
10. **Effective Java中文版** - Bloch (评分9.1)

---

## 注意事项

1. **当前使用Mock数据**：推荐和搜索功能返回预设的图书列表
2. **Token格式**：JWT经过AES加密，较长且包含特殊字符，使用时需要URL编码
3. **未来扩展**：可通过修改配置文件 `recommendation.enabled: true` 启用真实推荐系统
4. **端口配置**：默认8080，可在 `config.yaml` 中修改

---

## 使用Postman测试

导入提供的 Postman Collection:
```bash
scripts/BookCommunity.postman_collection.json
```

或者使用自动化测试脚本：
```bash
./scripts/test_recommendation_api.sh
```
