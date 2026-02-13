package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sylvia-ymlin/Coconut-book-community/config"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/collect"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/comment"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/follow"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/like"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/recommendation"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/review"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/user"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/middleware"
	"github.com/sylvia-ymlin/Coconut-book-community/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DouyinServer struct {
	Router       *gin.Engine
	PanicHandler gin.HandlerFunc
	// Service
	// ...
}

func NewDouyinServer() *DouyinServer {
	router := initDouyinRouter()
	return &DouyinServer{Router: router}
}

func (s *DouyinServer) Run(addr string) error {
	return s.Router.Run(addr)
}

func initPanicLogWriter() io.Writer {
	panicLogPath := filepath.Join(config.GetLogConfig().Path, config.GetLogConfig().PanicLogName)
	return utils.GetNewLazyFileWriter(panicLogPath)
}

// middleware -> handler -> service -> database
//
// middleware -> handler -> service -> message queue -> database
//
// middleware -> handler -> service -> cache -> database
func initDouyinRouter() *gin.Engine {
	router := gin.New()
	writer := io.MultiWriter(initPanicLogWriter(), os.Stdout)
	if config.IsDebug() {
		router.Use(gin.Logger(), gin.RecoveryWithWriter(writer))
	} else {
		router.Use(gin.RecoveryWithWriter(writer))
	}

	// CORS配置 - 支持前后端联调
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// 静态文件服务（图片）
	router.Static("/static/images", "./data/images")

	// API v1 路由组
	apiGroup := router.Group("/api")

	// ====================
	// 用户相关 API
	// ====================
	userGroup := apiGroup.Group("/users")
	{
		// 注册登录（不需要认证）
		apiGroup.POST("/register", user.UserRegisterHandler)
		apiGroup.POST("/login", middleware.UserLoginHandler)

		// 用户信息（需要认证）
		userGroup.GET("/:id", user.GetUserInfoHandler)

		// 用户的收藏列表
		userGroup.GET("/:user_id/collections", collect.GetUserCollectionsHandler)

		// 关注相关
		userGroup.POST("/:id/follow", middleware.JWTMiddleWare(), follow.PostFollowActionHandler)
		userGroup.GET("/:id/followers", follow.QueryFanListHandler)     // 粉丝列表
		userGroup.GET("/:id/following", follow.QueryFollowListHandler)  // 关注列表
	}

	// ====================
	// 书评相关 API（核心功能）
	// ====================
	reviewGroup := apiGroup.Group("/reviews")
	{
		// 书评 CRUD
		reviewGroup.POST("", middleware.JWTMiddleWare(), review.CreateReviewHandler)      // 创建书评
		reviewGroup.GET("", review.GetReviewListHandler)                                   // 查询列表
		reviewGroup.GET("/:id", review.GetReviewDetailHandler)                            // 查询详情
		reviewGroup.PUT("/:id", middleware.JWTMiddleWare(), review.UpdateReviewHandler)   // 更新书评
		reviewGroup.DELETE("/:id", middleware.JWTMiddleWare(), review.DeleteReviewHandler) // 删除书评

		// 点赞相关
		reviewGroup.POST("/:id/like", middleware.JWTMiddleWare(), like.LikeReviewHandler)   // 点赞
		reviewGroup.DELETE("/:id/like", middleware.JWTMiddleWare(), like.UnlikeReviewHandler) // 取消点赞
		reviewGroup.GET("/:id/likes", like.GetReviewLikesHandler)                          // 点赞列表

		// 评论相关
		reviewGroup.POST("/:id/comments", middleware.JWTMiddleWare(), comment.CreateCommentHandler)  // 发布评论
		reviewGroup.GET("/:id/comments", comment.GetCommentListHandler)                              // 评论列表

		// 收藏相关
		reviewGroup.POST("/:id/collect", middleware.JWTMiddleWare(), collect.CollectReviewHandler)   // 收藏
		reviewGroup.DELETE("/:id/collect", middleware.JWTMiddleWare(), collect.UncollectReviewHandler) // 取消收藏
	}

	// ====================
	// 评论（独立资源）
	// ====================
	commentGroup := apiGroup.Group("/comments")
	{
		commentGroup.DELETE("/:id", middleware.JWTMiddleWare(), comment.DeleteCommentHandler) // 删除评论
	}

	// ====================
	// Feed 流
	// ====================
	feedGroup := apiGroup.Group("/feed")
	{
		feedGroup.GET("", review.GetDiscoveryFeedHandler)                                    // 发现页
		feedGroup.GET("/following", middleware.JWTMiddleWare(), review.GetFollowingFeedHandler) // 关注页
	}

	// ====================
	// 图书推荐（代理到 Python 服务）
	// ====================
	bookGroup := apiGroup.Group("/books")
	{
		bookGroup.GET("/search", recommendation.SearchBooksHandler)           // 搜索图书
		bookGroup.GET("/recommendations", recommendation.GetRecommendationsHandler) // 个性化推荐
		bookGroup.GET("/:isbn", recommendation.GetBookDetailHandler)          // 图书详情
		// TODO: Chat with Book
		// bookGroup.POST("/:isbn/chat", recommendation.ChatWithBookHandler)
	}

	// ====================
	// 兼容旧路由（临时保留，逐步废弃）
	// ====================
	// 保留 /douyin 路由组用于向后兼容
	legacyGroup := router.Group("/douyin")
	{
		legacyGroup.POST("/user/register/", user.UserRegisterHandler)
		legacyGroup.POST("/user/login/", middleware.UserLoginHandler)
		legacyGroup.GET("/user/", middleware.JWTMiddleWare(), user.GetUserInfoHandler)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "BookCommunity API",
		})
	})

	// Swagger API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
