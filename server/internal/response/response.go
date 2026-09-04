package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是 Leslie Blog 统一 API 响应结构。
//
// 所有接口最终都应该返回类似：
//
//	{
//	    "code": 0,
//	    "message": "success",
//	    "data": {}
//	}
//
// 这样前端可以统一处理。
type Response struct {

	// Code 是业务状态码。
	//
	// 0 表示成功。
	//
	// 非 0 表示业务错误。
	Code int `json:"code"`

	// Message 是响应消息。
	Message string `json:"message"`

	// Data 是业务数据。
	//
	// 使用 any 表示这里可以放任意类型的数据。
	//
	// 例如：
	//
	// User
	// Post
	// []Post
	// Pagination
	//
	Data any `json:"data"`
}

// Success 返回一个成功响应。
//
// 使用方式：
//
// response.Success(c, user)
//
// 最终：
//
//	{
//	    "code": 0,
//	    "message": "success",
//	    "data": user
//	}
func Success(c *gin.Context, data any) {

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 返回一个错误响应。
//
// 使用方式：
//
// response.Error(c, http.StatusBadRequest, 40001, "invalid parameters")
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
