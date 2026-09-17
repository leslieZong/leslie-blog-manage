package dto

import (
	"leslie-blog-server/internal/modules/tag/model"
)

// =========================================================
// CreateTagRequest
// =========================================================
//
// 创建 Tag 时，前端提交的数据。
//
// 注意：
//
// Request DTO ≠ Model
//
// Request DTO 表示：
// “前端允许提交什么”
//
// Model 表示：
// “数据库保存什么”
//
// 两者不要混在一起。
//
// =========================================================

type CreateTagRequest struct {
	// Tag 名称。
	//
	// 例如：
	//
	// Vue3
	// Go
	// TypeScript
	Name string `json:"name"`

	// Tag URL Slug。
	//
	// 例如：
	//
	// Vue3 → vue3
	Slug string `json:"slug"`

	// Tag 描述。
	Description *string `json:"description"`
}

// =========================================================
// UpdateTagRequest
// =========================================================
//
// 更新 Tag。
// =========================================================

type UpdateTagRequest struct {
	// Tag 名称。
	Name string `json:"name"`

	// Tag Slug。
	Slug string `json:"slug"`

	// Tag 描述。
	Description *string `json:"description"`

	// Tag 状态。
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `json:"status"`
}

// =========================================================
// TagResponse
// =========================================================
//
// 返回给前端的数据。
//
// 注意：
// 不直接把 Model 返回给前端。
//
// Model
//   ↓
// TagResponse
//   ↓
// JSON
// =========================================================

type TagResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Description *string `json:"description"`

	Status int8 `json:"status"`

	CreatedAt string `json:"createdAt"`

	UpdatedAt string `json:"updatedAt"`
}

// =========================================================
// FromModel
// =========================================================
//
// Model → Response DTO。
//
// 这一步是数据转换层。
// =========================================================

func FromModel(tag *model.Tag) *TagResponse {

	// 防止传入 nil。
	if tag == nil {
		return nil
	}

	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		Status:      tag.Status,
		CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
