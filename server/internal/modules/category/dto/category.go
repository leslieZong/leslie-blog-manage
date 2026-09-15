package dto

import (
	"time"

	"leslie-blog-server/internal/modules/category/model"
)

// CreateCategoryRequest 是创建分类时的请求参数。
//
// 对应：
//
// POST /api/admin/v1/categories
//
// 前端 JSON：
//
//	{
//	  "name": "前端",
//	  "slug": "frontend",
//	  "description": "前端开发相关内容",
//	  "sort": 1
//	}
type CreateCategoryRequest struct {
	// Name 分类名称。
	Name string `json:"name"`

	// Slug 分类 URL 标识。
	Slug string `json:"slug"`

	// Description 分类描述。
	Description *string `json:"description"`

	// Sort 排序值。
	Sort int `json:"sort"`
}

// UpdateCategoryRequest 是修改分类时的请求参数。
//
// 对应：
//
// PUT /api/admin/v1/categories/:id
type UpdateCategoryRequest struct {
	// Name 分类名称。
	Name string `json:"name"`

	// Slug 分类 URL 标识。
	Slug string `json:"slug"`

	// Description 分类描述。
	Description *string `json:"description"`

	// Sort 排序值。
	Sort int `json:"sort"`

	// Status 分类状态。
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `json:"status"`
}

// CategoryResponse 是后台分类详情响应。
//
// 注意：
//
// 这里是 DTO，而不是 Model。
//
// 不应该直接把 GORM Model 返回给前端。
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	Sort        int       `json:"sort"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// FromModel 将数据库 Model 转换成 API DTO。
//
// 这个函数的作用非常简单：
//
// Model
//
//	↓
//
// # DTO
//
// 它可以避免 Handler 直接操作数据库模型。
func FromModel(category *model.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Sort:        category.Sort,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// FromModels 将多个 Category Model 转换成 DTO。
//
// 后台列表接口会用到这个函数。
func FromModels(
	categories []*model.Category,
) []*CategoryResponse {

	result := make([]*CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, FromModel(category))
	}

	return result
}
