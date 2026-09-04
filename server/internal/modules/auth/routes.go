package auth

import (
	"leslie-blog-server/internal/modules/auth/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 Auth 模块路由。
func RegisterRoutes(
	group *gin.RouterGroup,
	authHandler *handler.AuthHandler,
) {

	// POST /auth/login
	//
	// 因为上层 group 是：
	//
	// /api/admin/v1
	//
	// 所以最终完整路径：
	//
	// POST /api/admin/v1/auth/login
	group.POST(
		"/auth/login",
		authHandler.Login,
	)
}
