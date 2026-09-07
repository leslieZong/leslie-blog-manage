package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/user/dto"
	"leslie-blog-server/internal/modules/user/repository"
	"leslie-blog-server/internal/modules/user/service"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/pkg/validator"
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
		if errors.Is(err, appErrors.ErrUserNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				appErrors.ErrNotFound,
				err.Error(),
			)
			return
		}

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

// List 获取用户列表。
//
// GET /api/admin/v1/users
func (h *UserHandler) List(c *gin.Context) {

	// =========================================================
	// 第一步：解析分页参数
	// =========================================================

	pageParams := pagination.Parse(c)

	// =========================================================
	// 第二步：获取搜索条件
	// =========================================================

	keyword := c.Query("keyword")

	// =========================================================
	// 第三步：获取 status
	// =========================================================

	var status *int8

	statusString := c.Query("status")

	if statusString != "" {

		value, err := strconv.ParseInt(
			statusString,
			10,
			8,
		)

		if err != nil {

			response.Error(
				c,
				http.StatusBadRequest,
				appErrors.ErrInvalidParams,
				"invalid status",
			)

			return
		}

		statusValue := int8(value)

		status = &statusValue
	}

	// =========================================================
	// 第四步：调用 Service
	// =========================================================

	users, total, err := h.service.List(
		c.Request.Context(),
		repository.UserListParams{
			Keyword: keyword,
			Status:  status,
			Offset:  pageParams.Offset(),
			Limit:   pageParams.PageSize,
		},
	)

	if err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"get user list failed",
		)

		return
	}

	// =========================================================
	// 第五步：Model → DTO
	// =========================================================

	items := make(
		[]*dto.UserResponse,
		0,
		len(users),
	)

	for _, user := range users {

		items = append(
			items,
			dto.FromUser(user),
		)
	}

	// =========================================================
	// 第六步：构建分页结果
	// =========================================================

	result := pagination.Result[*dto.UserResponse]{
		Items:    items,
		Total:    total,
		Page:     pageParams.Page,
		PageSize: pageParams.PageSize,
	}

	response.Success(c, result)
}

func (h *UserHandler) Create(c *gin.Context) {

	var req dto.CreateUserRequest

	// =========================================================
	// 1. JSON → DTO
	// =========================================================

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	// =========================================================
	// 2. DTO 参数校验
	// =========================================================

	if err := validator.Validate.Struct(req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			validator.FormatErrors(err),
		)

		return
	}

	// =========================================================
	// 3. Service
	// =========================================================

	user, err := h.service.CreateFromRequest(
		c.Request.Context(),
		&req,
	)
	fmt.Println("============")
	fmt.Println(err)
	if err != nil {

		// 用户名重复。
		if errors.Is(
			err,
			appErrors.ErrUsernameExists,
		) {

			response.Error(
				c,
				http.StatusConflict,
				appErrors.ErrConflict,
				"username already exists",
			)

			return
		}

		// 其他未知错误。
		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"internal server error",
		)

		return
	}

	// =========================================================
	// 4. Model → Response DTO
	// =========================================================

	response.Success(
		c,
		dto.FromUser(user),
	)
}

// Update 修改用户。
//
// PUT /api/admin/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {

	// =========================================================
	// 1. 获取 URL Path 参数
	// =========================================================

	id := c.Param("id")

	if id == "" {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"id is required",
		)

		return
	}

	// =========================================================
	// 2. JSON → DTO
	// =========================================================

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)

		return
	}

	// =========================================================
	// 3. 参数校验
	// =========================================================

	if err := validator.Validate.Struct(req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			validator.FormatErrors(err),
		)

		return
	}

	// =========================================================
	// 4. Service
	// =========================================================

	user, err := h.service.Update(
		c.Request.Context(),
		id,
		&req,
	)

	if err != nil {

		// 用户不存在。
		if errors.Is(
			err,
			appErrors.ErrUserNotFound,
		) {

			response.Error(
				c,
				http.StatusNotFound,
				appErrors.ErrNotFound,
				"user not found",
			)

			return
		}

		// 其他错误。
		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"internal server error",
		)

		return
	}

	// =========================================================
	// 5. Model → DTO
	// =========================================================

	response.Success(
		c,
		dto.FromUser(user),
	)
}

func (h *UserHandler) Delete(c *gin.Context) {

	// 获取 URL Path 参数。
	//
	// DELETE /users/:id
	//
	// 例如：
	//
	// DELETE /users/01KABC...
	//
	// 那么：
	//
	// c.Param("id")
	// 就是：
	//
	// 01KABC...
	targetUserID := c.Param("id")

	// 获取当前登录用户 ID。
	//
	// 这个 ID 来自 JWT Middleware。
	operatorUserID := auth.GetUserID(c)

	// 调用 Service 执行删除。
	//
	// Handler 不直接操作数据库。
	err := h.service.Delete(
		c.Request.Context(),
		operatorUserID,
		targetUserID,
	)

	if err != nil {

		// 不能删除自己。
		if errors.Is(
			err,
			appErrors.ErrCannotDeleteSelf,
		) {
			response.Error(
				c,
				http.StatusForbidden,
				appErrors.ErrForbidden,
				err.Error(),
			)
			return
		}

		// 用户不存在。
		if errors.Is(
			err,
			appErrors.ErrUserNotFound,
		) {
			response.Error(
				c,
				http.StatusNotFound,
				appErrors.ErrNotFound,
				err.Error(),
			)
			return
		}

		// 其他错误。
		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"delete user failed",
		)
		return
	}

	// 删除成功。
	response.Success(c, nil)
}
