package initiate

import (
	"github.com/yourusername/bookcommunity/config"
	"github.com/yourusername/bookcommunity/internal/database"
	"github.com/yourusername/bookcommunity/internal/msgQueue"
	"github.com/yourusername/bookcommunity/internal/server"
	"github.com/yourusername/bookcommunity/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {
	initGlobalLogger()
	logrus.Info("hello world")
	initMysql()
	initVideoStorageServer()

	// init cache
	var cacheSize = 1000
	database.InitVideoInfoCacher(cacheSize)
	database.InitVideoCommentCacher(cacheSize)
	database.InitUserCacher(cacheSize)
	database.InitUserFavoriteCacher(cacheSize)

	// init message queue
	msgQueue.InitFavoriteMQ()
	msgQueue.InitCommentMQ()
	msgQueue.InitFollowMQ()

	// main logic
	runDouyinServer()
}

func initGlobalLogger() {
	logConfig := config.GetGlobalLoggerConfig()
	err := log.InitGlobalLogger(logConfig)
	if err != nil {
		panic("初始化日志失败, error:" + err.Error())
	}
	if log.PraseLevel(logConfig.LogLevel) < log.PraseLevel("debug") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func initMysql() {
	database.GetMysqlDB()
}

func initVideoStorageServer() {
	database.GetVideoSaver()
}

func runDouyinServer() {
	douyinServer := server.NewDouyinServer()
	err := douyinServer.Run(":" + config.GetServerPort())
	if err != nil {
		panic("启动服务失败, error:" + err.Error())
	}
}
