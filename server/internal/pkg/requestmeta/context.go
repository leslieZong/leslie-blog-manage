package requestmeta

import "context"

type contextKey struct{}

var metadataKey contextKey

// WithMetadata 将请求元数据写入 context。
func WithMetadata(
	ctx context.Context,
	metadata Metadata,
) context.Context {
	return context.WithValue(
		ctx,
		metadataKey,
		metadata,
	)
}

// FromContext 从 context 中读取请求元数据。
func FromContext(ctx context.Context) (Metadata, bool) {
	value := ctx.Value(metadataKey)

	if value == nil {
		return Metadata{}, false
	}

	metadata, ok := value.(Metadata)

	return metadata, ok
}
