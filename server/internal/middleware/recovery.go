package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Recovery 用来捕获程序运行过程中出现的 panic。
//
// 如果某个 Handler 出现 panic，
// 我们不希望整个 HTTP Server 直接崩溃。
//
// 所以需要 Recovery Middleware。
func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {

		// defer 表示：
		//
		// 不管后面的代码正常结束，
		// 还是发生 panic，
		// 都会执行这里。
		defer func() {

			// recover() 可以捕获 panic。
			if err := recover(); err != nil {

				// 输出错误日志。
				log.Printf(
					"panic recovered: %v",
					err,
				)

				// 返回 500。
				c.JSON(500, gin.H{
					"code":    500,
					"message": "internal server error",
					"data":    nil,
				})

				// 停止当前请求。
				c.Abort()
			}
		}()

		// 继续执行后面的 Handler。
		c.Next()
	}
}
