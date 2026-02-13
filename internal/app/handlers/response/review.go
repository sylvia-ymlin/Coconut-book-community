package response

import "github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"

// ReviewResponse 书评响应结构（单个书评）
type ReviewResponse struct {
	CommonResponse
	Review *ReviewInfo `json:"review,omitempty"`
}

// ReviewListResponse 书评列表响应
type ReviewListResponse struct {
	CommonResponse
	Reviews  []*ReviewInfo `json:"reviews,omitempty"`
	NextTime int64         `json:"next_time,omitempty"` // 下一页的时间戳（用于分页）
	Total    int64         `json:"total,omitempty"`     // 总数（可选）
}

// ReviewInfo 书评详细信息
type ReviewInfo struct {
	ID           uint        `json:"id"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	BookISBN     string      `json:"book_isbn,omitempty"`
	BookTitle    string      `json:"book_title,omitempty"`
	Images       []string    `json:"images,omitempty"`       // JSON 数组解析后的图片列表
	CoverURL     string      `json:"cover_url,omitempty"`
	Rating       float64     `json:"rating"`
	Tags         []string    `json:"tags,omitempty"`         // JSON 数组解析后的标签列表
	Author       *UserInfo   `json:"author"`                 // 作者信息
	LikeCount    uint        `json:"like_count"`
	CommentCount uint        `json:"comment_count"`
	ViewCount    uint        `json:"view_count"`
	CollectCount uint        `json:"collect_count"`
	CreatedAt    int64       `json:"created_at"`             // Unix 时间戳
	UpdatedAt    int64       `json:"updated_at"`

	// 当前用户与书评的关系（需要传入当前用户ID）
	IsLiked      bool        `json:"is_liked,omitempty"`     // 是否已点赞
	IsCollected  bool        `json:"is_collected,omitempty"` // 是否已收藏
}

// UserInfo 用户简要信息（用于书评作者）
type UserInfo struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	Avatar          string `json:"avatar,omitempty"`
	FollowerCount   uint   `json:"follower_count"`
	IsFollowed      bool   `json:"is_followed,omitempty"` // 当前用户是否关注了该作者
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	CommonResponse
	Users []*UserInfo `json:"users,omitempty"`
	Total int64       `json:"total,omitempty"`
}

// CreateReviewRequest 创建书评请求
type CreateReviewRequest struct {
	Title     string   `json:"title" binding:"required,min=1,max=200"`
	Content   string   `json:"content" binding:"required,min=1,max=5000"`
	BookISBN  string   `json:"book_isbn,omitempty"`
	BookTitle string   `json:"book_title,omitempty"`
	Images    []string `json:"images,omitempty" binding:"max=9"`       // 最多9张图片
	Rating    float64  `json:"rating" binding:"min=0,max=10"`
	Tags      []string `json:"tags,omitempty" binding:"max=10"`        // 最多10个标签
}

// UpdateReviewRequest 更新书评请求
type UpdateReviewRequest struct {
	Title     *string   `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Content   *string   `json:"content,omitempty" binding:"omitempty,min=1,max=5000"`
	Images    *[]string `json:"images,omitempty" binding:"omitempty,max=9"`
	Rating    *float64  `json:"rating,omitempty" binding:"omitempty,min=0,max=10"`
	Tags      *[]string `json:"tags,omitempty" binding:"omitempty,max=10"`
}

// FeedResponse 书评流响应（发现页、关注页）
type FeedResponse struct {
	CommonResponse
	Reviews  []*ReviewInfo `json:"reviews,omitempty"`
	NextTime int64         `json:"next_time,omitempty"` // 下一页的时间戳
	HasMore  bool          `json:"has_more"`            // 是否还有更多
}

// ConvertReviewToInfo 将数据库模型转换为响应结构
func ConvertReviewToInfo(review *models.BookReviewModel, currentUserID uint) *ReviewInfo {
	if review == nil {
		return nil
	}

	info := &ReviewInfo{
		ID:           review.ID,
		Title:        review.Title,
		Content:      review.Content,
		BookISBN:     review.BookISBN,
		BookTitle:    review.BookTitle,
		CoverURL:     review.CoverURL,
		Rating:       review.Rating,
		LikeCount:    review.LikeCount,
		CommentCount: review.CommentCount,
		ViewCount:    review.ViewCount,
		CollectCount: review.CollectCount,
		CreatedAt:    review.CreatedAt.Unix(),
		UpdatedAt:    review.UpdatedAt.Unix(),
	}

	// 解析 Images JSON
	if review.Images != "" && review.Images != "[]" {
		// TODO: 实际项目中应该使用 json.Unmarshal
		// 这里简化处理
		info.Images = []string{} // 临时为空数组
	}

	// 解析 Tags JSON
	if review.Tags != "" && review.Tags != "[]" {
		// TODO: 实际项目中应该使用 json.Unmarshal
		info.Tags = []string{} // 临时为空数组
	}

	// 作者信息
	if review.Author.ID != 0 {
		info.Author = &UserInfo{
			ID:            review.Author.ID,
			Username:      review.Author.Username,
			FollowerCount: review.Author.FollowerCount,
			// TODO: 判断当前用户是否关注了该作者
			IsFollowed: false,
		}
	}

	// TODO: 判断当前用户是否点赞、收藏
	info.IsLiked = false
	info.IsCollected = false

	return info
}
