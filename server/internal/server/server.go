package server

import (
	"context"
	"strconv"

	"leslie-blog-server/internal/bootstrap"
	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/auth/handler"
	authService "leslie-blog-server/internal/modules/auth/service"
	permissionRepository "leslie-blog-server/internal/modules/permission/repository"
	postHandler "leslie-blog-server/internal/modules/post/handler"
	postRepository "leslie-blog-server/internal/modules/post/repository"
	postService "leslie-blog-server/internal/modules/post/service"
	roleHandler "leslie-blog-server/internal/modules/role/handler"
	roleRepository "leslie-blog-server/internal/modules/role/repository"
	roleService "leslie-blog-server/internal/modules/role/service"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	userRepository "leslie-blog-server/internal/modules/user/repository"
	userService "leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server 表示 Leslie Blog API Server。
type Server struct {
	cfg    *config.Config
	engine *gin.Engine
	router *router.Router
	db     *gorm.DB
}

// New 创建 Server。
func New(cfg *config.Config) (*Server, error) {

	// ==================================================
	// 1. 创建 MySQL 连接
	// ==================================================

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.New(
		db,
		"./configs/casbin_model.conf",
	)

	if err != nil {
		return nil, err
	}

	// ==================================================
	// 2. 创建 Gin Engine
	// ==================================================

	engine := gin.New()

	// 初始化数据库。
	if err := bootstrap.SeedDatabase(context.Background(), db, enforcer); err != nil {
		return nil, err
	}

	// ==================================================
	// 3. 注册全局 Middleware
	// ==================================================

	engine.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Cors(),
	)

	// ==================================================
	// 4. 创建模块 Repository
	// ==================================================

	userRepo := userRepository.NewUserRepository(db)
	roleRepo := roleRepository.NewRoleRepository(db)
	permissionRepo := permissionRepository.NewPermissionRepository(db)
	postRepo := postRepository.NewPostRepository(db)

	// ==================================================
	// 5. 创建模块 Service
	// ==================================================

	userSvc := userService.NewUserService(
		userRepo,
		roleRepo,
		enforcer,
	)
	roleSvc := roleService.NewRoleService(
		roleRepo,
		permissionRepo,
		enforcer,
	)
	postSvc := postService.NewPostService(
		postRepo,
	)
	authSvc := authService.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.ExpireHours,
	)

	// ==================================================
	// 6. 创建模块 Handler
	// ==================================================

	userH := userHandler.NewUserHandler(userSvc)
	authH := handler.NewAuthHandler(authSvc, userSvc)
	roleH := roleHandler.NewRoleHandler(roleSvc)
	postH := postHandler.NewPostHandler(postSvc)

	// ==================================================
	// 9. 创建 JWT Middleware
	// ==================================================
	//
	// 注意：
	//
	// JWT Middleware 需要和 JWT Generate 使用相同的 Secret。
	//
	// 登录时：
	//
	// JWT.Generate(..., cfg.JWT.Secret, ...)
	//
	// 请求认证时：
	//
	// JWT(... cfg.JWT.Secret)
	//
	// 两边必须使用同一个 Secret。
	jwtMiddleware := middleware.JWT(
		cfg.JWT.Secret,
	)

	// ==================================================
	// 10. 创建 Router
	// ==================================================

	r := router.New(
		engine,
		userH,
		authH,
		roleH,
		postH,
		jwtMiddleware,
		enforcer,
	)

	// ==================================================
	// 11. 返回 Server
	// ==================================================

	return &Server{
		cfg:    cfg,
		engine: engine,
		router: r,
		db:     db,
	}, nil
}

// Run 启动 HTTP Server。
func (s *Server) Run() error {

	// 注册所有路由。
	s.router.Register()

	// 生成监听地址。
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
	return s.engine.Run(addr)
}
