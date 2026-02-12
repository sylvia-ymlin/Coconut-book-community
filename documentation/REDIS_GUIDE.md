# Redis 集成指南

## 📋 概述

BookCommunity 已集成 **Redis 7.0** 作为分布式缓存层，显著提升性能并支持高并发场景。

### ✅ Redis 优势

| 特性 | 内存LRU | Redis | 优势 |
|------|---------|-------|------|
| **持久化** | ❌ | ✅ | 重启不丢数据 |
| **分布式** | ❌ | ✅ | 多实例共享缓存 |
| **数据结构** | Key-Value | String/Hash/Set/ZSet/List | 丰富的数据类型 |
| **过期策略** | 手动 | 自动 | 自动清理过期键 |
| **原子操作** | ❌ | ✅ | 计数器、分布式锁 |
| **性能** | 极快 | 快 | 内存级访问 |

---

## 🚀 快速开始

### 1. 启动 Redis（Docker）

```bash
# 启动 Redis 开发环境
docker-compose -f docker-compose.dev.yaml up -d redis

# 查看日志
docker logs -f bookcommunity-redis

# 进入 Redis CLI
docker exec -it bookcommunity-redis redis-cli -a dev_redis_2024
```

**默认连接信息：**
```yaml
Host: localhost
Port: 6379
Password: dev_redis_2024
Database: 0
```

### 2. 配置 Redis

编辑 `config/conf/config.yaml`：

```yaml
redis:
  enabled: true                 # 启用Redis（false时自动降级到内存缓存）
  host: localhost
  port: 6379
  password: dev_redis_2024
  db: 0                         # 数据库索引 (0-15)
  pool_size: 100                # 连接池大小
  min_idle_conns: 10            # 最小空闲连接
  max_retries: 3                # 最大重试次数
  default_expiration: "1h"      # 默认过期时间
```

### 3. 使用 Redis Commander（可选）

访问 Web 管理界面：
```
http://localhost:8081
```

---

## 💻 代码示例

### 基础用法

```go
package main

import (
    "time"
    "github.com/Doraemonkeys/douyin2/internal/cache"
)

func main() {
    // 初始化混合缓存（Redis + 内存LRU）
    cache.InitHybridCache(5000) // 5000个条目的内存缓存

    // 设置缓存
    user := map[string]interface{}{
        "id":   1,
        "name": "John",
    }
    cache.Set("user:1", user, 10*time.Minute)

    // 获取缓存
    var cachedUser map[string]interface{}
    if err := cache.Get("user:1", &cachedUser); err == nil {
        fmt.Println("Cache hit:", cachedUser)
    }

    // 删除缓存
    cache.Delete("user:1")
}
```

### 用户缓存服务（推荐）

```go
import "github.com/Doraemonkeys/douyin2/internal/app/services"

// 获取用户（带缓存）
userService := services.GetUserCacheService()

user, err := userService.GetUserByID(123)
if err != nil {
    return err
}

// 更新用户后，使缓存失效
userService.InvalidateUserCache(123, "username")
```

### 计数器操作

```go
import "github.com/Doraemonkeys/douyin2/internal/cache"

// 增加点赞数
count, err := cache.Incr("video:123:like_count")
fmt.Println("New like count:", count)

// 减少点赞数
count, err = cache.Decr("video:123:like_count")

// 增加指定值
count, err = cache.IncrBy("video:123:view_count", 10)
```

### 集合操作（关注/点赞）

```go
import "github.com/Doraemonkeys/douyin2/internal/cache"

// 添加关注
cache.SAdd("user:1:following", 2, 3, 4)

// 检查是否已关注
isFollowing, _ := cache.SIsMember("user:1:following", 3)

// 获取关注数量
count, _ := cache.SCard("user:1:following")

// 获取所有关注的用户ID
following, _ := cache.SMembers("user:1:following")
```

### 有序集合（排行榜）

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/Doraemonkeys/douyin2/internal/cache"
)

// 添加到排行榜（分数 = 点赞数）
cache.ZAdd("video:hot_ranking", redis.Z{
    Score:  100,
    Member: "video:123",
})

// 获取排名前10的视频
top10, _ := cache.ZRevRange("video:hot_ranking", 0, 9)
```

---

## 🏗️ 架构设计

### 混合缓存策略

BookCommunity 使用 **L1 (内存LRU) + L2 (Redis)** 双层缓存：

```
用户请求
    ↓
L1缓存（内存LRU）
    ├─ 命中 → 立即返回（~1μs）
    └─ 未命中
         ↓
    L2缓存（Redis）
         ├─ 命中 → 回填L1 → 返回（~1ms）
         └─ 未命中
              ↓
         数据库查询
              ↓
         写入L2和L1 → 返回（~10ms）
```

### 缓存失效策略

1. **主动失效**（数据更新时）
   ```go
   // 更新用户信息后
   db.Save(&user)
   userService.InvalidateUserCache(user.ID, user.Username)
   ```

2. **被动失效**（过期时间）
   ```go
   // 热数据：15分钟
   cache.Set(key, value, 15*time.Minute)

   // 冷数据：1小时
   cache.Set(key, value, 1*time.Hour)
   ```

3. **LRU淘汰**（内存不足时）
   - 自动淘汰最少使用的数据

---

## 📊 性能优化

### 缓存键设计规范

```go
// 用户相关
user:{userID}                    // 用户基本信息
user:name:{username}             // 用户名索引
user:{userID}:follower_count     // 关注者计数
user:{userID}:following          // 关注列表（Set）

// 视频相关
video:{videoID}                  // 视频信息
video:{videoID}:like_count       // 点赞数
video:{videoID}:likers           // 点赞用户列表（Set）

// 排行榜
video:hot_ranking                // 热门视频排行（ZSet）
user:active_ranking              // 活跃用户排行（ZSet）
```

### 批量操作优化

```go
// ❌ 错误：N次查询
for _, userID := range userIDs {
    cache.Get(fmt.Sprintf("user:%d", userID), &user)
}

// ✅ 正确：1次批量查询
keys := make([]string, len(userIDs))
for i, id := range userIDs {
    keys[i] = fmt.Sprintf("user:%d", id)
}
results, _ := cache.MGet(keys)
```

### Pipeline（批量写入）

```go
// 使用原生Redis客户端实现Pipeline
client := cache.GetRedisClient()
pipe := client.Pipeline()

for i := 0; i < 100; i++ {
    pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, time.Hour)
}

_, err := pipe.Exec(ctx)
```

---

## 🛠️ 生产环境部署

### 方案A：托管Redis服务

**推荐服务商：**
- **AWS ElastiCache for Redis**（全球）
- **Google Cloud Memorystore**（全球）
- **Azure Cache for Redis**（欧洲数据中心）
- **Upstash**（Serverless Redis，欧洲支持）
- **Redis Cloud**（官方托管，欧洲节点）

**配置示例（生产环境）：**
```yaml
redis:
  enabled: true
  host: your-redis.cache.amazonaws.com
  port: 6379
  password: ${REDIS_PASSWORD}  # 使用环境变量
  db: 0
  pool_size: 200
  min_idle_conns: 50
  max_retries: 5
  default_expiration: "30m"
```

### 方案B：自托管Redis Cluster

```yaml
version: '3.8'

services:
  redis-node-1:
    image: redis:7-alpine
    command: redis-server --cluster-enabled yes --cluster-config-file nodes.conf
    ports:
      - "6379:6379"
    volumes:
      - redis1_data:/data

  redis-node-2:
    image: redis:7-alpine
    command: redis-server --cluster-enabled yes --cluster-config-file nodes.conf
    ports:
      - "6380:6379"
    volumes:
      - redis2_data:/data

  redis-node-3:
    image: redis:7-alpine
    command: redis-server --cluster-enabled yes --cluster-config-file nodes.conf
    ports:
      - "6381:6379"
    volumes:
      - redis3_data:/data
```

### Redis持久化配置

**RDB（快照）+ AOF（日志）双持久化：**

```conf
# redis.conf
save 900 1       # 900秒内至少1次写入
save 300 10      # 300秒内至少10次写入
save 60 10000    # 60秒内至少10000次写入

appendonly yes   # 启用AOF
appendfsync everysec  # 每秒同步一次
```

---

## 🔍 监控与调试

### 查看Redis状态

```bash
# 进入Redis CLI
docker exec -it bookcommunity-redis redis-cli -a dev_redis_2024

# 查看信息
INFO

# 查看内存使用
INFO memory

# 查看连接数
INFO clients

# 查看慢查询
SLOWLOG GET 10

# 查看所有键
KEYS *  # 生产环境禁用！使用SCAN代替

# 扫描键（安全）
SCAN 0 MATCH user:* COUNT 100
```

### 性能监控

```go
// 获取缓存统计
userService := services.GetUserCacheService()
stats := userService.GetCacheStats()

fmt.Printf("Redis enabled: %v\n", stats["redis_enabled"])
fmt.Printf("Memory cache size: %v\n", stats["memory_len"])
fmt.Printf("Redis pool stats: %+v\n", stats["redis_pool"])
```

### 常用命令

```bash
# 获取键值
GET user:1

# 查看过期时间
TTL user:1

# 删除键
DEL user:1

# 清空数据库（危险！）
FLUSHDB

# 查看内存占用最大的键
redis-cli --bigkeys
```

---

## 🐛 常见问题

### Q1: Redis连接失败

**症状：** `connection refused` 或 `NOAUTH Authentication required`

**解决方案：**
```bash
# 检查Redis是否运行
docker ps | grep redis

# 检查密码配置
cat config/conf/config.yaml | grep password

# 测试连接
redis-cli -h localhost -p 6379 -a dev_redis_2024 ping
```

### Q2: 缓存未生效

**原因：** `redis.enabled = false` 或 Redis未启动

**解决方案：**
```yaml
# config/conf/config.yaml
redis:
  enabled: true  # 确保启用
```

### Q3: 内存占用过高

**解决方案：**
```bash
# 设置最大内存限制
redis-cli CONFIG SET maxmemory 2gb
redis-cli CONFIG SET maxmemory-policy allkeys-lru

# 或在docker-compose.yaml中配置
command: redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru
```

### Q4: Redis崩溃重启后数据丢失

**解决方案：** 启用持久化
```yaml
# docker-compose.dev.yaml
command: redis-server --appendonly yes --appendfsync everysec
```

---

## 📚 最佳实践

### ✅ DO

1. **合理设置过期时间**
   ```go
   cache.Set(key, value, 15*time.Minute) // 热数据
   cache.Set(key, value, 1*time.Hour)    // 冷数据
   ```

2. **使用混合缓存**
   ```go
   hybridCache := cache.GetHybridCache()
   hybridCache.Get(key, &value) // 自动降级
   ```

3. **批量操作**
   ```go
   cache.MSet(map[string]interface{}{
       "user:1": user1,
       "user:2": user2,
   })
   ```

4. **主动失效**
   ```go
   db.Save(&user)
   cache.Delete(fmt.Sprintf("user:%d", user.ID))
   ```

### ❌ DON'T

1. **不要缓存敏感数据**
   ```go
   // ❌ 错误：密码不应缓存
   cache.Set("user:1:password", hashedPassword, time.Hour)
   ```

2. **不要使用KEYS命令**
   ```go
   // ❌ 错误：生产环境会阻塞Redis
   keys, _ := client.Keys(ctx, "*").Result()

   // ✅ 正确：使用SCAN
   iter := client.Scan(ctx, 0, "user:*", 100).Iterator()
   ```

3. **不要忘记设置过期时间**
   ```go
   // ❌ 错误：永不过期，内存泄漏
   cache.Set(key, value, 0)

   // ✅ 正确
   cache.Set(key, value, time.Hour)
   ```

---

## 📖 参考资源

- [Redis 官方文档](https://redis.io/docs/)
- [go-redis 文档](https://redis.uptrace.dev/)
- [Redis 最佳实践](https://redis.io/docs/management/optimization/)
- [Redis University](https://university.redis.com/)

---

## 🆘 支持

遇到问题？

1. 查看日志：`docker logs bookcommunity-redis`
2. 检查配置：`config/conf/config.yaml`
3. 测试连接：`redis-cli -h localhost -p 6379 -a dev_redis_2024 ping`
4. 查看统计：`userService.GetCacheStats()`
