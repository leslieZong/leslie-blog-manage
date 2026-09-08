package handler

import (
	"errors"
	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/role/dto"
	"leslie-blog-server/internal/modules/role/repository"
	roleService "leslie-blog-server/internal/modules/role/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/pkg/validator"
	"leslie-blog-server/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleHandler 负责处理 Role 相关 HTTP 请求。
type RoleHandler struct {
	roleService roleService.RoleService
}

// NewRoleHandler 创建 RoleHandler。
func NewRoleHandler(
	roleService roleService.RoleService,
) *RoleHandler {

	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) List(
	c *gin.Context,
) {

	// 解析分页参数。
	page := pagination.Parse(c)

	// 获取搜索关键词。
	keyword := c.DefaultQuery(
		"keyword",
		"",
	)

	// status 是可选参数。
	var status *int8

	statusString, exists := c.GetQuery("status")

	if exists {
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

	params := repository.RoleListParams{
		Keyword: keyword,
		Status:  status,
		Offset:  page.Offset(),
		Limit:   page.PageSize,
	}

	roles, total, err := h.roleService.List(
		c.Request.Context(),
		params,
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"get role list failed",
		)
		return
	}

	items := make(
		[]*dto.RoleResponse,
		0,
		len(roles),
	)

	for _, role := range roles {
		items = append(
			items,
			dto.FromRole(role),
		)
	}

	response.Success(
		c,
		pagination.Result[*dto.RoleResponse]{
			Items:    items,
			Total:    total,
			Page:     page.Page,
			PageSize: page.PageSize,
		},
	)
}

func (h *RoleHandler) Create(
	c *gin.Context,
) {

	var req dto.CreateRoleRequest

	// 解析 JSON Body。
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)
		return
	}

	// 参数校验。
	if err := validator.Validate.Struct(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			validator.FormatErrors(err),
		)
		return
	}

	role, err := h.roleService.Create(
		c.Request.Context(),
		&req,
	)

	if err != nil {
		if errors.Is(
			err,
			appErrors.ErrRoleNameExists,
		) {
			response.Error(
				c,
				http.StatusConflict,
				appErrors.ErrConflict,
				appErrors.ErrRoleNameExists.Error(),
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"create role failed",
		)
		return
	}

	response.Success(
		c,
		dto.FromRole(role),
	)
}

func (h *RoleHandler) Update(
	c *gin.Context,
) {

	id := c.Param("id")

	var req dto.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid request body",
		)
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			validator.FormatErrors(err),
		)
		return
	}

	role, err := h.roleService.Update(
		c.Request.Context(),
		id,
		&req,
	)

	if err != nil {
		if errors.Is(
			err,
			appErrors.ErrRoleNotFound,
		) {
			response.Error(
				c,
				http.StatusNotFound,
				appErrors.ErrNotFound,
				"role not found",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			appErrors.ErrInternalServer,
			"update role failed",
		)
		return
	}

	response.Success(
		c,
		dto.FromRole(role),
	)
}

func (h *RoleHandler) Delete(
	c *gin.Context,
) {

	id := c.Param("id")

	err := h.roleService.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			appErrors.ErrRoleNotFound,
		) {
			response.Error(
				c,
				http.StatusNotFound,
				appErrors.ErrNotFound,
				"role not found",
			)
			return
		}

		response.Error(
			c,
			http.StatusConflict,
			appErrors.ErrConflict,
			err.Error(),
		)
		return
	}

	response.Success(c, nil)
}
