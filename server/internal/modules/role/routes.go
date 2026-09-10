package role

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/role/handler"
	"leslie-blog-server/internal/pkg/casbin"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 Role API。
func RegisterRoutes(
	group *gin.RouterGroup,
	roleHandler *handler.RoleHandler,
	enforcer *casbin.Enforcer,
) {
	roles := group.Group("/roles")

	// GET /roles
	//
	// 查看角色列表。
	roles.GET(
		"/",
		middleware.Permission(
			enforcer,
			"role",
			"read",
		),
		roleHandler.List,
	)

	// POST /roles
	//
	// 创建角色。
	roles.POST(
		"/",
		middleware.Permission(
			enforcer,
			"role",
			"create",
		),
		roleHandler.Create,
	)

	// PUT /roles/:id
	//
	// 修改角色。
	roles.PUT(
		"/:id",
		middleware.Permission(
			enforcer,
			"role",
			"update",
		),
		roleHandler.Update,
	)

	// DELETE /roles/:id
	//
	// 删除角色。
	roles.DELETE(
		"/:id",
		middleware.Permission(
			enforcer,
			"role",
			"delete",
		),
		roleHandler.Delete,
	)

	roles.GET(
		"/:id/permissions",
		roleHandler.GetPermissions,
	)

	roles.PUT(
		"/:id/permissions",
		roleHandler.UpdatePermissions,
	)
}
