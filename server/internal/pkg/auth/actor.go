package auth

import "github.com/gin-gonic/gin"

// Actor 表示当前正在执行操作的用户。
//
// 它不是数据库 User Model。
// 它表示“当前请求的操作者是谁”。
type Actor struct {

	// 用户 ID。
	UserID string

	// 是否拥有管理员身份。
	IsAdmin bool
}

// ActorFromContext 从 Gin Context 获取当前操作者。
//
// 注意：
// 这里属于 HTTP / Auth 适配层。
// Service 不应该直接调用这个函数。
func ActorFromContext(
	c *gin.Context,
	isAdmin bool,
) (Actor, error) {

	userID := GetUserID(c)

	return Actor{
		UserID:  userID,
		IsAdmin: isAdmin,
	}, nil
}
