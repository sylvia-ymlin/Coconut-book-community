# PostgreSQL 迁移指南

## 📋 迁移概述

BookCommunity 已从 MySQL 迁移到 **PostgreSQL 15**，享受更强大的功能和性能。

### ✅ 迁移优势

| 特性 | MySQL | PostgreSQL | 优势 |
|------|-------|------------|------|
| **ACID 合规** | ✓ | ✓✓ | 更强的事务保证 |
| **全文搜索** | 基础 | 强大 | 内置 GIN 索引，性能更好 |
| **JSON 支持** | JSON | JSONB | JSONB 索引性能更高 |
| **并发控制** | 锁 | MVCC | 读写不阻塞 |
| **扩展性** | 有限 | 丰富 | PostGIS、pg_trgm 等 |
| **开源协议** | GPL | PostgreSQL License | 更自由 |
| **欧洲市场** | 常见 | **首选** | 金融/科技行业标准 |

---

## 🚀 快速开始

### 1. 使用 Docker 启动 PostgreSQL（推荐）

```bash
# 启动 PostgreSQL 开发环境
./scripts/db-manage.sh start

# 查看日志
./scripts/db-manage.sh logs

# 进入数据库 shell
./scripts/db-manage.sh shell
```

**默认连接信息：**
```yaml
Host: localhost
Port: 5432
Database: bookcommunity
Username: bookcommunity
Password: dev_password_2024
```

### 2. 更新配置文件

复制配置模板并修改：

```bash
cp config/conf/example.yaml config/conf/config.yaml
```

编辑 `config/conf/config.yaml`：

```yaml
database:
  driver: postgres              # 使用 PostgreSQL
  username: bookcommunity
  password: dev_password_2024
  host: localhost
  port: 5432
  dbname: bookcommunity
  sslmode: disable              # 开发环境禁用 SSL
  timezone: Asia/Shanghai
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 启动应用

```bash
go run main.go
```

应用会自动：
- 连接到 PostgreSQL
- 创建所有表结构（GORM AutoMigrate）
- 应用性能优化索引

---

## 🔧 生产环境部署

### 方案 A：使用托管 PostgreSQL

**推荐服务商：**
- **AWS RDS for PostgreSQL** （全球）
- **Google Cloud SQL** （全球）
- **Azure Database for PostgreSQL** （欧洲数据中心）
- **Supabase** （免费额度，欧洲节点）
- **Neon** （Serverless PostgreSQL，欧洲支持）

**配置示例（生产环境）：**

```yaml
database:
  driver: postgres
  username: bookcommunity_prod
  password: ${DB_PASSWORD}  # 使用环境变量
  host: your-rds-instance.eu-west-1.rds.amazonaws.com
  port: 5432
  dbname: bookcommunity_prod
  sslmode: require          # 生产环境启用 SSL
  timezone: Europe/London
  max_idle_conns: 25
  max_open_conns: 200
  conn_max_lifetime: "30m"
```

### 方案 B：自托管 PostgreSQL

**使用 Docker Compose（生产）：**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: bookcommunity
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
      POSTGRES_DB: bookcommunity
    volumes:
      - postgres_data:/var/lib/postgresql/data
    secrets:
      - db_password
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G

secrets:
  db_password:
    file: ./secrets/db_password.txt

volumes:
  postgres_data:
    driver: local
```

---

## 📊 PostgreSQL 特有功能

### 1. 全文搜索（已集成）

BookCommunity 已内置 PostgreSQL 全文搜索索引：

```sql
-- 搜索视频标题（示例）
SELECT * FROM video_models
WHERE to_tsvector('english', title) @@ to_tsquery('english', 'golang & tutorial');
```

**GORM 使用示例：**

```go
// 在 handlers 中使用
db.Where("to_tsvector('english', title) @@ to_tsquery('english', ?)", "golang & tutorial").
   Find(&videos)
```

### 2. JSONB 字段（可扩展）

如果需要存储灵活的元数据：

```go
type VideoModel struct {
    gorm.Model
    Title    string
    Metadata datatypes.JSON `gorm:"type:jsonb"` // PostgreSQL JSONB
}

// 查询 JSONB
db.Where("metadata->>'author' = ?", "John").Find(&videos)
```

### 3. 数组类型

```go
type BookModel struct {
    gorm.Model
    Tags pq.StringArray `gorm:"type:text[]"` // PostgreSQL 数组
}
```

---

## 🔍 性能优化

### 已应用的索引

BookCommunity 自动应用以下索引（见 `internal/database/migrate.go`）：

```sql
-- 用户表索引
CREATE INDEX idx_users_username ON users_models(username);

-- 视频表复合索引
CREATE INDEX idx_videos_author_created ON video_models(author_id, created_at DESC);

-- 评论表索引
CREATE INDEX idx_comments_video_created ON comment_models(video_id, created_at DESC);

-- 全文搜索索引
CREATE INDEX idx_videos_title_fulltext ON video_models USING gin(to_tsvector('english', title));
```

### 查询性能分析

```sql
-- 使用 EXPLAIN ANALYZE 分析查询
EXPLAIN ANALYZE
SELECT * FROM video_models WHERE author_id = 1 ORDER BY created_at DESC LIMIT 10;
```

### 连接池优化

**推荐配置（根据实例规格调整）：**

| 实例规格 | max_open_conns | max_idle_conns | conn_max_lifetime |
|---------|----------------|----------------|-------------------|
| 小型 (2核4G) | 50 | 10 | 1h |
| 中型 (4核8G) | 100 | 25 | 30m |
| 大型 (8核16G) | 200 | 50 | 15m |

---

## 🛠️ 数据库管理

### 备份数据库

```bash
# 使用脚本备份
./scripts/db-manage.sh backup

# 手动备份
docker exec bookcommunity-postgres pg_dump -U bookcommunity -d bookcommunity > backup.sql

# 生产环境（远程）
pg_dump -h your-host.com -U bookcommunity -d bookcommunity > backup.sql
```

### 恢复数据库

```bash
# 使用脚本恢复
./scripts/db-manage.sh restore backups/backup_20240101.sql

# 手动恢复
psql -h localhost -U bookcommunity -d bookcommunity < backup.sql
```

### 数据库维护

```sql
-- 清理死行（定期执行）
VACUUM ANALYZE;

-- 重建索引
REINDEX DATABASE bookcommunity;

-- 查看数据库大小
SELECT pg_size_pretty(pg_database_size('bookcommunity'));

-- 查看表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 🐛 常见问题

### Q1: 连接失败 "connection refused"

**解决方案：**
```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 查看日志
./scripts/db-manage.sh logs

# 重启数据库
./scripts/db-manage.sh restart
```

### Q2: SSL 错误

**开发环境：**
```yaml
database:
  sslmode: disable  # 开发环境禁用 SSL
```

**生产环境：**
```yaml
database:
  sslmode: require  # 生产环境启用 SSL
```

### Q3: 性能优化建议

1. **启用慢查询日志：**
```sql
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 记录超过1秒的查询
SELECT pg_reload_conf();
```

2. **查看慢查询：**
```bash
docker exec bookcommunity-postgres tail -f /var/log/postgresql/postgresql.log
```

3. **创建合适的索引：**
```sql
-- 找出缺失索引的查询
SELECT * FROM pg_stat_user_tables WHERE idx_scan = 0 AND seq_scan > 100;
```

---

## 📚 参考资源

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/15/)
- [GORM PostgreSQL 驱动](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [PostgreSQL 性能优化指南](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [PostgreSQL vs MySQL](https://www.postgresql.org/about/)

---

## 🆘 支持

遇到问题？

1. 查看日志：`./scripts/db-manage.sh logs`
2. 检查配置：`config/conf/config.yaml`
3. 重置数据库：`./scripts/db-manage.sh reset`

**欧洲 PostgreSQL 社区：**
- [PostgreSQL Europe](https://www.postgresql.eu/)
- [PGConf.EU](https://www.postgresql.eu/events/pgconfeu2024/)
