package project

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/middleware"
	"leslie-blog-server/internal/modules/project/handler"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/permission"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	projectHandler *handler.ProjectHandler,
	enforcer *casbin.Enforcer,
) {

	projects := router.Group("/projects")

	// ---------------------------------------------
	// List
	// ---------------------------------------------

	projects.GET(
		"",
		middleware.Permission(enforcer, permission.ProjectRead),
		projectHandler.List,
	)

	// ---------------------------------------------
	// Get
	// ---------------------------------------------

	projects.GET(
		"/:id",
		middleware.Permission(enforcer, permission.ProjectRead),
		projectHandler.GetByID,
	)

	// ---------------------------------------------
	// Create
	// ---------------------------------------------

	projects.POST(
		"",
		middleware.Permission(enforcer, permission.ProjectCreate),
		projectHandler.Create,
	)

	// ---------------------------------------------
	// Update
	// ---------------------------------------------

	projects.PUT(
		"/:id",
		middleware.Permission(enforcer, permission.ProjectUpdate),
		projectHandler.Update,
	)

	// ---------------------------------------------
	// Delete
	// ---------------------------------------------

	projects.DELETE(
		"/:id",
		middleware.Permission(enforcer, permission.ProjectDelete),
		projectHandler.Delete,
	)
}
