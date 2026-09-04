package middleware

import (
	"net/http"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// Permission 创建权限检查 Middleware。
//
// object 表示资源，例如：
//
// post
// user
// category
//
// action 表示操作，例如：
//
// read
// create
// update
// delete
func Permission(
	enforcer *casbin.Enforcer,
	object string,
	action string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// ==================================================
		// 1. 获取 JWT Middleware 设置的 userID
		// ==================================================

		userIDValue, exists := c.Get("userID")

		if !exists {
			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"user identity not found",
			)

			c.Abort()
			return
		}

		// ==================================================
		// 2. 类型断言
		// ==================================================
		//
		// c.Get() 返回的是 any。
		//
		// 我们知道 JWT Middleware 放进去的是 string。
		//
		// 所以需要：
		//
		// userID, ok := userIDValue.(string)

		userID, ok := userIDValue.(string)

		if !ok || userID == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"invalid user identity",
			)

			c.Abort()
			return
		}

		// ==================================================
		// 3. 调用 Casbin
		// ==================================================

		allowed, err := enforcer.Enforce(
			userID,
			object,
			action,
		)

		if err != nil {

			// Casbin 自身发生异常。
			response.Error(
				c,
				http.StatusInternalServerError,
				appErrors.ErrInternalServer,
				"permission check failed",
			)

			c.Abort()
			return
		}

		// ==================================================
		// 4. 权限不足
		// ==================================================

		if !allowed {

			response.Error(
				c,
				http.StatusForbidden,
				appErrors.ErrForbidden,
				"permission denied",
			)

			c.Abort()
			return
		}

		// ==================================================
		// 5. 权限通过
		// ==================================================

		c.Next()
	}
}
