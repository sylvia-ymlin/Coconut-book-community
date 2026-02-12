package services

import (
	"fmt"
	"time"

	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/cache"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
	"github.com/sirupsen/logrus"
)

// UserCacheService 用户缓存服务（展示Redis使用示例）
type UserCacheService struct {
	cache *cache.HybridCache
}

var userCacheService *UserCacheService

// InitUserCacheService 初始化用户缓存服务
func InitUserCacheService() {
	// 初始化混合缓存（5000个条目的内存LRU + Redis）
	if err := cache.InitHybridCache(5000); err != nil {
		logrus.Fatalf("Failed to init hybrid cache: %v", err)
	}

	userCacheService = &UserCacheService{
		cache: cache.GetHybridCache(),
	}

	logrus.Info("✅ UserCacheService initialized")
}

// GetUserCacheService 获取用户缓存服务实例
func GetUserCacheService() *UserCacheService {
	if userCacheService == nil {
		InitUserCacheService()
	}
	return userCacheService
}

// GetUserByID 获取用户（带缓存）
func (s *UserCacheService) GetUserByID(userID uint) (*models.UserCacheModel, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	// 1. 尝试从缓存获取
	var cachedUser models.UserCacheModel
	if err := s.cache.Get(cacheKey, &cachedUser); err == nil {
		logrus.Debugf("Cache hit for user %d", userID)
		return &cachedUser, nil
	}

	// 2. 缓存未命中，从数据库查询
	logrus.Debugf("Cache miss for user %d, querying database", userID)
	var user models.UserModel
	db := database.GetDB()

	if err := db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 3. 转换为缓存模型
	var userCache models.UserCacheModel
	userCache.SetValue(user)

	// 4. 写入缓存（15分钟过期）
	if err := s.cache.Set(cacheKey, userCache, 15*time.Minute); err != nil {
		logrus.Warnf("Failed to cache user %d: %v", userID, err)
	}

	return &userCache, nil
}

// GetUserByUsername 通过用户名获取用户（带缓存）
func (s *UserCacheService) GetUserByUsername(username string) (*models.UserCacheModel, error) {
	cacheKey := fmt.Sprintf("user:name:%s", username)

	// 1. 尝试从缓存获取
	var cachedUser models.UserCacheModel
	if err := s.cache.Get(cacheKey, &cachedUser); err == nil {
		logrus.Debugf("Cache hit for username %s", username)
		return &cachedUser, nil
	}

	// 2. 缓存未命中，从数据库查询
	logrus.Debugf("Cache miss for username %s, querying database", username)
	var user models.UserModel
	db := database.GetDB()

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 3. 转换为缓存模型
	var userCache models.UserCacheModel
	userCache.SetValue(user)

	// 4. 写入缓存（10分钟过期）
	if err := s.cache.Set(cacheKey, userCache, 10*time.Minute); err != nil {
		logrus.Warnf("Failed to cache username %s: %v", username, err)
	}

	return &userCache, nil
}

// InvalidateUserCache 使用户缓存失效（用户更新时调用）
func (s *UserCacheService) InvalidateUserCache(userID uint, username string) error {
	cacheKeyID := fmt.Sprintf("user:%d", userID)
	cacheKeyName := fmt.Sprintf("user:name:%s", username)

	if err := s.cache.Delete(cacheKeyID); err != nil {
		logrus.Warnf("Failed to delete cache for user %d: %v", userID, err)
	}

	if err := s.cache.Delete(cacheKeyName); err != nil {
		logrus.Warnf("Failed to delete cache for username %s: %v", username, err)
	}

	logrus.Infof("Invalidated cache for user %d (%s)", userID, username)
	return nil
}

// GetCacheStats 获取缓存统计信息
func (s *UserCacheService) GetCacheStats() map[string]interface{} {
	return s.cache.Stats()
}

// --- 计数器示例（点赞、关注等）---

// IncrementFollowerCount 增加关注者计数（使用Redis计数器）
func (s *UserCacheService) IncrementFollowerCount(userID uint) (int64, error) {
	if !cache.IsRedisEnabled() {
		return 0, fmt.Errorf("redis is not enabled")
	}

	cacheKey := fmt.Sprintf("user:%d:follower_count", userID)
	count, err := cache.Incr(cacheKey)
	if err != nil {
		return 0, err
	}

	// 设置30分钟过期（定期同步到数据库）
	cache.Expire(cacheKey, 30*time.Minute)
	return count, nil
}

// GetFollowerCount 获取关注者计数
func (s *UserCacheService) GetFollowerCount(userID uint) (int64, error) {
	if !cache.IsRedisEnabled() {
		// 从数据库获取
		var user models.UserModel
		db := database.GetDB()
		if err := db.First(&user, userID).Error; err != nil {
			return 0, err
		}
		return int64(user.FollowerCount), nil
	}

	cacheKey := fmt.Sprintf("user:%d:follower_count", userID)

	// 检查Redis中是否存在
	exists, err := cache.Exists(cacheKey)
	if err != nil || exists == 0 {
		// 从数据库加载
		var user models.UserModel
		db := database.GetDB()
		if err := db.First(&user, userID).Error; err != nil {
			return 0, err
		}

		// 写入Redis
		cache.Set(cacheKey, user.FollowerCount, 30*time.Minute)
		return int64(user.FollowerCount), nil
	}

	// 从Redis获取
	var count int64
	if err := cache.Get(cacheKey, &count); err != nil {
		return 0, err
	}

	return count, nil
}

// --- 集合操作示例（关注列表）---

// AddFollowing 添加关注（使用Redis Set）
func (s *UserCacheService) AddFollowing(userID uint, followingID uint) error {
	if !cache.IsRedisEnabled() {
		return fmt.Errorf("redis is not enabled")
	}

	cacheKey := fmt.Sprintf("user:%d:following", userID)
	return cache.SAdd(cacheKey, followingID)
}

// RemoveFollowing 取消关注
func (s *UserCacheService) RemoveFollowing(userID uint, followingID uint) error {
	if !cache.IsRedisEnabled() {
		return fmt.Errorf("redis is not enabled")
	}

	cacheKey := fmt.Sprintf("user:%d:following", userID)
	return cache.SRem(cacheKey, followingID)
}

// IsFollowing 检查是否已关注
func (s *UserCacheService) IsFollowing(userID uint, targetID uint) (bool, error) {
	if !cache.IsRedisEnabled() {
		// 从数据库查询
		db := database.GetDB()
		var count int64
		err := db.Model(&models.UserFollowerModel{}).
			Where("user_id = ? AND follower_id = ?", userID, targetID).
			Count(&count).Error
		return count > 0, err
	}

	cacheKey := fmt.Sprintf("user:%d:following", userID)
	return cache.SIsMember(cacheKey, targetID)
}

// GetFollowingCount 获取关注数量
func (s *UserCacheService) GetFollowingCount(userID uint) (int64, error) {
	if !cache.IsRedisEnabled() {
		db := database.GetDB()
		var user models.UserModel
		if err := db.First(&user, userID).Error; err != nil {
			return 0, err
		}
		return int64(user.FollowerCount), nil
	}

	cacheKey := fmt.Sprintf("user:%d:following", userID)
	return cache.SCard(cacheKey)
}
