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

		response.AppError(
			c,
			err,
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

// GetRoles 获取指定用户的角色列表
//
// HTTP:
//
//	GET /api/admin/v1/users/:id/roles
//
// 业务流程：
//  1. 从 URL 获取用户 ID
//  2. 获取 HTTP Context
//  3. 调用 UserService 查询角色
//  4. 返回角色列表
func (h *UserHandler) GetRoles(c *gin.Context) {
	fmt.Println("GetRoles")
	// ---------------------------------------------------------
	// 1. 获取 URL 参数
	//
	// 对于：
	//
	// GET /users/01KABC/roles
	//
	// c.Param("id") 得到：
	//
	// 01KABC
	// ---------------------------------------------------------
	userID := c.Param("id")

	// ---------------------------------------------------------
	// 2. 获取 Go 标准 context.Context
	//
	// c.Request.Context() 是当前 HTTP 请求对应的 Context。
	//
	// Repository / Service 层应该继续使用这个 Context，
	// 这样请求被取消、客户端断开连接等情况能够向下传递。
	// ---------------------------------------------------------
	ctx := c.Request.Context()

	// ---------------------------------------------------------
	// 3. 调用 Service
	//
	// Handler 不应该自己查询数据库。
	//
	// 这里真正的业务逻辑交给 UserService：
	//
	// UserService
	//    ↓
	// UserRepository
	//    ↓
	// Casbin
	//    ↓
	// RoleRepository
	// ---------------------------------------------------------
	roles, err := h.service.GetRoles(ctx, userID)
	if err != nil {
		// 这里沿用你当前项目已有的统一错误处理方式。
		//
		// 例如你现有 Handler 中如果使用：
		//
		// h.handleError(c, err)
		//
		// 就继续使用这个方法。
		//
		// 不要在这里重新实现一套错误处理。
		response.AppError(
			c,
			err,
		)
		return
	}

	// ---------------------------------------------------------
	// 4. 返回成功响应
	//
	// 同样沿用项目当前统一 Response 方法。
	// ---------------------------------------------------------
	response.Success(c, roles)
}

// UpdateRoles 修改指定用户的角色
//
// HTTP:
//
//	PUT /api/admin/v1/users/:id/roles
//
// Request:
//
//	{
//	    "roles": ["admin", "editor"]
//	}
func (h *UserHandler) UpdateRoles(c *gin.Context) {
	// ---------------------------------------------------------
	// 1. 获取用户 ID
	// ---------------------------------------------------------
	userID := c.Param("id")

	// ---------------------------------------------------------
	// 2. 定义请求 DTO
	//
	// DTO 的作用：
	//
	// HTTP JSON
	//    ↓
	// UpdateUserRolesRequest
	//    ↓
	// Service
	//
	// Handler 不应该直接拿 gin.Context 的数据往 Service
	// 里面传。
	// ---------------------------------------------------------
	var req dto.UpdateUserRolesRequest

	// ---------------------------------------------------------
	// 3. 解析 JSON
	//
	// 前端：
	//
	// {
	//     "roles": ["admin", "editor"]
	// }
	//
	// 会被解析成：
	//
	// req.Roles
	//
	// []string{"admin", "editor"}
	// ---------------------------------------------------------
	if err := c.ShouldBindJSON(&req); err != nil {
		// 这里继续使用你项目现有的参数错误处理方式。
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)
		return
	}

	// ---------------------------------------------------------
	// 4. 获取请求 Context
	// ---------------------------------------------------------
	ctx := c.Request.Context()

	// ---------------------------------------------------------
	// 5. 调用 Service
	//
	// Handler 到这里就应该停止处理业务。
	//
	// 它不需要知道：
	//
	// - Casbin 怎么保存角色
	// - g 表示什么
	// - roleRepo 怎么查询
	// - 删除旧角色还是增加新角色
	//
	// 这些全部由 UserService 负责。
	// ---------------------------------------------------------
	if err := h.service.UpdateRoles(ctx, userID, req.Roles); err != nil {
		response.AppError(
			c,
			err,
		)
		return
	}

	// ---------------------------------------------------------
	// 6. 返回成功
	//
	// PUT 修改成功以后，我们可以返回一个简单的成功响应。
	// ---------------------------------------------------------
	response.Success(c, nil)
}
