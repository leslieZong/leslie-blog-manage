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
	users := group.Group("/users")

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
	users.GET(
		"/",
		middleware.Permission(
			enforcer,
			"user",
			"read",
		),
		userHandler.List,
	)

	users.GET(
		"/:id",
		middleware.Permission(
			enforcer,
			"user",
			"read",
		),
		userHandler.GetByID,
	)

	users.POST(
		"/",
		middleware.Permission(
			enforcer,
			"user",
			"create",
		),
		userHandler.Create,
	)

	users.PUT(
		"/:id",
		middleware.Permission(
			enforcer,
			"user",
			"update",
		),
		userHandler.Update,
	)

	users.DELETE(
		"/:id",
		middleware.Permission(
			enforcer,
			"user",
			"delete",
		),
		userHandler.Delete,
	)

	// 查询用户角色
	users.GET(
		"/:id/roles",
		middleware.Permission(
			enforcer,
			"role",
			"read",
		),
		userHandler.GetRoles,
	)

	// 修改用户角色
	users.PUT(
		"/:id/roles",
		middleware.Permission(
			enforcer,
			"role",
			"update",
		),
		userHandler.UpdateRoles,
	)
}
