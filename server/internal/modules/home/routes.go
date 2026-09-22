package home

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/home/handler"
)

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	homeHandler *handler.HomeHandler,
) {

	router.GET(
		"/home",
		homeHandler.GetHome,
	)
}
