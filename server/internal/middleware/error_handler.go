package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/pkg/logger"
	"leslie-blog-server/internal/response"
)

// ErrorHandler 负责统一处理 HTTP 请求中的业务错误。
//
// 它应该放在 Router Middleware 中，
// 负责捕获 Handler 通过 c.Error(err)
// 传递出来的错误。
func ErrorHandler(
	appLogger *logger.Logger,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 先让后面的 Handler / Middleware 执行。
		c.Next()

		// 如果没有错误，
		// 直接结束。
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误。
		//
		// 一个请求理论上可能产生多个错误，
		// 第一版我们只处理最后一个。
		err := c.Errors.Last().Err

		// 将普通 error 转换成 AppError。
		appErr := appErrors.FromError(err)

		// 获取带 Request ID 的 Logger。
		log := logger.WithContext(
			c.Request.Context(),
			appLogger,
		)

		// 业务错误和系统错误的日志处理方式可以不同。
		//
		// 例如：
		//
		// 404 -> WARN
		// 401 -> WARN
		// 403 -> WARN
		// 500 -> ERROR
		if appErr.HTTPStatus >= http.StatusInternalServerError {

			log.Error(
				"http request failed",
				slog.Int(
					"error_code",
					int(appErr.Code),
				),
				slog.String(
					"error_message",
					appErr.Message,
				),
				slog.Any(
					"error",
					err,
				),
			)

		} else {

			log.Warn(
				"http request failed",
				slog.Int(
					"error_code",
					int(appErr.Code),
				),
				slog.String(
					"error_message",
					appErr.Message,
				),
			)
		}

		// 如果 Handler 已经写过响应，
		// 就不要再次写 JSON。
		if c.Writer.Written() {
			return
		}

		// 最终统一输出 JSON。
		response.Error(
			c,
			appErr.HTTPStatus,
			appErr.Code,
			appErr.Message,
		)
	}
}
