package models

// Book 图书信息（用于推荐系统返回）
type Book struct {
	ISBN      string  `json:"isbn"`       // ISBN号
	Title     string  `json:"title"`      // 书名
	Author    string  `json:"author"`     // 作者
	CoverURL  string  `json:"cover_url"`  // 封面图URL
	Rating    float32 `json:"rating"`     // 评分
	Reason    string  `json:"reason"`     // 推荐理由
	Publisher string  `json:"publisher"`  // 出版社（可选）
	PubDate   string  `json:"pub_date"`   // 出版日期（可选）
	Summary   string  `json:"summary"`    // 简介（可选）
}

// BookSearchRequest 图书搜索请求
type BookSearchRequest struct {
	Query string `json:"query" binding:"required,min=1"` // 搜索关键词
	TopK  int    `json:"top_k"`                          // 返回结果数量
}

// RecommendRequest 推荐请求
type RecommendRequest struct {
	UserID uint `json:"user_id"` // 用户ID
	TopK   int  `json:"top_k"`   // 返回结果数量
}
