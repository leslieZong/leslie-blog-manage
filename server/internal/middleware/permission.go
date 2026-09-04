package middleware

import (
	"net/http"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// Permission 创建一个权限检查 Middleware。
//
// object 表示资源，例如：
//
//	user
//	post
//	category
//
// action 表示操作，例如：
//
//	read
//	create
//	update
//	delete
//
// 最终检查：
//
// 当前用户
//
//	↓
//
// 当前用户角色
//
//	↓
//
// 是否拥有 object + action 权限
func Permission(
	enforcer *casbin.Enforcer,
	object string,
	action string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// =========================================================
		// 第一步：获取当前登录用户 ID
		// =========================================================

		userID := auth.GetUserID(c)

		if userID == "" {

			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"user identity not found",
			)

			c.Abort()

			return
		}

		// =========================================================
		// 第二步：让 Casbin 判断权限
		//
		// Enforce(
		//     谁,
		//     什么资源,
		//     什么操作,
		// )
		//
		// 例如：
		//
		// Enforce(
		//     "01KABC...",
		//     "user",
		//     "read",
		// )
		// =========================================================

		allowed, err := enforcer.Enforce(
			userID,
			object,
			action,
		)

		if err != nil {

			response.Error(
				c,
				http.StatusInternalServerError,
				appErrors.ErrInternalServer,
				"permission check failed",
			)

			c.Abort()

			return
		}

		// =========================================================
		// 第三步：没有权限
		// =========================================================

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

		// =========================================================
		// 第四步：权限验证成功
		// =========================================================

		c.Next()
	}
}
