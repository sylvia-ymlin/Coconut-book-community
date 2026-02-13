package review

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
)

// GetReviewListHandler 获取书评列表
// @Summary 获取书评列表
// @Description 获取书评列表（支持分页、筛选）
// @Tags Review
// @Accept json
// @Produce json
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20，最大100）"
// @Param user_id query int false "筛选指定用户的书评"
// @Param book_isbn query string false "筛选指定图书的书评"
// @Param order_by query string false "排序方式：latest(最新)、popular(最热)、rating(评分)"
// @Success 200 {object} response.ReviewListResponse
// @Router /api/reviews [get]
func GetReviewListHandler(c *gin.Context) {
	// 获取当前用户ID（可选，用于判断点赞/收藏状态）
	currentUserID, _ := c.Get("user_id")
	var userID uint = 0
	if currentUserID != nil {
		userID = currentUserID.(uint)
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filterUserID, _ := strconv.Atoi(c.Query("user_id"))
	filterBookISBN := c.Query("book_isbn")
	orderBy := c.DefaultQuery("order_by", "latest")

	// 构建查询
	db := database.GetMysqlDB()
	query := db.Model(&models.BookReviewModel{})

	// 筛选条件
	if filterUserID > 0 {
		query = query.Where("author_id = ?", filterUserID)
	}
	if filterBookISBN != "" {
		query = query.Where("book_isbn = ?", filterBookISBN)
	}

	// 排序
	switch orderBy {
	case "popular":
		// 按热度排序：综合点赞数、评论数、收藏数、浏览数
		query = query.Order("(like_count * 3 + comment_count * 2 + collect_count * 2 + view_count) DESC")
	case "rating":
		// 按评分排序
		query = query.Order("rating DESC, created_at DESC")
	default:
		// 按最新排序（默认）
		query = query.Order("created_at DESC")
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 分页查询
	var reviews []models.BookReviewModel
	offset := (page - 1) * pageSize
	if err := query.
		Preload("Author").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		logger.Printf("Failed to query reviews: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 转换为响应格式
	reviewInfos := make([]*response.ReviewInfo, 0, len(reviews))
	for i := range reviews {
		reviewInfos = append(reviewInfos, response.ConvertReviewToInfo(&reviews[i], userID))
	}

	// 返回结果
	c.JSON(http.StatusOK, response.ReviewListResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Reviews: reviewInfos,
		Total:   total,
	})
}

// GetReviewDetailHandler 获取书评详情
// @Summary 获取书评详情
// @Description 获取单条书评的详细信息
// @Tags Review
// @Accept json
// @Produce json
// @Param id path int true "书评ID"
// @Success 200 {object} response.ReviewResponse
// @Failure 404 {object} response.CommonResponse
// @Router /api/reviews/{id} [get]
func GetReviewDetailHandler(c *gin.Context) {
	// 获取当前用户ID（可选）
	currentUserID, _ := c.Get("user_id")
	var userID uint = 0
	if currentUserID != nil {
		userID = currentUserID.(uint)
	}

	// 获取书评ID
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无效的书评ID",
		})
		return
	}

	// 查询书评
	db := database.GetMysqlDB()
	var review models.BookReviewModel
	if err := db.Preload("Author").First(&review, reviewID).Error; err != nil {
		logger.Printf("Review not found: %d", reviewID)
		c.JSON(http.StatusNotFound, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "书评不存在",
		})
		return
	}

	// 异步更新浏览次数
	go func() {
		db.Model(&models.BookReviewModel{}).
			Where("id = ?", reviewID).
			UpdateColumn("view_count", db.Raw("view_count + 1"))
	}()

	// 返回结果
	c.JSON(http.StatusOK, response.ReviewResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Review: response.ConvertReviewToInfo(&review, userID),
	})
}
