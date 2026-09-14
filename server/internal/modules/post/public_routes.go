package post

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/post/handler"
)

// RegisterPublicRoutes
//
// 注册博客前台 Post API。
//
// 注意：
//
// 这里不添加 JWT Middleware。
// 也不添加 Casbin Permission Middleware。
//
// 因为博客前台允许游客访问。
func RegisterPublicRoutes(
	router *gin.RouterGroup,
	postHandler *handler.PublicPostHandler,
) {
	posts := router.Group("/posts")

	// 获取已发布文章列表
	posts.GET(
		"",
		postHandler.ListPublished,
	)

	// 根据 slug 获取文章
	posts.GET(
		"/slug/:slug",
		postHandler.GetBySlug,
	)

	// 根据 ID 获取文章
	posts.GET(
		"/:id",
		postHandler.GetByID,
	)

	// 增加阅读量
	posts.POST(
		"/:id/view",
		postHandler.IncrementViewCount,
	)
}
