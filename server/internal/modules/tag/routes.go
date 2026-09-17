package tag

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/tag/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 Tag 管理路由。
//
// 所有路由最终会挂载到：
//
// /api/admin/v1
//
// 因此这里的：
//
// /tags
//
// 最终就是：
//
// /api/admin/v1/tags
func RegisterRoutes(
	group *gin.RouterGroup,
	tagHandler *handler.TagHandler,
	enforcer *casbin.Enforcer,
) {

	// -----------------------------------------------------
	// Tag 路由组
	// -----------------------------------------------------

	tags := group.Group("/tags")

	// -----------------------------------------------------
	// 查询 Tag 列表
	// -----------------------------------------------------
	//
	// GET /tags
	//
	// 需要：
	//
	// tag:read
	//
	tags.GET(
		"",

		middleware.Permission(
			enforcer,
			permission.TagRead,
		),
		tagHandler.List,
	)

	// -----------------------------------------------------
	// 查询 Tag 详情
	// -----------------------------------------------------
	//
	// GET /tags/:id
	//
	// 需要：
	//
	// tag:read
	//
	tags.GET(
		"/:id",

		middleware.Permission(
			enforcer,
			permission.TagRead,
		),
		tagHandler.GetByID,
	)

	// -----------------------------------------------------
	// 创建 Tag
	// -----------------------------------------------------
	//
	// POST /tags
	//
	// 需要：
	//
	// tag:create
	//
	tags.POST(
		"",

		middleware.Permission(
			enforcer,
			permission.TagCreate,
		),
		tagHandler.Create,
	)

	// -----------------------------------------------------
	// 修改 Tag
	// -----------------------------------------------------
	//
	// PUT /tags/:id
	//
	// 需要：
	//
	// tag:update
	//
	tags.PUT(
		"/:id",

		middleware.Permission(
			enforcer,
			permission.TagUpdate,
		),
		tagHandler.Update,
	)

	// -----------------------------------------------------
	// 删除 Tag
	// -----------------------------------------------------
	//
	// DELETE /tags/:id
	//
	// 需要：
	//
	// tag:delete
	//
	tags.DELETE(
		"/:id",

		middleware.Permission(
			enforcer,
			permission.TagDelete,
		),
		tagHandler.Delete,
	)
}
