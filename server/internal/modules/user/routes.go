package user

import (
	"leslie-blog-server/internal/modules/user/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 User 模块路由。
//
// 这里负责的是：
//
// URL
// ↓
// Handler
//
// 它不负责：
//
// 数据库
// 业务逻辑
// JWT
// SQL
func RegisterRoutes(
	group *gin.RouterGroup,
	userHandler *handler.UserHandler,
) {

	// GET /users/:id
	//
	// 最终映射到：
	//
	// UserHandler.GetByID
	group.GET(
		"/users/:id",
		userHandler.GetByID,
	)
}
