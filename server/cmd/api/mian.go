package main

import (
	"log"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/server"
)

func main() {

	// --------------------------------------------------
	// 第一步：加载配置
	// --------------------------------------------------

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf(
			"load config failed: %v",
			err,
		)
	}

	// --------------------------------------------------
	// 第二步：创建 Server
	// --------------------------------------------------
	//
	// Server 创建过程中会：
	//
	//   1. 创建 MySQL
	//   2. 创建 Gin
	//   3. 创建 Middleware
	//   4. 创建 Router
	//
	srv, err := server.New(cfg)

	if err != nil {
		log.Fatalf(
			"create server failed: %v",
			err,
		)
	}

	// --------------------------------------------------
	// 第三步：启动 HTTP Server
	// --------------------------------------------------

	if err := srv.Run(); err != nil {
		log.Fatalf(
			"server run failed: %v",
			err,
		)
	}
}
