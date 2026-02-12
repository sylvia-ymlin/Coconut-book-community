package recommendation

import (
	"strconv"

	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var recommendService = services.NewRecommendationService()

// GetRecommendationsHandler 获取个性化推荐
// GET /api/v1/recommend?token=xxx&top_k=10
func GetRecommendationsHandler(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID := c.GetUint("userID")

	// 获取top_k参数
	topK := 10
	if k := c.Query("top_k"); k != "" {
		if parsedK, err := strconv.Atoi(k); err == nil && parsedK > 0 {
			topK = parsedK
		}
	}

	// 调用推荐服务
	books, err := recommendService.GetPersonalizedRecommendations(userID, topK)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Get recommendations failed")

		c.JSON(500, gin.H{
			"status_code": 3002,
			"message":     "推荐服务异常",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":     userID,
		"top_k":       topK,
		"result_size": len(books),
	}).Info("Get recommendations success")

	c.JSON(200, gin.H{
		"status_code": 0,
		"books":       books,
		"total":       len(books),
		"message":     "当前为模拟推荐数据，可对接真实推荐系统",
	})
}

// SearchBooksHandler 搜索图书
// GET /api/v1/search?q=计算机&top_k=10
func SearchBooksHandler(c *gin.Context) {
	// 获取查询关键词
	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{
			"status_code": 2003,
			"message":     "缺少查询参数 q",
		})
		return
	}

	// 获取top_k参数
	topK := 10
	if k := c.Query("top_k"); k != "" {
		if parsedK, err := strconv.Atoi(k); err == nil && parsedK > 0 {
			topK = parsedK
		}
	}

	// 调用搜索服务
	books, err := recommendService.SemanticSearch(query, topK)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"query": query,
			"error": err.Error(),
		}).Error("Search books failed")

		c.JSON(500, gin.H{
			"status_code": 3002,
			"message":     "搜索服务异常",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"query":       query,
		"top_k":       topK,
		"result_size": len(books),
	}).Info("Search books success")

	c.JSON(200, gin.H{
		"status_code": 0,
		"books":       books,
		"total":       len(books),
		"message":     "当前为模拟搜索结果，可对接RAG检索系统",
	})
}

// GetBookDetailHandler 获取图书详情（预留）
// GET /api/v1/book/:isbn
func GetBookDetailHandler(c *gin.Context) {
	isbn := c.Param("isbn")
	if isbn == "" {
		c.JSON(400, gin.H{
			"status_code": 2003,
			"message":     "缺少ISBN参数",
		})
		return
	}

	// TODO: 从数据库或外部API获取图书详情
	// 当前返回mock数据

	c.JSON(200, gin.H{
		"status_code": 0,
		"message":     "图书详情功能待实现",
		"book": gin.H{
			"isbn":  isbn,
			"title": "示例图书",
		},
	})
}
