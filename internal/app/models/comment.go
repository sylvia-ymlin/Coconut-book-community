package models

import (
	"sync"

	"gorm.io/gorm"
)

const CommentModelTableName = "comment_models"

const (
	CommentModelTable_ReviewID    = "review_id"  // 改为 review_id
	CommentModelTable_UserID      = "user_id"
	CommentModelTable_Content     = "content"
	CommentModelPreload_Commenter = "Commenter"
)

// CommentModel 评论模型
// 用户可以对书评进行评论
type CommentModel struct {
	gorm.Model
	ReviewID uint                 `gorm:"index;not null"` // 书评ID（原 VideoID）
	UserID   uint                 `gorm:"index;not null"` // 评论者ID
	Content  string               `gorm:"type:text;not null"` // 评论内容
	Review   BookReviewModel      `gorm:"foreignKey:ReviewID"` // 关联的书评
	// 使用前需确保里面有数据
	Commenter UserModel `gorm:"foreignKey:UserID"` // 评论者信息
}

type CommentCacheModel struct {
	//map[commentID]CommentModel
	CacheMap map[uint]CommentModel
	MapLock  *sync.RWMutex
}

func NewCommentCacheModel() CommentCacheModel {
	return CommentCacheModel{
		CacheMap: make(map[uint]CommentModel),
		MapLock:  new(sync.RWMutex),
	}
}
