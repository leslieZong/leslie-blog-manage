package dto

import "leslie-blog-server/internal/modules/category/model"

// SimpleCategoryResponse
//
// 用于 Post 返回分类信息。
//
// 这里只返回前台真正需要的字段。
type SimpleCategoryResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`
}

// FromSimpleModel
//
// 将 Category Model 转换为简单 DTO。
func FromSimpleModel(
	category *model.Category,
) *SimpleCategoryResponse {

	if category == nil {
		return nil
	}

	return &SimpleCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}
