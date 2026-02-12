package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sirupsen/logrus"
)

// HybridCache 混合缓存（Redis + 内存LRU）
// - 优先使用Redis（分布式缓存）
// - Redis不可用时自动降级到内存LRU
type HybridCache struct {
	memCache *lru.Cache[string, cacheItem]
	mu       sync.RWMutex
}

type cacheItem struct {
	Value      []byte
	ExpireTime time.Time
}

var globalHybridCache *HybridCache

// InitHybridCache 初始化混合缓存
func InitHybridCache(memCacheSize int) error {
	// 初始化Redis
	if err := InitRedis(); err != nil {
		logrus.Warnf("Redis initialization failed: %v, using memory cache only", err)
	}

	// 初始化内存缓存（作为fallback）
	memCache, err := lru.New[string, cacheItem](memCacheSize)
	if err != nil {
		return fmt.Errorf("failed to create memory cache: %w", err)
	}

	globalHybridCache = &HybridCache{
		memCache: memCache,
	}

	logrus.Infof("✅ Hybrid cache initialized (Redis: %v, Memory LRU: %d)", IsRedisEnabled(), memCacheSize)
	return nil
}

// GetHybridCache 获取全局混合缓存实例
func GetHybridCache() *HybridCache {
	return globalHybridCache
}

// Set 设置缓存（自动选择Redis或内存）
func (h *HybridCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// 尝试写入Redis
	if IsRedisEnabled() {
		if err := Set(key, value, expiration); err != nil {
			logrus.Warnf("Redis set failed for key %s: %v, falling back to memory", key, err)
		} else {
			// Redis写入成功，同时写入内存作为热数据缓存
			h.setMemory(key, data, expiration)
			return nil
		}
	}

	// Redis不可用，使用内存缓存
	h.setMemory(key, data, expiration)
	return nil
}

// Get 获取缓存（自动选择Redis或内存）
func (h *HybridCache) Get(key string, dest interface{}) error {
	// 先查内存（热数据，速度最快）
	if data, ok := h.getMemory(key); ok {
		return json.Unmarshal(data, dest)
	}

	// 内存未命中，尝试Redis
	if IsRedisEnabled() {
		if err := Get(key, dest); err == nil {
			// Redis命中，回填内存缓存
			data, _ := json.Marshal(dest)
			h.setMemory(key, data, time.Hour) // 默认1小时
			return nil
		}
	}

	return fmt.Errorf("cache miss: %s", key)
}

// Delete 删除缓存
func (h *HybridCache) Delete(key string) error {
	// 同时删除Redis和内存
	if IsRedisEnabled() {
		Delete(key) // 忽略错误
	}
	h.deleteMemory(key)
	return nil
}

// Exists 检查键是否存在
func (h *HybridCache) Exists(key string) bool {
	// 先查内存
	if _, ok := h.getMemory(key); ok {
		return true
	}

	// 再查Redis
	if IsRedisEnabled() {
		count, err := Exists(key)
		return err == nil && count > 0
	}

	return false
}

// --- 内存缓存操作 ---

func (h *HybridCache) setMemory(key string, data []byte, expiration time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	expireTime := time.Now().Add(expiration)
	h.memCache.Add(key, cacheItem{
		Value:      data,
		ExpireTime: expireTime,
	})
}

func (h *HybridCache) getMemory(key string) ([]byte, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	item, ok := h.memCache.Get(key)
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.ExpireTime) {
		h.memCache.Remove(key)
		return nil, false
	}

	return item.Value, true
}

func (h *HybridCache) deleteMemory(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.memCache.Remove(key)
}

// --- 统计信息 ---

// Stats 获取缓存统计信息
func (h *HybridCache) Stats() map[string]interface{} {
	stats := map[string]interface{}{
		"redis_enabled": IsRedisEnabled(),
		"memory_len":    h.memCache.Len(),
	}

	if IsRedisEnabled() {
		client := GetRedisClient()
		if client != nil {
			poolStats := client.PoolStats()
			stats["redis_pool"] = map[string]interface{}{
				"hits":       poolStats.Hits,
				"misses":     poolStats.Misses,
				"idle_conns": poolStats.IdleConns,
				"total_conns": poolStats.TotalConns,
			}
		}
	}

	return stats
}

// CleanupExpired 清理过期的内存缓存（定期调用）
func (h *HybridCache) CleanupExpired() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// LRU会自动淘汰，这里只是示例
	logrus.Debug("Memory cache cleanup triggered")
}
