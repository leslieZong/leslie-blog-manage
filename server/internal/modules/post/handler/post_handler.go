package handler

import (
	"leslie-blog-server/internal/modules/post/dto"
	"leslie-blog-server/internal/modules/post/service"
	"leslie-blog-server/internal/pkg/auth"
	"leslie-blog-server/internal/response"

	"github.com/gin-gonic/gin"
)

// PostHandler 文章 HTTP Handler。
//
// Handler 负责 HTTP 层：
//
//	HTTP Request
//	    ↓
//	Handler
//	    ↓
//	Service
//
// 它不直接操作数据库。
type PostHandler struct {
	// service 是文章业务层。
	service service.PostService
}

// NewPostHandler 创建 PostHandler。
//
// 依赖通过参数传入，而不是在 Handler 内部自己创建。
func NewPostHandler(
	postService service.PostService,
) *PostHandler {

	return &PostHandler{
		service: postService,
	}
}

func (h *PostHandler) Create(c *gin.Context) {

	// -------------------------------------------------------
	// 1. 创建请求 DTO
	// -------------------------------------------------------

	var req dto.CreatePostRequest

	// -------------------------------------------------------
	// 2. 解析 JSON
	// -------------------------------------------------------

	if err := c.ShouldBindJSON(&req); err != nil {

		// 这里应该使用你项目现有的统一错误响应方法。
		//
		// 不建议重新创建一套错误体系。
		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})

		return
	}

	// -------------------------------------------------------
	// 3. 获取当前登录用户
	// -------------------------------------------------------

	userID := auth.GetUserID(c)

	if userID == "" {

		c.JSON(401, gin.H{
			"message": "未登录",
		})

		return
	}

	// -------------------------------------------------------
	// 4. 获取 HTTP Context
	// -------------------------------------------------------

	ctx := c.Request.Context()

	// -------------------------------------------------------
	// 5. 调用 Service
	// -------------------------------------------------------

	post, err := h.service.Create(
		ctx,
		userID,
		req.Title,
		req.Slug,
		req.Summary,
		req.Content,
		req.Cover,
	)

	if err != nil {

		// 暂时使用简单响应。
		//
		// 后面需要统一接入项目已有错误码体系。
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	// -------------------------------------------------------
	// 6. Model → Response DTO
	// -------------------------------------------------------

	res := dto.FromModel(post)

	// -------------------------------------------------------
	// 7. 返回响应
	// -------------------------------------------------------

	response.Success(c, res)
}

// GetByID 获取文章详情。
func (h *PostHandler) GetByID(c *gin.Context) {

	// 从 URL 中获取文章 ID。
	//
	// 例如：
	// GET /posts/01KABC...
	//
	// c.Param("id")
	// 得到：
	// 01KABC...
	id := c.Param("id")

	// 获取请求 Context。
	ctx := c.Request.Context()

	// 调用 Service。
	post, err := h.service.GetByID(ctx, id)

	if err != nil {

		c.JSON(404, gin.H{
			"message": "文章不存在",
		})

		return
	}

	// Model → DTO。
	res := dto.FromModel(post)

	// 返回 JSON。
	response.Success(c, res)
}

// List 获取文章列表。
func (h *PostHandler) List(c *gin.Context) {

	ctx := c.Request.Context()

	posts, err := h.service.List(ctx)

	if err != nil {

		c.JSON(500, gin.H{
			"message": "获取文章列表失败",
		})

		return
	}

	res := dto.FromModelList(posts)

	response.Success(c, res)
}

// Update 更新文章。
func (h *PostHandler) Update(c *gin.Context) {

	// -------------------------------------------------------
	// 1. 获取文章 ID
	// -------------------------------------------------------

	id := c.Param("id")

	// -------------------------------------------------------
	// 2. 解析请求 JSON
	// -------------------------------------------------------

	var req dto.UpdatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{
			"message": "请求参数错误",
		})

		return
	}

	// -------------------------------------------------------
	// 3. 调用 Service
	// -------------------------------------------------------

	ctx := c.Request.Context()

	post, err := h.service.Update(
		ctx,
		id,
		req.Title,
		req.Slug,
		req.Summary,
		req.Content,
		req.Cover,
	)

	if err != nil {

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	// -------------------------------------------------------
	// 4. Model → DTO
	// -------------------------------------------------------

	res := dto.FromModel(post)

	// -------------------------------------------------------
	// 5. 返回
	// -------------------------------------------------------

	response.Success(c, res)
}

// Delete 删除文章。
//
// 当前是软删除。
func (h *PostHandler) Delete(c *gin.Context) {

	id := c.Param("id")

	ctx := c.Request.Context()

	if err := h.service.Delete(ctx, id); err != nil {

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	response.Success(c, nil)
}

// Publish 发布文章。
func (h *PostHandler) Publish(c *gin.Context) {

	id := c.Param("id")

	ctx := c.Request.Context()

	post, err := h.service.Publish(ctx, id)

	if err != nil {

		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	res := dto.FromModel(post)

	response.Success(c, res)
}

// GetBySlug 根据 slug 获取文章。
func (h *PostHandler) GetBySlug(c *gin.Context) {

	slug := c.Param("slug")

	ctx := c.Request.Context()

	post, err := h.service.GetBySlug(ctx, slug)

	if err != nil {

		c.JSON(404, gin.H{
			"message": "文章不存在",
		})

		return
	}

	res := dto.FromModel(post)

	response.Success(c, res)
}
