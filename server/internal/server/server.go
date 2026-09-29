package server

import (
	"context"
	"log/slog"
	"os"

	"leslie-blog-server/internal/bootstrap"
	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"

	"leslie-blog-server/internal/pkg/cache"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// New 创建 Server。
func New(cfg *config.Config) (*bootstrap.Server, error) {

	// 创建日志记录器。
	log := logger.New(
		cfg.App.Env == "development",
	)
	// ==================================================
	// 1. 创建 MySQL 连接
	// ==================================================

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Error(
			"failed to initialize mysql",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	// 创建 Casbin。
	enforcer, err := casbin.New(
		db,
		"./configs/casbin_model.conf",
	)
	if err != nil {
		log.Error(
			"failed to initialize casbin",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	// 创建 Redis 连接。
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
	)

	// 测试 Redis。
	if err := redisClient.Ping(
		context.Background(),
	).Err(); err != nil {

		log.Error(
			"failed to connect redis",
			slog.Any("error", err),
		)
	}

	// 创建 Redis 缓存。
	cacheStore := cache.NewRedisCache(
		redisClient,
	)

	// 初始化数据库。
	if err := bootstrap.SeedDatabase(context.Background(), db, enforcer); err != nil {
		return nil, err
	}

	// 3. 创建 Router
	// ==================================================
	r, err := bootstrap.NewRouter(
		cfg,
		db,
		log,
		enforcer,
		cacheStore,
		redisClient,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}
