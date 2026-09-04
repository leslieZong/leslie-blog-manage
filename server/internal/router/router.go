package router

import (
	"leslie-blog-server/internal/modules/auth"
	authHandler "leslie-blog-server/internal/modules/auth/handler"
	"leslie-blog-server/internal/modules/user"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine

	userHandler *userHandler.UserHandler

	authHandler *authHandler.AuthHandler

	// jwtMiddleware 是 JWT 认证中间件。
	//
	// Server 创建好 Middleware 后，
	// 注入到 Router。
	jwtMiddleware gin.HandlerFunc
	enforcer      *casbin.Enforcer
}

// New 创建 Router。
func New(
	engine *gin.Engine,
	userHandler *userHandler.UserHandler,
	authHandler *authHandler.AuthHandler,
	jwtMiddleware gin.HandlerFunc,
	enforcer *casbin.Enforcer,
) *Router {

	return &Router{
		engine:        engine,
		userHandler:   userHandler,
		authHandler:   authHandler,
		jwtMiddleware: jwtMiddleware,
		enforcer:      enforcer,
	}
}

// Register 注册整个项目的路由。
func (r *Router) Register() {

	// ==================================================
	// 健康检查
	// ==================================================

	r.engine.GET("/health", health)

	// ==================================================
	// Public API
	// ==================================================

	v1 := r.engine.Group("/api/v1")

	// 当前还没有公共 API。
	//
	// 所以这里先保留。
	_ = v1

	// ==================================================
	// Admin API
	// ==================================================

	admin := r.engine.Group("/api/admin/v1")

	// --------------------------------------------------
	// 登录接口
	// --------------------------------------------------
	//
	// 登录之前没有 JWT，
	// 所以这里不能使用 JWT Middleware。
	auth.RegisterRoutes(
		admin,
		r.authHandler,
		r.jwtMiddleware,
	)

	// --------------------------------------------------
	// 需要登录的接口
	// --------------------------------------------------
	//
	// 创建一个新的 Router Group。
	protected := admin.Group("")

	// 给 protected Group 添加 JWT Middleware。
	protected.Use(r.jwtMiddleware)

	// 所有注册到 protected 的接口，
	// 都必须先通过 JWT 验证。
	user.RegisterRoutes(
		protected,
		r.userHandler,
		r.enforcer,
	)
}

// health 是健康检查接口。
func health(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}
