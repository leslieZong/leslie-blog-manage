package testutil

import (
	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"log"
	"testing"

	"gorm.io/gorm"
)

func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// 后面这里连接测试数据库。
	// 当前先不要直接复制你项目中的数据库连接代码，
	// 应该复用现有 database/bootstrap 层。
	// 1. 加载配置。
	cfg, err := config.Load("configs/config.test.yaml", ".env.test")
	if err != nil {
		log.Fatalf(
			"load config failed: %v",
			err,
		)
	}

	// 2. 连接 MySQL。
	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf(
			"connect mysql failed: %v",
			err,
		)
	}
	return db
}
