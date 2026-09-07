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
		"/users",
		middleware.Permission(
			enforcer,
			"user",
			"read",
		),
		userHandler.List,
	)

	group.GET(
		"/users/:id",
		middleware.Permission(
			enforcer,
			"user",
			"read",
		),
		userHandler.GetByID,
	)

	group.POST(
		"/users",
		middleware.Permission(
			enforcer,
			"user",
			"create",
		),
		userHandler.Create,
	)

	group.PUT(
		"/users/:id",
		middleware.Permission(
			enforcer,
			"user",
			"update",
		),
		userHandler.Update,
	)

	group.DELETE(
		"/users/:id",
		middleware.Permission(
			enforcer,
			"user",
			"delete",
		),
		userHandler.Delete,
	)
}
