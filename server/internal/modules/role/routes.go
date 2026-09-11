package role

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/role/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"

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
		"",
		middleware.Permission(
			enforcer,
			permission.RoleRead,
		),
		roleHandler.List,
	)

	// POST /roles
	//
	// 创建角色。
	roles.POST(
		"",
		middleware.Permission(
			enforcer,
			permission.RoleCreate,
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
			permission.RoleUpdate,
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
			permission.RoleDelete,
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
