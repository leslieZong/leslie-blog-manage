package user

import (
	"leslie-blog-server/internal/modules/user/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 User 模块的 HTTP 路由。
//
// group 是上层 Router 创建好的路由分组。
//
// 例如：
//
//	/api/admin/v1
//
// userHandler 是 User 模块的 Handler。
func RegisterRoutes(
	group *gin.RouterGroup,
	userHandler *handler.UserHandler,
) {

	// GET /users/:id
	//
	// 最终完整路径：
	//
	// GET /api/admin/v1/users/:id
	group.GET(
		"/users/:id",
		userHandler.GetByID,
	)
}
