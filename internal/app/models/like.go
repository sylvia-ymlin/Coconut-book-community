package models

import (
	"sync"
	"time"
)

const (
	UserLikeModelTableName       = "user_like_models"
	UserLikeModelTable_UserID    = "user_id"
	UserLikeModelTable_ReviewID  = "review_id" // 改为 review_id
	UserLikeModelTable_CreatedAt = "created_at"
)

// UserLikeModel 点赞模型
// 用户可以点赞书评
type UserLikeModel struct {
	UserID    uint      `gorm:"primarykey;not null"` // 用户ID
	ReviewID  uint      `gorm:"primarykey;not null"` // 书评ID（原 VideoID）
	CreatedAt time.Time `gorm:"not null"`            // 点赞时间
}

// UserLikeCache 存放在sync.Map中, key为videoID
// type UserLikeCache struct {
// 	// 关注列表本来应该在前端缓存，但是青训营的API要求了，所以这里也缓存一份
// 	IsFollowed bool
// 	VideoCache VideoCacheModel
// }

// UserLike_ReviewAndAuthor 缓存结构：书评和作者信息
type UserLike_ReviewAndAuthor struct {
	ReviewID   uint // 书评ID（原 VideoID）
	AuthorID   uint // 作者ID
	// 书评的作者是否被此用户关注
	IsFollowed bool
}

// UserLikeCacheModel 用户点赞缓存模型
type UserLikeCacheModel struct {
	// map[reviewID]UserLike_ReviewAndAuthor
	ReviewIDMap sync.Map // 改为 ReviewIDMap
}
