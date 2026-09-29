package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// WithRequestID 将 Request ID 放入 Context。
func WithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {

	return context.WithValue(
		ctx,
		requestIDKey,
		requestID,
	)
}

// RequestID 从 Context 获取 Request ID。
func RequestID(
	ctx context.Context,
) string {

	value := ctx.Value(requestIDKey)

	if value == nil {
		return ""
	}

	requestID, ok := value.(string)

	if !ok {
		return ""
	}

	return requestID
}

// InfoContext 记录 INFO 日志。
func InfoContext(
	logger *Logger,
	ctx context.Context,
	message string,
	args ...any,
) {
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		message,
		buildAttrs(args...)...,
	)
}

// ErrorContext 记录 ERROR 日志。
func ErrorContext(
	logger *Logger,
	ctx context.Context,
	message string,
	args ...any,
) {
	logger.LogAttrs(
		ctx,
		slog.LevelError,
		message,
		buildAttrs(args...)...,
	)
}

// WarnContext 记录 WARN 日志。
func WarnContext(
	logger *Logger,
	ctx context.Context,
	message string,
	args ...any,
) {
	logger.LogAttrs(
		ctx,
		slog.LevelWarn,
		message,
		buildAttrs(args...)...,
	)
}

// DebugContext 记录 DEBUG 日志。
func DebugContext(
	logger *Logger,
	ctx context.Context,
	message string,
	args ...any,
) {
	logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		message,
		buildAttrs(args...)...,
	)
}

// buildAttrs 把 Context 中的请求级信息
// 转换成结构化日志字段。
func buildAttrs(
	args ...any,
) []slog.Attr {

	attrs := make([]slog.Attr, 0, len(args)/2+2)

	// ------------------------------------------------
	// 其他业务日志字段
	// ------------------------------------------------

	attrs = append(
		attrs,
		anyToAttrs(args...)...,
	)

	return attrs
}

// anyToAttrs 将 slog 参数转换为 Attr。
//
// 调用方式：
//
// logger.InfoContext(
//
//	logger,
//	ctx,
//	"post created",
//	slog.String("post_id", post.ID),
//
// )
func anyToAttrs(args ...any) []slog.Attr {

	attrs := make([]slog.Attr, 0, len(args))

	for _, arg := range args {

		if attr, ok := arg.(slog.Attr); ok {
			attrs = append(attrs, attr)
		}
	}

	return attrs
}
