package httpx

import (
	"github.com/gin-gonic/gin"
)

// Handler 是 Leslie Blog 自己定义的 Handler 类型。
//
// 和 Gin 原生 Handler 最大的区别：
//
// 我们允许 Handler 返回 error。
//
// 这样业务错误可以一路向上返回，
// 最后交给统一 Error Middleware 处理。
type Handler func(*gin.Context) error

// Adapt 将我们自己的 Handler 转换成 Gin Handler。
//
// Gin 需要：
//
//	func(*gin.Context)
//
// 而我们的 Handler 是：
//
//	func(*gin.Context) error
//
// Adapt 就负责把两者连接起来。
func Adapt(
	handler Handler,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 执行业务 Handler。
		err := handler(c)
		// 如果没有错误，
		// 什么都不做。
		if err == nil {
			return
		}

		// 将错误交给 Gin Context。
		//
		// 注意：
		// 这里暂时不直接返回 JSON。
		//
		// 真正的错误响应由后面的
		// Error Middleware 统一处理。
		_ = c.Error(err)
	}
}
