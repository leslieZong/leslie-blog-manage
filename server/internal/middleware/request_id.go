package middleware

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/pkg/logger"
	"leslie-blog-server/internal/pkg/ulid"
)

// RequestID 为每一个 HTTP 请求生成唯一 Request ID。
func RequestID() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 优先读取客户端已经传入的 Request ID。
		//
		// 这样如果未来：
		//
		// Nginx
		//   ↓
		// Go
		//
		// Nginx 已经生成 Request ID，
		// Go 可以继续沿用。
		requestID := c.GetHeader(
			"X-Request-ID",
		)

		// 如果客户端没有提供，
		// 则由 Go 服务生成。
		if requestID == "" {

			requestID = ulid.New()
		}

		// 写入 Gin Context。
		c.Set(
			"request_id",
			requestID,
		)

		// 同时写入标准 context。
		c.Request = c.Request.WithContext(
			logger.WithRequestID(
				c.Request.Context(),
				requestID,
			),
		)

		// 返回响应时也带上 Request ID。
		c.Header(
			"X-Request-ID",
			requestID,
		)

		c.Next()
	}
}
