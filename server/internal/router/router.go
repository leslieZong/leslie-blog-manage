package router

import (
	"leslie-blog-server/internal/modules/user"
	"leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// Router 是整个 HTTP 路由的管理对象。
//
// Router 本身不负责业务逻辑。
// 它主要负责：
//
// 1. 创建路由分组
// 2. 注册各个业务模块的路由
// 3. 注册一些全局接口，例如 /health
type Router struct {
	engine *gin.Engine

	// userHandler 是 User 模块的 HTTP Handler。
	//
	// Router 不需要知道 UserHandler 内部是怎么查询数据库的。
	// Router 只需要把 HTTP 请求交给它即可。
	userHandler *handler.UserHandler
}

// New 创建 Router。
//
// 参数：
//
//	engine：Gin HTTP 引擎
//	userHandler：User 模块 Handler
//
// 返回：
//
//	*Router：Router 对象
func New(
	engine *gin.Engine,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		engine:      engine,
		userHandler: userHandler,
	}
}

// Register 注册整个系统的 HTTP 路由。
//
// 可以把这个函数理解成：
//
// “告诉 Gin，整个系统有哪些 API。”
func (r *Router) Register() {

	// ----------------------------------------
	// 1. 健康检查
	// ----------------------------------------
	//
	// GET /health
	//
	// 这个接口通常用于：
	//
	// - Docker 健康检查
	// - Kubernetes Probe
	// - 负载均衡器检查
	// - 运维监控
	r.engine.GET("/health", health)

	// ----------------------------------------
	// 2. Public API
	// ----------------------------------------
	//
	// 例如未来：
	//
	// GET /api/v1/posts
	// GET /api/v1/categories
	// GET /api/v1/projects
	//
	// 这些接口主要服务于 Blog 前台。
	v1 := r.engine.Group("/api/v1")

	// User 模块未来如果有公开接口，
	// 可以注册到这里。
	_ = v1

	// ----------------------------------------
	// 3. Admin API
	// ----------------------------------------
	//
	// Admin API 专门给后台管理系统使用。
	//
	// 例如：
	//
	// POST /api/admin/v1/auth/login
	// GET  /api/admin/v1/users/:id
	// POST /api/admin/v1/posts
	//
	admin := r.engine.Group("/api/admin/v1")

	// ----------------------------------------
	// 4. 注册 User 模块自己的路由
	// ----------------------------------------
	//
	// Router 只负责调用 User 模块的路由注册函数。
	//
	// User 模块到底有哪些接口，
	// 由 user/routes.go 自己决定。
	user.RegisterRoutes(
		admin,
		r.userHandler,
	)
}

// health 是系统健康检查接口。
func health(c *gin.Context) {

	response.Success(c, gin.H{
		"status": "ok",
	})
}
