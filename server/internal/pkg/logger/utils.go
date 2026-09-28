package logger

import (
	"context"
	"log/slog"
)

// WithContext 给 Logger 自动增加 Request ID。
func WithContext(
	ctx context.Context,
	log *Logger,
) *slog.Logger {

	requestID := RequestID(ctx)

	if requestID == "" {
		return log.Logger
	}

	return log.With(
		"request_id",
		requestID,
	)
}
