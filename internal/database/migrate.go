package database

import (
	"fmt"
	"log"

	"github.com/yourusername/bookcommunity/config"
	"gorm.io/gorm"
)

// MigrationTask 迁移任务接口
type MigrationTask interface {
	Name() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db         *gorm.DB
	migrations []MigrationTask
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: []MigrationTask{},
	}
}

// Register 注册迁移任务
func (m *MigrationManager) Register(task MigrationTask) {
	m.migrations = append(m.migrations, task)
}

// RunUp 执行所有向上迁移
func (m *MigrationManager) RunUp() error {
	for _, migration := range m.migrations {
		log.Printf("Running migration: %s", migration.Name())
		if err := migration.Up(m.db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Name(), err)
		}
		log.Printf("Migration %s completed successfully", migration.Name())
	}
	return nil
}

// RunDown 执行所有向下迁移（回滚）
func (m *MigrationManager) RunDown() error {
	// 逆序执行
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		log.Printf("Rolling back migration: %s", migration.Name())
		if err := migration.Down(m.db); err != nil {
			return fmt.Errorf("rollback %s failed: %w", migration.Name(), err)
		}
		log.Printf("Rollback %s completed successfully", migration.Name())
	}
	return nil
}

// --- PostgreSQL 特定优化迁移示例 ---

// AddIndexesMigration 添加索引优化
type AddIndexesMigration struct{}

func (m *AddIndexesMigration) Name() string {
	return "add_performance_indexes"
}

func (m *AddIndexesMigration) Up(db *gorm.DB) error {
	dbConf := config.GetDatabaseConfig()

	// 仅针对 PostgreSQL 创建特定索引
	if dbConf.Driver == "postgres" || dbConf.Driver == "postgresql" {
		// 为用户表添加索引
		if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users_models(username)").Error; err != nil {
			return err
		}

		// 为视频表添加复合索引
		if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_author_created ON video_models(author_id, created_at DESC)").Error; err != nil {
			return err
		}

		// 为评论表添加索引
		if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_video_created ON comment_models(video_id, created_at DESC)").Error; err != nil {
			return err
		}

		log.Println("PostgreSQL performance indexes created")
	}

	return nil
}

func (m *AddIndexesMigration) Down(db *gorm.DB) error {
	dbConf := config.GetDatabaseConfig()

	if dbConf.Driver == "postgres" || dbConf.Driver == "postgresql" {
		db.Exec("DROP INDEX IF EXISTS idx_users_username")
		db.Exec("DROP INDEX IF EXISTS idx_videos_author_created")
		db.Exec("DROP INDEX IF EXISTS idx_comments_video_created")
		log.Println("PostgreSQL indexes dropped")
	}

	return nil
}

// AddFullTextSearchMigration 添加全文搜索支持（PostgreSQL）
type AddFullTextSearchMigration struct{}

func (m *AddFullTextSearchMigration) Name() string {
	return "add_fulltext_search"
}

func (m *AddFullTextSearchMigration) Up(db *gorm.DB) error {
	dbConf := config.GetDatabaseConfig()

	// PostgreSQL 全文搜索
	if dbConf.Driver == "postgres" || dbConf.Driver == "postgresql" {
		// 为视频标题添加全文搜索索引
		if err := db.Exec(`
			CREATE INDEX IF NOT EXISTS idx_videos_title_fulltext
			ON video_models
			USING gin(to_tsvector('english', title))
		`).Error; err != nil {
			return err
		}

		log.Println("PostgreSQL full-text search indexes created")
	}

	return nil
}

func (m *AddFullTextSearchMigration) Down(db *gorm.DB) error {
	dbConf := config.GetDatabaseConfig()

	if dbConf.Driver == "postgres" || dbConf.Driver == "postgresql" {
		db.Exec("DROP INDEX IF EXISTS idx_videos_title_fulltext")
		log.Println("PostgreSQL full-text search indexes dropped")
	}

	return nil
}

// RunMigrations 执行所有迁移（供外部调用）
func RunMigrations() error {
	manager := NewMigrationManager(GetDB())

	// 注册迁移任务
	manager.Register(&AddIndexesMigration{})
	manager.Register(&AddFullTextSearchMigration{})

	return manager.RunUp()
}
