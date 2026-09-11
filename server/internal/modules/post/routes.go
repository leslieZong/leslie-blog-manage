package post

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/post/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	postHandler *handler.PostHandler,
	enforcer *casbin.Enforcer,
) {
	// 创建 Admin 文章路由分组。
	//
	// router 本身应该已经是：
	//
	// /api/admin/v1
	//
	// 所以这里再加：
	//
	// /posts
	//
	// 最终：
	//
	// /api/admin/v1/posts
	posts := router.Group("/posts")

	// 创建文章。
	posts.POST(
		"",
		middleware.Permission(enforcer, permission.PostCreate),
		postHandler.Create,
	)

	// 获取文章列表。
	posts.GET(
		"",
		middleware.Permission(enforcer, permission.PostRead),
		postHandler.List,
	)

	// 获取文章详情。
	posts.GET(
		"/:id",
		middleware.Permission(enforcer, permission.PostRead),
		postHandler.GetByID,
	)

	posts.GET(
		"/slug/:slug",
		middleware.Permission(enforcer, permission.PostRead),
		postHandler.GetBySlug,
	)

	// 修改文章。
	posts.PUT(
		"/:id",
		middleware.Permission(enforcer, permission.PostUpdate),
		postHandler.Update,
	)

	// 删除文章。
	posts.DELETE(
		"/:id",
		middleware.Permission(enforcer, permission.PostDelete),
		postHandler.Delete,
	)

	// 发布文章。
	posts.POST(
		"/:id/publish",
		middleware.Permission(enforcer, permission.PostPublish),
		postHandler.Publish,
	)
}
