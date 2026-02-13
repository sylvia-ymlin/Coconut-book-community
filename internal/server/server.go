package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sylvia-ymlin/Coconut-book-community/config"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/comment"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/favorite"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/feed"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/follow"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/publish"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/handlers/recommendation"
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

	router.Static(config.GetVedioConfig().UrlPrefix, config.GetVedioConfig().BasePath)

	baseGroup := router.Group("/douyin")

	// basic api
	baseGroup.GET("/feed", middleware.JWTMiddleWare("/douyin/feed"), feed.FeedVideoListHandler)
	baseGroup.POST("/user/register/", user.UserRegisterHandler)
	baseGroup.POST("/user/login/", middleware.UserLoginHandler)
	baseGroup.GET("/user/", middleware.JWTMiddleWare(), user.GetUserInfoHandler)
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare(), publish.PublishVedioHandler)
	baseGroup.GET("/publish/list/", middleware.JWTMiddleWare(), publish.QueryPublishListHandler)

	//extend 1
	baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare(), favorite.PostFavorHandler)
	baseGroup.GET("/favorite/list/", middleware.JWTMiddleWare(), favorite.QueryFavorVideoListHandler)
	baseGroup.POST("/comment/action/", middleware.JWTMiddleWare(), comment.PostCommentHandler)
	baseGroup.GET("/comment/list/", middleware.JWTMiddleWare(), comment.QueryCommentListHandler)

	//extend 2
	baseGroup.POST("/relation/action/", middleware.JWTMiddleWare(), follow.PostFollowActionHandler)
	baseGroup.GET("/relation/follow/list/", middleware.JWTMiddleWare(), follow.QueryFollowListHandler)
	baseGroup.GET("/relation/follower/list/", middleware.JWTMiddleWare(), follow.QueryFanListHandler)

	// 推荐功能（预留接口，当前返回mock数据）
	baseGroup.GET("/recommend", middleware.JWTMiddleWare(), recommendation.GetRecommendationsHandler)
	baseGroup.GET("/search", recommendation.SearchBooksHandler)
	baseGroup.GET("/book/:isbn", recommendation.GetBookDetailHandler)

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
