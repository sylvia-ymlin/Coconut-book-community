package models

import (
	"time"
)

const UserCollectionModelTableName = "user_collection_models"

// UserCollectionModel 收藏模型
// 用户可以收藏书评，方便后续查看
type UserCollectionModel struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"index;not null"`    // 用户ID
	ReviewID  uint      `gorm:"index;not null"`    // 书评ID（原 VideoID）
	CreatedAt time.Time `gorm:"not null"`          // 收藏时间
}

func (u *UserCollectionModel) TableName() string {
	return UserCollectionModelTableName
}
