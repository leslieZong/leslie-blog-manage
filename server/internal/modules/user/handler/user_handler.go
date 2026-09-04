package handler

import (
	"net/http"

	"leslie-blog-server/internal/modules/user/dto"
	"leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

// GetByID 获取用户详情。
func (h *UserHandler) GetByID(
	c *gin.Context,
) {

	// ==================================================
	// 1. 从 URL 获取 ID
	// ==================================================

	id := c.Param("id")

	// ==================================================
	// 2. 调用 Service
	// ==================================================

	user, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	// ==================================================
	// 3. 统一处理错误
	// ==================================================

	if err != nil {
		response.AppError(c, err)
		return
	}

	// ==================================================
	// 4. Model → DTO
	// ==================================================

	userResponse := dto.FromUser(user)

	// ==================================================
	// 5. 返回成功结果
	// ==================================================

	response.Success(
		c,
		userResponse,
	)

	_ = http.StatusOK
}
