package router

import (
	"leslie-blog-server/internal/modules/auth"
	authHandler "leslie-blog-server/internal/modules/auth/handler"
	"leslie-blog-server/internal/modules/user"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// Router 管理整个 HTTP API。
type Router struct {
	engine *gin.Engine

	userHandler *userHandler.UserHandler

	authHandler *authHandler.AuthHandler
}

// New 创建 Router。
func New(
	engine *gin.Engine,
	userHandler *userHandler.UserHandler,
	authHandler *authHandler.AuthHandler,
) *Router {

	return &Router{
		engine:      engine,
		userHandler: userHandler,
		authHandler: authHandler,
	}
}

// Register 注册所有系统路由。
func (r *Router) Register() {

	// 健康检查。
	r.engine.GET("/health", health)

	// ----------------------------------------
	// Public API
	// ----------------------------------------
	v1 := r.engine.Group("/api/v1")

	_ = v1

	// ----------------------------------------
	// Admin API
	// ----------------------------------------
	admin := r.engine.Group("/api/admin/v1")

	// Auth 路由。
	//
	// 登录接口暂时不需要 JWT，
	// 因为用户此时还没有 Token。
	auth.RegisterRoutes(
		admin,
		r.authHandler,
	)

	// User 路由。
	//
	// 后面这里会增加 JWT Middleware。
	user.RegisterRoutes(
		admin,
		r.userHandler,
	)
}

// health 健康检查接口。
func health(c *gin.Context) {

	response.Success(c, gin.H{
		"status": "ok",
	})
}
