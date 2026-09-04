package handler

import (
	"net/http"

	"leslie-blog-server/internal/modules/auth/dto"
	"leslie-blog-server/internal/modules/auth/service"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: authService,
	}
}

// Login 处理用户登录请求。
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	// ==================================================
	// 1. 解析 JSON
	// ==================================================

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			40001,
			"invalid request body",
		)

		return
	}

	// ==================================================
	// 2. 调用 Service
	// ==================================================

	result, err := h.service.Login(
		c.Request.Context(),
		&req,
	)

	// ==================================================
	// 3. 统一错误处理
	// ==================================================

	if err != nil {
		response.AppError(c, err)
		return
	}

	// ==================================================
	// 4. 返回登录结果
	// ==================================================

	response.Success(c, result)
}
