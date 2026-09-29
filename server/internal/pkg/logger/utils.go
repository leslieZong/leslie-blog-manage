package logger

import (
	"leslie-blog-server/internal/pkg/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// WithContext 给 Logger 自动增加 Request ID。
func WithContext(
	ctx *gin.Context,
	log *Logger,
) *slog.Logger {
	c := ctx.Request.Context()
	requestID := RequestID(c)

	if requestID == "" {
		return log.Logger
	}

	return log.With(
		"request_id",
		requestID,
		"user_id",
		auth.GetUserID(ctx),
	)
}
