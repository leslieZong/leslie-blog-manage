package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 表示博客文章标签。
//
// 对应数据库：tags 表。
//
// 例如：
//
// ID:   01K...
// Name: Vue3
// Slug: vue3
//
// Tag 本身是独立的数据实体。
type Tag struct {
	// ID 是标签唯一 ID。
	//
	// 使用 ULID。
	ID string `gorm:"column:id;type:char(26);primaryKey"`

	// Name 是标签显示名称。
	//
	// 例如：
	//
	// Vue3
	// TypeScript
	// Go
	Name string `gorm:"column:name;size:100;not null;uniqueIndex"`

	// Slug 是标签 URL 标识。
	//
	// 例如：
	//
	// Vue3 -> vue3
	Slug string `gorm:"column:slug;size:100;not null;uniqueIndex"`

	// Description 是标签描述。
	//
	// 使用指针是因为数据库允许 NULL。
	Description *string `gorm:"column:description;size:255"`

	// Status：
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `gorm:"column:status;not null;default:1"`

	// CreatedAt 创建时间。
	CreatedAt time.Time `gorm:"column:created_at"`

	// UpdatedAt 更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// DeletedAt 软删除时间。
	//
	// nil 表示没有删除。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 明确告诉 GORM：
//
// Tag 对应 tags 表。
func (Tag) TableName() string {
	return "tags"
}
