package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/post/dto"
	"leslie-blog-server/internal/modules/post/repository"
	"leslie-blog-server/internal/modules/post/service"
	"leslie-blog-server/internal/pkg/pagination"
	"leslie-blog-server/internal/response"
)

// PublicPostHandler
//
// 专门处理博客前台 Post API。
//
// 为什么不直接复用 PostHandler？
//
// 因为 Admin 和 Public 的职责不同：
//
// AdminPostHandler
//
//	→ 管理文章
//
// PublicPostHandler
//
//	→ 展示文章
//
// 分开以后，权限边界更加清晰。
type PublicPostHandler struct {
	service service.PostService
}

// NewPublicPostHandler
//
// 创建 Public Post Handler。
func NewPublicPostHandler(
	service service.PostService,
) *PublicPostHandler {

	return &PublicPostHandler{
		service: service,
	}
}

// List
//
// GET /api/v1/posts
//
// 获取已经发布的文章列表。
// ListPublished 获取已发布文章列表。
func (h *PublicPostHandler) ListPublished(c *gin.Context) {

	ctx := c.Request.Context()
	query := repository.PostListQuery{
		Params:     pagination.Parse(c),
		CategoryID: c.Query("categoryId"),
		Status:     c.Query("status"),
	}

	posts, total, err := h.service.ListPublished(ctx, query)

	if err != nil {

		response.AppError(c, err)

		return
	}

	res := dto.FromPublicModelList(posts)
	// 第四步：
	// 构造统一分页结果。
	result := pagination.NewResult(
		res,
		query.Params,
		total,
	)

	response.Success(c, result)
}

// GetByID
//
// GET /api/v1/posts/:id
func (h *PublicPostHandler) GetByID(c *gin.Context) {

	id := c.Param("id")

	post, err := h.service.GetPublicByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		response.AppError(c, err)

		return
	}

	response.Success(
		c,
		dto.FromPublicModel(post),
	)
}

// GetBySlug
//
// GET /api/v1/posts/slug/:slug
func (h *PublicPostHandler) GetBySlug(c *gin.Context) {

	slug := c.Param("slug")

	post, err := h.service.GetPublicBySlug(
		c.Request.Context(),
		slug,
	)

	if err != nil {
		response.AppError(c, err)

		return
	}

	response.Success(
		c,
		dto.FromPublicModel(post),
	)
}

// IncrementViewCount
//
// POST /api/v1/posts/:id/view
//
// 这里暂时单独设计一个浏览量接口。
//
// 后续可以进一步优化成：
// 访问文章详情时自动增加浏览量。
func (h *PublicPostHandler) IncrementViewCount(c *gin.Context) {

	id := c.Param("id")

	// 先确认文章确实存在且已经公开。
	_, err := h.service.GetPublicByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		response.AppError(c, err)
		return
	}

	// 增加阅读量
	if err := h.service.IncrementViewCount(
		c.Request.Context(),
		id,
	); err != nil {

		response.AppError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
