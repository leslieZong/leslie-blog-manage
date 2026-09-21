package project

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/project/handler"
)

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	projectHandler *handler.ProjectPublicHandler,
) {

	projects := router.Group("/projects")

	// Project 列表。
	projects.GET(
		"",
		projectHandler.List,
	)

	// 注意：
	// /slug/:slug 应该放在 /:id 之前。
	projects.GET(
		"/slug/:slug",
		projectHandler.GetBySlug,
	)

	// Project 详情。
	projects.GET(
		"/:id",
		projectHandler.GetByID,
	)
}
