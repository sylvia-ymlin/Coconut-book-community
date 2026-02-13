# BookCommunity 图书社交平台改造方案

## 📋 改造目标

将 BookCommunity 从"抖音视频社交"改造为"图书社交平台"，保留完整的社交功能，同时为集成 Python 推荐服务做准备。

## 🎯 核心改动

### 1. 数据模型重命名与扩展

#### VideoModel → BookReviewModel
**改造前:**
```go
type VideoModel struct {
    gorm.Model
    Title        string
    StorageID    uint
    URL          string     // 视频URL
    CoverURL     string
    Author       UserModel
    AuthorID     uint
    LikeCount    uint
    CommentCount uint
    Comments     []CommentModel
    Likes        []UserModel
    Collections  []UserModel
}
```

**改造后:**
```go
type BookReviewModel struct {
    gorm.Model
    Title        string     // 书评标题
    Content      string     // 书评内容（可选，纯文字书评）
    BookISBN     string     // 关联图书 ISBN
    BookTitle    string     // 冗余存储，提升查询性能
    VideoURL     string     // 视频书评URL（可选）
    CoverURL     string     // 封面/缩略图
    ReviewType   string     // "video" | "text" | "mixed"
    Rating       float64    // 用户评分 (0-10)
    Author       UserModel
    AuthorID     uint
    LikeCount    uint
    CommentCount uint
    ViewCount    uint       // 新增：浏览次数
    Comments     []CommentModel
    Likes        []UserModel
    Collections  []UserModel
}
```

### 2. API 路由重构

#### 改造前路由 (Douyin 风格)
```
POST   /douyin/publish/action/       - 发布视频
GET    /douyin/publish/list/         - 查询发布列表
GET    /douyin/feed                  - 视频流
POST   /douyin/favorite/action/      - 点赞
GET    /douyin/favorite/list/        - 点赞列表
POST   /douyin/comment/action/       - 评论
GET    /douyin/comment/list/         - 评论列表
POST   /douyin/relation/action/      - 关注
GET    /douyin/relation/follow/list/ - 关注列表
GET    /douyin/relation/follower/list/ - 粉丝列表
```

#### 改造后路由 (BookCommunity 风格)
```
POST   /api/reviews                  - 发布书评
GET    /api/reviews                  - 获取书评列表（支持过滤）
GET    /api/reviews/:id              - 获取书评详情
PUT    /api/reviews/:id              - 更新书评
DELETE /api/reviews/:id              - 删除书评

GET    /api/books/:isbn/reviews      - 获取某本书的所有书评
GET    /api/users/:id/reviews        - 获取某用户的所有书评

POST   /api/reviews/:id/like         - 点赞书评
DELETE /api/reviews/:id/like         - 取消点赞
GET    /api/reviews/:id/likes        - 获取书评点赞列表

POST   /api/reviews/:id/comments     - 发布评论
GET    /api/reviews/:id/comments     - 获取评论列表
DELETE /api/comments/:id             - 删除评论

POST   /api/users/:id/follow         - 关注用户
DELETE /api/users/:id/follow         - 取消关注
GET    /api/users/:id/followers      - 粉丝列表
GET    /api/users/:id/following      - 关注列表

GET    /api/feed                     - 个性化书评流
GET    /api/feed/following           - 关注用户的书评流
```

### 3. 推荐服务集成

#### 新增 HTTP Client 调用 Python 服务

```go
// internal/services/recommendation_client.go
package services

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type RecommendationClient struct {
    baseURL string
    client  *http.Client
}

// 调用 Python RAG 搜索
func (r *RecommendationClient) SearchBooks(query string, topK int) ([]Book, error) {
    url := fmt.Sprintf("%s/search?q=%s&top_k=%d", r.baseURL, query, topK)
    // HTTP GET request
}

// 调用 Python 个性化推荐
func (r *RecommendationClient) GetPersonalRecommendations(userID uint, topK int) ([]Book, error) {
    url := fmt.Sprintf("%s/recommend/personal?user_id=%d&top_k=%d", r.baseURL, userID, topK)
    // HTTP GET request
}
```

#### 新增统一 API 路由
```
GET    /api/recommendations          - 个性化图书推荐（调用 Python）
GET    /api/search                   - 语义搜索图书（调用 Python）
GET    /api/books/:isbn              - 获取图书详情（调用 Python）
POST   /api/books/:isbn/chat         - Chat with Book（调用 Python LLM）
```

## 🗂️ 文件改造清单

### Phase 1: 数据模型层
- [x] `internal/app/models/vedio.go` → `book_review.go`
- [x] `internal/app/models/comment.go` - 更新字段名（VideoID → ReviewID）
- [x] `internal/app/models/like.go` - 更新字段名
- [x] `internal/app/models/collection.go` - 更新字段名
- [x] `internal/app/models/user.go` - 更新关联关系

### Phase 2: Handler 层
- [ ] `internal/app/handlers/publish/` → `review/`
- [ ] `internal/app/handlers/feed/` → `feed/` (更新逻辑)
- [ ] `internal/app/handlers/favorite/` - 更新字段引用
- [ ] `internal/app/handlers/comment/` - 更新字段引用
- [ ] `internal/app/handlers/follow/` - 保持不变

### Phase 3: Response 结构
- [ ] `internal/app/handlers/response/publish.go` → `review.go`
- [ ] `internal/app/handlers/response/feed.go` - 更新结构
- [ ] 其他 response 文件更新字段名

### Phase 4: 服务层
- [ ] `internal/app/services/` - 添加 `recommendation_client.go`
- [ ] `internal/app/services/` - 更新业务逻辑

### Phase 5: 路由层
- [ ] `internal/server/server.go` - 重构所有路由

### Phase 6: 配置与文档
- [ ] `config/conf/config.yaml` - 添加 Python 服务配置
- [ ] `README.md` - 更新项目描述
- [ ] `docs/swagger/` - 更新 API 文档
- [ ] 数据库迁移脚本

## 🚀 实施步骤

### Step 1: 数据模型改造（当前）
1. 重命名 `vedio.go` → `book_review.go`
2. 添加 BookISBN、BookTitle、ReviewType、Rating 等字段
3. 更新所有关联模型的字段名

### Step 2: 数据库迁移
1. 创建迁移脚本
2. 重命名表：`videos_models` → `book_reviews`
3. 重命名字段：所有 `video_id` → `review_id`
4. 添加新字段

### Step 3: Handler 和 API 改造
1. 重构路由（/douyin → /api）
2. 更新所有 Handler 函数
3. 更新 Response 结构

### Step 4: 集成推荐服务
1. 创建 HTTP Client
2. 添加配置（Python 服务地址）
3. 暴露统一 API

### Step 5: 测试与文档
1. 更新单元测试
2. 更新 Swagger 文档
3. 更新 README 和架构图

## 📊 数据库迁移 SQL

```sql
-- Step 1: Rename table
ALTER TABLE videos_models RENAME TO book_reviews;

-- Step 2: Add new columns
ALTER TABLE book_reviews ADD COLUMN book_isbn VARCHAR(20);
ALTER TABLE book_reviews ADD COLUMN book_title VARCHAR(200);
ALTER TABLE book_reviews ADD COLUMN review_type VARCHAR(20) DEFAULT 'mixed';
ALTER TABLE book_reviews ADD COLUMN rating DECIMAL(3,1) DEFAULT 0.0;
ALTER TABLE book_reviews ADD COLUMN view_count INTEGER DEFAULT 0;
ALTER TABLE book_reviews ADD COLUMN content TEXT;
ALTER TABLE book_reviews RENAME COLUMN url TO video_url;

-- Step 3: Update comment table
ALTER TABLE comment_models RENAME COLUMN video_id TO review_id;

-- Step 4: Update many-to-many tables
ALTER TABLE user_like RENAME COLUMN video_id TO review_id;
ALTER TABLE user_collection RENAME COLUMN video_id TO review_id;

-- Step 5: Add indexes
CREATE INDEX idx_book_reviews_isbn ON book_reviews(book_isbn);
CREATE INDEX idx_book_reviews_author ON book_reviews(author_id);
CREATE INDEX idx_book_reviews_created ON book_reviews(created_at);
```

## 🎨 前端改造

### 从视频流 → 书评流
- 视频卡片 → 书评卡片（支持视频/文字/混合）
- 显示关联图书信息（封面、标题、作者）
- 显示用户评分（星级）
- 支持筛选（按图书、按用户、按评分）

### 新增功能
- 图书详情页（来自 Python 服务）
- 图书推荐页（来自 Python 服务）
- Chat with Book（来自 Python LLM）

## 🔗 微服务架构

```
┌─────────────────────────────────────┐
│   Frontend (React)                  │
│   - 书评流展示                      │
│   - 图书搜索与推荐                  │
│   - 社交互动                        │
└──────────┬──────────────────────────┘
           │ REST API
           ↓
┌─────────────────────────────────────┐
│   BookCommunity (Go)                │
│   Port: 8080                        │
│   - 用户系统 (JWT Auth)             │
│   - 书评 CRUD                       │
│   - 社交功能 (点赞/评论/关注)       │
│   - PostgreSQL + Redis + RabbitMQ   │
└──────────┬──────────────────────────┘
           │ HTTP Client (内部调用)
           ↓
┌─────────────────────────────────────┐
│   book-rec-with-LLMs (Python)       │
│   Port: 6006                        │
│   - RAG 语义搜索                    │
│   - 个性化推荐                      │
│   - LLM Chat                        │
│   - ChromaDB + SQLite               │
└─────────────────────────────────────┘
```

## ✅ 验收标准

- [ ] 数据模型完成重命名和字段扩展
- [ ] API 路由符合 RESTful 规范
- [ ] 成功调用 Python 推荐服务
- [ ] Swagger 文档完整
- [ ] 所有测试通过
- [ ] 前端 Demo 可以展示完整流程
- [ ] README 准确描述项目定位

## 📝 后续优化

1. **性能优化**
   - Redis 缓存热门书评
   - 数据库查询优化（索引）
   - CDN 加速图片和视频

2. **功能扩展**
   - 书评草稿箱
   - 定时发布
   - 热榜排行
   - 标签系统

3. **监控与运维**
   - Prometheus 指标
   - Grafana 仪表盘
   - 日志聚合
