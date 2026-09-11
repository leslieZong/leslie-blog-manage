package user

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"

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
		"",
		middleware.Permission(
			enforcer,
			permission.UserRead,
		),
		userHandler.List,
	)

	users.GET(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.UserRead,
		),
		userHandler.GetByID,
	)

	users.POST(
		"",
		middleware.Permission(
			enforcer,
			permission.UserCreate,
		),
		userHandler.Create,
	)

	users.PUT(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.UserUpdate,
		),
		userHandler.Update,
	)

	users.DELETE(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.UserDelete,
		),
		userHandler.Delete,
	)

	// 查询用户角色
	users.GET(
		"/:id/roles",
		middleware.Permission(
			enforcer,
			permission.RoleRead,
		),
		userHandler.GetRoles,
	)

	// 修改用户角色
	users.PUT(
		"/:id/roles",
		middleware.Permission(
			enforcer,
			permission.RoleUpdate,
		),
		userHandler.UpdateRoles,
	)
}
