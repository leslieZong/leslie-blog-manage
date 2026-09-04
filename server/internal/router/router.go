package router

import (
	"net/http"

	"leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// Router 负责整个应用的 HTTP 路由。
type Router struct {

	// Gin Engine。
	engine *gin.Engine

	// UserHandler。
	//
	// Router 不创建 Handler。
	//
	// Handler 是外部注入进来的。
	userHandler *handler.UserHandler
}

// New 创建 Router。
//
// 注意：
//
// 现在 Router 不仅需要 Gin Engine，
// 还需要各个业务模块的 Handler。
func New(
	engine *gin.Engine,
	userHandler *handler.UserHandler,
) *Router {

	return &Router{
		engine:      engine,
		userHandler: userHandler,
	}
}

// Register 注册所有 HTTP 路由。
func (r *Router) Register() {

	// --------------------------------------------------
	// Health
	// --------------------------------------------------

	r.engine.GET(
		"/health",
		health,
	)

	// --------------------------------------------------
	// API V1
	// --------------------------------------------------

	v1 := r.engine.Group("/api/v1")

	_ = v1

	// --------------------------------------------------
	// Admin API V1
	// --------------------------------------------------

	admin := r.engine.Group("/api/admin/v1")

	// User
	//
	// /api/admin/v1/users/:id
	admin.GET(
		"/users/:id",
		r.userHandler.GetByID,
	)
}

// health 健康检查。
func health(c *gin.Context) {

	response.Success(
		c,
		gin.H{
			"status": "ok",
		},
	)
}

// exampleError 示例错误。
func exampleError(c *gin.Context) {

	response.Error(
		c,
		http.StatusBadRequest,
		errors.ErrInvalidParams,
		"invalid parameters",
	)
}
