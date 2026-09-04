package errors

// Leslie Blog 业务错误码。
//
// 错误码建议按照业务模块进行划分。
//
// 例如：
//
// 400xx → 参数相关
// 401xx → 认证相关
// 403xx → 权限相关
// 404xx → 资源不存在
// 500xx → 系统错误
const (

	// -----------------------------
	// 通用错误
	// -----------------------------

	// ErrInvalidParams 表示请求参数错误。
	ErrInvalidParams = 40001

	// ErrUnauthorized 表示用户没有登录。
	ErrUnauthorized = 40101

	// ErrForbidden 表示用户没有权限。
	ErrForbidden = 40301

	// ErrNotFound 表示资源不存在。
	ErrNotFound = 40401

	// ErrInternalServer 表示服务器内部错误。
	ErrInternalServer = 50001
)
