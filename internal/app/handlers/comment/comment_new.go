package comment

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

var logger = utils.NewLogger("comment_handler")

// CommentRequest 评论请求
type CommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// CommentResponse 单条评论响应
type CommentResponse struct {
	response.CommonResponse
	Comment *CommentInfo `json:"comment,omitempty"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	response.CommonResponse
	Comments []*CommentInfo `json:"comments,omitempty"`
	Total    int64          `json:"total,omitempty"`
}

// CommentInfo 评论详细信息
type CommentInfo struct {
	ID         uint                `json:"id"`
	Content    string              `json:"content"`
	User       *response.UserInfo  `json:"user"`      // 评论者信息
	CreatedAt  int64               `json:"created_at"`
}

// CreateCommentHandler 发布评论
// @Summary 发布评论
// @Description 对一条书评发布评论
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "书评ID"
// @Param comment body CommentRequest true "评论内容"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews/{id}/comments [post]
func CreateCommentHandler(c *gin.Context) {
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

	// 解析请求体
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "请求参数错误: " + err.Error(),
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

	// 创建评论
	comment := models.CommentModel{
		ReviewID: uint(reviewID),
		UserID:   userID.(uint),
		Content:  req.Content,
	}

	if err := db.Create(&comment).Error; err != nil {
		logger.Printf("Failed to create comment: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "评论失败",
		})
		return
	}

	// 更新书评的评论数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", reviewID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	// 更新用户的评论数统计
	db.Model(&models.UserModel{}).
		Where("id = ?", userID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	// 加载评论者信息
	db.Preload("Commenter").First(&comment, comment.ID)

	// 返回结果
	commentInfo := &CommentInfo{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Unix(),
	}

	if comment.Commenter.ID != 0 {
		commentInfo.User = &response.UserInfo{
			ID:            comment.Commenter.ID,
			Username:      comment.Commenter.Username,
			FollowerCount: comment.Commenter.FollowerCount,
		}
	}

	c.JSON(http.StatusOK, CommentResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "评论成功",
		},
		Comment: commentInfo,
	})

	logger.Printf("User %d commented on review %d", userID, reviewID)
}

// GetCommentListHandler 获取书评的评论列表
// @Summary 获取评论列表
// @Description 获取某条书评的所有评论
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "书评ID"
// @Param page query int false "页码（默认1）"
// @Param page_size query int false "每页数量（默认20）"
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} response.CommonResponse
// @Router /api/reviews/{id}/comments [get]
func GetCommentListHandler(c *gin.Context) {
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

	// 查询评论
	var comments []models.CommentModel
	offset := (page - 1) * pageSize
	if err := db.Where("review_id = ?", reviewID).
		Preload("Commenter").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "查询失败",
		})
		return
	}

	// 转换为响应格式
	commentInfos := make([]*CommentInfo, 0, len(comments))
	for i := range comments {
		info := &CommentInfo{
			ID:        comments[i].ID,
			Content:   comments[i].Content,
			CreatedAt: comments[i].CreatedAt.Unix(),
		}

		if comments[i].Commenter.ID != 0 {
			info.User = &response.UserInfo{
				ID:            comments[i].Commenter.ID,
				Username:      comments[i].Commenter.Username,
				FollowerCount: comments[i].Commenter.FollowerCount,
			}
		}

		commentInfos = append(commentInfos, info)
	}

	// 查询总数
	var total int64
	db.Model(&models.CommentModel{}).Where("review_id = ?", reviewID).Count(&total)

	c.JSON(http.StatusOK, CommentListResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "查询成功",
		},
		Comments: commentInfos,
		Total:    total,
	})
}

// DeleteCommentHandler 删除评论
// @Summary 删除评论
// @Description 删除自己发布的评论（只能删除自己的评论）
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "评论ID"
// @Success 200 {object} response.CommonResponse
// @Failure 403 {object} response.CommonResponse
// @Failure 404 {object} response.CommonResponse
// @Router /api/comments/{id} [delete]
func DeleteCommentHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.CommonResponse{
			StatusCode: response.ErrUserToken,
			StatusMsg:  "用户未登录",
		})
		return
	}

	// 获取评论ID
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无效的评论ID",
		})
		return
	}

	db := database.GetMysqlDB()

	// 查询评论
	var comment models.CommentModel
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "评论不存在",
		})
		return
	}

	// 检查权限（只能删除自己的评论）
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无权限删除此评论",
		})
		return
	}

	// 删除评论
	if err := db.Delete(&comment).Error; err != nil {
		logger.Printf("Failed to delete comment: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "删除失败",
		})
		return
	}

	// 更新书评的评论数
	db.Model(&models.BookReviewModel{}).
		Where("id = ?", comment.ReviewID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - 1"))

	// 更新用户的评论数统计
	db.Model(&models.UserModel{}).
		Where("id = ?", userID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - 1"))

	c.JSON(http.StatusOK, response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "删除成功",
	})

	logger.Printf("User %d deleted comment %d", userID, commentID)
}
