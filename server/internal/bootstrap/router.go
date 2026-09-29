package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/health"
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/auth/handler"
	authService "leslie-blog-server/internal/modules/auth/service"
	categoryHandler "leslie-blog-server/internal/modules/category/handler"
	categoryRepository "leslie-blog-server/internal/modules/category/repository"
	categoryService "leslie-blog-server/internal/modules/category/service"
	homeHandler "leslie-blog-server/internal/modules/home/handler"
	homeService "leslie-blog-server/internal/modules/home/service"
	permissionRepository "leslie-blog-server/internal/modules/permission/repository"
	postHandler "leslie-blog-server/internal/modules/post/handler"
	postRepository "leslie-blog-server/internal/modules/post/repository"
	postService "leslie-blog-server/internal/modules/post/service"
	projectHandler "leslie-blog-server/internal/modules/project/handler"
	projectRepository "leslie-blog-server/internal/modules/project/repository"
	projectService "leslie-blog-server/internal/modules/project/service"
	roleHandler "leslie-blog-server/internal/modules/role/handler"
	roleRepository "leslie-blog-server/internal/modules/role/repository"
	roleService "leslie-blog-server/internal/modules/role/service"
	tagHandler "leslie-blog-server/internal/modules/tag/handler"
	tagRepository "leslie-blog-server/internal/modules/tag/repository"
	tagService "leslie-blog-server/internal/modules/tag/service"
	techstackHandler "leslie-blog-server/internal/modules/techstack/handler"
	techstackRepository "leslie-blog-server/internal/modules/techstack/repository"
	techstackService "leslie-blog-server/internal/modules/techstack/service"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	userRepository "leslie-blog-server/internal/modules/user/repository"
	userService "leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/response"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/cache"
	"leslie-blog-server/internal/pkg/casbin"
	pkgDatabase "leslie-blog-server/internal/pkg/database"
	"leslie-blog-server/internal/pkg/logger"
	"leslie-blog-server/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Server 表示 Leslie Blog API Server。
type Server struct {
	cfg         *config.Config
	engine      *gin.Engine
	router      *router.Router
	db          *gorm.DB
	log         *logger.Logger
	redisClient *redis.Client
}

func NewRouter(
	cfg *config.Config,
	db *gorm.DB,
	log *logger.Logger,
	enforcer *casbin.Enforcer,
	cacheStore cache.Cache,
	redisClient *redis.Client,
) (*Server, error) {
	gin.SetMode(cfg.App.GinMode)
	engine := gin.New()

	// ==================================================
	// 3. 注册全局 Middleware
	// ==================================================
	engine.Use(
		middleware.RequestID(),
		middleware.RequestLogger(log),
		middleware.ErrorHandler(log),
		middleware.Recovery(),
		middleware.Cors(),
	)
	// 注册 404 处理函数。
	engine.NoRoute(func(c *gin.Context) {
		response.Error(
			c,
			http.StatusNotFound,
			appErrors.ErrNotFound,
			"resource not found",
		)
	})
	// 注册 405 处理函数。
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(func(c *gin.Context) {
		response.Error(
			c,
			http.StatusMethodNotAllowed,
			appErrors.ErrInvalidParams,
			"method not allowed",
		)
	})

	// ==================================================
	// 4. 创建模块 Repository
	// ==================================================

	userRepo := userRepository.NewUserRepository(db)
	roleRepo := roleRepository.NewRoleRepository(db)
	permissionRepo := permissionRepository.NewPermissionRepository(db)
	postRepo := postRepository.NewPostRepository(db)
	categoryRepo := categoryRepository.NewCategoryRepository(db)
	tagRepo := tagRepository.NewTagRepository(db)
	postTagRepo := postRepository.NewPostTagRepository(db)
	projectRepo := projectRepository.NewProjectRepository(db)
	techStackRepo := techstackRepository.NewTechStackRepository(db)

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
	tagSvc := tagService.NewTagService(
		tagRepo,
	)
	postSvc := postService.NewPostService(
		postRepo,
		postTagRepo,
		categoryRepo,
		tagSvc,
		pkgDatabase.NewTransactionManager(db),
		cacheStore,
		log,
	)
	authSvc := authService.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.ExpireHours,
	)
	categorySvc := categoryService.NewCategoryService(
		categoryRepo,
		cacheStore,
		log,
	)
	projectSvc := projectService.NewProjectService(
		projectRepo,
		techStackRepo,
		pkgDatabase.NewTransactionManager(db),
		cacheStore,
		log,
	)
	techstackSvc := techstackService.NewTechStackService(
		techStackRepo,
		cacheStore,
		log,
	)
	homeSvc := homeService.NewHomeService(
		postSvc,
		categorySvc,
		projectSvc,
		techstackSvc,
		cacheStore,
		log,
	)

	// ==================================================
	// 6. 创建模块 Handler
	// ==================================================

	userH := userHandler.NewUserHandler(userSvc)
	authH := handler.NewAuthHandler(authSvc, userSvc)
	roleH := roleHandler.NewRoleHandler(roleSvc)
	postH := postHandler.NewPostHandler(postSvc, enforcer)
	categoryH := categoryHandler.NewCategoryHandler(
		categorySvc,
	)
	tagH := tagHandler.NewTagHandler(
		tagSvc,
	)
	projectH := projectHandler.NewProjectHandler(
		projectSvc,
	)

	publicPostHandler := postHandler.NewPublicPostHandler(
		postSvc,
	)

	publicCategoryHandler := categoryHandler.NewPublicCategoryHandler(
		categorySvc,
	)
	projectPublicHandler := projectHandler.NewProjectPublicHandler(
		projectSvc,
	)
	techstackH := techstackHandler.NewTechStackHandler(
		techstackSvc,
	)
	homeH := homeHandler.NewHomeHandler(homeSvc)

	healthHandler := health.NewHandler()

	readyHandler := health.NewReadyHandler(
		health.NewMySQLChecker(db),
		health.NewRedisChecker(redisClient),
	)

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
		publicPostHandler,
		categoryH,
		publicCategoryHandler,
		tagH,
		projectH,
		projectPublicHandler,
		techstackH,
		homeH,
		healthHandler,
		readyHandler,
		jwtMiddleware,
		enforcer,
	)
	return &Server{
		cfg:         cfg,
		engine:      engine,
		router:      r,
		db:          db,
		log:         log,
		redisClient: redisClient,
	}, nil
}

// Run 启动 HTTP Server。
func (s *Server) Run() error {

	// 注册所有路由。
	s.router.Register()

	// // 生成监听地址。
	// //
	// // 例如：
	// //
	// // Host = 0.0.0.0
	// // Port = 8080
	// //
	// // 最终：
	// //
	// // 0.0.0.0:8080
	// addr := s.cfg.Server.Host +
	// 	":" +
	// 	strconv.Itoa(s.cfg.Server.Port)

	// // 启动 Gin HTTP Server。
	// return s.engine.Run(addr)

	server := &http.Server{
		Addr: s.cfg.Server.Host +
			":" +
			strconv.Itoa(s.cfg.Server.Port),
		Handler: s.engine,
	}

	// ------------------------------------------------
	// 7. 启动 HTTP Server
	// ------------------------------------------------

	go func() {

		s.log.Info(
			"http server started",
			slog.String(
				"addr",
				server.Addr,
			),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {

			s.log.Error(
				"http server stopped unexpectedly",
				slog.Any("error", err),
			)
		}

	}()
	// ------------------------------------------------
	// 8. 等待退出信号
	// ------------------------------------------------

	stop := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	s.log.Info(
		"shutdown signal received",
	)

	// ------------------------------------------------
	// 9. 创建 Shutdown Context
	// ------------------------------------------------

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	// ------------------------------------------------
	// 10. 停止 HTTP Server
	// ------------------------------------------------

	if err := server.Shutdown(ctx); err != nil {

		s.log.Error(
			"failed to shutdown http server",
			slog.Any("error", err),
		)
	}

	// ------------------------------------------------
	// 11. 关闭 Redis
	// ------------------------------------------------

	if err := s.redisClient.Close(); err != nil {

		s.log.Error(
			"failed to close redis",
			slog.Any("error", err),
		)
	}

	// ------------------------------------------------
	// 12. 关闭 MySQL
	// ------------------------------------------------

	if err := database.CloseMySQL(s.db); err != nil {

		s.log.Error(
			"failed to close mysql",
			slog.Any("error", err),
		)
	}

	s.log.Info(
		"application stopped",
	)
	return nil
}
