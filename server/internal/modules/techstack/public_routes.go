package techstack

import (
	"leslie-blog-server/internal/modules/techstack/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	handler *handler.TechStackHandler,
) {

	techStacks := router.Group("/tech-stacks")

	techStacks.GET(
		"",
		handler.PublicList,
	)

	techStacks.GET(
		"/:slug",
		handler.PublicGetBySlug,
	)
}
