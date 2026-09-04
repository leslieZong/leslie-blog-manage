package user

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/pkg/casbin"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 User 模块路由。
func RegisterRoutes(
	group *gin.RouterGroup,
	userHandler *handler.UserHandler,
	enforcer *casbin.Enforcer,
) {

	// ==================================================
	// 获取用户详情
	// ==================================================
	//
	// 请求：
	//
	// GET /users/:id
	//
	// 需要：
	//
	// JWT
	// +
	// user:read
	group.GET(
		"/users/:id",

		// 权限 Middleware。
		middleware.Permission(
			enforcer,
			"user",
			"read",
		),

		// 真正的业务 Handler。
		userHandler.GetByID,
	)
}
