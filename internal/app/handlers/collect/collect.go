package collect

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

var logger = utils.NewLogger("collect_handler")

// CollectReviewHandler 收藏书评
// @Summary 收藏书评
// @Description 收藏一条书评
// @Tags Collect
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Success 200 {object} response.CommonResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews/{id}/collect [post]
func CollectReviewHandler(c *gin.Context) {
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

	// 检查是否已经收藏
	var existingCollection models.UserCollectionModel
	result := db.Where("user_id = ? AND review_id = ?", userID, reviewID).First(&existingCollection)
	if result.Error == nil {
		c.JSON(http.StatusOK, response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "已经收藏过了",
		})
		return
	}

	// 创建收藏记录
	collection := models.UserCollectionModel{
		UserID:   userID.(uint),
		ReviewID: uint(reviewID),
	}

	if err := db.Create(&collection).Error; err != nil {
		logger.Printf("Failed to create collection: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "收藏失败",
		})
		return
	}

	// 更新书评的收藏数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", reviewID).
		UpdateColumn("collect_count", gorm.Expr("collect_count + 1"))

	// 更新用户的收藏数统计
	db.Model(&models.UserModel{}).
		Where("id = ?", userID).
		UpdateColumn("collections_count", gorm.Expr("collections_count + 1"))

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "收藏成功",
	})

	logger.Printf("User %d collected review %d", userID, reviewID)
}

// UncollectReviewHandler 取消收藏书评
// @Summary 取消收藏书评
// @Description 取消收藏一条书评
// @Tags Collect
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Success 200 {object} response.CommonResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews/{id}/collect [delete]
func UncollectReviewHandler(c *gin.Context) {
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

	// 删除收藏记录
	result := db.Where("user_id = ? AND review_id = ?", userID, reviewID).
		Delete(&models.UserCollectionModel{})

	if result.Error != nil {
		logger.Printf("Failed to delete collection: %v", result.Error)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "取消收藏失败",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "未收藏过此书评",
		})
		return
	}

	// 更新书评的收藏数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", reviewID).
		UpdateColumn("collect_count", gorm.Expr("collect_count - 1"))

	// 更新用户的收藏数统计
	db.Model(&models.UserModel{}).
		Where("id = ?", userID).
		UpdateColumn("collections_count", gorm.Expr("collections_count - 1"))

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "取消收藏成功",
	})

	logger.Printf("User %d uncollected review %d", userID, reviewID)
}

// GetUserCollectionsHandler 获取用户的收藏列表
// @Summary 获取用户收藏列表
// @Description 获取用户收藏的书评列表
// @Tags Collect
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20）"
// @Success 200 {object} response.ReviewListResponse
// @Failure 400 {object} response.CommonResponse
// @Router /api/users/{user_id}/collections [get]
func GetUserCollectionsHandler(c *gin.Context) {
	// 获取用户ID
	targetUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无效的用户ID",
		})
		return
	}

	// 获取当前用户ID（可选，用于判断点赞/收藏状态）
	currentUserID, _ := c.Get("user_id")
	var userID uint = 0
	if currentUserID != nil {
		userID = currentUserID.(uint)
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

	// 查询收藏记录
	var collections []models.UserCollectionModel
	offset := (page - 1) * pageSize
	if err := db.Where("user_id = ?", targetUserID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 获取书评ID列表
	reviewIDs := make([]uint, 0, len(collections))
	for _, collection := range collections {
		reviewIDs = append(reviewIDs, collection.ReviewID)
	}

	// 查询书评信息
	var reviews []models.BookReviewModel
	if len(reviewIDs) > 0 {
		db.Where("id IN ?", reviewIDs).Preload("Author").Find(&reviews)
	}

	// 转换为响应格式
	reviewInfos := make([]*response.ReviewInfo, 0, len(reviews))
	for i := range reviews {
		reviewInfos = append(reviewInfos, response.ConvertReviewToInfo(&reviews[i], userID))
	}

	// 查询总数
	var total int64
	db.Model(&models.UserCollectionModel{}).Where("user_id = ?", targetUserID).Count(&total)

	c.JSON(http.StatusOK, response.ReviewListResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Reviews: reviewInfos,
		Total:   total,
	})
}
