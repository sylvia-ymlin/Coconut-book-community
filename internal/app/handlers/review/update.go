package review

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
)

// UpdateReviewHandler 更新书评
// @Summary 更新书评
// @Description 更新自己发布的书评（只能更新自己的书评）
// @Tags Review
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Param review body response.UpdateReviewRequest true "更新的字段"
// @Success 200 {object} response.ReviewResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 403 {object} response.CommonResponse
// @Failure 404 {object} response.CommonResponse
// @Router /api/reviews/{id} [put]
func UpdateReviewHandler(c *gin.Context) {
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

	// 查询书评
	db := database.GetMysqlDB()
	var review models.BookReviewModel
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "书评不存在",
		})
		return
	}

	// 检查权限（只能更新自己的书评）
	if review.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无权限修改此书评",
		})
		return
	}

	// 解析请求体
	var req response.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Images != nil {
		if bytes, err := json.Marshal(*req.Images); err == nil {
			updates["images"] = string(bytes)
			// 更新封面图
			if len(*req.Images) > 0 {
				updates["cover_url"] = (*req.Images)[0]
			}
		}
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Tags != nil {
		if bytes, err := json.Marshal(*req.Tags); err == nil {
			updates["tags"] = string(bytes)
		}
	}

	// 执行更新
	if len(updates) > 0 {
		if err := db.Model(&review).Updates(updates).Error; err != nil {
			logger.Printf("Failed to update review: %v", err)
			c.JSON(http.StatusInternalServerError, response.CommonResponse{
				StatusCode: response.Failed,
				StatusMsg:  "更新失败",
			})
			return
		}
	}

	// 重新加载书评（包含作者信息）
	db.Preload("Author").First(&review, reviewID)

	// 返回结果
	c.JSON(http.StatusOK, response.ReviewResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "更新成功",
		},
		Review: response.ConvertReviewToInfo(&review, userID.(uint)),
	})

	logger.Printf("User %d updated review %d", userID, reviewID)
}

// DeleteReviewHandler 删除书评
// @Summary 删除书评
// @Description 删除自己发布的书评（只能删除自己的书评）
// @Tags Review
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Success 200 {object} response.CommonResponse
// @Failure 403 {object} response.CommonResponse
// @Failure 404 {object} response.CommonResponse
// @Router /api/reviews/{id} [delete]
func DeleteReviewHandler(c *gin.Context) {
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

	// 查询书评
	db := database.GetMysqlDB()
	var review models.BookReviewModel
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "书评不存在",
		})
		return
	}

	// 检查权限
	if review.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无权限删除此书评",
		})
		return
	}

	// 软删除
	if err := db.Delete(&review).Error; err != nil {
		logger.Printf("Failed to delete review: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "删除成功",
	})

	logger.Printf("User %d deleted review %d", userID, reviewID)

	// TODO: 异步任务
	// - 删除相关的点赞、评论、收藏记录
	// - 更新用户统计
}
