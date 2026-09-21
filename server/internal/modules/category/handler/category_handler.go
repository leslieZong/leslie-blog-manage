package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/category/dto"
	"leslie-blog-server/internal/modules/category/repository"
	"leslie-blog-server/internal/modules/category/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/pkg/utils"
	"leslie-blog-server/internal/response"
)

// CategoryHandler 是 Category HTTP Handler。
//
// Handler 的职责：
//
// 1. 接收 HTTP 请求
// 2. 解析参数
// 3. 调用 Service
// 4. 将结果转换成 DTO
// 5. 返回 HTTP 响应
//
// Handler 不应该：
//
// ❌ 直接操作数据库
// ❌ 编写复杂业务规则
// ❌ 调用 GORM
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler 创建 CategoryHandler。
func NewCategoryHandler(
	categoryService service.CategoryService,
) *CategoryHandler {

	return &CategoryHandler{
		service: categoryService,
	}
}

// Create 创建分类。
//
// POST /api/admin/v1/categories
func (h *CategoryHandler) Create(c *gin.Context) {

	var req dto.CreateCategoryRequest

	// ShouldBindJSON 会将 HTTP JSON 请求体
	// 转换成 Go struct。
	//
	// 例如：
	//
	// {
	//   "name": "前端",
	//   "slug": "frontend"
	// }
	//
	// 会转换成：
	//
	// req.Name = "前端"
	// req.Slug = "frontend"
	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid parameters",
		)

		return
	}

	// Handler 不自己判断：
	//
	// name 是否为空
	// slug 是否重复
	// slug 是否合法
	//
	// 这些属于业务规则。
	//
	// 所以统一交给 Service。
	category, err := h.service.Create(
		c.Request.Context(),
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.Slug),
		req.Description,
		req.Sort,
	)

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	// Model → DTO。
	result := dto.FromModel(category)

	response.Success(c, result)
}

// List 获取分类列表。
//
// GET /api/admin/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {

	ctx := c.Request.Context()
	statusValue := c.Query("status")
	status, err := utils.ParseStatus(statusValue)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid status",
		)
		return
	}

	query := repository.CategoryListQuery{
		Params:  pagination.Parse(c),
		Keyword: c.Query("keyword"),
		Status:  status,
	}

	categories, total, err := h.service.ListPage(
		ctx,
		query,
	)

	if err != nil {
		response.AppError(
			c,
			err,
		)

		return
	}

	// 将 Model 列表转换成 DTO 列表。
	res := dto.FromModels(categories)

	// 第四步：
	// 构造统一分页结果。
	result := pagination.NewResult(
		res,
		query.Params,
		total,
	)
	response.Success(c, result)
}

// GetByID 获取单个分类。
//
// GET /api/admin/v1/categories/:id
func (h *CategoryHandler) GetByID(c *gin.Context) {

	id := strings.TrimSpace(
		c.Param("id"),
	)

	if id == "" {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"category id is required",
		)

		return
	}

	category, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		response.AppError(
			c,
			err,
		)
		return
	}

	response.Success(
		c,
		dto.FromModel(category),
	)
}

// Update 修改分类。
//
// PUT /api/admin/v1/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {

	id := strings.TrimSpace(
		c.Param("id"),
	)

	if id == "" {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"category id is required",
		)

		return
	}

	var req dto.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid parameters",
		)

		return
	}

	category, err := h.service.Update(
		c.Request.Context(),
		id,
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.Slug),
		req.Description,
		req.Sort,
		req.Status,
	)

	if err != nil {
		response.AppError(
			c,
			err,
		)

		return
	}

	response.Success(
		c,
		dto.FromModel(category),
	)
}

// Delete 删除分类。
//
// DELETE /api/admin/v1/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {

	id := strings.TrimSpace(
		c.Param("id"),
	)

	if id == "" {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"category id is required",
		)

		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		id,
	); err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	response.Success(c, nil)
}
