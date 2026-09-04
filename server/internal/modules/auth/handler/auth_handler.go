package handler

import (
	"net/http"

	appErrors "leslie-blog-server/internal/errors"
	authDTO "leslie-blog-server/internal/modules/auth/dto"
	"leslie-blog-server/internal/modules/auth/service"
	userDTO "leslie-blog-server/internal/modules/user/dto"
	userService "leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 负责 HTTP 层的认证相关请求。
//
// 注意：
//
// Handler 不直接操作数据库。
// Handler 只负责：
//
// HTTP Request
//
//	↓
//
// 参数解析
//
//	↓
//
// 调用 Service
//
//	↓
//
// HTTP Response
type AuthHandler struct {

	// authService 负责登录业务。
	authService service.AuthService

	// userService 负责用户相关业务。
	//
	// /me 接口会使用它查询当前用户。
	userService userService.UserService
}

// NewAuthHandler 创建 AuthHandler。
//
// 这里使用依赖注入：
//
// authService
// userService
//
//	↓
//
// AuthHandler
func NewAuthHandler(
	authService service.AuthService,
	userService userService.UserService,
) *AuthHandler {

	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Login 登录接口。
//
// POST /api/admin/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {

	var req authDTO.LoginRequest

	// 将 JSON 请求体绑定到 LoginRequest。
	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	// 调用 Service 执行真正的登录业务。
	result, err := h.authService.Login(
		c.Request.Context(),
		&req,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusUnauthorized,
			appErrors.ErrUnauthorized,
			err.Error(),
		)

		return
	}

	response.Success(c, result)
}

// Me 获取当前登录用户。
//
// GET /api/admin/v1/auth/me
//
// 注意：
//
// 这个接口必须经过 JWT Middleware。
//
// JWT Middleware 会提前把：
//
//	userID
//	username
//
// 放入 Gin Context。
func (h *AuthHandler) Me(c *gin.Context) {

	// =========================================================
	// 第一步：从 Context 获取当前用户 ID
	// =========================================================

	userID := auth.GetUserID(c)

	if userID == "" {

		response.Error(
			c,
			http.StatusUnauthorized,
			appErrors.ErrUnauthorized,
			"user identity not found",
		)

		return
	}

	// =========================================================
	// 第二步：调用 UserService 查询用户
	// =========================================================

	user, err := h.userService.GetByID(
		c.Request.Context(),
		userID,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"get current user failed",
		)

		return
	}

	// =========================================================
	// 第三步：Model → DTO
	// =========================================================

	userResponse := userDTO.FromUser(user)

	// =========================================================
	// 第四步：返回 JSON
	// =========================================================

	response.Success(
		c,
		userResponse,
	)
}
