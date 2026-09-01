// Package repository 数据访问层：GORM 初始化 + 各域仓储
package repository

import (
	"fmt"
	"time"

	"leslie-blog-server/config"
	"leslie-blog-server/internal/model"
	"leslie-blog-server/pkg/logger"
	"leslie-blog-server/pkg/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

// InitDB 初始化 MySQL（GORM），并迁移 schema、注入默认管理员
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	level := gormlogger.Warn
	switch cfg.MySQL.LogLevel {
	case "silent":
		level = gormlogger.Silent
	case "error":
		level = gormlogger.Error
	case "warn":
		level = gormlogger.Warn
	case "info":
		level = gormlogger.Info
	}
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := seedAdmin(db, cfg); err != nil {
		logger.S().Warnf("seed admin: %v", err)
	}

	DB = db
	return db, nil
}

// InitRedis 初始化 Redis 客户端
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	Redis = rdb
	return rdb, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Post{},
		&model.PostTranslation{},
		&model.PostTag{},
		&model.Project{},
		&model.TechStack{},
		&model.Media{},
		&model.Comment{},
		&model.Settings{},
	)
}

// seedAdmin 在无用户时注入默认管理员 admin / admin123（首次启动）
func seedAdmin(db *gorm.DB, cfg *config.Config) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}
	admin := model.User{
		Username: "admin",
		Password: hash,
		Nickname: "Admin",
		Email:    "admin@leslie.blog",
		Role:     "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	logger.L().Info("default admin created", zap.String("username", "admin"))
	return nil
}

// NewContainer 构造仓储容器，集中持有各域 repo
func NewContainer(db *gorm.DB, rdb *redis.Client) *Container {
	return &Container{
		User:     NewUserRepo(db),
		Category: NewCategoryRepo(db),
		Post:     NewPostRepo(db),
		Project:  NewProjectRepo(db),
		TechStack: NewTechStackRepo(db),
		Media:    NewMediaRepo(db),
		Comment:  NewCommentRepo(db),
		Settings: NewSettingsRepo(db),
	}
}
