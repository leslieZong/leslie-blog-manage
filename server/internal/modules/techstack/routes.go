package techstack

import (
	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/techstack/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *handler.TechStackHandler,
	enforcer *casbin.Enforcer,
) {

	techStacks := router.Group("/tech-stacks")

	techStacks.GET(
		"",
		middleware.Permission(
			enforcer,
			permission.TechStackRead,
		),
		handler.List,
	)

	techStacks.GET(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.TechStackRead,
		),
		handler.GetByID,
	)

	techStacks.POST(
		"",
		middleware.Permission(
			enforcer,
			permission.TechStackCreate,
		),
		handler.Create,
	)

	techStacks.PUT(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.TechStackUpdate,
		),
		handler.Update,
	)

	techStacks.DELETE(
		"/:id",
		middleware.Permission(
			enforcer,
			permission.TechStackDelete,
		),
		handler.Delete,
	)
}
