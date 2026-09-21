package handler

import (
	"net/http"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/tag/dto"
	"leslie-blog-server/internal/modules/tag/repository"
	"leslie-blog-server/internal/modules/tag/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/pkg/utils"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// TagHandler 负责 Tag HTTP 请求。
//
// Handler 的职责非常明确：
//
// 1. 获取 HTTP 参数
// 2. 参数绑定
// 3. 调用 Service
// 4. 处理错误
// 5. Model → DTO
// 6. 返回 HTTP Response
//
// Handler 不负责 SQL。
// Handler 不负责核心业务。
type TagHandler struct {
	service service.TagService
}

// NewTagHandler 创建 TagHandler。
func NewTagHandler(
	tagService service.TagService,
) *TagHandler {

	return &TagHandler{
		service: tagService,
	}
}

// =========================================================
// Create
// =========================================================
//
// POST /tags
//
// 创建 Tag。
// =========================================================

func (h *TagHandler) Create(c *gin.Context) {

	// -----------------------------------------------------
	// 第一步：绑定 JSON
	// -----------------------------------------------------

	var req dto.CreateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid parameters",
		)

		return
	}

	// -----------------------------------------------------
	// 第二步：调用 Service
	// -----------------------------------------------------

	tag, err := h.service.Create(
		c.Request.Context(),
		req.Name,
		req.Slug,
		req.Description,
	)

	// -----------------------------------------------------
	// 第三步：处理错误
	// -----------------------------------------------------

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	// -----------------------------------------------------
	// 第四步：Model → DTO
	// -----------------------------------------------------

	response.Success(
		c,
		dto.FromModel(tag),
	)
}

// =========================================================
// GetByID
// =========================================================
//
// GET /tags/:id
// =========================================================

func (h *TagHandler) GetByID(c *gin.Context) {

	// 获取 URL 参数。
	id := c.Param("id")

	// 调用 Service。
	tag, err := h.service.GetByID(
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

	// 返回 DTO。
	response.Success(
		c,
		dto.FromModel(tag),
	)
}

// =========================================================
// List
// =========================================================
//
// GET /tags
// =========================================================

func (h *TagHandler) ListPage(c *gin.Context) {

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
	query := repository.TagListQuery{
		Params:  pagination.Parse(c),
		Keyword: c.Query("keyword"),
		Status:  status,
	}

	tags, total, err := h.service.ListPage(
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

	// Model [] → DTO []。
	res := make(
		[]*dto.TagResponse,
		0,
		len(tags),
	)

	for _, tag := range tags {

		res = append(
			res,
			dto.FromModel(tag),
		)
	}
	result := pagination.NewResult(
		res,
		query.Params,
		total,
	)

	response.Success(
		c,
		result,
	)
}

// =========================================================
// Update
// =========================================================
//
// PUT /tags/:id
// =========================================================

func (h *TagHandler) Update(c *gin.Context) {

	id := c.Param("id")

	// 绑定请求体。
	var req dto.UpdateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			appErrors.ErrInvalidParams,
			"invalid parameters",
		)

		return
	}

	// 调用 Service。
	tag, err := h.service.Update(
		c.Request.Context(),
		id,
		req.Name,
		req.Slug,
		req.Description,
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
		dto.FromModel(tag),
	)
}

// =========================================================
// Delete
// =========================================================
//
// DELETE /tags/:id
// =========================================================

func (h *TagHandler) Delete(c *gin.Context) {

	id := c.Param("id")

	err := h.service.Delete(
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
		nil,
	)
}
