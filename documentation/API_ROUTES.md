# BookCommunity API 路由文档

## 📋 API 总览

BookCommunity 采用 **RESTful 风格**的 API 设计，所有 API 都在 `/api` 路径下。

**Base URL**: `http://localhost:8080/api`

---

## 🔐 认证说明

### JWT Token 认证

大部分需要用户身份的 API 需要在 Header 中携带 JWT Token：

```http
Authorization: Bearer <your_jwt_token>
```

### 获取 Token

通过注册或登录接口获取：

```bash
# 注册
POST /api/register
{
  "username": "testuser",
  "password": "password123"
}

# 登录
POST /api/login
{
  "username": "testuser",
  "password": "password123"
}

# 响应
{
  "status_code": 0,
  "status_msg": "登录成功",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 📚 API 分类

### 1. 用户相关 `/api/users`

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/register` | ❌ | 用户注册 |
| POST | `/api/login` | ❌ | 用户登录 |
| GET | `/api/users/:id` | ✅ | 获取用户信息 |
| GET | `/api/users/:user_id/collections` | ❌ | 获取用户收藏列表 |
| POST | `/api/users/:id/follow` | ✅ | 关注用户 |
| GET | `/api/users/:id/followers` | ❌ | 获取粉丝列表 |
| GET | `/api/users/:id/following` | ❌ | 获取关注列表 |

---

### 2. 书评相关 `/api/reviews` ⭐ 核心功能

#### 2.1 书评 CRUD

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/reviews` | ✅ | 创建书评 |
| GET | `/api/reviews` | ❌ | 查询书评列表（支持分页、筛选、排序） |
| GET | `/api/reviews/:id` | ❌ | 查询书评详情 |
| PUT | `/api/reviews/:id` | ✅ | 更新书评（只能更新自己的） |
| DELETE | `/api/reviews/:id` | ✅ | 删除书评（只能删除自己的） |

#### 示例：创建书评

```bash
POST /api/reviews
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "《深入理解计算机系统》读后感",
  "content": "这本书深入浅出地讲解了计算机系统的本质...",
  "book_isbn": "9787111544937",
  "book_title": "深入理解计算机系统",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "rating": 9.5,
  "tags": ["计算机", "经典", "必读"]
}
```

#### 示例：查询书评列表

```bash
GET /api/reviews?page=1&page_size=20&order_by=popular&user_id=1

# 查询参数:
# - page: 页码（默认1）
# - page_size: 每页数量（默认20，最大100）
# - user_id: 筛选指定用户的书评
# - book_isbn: 筛选指定图书的书评
# - order_by: 排序方式（latest=最新, popular=最热, rating=评分）
```

#### 2.2 点赞相关

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/reviews/:id/like` | ✅ | 点赞书评 |
| DELETE | `/api/reviews/:id/like` | ✅ | 取消点赞 |
| GET | `/api/reviews/:id/likes` | ❌ | 获取点赞列表 |

#### 2.3 评论相关

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/reviews/:id/comments` | ✅ | 发布评论 |
| GET | `/api/reviews/:id/comments` | ❌ | 获取评论列表 |

#### 示例：发布评论

```bash
POST /api/reviews/123/comments
Content-Type: application/json
Authorization: Bearer <token>

{
  "content": "写得很好，受益匪浅！"
}
```

#### 2.4 收藏相关

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/reviews/:id/collect` | ✅ | 收藏书评 |
| DELETE | `/api/reviews/:id/collect` | ✅ | 取消收藏 |

---

### 3. 评论（独立资源） `/api/comments`

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| DELETE | `/api/comments/:id` | ✅ | 删除评论（只能删除自己的） |

---

### 4. Feed 流 `/api/feed`

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/feed` | ❌ | 发现页（热度推荐） |
| GET | `/api/feed/following` | ✅ | 关注页（关注用户的书评流） |

#### 示例：获取发现页

```bash
GET /api/feed?latest_time=0&page_size=20

# 查询参数:
# - latest_time: 最新一条的时间戳（用于下拉刷新，传0获取最新）
# - page_size: 每页数量（默认20）

# 响应
{
  "status_code": 0,
  "status_msg": "查询成功",
  "reviews": [
    {
      "id": 1,
      "title": "书评标题",
      "content": "书评内容...",
      "images": ["url1", "url2"],
      "rating": 9.5,
      "author": {
        "id": 1,
        "username": "张三"
      },
      "like_count": 123,
      "comment_count": 45,
      ...
    }
  ],
  "next_time": 1634567890,
  "has_more": true
}
```

---

### 5. 图书推荐 `/api/books` 🤖 AI 服务

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/books/search` | ❌ | 搜索图书（RAG 语义搜索） |
| GET | `/api/books/recommendations` | ✅ | 个性化推荐 |
| GET | `/api/books/:isbn` | ❌ | 获取图书详情 |

#### 示例：搜索图书

```bash
GET /api/books/search?q=深入理解计算机系统&top_k=10

# 查询参数:
# - q: 搜索关键词
# - top_k: 返回结果数量（默认10）
```

#### 示例：个性化推荐

```bash
GET /api/books/recommendations?top_k=20
Authorization: Bearer <token>

# 查询参数:
# - top_k: 返回结果数量（默认10）
```

---

## 🔄 响应格式

### 通用响应结构

所有 API 都返回统一的响应结构：

```json
{
  "status_code": 0,        // 0=成功, 非0=失败
  "status_msg": "操作成功",
  "data": {                // 具体数据（根据API不同）
    ...
  }
}
```

### 错误码

| 状态码 | 说明 |
|--------|------|
| 0 | 成功 |
| 500 | 服务器错误 |
| 1001 | 用户未登录 |
| 1002 | Token 无效 |
| 1003 | 用户名或密码错误 |

---

## 📊 完整 API 列表

### 用户相关（7个）

```
POST   /api/register                    - 注册
POST   /api/login                       - 登录
GET    /api/users/:id                   - 用户信息
GET    /api/users/:user_id/collections  - 收藏列表
POST   /api/users/:id/follow            - 关注
GET    /api/users/:id/followers         - 粉丝列表
GET    /api/users/:id/following         - 关注列表
```

### 书评相关（11个）

```
POST   /api/reviews                     - 创建书评
GET    /api/reviews                     - 查询列表
GET    /api/reviews/:id                 - 查询详情
PUT    /api/reviews/:id                 - 更新书评
DELETE /api/reviews/:id                 - 删除书评
POST   /api/reviews/:id/like            - 点赞
DELETE /api/reviews/:id/like            - 取消点赞
GET    /api/reviews/:id/likes           - 点赞列表
POST   /api/reviews/:id/comments        - 发布评论
GET    /api/reviews/:id/comments        - 评论列表
POST   /api/reviews/:id/collect         - 收藏书评
DELETE /api/reviews/:id/collect         - 取消收藏
```

### 评论相关（1个）

```
DELETE /api/comments/:id                - 删除评论
```

### Feed 流（2个）

```
GET    /api/feed                        - 发现页
GET    /api/feed/following              - 关注页
```

### 图书推荐（3个）

```
GET    /api/books/search                - 搜索图书
GET    /api/books/recommendations       - 个性化推荐
GET    /api/books/:isbn                 - 图书详情
```

**总计: 24 个 API**

---

## 🔧 开发工具

### Swagger UI

访问 Swagger UI 查看完整的 API 文档：

```
http://localhost:8080/swagger/index.html
```

### Postman Collection

导入 Postman 集合：

```
examples/BookCommunity.postman_collection.json
```

### 测试脚本

运行自动化测试脚本：

```bash
bash examples/test_api.sh
```

---

## 📝 RESTful 设计原则

### 1. 资源清晰

- `/reviews` - 书评资源
- `/users` - 用户资源
- `/comments` - 评论资源
- `/feed` - Feed 流资源

### 2. HTTP 方法语义

- `GET` - 查询资源
- `POST` - 创建资源
- `PUT` - 更新资源
- `DELETE` - 删除资源

### 3. URL 层级

- `/api/reviews/:id` - 单个书评
- `/api/reviews/:id/comments` - 书评的评论
- `/api/reviews/:id/likes` - 书评的点赞

### 4. 幂等性

- `GET`, `PUT`, `DELETE` 操作是幂等的
- 重复点赞、收藏会返回"已经点赞/收藏"

---

## 🚀 性能优化

### 分页策略

- **列表分页**: 使用 `page` 和 `page_size`
- **Feed 流分页**: 使用 `latest_time`（基于时间戳）

### 缓存策略

- 书评详情: Redis 缓存 1 小时
- 热门书评: Redis 缓存 10 分钟
- 用户信息: Redis 缓存 30 分钟

### 预加载

- 查询书评时自动预加载作者信息
- 查询评论时自动预加载评论者信息

---

## 📚 下一步

- [ ] 添加图片上传 API (`POST /api/upload/images`)
- [ ] 添加 Chat with Book API (`POST /api/books/:isbn/chat`)
- [ ] 添加标签系统 API
- [ ] 添加书单功能 API

---

**更新时间**: 2026-02-13
**API 版本**: v1.0
**文档状态**: ✅ 已完成
