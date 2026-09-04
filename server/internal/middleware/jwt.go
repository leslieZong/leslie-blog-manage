package middleware

import (
	"net/http"
	"strings"

	"leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/jwt"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// JWT 创建一个 JWT 认证中间件。
//
// secret 是生成 JWT 时使用的密钥。
// Middleware 在收到请求之后，会使用同一个 secret
// 验证客户端提交过来的 JWT。
func JWT(secret string) gin.HandlerFunc {

	// gin.HandlerFunc 本质上就是一个函数类型：
	//
	// func(c *gin.Context)
	//
	// 也就是说：
	//
	// JWT(secret)
	//     ↓
	// 返回一个函数
	//     ↓
	// Gin 在请求到达时执行这个函数
	return func(c *gin.Context) {

		// ==================================================
		// 第一步：获取 Authorization Header
		// ==================================================
		//
		// 前端请求应该类似：
		//
		// Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
		//
		// c.GetHeader() 用来读取 HTTP Header。
		authHeader := c.GetHeader("Authorization")

		// 如果没有 Authorization Header，
		// 说明客户端没有提供 Token。
		if authHeader == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				errors.ErrUnauthorized,
				"authorization header is required",
			)

			// Abort 非常重要。
			//
			// 它表示：
			//
			// “请求到这里就停止，不允许继续执行后面的 Handler。”
			c.Abort()

			return
		}

		// ==================================================
		// 第二步：检查 Bearer 格式
		// ==================================================
		//
		// 我们要求：
		//
		// Authorization: Bearer <token>
		//
		// 所以 Header 必须以：
		//
		// "Bearer "
		//
		// 开头。
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(
				c,
				http.StatusUnauthorized,
				errors.ErrUnauthorized,
				"invalid authorization header",
			)

			c.Abort()

			return
		}

		// ==================================================
		// 第三步：提取真正的 JWT Token
		// ==================================================
		//
		// authHeader：
		//
		// Bearer eyJhbGciOiJIUzI1NiIs...
		//
		// TrimPrefix 后：
		//
		// eyJhbGciOiJIUzI1NiIs...
		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		if tokenString == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				errors.ErrUnauthorized,
				"token is required",
			)

			c.Abort()

			return
		}

		// ==================================================
		// 第四步：解析 JWT
		// ==================================================
		//
		// Parse 会检查：
		//
		// 1. JWT 格式
		// 2. 签名
		// 3. Token 是否过期
		// 4. Claims 是否能够正确解析
		claims, err := jwt.Parse(
			tokenString,
			secret,
		)

		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				errors.ErrUnauthorized,
				"invalid or expired token",
			)

			c.Abort()

			return
		}

		// ==================================================
		// 第五步：把登录用户信息放进 Gin Context
		// ==================================================
		//
		// claims 中已经有：
		//
		// UserID
		// Username
		//
		// 我们把它们保存下来。
		//
		// 后面的 Handler / Service 就可以读取。
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// ==================================================
		// 第六步：继续执行后面的 Handler
		// ==================================================
		//
		// 到这里说明：
		//
		// JWT 验证通过。
		//
		// 所以允许请求继续执行。
		c.Next()
	}
}
