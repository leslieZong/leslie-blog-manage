package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 是 HTTP 请求日志 Middleware。
//
// 每当客户端请求我们的 API：
//
// 请求进入
//
//	↓
//
// Logger
//
//	↓
//
// Handler
//
//	↓
//
// 请求结束
//
// Logger 都可以记录相关信息。
func Logger() gin.HandlerFunc {

	// 返回一个 Gin Middleware 函数。
	return func(c *gin.Context) {

		// 记录请求开始时间。
		start := time.Now()

		// 继续执行后面的 Middleware / Handler。
		//
		// 这是 Middleware 最重要的一行。
		c.Next()

		// 请求执行结束。
		end := time.Now()

		// 计算请求耗时。
		latency := end.Sub(start)

		// 获取 HTTP 状态码。
		statusCode := c.Writer.Status()

		// 获取请求方法。
		method := c.Request.Method

		// 获取请求路径。
		path := c.Request.URL.Path

		// 输出日志。
		log.Printf(
			"method=%s path=%s status=%d latency=%s",
			method,
			path,
			statusCode,
			latency,
		)
	}
}
