package review

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
)

// GetDiscoveryFeedHandler 获取发现页书评流（个性化推荐）
// @Summary 发现页书评流
// @Description 获取个性化推荐的书评流（综合热度、时间、用户兴趣）
// @Tags Feed
// @Accept json
// @Produce json
// @Param latest_time query int false "最新一条的时间戳（用于下拉刷新，传0获取最新）"
// @Param page_size query int false "每页数量（默认20）"
// @Success 200 {object} response.FeedResponse
// @Router /api/feed [get]
func GetDiscoveryFeedHandler(c *gin.Context) {
	// 获取当前用户ID（可选，用于个性化推荐）
	currentUserID, _ := c.Get("user_id")
	var userID uint = 0
	if currentUserID != nil {
		userID = currentUserID.(uint)
	}

	// 解析查询参数
	latestTimeStr := c.DefaultQuery("latest_time", "0")
	latestTime, _ := strconv.ParseInt(latestTimeStr, 10, 64)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetMysqlDB()
	query := db.Model(&models.BookReviewModel{})

	// 如果传了 latest_time，获取该时间之前的书评
	if latestTime > 0 {
		t := time.Unix(latestTime, 0)
		query = query.Where("created_at < ?", t)
	}

	// 发现页策略：综合热度排序
	// 热度分 = 点赞数*3 + 评论数*2 + 收藏数*2 + 浏览数/100 - 发布时间衰减
	// 简化版：按热度 + 时间综合排序
	query = query.Order("(like_count * 3 + comment_count * 2 + collect_count * 2) DESC, created_at DESC")

	// 查询书评
	var reviews []models.BookReviewModel
	if err := query.
		Preload("Author").
		Limit(pageSize + 1). // 多查一条判断是否还有更多
		Find(&reviews).Error; err != nil {
		logger.Printf("Failed to query feed: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 判断是否还有更多
	hasMore := false
	if len(reviews) > pageSize {
		hasMore = true
		reviews = reviews[:pageSize]
	}

	// 获取下一页的时间戳
	var nextTime int64 = 0
	if len(reviews) > 0 {
		nextTime = reviews[len(reviews)-1].CreatedAt.Unix()
	}

	// 转换为响应格式
	reviewInfos := make([]*response.ReviewInfo, 0, len(reviews))
	for i := range reviews {
		reviewInfos = append(reviewInfos, response.ConvertReviewToInfo(&reviews[i], userID))
	}

	c.JSON(http.StatusOK, response.FeedResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Reviews:  reviewInfos,
		NextTime: nextTime,
		HasMore:  hasMore,
	})
}

// GetFollowingFeedHandler 获取关注页书评流
// @Summary 关注页书评流
// @Description 获取关注用户的最新书评（按时间倒序）
// @Tags Feed
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param latest_time query int false "最新一条的时间戳（用于下拉刷新，传0获取最新）"
// @Param page_size query int false "每页数量（默认20）"
// @Success 200 {object} response.FeedResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/feed/following [get]
func GetFollowingFeedHandler(c *gin.Context) {
	// 获取当前用户ID（必须登录）
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.CommonResponse{
			StatusCode: response.ErrUserToken,
			StatusMsg:  "请先登录",
		})
		return
	}
	userID := currentUserID.(uint)

	// 解析查询参数
	latestTimeStr := c.DefaultQuery("latest_time", "0")
	latestTime, _ := strconv.ParseInt(latestTimeStr, 10, 64)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetMysqlDB()

	// 查询当前用户关注的用户ID列表
	var user models.UserModel
	if err := db.Preload("Followers").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询关注列表失败",
		})
		return
	}

	// 获取关注的用户ID列表
	followerIDs := make([]uint, 0, len(user.Followers))
	for _, follower := range user.Followers {
		followerIDs = append(followerIDs, follower.ID)
	}

	// 如果没有关注任何人，返回空列表
	if len(followerIDs) == 0 {
		c.JSON(http.StatusOK, response.FeedResponse{
			CommonResponse: response.CommonResponse{
				StatusCode: response.Success,
				StatusMsg:  "暂无内容，去发现页看看吧",
			},
			Reviews:  []*response.ReviewInfo{},
			NextTime: 0,
			HasMore:  false,
		})
		return
	}

	// 查询关注用户的书评
	query := db.Model(&models.BookReviewModel{}).
		Where("author_id IN ?", followerIDs)

	// 如果传了 latest_time，获取该时间之前的书评
	if latestTime > 0 {
		t := time.Unix(latestTime, 0)
		query = query.Where("created_at < ?", t)
	}

	// 按时间倒序
	query = query.Order("created_at DESC")

	// 查询书评
	var reviews []models.BookReviewModel
	if err := query.
		Preload("Author").
		Limit(pageSize + 1).
		Find(&reviews).Error; err != nil {
		logger.Printf("Failed to query following feed: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 判断是否还有更多
	hasMore := false
	if len(reviews) > pageSize {
		hasMore = true
		reviews = reviews[:pageSize]
	}

	// 获取下一页的时间戳
	var nextTime int64 = 0
	if len(reviews) > 0 {
		nextTime = reviews[len(reviews)-1].CreatedAt.Unix()
	}

	// 转换为响应格式
	reviewInfos := make([]*response.ReviewInfo, 0, len(reviews))
	for i := range reviews {
		reviewInfos = append(reviewInfos, response.ConvertReviewToInfo(&reviews[i], userID))
	}

	c.JSON(http.StatusOK, response.FeedResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Reviews:  reviewInfos,
		NextTime: nextTime,
		HasMore:  hasMore,
	})
}
