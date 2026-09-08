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

	// GET /roles
	//
	// 查看角色列表。
	group.GET(
		"/roles",
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
	group.POST(
		"/roles",
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
	group.PUT(
		"/roles/:id",
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
	group.DELETE(
		"/roles/:id",
		middleware.Permission(
			enforcer,
			"role",
			"delete",
		),
		roleHandler.Delete,
	)
}
