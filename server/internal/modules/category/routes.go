package category

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/category/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"
)

// RegisterRoutes 注册 Category API。
//
// 这里统一负责：
//
// JWT
// +
// Casbin Permission
// +
// Handler
//
// Handler 本身不关心权限。
func RegisterRoutes(
	router *gin.RouterGroup,
	categoryHandler *handler.CategoryHandler,
	enforcer *casbin.Enforcer,
) {

	// 创建 categories 路由组。
	//
	// 外部最终可能是：
	//
	// /api/admin/v1/categories
	categories := router.Group("/categories")

	// 获取分类列表。
	categories.GET(
		"",
		middleware.Permission(enforcer, permission.CategoryRead),
		categoryHandler.List,
	)

	// 获取单个分类。
	categories.GET(
		"/:id",
		middleware.Permission(enforcer, permission.CategoryRead),
		categoryHandler.GetByID,
	)

	// 创建分类。
	categories.POST(
		"",
		middleware.Permission(enforcer, permission.CategoryCreate),
		categoryHandler.Create,
	)

	// 修改分类。
	categories.PUT(
		"/:id",
		middleware.Permission(enforcer, permission.CategoryUpdate),
		categoryHandler.Update,
	)

	// 删除分类。
	categories.DELETE(
		"/:id",
		middleware.Permission(enforcer, permission.CategoryDelete),
		categoryHandler.Delete,
	)
}

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	categoryHandler *handler.PublicCategoryHandler,
) {
	categories := router.Group("/categories")

	categories.GET(
		"",
		categoryHandler.List,
	)
}
