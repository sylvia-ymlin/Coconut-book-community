package like

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
	"github.com/sylvia-ymlin/Coconut-book-community/pkg/utils"
	"gorm.io/gorm"
)

var logger = utils.NewLogger("like_handler")

// LikeReviewHandler 点赞书评
// @Summary 点赞书评
// @Description 点赞一条书评
// @Tags Like
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Success 200 {object} response.CommonResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews/{id}/like [post]
func LikeReviewHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.CommonResponse{
			StatusCode: response.ErrUserToken,
			StatusMsg:  "用户未登录",
		})
		return
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

	db := database.GetMysqlDB()

	// 检查书评是否存在
	var review models.BookReviewModel
	if err := db.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.CommonResponse{
				StatusCode: response.Failed,
				StatusMsg:  "书评不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, response.CommonResponse{
				StatusCode: response.Failed,
				StatusMsg:  "查询失败",
			})
		}
		return
	}

	// 检查是否已经点赞
	var existingLike models.UserLikeModel
	result := db.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingLike)
	if result.Error == nil {
		// 已经点赞过
		c.JSON(http.StatusOK, response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "已经点赞过了",
		})
		return
	}

	// 创建点赞记录
	like := models.UserLikeModel{
		UserID:   userID.(uint),
		ReviewID: uint(reviewID),
	}

	if err := db.Create(&like).Error; err != nil {
		logger.Printf("Failed to create like: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "点赞失败",
		})
		return
	}

	// 更新书评的点赞数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", reviewID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1"))

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "点赞成功",
	})

	logger.Printf("User %d liked review %d", userID, reviewID)
}

// UnlikeReviewHandler 取消点赞书评
// @Summary 取消点赞书评
// @Description 取消点赞一条书评
// @Tags Like
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Success 200 {object} response.CommonResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews/{id}/like [delete]
func UnlikeReviewHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.CommonResponse{
			StatusCode: response.ErrUserToken,
			StatusMsg:  "用户未登录",
		})
		return
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

	db := database.GetMysqlDB()

	// 删除点赞记录
	result := db.Where("user_id = ? AND review_id = ?", userID, reviewID).
		Delete(&models.UserLikeModel{})

	if result.Error != nil {
		logger.Printf("Failed to delete like: %v", result.Error)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "取消点赞失败",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "未点赞过此书评",
		})
		return
	}

	// 更新书评的点赞数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", reviewID).
		UpdateColumn("like_count", gorm.Expr("like_count - 1"))

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "取消点赞成功",
	})

	logger.Printf("User %d unliked review %d", userID, reviewID)
}

// GetReviewLikesHandler 获取书评的点赞列表
// @Summary 获取书评点赞列表
// @Description 获取点赞某条书评的用户列表
// @Tags Like
// @Accept json
// @Produce json
// @Param id path int true "书评ID"
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20）"
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} response.CommonResponse
// @Router /api/reviews/{id}/likes [get]
func GetReviewLikesHandler(c *gin.Context) {
	// 获取书评ID
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无效的书评ID",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetMysqlDB()

	// 检查书评是否存在
	var review models.BookReviewModel
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "书评不存在",
		})
		return
	}

	// 查询点赞记录
	var likes []models.UserLikeModel
	offset := (page - 1) * pageSize
	if err := db.Where("review_id = ?", reviewID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 获取用户ID列表
	userIDs := make([]uint, 0, len(likes))
	for _, like := range likes {
		userIDs = append(userIDs, like.UserID)
	}

	// 查询用户信息
	var users []models.UserModel
	if len(userIDs) > 0 {
		db.Where("id IN ?", userIDs).Find(&users)
	}

	// 转换为响应格式
	userInfos := make([]*response.UserInfo, 0, len(users))
	for i := range users {
		userInfos = append(userInfos, &response.UserInfo{
			ID:            users[i].ID,
			Username:      users[i].Username,
			FollowerCount: users[i].FollowerCount,
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, response.UserListResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Users: userInfos,
		Total: int64(len(userInfos)),
	})
}
