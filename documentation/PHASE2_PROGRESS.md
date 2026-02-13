# Phase 2 Handler 层改造进度

## 🎯 目标

将所有 Handler 从 Video 相关改造为 BookReview 相关，采用 RESTful 风格的 API 设计。

## ✅ 已完成工作

### 1. Response 结构体 ✅

**文件**: `internal/app/handlers/response/review.go`

创建了完整的书评响应结构：

```go
// 核心结构
type ReviewResponse struct {
    CommonResponse
    Review *ReviewInfo
}

type ReviewListResponse struct {
    CommonResponse
    Reviews  []*ReviewInfo
    NextTime int64  // 分页时间戳
    Total    int64
}

type FeedResponse struct {
    CommonResponse
    Reviews  []*ReviewInfo
    NextTime int64
    HasMore  bool
}

// 详细信息
type ReviewInfo struct {
    ID, Title, Content, BookISBN, BookTitle
    Images       []string  // 图片列表（已解析）
    Tags         []string  // 标签列表（已解析）
    Rating       float64
    Author       *UserInfo
    LikeCount, CommentCount, ViewCount, CollectCount
    IsLiked, IsCollected  // 当前用户关系
}

// 请求结构
type CreateReviewRequest struct {
    Title, Content, BookISBN, BookTitle
    Images    []string  // 最多9张
    Rating    float64
    Tags      []string  // 最多10个
}

type UpdateReviewRequest struct {
    // 所有字段可选（指针类型）
    Title, Content, Images, Rating, Tags
}
```

### 2. 书评 Handler ✅

#### 2.1 创建书评
**文件**: `internal/app/handlers/review/create.go`

```go
POST /api/reviews
```

**功能**:
- ✅ JWT 认证（从中间件获取 user_id）
- ✅ 参数验证（title, content 必填，images 最多9张）
- ✅ JSON 序列化（Images, Tags 数组 → JSON字符串）
- ✅ 自动设置封面图（第一张图片）
- ✅ 保存到数据库
- ✅ 返回完整书评信息
- ⏳ TODO: 异步任务（更新用户统计、触发推荐更新、通知关注者）

#### 2.2 查询书评列表
**文件**: `internal/app/handlers/review/list.go`

```go
GET /api/reviews?page=1&page_size=20&user_id=1&book_isbn=xxx&order_by=latest
```

**功能**:
- ✅ 支持分页（page, page_size）
- ✅ 支持筛选（user_id, book_isbn）
- ✅ 支持排序（latest, popular, rating）
  - `latest`: 按时间倒序
  - `popular`: 按热度排序（点赞*3 + 评论*2 + 收藏*2 + 浏览数）
  - `rating`: 按评分排序
- ✅ 预加载作者信息
- ✅ 返回总数
- ⏳ TODO: 判断当前用户的点赞/收藏状态

#### 2.3 查询书评详情
**文件**: `internal/app/handlers/review/list.go`

```go
GET /api/reviews/:id
```

**功能**:
- ✅ 通过 ID 查询单条书评
- ✅ 预加载作者信息
- ✅ 异步更新浏览次数（view_count + 1）
- ✅ 返回完整书评信息
- ⏳ TODO: 判断当前用户的点赞/收藏状态

#### 2.4 更新书评
**文件**: `internal/app/handlers/review/update.go`

```go
PUT /api/reviews/:id
```

**功能**:
- ✅ JWT 认证（必须登录）
- ✅ 权限检查（只能更新自己的书评）
- ✅ 部分更新（只更新传入的字段）
- ✅ 自动更新封面图（如果更新了 images）
- ✅ 返回更新后的书评
- ⏳ TODO: 异步任务（可选）

#### 2.5 删除书评
**文件**: `internal/app/handlers/review/update.go`

```go
DELETE /api/reviews/:id
```

**功能**:
- ✅ JWT 认证（必须登录）
- ✅ 权限检查（只能删除自己的书评）
- ✅ 软删除（GORM DeletedAt）
- ⏳ TODO: 异步删除相关数据（点赞、评论、收藏记录）

### 3. Feed 流 Handler ✅

#### 3.1 发现页（个性化推荐）
**文件**: `internal/app/handlers/review/feed.go`

```go
GET /api/feed?latest_time=0&page_size=20
```

**功能**:
- ✅ 基于时间戳的分页（latest_time）
- ✅ 综合热度排序
  - 热度分 = 点赞数*3 + 评论数*2 + 收藏数*2
  - 综合时间倒序
- ✅ 返回 hasMore 标志
- ✅ 返回 nextTime（下一页时间戳）
- ⏳ TODO: 真正的个性化推荐（基于用户兴趣）

#### 3.2 关注页
**文件**: `internal/app/handlers/review/feed.go`

```go
GET /api/feed/following?latest_time=0&page_size=20
```

**功能**:
- ✅ JWT 认证（必须登录）
- ✅ 查询当前用户关注的用户列表
- ✅ 获取关注用户的书评（按时间倒序）
- ✅ 基于时间戳的分页
- ✅ 如果没有关注任何人，返回友好提示

## 📊 API 设计对比

### 改造前（Douyin 风格）
```
POST   /douyin/publish/action/       → 发布视频
GET    /douyin/publish/list/         → 查询发布列表
GET    /douyin/feed                  → 视频流
POST   /douyin/favorite/action/      → 点赞
```

### 改造后（RESTful 风格）
```
POST   /api/reviews                  → 创建书评
GET    /api/reviews                  → 查询书评列表
GET    /api/reviews/:id              → 查询书评详情
PUT    /api/reviews/:id              → 更新书评
DELETE /api/reviews/:id              → 删除书评
GET    /api/feed                     → 发现页
GET    /api/feed/following           → 关注页
```

**改进**:
- ✅ 符合 RESTful 规范
- ✅ URL 语义清晰
- ✅ HTTP 方法明确（GET/POST/PUT/DELETE）
- ✅ 资源路径一致（/api/reviews）

## 📁 文件结构

```
internal/app/handlers/
├── response/
│   ├── common.go           # 通用响应
│   ├── review.go           # 书评响应 ✅ 新增
│   ├── user.go             # 用户响应
│   └── ...
├── review/                 # ✅ 新增目录
│   ├── create.go           # ✅ 创建书评
│   ├── list.go             # ✅ 查询书评列表/详情
│   ├── update.go           # ✅ 更新/删除书评
│   └── feed.go             # ✅ Feed 流
├── user/
│   ├── register.go
│   └── user.go
├── comment/
│   └── comment.go          # ⏳ 待更新
├── favorite/
│   └── favorite.go         # ⏳ 待更新
└── follow/
    └── follow.go           # ✅ 无需改动
```

## ⏳ 待完成工作

### Phase 2.2: 更新社交 Handler

#### 1. 点赞 Handler
**文件**: `internal/app/handlers/like/` (新建)

需要改造：
- `/douyin/favorite/action/` → `/api/reviews/:id/like` (POST)
- `/douyin/favorite/list/` → `/api/reviews/:id/likes` (GET)

**改动**:
- `VideoID` → `ReviewID`
- 更新响应结构

#### 2. 评论 Handler
**文件**: `internal/app/handlers/comment/comment.go`

需要改造：
- `/douyin/comment/action/` → `/api/reviews/:id/comments` (POST)
- `/douyin/comment/list/` → `/api/reviews/:id/comments` (GET)
- 删除评论 → `/api/comments/:id` (DELETE)

**改动**:
- `VideoID` → `ReviewID`
- 更新响应结构

#### 3. 收藏 Handler
**文件**: `internal/app/handlers/collect/` (新建)

新增功能：
- `POST   /api/reviews/:id/collect` - 收藏书评
- `DELETE /api/reviews/:id/collect` - 取消收藏
- `GET    /api/users/:id/collections` - 查询收藏列表

### Phase 2.3: 路由重构

**文件**: `internal/server/server.go`

需要重构所有路由，采用 RESTful 风格：

```go
// 书评相关
reviewGroup := baseGroup.Group("/reviews")
{
    reviewGroup.POST("", middleware.JWTMiddleware(), review.CreateReviewHandler)
    reviewGroup.GET("", review.GetReviewListHandler)
    reviewGroup.GET("/:id", review.GetReviewDetailHandler)
    reviewGroup.PUT("/:id", middleware.JWTMiddleware(), review.UpdateReviewHandler)
    reviewGroup.DELETE("/:id", middleware.JWTMiddleware(), review.DeleteReviewHandler)

    // 书评的点赞
    reviewGroup.POST("/:id/like", middleware.JWTMiddleware(), like.LikeReviewHandler)
    reviewGroup.DELETE("/:id/like", middleware.JWTMiddleware(), like.UnlikeReviewHandler)
    reviewGroup.GET("/:id/likes", like.GetReviewLikesHandler)

    // 书评的评论
    reviewGroup.POST("/:id/comments", middleware.JWTMiddleware(), comment.CreateCommentHandler)
    reviewGroup.GET("/:id/comments", comment.GetCommentListHandler)

    // 书评的收藏
    reviewGroup.POST("/:id/collect", middleware.JWTMiddleware(), collect.CollectReviewHandler)
    reviewGroup.DELETE("/:id/collect", middleware.JWTMiddleware(), collect.UncollectReviewHandler)
}

// Feed 流
feedGroup := baseGroup.Group("/feed")
{
    feedGroup.GET("", review.GetDiscoveryFeedHandler)  // 发现页
    feedGroup.GET("/following", middleware.JWTMiddleware(), review.GetFollowingFeedHandler)  // 关注页
}

// 评论（独立资源）
commentGroup := baseGroup.Group("/comments")
{
    commentGroup.DELETE("/:id", middleware.JWTMiddleware(), comment.DeleteCommentHandler)
}
```

## 🎯 设计亮点

### 1. RESTful 风格
- 资源清晰：`/reviews`, `/comments`, `/likes`
- HTTP 方法语义明确：GET（查询）、POST（创建）、PUT（更新）、DELETE（删除）
- URL 层级合理：`/reviews/:id/comments`

### 2. 分页策略
- **列表分页**: 传统的 page/page_size
- **Feed 流分页**: 基于时间戳的 latest_time（适合实时更新）

### 3. 权限控制
- JWT 认证（必须登录才能创建、更新、删除）
- 权限检查（只能操作自己的资源）

### 4. 性能优化
- 预加载作者信息（`Preload("Author")`）
- 异步更新浏览次数（`go func()`）
- 索引优化（author_id, book_isbn, created_at）

### 5. 用户体验
- 友好的错误提示
- 热度排序算法（综合点赞、评论、收藏）
- hasMore 标志（前端无限滚动）

## 📝 下一步

1. **完成 Phase 2.2** - 更新点赞、评论、收藏 Handler
2. **完成 Phase 2.3** - 重构路由配置
3. **测试 API** - 确保所有接口正常工作
4. **Phase 3** - 添加图片上传功能
5. **Phase 4** - 更新 Swagger 文档

## 🎉 总结

Phase 2.1 已完成，创建了完整的书评 CRUD 和 Feed 流 Handler：

- ✅ Response 结构体（ReviewResponse, ReviewListResponse, FeedResponse）
- ✅ 创建书评 Handler（JWT认证、参数验证）
- ✅ 查询书评 Handler（分页、筛选、排序）
- ✅ 更新/删除书评 Handler（权限检查）
- ✅ Feed 流 Handler（发现页、关注页）

**代码质量**:
- ✅ 符合 RESTful 规范
- ✅ 完整的错误处理
- ✅ 清晰的注释和日志
- ✅ Swagger 注释（待生成文档）

**下一步**: 完成社交 Handler（点赞、评论、收藏）的改造。
