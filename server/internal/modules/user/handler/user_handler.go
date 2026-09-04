package handler

import (
	"errors"
	"net/http"

	"leslie-blog-server/internal/modules/user/dto"
	"leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler 负责处理 User 相关 HTTP 请求。
//
// Handler 是 HTTP 世界和业务世界之间的边界。
//
// 它主要负责：
//
// 1. 从 HTTP 请求中获取参数
// 2. 调用 Service
// 3. 将结果转换成 DTO
// 4. 返回 HTTP Response
//
// Handler 不应该：
//
// ❌ 直接操作 GORM
// ❌ 直接查询 MySQL
// ❌ 编写复杂业务逻辑
type UserHandler struct {

	// service 是 UserService。
	//
	// Handler 不关心数据库怎么查询。
	// 它只需要调用业务服务。
	service service.UserService
}

// NewUserHandler 创建 UserHandler。
func NewUserHandler(
	userService service.UserService,
) *UserHandler {

	return &UserHandler{
		service: userService,
	}
}

// GetByID 获取指定用户。
//
// 对应：
//
// GET /api/admin/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {

	// --------------------------------------------------
	// 第一步：从 URL 获取 ID
	// --------------------------------------------------
	//
	// 路由：
	//
	// /users/:id
	//
	// 如果请求：
	//
	// /users/01KABC
	//
	// 那么：
	//
	// c.Param("id")
	//
	// 得到：
	//
	// 01KABC
	id := c.Param("id")

	// 如果 ID 为空，
	// 说明请求参数不正确。
	if id == "" {

		response.Error(
			c,
			http.StatusBadRequest,
			40001,
			"id is required",
		)

		return
	}

	// --------------------------------------------------
	// 第二步：调用 Service
	// --------------------------------------------------
	//
	// 注意：
	//
	// Handler 没有：
	//
	// db.First(...)
	//
	// 而是：
	//
	// service.GetByID(...)
	//
	// 这就是分层架构。
	user, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	// --------------------------------------------------
	// 第三步：处理错误
	// --------------------------------------------------
	if err != nil {

		// 用户不存在。
		if errors.Is(err, gorm.ErrRecordNotFound) {

			response.Error(
				c,
				http.StatusNotFound,
				40401,
				"user not found",
			)

			return
		}

		// 其他未知错误。
		response.Error(
			c,
			http.StatusInternalServerError,
			50001,
			"internal server error",
		)

		return
	}

	// --------------------------------------------------
	// 第四步：Model → DTO
	// --------------------------------------------------
	//
	// 注意：
	//
	// 不直接返回 user。
	//
	// 而是：
	//
	// User Model
	//     ↓
	// UserResponse
	userResponse := dto.FromUser(user)

	// --------------------------------------------------
	// 第五步：返回统一响应
	// --------------------------------------------------
	response.Success(
		c,
		userResponse,
	)
}
