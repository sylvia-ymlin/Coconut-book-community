package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/bookcommunity/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	rdb               *redis.Client
	ctx               = context.Background()
	defaultExpiration time.Duration
	redisEnabled      bool
)

// InitRedis 初始化Redis连接
func InitRedis() error {
	conf := config.GetRedisConfig()

	if !conf.Enabled {
		logrus.Info("Redis is disabled, using in-memory cache fallback")
		redisEnabled = false
		return nil
	}

	// 创建Redis客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
		MaxRetries:   conf.MaxRetries,
	})

	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		logrus.Errorf("Failed to connect to Redis: %v", err)
		logrus.Warn("Falling back to in-memory cache")
		redisEnabled = false
		return err
	}

	// 解析默认过期时间
	var err error
	defaultExpiration, err = time.ParseDuration(conf.DefaultExpiration)
	if err != nil {
		defaultExpiration = 1 * time.Hour // 默认1小时
		logrus.Warnf("Invalid default_expiration, using 1h: %v", err)
	}

	redisEnabled = true
	logrus.Infof("✅ Redis connected successfully: %s:%d (DB %d)", conf.Host, conf.Port, conf.DB)
	return nil
}

// GetRedisClient 获取Redis客户端（供外部使用）
func GetRedisClient() *redis.Client {
	return rdb
}

// IsRedisEnabled 检查Redis是否启用
func IsRedisEnabled() bool {
	return redisEnabled
}

// --- 通用缓存操作 ---

// Set 设置缓存（自动序列化为JSON）
func Set(key string, value interface{}, expiration time.Duration) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	if expiration == 0 {
		expiration = defaultExpiration
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return rdb.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存（自动反序列化JSON）
func Get(key string, dest interface{}) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func Delete(keys ...string) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.Expire(ctx, key, expiration).Err()
}

// TTL 获取键的剩余生存时间
func TTL(key string) (time.Duration, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.TTL(ctx, key).Result()
}

// --- 批量操作 ---

// MSet 批量设置（不支持过期时间，适合临时数据）
func MSet(pairs map[string]interface{}) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	// 转换为Redis格式
	redisPairs := make([]interface{}, 0, len(pairs)*2)
	for k, v := range pairs {
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal %s: %w", k, err)
		}
		redisPairs = append(redisPairs, k, data)
	}

	return rdb.MSet(ctx, redisPairs...).Err()
}

// MGet 批量获取
func MGet(keys []string) ([]interface{}, error) {
	if !redisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}

	return rdb.MGet(ctx, keys...).Result()
}

// --- 计数器操作 ---

// Incr 自增
func Incr(key string) (int64, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.Decr(ctx, key).Result()
}

// IncrBy 增加指定值
func IncrBy(key string, value int64) (int64, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.IncrBy(ctx, key, value).Result()
}

// --- 集合操作（用于关注、点赞等场景）---

// SAdd 添加成员到集合
func SAdd(key string, members ...interface{}) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.SAdd(ctx, key, members...).Err()
}

// SRem 从集合移除成员
func SRem(key string, members ...interface{}) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.SRem(ctx, key, members...).Err()
}

// SIsMember 检查成员是否在集合中
func SIsMember(key string, member interface{}) (bool, error) {
	if !redisEnabled {
		return false, fmt.Errorf("redis is not enabled")
	}

	return rdb.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合成员数量
func SCard(key string) (int64, error) {
	if !redisEnabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	return rdb.SCard(ctx, key).Result()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	if !redisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}

	return rdb.SMembers(ctx, key).Result()
}

// --- 有序集合操作（用于排行榜等场景）---

// ZAdd 添加成员到有序集合
func ZAdd(key string, members ...redis.Z) error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.ZAdd(ctx, key, members...).Err()
}

// ZRangeByScore 按分数范围获取成员
func ZRangeByScore(key string, min, max string) ([]string, error) {
	if !redisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}

	return rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZRevRange 获取排名前N的成员（降序）
func ZRevRange(key string, start, stop int64) ([]string, error) {
	if !redisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}

	return rdb.ZRevRange(ctx, key, start, stop).Result()
}

// --- 过期键清理（可选，用于监控）---

// FlushDB 清空当前数据库（谨慎使用！）
func FlushDB() error {
	if !redisEnabled {
		return fmt.Errorf("redis is not enabled")
	}

	return rdb.FlushDB(ctx).Err()
}

// Close 关闭Redis连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
