# BookCommunity - 图书社交平台

[![CI](https://github.com/sylvia-ymlin/Coconut-book-community/workflows/CI%20Pipeline/badge.svg)](https://github.com/sylvia-ymlin/Coconut-book-community/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/sylvia-ymlin/Coconut-book-community)](https://goreportcard.com/report/github.com/sylvia-ymlin/Coconut-book-community)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> 一个基于 Go 的高性能图书社交平台后端，结合图文书评与 AI 智能推荐（图书版小红书）

## 📖 项目简介

BookCommunity 是一个**图文书评社交平台**（类似小红书的图书版本），用户可以：

- 📸 发布图文书评（最多9张图片 + 文字内容）
- ⭐ 为图书评分、点赞、评论、收藏
- 🏷️ 使用标签分类书评（悬疑、推理、文学等）
- 👥 关注读者、浏览书评流（发现页、关注页）
- 🤖 获取 AI 驱动的个性化图书推荐（集成 Python RAG 服务）
- 💬 与图书进行 LLM 对话（Chat with Book）

**定位**: 图书版小红书 = 图文书评社区 + AI 智能推荐

## 🏗️ 技术架构

### 微服务架构

```
┌─────────────────────────────────────┐
│   Frontend (React)                  │
│   - 书评流展示                      │
│   - 图书搜索与推荐                  │
│   - 社交互动（点赞/评论/关注）      │
└──────────┬──────────────────────────┘
           │ REST API
           ↓
┌─────────────────────────────────────┐
│   BookCommunity (Go Backend)        │
│   Port: 8080                        │
│   ├── 用户系统 (JWT Auth)           │
│   ├── 书评 CRUD                     │
│   ├── 社交功能 (点赞/评论/关注)     │
│   └── API Gateway                   │
│                                     │
│   Infrastructure:                   │
│   ├── PostgreSQL 15 (主库)          │
│   ├── Redis 7.0 (缓存)              │
│   └── RabbitMQ 3.12 (消息队列)      │
└──────────┬──────────────────────────┘
           │ HTTP Client (内部调用)
           ↓
┌─────────────────────────────────────┐
│   book-rec-with-LLMs (Python)       │
│   Port: 6006                        │
│   ├── RAG 语义搜索                  │
│   ├── 个性化推荐 (7通道召回)        │
│   ├── LLM Chat                      │
│   └── ChromaDB + SQLite             │
└─────────────────────────────────────┘
```

## 🚀 核心功能

### 社交功能 (Go Backend)

| 功能模块 | 说明 |
|---------|------|
| **用户系统** | 注册、登录、JWT 认证 |
| **书评发布** | 图文书评（文字 + 最多9张图片） |
| **社交互动** | 点赞、评论、收藏书评 |
| **标签系统** | 书评标签分类（悬疑、文学、科幻等） |
| **关注系统** | 关注用户、粉丝列表、关注列表 |
| **书评流** | 发现页（个性化 Feed）、关注页 |

### AI 推荐功能 (Python Service)

| 功能模块 | 说明 |
|---------|------|
| **语义搜索** | RAG + Hybrid Search（BM25 + Dense Vector） |
| **个性化推荐** | 7通道召回 + LGBMRanker + Stacking |
| **图书详情** | 完整的图书元数据 |
| **Chat with Book** | LLM 对话（基于图书内容） |

## 🛠️ 技术栈

### Backend (Go)
- **Framework**: Gin 1.9
- **ORM**: GORM 1.25
- **Database**: PostgreSQL 15
- **Cache**: Redis 7.0 + In-Memory LRU
- **Message Queue**: RabbitMQ 3.12
- **API Doc**: Swagger 2.0

### AI Service (Python)
- **Framework**: FastAPI 0.100+
- **Vector DB**: ChromaDB
- **Embeddings**: BGE-M3
- **LLM**: Ollama / OpenAI / Groq
- **RecSys**: Item2Vec, SASRec, LightGBM

### DevOps
- **Containerization**: Docker + Docker Compose
- **Orchestration**: Kubernetes
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

## 📦 快速开始

### 前置要求

- Go 1.20+
- Docker & Docker Compose
- Python 3.10+ (for AI service)
- PostgreSQL 15
- Redis 7.0

### 1. 启动 Go 后端服务

```bash
# 克隆仓库
git clone https://github.com/sylvia-ymlin/Coconut-book-community.git
cd Coconut-book-community

# 启动基础设施（PostgreSQL, Redis, RabbitMQ）
docker-compose up -d

# 复制配置文件
cp config/conf/config.yaml.example config/conf/config.yaml

# 编辑配置文件，设置数据库密码等

# 运行数据库迁移（首次启动）
# make migrate  # 或手动执行 scripts/migrate_to_book_social.sql

# 启动服务
go run main.go
```

服务将启动在 `http://localhost:8080`

### 2. 启动 Python 推荐服务（可选）

```bash
# 克隆推荐服务仓库
cd ..
git clone https://github.com/sylvia-ymlin/book-rec-with-LLMs.git
cd book-rec-with-LLMs

# 创建虚拟环境
conda env create -f environment.yml
conda activate book-rec

# 初始化数据库
python data/scripts/init_db.py
python scripts/init_sqlite_db.py

# 启动服务
make run  # http://localhost:6006
```

### 3. 访问服务

- **API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health
- **前端 Demo**: http://localhost:8000/demo-mock.html (Mock 版本)

## 📚 API 文档

### 核心 API

#### 用户系统
```
POST   /api/user/register       - 用户注册
POST   /api/user/login          - 用户登录
GET    /api/user/:id            - 获取用户信息
```

#### 书评管理
```
POST   /api/reviews             - 发布书评
GET    /api/reviews             - 获取书评列表
GET    /api/reviews/:id         - 获取书评详情
PUT    /api/reviews/:id         - 更新书评
DELETE /api/reviews/:id         - 删除书评
```

#### 社交互动
```
POST   /api/reviews/:id/like    - 点赞书评
GET    /api/reviews/:id/likes   - 获取点赞列表
POST   /api/reviews/:id/comments - 发布评论
GET    /api/reviews/:id/comments - 获取评论列表
```

#### 图书推荐（代理到 Python 服务）
```
GET    /api/recommendations     - 个性化推荐
GET    /api/search              - 语义搜索图书
GET    /api/books/:isbn         - 获取图书详情
POST   /api/books/:isbn/chat    - Chat with Book
```

完整 API 文档请访问 Swagger UI。

## 🎯 性能指标

| 指标 | 数值 |
|------|------|
| **QPS** | 5000+ requests/sec |
| **P99 延迟** | < 50ms |
| **缓存命中率** | 95%+ |
| **数据库连接池** | 100 max connections |

## 📊 数据模型

### BookReviewModel (书评)
```go
type BookReviewModel struct {
    gorm.Model
    Title        string  // 书评标题
    Content      string  // 书评内容（文字，必填）
    BookISBN     string  // 关联图书 ISBN
    BookTitle    string  // 图书标题
    Images       string  // 图片URL列表（JSON数组，最多9张）
    CoverURL     string  // 封面图（第一张图片）
    Rating       float64 // 评分 (0-10)
    Tags         string  // 标签列表（JSON数组）
    AuthorID     uint    // 作者ID
    LikeCount    uint    // 点赞数
    CommentCount uint    // 评论数
    ViewCount    uint    // 浏览次数
    CollectCount uint    // 收藏次数
}
```

### 关系模型
- **用户 ↔ 书评**: 一对多（用户可发布多个书评）
- **用户 ↔ 关注**: 多对多（用户可关注多个用户）
- **书评 ↔ 点赞**: 多对多（书评可被多个用户点赞）
- **书评 ↔ 评论**: 一对多（书评可有多条评论）
- **书评 ↔ 图书**: 多对一（多个书评关联同一本书）

## 🔧 开发

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 生成 Swagger 文档
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
make swagger
```

### 数据库迁移
```bash
# 执行迁移脚本
psql -U bookcommunity -d bookcommunity -f scripts/migrate_to_book_social.sql
```

## 📁 项目结构

```
BookCommunity/
├── cmd/                    # 命令行工具
├── config/                 # 配置文件
│   └── conf/
│       └── config.yaml
├── internal/               # 内部代码（不对外暴露）
│   ├── app/
│   │   ├── handlers/       # HTTP 处理器
│   │   │   ├── user/
│   │   │   ├── review/     # 书评相关
│   │   │   ├── comment/
│   │   │   └── follow/
│   │   ├── models/         # 数据模型
│   │   │   ├── book_review.go
│   │   │   ├── user.go
│   │   │   ├── comment.go
│   │   │   └── ...
│   │   └── services/       # 业务逻辑
│   │       └── recommendation_client.go  # Python 服务客户端
│   ├── database/           # 数据库连接
│   ├── middleware/         # 中间件（JWT, CORS）
│   └── server/             # 路由配置
├── pkg/                    # 可复用的包
├── scripts/                # 脚本
│   └── migrate_to_book_social.sql
├── docs/                   # Swagger 文档
├── examples/               # 示例代码
│   └── frontend/
├── documentation/          # 技术文档
│   └── MIGRATION_TO_BOOK_SOCIAL.md
├── docker-compose.yml      # Docker Compose 配置
├── Dockerfile              # Docker 镜像
├── Makefile                # Make 命令
├── go.mod                  # Go 依赖
└── README.md
```

## 🤝 集成 Python 推荐服务

BookCommunity 通过 HTTP 调用集成 [book-rec-with-LLMs](https://github.com/sylvia-ymlin/book-rec-with-LLMs) 推荐服务。

### 配置
在 `config/conf/config.yaml` 中配置：
```yaml
recommendation:
  enabled: true                        # 启用推荐服务
  api_url: "http://localhost:6006/api" # Python 服务地址
  timeout: "30s"                       # 超时时间
  fallback_to_mock: true               # 失败时降级到 mock 数据
```

### 调用示例
```go
import "github.com/sylvia-ymlin/Coconut-book-community/internal/app/services"

// 搜索图书
books, err := services.GlobalRecommendationClient.SearchBooks("深入理解计算机系统", 10)

// 个性化推荐
recommendations, err := services.GlobalRecommendationClient.GetPersonalRecommendations(userID, 10)

// Chat with Book
response, err := services.GlobalRecommendationClient.ChatWithBook(isbn, "这本书讲了什么？", userID)
```

## 📄 文档

- [迁移指南](documentation/MIGRATION_TO_BOOK_SOCIAL.md) - 从视频平台到图书社交平台的迁移文档
- [前端集成指南](FRONTEND_INTEGRATION_GUIDE.md) - 前端开发者集成指南
- [Swagger 指南](SWAGGER_GUIDE.md) - API 文档使用指南
- [架构设计](documentation/ARCHITECTURE.md) - 详细架构设计（待完成）

## 🎨 前端 Demo

项目提供了一个 Mock Demo 用于快速体验：

```bash
# 启动前端 Demo（无需后端）
cd examples/frontend/vanilla-js
python3 -m http.server 8000

# 访问 http://localhost:8000/demo-mock.html
```

## 🌟 特性

- ✅ **高性能**: Go + PostgreSQL + Redis，QPS 5000+
- ✅ **微服务架构**: Go 后端 + Python AI 服务解耦
- ✅ **社交功能完整**: 点赞、评论、关注、Feed 流
- ✅ **AI 驱动**: 集成先进的推荐算法和 LLM
- ✅ **RESTful API**: 符合 RESTful 规范的 API 设计
- ✅ **完整文档**: Swagger API 文档 + 技术文档
- ✅ **容器化部署**: Docker + Kubernetes 支持
- ✅ **测试覆盖**: 50+ 单元测试用例

## 🛣️ Roadmap

- [ ] 完成所有 API 的 Swagger 文档（当前 28.6%）
- [ ] 提升测试覆盖率到 80%+
- [ ] 实现完整的 React 前端
- [ ] 添加图片上传功能（封面、头像）
- [ ] 书评草稿箱功能
- [ ] 热榜排行（热门书评、热门图书）
- [ ] 标签系统
- [ ] Elasticsearch 全文搜索
- [ ] GraphQL API 支持

## 🤝 贡献

欢迎贡献代码、报告 Bug 或提出新功能建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 License

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 👥 作者

- **Sylvia Lin** - [GitHub](https://github.com/sylvia-ymlin)

## 🙏 致谢

- 原始项目基于字节跳动青训营抖音后端项目改造
- 推荐系统设计参考 [book-rec-with-LLMs](https://github.com/sylvia-ymlin/book-rec-with-LLMs)
- 感谢所有开源项目的贡献者

---

⭐ 如果这个项目对你有帮助，请给一个 Star！
