package server

import (
	"strconv"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/user/handler"
	userRepository "leslie-blog-server/internal/modules/user/repository"
	userService "leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server 表示整个 Leslie Blog API Server。
//
// Server 的核心职责不是处理具体业务，
// 而是把整个应用需要的组件组装起来。
type Server struct {
	// cfg 保存整个项目的配置。
	cfg *config.Config

	// engine 是 Gin HTTP Server。
	engine *gin.Engine

	// router 管理所有 HTTP 路由。
	router *router.Router

	// db 是 GORM 数据库连接。
	//
	// 后续其他 Repository 都会依赖它。
	db *gorm.DB
}

// New 创建并初始化 Server。
//
// 这个函数是整个应用的“依赖注入入口”。
//
// 可以理解成：
//
//	config
//	   ↓
//	database
//	   ↓
//	repository
//	   ↓
//	service
//	   ↓
//	handler
//	   ↓
//	router
//
// 最终组装成一个可以运行的 Server。
func New(cfg *config.Config) (*Server, error) {

	// ----------------------------------------
	// 1. 创建 MySQL 数据库连接
	// ----------------------------------------
	db, err := database.NewMySQL(cfg.MySQL)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------
	// 2. 创建 Gin Engine
	// ----------------------------------------
	engine := gin.New()

	// ----------------------------------------
	// 3. 注册全局 Middleware
	// ----------------------------------------
	engine.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	// ----------------------------------------
	// 4. 创建 User Repository
	// ----------------------------------------
	//
	// Repository 依赖数据库。
	//
	// 所以：
	//
	//	db
	//	 ↓
	//	UserRepository
	userRepo := userRepository.NewUserRepository(db)

	// ----------------------------------------
	// 5. 创建 User Service
	// ----------------------------------------
	//
	// Service 依赖 Repository。
	//
	//	userRepo
	//	   ↓
	//	UserService
	userSvc := userService.NewUserService(userRepo)

	// ----------------------------------------
	// 6. 创建 User Handler
	// ----------------------------------------
	//
	// Handler 依赖 Service。
	//
	//	userSvc
	//	   ↓
	//	UserHandler
	userH := handler.NewUserHandler(userSvc)

	// ----------------------------------------
	// 7. 创建 Router
	// ----------------------------------------
	//
	// Router 依赖 User Handler。
	//
	//	userH
	//	  ↓
	//	Router
	r := router.New(
		engine,
		userH,
	)

	// ----------------------------------------
	// 8. 返回完整 Server
	// ----------------------------------------
	return &Server{
		cfg:    cfg,
		engine: engine,
		router: r,
		db:     db,
	}, nil
}

// Run 启动 HTTP Server。
func (s *Server) Run() error {

	// 注册所有 HTTP 路由。
	s.router.Register()

	// 组装 HTTP Server 地址。
	//
	// 例如：
	//
	// Host = 0.0.0.0
	// Port = 8080
	//
	// 最终：
	//
	// 0.0.0.0:8080
	addr := s.cfg.Server.Host +
		":" +
		strconv.Itoa(s.cfg.Server.Port)

	// 启动 Gin HTTP Server。
	//
	// 这个函数会阻塞当前 goroutine，
	// 直到 Server 停止或者发生错误。
	return s.engine.Run(addr)
}
