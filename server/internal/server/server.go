package server

import (
	"strconv"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/middleware"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	userRepository "leslie-blog-server/internal/modules/user/repository"
	userService "leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server 表示整个 HTTP 应用。
type Server struct {

	// 应用配置。
	cfg *config.Config

	// Gin HTTP Engine。
	engine *gin.Engine

	// 路由管理器。
	router *router.Router

	// GORM 数据库连接。
	db *gorm.DB
}

// New 创建 Server。
func New(cfg *config.Config) (*Server, error) {

	// --------------------------------------------------
	// 1. 创建 MySQL
	// --------------------------------------------------

	db, err := database.NewMySQL(cfg.MySQL)

	if err != nil {
		return nil, err
	}

	// --------------------------------------------------
	// 2. 创建 Gin
	// --------------------------------------------------

	engine := gin.New()

	// --------------------------------------------------
	// 3. 注册全局 Middleware
	// --------------------------------------------------

	engine.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	// --------------------------------------------------
	// 4. 创建基础 Router
	// --------------------------------------------------

	// 创建 User Repository。
	userRepo := userRepository.NewUserRepository(db)

	// 创建 User Service。
	userSvc := userService.NewUserService(userRepo)

	// 创建 User Handler。
	userH := userHandler.NewUserHandler(userSvc)

	// 创建 Router。
	r := router.New(
		engine,
		userH,
	)

	// --------------------------------------------------
	// 9. 返回 Server
	// --------------------------------------------------

	return &Server{
		cfg:    cfg,
		engine: engine,
		router: r,
		db:     db,
	}, nil
}

// Run 启动 HTTP Server。
func (s *Server) Run() error {

	// 注册路由。
	s.router.Register()

	addr := s.cfg.Server.Host +
		":" +
		strconv.Itoa(s.cfg.Server.Port)

	return s.engine.Run(addr)
}
