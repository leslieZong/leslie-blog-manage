package errors

import "fmt"

// ============================================================
// 1. 定义业务错误码
// ============================================================
//
// 这里的 code 是返回给前端的业务错误码。
//
// 注意：
// HTTP Status 和业务 Code 是两套东西。
//
// 例如：
//
// HTTP 404
// Code 40401
//
// HTTP 401
// Code 40101
//
// 前端既可以根据 HTTP Status 判断大类别，
// 也可以根据 Code 判断具体业务错误。
const (
	// 通用错误
	ErrInvalidParams  = 40001
	ErrUnauthorized   = 40101
	ErrForbidden      = 40301
	ErrNotFound       = 40401
	ErrInternalServer = 50001

	// 用户相关错误
	ErrUserNotFound   = 40402
	ErrUsernameExists = 40002
	ErrUserDisabled   = 40302

	// 认证相关错误
	ErrInvalidCredentials = 40102
)

// AppError 是 Leslie Blog 项目统一的业务错误。
//
// 它比普通 errors.New() 多保存了一些信息：
//
// Code
//
//	→ 前端业务错误码
//
// HTTPStatus
//
//	→ HTTP 状态码
//
// Message
//
//	→ 给前端展示的错误信息
//
// Err
//
//	→ 原始错误，可选
type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

// Error 实现 Go 标准 error 接口。
//
// 只要一个类型实现：
//
// # Error() string
//
// 它就可以当作 error 使用。
func (e *AppError) Error() string {

	// 如果内部存在原始错误，
	// 返回“业务错误 + 原始错误”。
	if e.Err != nil {
		return fmt.Sprintf(
			"%s: %v",
			e.Message,
			e.Err,
		)
	}

	return e.Message
}

// Unwrap 返回底层原始错误。
//
// 这样 Go 的 errors.Is / errors.As
// 才可以继续向下查找原始错误。
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建一个业务错误。
func New(
	code int,
	httpStatus int,
	message string,
) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}
}

// Wrap 将一个底层 error 包装成业务错误。
func Wrap(
	code int,
	httpStatus int,
	message string,
	err error,
) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
		Err:        err,
	}
}
