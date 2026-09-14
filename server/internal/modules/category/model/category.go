package model

import "time"

// Category
//
// Category 是博客文章的分类。
//
// 数据库表：categories
type Category struct {
	// ID
	//
	// 使用 ULID 作为主键。
	ID string `gorm:"column:id;primaryKey"`

	// 分类名称
	//
	// 例如：
	// 前端
	// 后端
	// AI
	Name string `gorm:"column:name;size:100;uniqueIndex;not null"`

	// URL 标识
	//
	// 例如：
	// frontend
	// backend
	// ai
	Slug string `gorm:"column:slug;size:100;uniqueIndex;not null"`

	// 分类描述
	Description *string `gorm:"column:description;size:255"`

	// 排序
	Sort int `gorm:"column:sort;not null;default:0"`

	// 状态
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `gorm:"column:status;not null;default:1"`

	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at"`

	// 更新时间
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// 删除时间
	//
	// 使用软删除。
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

// TableName
//
// 显式指定数据库表名。
func (Category) TableName() string {
	return "categories"
}
