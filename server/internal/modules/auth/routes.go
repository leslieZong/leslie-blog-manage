package auth

import (
	"leslie-blog-server/internal/modules/auth/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	group *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	jwtMiddleware gin.HandlerFunc,
) {

	// =========================================================
	// 登录接口
	//
	// 登录的时候还没有 JWT。
	//
	// 所以不能使用 JWT Middleware。
	// =========================================================

	group.POST(
		"/auth/login",
		authHandler.Login,
	)

	// =========================================================
	// 当前用户接口
	//
	// /me 必须登录以后才能访问。
	// =========================================================

	group.GET(
		"/auth/me",
		jwtMiddleware,
		authHandler.Me,
	)
}
