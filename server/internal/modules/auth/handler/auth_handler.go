package handler

import (
	"net/http"

	"leslie-blog-server/internal/modules/auth/dto"
	"leslie-blog-server/internal/modules/auth/service"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 是认证模块的 HTTP Handler。
//
// Handler 只负责 HTTP 层工作：
//
// 1. 接收请求
// 2. 解析 JSON
// 3. 调用 Service
// 4. 返回 JSON
//
// 它不应该自己查询数据库。
type AuthHandler struct {

	// service 是 AuthService。
	service service.AuthService
}

// NewAuthHandler 创建 AuthHandler。
func NewAuthHandler(
	authService service.AuthService,
) *AuthHandler {

	return &AuthHandler{
		service: authService,
	}
}

// Login 处理登录请求。
//
// POST /api/admin/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {

	// ----------------------------------------
	// 1. 创建请求 DTO
	// ----------------------------------------

	var req dto.LoginRequest

	// ----------------------------------------
	// 2. 解析 JSON Body
	// ----------------------------------------

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			40001,
			"invalid request body",
		)

		return
	}

	// ----------------------------------------
	// 3. 调用 Auth Service
	// ----------------------------------------

	result, err := h.service.Login(
		c.Request.Context(),
		&req,
	)

	if err != nil {

		// 目前先简单处理。
		//
		// 后面我们会建立统一业务错误类型，
		// 到时候这里可以更加精确地区分：
		//
		// 参数错误
		// 用户不存在
		// 用户禁用
		// 密码错误
		// 系统错误
		response.Error(
			c,
			http.StatusUnauthorized,
			40101,
			err.Error(),
		)

		return
	}

	// ----------------------------------------
	// 4. 返回登录成功结果
	// ----------------------------------------

	response.Success(
		c,
		result,
	)
}
