package middleware

import (
	"leslie-blog-server/internal/pkg/requestmeta"

	"github.com/gin-gonic/gin"
)

// RequestMetaMiddleware 创建请求元数据 Middleware。
//
// 它负责收集：
//
// 1. Request ID
// 2. User ID
// 3. Client IP
// 4. User-Agent
//
// 然后把这些信息写入标准 context.Context。
func RequestMetaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// --------------------------------------------------
		// 3. 获取客户端 IP
		// --------------------------------------------------
		ip := c.ClientIP()

		// --------------------------------------------------
		// 4. 获取 User-Agent
		// --------------------------------------------------
		userAgent := c.Request.UserAgent()

		// --------------------------------------------------
		// 5. 组装 Metadata
		// --------------------------------------------------
		metadata := requestmeta.Metadata{
			IP:        ip,
			UserAgent: userAgent,
		}

		// --------------------------------------------------
		// 6. 将 Metadata 放入标准 context.Context
		// --------------------------------------------------
		ctx := requestmeta.WithMetadata(
			c.Request.Context(),
			metadata,
		)

		// --------------------------------------------------
		// 7. 将新的 context 设置回 HTTP Request
		// --------------------------------------------------
		c.Request = c.Request.WithContext(ctx)

		// --------------------------------------------------
		// 8. 继续执行后面的 Middleware / Handler
		// --------------------------------------------------
		c.Next()
	}
}
