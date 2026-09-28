package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/pkg/logger"
)

// RequestLogger 记录每个 HTTP 请求的基本信息。
func RequestLogger(
	appLogger *logger.Logger,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		// 请求执行完成之后，
		// 才能拿到最终 HTTP Status。
		latency := time.Since(start)

		log := logger.WithContext(
			c.Request.Context(),
			appLogger,
		)

		log.Info(
			"http request",
			slog.String(
				"method",
				c.Request.Method,
			),
			slog.String(
				"path",
				c.Request.URL.Path,
			),
			slog.Int(
				"status",
				c.Writer.Status(),
			),
			slog.Duration(
				"latency",
				latency,
			),
		)
	}
}
