package review

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/response"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/database"
	"github.com/sylvia-ymlin/Coconut-book-community/pkg/utils"
)

var logger = utils.NewLogger("review_handler")

// CreateReviewHandler 创建书评
// @Summary 创建书评
// @Description 用户发布一条图文书评
// @Tags Review
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param review body response.CreateReviewRequest true "书评信息"
// @Success 200 {object} response.ReviewResponse
// @Failure 400 {object} response.CommonResponse
// @Failure 401 {object} response.CommonResponse
// @Router /api/reviews [post]
func CreateReviewHandler(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.CommonResponse{
			StatusCode: response.ErrUserToken,
			StatusMsg:  "用户未登录",
		})
		return
	}

	// 解析请求体
	var req response.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 将图片数组转换为JSON字符串
	imagesJSON := "[]"
	if len(req.Images) > 0 {
		if bytes, err := json.Marshal(req.Images); err == nil {
			imagesJSON = string(bytes)
		}
	}

	// 将标签数组转换为JSON字符串
	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		if bytes, err := json.Marshal(req.Tags); err == nil {
			tagsJSON = string(bytes)
		}
	}

	// 封面图（取第一张图片）
	coverURL := ""
	if len(req.Images) > 0 {
		coverURL = req.Images[0]
	}

	// 创建书评模型
	review := models.BookReviewModel{
		Title:     req.Title,
		Content:   req.Content,
		BookISBN:  req.BookISBN,
		BookTitle: req.BookTitle,
		Images:    imagesJSON,
		CoverURL:  coverURL,
		Rating:    req.Rating,
		Tags:      tagsJSON,
		AuthorID:  userID.(uint),
	}

	// 保存到数据库
	db := database.GetMysqlDB()
	if err := db.Create(&review).Error; err != nil {
		logger.Printf("Failed to create review: %v", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  "创建书评失败",
		})
		return
	}

	// 预加载作者信息
	db.Preload("Author").First(&review, review.ID)

	// 返回结果
	c.JSON(http.StatusOK, response.ReviewResponse{
		CommonResponse: response.CommonResponse{
			StatusCode: response.Success,
			StatusMsg:  "创建成功",
		},
		Review: response.ConvertReviewToInfo(&review, userID.(uint)),
	})

	logger.Printf("User %d created review %d: %s", userID, review.ID, review.Title)

	// TODO: 异步任务
	// - 更新用户统计
	// - 触发推荐系统更新
	// - 通知关注者（如果需要）
}
