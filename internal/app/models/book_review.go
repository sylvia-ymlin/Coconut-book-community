package models

import "gorm.io/gorm"

const BookReviewTitleMaxByteLength = 199
const BookReviewTitleMaxRuneLength = 60
const BookReviewContentMaxLength = 5000

const (
	BookReviewModelTableName              = "book_reviews"
	BookReviewModelTable_LikeCount        = "like_count"
	BookReviewModelTable_CreatedAt        = "created_at"
	BookReviewModelTable_AuthorID         = "author_id"
	BookReviewModelTable_BookISBN         = "book_isbn"
	BookReviewModelTable_LikesSlice       = "Likes"
	BookReviewModelTable_CollectionsSlice = "Collections"
)

// BookReviewModel 图书书评模型
// 用户可以发布关于某本书的图文书评（类似小红书）
type BookReviewModel struct {
	gorm.Model

	// 基本信息
	Title   string `gorm:"size:200;not null"`              // 书评标题
	Content string `gorm:"type:text;not null"`             // 书评内容（必填，主要内容）

	// 关联图书信息
	BookISBN  string `gorm:"size:20;index"`                // 关联图书 ISBN（可选，用于推荐）
	BookTitle string `gorm:"size:200"`                     // 图书标题（冗余存储，提升查询性能）

	// 媒体资源（图片）
	Images   string `gorm:"type:text"`                     // 图片URL列表（JSON 数组，最多9张）
	CoverURL string `gorm:"size:200"`                      // 封面图（第一张图片）

	// 书评属性
	Rating     float64 `gorm:"type:decimal(3,1);default:0"` // 用户评分 (0.0-10.0)

	// 统计信息
	LikeCount    uint `gorm:"default:0"`                   // 点赞数
	CommentCount uint `gorm:"default:0"`                   // 评论数
	ViewCount    uint `gorm:"default:0"`                   // 浏览次数
	CollectCount uint `gorm:"default:0"`                   // 收藏次数

	// 作者信息
	Author   UserModel `gorm:"foreignKey:AuthorID"`
	AuthorID uint      `gorm:"index;not null"`

	// 社交关系
	Comments    []CommentModel `gorm:"foreignKey:ReviewID"`
	Likes       []UserModel    `gorm:"many2many:user_like;joinForeignKey:ReviewID;joinReferences:UserID"`
	Collections []UserModel    `gorm:"many2many:user_collection;joinForeignKey:ReviewID;joinReferences:UserID"`

	// 标签（可选，用于分类和搜索）
	Tags string `gorm:"size:500"` // 标签列表（JSON 数组，如 ["悬疑", "推理", "经典"]）
}

func (b *BookReviewModel) TableName() string {
	return BookReviewModelTableName
}

// BookReviewCacheModel 书评缓存模型（用于 Redis）
type BookReviewCacheModel struct {
	gorm.Model
	Title        string
	Content      string
	BookISBN     string
	BookTitle    string
	Images       string
	CoverURL     string
	Rating       float64
	AuthorID     uint
	LikeCount    uint
	CommentCount uint
	ViewCount    uint
	CollectCount uint
	Tags         string
}

func (b *BookReviewCacheModel) SetValue(other BookReviewModel) {
	b.ID = other.ID
	b.CreatedAt = other.CreatedAt
	b.UpdatedAt = other.UpdatedAt
	b.DeletedAt = other.DeletedAt
	b.Title = other.Title
	b.Content = other.Content
	b.BookISBN = other.BookISBN
	b.BookTitle = other.BookTitle
	b.Images = other.Images
	b.CoverURL = other.CoverURL
	b.Rating = other.Rating
	b.AuthorID = other.AuthorID
	b.LikeCount = other.LikeCount
	b.CommentCount = other.CommentCount
	b.ViewCount = other.ViewCount
	b.CollectCount = other.CollectCount
	b.Tags = other.Tags
}

func (b *BookReviewModel) SetValueFromCacheModel(other BookReviewCacheModel) {
	b.ID = other.ID
	b.CreatedAt = other.CreatedAt
	b.UpdatedAt = other.UpdatedAt
	b.DeletedAt = other.DeletedAt
	b.Title = other.Title
	b.Content = other.Content
	b.BookISBN = other.BookISBN
	b.BookTitle = other.BookTitle
	b.Images = other.Images
	b.CoverURL = other.CoverURL
	b.Rating = other.Rating
	b.AuthorID = other.AuthorID
	b.LikeCount = other.LikeCount
	b.CommentCount = other.CommentCount
	b.ViewCount = other.ViewCount
	b.CollectCount = other.CollectCount
	b.Tags = other.Tags
}

// HasImages 是否有配图
func (b *BookReviewModel) HasImages() bool {
	return b.Images != "" && b.Images != "[]"
}

// IncrementViewCount 增加浏览次数
func (b *BookReviewModel) IncrementViewCount(db *gorm.DB) error {
	return db.Model(b).UpdateColumn(BookReviewModelTable_LikeCount, gorm.Expr("view_count + ?", 1)).Error
}
