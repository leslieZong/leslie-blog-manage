package handler

import (
	"github.com/gin-gonic/gin"

	"leslie-blog-server/internal/modules/category/dto"
	"leslie-blog-server/internal/modules/category/service"
	"leslie-blog-server/internal/response"
)

// PublicCategoryHandler 负责 Blog 前台分类接口。
//
// 和 CategoryHandler 最大区别：
//
// CategoryHandler
//
//	↓
//
// Admin
//
//	↓
//
// 需要 JWT + Casbin
//
// PublicCategoryHandler
//
//	↓
//
// Blog
//
//	↓
//
// 不需要登录
type PublicCategoryHandler struct {
	service service.CategoryService
}

// NewPublicCategoryHandler 创建 PublicCategoryHandler。
func NewPublicCategoryHandler(
	categoryService service.CategoryService,
) *PublicCategoryHandler {

	return &PublicCategoryHandler{
		service: categoryService,
	}
}

// List 获取前台分类列表。
//
// GET /api/v1/categories
func (h *PublicCategoryHandler) List(c *gin.Context) {

	categories, err := h.service.List(
		c.Request.Context(),
	)

	if err != nil {

		response.AppError(
			c,
			err,
		)

		return
	}

	// Public API 不应该直接返回完整 Category Model。
	//
	// 这里只返回：
	//
	// ID
	// Name
	// Slug
	result := make(
		[]*dto.SimpleCategoryResponse,
		0,
		len(categories),
	)

	for _, category := range categories {

		// Public 页面只展示启用分类。
		if category.Status != 1 {
			continue
		}

		result = append(
			result,
			dto.FromSimpleModel(category),
		)
	}

	response.Success(c, result)
}
