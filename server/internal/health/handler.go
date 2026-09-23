package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 负责健康检查 HTTP 接口。
type Handler struct{}

// NewHandler 创建 Health Handler。
func NewHandler() *Handler {
	return &Handler{}
}

// Health 返回应用基础存活状态。
//
// 这个接口不检查 MySQL、Redis 等外部依赖。
// 只要 Go HTTP 服务本身能够正常响应，
// 就认为应用进程是 Alive 的。
func (h *Handler) Health(c *gin.Context) {

	c.JSON(
		http.StatusOK,
		Response{
			Status: StatusOK,
		},
	)
}

type ReadyHandler struct {
	checkers []Checker
}

func NewReadyHandler(
	checkers ...Checker,
) *ReadyHandler {

	return &ReadyHandler{
		checkers: checkers,
	}
}

func (h *ReadyHandler) Ready(
	c *gin.Context,
) {

	ctx := c.Request.Context()

	for _, checker := range h.checkers {

		if err := checker.Check(ctx); err != nil {

			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"status": "not_ready",
				},
			)

			return
		}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ready",
		},
	)
}
