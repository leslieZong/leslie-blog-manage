package middleware

import (
	"net/http"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/jwt"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// JWT 是 JWT 身份认证 Middleware。
//
// 它负责：
//
// 1. 获取 Authorization 请求头
// 2. 提取 Bearer Token
// 3. 解析 JWT
// 4. 验证 Token 是否有效
// 5. 获取当前用户身份
// 6. 把用户身份保存到 Gin Context
//
// 注意：
//
// JWT Middleware 只负责“你是谁”。
//
// 它不负责判断：
//
// “你能不能删除文章？”
//
// 权限判断由后面的 Permission Middleware + Casbin 完成。
func JWT(secret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// =========================================================
		// 第一步：获取 Authorization 请求头
		// =========================================================

		authorization := c.GetHeader("Authorization")

		if authorization == "" {

			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"authorization header is required",
			)

			// Abort 表示终止当前请求继续向后执行。
			c.Abort()

			return
		}

		// =========================================================
		// 第二步：检查 Bearer 格式
		//
		// HTTP 请求通常：
		//
		// Authorization: Bearer xxxxx
		// =========================================================

		if !strings.HasPrefix(authorization, "Bearer ") {

			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"invalid authorization header",
			)

			c.Abort()

			return
		}

		// =========================================================
		// 第三步：提取真正的 JWT Token
		// =========================================================

		tokenString := strings.TrimPrefix(
			authorization,
			"Bearer ",
		)

		if tokenString == "" {

			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"token is required",
			)

			c.Abort()

			return
		}

		// =========================================================
		// 第四步：解析 JWT
		// =========================================================

		claims, err := jwt.Parse(
			tokenString,
			secret,
		)

		if err != nil {

			response.Error(
				c,
				http.StatusUnauthorized,
				appErrors.ErrUnauthorized,
				"invalid or expired token",
			)

			c.Abort()

			return
		}

		// =========================================================
		// 第五步：把当前用户身份保存到 Gin Context
		// =========================================================
		//
		// 后面的 Handler 可以通过：
		//
		// auth.GetUserID(c)
		//
		// 获取用户 ID。
		//

		c.Set(
			auth.ContextKeyUserID,
			claims.UserID,
		)

		c.Set(
			auth.ContextKeyUsername,
			claims.Username,
		)

		// =========================================================
		// 第六步：认证成功
		//
		// c.Next() 表示：
		//
		// “JWT Middleware 已经处理完毕，
		// 继续执行后面的 Middleware / Handler。”
		// =========================================================

		c.Next()
	}
}
