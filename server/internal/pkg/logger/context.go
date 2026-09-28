package logger

import "context"

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
