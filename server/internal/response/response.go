package response

import (
	"errors"
	"net/http"

	appErrors "leslie-blog-server/internal/errors"

	"github.com/gin-gonic/gin"
)

// Response 是 Leslie Blog API 统一响应结构。
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 返回普通错误响应。
func Error(
	c *gin.Context,
	httpStatus int,
	code int,
	message string,
) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// AppError 处理统一业务错误。
//
// Handler 可以直接：
//
// response.AppError(c, err)
//
// 不需要自己判断 HTTP Status 和业务 Code。
func AppError(c *gin.Context, err error) {

	// --------------------------------------------------------
	// 如果是我们定义的 AppError
	// --------------------------------------------------------

	var appErr *appErrors.AppError

	if errors.As(err, &appErr) {

		c.JSON(
			appErr.HTTPStatus,
			Response{
				Code:    appErr.Code,
				Message: appErr.Message,
				Data:    nil,
			},
		)

		return
	}

	// --------------------------------------------------------
	// 如果不是我们定义的错误，
	// 就统一当作服务器内部错误。
	// --------------------------------------------------------

	c.JSON(
		http.StatusInternalServerError,
		Response{
			Code:    appErrors.ErrInternalServer,
			Message: "internal server error",
			Data:    nil,
		},
	)
}
