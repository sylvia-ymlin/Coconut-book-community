package database

import (
	"fmt"
	"time"

	"github.com/yourusername/bookcommunity/config"
	"github.com/yourusername/bookcommunity/internal/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// GetMysqlDB 兼容旧方法（已弃用）
// Deprecated: Use GetDB instead
func GetMysqlDB() *gorm.DB {
	return db
}

func init() {
	connectDatabase()
}

// connectDatabase 连接数据库（支持 PostgreSQL 和 MySQL）
func connectDatabase() {
	dbConf := config.GetDatabaseConfig()

	var dialector gorm.Dialector

	switch dbConf.Driver {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			dbConf.Host,
			dbConf.Username,
			dbConf.Password,
			dbConf.Dbname,
			dbConf.Port,
			getSSLMode(dbConf.SSLMode),
			getTimezone(dbConf.Timezone),
		)
		dialector = postgres.Open(dsn)
	case "mysql":
		// 保留 MySQL 支持（需要导入 gorm.io/driver/mysql）
		panic("MySQL driver not imported. Please use PostgreSQL or add mysql driver dependency.")
	default:
		panic(fmt.Sprintf("Unsupported database driver: %s. Supported: postgres, mysql", dbConf.Driver))
	}

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel()),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		// 禁用外键约束（根据需要调整）
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	if err != nil {
		panic("连接数据库失败, error:" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接池失败, error:" + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConns)

	// 设置连接最大生命周期
	if dbConf.ConnMaxLifetime != "" {
		lifetime, err := time.ParseDuration(dbConf.ConnMaxLifetime)
		if err == nil {
			sqlDB.SetConnMaxLifetime(lifetime)
		}
	}

	// 自动迁移表结构
	migrateTable()
}

// migrateTable 自动迁移表结构
func migrateTable() {
	// 设置多对多关联表
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FollowersSlice, &models.UserFollowerModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_FansSlice, &models.UserFollowerModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.UserModel{}, models.UserModelTable_CollectionsSlice, &models.UserCollectionModel{})

	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_LikesSlice, &models.UserLikeModel{})
	db.SetupJoinTable(&models.VideoModel{}, models.VideoModelTable_CollectionsSlice, &models.UserCollectionModel{})

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.UserModel{},
		&models.VideoModel{},
		&models.CommentModel{},
		&models.UserFollowerModel{},
		&models.UserLikeModel{},
		&models.UserCollectionModel{},
	)

	if err != nil {
		panic("数据库表迁移失败, error:" + err.Error())
	}
}

// getSSLMode 获取 SSL 模式（PostgreSQL）
func getSSLMode(mode string) string {
	if mode == "" {
		return "disable" // 开发环境默认禁用 SSL
	}
	return mode
}

// getTimezone 获取时区
func getTimezone(tz string) string {
	if tz == "" {
		return "Asia/Shanghai" // 默认时区
	}
	return tz
}

// getLogLevel 根据环境获取日志级别
func getLogLevel() logger.LogLevel {
	if config.IsDebug() {
		return logger.Info
	}
	return logger.Warn
}
