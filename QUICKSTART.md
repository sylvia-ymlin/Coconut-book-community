# BookCommunity 快速启动指南

## 🚀 3分钟快速启动

### 前置要求

- Docker Desktop (推荐) 或 Docker + Docker Compose
- Go 1.20+（如果要本地编译）

---

## 方式一：Docker 一键启动（推荐）

### 1. 启动完整技术栈

```bash
# 克隆项目后进入目录
cd bookcommunity

# 启动所有服务（PostgreSQL + Redis + RabbitMQ + 监控）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

**启动的服务：**
- ✅ PostgreSQL 15 (端口 5432)
- ✅ Redis 7.0 (端口 6379)
- ✅ RabbitMQ 3.12 (端口 5672, 管理界面 15672)
- ✅ Prometheus (端口 9090)
- ✅ Grafana (端口 3000)

### 2. 配置应用

```bash
# 复制配置模板
cp config/conf/example.yaml config/conf/config.yaml

# 编辑配置（可选，默认配置已可用）
vim config/conf/config.yaml
```

**默认配置摘要：**
```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  dbname: bookcommunity
  username: bookcommunity
  password: secure_password_2024

redis:
  enabled: true
  host: localhost
  port: 6379
  password: redis_password_2024
```

### 3. 运行应用

```bash
# 安装依赖
go mod tidy

# 编译运行
go run main.go

# 或者先编译再运行
go build -o bookcommunity main.go
./bookcommunity
```

**应用启动成功标志：**
```
✅ PostgreSQL connected successfully: localhost:5432 (DB 0)
✅ Redis connected successfully: localhost:6379 (DB 0)
✅ Hybrid cache initialized (Redis: true, Memory LRU: 5000)
✅ UserCacheService initialized
[GIN-debug] Listening and serving HTTP on :8080
```

### 4. 访问服务

| 服务 | URL | 用户名/密码 |
|------|-----|------------|
| **应用健康检查** | http://localhost:8080/health | - |
| **RabbitMQ 管理** | http://localhost:15672 | bookcommunity / rabbitmq_password_2024 |
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin / admin |

---

## 方式二：开发环境（包含管理工具）

```bash
# 启动开发环境（包含 pgAdmin + Redis Commander）
docker-compose --profile dev up -d

# 额外可访问：
# pgAdmin: http://localhost:5050
# Redis Commander: http://localhost:8081
```

---

## 📋 测试 API

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

**期望输出：**
```json
{
  "status": "healthy",
  "service": "BookCommunity API"
}
```

### 2. 注册用户

```bash
curl -X POST http://localhost:8080/douyin/user/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/douyin/user/login/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**保存返回的 token 用于后续请求**

### 4. 获取图书推荐

```bash
TOKEN="your_token_here"

curl -X GET "http://localhost:8080/douyin/recommend?top_k=5" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. 搜索图书

```bash
curl -X GET "http://localhost:8080/douyin/search?keyword=golang&limit=10"
```

---

## 🔍 管理服务

### 数据库管理

```bash
# 使用管理脚本
./scripts/db-manage.sh start      # 启动数据库
./scripts/db-manage.sh shell      # 进入 psql shell
./scripts/db-manage.sh backup     # 备份数据库
./scripts/db-manage.sh logs       # 查看日志
./scripts/db-manage.sh pgadmin    # 启动 pgAdmin

# 或直接使用 Docker
docker exec -it bookcommunity-postgres psql -U bookcommunity -d bookcommunity
```

### Redis 管理

```bash
# 进入 Redis CLI
docker exec -it bookcommunity-redis redis-cli -a redis_password_2024

# 常用命令
PING              # 测试连接
KEYS user:*       # 查看用户相关缓存
GET user:1        # 获取用户缓存
FLUSHDB           # 清空数据库（谨慎！）
```

### RabbitMQ 管理

访问管理界面：http://localhost:15672

**默认登录：**
- 用户名：`bookcommunity`
- 密码：`rabbitmq_password_2024`

---

## 📊 监控与可观测性

### Prometheus

访问：http://localhost:9090

**查询示例：**
```promql
# HTTP 请求总数
bookcommunity_http_requests_total

# 请求延迟 P95
histogram_quantile(0.95, bookcommunity_http_request_duration_seconds_bucket)

# 缓存命中率
bookcommunity_cache_hits_total / (bookcommunity_cache_hits_total + bookcommunity_cache_misses_total)
```

### Grafana

访问：http://localhost:3000

**默认登录：**
- 用户名：`admin`
- 密码：`admin`（首次登录后修改）

**添加数据源：**
1. Configuration → Data Sources → Add data source
2. 选择 Prometheus
3. URL: `http://prometheus:9090`
4. Save & Test

---

## 🛑 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（重置所有数据）
docker-compose down -v

# 停止单个服务
docker-compose stop postgres
```

---

## 🐛 常见问题

### Q1: 端口冲突

**错误：** `Bind for 0.0.0.0:5432 failed: port is already allocated`

**解决方案：**
```bash
# 查看占用端口的进程
lsof -i :5432

# 修改 docker-compose.yaml 端口映射
ports:
  - "15432:5432"  # 改为 15432
```

### Q2: Docker 连接失败

**错误：** `Cannot connect to Docker daemon`

**解决方案：**
```bash
# macOS/Windows: 确保 Docker Desktop 正在运行
# Linux: 启动 Docker 服务
sudo systemctl start docker
```

### Q3: 数据库连接失败

**错误：** `dial tcp [::1]:5432: connect: connection refused`

**解决方案：**
```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 查看日志
docker logs bookcommunity-postgres

# 重启数据库
docker-compose restart postgres
```

### Q4: Redis 认证失败

**错误：** `NOAUTH Authentication required`

**解决方案：**
检查 `config/conf/config.yaml` 中的密码配置：
```yaml
redis:
  password: redis_password_2024  # 确保与 docker-compose.yaml 一致
```

---

## 📚 进阶使用

### 性能调优

**PostgreSQL 连接池：**
```yaml
database:
  max_open_conns: 200      # 根据实例规格调整
  max_idle_conns: 50
  conn_max_lifetime: "30m"
```

**Redis 缓存策略：**
```yaml
redis:
  pool_size: 200           # 高并发场景
  default_expiration: "15m" # 热数据15分钟
```

### 生产环境部署

1. **使用环境变量存储敏感信息**
   ```bash
   export DB_PASSWORD=your_secure_password
   export REDIS_PASSWORD=your_redis_password
   ```

2. **启用 SSL**
   ```yaml
   database:
     sslmode: require  # 生产环境必须启用
   ```

3. **配置备份策略**
   ```bash
   # 定时备份数据库
   0 2 * * * /path/to/scripts/db-manage.sh backup
   ```

---

## 🎯 下一步

1. ✅ 启动服务 → http://localhost:8080/health
2. 📖 阅读文档 → `docs/`
3. 🔨 测试 API → Postman 或 cURL
4. 📊 查看监控 → Grafana 仪表盘
5. 🚀 开始开发 → 修改代码，热重载

---

## 📞 获取帮助

**文档位置：**
- PostgreSQL: `docs/POSTGRESQL_MIGRATION.md`
- Redis: `docs/REDIS_GUIDE.md`
- 总结: `docs/FINAL_SUMMARY.md`

**管理脚本：**
```bash
./scripts/db-manage.sh help
```

**检查服务状态：**
```bash
docker-compose ps
docker-compose logs -f [service_name]
```

---

**享受编码！🚀**
