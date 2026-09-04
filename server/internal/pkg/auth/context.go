package auth

import "github.com/gin-gonic/gin"

// ContextKeyUserID 是 Gin Context 中保存用户 ID 时使用的 key。
//
// JWT Middleware 在解析 Token 成功以后，会执行：
//
//	c.Set(ContextKeyUserID, claims.UserID)
//
// 后面的 Handler / Service 就可以通过 GetUserID 获取当前用户 ID。
const ContextKeyUserID = "userID"

// ContextKeyUsername 是 Gin Context 中保存用户名时使用的 key。
//
// JWT Middleware 会执行：
//
//	c.Set(ContextKeyUsername, claims.Username)
const ContextKeyUsername = "username"

// GetUserID 获取当前登录用户的 ID。
//
// 这里需要理解 Gin Context 中的数据类型是 any，
// 所以我们需要进行类型断言：
//
//	value.(string)
//
// 如果不存在，或者类型不是 string，返回空字符串。
func GetUserID(c *gin.Context) string {
	value, exists := c.Get(ContextKeyUserID)

	// Context 中没有 userID。
	if !exists {
		return ""
	}

	// 尝试把 any 转换成 string。
	userID, ok := value.(string)

	// 类型不是 string。
	if !ok {
		return ""
	}

	return userID
}

// GetUsername 获取当前登录用户的用户名。
func GetUsername(c *gin.Context) string {
	value, exists := c.Get(ContextKeyUsername)

	// Context 中没有 username。
	if !exists {
		return ""
	}

	// 尝试进行类型断言。
	username, ok := value.(string)

	// 类型不是 string。
	if !ok {
		return ""
	}

	return username
}
