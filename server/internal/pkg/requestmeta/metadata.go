package requestmeta

// Metadata 保存一次 HTTP 请求相关的元数据。
//
// 注意：
//
// 这些信息不是业务数据。
// 它们描述的是：
//
// “这一次请求是谁发起的、从哪里来、属于哪一次请求。”
type Metadata struct {
	// RequestID
	//
	// 当前 HTTP 请求唯一 ID。
	//
	// 例如：
	// 01K6ABC123...
	RequestID string

	// UserID
	//
	// 当前登录用户 ID。
	//
	// 未登录请求可能为空。
	UserID string

	// IP
	//
	// 客户端 IP 地址。
	IP string

	// UserAgent
	//
	// HTTP User-Agent。
	UserAgent string
}
